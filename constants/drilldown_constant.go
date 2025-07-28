package constants

const (
	COMPONENT_ACTIVITY      = "component_activity"
	AUTOMATION_ACTIVITY     = "automation_activity"
	AUTOMATION_RUN_ACTIVITY = "automation_run_activity"
	COMPONENT_ID            = "componentId"
	COMPONENT_NAME          = "componentName"
	REPOSITORY_URL          = "repositoryUrl"
	LAST_ACTIVE             = "lastActive"
	STATUS                  = "status"
	WORKFLOW_SOURCE         = "source"
	BRANCH                  = "branch"
	SOURCE_BRANCH           = "sourceBranch"
	TARGET_BRANCH           = "targetBranch"
	WORKFLOW                = "workflow"
	BUILD                   = "build"
	BUILDS                  = "builds"
	BUILD_ID                = "buildID"
	DEPLOYMENT_ID           = "deploymentID"
	START_TIME              = "startTime"
	END_TIME                = "endTime"
	DEPLOYMENTS             = "deployments"
	RUN_START_TIME          = "runStartTime"
	WORKFLOWS               = "workflows"
	ENVIRONMENT             = "environment"
	DURATION                = "duration"
	URL                     = "url"
	COMMITS                 = "commits"
	PULL_REQUESTS           = "pullrequests"
	COMMIT_ID               = "commitID"
	REPO                    = "repo"
	AUTHOR                  = "author"
	COMMIT_TIME             = "commitTime"
	PR_ID                   = "prID"
	CREATED_ON              = "createdOn"
	CREATED                 = "created"
	RUN_ID                  = "runID"
	RUN_ID_KEY              = "runId"
	COMMIT_DESCRIPTION      = "commitDescription"
	AUTOMATION_ID           = "automationId"
	BRANCH_ID               = "branchId"
	JOB_ID                  = "jobId"
	STEP_ID                 = "stepId"
	LEAD_TIME               = "leadTime"
	DEPLOYED_TIME           = "deployedTime"
	SUCCESS                 = "success"
	FAILURE                 = "failure"
	FAILURE_RATE            = "failureRate"
	FAILED_ON               = "failedOn"
	RECOVERED_ON            = "recoveredOn"
	RECOVERY_DURATION       = "recoveryDurationInMillis"
	FAILED_RUN              = "failedRun"
	RECOVERED_RUN           = "recoveredRun"

	COMPONENTS_DRILLDOWN              = "Components"
	WORKFLOWS_DRILLDOWN               = "Workflows"
	WORKFLOW_RUNS                     = "workflowRuns"
	TEST_INSIGHTS_WORKFLOWS_DRILLDOWN = "Workflows"

	SECURITY_INDEX          = "scan_results"
	OPEN_VUL_FILTER         = "severity"
	VUL_BY_SCANTYPE_COLUMNS = "scanner_type"
	VUL_OVERVIEW_COLUMN     = "bug_status"
	DRILLDOWNS              = "drilldowns"
	DRILLDOWN               = "drilldown"
	FLOW_METRICS_INDEX      = "flow_metrics"
	FLOW_METRICS_FILTER     = "flowItemType"
	SECURITY_ISSUE_COUNT    = "severityCounts"
	DISABLED                = "Disabled"
	NEW_VERSION_AVAILABLE   = "New version available"
	UP_TO_DATE              = "Up to date"
	ENABLE                  = "Enabled"
	START_TIME_CONVERTED    = "startTimeConverted"
	START_TIME_KEY          = "start_time"
	TEST_SUITE_INDEX        = "cb_test_suites"
	TEST_CASES_INDEX        = "cb_test_cases"
	RUN_STATUS              = "runStatus"
	CI_INSIGHTS_INDEX       = "cb_ci_tool_insight"
	DORA_METRICS_INDEX      = "deploy_data"
)

// Drilldown definition map keys
const (
	OPEN_VULNERABILITIES                             = "openVulnerabilities"
	NESTED_DRILLDOWN_VIEW_LOCATION                   = "nested-drilldown-view-location"
	OPEN_VULNERABILITIES_VIEW_LOCATION               = "open-vulnerabilities-view-location"
	VULNERABILITIES_SECURITY_SCAN_TYPE               = "vulnerabilitiesSecurityScanType"
	VULNERABILITIES_SECURITY_SCAN_TYPE_VIEW_LOCATION = "vulnerabilities-security-scan-type-view-location"
	VULNERABILITIES_OVERVIEW                         = "vulnerabilitiesOverview"
	VULNERABILITIES_OVERVIEW_VIEW_LOCATION           = "vulnerabilities-overview-view-location"
	MTTR_FOR_VULNERABILITIES                         = "mttrForVulnerabilities"
	CWE_TOP_25_VULNERABILITIES                       = "cweTop25Vulnerabilities"
	CWE_TOP_25_VULNERABILITIES_VIEW_LOCATION         = "cwe-top25-vulnerabilities-view-location"
	SECURITY_SLA_STATUS_OVERVIEW_OPEN                = "security-SLA-status-overview-open"
	SECURITY_SLA_STATUS_OVERVIEW_CLOSED              = "security-SLA-status-overview-closed"
	FLOW_METRICS_VELOCITY                            = "flowMetrics-velocity"
	FLOW_METRICS_WORK_ITEM_DISTRIBUTION              = "flowMetrics-workItemDistribution"
	FLOW_METRICS_CYCLE_TIME                          = "flowMetrics-cycleTime"
	FLOW_METRICS_WORK_EFFICIENCY                     = "flowMetrics-workEfficiency"
	FLOW_METRICS_WORK_LOAD                           = "flowMetrics-workLoad"
	ACTIVE_DEVELOPERS                                = "activeDevelopers"
	ACTIVE_DEVELOPERS_COMMITS                        = "activeDevelopersCommits"
	TRIVY_LICENSE_OCCURENCE                          = "trivy-license-occurence"
	OPEN_ISSUES_DRILL_DOWN                           = "open-issues-drill-down"
	LATEST_TEST_RESULTS                              = "latest-test-results"
	TEST_OVERVIEW_TOTAL_TEST_CASES                   = "test-overview-total-tests-cases"
	TEST_OVERVIEW_TOTAL_RUNS                         = "test-overview-total-runs"
	TEST_OVERVIEW_VIEW_RUN_ACTIVITY                  = "test-overview-view-run-activity"
	OPEN_VULNERABILITIES_SUBROWS                     = "openVulnerabilitiesSubRows"
	VULNERABILITIES_OVERVIEW_SUBROWS                 = "vulnerabilitiesOverviewSubRows"
	MTTR_FOR_VULNERABILITIES_SUBROWS                 = "mttrForVulnerabilitiesSubRows"
	CWE_TOP_25_VULNERABILITIES_SUBROWS               = "cweTop25VulnerabilitiesSubRows"
	RUN_DETAILS_TEST_RESULTS                         = "run-details-test-results"
	RUN_DETAILS_TOTAL_TEST_CASES                     = "run-details-total-test-cases"
	RUN_DETAILS_TEST_CASE_LOG                        = "run-details-test-case-log"
	RUN_DETAILS_TEST_RESULTS_INDICATORS              = "run-details-test-results-indicators"
	VULNERABILITIES_SECURITY_SCAN_TYPE_SUBROWS       = "vulnerabilitiesSecurityScanTypeSubRows"
)

const TestAutomationDrilldownQuery = `{
    "_source": false,
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
                            "format": "yyyy-MM-dd HH:mm:ss",
                            "time_zone": "{{.timeZone}}"
                        }
                    }
                }
            ]
        }
    },
    "aggs": {
        "component_activity": {
            "scripted_metric": {
                "init_script": "state.data_map = [:];",
                "map_script": "def map=state.data_map;def key=doc.component_id.value+'_'+doc.automation_id.value+'_'+doc.branch_id.value+'_'+doc.test_suite_name.value;def v=['org_id':doc.org_id.value,'component_id':doc.component_id.value,'automation_id':doc.automation_id.value,'branch_id':doc.branch_id.value,'test_suite_name':doc.test_suite_name.value];map.put(key,v);",
                "combine_script": "return state.data_map;",
                "reduce_script": "def resultMap=new HashMap();for(a in states){if(a!=null){for(i in a.keySet()){def record=a.get(i);def key=record.component_id+'_'+record.automation_id;if(resultMap.containsKey(key)){def lastRecord=resultMap.get(key);def testSuitesSet=lastRecord.get('test_suites_set');testSuitesSet.add(record.get('test_suite_name'));resultMap.put(key,lastRecord);}else{def set=new HashSet();set.add(record.get('test_suite_name'));record.put('test_suites_set',set);resultMap.put(key,record);}}}}return resultMap;"
            }
        }
    }
}`

const TestComponentDrilldownQuery = `{
    "_source": false,
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
                            "format": "yyyy-MM-dd HH:mm:ss",
                            "time_zone": "{{.timeZone}}"
                        }
                    }
                }
            ]
        }
    },
    "aggs": {
        "component_activity": {
            "scripted_metric": {
                "init_script": "state.data_map = [:];",
                "map_script": "def map=state.data_map;def key=doc.component_id.value+'_'+doc.test_suite_name.value;def v=['org_id':doc.org_id.value,'component_id':doc.component_id.value,'component_name':doc.component_name.value,'automation_id':doc.automation_id.value,'run_id':doc.run_id.value,'test_suite_name':doc.test_suite_name.value];map.put(key,v);",
                "combine_script": "return state.data_map;",
                "reduce_script": "def resultMap=new HashMap();for(a in states){if(a!=null){for(i in a.keySet()){def record=a.get(i);def key=record.component_id;if(resultMap.containsKey(key)){def lastRecord=resultMap.get(key);def testSuitesSet=lastRecord.get('test_suites_set');testSuitesSet.add(record.get('test_suite_name'));resultMap.put(key,lastRecord);}else{def set=new HashSet();set.add(record.get('test_suite_name'));record.put('test_suites_set',set);resultMap.put(key,record);}}}}return resultMap;"
            }
        }
    }
}`

const ComponentDrilldownQuery = `{
	"_source": false,
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
		  {
			"term": {
				"org_id": "{{.orgId}}"
			  }
			}
		  ]
		}
	},
	"aggs": {
	  "component_activity": {
		"scripted_metric": {
		  "init_script": "state.data_map=[:];",
		  "map_script": "def map = state.data_map;def key = doc.component_id.value + '_' + doc.last_active_time.value;def v = ['last_active_time': doc.last_active_time.value, 'component_id':doc.component_id.value,'component_name':doc.component_name.value,'repo_url':doc.repo_url.value];map.put(key, v);",
		  "combine_script": "return state.data_map;",
		  "reduce_script": "def tmpMap = [: ], resultMap = new HashMap();for (response in states) {if (response != null) {for (key in response.keySet()) {def record = response.get(key);if (tmpMap.containsKey(key)) {def mapRecord = tmpMap.get(key);if (mapRecord.last_active_time.getMillis() > record.last_active_time.getMillis()) {tmpMap.put(key, mapRecord);}} else {tmpMap.put(key, record);}}}}for (key in tmpMap.keySet()) {def mapRecord = tmpMap.get(key);def component_id = mapRecord.component_id;if (resultMap.containsKey(component_id)) {def record = resultMap.get(component_id);if (mapRecord.last_active_time.getMillis() > record.last_active_time.getMillis()) {resultMap.put(component_id, mapRecord);}} else {resultMap.put(component_id, mapRecord);}}return resultMap;"
		}
	  }
	}
  }`

const AutomationDrilldownQuery = `{
	"_source": false,
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
		  {
			"term": {
				"org_id": "{{.orgId}}"
			  }
			}
		  ]
		}
	},
	"aggs": {
	  "automation_activity": {
		"scripted_metric": {
		  "init_script": "state.data_map=[:];",
		  "map_script": "def map = state.data_map;def key = doc.automation_id.value + '_' + doc.last_active_time.value;def v = ['last_active_time': doc.last_active_time.value, 'component_id': doc.component_id.value, 'component_name': doc.component_name.value, 'workflow_name': doc.workflow_name.value, 'branch_name': doc.branch_name.value,'branch_id': doc.branch_id.value, 'automation_id': doc.automation_id.value];map.put(key, v);",
		  "combine_script": "return state.data_map;",
		  "reduce_script": "def tmpMap = [: ], resultMap = new HashMap();for (response in states) {if (response != null) {for (key in response.keySet()) {def record = response.get(key);if (tmpMap.containsKey(key)) {def mapRecord = tmpMap.get(key);if (mapRecord.last_active_time.getMillis() > record.last_active_time.getMillis()) {tmpMap.put(key, mapRecord);}} else {tmpMap.put(key, record);}}}}for (key in tmpMap.keySet()) {def mapRecord = tmpMap.get(key);def automation_id = mapRecord.automation_id;if (resultMap.containsKey(automation_id)) {def record = resultMap.get(automation_id);if (mapRecord.last_active_time.getMillis() > record.last_active_time.getMillis()) {resultMap.put(automation_id, mapRecord);}} else {resultMap.put(automation_id, mapRecord);}}return resultMap;"
		}
	  }
	}
  }`

const AutomationDrilldownQueryForBranch = `{
	"_source": false,
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
				}
		  ]
		}
	},
	"aggs": {
	  "automation_activity": {
		"scripted_metric": {
		  "init_script": "state.data_map=[:];",
		  "map_script": "def map = state.data_map;def key = doc.automation_id.value + '_' + doc.last_active_time.value;def v = ['last_active_time': doc.last_active_time.value, 'component_id': doc.component_id.value, 'component_name': doc.component_name.value, 'workflow_name': doc.workflow_name.value, 'branch_name': doc.branch_name.value,'branch_id': doc.branch_id.value, 'automation_id': doc.automation_id.value];map.put(key, v);",
		  "combine_script": "return state.data_map;",
		  "reduce_script": "def tmpMap = [: ], resultMap = new HashMap();for (response in states) {if (response != null) {for (key in response.keySet()) {def record = response.get(key);if (tmpMap.containsKey(key)) {def mapRecord = tmpMap.get(key);if (mapRecord.last_active_time.getMillis() > record.last_active_time.getMillis()) {tmpMap.put(key, mapRecord);}} else {tmpMap.put(key, record);}}}}for (key in tmpMap.keySet()) {def mapRecord = tmpMap.get(key);def automation_id = mapRecord.automation_id;if (resultMap.containsKey(automation_id)) {def record = resultMap.get(automation_id);if (mapRecord.last_active_time.getMillis() > record.last_active_time.getMillis()) {resultMap.put(automation_id, mapRecord);}} else {resultMap.put(automation_id, mapRecord);}}return resultMap;"
		}
	  }
	}
  }`

const AutomationRunDrilldownQuery = `{
	"_source": false,
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
				"status_timestamp": {
				  "gte": "{{.startDate}}",
				  "lte": "{{.endDate}}",
				  "format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis",
				  "time_zone":"{{.timeZone}}"
				}
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
		"automation_run_activity": {
		  "scripted_metric": {
			"init_script": "state.data_map=[:];",
			"map_script": "def map = state.data_map;def key = doc.org_id+'_'+doc.automation_id.value + '_' + doc.run_id.value + '_' + doc.status.value;def v = ['status_timestamp': doc.status_timestamp.value, 'component_id': doc.component_id.value,'component_name': doc.component_name.value, 'automation_id': doc.automation_id.value, 'run_id': doc.run_id.value, 'status': doc.status.value, 'duration': doc.duration.value, 'run_number':doc['workflow_info.run_number'].value];map.put(key, v);",
			"combine_script": "return state.data_map;",
			"reduce_script": "def tmpMap = [:], resultMap = new HashMap();for (response in states){if (response != null && response.size() > 0){for (key in response.keySet()){def record = response.get(key);def uniqueRunKey = record.automation_id + '_' + record.run_id;if (tmpMap.containsKey(uniqueRunKey)){def mapRecord = tmpMap.get(uniqueRunKey);if (mapRecord.status_timestamp.getMillis() > record.status_timestamp.getMillis()){if (mapRecord.status == 'SUCCEEDED'){mapRecord.status = 'Success';} else if (mapRecord.status == 'FAILED' || mapRecord.status == 'TIMED_OUT' || mapRecord.status == 'ABORTED'){mapRecord.status = 'Failure';}tmpMap.put(uniqueRunKey, mapRecord);}} else{if (record.status == 'SUCCEEDED'){record.status = 'Success';} else if (record.status == 'FAILED' || record.status == 'TIMED_OUT' || record.status == 'ABORTED'){record.status = 'Failure';}tmpMap.put(uniqueRunKey, record);}}}}for (key in tmpMap.keySet()){def mapRecord = tmpMap.get(key);def run_id = mapRecord.run_id;def automation_id = mapRecord.automation_id;def status = mapRecord.status;if (resultMap.containsKey(automation_id)){def runList = resultMap.get(automation_id);runList.add(mapRecord);} else{def runList = new ArrayList();runList.add(mapRecord);resultMap.put(automation_id, runList);}}return resultMap;"
		  }
		}
	}
}`

const TestSuiteQuery = `{
  "_source": false,
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
        }
      ]
    }
  },
 "aggs": {
    "test_workflow_drilldown": {
      "scripted_metric": {
        "init_script": "state.data_map=[:];",
        "map_script": "def map=state.data_map;def key=doc.component_id.value+'_'+doc.automation_id.value+'_'+doc.run_id.value;def v=['org_id':doc.org_id.value,'component_id':doc.component_id.value,'component_name':doc.component_name.value,'automation_id':doc.automation_id.value,'automation_name':doc.automation_name.value,'branch_id':doc.branch_id.value,'branch_name':doc.branch_name.value,'run_id':doc.run_id.value,'test_suite_name':doc.test_suite_name.value,'run_number':doc.run_number.value,'run_status':doc.run_status.value,'run_start_time':doc.run_start_time.value.toInstant().toEpochMilli(),'status':doc.status.value];map.put(key,v);",
        "combine_script": "return state.data_map;",
        "reduce_script": "def resultMap=new HashMap();for(a in states){if(a!=null){for(i in a.keySet()){def record=a.get(i);def key=record.component_id+'_'+record.automation_id+'_'+record.run_id;if(resultMap.containsKey(key)){def lastRecord=resultMap.get(key);def runsCount=lastRecord.get('runs');runsCount=runsCount+1;record.put('runs',runsCount);resultMap.put(key,record);}else{record.put('runs',1);resultMap.put(key,record);}}}}return resultMap;"
      }
    }
  }
}`

const TestAutomationRunDrilldownQuery = `{
	"_source": false,
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
				  "format": "yyyy-MM-dd HH:mm:ss",
				  "time_zone":"{{.timeZone}}"
				}
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
		"automation_run_activity": {
		  "scripted_metric": {
			"init_script": "state.data_map=[:];",
			"map_script": "def map=state.data_map;def key=doc.org_id+'_'+doc.automation_id.value+'_'+doc.run_id.value;def v=['status_timestamp':doc.status_timestamp.value,'component_id':doc.component_id.value,'component_name':doc.component_name.value,'automation_id':doc.automation_id.value,'run_id':doc.run_id.value,'status':doc.status.value,'duration':doc.duration.value,'run_number':doc['workflow_info.run_number'].value];map.put(key,v);",
			"combine_script": "return state.data_map;",
			"reduce_script": "def tmpMap = [: ], resultMap = new HashMap();for (response in states) {if (response != null) {for (key in response.keySet()) {def record = response.get(key);if (tmpMap.containsKey(key)) {def mapRecord = tmpMap.get(key);if (mapRecord.status_timestamp.getMillis() > record.status_timestamp.getMillis()) {if (mapRecord.status == 'SUCCEEDED') {mapRecord.status = 'Success';} else if(mapRecord.status == 'FAILED' || mapRecord.status == 'TIMED_OUT' || mapRecord.status == 'ABORTED') {mapRecord.status = 'Failure';}tmpMap.put(key, mapRecord);}} else {if (record.status == 'SUCCEEDED') {record.status = 'Success';} else if(record.status == 'FAILED' || record.status == 'TIMED_OUT' || record.status == 'ABORTED') {record.status = 'Failure';}tmpMap.put(key, record);}}}}for (key in tmpMap.keySet()) {def mapRecord = tmpMap.get(key);def run_id = mapRecord.run_id;def automation_id = mapRecord.automation_id;def status = mapRecord.status;if (resultMap.containsKey(automation_id)) {def runList = resultMap.get(automation_id);runList.add(mapRecord);} else {def runList = new ArrayList();runList.add(mapRecord);resultMap.put(automation_id, runList);}}return resultMap;"
		  }
		}
	}
}`

