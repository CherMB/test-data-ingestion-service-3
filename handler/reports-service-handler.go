package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/calculi-corp/config"

	api "github.com/calculi-corp/api/go"
	auth "github.com/calculi-corp/api/go/auth"
	"github.com/calculi-corp/api/go/endpoint"

	pb "github.com/calculi-corp/api/go/vsm/report"
	"github.com/calculi-corp/log"
	opensearchconfig "github.com/calculi-corp/opensearch-config"
	"github.com/calculi-corp/reports-service/constants"
	"github.com/calculi-corp/reports-service/db"
	"github.com/calculi-corp/reports-service/exceptions"
	helper "github.com/calculi-corp/reports-service/helper"
	"github.com/calculi-corp/reports-service/internal"
	"github.com/calculi-corp/reports-service/models"
	"github.com/opensearch-project/opensearch-go"

	"github.com/calculi-corp/common/defines"
	client "github.com/calculi-corp/grpc-client"
	handler "github.com/calculi-corp/grpc-handler"
	healthchecks "github.com/calculi-corp/grpc-handler/pb"
	hostflags "github.com/calculi-corp/grpc-hostflags"
	srvauth "github.com/calculi-corp/grpc-server/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

var errInvalidRequest = "ReportService Request is invalid"
var errMissingRequiredField = "ReportService Request is missing required field %s"
var errInvalidDateFormatRequest = "Invalid Date format in ReportService Request %s and %s"
var errResourceNotInOrg = "resource %s does not belong to the organization"
var errCiToolNotInOrg = "ci tool %s does not belong to the organization"

var (
	getAllEndpoints = helper.GetAllEndpoints
)

const (
	endpointService            = "api.endpoint.UserPreferencesService"
	updateUserPreferenceMethod = "UpdateUserPreferences"
	getUserPreferencesMethod   = "GetUserPreferences"
	widgetIdField              = "widgetId"
	tenantIdField              = "tenantId"
	applicationIdField         = "applicationId"
	environmentIdField         = "environmentId"
	userIdField                = "userId"
	timeLayout                 = "2006-01-02 15:04:03"
	timeLayoutDateHistogram    = "2006-01-02"
	maxDaysInMonth             = 31
	maxDaysInWeek              = 7
)

type ReportsHandler struct {
	pb.UnimplementedReportServiceHandlerServer
	metrics          *handler.Map
	client           client.GrpcClient
	endpointClient   endpoint.EndpointServiceClient
	orgServiceClient auth.OrganizationsServiceClient
	rbacClt          auth.RBACServiceClient
}

// Defining reverse map to fetch identifiers from scanner names
var scannerNameReverseMap = map[string]string{
	"Anchore":                    "anchore",
	"Aquasec":                    "aquasec",
	"Gosec":                      "gosec",
	"Snyk SCA":                   "snyksca",
	"Snyk SAST":                  "snyksast",
	"Mend SCA":                   "mendsca",
	"Mend SAST":                  "mendsast",
	"Checkmarx":                  "checkmarx",
	"SonarQube":                  "sonarqube",
	"Trivy":                      "trivy",
	"Find Security Bugs":         "findsecbugs",
	"GitHub Security Scanner":    "githubsecurity",
	"TruffleHog S3":              "trufflehogs3",
	"TruffleHog Container":       "trufflehogcontainer",
	"Snyk Container":             "snykcontainer",
	"JFrog Xray":                 "jfrog-xray",
	"Stackhawk":                  "stackhawk",
	"ZAP":                        "zap",
	"Nexus IQ":                   "nexusiq",
	"SonarQube bundled":          "sonarqube-bundled",
	"TruffleHog SAST":            "trufflehogsast",
	"Sonatype (Nexus) Container": "nexusiq-scan-container",
}

// byPassDurationValidation map holds report/widget ids for which duration validation is not required
var byPassDurationValidation = map[string]string{
	constants.RUN_DETAILS_TEST_RESULTS:                                          "",
	constants.RUN_DETAILS_TOTAL_TEST_CASES:                                      "",
	constants.RUN_DETAILS_TEST_CASE_LOG:                                         "",
	constants.RUN_DETAILS_TEST_RESULTS_INDICATORS:                               "",
	constants.OPEN_FINDINGS_BY_SEVERITY_WIDGET_ID:                               "",
	constants.FINDINGS_IDENTIFIED_WIDGET_ID:                                     "",
	constants.SLA_BREACHES_BY_ASSET_TYPE:                                        "",
	constants.OPEN_FINDINGS_DISTRIBUTION_BY_CATEGORY:                            "",
	constants.OPEN_FINDINGS_BY_SECURITY_TOOL:                                    "",
	constants.OPEN_FINDINGS_DISTRIBUTION_BY_SECURITY_TOOL:                       "",
	constants.RISK_ACCEPTED_FALSE_POSITIVE_FINDINGS_WIDGET_ID:                   "",
	constants.SLA_BREACHES_BY_SEVERITY_WIDGET_ID:                                "",
	constants.OPEN_FINDINGS_BY_SLA_STATUS:                                       "",
	constants.OPEN_FINDINGS_BY_REVIEW_STATUS_WIDGET_ID:                          "",
	constants.APPLICATION_SLA_BREACHES_BY_SEVERITY_WIDGET_ID:                    "",
	constants.APPLICATION_OPEN_FINDINGS_BY_REVIEW_STATUS_WIDGET_ID:              "",
	constants.APPLICATION_OPEN_FINDINGS_BY_SLA_STATUS_WIDGET_ID:                 "",
	constants.APPLICATION_OPEN_FINDINGS_BY_SEVERITY_WIDGET_ID:                   "",
	constants.APPLICATION_RISK_ACCEPTED_FALSE_POSITIVE_FINDINGS_WIDGET_ID:       "",
	constants.APPLICATION_FINDINGS_IDENTIFIED_WIDGET_ID:                         "",
	constants.APPLICATION_OPEN_FINDINGS_BY_SECURITY_TOOL_WIDGET_ID:              "",
	constants.APPLICATION_OPEN_FINDINGS_DISTRIBUTION_BY_CATEGORY_WIDGET_ID:      "",
	constants.APPLICATION_OPEN_FINDINGS_DISTRIBUTION_BY_SECURITY_TOOL_WIDGET_ID: "",
	constants.APPLICATION_SLA_BREACHES_BY_ASSET_TYPE_WIDGET_ID:                  "",
	constants.APPLICATION_COMPONENTS_WITH_MOST_OPEN_FINDINGS_WIDGET_ID:          "",
}

func NewDefaultReportsHandler(client client.GrpcClient) (*ReportsHandler, error) {

	epSvcGrpcCltConn, err := client.Connect(hostflags.EndpointServiceHost())
	if err != nil {
		return nil, err
	}

	epSvcClt := endpoint.NewEndpointServiceClient(epSvcGrpcCltConn)

	orgSvcCltConn, err := client.Connect(hostflags.RbacServiceHost())
	if err != nil {
		return nil, err
	}

	orgSvcClt := auth.NewOrganizationsServiceClient(orgSvcCltConn)

	// Wiring RBAC grpc client
	rbacSvcGrpcCltConn, err := client.Connect(hostflags.RbacServiceHost())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to rbac service: %w", err)
	}
	rbacClt := auth.NewRBACServiceClient(rbacSvcGrpcCltConn)

	rah := &ReportsHandler{
		client:           client,
		endpointClient:   epSvcClt,
		orgServiceClient: orgSvcClt,
		rbacClt:          rbacClt,
	}

	rah.metrics = handler.NewMap(rah.Description().Name)

	return rah, nil
}

// NewReportsHandler initialize the metrics and the description
func NewReportsHandler(clt client.GrpcClient, endpointClt endpoint.EndpointServiceClient) *ReportsHandler {
	reportsHandler := &ReportsHandler{}
	reportsHandler.client = clt
	reportsHandler.endpointClient = endpointClt
	reportsHandler.metrics = handler.NewMap(reportsHandler.Description().Name)

	return reportsHandler
}

// NewReportsHandler with the rbac
func NewReportsHandlerWithRbac(clt client.GrpcClient, rbacClt auth.RBACServiceClient) *ReportsHandler {
	reportsHandler := &ReportsHandler{}
	reportsHandler.client = clt
	reportsHandler.rbacClt = rbacClt
	reportsHandler.metrics = handler.NewMap(reportsHandler.Description().Name)

	return reportsHandler
}

func (rah *ReportsHandler) Description() *handler.ServiceDesc {
	return &handler.ServiceDesc{Name: pb.ReportServiceHandler_ServiceDesc.ServiceName, ProtoDesc: pb.ReportServiceHandler_ServiceDesc}
}

// MetricMap returns the metrics for this service
func (rah *ReportsHandler) MetricMap() *handler.Map {
	return rah.metrics
}

func (rah *ReportsHandler) Healthy() error {
	return nil
}

// Dependencies returns a list of services which are dependents of pull request service
func (rah *ReportsHandler) Dependencies() []string {
	return []string{hostflags.DbService}
}

// HealthDependency checks the service's own health and all of its dependencies healths based on the received depth
func (rah *ReportsHandler) HealthDependency(depth int32, task string) []*healthchecks.ServiceHealthResponse {
	return handler.HealthCheck.HealthDependency(rah, depth, task, rah.Dependencies())
}

func (rah *ReportsHandler) Stop() {
	log.Debugf("Closing reports-service clients")
	rah.client.Close()
}

