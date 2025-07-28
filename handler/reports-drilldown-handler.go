package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	api "github.com/calculi-corp/api/go"
	pb "github.com/calculi-corp/api/go/vsm/report"
	cutils "github.com/calculi-corp/common/utils"
	client "github.com/calculi-corp/grpc-client"
	"github.com/calculi-corp/log"
	opensearchconfig "github.com/calculi-corp/opensearch-config"
	"github.com/calculi-corp/reports-service/cache"
	"github.com/calculi-corp/reports-service/constants"
	"github.com/calculi-corp/reports-service/db"
	"github.com/calculi-corp/reports-service/exceptions"
	helper "github.com/calculi-corp/reports-service/helper"
	"github.com/calculi-corp/reports-service/internal"
	"github.com/opensearch-project/opensearch-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type ColorScheme struct {
	Color0 string `json:"color0"`
	Color1 string `json:"color1"`
}

type TestWorkflowReportResponse struct {
	Branch        string `json:"branch"`
	ComponentName string `json:"componentName"`
	ComponentId   string `json:"componentId"`
	WorkflowName  string `json:"workflow"`
	RunIdCount    int    `json:"workflowRuns"`
	TestSuiteType string `json:"testSuiteType"`
	TestSuites    int64  `json:"testSuites"`
	Source        string `json:"source"`
}

type TestComponentReportResponse struct {
	ComponentId   string `json:"componentId"`
	ComponentName string `json:"componentName"`
	RepositoryUrl string `json:"repositoryUrl"`
	TestSuiteType string `json:"testSuiteType"`
	TestSuites    int64  `json:"testSuites"`
}

type ComponentReportResponse struct {
	ComponentId        string `json:"componentId"`
	ComponentName      string `json:"componentName"`
	RepositoryUrl      string `json:"repositoryUrl"`
	Status             string `json:"status"`
	LastActive         string `json:"lastActive"`
	LastActiveInMillis int64  `json:"lastActiveInMillis"`
}

type WorkflowReportResponse struct {
	ComponentId        string `json:"componentId"`
	ComponentName      string `json:"componentName"`
	WorkflowName       string `json:"workflowName"`
	AutomationId       string `json:"automationId"`
	BranchId           string `json:"branchId"`
	Branch             string `json:"branch"`
	Status             string `json:"status"`
	LastActive         string `json:"lastActive"`
	LastActiveInMillis int64  `json:"lastActiveInMillis"`
	Source             string `json:"source"`
}

type WorkflowRunReportResponse struct {
	ComponentId          string  `json:"componentId"`
	ComponentName        string  `json:"componentName"`
	WorkflowName         string  `json:"workflowName"`
	AutomationId         string  `json:"automationId"`
	BranchId             string  `json:"branchId"`
	Branch               string  `json:"branch"`
	Build                float64 `json:"build"`
	RunId                string  `json:"runId"`
	Status               string  `json:"status"`
	RunStartTime         string  `json:"runStartTime"`
	RunStartTimeInMillis int64   `json:"runStartTimeInMillis"`
	Source               string  `json:"source"`
}

type CommitRunResponse struct {
	OrgName             string   `json:"org_name"`
	ComponentId         string   `json:"componentId"`
	ComponentName       string   `json:"componentName"`
	WorkflowName        string   `json:"workflowName"`
	AutomationId        string   `json:"automationId"`
	BranchId            string   `json:"branchId"`
	Branch              string   `json:"branch"`
	RunId               string   `json:"runId"`
	RunNumber           float64  `json:"runNumber"`
	CommitId            string   `json:"commitId"`
	CommitDescription   string   `json:"commitDescription"`
	Status              string   `json:"status"`
	Duration            float64  `json:"duration"`
	CreatedTime         string   `json:"createdTime"`
	CreatedTimeInMillis int64    `json:"createdTimeInMillis"`
	EndTime             string   `json:"endTime"`
	EndTimeInMillis     int64    `json:"endTimeInMillis"`
	Environment         []string `json:"environment"`
}

type BuildsResponse struct {
	OrgName           string  `json:"org_name"`
	ComponentId       string  `json:"componentId"`
	ComponentName     string  `json:"componentName"`
	WorkflowName      string  `json:"workflowName"`
	AutomationId      string  `json:"automationId"`
	BranchId          string  `json:"branchId"`
	Branch            string  `json:"branch"`
	RunId             string  `json:"runId"`
	Build             float64 `json:"build"`
	JobId             string  `json:"jobId"`
	StepId            string  `json:"stepId"`
	Status            string  `json:"status"`
	StartTime         string  `json:"startTime"`
	StartTimeInMillis int64   `json:"startTimeInMillis"`
	EndTime           string  `json:"endTime"`
	EndTimeInMillis   int64   `json:"endTimeInMillis"`
	Duration          float64 `json:"duration"`
	Environment       string  `json:"environment"`
	Source            string  `json:"source"`
}

type DeploymentResponse struct {
	OrgName           string  `json:"org_name"`
	ComponentId       string  `json:"componentId"`
	ComponentName     string  `json:"componentName"`
	WorkflowName      string  `json:"workflowName"`
	AutomationId      string  `json:"automationId"`
	BranchId          string  `json:"branchId"`
	Branch            string  `json:"branch"`
	RunId             string  `json:"runId"`
	Build             float64 `json:"build"`
	JobId             string  `json:"jobId"`
	StepId            string  `json:"stepId"`
	Status            string  `json:"status"`
	Duration          float64 `json:"duration"`
	StartTime         string  `json:"startTime"`
	StartTimeInMillis int64   `json:"startTimeInMillis"`
	EndTime           string  `json:"endTime"`
	EndTimeInMillis   int64   `json:"endTimeInMillis"`
	Environment       string  `json:"environment"`
}

type CommitsResponse struct {
	ComponentId             string `json:"componentId"`
	ComponentName           string `json:"componentName"`
	Author                  string `json:"author"`
	Branch                  string `json:"branch"`
	CommitId                string `json:"commitId"`
	RepositoryName          string `json:"repositoryName"`
	CommitTimestamp         string `json:"commitTimestamp"`
	CommitTimestampInMillis int64  `json:"commitTimestampInMillis"`
}

type PullRequestResponse struct {
	ComponentId         string `json:"componentId"`
	ComponentName       string `json:"componentName"`
	PrId                string `json:"prID"`
	Provider            string `json:"provider"`
	SourceBranch        string `json:"sourceBranch"`
	TargetBranch        string `json:"targetBranch"`
	RepositoryName      string `json:"repositoryName"`
	Status              string `json:"status"`
	CreatedTime         string `json:"createdTime"`
	CreatedTimeInMillis int64  `json:"createdTimeInMillis"`
}

type SecurityComponentReportResponse struct {
	ComponentId   string        `json:"componentId"`
	ComponentName string        `json:"componentName"`
	RepositoryUrl string        `json:"repositoryUrl"`
	Scanners      []interface{} `json:"scannerName"`
	ScannerType   string        `json:"scanners"`
}

type SecurityWorkflowReportResponse struct {
	ComponentId   string        `json:"componentId"`
	ComponentName string        `json:"component"`
	WorkflowName  string        `json:"workflow"`
	Branch        string        `json:"branch"`
	Scanners      []interface{} `json:"scannerName"`
	ScannerType   string        `json:"scanners"`
	RunIdCount    int           `json:"workflowRuns"`
}

type SecurityWorkflowRunReportResponse struct {
	ComponentId          string        `json:"componentId"`
	ComponentName        string        `json:"component"`
	WorkflowName         string        `json:"workflow"`
	AutomationId         string        `json:"automationId"`
	RunId                string        `json:"runId"`
	Branch               string        `json:"branch"`
	BranchId             string        `json:"branchId"`
	Build                float64       `json:"build"`
	Status               string        `json:"status"`
	ScanStatus           string        `json:"scanStatus"`
	Scanners             []interface{} `json:"scannerName"`
	ScannerType          string        `json:"scanners"`
	RunStartTimeInMillis int64         `json:"runStartTimeInMillis"`
}

type TestWorkflowRunReportResponse struct {
	ComponentId          string  `json:"componentId"`
	ComponentName        string  `json:"componentName"`
	WorkflowName         string  `json:"workflow"`
	AutomationId         string  `json:"automationId"`
	RunId                string  `json:"runId"`
	Branch               string  `json:"branch"`
	BranchId             string  `json:"branchId"`
	Build                float64 `json:"build"`
	Status               string  `json:"status"`
	RunStatus            string  `json:"runStatus"`
	TestSuiteType        string  `json:"testSuiteType"`
	RunStartTimeInMillis int     `json:"runStartTimeInMillis"`
	TestSuites           int     `json:"testSuites"`
	Source               string  `json:"source"`
}

type SecurityScanTypeWorkflowsReportResponse struct {
	ComponentId          string        `json:"componentId"`
	ComponentName        string        `json:"component"`
	WorkflowName         string        `json:"workflow"`
	AutomationId         string        `json:"automationId"`
	RunId                string        `json:"runId"`
	Branch               string        `json:"branch"`
	BranchId             string        `json:"branchId"`
	Build                float64       `json:"buildId"`
	Scanners             []interface{} `json:"scannerName"`
	ScannerType          []interface{} `json:"scanType"`
	RunStartTimeInMillis int64         `json:"runStartTimeInMillis"`
}

type SuccessfulBuildsResponse struct {
	ComponentId       string  `json:"componentId"`
	ComponentName     string  `json:"componentName"`
	WorkflowName      string  `json:"workflowName"`
	AutomationId      string  `json:"automationId"`
	BranchId          string  `json:"branchId"`
	Branch            string  `json:"branch"`
	RunId             string  `json:"runId"`
	JobId             string  `json:"jobId"`
	StepId            string  `json:"stepId"`
	RunNumber         float64 `json:"runNumber"`
	StartTime         string  `json:"startTime"`
	StartTimeInMillis int64   `json:"startTimeInMillis"`
	EndTime           string  `json:"endTime"`
	EndTimeInMillis   int64   `json:"endTimeInMillis"`
	Duration          float64 `json:"duration"`
	Environment       string  `json:"environment"`
	Source            string  `json:"source"`
}

type DoraMetricsDeploymentResponse struct {
	ComponentId          string  `json:"componentId"`
	ComponentName        string  `json:"componentName"`
	WorkflowName         string  `json:"workflowName"`
	AutomationId         string  `json:"automationId"`
	BranchId             string  `json:"branchId"`
	Branch               string  `json:"branch"`
	RunId                string  `json:"runId"`
	RunNumber            float64 `json:"runNumber"`
	RunStartTime         string  `json:"runStartTime"`
	RunStartTimeInMillis int64   `json:"runStartTimeInMillis"`
	LeadTimeInMillis     int64   `json:"leadTimeInMillis"`
	DeployTime           string  `json:"deployTime"`
	DeployTimeInMillis   int64   `json:"deployTimeInMillis"`
}

type FailureRateResponse struct {
	ComponentId   string  `json:"componentId"`
	ComponentName string  `json:"componentName"`
	Deployments   float64 `json:"deployments"`
	Success       float64 `json:"success"`
	Failure       float64 `json:"failure"`
	FailureRate   float64 `json:"failureRate"`
}

type DoraMttrResponse struct {
	ComponentId              string  `json:"componentId"`
	ComponentName            string  `json:"componentName"`
	FailedRunNumber          int64   `json:"failedRunNumber"`
	RecoveredRunNumber       int64   `json:"recoveredRunNumber"`
	FailedTime               string  `json:"failedTime"`
	FailedTimeInMillis       float64 `json:"failedTimeInMillis"`
	RecoveredTime            string  `json:"recoveredTime"`
	RecoveredTimeInMillis    float64 `json:"recoveredTimeInMillis"`
	RecoveryDurationInMillis float64 `json:"recoveryDurationInMillis"`
}

type RunDetailsTestResultsIndicatorsResponse struct {
	TestCasesFailed         int  `json:"testCasesFailed"`
	TestCasesPassed         int  `json:"testCasesPassed"`
	TestCasesSkipped        int  `json:"testCasesSkipped"`
	IsTestInsightsDataFound bool `json:"isTestInsightsDataFound"`
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
}

type Plugin struct {
	LongName            string              `json:"longName"`
	ShortName           string              `json:"shortName"`
	Version             string              `json:"version"`
	Enabled             string              `json:"enabled"`
	Active              bool                `json:"active"`
	HasUpdate           string              `json:"hasUpdate"`
	RequiredCoreVersion string              `json:"requiredCoreVersion"`
	MinimumJavaVersion  string              `json:"minimumJavaVersion"`
	Dependencies        *structpb.ListValue `json:"dependencies"`
	Status              string              `json:"status"`
}

type Dependency struct {
	ShortName string `json:"shortName"`
	Version   string `json:"version"`
	Optional  bool   `json:"optional"`
}

type CompletedRunData struct {
	EndpointId string  `json:"endpointId"`
	JobId      string  `json:"jobId"`
	RunId      float64 `json:"runId"`
	RunTime    float64 `json:"runTime"`
	StartTime  float64 `json:"startTime"`
	Result     string  `json:"result"`
	Type       string  `json:"type"`
	JenkinsUrl string  `json:"jenkinsUrl"`
}

type CompletedRunHeader struct {
	ProjectType   string  `json:"projectType"`
	SuccessRuns   float64 `json:"successRuns"`
	FailedRuns    float64 `json:"failedRuns"`
	AbortedRuns   float64 `json:"abortedRuns"`
	UnstableRuns  float64 `json:"unstableRuns"`
	TotalExecuted float64 `json:"totalExecuted"`
	TotalRunTime  float64 `json:"totalRunTime"`
	AvgRunTime    float64 `json:"avgRunTime"`
}

// Defining map to fetch actual names of scanners
var scannerNameMap = map[string]string{
	"anchore":                "Anchore",
	"aquasec":                "Aquasec",
	"gosec":                  "Gosec",
	"snyksca":                "Snyk SCA",
	"snyksast":               "Snyk SAST",
	"mendsca":                "Mend SCA",
	"mendsast":               "Mend SAST",
	"checkmarx":              "Checkmarx",
	"sonarqube":              "SonarQube",
	"trivy":                  "Trivy",
	"findsecbugs":            "Find Security Bugs",
	"githubsecurity":         "GitHub Security Scanner",
	"trufflehogs3":           "TruffleHog S3",
	"trufflehogcontainer":    "TruffleHog Container",
	"snykcontainer":          "Snyk Container",
	"jfrog-xray":             "JFrog Xray",
	"stackhawk":              "Stackhawk",
	"zap":                    "ZAP",
	"nexusiq":                "Nexus IQ",
	"sonarqube-bundled":      "SonarQube bundled",
	"trufflehogsast":         "TruffleHog SAST",
	"nexusiq-scan-container": "Sonatype (Nexus) Container",
}

var (
	openSearchClient        = opensearchconfig.GetOpensearchConnection
	getSearchResponse       = db.GetOpensearchData
	getOrganisationServices = helper.GetOrganisationServices
	multiSearchResponse     = helper.GetMultiQueryResponse
	getDocCount             = helper.GetIndexDocCount
)

func convertTimeFormat(inputTime string, timeFormat string) (string, error) {
	parsedTime, err := time.Parse("2006/01/02 15:04:05", inputTime)
	if err != nil {
		return "", fmt.Errorf("invalid time format: %w", err)
	}
	if timeFormat == "12h" {
		return parsedTime.Format("2006/01/02 03:04:05 PM"), nil
	} else if timeFormat == "24h" {
		return parsedTime.Format("2006/01/02 15:04:05"), nil
	}

	return "", fmt.Errorf("invalid time format specified")
}

func ComponentReport(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	log.Debugf("Time took to create opensearch client for component drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.ComponentDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	log.Debugf("Time took to replace placeholders for component drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	response, err := getSearchResponse(updatedJSON, constants.AUTOMATION_METADATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all component and last active time for component drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	componentActivityMap := make(map[string]ComponentReportResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.COMPONENT_ACTIVITY] != nil {
			automationRuns := aggsResult[constants.COMPONENT_ACTIVITY].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				for key, value := range values {
					valueMap := value.(map[string]interface{})
					componentResponse := ComponentReportResponse{
						ComponentId:   valueMap["component_id"].(string),
						ComponentName: valueMap["component_name"].(string),
						RepositoryUrl: valueMap["repo_url"].(string),
						LastActive:    valueMap["last_active_time"].(string),
					}
					componentActivityMap[key] = componentResponse
				}
			}
		}
	}
	log.Debugf("Time took to process all component for component drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}
	reports := &structpb.ListValue{}
	responseList := []ComponentReportResponse{}
	startTime = time.Now()
	serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
	if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
		return nil, err
	}
	log.Debugf("Time took to fetch all services : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	start, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN, replacements["startDate"].(string))
	if err != nil {
		log.Error("Exception while parsing start date", err)
		return nil, db.ErrParsingDate
	}
	endDateString := replacements["endDate"].(string)
	dateString := strings.Split(endDateString, " ")[0]
	end, err := time.Parse(timeLayoutDateHistogram, dateString)
	if err != nil {
		log.Error("Exception while parsing end date", err)
		return nil, db.ErrParsingDate
	}
	startDate := getStartOfTheDay(start)
	endDate := getEndOfTheDay(end)
	if serviceResponse != nil {
		for i := 0; i < len(serviceResponse.GetService()); i++ {
			service := serviceResponse.GetService()[i]
			found := false
			if len(components) > 0 {
				for _, component := range components {
					if service.Id == component {
						found = true
						break
					}
				}
			}
			if found || components == nil {
				componentResponse, ok := componentActivityMap[service.Id]
				var lastActiveInMillis int64
				lastActiveInMillis = 0
				status := constants.INACTIVE
				if ok {
					lastActive := componentResponse.LastActive
					if lastActive != "" {
						lastActiveTime, _ := time.Parse(constants.DATE_LAYOUT, lastActive)
						if (lastActiveTime.After(startDate) && lastActiveTime.Before(endDate)) || lastActiveTime.Equal(startDate) || lastActiveTime.Equal(endDate) {
							status = constants.ACTIVE
						}
						lastActiveInMillis = lastActiveTime.UnixMilli()
						lastActive = lastActiveTime.Format(constants.DATE_FORMAT)
					} else {
						lastActive = "-"
					}
					componentResponse.LastActive = lastActive
					componentResponse.LastActiveInMillis = lastActiveInMillis
					componentResponse.Status = status
				} else {
					componentResponse = ComponentReportResponse{
						ComponentName:      service.Name,
						RepositoryUrl:      service.RepositoryUrl,
						Status:             status,
						LastActive:         "-",
						LastActiveInMillis: lastActiveInMillis,
					}
				}
				responseList = append(responseList, componentResponse)
			}
		}
	}
	log.Debugf("Time took to filter and get component response for component drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].LastActiveInMillis > responseList[j].LastActiveInMillis
	})
	log.Debugf("Time took to sort the result for component drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, response := range responseList {
		lastActiveTime, err := helper.ConvertUTCtoTimeZone(response.LastActive, replacements["timeZone"].(string))
		convertedTime, _ := convertTimeFormat(lastActiveTime, replacements["timeFormat"].(string))
		lastActiveTime = convertedTime
		log.CheckErrorf(err, "Error converting last active time - %s to timezone - %s", response.LastActive, replacements["timeZone"].(string))
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT_ID:   structpb.NewStringValue(response.ComponentId),
						constants.COMPONENT_NAME: structpb.NewStringValue(response.ComponentName),
						constants.REPOSITORY_URL: structpb.NewStringValue(response.RepositoryUrl),
						constants.LAST_ACTIVE:    structpb.NewStringValue(lastActiveTime),
						constants.STATUS:         structpb.NewStringValue(response.Status),
					},
				},
			},
		})
	}
	return reports, nil
}

