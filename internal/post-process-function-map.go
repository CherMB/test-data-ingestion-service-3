package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	pb "github.com/calculi-corp/api/go/vsm/report"
	cutils "github.com/calculi-corp/common/utils"
	"github.com/calculi-corp/log"
	"github.com/calculi-corp/reports-service/constants"
	"github.com/calculi-corp/reports-service/db"
	"github.com/calculi-corp/reports-service/exceptions"
	"github.com/calculi-corp/reports-service/helper"
)

// PostProcess Function names referred in Widget Defintion are the Keys in PostProcessFunctionMap
// To add new entry, provide reference function name and its corresponding method
var PostProcessFunctionMap = map[string]interface{}{
	"Merged Default Branch Commits Header":  mergedDefaultBranchCommitsHeader,
	"Merged Default Branch Commits Section": mergedDefaultBranchCommitsSection,

	"Total Builds Header":  totalBuildsHeader,
	"Total Builds Section": totalBuildsSection,
	"Deployments Header":   deploymentsHeader,
	"Deployments Section":  deploymentsSection,

	"Average Development Cycle Time Header": averageDevelopmentCycleTimeHeader,
	"Development Cycle Chart Section":       developmentCycleChartSection,
	"Coding Time Footer Section":            codingTimeFooterSection,
	"Trivy Licenses Overview Section":       trivyLicenseOverviewSection,

	"component commits activity":            getCommitsActivity,
	"component runs activity":               getRunsActivity,
	"component builds metric":               getBuildsMetric,
	"component deployments":                 getDeployments,
	"component latest test results section": transformSummaryLatestTestResultsSection,

	"commit trends": getCommitTrends,

	"pull requests":                 getPullRequests,
	"successful build duration":     transformSuccessfulBuildDuration,
	"automation runs":               getAutomationRuns,
	"automation runs with scanner":  getAutomationRunsWithScanners,
	"vulnerabilities overview":      getVulnerabilitiesOverview,
	"open vulnerabilities overview": getOpenVulnerabilitiesOverview,
	"vulnerabilities by scan type":  getVulnerabilitiesByScanType,
	"sla status overview":           getSlaStatusOverview,
	"scan types in automation":      GetScanTypesInAutomation,
	"mttr for vulnerabilities":      getMttrForVulnerabilities,

	"work load":              getWorkload,
	"cycle time":             getCycleTime,
	"work efficiency":        getWorkEfficiency,
	"velocity":               getVelocity,
	"work item distribution": getWorkItemDistribution,

	"deployment frequency":               getDeploymentFrequency,
	"deployment lead time":               getDeploymentLeadTime,
	"failure rate":                       getFailureRate,
	"mean time to recovery":              getMttr,
	"deployment frequency and lead time": getDeploymentFrequencyAndLeadTime,
	"failure rate and mttr":              getFailureRateAndMttr,
	"code churn":                         getCodeChurn,

	"average deployment time": getAverageDeploymentTime,

	"automation runs for test suites": getAutomationRunsForTestSuites,

	"CWETM Top 25 Vulnerabilities":                            getCwetmTop25Vulnerabilities,
	"velocity component comparison":                           getVelocityComponentComparison,
	"cycle time component comparison":                         getCycleTimeComponentComparison,
	"commit trends component comparison":                      getCommitTrendsComponentComparison,
	"pull requests component comparison":                      getPullRequestComponentComparison,
	"workflow runs component comparison":                      getWorkflowRunsComponentComparison,
	"development cycle component comparison":                  getDevCycleTimeComponentComparison,
	"commits component comparison":                            getCommitsComponentComparison,
	"builds component comparison":                             getBuildsComponentComparison,
	"deployments component comparison":                        getDeploymentsComponentComparison,
	"components component comparison":                         getComponentsComponentComparison,
	"workflows component comparison":                          getWorkflowsComponentComparison,
	"active work time component comparison":                   getActiveFlowTimeComponentComparison,
	"work wait time component comparison":                     getWorkWaitTimeComponentComparison,
	"workload component comparison":                           getWorkloadComponentComparison,
	"vulnerabilities overview comparison":                     getVulnerabilitiesOverviewComponentComparison,
	"open vulnerabilities overview comparison":                getOpenVulnerabilitiesOverviewComponentComparison,
	"security workflow runs component comparison":             getSecurityWorkflowRunsComponentComparison,
	"security components component comparison":                getSecurityComponentsComponentComparison,
	"security workflow component comparison":                  getSecurityWorkflowsComponentComparison,
	"mttr very high component comparison":                     getMttrVeryHighComponentComparison,
	"mttr high component comparison":                          getMttrHighComponentComparison,
	"mttr medium component comparison":                        getMttrMediumComponentComparison,
	"mttr low component comparison":                           getMttrLowComponentComparison,
	"sast vulnerabilities scannner component comparison":      getSastVulnerabilitiesComponentComparison,
	"dast vulnerabilities scannner component comparison":      getDastVulnerabilitiesComponentComparison,
	"container vulnerabilities scannner component comparison": getContainerVulnerabilitiesComponentComparison,
	"sca vulnerabilities scannner component comparison":       getScaVulnerabilitiesComponentComparison,
	"deployment frequency component comparison":               getDeploymentFrequencyComponentComparison,
	"mean time to recovery component comparison":              getDoraMttrComponentComparison,
	"deployment lead time component comparison":               getDeploymentLeadTimeComponentComparison,
	"failure rate component comparison":                       getFailureRateComponentComparison,
	"test workflow component comparison":                      transformTestWorkflowsComponentComparison,
	"test suite workflow runs component comparison":           getTestSuiteWorkflowRunsComponentComparison,
	"test components component comparison":                    getTestComponentsComponentComparison,

	"findings remediation trend":                  transformFindingsRemediationTrend,
	"open findings by severity":                   transformOpenFindingsBySeverity,
	"open findings by security tool":              transformOpenFindingsBySecurityTool,
	"open findings distribution by category":      transformOpenFindingsDistributionByCategory,
	"sla breached by asset":                       transformSlasBreachedByAsset,
	"open findings distribution by security tool": transformOpenFindingsDistributionBySecurityTool,
	"findings identified since":                   transformFindingsIdentfiedSince,
	"risk accepted and false positive findings":   transformRiskAcceptedAndFalsePositiveFindings,
	"sla breaches by severity":                    transformSlaBreachesBySeverity,
	"open findings by SLA status":                 transformOpenFindingsBySlaStatus,
	"open findings by review status":              transformOpenFindingsByReviewStatus,
	"open findings by component":                  transformOpenFindingsByComponent,
}

// Widget Builder execute the function to get data for the widget
// Ensure the query keys used in Widget defintion is defined and match, as it is referred in the Function
func ExecutePostProcessFunction(k string, specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	// For any other function signature use the switch case
	switch k {
	default:
		f, ok := PostProcessFunctionMap[k]
		if !ok || f == nil {
			return nil, fmt.Errorf("post process function not found for %s", k)
		}
		res, err := f.(func(string, map[string]json.RawMessage, map[string]any) (json.RawMessage, error))(specKey, data, replacements)
		return res, err
	}
}

func ExecutePostProcessFunctionForComponentComparison(k string, specKey string, data map[string]json.RawMessage, replacements map[string]any, organization *constants.Organization) (json.RawMessage, error) {
	// For any other function signature use the switch case
	switch k {
	default:
		f, ok := PostProcessFunctionMap[k]
		if !ok || f == nil {
			return nil, fmt.Errorf("post process function not found for %s", k)
		}
		res, err := f.(func(string, map[string]json.RawMessage, map[string]any, *constants.Organization) (json.RawMessage, error))(specKey, data, replacements, organization)
		return res, err
	}
}

// Below Functions associated to the Analytics/Software_Delivery_Activity/Code_Progression_Snapshot_Section

func mergedDefaultBranchCommitsHeader(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	response, ok := data["automationRunsCount"]
	if !ok {
		return nil, db.ErrInternalServer
	}
	totalCount := 0.0
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.AUTOMATION_RUN] != nil {
			automationRuns := aggsResult[constants.AUTOMATION_RUN].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				if values[constants.DATA] != nil {
					count := values[constants.TOTAL_COUNT].(float64)
					totalCount += count
				}
			}
		}
	}
	outputResponse := make(map[string]interface{})
	outputResponse[constants.VALUE] = totalCount
	responseJson, err := json.Marshal(outputResponse)
	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}

}

func mergedDefaultBranchCommitsSection(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	response, ok := data["deployedAutomationCount"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	runQueryResponse, ok := data["automationRunsCount"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	resultObject := make(map[string]interface{})
	automationRunKey := make(map[string]interface{})
	json.Unmarshal([]byte(runQueryResponse), &resultObject)
	if resultObject[constants.AGGREGATION] != nil {
		aggsResult := resultObject[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.AUTOMATION_RUN] != nil {
			automationRuns := aggsResult[constants.AUTOMATION_RUN].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				if values[constants.DATA] != nil {
					data := values[constants.DATA].(map[string]interface{})
					if len(data) > 0 {
						automationRunKey = data
					}
				}
			}
		}
	}
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	envCountMap := make(map[string]int)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.AUTOMATION_RUN] != nil {
			automationRuns := aggsResult[constants.AUTOMATION_RUN].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				if values[constants.DATA] != nil {
					data := values[constants.DATA].(map[string]interface{})
					for key, value := range data {
						envMap := value.(map[string]interface{})
						for envKey, envValue := range envMap {
							val, ok := envCountMap[envKey]
							if ok {
								envCountMap[envKey] = val + int(envValue.(float64))
							} else {
								envCountMap[envKey] = int(envValue.(float64))
							}
						}
						_, ok := automationRunKey[key]
						if ok {
							delete(automationRunKey, key)
						}
					}
				}
			}
		}
	}
	outputResponse := make(map[string]interface{})
	dataMapList := make([]map[string]interface{}, 0)
	outData := make(map[string]interface{})
	outData["x"] = "Unspecified"
	outData["y"] = len(automationRunKey)
	colorList := make([]map[string]interface{}, 0)
	colorMap := make(map[string]interface{})
	colorMap["color0"] = "#FCC26C"
	colorMap["color1"] = "#FF8307"
	colorList = append(colorList, colorMap)
	outData["colorScheme"] = colorList
	colorList = make([]map[string]interface{}, 0)
	colorMap = make(map[string]interface{})
	colorMap["color0"] = "#FCC26C"
	colorMap["color1"] = "#FF8307"
	colorList = append(colorList, colorMap)
	outData["lightColorScheme"] = colorList
	dataMapList = append(dataMapList, outData)
	keys := make([]string, 0, len(envCountMap))
	for key := range envCountMap {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return envCountMap[keys[i]] < envCountMap[keys[j]]
	})
	for _, key := range keys {
		value := envCountMap[key]
		dataMap := make(map[string]interface{})
		dataMap["x"] = key
		dataMap["y"] = value
		dataMapList = append(dataMapList, dataMap)
	}
	outputResponse[constants.DATA] = dataMapList
	outputResponse[constants.ID] = "Run-initiating commits"

	responseJson, err := json.Marshal([]any{outputResponse})
	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}
}

func totalBuildsHeader(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	response, ok := data["totalBuildsHeader"]
	if !ok {
		return nil, db.ErrInternalServer
	}
	totalBuildsCount := 0.0
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult["total_builds"] != nil {
			totalBuilds := aggsResult["total_builds"].(map[string]interface{})
			if totalBuilds["value"] != nil {
				totalBuildsCount = totalBuilds["value"].(float64)
			}
		}
	}
	outputResponse := make(map[string]interface{})
	outputResponse[constants.VALUE] = totalBuildsCount
	responseJson, err := json.Marshal(outputResponse)

	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}
}

func totalBuildsSection(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type buildSection struct {
		Aggregations struct {
			BuildStatus struct {
				Value struct {
					Data []struct {
						Name  string `json:"name"`
						Value int    `json:"value"`
					} `json:"data"`
					Info []struct {
						DrillDown struct {
							ReportType  string `json:"reportType"`
							ReportID    string `json:"reportId"`
							ReportTitle string `json:"reportTitle"`
						} `json:"drillDown"`
						Title string `json:"title"`
						Value int    `json:"value"`
					} `json:"info"`
				} `json:"value"`
			} `json:"build_status"`
		} `json:"aggregations"`
	}

	result := buildSection{}
	buildsDataResponse := data["buildsData"]
	err := json.Unmarshal([]byte(buildsDataResponse), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in buildsDataResponse()") {
		return nil, err
	}

	b, err := json.Marshal(result.Aggregations.BuildStatus.Value)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getAutomationRunsWithScanners()") {
		return nil, err
	}
	return b, nil

}

func deploymentsHeader(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	response, ok := data["deploymentsHeader"]
	if !ok {
		return nil, db.ErrInternalServer
	}
	deployCount := 0.0
	result := make(map[string]interface{})
	json.Unmarshal([]byte(response), &result)
	if result[constants.AGGREGATION] != nil {
		aggsResult := result[constants.AGGREGATION].(map[string]interface{})
		if aggsResult["deploy_count"] != nil {
			deployments := aggsResult["deploy_count"].(map[string]interface{})
			if deployments[constants.VALUE] != nil {
				deployCount = deployments[constants.VALUE].(float64)
			}
		}
	}
	outputResponse := make(map[string]interface{})
	outputResponse[constants.VALUE] = deployCount
	responseJson, err := json.Marshal(outputResponse)

	if err != nil {
		return nil, err
	} else {
		return responseJson, nil
	}

}

func deploymentsSection(specKey string, data map[string]json.RawMessage, replacements map[string]interface{}) (json.RawMessage, error) {
	type deploySectionResponse struct {
		Aggregations struct {
			Deploys struct {
				Value struct {
					Data []struct {
						X string `json:"x"`
						Y int    `json:"y"`
					} `json:"data"`
				} `json:"value"`
			} `json:"deploys"`
		} `json:"aggregations"`
	}

	// Unmarshal the response data into the deploySectionResponse struct
	var result deploySectionResponse
	if err := json.Unmarshal(data["envDeploymentInfo"], &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response from OpenSearch envDeploymentInfo: %v", err)
	}

	responseWithID := struct {
		Data []struct {
			X string `json:"x"`
			Y int    `json:"y"`
		} `json:"data"`
		ID string `json:"id"`
	}{
		Data: result.Aggregations.Deploys.Value.Data,
		ID:   "Successful deployments",
	}
	responseArray := []interface{}{responseWithID}
	responseJSON, err := json.Marshal(responseArray)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response: %v", err)
	}
	return responseJSON, nil
}

// Below Functions associated to the Analytics/Software_Delivery_Activity/Development_Cycle_Section
func averageDevelopmentCycleTimeHeader(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	response, ok := data["avgDevelopmentHeader"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	type DevCycle struct {
		Aggregations struct {
			DevelopmentCycleTime struct {
				Value struct {
					Total string `json:"total"`
					Value int    `json:"value"`
				} `json:"value"`
			} `json:"developmentCycleTime"`
		} `json:"aggregations"`
	}

	result := DevCycle{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in trivyLicenseOverviewSection()") {
		return nil, err
	}

	b := struct {
		Value         string `json:"value"`
		ValueInMillis int    `json:"valueInMillis"`
	}{result.Aggregations.DevelopmentCycleTime.Value.Total, result.Aggregations.DevelopmentCycleTime.Value.Value}

	output, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func developmentCycleChartSection(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	response, ok := data["developmentCycleChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	type DevCycleChart struct {
		Aggregations struct {
			DevelopmentCycleTime struct {
				Value struct {
					CodingTime float64 `json:"coding_time"`
					ReviewTime float64 `json:"review_time"`
					PickupTime float64 `json:"pickup_time"`
					DeployTime float64 `json:"deploy_time"`
				} `json:"value"`
			} `json:"developmentCycleTime"`
		} `json:"aggregations"`
	}

	result := DevCycleChart{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in developmentCycleChartSection()") {
		return nil, err
	}

	b := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: "Coding time",
			Value: int(result.Aggregations.DevelopmentCycleTime.Value.CodingTime),
		}, {
			Title: "Code pickup time",
			Value: int(result.Aggregations.DevelopmentCycleTime.Value.PickupTime),
		}, {
			Title: "Code review time",
			Value: int(result.Aggregations.DevelopmentCycleTime.Value.ReviewTime),
		},
	}

	rawJson, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return rawJson, nil
}

func codingTimeFooterSection(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	response, ok := data["developmentTimeFooter"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	type codingTime struct {
		Aggregations struct {
			DevelopmentCycleTime struct {
				Value struct {
					CodingTime         string `json:"coding_time"`
					PickupTimeInMillis int    `json:"pickup_time_in_millis"`
					DeployTimeInMillis int    `json:"deploy_time_in_millis"`
					CodingTimeInMillis int    `json:"coding_time_in_millis"`
					ReviewTime         string `json:"review_time"`
					PickupTime         string `json:"pickup_time"`
					ReviewTimeInMillis int    `json:"review_time_in_millis"`
					DeployTime         string `json:"deploy_time"`
				} `json:"value"`
			} `json:"developmentCycleTime"`
		} `json:"aggregations"`
	}

	result := codingTime{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in codingTimeFooterSection()") {
		return nil, err
	}

	if specKey == "codingTimeSpec" {
		b := struct {
			Value         string `json:"value"`
			ValueInMillis int    `json:"valueInMillis"`
		}{result.Aggregations.DevelopmentCycleTime.Value.CodingTime, result.Aggregations.DevelopmentCycleTime.Value.CodingTimeInMillis}

		output, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}
		return output, nil

	} else if specKey == "codingPickupTimeSpec" {
		b := struct {
			Value         string `json:"value"`
			ValueInMillis int    `json:"valueInMillis"`
		}{result.Aggregations.DevelopmentCycleTime.Value.PickupTime, result.Aggregations.DevelopmentCycleTime.Value.PickupTimeInMillis}

		output, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}

		return output, nil

	} else if specKey == "codingReviewTimeSpec" {
		b := struct {
			Value         string `json:"value"`
			ValueInMillis int    `json:"valueInMillis"`
		}{result.Aggregations.DevelopmentCycleTime.Value.ReviewTime, result.Aggregations.DevelopmentCycleTime.Value.ReviewTimeInMillis}

		output, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}

		return output, nil
	}
	return nil, nil
}

// function to transform section data fetched for the Trivy License Overview widget
func trivyLicenseOverviewSection(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	response, ok := data["trivyLicensesOverviewSection"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in trivyLicenseOverviewSection()") {
		return nil, err
	}
	if result[constants.AGGREGATION] != nil {
		x := result[constants.AGGREGATION].(map[string]interface{})
		if x[constants.TRIVY_LICENSE_SECTION] != nil {
			y := x[constants.TRIVY_LICENSE_SECTION].(map[string]interface{})
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

// function to transform data fetched for the Component Summary Commits Activity
func getCommitsActivity(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	response, ok := data["commitsActivity"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	type commitsActivity struct {
		Aggregations struct {
			Commits struct {
				Value struct {
					Avg struct {
						Title string `json:"title"`
						Value int    `json:"value"`
					} `json:"avg"`
					Dev struct {
						Title string `json:"title"`
						Value int    `json:"value"`
					} `json:"dev"`
					CommitsCount int `json:"commits_count"`
				} `json:"value"`
			} `json:"commits"`
		} `json:"aggregations"`
	}

	result := commitsActivity{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getCommitsActivity()") {
		return nil, err
	}

	if specKey == "header" {
		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.Commits.Value.CommitsCount) + `}`)
		return b, nil
	} else if specKey == "section" {
		list := []any{}
		list = append(list, result.Aggregations.Commits.Value.Avg)
		list = append(list, result.Aggregations.Commits.Value.Dev)

		b, err := json.Marshal(list)
		if log.CheckErrorf(err, "Error unmarshaling response for section in getCommitsActivity()") {
			return nil, err
		}

		return b, err
	}

	return nil, nil
}

// function to transform data fetched for the Component Summary Runs Activity
func getRunsActivity(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	response, ok := data["runsActivity"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	type runsActivity struct {
		Aggregations struct {
			AutomationRun struct {
				Value struct {
					Data []struct {
						Title string `json:"title"`
						Value int    `json:"value"`
					} `json:"data"`
					TotalCount int `json:"totalCount"`
				} `json:"value"`
			} `json:"automation_run"`
		} `json:"aggregations"`
	}

	result := runsActivity{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getRunsActivity()") {
		return nil, err
	}

	if specKey == "header" {
		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.AutomationRun.Value.TotalCount) + `}`)
		return b, nil
	} else if specKey == "section" {
		output, err := json.Marshal(result.Aggregations.AutomationRun.Value.Data)
		if log.CheckErrorf(err, "Error unmarshaling section data in getRunsActivity()") {
			return nil, err
		}

		return output, nil
	}

	return nil, nil
}

// function to transform data fetched for the Component Summary Runs Activity
func getBuildsMetric(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type buildsData struct {
		Aggregations struct {
			BuildStatus struct {
				Value struct {
					Info []struct {
						DrillDown struct {
							ReportType  string `json:"reportType"`
							ReportID    string `json:"reportId"`
							ReportTitle string `json:"reportTitle"`
						} `json:"drillDown"`
						Title string `json:"title"`
						Value int    `json:"value"`
					} `json:"info"`
					TotalBuilds int `json:"total_builds"`
				} `json:"value"`
			} `json:"build_status"`
		} `json:"aggregations"`
	}

	type buildChartStruct struct {
		ID   string `json:"id"`
		Data []struct {
			X string `json:"x"`
			Y int    `json:"y"`
		} `json:"data"`
	}

	type buildChartDataStruct struct {
		Aggregations struct {
			RunsBuckets struct {
				Buckets []struct {
					KeyAsString   string `json:"key_as_string"`
					Key           int64  `json:"key"`
					DocCount      int    `json:"doc_count"`
					AutomationRun struct {
						Value struct {
							Success int `json:"Success"`
							Failure int `json:"Failure"`
						} `json:"value"`
					} `json:"automation_run"`
				} `json:"buckets"`
			} `json:"runs_buckets"`
		} `json:"aggregations"`
	}

	if specKey == "subHeader" || specKey == "header" {
		response, ok := data["buildsData"]
		if !ok {
			return nil, db.ErrInternalServer
		}
		result := buildsData{}
		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response buildsData in getBuildsMetric()") {
			return nil, err
		}

		if specKey == "header" {
			b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.BuildStatus.Value.TotalBuilds) + `}`)
			return b, nil
		} else if specKey == "subHeader" {

			outputResponse := make(map[string]interface{})
			outputResponse["subHeader"] = result.Aggregations.BuildStatus.Value.Info
			rawJson, err := json.Marshal(outputResponse)
			if err != nil {
				return nil, err
			}
			return rawJson, nil
		}
	} else if specKey == "sectionChart" {
		response, ok := data["buildsDataChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := buildChartDataStruct{}
		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response buildsDataChar in getBuildsMetric()") {
			return nil, err
		}

		successStruct := buildChartStruct{}
		successStruct.ID = "Success"
		failureStruct := buildChartStruct{}
		failureStruct.ID = "Failure"

		for _, v := range result.Aggregations.RunsBuckets.Buckets {

			sd := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: v.KeyAsString, Y: v.AutomationRun.Value.Success}

			fd := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: v.KeyAsString, Y: v.AutomationRun.Value.Failure}

			successStruct.Data = append(successStruct.Data, sd)
			failureStruct.Data = append(failureStruct.Data, fd)
		}

		output := []buildChartStruct{successStruct, failureStruct}
		b, err := json.Marshal(output)
		if log.CheckErrorf(err, "Error marshaling buildChartStruct in getBuildsMetric()") {
			return nil, err
		}
		return b, nil
	}

	return nil, nil
}

func getDeployments(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "deploymentSuccessRateHeaderSpec" || specKey == "deploymentSuccessRateSubHeaderSpec" {
		response, ok := data["deploymentSuccessRateHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type deploymentSuccessRateHeader struct {
			Aggregations struct {
				DeployData struct {
					Value struct {
						Total int `json:"total"`
						Data  []struct {
							Title string `json:"title"`
							Value int    `json:"value"`
						} `json:"data"`
						Value string `json:"value"`
					} `json:"value"`
				} `json:"deploy_data"`
			} `json:"aggregations"`
		}

		result := deploymentSuccessRateHeader{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getDeployments()") {
			return nil, err
		}

		if specKey == "deploymentSuccessRateHeaderSpec" {

			total := "0%"
			if result.Aggregations.DeployData.Value.Value != "" {
				total = result.Aggregations.DeployData.Value.Value
			}
			response := map[string]string{
				"value": total,
			}
			b, err := json.Marshal(response)
			if err != nil {
				return nil, err
			}
			return b, nil
		} else if specKey == "deploymentSuccessRateSubHeaderSpec" {

			type item struct {
				Value int    `json:"value"`
				Title string `json:"title"`
			}

			v1 := item{}
			v2 := item{}

			v1.Title = "Success"
			v2.Title = "Failure"

			if len(result.Aggregations.DeployData.Value.Data) > 0 {
				v1.Value = result.Aggregations.DeployData.Value.Data[0].Value
				v2.Value = result.Aggregations.DeployData.Value.Data[1].Value
			} else {
				v1.Value = 0
				v2.Value = 0
			}

			dataMapList := []map[string]interface{}{
				{
					"title": v1.Title,
					"value": v1.Value,
				},
				{
					"title": v2.Title,
					"value": v2.Value,
				},
			}
			outputResponse := make(map[string]interface{})
			outputResponse["subHeader"] = dataMapList
			rawJson, err := json.Marshal(outputResponse)
			if err != nil {
				return nil, err
			}
			return rawJson, nil
		}
	} else if specKey == "deploymentDataSpec" {
		response, ok := data["deploymentData"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetDeploymentTypeAssertion), "aggrBy is not a string in getDeployments()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetDeploymentTypeAssertion), "startDate is not a string in getDeployments()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getDeployments"), "endDate is not a string in getDeployments()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getDeployments(): ") {
			return nil, db.ErrInternalServer
		}

		type deploymentData struct {
			Aggregations struct {
				DeployBuckets struct {
					Buckets []struct {
						KeyAsString string `json:"key_as_string"`
						Key         int64  `json:"key"`
						DocCount    int    `json:"doc_count"`
						DeployData  struct {
							Value struct {
								Success int `json:"Success"`
								Failure int `json:"Failure"`
							} `json:"value"`
						} `json:"deploy_data"`
					} `json:"buckets"`
				} `json:"deploy_buckets"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}
		result := deploymentData{}
		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getDeployments()") {
			return nil, err
		}

		v1 := responseStruct{}
		v2 := responseStruct{}

		v1.ID = "Success"
		v2.ID = "Failure"

		for index, v := range result.Aggregations.DeployBuckets.Buckets {
			startDate := v.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			success := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(v.DeployData.Value.Success)}

			failure := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: v.DeployData.Value.Failure}

			v1.Data = append(v1.Data, success)
			v2.Data = append(v2.Data, failure)

		}

		outputStruct := []responseStruct{v1, v2}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getDeployments() chart") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getDeployments()") {
			return nil, err
		}
		return output, nil
	}
	return nil, nil
}

// transformSummaryLatestTestResultsSection transforms OpenSearch response for the Latest Test Results widget (cs11) in the Summary page into the API contract's structure.
func transformSummaryLatestTestResultsSection(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "summaryLatestTestResultsSectionSpec" {
		response, ok := data["summaryLatestTestResultsSection"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type SummaryLatestTestResultsInput struct {
			Aggregations struct {
				Workflows struct {
					Buckets []struct {
						Key    string `json:"key"`
						Suites struct {
							Buckets []struct {
								Key            string `json:"key"`
								LatestSuiteDoc struct {
									Hits struct {
										Hits []struct {
											Source struct {
												Duration     float64 `json:"duration"`
												Total        int     `json:"total"`
												RunId        string  `json:"run_id"`
												RunNumber    int     `json:"run_number"`
												Passed       int     `json:"passed"`
												Failed       int     `json:"failed"`
												Skipped      int     `json:"skipped"`
												WorkflowName string  `json:"automation_name"`
												Source       string  `json:"source"`
											} `json:"_source"`
											Fields struct {
												ZonedStartTime    []string  `json:"zoned_start_time"`
												StartTimeInMillis []float64 `json:"start_time_in_millis"`
											} `json:"fields"`
										} `json:"hits"`
									} `json:"hits"`
								} `json:"latest_suite_doc"`
							} `json:"buckets"`
						} `json:"suites"`
					} `json:"buckets"`
				} `json:"workflows"`
			} `json:"aggregations"`
		}

		type SummaryLatestTestResultsOutput struct {
			TestSuiteName    string                  `json:"testSuiteName"`
			WorkflowName     string                  `json:"workflow"`
			LastRun          string                  `json:"lastRun"`
			LastRunInMillis  float64                 `json:"lastRunInMillis"`
			TotalTestCases   int                     `json:"totalTestCases"`
			TestCasesPassed  int                     `json:"testCasesPassed"`
			TestCasesFailed  int                     `json:"testCasesFailed"`
			TestCasesSkipped int                     `json:"testCasesSkipped"`
			RunTime          float64                 `json:"runTime"`
			Source           string                  `json:"source"`
			DrillDown        DrillDownWithReportInfo `json:"drillDown"`
		}

		var SummaryLTRInput SummaryLatestTestResultsInput

		err := json.Unmarshal([]byte(response), &SummaryLTRInput)
		if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in transformSummaryLatestTestResultsSection: ") {
			return nil, err
		}

		var SummaryLTROutput []SummaryLatestTestResultsOutput

		if SummaryLTRInput.Aggregations.Workflows.Buckets != nil {
			for _, workflowBucket := range SummaryLTRInput.Aggregations.Workflows.Buckets {
				if workflowBucket.Suites.Buckets != nil {
					for _, suiteBucket := range workflowBucket.Suites.Buckets {
						hit := suiteBucket.LatestSuiteDoc.Hits.Hits[0]

						workflowName, sourceName, err := cutils.GetDisplayNameAndOrigin(hit.Source.WorkflowName)
						if log.CheckErrorf(err, "Error getting display name and origin for transformSummaryLatestTestResultsSection", err) {
							log.Infof("Error getting display name and origin for transformSummaryLatestTestResultsSection()")
						}
						SummaryLTROutput = append(SummaryLTROutput, SummaryLatestTestResultsOutput{
							TestSuiteName:    suiteBucket.Key,
							WorkflowName:     workflowName,
							TotalTestCases:   hit.Source.Total,
							TestCasesPassed:  hit.Source.Passed,
							TestCasesFailed:  hit.Source.Failed,
							TestCasesSkipped: hit.Source.Skipped,
							RunTime:          hit.Source.Duration,
							Source: func() string {
								if hit.Source.Source != "" {
									return hit.Source.Source
								}
								return "CloudBees"
							}(),
							DrillDown: DrillDownWithReportInfo{
								ReportId:    "latest-test-results",
								ReportType:  "status",
								ReportTitle: "Test cases - " + suiteBucket.Key,
								ReportInfo: pb.ReportInfo{
									RunId:         hit.Source.RunId,
									RunNumber:     strconv.Itoa(hit.Source.RunNumber),
									TestSuiteName: suiteBucket.Key,
									AutomationId:  workflowBucket.Key,
									WorkflowName:  workflowName,
									Source:        sourceName,
								},
							},
							LastRun:         hit.Fields.ZonedStartTime[0],
							LastRunInMillis: hit.Fields.StartTimeInMillis[0],
						})

					}
				}

			}
		}

		// Sort the list by LastRunInMillis in descending order
		sort.Slice(SummaryLTROutput, func(i, j int) bool {
			return SummaryLTROutput[i].LastRunInMillis > SummaryLTROutput[j].LastRunInMillis
		})

		output, err := json.Marshal(SummaryLTROutput)
		if err != nil {
			return nil, err
		}
		return output, nil

	}
	return nil, nil
}

// function to transform data fetched for the Commits Trend in SDA
func getCommitTrends(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "totalCommitsHeaderSpec" || specKey == "averageCommitsHeaderSpec" {
		response, ok := data["totalCommitsHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type commitsData struct {
			Aggregations struct {
				CommitsTrendsWidget struct {
					Buckets []struct {
						Key                 string `json:"key"`
						From                int64  `json:"from"`
						FromAsString        string `json:"from_as_string"`
						To                  int64  `json:"to"`
						ToAsString          string `json:"to_as_string"`
						DocCount            int    `json:"doc_count"`
						CommitsTrendHeaders struct {
							Value struct {
								UniqueAuthors    int     `json:"unique_authors"`
								CommitsCount     int     `json:"commits_count"`
								CommitsPerAuthor float64 `json:"commits-per-author"`
							} `json:"value"`
						} `json:"commits_trend_headers"`
					} `json:"buckets"`
				} `json:"commits_trends_widget"`
			} `json:"aggregations"`
		}

		result := commitsData{}
		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getCommitTrends()") {
			return nil, err
		}

		if specKey == "totalCommitsHeaderSpec" {

			if len(result.Aggregations.CommitsTrendsWidget.Buckets) > 0 {
				b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.CommitsTrendsWidget.Buckets[0].CommitsTrendHeaders.Value.CommitsCount) + `}`)
				return b, nil
			}
		} else if specKey == "averageCommitsHeaderSpec" {

			type item struct {
				Value int    `json:"value"`
				Title string `json:"title"`
			}

			v1 := item{}
			v2 := item{}

			v1.Title = "Active developers"

			t, ok := replacements["commitTitle"]
			if ok {
				v2.Title = t.(string)
			} else {
				v2.Title = "Weekly commits/active devs"
			}

			if len(result.Aggregations.CommitsTrendsWidget.Buckets) > 0 {
				v1.Value = result.Aggregations.CommitsTrendsWidget.Buckets[0].CommitsTrendHeaders.Value.UniqueAuthors
				v2.Value = int(result.Aggregations.CommitsTrendsWidget.Buckets[0].CommitsTrendHeaders.Value.CommitsPerAuthor)
			} else {
				v1.Value = 0
				v2.Value = 0
			}

			dataMapList := []map[string]interface{}{
				{
					"title": v1.Title,
					"value": v1.Value,
					"drillDown": map[string]interface{}{
						"reportType":  "status",
						"reportId":    "activeDevelopers",
						"reportTitle": "Active developers",
					},
				},
				{
					"title": v2.Title,
					"value": v2.Value,
				},
			}

			outputResponse := make(map[string]interface{})
			outputResponse["subHeader"] = dataMapList
			rawJson, err := json.Marshal(outputResponse)
			if err != nil {
				return nil, err
			}
			return rawJson, nil
		}
	} else if specKey == "commitsAndAverageChartSpec" {

		response, ok := data["commitsAndAverageChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetCommitTrendTypeAssertion), "aggrBy is not a string in getCommitTrends()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetCommitTrendTypeAssertion), "startDate is not a string in getCommitTrends()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getCommitTrends"), "endDate is not a string in getCommitTrends()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getCommitTrends(): ") {
			return nil, db.ErrInternalServer
		}

		type commitsDatastruct struct {
			Aggregations struct {
				CommitsTrendsWidget struct {
					Buckets []struct {
						KeyAsString  string `json:"key_as_string"`
						CommitsCount struct {
							Value int `json:"value"`
						} `json:"commits_count"`
						CommitsPerAuth struct {
							Value float64 `json:"value"`
						} `json:"commits-per-auth,omitempty"`
					} `json:"buckets"`
				} `json:"commits_trends_widget"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID             string `json:"id"`
			Type           string `json:"type"`
			IsClickDisable bool   `json:"isClickDisable,omitempty"`
			Data           []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		result := commitsDatastruct{}
		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getCommitTrends()") {
			return nil, err
		}

		v1 := responseStruct{}
		v2 := responseStruct{}

		t, ok := replacements["commitTitle"]
		if ok {
			v1.ID = t.(string)
		} else {
			v1.ID = "Weekly commits/active devs"
		}
		v2.ID = "Commits"

		v1.IsClickDisable = true
		v2.IsClickDisable = false

		v1.Type = "line"
		v2.Type = "line"

		for index, v := range result.Aggregations.CommitsTrendsWidget.Buckets {
			startDate := v.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			commitAuthor := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(v.CommitsPerAuth.Value)}

			commit := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: v.CommitsCount.Value}

			v1.Data = append(v1.Data, commitAuthor)
			v2.Data = append(v2.Data, commit)
		}

		outputStruct := []responseStruct{v2, v1}
		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getCommitTrends()") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getCommitTrends()") {
			return nil, err
		}

		return output, nil

	}

	return nil, nil
}

