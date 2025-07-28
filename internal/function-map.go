package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	api "github.com/calculi-corp/api/go"
	"github.com/calculi-corp/api/go/endpoint"
	scanner "github.com/calculi-corp/api/go/scanner"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/calculi-corp/api/go/vsm/report"
	cutils "github.com/calculi-corp/common/utils"
	"github.com/calculi-corp/config"
	client "github.com/calculi-corp/grpc-client"
	"github.com/calculi-corp/log"
	opensearchconfig "github.com/calculi-corp/opensearch-config"
	"github.com/calculi-corp/reports-service/cache"
	"github.com/calculi-corp/reports-service/constants"
	"github.com/calculi-corp/reports-service/db"
	"github.com/calculi-corp/reports-service/exceptions"
	helper "github.com/calculi-corp/reports-service/helper"
)

type ColorScheme struct {
	Color0 string `json:"color0"`
	Color1 string `json:"color1"`
}

type queryData struct {
	TOTAL     int64
	VERY_HIGH int64
	HIGH      int64
	MEDIUM    int64
	LOW       int64
}

type SonarHeaderData struct {
	Value string `json:"value"`
}

type CodeBaseOverview struct {
	FileName             string `json:"fileName"`
	TotalLines           int32  `json:"totalLinesOfCode"`
	CodeCoverage         string `json:"codeCoverage"`
	CoveredLines         int32  `json:"coveredLines"`
	LinesToCover         int32  `json:"linesToCover"`
	CyclomaticComplexity int32  `json:"cyclomaticComplexity"`
	CognitiveComplexity  int32  `json:"cognitiveComplexity"`
}
type IssueTypeData struct {
	CodeSmell        int32 `json:"codeSmell"`
	Bug              int32 `json:"bug"`
	Vulnerability    int32 `json:"vulnerability"`
	SecurityHotspots int32 `json:"securityHotspots"`
}
type CoverageInfoData struct {
	TotalLines     int32 `json:"totalLines"`
	TotalCodelInes int32 `json:"totalCodeLines"`
	LinesToCover   int32 `json:"linesToCover"`
	LinesCovered   int32 `json:"linesCovered"`
}
type DuplicateInfoData struct {
	TotalLines      int32 `json:"totalLines"`
	DuplicateFiles  int32 `json:"duplicateFiles"`
	DuplicateBlocks int32 `json:"duplicateBlocks"`
	DuplicateLines  int32 `json:"duplicateLines"`
}
type ComponentAndAutomationData struct {
	Value                int          `json:"value"`
	Raised               bool         `json:"raised"`
	DifferencePercentage int          `json:"differencePercentage"`
	SubTitle             SubTitleData `json:"subTitle"`
	Data                 []string     `json:"data"`
}

type ComponentAndAutomationBranchData struct {
	Value int `json:"value"`
}

type SubTitleData struct {
	Title string `json:"title"`
	Value int    `json:"value"`
}

type SubHeaderData struct {
	Title     string    `json:"title"`
	Value     int       `json:"value"`
	Drilldown DrillDown `json:"drillDown"`
}

type ChartData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type ChartInfo struct {
	Title                string    `json:"title"`
	Value                int       `json:"value"`
	Raised               bool      `json:"raised"`
	DifferencePercentage int       `json:"differencePercentage"`
	Drilldown            DrillDown `json:"drillDown"`
}

type DrillDown struct {
	ReportId    string `json:"reportId"`
	ReportTitle string `json:"reportTitle"`
	ReportType  string `json:"reportType"`
}

type DrillDownWithReportInfo struct {
	ReportId    string        `json:"reportId"`
	ReportTitle string        `json:"reportTitle"`
	ReportType  string        `json:"reportType"`
	ReportInfo  pb.ReportInfo `json:"reportInfo"`
}

type TestCasesOverview struct {
	TestSuiteName   string  `json:"testSuiteName"`
	TestCaseName    string  `json:"testCaseName"`
	ComponentName   string  `json:"componentName"`
	Workflow        string  `json:"workflow"`
	Source          string  `json:"source"`
	Branch          string  `json:"branch"`
	LastRun         string  `json:"lastRun"`
	LastRunInMillis int     `json:"lastRunInMillis"`
	AvgRunTime      float64 `json:"avgRunTime"`
	TotalRuns       struct {
		Value     int                     `json:"value"`
		DrillDown DrillDownWithReportInfo `json:"drillDown"`
	} `json:"totalRuns"`
	TotalRunsValue int `json:"totalRunsValue"`
	FailureRate    struct {
		Type        string `json:"type"`
		ColorScheme []struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		} `json:"colorScheme"`
		LightColorScheme []struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		} `json:"lightColorScheme"`
		Value string `json:"value"`
		Data  []struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		} `json:"data"`
	} `json:"failureRate"`
	FailureRateValue float32 `json:"failureRateValue"`
}

type TestSuitesOverviewResponse struct {
	Aggregations struct {
		TestSuitesOverview struct {
			Value map[string]struct {
				ComponentID           string  `json:"component_id"`
				ComponentName         string  `json:"component_name"`
				DurationInMillis      int     `json:"duration_in_millis"`
				StartTime             string  `json:"start_time"`
				AverageDuration       float64 `json:"average_duration"`
				AutomationID          string  `json:"automation_id"`
				RunId                 string  `json:"run_id"`
				Duration              int     `json:"duration"`
				Total                 int     `json:"total_cases"`
				TotalDurationInMillis float32 `json:"total_duration_in_millis"`
				BranchID              string  `json:"branch_id"`
				OrgID                 string  `json:"org_id"`
				BranchName            string  `json:"branch_name"`
				StartTimeInMillis     int64   `json:"start_time_in_millis"`
				AutomationName        string  `json:"automation_name"`
				Runs                  int     `json:"workflow_runs"`
				TestSuiteName         string  `json:"test_suite_name"`
				FailureRateForLastRun string  `json:"failure_rate_for_last_run"`
				SkippedCasesCount     float64 `json:"skipped_cases_count"`
				SuccessfulCasesCount  float64 `json:"successful_cases_count"`
				FailedCasesCount      float64 `json:"failed_cases_count"`
			} `json:"value"`
		} `json:"testSuitesOverview"`
	} `json:"aggregations"`
}

type TestSuitesOverview struct {
	TestSuiteName   string `json:"testSuiteName"`
	ComponentName   string `json:"componentName"`
	Workflow        string `json:"workflow"`
	Source          string `json:"source"`
	Branch          string `json:"branch"`
	DefaultBranch   string `json:"defaultBranch"`
	LastRun         string `json:"lastRun"`
	LastRunInMillis int    `json:"lastRunInMillis"`
	TotalTestCases  struct {
		Value     int                     `json:"value"`
		DrillDown DrillDownWithReportInfo `json:"drillDown"`
	} `json:"totalTestCases"`
	TotalTestCasesValue int     `json:"totalTestCasesValue"`
	AvgRunTime          float64 `json:"avgRunTime"`
	TotalRuns           struct {
		Value     int                     `json:"value"`
		DrillDown DrillDownWithReportInfo `json:"drillDown"`
	} `json:"totalRuns"`
	TotalRunsValue int `json:"totalRunsValue"`
	FailureRate    struct {
		ColorScheme []struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		} `json:"colorScheme"`
		Data []struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		} `json:"data"`
		LightColorScheme []struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		} `json:"lightColorScheme"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"failureRate"`
	FailureRateValue float32 `json:"failureRateValue"`
}

type TestCaseResponse struct {
	TotalExecCount        int     `json:"total_exec_count"`
	ComponentID           string  `json:"component_id"`
	ComponentName         string  `json:"component_name"`
	DurationInMillis      int     `json:"duration_in_millis"`
	StartTime             string  `json:"start_time"`
	AverageDuration       float64 `json:"average_duration"`
	AutomationID          string  `json:"automation_id"`
	TestCaseName          string  `json:"test_case_name"`
	Duration              int     `json:"duration"`
	FailureRate           string  `json:"failure_rate"`
	TotalDurationInMillis float32 `json:"total_duration_in_millis"`
	BranchID              string  `json:"branch_id"`
	OrgID                 string  `json:"org_id"`
	BranchName            string  `json:"branch_name"`
	StartTimeInMillis     int64   `json:"start_time_in_millis"`
	AutomationName        string  `json:"automation_name"`
	Runs                  int     `json:"runs"`
	FailureCount          int     `json:"failure_count"`
	SuccessCount          int     `json:"success_count"`
	SkippedCount          int     `json:"skipped_count"`
	TestSuiteName         string  `json:"test_suite_name"`
	Status                string  `json:"status"`
}

type TestCasesOverviewResponse struct {
	Aggregations struct {
		TestCasesOverview struct {
			Value map[string]TestCaseResponse
		} `json:"testCasesOverview"`
	} `json:"aggregations"`
}

type TestComponentsViewResponse struct {
	Aggregations struct {
		WorkflowBuckets struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key                string `json:"key"`
				DocCount           int    `json:"doc_count"`
				TotalTestCasesRuns struct {
					Value float64 `json:"value"`
				} `json:"total_test_cases_runs"`
				LatestDoc struct {
					Hits struct {
						Total struct {
							Value    int    `json:"value"`
							Relation string `json:"relation"`
						} `json:"total"`
						MaxScore any `json:"max_score"`
						Hits     []struct {
							Index  string `json:"_index"`
							ID     string `json:"_id"`
							Score  any    `json:"_score"`
							Source struct {
								AutomationID   string `json:"automation_id"`
								Total          int    `json:"total"`
								RunID          string `json:"run_id"`
								BranchID       string `json:"branch_id"`
								BranchName     string `json:"branch_name"`
								ComponentName  string `json:"component_name"`
								RunStartTime   string `json:"run_start_time"`
								AutomationName string `json:"automation_name"`
							} `json:"_source"`
							Fields struct {
								ZonedRunStartTime    []string `json:"zoned_run_start_time"`
								RunStartTimeInMillis []int64  `json:"run_start_time_in_millis"`
							} `json:"fields"`
							Sort []int64 `json:"sort"`
						} `json:"hits"`
					} `json:"hits"`
				} `json:"latest_doc"`
				TotalDuration struct {
					Value float64 `json:"value"`
				} `json:"total_duration"`
				SkippedCount struct {
					Value float64 `json:"value"`
				} `json:"skipped_count"`
				SuccessCount struct {
					Value float64 `json:"value"`
				} `json:"success_count"`
				FailureCount struct {
					Value float64 `json:"value"`
				} `json:"failure_count"`
				FailureRate struct {
					Value float32 `json:"value"`
				} `json:"failure_rate"`
				AvgRunTime struct {
					Value float64 `json:"value"`
				} `json:"avg_run_time"`
				TotalTestCasesCount struct {
					Value float64 `json:"value"`
				} `json:"total_test_cases_count"`
			} `json:"buckets"`
		} `json:"workflow_buckets"`
	} `json:"aggregations"`
}

type TestComponentsView struct {
	AvgRunTime    float64 `json:"avgRunTime"`
	DefaultBranch string  `json:"defaultBranch"`
	ComponentName string  `json:"componentName"`
	FailureRate   struct {
		ColorScheme []struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		} `json:"colorScheme"`
		Data []struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		} `json:"data"`
		LightColorScheme []struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		} `json:"lightColorScheme"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"failureRate"`
	FailureRateValue    float32 `json:"failureRateValue"`
	LastRun             string  `json:"lastRun"`
	LastRunInMillis     int64   `json:"lastRunInMillis"`
	TotalTestCasesValue int     `json:"totalTestCasesValue"`
	TotalTestCases      struct {
		Value     int                     `json:"value"`
		DrillDown DrillDownWithReportInfo `json:"drillDown"`
	} `json:"totalTestCases"`
	Workflow string `json:"workflow"`
	Source   string `json:"source"`
}

// Function names referred in Widget Defintion are the Keys in FunctionMap
// To add new entry, provide reference function name and its corresponding method
var FunctionMap = map[string]interface{}{
	"mockHeader1":                                 mockContent,
	"mockHeader2":                                 mockContent,
	"mockHeader3":                                 mockContent,
	"mockSection1":                                mockContent,
	"mockSection2":                                mockContent,
	"mockSection3":                                mockContent,
	"mockContent1":                                mockContent,
	"mockContent2":                                mockContent,
	"mockContent3":                                mockContent,
	"Component Widget Header":                     componentWidgetHeader,
	"Component Widget Section":                    componentWidgetSection,
	"Automation Widget Header":                    automationWidgetHeader,
	"Automation Widget Section":                   automationWidgetSection,
	"Summary Automation Widget Section":           summaryAutomationWidgetSection,
	"Security Widget Header":                      componentWidgetHeader,
	"Security Widget Section":                     securityComponentWidgetSection,
	"Security Automation Widget Header":           automationWidgetHeader,
	"Security Automation Widget Section":          securityAutomationWidgetSection,
	"Code Coverage Header":                        getCoverageDataHeader,
	"Code Coverage Section1":                      getCoverageDataSection1,
	"Code Coverage Section2":                      getCoverageDataSection2,
	"Duplication Header":                          getDuplicationHeader,
	"Duplication  Section1":                       getDuplicationDataSection1,
	"Duplication  Section2":                       getDuplicationDataSection2,
	"Issue Type Header1":                          getIssueTypeHeader1,
	"Issue Type Header2":                          getIssueTypeHeader2,
	"Issue Type Section":                          getIssueTypeSection,
	"Code Base Overview":                          getCodeBaseOverviewWidget,
	"Open Issue Section":                          getOpenIssuesSection,
	"Component Open Issues Header":                getComponentOpenIssuesHeader,
	"Insight Project Types Widget":                getInsightProjectTypes,
	"Insight System Information Widget":           getInsightSystemInformation,
	"Insight System Health Widget":                getInsightSystemHealth,
	"Insight Completed Runs Widget":               getInsightCompletedRuns,
	"Insight Projects Activity Widget":            getInsightProjectsActivity,
	"Insight CJOC Controllers":                    getInsightCjocControllers,
	"Insight Runs Overview Widget":                getInsightRunsOverview,
	"Insight Usage Patterns Widget":               getInsightUsagePatterns,
	"workflow component comparison data":          workflowWidgetComponentComparison,
	"security workflow component comparison data": workflowComponentComparisonSI,
	"security widget component comparison":        securityComponentComponentComparison,
	"Test Workflows Widget Header":                automationWidgetHeader,
	"Test Workflows Widget Section":               testInsightsWidgetSection,
	"Test Components Widget Header":               componentWidgetHeader,
	"Test Components Widget Section":              testComponentWidgetSection,
	"Tests overview":                              getTestsOverview,
	"test workflow component comparison data":     testInsightsWorkflowsComponentComparison,
}

