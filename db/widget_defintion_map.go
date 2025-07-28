package db

import "github.com/calculi-corp/reports-service/constants"

// Map of widgetId to Widget definition json file name
//
//	TBD : This will be moved to DB
var WidgetDefinitionMap = map[string]string{
	"e1":     "swdelivery/components.json",
	"e2":     "swdelivery/automations.json",
	"e3":     "swdelivery/automationsRuns.json",
	"e4":     "swdelivery/commitTrends.json",
	"e5":     "swdelivery/pullRequests.json",
	"e6":     "swdelivery/commitsCodeChurn.json",
	"e7":     "swdelivery/codeProgressionSnapshot.json",
	"e8":     "swdelivery/successfulBuildDurationMainWidget.json",
	"e8-sub": "swdelivery/successfulBuildDuration.json",
	"e9":     "swdelivery/deploymentSuccessRateMainWidget.json",
	"e9-sub": "swdelivery/deploymentSuccessRate.json",
	"e10":    "swdelivery/developmentCycleTime.json",
	"e11":    "swdelivery/averageDeploymentTime.json",
	"e12":    "swdelivery/codeChurn.json",

	"s1":  "security/components.json",
	"s2":  "security/automations.json",
	"s3":  "security/automationsRuns.json",
	"s4":  "security/openFixed.json",
	"s5":  "security/openVulnerabilitiesOverview.json",
	"s6":  "security/scanTypesInAutomation.json",
	"s7":  "security/vulbyscannertype.json",
	"s8":  "security/SLAStatusOverview.json",
	"s9":  "security/MTTR.json",
	"s10": "security/topVulnerabilities.json",
	"s33": "vulbyscannertype.json",
	"s44": "scansInAutomations.json",

	"f1": "flow/components.json",
	"f2": "flow/pullRequests.json",
	"f3": "flow/tickets.json",
	"f4": "flow/workload.json",
	"f5": "flow/workitemDistribution.json",
	"f6": "flow/velocity.json",
	"f7": "flow/cycleTime.json",
	"f8": "flow/workEfficiency.json",
	"d1": "dorametrics/deploymentFrequency.json",
	"d2": "dorametrics/deploymentLeadTime.json",
	"d3": "dorametrics/failureRate.json",
	"d4": "dorametrics/mttr.json",
	"d5": "dorametrics/depFrequencyAndLeadTimeTrend.json",
	"d6": "dorametrics/failureRateAndMttrTrend.json",

	"cs1":     "componentSummary/components-activity.json",
	"cs2":     "componentSummary/components-builds.json",
	"cs3":     "componentSummary/components-deployments.json",
	"cs3-sub": "componentSummary/deploymentSuccessRate.json",

	"cs4":  "componentSummary/components-code-coverage.json",
	"cs5":  "componentSummary/components-issue-types.json",
	"cs6":  "componentSummary/components-duplications.json",
	"cs7":  "componentSummary/components-codebase-overview.json",
	"cs8":  "componentSummary/components-open-issues.json",
	"cs9":  "componentSummary/components-no-scanners-configured.json",
	"cs10": "componentSummary/components-trivy-licenses-overview.json",
	"cs11": "componentSummary/components-latest-test-results.json",

	"ci1": "ciInsight/projectTypes.json",
	"ci2": "ciInsight/systemInformation.json",
	"ci4": "ciInsight/runsOverview.json",
	"ci3": "ciInsight/systemHealth.json",
	"ci5": "ciInsight/completedRuns.json",
	"ci6": "ciInsight/usagePatterns.json",
	"ci7": "ciInsight/projectsActivity.json",

	"ci01": "ciInsight/cjocControllers.json",

	"ti1": "testInsight/components.json",
	"ti2": "testInsight/automations.json",
	"ti3": "testInsight/automationsRuns.json",
	"ti4": "testInsight/testsOverview.json",

	"so1":  "componentSecurity/findingsRemediationTrend.json",
	"so2":  "componentSecurity/openFindingsBySeverity.json",
	"so3":  "componentSecurity/riskAcceptedFalsePositiveFindings.json",
	"so4":  "componentSecurity/SLABreachesBySeverity.json",
	"so5":  "componentSecurity/findingsIdentifiedSince.json",
	"so6":  "componentSecurity/openFindingsByReviewStatus.json",
	"so7":  "componentSecurity/openFindingsBySLAStatus.json",
	"so8":  "componentSecurity/openFindingsBySecurityTool.json",
	"so9":  "componentSecurity/openFindingsDistributionBySecurityTool.json",
	"so10": "componentSecurity/openFindingsDistributionByCategory.json",
	"so11": "componentSecurity/SLABreachedbyAssets.json",

	"aso1":  "applicationSecurity/findingsRemediationTrend.json",
	"aso2":  "applicationSecurity/openFindingsBySeverity.json",
	"aso3":  "applicationSecurity/openFindingsByComponent.json",
	"aso4":  "applicationSecurity/riskAcceptedFalsePositiveFindings.json",
	"aso5":  "applicationSecurity/SLABreachesBySeverity.json",
	"aso6":  "applicationSecurity/findingsIdentifiedSince.json",
	"aso7":  "applicationSecurity/openFindingsByReviewStatus.json",
	"aso8":  "applicationSecurity/openFindingsBySLAStatus.json",
	"aso9":  "applicationSecurity/openFindingsBySecurityTool.json",
	"aso10": "applicationSecurity/openFindingsDistributionBySecurityTool.json",
	"aso11": "applicationSecurity/openFindingsDistributionByCategory.json",
	"aso12": "applicationSecurity/SLABreachedbyAssets.json",

	"velocity-compare":                 "componentComparison/flow/velocityComponentComparison.json",
	"cycle-time-compare":               "componentComparison/flow/cycleTimeComponentComparison.json",
	"average-active-work-time-compare": "componentComparison/flow/activeWorkTimeComponentComparison.json",
	"work-wait-time-compare":           "componentComparison/flow/workWaitTimeComponentComparison.json",
	"work-load-compare":                "componentComparison/flow/workLoadComponentComparison.json",

	"components-compare":                             "componentComparison/swdelivery/componentsComponentComparison.json",
	"commit-trends-compare":                          "componentComparison/swdelivery/commitTrendsComponentComparison.json",
	"pull-requests-trend-compare":                    "componentComparison/swdelivery/pullRequestsComponentComparison.json",
	"deployment-success-rate-compare":                "componentComparison/swdelivery/deploymentOverviewComponentComparison.json",
	"development-cycle-compare":                      "componentComparison/swdelivery/developmentCycleTimeComponentComparison.json",
	"workflow-runs-compare":                          "componentComparison/swdelivery/workflowRunsComponentComparison.json",
	"run-initiating-commits-compare":                 "componentComparison/swdelivery/commitsComponentComparison.json",
	"builds-compare":                                 "componentComparison/swdelivery/buildsComponentComparison.json",
	"successful-deployments-compare":                 "componentComparison/swdelivery/deploymentsComponentComparison.json",
	"workflows-compare":                              "componentComparison/swdelivery/workflowsComponentComparison.json",
	"security-workflows-compare":                     "componentComparison/security/workflowsCC.json",
	"vulnerabilities-overview-compare":               "componentComparison/security/vulnerabilitiesOverviewComponentComparison.json",
	"open-vulnerabilities-compare":                   "componentComparison/security/openVulnerabilitiesOverviewComponentComparison.json",
	"security-workflow-runs-compare":                 "componentComparison/security/workflowRunsComponentComparison.json",
	"security-components-compare":                    "componentComparison/security/componentsComponentComparison.json",
	"mttr-vulnerabilities-very-high-compare":         "componentComparison/security/mttrVeryHighComponentComparison.json",
	"mttr-vulnerabilities-high-compare":              "componentComparison/security/mttrHighComponentComparison.json",
	"mttr-vulnerabilities-medium-compare":            "componentComparison/security/mttrMediumComponentComparison.json",
	"mttr-vulnerabilities-low-compare":               "componentComparison/security/mttrLowComponentComparison.json",
	"vulnerabilities-scanner-type-SAST-compare":      "componentComparison/security/sastVulnerabilitiesComponentComparison.json",
	"vulnerabilities-scanner-type-DAST-compare":      "componentComparison/security/dastVulnerabilitiesComponentComparison.json",
	"vulnerabilities-scanner-type-container-compare": "componentComparison/security/containerVulnerabilitiesComponentComparison.json",
	"vulnerabilities-scanner-type-SCA-compare":       "componentComparison/security/scaVulnerabilitiesComponentComparison.json",
	"deployment-frequency-compare":                   "componentComparison/dorametrics/deploymentFrequencyComponentComparison.json",
	"mttr-compare":                                   "componentComparison/dorametrics/mttrComponentComparison.json",
	"deployment-lead-time-compare":                   "componentComparison/dorametrics/deploymentLeadTimeComponentComparison.json",
	"failure-rate-compare":                           "componentComparison/dorametrics/failureRateComponentComparison.json",
	"test-insights-workflows-compare":                "componentComparison/testInsight/workflowsComponentComparison.json",
	"test-insights-workflow-runs-compare":            "componentComparison/testInsight/workflowRunsComponentComparison.json",
	"test-insights-components-compare":               "componentComparison/testInsight/componentsComponentComparison.json",

	"99":   "funcSample.json",
	"test": "automations.json",
}

