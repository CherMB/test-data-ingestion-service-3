package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"sort"
	"time"

	"strings"

	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/calculi-corp/api/go/endpoint"
	pb "github.com/calculi-corp/api/go/vsm/report"
	client "github.com/calculi-corp/grpc-client"
	"github.com/calculi-corp/reports-service/constants"
	db "github.com/calculi-corp/reports-service/db"
	"github.com/calculi-corp/reports-service/exceptions"
	"github.com/calculi-corp/reports-service/helper"

	"github.com/calculi-corp/log"
)

func loadWidgetData(wb *WidgetBuilder, segment string, data map[string]json.RawMessage, fdata map[string]json.RawMessage, replacements map[string]any) error {

	switch segment {

	case "header":
		err := wb.setHeaders(data, fdata)
		if err != nil {
			if err == db.ErrNoDataFound {
				log.Debugf("No data found in set headers")
				return err
			} else {
				log.Errorf(err, "Error in set headers: ")
				return err
			}
		}

		return nil

	case "section":
		err := wb.setSection(data, fdata, replacements)
		if err != nil {
			if err == db.ErrNoDataFound {
				log.Debugf("No data found in set section")
				return err
			} else {
				log.Errorf(err, "Error in set section:")
				return err
			}
		}

		return nil

	case "footer":
		err := wb.setFooter(data, fdata)
		if err != nil {
			if err == db.ErrNoDataFound {
				log.Debugf("No data found in set footer")
				return err
			} else {
				log.Errorf(err, "Error in set footer:")
				return err
			}
		}

		return nil

	case "data":
		err := wb.setData(data, fdata)
		if err != nil {
			if err == db.ErrNoDataFound {
				log.Debugf("No data found in set data")
				return err
			} else {
				log.Errorf(err, "Error in set data:")
				return err
			}
		}

		return nil

	default:
		log.Error("No default action for loadWidgetData()", db.ErrInternalServer)
		return db.ErrInternalServer

	}

}

type WidgetTransform struct {
	segment string
	widget  pb.Widget
}

// Create widget with the reponse from Opensearch
// Apply spec for data transformation as per Widget definition
func CreateWidget(widgetId string, data map[string]json.RawMessage, replacements map[string]any, replacementsSpec map[string]any, fdata map[string]json.RawMessage) (*pb.Widget, error) {

	allReplacment := maps.Clone(replacements)
	maps.Copy(allReplacment, replacementsSpec)
	we, err := db.GetWidgetEntity(widgetId, allReplacment)
	if log.CheckErrorf(err, "could not get data from db in CreateWidget()") {
		return nil, err
	}

	//Call widget builder to construct the Widget
	wb := newWidgetBuilder()
	wb.setWidgetInfo(we.Widget)

	var segments = []string{"header", "section", "footer", "data"}
	for _, segment := range segments {
		log.Debugf("Processing %s for component id:%s widgetId:%s", segment, replacements["component"], widgetId)

		err = loadWidgetData(wb, segment, data, fdata, allReplacment)
		if err != nil {
			if err == db.ErrNoDataFound {
				log.Debugf("No data found for segment: %s,  componentID: %s and widgetID: %s", segment, replacements["component"], widgetId)
				return nil, err
			} else {
				log.Errorf(err, "Error in loadWidgetData() for segment: %s,  componentID: %s and widgetID: %s", segment, replacements["component"], widgetId)
				return nil, err
			}
		}

		log.Debugf("Process done %s for component id:%s widgetId:%s", segment, replacements["component"], widgetId)
	}
	startTime := time.Now()
	log.Debugf("Time took to set footer for widget %s : %v in milliseconds", widgetId, time.Since(startTime).Milliseconds())

	return &wb.widget, nil
}