var PageMainFunctionMap = map[string]interface{}{
	"e8":  fetchAllBuildedComponents,
	"e9":  FetchAllDeployedEnvironments,
	"cs3": FetchAllDeployedEnvironments,
}

var (
	openSearchClient        = opensearchconfig.GetOpensearchConnection
	searchResponse          = db.GetOpensearchData
	countResponse           = db.GetOpensearchCount
	multiSearchResponse     = helper.GetMultiQueryResponse
	getOrganisationServices = helper.GetOrganisationServices
	GetAutomationMap        = getAutomationMap
	getAutomationsForBranch = GetAutomationsForBranch
	getBranchNameForId      = GetBranchNameForId
)

func getOpenIssuesSection(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	client, err := openSearchClient()
	if log.CheckErrorf(err, "Error establishing connection with OpenSearch in getOpenIssuesSection()") {
		return nil, err
	}
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.OpenIssuesSectionQuery)
	if log.CheckErrorf(err, "could not replace json placeholders in getOpenIssuesSection()", replacements) {
		return nil, err
	}

	branch, ok := replacements[constants.BRANCH]
	var modifiedJson string

	if ok && branch != nil && len(branch.(string)) > 0 {
		modifiedJson = UpdateFiltersForDrilldown(updatedJSON, replacements, true, false)

	} else {
		modifiedJson = UpdateFilters(updatedJSON, replacements)

	}

	response1, err := searchResponse(modifiedJson, constants.SECURITY_INDEX, client)
	if log.CheckErrorf(err, "Error fetching response from OpenSearch in getOpenIssuesSection()") {
		return nil, err
	}

	if ok := helper.HasNoHits(response1); ok {
		return nil, nil
	}

	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(response1), &data)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getOpenIssuesSection()") {
		return nil, err
	}
	if data[constants.AGGREGATION] != nil {
		x := data[constants.AGGREGATION].(map[string]interface{})
		if x[constants.DRILLDOWNS] != nil {
			y := x[constants.DRILLDOWNS].(map[string]interface{})
			if y[constants.VALUE] != nil {
				values := y[constants.VALUE].([]interface{})
				responseJson, err := json.Marshal(values)
				if err != nil {
					return nil, err
				} else {
					return responseJson, nil
				}
			}
		}
	}

	return nil, nil
}

func getComponentOpenIssuesHeader(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	client, err := openSearchClient()
	if log.CheckErrorf(err, "Error getting OpenSearch client in getComponentOpenIssuesHeader") {
		return nil, err
	}
	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.ComponentOpenIssueHeaderQuery)
	if log.CheckErrorf(err, "could not replace json placeholders in getComponentOpenIssuesHeader()", replacements) {
		return nil, err
	}

	branch, ok := replacements[constants.BRANCH]
	var modifiedJson string

	if ok && branch != nil && len(branch.(string)) > 0 {
		modifiedJson = UpdateFiltersForDrilldown(updatedJSON, replacements, true, false)

	} else {
		modifiedJson = UpdateFilters(updatedJSON, replacements)
	}

	response, err := searchResponse(modifiedJson, constants.SECURITY_INDEX, client)
	if log.CheckErrorf(err, "Error fetching response from OpenSearch in getComponentOpenIssuesHeader()") {
		return nil, err
	}

	if ok := helper.HasNoHits(response); ok {
		return nil, nil
	}

	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(response), &data)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getComponentOpenIssuesHeader()") {
		return nil, err
	}
	if data[constants.AGGREGATION] != nil {
		x := data[constants.AGGREGATION].(map[string]interface{})
		if x[constants.SECURITY_ISSUE_COUNT] != nil {
			y := x[constants.SECURITY_ISSUE_COUNT].(map[string]interface{})
			if y[constants.VALUE] != nil {
				jsonStr, err := json.Marshal(y[constants.VALUE])
				if err != nil {
					return nil, err
				}
				var qd queryData
				err = json.Unmarshal(jsonStr, &qd)
				if err != nil {
					return nil, err
				}

				responseData := map[string]interface{}{
					"TOTAL":     map[string]interface{}{"value": qd.TOTAL},
					"VERY_HIGH": map[string]interface{}{"value": qd.VERY_HIGH},
					"HIGH":      map[string]interface{}{"value": qd.HIGH},
					"MEDIUM":    map[string]interface{}{"value": qd.MEDIUM},
					"LOW":       map[string]interface{}{"value": qd.LOW},
				}
				responseJson, err := json.Marshal(responseData)
				if err != nil {
					return nil, err
				} else {
					return responseJson, nil
				}
			}
		}
	}
	return nil, nil
}

func getCodeBaseOverviewWidget(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	sonarArray, err := getSonarQubeResponse(replacements, ctx)
	log.CheckErrorf(err, exceptions.ErrSonarGetData)
	if len(sonarArray) > 0 {
		currentScan := &sonarArray[0]
		codeBaseOverviewArr := getCodeBaseOverview(currentScan)
		responseJson, err := json.Marshal(codeBaseOverviewArr)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	}
	return nil, nil

}

func getIssueTypeHeader1(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	sonarArray, err := getSonarQubeResponse(replacements, ctx)

	log.CheckErrorf(err, exceptions.ErrSonarGetData)
	if len(sonarArray) > 0 {
		currentScan := &sonarArray[0]

		issueTypeMap := getIssueTypeCount(currentScan)

		outputResponse := make(map[string]interface{})
		outputResponse[constants.SUB_HEADER] = []SubHeaderData{
			{
				Title: constants.CODE_SMELL,
				Value: int(issueTypeMap[constants.CODE_SMELL]),
			}, {
				Title: constants.BUG,
				Value: int(issueTypeMap[constants.BUG]),
			},
		}
		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}

	}
	return nil, nil
}

func getIssueTypeHeader2(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	sonarArray, err := getSonarQubeResponse(replacements, ctx)

	log.CheckErrorf(err, exceptions.ErrSonarGetData)
	if len(sonarArray) > 0 {
		currentScan := &sonarArray[0]

		issueTypeMap := getIssueTypeCount(currentScan)

		outputResponse := make(map[string]interface{})
		outputResponse[constants.SUB_HEADER] = []SubHeaderData{
			{
				Title: constants.VULNERABILITIES,
				Value: int(issueTypeMap[constants.VULNERABILITY]),
			}, {
				Title: constants.SECURITY_HOTSPOTS,
				Value: int(issueTypeMap[constants.SECURITY_HOTSPOTS]),
			},
		}
		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}

	}
	return nil, nil
}

