package constants

const (
	PLUGINS                  = "plugins"
	VERSION                  = "version"
	LATEST_VERSION           = "latestVersion"
	VERSION_UPDATE_AVAILABLE = "versionUpdateAvailable"
	VERSION_UPDATE_MESSAGE   = "versionUpdateMessage"
	VERSION_UPDATE_HINT      = "versionUpdateHint"

	METRICS                              = "metrics"
	METRICS_DATA                         = "metricsData"
	TOTAL_EXECUTORS                      = "totalExecutors"
	TOTAL_EXECUTOR_KEY                   = "jenkins.executor.count.value"
	FREE_EXECUTORS                       = "freeExecutors"
	FREE_EXECUTOR_KEY                    = "jenkins.executor.free.value"
	TOTAL_NODES                          = "totalNodes"
	TOTAL_NODES_KEY                      = "jenkins.node.count.value"
	PLUGIN_INFO                          = "pluginsInfo"
	PLUGIN_INFO_TITLE                    = "Installed plugin information"
	COUNT                                = "count"
	EMAIL                                = "email"
	NAME                                 = "name"
	LONG_NAME                            = "longName"
	SHORT_NAME                           = "shortName"
	ENABLED                              = "enabled"
	HAS_UPDATE                           = "hasUpdate"
	UPDATES                              = "updates"
	REQUIRED_CORE_VERSION                = "requiredCoreVersion"
	MINIMUM_JAVA_VERSION                 = "minimumJavaVersion"
	DEPENDENCIES                         = "dependencies"
	COMPLETED_RUNS                       = "completedRuns"
	PROJECT_TYPES                        = "projectTypes"
	TYPE                                 = "type"
	SCATTER_TYPE                         = "SCATTER"
	COLOR_SCHEME                         = "colorScheme"
	LIGHT_COLOR_SCHEME                   = "lightColorScheme"
	RUN_INFORMATION                      = "Run information"
	JOBS                                 = "jobs"
	HEADER_DETAILS                       = "headerDetails"
	COMPLETED_RUNS_DATA                  = "completedRunsData"
	CONNECTED_CONTROLLERS                = "connectedControllers"
	TOTAL_CONTROLLERS                    = "totalControllers"
	ENDPOINT_JOBS                        = "endpointJobs"
	JENKINS_STABLE_VERSION_URL           = "https://updates.jenkins.io/stable/latestCore.txt"
	JENKINS_LATEST_VERSION_URL           = "https://updates.jenkins.io/current/latestCore.txt"
	CBCI_CJOC_LATEST_VERSION_URL         = "https://downloads.cloudbees.com/cloudbees-core/traditional/operations-center/rolling/war/"
	FREESTYLE_JOB                        = "hudson.model.FreeStyleProject"
	WORKFLOW_JOB                         = "org.jenkinsci.plugins.workflow.job.WorkflowJob"
	BACKUP_PROJECT                       = "com.infradna.hudson.plugins.backup.BackupProject"
	JIRA_ENDPOINT_CONTRIBUTION_ID        = "cb.jira.jira-token-endpoint-type"
	ENVIRONMENT_ENDPOINT_CONTRIBUTION_ID = "cb.configuration.basic-environment"
	MATRIX_JOB                           = "hudson.matrix.MatrixProject"
	FREESTYLE_TYPE                       = "Freestyle"
	PIPELINE_TYPE                        = "Pipeline"
	BACKUP_PROJECT_TYPE                  = "Backup project"
	MULTI_CONFIG_TYPE                    = "Multi-config"
	JENKINS_FOLDER                       = "com.cloudbees.hudson.plugins.folder.Folder"
	MULTI_FOLDER                         = "Multi-folder"
	JENKINS_BRANCH                       = "org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject"
	MULTI_BRANCH                         = "Multi-branch"
	JENKINS_MULTI_JOB                    = "com.tikal.jenkins.plugins.multijob.MultiJobProject"
	JENKINS_BLUE_STEEL_FOLDER            = "com.cloudbees.opscenter.bluesteel.folder.BlueSteelTeamFolder"
	MULTI_JOB                            = "Multi-job"
	PIPELINE_JOB_TEMPLATE                = "com.cloudbees.pipeline.governance.templates.classic.standalone.GovernancePipelineJobTemplate"
	PIPELINE_TEMPLATE                    = "Pipeline template"
	CJOC                                 = "CJOC"
	CBCI                                 = "CBCI"
	JAAS                                 = "JAAS"
	JENKINS                              = "JENKINS"
	JENKINS_ENDPOINT                     = "cb.jenkins.jenkins-private-endpoint-type"
	CJOC_ENDPOINT                        = "cb.cjoc.cjoc-private-endpoint-type"
	CBCI_ENDPOINT                        = "cb.cbci.cbci-private-endpoint-type"
	JAAS_ENDPOINT                        = "cb.cbci.cbci-jaas-endpoint-type"
	CI_TOOL_ID                           = "ciToolId"
	CI_TOOL_TYPE                         = "ciToolType"
	ENDPOINT_IDS                         = "endpointIds"
	HITS                                 = "hits"
	SOURCE                               = "_source"
	SOURCE_PROVIDER                      = "source"
	SYSTEM_HEALTH                        = "system_health"
	DESCRIPTION                          = "description"
	DISK_SPACE                           = "disk-space"
	TEMPORARY_SPACE                      = "temporary-space"
	THREAD_DEADLOCK                      = "thread-deadlock"
	DISK_SPACE_DESCRIPTION               = "Monitoring disk space health by configured threshold"
	PLUGIN_DESCRIPTION                   = "Monitoring plugin health"
	TEMPORARY_SPACE_DESCRIPTION          = "Monitoring temporary space health by configured threshold"
	THREAD_DEADLOCK_DESCRIPTION          = "Monitoring thread deadlock in JVM"
	HEALTHY                              = "healthy"
	HEALTH_LIST                          = "healthList"
	HEALTH_SCORE                         = "healthScore"
	HEALTH_STATUS                        = "healthStatus"
	WARNING                              = "warning"
	FAILED                               = "failed"
	JOBID                                = "job_id"
	RUNID                                = "run_id"
	ENDPOINT_ID                          = "endpoint_id"
	START_TIME_MILLIS                    = "start_time_in_millis"
	RESULT                               = "result"
	COLOR_0                              = "color0"
	COLOR_SCHEME_0                       = "#8BC34A"
	COLOR_SCHEME_1                       = "#EF5350"
	COLOR_SCHEME_2                       = "#FFA726"
	COLOR_SCHEME_3                       = "#BDBDBD"
	LIGHT_COLOR_0                        = "#8BC34A"
	LIGHT_COLOR_1                        = "#EF5350"
	LIGHT_COLOR_2                        = "#FF970A"
	LIGHT_COLOR_3                        = "#424242"
	JOB_NAME                             = "job_name"
	DISPLAY_NAME                         = "display_name"
	FILTER_TYPE                          = "filterType"
	JOB_IDS                              = "jobIds"
	ENDPOINT_ID_KEY                      = "endpointId"
	EXECUTED                             = "executed"
	ABORTED                              = "aborted"
	UNSTABLE                             = "unstable"
	NOT_BUILT                            = "notBuilt"
	TOTAL_DURATION                       = "totalDuration"
	HYPHEN                               = "-"
	MONTH_STRING                         = "mo "
	DAY_STRING                           = "d "
	HOUR_STRING                          = "h "
	MINUTE_STRING                        = "m "
	SECOND_STRING                        = "s "
	TOTAL                                = "total"
	NOT_INSTALLED                        = "NOT_INSTALLED"
	START_DATE                           = "startDate"
	END_DATE                             = "endDate"
	EMPTY_SECOND                         = "0s"
	DAY_FORMAT                           = "Mon"
	ZERO                                 = "0"
	AM_12                                = "12am"
	PM_12                                = "12pm"
	SUNDAY                               = "Sun"
	MONDAY                               = "Mon"
	TUESDAY                              = "Tue"
	WEDNESDAY                            = "Wed"
	THURSDAY                             = "Thu"
	FRIDAY                               = "Fri"
	SATURDAY                             = "Sat"
	ACTUAL                               = "actual"
	EXPECTED                             = "expected"
	ACTIVE_RUNS                          = "Active runs"
	ACTIVE_RUNS_DESCRIPTION              = "Active runs shows how many active job/pipeline runs are executing."
	IDLE_EXECUTORS                       = "Idle executors"
	IDLE_EXECUTORS_DESCRIPTION           = "Idle executors shows how many executors are online and idle."
	WAITING_RUNS                         = "Runs waiting to start"
	WAITING_RUNS_DESCRIPTION             = "Runs waiting to start shows how many job/pipeline runs are in the queue waiting to start."
	WAITING_TIME                         = "Average time waiting to start"
	WAITING_TIME_DESCRIPTION             = "Average time waiting to start shows the average time in the queue for all queued build requests."
	IDLE_TIME                            = "Current time to idle"
	IDLE_TIME_DESCRIPTOR                 = "Current time to idle shows an estimation of how long before new builds can start.."
	HOUR                                 = "hour"
	DAY                                  = "day"
	LATEST_ORG                           = "latest_per_org"
	BUCKETS                              = "buckets"
	LATEST_DOC                           = "latest_doc"
	ACTIVE_RUNS_KEY                      = "active_runs"
	IDLE_EXECUTORS_KEY                   = "idle_executor"
	WAITING_RUNS_KEY                     = "runs_waiting_to_start"
	WAITING_TIME_KEY                     = "avg_time_waiting_to_start"
	IDLE_TIME_KEY                        = "current_time_to_idle"
	VIEW_OPTION                          = "viewOption"
	IDLE_TIME_VIEW                       = "CurrentTimeToIdle"
	WAITING_TIME_VIEW                    = "AvgTimeWaitingToStart"
	VALUES                               = "values"
	UP_COLOR_0                           = "#4FC3F7"
	UP_COLOR_1                           = "#2196f3"
	SCATTER                              = "SCATTER_WITH_SINGLE_COLOR"
	ACTIVITY_DAY_KEY                     = "activity_day"
	ACTIVITY_TIME_KEY                    = "activity_time"
	ZERO_PERCENT                         = "0%"
	DATE_PARSE                           = "2006-01-02 15:04:05"
	CJOC_ENDPOINT_ID                     = "cjocEndpointId"
	RUN_INFORMATION_KEY                  = "runInformation"
	EMPTY_STRING                         = ""
	HREF_STRING                          = "a href=\""
	SLASH                                = "/"
	NINE                                 = "9"
	DOT                                  = "."
	TOOL_URL                             = "toolUrl"
	CJOC_CONTROLLER_INFO                 = "cjocControllerInfo"
	ACTIVITIES                           = "activities"
	RUN_TIME                             = "runTime"
	SUCCESS_KEY                          = "SUCCESS"
	FAILED_KEY                           = "FAILED"
	FAILURE_KEY                          = "FAILURE"
	ABORTED_KEY                          = "ABORTED"
	UNSTABLE_KEY                         = "UNSTABLE"
	NOT_BUILT_KEY                        = "NOT_BUILT"
	AVERAGE_RUN_TIME                     = "avgRunTime"
	BAR_WITH_BOTH_AXIS                   = "BAR_WITH_BOTH_AXIS"
	COLOR_1                              = "color1"
	COLOR_SCHEME_0_0                     = "#00BFA8"
	COLOR_SCHEME_0_1                     = "#056459"
	COLOR_SCHEME_1_0                     = "#ED5252"
	COLOR_SCHEME_1_1                     = "#640505"
	COLOR_SCHEME_2_0                     = "#FFA726"
	COLOR_SCHEME_2_1                     = "#B96E00"
	COLOR_SCHEME_3_0                     = "#969696"
	COLOR_SCHEME_3_1                     = "#606060"
	SUCCESS_RUNS                         = "successRuns"
	FAILED_RUNS                          = "failedRuns"
	ABORTED_RUNS                         = "abortedRuns"
	UNSTABLE_RUNS                        = "unstableRuns"
	NOT_BUILT_RUNS                       = "notBuiltRuns"
	TOTAL_RUN_TIME                       = "totalRunTime"
	TOTAL_EXECUTED                       = "totalExecuted"
	CREATED_AT                           = "created_at"
	PLUGIN_KEY                           = "Plugins"
	DISK_SPACE_KEY                       = "Disk space"
	TEMPORARY_SPACE_KEY                  = "Temporary space"
	THREAD_DEADLOCK_KEY                  = "Thread deadlock"
	NUM_WORKERS                          = 20
	BATCH_SIZE                           = 20000
	MAX_SIZE_REACHED                     = "maxSizeReached"
	PLUGIN_COUNT                         = "count"
	TEST_SUITE_VIEW                      = "testSuite"
	TEST_CASE_VIEW                       = "testCase"
	COMPONENTS_VIEW                      = "componentsView"
	RESULT_COUNTS                        = "result_counts"
	SUCCESSFUL                           = "SUCCESSFUL"
	CANCELED                             = "CANCELED"
	COUNT_INFO                           = "countInfo"
	DOC_COUNT                            = "doc_count"
	TOTAL_RUNS                           = "Total runs"
	SUCCESSFUL_HEADER                    = "Successful"
	CANCELED_HEADER                      = "Canceled"
	FAILED_HEADER                        = "Failed"
	UNSTABLE_HEADER                      = "Unstable"
)

