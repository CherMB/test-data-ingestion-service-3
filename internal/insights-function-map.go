package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/calculi-corp/api/go/vsm/report"

	"github.com/gonum/floats"
	"github.com/hashicorp/go-version"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/calculi-corp/api/go"
	"github.com/calculi-corp/api/go/endpoint"
	client "github.com/calculi-corp/grpc-client"
	"github.com/calculi-corp/log"
	"github.com/calculi-corp/reports-service/constants"
	"github.com/calculi-corp/reports-service/db"
	"github.com/calculi-corp/reports-service/exceptions"
	helper "github.com/calculi-corp/reports-service/helper"
	"github.com/calculi-corp/reports-service/models"
	"github.com/opensearch-project/opensearch-go"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type CompletedRuns struct {
	JobId          string  `json:"jobId,omitempty"`
	RunId          float64 `json:"runId,omitempty"`
	RunTime        float64 `json:"runTime,omitempty"`
	StartTime      float64 `json:"startTime,omitempty"`
	Result         string  `json:"result,omitempty"`
	Name           string  `json:"name,omitempty"`
	ControllerName string  `json:"controllerName,omitempty"`
}

type ResultCount struct {
	Result string  `json:"result,omitempty"`
	Count  float64 `json:"count"`
}

type CiJobInfo struct {
	EndpointId  string `json:"endpointId"`
	JobId       string `json:"jobId"`
	JobName     string `json:"jobName"`
	Type        string `json:"type"`
	DisplayName string `json:"displayName"`
}

type JobExecutionInfo struct {
	EndpointId         string  `json:"endpointId"`
	JobId              string  `json:"jobId"`
	JobName            string  `json:"jobName"`
	DisplayName        string  `json:"displayName"`
	ControllerName     string  `json:"controllerName"`
	LastActive         string  `json:"lastActive"`
	LastActiveDuration int64   `json:"lastActiveDuration"`
	Type               string  `json:"type"`
	TotalExecuted      float64 `json:"totalExecuted"`
	Success            float64 `json:"success"`
	Failed             float64 `json:"failed"`
	Aborted            float64 `json:"aborted"`
	Unstable           float64 `json:"unstable"`
	NotBuilt           float64 `json:"notBuilt"`
	TotalDuration      float64 `json:"totalDuration"`
	AvgRunTime         float64 `json:"avgRunTime"`
	Result             string  `json:"result"`
}

type RunsInfo struct {
	TotalNumberofProjects       int64 `json:"totalNumberOfProjects"`
	TotalNumberofRuns           int64 `json:"totalNumberOfRuns"`
	TotalNumberOfActiveProjects int64 `json:"totalNumberOfActiveProjects"`
	TotalNumberOfIdleProjects   int64 `json:"totalNumberOfIdleProjects"`
}

var (
	getAllEndpoints = helper.GetAllEndpoints
)

func getInsightProjectTypes(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	outputResponse := make(map[string]interface{})
	contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
	parentIds := []string{}
	endPointsResponse, err := getAllEndpoints(ctx, epClt, replacements[constants.SUB_ORG_ID].(string), contributionIds, true)
	if len(endPointsResponse.Endpoints) > 0 {
		for _, endpoint := range endPointsResponse.Endpoints {
			parentIds = append(parentIds, endpoint.ResourceId)
		}
	}
	replacements[constants.PARENTS_IDS] = parentIds
	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	ciToolType := replacements[constants.CI_TOOL_TYPE].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	if ciToolType == constants.CJOC {
		_, _, _, err := updateReplacementsAndGetControllerInfo(ctx, epClt, replacements)
		if err != nil {
			return nil, err
		}
	}
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiProjectTypesQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	response, err := searchResponse(modifiedJson, constants.CB_CI_JOB_INFO_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	projectTypesResult := []map[string]any{}
	result := make(map[string]interface{})
	runInfo := RunsInfo{}
	json.Unmarshal([]byte(response), &result)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.PROJECT_TYPES] != nil {
			projectTypes := aggsResult[constants.PROJECT_TYPES].(map[string]interface{})
			if projectTypes[constants.VALUE] != nil {
				values := projectTypes[constants.VALUE].([]interface{})
				for _, value := range values {
					dataMap := value.(map[string]any)
					dataMap[constants.NAME] = GetProcessedJobType(dataMap[constants.NAME].(string))
					projectTypesResult = append(projectTypesResult, dataMap)
					runInfo.TotalNumberofProjects = runInfo.TotalNumberofProjects + int64(dataMap[constants.VALUE].(float64))
				}
			}
		}
	}
	outputResponse[constants.DATA] = projectTypesResult
	outputResponse[constants.RUNS_INFO] = runInfo
	responseJson, err := json.Marshal(outputResponse)
	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}
}

func updateReplacementsAndGetControllerInfo(ctx context.Context, epClt endpoint.EndpointServiceClient, replacements map[string]any) (int, int, map[string]ControllerInfo, error) {
	controllerMap := make(map[string]ControllerInfo)
	var cbciList []string
	totalControllers, connectedControllers := 0, 0
	latestUpdateTime := &timestamppb.Timestamp{}
	contributionIds := []string{constants.CBCI_ENDPOINT}
	endPointsResponse, err := getAllEndpoints(ctx, epClt, replacements[constants.SUB_ORG_ID].(string), contributionIds, true)
	if log.CheckErrorf(err, "Exception while fetching endpoints") {
		return 0, 0, nil, errors.New("endpoint api failed")
	}
	parentIds := []string{}
	if len(endPointsResponse.Endpoints) > 0 {
		for _, endpoint := range endPointsResponse.Endpoints {
			parentIds = append(parentIds, endpoint.ResourceId)
		}
	}
	replacements[constants.PARENTS_IDS] = parentIds
	if endPointsResponse != nil && len(endPointsResponse.Endpoints) > 0 {
		endpoints := endPointsResponse.Endpoints
		controllerUrls, err := GetControllerUrlsByEndpoint(replacements, ctx)
		totalControllers = len(controllerUrls)
		log.CheckErrorf(err, "Controller fetch from opensearch failed")
		if len(controllerUrls) > 0 {
			for _, endpoint := range endpoints {
				toolUrl, status := constants.EMPTY_STRING, constants.EMPTY_STRING
				for _, property := range endpoint.Properties {
					if property.Name == constants.TOOL_URL {
						if stringValue, ok := property.Value.(*api.Property_String_); ok {
							toolUrl = strings.TrimSpace(stringValue.String_)
						}
					} else if property.Name == constants.STATUS {
						if stringValue, ok := property.Value.(*api.Property_String_); ok {
							status = stringValue.String_
							if property.Audit != nil {
								latestUpdateTime = property.Audit.When
							}
						}
					}
				}
				if toolUrl != constants.EMPTY_STRING {
					for _, controllerUrl := range controllerUrls {
						if controllerUrl == toolUrl {
							info := ControllerInfo{
								Name:   endpoint.Name,
								ToolId: endpoint.Id,
							}
							if status == constants.EMPTY_STRING {
								info.Status = constants.NOT_INSTALLED
								info.LastUpdatedTime = latestUpdateTime.AsTime().UnixMilli()
							} else {
								info.Status = status
								info.LastUpdatedTime = latestUpdateTime.AsTime().UnixMilli()
							}
							connectedControllers++
							controllerMap[info.ToolId] = info
							cbciList = append(cbciList, endpoint.Id)
							break
						}
					}
				}
			}
		}
	}
	if len(cbciList) > 0 {
		replacements[constants.ENDPOINT_IDS] = cbciList
	}
	return totalControllers, connectedControllers, controllerMap, nil
}

func getInsightSystemInformation(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	outputResponse := make(map[string]interface{})
	contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
	parentIds := []string{}
	endPointsResponse, err := getAllEndpoints(ctx, epClt, replacements[constants.SUB_ORG_ID].(string), contributionIds, true)
	if len(endPointsResponse.Endpoints) > 0 {
		for _, endpoint := range endPointsResponse.Endpoints {
			parentIds = append(parentIds, endpoint.ResourceId)
		}
	}
	replacements[constants.PARENTS_IDS] = parentIds
	ciToolType := replacements[constants.CI_TOOL_TYPE].(string)
	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiToolInsightFetchQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	response, err := searchResponse(modifiedJson, constants.CB_CI_TOOL_INSIGHT_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	hits := result[constants.HITS].(map[string]interface{})[constants.HITS].([]interface{})
	if len(hits) > 0 {
		source := hits[0].(map[string]interface{})[constants.SOURCE].(map[string]interface{})
		for key, value := range source {
			if key == constants.VERSION {
				outputResponse[key] = value
				getAndUpdateLatestVersion(outputResponse, value, ciToolType)
			} else if key == constants.PLUGINS {
				plugins := value.([]interface{})
				pluginInfo := make(map[string]interface{})
				pluginInfo[constants.COUNT] = len(plugins)
				pluginInfo[constants.DRILLDOWN] = DrillDown{
					ReportId:    constants.PLUGIN_INFO,
					ReportTitle: constants.PLUGIN_INFO_TITLE,
				}
				outputResponse[constants.PLUGIN_INFO] = pluginInfo
			} else if key == constants.METRICS {
				metrics := value.([]interface{})
				for _, metric := range metrics {
					data := metric.(map[string]interface{})
					if metricValue, ok := data[constants.METRICS_DATA]; ok {
						metricData := metricValue.(map[string]interface{})
						if countValue, ok := metricData[constants.TOTAL_EXECUTOR_KEY]; ok {
							updateMetricsInResponse(countValue, outputResponse, constants.TOTAL_EXECUTORS)
						}
						if countValue, ok := metricData[constants.FREE_EXECUTOR_KEY]; ok {
							updateMetricsInResponse(countValue, outputResponse, constants.FREE_EXECUTORS)
						}
						if countValue, ok := metricData[constants.TOTAL_NODES_KEY]; ok {
							updateMetricsInResponse(countValue, outputResponse, constants.TOTAL_NODES)
						}
					}
				}
			}
		}
	}
	responseJson, err := json.Marshal(outputResponse)
	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}
}