func getIssueTypeSection(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {

	sonarArray, err := getSonarQubeResponse(replacements, ctx)
	log.CheckErrorf(err, exceptions.ErrSonarGetData)
	if len(sonarArray) > 0 {
		currentScan := &sonarArray[0]
		currentScanMap := getIssueTypeCount(currentScan)
		previousScanMap := make(map[string]int32)

		if len(sonarArray) > 1 {
			previousScan := &sonarArray[1]
			previousScanMap = getIssueTypeCount(previousScan)
		} else {
			previousScanMap = getIssueTypeCount(nil)

		}
		outputResponseList := make([]map[string]interface{}, 0)

		for key, value := range currentScanMap {
			outputResponse := make(map[string]interface{})
			dataMapList := make([]map[string]interface{}, 0)
			dataMap := make(map[string]interface{})
			dataMap["x"] = constants.CURRENT_SCAN
			dataMap["y"] = value
			dataMapList = append(dataMapList, dataMap)
			previousDataMap := make(map[string]interface{})

			previousDataMap["x"] = constants.PREVIOUS_SCAN
			previousDataMap["y"] = previousScanMap[key]
			dataMapList = append(dataMapList, previousDataMap)

			outputResponse[constants.DATA] = dataMapList

			outputResponse[constants.ID] = key
			outputResponseList = append(outputResponseList, outputResponse)

		}
		responseJson, err := json.Marshal(outputResponseList)
		if err != nil {
			return nil, err
		}
		return responseJson, nil

	}

	return nil, nil
}

func getCodeBaseOverview(scan *scanner.Report) []CodeBaseOverview {
	codeBaseOverviewArr := make([]CodeBaseOverview, 0)
	if scan != nil && scan.Files != nil && len(scan.Files) > 0 {
		for _, file := range scan.Files {
			if file != nil {
				codeBaseOverview := CodeBaseOverview{}
				codeBaseOverview.FileName = file.GetFile()
				if file.GetCoverage() != nil {
					codeBaseOverview.CoveredLines = file.GetCoverage().GetCoveredLines()
					codeBaseOverview.LinesToCover = file.GetCoverage().GetCoveredLines()
					codeBaseOverview.CodeCoverage = fmt.Sprintf("%d%s", int(file.GetCoverage().GetCoveragePct()), "%")
				}
				if file.GetComplexity() != nil {
					codeBaseOverview.CognitiveComplexity = file.GetComplexity().GetCognitive()
					codeBaseOverview.CyclomaticComplexity = file.GetComplexity().GetCyclomatic()
				}

				if file.GetSize() != nil {
					codeBaseOverview.TotalLines = file.GetSize().GetCodeLines()

				}
				codeBaseOverviewArr = append(codeBaseOverviewArr, codeBaseOverview)

			}

		}
	}
	return codeBaseOverviewArr
}

func IsSonarWidgetsApplicable(replacements map[string]any, ctx context.Context, securityIndex, rawIndex string) bool {
	sonarFound, _ := isSonarOrSecurityDataFound(replacements, ctx, false, securityIndex, rawIndex)
	return sonarFound
}

func IsSecurityWidgetsApplicable(replacements map[string]any, ctx context.Context, securityIndex, rawIndex string) bool {
	sonarFound, _ := isSonarOrSecurityDataFound(replacements, ctx, true, securityIndex, rawIndex)
	return sonarFound
}

func getIssueTypeCount(scan *scanner.Report) map[string]int32 {
	issueTypeMap := make(map[string]int32)
	issueTypeMap[constants.CODE_SMELL] = 0
	issueTypeMap[constants.BUG] = 0
	issueTypeMap[constants.VULNERABILITY] = 0
	issueTypeMap[constants.SECURITY_HOTSPOTS] = 0

	if scan != nil && scan.Files != nil && len(scan.Files) > 0 {
		for _, file := range scan.Files {
			if file.Issues != nil && len(file.Issues) > 0 {
				for _, issue := range file.Issues {
					if issue != nil {

						if strings.HasPrefix(issue.Code, constants.CODE_SMELL_PREFIX) {
							issueTypeMap[constants.CODE_SMELL] = (issueTypeMap[constants.CODE_SMELL] + 1)
						}

						if strings.HasPrefix(issue.Code, constants.BUG_PREFIX) {
							issueTypeMap[constants.BUG] = (issueTypeMap[constants.BUG] + 1)
						}
						if strings.HasPrefix(issue.Code, constants.SECURITY_HOTSPOTS_PREFIX) {
							issueTypeMap[constants.SECURITY_HOTSPOTS] = (issueTypeMap[constants.SECURITY_HOTSPOTS] + 1)
						}
						if strings.HasPrefix(issue.Code, constants.VULNERABILITY_PREFIX) {
							issueTypeMap[constants.VULNERABILITY] = (issueTypeMap[constants.VULNERABILITY] + 1)
						}

					}
				}
			}
		}

	}

	return issueTypeMap
}

func getCoverageDataHeader(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	sonarArray, err := getLatestSonarQubeResponse(replacements)

	log.CheckErrorf(err, exceptions.ErrSonarGetData)
	if len(sonarArray) > 0 {
		currentScan := &sonarArray[0]
		if currentScan.Coverage != nil && currentScan.Size != nil {

			codeCoverage := getCoveragePercetange(currentScan)
			codeCoverageStr := fmt.Sprintf("%d%s", codeCoverage, "%")

			coverageData := &SonarHeaderData{
				Value: codeCoverageStr,
			}

			response, err := json.Marshal(coverageData)
			if err != nil {
				return nil, err
			}
			return response, nil

		}
	}
	return nil, nil

}

func getCoverageDataSection1(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	sonarArray, err := getLatestSonarQubeResponse(replacements)

	log.CheckErrorf(err, exceptions.ErrSonarGetData)
	if len(sonarArray) > 0 {
		currentScan := &sonarArray[0]
		coverageMap := make(map[string]int)
		if currentScan != nil && currentScan.Coverage != nil && currentScan.Size != nil {
			codeCoverage := getCoveragePercetange(currentScan)
			coverageMap["Current scan"] = codeCoverage

		}

		if len(sonarArray) > 1 {
			previousScan := &sonarArray[1]

			codeCoverage := getCoveragePercetange(previousScan)
			coverageMap["Previous scan"] = codeCoverage
		}

		outputResponse := make(map[string]interface{})
		dataMapList := make([]map[string]interface{}, 0)
		for key, value := range coverageMap {
			dataMap := make(map[string]interface{})
			dataMap["x"] = key
			dataMap["y"] = value
			dataMapList = append(dataMapList, dataMap)
		}
		outputResponse[constants.DATA] = dataMapList
		outputResponse[constants.ID] = "branch"
		responseJson, err := json.Marshal([]any{outputResponse})
		if err != nil {
			return nil, err
		}
		return responseJson, nil
	}
	return nil, nil

}

func getCoverageDataSection2(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	sonarArray, err := getLatestSonarQubeResponse(replacements)

	log.CheckErrorf(err, exceptions.ErrSonarGetData)
	if len(sonarArray) > 0 {
		currentScan := &sonarArray[0]
		if currentScan != nil && currentScan.Coverage != nil && currentScan.Size != nil {
			codeCoverageInfo := getCoverageInfo(currentScan)

			outputResponse := []SubTitleData{
				{
					Title: constants.TOTAL_LINES,
					Value: int(codeCoverageInfo.TotalLines),
				}, {
					Title: constants.TOTAL_CODE_LINES,
					Value: int(codeCoverageInfo.TotalCodelInes),
				}, {
					Title: constants.LINES_COVERED,
					Value: int(codeCoverageInfo.LinesCovered),
				}, {
					Title: constants.LINES_TO_COVER,
					Value: int(codeCoverageInfo.LinesToCover),
				},
			}

			responseJson, err := json.Marshal(outputResponse)

			if err != nil {
				return nil, err
			}
			return responseJson, nil

		}
	}
	return nil, nil

}

func getDuplicationDataSection2(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	sonarArray, err := getSonarQubeResponse(replacements, ctx)

	log.CheckErrorf(err, exceptions.ErrSonarGetData)
	if len(sonarArray) > 0 {
		currentScan := &sonarArray[0]
		if currentScan != nil && currentScan.Coverage != nil && currentScan.Size != nil {
			codeCoverageInfo := getDuplicateInfo(currentScan)

			outputResponse := []SubTitleData{
				{
					Title: constants.DUPLICATE_FILES,
					Value: int(codeCoverageInfo.DuplicateFiles),
				}, {
					Title: constants.DUPLICATE_BLOCKS,
					Value: int(codeCoverageInfo.DuplicateBlocks),
				}, {
					Title: constants.TOTAL_LINES,
					Value: int(codeCoverageInfo.TotalLines),
				}, {
					Title: constants.DUPLICATE_LINES,
					Value: int(codeCoverageInfo.DuplicateLines),
				},
			}

			responseJson, err := json.Marshal(outputResponse)

			if err != nil {
				return nil, err
			}
			return responseJson, nil

		}
	}
	return nil, nil

}

func getDuplicationHeader(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	sonarArray, err := getSonarQubeResponse(replacements, ctx)

	log.CheckErrorf(err, exceptions.ErrSonarGetData)
	if len(sonarArray) > 0 {
		currentScan := &sonarArray[0]
		if currentScan != nil && currentScan.Duplication != nil && currentScan.Size != nil {

			duplicationDensity := getDuplicationDensity(currentScan)
			duplicationDensityStr := fmt.Sprintf("%d%s", int(duplicationDensity), "%")

			duplicationData := &SonarHeaderData{
				Value: duplicationDensityStr,
			}

			response, err := json.Marshal(duplicationData)
			if err != nil {
				return nil, err
			}
			return response, nil

		}
	}
	return nil, nil

}

func getDuplicationDataSection1(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	sonarArray, err := getSonarQubeResponse(replacements, ctx)

	log.CheckErrorf(err, exceptions.ErrSonarGetData)
	if len(sonarArray) > 0 {
		currentScan := &sonarArray[0]
		coverageMap := make(map[string]any)
		if currentScan != nil && currentScan.Coverage != nil && currentScan.Size != nil {
			codeCoverage := getDuplicationDensity(currentScan)
			coverageMap[constants.CURRENT_SCAN] = int(codeCoverage)

		}

		if len(sonarArray) > 1 {
			previousScan := &sonarArray[1]

			codeCoverage := getDuplicationDensity(previousScan)
			coverageMap[constants.PREVIOUS_SCAN] = int(codeCoverage)
		}

		outputResponse := make(map[string]interface{})
		dataMapList := make([]map[string]interface{}, 0)
		for key, value := range coverageMap {
			dataMap := make(map[string]interface{})
			dataMap["x"] = key
			dataMap["y"] = value
			dataMapList = append(dataMapList, dataMap)
		}
		outputResponse[constants.DATA] = dataMapList
		outputResponse[constants.ID] = "Duplicated line density"
		responseJson, err := json.Marshal([]any{outputResponse})
		if err != nil {
			return nil, err
		}
		return responseJson, nil
	}
	return nil, nil

}

func getCoveragePercetange(scan *scanner.Report) int {
	if scan != nil && scan.Coverage != nil {

		codeCoverage := scan.Coverage.GetCoveragePct()
		return int(codeCoverage)
	}
	return 0

}
func getDuplicationDensity(scan *scanner.Report) float64 {
	if scan != nil && scan.Duplication != nil && scan.Size != nil {
		duplicationDensity := scan.Duplication.GetLineDensity()
		return math.Floor(float64(duplicationDensity)*100) / 100
	}
	return 0

}
func getCoverageInfo(scan *scanner.Report) CoverageInfoData {
	coverageInfo := CoverageInfoData{}
	if scan != nil && scan.Coverage != nil && scan.Size != nil {
		coverageInfo.LinesToCover = scan.Coverage.GetLinesToCover()
		coverageInfo.TotalCodelInes = scan.Size.GetCodeLines()
		coverageInfo.TotalLines = scan.Size.GetLines()
		coverageInfo.LinesCovered = scan.Coverage.GetCoveredLines()

	}
	return coverageInfo

}

func getDuplicateInfo(scan *scanner.Report) DuplicateInfoData {
	duplicateInfo := DuplicateInfoData{}
	if scan != nil && scan.Duplication != nil && scan.Size != nil {
		duplicateInfo.DuplicateLines = scan.Duplication.GetLines()
		duplicateInfo.DuplicateFiles = scan.Duplication.GetFiles()
		duplicateInfo.TotalLines = scan.Size.GetLines()
		duplicateInfo.DuplicateBlocks = scan.Duplication.GetBlocks()

	}
	return duplicateInfo

}

func getLatestSonarQubeResponse(replacements map[string]any) ([]scanner.Report, error) {
	client, err := openSearchClient()
	var sonarArray []scanner.Report
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.Latest_Sonar_Query)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := UpdateFiltersForSonar(updatedJSON, replacements)

	response, err := searchResponse(modifiedJson, constants.RAW_SCAN_RESULTS_INDEX, client)
	if log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure) {
		return nil, err
	}

	var result map[string]interface{}
	err1 := json.Unmarshal([]byte(response), &result)
	if err1 == nil {
		if _, ok := result["hits"]; ok {
			if len(result["hits"].(map[string]interface{})["hits"].([]interface{})) > 0 {
				for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
					if _, ok := hit.(map[string]interface{})["_source"]; ok {
						sonar := hit.(map[string]interface{})["_source"].(map[string]interface{})
						if sonar != nil {
							if _, ok := sonar["action_raw_result"]; ok {
								resultByte, err := json.Marshal(sonar["action_raw_result"])
								if err == nil {
									result := scanner.Report{}
									json.Unmarshal(resultByte, &result)
									sonarArray = append(sonarArray, result)
								} else {
									log.Errorf(err, "Error in conversion of raw result : %s - ", err.Error())
								}

							} else {
								log.Debugf("no action raw_result object")
							}
						} else {
							log.Debugf(exceptions.DebugNilSonarObject)
						}

					}
				}
			}
		}
	}

	return sonarArray, nil

}

func getSonarQubeResponse(replacements map[string]any, ctx context.Context) ([]scanner.Report, error) {
	client, err := openSearchClient()
	var sonarArray []scanner.Report
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.Sonar_Query)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, err
	}
	modifiedJson := UpdateFiltersForSonar(updatedJSON, replacements)

	response, err := searchResponse(modifiedJson, constants.RAW_SCAN_RESULTS_INDEX, client)
	if log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure) {
		return nil, err
	}

	var result map[string]interface{}
	err1 := json.Unmarshal([]byte(response), &result)
	if err1 == nil {
		if _, ok := result["hits"]; ok {
			if len(result["hits"].(map[string]interface{})["hits"].([]interface{})) > 0 {
				for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
					if _, ok := hit.(map[string]interface{})["_source"]; ok {
						sonar := hit.(map[string]interface{})["_source"].(map[string]interface{})
						if sonar != nil {
							if _, ok := sonar["action_raw_result"]; ok {
								resultByte, err := json.Marshal(sonar["action_raw_result"])
								if err == nil {
									result := scanner.Report{}
									json.Unmarshal(resultByte, &result)
									sonarArray = append(sonarArray, result)
								} else {
									log.Errorf(err, "Error in conversion of raw result : %s - ", err.Error())
								}

							} else {
								log.Debugf("no action raw_result object")
							}
						} else {
							log.Debugf(exceptions.DebugNilSonarObject)
						}

					}
				}
			}
		}
	}

	return sonarArray, nil

}

func isSonarOrSecurityDataFound(replacements map[string]any, ctx context.Context, isSecurityRequest bool, securityIndex string, rawIndex string) (bool, error) {
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)
	var indexname string
	if isSecurityRequest {
		indexname = securityIndex
	} else {
		indexname = rawIndex
	}

	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.Sonar_Base_Query)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return false, err
	}

	response, err := searchResponse(updatedJSON, indexname, client)

	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	var result map[string]interface{}
	err1 := json.Unmarshal([]byte(response), &result)
	if err1 == nil {
		if _, ok := result["hits"]; ok {
			if len(result["hits"].(map[string]interface{})["hits"].([]interface{})) > 0 {
				for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
					if _, ok := hit.(map[string]interface{})["_source"]; ok {
						sonar := hit.(map[string]interface{})["_source"].(map[string]interface{})
						if sonar != nil {
							return true, nil
						} else {
							log.Debugf(exceptions.DebugNilSonarObject)
						}

					}
				}
			}
		}
	}

	return false, nil

}
func FetchAllDeployedEnvironments(widgetId string, replacements map[string]any, ctx context.Context) ([]string, map[string]any, error) {
	var environmentArray []string
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.FetchAllDeployEnv)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, nil, err
	}
	modifiedJson := UpdateFilters(updatedJSON, replacements)

	response, err := searchResponse(modifiedJson, constants.DEPLOY_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)

	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.ENVIRONMENTS] != nil {
			environments := aggsResult[constants.ENVIRONMENTS].(map[string]interface{})
			if environments[constants.VALUE] != nil {
				values := environments[constants.VALUE].([]interface{})
				for _, value := range values {
					environmentArray = append(environmentArray, fmt.Sprint(value))
				}
			}
		}
	}
	return environmentArray, nil, err
}

func fetchAllBuildedComponents(widgetId string, replacements map[string]any, ctx context.Context) ([]string, map[string]any, error) {
	var componentsArray []string
	client, err := openSearchClient()
	log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

	updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.FetchAllBuildComponents)
	if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
		return nil, nil, err
	}
	modifiedJson := UpdateFilters(updatedJSON, replacements)

	response, err := searchResponse(modifiedJson, constants.BUILD_DATA_INDEX, client)
	log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
	result := make(map[string]interface{})
	dataMap := make(map[string]any)
	json.Unmarshal([]byte(response), &result)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.COMPONENTS] != nil {
			components := aggsResult[constants.COMPONENTS].(map[string]interface{})
			if components[constants.VALUE] != nil {
				value := components[constants.VALUE].(map[string]interface{})
				if value != nil {
					componentList := value[constants.COMPONENTS].([]interface{})
					if len(componentList) > 0 {
						for _, value := range componentList {
							componentsArray = append(componentsArray, fmt.Sprint(value))
						}
					}
					min, ok := value["min"]
					if ok {
						data, err := strconv.ParseInt(min.(string), 10, 64)
						if err == nil {
							dataMap["min"] = data
						}
					}
					max, ok := value["max"]
					if ok {
						data, err := strconv.ParseInt(max.(string), 10, 64)
						if err == nil {
							dataMap["max"] = data
						}
					}
				}
			}
		}
	}
	return componentsArray, dataMap, err
}