var INSIGHTS_SUPPORTED_TYPE = []string{JENKINS_ENDPOINT, CBCI_ENDPOINT, CJOC_ENDPOINT, JAAS_ENDPOINT}

const CiProjectTypesQuery = `{
	"size": 0,
	"query": {
		"bool": {
		}
	},
	"aggs": {
		"projectTypes": {
			"scripted_metric": {
				"init_script": "state.dataMap = [:];",
				"map_script": "def map = state.dataMap; def key = doc.org_id.value + '_' + doc.endpoint_id.value + '_' + doc.job_id.value; def v = ['org_id': doc.org_id.value, 'endpoint_id': doc.endpoint_id.value, 'job_id': doc.job_id.value, 'type': doc.type.value]; map.put(key, v);",
				"combine_script": "return state.dataMap;",
				"reduce_script": "def tmpMap = [: ], resultList = []; for (response in states) { if (response != null) { for (key in response.keySet()) { tmpMap.put(key, response.get(key)); } } } def jobTypeMap = new HashMap(); for (key in tmpMap.keySet()) { def record = tmpMap.get(key); if (jobTypeMap.containsKey(record.type)) { jobTypeMap.put(record.type, jobTypeMap.get(record.type) + 1); } else { jobTypeMap.put(record.type, 1); } } def sortedEntries = jobTypeMap.entrySet().stream().sorted((a, b) -> a.getValue().compareTo(b.getValue())).collect(Collectors.toList()); for (entry in sortedEntries) { def resultMap = [: ]; resultMap.put('name', entry.getKey()); resultMap.put('value', entry.getValue()); resultList.add(resultMap); } if (resultList.size() == 0) { def resultMap = [: ], resultMap1 = [: ]; resultMap.put('name', 'Freestyle'); resultMap.put('value', 0); resultList.add(resultMap); resultMap1.put('name', 'Pipeline'); resultMap1.put('value', 0); resultList.add(resultMap1); } return resultList;"
			}
		}
	}
}`