func getAndUpdateLatestVersion(outputResponse map[string]interface{}, value interface{}, toolType string) {
	var latestVersion string
	var isStableVersion bool
	if toolType == constants.JENKINS {
		versionSplit := strings.Split(value.(string), ".")

		if len(versionSplit) == 3 {
			isStableVersion = true
		}
		if isStableVersion {
			latestVersionResponse, err := helper.Get(constants.JENKINS_STABLE_VERSION_URL, constants.EMPTY_STRING)
			if err == nil {
				latestVersion = convertResponseToString(latestVersionResponse)
			}
		} else {
			latestVersionResponse, err := helper.Get(constants.JENKINS_LATEST_VERSION_URL, constants.EMPTY_STRING)
			if err == nil {
				latestVersion = convertResponseToString(latestVersionResponse)
			}
		}

	} else if toolType == constants.CJOC || toolType == constants.CBCI || toolType == constants.JAAS {
		latestVersionResponse, err := helper.Get(constants.CBCI_CJOC_LATEST_VERSION_URL, constants.EMPTY_STRING)
		if latestVersionResponse != nil && err == nil {
			responseString := convertResponseToString(latestVersionResponse)
			if responseString != constants.EMPTY_STRING {
				var output string = ""
				for i := strings.LastIndex(responseString, constants.HREF_STRING) + 8; i < len(responseString); i++ {
					charString := string([]byte{responseString[i]})
					if charString == constants.SLASH {
						break
					}
					if (charString >= constants.ZERO && charString <= constants.NINE) || charString == constants.DOT {
						output += charString
					}
				}
				latestVersion = output
			}
		}
	}
	if latestVersion != constants.EMPTY_STRING && value != constants.EMPTY_STRING {
		outputResponse[constants.LATEST_VERSION] = latestVersion
		v1, v1Err := version.NewVersion(value.(string))
		v2, v2Err := version.NewVersion(latestVersion)
		if v1Err != nil || v2Err != nil {
			log.Debugf("Exception while creating version object", v1Err, v2Err)
			outputResponse[constants.VERSION_UPDATE_AVAILABLE] = false
			return
		}
		if v1.LessThan(v2) {
			outputResponse[constants.VERSION_UPDATE_AVAILABLE] = true
			if isStableVersion {
				outputResponse[constants.VERSION_UPDATE_MESSAGE] = "New stable version available"
				outputResponse[constants.VERSION_UPDATE_HINT] = "Latest stable version - " + latestVersion
			} else {
				outputResponse[constants.VERSION_UPDATE_MESSAGE] = "New version available"
				outputResponse[constants.VERSION_UPDATE_HINT] = "Latest version - " + latestVersion
			}
		} else {
			outputResponse[constants.VERSION_UPDATE_AVAILABLE] = false
			outputResponse[constants.VERSION_UPDATE_MESSAGE] = "Up to Date"
			outputResponse[constants.VERSION_UPDATE_HINT] = ""
		}
	} else {
		outputResponse[constants.LATEST_VERSION] = nil
		outputResponse[constants.VERSION_UPDATE_AVAILABLE] = false
		outputResponse[constants.VERSION_UPDATE_MESSAGE] = ""
		outputResponse[constants.VERSION_UPDATE_HINT] = ""
	}
}

func convertResponseToString(resp *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respBytes := buf.String()
	return string(respBytes)
}

func updateMetricsInResponse(countValue interface{}, outputResponse map[string]interface{}, key string) {
	value := countValue.(map[string]interface{})[constants.VALUE]
	if value != nil {
		outputResponse[key] = value
	} else {
		outputResponse[key] = 0
	}
}

func getInsightSystemHealth(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	outputResponse := make(map[string]interface{})
	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
	parentIds := []string{}
	endPointsResponse, err := getAllEndpoints(ctx, epClt, replacements[constants.SUB_ORG_ID].(string), contributionIds, true)
	if len(endPointsResponse.Endpoints) > 0 {
		for _, endpoint := range endPointsResponse.Endpoints {
			parentIds = append(parentIds, endpoint.ResourceId)
		}
	}
	replacements[constants.PARENTS_IDS] = parentIds
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiToolInsightFetchQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	response, err := searchResponse(modifiedJson, constants.CB_CI_TOOL_INSIGHT_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	hits := result[constants.HITS].(map[string]interface{})[constants.HITS].([]interface{})
	if len(hits) > 0 {
		source := hits[0].(map[string]interface{})[constants.SOURCE].(map[string]interface{})
		for key, value := range source {
			if key == constants.SYSTEM_HEALTH {
				healthList := value.([]interface{})
				list := []map[string]interface{}{}
				healthyCount := 0.0
				for _, health := range healthList {
					data := health.(map[string]interface{})
					if data[constants.NAME] == constants.DISK_SPACE {
						data[constants.NAME] = constants.DISK_SPACE_KEY
						data[constants.DESCRIPTION] = constants.DISK_SPACE_DESCRIPTION
					} else if data[constants.NAME] == constants.PLUGINS {
						data[constants.NAME] = constants.PLUGIN_KEY
						data[constants.DESCRIPTION] = constants.PLUGIN_DESCRIPTION
					} else if data[constants.NAME] == constants.TEMPORARY_SPACE {
						data[constants.NAME] = constants.TEMPORARY_SPACE_KEY
						data[constants.DESCRIPTION] = constants.TEMPORARY_SPACE_DESCRIPTION
					} else if data[constants.NAME] == constants.THREAD_DEADLOCK {
						data[constants.NAME] = constants.THREAD_DEADLOCK_KEY
						data[constants.DESCRIPTION] = constants.THREAD_DEADLOCK_DESCRIPTION
					}
					list = append(list, data)
					if healthy, ok := data[constants.HEALTHY]; ok {
						if healthy.(bool) {
							healthyCount++
						}
					}
				}
				outputResponse[constants.HEALTH_LIST] = list
				if healthyCount == 0 {
					outputResponse[constants.HEALTH_SCORE] = constants.ZERO_PERCENT
				} else {
					percent := int((healthyCount / float64(len(healthList))) * 100)
					outputResponse[constants.HEALTH_SCORE] = fmt.Sprint(percent) + `%`
					if percent > 80 {
						outputResponse[constants.HEALTH_STATUS] = constants.SUCCESS
					} else if percent > 50 {
						outputResponse[constants.HEALTH_STATUS] = constants.WARNING
					} else {
						outputResponse[constants.HEALTH_STATUS] = constants.FAILED
					}
				}
			}
		}
	}
	responseJson, err := json.Marshal(outputResponse)
	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}
}

func processBatch(batch []interface{}, wg *sync.WaitGroup, ch chan []CompletedRuns, controllerMap map[string]ControllerInfo, jobInfoMap map[string]CiJobInfo) {
	defer wg.Done()
	var data []CompletedRuns
	for _, value := range batch {
		run := value.(map[string]interface{})
		runID := run[constants.RUNID].(float64)
		startTime := run[constants.START_TIME_MILLIS].(float64)
		runTime := math.Max(run[constants.DURATION].(float64), 1000)
		result := strings.ToUpper(run[constants.RESULT].(string))
		jobID := run[constants.JOBID].(string)
		switch result {
		case constants.SUCCESS_KEY:
			result = "SUCCESSFUL"
		case constants.ABORTED_KEY:
			result = "CANCELED"
		case constants.FAILURE_KEY:
			result = "FAILED"
		}
		controllerID := run[constants.ENDPOINT_ID].(string)
		controller, controllerExists := controllerMap[controllerID]
		job, jobExists := jobInfoMap[jobID]
		completedRun := CompletedRuns{
			JobId:     jobID,
			RunId:     runID,
			StartTime: startTime,
			RunTime:   runTime,
			Result:    result,
		}
		if controllerExists {
			completedRun.ControllerName = controller.Name
		}
		if jobExists {
			completedRun.Name = job.JobName
		}
		data = append(data, completedRun)
	}
	ch <- data
}