// Widget Builder excuete the function to get data for the widget
func ExecuteFunction(k string, widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	// For any other function signature use the switch case
	switch k {
	default:
		replacements["FunctionName"] = k // TBD Do this only for mock
		f, ok := FunctionMap[k]
		if !ok || f == nil {
			return nil, fmt.Errorf("function not found for %s:%s", k, widgetId)
		}
		res, err := f.(func(string, map[string]any, context.Context, client.GrpcClient, endpoint.EndpointServiceClient) (json.RawMessage, error))(widgetId, replacements, ctx, clt, epClt)
		return res, err
	}
}

// Widget Builder excuete the function to get data for the widget
func ExecuteMultiPageBaseFunction(k string, widgetId string, replacements map[string]any, ctx context.Context) ([]string, map[string]any, error) {
	// For any other function signature use the switch case
	switch k {
	default:
		f, ok := PageMainFunctionMap[widgetId]
		if !ok || f == nil {
			return nil, nil, fmt.Errorf("function not found for %s:%s", k, widgetId)
		}
		res, dataMap, err := f.(func(string, map[string]any, context.Context) ([]string, map[string]any, error))(widgetId, replacements, ctx)
		return res, dataMap, err
	}
}

// Mock functions to get mock data.
func mockContent(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {

	var md map[string]interface{}

	fileName, ok := db.WidgetDefinitionMap[widgetId]
	if !ok {
		log.Errorf(errors.New("widget definition not found"), "widget definition not found for Id : ", widgetId)
		return nil, db.ErrInternalServer
	}

	// open the JSON file
	jsonFile, err := os.Open(config.Config.GetString("report.definition.filepath") + fileName)
	if log.CheckErrorf(err, "error opening widget definition json file ") {
		return nil, db.ErrFileNotFound
	}
	defer jsonFile.Close()

	// read from the JSON file
	byteValue, err := io.ReadAll(jsonFile)
	if log.CheckErrorf(err, "error reading widget definition json file: ") {
		return nil, db.ErrFileNotFound
	}
	json.Unmarshal([]byte(byteValue), &md)
	myMap := md["widget"].(map[string]interface{})

	name := replacements["FunctionName"].(string)
	jsonBytes, err := json.MarshalIndent(myMap[string(name)], "", "   ")
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

func componentWidgetHeader(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	log.Debugf("Component widget header params ", widgetId, replacements)
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
	headerResponse, err := getComponentWidgetData(ctx, clt, orgId, components)
	if err != nil {
		return nil, err
	}
	return headerResponse, nil
}

func componentWidgetSection(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	log.Debugf("Component widget section params ", widgetId, replacements)
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
	componentMap := map[string]struct{}{}
	serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
	if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
		return nil, err
	}
	if serviceResponse != nil {
		for i := 0; i < len(serviceResponse.GetService()); i++ {
			service := serviceResponse.GetService()[i]
			if components != nil {
				for _, component := range components {
					if service.Id == component {
						componentMap[service.Id] = struct{}{}
					}
				}
			} else {
				componentMap[service.Id] = struct{}{}
			}
		}
	}

	componentCount := float64(len(componentMap))
	if componentCount > 0 {
		client := db.GetOpenSearchClient()
		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.ComponentFilterQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson := UpdateFilters(updatedJSON, replacements)

		response, err := searchResponse(modifiedJson, constants.AUTOMATION_METADATA_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		activeCount := 0.0
		inactiveCount := 0.0
		result := make(map[string]interface{})
		actives := []string{}
		json.Unmarshal([]byte(response), &result)
		if result[constants.AGGREGATION] != nil {
			aggsResult := result[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.DISTINCT_COMPONENT] != nil {
				distinctComponent := aggsResult[constants.DISTINCT_COMPONENT].(map[string]interface{})
				if distinctComponent[constants.VALUE] != nil {
					values := distinctComponent[constants.VALUE].([]interface{})
					active := make([]string, len(values))
					for i, v := range values {
						_, ok := componentMap[fmt.Sprint(v)]
						if ok {
							active[i] = fmt.Sprint(v)
							activeCount++
						}
					}
					actives = append(actives, active...)
				}
			}
		}
		log.Debugf("Active components : %v", actives)
		inactiveCount = componentCount - activeCount
		outputResponse := make(map[string]interface{})
		drilldown := DrillDown{
			ReportId:    constants.COMPONENT,
			ReportTitle: constants.COMPONENTS_DRILLDOWN,
			ReportType:  constants.STATUS,
		}
		outputResponse[constants.DATA] = []ChartData{
			{
				Name:  constants.ACTIVE,
				Value: int(math.Round((activeCount / componentCount) * 100)),
			}, {
				Name:  constants.INACTIVE,
				Value: int(math.Round((inactiveCount / componentCount) * 100)),
			},
		}
		outputResponse[constants.INFO] = []ChartInfo{
			{
				Title:     constants.ACTIVE,
				Value:     int(activeCount),
				Drilldown: drilldown,
			}, {
				Title:     constants.INACTIVE,
				Value:     int(inactiveCount),
				Drilldown: drilldown,
			},
		}
		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Debugf(exceptions.DebugEmptyComponentData)
	}
	return nil, nil
}

func securityComponentWidgetSection(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {

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
	componentMap := map[string]struct{}{}
	serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
	if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
		return nil, err
	}
	if serviceResponse != nil {
		for i := 0; i < len(serviceResponse.GetService()); i++ {
			service := serviceResponse.GetService()[i]
			if components != nil {
				for _, component := range components {
					if service.Id == component {
						componentMap[service.Id] = struct{}{}
					}
				}
			} else {
				componentMap[service.Id] = struct{}{}
			}
		}
	}
	componentCount := float64(len(componentMap))

	if componentCount > 0 {

		rawScanQuery, err := db.ReplaceJSONplaceholders(replacements, constants.SecurityComponentFilterQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		rawScanQueryUpdated := UpdateFilters(rawScanQuery, replacements)

		scanQuery, err := db.ReplaceJSONplaceholders(replacements, constants.SecurityComponentFilterQueryByScanTime)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		scanQueryUpdated := UpdateFilters(scanQuery, replacements)

		queryMap := make(map[string]db.DbQuery)
		queryMap["scan"] = db.DbQuery{AliasName: constants.SECURITY_INDEX, QueryString: scanQueryUpdated}
		queryMap["rawScan"] = db.DbQuery{AliasName: constants.RAW_SCAN_RESULTS_INDEX, QueryString: rawScanQueryUpdated}

		responseMap, err := multiSearchResponse(queryMap)
		if log.CheckErrorf(err, "multi search query failed in securityComponentWidgetSection()") {
			return nil, err
		}

		type multiSearchResponse struct {
			Aggregations struct {
				DistinctComponent struct {
					Value []string `json:"value"`
				} `json:"distinct_component"`
			} `json:"aggregations"`
		}

		scanResult := multiSearchResponse{}
		err = json.Unmarshal(responseMap["scan"], &scanResult)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallRespGetCommitTrends) {
			return nil, err
		}

		rawScanResult := multiSearchResponse{}
		err = json.Unmarshal(responseMap["rawScan"], &rawScanResult)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallRespGetCommitTrends) {
			return nil, err
		}

		//Remove duplicates from each query response
		check := make(map[string]int)
		d := append(scanResult.Aggregations.DistinctComponent.Value, rawScanResult.Aggregations.DistinctComponent.Value...)
		combinedResult := make([]string, 0)
		for _, val := range d {
			check[val] = 1
		}

		for k, _ := range check {
			combinedResult = append(combinedResult, k)
		}

		activeCount := float64(len(combinedResult))
		inactiveCount := componentCount - activeCount
		if activeCount > componentCount {
			inactiveCount = componentCount
		}

		outputResponse := make(map[string]interface{})
		drilldown := DrillDown{
			ReportId:    constants.SECURITY_COMPONENT,
			ReportTitle: constants.COMPONENTS_DRILLDOWN,
			ReportType:  constants.SCANNERS_LOWERCASE,
		}
		outputResponse[constants.DATA] = []ChartData{
			{
				Name:  constants.WITHSCANNERS,
				Value: int(math.Round((activeCount / componentCount) * 100)),
			}, {
				Name:  constants.WITHOUTSCANNERS,
				Value: int(math.Round((inactiveCount / componentCount) * 100)),
			},
		}
		outputResponse[constants.INFO] = []ChartInfo{
			{
				Title:     constants.WITHSCANNERS,
				Value:     int(activeCount),
				Drilldown: drilldown,
			}, {
				Title:     constants.WITHOUTSCANNERS,
				Value:     int(inactiveCount),
				Drilldown: drilldown,
			},
		}
		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Debugf(exceptions.DebugEmptyComponentData)
	}
	return nil, nil
}

func securityComponentComponentComparison(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {

	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, exceptions.ErrOrgIdNotSpecified)
	}

	components, ok := replacements[constants.COMPONENT].([]string)
	if !ok || len(components) == 0 || (len(components) == 1 && components[0] == constants.ALL) {
		components = nil
	}
	componentMap := map[string]struct{}{}
	serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
	if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
		return nil, err
	}
	if serviceResponse != nil {
		for i := 0; i < len(serviceResponse.GetService()); i++ {
			service := serviceResponse.GetService()[i]
			if components != nil {
				for _, component := range components {
					if service.Id == component {
						componentMap[service.Id] = struct{}{}
					}
				}
			} else {
				componentMap[service.Id] = struct{}{}
			}
		}
	}
	componentCount := float64(len(componentMap))

	if componentCount > 0 {

		rawScanQuery, err := db.ReplaceJSONplaceholders(replacements, constants.SecurityComponentFilterQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}

		scanQuery, err := db.ReplaceJSONplaceholders(replacements, constants.SecurityComponentFilterQueryByScanTime)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}

		queryMap := make(map[string]db.DbQuery)
		queryMap["scan"] = db.DbQuery{AliasName: constants.SECURITY_INDEX, QueryString: scanQuery}
		queryMap["rawScan"] = db.DbQuery{AliasName: constants.RAW_SCAN_RESULTS_INDEX, QueryString: rawScanQuery}

		responseMap, err := multiSearchResponse(queryMap)
		if log.CheckErrorf(err, "multi search query failed in securityComponentWidgetSection()") {
			return nil, err
		}

		type multiSearchResponse struct {
			Aggregations struct {
				DistinctComponent struct {
					Value []string `json:"value"`
				} `json:"distinct_component"`
			} `json:"aggregations"`
		}

		scanResult := multiSearchResponse{}
		err = json.Unmarshal(responseMap["scan"], &scanResult)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallRespGetCommitTrends) {
			return nil, err
		}

		rawScanResult := multiSearchResponse{}
		err = json.Unmarshal(responseMap["rawScan"], &rawScanResult)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallRespGetCommitTrends) {
			return nil, err
		}

		//Remove duplicates from each query response
		check := make(map[string]int)
		d := append(scanResult.Aggregations.DistinctComponent.Value, rawScanResult.Aggregations.DistinctComponent.Value...)
		combinedResult := make([]string, 0)
		for _, val := range d {
			check[val] = 1
		}

		for k, _ := range check {
			combinedResult = append(combinedResult, k)
		}

		responseJson, err := json.Marshal(combinedResult)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Debugf(exceptions.DebugEmptyComponentData)
	}
	return nil, nil
}

func UpdateFilters(updatedJSON string, replacements map[string]any) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(updatedJSON), &data)

	resultJson := updatedJSON
	if err == nil {
		filterArray := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})

		if replacements["component"].([]string) != nil {
			if replacements["component"].([]string)[0] != "All" {
				filter3 := helper.AddTermsFilter("component_id", replacements["component"].([]string))
				filterArray = append(filterArray, filter3)
			}
		}

		data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterArray
		modifiedData, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Warn(err.Error())
			modifiedData = []byte(updatedJSON)
		}
		resultJson = string(modifiedData)
	}
	return resultJson
}

func UpdateMustNotFilters(updatedJSON string, replacements map[string]any) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(updatedJSON), &data)

	resultJson := updatedJSON
	if err == nil {
		filterArray := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["must_not"].([]interface{})

		if replacements["jobIds"].([]string) != nil {
			filter3 := helper.AddTermsFilter("job_id", replacements["jobIds"].([]string))
			filterArray = append(filterArray, filter3)
		}

		data["query"].(map[string]interface{})["bool"].(map[string]interface{})["must_not"] = filterArray
		modifiedData, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Warn(err.Error())
			modifiedData = []byte(updatedJSON)
		}
		resultJson = string(modifiedData)
	}
	return resultJson
}

func UpdateMustFilters(updatedJSON string, replacements map[string]interface{}) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(updatedJSON), &data)
	resultJSON := updatedJSON
	if err == nil {
		boolObj, ok := data["query"].(map[string]interface{})["bool"].(map[string]interface{})
		if ok {
			mustArray, mustExists := boolObj["must"].([]interface{})
			if !mustExists {
				mustArray = make([]interface{}, 0)
			}
			if endpointIDs, ok := replacements["endpointIds"].([]string); ok && len(endpointIDs) > 0 {
				filter3 := helper.AddTermsFilter("endpoint_id", endpointIDs)
				mustArray = append(mustArray, filter3)
			}
			boolObj["must"] = mustArray
			filterArray, filterExists := boolObj["filter"].([]interface{})
			if !filterExists {
				filterArray = make([]interface{}, 0)
			}
			if parentIDs, ok := replacements["parentIds"].([]string); ok && len(parentIDs) > 0 {
				filter4 := helper.AddTermsFilter("org_id", parentIDs)
				filterArray = append(filterArray, filter4)
			}
			boolObj["filter"] = filterArray
			modifiedData, err := json.MarshalIndent(data, "", " ")
			if err != nil {
				log.Warn(err.Error())
				modifiedData = []byte(updatedJSON)
			}
			resultJSON = string(modifiedData)
		}
	}
	return resultJSON
}