const CiToolInsightFetchQuery = `{
	"query": {
	  "bool": {
	  }
	}
  }`

const CiToolVersionAndPluginCountQuery = `{
  "size": 0,
  "query": {
    "bool": {
      
    }
  },
  "aggs": {
    "count": {
      "scripted_metric": {
        "init_script": "state.dataMap = [:];",
        "map_script": "def map = state.dataMap; def key = doc.endpoint_id.value ; def v = []; if (params._source.containsKey('plugins')) { def tags = params._source['plugins'];  v = ['version': doc.version.value, 'count': tags.size()];  } else {   v = ['version': doc.version.value, 'count': 0]} map.put(key, v);",
        "combine_script": "return state.dataMap;",
        "reduce_script": "def tmpMap = [: ]; for (response in states) { if (response != null) { for (key in response.keySet()) { tmpMap.put(key, response.get(key)); } } } return tmpMap;"
      }
    }
  }
}`

const CiActivityOverviewActualFetchQuery = `
  {
	"size":0,
	"query": {
	  "bool": {
		"filter": [
		]
	  }
	},
	"aggs": {
		"latest_per_org": {
		  "terms": {
			"field": "endpoint_id",
			"order": {
			  "latest_timestamp": "desc"
			}
		  },
		  "aggs": {
			"latest_timestamp": {
			  "max": {
				"field": "timestamp"
			  }
			},
			"latest_doc": {
			  "top_hits": {
				"size": 1,
				"sort": [
				  {
					"timestamp": {
					  "order": "desc"
					}
				  }
				]
			  }
			}
		  }
		}
	  }
  }`