func GetInsightCompletedRunsStream(replacements map[string]any, ctx context.Context, epClt endpoint.EndpointServiceClient, srv pb.ReportServiceHandler_StreamCIInsightsCompletedRunServer) error {
	var controllerMap map[string]ControllerInfo
	var err error
	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	ciToolType := replacements[constants.CI_TOOL_TYPE].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	if ciToolType == constants.CJOC {
		_, _, controllerMap, err = updateReplacementsAndGetControllerInfo(ctx, epClt, replacements)
		if err != nil {
			return err
		}
	}
	contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
	parentIds := []string{}
	endPointsResponse, err := getAllEndpoints(ctx, epClt, replacements[constants.SUB_ORG_ID].(string), contributionIds, true)
	if len(endPointsResponse.Endpoints) > 0 {
		for _, endpoints := range endPointsResponse.Endpoints {
			parentIds = append(parentIds, endpoints.ResourceId)
		}
	}
	replacements[constants.PARENTS_IDS] = parentIds
	osclient, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiCompletedRunsFetchQuery)
	if log.CheckErrorf(err, "could not replace json placeholders : %s", replacements) {
		return err
	}

	jobInfoMap, err := GetJobInfosForEndpoint(replacements, osclient)
	if log.CheckErrorf(err, "error fetching job info for endpoint: %s", replacements) {
		return err
	}

	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	response, err := searchResponse(modifiedJson, constants.CB_CI_RUN_INFO_INDEX, osclient)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	result := make(map[string]interface{})
	var data []*pb.CIInsightsCompletedRuns
	err = json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error Unmarshal Completed Runs response from Open Search") {
		return err
	}
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.COMPLETED_RUNS] != nil {
			completedRuns := aggsResult[constants.COMPLETED_RUNS].(map[string]interface{})
			if completedRuns[constants.VALUE] != nil {
				values := completedRuns[constants.VALUE].([]interface{})
				numBatches := (len(values) + constants.BATCH_SIZE - 1) / constants.BATCH_SIZE
				for i := 0; i < numBatches; i++ {
					start := i * constants.BATCH_SIZE
					end := start + constants.BATCH_SIZE
					if end > len(values) {
						end = len(values)
					}
					for _, value := range values[start:end] {
						run := value.(map[string]interface{})
						runID := float32(run[constants.RUNID].(float64))
						startTime := float32(run[constants.START_TIME_MILLIS].(float64))
						runTime := float32(math.Max(run[constants.DURATION].(float64), 1000))
						result := strings.ToUpper(run[constants.RESULT].(string))
						jobID := run[constants.JOBID].(string)
						switch result {
						case constants.SUCCESS_KEY:
							result = "SUCCESSFUL"
						case constants.ABORTED_KEY:
							result = "CANCELED"
						case constants.FAILURE_KEY:
							result = "FAILED"
						}
						controllerID := run[constants.ENDPOINT_ID].(string)
						controller, controllerExists := controllerMap[controllerID]
						job, jobExists := jobInfoMap[jobID]
						completedRun := &pb.CIInsightsCompletedRuns{
							JobId:     jobID,
							RunId:     runID,
							StartTime: startTime,
							RunTime:   runTime,
							Result:    result,
						}
						if controllerExists {
							completedRun.ControllerName = controller.Name
						}
						if jobExists {
							completedRun.Name = job.JobName
						}
						data = append(data, completedRun)
					}
					completedRunsResponse := pb.StreamCIInsightsCompletedRunsResponse{
						CompletedRuns: data,
					}
					err := srv.Send(&completedRunsResponse)
					if log.CheckErrorf(err, "Error Sending the Completed runs Stream data for batch %v", numBatches) {
						continue
					}
				}
			}
		}
	}

	return nil
}
func getInsightCompletedRuns(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	outputResponse := make(map[string]interface{})
	var controllerMap map[string]ControllerInfo
	var err error
	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	ciToolType := replacements[constants.CI_TOOL_TYPE].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	if ciToolType == constants.CJOC {
		_, _, controllerMap, err = updateReplacementsAndGetControllerInfo(ctx, epClt, replacements)
		if err != nil {
			return nil, err
		}
	}
	contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
	parentIds := []string{}
	endPointsResponse, err := getAllEndpoints(ctx, epClt, replacements[constants.SUB_ORG_ID].(string), contributionIds, true)
	if len(endPointsResponse.Endpoints) > 0 {
		for _, endpoint := range endPointsResponse.Endpoints {
			parentIds = append(parentIds, endpoint.ResourceId)
		}
	}

	replacements[constants.PARENTS_IDS] = parentIds
	client, err := openSearchClient()

	statusResultCount, err := db.ReplaceJSONplaceholders(replacements, constants.CiCompletedRunsResultQueryCount)
	modifiedStatusResultJsonCount := UpdateMustFilters(statusResultCount, replacements)
	responseStatusCount, err := searchResponse(modifiedStatusResultJsonCount, constants.CB_CI_RUN_INFO_INDEX, client)
	resultStatusCount := make(map[string]interface{})

	err = json.Unmarshal([]byte(responseStatusCount), &resultStatusCount)
	var allResultCountData []ResultCount
	if resultStatusCount[constants.AGGREGATION] != nil {
		aggsResult := resultStatusCount[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.RESULT_COUNTS] != nil {
			jobs := aggsResult[constants.RESULT_COUNTS].(map[string]interface{})
			if jobs[constants.BUCKETS] != nil {
				values := jobs[constants.BUCKETS].([]interface{})
				failureCount := float64(0)
				successCount := float64(0)
				cancelledCount := float64(0)
				unstableCount := float64(0)
				for _, value := range values {
					runCount := value.(map[string]interface{})
					resultValue := runCount[constants.KEY].(string)
					if resultValue == constants.SUCCESS_KEY {
						successCount = successCount + runCount[constants.DOC_COUNT].(float64)
					} else if resultValue == constants.FAILED_KEY || resultValue == constants.FAILURE_KEY {
						failureCount = failureCount + runCount[constants.DOC_COUNT].(float64)
					} else if resultValue == constants.ABORTED_KEY {
						cancelledCount = cancelledCount + runCount[constants.DOC_COUNT].(float64)
					} else if resultValue == constants.UNSTABLE_KEY {
						unstableCount = unstableCount + runCount[constants.DOC_COUNT].(float64)
					} else {
						failureCount = failureCount + runCount[constants.DOC_COUNT].(float64)
					}
				}
				totalResultCountData := ResultCount{}
				totalResultCountData.Result = constants.TOTAL_RUNS
				totalResultCountData.Count = successCount + failureCount + cancelledCount + unstableCount
				allResultCountData = append(allResultCountData, totalResultCountData)
				successResultCountData := ResultCount{}
				successResultCountData.Result = constants.SUCCESSFUL_HEADER
				successResultCountData.Count = successCount
				allResultCountData = append(allResultCountData, successResultCountData)
				failureResultCountData := ResultCount{}
				failureResultCountData.Result = constants.FAILED_HEADER
				failureResultCountData.Count = failureCount
				allResultCountData = append(allResultCountData, failureResultCountData)
				unstableResultCountData := ResultCount{}
				unstableResultCountData.Result = constants.UNSTABLE_HEADER
				unstableResultCountData.Count = unstableCount
				allResultCountData = append(allResultCountData, unstableResultCountData)
				cancelledResultCountData := ResultCount{}
				cancelledResultCountData.Result = constants.CANCELED_HEADER
				cancelledResultCountData.Count = cancelledCount
				allResultCountData = append(allResultCountData, cancelledResultCountData)
			}
		}
	}
	updatedJSONCount, err := db.ReplaceJSONplaceholders(replacements, constants.CiCompletedRunsFetchQueryCount)
	modifiedJsonCount := UpdateMustFilters(updatedJSONCount, replacements)
	responseCount, err := countResponse(modifiedJsonCount, constants.CB_CI_RUN_INFO_INDEX, client)
	resultCount := make(map[string]interface{})
	json.Unmarshal([]byte(responseCount), &resultCount)
	if resultCount["count"] != nil && resultCount["count"].(float64) > 20000 {
		outputResponse[constants.MAX_SIZE_REACHED] = true
		outputResponse[constants.DATA] = []map[string]string{}
		outputResponse[constants.LIGHT_COLOR_SCHEME] = []map[string]string{}
		outputResponse[constants.COLOR_SCHEME] = []map[string]string{}
		outputResponse[constants.TYPE] = constants.SCATTER_TYPE
		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	}
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiCompletedRunsFetchQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}

	jobInfoMap, err := GetJobInfosForEndpoint(replacements, client)
	if log.CheckErrorf(err, "error fetching job info for endpoint: %s", replacements) {
		return nil, err
	}

	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	response, err := searchResponse(modifiedJson, constants.CB_CI_RUN_INFO_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	result := make(map[string]interface{})
	data := []CompletedRuns{}
	json.Unmarshal([]byte(response), &result)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.COMPLETED_RUNS] != nil {
			completedRuns := aggsResult[constants.COMPLETED_RUNS].(map[string]interface{})
			if completedRuns[constants.VALUE] != nil {
				values := completedRuns[constants.VALUE].([]interface{})
				var allResults []CompletedRuns
				var wg sync.WaitGroup
				ch := make(chan []CompletedRuns, constants.NUM_WORKERS)
				numBatches := (len(values) + constants.BATCH_SIZE - 1) / constants.BATCH_SIZE
				for i := 0; i < numBatches; i++ {
					start := i * constants.BATCH_SIZE
					end := start + constants.BATCH_SIZE
					if end > len(values) {
						end = len(values)
					}
					wg.Add(1)
					go processBatch(values[start:end], &wg, ch, controllerMap, jobInfoMap)
				}
				wg.Wait()
				close(ch)
				for batchResult := range ch {
					allResults = append(allResults, batchResult...)
				}
				data = allResults
			}
		}
	}

	if len(data) > 0 {
		outputResponse[constants.MAX_SIZE_REACHED] = false
		outputResponse[constants.DATA] = data
		outputResponse[constants.COUNT_INFO] = allResultCountData
		outputResponse[constants.TYPE] = constants.SCATTER_TYPE
		drilldown := DrillDown{
			ReportId:    constants.RUN_INFORMATION_KEY,
			ReportTitle: constants.RUN_INFORMATION,
		}
		outputResponse[constants.DRILLDOWN] = drilldown
		colorSchemes := []map[string]string{}
		colorScheme := map[string]string{}
		colorScheme[constants.COLOR_0] = constants.COLOR_SCHEME_0
		colorSchemes = append(colorSchemes, colorScheme)
		colorScheme1 := map[string]string{}
		colorScheme1[constants.COLOR_0] = constants.COLOR_SCHEME_1
		colorSchemes = append(colorSchemes, colorScheme1)
		colorScheme2 := map[string]string{}
		colorScheme2[constants.COLOR_0] = constants.COLOR_SCHEME_2
		colorSchemes = append(colorSchemes, colorScheme2)
		colorScheme3 := map[string]string{}
		colorScheme3[constants.COLOR_0] = constants.COLOR_SCHEME_3
		colorSchemes = append(colorSchemes, colorScheme3)
		outputResponse[constants.COLOR_SCHEME] = colorSchemes

		lightColorSchemes := []map[string]string{}
		lightColor := map[string]string{}
		lightColor[constants.COLOR_0] = constants.LIGHT_COLOR_0
		lightColorSchemes = append(lightColorSchemes, lightColor)
		lightColor1 := map[string]string{}
		lightColor1[constants.COLOR_0] = constants.LIGHT_COLOR_1
		lightColorSchemes = append(lightColorSchemes, lightColor1)
		lightColor2 := map[string]string{}
		lightColor2[constants.COLOR_0] = constants.LIGHT_COLOR_2
		lightColorSchemes = append(lightColorSchemes, lightColor2)
		lightColor3 := map[string]string{}
		lightColor3[constants.COLOR_0] = constants.LIGHT_COLOR_3
		lightColorSchemes = append(lightColorSchemes, lightColor3)
		outputResponse[constants.LIGHT_COLOR_SCHEME] = lightColorSchemes
	}
	responseJson, err := json.Marshal(outputResponse)

	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}
}