// Gets Latest & data for Trend comparison from OpenSearch along with data from functions
func GetData(widgetId string, replacements map[string]any, baseDataReplacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (map[string]json.RawMessage, map[string]json.RawMessage, map[string]json.RawMessage, error) {
	// Get all queries from Widget Definition using Widget Id.
	// Apply the filters in queries and get data
	qm, fl, err := getWidgetQueries(widgetId, replacements, baseDataReplacements)
	if log.CheckErrorf(err, "failed to get widget queries for widget id %s", widgetId) {
		return nil, nil, nil, err // log.CheckErrorf returns nil if error
	}

	var resp map[string]json.RawMessage
	if qm != nil {
		startTime := time.Now()
		resp, err = helper.GetMultiQueryResponse(qm)
		//resp, err := helper.GetQueryResponse(qm)
		log.Infof("TIME TAKEN GetMultiQueryResponse call widget %s : %v in milliseconds", widgetId, time.Since(startTime).Milliseconds())
	}

	if log.CheckErrorf(err, exceptions.ErrQueryFailureToGetResponse) {
		return nil, nil, nil, err // log.CheckErrorf returns nil if error
	}

	//Get all Function data
	var fResp = make(map[string]json.RawMessage)
	startTime := time.Now()
	for _, s := range fl {
		fd, err := ExecuteFunction(s, widgetId, replacements, ctx, clt, epClt)
		if log.CheckErrorf(err, exceptions.ErrQueryFailureToGetResponse) {
			return nil, nil, nil, err
		}
		if fd != nil {
			fResp[s] = fd
		}

	}

	if len(fResp) == 0 {
		fResp = nil
		log.Debugf("No data returned by function for widget: %s", widgetId)
	} else {
		log.Infof("TIME TAKEN ExecuteFunction call  widget: %s : %v in milliseconds", widgetId, time.Since(startTime).Milliseconds())
	}

	// Get past duration queries
	p_qm, err := getWidgetPastDurationQueries(widgetId, replacements)
	if log.CheckErrorf(err, "failed to get widget's past duration queries") {
		return nil, nil, nil, err // log.CheckErrorf returns nil if error
	}

	var pResp = make(map[string]json.RawMessage)

	if p_qm != nil {
		//Get data for Trend comparison
		pResp, err = helper.GetQueryResponse(p_qm)
		if log.CheckErrorf(err, "failed to get past duration query response") {
			return nil, nil, nil, err // log.CheckErrorf returns nil if error
		}
	}
	return resp, pResp, fResp, nil

}

// getWidgetQueries returns a map of Opensearch queries and a list of functions for a widget, as defined in the widget definition JSON
func getWidgetQueries(widgetId string, replacements map[string]any, baseDataReplacements map[string]any) (map[string]db.DbQuery, []string, error) {

	queryMap := make(map[string]db.DbQuery)

	we, err := db.GetWidgetEntity(widgetId, replacements)
	if log.CheckErrorf(err, "could not get data from db in GetQueryFromDB()") {
		return nil, nil, err
	}
	var branchFilterQueries = []string{"totalCommitsHeader"}
	//build queries for current data
	for k, v := range we.Queries {
		b, err := json.Marshal(v.Query)
		if log.CheckErrorf(err, exceptions.ErrMarshallQuery, k) {
			return nil, nil, err
		}

		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, string(b))
		if log.CheckErrorf(err, "could not replace json placeholders :", k) {
			return nil, nil, err
		}
		reader := strings.NewReader(updatedJSON)

		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		queryBytes := buf.String()

		var data map[string]interface{}
		err = json.Unmarshal([]byte(queryBytes), &data)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallQuery) {
			return nil, nil, err
		}
		filterArray, isSingleQueryMap := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})

		if !isSingleQueryMap {
			criteriaArray := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]interface{})

			for i, criteria := range criteriaArray {
				filterArray := criteria.(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{})

				// Application filter (only if widget is in the list)
				if widgets, ok := db.FilterWidgetMap["application_dashboard_widgets"]; ok {
					if _, found := widgets[widgetId]; found {
						if apps, ok := replacements["application"].([]string); ok && len(apps) > 0 && apps[0] != "All" {
							filter := helper.AddTermsFilter("application_id", apps)
							filterArray = append(filterArray, filter)
						}
					} else {
						// Default: apply components filter
						if comps, ok := replacements["component"].([]string); ok && len(comps) > 0 && comps[0] != "All" {
							filter := helper.AddTermsFilter("component_id", comps)
							filterArray = append(filterArray, filter)
						}
					}
				}
				data["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]interface{})[i].(map[string]interface{})["bool"].(map[string]interface{})["must"] = filterArray
			}

		} else {

			// Application filter (only if widget is in the list)
			if widgets, ok := db.FilterWidgetMap["application_dashboard_widgets"]; ok {
				if _, found := widgets[widgetId]; found {
					if apps, ok := replacements["application"].([]string); ok && len(apps) > 0 && apps[0] != "All" {
						filter := helper.AddTermsFilter("application_id", apps)
						filterArray = append(filterArray, filter)
					}
				} else {
					// Default: apply components filter
					if comps, ok := replacements["component"].([]string); ok && len(comps) > 0 && comps[0] != "All" {
						filter := helper.AddTermsFilter("component_id", comps)
						filterArray = append(filterArray, filter)
					}
				}
			}
		}

		branch, ok := replacements[constants.REQUEST_BRANCH]
		if ok && branch != nil && isSingleQueryMap {
			// Check if the index being queried has the branch id in its mapping
			if ok1, err := helper.HasField("branch_id", v.Alias, db.FieldIndexMap); ok1 && err == nil {
				log.Debugf(exceptions.DebugAddingBranchFilterWithBranchId, branch)
				filterBranchId := helper.AddTermFilter("branch_id", branch.(string))
				filterArray = append(filterArray, filterBranchId)
			} else if widgetId == "cs1" && v.Alias == constants.COMMIT_DATA_INDEX {
				log.Debugf("Inside ResourceId Filter : %s", branch.(string))
				filterResourceId := helper.AddTermFilter("resource_id", branch.(string))
				filterArray = append(filterArray, filterResourceId)

			} else if !contains(branchFilterQueries, k) {
				log.Debugf("Inside Automation Filter")
				automations := GetAutomationsForBranch(branch.(string))
				if len(automations) > 0 {
					log.Debugf("Adding  Automation Filter for branch:%s, Automations:%v", branch, automations)

					filter4 := helper.AddTermsFilter("automation_id", automations)
					filterArray = append(filterArray, filter4)
				}
			} else {
				branchName := GetBranchNameForId(branch.(string))
				log.Debugf("Inside Branch Filter")

				if len(branchName) > 0 {
					log.Debugf("Adding  Branch Filter for branch:%s, BranchName:%v", branch, branchName)

					var filter5 map[string]interface{}
					if widgetId != "cs10" && widgetId != "css10" {
						filter5 = helper.AddTermFilter("branch", branchName)
					}

					filterArray = append(filterArray, filter5)

				}
			}
		}

		// Add "filter" critieria for tools
		if ok1, err := helper.HasField("tools", widgetId, db.FilterWidgetMap); ok1 && err == nil {
			tools, ok := replacements[constants.REQUEST_TOOLS]
			if ok && tools != nil {
				filter := helper.AddTermsFilter("tool_id", replacements["tools"].([]string))
				filterArray = append(filterArray, filter)
			}
		}

		// Add "filter" critieria for severities
		if ok1, err := helper.HasField("severities", widgetId, db.FilterWidgetMap); ok1 && err == nil {
			severities, ok := replacements[constants.REQUEST_SEVERITIES]
			if ok && severities != nil {
				filter := helper.AddTermsFilter("severity", replacements["severities"].([]string))
				filterArray = append(filterArray, filter)
			}
		}

		if ok1, err := helper.HasField("sla", widgetId, db.FilterWidgetMap); ok1 && err == nil {
			sla, ok := replacements["sla"]
			if ok && sla != nil && sla != "" {
				if sla.(pb.SlaStatus) == pb.SlaStatus_WITHIN_SLA {
					// Construct the JSON structure for the "range" query
					rangeQuery := map[string]interface{}{
						"range": map[string]interface{}{
							"sla_breach_time": map[string]interface{}{
								"gt":        "now",
								"time_zone": replacements["timeZone"],
							},
						},
					}
					filterArray = append(filterArray, rangeQuery)
				} else if sla.(pb.SlaStatus) == pb.SlaStatus_BREACHED_SLA {
					// Construct the JSON structure for the "range" query
					rangeQuery := map[string]interface{}{
						"range": map[string]interface{}{
							"sla_breach_time": map[string]interface{}{
								"lte":       "now",
								"time_zone": replacements["timeZone"],
							},
						},
					}
					// Append the constructed range query to the filterArray
					filterArray = append(filterArray, rangeQuery)
				}

			}
		}

		if len(baseDataReplacements) != 0 {
			for key, data := range baseDataReplacements {
				baseFilter := helper.AddTermsFilter(key, data.([]string))
				filterArray = append(filterArray, baseFilter)
			}
		}

		if isSingleQueryMap {
			data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterArray
		}

		modifiedData, err := json.MarshalIndent(data, "", " ")
		if log.CheckErrorf(err, exceptions.ErrJsonConversion) {
			return nil, nil, err // log.CheckErrorf returns nil if error
		}

		log.Debugf("Query string : %s", string(modifiedData))
		queryMap[k] = db.DbQuery{AliasName: v.Alias, QueryString: string(modifiedData)}
	}

	if len(queryMap) == 0 {
		queryMap = nil
	}

	//Check if queries and funtions is empty
	if queryMap == nil && len(we.Functions) == 0 && widgetId != "cs9" {
		return nil, nil, db.ErrEmptyDbData
	}

	return queryMap, we.Functions, nil
}