const CiActivityOverviewFetchQuery = `{
	"query": {
	  "bool": {
		"filter": [
			{
				"range": {
				  "activity_date": {
					"gte": "{{.startDate}}",
					"lte": "{{.endDate}}",
					"format": "yyyy-MM-dd HH:mm:ss",
					"time_zone": "{{.timeZone}}"
				  }
				}
			  },
		  {
			"term": {
			  "activity_time": {{.hour}}
			}
		  },
		  {
			"term": {
			  "activity_day": "{{.day}}"
			}
		  }
		]
	  }
	}
  }`

const CiUsagePatternsFetchQuery = `{
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
		  }
		]
	  }
	},
	"aggs": {
	  "activities": {
		"scripted_metric": {
		  "combine_script": "return state.dataMap;",
		  "init_script": "state.dataMap = [:];",
		  "map_script": "def map = state.dataMap; def key = doc.org_id.value + '_' + doc.endpoint_id.value + '_' + doc.timestamp.value; def v = ['org_id': doc.org_id.value, 'endpoint_id': doc.endpoint_id.value, 'activity_day': doc.activity_day.value, 'activity_time': doc.activity_time.value,'active_runs': doc.active_runs.value,'idle_executor': doc.idle_executor.value,'runs_waiting_to_start': doc.runs_waiting_to_start.value,'current_time_to_idle': doc.current_time_to_idle.value,'avg_time_waiting_to_start': doc.avg_time_waiting_to_start.value]; map.put(key, v);",
		  "reduce_script": "def tmpMap = [: ], resultList = new ArrayList();for (response in states) {if (response != null) {for (key in response.keySet()) {tmpMap.put(key, response.get(key));}}}def resultMap = new HashMap();for (key in tmpMap.keySet()) {def record = tmpMap.get(key);if (resultMap.containsKey(record.activity_day)) {resultMap.get(record.activity_day).add(record);} else {def activities = new ArrayList();activities.add(record);resultMap.put(record.activity_day, activities);}}return resultMap;"
		}
	  }
	}
  }`

