package constants

type Organization struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	SubOrgs    []*Organization `json:"sub_orgs"`
	Components []*Component    `json:"components"`
}

type Component struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// ParentOrg *Organization
}

type DeploymentFrequencyComponentComparison struct {
	Aggregations struct {
		DeploymentFrequencyComponentComparison struct {
			Buckets []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				DeployData struct {
					Value struct {
						Average        float64 `json:"average"`
						Deployments    int     `json:"deployments"`
						DifferenceDays float64 `json:"differenceDays"`
					} `json:"value"`
				} `json:"deploy_data"`
			} `json:"buckets"`
		} `json:"deployment_frequency_component_comparison"`
	} `json:"aggregations"`
}

type DeploymentLeadTimeComponentComparison struct {
	Aggregations struct {
		DeploymentLeadTimeComponentComparison struct {
			Buckets []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				DeployData struct {
					Value struct {
						TotalDuration int     `json:"totalDuration"`
						Average       float64 `json:"average"`
						Deployments   int     `json:"deployments"`
					} `json:"value"`
				} `json:"deploy_data"`
			} `json:"buckets"`
		} `json:"deployment_lead_time_component_comparison"`
	} `json:"aggregations"`
}

type DoraMttrComponentComparison struct {
	Aggregations struct {
		MttrComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key         string `json:"key"`
				DocCount    int    `json:"doc_count"`
				Deployments struct {
					Value struct {
						RecoveredTotalDuration int `json:"recoveredTotalDuration"`
						RecoveredCount         int `json:"recoveredCount"`
					} `json:"value"`
				} `json:"deployments"`
			} `json:"buckets"`
		} `json:"mttr_component_comparison"`
	} `json:"aggregations"`
}

type FailureRateComponentComparison struct {
	Aggregations struct {
		FailureRateComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key        string `json:"key"`
				DocCount   int    `json:"doc_count"`
				DeployData struct {
					Value struct {
						Average           string `json:"average"`
						Deployments       int    `json:"deployments"`
						FailedDeployments int    `json:"failedDeployments"`
					} `json:"value"`
				} `json:"deploy_data"`
			} `json:"buckets"`
		} `json:"failure_rate_component_comparison"`
	} `json:"aggregations"`
}

type OpenVulnerabilitiesOverviewComponentComparison struct {
	Aggregations struct {
		OpenVulnerabilitiesOverviewComponentComparison struct {
			Buckets []struct {
				Key                  string `json:"key"`
				DocCount             int    `json:"doc_count"`
				OpenVulSeverityCount struct {
					Value struct {
						VeryHigh int `json:"VERY_HIGH"`
						High     int `json:"HIGH"`
						Medium   int `json:"MEDIUM"`
						Low      int `json:"LOW"`
					} `json:"value"`
				} `json:"openVulSeverityCount"`
			} `json:"buckets"`
		} `json:"open_vulnerabilities_overview_component_comparison"`
	} `json:"aggregations"`
}

type VulnerabilitiesOverviewComponentComparison struct {
	Aggregations struct {
		VulnerabilitiesOverviewComponentComparison struct {
			Buckets []struct {
				Key                       string `json:"key"`
				DocCount                  int    `json:"doc_count"`
				VulnerabilityStatusCounts struct {
					Value struct {
						Reopened int `json:"Reopened"`
						Resolved int `json:"Resolved"`
						Found    int `json:"Found"`
						Open     int `json:"Open"`
					} `json:"value"`
				} `json:"vulnerabilityStatusCounts"`
			} `json:"buckets"`
		} `json:"vulnerabilities_overview_component_comparison"`
	} `json:"aggregations"`
}

type VelocityComponentComparison struct {
	Aggregations struct {
		FlowVelocityComponentComparison struct {
			Buckets []struct {
				Key               string `json:"key"`
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
		} `json:"flow_velocity_component_comparison"`
	} `json:"aggregations"`
}

type CompareReports struct {
	IsSubOrg       bool             `json:"is_sub_org"`
	SubOrgID       string           `json:"sub_org_id"`
	CompareTitle   string           `json:"compare_title"`
	SubOrgCount    int              `json:"sub_org_count"`
	ComponentCount int              `json:"component_count"`
	TotalValue     int              `json:"total_value"`
	ValueInMillis  float64          `json:"value_in_millis"`
	CompareReports []CompareReports `json:"compare_reports"`
	Section        struct {
		Data []struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		} `json:"data"`
	} `json:"section"`
}