// Helper function to ensure that the "should" field exists
func ensureShouldField(data map[string]interface{}) {
	query, ok := data["query"].(map[string]interface{})
	if !ok {
		query = map[string]interface{}{}
		data["query"] = query
	}

	boolField, ok := query["bool"].(map[string]interface{})
	if !ok {
		boolField = map[string]interface{}{}
		query["bool"] = boolField
	}

	// Add an empty "should" field if it doesn't exist
	if _, ok := boolField["should"]; !ok {
		boolField["should"] = []interface{}{}
	}
}

// Create component comparison widget with the reponse from Opensearch
func CreateComponentComparisonWidget(widgetId string, data map[string]json.RawMessage, fd map[string]json.RawMessage, replacements map[string]any, organization *constants.Organization) (*pb.ComponentComparisonData, error) {

	widgetConfig, err := db.GetComponentComparisonConfig(widgetId, replacements)
	if log.CheckErrorf(err, "could not get data from db in CreateComponentComparisonWidget(), for widgetId: ", widgetId) {
		return nil, err
	}

	// Building final response proto structure
	componentComparisonData := pb.ComponentComparisonData{}
	definition := widgetConfig.Definition

	componentComparisonData.Title = definition.Title
	componentComparisonData.SubTitle = definition.SubTitle
	componentComparisonData.ColumnDetails = definition.ColumnDetails
	componentComparisonData.BreadCrumbTitle = definition.BreadCrumbTitle
	componentComparisonData.CompareCommonSectionDetails = definition.CompareCommonSectionDetails
	componentComparisonData.HeaderField = definition.HeaderField

	var reportData json.RawMessage

	if len(data) > 0 {
		reportData, err = ExecutePostProcessFunctionForComponentComparison(definition.PostProcessFunctionName, "", data, replacements, organization)
	} else if len(fd) > 0 {
		reportData, err = ExecutePostProcessFunctionForComponentComparison(definition.PostProcessFunctionName, "", fd, replacements, organization)
	}

	if log.CheckErrorf(err, "Error in ExecutePostProcessFunctionForComponentComparison for widgetId: %s : ", widgetId) {
		return nil, err
	}

	// Unmarshalling post-transformation data into a slice of interface{}
	var protoDataList []interface{}
	err = json.Unmarshal(reportData, &protoDataList)
	if log.CheckErrorf(err, "Failed to unmarshal post-transformation into slice, for widgetId: ", widgetId) {
		return nil, err
	}

	// Creating a slice to hold the proto messages
	protoMessages := make([]*pb.CompareReports, len(protoDataList))

	// Unmarshalling each interface{} in the slice into a CompareReports, adding it to []CompareReports
	for i, m := range protoDataList {

		marshalledData, err := json.Marshal(m.(map[string]interface{}))
		if log.CheckErrorf(err, "could not marshal data in CreateComponentComparisonWidget()") {
			return nil, err
		}

		protoMessage := &pb.CompareReports{}
		if err := protojson.Unmarshal(marshalledData, protoMessage); err != nil {
			log.Error("Failed to unmarshal JSON into proto message:", err)
		}

		protoMessages[i] = protoMessage
	}

	if len(protoMessages) > 0 {
		sortCompareReportsAlphabetically(protoMessages)
	}

	componentComparisonData.CompareReports = protoMessages

	return &componentComparisonData, nil
}