const CiUsagePatterns = `{
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
		  }
		]
	  }
	},
	"aggs": {
	  "activities": {
		"date_histogram": {
		  "field": "timestamp",
		  "calendar_interval": "hour",
		  "format": "yyyy-MM-dd HH:mm:ss",
		  "time_zone": "{{.timeZone}}"
		},
		"aggs": {
		  "endpoint_ids": {
			"terms": {
			  "field": "endpoint_id",
			  "size": 10
			},
			"aggs": {
			  "avg_idle_executor": {
				"avg": {
				  "field": "idle_executor"
				}
			  },
			  "avg_active_runs": {
				"avg": {
				  "field": "active_runs"
				}
			  },
			  "avg_runs_waiting_to_start": {
				"avg": {
				  "field": "runs_waiting_to_start"
				}
			  },
			  "avg_time_waiting_to_start": {
				"avg": {
				  "field": "avg_time_waiting_to_start"
				}
			  },
			  "avg_current_time_to_idle": {
				"avg": {
				  "field": "current_time_to_idle"
				}
			  }
			}
		  }
		}
	  }
	}
  }`
const CiToolInsightPluginsFetchQuery = `{
	"_source": "plugins", 
	"query": {
	  "bool": {
		"filter": [
		]
	  }
	}
  }`

const CiCompletedRunsFetchQuery = `{
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
		  }
		],
		"must_not": {
		  "term": {
			"result": "IN_PROGRESS"
		  }
		}
	  }
	}, "_source": [ "job_id","endpoint_id", "run_id", "result", "duration", "start_time_in_millis"],
	"aggs": {
		"completedRuns": {
		  "scripted_metric": {
			"init_script": "state.dataMap = [:];",
			"map_script": "def map = state.dataMap; def key = doc.endpoint_id.value + '_' + doc.job_id.value + '_'+doc.run_id; def v = ['endpoint_id': doc.endpoint_id.value, 'job_id': doc.job_id.value,'run_id': doc.run_id.value, 'result': doc.result.value,'duration': doc.duration.value,  'start_time_in_millis': doc.start_time_in_millis.value]; map.put(key, v);",
		  	"combine_script": "return state.dataMap;",
		  	"reduce_script": "def tmpMap = [:]; for (response in states) { if (response != null) { tmpMap.putAll(response); } } return new ArrayList(tmpMap.values());"
		  }
		}
	  }
  }`

const CiJobInfoFetchQuery = `{
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
		]
	  }
	},
	"aggs": {
	  "jobs": {
		"scripted_metric": {
		  "init_script": "state.dataMap = [:];",
		  "map_script": "def map = state.dataMap; def key = doc.org_id.value + '_' + doc.endpoint_id.value + '_' + doc.job_id.value; def v = ['org_id': doc.org_id.value, 'endpoint_id': doc.endpoint_id.value, 'job_id': doc.job_id.value, 'type': doc.type.value,'job_name': doc.job_name.value,'display_name': doc.display_name.value,'last_completed_run_id': doc.last_completed_run_id.value]; map.put(key, v);",
		  "combine_script": "return state.dataMap;",
		  "reduce_script": "def tmpMap = [: ], resultList = new ArrayList(); for (response in states) { if (response != null) { for (key in response.keySet()) { tmpMap.put(key, response.get(key)); } } } for (key in tmpMap.keySet()) { resultList.add(tmpMap.get(key)); } return resultList;"
		}
	  }
	}
  }`