type CycleTimeCompareReports struct {
	IsSubOrg       bool                      `json:"is_sub_org"`
	SubOrgID       string                    `json:"sub_org_id"`
	CompareTitle   string                    `json:"compare_title"`
	SubOrgCount    int                       `json:"sub_org_count"`
	ComponentCount int                       `json:"component_count"`
	TotalValue     int                       `json:"total_value"`
	ValueInMillis  float64                   `json:"value_in_millis"`
	CompareReports []CycleTimeCompareReports `json:"compare_reports"`
	Section        struct {
		Data []struct {
			Title string  `json:"title"`
			Value float64 `json:"value"`
			Time  int     `json:"time"`
			Count int     `json:"count"`
		} `json:"data"`
	} `json:"section"`
}

type DevCycleTimeCompareReports struct {
	IsSubOrg       bool                         `json:"is_sub_org"`
	SubOrgID       string                       `json:"sub_org_id"`
	CompareTitle   string                       `json:"compare_title"`
	SubOrgCount    int                          `json:"sub_org_count"`
	ComponentCount int                          `json:"component_count"`
	TotalValue     int                          `json:"total_value"`
	ValueInMillis  float64                      `json:"value_in_millis"`
	CompareReports []DevCycleTimeCompareReports `json:"compare_reports"`
	Section        struct {
		Data []struct {
			Title string  `json:"title"`
			Value float64 `json:"value"`
			Time  int     `json:"time"`
			Count int     `json:"count"`
		} `json:"data"`
	} `json:"section"`
}

type DeploymentFrequencyCompareReports struct {
	IsSubOrg       bool                                `json:"is_sub_org"`
	SubOrgID       string                              `json:"sub_org_id"`
	CompareTitle   string                              `json:"compare_title"`
	SubOrgCount    int                                 `json:"sub_org_count"`
	ComponentCount int                                 `json:"component_count"`
	TotalValue     int                                 `json:"total_value"`
	ValueInMillis  float64                             `json:"value_in_millis"`
	CompareReports []DeploymentFrequencyCompareReports `json:"compare_reports"`
	Section        struct {
		Data []struct {
			Title string  `json:"title"`
			Value float64 `json:"value"`
		} `json:"data"`
	} `json:"section"`
}

type ActiveTimeCompareReports struct {
	IsSubOrg       bool                       `json:"is_sub_org"`
	SubOrgID       string                     `json:"sub_org_id"`
	CompareTitle   string                     `json:"compare_title"`
	SubOrgCount    int                        `json:"sub_org_count"`
	ComponentCount int                        `json:"component_count"`
	TotalValue     int                        `json:"total_value"`
	ValueInMillis  float64                    `json:"value_in_millis"`
	CompareReports []ActiveTimeCompareReports `json:"compare_reports"`
	Section        struct {
		Data []struct {
			Title      string `json:"title"`
			Value      int    `json:"value"`
			ActiveTime int    `json:"active_time"`
			FlowTime   int    `json:"flow_time"`
		} `json:"data"`
	} `json:"section"`
}

type CycleTimeComponentComparison struct {
	Aggregations struct {
		CycleTimeComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key                string `json:"key"`
				DocCount           int    `json:"doc_count"`
				FlowCycleTimeCount struct {
					Value struct {
						FeatureTime   int `json:"FEATURE_TIME"`
						DefectTime    int `json:"DEFECT_TIME"`
						DefectCount   int `json:"DEFECT_COUNT"`
						TechDebtCount int `json:"TECH_DEBT_COUNT"`
						RiskTime      int `json:"RISK_TIME"`
						TechDebtTime  int `json:"TECH_DEBT_TIME"`
						FeatureCount  int `json:"FEATURE_COUNT"`
						RiskCount     int `json:"RISK_COUNT"`
					} `json:"value"`
				} `json:"flow_cycle_time_count"`
			} `json:"buckets"`
		} `json:"cycle_time_component_comparison"`
	} `json:"aggregations"`
}

type ActiveWorkTimeComponentComparison struct {
	Aggregations struct {
		FlowVelocityComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key                 string `json:"key"`
				DocCount            int    `json:"doc_count"`
				FlowEfficiencyCount struct {
					Value struct {
						TechDebt struct {
							ActiveTime int `json:"activeTime"`
							FlowTime   int `json:"flowTimeTime"`
						} `json:"TECH_DEBT"`
						Defect struct {
							ActiveTime int `json:"activeTime"`
							FlowTime   int `json:"flowTime"`
						} `json:"DEFECT"`
						Feature struct {
							ActiveTime int `json:"activeTime"`
							FlowTime   int `json:"flowTime"`
						} `json:"FEATURE"`
						Risk struct {
							ActiveTime int `json:"activeTime"`
							FlowTime   int `json:"flowTimeTime"`
						} `json:"RISK"`
					} `json:"value"`
				} `json:"flow_efficiency_count"`
			} `json:"buckets"`
		} `json:"active_work_time_component_comparison"`
	} `json:"aggregations"`
}