func sortCompareReports(compareReports []*pb.CompareReports) {
	sort.Slice(compareReports, func(i, j int) bool {
		return compareReports[i].IsSubOrg && !compareReports[j].IsSubOrg
	})

	for _, compareReports := range compareReports {
		if compareReports.CompareReports != nil {
			sortCompareReports(compareReports.CompareReports)
		}
	}
}

func sortCompareReportsAlphabetically(compareReports []*pb.CompareReports) {
	sort.SliceStable(compareReports, func(i, j int) bool {
		// If is_sub_org is different, sort by is_sub_org
		if compareReports[i].IsSubOrg != compareReports[j].IsSubOrg {
			return compareReports[i].IsSubOrg
		}

		// If both are sub-orgs, or both are not sub-orgs, sort alphabetically by compare_title
		return compareReports[i].CompareTitle < compareReports[j].CompareTitle
	})

	// Recursively sort sub-reports
	for _, compareReport := range compareReports {
		if compareReport.CompareReports != nil {
			sortCompareReportsAlphabetically(compareReport.CompareReports)
		}
	}
}

// Gets data directly from OpenSearch or by invoking functions for Component Comparison
func GetComponentComparisonData(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (map[string]json.RawMessage, map[string]json.RawMessage, error) {
	// Get all queries from Widget Definition using Widget Id.
	// Apply the filters in queries and get data
	queryMap, fl, err := getComponentComparisonQueries(widgetId, replacements)
	if log.CheckErrorf(err, "failed to get widget queries") {
		return nil, nil, err
	}
	var resp map[string]json.RawMessage
	if queryMap != nil {
		startTime := time.Now()
		resp, err = helper.GetMultiQueryResponse(queryMap)
		if log.CheckErrorf(err, "failed to get query response") {
			return nil, nil, err
		}
		log.Infof("TIME TAKEN GetMultiQueryResponse call widget %s : %v in milliseconds", widgetId, time.Since(startTime).Milliseconds())
	}

	//Get all Function data
	var fResp = make(map[string]json.RawMessage)
	startTime := time.Now()
	for _, s := range fl {
		fd, err := ExecuteFunction(s, widgetId, replacements, ctx, clt, epClt)
		if log.CheckErrorf(err, "failed to get query response") {
			return nil, nil, err
		}
		fResp[s] = fd
	}

	if len(fResp) == 0 {
		fResp = nil
	} else {
		log.Infof("TIME TAKEN ExecuteFunction call  widget: %s : %v in milliseconds", widgetId, time.Since(startTime).Milliseconds())
	}

	return resp, fResp, nil

}

// Get map of Opensearch queries and list of functions as per widget definition
func getComponentComparisonQueries(widgetId string, replacements map[string]any) (map[string]db.DbQuery, []string, error) {

	queryMap := make(map[string]db.DbQuery)

	config, err := db.GetComponentComparisonConfig(widgetId, replacements)
	if log.CheckErrorf(err, "could not get data from db in getComponentComparisonQueries()") {
		return nil, nil, err
	}

	// replace placeholders in the query and form the DbQquery data structure
	for k, v := range config.Queries {
		b, err := json.Marshal(v.Query)
		if log.CheckErrorf(err, exceptions.ErrMarshallQuery, k) {
			return nil, nil, err
		}

		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, string(b))
		if log.CheckErrorf(err, "could not replace json placeholders :", k) {
			return nil, nil, err
		}

		reader := strings.NewReader(updatedJSON)

		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		queryBytes := buf.String()

		var data map[string]interface{}
		err = json.Unmarshal([]byte(queryBytes), &data)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallQuery) {
			return nil, nil, err
		}
		filterArray, isSingleQueryMap := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})

		if !isSingleQueryMap {
			criteriaArray := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]interface{})

			for i, criteria := range criteriaArray {
				filterArray := criteria.(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{})

				// Application filter (only if widget is in the list)
				if widgets, ok := db.FilterWidgetMap["application_dashboard_widgets"]; ok {
					if _, found := widgets[widgetId]; found {
						if apps, ok := replacements["application"].([]string); ok && len(apps) > 0 && apps[0] != "All" {
							filter := helper.AddTermsFilter("application_id", apps)
							filterArray = append(filterArray, filter)
						}
					} else {
						// Default: apply components filter
						if comps, ok := replacements["component"].([]string); ok && len(comps) > 0 && comps[0] != "All" {
							filter := helper.AddTermsFilter("component_id", comps)
							filterArray = append(filterArray, filter)
						}
					}
				}
				data["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]interface{})[i].(map[string]interface{})["bool"].(map[string]interface{})["must"] = filterArray
			}
		} else {
			// Application filter (only if widget is in the list)
			if widgets, ok := db.FilterWidgetMap["application_dashboard_widgets"]; ok {
				if _, found := widgets[widgetId]; found {
					if apps, ok := replacements["application"].([]string); ok && len(apps) > 0 && apps[0] != "All" {
						filter := helper.AddTermsFilter("application_id", apps)
						filterArray = append(filterArray, filter)
					}
				} else {
					// Default: apply components filter
					if comps, ok := replacements["component"].([]string); ok && len(comps) > 0 && comps[0] != "All" {
						filter := helper.AddTermsFilter("component_id", comps)
						filterArray = append(filterArray, filter)
					}
				}
			}
			data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterArray
		}

		queryMap[k] = db.DbQuery{AliasName: v.Alias, QueryString: string(updatedJSON)}
	}

	if len(queryMap) == 0 {
		queryMap = nil
	}

	//Check if queries and funtions is empty
	if queryMap == nil && len(config.Functions) == 0 && widgetId != "cs9" {
		return nil, nil, db.ErrEmptyDbData
	}

	return queryMap, config.Functions, nil
}

