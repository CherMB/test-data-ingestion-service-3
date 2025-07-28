package constants

const (
	AGGREGATION               = "aggregations"
	DISTINCT_AUTOMATION       = "distinct_automation"
	DISTINCT_COMPONENT        = "distinct_component"
	DISTINCT_RUN              = "distinct_run"
	TEST_WORKFLOW_DRILLDOWN   = "test_workflow_drilldown"
	ENVIRONMENTS              = "environments"
	COMPONENTS                = "components"
	VALUE                     = "value"
	KEY                       = "key"
	DATA                      = "data"
	RUNS_INFO                 = "runsInfo"
	ID                        = "id"
	INFO                      = "info"
	ACTIVE                    = "Active"
	ACTIVE_KEY                = "active"
	INACTIVE                  = "Inactive"
	WITHSCANNERS              = "With scanners"
	WITHOUTSCANNERS           = "Without scanners"
	BUNDLED_SONARQUBE         = "SonarQube"
	AUTOMATION_STATUS_SUCCESS = "Success"
	NO_SCANNERS               = "No Scanners"
	SCANNERS                  = "Scanners"
	SCANNERS_LOWERCASE        = "scanners"
	SCANNER_TYPE              = "ScannerType"
	SCAN_STATUS               = "scanStatus"
	SCANNER_NAME              = "ScannerName"
	SEC_SCANNER_TYPE          = "scannerType"
	SCANNED                   = "Scanned"
	NOT_SCANNED               = "Not Scanned"
	NOT_APPLICABLE            = "N/A"
	SECURITY_COMPONENT        = "security-components"
	SECURITY_WORKFLOWS        = "security-workflows"
	ORG_ID                    = "orgId"
	SUB_ORG_ID                = "subOrgId"
	BRANCHES                  = "branches"
	REPOS                     = "repositories"
	COMPONENT                 = "component"
	USER_ID                   = "userId"
	ALL                       = "All"
	REPLACE_HEADER_WIDGET     = "e9-sub"
	REPLACE_SUMMARY_WIDGET    = "cs3-sub"
	TRIVY_LICENSE_SECTION     = "trivyLicenseSection"
	TEST_INSIGHTS_AUTOMATIONS = "test_insights_automations"
	WITH_TEST_SUITES          = "With test suites"
	WITHOUT_TEST_SUITES       = "Without test suites"
	TEST_SUITE                = "testSuiteType"
	TEST_INSIGHTS_WORKFLOWS   = "test-insights-workflows"
	TEST_INSIGHTS_COMPONENT   = "test-insights-components"

	REPLACE_HEADER_KEY              = "environment"
	DEPLOY_DATA_INDEX               = "deploy_data"
	BUILD_DATA_INDEX                = "build_data"
	AUTOMATION_RUN                  = "automation_run"
	AUTOMATION_RUN_STATUS_INDEX     = "automation_run_status"
	AUTOMATION_METADATA_INDEX       = "automation_metadata"
	COMMIT_DATA_INDEX               = "commit_data"
	PULL_REQUESTS_REVIEW_DATA_INDEX = "pull_requests_review_data"
	COMPUTED_INDEX                  = "compute_data"
	RAW_SCAN_RESULTS_INDEX          = "raw_scan_result"
	CB_CI_TOOL_INSIGHT_INDEX        = "cb_ci_tool_insight"
	CB_CI_JOB_INFO_INDEX            = "cb_ci_job_info"
	CB_CI_RUN_INFO_INDEX            = "cb_ci_run_info"
	CB_CI_RUNS_ACTIVITY             = "cb_ci_runs_activity"
	CB_CI_ACTIVITY_OVERVIEW         = "cb_ci_activity_overview"
	CB_CI_CJOC_CONTROLLER_INFO      = "cb_ci_cjoc_controller_info"
	CB_VSM_DASHBOARD_LAYOUT         = "cb_vsm_dashboard_layout"

	TOTAL_COUNT              = "totalCount"
	REQUEST_BRANCH           = "branch"
	DUPLICATE_FILES          = "Files with duplication"
	DUPLICATE_LINES          = "Duplicate lines"
	DUPLICATE_BLOCKS         = "Duplicate blocks"
	TOTAL_LINES              = "Total lines"
	CURRENT_SCAN             = "Current scan"
	PREVIOUS_SCAN            = "Previous scan"
	LINES_COVERED            = "Lines covered"
	TOTAL_CODE_LINES         = "Total code lines"
	LINES_TO_COVER           = "Lines to cover"
	CODE_SMELL               = "Code Smell"
	BUG                      = "Bug"
	VULNERABILITY            = "Vulnerability"
	SECURITY_HOTSPOTS        = "Security hotspots"
	CODE_SMELL_PREFIX        = "CODE_SMELL"
	SECURITY_HOTSPOTS_PREFIX = "SECURITY_HOTSPOT"
	BUG_PREFIX               = "BUG"
	VULNERABILITY_PREFIX     = "VULNERABILITY"
	SUB_HEADER               = "subHeader"
	VULNERABILITIES          = "Vulnerabilities"
	PROVIDER                 = "provider"
	PROJECT_TYPE             = "projectType"
	REQUEST_RUN_NUMBER       = "runNumber"

	SDA_DASHBOARD                  = "software-delivery-activity"
	SECURITY_INSIGHTS_DASHBOARD    = "security-insights"
	DORA_METRICS_DASHBOARD         = "dora-metrics"
	FLOW_METRICS_DASHBOARD         = "flow-metrics"
	COMPONENT_SUMMARY_DASHBOARD    = "component-summary"
	CI_INSIGHTS_DASHBOARD          = "ci-insights"
	COMPONENT_SECURITY_DASHBOARD   = "component-security-overview"
	APPLICATION_SECURITY_DASHBOARD = "application-security-overview"
	PARENTS_IDS                    = "parentIds"
	ENVIRONMENT_ENDPOINT           = "cb.configuration.basic-environment"
	TEST_INSIGHTS_DASHBOARD        = "test-insights"

	COLUMN_TYPE = "columnType"

	DATE_LAYOUT                          = "2006-01-02T15:04:05Z"
	DATE_LAYOUT_TZ                       = "2006-01-02T15:04:05.000Z07:00"
	DATE_FORMAT                          = "2006/01/02 15:04:05"
	DATE_FORMAT_TZ                       = "2006/01/02 15:04:05"
	DATE_FORMAT_WITH_HYPHEN              = "2006-01-02 15:04:05"
	DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME = "2006-01-02"
	DB_REPORT_DEF_FILEPATH_TEMPLATE      = "report.definition.filepath"

	REQUEST_TOOLS      = "tools"
	REQUEST_SEVERITIES = "severities"
	REQUEST_SLA        = "sla"

	OPEN_FINDINGS_BY_SEVERITY_WIDGET_ID             = "so2"
	RISK_ACCEPTED_FALSE_POSITIVE_FINDINGS_WIDGET_ID = "so3"
	SLA_BREACHES_BY_SEVERITY_WIDGET_ID              = "so4"
	FINDINGS_IDENTIFIED_WIDGET_ID                   = "so5"
	OPEN_FINDINGS_BY_REVIEW_STATUS_WIDGET_ID        = "so6"
	OPEN_FINDINGS_BY_SLA_STATUS                     = "so7"
	OPEN_FINDINGS_BY_SECURITY_TOOL                  = "so8"
	OPEN_FINDINGS_DISTRIBUTION_BY_SECURITY_TOOL     = "so9"
	OPEN_FINDINGS_DISTRIBUTION_BY_CATEGORY          = "so10"
	SLA_BREACHES_BY_ASSET_TYPE                      = "so11"

	APPLICATION_VULNERABILITY_FINDINGS_REMEDIATION_TREND_WIDGET_ID    = "aso1"
	APPLICATION_OPEN_FINDINGS_BY_SEVERITY_WIDGET_ID                   = "aso2"
	APPLICATION_COMPONENTS_WITH_MOST_OPEN_FINDINGS_WIDGET_ID          = "aso3"
	APPLICATION_RISK_ACCEPTED_FALSE_POSITIVE_FINDINGS_WIDGET_ID       = "aso4"
	APPLICATION_SLA_BREACHES_BY_SEVERITY_WIDGET_ID                    = "aso5"
	APPLICATION_FINDINGS_IDENTIFIED_WIDGET_ID                         = "aso6"
	APPLICATION_OPEN_FINDINGS_BY_REVIEW_STATUS_WIDGET_ID              = "aso7"
	APPLICATION_OPEN_FINDINGS_BY_SLA_STATUS_WIDGET_ID                 = "aso8"
	APPLICATION_OPEN_FINDINGS_BY_SECURITY_TOOL_WIDGET_ID              = "aso9"
	APPLICATION_OPEN_FINDINGS_DISTRIBUTION_BY_CATEGORY_WIDGET_ID      = "aso10"
	APPLICATION_OPEN_FINDINGS_DISTRIBUTION_BY_SECURITY_TOOL_WIDGET_ID = "aso11"
	APPLICATION_SLA_BREACHES_BY_ASSET_TYPE_WIDGET_ID                  = "aso12"
)

var ComputedDrilldown = []string{"test-insights-workflows", "test-insights-components", "component", "workflows", "workflowRuns", "commits", "pullrequests", "runInitiatingCommits", "builds", "deployments", "successfulBuildsDuration",
	"security-components", "security-workflows", "security-workflowRuns", "vulnerabilitiesOverview", "openVulnerabilities", "security-scan-type-workflows", "vulnerabilitiesSecurityScanType", "mttrForVulnerabilities",
}

// CWETop25IDsList is the latest list of top 25 CWE IDs as reported by MITRE
var CWETop25IDsList = []string{
	"CWE-787",
	"CWE-79",
	"CWE-89",
	"CWE-416",
	"CWE-78",
	"CWE-20",
	"CWE-125",
	"CWE-22",
	"CWE-352",
	"CWE-434",
	"CWE-862",
	"CWE-476",
	"CWE-287",
	"CWE-190",
	"CWE-502",
	"CWE-77",
	"CWE-119",
	"CWE-798",
	"CWE-918",
	"CWE-306",
	"CWE-362",
	"CWE-269",
	"CWE-94",
	"CWE-863",
	"CWE-276",
}