type WorkWaitTimeComponentComparison struct {
	Aggregations struct {
		FlowVelocityComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key                 string `json:"key"`
				DocCount            int    `json:"doc_count"`
				FlowEfficiencyCount struct {
					Value struct {
						TechDebt struct {
							WaitingTime int `json:"waitingTime"`
							FlowTime    int `json:"flowTime"`
						} `json:"TECH_DEBT"`
						Defect struct {
							WaitingTime int `json:"waitingTime"`
							FlowTime    int `json:"flowTime"`
						} `json:"DEFECT"`
						Feature struct {
							WaitingTime int `json:"waitingTime"`
							FlowTime    int `json:"flowTime"`
						} `json:"FEATURE"`
						Risk struct {
							WaitingTime int `json:"waitingTime"`
							FlowTime    int `json:"flowTime"`
						} `json:"RISK"`
					} `json:"value"`
				} `json:"flow_efficiency_count"`
			} `json:"buckets"`
		} `json:"work_wait_time_component_comparison"`
	} `json:"aggregations"`
}

type WorkloadComponentComparison struct {
	Aggregations struct {
		WorkloadComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key            string `json:"key"`
				DocCount       int    `json:"doc_count"`
				WorkLoadCounts struct {
					Value struct {
						HeaderValue int             `json:"headerValue"`
						Dates       map[string]Data `json:"dates"`
					} `json:"value"`
				} `json:"work_load_counts"`
			} `json:"buckets"`
		} `json:"workload_component_comparison"`
	} `json:"aggregations"`
}

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

type CommitsTrendsComponentComparison struct {
	Aggregations struct {
		CommitsTrendsComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key           string `json:"key"`
				DocCount      int    `json:"doc_count"`
				UniqueAuthors struct {
					Value int `json:"value"`
				} `json:"unique_authors"`
				CommitsCount struct {
					Value int `json:"value"`
				} `json:"commits_count"`
			} `json:"buckets"`
		} `json:"commits_trends_component_comparison"`
	} `json:"aggregations"`
}

type PullRequestComponentComparison struct {
	Aggregations struct {
		PullRequestsComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key          string `json:"key"`
				DocCount     int    `json:"doc_count"`
				Pullrequests struct {
					Value struct {
						ChangesRequested int `json:"CHANGES_REQUESTED"`
						Approved         int `json:"APPROVED"`
						Open             int `json:"OPEN"`
						Rejected         int `json:"REJECTED"`
					} `json:"value"`
				} `json:"pullrequests"`
			} `json:"buckets"`
		} `json:"pull_requests_component_comparison"`
	} `json:"aggregations"`
}

type WorkflowRunsComponentComparison struct {
	Aggregations struct {
		WorkflowRunsComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key           string `json:"key"`
				DocCount      int    `json:"doc_count"`
				AutomationRun struct {
					Value struct {
						Success int `json:"Success"`
						Failure int `json:"Failure"`
					} `json:"value"`
				} `json:"automation_run"`
			} `json:"buckets"`
		} `json:"workflow_runs_component_comparison"`
	} `json:"aggregations"`
}

type DevCycleTimeComponentComparison struct {
	Aggregations struct {
		DevelopmentCycleTimeComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key                  string `json:"key"`
				DocCount             int    `json:"doc_count"`
				DevelopmentCycleTime struct {
					Value struct {
						CodingTimeCount         int    `json:"coding_time_count"`
						CodingTime              string `json:"coding_time"`
						ReviewTimeCount         int    `json:"review_time_count"`
						ReviewTimeValueInMillis int    `json:"review_time_value_in_millis"`
						PickupTimeValueInMillis int    `json:"pickup_time_value_in_millis"`
						ReviewTime              string `json:"review_time"`
						CodingTimeValueInMillis int    `json:"coding_time_value_in_millis"`
						PickupTime              string `json:"pickup_time"`
						DeployTime              string `json:"deploy_time"`
						DeployTimeCount         int    `json:"deploy_time_count"`
						DeployTimeValueInMillis int    `json:"deploy_time_value_in_millis"`
						PickupTimeCount         int    `json:"pickup_time_count"`
					} `json:"value"`
				} `json:"developmentCycleTime"`
			} `json:"buckets"`
		} `json:"development_cycle_time_component_comparison"`
	} `json:"aggregations"`
}