// Get map of Opensearch queries for trend computation from past records as per widget definition
func getWidgetPastDurationQueries(widgetId string, replacements map[string]any) (map[string]db.DbQuery, error) {

	queryMap := make(map[string]db.DbQuery)

	q, err := db.GetWidgetEntity(widgetId, replacements)
	if log.CheckErrorf(err, "could not get data from db in GetQueryFromDB()") {
		return nil, err
	}

	if q.PastDurationQueries == nil {
		return nil, nil
	}

	//build queries for past data
	for k, v := range q.PastDurationQueries {
		// add filters here and append it to queryMap
		b, err := json.Marshal(v.Query)
		if log.CheckErrorf(err, exceptions.ErrMarshallQuery, k) {
			return nil, err
		}

		p_startDate, p_endDate, err := helper.CalculatePreviousDates(replacements["endDate"].(string), replacements["duration"].(string))
		if log.CheckErrorf(err, "could not calculate previous week start and end dates :") {
			return nil, err
		}

		replacements["metricName"] = k
		replacements["p_startDate"] = p_startDate
		replacements["p_endDate"] = p_endDate

		// queryWithDates := strings.NewReader(replaceJSONplaceholders(replacements, string(b)))

		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, string(b))
		if log.CheckErrorf(err, exceptions.ErrMarshallQuery, k) {
			return nil, err
		}
		reader := strings.NewReader(updatedJSON)

		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		queryBytes := buf.String()

		var data map[string]interface{}
		err = json.Unmarshal([]byte(queryBytes), &data)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallQuery) {
			return nil, err
		}

		filterArray := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})

		// Application filter (only if widget is in the list)
		if widgets, ok := db.FilterWidgetMap["application_dashboard_widgets"]; ok {
			if _, found := widgets[widgetId]; found {
				if apps, ok := replacements["application"].([]string); ok && len(apps) > 0 && apps[0] != "All" {
					filter := helper.AddTermsFilter("application_id", apps)
					filterArray = append(filterArray, filter)
				}
			} else {
				// Default: apply components filter
				if comps, ok := replacements["component"].([]string); ok && len(comps) > 0 && comps[0] != "All" {
					filter := helper.AddTermsFilter("component_id", comps)
					filterArray = append(filterArray, filter)
				}
			}
		}

		data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterArray

		modifiedData, err := json.MarshalIndent(data, "", " ")
		if log.CheckErrorf(err, exceptions.ErrJsonConversion) {
			return nil, err // log.CheckErrorf returns nil if error
		}

		queryMap[k] = db.DbQuery{AliasName: v.Alias, QueryString: string(modifiedData)}
	}

	//Check for queries empty
	if len(queryMap) == 0 {
		return nil, db.ErrEmptyDbData
	}

	return queryMap, nil
}