const SecurityAutomationRunDrilldownQuery = `{
	"_source": false,
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
					"job_id": ""
				}
			},
			{
				"term": {
					"step_id": ""
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
		"automation_run_activity": {
		  "scripted_metric": {
			"init_script": "state.data_map=[:];",
			"map_script": "def map=state.data_map;def key=doc.org_id+'_'+doc.automation_id.value+'_'+doc.run_id.value+'_'+doc.status.value;def v=['status_timestamp':doc.status_timestamp.value,'component_id':doc.component_id.value,'component_name':doc.component_name.value,'automation_id':doc.automation_id.value,'run_id':doc.run_id.value,'status':doc.status.value,'duration':doc.duration.value,'run_number':doc['workflow_info.run_number'].value];map.put(key,v);",
			"combine_script": "return state.data_map;",
			"reduce_script": "def tmpMap = [: ], resultMap = new HashMap();for (response in states) {if (response != null) {for (key in response.keySet()) {def record = response.get(key);if (tmpMap.containsKey(key)) {def mapRecord = tmpMap.get(key);if (mapRecord.status_timestamp.getMillis() > record.status_timestamp.getMillis()) {if (mapRecord.status == 'SUCCEEDED') {mapRecord.status = 'Success';} else if(mapRecord.status == 'FAILED' || mapRecord.status == 'TIMED_OUT' || mapRecord.status == 'ABORTED') {mapRecord.status = 'Failure';}tmpMap.put(key, mapRecord);}} else {if (record.status == 'SUCCEEDED') {record.status = 'Success';} else if(record.status == 'FAILED' || record.status == 'TIMED_OUT' || record.status == 'ABORTED') {record.status = 'Failure';}tmpMap.put(key, record);}}}}for (key in tmpMap.keySet()) {def mapRecord = tmpMap.get(key);def run_id = mapRecord.run_id;def automation_id = mapRecord.automation_id;def status = mapRecord.status;if (resultMap.containsKey(automation_id)) {def runList = resultMap.get(automation_id);runList.add(mapRecord);} else {def runList = new ArrayList();runList.add(mapRecord);resultMap.put(automation_id, runList);}}return resultMap;"
		  }
		}
	}
}`

const CommitsDrilldownQuery = `{
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
		  {
			"range": {
			  "commit_timestamp": {
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
			"bool": {
			  "must_not": {
				"term": {
				  "author": "github-actions[bot]"
				}
			  }
			}
		  },
		  {
			"bool": {
			  "must_not": {
				"prefix": {
				  "branch": {
					"value": "dependabot"
				  }
				}
			  }
			}
		  }
		]
	  }
	},
	"aggs": {
		"commits": {
		  "scripted_metric": {
			"params": {
				"timeZone": "{{.timeZone}}"
			},
			"init_script": "state.statusMap = [:];",
			"map_script": "def map=state.statusMap;def key=doc.repository_name.value+'_'+doc.commit_id.value+'_'+doc.branch.value;def repoUrl=\"\";if(doc.containsKey('repository_url')&&!doc.repository_url.empty){repoUrl=doc.repository_url.value;}def v=['branch':doc.branch.value,'component_id':doc.component_id.value,'commit_id':doc.commit_id.value,'org_id':doc.org_id.value,'commit_timestamp':doc.commit_timestamp.value,'component_name':doc.component_name.value,'repository_name':doc.repository_name.value,'repository_url':repoUrl,'author':doc.author.value];map.put(key,v);",
			"combine_script": "return state.statusMap;",
			"reduce_script": "def tmpMap=[:],resultMap=new HashMap(),commits=new ArrayList();for(response in states){if(response!=null){for(key in response.keySet()){tmpMap.put(key,response.get(key));}}}for(key in tmpMap.keySet()){def valueMap=tmpMap.get(key);def rd=valueMap.commit_timestamp;valueMap.commit_timestamp_UTC=rd;valueMap.commit_timestamp=rd.withZoneSameInstant(ZoneId.of(params.timeZone));commits.add(valueMap);}return commits;"
		  }
		}
	}
  }`

const ActiveDeveloperDrillDownQuery = `{
	"_source": false,
  "size": 0,
  "query": {
	  "bool": {
		"filter": [
			{
			  "range": {
				"commit_timestamp": {
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
			  "bool": {
				"must_not": {
				  "term": {
					"author": "github-actions[bot]"
				  }
				}
			  }
			},
			{
			  "bool": {
				"must_not": {
				  "prefix": {
					"branch": {
					  "value": "dependabot"
					}
				  }
				}
			  }
			}
		  ]
	  }
  },
   "aggs": {
		"drilldowns": {
		  "scripted_metric": {
			"init_script": "state.statusMap = [:];",
			"map_script": "def map = state.statusMap;def key = doc.repository_name.value + '_' + doc.commit_id.value + '_' + doc.branch.value + '_' + doc.author.value;def v = ['branch':doc.branch.value, 'component_id':doc.component_id.value, 'commit_id':doc.commit_id.value, 'org_id':doc.org_id.value, 'commit_timestamp':doc.commit_timestamp.value, 'component_name':doc.component_name.value, 'repository_name':doc.repository_name.value, 'author':doc.author.value];map.put(key, v);",
			"combine_script": "return state.statusMap;",
			"reduce_script": "def tmpMap = [:], resultMap = new HashMap(), commits = new ArrayList();def statusMap = new HashMap();def countMap = new HashMap();def resultList = new ArrayList();for (response in states){if (response != null){for (key in response.keySet()){statusMap.put(key, response.get(key));}}}if (statusMap.size() > 0){for (uniqueKey in statusMap.keySet()){def item = statusMap.get(uniqueKey);def map = [:];def countAuthor = item.author;if (countMap.containsKey(countAuthor)){def count = countMap.get(countAuthor);countMap.put(countAuthor, count + 1);} else{countMap.put(countAuthor, 1);}def reportInfoMap = new HashMap();reportInfoMap.put('author', item.author);def drillDownInfoMap = new HashMap();drillDownInfoMap.put('reportId', 'activeDevelopersCommits');drillDownInfoMap.put('reportTitle', 'Commits by '+item.author);drillDownInfoMap.put('reportInfo', reportInfoMap);map.put('drillDown', drillDownInfoMap);map.put('author', item.author);map.put('totalCommits', countMap.get(item.author));resultMap.put(item.author,map)}}for (x in resultMap.keySet()){resultList.add(resultMap.get(x))}return resultList;"
		  }
		}
	  }
	}`

const NestedActiveDeveloperCommitsInfoQuery = `{
		"_source": false,
	  "size": 0,
	  "query": {
		"bool": {
		  "filter": [
				  {
					"range": {
					  "commit_timestamp": {
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
					  "author": "{{.author}}"
					}
				  },
				  {
					"bool": {
					  "must_not": {
						"term": {
						  "author": "github-actions[bot]"
						}
					  }
					}
				  },
				  {
					"bool": {
					  "must_not": {
						"prefix": {
						  "branch": {
							"value": "dependabot"
						  }
						}
					  }
					}
				  }
				]
		}
	  },
	   "aggs": {
			  "drilldowns": {
				"scripted_metric": {
					"params": {
						"timeZone": "{{.timeZone}}",
						"timeFormat":"{{.timeFormat}}"
					},
				  "init_script": "state.statusMap = [:];",
				  "map_script": "def map=state.statusMap;def key=doc.repository_name.value+'_'+doc.commit_id.value+'_'+doc.branch.value+'_'+doc.author.value;def repoUrl=\"\";if(doc.containsKey('repository_url')&&!doc.repository_url.empty){repoUrl=doc.repository_url.value;}def v=['branch':doc.branch.value,'component_id':doc.component_id.value,'commit_id':doc.commit_id.value,'org_id':doc.org_id.value,'commit_timestamp':doc.commit_timestamp.value,'component_name':doc.component_name.value,'repository_name':doc.repository_name.value,'repository_url':repoUrl,'author':doc.author.value];map.put(key,v);",
				  "combine_script": "return state.statusMap;",
				  "reduce_script": "def is24HourFormat=params.timeFormat=='24h';DateTimeFormatter formatter=DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone));DateTimeFormatter twelveHourFormatter=DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone));def tmpMap=[:],resultMap=new HashMap(),commits=new ArrayList();def statusMap=new HashMap();def countMap=new HashMap();def resultList=new ArrayList();for(response in states){if(response!=null){for(key in response.keySet()){statusMap.put(key,response.get(key));}}}if(statusMap.size()>0){for(uniqueKey in statusMap.keySet()){def item=statusMap.get(uniqueKey);def map=[:];map.put('Author',item.author);map.put('branch',item.branch);map.put('commitID',item.commit_id);if(is24HourFormat){map.put('commitTime',formatter.format(item.commit_timestamp));}else{map.put('commitTime',twelveHourFormatter.format(item.commit_timestamp));}map.put('component',item.component_name);map.put('componentId',item.component_id);map.put('repo',item.repository_url);resultList.add(map);}}return resultList;"
				}
			  }
			}
		  }`

const PullRequestDrilldownQuery = `{
	"size": 0, 
	"query": {
	  "bool": {
		"filter": [
		  {
			"range": {
			  "pr_created_time": {
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
		  }
		]
	  }
	},
	"aggs": {
		"pullrequests": {
		  "scripted_metric": {
			"params": {
				"timeZone": "{{.timeZone}}"
			},
			"init_script": "state.statusMap = [:];",
			"map_script": "def map=state.statusMap;def key=doc['org_id'].value+'_'+doc['repository_name'].value+'_'+doc['component_id'].value+'_'+doc['pull_request_id'].value+'_'+(doc.timestamp.getValue().toEpochSecond()*1000);def repoUrl=\"\";if(doc.containsKey('repository_url')&&!doc.repository_url.empty){repoUrl=doc.repository_url.value;}def v=['pull_request_id':doc.pull_request_id.value,'provider':doc.provider.value,'component_id':doc.component_id.value,'component_name':doc.component_name.value,'review_status':doc.review_status.value,'org_id':doc.org_id.value,'pr_created_time':doc.pr_created_time.value,'repository_name':doc.repository_name.value,'repository_url':repoUrl,'timestamp':doc.timestamp.value];if(doc['source_branch'].size()!=0){v['source_branch']=doc.source_branch.value;}if(doc['target_branch'].size()!=0){v['target_branch']=doc.target_branch.value;}map.put(key,v);",
			"combine_script": "return state.statusMap;",
			"reduce_script": "def tmpMap=[:],pullrequests=new ArrayList();def statusMap=new HashMap();statusMap.put('OPEN',0);statusMap.put('CHANGES_REQUESTED',1);statusMap.put('REJECTED',2);statusMap.put('APPROVED',3);statusMap.put('MERGED',4);statusMap.put('CLOSED',5);for(agg in states){if(agg!=null){for(key in agg.keySet()){def record=agg.get(key);def devkey=record.org_id+'_'+record.provider+'_'+record.pull_request_id+'_'+record.component_id+'_'+record.repository_name;if(tmpMap.containsKey(devkey)){def mapRecord=tmpMap.get(devkey);def currentStatusVal=statusMap.get(mapRecord.review_status);def statusVal=statusMap.get(record.review_status);if(currentStatusVal<statusVal){tmpMap.put(devkey,record);}}else{tmpMap.put(devkey,record);}}}}if(tmpMap.size()>0){for(uniqueKey in tmpMap.keySet()){def valueMap=tmpMap.get(uniqueKey);def revStatus=valueMap.review_status;def rd=valueMap.pr_created_time;valueMap.pr_created_time_UTC=rd;valueMap.pr_created_time=rd.withZoneSameInstant(ZoneId.of(params.timeZone));if(revStatus=='APPROVED'){valueMap.review_status='Approved';}else if(revStatus=='OPEN'){valueMap.review_status='Open';}else if(revStatus=='CHANGES_REQUESTED'){valueMap.review_status='Changes requested';}else if(revStatus=='REJECTED'){valueMap.review_status='Rejected';}else if(revStatus=='MERGED'){valueMap.review_status='Merged';}else if(revStatus=='CLOSED'){valueMap.review_status='Closed';}pullrequests.add(valueMap);}}return pullrequests;"
		  }
		}
	}
  }`

const CPSRunInitiatingCommitsQuery = `{
	"_source": false,
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
			{
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
		  	}
		]
	  }
	},
	"aggs": {
		"automation_run": {
		  "scripted_metric": {
			"init_script": "state.data_map=[:];",
			"map_script": "def map = state.data_map;def key = doc.org_id+'_'+doc.automation_id.value + '_' + doc.run_id.value+ '_'+doc.status;def v = ['component_id':doc.component_id.value,'component_name':doc.component_name.value,'automation_id': doc.automation_id.value, 'run_id': doc.run_id.value, 'status': doc.status.value, 'status_timestamp': doc.status_timestamp.value, 'start_time': doc.start_time.value,'run_number':doc['workflow_info.run_number'].value, 'org_name':doc.org_name.value];if (doc.commit_sha.size() > 0) {v['commit_sha'] = doc.commit_sha.value;}if (doc.commit_description.size() > 0) {v['commit_description'] = doc.commit_description.value;}if (map.containsKey(key)) {def record = map.get(key);if (v.status_timestamp.getMillis() > record.status_timestamp.getMillis()) {map.put(key, v);}} else {map.put(key, v);}",
			"combine_script": "return state.data_map;",
			"reduce_script": "def tmpMap = new HashMap();def dataMap = new HashMap();for (response in states) {if (response != null) {for (key in response.keySet()) {def tmpRecord = response.get(key);if (tmpMap.containsKey(key)) {def record = tmpMap.get(key);if (tmpRecord.status_timestamp.getMillis() > record.status_timestamp.getMillis()) {tmpMap.put(key, tmpRecord);}} else {tmpMap.put(key, tmpRecord);}}}}for (key in tmpMap.keySet()) {def record = tmpMap.get(key);def runKey = record.automation_id+'_'+record.run_id;if (record.status == 'SUCCEEDED') {record.status = 'Success';} else if (record.status == 'FAILED' || record.status == 'TIMED_OUT' || record.status == 'ABORTED') {record.status = 'Failure';}if (dataMap.containsKey(runKey)) {def tmpRecord = dataMap.get(runKey);if (tmpRecord.status_timestamp.getMillis() < record.status_timestamp.getMillis()) {dataMap.put(runKey, record);}} else {dataMap.put(runKey, record);}}return dataMap;"
		  }
		}
	}
  }`

const CPSRunCommitsDeployedEnvQuery = `{
	"_source": false,
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
		  {
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
				"org_id": "{{.orgId}}"
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
	  "deployments": {
		"scripted_metric": {
			"init_script": "state.data_map=[:];",
			"map_script": "def map = state.data_map;def key = doc.automation_id.value + '_' + doc.run_id.value + '_' +doc.job_id + '_' + doc.step_id + '_' + doc.target_env + '_' + doc.status.value + '_' + doc.status_timestamp.value;def v = ['component_id':doc.component_id.value, 'component_name':doc.component_name.value, 'automation_id':doc.automation_id.value, 'run_id':doc.run_id.value, 'status':doc.status.value, 'status_timestamp':doc.status_timestamp.value, 'target_env':doc.target_env.value];map.put(key, v);",
			"combine_script": "return state.data_map;",
			"reduce_script": "def tmpMap = new HashMap();def dataMap = new HashMap();def jobStepDedupMap = new HashMap();for (response in states){if (response != null){for (key in response.keySet()){def tmpRecord = response.get(key);if (tmpMap.containsKey(key)){def record = tmpMap.get(key);if (tmpRecord.status_timestamp.getMillis() > record.status_timestamp.getMillis()){tmpMap.put(key, tmpRecord);}} else{tmpMap.put(key, tmpRecord);}}}}for (key in tmpMap.keySet()){def currRecord = tmpMap.get(key);if (currRecord.step_id == ''){jobStepDedupMap.put(key, currRecord);} else{def jobLevelRecordKey = currRecord.automation_id + '_' + currRecord.run_id + '_' + currRecord.job_id + '_' + '' + '_' + currRecord.target_env + '_' + currRecord.status;if (!tmpMap.containsKey(jobLevelRecordKey)){jobStepDedupMap.put(key, currRecord);}}}for (key in jobStepDedupMap.keySet()){def record = jobStepDedupMap.get(key);def runKey = record.run_id;if (dataMap.containsKey(runKey)){def envs = dataMap.get(runKey);envs.add(record.target_env);dataMap.put(runKey, envs);} else{def envs = new HashSet();envs.add(record.target_env);dataMap.put(runKey, envs);}}return dataMap;"
		}
	  }
	}
  }`

const CodeProgressionSnapshotBuild = `{
	"_source": false,
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
				"status_timestamp": {
					"gte": "{{.startDate}}",
					"lte": "{{.endDate}}",
					"format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis",
					"time_zone":"{{.timeZone}}"
				}
			}
		  },
		  {
			"term": {
			  "step_kind": "build"
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
		"builds": {
		  "scripted_metric": {
			"init_script": "state.data_map=[:];",
			"map_script": "def map=state.data_map;def key=doc.run_id.value+'_'+doc.status.value+'_'+doc.step_id.value+'_'+doc.job_id.value+'_'+doc.component_id.value+'_'+doc._id.value;def v=['automation_id':doc.automation_id.value,'component_id':doc.component_id.value,'org_id':doc.org_id.value,'component_name':doc.component_name.value,'target_env':doc.target_env.value,'step_kind':doc.step_kind.value,'run_number':doc.run_number.value,'workflow_name':doc.workflow_name.value,'run_id':doc.run_id.value,'status':doc.status.value,'duration':doc.duration.value,'status_timestamp':doc.status_timestamp.value,'workflow_name':doc.workflow_name.value,'org_name':doc.org_name.value,'job_id':doc.job_id.value,'step_id':doc.step_id.value,'source':doc['source'].size()==0?'':doc.source.value];map.put(key,v);",
			"combine_script": "return state.data_map;",
			"reduce_script": "def allDataMap=new HashMap(),statusMap=new HashMap();for(response in states){if(response!=null){for(key in response.keySet()){def record=response.get(key);if(record.status=='SUCCEEDED'){record.status='Success';}else if(record.status=='FAILED'||record.status=='TIMED_OUT'||record.status=='ABORTED'){record.status='Failure';}allDataMap.put(key,record);}}}def jobLevelKeys=new HashSet();for(record in allDataMap.values()){if(record.step_id==''){def dedupKey=record.component_id+'_'+record.run_id+'_'+record.job_id;jobLevelKeys.add(dedupKey);}}def dedupKeys=new HashSet();for(entry in allDataMap.entrySet()){def record=entry.getValue();def dedupKey=record.component_id+'_'+record.run_id+'_'+record.job_id;def fullKey=entry.getKey();if(jobLevelKeys.contains(dedupKey)){if(record.step_id==''){statusMap.put(fullKey,record);}}else{statusMap.put(fullKey,record);}}return statusMap;"
		  }
		}
	}
  }`

const CodeProgressionSnapshotDeploy = `{
	"_source": false,
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
			  "status_timestamp": {
				"gte": "{{.startDate}}",
				"lte": "{{.endDate}}",
				"format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis",
				"time_zone":"{{.timeZone}}"
			  }
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
				}
			  ]
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
		"deployments": {
		  "scripted_metric": {
			"init_script": "state.data_map=[:];",
			"map_script": "def map = state.data_map;def key = doc.run_id.value + '_' + doc.job_id.value + '_' + doc.step_id.value + '_' + doc.target_env.value + '_' + doc.status.value;def v = ['automation_id':doc.automation_id.value,'component_id':doc.component_id.value, 'org_id':doc.org_id.value,'component_name':doc.component_name.value, 'target_env':doc.target_env.value, 'step_kind':doc.step_kind.value,'run_number':doc.run_number.value,'workflow_name':doc.workflow_name.value, 'run_id': doc.run_id.value, 'status': doc.status.value,'duration':doc.duration.value,'status_timestamp':doc.status_timestamp.value,'workflow_name':doc.workflow_name.value,'org_name':doc.org_name.value,'job_id':doc.job_id.value,'step_id':doc.step_id.value];map.put(key, v);",
			"combine_script": "return state.data_map;",
			"reduce_script": "def tmpMap = [:], out = [:], resultMap = new HashMap(), jobStepDedupMap = new HashMap();for (response in states){if (response != null){for (key in response.keySet()){resultMap.put(key, response.get(key));}}}for (key in resultMap.keySet()){def currRecord = resultMap.get(key);if (currRecord.step_id == ''){jobStepDedupMap.put(key, currRecord);} else{def jobLevelRecordKey = currRecord.run_id + '_' + currRecord.job_id + '_' + '' + '_' + currRecord.target_env + '_' + currRecord.status;if (!resultMap.containsKey(jobLevelRecordKey)){jobStepDedupMap.put(key, currRecord);}}}return jobStepDedupMap;"
		  }
		}
	}
  }`

