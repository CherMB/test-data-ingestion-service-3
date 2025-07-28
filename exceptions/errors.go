package exceptions

import "errors"

// error description
const (
	ErrGetDeploymentTypeAssertion              string = "error in type assertion in getDeployments"
	ErrGetCommitTrendTypeAssertion             string = "error in type assertion in getCommitTrends"
	ErrGetAutomationRunTypeAssertion           string = "error in type assertion in getAutomationRuns"
	ErrGetPullRequestTypeAssertion             string = "error in type assertion in getPullRequests"
	ErrGetWorkloadTypeAssertion                string = "error in type assertion in getWorkload"
	ErrGetCycleTimeTypeAssertion               string = "error in type assertion in getCycleTime"
	ErrGetVulnerabilitiesOverviewTypeAssertion string = "error in type assertion in getVulnerabilitiesOverview"
	ErrGetWorkEfficiencyTypeAssertion          string = "error in type assertion in getWorkEfficiency"
	ErrGetVelocityTypeAssertion                string = "error in type assertion in getVelocity"
	ErrGetWorkItemDistributionTypeAssertion    string = "error in type assertion in getWorkItemDistribution"
	ErrGetMttrForVulnerabilitiesTypeAssertion  string = "error in type assertion in getMttrForVulnerabilities"
	ErrGetCodeChurnTypeAssertion               string = "error in type assertion in getCodeChurn"
	ErrEndpointAPIFailure                      string = "endpoint api failed"

	ErrUnmarshallingGetVulnerabilitiesOverview                    string = "Error unmarshaling response getVulnerabilitiesOverview()"
	ErrUnmarshallingGetWorkEfficiency                             string = "Error unmarshaling response getWorkEfficiency()"
	ErrUnmarshallingGetVulnerablitiesByScanType                   string = "Error unmarshaling response getVulnerabilitiesByScanType()"
	ErrUnmarshallingGetMttrForVulnerabilities                     string = "Error unmarshaling response getMttrForVulnerabilities()"
	ErrMarshallingGetMttrForVulnerabilities                       string = "Error marshaling responseStruct in getMttrForVulnerabilities() chart"
	ErrUnmarshallingOpenSearchRespInTestAutomationDrillDown       string = "Error unmarshaling response from OpenSearch in TestAutomationDrilldown()"
	ErrMarshallingRespInSecurityComponentDrillDown                string = "error marshaling reponse in SecurityComponentDrillDown() :"
	ErrUnmarshallingRespInSecurityComponentDrillDown              string = "error unmarshaling reponse in SecurityComponentDrillDown() :"
	ErrUnmarshallingTransformFindingsIdentifiedSince              string = "Error unmarshaling response transformFindingsIdentifiedSince()"
	ErrUnmarshallingTransformRiskAcceptedAndFalsePositiveFindings string = "Error unmarshaling response transformRiskAcceptedAndFalsePositiveFindings()"
	ErrUnmarshallingTransformSlaBreachesBySeverity                string = "Error unmarshaling response transformSlaBreachesBySeverity()"
	ErrUnmarshallingTransformOpenFindingsByReviewStatus           string = "Error unmarshaling response transformOpenFindingsByReviewStatus()"

	ErrQueryFailureToGetResponse                   string = "failed to get query response"
	ErrMarshallQuery                               string = "could not marshal query :"
	ErrUnmarshallQuery                             string = "could not unmarshal queryBytes :"
	ErrJsonConversion                              string = "error converting to json"
	ErrExecutePostProcess                          string = "error in ExecutePostProcessFunction : %s"
	ErrOpenSearchConnection                        string = "Error establishing connection with OpenSearch in getQueryResponse(). Connection error - "
	ErrOpenSearchConnectionInProcessDrillDownQuery string = "Error establishing connection with OpenSearch in processDrilldownQueryAndSpec()"
	ErrJsonPlaceholderNotReplaceable               string = "could not replace json placeholders :"
	ErrOpenSearchFetchDataFailure                  string = "Error fetching Opensearch data in service.getResponse(). Fetch failed - "
	ErrParsingStartDate                            string = "Error parsing start date :"
	ErrDefaultRespTemplate                         string = "Error response :"
	ErrFetchingServiceTemplate                     string = "Exception while fetching services : %s : "
	ErrSonarGetData                                string = "Error in sonar get data : "
	ErrUnmarshallRespGetCommitTrends               string = "Error unmarshaling response getCommitTrends()"
	ErrOrgIdNotSpecified                           string = "Org id must be specified"
	ErrReportServiceReqValidationFailure           string = "ReportServiceRequest Validation failed: %s"

	ErrFetchingSubOrgForOrg                       string = "Exception while fetching sub orgs for org : %s : %s"
	ErrFetchingServicesForOrg                     string = "Exception while fetching services for org : %s : %s"
	ErrFetchingServicesForSubOrg                  string = "Exception while fetching services for sub org : %s : %s"
	ErrFetchingServicesForSubOrgWithoutFormatting string = "Exception while fetching sub org services"
	ErrFetchingEndpoints                          string = "Exception while fetching endpoints"
	ErrInvalidFormatStartDate                     string = "Start date invalid format"
	ErrInvalidFormatEndDate                       string = "End date invalid format"
	ErrConvertingStartTimeToTZ                    string = "Error converting start time - %s to timezone - %s"
	ErrConvertingEndTimeToTZ                      string = "Error converting end time - %s to timezone - %s"

	// helper errors
	ErrHelperAddBuckets                      string = "error in helper.AddDateBuckets()"
	ErrReadingWidgetDefJsonFile              string = "error reading widget definition json file: "
	ErrOpeningWidgetDefJsonFile              string = "error opening widget definition json file "
	ErrSearchingDocInHelperGetOpenSearchData string = "error searching document in helpers.GetOpenSearchData()"
)