var PaginationReportMap = map[string]string{
	"e9":  "deploymentSuccessRateMainWidget.json",
	"cs3": "deploymentSuccessRateMainWidget.json",

	"e8": "successfulBuildDurationMainWidget.json",
}

var PaginationSubReportMap = map[string]string{
	"e9":   "e9-sub",
	"e8":   "e8-sub",
	"cs3":  "cs3-sub",
	"es9":  "es9-sub",
	"es8":  "es8-sub",
	"css3": "css3-sub",
}

var PaginationReportFilterMap = map[string]string{
	"e8":   "component_id",
	"e9":   "target_env",
	"cs3":  "target_env",
	"es8":  "component_id",
	"es9":  "target_env",
	"css3": "target_env",
}

var PaginationReportFilterCountMap = map[string]int{
	"e8":  5,
	"e9":  1,
	"cs3": 1,
}

var ReportDefinitionMap = map[int64]string{
	1: "scanReport.json",
}

var DrillDownQueryDefinitionMap = map[string]string{
	constants.OPEN_VULNERABILITIES:                             constants.OpenVulnerabilitiesDrillDownQuery,
	constants.NESTED_DRILLDOWN_VIEW_LOCATION:                   constants.ViewLocationsNestedDrillDownQuery,
	constants.OPEN_VULNERABILITIES_VIEW_LOCATION:               constants.ViewLocationsNestedDrillDownQuery,
	constants.VULNERABILITIES_SECURITY_SCAN_TYPE:               constants.VulnerabiltyByScannerTypeDrillDownQuery,
	constants.VULNERABILITIES_SECURITY_SCAN_TYPE_VIEW_LOCATION: constants.ViewLocationsNestedDrillDownQuery,
	constants.VULNERABILITIES_OVERVIEW:                         constants.VulnerabilitiesOverviewDrillDownQuery,
	constants.VULNERABILITIES_OVERVIEW_VIEW_LOCATION:           constants.ViewLocationsNestedDrillDownQuery,
	constants.MTTR_FOR_VULNERABILITIES:                         constants.MTTRDrillDownQuery,
	constants.CWE_TOP_25_VULNERABILITIES:                       constants.Top25OpenVulnerabilitiesDrillDownQuery,
	constants.CWE_TOP_25_VULNERABILITIES_VIEW_LOCATION:         constants.ViewLocationsNestedDrillDownQuery,
	constants.SECURITY_SLA_STATUS_OVERVIEW_OPEN:                constants.SLAStatusOverviewOpenDrillDownQuery,
	constants.SECURITY_SLA_STATUS_OVERVIEW_CLOSED:              constants.SLAStatusOverviewClosedDrillDownQuery,
	constants.FLOW_METRICS_VELOCITY:                            constants.VelocityDrilldownQuery,
	constants.FLOW_METRICS_WORK_ITEM_DISTRIBUTION:              constants.DistributionDrilldownQuery,
	constants.FLOW_METRICS_CYCLE_TIME:                          constants.CycleTimeDrilldownQuery,
	constants.FLOW_METRICS_WORK_EFFICIENCY:                     constants.WorkEfficiencyDrilldownQuery,
	constants.FLOW_METRICS_WORK_LOAD:                           constants.WorkLoadDrilldownQuery,
	constants.ACTIVE_DEVELOPERS:                                constants.ActiveDeveloperDrillDownQuery,
	constants.ACTIVE_DEVELOPERS_COMMITS:                        constants.NestedActiveDeveloperCommitsInfoQuery,
	constants.TRIVY_LICENSE_OCCURENCE:                          constants.TrivyLicenseDrilldownQuery,
	constants.OPEN_ISSUES_DRILL_DOWN:                           constants.OpenIssuesDrillDownQuery,
	constants.LATEST_TEST_RESULTS:                              constants.SummaryLatestTestResultsDrilldownQuery,
	constants.OPEN_VULNERABILITIES_SUBROWS:                     constants.OpenVulnerabilitiesSubRowsQuery,
	constants.VULNERABILITIES_OVERVIEW_SUBROWS:                 constants.VulnerabilitiesOverviewSubRowsQuery,
	constants.MTTR_FOR_VULNERABILITIES_SUBROWS:                 constants.MTTRDrillDownSubRowsQuery,
	constants.CWE_TOP_25_VULNERABILITIES_SUBROWS:               constants.Top25OpenVulnerabilitiesSubRowsQuery,
	constants.VULNERABILITIES_SECURITY_SCAN_TYPE_SUBROWS:       constants.VulnerabiltyByScannerTypeDrillDownSubRowsQuery,
}