const CiRunsExecutionInfoQuery = `{
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
		  }
		],
		"must_not": {
		  "term": {
			"result": "IN_PROGRESS"
		  }
		}
	  }
	},
	"aggs": {
		"completedRuns": {
		  "terms": {
			"field": "job_id",
			"size": 20000000
		  },
		  "aggs": {
			"endpoint_id": {
			  "terms": {
				"field": "endpoint_id",
				"size": 1
			  }
			},
			"result_buckets": {
			  "filters": {
				"filters": {
				  "SUCCESS": { "term": { "result": "SUCCESS" } },
				  "FAILED": { "term": { "result": "FAILURE" } },
				  "ABORTED": { "term": { "result": "ABORTED" } },
				  "UNSTABLE": { "term": { "result": "UNSTABLE" } },
				  "NOT_BUILT" : { "term": { "result": "NOT_BUILT" } }
				}
			  },
			  "aggs": {
				"total_duration": {
				  "sum": {
					"field": "duration"
				  }
				},
				"last_active": {
				  "scripted_metric": {
					"init_script": "state.lastActive = [];",
					"map_script": "def endTime = doc['start_time_in_millis'].value + doc['duration'].value; state.lastActive.add(endTime);",
					"combine_script": "return state.lastActive;",
					"reduce_script": "def maxEndTime = 0; for (lastActive in states) { for (endTime in lastActive) { maxEndTime = Math.max(maxEndTime, endTime); } } return maxEndTime;"
				  }
				}
			  }
			}
		  }
		}
	  }
  }`

const CiRunsExecutionWithoutRangeQuery = `{
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
		],
		"must_not": [
			{
		    	"term": {
					"result": "IN_PROGRESS"
		  		}
			}
		]
	  }
	},
	"aggs": {
		"completedRuns": {
		  "terms": {
			"field": "job_id",
			"size": 20000000
		  },
		  "aggs": {
			"endpoint_id": {
			  "terms": {
				"field": "endpoint_id",
				"size": 1
			  }
			},
			"result_buckets": {
			  "filters": {
				"filters": {
				  "SUCCESS": { "term": { "result": "SUCCESS" } },
				  "FAILED": { "term": { "result": "FAILURE" } },
				  "ABORTED": { "term": { "result": "ABORTED" } },
				  "UNSTABLE": { "term": { "result": "UNSTABLE" } },
				  "NOT_BUILT" : { "term": { "result": "NOT_BUILT" } }
				}
			  },
			  "aggs": {
				"total_duration": {
				  "sum": {
					"field": "duration"
				  }
				},
				"last_active": {
				  "scripted_metric": {
					"init_script": "state.lastActive = [];",
					"map_script": "def endTime = doc['start_time_in_millis'].value + doc['duration'].value; state.lastActive.add(endTime);",
					"combine_script": "return state.lastActive;",
					"reduce_script": "def maxEndTime = 0; for (lastActive in states) { for (endTime in lastActive) { maxEndTime = Math.max(maxEndTime, endTime); } } return maxEndTime;"
				  }
				}
			  }
			}
		  }
		}
	  }
  }`

const CiFragileJobRunsQuery = `{
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
		  }
		],
		"must_not": {
		  "term": {
			"result": "IN_PROGRESS"
		  }
		}
	  }
	},
	"aggs": {
		"completedRuns": {
		  "terms": {
			"field": "job_id",
			"size": 20000000
		  },
		  "aggs": {
			"endpoint_id": {
			  "terms": {
				"field": "endpoint_id",
				"size": 1
			  }
			},
			"result_buckets": {
			  "filters": {
				"filters": {
				  "SUCCESS": { "term": { "result": "SUCCESS" } },
				  "FAILED": { "term": { "result": "FAILURE" } },
				  "ABORTED": { "term": { "result": "ABORTED" } },
				  "UNSTABLE": { "term": { "result": "UNSTABLE" } },
				  "NOT_BUILT" : { "term": { "result": "NOT_BUILT" } }
				}
			  },
			  "aggs": {
				"total_duration": {
				  "sum": {
					"field": "duration"
				  }
				},
				"last_active": {
				  "scripted_metric": {
					"init_script": "state.lastActive = [];",
					"map_script": "def endTime = doc['start_time_in_millis'].value + doc['duration'].value; state.lastActive.add(endTime);",
					"combine_script": "return state.lastActive;",
					"reduce_script": "def maxEndTime = 0; for (lastActive in states) { for (endTime in lastActive) { maxEndTime = Math.max(maxEndTime, endTime); } } return maxEndTime;"
				  }
				}
			  }
			}
		  }
		}
	  }
  }`