// function to transform data fetched for the automation runs in SDA
func getAutomationRuns(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "totalRunsHeaderSpec" || specKey == "totalRunsSubHeaderSpec" {
		response, ok := data["totalRunsSubHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type automationRuns struct {
			Aggregations struct {
				AutomationRun struct {
					Value struct {
						Data []struct {
							Title string `json:"title"`
							Value int    `json:"value"`
						} `json:"data"`
					} `json:"value"`
				} `json:"automation_run"`
			} `json:"aggregations"`
		}
		result := automationRuns{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getAutomationRuns()") {
			return nil, err
		}

		if specKey == "totalRunsHeaderSpec" {
			total := 0
			if len(result.Aggregations.AutomationRun.Value.Data) > 0 {
				total = result.Aggregations.AutomationRun.Value.Data[0].Value + result.Aggregations.AutomationRun.Value.Data[1].Value
			}

			b := []byte(`{"value":` + fmt.Sprint(total) + `}`)
			return b, nil
		} else if specKey == "totalRunsSubHeaderSpec" {

			type item struct {
				Value int    `json:"value"`
				Title string `json:"title"`
			}

			v1 := item{}
			v2 := item{}

			v1.Title = "Success"
			v2.Title = "Failure"

			if len(result.Aggregations.AutomationRun.Value.Data) > 0 {
				v1.Value = result.Aggregations.AutomationRun.Value.Data[0].Value
				v2.Value = result.Aggregations.AutomationRun.Value.Data[1].Value
			} else {
				v1.Value = 0
				v2.Value = 0
			}

			dataMapList := []map[string]interface{}{
				{
					"title": v1.Title,
					"value": v1.Value,
					"drillDown": map[string]interface{}{
						"reportId":    "workflowRuns",
						"reportTitle": constants.WORKFLOW_RUNS_REPORT_TITLE,
						"reportType":  "status",
					},
				},
				{
					"title": v2.Title,
					"value": v2.Value,
					"drillDown": map[string]interface{}{
						"reportId":    "workflowRuns",
						"reportTitle": constants.WORKFLOW_RUNS_REPORT_TITLE,
						"reportType":  "status",
					},
				},
			}

			outputResponse := make(map[string]interface{})
			outputResponse["subHeader"] = dataMapList
			rawJson, err := json.Marshal(outputResponse)
			if err != nil {
				return nil, err
			}
			return rawJson, nil
		}
	} else if specKey == "runsStatusChartSpec" {

		response, ok := data["runsStatusChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetAutomationRunTypeAssertion), "aggrBy is not a string in getAutomationRuns()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetAutomationRunTypeAssertion), "startDate is not a string in getAutomationRuns()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getAutomationRuns"), "endDate is not a string in getAutomationRuns()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getAutomationRuns(): ") {
			return nil, db.ErrInternalServer
		}

		type runStatus struct {
			Aggregations struct {
				RunsBuckets struct {
					Buckets []struct {
						KeyAsString   string `json:"key_as_string"`
						Key           int64  `json:"key"`
						DocCount      int    `json:"doc_count"`
						AutomationRun struct {
							Value struct {
								Success int `json:"Success"`
								Failure int `json:"Failure"`
							} `json:"value"`
						} `json:"automation_run"`
					} `json:"buckets"`
				} `json:"runs_buckets"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		result := runStatus{}
		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getAutomationRuns()") {
			return nil, err
		}

		v1 := responseStruct{}
		v2 := responseStruct{}

		v1.ID = "Success"
		v2.ID = "Failure"

		for index, v := range result.Aggregations.RunsBuckets.Buckets {
			startDate := v.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			success := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(v.AutomationRun.Value.Success)}

			failure := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: v.AutomationRun.Value.Failure}

			v1.Data = append(v1.Data, success)
			v2.Data = append(v2.Data, failure)

		}

		outputStruct := []responseStruct{v1, v2}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getAutomationRuns() chart") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getAutomationRuns()") {
			return nil, err
		}
		return output, nil

	}

	return nil, nil
}

// function to transform data fetched for the automation runs in SDA
func getPullRequests(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "totalPullRequestsHeaderSpec" {
		response, ok := data["totalPullRequestsHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type pullRequests struct {
			Aggregations struct {
				ByRepos struct {
					DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
					SumOtherDocCount        int `json:"sum_other_doc_count"`
					Buckets                 []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
						PrCount  struct {
							Value int `json:"value"`
						} `json:"pr_count"`
					} `json:"buckets"`
				} `json:"by_repos"`
				SumPrs struct {
					Value int `json:"value"`
				} `json:"sum_prs"`
			} `json:"aggregations"`
		}

		result := pullRequests{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getPullRequests()") {
			return nil, err
		}

		if specKey == "totalPullRequestsHeaderSpec" {
			b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.SumPrs.Value) + `}`)
			return b, nil

		}
	} else if specKey == "pullRequestsChartSpec" {

		response, ok := data["pullRequestsChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetPullRequestTypeAssertion), "aggrBy is not a string in getPullRequests()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetPullRequestTypeAssertion), "startDate is not a string in getPullRequests()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getPullRequests"), "endDate is not a string in getPullRequests()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getPullRequests(): ") {
			return nil, db.ErrInternalServer
		}

		type pullRequestBuckets struct {
			Aggregations struct {
				DateCounts struct {
					Buckets []struct {
						KeyAsString string `json:"key_as_string"`
						Key         int64  `json:"key"`
						DocCount    int    `json:"doc_count"`
						ByStatus    struct {
							DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
							SumOtherDocCount        int `json:"sum_other_doc_count"`
							Buckets                 []struct {
								Key         string `json:"key"`
								DocCount    int    `json:"doc_count"`
								UniquePrIds struct {
									DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
									SumOtherDocCount        int `json:"sum_other_doc_count"`
									Buckets                 []struct {
										Key      string `json:"key"`
										DocCount int    `json:"doc_count"`
									} `json:"buckets"`
								} `json:"unique_pr_ids"`
							} `json:"buckets"`
						} `json:"by_status"`
					} `json:"buckets"`
				} `json:"date_counts"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		result := pullRequestBuckets{}
		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getPullRequests()") {
			return nil, err
		}

		v1 := responseStruct{}
		v2 := responseStruct{}
		v3 := responseStruct{}
		v4 := responseStruct{}

		v1.ID = constants.APPROVED_TITLE
		v2.ID = constants.CHANGES_REQUESTED_TITLE
		v3.ID = constants.OPEN_TITLE
		v4.ID = constants.REJECTED_TITLE

		for index, v := range result.Aggregations.DateCounts.Buckets {
			startDate := v.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			if len(v.ByStatus.Buckets) > 0 {
				approvedExists := false
				changesrequestedExists := false
				openExists := false
				rejectedExists := false
				for _, value := range v.ByStatus.Buckets {

					axisValue := struct {
						X string `json:"x"`
						Y int    `json:"y"`
					}{X: startDate, Y: len(value.UniquePrIds.Buckets)}

					if value.Key == strings.ToUpper(constants.APPROVED_TITLE) {
						v1.Data = append(v1.Data, axisValue)
						approvedExists = true
					} else if value.Key == strings.ToUpper(constants.CHANGES_REQUESTED_WITH_UNDERSCORE_TITLE) {
						v2.Data = append(v2.Data, axisValue)
						changesrequestedExists = true
					} else if value.Key == strings.ToUpper(constants.OPEN_TITLE) {
						v3.Data = append(v3.Data, axisValue)
						openExists = true
					} else if value.Key == strings.ToUpper(constants.REJECTED_TITLE) {
						v4.Data = append(v4.Data, axisValue)
						rejectedExists = true
					}
				}

				if !approvedExists {
					axisValue := struct {
						X string `json:"x"`
						Y int    `json:"y"`
					}{X: startDate, Y: 0}

					v1.Data = append(v1.Data, axisValue)
				}
				if !changesrequestedExists {
					axisValue := struct {
						X string `json:"x"`
						Y int    `json:"y"`
					}{X: startDate, Y: 0}

					v2.Data = append(v2.Data, axisValue)
				}

				if !openExists {
					axisValue := struct {
						X string `json:"x"`
						Y int    `json:"y"`
					}{X: startDate, Y: 0}
					v3.Data = append(v3.Data, axisValue)
				}

				if !rejectedExists {
					axisValue := struct {
						X string `json:"x"`
						Y int    `json:"y"`
					}{X: startDate, Y: 0}
					v4.Data = append(v4.Data, axisValue)
				}

			} else {
				axisValue := struct {
					X string `json:"x"`
					Y int    `json:"y"`
				}{X: startDate, Y: 0}

				v1.Data = append(v1.Data, axisValue)
				v2.Data = append(v2.Data, axisValue)
				v3.Data = append(v3.Data, axisValue)
				v4.Data = append(v4.Data, axisValue)
			}
		}

		outputStruct := []responseStruct{v1, v2, v3, v4}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getPullRequests() chart") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getPullRequests()") {
			return nil, err
		}

		return output, nil
	}
	return nil, nil
}

// Transforms output from the getSubReportWidget function for the Successful Build Duration widget in SDA
func transformSuccessfulBuildDuration(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type SuccessfulBuildDurationInput struct {
		Aggregations struct {
			Builds struct {
				Value map[string][]int `json:"value"`
			} `json:"builds"`
		} `json:"aggregations"`
	}

	type SuccessfulBuildDurationOutput struct {
		Data []struct {
			X string `json:"x"`
			Y []int  `json:"y"`
		} `json:"data"`
		ID  string `json:"id"`
		Min string `json:"min"`
		Max string `json:"max"`
	}

	if specKey == "successfulBuildDurationSpec" {
		response, ok := data["successfulBuildDuration"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		var input SuccessfulBuildDurationInput

		err := json.Unmarshal([]byte(response), &input)
		if log.CheckErrorf(err, "Error unmarshaling response transformSuccessfulBuildDuration()") {
			return nil, err
		}

		var output SuccessfulBuildDurationOutput
		output.ID = "Build Duration"
		if min, ok := replacements["min"]; ok {
			output.Min = fmt.Sprint(min)
		}
		if max, ok := replacements["max"]; ok {
			output.Max = fmt.Sprint(max)
		}

		for key, value := range input.Aggregations.Builds.Value {
			output.Data = append(output.Data, struct {
				X string `json:"x"`
				Y []int  `json:"y"`
			}{
				X: key,
				Y: value,
			})
		}

		if output.Data == nil {
			output.Data = []struct {
				X string `json:"x"`
				Y []int  `json:"y"`
			}{}
		}

		outputJSON, err := json.Marshal([]SuccessfulBuildDurationOutput{output})
		if log.CheckErrorf(err, "Error marshaling outputStruct in transformSuccessfulBuildDuration()") {
			return nil, err
		}
		return outputJSON, nil

	}

	return nil, nil

}

func getAutomationRunsWithScanners(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type automationRuns struct {
		Aggregations struct {
			RunStatus struct {
				Value struct {
					ChartData struct {
						Data []struct {
							Name  string `json:"name"`
							Value int    `json:"value"`
						} `json:"data"`
						Info []struct {
							DrillDown struct {
								ReportType  string `json:"reportType"`
								ReportID    string `json:"reportId"`
								ReportTitle string `json:"reportTitle"`
							} `json:"drillDown"`
							Title string `json:"title"`
							Value int    `json:"value"`
						} `json:"info"`
					} `json:"chartData"`
					Total struct {
						Value int    `json:"value"`
						Key   string `json:"key"`
					} `json:"Total"`
				} `json:"value"`
			} `json:"run_status"`
		} `json:"aggregations"`
	}

	if specKey == "totalRunsSpec" {
		response, ok := data["totalAutomationRuns"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := automationRuns{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getAutomationRunsWithScanners()") {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.RunStatus.Value.Total.Value) + `}`)
		return b, nil
	} else if specKey == "runsStatusChartSpec" {

		response, ok := data["totalAutomationRuns"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := automationRuns{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getAutomationRunsWithScanners()") {
			return nil, err
		}

		b, err := json.Marshal(result.Aggregations.RunStatus.Value.ChartData)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getAutomationRunsWithScanners()") {
			return nil, err
		}
		return b, nil

	}

	return nil, nil
}

func getVulnerabilitiesOverview(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type vulnerabiltiesOverview struct {
		Aggregations struct {
			VulnerabilityStatusCounts struct {
				Value struct {
					Reopened int `json:"Reopened"`
					Resolved int `json:"Resolved"`
					Found    int `json:"Found"`
					Open     int `json:"Open"`
				} `json:"value"`
			} `json:"vulnerabilityStatusCounts"`
		} `json:"aggregations"`
	}

	if specKey == "foundVulHeaderSpec" {
		response, ok := data["vulnerabilityStatusCounts"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesOverview{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetVulnerabilitiesOverview) {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.VulnerabilityStatusCounts.Value.Found) + `}`)
		return b, nil
	} else if specKey == "reopenedVulHeaderSpec" {
		response, ok := data["vulnerabilityStatusCounts"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesOverview{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetVulnerabilitiesOverview) {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.VulnerabilityStatusCounts.Value.Reopened) + `}`)
		return b, nil
	} else if specKey == "resolvedVulHeaderSpec" {
		response, ok := data["vulnerabilityStatusCounts"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesOverview{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getVulnerabilitiesOverview()") {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.VulnerabilityStatusCounts.Value.Resolved) + `}`)
		return b, nil
	} else if specKey == "openVulHeaderSpec" {
		response, ok := data["vulnerabilityStatusCounts"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesOverview{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetVulnerabilitiesOverview) {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.VulnerabilityStatusCounts.Value.Open) + `}`)
		return b, nil
	} else if specKey == "vulOverviewChartSpec" {
		response, ok := data["vulOverviewChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetVulnerabilitiesOverviewTypeAssertion), "aggrBy is not a string in getVulnerabilitiesOverview()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetVulnerabilitiesOverviewTypeAssertion), "startDate is not a string in getVulnerabilitiesOverview()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getVulnerabilitiesOverview"), "endDate is not a string in getVulnerabilitiesOverview()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getVulnerabilitiesOverview(): ") {
			return nil, db.ErrInternalServer
		}

		type vulnerabiltiesOverViewChart struct {
			Aggregations struct {
				VulOverviewBuckets struct {
					Buckets []struct {
						KeyAsString      string `json:"key_as_string"`
						Key              int64  `json:"key"`
						DocCount         int    `json:"doc_count"`
						VulOverviewChart struct {
							Value struct {
								Found    int `json:"Found"`
								Reopened int `json:"Reopened"`
								Resolved int `json:"Resolved"`
								Open     int `json:"Open"`
							} `json:"value"`
						} `json:"vul_overview_chart"`
					} `json:"buckets"`
				} `json:"vul_overview_buckets"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		result := vulnerabiltiesOverViewChart{}
		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetVulnerabilitiesOverview) {
			return nil, err
		}

		v1 := responseStruct{}
		v2 := responseStruct{}
		v3 := responseStruct{}
		v4 := responseStruct{}

		v1.ID = "Found"
		v2.ID = "Open"
		v3.ID = "Reopened"
		v4.ID = "Resolved"

		for index, value := range result.Aggregations.VulOverviewBuckets.Buckets {
			startDate := value.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			found := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.VulOverviewChart.Value.Found)}
			open := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.VulOverviewChart.Value.Open)}

			reopened := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.VulOverviewChart.Value.Reopened)}

			resolved := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.VulOverviewChart.Value.Resolved)}

			v1.Data = append(v1.Data, found)
			v2.Data = append(v2.Data, open)
			v3.Data = append(v3.Data, reopened)
			v4.Data = append(v4.Data, resolved)
		}
		outputStruct := []responseStruct{v1, v2, v3, v4}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getVulnerabilitiesOverview() chart") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getVulnerabilitiesOverview()") {
			return nil, err
		}
		return output, nil
	}

	return nil, nil
}

func getOpenVulnerabilitiesOverview(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type vulnerabiltiesOverview struct {
		Aggregations struct {
			SeverityCounts struct {
				Value struct {
					VeryHigh int `json:"VERY_HIGH"`
					High     int `json:"HIGH"`
					Medium   int `json:"MEDIUM"`
					Low      int `json:"LOW"`
				} `json:"value"`
			} `json:"severityCounts"`
		} `json:"aggregations"`
	}

	if specKey == "veryHighSeverityHeaderSpec" {
		response, ok := data["openVulSeverityCount"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesOverview{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetVulnerabilitiesOverview) {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.SeverityCounts.Value.VeryHigh) + `}`)
		return b, nil
	} else if specKey == "highSeverityHeaderSpec" {
		response, ok := data["openVulSeverityCount"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesOverview{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetVulnerabilitiesOverview) {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.SeverityCounts.Value.High) + `}`)
		return b, nil
	} else if specKey == "mediumSeverityHeaderSpec" {
		response, ok := data["openVulSeverityCount"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesOverview{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetVulnerabilitiesOverview) {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.SeverityCounts.Value.Medium) + `}`)
		return b, nil
	} else if specKey == "lowSeverityHeaderSpec" {
		response, ok := data["openVulSeverityCount"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesOverview{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getOpenVulnerabilitiesOverview()") {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.SeverityCounts.Value.Low) + `}`)
		return b, nil
	} else if specKey == "openVulAgeChartSpec" {
		response, ok := data["openVulAgeChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type vulnerabiltiesOverViewChart struct {
			Aggregations struct {
				AgeCounts struct {
					Value []struct {
						ID    string `json:"id"`
						Value []any  `json:"value"`
					} `json:"value"`
				} `json:"ageCounts"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		result := vulnerabiltiesOverViewChart{}
		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getOpenVulnerabilitiesOverview()") {
			return nil, err
		}

		b, err := json.Marshal(result.Aggregations.AgeCounts.Value)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getOpenVulnerabilitiesOverview()") {
			return nil, err
		}
		return b, nil
	}

	return nil, nil
}

func getWorkload(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type Data struct {
		Defect      int      `json:"DEFECT"`
		Feature     int      `json:"FEATURE"`
		Risk        int      `json:"RISK"`
		TechDebt    int      `json:"TECH_DEBT"`
		DefectSet   []string `json:"DEFECT_SET"`
		FeatureSet  []string `json:"FEATURE_SET"`
		RiskSet     []string `json:"RISK_SET"`
		TechDebtSet []string `json:"TECH_DEBT_SET"`
	}

	type workload struct {
		Aggregations struct {
			WorkLoadCounts struct {
				Value struct {
					HeaderValue int             `json:"headerValue"`
					Dates       map[string]Data `json:"dates"`
				} `json:"value"`
			} `json:"work_load_counts"`
		} `json:"aggregations"`
	}

	if specKey == "flowWorkLoadHeaderSpec" {
		response, ok := data["flowWorkLoad"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := workload{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getWorkload()") {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.WorkLoadCounts.Value.HeaderValue) + `}`)
		return b, nil
	} else if specKey == "flowWorkLoadChartSpec" {
		response, ok := data["flowWorkLoad"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetWorkloadTypeAssertion), "aggrBy is not a string in getWorkload()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetWorkloadTypeAssertion), "startDate is not a string in getWorkload()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getWorkload"), "endDate is not a string in getWorkload()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getWorkload(): ") {
			return nil, db.ErrInternalServer
		}

		result := workload{}

		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getWorkload()") {
			return nil, err
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		v1 := responseStruct{}
		v2 := responseStruct{}
		v3 := responseStruct{}
		v4 := responseStruct{}

		v1.ID = constants.DEFECT_TITLE
		v2.ID = constants.FEATURE_TITLE
		v3.ID = constants.RISK_TITLE
		v4.ID = constants.TECH_LOWER_CASE_DEBT_TITLE

		// Extract keys (dates) from the map and convert them to time.Time
		var dates []time.Time
		for key := range result.Aggregations.WorkLoadCounts.Value.Dates {
			date, err := time.Parse("2006-01-02", key)
			if err != nil {
				return nil, err
			}
			dates = append(dates, date)
		}

		// Sort the dates
		sort.Slice(dates, func(i, j int) bool {
			return dates[i].Before(dates[j])
		})

		// Iterate over sorted dates and access values
		for _, date := range dates {
			key := date.Format("2006-01-02")
			value := result.Aggregations.WorkLoadCounts.Value.Dates[key]

			startDate := key

			defect := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.Defect)}
			feature := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.Feature)}

			risk := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.Risk)}

			techDebt := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.TechDebt)}

			v1.Data = append(v1.Data, defect)
			v2.Data = append(v2.Data, feature)
			v3.Data = append(v3.Data, risk)
			v4.Data = append(v4.Data, techDebt)
		}

		outputStruct := []responseStruct{v1, v2, v3, v4}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkload() chart") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkload()") {
			return nil, err
		}
		return output, nil
	}

	return nil, nil
}

func getCycleTime(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	if specKey == "flowCycleTimeHeaderSpec" {

		type cycleTime struct {
			Aggregations struct {
				FlowCycleTimeCount struct {
					Value struct {
						Value         string `json:"value"`
						ValueInMillis int    `json:"valueInMillis"`
					} `json:"value"`
				} `json:"flow_cycle_time_count"`
			} `json:"aggregations"`
		}
		response, ok := data["flowCycleTimeHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := cycleTime{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getCycleTime()") {
			return nil, err
		}

		b := []byte(`{"value":"` + result.Aggregations.FlowCycleTimeCount.Value.Value + `"}`)
		return b, nil
	} else if specKey == "flowCycleTimeChartSpec" {

		type cycleTimeChart struct {
			Aggregations struct {
				FlowCycleTimeBuckets struct {
					Buckets []struct {
						KeyAsString        string `json:"key_as_string"`
						Key                int64  `json:"key"`
						DocCount           int    `json:"doc_count"`
						FlowCycleTimeCount struct {
							Value struct {
								TechDebt int `json:"TECH_DEBT"`
								Defect   int `json:"DEFECT"`
								Feature  int `json:"FEATURE"`
								Risk     int `json:"RISK"`
							} `json:"value"`
						} `json:"flow_cycle_time_count"`
					} `json:"buckets"`
				} `json:"flow_cycle_time_buckets"`
			} `json:"aggregations"`
		}

		response, ok := data["flowCycleTimeChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetCycleTimeTypeAssertion), "aggrBy is not a string in getCycleTime()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetCycleTimeTypeAssertion), "startDate is not a string in getCycleTime()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getCycleTime"), "endDate is not a string in getCycleTime()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getCycleTime(): ") {
			return nil, db.ErrInternalServer
		}

		result := cycleTimeChart{}

		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getCycleTime()") {
			return nil, err
		}

		type responseStruct struct {
			ID             string `json:"id"`
			YAxisFormatter struct {
				Type string `json:"type"`
			} `json:"yAxisFormatter"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		v1 := responseStruct{}
		v2 := responseStruct{}
		v3 := responseStruct{}
		v4 := responseStruct{}

		v1.ID = constants.DEFECT_TITLE
		v2.ID = constants.FEATURE_TITLE
		v3.ID = constants.RISK_TITLE
		v4.ID = constants.TECH_LOWER_CASE_DEBT_TITLE
		v1.YAxisFormatter.Type = "TIME_DURATION"
		v2.YAxisFormatter.Type = "TIME_DURATION"
		v3.YAxisFormatter.Type = "TIME_DURATION"
		v4.YAxisFormatter.Type = "TIME_DURATION"

		for index, value := range result.Aggregations.FlowCycleTimeBuckets.Buckets {
			startDate := value.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			defect := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowCycleTimeCount.Value.Defect)}
			feature := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowCycleTimeCount.Value.Feature)}

			risk := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowCycleTimeCount.Value.Risk)}

			techDebt := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowCycleTimeCount.Value.TechDebt)}

			v1.Data = append(v1.Data, defect)
			v2.Data = append(v2.Data, feature)
			v3.Data = append(v3.Data, risk)
			v4.Data = append(v4.Data, techDebt)
		}
		outputStruct := []responseStruct{v1, v2, v3, v4}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getCycleTime() chart") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getCycleTime()") {
			return nil, err
		}
		return output, nil
	}

	return nil, nil
}

func getWorkEfficiency(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	if specKey == "flowEfficiencyHeaderSpec" {
		type activeWorkTime struct {
			Aggregations struct {
				FlowEfficiencyCount struct {
					Value string `json:"value"`
				} `json:"flow_efficiency_count"`
			} `json:"aggregations"`
		}
		response, ok := data["flowEfficiencyHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := activeWorkTime{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetWorkEfficiency) {
			return nil, err
		}

		b := []byte(`{"value":"` + result.Aggregations.FlowEfficiencyCount.Value + `"}`)
		return b, nil
	} else if specKey == "flowEfficiencyChartSpec" {

		type activeWorkTimeChart struct {
			Aggregations struct {
				FlowEffTimeBuckets struct {
					Buckets []struct {
						KeyAsString         string `json:"key_as_string"`
						Key                 int64  `json:"key"`
						DocCount            int    `json:"doc_count"`
						FlowEfficiencyCount struct {
							Value struct {
								TechDebt int `json:"TECH_DEBT"`
								Defect   int `json:"DEFECT"`
								Feature  int `json:"FEATURE"`
								Risk     int `json:"RISK"`
							} `json:"value"`
						} `json:"flow_efficiency_count"`
					} `json:"buckets"`
				} `json:"flow_eff_time_buckets"`
			} `json:"aggregations"`
		}

		response, ok := data["flowEfficiencyChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetWorkEfficiencyTypeAssertion), "aggrBy is not a string in getWorkEfficiency()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetWorkEfficiencyTypeAssertion), "startDate is not a string in getWorkEfficiency()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getWorkEfficiency"), "endDate is not a string in getWorkEfficiency()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getWorkEfficiency(): ") {
			return nil, db.ErrInternalServer
		}

		result := activeWorkTimeChart{}

		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetWorkEfficiency) {
			return nil, err
		}

		type responseStruct struct {
			ID             string `json:"id"`
			YAxisFormatter struct {
				Type            string `json:"type"`
				AppendUnitValue string `json:"appendUnitValue"`
			} `json:"yAxisFormatter"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		v1 := responseStruct{}
		v2 := responseStruct{}
		v3 := responseStruct{}
		v4 := responseStruct{}

		v1.ID = constants.DEFECT_TITLE
		v2.ID = constants.FEATURE_TITLE
		v3.ID = constants.RISK_TITLE
		v4.ID = constants.TECH_LOWER_CASE_DEBT_TITLE
		v1.YAxisFormatter.Type = "APPEND_UNIT"
		v2.YAxisFormatter.Type = "APPEND_UNIT"
		v3.YAxisFormatter.Type = "APPEND_UNIT"
		v4.YAxisFormatter.Type = "APPEND_UNIT"
		v1.YAxisFormatter.AppendUnitValue = "%"
		v2.YAxisFormatter.AppendUnitValue = "%"
		v3.YAxisFormatter.AppendUnitValue = "%"
		v4.YAxisFormatter.AppendUnitValue = "%"

		for index, value := range result.Aggregations.FlowEffTimeBuckets.Buckets {
			startDate := value.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			defect := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowEfficiencyCount.Value.Defect)}
			feature := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowEfficiencyCount.Value.Feature)}

			risk := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowEfficiencyCount.Value.Risk)}

			techDebt := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowEfficiencyCount.Value.TechDebt)}

			v1.Data = append(v1.Data, defect)
			v2.Data = append(v2.Data, feature)
			v3.Data = append(v3.Data, risk)
			v4.Data = append(v4.Data, techDebt)
		}
		outputStruct := []responseStruct{v1, v2, v3, v4}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkEfficiency() chart") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkEfficiency()") {
			return nil, err
		}
		return output, nil
	} else if specKey == "flowWaitTimeHeaderSpec" {
		type workWaitTime struct {
			Aggregations struct {
				FlowWaitTimeCount struct {
					Value string `json:"value"`
				} `json:"flow_wait_time_count"`
			} `json:"aggregations"`
		}
		response, ok := data["flowWaitTimeHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := workWaitTime{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getWorkEfficiency()") {
			return nil, err
		}

		b := []byte(`{"value":"` + result.Aggregations.FlowWaitTimeCount.Value + `"}`)
		return b, nil
	} else if specKey == "flowWaitTimeChartSpec" {

		type workWaitTimeChart struct {
			Aggregations struct {
				FlowWaitTimeCount struct {
					Value []struct {
						X string `json:"x"`
						Y int    `json:"y"`
					} `json:"value"`
				} `json:"flow_wait_time_count"`
			} `json:"aggregations"`
		}

		response, ok := data["flowWaitTimeChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := workWaitTimeChart{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getWorkEfficiency()") {
			return nil, err
		}

		type responseStruct struct {
			ID             string `json:"id"`
			YAxisFormatter struct {
				Type            string `json:"type"`
				AppendUnitValue string `json:"appendUnitValue"`
			} `json:"yAxisFormatter"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		v1 := responseStruct{}
		v1.ID = "work wait time"
		v1.YAxisFormatter.Type = "APPEND_UNIT"
		v1.YAxisFormatter.AppendUnitValue = "%"

		for _, value := range result.Aggregations.FlowWaitTimeCount.Value {

			someData := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: value.X, Y: value.Y}

			v1.Data = append(v1.Data, someData)
		}
		outputStruct := []responseStruct{v1}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkEfficiency() chart") {
			return nil, err
		}

		return b, nil
	}

	return nil, nil
}

// function to transform data fetched for the velocity in flow metrics
func getVelocity(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "flowVelocityHeaderSpec" {
		response, ok := data["flowVelocityHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type velocity struct {
			Aggregations struct {
				Velocity struct {
					Value int `json:"value"`
				} `json:"velocity"`
			} `json:"aggregations"`
		}
		result := velocity{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getVelocity()") {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.Velocity.Value) + `}`)
		return b, nil
	} else if specKey == "flowVelocityChartSpec" {

		response, ok := data["flowVelocityChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetVelocityTypeAssertion), "aggrBy is not a string in getVelocity()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetVelocityTypeAssertion), "startDate is not a string in getVelocity()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getVelocity"), "endDate is not a string in getVelocity()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getVelocity(): ") {
			return nil, db.ErrInternalServer
		}

		type velocityChart struct {
			Aggregations struct {
				FlowVelocityBuckets struct {
					Buckets []struct {
						KeyAsString       string `json:"key_as_string"`
						Key               int64  `json:"key"`
						DocCount          int    `json:"doc_count"`
						FlowVelocityCount struct {
							Value struct {
								TechDebt int `json:"TECH_DEBT"`
								Defect   int `json:"DEFECT"`
								Feature  int `json:"FEATURE"`
								Risk     int `json:"RISK"`
							} `json:"value"`
						} `json:"flow_velocity_count"`
					} `json:"buckets"`
				} `json:"flow_velocity_buckets"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		result := velocityChart{}
		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getVelocity()") {
			return nil, err
		}

		v1 := responseStruct{}
		v2 := responseStruct{}
		v3 := responseStruct{}
		v4 := responseStruct{}

		v1.ID = constants.DEFECT_TITLE
		v2.ID = constants.FEATURE_TITLE
		v3.ID = constants.RISK_TITLE
		v4.ID = constants.TECH_LOWER_CASE_DEBT_TITLE

		for index, value := range result.Aggregations.FlowVelocityBuckets.Buckets {
			startDate := value.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			defect := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowVelocityCount.Value.Defect)}
			feature := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowVelocityCount.Value.Feature)}

			risk := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowVelocityCount.Value.Risk)}

			techDebt := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowVelocityCount.Value.TechDebt)}

			v1.Data = append(v1.Data, defect)
			v2.Data = append(v2.Data, feature)
			v3.Data = append(v3.Data, risk)
			v4.Data = append(v4.Data, techDebt)
		}
		outputStruct := []responseStruct{v1, v2, v3, v4}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getVelocity() chart") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getVelocity()") {
			return nil, err
		}
		return output, nil

	}

	return nil, nil
}