const ComponentFilterQuery = `{
    "size": 0,
    "query": {
        "bool": {
          "filter": [
            {
              "term": {
                "org_id": "{{.orgId}}"
              }
            },
            {
              "range": {
                "last_active_time": {
                  "gte": "{{.startDate}}",
                  "lte": "{{.endDate}}",
                  "format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis"
                }
              }
            }
          ]
        }
      },
      "aggs": {
        "distinct_component": {
          "scripted_metric": {
            "init_script": "state.data_map=[:];",
            "map_script": "def map = state.data_map;def key = doc.component_id.value;def v = ['component_id': doc.component_id.value];map.put(key, v);",
            "combine_script": "return state.data_map;",
            "reduce_script": "def tmpMap = [: ]; def resultMap = new HashMap(); for (response in states) { if (response != null) { for (key in response.keySet()) { def record = response.get(key); def compKey = record.component_id; tmpMap.put(compKey, record); } } } def components = []; for (component in tmpMap.keySet()) { components.add(component) } return components;"
          }
        }
      }
  }`

const TestComponentFilterQuery = `{
    "size": 0,
    "query": {
        "bool": {
            "filter": [
                {
                    "term": {
                        "org_id": "{{.orgId}}"
                    }
                },
                {
                    "range": {
                        "run_start_time": {
                            "gte": "{{.startDate}}",
                            "lte": "{{.endDate}}",
                            "format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis",
                            "time_zone":  "{{.timeZone}}"
                        }
                    }
                }
            ]
        }
    },
    "aggs": {
        "distinct_component": {
            "scripted_metric": {
                "init_script": "state.data_map=[:];",
                "map_script": "def map = state.data_map;def key = doc.component_id.value;def v = ['component_id': doc.component_id.value];map.put(key, v);",
                "combine_script": "return state.data_map;",
                "reduce_script": "def tmpMap=[:];def resultMap=new HashMap();def components=new HashSet();for(response in states){if(response!=null){for(key in response.keySet()){def record=response.get(key);components.add(record.component_id);}}}return components;"
            }
        }
    }
}`

const SecurityComponentFilterQuery = `{
    "size": 0,
    "query": {
        "bool": {
          "filter": [
            {
              "term": {
                "org_id": "{{.orgId}}"
              }
            },
            {
              "range": {
                "timestamp": {
                  "gte": "{{.startDate}}",
                  "lte": "{{.endDate}}",
                  "format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis"
                }
              }
            }
          ]
        }
      },
      "aggs": {
        "distinct_component": {
          "scripted_metric": {
            "init_script": "state.data_map=[:];",
            "map_script": "def map = state.data_map;def key = doc.component_id.value;def v = ['component_id': doc.component_id.value];map.put(key, v);",
            "combine_script": "return state.data_map;",
            "reduce_script": "def tmpMap = [: ]; def resultMap = new HashMap(); for (response in states) { if (response != null) { for (key in response.keySet()) { def record = response.get(key); def compKey = record.component_id; tmpMap.put(compKey, record); } } } def components = []; for (component in tmpMap.keySet()) { components.add(component) } return components;"
          }
        }
      }
  }`

const SecurityComponentFilterQueryByScanTime = `{
    "size": 0,
    "query": {
        "bool": {
          "filter": [
            {
              "term": {
                "org_id": "{{.orgId}}"
              }
            },
            {
              "range": {
                "scan_time": {
                  "gte": "{{.startDate}}",
                  "lte": "{{.endDate}}",
                  "format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis"
                }
              }
            }
          ]
        }
      },
      "aggs": {
        "distinct_component": {
          "scripted_metric": {
            "init_script": "state.data_map=[:];",
            "map_script": "def map = state.data_map;def key = doc.component_id.value;def v = ['component_id': doc.component_id.value];map.put(key, v);",
            "combine_script": "return state.data_map;",
            "reduce_script": "def tmpMap = [: ]; def resultMap = new HashMap(); for (response in states) { if (response != null) { for (key in response.keySet()) { def record = response.get(key); def compKey = record.component_id; tmpMap.put(compKey, record); } } } def components = []; for (component in tmpMap.keySet()) { components.add(component) } return components;"
          }
        }
      }
  }`

const AutomationFilterQuery = `{
    "size": 0,
    "query": {
        "bool": {
          "filter": [
            {
              "term": {
                "org_id": "{{.orgId}}"
              }
            },
            {
              "range": {
                "last_active_time": {
                  "gte": "{{.startDate}}",
                  "lte": "{{.endDate}}",
                  "format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis"
                }
              }
            }
          ]
        }
      },
      "aggs": {
        "distinct_automation": {
          "scripted_metric": {
            "init_script": "state.data_map=[:];",
            "map_script": "def map = state.data_map;def key = doc.automation_id.value;def v = ['automation_id': doc.automation_id.value];map.put(key, v);",
            "combine_script": "return state.data_map;",
            "reduce_script": "def tmpMap = [: ]; def resultMap = new HashMap(); for (response in states) { if (response != null) { for (key in response.keySet()) { def record = response.get(key); def autKey = record.automation_id; tmpMap.put(autKey, record); } } } def automations = []; for (automation in tmpMap.keySet()) { automations.add(automation) } return automations;"
          }
        }
      }
  }`

const AutomationFilterQueryWithBranch = `{
    "size": 0,
    "query": {
        "bool": {
          "filter": [
            {
              "term": {
                "org_id": "{{.orgId}}"
              }
            },
            {
                "term": {
                  "branch_id": "{{.branch}}"
                }
            },
            {
              "range": {
                "last_active_time": {
                  "gte": "{{.startDate}}",
                  "lte": "{{.endDate}}",
                  "format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis"
                }
              }
            }
          ]
        }
      },
      "aggs": {
        "distinct_automation": {
          "scripted_metric": {
            "init_script": "state.data_map=[:];",
            "map_script": "def map = state.data_map;def key = doc.automation_id.value;def v = ['automation_id': doc.automation_id.value];map.put(key, v);",
            "combine_script": "return state.data_map;",
            "reduce_script": "def tmpMap = [: ]; def resultMap = new HashMap(); for (response in states) { if (response != null) { for (key in response.keySet()) { def record = response.get(key); def autKey = record.automation_id; tmpMap.put(autKey, record); } } } def automations = []; for (automation in tmpMap.keySet()) { automations.add(automation) } return automations;"
          }
        }
      }
  }`
const SecurityAutomationFilterQuery = `{
    "size": 0,
    "query": {
        "bool": {
          "filter": [
            {
              "term": {
                "org_id": "{{.orgId}}"
              }
            },
            {
              "range": {
                "timestamp": {
                  "gte": "{{.startDate}}",
                  "lte": "{{.endDate}}",
                  "format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis"
                }
              }
            }
          ]
        }
      },
      "aggs": {
        "distinct_automation": {
          "scripted_metric": {
            "init_script": "state.data_map=[:];",
            "map_script": "def map = state.data_map;def key = doc.automation_id.value;def v = ['automation_id': doc.automation_id.value];map.put(key, v);",
            "combine_script": "return state.data_map;",
            "reduce_script": "def tmpMap = [: ]; def resultMap = new HashMap(); for (response in states) { if (response != null) { for (key in response.keySet()) { def record = response.get(key); def autKey = record.automation_id; tmpMap.put(autKey, record); } } } def automations = []; for (automation in tmpMap.keySet()) { automations.add(automation) } return automations;"
          }
        }
      }
  }`

const TestInsightsAutomationFilterQuery = `{
    "size": 0,
    "query": {
        "bool": {
          "filter": [
            {
              "term": {
                "org_id": "{{.orgId}}"
              }
            },
            {
              "range": {
                "run_start_time": {
                  "gte": "{{.startDate}}",
                  "lte": "{{.endDate}}",
                  "format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis",
                  "time_zone": "{{.timeZone}}"
                }
              }
            }
          ]
        }
      },
      "aggs": {
        "test_insights_automations": {
          "scripted_metric": {
            "init_script": "state.data_map=[:];",
            "map_script": "def map = state.data_map;def key = doc.automation_id.value;def v = ['automation_id': doc.automation_id.value];map.put(key, v);",
            "combine_script": "return state.data_map;",
            "reduce_script": "def tmpMap=[:];def resultMap=new HashMap();def automations=new HashSet();for(response in states){if(response!=null){for(key in response.keySet()){def record=response.get(key);automations.add(record.automation_id)}}}return automations;"
          }
        }
      }
  }`

const FetchAllBuildComponents = `{
    "size": 0,
    "query": {
      "bool": {
        "filter": [
         {
                    "range": {
                      "status_timestamp": {
                        "gte": "{{.startDate}}",
                        "lte": "{{.endDate}}",
                        "format": "yyyy-MM-dd HH:mm:ss",
                        "time_zone":"{{.timeZone}}"
                      }
                    }
                  },
                  {
                    "term": {
                      "org_id": "{{.orgId}}"
  
            }
          },
          {
            "term": {
              "status": "SUCCEEDED"
            }
          },
          {
		  	"term": {
				"data_type": 2
			}
		  }
        ]
      }
    },
    "aggs": {
        "components": {
          "scripted_metric": {
            "init_script": "state.statusMap = [:];",
              "map_script": "def map = state.statusMap;def key = doc.component_id.value + '_' + doc.run_id.value + '_' + doc.job_id.value + '_' + doc.step_id.value + '_' + doc.status_timestamp.value;def v = ['component_id': doc.component_id.value, 'run_id': doc.run_id.value, 'job_id': doc.job_id.value, 'step_id': doc.step_id.value, 'step_kind': doc.step_kind.value, 'target_env': doc.target_env.value, 'status': doc.status.value, 'status_timestamp': doc.status_timestamp.value, 'start_time': doc.start_time.value, 'completed_time': doc.completed_time.value];map.put(key, v);",
              "combine_script": "return state.statusMap;",
              "reduce_script": "float getMedian(def input) {def q2;def count = input.size();if (count % 2 == 0) {q2 = (float)(input.get((count / 2) - 1) + input.get(count / 2)) / 2;} else {q2 = input.get(count / 2);}return (float) q2;}def statusMap = new HashMap();def durationCountMap = new HashMap();def resultMap = new HashMap();Instant Currentdate = Instant.ofEpochMilli(new Date().getTime());for (response in states) {if (response != null) {for (key in response.keySet()) {statusMap.put(key, response.get(key));}}}for (uniqueKey in statusMap.keySet()) {def build = statusMap.get(uniqueKey);if (build.start_time != 0 && build.completed_time != 0) {Instant startDate = Instant.ofEpochMilli(build.start_time);Instant completedDate = Instant.ofEpochMilli(build.completed_time);def duration = ChronoUnit.MILLIS.between(startDate, completedDate);if (durationCountMap.containsKey(build.component_id)) {def durationList = durationCountMap.get(build.component_id);durationList.add(duration);durationCountMap.put(build.component_id, durationList);} else {def durationList = new ArrayList();durationList.add(duration);durationCountMap.put(build.component_id, durationList);}}}def sortedComponents = new ArrayList(), min = -1, max = 0;def boxPlotMap = new HashMap();for (uniqueComponent in durationCountMap.keySet()) {def sortedValues = durationCountMap.get(uniqueComponent);Collections.sort(sortedValues);def firstHalf = new ArrayList();def secondHalf = new ArrayList();def q1, q2, q3, firstHalfToIndex, secondHalfFromIndex;int count = sortedValues.size();if (count < 2) {def tempMap = [sortedValues[0], sortedValues[0], sortedValues[0], sortedValues[0], sortedValues[0]];boxPlotMap.put(uniqueComponent, sortedValues[0]);if (min == -1 || sortedValues[0] < min) {min = sortedValues[0];}if (sortedValues[0] > max) {max = sortedValues[0];}} else {if (count % 2 == 0) {firstHalfToIndex = (count / 2);secondHalfFromIndex = (count / 2);} else {firstHalfToIndex = (count / 2);secondHalfFromIndex = (count / 2) + 1;}q2 = getMedian(sortedValues);firstHalf = sortedValues.subList(0, firstHalfToIndex);secondHalf = sortedValues.subList(secondHalfFromIndex, count);q1 = getMedian(firstHalf);q3 = getMedian(secondHalf);def iqr = q3 - q1;def whiskerMin = q1 - 1.5 * iqr;def whiskerMax = q3 + 1.5 * iqr;for (val in sortedValues) {if (val >= whiskerMin) {whiskerMin = val;break;}}for (int i = sortedValues.size() - 1; i >= 0; i--) {if (sortedValues.get(i) <= whiskerMax) {whiskerMax = sortedValues.get(i);break;}}boxPlotMap.put(uniqueComponent, whiskerMax);if (min == -1 || whiskerMin < min) {min = whiskerMin;}if (whiskerMax > max) {max = whiskerMax;}}}def sortedEntries = boxPlotMap.entrySet().stream().sorted((a, b) -> -a.getValue().compareTo(b.getValue())).collect(Collectors.toList());for (entry in sortedEntries) {sortedComponents.add(entry.getKey());}resultMap.put('componentsInfo', boxPlotMap);resultMap.put('components', sortedComponents);if (min == -1) {min = 0;}resultMap.put('min', '' + min);resultMap.put('max', '' + max);return resultMap;"
          }
        }
    }
  }`