func UpdateFiltersForDrilldown(updatedJSON string, replacements map[string]any, addAutomations bool, addBranchName bool) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(updatedJSON), &data)

	resultJson := updatedJSON
	if err == nil {
		filterArray := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})

		if replacements["component"].([]string) != nil {
			if replacements["component"].([]string)[0] != "All" {
				filter3 := helper.AddTermsFilter("component_id", replacements["component"].([]string))
				filterArray = append(filterArray, filter3)
			}
		}
		branch, ok := replacements[constants.REQUEST_BRANCH]
		if ok && branch != nil && addAutomations {
			log.Debugf("Inside Automation Filter")
			automations := getAutomationsForBranch(branch.(string))
			if len(automations) > 0 {
				log.Debugf("Adding  Automation Filter for branch:%s, Automations:%v", branch, automations)

				filter4 := helper.AddTermsFilter("automation_id", automations)
				filterArray = append(filterArray, filter4)
			}
		} else if ok && branch != nil && addBranchName {
			branchName := getBranchNameForId(branch.(string))
			log.Debugf("Inside Branch Filter")

			if len(branchName) > 0 {
				log.Debugf(exceptions.DebugAddingBranchFilterWithBranchName, branch, branchName)

				filter5 := helper.AddTermFilter("branch", branchName)
				filterArray = append(filterArray, filter5)

			}
		} else if ok && branch != nil {
			branchId := branch.(string)
			log.Debugf("Inside branch ID filter")

			if len(branchId) > 0 {
				log.Debugf(exceptions.DebugAddingBranchFilterWithBranchId, branch)

				filter5 := helper.AddTermFilter("branch_id", branchId)
				filterArray = append(filterArray, filter5)

			}

		}

		data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterArray
		modifiedData, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Warn(err.Error())
			modifiedData = []byte(updatedJSON)
		}
		resultJson = string(modifiedData)
	}
	return resultJson
}

func UpdateFiltersForSonar(updatedJSON string, replacements map[string]any) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(updatedJSON), &data)
	resultJson := updatedJSON
	if err == nil {
		filterArray := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})

		if replacements["component"].([]string) != nil {
			if replacements["component"].([]string)[0] != "All" {
				filter3 := helper.AddTermsFilter("component_id", replacements["component"].([]string))
				filterArray = append(filterArray, filter3)
			}
		}

		branch, ok := replacements[constants.REQUEST_BRANCH]
		if ok && branch != nil {

			if len(branch.(string)) > 0 {
				log.Debugf(exceptions.DebugAddingBranchFilterWithBranchName, branch, branch.(string))

				filter5 := helper.AddTermFilter("github_branch_id", branch.(string))
				filterArray = append(filterArray, filter5)

			}
		}

		data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterArray
		modifiedData, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Warn(err.Error())
			modifiedData = []byte(updatedJSON)
		}
		resultJson = string(modifiedData)
	}
	return resultJson
}

func automationWidgetHeader(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	log.Debugf("Automation widget header params ", widgetId, replacements)
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
	branch, ok := replacements[constants.REQUEST_BRANCH]

	if ok {
		log.Debugf("Changes for cs1 ")
		headerResponse, err := getAutomationWidgetDataForBranch(ctx, clt, orgId, components, branch.(string))
		if err != nil {
			return nil, err
		}
		return headerResponse, nil
	} else {
		headerResponse, err := getAutomationWidgetData(ctx, clt, orgId, components)
		if err != nil {
			return nil, err
		}
		return headerResponse, nil

	}

}

func automationWidgetSection(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	log.Debugf(exceptions.DebugAutomationWidgetSectionParams, widgetId, replacements)
	var updatedJSON string
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
	branch, ok := replacements[constants.REQUEST_BRANCH].(string)

	automationMap := GetAutomationMap(ctx, clt, orgId, components, branch)

	// testResource(ctx, clt, orgId)

	automationCount := float64(len(automationMap))
	if automationCount > 0 {
		client, err := openSearchClient()
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

		if ok {
			updatedJSON, err = db.ReplaceJSONplaceholders(replacements, constants.AutomationFilterQueryWithBranch)
			if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
				return nil, err
			}
		} else {
			updatedJSON, err = db.ReplaceJSONplaceholders(replacements, constants.AutomationFilterQuery)
			if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
				return nil, err
			}
		}

		response, err := searchResponse(updatedJSON, constants.AUTOMATION_METADATA_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		activeCount := 0.0
		inactiveCount := 0.0
		actives := []string{}
		result := make(map[string]interface{})
		json.Unmarshal([]byte(response), &result)
		if result[constants.AGGREGATION] != nil {
			aggsResult := result[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.DISTINCT_AUTOMATION] != nil {
				distinctAutomation := aggsResult[constants.DISTINCT_AUTOMATION].(map[string]interface{})
				if distinctAutomation[constants.VALUE] != nil {
					values := distinctAutomation[constants.VALUE].([]interface{})
					for _, v := range values {
						automationId := fmt.Sprint(v)
						_, ok := automationMap[automationId]
						if ok {
							actives = append(actives, automationId)
							activeCount++
						}
					}
				}
			}
		}
		log.Debugf(exceptions.DebugActiveAutomations, actives)
		inactiveCount = automationCount - activeCount
		outputResponse := make(map[string]interface{})
		drilldown := DrillDown{
			ReportId:    constants.WORKFLOWS,
			ReportTitle: constants.WORKFLOWS_DRILLDOWN,
			ReportType:  constants.STATUS,
		}
		outputResponse[constants.DATA] = []ChartData{
			{
				Name:  constants.ACTIVE,
				Value: int(math.Round((activeCount / automationCount) * 100)),
			}, {
				Name:  constants.INACTIVE,
				Value: int(math.Round((inactiveCount / automationCount) * 100)),
			},
		}
		outputResponse[constants.INFO] = []ChartInfo{
			{
				Title:     constants.ACTIVE,
				Value:     int(activeCount),
				Drilldown: drilldown,
			}, {
				Title:     constants.INACTIVE,
				Value:     int(inactiveCount),
				Drilldown: drilldown,
			},
		}
		responseJson, err := json.Marshal(outputResponse)

		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Debugf(exceptions.DebugEmptyAutomationData)
	}
	return nil, nil
}

// Fetches Component Comparison data for the Workflows widget (Associated IDs: e2, workflow-compare)
func workflowWidgetComponentComparison(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	log.Debugf("Automation widget component comparison params ", widgetId, replacements)
	var updatedJSON string
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, exceptions.ErrOrgIdNotSpecified)
	}

	componentAutomationMap, err := getAutomationMapComponentComparison(ctx, clt, orgId)
	if log.CheckErrorf(err, "failed to get componentAutomationMap from coreDataCache ") {
		return nil, err
	}

	componentCount := len(componentAutomationMap)
	if componentCount > 0 {
		client, err := openSearchClient()
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

		updatedJSON, err = db.ReplaceJSONplaceholders(replacements, constants.AutomationFilterQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}

		response, err := searchResponse(updatedJSON, constants.AUTOMATION_METADATA_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

		// create a map and store all the active automation IDs
		actives := make(map[string]struct{})
		result := make(map[string]interface{})
		json.Unmarshal([]byte(response), &result)
		if result[constants.AGGREGATION] != nil {
			aggsResult := result[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.DISTINCT_AUTOMATION] != nil {
				distinctAutomation := aggsResult[constants.DISTINCT_AUTOMATION].(map[string]interface{})
				if distinctAutomation[constants.VALUE] != nil {
					values := distinctAutomation[constants.VALUE].([]interface{})
					for _, v := range values {
						automationId := fmt.Sprint(v)
						actives[automationId] = struct{}{}

					}
				}
			}
		}
		log.Debugf(exceptions.DebugActiveAutomations, actives)

		outputMap := make(map[string]constants.WorkflowsComponentComparison)
		for componentID, automationMap := range componentAutomationMap {
			distinctAutomations := constants.WorkflowsComponentComparison{
				Active:   0,
				Inactive: 0,
			}
			for automationID := range automationMap {
				if _, ok := actives[automationID]; !ok {
					distinctAutomations.Inactive++
				} else {
					distinctAutomations.Active++
				}
			}
			outputMap[componentID] = distinctAutomations
		}

		responseJson, err := json.Marshal(outputMap)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Debugf(exceptions.DebugEmptyComponentDataForWidgetID, widgetId)
	}
	return nil, nil
}

func summaryAutomationWidgetSection(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	log.Debugf(exceptions.DebugAutomationWidgetSectionParams, widgetId, replacements)
	var updatedJSON string
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
	branch, ok := replacements[constants.REQUEST_BRANCH].(string)

	automationMap := GetAutomationMap(ctx, clt, orgId, components, branch)

	// testResource(ctx, clt, orgId)

	automationCount := float64(len(automationMap))
	log.Debugf("Count for branch:%v, automation count:%v", branch, automationCount)
	if automationCount > 0 {
		client, err := openSearchClient()
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

		if ok && len(branch) > 0 {
			log.Debugf("Changes for summary automation widget with branch")
			updatedJSON, err = db.ReplaceJSONplaceholders(replacements, constants.AutomationFilterQueryWithBranch)
			if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
				return nil, err
			}
		} else {
			log.Debugf("Changes for summary automation widget without branch")

			updatedJSON, err = db.ReplaceJSONplaceholders(replacements, constants.AutomationFilterQuery)
			if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
				return nil, err
			}
		}

		response, err := searchResponse(updatedJSON, constants.AUTOMATION_METADATA_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		activeCount := 0.0
		inactiveCount := 0.0
		actives := []string{}
		result := make(map[string]interface{})
		json.Unmarshal([]byte(response), &result)
		if result[constants.AGGREGATION] != nil {
			aggsResult := result[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.DISTINCT_AUTOMATION] != nil {
				distinctAutomation := aggsResult[constants.DISTINCT_AUTOMATION].(map[string]interface{})
				if distinctAutomation[constants.VALUE] != nil {
					values := distinctAutomation[constants.VALUE].([]interface{})
					for _, v := range values {
						automationId := fmt.Sprint(v)
						_, ok := automationMap[automationId]
						if ok {
							actives = append(actives, automationId)
							activeCount++
						}
					}
				}
			}
		}
		log.Debugf(exceptions.DebugActiveAutomations, actives)
		inactiveCount = automationCount - activeCount
		outputResponse := make(map[string]interface{})

		outputResponse[constants.DATA] = []SubTitleData{
			{
				Title: constants.ACTIVE,
				Value: int(activeCount),
			}, {
				Title: constants.INACTIVE,
				Value: int(inactiveCount),
			},
		}

		responseJson, err := json.Marshal(outputResponse[constants.DATA])
		if log.CheckErrorf(err, "Error in summary automation widget : ") {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Debugf(exceptions.DebugEmptyAutomationData)
	}
	return nil, nil
}