func getWorkItemDistribution(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "flowDistributionChartAvgSpec" {
		response, ok := data["flowDistributionAvgChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type workItemDistribution struct {
			Aggregations struct {
				FlowDistributionAvgCount struct {
					Value []struct {
						Title string `json:"title"`
						Value int    `json:"value"`
					} `json:"value"`
				} `json:"flow_distribution_avg_count"`
			} `json:"aggregations"`
		}
		result := workItemDistribution{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getWorkItemDistribution()") {
			return nil, err
		}

		type responseStruct struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		}

		v1 := responseStruct{}
		v2 := responseStruct{}
		v3 := responseStruct{}
		v4 := responseStruct{}

		for _, value := range result.Aggregations.FlowDistributionAvgCount.Value {

			if value.Title == constants.DEFECT_TITLE {
				v1.Title = constants.DEFECT_TITLE
				v1.Value = value.Value
			}
			if value.Title == constants.FEATURE_TITLE {
				v2.Title = constants.FEATURE_TITLE
				v2.Value = value.Value
			}
			if value.Title == constants.RISK_TITLE {
				v3.Title = constants.RISK_TITLE
				v3.Value = value.Value
			}
			if value.Title == constants.TECH_LOWER_CASE_DEBT_TITLE {
				v4.Title = constants.TECH_LOWER_CASE_DEBT_TITLE
				v4.Value = value.Value
			}
		}
		outputStruct := []responseStruct{v1, v2, v3, v4}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkItemDistribution() chart") {
			return nil, err
		}
		return b, nil
	} else if specKey == "flowDistributionChartSpec" {

		response, ok := data["flowDistributionChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetWorkItemDistributionTypeAssertion), "aggrBy is not a string in getWorkItemDistribution()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetWorkItemDistributionTypeAssertion), "startDate is not a string in getWorkItemDistribution()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getWorkItemDistribution"), "endDate is not a string in getWorkItemDistribution()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getWorkItemDistribution(): ") {
			return nil, db.ErrInternalServer
		}

		type flowDistributionChart struct {
			Aggregations struct {
				FlowDistributionBuckets struct {
					Buckets []struct {
						KeyAsString           string `json:"key_as_string"`
						Key                   int64  `json:"key"`
						DocCount              int    `json:"doc_count"`
						FlowDistributionCount struct {
							Value struct {
								TechDebt int `json:"TECH_DEBT"`
								Defect   int `json:"DEFECT"`
								Feature  int `json:"FEATURE"`
								Risk     int `json:"RISK"`
							} `json:"value"`
						} `json:"flow_distribution_count"`
					} `json:"buckets"`
				} `json:"flow_distribution_buckets"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID             string `json:"id"`
			YAxisFormatter struct {
				Type            string `json:"type"`
				AppendUnitValue string `json:"appendUnitValue"`
			} `json:"yAxisFormatter"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		result := flowDistributionChart{}
		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getWorkItemDistribution()") {
			return nil, err
		}

		v1 := responseStruct{}
		v2 := responseStruct{}
		v3 := responseStruct{}
		v4 := responseStruct{}

		v1.ID = constants.DEFECT_TITLE
		v2.ID = constants.FEATURE_TITLE
		v3.ID = constants.RISK_TITLE
		v4.ID = constants.TECH_LOWER_CASE_DEBT_TITLE
		v1.YAxisFormatter.Type = "APPEND_UNIT"
		v2.YAxisFormatter.Type = "APPEND_UNIT"
		v3.YAxisFormatter.Type = "APPEND_UNIT"
		v4.YAxisFormatter.Type = "APPEND_UNIT"
		v1.YAxisFormatter.AppendUnitValue = "%"
		v2.YAxisFormatter.AppendUnitValue = "%"
		v3.YAxisFormatter.AppendUnitValue = "%"
		v4.YAxisFormatter.AppendUnitValue = "%"

		for index, value := range result.Aggregations.FlowDistributionBuckets.Buckets {
			startDate := value.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			defect := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowDistributionCount.Value.Defect)}
			feature := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowDistributionCount.Value.Feature)}

			risk := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowDistributionCount.Value.Risk)}

			techDebt := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.FlowDistributionCount.Value.TechDebt)}

			v1.Data = append(v1.Data, defect)
			v2.Data = append(v2.Data, feature)
			v3.Data = append(v3.Data, risk)
			v4.Data = append(v4.Data, techDebt)
		}
		outputStruct := []responseStruct{v1, v2, v3, v4}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkItemDistribution() chart") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkItemDistribution()") {
			return nil, err
		}
		return output, nil

	}

	return nil, nil
}

func getDeploymentFrequency(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "deploymentFrequencyHeaderSpec" {
		response, ok := data["deploymentFrequencyHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type deploymentFrequency struct {
			Aggregations struct {
				DeployData struct {
					Value struct {
						Average        float64 `json:"average"`
						Deployments    int     `json:"deployments"`
						DifferenceDays int     `json:"differenceDays"`
					} `json:"value"`
				} `json:"deploy_data"`
			} `json:"aggregations"`
		}
		result := deploymentFrequency{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getDeploymentFrequency()") {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.DeployData.Value.Average) + `}`)
		return b, nil
	}
	return nil, nil
}

func getDeploymentLeadTime(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "deploymentLeadTimeHeaderSpec" {
		response, ok := data["deploymentLeadTimeHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type deploymentLeadTime struct {
			Aggregations struct {
				DeployData struct {
					Value struct {
						TotalDuration int `json:"totalDuration"`
						Average       int `json:"average"`
						Deployments   int `json:"deployments"`
					} `json:"value"`
				} `json:"deploy_data"`
			} `json:"aggregations"`
		}
		result := deploymentLeadTime{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getDeploymentLeadTime()") {
			return nil, err
		}

		b := []byte(`{"valueInMillis":` + fmt.Sprint(result.Aggregations.DeployData.Value.Average) + `}`)
		return b, nil
	}
	return nil, nil
}

func getFailureRate(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "averageFailureRateHeaderSpec" {
		response, ok := data["averageFailureRateHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type failureRate struct {
			Aggregations struct {
				DeployData struct {
					Value struct {
						Average           string `json:"average"`
						Deployments       int    `json:"deployments"`
						FailedDeployments int    `json:"failedDeployments"`
					} `json:"value"`
				} `json:"deploy_data"`
			} `json:"aggregations"`
		}
		result := failureRate{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getFailureRate()") {
			return nil, err
		}

		b := []byte(`{"value":"` + fmt.Sprint(result.Aggregations.DeployData.Value.Average) + `"}`)
		return b, nil
	}
	return nil, nil
}

func getMttr(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "mttrHeaderSpec" {
		response, ok := data["mttrHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type mttr struct {
			Aggregations struct {
				Deployments struct {
					Value float64 `json:"value"`
				} `json:"deployments"`
			} `json:"aggregations"`
		}
		result := mttr{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getMttr()") {
			return nil, err
		}

		b := []byte(`{"valueInMillis":` + fmt.Sprint(int(result.Aggregations.Deployments.Value)) + `}`)
		return b, nil
	}
	return nil, nil
}

func getDeploymentFrequencyAndLeadTime(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "frequencyAndLeadTimeTrendSpec" {
		response, ok := data["frequencyAndLeadTimeTrend"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getDeploymentFrequencyAndLeadTime"), "startDate is not a string in getDeploymentFrequencyAndLeadTime()")
			return nil, db.ErrInternalServer
		}

		type deploymentFreqAndLeadTime struct {
			Aggregations struct {
				DeployBuckets struct {
					Buckets []struct {
						KeyAsString string `json:"key_as_string"`
						Key         int64  `json:"key"`
						DocCount    int    `json:"doc_count"`
						Deployments struct {
							Value struct {
								TotalDuration int `json:"totalDuration"`
								Average       int `json:"average"`
								Deployments   int `json:"deployments"`
							} `json:"value"`
						} `json:"deployments"`
					} `json:"buckets"`
				} `json:"deploy_buckets"`
			} `json:"aggregations"`
		}
		result := deploymentFreqAndLeadTime{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getDeploymentFrequencyAndLeadTime()") {
			return nil, err
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
			YAxisFormatter struct {
				Type string `json:"type"`
			} `json:"yAxisFormatter"`
		}

		v1 := responseStruct{}
		v2 := responseStruct{}

		v1.ID = "Successful deployments"
		v1.Type = "bar"
		v2.ID = "Deployment lead time"
		v2.Type = "line"
		v2.YAxisFormatter.Type = "TIME_DURATION"

		for index, value := range result.Aggregations.DeployBuckets.Buckets {

			startDate := value.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}
			deployments := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.Deployments.Value.Deployments)}
			leadTime := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.Deployments.Value.Average)}

			v1.Data = append(v1.Data, deployments)
			v2.Data = append(v2.Data, leadTime)
		}
		outputStruct := []responseStruct{v1, v2}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getDeploymentFrequencyAndLeadTime() chart") {
			return nil, err
		}
		return b, nil
	}
	return nil, nil
}

func getFailureRateAndMttr(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "failureRateAndMttrTrendSpec" {
		response, ok := data["failureRateAndMttrTrend"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getFailureRateAndMttr"), "startDate is not a string in getFailureRateAndMttr()")
			return nil, db.ErrInternalServer
		}

		type failureRateAndMttr struct {
			Aggregations struct {
				DeployBuckets struct {
					Buckets []struct {
						KeyAsString string `json:"key_as_string"`
						Key         int64  `json:"key"`
						DocCount    int    `json:"doc_count"`
						Deployments struct {
							Value struct {
								FailureRate float64 `json:"failureRate"`
								Total       int     `json:"total"`
								Mttr        float64 `json:"mttr"`
								Failed      int     `json:"failed"`
							} `json:"value"`
						} `json:"deployments"`
					} `json:"buckets"`
				} `json:"deploy_buckets"`
			} `json:"aggregations"`
		}
		result := failureRateAndMttr{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getFailureRateAndMttr()") {
			return nil, err
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
				Z string `json:"z"`
			} `json:"data"`
			YAxisFormatter struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			} `json:"yAxisFormatter"`
		}

		v1 := responseStruct{}
		v2 := responseStruct{}

		v1.ID = "Failure rate"
		v1.Type = "bar"
		v1.YAxisFormatter.AppendUnitValue = "%"
		v1.YAxisFormatter.Type = "APPEND_UNIT"
		v2.ID = "Mean time to recovery"
		v2.Type = "line"
		v2.YAxisFormatter.Type = "TIME_DURATION"

		for index, value := range result.Aggregations.DeployBuckets.Buckets {

			startDate := value.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}
			zBytes := []byte(fmt.Sprint(value.Deployments.Value.FailureRate) + `% (` + fmt.Sprint(value.Deployments.Value.Failed) + ` of ` + fmt.Sprint(value.Deployments.Value.Total) + ` failed)`)
			// zValue := `"` + fmt.Sprint(value.Deployments.Value.FailureRate) + "% (" + fmt.Sprint(value.Deployments.Value.Failed) + " of " + fmt.Sprint(value.Deployments.Value.Total) + " failed)" + `"`
			failureRate := struct {
				X string `json:"x"`
				Y int    `json:"y"`
				Z string `json:"z"`
				// "100% (3 of 3 failed)"
			}{X: startDate, Y: int(value.Deployments.Value.FailureRate), Z: string(zBytes)}
			mttr := struct {
				X string `json:"x"`
				Y int    `json:"y"`
				Z string `json:"z"`
			}{X: startDate, Y: int(value.Deployments.Value.Mttr)}

			v1.Data = append(v1.Data, failureRate)
			v2.Data = append(v2.Data, mttr)
		}
		outputStruct := []responseStruct{v1, v2}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getFailureRateAndMttr() chart") {
			return nil, err
		}
		return b, nil
	}
	return nil, nil
}

func getVulnerabilitiesByScanType(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type vulnerabiltiesByScanType struct {
		Aggregations struct {
			ScannerTypeCount struct {
				Value struct {
					Sca       int `json:"SCA"`
					Dast      int `json:"DAST"`
					Container int `json:"Container"`
					Sast      int `json:"SAST"`
				} `json:"value"`
			} `json:"scanner_type_count"`
		} `json:"aggregations"`
	}

	if specKey == "SASTHeaderSpec" {
		response, ok := data["vulnerabilityByScannerTypeHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesByScanType{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetVulnerablitiesByScanType) {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.ScannerTypeCount.Value.Sast) + `}`)
		return b, nil
	} else if specKey == "DASTHeaderSpec" {
		response, ok := data["vulnerabilityByScannerTypeHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesByScanType{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetVulnerablitiesByScanType) {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.ScannerTypeCount.Value.Dast) + `}`)
		return b, nil
	} else if specKey == "ContainerHeaderSpec" {
		response, ok := data["vulnerabilityByScannerTypeHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesByScanType{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetVulnerablitiesByScanType) {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.ScannerTypeCount.Value.Container) + `}`)
		return b, nil
	} else if specKey == "SCAHeaderSpec" {
		response, ok := data["vulnerabilityByScannerTypeHeader"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := vulnerabiltiesByScanType{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getVulnerabilitiesByScanType()") {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.ScannerTypeCount.Value.Sca) + `}`)
		return b, nil
	} else if specKey == "vulnerabilitybyscannertypechartSpec" {
		response, ok := data["vulnerabilitybyscannertypechart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type vulnerabilitiesByScanTypeChart struct {
			Aggregations struct {
				VulByScannerTypeCounts struct {
					Value struct {
						VeryHigh []struct {
							X string `json:"x"`
							Y int    `json:"y"`
						} `json:"VERY_HIGH"`
						High []struct {
							X string `json:"x"`
							Y int    `json:"y"`
						} `json:"HIGH"`
						Medium []struct {
							X string `json:"x"`
							Y int    `json:"y"`
						} `json:"MEDIUM"`
						Low []struct {
							X string `json:"x"`
							Y int    `json:"y"`
						} `json:"LOW"`
					} `json:"value"`
				} `json:"vulByScannerTypeCounts"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		result := vulnerabilitiesByScanTypeChart{}
		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getVulnerabilitiesByScanType()") {
			return nil, err
		}

		v1 := responseStruct{}
		v2 := responseStruct{}
		v3 := responseStruct{}
		v4 := responseStruct{}

		v1.Data = result.Aggregations.VulByScannerTypeCounts.Value.VeryHigh
		v2.Data = result.Aggregations.VulByScannerTypeCounts.Value.High
		v3.Data = result.Aggregations.VulByScannerTypeCounts.Value.Medium
		v4.Data = result.Aggregations.VulByScannerTypeCounts.Value.Low

		v1.ID = constants.VERY_HIGH_TITLE
		v2.ID = constants.HIGH_TITLE
		v3.ID = constants.MEDIUM_TITLE
		v4.ID = constants.LOW_TITLE

		outputStruct := []responseStruct{v1, v2, v3, v4}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getVulnerabilitiesByScanType() chart") {
			return nil, err
		}
		return b, nil
	}

	return nil, nil
}

func getSlaStatusOverview(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	response, ok := data["slaStatusOverview"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	type slaStatusOverview struct {
		Aggregations struct {
			SLAStatusOverview struct {
				Value struct {
					OpenSLACounts []struct {
						X string `json:"x"`
						Y int    `json:"y"`
					} `json:"openSlaCounts"`
					ClosedSLACounts []struct {
						X string `json:"x"`
						Y int    `json:"y"`
					} `json:"closedSlaCounts"`
					CloseSLAKey string `json:"closeSlaKey"`
					OpenSLAKey  string `json:"openSlaKey"`
				} `json:"value"`
			} `json:"slaStatusOverview"`
		} `json:"aggregations"`
	}
	result := slaStatusOverview{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getSlaStatusOverview()") {
		return nil, err
	}

	if specKey == "slaStatusOverviewOpenSpec" {

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		v1 := responseStruct{}

		v1.ID = "Open vulnerabilities"
		v1.Data = result.Aggregations.SLAStatusOverview.Value.OpenSLACounts

		outputStruct := []responseStruct{v1}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getSlaStatusOverview() chart") {
			return nil, err
		}
		return b, nil
	} else if specKey == "slaStatusOverviewClosedSpec" {

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		v1 := responseStruct{}

		v1.ID = "Resolved vulnerabilites"
		v1.Data = result.Aggregations.SLAStatusOverview.Value.ClosedSLACounts

		outputStruct := []responseStruct{v1}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getSlaStatusOverview() chart") {
			return nil, err
		}
		return b, nil
	}
	return nil, nil
}

func GetScanTypesInAutomation(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	scanAutomationResp, ok := data["scanAutomationResp"]
	if !ok {
		return nil, db.ErrInternalServer
	}
	automationRunResp, ok := data["automationRunResp"]
	if !ok {
		return nil, db.ErrInternalServer
	}
	scannerTypeResp, ok := data["scannerTypeResp"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	var sastCount, dastCount, containerCount, scaCount int

	automationResult := make(map[string]interface{})
	json.Unmarshal([]byte(automationRunResp), &automationResult)

	securityResult := make(map[string]interface{})
	json.Unmarshal([]byte(scannerTypeResp), &securityResult)

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

	if automationResult[constants.AGGREGATION] != nil {
		aggsResult := automationResult[constants.AGGREGATION].(map[string]interface{})
		if aggsResult[constants.AUTOMATION_RUN_ACTIVITY] != nil {
			automationRuns := aggsResult[constants.AUTOMATION_RUN_ACTIVITY].(map[string]interface{})
			if automationRuns[constants.VALUE] != nil {
				values := automationRuns[constants.VALUE].(map[string]interface{})
				for _, value := range values {
					runs := value.([]interface{})
					for _, runValueMap := range runs {
						runValue := runValueMap.(map[string]interface{})
						curRunID := runValue["run_id"].(string)
						var scannerTypeList []interface{}
						if val, ok := scannerRuns[curRunID]; ok {
							valueMap := val.(map[string]interface{})
							scannerTypeList = (valueMap["scanner_types"]).([]interface{})
							for _, scannerType := range scannerTypeList {
								if scannerType == "SAST" {
									sastCount++
								}
								if scannerType == "DAST" {
									dastCount++
								}
								if scannerType == "Container" {
									containerCount++
								}
								if scannerType == "SCA" {
									scaCount++
								}
							}
						}
					}
				}
			}
		}
	}

	type scanTypesInAutomation struct {
		Aggregations struct {
			ScanTypesInAutomation struct {
				Value struct {
					AutomationResult []struct {
						X string `json:"x"`
						Y int    `json:"y"`
					} `json:"automationResult"`
					RunKey        string `json:"runKey"`
					AutomationKey string `json:"automationKey"`
					RunResult     []struct {
						X string `json:"x"`
						Y int    `json:"y"`
					} `json:"runResult"`
				} `json:"value"`
			} `json:"scanTypesInAutomation"`
		} `json:"aggregations"`
	}
	result := scanTypesInAutomation{}

	err := json.Unmarshal([]byte(scanAutomationResp), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in GetScanTypesInAutomation()") {
		return nil, err
	}

	if specKey == "scanTypesInAutomationsSpec" {

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		v1 := responseStruct{}
		v2 := responseStruct{}

		v1.ID = "Workflows"
		v1.Data = result.Aggregations.ScanTypesInAutomation.Value.AutomationResult

		runResult := []struct {
			X string `json:"x"`
			Y int    `json:"y"`
		}{}

		if scaCount > 0 {
			runResult = append(runResult, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: "SCA", Y: scaCount})
		}

		if dastCount > 0 {
			runResult = append(runResult, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: "DAST", Y: dastCount})
		}

		if containerCount > 0 {
			runResult = append(runResult, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: "Container", Y: containerCount})
		}

		if sastCount > 0 {
			runResult = append(runResult, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: "SAST", Y: sastCount})
		}
		v2.ID = "Workflow Runs"
		v2.Data = runResult

		outputStruct := []responseStruct{v1, v2}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in GetScanTypesInAutomation() chart") {
			return nil, err
		}
		return b, nil
	}
	return nil, nil
}

func getMttrForVulnerabilities(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type mttr struct {
		Aggregations struct {
			AvgTTR struct {
				Value struct {
					VeryHigh              string `json:"VERY_HIGH"`
					High                  string `json:"HIGH"`
					Medium                string `json:"MEDIUM"`
					Low                   string `json:"LOW"`
					HighResolvedCount     int    `json:"HIGH_RESOLVED_COUNT"`
					LowResolvedCount      int    `json:"LOW_RESOLVED_COUNT"`
					MediumResolvedCount   int    `json:"MEDIUM_RESOLVED_COUNT"`
					VeryHighResolvedCount int    `json:"VERY_HIGH_RESOLVED_COUNT"`
				} `json:"value"`
			} `json:"Avg_TTR"`
		} `json:"aggregations"`
	}

	type mttrHeader struct {
		TitleCount int    `json:"titleCount"`
		Value      string `json:"value"`
	}

	if specKey == "MTTRHeaderSpecVeryHigh" {

		response, ok := data["MTTRHeaders"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := mttr{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getMttrForVulnerabilities()") {
			return nil, err
		}

		v1 := mttrHeader{}

		v1.TitleCount = result.Aggregations.AvgTTR.Value.VeryHighResolvedCount
		v1.Value = result.Aggregations.AvgTTR.Value.VeryHigh

		b, err := json.Marshal(v1)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getMttrForVulnerabilities() chart") {
			return nil, err
		}
		return b, nil
	} else if specKey == "MTTRHeaderSpecHigh" {

		response, ok := data["MTTRHeaders"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := mttr{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetMttrForVulnerabilities) {
			return nil, err
		}
		v1 := mttrHeader{}

		v1.TitleCount = result.Aggregations.AvgTTR.Value.HighResolvedCount
		v1.Value = result.Aggregations.AvgTTR.Value.High

		b, err := json.Marshal(v1)
		if log.CheckErrorf(err, exceptions.ErrMarshallingGetMttrForVulnerabilities) {
			return nil, err
		}
		return b, nil
	} else if specKey == "MTTRHeaderSpecMedium" {

		response, ok := data["MTTRHeaders"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := mttr{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetMttrForVulnerabilities) {
			return nil, err
		}
		v1 := mttrHeader{}

		v1.TitleCount = result.Aggregations.AvgTTR.Value.MediumResolvedCount
		v1.Value = result.Aggregations.AvgTTR.Value.Medium

		b, err := json.Marshal(v1)
		if log.CheckErrorf(err, exceptions.ErrMarshallingGetMttrForVulnerabilities) {
			return nil, err
		}
		return b, nil
	} else if specKey == "MTTRHeaderSpecLow" {

		response, ok := data["MTTRHeaders"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := mttr{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetMttrForVulnerabilities) {
			return nil, err
		}
		v1 := mttrHeader{}

		v1.TitleCount = result.Aggregations.AvgTTR.Value.LowResolvedCount
		v1.Value = result.Aggregations.AvgTTR.Value.Low

		b, err := json.Marshal(v1)
		if log.CheckErrorf(err, exceptions.ErrMarshallingGetMttrForVulnerabilities) {
			return nil, err
		}
		return b, nil
	} else if specKey == "MTTRChartSpec" {
		response, ok := data["MTTRChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetMttrForVulnerabilitiesTypeAssertion), "aggrBy is not a string in getMttrForVulnerabilities()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetMttrForVulnerabilitiesTypeAssertion), "startDate is not a string in getMttrForVulnerabilities()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getMttrForVulnerabilities"), "endDate is not a string in getMttrForVulnerabilities()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getMttrForVulnerabilities(): ") {
			return nil, db.ErrInternalServer
		}

		type mttrChart struct {
			Aggregations struct {
				MTTRChartDateBuckets struct {
					Buckets []struct {
						KeyAsString string `json:"key_as_string"`
						Key         int64  `json:"key"`
						DocCount    int    `json:"doc_count"`
						AvgTTR      struct {
							Value struct {
								VeryHigh int `json:"VERY_HIGH"`
								High     int `json:"HIGH"`
								Medium   int `json:"MEDIUM"`
								Low      int `json:"LOW"`
							} `json:"value"`
						} `json:"Avg_TTR"`
					} `json:"buckets"`
				} `json:"MTTR_chart_date_buckets"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID             string `json:"id"`
			YAxisFormatter struct {
				Type string `json:"type"`
			} `json:"yAxisFormatter"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
		}

		result := mttrChart{}
		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, exceptions.ErrUnmarshallingGetMttrForVulnerabilities) {
			return nil, err
		}

		v1 := responseStruct{}
		v2 := responseStruct{}
		v3 := responseStruct{}
		v4 := responseStruct{}

		v1.ID = constants.VERY_HIGH_TITLE
		v1.YAxisFormatter.Type = "TIME_DURATION"
		v2.ID = constants.HIGH_TITLE
		v2.YAxisFormatter.Type = "TIME_DURATION"
		v3.ID = constants.MEDIUM_TITLE
		v3.YAxisFormatter.Type = "TIME_DURATION"
		v4.ID = constants.LOW_TITLE
		v4.YAxisFormatter.Type = "TIME_DURATION"

		for index, value := range result.Aggregations.MTTRChartDateBuckets.Buckets {
			startDate := value.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			veryHigh := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.AvgTTR.Value.VeryHigh)}
			high := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.AvgTTR.Value.High)}

			medium := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.AvgTTR.Value.Medium)}

			low := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.AvgTTR.Value.Low)}

			v1.Data = append(v1.Data, veryHigh)
			v2.Data = append(v2.Data, high)
			v3.Data = append(v3.Data, medium)
			v4.Data = append(v4.Data, low)
		}
		outputStruct := []responseStruct{v1, v2, v3, v4}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, exceptions.ErrMarshallingGetMttrForVulnerabilities) {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getMttrForVulnerabilities()") {
			return nil, err
		}
		return output, nil
	}

	return nil, nil
}

func getCodeChurn(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "codeChurnChartSpec" {

		response, ok := data["codeChurnChart"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		aggrBy, ok := replacements["aggrBy"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetCodeChurnTypeAssertion), "aggrBy is not a string in getCodeChurn()")
			return nil, db.ErrInternalServer
		}

		replacementStartDate, ok := replacements["dateHistogramMin"].(string)
		if !ok {
			log.Errorf(exceptions.GetExceptionByCode(exceptions.ErrGetCodeChurnTypeAssertion), "startDate is not a string in getCodeChurn()")
			return nil, db.ErrInternalServer
		}

		//Replace replacementStartDate with normalizedStartDate if specified. This field is not mandatory(no errors)
		normalizedStartDate, ok := replacements["normalizeMonthInSpec"].(string)
		if ok && normalizedStartDate != "@x" {
			replacementStartDate = normalizedStartDate[0:10]
		}

		replacemenEndDate, ok := replacements["dateHistogramMax"].(string)
		if !ok {
			log.Errorf(errors.New("error in type assertion in getCodeChurn"), "endDate is not a string in getCodeChurn()")
			return nil, db.ErrInternalServer
		}

		bucketDates, err := helper.GetBucketDates(aggrBy, replacementStartDate, replacemenEndDate)
		if log.CheckErrorf(err, "error fetching output from helper.GetBucketDates() in internal.getCodeChurn(): ") {
			return nil, db.ErrInternalServer
		}

		type codeChurnChart struct {
			Aggregations struct {
				CodeChurnBuckets struct {
					Buckets []struct {
						KeyAsString string `json:"key_as_string"`
						Key         int64  `json:"key"`
						DocCount    int    `json:"doc_count"`
						CodeChurn   struct {
							Value struct {
								LinesDeleted int `json:"lines_deleted"`
								LinesAdded   int `json:"lines_added"`
							} `json:"value"`
						} `json:"code_churn"`
					} `json:"buckets"`
				} `json:"code_churn_buckets"`
			} `json:"aggregations"`
		}

		type responseStruct struct {
			ID   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
			IsClickDisabled bool `json:"isClickDisable"`
		}

		result := codeChurnChart{}
		err = json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getCodeChurn()") {
			return nil, err
		}

		v1 := responseStruct{}
		v2 := responseStruct{}

		v1.ID = "Additions"
		v1.IsClickDisabled = true

		v2.ID = "Deletions"
		v2.IsClickDisabled = true

		for index, value := range result.Aggregations.CodeChurnBuckets.Buckets {
			startDate := value.KeyAsString

			if index == 0 {
				startDate = replacementStartDate
			}

			additions := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.CodeChurn.Value.LinesAdded)}
			deletions := struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{X: startDate, Y: int(value.CodeChurn.Value.LinesDeleted)}

			v1.Data = append(v1.Data, additions)
			v2.Data = append(v2.Data, deletions)

		}
		outputStruct := []responseStruct{v1, v2}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getCodeChurn() chart") {
			return nil, err
		}

		//Add Start and End dates for hovering
		output, err := helper.AddDateBuckets(b, bucketDates)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getCodeChurn()") {
			return nil, err
		}
		return output, nil

	}

	return nil, nil
}

func getAverageDeploymentTime(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type averageDeploymentTime struct {
		Aggregations struct {
			DeployData struct {
				Value []struct {
					Average int    `json:"average"`
					Count   int    `json:"count"`
					From    string `json:"from"`
					To      string `json:"to"`
					Value   int    `json:"value"`
				} `json:"value"`
			} `json:"deploy_data"`
		} `json:"aggregations"`
	}

	type avgDeploymentTimeResponse struct {
		FromTitle     string `json:"fromTitle"`
		ToTitle       string `json:"toTitle"`
		TotalCount    int    `json:"totalCount"`
		TotalDuration int    `json:"totalDuration"`
		Value         int    `json:"value"`
	}

	if specKey == "averageDeploymentTimeSpec" {

		response, ok := data["averageDeploymentTime"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := averageDeploymentTime{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getAverageDeploymentTime()") {
			return nil, err
		}

		outputStruct := []avgDeploymentTimeResponse{}
		for _, value := range result.Aggregations.DeployData.Value {
			v1 := avgDeploymentTimeResponse{}
			v1.FromTitle = value.From
			v1.ToTitle = value.To
			v1.Value = value.Average
			v1.TotalCount = value.Count
			v1.TotalDuration = value.Value
			outputStruct = append(outputStruct, v1)
		}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getAverageDeploymentTime() chart") {
			return nil, err
		}
		return b, nil
	}

	return nil, nil
}

func getCwetmTop25Vulnerabilities(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	type cwetmTop25VulnerabilitiesOverview struct {
		Aggregations struct {
			Top25CWE struct {
				Value struct {
					Top25Table struct {
						CWE787 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-787"`
						CWE79 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-79"`
						CWE89 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-89"`
						CWE416 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-416"`
						CWE78 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-78"`
						CWE20 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-20"`
						CWE125 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-125"`
						CWE22 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-22"`
						CWE352 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-352"`
						CWE434 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-434"`
						CWE862 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-862"`
						CWE476 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-476"`
						CWE287 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-287"`
						CWE190 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-190"`
						CWE502 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-502"`
						CWE77 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-77"`
						CWE119 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-119"`
						CWE798 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-798"`
						CWE918 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-918"`
						CWE306 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-306"`
						CWE362 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-362"`
						CWE269 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-269"`
						CWE94 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-94"`
						CWE863 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-863"`
						CWE276 struct {
							Name        string `json:"name"`
							IssuesFound int    `json:"issuesFound"`
						} `json:"CWE-276"`
					} `json:"top25Table"`
					Top25TotalCount int `json:"top25TotalCount"`
				} `json:"value"`
			} `json:"top25CWE"`
		} `json:"aggregations"`
	}

	response, ok := data["top25Vul"]
	if !ok {
		return nil, db.ErrInternalServer
	}
	result := cwetmTop25VulnerabilitiesOverview{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response getCwetmTop25Vulnerabilities()") {
		return nil, err
	}

	if result.Aggregations.Top25CWE.Value.Top25TotalCount == 0 {
		return nil, db.ErrNoDataFound
	}

	if specKey == "top25VulHeaderSpec" {

		type responsewithSubTitle struct {
			SubTitle struct {
				Title string `json:"title"`
			} `json:"subTitle"`
			Value    int `json:"value"`
			SubValue int `json:"subValue"`
		}

		responseArray := responsewithSubTitle{
			SubTitle: struct {
				Title string `json:"title"`
			}{
				Title: "CWE<sup>TM</sup> top 25 vulnerabilities",
			},
			Value:    result.Aggregations.Top25CWE.Value.Top25TotalCount,
			SubValue: 25,
		}
		responseJSON, err := json.Marshal(responseArray)

		if log.CheckErrorf(err, "error marshaling response: %v") {
			return nil, err
		}

		return responseJSON, nil
	} else if specKey == "top25VulChartSpec" {

		type responsewithTable struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			IssuesFound int    `json:"issuesFound"`
		}

		var indivMap []responsewithTable
		top25TableMap := map[string]struct {
			Name        string `json:"name"`
			IssuesFound int    `json:"issuesFound"`
		}{
			"CWE-787": result.Aggregations.Top25CWE.Value.Top25Table.CWE787,
			"CWE-79":  result.Aggregations.Top25CWE.Value.Top25Table.CWE79,
			"CWE-89":  result.Aggregations.Top25CWE.Value.Top25Table.CWE89,
			"CWE-416": result.Aggregations.Top25CWE.Value.Top25Table.CWE416,
			"CWE-78":  result.Aggregations.Top25CWE.Value.Top25Table.CWE78,
			"CWE-20":  result.Aggregations.Top25CWE.Value.Top25Table.CWE20,
			"CWE-125": result.Aggregations.Top25CWE.Value.Top25Table.CWE125,
			"CWE-22":  result.Aggregations.Top25CWE.Value.Top25Table.CWE22,
			"CWE-352": result.Aggregations.Top25CWE.Value.Top25Table.CWE352,
			"CWE-434": result.Aggregations.Top25CWE.Value.Top25Table.CWE434,
			"CWE-862": result.Aggregations.Top25CWE.Value.Top25Table.CWE862,
			"CWE-476": result.Aggregations.Top25CWE.Value.Top25Table.CWE476,
			"CWE-287": result.Aggregations.Top25CWE.Value.Top25Table.CWE287,
			"CWE-190": result.Aggregations.Top25CWE.Value.Top25Table.CWE190,
			"CWE-502": result.Aggregations.Top25CWE.Value.Top25Table.CWE502,
			"CWE-77":  result.Aggregations.Top25CWE.Value.Top25Table.CWE77,
			"CWE-119": result.Aggregations.Top25CWE.Value.Top25Table.CWE119,
			"CWE-798": result.Aggregations.Top25CWE.Value.Top25Table.CWE798,
			"CWE-918": result.Aggregations.Top25CWE.Value.Top25Table.CWE918,
			"CWE-306": result.Aggregations.Top25CWE.Value.Top25Table.CWE306,
			"CWE-362": result.Aggregations.Top25CWE.Value.Top25Table.CWE362,
			"CWE-269": result.Aggregations.Top25CWE.Value.Top25Table.CWE269,
			"CWE-94":  result.Aggregations.Top25CWE.Value.Top25Table.CWE94,
			"CWE-863": result.Aggregations.Top25CWE.Value.Top25Table.CWE863,
			"CWE-276": result.Aggregations.Top25CWE.Value.Top25Table.CWE276,
		}

		for _, id := range constants.CWETop25IDsList {
			data := top25TableMap[id]
			if data.IssuesFound > 0 {
				entry := responsewithTable{
					ID:          id,
					Name:        data.Name,
					IssuesFound: data.IssuesFound,
				}
				indivMap = append(indivMap, entry)
			}
		}
		responseJSON, err := json.Marshal(indivMap)
		if log.CheckErrorf(err, "Error marshaling response in getCwetmTop25Vulnerabilities:") {
			return nil, err
		}
		return responseJSON, nil
	}
	return nil, nil
}

func getDeploymentFrequencyComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {
	response, ok := data["deploymentFrequencyHeader"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.DeploymentFrequencyComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getDeploymentFrequencyComponentComparison()") {
		return nil, err
	}
	compareReports := createDeploymentFrequencyCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getDeploymentFrequencyComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createDeploymentFrequencyCompareReports(orgData *constants.Organization, compComparison *constants.DeploymentFrequencyComponentComparison) *constants.DeploymentFrequencyCompareReports {
	var compareReports constants.DeploymentFrequencyCompareReports
	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalDeployments float64
	var totalDifferenceDays float64

	// Count components
	for _, bucket := range compComparison.Aggregations.DeploymentFrequencyComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				totalDeployments += float64(bucket.DeployData.Value.Deployments)

				totalDifferenceDays = bucket.DeployData.Value.DifferenceDays

				data := []struct {
					Title string  `json:"title"`
					Value float64 `json:"value"`
				}{
					{
						Title: "average per day",
						Value: calculateAverageDeployments(float64(bucket.DeployData.Value.Deployments), bucket.DeployData.Value.DifferenceDays),
					},
					{
						Title: "deployments",
						Value: float64(bucket.DeployData.Value.Deployments),
					},
					{
						Title: "differenceDays",
						Value: bucket.DeployData.Value.DifferenceDays,
					},
				}

				deploymentData := bucket.DeployData.Value.Deployments
				differenceDaysData := bucket.DeployData.Value.DifferenceDays

				compareReports.ComponentCount++
				var compCompareReports constants.DeploymentFrequencyCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data

				// Roundoff the value so that instead of 0 it will show 1 on UI
				calculateAverageDeploymentsVal := calculateAverageDeployments(float64(deploymentData), differenceDaysData)
				roundOffcalculateAverageDeployments := math.Round(calculateAverageDeploymentsVal)
				compCompareReports.TotalValue = int(roundOffcalculateAverageDeployments)

				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string  `json:"title"`
		Value float64 `json:"value"`
	}{
		{
			Title: "average per day",
			Value: calculateAverageDeployments(totalDeployments, totalDifferenceDays),
		},
		{
			Title: "deployments",
			Value: float64(totalDeployments),
		},
		{
			Title: "differenceDays",
			Value: totalDifferenceDays,
		},
	}

	compareReports.Section.Data = data
	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createDeploymentFrequencyCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalDeployments += subCompareReports.Section.Data[1].Value
		if totalDifferenceDays == 0 {
			totalDifferenceDays += subCompareReports.Section.Data[2].Value
		}

	}
	compareReports.Section.Data[0].Value = calculateAverageDeployments(totalDeployments, totalDifferenceDays)
	compareReports.Section.Data[1].Value = float64(totalDeployments)
	if compareReports.Section.Data[2].Value == 0 {
		compareReports.Section.Data[2].Value += totalDifferenceDays
	}

	compareReports.TotalValue = int(calculateAverageDeployments(totalDeployments, totalDifferenceDays))
	return &compareReports

}

func calculateAverageDeployments(deployments float64, differenceDays float64) float64 {
	if differenceDays > 0 {
		result := float64(deployments) / differenceDays
		roundedResult := math.Round(result*100) / 100
		return roundedResult
	}
	return 0
}

func getOpenVulnerabilitiesOverviewComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {
	response, ok := data["openVulAgeChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.OpenVulnerabilitiesOverviewComponentComparison{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getOpenVulnerabilitiesOverviewComponentComparison()") {
		return nil, err
	}
	compareReports := createOpenVulnerabilitiesOverviewCompareReports(organisation, &result)
	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getOpenVulnerabilitiesOverviewComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createOpenVulnerabilitiesOverviewCompareReports(orgData *constants.Organization, compComparison *constants.OpenVulnerabilitiesOverviewComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports
	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true
	var totalVeryHigh, totalHigh, totalMedium, totalLow int
	// Count components
	for _, bucket := range compComparison.Aggregations.OpenVulnerabilitiesOverviewComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				totalVeryHigh += bucket.OpenVulSeverityCount.Value.VeryHigh
				totalHigh += bucket.OpenVulSeverityCount.Value.High
				totalMedium += bucket.OpenVulSeverityCount.Value.Medium
				totalLow += bucket.OpenVulSeverityCount.Value.Low

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE),
						Value: bucket.OpenVulSeverityCount.Value.VeryHigh,
					}, {
						Title: strings.ToUpper(constants.HIGH_TITLE),
						Value: bucket.OpenVulSeverityCount.Value.High,
					}, {
						Title: strings.ToUpper(constants.MEDIUM_TITLE),
						Value: bucket.OpenVulSeverityCount.Value.Medium,
					},
					{
						Title: strings.ToUpper(constants.LOW_TITLE),
						Value: bucket.OpenVulSeverityCount.Value.Low,
					},
				}

				value := bucket.OpenVulSeverityCount.Value.VeryHigh + bucket.OpenVulSeverityCount.Value.High + bucket.OpenVulSeverityCount.Value.Medium + bucket.OpenVulSeverityCount.Value.Low

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break

			}
		}
	}
	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE),
			Value: totalVeryHigh,
		}, {
			Title: strings.ToUpper(constants.HIGH_TITLE),
			Value: totalHigh,
		}, {
			Title: strings.ToUpper(constants.MEDIUM_TITLE),
			Value: totalMedium,
		},
		{
			Title: strings.ToUpper(constants.LOW_TITLE),
			Value: totalLow,
		},
	}
	compareReports.Section.Data = data
	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createOpenVulnerabilitiesOverviewCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalVeryHigh += subCompareReports.Section.Data[0].Value
		totalHigh += subCompareReports.Section.Data[1].Value
		totalMedium += subCompareReports.Section.Data[2].Value
		totalLow += subCompareReports.Section.Data[3].Value
	}
	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalVeryHigh
	compareReports.Section.Data[1].Value = totalHigh
	compareReports.Section.Data[2].Value = totalMedium
	compareReports.Section.Data[3].Value = totalLow

	compareReports.TotalValue = totalVeryHigh + totalHigh + totalMedium + totalLow
	return &compareReports
}

func getVulnerabilitiesOverviewComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {
	response, ok := data["vulOverviewChart"]

	if !ok {
		return nil, db.ErrInternalServer
	}
	result := constants.VulnerabilitiesOverviewComponentComparison{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getVulnerabilitiesOverviewComponentComparison()") {
		return nil, err
	}
	compareReports := createVulnerabilitiesOverviewCompareReports(organisation, &result)
	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getVulnerabilitiesOverviewComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createVulnerabilitiesOverviewCompareReports(orgData *constants.Organization, compComparison *constants.VulnerabilitiesOverviewComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports
	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true
	var totalReopened, totalResolved, totalFound, totalOpen int
	// Count components
	for _, bucket := range compComparison.Aggregations.VulnerabilitiesOverviewComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				totalFound += bucket.VulnerabilityStatusCounts.Value.Found
				totalOpen += bucket.VulnerabilityStatusCounts.Value.Open
				totalReopened += bucket.VulnerabilityStatusCounts.Value.Reopened
				totalResolved += bucket.VulnerabilityStatusCounts.Value.Resolved

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: "Found",
						Value: bucket.VulnerabilityStatusCounts.Value.Found,
					},
					{
						Title: "Open",
						Value: bucket.VulnerabilityStatusCounts.Value.Open,
					},
					{
						Title: "Reopened",
						Value: bucket.VulnerabilityStatusCounts.Value.Reopened,
					}, {
						Title: "Resolved",
						Value: bucket.VulnerabilityStatusCounts.Value.Resolved,
					},
				}

				value := bucket.VulnerabilityStatusCounts.Value.Reopened + bucket.VulnerabilityStatusCounts.Value.Resolved + bucket.VulnerabilityStatusCounts.Value.Found + bucket.VulnerabilityStatusCounts.Value.Open

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break

			}
		}
	}
	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: "Found",
			Value: totalFound,
		},
		{
			Title: "Open",
			Value: totalOpen,
		},
		{
			Title: "Reopened",
			Value: totalReopened,
		}, {
			Title: "Resolved",
			Value: totalResolved,
		},
	}
	compareReports.Section.Data = data
	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createVulnerabilitiesOverviewCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalFound += subCompareReports.Section.Data[0].Value
		totalOpen += subCompareReports.Section.Data[1].Value
		totalReopened += subCompareReports.Section.Data[2].Value
		totalResolved += subCompareReports.Section.Data[3].Value

	}
	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalFound
	compareReports.Section.Data[1].Value = totalOpen
	compareReports.Section.Data[2].Value = totalReopened
	compareReports.Section.Data[3].Value = totalResolved

	compareReports.TotalValue = totalReopened + totalResolved + totalFound + totalOpen
	return &compareReports
}

func getVelocityComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {
	response, ok := data["flowVelocityChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.VelocityComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getVelocityComponentComparison()") {
		return nil, err
	}

	compareReports := createVelocityCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getVelocityComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createVelocityCompareReports(orgData *constants.Organization, compComparison *constants.VelocityComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalBugs, totalFeature, totalRisk, totalTechDebt int

	// Count components
	for _, bucket := range compComparison.Aggregations.FlowVelocityComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				totalBugs += bucket.FlowVelocityCount.Value.Defect
				totalFeature += bucket.FlowVelocityCount.Value.Feature
				totalRisk += bucket.FlowVelocityCount.Value.Risk
				totalTechDebt += bucket.FlowVelocityCount.Value.TechDebt

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.BUGS_TITLE,
						Value: bucket.FlowVelocityCount.Value.Defect,
					}, {
						Title: constants.FEATURE_TITLE,
						Value: bucket.FlowVelocityCount.Value.Feature,
					}, {
						Title: constants.RISK_TITLE,
						Value: bucket.FlowVelocityCount.Value.Risk,
					},
					{
						Title: constants.TECH_DEBT_TITLE,
						Value: bucket.FlowVelocityCount.Value.TechDebt,
					},
				}

				value := bucket.FlowVelocityCount.Value.TechDebt + bucket.FlowVelocityCount.Value.Defect + bucket.FlowVelocityCount.Value.Feature + bucket.FlowVelocityCount.Value.Risk

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.BUGS_TITLE,
			Value: totalBugs,
		}, {
			Title: constants.FEATURE_TITLE,
			Value: totalFeature,
		}, {
			Title: constants.RISK_TITLE,
			Value: totalRisk,
		},
		{
			Title: constants.TECH_DEBT_TITLE,
			Value: totalTechDebt,
		},
	}
	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createVelocityCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalBugs += subCompareReports.Section.Data[0].Value
		totalFeature += subCompareReports.Section.Data[1].Value
		totalRisk += subCompareReports.Section.Data[2].Value
		totalTechDebt += subCompareReports.Section.Data[3].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalBugs
	compareReports.Section.Data[1].Value = totalFeature
	compareReports.Section.Data[2].Value = totalRisk
	compareReports.Section.Data[3].Value = totalTechDebt

	compareReports.TotalValue = totalBugs + totalFeature + totalRisk + totalTechDebt

	return &compareReports
}

func sumCounts(compareReports *constants.CompareReports, featureSum, riskSum, techDebtSum, bugSum int64) (int64, int64, int64, int64) {
	// Sum counts of feature, risk, tech debt, and bugs for the current level
	if !compareReports.IsSubOrg {
		for _, item := range compareReports.Section.Data {
			switch item.Title {
			case constants.BUGS_TITLE:
				bugSum += int64(item.Value)
			case constants.FEATURE_TITLE:
				featureSum += int64(item.Value)
			case constants.RISK_TITLE:
				riskSum += int64(item.Value)
			case constants.TECH_DEBT_TITLE:
				techDebtSum += int64(item.Value)
			}
		}

	}
	// Recursively sum counts for each sub-org
	for _, report := range compareReports.CompareReports {

		featureSum, riskSum, techDebtSum, bugSum = sumCounts(&report, featureSum, riskSum, techDebtSum, bugSum)

	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.BUGS_TITLE,
			Value: int(bugSum),
		}, {
			Title: constants.FEATURE_TITLE,
			Value: int(featureSum),
		}, {
			Title: constants.RISK_TITLE,
			Value: int(riskSum),
		},
		{
			Title: constants.TECH_DEBT_TITLE,
			Value: int(techDebtSum),
		},
	}

	compareReports.Section.Data = data
	return featureSum, riskSum, techDebtSum, bugSum
}

func sumTotalValue(compareReports *constants.CompareReports, totalValue int) int {
	// Sum counts of feature, risk, tech debt, and bugs for the current level
	if !compareReports.IsSubOrg {
		totalValue += compareReports.TotalValue

	}
	// Recursively sum counts for each sub-org
	for _, report := range compareReports.CompareReports {

		totalValue = sumTotalValue(&report, totalValue)

	}

	compareReports.TotalValue = totalValue
	return totalValue
}

func getCycleTimeComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["flowCycleTimeChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.CycleTimeComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getCycleTimeComponentComparison()") {
		return nil, err
	}

	compareReports := createCycleTimeCompareReports(organisation, &result)

	calculateFmPercentage(compareReports)

	// // To be removed - start
	// bData, err := json.MarshalIndent(compareReports.CompareReports, "", "    ")
	// if log.CheckErrorf(err, "Error marshaling compareReports : ") {
	// 	return nil, err
	// }
	// log.Infof(string(bData))
	// // To be removed - end

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getCycleTimeComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createCycleTimeCompareReports(orgData *constants.Organization, compComparison *constants.CycleTimeComponentComparison) *constants.CycleTimeCompareReports {
	var compareReports constants.CycleTimeCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalBugsValueInMillis, totalFeatureValueInMillis, totalRiskValueInMillis, totalTechDebtValueInMillis int
	var totalBugsCount, totalFeatureCount, totalRiskCount, totalTechDebtCount int

	// Count components
	for _, bucket := range compComparison.Aggregations.CycleTimeComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				totalBugsValueInMillis += bucket.FlowCycleTimeCount.Value.DefectTime
				totalFeatureValueInMillis += bucket.FlowCycleTimeCount.Value.FeatureTime
				totalRiskValueInMillis += bucket.FlowCycleTimeCount.Value.RiskTime
				totalTechDebtValueInMillis += bucket.FlowCycleTimeCount.Value.TechDebtTime
				totalBugsCount += bucket.FlowCycleTimeCount.Value.DefectCount
				totalFeatureCount += bucket.FlowCycleTimeCount.Value.FeatureCount
				totalRiskCount += bucket.FlowCycleTimeCount.Value.RiskCount
				totalTechDebtCount += bucket.FlowCycleTimeCount.Value.TechDebtCount

				data := []struct {
					Title string  `json:"title"`
					Value float64 `json:"value"`
					Time  int     `json:"time"`
					Count int     `json:"count"`
				}{
					{
						Title: constants.BUGS_TITLE,
						Time:  bucket.FlowCycleTimeCount.Value.DefectTime,
						Count: bucket.FlowCycleTimeCount.Value.DefectCount,
					}, {
						Title: constants.FEATURE_TITLE,
						Time:  bucket.FlowCycleTimeCount.Value.FeatureTime,
						Count: bucket.FlowCycleTimeCount.Value.FeatureCount,
					}, {
						Title: constants.RISK_TITLE,
						Time:  bucket.FlowCycleTimeCount.Value.RiskTime,
						Count: bucket.FlowCycleTimeCount.Value.RiskCount,
					}, {
						Title: constants.TECH_DEBT_TITLE,
						Time:  bucket.FlowCycleTimeCount.Value.TechDebtTime,
						Count: bucket.FlowCycleTimeCount.Value.TechDebtCount,
					},
				}

				if data[0].Count > 0 {
					data[0].Value = float64(data[0].Time) / float64(data[0].Count)
				} else {
					data[0].Value = 0
				}
				if data[1].Count > 0 {
					data[1].Value = float64(data[1].Time) / float64(data[1].Count)
				} else {
					data[1].Value = 0
				}
				if data[2].Count > 0 {
					data[2].Value = float64(data[2].Time) / float64(data[2].Count)
				} else {
					data[2].Value = 0
				}
				if data[3].Count > 0 {
					data[3].Value = float64(data[3].Time) / float64(data[3].Count)
				} else {
					data[3].Value = 0
				}

				valueInMillis := math.Round(data[0].Value + data[1].Value + data[2].Value + data[3].Value)

				compareReports.ComponentCount++
				var compCompareReports constants.CycleTimeCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.ValueInMillis = valueInMillis
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string  `json:"title"`
		Value float64 `json:"value"`
		Time  int     `json:"time"`
		Count int     `json:"count"`
	}{
		{
			Title: constants.BUGS_TITLE,
			Time:  totalBugsValueInMillis,
			Count: totalBugsCount,
		}, {
			Title: constants.FEATURE_TITLE,
			Time:  totalFeatureValueInMillis,
			Count: totalFeatureCount,
		}, {
			Title: constants.RISK_TITLE,
			Time:  totalRiskValueInMillis,
			Count: totalRiskCount,
		}, {
			Title: constants.TECH_DEBT_TITLE,
			Time:  totalTechDebtValueInMillis,
			Count: totalTechDebtCount,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createCycleTimeCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalBugsValueInMillis += subCompareReports.Section.Data[0].Time
		totalFeatureValueInMillis += subCompareReports.Section.Data[1].Time
		totalRiskValueInMillis += subCompareReports.Section.Data[2].Time
		totalTechDebtValueInMillis += subCompareReports.Section.Data[3].Time

		totalBugsCount += subCompareReports.Section.Data[0].Count
		totalFeatureCount += subCompareReports.Section.Data[1].Count
		totalRiskCount += subCompareReports.Section.Data[2].Count
		totalTechDebtCount += subCompareReports.Section.Data[3].Count
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Time = totalBugsValueInMillis
	compareReports.Section.Data[1].Time = totalFeatureValueInMillis
	compareReports.Section.Data[2].Time = totalRiskValueInMillis
	compareReports.Section.Data[3].Time = totalTechDebtValueInMillis

	compareReports.Section.Data[0].Count = totalBugsCount
	compareReports.Section.Data[1].Count = totalFeatureCount
	compareReports.Section.Data[2].Count = totalRiskCount
	compareReports.Section.Data[3].Count = totalTechDebtCount

	if totalBugsCount > 0 {
		compareReports.Section.Data[0].Value = float64(totalBugsValueInMillis) / float64(totalBugsCount)
	} else {
		compareReports.Section.Data[0].Value = float64(totalBugsValueInMillis)
	}

	if totalFeatureCount > 0 {
		compareReports.Section.Data[1].Value = float64(totalFeatureValueInMillis) / float64(totalFeatureCount)
	} else {
		compareReports.Section.Data[1].Value = float64(totalFeatureValueInMillis)
	}

	if totalRiskCount > 0 {
		compareReports.Section.Data[2].Value = float64(totalRiskValueInMillis) / float64(totalRiskCount)
	} else {
		compareReports.Section.Data[2].Value = float64(totalRiskValueInMillis)
	}

	if totalTechDebtCount > 0 {
		compareReports.Section.Data[3].Value = float64(totalTechDebtValueInMillis) / float64(totalTechDebtCount)
	} else {
		compareReports.Section.Data[3].Value = float64(totalTechDebtValueInMillis)
	}

	compareReports.ValueInMillis = math.Round(compareReports.Section.Data[0].Value + compareReports.Section.Data[1].Value + compareReports.Section.Data[2].Value + compareReports.Section.Data[3].Value)

	return &compareReports
}

func calculateFmPercentage(compareReports *constants.CycleTimeCompareReports) {

	// Recursively sum counts for each sub-org
	for _, report := range compareReports.CompareReports {
		if report.ValueInMillis > 0 {
			valueInMillis := float64(report.ValueInMillis)
			bugs := float64(report.Section.Data[0].Value)
			feature := float64(report.Section.Data[1].Value)
			risk := float64(report.Section.Data[2].Value)
			techDebt := float64(report.Section.Data[3].Value)

			// Calculate percentages
			bugsPercentage := math.Round((bugs / valueInMillis) * 100)
			featurePercentage := math.Round((feature / valueInMillis) * 100)
			riskPercentage := math.Round((risk / valueInMillis) * 100)
			techDebtPercentage := math.Round((techDebt / valueInMillis) * 100)

			// Check if sum exceeds 100%, and adjust percentages accordingly
			sumPercentages := bugsPercentage + featurePercentage + riskPercentage
			if sumPercentages > 100 {
				diff := sumPercentages - 100
				// Find the maximum percentage and reduce it by the difference
				maxPercentage := math.Max(math.Max(math.Max(float64(bugsPercentage), float64(featurePercentage)), float64(riskPercentage)), float64(techDebtPercentage))
				if maxPercentage == float64(bugsPercentage) {
					bugsPercentage -= diff
				} else if maxPercentage == float64(featurePercentage) {
					featurePercentage -= diff
				} else if maxPercentage == float64(riskPercentage) {
					riskPercentage -= diff
				} else {
					techDebtPercentage -= diff
				}
			}

			report.Section.Data[0].Value = bugsPercentage
			report.Section.Data[1].Value = featurePercentage
			report.Section.Data[2].Value = riskPercentage
			report.Section.Data[3].Value = techDebtPercentage
		} else {
			report.Section.Data[0].Value = 0
			report.Section.Data[1].Value = 0
			report.Section.Data[2].Value = 0
			report.Section.Data[3].Value = 0
		}

		calculateFmPercentage(&report)

	}
}

func getCommitTrendsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["commitsAndAverageChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.CommitsTrendsComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getCommitTrendsComponentComparison()") {
		return nil, err
	}

	compareReports := createCommitTrendsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getCommitTrendsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createCommitTrendsCompareReports(orgData *constants.Organization, compComparison *constants.CommitsTrendsComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	totalValue := 0

	// Count components
	for _, bucket := range compComparison.Aggregations.CommitsTrendsComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				value := bucket.CommitsCount.Value
				totalValue += value

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				// compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createCommitTrendsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++
		totalValue += subCompareReports.TotalValue

	}

	compareReports.TotalValue = totalValue

	return &compareReports
}

func getPullRequestComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["pullRequestsChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.PullRequestComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getPullRequestComponentComparison()") {
		return nil, err
	}

	compareReports := createPullRequestCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getPullRequestComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createPullRequestCompareReports(orgData *constants.Organization, compComparison *constants.PullRequestComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalApproved, totalChangesRequested, totalOpen, totalRejected int

	// Count components
	for _, bucket := range compComparison.Aggregations.PullRequestsComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {
				totalApproved += bucket.Pullrequests.Value.Approved
				totalChangesRequested += bucket.Pullrequests.Value.ChangesRequested
				totalOpen += bucket.Pullrequests.Value.Open
				totalRejected += bucket.Pullrequests.Value.Rejected

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.APPROVED_TITLE,
						Value: bucket.Pullrequests.Value.Approved,
					}, {
						Title: constants.CHANGES_REQUESTED_TITLE,
						Value: bucket.Pullrequests.Value.ChangesRequested,
					}, {
						Title: constants.OPEN_TITLE,
						Value: bucket.Pullrequests.Value.Open,
					},
					{
						Title: constants.REJECTED_TITLE,
						Value: bucket.Pullrequests.Value.Rejected,
					},
				}

				value := bucket.Pullrequests.Value.Approved + bucket.Pullrequests.Value.ChangesRequested + bucket.Pullrequests.Value.Open + bucket.Pullrequests.Value.Rejected

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.APPROVED_TITLE,
			Value: totalApproved,
		}, {
			Title: constants.CHANGES_REQUESTED_TITLE,
			Value: totalChangesRequested,
		}, {
			Title: constants.OPEN_TITLE,
			Value: totalOpen,
		},
		{
			Title: constants.REJECTED_TITLE,
			Value: totalRejected,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs and accumulate total values
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createPullRequestCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalApproved += subCompareReports.Section.Data[0].Value
		totalChangesRequested += subCompareReports.Section.Data[1].Value
		totalOpen += subCompareReports.Section.Data[2].Value
		totalRejected += subCompareReports.Section.Data[3].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalApproved
	compareReports.Section.Data[1].Value = totalChangesRequested
	compareReports.Section.Data[2].Value = totalOpen
	compareReports.Section.Data[3].Value = totalRejected

	compareReports.TotalValue = totalApproved + totalChangesRequested + totalOpen + totalRejected

	return &compareReports
}

func getWorkflowRunsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["runsStatusChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.WorkflowRunsComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getWorkflowRunsComponentComparison()") {
		return nil, err
	}

	compareReports := createWorkflowRunsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkflowRunsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createWorkflowRunsCompareReports(orgData *constants.Organization, compComparison *constants.WorkflowRunsComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalSuccess, totalFailure int

	// Count components
	for _, bucket := range compComparison.Aggregations.WorkflowRunsComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				totalSuccess += bucket.AutomationRun.Value.Success
				totalFailure += bucket.AutomationRun.Value.Failure

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: "Success",
						Value: bucket.AutomationRun.Value.Success,
					}, {
						Title: "Failure",
						Value: bucket.AutomationRun.Value.Failure,
					},
				}

				value := bucket.AutomationRun.Value.Success + bucket.AutomationRun.Value.Failure

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: "Success",
			Value: totalSuccess,
		}, {
			Title: "Failure",
			Value: totalFailure,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createWorkflowRunsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		totalSuccess += subCompareReports.Section.Data[0].Value
		totalFailure += subCompareReports.Section.Data[1].Value
	}

	compareReports.Section.Data[0].Value = totalSuccess
	compareReports.Section.Data[1].Value = totalFailure

	compareReports.TotalValue = totalSuccess + totalFailure

	return &compareReports
}

func getDevCycleTimeComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["developmentCycleChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.DevCycleTimeComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getDevCycleTimeComponentComparison()") {
		return nil, err
	}

	compareReports := createDevCycleTimeCompareReports(organisation, &result)

	calculatePercentage(compareReports)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getDevCycleTimeComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func calculatePercentage(compareReports *constants.DevCycleTimeCompareReports) {

	// Recursively sum counts for each sub-org
	for _, report := range compareReports.CompareReports {
		if report.ValueInMillis > 0 {

			valueInMillis := float64(report.ValueInMillis)
			codingTime := float64(report.Section.Data[0].Value)
			codePickupTime := float64(report.Section.Data[1].Value)
			codeReviewTime := float64(report.Section.Data[2].Value)

			// Calculate percentages
			codingTimePercentage := math.Round((codingTime / valueInMillis) * 100)
			codePickupTimePercentage := math.Round((codePickupTime / valueInMillis) * 100)
			codeReviewTimePercentage := math.Round((codeReviewTime / valueInMillis) * 100)

			// Check if sum exceeds 100%, and adjust percentages accordingly
			sumPercentages := codingTimePercentage + codePickupTimePercentage + codeReviewTimePercentage
			if sumPercentages > 100 {
				diff := sumPercentages - 100
				// Find the maximum percentage and reduce it by the difference
				maxPercentage := math.Max(math.Max(float64(codingTimePercentage), float64(codePickupTimePercentage)), float64(codeReviewTimePercentage))
				if maxPercentage == float64(codingTimePercentage) {
					codingTimePercentage -= diff
				} else if maxPercentage == float64(codePickupTimePercentage) {
					codePickupTimePercentage -= diff
				} else {
					codeReviewTimePercentage -= diff
				}
			}

			report.Section.Data[0].Value = codingTimePercentage
			report.Section.Data[1].Value = codePickupTimePercentage
			report.Section.Data[2].Value = codeReviewTimePercentage
		} else {
			report.Section.Data[0].Value = 0
			report.Section.Data[1].Value = 0
			report.Section.Data[2].Value = 0
		}

		calculatePercentage(&report)

	}
}

func createDevCycleTimeCompareReports(orgData *constants.Organization, compComparison *constants.DevCycleTimeComponentComparison) *constants.DevCycleTimeCompareReports {
	var compareReports constants.DevCycleTimeCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalCodingTimeValueInMillis, totalPickupTimeValueInMillis, totalReviewTimeValueInMillis int
	var totalCodingTimeCount, totalPickupTimeCount, totalReviewTimeCount int

	// Count components
	for _, bucket := range compComparison.Aggregations.DevelopmentCycleTimeComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				totalCodingTimeValueInMillis += bucket.DevelopmentCycleTime.Value.CodingTimeValueInMillis
				totalPickupTimeValueInMillis += bucket.DevelopmentCycleTime.Value.PickupTimeValueInMillis
				totalReviewTimeValueInMillis += bucket.DevelopmentCycleTime.Value.ReviewTimeValueInMillis
				totalCodingTimeCount += bucket.DevelopmentCycleTime.Value.CodingTimeCount
				totalPickupTimeCount += bucket.DevelopmentCycleTime.Value.PickupTimeCount
				totalReviewTimeCount += bucket.DevelopmentCycleTime.Value.ReviewTimeCount

				data := []struct {
					Title string  `json:"title"`
					Value float64 `json:"value"`
					Time  int     `json:"time"`
					Count int     `json:"count"`
				}{
					{
						Title: constants.COMPARE_REPORT_CODING_TIME,
						Time:  bucket.DevelopmentCycleTime.Value.CodingTimeValueInMillis,
						Count: bucket.DevelopmentCycleTime.Value.CodingTimeCount,
					}, {
						Title: constants.COMPARE_REPORT_CODE_PICKUP_TIME,
						Time:  bucket.DevelopmentCycleTime.Value.PickupTimeValueInMillis,
						Count: bucket.DevelopmentCycleTime.Value.PickupTimeCount,
					}, {
						Title: constants.COMPARE_REPORT_CODE_REVIEW_TIME,
						Time:  bucket.DevelopmentCycleTime.Value.ReviewTimeValueInMillis,
						Count: bucket.DevelopmentCycleTime.Value.ReviewTimeCount,
					},
				}

				if data[0].Count > 0 {
					data[0].Value = float64(data[0].Time) / float64(data[0].Count)
				} else {
					data[0].Value = float64(data[0].Time)
				}
				if data[1].Count > 0 {
					data[1].Value = float64(data[1].Time) / float64(data[1].Count)
				} else {
					data[1].Value = float64(data[1].Time)
				}
				if data[2].Count > 0 {
					data[2].Value = float64(data[2].Time) / float64(data[2].Count)
				} else {
					data[2].Value = float64(data[2].Time)
				}

				valueInMillis := data[0].Value + data[1].Value + data[2].Value

				compareReports.ComponentCount++
				var compCompareReports constants.DevCycleTimeCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.ValueInMillis = valueInMillis
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string  `json:"title"`
		Value float64 `json:"value"`
		Time  int     `json:"time"`
		Count int     `json:"count"`
	}{
		{
			Title: constants.SECTION_CODING_TIME,
			Time:  totalCodingTimeValueInMillis,
			Count: totalCodingTimeCount,
		}, {
			Title: constants.SECTION_CODE_PICKUP_TIME,
			Time:  totalPickupTimeValueInMillis,
			Count: totalPickupTimeCount,
		}, {
			Title: constants.SECTION_CODE_REVIEW_TIME,
			Time:  totalReviewTimeValueInMillis,
			Count: totalReviewTimeCount,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createDevCycleTimeCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalCodingTimeValueInMillis += subCompareReports.Section.Data[0].Time
		totalPickupTimeValueInMillis += subCompareReports.Section.Data[1].Time
		totalReviewTimeValueInMillis += subCompareReports.Section.Data[2].Time

		// Accumulate total values recursively
		totalCodingTimeCount += subCompareReports.Section.Data[0].Count
		totalPickupTimeCount += subCompareReports.Section.Data[1].Count
		totalReviewTimeCount += subCompareReports.Section.Data[2].Count
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Time = totalCodingTimeValueInMillis
	compareReports.Section.Data[1].Time = totalPickupTimeValueInMillis
	compareReports.Section.Data[2].Time = totalReviewTimeValueInMillis

	compareReports.Section.Data[0].Count = totalCodingTimeCount
	compareReports.Section.Data[1].Count = totalPickupTimeCount
	compareReports.Section.Data[2].Count = totalReviewTimeCount

	if totalCodingTimeCount > 0 {
		compareReports.Section.Data[0].Value = float64(totalCodingTimeValueInMillis) / float64(totalCodingTimeCount)
	} else {
		compareReports.Section.Data[0].Value = float64(totalCodingTimeValueInMillis)
	}

	if totalPickupTimeCount > 0 {
		compareReports.Section.Data[1].Value = float64(totalPickupTimeValueInMillis) / float64(totalPickupTimeCount)
	} else {
		compareReports.Section.Data[1].Value = float64(totalPickupTimeValueInMillis)
	}

	if totalReviewTimeCount > 0 {
		compareReports.Section.Data[2].Value = float64(totalReviewTimeValueInMillis) / float64(totalReviewTimeCount)
	} else {
		compareReports.Section.Data[2].Value = float64(totalReviewTimeValueInMillis)
	}

	compareReports.ValueInMillis = math.Round(compareReports.Section.Data[0].Value + compareReports.Section.Data[1].Value + compareReports.Section.Data[2].Value)

	return &compareReports
}

func getCommitsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["commitsChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.CommitsComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getCommitsComponentComparison()") {
		return nil, err
	}

	compareReports := createCommitsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getCommitsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createCommitsCompareReports(orgData *constants.Organization, compComparison *constants.CommitsComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	totalValue := 0

	// Count components
	for _, bucket := range compComparison.Aggregations.CommitsComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				totalValue += bucket.AutomationRun.Value.TotalCount
				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				// compCompareReports.Section.Data = data
				compCompareReports.TotalValue = bucket.AutomationRun.Value.TotalCount
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	compareReports.TotalValue = totalValue

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createCommitsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++
		totalValue += subCompareReports.TotalValue

	}
	compareReports.TotalValue = totalValue

	return &compareReports
}

func getBuildsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["buildsChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.BuildsComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getBuildsComponentComparison()") {
		return nil, err
	}

	compareReports := createBuildsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getBuildsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createBuildsCompareReports(orgData *constants.Organization, compComparison *constants.BuildsComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalSuccess, totalFailure int

	// Count components
	for _, bucket := range compComparison.Aggregations.BuildsComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				success := 0
				failure := 0

				for _, currentBucket := range bucket.BuildStatus.Value.Info {
					if currentBucket.Title == "Success" {
						totalSuccess += currentBucket.Value
						success = currentBucket.Value
					}
					if currentBucket.Title == "Failure" {
						totalFailure += currentBucket.Value
						failure = currentBucket.Value
					}
				}

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: "Success",
						Value: success,
					}, {
						Title: "Failure",
						Value: failure,
					},
				}

				value := success + failure

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: "Success",
			Value: totalSuccess,
		}, {
			Title: "Failure",
			Value: totalFailure,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createBuildsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalSuccess += subCompareReports.Section.Data[0].Value
		totalFailure += subCompareReports.Section.Data[1].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalSuccess
	compareReports.Section.Data[1].Value = totalFailure

	compareReports.TotalValue = totalSuccess + totalFailure

	return &compareReports
}

func getDeploymentsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["deploymentsChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.DeploymentsComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getDeploymentsComponentComparison()") {
		return nil, err
	}

	compareReports := createDeploymentsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getDeploymentsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createDeploymentsCompareReports(orgData *constants.Organization, compComparison *constants.DeploymentsComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalValue int

	// Count components
	for _, bucket := range compComparison.Aggregations.DeploymentsComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				totalValue += bucket.Deploys.Value

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.TotalValue = bucket.Deploys.Value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createDeploymentsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++
		totalValue += subCompareReports.TotalValue
	}
	compareReports.TotalValue = totalValue
	return &compareReports
}

func getComponentsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["automationsChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.ComponentsComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getComponentsComponentComparison()") {
		return nil, err
	}

	compareReports := createComponentsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getComponentsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createComponentsCompareReports(orgData *constants.Organization, compComparison *constants.ComponentsComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalActive, totalInactive int

	// Count components
	for _, component := range orgData.Components {
		found := false
		for _, componentId := range compComparison.Aggregations.DistinctComponent.Value {
			if component.ID == componentId {
				found = true
				totalActive += 1

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: "Active",
						Value: 1,
					}, {
						Title: "Inactive",
						Value: 0,
					},
				}

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = 1
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
		if !found {
			totalInactive += 1

			data := []struct {
				Title string `json:"title"`
				Value int    `json:"value"`
			}{
				{
					Title: "Active",
					Value: 0,
				}, {
					Title: "Inactive",
					Value: 1,
				},
			}

			compareReports.ComponentCount++
			var compCompareReports constants.CompareReports

			compCompareReports.SubOrgID = component.ID
			compCompareReports.CompareTitle = component.Name
			compCompareReports.IsSubOrg = false
			compCompareReports.ComponentCount = 1
			compCompareReports.Section.Data = data
			compCompareReports.TotalValue = 1
			compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: "Active",
			Value: totalActive,
		}, {
			Title: "Inactive",
			Value: totalInactive,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createComponentsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalActive += subCompareReports.Section.Data[0].Value
		totalInactive += subCompareReports.Section.Data[1].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalActive
	compareReports.Section.Data[1].Value = totalInactive

	compareReports.TotalValue = totalActive + totalInactive

	return &compareReports
}

func calculateAverage(activeTime, flowTime int) float64 {
	if flowTime > 0 {
		return math.Round(float64(activeTime) / float64(flowTime))
	}
	return 0
}

func calculatePercentageValue(input1, input2 int) int {
	if input2 > 0 {
		return int(math.Round(float64(input1) / float64(input2) * 100))
	}
	return 0
}

func getWorkflowsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["workflow component comparison data"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	var result map[string]constants.WorkflowsComponentComparison

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getWorkflowsComponentComparison()") {
		return nil, err
	}

	compareReports := createWorkflowsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkflowsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createWorkflowsCompareReports(orgData *constants.Organization, compComparison *map[string]constants.WorkflowsComponentComparison) *constants.CompareReports {

	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalActive, totalInactive int

	// Count components
	for key, value := range *compComparison {
		for _, component := range orgData.Components {
			if component.ID == key {

				totalActive += value.Active
				totalInactive += value.Inactive

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: "Active",
						Value: value.Active,
					}, {
						Title: "Inactive",
						Value: value.Inactive,
					},
				}

				value := value.Active + value.Inactive

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: "Active",
			Value: totalActive,
		}, {
			Title: "Inactive",
			Value: totalInactive,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createWorkflowsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalActive += subCompareReports.Section.Data[0].Value
		totalInactive += subCompareReports.Section.Data[1].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalActive
	compareReports.Section.Data[1].Value = totalInactive

	compareReports.TotalValue = totalActive + totalInactive

	return &compareReports
}

func transformTestWorkflowsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {
	response, ok := data["test workflow component comparison data"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	var result map[string]constants.TestWorkflowsComponentComparison

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in transformTestWorkflowsComponentComparison()") {
		return nil, err
	}

	compareReports := createTestWorkflowsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in transformTestWorkflowsComponentComparison() chart") {
		return nil, err
	}
	return b, err

}

func createTestWorkflowsCompareReports(orgData *constants.Organization, compComparison *map[string]constants.TestWorkflowsComponentComparison) *constants.CompareReports {

	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalWithTestSuites, totalWithoutTestSuites int

	// Count components
	for key, value := range *compComparison {
		for _, component := range orgData.Components {
			if component.ID == key {

				totalWithTestSuites += value.WithTestSuites
				totalWithoutTestSuites += value.WithoutTestSuites

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.WITH_TEST_SUITES,
						Value: value.WithTestSuites,
					}, {
						Title: constants.WITHOUT_TEST_SUITES,
						Value: value.WithoutTestSuites,
					},
				}

				value := value.WithTestSuites + value.WithoutTestSuites

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.WITH_TEST_SUITES,
			Value: totalWithTestSuites,
		}, {
			Title: constants.WITHOUT_TEST_SUITES,
			Value: totalWithoutTestSuites,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createTestWorkflowsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalWithTestSuites += subCompareReports.Section.Data[0].Value
		totalWithoutTestSuites += subCompareReports.Section.Data[1].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalWithTestSuites
	compareReports.Section.Data[1].Value = totalWithoutTestSuites

	compareReports.TotalValue = totalWithTestSuites + totalWithoutTestSuites

	return &compareReports
}

func getSecurityWorkflowsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["security workflow component comparison data"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	var result map[string]constants.SecurityWorkflowsComponentComparison

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getSecurityWorkflowsComponentComparison()") {
		return nil, err
	}

	compareReports := createSecurityWorkflowsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getSecurityWorkflowsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createSecurityWorkflowsCompareReports(orgData *constants.Organization, compComparison *map[string]constants.SecurityWorkflowsComponentComparison) *constants.CompareReports {

	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalWithScanners, totalWithoutScanners int

	// Count components
	for key, value := range *compComparison {
		for _, component := range orgData.Components {
			if component.ID == key {

				totalWithScanners += value.WithScanners
				totalWithoutScanners += value.WithoutScanners

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.WITHSCANNERS,
						Value: value.WithScanners,
					}, {
						Title: constants.WITHOUTSCANNERS,
						Value: value.WithoutScanners,
					},
				}

				value := value.WithScanners + value.WithoutScanners

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.WITHSCANNERS,
			Value: totalWithScanners,
		}, {
			Title: constants.WITHOUTSCANNERS,
			Value: totalWithoutScanners,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createSecurityWorkflowsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalWithScanners += subCompareReports.Section.Data[0].Value
		totalWithoutScanners += subCompareReports.Section.Data[1].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalWithScanners
	compareReports.Section.Data[1].Value = totalWithoutScanners

	compareReports.TotalValue = totalWithScanners + totalWithoutScanners

	return &compareReports
}

func getActiveFlowTimeComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["flowEfficiencyChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.ActiveWorkTimeComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getActiveFlowTimeComponentComparison()") {
		return nil, err
	}

	compareReports := createActiveFlowTimeCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getActiveFlowTimeComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createActiveFlowTimeCompareReports(orgData *constants.Organization, compComparison *constants.ActiveWorkTimeComponentComparison) *constants.ActiveTimeCompareReports {
	var compareReports constants.ActiveTimeCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var activeBugs, activeFeature, activeRisk, activeTechDebt int

	var totalBugs, totalFeature, totalRisk, totalTechDebt int

	// Count components
	for _, bucket := range compComparison.Aggregations.FlowVelocityComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				activeBugs += int(bucket.FlowEfficiencyCount.Value.Defect.ActiveTime)
				activeFeature += int(bucket.FlowEfficiencyCount.Value.Feature.ActiveTime)
				activeRisk += int(bucket.FlowEfficiencyCount.Value.Risk.ActiveTime)
				activeTechDebt += int(bucket.FlowEfficiencyCount.Value.TechDebt.ActiveTime)

				totalBugs += int(bucket.FlowEfficiencyCount.Value.Defect.FlowTime)
				totalFeature += int(bucket.FlowEfficiencyCount.Value.Feature.FlowTime)
				totalRisk += int(bucket.FlowEfficiencyCount.Value.Risk.FlowTime)
				totalTechDebt += int(bucket.FlowEfficiencyCount.Value.TechDebt.FlowTime)

				data := []struct {
					Title      string `json:"title"`
					Value      int    `json:"value"`
					ActiveTime int    `json:"active_time"`
					FlowTime   int    `json:"flow_time"`
				}{
					{
						Title:      constants.BUGS_TITLE,
						Value:      calculateValue(bucket.FlowEfficiencyCount.Value.Defect.ActiveTime, bucket.FlowEfficiencyCount.Value.Defect.FlowTime),
						ActiveTime: bucket.FlowEfficiencyCount.Value.Defect.ActiveTime,
						FlowTime:   bucket.FlowEfficiencyCount.Value.Defect.FlowTime,
					}, {
						Title:      constants.FEATURE_TITLE,
						Value:      calculateValue(bucket.FlowEfficiencyCount.Value.Feature.ActiveTime, bucket.FlowEfficiencyCount.Value.Feature.FlowTime),
						ActiveTime: bucket.FlowEfficiencyCount.Value.Feature.ActiveTime,
						FlowTime:   bucket.FlowEfficiencyCount.Value.Feature.FlowTime,
					}, {
						Title:      constants.RISK_TITLE,
						Value:      calculateValue(bucket.FlowEfficiencyCount.Value.Risk.ActiveTime, bucket.FlowEfficiencyCount.Value.Risk.FlowTime),
						ActiveTime: bucket.FlowEfficiencyCount.Value.Risk.ActiveTime,
						FlowTime:   bucket.FlowEfficiencyCount.Value.Risk.FlowTime,
					}, {
						Title:      constants.TECH_DEBT_TITLE,
						Value:      calculateValue(bucket.FlowEfficiencyCount.Value.TechDebt.ActiveTime, bucket.FlowEfficiencyCount.Value.TechDebt.FlowTime),
						ActiveTime: bucket.FlowEfficiencyCount.Value.TechDebt.ActiveTime,
						FlowTime:   bucket.FlowEfficiencyCount.Value.TechDebt.FlowTime,
					},
				}

				activeTime := bucket.FlowEfficiencyCount.Value.Defect.ActiveTime + bucket.FlowEfficiencyCount.Value.Feature.ActiveTime + bucket.FlowEfficiencyCount.Value.Risk.ActiveTime + bucket.FlowEfficiencyCount.Value.TechDebt.ActiveTime
				flowTime := bucket.FlowEfficiencyCount.Value.Defect.FlowTime + bucket.FlowEfficiencyCount.Value.Feature.FlowTime + bucket.FlowEfficiencyCount.Value.Risk.FlowTime + bucket.FlowEfficiencyCount.Value.TechDebt.FlowTime

				compareReports.ComponentCount++
				var compCompareReports constants.ActiveTimeCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = calculateValue(activeTime, flowTime)
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title      string `json:"title"`
		Value      int    `json:"value"`
		ActiveTime int    `json:"active_time"`
		FlowTime   int    `json:"flow_time"`
	}{
		{
			Title:      constants.BUGS_TITLE,
			Value:      calculateValue(activeBugs, totalBugs),
			ActiveTime: activeBugs,
			FlowTime:   totalBugs,
		}, {
			Title:      constants.FEATURE_TITLE,
			Value:      calculateValue(activeFeature, totalFeature),
			ActiveTime: activeFeature,
			FlowTime:   totalFeature,
		}, {
			Title:      constants.RISK_TITLE,
			Value:      calculateValue(activeRisk, totalRisk),
			ActiveTime: activeRisk,
			FlowTime:   totalRisk,
		}, {
			Title:      constants.TECH_DEBT_TITLE,
			Value:      calculateValue(activeTechDebt, totalTechDebt),
			ActiveTime: activeTechDebt,
			FlowTime:   totalTechDebt,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createActiveFlowTimeCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively

		activeBugs += subCompareReports.Section.Data[0].ActiveTime
		activeFeature += subCompareReports.Section.Data[1].ActiveTime
		activeRisk += subCompareReports.Section.Data[2].ActiveTime
		activeTechDebt += subCompareReports.Section.Data[3].ActiveTime

		totalBugs += subCompareReports.Section.Data[0].FlowTime
		totalFeature += subCompareReports.Section.Data[1].FlowTime
		totalRisk += subCompareReports.Section.Data[2].FlowTime
		totalTechDebt += subCompareReports.Section.Data[3].FlowTime
	}

	compareReports.Section.Data[0].Value = calculateValue(activeBugs, totalBugs)
	compareReports.Section.Data[1].Value = calculateValue(activeFeature, totalFeature)
	compareReports.Section.Data[2].Value = calculateValue(activeRisk, totalRisk)
	compareReports.Section.Data[3].Value = calculateValue(activeTechDebt, totalTechDebt)

	totalActiveTime := activeBugs + activeFeature + activeRisk + activeTechDebt
	totalFlowTime := totalBugs + totalFeature + totalRisk + totalTechDebt
	compareReports.TotalValue = calculateValue(totalActiveTime, totalFlowTime)

	return &compareReports
}

func calculateValue(activeTime, flowTime int) int {
	if flowTime > 0 {
		return int(math.Round((float64(activeTime) / float64(flowTime)) * 100))
	}
	return 0
}

func getWorkWaitTimeComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["flowWaitTimeChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.WorkWaitTimeComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getWorkWaitTimeComponentComparison()") {
		return nil, err
	}

	compareReports := createWorkWaitTimeCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkWaitTimeComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createWorkWaitTimeCompareReports(orgData *constants.Organization, compComparison *constants.WorkWaitTimeComponentComparison) *constants.ActiveTimeCompareReports {
	var compareReports constants.ActiveTimeCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var activeBugs, activeFeature, activeRisk, activeTechDebt int

	var totalBugs, totalFeature, totalRisk, totalTechDebt int

	// Count components
	for _, bucket := range compComparison.Aggregations.FlowVelocityComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				activeBugs += int(bucket.FlowEfficiencyCount.Value.Defect.WaitingTime)
				activeFeature += int(bucket.FlowEfficiencyCount.Value.Feature.WaitingTime)
				activeRisk += int(bucket.FlowEfficiencyCount.Value.Risk.WaitingTime)
				activeTechDebt += int(bucket.FlowEfficiencyCount.Value.TechDebt.WaitingTime)

				totalBugs += int(bucket.FlowEfficiencyCount.Value.Defect.FlowTime)
				totalFeature += int(bucket.FlowEfficiencyCount.Value.Feature.FlowTime)
				totalRisk += int(bucket.FlowEfficiencyCount.Value.Risk.FlowTime)
				totalTechDebt += int(bucket.FlowEfficiencyCount.Value.TechDebt.FlowTime)

				data := []struct {
					Title      string `json:"title"`
					Value      int    `json:"value"`
					ActiveTime int    `json:"active_time"`
					FlowTime   int    `json:"flow_time"`
				}{
					{
						Title:      constants.BUGS_TITLE,
						Value:      calculateValue(bucket.FlowEfficiencyCount.Value.Defect.WaitingTime, bucket.FlowEfficiencyCount.Value.Defect.FlowTime),
						ActiveTime: bucket.FlowEfficiencyCount.Value.Defect.WaitingTime,
						FlowTime:   bucket.FlowEfficiencyCount.Value.Defect.WaitingTime,
					}, {
						Title:      constants.FEATURE_TITLE,
						Value:      calculateValue(bucket.FlowEfficiencyCount.Value.Feature.WaitingTime, bucket.FlowEfficiencyCount.Value.Feature.FlowTime),
						ActiveTime: bucket.FlowEfficiencyCount.Value.Feature.WaitingTime,
						FlowTime:   bucket.FlowEfficiencyCount.Value.Feature.WaitingTime,
					}, {
						Title:      constants.RISK_TITLE,
						Value:      calculateValue(bucket.FlowEfficiencyCount.Value.Risk.WaitingTime, bucket.FlowEfficiencyCount.Value.Risk.FlowTime),
						ActiveTime: bucket.FlowEfficiencyCount.Value.Risk.WaitingTime,
						FlowTime:   bucket.FlowEfficiencyCount.Value.Risk.WaitingTime,
					}, {
						Title:      constants.TECH_DEBT_TITLE,
						Value:      calculateValue(bucket.FlowEfficiencyCount.Value.TechDebt.WaitingTime, bucket.FlowEfficiencyCount.Value.TechDebt.FlowTime),
						ActiveTime: bucket.FlowEfficiencyCount.Value.TechDebt.WaitingTime,
						FlowTime:   bucket.FlowEfficiencyCount.Value.TechDebt.WaitingTime,
					},
				}

				activeTime := bucket.FlowEfficiencyCount.Value.Defect.WaitingTime + bucket.FlowEfficiencyCount.Value.Feature.WaitingTime + bucket.FlowEfficiencyCount.Value.Risk.WaitingTime + bucket.FlowEfficiencyCount.Value.TechDebt.WaitingTime
				flowTime := bucket.FlowEfficiencyCount.Value.Defect.FlowTime + bucket.FlowEfficiencyCount.Value.Feature.FlowTime + bucket.FlowEfficiencyCount.Value.Risk.FlowTime + bucket.FlowEfficiencyCount.Value.TechDebt.FlowTime

				compareReports.ComponentCount++
				var compCompareReports constants.ActiveTimeCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = calculateValue(activeTime, flowTime)
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title      string `json:"title"`
		Value      int    `json:"value"`
		ActiveTime int    `json:"active_time"`
		FlowTime   int    `json:"flow_time"`
	}{
		{
			Title:      constants.BUGS_TITLE,
			Value:      calculateValue(activeBugs, totalBugs),
			ActiveTime: activeBugs,
			FlowTime:   totalBugs,
		}, {
			Title:      constants.FEATURE_TITLE,
			Value:      calculateValue(activeFeature, totalFeature),
			ActiveTime: activeFeature,
			FlowTime:   totalFeature,
		}, {
			Title:      constants.RISK_TITLE,
			Value:      calculateValue(activeRisk, totalRisk),
			ActiveTime: activeRisk,
			FlowTime:   totalRisk,
		}, {
			Title:      constants.TECH_DEBT_TITLE,
			Value:      calculateValue(activeTechDebt, totalTechDebt),
			ActiveTime: activeTechDebt,
			FlowTime:   totalTechDebt,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createWorkWaitTimeCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively

		activeBugs += subCompareReports.Section.Data[0].ActiveTime
		activeFeature += subCompareReports.Section.Data[1].ActiveTime
		activeRisk += subCompareReports.Section.Data[2].ActiveTime
		activeTechDebt += subCompareReports.Section.Data[3].ActiveTime

		totalBugs += subCompareReports.Section.Data[0].FlowTime
		totalFeature += subCompareReports.Section.Data[1].FlowTime
		totalRisk += subCompareReports.Section.Data[2].FlowTime
		totalTechDebt += subCompareReports.Section.Data[3].FlowTime
	}

	compareReports.Section.Data[0].Value = calculateValue(activeBugs, totalBugs)
	compareReports.Section.Data[1].Value = calculateValue(activeFeature, totalFeature)
	compareReports.Section.Data[2].Value = calculateValue(activeRisk, totalRisk)
	compareReports.Section.Data[3].Value = calculateValue(activeTechDebt, totalTechDebt)

	totalActiveTime := activeBugs + activeFeature + activeRisk + activeTechDebt
	totalFlowTime := totalBugs + totalFeature + totalRisk + totalTechDebt
	compareReports.TotalValue = calculateValue(totalActiveTime, totalFlowTime)

	return &compareReports
}

func getWorkloadComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["flowWorkLoad"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.WorkloadComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getWorkloadComponentComparison()") {
		return nil, err
	}

	compareReports := createWorkloadCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getWorkloadComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createWorkloadCompareReports(orgData *constants.Organization, compComparison *constants.WorkloadComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalBugs, totalFeature, totalRisk, totalTechDebt int

	// Count components
	for _, bucket := range compComparison.Aggregations.WorkloadComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				for _, flowItem := range bucket.WorkLoadCounts.Value.Dates {
					totalBugs += flowItem.Defect
					totalFeature += flowItem.Feature
					totalRisk += flowItem.Risk
					totalTechDebt += flowItem.TechDebt

					data := []struct {
						Title string `json:"title"`
						Value int    `json:"value"`
					}{
						{
							Title: "Bugs",
							Value: flowItem.Defect,
						}, {
							Title: "Feature",
							Value: flowItem.Feature,
						}, {
							Title: "Risk",
							Value: flowItem.Risk,
						},
						{
							Title: "Tech Debt",
							Value: flowItem.TechDebt,
						},
					}

					value := flowItem.Defect + flowItem.Feature + flowItem.Risk + flowItem.TechDebt

					compareReports.ComponentCount++
					var compCompareReports constants.CompareReports

					compCompareReports.SubOrgID = component.ID
					compCompareReports.CompareTitle = component.Name
					compCompareReports.IsSubOrg = false
					compCompareReports.ComponentCount = 1
					compCompareReports.Section.Data = data
					compCompareReports.TotalValue = value
					compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
					break
				}
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.BUGS_TITLE,
			Value: totalBugs,
		}, {
			Title: constants.FEATURE_TITLE,
			Value: totalFeature,
		}, {
			Title: constants.RISK_TITLE,
			Value: totalRisk,
		},
		{
			Title: constants.TECH_DEBT_TITLE,
			Value: totalTechDebt,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createWorkloadCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalBugs += subCompareReports.Section.Data[0].Value
		totalFeature += subCompareReports.Section.Data[1].Value
		totalRisk += subCompareReports.Section.Data[2].Value
		totalTechDebt += subCompareReports.Section.Data[3].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalBugs
	compareReports.Section.Data[1].Value = totalFeature
	compareReports.Section.Data[2].Value = totalRisk
	compareReports.Section.Data[3].Value = totalTechDebt

	compareReports.TotalValue = totalBugs + totalFeature + totalRisk + totalTechDebt

	return &compareReports
}

func getSecurityWorkflowRunsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["runsStatusChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.SecurityWorkflowRunsComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getSecurityWorkflowRunsComponentComparison()") {
		return nil, err
	}

	compareReports := createSecurityWorkflowRunsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getSecurityWorkflowRunsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createSecurityWorkflowRunsCompareReports(orgData *constants.Organization, compComparison *constants.SecurityWorkflowRunsComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalScanners, totalWithoutScanners int

	// Count components
	for _, bucket := range compComparison.Aggregations.WorkflowRunsComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {
				var scannersCount, withoutScannersCount int
				for _, currentBucket := range bucket.RunStatus.Value.ChartData.Info {
					switch currentBucket.Title {
					case constants.WITHSCANNERS:
						totalScanners += currentBucket.Value
						scannersCount = currentBucket.Value
					case constants.WITHOUTSCANNERS:
						totalWithoutScanners += currentBucket.Value
						withoutScannersCount = currentBucket.Value
					}
				}

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.WITHSCANNERS,
						Value: scannersCount,
					}, {
						Title: constants.WITHOUTSCANNERS,
						Value: withoutScannersCount,
					},
				}

				value := scannersCount + withoutScannersCount

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.WITHSCANNERS,
			Value: totalScanners,
		}, {
			Title: constants.WITHOUTSCANNERS,
			Value: totalWithoutScanners,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createSecurityWorkflowRunsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		totalScanners += subCompareReports.Section.Data[0].Value
		totalWithoutScanners += subCompareReports.Section.Data[1].Value
	}

	compareReports.Section.Data[0].Value = totalScanners
	compareReports.Section.Data[1].Value = totalWithoutScanners

	compareReports.TotalValue = totalScanners + totalWithoutScanners

	return &compareReports
}

func getSecurityComponentsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["security widget component comparison"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := []string{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getSecurityComponentsComponentComparison() : %+v : ", response) {
		return nil, err
	}

	compareReports := createSecurityComponentsCompareReports(organisation, result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getSecurityComponentsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createSecurityComponentsCompareReports(orgData *constants.Organization, components []string) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var withScannerTotal, withoutScannerTotal int

	// Count components
	for _, component := range orgData.Components {
		found := false
		for _, componentId := range components {
			if component.ID == componentId {
				found = true
				withScannerTotal += 1

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.WITHSCANNERS,
						Value: 1,
					}, {
						Title: constants.WITHOUTSCANNERS,
						Value: 0,
					},
				}

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = 1
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
		if !found {
			withoutScannerTotal += 1

			data := []struct {
				Title string `json:"title"`
				Value int    `json:"value"`
			}{
				{
					Title: constants.WITHSCANNERS,
					Value: 0,
				}, {
					Title: constants.WITHOUTSCANNERS,
					Value: 1,
				},
			}

			compareReports.ComponentCount++
			var compCompareReports constants.CompareReports

			compCompareReports.SubOrgID = component.ID
			compCompareReports.CompareTitle = component.Name
			compCompareReports.IsSubOrg = false
			compCompareReports.ComponentCount = 1
			compCompareReports.Section.Data = data
			compCompareReports.TotalValue = 1
			compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.WITHSCANNERS,
			Value: withScannerTotal,
		}, {
			Title: constants.WITHOUTSCANNERS,
			Value: withoutScannerTotal,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createSecurityComponentsCompareReports(subOrg, components)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		withScannerTotal += subCompareReports.Section.Data[0].Value
		withoutScannerTotal += subCompareReports.Section.Data[1].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = withScannerTotal
	compareReports.Section.Data[1].Value = withoutScannerTotal

	compareReports.TotalValue = withScannerTotal + withoutScannerTotal

	return &compareReports
}

func getMttrVeryHighComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["mttrVeryHighChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.MttrComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getMttrVeryHighComponentComparison()") {
		return nil, err
	}

	compareReports := createMttrVeryHighCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getMttrVeryHighComponentComparison() chart") {
		return nil, err
	}
	return b, err

}

func createMttrVeryHighCompareReports(orgData *constants.Organization, compComparison *constants.MttrComponentComparison) *constants.MttrCompareReports {
	var compareReports constants.MttrCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var veryHighAvg int

	var veryHightCount int

	// Count components
	for _, bucket := range compComparison.Aggregations.MttrComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				veryHighAvg += int(bucket.AvgTTR.Value.VeryHigh)
				veryHightCount++

				compareReports.ComponentCount++
				var compCompareReports constants.MttrCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.AverageValue = bucket.AvgTTR.Value.VeryHigh
				compCompareReports.Count = 1
				compCompareReports.ValueInMillis = float64(bucket.AvgTTR.Value.VeryHigh)
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createMttrVeryHighCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively

		veryHighAvg += subCompareReports.AverageValue
		veryHightCount += subCompareReports.Count
	}

	compareReports.ValueInMillis = calculateAverage(veryHighAvg, veryHightCount)

	return &compareReports
}

func getMttrHighComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["mttrHighChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.MttrComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getMttrHighComponentComparison()") {
		return nil, err
	}

	compareReports := createMttrHighCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getMttrHighComponentComparison() chart") {
		return nil, err
	}
	return b, err

}