const FetchAllDeployEnv = `{
    "size": 0,
    "query": {
      "bool": {
        "filter": [
          {
            "range": {
                "status_timestamp": {
                    "gte": "{{.startDate}}",
                    "lte": "{{.endDate}}",
                    "format": "yyyy-MM-dd HH:mm:ss",
                    "time_zone": "{{.timeZone}}"
                }
            }
          },
          {
            "term": {
                "org_id": "{{.orgId}}"
            }
          },
          {
            "term": {
              "status": "SUCCEEDED"
            }
          },
          {
            "bool": {
                "must_not": [
                    {
                        "term": {
                            "target_env": ""
                        }
                    }
                ]
            }
          }
        ]
      }
    },
    "aggs": {
      "environments": {
        "scripted_metric": {
          "init_script": "state.data_set=new HashSet();",
          "map_script": "def set = state.data_set; if (doc['target_env'].size() != 0) {set.add(doc.target_env.value);}",
          "combine_script": "return state.data_set;",
          "reduce_script": "def resultSet = new HashSet(); for (responses in states) { if (responses != null) { for (response in responses) { resultSet.add(response); } } } return resultSet; "
        }
      }
    }
  }`

const AutomationRunCountQuery = `{
    "_source": false,
    "size": 0,
    "query": {
      "bool": {
        "filter": [
            {
                "range": {
                  "status_timestamp": {
                    "gte": "{{.startDate}}",
                    "lte": "{{.endDate}}",
                    "format": "yyyy-MM-dd HH:mm:ss"
                  }
                }
            },
            {
                "term": {
                  "org_id": "{{.orgId}}"
                }
            },
            {
                "term": {
                    "job_id": ""
                }
            },
            {
                "term": {
                    "step_id": ""
                }
            },
            {
                "term": {
                    "data_type": 2
                }
            },
            {
                "bool": {
                  "should": [
                    {
                      "term": {
                        "status": "SUCCEEDED"
                      }
                    },
                    {
                      "term": {
                        "status": "FAILED"
                      }
                    },
                    {
                      "term": {
                        "status": "TIMED_OUT"
                      }
                    },
                    {
                      "term": {
                        "status": "ABORTED"
                      }
                    }
                  ]
                }
            }
        ]
      }
    },
    "aggs": {
        "automation_run": {
          "scripted_metric": {
            "init_script": "state.data_map=[:];",
            "map_script": "def map = state.data_map;def key = doc.org_id+'_'+doc.automation_id.value + '_' + doc.run_id.value+'_'+doc.status.value;def v = ['run_id': doc.run_id.value, 'status': doc.status.value,'automation_id':doc.automation_id.value];map.put(key, v);",
            "combine_script": "return state.data_map;",
            "reduce_script": "def resultMap = new HashMap();def tmpMap = new HashMap();def dataMap = new HashMap();for (response in states) {if (response != null) {for (key in response.keySet()) {tmpMap.put(key, response.get(key));}}}for (key in tmpMap.keySet()) {def record = tmpMap.get(key);def runKey = record.automation_id+'_'+record.run_id;if (dataMap.containsKey(runKey)) {def count = dataMap.get(runKey);dataMap.put(runKey, count + 1);} else {dataMap.put(runKey, 1);}}resultMap.put('data', dataMap);resultMap.put('totalCount', dataMap.size());return resultMap;"
          }
        }
    }
  }`
const DeployedAutomationCountQuery = `{
    "_source": false,
    "size": 0,
    "query": {
      "bool": {
        "filter": [
            {
                "range": {
                  "status_timestamp": {
                    "gte": "{{.startDate}}",
                    "lte": "{{.endDate}}",
                    "format": "yyyy-MM-dd HH:mm:ss"
                  }
                }
            },
            {
                "term": {
                  "org_id": "{{.orgId}}"
                }
            },
            {
                "term": {
                  "status": "SUCCEEDED"
                }
            },
            {
                "term": {
                    "data_type": 2
                }
            },
            {
                "bool": {
                    "must_not": [
                        {
                            "term": {
                              "target_env": ""
                            }
                        }
                    ]
                }
            }
        ]
      }
    },
    "aggs": {
        "automation_run": {
          "scripted_metric": {
            "init_script": "state.data_map=[:];",
            "map_script": "def map = state.data_map;def key = doc.automation_id.value + '_' + doc.run_id.value + '_' + doc.job_id.value + '_' + doc.step_id.value + '_' + doc.status.value + '_' + doc.status_timestamp.value;def v = ['automation_id': doc.automation_id.value, 'run_id': doc.run_id.value, 'job_id': doc.job_id.value, 'step_id': doc.step_id.value, 'target_env': doc.target_env.value, 'status': doc.status.value];map.put(key, v);",
            "combine_script": "return state.data_map;",
            "reduce_script": "def resultMap = new HashMap();def tmpMap = new HashMap();def dataMap = new HashMap();for (response in states) {if (response != null) {for (key in response.keySet()) {tmpMap.put(key, response.get(key));}}}for (key in tmpMap.keySet()) {def record = tmpMap.get(key);def automationKey = record.automation_id + '_'+record.run_id;def envKey = record.target_env;if (dataMap.containsKey(automationKey)) {def envMap = dataMap.get(automationKey);if (envMap.containsKey(envKey)) {envMap.put(envKey, envMap.get(envKey) + 1);} else {envMap.put(envKey, 1);}} else {def envMap = new HashMap();envMap.put(envKey, 1);dataMap.put(automationKey, envMap);}}resultMap.put('data', dataMap);return resultMap;"
          }
        }
      }
  }`

const ComputedWidgetQuery = `{     
    "size": 1,
    "_source": ["metric_value"], 
    "query": {
        "bool": {
            "must": [
                {
                    "term" : {
                        "org_id" : "{{.orgId}}"
                    }
                },
                {
                    "term" : {
                        "metric_key": "{{.metricsKey}}"
                    }
                }, 
                  {
                    "term": {
                      "start_date": {
                        "value": "{{.startDate}}"
                      }
                    }
                  },
                  {
                    "term": {
                      "end_date": {
                        "value": "{{.endDate}}"
                      }
                    }
                  }
                  
              
      
            ]
        }
    }   
}`

var ComponentSummaryWidgetsLayout = map[string]string{
	"cs1":  `{"widgetId":"cs1","widgetName":"components-activity","widgetWidth":12,"widgetHeight":2,"mockData":false}`,
	"cs2":  `{"widgetId":"cs2","widgetName":"components-builds","widgetWidth":6,"widgetHeight":3,"mockData":false}`,
	"cs3":  `{"widgetId":"cs3","widgetName":"components-deployments","widgetWidth":6,"widgetHeight":3,"mockData":false}`,
	"cs4":  `{"widgetId":"cs4","widgetName":"components-code-coverage","widgetWidth":4,"widgetHeight":3,"mockData":false}`,
	"cs5":  `{"widgetId":"cs5","widgetName":"components-issue-types","widgetWidth":4,"widgetHeight":3,"mockData":false}`,
	"cs6":  `{"widgetId":"cs6","widgetName":"components-duplications","widgetWidth":4,"widgetHeight":3,"mockData":false}`,
	"cs7":  `{"widgetId":"cs7","widgetName":"components-codebase-overview","widgetWidth":12,"widgetHeight":4,"mockData":false}`,
	"cs8":  `{"widgetId":"cs8","widgetName":"components-open-issues","widgetWidth":12,"widgetHeight":4,"mockData":false}`,
	"cs9":  `{"widgetId":"cs9","widgetName":"components-no-scanners-configured","widgetWidth":12,"widgetHeight":2,"mockData":false}`,
	"cs10": `{"widgetId":"cs10","widgetName":"components-trivy-licenses-overview","widgetWidth":12,"widgetHeight":3,"mockData":false}`,
}

var TransformValidationsMap = map[string]string{
	"f8_flowWaitTimeChart":        `"aggregations":{"flow_wait_time_count":{"value":[]`,
	"f8_flowEfficiencyChart":      `"aggregations":{"flow_eff_time_buckets":{"buckets":[]}}`,
	"f5_flowDistributionAvgChart": `"aggregations":{"flow_distribution_avg_count":{"value":[]`,
	"f5_flowDistributionChart":    `"aggregations":{"flow_distribution_buckets":{"buckets":[]`,
}