const GetExecutedJobIds = `{
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
		  }
		],
		"must_not": {
		  "term": {
			"result": "IN_PROGRESS"
		  }
		}
	  }
	},
	"aggs": {
		"unique_jobs": {
		  "terms": {
			"field": "job_id",
			"size": 20000000 
		  }
		}
	  }
  }`

const CiJobInfoByJobIdQuery = `{
	"size": 1,
	"query": {
	  "bool": {
		"filter": [
		  {
			"term": {
			  "job_id": "{{.jobId}}"
			}
		  }
		]
	  }
	}
  }`

const CiRunsByJobIdQuery = `{
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
				"time_zone":"{{.timeZone}}"
			  }
			}
		  },
		  {
			"term": {
			  "job_id": "{{.jobId}}"
			}
		  }
		],
		"must_not": {
		  "term": {
			"result": "IN_PROGRESS"
		  }
		}
	  }
	},
	"aggs": {
	  "completedRuns": {
		"scripted_metric": {
			"params": {
				"timeZone": "{{.timeZone}}"
			},
		  "init_script": "state.dataMap = [:];",
		  "map_script": "def map = state.dataMap; def key = doc.org_id.value + '_' + doc.endpoint_id.value + '_' + doc.job_id.value + '_'+doc.run_id; def v = ['org_id': doc.org_id.value, 'endpoint_id': doc.endpoint_id.value, 'job_id': doc.job_id.value,'run_id': doc.run_id.value, 'result': doc.result.value,'duration': doc.duration.value, 'start_time': doc.start_time.value, 'start_time_in_millis': doc.start_time_in_millis.value,'timestamp':doc.timestamp.value,'url':doc.url.value]; map.put(key, v);",
		  "combine_script": "return state.dataMap;",
		  "reduce_script": "def tmpMap = [: ], resultList = new ArrayList();for (response in states) {if (response != null) {for (key in response.keySet()) {def record = response.get(key);if (record.result == 'FAILURE') {record.result = 'FAILED';}tmpMap.put(key, record);}}}for (key in tmpMap.keySet()) {def valueMap = tmpMap.get(key);def rd = valueMap.timestamp;valueMap.timestamp = rd.withZoneSameInstant(ZoneId.of(params.timeZone));def rd1 = valueMap.start_time;valueMap.start_time = rd1.withZoneSameInstant(ZoneId.of(params.timeZone));resultList.add(valueMap);}return resultList;"
		}
	  }
	}
  }`

const CiEndpointJobsQuery = `{
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
		]
	  }
	},
	"aggs": {
	  "endpointJobs": {
		"scripted_metric": {
		  "init_script": "state.dataMap = [:];",
		  "map_script": "def map = state.dataMap;def key = doc.org_id.value + '_' + doc.endpoint_id.value + '_' + doc.job_id.value + '_' + doc.job_name.value;def v = ['orgId': doc.org_id.value, 'endpointId': doc.endpoint_id.value, 'jobId': doc.job_id.value, 'jobName': doc.job_name.value, 'type': doc.type.value];map.put(key, v);",
		  "combine_script": "return state.dataMap;",
		  "reduce_script": "def tmpMap = [: ], resultList = new ArrayList();for (response in states) {if (response != null) {for (key in response.keySet()) {tmpMap.put(key, response.get(key));}}}def endpointJobMap = new HashMap();for (key in tmpMap.keySet()) {def record = tmpMap.get(key);if (endpointJobMap.containsKey(record.endpointId)) {def jobsList = endpointJobMap.get(record.endpointId);jobsList.add(record.jobId);endpointJobMap.put(record.endpointId, jobsList);} else {def jobsList = new HashSet();jobsList.add(record.jobId);endpointJobMap.put(record.endpointId, jobsList);}}return endpointJobMap;"
		}
	  }
	}
  }`