func GetJobInfosForEndpoint(replacements map[string]any, client *opensearch.Client) (map[string]CiJobInfo, error) {
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiJobInfoFetchQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	response, err := searchResponse(modifiedJson, constants.CB_CI_JOB_INFO_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	result := make(map[string]interface{})
	data := map[string]CiJobInfo{}
	json.Unmarshal([]byte(response), &result)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.JOBS] != nil {
			jobs := aggsResult[constants.JOBS].(map[string]interface{})
			if jobs[constants.VALUE] != nil {
				values := jobs[constants.VALUE].([]interface{})
				for _, value := range values {
					jobInfo := CiJobInfo{}
					run := value.(map[string]interface{})
					jobInfo.JobId = run[constants.JOBID].(string)
					jobInfo.JobName = run[constants.JOB_NAME].(string)
					jobInfo.EndpointId = run[constants.ENDPOINT_ID].(string)
					jobInfo.DisplayName = run[constants.DISPLAY_NAME].(string)
					jobInfo.Type = run[constants.TYPE].(string)
					if _, ok := data[jobInfo.JobId]; !ok {
						data[jobInfo.JobId] = jobInfo
					}
				}
			}
		}
	}
	return data, nil
}

func getInsightProjectsActivity(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	outputResponse := make(map[string]interface{})
	client, err := openSearchClient()
	var controllerMap map[string]ControllerInfo
	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	ciToolType := replacements[constants.CI_TOOL_TYPE].(string)
	filterType, _ := replacements[constants.FILTER_TYPE].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	if ciToolType == constants.CJOC {
		_, _, controllerMap, err = updateReplacementsAndGetControllerInfo(ctx, epClt, replacements)
		if err != nil {
			return nil, err
		}
	}
	contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
	parentIds := []string{}
	endPointsResponse, err := getAllEndpoints(ctx, epClt, replacements[constants.SUB_ORG_ID].(string), contributionIds, true)
	if len(endPointsResponse.Endpoints) > 0 {
		for _, endpoint := range endPointsResponse.Endpoints {
			parentIds = append(parentIds, endpoint.ResourceId)
		}
	}
	replacements[constants.PARENTS_IDS] = parentIds
	if filterType == "" {
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiRunsExecutionInfoQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson := UpdateMustFilters(updatedJSON, replacements)
		response, err := searchResponse(modifiedJson, constants.CB_CI_RUN_INFO_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		updateProjectActivity := models.UpdateProjectActivityResponse{
			Response:       response,
			ReportId:       constants.RUN_INFORMATION_KEY,
			Replacements:   replacements,
			Client:         client,
			OutputResponse: outputResponse,
			IsIdle:         false,
			JobIds:         nil,
			IsFragile:      false,
		}
		updateProjectActivityResponse(updateProjectActivity, controllerMap)
		runsInfo := GetJobandRunsCount(client, replacements)
		jobExecutions := updateProjectActivity.OutputResponse[constants.DATA].([]JobExecutionInfo)
		runsInfo.TotalNumberOfActiveProjects = int64(len(jobExecutions))
		runsInfo.TotalNumberOfIdleProjects = runsInfo.TotalNumberofProjects - runsInfo.TotalNumberOfActiveProjects

		updateProjectActivity.OutputResponse[constants.RUNS_INFO] = runsInfo
	} else if filterType == "IdleFilter" {
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.GetExecutedJobIds)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson := UpdateMustFilters(updatedJSON, replacements)
		response, err := searchResponse(modifiedJson, constants.CB_CI_RUN_INFO_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		filterIdleProjects(response, constants.RUN_INFORMATION_KEY, replacements, err, client, outputResponse, controllerMap)

		runsInfo := GetJobandRunsCount(client, replacements)
		jobExecutions := outputResponse[constants.DATA].([]JobExecutionInfo)
		runsInfo.TotalNumberOfIdleProjects = int64(len(jobExecutions))
		runsInfo.TotalNumberOfActiveProjects = runsInfo.TotalNumberofProjects - int64(len(jobExecutions))
		outputResponse[constants.RUNS_INFO] = runsInfo

	} else if filterType == "FragileFilter" {
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiFragileJobRunsQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson := UpdateMustFilters(updatedJSON, replacements)
		response, err := searchResponse(modifiedJson, constants.CB_CI_RUN_INFO_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		updateProjectActivity := models.UpdateProjectActivityResponse{
			Response:       response,
			ReportId:       constants.RUN_INFORMATION_KEY,
			Replacements:   replacements,
			Client:         client,
			OutputResponse: outputResponse,
			IsIdle:         false,
			JobIds:         nil,
			IsFragile:      true,
		}
		updateProjectActivityResponse(updateProjectActivity, controllerMap)
		outputActiveProjectResponse := make(map[string]interface{})
		updateActiveProjectActivity, err1 := getActiveProjects(replacements, client, outputActiveProjectResponse, controllerMap)
		log.CheckErrorf(err1, "Error fetching the Active projects")

		runsInfo := GetJobandRunsCount(client, replacements)
		jobExecutions := updateActiveProjectActivity.OutputResponse[constants.DATA].([]JobExecutionInfo)
		runsInfo.TotalNumberOfActiveProjects = int64(len(jobExecutions))
		runsInfo.TotalNumberOfIdleProjects = runsInfo.TotalNumberofProjects - runsInfo.TotalNumberOfActiveProjects
		updateProjectActivity.OutputResponse[constants.RUNS_INFO] = runsInfo
	}

	responseJson, err := json.Marshal(outputResponse)
	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}
}

func getActiveProjects(replacements map[string]any, client *opensearch.Client, outputActiveProjectResponse map[string]interface{}, controllerMap map[string]ControllerInfo) (models.UpdateProjectActivityResponse, error) {
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiRunsExecutionInfoQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return models.UpdateProjectActivityResponse{}, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	response, err := searchResponse(modifiedJson, constants.CB_CI_RUN_INFO_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	updateProjectActivity := models.UpdateProjectActivityResponse{
		Response:       response,
		ReportId:       constants.RUN_INFORMATION_KEY,
		Replacements:   replacements,
		Client:         client,
		OutputResponse: outputActiveProjectResponse,
		IsIdle:         false,
		JobIds:         nil,
		IsFragile:      false,
	}
	updateProjectActivityResponse(updateProjectActivity, controllerMap)

	return updateProjectActivity, nil
}

func filterIdleProjects(response string, reportId string, replacements map[string]any, err error, client *opensearch.Client, outputResponse map[string]interface{}, controllerMap map[string]ControllerInfo) {
	var modifiedJson string
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	jobIds := []string{}
	if result["aggregations"] != nil {
		aggsResult := result["aggregations"].(map[string]interface{})
		if aggsResult["unique_jobs"] != nil {
			uniqueJobs := aggsResult["unique_jobs"].(map[string]interface{})
			if uniqueJobs["buckets"] != nil {
				buckets := uniqueJobs["buckets"].([]interface{})
				for _, bucket := range buckets {
					job := bucket.(map[string]interface{})
					jobId := job["key"].(string)
					jobIds = append(jobIds, jobId)
				}
			}
		}
	}

	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiRunsExecutionWithoutRangeQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		log.Error("Exception while fetching filtered job ids", err)
	}
	if len(jobIds) > 0 {
		replacements[constants.JOB_IDS] = jobIds
		modifiedJson = UpdateMustNotFilters(updatedJSON, replacements)
	} else {
		modifiedJson = updatedJSON
	}
	mustModifiedJson := UpdateMustFilters(modifiedJson, replacements)
	response, err = searchResponse(mustModifiedJson, constants.CB_CI_RUN_INFO_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	updateProjectActivity := models.UpdateProjectActivityResponse{
		Response:       response,
		ReportId:       reportId,
		Replacements:   replacements,
		Client:         client,
		OutputResponse: outputResponse,
		IsIdle:         true,
		JobIds:         jobIds,
		IsFragile:      false,
	}
	updateProjectActivityResponse(updateProjectActivity, controllerMap)
}

type lastActiveRecord struct {
	resultType string
	lastActive float64
}

func updateProjectActivityResponse(updateProjectActivity models.UpdateProjectActivityResponse, controllerMap map[string]ControllerInfo) {
	result := make(map[string]interface{})
	data := []JobExecutionInfo{}
	json.Unmarshal([]byte(updateProjectActivity.Response), &result)
	if result["aggregations"] != nil {
		aggsResult := result["aggregations"].(map[string]interface{})
		if aggsResult["completedRuns"] != nil {
			completedRuns := aggsResult["completedRuns"].(map[string]interface{})
			if completedRuns["buckets"] != nil {
				buckets := completedRuns["buckets"].([]interface{})
				for _, bucket := range buckets {
					var lastActiveRecords []lastActiveRecord
					run := bucket.(map[string]interface{})
					jobId := run["key"].(string)
					executionInfo := JobExecutionInfo{}
					executionInfo.JobId = jobId
					executionInfo.TotalExecuted = run["doc_count"].(float64)
					resultBuckets := run["result_buckets"].(map[string]interface{})["buckets"].(map[string]interface{})
					for resultType, resultData := range resultBuckets {
						result := resultData.(map[string]interface{})
						switch resultType {
						case "SUCCESS":
							executionInfo.Success = result["doc_count"].(float64)
							lastActiveRecords = append(lastActiveRecords, lastActiveRecord{resultType, result["last_active"].(map[string]interface{})["value"].(float64)})
						case "FAILED":
							executionInfo.Failed = result["doc_count"].(float64)
							lastActiveRecords = append(lastActiveRecords, lastActiveRecord{resultType, result["last_active"].(map[string]interface{})["value"].(float64)})
						case "ABORTED":
							executionInfo.Aborted = result["doc_count"].(float64)
							lastActiveRecords = append(lastActiveRecords, lastActiveRecord{resultType, result["last_active"].(map[string]interface{})["value"].(float64)})
						case "UNSTABLE":
							executionInfo.Unstable = result["doc_count"].(float64)
							lastActiveRecords = append(lastActiveRecords, lastActiveRecord{resultType, result["last_active"].(map[string]interface{})["value"].(float64)})
						case "NOT_BUILT":
							executionInfo.NotBuilt = result["doc_count"].(float64)
							lastActiveRecords = append(lastActiveRecords, lastActiveRecord{resultType, result["last_active"].(map[string]interface{})["value"].(float64)})
						}
						executionInfo.TotalDuration += result["total_duration"].(map[string]interface{})["value"].(float64)
					}

					include := false
					if updateProjectActivity.IsFragile && (executionInfo.Failed > 0 || executionInfo.Unstable > 0) {
						fragileRatio := (executionInfo.Failed + executionInfo.Unstable) / executionInfo.TotalExecuted
						if fragileRatio > 0.3 && fragileRatio < 0.7 {
							include = true
						}
					}
					if len(lastActiveRecords) > 0 {
						sort.Slice(lastActiveRecords, func(i, j int) bool {
							return lastActiveRecords[i].lastActive > lastActiveRecords[j].lastActive
						})
						executionInfo.Result = lastActiveRecords[0].resultType
						maxLastActive := lastActiveRecords[0].lastActive
						if maxLastActive > 0 {
							executionInfo.LastActiveDuration = time.Now().UnixMilli() - int64(maxLastActive)
							executionInfo.LastActive = getReadableDuration(executionInfo.LastActiveDuration)
						} else {
							executionInfo.LastActiveDuration = 0
							executionInfo.LastActive = constants.HYPHEN
						}
					}
					if executionInfo.TotalDuration > 0 && executionInfo.TotalExecuted > 0 {
						executionInfo.AvgRunTime = math.Ceil(executionInfo.TotalDuration / executionInfo.TotalExecuted)
						if executionInfo.AvgRunTime != 0 && executionInfo.AvgRunTime < 1000 {
							executionInfo.AvgRunTime = 1000
						}
					}
					endpointIdBuckets := run["endpoint_id"].(map[string]interface{})["buckets"].([]interface{})
					for _, endpointIdBucket := range endpointIdBuckets {
						endpointId := endpointIdBucket.(map[string]interface{})["key"].(string)
						executionInfo.EndpointId = endpointId
					}
					controller, ok := controllerMap[executionInfo.EndpointId]
					if ok {
						executionInfo.ControllerName = controller.Name
					}
					if !updateProjectActivity.IsFragile || (include && updateProjectActivity.IsFragile) {
						data = append(data, executionInfo)
					}
				}
			}
		}
	}

	jobExecutions := []JobExecutionInfo{}
	jobInfoMap, err := GetJobInfosForEndpoint(updateProjectActivity.Replacements, updateProjectActivity.Client)
	if err == nil {
		for _, id := range updateProjectActivity.JobIds {
			delete(jobInfoMap, id)
		}
		for _, run := range data {
			job, ok := jobInfoMap[run.JobId]
			if ok {
				run.JobName = job.JobName
				run.DisplayName = job.JobName
				run.Type = GetProcessedJobType(job.Type)
			}
			jobExecutions = append(jobExecutions, run)
			delete(jobInfoMap, run.JobId)
		}
		if updateProjectActivity.IsIdle && len(jobInfoMap) > 0 {
			for key, value := range jobInfoMap {
				executionInfo := JobExecutionInfo{
					JobId:       key,
					JobName:     value.JobName,
					DisplayName: value.JobName,
					EndpointId:  value.EndpointId,
					LastActive:  constants.HYPHEN,
					Type:        GetProcessedJobType(value.Type),
				}
				controller, ok := controllerMap[executionInfo.EndpointId]
				if ok {
					executionInfo.ControllerName = controller.Name
				}
				jobExecutions = append(jobExecutions, executionInfo)
			}
		}
	}
	sort.Slice(jobExecutions, func(i, j int) bool {
		return jobExecutions[i].LastActiveDuration < jobExecutions[j].LastActiveDuration
	})
	updateProjectActivity.OutputResponse[constants.DATA] = jobExecutions
	if !updateProjectActivity.IsIdle {
		drilldown := DrillDown{
			ReportId:    updateProjectActivity.ReportId,
			ReportTitle: constants.RUN_INFORMATION,
		}
		updateProjectActivity.OutputResponse[constants.DRILLDOWN] = drilldown
	}
}

func GetJobandRunsCount(client *opensearch.Client, replacements map[string]any) RunsInfo {
	runsInfo := RunsInfo{}
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiJobCount)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return runsInfo
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	response, err := countResponse(modifiedJson, constants.CB_CI_JOB_INFO_INDEX, client)
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)

	if result["count"] != nil {
		runsInfo.TotalNumberofProjects = int64(result["count"].(float64))
	} else {
		runsInfo.TotalNumberofProjects = 0
	}
	runsInfo.TotalNumberofRuns = 0
	return runsInfo
}