func createMttrHighCompareReports(orgData *constants.Organization, compComparison *constants.MttrComponentComparison) *constants.MttrCompareReports {
	var compareReports constants.MttrCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var HighAvg int

	var HightCount int

	// Count components
	for _, bucket := range compComparison.Aggregations.MttrComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				HighAvg += int(bucket.AvgTTR.Value.High)
				HightCount++

				compareReports.ComponentCount++
				var compCompareReports constants.MttrCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.AverageValue = bucket.AvgTTR.Value.High
				compCompareReports.Count = 1
				compCompareReports.ValueInMillis = float64(bucket.AvgTTR.Value.High)
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createMttrHighCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively

		HighAvg += subCompareReports.AverageValue
		HightCount += subCompareReports.Count
	}

	compareReports.ValueInMillis = calculateAverage(HighAvg, HightCount)

	return &compareReports
}

func getMttrMediumComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["mttrMediumChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.MttrComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getMttrMediumComponentComparison()") {
		return nil, err
	}

	compareReports := createMttrMediumCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getMttrMediumComponentComparison() chart") {
		return nil, err
	}
	return b, err

}

func createMttrMediumCompareReports(orgData *constants.Organization, compComparison *constants.MttrComponentComparison) *constants.MttrCompareReports {
	var compareReports constants.MttrCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var mediumAvg int

	var mediumCount int

	// Count components
	for _, bucket := range compComparison.Aggregations.MttrComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				mediumAvg += int(bucket.AvgTTR.Value.Medium)
				mediumCount++

				compareReports.ComponentCount++
				var compCompareReports constants.MttrCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.AverageValue = bucket.AvgTTR.Value.Medium
				compCompareReports.Count = 1
				compCompareReports.ValueInMillis = float64(bucket.AvgTTR.Value.Medium)
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createMttrMediumCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively

		mediumAvg += subCompareReports.AverageValue
		mediumCount += subCompareReports.Count
	}

	compareReports.ValueInMillis = calculateAverage(mediumAvg, mediumCount)

	return &compareReports
}

func getMttrLowComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["mttrLowChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.MttrComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getMttrLowComponentComparison()") {
		return nil, err
	}

	compareReports := createMttrLowCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getMttrLowComponentComparison() chart") {
		return nil, err
	}
	return b, err

}

func createMttrLowCompareReports(orgData *constants.Organization, compComparison *constants.MttrComponentComparison) *constants.MttrCompareReports {
	var compareReports constants.MttrCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var lowAvg int

	var lowCount int

	// Count components
	for _, bucket := range compComparison.Aggregations.MttrComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				lowAvg += int(bucket.AvgTTR.Value.Low)
				lowCount++

				compareReports.ComponentCount++
				var compCompareReports constants.MttrCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.AverageValue = bucket.AvgTTR.Value.Low
				compCompareReports.Count = 1
				compCompareReports.ValueInMillis = float64(bucket.AvgTTR.Value.Low)
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createMttrLowCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively

		lowAvg += subCompareReports.AverageValue
		lowCount += subCompareReports.Count
	}

	compareReports.ValueInMillis = calculateAverage(lowAvg, lowCount)

	return &compareReports
}

func getSastVulnerabilitiesComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {
	response, ok := data["sastVulnerabilityScannerChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.VulnerabilitesByScannerType{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getSastVulnerabilitiesComponentComparison()") {
		return nil, err
	}

	compareReports := createSastVulnerabilitiesCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getSastVulnerabilitiesComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createSastVulnerabilitiesCompareReports(orgData *constants.Organization, compComparison *constants.VulnerabilitesByScannerType) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalVeryHigh, totalHigh, totalMedium, totalLow int

	// Count components
	for _, bucket := range compComparison.Aggregations.VulByScannerTypeComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {
				var veryHigh, high, medium, low int
				if len(bucket.VulByScannerTypeCounts.Value.VeryHigh) > 0 {
					veryHigh = bucket.VulByScannerTypeCounts.Value.VeryHigh[0].Y
				} else {
					veryHigh = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.High) > 0 {
					high = bucket.VulByScannerTypeCounts.Value.High[0].Y
				} else {
					high = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.Medium) > 0 {
					medium = bucket.VulByScannerTypeCounts.Value.Medium[0].Y
				} else {
					medium = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.Low) > 0 {
					low = bucket.VulByScannerTypeCounts.Value.Low[0].Y
				} else {
					low = 0
				}

				totalVeryHigh += veryHigh
				totalHigh += high
				totalMedium += medium
				totalLow += low

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.VERY_HIGH_TITLE,
						Value: veryHigh,
					}, {
						Title: constants.HIGH_TITLE,
						Value: high,
					}, {
						Title: constants.MEDIUM_TITLE,
						Value: medium,
					},
					{
						Title: constants.LOW_TITLE,
						Value: low,
					},
				}

				value := veryHigh + high + medium + low

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.VERY_HIGH_TITLE,
			Value: totalVeryHigh,
		}, {
			Title: constants.HIGH_TITLE,
			Value: totalHigh,
		}, {
			Title: constants.MEDIUM_TITLE,
			Value: totalMedium,
		},
		{
			Title: constants.LOW_TITLE,
			Value: totalLow,
		},
	}
	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createSastVulnerabilitiesCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalVeryHigh += subCompareReports.Section.Data[0].Value
		totalHigh += subCompareReports.Section.Data[1].Value
		totalMedium += subCompareReports.Section.Data[2].Value
		totalLow += subCompareReports.Section.Data[3].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalVeryHigh
	compareReports.Section.Data[1].Value = totalHigh
	compareReports.Section.Data[2].Value = totalMedium
	compareReports.Section.Data[3].Value = totalLow

	compareReports.TotalValue = totalVeryHigh + totalHigh + totalMedium + totalLow

	return &compareReports
}

func getDastVulnerabilitiesComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {
	response, ok := data["dastVulnerabilityScannerChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.VulnerabilitesByScannerType{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getDastVulnerabilitiesComponentComparison()") {
		return nil, err
	}

	compareReports := createDastVulnerabilitiesCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getDastVulnerabilitiesComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createDastVulnerabilitiesCompareReports(orgData *constants.Organization, compComparison *constants.VulnerabilitesByScannerType) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalVeryHigh, totalHigh, totalMedium, totalLow int

	// Count components
	for _, bucket := range compComparison.Aggregations.VulByScannerTypeComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				var veryHigh, high, medium, low int
				if len(bucket.VulByScannerTypeCounts.Value.VeryHigh) > 0 {
					veryHigh = bucket.VulByScannerTypeCounts.Value.VeryHigh[1].Y
				} else {
					veryHigh = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.High) > 0 {
					high = bucket.VulByScannerTypeCounts.Value.High[1].Y
				} else {
					high = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.Medium) > 0 {
					medium = bucket.VulByScannerTypeCounts.Value.Medium[1].Y
				} else {
					medium = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.Low) > 0 {
					low = bucket.VulByScannerTypeCounts.Value.Low[1].Y
				} else {
					low = 0
				}

				totalVeryHigh += veryHigh
				totalHigh += high
				totalMedium += medium
				totalLow += low

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.VERY_HIGH_TITLE,
						Value: veryHigh,
					}, {
						Title: constants.HIGH_TITLE,
						Value: high,
					}, {
						Title: constants.MEDIUM_TITLE,
						Value: medium,
					},
					{
						Title: constants.LOW_TITLE,
						Value: low,
					},
				}

				value := veryHigh + high + medium + low

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.VERY_HIGH_TITLE,
			Value: totalVeryHigh,
		}, {
			Title: constants.HIGH_TITLE,
			Value: totalHigh,
		}, {
			Title: constants.MEDIUM_TITLE,
			Value: totalMedium,
		},
		{
			Title: constants.LOW_TITLE,
			Value: totalLow,
		},
	}
	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createDastVulnerabilitiesCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalVeryHigh += subCompareReports.Section.Data[0].Value
		totalHigh += subCompareReports.Section.Data[1].Value
		totalMedium += subCompareReports.Section.Data[2].Value
		totalLow += subCompareReports.Section.Data[3].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalVeryHigh
	compareReports.Section.Data[1].Value = totalHigh
	compareReports.Section.Data[2].Value = totalMedium
	compareReports.Section.Data[3].Value = totalLow

	compareReports.TotalValue = totalVeryHigh + totalHigh + totalMedium + totalLow

	return &compareReports
}

func getContainerVulnerabilitiesComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {
	response, ok := data["containerVulnerabilityScannerChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.VulnerabilitesByScannerType{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getContainerVulnerabilitiesComponentComparison()") {
		return nil, err
	}

	compareReports := createContainerVulnerabilitiesCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getContainerVulnerabilitiesComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createContainerVulnerabilitiesCompareReports(orgData *constants.Organization, compComparison *constants.VulnerabilitesByScannerType) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalVeryHigh, totalHigh, totalMedium, totalLow int

	// Count components
	for _, bucket := range compComparison.Aggregations.VulByScannerTypeComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				var veryHigh, high, medium, low int
				if len(bucket.VulByScannerTypeCounts.Value.VeryHigh) > 0 {
					veryHigh = bucket.VulByScannerTypeCounts.Value.VeryHigh[2].Y
				} else {
					veryHigh = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.High) > 0 {
					high = bucket.VulByScannerTypeCounts.Value.High[2].Y
				} else {
					high = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.Medium) > 0 {
					medium = bucket.VulByScannerTypeCounts.Value.Medium[2].Y
				} else {
					medium = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.Low) > 0 {
					low = bucket.VulByScannerTypeCounts.Value.Low[2].Y
				} else {
					low = 0
				}

				totalVeryHigh += veryHigh
				totalHigh += high
				totalMedium += medium
				totalLow += low

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.VERY_HIGH_TITLE,
						Value: veryHigh,
					}, {
						Title: constants.HIGH_TITLE,
						Value: high,
					}, {
						Title: constants.MEDIUM_TITLE,
						Value: medium,
					},
					{
						Title: constants.LOW_TITLE,
						Value: low,
					},
				}

				value := veryHigh + high + medium + low
				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.VERY_HIGH_TITLE,
			Value: totalVeryHigh,
		}, {
			Title: constants.HIGH_TITLE,
			Value: totalHigh,
		}, {
			Title: constants.MEDIUM_TITLE,
			Value: totalMedium,
		},
		{
			Title: constants.LOW_TITLE,
			Value: totalLow,
		},
	}
	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createContainerVulnerabilitiesCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalVeryHigh += subCompareReports.Section.Data[0].Value
		totalHigh += subCompareReports.Section.Data[1].Value
		totalMedium += subCompareReports.Section.Data[2].Value
		totalLow += subCompareReports.Section.Data[3].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalVeryHigh
	compareReports.Section.Data[1].Value = totalHigh
	compareReports.Section.Data[2].Value = totalMedium
	compareReports.Section.Data[3].Value = totalLow

	compareReports.TotalValue = totalVeryHigh + totalHigh + totalMedium + totalLow

	return &compareReports
}

func getScaVulnerabilitiesComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {
	response, ok := data["scaVulnerabilityScannerChart"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.VulnerabilitesByScannerType{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getScaVulnerabilitiesComponentComparison()") {
		return nil, err
	}

	compareReports := createScaVulnerabilitiesCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getScaVulnerabilitiesComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createScaVulnerabilitiesCompareReports(orgData *constants.Organization, compComparison *constants.VulnerabilitesByScannerType) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalVeryHigh, totalHigh, totalMedium, totalLow int

	// Count components
	for _, bucket := range compComparison.Aggregations.VulByScannerTypeComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				var veryHigh, high, medium, low int
				if len(bucket.VulByScannerTypeCounts.Value.VeryHigh) > 0 {
					veryHigh = bucket.VulByScannerTypeCounts.Value.VeryHigh[3].Y
				} else {
					veryHigh = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.High) > 0 {
					high = bucket.VulByScannerTypeCounts.Value.High[3].Y
				} else {
					high = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.Medium) > 0 {
					medium = bucket.VulByScannerTypeCounts.Value.Medium[3].Y
				} else {
					medium = 0
				}
				if len(bucket.VulByScannerTypeCounts.Value.Low) > 0 {
					low = bucket.VulByScannerTypeCounts.Value.Low[3].Y
				} else {
					low = 0
				}

				totalVeryHigh += veryHigh
				totalHigh += high
				totalMedium += medium
				totalLow += low

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.VERY_HIGH_TITLE,
						Value: veryHigh,
					}, {
						Title: constants.HIGH_TITLE,
						Value: high,
					}, {
						Title: constants.MEDIUM_TITLE,
						Value: medium,
					},
					{
						Title: constants.LOW_TITLE,
						Value: low,
					},
				}

				value := veryHigh + high + medium + low

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = value
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.VERY_HIGH_TITLE,
			Value: totalVeryHigh,
		}, {
			Title: constants.HIGH_TITLE,
			Value: totalHigh,
		}, {
			Title: constants.MEDIUM_TITLE,
			Value: totalMedium,
		},
		{
			Title: constants.LOW_TITLE,
			Value: totalLow,
		},
	}
	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createScaVulnerabilitiesCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		totalVeryHigh += subCompareReports.Section.Data[0].Value
		totalHigh += subCompareReports.Section.Data[1].Value
		totalMedium += subCompareReports.Section.Data[2].Value
		totalLow += subCompareReports.Section.Data[3].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = totalVeryHigh
	compareReports.Section.Data[1].Value = totalHigh
	compareReports.Section.Data[2].Value = totalMedium
	compareReports.Section.Data[3].Value = totalLow

	compareReports.TotalValue = totalVeryHigh + totalHigh + totalMedium + totalLow

	return &compareReports
}

func getDeploymentLeadTimeComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {
	response, ok := data["deploymentLeadTimeHeader"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.DeploymentLeadTimeComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getDeploymentLeadTimeComponentComparison()") {
		return nil, err
	}

	compareReports := createDeploymentLeadTimeCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getDeploymentLeadTimeComponentComparison() chart") {
		return nil, err
	}
	return b, err

}

func createDeploymentLeadTimeCompareReports(orgData *constants.Organization, compComparison *constants.DeploymentLeadTimeComponentComparison) *constants.MttrCompareReports {
	var compareReports constants.MttrCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var deployments, duration int

	// Count components
	for _, bucket := range compComparison.Aggregations.DeploymentLeadTimeComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				deployments += bucket.DeployData.Value.Deployments
				duration += bucket.DeployData.Value.TotalDuration

				compareReports.ComponentCount++
				var compCompareReports constants.MttrCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.AverageValue = bucket.DeployData.Value.Deployments
				compCompareReports.Count = bucket.DeployData.Value.TotalDuration
				compCompareReports.ValueInMillis = calculateAverage(bucket.DeployData.Value.TotalDuration, bucket.DeployData.Value.Deployments)
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createDeploymentLeadTimeCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively

		deployments += subCompareReports.AverageValue
		duration += subCompareReports.Count
	}

	compareReports.ValueInMillis = calculateAverage(duration, deployments)

	return &compareReports
}

func getDoraMttrComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["mttrHeader"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.DoraMttrComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getDoraMttrComponentComparison()") {
		return nil, err
	}

	compareReports := createDoraMttrCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getDoraMttrComponentComparison() chart") {
		return nil, err
	}
	return b, err

}

func createDoraMttrCompareReports(orgData *constants.Organization, compComparison *constants.DoraMttrComponentComparison) *constants.MttrCompareReports {
	var compareReports constants.MttrCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalRecoveredTotalDuration, totalRecoveredCount int

	// Count components
	for _, bucket := range compComparison.Aggregations.MttrComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				totalRecoveredTotalDuration += bucket.Deployments.Value.RecoveredTotalDuration
				totalRecoveredCount += bucket.Deployments.Value.RecoveredCount

				compareReports.ComponentCount++
				var compCompareReports constants.MttrCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.AverageValue = bucket.Deployments.Value.RecoveredTotalDuration
				compCompareReports.Count = bucket.Deployments.Value.RecoveredCount
				compCompareReports.ValueInMillis = calculateAverage(bucket.Deployments.Value.RecoveredTotalDuration, bucket.Deployments.Value.RecoveredCount)
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createDoraMttrCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively

		totalRecoveredTotalDuration += subCompareReports.AverageValue
		totalRecoveredCount += subCompareReports.Count
	}

	compareReports.ValueInMillis = calculateAverage(totalRecoveredTotalDuration, totalRecoveredCount)

	return &compareReports
}

func getFailureRateComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["averageFailureRateHeader"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.FailureRateComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getFailureRateComponentComparison()") {
		return nil, err
	}

	compareReports := createFailureRateCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getFailureRateComponentComparison() chart") {
		return nil, err
	}
	return b, err

}