func AutomationReport(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.AutomationDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	response, err := getSearchResponse(updatedJSON, constants.AUTOMATION_METADATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	automationActivityMap := make(map[string]WorkflowReportResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.AUTOMATION_ACTIVITY] != nil {
			automationRuns := aggsResult[constants.AUTOMATION_ACTIVITY].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				for key, value := range values {
					valueMap := value.(map[string]interface{})
					automationResponse := WorkflowReportResponse{
						ComponentId:   valueMap["component_id"].(string),
						ComponentName: valueMap["component_name"].(string),
						WorkflowName:  valueMap["workflow_name"].(string),
						AutomationId:  valueMap["automation_id"].(string),
						BranchId:      valueMap["branch_id"].(string),
						Branch:        valueMap["branch_name"].(string),
						LastActive:    valueMap["last_active_time"].(string),
					}

					// resourceName, source, err := cutils.GetDisplayNameAndOrigin(valueMap["workflow_name"].(string))
					// if err != nil {
					// 	log.Error("Error getting display name and origin for AutomationReport()", err)
					// } else {
					// 	automationResponse.Source = source
					// 	automationResponse.WorkflowName = resourceName
					// }
					automationActivityMap[key] = automationResponse
				}
			}
		}
	}
	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationMilliSec, time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}
	reports := &structpb.ListValue{}
	responseList := []WorkflowReportResponse{}
	startTime = time.Now()
	automationResponse := getEnabledAutomationResponses(ctx, clt, orgId, components, nil, true)
	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationInfoFromCacheMilliSec, time.Since(startTime).Milliseconds())

	start, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN, replacements["startDate"].(string))
	if err != nil {
		log.Error("Exception while parsing start date in Automation report", err)
	}
	endDateString := replacements["endDate"].(string)
	dateString := strings.Split(endDateString, " ")[0]
	end, err := time.Parse(timeLayoutDateHistogram, dateString)
	if err != nil {
		log.Error("Exception while parsing start date in Automation report", err)
	}
	startDate := getStartOfTheDay(start)
	endDate := getEndOfTheDay(end)

	if len(automationResponse) > 0 {
		for id, automation := range automationResponse {
			found := false
			if len(components) > 0 {
				for _, component := range components {
					if automation.ComponentId == component {
						found = true
						break
					}
				}
			}
			if found || components == nil {
				workflowReponse, ok := automationActivityMap[id]
				var lastActiveInMillis int64
				lastActiveInMillis = 0
				status := constants.INACTIVE
				if ok {
					lastActive := workflowReponse.LastActive
					if lastActive != "" {
						lastActiveTime, _ := time.Parse(constants.DATE_LAYOUT, lastActive)

						if (lastActiveTime.After(startDate) && lastActiveTime.Before(endDate)) || lastActiveTime.Equal(startDate) || lastActiveTime.Equal(endDate) {
							status = constants.ACTIVE
						}
						lastActiveInMillis = lastActiveTime.UnixMilli()
						lastActive = lastActiveTime.Format(constants.DATE_FORMAT)
					} else {
						lastActive = "-"
					}
					workflowReponse.LastActive = lastActive
					workflowReponse.LastActiveInMillis = lastActiveInMillis
					workflowReponse.Status = status
					workflowReponse.WorkflowName = automation.WorkflowName
					workflowReponse.Source = automation.Source
				} else {
					workflowReponse = WorkflowReportResponse{
						ComponentId:        automation.ComponentId,
						ComponentName:      automation.ComponentName,
						WorkflowName:       automation.WorkflowName,
						AutomationId:       automation.AutomationId,
						BranchId:           automation.BranchId,
						Branch:             automation.Branch,
						Status:             status,
						LastActive:         "-",
						LastActiveInMillis: lastActiveInMillis,
						Source:             automation.Source,
					}
				}
				responseList = append(responseList, workflowReponse)
			}
		}
	}
	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationFromCacheAndOpensearchMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].LastActiveInMillis > responseList[j].LastActiveInMillis
	})
	log.Debugf(exceptions.DebugTimeTookToSortAllResultForAutomationMilliSec, time.Since(startTime).Milliseconds())
	for _, response := range responseList {
		lastActiveTime, err := helper.ConvertUTCtoTimeZone(response.LastActive, replacements["timeZone"].(string))
		convertedTime, _ := convertTimeFormat(lastActiveTime, replacements["timeFormat"].(string))
		lastActiveTime = convertedTime
		log.CheckErrorf(err, "Error converting last active time - %s to timezone - %s", response.LastActive, replacements["timeZone"].(string))

		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:       structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:    structpb.NewStringValue(response.ComponentId),
						constants.AUTOMATION_ID:   structpb.NewStringValue(response.AutomationId),
						constants.WORKFLOW:        structpb.NewStringValue(response.WorkflowName),
						constants.BRANCH:          structpb.NewStringValue(response.Branch),
						constants.BRANCH_ID:       structpb.NewStringValue(response.BranchId),
						constants.LAST_ACTIVE:     structpb.NewStringValue(lastActiveTime),
						constants.STATUS:          structpb.NewStringValue(response.Status),
						constants.WORKFLOW_SOURCE: structpb.NewStringValue(response.Source),
					},
				},
			},
		})
	}
	return reports, nil
}
func AutomationReportForBranch(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client, err := opensearchconfig.GetOpensearchConnection()
	branch, ok := replacements[constants.BRANCH]
	reports := &structpb.ListValue{}
	if ok {
		log.CheckErrorf(err, "Error establishing connection with OpenSearch in getQueryResponse(). Connection error - ")
		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.AutomationDrilldownQueryForBranch)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		response, err := db.GetOpensearchData(updatedJSON, constants.AUTOMATION_METADATA_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationMilliSec, time.Since(startTime).Milliseconds())
		startTime = time.Now()
		result := make(map[string]interface{})
		json.Unmarshal([]byte(response), &result)
		automationActivityMap := make(map[string]WorkflowReportResponse)
		if result[constants.AGGREGATION] != nil {
			aggsResult := result[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.AUTOMATION_ACTIVITY] != nil {
				automationRuns := aggsResult[constants.AUTOMATION_ACTIVITY].(map[string]interface{})
				if automationRuns[constants.VALUE] != nil {
					values := automationRuns[constants.VALUE].(map[string]interface{})
					for key, value := range values {
						valueMap := value.(map[string]interface{})
						automationResponse := WorkflowReportResponse{
							ComponentId:   valueMap["component_id"].(string),
							ComponentName: valueMap["component_name"].(string),
							WorkflowName:  valueMap["workflow_name"].(string),
							AutomationId:  valueMap["automation_id"].(string),
							BranchId:      valueMap["branch_id"].(string),
							Branch:        valueMap["branch_name"].(string),
							LastActive:    valueMap["last_active_time"].(string),
						}
						automationActivityMap[key] = automationResponse
					}
				}
			}
		}
		log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationMilliSec, time.Since(startTime).Milliseconds())
		orgId, ok := replacements[constants.ORG_ID].(string)
		if !ok {
			orgId = ""
		}
		subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
		if ok && orgId != subOrgId && subOrgId != "" {
			orgId = subOrgId
		}
		components, ok := replacements[constants.COMPONENT].([]string)
		if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
			components = nil
		}
		responseList := []WorkflowReportResponse{}
		startTime = time.Now()
		automationResponse := getAutomationResponseMapForBranch(ctx, clt, orgId, components, nil, true, branch.(string))
		log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationInfoFromCacheMilliSec, time.Since(startTime).Milliseconds())

		start, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN, replacements["startDate"].(string))
		if err != nil {
			log.Error("Exception while parsing start date", err)
			return nil, db.ErrParsingDate
		}
		endDateString := replacements["endDate"].(string)
		dateString := strings.Split(endDateString, " ")[0]
		end, err := time.Parse(timeLayoutDateHistogram, dateString)
		if err != nil {
			log.Error("Exception while parsing end date", err)
			return nil, db.ErrParsingDate
		}
		startDate := getStartOfTheDay(start)
		endDate := getEndOfTheDay(end)

		if len(automationResponse) > 0 {
			for id, automation := range automationResponse {
				found := false
				if len(components) > 0 {
					for _, component := range components {
						if automation.ComponentId == component {
							found = true
							break
						}
					}
				}
				if found || components == nil {
					workflowReponse, ok := automationActivityMap[id]
					var lastActiveInMillis int64
					lastActiveInMillis = 0
					status := constants.INACTIVE
					if ok {
						lastActive := workflowReponse.LastActive
						if lastActive != "" {
							lastActiveTime, _ := time.Parse(constants.DATE_LAYOUT, lastActive)

							if (lastActiveTime.After(startDate) && lastActiveTime.Before(endDate)) || lastActiveTime.Equal(startDate) || lastActiveTime.Equal(endDate) {
								status = constants.ACTIVE
							}
							lastActiveInMillis = lastActiveTime.UnixMilli()
							lastActive = lastActiveTime.Format(constants.DATE_FORMAT)
						} else {
							lastActive = "-"
						}
						workflowReponse.LastActive = lastActive
						workflowReponse.LastActiveInMillis = lastActiveInMillis
						workflowReponse.Status = status
						workflowReponse.WorkflowName = automation.WorkflowName
						workflowReponse.Source = automation.Source
					} else {
						workflowReponse = WorkflowReportResponse{

							WorkflowName:       automation.WorkflowName,
							AutomationId:       automation.AutomationId,
							BranchId:           automation.BranchId,
							Branch:             automation.Branch,
							Status:             status,
							LastActive:         "-",
							LastActiveInMillis: lastActiveInMillis,
							Source:             automation.Source,
						}
					}
					responseList = append(responseList, workflowReponse)
				}
			}
		}
		log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationFromCacheAndOpensearchMilliSec, time.Since(startTime).Milliseconds())
		startTime = time.Now()
		sort.Slice(responseList, func(i, j int) bool {
			return responseList[i].LastActiveInMillis > responseList[j].LastActiveInMillis
		})
		log.Debugf(exceptions.DebugTimeTookToSortAllResultForAutomationMilliSec, time.Since(startTime).Milliseconds())
		for _, response := range responseList {
			reports.Values = append(reports.Values, &structpb.Value{
				Kind: &structpb.Value_StructValue{
					StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							constants.COMPONENT:       structpb.NewStringValue(response.ComponentName),
							constants.COMPONENT_ID:    structpb.NewStringValue(response.ComponentId),
							constants.AUTOMATION_ID:   structpb.NewStringValue(response.AutomationId),
							constants.WORKFLOW:        structpb.NewStringValue(response.WorkflowName),
							constants.BRANCH:          structpb.NewStringValue(response.Branch),
							constants.BRANCH_ID:       structpb.NewStringValue(response.BranchId),
							constants.LAST_ACTIVE:     structpb.NewStringValue(response.LastActive),
							constants.STATUS:          structpb.NewStringValue(response.Status),
							constants.WORKFLOW_SOURCE: structpb.NewStringValue(response.Source),
						},
					},
				},
			})
		}
	}
	return reports, nil
}
func getAutomationResponseMap(ctx context.Context, clt client.GrpcClient, orgId string, components []string, automationSet map[string]struct{}, excludeDisabledBranch bool) map[string]WorkflowReportResponse {
	automationMap := map[string]WorkflowReportResponse{}
	coreDataCache := cache.GetCoreDataCache()
	if coreDataCache != nil {
		startTime := time.Now()
		serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
		if err != nil {
			return automationMap
		}
		log.Debugf("Time took to fectch all services : %v in milliseconds", time.Since(startTime).Microseconds())
		startTime = time.Now()
		if serviceResponse != nil {
			for i := 0; i < len(serviceResponse.GetService()); i++ {
				service := serviceResponse.GetService()[i]
				found := false
				if len(components) > 0 {
					for _, component := range components {
						if service.Id == component {
							found = true
							break
						}
					}
				}
				if found || components == nil {
					automation := WorkflowReportResponse{}
					for _, child := range coreDataCache.GetChildrenOfType(service.Id, api.ResourceType_RESOURCE_TYPE_BRANCH) {
						childResource := coreDataCache.Get(child)
						if !excludeDisabledBranch || !childResource.IsDisabled {
							automations := coreDataCache.GetChildrenOfType(child, api.ResourceType_RESOURCE_TYPE_AUTOMATION)
							if len(automations) > 0 {
								for _, id := range automations {
									found := false
									if automationSet != nil {
										_, ok := automationSet[id]
										if ok {
											found = true
										}
									}
									if found || automationSet == nil {
										automationResource := coreDataCache.Get(id)
										automation.ComponentId = service.Id
										automation.ComponentName = service.Name
										automation.AutomationId = id
										workflowName, sourceName, err := cutils.GetDisplayNameAndOrigin(automationResource.Name)
										if log.CheckErrorf(err, "Error getting display name and origin for getAutomationResponseMap()", err) {
											log.Infof("Error getting display name and origin for getAutomationResponseMap()")
										}
										automation.WorkflowName = workflowName
										automation.Source = sourceName
										automation.BranchId = childResource.Id
										automation.Branch = childResource.Name
										automationMap[id] = automation
									}
								}
							}
						}
					}
				}
			}
		}
		log.Debugf("Time took to fectch all automations : %v in milliseconds", time.Since(startTime).Microseconds())
	} else {
		log.Info("Core data cache is null.")
	}
	return automationMap
}
func getEnabledAutomationResponses(ctx context.Context, clt client.GrpcClient, orgId string, components []string, automationSet map[string]struct{}, excludeDisabledBranch bool) map[string]WorkflowReportResponse {
	automationMap := map[string]WorkflowReportResponse{}
	coreDataCache := cache.GetCoreDataCache()
	if coreDataCache != nil {
		startTime := time.Now()
		serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
		if err != nil {
			return automationMap
		}
		log.Debugf("Time took to fectch all services : %v in milliseconds", time.Since(startTime).Microseconds())
		startTime = time.Now()
		if serviceResponse != nil {
			for i := 0; i < len(serviceResponse.GetService()); i++ {
				service := serviceResponse.GetService()[i]
				found := false
				if len(components) > 0 {
					for _, component := range components {
						if service.Id == component {
							found = true
							break
						}
					}
				}
				if found || components == nil {
					automation := WorkflowReportResponse{}
					for _, child := range coreDataCache.GetChildrenOfType(service.Id, api.ResourceType_RESOURCE_TYPE_BRANCH) {
						childResource := coreDataCache.Get(child)
						if !excludeDisabledBranch || !childResource.IsDisabled {
							automations := coreDataCache.GetChildrenOfType(child, api.ResourceType_RESOURCE_TYPE_AUTOMATION)
							if len(automations) > 0 {
								for _, id := range automations {
									found := false
									if automationSet != nil {
										_, ok := automationSet[id]
										if ok {
											found = true
										}
									}
									if found || automationSet == nil {
										automationResource := coreDataCache.Get(id)
										automation.ComponentId = service.Id
										automation.ComponentName = service.Name
										automation.AutomationId = id
										workflowName, sourceName, err := cutils.GetDisplayNameAndOrigin(automationResource.Name)
										if log.CheckErrorf(err, "Error getting display name and origin for getEnabledAutomationResponses()", err) {
											log.Infof("Error getting display name and origin for getEnabledAutomationResponses()")
										}
										automation.WorkflowName = workflowName
										automation.Source = sourceName
										automation.BranchId = childResource.Id
										automation.Branch = childResource.Name
										if !automationResource.IsDisabled {
											automationMap[id] = automation
										}
									}
								}
							}
						}
					}
				}
			}
		}
		log.Debugf("Time took to fectch all automations : %v in milliseconds", time.Since(startTime).Microseconds())
	} else {
		log.Info("Core data cache is null.")
	}
	return automationMap
}

func getAutomationResponseMapForBranch(ctx context.Context, clt client.GrpcClient, orgId string, components []string, automationSet map[string]struct{}, excludeDisabledBranch bool, branchId string) map[string]WorkflowReportResponse {
	automationMap := map[string]WorkflowReportResponse{}
	coreDataCache := cache.GetCoreDataCache()
	if coreDataCache != nil {
		startTime := time.Now()
		serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
		if err != nil {
			return automationMap
		}
		log.Debugf("Time took to fectch all services : %v in milliseconds", time.Since(startTime).Microseconds())
		startTime = time.Now()
		if serviceResponse != nil {
			for i := 0; i < len(serviceResponse.GetService()); i++ {
				service := serviceResponse.GetService()[i]
				found := false
				if len(components) > 0 {
					for _, component := range components {
						if service.Id == component {
							found = true
							break
						}
					}
				}
				if found || components == nil {
					automation := WorkflowReportResponse{}
					for _, child := range coreDataCache.GetChildrenOfType(service.Id, api.ResourceType_RESOURCE_TYPE_BRANCH) {
						childResource := coreDataCache.Get(child)
						if !excludeDisabledBranch || !childResource.IsDisabled && childResource.Id == branchId {
							automations := coreDataCache.GetChildrenOfType(child, api.ResourceType_RESOURCE_TYPE_AUTOMATION)
							if len(automations) > 0 {
								for _, id := range automations {
									found := false
									if automationSet != nil {
										_, ok := automationSet[id]
										if ok {
											found = true
										}
									}
									if found || automationSet == nil {
										automationResource := coreDataCache.Get(id)
										automation.ComponentId = service.Id
										automation.ComponentName = service.Name
										automation.AutomationId = id
										workflowName, sourceName, err := cutils.GetDisplayNameAndOrigin(automationResource.Name)
										if log.CheckErrorf(err, "Error getting display name and origin for getAutomationResponseMapForBranch()", err) {
											log.Infof("Error getting display name and origin for getAutomationResponseMapForBranch()")
										}
										automation.WorkflowName = workflowName
										automation.Source = sourceName
										automation.BranchId = childResource.Id
										automation.Branch = childResource.Name
										automationMap[id] = automation
									}
								}
							}
						}
					}
				}
			}
		}
		log.Debugf("Time took to fectch all automations : %v in milliseconds", time.Since(startTime).Microseconds())
	} else {
		log.Info("Core data cache is null.")
	}
	return automationMap
}

func AutomationRunReport(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.AutomationRunDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	branch, ok := replacements[constants.BRANCH]
	var modifiedJson string
	if ok && branch != nil && len(branch.(string)) > 0 {
		modifiedJson = internal.UpdateFiltersForDrilldown(updatedJSON, replacements, true, false)
	} else {
		modifiedJson = internal.UpdateFilters(updatedJSON, replacements)
	}
	response, err := getSearchResponse(modifiedJson, constants.AUTOMATION_RUN_STATUS_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationRunMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	automationSet := make(map[string]struct{})
	automationRunMap := make(map[string][]WorkflowRunReportResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.AUTOMATION_RUN_ACTIVITY] != nil {
			automationRuns := aggsResult[constants.AUTOMATION_RUN_ACTIVITY].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				for key, value := range values {
					automationSet[key] = struct{}{}
					runs := value.([]interface{})
					workflowRuns := []WorkflowRunReportResponse{}
					for _, runValueMap := range runs {
						runValue := runValueMap.(map[string]interface{})
						statusTime, err := time.Parse(constants.DATE_LAYOUT, runValue["status_timestamp"].(string))
						runInfo := WorkflowRunReportResponse{
							ComponentId:   runValue["component_id"].(string),
							ComponentName: runValue["component_name"].(string),
							Status:        runValue["status"].(string),
							AutomationId:  runValue["automation_id"].(string),
							RunId:         runValue["run_id"].(string),
							Build:         runValue["run_number"].(float64),
						}
						if err == nil {
							runInfo.RunStartTimeInMillis = int64(float64(statusTime.UnixMilli()) - runValue["duration"].(float64))
							runStartDate := time.UnixMilli(runInfo.RunStartTimeInMillis).Format(constants.DATE_FORMAT)
							runInfo.RunStartTime = runStartDate
						}
						workflowRuns = append(workflowRuns, runInfo)
					}
					automationRunMap[key] = workflowRuns
				}
			}
		}
	}
	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationRunMilliSec, time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}
	reports := &structpb.ListValue{}
	responseList := []WorkflowRunReportResponse{}
	startTime = time.Now()

	var automationResponse map[string]WorkflowReportResponse
	if branch != nil && len(branch.(string)) > 0 {
		automationResponse = getAutomationResponseMapForBranch(ctx, clt, orgId, components, automationSet, false, branch.(string))

	} else {
		automationResponse = getAutomationResponseMap(ctx, clt, orgId, components, automationSet, false)

	}
	autorunIDs := make(map[string]bool)

	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationRunInfoFromCacheMilliSec, time.Since(startTime).Milliseconds())
	for id, runs := range automationRunMap {
		automation, ok := automationResponse[id]
		for _, workflowRun := range runs {
			if ok {
				workflowRun.Branch = automation.Branch
				workflowRun.BranchId = automation.BranchId
				workflowRun.WorkflowName = automation.WorkflowName
				workflowRun.Source = automation.Source
			} else {
				log.Infof("Automation data not found in coredatacache for automationid : %s", id)
			}
			if !autorunIDs[workflowRun.RunId] {
				autorunIDs[workflowRun.RunId] = true
				responseList = append(responseList, workflowRun)
			}
		}
	}
	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationRunFromCacheAndOpensearchMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].RunStartTimeInMillis > responseList[j].RunStartTimeInMillis
	})
	log.Debugf(exceptions.DebugTimeTookToSortAllResultForAutomationRunMilliSec, time.Since(startTime).Milliseconds())
	for _, response := range responseList {
		runStartTime, err := helper.ConvertUTCtoTimeZone(response.RunStartTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, "Error converting run start time - %s to timezone - %s", response.RunStartTime, replacements["timeZone"].(string))
		convertedTime, _ := convertTimeFormat(runStartTime, replacements["timeFormat"].(string))
		runStartTime = convertedTime
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:       structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:    structpb.NewStringValue(response.ComponentId),
						constants.WORKFLOW:        structpb.NewStringValue(response.WorkflowName),
						constants.AUTOMATION_ID:   structpb.NewStringValue(response.AutomationId),
						constants.BRANCH:          structpb.NewStringValue(response.Branch),
						constants.BRANCH_ID:       structpb.NewStringValue(response.BranchId),
						constants.BUILD:           structpb.NewNumberValue(response.Build),
						constants.RUN_ID_KEY:      structpb.NewStringValue(response.RunId),
						constants.RUN_START_TIME:  structpb.NewStringValue(runStartTime),
						constants.STATUS:          structpb.NewStringValue(response.Status),
						constants.WORKFLOW_SOURCE: structpb.NewStringValue(response.Source),
					},
				},
			},
		})
	}
	return reports, nil
}

func PullRequestsReport(replacements map[string]any, ctx context.Context, grpcClient client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.PullRequestDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)
	response, err := getSearchResponse(modifiedJson, constants.PULL_REQUESTS_REVIEW_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all commit data for pullrequests drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	pullrequestsResponse := []PullRequestResponse{}
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.PULL_REQUESTS] != nil {
			pullrequests := aggsResult[constants.PULL_REQUESTS].(map[string]interface{})
			if pullrequests[constants.VALUE] != nil {
				values := pullrequests[constants.VALUE].([]interface{})
				for _, value := range values {
					pullRequestsMap := value.(map[string]interface{})
					createdTimestamp, err := time.Parse(constants.DATE_LAYOUT_TZ, pullRequestsMap["pr_created_time"].(string))
					prInfo := PullRequestResponse{
						ComponentId:    pullRequestsMap["component_id"].(string),
						ComponentName:  pullRequestsMap["component_name"].(string),
						Status:         pullRequestsMap["review_status"].(string),
						RepositoryName: pullRequestsMap["repository_url"].(string),
						PrId:           pullRequestsMap["pull_request_id"].(string),
						Provider:       pullRequestsMap["provider"].(string),
					}
					if err == nil {
						prInfo.CreatedTimeInMillis = int64(float64(createdTimestamp.UnixMilli()))
						createdTime := createdTimestamp.Format(constants.DATE_FORMAT_TZ)
						convertedTime, _ := convertTimeFormat(createdTime, replacements["timeFormat"].(string))
						createdTime = convertedTime
						prInfo.CreatedTime = createdTime
					}
					sourceBranch, ok := pullRequestsMap["source_branch"]
					if ok {
						prInfo.SourceBranch = sourceBranch.(string)
					}
					targetBranch, ok := pullRequestsMap["target_branch"]
					if ok {
						prInfo.TargetBranch = targetBranch.(string)
					}
					pullrequestsResponse = append(pullrequestsResponse, prInfo)
				}
			}
		}
	}
	log.Debugf("Time took to process all data for pullrequests drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	reports := &structpb.ListValue{}
	startTime = time.Now()
	sort.Slice(pullrequestsResponse, func(i, j int) bool {
		return pullrequestsResponse[i].CreatedTimeInMillis > pullrequestsResponse[j].CreatedTimeInMillis
	})
	log.Debugf("Time took to sort the result for pullrequests drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, response := range pullrequestsResponse {
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:     structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:  structpb.NewStringValue(response.ComponentId),
						constants.PR_ID:         structpb.NewStringValue(response.PrId),
						constants.SOURCE_BRANCH: structpb.NewStringValue(response.SourceBranch),
						constants.REPO:          structpb.NewStringValue(response.RepositoryName),
						constants.TARGET_BRANCH: structpb.NewStringValue(response.TargetBranch),
						constants.STATUS:        structpb.NewStringValue(response.Status),
						constants.PROVIDER:      structpb.NewStringValue(response.Provider),
						constants.CREATED_ON:    structpb.NewStringValue(response.CreatedTime),
					},
				},
			},
		})
	}
	return reports, nil
}