var ComponentSummarySonarWidgets = []string{"cs1", "cs2", "cs3", "cs4", "cs5", "cs6", "cs7", "cs8", "cs10"}
var ComponentSummaryNoSonarWidgets = []string{"cs1", "cs2", "cs3", "cs8"}
var ComponentSummaryNoScannerWidgets = []string{"cs1", "cs2", "cs3", "cs9"}

const CURRENT_MONTH = "CURRENT_MONTH"
const CURRENT_WEEK = "CURRENT_WEEK"
const COMPUTE_INDEX = "compute_data"

var Durations = []string{"CURRENT_WEEK", "CURRENT_MONTH"}

const Latest_Sonar_Query = `{
    "size": 2,
    "sort": [{
        "timestamp": {
            "order": "desc"
        }
    }],
    "query": {
        "bool": {
            "filter": [
                {
                    "term": {
                        "org_id": {
                            "value": "{{.orgId}}"
                        }
                    }
                }
            ]
        }
    }
}`

const Sonar_Query = `{
    "size": 2,
    "sort": [{
        "timestamp": {
            "order": "desc"
        }
    }],
    "query": {
        "bool": {
            "filter": [{
                "range": {
                  "timestamp": {
                    "gte": "{{.startDate}}",
                    "lte": "{{.endDate}}",
                    "format": "yyyy-MM-dd HH:mm:ss",
                    "time_zone":"{{.timeZone}}"
                  }
                }
            },
                {
                    "term": {
                        "org_id": {
                            "value": "{{.orgId}}"
                        }
                    }
                }
            ]
        }
    }
}`

const Sonar_Base_Query = `{
    "size": 2,
    "sort": [{
        "timestamp": {
            "order": "desc"
        }
    }],
    "query": {
        "bool": {
            "filter": [
                {
                    "term": {
                        "org_id": {
                            "value": "{{.orgId}}"
                        }
                    }
                },
                {
                    "term":{
                        "component_id":{
                            "value":"{{.component}}"
                        }
                    }
                }
            ]
        }
    }
}`

const OpenIssuesSectionQuery = `{
    "size": 0,
    "query": {
        "bool": {
            "filter": [
                {
                    "range": {
                        "scan_time": {
                            "gte": "{{.startDate}}",
                            "lte": "{{.endDate}}",
                            "format": "yyyy-MM-dd HH:mm:ss",
                            "time_zone": "{{.timeZone}}"
                        }
                    }
                },
                {
                    "term": {
                        "org_id": "{{.orgId}}"
                    }
                },
                {
                    "exists": {
                        "field": "date_of_discovery"
                    }
                },
                {
                    "exists": {
                        "field": "github_branch"
                    }
                },
                {
                    "exists": {
                        "field": "scanner_name"
                    }
                },
                {
                    "exists": {
                        "field": "code"
                    }
                },
                {
                    "terms": {
                        "severity": [
                            "MEDIUM",
                            "HIGH",
                            "LOW",
                            "VERY_HIGH"
                        ]
                    }
                }
            ]
        }
    },
    "aggs": {
        "drilldowns": {
            "scripted_metric": {
                "params": {
                    "timeZone": "{{.timeZone}}"
                },
                "init_script": "state.statusMap = [: ];",
                "map_script": "if (doc['standard'].size()==0 || doc['standard'].value == 'STANDARD' || doc['standard'].value == 'SECRET'){def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value,'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':'', 'scanner_name':doc.scanner_name.value, 'run_id':doc.run_id.value,'recurrences':params['_source']['failure_files'].size(),'issue_type':doc['standard'].size()==0?'':doc.standard.value];map.put(key, v);}",
                "combine_script": "return state.statusMap;",
                "reduce_script": "Instant Currentdate = Instant.ofEpochMilli(new Date().getTime());def statusMap = new HashMap();def resultMap = new HashMap();def resultList = new ArrayList();def slaRules = new HashMap();slaRules.put('Breached', 3);slaRules.put('AtRisk', 2);DateTimeFormatter formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone));for (a in states){if (a != null){for (i in a.keySet()){def record = a.get(i);def key = record.org_id + '_' + record.component_id + '_' + record.code;if (statusMap.containsKey(key)){def vulDetailsMap = statusMap.get(key);def vulKey = record.org_id + '_' + record.component_id + '_' + record.code + '_' + record.scanner_name;if (vulDetailsMap.containsKey(vulKey)){def lastRecord = vulDetailsMap.get(vulKey);if (lastRecord.timestamp < record.timestamp){vulDetailsMap.put(vulKey, record);}} else{vulDetailsMap.put(vulKey, record);}} else{def vulDetailsMap = new HashMap();def vulKey = record.org_id + '_' + record.component_id + '_' + record.code + '_' + record.scanner_name;vulDetailsMap.put(vulKey, record);statusMap.put(key, vulDetailsMap);}}}}if (statusMap.size() > 0){for (uniqueVul in statusMap.keySet()){def vulScannersMap = statusMap.get(uniqueVul);for (vulKey in vulScannersMap.keySet()){def curVul = vulScannersMap.get(vulKey);if (curVul.bug_status == 'Open' || curVul.bug_status == 'Reopened'){if (resultMap.containsKey(uniqueVul)){def vulFinalInfo = resultMap.get(uniqueVul);if (curVul.date_of_discovery.getMillis() < vulFinalInfo.date_of_discovery.getMillis()){vulFinalInfo.put('date_of_discovery', curVul.date_of_discovery);}def issueTypeSet = vulFinalInfo.get('issueType');issueTypeSet.add(curVul.issue_type);def scannersList = vulFinalInfo.get('scannerNames');def scannerInfoMap = new HashMap();scannerInfoMap.put('name', curVul.scanner_name);scannersList.add(scannerInfoMap);def recurrences = vulFinalInfo.recurrences;recurrences += curVul.recurrences;vulFinalInfo.put('recurrences', recurrences);def drillDownInfoMap = vulFinalInfo.get('drillDown');def reportInfo = drillDownInfoMap.get('reportInfo');def scannersSet = reportInfo.get('scanner_name_list');scannersSet.add(curVul.scanner_name);def runIDSet = reportInfo.get('run_id_list');runIDSet.add(curVul.run_id);} else{def severityCode = 0;def curSeverity = curVul.severity;if (curSeverity == 'VERY_HIGH'){curVul.severity = 'Very high';severityCode = 4;} else if (curSeverity == 'HIGH'){curVul.severity = 'High';severityCode = 3;} else if (curSeverity == 'MEDIUM'){curVul.severity = 'Medium';severityCode = 2;} else if (curSeverity == 'LOW'){curVul.severity = 'Low';severityCode = 1;}def map = new HashMap();def runIDSet = new HashSet();def scannersSet = new HashSet();def issueTypeSet = new HashSet();def scannersList = new ArrayList();def scannerInfoMap = new HashMap();def reportInfoMap = new HashMap();runIDSet.add(curVul.run_id);reportInfoMap.put('run_id_list', runIDSet);scannersSet.add(curVul.scanner_name);reportInfoMap.put('scanner_name_list', scannersSet);reportInfoMap.put('code', curVul.code);reportInfoMap.put('component_id', curVul.component_id);def drillDownInfoMap = new HashMap();drillDownInfoMap.put('reportId', 'open-issues-drill-down');drillDownInfoMap.put('reportTitle', curVul.code);drillDownInfoMap.put('reportInfo', reportInfoMap);map.put('drillDown', drillDownInfoMap);map.put('date_of_discovery', curVul.date_of_discovery);map.put('vulnerabilityId', curVul.code);map.put('scannerName', curVul.scanner_name);map.put('sla', '');scannerInfoMap.put('name', curVul.scanner_name);scannersList.add(scannerInfoMap);map.put('scannerNames', scannersList);map.put('severity', curVul.severity);map.put('severityCode', severityCode);issueTypeSet.add(curVul.issue_type);map.put('issueType', issueTypeSet);map.put('recurrences', curVul.recurrences);resultMap.put(uniqueVul, map);}}}def finalVulInfo = resultMap.get(uniqueVul);if (finalVulInfo != null){Instant Startdate = Instant.ofEpochMilli(finalVulInfo.date_of_discovery.getMillis());def diffAge = ChronoUnit.DAYS.between(Startdate, Currentdate);if (diffAge >= slaRules.Breached){finalVulInfo.sla = 'Breached';} else if (diffAge >= slaRules.AtRisk){finalVulInfo.sla = 'At risk';} else{finalVulInfo.sla = 'On track';}finalVulInfo.put('dateOfDiscovery', formatter.format(finalVulInfo.date_of_discovery));finalVulInfo.remove('date_of_discovery');resultList.add(finalVulInfo);}}}return resultList;"
            }
        }
    }
}`

const ComponentOpenIssueHeaderQuery = `{
    "size": 0,
    "query": {
        "bool": {
            "filter": [
                {
                    "range": {
                        "scan_time": {
                            "gte": "{{.startDate}}",
                            "lte": "{{.endDate}}",
                            "format": "yyyy-MM-dd HH:mm:ss",
                            "time_zone":"{{.timeZone}}"
                        }
                    }
                },
                {
                    "term": {
                        "org_id": "{{.orgId}}"
                    }
                },
                {
                    "exists": {
                        "field": "date_of_discovery"
                    }
                },
                {
                    "terms": {
                        "severity": [
                            "MEDIUM",
                            "HIGH",
                            "LOW",
                            "VERY_HIGH"
                        ]
                    }
                }
            ]
        }
    },
    "aggs": {
        "severityCounts": {
            "scripted_metric": {
                "init_script": "state.statusMap = [:];",
                "map_script": "if (doc['standard'].size()==0 || doc['standard'].value == 'STANDARD' || doc['standard'].value == 'SECRET'){def map=state.statusMap;def key=doc.org_id.value+'_'+doc.component_id.value+ '_' +doc.github_branch.value +'_'+doc.code.value+'_'+doc.scanner_name.value+'_'+doc['timestamp'].getValue().toEpochSecond()*1000;def v=['org_id': doc.org_id.value, 'component_id': doc.component_id.value,'branch': doc.github_branch.value,'scanner_name': doc.scanner_name.value,'timestamp':doc['timestamp'].getValue().toEpochSecond()*1000,'code':doc.code.value,'bug_status':doc.bug_status.value,'name':doc.name.value, 'component_name':doc.component_name.value, 'severity':doc.severity.value];map.put(key,v);}",
                "combine_script": "return state.statusMap;",
                "reduce_script": "def statusMap = new HashMap();def severityMap = new HashMap();severityMap.put('VERY_HIGH',0);severityMap.put('HIGH',0);severityMap.put('MEDIUM',0);severityMap.put('LOW',0);for (a in states){if (a != null){for (i in a.keySet()){def record = a.get(i);def key = record.org_id + '_' + record.code;if (statusMap.containsKey(key)){def vulDetailsMap = statusMap.get(key);def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record.code + '_' + record.scanner_name;if (vulDetailsMap.containsKey(vulKey)){def lastRecord = vulDetailsMap.get(vulKey);if (lastRecord.timestamp < record.timestamp){vulDetailsMap.put(vulKey, record);}} else{vulDetailsMap.put(vulKey, record);}} else{def vulDetailsMap = new HashMap();def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record.code + '_' + record.scanner_name;vulDetailsMap.put(vulKey, record);statusMap.put(key, vulDetailsMap);}}}}if (statusMap.size() > 0){for (uniqueKey in statusMap.keySet()){def vulMapBranchLevel = statusMap.get(uniqueKey);for (vulKey in vulMapBranchLevel.keySet()){def vul = vulMapBranchLevel.get(vulKey);if (vul.bug_status == 'Open' || vul.bug_status == 'Reopened'){if (severityMap.containsKey(vul.severity)){def count = severityMap.get(vul.severity);severityMap.put(vul.severity, count + 1);break;} else{severityMap.put(vul.severity, 1);break;}}}}}def totalCount = 0;if (severityMap.size() > 0){for (key in severityMap.keySet()){totalCount = totalCount + severityMap.get(key);}}severityMap.put('TOTAL',totalCount);return severityMap;"
            }
        }
    }
}`