func securityAutomationWidgetSection(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {

	log.Debugf(exceptions.DebugAutomationWidgetSectionParams, widgetId, replacements)
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

	automationMap := GetAutomationMap(ctx, clt, orgId, components, "")

	automationCount := float64(len(automationMap))
	if automationCount > 0 {
		client, err := openSearchClient()
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.SecurityAutomationFilterQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		response, err := searchResponse(updatedJSON, constants.SECURITY_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		activeCount := 0.0
		inactiveCount := 0.0
		result := make(map[string]interface{})
		json.Unmarshal([]byte(response), &result)
		if result[constants.AGGREGATION] != nil {
			aggsResult := result[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.DISTINCT_AUTOMATION] != nil {
				distinctAutomation := aggsResult[constants.DISTINCT_AUTOMATION].(map[string]interface{})
				if distinctAutomation[constants.VALUE] != nil {
					values := distinctAutomation[constants.VALUE].([]interface{})
					for _, v := range values {
						automationId := fmt.Sprint(v)
						_, ok := automationMap[automationId]
						if ok {
							activeCount++
						}
					}
				}
			}
		}
		inactiveCount = automationCount - activeCount
		outputResponse := make(map[string]interface{})
		drilldown := DrillDown{
			ReportId:    constants.SECURITY_WORKFLOWS,
			ReportTitle: constants.WORKFLOWS_DRILLDOWN,
			ReportType:  constants.SCANNERS_LOWERCASE,
		}
		outputResponse[constants.DATA] = []ChartData{
			{
				Name:  constants.WITHSCANNERS,
				Value: int(math.Round((activeCount / automationCount) * 100)),
			}, {
				Name:  constants.WITHOUTSCANNERS,
				Value: int(math.Round((inactiveCount / automationCount) * 100)),
			},
		}
		outputResponse[constants.INFO] = []ChartInfo{
			{
				Title:     constants.WITHSCANNERS,
				Value:     int(activeCount),
				Drilldown: drilldown,
			}, {
				Title:     constants.WITHOUTSCANNERS,
				Value:     int(inactiveCount),
				Drilldown: drilldown,
			},
		}
		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Debugf(exceptions.DebugEmptyAutomationData)
	}
	return nil, nil
}

// Fetches Component Comparison data for the Workflows widget in the Test Insights dashboard (Associated IDs: ti2, test-insights-workflows-compare)
func testInsightsWorkflowsComponentComparison(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	log.Debugf("test insights workflows widget component comparison params ", widgetId, replacements)
	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, exceptions.ErrOrgIdNotSpecified)
	}
	testComponentAutomationMap, err := getAutomationMapComponentComparison(ctx, clt, orgId)
	if log.CheckErrorf(err, "failed to get testComponentAutomationMap from coreDataCache ") {
		return nil, err
	}
	testComponentCount := len(testComponentAutomationMap)
	if testComponentCount > 0 {
		client, err := openSearchClient()
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TestInsightsAutomationFilterQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		response, err := searchResponse(updatedJSON, constants.TEST_SUITE_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

		// create a map and store all the active automation IDs
		actives := make(map[string]struct{})
		result := make(map[string]interface{})
		json.Unmarshal([]byte(response), &result)
		if result[constants.AGGREGATION] != nil {
			aggsResult := result[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.TEST_INSIGHTS_AUTOMATIONS] != nil {
				testInsightsAutomation := aggsResult[constants.TEST_INSIGHTS_AUTOMATIONS].(map[string]interface{})
				if testInsightsAutomation[constants.VALUE] != nil {
					values := testInsightsAutomation[constants.VALUE].([]interface{})
					for _, v := range values {
						automationId := fmt.Sprint(v)
						actives[automationId] = struct{}{}

					}
				}
			}
		}
		log.Debugf(exceptions.DebugActiveAutomations, actives)
		outputMap := make(map[string]constants.TestWorkflowsComponentComparison)
		for componentID, automationMap := range testComponentAutomationMap {
			testInsightsAutomations := constants.TestWorkflowsComponentComparison{
				WithTestSuites:    0,
				WithoutTestSuites: 0,
			}
			for automationID := range automationMap {
				if _, ok := actives[automationID]; !ok {
					testInsightsAutomations.WithoutTestSuites++
				} else {
					testInsightsAutomations.WithTestSuites++
				}
			}
			outputMap[componentID] = testInsightsAutomations
		}

		responseJson, err := json.Marshal(outputMap)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Debugf(exceptions.DebugEmptyComponentDataForWidgetID, widgetId)
	}
	return nil, nil

}

// Fetches Component Comparison data for the Workflows widget in the Security Insights dashboard (Associated IDs: s2, workflow-compare)
func workflowComponentComparisonSI(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	log.Debugf("Automation widget component comparison params ", widgetId, replacements)

	orgId, ok := replacements[constants.ORG_ID].(string)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, exceptions.ErrOrgIdNotSpecified)
	}

	componentAutomationMap, err := getAutomationMapComponentComparison(ctx, clt, orgId)
	if log.CheckErrorf(err, "failed to get componentAutomationMap from coreDataCache ") {
		return nil, err
	}

	componentCount := len(componentAutomationMap)
	if componentCount > 0 {
		client, err := openSearchClient()
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.SecurityAutomationFilterQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}

		response, err := searchResponse(updatedJSON, constants.SECURITY_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)

		// create a map and store all the active automation IDs
		actives := make(map[string]struct{})
		result := make(map[string]interface{})
		json.Unmarshal([]byte(response), &result)
		if result[constants.AGGREGATION] != nil {
			aggsResult := result[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.DISTINCT_AUTOMATION] != nil {
				distinctAutomation := aggsResult[constants.DISTINCT_AUTOMATION].(map[string]interface{})
				if distinctAutomation[constants.VALUE] != nil {
					values := distinctAutomation[constants.VALUE].([]interface{})
					for _, v := range values {
						automationId := fmt.Sprint(v)
						actives[automationId] = struct{}{}

					}
				}
			}
		}
		log.Debugf(exceptions.DebugActiveAutomations, actives)

		outputMap := make(map[string]constants.SecurityWorkflowsComponentComparison)
		for componentID, automationMap := range componentAutomationMap {
			distinctAutomations := constants.SecurityWorkflowsComponentComparison{
				WithScanners:    0,
				WithoutScanners: 0,
			}
			for automationID := range automationMap {
				if _, ok := actives[automationID]; !ok {
					distinctAutomations.WithoutScanners++
				} else {
					distinctAutomations.WithScanners++
				}
			}
			outputMap[componentID] = distinctAutomations
		}

		responseJson, err := json.Marshal(outputMap)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Debugf(exceptions.DebugEmptyComponentDataForWidgetID, widgetId)
	}
	return nil, nil
}