type CommitsComponentComparison struct {
	Aggregations struct {
		CommitsComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key           string `json:"key"`
				DocCount      int    `json:"doc_count"`
				AutomationRun struct {
					Value struct {
						TotalCount int `json:"totalCount"`
					} `json:"value"`
				} `json:"automation_run"`
			} `json:"buckets"`
		} `json:"commits_component_comparison"`
	} `json:"aggregations"`
}

type BuildsComponentComparison struct {
	Aggregations struct {
		BuildsComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key         string `json:"key"`
				DocCount    int    `json:"doc_count"`
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
			} `json:"buckets"`
		} `json:"builds_component_comparison"`
	} `json:"aggregations"`
}

type DeploymentsComponentComparison struct {
	Aggregations struct {
		DeploymentsComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
				Deploys  struct {
					Value int `json:"value"`
				} `json:"deploys"`
			} `json:"buckets"`
		} `json:"deployments_component_comparison"`
	} `json:"aggregations"`
}

type ComponentsComponentComparison struct {
	Aggregations struct {
		DistinctComponent struct {
			Value []string `json:"value"`
		} `json:"distinct_component"`
	} `json:"aggregations"`
}

type WorkflowsComponentComparison struct {
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
}

type SecurityWorkflowsComponentComparison struct {
	WithScanners    int `json:"withScanners"`
	WithoutScanners int `json:"withoutScanners"`
}

type TestWorkflowsComponentComparison struct {
	WithTestSuites    int `json:"withTestSuites"`
	WithoutTestSuites int `json:"withoutTestSuites"`
}

type SecurityWorkflowRunsComponentComparison struct {
	Aggregations struct {
		WorkflowRunsComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key       string `json:"key"`
				DocCount  int    `json:"doc_count"`
				RunStatus struct {
					Value struct {
						ChartData struct {
							Info []struct {
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
			} `json:"buckets"`
		} `json:"workflow_runs_component_comparison"`
	} `json:"aggregations"`
}

type MttrComponentComparison struct {
	Aggregations struct {
		MttrComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
				AvgTTR   struct {
					Value struct {
						VeryHigh int `json:"VERY_HIGH"`
						High     int `json:"HIGH"`
						Medium   int `json:"MEDIUM"`
						Low      int `json:"LOW"`
					} `json:"value"`
				} `json:"Avg_TTR"`
			} `json:"buckets"`
		} `json:"mttr_component_comparison"`
	} `json:"aggregations"`
}

type MttrCompareReports struct {
	IsSubOrg       bool                 `json:"is_sub_org"`
	SubOrgID       string               `json:"sub_org_id"`
	CompareTitle   string               `json:"compare_title"`
	SubOrgCount    int                  `json:"sub_org_count"`
	ComponentCount int                  `json:"component_count"`
	TotalValue     int                  `json:"total_value"`
	ValueInMillis  float64              `json:"value_in_millis"`
	CompareReports []MttrCompareReports `json:"compare_reports"`
	AverageValue   int                  `json:"-"`
	Count          int                  `json:"-"`
}

type VulnerabilitesByScannerType struct {
	Aggregations struct {
		VulByScannerTypeComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key                    string `json:"key"`
				DocCount               int    `json:"doc_count"`
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
			} `json:"buckets"`
		} `json:"vul_by_scanner_type_component_comparison"`
	} `json:"aggregations"`
}

type TestSuiteAutoRun struct {
	Aggregations struct {
		WorkflowRunsComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key       string `json:"key"`
				DocCount  int    `json:"doc_count"`
				RunStatus struct {
					Value struct {
						Runs int `json:"runs"`
					} `json:"value"`
				} `json:"run_status"`
			} `json:"buckets"`
		} `json:"workflow_runs_component_comparison"`
	} `json:"aggregations"`
}

type TestSuitesRun struct {
	Aggregations struct {
		WorkflowRunsComponentComparison struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key               string `json:"key"`
				DocCount          int    `json:"doc_count"`
				ComponentActivity struct {
					Value struct {
						Runs int `json:"runs"`
					} `json:"value"`
				} `json:"component_activity"`
			} `json:"buckets"`
		} `json:"workflow_runs_component_comparison"`
	} `json:"aggregations"`
}

type TestSuite struct {
	WithTestSuite    int
	WithoutTestSuite int
}

type TestSuiteComponentComparison struct {
	Val []map[string]TestSuite
}

type TestComponentsComponentComparison struct {
	Aggregations struct {
		Components struct {
			Value []string `json:"value"`
		} `json:"components"`
	} `json:"aggregations"`
}