const CompareReportsMock1 = `{
  "is_sub_org": true,
  "sub_org_id": "sub-org-1",
  "compare_title": "Sub org 1",
  "sub_org_count": 2,
  "component_count": 4,
  "total_value": 2545,
  "compare_reports": [
      {
          "is_sub_org": true,
          "sub_org_id": "sub-org-3",
          "compare_title": "Sub org 3",
          "sub_org_count": 0,
          "component_count": 2,
          "total_value": 2569,
          "compare_reports": [
              {
                  "is_sub_org": true,
                  "sub_org_id": "sub-org-4",
                  "compare_title": "Suborg 4",
                  "sub_org_count": 1,
                  "component_count": 1,
                  "total_value": 234,
                  "compare_reports": [
                      {
                          "is_sub_org": true,
                          "sub_org_id": "sub-org-5",
                          "compare_title": "Suborg 5",
                          "sub_org_count": 0,
                          "component_count": 1,
                          "total_value": 2545,
                          "compare_reports": [
                              {
                                  "is_sub_org": false,
                                  "compare_title": "Suborg5 - component",
                                  "sub_org_count": 0,
                                  "component_count": 0,
                                  "total_value": 2545,
                                  "section": {
                                      "data": [
                                          {
                                              "title": "Bugs",
                                              "value": 30
                                          },
                                          {
                                              "title": "Feature",
                                              "value": 52
                                          },
                                          {
                                              "title": "Risk",
                                              "value": 52
                                          },
                                          {
                                              "title": "Tech Debt",
                                              "value": 52
                                          }
                                      ]
                                  }
                              }
                          ],
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      },
                      {
                          "is_sub_org": false,
                          "compare_title": "Suborg4 - component",
                          "sub_org_count": 0,
                          "component_count": 0,
                          "total_value": 2545,
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      }
                  ],
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 3 - Component 1",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "total_value": 234,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 3 - Component 2",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "total_value": 234,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              }
          ],
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": true,
          "sub_org_id": "sub-org-9",
          "compare_title": "Sub org 9",
          "sub_org_count": 0,
          "component_count": 2,
          "total_value": 2569,
          "compare_reports": [
              {
                  "is_sub_org": true,
                  "sub_org_id": "sub-org-10",
                  "compare_title": "Suborg 10",
                  "sub_org_count": 1,
                  "component_count": 1,
                  "total_value": 234,
                  "compare_reports": [
                      {
                          "is_sub_org": true,
                          "sub_org_id": "sub-org-11",
                          "compare_title": "Suborg 11",
                          "sub_org_count": 0,
                          "component_count": 1,
                          "total_value": 2545,
                          "compare_reports": [
                              {
                                  "is_sub_org": false,
                                  "compare_title": "Suborg11 - component",
                                  "sub_org_count": 0,
                                  "component_count": 0,
                                  "total_value": 2545,
                                  "section": {
                                      "data": [
                                          {
                                              "title": "Bugs",
                                              "value": 30
                                          },
                                          {
                                              "title": "Feature",
                                              "value": 52
                                          },
                                          {
                                              "title": "Risk",
                                              "value": 52
                                          },
                                          {
                                              "title": "Tech Debt",
                                              "value": 52
                                          }
                                      ]
                                  }
                              }
                          ],
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      },
                      {
                          "is_sub_org": false,
                          "compare_title": "Suborg10 - component",
                          "sub_org_count": 0,
                          "component_count": 0,
                          "total_value": 2545,
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      }
                  ],
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 9 - Component 1",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "total_value": 234,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 9 - Component 2",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "total_value": 234,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              }
          ],
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 1",
          "sub_org_count": 0,
          "component_count": 0,
          "total_value": 234,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 2",
          "sub_org_count": 0,
          "component_count": 0,
          "total_value": 234,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 3",
          "sub_org_count": 0,
          "component_count": 0,
          "total_value": 234,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 4",
          "sub_org_count": 0,
          "component_count": 0,
          "total_value": 234,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      }
  ],
  "section": {
      "data": [
          {
              "title": "Bugs",
              "value": 30
          },
          {
              "title": "Feature",
              "value": 52
          },
          {
              "title": "Risk",
              "value": 52
          },
          {
              "title": "Tech Debt",
              "value": 52
          }
      ]
  }
}`

const CompareReportsMock2 = `{
  "is_sub_org": true,
  "sub_org_id": "sub-org-2",
  "compare_title": "Sub org 2",
  "sub_org_count": 1,
  "component_count": 5,
  "total_value": 2569,
  "compare_reports": [
      {
          "is_sub_org": true,
          "sub_org_id": "sub-org-6",
          "compare_title": "Sub org 6",
          "sub_org_count": 0,
          "component_count": 2,
          "total_value": 2569,
          "compare_reports": [
              {
                  "is_sub_org": true,
                  "sub_org_id": "sub-org-7",
                  "compare_title": "Suborg 7",
                  "sub_org_count": 1,
                  "component_count": 1,
                  "total_value": 234,
                  "compare_reports": [
                      {
                          "is_sub_org": true,
                          "sub_org_id": "sub-org-8",
                          "compare_title": "Suborg 8",
                          "sub_org_count": 0,
                          "component_count": 1,
                          "total_value": 2545,
                          "compare_reports": [
                              {
                                  "is_sub_org": false,
                                  "compare_title": "Suborg8 - component",
                                  "sub_org_count": 0,
                                  "component_count": 0,
                                  "total_value": 2545,
                                  "section": {
                                      "data": [
                                          {
                                              "title": "Bugs",
                                              "value": 30
                                          },
                                          {
                                              "title": "Feature",
                                              "value": 52
                                          },
                                          {
                                              "title": "Risk",
                                              "value": 52
                                          },
                                          {
                                              "title": "Tech Debt",
                                              "value": 52
                                          }
                                      ]
                                  }
                              }
                          ],
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      },
                      {
                          "is_sub_org": false,
                          "compare_title": "Suborg7 - component",
                          "sub_org_count": 0,
                          "component_count": 0,
                          "total_value": 2545,
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      }
                  ],
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 6 - Component 1",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "total_value": 234,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 6 - Component 2",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "total_value": 234,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              }
          ],
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 1",
          "sub_org_count": 0,
          "component_count": 0,
          "total_value": 234,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 2",
          "sub_org_count": 0,
          "component_count": 0,
          "total_value": 234,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 3",
          "sub_org_count": 0,
          "component_count": 0,
          "total_value": 234,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 4",
          "sub_org_count": 0,
          "component_count": 0,
          "total_value": 234,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      }
  ],
  "section": {
      "data": [
          {
              "title": "Bugs",
              "value": 30
          },
          {
              "title": "Feature",
              "value": 52
          },
          {
              "title": "Risk",
              "value": 52
          },
          {
              "title": "Tech Debt",
              "value": 52
          }
      ]
  }
}`

const CompareReportsMock3 = `{"is_sub_org":false,"compare_title":"Component","sub_org_count":0,"component_count":0,"total_value":234,"section":{"data":[{"title":"Bugs","value":30},{"title":"Feature","value":52},{"title":"Risk","value":52},{"title":"Tech Debt","value":52}]}}`