func GetProcessedJobType(jobType string) string {
	if constants.FREESTYLE_JOB == jobType {
		jobType = constants.FREESTYLE_TYPE
	} else if constants.WORKFLOW_JOB == jobType {
		jobType = constants.PIPELINE_TYPE
	} else if constants.MATRIX_JOB == jobType {
		jobType = constants.MULTI_CONFIG_TYPE
	} else if constants.JENKINS_FOLDER == jobType || constants.JENKINS_BLUE_STEEL_FOLDER == jobType {
		jobType = constants.MULTI_FOLDER
	} else if constants.JENKINS_BRANCH == jobType {
		jobType = constants.MULTI_BRANCH
	} else if constants.JENKINS_MULTI_JOB == jobType {
		jobType = constants.MULTI_JOB
	} else if constants.PIPELINE_JOB_TEMPLATE == jobType {
		jobType = constants.PIPELINE_TEMPLATE
	} else if constants.BACKUP_PROJECT == jobType {
		jobType = constants.BACKUP_PROJECT_TYPE
	}
	return jobType
}

func getReadableDuration(duration int64) string {
	durationStr, count := "", 0
	months := duration / (30 * 24 * 60 * 60 * 1000)
	duration = duration % (30 * 24 * 60 * 60 * 1000)
	if months > 0 {
		durationStr = fmt.Sprint(months) + constants.MONTH_STRING
		count++
	}
	days := duration / (24 * 60 * 60 * 1000)
	duration = duration % (24 * 60 * 60 * 1000)
	if days > 0 {
		durationStr += fmt.Sprint(days) + constants.DAY_STRING
		count++
	}
	if count < 2 {
		hrs := duration / (60 * 60 * 1000)
		duration = duration % (60 * 60 * 1000)
		if hrs > 0 {
			durationStr += fmt.Sprint(hrs) + constants.HOUR_STRING
			count++
		}
	}
	if count < 2 {
		min := duration / (60 * 1000)
		duration = duration % (60 * 1000)
		if min > 0 {
			durationStr += fmt.Sprint(min) + constants.MINUTE_STRING
			count++
		}
	}
	if count < 2 {
		sec := duration / 1000
		if sec > 0 {
			durationStr += fmt.Sprint(sec) + constants.SECOND_STRING
			count++
		}
	}
	if count == 0 {
		durationStr = constants.HYPHEN
	}
	return durationStr
}

type ControllerInfo struct {
	Name              string  `json:"name"`
	ToolId            string  `json:"toolId"`
	Status            string  `json:"status"`
	TotalNumberOfJobs float64 `json:"totalNumberOfJobs"`
	TotalNumberOfRuns float64 `json:"totalNumberOfRuns"`
	Passed            float64 `json:"passed"`
	Failed            float64 `json:"failed"`
	Aborted           float64 `json:"aborted"`
	Unstable          float64 `json:"unstable"`
	NotBuilt          float64 `json:"notBuilt"`
	FailureRate       int32   `json:"failureRate"`
	LastActive        float64 `json:"lastActive"`
	LastUpdatedTime   int64   `json:"lastUpdatedTime"`
}