func CommitsReport(replacements map[string]any, ctx context.Context, grpcClient client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CommitsDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	branch, ok := replacements[constants.BRANCH]
	var modifiedJson string

	if ok && branch != nil && len(branch.(string)) > 0 {
		modifiedJson = internal.UpdateFiltersForDrilldown(updatedJSON, replacements, false, true)

	} else {
		modifiedJson = internal.UpdateFilters(updatedJSON, replacements)

	}
	response, err := getSearchResponse(modifiedJson, constants.COMMIT_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all commit data for commits drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	commits := []CommitsResponse{}
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.COMMITS] != nil {
			pullrequests := aggsResult[constants.COMMITS].(map[string]interface{})
			if pullrequests[constants.VALUE] != nil {
				values := pullrequests[constants.VALUE].([]interface{})
				for _, value := range values {
					commitMap := value.(map[string]interface{})
					commitTimestamp, err := time.Parse(constants.DATE_LAYOUT_TZ, commitMap["commit_timestamp"].(string))
					commitInfo := CommitsResponse{
						ComponentId:    commitMap["component_id"].(string),
						ComponentName:  commitMap["component_name"].(string),
						Author:         commitMap["author"].(string),
						Branch:         commitMap["branch"].(string),
						RepositoryName: commitMap["repository_url"].(string),
						CommitId:       commitMap["commit_id"].(string),
					}
					if err == nil {
						commitInfo.CommitTimestampInMillis = int64(float64(commitTimestamp.UnixMilli()))
						commitDate := commitTimestamp.Format(constants.DATE_FORMAT_TZ)
						convertedTimeCommitDate, _ := convertTimeFormat(commitDate, replacements["timeFormat"].(string))
						commitInfo.CommitTimestamp = convertedTimeCommitDate
					}
					commits = append(commits, commitInfo)
				}
			}
		}
	}
	log.Debugf("Time took to process all data for commits drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	reports := &structpb.ListValue{}
	startTime = time.Now()
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].CommitTimestampInMillis > commits[j].CommitTimestampInMillis
	})
	log.Debugf("Time took to sort the result for commits drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, response := range commits {
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:      structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_NAME: structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:   structpb.NewStringValue(response.ComponentId),
						constants.COMMIT_ID:      structpb.NewStringValue(response.CommitId),
						constants.BRANCH:         structpb.NewStringValue(response.Branch),
						constants.REPO:           structpb.NewStringValue(response.RepositoryName),
						constants.AUTHOR:         structpb.NewStringValue(response.Author),
						constants.COMMIT_TIME:    structpb.NewStringValue(response.CommitTimestamp),
					},
				},
			},
		})
	}
	return reports, nil
}

func CPSRunInitiatingCommits(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CPSRunInitiatingCommitsQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)
	response, err := getSearchResponse(modifiedJson, constants.AUTOMATION_RUN_STATUS_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all builds for run-initiating commits drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	componentSet := make(map[string]struct{})
	commitRunMap := getCommitRunMap(result, componentSet)
	log.Debugf("Time took to process all commit and runs for run-initiating commits drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	updatedJSON, err = db.ReplaceJSONplaceholders(replacements, constants.CPSRunCommitsDeployedEnvQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson = internal.UpdateFilters(updatedJSON, replacements)
	response, err = getSearchResponse(modifiedJson, constants.DEPLOY_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	runsEnvMap := getRunAndDeployedEnvMap(result, response)
	log.Debugf("Time took to process all run and deployed env for run-initiating commits drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	reports := &structpb.ListValue{}
	responseList := []CommitRunResponse{}
	startTime = time.Now()
	coreDataCache := cache.GetCoreDataCache()
	if coreDataCache != nil {
		for key := range componentSet {
			for _, child := range coreDataCache.GetChildrenOfType(key, api.ResourceType_RESOURCE_TYPE_BRANCH) {
				childResource := coreDataCache.Get(child)
				automations := coreDataCache.GetChildrenOfType(child, api.ResourceType_RESOURCE_TYPE_AUTOMATION)
				if len(automations) > 0 {
					for _, id := range automations {
						commitRuns, ok := commitRunMap[id]
						if ok {
							automationResource := coreDataCache.Get(id)
							for _, commitRun := range commitRuns {
								commitRun.WorkflowName = automationResource.Name
								commitRun.BranchId = childResource.Id
								commitRun.Branch = childResource.Name
							}
						}
					}
				}
			}
		}
	}
	log.Debugf("Time took to automation info from cache for run-initiating commits drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, commitRuns := range commitRunMap {
		for _, commitRun := range commitRuns {
			envs, ok := runsEnvMap[commitRun.RunId]
			if ok {
				commitRun.Environment = envs
			} else {
				commitRun.Environment = []string{"Unspecified"}
			}
			if commitRun.CommitId != "" {
				responseList = append(responseList, *commitRun)
			}
		}
	}
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].CreatedTimeInMillis > responseList[j].CreatedTimeInMillis
	})
	log.Debugf("Time took to sort the result for run-initiating commits drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, response := range responseList {
		createdTime, err := helper.ConvertUTCtoTimeZone(response.CreatedTime, replacements["timeZone"].(string))
		envList := &structpb.ListValue{}
		for _, env := range response.Environment {
			envList.Values = append(envList.Values, structpb.NewStringValue(env))
		}

		log.CheckErrorf(err, "Error converting create time - %s to timezone - %s", response.CreatedTime, replacements["timeZone"].(string))
		convertedTime, _ := convertTimeFormat(createdTime, replacements["timeFormat"].(string))
		createdTime = convertedTime
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:          structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:       structpb.NewStringValue(response.ComponentId),
						constants.COMMIT_ID:          structpb.NewStringValue(response.CommitId),
						constants.COMMIT_DESCRIPTION: structpb.NewStringValue(response.CommitDescription),
						constants.RUN_ID:             structpb.NewNumberValue(response.RunNumber),
						constants.RUN_ID_KEY:         structpb.NewStringValue(response.RunId),
						constants.AUTOMATION_ID:      structpb.NewStringValue(response.AutomationId),
						constants.BRANCH_ID:          structpb.NewStringValue(response.BranchId),
						constants.CREATED:            structpb.NewStringValue(createdTime),
						constants.STATUS:             structpb.NewStringValue(response.Status),
						constants.ENVIRONMENT:        structpb.NewListValue(envList),
					},
				},
			},
		})
	}
	return reports, nil
}

func getRunAndDeployedEnvMap(result map[string]interface{}, response string) map[string][]string {
	result = make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	runsEnvMap := make(map[string][]string)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.DEPLOYMENTS] != nil {
			runEnvs := aggsResult[constants.DEPLOYMENTS].(map[string]interface{})
			if runEnvs[constants.VALUE] != nil {
				values := runEnvs[constants.VALUE].(map[string]interface{})
				for key, value := range values {
					envs := value.([]interface{})
					envStrs := make([]string, len(envs))
					for i, v := range envs {
						envStrs[i] = fmt.Sprint(v)
					}
					runsEnvMap[key] = envStrs
				}
			}
		}
	}
	return runsEnvMap
}

func getCommitRunMap(result map[string]interface{}, componentSet map[string]struct{}) map[string][]*CommitRunResponse {
	commitRunMap := make(map[string][]*CommitRunResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.AUTOMATION_RUN] != nil {
			automationRun := aggsResult[constants.AUTOMATION_RUN].(map[string]interface{})
			if automationRun[constants.VALUE] != nil {
				values := automationRun[constants.VALUE].(map[string]interface{})
				for _, value := range values {
					runValue := value.(map[string]interface{})
					statusTime, err := time.Parse(constants.DATE_LAYOUT, runValue["status_timestamp"].(string))
					commitRunResponse := CommitRunResponse{
						ComponentId:         runValue["component_id"].(string),
						ComponentName:       runValue["component_name"].(string),
						AutomationId:        runValue["automation_id"].(string),
						Status:              runValue["status"].(string),
						RunNumber:           runValue["run_number"].(float64),
						RunId:               runValue["run_id"].(string),
						OrgName:             runValue["org_name"].(string),
						CreatedTimeInMillis: int64(runValue["start_time"].(float64)),
					}
					if err == nil {
						commitRunResponse.EndTimeInMillis = int64(float64(statusTime.UnixMilli()))
						commitRunResponse.EndTime = time.UnixMilli(commitRunResponse.EndTimeInMillis).Format(constants.DATE_FORMAT)
					}
					commitSha, ok := runValue["commit_sha"]
					if ok {
						commitRunResponse.CommitId = commitSha.(string)
					}
					commitDescription, ok := runValue["commit_description"]
					if ok {
						commitRunResponse.CommitDescription = commitDescription.(string)
					}
					componentSet[commitRunResponse.ComponentId] = struct{}{}
					if commitRunResponse.CreatedTimeInMillis == 0 {
						commitRunResponse.CreatedTimeInMillis = commitRunResponse.EndTimeInMillis
					}
					commitRunResponse.CreatedTime = time.UnixMilli(commitRunResponse.CreatedTimeInMillis).Format(constants.DATE_FORMAT)
					commits, automationPresent := commitRunMap[commitRunResponse.AutomationId]
					if automationPresent {
						commits = append(commits, &commitRunResponse)
						commitRunMap[commitRunResponse.AutomationId] = commits
					} else {
						commits = []*CommitRunResponse{}
						commits = append(commits, &commitRunResponse)
						commitRunMap[commitRunResponse.AutomationId] = commits
					}
				}
			}
		}
	}
	return commitRunMap
}

func CodeProgressionSnapshotBuilds(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CodeProgressionSnapshotBuild)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	var modifiedJson string
	branch, ok := replacements[constants.BRANCH]
	if ok && branch != nil && len(branch.(string)) > 0 {
		modifiedJson = internal.UpdateFiltersForDrilldown(updatedJSON, replacements, true, false)
	} else {
		modifiedJson = internal.UpdateFilters(updatedJSON, replacements)
	}
	response, err := getSearchResponse(modifiedJson, constants.BUILD_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all builds for build drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	automationSet := make(map[string]struct{})
	buildMap := make(map[string][]BuildsResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.BUILDS] != nil {
			builds := aggsResult[constants.BUILDS].(map[string]interface{})
			if builds[constants.VALUE] != nil {
				values := builds[constants.VALUE].(map[string]interface{})
				for _, value := range values {
					runValue := value.(map[string]interface{})
					statusTime, err := time.Parse(constants.DATE_LAYOUT, runValue["status_timestamp"].(string))
					buildResponse := BuildsResponse{
						OrgName:       runValue["org_name"].(string),
						ComponentId:   runValue["component_id"].(string),
						ComponentName: runValue["component_name"].(string),
						AutomationId:  runValue["automation_id"].(string),
						WorkflowName:  runValue["workflow_name"].(string),
						Status:        runValue["status"].(string),
						Duration:      runValue["duration"].(float64),
						Build:         runValue["run_number"].(float64),
						Environment:   runValue["target_env"].(string),
						RunId:         runValue["run_id"].(string),
						JobId:         runValue["job_id"].(string),
						StepId:        runValue["step_id"].(string),
						Source:        runValue["source"].(string),
					}
					automationSet[buildResponse.AutomationId] = struct{}{}
					if err == nil {
						buildResponse.StartTimeInMillis = int64(float64(statusTime.UnixMilli()) - runValue["duration"].(float64))
						buildResponse.EndTimeInMillis = int64(float64(statusTime.UnixMilli()))
						buildResponse.StartTime = time.UnixMilli(buildResponse.StartTimeInMillis).Format(constants.DATE_FORMAT)
						buildResponse.EndTime = time.UnixMilli(buildResponse.EndTimeInMillis).Format(constants.DATE_FORMAT)
					}
					buildData, ok := buildMap[buildResponse.AutomationId]
					if ok {
						buildData = append(buildData, buildResponse)
						buildMap[buildResponse.AutomationId] = buildData
					} else {
						buildResponses := []BuildsResponse{}
						buildResponses = append(buildResponses, buildResponse)
						buildMap[buildResponse.AutomationId] = buildResponses
					}
				}
			}
		}
	}
	log.Debugf("Time took to process all auitomation and last active time for build drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}
	reports := &structpb.ListValue{}
	responseList := []BuildsResponse{}
	startTime = time.Now()
	var automationResponse map[string]WorkflowReportResponse
	if ok && branch != nil && len(branch.(string)) > 0 {
		automationResponse = getAutomationResponseMapForBranch(ctx, clt, orgId, components, automationSet, false, branch.(string))
	} else {
		automationResponse = getAutomationResponseMap(ctx, clt, orgId, components, automationSet, false)
	}
	log.Debugf("Time took to automation info from cache for build drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for id, builds := range buildMap {
		automation, ok := automationResponse[id]
		for _, build := range builds {
			if ok {
				build.Branch = automation.Branch
				build.BranchId = automation.BranchId
				build.WorkflowName = automation.WorkflowName
			}
			responseList = append(responseList, build)
		}
	}
	log.Debugf("Time took to process all automation from cache and opensearch for build drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].EndTimeInMillis > responseList[j].EndTimeInMillis
	})
	log.Debugf("Time took to sort the result for build drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, response := range responseList {
		startTime, err := helper.ConvertUTCtoTimeZone(response.StartTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, exceptions.ErrConvertingStartTimeToTZ, response.StartTime, replacements["timeZone"].(string))
		convertedTimeStart, _ := convertTimeFormat(startTime, replacements["timeFormat"].(string))
		startTime = convertedTimeStart
		endTime, err := helper.ConvertUTCtoTimeZone(response.EndTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, exceptions.ErrConvertingEndTimeToTZ, response.EndTime, replacements["timeZone"].(string))
		convertedTimeEnd, _ := convertTimeFormat(endTime, replacements["timeFormat"].(string))
		endTime = convertedTimeEnd
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:       structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:    structpb.NewStringValue(response.ComponentId),
						constants.WORKFLOW:        structpb.NewStringValue(response.WorkflowName),
						constants.AUTOMATION_ID:   structpb.NewStringValue(response.AutomationId),
						constants.BRANCH:          structpb.NewStringValue(response.Branch),
						constants.BRANCH_ID:       structpb.NewStringValue(response.BranchId),
						constants.BUILD_ID:        structpb.NewNumberValue(response.Build),
						constants.RUN_ID_KEY:      structpb.NewStringValue(response.RunId),
						constants.START_TIME:      structpb.NewStringValue(startTime),
						constants.END_TIME:        structpb.NewStringValue(endTime),
						constants.STATUS:          structpb.NewStringValue(response.Status),
						constants.ENVIRONMENT:     structpb.NewStringValue(response.Environment),
						constants.DURATION:        structpb.NewNumberValue(response.Duration),
						constants.JOB_ID:          structpb.NewStringValue(response.JobId),
						constants.STEP_ID:         structpb.NewStringValue(response.StepId),
						constants.SOURCE_PROVIDER: structpb.NewStringValue(response.Source),
					},
				},
			},
		})
	}
	return reports, nil
}

func CodeProgressionSnapshotDeployments(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CodeProgressionSnapshotDeploy)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	branch, ok := replacements[constants.BRANCH]
	var modifiedJson string

	if ok && branch != nil && len(branch.(string)) > 0 {
		modifiedJson = internal.UpdateFiltersForDrilldown(updatedJSON, replacements, true, false)
	} else {
		modifiedJson = internal.UpdateFilters(updatedJSON, replacements)
	}
	response, err := getSearchResponse(modifiedJson, constants.DEPLOY_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all deployments for deployment drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	automationSet := make(map[string]struct{})
	deployMap := make(map[string][]DeploymentResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.DEPLOYMENTS] != nil {
			deployments := aggsResult[constants.DEPLOYMENTS].(map[string]interface{})
			if deployments[constants.VALUE] != nil {
				values := deployments[constants.VALUE].(map[string]interface{})
				for _, value := range values {
					runValue := value.(map[string]interface{})
					statusTime, err := time.Parse(constants.DATE_LAYOUT, runValue["status_timestamp"].(string))
					deployResponse := DeploymentResponse{
						OrgName:       runValue["org_name"].(string),
						ComponentId:   runValue["component_id"].(string),
						ComponentName: runValue["component_name"].(string),
						AutomationId:  runValue["automation_id"].(string),
						WorkflowName:  runValue["workflow_name"].(string),
						Status:        runValue["status"].(string),
						Duration:      runValue["duration"].(float64),
						Build:         runValue["run_number"].(float64),
						Environment:   runValue["target_env"].(string),
						RunId:         runValue["run_id"].(string),
						JobId:         runValue["job_id"].(string),
						StepId:        runValue["step_id"].(string),
					}
					if err == nil {
						deployResponse.StartTimeInMillis = int64(float64(statusTime.UnixMilli()) - runValue["duration"].(float64))
						deployResponse.EndTimeInMillis = int64(float64(statusTime.UnixMilli()))
						deployResponse.StartTime = time.UnixMilli(deployResponse.StartTimeInMillis).Format(constants.DATE_FORMAT)
						deployResponse.EndTime = time.UnixMilli(deployResponse.EndTimeInMillis).Format(constants.DATE_FORMAT)
					}
					if deployResponse.Duration == 0 {
						deployResponse.Duration = 1000
					}
					automationSet[deployResponse.AutomationId] = struct{}{}
					deployData, ok := deployMap[deployResponse.AutomationId]
					if ok {
						deployData = append(deployData, deployResponse)
						deployMap[deployResponse.AutomationId] = deployData
					} else {
						deployResponses := []DeploymentResponse{}
						deployResponses = append(deployResponses, deployResponse)
						deployMap[deployResponse.AutomationId] = deployResponses
					}
				}
			}
		}
	}
	log.Debugf("Time took to process all auitomation and last active time for deployment drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}
	reports := &structpb.ListValue{}
	responseList := []DeploymentResponse{}
	startTime = time.Now()
	var automationResponse map[string]WorkflowReportResponse
	if ok && branch != nil && len(branch.(string)) > 0 {
		automationResponse = getAutomationResponseMapForBranch(ctx, clt, orgId, components, automationSet, false, branch.(string))

	} else {
		automationResponse = getAutomationResponseMap(ctx, clt, orgId, components, automationSet, false)

	}
	log.Debugf("Time took to automation info from cache for deployment drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for id, deployments := range deployMap {
		automation, ok := automationResponse[id]
		for _, deployment := range deployments {
			if ok {
				deployment.Branch = automation.Branch
				deployment.BranchId = automation.BranchId
				deployment.WorkflowName = automation.WorkflowName
			}
			responseList = append(responseList, deployment)
		}
	}
	log.Debugf("Time took to process all automation from cache and opensearch for deployment drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].EndTimeInMillis > responseList[j].EndTimeInMillis
	})
	log.Debugf("Time took to sort the result for deployment drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, response := range responseList {
		startTime, err := helper.ConvertUTCtoTimeZone(response.StartTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, exceptions.ErrConvertingStartTimeToTZ, response.StartTime, replacements["timeZone"].(string))
		convertedTimeStart, _ := convertTimeFormat(startTime, replacements["timeFormat"].(string))
		startTime = convertedTimeStart

		endTime, err := helper.ConvertUTCtoTimeZone(response.EndTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, exceptions.ErrConvertingEndTimeToTZ, response.EndTime, replacements["timeZone"].(string))
		convertedTimeEnd, _ := convertTimeFormat(endTime, replacements["timeFormat"].(string))
		endTime = convertedTimeEnd

		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:     structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:  structpb.NewStringValue(response.ComponentId),
						constants.WORKFLOW:      structpb.NewStringValue(response.WorkflowName),
						constants.AUTOMATION_ID: structpb.NewStringValue(response.AutomationId),
						constants.BRANCH:        structpb.NewStringValue(response.Branch),
						constants.BRANCH_ID:     structpb.NewStringValue(response.BranchId),
						constants.DEPLOYMENT_ID: structpb.NewNumberValue(response.Build),
						constants.RUN_ID_KEY:    structpb.NewStringValue(response.RunId),
						constants.DURATION:      structpb.NewNumberValue(response.Duration),
						constants.START_TIME:    structpb.NewStringValue(startTime),
						constants.END_TIME:      structpb.NewStringValue(endTime),
						constants.STATUS:        structpb.NewStringValue(response.Status),
						constants.ENVIRONMENT:   structpb.NewStringValue(response.Environment),
						constants.JOB_ID:        structpb.NewStringValue(response.JobId),
						constants.STEP_ID:       structpb.NewStringValue(response.StepId),
					},
				},
			},
		})
	}
	return reports, nil
}