const CompareReportsCycleTimeMock1 = `{
  "is_sub_org": true,
  "sub_org_id": "sub-org-1",
  "compare_title": "Sub org 1",
  "sub_org_count": 2,
  "component_count": 4,
  "value_in_millis": 168429000,
  "compare_reports": [
      {
          "is_sub_org": true,
          "sub_org_id": "sub-org-3",
          "compare_title": "Sub org 3",
          "sub_org_count": 0,
          "component_count": 2,
          "value_in_millis": 168429000,
          "compare_reports": [
              {
                  "is_sub_org": true,
                  "sub_org_id": "sub-org-4",
                  "compare_title": "Suborg 4",
                  "sub_org_count": 1,
                  "component_count": 1,
                  "value_in_millis": 168429000,
                  "compare_reports": [
                      {
                          "is_sub_org": true,
                          "sub_org_id": "sub-org-5",
                          "compare_title": "Suborg 5",
                          "sub_org_count": 0,
                          "component_count": 1,
                          "value_in_millis": 168429000,
                          "compare_reports": [
                              {
                                  "is_sub_org": false,
                                  "compare_title": "Suborg5 - component",
                                  "sub_org_count": 0,
                                  "component_count": 0,
                                  "value_in_millis": 168429000,
                                  "section": {
                                      "data": [
                                          {
                                              "title": "Bugs",
                                              "value": 30
                                          },
                                          {
                                              "title": "Feature",
                                              "value": 52
                                          },
                                          {
                                              "title": "Risk",
                                              "value": 52
                                          },
                                          {
                                              "title": "Tech Debt",
                                              "value": 52
                                          }
                                      ]
                                  }
                              }
                          ],
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      },
                      {
                          "is_sub_org": false,
                          "compare_title": "Suborg4 - component",
                          "sub_org_count": 0,
                          "component_count": 0,
                          "value_in_millis": 168429000,
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      }
                  ],
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 3 - Component 1",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "value_in_millis": 168429000,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 3 - Component 2",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "value_in_millis": 168429000,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              }
          ],
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": true,
          "sub_org_id": "sub-org-9",
          "compare_title": "Sub org 9",
          "sub_org_count": 0,
          "component_count": 2,
          "value_in_millis": 168429000,
          "compare_reports": [
              {
                  "is_sub_org": true,
                  "sub_org_id": "sub-org-10",
                  "compare_title": "Suborg 10",
                  "sub_org_count": 1,
                  "component_count": 1,
                  "value_in_millis": 168429000,
                  "compare_reports": [
                      {
                          "is_sub_org": true,
                          "sub_org_id": "sub-org-11",
                          "compare_title": "Suborg 11",
                          "sub_org_count": 0,
                          "component_count": 1,
                          "value_in_millis": 168429000,
                          "compare_reports": [
                              {
                                  "is_sub_org": false,
                                  "compare_title": "Suborg11 - component",
                                  "sub_org_count": 0,
                                  "component_count": 0,
                                  "value_in_millis": 168429000,
                                  "section": {
                                      "data": [
                                          {
                                              "title": "Bugs",
                                              "value": 30
                                          },
                                          {
                                              "title": "Feature",
                                              "value": 52
                                          },
                                          {
                                              "title": "Risk",
                                              "value": 52
                                          },
                                          {
                                              "title": "Tech Debt",
                                              "value": 52
                                          }
                                      ]
                                  }
                              }
                          ],
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      },
                      {
                          "is_sub_org": false,
                          "compare_title": "Suborg10 - component",
                          "sub_org_count": 0,
                          "component_count": 0,
                          "value_in_millis": 168429000,
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      }
                  ],
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 9 - Component 1",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "value_in_millis": 168429000,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 9 - Component 2",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "value_in_millis": 168429000,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              }
          ],
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 1",
          "sub_org_count": 0,
          "component_count": 0,
          "value_in_millis": 168429000,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 2",
          "sub_org_count": 0,
          "component_count": 0,
          "value_in_millis": 168429000,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 3",
          "sub_org_count": 0,
          "component_count": 0,
          "value_in_millis": 168429000,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 4",
          "sub_org_count": 0,
          "component_count": 0,
          "value_in_millis": 168429000,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      }
  ],
  "section": {
      "data": [
          {
              "title": "Bugs",
              "value": 30
          },
          {
              "title": "Feature",
              "value": 52
          },
          {
              "title": "Risk",
              "value": 52
          },
          {
              "title": "Tech Debt",
              "value": 52
          }
      ]
  }
}`

const CompareReportsCycleTimeMock2 = `{
  "is_sub_org": true,
  "sub_org_id": "sub-org-2",
  "compare_title": "Sub org 2",
  "sub_org_count": 1,
  "component_count": 5,
  "value_in_millis": 168429000,
  "compare_reports": [
      {
          "is_sub_org": true,
          "sub_org_id": "sub-org-6",
          "compare_title": "Sub org 6",
          "sub_org_count": 0,
          "component_count": 2,
          "value_in_millis": 168429000,
          "compare_reports": [
              {
                  "is_sub_org": true,
                  "sub_org_id": "sub-org-7",
                  "compare_title": "Suborg 7",
                  "sub_org_count": 1,
                  "component_count": 1,
                  "value_in_millis": 168429000,
                  "compare_reports": [
                      {
                          "is_sub_org": true,
                          "sub_org_id": "sub-org-8",
                          "compare_title": "Suborg 8",
                          "sub_org_count": 0,
                          "component_count": 1,
                          "value_in_millis": 168429000,
                          "compare_reports": [
                              {
                                  "is_sub_org": false,
                                  "compare_title": "Suborg8 - component",
                                  "sub_org_count": 0,
                                  "component_count": 0,
                                  "value_in_millis": 168429000,
                                  "section": {
                                      "data": [
                                          {
                                              "title": "Bugs",
                                              "value": 30
                                          },
                                          {
                                              "title": "Feature",
                                              "value": 52
                                          },
                                          {
                                              "title": "Risk",
                                              "value": 52
                                          },
                                          {
                                              "title": "Tech Debt",
                                              "value": 52
                                          }
                                      ]
                                  }
                              }
                          ],
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      },
                      {
                          "is_sub_org": false,
                          "compare_title": "Suborg7 - component",
                          "sub_org_count": 0,
                          "component_count": 0,
                          "value_in_millis": 168429000,
                          "section": {
                              "data": [
                                  {
                                      "title": "Bugs",
                                      "value": 30
                                  },
                                  {
                                      "title": "Feature",
                                      "value": 52
                                  },
                                  {
                                      "title": "Risk",
                                      "value": 52
                                  },
                                  {
                                      "title": "Tech Debt",
                                      "value": 52
                                  }
                              ]
                          }
                      }
                  ],
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 6 - Component 1",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "value_in_millis": 168429000,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              },
              {
                  "is_sub_org": false,
                  "compare_title": "Suborg 6 - Component 2",
                  "sub_org_count": 0,
                  "component_count": 0,
                  "value_in_millis": 168429000,
                  "section": {
                      "data": [
                          {
                              "title": "Bugs",
                              "value": 30
                          },
                          {
                              "title": "Feature",
                              "value": 52
                          },
                          {
                              "title": "Risk",
                              "value": 52
                          },
                          {
                              "title": "Tech Debt",
                              "value": 52
                          }
                      ]
                  }
              }
          ],
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 1",
          "sub_org_count": 0,
          "component_count": 0,
          "value_in_millis": 168429000,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 2",
          "sub_org_count": 0,
          "component_count": 0,
          "value_in_millis": 168429000,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 3",
          "sub_org_count": 0,
          "component_count": 0,
          "value_in_millis": 168429000,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      },
      {
          "is_sub_org": false,
          "compare_title": "Component 4",
          "sub_org_count": 0,
          "component_count": 0,
          "value_in_millis": 168429000,
          "section": {
              "data": [
                  {
                      "title": "Bugs",
                      "value": 30
                  },
                  {
                      "title": "Feature",
                      "value": 52
                  },
                  {
                      "title": "Risk",
                      "value": 52
                  },
                  {
                      "title": "Tech Debt",
                      "value": 52
                  }
              ]
          }
      }
  ],
  "section": {
      "data": [
          {
              "title": "Bugs",
              "value": 30
          },
          {
              "title": "Feature",
              "value": 52
          },
          {
              "title": "Risk",
              "value": 52
          },
          {
              "title": "Tech Debt",
              "value": 52
          }
      ]
  }
}`

const CompareReportsCycleTimeMock3 = `{
  "is_sub_org": false,
  "compare_title": "Component",
  "sub_org_count": 0,
  "component_count": 0,
  "value_in_millis": 168429000,
  "section": {
      "data": [
          {
              "title": "Bugs",
              "value": 30
          },
          {
              "title": "Feature",
              "value": 52
          },
          {
              "title": "Risk",
              "value": 52
          },
          {
              "title": "Tech Debt",
              "value": 52
          }
      ]
  }
}`

const GetCustomLayout = `{
    "size": 1,
       "sort": [
         {
           "timestamp": {
             "order": "desc"
           }
         }
       ],
    "query": {
       "bool": {
          "filter": [
            {
               "term": {
                 "org_id": "{{.orgId}}"
               }
            },
            {
               "term": {
                  "user_id" : "{{.userId}}"
               }
            },
            {
               "term": {
                 "dashboard_name": "{{.dashboardName}}"
               }
            }
          ]
       }
    }
  }`

const GetWorkflowRunsCount = `{
    "query": {
       "bool": {
          "filter": []
       }
    }
  }`

const DocsCountQuery = `{
    "query": {
        "bool": {
            "filter": [
                {
                    "term":
                    {
                        "org_id": "{{.orgId}}"
                    }
                }
            ]
        }
    }
}`

const TestSuitesOverviewQuery = `{
    "size": 0,
    "query":
    {
        "bool":
        {
            "filter":
            [
                {
                    "range":
                    {
                        "start_time":
                        {
                            "gte":"{{.startDate}}",
                            "lte": "{{.endDate}}",
                            "format": "yyyy-MM-dd HH:mm:ss",
                            "time_zone":  "{{.timeZone}}"
                        }
                    }
                },
                {
                    "term":
                    {
                        "org_id": "{{.orgId}}"
                    }
                },
                {
                    "term":
                    {
                        "default_branch": "true"
                    }
                }
            ]
        }
    },
    "aggs": {
    "testSuitesOverview": {
      "scripted_metric": {
        "params": {
          "timeZone": "{{.timeZone}}",
          "timeFormat": "{{.timeFormat}}"
        },
        "combine_script": "return state.dataMap;",
        "init_script": "state.dataMap = [:];",
        "map_script": "def map = state.dataMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.automation_id.value + '_' + doc.branch_id.value + '_' + doc.test_suite_name.value + '_' + doc.start_time.getValue().toEpochSecond() * 1000;def v = ['org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'component_name':doc.component_name.value, 'automation_id':doc.automation_id.value, 'automation_name':doc.automation_name.value, 'branch_id':doc.branch_id.value, 'branch_name':doc.branch_name.value, 'run_id':doc.run_id.value, 'test_suite_name':doc.test_suite_name.value, 'start_time':doc.start_time.value, 'total_cases':doc.total.value, 'duration':doc.duration.value, 'start_time_in_millis':doc['start_time'].getValue().toEpochSecond() * 1000, 'duration_in_millis':doc['duration'].getValue(), 'failed_cases_count':doc.failed.value, 'successful_cases_count':doc.passed.value, 'skipped_cases_count':doc.skipped.value];map.put(key, v);",
        "reduce_script": "def is24HourFormat = params.timeFormat == '24h'; def resultMap = new HashMap(); DateTimeFormatter formatterZoned = DateTimeFormatter.ofPattern( 'yyyy/MM/dd HH:mm').withZone(ZoneId.of(params.timeZone)); DateTimeFormatter twelveHourFormatter = DateTimeFormatter.ofPattern( 'yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); for (a in states) { if (a != null && a.size() > 0) { for (i in a.keySet()) { def record = a.get(i); def key = record.component_id + '_' + record.automation_id + '_' + record.branch_id + '_' + record.test_suite_name; if (resultMap.containsKey(key)) { def lastRecord = resultMap.get(key); if (lastRecord.start_time_in_millis < record .start_time_in_millis) { def runIdSet = lastRecord.get('run_id_set'); runIdSet.add(record.get('run_id')); record.put('run_id_set', runIdSet); record.put('workflow_runs', runIdSet.size()); def testSuiteRunCount = lastRecord.get('test_suite_runs'); record.put('test_suite_runs', ++testSuiteRunCount); double totalDurationInMillis = lastRecord.get( 'total_duration_in_millis'); double durationInMillis = record.get('duration_in_millis') < 10 ? 10.0 : record.get('duration_in_millis'); record.put('total_duration_in_millis', totalDurationInMillis + durationInMillis); double result = Math.round((totalDurationInMillis + durationInMillis) / testSuiteRunCount); record.put('average_duration', result); if (is24HourFormat) { record.put('start_time', formatterZoned.format(record.get( 'start_time'))); } else { record.put('start_time', twelveHourFormatter.format(record.get( 'start_time'))); } float failureRate = Math.round(record.get('total_cases') == 0 ? 0 : record.get('failed_cases_count') * 100 / record.get('total_cases')); record.put('failure_rate_for_last_run', failureRate + '%'); resultMap.put(key, record); } else { def runIdSet = lastRecord.get('run_id_set'); runIdSet.add(record.get('run_id')); lastRecord.put('workflow_runs', runIdSet.size()); def testSuiteRunCount = lastRecord.get('test_suite_runs'); lastRecord.put('test_suite_runs', ++testSuiteRunCount); double totalDurationInMillis = lastRecord.get( 'total_duration_in_millis'); double durationInMillis = record.get('duration_in_millis') < 10 ? 10.0 : record.get('duration_in_millis'); lastRecord.put('total_duration_in_millis', totalDurationInMillis + durationInMillis); double result = Math.round((totalDurationInMillis + durationInMillis) / testSuiteRunCount); lastRecord.put('average_duration', result); } } else { def runIdSet = new HashSet(); runIdSet.add(record.get('run_id')); record.put('run_id_set', runIdSet); record.put('workflow_runs', 1); record.put('test_suite_runs', 1); double durationInMillis = record.get('duration_in_millis') < 10 ? 10.0 : record.get('duration_in_millis'); record.put('total_duration_in_millis', durationInMillis); record.put('average_duration', durationInMillis); if (is24HourFormat) { record.put('start_time', formatterZoned.format(record.get( 'start_time'))); } else { record.put('start_time', twelveHourFormatter.format(record.get( 'start_time'))); } float failureRate = Math.round(record.get('total_cases') == 0 ? 0 : record.get('failed_cases_count') * 100 / record.get( 'total_cases')); record.put('failure_rate_for_last_run', failureRate + '%'); resultMap.put(key, record); } } } } return resultMap;"
      }
    }
  }
}`