const SuccessfulBuildDuration = `{
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
			  "status_timestamp": {
				"gte": "{{.startDate}}",
				"lte": "{{.endDate}}",
				"format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis",
				"time_zone":"{{.timeZone}}"
			  }
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
	  "builds": {
		"scripted_metric": {
		  "init_script": "state.statusMap = [:];",
		  "map_script": "def map=state.statusMap;def key=doc.component_id.value+'_'+doc.run_id.value+'_'+doc.job_id.value+'_'+doc.step_id.value+'_'+doc._id.value;def v=['component_id':doc.component_id.value,'component_name':doc.component_name.value,'run_id':doc.run_id.value,'job_id':doc.job_id.value,'step_id':doc.step_id.value,'step_kind':doc.step_kind.value,'target_env':doc.target_env.value,'status':doc.status.value,'status_timestamp':doc.status_timestamp.value,'start_time':doc.start_time.value,'completed_time':doc.completed_time.value,'automation_id':doc.automation_id.value,'duration':doc.duration.value,'run_number':doc.run_number.value,'workflow_name':doc.workflow_name.value,'source':doc['source'].size()==0?'':doc.source.value];map.put(key,v);",
		  "combine_script": "return state.statusMap;",
		  "reduce_script": "def statusMap=new HashMap();def allDataMap=new HashMap();for(response in states){if(response!=null){for(key in response.keySet()){def record=response.get(key);allDataMap.put(key,record);}}}def jobLevelKeys=new HashSet();for(record in allDataMap.values()){if(record.step_id==''){def dedupKey=record.component_id+'_'+record.run_id+'_'+record.job_id;jobLevelKeys.add(dedupKey);}}def dedupKeys=new HashSet();for(entry in allDataMap.entrySet()){def record=entry.getValue();def dedupKey=record.component_id+'_'+record.run_id+'_'+record.job_id;def fullKey=entry.getKey();if(jobLevelKeys.contains(dedupKey)){if(record.step_id==''){statusMap.put(fullKey,record);}}else{statusMap.put(fullKey,record);}}return statusMap;"
		}
	  }
	}
  }`

const DeploymentOverview = `{
	"_source": false,
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
			  "status_timestamp": {
				"gte": "{{.startDate}}",
				"lte": "{{.endDate}}",
				"format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis",
				"time_zone":"{{.timeZone}}"
			  }
			}
		  },
		  {
			"term": {
			  "data_type": 2
			}
		  },
		  {
			"term": {
			  "target_env": "{{.targetEnv}}"
			}
		  }
		]
	  }
	},
	"aggs": {
		"deployments": {
		  "scripted_metric": {
			"init_script": "state.data_map=[:];",
			"map_script": "def map = state.data_map;def key = doc.component_id.value + '_' + doc.run_id.value + '_' + doc.job_id.value + '_' + doc.step_id.value + '_' + doc.target_env.value + '_' + doc.status.value;def v = ['automation_id':doc.automation_id.value,'component_id':doc.component_id.value, 'org_id':doc.org_id.value,'component_name':doc.component_name.value, 'target_env':doc.target_env.value, 'step_kind':doc.step_kind.value,'run_number':doc.run_number.value,'workflow_name':doc.workflow_name.value, 'run_id': doc.run_id.value, 'status': doc.status.value,'duration':doc.duration.value,'status_timestamp':doc.status_timestamp.value,'workflow_name':doc.workflow_name.value,'org_name':doc.org_name.value,'job_id':doc.job_id.value,'step_id':doc.step_id.value];map.put(key, v);",
			"combine_script": "return state.data_map;",
			"reduce_script": "def tmpMap = [:], out = [:], resultMap = new HashMap(), jobStepDedupMap = new HashMap();for (response in states){if (response != null){for (key in response.keySet()){def record = response.get(key);if (record.status == 'SUCCEEDED'){record.status = 'Success';} else if (record.status == 'FAILED' || record.status == 'ABORTED' || record.status == 'TIMED_OUT'){record.status = 'Failure';}resultMap.put(key, record);}}}for (key in resultMap.keySet()){def currRecord = resultMap.get(key);if (currRecord.step_id == ''){jobStepDedupMap.put(key, currRecord);} else{def jobLevelRecordKey = currRecord.component_id + '_' + currRecord.run_id + '_' + currRecord.job_id + '_' + '' + '_' + currRecord.target_env + '_' + currRecord.status;if (!resultMap.containsKey(jobLevelRecordKey)){jobStepDedupMap.put(key, currRecord);}}}return jobStepDedupMap;"
		  }
		}
	}
  }`

const SecurityComponentDrilldownQuery = `{
	"_source": false,
	"size": 0,
	"query": {
	  "bool": {
		  "filter": [
			{
			  "term": {
				"org_id":  "{{.orgId}}"
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
			"map_script": "def map = state.data_map;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.scanner_name.value ;def v = ['component_id': doc.component_id.value, 'component_name': doc.component_name.value, 'scanner_name': doc.scanner_name.value ];map.put(key, v);",
      		"combine_script": "return state.data_map;",
     		"reduce_script": "def tmpMap = [:];def resultMap = new HashMap();for (response in states){if (response != null){for (key in response.keySet()){def record = response.get(key);def compKey = record.component_id;def scannerName =record.scanner_name;if(tmpMap.containsKey(compKey)){ def scannerNamesList = tmpMap.get(compKey);if (scannerNamesList.contains(scannerName)){ continue;}else{ scannerNamesList.add(scannerName);}tmpMap.put(compKey, scannerNamesList);}else{ def scannerNamesList = new ArrayList(); scannerNamesList.add(scannerName); tmpMap.put(compKey, scannerNamesList);}}}}return tmpMap;"
		  }
		}
	  }
	}`

const SecurityAutomationRunsDrilldownQuery = `{
	"_source": false,
  
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
  
				"format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis",

				"time_zone":"{{.timeZone}}"
			  }
			}
		  }
		]
	  }
	},
	"aggs": {
	  "distinct_run": {
		"scripted_metric": {
		  	"init_script": "state.data_map=[:];",
		 	"map_script": "def map = state.data_map;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.automation_id.value + '_' + doc.scanner_name.value + '_' + doc.run_id.value; def v = ['automation_id': doc.automation_id.value, 'scanner_name': doc.scanner_name.value, 'run_id': doc.run_id.value, 'scan_status': doc.scan_status.value];map.put(key, v);",
        	"combine_script": "return state.data_map;",
        	"reduce_script": "def tmpMap = [: ]; for (response in states) { if (response != null) { for (key in response.keySet()) { def record = response.get(key); def runIDKey = record.run_id; def scanStatus = record.scan_status; def scannerName = record.scanner_name; if (tmpMap.containsKey(runIDKey)) { def scannerNamesList = tmpMap.get(runIDKey); def map = scannerNamesList[0]; if (!map.containsKey(scannerName)) { scannerNamesList[0].put(scannerName, scanStatus); } } else { def scannerNamesList = []; def scanStatusMap = new HashMap(); scanStatusMap[scannerName] = scanStatus; scannerNamesList.add(scanStatusMap); tmpMap.put(runIDKey, scannerNamesList); } } } } return tmpMap;"
		}
	  }
	}
  }`

const SecurityScanTypeWorkflowsDrilldownQuery = `{
	"_source": false,
  
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
	  "distinct_run": {
		"scripted_metric": {
			"init_script": "state.data_map=[:];",
			"map_script": "def map = state.data_map;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.automation_id.value + '_' + doc.scanner_name.value + '_' + doc.scanner_type.value + '_' + doc.run_id.value; def v = ['automation_id': doc.automation_id.value, 'scanner_name': doc.scanner_name.value, 'scanner_type': doc.scanner_type.value,'run_id': doc.run_id.value];map.put(key, v);",
			"combine_script": "return state.data_map;",
			"reduce_script": "def tmpMap = [:];def resultMap = new HashMap();for (response in states){if (response != null){for (key in response.keySet()){def record = response.get(key);def runIDKey = record.run_id;def scannerName = record.scanner_name;def scannerType = record.scanner_type;if (tmpMap.containsKey(runIDKey)){def extraMap = tmpMap.get(runIDKey);extraMap.get('scanner_names').add(scannerName);extraMap.get('scanner_types').add(scannerType);tmpMap.put(runIDKey, extraMap);} else{def scannerNamesList = new HashSet();def scannerTypes = new HashSet();def extraMap = new HashMap();scannerNamesList.add(scannerName);scannerTypes.add(scannerType);extraMap.put('scanner_names', scannerNamesList);extraMap.put('scanner_types', scannerTypes);tmpMap.put(runIDKey, extraMap);}}}}return tmpMap;"
		  }
	  }
	}
  }`

const OpenVulnerabilitiesDrillDownQuery = `{
	"_source": false,
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
		  {
			"range": {
			  "scan_time": {
				"gte": "{{.startDate}}",
				"lte": "{{.endDate}}",
				"format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time ||epoch_millis",
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
			"terms": {
			  "severity": ["MEDIUM", "HIGH", "LOW", "VERY_HIGH"]
			}
		  }
		]
	  }
	},
	"aggs": {
	  "drilldowns": {
		"scripted_metric": {
			"params": {
				"timeZone": "{{.timeZone}}",
				"timeFormat":"{{.timeFormat}}"
			},
			"init_script": "state.statusMap = [: ];",
			"map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'recurrences':params['_source']['failure_files'].size(), 'scanner_name':doc.scanner_name.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value, 'scan_time':doc.scan_time.value, 'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':'', 'run_id':doc.run_id.value];map.put(key, v);",
			"combine_script": "return state.statusMap;",
			"reduce_script": "Instant currentDate = Instant.ofEpochMilli(new Date().getTime()); def statusMap = new HashMap(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.code; if (statusMap.containsKey(key)) { def vulDetailsMap = statusMap.get(key); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; if (vulDetailsMap.containsKey(vulKey)) { def lastRecord = vulDetailsMap.get(vulKey); if (lastRecord.timestamp < record.timestamp) { vulDetailsMap.put(vulKey, record) } } else { vulDetailsMap.put(vulKey, record) } } else { def vulDetailsMap = new HashMap(); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; vulDetailsMap.put(vulKey, record); statusMap.put(key, vulDetailsMap) } } } } def resultList = new ArrayList(); DateTimeFormatter formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params .timeZone)); DateTimeFormatter twelveHourFormatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId .of(params.timeZone)); def is24HourFormat = params.timeFormat == '24h'; if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vulMapBranchLevel = statusMap.get(uniqueKey); def vulMapOrgLevel = new HashMap(); def uniqueComponents = new HashSet(); def hasVulDetailsAtOrgLevel = false; def isVulOpen = false; for (vulKey in vulMapBranchLevel.keySet()) { def vul = vulMapBranchLevel.get(vulKey); if (vul.bug_status == 'Open' || vul.bug_status == 'Reopened') { uniqueComponents.add(vul.component_name); isVulOpen = true; Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); def curSeverity = vul.severity; def severityCode = 0; if (curSeverity == 'VERY_HIGH') { vul.severity = 'Very high'; severityCode = 4; } else if (curSeverity == 'HIGH') { vul.severity = 'High'; severityCode = 3; } else if (curSeverity == 'MEDIUM') { vul.severity = 'Medium'; severityCode = 2; } else if (curSeverity == 'LOW') { vul.severity = 'Low'; severityCode = 1; } def diffAge = ChronoUnit.DAYS.between(startDate, currentDate); def SLAToolTip = ''; if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else if (diffAge >= slaRules.AtRisk) { vul.sla = 'At risk'; Instant willBreachOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Will breach on: ' + formatter.format(willBreachOn); } else { vul.sla = 'On track'; } def map = [: ]; map.put('lastDiscovered', formatter.format(vul.scan_time)); map.put('component', vul.component_name); map.put('branch', vul.branch); map.put('scannerName', vul.scanner_name); map.put('recurrences', vul.recurrences); map.put('componentId', vul.component_id); map.put('sla', vul.sla); if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } if (!hasVulDetailsAtOrgLevel) { vulMapOrgLevel.put('vulnerabilityId', vul.code); vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery); vulMapOrgLevel.put('vulnerabilityName', vul.name); vulMapOrgLevel.put('severity', vul.severity); vulMapOrgLevel.put('severityCode', severityCode); hasVulDetailsAtOrgLevel = true } if (vul.date_of_discovery.getMillis() < vulMapOrgLevel.get('firstDiscovered').getMillis()) { vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery) } } } if (isVulOpen) { def fd = vulMapOrgLevel.get('firstDiscovered'); vulMapOrgLevel.put('firstDiscovered', is24HourFormat ? formatter.format(fd) : twelveHourFormatter .format(fd)); vulMapOrgLevel.put('openLocations', uniqueComponents.size()); vulMapOrgLevel.put('customSubRowsInfo', ['enabled': true, 'report_id': 'openVulnerabilitiesSubRows', 'reportInfo': ['code': vulMapOrgLevel.get( 'vulnerabilityId')] ]); resultList.add(vulMapOrgLevel) } } } return resultList;"
		}
	  }
	}
  }`

const VulnerabilitiesOverviewDrillDownQuery = `{
    "_source": false,
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
                    "timeZone": "{{.timeZone}}",
					"timeFormat":"{{.timeFormat}}"
                },
                "init_script": "state.statusMap = [: ];",
                "map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'recurrences':params['_source']['failure_files'].size(), 'scanner_name':doc.scanner_name.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value, 'scan_time':doc.scan_time.value, 'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':'', 'run_id':doc.run_id.value];map.put(key, v);",
                "combine_script": "return state.statusMap;",
                "reduce_script": "Instant currentDate = Instant.ofEpochMilli(new Date().getTime()); def statusMap = new HashMap(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.code; if (statusMap.containsKey(key)) { def vulDetailsMap = statusMap.get(key); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record.code + '_' + record.scanner_name; if (vulDetailsMap.containsKey(vulKey)) { def lastRecord = vulDetailsMap.get(vulKey); if (lastRecord.timestamp < record.timestamp) { vulDetailsMap.put(vulKey, record); } } else { vulDetailsMap.put(vulKey, record); } } else { def vulDetailsMap = new HashMap(); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record.code + '_' + record.scanner_name; vulDetailsMap.put(vulKey, record); statusMap.put(key, vulDetailsMap); } } } } def resultList = new ArrayList(); DateTimeFormatter formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); DateTimeFormatter twelveHourFormatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); def is24HourFormat = params.timeFormat == '24h'; if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vulMapBranchLevel = statusMap.get(uniqueKey); def vulMapOrgLevel = new HashMap(); def uniqueComponents = new HashSet(); def hasVulDetailsAtOrgLevel = false; def isVulReopened = false; for (vulKey in vulMapBranchLevel.keySet()) { def vul = vulMapBranchLevel.get(vulKey); uniqueComponents.add(vul.component_name); Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); def severityCode = 0; def curSeverity = vul.severity; if (curSeverity == 'VERY_HIGH') { vul.severity = 'Very high'; severityCode = 4; } else if (curSeverity == 'HIGH') { vul.severity = 'High'; severityCode = 3; } else if (curSeverity == 'MEDIUM') { vul.severity = 'Medium'; severityCode = 2; } else if (curSeverity == 'LOW') { vul.severity = 'Low'; severityCode = 1; } def diffAge = ChronoUnit.DAYS.between(startDate, currentDate); def SLAToolTip = ''; def statusToolTip = ''; if (vul.bug_status == 'Resolved') { Instant resolutionDate = Instant.ofEpochMilli(vul.scan_time.getMillis()); diffAge = ChronoUnit.DAYS.between(startDate, resolutionDate); statusToolTip = 'Date of resolution: ' + formatter.format(vul.scan_time); if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else { vul.sla = 'Within SLA'; } } else { if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else if (diffAge >= slaRules.AtRisk) { vul.sla = 'At risk'; Instant willBreachOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Will breach on: ' + formatter.format(willBreachOn); } else { vul.sla = 'On track'; } } def map = [:]; map.put('lastDiscovered', is24HourFormat ? formatter.format(vul.scan_time) : twelveHourFormatter.format(vul.scan_time)); map.put('component', vul.component_name); map.put('branch', vul.branch); map.put('scannerName', vul.scanner_name); map.put('recurrences', vul.recurrences); map.put('componentId', vul.component_id); map.put('sla', vul.sla); map.put('status', vul.bug_status); if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } if (statusToolTip != '') { map.put('statusToolTipContent', statusToolTip); } if (!hasVulDetailsAtOrgLevel) { vulMapOrgLevel.put('vulnerabilityId', vul.code); vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery); vulMapOrgLevel.put('vulnerabilityName', vul.name); vulMapOrgLevel.put('severity', vul.severity); vulMapOrgLevel.put('severityCode', severityCode); vulMapOrgLevel.put('status', 'Resolved'); hasVulDetailsAtOrgLevel = true; } if (vul.date_of_discovery.getMillis() < vulMapOrgLevel.get('firstDiscovered').getMillis()) { vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery); } if (!isVulReopened) { if (vul.bug_status == 'Open') { vulMapOrgLevel.put('status', 'Open'); } else if (vul.bug_status == 'Reopened') { isVulReopened = true; vulMapOrgLevel.put('status', 'Reopened'); } } } def fd = vulMapOrgLevel.get('firstDiscovered'); vulMapOrgLevel.put('firstDiscovered', is24HourFormat ? formatter.format(fd) : twelveHourFormatter.format(fd)); vulMapOrgLevel.put('identifiedComponents', uniqueComponents.size()); vulMapOrgLevel.put('customSubRowsInfo', ['enabled': true, 'report_id': 'vulnerabilitiesOverviewSubRows', 'reportInfo': ['code': vulMapOrgLevel.get('vulnerabilityId')]]); resultList.add(vulMapOrgLevel); } } return resultList;"
            }
        }
    }
}`

const SLAStatusOverviewOpenDrillDownQuery = `{
"_source": false,
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
								"timeZone": "{{.timeZone}}",
								"timeFormat":"{{.timeFormat}}"
							},
							"init_script": "state.statusMap = [: ];",
							"map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value,'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':''];map.put(key, v);",
							"combine_script": "return state.statusMap;",
							"reduce_script": "Instant currentDate = Instant.ofEpochMilli(new Date().getTime()); def statusMap = new HashMap(); def resultList = new ArrayList(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record.code; if (statusMap.containsKey(key)) { def lastRecord = statusMap.get(key); if (lastRecord.timestamp < record.timestamp) { statusMap.put(key, record); } } else { statusMap.put(key, record); } } } } DateTimeFormatter formatter; if (params.timeFormat == '12h') { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); } else { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); } if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vul = statusMap.get(uniqueKey); if (vul.bug_status == 'Open' || vul.bug_status == 'Reopened') { Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); def diffAge = ChronoUnit.DAYS.between(startDate, currentDate); def severityCode = 0; def curSeverity = vul.severity; if (curSeverity == 'VERY_HIGH') { vul.severity = 'Very high'; severityCode = 4; } else if (curSeverity == 'HIGH') { vul.severity = 'High'; severityCode = 3; } else if (curSeverity == 'MEDIUM') { vul.severity = 'Medium'; severityCode = 2; } else if (curSeverity == 'LOW') { vul.severity = 'Low'; severityCode = 1; } def SLAToolTip = ''; if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else if (diffAge >= slaRules.AtRisk) { vul.sla = 'At risk'; Instant willBreachOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Will breach on: ' + formatter.format(willBreachOn); } else { vul.sla = 'On track'; } def map = [: ]; if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } map.put('dateOfDiscovery', formatter.format(vul.date_of_discovery)); map.put('vulnerabilityName', vul.name); map.put('ComponentId', vul.component_id); map.put('componentName', vul.component_name); map.put('severity', vul.severity); map.put('severityCode', severityCode); map.put('sla', vul.sla); map.put('status', vul.bug_status); resultList.add(map); } } } return resultList;"
						}
					}
				}
			}`

const SLAStatusOverviewClosedDrillDownQuery = `{
	"_source": false,
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
							"timeZone": "{{.timeZone}}",
							"timeFormat":"{{.timeFormat}}"
						},
						"init_script": "state.statusMap = [: ];",
						"map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value,'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':'','scan_time':doc.scan_time.value];map.put(key, v);",
						"combine_script": "return state.statusMap;",
						"reduce_script": "def statusMap = new HashMap(); def resultList = new ArrayList(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record.code; if (statusMap.containsKey(key)) { def lastRecord = statusMap.get(key); if (lastRecord.timestamp < record.timestamp) { statusMap.put(key, record); } } else { statusMap.put(key, record); } } } } DateTimeFormatter formatter; if (params.timeFormat == '12h') { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); } else { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); } if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vul = statusMap.get(uniqueKey); if (vul.bug_status == 'Resolved') { Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); Instant scanDate = Instant.ofEpochMilli(vul.scan_time.getMillis()); def diffAge = ChronoUnit.DAYS.between(startDate, scanDate); def severityCode = 0; def curSeverity = vul.severity; if (curSeverity == 'VERY_HIGH') { vul.severity = 'Very high'; severityCode = 4; } else if (curSeverity == 'HIGH') { vul.severity = 'High'; severityCode = 3; } else if (curSeverity == 'MEDIUM') { vul.severity = 'Medium'; severityCode = 2; } else if (curSeverity == 'LOW') { vul.severity = 'Low'; severityCode = 1; } def SLAToolTip = ''; def statusToolTip = ''; statusToolTip = 'Date of resolution: ' + formatter.format(vul.scan_time); if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else { vul.sla = 'Within SLA'; } def map = [: ]; if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } if (statusToolTip != '') { map.put('statusToolTipContent', statusToolTip); } map.put('dateOfDiscovery', formatter.format(vul.date_of_discovery)); map.put('vulnerabilityName', vul.name); map.put('ComponentId', vul.component_id); map.put('componentName', vul.component_name); map.put('severity', vul.severity); map.put('severityCode', severityCode); map.put('sla', vul.sla); map.put('status', vul.bug_status); resultList.add(map); } } } return resultList;"
					}
				}
			}
		}`