func getInsightCjocControllers(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	cbciInfoList := []map[string]any{}
	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	totalControllers, connectedControllers, controllerMap, err := updateReplacementsAndGetControllerInfo(ctx, epClt, replacements)
	if err != nil {
		return nil, err
	}
	contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
	parentIds := []string{}
	endPointsResponse, err := getAllEndpoints(ctx, epClt, replacements[constants.SUB_ORG_ID].(string), contributionIds, true)
	if len(endPointsResponse.Endpoints) > 0 {
		for _, endpoint := range endPointsResponse.Endpoints {
			parentIds = append(parentIds, endpoint.ResourceId)
		}
	}
	replacements[constants.PARENTS_IDS] = parentIds
	outputResponse := make(map[string]interface{})
	if len(controllerMap) > 0 {
		client, err := openSearchClient()
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiEndpointJobsQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson := UpdateMustFilters(updatedJSON, replacements)
		jobResponse, err := searchResponse(modifiedJson, constants.CB_CI_JOB_INFO_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

		updatedJSON, err = db.ReplaceJSONplaceholders(replacements, constants.CiJobRunsQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson = UpdateMustFilters(updatedJSON, replacements)
		response, err := searchResponse(modifiedJson, constants.CB_CI_RUN_INFO_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

		result := make(map[string]interface{})
		json.Unmarshal([]byte(response), &result)
		jobRunMap := map[string]map[string]interface{}{}
		if result["aggregations"] != nil {
			aggsResult := result["aggregations"].(map[string]interface{})
			if aggsResult["completedRuns"] != nil {
				completedRuns := aggsResult["completedRuns"].(map[string]interface{})
				if completedRuns["buckets"] != nil {
					buckets := completedRuns["buckets"].([]interface{})
					for _, bucket := range buckets {
						lastActiveRecords := []float64{}
						run := bucket.(map[string]interface{})
						jobId := run["key"].(string)
						runMap := make(map[string]interface{})
						runMap[constants.TOTAL] = run["doc_count"].(float64)
						resultBuckets := run["result_buckets"].(map[string]interface{})["buckets"].(map[string]interface{})
						for resultType, resultData := range resultBuckets {
							result := resultData.(map[string]interface{})
							switch resultType {
							case "SUCCESS":
								runMap[constants.SUCCESS] = result["doc_count"].(float64)
								lastActiveRecords = append(lastActiveRecords, result["last_active"].(map[string]interface{})["value"].(float64))
							case "FAILED":
								runMap[constants.FAILED] = result["doc_count"].(float64)
								lastActiveRecords = append(lastActiveRecords, result["last_active"].(map[string]interface{})["value"].(float64))
							case "ABORTED":
								runMap[constants.ABORTED] = result["doc_count"].(float64)
								lastActiveRecords = append(lastActiveRecords, result["last_active"].(map[string]interface{})["value"].(float64))
							case "UNSTABLE":
								runMap[constants.UNSTABLE] = result["doc_count"].(float64)
								lastActiveRecords = append(lastActiveRecords, result["last_active"].(map[string]interface{})["value"].(float64))
							case "NOT_BUILT":
								runMap[constants.NOT_BUILT] = result["doc_count"].(float64)
								lastActiveRecords = append(lastActiveRecords, result["last_active"].(map[string]interface{})["value"].(float64))
							}
						}
						if len(lastActiveRecords) > 0 {
							runMap[constants.LAST_ACTIVE] = slices.Max(lastActiveRecords)
						}
						jobRunMap[jobId] = runMap
					}
				}
			}
		}
		endpointJobResult := make(map[string]interface{})
		json.Unmarshal([]byte(jobResponse), &endpointJobResult)
		if endpointJobResult[constants.AGGREGATION] != nil {
			aggsResult := endpointJobResult[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.ENDPOINT_JOBS] != nil {
				endpointJobs := aggsResult[constants.ENDPOINT_JOBS].(map[string]interface{})
				if endpointJobs[constants.VALUE] != nil {
					values := endpointJobs[constants.VALUE].(map[string]interface{})
					for key, value := range values {
						lastActives := []float64{}
						controller, ok := controllerMap[key]
						if ok {
							jobIds := value.([]interface{})
							controller.TotalNumberOfJobs += float64(len(jobIds))
							for _, jobId := range jobIds {
								jobRunMap, runOk := jobRunMap[jobId.(string)]
								if runOk {
									if total, ok := jobRunMap[constants.TOTAL]; ok {
										controller.TotalNumberOfRuns += total.(float64)
									}
									if success, ok := jobRunMap[constants.SUCCESS]; ok {
										controller.Passed += success.(float64)
									}
									if failed, ok := jobRunMap[constants.FAILED]; ok {
										controller.Failed += failed.(float64)
									}
									if aborted, ok := jobRunMap[constants.ABORTED]; ok {
										controller.Aborted += aborted.(float64)
									}
									if unstable, ok := jobRunMap[constants.UNSTABLE]; ok {
										controller.Unstable += unstable.(float64)
									}
									if notBuilt, ok := jobRunMap[constants.NOT_BUILT]; ok {
										controller.NotBuilt += notBuilt.(float64)
									}
									if lastActive, ok := jobRunMap[constants.LAST_ACTIVE]; ok {
										lastActives = append(lastActives, lastActive.(float64))
									}
								}
							}
							if len(lastActives) > 0 {
								controller.LastActive = slices.Max(lastActives)
							}
							controllerMap[key] = controller
						}
					}
				}
			}
		}
		for _, info := range controllerMap {
			if info.TotalNumberOfRuns > 0 && info.Failed > 0 {
				info.FailureRate = int32((info.Failed / info.TotalNumberOfRuns) * 100)
			}
			controllerInfo, err := json.Marshal(info)
			if err != nil {
				log.Debugf("Exception while marshaling controller ", info)
			} else {
				controllerMap := map[string]any{}
				json.Unmarshal(controllerInfo, &controllerMap)
				cbciInfoList = append(cbciInfoList, controllerMap)
			}
		}
	}
	outputResponse[constants.DATA] = cbciInfoList
	outputResponse[constants.CONNECTED_CONTROLLERS] = connectedControllers
	outputResponse[constants.TOTAL_CONTROLLERS] = totalControllers

	responseJson, err := json.Marshal(outputResponse)
	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}
}

func getInsightRunsOverview(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	runsdata := make(map[string]interface{})
	runswaitingdata := make(map[string]interface{})
	idledata := make(map[string]interface{})
	waitingTime := make(map[string]interface{})
	idleTime := make(map[string]interface{})
	//valueMap := make(map[string][]float64)
	actualValueMap := make(map[string][]float64)
	keys := []string{}

	list := []map[string]interface{}{}
	runsdata[constants.NAME] = constants.ACTIVE_RUNS
	runsdata[constants.DESCRIPTION] = constants.ACTIVE_RUNS_DESCRIPTION
	runsdata[constants.ACTUAL] = constants.ZERO
	//runsdata[constants.EXPECTED] = constants.ZERO
	list = append(list, runsdata)

	idledata[constants.NAME] = constants.IDLE_EXECUTORS
	idledata[constants.DESCRIPTION] = constants.IDLE_EXECUTORS_DESCRIPTION
	idledata[constants.ACTUAL] = constants.ZERO
	//idledata[constants.EXPECTED] = constants.ZERO
	list = append(list, idledata)

	runswaitingdata[constants.NAME] = constants.WAITING_RUNS
	runswaitingdata[constants.DESCRIPTION] = constants.WAITING_RUNS_DESCRIPTION
	runswaitingdata[constants.ACTUAL] = constants.ZERO
	//runswaitingdata[constants.EXPECTED] = constants.ZERO
	list = append(list, runswaitingdata)

	waitingTime[constants.NAME] = constants.WAITING_TIME
	waitingTime[constants.DESCRIPTION] = constants.WAITING_TIME_DESCRIPTION
	waitingTime[constants.ACTUAL] = constants.ZERO
	//waitingTime[constants.EXPECTED] = constants.ZERO
	list = append(list, waitingTime)

	idleTime[constants.NAME] = constants.IDLE_TIME
	idleTime[constants.DESCRIPTION] = constants.IDLE_TIME_DESCRIPTOR
	idleTime[constants.ACTUAL] = constants.ZERO
	//idleTime[constants.EXPECTED] = constants.ZERO
	list = append(list, idleTime)

	contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
	parentIds := []string{}
	endPointsResponse, err := getAllEndpoints(ctx, epClt, replacements[constants.SUB_ORG_ID].(string), contributionIds, true)
	if len(endPointsResponse.Endpoints) > 0 {
		for _, endpoint := range endPointsResponse.Endpoints {
			parentIds = append(parentIds, endpoint.ResourceId)
		}
	}
	replacements[constants.PARENTS_IDS] = parentIds

	hour, day, _, _, _, err := getHourDayAndCount(replacements)
	if err != nil {
		return nil, err
	}
	replacements[constants.HOUR] = hour
	replacements[constants.DAY] = day

	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	ciToolType := replacements[constants.CI_TOOL_TYPE].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	if ciToolType == constants.CJOC {
		_, _, _, err = updateReplacementsAndGetControllerInfo(ctx, epClt, replacements)
		if err != nil {
			return nil, err
		}
	}

	outputResponse := make(map[string]interface{})
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	//allReplacement := getActivityOverviewReplacements(replacements, endDate, startDate)
	//updatedJSON, err := db.ReplaceJSONplaceholders(allReplacement, constants.CiActivityOverviewFetchQuery)
	//if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
	//	return nil, err
	//}
	//
	//modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	//response, err := searchResponse(modifiedJson, constants.CB_CI_ACTIVITY_OVERVIEW, client)
	//log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiActivityOverviewActualFetchQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	actualResponse, err := searchResponse(modifiedJson, constants.CB_CI_RUNS_ACTIVITY, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	actualResult := make(map[string]interface{})
	json.Unmarshal([]byte(actualResponse), &actualResult)
	if aggregations, ok := actualResult[constants.AGGREGATION].(map[string]interface{}); ok {
		if latestPerOrg, ok := aggregations[constants.LATEST_ORG].(map[string]interface{}); ok {
			if buckets, ok := latestPerOrg[constants.BUCKETS].([]interface{}); ok {
				for _, bucket := range buckets {
					if bucketMap, ok := bucket.(map[string]interface{}); ok {
						if latestDoc, ok := bucketMap[constants.LATEST_DOC].(map[string]interface{}); ok {
							if hits, ok := latestDoc[constants.HITS].(map[string]interface{}); ok {
								if hitsList, ok := hits[constants.HITS].([]interface{}); ok {
									for _, hit := range hitsList {
										if hitMap, ok := hit.(map[string]interface{}); ok {
											source := hitMap[constants.SOURCE].(map[string]interface{})
											for key, value := range source {
												if key == constants.ACTIVE_RUNS_KEY {
													actualValueMap[key] = append(actualValueMap[key], value.(float64))
													keys = append(keys, key)
												}
												if key == constants.WAITING_RUNS_KEY {
													actualValueMap[key] = append(actualValueMap[key], value.(float64))
													keys = append(keys, key)
												}
												if key == constants.IDLE_EXECUTORS_KEY {
													actualValueMap[key] = append(actualValueMap[key], value.(float64))
													keys = append(keys, key)
												}
												if key == constants.WAITING_TIME_KEY {
													actualValueMap[key] = append(actualValueMap[key], value.(float64))
													keys = append(keys, key)
												}
												if key == constants.IDLE_TIME_KEY {
													actualValueMap[key] = append(actualValueMap[key], value.(float64))
													keys = append(keys, key)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	//result := make(map[string]interface{})
	//json.Unmarshal([]byte(response), &result)
	//hits := result[constants.HITS].(map[string]interface{})[constants.HITS].([]interface{})
	//if len(hits) > 0 {
	//	for _, i := range hits {
	//		source := i.(map[string]interface{})[constants.SOURCE].(map[string]interface{})
	//		for key, value := range source {
	//			if key == constants.ACTIVE_RUNS_KEY {
	//				valueMap[key] = append(valueMap[key], value.(float64))
	//			}
	//			if key == constants.WAITING_RUNS_KEY {
	//				valueMap[key] = append(valueMap[key], value.(float64))
	//			}
	//			if key == constants.IDLE_EXECUTORS_KEY {
	//				valueMap[key] = append(valueMap[key], value.(float64))
	//			}
	//			if key == constants.WAITING_TIME_KEY {
	//				valueMap[key] = append(valueMap[key], value.(float64))
	//			}
	//			if key == constants.IDLE_TIME_KEY {
	//				valueMap[key] = append(valueMap[key], value.(float64))
	//			}
	//		}
	//	}
	//}

	if slices.Contains(keys, constants.ACTIVE_RUNS_KEY) {
		//sum := floats.Sum(valueMap[constants.ACTIVE_RUNS_KEY])
		//runsdata[constants.EXPECTED] = strconv.FormatFloat((math.Ceil(sum / count)), 'f', -1, 64)
		runsdata[constants.ACTUAL] = strconv.FormatFloat(floats.Sum(actualValueMap[constants.ACTIVE_RUNS_KEY]), 'f', -1, 64)
	}
	if slices.Contains(keys, constants.WAITING_RUNS_KEY) {
		//sum := floats.Sum(valueMap[constants.WAITING_RUNS_KEY])
		//runswaitingdata[constants.EXPECTED] = strconv.FormatFloat((math.Ceil(sum / count)), 'f', -1, 64)
		runswaitingdata[constants.ACTUAL] = strconv.FormatFloat(floats.Sum(actualValueMap[constants.WAITING_RUNS_KEY]), 'f', -1, 64)
	}
	if slices.Contains(keys, constants.IDLE_EXECUTORS_KEY) {
		//sum := floats.Sum(valueMap[constants.IDLE_EXECUTORS_KEY])
		//idledata[constants.EXPECTED] = strconv.FormatFloat((math.Ceil(sum / count)), 'f', -1, 64)
		idledata[constants.ACTUAL] = strconv.FormatFloat(floats.Sum(actualValueMap[constants.IDLE_EXECUTORS_KEY]), 'f', -1, 64)
	}
	if slices.Contains(keys, constants.WAITING_TIME_KEY) {
		//sum := floats.Sum(valueMap[constants.WAITING_TIME_KEY])
		//waitingTime[constants.EXPECTED] = convertMilliSecToTime(math.Ceil(sum / count))
		waitingTime[constants.ACTUAL] = convertMilliSecToTime(floats.Sum(actualValueMap[constants.WAITING_TIME_KEY]))
	}
	if slices.Contains(keys, constants.IDLE_TIME_KEY) {
		//sum := floats.Sum(valueMap[constants.IDLE_TIME_KEY])
		//idleTime[constants.EXPECTED] = convertMilliSecToTime(math.Ceil(sum / count))
		idleTime[constants.ACTUAL] = convertMilliSecToTime(floats.Sum(actualValueMap[constants.IDLE_TIME_KEY]))
	}

	outputResponse[constants.DATA] = list
	responseJson, err := json.Marshal(outputResponse)
	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}
}

func getActivityOverviewReplacements(replacements map[string]any, endDate time.Time, startDate time.Time) map[string]any {
	allReplacement := maps.Clone(replacements)
	difference := endDate.Sub(startDate)
	allReplacement[constants.START_DATE] = startDate.Add(-difference).Format(constants.DATE_FORMAT_WITH_HYPHEN)
	allReplacement[constants.END_DATE] = endDate.Add(-difference).Format(constants.DATE_FORMAT_WITH_HYPHEN)
	return allReplacement
}

func convertMilliSecToTime(milliseconds float64) string {
	duration := time.Duration(int64(milliseconds)) * time.Millisecond
	days := int(duration / (24 * time.Hour))
	duration = duration % (24 * time.Hour)

	hours := int(duration / time.Hour)
	duration = duration % time.Hour

	minutes := int(duration / time.Minute)
	duration = duration % time.Minute

	seconds := int(duration / time.Second)

	var timeString string

	if days > 0 {
		timeString += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 {
		timeString += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		timeString += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 {
		timeString += fmt.Sprintf("%ds ", seconds)
	}

	if days == 0 && hours == 0 && minutes == 0 && seconds == 0 {
		timeString += constants.ZERO
	}
	log.Debugf("The timeString value that is going to get trimmed is %s", timeString)
	result := strings.TrimSuffix(timeString, " ")
	if result == "" || result == constants.ZERO {
		result = constants.EMPTY_SECOND
	}
	return result
}

func getHourDayAndCount(replacements map[string]any) (int, string, float64, time.Time, time.Time, error) {
	currentTime := time.Now()
	hour := currentTime.Hour()
	day := currentTime.Format(constants.DAY_FORMAT)
	startDate, err := time.Parse(constants.DATE_PARSE, replacements[constants.START_DATE].(string))
	if err != nil {
		log.CheckErrorf(err, exceptions.ErrParsingStartDate, startDate)
		return hour, day, 0, time.Now(), time.Now(), err
	}
	endDate, err := time.Parse(constants.DATE_PARSE, replacements[constants.END_DATE].(string))
	if err != nil {
		log.CheckErrorf(err, "Error parsing end date  :", endDate)
		return hour, day, 0, time.Now(), time.Now(), err
	}

	var count float64
	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		if current.Weekday() == currentTime.Weekday() {
			count++
		}
		current = current.Add(time.Hour * 24)
	}

	return hour, day, count, startDate, endDate, nil
}

type ActivityOverview struct {
	ActivityTime          float64 `json:"activity_time"`
	CurrentTimeToIdle     float64 `json:"current_time_to_idle"`
	ActiveRuns            float64 `json:"active_runs"`
	ActivityDay           string  `json:"activity_day"`
	IdleExecutor          float64 `json:"idle_executor"`
	RunsWaitingToStart    float64 `json:"runs_waiting_to_start"`
	AvgTimeWaitingToStart float64 `json:"avg_time_waiting_to_start"`
	EndpointId            string  `json:"endpoint_id"`
}

type ActivityInfo struct {
	Name   string
	Values map[string]interface{}
}

func ConvertTime(input string, timezone string) (string, error) {
	// Parse the input time string
	inputTime, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN, input)
	if err != nil {
		return "", err
	}

	// Specify the target timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	// Convert to target timezone
	localTime := inputTime.In(loc)

	// Format the local time
	outputStr := localTime.Format(constants.DATE_FORMAT_WITH_HYPHEN)

	return outputStr, nil
}
func getInsightUsagePatterns(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	ciToolType := replacements[constants.CI_TOOL_TYPE].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
	parentIds := []string{}
	endPointsResponse, err := getAllEndpoints(ctx, epClt, replacements[constants.SUB_ORG_ID].(string), contributionIds, true)
	if len(endPointsResponse.Endpoints) > 0 {
		for _, endpoint := range endPointsResponse.Endpoints {
			parentIds = append(parentIds, endpoint.ResourceId)
		}
	}
	replacements[constants.PARENTS_IDS] = parentIds
	if ciToolType == constants.CJOC {
		_, _, _, err = updateReplacementsAndGetControllerInfo(ctx, epClt, replacements)
		if err != nil {
			return nil, err
		}
	}
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiUsagePatterns)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	responses, err := searchResponse(modifiedJson, constants.CB_CI_RUNS_ACTIVITY, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	dayAndTimeMap := make(map[string]map[float64][]ActivityOverview)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(responses), &result); err != nil {
		log.CheckErrorf(err, exceptions.ErrDefaultRespTemplate)
		return nil, err
	}
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.ACTIVITIES] != nil {
			activities := aggsResult[constants.ACTIVITIES].(map[string]interface{})
			if activities[constants.BUCKETS] != nil {
				buckets := activities[constants.BUCKETS].([]interface{})
				for _, bucket := range buckets {
					bucketMap := bucket.(map[string]interface{})
					if bucketMap["endpoint_ids"] != nil {
						endpointIds := bucketMap["endpoint_ids"].(map[string]interface{})
						if endpointIds["buckets"] != nil {
							endpointBuckets := endpointIds["buckets"].([]interface{})
							for _, endpointBucket := range endpointBuckets {
								endpointBucketMap := endpointBucket.(map[string]interface{})
								createdAt := bucketMap["key_as_string"].(string)
								createdAtTime, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN, createdAt)
								if err != nil {
									log.CheckErrorf(err, "Error parsing created at:", err)
									return nil, err
								}
								activityDay := createdAtTime.Format("Mon")
								activityHour := float64(createdAtTime.Hour())
								overview := ActivityOverview{
									ActiveRuns:            math.Ceil(endpointBucketMap["avg_active_runs"].(map[string]interface{})["value"].(float64)),
									ActivityDay:           activityDay,
									ActivityTime:          activityHour,
									CurrentTimeToIdle:     math.Ceil(endpointBucketMap["avg_current_time_to_idle"].(map[string]interface{})["value"].(float64)),
									IdleExecutor:          math.Ceil(endpointBucketMap["avg_idle_executor"].(map[string]interface{})["value"].(float64)),
									RunsWaitingToStart:    math.Ceil(endpointBucketMap["avg_runs_waiting_to_start"].(map[string]interface{})["value"].(float64)),
									AvgTimeWaitingToStart: math.Ceil(endpointBucketMap["avg_time_waiting_to_start"].(map[string]interface{})["value"].(float64)),
									EndpointId:            endpointBucketMap["key"].(string),
								}
								if _, ok := dayAndTimeMap[activityDay]; !ok {
									dayAndTimeMap[activityDay] = make(map[float64][]ActivityOverview)
								}
								dayAndTimeMap[activityDay][activityHour] = append(dayAndTimeMap[activityDay][activityHour], overview)
							}
						}
					}
				}
			}
		}
	}

	start, err := time.Parse(constants.DATE_PARSE, replacements[constants.START_DATE].(string))
	if err != nil {
		log.CheckErrorf(err, exceptions.ErrParsingStartDate, start)
		return nil, err
	}
	end, err := time.Parse(constants.DATE_PARSE, replacements[constants.END_DATE].(string))
	if err != nil {
		log.CheckErrorf(err, exceptions.ErrParsingStartDate, end)
		return nil, err

	}
	dayCountMap := make(map[int]int)
	oneDay := 24 * time.Hour
	currentDate := start
	for currentDate.Before(end) || currentDate.Equal(end) {
		dayOfWeek := int(currentDate.Weekday()) + 1
		dayCountMap[dayOfWeek]++
		currentDate = currentDate.Add(oneDay)
	}
	weekDays := []string{constants.SUNDAY, constants.MONDAY, constants.TUESDAY, constants.WEDNESDAY, constants.THURSDAY, constants.FRIDAY, constants.SATURDAY}
	var activityInfoList []ActivityInfo
	for _, day := range weekDays {
		activityInfo := ActivityInfo{Name: day}
		valueMap := make(map[string]interface{})
		timeOverviewMap, ok := dayAndTimeMap[day]
		if ok {
			for _, hour := range allHours() {
				hourAMPM := convertToAMPM(hour, replacements["timeFormat"].(string))
				valueMap[hourAMPM] = 0
				activityList, ok := timeOverviewMap[float64(hour)]
				if ok {
					if ciToolType == constants.CJOC {
						endpointMap := make(map[string][]ActivityOverview)
						for _, activity := range activityList {
							if _, ok := endpointMap[activity.EndpointId]; !ok {
								endpointMap[activity.EndpointId] = []ActivityOverview{}
							}
							endpointMap[activity.EndpointId] = append(endpointMap[activity.EndpointId], activity)
						}
						var activityData int64 = 0
						for _, list := range endpointMap {
							dataValue := getCalculatedValue(list, replacements, day, dayCountMap)
							activityData += int64(dataValue)
						}
						valueMap[hourAMPM] = activityData
					} else {
						valueMap[hourAMPM] = getCalculatedValue(activityList, replacements, day, dayCountMap)
					}
				}
			}
		} else {
			for _, hour := range allHours() {
				hourAMPM := convertToAMPM(hour, replacements["timeFormat"].(string))
				valueMap[hourAMPM] = 0
			}
		}
		activityInfo.Values = valueMap
		activityInfoList = append(activityInfoList, activityInfo)
	}
	outputResponse := make(map[string]interface{})
	color := make(map[string]interface{})
	color1 := make(map[string]interface{})
	color[constants.COLOR_0] = constants.UP_COLOR_0
	color1[constants.COLOR_0] = constants.UP_COLOR_1
	colorschemeList := []map[string]interface{}{}
	lightcolorschemeList := []map[string]interface{}{}
	colorschemeList = append(colorschemeList, color)
	lightcolorschemeList = append(lightcolorschemeList, color1)
	outputResponse[constants.TYPE] = constants.SCATTER
	outputResponse[constants.COLOR_SCHEME] = colorschemeList
	outputResponse[constants.LIGHT_COLOR_SCHEME] = lightcolorschemeList
	list := []map[string]interface{}{}
	for _, i := range activityInfoList {
		daylist := make(map[string]interface{})
		daylist[constants.NAME] = i.Name
		daylist[constants.VALUES] = i.Values
		list = append(list, daylist)
	}
	outputResponse[constants.DATA] = list
	responseJson, err := json.Marshal(outputResponse)
	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}
}

func getCalculatedValue(activityList []ActivityOverview, replacements map[string]any, day string, dayCountMap map[int]int) int {
	totalValue := 0.0
	for _, activity := range activityList {
		val := reflect.ValueOf(activity)
		field := val.FieldByName(replacements[constants.VIEW_OPTION].(string))
		if field.IsValid() && field.Float() > 0 {
			totalValue += field.Float()
		}
	}
	dayCount := len(activityList)
	weekDay := getWeekDayByName(day)
	if weekDay != 0 {
		dayCount = dayCountMap[int(weekDay)]
	}
	dataValue := 0
	if totalValue > 0 && dayCount > 0 {
		dataValue = int(math.Ceil(totalValue / float64(dayCount)))
		if (replacements[constants.VIEW_OPTION].(string) == constants.IDLE_TIME_VIEW || replacements[constants.VIEW_OPTION].(string) == constants.WAITING_TIME_VIEW) && dataValue < 1000 && dataValue > 0 {
			dataValue = 1000
		}
	}
	return dataValue
}

const (
	SUNDAY int = iota + 1
	MONDAY
	TUESDAY
	WEDNESDAY
	THURSDAY
	FRIDAY
	SATURDAY
)

func getWeekDayByName(dayName string) int {
	switch dayName {
	case constants.SUNDAY:
		return SUNDAY
	case constants.MONDAY:
		return MONDAY
	case constants.TUESDAY:
		return TUESDAY
	case constants.WEDNESDAY:
		return WEDNESDAY
	case constants.THURSDAY:
		return THURSDAY
	case constants.FRIDAY:
		return FRIDAY
	case constants.SATURDAY:
		return SATURDAY
	default:
		return 0
	}
}

func allHours() []int {
	return []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
}

func convertToAMPM(hour int, format string) string {
	if format == "24h" {
		if hour >= 0 && hour < 10 {
			return fmt.Sprintf("0%d:00", hour)
		}
		return fmt.Sprintf("%d:00", hour)
	}
	if hour >= 1 && hour < 12 {
		return fmt.Sprintf("%dam", hour)
	} else if hour == 12 {
		return constants.PM_12
	} else if hour == 0 {
		return constants.AM_12
	} else {
		return fmt.Sprintf("%dpm", hour-12)
	}
}

func GetControllerUrlMap(replacements map[string]any, ctx context.Context) (map[string][]string, error) {
	resultMap := make(map[string][]string)
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiCjocControllersFetchQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	responses, err := searchResponse(modifiedJson, constants.CB_CI_CJOC_CONTROLLER_INFO, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responses), &response); err != nil {
		log.CheckErrorf(err, exceptions.ErrDefaultRespTemplate)
		return nil, err
	}
	result := make(map[string]interface{})
	json.Unmarshal([]byte(responses), &result)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.CJOC_CONTROLLER_INFO] != nil {
			cjocControllerInfo := aggsResult[constants.CJOC_CONTROLLER_INFO].(map[string]interface{})
			if cjocControllerInfo[constants.VALUE] != nil {
				values := cjocControllerInfo[constants.VALUE].(map[string]interface{})
				for key, value := range values {
					controllers := value.([]interface{})
					urlList := []string{}
					for _, controllerUrl := range controllers {
						urlList = append(urlList, controllerUrl.(string))
					}
					resultMap[key] = urlList
				}
			}
		}
	}
	return resultMap, nil
}

func GetControllerUrlsByEndpoint(replacements map[string]any, ctx context.Context) ([]string, error) {
	controllersList := []string{}
	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiCjocControllersFetchByEndpoint)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	responses, err := searchResponse(modifiedJson, constants.CB_CI_CJOC_CONTROLLER_INFO, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responses), &response); err != nil {
		log.CheckErrorf(err, exceptions.ErrDefaultRespTemplate)
		return nil, err
	}
	result := make(map[string]interface{})
	json.Unmarshal([]byte(responses), &result)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.CJOC_CONTROLLER_INFO] != nil {
			cjocControllerInfo := aggsResult[constants.CJOC_CONTROLLER_INFO].(map[string]interface{})
			if cjocControllerInfo[constants.VALUE] != nil {
				values := cjocControllerInfo[constants.VALUE].([]interface{})
				for _, value := range values {
					controllersList = append(controllersList, value.(string))
				}
			}
		}
	}
	return controllersList, nil
}

func JobAndRunCount(replacements map[string]any, ctx context.Context) (map[string]int, map[string]int, error) {
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiAllJobCount)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, nil, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	jobResponses, err := searchResponse(modifiedJson, constants.CB_CI_JOB_INFO_INDEX, client)
	runResponses, err := searchResponse(modifiedJson, constants.CB_CI_RUN_INFO_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	var jobResponse map[string]interface{}
	var runResponse map[string]interface{}
	if err := json.Unmarshal([]byte(jobResponses), &jobResponse); err != nil {
		log.CheckErrorf(err, exceptions.ErrDefaultRespTemplate)
		return nil, nil, err
	}
	if err := json.Unmarshal([]byte(runResponses), &runResponse); err != nil {
		log.CheckErrorf(err, exceptions.ErrDefaultRespTemplate)
		return nil, nil, err
	}
	jobCounts := make(map[string]int)
	if jobResponse[constants.AGGREGATION] != nil {
		aggsResult := jobResponse[constants.AGGREGATION].(map[string]interface{})
		if aggsResult["jobs_per_endpoint"] != nil {
			jobInfo := aggsResult["jobs_per_endpoint"].(map[string]interface{})
			if jobInfo["buckets"] != nil {
				buckets := jobInfo["buckets"].([]interface{})
				for _, bucket := range buckets {
					endpointID := bucket.(map[string]interface{})["key"].(string)
					docCount := int(bucket.(map[string]interface{})["doc_count"].(float64))
					jobCounts[endpointID] = docCount
				}
			}
		}
	}
	runCounts := make(map[string]int)
	if runResponse[constants.AGGREGATION] != nil {
		aggsResult := runResponse[constants.AGGREGATION].(map[string]interface{})
		if aggsResult["jobs_per_endpoint"] != nil {
			jobInfo := aggsResult["jobs_per_endpoint"].(map[string]interface{})
			if jobInfo["buckets"] != nil {
				buckets := jobInfo["buckets"].([]interface{})
				for _, bucket := range buckets {
					endpointID := bucket.(map[string]interface{})["key"].(string)
					docCount := int(bucket.(map[string]interface{})["doc_count"].(float64))
					runCounts[endpointID] = docCount
				}
			}
		}
	}
	return jobCounts, runCounts, nil
}

func GetVersionAndPluginCount(replacements map[string]any, ctx context.Context) (map[string]interface{}, error) {
	osclient, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiToolVersionAndPluginCountQuery)
	if log.CheckErrorf(err, "could not replace json placeholders : %v", replacements) {
		return nil, err
	}
	modifiedJson := UpdateMustFilters(updatedJSON, replacements)
	response, err := searchResponse(modifiedJson, constants.CB_CI_TOOL_INSIGHT_INDEX, osclient)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure) {
		return nil, err
	}
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.PLUGIN_COUNT] != nil {
			completedRuns := aggsResult[constants.PLUGIN_COUNT].(map[string]interface{})
			if completedRuns[constants.VALUE] != nil {
				values := completedRuns[constants.VALUE].(map[string]interface{})
				return values, nil
			}
		}
	}
	return nil, err
}