func createFailureRateCompareReports(orgData *constants.Organization, compComparison *constants.FailureRateComponentComparison) *constants.MttrCompareReports {
	var compareReports constants.MttrCompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var failedDeployments, deployments int

	// Count components
	for _, bucket := range compComparison.Aggregations.FailureRateComponentComparison.Buckets {
		for _, component := range orgData.Components {
			if component.ID == bucket.Key {

				failedDeployments += bucket.DeployData.Value.FailedDeployments
				deployments += bucket.DeployData.Value.Deployments

				compareReports.ComponentCount++
				var compCompareReports constants.MttrCompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.AverageValue = bucket.DeployData.Value.FailedDeployments
				compCompareReports.Count = bucket.DeployData.Value.Deployments
				compCompareReports.TotalValue = calculatePercentageValue(bucket.DeployData.Value.FailedDeployments, bucket.DeployData.Value.Deployments)
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
	}

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createFailureRateCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively

		failedDeployments += subCompareReports.AverageValue
		deployments += subCompareReports.Count
	}

	compareReports.TotalValue = calculatePercentageValue(failedDeployments, deployments)

	return &compareReports
}

func getAutomationRunsForTestSuites(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "totalTestRunsSpec" {
		response, ok := data["automationRuns"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := constants.AutomationRunsCount{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getAutomationRunsForTestSuites()") {
			return nil, err
		}

		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.ComponentActivity.Value.Runs) + `}`)
		return b, nil
	} else if specKey == "testRunsChartSpec" {

		response, ok := data["automationRuns"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		result := constants.AutomationRunsCount{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getAutomationRunsForTestSuites()") {
			return nil, err
		}

		tsResponse, ok := data["testSuites"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		tsResult := constants.TestSuitesCount{}

		err = json.Unmarshal([]byte(tsResponse), &tsResult)
		if log.CheckErrorf(err, "Error unmarshaling responseStruct in getAutomationRunsForTestSuites()") {
			return nil, err
		}

		withTestSuitesPercentage := 0.0
		withoutTestSuitesPercentage := 0.0
		totalRunsCount := result.Aggregations.ComponentActivity.Value.Runs
		withTestSuitesCount := tsResult.Aggregations.ComponentActivity.Value.Runs
		withoutTestSuitesCount := totalRunsCount - withTestSuitesCount

		if totalRunsCount > 0 {
			withTestSuitesPercentage = math.Round((float64(withTestSuitesCount) / float64(totalRunsCount)) * 100)
			withoutTestSuitesPercentage = math.Round((float64(withoutTestSuitesCount) / float64(totalRunsCount)) * 100)

		}

		chart := constants.TestAutomationRunChart{}

		chart.Data = []struct {
			Name  string  `json:"name"`
			Value float64 `json:"value"`
		}{
			{Name: constants.WITH_TEST_SUITES, Value: withTestSuitesPercentage},
			{Name: constants.WITHOUT_TEST_SUITES, Value: withoutTestSuitesPercentage},
		}

		chart.Info = []struct {
			DrillDown struct {
				ReportType  string `json:"reportType"`
				ReportID    string `json:"reportId"`
				ReportTitle string `json:"reportTitle"`
			} `json:"drillDown"`
			Title string `json:"title"`
			Value int    `json:"value"`
		}{
			{
				DrillDown: struct {
					ReportType  string `json:"reportType"`
					ReportID    string `json:"reportId"`
					ReportTitle string `json:"reportTitle"`
				}{
					ReportType:  "testSuiteType",
					ReportID:    "test-insights-workflowRuns",
					ReportTitle: constants.WORKFLOW_RUNS_REPORT_TITLE,
				},
				Title: constants.WITH_TEST_SUITES,
				Value: withTestSuitesCount,
			},
			{
				DrillDown: struct {
					ReportType  string `json:"reportType"`
					ReportID    string `json:"reportId"`
					ReportTitle string `json:"reportTitle"`
				}{
					ReportType:  "testSuiteType",
					ReportID:    "test-insights-workflowRuns",
					ReportTitle: constants.WORKFLOW_RUNS_REPORT_TITLE,
				},
				Title: constants.WITHOUT_TEST_SUITES,
				Value: withoutTestSuitesCount,
			},
		}

		b, err := json.Marshal(chart)
		if log.CheckErrorf(err, "Error marshaling responseStruct in getAutomationRunsForTestSuites()") {
			return nil, err
		}
		return b, nil

	}

	return nil, nil
}

func getTestSuiteWorkflowRunsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	testSuitesResp, ok := data["testSuiteRuns"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	autoRunResp, ok := data["automationRuns"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	testSuiteResult := constants.TestSuitesRun{}
	autoRunResult := constants.TestSuiteAutoRun{}

	testSuiteErr := json.Unmarshal([]byte(testSuitesResp), &testSuiteResult)
	if log.CheckErrorf(testSuiteErr, "Error unmarshaling response from OpenSearch for test suite in getTestSuiteWorkflowRunsComponentComparison()") {
		return nil, testSuiteErr
	}

	autoRunErr := json.Unmarshal([]byte(autoRunResp), &autoRunResult)
	if log.CheckErrorf(autoRunErr, "Error unmarshaling response from OpenSearch for auto runs in getTestSuiteWorkflowRunsComponentComparison()") {
		return nil, autoRunErr
	}

	testSuiteStructMap := make(map[string]constants.TestSuite)

	for _, testBucket := range testSuiteResult.Aggregations.WorkflowRunsComponentComparison.Buckets {
		testSuiteStructMap[testBucket.Key] = constants.TestSuite{
			WithTestSuite: testBucket.ComponentActivity.Value.Runs,
		}
	}

	for _, autoRunBucket := range autoRunResult.Aggregations.WorkflowRunsComponentComparison.Buckets {
		if val, ok := testSuiteStructMap[autoRunBucket.Key]; ok {
			val.WithoutTestSuite = autoRunBucket.RunStatus.Value.Runs - val.WithTestSuite
			testSuiteStructMap[autoRunBucket.Key] = val
		} else {
			testSuiteStructMap[autoRunBucket.Key] = constants.TestSuite{
				WithTestSuite:    0,
				WithoutTestSuite: autoRunBucket.RunStatus.Value.Runs,
			}
		}
	}

	var resultList []map[string]constants.TestSuite
	for key, val := range testSuiteStructMap {
		entry := map[string]constants.TestSuite{
			key: val,
		}
		resultList = append(resultList, entry)
	}
	result := constants.TestSuiteComponentComparison{Val: resultList}
	compareReports := createTestSuiteWorkflowRunsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getTestSuiteWorkflowRunsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createTestSuiteWorkflowRunsCompareReports(orgData *constants.Organization, compComparison *constants.TestSuiteComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var totalWithTestSuites, totalWithoutTestSuites int

	// Count components
	for _, val := range compComparison.Val {
		for _, component := range orgData.Components {
			for mapKey, mapVal := range val {
				if component.ID == mapKey {
					var withTestSuiteCount, withoutTestSuiteCount int
					totalWithTestSuites += mapVal.WithTestSuite
					withTestSuiteCount = mapVal.WithTestSuite
					totalWithoutTestSuites += mapVal.WithoutTestSuite
					withoutTestSuiteCount = mapVal.WithoutTestSuite

					data := []struct {
						Title string `json:"title"`
						Value int    `json:"value"`
					}{
						{
							Title: constants.WITH_TEST_SUITES,
							Value: withTestSuiteCount,
						}, {
							Title: constants.WITHOUT_TEST_SUITES,
							Value: withoutTestSuiteCount,
						},
					}

					value := withTestSuiteCount + withoutTestSuiteCount

					compareReports.ComponentCount++
					var compCompareReports constants.CompareReports

					compCompareReports.SubOrgID = component.ID
					compCompareReports.CompareTitle = component.Name
					compCompareReports.IsSubOrg = false
					compCompareReports.ComponentCount = 1
					compCompareReports.Section.Data = data
					compCompareReports.TotalValue = value
					compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
					break
				}
			}
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.WITH_TEST_SUITES,
			Value: totalWithTestSuites,
		}, {
			Title: constants.WITHOUT_TEST_SUITES,
			Value: totalWithoutTestSuites,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createTestSuiteWorkflowRunsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		totalWithTestSuites += subCompareReports.Section.Data[0].Value
		totalWithoutTestSuites += subCompareReports.Section.Data[1].Value
	}

	compareReports.Section.Data[0].Value = totalWithTestSuites
	compareReports.Section.Data[1].Value = totalWithoutTestSuites

	compareReports.TotalValue = totalWithTestSuites + totalWithoutTestSuites

	return &compareReports
}

func getTestComponentsComponentComparison(specKey string, data map[string]json.RawMessage, replacements map[string]any, organisation *constants.Organization) (json.RawMessage, error) {

	response, ok := data["components"]
	if !ok {
		return nil, db.ErrInternalServer
	}

	result := constants.TestComponentsComponentComparison{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in getTestComponentsComponentComparison()") {
		return nil, err
	}

	compareReports := createTestComponentsCompareReports(organisation, &result)

	b, err := json.Marshal(compareReports.CompareReports)
	if log.CheckErrorf(err, "Error marshaling responseStruct in getTestComponentsComponentComparison() chart") {
		return nil, err
	}
	return b, err
}

func createTestComponentsCompareReports(orgData *constants.Organization, compComparison *constants.TestComponentsComponentComparison) *constants.CompareReports {
	var compareReports constants.CompareReports

	compareReports.SubOrgID = orgData.ID
	compareReports.CompareTitle = orgData.Name
	compareReports.IsSubOrg = true

	var withTestSuites, withoutTestSuites int

	// Count components
	for _, component := range orgData.Components {
		found := false
		for _, componentId := range compComparison.Aggregations.Components.Value {
			if component.ID == componentId {
				found = true
				withTestSuites += 1

				data := []struct {
					Title string `json:"title"`
					Value int    `json:"value"`
				}{
					{
						Title: constants.WITH_TEST_SUITES,
						Value: 1,
					}, {
						Title: constants.WITHOUT_TEST_SUITES,
						Value: 0,
					},
				}

				compareReports.ComponentCount++
				var compCompareReports constants.CompareReports

				compCompareReports.SubOrgID = component.ID
				compCompareReports.CompareTitle = component.Name
				compCompareReports.IsSubOrg = false
				compCompareReports.ComponentCount = 1
				compCompareReports.Section.Data = data
				compCompareReports.TotalValue = 1
				compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
				break
			}
		}
		if !found {
			withoutTestSuites += 1

			data := []struct {
				Title string `json:"title"`
				Value int    `json:"value"`
			}{
				{
					Title: constants.WITH_TEST_SUITES,
					Value: 0,
				}, {
					Title: constants.WITHOUT_TEST_SUITES,
					Value: 1,
				},
			}

			compareReports.ComponentCount++
			var compCompareReports constants.CompareReports

			compCompareReports.SubOrgID = component.ID
			compCompareReports.CompareTitle = component.Name
			compCompareReports.IsSubOrg = false
			compCompareReports.ComponentCount = 1
			compCompareReports.Section.Data = data
			compCompareReports.TotalValue = 1
			compareReports.CompareReports = append(compareReports.CompareReports, compCompareReports)
		}
	}

	data := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{
			Title: constants.WITH_TEST_SUITES,
			Value: withTestSuites,
		}, {
			Title: constants.WITHOUT_TEST_SUITES,
			Value: withoutTestSuites,
		},
	}

	compareReports.Section.Data = data

	// Recursively process sub-orgs
	for _, subOrg := range orgData.SubOrgs {
		subCompareReports := createTestComponentsCompareReports(subOrg, compComparison)
		compareReports.CompareReports = append(compareReports.CompareReports, *subCompareReports)
		compareReports.SubOrgCount++

		// Accumulate total values recursively
		withTestSuites += subCompareReports.Section.Data[0].Value
		withoutTestSuites += subCompareReports.Section.Data[1].Value
	}

	// Update the total value and section data with cumulative values
	compareReports.Section.Data[0].Value = withTestSuites
	compareReports.Section.Data[1].Value = withoutTestSuites

	compareReports.TotalValue = withTestSuites + withoutTestSuites

	return &compareReports
}

// Function to transform data fetched for the Open Findings By Security Tool in the security dashboard
func transformOpenFindingsBySecurityTool(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	if specKey == "openFindingsBySecurityToolSpec" {

		var response json.RawMessage
		var ok bool
		var isApplication bool

		if response, ok = data["openFindingsBySecurityTool"]; !ok {
			if response, ok = data["openFindingsBySecurityToolApplication"]; !ok {
				return nil, db.ErrInternalServer
			}
			isApplication = true
		}

		result := constants.OpenFindingsBySecurityTools{}
		err := json.Unmarshal([]byte(response), &result)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling open_findings_by_security_tool: %v", err)
		}

		// Redirection Info to be used in drill down inside section data
		type RedirectionInfo struct {
			Url string `json:"url"`
		}

		type DrillDown struct {
			ReportID        string            `json:"reportId"`
			RedirectionInfo []RedirectionInfo `json:"redirectionInfo"`
		}

		var output []struct {
			Total              int     `json:"total"`
			FindingsPercentage float64 `json:"findingsPercentage"`
			SecurityToolName   string  `json:"securityToolName"`
			ToolId             string  `json:"toolId"`
			ColorScheme        []struct {
				Color0 string `json:"color0"`
				Color1 string `json:"color1"`
			} `json:"colorScheme"`
			Drilldown DrillDown `json:"drillDown"`
		}

		// total number of findings (sum of all counts for all tools)
		var totalFindings int
		for _, bucket := range result.Aggregations.ToolName.Buckets {
			totalFindings += len(bucket.TrackingID.Buckets)
		}

		for _, bucket := range result.Aggregations.ToolName.Buckets {
			// Calculate the findings percentage for each tool
			findingsPercentage := 0.0
			if totalFindings > 0 {
				findingsPercentage = float64(len(bucket.TrackingID.Buckets)) / float64(totalFindings) * 100
				// rounding up to nearest whole number
				findingsPercentage = math.Round(findingsPercentage)
			}

			// unique tool_id for each security tool
			var toolId string
			if len(bucket.TrackingID.Buckets) > 0 && len(bucket.TrackingID.Buckets[0].ToolID.Buckets) > 0 {
				toolId = bucket.TrackingID.Buckets[0].ToolID.Buckets[0].Key
			}

			// Color Scheme adding from widget definition - to do
			colorScheme := []struct {
				Color0 string `json:"color0"`
				Color1 string `json:"color1"`
			}{
				{
					Color0: "#4696E5",
					Color1: "#0963BD",
				},
			}
			var drilldown DrillDown
			if !isApplication {
				drilldown = DrillDown{
					ReportID: "redirect-url",
					RedirectionInfo: []RedirectionInfo{
						{
							Url: fmt.Sprintf("tools=%s&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL", toolId),
						},
					},
				}
			}

			// Create the output object for each security tool
			output = append(output, struct {
				Total              int     `json:"total"`
				FindingsPercentage float64 `json:"findingsPercentage"`
				SecurityToolName   string  `json:"securityToolName"`
				ToolId             string  `json:"toolId"`
				ColorScheme        []struct {
					Color0 string `json:"color0"`
					Color1 string `json:"color1"`
				} `json:"colorScheme"`
				Drilldown DrillDown `json:"drillDown"`
			}{
				Total:              len(bucket.TrackingID.Buckets),
				FindingsPercentage: findingsPercentage,
				SecurityToolName:   bucket.Key,
				ToolId:             toolId,
				ColorScheme:        colorScheme,
				Drilldown:          drilldown,
			})
		}

		// Sort the output by FindingsPercentage
		sort.Slice(output, func(i, j int) bool {
			return output[i].Total > output[j].Total
		})

		// Use Buffer + json.Encoder to prevent escaping
		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(output) // Encode struct to JSON

		if err != nil {
			return nil, err
		}
		return bytes.TrimRight(buffer.Bytes(), "\n"), nil
	}
	return nil, nil
}

// function to transform data fetched for the Open Findings By Severity in new component security dashboard
func transformOpenFindingsBySeverity(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	if specKey == "openFindingsBySeveritySpec" {
		var response json.RawMessage
		var ok bool

		if response, ok = data["openFindingsBySeverity"]; !ok {
			if response, ok = data["openFindingsBySeverityApplication"]; !ok {
				return nil, db.ErrInternalServer
			}
		}

		type responseStruct []struct {
			Id             string  `json:"id"`
			Value          int     `json:"value"`
			Percentage     float64 `json:"percentage"`
			YAxisFormatter struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			} `json:"yAxisFormatter"`
		}

		result := constants.OpenFindingsBySeverity{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getVelocity()") {
			return nil, err
		}

		// Map to hold defaults
		severities := map[string]*struct {
			Id    string `json:"id"`
			Value int    `json:"value"`
		}{
			strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE): {Id: constants.VERY_HIGH_TITLE, Value: 0},
			strings.ToUpper(constants.HIGH_TITLE):                      {Id: constants.HIGH_TITLE, Value: 0},
			strings.ToUpper(constants.MEDIUM_TITLE):                    {Id: constants.MEDIUM_TITLE, Value: 0},
			strings.ToUpper(constants.LOW_TITLE):                       {Id: constants.LOW_TITLE, Value: 0},
		}

		// Initialize total value for ONLY VERYHIGH, HIGH, MEDIUM, LOW severities
		mainSeverityTotalValue := 0

		// Process input data
		for _, value := range result.Aggregations.OpenFindingsBySeverity.Buckets {
			severityKey := strings.ToUpper(value.Key) // Normalize key to uppercase
			if severity, exists := severities[severityKey]; exists {
				severity.Value = len(value.TrackingID.Buckets) // Update the value if key exists in the map
				mainSeverityTotalValue += severity.Value       // Only add to the main severities total (VERYHIGH, HIGH, MEDIUM, LOW)
			}
		}

		// Define an ordered slice of keys to retain the required order
		orderedKeys := []string{
			strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE),
			strings.ToUpper(constants.HIGH_TITLE),
			strings.ToUpper(constants.MEDIUM_TITLE),
			strings.ToUpper(constants.LOW_TITLE),
		}

		// Final output structure
		var output responseStruct

		// Iterate over the ordered keys and append to output in the correct order
		for _, key := range orderedKeys {
			if severity, exists := severities[key]; exists {
				// Calculate percentage
				percentage := math.Round((float64(severity.Value)/float64(mainSeverityTotalValue))*100*100) / 100

				// Append to the output with the percentage
				output = append(output, struct {
					Id             string  `json:"id"`
					Value          int     `json:"value"`
					Percentage     float64 `json:"percentage"`
					YAxisFormatter struct {
						AppendUnitValue string `json:"appendUnitValue"`
						Type            string `json:"type"`
					} `json:"yAxisFormatter"`
				}{
					Id:         severity.Id,
					Value:      severity.Value,
					Percentage: percentage,
					YAxisFormatter: struct {
						AppendUnitValue string `json:"appendUnitValue"`
						Type            string `json:"type"`
					}{
						AppendUnitValue: "Findings",
						Type:            "APPEND_TEXT",
					},
				})
			}
		}

		// Marshal the output to JSON
		b, err := json.Marshal(output)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, nil
}

func transformOpenFindingsDistributionByCategory(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	if specKey == "openFindingsDistributionByCategorySpec" {
		var response json.RawMessage
		var ok bool
		var isApplication bool

		if response, ok = data["openFindingsDistributionByCategory"]; !ok {
			if response, ok = data["openFindingsDistributionByCategoryApplication"]; !ok {
				return nil, db.ErrInternalServer
			}
			isApplication = true
		}

		result := constants.OpenFindingsDistributionByCategory{}
		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling openFindingsDistributionByCategory response") {
			return nil, err
		}

		colourScheme := []struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			{Color0: "#FA6D71", Color1: "#D5252A"},
			{Color0: "#FCC16C", Color1: "#FF8307"},
			{Color0: "#FCFF89", Color1: "#FDC913"},
			{Color0: "#9FB6C1", Color1: "#648192"},
		}

		// Define an ordered slice of keys for severity
		orderedKeys := []string{
			strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE),
			strings.ToUpper(constants.HIGH_TITLE),
			strings.ToUpper(constants.MEDIUM_TITLE),
			strings.ToUpper(constants.LOW_TITLE),
		}
		categories := []constants.OpenFindingsByCategoryData{}
		for _, bucket := range result.Aggregations.Category.Buckets {
			category := constants.OpenFindingsByCategoryData{}
			category.CategoryName = capitalizeFirst(bucket.Key)

			// Map for default severity values
			severities := map[string]*constants.SeverityData{
				strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE): {Title: constants.VERY_HIGH_TITLE, Value: 0},
				strings.ToUpper(constants.HIGH_TITLE):                      {Title: constants.HIGH_TITLE, Value: 0},
				strings.ToUpper(constants.MEDIUM_TITLE):                    {Title: constants.MEDIUM_TITLE, Value: 0},
				strings.ToUpper(constants.LOW_TITLE):                       {Title: constants.LOW_TITLE, Value: 0},
			}

			// Create a reverse lookup map to normalize inputs
			normalizedSeverities := map[string]string{
				constants.VERY_HIGH_TITLE: strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE), // "very high" -> "VERY_HIGH"
				constants.HIGH_TITLE:      strings.ToUpper(constants.HIGH_TITLE),                      // "high" -> "HIGH"
				constants.MEDIUM_TITLE:    strings.ToUpper(constants.MEDIUM_TITLE),                    // "medium" -> "MEDIUM"
				constants.LOW_TITLE:       strings.ToUpper(constants.LOW_TITLE),                       // "low" -> "LOW"
			}
			for _, key := range orderedKeys {
				severities[key] = &constants.SeverityData{Title: capitalizeFirst(key), Value: 0}
			}
			for _, val := range bucket.Severity.Buckets {
				severityKey := strings.ToUpper(val.Key) // Normalize key
				if severity, exists := severities[severityKey]; exists {
					severity.Title = capitalizeFirst(val.Key)
					severity.Value = len(val.TrackingID.Buckets)
					category.Total += len(val.TrackingID.Buckets)
				}

			}
			redirectionInfos := []constants.RedirectionInfo{}
			for _, key := range orderedKeys {
				if severity, exists := severities[key]; exists {
					if severity.Title == constants.VERY_HIGH_WITH_UNDERSCORE_TITLE {
						severity.Title = capitalizeFirst(constants.VERY_HIGH_TITLE)
					}
					category.SeverityDistribution.Data = append(category.SeverityDistribution.Data, *severity)
					category.SeverityDistribution.ColorScheme = colourScheme
					// Only add drilldown if not application
					if !isApplication {
						if severityLevel, exists := normalizedSeverities[severity.Title]; exists {
							redirectionInfo := constants.RedirectionInfo{
								Id: severity.Title,
								Url: fmt.Sprintf("categories=%s&severities=%s&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL",
									normalizeToConstantFormat(category.CategoryName), severityLevel),
							}
							redirectionInfos = append(redirectionInfos, redirectionInfo)
						}
					}
				}

			}

			// Severity DrillDown only if not application
			if !isApplication {
				category.SeverityDistribution.DrillDown = constants.DrillDown{
					ReportID:        "redirect-url",
					RedirectionInfo: redirectionInfos,
				}
			}

			// Top-level drilldown (category-level) only if not application
			if !isApplication {
				totalRedirectionInfo := constants.RedirectionInfo{
					Url: fmt.Sprintf("categories=%s&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL", normalizeToConstantFormat(category.CategoryName)),
				}
				category.DrillDown = constants.DrillDown{
					ReportID:        "redirect-url",
					RedirectionInfo: []constants.RedirectionInfo{totalRedirectionInfo},
				}
			}
			categories = append(categories, category)

		}

		// Sort categories by total count
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].Total > categories[j].Total
		})

		// Use Buffer + json.Encoder to prevent escaping
		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(categories) // Encode struct to JSON

		if err != nil {
			return nil, err
		}
		return bytes.TrimRight(buffer.Bytes(), "\n"), nil

	}
	return nil, nil
}

func normalizeToConstantFormat(input string) string {
	// Convert to uppercase
	upperStr := strings.ToUpper(input)

	// Replace hyphens and spaces with underscores
	normalizedStr := strings.ReplaceAll(upperStr, "-", "_")
	normalizedStr = strings.ReplaceAll(normalizedStr, " ", "_")

	return normalizedStr
}

// Capitalize first letter and replace underscore with space for UI to display
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	s = strings.ReplaceAll(s, "_", " ")
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func transformSlasBreachedByAsset(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	if specKey == "slaBreachedByAssetSpec" {

		var response json.RawMessage
		var ok bool
		var isApplication bool

		if response, ok = data["slaBreachedByAssetType"]; !ok {
			if response, ok = data["slaBreachedByAssetTypeApplication"]; !ok {
				return nil, db.ErrInternalServer
			}
			isApplication = true
		}

		type slaColorScheme struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}

		// Redirection Info to be used in drill down inside section data
		type RedirectionInfo struct {
			Id  string `json:"id"`
			Url string `json:"url"`
		}

		type DrillDown struct {
			ReportID        string            `json:"reportId"`
			RedirectionInfo []RedirectionInfo `json:"redirectionInfo"`
		}

		type responseStruct []struct {
			AssetType          string           `json:"assetType"`
			FindingsPercentage float64          `json:"findingsPercentage"`
			Total              int              `json:"total"`
			ColorScheme        []slaColorScheme `json:"colorScheme"`
			DrillDown          DrillDown        `json:"drillDown"`
		}

		result := constants.SlaBreachedByAssetType{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response SlaBreachedByAssetType()") {
			return nil, err
		}

		var totalFindings int
		for _, bucket := range result.Aggregations.RemediationKey.Buckets {
			totalFindings += len(bucket.TrackingID.Buckets)
		}

		var output responseStruct

		for _, bucket := range result.Aggregations.RemediationKey.Buckets {

			drilldown := DrillDown{}
			if !isApplication {
				drilldown = DrillDown{
					ReportID: "redirect-url",
					RedirectionInfo: []RedirectionInfo{{
						Url: fmt.Sprintf("assetTypes=%s&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL&sla=BREACHED", strings.ToUpper(bucket.Key)),
					}},
				}
			}

			output = append(output, struct {
				AssetType          string           `json:"assetType"`
				FindingsPercentage float64          `json:"findingsPercentage"`
				Total              int              `json:"total"`
				ColorScheme        []slaColorScheme `json:"colorScheme"`
				DrillDown          DrillDown        `json:"drillDown"`
			}{
				AssetType:          strings.ToUpper(bucket.Key),
				FindingsPercentage: math.Round((float64(len(bucket.TrackingID.Buckets)) / float64(totalFindings)) * 100),
				Total:              len(bucket.TrackingID.Buckets),
				ColorScheme:        []slaColorScheme{{Color0: "#4696E5", Color1: "#0963BD"}},
				DrillDown:          drilldown,
			})
		}

		// Sort the output by FindingsPercentage
		sort.Slice(output, func(i, j int) bool {
			return output[i].FindingsPercentage > output[j].FindingsPercentage
		})

		// Use Buffer + json.Encoder to prevent escaping
		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(output) // Encode struct to JSON

		if err != nil {
			return nil, err
		}
		return bytes.TrimRight(buffer.Bytes(), "\n"), nil
	}

	return nil, nil
}

// function to transform data fetched for the Open Findings Distribution By Severity in new component security dashboard
func transformOpenFindingsDistributionBySecurityTool(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	if specKey == "openFindingsDistributionBySecurityToolSpec" {
		var response json.RawMessage
		var ok bool
		var isApplication bool

		if response, ok = data["openFindingsDistributionBySecurityTool"]; !ok {
			if response, ok = data["openFindingsDistributionBySecurityToolApplication"]; !ok {
				return nil, db.ErrInternalServer
			}
			isApplication = true
		}

		result := constants.OpenFindingsDistributionBySecurityTool{}
		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling openFindingsDistributionBySecTool response") {
			return nil, err
		}

		// output structs
		type SeverityValue struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		}

		type SeverityColourScheme struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}

		// Redirection Info to be used in drill down inside section data
		type RedirectionInfo struct {
			Id  string `json:"id"`
			Url string `json:"url"`
		}

		type DrillDown struct {
			ReportID        string            `json:"reportId"`
			RedirectionInfo []RedirectionInfo `json:"redirectionInfo"`
		}

		type SeverityDistribution struct {
			ColorScheme []SeverityColourScheme `json:"colorScheme"`
			Data        []SeverityValue        `json:"data"`
			DrillDown   DrillDown              `json:"drillDown"`
		}

		type responseStruct []struct {
			TotalCount           int                  `json:"total"`
			SeverityDistribution SeverityDistribution `json:"severityDistribution"`
			SecToolName          string               `json:"securityToolName"`
			ToolId               string               `json:"toolId"`
			DrillDown            DrillDown            `json:"drillDown"`
		}

		orderedKeys := []string{
			strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE),
			strings.ToUpper(constants.HIGH_TITLE),
			strings.ToUpper(constants.MEDIUM_TITLE),
			strings.ToUpper(constants.LOW_TITLE),
		}

		keyConstants := map[string]string{
			strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE): constants.VERY_HIGH_TITLE,
			strings.ToUpper(constants.HIGH_TITLE):                      constants.HIGH_TITLE,
			strings.ToUpper(constants.MEDIUM_TITLE):                    constants.MEDIUM_TITLE,
			strings.ToUpper(constants.LOW_TITLE):                       constants.LOW_TITLE,
		}

		// Map for default severity values
		severitySec := map[string]SeverityValue{
			strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE): {Title: constants.VERY_HIGH_TITLE, Value: 0},
			strings.ToUpper(constants.HIGH_TITLE):                      {Title: constants.HIGH_TITLE, Value: 0},
			strings.ToUpper(constants.MEDIUM_TITLE):                    {Title: constants.MEDIUM_TITLE, Value: 0},
			strings.ToUpper(constants.LOW_TITLE):                       {Title: constants.LOW_TITLE, Value: 0},
		}

		colourSchemes := map[string]*SeverityColourScheme{
			strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE): {Color0: "#FA6D71", Color1: "#D5252A"},
			strings.ToUpper(constants.HIGH_TITLE):                      {Color0: "#FCC16C", Color1: "#FF8307"},
			strings.ToUpper(constants.MEDIUM_TITLE):                    {Color0: "#FCFF89", Color1: "#FDC913"},
			strings.ToUpper(constants.LOW_TITLE):                       {Color0: "#9FB6C1", Color1: "#648192"},
		}

		findingsBySecTool := make(map[string]struct {
			Total       int
			Severity    map[string]SeverityValue
			ColorScheme map[string]*SeverityColourScheme
			SecToolName string
			SecToolId   string
		})

		for _, bucket := range result.Aggregations.ToolID.Buckets {
			toolIdKey := bucket.Key
			severities := make(map[string]SeverityValue)
			for _, severityRec := range bucket.Severity.Buckets {
				severityKey := strings.ToUpper(severityRec.Key)
				if severity, exists := severitySec[severityKey]; exists {
					severity.Value = severity.Value + len(severityRec.TrackingID.Buckets)
					var toolName string
					if len(severityRec.TrackingID.Buckets) > 0 && len(severityRec.TrackingID.Buckets[0].ToolDisplayName.Buckets) > 0 {
						toolName = severityRec.TrackingID.Buckets[0].ToolDisplayName.Buckets[0].Key
					}

					severities[severityKey] = severity
					findingsBySecTool[toolIdKey] = struct {
						Total       int
						Severity    map[string]SeverityValue
						ColorScheme map[string]*SeverityColourScheme
						SecToolName string
						SecToolId   string
					}{
						Total:       0,
						Severity:    severities,
						ColorScheme: colourSchemes,
						SecToolName: toolName,
						SecToolId:   toolIdKey,
					}
				}
			}

		}

		var output responseStruct
		for _, finding := range findingsBySecTool {
			var TotalCount int
			severityDistribution := SeverityDistribution{}
			redirectionInfos := []RedirectionInfo{}
			for _, key := range orderedKeys {
				if severity, exists := finding.Severity[key]; exists {
					severityDistribution.Data = append(severityDistribution.Data, severity)
					TotalCount += severity.Value
				} else {
					severityDistribution.Data = append(severityDistribution.Data, SeverityValue{Title: keyConstants[key], Value: 0})
				}
				if colourScheme, exists := finding.ColorScheme[key]; exists {
					severityDistribution.ColorScheme = append(severityDistribution.ColorScheme, *colourScheme)
				}

				if !isApplication {
					redirectionInfos = append(redirectionInfos, RedirectionInfo{
						Id:  keyConstants[key],
						Url: fmt.Sprintf("tools=%s&severities=%s&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL", finding.SecToolId, key),
					})
				}
			}

			if !isApplication {
				severityDistribution.DrillDown = DrillDown{
					ReportID:        "redirect-url",
					RedirectionInfo: redirectionInfos,
				}
			}

			var totalDrilldown DrillDown
			if !isApplication {
				totalDrilldown = DrillDown{
					ReportID: "redirect-url",
					RedirectionInfo: []RedirectionInfo{
						{
							Url: fmt.Sprintf("tools=%s&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL", finding.SecToolId),
						},
					},
				}
			}

			output = append(output, struct {
				TotalCount           int                  `json:"total"`
				SeverityDistribution SeverityDistribution `json:"severityDistribution"`
				SecToolName          string               `json:"securityToolName"`
				ToolId               string               `json:"toolId"`
				DrillDown            DrillDown            `json:"drillDown"`
			}{
				TotalCount:           TotalCount,
				SeverityDistribution: severityDistribution,
				SecToolName:          finding.SecToolName,
				ToolId:               finding.SecToolId,
				DrillDown:            totalDrilldown,
			})
		}

		// Sort the output by FindingsPercentage
		sort.Slice(output, func(i, j int) bool {
			return output[i].TotalCount > output[j].TotalCount
		})

		// Use Buffer + json.Encoder to prevent escaping
		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(output) // Encode struct to JSON

		if err != nil {
			return nil, err
		}
		return bytes.TrimRight(buffer.Bytes(), "\n"), nil
	}
	return nil, nil
}

// function to transform data fetched for the findings identified widget in new component security dashboard
func transformFindingsIdentfiedSince(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	var response json.RawMessage
	var ok bool

	if response, ok = data["findingsIdentifiedCount"]; !ok {
		if response, ok = data["findingsIdentifiedCountApplication"]; !ok {
			return nil, db.ErrInternalServer
		}
	}

	result := constants.FindingsIdentifiedSince{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallingTransformFindingsIdentifiedSince) {
		return nil, err
	}

	// Values to display in the response
	open := findValue(result, constants.REMEDIATION_STATUS_OPEN)
	resolved := findValue(result, constants.REMEDIATION_STATUS_RESOLVED)
	riskAccepted := findValue(result, constants.REMEDIATION_STATUS_RISK_ACCEPTED)
	inProgress := findValue(result, constants.REMEDIATION_STATUS_IN_PROGRESS)
	falsePositive := findValue(result, constants.REMEDIATION_STATUS_FALSE_POSITIVE)

	// Calculate percentages
	total := open + resolved + riskAccepted + inProgress + falsePositive

	var openPercentage, resolvedPercentage, riskAcceptedPercentage, inProgressPercentage, falsePositivePercentage float64
	if total > 0 {
		openPercentage = math.Round(float64(open)/float64(total)*100*100) / 100
		resolvedPercentage = math.Round(float64(resolved)/float64(total)*100*100) / 100
		riskAcceptedPercentage = math.Round(float64(riskAccepted)/float64(total)*100*100) / 100
		inProgressPercentage = math.Round(float64(inProgress)/float64(total)*100*100) / 100
		falsePositivePercentage = math.Round(float64(falsePositive)/float64(total)*100*100) / 100
	}
	if specKey == "openFindingsSpec" {
		b := []byte(`{"value":` + fmt.Sprint(findValue(result, constants.REMEDIATION_STATUS_OPEN)) + `}`)
		return b, nil
	} else if specKey == "resolvedFindingsSpec" {
		b := []byte(`{"value":` + fmt.Sprint(findValue(result, constants.REMEDIATION_STATUS_RESOLVED)) + `}`)
		return b, nil
	} else if specKey == "riskAcceptedFindingsSpec" {
		b := []byte(`{"value":` + fmt.Sprint(findValue(result, constants.REMEDIATION_STATUS_RISK_ACCEPTED)) + `}`)
		return b, nil
	} else if specKey == "inProgressFindingsSpec" {
		b := []byte(`{"value":` + fmt.Sprint(findValue(result, constants.REMEDIATION_STATUS_IN_PROGRESS)) + `}`)
		return b, nil
	} else if specKey == "falsePositiveFindingsSpec" {
		b := []byte(`{"value":` + fmt.Sprint(findValue(result, constants.REMEDIATION_STATUS_FALSE_POSITIVE)) + `}`)
		return b, nil
	} else if specKey == "findingsIndentifiedChartSpec" {
		type YAxisFormatter struct {
			AppendUnitValue string `json:"appendUnitValue"`
			Type            string `json:"type"`
		}

		type responseStruct struct {
			Id             string         `json:"id"`
			Value          int            `json:"value"`
			Percentage     float64        `json:"percentage"`
			YAxisFormatter YAxisFormatter `json:"yAxisFormatter"`
		}

		// response structure for chart data
		outputStruct := []responseStruct{
			{
				Id:             constants.REMEDIATION_STATUS_OPEN_LABEL,
				Value:          open,
				Percentage:     openPercentage,
				YAxisFormatter: YAxisFormatter{"Findings", "APPEND_TEXT"},
			},
			{
				Id:             constants.REMEDIATION_STATUS_RESOLVED_LABEL,
				Value:          resolved,
				Percentage:     resolvedPercentage,
				YAxisFormatter: YAxisFormatter{"Findings", "APPEND_TEXT"},
			},
			{
				Id:             constants.REMEDIATION_STATUS_RISK_ACCEPTED_LABEL,
				Value:          riskAccepted,
				Percentage:     riskAcceptedPercentage,
				YAxisFormatter: YAxisFormatter{"Findings", "APPEND_TEXT"},
			},
			{
				Id:             constants.REMEDIATION_STATUS_IN_PROGRESS_LABEL,
				Value:          inProgress,
				Percentage:     inProgressPercentage,
				YAxisFormatter: YAxisFormatter{"Findings", "APPEND_TEXT"},
			},
			{
				Id:             constants.REMEDIATION_STATUS_FALSE_POSITIVE_LABEL,
				Value:          falsePositive,
				Percentage:     falsePositivePercentage,
				YAxisFormatter: YAxisFormatter{"Findings", "APPEND_TEXT"},
			},
		}

		b, err := json.Marshal(outputStruct)
		if log.CheckErrorf(err, "Error marshaling responseStruct in transformFindingsIdentifiedSince() chart") {
			return nil, err
		}
		return b, nil
	}

	return nil, nil
}

// findOpenValue iterates over the buckets and returns the value of Key, or 0 if not found.
func findValue(findings constants.FindingsIdentifiedSince, key string) int {
	for _, bucket := range findings.Aggregations.RemediationStatus.Buckets {
		if bucket.Key == key {
			return len(bucket.TrackingID.Buckets)
		}
	}
	return 0
}

// findRAFPValue iterates over the buckets and returns the value of Key, or 0 if not found.
func findRAFPValue(findings constants.RiskAcceptedFalsePositiveFindings, key string) int {
	for _, bucket := range findings.Aggregations.RemediationStatus.Buckets {
		if bucket.Key == key {
			return len(bucket.TrackingID.Buckets)
		}
	}
	return 0
}