const VulnerabiltyByScannerTypeDrillDownQuery = `{
    "_source": false,
    "size": 0,
    "query": {
        "bool": {
            "filter": [
                {
                    "range": {
                        "timestamp": {
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
                    "timeZone": "{{.timeZone}}",
					"timeFormat":"{{.timeFormat}}"
                },
                "init_script": "state.statusMap = [: ];",
                "map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'recurrences':params['_source']['failure_files'].size(), 'scanner_name':doc.scanner_name.value,'scanner_type':doc.scanner_type.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value, 'scan_time':doc.scan_time.value, 'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':'', 'run_id':doc.run_id.value];map.put(key, v);",
                "combine_script": "return state.statusMap;",
                "reduce_script": "Instant currentDate = Instant.ofEpochMilli(new Date().getTime()); def statusMap = new HashMap(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.code; if (statusMap.containsKey(key)) { def vulDetailsMap = statusMap.get(key); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; if (vulDetailsMap.containsKey(vulKey)) { def lastRecord = vulDetailsMap.get(vulKey); if (lastRecord.timestamp < record.timestamp) { vulDetailsMap.put(vulKey, record) } } else { vulDetailsMap.put(vulKey, record) } } else { def vulDetailsMap = new HashMap(); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; vulDetailsMap.put(vulKey, record); statusMap.put(key, vulDetailsMap) } } } } def resultList = new ArrayList(); DateTimeFormatter formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params .timeZone)); DateTimeFormatter twelveHourFormatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId .of(params.timeZone)); def is24HourFormat = params.timeFormat == '24h'; if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vulMapBranchLevel = statusMap.get(uniqueKey); def vulDetailsListBranchLevel = new ArrayList(); def vulMapOrgLevel = new HashMap(); def uniqueComponents = new HashSet(); def scannerTypes = new HashSet(); def hasVulDetailsAtOrgLevel = false; for (vulKey in vulMapBranchLevel.keySet()) { def vul = vulMapBranchLevel.get(vulKey); uniqueComponents.add(vul.component_name); scannerTypes.add(vul.scanner_type); Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); def severityCode = 0; def curSeverity = vul.severity; if (curSeverity == 'VERY_HIGH') { vul.severity = 'Very high'; severityCode = 4; } else if (curSeverity == 'HIGH') { vul.severity = 'High'; severityCode = 3; } else if (curSeverity == 'MEDIUM') { vul.severity = 'Medium'; severityCode = 2; } else if (curSeverity == 'LOW') { vul.severity = 'Low'; severityCode = 1; } def diffAge = ChronoUnit.DAYS.between(startDate, currentDate); def SLAToolTip = ''; def statusToolTip = ''; if (vul.bug_status == 'Resolved') { statusToolTip = 'Date of resolution: ' + formatter.format(vul.scan_time); if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else { vul.sla = 'Within SLA'; } } else { if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else if (diffAge >= slaRules.AtRisk) { vul.sla = 'At risk'; Instant willBreachOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Will breach on: ' + formatter.format(willBreachOn); } else { vul.sla = 'On track'; } } def map = [: ]; if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } if (statusToolTip != '') { map.put('statusToolTipContent', statusToolTip); } map.put('lastDiscovered', formatter.format(vul.scan_time)); map.put('component', vul.component_name); map.put('componentId', vul.component_id); map.put('branch', vul.branch); map.put('scannerName', vul.scanner_name); map.put('recurrences', vul.recurrences); map.put('sla', vul.sla); map.put('status', vul.bug_status); if (!hasVulDetailsAtOrgLevel) { vulMapOrgLevel.put('vulnerabilityId', vul.code); vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery); vulMapOrgLevel.put('vulnerabilityName', vul.name); vulMapOrgLevel.put('severity', vul.severity); vulMapOrgLevel.put('severityCode', severityCode); hasVulDetailsAtOrgLevel = true } if (vul.date_of_discovery.getMillis() < vulMapOrgLevel.get('firstDiscovered').getMillis()) { vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery) } } def fd = vulMapOrgLevel.get('firstDiscovered'); vulMapOrgLevel.put('firstDiscovered', is24HourFormat ? formatter.format(fd) : twelveHourFormatter .format(fd)); vulMapOrgLevel.put('foundLocations', uniqueComponents.size()); vulMapOrgLevel.put('scanType', scannerTypes); vulMapOrgLevel.put('customSubRowsInfo', ['enabled': true, 'report_id': 'vulnerabilitiesSecurityScanTypeSubRows', 'reportInfo': ['code': vulMapOrgLevel.get( 'vulnerabilityId')] ]); resultList.add(vulMapOrgLevel); } } return resultList;"
            }
        }
    }
}`

const Top25OpenVulnerabilitiesDrillDownQuery = `{
	"_source": false,
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
			"terms": {
			  "severity": ["MEDIUM", "HIGH", "LOW", "VERY_HIGH"]
			}
		  },
		  {
			"terms": {
			  "code": [
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
				"CWE-276"
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
				"timeZone": "{{.timeZone}}",
				"timeFormat":"{{.timeFormat}}"
			},
			"init_script": "state.statusMap = [: ];",
			"map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'recurrences':params['_source']['failure_files'].size(), 'scanner_name':doc.scanner_name.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value, 'scan_time':doc.scan_time.value, 'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':'', 'run_id':doc.run_id.value];map.put(key, v);",
			"combine_script": "return state.statusMap;",
			"reduce_script": "Instant currentDate = Instant.ofEpochMilli(new Date().getTime()); def statusMap = new HashMap(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.code; if (statusMap.containsKey(key)) { def vulDetailsMap = statusMap.get(key); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; if (vulDetailsMap.containsKey(vulKey)) { def lastRecord = vulDetailsMap.get(vulKey); if (lastRecord.timestamp < record.timestamp) { vulDetailsMap.put(vulKey, record) } } else { vulDetailsMap.put(vulKey, record) } } else { def vulDetailsMap = new HashMap(); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; vulDetailsMap.put(vulKey, record); statusMap.put(key, vulDetailsMap) } } } } def resultList = new ArrayList(); DateTimeFormatter formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params .timeZone)); DateTimeFormatter twelveHourFormatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId .of(params.timeZone)); def is24HourFormat = params.timeFormat == '24h'; if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vulMapBranchLevel = statusMap.get(uniqueKey); def vulMapOrgLevel = new HashMap(); def uniqueComponents = new HashSet(); def hasVulDetailsAtOrgLevel = false; def isVulOpen = false; for (vulKey in vulMapBranchLevel.keySet()) { def vul = vulMapBranchLevel.get(vulKey); if (vul.bug_status == 'Open' || vul.bug_status == 'Reopened') { uniqueComponents.add(vul.component_name); isVulOpen = true; Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); def severityCode = 0; def curSeverity = vul.severity; if (curSeverity == 'VERY_HIGH') { vul.severity = 'Very high'; severityCode = 4; } else if (curSeverity == 'HIGH') { vul.severity = 'High'; severityCode = 3; } else if (curSeverity == 'MEDIUM') { vul.severity = 'Medium'; severityCode = 2; } else if (curSeverity == 'LOW') { vul.severity = 'Low'; severityCode = 1; } def diffAge = ChronoUnit.DAYS.between(startDate, currentDate); def SLAToolTip = ''; if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else if (diffAge >= slaRules.AtRisk) { vul.sla = 'At risk'; Instant willBreachOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Will breach on: ' + formatter.format(willBreachOn); } else { vul.sla = 'On track'; } def map = [: ]; if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } map.put('lastDiscovered', is24HourFormat ? formatter.format(vul.scan_time) : twelveHourFormatter .format(vul.scan_time)); map.put('component', vul.component_name); map.put('branch', vul.branch); map.put('scannerName', vul.scanner_name); map.put('recurrences', vul.recurrences); map.put('componentId', vul.component_id); map.put('sla', vul.sla); if (!hasVulDetailsAtOrgLevel) { vulMapOrgLevel.put('vulnerabilityId', vul.code); vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery); vulMapOrgLevel.put('vulnerabilityName', vul.name); vulMapOrgLevel.put('severity', vul.severity); vulMapOrgLevel.put('severityCode', severityCode); hasVulDetailsAtOrgLevel = true } if (vul.date_of_discovery.getMillis() < vulMapOrgLevel.get('firstDiscovered').getMillis()) { vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery) } } } if (isVulOpen) { def fd = vulMapOrgLevel.get('firstDiscovered'); vulMapOrgLevel.put('firstDiscovered',is24HourFormat ? formatter.format(fd) : twelveHourFormatter .format(fd)); vulMapOrgLevel.put('openLocations', uniqueComponents.size()); vulMapOrgLevel.put('customSubRowsInfo', ['enabled': true, 'report_id': 'cweTop25VulnerabilitiesSubRows', 'reportInfo': ['code': vulMapOrgLevel.get( 'vulnerabilityId')] ]); resultList.add(vulMapOrgLevel) } } } return resultList;"
		}
	  }
	}
  }`
const ViewLocationsNestedDrillDownQuery = `{
    "query": {
        "bool": {
          "filter": [
            {
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
            "org_id": "{{.orgId}}"
          } 
        },
        {
          "term": {
            "component_id": "{{.componentIdForNestedDrillDown}}"
          } 
        },
        {
          "term": {
            "github_branch": "{{.branch}}"
          } 
        },
        {
          "term": {
            "code": "{{.vulCode}}"
          } 
        },
        {
          "term": {
            "scanner_name": "{{.scannerName}}"
          } 
        },
        {
          "term": {
            "run_id": "{{.runId}}"
          } 
        },
            {
              "exists": {
                "field": "date_of_discovery"
              }
            },
            {
              "terms": {
                "severity": ["MEDIUM", "HIGH", "LOW", "VERY_HIGH"]
              }
            }
          ]
        }
      },
      "aggs": {
        "drilldowns": {
          "scripted_metric": {
            "init_script": "state.statusMap = [: ];",
            "map_script": "def map = state.statusMap;map.put('current_sub_row_failures',params['_source']['failure_files']);map.put('scanner_name',doc.scanner_name.value);map.put('repo_link',doc.github_branch.value);",
            "combine_script": "return state.statusMap;",
            "reduce_script": "def scannerLabelsMetadataMap = new HashMap();def locationLabelsMap = new HashMap();locationLabelsMap.put('snyksca','Issue Found In');locationLabelsMap.put('trivy','Package Name');locationLabelsMap.put('checkmarx','Location');locationLabelsMap.put('sonarqube','Location');locationLabelsMap.put('xray','Path');locationLabelsMap.put('findsecbugs','Class name');locationLabelsMap.put('mendsast','Location');locationLabelsMap.put('snykcontainer','Package Name');locationLabelsMap.put('trufflehogcontainer','Package Name');locationLabelsMap.put('snyksast','Location');locationLabelsMap.put('trufflehogsast','Location');locationLabelsMap.put('trufflehogs3','Location');locationLabelsMap.put('zap','Paths');locationLabelsMap.put('anchore','Package Name');locationLabelsMap.put('mendsca','Package Name');def messageLabelsMap = new HashMap();messageLabelsMap.put('checkmarx','Description');messageLabelsMap.put('sonarqube','Message');messageLabelsMap.put('xray','Description');messageLabelsMap.put('snyksast','Message');messageLabelsMap.put('trufflehogsast','Message');messageLabelsMap.put('trufflehogs3','Message');messageLabelsMap.put('trufflehogcontainer','Message');messageLabelsMap.put('zap','Solution');messageLabelsMap.put('mendsca','Fix Resolution');scannerLabelsMetadataMap.put('vul_location_labels',locationLabelsMap);scannerLabelsMetadataMap.put('vul_message_labels',messageLabelsMap);def resultArray = new ArrayList();for (a in states){if (a.size()>0){def scannerName = a.get('scanner_name');def locationLabel = scannerLabelsMetadataMap.get('vul_location_labels').get(scannerName);def messageLabel = scannerLabelsMetadataMap.get('vul_message_labels').get(scannerName);def failureArray = a.get('current_sub_row_failures');for (int i = 0; i<failureArray.size(); i++){def failures = new HashMap();def curFailureLocation;def curFailureMessage;if (locationLabel != null){curFailureLocation = failureArray[i].get(locationLabel);} else{curFailureLocation ='No Failure Location Data Found';}if (messageLabel != null){curFailureMessage = failureArray[i].get(messageLabel);} else{curFailureMessage ='No Message Data Found';}failures.put('repo',a.get('repo_link'));failures.put('locations',curFailureLocation);failures.put('message',curFailureMessage);resultArray.add(failures);}break;}}return resultArray;"
          }
        }
      }
  }`

const OpenIssuesDrillDownQuery = `{
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
                    "term": {
                        "component_id": "{{.componentIdForNestedDrillDown}}"
                    }
                },
                {
                    "term": {
                        "code": "{{.vulCode}}"
                    }
                },
                {
                    "terms": {
                        "scanner_name": {{json .scannerNameList}}
                    }
                },
                {
                    "terms": {
                        "run_id": {{json .runIDList}}
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
        "drilldowns": {
            "scripted_metric": {
                "params": {
                    "timeZone": "{{.timeZone}}"
                },
                "init_script": "state.statusMap = [: ];",
                "map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc.run_id.value;def v = ['org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'branch':doc.github_branch.value, 'scanner_name':doc.scanner_name.value, 'code':doc.code.value, 'current_sub_row_failures':params['_source']['failure_files'], 'date_of_discovery':doc.date_of_discovery.value];map.put(key, v);",
                "combine_script": "return state.statusMap;",
                "reduce_script": "Instant Currentdate = Instant.ofEpochMilli(new Date().getTime()); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); def scannerLabelsMetadataMap = new HashMap(); def locationLabelsMap = new HashMap(); locationLabelsMap.put('snyksca', 'Issue Found In'); locationLabelsMap.put('trivy', 'PkgName'); locationLabelsMap.put('checkmarx', 'Location'); locationLabelsMap.put('sonarqube', 'Location'); locationLabelsMap.put('xray', 'Path'); locationLabelsMap.put('findsecbugs', 'Class name'); locationLabelsMap.put('mendsast', 'Location'); locationLabelsMap.put('snykcontainer', 'Package Name'); locationLabelsMap.put('trufflehogcontainer', 'Package Name'); locationLabelsMap.put('snyksast', 'Location'); locationLabelsMap.put('trufflehogsast', 'Location'); locationLabelsMap.put('zap', 'Paths'); locationLabelsMap.put('anchore', 'Package Name'); locationLabelsMap.put('mendsca', 'Package Name'); locationLabelsMap.put('gitleaks', 'Location'); def messageLabelsMap = new HashMap(); messageLabelsMap.put('checkmarx', 'Description'); messageLabelsMap.put('sonarqube', 'Message'); messageLabelsMap.put('xray', 'Description'); messageLabelsMap.put('snyksast', 'Message'); messageLabelsMap.put('trufflehogsast', 'Message'); messageLabelsMap.put('trufflehogcontainer', 'Message'); messageLabelsMap.put('zap', 'Solution'); messageLabelsMap.put('mendsca', 'Fix Resolution'); messageLabelsMap.put('gitleaks', 'Rule ID'); scannerLabelsMetadataMap.put('vul_location_labels', locationLabelsMap); scannerLabelsMetadataMap.put('vul_message_labels', messageLabelsMap); DateTimeFormatter formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); def resultArray = new ArrayList(); def uniqueRunsSet = new HashSet(); for (a in states) { if (a.size() > 0) { for (key in a.keySet()) { if (!uniqueRunsSet.contains(key)) { uniqueRunsSet.add(key); def tempMap = a.get(key); def date_of_dis = tempMap.get('date_of_discovery'); def failureArray = tempMap.get('current_sub_row_failures'); def scannerName = tempMap.get('scanner_name'); def assetID = tempMap.get('branch'); def sla = ''; Instant Startdate = Instant.ofEpochMilli(date_of_dis.getMillis()); def diffAge = ChronoUnit.DAYS.between(Startdate, Currentdate); if (diffAge >= slaRules.Breached) { sla = 'Breached' } else if (diffAge >= slaRules.AtRisk) { sla = 'At Risk' } else { sla = 'On Track' } def locationLabel = scannerLabelsMetadataMap.get('vul_location_labels').get(scannerName); def messageLabel = scannerLabelsMetadataMap.get('vul_message_labels').get(scannerName); for (int i = 0; i < failureArray.size(); i++) { def failures = new HashMap(); def curFailureLocation; def curFailureMessage; if (locationLabel != null) { curFailureLocation = failureArray[i].get(locationLabel); } else { curFailureLocation = 'No Failure Location Data Found'; } if (messageLabel != null) { curFailureMessage = failureArray[i].get(messageLabel); } else { curFailureMessage = 'No Message Data Found'; } failures.put('repo', assetID); failures.put('discoveredOn', formatter.format(date_of_dis)); failures.put('scannerName', scannerName); failures.put('locations', curFailureLocation); failures.put('message', curFailureMessage); failures.put('sla', sla); resultArray.add(failures); } } } } } return resultArray;"            
			}
        }
    }
}`

const SummaryLatestTestResultsDrilldownQuery = `{
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
                        "automation_id": "{{.automationId}}"
                    }
                },
                {
                    "term": {
                        "branch_id": "{{.branch}}"
                    }
                },
                {
                    "term": {
                        "test_suite_name": "{{.testSuiteName}}"
                    }
                }
            ]
        }
    },
    "aggs": {
        "drilldowns": {
            "scripted_metric": {
                "params": {
                    "timeZone": "{{.timeZone}}",
                    "runId": "{{.runId}}",
					"runNumber": "{{.runNumber}}"
                },
                "init_script": "state.map = new HashMap(); def allTestCasesMap = new HashMap(); def latestTestCasesList = new HashSet(); state.map.put('allTestCasesMap', allTestCasesMap); state.map.put('latestTestCasesList',latestTestCasesList);",
                "map_script": "def map=state.map;def uniqueTestCaseKey=doc['component_id'].size()>0&&doc['test_suite_name'].size()>0&&doc['test_case_name'].size()>0&&doc['run_id'].size()>0&&doc['run_number'].size()>0?doc['component_id'].value+'_'+doc['test_suite_name'].value+'_'+doc['test_case_name'].value+'_'+doc['run_id'].value+'_'+doc['run_number'].value:null;if(uniqueTestCaseKey!=null){def v=['component_id':doc['component_id'].size()>0?doc['component_id'].value:null,'automation_id':doc['automation_id'].size()>0?doc['automation_id'].value:null,'test_suite_name':doc['test_suite_name'].size()>0?doc['test_suite_name'].value:null,'run_id':doc['run_id'].size()>0?doc['run_id'].value:null,'run_number':doc['run_number'].size()>0?doc['run_number'].value:null,'test_case_name':doc['test_case_name'].size()>0?doc['test_case_name'].value:null,'test_case_status':doc['status'].size()>0?doc['status'].value:null,'source':doc['source'].size()>0?doc['source'].value:'CloudBees','test_case_duration':doc['duration'].size()>0?doc['duration'].value:null,'test_case_start_time':doc['start_time'].size()>0?doc['start_time'].value:null];map.get('allTestCasesMap').put(uniqueTestCaseKey,v);if(params.runId==doc['run_id'].value&&params.runNumber==doc['run_number'].value){map.get('latestTestCasesList').add(doc['test_case_name'].value);}}",
                "combine_script": "return state.map;",
                "reduce_script": "def resultMap=new HashMap();def consolidatedMap=new HashMap();def testCaseMeta=new HashMap();def resultList=new ArrayList();DateTimeFormatter formatter=DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone));for(state in states){if(state!=null){def allTestCasesMap=state.get('allTestCasesMap');def latestTestCasesList=state.get('latestTestCasesList');if(allTestCasesMap.size()>0){for(key in allTestCasesMap.keySet()){def testCaseRecord=allTestCasesMap.get(key);consolidatedMap.put(key,testCaseRecord);}}}}if(consolidatedMap.size()>0){for(uniqueKey in consolidatedMap.keySet()){def item=consolidatedMap.get(uniqueKey);def curTestCase=item.get('test_case_name');def testCaseHistory=testCaseMeta.get(curTestCase);if(testCaseHistory!=null){def curTestCaseStatus=item.get('test_case_status');def statusCount=testCaseHistory.get(curTestCaseStatus);testCaseHistory.put(curTestCaseStatus,statusCount+1);}else{def newTestCaseHistory=new HashMap();newTestCaseHistory.put('PASSED',0);newTestCaseHistory.put('FAILED',0);newTestCaseHistory.put('SKIPPED',0);def curTestCaseStatus=item.get('test_case_status');newTestCaseHistory.put(curTestCaseStatus,1);testCaseMeta.put(curTestCase,newTestCaseHistory);}if(item.run_id==params.runId&&item.run_number==params.runNumber){def map=[:];def reportInfoMap=new HashMap();reportInfoMap.put('component_id',item.component_id);reportInfoMap.put('automation_id',item.automation_id);reportInfoMap.put('test_suite_name',item.test_suite_name);reportInfoMap.put('test_case_name',item.test_case_name);def drillDownInfoMap=new HashMap();drillDownInfoMap.put('reportId','test-overview-view-run-activity');drillDownInfoMap.put('reportTitle','Test case activity - '+item.test_case_name);drillDownInfoMap.put('reportInfo',reportInfoMap);def viewRunActivityMap=new HashMap();viewRunActivityMap.put('drillDown',drillDownInfoMap);map.put('viewRunActivity',viewRunActivityMap);map.put('testCaseName',item.test_case_name);map.put('lastRun',formatter.format(item.test_case_start_time));map.put('lastRunEpochMillis',item.test_case_start_time!=null?item.test_case_start_time.getMillis():0);map.put('status',item.test_case_status);map.put('source',item.source);map.put('runDuration',item.test_case_duration);resultMap.put(item.test_case_name,map);}}}for(case in resultMap.keySet()){def tempMap=resultMap.get(case);def caseHistory=testCaseMeta.get(case);def totalRuns=caseHistory.get('PASSED')+caseHistory.get('FAILED')+caseHistory.get('SKIPPED');def failureRate=totalRuns>0?caseHistory.get('FAILED')*100/totalRuns:0;tempMap.put('failureRateValue',failureRate);def failureRateMap=new HashMap();def colorSchemeList=new ArrayList();def lightColorSchemeList=new ArrayList();def color1=[:];def color2=[:];def color3=[:];def lightColor1=[:];def lightColor2=[:];def lightColor3=[:];color1.put('color0','#009C5B');color1.put('color1','#62CA9D');colorSchemeList.add(color1);color2.put('color0','#D32227');color2.put('color1','#FB6E72');colorSchemeList.add(color2);color3.put('color0','#F2A414');color3.put('color1','#FFE6C1');colorSchemeList.add(color3);lightColor1.put('color0','#0C9E61');lightColor1.put('color1','#79CAA8');lightColorSchemeList.add(lightColor1);lightColor2.put('color0','#E83D39');lightColor2.put('color1','#F39492');lightColorSchemeList.add(lightColor2);lightColor3.put('color0','#F2A414');lightColor3.put('color1','#FFE6C1');lightColorSchemeList.add(lightColor3);def data=new ArrayList();def data1=[:];def data2=[:];def data3=[:];data1.put('title','Succesful runs');data1.put('value',caseHistory.get('PASSED'));data.add(data1);data2.put('title','Failed runs');data2.put('value',caseHistory.get('FAILED'));data.add(data2);data3.put('title','Skipped runs');data3.put('value',caseHistory.get('SKIPPED'));data.add(data3);failureRateMap.put('type','SINGLE_BAR');failureRateMap.put('colorScheme',colorSchemeList);failureRateMap.put('lightColorScheme',lightColorSchemeList);failureRateMap.put('data',data);failureRateMap.put('totalRuns',totalRuns);tempMap.put('failureRate',failureRateMap);String testCaseStatus=tempMap.get('status');String cap='';if(testCaseStatus.length()>1){cap=testCaseStatus.substring(0,1).toUpperCase()+testCaseStatus.substring(1).toLowerCase();}tempMap.put('status',cap);resultList.add(tempMap);}resultList.sort((a,b)->{def aValue=a.get('failureRateValue');def bValue=b.get('failureRateValue');return Double.compare(bValue,aValue);});return resultList;"
            }
        }
    }
}`