const TestCasesOverviewQuery = `{
    "size": 0,
    "query":
    {
        "bool":
        {
            "filter":
            [
                {
                    "range":
                    {
                        "start_time":
                        {
                            "gte":"{{.startDate}}",
                            "lte": "{{.endDate}}",
                            "format": "yyyy-MM-dd HH:mm:ss",
                            "time_zone":  "{{.timeZone}}"
                        }
                    }
                },
                {
                    "term":
                    {
                        "org_id": "{{.orgId}}"
                    }
                }
            ]
        }
    },
    "aggs":
    {
        "testCasesOverview":
        {
            "scripted_metric":
            {
                "params":
                {
                    "timeZone": "{{.timeZone}}",
                    "timeFormat": "{{.timeFormat}}"
                },
                "combine_script": "return state.dataMap;",
                "init_script": "state.dataMap = [:];",
                "map_script": "def map=state.dataMap;def key=doc.org_id.value+'_'+doc.component_id.value+'_'+doc.automation_id.value+'_'+doc.branch_id.value+'_'+doc.test_suite_name.value+'_'+doc.test_case_name.value+'_'+doc.start_time.getValue().toEpochSecond()*1000;def v=['org_id':doc.org_id.value,'component_id':doc.component_id.value,'component_name':doc.component_name.value,'automation_id':doc.automation_id.value,'automation_name':doc.automation_name.value,'branch_id':doc.branch_id.value,'branch_name':doc.branch_name.value,'test_suite_name':doc.test_suite_name.value,'test_case_name':doc.test_case_name.value,'start_time':doc.start_time.value,'status':doc.status.value,'duration':doc.duration.value,'start_time_in_millis':doc['start_time'].getValue().toEpochSecond()*1000,'duration_in_millis':doc['duration'].getValue()];map.put(key,v);",
                "reduce_script": "def resultMap = new HashMap(); def is24HourFormat = params.timeFormat == '24h'; DateTimeFormatter formatterZoned = DateTimeFormatter.ofPattern( 'yyyy/MM/dd HH:mm').withZone(ZoneId.of(params.timeZone)); DateTimeFormatter twelveHourFormatter = DateTimeFormatter.ofPattern( 'yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.component_id + '_' + record.automation_id + '_' + record.branch_id + '_' + record.test_suite_name + '_' + record .test_case_name; if (resultMap.containsKey(key)) { def lastRecord = resultMap.get(key); if (lastRecord.start_time_in_millis < record .start_time_in_millis) { resultMap.put(key, record); } record.put('failure_count', lastRecord.get('failure_count')); record.put('success_count', lastRecord.get('success_count')); record.put('skipped_count', lastRecord.get('skipped_count')); record.put('total_exec_count', lastRecord.get( 'total_exec_count')); record.put('failure_rate', lastRecord.get('failure_rate')); def runsCount = lastRecord.get('runs'); runsCount = runsCount + 1; record.put('runs', runsCount); double totalDurationInMillis = lastRecord.get( 'total_duration_in_millis'); double durationInMillis = record.get('duration_in_millis'); record.put('total_duration_in_millis', totalDurationInMillis + durationInMillis); double result = (totalDurationInMillis + durationInMillis) / runsCount; result = Math.round(result * 100) / 100.0; record.put('average_duration', result); def successCount = record.get('success_count'); def failureCount = record.get('failure_count'); def skippedCount = record.get('skipped_count'); def totalExecCount = record.get('total_exec_count'); if (record.get('status') == 'FAILED') { failureCount = failureCount + 1; record.put('failure_count', failureCount); } else if (record.get('status') == 'PASSED') { successCount = successCount + 1; record.put('success_count', successCount); } else { skippedCount = skippedCount + 1; record.put('skipped_count', skippedCount); } totalExecCount = totalExecCount + 1; record.put('total_exec_count', totalExecCount); double failureRate = Math.round((failureCount * 100) / totalExecCount); record.put('failure_rate', failureRate + '%'); if (is24HourFormat) { record.put('start_time', formatterZoned.format(record.get( 'start_time'))); } else { record.put('start_time', twelveHourFormatter.format(record.get( 'start_time'))); } resultMap.put(key, record); } else { record.put('runs', 1); record.put('total_duration_in_millis', record.get( 'duration_in_millis')); record.put('average_duration', record.get( 'duration_in_millis')); if (record.get('status') == 'FAILED') { record.put('failure_count', 1); record.put('success_count', 0); record.put('skipped_count', 0); record.put('failure_rate', '100.0%'); } else if (record.get('status') == 'PASSED') { record.put('success_count', 1); record.put('failure_count', 0); record.put('skipped_count', 0); record.put('failure_rate', '0.0%'); } else { record.put('success_count', 0); record.put('failure_count', 0); record.put('skipped_count', 1); record.put('failure_rate', '0.0%'); } record.put('total_exec_count', 1); if (is24HourFormat) { record.put('start_time', formatterZoned.format(record.get( 'start_time'))); } else { record.put('start_time', twelveHourFormatter.format(record.get( 'start_time'))); } resultMap.put(key, record); } } } } return resultMap;"
            }
        }
    }
}`

const TestOverviewComponentsViewQuery = `{
    "size": 0,
    "query": {
        "bool": {
            "filter": [
                {
                    "range": {
                        "run_start_time": {
                            "gte": "{{.startDate}}",
                            "lte": "{{.endDate}}",
                            "format": "yyyy-MM-dd HH:mm:ss",
                            "time_zone": "{{.timeZone}}"
                        }
                    }
                },
                {
                    "term": {
                        "org_id": "{{.orgId}}"
                    }
                },
                {
                    "term": {
                        "default_branch": true
                    }
                }
            ]
        }
    },
    "aggs": {
        "workflow_buckets": {
            "terms": {
                "size": 65000,
                "script": {
                    "source": "doc['component_id'].value + '_' + doc['automation_id'].value"
                }
            },
            "aggs": {
                "failure_count": {
                    "sum": {
                        "field": "failed"
                    }
                },
                "success_count": {
                    "sum": {
                        "field": "passed"
                    }
                },
                "skipped_count": {
                    "sum": {
                        "field": "skipped"
                    }
                },
                "total_test_cases_runs": {
                    "sum": {
                        "field": "total"
                    }
                },
                "total_duration": {
                    "sum": {
                        "script": "doc['duration'].value == 0? 10: doc['duration'].value"
                    }
                },
                "failure_rate": {
                    "bucket_script": {
                        "buckets_path": {
                            "failures": "failure_count",
                            "total": "total_test_cases_runs"
                        },
                        "script": "params.failures*100/params.total"
                    }
                },
                "avg_run_time": {
                    "bucket_script": {
                        "buckets_path": {
                            "total_duration": "total_duration",
                            "total": "total_test_cases_runs"
                        },
                        "script": "params.total_duration/params.total"
                    }
                },
                "sorted_buckets": {
                    "bucket_sort": {
                        "sort": [
                            {
                                "failure_rate": {
                                    "order": "desc"
                                }
                            }
                        ]
                    }
                },
                "latest_doc": {
                    "top_hits": {
                        "sort": [
                            {
                                "run_start_time": {
                                    "order": "desc"
                                }
                            }
                        ],
                        "size": 1,
                        "script_fields": {
                            "run_start_time_in_millis": {
                                "script": {
                                    "source": "doc['run_start_time'].value.getMillis()"
                                }
                            },
                            "zoned_run_start_time": {
                                "script": {
                                    "params": {
                                        "timeZone": "{{.timeZone}}",
                                        "timeFormat": "{{.timeFormat}}"
                                    },
                                    "source": "def is24HourFormat = params.timeFormat == '24h'; DateTimeFormatter formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm') .withZone(ZoneId.of(params.timeZone)); DateTimeFormatter twelveHourFormatter = DateTimeFormatter.ofPattern( 'yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); if (is24HourFormat) { return formatter.format(doc['run_start_time'].value); } else { return twelveHourFormatter.format(doc['run_start_time'].value); }"
                                }
                            }
                        },
                        "_source": [
                            "total",
                            "run_id",
                            "branch_id",
                            "branch_name",
                            "automation_name",
                            "automation_id",
                            "run_start_time",
                            "component_name"
                        ]
                    }
                },
                "test_suite_buckets": {
                    "terms": {
                        "field": "test_suite_name",
                        "size": 10000
                    },
                    "aggs": {
                        "test_cases_count": {
                            "scripted_metric": {
                                "init_script": "state.latest = null;",
                                "map_script": "if (state.latest == null || doc['timestamp'].value.getMillis() > state.latest.timestamp){state.latest = ['timestamp':doc['timestamp'].value.getMillis(), 'test_cases_count':doc.total.value];}",
                                "combine_script": "return state.latest;",
                                "reduce_script": "long latest_timestamp = 0;long latest_test_cases_val = 0;for (s in states){if (s != null && s.timestamp > latest_timestamp){latest_timestamp = s.timestamp;latest_test_cases_val = s.test_cases_count;}}return latest_test_cases_val;"
                            }
                        }
                    }
                },
                "total_test_cases_count": {
                    "sum_bucket": {
                        "buckets_path": "test_suite_buckets>test_cases_count.value"
                    }
                }
            }
        }
    }
}`