// debug logs
const (
	DebugNilSonarObject                    string = "sonar object is nil"
	DebugEmptyComponentData                string = "Component Data is not present"
	DebugAutomationWidgetSectionParams     string = "Automation widget section params "
	DebugEmptyAutomationData               string = "Automation Data is not present"
	DebugActiveAutomations                 string = "Active automations : %v"
	DebugEmptyComponentDataForWidgetID     string = "Component Data is not present for widget Id: %s"
	DebugAddingBranchFilterWithBranchName  string = "Adding Branch Filter for branch:%s, BranchName:%v"
	DebugAddingBranchFilterWithBranchId    string = "Adding branch filter with branch ID: %s"
	DebugTimeTookToFetchAllServiceMilliSec string = "Time took to fectch all services : %v in milliseconds"
	DebugApplyingComponentsFilter          string = "Applying components filter : %s"
	DebugComponentInReqForWidget           string = "Components in request %v for widget %s"
	DebugEndpointList                      string = "Endpoints list : %v"
	DebugRequestList                       string = "Request  :"

	//handler debug
	DebugTimeTookToFetchAllAutomationMilliSec                            string = "Time took to fetch all automation and last active time for automation drilldown : %v in milliseconds"
	DebugTimeTookToFetchAllAutomationRunMilliSec                         string = "Time took to fetch all automation and last active time for automation run drilldown : %v in milliseconds"
	DebugTimeTookToProcessAllAutomationMilliSec                          string = "Time took to process all auitomation and last active time for automation drilldown : %v in milliseconds"
	DebugTimeTookToProcessAllAutomationRunMilliSec                       string = "Time took to process all auitomation and last active time for automation run drilldown : %v in milliseconds"
	DebugTimeTookToFetchAllAutomationInfoFromCacheMilliSec               string = "Time took to automation info from cache for automation drilldown : %v in milliseconds"
	DebugTimeTookToFetchAllAutomationRunInfoFromCacheMilliSec            string = "Time took to automation info from cache for automation run drilldown : %v in milliseconds"
	DebugTimeTookToProcessAllAutomationFromCacheAndOpensearchMilliSec    string = "Time took to process all automation from cache and opensearch for automation drilldown : %v in milliseconds"
	DebugTimeTookToProcessAllAutomationRunFromCacheAndOpensearchMilliSec string = "Time took to process all automation from cache and opensearch for automation run drilldown : %v in milliseconds"
	DebugTimeTookToSortAllResultForAutomationMilliSec                    string = "Time took to sort the result for automation drilldown : %v in milliseconds"
	DebugTimeTookToSortAllResultForAutomationRunMilliSec                 string = "Time took to sort the result for automation run drilldown : %v in milliseconds"
)

var errMap map[string]error = map[string]error{

	ErrGetDeploymentTypeAssertion:              errors.New(ErrGetDeploymentTypeAssertion),
	ErrGetCommitTrendTypeAssertion:             errors.New(ErrGetCommitTrendTypeAssertion),
	ErrGetAutomationRunTypeAssertion:           errors.New(ErrGetAutomationRunTypeAssertion),
	ErrGetPullRequestTypeAssertion:             errors.New(ErrGetPullRequestTypeAssertion),
	ErrGetVulnerabilitiesOverviewTypeAssertion: errors.New(ErrGetVulnerabilitiesOverviewTypeAssertion),
	ErrGetWorkloadTypeAssertion:                errors.New(ErrGetWorkloadTypeAssertion),
	ErrGetCycleTimeTypeAssertion:               errors.New(ErrGetCycleTimeTypeAssertion),
	ErrGetWorkEfficiencyTypeAssertion:          errors.New(ErrGetWorkEfficiencyTypeAssertion),
	ErrGetVelocityTypeAssertion:                errors.New(ErrGetVelocityTypeAssertion),
	ErrGetWorkItemDistributionTypeAssertion:    errors.New(ErrGetWorkItemDistributionTypeAssertion),
	ErrGetMttrForVulnerabilitiesTypeAssertion:  errors.New(ErrGetMttrForVulnerabilitiesTypeAssertion),
	ErrGetCodeChurnTypeAssertion:               errors.New(ErrGetCodeChurnTypeAssertion),
	ErrEndpointAPIFailure:                      errors.New(ErrEndpointAPIFailure),
}

func GetExceptionByCode(key string) error {
	val, ok := errMap[key]
	if ok {
		return val
	}
	return errors.ErrUnsupported
}