const SecurityAutomationDrilldownQuery = `{
	"_source": false,
	"size": 0,
	"query": {
	  "bool": {
		  "filter": [
			{
			  "term": {
				"org_id":  "{{.orgId}}"
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
		
	},"aggs": {
		"distinct_automation": {
		  "scripted_metric": {
			"init_script": "state.data_map=[:];",
			"map_script": "def map = state.data_map;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.automation_id.value + '_' + doc.scanner_name.value + '_' + doc.run_id.value; def v = ['automation_id': doc.automation_id.value, 'scanner_name': doc.scanner_name.value, 'run_id': doc.run_id.value];map.put(key, v);",
      		"combine_script": "return state.data_map;",
     		"reduce_script": "def tmpMap = [:];def resultMap = new HashMap();for (response in states){if (response != null){for (key in response.keySet()){def record = response.get(key);def autKey = record.automation_id;def scannerName = record.scanner_name;def runID = record.run_id;if (tmpMap.containsKey(autKey)){def subMap = tmpMap.get(autKey);subMap.get('Scanner_List').add(scannerName);subMap.get('run_ids').add(runID);def count = subMap.get('run_ids');subMap.put('run_count',count.size());tmpMap.put(autKey, subMap);} else{def subMap = new HashMap();def scannerNamesList = new HashSet();def runidList = new HashSet();scannerNamesList.add(scannerName);runidList.add(runID);subMap.put('Scanner_List',scannerNamesList);subMap.put('run_ids',runidList);subMap.put('run_count',runidList.size());tmpMap.put(autKey, subMap);}}}}return tmpMap;"
		}
		}
	}
}`

const MTTRDrillDownQuery = `{
	"_source": false,
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
			"terms": {
			  "severity": ["MEDIUM", "HIGH", "LOW", "VERY_HIGH"]
			}
		  }
		]
	  }
	},
	"aggs": {
	  "drilldowns": {
		"scripted_metric": {
			"params": {
				"timeZone": "{{.timeZone}}",
				"timeFormat": "{{.timeFormat}}"
			},
			"init_script": "state.statusMap = [: ];",
			"map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'recurrences':params['_source']['failure_files'].size(), 'scanner_name':doc.scanner_name.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value, 'scan_time':doc.scan_time.value, 'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':''];map.put(key, v);",
			"combine_script": "return state.statusMap;",
			"reduce_script": "def statusMap = new HashMap(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.code; if (statusMap.containsKey(key)) { def vulDetailsMap = statusMap.get(key); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; if (vulDetailsMap.containsKey(vulKey)) { def lastRecord = vulDetailsMap.get(vulKey); if (lastRecord.timestamp < record.timestamp) { vulDetailsMap.put(vulKey, record); } } else { vulDetailsMap.put(vulKey, record); } } else { def vulDetailsMap = new HashMap(); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; vulDetailsMap.put(vulKey, record); statusMap.put(key, vulDetailsMap); } } } } def resultList = new ArrayList(); DateTimeFormatter formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params .timeZone)); DateTimeFormatter twelveHourFormatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId .of(params.timeZone)); def is24HourFormat = params.timeFormat == '24h'; if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vulMapBranchLevel = statusMap.get(uniqueKey); def vulDetailsListBranchLevel = new ArrayList(); def vulMapOrgLevel = new HashMap(); def uniqueOccurrences = 0; def hasVulDetailsAtOrgLevel = false; def isVulResolved = false; Duration sumOfTTR = null; for (vulKey in vulMapBranchLevel.keySet()) { def vul = vulMapBranchLevel.get(vulKey); if (vul.bug_status == 'Resolved' && vul.date_of_discovery.getMillis() > 0) { isVulResolved = true; uniqueOccurrences++; Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); Instant resolutionDate = Instant.ofEpochMilli(vul.scan_time.getMillis()); Duration resolutionTime = Duration.between(startDate, resolutionDate); if (sumOfTTR != null) { sumOfTTR = sumOfTTR.plus(resolutionTime); } else { sumOfTTR = resolutionTime; } def days = resolutionTime.toDays(); def hours = resolutionTime.minusDays(days).toHours(); def mins = resolutionTime.minusDays(days).minusHours(hours).toMinutes(); def TTR = days + 'd ' + hours + 'h ' + mins + 'm'; def severityCode = 0; def curSeverity = vul.severity; if (curSeverity == 'VERY_HIGH') { vul.severity = 'Very high'; severityCode = 4; } else if (curSeverity == 'HIGH') { vul.severity = 'High'; severityCode = 3; } else if (curSeverity == 'MEDIUM') { vul.severity = 'Medium'; severityCode = 2; } else if (curSeverity == 'LOW') { vul.severity = 'Low'; severityCode = 1; } def diffAge = ChronoUnit.DAYS.between(startDate, resolutionDate); def SLAToolTip = ''; def statusToolTip = ''; statusToolTip = 'Date of resolution: ' + formatter.format(vul.scan_time); if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else { vul.sla = 'Within SLA'; } def map = [: ]; if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } if (statusToolTip != '') { map.put('statusToolTipContent', statusToolTip); } map.put('lastDiscovered', is24HourFormat ? formatter.format(vul.scan_time) : twelveHourFormatter .format(vul.scan_time)); map.put('component', vul.component_name); map.put('componentId', vul.component_id); map.put('branch', vul.branch); map.put('scannerName', vul.scanner_name); map.put('resolutionTime', TTR); map.put('sla', vul.sla); if (!hasVulDetailsAtOrgLevel) { vulMapOrgLevel.put('vulnerabilityId', vul.code); vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery); vulMapOrgLevel.put('vulnerabilityName', vul.name); vulMapOrgLevel.put('severity', vul.severity); vulMapOrgLevel.put('severityCode', severityCode); hasVulDetailsAtOrgLevel = true; } if (vul.date_of_discovery.getMillis() < vulMapOrgLevel.get('firstDiscovered').getMillis()) { vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery); } } } if (isVulResolved) { def size = vulMapBranchLevel.size(); Duration avg = sumOfTTR.dividedBy(size); def days = avg.toDays(); def hours = avg.minusDays(days).toHours(); def mins = avg.minusDays(days).minusHours(hours).toMinutes(); def avgResolutionTime = days + 'd ' + hours + 'h ' + mins + 'm'; def fd = vulMapOrgLevel.get('firstDiscovered'); vulMapOrgLevel.put('firstDiscovered', is24HourFormat ? formatter.format(fd) : twelveHourFormatter .format(fd)); vulMapOrgLevel.put('resolvedAreas', uniqueOccurrences); vulMapOrgLevel.put('averageResolutionTime', avgResolutionTime); vulMapOrgLevel.put('customSubRowsInfo', ['enabled': true, 'report_id': 'mttrForVulnerabilitiesSubRows', 'reportInfo': ['code': vulMapOrgLevel.get( 'vulnerabilityId')] ]); resultList.add(vulMapOrgLevel); } } } return resultList;"			}
		}
	}
  }`

const MTTRDrillDownSubRowsQuery = `{
	"_source": false,
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
          	"term": {
            	"code": "{{.vulCode}}"
          	}
          },
		  {
			"exists": {
			  "field": "date_of_discovery"
			}
		  },
		  {
			"terms": {
			  "severity": ["MEDIUM", "HIGH", "LOW", "VERY_HIGH"]
			}
		  }
		]
	  }
	},
	"aggs": {
	  "drilldowns": {
		"scripted_metric": {
			"params": {
				"timeZone": "{{.timeZone}}",
				"timeFormat": "{{.timeFormat}}"
			},
			"init_script": "state.statusMap = [: ];",
			"map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'recurrences':params['_source']['failure_files'].size(), 'scanner_name':doc.scanner_name.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value, 'scan_time':doc.scan_time.value, 'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':''];map.put(key, v);",
			"combine_script": "return state.statusMap;",
			"reduce_script": "def statusMap = new HashMap(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.code; if (statusMap.containsKey(key)) { def vulDetailsMap = statusMap.get(key); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; if (vulDetailsMap.containsKey(vulKey)) { def lastRecord = vulDetailsMap.get(vulKey); if (lastRecord.timestamp < record.timestamp) { vulDetailsMap.put(vulKey, record); } } else { vulDetailsMap.put(vulKey, record); } } else { def vulDetailsMap = new HashMap(); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; vulDetailsMap.put(vulKey, record); statusMap.put(key, vulDetailsMap); } } } } def resultList = new ArrayList(); DateTimeFormatter formatter; if (params.timeFormat == '12h') { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); } else { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); } if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vulMapBranchLevel = statusMap.get(uniqueKey); def vulDetailsListBranchLevel = new ArrayList(); def vulMapOrgLevel = new HashMap(); def hasVulDetailsAtOrgLevel = false; def isVulResolved = false; Duration sumOfTTR = null; for (vulKey in vulMapBranchLevel.keySet()) { def vul = vulMapBranchLevel.get(vulKey); if (vul.bug_status == 'Resolved' && vul.date_of_discovery.getMillis() > 0) { isVulResolved = true; Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); Instant resolutionDate = Instant.ofEpochMilli(vul.scan_time.getMillis()); Duration resolutionTime = Duration.between(startDate, resolutionDate); if (sumOfTTR != null) { sumOfTTR = sumOfTTR.plus(resolutionTime); } else { sumOfTTR = resolutionTime; } def days = resolutionTime.toDays(); def hours = resolutionTime.minusDays(days).toHours(); def mins = resolutionTime.minusDays(days).minusHours(hours).toMinutes(); def TTR = days + 'd ' + hours + 'h ' + mins + 'm'; def severityCode = 0; def curSeverity = vul.severity; if (curSeverity == 'VERY_HIGH') { vul.severity = 'Very high'; severityCode = 4; } else if (curSeverity == 'HIGH') { vul.severity = 'High'; severityCode = 3; } else if (curSeverity == 'MEDIUM') { vul.severity = 'Medium'; severityCode = 2; } else if (curSeverity == 'LOW') { vul.severity = 'Low'; severityCode = 1; } def diffAge = ChronoUnit.DAYS.between(startDate, resolutionDate); def SLAToolTip = ''; def statusToolTip = ''; statusToolTip = 'Date of resolution: ' + formatter.format(vul.scan_time); if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else { vul.sla = 'Within SLA'; } def map = [: ]; if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } if (statusToolTip != '') { map.put('statusToolTipContent', statusToolTip); } map.put('lastDiscovered', formatter.format(vul.scan_time)); map.put('component', vul.component_name); map.put('componentId', vul.component_id); map.put('branch', vul.branch); map.put('scannerName', vul.scanner_name); map.put('resolutionTime', TTR); map.put('sla', vul.sla); vulDetailsListBranchLevel.add(map); if (!hasVulDetailsAtOrgLevel) { vulMapOrgLevel.put('vulnerabilityId', vul.code); vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery); vulMapOrgLevel.put('vulnerabilityName', vul.name); vulMapOrgLevel.put('severity', vul.severity); vulMapOrgLevel.put('severityCode', severityCode); hasVulDetailsAtOrgLevel = true; } if (vul.date_of_discovery.getMillis() < vulMapOrgLevel.get('firstDiscovered').getMillis()) { vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery); } } } if (isVulResolved) { resultList.add(vulDetailsListBranchLevel); } } } return resultList[0];"
			}
    	}
    }
}`

const VelocityDrilldownQuery = `{
    "_source": false,
    "size": 0,
    "query": {
        "bool": {
            "filter": [{
                    "range": {
                        "completed_at": {
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
                        "deleted": "false"
                    }
                }
            ]
        }
    },
    "aggs": {
        "drilldowns": {
            "scripted_metric": {
				"params": {
					"timeZone": "{{.timeZone}}",
					"timeFormat": "{{.timeFormat}}"
				},
                "init_script": "state.statusMap = [:];",
                "map_script": "def map=state.statusMap;def key=doc.automation_id.value+'_'+doc.component_id.value+'_'+doc.flow_item.value+'_'+doc.issue_key.value+'_'+doc.run_id.value+'_'+doc.org_id.value;def v=['updatedAt':doc['updated_at'].getValue().toEpochSecond()*1000,'automation_id':doc.automation_id.value,'component_id':doc.component_id.value,'run_id':doc.run_id.value,'org_id':doc.org_id.value,'issue_key':doc.issue_key.value,'flow_item':doc.flow_item.value,'created_at':doc.created_at.value,'completed_at':doc['completed_at'].size()==0?null:doc.completed_at.value,'issue_url':doc['issue_url'].size()==0?'':doc.issue_url.value,'issue_type':doc['issue_type'].size()==0?'':doc.issue_type.value,'summary':doc['summary'].size()==0?'':doc.summary.value,'assignee':doc['assignee_name'].size()==0?'':doc.assignee_name.value];map.put(key,v);",
                "combine_script": "return state.statusMap;",
                "reduce_script": "def resultMap = new HashMap(); def valueList = new ArrayList(); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); if (record.flow_item.toUpperCase() == 'FEATURE') { record.flow_item = 'Feature'; } else if (record.flow_item.toUpperCase() == 'DEFECT') { record.flow_item = 'Defect'; } else if (record.flow_item.toUpperCase() == 'TECH_DEBT') { record.flow_item = 'Tech debt'; } else if (record.flow_item.toUpperCase() == 'RISK') { record.flow_item = 'Risk'; } def key = record.issue_key; if (resultMap.containsKey(key)) { def lastRecord = resultMap.get(key); if (lastRecord.updatedAt < record.updatedAt) { resultMap.put(key, record); } } else { resultMap.put(key, record); } } } } DateTimeFormatter formatterZoned; if (params.timeFormat == '12h') { formatterZoned = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); } else { formatterZoned = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); } for (key in resultMap.keySet()) { def record = resultMap.get(key); def v = ['issueId': record.issue_key, 'issueUrl': record.issue_url, 'issueType': record.issue_type, 'summary': record.summary, 'assignedTo': record.assignee, 'issueCompletedOn': formatterZoned .format(record.completed_at), 'flowItemType': record.flow_item ]; valueList.add(v) } return valueList;"
            }
        }
    }
}`

const DistributionDrilldownQuery = `{
    "_source": false,
    "size": 0,
    "query": {
        "bool": {
            "filter": [{
                    "range": {
                        "completed_at": {
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
                        "deleted": "false"
                    }
                }
            ]
        }
    },
    "aggs": {
        "drilldowns": {
            "scripted_metric": {
				"params": {
					"timeZone": "{{.timeZone}}",
					"timeFormat": "{{.timeFormat}}"
				},
                "init_script": "state.statusMap = [:];",
                "map_script": "def map=state.statusMap;def key=doc.automation_id.value+'_'+doc.component_id.value+'_'+doc.flow_item.value+'_'+doc.issue_key.value+'_'+doc.run_id.value+'_'+doc.org_id.value;def v=['updatedAt':doc['updated_at'].getValue().toEpochSecond()*1000,'automation_id':doc.automation_id.value,'component_id':doc.component_id.value,'run_id':doc.run_id.value,'org_id':doc.org_id.value,'issue_key':doc.issue_key.value,'flow_item':doc.flow_item.value,'created_at':doc.created_at.value,'completed_at':doc['completed_at'].size()==0?null:doc.completed_at.value,'issue_url':doc['issue_url'].size()==0?'':doc.issue_url.value,'issue_type':doc['issue_type'].size()==0?'':doc.issue_type.value,'summary':doc['summary'].size()==0?'':doc.summary.value,'assignee':doc['assignee_name'].size()==0?'':doc.assignee_name.value];map.put(key,v);",
                "combine_script": "return state.statusMap;",
                "reduce_script": "def resultMap = new HashMap(); def valueList = new ArrayList(); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); if (record.flow_item.toUpperCase() == 'FEATURE') { record.flow_item = 'Feature'; } else if (record.flow_item.toUpperCase() == 'DEFECT') { record.flow_item = 'Defect'; } else if (record.flow_item.toUpperCase() == 'TECH_DEBT') { record.flow_item = 'Tech debt'; } else if (record.flow_item.toUpperCase() == 'RISK') { record.flow_item = 'Risk'; } def key = record.issue_key; if (resultMap.containsKey(key)) { def lastRecord = resultMap.get(key); if (lastRecord.updatedAt < record.updatedAt) { resultMap.put(key, record); } } else { resultMap.put(key, record); } } } } DateTimeFormatter formatterZoned; if (params.timeFormat == '12h') { formatterZoned = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); } else { formatterZoned = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); } for (key in resultMap.keySet()) { def record = resultMap.get(key); def v = ['issueId': record.issue_key, 'issueUrl': record.issue_url, 'issueType': record.issue_type, 'summary': record.summary, 'assignedTo': record.assignee, 'issueCompletedOn': formatterZoned .format(record.completed_at), 'flowItemType': record.flow_item ]; valueList.add(v) } return valueList;"
            }
        }
    }
}`

const CycleTimeDrilldownQuery = `{
    "_source": false,
    "size": 0,
    "query": {
        "bool": {
            "filter": [{
                    "range": {
                        "completed_at": {
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
                        "deleted": "false"
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
                "init_script": "state.statusMap = [:];",
                "map_script": "def map=state.statusMap;def key=doc.automation_id.value+'_'+doc.component_id.value+'_'+doc.flow_item.value+'_'+doc.issue_key.value+'_'+doc.run_id.value+'_'+doc.org_id.value;def v=['updatedAt':doc['updated_at'].getValue().toEpochSecond()*1000,'automation_id':doc.automation_id.value,'component_id':doc.component_id.value,'run_id':doc.run_id.value,'org_id':doc.org_id.value,'issue_key':doc.issue_key.value,'flow_item':doc.flow_item.value,'created_at':doc.created_at.value,'completed_at':doc['completed_at'].size()==0?null:doc.completed_at.value,'issue_url':doc['issue_url'].size()==0?'':doc.issue_url.value,'issue_type':doc['issue_type'].size()==0?'':doc.issue_type.value,'summary':doc['summary'].size()==0?'':doc.summary.value,'flow_time':doc['flow_time'].size()==0?0:doc.flow_time.value,'assignee':doc['assignee_name'].size()==0?'':doc.assignee_name.value];map.put(key,v);",
                "combine_script": "return state.statusMap;",
                "reduce_script": "def resultMap = new HashMap();def valueList = new ArrayList();for (a in states){if (a != null){for (i in a.keySet()){def record = a.get(i);if (record.flow_item.toUpperCase() == 'FEATURE'){record.flow_item = 'Feature';} else if (record.flow_item.toUpperCase() == 'DEFECT'){record.flow_item = 'Defect';} else if (record.flow_item.toUpperCase() == 'TECH_DEBT'){record.flow_item = 'Tech debt';} else if (record.flow_item.toUpperCase() == 'RISK'){record.flow_item = 'Risk';}def key = record.issue_key;if (resultMap.containsKey(key)){def lastRecord = resultMap.get(key);if (lastRecord.updatedAt < record.updatedAt){resultMap.put(key, record);}} else{resultMap.put(key, record);}}}}DateTimeFormatter formatterZoned = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm').withZone(ZoneId.of(params.timeZone));for (key in resultMap.keySet()){def record = resultMap.get(key);def v = ['issueId':record.issue_key, 'issueUrl':record.issue_url, 'issueType':record.issue_type, 'summary':record.summary, 'assignedTo':record.assignee, 'issueCompletedOn':formatterZoned.format(record.completed_at), 'flowItemType':record.flow_item, 'flowTime':record.flow_time];valueList.add(v)}return valueList;"
            }
        }
    }
}`