var DrillDownAliasDefinitionMap = map[string]string{
	constants.OPEN_VULNERABILITIES:                             constants.SECURITY_INDEX,
	constants.NESTED_DRILLDOWN_VIEW_LOCATION:                   constants.SECURITY_INDEX,
	constants.OPEN_VULNERABILITIES_VIEW_LOCATION:               constants.SECURITY_INDEX,
	constants.VULNERABILITIES_SECURITY_SCAN_TYPE:               constants.SECURITY_INDEX,
	constants.VULNERABILITIES_SECURITY_SCAN_TYPE_VIEW_LOCATION: constants.SECURITY_INDEX,
	constants.VULNERABILITIES_OVERVIEW:                         constants.SECURITY_INDEX,
	constants.VULNERABILITIES_OVERVIEW_VIEW_LOCATION:           constants.SECURITY_INDEX,
	constants.MTTR_FOR_VULNERABILITIES:                         constants.SECURITY_INDEX,
	constants.CWE_TOP_25_VULNERABILITIES:                       constants.SECURITY_INDEX,
	constants.CWE_TOP_25_VULNERABILITIES_VIEW_LOCATION:         constants.SECURITY_INDEX,
	constants.SECURITY_SLA_STATUS_OVERVIEW_OPEN:                constants.SECURITY_INDEX,
	constants.SECURITY_SLA_STATUS_OVERVIEW_CLOSED:              constants.SECURITY_INDEX,
	constants.FLOW_METRICS_VELOCITY:                            constants.FLOW_METRICS_INDEX,
	constants.FLOW_METRICS_WORK_ITEM_DISTRIBUTION:              constants.FLOW_METRICS_INDEX,
	constants.FLOW_METRICS_CYCLE_TIME:                          constants.FLOW_METRICS_INDEX,
	constants.FLOW_METRICS_WORK_EFFICIENCY:                     constants.FLOW_METRICS_INDEX,
	constants.FLOW_METRICS_WORK_LOAD:                           constants.FLOW_METRICS_INDEX,
	constants.ACTIVE_DEVELOPERS:                                constants.COMMIT_DATA_INDEX,
	constants.ACTIVE_DEVELOPERS_COMMITS:                        constants.COMMIT_DATA_INDEX,
	constants.TRIVY_LICENSE_OCCURENCE:                          constants.SECURITY_INDEX,
	constants.OPEN_ISSUES_DRILL_DOWN:                           constants.SECURITY_INDEX,
	constants.LATEST_TEST_RESULTS:                              constants.TEST_CASES_INDEX,
	constants.OPEN_VULNERABILITIES_SUBROWS:                     constants.SECURITY_INDEX,
	constants.VULNERABILITIES_OVERVIEW_SUBROWS:                 constants.SECURITY_INDEX,
	constants.MTTR_FOR_VULNERABILITIES_SUBROWS:                 constants.SECURITY_INDEX,
	constants.CWE_TOP_25_VULNERABILITIES_SUBROWS:               constants.SECURITY_INDEX,
	constants.VULNERABILITIES_SECURITY_SCAN_TYPE_SUBROWS:       constants.SECURITY_INDEX,
}
var DrillDownFilterDefinitionMap = map[string]string{
	constants.OPEN_VULNERABILITIES:                             constants.OPEN_VUL_FILTER,
	constants.NESTED_DRILLDOWN_VIEW_LOCATION:                   constants.OPEN_VUL_FILTER,
	constants.OPEN_VULNERABILITIES_VIEW_LOCATION:               constants.OPEN_VUL_FILTER,
	constants.VULNERABILITIES_SECURITY_SCAN_TYPE:               constants.VUL_BY_SCANTYPE_COLUMNS,
	constants.VULNERABILITIES_SECURITY_SCAN_TYPE_VIEW_LOCATION: constants.VUL_BY_SCANTYPE_COLUMNS,
	constants.VULNERABILITIES_OVERVIEW:                         constants.VUL_OVERVIEW_COLUMN,
	constants.VULNERABILITIES_OVERVIEW_VIEW_LOCATION:           constants.VUL_OVERVIEW_COLUMN,
	constants.MTTR_FOR_VULNERABILITIES:                         constants.OPEN_VUL_FILTER,
	constants.CWE_TOP_25_VULNERABILITIES:                       constants.OPEN_VUL_FILTER,
	constants.CWE_TOP_25_VULNERABILITIES_VIEW_LOCATION:         constants.OPEN_VUL_FILTER,
	constants.SECURITY_SLA_STATUS_OVERVIEW_OPEN:                constants.VUL_OVERVIEW_COLUMN,
	constants.SECURITY_SLA_STATUS_OVERVIEW_CLOSED:              constants.VUL_OVERVIEW_COLUMN,
	constants.FLOW_METRICS_VELOCITY:                            constants.FLOW_METRICS_FILTER,
	constants.FLOW_METRICS_WORK_ITEM_DISTRIBUTION:              constants.FLOW_METRICS_FILTER,
	constants.FLOW_METRICS_CYCLE_TIME:                          constants.FLOW_METRICS_FILTER,
	constants.FLOW_METRICS_WORK_EFFICIENCY:                     constants.FLOW_METRICS_FILTER,
	constants.FLOW_METRICS_WORK_LOAD:                           constants.FLOW_METRICS_FILTER,
	constants.ACTIVE_DEVELOPERS:                                constants.COMPONENT,
	constants.ACTIVE_DEVELOPERS_COMMITS:                        constants.COMPONENT,
	constants.TRIVY_LICENSE_OCCURENCE:                          constants.OPEN_VUL_FILTER,
	constants.OPEN_ISSUES_DRILL_DOWN:                           constants.OPEN_VUL_FILTER,
	constants.LATEST_TEST_RESULTS:                              "",
}