func SuccessfulBuildDuration(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.SuccessfulBuildDuration)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)
	response, err := getSearchResponse(modifiedJson, constants.BUILD_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all builds for successful build duration drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	automationSet := make(map[string]struct{})
	buildMap := make(map[string][]SuccessfulBuildsResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.BUILDS] != nil {
			builds := aggsResult[constants.BUILDS].(map[string]interface{})
			if builds[constants.VALUE] != nil {
				values := builds[constants.VALUE].(map[string]interface{})
				for _, value := range values {
					runValue := value.(map[string]interface{})
					buildResponse := SuccessfulBuildsResponse{
						ComponentId:       runValue["component_id"].(string),
						ComponentName:     runValue["component_name"].(string),
						AutomationId:      runValue["automation_id"].(string),
						WorkflowName:      runValue["workflow_name"].(string),
						Duration:          runValue["duration"].(float64),
						StartTimeInMillis: int64(runValue["start_time"].(float64)),
						EndTimeInMillis:   int64(runValue["completed_time"].(float64)),
						RunNumber:         runValue["run_number"].(float64),
						Environment:       runValue["target_env"].(string),
						RunId:             runValue["run_id"].(string),
						JobId:             runValue["job_id"].(string),
						StepId:            runValue["step_id"].(string),
						Source:            runValue["source"].(string),
					}
					automationSet[buildResponse.AutomationId] = struct{}{}
					buildResponse.StartTime = time.UnixMilli(buildResponse.StartTimeInMillis).Format(constants.DATE_FORMAT)
					buildResponse.EndTime = time.UnixMilli(buildResponse.EndTimeInMillis).Format(constants.DATE_FORMAT)
					buildData, ok := buildMap[buildResponse.AutomationId]
					if ok {
						buildData = append(buildData, buildResponse)
						buildMap[buildResponse.AutomationId] = buildData
					} else {
						buildResponses := []SuccessfulBuildsResponse{}
						buildResponses = append(buildResponses, buildResponse)
						buildMap[buildResponse.AutomationId] = buildResponses
					}
				}
			}
		}
	}
	log.Debugf("Time took to process all auitomation and last active time for successful build duration drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}
	reports := &structpb.ListValue{}
	responseList := []SuccessfulBuildsResponse{}
	startTime = time.Now()
	automationResponse := getAutomationResponseMap(ctx, clt, orgId, components, automationSet, false)
	log.Debugf("Time took to automation info from cache for successful build duration drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for id, builds := range buildMap {
		for _, build := range builds {
			automation, ok := automationResponse[id]
			if ok {
				build.Branch = automation.Branch
				build.BranchId = automation.BranchId
				build.WorkflowName = automation.WorkflowName
			}
			responseList = append(responseList, build)
		}
	}
	log.Debugf("Time took to process all automation from cache and opensearch for build drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].EndTimeInMillis > responseList[j].EndTimeInMillis
	})
	log.Debugf("Time took to sort the result for build drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, response := range responseList {
		startTime, err := helper.ConvertUTCtoTimeZone(response.StartTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, exceptions.ErrConvertingStartTimeToTZ, response.StartTime, replacements["timeZone"].(string))
		convertedTimeStart, _ := convertTimeFormat(startTime, replacements["timeFormat"].(string))
		startTime = convertedTimeStart

		endTime, err := helper.ConvertUTCtoTimeZone(response.EndTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, exceptions.ErrConvertingEndTimeToTZ, response.EndTime, replacements["timeZone"].(string))
		convertedTimeEnd, _ := convertTimeFormat(endTime, replacements["timeFormat"].(string))
		endTime = convertedTimeEnd

		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:       structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_NAME:  structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:    structpb.NewStringValue(response.ComponentId),
						constants.WORKFLOW:        structpb.NewStringValue(response.WorkflowName),
						constants.AUTOMATION_ID:   structpb.NewStringValue(response.AutomationId),
						constants.BRANCH:          structpb.NewStringValue(response.Branch),
						constants.BRANCH_ID:       structpb.NewStringValue(response.BranchId),
						constants.RUN_ID:          structpb.NewNumberValue(response.RunNumber),
						constants.RUN_ID_KEY:      structpb.NewStringValue(response.RunId),
						constants.START_TIME:      structpb.NewStringValue(startTime),
						constants.END_TIME:        structpb.NewStringValue(endTime),
						constants.ENVIRONMENT:     structpb.NewStringValue(response.Environment),
						constants.DURATION:        structpb.NewNumberValue(response.Duration),
						constants.JOB_ID:          structpb.NewStringValue(response.JobId),
						constants.STEP_ID:         structpb.NewStringValue(response.StepId),
						constants.SOURCE_PROVIDER: structpb.NewStringValue(response.Source),
					},
				},
			},
		})
	}
	return reports, nil
}

func DeploymentOverviewDrilldown(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.DeploymentOverview)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	var modifiedJson string
	branch, ok := replacements[constants.BRANCH]
	if ok && branch != nil && len(branch.(string)) > 0 {
		modifiedJson = internal.UpdateFiltersForDrilldown(updatedJSON, replacements, true, false)
	} else {
		modifiedJson = internal.UpdateFilters(updatedJSON, replacements)
	}
	response, err := getSearchResponse(modifiedJson, constants.DEPLOY_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all deployments for deployment overview drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	automationSet := make(map[string]struct{})
	deployMap := make(map[string][]DeploymentResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.DEPLOYMENTS] != nil {
			deployments := aggsResult[constants.DEPLOYMENTS].(map[string]interface{})
			if deployments[constants.VALUE] != nil {
				values := deployments[constants.VALUE].(map[string]interface{})
				for _, value := range values {
					runValue := value.(map[string]interface{})
					statusTime, err := time.Parse(constants.DATE_LAYOUT, runValue["status_timestamp"].(string))
					deployResponse := DeploymentResponse{
						OrgName:       runValue["org_name"].(string),
						ComponentId:   runValue["component_id"].(string),
						ComponentName: runValue["component_name"].(string),
						AutomationId:  runValue["automation_id"].(string),
						WorkflowName:  runValue["workflow_name"].(string),
						Status:        runValue["status"].(string),
						Duration:      runValue["duration"].(float64),
						Build:         runValue["run_number"].(float64),
						Environment:   runValue["target_env"].(string),
						RunId:         runValue["run_id"].(string),
						JobId:         runValue["job_id"].(string),
						StepId:        runValue["step_id"].(string),
					}
					if deployResponse.Duration == 0 {
						deployResponse.Duration = 1000
					}
					automationSet[deployResponse.AutomationId] = struct{}{}
					if err == nil {
						deployResponse.StartTimeInMillis = int64(float64(statusTime.UnixMilli()) - runValue["duration"].(float64))
						deployResponse.EndTimeInMillis = int64(float64(statusTime.UnixMilli()))
						deployResponse.StartTime = time.UnixMilli(deployResponse.StartTimeInMillis).Format(constants.DATE_FORMAT)
						deployResponse.EndTime = time.UnixMilli(deployResponse.EndTimeInMillis).Format(constants.DATE_FORMAT)
					}
					deployData, ok := deployMap[deployResponse.AutomationId]
					if ok {
						deployData = append(deployData, deployResponse)
						deployMap[deployResponse.AutomationId] = deployData
					} else {
						deployResponses := []DeploymentResponse{}
						deployResponses = append(deployResponses, deployResponse)
						deployMap[deployResponse.AutomationId] = deployResponses
					}
				}
			}
		}
	}
	log.Debugf("Time took to process all auitomation and last active time for deployment overview drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}
	reports := &structpb.ListValue{}
	responseList := []DeploymentResponse{}
	startTime = time.Now()
	automationResponse := getAutomationResponseMap(ctx, clt, orgId, components, automationSet, false)
	log.Debugf("Time took to automation info from cache for deployment drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for id, deployments := range deployMap {
		automation, ok := automationResponse[id]
		for _, deployment := range deployments {
			if ok {
				deployment.Branch = automation.Branch
				deployment.BranchId = automation.BranchId
				deployment.WorkflowName = automation.WorkflowName
			}
			responseList = append(responseList, deployment)
		}
	}
	log.Debugf("Time took to process all automation from cache and opensearch for deployment overview drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].EndTimeInMillis > responseList[j].EndTimeInMillis
	})
	log.Debugf("Time took to sort the result for deployment drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, response := range responseList {
		startTime, err := helper.ConvertUTCtoTimeZone(response.StartTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, exceptions.ErrConvertingStartTimeToTZ, response.StartTime, replacements["timeZone"].(string))
		convertedTimeStart, _ := convertTimeFormat(startTime, replacements["timeFormat"].(string))
		startTime = convertedTimeStart

		endTime, err := helper.ConvertUTCtoTimeZone(response.EndTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, exceptions.ErrConvertingEndTimeToTZ, response.EndTime, replacements["timeZone"].(string))
		convertedTimeEnd, _ := convertTimeFormat(endTime, replacements["timeFormat"].(string))
		endTime = convertedTimeEnd

		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:     structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:  structpb.NewStringValue(response.ComponentId),
						constants.WORKFLOW:      structpb.NewStringValue(response.WorkflowName),
						constants.AUTOMATION_ID: structpb.NewStringValue(response.AutomationId),
						constants.BRANCH:        structpb.NewStringValue(response.Branch),
						constants.BRANCH_ID:     structpb.NewStringValue(response.BranchId),
						constants.RUN_ID:        structpb.NewNumberValue(response.Build),
						constants.RUN_ID_KEY:    structpb.NewStringValue(response.RunId),
						constants.DURATION:      structpb.NewNumberValue(response.Duration),
						constants.START_TIME:    structpb.NewStringValue(startTime),
						constants.END_TIME:      structpb.NewStringValue(endTime),
						constants.STATUS:        structpb.NewStringValue(response.Status),
						constants.ENVIRONMENT:   structpb.NewStringValue(response.Environment),
						constants.JOB_ID:        structpb.NewStringValue(response.JobId),
						constants.STEP_ID:       structpb.NewStringValue(response.StepId),
					},
				},
			},
		})
	}
	return reports, nil
}

func processDrilldownQueryAndSpec(replacements map[string]any, query string, reportId string, response *pb.DrilldownResponse) (*structpb.ListValue, error) {
	client := db.GetOpenSearchClient()
	if client == nil {
		return nil, fmt.Errorf("Failed to establish opensearch connection")
	}
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, query)
	if log.CheckErrorf(err, "could not replace json placeholders in processDrilldownQueryAndSpec()", replacements) {
		return nil, err
	}

	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)

	response1, err := getSearchResponse(modifiedJson, db.DrillDownAliasDefinitionMap[reportId], client)
	if log.CheckErrorf(err, "Error fetching response from OpenSearch in processDrilldownQueryAndSpec() for reportId: %s", reportId) {
		return nil, err
	}

	var listValue *structpb.ListValue
	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(response1), &data)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in processDrilldownQueryAndSpec()") {
		return nil, err
	}

	if data[constants.AGGREGATION] != nil {
		x := data[constants.AGGREGATION].(map[string]interface{})
		if x[constants.DRILLDOWNS] != nil {
			y := x[constants.DRILLDOWNS].(map[string]interface{})
			if y[constants.VALUE] != nil {
				values := y[constants.VALUE].([]interface{})
				var mappedValues []interface{}
				for _, value := range values {
					valueMap := value.(map[string]interface{})
					// change letter cases based on scannerNameMap for subrow extraction
					if scannerName, ok := valueMap["scannerName"].(string); ok {
						if mappedName, exists := scannerNameMap[scannerName]; exists {
							valueMap["scannerName"] = mappedName
						}
					}

					// Handle case where scanner_name is nested within drillDown
					if drillDown, ok := valueMap["drillDown"].(map[string]interface{}); ok {
						if reportInfo, ok := drillDown["reportInfo"].(map[string]interface{}); ok {
							if scannerName, ok := reportInfo["scanner_name"].(string); ok {
								if mappedName, exists := scannerNameMap[scannerName]; exists {
									reportInfo["scanner_name"] = mappedName
								}
							}
						}
					}

					// Check if subRows exist and map scanner names in each subRow
					// Remove the following if block after UI has implemented subrows extraction on every widget
					if subRows, ok := valueMap["subRows"].([]interface{}); ok {
						var mappedSubRows []interface{}
						for _, subRow := range subRows {
							subRowMap := subRow.(map[string]interface{})

							// Handle case where scannerName is present directly in subRow
							if scannerName, ok := subRowMap["scannerName"].(string); ok {
								if mappedName, exists := scannerNameMap[scannerName]; exists {
									subRowMap["scannerName"] = mappedName
								}
							}

							// Handle case where scanner_name is nested within drillDown
							if drillDown, ok := subRowMap["drillDown"].(map[string]interface{}); ok {
								if reportInfo, ok := drillDown["reportInfo"].(map[string]interface{}); ok {
									if scannerName, ok := reportInfo["scanner_name"].(string); ok {
										if mappedName, exists := scannerNameMap[scannerName]; exists {
											reportInfo["scanner_name"] = mappedName
										}
									}
								}
							}
							mappedSubRows = append(mappedSubRows, subRowMap)
						}

						valueMap["subRows"] = mappedSubRows
					}
					mappedValues = append(mappedValues, valueMap)
				}

				listValue, err = structpb.NewList(mappedValues)
				if log.CheckErrorf(err, "Error forming drilldown response in processDrilldownQueryAndSpec()") {
					return nil, err
				}
			}
		}
	}
	return listValue, nil

}

func TestAutomationDrilldown(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TestAutomationDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	testAutomationQueryUpdated := internal.UpdateFilters(updatedJSON, replacements)

	testAutomationRunsQuery, err := db.ReplaceJSONplaceholders(replacements, constants.TestAutomationRunsQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	testAutomationRunsQueryUpdated := internal.UpdateFilters(testAutomationRunsQuery, replacements)

	queryMap := make(map[string]db.DbQuery)
	queryMap["cb_test_suites"] = db.DbQuery{AliasName: constants.TEST_SUITE_INDEX, QueryString: testAutomationQueryUpdated}
	queryMap["automation_run_status"] = db.DbQuery{AliasName: constants.AUTOMATION_RUN_STATUS_INDEX, QueryString: testAutomationRunsQueryUpdated}

	responseMap, err := multiSearchResponse(queryMap)
	if log.CheckErrorf(err, "multi search query failed in testInsights workflows drilldown widget()") {
		return nil, err
	}

	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	testAutomationResponse := constants.TestAutomationDrillDownResponse{}
	err = json.Unmarshal(responseMap["cb_test_suites"], &testAutomationResponse)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallingOpenSearchRespInTestAutomationDrillDown) {
		return nil, err
	}

	testAutomationRunsResponse := constants.TestAutomationRunsResponse{}
	err = json.Unmarshal(responseMap["automation_run_status"], &testAutomationRunsResponse)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallingOpenSearchRespInTestAutomationDrillDown) {
		return nil, err
	}

	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}

	testResponseList := []TestWorkflowReportResponse{}
	startTime = time.Now()
	automationResponse := getAutomationResponseMap(ctx, clt, orgId, components, nil, true)
	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationInfoFromCacheMilliSec, time.Since(startTime).Milliseconds())
	if len(automationResponse) > 0 {
		for _, automation := range automationResponse {
			testResponse := testAutomationResponse.Aggregations.ComponentActivity.Value[automation.ComponentId+"_"+automation.AutomationId]
			workflowResponse := TestWorkflowReportResponse{}
			if testResponse.AutomationID != "" {
				workflowResponse.ComponentId = automation.ComponentId
				workflowResponse.ComponentName = automation.ComponentName
				workflowResponse.WorkflowName = automation.WorkflowName
				workflowResponse.Source = automation.Source
				workflowResponse.Branch = automation.Branch
				workflowResponse.TestSuiteType = constants.WITH_TEST_SUITES
				workflowResponse.TestSuites = int64(len(testResponse.TestSuitesSet))
			} else {
				workflowResponse.ComponentId = automation.ComponentId
				workflowResponse.ComponentName = automation.ComponentName
				workflowResponse.WorkflowName = automation.WorkflowName
				workflowResponse.Source = automation.Source
				workflowResponse.Branch = automation.Branch
				workflowResponse.TestSuiteType = constants.WITHOUT_TEST_SUITES
				workflowResponse.TestSuites = 0
			}

			testRunsResponse := testAutomationRunsResponse.Aggregations.ComponentActivity.Value[automation.ComponentId+"_"+automation.AutomationId]
			if testRunsResponse.AutomationID != "" {
				workflowResponse.RunIdCount = len(testRunsResponse.RunIds)
			}

			testResponseList = append(testResponseList, workflowResponse)
		}
	}

	// Sort testResponseList by RunIdCount in descending order
	sort.Slice(testResponseList, func(i, j int) bool {
		return testResponseList[i].RunIdCount > testResponseList[j].RunIdCount
	})

	log.Debugf("Time took to process all automation from cache and opensearch for test automation drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()

	log.Debugf("Time took to sort the result for test automation drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())

	var reports structpb.ListValue
	byteValue, err := json.Marshal(testResponseList)
	if log.CheckErrorf(err, "error marshaling reponse in TestAutomationDrilldown() :") {
		return nil, err
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, "error unmarshaling reponse in TestAutomationDrilldown() :") {
		return nil, err
	}
	return &reports, nil
}

func TestComponentDrillDown(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {

	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TestComponentDrilldownQuery)
	testQueryUpdated := internal.UpdateFilters(updatedJSON, replacements)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	queryMap := make(map[string]db.DbQuery)
	queryMap["cb_test_suites"] = db.DbQuery{AliasName: constants.TEST_SUITE_INDEX, QueryString: testQueryUpdated}
	responseMap, err := multiSearchResponse(queryMap)
	if log.CheckErrorf(err, "multi search query failed in testComponentWidgetSection()") {
		return nil, err
	}

	testComponentDrillDownResponse := constants.TestComponentDrillDownResponse{}
	err = json.Unmarshal(responseMap["cb_test_suites"], &testComponentDrillDownResponse)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallingOpenSearchRespInTestAutomationDrillDown) {
		return nil, err
	}

	testComponentActivityMap := make(map[string]TestComponentReportResponse)

	for _, value := range testComponentDrillDownResponse.Aggregations.ComponentActivity.Value {
		testComponentResponse := TestComponentReportResponse{
			ComponentId:   value.ComponentID,
			ComponentName: value.ComponentName,
			TestSuites:    int64(len(value.TestSuitesSet)),
		}
		testComponentActivityMap[value.ComponentID] = testComponentResponse
	}

	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}

	testResponseList := []TestComponentReportResponse{}
	serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
	if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
		return nil, err
	}
	if serviceResponse != nil {
		for i := 0; i < len(serviceResponse.GetService()); i++ {
			service := serviceResponse.GetService()[i]
			found := false
			if len(components) > 0 {
				for _, component := range components {
					if service.Id == component {
						found = true
						break
					}
				}
			}
			if found || components == nil {
				componentResponse, ok := testComponentActivityMap[service.Id]

				if ok {
					componentResponse.ComponentId = service.Id
					componentResponse.ComponentName = service.Name
					componentResponse.RepositoryUrl = service.RepositoryUrl
					componentResponse.TestSuiteType = constants.WITH_TEST_SUITES
				} else {
					componentResponse = TestComponentReportResponse{

						ComponentId:   service.Id,
						ComponentName: service.Name,
						RepositoryUrl: service.RepositoryUrl,
						TestSuiteType: constants.WITHOUT_TEST_SUITES,
					}
				}
				testResponseList = append(testResponseList, componentResponse)
			}
		}
	}

	// Sort testResponseList by TestSuites in descending order
	sort.Slice(testResponseList, func(i, j int) bool {
		return testResponseList[i].TestSuites > testResponseList[j].TestSuites
	})

	var reports structpb.ListValue
	byteValue, err := json.Marshal(testResponseList)
	if log.CheckErrorf(err, "error marshaling reponse in TestComponentDrillDown() :") {
		return nil, err
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, "error unmarshaling reponse in TestComponentDrillDown() :") {
		return nil, err
	}
	return &reports, nil
}

func SecurityComponentDrillDown(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.SecurityComponentDrilldownQuery)
	scanQueryUpdated := internal.UpdateFilters(updatedJSON, replacements)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}

	rawScanQuery, err := db.ReplaceJSONplaceholders(replacements, constants.SecurityComponentFilterQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	rawScanQueryUpdated := internal.UpdateFilters(rawScanQuery, replacements)

	queryMap := make(map[string]db.DbQuery)
	queryMap["scan"] = db.DbQuery{AliasName: constants.SECURITY_INDEX, QueryString: scanQueryUpdated}
	queryMap["rawScan"] = db.DbQuery{AliasName: constants.RAW_SCAN_RESULTS_INDEX, QueryString: rawScanQueryUpdated}

	responseMap, err := multiSearchResponse(queryMap)
	if log.CheckErrorf(err, "multi search query failed in securityComponentWidgetSection()") {
		return nil, err
	}

	result := make(map[string]interface{})
	json.Unmarshal(responseMap["scan"], &result)
	componentActivityMap := make(map[string]SecurityComponentReportResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.DISTINCT_COMPONENT] != nil {
			automationRuns := aggsResult[constants.DISTINCT_COMPONENT].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				for key, value := range values {
					scannerList := value.([]interface{})
					// Mapping correct names for each scanner
					for i, scanner := range scannerList {
						if name, ok := scannerNameMap[scanner.(string)]; ok {
							scannerList[i] = name
						}
					}
					componentResponse := SecurityComponentReportResponse{
						Scanners: scannerList,
					}
					componentActivityMap[key] = componentResponse
				}
			}
		}
	}

	type multiSearchResponse struct {
		Aggregations struct {
			DistinctComponent struct {
				Value []string `json:"value"`
			} `json:"distinct_component"`
		} `json:"aggregations"`
	}

	rawScanResult := multiSearchResponse{}
	err = json.Unmarshal(responseMap["rawScan"], &rawScanResult)
	if log.CheckErrorf(err, "Error unmarshaling response getCommitTrends()") {
		return nil, err
	}

	for _, value := range rawScanResult.Aggregations.DistinctComponent.Value {
		k, ok := componentActivityMap[value]
		if !ok {
			scannerList := []interface{}{constants.BUNDLED_SONARQUBE}
			componentResponse := SecurityComponentReportResponse{
				Scanners: scannerList,
			}
			componentActivityMap[value] = componentResponse
		} else {
			sc := k.Scanners
			if !slices.Contains(sc, constants.BUNDLED_SONARQUBE) {
				sc = append(sc, constants.BUNDLED_SONARQUBE)
				k.Scanners = sc
				componentActivityMap[value] = k
			}
		}
	}

	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}

	responseList := []SecurityComponentReportResponse{}
	serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
	if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
		return nil, err
	}

	if serviceResponse != nil {
		for i := 0; i < len(serviceResponse.GetService()); i++ {
			service := serviceResponse.GetService()[i]
			found := false
			if len(components) > 0 {
				for _, component := range components {
					if service.Id == component {
						found = true
						break
					}
				}
			}
			if found || components == nil {
				componentResponse, ok := componentActivityMap[service.Id]

				if ok {
					componentResponse.ComponentId = service.Id
					componentResponse.ComponentName = service.Name
					componentResponse.RepositoryUrl = service.RepositoryUrl
					componentResponse.ScannerType = constants.WITHSCANNERS
				} else {
					componentResponse = SecurityComponentReportResponse{

						ComponentId:   service.Id,
						ComponentName: service.Name,
						RepositoryUrl: service.RepositoryUrl,
						Scanners:      []interface{}{constants.NO_SCANNERS},
						ScannerType:   constants.WITHOUTSCANNERS,
					}
				}
				responseList = append(responseList, componentResponse)
			}
		}
	}

	var reports structpb.ListValue
	byteValue, err := json.Marshal(responseList)
	if log.CheckErrorf(err, exceptions.ErrMarshallingRespInSecurityComponentDrillDown) {
		return nil, err
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallingRespInSecurityComponentDrillDown) {
		return nil, err
	}
	return &reports, nil
}