const WorkEfficiencyDrilldownQuery = `{
    "_source": false,
    "size": 0,
    "query": {
        "bool": {
            "filter": [{
                    "range": {
                        "completed_at": {
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
                        "deleted": "false"
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
                "init_script": "state.statusMap = [:];",
                "map_script": "def map=state.statusMap;def key=doc.automation_id.value+'_'+doc.component_id.value+'_'+doc.flow_item.value+'_'+doc.issue_key.value+'_'+doc.run_id.value+'_'+doc.org_id.value;def v=['updatedAt':doc['updated_at'].getValue().toEpochSecond()*1000,'automation_id':doc.automation_id.value,'component_id':doc.component_id.value,'run_id':doc.run_id.value,'org_id':doc.org_id.value,'issue_key':doc.issue_key.value,'flow_item':doc.flow_item.value,'created_at':doc.created_at.value,'completed_at':doc['completed_at'].size()==0?null:doc.completed_at.value,'issue_url':doc['issue_url'].size()==0?'':doc.issue_url.value,'issue_type':doc['issue_type'].size()==0?'':doc.issue_type.value,'summary':doc['summary'].size()==0?'':doc.summary.value,'active_time':doc['active_time'].size()==0?0:doc.active_time.value,'waiting_time':doc['waiting_time'].size()==0?0:doc.waiting_time.value,'flow_efficiency':doc['flow_efficiency'].size()==0?0:doc.flow_efficiency.value,'assignee':doc['assignee_name'].size()==0?'':doc.assignee_name.value];map.put(key,v);",
                "combine_script": "return state.statusMap;",
                "reduce_script": "def resultMap = new HashMap();def valueList = new ArrayList();for (a in states){if (a != null){for (i in a.keySet()){def record = a.get(i);if (record.flow_item.toUpperCase() == 'FEATURE'){record.flow_item = 'Feature';} else if (record.flow_item.toUpperCase() == 'DEFECT'){record.flow_item = 'Defect';} else if (record.flow_item.toUpperCase() == 'TECH_DEBT'){record.flow_item = 'Tech debt';} else if (record.flow_item.toUpperCase() == 'RISK'){record.flow_item = 'Risk';}def key = record.issue_key;if (resultMap.containsKey(key)){def lastRecord = resultMap.get(key);if (lastRecord.updatedAt < record.updatedAt){resultMap.put(key, record);}} else{resultMap.put(key, record);}}}}DateTimeFormatter formatterZoned = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm').withZone(ZoneId.of(params.timeZone));for (key in resultMap.keySet()){def record = resultMap.get(key);def v = ['issueId':record.issue_key, 'issueUrl':record.issue_url, 'issueType':record.issue_type, 'summary':record.summary, 'assignedTo':record.assignee, 'issueCompletedOn':formatterZoned.format(record.completed_at), 'flowItemType':record.flow_item, 'activeTime':record.active_time, 'waitingTime':record.waiting_time, 'efficiency':record.flow_efficiency];valueList.add(v)}return valueList;"
            }
        }
    }
}`

const WorkLoadDrilldownQuery = `{
	"_source": false,
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
			  "deleted": "false"
			}
		  }
		]
	  }
	},
	"aggs": {
	  "drilldowns": {
		"scripted_metric": {
		  "params": {
			"startDate": "{{.dateHistogramMin}}",
			"endDate": "{{.dateHistogramMax}}",
			"aggrBy": "{{.aggrBy}}",
			"timeZone": "{{.timeZone}}",
			"timeFormat": "{{.timeFormat}}"
		  },
		  "init_script": "state.uniqueIssuesMap = [:];",
		  "map_script": "def map = state.uniqueIssuesMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.automation_id.value + '_' + doc.run_id.value + '_' + doc.issue_key.value + doc.flow_item.value;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'created_at':doc['created_at'].getValue().toEpochSecond() * 1000, 'updated_at':doc['updated_at'].getValue().toEpochSecond() * 1000, 'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'issue_key':doc.issue_key.value, 'flow_item':doc.flow_item.value, 'completed_at':doc['completed_at'].size() == 0 ? null:doc['completed_at'].getValue().toEpochSecond() * 1000,'clh':params['_source']['change_log_history'],'completed_at_date': doc['completed_at'].size() == 0 ? null : doc.completed_at.value, 'created_at_date': doc['created_at'].size() == 0 ? null : doc.created_at.value, 'issue_url': doc['issue_url'].size() == 0 ? '' : doc.issue_url.value, 'issue_type': doc['issue_type'].size() == 0 ? '' : doc.issue_type.value, 'summary': doc['summary'].size() == 0 ? '' : doc.summary.value, 'assignee': doc['assignee_name'].size() == 0 ? '' : doc.assignee_name.value];map.put(key, v);",
		  "combine_script": "return state.uniqueIssuesMap;",
		  "reduce_script": "HashMap getDateBuckets(LocalDate startDate, LocalDate endDate, String aggrBy) { def OutputMap = new HashMap(); def dateIntervalsMap = new HashMap(); TreeMap dates = new TreeMap(); if (aggrBy.equals('week')) { def prevDate = startDate; dates.put(startDate.toString(), new LinkedHashMap()); LocalDate firstMonday = startDate.with(TemporalAdjusters.next(DayOfWeek.MONDAY)); while (firstMonday.compareTo(endDate) <= 0) { def curDate = firstMonday; long diffMillis = Duration.between(prevDate.atStartOfDay(), curDate.atStartOfDay()) .toMillis(); dateIntervalsMap.put(prevDate.toString(), diffMillis); prevDate = curDate; dates.put(firstMonday.toString(), new LinkedHashMap()); firstMonday = firstMonday.plusDays(7); } long diffForLastDate = Duration.between(prevDate.atStartOfDay(), endDate.atStartOfDay()) .toMillis(); dateIntervalsMap.put(prevDate.toString(), diffForLastDate); } else if (aggrBy.equals('day')) { long days = ChronoUnit.DAYS.between(startDate, endDate); long oneDayInMilli = 86400000; for (long i = 0; i <= days; i++) { LocalDate date = startDate.plusDays(i); dates.put(date.toString(), new LinkedHashMap()); dateIntervalsMap.put(date.toString(), oneDayInMilli); } } else if (aggrBy.equals('month')) { def prevDate = startDate; dates.put(startDate.toString(), new LinkedHashMap()); LocalDate firstDayOfMonth = startDate.with(TemporalAdjusters.firstDayOfMonth()); if (startDate.isAfter(firstDayOfMonth)) { firstDayOfMonth = firstDayOfMonth.plusMonths(1); } while (firstDayOfMonth.compareTo(endDate) <= 0) { def curDate = firstDayOfMonth; long diffMillis = Duration.between(prevDate.atStartOfDay(), curDate.atStartOfDay()) .toMillis(); dateIntervalsMap.put(prevDate.toString(), diffMillis); prevDate = curDate; dates.put(firstDayOfMonth.toString(), new LinkedHashMap()); firstDayOfMonth = firstDayOfMonth.plusMonths(1); } long diffForLastDate = Duration.between(prevDate.atStartOfDay(), endDate.atStartOfDay()) .toMillis(); dateIntervalsMap.put(prevDate.toString(), diffForLastDate); } OutputMap.put('dates', dates); OutputMap.put('intervals', dateIntervalsMap); return OutputMap } def issuesMap = new HashMap(); def workLoadIssuesTotal = new HashSet(); def valueList = new ArrayList(); def formatterISO = DateTimeFormatter.ISO_LOCAL_DATE; LocalDate stDt = LocalDate.parse(params.startDate, formatterISO); LocalDate endDt = LocalDate.parse(params.endDate, formatterISO); def dateInfoMap = getDateBuckets(stDt, endDt, params.aggrBy); def dateBuckets = dateInfoMap.get('dates'); def dateIntervals = dateInfoMap.get('intervals'); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); if (record.flow_item.toUpperCase() == 'FEATURE') { record.flow_item = 'Feature'; } else if (record.flow_item.toUpperCase() == 'DEFECT') { record.flow_item = 'Defect'; } else if (record.flow_item.toUpperCase() == 'TECH_DEBT') { record.flow_item = 'Tech debt'; } else if (record.flow_item.toUpperCase() == 'RISK') { record.flow_item = 'Risk'; } def key = record.org_id + '_' + record.issue_key; if (issuesMap.containsKey(key)) { def lastRecord = issuesMap.get(key); if (lastRecord.updated_at < record.updated_at) { issuesMap.put(key, record); } } else { issuesMap.put(key, record); } } } } SimpleDateFormat dfDateTime = new SimpleDateFormat('yyyy-MM-dd HH:mm:ss'); DateTimeFormatter formatterZoned; if (params.timeFormat == '12h') { formatterZoned = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); } else { formatterZoned = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); } for (ele in issuesMap.keySet()) { def curIssue = issuesMap.get(ele); def clh = curIssue.clh; def statusTimestamps = clh.status_timestamp; if (statusTimestamps != null) { for (int i = 0; i < statusTimestamps.size(); i++) { def curStatusTimestampMap = statusTimestamps[i]; if (curStatusTimestampMap.get('type') == 'IN_PROGRESS') { def wipStartTimeString = curStatusTimestampMap.get('start_time'); def wipEndTimeString = curStatusTimestampMap.get('end_time'); Date wipST = dfDateTime.parse(wipStartTimeString); long wipStartTime = wipST.getTime(); Date wipET; long wipEndTime; if (wipEndTimeString != '') { wipET = dfDateTime.parse(wipEndTimeString); wipEndTime = wipET.getTime(); } def issueFlowItem = curIssue.flow_item; def issueKey = curIssue.issue_key; for (date in dateBuckets.keySet()) { LocalDate bucketDate = LocalDate.parse(date, formatterISO); ZonedDateTime zonedBucketDate = bucketDate.atStartOfDay(ZoneId.of(params.timeZone)); Instant zonedBucketDateInstant = zonedBucketDate.toInstant(); long bucketDateEpoch = zonedBucketDateInstant.toEpochMilli(); if ((wipStartTime < bucketDateEpoch && (wipEndTimeString == '' || (wipEndTime >= bucketDateEpoch))) || (wipStartTime >= bucketDateEpoch && wipStartTime < ( bucketDateEpoch + dateIntervals.get(date)))) { if (!workLoadIssuesTotal.contains(issueKey)) { workLoadIssuesTotal.add(issueKey); def v = ['issueId': curIssue.issue_key, 'issueUrl': curIssue.issue_url, 'issueType': curIssue.issue_type, 'flowItemType': curIssue.flow_item, 'summary': curIssue.summary, 'assignedTo': curIssue.assignee, 'issueCompletedOn': curIssue.completed_at_date == null ? '' : formatterZoned.format(curIssue.completed_at_date), 'issueCreatedOn': curIssue.created_at_date == null ? '' : formatterZoned.format(curIssue .created_at_date) ]; valueList.add(v); } } } } } } } return valueList;"
		}
	  }
	}
  }`

const DeploymentFrequencyAndLeadTime = `{
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
                    "term": {
                        "data_type": 2
                    }
                },
                {
                    "term": {
                        "target_env": "{{.targetEnv}}"
                    }
                }
            ]
        }
    },
    "aggs": {
        "deployments": {
            "scripted_metric": {
                "params": {
                    "timeZone": "{{.timeZone}}"
                },
                "combine_script": "return state.data_map;",
                "init_script": "state.data_map=[:];",
                "map_script": "def map = state.data_map, runStartTime = 0, runStartTimeZonedString = '';def key = doc.component_id.value + '_' + doc.run_id.value + '_' + doc.job_id.value + '_' + doc.step_id.value + '_' + doc.target_env.value + '_' + doc.status.value;def v = ['run_id':doc.run_id.value, 'run_number':doc.run_number.value, 'job_id':doc.job_id.value, 'step_id':doc.step_id.value, 'target_env':doc.target_env.value, 'step_kind':doc.step_kind.value, 'status_timestamp':doc.status_timestamp.value, 'status':doc.status.value, 'component_id':doc.component_id.value, 'automation_id':doc.automation_id.value, 'workflow_name':doc.workflow_name.value, 'component_name':doc.component_name.value, 'status_timestamp_zoned':doc.status_timestamp.value.withZoneSameInstant(ZoneId.of(params.timeZone))];if (doc['run_start_time'].size() != 0){runStartTime = doc.run_start_time.value;ZonedDateTime zdt = Instant.ofEpochMilli(runStartTime).atZone(ZoneId.of(params.timeZone));runStartTimeZonedString = zdt.format(DateTimeFormatter.ofPattern('yyyy-MM-dd HH:mm:ss'))}v['run_start_time'] = runStartTime;v['run_start_time_string_zoned'] = runStartTimeZonedString;map.put(key, v);",
                "reduce_script": "def allDataMap = [:], resultArray = [], jobStepDedupMap = new HashMap();for (response in states){if (response != null){for (key in response.keySet()){allDataMap.put(key, response.get(key));}}}for (key in allDataMap.keySet()){def currRecord = allDataMap.get(key);if (currRecord.step_id == ''){jobStepDedupMap.put(key, currRecord);} else{def jobLevelRecordKey = currRecord.component_id + '_' + currRecord.run_id + '_' + currRecord.job_id + '_' + '' + '_' + currRecord.target_env + '_' + currRecord.status;if (!allDataMap.containsKey(jobLevelRecordKey)){jobStepDedupMap.put(key, currRecord);}}}for (key in jobStepDedupMap.keySet()){resultArray.add(jobStepDedupMap.get(key))}return resultArray;"
            }
        }
    }
}`

const FailureRate = `{
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
                        "data_type": 2
                    }
                },
				{
					"term": {
						"target_env": "{{.targetEnv}}"
					}
				}
			]
		}
	},
	"aggs": {
		"deployments": {
		  "scripted_metric": {
			"combine_script": "return state.data_map;",
			"init_script": "state.data_map=[:];",
			"map_script": "def map = state.data_map,runStartTime=0;def key = doc.component_id.value + '_' + doc.run_id.value + '_' + doc.job_id.value + '_' + doc.step_id.value + '_' + doc.target_env.value + '_' + doc.status.value;def v = ['run_id': doc.run_id.value,'run_number': doc.run_number.value, 'job_id': doc.job_id.value, 'step_id': doc.step_id.value, 'target_env': doc.target_env.value, 'step_kind': doc.step_kind.value,  'status_timestamp': doc.status_timestamp.value, 'status': doc.status.value, 'component_id': doc.component_id.value,'automation_id':doc.automation_id.value,'workflow_name':doc.workflow_name.value,'component_name':doc.component_name.value];if (doc['run_start_time'].size() != 0) {runStartTime=doc.run_start_time.value;}v['run_start_time']=runStartTime;map.put(key, v);",
			"reduce_script": "def allDataMap = [:], componentMap = new HashMap(), jobStepDedupMap = new HashMap();for (response in states){if (response != null){for (key in response.keySet()){allDataMap.put(key, response.get(key));}}}for (key in allDataMap.keySet()){def currRecord = allDataMap.get(key);if (currRecord.step_id == ''){jobStepDedupMap.put(key, currRecord);} else{def jobLevelRecordKey = currRecord.component_id + '_' + currRecord.run_id + '_' + currRecord.job_id + '_' + '' + '_' + currRecord.target_env + '_' + currRecord.status;if (!allDataMap.containsKey(jobLevelRecordKey)){jobStepDedupMap.put(key, currRecord);}}}for (key in jobStepDedupMap.keySet()){def record = jobStepDedupMap.get(key);if (componentMap.containsKey(record.component_id)){def infoMap = componentMap.get(record.component_id);infoMap['deployments'] += 1;if (record.status == 'SUCCEEDED'){infoMap['success'] += 1;} else{infoMap['failure'] += 1;}componentMap.put(record.component_id, infoMap);} else{def infoMap = new HashMap();infoMap['component_id'] = record.component_id;infoMap['component_name'] = record.component_name;infoMap['deployments'] = 1;if (record.status == 'SUCCEEDED'){infoMap['success'] = 1;infoMap['failure'] = 0;} else if (record.status == 'FAILED' || record.status == 'ABORTED' || record.status == 'TIMED_OUT'){infoMap['success'] = 0;infoMap['failure'] = 1;}componentMap.put(record.component_id, infoMap);}}return componentMap;"
		  }
		}
	}
}`

const DoraMetricsMttr = `{
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
                        "data_type": 2
                    }
                },
				{
					"term": {
						"target_env": "{{.targetEnv}}"
					}
				}
			]
		}
	},
	"aggs": {
		"deployments": {
		  "scripted_metric": {
			"combine_script": "return state.data_map;",
			"init_script": "state.data_map=[:];",
			"map_script": "def map = state.data_map,runStartTime=0;def key = doc.component_id.value + '_' + doc.run_id.value + '_' + doc.job_id.value + '_' + doc.step_id.value + '_' + doc.target_env.value + '_' + doc.status.value;def v = ['run_id': doc.run_id.value,'run_number': doc.run_number.value, 'job_id': doc.job_id.value, 'step_id': doc.step_id.value, 'target_env': doc.target_env.value, 'step_kind': doc.step_kind.value,  'status_timestamp': doc.status_timestamp.value, 'status': doc.status.value, 'component_id': doc.component_id.value,'automation_id':doc.automation_id.value,'workflow_name':doc.workflow_name.value,'component_name':doc.component_name.value];if (doc['run_start_time'].size() != 0) {runStartTime=doc.run_start_time.value;}v['run_start_time']=runStartTime;map.put(key, v);",
			"reduce_script": "def allDataMap = [:], componentMap = new HashMap(), jobStepDedupMap = new HashMap();for (response in states){if (response != null){for (key in response.keySet()){allDataMap.put(key, response.get(key));}}}for (key in allDataMap.keySet()){def currRecord = allDataMap.get(key);if (currRecord.step_id == ''){jobStepDedupMap.put(key, currRecord);} else{def jobLevelRecordKey = currRecord.component_id + '_' + currRecord.run_id + '_' + currRecord.job_id + '_' + '' + '_' + currRecord.target_env + '_' + currRecord.status;if (!allDataMap.containsKey(jobLevelRecordKey)){jobStepDedupMap.put(key, currRecord);}}}for (key in jobStepDedupMap.keySet()){def record = jobStepDedupMap.get(key);def compArray = new ArrayList();if (componentMap.containsKey(record.component_id)){compArray = componentMap.get(record.component_id);}compArray.add(record);componentMap.put(record.component_id, compArray)}def recoveredCount = 0.0, recoveredTotalDuration = 0.0, resultArray = new ArrayList();for (key in componentMap.keySet()){def componentArray = componentMap.get(key);for (def i = 0; i < componentArray.size(); i++){for (def j = i + 1; j < componentArray.size(); j++){if (componentArray[i].status_timestamp.getMillis() > componentArray[j].status_timestamp.getMillis()){def temp = componentArray[i];componentArray[i] = componentArray[j];componentArray[j] = temp;}}}def failedTime = 0, failedRunId, failedRunNumber;for (def i = 0; i < componentArray.size(); i++){def status = componentArray[i].status;if ((status == 'FAILED' || status == 'TIMED_OUT' || status == 'ABORTED') && failedTime == 0){failedTime = componentArray[i].status_timestamp.getMillis();failedRunId = componentArray[i].run_id;failedRunNumber = componentArray[i].run_number;}if (failedTime != 0 && status == 'SUCCEEDED' && failedRunNumber < componentArray[i].run_number){def successTime = componentArray[i].status_timestamp.getMillis();def result = new HashMap();result.put('component_id', componentArray[i].component_id);result.put('component_name', componentArray[i].component_name);result.put('failed_run', failedRunId);result.put('failed_on', failedTime);result.put('recovered_run', componentArray[i].run_id);result.put('recovered_on', successTime);result.put('recovered_duration', successTime - failedTime);result.put('failed_run_number', failedRunNumber);result.put('recovered_run_number', componentArray[i].run_number);recoveredTotalDuration += (successTime - failedTime);recoveredCount += 1;failedTime = 0;resultArray.add(result);}}}return resultArray;"
		  }
		}
	}
}`