var DashboardWidgetsDefinitionMap = map[string]string{
	"software-delivery-activity":    "swdelivery/widgets.json",
	"security-insights":             "security/widgets.json",
	"dora-metrics":                  "dorametrics/widgets.json",
	"flow-metrics":                  "flow/widgets.json",
	"component-security-overview":   "componentSecurity/widgets.json",
	"application-security-overview": "applicationSecurity/widgets.json",
}

// FieldIndexMap holds information as to which indices, fields commonly used for filtering could be found in.
var FieldIndexMap = map[string]map[string]string{
	"branch_id": {"cb_test_suites": "", "cb_test_cases": "", "cb_security_findings_remediation_trend": "", "cb_security_findings": ""},
}

// FilterWidgetMap holds information as to which widgets, a particular filter is applicable for.
var FilterWidgetMap = map[string]map[string]string{
	"severities":                    {"so3": "", "so5": "", "so6": "", "so7": "", "so8": "", "so9": "", "so10": "", "so11": "", "aso3": "", "aso4": "", "aso6": "", "aso7": "", "aso8": "", "aso9": "", "aso10": "", "aso11": "", "aso12": ""},
	"tools":                         {"so2": "", "so3": "", "so4": "", "so5": "", "so6": "", "so7": "", "so8": "", "so9": "", "so10": "", "so11": "", "aso2": "", "aso3": "", "aso4": "", "aso5": "", "aso6": "", "aso7": "", "aso8": "", "aso9": "", "aso10": "", "aso11": "", "aso12": ""},
	"sla":                           {"so2": "", "so3": "", "so6": "", "so5": "", "so8": "", "so9": "", "so10": "", "aso2": "", "aso3": "", "aso4": "", "aso6": "", "aso7": "", "aso9": "", "aso10": "", "aso11": ""},
	"application_dashboard_widgets": {"aso1": "", "aso2": "", "aso3": "", "aso4": "", "aso5": "", "aso6": "", "aso7": "", "aso8": "", "aso9": "", "aso10": "", "aso11": "", "aso12": ""},
}

// CustomAggrMap lets you override the default aggregation duration for a particular widget
var CustomAggrMap = map[string]string{
	"so1":  constants.DURATION_DAY,
	"aso1": constants.DURATION_DAY,
}