func getAutomationMap(ctx context.Context, clt client.GrpcClient, orgId string, components []string, branch string) map[string]struct{} {
	automationMap := map[string]struct{}{}
	coreDataCache := cache.GetCoreDataCache()
	if coreDataCache != nil {
		startTime := time.Now()
		serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
		if err != nil {
			return automationMap
		}
		log.Debugf(exceptions.DebugTimeTookToFetchAllServiceMilliSec, time.Since(startTime).Microseconds())
		startTime = time.Now()
		if serviceResponse != nil {
			for i := 0; i < len(serviceResponse.GetService()); i++ {
				found := false
				service := serviceResponse.GetService()[i]
				if len(components) > 0 {
					for _, component := range components {
						if service.Id == component {
							found = true
							break
						}
					}
				}
				if found || components == nil {
					for _, child := range coreDataCache.GetChildrenOfType(service.Id, api.ResourceType_RESOURCE_TYPE_BRANCH) {
						childResource := coreDataCache.Get(child)
						if !childResource.IsDisabled && (len(branch) == 0 || child == branch) {
							automations := coreDataCache.GetChildrenOfType(child, api.ResourceType_RESOURCE_TYPE_AUTOMATION)
							if len(automations) > 0 {
								for _, id := range automations {
									automationResource := coreDataCache.Get(id)
									if !automationResource.IsDisabled {
										automationMap[id] = struct{}{}
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

func GetAutomationsForBranch(branch string) []string {

	if len(branch) > 0 {
		coreDataCache := cache.GetCoreDataCache()

		automations := coreDataCache.GetChildrenOfType(branch, api.ResourceType_RESOURCE_TYPE_AUTOMATION)
		log.Debugf("For Branch:%s, Automations:%v", branch, automations)
		return automations
	}
	return nil
}
func GetBranchNameForId(branch string) string {

	if len(branch) > 0 {
		coreDataCache := cache.GetCoreDataCache()

		childResource := coreDataCache.Get(branch)
		if childResource != nil {
			return childResource.Name

		}
	}
	return ""
}

func getComponentWidgetData(ctx context.Context, clt client.GrpcClient, orgId string, components []string) ([]byte, error) {
	componentCount := 0
	repoCount := 0
	reposUrls := map[string]struct{}{}
	duplicateUrls := map[string]struct{}{}
	serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
	log.Debugf("Org Id : %s Components : %v", orgId, components)
	log.Debugf("Service Response : %v", serviceResponse)
	if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
		return nil, err
	}
	if serviceResponse != nil {
		for i := 0; i < len(serviceResponse.GetService()); i++ {
			service := serviceResponse.GetService()[i]
			url := service.GetRepositoryUrl()
			if _, ok := reposUrls[url]; ok {
				duplicateUrls[url] = struct{}{}
			}
			if components != nil {
				for _, component := range components {
					if service.Id == component {
						componentCount++
						reposUrls[url] = struct{}{}
					}
				}
			} else {
				componentCount++
				reposUrls[url] = struct{}{}
			}
		}
	}
	repoCount = len(reposUrls)
	duplicateUrlString := "Duplicate Repos : "
	for k := range duplicateUrls {
		duplicateUrlString += (k + " ")
	}
	log.Debugf(duplicateUrlString)
	componentData := &ComponentAndAutomationData{
		Value: componentCount,
		SubTitle: SubTitleData{
			Title: constants.REPOS,
			Value: repoCount,
		},
	}
	response, err := json.Marshal(componentData)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getAutomationWidgetData(ctx context.Context, clt client.GrpcClient, orgId string, components []string) (json.RawMessage, error) {
	automationCount := 0
	branchesCount := 0
	automationMap := map[string]struct{}{}
	coreDataCache := cache.GetCoreDataCache()
	var multipleAutomation []string

	startTime := time.Now()
	serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
	log.Debugf("Org Id : %s Components for automation widget : %v", orgId, components)
	log.Debugf("Service Response for automation widget : %v", serviceResponse)
	if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
		return nil, err
	}
	log.Debugf(exceptions.DebugTimeTookToFetchAllServiceMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()
	if serviceResponse != nil {
		for i := 0; i < len(serviceResponse.GetService()); i++ {
			found := false
			service := serviceResponse.GetService()[i]
			if len(components) > 0 {
				for _, component := range components {
					if service.Id == component {
						found = true
						break
					}
				}
			}
			if found || components == nil {
				for _, child := range coreDataCache.GetChildrenOfType(service.Id, api.ResourceType_RESOURCE_TYPE_BRANCH) {
					childResource := coreDataCache.Get(child)
					if !childResource.IsDisabled {
						automations := coreDataCache.GetChildrenOfType(child, api.ResourceType_RESOURCE_TYPE_AUTOMATION)
						var automationWithoutDeletedWorkflows []string
						for _, automation := range automations {
							automationResource := coreDataCache.Get(automation)
							log.Debugf("The automation resource for component %+v is %+v", automationResource, service.Id)
							if !automationResource.IsDisabled {
								automationWithoutDeletedWorkflows = append(automationWithoutDeletedWorkflows, automation)
							}
						}
						count := len(automationWithoutDeletedWorkflows)
						if count > 0 {
							for _, id := range automationWithoutDeletedWorkflows {
								automationMap[id] = struct{}{}
							}
							branchesCount++
							if count > 1 {
								branchInfo := fmt.Sprintf("%s%s%s%s%v", service.Name, " ", childResource.Name, " ", len(automationWithoutDeletedWorkflows))
								multipleAutomation = append(multipleAutomation, branchInfo)
							}
							automationCount += count
						}
					}
				}
			}
		}
	}
	log.Debugf("Total Time took to process and get automation and branch info : %v in milliseconds", time.Since(startTime).Milliseconds())
	log.Debugf("Multiple Automation information : ", multipleAutomation)

	automationData := &ComponentAndAutomationData{
		Value: automationCount,
		SubTitle: SubTitleData{
			Title: constants.BRANCHES,
			Value: branchesCount,
		},
	}
	response, err := json.Marshal(automationData)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getAutomationWidgetDataForBranch(ctx context.Context, clt client.GrpcClient, orgId string, components []string, branch string) (json.RawMessage, error) {
	automationCount := 0
	branchesCount := 0
	automationMap := map[string]struct{}{}
	coreDataCache := cache.GetCoreDataCache()
	var multipleAutomation []string

	startTime := time.Now()
	serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
	log.Infof("Org Id : %s Components for automation widget : %v", orgId, components)
	log.Infof("Service Response for automation widget : %v", serviceResponse)
	if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
		return nil, err
	}
	log.Debugf(exceptions.DebugTimeTookToFetchAllServiceMilliSec, time.Since(startTime).Milliseconds())
	startTime = time.Now()
	if serviceResponse != nil {
		for i := 0; i < len(serviceResponse.GetService()); i++ {
			found := false
			service := serviceResponse.GetService()[i]
			if len(components) > 0 {
				for _, component := range components {
					if service.Id == component {
						found = true
						break
					}
				}
			}
			if found || components == nil {
				for _, child := range coreDataCache.GetChildrenOfType(service.Id, api.ResourceType_RESOURCE_TYPE_BRANCH) {
					childResource := coreDataCache.Get(child)
					if childResource.Id == branch {
						automations := coreDataCache.GetChildrenOfType(child, api.ResourceType_RESOURCE_TYPE_AUTOMATION)
						var automationWithoutDeletedWorkflows []string
						for _, automation := range automations {
							automationResource := coreDataCache.Get(automation)
							if !automationResource.IsDisabled {
								automationWithoutDeletedWorkflows = append(automationWithoutDeletedWorkflows, automation)
							}
						}
						count := len(automationWithoutDeletedWorkflows)
						if count > 0 {
							for _, id := range automationWithoutDeletedWorkflows {
								automationMap[id] = struct{}{}
							}
							branchesCount++
							if count > 1 {
								branchInfo := fmt.Sprintf("%s%s%s%s%v", service.Name, " ", childResource.Name, " ", len(automationWithoutDeletedWorkflows))
								multipleAutomation = append(multipleAutomation, branchInfo)
							}
						}
					}
				}
			}
		}
		automationCount = len(automationMap)
	}
	log.Debugf("Total Time took to process and get automation and branch info : %v in milliseconds", time.Since(startTime).Milliseconds())
	log.Debugf("Multiple Automation information : ", multipleAutomation)

	automationData := &ComponentAndAutomationBranchData{
		Value: automationCount,
	}
	response, err := json.Marshal(automationData)
	if err != nil {
		return nil, err
	}
	return response, nil
}

/* Fetches and groups all the automation IDs in an org by Component ID. Returns a map[string]map[string]struct{} where the keys at the first level are component IDs and the second level are automation IDs.*/
func getAutomationMapComponentComparison(ctx context.Context, clt client.GrpcClient, orgId string) (map[string]map[string]struct{}, error) {

	componentAutomationsMap := map[string]map[string]struct{}{}
	coreDataCache := cache.GetCoreDataCache()
	if coreDataCache != nil {
		startTime := time.Now()
		serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
		if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
			return nil, err
		}

		log.Debugf(exceptions.DebugTimeTookToFetchAllServiceMilliSec, time.Since(startTime).Microseconds())
		startTime = time.Now()
		if serviceResponse != nil {
			for i := 0; i < len(serviceResponse.GetService()); i++ {

				service := serviceResponse.GetService()[i]

				for _, child := range coreDataCache.GetChildrenOfType(service.Id, api.ResourceType_RESOURCE_TYPE_BRANCH) {
					childResource := coreDataCache.Get(child)
					if !childResource.IsDisabled {
						automations := coreDataCache.GetChildrenOfType(child, api.ResourceType_RESOURCE_TYPE_AUTOMATION)
						if len(automations) > 0 {
							if _, ok := componentAutomationsMap[service.Id]; !ok {
								componentAutomationsMap[service.Id] = map[string]struct{}{}
							}
							automationMap := componentAutomationsMap[service.Id]
							for _, id := range automations {
								automationMap[id] = struct{}{}
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
	return componentAutomationsMap, nil
}

func UpdateMustParentFilters(updatedJSON string, replacements map[string]any) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(updatedJSON), &data)
	resultJson := updatedJSON
	if err == nil {
		boolObj := data["query"].(map[string]interface{})["bool"].(map[string]interface{})
		mustArray, ok := boolObj["must"]
		filterArray := []interface{}{}
		if ok {
			filterArray = mustArray.([]interface{})
		}
		if replacements["parentIds"].([]string) != nil {
			filter3 := helper.AddTermsFilter("org_id", replacements["parentIds"].([]string))
			filterArray = append(filterArray, filter3)
		}
		data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterArray
		modifiedData, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Warn(err.Error())
			modifiedData = []byte(updatedJSON)
		}
		resultJson = string(modifiedData)
	}
	return resultJson
}

func testInsightsWidgetSection(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {

	log.Debugf("Test Insights Automation widget section params ", widgetId, replacements)
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

	automationMap := GetAutomationMap(ctx, clt, orgId, components, "")

	automationCount := float64(len(automationMap))
	if automationCount > 0 {
		client, err := openSearchClient()
		log.CheckErrorf(err, exceptions.ErrOpenSearchConnection)

		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TestInsightsAutomationFilterQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		response, err := searchResponse(updatedJSON, constants.TEST_SUITE_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		withTestSuitesCount := 0.0
		withoutTestSuitesCount := 0.0
		result := make(map[string]interface{})
		json.Unmarshal([]byte(response), &result)
		if result[constants.AGGREGATION] != nil {
			aggsResult := result[constants.AGGREGATION].(map[string]interface{})
			if aggsResult[constants.TEST_INSIGHTS_AUTOMATIONS] != nil {
				testInsightsAutomation := aggsResult[constants.TEST_INSIGHTS_AUTOMATIONS].(map[string]interface{})
				if testInsightsAutomation[constants.VALUE] != nil {
					values := testInsightsAutomation[constants.VALUE].([]interface{})
					for _, v := range values {
						automationId := fmt.Sprint(v)
						_, ok := automationMap[automationId]
						if ok {
							withTestSuitesCount++
						}
					}
				}
			}
		}
		withoutTestSuitesCount = automationCount - withTestSuitesCount
		outputResponse := make(map[string]interface{})
		drilldown := DrillDown{
			ReportId:    constants.TEST_INSIGHTS_WORKFLOWS,
			ReportTitle: constants.TEST_INSIGHTS_WORKFLOWS_DRILLDOWN,
			ReportType:  constants.TEST_SUITE,
		}
		// For calculating the percentage in Donut chart
		outputResponse[constants.DATA] = []ChartData{
			{
				Name:  constants.WITH_TEST_SUITES,
				Value: int(math.Round((withTestSuitesCount / automationCount) * 100)),
			}, {
				Name:  constants.WITHOUT_TEST_SUITES,
				Value: int(math.Round((withoutTestSuitesCount / automationCount) * 100)),
			},
		}
		// For section data
		outputResponse[constants.INFO] = []ChartInfo{
			{
				Title:     constants.WITH_TEST_SUITES,
				Value:     int(withTestSuitesCount),
				Drilldown: drilldown,
			}, {
				Title:     constants.WITHOUT_TEST_SUITES,
				Value:     int(withoutTestSuitesCount),
				Drilldown: drilldown,
			},
		}
		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Info("Test Insights Automation Data is not present")
	}
	return nil, nil
}

func testComponentWidgetSection(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	log.Debugf("Test Insights component widget section params ", widgetId, replacements)
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
	componentMap := map[string]struct{}{}
	serviceResponse, err := getOrganisationServices(ctx, clt, orgId)
	if log.CheckErrorf(err, exceptions.ErrFetchingServiceTemplate, orgId) {
		return nil, err
	}

	if serviceResponse != nil {
		for i := 0; i < len(serviceResponse.GetService()); i++ {
			service := serviceResponse.GetService()[i]
			if components != nil {
				for _, component := range components {
					if service.Id == component {
						componentMap[service.Id] = struct{}{}
					}
				}
			} else {
				componentMap[service.Id] = struct{}{}
			}
		}
	}

	componentCount := float64(len(componentMap))
	if componentCount > 0 {
		client := db.GetOpenSearchClient()
		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TestComponentFilterQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson := UpdateFilters(updatedJSON, replacements)
		response, err := searchResponse(modifiedJson, constants.TEST_SUITE_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		withTestSuitesCount := 0.0
		withoutTestSuitesCount := 0.0

		testComponentResponse := constants.TestComponentResponse{}
		json.Unmarshal([]byte(response), &testComponentResponse)
		withTestSuitesCount = float64(len(testComponentResponse.Aggregations.DistinctComponent.Value))

		withoutTestSuitesCount = componentCount - withTestSuitesCount
		outputResponse := make(map[string]interface{})
		drilldown := DrillDown{
			ReportId:    constants.TEST_INSIGHTS_COMPONENT,
			ReportTitle: constants.COMPONENTS_DRILLDOWN,
			ReportType:  constants.TEST_SUITE,
		}

		outputResponse[constants.DATA] = []ChartData{
			{
				Name:  constants.WITH_TEST_SUITES,
				Value: int(math.Round((withTestSuitesCount / componentCount) * 100)),
			}, {
				Name:  constants.WITHOUT_TEST_SUITES,
				Value: int(math.Round((withoutTestSuitesCount / componentCount) * 100)),
			},
		}
		outputResponse[constants.INFO] = []ChartInfo{
			{
				Title:     constants.WITH_TEST_SUITES,
				Value:     int(withTestSuitesCount),
				Drilldown: drilldown,
			}, {
				Title:     constants.WITHOUT_TEST_SUITES,
				Value:     int(withoutTestSuitesCount),
				Drilldown: drilldown,
			},
		}

		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		log.Info(exceptions.DebugEmptyComponentData)
	}
	return nil, nil
}

func getTestsOverview(widgetId string, replacements map[string]any, ctx context.Context, clt client.GrpcClient, epClt endpoint.EndpointServiceClient) (json.RawMessage, error) {
	client, err := openSearchClient()
	if log.CheckErrorf(err, exceptions.ErrOpenSearchConnection) {
		return nil, err
	}

	viewOption := replacements[constants.VIEW_OPTION].(string)

	if viewOption == constants.TEST_SUITE_VIEW || viewOption == "" {
		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TestSuitesOverviewQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson := UpdateFilters(updatedJSON, replacements)

		response, err := searchResponse(modifiedJson, constants.TEST_SUITE_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		outputResponse, err := updateTestSuitesOverviewResponse(response)
		if log.CheckErrorf(err, "error in transforming query response : ") {
			return nil, err
		}
		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else if viewOption == constants.COMPONENTS_VIEW {
		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TestOverviewComponentsViewQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson := UpdateFilters(updatedJSON, replacements)
		response, err := searchResponse(modifiedJson, constants.TEST_SUITE_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		outputResponse, err := updateTestComponentsViewResponse(response)
		if log.CheckErrorf(err, "error in transforming query response : ") {
			return nil, err
		}
		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else if viewOption == constants.TEST_CASE_VIEW {
		updatedJSON, err := db.ReplaceJSONplaceholders(replacements, constants.TestCasesOverviewQuery)
		if log.CheckErrorf(err, exceptions.ErrJsonPlaceholderNotReplaceable, replacements) {
			return nil, err
		}
		modifiedJson := UpdateFilters(updatedJSON, replacements)
		response, err := searchResponse(modifiedJson, constants.TEST_CASES_INDEX, client)
		log.CheckErrorf(err, exceptions.ErrOpenSearchFetchDataFailure)
		outputResponse, err := updateTestCasesOverviewResponse(response)
		if log.CheckErrorf(err, "error in transforming query response : ") {
			return nil, err
		}
		responseJson, err := json.Marshal(outputResponse)
		if err != nil {
			return nil, err
		} else {
			return responseJson, nil
		}
	} else {
		return nil, fmt.Errorf("invalid view option - %s", viewOption)
	}
}

func updateTestSuitesOverviewResponse(response string) ([]TestSuitesOverview, error) {
	result := TestSuitesOverviewResponse{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in updateTestSuitesOverviewResponse()") {
		return nil, err
	}

	testSuitesOverview := []TestSuitesOverview{}

	for _, value := range result.Aggregations.TestSuitesOverview.Value {
		if value.TestSuiteName == "" {
			continue
		}
		testSuiteInfo := TestSuitesOverview{}
		testSuiteInfo.ComponentName = value.ComponentName
		testSuiteInfo.TestSuiteName = value.TestSuiteName
		workflowName, sourceName, err := cutils.GetDisplayNameAndOrigin(value.AutomationName)
		if log.CheckErrorf(err, "Error getting display name and origin for updateTestSuitesOverviewResponse()", err) {
			log.Infof("Error getting display name and origin for updateTestSuitesOverviewResponse()")
		}
		testSuiteInfo.Workflow = workflowName
		testSuiteInfo.Source = sourceName
		testSuiteInfo.Branch = value.BranchName
		testSuiteInfo.DefaultBranch = value.BranchName
		testSuiteInfo.LastRun = value.StartTime
		testSuiteInfo.LastRunInMillis = int(value.StartTimeInMillis)
		testSuiteInfo.TotalTestCases.Value = value.Total
		testSuiteInfo.TotalTestCases.DrillDown.ReportId = "test-overview-total-tests-cases"
		testSuiteInfo.TotalTestCases.DrillDown.ReportTitle = "Test cases - " + testSuiteInfo.TestSuiteName
		testSuiteInfo.TotalTestCases.DrillDown.ReportInfo = pb.ReportInfo{
			TestSuiteName: value.TestSuiteName,
			Branch:        value.BranchID,
			AutomationId:  value.AutomationID,
			ComponentName: value.ComponentName,
			BranchName:    value.BranchName,
			WorkflowName:  workflowName,
			Source:        sourceName,
		}
		testSuiteInfo.TotalTestCasesValue = value.Total
		testSuiteInfo.AvgRunTime = value.AverageDuration
		testSuiteInfo.TotalRuns.Value = value.Runs
		testSuiteInfo.TotalRuns.DrillDown.ReportId = "test-overview-total-runs"
		testSuiteInfo.TotalRuns.DrillDown.ReportTitle = "Runs - " + value.TestSuiteName
		testSuiteInfo.TotalRuns.DrillDown.ReportInfo = pb.ReportInfo{
			TestSuiteName: value.TestSuiteName,
			Branch:        value.BranchID,
			AutomationId:  value.AutomationID,
		}
		testSuiteInfo.TotalRunsValue = value.Runs

		// Dark color scheme for the failure rate bar
		testSuiteInfo.FailureRate.ColorScheme = append(testSuiteInfo.FailureRate.ColorScheme, ColorScheme{
			Color0: "#009C5B",
			Color1: "#62CA9D",
		})

		testSuiteInfo.FailureRate.ColorScheme = append(testSuiteInfo.FailureRate.ColorScheme, ColorScheme{
			Color0: "#D32227",
			Color1: "#FB6E72",
		})

		testSuiteInfo.FailureRate.ColorScheme = append(testSuiteInfo.FailureRate.ColorScheme, ColorScheme{
			Color0: "#F2A414",
			Color1: "#FFE6C1",
		})

		// Light color scheme for the failure rate bar
		testSuiteInfo.FailureRate.LightColorScheme = append(testSuiteInfo.FailureRate.LightColorScheme, ColorScheme{
			Color0: "#0C9E61",
			Color1: "#79CAA8",
		})

		testSuiteInfo.FailureRate.LightColorScheme = append(testSuiteInfo.FailureRate.LightColorScheme, ColorScheme{
			Color0: "#E83D39",
			Color1: "#F39492",
		})

		testSuiteInfo.FailureRate.LightColorScheme = append(testSuiteInfo.FailureRate.LightColorScheme, ColorScheme{
			Color0: "#F2A414",
			Color1: "#FFE6C1",
		})

		// Fialure rate tool tip data
		testSuiteInfo.FailureRate.Data = append(testSuiteInfo.FailureRate.Data, struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		}{
			Title: "Successful test cases",
			Value: int(value.SuccessfulCasesCount),
		})
		testSuiteInfo.FailureRate.Data = append(testSuiteInfo.FailureRate.Data, struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		}{
			Title: "Failed test cases",
			Value: int(value.FailedCasesCount),
		})
		testSuiteInfo.FailureRate.Data = append(testSuiteInfo.FailureRate.Data, struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		}{
			Title: "Skipped test cases",
			Value: int(value.SkippedCasesCount),
		})

		failureRateFloat := helper.ConvertPercentageToFloat(value.FailureRateForLastRun)
		testSuiteInfo.FailureRateValue = failureRateFloat

		testSuiteInfo.FailureRate.Value = value.FailureRateForLastRun

		testSuiteInfo.FailureRate.Type = pb.ChartType_SINGLE_BAR.String()

		testSuitesOverview = append(testSuitesOverview, testSuiteInfo)
	}

	// Sorting test suites based on the failure rate in descending order
	sort.Slice(testSuitesOverview, func(i, j int) bool {
		return testSuitesOverview[i].FailureRateValue > testSuitesOverview[j].FailureRateValue
	})

	return testSuitesOverview, nil
}

func updateTestComponentsViewResponse(response string) ([]TestComponentsView, error) {
	result := TestComponentsViewResponse{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in updateTestSuitesOverviewResponse()") {
		return nil, err
	}

	testComponentsView := []TestComponentsView{}

	for _, value := range result.Aggregations.WorkflowBuckets.Buckets {
		testComponentViewInfo := TestComponentsView{}
		testComponentViewInfo.ComponentName = value.LatestDoc.Hits.Hits[0].Source.ComponentName
		testComponentViewInfo.AvgRunTime = math.Round(value.AvgRunTime.Value*100) / 100
		testComponentViewInfo.DefaultBranch = value.LatestDoc.Hits.Hits[0].Source.BranchName

		truncatedValue := helper.TruncateFloat(value.FailureRate.Value)
		testComponentViewInfo.FailureRateValue = truncatedValue

		testComponentViewInfo.LastRun = value.LatestDoc.Hits.Hits[0].Fields.ZonedRunStartTime[0]
		testComponentViewInfo.LastRunInMillis = value.LatestDoc.Hits.Hits[0].Fields.RunStartTimeInMillis[0]
		// testSuiteInfo.AvgRunTime = millisecondsToSecondsString(value.AverageDuration)
		testComponentViewInfo.TotalTestCasesValue = int(value.TotalTestCasesCount.Value)
		workflowName, sourceName, err := cutils.GetDisplayNameAndOrigin(value.LatestDoc.Hits.Hits[0].Source.AutomationName)
		if log.CheckErrorf(err, "Error getting display name and origin for updateTestComponentsViewResponse()", err) {
			log.Infof("Error getting display name and origin for updateTestComponentsViewResponse()")
		}
		testComponentViewInfo.Workflow = workflowName
		testComponentViewInfo.Source = sourceName

		testComponentViewInfo.FailureRate.Type = pb.ChartType_SINGLE_BAR.String()
		testComponentViewInfo.FailureRate.ColorScheme = append(testComponentViewInfo.FailureRate.ColorScheme, struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			Color0: "#009C5B",
			Color1: "#62CA9D",
		})
		testComponentViewInfo.FailureRate.ColorScheme = append(testComponentViewInfo.FailureRate.ColorScheme, struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			Color0: "#D32227",
			Color1: "#FB6E72",
		})
		testComponentViewInfo.FailureRate.ColorScheme = append(testComponentViewInfo.FailureRate.ColorScheme, struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			Color0: "#F2A414",
			Color1: "#FFE6C1",
		})
		testComponentViewInfo.FailureRate.LightColorScheme = append(testComponentViewInfo.FailureRate.LightColorScheme, struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			Color0: "#0C9E61",
			Color1: "#79CAA8",
		})
		testComponentViewInfo.FailureRate.LightColorScheme = append(testComponentViewInfo.FailureRate.LightColorScheme, struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			Color0: "#E83D39",
			Color1: "#F39492",
		})
		testComponentViewInfo.FailureRate.LightColorScheme = append(testComponentViewInfo.FailureRate.LightColorScheme, struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			Color0: "#F2A414",
			Color1: "#FFE6C1",
		})
		// Populate the Data array
		testComponentViewInfo.FailureRate.Data = append(testComponentViewInfo.FailureRate.Data, struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		}{
			Title: "Successful test case runs",
			Value: int(value.SuccessCount.Value),
		})
		testComponentViewInfo.FailureRate.Data = append(testComponentViewInfo.FailureRate.Data, struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		}{
			Title: "Failed test case runs",
			Value: int(value.FailureCount.Value),
		})
		testComponentViewInfo.FailureRate.Data = append(testComponentViewInfo.FailureRate.Data, struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		}{
			Title: "Skipped test case runs",
			Value: int(value.SkippedCount.Value),
		})

		testComponentViewInfo.FailureRate.Value = fmt.Sprintf("%.1f%%", value.FailureRate.Value)

		testComponentViewInfo.TotalTestCases.DrillDown.ReportId = "test-overview-total-tests-cases"
		testComponentViewInfo.TotalTestCases.DrillDown.ReportTitle = "Test cases - " + value.LatestDoc.Hits.Hits[0].Source.ComponentName
		testComponentViewInfo.TotalTestCases.DrillDown.ReportInfo = pb.ReportInfo{
			Branch:        value.LatestDoc.Hits.Hits[0].Source.BranchID,
			AutomationId:  value.LatestDoc.Hits.Hits[0].Source.AutomationID,
			ComponentName: value.LatestDoc.Hits.Hits[0].Source.ComponentName,
			WorkflowName:  workflowName,
			Source:        sourceName,
			BranchName:    value.LatestDoc.Hits.Hits[0].Source.BranchName,
		}
		testComponentViewInfo.TotalTestCases.Value = int(value.TotalTestCasesCount.Value)

		testComponentsView = append(testComponentsView, testComponentViewInfo)
	}

	return testComponentsView, nil
}

func updateTestCasesOverviewResponse(response string) ([]TestCasesOverview, error) {

	result := TestCasesOverviewResponse{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in updateTestCasesOverviewResponse(/.)") {
		return nil, err
	}

	// Preallocate slice capacity
	testCasesOverview := make([]TestCasesOverview, 0, len(result.Aggregations.TestCasesOverview.Value))

	// Use WaitGroup for parallel processing
	// var wg sync.WaitGroup
	// var mu sync.Mutex // Mutex to safely append to the slice

	// // Process each value in parallel
	// for _, value := range result.Aggregations.TestCasesOverview.Value {
	// 	wg.Add(1)
	// 	go func(value TestCaseResponse) {
	// 		defer wg.Done()

	// 		testCaseInfo := TestCasesOverview{
	// 			TestSuiteName:   value.TestSuiteName,
	// 			ComponentName:   value.ComponentName,
	// 			Workflow:        value.AutomationName,
	// 			Branch:          value.BranchName,
	// 			LastRun:         value.StartTime,
	// 			LastRunInMillis: int(value.StartTimeInMillis),
	// 			AvgRunTime:      value.AverageDuration,
	// 			TestCaseName:    value.TestCaseName,
	// 		}

	// 		testCaseInfo.TotalRuns.Value = value.Runs
	// 		testCaseInfo.TotalRuns.DrillDown.ReportId = "test-overview-view-run-activity"
	// 		testCaseInfo.TotalRuns.DrillDown.ReportTitle = "Runs"
	// 		testCaseInfo.FailureRate.Value = value.FailureRate
	// 		testCaseInfo.FailureRate.Type = "SINGLE_BAR"
	// 		testCaseInfo.FailureRate.ColorScheme = append(testCaseInfo.FailureRate.ColorScheme, struct {
	// 			Color0 string `json:"color0"`
	// 			Color1 string `json:"color1"`
	// 		}{
	// 			Color0: "#009C5B",
	// 			Color1: "#62CA9D",
	// 		})
	// 		testCaseInfo.FailureRate.ColorScheme = append(testCaseInfo.FailureRate.ColorScheme, struct {
	// 			Color0 string `json:"color0"`
	// 			Color1 string `json:"color1"`
	// 		}{
	// 			Color0: "#D32227",
	// 			Color1: "#FB6E72",
	// 		})
	// 		testCaseInfo.FailureRate.LightColorScheme = append(testCaseInfo.FailureRate.LightColorScheme, struct {
	// 			Color0 string `json:"color0"`
	// 			Color1 string `json:"color1"`
	// 		}{
	// 			Color0: "#0C9E61",
	// 			Color1: "#79CAA8",
	// 		})
	// 		testCaseInfo.FailureRate.LightColorScheme = append(testCaseInfo.FailureRate.LightColorScheme, struct {
	// 			Color0 string `json:"color0"`
	// 			Color1 string `json:"color1"`
	// 		}{
	// 			Color0: "#E83D39",
	// 			Color1: "#F39492",
	// 		})
	// 		// Populate the Data array
	// 		testCaseInfo.FailureRate.Data = append(testCaseInfo.FailureRate.Data, struct {
	// 			Title string `json:"title"`
	// 			Value int    `json:"value"`
	// 		}{
	// 			Title: "Successful runs",
	// 			Value: value.SuccessCount,
	// 		})
	// 		testCaseInfo.FailureRate.Data = append(testCaseInfo.FailureRate.Data, struct {
	// 			Title string `json:"title"`
	// 			Value int    `json:"value"`
	// 		}{
	// 			Title: "Failed runs",
	// 			Value: value.FailureCount,
	// 		})
	// 		// Safely append to the slice
	// 		mu.Lock()
	// 		testCasesOverview = append(testCasesOverview, testCaseInfo)
	// 		mu.Unlock()
	// 	}(value)
	// }

	// // Wait for all goroutines to finish
	// wg.Wait()

	// // Sort the slice
	// sort.Slice(testCasesOverview, func(i, j int) bool {
	// 	return testCasesOverview[i].LastRunInMillis > testCasesOverview[j].LastRunInMillis
	// })

	for _, value := range result.Aggregations.TestCasesOverview.Value {
		testCaseInfo := TestCasesOverview{}
		testCaseInfo.ComponentName = value.ComponentName
		testCaseInfo.TestSuiteName = value.TestSuiteName
		workflowName, sourceName, err := cutils.GetDisplayNameAndOrigin(value.AutomationName)
		if log.CheckErrorf(err, "Error getting display name and origin for updateTestCasesOverviewResponse()", err) {
			log.Infof("Error getting display name and origin for updateTestCasesOverviewResponse()")
		}
		testCaseInfo.Workflow = workflowName
		testCaseInfo.Source = sourceName
		testCaseInfo.TestCaseName = value.TestCaseName
		testCaseInfo.Branch = value.BranchName
		testCaseInfo.LastRun = value.StartTime
		testCaseInfo.LastRunInMillis = int(value.StartTimeInMillis)
		// testSuiteInfo.AvgRunTime = millisecondsToSecondsString(value.AverageDuration)
		testCaseInfo.AvgRunTime = value.AverageDuration
		testCaseInfo.TotalRuns.Value = value.Runs
		testCaseInfo.TotalRuns.DrillDown.ReportId = "test-overview-view-run-activity"
		testCaseInfo.TotalRuns.DrillDown.ReportTitle = "Runs - " + value.TestCaseName
		testCaseInfo.TotalRuns.DrillDown.ReportInfo = pb.ReportInfo{
			TestSuiteName: value.TestSuiteName,
			TestCaseName:  value.TestCaseName,
			AutomationId:  value.AutomationID,
			Branch:        value.BranchID,
			ComponentId:   value.ComponentID,
			WorkflowName:  workflowName,
			Source:        sourceName,
		}
		testCaseInfo.TotalRunsValue = value.Runs
		testCaseInfo.FailureRate.Value = value.FailureRate

		failureRateFloat := helper.ConvertPercentageToFloat(value.FailureRate)
		testCaseInfo.FailureRateValue = failureRateFloat

		testCaseInfo.FailureRate.Type = "SINGLE_BAR"
		testCaseInfo.FailureRate.ColorScheme = append(testCaseInfo.FailureRate.ColorScheme, struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			Color0: "#009C5B",
			Color1: "#62CA9D",
		})
		testCaseInfo.FailureRate.ColorScheme = append(testCaseInfo.FailureRate.ColorScheme, struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			Color0: "#D32227",
			Color1: "#FB6E72",
		})
		testCaseInfo.FailureRate.LightColorScheme = append(testCaseInfo.FailureRate.LightColorScheme, struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			Color0: "#0C9E61",
			Color1: "#79CAA8",
		})
		testCaseInfo.FailureRate.LightColorScheme = append(testCaseInfo.FailureRate.LightColorScheme, struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			Color0: "#E83D39",
			Color1: "#F39492",
		})
		// Populate the Data array
		testCaseInfo.FailureRate.Data = append(testCaseInfo.FailureRate.Data, struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		}{
			Title: "Successful runs",
			Value: value.SuccessCount,
		})
		testCaseInfo.FailureRate.Data = append(testCaseInfo.FailureRate.Data, struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		}{
			Title: "Failed runs",
			Value: value.FailureCount,
		})
		testCasesOverview = append(testCasesOverview, testCaseInfo)
	}

	sort.Slice(testCasesOverview, func(i, j int) bool {
		return testCasesOverview[i].LastRunInMillis > testCasesOverview[j].LastRunInMillis
	})

	return testCasesOverview, nil
}

func millisecondsToSecondsString(milliseconds float64) string {
	seconds := milliseconds / 1000.0
	secondsString := strconv.FormatFloat(seconds, 'f', -1, 64) // Convert float to string
	return secondsString + "s"
}