const TrivyLicenseDrilldownQuery = `{
    "_source": false,
    "size": 0,
    "query": {
        "bool": {
          "filter": [
            {
          "range": {
            "timestamp": {
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
            "component_id": "{{.componentIdForNestedDrillDown}}"
          } 
        },
        {
          "term": {
            "github_branch": "{{.branch}}"
          } 
        },
        {
          "term": {
            "standard": "LICENSE"
          } 
        },
        {
          "term": {
            "license_type": "{{.licenseType}}"
          } 
        },
        {
          "term": {
            "scanner_name": "trivy"
          } 
        },
        {
          "term": {
            "run_id": "{{.runId}}"
          } 
        },
            {
              "exists": {
                "field": "date_of_discovery"
              }
            },
            {
              "terms": {
                "severity": ["MEDIUM", "HIGH", "LOW", "VERY_HIGH"]
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
            "map_script": "def map = state.statusMap;def key=doc.org_id.value+'_'+doc.component_id.value+ '_' +doc.github_branch.value +'_'+doc.code.value;def v=['org_id': doc.org_id.value, 'component_id': doc.component_id.value,'branch': doc.github_branch.value,'scanner_name': doc.scanner_name.value,'code':doc.code.value,'current_sub_row_failures':params['_source']['failure_files'],'date_of_discovery':doc.date_of_discovery.value];map.put(key,v);",
            "combine_script": "return state.statusMap;",
            "reduce_script": "def locationField = 'FilePath';def packageLabel = 'PkgName';DateTimeFormatter formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone));def resultArray = new ArrayList();for (a in states){if (a.size() > 0){for ( key in a.keySet()){def tempMap = a.get(key);def date_of_dis = tempMap.get('date_of_discovery');def failureArray = tempMap.get('current_sub_row_failures');for (int i = 0; i < failureArray.size(); i++){def failures = new HashMap();def curFailureLocation;def filePath = failureArray[i].get(locationField);if (filePath != null && filePath != ''){curFailureLocation = filePath;} else{curFailureLocation = 'No Failure Location Data Found';}def pkgName = failureArray[i].get(packageLabel);def curPackage;if (pkgName != null && pkgName != ''){curPackage = pkgName;} else{curPackage = 'No Package Name Found';}failures.put('usedSince', formatter.format(date_of_dis));failures.put('locations', curFailureLocation);failures.put('packageName', curPackage);resultArray.add(failures);}}}}return resultArray;"
		}
        }
      }
  }`

const TotalTestCasesDrilldownQuery = `{
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
                            "gte": "{{.startDate}}",
                            "lte": "{{.endDate}}",
                            "format": "yyyy-MM-dd HH:mm:ss",
                            "time_zone": "{{.timeZone}}"
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
                        "branch_id": "{{.branch}}"
                    }
                },
				{
                    "term":
                    {
                        "automation_id": "{{.automationId}}"
                    }
                }
            ]
        }
    },
    "aggs":
    {
        "drilldowns":
        {
            "scripted_metric":
            {
                "combine_script": "return state.dataMap;",
                "init_script": "state.dataMap = [:];",
                "map_script": "def map=state.dataMap;def key=doc.org_id.value+'_'+doc.component_id.value+'_'+doc.automation_id.value+'_'+doc.branch_id.value+'_'+doc.test_suite_name.value+'_'+doc.test_case_name.value+'_'+doc.start_time.getValue().toEpochSecond()*1000;def v=['org_id':doc.org_id.value,'component_id':doc.component_id.value,'component_name':doc.component_name.value,'automation_id':doc.automation_id.value,'automation_name':doc.automation_name.value,'branch_id':doc.branch_id.value,'branch_name':doc.branch_name.value,'run_id':doc.run_id.value,'test_suite_name':doc.test_suite_name.value,'test_case_name':doc.test_case_name.value,'start_time':doc.start_time.value,'status':doc.status.value,'duration':doc.duration.value,'start_time_in_millis':doc['start_time'].getValue().toEpochSecond()*1000,'duration_in_millis':doc['duration'].getValue()];map.put(key,v);",
                "reduce_script": "def resultMap=new HashMap();def resultList=new ArrayList();for(a in states){if(a!=null&&a.size()>0){for(i in a.keySet()){def record=a.get(i);def key=record.component_id+'_'+record.automation_id+'_'+record.branch_id+'_'+record.test_case_name+'_'+record.test_suite_name;if(resultMap.containsKey(key)){def lastRecord=resultMap.get(key);if(lastRecord.start_time_in_millis<record.start_time_in_millis){resultMap.put(key,record);}record.put('failure_count',lastRecord.get('failure_count'));record.put('success_count',lastRecord.get('success_count'));record.put('skipped_count',lastRecord.get('skipped_count'));record.put('total_exec_count',lastRecord.get('total_exec_count'));record.put('failure_rate',lastRecord.get('failure_rate'));def runsCount=lastRecord.get('runs');runsCount=runsCount+1;record.put('runs',runsCount);double totalDurationInMillis=lastRecord.get('total_duration_in_millis');double durationInMillis=record.get('duration_in_millis');record.put('total_duration_in_millis',totalDurationInMillis+durationInMillis);double result=(totalDurationInMillis+durationInMillis)/runsCount;result=Math.round(result*100)/100.0;record.put('average_duration',result);def successCount=record.get('success_count');def failureCount=record.get('failure_count');def skippedCount=record.get('skipped_count');def totalExecCount=record.get('total_exec_count');if(record.get('status')=='FAILED'){failureCount=failureCount+1;record.put('failure_count',failureCount);}else if(record.get('status')=='PASSED'){successCount=successCount+1;record.put('success_count',successCount);}else{skippedCount=skippedCount+1;record.put('skipped_count',skippedCount);}totalExecCount=totalExecCount+1;record.put('total_exec_count',totalExecCount);double failureRate=Math.round((failureCount*100)/totalExecCount);record.put('failure_rate',failureRate);resultMap.put(key,record);}else{record.put('runs',1);record.put('total_duration_in_millis',record.get('duration_in_millis'));record.put('average_duration',record.get('duration_in_millis'));if(record.get('status')=='FAILED'){record.put('failure_count',1);record.put('success_count',0);record.put('skipped_count',0);record.put('failure_rate',100.0);}else if(record.get('status')=='PASSED'){record.put('success_count',1);record.put('failure_count',0);record.put('skipped_count',0);record.put('failure_rate',0.0);}else{record.put('success_count',0);record.put('failure_count',0);record.put('skipped_count',1);record.put('failure_rate',0.0);}record.put('total_exec_count',1);resultMap.put(key,record);}}}}for(case in resultMap.keySet()){def tempMap=resultMap.get(case);def testCaseMap=new HashMap();testCaseMap.put('testName',tempMap.get('test_case_name'));double avgRunDuration=tempMap.get('average_duration');if(avgRunDuration<10){avgRunDuration=10;}testCaseMap.put('avgRunDuration',avgRunDuration);testCaseMap.put('totalRuns',tempMap.get('runs'));testCaseMap.put('branchName',tempMap.get('branch_name'));def failureRateMap=new HashMap();failureRateMap.put('value',tempMap.get('failure_rate'));def successfulRunsMap=new HashMap();successfulRunsMap.put('title','Successful runs');successfulRunsMap.put('value',tempMap.get('success_count'));def failedRunsMap=new HashMap();failedRunsMap.put('title','Failed runs');failedRunsMap.put('value',tempMap.get('failure_count'));def skippedRunsMap=new HashMap();skippedRunsMap.put('title','Skipped runs');skippedRunsMap.put('value',tempMap.get('skipped_count'));failureRateMap.put('type','SINGLE_BAR');def colorSchemeList=new ArrayList();def colorMap1=new HashMap();colorMap1.put('color0','#009C5B');colorMap1.put('color1','#62CA9D');def colorMap2=new HashMap();colorMap2.put('color0','#D32227');colorMap2.put('color1','#FB6E72');def colorMap3=new HashMap();colorMap3.put('color0','#F2A414');colorMap3.put('color1','#FFE6C1');colorSchemeList.add(colorMap1);colorSchemeList.add(colorMap2);colorSchemeList.add(colorMap3);failureRateMap.put('colorScheme',colorSchemeList);def lightColorSchemeList=new ArrayList();def lightColorMap1=new HashMap();lightColorMap1.put('color0','#0C9E61');lightColorMap1.put('color1','#79CAA8');def lightColorMap2=new HashMap();lightColorMap2.put('color0','#E83D39');lightColorMap2.put('color1','#F39492');def lightColorMap3=new HashMap();lightColorMap3.put('color0','#F2A414');lightColorMap3.put('color1','#FFE6C1');lightColorSchemeList.add(lightColorMap1);lightColorSchemeList.add(lightColorMap2);lightColorSchemeList.add(lightColorMap3);failureRateMap.put('lightColorScheme',lightColorSchemeList);def dataList=new ArrayList();dataList.add(successfulRunsMap);dataList.add(failedRunsMap);dataList.add(skippedRunsMap);failureRateMap.put('data',dataList);testCaseMap.put('failureRate',failureRateMap);testCaseMap.put('failureRateValue',tempMap.get('failure_rate'));def drilldownMap=new HashMap();drilldownMap.put('reportId','test-overview-view-run-activity');drilldownMap.put('reportTitle','Test case activity - '+tempMap.get('test_case_name'));def reportInfoMap=new HashMap();reportInfoMap.put('component_id',tempMap.get('component_id'));reportInfoMap.put('test_case_name',tempMap.get('test_case_name'));reportInfoMap.put('test_suite_name',tempMap.get('test_suite_name'));reportInfoMap.put('branch',tempMap.get('branch_id'));reportInfoMap.put('automation_id',tempMap.get('automation_id'));drilldownMap.put('reportInfo',reportInfoMap);def viewRunActivityMap=new HashMap();viewRunActivityMap.put('drillDown',drilldownMap);testCaseMap.put('viewRunActivity',viewRunActivityMap);resultList.add(testCaseMap)}resultList.sort((a,b)->{def aValue=a.get('failureRateValue');def bValue=b.get('failureRateValue');return Double.compare(bValue,aValue)});return resultList;"
            }
        }
    }
}`

const TestOverviewTotalRunsDrilldownQuery = `{
    "size": 0,
    "query": {
        "bool": {
            "filter": [
                {
                    "range": {
                        "start_time": {
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
                        "test_suite_name": "{{.testSuiteName}}"
                    }
                },
                {
                    "term": {
                        "branch_id": "{{.branch}}"
                    }
                },
                {
                    "term": {
                        "automation_id": "{{.automationId}}"
                    }
                }
            ]
        }
    },
    "aggs": {
        "runs": {
            "terms": {
                "field": "run_id",
                "size": 65000
            },
            "aggs": {
                "failed_docs": {
                    "filter": {
                        "term": {
                            "status": "FAILED"
                        }
                    },
                    "aggs": {
                        "failed_test_cases": {
                            "terms": {
                                "script": "doc['test_suite_name'].value + '_' + doc['test_case_name'].value",
                                "size": 65000
                            }
                        }
                    }
                },
                "total_test_cases": {
                    "terms": {
                        "script": "doc['test_suite_name'].value + '_' + doc['test_case_name'].value",
                        "size": 65000
                    }
                },
                "total_test_cases_count": {
                    "bucket_script": {
                        "buckets_path": {
                            "count": "total_test_cases._bucket_count"
                        },
                        "script": "params.count"
                    }
                },
                "run_details": {
                    "top_hits": {
                        "size": 1,
                        "_source": [
                            "org_id",
                            "component_id",
                            "automation_id",
                            "branch_id",
                            "run_id",
                            "run_number",
                            "run_status"
                        ],
                        "script_fields": {
                            "zoned_run_start_time": {
                                "script": {
                                    "params": {
                                        "timeZone": "{{.timeZone}}",
										"timeFormat": "{{.timeFormat}}"
                                    },
                                    "source": "def is24HourFormat = params.timeFormat == '24h'; DateTimeFormatter formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss') .withZone(ZoneId.of(params.timeZone)); DateTimeFormatter twelveHourFormatter = DateTimeFormatter.ofPattern( 'yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); if (is24HourFormat) { return formatter.format(doc['run_start_time'].value); } else { return twelveHourFormatter.format(doc['run_start_time'].value); }"
                                }
                            }
                        }
                    }
                }
            }
        },
        "test_cases_that_failed_at_least_once": {
            "terms": {
                "script": "doc['test_suite_name'].value + '_' + doc['test_case_name'].value",
                "size": 65000
            },
            "aggs": {
                "failed_docs": {
                    "filter": {
                        "term": {
                            "status": "FAILED"
                        }
                    }
                },
                "test_case_status_history": {
                    "terms": {
                        "field": "status",
                        "size": 65000
                    }
                },
                "test_case_name": {
                    "top_hits": {
                        "size": 1,
                        "_source": [
                            "test_case_name",
                            "test_suite_name"
                        ]
                    }
                },
                "failed_buckets": {
                    "bucket_selector": {
                        "buckets_path": {
                            "failed_count": "failed_docs._count"
                        },
                        "script": "params.failed_count > 0"
                    }
                }
            }
        }
    }
}`

const ViewRunActivityDrilldownQuery = `{
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
                        "component_id": "{{.componentIdForNestedDrillDown}}"
                    }
                },
                {
                    "term": {
                        "branch_id": "{{.branch}}"
                    }
                },
                {
                    "term": {
                        "automation_id": "{{.automationId}}"
                    }
                },
                {
                    "term": {
                        "test_suite_name": "{{.testSuiteName}}"
                    }
                },
                {
                    "term": {
                        "test_case_name": "{{.testCaseName}}"
                    }
                }
            ]
        }
    },
    "aggs": {
        "viewRunActivity": {
            "scripted_metric": {
                "init_script": "state.map = new HashMap();",
                "map_script": "def map = state.map; def uniqueTestCaseRun =doc.test_suite_name.value + '_' + doc.test_case_name.value + '_' + doc.run_id.value;def v = ['run_id':doc.run_id.value, 'test_suite_name':doc.test_suite_name.value, 'test_case_name':doc.test_case_name.value, 'test_case_status':doc.status.value, 'test_case_duration':doc.duration.value<10?10:doc.duration.value, 'job_id':doc.job_id.value, 'build_id':doc.run_number.value, 'start_time':doc['run_start_time'].value.getMillis(),'branch_id': doc.branch_id.value, 'workflow_name': doc.automation_name.value, 'source': doc.containsKey('source') ? doc.source.value : 'CloudBees'];map.put(uniqueTestCaseRun, v);",
                "combine_script": "return state.map;",
                "reduce_script": "def resultMap=new HashMap();def consolidatedMap=new HashMap();def resultList=new ArrayList();def outputMap=new HashMap();def resultMapKey='';for(state in states){if(state!=null&&state.size()>0){for(key in state.keySet()){def testCaseRecord=state.get(key);consolidatedMap.put(key,testCaseRecord)}}}if(consolidatedMap.size()>0){for(uniqueKey in consolidatedMap.keySet()){def item=consolidatedMap.get(uniqueKey);resultList.add(item);if(resultMap.containsKey(resultMapKey)){def map=resultMap.get(resultMapKey);def statusCount=map.get(item.test_case_status);map.put(item.test_case_status,++statusCount);def total=map.get('total');map.put('total',++total);def durationList=map.get('duration_array');durationList.add(item.test_case_duration);map.put('duration_array',durationList);if(total>=consolidatedMap.size()){double sum=0,size=durationList.size();for(int i=0;i<size;i++){sum+=durationList.get(i)}def avgDuration=(double)sum/size;map.put('avg_duration',avgDuration)}resultMap.put(resultMapKey,map)}else{def map=[:];def durationList=new ArrayList();map.put('SKIPPED',0);map.put('FAILED',0);map.put('PASSED',0);map.put('total',1);durationList.add(item.test_case_duration);map.put('duration_array',durationList);map.put('avg_duration',item.test_case_duration);map.put('workflow',item.workflow_name);map.put('source',item.source);map.put(item.test_case_status,1);resultMapKey=item.test_suite_name+'_'+item.test_case_name;resultMap.put(resultMapKey,map)}}}resultList.sort((a,b)->{def aValue=(Long)(a.get('start_time'));def bValue=(Long)(b.get('start_time'));return bValue.compareTo(aValue);});outputMap.put('headers',resultMap.get(resultMapKey));outputMap.put('section',resultList);return outputMap;"
            }
        }
    }
}`

const ViewRunActivityLogsDrilldownQuery = `{
    "size": 10,
    "_source": [
        "std_out",
        "std_err",
        "error_trace",
        "run_number",
        "automation_id",
        "branch_id",
        "component_id",
        "run_id"
    ],
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
                        "branch_id": "{{.branch}}"
                    }
                },
                {
                    "term": {
                        "test_suite_name": "{{.testSuiteName}}"
                    }
                },
                {
                    "term": {
                        "test_case_name": "{{.testCaseName}}"
                    }
                },
                {
                    "term": {
                        "run_id": "{{.runId}}"
                    }
                },
				{
                    "term": {
                        "run_number": "{{.runNumber}}"
                    }
                }
            ]
        }
    }
}`

const TestAutomationRunsQuery = `{
    "_source": false,
    "size": 0,
    "query":
    {
        "bool":
        {
            "filter":
            [
                {
                    "term":
                    {
                        "org_id": "{{.orgId}}"
                    }
                },
                {
                    "range":
                    {
                        "timestamp":
                        {
							"gte": "{{.startDate}}",
                            "lte": "{{.endDate}}",
                            "format": "yyyy-MM-dd HH:mm:ss",
                            "time_zone": "{{.timeZone}}"
                        }
                    }
                }
            ]
        }
    },
    "aggs":
    {
        "component_activity":
        {
            "scripted_metric":
            {
                "init_script": "state.data_map = [:];",
                "map_script": "def map=state.data_map;def key=doc.component_id.value+'_'+doc.automation_id.value+'_'+doc.run_id.value;def v=['org_id':doc.org_id.value,'component_id':doc.component_id.value,'automation_id':doc.automation_id.value,'run_id':doc.run_id.value];map.put(key,v);",
                "combine_script": "return state.data_map;",
                "reduce_script": "def resultMap=new HashMap();for(a in states){if(a!=null){for(i in a.keySet()){def record=a.get(i);def key=record.component_id+'_'+record.automation_id;if(resultMap.containsKey(key)){def lastRecord=resultMap.get(key);def runIdsSet=lastRecord.get('run_ids');runIdsSet.add(record.get('run_id'));resultMap.put(key,lastRecord);}else{def set=new HashSet();set.add(record.get('run_id'));record.put('run_ids',set);resultMap.put(key,record);}}}}return resultMap;"
            }
        }
    }
}`