func (rah *ReportsHandler) verifyUserPreferencesRequest(ctx context.Context, userId string) error {
	if userId == "" {
		return defines.ErrIDMissing
	}

	resp, err := rah.rbacClt.GetContextUserId(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	if userId != resp.GetUserId() || resp.GetUserId() == defines.SystemUser {
		log.Warnf("user id %s is not authorized to update user dashboard preferences", resp.GetUserId())
		return status.Error(codes.PermissionDenied, "not authorized")
	}
	return nil
}

func (rah *ReportsHandler) BuildReport(ctx context.Context, req *pb.ReportServiceRequest) (*pb.ReportServiceResponse, error) {
	log.Debugf("BuildReport started - Components in request %v for widget %s", req.Component, req.WidgetId)
	response := pb.ReportServiceResponse{
		Status:  pb.Status_success,
		Error:   "",
		Widget:  nil,
		Message: "",
	}

	log.Debugf("BuildReport %v", req)
	err := ValidateDataRequest(req, ctx, rah.endpointClient)
	if err != nil {
		// convert error to GRPC status error so that api-gateway can handle it
		return nil, status.Errorf(codes.InvalidArgument, exceptions.ErrReportServiceReqValidationFailure, err.Error())
	}

	// validate access if the request is intended for CI Insights report
	if err := helper.ValidateCIInsightsReportAccess(ctx, rah.rbacClt, req); err != nil {
		return nil, err
	}

	components := []string{}

	// 1. Using user id - fetch all the sub orgs under the org

	// 2.  helper.GetOrganisationServices - fetch all the components for the parent org

	// 3. Filter the components which are mapping with the sub org id from #1.sub org id == #2. organisationId

	if req.UserId != "" {

		if len(req.Component) > 1 || (len(req.Component) == 1 && req.Component[0] != "All") {
			log.Debugf(exceptions.DebugApplyingComponentsFilter, req.Component)
		} else {
			if req.OrgId == req.SubOrgId {
				getOrganizationByIdResponse, err := helper.GetOrganisationsById(ctx, rah.orgServiceClient, req.OrgId, req.UserId)
				if err != nil {
					// convert error to GRPC status error so that api-gateway can handle it
					return nil, status.Errorf(codes.FailedPrecondition, exceptions.ErrFetchingSubOrgForOrg, req.OrgId, err.Error())
				}
				serviceResponse, err := helper.GetOrganisationServices(ctx, rah.client, req.OrgId)
				if err != nil {
					return nil, status.Errorf(codes.FailedPrecondition, exceptions.ErrFetchingServicesForOrg, req.OrgId, err.Error())
				}
				for _, service := range serviceResponse.GetService() {
					if req.OrgId == service.OrganizationId {
						components = append(components, service.Id)
					}
				}
				helper.GetComponentsRecursively(getOrganizationByIdResponse.GetOrganization(), serviceResponse, &components)
				req.Component = components
			} else if req.OrgId != req.SubOrgId && req.SubOrgId != "" {
				getOrganizationByIdResponse, err := helper.GetOrganisationsById(ctx, rah.orgServiceClient, req.SubOrgId, req.UserId)

				if err != nil {
					// convert error to GRPC status error so that api-gateway can handle it
					return nil, status.Errorf(codes.FailedPrecondition, "Exception while fetching sub orgs for sub org : %s : %s", req.SubOrgId, err.Error())
				}

				serviceResponse, err := getOrganisationServices(ctx, rah.client, req.SubOrgId)

				if err != nil {

					return nil, status.Errorf(codes.FailedPrecondition, exceptions.ErrFetchingServicesForSubOrg, req.SubOrgId, err.Error())
				}

				for _, service := range serviceResponse.GetService() {

					if req.SubOrgId == service.OrganizationId {
						components = append(components, service.Id)
					}
				}

				helper.GetComponentsRecursively(getOrganizationByIdResponse.GetOrganization(), serviceResponse, &components)
				req.Component = components
			}
		}
	} else {
		if req.OrgId != req.SubOrgId && req.SubOrgId != "" {
			serviceResponse, err := helper.GetOrganisationServices(ctx, rah.client, req.SubOrgId)
			if err != nil {
				return nil, status.Errorf(codes.FailedPrecondition, exceptions.ErrFetchingServicesForSubOrg, req.SubOrgId, err.Error())
			}
			if len(serviceResponse.GetService()) > 0 {
				for i := 0; i < len(serviceResponse.GetService()); i++ {
					service := serviceResponse.GetService()[i]
					if len(req.Component) > 1 || (len(req.Component) == 1 && req.Component[0] != "All") {
						for _, comp := range req.Component {
							if comp == service.Id {
								components = append(components, service.Id)
							}
						}
					} else {
						components = append(components, service.Id)
					}
				}
			}
			req.Component = components
		}
	}

	if len(req.Component) == 0 && !strings.HasPrefix(req.WidgetId, "ci") {
		response.Message = constants.NO_DATA_FOUND
		return &response, nil
	}

	var aggrBy string           // aggregation interval for date histogram - week when duration is a month; day when duration is week
	var duration string         // duration type - week or month
	var normalizeMonthFlag bool // flag to check if the final dates in the x axis for charts should be normalized, since aggr. by week interval exceeds month range
	var weekOrDayInMilli int64  // duration for a week or day in milliseconds - used in the query for the Flow Work Load widget
	var commitTitle string
	_, byPassDurationValidation := byPassDurationValidation[req.WidgetId]
	if !byPassDurationValidation {
		if req.DurationType == pb.DurationType_CURRENT_WEEK ||
			req.DurationType == pb.DurationType_PREVIOUS_WEEK ||
			req.DurationType == pb.DurationType_TWO_WEEKS_BACK ||
			req.DurationType == pb.DurationType_LAST_7_DAYS {
			duration = constants.DURATION_WEEK
			aggrBy = constants.DURATION_DAY
			normalizeMonthFlag = false
			weekOrDayInMilli = constants.DURATION_DAY_IN_MILLISEC // 24 hr
			commitTitle = constants.DAILY_COMMITS_OR_ACTIVE_DEVS
		} else if req.DurationType == pb.DurationType_CURRENT_MONTH ||
			req.DurationType == pb.DurationType_PREVIOUS_MONTH ||
			req.DurationType == pb.DurationType_TWO_MONTHS_BACK ||
			req.DurationType == pb.DurationType_LAST_30_DAYS {
			duration = constants.DURATION_MONTH
			aggrBy = constants.DURATION_WEEK
			if aggr, ok := db.CustomAggrMap[req.WidgetId]; ok {
				aggrBy = aggr
			}
			normalizeMonthFlag = true
			weekOrDayInMilli = constants.DURATION_WEEK_IN_MILLISEC // week
			commitTitle = constants.WEEKLY_COMMITS_OR_ACTIVE_DEVS
		} else if req.DurationType == pb.DurationType_LAST_90_DAYS {
			duration = constants.DURATION_MONTH
			aggrBy = constants.DURATION_MONTH
			normalizeMonthFlag = true
			weekOrDayInMilli = constants.DURATION_30_DAY_IN_MILLISEC //30day
			commitTitle = constants.WEEKLY_COMMITS_OR_ACTIVE_DEVS
		} else if req.DurationType == pb.DurationType_LAST_YEAR {
			duration = constants.DURATION_YEAR
			aggrBy = constants.DURATION_MONTH
			normalizeMonthFlag = true
		} else if req.DurationType == pb.DurationType_CUSTOM_RANGE {
			d, err := getDateDiffInDays(req.StartDate, req.EndDate)
			if err != nil {
				return nil, err
			}
			if d > maxDaysInMonth {
				duration = constants.DURATION_MONTH
				aggrBy = constants.DURATION_MONTH
				normalizeMonthFlag = true
				weekOrDayInMilli = constants.DURATION_30_DAY_IN_MILLISEC
				commitTitle = constants.WEEKLY_COMMITS_OR_ACTIVE_DEVS
			} else if d > maxDaysInWeek {
				duration = constants.DURATION_MONTH
				aggrBy = constants.DURATION_WEEK
				normalizeMonthFlag = true
				weekOrDayInMilli = constants.DURATION_WEEK_IN_MILLISEC
				commitTitle = constants.WEEKLY_COMMITS_OR_ACTIVE_DEVS
			} else {
				duration = constants.DURATION_WEEK
				aggrBy = constants.DURATION_DAY
				normalizeMonthFlag = false
				weekOrDayInMilli = constants.DURATION_DAY_IN_MILLISEC
				commitTitle = constants.DAILY_COMMITS_OR_ACTIVE_DEVS
			}
		}
	}

	//Get the filters from request and create map of replacement placeholders
	replacements := map[string]any{
		"startDate":        req.StartDate,
		"endDate":          req.EndDate,
		"orgId":            req.OrgId,
		"subOrgId":         req.SubOrgId,
		"component":        req.Component,
		"aggrBy":           aggrBy,
		"duration":         duration,
		"weekOrDayInMilli": weekOrDayInMilli,
		"targetEnv":        req.Environment,
		"application":      req.Application,
		"ciToolId":         req.CiToolId,
		"ciToolType":       req.CiToolType,
		"sortBy":           req.SortBy,
		"filterType":       req.FilterType,
		"viewOption":       req.ViewOption,
		"timeFormat":       req.TimeFormat,
		"timeZone":         req.TimeZone,
	}

	if len(req.Tools) > 0 {
		replacements["tools"] = req.Tools
	}

	if len(req.Severities) > 0 {
		replacements["severities"] = ConvertSeveritiesToStrings(req.Severities)
	}

	if IsValidSlaStatus(req.Sla) {
		replacements["sla"] = req.Sla
	}

	if req.TimeZone == "" {
		replacements["timeZone"] = constants.LOCATION_EUROPE_OR_LONDON
	}

	if req.GetTimeFormat() == "" {
		replacements["timeFormat"] = constants.DURATION_TWELVE_HOURS
	} else {
		replacements["timeFormat"] = req.GetTimeFormat()
	}

	if len(req.Branch) > 0 {
		replacements[constants.REQUEST_BRANCH] = req.Branch
	}

	//replacements for placeholders inside JOLT specs
	replacementsSpec := map[string]any{
		"normalizeMonthInSpec": "@x",
		"commitTitle":          commitTitle,
	}

	if !byPassDurationValidation {

		dateBydurationType := models.CalculateDateBydurationType{
			StartDateStr:       req.StartDate,
			EndDateStr:         req.EndDate,
			DurationType:       req.DurationType,
			Replacements:       replacements,
			ReplacementsSpec:   replacementsSpec,
			NormalizeMonthFlag: normalizeMonthFlag,
			IsComputFlag:       false,
			CurrentTime:        time.Now().UTC(),
		}
		calculateDateBydurationType(dateBydurationType)
	}

	_, ok := db.PaginationReportMap[req.WidgetId]
	if ok {
		startTime := time.Now()
		parentWidget, err := internal.CreateWidget(req.WidgetId, nil, replacements, replacementsSpec, nil)
		if err != nil {
			response.Status = pb.Status_error
			response.Error = fmt.Sprintf("ReportServiceRequest failed to get parent widget : %s", err.Error())
			return &response, nil
		}
		log.Debugf("Time took to get parent widget for widget %s : %v in milliseconds", req.WidgetId, time.Since(startTime).Milliseconds())
		startTime = time.Now()
		baseData, dataMap, err := internal.ExecuteMultiPageBaseFunction("", req.WidgetId, replacements, ctx)
		if err != nil {
			response.Status = pb.Status_error
			response.Error = fmt.Sprintf("ReportServiceRequest failed to get base pagination data : %s", err.Error())
			return &response, nil
		}
		log.Debugf("Time took to get subwidget filter data for widget %s : %v in milliseconds", req.WidgetId, time.Since(startTime).Milliseconds())
		startTime = time.Now()
		count := 1
		filterCount, ok := db.PaginationReportFilterCountMap[req.WidgetId]
		if ok {
			count = filterCount
		}
		if len(dataMap) > 0 {
			for key, value := range dataMap {
				replacementsSpec[key] = value
			}
		}
		size := len(baseData)
		for index := 0; index < size; index = index + count {
			endIndex := index + count
			var filterData []string
			if endIndex > size {
				filterData = baseData[index:size]
			} else {
				filterData = baseData[index:endIndex]
			}

			subReportWidget := models.GetSubReportWidget{
				FilterData:       filterData,
				BaseData:         baseData,
				Req:              req,
				Replacements:     replacements,
				Ctx:              ctx,
				ReplacementsSpec: replacementsSpec,
				ParentWidget:     parentWidget,
			}

			// Get sub report widget for each filter value
			getSubReportWidget(subReportWidget, rah)
		}
		if len(parentWidget.Content) > 0 {
			response.Widget = parentWidget
		} else {
			response.Message = constants.NO_DATA_FOUND
		}
		log.Debugf("Total time took to get all subwidget for %s : %v in milliseconds", req.WidgetId, time.Since(startTime).Milliseconds())
	} else {
		//Get all queries for the widget, apply filters, run and get the response for OpenSearch
		d, prev_d, fd, err := internal.GetData(req.WidgetId, replacements, nil, ctx, rah.client, rah.endpointClient)
		if err != nil {
			response.Status = pb.Status_error
			response.Error = fmt.Sprintf("ReportServiceRequest failed to get data: %s", err.Error())
			return &response, nil
		}

		if helper.IsResponseEmpty(d) && fd == nil && req.WidgetId != "cs9" {
			response.Message = constants.NO_DATA_FOUND
		} else {
			// append previous duration data with current data before building the widget
			for key, val := range prev_d {
				d[key] = val
			}
			startTime := time.Now()
			// Build the widget
			w, err := internal.CreateWidget(req.WidgetId, d, replacements, replacementsSpec, fd)
			log.Debugf("Time took to create widget for id %s : %v in milliseconds", req.WidgetId, time.Since(startTime).Milliseconds())
			if err != nil {
				if err == db.ErrNoDataFound {
					response.Message = constants.NO_DATA_FOUND
					return &response, nil
				} else {
					response.Status = pb.Status_error
					response.Error = db.ErrInternalServer.Error()
					return &response, nil
				}

			}
			response.Widget = w
		}
	}
	log.Debugf("BuildReport completed for widget : %s for component - %s", req.WidgetId, req.Component)
	return &response, nil
}

// Convert a slice of Severity enums into a slice of strings
func ConvertSeveritiesToStrings(severities []pb.Severity) []string {
	// Create a new slice with the same length as the input slice
	strings := make([]string, len(severities))

	// Iterate over the severities and convert each one to its string representation
	for i, severity := range severities {
		strings[i] = severity.String()
	}

	return strings
}

// Define a helper function to validate the SlaStatus value
func IsValidSlaStatus(status pb.SlaStatus) bool {
	switch status {
	case pb.SlaStatus_BREACHED_SLA, pb.SlaStatus_WITHIN_SLA:
		return true
	default:
		return false
	}
}

func (rah *ReportsHandler) BuildComponentComparisonReport(ctx context.Context, req *pb.ComponentComparisonRequest) (*pb.ComponentComparisonResponse, error) {
	log.Debugf("BuildComponentComparisonReport started - for widget %s", req.WidgetId)
	response := pb.ComponentComparisonResponse{
		Status:                  pb.Status_success,
		Error:                   "",
		ComponentComparisonData: nil,
		Message:                 "",
	}

	log.Debugf("BuildComponentComparisonReport %v", req)
	err := ValidateComponentComparisonDataRequest(req)
	if err != nil {
		// convert error to GRPC status error so that api-gateway can handle it
		return nil, status.Errorf(codes.InvalidArgument, "ComponentComparisonRequest validation failed: %s", err.Error())
	}

	err = rah.verifyUserPreferencesRequest(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	var environmentId string
	if req.WidgetId == "deployment-frequency-compare" || req.WidgetId == "deployment-lead-time-compare" || req.WidgetId == "failure-rate-compare" || req.WidgetId == "mttr-compare" {
		environmentId = req.Environment
		if environmentId == "" {
			return nil, status.Errorf(codes.InvalidArgument, "ComponentComparisonRequest validation failed: Environment id is missing in request : ")
		} else {
			if req.SubOrgId != "" {
				contributionIds := []string{constants.ENVIRONMENT_ENDPOINT}
				endPointsResponse, err := helper.GetAllEndpoints(ctx, rah.endpointClient, req.SubOrgId, contributionIds, true)
				if log.CheckErrorf(err, exceptions.ErrFetchingEndpoints) {
					return nil, exceptions.GetExceptionByCode(exceptions.ErrEndpointAPIFailure)
				} else {
					endpoints := endPointsResponse.Endpoints
					log.Debugf(exceptions.DebugEndpointList, endpoints)
					if len(endpoints) > 0 {
						for _, endpoint := range endpoints {
							if endpoint.Id == environmentId {
								req.Environment = endpoint.Name
								break
							}
						}
					}
				}
			}
		}

		if req.Environment == environmentId {
			return nil, status.Errorf(codes.InvalidArgument, "ComponentComparisonRequest validation failed: Invalid Environment id in request : %s", req.Environment)
		}
	}

	var aggrBy string           // currently only being used in the Flow Work Load query as it uses a custom implementaion instead of the native date histogram aggr, and none of the other widgets in Component Comparison warrant the use of a date histogram
	var duration string         // duration type - week or month
	var normalizeMonthFlag bool // flag to check if the final dates in the x axis for charts should be normalized, since aggr. by week interval exceeds month range
	var weekOrDayInMilli int64  // duration for a week or day in milliseconds - used in the query for the Flow Work Load widget
	var targetEnv string
	var commitTitle string
	if req.DurationType == pb.DurationType_CURRENT_WEEK ||
		req.DurationType == pb.DurationType_PREVIOUS_WEEK ||
		req.DurationType == pb.DurationType_TWO_WEEKS_BACK ||
		req.DurationType == pb.DurationType_LAST_7_DAYS {
		duration = "week"
		normalizeMonthFlag = false
		weekOrDayInMilli = constants.DURATION_DAY_IN_MILLISEC
		commitTitle = constants.DAILY_COMMITS_OR_ACTIVE_DEVS
	} else if req.DurationType == pb.DurationType_CURRENT_MONTH ||
		req.DurationType == pb.DurationType_PREVIOUS_MONTH ||
		req.DurationType == pb.DurationType_TWO_MONTHS_BACK ||
		req.DurationType == pb.DurationType_LAST_30_DAYS {
		duration = "month"
		normalizeMonthFlag = true
		weekOrDayInMilli = constants.DURATION_WEEK_IN_MILLISEC
		commitTitle = constants.WEEKLY_COMMITS_OR_ACTIVE_DEVS
	} else if req.DurationType == pb.DurationType_LAST_90_DAYS {
		duration = "month"
		normalizeMonthFlag = true
		weekOrDayInMilli = constants.DURATION_30_DAY_IN_MILLISEC
		commitTitle = constants.WEEKLY_COMMITS_OR_ACTIVE_DEVS
	} else if req.DurationType == pb.DurationType_CUSTOM_RANGE {
		d, err := getDateDiffInDays(req.StartDate, req.EndDate)
		if err != nil {
			return nil, err
		}
		if d > maxDaysInMonth {
			duration = "month"
			normalizeMonthFlag = true
			weekOrDayInMilli = constants.DURATION_30_DAY_IN_MILLISEC
			commitTitle = constants.WEEKLY_COMMITS_OR_ACTIVE_DEVS
		} else if d > maxDaysInWeek {
			duration = "month"
			normalizeMonthFlag = true
			weekOrDayInMilli = constants.DURATION_WEEK_IN_MILLISEC
			commitTitle = constants.WEEKLY_COMMITS_OR_ACTIVE_DEVS
		} else {
			duration = "week"
			normalizeMonthFlag = false
			weekOrDayInMilli = constants.DURATION_DAY_IN_MILLISEC
			commitTitle = constants.DAILY_COMMITS_OR_ACTIVE_DEVS
		}
	}
	targetEnv = req.Environment
	aggrBy = "duration"

	//Get the filters from request and create map of replacement placeholders
	replacements := map[string]any{
		"startDate":        req.StartDate,
		"endDate":          req.EndDate,
		"orgId":            req.OrgId,
		"subOrgId":         req.SubOrgId,
		"aggrBy":           aggrBy,
		"duration":         duration,
		"weekOrDayInMilli": weekOrDayInMilli,
		"targetEnv":        targetEnv,
		"timeZone":         req.TimeZone,
		"userId":           req.UserId,
	}

	if req.TimeZone == "" {
		replacements["timeZone"] = constants.LOCATION_EUROPE_OR_LONDON
	}

	//replacements for placeholders inside JOLT specs
	replacementsSpec := map[string]any{
		"normalizeMonthInSpec": "@x",
		"commitTitle":          commitTitle,
	}

	organization, components, err := helper.FetchOrganizationAndServices(ctx, rah.client, rah.orgServiceClient, req.OrgId, req.UserId)
	replacements["component"] = components
	if log.CheckErrorf(err, "Error fetching organization and services in BuildComponentComparisonReport : ") {
		return nil, err
	}

	if organization != nil && organization.Name != "" {
		replacements["orgName"] = organization.Name
	} else {
		response.Status = pb.Status_error
		response.Error = fmt.Sprintf("ComponentComparisonRequest failed to get organization: %s", err.Error())
		return &response, nil
	}
	dateBydurationType := models.CalculateDateBydurationType{
		StartDateStr:       req.StartDate,
		EndDateStr:         req.EndDate,
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   replacementsSpec,
		NormalizeMonthFlag: normalizeMonthFlag,
		IsComputFlag:       false,
		CurrentTime:        time.Now().UTC(),
	}
	calculateDateBydurationType(dateBydurationType)

	//Get all queries for the widget, apply filters, run and get the response for OpenSearch
	d, fd, err := internal.GetComponentComparisonData(req.WidgetId, replacements, ctx, rah.client, rah.endpointClient)
	if err != nil {
		response.Status = pb.Status_error
		response.Error = fmt.Sprintf("ComponentComparisonRequest failed to get data: %s", err.Error())
		return &response, nil
	}

	if helper.IsResponseEmpty(d) && len(fd) == 0 {
		response.Message = constants.NO_DATA_FOUND
	} else {

		startTime := time.Now()
		// Build the widget
		w, err := internal.CreateComponentComparisonWidget(req.WidgetId, d, fd, replacements, organization)
		log.Debugf("Time took to create widget for id %s : %v in milliseconds", req.WidgetId, time.Since(startTime).Milliseconds())
		if err != nil {
			response.Status = pb.Status_error
			response.Error = fmt.Sprintf("ComponentComparisonRequest failed to create widget: %s", err.Error())
			return &response, nil
		}
		w.BreadCrumbTitle = organization.Name
		response.ComponentComparisonData = w
	}

	log.Debugf("BuildComponentComparisonReport completed for widget : %s ", req.WidgetId)
	return &response, nil

}

func getSubReportWidget(subReportWidget models.GetSubReportWidget, rah *ReportsHandler) {
	replaceMentKey, ok := db.PaginationReportFilterMap[subReportWidget.Req.WidgetId]
	var baseDataReplacements map[string]any
	if ok {
		baseDataReplacements = make(map[string]any)
		baseDataReplacements[replaceMentKey] = subReportWidget.FilterData
	}
	subReportWidgetId, ok := db.PaginationSubReportMap[subReportWidget.Req.WidgetId]
	startTime := time.Now()
	if ok {
		if constants.REPLACE_HEADER_WIDGET == subReportWidgetId {
			subReportWidget.Replacements[constants.REPLACE_HEADER_KEY] = subReportWidget.FilterData[0]
		} else if constants.REPLACE_SUMMARY_WIDGET == subReportWidgetId {
			subReportWidget.Replacements[constants.REPLACE_HEADER_KEY] = subReportWidget.FilterData[0]
		}
		startTime1 := time.Now()
		d, prev_d, fd, err := internal.GetData(subReportWidgetId, subReportWidget.Replacements, baseDataReplacements, subReportWidget.Ctx, rah.client, nil)
		log.Debugf("Time took to get data for widget %s : %v in milliseconds", subReportWidgetId, time.Since(startTime1).Milliseconds())
		if err != nil {
			log.Error("Error while getting sub widget data", err)
		} else {
			if !(helper.IsResponseEmpty(d) && len(fd) == 0) {
				for key, val := range prev_d {
					d[key] = val
				}
				startTime2 := time.Now()
				w, err := internal.CreateWidget(subReportWidgetId, d, subReportWidget.Replacements, subReportWidget.ReplacementsSpec, fd)
				log.Debugf("Time took to CreateWidget for widget %s : %v in milliseconds", subReportWidgetId, time.Since(startTime2).Milliseconds())
				if err != nil {
					log.Error("Error while creating sub widget ", err)
				} else {
					subReportWidget.ParentWidget.Content = append(subReportWidget.ParentWidget.Content, w.Content...)
				}
			}
		}
	}
	log.Debugf("Total time took for sub widget %s : %v in milliseconds", subReportWidgetId, time.Since(startTime).Milliseconds())
}

func (rah *ReportsHandler) FetchReportData(ctx context.Context, req *pb.ReportServiceRequest) (*pb.ReportServiceResponse, error) {
	return rah.GetReportData(ctx, req)
}

func (rah *ReportsHandler) GetReportData(ctx context.Context, req *pb.ReportServiceRequest) (*pb.ReportServiceResponse, error) {
	res := &pb.ReportServiceResponse{}
	var err error
	startTime := time.Now()
	if strings.HasPrefix(req.WidgetId, "d") {
		environmentId := req.Environment
		if environmentId == "" {
			res.Message = constants.NO_DATA_FOUND
			return res, nil
		} else {
			if req.SubOrgId != "" {
				contributionIds := []string{constants.ENVIRONMENT_ENDPOINT}
				endPointsResponse, err := helper.GetAllEndpoints(ctx, rah.endpointClient, req.SubOrgId, contributionIds, true)
				if log.CheckErrorf(err, exceptions.ErrFetchingEndpoints) {
					return nil, exceptions.GetExceptionByCode(exceptions.ErrEndpointAPIFailure)
				} else {
					endpoints := endPointsResponse.Endpoints
					log.Debugf(exceptions.DebugEndpointList, endpoints)
					if len(endpoints) > 0 {
						for _, endpoint := range endpoints {
							if endpoint.Id == environmentId {
								req.Environment = endpoint.Name
								break
							}
						}
					}
				}
			}
			res, err = rah.BuildReport(ctx, req)
			if err != nil {
				return nil, err
			}
		}
	} else if req.WidgetId != "" {
		//else if strings.HasPrefix(req.WidgetId, "cs") || strings.HasPrefix(req.WidgetId, "ci") {
		// Disabling call to BuildComputedReport for now.
		res, err = rah.BuildReport(ctx, req)
		if err != nil {
			return nil, err
		}
	} else {
		log.Debugf("BuildComputedReport started - Components in request %v for widget %s", req.Component, req.WidgetId)
		res, err = rah.BuildComputedReport(ctx, req)
		if err != nil {
			return nil, err
		}
		log.Debugf("BuildComputedReport completed - Components in request %v for widget %s", req.Component, req.WidgetId)
	}
	log.Infof("TIME TAKEN GetReportData : %v Milliseconds, Comp : %v , Widget : %s", time.Since(startTime).Milliseconds(), req.Component, req.WidgetId)
	return res, nil
}

func isValidDate(date string, format string) bool {
	if len(date) > 0 {
		_, err := time.Parse(format, date)
		return err == nil
	} else {
		return false
	}
}

func calculateDateBydurationType(dateBydurationType models.CalculateDateBydurationType) {
	weekdayVal := int(dateBydurationType.CurrentTime.Weekday())
	dayDifference := weekdayVal - 1
	var startDate, endDate time.Time
	if isValidDate(dateBydurationType.StartDateStr, timeLayoutDateHistogram) && isValidDate(dateBydurationType.EndDateStr, timeLayoutDateHistogram) {
		dateBydurationType.Replacements["startDate"] = dateBydurationType.StartDateStr + " 00:00:00"
		dateBydurationType.Replacements["endDate"] = dateBydurationType.EndDateStr + " 23:59:59"
		dateBydurationType.Replacements["dateHistogramMin"] = dateBydurationType.StartDateStr
		dateBydurationType.Replacements["dateHistogramMax"] = dateBydurationType.EndDateStr
		startDate, err := time.Parse(timeLayoutDateHistogram, dateBydurationType.StartDateStr)
		if err == nil {
			dateBydurationType.Replacements["startDateInMillis"] = getStartOfTheDay(startDate).UnixMilli()
		}
		endDate, err := time.Parse(timeLayoutDateHistogram, dateBydurationType.EndDateStr)
		if err == nil {
			dateBydurationType.Replacements["endDateInMillis"] = getEndOfTheDay(endDate).UnixMilli()
		}
		if dateBydurationType.NormalizeMonthFlag {
			dateBydurationType.ReplacementsSpec["normalizeMonthInSpec"] = startDate.Format(timeLayoutDateHistogram)
		}
	} else {
		switch dateBydurationType.DurationType.Number() {
		case 1:
			if dayDifference > 0 {
				startIndex := -(dayDifference)
				start := dateBydurationType.CurrentTime.AddDate(0, 0, startIndex)
				startDate = getStartOfTheDay(start)
				endDate = getEndOfTheDay(start.AddDate(0, 0, 6))
			} else if dayDifference == 0 {
				startDate = getStartOfTheDay(dateBydurationType.CurrentTime)
				endDate = getEndOfTheDay(dateBydurationType.CurrentTime.AddDate(0, 0, 6))
			} else {
				startDate = getStartOfTheDay(dateBydurationType.CurrentTime.AddDate(0, 0, -6))
				endDate = getEndOfTheDay(dateBydurationType.CurrentTime)
			}
		case 2:
			if dayDifference > 0 {
				startIndex := -(dayDifference + 7)
				start := dateBydurationType.CurrentTime.AddDate(0, 0, startIndex)
				startDate = getStartOfTheDay(start)
				endDate = getEndOfTheDay(start.AddDate(0, 0, 6))
			} else if dayDifference == 0 {
				startDate = getStartOfTheDay(dateBydurationType.CurrentTime.AddDate(0, 0, -7))
				endDate = getEndOfTheDay(dateBydurationType.CurrentTime.AddDate(0, 0, -1))
			} else {
				startDate = getStartOfTheDay(dateBydurationType.CurrentTime.AddDate(0, 0, -13))
				endDate = getEndOfTheDay(dateBydurationType.CurrentTime.AddDate(0, 0, -7))
			}
		case 3:
			if dayDifference > 0 {
				startIndex := -(dayDifference + 14)
				start := dateBydurationType.CurrentTime.AddDate(0, 0, startIndex)
				startDate = getStartOfTheDay(start)
				endDate = getEndOfTheDay(start.AddDate(0, 0, 6))
			} else if dayDifference == 0 {
				startDate = getStartOfTheDay(dateBydurationType.CurrentTime.AddDate(0, 0, -14))
				endDate = getEndOfTheDay(dateBydurationType.CurrentTime.AddDate(0, 0, -8))
			} else {
				startDate = getStartOfTheDay(dateBydurationType.CurrentTime.AddDate(0, 0, -20))
				endDate = getEndOfTheDay(dateBydurationType.CurrentTime.AddDate(0, 0, -14))
			}
		case 4:
			currentYear, currentMonth, _ := dateBydurationType.CurrentTime.Date()
			currentLocation := dateBydurationType.CurrentTime.Location()
			startDate = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
			endDate = getEndOfTheDay(startDate.AddDate(0, 1, -1))
		case 5:
			year, currentMonth, _ := dateBydurationType.CurrentTime.Date()
			currentLocation := dateBydurationType.CurrentTime.Location()
			month := currentMonth - 1
			if currentMonth == time.January {
				month = time.December
				year = year - 1
			}
			startDate = time.Date(year, month, 1, 0, 0, 0, 0, currentLocation)
			endDate = getEndOfTheDay(startDate.AddDate(0, 1, -1))
		case 6:
			year, currentMonth, _ := dateBydurationType.CurrentTime.Date()
			currentLocation := dateBydurationType.CurrentTime.Location()
			month := currentMonth - 2
			if currentMonth == time.January {
				month = time.November
				year = year - 1
			} else if currentMonth == time.February {
				month = time.December
				year = year - 1
			}
			startDate = time.Date(year, month, 1, 0, 0, 0, 0, currentLocation)
			endDate = getEndOfTheDay(startDate.AddDate(0, 1, -1))
		default:
			log.Debugf("Duration type is not supported.")
			startDate, _ = time.Parse(timeLayout, dateBydurationType.Replacements["startDate"].(string))
			endDate, _ = time.Parse(timeLayout, dateBydurationType.Replacements["startDate"].(string))
		}
		// start and end dates for query filters
		dateBydurationType.Replacements["startDate"] = startDate.Format(timeLayout)
		dateBydurationType.Replacements["startDateInMillis"] = startDate.UnixMilli()
		dateBydurationType.Replacements["endDate"] = endDate.Format(timeLayout)
		dateBydurationType.Replacements["endDateInMillis"] = endDate.UnixMilli()

		if dateBydurationType.IsComputFlag {
			dateBydurationType.Replacements["startDate"] = startDate.Format(timeLayoutDateHistogram) + " 00:00:00"
			dateBydurationType.Replacements["endDate"] = endDate.Format(timeLayoutDateHistogram) + " 23:59:59"
		}

		//date histogram bounds
		//minimum
		dateBydurationType.Replacements["dateHistogramMin"] = startDate.Format(timeLayoutDateHistogram)
		//maximum
		if dateBydurationType.DurationType == pb.DurationType_CURRENT_MONTH ||
			dateBydurationType.DurationType == pb.DurationType_CURRENT_WEEK {
			currentTime := time.Now()
			dateBydurationType.Replacements["dateHistogramMax"] = currentTime.Format(timeLayoutDateHistogram)
		} else {
			dateBydurationType.Replacements["dateHistogramMax"] = endDate.Format(timeLayoutDateHistogram)
		}

		//normalizing all x-axis labels (dates) to the month when aggregating by week
		if dateBydurationType.NormalizeMonthFlag {
			dateBydurationType.ReplacementsSpec["normalizeMonthInSpec"] = startDate.Format(timeLayoutDateHistogram)
		}
	}
}

func getStartOfTheDay(startDate time.Time) time.Time {
	year, month, day := startDate.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, startDate.Location())
}

func getEndOfTheDay(endDate time.Time) time.Time {
	year, month, day := endDate.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, endDate.Location())
}