func combineJSONs(json1, json2 json.RawMessage) (json.RawMessage, error) {
	out := struct {
		Before json.RawMessage
		After  json.RawMessage
	}{
		Before: json1,
		After:  json2,
	}

	combinedJSON, err := json.Marshal(out)
	if log.CheckErrorf(err, "error marshaling output from combining JSON :") {
		return nil, err
	}

	return combinedJSON, nil
}

func GetComputedData(widgetId string, replacements map[string]any) (json.RawMessage, error) {

	baseQuery, err := db.ReplaceJSONplaceholders(replacements, constants.ComputedWidgetQuery)
	if log.CheckErrorf(err, "could not replace json placeholders in GetComputedData() for %s :", widgetId) {
		return nil, err
	}

	reader := strings.NewReader(baseQuery)

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	queryBytes := buf.String()

	var data map[string]interface{}
	err = json.Unmarshal([]byte(queryBytes), &data)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallQuery) {
		return nil, err
	}

	filterArray := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{})

	// Application filter (only if widget is in the list)
	if widgets, ok := db.FilterWidgetMap["application_dashboard_widgets"]; ok {
		if _, found := widgets[widgetId]; found {
			if apps, ok := replacements["application"].([]string); ok && len(apps) > 0 && apps[0] != "All" {
				filter := helper.AddTermsFilter("application_id", apps)
				filterArray = append(filterArray, filter)
			}
		} else {
			// Default: apply components filter
			if comps, ok := replacements["component"].([]string); ok && len(comps) > 0 && comps[0] != "All" {
				filter := helper.AddTermsFilter("component_id", comps)
				filterArray = append(filterArray, filter)
			}
		}
	}

	data["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = filterArray

	updatedQuery, err := json.MarshalIndent(data, "", " ")
	if log.CheckErrorf(err, exceptions.ErrJsonConversion) {
		return nil, err // log.CheckErrorf returns nil if error
	}

	queryMap := make(map[string]db.DbQuery)
	queryMap["widgetData"] = db.DbQuery{AliasName: constants.COMPUTED_INDEX, QueryString: string(updatedQuery)}

	resp, err := helper.GetQueryResponse(queryMap)
	if log.CheckErrorf(err, "could not get computed data") {
		return nil, err
	}

	return resp["widgetData"], nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