const OpenVulnerabilitiesSubRowsQuery = `{
  "_source": false,
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
          "term": {
            "code": "{{.vulCode}}"
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
    "drilldowns": {
      "scripted_metric": {
        "params": {
          "timeZone": "{{.timeZone}}",
		  "timeFormat": "{{.timeFormat}}"
        },
        "init_script": "state.statusMap = [: ];",
        "map_script": "def map=state.statusMap;def key=doc.org_id.value+'_'+doc.component_id.value+'_'+doc.github_branch.value+'_'+doc.code.value+'_'+doc.scanner_name.value+'_'+doc['timestamp'].getValue().toEpochSecond()*1000;def v=['timestamp':doc['timestamp'].getValue().toEpochSecond()*1000,'bug_status':doc.bug_status.value,'code':doc.code.value,'branch':doc.github_branch.value,'severity':doc.severity.value,'recurrences':params['_source']['failure_files'].size(),'scanner_name':doc.scanner_name.value,'name':doc.name.value,'component_name':doc.component_name.value,'date_of_discovery':doc.date_of_discovery.value,'scan_time':doc.scan_time.value,'org_id':doc.org_id.value,'component_id':doc.component_id.value,'sla':'','run_id':doc.run_id.value];map.put(key,v);",
        "combine_script": "return state.statusMap;",
        "reduce_script": "Instant currentDate=Instant.ofEpochMilli(new Date().getTime());def statusMap=new HashMap();def slaNames=['Breached','At Risk','On Track'];def resultMap=new HashMap();def slaRules=new HashMap();slaRules.put('Breached',3);slaRules.put('AtRisk',2);slaRules.put('OnTrack',1);for(a in states){if(a!=null){for(i in a.keySet()){def record=a.get(i);def key=record.org_id+'_'+record.code;if(statusMap.containsKey(key)){def vulDetailsMap=statusMap.get(key);def vulKey=record.org_id+'_'+record.component_id+'_'+record.branch+'_'+record.code+'_'+record.scanner_name;if(vulDetailsMap.containsKey(vulKey)){def lastRecord=vulDetailsMap.get(vulKey);if(lastRecord.timestamp<record.timestamp){vulDetailsMap.put(vulKey,record)}}else{vulDetailsMap.put(vulKey,record)}}else{def vulDetailsMap=new HashMap();def vulKey=record.org_id+'_'+record.component_id+'_'+record.branch+'_'+record.code+'_'+record.scanner_name;vulDetailsMap.put(vulKey,record);statusMap.put(key,vulDetailsMap)}}}}def resultList=new ArrayList();DateTimeFormatter formatter;if(params.timeFormat=='12h'){formatter=DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone));}else{formatter=DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone));}if(statusMap.size()>0){for(uniqueKey in statusMap.keySet()){def vulMapBranchLevel=statusMap.get(uniqueKey);def vulDetailsListBranchLevel=new ArrayList();def uniqueComponents=new HashSet();def hasVulDetailsAtOrgLevel=false;def isVulOpen=false;for(vulKey in vulMapBranchLevel.keySet()){def vul=vulMapBranchLevel.get(vulKey);if(vul.bug_status=='Open'||vul.bug_status=='Reopened'){uniqueComponents.add(vul.component_name);isVulOpen=true;Instant startDate=Instant.ofEpochMilli(vul.date_of_discovery.getMillis());def curSeverity=vul.severity;def severityCode=0;if(curSeverity=='VERY_HIGH'){vul.severity='Very High';severityCode=4;}else if(curSeverity=='HIGH'){vul.severity='High';severityCode=3;}else if(curSeverity=='MEDIUM'){vul.severity='Medium';severityCode=2;}else if(curSeverity=='LOW'){vul.severity='Low';severityCode=1;}def diffAge=ChronoUnit.DAYS.between(startDate,currentDate);def SLAToolTip='';if(diffAge>=slaRules.Breached){vul.sla='Breached';Instant breachedOn=startDate.plus(Duration.ofDays(slaRules.get('Breached')));SLAToolTip='Breached on: '+formatter.format(breachedOn);}else if(diffAge>=slaRules.AtRisk){vul.sla='At Risk';Instant willBreachOn=startDate.plus(Duration.ofDays(slaRules.get('Breached')));SLAToolTip='Will breach on: '+formatter.format(willBreachOn);}else{vul.sla='On Track';}def map=[:];map.put('lastDiscovered',formatter.format(vul.scan_time));map.put('component',vul.component_name);map.put('branch',vul.branch);map.put('scannerName',vul.scanner_name);map.put('recurrences',vul.recurrences);map.put('componentId',vul.component_id);map.put('sla',vul.sla);if(SLAToolTip!=''){map.put('slaToolTipContent',SLAToolTip);}def reportInfoMap=new HashMap();reportInfoMap.put('code',vul.code);reportInfoMap.put('branch',vul.branch);reportInfoMap.put('scanner_name',vul.scanner_name);reportInfoMap.put('run_id',vul.run_id);reportInfoMap.put('run_number',vul.run_number);reportInfoMap.put('component_id',vul.component_id);def drillDownInfoMap=new HashMap();drillDownInfoMap.put('reportId','cwe-top25-vulnerabilities-view-location');drillDownInfoMap.put('reportTitle','Open Vulnerabilities');drillDownInfoMap.put('reportInfo',reportInfoMap);map.put('drillDown',drillDownInfoMap);if(vulDetailsListBranchLevel.size()<20){vulDetailsListBranchLevel.add(map);}}}if(isVulOpen){resultList.add(vulDetailsListBranchLevel)}}}if(!resultList.isEmpty()){return resultList[0];}else{return[];}"
      }
    }
  }
}`

const VulnerabilitiesOverviewSubRowsQuery = `{
  "_source": false,
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
          "term": {
            "code": "{{.vulCode}}"
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
    "drilldowns": {
      "scripted_metric": {
        "params": {
          "timeZone": "{{.timeZone}}",
		  "timeFormat": "{{.timeFormat}}"
        },
        "init_script": "state.statusMap = [: ];",
        "map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'recurrences':params['_source']['failure_files'].size(), 'scanner_name':doc.scanner_name.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value, 'scan_time':doc.scan_time.value, 'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':'', 'run_id':doc.run_id.value];map.put(key, v);",
        "combine_script": "return state.statusMap;",
        "reduce_script": "Instant currentDate = Instant.ofEpochMilli(new Date().getTime()); def statusMap = new HashMap(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); DateTimeFormatter formatter; if (params.timeFormat == '12h') { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); } else { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); } for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.code; if (statusMap.containsKey(key)) { def vulDetailsMap = statusMap.get(key); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record.code + '_' + record.scanner_name; if (vulDetailsMap.containsKey(vulKey)) { def lastRecord = vulDetailsMap.get(vulKey); if (lastRecord.timestamp < record.timestamp) { vulDetailsMap.put(vulKey, record); } } else { vulDetailsMap.put(vulKey, record); } } else { def vulDetailsMap = new HashMap(); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record.code + '_' + record.scanner_name; vulDetailsMap.put(vulKey, record); statusMap.put(key, vulDetailsMap); } } } } def resultList = new ArrayList(); if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vulMapBranchLevel = statusMap.get(uniqueKey); def vulDetailsListBranchLevel = new ArrayList(); for (vulKey in vulMapBranchLevel.keySet()) { def vul = vulMapBranchLevel.get(vulKey); Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); def diffAge = ChronoUnit.DAYS.between(startDate, currentDate); def SLAToolTip = ''; if (vul.bug_status == 'Resolved') { Instant resolutionDate = Instant.ofEpochMilli(vul.scan_time.getMillis()); diffAge = ChronoUnit.DAYS.between(startDate, resolutionDate); SLAToolTip = 'Date of resolution: ' + formatter.format(vul.scan_time); if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; } else { vul.sla = 'Within SLA'; } } else { if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else if (diffAge >= slaRules.AtRisk) { vul.sla = 'At risk'; } else { vul.sla = 'On track'; } } def map = [:]; map.put('lastDiscovered', formatter.format(vul.scan_time)); map.put('component', vul.component_name); map.put('branch', vul.branch); map.put('scannerName', vul.scanner_name); map.put('recurrences', vul.recurrences); map.put('componentId', vul.component_id); map.put('sla', vul.sla); map.put('status', vul.bug_status); if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } if (vul.bug_status != 'Resolved') { def reportInfoMap = new HashMap(); reportInfoMap.put('code', vul.code); reportInfoMap.put('branch', vul.branch); reportInfoMap.put('scanner_name', vul.scanner_name); reportInfoMap.put('run_id', vul.run_id); reportInfoMap.put('component_id', vul.component_id); def drillDownInfoMap = new HashMap(); drillDownInfoMap.put('reportId', 'cwe-top25-vulnerabilities-view-location'); drillDownInfoMap.put('reportTitle', 'Vulnerabilities overview'); drillDownInfoMap.put('reportInfo', reportInfoMap); map.put('drillDown', drillDownInfoMap); } vulDetailsListBranchLevel.add(map); } resultList.add(vulDetailsListBranchLevel); } } if (!resultList.isEmpty()) { return resultList[0]; } else { return []; }"
      }
    }
  }
}`

const Top25OpenVulnerabilitiesSubRowsQuery = `{
    "_source": false,
    "size": 0,
    "query": {
        "bool": {
            "filter": [{
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
                    "term": {
                        "code": "{{.vulCode}}"

                    }
                },

                {
                    "exists": {
                        "field": "date_of_discovery"
                    }
                },
                {
                    "terms": {
                        "severity": ["MEDIUM", "HIGH", "LOW", "VERY_HIGH"]
                    }
                },
                {
                    "terms": {
                        "code": [
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
                            "CWE-276"
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
                    "timeZone": "{{.timeZone}}",
					"timeFormat": "{{.timeFormat}}"
                },
                "init_script": "state.statusMap = [: ];",
                "map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'recurrences':params['_source']['failure_files'].size(), 'scanner_name':doc.scanner_name.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value, 'scan_time':doc.scan_time.value, 'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':'', 'run_id':doc.run_id.value];map.put(key, v);",
                "combine_script": "return state.statusMap;",
                "reduce_script": "Instant currentDate = Instant.ofEpochMilli(new Date().getTime()); def statusMap = new HashMap(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.code; if (statusMap.containsKey(key)) { def vulDetailsMap = statusMap.get(key); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; if (vulDetailsMap.containsKey(vulKey)) { def lastRecord = vulDetailsMap.get(vulKey); if (lastRecord.timestamp < record.timestamp) { vulDetailsMap.put(vulKey, record) } } else { vulDetailsMap.put(vulKey, record) } } else { def vulDetailsMap = new HashMap(); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; vulDetailsMap.put(vulKey, record); statusMap.put(key, vulDetailsMap) } } } } def resultList = new ArrayList(); DateTimeFormatter formatter; if (params.timeFormat == '12h') { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); } else { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); } if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vulMapBranchLevel = statusMap.get(uniqueKey); def vulDetailsListBranchLevel = new ArrayList(); def uniqueComponents = new HashSet(); def hasVulDetailsAtOrgLevel = false; def isVulOpen = false; for (vulKey in vulMapBranchLevel.keySet()) { def vul = vulMapBranchLevel.get(vulKey); if (vul.bug_status == 'Open' || vul.bug_status == 'Reopened') { uniqueComponents.add(vul.component_name); isVulOpen = true; Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); def severityCode = 0; def curSeverity = vul.severity; if (curSeverity == 'VERY_HIGH') { vul.severity = 'Very high'; severityCode = 4; } else if (curSeverity == 'HIGH') { vul.severity = 'High'; severityCode = 3; } else if (curSeverity == 'MEDIUM') { vul.severity = 'Medium'; severityCode = 2; } else if (curSeverity == 'LOW') { vul.severity = 'Low'; severityCode = 1; } def diffAge = ChronoUnit.DAYS.between(startDate, currentDate); def SLAToolTip = ''; if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else if (diffAge >= slaRules.AtRisk) { vul.sla = 'At risk'; Instant willBreachOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Will breach on: ' + formatter.format(willBreachOn); } else { vul.sla = 'On track'; } def map = [: ]; if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } map.put('lastDiscovered', formatter.format(vul.scan_time)); map.put('component', vul.component_name); map.put('branch', vul.branch); map.put('scannerName', vul.scanner_name); map.put('recurrences', vul.recurrences); map.put('componentId', vul.component_id); map.put('sla', vul.sla); def reportInfoMap = new HashMap(); reportInfoMap.put('code', vul.code); reportInfoMap.put('branch', vul.branch); reportInfoMap.put('scanner_name', vul.scanner_name); reportInfoMap.put('run_id', vul.run_id); reportInfoMap.put('component_id', vul.component_id); def drillDownInfoMap = new HashMap(); drillDownInfoMap.put('reportId', 'cwe-top25-vulnerabilities-view-location'); drillDownInfoMap.put('reportTitle', 'CWETM top 25 vulnerabilities'); drillDownInfoMap.put('reportInfo', reportInfoMap); map.put('drillDown', drillDownInfoMap); if (vulDetailsListBranchLevel.size() < 20) { vulDetailsListBranchLevel.add(map); } } } if (isVulOpen) { resultList.add(vulDetailsListBranchLevel) } } } if (!resultList.isEmpty()) { return resultList[0]; } else { return []; }"
            }
        }
    }
}`

const RunDetailsTestResultsTestSuitesViewQuery = `{
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
                        "component_id": "{{.componentIdForNestedDrillDown}}"
                    }
                },
                {
                    "term": {
                        "run_id": "{{.runId}}"
                    }
                },
				{
                    "term": {
                        "run_number": "{{.runNumber}}"
                    }
                }
            ]
        }
    },
    "aggs": {
        "test_suite_buckets": {
            "terms": {
                "field": "test_suite_name",
                "size": 10000
            },
            "aggs": {
                "test_suite_doc": {
                    "top_hits": {
                        "size": 1,
                        "_source": [
                            "total",
                            "passed",
                            "failed",
                            "skipped",
                            "duration",
                            "run_id",
							"run_number",
                            "component_id",
                            "test_suite_name"
                        ]
                    }
                }
            }
        }
    }
}`

const RunDetailsTestResultsTestCasesViewQuery = `{
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
                        "component_id": "{{.componentIdForNestedDrillDown}}"
                    }
                },
                {
                    "term": {
                        "run_id": "{{.runId}}"
                    }
                },
				{
                    "term": {
                        "run_number": "{{.runNumber}}"
                    }
                }
            ]
        }
    },
    "aggs": {
        "test_case_buckets": {
            "terms": {
                "script": "doc['test_suite_name'].value + '_' + doc['test_case_name'].value",
                "size": 10000
            },
            "aggs": {
                "test_case_doc": {
                    "top_hits": {
                        "size": 1,
                        "_source": [
                            "status",
                            "test_case_name",
                            "test_suite_name",
                            "duration",
                            "run_id",
							"run_number",
                            "component_id",
                            "std_out",
                            "std_err",
                            "error_trace"
                        ]
                    }
                }
            }
        }
    }
}`

const RunDetailsTotalTestCasesDrillDownQuery = `{
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
                        "component_id": "{{.componentIdForNestedDrillDown}}"
                    }
                },
                {
                    "term": {
                        "run_id": "{{.runId}}"
                    }
                },
				{
                    "term": {
                        "run_number": "{{.runNumber}}"
                    }
                },
                {
                    "term": {
                        "test_suite_name": "{{.testSuiteName}}"
                    }
                }
            ]
        }
    },
    "aggs": {
        "test_case_buckets": {
            "terms": {
                "field": "test_case_name",
                "size": 10000
            },
            "aggs": {
                "test_case_doc": {
                    "top_hits": {
                        "size": 1,
                        "_source": [
                            "status",
                            "test_case_name",
                            "test_suite_name",
                            "duration",
                            "run_id",
							"run_number",
                            "component_id",
                            "std_out",
                            "std_err",
                            "error_trace"
                        ]
                    }
                }
            }
        }
    }
}`

const RunDetailsTestCaseLogDrillDownQuery = `{
    "size": 1,
    "_source": [
        "std_out",
        "std_err",
        "error_trace"
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
                        "component_id": "{{.componentIdForNestedDrillDown}}"
                    }
                },
                {
                    "term": {
                        "run_id": "{{.runId}}"
                    }
                },
				{
                    "term": {
                        "run_number": "{{.runNumber}}"
                    }
                },
                {
                    "term": {
                        "test_suite_name": "{{.testSuiteName}}"
                    }
                },
                {
                    "term": {
                        "test_case_name": "{{.testCaseName}}"
                    }
                }
            ]
        }
    }
}`

const RunDetailsTestResultsIndicatorsQuery = `{
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
                        "component_id": "{{.componentIdForNestedDrillDown}}"
                    }
                },
                {
                    "term": {
                        "run_id": "{{.runId}}"
                    }
                },
				{
                    "term": {
                        "run_number": "{{.runNumber}}"
                    }
                }
            ]
        }
    },
    "aggs": {
        "test_suites": {
            "terms": {
                "field": "test_suite_name",
                "size": 65000
            },
            "aggs": {
                "statuses": {
                    "top_hits": {
                        "size": 1,
                        "_source": [
                            "passed",
                            "skipped",
                            "failed"
                        ]
                    }
                }
            }
        }
    }
}`

const VulnerabiltyByScannerTypeDrillDownSubRowsQuery = `{
    "_source": false,
    "size": 0,
    "query": {
        "bool": {
            "filter": [
                {
                    "range": {
                        "timestamp": {
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
                        "code": "{{.vulCode}}"

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
                    "timeZone": "{{.timeZone}}",
					"timeFormat": "{{.timeFormat}}"
                },
                "init_script": "state.statusMap = [: ];",
                "map_script": "def map = state.statusMap;def key = doc.org_id.value + '_' + doc.component_id.value + '_' + doc.github_branch.value + '_' + doc.code.value + '_' + doc.scanner_name.value + '_' + doc['timestamp'].getValue().toEpochSecond() * 1000;def v = ['timestamp':doc['timestamp'].getValue().toEpochSecond() * 1000, 'bug_status':doc.bug_status.value,'code':doc.code.value, 'branch':doc.github_branch.value,'severity':doc.severity.value, 'recurrences':params['_source']['failure_files'].size(), 'scanner_name':doc.scanner_name.value,'scanner_type':doc.scanner_type.value, 'name':doc.name.value, 'component_name':doc.component_name.value, 'date_of_discovery':doc.date_of_discovery.value, 'scan_time':doc.scan_time.value, 'org_id':doc.org_id.value, 'component_id':doc.component_id.value, 'sla':'', 'run_id':doc.run_id.value];map.put(key, v);",
                "combine_script": "return state.statusMap;",
                "reduce_script": "Instant currentDate = Instant.ofEpochMilli(new Date().getTime()); def statusMap = new HashMap(); def slaNames = ['Breached', 'At risk', 'On track']; def resultMap = new HashMap(); def slaRules = new HashMap(); slaRules.put('Breached', 3); slaRules.put('AtRisk', 2); slaRules.put('OnTrack', 1); for (a in states) { if (a != null) { for (i in a.keySet()) { def record = a.get(i); def key = record.org_id + '_' + record.code; if (statusMap.containsKey(key)) { def vulDetailsMap = statusMap.get(key); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; if (vulDetailsMap.containsKey(vulKey)) { def lastRecord = vulDetailsMap.get(vulKey); if (lastRecord.timestamp < record.timestamp) { vulDetailsMap.put(vulKey, record) } } else { vulDetailsMap.put(vulKey, record) } } else { def vulDetailsMap = new HashMap(); def vulKey = record.org_id + '_' + record.component_id + '_' + record.branch + '_' + record .code + '_' + record.scanner_name; vulDetailsMap.put(vulKey, record); statusMap.put(key, vulDetailsMap) } } } } def resultList = new ArrayList(); DateTimeFormatter formatter; if (params.timeFormat == '12h') { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd h:mm:ss a').withZone(ZoneId.of(params.timeZone)); } else { formatter = DateTimeFormatter.ofPattern('yyyy/MM/dd HH:mm:ss').withZone(ZoneId.of(params.timeZone)); } if (statusMap.size() > 0) { for (uniqueKey in statusMap.keySet()) { def vulMapBranchLevel = statusMap.get(uniqueKey); def vulDetailsListBranchLevel = new ArrayList(); def vulMapOrgLevel = new HashMap(); def hasVulDetailsAtOrgLevel = false; for (vulKey in vulMapBranchLevel.keySet()) { def vul = vulMapBranchLevel.get(vulKey); Instant startDate = Instant.ofEpochMilli(vul.date_of_discovery.getMillis()); def severityCode = 0; def curSeverity = vul.severity; if (curSeverity == 'VERY_HIGH') { vul.severity = 'Very high'; severityCode = 4; } else if (curSeverity == 'HIGH') { vul.severity = 'High'; severityCode = 3; } else if (curSeverity == 'MEDIUM') { vul.severity = 'Medium'; severityCode = 2; } else if (curSeverity == 'LOW') { vul.severity = 'Low'; severityCode = 1; } def diffAge = ChronoUnit.DAYS.between(startDate, currentDate); def SLAToolTip = ''; def statusToolTip = ''; if (vul.bug_status == 'Resolved') { statusToolTip = 'Date of resolution: ' + formatter.format(vul.scan_time); if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else { vul.sla = 'Within SLA'; } } else { if (diffAge >= slaRules.Breached) { vul.sla = 'Breached'; Instant breachedOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Breached on: ' + formatter.format(breachedOn); } else if (diffAge >= slaRules.AtRisk) { vul.sla = 'At risk'; Instant willBreachOn = startDate.plus(Duration.ofDays(slaRules.get('Breached'))); SLAToolTip = 'Will breach on: ' + formatter.format(willBreachOn); } else { vul.sla = 'On track'; } } def map = [: ]; if (SLAToolTip != '') { map.put('slaToolTipContent', SLAToolTip); } if (statusToolTip != '') { map.put('statusToolTipContent', statusToolTip); } map.put('lastDiscovered', formatter.format(vul.scan_time)); map.put('component', vul.component_name); map.put('componentId', vul.component_id); map.put('branch', vul.branch); map.put('scannerName', vul.scanner_name); map.put('recurrences', vul.recurrences); map.put('sla', vul.sla); map.put('status', vul.bug_status); if (vul.bug_status != 'Resolved') { def reportInfoMap = new HashMap(); reportInfoMap.put('code', vul.code); reportInfoMap.put('branch', vul.branch); reportInfoMap.put('scanner_name', vul.scanner_name); reportInfoMap.put('run_id', vul.run_id); reportInfoMap.put('component_id', vul.component_id); def drillDownInfoMap = new HashMap(); drillDownInfoMap.put('reportId', 'cwe-top25-vulnerabilities-view-location'); drillDownInfoMap.put('reportTitle', 'Vulnerabilities by security scan type'); drillDownInfoMap.put('reportInfo', reportInfoMap); map.put('drillDown', drillDownInfoMap); } vulDetailsListBranchLevel.add(map); if (!hasVulDetailsAtOrgLevel) { vulMapOrgLevel.put('vulnerabilityId', vul.code); vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery); vulMapOrgLevel.put('vulnerabilityName', vul.name); vulMapOrgLevel.put('severity', vul.severity); vulMapOrgLevel.put('severityCode', severityCode); hasVulDetailsAtOrgLevel = true } if (vul.date_of_discovery.getMillis() < vulMapOrgLevel.get('firstDiscovered').getMillis()) { vulMapOrgLevel.put('firstDiscovered', vul.date_of_discovery) } } resultList.add(vulDetailsListBranchLevel); } } if (!resultList.isEmpty()) { return resultList[0]; } else { return []; }"
            }
        }
    }
}`