func ValidateDataRequest(req *pb.ReportServiceRequest, ctx context.Context, epClt endpoint.EndpointServiceClient) error {
	log.Debugf(exceptions.DebugRequestList, req)
	if req.WidgetId == "" {
		return fmt.Errorf(errMissingRequiredField, widgetIdField)
	}

	if req.OrgId == "" {
		return fmt.Errorf(errMissingRequiredField, tenantIdField)
	}

	if len(req.Component) > 0 && (len(req.Component) == 1 && req.Component[0] != constants.ALL) {
		for _, resource := range req.Component {
			if !helper.IsResourceInOrg(resource, req.OrgId) {
				return fmt.Errorf(errResourceNotInOrg, resource)
			}
		}
	}

	if _, exists := db.FilterWidgetMap["application_dashboard_widgets"][req.WidgetId]; exists && len(req.Application) == 0 {
		return fmt.Errorf(errMissingRequiredField, applicationIdField)
	}

	if len(req.Application) > 0 {
		for _, app := range req.Application {
			if !helper.IsResourceInOrg(app, req.OrgId) {
				return fmt.Errorf(errResourceNotInOrg, app)
			}
		}

		if req.Environment == "" {
			return fmt.Errorf(errMissingRequiredField, environmentIdField)
		}

		contributionIds := []string{constants.ENVIRONMENT_ENDPOINT}
		orgIdToUse := req.SubOrgId
		if orgIdToUse == "" {
			orgIdToUse = req.OrgId
		}
		endPointsResponse, err := helper.GetAllEndpoints(ctx, epClt, orgIdToUse, contributionIds, true)
		if err != nil {
			log.Errorf(err, exceptions.ErrFetchingEndpoints)
			return fmt.Errorf("failed to fetch environments for validation: %v", err)
		}
		found := false
		for _, endpoint := range endPointsResponse.GetEndpoints() {
			if endpoint.Id == req.Environment {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid environment ID '%s' for applications, ensure the environment exists and is valid within the selected applications", req.Environment)
		}

	}

	if req.OrgId != req.SubOrgId && req.SubOrgId != "" {
		if !helper.IsResourceInOrg(req.SubOrgId, req.OrgId) {
			return fmt.Errorf(errResourceNotInOrg, req.SubOrgId)
		}
	}

	if req.CiToolId != "" {
		contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
		endPointsResponse, _ := getAllEndpoints(ctx, epClt, req.GetSubOrgId(), contributionIds, true)
		endpoints := endPointsResponse.GetEndpoints()
		endpointIds := make([]string, len(endpoints))
		for _, endpoint := range endpoints {
			endpointIds = append(endpointIds, endpoint.GetId())
		}
		if !slices.Contains(endpointIds, req.CiToolId) {
			return fmt.Errorf(errCiToolNotInOrg, req.CiToolId)
		}
	}
	_, ok := byPassDurationValidation[req.WidgetId]
	if req.DurationType.Number() == 0 && !ok {

		st, err := time.Parse(timeLayout, req.StartDate)

		if log.CheckErrorf(err, exceptions.ErrInvalidFormatStartDate) {
			return fmt.Errorf(errInvalidDateFormatRequest, req.StartDate, req.DurationType.String())
		}

		en, e := time.Parse(timeLayout, req.EndDate)
		if log.CheckErrorf(e, exceptions.ErrInvalidFormatEndDate) {
			return fmt.Errorf(errInvalidDateFormatRequest, req.EndDate, req.DurationType.String())
		}

		if st.After(en) {
			return errors.New(errInvalidRequest)
		}
	}

	return nil
}

func ValidateDrilldownDataRequest(req *pb.DrilldownRequest, ctx context.Context, epClt endpoint.EndpointServiceClient) error {
	log.Debugf(exceptions.DebugRequestList, req)
	// TBD Check if Widget Id is valid
	if req.ReportId == "" {
		return fmt.Errorf(errMissingRequiredField, widgetIdField)
	}

	if req.OrgId == "" {
		return fmt.Errorf(errMissingRequiredField, tenantIdField)
	}

	if len(req.Component) > 0 && (len(req.Component) == 1 && req.Component[0] != constants.ALL) {
		for _, resource := range req.Component {
			if !helper.IsResourceInOrg(resource, req.OrgId) {
				return fmt.Errorf(errResourceNotInOrg, resource)
			}
		}
	}

	if req.OrgId != req.SubOrgId && req.SubOrgId != "" {
		if !helper.IsResourceInOrg(req.SubOrgId, req.OrgId) {
			return fmt.Errorf(errResourceNotInOrg, req.SubOrgId)
		}
	}

	if req.CiToolId != "" {
		contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
		endPointsResponse, _ := getAllEndpoints(ctx, epClt, req.GetSubOrgId(), contributionIds, true)
		endpoints := endPointsResponse.GetEndpoints()
		endpointIds := make([]string, len(endpoints))
		for _, endpoint := range endpoints {
			endpointIds = append(endpointIds, endpoint.GetId())
		}
		if !slices.Contains(endpointIds, req.CiToolId) {
			return fmt.Errorf(errCiToolNotInOrg, req.CiToolId)
		}
	}

	if _, ok := byPassDurationValidation[req.ReportId]; !ok {
		if req.DurationType.Number() == 0 {
			st, err := time.Parse(timeLayout, req.StartDate)
			if log.CheckErrorf(err, exceptions.ErrInvalidFormatStartDate) {
				return fmt.Errorf(errInvalidDateFormatRequest, req.StartDate, req.DurationType.String())
			}
			en, e := time.Parse(timeLayout, req.EndDate)
			if log.CheckErrorf(e, exceptions.ErrInvalidFormatEndDate) {
				return fmt.Errorf(errInvalidDateFormatRequest, req.EndDate, req.DurationType.String())
			}

			if st.After(en) {
				return errors.New(errInvalidRequest)
			}
		}
	} else {
		if req.StartDate != "" {
			return fmt.Errorf(errInvalidDateFormatRequest, req.StartDate, req.DurationType.String())
		} else if req.EndDate != "" {
			return fmt.Errorf(errInvalidDateFormatRequest, req.EndDate, req.DurationType.String())
		}
	}

	return nil
}

func ValidateComponentComparisonDataRequest(req *pb.ComponentComparisonRequest) error {
	log.Debugf(exceptions.DebugRequestList, req)
	// TBD Check if Widget Id is valid
	if req.WidgetId == "" {
		return fmt.Errorf(errMissingRequiredField, widgetIdField)
	}

	if req.OrgId == "" {
		return fmt.Errorf(errMissingRequiredField, tenantIdField)
	}

	if req.UserId == "" {
		return fmt.Errorf(errMissingRequiredField, userIdField)
	}

	if req.OrgId != req.SubOrgId && req.SubOrgId != "" {
		if !helper.IsResourceInOrg(req.SubOrgId, req.OrgId) {
			return fmt.Errorf(errResourceNotInOrg, req.SubOrgId)
		}
	}

	if req.DurationType.Number() == 0 {
		st, err := time.Parse(timeLayout, req.StartDate)
		if log.CheckErrorf(err, exceptions.ErrInvalidFormatStartDate) {
			return fmt.Errorf(errInvalidDateFormatRequest, req.StartDate, req.DurationType.String())
		}
		en, e := time.Parse(timeLayout, req.EndDate)
		if log.CheckErrorf(e, exceptions.ErrInvalidFormatEndDate) {
			return fmt.Errorf(errInvalidDateFormatRequest, req.EndDate, req.DurationType.String())
		}

		if st.After(en) {
			return errors.New(errInvalidRequest)
		}
	}

	return nil
}

// Drilldown reports for widgets
func (rah *ReportsHandler) BuildDrilldownReport(ctx context.Context, req *pb.DrilldownRequest) (*pb.DrilldownResponse, error) {
	log.Debugf("BuildDrilldownReport started for drilldown - %s", req.ReportId)
	response := pb.DrilldownResponse{
		Status:  pb.Status_success,
		Error:   "",
		Message: "",
	}
	startTime := time.Now()

	err := ValidateDrilldownDataRequest(req, ctx, rah.endpointClient)
	if err != nil {
		// convert error to GRPC status error so that api-gateway can handle it
		return nil, status.Errorf(codes.InvalidArgument, exceptions.ErrReportServiceReqValidationFailure, err.Error())
	}

	components := []string{}
	if req.UserId != "" {

		if len(req.Component) > 1 || (len(req.Component) == 1 && req.Component[0] != "All") {
			log.Debugf(exceptions.DebugApplyingComponentsFilter, req.Component)
		} else {
			if req.OrgId == req.SubOrgId {
				getOrganizationByIdResponse, err := helper.GetOrganisationsById(ctx, rah.orgServiceClient, req.OrgId, req.UserId)
				if err != nil {
					// convert error to GRPC status error so that api-gateway can handle it
					return nil, status.Errorf(codes.FailedPrecondition, exceptions.ErrFetchingSubOrgForOrg, req.OrgId, err.Error())
				}
				serviceResponse, err := helper.GetOrganisationServices(ctx, rah.client, req.OrgId)
				if err != nil {
					return nil, status.Errorf(codes.FailedPrecondition, exceptions.ErrFetchingServicesForOrg, req.OrgId, err.Error())
				}
				for _, service := range serviceResponse.GetService() {
					if req.OrgId == service.OrganizationId {
						components = append(components, service.Id)
					}
				}
				helper.GetComponentsRecursively(getOrganizationByIdResponse.GetOrganization(), serviceResponse, &components)
				req.Component = components
			} else if req.OrgId != req.SubOrgId && req.SubOrgId != "" {
				getOrganizationByIdResponse, err := helper.GetOrganisationsById(ctx, rah.orgServiceClient, req.SubOrgId, req.UserId)
				if err != nil {
					// convert error to GRPC status error so that api-gateway can handle it
					return nil, status.Errorf(codes.FailedPrecondition, "Exception while fetching sub orgs for sub org : %s : %s", req.SubOrgId, err.Error())
				}
				serviceResponse, err := helper.GetOrganisationServices(ctx, rah.client, req.SubOrgId)
				if err != nil {
					return nil, status.Errorf(codes.FailedPrecondition, exceptions.ErrFetchingServicesForSubOrg, req.SubOrgId, err.Error())
				}
				for _, service := range serviceResponse.GetService() {
					if req.SubOrgId == service.OrganizationId {
						components = append(components, service.Id)
					}
				}
				helper.GetComponentsRecursively(getOrganizationByIdResponse.GetOrganization(), serviceResponse, &components)
				req.Component = components
			}
		}
	} else {
		if req.OrgId != req.SubOrgId && req.SubOrgId != "" {
			serviceResponse, err := getOrganisationServices(ctx, rah.client, req.SubOrgId)
			if err != nil {
				return nil, status.Errorf(codes.FailedPrecondition, exceptions.ErrFetchingServicesForSubOrg, req.SubOrgId, err.Error())
			}
			if len(serviceResponse.GetService()) > 0 {
				for i := 0; i < len(serviceResponse.GetService()); i++ {
					service := serviceResponse.GetService()[i]
					if len(req.Component) > 1 || (len(req.Component) == 1 && req.Component[0] != "All") {
						for _, comp := range req.Component {
							if comp == service.Id {
								components = append(components, service.Id)
							}
						}
					} else {
						components = append(components, service.Id)
					}
				}
			}
			req.Component = components
		}
	}

	log.Debugf("Components in request %v for widget %s took %v", req.Component, req.ReportId, time.Since(startTime).Milliseconds())

	response, err = getDrilldownReport(req, response, ctx, rah)
	log.Debugf("BuildDrilldownReport completed for drilldown - %s", req.ReportId)
	return &response, err

}

func getDrilldownReport(req *pb.DrilldownRequest, response pb.DrilldownResponse, ctx context.Context, rah *ReportsHandler) (pb.DrilldownResponse, error) {
	if len(req.Component) == 0 && req.CiToolId == "" {
		response.Reports = &structpb.ListValue{}
		return response, nil
	}

	var aggrBy string // aggregation interval for date histogram - week when duration is a month; day when duration is week

	if req.DurationType == pb.DurationType_CURRENT_WEEK ||
		req.DurationType == pb.DurationType_PREVIOUS_WEEK ||
		req.DurationType == pb.DurationType_LAST_7_DAYS ||
		req.DurationType == pb.DurationType_TWO_WEEKS_BACK {
		aggrBy = "day"
	} else if req.DurationType == pb.DurationType_CURRENT_MONTH ||
		req.DurationType == pb.DurationType_PREVIOUS_MONTH ||
		req.DurationType == pb.DurationType_LAST_30_DAYS ||
		req.DurationType == pb.DurationType_TWO_MONTHS_BACK {
		aggrBy = "week"
	} else if req.DurationType == pb.DurationType_LAST_90_DAYS {
		aggrBy = "month"
	} else if req.DurationType == pb.DurationType_CUSTOM_RANGE {
		d, _ := getDateDiffInDays(req.StartDate, req.EndDate)

		if d > maxDaysInMonth {
			aggrBy = "month"
		} else if d > maxDaysInWeek {
			aggrBy = "week"
		} else {
			aggrBy = "day"
		}
	}

	replacements := map[string]any{
		"startDate":  req.StartDate,
		"endDate":    req.EndDate,
		"orgId":      req.OrgId,
		"subOrgId":   req.SubOrgId,
		"component":  req.Component,
		"aggrBy":     aggrBy,
		"timeFormat": req.TimeFormat,
		"timeZone":   req.TimeZone,
		"viewOption": req.ViewOption.String(),
	}

	if req.CiToolId != "" {
		replacements["ciToolId"] = req.CiToolId
	}

	if req.GetTimeFormat() == "" {
		replacements["timeFormat"] = constants.DURATION_TWELVE_HOURS
	} else {
		replacements["timeFormat"] = req.GetTimeFormat()
	}

	if req.TimeZone == "" {
		replacements["timeZone"] = constants.LOCATION_EUROPE_OR_LONDON
	}

	if req.ReportInfo != nil {
		if req.ReportInfo.DeploymentEnv != "" {
			replacements["targetEnv"] = req.ReportInfo.DeploymentEnv
		}
		if req.ReportInfo.ComponentId != "" {
			replacements["componentIdForNestedDrillDown"] = req.ReportInfo.ComponentId
		}
		if req.ReportInfo.Branch != "" {
			replacements[constants.REQUEST_BRANCH] = req.ReportInfo.Branch
		}
		if req.ReportInfo.Code != "" {
			replacements["vulCode"] = req.ReportInfo.Code
		}
		if req.ReportInfo.ScannerName != "" {
			if key, exists := scannerNameReverseMap[req.ReportInfo.ScannerName]; exists {
				replacements["scannerName"] = key
			} else {
				replacements["scannerName"] = req.ReportInfo.ScannerName
			}
		}
		if req.ReportInfo.RunId != "" {
			replacements["runId"] = req.ReportInfo.RunId
		}
		if req.ReportInfo.Author != "" {
			replacements["author"] = req.ReportInfo.Author
		}
		if req.ReportInfo.JobId != "" {
			replacements["jobId"] = req.ReportInfo.JobId
		}
		if req.ReportInfo.LicenseType != "" {
			replacements["licenseType"] = req.ReportInfo.LicenseType
		}
		if req.ReportInfo.ScannerNameList != nil {
			var mappedScannerNames []string
			for _, name := range req.ReportInfo.ScannerNameList {
				if key, exists := scannerNameReverseMap[name]; exists {
					mappedScannerNames = append(mappedScannerNames, key)
				} else {
					mappedScannerNames = append(mappedScannerNames, name)
				}
			}
			replacements["scannerNameList"] = mappedScannerNames
		}
		if req.ReportInfo.RunIdList != nil {
			replacements["runIDList"] = req.ReportInfo.RunIdList
		}

		if req.ReportInfo.TestSuiteName != "" {
			replacements["testSuiteName"] = req.ReportInfo.TestSuiteName
		}

		if req.ReportInfo.TestCaseName != "" {
			replacements["testCaseName"] = req.ReportInfo.TestCaseName
		}

		if req.ReportInfo.AutomationId != "" {
			replacements["automationId"] = req.ReportInfo.AutomationId
		}

		if req.ReportInfo.Branch != "" {
			replacements[constants.REQUEST_BRANCH] = req.ReportInfo.Branch
		}

		if req.ReportInfo.RunNumber != "" {
			replacements[constants.REQUEST_RUN_NUMBER] = req.ReportInfo.RunNumber
		}

	}
	//replacements for placeholders inside JOLT specs
	replacementsSpec := map[string]any{
		"normalizeMonthInSpec": "@x",
	}
	// This branch field outside report_info is used in Component Summary where a branch filter is mandatory.
	_, ok := replacements[constants.REQUEST_BRANCH]
	if len(req.Branch) > 0 && !ok {
		replacements[constants.REQUEST_BRANCH] = req.Branch
	}

	dateBydurationType := models.CalculateDateBydurationType{
		StartDateStr:       req.StartDate,
		EndDateStr:         req.EndDate,
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   replacementsSpec,
		NormalizeMonthFlag: false,
		IsComputFlag:       false,
		CurrentTime:        time.Now().UTC(),
	}
	calculateDateBydurationType(dateBydurationType)

	if query, ok := db.DrillDownQueryDefinitionMap[req.ReportId]; ok {
		result, err := processDrilldownQueryAndSpec(replacements, query, req.ReportId, &response)
		if err == nil {
			response.Reports = result
			response.ColumnAttributes = map[string]string{"column": db.DrillDownFilterDefinitionMap[req.ReportId]}
		} else {
			response.Status = pb.Status_error
			response.Error = err.Error()
		}
	} else {
		switch req.ReportId {
		case "component":
			result, err := ComponentReport(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "componentName"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "workflows":
			result, err := AutomationReport(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "workflowRuns", "component-summary-workflowRuns":
			result, err := AutomationRunReport(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "component-summary-workflows":
			result, err := AutomationReportForBranch(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "security-components":
			result, err := SecurityComponentDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "componentName"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "security-workflows":
			result, err := SecurityAutomationDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "security-workflowRuns":
			result, err := SecurityAutomationRunDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "componentName"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "security-scan-type-workflows":
			result, err := SecurityScanTypeWorkflowsDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "componentName"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "commits", "component-summary-commits":
			result, err := CommitsReport(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "pullrequests":
			result, err := PullRequestsReport(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "runInitiatingCommits":
			result, err := CPSRunInitiatingCommits(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "builds", "component-summary-builds":
			result, err := CodeProgressionSnapshotBuilds(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "deployments":
			result, err := CodeProgressionSnapshotDeployments(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "successfulBuildsDuration":
			result, err := SuccessfulBuildDuration(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "deploymentOverview", "component-summary-deploymentOverview":
			result, err := DeploymentOverviewDrilldown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "doraMetrics-deploymentFrequency", "doraMetrics-deploymentLeadTime":
			result, err := DeploymentFrequencyAndLeadTime(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "doraMetrics-failureRate":
			result, err := FailureRate(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "doraMetrics-mttr":
			result, err := DoraMetricsMttr(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "pluginsInfo":
			result, err := CiInsightsPluginsInfo(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "runInformation":
			result, err := CiInsightsCompletedRunsAndTime(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "test-insights-components":
			result, err := TestComponentDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "componentName"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "test-insights-workflowRuns":
			result, err := TestAutomationRunDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "componentName"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "test-insights-workflows":
			result, err := TestAutomationDrilldown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "componentName"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case constants.TEST_OVERVIEW_VIEW_RUN_ACTIVITY:
			result, err := TestInsightsViewRunActivityDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case "test-overview-view-run-activity-logs":
			result, err := TestInsightsViewRunActivityLogsDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "component"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case constants.TEST_OVERVIEW_TOTAL_TEST_CASES:
			result, err := TotalTestCasesDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "componentName"}
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case constants.TEST_OVERVIEW_TOTAL_RUNS:
			result, err := TestInsightsTotalRunsDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": "runStatus"}
			} else if err == db.ErrNoDataFound {
				response.Message = db.ErrNoDataFound.Error()
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case constants.RUN_DETAILS_TEST_RESULTS:
			result, err := RunDetailsTestResults(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": ""}
			} else if err == db.ErrNoDataFound {
				response.Message = db.ErrNoDataFound.Error()
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case constants.RUN_DETAILS_TOTAL_TEST_CASES:
			result, err := RunDetailsTotalTestCasesDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": ""}
			} else if err == db.ErrNoDataFound {
				response.Message = db.ErrNoDataFound.Error()
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case constants.RUN_DETAILS_TEST_CASE_LOG:
			result, err := RunDetailsTestCaseLogDrillDown(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": ""}
			} else if err == db.ErrNoDataFound {
				response.Message = db.ErrNoDataFound.Error()
			} else {
				response.Status = pb.Status_error
				response.Error = err.Error()
			}
		case constants.RUN_DETAILS_TEST_RESULTS_INDICATORS:
			result, err := RunDetailsTestResultsIndicators(replacements, ctx, rah.client)
			if err == nil {
				response.Reports = result
				response.ColumnAttributes = map[string]string{"column": ""}
			} else {
				response.Status = pb.Status_error
				response.Error = db.ErrInternalServer.Error()
			}
		default:
			log.Warn("Drilldown Report Id : " + req.ReportId + " is not applicable")
		}
	}

	log.Debugf(exceptions.DebugComponentInReqForWidget, req.Component, req.ReportId)

	return response, nil
}

// Computed Drilldown reports for widgets
func (rah *ReportsHandler) BuildComputedDrilldownReport(ctx context.Context, req *pb.DrilldownRequest) (*pb.DrilldownResponse, error) {
	log.Debugf("BuildComputedDrilldownReport started for drilldown - %s", req.ReportId)
	response := pb.DrilldownResponse{
		Status:  pb.Status_success,
		Error:   "",
		Message: "",
	}

	err := ValidateDrilldownDataRequest(req, ctx, rah.endpointClient)
	if err != nil {
		// convert error to GRPC status error so that api-gateway can handle it
		return nil, status.Errorf(codes.InvalidArgument, exceptions.ErrReportServiceReqValidationFailure, err.Error())
	}

	// validate access if the request is intended for CI Insights drilldown report
	if err := helper.ValidateCIInsightsDrilldownReportAccess(ctx, rah.rbacClt, req); err != nil {
		return nil, err
	}

	components := []string{}

	if len(req.Component) > 1 || (len(req.Component) == 1 && req.Component[0] != "All") {
		log.Debugf(exceptions.DebugApplyingComponentsFilter, req.Component)
	} else {
		if req.OrgId == req.SubOrgId {
			getOrganizationByIdResponse, err := helper.GetOrganisationsById(ctx, rah.orgServiceClient, req.OrgId, req.UserId)
			if err != nil {
				// convert error to GRPC status error so that api-gateway can handle it
				return nil, status.Errorf(codes.FailedPrecondition, exceptions.ErrFetchingSubOrgForOrg, req.OrgId, err.Error())
			}
			serviceResponse, err := getOrganisationServices(ctx, rah.client, req.OrgId)
			if err != nil {
				return nil, status.Errorf(codes.FailedPrecondition, exceptions.ErrFetchingServicesForOrg, req.OrgId, err.Error())
			}
			for _, service := range serviceResponse.GetService() {
				if req.OrgId == service.OrganizationId {
					components = append(components, service.Id)
				}
			}
			helper.GetComponentsRecursively(getOrganizationByIdResponse.GetOrganization(), serviceResponse, &components)
			req.Component = components
		} else if req.OrgId != req.SubOrgId && req.SubOrgId != "" {
			serviceResponse, err := getOrganisationServices(ctx, rah.client, req.SubOrgId)
			if err != nil {
				log.Errorf(err, exceptions.ErrFetchingServicesForSubOrgWithoutFormatting)
			} else {
				if len(serviceResponse.GetService()) > 0 {
					for i := 0; i < len(serviceResponse.GetService()); i++ {
						service := serviceResponse.GetService()[i]
						if len(req.Component) > 1 || (len(req.Component) == 1 && req.Component[0] != "All") {
							for _, comp := range req.Component {
								if comp == service.Id {
									components = append(components, service.Id)
								}
							}
						} else {
							components = append(components, service.Id)
						}
					}
				}
			}
			req.Component = components
		}
	}
	log.Debugf(exceptions.DebugComponentInReqForWidget, req.Component, req.ReportId)

	found := false
	// for _, reportsId := range constants.ComputedDrilldown {
	// 	if req.ReportId == reportsId {
	// 		found = true
	// 	}
	// }

	if !found {
		response, err := getDrilldownReport(req, response, ctx, rah)
		return &response, err
	}

	if len(req.Component) > 1 {
		res, err := rah.BuildDrilldownReport(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else if len(req.Component) == 0 {
		response.Reports = &structpb.ListValue{}
		return &response, nil
	}
	//Get the filters from request and create map of replacement placeholders
	replacements := map[string]any{
		"metricsKey": req.ReportId,
		"startDate":  req.StartDate,
		"endDate":    req.EndDate,
		"orgId":      req.OrgId,
		"subOrgId":   req.SubOrgId,
		"component":  req.Component,
		"ciToolId":   req.CiToolId,
		"timeFormat": req.TimeFormat,
		"timeZone":   req.TimeZone,
	}

	dateBydurationType := models.CalculateDateBydurationType{
		StartDateStr:       req.StartDate,
		EndDateStr:         req.EndDate,
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   nil,
		NormalizeMonthFlag: false,
		IsComputFlag:       true,
		CurrentTime:        time.Now().UTC(),
	}
	calculateDateBydurationType(dateBydurationType)

	res, err := internal.GetComputedData(req.ReportId, replacements)
	if err != nil {
		response.Status = pb.Status_error
		response.Error = fmt.Sprintf("failed to create widget in BuildComputedDrilldownReport : %s", err.Error())
		return &response, nil
	}
	var mapResp map[string]interface{}
	json.Unmarshal([]byte(res), &mapResp)
	mv := mapResp["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(mv) > 0 {
		//TBD.  Handle multiple components here.
		jmv, err := json.Marshal(mv[0].(map[string]interface{})["_source"].(map[string]interface{})["metric_value"].(map[string]interface{})["reports"])
		if err == nil {
			reports := structpb.ListValue{}
			err = protojson.Unmarshal([]byte(jmv), &reports)
			response.Reports = &reports
		}
		if err != nil {
			response.Status = pb.Status_error
			response.Error = fmt.Sprintf("failed to create widget in BuildComputedDrilldownReport : %s", err.Error())
			return &response, nil
		}
	} else {
		response.Message = constants.NO_DATA_FOUND
	}
	return &response, nil
}

func (rah *ReportsHandler) BuildComputedReport(ctx context.Context, req *pb.ReportServiceRequest) (*pb.ReportServiceResponse, error) {
	response := pb.ReportServiceResponse{
		Status:  pb.Status_success,
		Error:   "",
		Widget:  nil,
		Message: "",
	}
	err := ValidateDataRequest(req, ctx, rah.endpointClient)
	if err != nil {
		// convert error to GRPC status error so that api-gateway can handle it
		return nil, status.Errorf(codes.InvalidArgument, exceptions.ErrReportServiceReqValidationFailure, err.Error())
	}
	components := []string{}
	if req.OrgId != req.SubOrgId && req.SubOrgId != "" {
		serviceResponse, err := getOrganisationServices(ctx, rah.client, req.OrgId)
		if err != nil {
			log.Errorf(err, exceptions.ErrFetchingServicesForSubOrgWithoutFormatting)
		} else {
			if len(serviceResponse.GetService()) > 0 {
				for i := 0; i < len(serviceResponse.GetService()); i++ {
					service := serviceResponse.GetService()[i]
					if len(req.Component) > 1 || (len(req.Component) == 1 && req.Component[0] != "All") {
						for _, comp := range req.Component {
							if comp == service.Id {
								components = append(components, service.Id)
							}
						}
					} else {
						components = append(components, service.Id)
					}
				}
			}
		}
		req.Component = components
	}
	//Check for Components. If Component is multiple(Apart from "All"), then redirect to call
	if len(req.Component) > 1 {
		res, err := rah.BuildReport(ctx, req)
		if err != nil {
			return nil, err
		}

		return res, nil
	} else if len(req.Component) == 0 {
		response.Message = constants.NO_DATA_FOUND
		return &response, nil
	}

	replacements := map[string]any{
		"metricsKey": req.WidgetId,
		"orgId":      req.OrgId,
		"subOrgId":   req.SubOrgId,
		"component":  req.Component,
		"startDate":  req.StartDate,
		"endDate":    req.EndDate,
		"ciToolId":   req.CiToolId,
		"ciToolType": req.CiToolType,
		"sortBy":     req.SortBy,
	}

	dateBydurationType := models.CalculateDateBydurationType{
		StartDateStr:       req.StartDate,
		EndDateStr:         req.EndDate,
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   nil,
		NormalizeMonthFlag: false,
		IsComputFlag:       true,
		CurrentTime:        time.Now().UTC(),
	}
	calculateDateBydurationType(dateBydurationType)

	log.Debugf("REPLACEMNETS", replacements)

	res, err := internal.GetComputedData(req.WidgetId, replacements)
	if err != nil {
		response.Status = pb.Status_error
		response.Error = fmt.Sprintf("failed to create widget in BuildComputedReport : %s", err.Error())
		return &response, nil
	}

	var mapResp map[string]interface{}
	json.Unmarshal([]byte(res), &mapResp)

	mv := mapResp["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(mv) > 0 {
		//TBD.  Handle multiple components here.
		jmv, err := json.Marshal(mv[0].(map[string]interface{})["_source"].(map[string]interface{})["metric_value"])
		if err == nil {
			w := pb.Widget{}
			err = protojson.Unmarshal([]byte(jmv), &w)
			response.Widget = &w
		}

		if err != nil {
			response.Status = pb.Status_error
			response.Error = fmt.Sprintf("failed to create widget in BuildComputedReport : %s", err.Error())
			return &response, nil
		}
	} else {
		response.Message = constants.NO_DATA_FOUND
	}

	return &response, nil
}

func ValidateWidgetDataRequest(req *pb.ManageWidgetRequest) error {
	log.Debugf(exceptions.DebugRequestList, req)
	if req.OrgId == "" {
		return fmt.Errorf(errMissingRequiredField, req.OrgId)
	}
	if req.DashboardName == "" {
		return fmt.Errorf(errMissingRequiredField, req.DashboardName)
	}
	return nil
}

func (rah *ReportsHandler) GetWidgets(ctx context.Context, req *pb.ManageWidgetRequest) (*pb.ManageWidgetResponse, error) {
	log.Debugf("Manage Widgets api started - dashboard in request %v for widget", req.DashboardName)
	err := ValidateWidgetDataRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, exceptions.ErrReportServiceReqValidationFailure, err.Error())
	}
	response := pb.ManageWidgetResponse{
		Widgets: make([]*pb.ReportWidget, 0),
	}
	dn := req.DashboardName
	switch dn {
	case constants.SDA_DASHBOARD:
		response.Widgets, err = db.GetWidgetDefinitionList(constants.SDA_DASHBOARD)
	case constants.SECURITY_INSIGHTS_DASHBOARD:
		response.Widgets, err = db.GetWidgetDefinitionList(constants.SECURITY_INSIGHTS_DASHBOARD)
	case constants.FLOW_METRICS_DASHBOARD:
		response.Widgets, err = db.GetWidgetDefinitionList(constants.FLOW_METRICS_DASHBOARD)
	case constants.DORA_METRICS_DASHBOARD:
		response.Widgets, err = db.GetWidgetDefinitionList(constants.DORA_METRICS_DASHBOARD)
	case constants.COMPONENT_SECURITY_DASHBOARD:
		response.Widgets, err = db.GetWidgetDefinitionList(constants.COMPONENT_SECURITY_DASHBOARD)
	case constants.APPLICATION_SECURITY_DASHBOARD:
		response.Widgets, err = db.GetWidgetDefinitionList(constants.APPLICATION_SECURITY_DASHBOARD)
	case constants.COMPONENT_SUMMARY_DASHBOARD:
		response.Widgets = nil // needs to be configured
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Invalid dashboard name: %s", req.DashboardName)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Reading widget definition failed: %s", err.Error())
	}
	return &response, nil
}

func (rah *ReportsHandler) BuildReportLayout(ctx context.Context, req *pb.ReportLayoutRequest) (*pb.ReportLayoutResponse, error) {

	replacements := map[string]any{
		"orgId":     req.OrgId,
		"component": req.Component,
	}
	var reportLayoutArray []*pb.ReportLayout
	response := pb.ReportLayoutResponse{Status: pb.Status_success, ComponentId: req.Component, Message: "", Widgets: nil, Error: ""}

	isSonarData := internal.IsSonarWidgetsApplicable(replacements, ctx, constants.SECURITY_INDEX, constants.RAW_SCAN_RESULTS_INDEX)
	isSecurity := internal.IsSecurityWidgetsApplicable(replacements, ctx, constants.SECURITY_INDEX, constants.RAW_SCAN_RESULTS_INDEX)
	if isSonarData {
		widgets := constants.ComponentSummarySonarWidgets
		reportLayoutArray = getReportLayouts(widgets, constants.ComponentSummaryWidgetsLayout)
	} else if isSecurity {
		widgets := constants.ComponentSummaryNoSonarWidgets

		reportLayoutArray = getReportLayouts(widgets, constants.ComponentSummaryWidgetsLayout)
	} else {
		widgets := constants.ComponentSummaryNoScannerWidgets

		reportLayoutArray = getReportLayouts(widgets, constants.ComponentSummaryWidgetsLayout)
	}
	response.Widgets = reportLayoutArray

	return &response, nil

}

func getReportLayouts(widgets []string, widgetLayout map[string]string) []*pb.ReportLayout {
	reportLayoutArray := []*pb.ReportLayout{}

	for _, widgetId := range widgets {
		widgetLayoutStr := widgetLayout[widgetId]
		reportLayout := pb.ReportLayout{}
		err := json.Unmarshal([]byte(string(widgetLayoutStr)), &reportLayout)
		log.CheckErrorf(err, "Error unmarshalling report layout from : %s", widgetLayoutStr)

		reportLayoutArray = append(reportLayoutArray, &reportLayout)

	}
	return reportLayoutArray
}

func (rah *ReportsHandler) GetEnvironments(ctx context.Context, req *pb.EnvironmentRequest) (*pb.EnvironmentResponse, error) {
	response := pb.EnvironmentResponse{
		Environments: []*pb.Environment{},
	}
	components := []string{}
	if req.OrgId != req.SubOrgId && req.SubOrgId != "" {
		serviceResponse, err := helper.GetOrganisationServices(ctx, rah.client, req.SubOrgId)
		if err != nil {
			log.Errorf(err, exceptions.ErrFetchingServicesForSubOrgWithoutFormatting)
		} else {
			if len(serviceResponse.GetService()) > 0 {
				for i := 0; i < len(serviceResponse.GetService()); i++ {
					service := serviceResponse.GetService()[i]
					components = append(components, service.Id)
				}
			}
		}
	}

	if len(components) == 0 {
		components = append(components, "All")
	}
	replacements := map[string]any{
		"orgId":     req.OrgId,
		"component": components,
		"startDate": getStartOfTheDay(time.Now().AddDate(0, -5, 0)).Format(timeLayout),
		"endDate":   getEndOfTheDay(time.Now()).Format(timeLayout),
		"timeZone":  req.TimeZone,
	}
	environments, _, err := internal.FetchAllDeployedEnvironments("", replacements, ctx)
	log.CheckErrorf(err, "Enviroment endpoint api failing")

	if req.SubOrgId != "" {
		contributionIds := []string{constants.ENVIRONMENT_ENDPOINT}
		endPointsResponse, err := helper.GetAllEndpoints(ctx, rah.endpointClient, req.SubOrgId, contributionIds, true)
		if err != nil {
			log.Errorf(err, exceptions.ErrFetchingEndpoints)
			return nil, exceptions.GetExceptionByCode(exceptions.ErrEndpointAPIFailure)
		} else {
			endpoints := endPointsResponse.Endpoints
			log.Debugf(exceptions.DebugEndpointList, endpoints)
			if len(endpoints) > 0 {
				for _, environment := range environments {
					for _, endpoint := range endpoints {
						if endpoint.ContributionType == "cb.platform.environment" && endpoint.Name == environment {
							env := pb.Environment{
								Id:         endpoint.Id,
								Name:       endpoint.Name,
								IsDisabled: endpoint.IsDisabled,
								ResourceId: endpoint.ResourceId,
							}
							response.Environments = append(response.Environments, &env)
							break
						}
					}
				}
			}
		}
	}
	return &response, nil
}

func (rah *ReportsHandler) GetEnvironmentsv2(ctx context.Context, req *pb.EnvironmentRequest) (*pb.EnvironmentResponse, error) {
	response := pb.EnvironmentResponse{
		Environments: []*pb.Environment{},
	}
	contributionIds := []string{constants.ENVIRONMENT_ENDPOINT}
	endPointsResponse, err := helper.GetAllEndpoints(ctx, rah.endpointClient, req.OrgId, contributionIds, true)
	if err != nil {
		log.Errorf(err, exceptions.ErrFetchingEndpoints)
		return nil, exceptions.GetExceptionByCode(exceptions.ErrEndpointAPIFailure)
	}
	endpoints := endPointsResponse.Endpoints
	log.Debugf(exceptions.DebugEndpointList, endpoints)
	if req.Name != "" {
		for _, endpoint := range endpoints {
			if endpoint.ContributionType == "cb.platform.environment" && endpoint.Name == req.Name {
				env := pb.Environment{Id: endpoint.Id, Name: endpoint.Name, IsDisabled: endpoint.IsDisabled, ResourceId: endpoint.ResourceId}
				response.Environments = append(response.Environments, &env)
				break
			}
		}
	} else {
		for _, endpoint := range endpoints {
			if endpoint.ContributionType == "cb.platform.environment" {
				env := pb.Environment{Id: endpoint.Id, Name: endpoint.Name, IsDisabled: endpoint.IsDisabled, ResourceId: endpoint.ResourceId}
				response.Environments = append(response.Environments, &env)
			}
		}
	}
	return &response, nil
}

func getReportRequest(req *pb.ComputeUpdateRequest, duration string) *pb.ReportServiceRequest {
	request := pb.ReportServiceRequest{}
	request.OrgId = req.OrgId
	request.Component = strings.Fields(req.ComponentId)
	if duration == constants.CURRENT_MONTH {
		request.DurationType = pb.DurationType_CURRENT_MONTH
	} else {
		request.DurationType = pb.DurationType_CURRENT_WEEK
	}
	request.WidgetId = req.MetricKey
	return &request

}

func checkNullResponse(response []*pb.MetricInfo) error {

	if len(response) > 0 {
		for _, widgetElement := range response {
			if widgetElement != nil && widgetElement.Data == nil {
				return db.ErrInternalServer
			}
		}
	}
	return nil
}

func checkChartNullResponse(response []*pb.ChartInfo) error {

	if len(response) > 0 {
		for _, widgetElement := range response {
			if widgetElement != nil && widgetElement.Data == nil {
				return db.ErrInternalServer
			}
		}
	}
	return nil
}

func constructComputeResponseQuery(dev ComputeInfo) string {
	var bulkInsertQuery string

	if dev.MetricValue != nil {
		data, error := json.Marshal(dev)
		if log.CheckErrorf(error, "could not unmarshal potato proto") {
			// The data can't be marshalled, return nil since this is not a recoverable error
			return ""
		} else {

			bulkInsertQuery += fmt.Sprintf("%s%s%s%s%s%s%v%s%v%s%s%s\n", "{\"index\":{\"_index\":\""+constants.COMPUTE_INDEX+"\",\"_id\": \"", dev.OrgId, "_",
				dev.ComponentId, "_", dev.StartDate, "_", dev.EndDate, "_", dev.MetricKey, "\"}}\n", string(data))
		}
	} else {
		log.Debugf("No data found for widget %s and component %s", dev.MetricKey, dev.ComponentId)
	}

	return bulkInsertQuery
}

func constructComputeScheduleQuery(dev ComputeSchedule) string {
	var bulkInsertQuery string

	if len(dev.JobName) > 0 {
		data, error := json.Marshal(dev)
		if log.CheckErrorf(error, "could not unmarshal potato proto") {
			// The data can't be marshalled, return nil since this is not a recoverable error
			return ""
		} else {

			bulkInsertQuery += fmt.Sprintf("%s%s%s%v\n", "{\"index\":{\"_index\":\""+constants.COMPUTE_INDEX+"\",\"_id\": \"", dev.JobName, "\"}}\n", string(data))
		}
	} else {
		log.Debugf("No job name found for compute schedule insertion")
	}

	return bulkInsertQuery
}
func weekStartDate(date time.Time) time.Time {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}
func CurrentWeekStartAndEndDate() (string, string) {
	now := time.Now()
	DDMMYYYYhhmmss := "2006-01-02 15:04:05"

	currentYear, currentMonth, currentDay := now.Date()
	currentLocation := now.Location()

	currentDayTime := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation)

	startWeek := weekStartDate(currentDayTime)
	endWeek := startWeek.AddDate(0, 0, 6)
	endWeekTime := endWeek.Add(time.Hour*time.Duration(23) +
		time.Minute*time.Duration(59) +
		time.Second*time.Duration(59))

	weekStart := startWeek.Format(DDMMYYYYhhmmss)
	weekEnd := endWeekTime.Format(DDMMYYYYhhmmss)
	log.Debugf(" Current Week StartDate:%s, EndDate:%s", weekStart, weekEnd)
	return weekStart, weekEnd

}
func CurrentMonthStartAndEndDate() (string, string) {
	now := time.Now()
	DDMMYYYYhhmmss := "2006-01-02 15:04:05"

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	monthStart := firstOfMonth.Format(DDMMYYYYhhmmss)
	endMonthTime := lastOfMonth.Add(time.Hour*time.Duration(23) +
		time.Minute*time.Duration(59) +
		time.Second*time.Duration(59))
	monthEnd := endMonthTime.Format(DDMMYYYYhhmmss)
	log.Debugf(" Current Month StartDate:%s, EndDate:%s", monthStart, monthEnd)
	return monthStart, monthEnd

}

func getDateDiffInDays(startDateStr string, endDateStr string) (int64, error) {
	startDate, err := time.Parse(timeLayoutDateHistogram, startDateStr)
	if err != nil {
		return 0, err
	}

	endDate, err := time.Parse(timeLayoutDateHistogram, endDateStr)
	if err != nil {
		return 0, err
	}

	diff := int64(endDate.Sub(startDate).Hours() / 24)

	return diff, nil
}

func performComputeInfo(oclient *opensearch.Client, response *pb.ReportServiceResponse, request *pb.ReportServiceRequest, computeResponse ComputeInfo) (ComputeInfo, error) {
	w := response.Widget
	if len(w.GetContent()) > 0 {
		for _, content := range w.Content {
			if content != nil {
				if content.Header != nil && len(content.Header) > 0 {
					err := checkNullResponse(content.Header)
					if err != nil {
						log.Debugf("Ignoring computation for widget :%s since header is null", w.Id)
						return computeResponse, err
					}
				}
				if content.Section != nil && len(content.Section) > 0 {
					err := checkChartNullResponse(content.Section)
					if err != nil {
						log.Debugf("Ignoring computation for widget :%s since section is null", w.Id)
						return computeResponse, err
					}
				}
				if content.Footer != nil && len(content.Footer) > 0 {
					err := checkNullResponse(content.Footer)
					if err != nil {
						log.Debugf("Ignoring computation for widget :%s since footer is null", w.Id)
						return computeResponse, err
					}
				}
			}
		}

		computeResponse.MetricValue = w

		bulkQuery := constructComputeResponseQuery(computeResponse)
		if bulkQuery != "" && len(bulkQuery) > 0 {
			err1 := db.InsertBulkData(oclient, bulkQuery)
			if log.CheckErrorf(err1, "Could not insert compute data") {
				return computeResponse, err1
			} else {
				log.Debugf("Compute data inserted for %s : %s : %s", computeResponse.MetricKey, request.DurationType, request.Component)
			}
		} else {
			log.Infof("compute data bulk query failed for %s : %s : %s", computeResponse.MetricKey, request.DurationType, request.Component)
		}
	}

	return computeResponse, nil
}

type ComputeInfo struct {
	OrgId         string     `json:"org_id"`
	OrgName       string     `json:"org_name"`
	SubOrgName    string     `json:"sub_org_name"`
	ComponentId   string     `json:"component_id"`
	ComponentName string     `json:"component_name"`
	StartDate     string     `json:"start_date"`
	EndDate       string     `json:"end_date"`
	MetricKey     string     `json:"metric_key"`
	MetricValue   *pb.Widget `json:"metric_value"`
}

func (rah *ReportsHandler) UpdateComputeData(ctx context.Context, req *pb.ComputeUpdateRequest) (*pb.ComputeUpdateResponse, error) {
	computeResponse := pb.ComputeUpdateResponse{}
	if len(req.OrgId) > 0 && len(req.ComponentId) > 0 {
		client, err := opensearchconfig.GetOpensearchConnection()
		if log.CheckErrorf(err, exceptions.ErrOpenSearchConnectionInProcessDrillDownQuery) {
			return nil, err
		}
		for _, duration := range constants.Durations {
			reportRequest := getReportRequest(req, duration)
			response, err := rah.BuildReport(ctx, reportRequest)
			if err == nil {
				computeInfo := ComputeInfo{}
				computeInfo.MetricKey = req.MetricKey
				computeInfo.ComponentId = req.ComponentId
				computeInfo.OrgId = req.OrgId
				if duration == constants.CURRENT_MONTH {
					startDate, endDate := CurrentMonthStartAndEndDate()
					computeInfo.StartDate = startDate
					computeInfo.EndDate = endDate
				} else {
					startDate, endDate := CurrentWeekStartAndEndDate()
					computeInfo.StartDate = startDate
					computeInfo.EndDate = endDate

				}
				performComputeInfo(client, response, reportRequest, computeInfo)
				computeResponse.Response = "Updated compute data successfully"
			} else {
				log.Errorf(err, "Error in updating response : %s", err.Error())
				computeResponse.Response = "Error in updating compute response"
			}
		}

	}

	return &computeResponse, nil
}

func (rah *ReportsHandler) GetRawData(ctx context.Context, req *pb.DataRequest) (*pb.DataResponse, error) {
	dataResponse := pb.DataResponse{}
	if req != nil && len(req.Query) > 0 {
		client, err := opensearchconfig.GetOpensearchConnection()
		if log.CheckErrorf(err, exceptions.ErrOpenSearchConnectionInProcessDrillDownQuery) {
			return nil, err
		} else {
			if req.IsMappingRequest {
				mappingResponse, respErr := db.GetOpensearchMappingData(req.IndexName, client)
				if respErr != nil {
					log.Errorf(err, "Error in getting raw data : %s", respErr.Error())
					return nil, respErr
				} else {
					dataResponse.Response = mappingResponse

				}
			} else {
				response, respErr := db.GetOpensearchData(req.Query, req.IndexName, client)
				if respErr != nil {
					log.Errorf(err, "Error in getting raw data : %s", respErr.Error())
					return nil, respErr
				} else {
					dataResponse.Response = response

				}
			}

		}

	}
	return &dataResponse, nil
}

type ComputeSchedule struct {
	JobName     string  `json:"job_name"`
	LastRunTime *string `json:"last_run_start_time"`
	Timestamp   *string `json:"timestamp"`
}

func (rah *ReportsHandler) UpdateRawData(ctx context.Context, req *pb.DataRequest) (*pb.DataResponse, error) {
	dataResponse := pb.DataResponse{}
	if req != nil && len(req.Query) > 0 {
		client, err := opensearchconfig.GetOpensearchConnection()
		if log.CheckErrorf(err, exceptions.ErrOpenSearchConnectionInProcessDrillDownQuery) {
			return nil, err
		} else {
			computeSchedule := ComputeSchedule{}
			queryBytes, err := json.Marshal([]byte(req.Query))
			if err == nil {
				err = json.Unmarshal(queryBytes, &computeSchedule)
				if log.CheckErrorf(err, exceptions.ErrUnmarshallQuery) {
					return nil, err
				} else {
					bulkQuery := constructComputeScheduleQuery(computeSchedule)
					err := db.InsertBulkData(client, bulkQuery)
					if err != nil {
						log.Errorf(err, "Error in getting raw data :%s", err.Error())
						return nil, err
					} else {
						dataResponse.Response = "Updated the compute schedule index successfully"

					}
				}

			} else {
				log.Errorf(err, "Error in parsing the query :%s", err.Error())
			}
		}

	}
	return &dataResponse, nil
}

func (rah *ReportsHandler) GetInsightsIntegration(ctx context.Context, req *pb.CiInsightIntegrationRequest) (*pb.CiInsightIntegrationResponse, error) {
	response := pb.CiInsightIntegrationResponse{
		Endpoints: []*endpoint.Endpoint{},
	}
	resourceIds := []string{}
	if req.OrgId != constants.EMPTY_STRING {
		if len(req.ContributionIds) > 0 {
			endPointsResponse, err := helper.GetAllEndpoints(ctx, rah.endpointClient, req.OrgId, req.ContributionIds, true)
			if err != nil {
				log.CheckErrorf(err, exceptions.ErrEndpointAPIFailure)
				return &response, err
			}
			response.Endpoints = append(response.Endpoints, endPointsResponse.GetEndpoints()...)
		} else {
			cbciMap, cjocEndpointIds := make(map[string]*endpoint.Endpoint), []string{}
			contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
			endPointsResponse, err := helper.GetAllEndpoints(ctx, rah.endpointClient, req.OrgId, contributionIds, true)
			if err != nil {
				log.CheckErrorf(err, exceptions.ErrEndpointAPIFailure)
				return &response, err
			}
			endpoints := endPointsResponse.GetEndpoints()
			if len(endpoints) > 0 {
				for _, endpoint := range endpoints {
					resourceIds = append(resourceIds, endpoint.ResourceId)
					if endpoint.ContributionId == constants.CJOC_ENDPOINT {
						cjocEndpointIds = append(cjocEndpointIds, endpoint.Id)
						response.Endpoints = append(response.Endpoints, endpoint)
					} else if endpoint.ContributionId == constants.CBCI_ENDPOINT {
						found := false
						for _, property := range endpoint.Properties {
							if stringValue, ok := property.Value.(*api.Property_String_); ok {
								if property.Name == constants.TOOL_URL {
									toolUrl := stringValue.String_
									if strings.TrimSpace(toolUrl) != constants.EMPTY_STRING {
										cbciMap[toolUrl] = endpoint
										found = true
									}
								}
							}
						}
						if !found {
							response.Endpoints = append(response.Endpoints, endpoint)
						}
					} else {
						response.Endpoints = append(response.Endpoints, endpoint)
					}
				}
			}
			if len(cjocEndpointIds) != 0 {
				replacements := map[string]any{
					"subOrgId":  req.OrgId,
					"parentIds": resourceIds,
				}
				controllerUrlMap, err := getControllerUrlMap(replacements, ctx)
				log.CheckErrorf(err, "Controller fetch from opensearch failed")
				if len(controllerUrlMap) > 0 {
					for _, cjocEndpointId := range cjocEndpointIds {
						if controllerUrls, ok := controllerUrlMap[cjocEndpointId]; ok {
							for _, controllerUrl := range controllerUrls {
								delete(cbciMap, controllerUrl)
							}
						}
					}
				}
			}
			for _, endpoint := range cbciMap {
				response.Endpoints = append(response.Endpoints, endpoint)
			}
		}
	}
	sort.Slice(response.Endpoints, func(i, j int) bool {
		return response.Endpoints[i].Name < response.Endpoints[j].Name
	})
	return &response, nil
}

func (rah *ReportsHandler) sendGetUserPreferences(ctx context.Context, req *endpoint.GetUserPreferencesRequest) (*endpoint.GetUserPreferencesResponse, error) {
	resp := &endpoint.GetUserPreferencesResponse{}
	ctx = srvauth.SystemUserCtx(ctx, "reports-service.sendGetUserPreferences")
	err := rah.client.SendGrpcCtx(ctx, hostflags.EndpointServiceHost(), endpointService, getUserPreferencesMethod, req, resp)
	return resp, err
}

func (rah *ReportsHandler) sendUpdateUserPreference(ctx context.Context, ep *endpoint.UpdateUserPreferencesRequest) (*emptypb.Empty, error) {
	upResp := &emptypb.Empty{}
	err := rah.client.SendGrpcCtx(ctx, hostflags.EndpointServiceHost(), endpointService, updateUserPreferenceMethod, ep, upResp)
	return upResp, err
}

func (rah *ReportsHandler) UpdateDashboardLayout(ctx context.Context, req *pb.DashboardLayoutRequest) (*pb.DashboardLayoutResponse, error) {
	err := rah.verifyUserPreferencesRequest(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	response := pb.DashboardLayoutResponse{
		Status:          pb.Status_success,
		Error:           "",
		DashboardName:   req.DashboardName,
		DashboardLayout: nil,
	}

	if req.DashboardLayout == nil {
		return nil, db.ErrInvalidRequest
	}

	getReq := endpoint.GetUserPreferencesRequest{}
	getReq.Id = req.UserId

	getResp, err := rah.sendGetUserPreferences(ctx, &getReq)
	if err != nil {
		log.Errorf(err, "failed to get user preferences for request %+v", getResp)
	}

	customLayout := db.VsmDashboardLayout{OrgID: req.OrgId, UserID: req.UserId, DashboardName: req.DashboardName, DashboardLayout: pb.DashboardLayout{Xl: req.DashboardLayout.Xl, Lg: req.DashboardLayout.Lg, Md: req.DashboardLayout.Md, Sm: req.DashboardLayout.Sm, Xs: req.DashboardLayout.Xs}}

	log.Debugf("Project mappings : %+v", customLayout)

	jsonData, err := json.Marshal(customLayout)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unable to marshal layout")
	}

	binaryData := &api.BinaryData{
		ContentType: "application/json",
		Bytes:       jsonData,
	}

	property := &api.Property{
		Name:  req.DashboardName,
		Value: &api.Property_Data{binaryData},
	}
	epReq := &endpoint.UpdateUserPreferencesRequest{}
	epReq.Id = req.UserId

	propFound := false

	for _, v := range getResp.GetProperties() {
		if v.GetName() != req.DashboardName {
			epReq.Properties = append(epReq.Properties, v)
		} else {
			epReq.Properties = append(epReq.Properties, property)
			propFound = true
		}
	}

	if !propFound {
		epReq.Properties = append(epReq.Properties, property)
	}

	resp, err := rah.sendUpdateUserPreference(ctx, epReq)
	if err != nil {
		log.Errorf(err, "failed to update user preference endpoint %+v", resp)
		return nil, status.Error(codes.InvalidArgument, "unable to update user preference")
	}

	return &response, nil
}

func (rah *ReportsHandler) GetDashboardLayout(ctx context.Context, req *pb.DashboardLayoutRequest) (*pb.DashboardLayoutResponse, error) {
	err := rah.verifyUserPreferencesRequest(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	response := pb.DashboardLayoutResponse{
		Status:            pb.Status_success,
		Error:             "",
		DashboardName:     req.DashboardName,
		DashboardLayout:   nil,
		DefaultLayout:     false,
		DisplayTransition: false,
		ComponentId:       req.Component,
	}

	getReq := endpoint.GetUserPreferencesRequest{}
	getReq.Id = req.UserId

	// Set Transition data
	switch req.DashboardName {
	case constants.SDA_DASHBOARD:
		response.TransitionField = "softwareDeliveryActivity"
		response.DisplayTransition, response.TransitionData = rah.GetSDATransitionConfig(ctx, req)
	case constants.SECURITY_INSIGHTS_DASHBOARD:
		response.TransitionField = "securityInsights"
		response.DisplayTransition, response.TransitionData = rah.GetSiTransitionConfig(ctx, req)
	case constants.CI_INSIGHTS_DASHBOARD:
		response.TransitionField = "ciInsights"
		response.DisplayTransition, response.TransitionData = rah.GetCiTransitionConfig(ctx, req)
	case constants.FLOW_METRICS_DASHBOARD:
		response.TransitionField = "flowMetrics"
		response.DisplayTransition, response.TransitionData = rah.GetFlowMetricsTransitionConfig(ctx, req)
	case constants.DORA_METRICS_DASHBOARD:
		response.TransitionField = "doraMetrics"
		response.DisplayTransition, response.TransitionData = rah.GetDoraMetricsTransitionConfig(ctx, req)
	case constants.TEST_INSIGHTS_DASHBOARD:
		response.TransitionField = "testInsights"
		response.DisplayTransition, response.TransitionData = rah.GetTestInsightsTransitionConfig(ctx, req)
	case constants.COMPONENT_SUMMARY_DASHBOARD:
		response.TransitionField = "componentSummary"
		response.DisplayTransition, response.TransitionData = rah.GetComponentSummaryTransitionConfig(ctx, req)
	}

	resp, err := rah.sendGetUserPreferences(ctx, &getReq)
	if err != nil {
		log.Errorf(err, "failed to get user preferences for request %+v", resp)
	} else {

		dashboardLayout := db.VsmDashboardLayout{}

		layoutPreference := helper.GetDashboardLayoutForEndpoint(resp, req.DashboardName)

		if layoutPreference != nil {

			err = json.Unmarshal(layoutPreference.GetBytes(), &dashboardLayout)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "unable to unmarshal layout binary data")
			} else {
				// Check if layout of screens are empty and return default layout if true
				// This condition handles restroe default layout functionality
				if !helper.CheckEmptyDashboardLayout(&dashboardLayout) {
					response.DashboardLayout = &pb.DashboardLayout{Xl: dashboardLayout.DashboardLayout.Xl, Lg: dashboardLayout.DashboardLayout.Lg, Xs: dashboardLayout.DashboardLayout.Xs, Md: dashboardLayout.DashboardLayout.Md}
					return &response, nil
				}
			}
		}
	}

	//Get the Default layout if no custom layout found

	fileName := ""

	if strings.Compare(req.DashboardName, constants.SDA_DASHBOARD) == 0 {
		fileName = "customDashboard/sdaDashboard.json"
	} else if strings.Compare(req.DashboardName, constants.SECURITY_INSIGHTS_DASHBOARD) == 0 {
		fileName = "customDashboard/securityInsightsDashboard.json"
	} else if strings.Compare(req.DashboardName, constants.FLOW_METRICS_DASHBOARD) == 0 {
		fileName = "customDashboard/flowMetricsDashboard.json"
	} else if strings.Compare(req.DashboardName, constants.DORA_METRICS_DASHBOARD) == 0 {
		fileName = "customDashboard/doraMetricsDashboard.json"
	} else if strings.Compare(req.DashboardName, constants.TEST_INSIGHTS_DASHBOARD) == 0 {
		fileName = "customDashboard/testInsightsDashboard.json"
	} else if strings.Compare(req.DashboardName, constants.COMPONENT_SUMMARY_DASHBOARD) == 0 {
		fileName = "customDashboard/componentSummarySonarQube.json"
	} else if strings.Compare(req.DashboardName, constants.CI_INSIGHTS_DASHBOARD) == 0 {
		fileName = "customDashboard/sdaDashboard.json"
	} else if strings.Compare(req.DashboardName, constants.COMPONENT_SECURITY_DASHBOARD) == 0 {
		fileName = "customDashboard/componentSecurityDashboard.json"
	} else if strings.Compare(req.DashboardName, constants.APPLICATION_SECURITY_DASHBOARD) == 0 {
		fileName = "customDashboard/applicationSecurityDashboard.json"
	} else {
		return nil, db.ErrInvalidRequest
	}

	// open the JSON file
	jsonFile, err := os.Open(config.Config.GetString("report.definition.filepath") + fileName)
	if err != nil {
		log.CheckErrorf(err, "error opening widget definition json file ")
		return nil, db.ErrFileNotFound
	}
	defer jsonFile.Close()

	// read from the JSON file
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.CheckErrorf(err, "error reading widget definition json file: ")
		return nil, db.ErrFileNotFound
	}

	layout := pb.DashboardLayout{}
	json.Unmarshal(byteValue, &layout)

	response.Error = "default layout"
	response.DefaultLayout = true
	response.DashboardLayout = &layout

	return &response, nil
}

// Get transition flags for SDA dashboard
func (rah *ReportsHandler) GetSDATransitionConfig(ctx context.Context, req *pb.DashboardLayoutRequest) (bool, *pb.DashboardLayoutResponse_SoftwareDeliveryActivity) {
	sda := &pb.DashboardLayoutResponse_SoftwareDeliveryActivity{
		SoftwareDeliveryActivity: &pb.SoftwareDeliveryActivity{},
	}
	displayTransition := true

	// Get components flag
	components := helper.GetComponents(ctx, rah.client, rah.orgServiceClient, req.OrgId, req.UserId)
	if len(components) > 0 {
		sda.SoftwareDeliveryActivity.IsComponentCreated = true
	} else {
		return displayTransition, sda
	}

	if req.Component != "" {
		components = []string{req.Component}
	}
	wg := sync.WaitGroup{}
	wg.Add(2)

	// Get Workflow flag
	go func() {
		defer wg.Done()
		workflowsCount := helper.GetWorkflowsCount(ctx, req.OrgId, components)
		if workflowsCount > 0 {
			sda.SoftwareDeliveryActivity.IsWorkflowsFound = true
		}
	}()

	// Get Workflow runs flag
	go func() {
		defer wg.Done()
		automationRunsCount := helper.GetIndexDocCount(constants.AUTOMATION_RUN_STATUS_INDEX, req.OrgId, components)
		if automationRunsCount > 0 {
			sda.SoftwareDeliveryActivity.IsWorkflowRunsFound = true
		}
	}()

	wg.Wait()
	displayTransition = sda.SoftwareDeliveryActivity.IsComponentCreated && sda.SoftwareDeliveryActivity.IsWorkflowsFound && sda.SoftwareDeliveryActivity.IsWorkflowRunsFound

	return !displayTransition, sda
}

// Get transition flags for Security Insights dashboard
func (rah *ReportsHandler) GetSiTransitionConfig(ctx context.Context, req *pb.DashboardLayoutRequest) (bool, *pb.DashboardLayoutResponse_SecurityInsights) {
	si := &pb.DashboardLayoutResponse_SecurityInsights{
		SecurityInsights: &pb.SecurityInsights{},
	}
	displayTransition := true

	// Get components flag
	components := helper.GetComponents(ctx, rah.client, rah.orgServiceClient, req.OrgId, req.UserId)
	if len(components) > 0 {
		si.SecurityInsights.IsComponentCreated = true
	} else {
		return displayTransition, si
	}

	if req.Component != "" {
		components = []string{req.Component}
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Get Workflow flag
	go func() {
		defer wg.Done()
		workflowsCount := helper.GetWorkflowsCount(ctx, req.OrgId, components)
		if workflowsCount > 0 {
			si.SecurityInsights.IsWorkflowsFound = true
		}
	}()

	// Get Workflow runs flag
	go func() {
		defer wg.Done()
		automationRunsCount := helper.GetIndexDocCount(constants.AUTOMATION_RUN_STATUS_INDEX, req.OrgId, components)
		if automationRunsCount > 0 {
			si.SecurityInsights.IsWorkflowRunsFound = true
		}
	}()

	wg.Wait()
	displayTransition = si.SecurityInsights.IsComponentCreated && si.SecurityInsights.IsWorkflowsFound && si.SecurityInsights.IsWorkflowRunsFound

	return !displayTransition, si
}

// Get transition flags for CI Insights dashboard
func (rah *ReportsHandler) GetCiTransitionConfig(ctx context.Context, req *pb.DashboardLayoutRequest) (bool, *pb.DashboardLayoutResponse_CiInsights) {
	ci := &pb.DashboardLayoutResponse_CiInsights{
		CiInsights: &pb.CiInsights{},
	}
	displayTransition := true

	// Get ci tools flag
	contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
	endPointsResponse, err := helper.GetAllEndpoints(ctx, rah.endpointClient, req.OrgId, contributionIds, true)
	log.CheckErrorf(err, exceptions.ErrEndpointAPIFailure)

	if err == nil && len(endPointsResponse.Endpoints) > 0 {
		ci.CiInsights.IsCiToolsFound = true
	}

	// Get ci insights data flag
	// ciInsightsDataCount := helper.GetIndexDocCountByOrgId(constants.CI_INSIGHTS_INDEX, req.OrgId)
	// if ciInsightsDataCount > 0 {
	// 	ci.CiInsights.IsCiToolsFound = ci.CiInsights.IsCiToolsFound && true
	// } else {
	// 	ci.CiInsights.IsCiToolsFound = ci.CiInsights.IsCiToolsFound && false
	// }

	displayTransition = ci.CiInsights.IsCiToolsFound

	return !displayTransition, ci
}

// Get transition flags for Flow metrics dashboard
func (rah *ReportsHandler) GetFlowMetricsTransitionConfig(ctx context.Context, req *pb.DashboardLayoutRequest) (bool, *pb.DashboardLayoutResponse_FlowMetrics) {
	fm := &pb.DashboardLayoutResponse_FlowMetrics{
		FlowMetrics: &pb.FlowMetrics{},
	}
	displayTransition := true

	// Get components flag
	components := helper.GetComponents(ctx, rah.client, rah.orgServiceClient, req.OrgId, req.UserId)
	if len(components) > 0 {
		fm.FlowMetrics.IsComponentCreated = true
	}

	if req.Component != "" {
		components = []string{req.Component}
	}

	wg := sync.WaitGroup{}
	wg.Add(5)

	// Get Workflow flag
	go func() {
		defer wg.Done()
		workflowsCount := helper.GetWorkflowsCount(ctx, req.OrgId, components)
		if workflowsCount > 0 {
			fm.FlowMetrics.IsWorkflowsFound = true
		}
	}()

	// Get Workflow runs flag
	go func() {
		defer wg.Done()
		automationRunsCount := helper.GetIndexDocCount(constants.AUTOMATION_RUN_STATUS_INDEX, req.OrgId, components)
		if automationRunsCount > 0 {
			fm.FlowMetrics.IsWorkflowRunsFound = true
		}
	}()

	// Get Jira Integrations flag
	go func() {
		defer wg.Done()
		contributionIds := []string{constants.JIRA_ENDPOINT_CONTRIBUTION_ID}
		endPointsResponse, err := helper.GetAllEndpoints(ctx, rah.endpointClient, req.OrgId, contributionIds, true)
		log.CheckErrorf(err, exceptions.ErrEndpointAPIFailure)

		if err == nil && len(endPointsResponse.Endpoints) > 0 {
			displayTransition = true
			fm.FlowMetrics.IsIntegrationFound = true
		}
	}()

	// Get Analytics configuration flag
	go func() {
		defer wg.Done()
		flowMappings := helper.GetFlowItemMappings(ctx, rah.endpointClient, req.OrgId)
		if flowMappings > 0 {
			displayTransition = true
			fm.FlowMetrics.IsAnalyticsConfigFound = true
		}
	}()

	// Get flow metrics data configuration flag
	go func() {
		defer wg.Done()
		flowMetricsDataCount := helper.GetIndexDocCount(constants.FLOW_METRICS_INDEX, req.OrgId, components)
		if flowMetricsDataCount > 0 {
			fm.FlowMetrics.IsFlowMetricsDataFound = true
		}
	}()

	wg.Wait()
	displayTransition = fm.FlowMetrics.IsComponentCreated && fm.FlowMetrics.IsWorkflowsFound && fm.FlowMetrics.IsWorkflowRunsFound && fm.FlowMetrics.IsAnalyticsConfigFound && fm.FlowMetrics.IsIntegrationFound && fm.FlowMetrics.IsFlowMetricsDataFound

	return !displayTransition, fm
}

// Get transition flags for Dora metrics dashboard
func (rah *ReportsHandler) GetDoraMetricsTransitionConfig(ctx context.Context, req *pb.DashboardLayoutRequest) (bool, *pb.DashboardLayoutResponse_DoraMetrics) {
	dm := &pb.DashboardLayoutResponse_DoraMetrics{
		DoraMetrics: &pb.DoraMetrics{},
	}
	displayTransition := true

	// Get components flag
	components := helper.GetComponents(ctx, rah.client, rah.orgServiceClient, req.OrgId, req.UserId)
	if len(components) > 0 {
		dm.DoraMetrics.IsComponentCreated = true
	}

	if req.Component != "" {
		components = []string{req.Component}
	}

	var wg sync.WaitGroup
	wg.Add(4)

	// Get Workflow flag
	go func() {
		defer wg.Done()
		workflowsCount := helper.GetWorkflowsCount(ctx, req.OrgId, components)
		if workflowsCount > 0 {
			dm.DoraMetrics.IsWorkflowsFound = true
		}
	}()

	// Get Workflow runs flag
	go func() {
		defer wg.Done()
		automationRunsCount := helper.GetIndexDocCount(constants.AUTOMATION_RUN_STATUS_INDEX, req.OrgId, components)
		if automationRunsCount > 0 {
			dm.DoraMetrics.IsWorkflowRunsFound = true
		}
	}()

	// Get Environments flag
	go func() {
		defer wg.Done()
		contributionIds := []string{constants.ENVIRONMENT_ENDPOINT_CONTRIBUTION_ID}
		endPointsResponse, err := helper.GetAllEndpoints(ctx, rah.endpointClient, req.OrgId, contributionIds, true)
		log.CheckErrorf(err, exceptions.ErrEndpointAPIFailure)

		if err == nil && len(endPointsResponse.Endpoints) > 0 {
			displayTransition = true
			dm.DoraMetrics.IsEnvironmentFound = true
		}
	}()

	// Get DORA metrics data configuration flag
	go func() {
		defer wg.Done()
		doraMetricsDataCount := helper.GetIndexDocCount(constants.DORA_METRICS_INDEX, req.OrgId, components)
		if doraMetricsDataCount > 0 {
			dm.DoraMetrics.IsDeployDataFound = true
		}
	}()

	wg.Wait()
	displayTransition = dm.DoraMetrics.IsComponentCreated && dm.DoraMetrics.IsWorkflowsFound && dm.DoraMetrics.IsWorkflowRunsFound && dm.DoraMetrics.IsEnvironmentFound && dm.DoraMetrics.IsDeployDataFound

	return !displayTransition, dm
}

// Get transition flags for Test Insights dashboard
func (rah *ReportsHandler) GetTestInsightsTransitionConfig(ctx context.Context, req *pb.DashboardLayoutRequest) (bool, *pb.DashboardLayoutResponse_TestInsights) {

	ti := &pb.DashboardLayoutResponse_TestInsights{
		TestInsights: &pb.TestInsights{},
	}

	displayTransition := true

	// Get components flag
	components := helper.GetComponents(ctx, rah.client, rah.orgServiceClient, req.OrgId, req.UserId)

	if len(components) > 0 {
		ti.TestInsights.IsComponentCreated = true
	}

	if req.Component != "" {
		components = []string{req.Component}
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// Get Workflow flag
	go func() {
		defer wg.Done()
		workflowsCount := helper.GetWorkflowsCount(ctx, req.OrgId, components)
		if workflowsCount > 0 {
			ti.TestInsights.IsWorkflowsFound = true
		}
	}()

	// Get Workflow runs flag
	go func() {
		defer wg.Done()
		automationRunsCount := helper.GetIndexDocCount(constants.AUTOMATION_RUN_STATUS_INDEX, req.OrgId, components)
		if automationRunsCount > 0 {
			ti.TestInsights.IsWorkflowRunsFound = true
		}
	}()

	// Get TestInsightsData flag
	go func() {
		defer wg.Done()
		testInsightsCount := helper.GetIndexDocCount(constants.TEST_SUITE_INDEX, req.OrgId, components)
		if testInsightsCount > 0 {
			ti.TestInsights.IsTestInsightsDataFound = true
		}
	}()

	wg.Wait()
	displayTransition = ti.TestInsights.IsComponentCreated && ti.TestInsights.IsWorkflowsFound && ti.TestInsights.IsWorkflowRunsFound && ti.TestInsights.IsTestInsightsDataFound

	return !displayTransition, ti

}

// Get transition flags for Component summary dashboard
func (rah *ReportsHandler) GetComponentSummaryTransitionConfig(ctx context.Context, req *pb.DashboardLayoutRequest) (bool, *pb.DashboardLayoutResponse_ComponentSummary) {
	cs := &pb.DashboardLayoutResponse_ComponentSummary{
		ComponentSummary: &pb.ComponentSummary{},
	}
	displayTransition := true

	// Get components flag
	cs.ComponentSummary.IsComponentCreated = true
	components := make([]string, 0)

	if req.Component != "" {
		components = []string{req.Component}
	} else {
		return displayTransition, cs
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Get Workflow flag
	go func() {
		defer wg.Done()
		workflowsCount := helper.GetWorkflowsCount(ctx, req.OrgId, components)
		if workflowsCount > 0 {
			cs.ComponentSummary.IsWorkflowsFound = true
		}
	}()

	// Get Workflow runs flag
	go func() {
		defer wg.Done()
		automationRunsCount := helper.GetIndexDocCount(constants.AUTOMATION_RUN_STATUS_INDEX, req.OrgId, components)
		if automationRunsCount > 0 {
			cs.ComponentSummary.IsWorkflowRunsFound = true
		}
	}()

	wg.Wait()
	displayTransition = cs.ComponentSummary.IsComponentCreated && cs.ComponentSummary.IsWorkflowsFound && cs.ComponentSummary.IsWorkflowRunsFound

	return !displayTransition, cs
}

func (rah *ReportsHandler) StreamCIInsightsCompletedRun(req *pb.ReportServiceRequest, srv pb.ReportServiceHandler_StreamCIInsightsCompletedRunServer) error {
	defer rah.RecordTiming("StreamCIInsightsCompletedRun", time.Now())

	ctx := context.Background()

	replacements := map[string]any{
		"startDate":  req.StartDate,
		"endDate":    req.EndDate,
		"orgId":      req.OrgId,
		"subOrgId":   req.SubOrgId,
		"component":  req.Component,
		"ciToolId":   req.CiToolId,
		"ciToolType": req.CiToolType,
		"sortBy":     req.SortBy,
		"filterType": req.FilterType,
		"viewOption": req.ViewOption,
		"timeZone":   req.TimeZone,
	}

	if req.TimeZone == "" {
		replacements["timeZone"] = constants.LOCATION_EUROPE_OR_LONDON
	}

	dateBydurationType := models.CalculateDateBydurationType{
		StartDateStr:       req.StartDate,
		EndDateStr:         req.EndDate,
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   map[string]any{},
		NormalizeMonthFlag: false,
		IsComputFlag:       false,
		CurrentTime:        time.Now().UTC(),
	}
	calculateDateBydurationType(dateBydurationType)

	err := internal.GetInsightCompletedRunsStream(replacements, ctx, rah.endpointClient, srv)
	log.CheckErrorf(err, "could not process the completed runs stream :")
	return err
}

func (rah *ReportsHandler) RecordTiming(metric string, started time.Time) {
	m := rah.MetricMap()
	us := float64(time.Now().Sub(started) / time.Microsecond) // microseconds
	avg := m.NewAverage(metric+".avg", "ms", 30, true)
	avg.AddValue(us / 1000.0) // milliseconds
}

var (
	getControllerUrlMap      = internal.GetControllerUrlMap
	jobAndRunCount           = internal.JobAndRunCount
	getVersionAndPluginCount = internal.GetVersionAndPluginCount
)

func (rah *ReportsHandler) GetControllersInfo(ctx context.Context, req *pb.CiControllerInfoRequest) (*pb.CiControllerInfoResponse, error) {
	response := &pb.CiControllerInfoResponse{}
	endPointsResponse := &endpoint.EndpointsResponse{}
	resourceIds := []string{}
	endpointIds := []string{}
	cjocMap := make(map[string][]string)
	cbciMap, cjocEndpointIds := make(map[string]*pb.CiControllerInfo), []string{}
	var endpoints []*endpoint.Endpoint
	var err error
	if req.OrgId != constants.EMPTY_STRING {
		if len(req.ContributionIds) > 0 {
			endPointsResponse, err = helper.GetAllEndpoints(ctx, rah.endpointClient, req.OrgId, req.ContributionIds, true)
			if err != nil {
				log.CheckErrorf(err, exceptions.ErrEndpointAPIFailure)
				return response, err
			}
			endpoints = endPointsResponse.GetEndpoints()
		} else {
			contributionIds := constants.INSIGHTS_SUPPORTED_TYPE
			endPointsResponse, err = helper.GetAllEndpoints(ctx, rah.endpointClient, req.OrgId, contributionIds, true)
			if err != nil {
				log.CheckErrorf(err, exceptions.ErrEndpointAPIFailure)
				return response, err
			}
			endpoints = endPointsResponse.GetEndpoints()
		}
		if len(endpoints) > 0 {
			for _, endpoint := range endpoints {
				endpointIds = append(endpointIds, endpoint.GetId())
				resourceIds = append(resourceIds, endpoint.GetResourceId())
				controllerInfo := &pb.CiControllerInfo{
					Endpoint:   endpoint,
					RunsInfo:   &pb.RunsInfo{},
					PluginInfo: &pb.PluginInfo{},
				}
				if req.GetOrgId() == endpoint.GetResourceId() {
					controllerInfo.Source = "Original"
				} else {
					controllerInfo.Source = "Inherited"
				}
				if endpoint.ContributionId != constants.CBCI_ENDPOINT {
					response.Controllers = append(response.Controllers, controllerInfo)
				}
				if endpoint.ContributionId == constants.CJOC_ENDPOINT {
					cjocEndpointIds = append(cjocEndpointIds, endpoint.Id)
				} else if endpoint.ContributionId == constants.CBCI_ENDPOINT {
					for _, property := range endpoint.GetProperties() {
						if stringValue, ok := property.Value.(*api.Property_String_); ok {
							if property.Name == constants.TOOL_URL {
								toolUrl := stringValue.String_
								if strings.TrimSpace(toolUrl) != constants.EMPTY_STRING {
									cbciMap[toolUrl] = controllerInfo
								} else {
									response.Controllers = append(response.Controllers, controllerInfo)
								}
							}
						}
					}
				}
			}
		}
		if len(cjocEndpointIds) != 0 {
			replacements := map[string]any{
				"subOrgId":  req.OrgId,
				"parentIds": resourceIds,
			}
			controllerUrlMap, err := getControllerUrlMap(replacements, ctx)
			log.CheckErrorf(err, "Controller fetch from opensearch failed")
			if len(controllerUrlMap) > 0 {
				for _, cjocEndpointId := range cjocEndpointIds {
					if controllerUrls, ok := controllerUrlMap[cjocEndpointId]; ok {
						for _, controllerUrl := range controllerUrls {
							if cbciMap[controllerUrl] != nil {
								resp := cbciMap[controllerUrl]
								resp.ParentId = cjocEndpointId
								cbciMap[controllerUrl] = resp
								cjocMap[cjocEndpointId] = append(cjocMap[cjocEndpointId], resp.Endpoint.Id)
							}
						}
					}
				}
			}
		}
		for _, value := range cbciMap {
			response.Controllers = append(response.Controllers, value)
		}
	}

	replacements := map[string]any{
		"endpointIds": endpointIds,
		"parentIds":   resourceIds,
	}
	jobCountMap, runCountMap, err := jobAndRunCount(replacements, ctx)
	log.CheckErrorf(err, "Job and run info fetch from opensearch failed")

	endpointVersionPluginCount, errPluginCount := getVersionAndPluginCount(replacements, ctx)
	log.CheckErrorf(errPluginCount, "Version and plugin count fetch from opensearch failed")
	for _, controller := range response.Controllers {
		if controller.Endpoint.ContributionId == constants.CJOC_ENDPOINT {
			resp := cjocMap[controller.Endpoint.Id]
			if resp != nil {
				jobCount, runCount := 0, 0
				for _, cbci := range resp {
					job, jobExists := jobCountMap[cbci]
					run, runExists := runCountMap[cbci]
					if jobExists {
						jobCount += job
					}
					if runExists {
						runCount += run
					}
				}
				controller.RunsInfo.TotalNumberOfProjects = float64(jobCount)
				controller.RunsInfo.TotalNumberOfRuns = float64(runCount)
			}
		}
		if endpointVersionPluginCount != nil && endpointVersionPluginCount[controller.Endpoint.Id] != nil {
			versionPluginCount := endpointVersionPluginCount[controller.Endpoint.Id].(map[string]interface{})
			if versionPluginCount != nil {
				controller.Version = versionPluginCount[constants.VERSION].(string)
				controller.PluginInfo.Count = versionPluginCount[constants.PLUGIN_COUNT].(float64)
			}
		}
		jobCount, jobExists := jobCountMap[controller.Endpoint.Id]
		runCount, runExists := runCountMap[controller.Endpoint.Id]
		if jobExists {
			controller.RunsInfo.TotalNumberOfProjects = float64(jobCount)
		}
		if runExists {
			controller.RunsInfo.TotalNumberOfRuns = float64(runCount)
		}
	}
	sort.Slice(response.Controllers, func(i, j int) bool {
		return response.Controllers[i].Endpoint.Audit.When.AsTime().After(response.Controllers[j].Endpoint.Audit.When.AsTime())
	})
	return response, nil
}