func SecurityAutomationDrillDown(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.SecurityAutomationDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	response, err := getSearchResponse(updatedJSON, constants.SECURITY_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	automationActivityMap := make(map[string]SecurityWorkflowReportResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.DISTINCT_AUTOMATION] != nil {
			automationRuns := aggsResult[constants.DISTINCT_AUTOMATION].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				for key, value := range values {
					valueMap := value.(map[string]interface{})
					scannerList := (valueMap["Scanner_List"]).([]interface{})
					// Mapping correct names for each scanner
					for i, scanner := range scannerList {
						if name, ok := scannerNameMap[scanner.(string)]; ok {
							scannerList[i] = name
						}
					}
					runCount := int((valueMap["run_count"]).(float64))
					automationResponse := SecurityWorkflowReportResponse{
						Scanners:   scannerList,
						RunIdCount: runCount,
					}
					automationActivityMap[key] = automationResponse
				}
			}
		}
	}
	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationMilliSec, time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}

	responseList := []SecurityWorkflowReportResponse{}
	startTime = time.Now()
	automationResponse := getAutomationResponseMap(ctx, clt, orgId, components, nil, true)
	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationInfoFromCacheMilliSec, time.Since(startTime).Milliseconds())
	if len(automationResponse) > 0 {
		for id, automation := range automationResponse {
			found := false
			if len(components) > 0 {
				for _, component := range components {
					if automation.ComponentId == component {
						found = true
						break
					}
				}
			}
			// for one or all components
			if found || components == nil {
				coreDataCache := cache.GetCoreDataCache()
				if coreDataCache != nil {
					automationResource := coreDataCache.Get(id)
					if !automationResource.IsDisabled {
						workflowReponse, ok := automationActivityMap[id]
						if ok {
							workflowReponse.ComponentId = automation.ComponentId
							workflowReponse.ComponentName = automation.ComponentName
							workflowReponse.WorkflowName = automation.WorkflowName
							workflowReponse.Branch = automation.Branch
							workflowReponse.ScannerType = constants.WITHSCANNERS
						} else {
							workflowReponse = SecurityWorkflowReportResponse{
								ComponentId:   automation.ComponentId,
								ComponentName: automation.ComponentName,
								WorkflowName:  automation.WorkflowName,
								Branch:        automation.Branch,
								Scanners:      []interface{}{constants.NO_SCANNERS},
								RunIdCount:    '0',
								ScannerType:   constants.WITHOUTSCANNERS,
							}
						}
						responseList = append(responseList, workflowReponse)
					}
				}
			}
		}
	}
	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationFromCacheAndOpensearchMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()

	log.Debugf(exceptions.DebugTimeTookToSortAllResultForAutomationMilliSec, time.Since(startTime).Milliseconds())

	var reports structpb.ListValue
	byteValue, err := json.Marshal(responseList)
	if log.CheckErrorf(err, "error marshaling reponse in SecurityAutomationDrillDown() :") {
		return nil, err
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, "error unmarshaling reponse in SecurityAutomationDrillDown() :") {
		return nil, err
	}
	return &reports, nil
}

func SecurityAutomationRunDrillDown(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	//fetch automation run status index repsone
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.SecurityAutomationRunDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)

	response, err := getSearchResponse(modifiedJson, constants.AUTOMATION_RUN_STATUS_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	//unmarshal reponse - contains map with automationID as key, and a []struct as value (each struct corresponds to a runID)
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)

	//fetch security index response
	updatedJSON, err = db.ReplaceJSONplaceholders(replacements, constants.SecurityAutomationRunsDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson = internal.UpdateFilters(updatedJSON, replacements)

	securityResponse, err := getSearchResponse(modifiedJson, constants.SECURITY_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	//unmarshal security index reponse - contains map with runID as key, and a map of scanner names as keys with scanned status as values
	securityResult := make(map[string]interface{})
	json.Unmarshal([]byte(securityResponse), &securityResult)

	scannerRuns := make(map[string]interface{})

	if securityResult[constants.AGGREGATION] != nil {
		aggsResult := securityResult[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.DISTINCT_RUN] != nil {
			automationRuns := aggsResult[constants.DISTINCT_RUN].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				scannerRuns = automationRuns[constants.VALUE].(map[string]interface{})

			}
		}
	}

	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationRunMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()

	automationSet := make(map[string]struct{})
	automationRunMap := make(map[string][]SecurityWorkflowRunReportResponse)

	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.AUTOMATION_RUN_ACTIVITY] != nil {
			automationRuns := aggsResult[constants.AUTOMATION_RUN_ACTIVITY].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				for key, value := range values {
					automationSet[key] = struct{}{}
					runs := value.([]interface{})
					workflowRuns := []SecurityWorkflowRunReportResponse{}
					for _, runValueMap := range runs {
						runValue := runValueMap.(map[string]interface{})
						curRunID := runValue["run_id"].(string)
						statusTime, err := time.Parse(constants.DATE_LAYOUT, runValue["status_timestamp"].(string))
						var scannerList []interface{}
						var scannerType string
						var scanStatus string
						securityScanFlag := true
						if val, ok := scannerRuns[curRunID]; !ok {
							scannerList = []interface{}{constants.NO_SCANNERS}
							scannerType = constants.WITHOUTSCANNERS
							scanStatus = constants.NOT_APPLICABLE
						} else {
							scannersMapList := val.([]interface{})
							// Map the scanner names
							for _, scanner := range scannersMapList {
								scannerMap := scanner.(map[string]interface{})
								for scannerName, scannedStatus := range scannerMap {
									if name, ok := scannerNameMap[scannerName]; ok {
										scannerList = append(scannerList, name)
										if scannedStatus.(string) != constants.AUTOMATION_STATUS_SUCCESS {
											securityScanFlag = false
										}
									}
								}
							}
							scannerType = constants.WITHSCANNERS
							if securityScanFlag {
								scanStatus = constants.SCANNED
							} else {
								scanStatus = constants.NOT_SCANNED
							}
						}
						runInfo := SecurityWorkflowRunReportResponse{
							ComponentId:   runValue["component_id"].(string),
							ComponentName: runValue["component_name"].(string),
							AutomationId:  runValue["automation_id"].(string),
							RunId:         runValue["run_id"].(string),
							Build:         runValue["run_number"].(float64),
							Scanners:      scannerList,
							ScanStatus:    scanStatus,
							ScannerType:   scannerType,
						}
						if err == nil {
							runInfo.RunStartTimeInMillis = int64(float64(statusTime.UnixMilli()) - runValue["duration"].(float64))
						}
						workflowRuns = append(workflowRuns, runInfo)
					}
					automationRunMap[key] = workflowRuns
				}
			}
		}
	}
	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationRunMilliSec, time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}

	responseList := []SecurityWorkflowRunReportResponse{}
	autorunIDs := make(map[string]bool)

	startTime = time.Now()
	automationResponse := getAutomationResponseMap(ctx, clt, orgId, components, automationSet, false)
	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationRunInfoFromCacheMilliSec, time.Since(startTime).Milliseconds())
	for id, runs := range automationRunMap {
		automation, ok := automationResponse[id]
		for _, workflowRun := range runs {
			if ok {
				workflowRun.Branch = automation.Branch
				workflowRun.BranchId = automation.BranchId
				workflowRun.WorkflowName = automation.WorkflowName
			}
			if !autorunIDs[workflowRun.RunId] {
				autorunIDs[workflowRun.RunId] = true
				responseList = append(responseList, workflowRun)
			}
		}
	}

	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationRunFromCacheAndOpensearchMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].RunStartTimeInMillis > responseList[j].RunStartTimeInMillis
	})
	log.Debugf(exceptions.DebugTimeTookToSortAllResultForAutomationRunMilliSec, time.Since(startTime).Milliseconds())

	var reports structpb.ListValue

	byteValue, err := json.Marshal(responseList)
	if log.CheckErrorf(err, exceptions.ErrMarshallingRespInSecurityComponentDrillDown) {
		return nil, err
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallingRespInSecurityComponentDrillDown) {
		return nil, err
	}
	return &reports, nil
}

func SecurityScanTypeWorkflowsDrillDown(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	//fetch automation run status index repsone
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.AutomationRunDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)

	response, err := getSearchResponse(modifiedJson, constants.AUTOMATION_RUN_STATUS_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	//unmarshal reponse - contains map with automationID as key, and a []struct as value (each struct corresponds to a runID)
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)

	//fetch security index response
	updatedJSON, err = db.ReplaceJSONplaceholders(replacements, constants.SecurityScanTypeWorkflowsDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson = internal.UpdateFilters(updatedJSON, replacements)

	securityResponse, err := getSearchResponse(modifiedJson, constants.SECURITY_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	//unmarshal security index reponse - contains map with runID as key, and a []string of scanner names as value
	securityResult := make(map[string]interface{})
	json.Unmarshal([]byte(securityResponse), &securityResult)

	scannerRuns := make(map[string]interface{})

	if securityResult[constants.AGGREGATION] != nil {
		aggsResult := securityResult[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.DISTINCT_RUN] != nil {
			automationRuns := aggsResult[constants.DISTINCT_RUN].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				scannerRuns = automationRuns[constants.VALUE].(map[string]interface{})

			}
		}
	}

	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationRunMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()

	automationSet := make(map[string]struct{})
	automationRunMap := make(map[string][]SecurityScanTypeWorkflowsReportResponse)

	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.AUTOMATION_RUN_ACTIVITY] != nil {
			automationRuns := aggsResult[constants.AUTOMATION_RUN_ACTIVITY].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				for key, value := range values {
					automationSet[key] = struct{}{}
					runs := value.([]interface{})
					workflowRuns := []SecurityScanTypeWorkflowsReportResponse{}
					for _, runValueMap := range runs {
						runValue := runValueMap.(map[string]interface{})
						curRunID := runValue["run_id"].(string)
						statusTime, err := time.Parse(constants.DATE_LAYOUT, runValue["status_timestamp"].(string))
						var scannerList []interface{}
						var scannerTypeList []interface{}
						if val, ok := scannerRuns[curRunID]; ok {
							valueMap := val.(map[string]interface{})
							scannerList = (valueMap["scanner_names"]).([]interface{})
							// Retrieve and map the scanner names
							for i, scanner := range scannerList {
								if name, ok := scannerNameMap[scanner.(string)]; ok {
									scannerList[i] = name
								}
							}
							scannerTypeList = (valueMap["scanner_types"]).([]interface{})
							runInfo := SecurityScanTypeWorkflowsReportResponse{
								ComponentId:   runValue["component_id"].(string),
								ComponentName: runValue["component_name"].(string),
								AutomationId:  runValue["automation_id"].(string),
								RunId:         runValue["run_id"].(string),
								Build:         runValue["run_number"].(float64),
								Scanners:      scannerList,
								ScannerType:   scannerTypeList,
							}
							if err == nil {
								runInfo.RunStartTimeInMillis = int64(float64(statusTime.UnixMilli()) - runValue["duration"].(float64))
							}
							workflowRuns = append(workflowRuns, runInfo)
						}
					}
					automationRunMap[key] = workflowRuns
				}
			}
		}
	}
	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationRunMilliSec, time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}

	responseList := []SecurityScanTypeWorkflowsReportResponse{}

	startTime = time.Now()
	automationResponse := getAutomationResponseMap(ctx, clt, orgId, components, automationSet, false)
	log.Debugf(exceptions.DebugTimeTookToFetchAllAutomationRunInfoFromCacheMilliSec, time.Since(startTime).Milliseconds())
	for id, runs := range automationRunMap {
		automation, ok := automationResponse[id]
		for _, workflowRun := range runs {
			if ok {
				workflowRun.Branch = automation.Branch
				workflowRun.BranchId = automation.BranchId
				workflowRun.WorkflowName = automation.WorkflowName
			}
			responseList = append(responseList, workflowRun)
		}
	}

	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationRunFromCacheAndOpensearchMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].RunStartTimeInMillis > responseList[j].RunStartTimeInMillis
	})
	log.Debugf(exceptions.DebugTimeTookToSortAllResultForAutomationRunMilliSec, time.Since(startTime).Milliseconds())

	var reports structpb.ListValue

	byteValue, err := json.Marshal(responseList)
	if log.CheckErrorf(err, exceptions.ErrMarshallingRespInSecurityComponentDrillDown) {
		return nil, err
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallingRespInSecurityComponentDrillDown) {
		return nil, err
	}
	return &reports, nil
}

func DeploymentFrequencyAndLeadTime(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.DeploymentFrequencyAndLeadTime)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)
	response, err := getSearchResponse(modifiedJson, constants.DEPLOY_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all deployments for deployment frequency and lead time drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	automationSet := make(map[string]struct{})
	deployMap := make(map[string][]DoraMetricsDeploymentResponse)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.DEPLOYMENTS] != nil {
			deployments := aggsResult[constants.DEPLOYMENTS].(map[string]interface{})
			if deployments[constants.VALUE] != nil {
				values := deployments[constants.VALUE].([]interface{})
				for _, value := range values {
					runValue := value.(map[string]interface{})
					originalRunStartTime := runValue["run_start_time_string_zoned"].(string)
					parsedTime, _ := time.Parse("2006-01-02 15:04:05", originalRunStartTime)
					formattedRunStartTime := parsedTime.Format(constants.DATE_FORMAT)
					convertedRunStartTime, _ := convertTimeFormat(formattedRunStartTime, replacements["timeFormat"].(string))
					statusTime, err := time.Parse(constants.DATE_LAYOUT_TZ, runValue["status_timestamp_zoned"].(string))
					deployResponse := DoraMetricsDeploymentResponse{
						ComponentId:          runValue["component_id"].(string),
						ComponentName:        runValue["component_name"].(string),
						AutomationId:         runValue["automation_id"].(string),
						WorkflowName:         runValue["workflow_name"].(string),
						RunStartTimeInMillis: int64(runValue["run_start_time"].(float64)),
						RunStartTime:         convertedRunStartTime,
						RunNumber:            runValue["run_number"].(float64),
						RunId:                runValue["run_id"].(string),
					}
					automationSet[deployResponse.AutomationId] = struct{}{}
					if err == nil {
						deployResponse.DeployTimeInMillis = int64(float64(statusTime.UnixMilli()))
						deployResponse.DeployTime = statusTime.Format(constants.DATE_FORMAT_TZ)
						convertedTimeDeploy, _ := convertTimeFormat(deployResponse.DeployTime, replacements["timeFormat"].(string))
						deployResponse.DeployTime = convertedTimeDeploy
						deployResponse.LeadTimeInMillis = deployResponse.DeployTimeInMillis - deployResponse.RunStartTimeInMillis
					}
					deployData, ok := deployMap[deployResponse.AutomationId]
					if ok {
						deployData = append(deployData, deployResponse)
						deployMap[deployResponse.AutomationId] = deployData
					} else {
						deployResponses := []DoraMetricsDeploymentResponse{}
						deployResponses = append(deployResponses, deployResponse)
						deployMap[deployResponse.AutomationId] = deployResponses
					}
				}
			}
		}
	}
	log.Debugf("Time took to process all auitomation and last active time for deployment frequency and lead time drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}
	reports := &structpb.ListValue{}
	responseList := []DoraMetricsDeploymentResponse{}
	startTime = time.Now()
	automationResponse := getAutomationResponseMap(ctx, clt, orgId, components, automationSet, false)
	log.Debugf("Time took to automation info from cache for deployment frequency and lead time drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for id, deployments := range deployMap {
		automation, ok := automationResponse[id]
		for _, deployment := range deployments {
			if ok {
				deployment.Branch = automation.Branch
				deployment.BranchId = automation.BranchId
				deployment.WorkflowName = automation.WorkflowName
			}
			responseList = append(responseList, deployment)
		}
	}
	log.Debugf("Time took to process all automation from cache and opensearch for deployment frequency and lead time drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].DeployTimeInMillis > responseList[j].DeployTimeInMillis
	})
	log.Debugf("Time took to sort the result for deployment frequency and lead time drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	// convertedTimeEnd, _ := convertTimeFormat(response.RunStartTime, replacements["timeFormat"].(string))
	// endTime = convertedTimeEnd
	for _, response := range responseList {
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:      structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_NAME: structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:   structpb.NewStringValue(response.ComponentId),
						constants.WORKFLOW:       structpb.NewStringValue(response.WorkflowName),
						constants.AUTOMATION_ID:  structpb.NewStringValue(response.AutomationId),
						constants.BRANCH:         structpb.NewStringValue(response.Branch),
						constants.BRANCH_ID:      structpb.NewStringValue(response.BranchId),
						constants.RUN_ID:         structpb.NewNumberValue(response.RunNumber),
						constants.RUN_ID_KEY:     structpb.NewStringValue(response.RunId),
						constants.LEAD_TIME:      structpb.NewNumberValue(float64(response.LeadTimeInMillis)),
						constants.RUN_START_TIME: structpb.NewStringValue(response.RunStartTime),
						constants.DEPLOYED_TIME:  structpb.NewStringValue(response.DeployTime),
					},
				},
			},
		})
	}
	return reports, nil
}

func FailureRate(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.FailureRate)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)
	response, err := getSearchResponse(modifiedJson, constants.DEPLOY_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all deployments for failure rate drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	responseList := []FailureRateResponse{}
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.DEPLOYMENTS] != nil {
			deployments := aggsResult[constants.DEPLOYMENTS].(map[string]interface{})
			if deployments[constants.VALUE] != nil {
				values := deployments[constants.VALUE].(map[string]interface{})
				for _, value := range values {
					runValue := value.(map[string]interface{})
					deployResponse := FailureRateResponse{
						ComponentId:   runValue["component_id"].(string),
						ComponentName: runValue["component_name"].(string),
						Deployments:   runValue["deployments"].(float64),
						Success:       runValue["success"].(float64),
						Failure:       runValue["failure"].(float64),
					}
					deployResponse.FailureRate = math.Round((deployResponse.Failure/deployResponse.Deployments*100)*100.0) / 100.0
					responseList = append(responseList, deployResponse)
				}
			}
		}
	}
	log.Debugf("Time took to process all auitomation and last active time for failure rate drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	reports := &structpb.ListValue{}
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].Deployments > responseList[j].Deployments
	})
	log.Debugf("Time took to sort the result for failure rate drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, response := range responseList {
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:      structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_NAME: structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:   structpb.NewStringValue(response.ComponentId),
						constants.DEPLOYMENTS:    structpb.NewNumberValue(float64(response.Deployments)),
						constants.SUCCESS:        structpb.NewNumberValue(float64(response.Success)),
						constants.FAILURE:        structpb.NewNumberValue(float64(response.Failure)),
						constants.FAILURE_RATE:   structpb.NewNumberValue(float64(response.FailureRate)),
					},
				},
			},
		})
	}
	return reports, nil
}

func DoraMetricsMttr(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.DoraMetricsMttr)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)
	response, err := getSearchResponse(modifiedJson, constants.DEPLOY_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all deployments for mttr drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	responseList := []DoraMttrResponse{}
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.DEPLOYMENTS] != nil {
			deployments := aggsResult[constants.DEPLOYMENTS].(map[string]interface{})
			if deployments[constants.VALUE] != nil {
				values := deployments[constants.VALUE].([]interface{})
				for _, value := range values {
					runValue := value.(map[string]interface{})
					deployResponse := DoraMttrResponse{
						ComponentId:              runValue["component_id"].(string),
						ComponentName:            runValue["component_name"].(string),
						FailedTimeInMillis:       runValue["failed_on"].(float64),
						RecoveredTimeInMillis:    runValue["recovered_on"].(float64),
						RecoveryDurationInMillis: runValue["recovered_duration"].(float64),
						FailedRunNumber:          int64(runValue["failed_run_number"].(float64)),
						RecoveredRunNumber:       int64(runValue["recovered_run_number"].(float64)),
					}
					deployResponse.FailedTime = time.UnixMilli(int64(deployResponse.FailedTimeInMillis)).Format(constants.DATE_FORMAT)
					deployResponse.RecoveredTime = time.UnixMilli(int64(deployResponse.RecoveredTimeInMillis)).Format(constants.DATE_FORMAT)
					responseList = append(responseList, deployResponse)
				}
			}
		}
	}
	log.Debugf("Time took to process all auitomation for mttr drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	reports := &structpb.ListValue{}
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].RecoveredTimeInMillis > responseList[j].RecoveredTimeInMillis
	})
	log.Debugf("Time took to sort the result for mttr drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for _, response := range responseList {
		failedTime, err := helper.ConvertUTCtoTimeZone(response.FailedTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, "Error converting failed time - %s to timezone - %s", response.FailedTime, replacements["timeZone"].(string))
		recoveredTime, err := helper.ConvertUTCtoTimeZone(response.RecoveredTime, replacements["timeZone"].(string))
		log.CheckErrorf(err, "Error converting recovered time - %s to timezone - %s", response.RecoveredTime, replacements["timeZone"].(string))
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPONENT:         structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_NAME:    structpb.NewStringValue(response.ComponentName),
						constants.COMPONENT_ID:      structpb.NewStringValue(response.ComponentId),
						constants.FAILED_ON:         structpb.NewStringValue(failedTime),
						constants.RECOVERED_ON:      structpb.NewStringValue(recoveredTime),
						constants.RECOVERY_DURATION: structpb.NewNumberValue(float64(response.RecoveryDurationInMillis)),
						constants.FAILED_RUN:        structpb.NewNumberValue(float64(response.FailedRunNumber)),
						constants.RECOVERED_RUN:     structpb.NewNumberValue(float64(response.RecoveredRunNumber)),
					},
				},
			},
		})
	}
	return reports, nil
}