type AutomationRunsCount struct {
	Aggregations struct {
		ComponentActivity struct {
			Value struct {
				Runs int `json:"runs"`
			} `json:"value"`
		} `json:"component_activity"`
	} `json:"aggregations"`
}

type TestSuitesCount struct {
	Aggregations struct {
		ComponentActivity struct {
			Value struct {
				Runs int `json:"runs"`
			} `json:"value"`
		} `json:"component_activity"`
	} `json:"aggregations"`
}

type TestAutomationRunChart struct {
	Data []struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
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
}

type AutomationRuns struct {
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

type TestComponentResponse struct {
	Aggregations struct {
		DistinctComponent struct {
			Value []string `json:"value"`
		} `json:"distinct_component"`
	} `json:"aggregations"`
}

type TestComponentDrillDownResponse struct {
	Aggregations struct {
		ComponentActivity struct {
			Value map[string]struct {
				AutomationID  string   `json:"automation_id"`
				ComponentID   string   `json:"component_id"`
				ComponentName string   `json:"component_name"`
				OrgID         string   `json:"org_id"`
				TestSuitesSet []string `json:"test_suites_set"`
				TestSuiteName string   `json:"test_suite_name"`
			} `json:"value"`
		} `json:"component_activity"`
	} `json:"aggregations"`
}

type TestAutomationDrillDownResponse struct {
	Aggregations struct {
		ComponentActivity struct {
			Value map[string]struct {
				AutomationID  string   `json:"automation_id"`
				ComponentID   string   `json:"component_id"`
				BranchID      string   `json:"branch_id"`
				OrgID         string   `json:"org_id"`
				TestSuitesSet []string `json:"test_suites_set"`
				TestSuiteName string   `json:"test_suite_name"`
			} `json:"value"`
		} `json:"component_activity"`
	} `json:"aggregations"`
}

type TestAutomationRunsResponse struct {
	Aggregations struct {
		ComponentActivity struct {
			Value map[string]struct {
				AutomationID string   `json:"automation_id"`
				ComponentID  string   `json:"component_id"`
				OrgID        string   `json:"org_id"`
				RunIds       []string `json:"run_ids"`
			} `json:"value"`
		} `json:"component_activity"`
	} `json:"aggregations"`
}

type TestAutomationRunDrillDownResponse struct {
	Aggregations struct {
		TestWorkflowDrilldown struct {
			Value map[string]struct {
				ComponentID    string `json:"component_id"`
				RunID          string `json:"run_id"`
				ComponentName  string `json:"component_name"`
				AutomationID   string `json:"automation_id"`
				BranchID       string `json:"branch_id"`
				OrgID          string `json:"org_id"`
				BranchName     string `json:"branch_name"`
				RunNumber      string `json:"run_number"`
				RunStatus      string `json:"run_status"`
				AutomationName string `json:"automation_name"`
				Runs           int    `json:"runs"`
				TestSuiteName  string `json:"test_suite_name"`
				RunStartTime   int    `json:"run_start_time"`
				Status         string `json:"status"`
			} `json:"value"`
		} `json:"test_workflow_drilldown"`
	} `json:"aggregations"`
}

type OpenFindingsBySeverity struct {
	Aggregations struct {
		OpenFindingsBySeverity struct {
			Buckets []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				TrackingID struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
						ToolID   struct {
							Buckets []struct {
								Key      string `json:"key"`
								DocCount int    `json:"doc_count"`
							} `json:"buckets"`
						} `json:"tool_id"`
					} `json:"buckets"`
				} `json:"tracking_id"`
			} `json:"buckets"`
		} `json:"open_findings_by_severity"`
	} `json:"aggregations"`
}

type SlaBreachedByAssetType struct {
	Aggregations struct {
		RemediationKey struct {
			Buckets []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				TrackingID struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"tracking_id"`
			} `json:"buckets"`
		} `json:"remediation_key"`
	} `json:"aggregations"`
}

type OpenFindingsBySecurityTools struct {
	Aggregations struct {
		ToolName struct {
			Buckets []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				TrackingID struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
						ToolID   struct {
							Buckets []struct {
								Key      string `json:"key"`
								DocCount int    `json:"doc_count"`
							} `json:"buckets"`
						} `json:"tool_id"`
					} `json:"buckets"`
				} `json:"tracking_id"`
			} `json:"buckets"`
		} `json:"tool_name"`
	} `json:"aggregations"`
}

type OpenFindingsDistributionByCategory struct {
	Aggregations struct {
		Category struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
				Severity struct {
					Buckets []struct {
						Key        string `json:"key"`
						DocCount   int    `json:"doc_count"`
						TrackingID struct {
							Buckets []struct {
								Key      string `json:"key"`
								DocCount int    `json:"doc_count"`
							} `json:"buckets"`
						} `json:"tracking_id"`
					} `json:"buckets"`
				} `json:"severity"`
			} `json:"buckets"`
		} `json:"category"`
	} `json:"aggregations"`
}

type SeverityData struct {
	Title string `json:"title"`
	Value int    `json:"value"`
}

type SeverityDistribution struct {
	ColorScheme []struct {
		Color0 string `json:"color0"`
		Color1 string `json:"color1"`
	} `json:"colorScheme"`
	Data      []SeverityData `json:"data"`
	DrillDown DrillDown      `json:"drillDown"`
}

type OpenFindingsByCategoryData struct {
	Total                int                  `json:"total"`
	SeverityDistribution SeverityDistribution `json:"severityDistribution"`
	CategoryName         string               `json:"categoryName"`
	DrillDown            DrillDown            `json:"drillDown"`
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

type OpenFindingsDistributionBySecurityTool struct {
	Aggregations struct {
		ToolID struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
				Severity struct {
					Buckets []struct {
						Key        string `json:"key"`
						DocCount   int    `json:"doc_count"`
						TrackingID struct {
							Buckets []struct {
								Key             string `json:"key"`
								DocCount        int    `json:"doc_count"`
								ToolDisplayName struct {
									Buckets []struct {
										Key      string `json:"key"`
										DocCount int    `json:"doc_count"`
									} `json:"buckets"`
								} `json:"tool_display_name"`
							} `json:"buckets"`
						} `json:"tracking_id"`
					} `json:"buckets"`
				} `json:"severity"`
			} `json:"buckets"`
		} `json:"tool_id"`
	} `json:"aggregations"`
}

type FindingsIdentifiedSince struct {
	Aggregations struct {
		RemediationStatus struct {
			Buckets []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				TrackingID struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"tracking_id"`
			} `json:"buckets"`
		} `json:"remediation_status"`
	} `json:"aggregations"`
}

type RiskAcceptedFalsePositiveFindings struct {
	Aggregations struct {
		RiskAcceptedCounts struct {
			DocCount                   int `json:"doc_count"`
			RA_NOT_EXPIRING_IN_30_DAYS struct {
				DocCount int `json:"doc_count"`
			} `json:"RA_NOT_EXPIRING_IN_30_DAYS"`
			RA_EXPIRING_IN_30_DAYS struct {
				DocCount int `json:"doc_count"`
			} `json:"RA_EXPIRING_IN_30_DAYS"`
			RA_EXPIRED struct {
				DocCount int `json:"doc_count"`
			} `json:"RA_EXPIRED"`
		} `json:"risk_accepted_counts"`

		RemediationStatus struct {
			Buckets []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				TrackingID struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"tracking_id"`
			} `json:"buckets"`
		} `json:"remediation_status"`
	} `json:"aggregations"`
}

type SLABreachesBySeverity struct {
	Aggregations struct {
		SLABreachesBySeverity struct {
			Buckets []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				TrackingID struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
						ToolID   struct {
							Buckets []struct {
								Key      string `json:"key"`
								DocCount int    `json:"doc_count"`
							} `json:"buckets"`
						} `json:"tool_id"`
					} `json:"buckets"`
				} `json:"tracking_id"`
			} `json:"buckets"`
		} `json:"sla_breaches_by_severity"`
	} `json:"aggregations"`
}

type OpenFindingsBySLAStatus struct {
	Aggregations struct {
		OpenFindingsBySLAStatus struct {
			Buckets struct {
				NonSLABreached struct {
					DocCount   int `json:"doc_count"`
					TrackingID struct {
						Buckets []struct {
							Key      string `json:"key"`
							DocCount int    `json:"doc_count"`
						} `json:"buckets"`
					} `json:"tracking_id"`
				} `json:"non_sla_breached"`
				SLABreached struct {
					DocCount   int `json:"doc_count"`
					TrackingID struct {
						Buckets []struct {
							Key      string `json:"key"`
							DocCount int    `json:"doc_count"`
						} `json:"buckets"`
					} `json:"tracking_id"`
				} `json:"sla_breached"`
			} `json:"buckets"`
		} `json:"open_findings_by_sla_status"`
	} `json:"aggregations"`
}

type OpenFindingsByReviewStatus struct {
	Aggregations struct {
		TriageStatus struct {
			Buckets []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				TrackingID struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"tracking_id"`
			} `json:"buckets"`
		} `json:"triage_status"`
	} `json:"aggregations"`
}

type FindingsRemediationTrend struct {
	Aggregations struct {
		FindingsRemediationTrend struct {
			Value map[string]struct {
				Open            int `json:"Open"`
				BreachedSLA     int `json:"BreachedSLA"`
				ClosedWithinSLA int `json:"ClosedWithinSLA"`
			} `json:"value"`
		} `json:"findings_remediation_trend"`
	} `json:"aggregations"`
}

type FindingsRemediationTrendAppSec struct {
	Aggregations struct {
		FindingsRemediationTrend struct {
			Value map[string]struct {
				Open            int `json:"Open"`
				BreachedSLA     int `json:"BreachedSLA"`
				ClosedWithinSLA int `json:"ClosedWithinSLA"`
				New             int `json:"New"`
			} `json:"value"`
		} `json:"findings_remediation_trend"`
	} `json:"aggregations"`
}