const CiJobRunsQuery = `{
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
		  }
		],
		"must_not": {
		  "term": {
			"result": "IN_PROGRESS"
		  }
		}
	  }
	},
	"aggs": {
		"completedRuns": {
		  "terms": {
			"field": "job_id",
			"size": 20000000
		  },
		  "aggs": {
			"endpoint_id": {
			  "terms": {
				"field": "endpoint_id",
				"size": 1
			  }
			},
			"result_buckets": {
			  "filters": {
				"filters": {
				  "SUCCESS": { "term": { "result": "SUCCESS" } },
				  "FAILED": { "term": { "result": "FAILURE" } },
				  "ABORTED": { "term": { "result": "ABORTED" } },
				  "UNSTABLE": { "term": { "result": "UNSTABLE" } },
				  "NOT_BUILT" : { "term": { "result": "NOT_BUILT" } }
				}
			  },
			  "aggs": {
				"total_duration": {
				  "sum": {
					"field": "duration"
				  }
				},
				"last_active": {
				  "scripted_metric": {
					"init_script": "state.lastActive = [];",
					"map_script": "def endTime = doc['start_time_in_millis'].value + doc['duration'].value; state.lastActive.add(endTime);",
					"combine_script": "return state.lastActive;",
					"reduce_script": "def maxEndTime = 0; for (lastActive in states) { for (endTime in lastActive) { maxEndTime = Math.max(maxEndTime, endTime); } } return maxEndTime;"
				  }
				}
			  }
			}
		  }
		}
	  }
  }`

const CiCjocControllersFetchQuery = `{
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
		]
	  }
	},
	"aggs": {
	  "cjocControllerInfo": {
		"scripted_metric": {
		  "init_script": "state.dataMap = [:];",
		  "map_script": "def map = state.dataMap; def key = doc.org_id.value + '_' + doc.endpoint_id.value + '_' + doc.url.value; def v = ['orgId': doc.org_id.value, 'endpointId': doc.endpoint_id.value, 'url': doc.url.value]; map.put(key, v);",
		  "combine_script": "return state.dataMap;",
		  "reduce_script": "def tmpMap = [: ];for (response in states) {if (response != null) {for (key in response.keySet()) {tmpMap.put(key, response.get(key));}}}def cjocControllerMap = new HashMap();for (key in tmpMap.keySet()) {def record = tmpMap.get(key);def controllerList = new ArrayList();if (cjocControllerMap.containsKey(record.endpointId)) {controllerList = cjocControllerMap.get(record.endpointId);}controllerList.add(record.url);cjocControllerMap.put(record.endpointId, controllerList);}return cjocControllerMap;"
		}
	  }
	}
  }`

const CiCjocControllersFetchByEndpoint = `{
	"size": 0,
	"query": {
	  "bool": {
		"filter": [
		]
	  }
	},
	"aggs": {
	  "cjocControllerInfo": {
		"scripted_metric": {
		  "init_script": "state.dataMap = [:];",
		  "map_script": "def map = state.dataMap; def key = doc.org_id.value + '_' + doc.endpoint_id.value + '_' + doc.url.value; def v = ['orgId': doc.org_id.value, 'endpointId': doc.endpoint_id.value, 'url': doc.url.value]; map.put(key, v);",
		  "combine_script": "return state.dataMap;",
		  "reduce_script": "def tmpMap = [: ];for (response in states) {if (response != null) {for (key in response.keySet()) {tmpMap.put(key, response.get(key));}}}def controllerUrlList = new ArrayList();for (key in tmpMap.keySet()) {def record = tmpMap.get(key);controllerUrlList.add(record.url);}return controllerUrlList;"
		}
	  }
	}
  }`

const CiJobCount = `{
	"query": {
		"bool": {
		  "filter": [
		  ]
		}
	  }
  }`

const CiRunCount = `{
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
			}
		  ],
		  "must_not": {
			"term": {
			  "result": "IN_PROGRESS"
			}
		  }
		}
	  }
  }`

const CiCompletedRunsFetchQueryCount = `{
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
				}
			],
			"must_not": {
				"term": {
					"result": "IN_PROGRESS"
				}
			}
		}
	}
}`

const CiAllJobCount = `{
  "size": 0,
  	"query": {
	  "bool": {
	  }
	},
  "aggs": {
    "jobs_per_endpoint": {
      "terms": {
        "field": "endpoint_id",
        "size": 20000000
      }
    }
  }
}`

const CiCompletedRunsResultQueryCount = `{
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
        }
      ],
      "must_not": {
        "term": {
          "result": "IN_PROGRESS"
        }
      }
    }
  },
  "aggs": {
    "result_counts": {
      "terms": {
        "field": "result",
        "size": 10
      }
    }
  },
  "size": 0
}`