func CiInsightsPluginsInfo(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiToolInsightPluginsFetchQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	ciToolId := replacements[constants.CI_TOOL_ID].(string)
	replacements[constants.ENDPOINT_IDS] = []string{ciToolId}
	modifiedJson := internal.UpdateMustFilters(updatedJSON, replacements)
	response, err := getSearchResponse(modifiedJson, constants.CB_CI_TOOL_INSIGHT_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	log.Debugf("Time took to fetch all plugins for ci insight drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	responseList := []Plugin{}
	hits := result[constants.HITS].(map[string]interface{})[constants.HITS].([]interface{})
	if len(hits) > 0 {
		source := hits[0].(map[string]interface{})[constants.SOURCE].(map[string]interface{})
		if pluginList, ok := source[constants.PLUGINS]; ok {
			for _, plugin := range pluginList.([]interface{}) {
				pluginObject := plugin.(map[string]interface{})
				plugin := Plugin{}
				if pluginObject[constants.LONG_NAME] != nil {
					plugin.LongName = pluginObject[constants.LONG_NAME].(string)
				}
				if pluginObject[constants.SHORT_NAME] != nil {
					plugin.ShortName = pluginObject[constants.SHORT_NAME].(string)
				}
				if pluginObject[constants.VERSION] != nil {
					plugin.Version = pluginObject[constants.VERSION].(string)
				}
				if pluginObject[constants.ENABLED] != nil {
					if pluginObject[constants.ENABLED].(bool) {
						plugin.Enabled = constants.ENABLE
					} else {
						plugin.Enabled = constants.DISABLED
					}
				} else {
					plugin.Enabled = constants.DISABLED
				}
				if pluginObject[constants.ACTIVE_KEY] != nil {
					plugin.Active = pluginObject[constants.ACTIVE_KEY].(bool)
				} else {
					plugin.Active = false
				}
				if pluginObject[constants.HAS_UPDATE] != nil {
					if pluginObject[constants.HAS_UPDATE].(bool) {
						plugin.HasUpdate = constants.NEW_VERSION_AVAILABLE
					} else {
						plugin.HasUpdate = constants.UP_TO_DATE
					}
				} else {
					plugin.HasUpdate = constants.UP_TO_DATE
				}
				if pluginObject[constants.REQUIRED_CORE_VERSION] != nil {
					plugin.RequiredCoreVersion = pluginObject[constants.REQUIRED_CORE_VERSION].(string)
				}
				if pluginObject[constants.MINIMUM_JAVA_VERSION] != nil {
					plugin.MinimumJavaVersion = pluginObject[constants.MINIMUM_JAVA_VERSION].(string)
				}
				if pluginObject[constants.STATUS] != nil {
					plugin.Status = pluginObject[constants.STATUS].(string)
				}
				if pluginObject[constants.DEPENDENCIES] != nil {
					plugin.Dependencies = &structpb.ListValue{
						Values: []*structpb.Value{},
					}
					for _, dependency := range pluginObject[constants.DEPENDENCIES].([]interface{}) {
						value, err := structpb.NewValue(dependency.(map[string]interface{}))
						if err == nil {
							plugin.Dependencies.Values = append(plugin.Dependencies.Values, value)
						}
					}
				}
				responseList = append(responseList, plugin)
			}
		}
	}
	log.Debugf("Time took to process all plugins for ci insight drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	reports := &structpb.ListValue{}
	for _, response := range responseList {
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.LONG_NAME:             structpb.NewStringValue(response.LongName),
						constants.SHORT_NAME:            structpb.NewStringValue(response.ShortName),
						constants.VERSION:               structpb.NewStringValue(response.Version),
						constants.ACTIVE:                structpb.NewBoolValue(response.Active),
						constants.UPDATES:               structpb.NewStringValue(response.HasUpdate),
						constants.REQUIRED_CORE_VERSION: structpb.NewStringValue(response.RequiredCoreVersion),
						constants.MINIMUM_JAVA_VERSION:  structpb.NewStringValue(response.MinimumJavaVersion),
						constants.DEPENDENCIES:          structpb.NewListValue(response.Dependencies),
						constants.STATUS:                structpb.NewStringValue(response.Enabled),
						// constants.STATUS:                structpb.NewStringValue(response.Status),
					},
				},
			},
		})
	}
	return reports, nil
}

func CiInsightsCompletedRunsAndTime(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	client := db.GetOpenSearchClient()
	jobInfo, err := GetJobInfoByJobId(replacements, client)
	reports := &structpb.ListValue{}
	if jobInfo.JobId != "" && err == nil {
		completedHeader := map[string]interface{}{}
		completedHeader[constants.PROJECT_TYPE] = internal.GetProcessedJobType(jobInfo.Type)

		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiRunsByJobIdQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson := internal.UpdateMustFilters(updatedJSON, replacements)
		response, err := getSearchResponse(modifiedJson, constants.CB_CI_RUN_INFO_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		result := make(map[string]interface{})
		data := []interface{}{}
		totalExecuted, success, failed, aborted, unstable, totalRunTime, notBuilt := 0.0, 0, 0, 0, 0, 0.0, 0
		json.Unmarshal([]byte(response), &result)
		if result[constants.AGGREGATION] != nil {
			aggsResult := result[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.COMPLETED_RUNS] != nil {
				completedRuns := aggsResult[constants.COMPLETED_RUNS].(map[string]interface{})
				if completedRuns[constants.VALUE] != nil {
					values := completedRuns[constants.VALUE].([]interface{})
					for _, value := range values {
						completedRuns := map[string]interface{}{}
						run := value.(map[string]interface{})
						completedRuns[constants.JOB_ID] = run[constants.JOBID].(string)
						completedRuns[constants.RUN_ID_KEY] = run[constants.RUNID].(float64)
						completedRuns[constants.ENDPOINT_ID_KEY] = run[constants.ENDPOINT_ID].(string)
						completedRuns[constants.URL] = run[constants.URL].(string)
						completedRuns[constants.START_TIME] = run[constants.START_TIME_MILLIS].(float64)
						completedRuns[constants.START_TIME_CONVERTED] = run[constants.START_TIME_KEY].(string)
						duration := run[constants.DURATION].(float64)
						if duration < 1000 {
							duration = 1000
						}
						result := run[constants.RESULT].(string)
						completedRuns[constants.RESULT] = strings.ToUpper(result)
						completedRuns[constants.RUN_TIME] = duration
						totalRunTime += duration
						totalExecuted += 1
						if result == constants.SUCCESS_KEY {
							success += 1
						} else if result == constants.FAILED_KEY || result == constants.FAILURE_KEY {
							failed += 1
						} else if result == constants.ABORTED_KEY {
							aborted += 1
						} else if result == constants.UNSTABLE_KEY {
							unstable += 1
						} else if result == constants.NOT_BUILT_KEY {
							notBuilt += 1
						}
						if completedRuns[constants.RESULT] == constants.SUCCESS_KEY {
							completedRuns[constants.RESULT] = "SUCCESSFUL"
						} else if completedRuns[constants.RESULT] == constants.ABORTED_KEY {
							completedRuns[constants.RESULT] = "CANCELED"
						}
						data = append(data, completedRuns)
					}
				}
			}
		}
		if totalExecuted > 0 && totalRunTime > 0 {
			completedHeader[constants.AVERAGE_RUN_TIME] = (totalRunTime / totalExecuted)
		}
		if completedHeader[constants.AVERAGE_RUN_TIME] != nil && completedHeader[constants.AVERAGE_RUN_TIME].(float64) < 1000 {
			completedHeader[constants.AVERAGE_RUN_TIME] = 1000
		}
		completedRunDataMap := make(map[string]interface{})
		completedRunDataMap[constants.TYPE] = constants.BAR_WITH_BOTH_AXIS
		completedRunDataMap[constants.DATA] = data
		colorSchemes := []interface{}{}
		colorScheme := map[string]interface{}{}
		colorScheme[constants.COLOR_0] = constants.COLOR_SCHEME_0_0
		colorScheme[constants.COLOR_1] = constants.COLOR_SCHEME_0_1
		colorSchemes = append(colorSchemes, colorScheme)
		colorScheme1 := map[string]interface{}{}
		colorScheme1[constants.COLOR_0] = constants.COLOR_SCHEME_1_0
		colorScheme1[constants.COLOR_1] = constants.COLOR_SCHEME_1_1
		colorSchemes = append(colorSchemes, colorScheme1)
		colorScheme2 := map[string]interface{}{}
		colorScheme2[constants.COLOR_0] = constants.COLOR_SCHEME_2_0
		colorScheme2[constants.COLOR_1] = constants.COLOR_SCHEME_2_1
		colorSchemes = append(colorSchemes, colorScheme2)
		colorScheme3 := map[string]interface{}{}
		colorScheme3[constants.COLOR_0] = constants.COLOR_SCHEME_3_0
		colorScheme3[constants.COLOR_1] = constants.COLOR_SCHEME_3_1
		colorSchemes = append(colorSchemes, colorScheme3)
		completedRunDataMap[constants.COLOR_SCHEME] = colorSchemes
		completedRunDataMap[constants.LIGHT_COLOR_SCHEME] = colorSchemes
		completedHeader[constants.SUCCESS_RUNS] = success
		completedHeader[constants.FAILED_RUNS] = failed
		completedHeader[constants.ABORTED_RUNS] = aborted
		completedHeader[constants.UNSTABLE_RUNS] = unstable
		completedHeader[constants.NOT_BUILT_RUNS] = notBuilt
		completedHeader[constants.TOTAL_RUN_TIME] = totalRunTime
		completedHeader[constants.TOTAL_EXECUTED] = totalExecuted
		completedHeaderMap, err := structpb.NewStruct(completedHeader)
		if err != nil {
			log.Infof("Exception while forming header struct in runInformation drilldown", err)
		}
		completedRunsMap, err := structpb.NewStruct(completedRunDataMap)
		if err != nil {
			log.Infof("Exception while forming runInformation struct in runInformation drilldown", err)
		}
		reports.Values = append(reports.Values, &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						constants.COMPLETED_RUNS_DATA: structpb.NewStructValue(completedRunsMap),
						constants.HEADER_DETAILS:      structpb.NewStructValue(completedHeaderMap),
					},
				},
			},
		})
	}
	return reports, nil
}

func GetJobInfoByJobId(replacements map[string]any, client *opensearch.Client) (internal.CiJobInfo, error) {
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.CiJobInfoByJobIdQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return internal.CiJobInfo{}, err
	}
	modifiedJson := internal.UpdateMustFilters(updatedJSON, replacements)
	response, err := getSearchResponse(modifiedJson, constants.CB_CI_JOB_INFO_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	hits := result[constants.HITS].(map[string]interface{})[constants.HITS].([]interface{})
	jobInfo := internal.CiJobInfo{}
	if len(hits) > 0 {
		source := hits[0].(map[string]interface{})[constants.SOURCE].(map[string]interface{})
		jobInfo.JobId = source[constants.JOBID].(string)
		jobInfo.JobName = source[constants.JOB_NAME].(string)
		jobInfo.EndpointId = source[constants.ENDPOINT_ID].(string)
		jobInfo.DisplayName = source[constants.DISPLAY_NAME].(string)
		jobInfo.Type = source[constants.TYPE].(string)
	}
	return jobInfo, nil
}