// function to transform data fetched for the RA and FP widget in new component security dashboard
func transformRiskAcceptedAndFalsePositiveFindings(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	var response json.RawMessage
	var ok bool

	if response, ok = data["riskAcceptedFalsePositiveFindingsCount"]; !ok {
		if response, ok = data["riskAcceptedFalsePositiveFindingsCountApplication"]; !ok {
			return nil, db.ErrInternalServer
		}
	}

	result := constants.RiskAcceptedFalsePositiveFindings{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallingTransformRiskAcceptedAndFalsePositiveFindings) {
		return nil, err
	}

	if specKey == "totalFindingsSpec" {
		b := []byte(`{"value":` + fmt.Sprint(findRAFPValue(result, constants.REMEDIATION_STATUS_OPEN)+findRAFPValue(result, constants.REMEDIATION_STATUS_IN_PROGRESS)+findRAFPValue(result, constants.REMEDIATION_STATUS_RISK_ACCEPTED)+findRAFPValue(result, constants.REMEDIATION_STATUS_FALSE_POSITIVE)) + `}`)
		return b, nil
	} else if specKey == "riskAcceptedFindingsSpec" {
		b := []byte(`{"value":` + fmt.Sprint(findRAFPValue(result, constants.REMEDIATION_STATUS_RISK_ACCEPTED)) + `}`)
		return b, nil
	} else if specKey == "falsePositiveFindingsSpec" {
		b := []byte(`{"value":` + fmt.Sprint(findRAFPValue(result, constants.REMEDIATION_STATUS_FALSE_POSITIVE)) + `}`)
		return b, nil
	} else if specKey == "raExpiringIn30DaysFindingsSpec" {
		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.RiskAcceptedCounts.RA_EXPIRING_IN_30_DAYS.DocCount) + `}`)
		return b, nil
	} else if specKey == "raExpiredFindingsSpec" {
		b := []byte(`{"value":` + fmt.Sprint(result.Aggregations.RiskAcceptedCounts.RA_EXPIRED.DocCount) + `}`)
		return b, nil
	} else if specKey == "riskAcceptedChartSpec" || specKey == "falsePositiveChartSpec" {
		type responseStruct struct {
			ColorScheme []struct {
				Color0 string `json:"color0"`
				Color1 string `json:"color1"`
			} `json:"colorScheme"`
			FindingsPercentage float64 `json:"findingsPercentage"`
		}

		var output []responseStruct
		var percentage float64
		if specKey == "riskAcceptedChartSpec" {
			percentage = (float64(findRAFPValue(result, constants.REMEDIATION_STATUS_RISK_ACCEPTED)) /
				(float64(findRAFPValue(result, constants.REMEDIATION_STATUS_OPEN)) +
					float64(findRAFPValue(result, constants.REMEDIATION_STATUS_IN_PROGRESS)) +
					float64(findRAFPValue(result, constants.REMEDIATION_STATUS_FALSE_POSITIVE)) +
					float64(findRAFPValue(result, constants.REMEDIATION_STATUS_RISK_ACCEPTED))) * 100)
		} else if specKey == "falsePositiveChartSpec" {
			percentage = (float64(findRAFPValue(result, constants.REMEDIATION_STATUS_FALSE_POSITIVE)) /
				(float64(findRAFPValue(result, constants.REMEDIATION_STATUS_OPEN)) +
					float64(findRAFPValue(result, constants.REMEDIATION_STATUS_IN_PROGRESS)) +
					float64(findRAFPValue(result, constants.REMEDIATION_STATUS_FALSE_POSITIVE)) +
					float64(findRAFPValue(result, constants.REMEDIATION_STATUS_RISK_ACCEPTED))) * 100)
		}

		percentage = math.Round(percentage*100) / 100

		colorScheme := []struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		}{
			{
				Color0: "#4696E5",
				Color1: "#0963BD",
			},
		}
		output = append(output, responseStruct{
			ColorScheme:        colorScheme,
			FindingsPercentage: percentage,
		})

		b, err := json.Marshal(output)
		if log.CheckErrorf(err, "Error marshaling responseStruct in transformRiskAcceptedAndFalsePositiveFindings() chart") {
			return nil, err
		}
		return b, nil
	}

	return nil, nil
}

type SeverityData struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type RedirectionInfoWithUrlOnly struct {
	Url string `json:"url"`
}

type DrillDownWithRedirectionInfo struct {
	ReportID        string                       `json:"reportId"`
	RedirectionInfo []RedirectionInfoWithUrlOnly `json:"redirectionInfo"`
}

type SeveritySubHeader struct {
	DrillDown DrillDownWithRedirectionInfo `json:"drillDown"`
	Title     string                       `json:"title"`
	Value     int                          `json:"value"`
	Color     string                       `json:"color"`
}

func calculatePercentagesforSeverityData(severities []SeverityData, total float64) []SeverityData {
	for i := range severities {
		percentage := (severities[i].Value / total) * 100
		severities[i].Value = math.Round(percentage*100) / 100 // Round to two decimal places
	}
	return severities
}

func transformSlaBreachesBySeverity(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	var response json.RawMessage
	var ok bool
	var isApplication bool

	if response, ok = data["slaBreachesBySeverity"]; !ok {
		if response, ok = data["slaBreachesBySeverityApplication"]; !ok {
			return nil, db.ErrInternalServer
		}
		isApplication = true
	}

	result := constants.SLABreachesBySeverity{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallingTransformSlaBreachesBySeverity) {
		return nil, err
	}

	// Define the order of keys
	orderedKeys := []string{
		constants.VERY_HIGH_TITLE,
		constants.HIGH_TITLE,
		constants.MEDIUM_TITLE,
		constants.LOW_TITLE,
	}

	// Map to hold defaults
	severities := make(map[string]*SeverityData)

	for _, bucket := range result.Aggregations.SLABreachesBySeverity.Buckets {
		if bucket.Key == strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE) {
			veryHigh := bucket.TrackingID.Buckets
			veryHighCount := len(veryHigh)
			severities[constants.VERY_HIGH_TITLE] = &SeverityData{Name: constants.VERY_HIGH_TITLE, Value: float64(veryHighCount)}
		} else if bucket.Key == strings.ToUpper(constants.HIGH_TITLE) {
			high := bucket.TrackingID.Buckets
			highCount := len(high)
			severities[constants.HIGH_TITLE] = &SeverityData{Name: constants.HIGH_TITLE, Value: float64(highCount)}
		} else if bucket.Key == strings.ToUpper(constants.MEDIUM_TITLE) {
			medium := bucket.TrackingID.Buckets
			mediumCount := len(medium)
			severities[constants.MEDIUM_TITLE] = &SeverityData{Name: constants.MEDIUM_TITLE, Value: float64(mediumCount)}
		} else if bucket.Key == strings.ToUpper(constants.LOW_TITLE) {
			low := bucket.TrackingID.Buckets
			lowCount := len(low)
			severities[constants.LOW_TITLE] = &SeverityData{Name: constants.LOW_TITLE, Value: float64(lowCount)}
		}
	}

	// Calculate total by checking if the key exists
	veryHighValue := 0.0
	if val, exists := severities[constants.VERY_HIGH_TITLE]; exists {
		veryHighValue = val.Value
	}

	highValue := 0.0
	if val, exists := severities[constants.HIGH_TITLE]; exists {
		highValue = val.Value
	}

	mediumValue := 0.0
	if val, exists := severities[constants.MEDIUM_TITLE]; exists {
		mediumValue = val.Value
	}

	lowValue := 0.0
	if val, exists := severities[constants.LOW_TITLE]; exists {
		lowValue = val.Value
	}

	total := int(veryHighValue + highValue + mediumValue + lowValue)

	switch specKey {
	case "slaBreachesBySeverityHeaderSpec":

		b := []byte(`{"value":` + fmt.Sprint(total) + `}`)
		return b, nil

	case "slaBreachesBySeveritySubHeaderSpec":
		remidiationUrl := "&sla=BREACHED&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"
		SubHeaderData := []SeveritySubHeader{}

		for _, key := range orderedKeys {
			value, exists := severities[key]
			if !exists {
				continue
			}

			var drillDown DrillDownWithRedirectionInfo
			var hasDrillDown bool

			var color string

			// Drilldown within subheader disabled for application security temporarily
			if !isApplication {
				switch key {
				case constants.VERY_HIGH_TITLE:
					drillDown = DrillDownWithRedirectionInfo{
						ReportID: "redirect-url",
						RedirectionInfo: []RedirectionInfoWithUrlOnly{
							{Url: fmt.Sprintf("severities=%s%s", strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE), remidiationUrl)},
						},
					}
					color = "#EA4F54"
					hasDrillDown = true
				case constants.HIGH_TITLE:
					drillDown = DrillDownWithRedirectionInfo{
						ReportID: "redirect-url",
						RedirectionInfo: []RedirectionInfoWithUrlOnly{
							{Url: fmt.Sprintf("severities=%s%s", strings.ToUpper(constants.HIGH_TITLE), remidiationUrl)},
						},
					}
					color = "#FE9D33"
					hasDrillDown = true
				case constants.MEDIUM_TITLE:
					drillDown = DrillDownWithRedirectionInfo{
						ReportID: "redirect-url",
						RedirectionInfo: []RedirectionInfoWithUrlOnly{
							{Url: fmt.Sprintf("severities=%s%s", strings.ToUpper(constants.MEDIUM_TITLE), remidiationUrl)},
						},
					}
					color = "#FCE44E"
					hasDrillDown = true
				case constants.LOW_TITLE:
					drillDown = DrillDownWithRedirectionInfo{
						ReportID: "redirect-url",
						RedirectionInfo: []RedirectionInfoWithUrlOnly{
							{Url: fmt.Sprintf("severities=%s%s", strings.ToUpper(constants.LOW_TITLE), remidiationUrl)},
						},
					}
					color = "#738E9D"
					hasDrillDown = true
				}
			} else {
				// Still assign color even if drilldown is nil for application security case
				switch key {
				case constants.VERY_HIGH_TITLE:
					color = "#EA4F54"
				case constants.HIGH_TITLE:
					color = "#FE9D33"
				case constants.MEDIUM_TITLE:
					color = "#FCE44E"
				case constants.LOW_TITLE:
					color = "#738E9D"
				}
			}
			SubHeaderData = append(SubHeaderData, SeveritySubHeader{
				DrillDown: func() DrillDownWithRedirectionInfo {
					if hasDrillDown {
						return drillDown
					}
					return DrillDownWithRedirectionInfo{}
				}(),
				Title: key,
				Value: int(value.Value),
				Color: color,
			})
		}

		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		encoder.SetEscapeHTML(false) // Handling ambersand in URL
		err = encoder.Encode(SubHeaderData)
		if err != nil {
			return nil, err
		}

		b := []byte(`{"subHeader":` + buffer.String() + `}`)
		return b, nil

	case "slaBreachesBySeveritySpec":
		sevList := []SeverityData{}
		for _, key := range orderedKeys {
			value, exists := severities[key]
			if !exists {
				missingValue := SeverityData{Name: key, Value: 0}
				sevList = append(sevList, missingValue)
				continue
			}
			sevList = append(sevList, *value)

		}

		// Calculate percentages
		newList := calculatePercentagesforSeverityData(sevList, float64(total))

		b, err := json.Marshal(newList)
		if log.CheckErrorf(err, "Error marshaling responseStruct in transformSlaBreachesBySeverity()") {
			return nil, err
		}

		return b, nil
	}
	return nil, nil
}

func transformOpenFindingsBySlaStatus(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	var response json.RawMessage
	var ok bool

	if response, ok = data["openFindingsBySLAStatus"]; !ok {
		if response, ok = data["openFindingsBySLAStatusApplication"]; !ok {
			return nil, db.ErrInternalServer
		}
	}

	result := constants.OpenFindingsBySLAStatus{}
	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, "Error unmarshaling openFindingsBySlaStatus response") {
		return nil, err
	}

	WithinSlaCount := len(result.Aggregations.OpenFindingsBySLAStatus.Buckets.NonSLABreached.TrackingID.Buckets)
	BreachedSlaCount := len(result.Aggregations.OpenFindingsBySLAStatus.Buckets.SLABreached.TrackingID.Buckets)
	TotalCount := WithinSlaCount + BreachedSlaCount

	switch specKey {
	case "openFindingsSpec":
		b := []byte(`{"value":` + fmt.Sprint(TotalCount) + `}`)
		return b, nil

	case "withinSLASpec":
		b := []byte(`{"value":` + fmt.Sprint(WithinSlaCount) + `}`)
		return b, nil

	case "breachedSLASpec":
		b := []byte(`{"value":` + fmt.Sprint(BreachedSlaCount) + `}`)
		return b, nil

	case "openFindingsBySLAStatusChartSpec":
		responseStruct := []SeverityData{
			{Name: constants.WITHIN_SLA_TITLE, Value: float64(WithinSlaCount)},
			{Name: constants.BREACHED_SLA_TITLE, Value: float64(BreachedSlaCount)},
		}

		// Calculate percentages
		percentageList := calculatePercentagesforSeverityData(responseStruct, float64(TotalCount))

		b, err := json.Marshal(percentageList)
		if log.CheckErrorf(err, "Error marshaling responseStruct in transformOpenFindingsBySlaStatus() chart") {
			return nil, err
		}
		return b, nil

	}

	return nil, nil

}

func transformOpenFindingsByReviewStatus(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {

	var response json.RawMessage
	var ok bool

	if response, ok = data["openFindingsByReviewStatus"]; !ok {
		if response, ok = data["openFindingsByReviewStatusApplication"]; !ok {
			return nil, db.ErrInternalServer
		}
	}

	result := constants.OpenFindingsByReviewStatus{}

	err := json.Unmarshal([]byte(response), &result)
	if log.CheckErrorf(err, exceptions.ErrUnmarshallingTransformOpenFindingsByReviewStatus) {
		return nil, err
	}

	orderedKeys := []string{
		constants.TRIAGE_STATUS_UNREVIEWED_LABEL,
		constants.TRIAGE_STATUS_FIX_REQUIRED_LABEL,
		constants.TRIAGE_STATUS_AWAITING_APPROVAL_LABEL,
	}

	severities := make(map[string]*SeverityData)
	for _, bucket := range result.Aggregations.TriageStatus.Buckets {

		if bucket.Key == constants.TRIAGE_STATUS_IN_REVIEW {
			inReviewTriage := bucket.TrackingID.Buckets
			inReviewTriageCount := len(inReviewTriage)
			severities[constants.TRIAGE_STATUS_AWAITING_APPROVAL_LABEL] = &SeverityData{Name: constants.TRIAGE_STATUS_AWAITING_APPROVAL_LABEL, Value: float64(inReviewTriageCount)}
		} else if bucket.Key == constants.TRIAGE_STATUS_FIX_REQUIRED {
			fixRequiredTriage := bucket.TrackingID.Buckets
			fixRequiredTriageCount := len(fixRequiredTriage)
			severities[constants.TRIAGE_STATUS_FIX_REQUIRED_LABEL] = &SeverityData{Name: constants.TRIAGE_STATUS_FIX_REQUIRED_LABEL, Value: float64(fixRequiredTriageCount)}
		} else if bucket.Key == constants.TRIAGE_STATUS_UNREVIEWED {
			unreviewedTriage := bucket.TrackingID.Buckets
			unreviewedTriageCount := len(unreviewedTriage)
			severities[constants.TRIAGE_STATUS_UNREVIEWED_LABEL] = &SeverityData{Name: constants.TRIAGE_STATUS_UNREVIEWED_LABEL, Value: float64(unreviewedTriageCount)}
		}
	}
	total := 0.0
	if severity, ok := severities[constants.TRIAGE_STATUS_AWAITING_APPROVAL_LABEL]; ok && severity != nil {
		total += severity.Value
	}

	if severity, ok := severities[constants.TRIAGE_STATUS_FIX_REQUIRED_LABEL]; ok && severity != nil {
		total += severity.Value
	}

	if severity, ok := severities[constants.TRIAGE_STATUS_UNREVIEWED_LABEL]; ok && severity != nil {
		total += severity.Value
	}

	switch specKey {
	case "unreviewedSpec":
		if severity, ok := severities[constants.TRIAGE_STATUS_UNREVIEWED_LABEL]; ok && severity != nil {
			b := []byte(`{"value":` + fmt.Sprint(severity.Value) + `}`)
			return b, nil
		}
		b := []byte(`{"value":0}`)
		return b, nil
	case "awaitingApprovalSpec":
		if severity, ok := severities[constants.TRIAGE_STATUS_AWAITING_APPROVAL_LABEL]; ok && severity != nil {
			b := []byte(`{"value":` + fmt.Sprint(severity.Value) + `}`)
			return b, nil
		}
		b := []byte(`{"value":0}`)
		return b, nil
	case "fixRequiredSpec":
		if severity, ok := severities[constants.TRIAGE_STATUS_FIX_REQUIRED_LABEL]; ok && severity != nil {
			b := []byte(`{"value":` + fmt.Sprint(severity.Value) + `}`)
			return b, nil
		}
		b := []byte(`{"value":0}`)
		return b, nil
	case "openFindingsByReviewStatusSectionSpec":
		sevList := []SeverityData{}
		for _, key := range orderedKeys {
			value, exists := severities[key]
			if !exists {
				continue
			}
			sevList = append(sevList, *value)

		}

		newList := calculatePercentagesforSeverityData(sevList, float64(total))

		b, err := json.Marshal(newList)
		if log.CheckErrorf(err, "Error marshaling responseStruct in transformOpenFindingsByReviewStatus()") {
			return nil, err
		}

		return b, nil
	}
	return nil, nil
}

func transformFindingsRemediationTrend(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	if specKey == "findingsRemediationTrendSpec" {
		response, ok := data["findingsRemediationTrend"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type responseStruct struct {
			Id   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
			YAxisFormatter struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			} `json:"yAxisFormatter"`
		}

		result := constants.FindingsRemediationTrend{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getVelocity()") {
			return nil, err
		}

		// Form Open, Closed within SLA and Breached SLA response structs
		openData := responseStruct{
			Id: "Open",
			YAxisFormatter: struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			}{
				AppendUnitValue: "Findings",
				Type:            "APPEND_TEXT",
			},
		}

		closedWithinSLAData := responseStruct{
			Id: "Closed within SLA",
			YAxisFormatter: struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			}{
				AppendUnitValue: "Findings",
				Type:            "APPEND_TEXT",
			},
		}

		breachedSlaData := responseStruct{
			Id: "Breached SLA",
			YAxisFormatter: struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			}{
				AppendUnitValue: "Findings",
				Type:            "APPEND_TEXT",
			},
		}

		if len(result.Aggregations.FindingsRemediationTrend.Value) == 0 {
			return nil, errors.New("No data found for Findings Remediation Trend for duration type " + replacements["duration"].(string))
		}

		// sort the date keys in the OpenSearch response map in ascending order
		keys := make([]string, 0, len(result.Aggregations.FindingsRemediationTrend.Value))
		for k := range result.Aggregations.FindingsRemediationTrend.Value {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, key := range keys {
			remediationStruct := result.Aggregations.FindingsRemediationTrend.Value[key]
			date := key

			// Convert date to the right format based on the duration type
			dateString := ""
			if replacements["duration"].(string) == constants.DURATION_WEEK {
				// Convert date to day of week
				t, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME, date)
				if err != nil {
					return nil, err
				}
				dateString = t.Weekday().String()

			} else if replacements["duration"].(string) == constants.DURATION_MONTH {
				// Convert date to DD MMM format
				t, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME, date)
				if err != nil {
					return nil, err
				}
				dateString = t.Format("02 Jan")
			} else if replacements["duration"].(string) == constants.DURATION_YEAR {
				// for the last element in the dates map, set the date to "Latest"
				if i == len(keys)-1 {
					dateString = "Latest"
				} else {
					// Convert date to MMM YY format
					t, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME, date)
					if err != nil {
						return nil, err
					}
					dateString = t.Format("Jan 06")
				}
			} else {
				return nil, errors.New("Invalid duration type")
			}

			openData.Data = append(openData.Data, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{
				X: dateString,
				Y: remediationStruct.Open,
			})

			closedWithinSLAData.Data = append(closedWithinSLAData.Data, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{
				X: dateString,
				Y: remediationStruct.ClosedWithinSLA,
			})

			breachedSlaData.Data = append(breachedSlaData.Data, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{
				X: dateString,
				Y: remediationStruct.BreachedSLA,
			})
		}

		output := []responseStruct{openData, closedWithinSLAData, breachedSlaData}

		// Marshal the output to JSON
		b, err := json.Marshal(output)
		if err != nil {
			return nil, err
		}
		return b, nil
	} else if specKey == "findingsRemediationTrendAppSec" {
		response, ok := data["findingsRemediationTrend"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		type responseStruct struct {
			Id   string `json:"id"`
			Data []struct {
				X string `json:"x"`
				Y int    `json:"y"`
			} `json:"data"`
			YAxisFormatter struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			} `json:"yAxisFormatter"`
		}

		result := constants.FindingsRemediationTrendAppSec{}

		err := json.Unmarshal([]byte(response), &result)
		if log.CheckErrorf(err, "Error unmarshaling response getVelocity()") {
			return nil, err
		}

		// Form Open, Closed within SLA, Breached SLA and New response structs
		openData := responseStruct{
			Id: "Open",
			YAxisFormatter: struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			}{
				AppendUnitValue: "Findings",
				Type:            "APPEND_TEXT",
			},
		}

		closedWithinSLAData := responseStruct{
			Id: "Closed within SLA",
			YAxisFormatter: struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			}{
				AppendUnitValue: "Findings",
				Type:            "APPEND_TEXT",
			},
		}

		breachedSlaData := responseStruct{
			Id: "Breached SLA",
			YAxisFormatter: struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			}{
				AppendUnitValue: "Findings",
				Type:            "APPEND_TEXT",
			},
		}

		newData := responseStruct{
			Id: "New",
			YAxisFormatter: struct {
				AppendUnitValue string `json:"appendUnitValue"`
				Type            string `json:"type"`
			}{
				AppendUnitValue: "Findings",
				Type:            "APPEND_TEXT",
			},
		}

		if len(result.Aggregations.FindingsRemediationTrend.Value) == 0 {
			return nil, errors.New("No data found for Findings Remediation Trend for duration type " + replacements["duration"].(string))
		}

		// sort the date keys in the OpenSearch response map in ascending order
		keys := make([]string, 0, len(result.Aggregations.FindingsRemediationTrend.Value))
		for k := range result.Aggregations.FindingsRemediationTrend.Value {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, key := range keys {
			remediationStruct := result.Aggregations.FindingsRemediationTrend.Value[key]
			date := key

			// Convert date to the right format based on the duration type
			dateString := ""
			if replacements["duration"].(string) == constants.DURATION_WEEK {
				// Convert date to day of week
				t, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME, date)
				if err != nil {
					return nil, err
				}
				dateString = t.Weekday().String()

			} else if replacements["duration"].(string) == constants.DURATION_MONTH {
				// Convert date to DD MMM format
				t, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME, date)
				if err != nil {
					return nil, err
				}
				dateString = t.Format("02 Jan")
			} else if replacements["duration"].(string) == constants.DURATION_YEAR {
				// for the last element in the dates map, set the date to "Latest"
				if i == len(keys)-1 {
					dateString = "Latest"
				} else {
					// Convert date to MMM YY format
					t, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME, date)
					if err != nil {
						return nil, err
					}
					dateString = t.Format("Jan 06")
				}
			} else {
				return nil, errors.New("Invalid duration type")
			}

			openData.Data = append(openData.Data, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{
				X: dateString,
				Y: remediationStruct.Open,
			})

			closedWithinSLAData.Data = append(closedWithinSLAData.Data, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{
				X: dateString,
				Y: remediationStruct.ClosedWithinSLA,
			})

			breachedSlaData.Data = append(breachedSlaData.Data, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{
				X: dateString,
				Y: remediationStruct.BreachedSLA,
			})

			newData.Data = append(newData.Data, struct {
				X string `json:"x"`
				Y int    `json:"y"`
			}{
				X: dateString,
				Y: remediationStruct.New,
			})
		}

		output := []responseStruct{openData, closedWithinSLAData, newData, breachedSlaData}

		// Marshal the output to JSON
		b, err := json.Marshal(output)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	return nil, nil
}

func transformOpenFindingsByComponent(specKey string, data map[string]json.RawMessage, replacements map[string]any) (json.RawMessage, error) {
	if specKey == "openFindingsByComponentSpec" {

		responseRAW, ok := data["openFindingsByComponent"]
		if !ok {
			return nil, db.ErrInternalServer
		}

		// Defining the structure for Open Findings By Component query response
		type findingsByComponentResponse struct {
			Hits struct {
				Total struct {
					Value int `json:"value"`
				} `json:"total"`
			} `json:"hits"`
			Aggregations struct {
				Components struct {
					DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
					SumOtherDocCount        int `json:"sum_other_doc_count"`
					Buckets                 []struct {
						Key                     string `json:"key"`
						DocCount                int    `json:"doc_count"`
						FindingsByEachComponent struct {
							Value int `json:"value"`
						} `json:"findings_by_each_component"`
						Categories struct {
							DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
							SumOtherDocCount        int `json:"sum_other_doc_count"`
							Buckets                 []struct {
								Key                    string `json:"key"`
								DocCount               int    `json:"doc_count"`
								FindingsByEachCategory struct {
									DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
									SumOtherDocCount        int `json:"sum_other_doc_count"`
									Buckets                 []struct {
										Key      string `json:"key"`
										DocCount int    `json:"doc_count"`
									} `json:"buckets"`
								} `json:"findings_by_each_category"`
							} `json:"buckets"`
						} `json:"categories"`
						Severities struct {
							DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
							SumOtherDocCount        int `json:"sum_other_doc_count"`
							Buckets                 []struct {
								Key                    string `json:"key"`
								DocCount               int    `json:"doc_count"`
								FindingsByEachSeverity struct {
									DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
									SumOtherDocCount        int `json:"sum_other_doc_count"`
									Buckets                 []struct {
										Key      string `json:"key"`
										DocCount int    `json:"doc_count"`
									} `json:"buckets"`
								} `json:"findings_by_each_severity"`
							} `json:"buckets"`
						} `json:"severities"`
					} `json:"buckets"`
				} `json:"components"`
			} `json:"aggregations"`
		}

		response := findingsByComponentResponse{}

		err := json.Unmarshal([]byte(responseRAW), &response)
		if log.CheckErrorf(err, "Error unmarshaling response in transformOpenFindingsByComponent") {
			return nil, err
		}
		if response.Hits.Total.Value == 0 {
			return nil, errors.New("No data found for Open Findings By Component")
		}

		// Defining the order of severity levels
		severityOrder := []string{
			constants.VERY_HIGH_TITLE,
			constants.HIGH_TITLE,
			constants.MEDIUM_TITLE,
			constants.LOW_TITLE,
		}

		// Map to hold which severities have non-zero findings
		severitiesWithNonZeroFindings := map[string]bool{
			constants.VERY_HIGH_TITLE: false,
			constants.HIGH_TITLE:      false,
			constants.MEDIUM_TITLE:    false,
			constants.LOW_TITLE:       false,
		}

		// Map to transform the severity strings in the data store to corresponding display labels
		severityLabelsMap := map[string]string{
			strings.ToUpper(constants.VERY_HIGH_WITH_UNDERSCORE_TITLE): constants.VERY_HIGH_TITLE,
			strings.ToUpper(constants.HIGH_TITLE):                      constants.HIGH_TITLE,
			strings.ToUpper(constants.MEDIUM_TITLE):                    constants.MEDIUM_TITLE,
			strings.ToUpper(constants.LOW_TITLE):                       constants.LOW_TITLE,
		}

		// Map to hold color codes for each severity
		severityColorsMap := map[string]string{
			constants.VERY_HIGH_TITLE: "#EA4F54",
			constants.HIGH_TITLE:      "#FE9D33",
			constants.MEDIUM_TITLE:    "#FCE44E",
			constants.LOW_TITLE:       "#738E9D",
		}

		// Defining the structure for final transformed response

		type tooltipDataWithColor struct {
			Name               string  `json:"name"`
			Findings           int     `json:"findings"`
			FindingsPercentage float32 `json:"findingsPercentage"`
			Color              string  `json:"color"`
		}

		type tooltipDistributionWithColor struct {
			Title string                 `json:"title"`
			Data  []tooltipDataWithColor `json:"data"`
		}

		type findingsByComponentTransformedResponse struct {
			Id          string                         `json:"id"`
			Value       int                            `json:"value"`
			TooltipInfo []tooltipDistributionWithColor `json:"tooltipInfo"`
		}

		transformedResponse := []findingsByComponentTransformedResponse{}

		// Sorting the Aggregations.Components.Buckets slice in the response struct in descending order based on the value in FindingsByEachComponent.Value
		sort.Slice(response.Aggregations.Components.Buckets, func(i, j int) bool {
			return response.Aggregations.Components.Buckets[i].FindingsByEachComponent.Value > response.Aggregations.Components.Buckets[j].FindingsByEachComponent.Value
		})

		// Getingt the top 10 components
		topComponents := response.Aggregations.Components.Buckets
		if len(topComponents) > 10 {
			topComponents = topComponents[:10]
		}

		// Transforming the response to the required format

		for _, component := range topComponents {
			totalFindings := float64(component.FindingsByEachComponent.Value)
			if totalFindings == 0 {
				continue // Skip if there are no findings for this component
			}

			transformedComponent := findingsByComponentTransformedResponse{
				Id:          component.Key,
				Value:       component.FindingsByEachComponent.Value,
				TooltipInfo: make([]tooltipDistributionWithColor, 2), // Two distributions: Severity and Category
			}

			// The first index will hold severity distribution and the second index will hold category distribution
			transformedComponent.TooltipInfo[0] = tooltipDistributionWithColor{
				Title: "Severity distribution",
				Data:  []tooltipDataWithColor{},
			}
			transformedComponent.TooltipInfo[1] = tooltipDistributionWithColor{
				Title: "Category distribution",
				Data:  []tooltipDataWithColor{},
			}

			// Calculate severity distribution
			for _, severityBucket := range component.Severities.Buckets {
				findingsForCurrentSeverity := len(severityBucket.FindingsByEachSeverity.Buckets)
				percentage := float64(findingsForCurrentSeverity*100) / totalFindings

				// Use the SeverityLabelsMap to get the corresponding label for the severity
				severityLabel, ok := severityLabelsMap[strings.ToUpper(severityBucket.Key)]
				if !ok {
					severityLabel = severityBucket.Key
				}

				severitiesWithNonZeroFindings[severityLabel] = true // Mark that this severity has findings

				severityData := tooltipDataWithColor{
					Name:               severityLabel,
					Findings:           findingsForCurrentSeverity,
					FindingsPercentage: float32(math.Round(percentage*100) / 100), // Round to two decimal places
					Color:              severityColorsMap[severityLabel],
				}

				transformedComponent.TooltipInfo[0].Data = append(transformedComponent.TooltipInfo[0].Data, severityData)
			}

			// every severity in the severityOrder should be present in the final response, if not, it should be added with 0 findings
			for _, severity := range severityOrder {
				if !severitiesWithNonZeroFindings[severity] {
					severityData := tooltipDataWithColor{
						Name:               severity,
						Findings:           0,
						FindingsPercentage: 0.00,
						Color:              severityColorsMap[severity],
					}
					transformedComponent.TooltipInfo[0].Data = append(transformedComponent.TooltipInfo[0].Data, severityData)
				}
			}

			// sort the severity distribution data based on the order in severityOrder.
			sort.Slice(transformedComponent.TooltipInfo[0].Data, func(i, j int) bool {
				// Get the severity names from the data
				severityNameI := transformedComponent.TooltipInfo[0].Data[i].Name
				severityNameJ := transformedComponent.TooltipInfo[0].Data[j].Name
				// Find the index of the severity names in the orderedKeys
				indexI := slices.Index(severityOrder, severityNameI)
				indexJ := slices.Index(severityOrder, severityNameJ)
				// Compare the indices to determine the order
				return indexI < indexJ
			})

			// Calculate category distribution
			for _, categoryBucket := range component.Categories.Buckets {
				findingsForCurrentCategory := len(categoryBucket.FindingsByEachCategory.Buckets)
				percentage := float64(findingsForCurrentCategory*100) / totalFindings

				categoryData := tooltipDataWithColor{
					Name:               capitalizeFirst(categoryBucket.Key),
					Findings:           findingsForCurrentCategory,
					FindingsPercentage: float32(math.Round(percentage*100) / 100), // Round to two decimal places
				}

				transformedComponent.TooltipInfo[1].Data = append(transformedComponent.TooltipInfo[1].Data, categoryData)
			}

			transformedResponse = append(transformedResponse, transformedComponent)

		}
		// Marshal the transformed response to JSON
		b, err := json.Marshal(transformedResponse)
		if log.CheckErrorf(err, "Error marshaling responseStruct in transformOpenFindingsByComponent()") {
			return nil, err
		}
		return b, nil

	}
	return nil, nil

}