func TestAutomationRunDrillDown(replacements map[string]any, ctx context.Context, clt client.GrpcClient) (*structpb.ListValue, error) {
	startTime := time.Now()
	client := db.GetOpenSearchClient()
	//fetch automation run status index repsone
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TestAutomationRunDrilldownQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)

	response, err := getSearchResponse(modifiedJson, constants.AUTOMATION_RUN_STATUS_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	//unmarshal reponse - contains map with automationID as key, and a []struct as value (each struct corresponds to a runID)
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)

	// fetch test suite index response
	updatedJSON, err = db.ReplaceJSONplaceholders(replacements, constants.TestSuiteQuery)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson = internal.UpdateFilters(updatedJSON, replacements)

	testSuiteresponse, err := getSearchResponse(modifiedJson, constants.TEST_SUITE_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	testAutomationRunResponse := constants.TestAutomationRunDrillDownResponse{}
	json.Unmarshal([]byte(testSuiteresponse), &testAutomationRunResponse)

	log.Debugf("Time took to fetch all automation and last active time for test insights automation run drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	startTime = time.Now()

	automationSet := make(map[string]struct{})
	automationRunMap := make(map[string][]TestWorkflowRunReportResponse)

	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.AUTOMATION_RUN_ACTIVITY] != nil {
			automationRuns := aggsResult[constants.AUTOMATION_RUN_ACTIVITY].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				for key, value := range values {
					automationSet[key] = struct{}{}
					runs := value.([]interface{})
					workflowRuns := []TestWorkflowRunReportResponse{}
					for _, runValueMap := range runs {
						runValue := runValueMap.(map[string]interface{})
						curRunID := runValue["run_id"].(string)
						curAutomationID := runValue["automation_id"].(string)
						curComponentID := runValue["component_id"].(string)

						testResponse := testAutomationRunResponse.Aggregations.TestWorkflowDrilldown.Value[curComponentID+"_"+curAutomationID+"_"+curRunID]

						runInfo := TestWorkflowRunReportResponse{
							ComponentId:   runValue["component_id"].(string),
							ComponentName: runValue["component_name"].(string),
							AutomationId:  runValue["automation_id"].(string),
							RunId:         runValue["run_id"].(string),
							Build:         runValue["run_number"].(float64),
							RunStatus:     runValue["status"].(string),
						}

						runInfo.RunStartTimeInMillis = testResponse.RunStartTime
						if testResponse.RunID == curRunID {
							runInfo.TestSuiteType = constants.WITH_TEST_SUITES
						} else {
							runInfo.TestSuiteType = constants.WITHOUT_TEST_SUITES
						}
						if testResponse.AutomationID != "" {
							runInfo.TestSuites = testResponse.Runs
						} else {
							runInfo.TestSuites = 0
						}

						workflowRuns = append(workflowRuns, runInfo)
					}
					automationRunMap[key] = workflowRuns
				}
			}
		}
	}
	log.Debugf("Time took to process all automation and last active time for test insights automation run drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		orgId = ""
	}
	subOrgId, ok := replacements[constants.SUB_ORG_ID].(string)
	if ok && orgId != subOrgId && subOrgId != "" {
		orgId = subOrgId
	}
	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}

	responseList := []TestWorkflowRunReportResponse{}

	startTime = time.Now()
	automationResponse := getAutomationResponseMap(ctx, clt, orgId, components, automationSet, false)
	autorunIDs := make(map[string]bool)
	log.Debugf("Time took to automation info from cache for test insights automation run drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())
	for id, runs := range automationRunMap {
		automation, ok := automationResponse[id]
		for _, workflowRun := range runs {
			if ok {
				workflowRun.Branch = automation.Branch
				workflowRun.BranchId = automation.BranchId
				workflowRun.WorkflowName = automation.WorkflowName
				workflowRun.Source = automation.Source
			}
			if !autorunIDs[workflowRun.RunId] {
				autorunIDs[workflowRun.RunId] = true
				responseList = append(responseList, workflowRun)
			}
		}
	}

	log.Debugf(exceptions.DebugTimeTookToProcessAllAutomationRunFromCacheAndOpensearchMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()
	sort.Slice(responseList, func(i, j int) bool {
		return responseList[i].RunStartTimeInMillis > responseList[j].RunStartTimeInMillis
	})
	log.Debugf("Time took to sort the result for test insights automation run drilldown : %v in milliseconds", time.Since(startTime).Milliseconds())

	var reports structpb.ListValue

	byteValue, err := json.Marshal(responseList)
	if log.CheckErrorf(err, "error marshaling reponse in TestAutomationRunDrillDown() :") {
		return nil, err
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, "error unmarshaling reponse in TestAutomationRunDrillDown() :") {
		return nil, err
	}
	return &reports, nil
}

func TestInsightsViewRunActivityDrillDown(replacements map[string]any, ctx context.Context, grpcClient client.GrpcClient) (*structpb.ListValue, error) {

	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.ViewRunActivityDrilldownQuery)
	if log.CheckErrorf(err, "could not replace json placeholders for View Run Activity DrillDown:", replacements) {
		return nil, db.ErrInternalServer
	}

	response, err := getSearchResponse(updatedJSON, constants.TEST_CASES_INDEX, client)
	if log.CheckErrorf(err, "Error fetching Opensearch data for View Run Activity DrillDown. Fetch failed - ") {
		return nil, db.ErrInternalServer
	}

	type ViewRunActivityInput struct {
		Aggregations struct {
			ViewRunActivity struct {
				Value struct {
					Headers struct {
						Total       int     `json:"total"`
						Failed      int     `json:"FAILED"`
						Skipped     int     `json:"SKIPPED"`
						Passed      int     `json:"PASSED"`
						AvgDuration float64 `json:"avg_duration"`
						Workflow    string  `json:"workflow"`
						Source      string  `json:"source"`
					} `json:"headers"`
					Section []struct {
						JobId            string  `json:"job_id"`
						TestCaseStatus   string  `json:"test_case_status"`
						BuildId          string  `json:"build_id"`
						TestCaseDuration int     `json:"test_case_duration"`
						StartTime        float64 `json:"start_time"`
						RunId            string  `json:"run_id"`
						TestCaseName     string  `json:"test_case_name"`
						TestSuiteName    string  `json:"test_suite_name"`
						BranchId         string  `json:"branch_id"`
					} `json:"section"`
				} `json:"value"`
			} `json:"viewRunActivity"`
		} `json:"aggregations"`
	}

	type HeaderData struct {
		TotalRuns      int     `json:"totalRuns"`
		SuccessfulRuns int     `json:"successfulRuns"`
		FailedRuns     int     `json:"failedRuns"`
		SkippedRuns    int     `json:"skippedRuns"`
		AvgRunTime     float64 `json:"avgRunTime"`
		Workflow       string  `json:"workflow"`
		Source         string  `json:"source"`
	}

	type SectionData struct {
		JobId         string                           `json:"jobId"`
		Status        string                           `json:"status"`
		BuildId       string                           `json:"buildId"`
		RunTime       int                              `json:"runTime"`
		StartTime     float64                          `json:"startTime"`
		DrillDownInfo internal.DrillDownWithReportInfo `json:"drillDown"`
	}

	type ViewRunActivityOutput struct {
		TestCasesActivityData struct {
			Data             []SectionData `json:"data"`
			Type             string        `json:"type"`
			ColorScheme      []ColorScheme `json:"colorScheme"`
			LightColorScheme []ColorScheme `json:"lightColorScheme"`
			ShowLegends      bool          `json:"showLegends"`
		} `json:"testCasesActivityData"`
		HeaderDetails HeaderData `json:"headerDetails"`
	}

	var input ViewRunActivityInput
	err = json.Unmarshal([]byte(response), &input)
	if log.CheckErrorf(err, "Error unmarshalling OpenSearch response into struct in View Run Activity drill down ") {
		return nil, db.ErrInternalServer
	}

	var output []ViewRunActivityOutput
	var section []SectionData

	var workflowName, sourceName string

	if input.Aggregations.ViewRunActivity.Value.Headers.Workflow != "" {
		workflowName, sourceName, err = cutils.GetDisplayNameAndOrigin(input.Aggregations.ViewRunActivity.Value.Headers.Workflow)
		if log.CheckErrorf(err, "Error getting display name and origin for TestInsightsViewRunActivityDrillDown") {
			log.Infof("Error getting display name and origin for workflow")
			return nil, db.ErrInternalServer
		}
	} else {
		log.Infof("Workflow name is nil or empty in TestInsightsViewRunActivityDrillDown, skipping workflow parsing")
	}

	if input.Aggregations.ViewRunActivity.Value.Section != nil {
		for _, sectionEle := range input.Aggregations.ViewRunActivity.Value.Section {
			section = append(section, SectionData{
				JobId:     sectionEle.JobId,
				Status:    sectionEle.TestCaseStatus,
				BuildId:   sectionEle.BuildId,
				RunTime:   sectionEle.TestCaseDuration,
				StartTime: sectionEle.StartTime,
				DrillDownInfo: internal.DrillDownWithReportInfo{
					ReportId: "test-overview-view-run-activity-logs",
					ReportInfo: pb.ReportInfo{
						TestSuiteName: sectionEle.TestSuiteName,
						TestCaseName:  sectionEle.TestCaseName,
						RunId:         sectionEle.RunId,
						RunNumber:     sectionEle.BuildId,
						Branch:        sectionEle.BranchId,
					},
				},
			})
		}

		output = append(output, ViewRunActivityOutput{
			TestCasesActivityData: struct {
				Data             []SectionData `json:"data"`
				Type             string        `json:"type"`
				ColorScheme      []ColorScheme `json:"colorScheme"`
				LightColorScheme []ColorScheme `json:"lightColorScheme"`
				ShowLegends      bool          `json:"showLegends"`
			}{
				Data: section,
				Type: pb.ChartType_BAR_WITH_BOTH_AXIS.String(),
				ColorScheme: []ColorScheme{
					{
						Color0: "#00BFA8",
						Color1: "#056459",
					},
					{
						Color0: "#ED5252",
						Color1: "#640505",
					},
					{
						Color0: "#FFA726",
						Color1: "#B96E00",
					},
				},
				LightColorScheme: []ColorScheme{
					{
						Color0: "#00BFA8",
						Color1: "#056459",
					},
					{
						Color0: "#ED5252",
						Color1: "#640505",
					},
					{
						Color0: "#FFA726",
						Color1: "#B96E00",
					},
				},
				ShowLegends: true,
			},
			HeaderDetails: HeaderData{
				TotalRuns:      input.Aggregations.ViewRunActivity.Value.Headers.Total,
				SuccessfulRuns: input.Aggregations.ViewRunActivity.Value.Headers.Passed,
				FailedRuns:     input.Aggregations.ViewRunActivity.Value.Headers.Failed,
				SkippedRuns:    input.Aggregations.ViewRunActivity.Value.Headers.Skipped,
				AvgRunTime:     input.Aggregations.ViewRunActivity.Value.Headers.AvgDuration,
				Workflow:       workflowName,
				Source:         sourceName,
			},
		})

		var reports structpb.ListValue
		byteValue, err := json.Marshal(output)
		if log.CheckErrorf(err, "error marshaling response in View Run Activity DrillDown:") {
			return nil, db.ErrInternalServer
		}
		err = protojson.Unmarshal(byteValue, &reports)
		if log.CheckErrorf(err, "error unmarshaling response in View Run Activity DrillDown:") {
			return nil, db.ErrInternalServer
		}
		return &reports, nil
	}

	log.Error("", fmt.Errorf("OpenSearch response unmarshalled into struct in View Run Activity drill down is incomplete"))
	return nil, db.ErrInternalServer

}
func TestInsightsViewRunActivityLogsDrillDown(replacements map[string]any, ctx context.Context, grpcClient client.GrpcClient) (*structpb.ListValue, error) {

	client := db.GetOpenSearchClient()
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.ViewRunActivityLogsDrilldownQuery)
	if log.CheckErrorf(err, "could not replace json placeholders in View Run Activity Logs DrillDown:", replacements) {
		return nil, db.ErrInternalServer
	}

	// adding components term filter
	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)

	response, err := getSearchResponse(modifiedJson, constants.TEST_CASES_INDEX, client)
	if log.CheckErrorf(err, "Error fetching Opensearch data for View Run Activity Logs DrillDown. Fetch failed - ") {
		return nil, db.ErrInternalServer
	}

	type ViewRunActivityLogsInput struct {
		Hits struct {
			Hits []struct {
				Source struct {
					StdErr       string `json:"std_err"`
					BuildId      int    `json:"run_number"`
					ErrorTrace   string `json:"error_trace"`
					StdOut       string `json:"std_out"`
					BranchId     string `json:"branch_id"`
					AutomationId string `json:"automation_id"`
					ComponentId  string `json:"component_id"`
					RunId        string `json:"run_id"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	type ViewRunActivityLogsOutput struct {
		BuildId      string `json:"buildId"`
		BranchId     string `json:"branchId"`
		AutomationId string `json:"automationId"`
		ComponentId  string `json:"componentId"`
		RunId        string `json:"runId"`
		Message      string `json:"message"`
	}

	var input ViewRunActivityLogsInput
	err = json.Unmarshal([]byte(response), &input)
	if log.CheckErrorf(err, "Error unmarshalling OpenSearch response into struct in View Run Activity Logs drill down") {
		return nil, db.ErrInternalServer
	}

	if input.Hits.Hits != nil && len(input.Hits.Hits) > 0 {
		source := input.Hits.Hits[0].Source
		buildIdString := fmt.Sprintf("%d", source.BuildId)
		var logMessage string

		hasStdOut := source.StdOut != ""
		hasStdErr := source.StdErr != ""
		hasErrorTrace := source.ErrorTrace != ""

		if !hasStdOut && !hasStdErr && !hasErrorTrace {
			logMessage = "This test case did not report any output."
		} else {
			if hasStdOut {
				logMessage = logMessage + source.StdOut + "\n"
			}
			if hasStdErr {
				logMessage = logMessage + source.StdErr + "\n"
			}
			if hasErrorTrace {
				logMessage = logMessage + source.ErrorTrace + "\n"
			}

		}

		output := []struct {
			LogDetails ViewRunActivityLogsOutput `json:"logDetails"`
		}{
			{
				LogDetails: ViewRunActivityLogsOutput{
					BuildId:      buildIdString,
					BranchId:     source.BranchId,
					AutomationId: source.AutomationId,
					ComponentId:  source.ComponentId,
					RunId:        source.RunId,
					Message:      logMessage,
				},
			},
		}

		var reports structpb.ListValue
		byteValue, err := json.Marshal(output)
		if log.CheckErrorf(err, "error marshaling response in View Run Activity Logs DrillDown") {
			return nil, db.ErrInternalServer
		}
		err = protojson.Unmarshal(byteValue, &reports)
		if log.CheckErrorf(err, "error unmarshaling response in View Run Activity Logs DrillDown") {
			return nil, db.ErrInternalServer
		}
		return &reports, nil
	}

	log.Error("", fmt.Errorf("OpenSearch response unmarshalled into struct in View Run Activity Logs drill down is incomplete"))
	return nil, db.ErrInternalServer

}

/*
TotalTestCasesDrillDown fetches all the test cases under an automation_id,
or a test suite under an automation_id.
*/
func TotalTestCasesDrillDown(replacements map[string]any, ctx context.Context, grpcClient client.GrpcClient) (*structpb.ListValue, error) {

	client := db.GetOpenSearchClient()
	if client == nil {
		return nil, fmt.Errorf("Failed to establish opensearch connection")
	}
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TotalTestCasesDrilldownQuery)
	if log.CheckErrorf(err, "could not replace json placeholders in processDrilldownQueryAndSpec()", replacements) {
		return nil, err
	}

	modifiedJson := internal.UpdateFilters(updatedJSON, replacements)
	resultJson := modifiedJson

	// adding test suite filter dynamically
	if testSuiteName, ok := replacements["testSuiteName"]; ok {
		var data map[string]interface{}
		err = json.Unmarshal([]byte(modifiedJson), &data)
		if log.CheckErrorf(err, "Error unmarshalling modifiedJson in TotalTestCasesDrillDown()") {
			return nil, err
		}

		filterArray := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})

		filter := helper.AddTermFilter("test_suite_name", testSuiteName.(string))
		filterArray = append(filterArray, filter)

		data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterArray
		modifiedData, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Warn(err.Error())
			modifiedData = []byte(modifiedJson)
		}
		resultJson = string(modifiedData)
	}

	response1, err := getSearchResponse(resultJson, constants.TEST_CASES_INDEX, client)
	if log.CheckErrorf(err, "Error fetching response from OpenSearch in TotalTestCasesDrillDown()") {
		return nil, err
	}

	var listValue *structpb.ListValue
	queryResponse := make(map[string]interface{})
	err = json.Unmarshal([]byte(response1), &queryResponse)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in processDrilldownQueryAndSpec()") {
		return nil, err
	}

	if queryResponse[constants.AGGREGATION] != nil {
		x := queryResponse[constants.AGGREGATION].(map[string]interface{})
		if x[constants.DRILLDOWNS] != nil {
			y := x[constants.DRILLDOWNS].(map[string]interface{})
			if y[constants.VALUE] != nil {
				values := y[constants.VALUE].([]interface{})

				listValue, err = structpb.NewList(values)
				if log.CheckErrorf(err, "Error forming drilldown response in TotalTestCasesDrillDown()") {
					return nil, err
				}
			}
		}
	}
	return listValue, nil
}

func TestInsightsTotalRunsDrillDown(replacements map[string]any, ctx context.Context, grpcClient client.GrpcClient) (*structpb.ListValue, error) {

	client := db.GetOpenSearchClient()
	if client == nil {
		log.Error("", fmt.Errorf("Failed to establish OpenSearch connection in TestInsightsTotalRunsDrillDown"))
		return nil, db.ErrInternalServer
	}

	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TestOverviewTotalRunsDrilldownQuery)
	if log.CheckErrorf(err, "could not replace json placeholders in TestInsightsTotalRunsDrillDown:", replacements) {
		return nil, db.ErrInternalServer
	}

	response, err := getSearchResponse(updatedJSON, constants.TEST_CASES_INDEX, client)
	if log.CheckErrorf(err, "Error fetching Opensearch data in TestInsightsTotalRunsDrillDown") {
		return nil, db.ErrInternalServer
	}

	reports, err := transformTestInsightsTotalRunsDrillDown(response)
	if err != nil {
		if err == db.ErrNoDataFound {
			log.Info("No data found in TestInsightsTotalRunsDrillDown")
			return nil, db.ErrNoDataFound
		} else {
			log.Errorf(err, "Error transforming response in TestInsightsTotalRunsDrillDown")
			return nil, db.ErrInternalServer
		}
	}
	return reports, nil

}

func transformTestInsightsTotalRunsDrillDown(response string) (*structpb.ListValue, error) {

	var totalRunsDrillDownQueryResponse struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
		} `json:"hits"`
		Aggregations struct {
			TestCasesThatFailedAtLeastOnce struct {
				DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
				SumOtherDocCount        int `json:"sum_other_doc_count"`
				Buckets                 []struct {
					Key          string `json:"key"`
					DocCount     int    `json:"doc_count"`
					TestCaseName struct {
						Hits struct {
							Total struct {
								Value int `json:"value"`
							} `json:"total"`
							Hits []struct {
								Source struct {
									TestCaseName  string `json:"test_case_name"`
									TestSuiteName string `json:"test_suite_name"`
								} `json:"_source"`
							} `json:"hits"`
						} `json:"hits"`
					} `json:"test_case_name"`
					TestCaseStatusHistory struct {
						DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
						SumOtherDocCount        int `json:"sum_other_doc_count"`
						Buckets                 []struct {
							Key      string `json:"key"`
							DocCount int    `json:"doc_count"`
						} `json:"buckets"`
					} `json:"test_case_status_history"`
					FailedDocs struct {
						DocCount int `json:"doc_count"`
					} `json:"failed_docs"`
				} `json:"buckets"`
			} `json:"test_cases_that_failed_at_least_once"`
			Runs struct {
				DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
				SumOtherDocCount        int `json:"sum_other_doc_count"`
				Buckets                 []struct {
					Key        string `json:"key"`
					DocCount   int    `json:"doc_count"`
					RunDetails struct {
						Hits struct {
							Total struct {
								Value int `json:"value"`
							} `json:"total"`
							Hits []struct {
								Source struct {
									AutomationID string `json:"automation_id"`
									ComponentID  string `json:"component_id"`
									RunID        string `json:"run_id"`
									BranchID     string `json:"branch_id"`
									OrgID        string `json:"org_id"`
									RunNumber    int    `json:"run_number"`
									RunStatus    string `json:"run_status"`
								} `json:"_source"`
								Fields struct {
									ZonedRunStartTime []string `json:"zoned_run_start_time"`
								} `json:"fields"`
							} `json:"hits"`
						} `json:"hits"`
					} `json:"run_details"`
					FailedDocs struct {
						DocCount        int `json:"doc_count"`
						FailedTestCases struct {
							DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
							SumOtherDocCount        int `json:"sum_other_doc_count"`
							Buckets                 []struct {
								Key      string `json:"key"`
								DocCount int    `json:"doc_count"`
							} `json:"buckets"`
						} `json:"failed_test_cases"`
					} `json:"failed_docs"`
					TotalTestCasesCount struct {
						Value float64 `json:"value"`
					} `json:"total_test_cases_count"`
				} `json:"buckets"`
			} `json:"runs"`
		} `json:"aggregations"`
	}

	err := json.Unmarshal([]byte(response), &totalRunsDrillDownQueryResponse)
	if log.CheckErrorf(err, "Error unmarshalling OpenSearch response into struct in TestInsightsTotalRunsDrillDown") {
		return nil, db.ErrInternalServer
	}

	if totalRunsDrillDownQueryResponse.Hits.Total.Value == 0 {
		return nil, db.ErrNoDataFound
	}

	type failureRateTotalRunsDrillDown struct {
		ColorScheme      []internal.ColorScheme `json:"colorScheme"`
		LightColorScheme []internal.ColorScheme `json:"lightColorScheme"`
		Type             string                 `json:"type"`
		Value            string                 `json:"value"`
		Data             []struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		} `json:"data"`
	}

	type subRowsTotalRunsDrillDown struct {
		FailedTestName   string                        `json:"failedTestName"`
		FailureRate      failureRateTotalRunsDrillDown `json:"failureRate"`
		FailureRateValue float32                       `json:"failureRateValue"`
		ViewRunActivity  struct {
			DrillDown internal.DrillDownWithReportInfo `json:"drillDown"`
		} `json:"viewRunActivity"`
	}

	type totalRunsDrillDownOutput struct {
		AutomationID   string                      `json:"automationId"`
		BranchID       string                      `json:"branchId"`
		Build          string                      `json:"build"`
		ComponentID    string                      `json:"componentId"`
		TotalTests     int                         `json:"totalTests"`
		FailedTests    int                         `json:"failedTests"`
		OrganizationID string                      `json:"organizationId"`
		RunID          string                      `json:"runId"`
		RunStatus      string                      `json:"runStatus"`
		RunTime        string                      `json:"runTime"`
		SubRows        []subRowsTotalRunsDrillDown `json:"subRows"`
	}

	type testCaseStatusHistory struct {
		Failed  int
		Skipped int
		Passed  int
	}

	type testCaseInfo struct {
		TestCaseStatusHistory testCaseStatusHistory
		TestCaseName          string
		TestSuiteName         string
		FailurRate            float32
	}

	// Map to hold information related to each test case under the filters selected. Key is test_suite_name + '_' + test_case_name
	allFailedTestCases := make(map[string]testCaseInfo)

	var output []totalRunsDrillDownOutput

	// Forming one map with information related to all test cases instead of looping through the Buckets array in the response every time
	for _, uniqueFailedTestCase := range totalRunsDrillDownQueryResponse.Aggregations.TestCasesThatFailedAtLeastOnce.Buckets {

		var failed, skipped, passed int
		var testCaseName, testSuiteName string

		for _, testCaseStatus := range uniqueFailedTestCase.TestCaseStatusHistory.Buckets {
			if testCaseStatus.Key == "FAILED" {
				failed = testCaseStatus.DocCount
			} else if testCaseStatus.Key == "SKIPPED" {
				skipped = testCaseStatus.DocCount
			} else if testCaseStatus.Key == "PASSED" {
				passed = testCaseStatus.DocCount
			}
		}

		failureRate := float32(failed) * 100 / float32(failed+skipped+passed)
		if uniqueFailedTestCase.TestCaseName.Hits.Hits != nil && len(uniqueFailedTestCase.TestCaseName.Hits.Hits) > 0 {
			testCaseName = uniqueFailedTestCase.TestCaseName.Hits.Hits[0].Source.TestCaseName
			testSuiteName = uniqueFailedTestCase.TestCaseName.Hits.Hits[0].Source.TestSuiteName
		}

		allFailedTestCases[uniqueFailedTestCase.Key] = testCaseInfo{
			TestCaseStatusHistory: testCaseStatusHistory{
				Failed:  failed,
				Skipped: skipped,
				Passed:  passed,
			},
			TestCaseName:  testCaseName,
			TestSuiteName: testSuiteName,
			FailurRate:    failureRate,
		}
	}

	// Forming the output array where each element corresponds to a run
	for _, uniqueRun := range totalRunsDrillDownQueryResponse.Aggregations.Runs.Buckets {

		var subRows []subRowsTotalRunsDrillDown = []subRowsTotalRunsDrillDown{}
		var runStatus string

		// Using the first element directly as top_hits aggregation is set to return exactly one doc in the query
		source := uniqueRun.RunDetails.Hits.Hits[0].Source

		failedTestsCount := len(uniqueRun.FailedDocs.FailedTestCases.Buckets)

		if failedTestsCount > 0 {

			for _, failedTest := range uniqueRun.FailedDocs.FailedTestCases.Buckets {

				if testCaseInfo, ok := allFailedTestCases[failedTest.Key]; ok {
					subRows = append(subRows, subRowsTotalRunsDrillDown{
						FailedTestName: testCaseInfo.TestCaseName,
						FailureRate: failureRateTotalRunsDrillDown{
							ColorScheme: []internal.ColorScheme{
								{
									Color0: "#009C5B",
									Color1: "#62CA9D",
								},
								{
									Color0: "#D32227",
									Color1: "#FB6E72",
								},
								{
									Color0: "#F2A414",
									Color1: "#FFE6C1",
								},
							},
							LightColorScheme: []internal.ColorScheme{
								{
									Color0: "#0C9E61",
									Color1: "#79CAA8",
								},
								{
									Color0: "#E83D39",
									Color1: "#F39492",
								},
								{
									Color0: "#F2A414",
									Color1: "#FFE6C1",
								},
							},
							Type:  pb.ChartType_SINGLE_BAR.String(),
							Value: fmt.Sprintf("%.1f%%", testCaseInfo.FailurRate),
							Data: []struct {
								Title string `json:"title"`
								Value int    `json:"value"`
							}{
								{
									Title: "Successful runs",
									Value: testCaseInfo.TestCaseStatusHistory.Passed,
								},
								{
									Title: "Failed runs",
									Value: testCaseInfo.TestCaseStatusHistory.Failed,
								},
								{
									Title: "Skipped runs",
									Value: testCaseInfo.TestCaseStatusHistory.Skipped,
								},
							},
						},
						FailureRateValue: helper.TruncateFloat(testCaseInfo.FailurRate),

						ViewRunActivity: struct {
							DrillDown internal.DrillDownWithReportInfo `json:"drillDown"`
						}{
							DrillDown: internal.DrillDownWithReportInfo{
								ReportId: constants.TEST_OVERVIEW_VIEW_RUN_ACTIVITY,
								ReportInfo: pb.ReportInfo{
									AutomationId:  source.AutomationID,
									Branch:        source.BranchID,
									RunId:         source.RunID,
									RunNumber:     strconv.Itoa(source.RunNumber),
									ComponentId:   source.ComponentID,
									TestSuiteName: testCaseInfo.TestSuiteName,
									TestCaseName:  testCaseInfo.TestCaseName,
								},
								ReportTitle: "Test case activity - " + testCaseInfo.TestCaseName,
							},
						},
					})
				} else {
					return nil, fmt.Errorf("Failed to find test case information for %s", failedTest.Key)
				}

			}
		}

		// Sort the subRows slice based on the failureRateValue field
		if len(subRows) > 0 {
			sort.Slice(subRows, func(i, j int) bool {
				return subRows[i].FailureRateValue > subRows[j].FailureRateValue
			})
		}

		if utf8.RuneCountInString(source.RunStatus) > 1 {
			if source.RunStatus == "SUCCEEDED" {
				runStatus = "Success"
			} else if source.RunStatus == "FAILED" || source.RunStatus == "ABORTED" || source.RunStatus == "TIMED_OUT" {
				runStatus = "Failure"
			} else {
				runStatus = strings.ToUpper(source.RunStatus[:1]) + strings.ToLower(source.RunStatus[1:])
			}

		}

		output = append(output, totalRunsDrillDownOutput{
			RunTime:        uniqueRun.RunDetails.Hits.Hits[0].Fields.ZonedRunStartTime[0],
			AutomationID:   source.AutomationID,
			BranchID:       source.BranchID,
			Build:          fmt.Sprintf("%d", source.RunNumber),
			ComponentID:    source.ComponentID,
			OrganizationID: source.OrgID,
			RunID:          source.RunID,
			RunStatus:      runStatus,
			TotalTests:     int(uniqueRun.TotalTestCasesCount.Value),
			FailedTests:    failedTestsCount,
			SubRows:        subRows,
		})

	}

	// Sort the output slice based on the Build field
	if len(output) > 0 {
		sort.Slice(output, func(i, j int) bool {
			build1, _ := strconv.Atoi(output[i].Build)
			build2, _ := strconv.Atoi(output[j].Build)
			return build1 > build2
		})
	} else {
		return nil, db.ErrNoDataFound
	}

	var reports structpb.ListValue

	byteValue, err := json.Marshal(output)
	if log.CheckErrorf(err, "error marshaling response in TestInsightsTotalRunsDrillDown") {
		return nil, db.ErrInternalServer
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, "error unmarshaling response in TestInsightsTotalRunsDrillDown") {
		return nil, db.ErrInternalServer
	}
	return &reports, nil
}

/*
RunDetailsTestResults fetches the test results reported in a workflow run based on the view option selected
*/
func RunDetailsTestResults(replacements map[string]any, ctx context.Context, grpcClient client.GrpcClient) (*structpb.ListValue, error) {

	client := db.GetOpenSearchClient()
	if client == nil {
		log.Error("", fmt.Errorf("Failed to establish OpenSearch connection in RunDetailsTestResults()"))
		return nil, db.ErrInternalServer
	}

	if viewOption, ok := replacements["viewOption"].(string); ok {
		switch viewOption {
		case pb.ViewOptions_TEST_SUITES_VIEW.String():
			updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.RunDetailsTestResultsTestSuitesViewQuery)
			if log.CheckErrorf(err, "could not replace json placeholders for Run Details Test Results Test Suites View:", replacements) {
				return nil, db.ErrInternalServer
			}

			response, err := getSearchResponse(updatedJSON, constants.TEST_SUITE_INDEX, client)
			if log.CheckErrorf(err, "Error fetching Opensearch data for Run Details Test Results Test Suites View. Fetch failed - ") {
				return nil, db.ErrInternalServer
			}

			reports, err := transformRunDetailsTestResultsTestSuitesView(response)
			if err != nil {
				if err == db.ErrNoDataFound {
					log.Info("No data found in RunDetailsTestResults() for Test Suites view")
					return nil, db.ErrNoDataFound
				} else {
					log.Errorf(err, "Error transforming Test Suites view response in RunDetailsTestResults()")
					return nil, db.ErrInternalServer
				}
			}
			return reports, nil

		case pb.ViewOptions_TEST_CASES_VIEW.String():
			updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.RunDetailsTestResultsTestCasesViewQuery)
			if log.CheckErrorf(err, "could not replace json placeholders for Run Details Test Results Test Cases View:", replacements) {
				return nil, db.ErrInternalServer
			}

			response, err := getSearchResponse(updatedJSON, constants.TEST_CASES_INDEX, client)
			if log.CheckErrorf(err, "Error fetching Opensearch data for Run Details Test Results Test Cases View. Fetch failed - ") {
				return nil, db.ErrInternalServer
			}

			reports, err := transformRunDetailsTestResultsTestCasesView(response)
			if err != nil {
				if err == db.ErrNoDataFound {
					log.Info("No data found in RunDetailsTestResults() for Test Cases view")
					return nil, db.ErrNoDataFound
				} else {
					log.Errorf(err, "Error transforming Test Cases view response in RunDetailsTestResults()")
					return nil, db.ErrInternalServer
				}
			}
			return reports, nil

		default:
			return nil, fmt.Errorf("Invalid view option")
		}

	}

	return nil, fmt.Errorf("Invalid view option")

}

/*
transformRunDetailsTestResultsTestSuitesView transforms OpenSearch response for the "Test Suites" view in the Test Results widget in
the Runs page into the API contract format
*/
func transformRunDetailsTestResultsTestSuitesView(response string) (*structpb.ListValue, error) {

	var testSuitesViewResponse struct {
		Hits struct {
			Total struct {
				Value    int    `json:"value"`
				Relation string `json:"relation"`
			} `json:"total"`
		} `json:"hits"`
		Aggregations struct {
			TestSuitesBuckets struct {
				Buckets []struct {
					Key          string `json:"key"`
					DocCount     int    `json:"doc_count"`
					TestSuiteDoc struct {
						Hits struct {
							Hits []struct {
								Source struct {
									Duration      float64 `json:"duration"`
									Total         int     `json:"total"`
									ComponentId   string  `json:"component_id"`
									RunId         string  `json:"run_id"`
									RunNumber     int     `json:"run_number"`
									Passed        int     `json:"passed"`
									Failed        int     `json:"failed"`
									Skipped       int     `json:"skipped"`
									TestSuiteName string  `json:"test_suite_name"`
								} `json:"_source"`
							} `json:"hits"`
						} `json:"hits"`
					} `json:"test_suite_doc"`
				} `json:"buckets"`
			} `json:"test_suite_buckets"`
		} `json:"aggregations"`
	}

	err := json.Unmarshal([]byte(response), &testSuitesViewResponse)
	if log.CheckErrorf(err, "Error unmarshalling OpenSearch response into struct in transformRunDetailsTestResultsTestSuitesView()") {
		return nil, db.ErrInternalServer
	}

	if testSuitesViewResponse.Hits.Total.Value == 0 {
		return nil, db.ErrNoDataFound
	}

	type TestSuiteInfo struct {
		DrillDown        internal.DrillDownWithReportInfo `json:"drillDown"`
		RunTime          int64                            `json:"runTime"`
		TestCasesFailed  int                              `json:"testCasesFailed"`
		TestCasesPassed  int                              `json:"testCasesPassed"`
		TestCasesSkipped int                              `json:"testCasesSkipped"`
		TotalTestCases   int                              `json:"totalTestCases"`
		TestSuiteName    string                           `json:"testSuiteName"`
	}

	var output []TestSuiteInfo

	for _, value := range testSuitesViewResponse.Aggregations.TestSuitesBuckets.Buckets {

		// using the first element directly as top_hits agg. is set to return exactly one doc in the query
		source := value.TestSuiteDoc.Hits.Hits[0].Source

		output = append(output, TestSuiteInfo{
			DrillDown: internal.DrillDownWithReportInfo{
				ReportId: constants.RUN_DETAILS_TOTAL_TEST_CASES,
				ReportInfo: pb.ReportInfo{
					TestSuiteName: source.TestSuiteName,
					RunId:         source.RunId,
					RunNumber:     strconv.Itoa(source.RunNumber),
					ComponentId:   source.ComponentId,
				},
				ReportTitle: "Test cases - " + source.TestSuiteName,
				ReportType:  "status",
			},
			RunTime:          int64(source.Duration),
			TestCasesFailed:  source.Failed,
			TestCasesPassed:  source.Passed,
			TestCasesSkipped: source.Skipped,
			TestSuiteName:    source.TestSuiteName,
			TotalTestCases:   source.Total,
		})

	}

	// Sorting the output based on the number of failed, skipped and passed test cases
	sort.Slice(output, func(i, j int) bool {
		// Sort by failed (positive values first, larger values first)
		if output[i].TestCasesFailed > 0 && output[j].TestCasesFailed > 0 {
			return output[i].TestCasesFailed > output[j].TestCasesFailed
		}
		if output[i].TestCasesFailed > 0 {
			return true
		}
		if output[j].TestCasesFailed > 0 {
			return false
		}

		// Sort by skipped (positive values first, larger values first)
		if output[i].TestCasesSkipped > 0 && output[j].TestCasesSkipped > 0 {
			return output[i].TestCasesSkipped > output[j].TestCasesSkipped
		}
		if output[i].TestCasesSkipped > 0 {
			return true
		}
		if output[j].TestCasesSkipped > 0 {
			return false
		}

		// Sort by passed (positive values first, larger values first)
		return output[i].TestCasesPassed > output[j].TestCasesPassed
	})

	var reports structpb.ListValue

	byteValue, err := json.Marshal(output)
	if log.CheckErrorf(err, "error marshaling response in transformRunDetailsTestResultsTestSuitesView()") {
		return nil, db.ErrInternalServer
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, "error unmarshaling response in transformRunDetailsTestResultsTestSuitesView()") {
		return nil, db.ErrInternalServer
	}
	return &reports, nil
}

/*
transformRunDetailsTestResultsTestCasesView transforms OpenSearch response for the "Test Cases" view in the Test Results widget in
the Runs page into the API contract format
*/
func transformRunDetailsTestResultsTestCasesView(response string) (*structpb.ListValue, error) {

	var testCasesViewResponse struct {
		Hits struct {
			Total struct {
				Value    int    `json:"value"`
				Relation string `json:"relation"`
			} `json:"total"`
		} `json:"hits"`
		Aggregations struct {
			TestCaseBuckets struct {
				Buckets []struct {
					Key         string `json:"key"`
					DocCount    int    `json:"doc_count"`
					TestCaseDoc struct {
						Hits struct {
							Hits []struct {
								Source struct {
									TestCaseName  string  `json:"test_case_name"`
									Duration      float64 `json:"duration"`
									ComponentId   string  `json:"component_id"`
									RunId         string  `json:"run_id"`
									RunNumber     int     `json:"run_number"`
									TestSuiteName string  `json:"test_suite_name"`
									Status        string  `json:"status"`
									StdOut        string  `json:"std_out"`
									StdErr        string  `json:"std_err"`
									ErrorTrace    string  `json:"error_trace"`
								} `json:"_source"`
							} `json:"hits"`
						} `json:"hits"`
					} `json:"test_case_doc"`
				} `json:"buckets"`
			} `json:"test_case_buckets"`
		} `json:"aggregations"`
	}

	err := json.Unmarshal([]byte(response), &testCasesViewResponse)
	if log.CheckErrorf(err, "Error unmarshalling OpenSearch response into struct in transformRunDetailsTestResultsTestCasesView()") {
		return nil, db.ErrInternalServer
	}

	if testCasesViewResponse.Hits.Total.Value == 0 {
		return nil, db.ErrNoDataFound
	}

	type TestCaseInfo struct {
		DrillDown     internal.DrillDownWithReportInfo `json:"drillDown"`
		RunTime       int64                            `json:"runTime"`
		Status        string                           `json:"status"`
		TestSuiteName string                           `json:"testSuiteName"`
		TestCaseName  string                           `json:"testCaseName"`
		IsLogReported bool                             `json:"isLogReported"`
	}

	var output []TestCaseInfo

	for _, value := range testCasesViewResponse.Aggregations.TestCaseBuckets.Buckets {

		// using the first element directly as top_hits agg. is set to return exactly one doc in the query
		source := value.TestCaseDoc.Hits.Hits[0].Source

		var testCaseStatus string
		if utf8.RuneCountInString(source.Status) > 1 {
			testCaseStatus = strings.ToUpper(source.Status[:1]) + strings.ToLower(source.Status[1:])
		} else {
			testCaseStatus = source.Status
		}

		isLogReported := false
		if source.StdOut != "" || source.StdErr != "" || source.ErrorTrace != "" {
			isLogReported = true
		}

		output = append(output, TestCaseInfo{
			DrillDown: internal.DrillDownWithReportInfo{
				ReportId: constants.RUN_DETAILS_TEST_CASE_LOG,
				ReportInfo: pb.ReportInfo{
					TestSuiteName: source.TestSuiteName,
					TestCaseName:  source.TestCaseName,
					RunId:         source.RunId,
					RunNumber:     strconv.Itoa(source.RunNumber),
					ComponentId:   source.ComponentId,
				},
				ReportTitle: "Test case log - " + source.TestCaseName,
			},
			RunTime:       int64(source.Duration),
			Status:        testCaseStatus,
			TestSuiteName: source.TestSuiteName,
			TestCaseName:  source.TestCaseName,
			IsLogReported: isLogReported,
		})

	}

	// Sort the output slice based on the status field
	sort.Slice(output, func(i, j int) bool {
		// Custom sort order
		statusOrder := map[string]int{
			"Failed":  0,
			"Skipped": 1,
			"Passed":  2,
		}
		// Check if status is in the map, if not, assign it a higher value
		iStatus, iExists := statusOrder[output[i].Status]
		jStatus, jExists := statusOrder[output[j].Status]
		if !iExists {
			iStatus = len(statusOrder)
		}
		if !jExists {
			jStatus = len(statusOrder)
		}
		return iStatus < jStatus
	})

	var reports structpb.ListValue

	byteValue, err := json.Marshal(output)
	if log.CheckErrorf(err, "error marshaling response in transformRunDetailsTestResultsTestCasesView()") {
		return nil, db.ErrInternalServer
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, "error unmarshaling response in transformRunDetailsTestResultsTestCasesView()") {
		return nil, db.ErrInternalServer
	}
	return &reports, nil
}

/*
RunDetailsTotalTestCasesDrillDown returns test case data under a paritcular test suite in a workflow run
*/
func RunDetailsTotalTestCasesDrillDown(replacements map[string]any, ctx context.Context, grpcClient client.GrpcClient) (*structpb.ListValue, error) {

	client := db.GetOpenSearchClient()
	if client == nil {
		log.Error("", fmt.Errorf("Failed to establish OpenSearch connection in RunDetailsTotalTestCasesDrillDown()"))
		return nil, db.ErrInternalServer
	}

	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.RunDetailsTotalTestCasesDrillDownQuery)
	if log.CheckErrorf(err, "could not replace json placeholders in RunDetailsTotalTestCasesDrillDown():", replacements) {
		return nil, db.ErrInternalServer
	}

	response, err := getSearchResponse(updatedJSON, constants.TEST_CASES_INDEX, client)
	if log.CheckErrorf(err, "Error fetching Opensearch data in RunDetailsTotalTestCasesDrillDown()") {
		return nil, db.ErrInternalServer
	}

	reports, err := transformRunDetailsTotalTestCasesDrillDown(response)
	if log.CheckErrorf(err, "Error transforming OpenSearch response in RunDetailsTotalTestCasesDrillDown()") {
		if err == db.ErrNoDataFound {
			return nil, db.ErrNoDataFound
		} else {
			return nil, db.ErrInternalServer
		}
	}
	return reports, nil

}

/*
transformRunDetailsTotalTestCasesDrillDown transforms OpenSearch response for the "Total Test Cases" drill down in the Test Results widget
in the Runs page into the API contract format
*/
func transformRunDetailsTotalTestCasesDrillDown(response string) (*structpb.ListValue, error) {

	var totalTestCasesDrillDownResponse struct {
		Hits struct {
			Total struct {
				Value    int    `json:"value"`
				Relation string `json:"relation"`
			} `json:"total"`
		} `json:"hits"`
		Aggregations struct {
			TestCaseBuckets struct {
				Buckets []struct {
					Key         string `json:"key"`
					DocCount    int    `json:"doc_count"`
					TestCaseDoc struct {
						Hits struct {
							Hits []struct {
								Source struct {
									TestCaseName  string  `json:"test_case_name"`
									Duration      float64 `json:"duration"`
									ComponentId   string  `json:"component_id"`
									RunId         string  `json:"run_id"`
									RunNumber     int     `json:"run_number"`
									TestSuiteName string  `json:"test_suite_name"`
									Status        string  `json:"status"`
									StdOut        string  `json:"std_out"`
									StdErr        string  `json:"std_err"`
									ErrorTrace    string  `json:"error_trace"`
								} `json:"_source"`
							} `json:"hits"`
						} `json:"hits"`
					} `json:"test_case_doc"`
				} `json:"buckets"`
			} `json:"test_case_buckets"`
		} `json:"aggregations"`
	}

	err := json.Unmarshal([]byte(response), &totalTestCasesDrillDownResponse)
	if log.CheckErrorf(err, "Error unmarshalling OpenSearch response into struct in transformRunDetailsTotalTestCasesDrillDown()") {
		return nil, db.ErrInternalServer
	}

	if totalTestCasesDrillDownResponse.Hits.Total.Value == 0 {
		return nil, db.ErrNoDataFound
	}

	type TestCaseInfo struct {
		DrillDown     internal.DrillDownWithReportInfo `json:"drillDown"`
		RunTime       int64                            `json:"runTime"`
		Status        string                           `json:"status"`
		TestCaseName  string                           `json:"testCaseName"`
		IsLogReported bool                             `json:"isLogReported"`
	}

	var output []TestCaseInfo

	for _, value := range totalTestCasesDrillDownResponse.Aggregations.TestCaseBuckets.Buckets {

		// using the first element directly as top_hits agg. is set to return exactly one doc in the query
		source := value.TestCaseDoc.Hits.Hits[0].Source
		var testCaseStatus string
		if utf8.RuneCountInString(source.Status) > 1 {
			testCaseStatus = strings.ToUpper(source.Status[:1]) + strings.ToLower(source.Status[1:])
		} else {
			testCaseStatus = source.Status
		}

		isLogReported := false
		if source.StdOut != "" || source.StdErr != "" || source.ErrorTrace != "" {
			isLogReported = true
		}

		output = append(output, TestCaseInfo{
			DrillDown: internal.DrillDownWithReportInfo{
				ReportId: constants.RUN_DETAILS_TEST_CASE_LOG,
				ReportInfo: pb.ReportInfo{
					TestSuiteName: source.TestSuiteName,
					TestCaseName:  source.TestCaseName,
					RunId:         source.RunId,
					RunNumber:     strconv.Itoa(source.RunNumber),
					ComponentId:   source.ComponentId,
				},
				ReportTitle: "Test case log - " + source.TestCaseName,
			},
			RunTime:       int64(source.Duration),
			Status:        testCaseStatus,
			TestCaseName:  source.TestCaseName,
			IsLogReported: isLogReported,
		})

	}

	// Sort the output slice based on the status field
	sort.Slice(output, func(i, j int) bool {
		// Custom sort order
		statusOrder := map[string]int{
			"Failed":  0,
			"Skipped": 1,
			"Passed":  2,
		}
		// Check if status is in the map, if not, assign it a higher value
		iStatus, iExists := statusOrder[output[i].Status]
		jStatus, jExists := statusOrder[output[j].Status]
		if !iExists {
			iStatus = len(statusOrder)
		}
		if !jExists {
			jStatus = len(statusOrder)
		}
		return iStatus < jStatus
	})

	var reports structpb.ListValue

	byteValue, err := json.Marshal(output)
	if log.CheckErrorf(err, "error marshaling response in transformRunDetailsTotalTestCasesDrillDown()") {
		return nil, db.ErrInternalServer
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, "error unmarshaling response in transformRunDetailsTotalTestCasesDrillDown()") {
		return nil, db.ErrInternalServer
	}
	return &reports, nil
}

/*
RunDetailsTestCaseLogDrillDown returns the logs reported by a test case in a workflow run
*/
func RunDetailsTestCaseLogDrillDown(replacements map[string]any, ctx context.Context, grpcClient client.GrpcClient) (*structpb.ListValue, error) {

	client := db.GetOpenSearchClient()
	if client == nil {
		log.Error("", fmt.Errorf("Failed to establish OpenSearch connection in RunDetailsTestCaseLogDrillDown()"))
		return nil, db.ErrInternalServer
	}

	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.RunDetailsTestCaseLogDrillDownQuery)
	if log.CheckErrorf(err, "could not replace json placeholders in RunDetailsTestCaseLogDrillDown():", replacements) {
		return nil, db.ErrInternalServer
	}

	response, err := getSearchResponse(updatedJSON, constants.TEST_CASES_INDEX, client)
	if log.CheckErrorf(err, "Error fetching Opensearch data in RunDetailsTestCaseLogDrillDown(). Fetch failed - ") {
		return nil, db.ErrInternalServer
	}

	reports, err := transformRunDetailsTestCaseLogDrillDown(response)
	if log.CheckErrorf(err, "Error transforming OpenSearch response in RunDetailsTestCaseLogDrillDown()") {
		if err == db.ErrNoDataFound {
			return nil, db.ErrNoDataFound
		} else {
			return nil, db.ErrInternalServer
		}
	}
	return reports, nil

}

/*
transformRunDetailsTestCaseLogDrillDown transforms OpenSearch response for the "Test Case Log" drill down in the Test Results widget
in the Runs page into the API contract format
*/
func transformRunDetailsTestCaseLogDrillDown(response string) (*structpb.ListValue, error) {
	var testCaseLogInput struct {
		Hits struct {
			Hits []struct {
				Source struct {
					StdErr     string `json:"std_err"`
					ErrorTrace string `json:"error_trace"`
					StdOut     string `json:"std_out"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	type testCaseLogOutput struct {
		Message string `json:"message"`
	}

	err := json.Unmarshal([]byte(response), &testCaseLogInput)
	if log.CheckErrorf(err, "Error unmarshalling OpenSearch response into struct in transformRunDetailsTestCaseLogDrillDown()") {
		return nil, db.ErrInternalServer
	}

	if testCaseLogInput.Hits.Hits != nil && len(testCaseLogInput.Hits.Hits) > 0 {
		source := testCaseLogInput.Hits.Hits[0].Source
		var logMessage string

		hasStdOut := source.StdOut != ""
		hasStdErr := source.StdErr != ""
		hasErrorTrace := source.ErrorTrace != ""

		if !hasStdOut && !hasStdErr && !hasErrorTrace {
			logMessage = "This test case did not report any output."
		} else {
			if hasStdOut {
				logMessage = logMessage + source.StdOut + "\n"
			}
			if hasStdErr {
				logMessage = logMessage + source.StdErr + "\n"
			}
			if hasErrorTrace {
				logMessage = logMessage + source.ErrorTrace + "\n"
			}

		}

		output := []struct {
			LogDetails testCaseLogOutput `json:"logDetails"`
		}{
			{
				LogDetails: testCaseLogOutput{
					Message: logMessage,
				},
			},
		}

		var reports structpb.ListValue
		byteValue, err := json.Marshal(output)
		if log.CheckErrorf(err, "error marshaling response in transformRunDetailsTestCaseLogDrillDown()") {
			return nil, db.ErrInternalServer
		}
		err = protojson.Unmarshal(byteValue, &reports)
		if log.CheckErrorf(err, "error unmarshaling response in transformRunDetailsTestCaseLogDrillDown()") {
			return nil, db.ErrInternalServer
		}
		return &reports, nil
	}

	log.Error("", fmt.Errorf("Struct with OpenSearch response in transformRunDetailsTestCaseLogDrillDown is empty"))
	return nil, db.ErrNoDataFound
}

func RunDetailsTestResultsIndicators(replacements map[string]any, ctx context.Context, grpcClient client.GrpcClient) (*structpb.ListValue, error) {

	isTestInsightsDataFound := false

	var testResultsIndicatorsOutput []RunDetailsTestResultsIndicatorsResponse

	testResultsIndicatorsOutput = append(testResultsIndicatorsOutput, RunDetailsTestResultsIndicatorsResponse{})

	orgID, ok := replacements["orgId"].(string)
	if !ok {
		log.Errorf(fmt.Errorf("org ID not found"), "org ID not found in RunDetailsTestResultsIndicators()")
		return nil, fmt.Errorf("org ID not found")
	}

	component, ok := replacements["component"].([]string)
	if !ok {
		log.Errorf(fmt.Errorf("component not found"), "component ID not found in RunDetailsTestResultsIndicators()")
		return nil, fmt.Errorf("component ID not found")
	}

	// Check if test data exists for the requested component
	testInsightsCount := getDocCount(constants.TEST_SUITE_INDEX, orgID, component)
	if testInsightsCount > 0 {
		isTestInsightsDataFound = true
	}

	testResultsIndicatorsOutput[0].IsTestInsightsDataFound = isTestInsightsDataFound

	if isTestInsightsDataFound {
		client := db.GetOpenSearchClient()
		if client == nil {
			log.Error("", fmt.Errorf("Failed to establish OpenSearch connection in RunDetailsTestResultsIndicators()"))
			return nil, db.ErrInternalServer
		}

		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.RunDetailsTestResultsIndicatorsQuery)
		if log.CheckErrorf(err, "could not replace json placeholders in RunDetailsTestResultsIndicators():", replacements) {
			return nil, err
		}

		response, err := getSearchResponse(updatedJSON, constants.TEST_SUITE_INDEX, client)
		if log.CheckErrorf(err, "Error fetching Opensearch data in RunDetailsTestResultsIndicators(). ") {
			return nil, err
		}

		err = transformRunDetailsTestResultsIndicators(response, &testResultsIndicatorsOutput[0])
		if log.CheckErrorf(err, "Error transforming OpenSearch response in RunDetailsTestResultsIndicators()") {
			return nil, err
		}

	}

	var reports structpb.ListValue
	byteValue, err := json.Marshal(testResultsIndicatorsOutput)
	if log.CheckErrorf(err, "error marshaling response in RunDetailsTestResultsIndicators()") {
		return nil, err
	}
	err = protojson.Unmarshal(byteValue, &reports)
	if log.CheckErrorf(err, "error unmarshaling response in RunDetailsTestResultsIndicators()") {
		return nil, err
	}
	return &reports, nil

}

func transformRunDetailsTestResultsIndicators(queryResponse string, output *RunDetailsTestResultsIndicatorsResponse) error {

	var indicatorsQueryResponse struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
		} `json:"hits"`
		Aggregations struct {
			TestSuites struct {
				DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
				SumOtherDocCount        int `json:"sum_other_doc_count"`
				Buckets                 []struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
					Statuses struct {
						Hits struct {
							Hits []struct {
								Source struct {
									Passed  int `json:"passed"`
									Failed  int `json:"failed"`
									Skipped int `json:"skipped"`
								} `json:"_source"`
							} `json:"hits"`
						} `json:"hits"`
					} `json:"statuses"`
				} `json:"buckets"`
			} `json:"test_suites"`
		} `json:"aggregations"`
	}

	var passed, failed, skipped int

	err := json.Unmarshal([]byte(queryResponse), &indicatorsQueryResponse)
	if log.CheckErrorf(err, "Error unmarshalling OpenSearch response into struct in transformRunDetailsTestResultsIndicators()") {
		return err
	}

	if indicatorsQueryResponse.Hits.Total.Value == 0 {
		return nil
	}

	for _, value := range indicatorsQueryResponse.Aggregations.TestSuites.Buckets {
		source := value.Statuses.Hits.Hits[0].Source
		passed += source.Passed
		failed += source.Failed
		skipped += source.Skipped
	}

	output.TestCasesPassed = passed
	output.TestCasesFailed = failed
	output.TestCasesSkipped = skipped

	return nil

}
