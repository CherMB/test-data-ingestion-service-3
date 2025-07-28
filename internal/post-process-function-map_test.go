package internal

import (
	"encoding/json"
	"sort"

	"testing"

	"github.com/calculi-corp/log"
	"github.com/calculi-corp/reports-service/constants"
	db "github.com/calculi-corp/reports-service/db"
	"github.com/opensearch-project/opensearch-go"
	"github.com/stretchr/testify/assert"
)

func TestExecutePostProcessFunctionForComponentComparison(t *testing.T) {
	organization := &constants.Organization{}

	originalPostProcessFunctionMap := PostProcessFunctionMap
	defer func() { PostProcessFunctionMap = originalPostProcessFunctionMap }()

	tests := []struct {
		name          string
		k             string
		specKey       string
		data          map[string]json.RawMessage
		replacements  map[string]interface{}
		organization  *constants.Organization
		wantRawResult json.RawMessage
		wantErr       bool
	}{
		{
			name: "Successful execution of mockFunction",
			k:    "mockFunction",
			data: map[string]json.RawMessage{
				"exampleData": json.RawMessage(`{"example": "data"}`),
			},
			replacements: map[string]interface{}{
				"exampleReplacement": "value",
			},
			organization:  organization,
			wantRawResult: json.RawMessage(`{"result": "mocked data"}`),
			wantErr:       false,
		},
		{
			name:          "Non-existent function",
			k:             "nonExistentFunction",
			wantRawResult: nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExecutePostProcessFunctionForComponentComparison(tt.k, tt.specKey, tt.data, tt.replacements, tt.organization)
		})
	}
}

func TestSumTotalValue(t *testing.T) {

	mockCompareReports := &constants.CompareReports{
		IsSubOrg:   false,
		TotalValue: 100,
		CompareReports: []constants.CompareReports{
			{
				IsSubOrg:   true,
				TotalValue: 50,
			},
			{
				IsSubOrg:   false,
				TotalValue: 25,
				CompareReports: []constants.CompareReports{
					{
						IsSubOrg:   true,
						TotalValue: 10,
					},
					{
						IsSubOrg:   false,
						TotalValue: 15,
					},
				},
			},
		},
	}

	totalValue := sumTotalValue(mockCompareReports, 0)

	assert.Equal(t, 140, totalValue, "Expected totalValue to be 200")

	assert.Equal(t, 140, mockCompareReports.TotalValue, "Expected TotalValue to be updated correctly")
}

func TestSumCounts(t *testing.T) {
	mockData := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{Title: "Bugs", Value: 5},
		{Title: "Feature", Value: 10},
		{Title: "Risk", Value: 3},
		{Title: "Tech Debt", Value: 7},
	}
	mockSection := struct {
		Data []struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		} `json:"data"`
	}{
		Data: mockData,
	}
	mockCompareReports := &constants.CompareReports{
		IsSubOrg: false,
		Section:  mockSection,
	}

	featureSum, riskSum, techDebtSum, bugSum := sumCounts(mockCompareReports, 0, 0, 0, 0)

	assert.Equal(t, int64(10), featureSum, "Expected featureSum to be 10")
	assert.Equal(t, int64(3), riskSum, "Expected riskSum to be 3")
	assert.Equal(t, int64(7), techDebtSum, "Expected techDebtSum to be 7")
	assert.Equal(t, int64(5), bugSum, "Expected bugSum to be 5")

	expectedData := []struct {
		Title string `json:"title"`
		Value int    `json:"value"`
	}{
		{Title: "Bugs", Value: 5},
		{Title: "Feature", Value: 10},
		{Title: "Risk", Value: 3},
		{Title: "Tech Debt", Value: 7},
	}
	assert.Equal(t, expectedData, mockCompareReports.Section.Data, "Expected Section.Data to be updated correctly")
}

func Test_getCommitsActivity(t *testing.T) {

	t.Run("Case 1: Successful execution of component commit activity", func(t *testing.T) {

		responseString := `{
			"commitsActivity":{
				"aggregations": {
			  		"commits": {
						"value": {
				  			"avg": {
								"title": "Commits / active dev",
								"value": 45
				  			},
				  			"dev": {
								"title": "Active Developers",
								"value": 92
				  			},
				 			 "commits_count": 4165
						}
			  		}
				}
			}
		}`

		x := map[string]json.RawMessage{}

		json.Unmarshal([]byte(responseString), &x)

		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		expectResult := []byte(`{"value":4165}`)
		b, err := getCommitsActivity("header", x, nil)
		assert.Nil(t, err, "error processing getCommitsActivity header")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`[{"title":"Commits / active dev","value":45},{"title":"Active Developers","value":92}]`)
		b, err = getCommitsActivity("section", x, nil)
		assert.Nil(t, err, "error processing getCommitsActivity section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_getRunsActivity(t *testing.T) {

	t.Run("Case 1: Successful execution of component runs activity", func(t *testing.T) {

		responseString := `{"runsActivity":{"took":56,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":4797,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"automation_run":{"value":{"data":[{"title":"Success","value":3886},{"title":"Failure","value":619}],"totalCount":4505}}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":4505}`)
		b, err := getRunsActivity("header", x, nil)
		assert.Nil(t, err, "error processing getRunsActivity header")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`[{"title":"Success","value":3886},{"title":"Failure","value":619}]`)
		b, err = getRunsActivity("section", x, nil)
		assert.Nil(t, err, "error processing getRunsActivity section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_getBuildsMetric(t *testing.T) {

	t.Run("Case 1: Successful execution of Builds Info", func(t *testing.T) {
		responseString := `{"buildsData":{"took":48,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":9788,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"build_status":{"value":{"total_builds":3194,"info":[{"drillDown":{"reportType":"status","reportId":"component-summary-builds","reportTitle":"Builds"},"title":"Success","value":3140},{"drillDown":{"reportType":"status","reportId":"component-summary-builds","reportTitle":"Builds"},"title":"Failure","value":54}]}}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":3194}`)
		b, err := getBuildsMetric("header", x, nil)
		assert.Nil(t, err, "error processing getBuildsMetric header")
		assert.Equal(t, expectResult, []byte(b))

		//expectResult = []byte(`{"subHeader":[{"drillDown":{"reportType":"status","reportId":"component-summary-builds","reportTitle":"Builds"},"title":"Success","value":3140},{"drillDown":{"reportType":"status","reportId":"component-summary-builds","reportTitle":"Builds"},"title":"Failure","value":54}]}`)
		_, err = getBuildsMetric("subHeader", x, nil)
		assert.Nil(t, err, "error processing getBuildsMetric subHeader")
	})

	t.Run("Case 1: Successful execution of Builds Chart", func(t *testing.T) {

		responseString := `{"buildsDataChart":{"took":1,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":63,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"runs_buckets":{"buckets":[{"key_as_string":"2024-05-27","key":1716768000000,"doc_count":0,"automation_run":{"value":{"Success":0,"Failure":0}}},{"key_as_string":"2024-06-03","key":1717372800000,"doc_count":63,"automation_run":{"value":{"Success":52,"Failure":11}}},{"key_as_string":"2024-06-10","key":1717977600000,"doc_count":0,"automation_run":{"value":{"Success":0,"Failure":0}}},{"key_as_string":"2024-06-17","key":1718582400000,"doc_count":0,"automation_run":{"value":{"Success":0,"Failure":0}}},{"key_as_string":"2024-06-24","key":1719187200000,"doc_count":0,"automation_run":{"value":{"Success":0,"Failure":0}}}]}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"id":"Success","data":[{"x":"2024-05-27","y":0},{"x":"2024-06-03","y":52},{"x":"2024-06-10","y":0},{"x":"2024-06-17","y":0},{"x":"2024-06-24","y":0}]},{"id":"Failure","data":[{"x":"2024-05-27","y":0},{"x":"2024-06-03","y":11},{"x":"2024-06-10","y":0},{"x":"2024-06-17","y":0},{"x":"2024-06-24","y":0}]}]`)
		b, err := getBuildsMetric("sectionChart", x, nil)
		assert.Nil(t, err, "error processing getBuildsMetric subHeader")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetDeployments(t *testing.T) {
	t.Run("Case 1: Successful execution of component summary deployments", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"deploymentSuccessRateHeader":{"took":1,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":57,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deploy_data":{"value":{"total":57,"data":[{"title":"Success","value":34},{"title":"Failure","value":23}],"value":"60%"}}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)
		expectResult := []byte(`{"value":"60%"}`)

		b, err := getDeployments("deploymentSuccessRateHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetDeployments header")
		assert.JSONEq(t, string(expectResult), string(b), "Header response did not match the expected result")

		//SubHeader
		expectResult = []byte(`{"subHeader":[{"title":"Success","value":34},{"title":"Failure","value":23}]}`)

		b, err = getDeployments("deploymentSuccessRateSubHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetDeployments sub header")
		assert.Equal(t, expectResult, []byte(b))

		//Section
		responseString = `{"deploymentData":{"took":2,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":57,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deploy_buckets":{"buckets":[{"key_as_string":"2024-05-27","key":1716748200000,"doc_count":0,"deploy_data":{"value":{"Success":0,"Failure":0}}},{"key_as_string":"2024-06-03","key":1717353000000,"doc_count":56,"deploy_data":{"value":{"Success":34,"Failure":22}}},{"key_as_string":"2024-06-10","key":1717957800000,"doc_count":0,"deploy_data":{"value":{"Success":0,"Failure":0}}},{"key_as_string":"2024-06-17","key":1718562600000,"doc_count":0,"deploy_data":{"value":{"Success":0,"Failure":0}}},{"key_as_string":"2024-06-24","key":1719167400000,"doc_count":1,"deploy_data":{"value":{"Success":0,"Failure":1}}}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)
		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"x":"2024-06-03","y":34},{"x":"2024-06-10","y":0},{"x":"2024-06-17","y":0},{"x":"2024-06-24","y":0}],"id":"Success"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"x":"2024-06-03","y":22},{"x":"2024-06-10","y":0},{"x":"2024-06-17","y":0},{"x":"2024-06-24","y":1}],"id":"Failure"}]`)
		b, err = getDeployments("deploymentDataSpec", x, replacements)

		assert.Nil(t, err, "error processing GetDeployments section")
		assert.Equal(t, expectResult, []byte(b))

	})
}

func Test_DevelopmentCycleTime(t *testing.T) {
	t.Run("Case 1: Successful execution of Development Cycle Time", func(t *testing.T) {
		responseString := `{"avgDevelopmentHeader":{"took":33,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":3194,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"developmentCycleTime":{"value":{"total":"9d 7h 28m ","value":804531000}}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":"9d 7h 28m ","valueInMillis":804531000}`)
		b, err := averageDevelopmentCycleTimeHeader("", x, nil)
		assert.Nil(t, err, "error processing averageDevelopmentCycleTimeHeader ")
		assert.Equal(t, expectResult, []byte(b))

		//Chart
		responseString = `{"developmentCycleChart":{"took":258,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":3207,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"developmentCycleTime":{"value":{"coding_time":14,"review_time":5,"pickup_time":81,"deploy_time":0}}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"title":"Coding time","value":14},{"title":"Code pickup time","value":81},{"title":"Code review time","value":5}]`)
		b, err = developmentCycleChartSection("", x, nil)
		assert.Nil(t, err, "error processing developmentCycleChartSection")
		assert.Equal(t, expectResult, []byte(b))

		//Footer
		responseString = `{"developmentTimeFooter":{"took":32,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":3218,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"developmentCycleTime":{"value":{"coding_time":"1d 8h 16m ","pickup_time_in_millis":653659000,"deploy_time_in_millis":0,"coding_time_in_millis":116198000,"review_time":"11h 1m ","pickup_time":"7d 13h 34m ","review_time_in_millis":39684000,"deploy_time":""}}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`{"value":"1d 8h 16m ","valueInMillis":116198000}`)
		b, err = codingTimeFooterSection("codingTimeSpec", x, nil)
		assert.Nil(t, err, "error processing Coding time ")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":"7d 13h 34m ","valueInMillis":653659000}`)
		b, err = codingTimeFooterSection("codingPickupTimeSpec", x, nil)
		assert.Nil(t, err, "error processing Coding time ")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":"11h 1m ","valueInMillis":39684000}`)
		b, err = codingTimeFooterSection("codingReviewTimeSpec", x, nil)
		assert.Nil(t, err, "error processing Review time ")
		assert.Equal(t, expectResult, []byte(b))

	})
}

func Test_getCommitTrends(t *testing.T) {
	t.Run("Case 1: Successful execution of Commits Trend", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-30",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-30",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"normalizeMonthInSpec": "2024-02-28",
			"commitTitle":          "Weekly commits/active devs",
		}

		responseString := `{"totalCommitsHeader":{"took":31,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":5712,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"commits_trends_widget":{"buckets":[{"key":"2024-05-01 00:00:00-2024-05-31 23:59:59","from":1714501800000,"from_as_string":"2024-05-01 00:00:00","to":1717180199000,"to_as_string":"2024-05-31 23:59:59","doc_count":5712,"commits_trend_headers":{"value":{"commits-per-author":44.929825,"unique_authors":114,"commits_count":5122}}}]}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":5122}`)
		b, err := getCommitTrends("totalCommitsHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing getCommitTrends")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"subHeader":[{"drillDown":{"reportId":"activeDevelopers","reportTitle":"Active developers","reportType":"status"},"title":"Active developers","value":114},{"title":"Weekly commits/active devs","value":44}]}`)
		b, err = getCommitTrends("averageCommitsHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing getCommitTrends")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"commitsAndAverageChart":{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":7472,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"commits_trends_widget":{"buckets":[{"key_as_string":"2024-02-26","key":1708885800000,"doc_count":1654,"unique_authors":{"value":44},"commits_count":{"value":1654},"commits-per-auth":{"value":37.59090909090909}},{"key_as_string":"2024-03-04","key":1709490600000,"doc_count":4735,"unique_authors":{"value":89},"commits_count":{"value":4735},"commits-per-auth":{"value":53.20224719101124}},{"key_as_string":"2024-03-11","key":1710095400000,"doc_count":1083,"unique_authors":{"value":49},"commits_count":{"value":1083},"commits-per-auth":{"value":22.102040816326532}},{"key_as_string":"2024-03-18","key":1710700200000,"doc_count":0,"unique_authors":{"value":0},"commits_count":{"value":0}},{"key_as_string":"2024-03-25","key":1711305000000,"doc_count":0,"unique_authors":{"value":0},"commits_count":{"value":0}}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-02-28"},"x":"2024-02-28","y":1654},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":4735},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":1083},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-30","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Commits","type":"line"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-02-28"},"x":"2024-02-28","y":37},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":53},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":22},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-30","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Weekly commits/active devs","isClickDisable":true,"type":"line"}]`)
		b, err = getCommitTrends("commitsAndAverageChartSpec", x, replacements)
		assert.Nil(t, err, "error processing getCommitTrends")
		assert.Equal(t, expectResult, []byte(b))

		//Check with Normalized date set a "@x"
		replacements["normalizeMonthInSpec"] = "@x"
		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":1654},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":4735},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":1083},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-30","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Commits","type":"line"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":37},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":53},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":22},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-30","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Weekly commits/active devs","isClickDisable":true,"type":"line"}]`)
		b, err = getCommitTrends("commitsAndAverageChartSpec", x, replacements)
		assert.Nil(t, err, "error processing getCommitTrends")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetAutomationRuns(t *testing.T) {
	t.Run("Case 1: Successful execution of Automation runs", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"totalRunsSubHeader":{"took":248,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":6676,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"automation_run":{"value":{"data":[{"title":"Success","value":5354},{"title":"Failure","value":897}]}}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":6251}`)
		b, err := getAutomationRuns("totalRunsHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetAutomationRuns header")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"subHeader":[{"drillDown":{"reportId":"workflowRuns","reportTitle":"Workflow runs","reportType":"status"},"title":"Success","value":5354},{"drillDown":{"reportId":"workflowRuns","reportTitle":"Workflow runs","reportType":"status"},"title":"Failure","value":897}]}`)
		b, err = getAutomationRuns("totalRunsSubHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetAutomationRuns sub header")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"runsStatusChart":{"took":101,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":6679,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"runs_buckets":{"buckets":[{"key_as_string":"2024-02-26","key":1708885800000,"doc_count":1328,"automation_run":{"value":{"Success":1095,"Failure":162}}},{"key_as_string":"2024-03-04","key":1709490600000,"doc_count":3664,"automation_run":{"value":{"Success":2953,"Failure":479}}},{"key_as_string":"2024-03-11","key":1710095400000,"doc_count":1687,"automation_run":{"value":{"Success":1309,"Failure":257}}},{"key_as_string":"2024-03-18","key":1710700200000,"doc_count":0,"automation_run":{"value":{"Success":0,"Failure":0}}},{"key_as_string":"2024-03-25","key":1711305000000,"doc_count":0,"automation_run":{"value":{"Success":0,"Failure":0}}}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":1095},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":2953},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":1309},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Success"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":162},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":479},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":257},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Failure"}]`)
		b, err = getAutomationRuns("runsStatusChartSpec", x, replacements)
		assert.Nil(t, err, "error processing GetAutomationRuns section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetPullRequests(t *testing.T) {
	t.Run("Case 1: Successful execution of Pull requests", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		//Header
		responseString := `{"totalPullRequestsHeader":{"took":8,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1158,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"by_repos":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"cloudbees/nextgen-ui","doc_count":237,"pr_count":{"value":61}}]},"sum_prs":{"value":416}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":416}`)
		b, err := getPullRequests("totalPullRequestsHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetPullRequests header")
		assert.Equal(t, expectResult, []byte(b))

		//Section
		responseString = `{"pullRequestsChart":{"took":2,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"aggregations":{"date_counts":{"buckets":[{"key_as_string":"2025-01-25","key":1737743400000,"doc_count":85,"by_status":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"OPEN","doc_count":81,"unique_pr_ids":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"1442","doc_count":22},{"key":"1440","doc_count":21},{"key":"1441","doc_count":20},{"key":"1443","doc_count":18}]}},{"key":"APPROVED","doc_count":4,"unique_pr_ids":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"1440","doc_count":2},{"key":"1442","doc_count":1},{"key":"1443","doc_count":1}]}}]}}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":3}],"id":"Approved"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0}],"id":"Changes requested"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":4}],"id":"Open"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0}],"id":"Rejected"}]`)
		b, err = getPullRequests("pullRequestsChartSpec", x, replacements)
		assert.Nil(t, err, "error processing GetPullRequests section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetAutomationRunsWithScanner(t *testing.T) {
	t.Run("Case 1: Successful execution of Automation runs with scanners", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"totalAutomationRuns":{"took":205,"timed_out":false,"_shards":{"total":4,"successful":4,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"run_status":{"value":{"chartData":{"data":[{"name":"With Scanners","value":9},{"name":"Without Scanners","value":90}],"info":[{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"With Scanners","value":567},{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"Without Scanners","value":5723}]},"Total":{"value":6290,"key":"Total"}}}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":6290}`)
		b, err := getAutomationRunsWithScanners("totalRunsSpec", x, replacements)
		assert.Nil(t, err, "error processing GetAutomationRuns with scanner header")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"totalAutomationRuns":{"took":205,"timed_out":false,"_shards":{"total":4,"successful":4,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"run_status":{"value":{"chartData":{"data":[{"name":"With Scanners","value":9},{"name":"Without Scanners","value":90}],"info":[{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"With Scanners","value":567},{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"Without Scanners","value":5723}]},"Total":{"value":6290,"key":"Total"}}}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`{"data":[{"name":"With Scanners","value":9},{"name":"Without Scanners","value":90}],"info":[{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"With Scanners","value":567},{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"Without Scanners","value":5723}]}`)
		b, err = getAutomationRunsWithScanners("runsStatusChartSpec", x, replacements)
		assert.Nil(t, err, "error processing GetAutomationRuns section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetVulnerabilitiesOverview(t *testing.T) {
	t.Run("Case 1: Successful execution of Vulnerabilities overview", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		//Header
		responseString := `{"vulnerabilityStatusCounts":{"took":734,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"vulnerabilityStatusCounts":{"value":{"Reopened":4,"Resolved":24,"Found":1989,"Open":1961}}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":1989}`)
		b, err := getVulnerabilitiesOverview("foundVulHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetVulnerabilitiesOverview header - found")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":4}`)
		b, err = getVulnerabilitiesOverview("reopenedVulHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetVulnerabilitiesOverview header - reopened")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":24}`)
		b, err = getVulnerabilitiesOverview("resolvedVulHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetVulnerabilitiesOverview header - resolved")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":1961}`)
		b, err = getVulnerabilitiesOverview("openVulHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetVulnerabilitiesOverview header - open")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"vulOverviewChart":{"took":773,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"vul_overview_buckets":{"buckets":[{"key_as_string":"2024-02-26","key":1708885800000,"doc_count":9421,"vul_overview_chart":{"value":{"Reopened":0,"Resolved":11,"Found":1727,"Open":1716}}},{"key_as_string":"2024-03-04","key":1709490600000,"doc_count":13816,"vul_overview_chart":{"value":{"Reopened":0,"Resolved":1,"Found":1349,"Open":1348}}},{"key_as_string":"2024-03-11","key":1710095400000,"doc_count":3643,"vul_overview_chart":{"value":{"Reopened":4,"Resolved":16,"Found":194,"Open":174}}},{"key_as_string":"2024-03-18","key":1710700200000,"doc_count":0,"vul_overview_chart":{"value":{"Reopened":0,"Resolved":0,"Found":0,"Open":0}}},{"key_as_string":"2024-03-25","key":1711305000000,"doc_count":0,"vul_overview_chart":{"value":{"Reopened":0,"Resolved":0,"Found":0,"Open":0}}}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":1727},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":1349},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":194},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Found"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":1716},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":1348},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":174},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Open"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":4},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Reopened"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":11},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":1},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":16},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Resolved"}]`)
		b, err = getVulnerabilitiesOverview("vulOverviewChartSpec", x, replacements)
		assert.Nil(t, err, "error processing GetVulnerabilitiesOverview section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetOpenVulnerabilitiesOverview(t *testing.T) {
	t.Run("Case 1: Successful execution of open vulnerabilities overview", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		//Header
		responseString := `{"openVulSeverityCount":{"took":1201,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"severityCounts":{"value":{"VERY_HIGH":41,"HIGH":90,"MEDIUM":181,"LOW":1653}}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":41}`)
		b, err := getOpenVulnerabilitiesOverview("veryHighSeverityHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetOpenVulnerabilitiesOverview header - very high")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":90}`)
		b, err = getOpenVulnerabilitiesOverview("highSeverityHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetOpenVulnerabilitiesOverview header - high")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":181}`)
		b, err = getOpenVulnerabilitiesOverview("mediumSeverityHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetOpenVulnerabilitiesOverview header - medium")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":1653}`)
		b, err = getOpenVulnerabilitiesOverview("lowSeverityHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetOpenVulnerabilitiesOverview header - low")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"openVulAgeChart":{"took":648,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"ageCounts":{"value":[{"id":"Very high","value":[3847936,109357440,269859936,5795293700,12194959936]},{"id":"High","value":[3847936,591459970,1143053950,5452564000,12197734936]},{"id":"Medium","value":[176130936,710921920,1143081980,1143184900,1322826936]},{"id":"Low","value":[3758936,526933952,756313920,1143081980,1817366936]}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"id":"Very high","value":[3847936,109357440,269859936,5795293700,12194959936]},{"id":"High","value":[3847936,591459970,1143053950,5452564000,12197734936]},{"id":"Medium","value":[176130936,710921920,1143081980,1143184900,1322826936]},{"id":"Low","value":[3758936,526933952,756313920,1143081980,1817366936]}]`)
		b, err = getOpenVulnerabilitiesOverview("openVulAgeChartSpec", x, replacements)
		assert.Nil(t, err, "error processing GetOpenVulnerabilitiesOverview section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetWorkload(t *testing.T) {
	t.Run("Case 1: Successful execution of workload", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		//Header
		responseString := `{"flowWorkLoad":{"took":874,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"work_load_counts":{"value":{"headerValue":90,"dates":{"2024-03-01":{"DEFECT":11,"FEATURE":29,"RISK":0,"TECH_DEBT":0,"DEFECT_SET":["SDP-11950","SDP-10158","SDP-13784","SDP-9573","SDP-14056","SDP-13175","SDP-13172","SDP-10429","SDP-11847","SDP-13171","SDP-14180"],"FEATURE_SET":["SDP-14206","SDP-14204","SDP-14267","SDP-14266","SDP-14287","SDP-7485","SDP-12706","SDP-14509","SDP-12705","SDP-13759","SDP-9846","SDP-11996","SDP-14274","SDP-14273","SDP-14071","SDP-14414","SDP-14258","SDP-6382","SDP-12233","SDP-9037","SDP-13627","SDP-13208","SDP-8865","SDP-13209","SDP-9418","SDP-14264","SDP-14120","SDP-14262","SDP-13050"],"RISK_SET":[],"TECH_DEBT_SET":[]},"2024-03-04":{"DEFECT":24,"FEATURE":39,"RISK":0,"TECH_DEBT":0,"DEFECT_SET":["SDP-14668","SDP-11950","SDP-14612","SDP-10158","SDP-14610","SDP-14663","SDP-14410","SDP-14597","SDP-14661","SDP-14452","SDP-9573","SDP-14056","SDP-14609","SDP-10429","SDP-11847","SDP-6554","SDP-14703","SDP-14625","SDP-13175","SDP-14252","SDP-13172","SDP-13171","SDP-14180","SDP-13191"],"FEATURE_SET":["SDP-14646","SDP-14206","SDP-11794","SDP-11970","SDP-14268","SDP-14266","SDP-14287","SDP-7485","SDP-12706","SDP-7686","SDP-12705","SDP-13759","SDP-14726","SDP-9846","SDP-11996","SDP-14274","SDP-14273","SDP-14071","SDP-14414","SDP-14238","SDP-12235","SDP-14236","SDP-14258","SDP-12233","SDP-14311","SDP-14233","SDP-13660","SDP-9037","SDP-13627","SDP-13208","SDP-8865","SDP-13209","SDP-9418","SDP-14242","SDP-14121","SDP-14264","SDP-14120","SDP-14262","SDP-13050"],"RISK_SET":[],"TECH_DEBT_SET":[]},"2024-03-11":{"DEFECT":18,"FEATURE":27,"RISK":0,"TECH_DEBT":0,"DEFECT_SET":["SDP-14448","SDP-14612","SDP-14875","SDP-10158","SDP-11971","SDP-14533","SDP-14872","SDP-14740","SDP-9573","SDP-14056","SDP-10429","SDP-14948","SDP-11847","SDP-13704","SDP-14727","SDP-14627","SDP-14693","SDP-13191"],"FEATURE_SET":["SDP-14800","SDP-14722","SDP-11794","SDP-11970","SDP-14268","SDP-7485","SDP-12706","SDP-14805","SDP-12705","SDP-14804","SDP-13759","SDP-6477","SDP-14803","SDP-9846","SDP-11996","SDP-14274","SDP-14071","SDP-14019","SDP-14798","SDP-14874","SDP-9037","SDP-13208","SDP-8865","SDP-14416","SDP-13800","SDP-14120","SDP-13050"],"RISK_SET":[],"TECH_DEBT_SET":[]},"2024-03-18":{"DEFECT":5,"FEATURE":13,"RISK":0,"TECH_DEBT":0,"DEFECT_SET":["SDP-14875","SDP-10158","SDP-9573","SDP-10429","SDP-11847"],"FEATURE_SET":["SDP-14019","SDP-11794","SDP-7485","SDP-12706","SDP-12705","SDP-13759","SDP-13208","SDP-8865","SDP-9846","SDP-11996","SDP-13800","SDP-14120","SDP-13050"],"RISK_SET":[],"TECH_DEBT_SET":[]},"2024-03-25":{"DEFECT":5,"FEATURE":13,"RISK":0,"TECH_DEBT":0,"DEFECT_SET":["SDP-14875","SDP-10158","SDP-9573","SDP-10429","SDP-11847"],"FEATURE_SET":["SDP-14019","SDP-11794","SDP-7485","SDP-12706","SDP-12705","SDP-13759","SDP-13208","SDP-8865","SDP-9846","SDP-11996","SDP-13800","SDP-14120","SDP-13050"],"RISK_SET":[],"TECH_DEBT_SET":[]}}}}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":90}`)
		b, err := getWorkload("flowWorkLoadHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing workload header")
		assert.Equal(t, expectResult, []byte(b))

		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":11},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":24},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":18},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":5},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":5}],"id":"Defect"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":29},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":39},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":27},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":13},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":13}],"id":"Feature"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Risk"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Tech debt"}]`)
		b, err = getWorkload("flowWorkLoadChartSpec", x, replacements)
		assert.Nil(t, err, "error processing workload section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetCycleTime(t *testing.T) {
	t.Run("Case 1: Successful execution of cycle time", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		//Header

		responseString := `{"flowCycleTimeHeader":{"took":35,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":837,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"flow_cycle_time_count":{"value":{"value":"  15d 18h 43m","valueInMillis":1363420806}}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":"  15d 18h 43m"}`)
		b, err := getCycleTime("flowCycleTimeHeaderSpec", x, replacements)

		assert.Nil(t, err, "error processing cycle time header")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"flowCycleTimeChart":{"took":205,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":837,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"flow_cycle_time_buckets":{"buckets":[{"key_as_string":"2024-02-26","key":1708885800000,"doc_count":37,"flow_cycle_time_count":{"value":{"TECH_DEBT":0,"DEFECT":1438132666,"FEATURE":29462000,"RISK":0}}},{"key_as_string":"2024-03-04","key":1709490600000,"doc_count":520,"flow_cycle_time_count":{"value":{"TECH_DEBT":0,"DEFECT":361104666,"FEATURE":1826674833,"RISK":0}}},{"key_as_string":"2024-03-11","key":1710095400000,"doc_count":280,"flow_cycle_time_count":{"value":{"TECH_DEBT":0,"DEFECT":1321270000,"FEATURE":1609742222,"RISK":0}}},{"key_as_string":"2024-03-18","key":1710700200000,"doc_count":0,"flow_cycle_time_count":{"value":{"TECH_DEBT":0,"DEFECT":0,"FEATURE":0,"RISK":0}}},{"key_as_string":"2024-03-25","key":1711305000000,"doc_count":0,"flow_cycle_time_count":{"value":{"TECH_DEBT":0,"DEFECT":0,"FEATURE":0,"RISK":0}}}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":1438132666},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":361104666},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":1321270000},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Defect","yAxisFormatter":{"type":"TIME_DURATION"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":29462000},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":1826674833},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":1609742222},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Feature","yAxisFormatter":{"type":"TIME_DURATION"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Risk","yAxisFormatter":{"type":"TIME_DURATION"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Tech debt","yAxisFormatter":{"type":"TIME_DURATION"}}]`)
		b, err = getCycleTime("flowCycleTimeChartSpec", x, replacements)
		assert.Nil(t, err, "error processing cycle time section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetWorkEfficiency(t *testing.T) {
	t.Run("Case 1: Successful execution of work wait time", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		//Header

		responseString := `{"flowEfficiencyHeader":{"took":212,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":844,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"flow_efficiency_count":{"value":"82%"}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":"82%"}`)
		b, err := getWorkEfficiency("flowEfficiencyHeaderSpec", x, replacements)

		assert.Nil(t, err, "error processing active work time header")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"flowEfficiencyChart":{"took":346,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":844,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"flow_eff_time_buckets":{"buckets":[{"key_as_string":"2024-02-26","key":1708885800000,"doc_count":37,"flow_efficiency_count":{"value":{"TECH_DEBT":0,"DEFECT":93,"FEATURE":100,"RISK":0}}},{"key_as_string":"2024-03-04","key":1709490600000,"doc_count":520,"flow_efficiency_count":{"value":{"TECH_DEBT":0,"DEFECT":71,"FEATURE":76,"RISK":0}}},{"key_as_string":"2024-03-11","key":1710095400000,"doc_count":287,"flow_efficiency_count":{"value":{"TECH_DEBT":0,"DEFECT":85,"FEATURE":94,"RISK":0}}},{"key_as_string":"2024-03-18","key":1710700200000,"doc_count":0,"flow_efficiency_count":{"value":{"TECH_DEBT":0,"DEFECT":0,"FEATURE":0,"RISK":0}}},{"key_as_string":"2024-03-25","key":1711305000000,"doc_count":0,"flow_efficiency_count":{"value":{"TECH_DEBT":0,"DEFECT":0,"FEATURE":0,"RISK":0}}}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":93},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":71},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":85},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Defect","yAxisFormatter":{"appendUnitValue":"%","type":"APPEND_UNIT"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":100},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":76},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":94},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Feature","yAxisFormatter":{"appendUnitValue":"%","type":"APPEND_UNIT"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Risk","yAxisFormatter":{"appendUnitValue":"%","type":"APPEND_UNIT"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Tech debt","yAxisFormatter":{"appendUnitValue":"%","type":"APPEND_UNIT"}}]`)
		b, err = getWorkEfficiency("flowEfficiencyChartSpec", x, replacements)
		assert.Nil(t, err, "error processing active work time section")
		assert.Equal(t, expectResult, []byte(b))

		//Header

		responseString = `{"flowWaitTimeHeader":{"took":37,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":844,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"flow_wait_time_count":{"value":"18%"}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`{"value":"18%"}`)
		b, err = getWorkEfficiency("flowWaitTimeHeaderSpec", x, replacements)

		assert.Nil(t, err, "error processing work wait time header")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"flowWaitTimeChart":{"took":33,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":844,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"flow_wait_time_count":{"value":[{"x":"BLOCKED","y":7},{"x":"CODE REVIEW","y":11}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"id":"work wait time","yAxisFormatter":{"type":"APPEND_UNIT","appendUnitValue":"%"},"data":[{"x":"BLOCKED","y":7},{"x":"CODE REVIEW","y":11}]}]`)
		b, err = getWorkEfficiency("flowWaitTimeChartSpec", x, replacements)
		assert.Nil(t, err, "error processing work wait time section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetVelocity(t *testing.T) {
	t.Run("Case 1: Successful execution of velocity", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"flowVelocityHeader":{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":879,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"velocity":{"value":71}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":71}`)
		b, err := getVelocity("flowVelocityHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing velocity header")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"flowVelocityChart":{"took":21,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":879,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"flow_velocity_buckets":{"buckets":[{"key_as_string":"2024-02-26","key":1708885800000,"doc_count":37,"flow_velocity_count":{"value":{"TECH_DEBT":0,"DEFECT":6,"FEATURE":1,"RISK":0}}},{"key_as_string":"2024-03-04","key":1709490600000,"doc_count":523,"flow_velocity_count":{"value":{"TECH_DEBT":0,"DEFECT":12,"FEATURE":24,"RISK":0}}},{"key_as_string":"2024-03-11","key":1710095400000,"doc_count":319,"flow_velocity_count":{"value":{"TECH_DEBT":0,"DEFECT":11,"FEATURE":17,"RISK":0}}},{"key_as_string":"2024-03-18","key":1710700200000,"doc_count":0,"flow_velocity_count":{"value":{"TECH_DEBT":0,"DEFECT":0,"FEATURE":0,"RISK":0}}},{"key_as_string":"2024-03-25","key":1711305000000,"doc_count":0,"flow_velocity_count":{"value":{"TECH_DEBT":0,"DEFECT":0,"FEATURE":0,"RISK":0}}}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":6},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":12},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":11},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Defect"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":1},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":24},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":17},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Feature"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Risk"},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Tech debt"}]`)
		b, err = getVelocity("flowVelocityChartSpec", x, replacements)
		assert.Nil(t, err, "error processing velocity section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetWorkItemDistribution(t *testing.T) {
	t.Run("Case 1: Successful execution of work item distribution", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"flowDistributionAvgChart":{"took":22,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":879,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"flow_distribution_avg_count":{"value":[{"title":"Defect","value":41},{"title":"Feature","value":59},{"title":"Risk","value":0},{"title":"Tech debt","value":0}]}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"title":"Defect","value":41},{"title":"Feature","value":59},{"title":"Risk","value":0},{"title":"Tech debt","value":0}]`)
		b, err := getWorkItemDistribution("flowDistributionChartAvgSpec", x, replacements)
		assert.Nil(t, err, "error processing work item distribution header")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"flowDistributionChart":{"took":22,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":879,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"flow_distribution_buckets":{"buckets":[{"key_as_string":"2024-02-26","key":1708885800000,"doc_count":37,"flow_distribution_count":{"value":{"TECH_DEBT":0,"DEFECT":86,"FEATURE":14,"RISK":0}}},{"key_as_string":"2024-03-04","key":1709490600000,"doc_count":523,"flow_distribution_count":{"value":{"TECH_DEBT":0,"DEFECT":33,"FEATURE":67,"RISK":0}}},{"key_as_string":"2024-03-11","key":1710095400000,"doc_count":319,"flow_distribution_count":{"value":{"TECH_DEBT":0,"DEFECT":39,"FEATURE":61,"RISK":0}}},{"key_as_string":"2024-03-18","key":1710700200000,"doc_count":0,"flow_distribution_count":{"value":{"TECH_DEBT":0,"DEFECT":0,"FEATURE":0,"RISK":0}}},{"key_as_string":"2024-03-25","key":1711305000000,"doc_count":0,"flow_distribution_count":{"value":{"TECH_DEBT":0,"DEFECT":0,"FEATURE":0,"RISK":0}}}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":86},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":33},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":39},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Defect","yAxisFormatter":{"appendUnitValue":"%","type":"APPEND_UNIT"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":14},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":67},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":61},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Feature","yAxisFormatter":{"appendUnitValue":"%","type":"APPEND_UNIT"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Risk","yAxisFormatter":{"appendUnitValue":"%","type":"APPEND_UNIT"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":0},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Tech debt","yAxisFormatter":{"appendUnitValue":"%","type":"APPEND_UNIT"}}]`)
		b, err = getWorkItemDistribution("flowDistributionChartSpec", x, replacements)

		assert.Nil(t, err, "error processing work item distribution section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetDeploymentFrequency(t *testing.T) {
	t.Run("Case 1: Successful execution of deployment frequency", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"deploymentFrequencyHeader":{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deploy_data":{"value":{"average":0.03,"deployments":1,"differenceDays":31}}}}}`
		// need to change test case

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":0.03}`)
		b, err := getDeploymentFrequency("deploymentFrequencyHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing deployment frequency header")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetDeploymentLeadTime(t *testing.T) {
	t.Run("Case 1: Successful execution of deployment lead time", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"deploymentLeadTimeHeader":{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deploy_data":{"value":{"totalDuration":57000,"average":57000,"deployments":1}}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"valueInMillis":57000}`)
		b, err := getDeploymentLeadTime("deploymentLeadTimeHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing deployment lead time header")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetFailureRate(t *testing.T) {
	t.Run("Case 1: Successful execution of failure rate", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"averageFailureRateHeader":{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deploy_data":{"value":{"average":"0.0%","deployments":1,"failedDeployments":0}}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":"0.0%"}`)
		b, err := getFailureRate("averageFailureRateHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing failure rate header")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetMttr(t *testing.T) {
	t.Run("Case 1: Successful execution of mttr", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"mttrHeader":{"took":3,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":32,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deployments":{"value":2379615.3846153845}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"valueInMillis":2379615}`)
		b, err := getMttr("mttrHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing mttr")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_totalBuildsHeader(t *testing.T) {
	replacements := map[string]any{
		"startDate":            "2024-03-01",
		"endDate":              "2024-03-31",
		"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		"aggrBy":               "week",
		"duration":             "month",
		"dateHistogramMin":     "2024-03-01",
		"dateHistogramMax":     "2024-03-31",
		"normalizeMonthInSpec": "2024-03-01",
	}
	t.Run("Case 1: Successful execution of TotalBuilds Header", func(t *testing.T) {
		responseString := `{"totalBuildsHeader":{"_shards":{"failed":0,"skipped":0,"successful":2,"total":2},"aggregations":{"total_builds":{"value":3932}},"hits":{"hits":[],"max_score":null,"total":{"relation":"gte","value":10000}},"status":200,"timed_out":false,"took":3}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":3932}`)
		b, err := totalBuildsHeader("", x, replacements)
		assert.Nil(t, err, "error processing totalBuildsHeader ")
		assert.Equal(t, expectResult, []byte(b))

	})
}

func Test_totalBuildsSection(t *testing.T) {

	replacements := map[string]any{
		"startDate":            "2024-03-01",
		"endDate":              "2024-03-31",
		"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		"aggrBy":               "week",
		"duration":             "month",
		"dateHistogramMin":     "2024-03-01",
		"dateHistogramMax":     "2024-03-31",
		"normalizeMonthInSpec": "2024-03-01",
	}

	responseStringSection := `{"buildsData":{"took":55,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"build_status":{"value":{"data":[{"name":"Success","value":98},{"name":"Failure","value":2}],"info":[{"drillDown":{"reportType":"status","reportId":"builds","reportTitle":"Builds"},"title":"Success","value":4189},{"drillDown":{"reportType":"status","reportId":"builds","reportTitle":"Builds"},"title":"Failure","value":77}]}}}}}`

	x := map[string]json.RawMessage{}
	json.Unmarshal([]byte(responseStringSection), &x)

	expectResult := []byte(`{"data":[{"name":"Success","value":98},{"name":"Failure","value":2}],"info":[{"drillDown":{"reportType":"status","reportId":"builds","reportTitle":"Builds"},"title":"Success","value":4189},{"drillDown":{"reportType":"status","reportId":"builds","reportTitle":"Builds"},"title":"Failure","value":77}]}`)
	b, err := totalBuildsSection("", x, replacements)
	assert.Nil(t, err, "error processing totalBuildsSection")
	assert.Equal(t, expectResult, []byte(b))
}

func Test_deploymentsHeader(t *testing.T) {
	replacements := map[string]any{
		"startDate":            "2024-03-01",
		"endDate":              "2024-03-31",
		"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		"aggrBy":               "week",
		"duration":             "month",
		"dateHistogramMin":     "2024-03-01",
		"dateHistogramMax":     "2024-03-31",
		"normalizeMonthInSpec": "2024-03-01",
	}
	t.Run("Case 1: Successful execution of deploymentsHeader", func(t *testing.T) {
		responseString := `{"deploymentsHeader":{"_shards":{"failed":0,"skipped":0,"successful":2,"total":2},"aggregations":{"deploy_count":{"value":1}},"hits":{"hits":[],"max_score":null,"total":{"relation":"eq","value":1}},"status":200,"timed_out":false,"took":3}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":1}`)
		b, err := deploymentsHeader("", x, replacements)
		assert.Nil(t, err, "error processing deploymentsHeader ")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_deploymentsSection(t *testing.T) {

	replacements := map[string]any{
		"startDate":            "2024-03-01",
		"endDate":              "2024-03-31",
		"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		"aggrBy":               "week",
		"duration":             "month",
		"dateHistogramMin":     "2024-03-01",
		"dateHistogramMax":     "2024-03-31",
		"normalizeMonthInSpec": "2024-03-01",
	}

	responseStringSection := `{"envDeploymentInfo":{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deploys":{"value":{"data":[{"x":"staging","y":1}]}}}}}`

	x := map[string]json.RawMessage{}
	json.Unmarshal([]byte(responseStringSection), &x)

	expectResult := []byte(`[{"data":[{"x":"staging","y":1}],"id":"Successful deployments"}]`)

	b, err := deploymentsSection("", x, replacements)
	assert.Nil(t, err, "error processing deploymentsSection")
	assert.Equal(t, expectResult, []byte(b))

}

func Test_GetDeploymentFrequencyAndLeadTime(t *testing.T) {
	t.Run("Case 1: Successful execution of deployment frequency and lead time", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"frequencyAndLeadTimeTrend":{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deploy_buckets":{"buckets":[{"key_as_string":"2024-02-26","key":1708885800000,"doc_count":0,"deployments":{"value":{"totalDuration":0,"average":0,"deployments":0}}},{"key_as_string":"2024-03-04","key":1709490600000,"doc_count":0,"deployments":{"value":{"totalDuration":0,"average":0,"deployments":0}}},{"key_as_string":"2024-03-11","key":1710095400000,"doc_count":1,"deployments":{"value":{"totalDuration":57000,"average":57000,"deployments":1}}},{"key_as_string":"2024-03-18","key":1710700200000,"doc_count":0,"deployments":{"value":{"totalDuration":0,"average":0,"deployments":0}}},{"key_as_string":"2024-03-25","key":1711305000000,"doc_count":0,"deployments":{"value":{"totalDuration":0,"average":0,"deployments":0}}}]}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"id":"Successful deployments","type":"bar","data":[{"x":"2024-03-01","y":0},{"x":"2024-03-04","y":0},{"x":"2024-03-11","y":1},{"x":"2024-03-18","y":0},{"x":"2024-03-25","y":0}],"yAxisFormatter":{"type":""}},{"id":"Deployment lead time","type":"line","data":[{"x":"2024-03-01","y":0},{"x":"2024-03-04","y":0},{"x":"2024-03-11","y":57000},{"x":"2024-03-18","y":0},{"x":"2024-03-25","y":0}],"yAxisFormatter":{"type":"TIME_DURATION"}}]`)
		b, err := getDeploymentFrequencyAndLeadTime("frequencyAndLeadTimeTrendSpec", x, replacements)
		assert.Nil(t, err, "error processing deploymentFrequencyAndLeadTime")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetFailureRateAndMttr(t *testing.T) {
	t.Run("Case 1: Successful execution of failure rate and mttr", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-02-01",
			"endDate":              "2024-02-29",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-02-01",
			"dateHistogramMax":     "2024-02-29",
			"normalizeMonthInSpec": "2024-02-01",
		}

		responseString := `{"failureRateAndMttrTrend":{"took":2,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":32,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deploy_buckets":{"buckets":[{"key_as_string":"2024-06-04","key":1717439400000,"doc_count":0,"deployments":{"value":{"failureRate":0,"total":0,"mttr":0,"failed":0}}},{"key_as_string":"2024-06-05","key":1717525800000,"doc_count":32,"deployments":{"value":{"failureRate":46.88,"total":32,"mttr":2379615.3846153845,"failed":15}}}]}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"id":"Failure rate","type":"bar","data":[{"x":"2024-02-01","y":0,"z":"0% (0 of 0 failed)"},{"x":"2024-06-05","y":46,"z":"46.88% (15 of 32 failed)"}],"yAxisFormatter":{"appendUnitValue":"%","type":"APPEND_UNIT"}},{"id":"Mean time to recovery","type":"line","data":[{"x":"2024-02-01","y":0,"z":""},{"x":"2024-06-05","y":2379615,"z":""}],"yAxisFormatter":{"appendUnitValue":"","type":"TIME_DURATION"}}]`)
		b, err := getFailureRateAndMttr("failureRateAndMttrTrendSpec", x, replacements)
		assert.Nil(t, err, "error processing failure rate and mttr")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_trivyLicenseOverviewSection(t *testing.T) {

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {

		responseString := `{
            "trivyLicensesOverviewSection": {
                "aggregations": {
                    "trivyLicenseSection": {
                        "value": [{
                            "drillDown": {
                                "reportId": "testID",
                                "reportTitle": "testTitle",
                                "reportInfo": {
                                    "component_id": "testComponentID",
                                    "run_id": "testRunID",
                                    "license_type": "testType",
                                    "branch": "testBranch"
                                }
                            },
                            "severity": "Low",
                            "licenseType": "testType",
                            "occurences": 1,
                            "classification": "Restricted",
                            "firstDiscovered": "2024/03/01 12:21:31"
                        }]
                    }
                }
            }
        }`

		x := map[string]json.RawMessage{}

		err := json.Unmarshal([]byte(responseString), &x)
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		reports, err := trivyLicenseOverviewSection("", x, nil)
		assert.Equal(t, 307, len(reports), "Validating response count for automation run")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 2: Error on invalid query key", func(t *testing.T) {

		responseString := `{
            "trivyLicensesOverviewSectionFail": {
                "aggregations": {
                    "trivyLicenseSection": {
                        "value": [{
                            "drillDown": {
                                "reportId": "testID",
                                "reportTitle": "testTitle",
                                "reportInfo": {
                                    "component_id": "testComponentID",
                                    "run_id": "testRunID",
                                    "license_type": "testType",
                                    "branch": "testBranch"
                                }
                            },
                            "severity": "Low",
                            "licenseType": "testType",
                            "occurences": 1,
                            "classification": "Restricted",
                            "firstDiscovered": "2024/03/01 12:21:31"
                        }]
                    }
                }
            }
        }`

		x := map[string]json.RawMessage{}

		json.Unmarshal([]byte(responseString), &x)

		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		_, err := mergedDefaultBranchCommitsSection("", x, nil)
		assert.Error(t, err, db.ErrInternalServer)

	})

}

func Test_GetVulnerabilitiesByScanType(t *testing.T) {
	t.Run("Case 1: Successful execution of Vulnerabilities by scan type", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		//Header
		responseString := `{"vulnerabilityByScannerTypeHeader":{"took":151,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"scanner_type_count":{"value":{"SCA":88,"DAST":0,"Container":1838,"SAST":92}}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":92}`)
		b, err := getVulnerabilitiesByScanType("SASTHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetVulnerabilitiesByScanType header - found")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":0}`)
		b, err = getVulnerabilitiesByScanType("DASTHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetVulnerabilitiesByScanType header - reopened")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":1838}`)
		b, err = getVulnerabilitiesByScanType("ContainerHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetVulnerabilitiesByScanType header - resolved")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":88}`)
		b, err = getVulnerabilitiesByScanType("SCAHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetVulnerabilitiesByScanType header - open")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"vulnerabilitybyscannertypechart":{"took":126,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"vulByScannerTypeCounts":{"value":{"VERY_HIGH":[{"x":"SAST","y":26},{"x":"DAST","y":0},{"x":"Container","y":26},{"x":"SCA","y":4}],"HIGH":[{"x":"SAST","y":27},{"x":"DAST","y":0},{"x":"Container","y":43},{"x":"SCA","y":44}],"MEDIUM":[{"x":"SAST","y":23},{"x":"DAST","y":0},{"x":"Container","y":126},{"x":"SCA","y":38}],"LOW":[{"x":"SAST","y":16},{"x":"DAST","y":0},{"x":"Container","y":1643},{"x":"SCA","y":2}]}}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"id":"Very high","data":[{"x":"SAST","y":26},{"x":"DAST","y":0},{"x":"Container","y":26},{"x":"SCA","y":4}]},{"id":"High","data":[{"x":"SAST","y":27},{"x":"DAST","y":0},{"x":"Container","y":43},{"x":"SCA","y":44}]},{"id":"Medium","data":[{"x":"SAST","y":23},{"x":"DAST","y":0},{"x":"Container","y":126},{"x":"SCA","y":38}]},{"id":"Low","data":[{"x":"SAST","y":16},{"x":"DAST","y":0},{"x":"Container","y":1643},{"x":"SCA","y":2}]}]`)
		b, err = getVulnerabilitiesByScanType("vulnerabilitybyscannertypechartSpec", x, replacements)
		assert.Nil(t, err, "error processing GetVulnerabilitiesByScanType section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetSlaStatusOverview(t *testing.T) {
	t.Run("Case 1: Successful execution of SLA status overview", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"slaStatusOverview":{"took":656,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"slaStatusOverview":{"value":{"openSlaCounts":[{"x":"Breached","y":12911},{"x":"At risk","y":528},{"x":"On track","y":4921}],"closedSlaCounts":[{"x":"Breached","y":20},{"x":"Within SLA","y":11}],"closeSlaKey":"Resolved vulnerabilites","openSlaKey":"Open vulnerabilities"}}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"id":"Open vulnerabilities","data":[{"x":"Breached","y":12911},{"x":"At risk","y":528},{"x":"On track","y":4921}]}]`)
		b, err := getSlaStatusOverview("slaStatusOverviewOpenSpec", x, replacements)
		assert.Nil(t, err, "error processing GetSlaStatusOverview section")
		assert.Equal(t, expectResult, []byte(b))

		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"id":"Resolved vulnerabilites","data":[{"x":"Breached","y":20},{"x":"Within SLA","y":11}]}]`)
		b, err = getSlaStatusOverview("slaStatusOverviewClosedSpec", x, replacements)
		assert.Nil(t, err, "error processing GetSlaStatusOverview section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetScanTypesInAutomation(t *testing.T) {
	t.Run("Case 1: Successful execution of scan types in automation", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"scanAutomationResp":{"took":1,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":141,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"scanTypesInAutomation":{"value":{"automationResult":[{"x":"SAST","y":10},{"x":"SCA","y":10},{"x":"Container","y":10}],"runKey":"Workflow Runs","automationKey":"Workflows","runResult":[{"x":"SAST","y":16}]}}}},"automationRunResp":{"took":6,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":301,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"automation_run_activity":{"value":{"95da8581-a320-4844-8e5b-dd5db5e462f5":[{"automation_id":"95da8581-a320-4844-8e5b-dd5db5e462f5","duration":0,"component_id":"62b6124f-4ba6-44d4-a83a-2dbd9aae97ce","run_id":"704e73f9-658b-4cf7-b5be-686203c68972","component_name":"asset-service","run_number":447,"status_timestamp":"2024-09-25T11:09:11.000Z","status":"Success"},{"automation_id":"95da8581-a320-4844-8e5b-dd5db5e462f5","duration":68000,"component_id":"62b6124f-4ba6-44d4-a83a-2dbd9aae97ce","run_id":"98493930-8f0b-4eaa-8087-2edb4642e233","component_name":"asset-service","run_number":445,"status_timestamp":"2024-09-25T09:27:01.000Z","status":"Success"}],"05b0578d-99b3-4c5d-a996-d2d67675f574":[{"automation_id":"05b0578d-99b3-4c5d-a996-d2d67675f574","duration":0,"component_id":"22ba7f3d-1944-456a-a4a4-b32ab6b0f83b","run_id":"d2d0e768-951e-418b-955f-e1712a3f3c29","component_name":"api-proto","run_number":824,"status_timestamp":"2024-09-25T10:50:40.000Z","status":"Success"}],"f8d603e2-6af8-4c69-9e05-ec8586444c2b":[{"automation_id":"f8d603e2-6af8-4c69-9e05-ec8586444c2b","duration":119000,"component_id":"22ba7f3d-1944-456a-a4a4-b32ab6b0f83b","run_id":"4fc7efdd-99ff-42ee-95af-145982c0a306","component_name":"api-proto","run_number":230,"status_timestamp":"2024-09-25T10:02:54.000Z","status":"Success"}]}}}},"scannerTypeResp":{"took":2,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":143,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"distinct_run":{"value":{"704e73f9-658b-4cf7-b5be-686203c68972":{"scanner_names":["snyksast","sonarqube"],"scanner_types":["SAST"]},"d2d0e768-951e-418b-955f-e1712a3f3c29":{"scanner_names":["sonarqube"],"scanner_types":["SCA","Container"]},"4fc7efdd-99ff-42ee-95af-145982c0a306":{"scanner_names":["sonarqube"],"scanner_types":["Container"]}}}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"id":"Workflows","data":[{"x":"SAST","y":10},{"x":"SCA","y":10},{"x":"Container","y":10}]},{"id":"Workflow Runs","data":[{"x":"SCA","y":1},{"x":"Container","y":2},{"x":"SAST","y":1}]}]`)
		b, err := GetScanTypesInAutomation("scanTypesInAutomationsSpec", x, replacements)

		assert.Nil(t, err, "error processing GetScanTypesInAutomation section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetMttrForVulnerabilities(t *testing.T) {
	t.Run("Case 1: Successful execution of MTTR for Vulnerabilities", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		//Header
		responseString := `{"MTTRHeaders":{"took":676,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"Avg_TTR":{"value":{"VERY_HIGH":"48d 6h","HIGH":"33d 15h","MEDIUM":"17d 2h","LOW":"46d 16h","HIGH_RESOLVED_COUNT":13,"LOW_RESOLVED_COUNT":7,"MEDIUM_RESOLVED_COUNT":2,"VERY_HIGH_RESOLVED_COUNT":9}}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"titleCount":9,"value":"48d 6h"}`)
		b, err := getMttrForVulnerabilities("MTTRHeaderSpecVeryHigh", x, replacements)
		assert.Nil(t, err, "error processing GetMttrForVulnerabilities header - found")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"titleCount":13,"value":"33d 15h"}`)
		b, err = getMttrForVulnerabilities("MTTRHeaderSpecHigh", x, replacements)
		assert.Nil(t, err, "error processing GetMttrForVulnerabilities header - reopened")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"titleCount":2,"value":"17d 2h"}`)
		b, err = getMttrForVulnerabilities("MTTRHeaderSpecMedium", x, replacements)
		assert.Nil(t, err, "error processing GetMttrForVulnerabilities header - resolved")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"titleCount":7,"value":"46d 16h"}`)
		b, err = getMttrForVulnerabilities("MTTRHeaderSpecLow", x, replacements)
		assert.Nil(t, err, "error processing GetMttrForVulnerabilities header - open")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"MTTRChart":{"took":774,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"MTTR_chart_date_buckets":{"buckets":[{"key_as_string":"2024-02-26","key":1708885800000,"doc_count":9421,"Avg_TTR":{"value":{"VERY_HIGH":1477444000,"HIGH":2145088000,"MEDIUM":1477444000,"LOW":5649160400}}},{"key_as_string":"2024-03-04","key":1709490600000,"doc_count":13816,"Avg_TTR":{"value":{"VERY_HIGH":50337000,"HIGH":82993000,"MEDIUM":0,"LOW":0}}},{"key_as_string":"2024-03-11","key":1710095400000,"doc_count":9704,"Avg_TTR":{"value":{"VERY_HIGH":4505523750,"HIGH":3387185875,"MEDIUM":0,"LOW":0}}},{"key_as_string":"2024-03-18","key":1710700200000,"doc_count":0,"Avg_TTR":{"value":{"VERY_HIGH":0,"HIGH":0,"MEDIUM":0,"LOW":0}}},{"key_as_string":"2024-03-25","key":1711305000000,"doc_count":0,"Avg_TTR":{"value":{"VERY_HIGH":0,"HIGH":0,"MEDIUM":0,"LOW":0}}}]}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":1477444000},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":50337000},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":4505523750},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Very high","yAxisFormatter":{"type":"TIME_DURATION"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":2145088000},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":82993000},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":3387185875},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"High","yAxisFormatter":{"type":"TIME_DURATION"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":1477444000},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Medium","yAxisFormatter":{"type":"TIME_DURATION"}},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":5649160400},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":0},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":0},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Low","yAxisFormatter":{"type":"TIME_DURATION"}}]`)
		b, err = getMttrForVulnerabilities("MTTRChartSpec", x, replacements)
		assert.Nil(t, err, "error processing GetMttrForVulnerabilities section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetCodeChurn(t *testing.T) {
	t.Run("Case 1: Successful execution of code churn", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"codeChurnChart":{"took":146,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":4171,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"code_churn_buckets":{"buckets":[{"key_as_string":"2024-02-26","key":1708885800000,"doc_count":353,"code_churn":{"value":{"lines_deleted":13365,"lines_added":24368}}},{"key_as_string":"2024-03-04","key":1709490600000,"doc_count":2026,"code_churn":{"value":{"lines_deleted":94034,"lines_added":175487}}},{"key_as_string":"2024-03-11","key":1710095400000,"doc_count":1792,"code_churn":{"value":{"lines_deleted":33624,"lines_added":59097}}},{"key_as_string":"2024-03-18","key":1710700200000,"doc_count":0,"code_churn":{"value":{}}},{"key_as_string":"2024-03-25","key":1711305000000,"doc_count":0,"code_churn":{"value":{}}}]}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":24368},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":175487},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":59097},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Additions","isClickDisable":true},{"data":[{"date":{"endDate":"2024-03-03","startDate":"2024-03-01"},"x":"2024-03-01","y":13365},{"date":{"endDate":"2024-03-10","startDate":"2024-03-04"},"x":"2024-03-04","y":94034},{"date":{"endDate":"2024-03-17","startDate":"2024-03-11"},"x":"2024-03-11","y":33624},{"date":{"endDate":"2024-03-24","startDate":"2024-03-18"},"x":"2024-03-18","y":0},{"date":{"endDate":"2024-03-31","startDate":"2024-03-25"},"x":"2024-03-25","y":0}],"id":"Deletions","isClickDisable":true}]`)
		b, err := getCodeChurn("codeChurnChartSpec", x, replacements)
		assert.Nil(t, err, "error processing GetCodeChurn section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetAverageDeploymentTime(t *testing.T) {
	t.Run("Case 1: Successful execution of average deployment time", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-12-16",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-12-16",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		//Header
		responseString := `{"averageDeploymentTime":{"took":5,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":5,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deploy_data":{"value":[{"average":64000,"count":1,"from":"PR","to":"staging","value":64000}]}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"fromTitle":"PR","toTitle":"staging","totalCount":1,"totalDuration":64000,"value":64000}]`)
		b, err := getAverageDeploymentTime("averageDeploymentTimeSpec", x, replacements)
		assert.Nil(t, err, "error processing GetAverageDeploymentTime header - found")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetCwetmTop25VulnerabilitiesHeader(t *testing.T) {
	t.Run("Case 1: Successful execution of GetCwetmTop25VulnerabilitiesHeader", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-12-16",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-12-16",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"top25Vul":{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":154,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"top25CWE":{"value":{"top25Table":{"CWE-79":{"name":"Cross-site Scripting (XSS)","issuesFound":2},"CWE-77":{"name":"Command_Injection","issuesFound":1},"CWE-798":{"name":"Use of Hardcoded Credentials","issuesFound":12},"CWE-89":{"name":"CWE-89: Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')","issuesFound":1},"CWE-362":{"name":"Race_Condition_In_Cross_Functionality","issuesFound":1}},"top25TotalCount":5}}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"subTitle":{"title":"CWE\u003csup\u003eTM\u003c/sup\u003e top 25 vulnerabilities"},"value":5,"subValue":25}`)

		b, err := getCwetmTop25Vulnerabilities("top25VulHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing GetCwetmTop25VulnerabilitiesHeader")
		assert.Equal(t, expectResult, []byte(b))
	})

	t.Run("Case 2: Successful execution of GetCwetmTop25VulnerabilitiesSection", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-12-16",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-12-16",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"top25Vul":{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":154,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"top25CWE":{"value":{"top25Table":{"CWE-79":{"name":"Cross-site Scripting (XSS)","issuesFound":2},"CWE-77":{"name":"Command_Injection","issuesFound":1},"CWE-798":{"name":"Use of Hardcoded Credentials","issuesFound":12},"CWE-89":{"name":"CWE-89: Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')","issuesFound":1},"CWE-362":{"name":"Race_Condition_In_Cross_Functionality","issuesFound":1}},"top25TotalCount":5}}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := `[{"id":"CWE-79","name":"Cross-site Scripting (XSS)","issuesFound":2},{"id":"CWE-89","name":"CWE-89: Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')","issuesFound":1},{"id":"CWE-77","name":"Command_Injection","issuesFound":1},{"id":"CWE-798","name":"Use of Hardcoded Credentials","issuesFound":12},{"id":"CWE-362","name":"Race_Condition_In_Cross_Functionality","issuesFound":1}]`

		b, err := getCwetmTop25Vulnerabilities("top25VulChartSpec", x, replacements)

		assert.Nil(t, err, "error processing GetCwetmTop25VulnerabilitiesSection")
		assert.Equal(t, expectResult, string(b))
	})

	t.Run("Case 3: Regression test in case the function starts returning vulnerabilities with 0 affected components", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-12-16",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-12-16",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
		}

		responseString := `{"top25Vul":{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":154,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"top25CWE":{"value":{"top25Table":{"CWE-79":{"name":"Cross-site Scripting (XSS)","issuesFound":0},"CWE-77":{"name":"Command_Injection","issuesFound":1},"CWE-798":{"name":"Use of Hardcoded Credentials","issuesFound":0},"CWE-89":{"name":"CWE-89: Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')","issuesFound":1},"CWE-362":{"name":"Race_Condition_In_Cross_Functionality","issuesFound":1}},"top25TotalCount":5}}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := `[{"id":"CWE-89","name":"CWE-89: Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')","issuesFound":1},{"id":"CWE-77","name":"Command_Injection","issuesFound":1},{"id":"CWE-362","name":"Race_Condition_In_Cross_Functionality","issuesFound":1}]`

		b, err := getCwetmTop25Vulnerabilities("top25VulChartSpec", x, replacements)

		assert.Nil(t, err, "error processing GetCwetmTop25VulnerabilitiesSection")
		assert.Equal(t, expectResult, string(b))
	})

}

func Test_GetDeploymentFrequencyComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of deployment frequency component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"deploymentFrequencyHeader":{"took":1,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":34,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deployment_frequency_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"020c6a27-5680-4d21-b1f2-3fac1a15053e","doc_count":18,"deploy_data":{"value":{"average":0.6,"deployments":18,"differenceDays":30}}},{"key":"e389272a-3ad0-4766-bed8-1eff2211ed70","doc_count":16,"deploy_data":{"value":{"average":0.53,"deployments":16,"differenceDays":30}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":1,"component_count":0,"total_value":1,"value_in_millis":0,"compare_reports":[{"is_sub_org":true,"sub_org_id":"sub-org-2","compare_title":"sub-org-2","sub_org_count":1,"component_count":0,"total_value":1,"value_in_millis":0,"compare_reports":[{"is_sub_org":true,"sub_org_id":"sub-org-3","compare_title":"sub-org-3","sub_org_count":0,"component_count":2,"total_value":1,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"020c6a27-5680-4d21-b1f2-3fac1a15053e","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"average per day","value":0.6},{"title":"deployments","value":18},{"title":"differenceDays","value":30}]}},{"is_sub_org":false,"sub_org_id":"e389272a-3ad0-4766-bed8-1eff2211ed70","compare_title":"component 2","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"average per day","value":0.53},{"title":"deployments","value":16},{"title":"differenceDays","value":30}]}}],"section":{"data":[{"title":"average per day","value":1.13},{"title":"deployments","value":34},{"title":"differenceDays","value":30}]}}],"section":{"data":[{"title":"average per day","value":1.13},{"title":"deployments","value":34},{"title":"differenceDays","value":30}]}}],"section":{"data":[{"title":"average per day","value":1.13},{"title":"deployments","value":34},{"title":"differenceDays","value":30}]}}]`)

		organisation, err := getMultiLevelSubOrg()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getDeploymentFrequencyComponentComparison("deploymentFrequencySpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing deployment frequency component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetOpenVulnerabilitiesOverviewComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of open vulnerabilities overview component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"openVulAgeChart":{"took":1573,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"open_vulnerabilities_overview_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":21022,"openVulSeverityCount":{"value":{"VERY_HIGH":3,"HIGH":4,"MEDIUM":4,"LOW":85}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":2554,"openVulSeverityCount":{"value":{"VERY_HIGH":9,"HIGH":42,"MEDIUM":18,"LOW":1}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":70,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"VERY_HIGH","value":9},{"title":"HIGH","value":42},{"title":"MEDIUM","value":18},{"title":"LOW","value":1}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":96,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":96,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"VERY_HIGH","value":3},{"title":"HIGH","value":4},{"title":"MEDIUM","value":4},{"title":"LOW","value":85}]}}],"section":{"data":[{"title":"VERY_HIGH","value":3},{"title":"HIGH","value":4},{"title":"MEDIUM","value":4},{"title":"LOW","value":85}]}}]`)

		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getOpenVulnerabilitiesOverviewComponentComparison("openVulnerabilitiesOverviewSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing open vulnerabilities overview component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetVulnerabilitiesOverviewComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of vulnerabilities overview component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"vulOverviewChart":{"took":458,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1750,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"vulnerabilities_overview_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":21022,"vulnerabilityStatusCounts":{"value":{"Reopened":0,"Resolved":0,"Found":42,"Open":96}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":2554,"vulnerabilityStatusCounts":{"value":{"Reopened":0,"Resolved":0,"Found":65,"Open":70}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":135,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Found","value":65},{"title":"Open","value":70},{"title":"Reopened","value":0},{"title":"Resolved","value":0}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":138,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":138,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Found","value":42},{"title":"Open","value":96},{"title":"Reopened","value":0},{"title":"Resolved","value":0}]}}],"section":{"data":[{"title":"Found","value":42},{"title":"Open","value":96},{"title":"Reopened","value":0},{"title":"Resolved","value":0}]}}]`)

		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getVulnerabilitiesOverviewComponentComparison("vulnerabilitiesOverviewSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing vulnerabilities overview component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetVelocityComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of velocity component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"flowVelocityChart":{"took":458,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1750,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"flow_velocity_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":207,"flow_velocity_count":{"value":{"TECH_DEBT":0,"DEFECT":24,"FEATURE":14,"RISK":0}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":143,"flow_velocity_count":{"value":{"TECH_DEBT":0,"DEFECT":4,"FEATURE":30,"RISK":0}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":38,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Bugs","value":24},{"title":"Feature","value":14},{"title":"Risk","value":0},{"title":"Tech Debt","value":0}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":34,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":34,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Bugs","value":4},{"title":"Feature","value":30},{"title":"Risk","value":0},{"title":"Tech Debt","value":0}]}}],"section":{"data":[{"title":"Bugs","value":4},{"title":"Feature","value":30},{"title":"Risk","value":0},{"title":"Tech Debt","value":0}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getVelocityComponentComparison("velocityComponentSpec", x, replacements, organisation)
		// bString, _ := json.Marshal(b)
		// fmt.Println(string(bString))
		assert.Nil(t, err, "error processing velocity component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetCycleTimeComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of cycle time component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"flowCycleTimeChart":{"took":11,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":126,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"cycle_time_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":16,"flow_cycle_time_count":{"value":{"FEATURE_TIME":0,"DEFECT_TIME":477737000,"DEFECT_COUNT":5,"TECH_DEBT_COUNT":0,"RISK_TIME":0,"TECH_DEBT_TIME":0,"FEATURE_COUNT":0,"RISK_COUNT":0}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":15,"flow_cycle_time_count":{"value":{"FEATURE_TIME":0,"DEFECT_TIME":770734000,"DEFECT_COUNT":2,"TECH_DEBT_COUNT":0,"RISK_TIME":0,"TECH_DEBT_TIME":0,"FEATURE_COUNT":0,"RISK_COUNT":0}}},{"key":"3891d690-a504-4889-4ab1-e469db088597","doc_count":12,"flow_cycle_time_count":{"value":{"FEATURE_TIME":0,"DEFECT_TIME":969051000,"DEFECT_COUNT":2,"TECH_DEBT_COUNT":0,"RISK_TIME":0,"TECH_DEBT_TIME":0,"FEATURE_COUNT":0,"RISK_COUNT":0}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":95547400,"compare_reports":null,"section":{"data":[{"title":"Bugs","value":100,"time":477737000,"count":5},{"title":"Feature","value":0,"time":0,"count":0},{"title":"Risk","value":0,"time":0,"count":0},{"title":"Tech Debt","value":0,"time":0,"count":0}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":385367000,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":385367000,"compare_reports":null,"section":{"data":[{"title":"Bugs","value":100,"time":770734000,"count":2},{"title":"Feature","value":0,"time":0,"count":0},{"title":"Risk","value":0,"time":0,"count":0},{"title":"Tech Debt","value":0,"time":0,"count":0}]}}],"section":{"data":[{"title":"Bugs","value":100,"time":770734000,"count":2},{"title":"Feature","value":0,"time":0,"count":0},{"title":"Risk","value":0,"time":0,"count":0},{"title":"Tech Debt","value":0,"time":0,"count":0}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getCycleTimeComponentComparison("cycleTimeComponentSpec", x, replacements, organisation)
		// bString, _ := json.Marshal(b)
		// fmt.Println(string(bString))
		assert.Nil(t, err, "error processing cycle time component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetActiveFlowTimeComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of active flow time component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"flowEfficiencyChart":{"took":46,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1750,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"active_work_time_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":207,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":5256294000,"flowTime":12151621000},"FEATURE":{"activeTime":5256294000,"flowTime":12151621000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":143,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":10285326000,"flowTime":10679960000},"FEATURE":{"activeTime":10285326000,"flowTime":10679960000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"cdcfa229-facc-4e9a-80a1-d82211f910b4","doc_count":109,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":4556335000,"flowTime":6768296000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"d2ab028a-4902-4207-4352-23e412c24884","doc_count":103,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":5865702000,"flowTime":6181037000},"FEATURE":{"activeTime":5865702000,"flowTime":6181037000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"3891d690-a504-4889-4ab1-e469db088597","doc_count":101,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":6409146000,"flowTime":6603123000},"FEATURE":{"activeTime":6409146000,"flowTime":6603123000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"ae694ee5-6234-419c-5035-f3d59622565d","doc_count":78,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":9932447000,"flowTime":10478205000},"FEATURE":{"activeTime":9932447000,"flowTime":10478205000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"9b707c3e-b9a2-4705-6c44-6f9f250f61c6","doc_count":77,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":7552787000,"flowTime":7666114000},"FEATURE":{"activeTime":7552787000,"flowTime":7666114000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"6cbb8b60-16d6-48a4-61fe-3a33c6c1169f","doc_count":69,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":6282756000,"flowTime":6596759000},"FEATURE":{"activeTime":6282756000,"flowTime":6596759000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"785894c9-91cf-47b9-734d-ae73377cc29c","doc_count":62,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":2874697000,"flowTime":4070111000},"FEATURE":{"activeTime":2874697000,"flowTime":4070111000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"65240b81-48d8-44a7-4592-82a9e7a1a132","doc_count":60,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":4736522000,"flowTime":6432075000},"FEATURE":{"activeTime":4736522000,"flowTime":6432075000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"55f26650-b259-4453-6ff3-c2843d096780","doc_count":59,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":22395544000,"flowTime":23194197000},"FEATURE":{"activeTime":22395544000,"flowTime":23194197000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"6e5a4968-1806-4ff4-5052-0cd9859c4d93","doc_count":50,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":7160520000,"flowTime":7312308000},"FEATURE":{"activeTime":7160520000,"flowTime":7312308000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"94a81dd1-3f52-4520-891e-b2440f660945","doc_count":49,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":5708209000,"flowTime":5747978000},"FEATURE":{"activeTime":5708209000,"flowTime":5747978000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"bdd99ab8-12b8-49c8-61b5-087796f009f3","doc_count":49,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":23476637000,"flowTime":31224504000},"FEATURE":{"activeTime":23476637000,"flowTime":31224504000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"363c6b99-4c67-4e02-9bda-10585eb25d72","doc_count":48,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":5708209000,"flowTime":5747978000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"7ee1c5e5-0564-43d9-7af7-ba46907f5dd0","doc_count":46,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":9207707000,"flowTime":9841050000},"FEATURE":{"activeTime":9207707000,"flowTime":9841050000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"1063a418-d663-4503-af48-fe5c6ff4eb68","doc_count":39,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":1518203000,"flowTime":2121393000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"dda69191-5492-4b7e-88b2-9d9d42f61899","doc_count":38,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":7520773000,"flowTime":7687596000},"FEATURE":{"activeTime":7520773000,"flowTime":7687596000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"238ffe68-8cb4-459d-64ac-2e4f752fe8dc","doc_count":37,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":1755840000,"flowTime":2954673000},"FEATURE":{"activeTime":1755840000,"flowTime":2954673000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"31e661ac-6553-4d6c-6f3c-6528f64b8fcd","doc_count":37,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":3287844000,"flowTime":5134119000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"746a9c25-df05-420e-5f4e-f65f687d777e","doc_count":28,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":6716639500,"flowTime":7240295500},"FEATURE":{"activeTime":6716639500,"flowTime":7240295500},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"e2f8fef6-5041-4843-b37e-6cdae38099bc","doc_count":22,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":7520773000,"flowTime":7687596000},"FEATURE":{"activeTime":7520773000,"flowTime":7687596000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"743ddab4-955d-4282-7230-eb2cf50886da","doc_count":18,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":-59762000,"flowTime":10555000},"FEATURE":{"activeTime":-59762000,"flowTime":10555000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"e0fb1caf-4c53-4be5-644e-1836524bab2a","doc_count":16,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":5708209000,"flowTime":5747978000},"FEATURE":{"activeTime":5708209000,"flowTime":5747978000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"e8f3a95e-2387-40d0-b727-5a542584e4b0","doc_count":16,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":399341000,"flowTime":749871000},"FEATURE":{"activeTime":0,"flowTimeTime":0},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"16b2c653-5f73-41d1-4765-5b8ae36fcfe2","doc_count":15,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":6504768000,"flowTime":6544537000},"FEATURE":{"activeTime":6504768000,"flowTime":6544537000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"2a315f19-88de-447c-8ad7-f3bbbd0a64c4","doc_count":13,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":13786000,"flowTime":263215000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"3a04548e-9af3-4c45-6aa2-937c13baaa40","doc_count":13,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":517658000,"flowTime":594249000},"FEATURE":{"activeTime":517658000,"flowTime":594249000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"9abaae9f-4e37-4c36-8424-53cfd9650a13","doc_count":13,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":5708209000,"flowTime":5747978000},"FEATURE":{"activeTime":5708209000,"flowTime":5747978000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"3c09cd3a-d954-4548-7d9a-f1807fb00daf","doc_count":12,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":12899329000,"flowTime":19910701000},"FEATURE":{"activeTime":12899329000,"flowTime":19910701000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"d5c3f808-d0a1-4e9d-4ca3-47d82aca3257","doc_count":11,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":5962515000,"flowTime":6089869000},"FEATURE":{"activeTime":5962515000,"flowTime":6089869000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"3192697f-a073-49d0-ae6f-a1159152ec76","doc_count":10,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":5879875000,"flowTime":6515610000},"FEATURE":{"activeTime":5879875000,"flowTime":6515610000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"aab42c4f-1d6a-476b-a0f8-d3d78671134a","doc_count":10,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":83487000,"flowTime":83824000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"054d5f8b-c1bf-4c2f-41d0-73e33dfcad8b","doc_count":9,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":7262000,"flowTime":99341000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"89e389cc-3f56-429c-6ee1-1a8bfc8ef5ae","doc_count":8,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":254306000,"flowTime":341891000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"986570f1-d595-4e17-bdd8-f141b181a798","doc_count":8,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":1247936000,"flowTime":1485319000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"ab21947d-d035-4794-b4f6-a90421fc2bf7","doc_count":8,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":941171000,"flowTime":941171000},"FEATURE":{"activeTime":941171000,"flowTime":941171000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"dfa63f63-d55d-4fe2-7c6e-dc429133593f","doc_count":8,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":6074863000,"flowTime":6181255000},"FEATURE":{"activeTime":6074863000,"flowTime":6181255000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"e69dab25-77dd-43df-9a4f-2c66fbc9b83a","doc_count":6,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":7151000,"flowTime":89175000},"FEATURE":{"activeTime":0,"flowTimeTime":0},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"f783334c-97b3-478e-8788-2d8e76c91d91","doc_count":6,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":266140000,"flowTime":514108000},"FEATURE":{"activeTime":0,"flowTimeTime":0},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"7208ab11-afa5-4145-57b3-e3d6df5b09c6","doc_count":5,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":1844578000,"flowTime":1918136000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"81d93aec-19a0-4258-59a5-55b7f7140947","doc_count":5,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":12866048000,"flowTime":19845280000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"92df4461-d43f-4d19-5e00-4b6346f747d9","doc_count":5,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":5708209000,"flowTime":5747978000},"FEATURE":{"activeTime":5708209000,"flowTime":5747978000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"20881a53-1dc9-4c49-98a2-31a08a3824e3","doc_count":4,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":560320000,"flowTime":625402000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"716d4922-add7-474b-7f21-240a8e253c38","doc_count":4,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":446969000,"flowTime":551628000},"FEATURE":{"activeTime":0,"flowTimeTime":0},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"c0fb4e57-38e7-4522-6b07-43353312800a","doc_count":4,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":5962515000,"flowTime":6089869000},"FEATURE":{"activeTime":5962515000,"flowTime":6089869000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"34df94d2-2892-4a53-7b0c-f778b6842bf0","doc_count":2,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":254306000,"flowTime":341891000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"3f9855a8-308e-455e-5601-7e967cb3f268","doc_count":2,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":171666000,"flowTime":767632000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"61e5dc5a-d59c-450c-99ee-009414a22d7a","doc_count":2,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":266140000,"flowTime":514108000},"FEATURE":{"activeTime":0,"flowTimeTime":0},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"05920b1f-02d6-47de-af20-d4dd209e8b17","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":266140000,"flowTime":514108000},"FEATURE":{"activeTime":0,"flowTimeTime":0},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"3cd7e6f3-8fcb-426f-791e-6285150c6994","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":266140000,"flowTime":514108000},"FEATURE":{"activeTime":0,"flowTimeTime":0},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"8f423498-2f04-4ac3-b948-bdeb1e824a09","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":366654000,"flowTime":433277000},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"a1d40803-7978-40ed-6c6e-d6fa114cbfa7","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":266140000,"flowTime":514108000},"FEATURE":{"activeTime":0,"flowTimeTime":0},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"d15bfc5b-1b01-4362-9851-f51388121e0f","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":266140000,"flowTime":514108000},"FEATURE":{"activeTime":0,"flowTimeTime":0},"RISK":{"activeTime":0,"flowTimeTime":0}}}},{"key":"e6c3a25b-e576-4d28-4ca5-d8c4af3d0113","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"activeTime":0,"flowTimeTime":0},"DEFECT":{"activeTime":0,"flowTimeTime":0},"FEATURE":{"activeTime":5708209000,"flowTime":5747978000},"RISK":{"activeTime":0,"flowTimeTime":0}}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":43,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Bugs","value":43,"active_time":5256294000,"flow_time":12151621000},{"title":"Feature","value":43,"active_time":5256294000,"flow_time":12151621000},{"title":"Risk","value":0,"active_time":0,"flow_time":0},{"title":"Tech Debt","value":0,"active_time":0,"flow_time":0}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":96,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":96,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Bugs","value":96,"active_time":10285326000,"flow_time":10679960000},{"title":"Feature","value":96,"active_time":10285326000,"flow_time":10679960000},{"title":"Risk","value":0,"active_time":0,"flow_time":0},{"title":"Tech Debt","value":0,"active_time":0,"flow_time":0}]}}],"section":{"data":[{"title":"Bugs","value":96,"active_time":10285326000,"flow_time":10679960000},{"title":"Feature","value":96,"active_time":10285326000,"flow_time":10679960000},{"title":"Risk","value":0,"active_time":0,"flow_time":0},{"title":"Tech Debt","value":0,"active_time":0,"flow_time":0}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getActiveFlowTimeComponentComparison("activeFlowTimeComponentSpec", x, replacements, organisation)
		// bString, _ := json.Marshal(b)
		// fmt.Println(string(bString))
		assert.Nil(t, err, "error processing active flow time component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetWorkWaitTimeComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of work wait time component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"flowWaitTimeChart":{"took":15,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1750,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"work_wait_time_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":207,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":6895327000,"flowTime":12151621000},"FEATURE":{"waitingTime":6895327000,"flowTime":12151621000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":143,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":394634000,"flowTime":10679960000},"FEATURE":{"waitingTime":394634000,"flowTime":10679960000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"cdcfa229-facc-4e9a-80a1-d82211f910b4","doc_count":109,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":2211961000,"flowTime":6768296000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"d2ab028a-4902-4207-4352-23e412c24884","doc_count":103,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":315335000,"flowTime":6181037000},"FEATURE":{"waitingTime":315335000,"flowTime":6181037000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"3891d690-a504-4889-4ab1-e469db088597","doc_count":101,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":193977000,"flowTime":6603123000},"FEATURE":{"waitingTime":193977000,"flowTime":6603123000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"ae694ee5-6234-419c-5035-f3d59622565d","doc_count":78,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":545758000,"flowTime":10478205000},"FEATURE":{"waitingTime":545758000,"flowTime":10478205000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"9b707c3e-b9a2-4705-6c44-6f9f250f61c6","doc_count":77,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":113327000,"flowTime":7666114000},"FEATURE":{"waitingTime":113327000,"flowTime":7666114000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"6cbb8b60-16d6-48a4-61fe-3a33c6c1169f","doc_count":69,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":314003000,"flowTime":6596759000},"FEATURE":{"waitingTime":314003000,"flowTime":6596759000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"785894c9-91cf-47b9-734d-ae73377cc29c","doc_count":62,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":1195414000,"flowTime":4070111000},"FEATURE":{"waitingTime":1195414000,"flowTime":4070111000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"65240b81-48d8-44a7-4592-82a9e7a1a132","doc_count":60,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":1695553000,"flowTime":6432075000},"FEATURE":{"waitingTime":1695553000,"flowTime":6432075000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"55f26650-b259-4453-6ff3-c2843d096780","doc_count":59,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":798653000,"flowTime":23194197000},"FEATURE":{"waitingTime":798653000,"flowTime":23194197000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"6e5a4968-1806-4ff4-5052-0cd9859c4d93","doc_count":50,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":151788000,"flowTime":7312308000},"FEATURE":{"waitingTime":151788000,"flowTime":7312308000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"94a81dd1-3f52-4520-891e-b2440f660945","doc_count":49,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":39769000,"flowTime":5747978000},"FEATURE":{"waitingTime":39769000,"flowTime":5747978000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"bdd99ab8-12b8-49c8-61b5-087796f009f3","doc_count":49,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":7747867000,"flowTime":31224504000},"FEATURE":{"waitingTime":7747867000,"flowTime":31224504000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"363c6b99-4c67-4e02-9bda-10585eb25d72","doc_count":48,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":39769000,"flowTime":5747978000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"7ee1c5e5-0564-43d9-7af7-ba46907f5dd0","doc_count":46,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":633343000,"flowTime":9841050000},"FEATURE":{"waitingTime":633343000,"flowTime":9841050000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"1063a418-d663-4503-af48-fe5c6ff4eb68","doc_count":39,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":603190000,"flowTime":2121393000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"dda69191-5492-4b7e-88b2-9d9d42f61899","doc_count":38,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":166823000,"flowTime":7687596000},"FEATURE":{"waitingTime":166823000,"flowTime":7687596000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"238ffe68-8cb4-459d-64ac-2e4f752fe8dc","doc_count":37,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":1198833000,"flowTime":2954673000},"FEATURE":{"waitingTime":1198833000,"flowTime":2954673000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"31e661ac-6553-4d6c-6f3c-6528f64b8fcd","doc_count":37,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":1846275000,"flowTime":5134119000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"746a9c25-df05-420e-5f4e-f65f687d777e","doc_count":28,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":523656000,"flowTime":7240295500},"FEATURE":{"waitingTime":523656000,"flowTime":7240295500},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"e2f8fef6-5041-4843-b37e-6cdae38099bc","doc_count":22,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":166823000,"flowTime":7687596000},"FEATURE":{"waitingTime":166823000,"flowTime":7687596000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"743ddab4-955d-4282-7230-eb2cf50886da","doc_count":18,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":70317000,"flowTime":10555000},"FEATURE":{"waitingTime":70317000,"flowTime":10555000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"e0fb1caf-4c53-4be5-644e-1836524bab2a","doc_count":16,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":39769000,"flowTime":5747978000},"FEATURE":{"waitingTime":39769000,"flowTime":5747978000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"e8f3a95e-2387-40d0-b727-5a542584e4b0","doc_count":16,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":350530000,"flowTime":749871000},"FEATURE":{"waitingTime":0,"flowTime":0},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"16b2c653-5f73-41d1-4765-5b8ae36fcfe2","doc_count":15,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":39769000,"flowTime":6544537000},"FEATURE":{"waitingTime":39769000,"flowTime":6544537000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"2a315f19-88de-447c-8ad7-f3bbbd0a64c4","doc_count":13,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":249429000,"flowTime":263215000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"3a04548e-9af3-4c45-6aa2-937c13baaa40","doc_count":13,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":76591000,"flowTime":594249000},"FEATURE":{"waitingTime":76591000,"flowTime":594249000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"9abaae9f-4e37-4c36-8424-53cfd9650a13","doc_count":13,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":39769000,"flowTime":5747978000},"FEATURE":{"waitingTime":39769000,"flowTime":5747978000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"3c09cd3a-d954-4548-7d9a-f1807fb00daf","doc_count":12,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":7011372000,"flowTime":19910701000},"FEATURE":{"waitingTime":7011372000,"flowTime":19910701000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"d5c3f808-d0a1-4e9d-4ca3-47d82aca3257","doc_count":11,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":127354000,"flowTime":6089869000},"FEATURE":{"waitingTime":127354000,"flowTime":6089869000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"3192697f-a073-49d0-ae6f-a1159152ec76","doc_count":10,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":635735000,"flowTime":6515610000},"FEATURE":{"waitingTime":635735000,"flowTime":6515610000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"aab42c4f-1d6a-476b-a0f8-d3d78671134a","doc_count":10,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":337000,"flowTime":83824000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"054d5f8b-c1bf-4c2f-41d0-73e33dfcad8b","doc_count":9,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":92079000,"flowTime":99341000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"89e389cc-3f56-429c-6ee1-1a8bfc8ef5ae","doc_count":8,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":87585000,"flowTime":341891000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"986570f1-d595-4e17-bdd8-f141b181a798","doc_count":8,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":237383000,"flowTime":1485319000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"ab21947d-d035-4794-b4f6-a90421fc2bf7","doc_count":8,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":941171000},"FEATURE":{"waitingTime":0,"flowTime":941171000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"dfa63f63-d55d-4fe2-7c6e-dc429133593f","doc_count":8,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":106392000,"flowTime":6181255000},"FEATURE":{"waitingTime":106392000,"flowTime":6181255000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"e69dab25-77dd-43df-9a4f-2c66fbc9b83a","doc_count":6,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":82024000,"flowTime":89175000},"FEATURE":{"waitingTime":0,"flowTime":0},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"f783334c-97b3-478e-8788-2d8e76c91d91","doc_count":6,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":247968000,"flowTime":514108000},"FEATURE":{"waitingTime":0,"flowTime":0},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"7208ab11-afa5-4145-57b3-e3d6df5b09c6","doc_count":5,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":73558000,"flowTime":1918136000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"81d93aec-19a0-4258-59a5-55b7f7140947","doc_count":5,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":6979232000,"flowTime":19845280000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"92df4461-d43f-4d19-5e00-4b6346f747d9","doc_count":5,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":39769000,"flowTime":5747978000},"FEATURE":{"waitingTime":39769000,"flowTime":5747978000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"20881a53-1dc9-4c49-98a2-31a08a3824e3","doc_count":4,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":65082000,"flowTime":625402000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"716d4922-add7-474b-7f21-240a8e253c38","doc_count":4,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":104659000,"flowTime":551628000},"FEATURE":{"waitingTime":0,"flowTime":0},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"c0fb4e57-38e7-4522-6b07-43353312800a","doc_count":4,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":127354000,"flowTime":6089869000},"FEATURE":{"waitingTime":127354000,"flowTime":6089869000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"34df94d2-2892-4a53-7b0c-f778b6842bf0","doc_count":2,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":87585000,"flowTime":341891000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"3f9855a8-308e-455e-5601-7e967cb3f268","doc_count":2,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":595966000,"flowTime":767632000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"61e5dc5a-d59c-450c-99ee-009414a22d7a","doc_count":2,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":247968000,"flowTime":514108000},"FEATURE":{"waitingTime":0,"flowTime":0},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"05920b1f-02d6-47de-af20-d4dd209e8b17","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":247968000,"flowTime":514108000},"FEATURE":{"waitingTime":0,"flowTime":0},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"3cd7e6f3-8fcb-426f-791e-6285150c6994","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":247968000,"flowTime":514108000},"FEATURE":{"waitingTime":0,"flowTime":0},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"8f423498-2f04-4ac3-b948-bdeb1e824a09","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":66623000,"flowTime":433277000},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"a1d40803-7978-40ed-6c6e-d6fa114cbfa7","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":247968000,"flowTime":514108000},"FEATURE":{"waitingTime":0,"flowTime":0},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"d15bfc5b-1b01-4362-9851-f51388121e0f","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":247968000,"flowTime":514108000},"FEATURE":{"waitingTime":0,"flowTime":0},"RISK":{"waitingTime":0,"flowTime":0}}}},{"key":"e6c3a25b-e576-4d28-4ca5-d8c4af3d0113","doc_count":1,"flow_efficiency_count":{"value":{"TECH_DEBT":{"waitingTime":0,"flowTime":0},"DEFECT":{"waitingTime":0,"flowTime":0},"FEATURE":{"waitingTime":39769000,"flowTime":5747978000},"RISK":{"waitingTime":0,"flowTime":0}}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":57,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Bugs","value":57,"active_time":6895327000,"flow_time":6895327000},{"title":"Feature","value":57,"active_time":6895327000,"flow_time":6895327000},{"title":"Risk","value":0,"active_time":0,"flow_time":0},{"title":"Tech Debt","value":0,"active_time":0,"flow_time":0}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":4,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":4,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Bugs","value":4,"active_time":394634000,"flow_time":394634000},{"title":"Feature","value":4,"active_time":394634000,"flow_time":394634000},{"title":"Risk","value":0,"active_time":0,"flow_time":0},{"title":"Tech Debt","value":0,"active_time":0,"flow_time":0}]}}],"section":{"data":[{"title":"Bugs","value":4,"active_time":394634000,"flow_time":10679960000},{"title":"Feature","value":4,"active_time":394634000,"flow_time":10679960000},{"title":"Risk","value":0,"active_time":0,"flow_time":0},{"title":"Tech Debt","value":0,"active_time":0,"flow_time":0}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getWorkWaitTimeComponentComparison("workWaitTimeComponentSpec", x, replacements, organisation)
		// bString, _ := json.Marshal(b)
		// fmt.Println(string(bString))
		assert.Nil(t, err, "error processing active work wait time component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetWorkloadComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of workload component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"flowWorkLoad":{"took":1676,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"workload_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":4441,"work_load_counts":{"value":{"headerValue":46,"dates":{"2024-03-01":{"DEFECT":25,"FEATURE":21,"RISK":0,"TECH_DEBT":0,"DEFECT_SET":["SDP-11950","SDP-15413","SDP-14423","SDP-14981","SDP-14727","SDP-14627","SDP-14625","SDP-14108","SDP-15122","SDP-14954","SDP-14533","SDP-14896","SDP-13784","SDP-14475","SDP-15069","SDP-14452","SDP-14056","SDP-15068","SDP-13704","SDP-14979","SDP-13844","SDP-15130","SDP-15171","SDP-15193","SDP-13191"],"FEATURE_SET":["SDP-14943","SDP-14646","SDP-12236","SDP-14876","SDP-13414","SDP-15129","SDP-11970","SDP-10025","SDP-12235","SDP-12332","SDP-14982","SDP-12233","SDP-15235","SDP-13660","SDP-12706","SDP-14726","SDP-12437","SDP-9418","SDP-14274","SDP-14373","SDP-14372"],"RISK_SET":[],"TECH_DEBT_SET":[]}}}}},{"key":"6cbb8b60-16d6-48a4-61fe-3a33c6c1169f","doc_count":716,"work_load_counts":{"value":{"headerValue":17,"dates":{"2024-03-01":{"DEFECT":8,"FEATURE":9,"RISK":0,"TECH_DEBT":0,"DEFECT_SET":["SDP-14822","SDP-14612","SDP-14875","SDP-14610","SDP-14663","SDP-14948","SDP-15408","SDP-14703"],"FEATURE_SET":["SDP-14874","SDP-14950","SDP-13355","SDP-13575","SDP-12706","SDP-12705","SDP-13208","SDP-13050","SDP-13209"],"RISK_SET":[],"TECH_DEBT_SET":[]}}}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":662,"work_load_counts":{"value":{"headerValue":43,"dates":{"2024-03-01":{"DEFECT":6,"FEATURE":37,"RISK":0,"TECH_DEBT":0,"DEFECT_SET":["SDP-14663","SDP-15059","SDP-15145","SDP-14749","SDP-14703","SDP-15060"],"FEATURE_SET":["SDP-14800","SDP-14568","SDP-14821","SDP-14963","SDP-14287","SDP-15057","SDP-14808","SDP-12706","SDP-14807","SDP-14806","SDP-14805","SDP-12705","SDP-14804","SDP-14803","SDP-11996","SDP-14809","SDP-14274","SDP-14811","SDP-14810","SDP-15327","SDP-14798","SDP-14258","SDP-15049","SDP-15324","SDP-14311","SDP-15024","SDP-14819","SDP-14817","SDP-14816","SDP-14815","SDP-14859","SDP-14814","SDP-14813","SDP-14559","SDP-14812","SDP-14561","SDP-15050"],"RISK_SET":[],"TECH_DEBT_SET":[]}}}}},{"key":"6e5a4968-1806-4ff4-5052-0cd9859c4d93","doc_count":504,"work_load_counts":{"value":{"headerValue":15,"dates":{"2024-03-01":{"DEFECT":10,"FEATURE":5,"RISK":0,"TECH_DEBT":0,"DEFECT_SET":["SDP-14448","SDP-14687","SDP-14663","SDP-15432","SDP-15166","SDP-14693","SDP-15187","SDP-11847","SDP-14703","SDP-14180"],"FEATURE_SET":["SDP-14722","SDP-14121","SDP-12706","SDP-12705","SDP-8865"],"RISK_SET":[],"TECH_DEBT_SET":[]}}}}},{"key":"d2ab028a-4902-4207-4352-23e412c24884","doc_count":455,"work_load_counts":{"value":{"headerValue":5,"dates":{"2024-03-01":{"DEFECT":2,"FEATURE":3,"RISK":0,"TECH_DEBT":0,"DEFECT_SET":["SDP-14663","SDP-14703"],"FEATURE_SET":["SDP-14206","SDP-12706","SDP-12705"],"RISK_SET":[],"TECH_DEBT_SET":[]}}}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":46,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Bugs","value":25},{"title":"Feature","value":21},{"title":"Risk","value":0},{"title":"Tech Debt","value":0}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":43,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":43,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Bugs","value":6},{"title":"Feature","value":37},{"title":"Risk","value":0},{"title":"Tech Debt","value":0}]}}],"section":{"data":[{"title":"Bugs","value":6},{"title":"Feature","value":37},{"title":"Risk","value":0},{"title":"Tech Debt","value":0}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getWorkloadComponentComparison("workloadComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing workload component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func getMultipleOrganisation() (*constants.Organization, error) {
	jsonData := `
	{
	    "id": "Org-1",
	    "name": "Org-1",
	    "sub_orgs":
	    [
	        {
	            "id": "sub-org-1",
	            "name": "sub-org-1",
	            "sub_orgs":
	            [],
	            "components":
	            [
	                {
	                    "id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
	                    "name": "component 11"
	                },
					{
	                    "id": "cb4ffec3-1a8d-4c8a-9fc1-b4d9639b268c",
	                    "name": "component 12"
	                }
	            ]
	        }
	    ],
	    "components":
	    [
	        {
	            "id": "f4180826-bb76-421e-5410-791408daadeb",
	            "name": "component 1"
	        }
	    ]
	}`

	// Unmarshal JSON into Organization struct
	var orgData constants.Organization
	if err := json.Unmarshal([]byte(jsonData), &orgData); err != nil {
		log.Errorf(err, "Error parsing JSON : ")
		return nil, err
	}
	return &orgData, nil
}

func getOrganisation() (*constants.Organization, error) {
	jsonData := `
	{
	    "id": "Org-1",
	    "name": "Org-1",
	    "sub_orgs":
	    [
	        {
	            "id": "sub-org-1",
	            "name": "sub-org-1",
	            "sub_orgs":
	            [],
	            "components":
	            [
	                {
	                    "id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
	                    "name": "component 11"
	                }
	            ]
	        }
	    ],
	    "components":
	    [
	        {
	            "id": "f4180826-bb76-421e-5410-791408daadeb",
	            "name": "component 1"
	        }
	    ]
	}`

	// Unmarshal JSON into Organization struct
	var orgData constants.Organization
	if err := json.Unmarshal([]byte(jsonData), &orgData); err != nil {
		log.Errorf(err, "Error parsing JSON : ")
		return nil, err
	}
	return &orgData, nil
}

func getMultiLevelSubOrg() (*constants.Organization, error) {
	jsonData := `
	{
	    "id": "Org-1",
	    "name": "Org-1",
	    "sub_orgs":
	    [
	        {
	            "id": "sub-org-1",
	            "name": "sub-org-1",
	            "sub_orgs":
	            [{
					"id": "sub-org-2",
					"name": "sub-org-2",
					"sub_orgs":
					[{
						"id": "sub-org-3",
						"name": "sub-org-3",
						"sub_orgs":
						[],
						"components":
						[
							{
								"id": "020c6a27-5680-4d21-b1f2-3fac1a15053e",
								"name": "component 1"
							},
							{
								"id": "e389272a-3ad0-4766-bed8-1eff2211ed70",
								"name": "component 2"
							},
							{
								"id": "e389272a-3ad0-4766-bed8-1eff2211ed71",
								"name": "component 3"
							}
						]
					}],
					"components":
					[
					]
				}],
	            "components":
	            [
	            ]
	        }
	    ],
	    "components":
	    [{
			"id": "e389272a-3ad0-4766-bed8-1eff2211ed71",
			"name": "component 2"
		}
	    ]
	}`

	// Unmarshal JSON into Organization struct
	var orgData constants.Organization
	if err := json.Unmarshal([]byte(jsonData), &orgData); err != nil {
		log.Errorf(err, "Error parsing JSON : ")
		return nil, err
	}
	return &orgData, nil
}

func getOrganisationInactive() (*constants.Organization, error) {
	jsonData := `
	{
	    "id": "Org-1",
	    "name": "Org-1",
	    "sub_orgs":
	    [
	        {
	            "id": "sub-org-1",
	            "name": "sub-org-1",
	            "sub_orgs":
	            [],
	            "components":
	            [
	                {
	                    "id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
	                    "name": "component 11"
	                },
					{
						"id": "test-inactive-sub-org-comp",
						"name": "test inactive sub org component"
					}
	            ]
	        }
	    ],
	    "components":
	    [
	        {
	            "id": "f4180826-bb76-421e-5410-791408daadeb",
	            "name": "component 1"
	        },
			{
	            "id": "test-inactive-comp",
	            "name": "test inactive component"
	        }
	    ]
	}`

	// Unmarshal JSON into Organization struct
	var orgData constants.Organization
	if err := json.Unmarshal([]byte(jsonData), &orgData); err != nil {
		log.Errorf(err, "Error parsing JSON : ")
		return nil, err
	}
	return &orgData, nil
}

func Test_GetCommitTrendsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of commit trends component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"commitsAndAverageChart":{"took":18,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"commits_trends_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"162239c9-e4dd-4c7c-9c79-a97841e625e3","doc_count":13363,"unique_authors":{"value":1},"commits_count":{"value":13363},"commits-per-auth":{"value":13363}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":1390,"unique_authors":{"value":35},"commits_count":{"value":1390},"commits-per-auth":{"value":39.714285714285715}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":331,"unique_authors":{"value":18},"commits_count":{"value":331},"commits-per-auth":{"value":18.38888888888889}},{"key":"545c6a93-7b98-49b4-707e-9c9ab52ce617","doc_count":280,"unique_authors":{"value":8},"commits_count":{"value":280},"commits-per-auth":{"value":35}},{"key":"94a81dd1-3f52-4520-891e-b2440f660945","doc_count":247,"unique_authors":{"value":6},"commits_count":{"value":247},"commits-per-auth":{"value":41.166666666666664}},{"key":"3891d690-a504-4889-4ab1-e469db088597","doc_count":179,"unique_authors":{"value":15},"commits_count":{"value":179},"commits-per-auth":{"value":11.933333333333334}},{"key":"6cbb8b60-16d6-48a4-61fe-3a33c6c1169f","doc_count":158,"unique_authors":{"value":12},"commits_count":{"value":158},"commits-per-auth":{"value":13.166666666666666}},{"key":"6e5a4968-1806-4ff4-5052-0cd9859c4d93","doc_count":150,"unique_authors":{"value":8},"commits_count":{"value":150},"commits-per-auth":{"value":18.75}},{"key":"55f26650-b259-4453-6ff3-c2843d096780","doc_count":132,"unique_authors":{"value":31},"commits_count":{"value":132},"commits-per-auth":{"value":4.258064516129032}},{"key":"363c6b99-4c67-4e02-9bda-10585eb25d72","doc_count":130,"unique_authors":{"value":2},"commits_count":{"value":130},"commits-per-auth":{"value":65}},{"key":"238ffe68-8cb4-459d-64ac-2e4f752fe8dc","doc_count":123,"unique_authors":{"value":7},"commits_count":{"value":123},"commits-per-auth":{"value":17.571428571428573}},{"key":"e0fb1caf-4c53-4be5-644e-1836524bab2a","doc_count":122,"unique_authors":{"value":13},"commits_count":{"value":122},"commits-per-auth":{"value":9.384615384615385}},{"key":"785894c9-91cf-47b9-734d-ae73377cc29c","doc_count":116,"unique_authors":{"value":6},"commits_count":{"value":116},"commits-per-auth":{"value":19.333333333333332}},{"key":"bdd99ab8-12b8-49c8-61b5-087796f009f3","doc_count":112,"unique_authors":{"value":8},"commits_count":{"value":112},"commits-per-auth":{"value":14}},{"key":"9b707c3e-b9a2-4705-6c44-6f9f250f61c6","doc_count":107,"unique_authors":{"value":8},"commits_count":{"value":107},"commits-per-auth":{"value":13.375}},{"key":"cdcfa229-facc-4e9a-80a1-d82211f910b4","doc_count":107,"unique_authors":{"value":3},"commits_count":{"value":107},"commits-per-auth":{"value":35.666666666666664}},{"key":"ae694ee5-6234-419c-5035-f3d59622565d","doc_count":105,"unique_authors":{"value":20},"commits_count":{"value":105},"commits-per-auth":{"value":5.25}},{"key":"65240b81-48d8-44a7-4592-82a9e7a1a132","doc_count":97,"unique_authors":{"value":12},"commits_count":{"value":97},"commits-per-auth":{"value":8.083333333333334}},{"key":"839a6ffc-3c14-41ca-8621-532e768ea9d5","doc_count":97,"unique_authors":{"value":2},"commits_count":{"value":97},"commits-per-auth":{"value":48.5}},{"key":"7ee1c5e5-0564-43d9-7af7-ba46907f5dd0","doc_count":86,"unique_authors":{"value":6},"commits_count":{"value":86},"commits-per-auth":{"value":14.333333333333334}},{"key":"1063a418-d663-4503-af48-fe5c6ff4eb68","doc_count":73,"unique_authors":{"value":6},"commits_count":{"value":73},"commits-per-auth":{"value":12.166666666666666}},{"key":"746a9c25-df05-420e-5f4e-f65f687d777e","doc_count":68,"unique_authors":{"value":15},"commits_count":{"value":68},"commits-per-auth":{"value":4.533333333333333}},{"key":"16b2c653-5f73-41d1-4765-5b8ae36fcfe2","doc_count":63,"unique_authors":{"value":6},"commits_count":{"value":63},"commits-per-auth":{"value":10.5}},{"key":"743ddab4-955d-4282-7230-eb2cf50886da","doc_count":59,"unique_authors":{"value":5},"commits_count":{"value":59},"commits-per-auth":{"value":11.8}},{"key":"dda69191-5492-4b7e-88b2-9d9d42f61899","doc_count":57,"unique_authors":{"value":7},"commits_count":{"value":57},"commits-per-auth":{"value":8.142857142857142}},{"key":"e2f8fef6-5041-4843-b37e-6cdae38099bc","doc_count":56,"unique_authors":{"value":10},"commits_count":{"value":56},"commits-per-auth":{"value":5.6}},{"key":"05920b1f-02d6-47de-af20-d4dd209e8b17","doc_count":53,"unique_authors":{"value":3},"commits_count":{"value":53},"commits-per-auth":{"value":17.666666666666668}},{"key":"d2ab028a-4902-4207-4352-23e412c24884","doc_count":52,"unique_authors":{"value":4},"commits_count":{"value":52},"commits-per-auth":{"value":13}},{"key":"cf208105-08b5-44ad-842c-5fc1dfae3c76","doc_count":50,"unique_authors":{"value":3},"commits_count":{"value":50},"commits-per-auth":{"value":16.666666666666668}},{"key":"590947f7-4640-422b-a8a9-36aae8b85ee6","doc_count":48,"unique_authors":{"value":4},"commits_count":{"value":48},"commits-per-auth":{"value":12}},{"key":"9af098ae-d223-4a2f-a185-e6557db9d8a1","doc_count":44,"unique_authors":{"value":3},"commits_count":{"value":44},"commits-per-auth":{"value":14.666666666666666}},{"key":"dfa63f63-d55d-4fe2-7c6e-dc429133593f","doc_count":37,"unique_authors":{"value":4},"commits_count":{"value":37},"commits-per-auth":{"value":9.25}},{"key":"06383cef-97b7-47b2-b638-d51184f55b6f","doc_count":34,"unique_authors":{"value":2},"commits_count":{"value":34},"commits-per-auth":{"value":17}},{"key":"d5c3f808-d0a1-4e9d-4ca3-47d82aca3257","doc_count":33,"unique_authors":{"value":8},"commits_count":{"value":33},"commits-per-auth":{"value":4.125}},{"key":"31e661ac-6553-4d6c-6f3c-6528f64b8fcd","doc_count":32,"unique_authors":{"value":7},"commits_count":{"value":32},"commits-per-auth":{"value":4.571428571428571}},{"key":"ab21947d-d035-4794-b4f6-a90421fc2bf7","doc_count":32,"unique_authors":{"value":4},"commits_count":{"value":32},"commits-per-auth":{"value":8}},{"key":"c4398680-af38-4e76-bab8-0f333042a891","doc_count":32,"unique_authors":{"value":1},"commits_count":{"value":32},"commits-per-auth":{"value":32}},{"key":"3c09cd3a-d954-4548-7d9a-f1807fb00daf","doc_count":31,"unique_authors":{"value":4},"commits_count":{"value":31},"commits-per-auth":{"value":7.75}},{"key":"d000292d-bd9c-44c9-a3f5-598ec3a561f4","doc_count":31,"unique_authors":{"value":1},"commits_count":{"value":31},"commits-per-auth":{"value":31}},{"key":"3a04548e-9af3-4c45-6aa2-937c13baaa40","doc_count":30,"unique_authors":{"value":6},"commits_count":{"value":30},"commits-per-auth":{"value":5}},{"key":"61751a40-2830-4783-8528-6ce7fad549b9","doc_count":24,"unique_authors":{"value":1},"commits_count":{"value":24},"commits-per-auth":{"value":24}},{"key":"69b3cf79-c1bd-47a6-934c-d2a83c7741e9","doc_count":24,"unique_authors":{"value":3},"commits_count":{"value":24},"commits-per-auth":{"value":8}},{"key":"81d93aec-19a0-4258-59a5-55b7f7140947","doc_count":24,"unique_authors":{"value":3},"commits_count":{"value":24},"commits-per-auth":{"value":8}},{"key":"9abaae9f-4e37-4c36-8424-53cfd9650a13","doc_count":23,"unique_authors":{"value":5},"commits_count":{"value":23},"commits-per-auth":{"value":4.6}},{"key":"f783334c-97b3-478e-8788-2d8e76c91d91","doc_count":23,"unique_authors":{"value":5},"commits_count":{"value":23},"commits-per-auth":{"value":4.6}},{"key":"896a0752-03bf-4531-946f-ab2c88b14a02","doc_count":21,"unique_authors":{"value":3},"commits_count":{"value":21},"commits-per-auth":{"value":7}},{"key":"054d5f8b-c1bf-4c2f-41d0-73e33dfcad8b","doc_count":20,"unique_authors":{"value":3},"commits_count":{"value":20},"commits-per-auth":{"value":6.666666666666667}},{"key":"e8f3a95e-2387-40d0-b727-5a542584e4b0","doc_count":19,"unique_authors":{"value":3},"commits_count":{"value":19},"commits-per-auth":{"value":6.333333333333333}},{"key":"92df4461-d43f-4d19-5e00-4b6346f747d9","doc_count":18,"unique_authors":{"value":2},"commits_count":{"value":18},"commits-per-auth":{"value":9}},{"key":"a677dedb-dd24-42b5-8f55-16594562f0de","doc_count":18,"unique_authors":{"value":2},"commits_count":{"value":18},"commits-per-auth":{"value":9}},{"key":"a1d40803-7978-40ed-6c6e-d6fa114cbfa7","doc_count":17,"unique_authors":{"value":2},"commits_count":{"value":17},"commits-per-auth":{"value":8.5}},{"key":"c0fb4e57-38e7-4522-6b07-43353312800a","doc_count":16,"unique_authors":{"value":3},"commits_count":{"value":16},"commits-per-auth":{"value":5.333333333333333}},{"key":"89e389cc-3f56-429c-6ee1-1a8bfc8ef5ae","doc_count":15,"unique_authors":{"value":3},"commits_count":{"value":15},"commits-per-auth":{"value":5}},{"key":"3192697f-a073-49d0-ae6f-a1159152ec76","doc_count":14,"unique_authors":{"value":4},"commits_count":{"value":14},"commits-per-auth":{"value":3.5}},{"key":"3cd7e6f3-8fcb-426f-791e-6285150c6994","doc_count":14,"unique_authors":{"value":1},"commits_count":{"value":14},"commits-per-auth":{"value":14}},{"key":"ac44a88d-a3a5-44d5-90a3-01e3fdb0d535","doc_count":14,"unique_authors":{"value":6},"commits_count":{"value":14},"commits-per-auth":{"value":2.3333333333333335}},{"key":"d1e786b7-c2f0-48f4-9055-741e7ad2b171","doc_count":13,"unique_authors":{"value":2},"commits_count":{"value":13},"commits-per-auth":{"value":6.5}},{"key":"d15bfc5b-1b01-4362-9851-f51388121e0f","doc_count":12,"unique_authors":{"value":2},"commits_count":{"value":12},"commits-per-auth":{"value":6}},{"key":"e6c3a25b-e576-4d28-4ca5-d8c4af3d0113","doc_count":12,"unique_authors":{"value":2},"commits_count":{"value":12},"commits-per-auth":{"value":6}},{"key":"0f109db8-4f67-429a-8485-4b43b3c94d74","doc_count":11,"unique_authors":{"value":2},"commits_count":{"value":11},"commits-per-auth":{"value":5.5}},{"key":"2a315f19-88de-447c-8ad7-f3bbbd0a64c4","doc_count":11,"unique_authors":{"value":1},"commits_count":{"value":11},"commits-per-auth":{"value":11}},{"key":"8f423498-2f04-4ac3-b948-bdeb1e824a09","doc_count":11,"unique_authors":{"value":2},"commits_count":{"value":11},"commits-per-auth":{"value":5.5}},{"key":"3ad579c8-a286-4df3-acc9-548003ffc8d1","doc_count":9,"unique_authors":{"value":3},"commits_count":{"value":9},"commits-per-auth":{"value":3}},{"key":"c497a177-931e-4ab9-8822-0c4d39d2acfc","doc_count":9,"unique_authors":{"value":1},"commits_count":{"value":9},"commits-per-auth":{"value":9}},{"key":"d87f472d-629e-487e-ba51-948cff7cbf10","doc_count":9,"unique_authors":{"value":1},"commits_count":{"value":9},"commits-per-auth":{"value":9}},{"key":"e69dab25-77dd-43df-9a4f-2c66fbc9b83a","doc_count":9,"unique_authors":{"value":3},"commits_count":{"value":9},"commits-per-auth":{"value":3}},{"key":"edb4138a-4fb8-4bdb-8397-2b64bdad3851","doc_count":9,"unique_authors":{"value":1},"commits_count":{"value":9},"commits-per-auth":{"value":9}},{"key":"bf3322a8-ab44-4e9b-a75f-863ac11737f3","doc_count":8,"unique_authors":{"value":1},"commits_count":{"value":8},"commits-per-auth":{"value":8}},{"key":"e6233eaa-c7dd-43be-8f29-7d9078673392","doc_count":7,"unique_authors":{"value":1},"commits_count":{"value":7},"commits-per-auth":{"value":7}},{"key":"0af5ff0c-c697-4576-a717-0319264e231b","doc_count":6,"unique_authors":{"value":1},"commits_count":{"value":6},"commits-per-auth":{"value":6}},{"key":"3f9855a8-308e-455e-5601-7e967cb3f268","doc_count":6,"unique_authors":{"value":2},"commits_count":{"value":6},"commits-per-auth":{"value":3}},{"key":"4effa6fb-6d72-497f-a730-ae8eb763d6d9","doc_count":6,"unique_authors":{"value":1},"commits_count":{"value":6},"commits-per-auth":{"value":6}},{"key":"50f622f9-4f11-4984-9a97-0a95442e34d3","doc_count":6,"unique_authors":{"value":2},"commits_count":{"value":6},"commits-per-auth":{"value":3}},{"key":"6a3344bb-d280-4139-9b45-c709a7c572de","doc_count":6,"unique_authors":{"value":1},"commits_count":{"value":6},"commits-per-auth":{"value":6}},{"key":"986570f1-d595-4e17-bdd8-f141b181a798","doc_count":6,"unique_authors":{"value":1},"commits_count":{"value":6},"commits-per-auth":{"value":6}},{"key":"a90a60c5-b39d-438b-849c-6b333f765f38","doc_count":6,"unique_authors":{"value":3},"commits_count":{"value":6},"commits-per-auth":{"value":2}},{"key":"cc7041e3-0260-4d9f-b941-87fb440d6aee","doc_count":6,"unique_authors":{"value":1},"commits_count":{"value":6},"commits-per-auth":{"value":6}},{"key":"ebd297e1-afcb-4837-9bc1-2774cc6eb8a0","doc_count":6,"unique_authors":{"value":2},"commits_count":{"value":6},"commits-per-auth":{"value":3}},{"key":"4f3884a6-195d-4d7e-84c8-039665eacd8b","doc_count":5,"unique_authors":{"value":1},"commits_count":{"value":5},"commits-per-auth":{"value":5}},{"key":"6982e936-d26f-45bf-9155-9aeed85a5137","doc_count":5,"unique_authors":{"value":1},"commits_count":{"value":5},"commits-per-auth":{"value":5}},{"key":"7208ab11-afa5-4145-57b3-e3d6df5b09c6","doc_count":5,"unique_authors":{"value":2},"commits_count":{"value":5},"commits-per-auth":{"value":2.5}},{"key":"876797e6-cda2-481b-ba14-67d96f957312","doc_count":5,"unique_authors":{"value":2},"commits_count":{"value":5},"commits-per-auth":{"value":2.5}},{"key":"aab42c4f-1d6a-476b-a0f8-d3d78671134a","doc_count":5,"unique_authors":{"value":1},"commits_count":{"value":5},"commits-per-auth":{"value":5}},{"key":"b4b713e1-b51f-469b-6a2b-b02b074e7271","doc_count":5,"unique_authors":{"value":3},"commits_count":{"value":5},"commits-per-auth":{"value":1.6666666666666667}},{"key":"fb2429a2-3299-4e85-bb30-a394cd252f7c","doc_count":5,"unique_authors":{"value":2},"commits_count":{"value":5},"commits-per-auth":{"value":2.5}},{"key":"1f7f9317-09b8-44ba-b660-aa9d984d707b","doc_count":4,"unique_authors":{"value":2},"commits_count":{"value":4},"commits-per-auth":{"value":2}},{"key":"61e5dc5a-d59c-450c-99ee-009414a22d7a","doc_count":4,"unique_authors":{"value":4},"commits_count":{"value":4},"commits-per-auth":{"value":1}},{"key":"716d4922-add7-474b-7f21-240a8e253c38","doc_count":4,"unique_authors":{"value":1},"commits_count":{"value":4},"commits-per-auth":{"value":4}},{"key":"80a5fe62-4d9c-4385-b2ad-fb6ad99e9b2d","doc_count":4,"unique_authors":{"value":1},"commits_count":{"value":4},"commits-per-auth":{"value":4}},{"key":"a1c741a5-f653-4f46-b082-f121924d50b3","doc_count":4,"unique_authors":{"value":2},"commits_count":{"value":4},"commits-per-auth":{"value":2}},{"key":"a9199a76-90b3-46f4-a2d4-4a4b8c643e59","doc_count":4,"unique_authors":{"value":3},"commits_count":{"value":4},"commits-per-auth":{"value":1.3333333333333333}},{"key":"1381e865-fca5-49ff-a0cf-d67c1d743a76","doc_count":3,"unique_authors":{"value":1},"commits_count":{"value":3},"commits-per-auth":{"value":3}},{"key":"20881a53-1dc9-4c49-98a2-31a08a3824e3","doc_count":3,"unique_authors":{"value":1},"commits_count":{"value":3},"commits-per-auth":{"value":3}},{"key":"40df6d5b-2401-4c1c-a868-1a5a4ea1bbbe","doc_count":3,"unique_authors":{"value":2},"commits_count":{"value":3},"commits-per-auth":{"value":1.5}},{"key":"907808f7-ea55-461c-8a3c-06c79befc919","doc_count":3,"unique_authors":{"value":2},"commits_count":{"value":3},"commits-per-auth":{"value":1.5}},{"key":"f998e0c4-730c-48ea-a816-81ea48d4e3ae","doc_count":3,"unique_authors":{"value":1},"commits_count":{"value":3},"commits-per-auth":{"value":3}},{"key":"0014b65b-e3a5-41d8-54a0-77392c032097","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"087cae5a-db8f-4d34-bef2-1175c7e0d394","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"1113685f-fcbb-4505-b6cf-8f06b36e72a8","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"135e13a4-1017-415c-84c6-204c0b7b42bd","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"1433218b-f6c0-43f4-9f56-34577d905909","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"32648e2c-256d-4b38-979f-5e089f6c34d8","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"34df94d2-2892-4a53-7b0c-f778b6842bf0","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"3c1685a8-734d-4e43-a86c-6f6541a77ad1","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"55ad025f-157a-4af5-98c7-bab36ba1b1b5","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"6bdeb8dc-4a2a-4b20-ad76-64b0e5926949","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"8451aa86-ffe3-4496-81c6-390a49bf92aa","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"9a2e490b-63d1-4dfd-830f-c26a8eccfcb0","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"a82a2754-5c31-45fa-9605-453d56a5aff3","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"aaeae1d8-5e7f-4b44-834a-96ac37ae7bc4","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"b0fe4960-2ad8-4a08-6706-c9460a433c99","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"cb4ffec3-1a8d-4c8a-9fc1-b4d9639b268c","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"cbb6e50b-ff80-4d90-6a61-a75e9c4bb7ce","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"cbe2e16b-59e3-4b75-aa9a-01ecd12372b9","doc_count":2,"unique_authors":{"value":1},"commits_count":{"value":2},"commits-per-auth":{"value":2}},{"key":"4c1ce669-7acc-475f-b5ee-b5c826ff5c3c","doc_count":1,"unique_authors":{"value":1},"commits_count":{"value":1},"commits-per-auth":{"value":1}},{"key":"709da9f4-c61d-4443-9cdd-8dce20f65052","doc_count":1,"unique_authors":{"value":1},"commits_count":{"value":1},"commits-per-auth":{"value":1}},{"key":"77e6f749-d45d-40f0-b0fa-a63c6b5bcfcd","doc_count":1,"unique_authors":{"value":1},"commits_count":{"value":1},"commits-per-auth":{"value":1}},{"key":"78ccc7ea-9c30-4a90-4aea-d3a07df8c382","doc_count":1,"unique_authors":{"value":1},"commits_count":{"value":1},"commits-per-auth":{"value":1}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":1390,"value_in_millis":0,"compare_reports":null,"section":{"data":null}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":331,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":331,"value_in_millis":0,"compare_reports":null,"section":{"data":null}}],"section":{"data":null}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getCommitTrendsComponentComparison("commitTrendsComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing commit trends component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetPullRequestComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of pull request component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"pullRequestsChart":{"took":742,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":626,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"pull_requests_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":112,"pullrequests":{"value":{"CHANGES_REQUESTED":0,"APPROVED":17,"OPEN":4,"REJECTED":0}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":94,"pullrequests":{"value":{"CHANGES_REQUESTED":2,"APPROVED":20,"OPEN":7,"REJECTED":0}}},{"key":"2281041f-7381-491b-aec0-ca2e5cd9061d","doc_count":39,"pullrequests":{"value":{"CHANGES_REQUESTED":0,"APPROVED":8,"OPEN":5,"REJECTED":0}}},{"key":"951c94c4-dd3a-4a24-8d7f-4c0f22df6838","doc_count":35,"pullrequests":{"value":{"CHANGES_REQUESTED":1,"APPROVED":20,"OPEN":0,"REJECTED":0}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":29,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Approved","value":20},{"title":"Changes requested","value":2},{"title":"Open","value":7},{"title":"Rejected","value":0}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":21,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":21,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Approved","value":17},{"title":"Changes requested","value":0},{"title":"Open","value":4},{"title":"Rejected","value":0}]}}],"section":{"data":[{"title":"Approved","value":17},{"title":"Changes requested","value":0},{"title":"Open","value":4},{"title":"Rejected","value":0}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getPullRequestComponentComparison("pullRequestComponentSpec", x, replacements, organisation)
		// bString, _ := json.Marshal(b)
		// fmt.Println(string(bString))
		assert.Nil(t, err, "error processing pull request component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetWorkflowRunsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of workflow runs component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"runsStatusChart":{"took":147,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"workflow_runs_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"162239c9-e4dd-4c7c-9c79-a97841e625e3","doc_count":8626,"automation_run":{"value":{"Success":8310,"Failure":145}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":1346,"automation_run":{"value":{"Success":608,"Failure":495}}},{"key":"545c6a93-7b98-49b4-707e-9c9ab52ce617","doc_count":308,"automation_run":{"value":{"Success":93,"Failure":160}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":283,"automation_run":{"value":{"Success":227,"Failure":50}}},{"key":"94a81dd1-3f52-4520-891e-b2440f660945","doc_count":233,"automation_run":{"value":{"Success":210,"Failure":21}}},{"key":"238ffe68-8cb4-459d-64ac-2e4f752fe8dc","doc_count":178,"automation_run":{"value":{"Success":122,"Failure":37}}},{"key":"6cbb8b60-16d6-48a4-61fe-3a33c6c1169f","doc_count":143,"automation_run":{"value":{"Success":85,"Failure":39}}},{"key":"9b707c3e-b9a2-4705-6c44-6f9f250f61c6","doc_count":138,"automation_run":{"value":{"Success":40,"Failure":71}}},{"key":"d2ab028a-4902-4207-4352-23e412c24884","doc_count":134,"automation_run":{"value":{"Success":92,"Failure":26}}},{"key":"590947f7-4640-422b-a8a9-36aae8b85ee6","doc_count":131,"automation_run":{"value":{"Success":76,"Failure":26}}},{"key":"3891d690-a504-4889-4ab1-e469db088597","doc_count":119,"automation_run":{"value":{"Success":99,"Failure":20}}},{"key":"6e5a4968-1806-4ff4-5052-0cd9859c4d93","doc_count":109,"automation_run":{"value":{"Success":78,"Failure":28}}},{"key":"e0fb1caf-4c53-4be5-644e-1836524bab2a","doc_count":107,"automation_run":{"value":{"Success":39,"Failure":65}}},{"key":"55f26650-b259-4453-6ff3-c2843d096780","doc_count":104,"automation_run":{"value":{"Success":103,"Failure":1}}},{"key":"785894c9-91cf-47b9-734d-ae73377cc29c","doc_count":103,"automation_run":{"value":{"Success":80,"Failure":15}}},{"key":"cdcfa229-facc-4e9a-80a1-d82211f910b4","doc_count":103,"automation_run":{"value":{"Success":64,"Failure":31}}},{"key":"ae694ee5-6234-419c-5035-f3d59622565d","doc_count":97,"automation_run":{"value":{"Success":78,"Failure":16}}},{"key":"839a6ffc-3c14-41ca-8621-532e768ea9d5","doc_count":96,"automation_run":{"value":{"Success":42,"Failure":43}}},{"key":"bdd99ab8-12b8-49c8-61b5-087796f009f3","doc_count":94,"automation_run":{"value":{"Success":54,"Failure":35}}},{"key":"363c6b99-4c67-4e02-9bda-10585eb25d72","doc_count":87,"automation_run":{"value":{"Success":73,"Failure":8}}},{"key":"7ee1c5e5-0564-43d9-7af7-ba46907f5dd0","doc_count":80,"automation_run":{"value":{"Success":44,"Failure":28}}},{"key":"65240b81-48d8-44a7-4592-82a9e7a1a132","doc_count":79,"automation_run":{"value":{"Success":76,"Failure":2}}},{"key":"743ddab4-955d-4282-7230-eb2cf50886da","doc_count":74,"automation_run":{"value":{"Success":36,"Failure":25}}},{"key":"16b2c653-5f73-41d1-4765-5b8ae36fcfe2","doc_count":63,"automation_run":{"value":{"Success":34,"Failure":24}}},{"key":"746a9c25-df05-420e-5f4e-f65f687d777e","doc_count":57,"automation_run":{"value":{"Success":51,"Failure":6}}},{"key":"05920b1f-02d6-47de-af20-d4dd209e8b17","doc_count":52,"automation_run":{"value":{"Success":32,"Failure":17}}},{"key":"1063a418-d663-4503-af48-fe5c6ff4eb68","doc_count":52,"automation_run":{"value":{"Success":48,"Failure":1}}},{"key":"cf208105-08b5-44ad-842c-5fc1dfae3c76","doc_count":51,"automation_run":{"value":{"Success":19,"Failure":21}}},{"key":"e2f8fef6-5041-4843-b37e-6cdae38099bc","doc_count":49,"automation_run":{"value":{"Success":39,"Failure":9}}},{"key":"dda69191-5492-4b7e-88b2-9d9d42f61899","doc_count":48,"automation_run":{"value":{"Success":36,"Failure":10}}},{"key":"9af098ae-d223-4a2f-a185-e6557db9d8a1","doc_count":45,"automation_run":{"value":{"Success":15,"Failure":21}}},{"key":"ab21947d-d035-4794-b4f6-a90421fc2bf7","doc_count":41,"automation_run":{"value":{"Success":36,"Failure":2}}},{"key":"054d5f8b-c1bf-4c2f-41d0-73e33dfcad8b","doc_count":37,"automation_run":{"value":{"Success":15,"Failure":15}}},{"key":"31e661ac-6553-4d6c-6f3c-6528f64b8fcd","doc_count":37,"automation_run":{"value":{"Success":37,"Failure":0}}},{"key":"d87f472d-629e-487e-ba51-948cff7cbf10","doc_count":31,"automation_run":{"value":{"Success":3,"Failure":21}}},{"key":"d5c3f808-d0a1-4e9d-4ca3-47d82aca3257","doc_count":30,"automation_run":{"value":{"Success":24,"Failure":5}}},{"key":"0f109db8-4f67-429a-8485-4b43b3c94d74","doc_count":29,"automation_run":{"value":{"Success":0,"Failure":19}}},{"key":"f783334c-97b3-478e-8788-2d8e76c91d91","doc_count":28,"automation_run":{"value":{"Success":12,"Failure":12}}},{"key":"3a04548e-9af3-4c45-6aa2-937c13baaa40","doc_count":27,"automation_run":{"value":{"Success":18,"Failure":7}}},{"key":"61751a40-2830-4783-8528-6ce7fad549b9","doc_count":27,"automation_run":{"value":{"Success":6,"Failure":14}}},{"key":"896a0752-03bf-4531-946f-ab2c88b14a02","doc_count":27,"automation_run":{"value":{"Success":6,"Failure":13}}},{"key":"dfa63f63-d55d-4fe2-7c6e-dc429133593f","doc_count":27,"automation_run":{"value":{"Success":23,"Failure":4}}},{"key":"c4398680-af38-4e76-bab8-0f333042a891","doc_count":25,"automation_run":{"value":{"Success":11,"Failure":9}}},{"key":"a677dedb-dd24-42b5-8f55-16594562f0de","doc_count":24,"automation_run":{"value":{"Success":5,"Failure":12}}},{"key":"d000292d-bd9c-44c9-a3f5-598ec3a561f4","doc_count":24,"automation_run":{"value":{"Success":8,"Failure":10}}},{"key":"3ad579c8-a286-4df3-acc9-548003ffc8d1","doc_count":22,"automation_run":{"value":{"Success":5,"Failure":8}}},{"key":"3c09cd3a-d954-4548-7d9a-f1807fb00daf","doc_count":21,"automation_run":{"value":{"Success":13,"Failure":8}}},{"key":"81d93aec-19a0-4258-59a5-55b7f7140947","doc_count":21,"automation_run":{"value":{"Success":8,"Failure":10}}},{"key":"a1d40803-7978-40ed-6c6e-d6fa114cbfa7","doc_count":20,"automation_run":{"value":{"Success":9,"Failure":8}}},{"key":"e8f3a95e-2387-40d0-b727-5a542584e4b0","doc_count":20,"automation_run":{"value":{"Success":5,"Failure":12}}},{"key":"3192697f-a073-49d0-ae6f-a1159152ec76","doc_count":18,"automation_run":{"value":{"Success":15,"Failure":2}}},{"key":"9abaae9f-4e37-4c36-8424-53cfd9650a13","doc_count":18,"automation_run":{"value":{"Success":14,"Failure":3}}},{"key":"c0fb4e57-38e7-4522-6b07-43353312800a","doc_count":18,"automation_run":{"value":{"Success":12,"Failure":4}}},{"key":"2a315f19-88de-447c-8ad7-f3bbbd0a64c4","doc_count":17,"automation_run":{"value":{"Success":9,"Failure":5}}},{"key":"89e389cc-3f56-429c-6ee1-1a8bfc8ef5ae","doc_count":16,"automation_run":{"value":{"Success":16,"Failure":0}}},{"key":"92df4461-d43f-4d19-5e00-4b6346f747d9","doc_count":16,"automation_run":{"value":{"Success":12,"Failure":4}}},{"key":"e6233eaa-c7dd-43be-8f29-7d9078673392","doc_count":16,"automation_run":{"value":{"Success":7,"Failure":7}}},{"key":"4d72bf77-ce60-4db7-6d26-7f3adff30fa6","doc_count":15,"automation_run":{"value":{"Success":4,"Failure":7}}},{"key":"aab42c4f-1d6a-476b-a0f8-d3d78671134a","doc_count":14,"automation_run":{"value":{"Success":8,"Failure":3}}},{"key":"e6c3a25b-e576-4d28-4ca5-d8c4af3d0113","doc_count":14,"automation_run":{"value":{"Success":5,"Failure":7}}},{"key":"6a3344bb-d280-4139-9b45-c709a7c572de","doc_count":12,"automation_run":{"value":{"Success":12,"Failure":0}}},{"key":"3cd7e6f3-8fcb-426f-791e-6285150c6994","doc_count":11,"automation_run":{"value":{"Success":11,"Failure":0}}},{"key":"986570f1-d595-4e17-bdd8-f141b181a798","doc_count":11,"automation_run":{"value":{"Success":7,"Failure":2}}},{"key":"d15bfc5b-1b01-4362-9851-f51388121e0f","doc_count":11,"automation_run":{"value":{"Success":8,"Failure":2}}},{"key":"b4b713e1-b51f-469b-6a2b-b02b074e7271","doc_count":10,"automation_run":{"value":{"Success":0,"Failure":8}}},{"key":"bf3322a8-ab44-4e9b-a75f-863ac11737f3","doc_count":10,"automation_run":{"value":{"Success":7,"Failure":2}}},{"key":"cc7041e3-0260-4d9f-b941-87fb440d6aee","doc_count":10,"automation_run":{"value":{"Success":4,"Failure":4}}},{"key":"d1e786b7-c2f0-48f4-9055-741e7ad2b171","doc_count":10,"automation_run":{"value":{"Success":0,"Failure":7}}},{"key":"e69dab25-77dd-43df-9a4f-2c66fbc9b83a","doc_count":10,"automation_run":{"value":{"Success":9,"Failure":1}}},{"key":"3f9855a8-308e-455e-5601-7e967cb3f268","doc_count":9,"automation_run":{"value":{"Success":8,"Failure":1}}},{"key":"8f423498-2f04-4ac3-b948-bdeb1e824a09","doc_count":9,"automation_run":{"value":{"Success":9,"Failure":0}}},{"key":"80a5fe62-4d9c-4385-b2ad-fb6ad99e9b2d","doc_count":8,"automation_run":{"value":{"Success":2,"Failure":4}}},{"key":"1433218b-f6c0-43f4-9f56-34577d905909","doc_count":7,"automation_run":{"value":{"Success":0,"Failure":4}}},{"key":"20881a53-1dc9-4c49-98a2-31a08a3824e3","doc_count":7,"automation_run":{"value":{"Success":7,"Failure":0}}},{"key":"4f3884a6-195d-4d7e-84c8-039665eacd8b","doc_count":6,"automation_run":{"value":{"Success":6,"Failure":0}}},{"key":"50f622f9-4f11-4984-9a97-0a95442e34d3","doc_count":6,"automation_run":{"value":{"Success":5,"Failure":1}}},{"key":"69b3cf79-c1bd-47a6-934c-d2a83c7741e9","doc_count":6,"automation_run":{"value":{"Success":0,"Failure":4}}},{"key":"7208ab11-afa5-4145-57b3-e3d6df5b09c6","doc_count":6,"automation_run":{"value":{"Success":6,"Failure":0}}},{"key":"a1c741a5-f653-4f46-b082-f121924d50b3","doc_count":6,"automation_run":{"value":{"Success":0,"Failure":4}}},{"key":"cbe2e16b-59e3-4b75-aa9a-01ecd12372b9","doc_count":6,"automation_run":{"value":{"Success":0,"Failure":3}}},{"key":"1f7f9317-09b8-44ba-b660-aa9d984d707b","doc_count":5,"automation_run":{"value":{"Success":3,"Failure":1}}},{"key":"34df94d2-2892-4a53-7b0c-f778b6842bf0","doc_count":5,"automation_run":{"value":{"Success":5,"Failure":0}}},{"key":"61e5dc5a-d59c-450c-99ee-009414a22d7a","doc_count":5,"automation_run":{"value":{"Success":5,"Failure":0}}},{"key":"716d4922-add7-474b-7f21-240a8e253c38","doc_count":5,"automation_run":{"value":{"Success":3,"Failure":1}}},{"key":"876797e6-cda2-481b-ba14-67d96f957312","doc_count":5,"automation_run":{"value":{"Success":4,"Failure":1}}},{"key":"c497a177-931e-4ab9-8822-0c4d39d2acfc","doc_count":5,"automation_run":{"value":{"Success":2,"Failure":3}}},{"key":"ebd297e1-afcb-4837-9bc1-2774cc6eb8a0","doc_count":5,"automation_run":{"value":{"Success":0,"Failure":4}}},{"key":"36daa845-8804-4424-4bf7-9b67c1a408d4","doc_count":4,"automation_run":{"value":{"Success":0,"Failure":2}}},{"key":"40df6d5b-2401-4c1c-a868-1a5a4ea1bbbe","doc_count":4,"automation_run":{"value":{"Success":2,"Failure":1}}},{"key":"6982e936-d26f-45bf-9155-9aeed85a5137","doc_count":4,"automation_run":{"value":{"Success":1,"Failure":2}}},{"key":"a9199a76-90b3-46f4-a2d4-4a4b8c643e59","doc_count":4,"automation_run":{"value":{"Success":4,"Failure":0}}},{"key":"f998e0c4-730c-48ea-a816-81ea48d4e3ae","doc_count":4,"automation_run":{"value":{"Success":0,"Failure":3}}},{"key":"0af5ff0c-c697-4576-a717-0319264e231b","doc_count":3,"automation_run":{"value":{"Success":3,"Failure":0}}},{"key":"cb4ffec3-1a8d-4c8a-9fc1-b4d9639b268c","doc_count":3,"automation_run":{"value":{"Success":1,"Failure":2}}},{"key":"08e2bfd6-4c47-45e6-9578-d1b11d8dfbd4","doc_count":2,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"0e80c361-f9de-492d-a751-1f1af6632306","doc_count":2,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"135e13a4-1017-415c-84c6-204c0b7b42bd","doc_count":2,"automation_run":{"value":{"Success":0,"Failure":2}}},{"key":"1381e865-fca5-49ff-a0cf-d67c1d743a76","doc_count":2,"automation_run":{"value":{"Success":1,"Failure":1}}},{"key":"3c1685a8-734d-4e43-a86c-6f6541a77ad1","doc_count":2,"automation_run":{"value":{"Success":2,"Failure":0}}},{"key":"56c4fb46-d533-4309-8452-74b1c6a90063","doc_count":2,"automation_run":{"value":{"Success":2,"Failure":0}}},{"key":"77e6f749-d45d-40f0-b0fa-a63c6b5bcfcd","doc_count":2,"automation_run":{"value":{"Success":2,"Failure":0}}},{"key":"7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a","doc_count":2,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"7cd1d3eb-eb08-4fa3-a19d-24a8986dbe5b","doc_count":2,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"7dfaabc0-5eb0-40c6-427c-a2ac4cc01da3","doc_count":2,"automation_run":{"value":{"Success":2,"Failure":0}}},{"key":"85b1f8b6-421d-403d-9bb1-ebbd05732958","doc_count":2,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"9acd89fd-6408-4ff7-ab1e-d1af2ed3e206","doc_count":2,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"b7bd2323-4771-4810-6025-646a5c0fb238","doc_count":2,"automation_run":{"value":{"Success":0,"Failure":2}}},{"key":"c2daa34d-2432-4de0-b98b-1c4a0c63abbb","doc_count":2,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"faa35709-87d7-4e8a-85f4-feef324e170b","doc_count":2,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"0014b65b-e3a5-41d8-54a0-77392c032097","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"06c42a84-b326-4614-b10d-85d8db59bed1","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"087cae5a-db8f-4d34-bef2-1175c7e0d394","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"0c3c691a-28a4-4779-afb5-da0f8ca4e6d3","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"1415fbd8-3c58-4363-9ce7-f2483277f6b8","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"24c2c304-75b1-4076-9117-30d5d375ab15","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"2b206de1-0ad9-4798-920e-08f67c61128b","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"33b69515-ba92-4519-a309-731b44a7fdf8","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"37534a1f-25ac-4b05-5f01-2f831b213ac3","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"3d64036f-1b28-400a-8f6d-d0f4ecb81c86","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"3e280b19-e453-410b-be85-a909213811ee","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"4191e255-6957-4f78-80f4-6fa5062ad336","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"459144b0-bd7e-40d4-85d7-dff9e7d628e6","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"45abf2e8-65c8-4787-adef-43c90144edd6","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"4951e6cc-75c9-4324-a51c-51b8a3d0de5e","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"4baf114a-b71d-4077-ac84-31069e7d32a0","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"4c1ce669-7acc-475f-b5ee-b5c826ff5c3c","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"4fe12a6d-9329-468f-ae6e-e8591f5b759b","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"55ad025f-157a-4af5-98c7-bab36ba1b1b5","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"5cc63bb0-f3d4-41c7-a38c-48a4bbf3ed57","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"64dfe4fa-47f3-4d50-a35b-91b7dac6f506","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"709da9f4-c61d-4443-9cdd-8dce20f65052","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"78800f06-afee-4771-96c2-29464d3833ed","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"78ccc7ea-9c30-4a90-4aea-d3a07df8c382","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"8451aa86-ffe3-4496-81c6-390a49bf92aa","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"956130b8-761d-4033-9968-3185b09cb24e","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"99bd9baf-d645-4a70-5848-db417e8c7301","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"9a2e490b-63d1-4dfd-830f-c26a8eccfcb0","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"9b292194-a17c-4417-9a6b-d4fc5bab3b9d","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"a42277e3-c8a3-41f1-a970-26aee1296d80","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"a82a2754-5c31-45fa-9605-453d56a5aff3","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"a90a60c5-b39d-438b-849c-6b333f765f38","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"addf84e3-5c30-49d0-aaa7-1564ee109c0b","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"bc8994ee-9bed-4a9f-bc4c-9fd32b97d326","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"c432f062-ee61-4e9a-a21c-b3e83d25486a","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"cb0be3c6-d203-4a26-a9bd-4dc803b793aa","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"d153c6f9-3809-4be9-b0db-113e8e4a0fe4","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"d569c6f1-13a5-40ac-8d8f-d4b9a00ac44a","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"daf7382a-a3d8-4b37-ae2e-03e8f4c7eb6e","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"e0674e0a-4576-4193-8bb4-6d529db962fe","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"e48d754f-46bf-4fab-9b82-98c270fa0d42","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"e58ab30d-069a-4986-8ce1-22e4a0ff2d93","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"ed35ceee-42b2-40b2-85fe-d3407746eb72","doc_count":1,"automation_run":{"value":{"Success":1,"Failure":0}}},{"key":"f1fa3c00-603a-4621-9e71-b5892b895ed0","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}},{"key":"f3008774-02f0-45d4-8b8c-facff407890e","doc_count":1,"automation_run":{"value":{"Success":0,"Failure":1}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":1103,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Success","value":608},{"title":"Failure","value":495}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":277,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":277,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Success","value":227},{"title":"Failure","value":50}]}}],"section":{"data":[{"title":"Success","value":227},{"title":"Failure","value":50}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getWorkflowRunsComponentComparison("workflowRunsComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing workflow runs component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetDevCycleTimeComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of dev cycle time component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"developmentCycleChart":{"took":26,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":1903,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"development_cycle_time_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"4180826-bb76-421e-5410-791408daadeb","doc_count":336,"developmentCycleTime":{"value":{"coding_time_count":43,"coding_time":"17h 6m ","review_time_count":8,"deploy_time_in_millis":0,"coding_time_in_millis":61594000,"pickup_time_in_millis":82716000,"review_time_value_in_millis":2194192467,"pickup_time_value_in_millis":827164000,"review_time":"3h 5m ","coding_time_value_in_millis":2648578000,"pickup_time":"22h 58m ","review_time_in_millis":11120000,"deploy_time":"","deploy_time_count":0,"deploy_time_value_in_millis":0,"pickup_time_count":10}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":230,"developmentCycleTime":{"value":{"coding_time_count":20,"coding_time":"12h 9m ","review_time_count":60,"deploy_time_in_millis":2194192467,"coding_time_in_millis":10000,"pickup_time_in_millis":20000,"review_time_value_in_millis":2194192467,"pickup_time_value_in_millis":20000,"review_time":"49m ","coding_time_value_in_millis":10000,"pickup_time":"4h 10m ","review_time_in_millis":2993000,"deploy_time":"","deploy_time_count":0,"deploy_time_value_in_millis":0,"pickup_time_count":10}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":130,"developmentCycleTime":{"value":{"coding_time_count":50,"coding_time":"3h 37m ","review_time_count":30,"deploy_time_in_millis":0,"coding_time_in_millis":10000,"pickup_time_in_millis":20000,"review_time_value_in_millis":2194192467,"pickup_time_value_in_millis":20000,"review_time":"7h 23m ","coding_time_value_in_millis":10000,"pickup_time":"1d 18m ","review_time_in_millis":26623000,"deploy_time":"","deploy_time_count":0,"deploy_time_value_in_millis":0,"pickup_time_count":20}}},{"key":"94a81dd1-3f52-4520-891e-b2440f660945","doc_count":124,"developmentCycleTime":{"value":{"coding_time_count":22,"coding_time":"3h 31m ","review_time_count":14,"deploy_time_in_millis":0,"coding_time_in_millis":12660000,"pickup_time_in_millis":39756000,"review_time_value_in_millis":2194192467,"pickup_time_value_in_millis":516834000,"review_time":"15m ","coding_time_value_in_millis":278537000,"pickup_time":"11h 2m ","review_time_in_millis":904000,"deploy_time":"","deploy_time_count":0,"deploy_time_value_in_millis":0,"pickup_time_count":13}}},{"key":"cb4ffec3-1a8d-4c8a-9fc1-b4d9639b268c","doc_count":90,"developmentCycleTime":{"value":{"coding_time_count":11,"coding_time":"10h 7m ","review_time_count":9,"deploy_time_in_millis":0,"coding_time_in_millis":36424000,"pickup_time_in_millis":3422000,"review_time_value_in_millis":2194192467,"pickup_time_value_in_millis":27378000,"review_time":"8h 14m ","coding_time_value_in_millis":400674000,"pickup_time":"57m ","review_time_in_millis":29687000,"deploy_time":"","deploy_time_count":0,"deploy_time_value_in_millis":0,"pickup_time_count":8}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":73140948.9,"compare_reports":null,"section":{"data":[{"title":"Coding time","value":0,"time":10000,"count":50},{"title":"Code pickup time","value":0,"time":20000,"count":20},{"title":"Code review time","value":100,"time":2194192467,"count":30}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":36572374,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":36572374.45,"compare_reports":null,"section":{"data":[{"title":"Coding time","value":0,"time":10000,"count":20},{"title":"Code pickup time","value":0,"time":20000,"count":10},{"title":"Code review time","value":100,"time":2194192467,"count":60}]}}],"section":{"data":[{"title":"Coding time","value":0,"time":10000,"count":20},{"title":"Code pickup time","value":0,"time":20000,"count":10},{"title":"Code review time","value":100,"time":2194192467,"count":60}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getDevCycleTimeComponentComparison("devCycleTimeComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing dev cycle time component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetCommitsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of commits component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"commitsChart":{"took":51,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":279,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"commits_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":214,"automation_run":{"value":{"data":{"a4a6e438-7112-42b2-81d1-93ee5599dc63_27aa76c3-4b6b-4d98-bd15-7a2a907a6dcc":1,"0ae35b79-60ca-41e1-ab04-2a9d5954b0b7_f029a967-65db-4094-b7b0-4c72853f8fbb":1,"a4a6e438-7112-42b2-81d1-93ee5599dc63_87d83a45-52dc-4f08-9c26-16641fa5a7db":1,"34de117c-4934-4829-8346-3c114bc43b88_2d8bc5c5-81cc-4ec6-a179-c67ce6faa8bd":1,"e1b92ded-9edc-4400-8326-13e37c1e5ed5_76972243-f9d0-423f-b37d-865a6243c1e6":1,"0ae35b79-60ca-41e1-ab04-2a9d5954b0b7_028903ca-15b9-402d-a6fc-85cf734d7466":1,"b4f650cc-c297-4bce-8104-0ee1d0bae006_10f5b4a2-ea20-4923-b93f-290189a31dc9":1,"b4f650cc-c297-4bce-8104-0ee1d0bae006_e681afba-00cc-4fb9-b589-97e69e157b3d":1,"17a86fbc-1b4d-4976-967f-ea3f66bcf6aa_543f37e9-59a3-4543-9725-b67c368c29b5":1,"e2b3f55d-139f-4831-854f-b20b86eb9307_a386b49b-21ae-4f24-bbf9-9698e4c1f19e":1,"34de117c-4934-4829-8346-3c114bc43b88_2a2df291-d480-45b4-b347-fb3431d8db19":1,"92dd1bc5-086c-4eb5-a3d3-83548285357e_93a5f6cb-df39-4e80-94da-c586fc04bc5c":1},"totalCount":180}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":100,"automation_run":{"value":{"data":{"7a0584ed-6ad5-4e00-8218-d5af44f730fc_f0360509-f756-40e2-9b24-beedae08cc5a":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_5f1cda45-5fc2-481d-8c64-6396ea2e8d35":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_595b6881-0d4d-4843-ba31-613f6cdc1a0c":1,"0a025d96-f0c9-49e9-8e10-5c8f4bb39949_338753d9-aa21-4c12-a546-3c50faec8e50":1,"e7a03e3c-19e7-4169-b039-351164858a7c_56a07bd7-09d4-46c0-a4f6-daf580a0b83e":1,"915a8eb8-4ccc-46cf-bdc0-54da661310fc_32a9c2a4-8eff-42e9-b1e0-b666c7fa3622":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_2b87088f-8166-4415-a675-08a343dd2b36":1,"bdb5abed-4f93-4217-b8ec-3f14ebf541f5_680d9499-7054-4fd9-8146-c4731a975db1":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_a08f48e7-80bd-46fb-9d6e-6ed08aae6103":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_21b38a82-4a89-4ad8-b614-cf06b421bdfb":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_1e61d2da-d240-49e7-a417-13fc736d038c":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_4e792bb1-5d7e-4373-bfb1-d6e9acc51849":1,"915a8eb8-4ccc-46cf-bdc0-54da661310fc_4e4fa03a-1995-4cba-a9f0-a0171e95b2ea":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_d74fcdbf-d67d-4ed2-a63e-0b43ce1bf8e9":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_b4eb4025-62d1-4ad2-a3ce-a30a57779bea":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_8f32ce94-2b49-4910-b45c-403f81695867":1,"ddc27082-7c1d-446c-8e06-23eaa8156d06_5d8288c1-fc42-44a3-aced-8de9b3b1abdc":1,"bdb5abed-4f93-4217-b8ec-3f14ebf541f5_39667f01-6288-444d-b1ae-70e771cf0630":1,"b1cd4be3-979c-4334-80b8-bc4b7e7af6be_b6952413-1ba0-40c6-98db-78efa4ad41df":1,"0bd44a32-3d52-4273-a808-dd9d32adcaab_7dd9516b-de0c-4c6d-8ec1-803c9d0320c5":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_91cc4cfd-7149-498a-a82a-83921cd3764f":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_ad8ef1bd-098a-4eeb-bb5e-f4be0e2cdc34":1,"fff50234-836e-4962-9287-ad36b1ff6af0_c6d04c92-6112-451d-991e-ec88c90681ef":1,"e7a03e3c-19e7-4169-b039-351164858a7c_64b10008-35aa-467c-b259-ef1e81424041":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_efca7392-3ff0-4e77-a699-d4d359730346":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_3bb8b65e-6636-4ed5-8fca-90d251a0a677":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_f5a07fd9-e24b-4752-bbb7-39cfd3f55ed4":1,"e7a03e3c-19e7-4169-b039-351164858a7c_67d30afc-850e-4096-9005-937db27e5c03":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_67cfe5db-5824-4261-a412-825373754a6a":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_4467583b-fd63-4e09-82ac-2ce60054c78b":1,"7a0584ed-6ad5-4e00-8218-d5af44f730fc_b8a43010-05f2-43f1-9c11-32e8fdf962aa":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_9819dea5-e348-493b-82f0-bf1b0e04ed51":1,"fd15e407-c680-4ba1-8a1f-17c3bbfe6d57_509c6893-fa84-41cc-b5f1-ccba9ae6ce03":1,"5c6e6a5a-fba0-4c6d-9697-26e3dd8229aa_2da03f84-33d8-4bef-88b4-3d58a8f84002":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_0c2888b3-bcc8-4a60-9b80-1c8dc3c67efb":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_359ab23c-c5e2-4961-a363-d0c3d19c4f8f":1,"6edae5a4-9af5-46c4-80f2-45c77d672685_3868de76-9843-4b30-9397-3fc83941a12b":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_122cd117-4722-4a4d-aef0-9b1011aef706":1,"e3b1d6c6-7f43-40a6-9b94-342cf3c07e43_155aa3fb-2f54-4bea-9c0f-232c0dd456c7":1,"fd15e407-c680-4ba1-8a1f-17c3bbfe6d57_663edb28-b287-4f59-a19a-c22d8fc4f846":1,"6b7a73d9-cafe-416d-b75c-349c5f8b6994_9daf334d-967e-444b-9bd2-29efa3dabfe4":1,"961ab526-4796-4ba7-b9e6-278bc62a87e9_d490758c-326e-4e4c-9757-b5fdccecae9b":1,"e3b1d6c6-7f43-40a6-9b94-342cf3c07e43_38bea2ee-b0b3-492c-be3e-cb7668d8418b":1,"961ab526-4796-4ba7-b9e6-278bc62a87e9_f362b4a1-cc33-472c-a91c-897497ab6cf8":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_6e37f223-d088-422e-9d43-1c93b90a8a7e":1,"5c6e6a5a-fba0-4c6d-9697-26e3dd8229aa_79452aef-390c-40c1-a164-cc2e3e8f4434":1,"6b7a73d9-cafe-416d-b75c-349c5f8b6994_b5e99c7b-8e8f-42bb-8fda-129c285a81bc":1,"bdb5abed-4f93-4217-b8ec-3f14ebf541f5_1bf0d5f8-68d4-408d-8b33-c9bbd3bfd4a1":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_ae62a18a-e507-4bbb-b58f-cbdd31a66f6d":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_a05996d2-0905-4d1e-9c0d-2a7fafb7c288":1,"fff50234-836e-4962-9287-ad36b1ff6af0_88a95541-d7b2-4d43-b114-4fb0240badbd":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_fc27d601-0a77-4e4c-9b68-626e65095abf":1,"fff50234-836e-4962-9287-ad36b1ff6af0_c5b3a934-6b38-45ee-91cb-5285bbb503db":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_6d081747-a3a6-400b-b762-1c1add671b24":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_cf586387-5c2e-43ea-96ef-233417094ff4":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_6554f464-0370-4e89-b1a8-d08ba07a74e2":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_49fe411b-88ae-40db-ae29-96568f1c6de0":1,"915a8eb8-4ccc-46cf-bdc0-54da661310fc_87ea3640-b304-4c7b-a05d-9fa336fe5ba3":1,"fd15e407-c680-4ba1-8a1f-17c3bbfe6d57_e391f548-3d43-43e7-94b5-3de1989ea7c0":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_cf7cf0e0-f578-454d-9bdf-4e5087261418":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_27567cbb-4b8e-48f0-b024-2025813a2d9a":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_1a70c2fc-63db-4af5-89de-429ba2518a84":1,"e7a03e3c-19e7-4169-b039-351164858a7c_9417ad8a-7ef8-424d-82bc-f68416124cc8":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_99be2fd0-907e-43ca-a3d1-d0986d84bb19":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_00ffeb8b-16f1-40d4-984f-31225d66188a":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_2a32e47b-dd63-450f-8a3f-d70e514fcae9":1,"0a025d96-f0c9-49e9-8e10-5c8f4bb39949_1c46955d-6cd8-44c8-aeaf-4f44bd752ec2":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_a63a843e-0292-432a-b1b7-6382b54293e4":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_ef555473-a0ac-4adf-bb3f-4aca1bf3f6f2":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_94ae726b-883b-45f5-8068-d508d03f5331":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_51d77dc7-8cf6-4fa6-be68-bdd6a5a87f4f":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_d10d5f05-4fdf-4804-954a-49d5bfbbc527":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_dce2940b-51f3-4ceb-8c47-5888672058da":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_3b5cd1b2-a1cc-40b2-89c0-6d3e0f7cdcbc":1,"fd15e407-c680-4ba1-8a1f-17c3bbfe6d57_3d2db8e1-a37d-4dfd-a000-61ea0c5434bd":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_d480e951-21c3-4422-9dc6-f50ae24f51f7":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_f29b9f1c-98c8-48a6-a239-beda9d29f8ef":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_5ca7a8a4-851e-4263-874a-f66672c14cc6":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_ae15847c-c729-41a3-8801-5611d385da59":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_26b9dbfe-57d5-43f3-a290-5692c34ff576":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_9ed71933-2efd-40c4-97b2-620cbdb90feb":1,"6b7a73d9-cafe-416d-b75c-349c5f8b6994_65b547f2-f072-424b-b2cc-541342a561a6":1,"0bd44a32-3d52-4273-a808-dd9d32adcaab_df8be8df-58c1-46ae-84e7-167a1e918221":1,"a667d3d4-be4c-413d-8b9e-c978cfe7ba19_c6d322a4-18d7-4b46-92d7-c1e05bc0815c":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_f36d0d6d-9017-4d91-8882-aea40e55f17a":1,"7a0584ed-6ad5-4e00-8218-d5af44f730fc_61a17448-5960-4abe-9a92-1f6e3fbd68b8":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_da82034e-dd68-45ac-a15a-f2d78ae0ab02":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_e99a073a-919c-44ca-9faf-f39311ce29c5":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_1acbb842-7daf-43ff-8a40-ccb1721ed436":1,"bdb5abed-4f93-4217-b8ec-3f14ebf541f5_b614e707-e8f6-494d-b3c0-abc15b770672":1,"ccd63de3-e7a4-42c1-b358-b79b9e8fffc7_035bed76-5120-40a9-8abb-b73a0ffc2ddd":1,"35ae0c75-f228-4daf-80a9-9d01a6ad77c5_97cce0d5-9a95-4c24-b031-781abe0e5a92":1,"d2dd323a-985e-48f1-9537-1d97ffecacc9_cb9dd766-410b-496e-bead-579a92d36ef1":1,"f6574f20-16f1-44e7-9385-bd3e3279dc37_801bf153-ed58-4a21-b9d6-f50205daf7aa":1,"7a0584ed-6ad5-4e00-8218-d5af44f730fc_d54d768c-3b79-4db6-98db-0abf7f2f0bca":1,"915a8eb8-4ccc-46cf-bdc0-54da661310fc_4003d75e-b6da-42b1-8b31-4c23755c030e":1,"0bd44a32-3d52-4273-a808-dd9d32adcaab_8b4dd7df-9f85-4a76-aa1c-373435739844":1,"a73f77af-a0ce-4f54-a1b6-ac416be813f4_e4b3c744-72d8-457a-859b-f6d3e6208bb6":1,"fff50234-836e-4962-9287-ad36b1ff6af0_c93c9d0d-46d6-4d4c-8811-004c117a52b2":1},"totalCount":99}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":180,"value_in_millis":0,"compare_reports":null,"section":{"data":null}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":99,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":99,"value_in_millis":0,"compare_reports":null,"section":{"data":null}}],"section":{"data":null}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getCommitsComponentComparison("commitsComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing commits component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetBuildsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of builds component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"buildsChart":{"took":235,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"builds_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":26191,"build_status":{"value":{"data":[{"name":"Success","value":98},{"name":"Failure","value":2}],"info":[{"drillDown":{"reportType":"status","reportId":"builds","reportTitle":"Builds"},"title":"Success","value":8314},{"drillDown":{"reportType":"status","reportId":"builds","reportTitle":"Builds"},"title":"Failure","value":141}]}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":92,"build_status":{"value":{"data":[{"name":"Success","value":99},{"name":"Failure","value":1}],"info":[{"drillDown":{"reportType":"status","reportId":"builds","reportTitle":"Builds"},"title":"Success","value":86},{"drillDown":{"reportType":"status","reportId":"builds","reportTitle":"Builds"},"title":"Failure","value":1}]}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":8455,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Success","value":8314},{"title":"Failure","value":141}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":87,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":87,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Success","value":86},{"title":"Failure","value":1}]}}],"section":{"data":[{"title":"Success","value":86},{"title":"Failure","value":1}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getBuildsComponentComparison("buildsComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing builds component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetDeploymentsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of deployments component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"deploymentsChart":{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":6,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deployments_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":6,"deploys":{"value":70}},{"key":"cb4ffec3-1a8d-4c8a-9fc1-b4d9639b268c","doc_count":6,"deploys":{"value":60}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":6,"deploys":{"value":69}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":69,"value_in_millis":0,"compare_reports":null,"section":{"data":null}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":130,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":70,"value_in_millis":0,"compare_reports":null,"section":{"data":null}},{"is_sub_org":false,"sub_org_id":"cb4ffec3-1a8d-4c8a-9fc1-b4d9639b268c","compare_title":"component 12","sub_org_count":0,"component_count":1,"total_value":60,"value_in_millis":0,"compare_reports":null,"section":{"data":null}}],"section":{"data":null}}]`)
		organisation, err := getMultipleOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getDeploymentsComponentComparison("deploymentsComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing deployments component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetComponentsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of components component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"automationsChart":{"took":890,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1047,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"distinct_component":{"value":["f4180826-bb76-421e-5410-791408daadeb","238ffe68-8cb4-459d-64ac-2e4f752fe8dc","80a5fe62-4d9c-4385-b2ad-fb6ad99e9b2d","78ccc7ea-9c30-4a90-4aea-d3a07df8c382","ed35ceee-42b2-40b2-85fe-d3407746eb72","a90a60c5-b39d-438b-849c-6b333f765f38","0c3c691a-28a4-4779-afb5-da0f8ca4e6d3","4951e6cc-75c9-4324-a51c-51b8a3d0de5e","54d882eb-06ad-4528-9b7f-7a5db6b7bd7a"]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Active","value":1},{"title":"Inactive","value":0}]}},{"is_sub_org":false,"sub_org_id":"test-inactive-comp","compare_title":"test inactive component","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Active","value":0},{"title":"Inactive","value":1}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":2,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Active","value":1},{"title":"Inactive","value":0}]}},{"is_sub_org":false,"sub_org_id":"test-inactive-sub-org-comp","compare_title":"test inactive sub org component","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Active","value":0},{"title":"Inactive","value":1}]}}],"section":{"data":[{"title":"Active","value":1},{"title":"Inactive","value":1}]}}]`)
		organisation, err := getOrganisationInactive()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getComponentsComponentComparison("componentsComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing components component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetWorkflowsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of workflows component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"workflow component comparison data":{"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a":{"active":0,"inactive":12},"002bac07-6f60-43da-81c3-86d72f0610c4":{"active":0,"inactive":1},"f4180826-bb76-421e-5410-791408daadeb":{"active":0,"inactive":3},"013f8c84-fde9-4172-9a45-f33295ebd09d":{"active":0,"inactive":7}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":3,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Active","value":0},{"title":"Inactive","value":3}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":12,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":12,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Active","value":0},{"title":"Inactive","value":12}]}}],"section":{"data":[{"title":"Active","value":0},{"title":"Inactive","value":12}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getWorkflowsComponentComparison("workflowsComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing workflows component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetSecurityWorkflowRunsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of workflow runs component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"runsStatusChart":{"took":413,"timed_out":false,"_shards":{"total":4,"successful":4,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"workflow_runs_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":161961,"run_status":{"value":{"chartData":{"data":[{"name":"With scanners","value":3},{"name":"Without scanners","value":96}],"info":[{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"With scanners","value":35},{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"Without scanners","value":1068}]},"Total":{"value":1103,"key":"Total"}}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":64434,"run_status":{"value":{"chartData":{"data":[{"name":"With scanners","value":80},{"name":"Without scanners","value":19}],"info":[{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"With scanners","value":220},{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"Without scanners","value":54}]},"Total":{"value":274,"key":"Total"}}}},{"key":"bc8994ee-9bed-4a9f-bc4c-9fd32b97d326","doc_count":5,"run_status":{"value":{"chartData":{"data":[{"name":"With scanners","value":0},{"name":"Without scanners","value":100}],"info":[{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"With scanners","value":0},{"drillDown":{"reportType":"scannerType","reportId":"security-workflowRuns","reportTitle":"Successful workflow runs"},"title":"Without scanners","value":1}]},"Total":{"value":1,"key":"Total"}}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":1103,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With scanners","value":35},{"title":"Without scanners","value":1068}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":274,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":274,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With scanners","value":220},{"title":"Without scanners","value":54}]}}],"section":{"data":[{"title":"With scanners","value":220},{"title":"Without scanners","value":54}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getSecurityWorkflowRunsComponentComparison("workflowRunsComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing workflow runs component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetSecurityWorkflowsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of security workflows component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"security workflow component comparison data":{"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a":{"withScanners":0,"withoutScanners":12},"002bac07-6f60-43da-81c3-86d72f0610c4":{"withScanners":0,"withoutScanners":1},"f4180826-bb76-421e-5410-791408daadeb":{"withScanners":0,"withoutScanners":3},"013f8c84-fde9-4172-9a45-f33295ebd09d":{"withScanners":0,"withoutScanners":7}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":3,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With scanners","value":0},{"title":"Without scanners","value":3}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":12,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":12,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With scanners","value":0},{"title":"Without scanners","value":12}]}}],"section":{"data":[{"title":"With scanners","value":0},{"title":"Without scanners","value":12}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getSecurityWorkflowsComponentComparison("securityWorkflowsComponentSpec", x, replacements, organisation)

		assert.Nil(t, err, "error processing workflows component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetSecurityComponentsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of components component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"security widget component comparison":["54d882eb-06ad-4528-9b7f-7a5db6b7bd7a", "f4180826-bb76-421e-5410-791408daadeb", "test-bb76-421e-5410-791408daadeb"]}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With scanners","value":1},{"title":"Without scanners","value":0}]}},{"is_sub_org":false,"sub_org_id":"test-inactive-comp","compare_title":"test inactive component","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With scanners","value":0},{"title":"Without scanners","value":1}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":2,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With scanners","value":1},{"title":"Without scanners","value":0}]}},{"is_sub_org":false,"sub_org_id":"test-inactive-sub-org-comp","compare_title":"test inactive sub org component","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With scanners","value":0},{"title":"Without scanners","value":1}]}}],"section":{"data":[{"title":"With scanners","value":1},{"title":"Without scanners","value":1}]}}]`)
		organisation, err := getOrganisationInactive()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getSecurityComponentsComponentComparison("componentsComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing components component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_getMttrVeryHighComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of MTTR very high component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"mttrVeryHighChart":{"took":244,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":3243,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"mttr_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":661,"Avg_TTR":{"value":{"VERY_HIGH":25}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":608,"Avg_TTR":{"value":{"VERY_HIGH":13}}},{"key":"f4180826-bb76-421e-5412-791408daadeb","doc_count":608,"Avg_TTR":{"value":{"VERY_HIGH":15}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":13,"compare_reports":null},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":0,"value_in_millis":20,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":25,"compare_reports":null},{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5412-791408daadeb","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":15,"compare_reports":null}]}]`)
		organisation, err := getOrganisationMttr()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getMttrVeryHighComponentComparison("mttrVeryHighComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing mttr very high component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_getMttrHighComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of MTTR high component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"mttrHighChart":{"took":190,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":5107,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"mttr_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":2790,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":963,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"f4180826-bb76-421e-5412-791408daadeb","doc_count":344,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"40df6d5b-2401-4c1c-a868-1a5a4ea1bbbe","doc_count":176,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"545c6a93-7b98-49b4-707e-9c9ab52ce617","doc_count":163,"Avg_TTR":{"value":{"VERY_HIGH":17344000}}},{"key":"e2f8fef6-5041-4843-b37e-6cdae38099bc","doc_count":120,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"956130b8-761d-4033-9968-3185b09cb24e","doc_count":70,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"bf3322a8-ab44-4e9b-a75f-863ac11737f3","doc_count":60,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"92df4461-d43f-4d19-5e00-4b6346f747d9","doc_count":45,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"55ad025f-157a-4af5-98c7-bab36ba1b1b5","doc_count":44,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"cdcfa229-facc-4e9a-80a1-d82211f910b4","doc_count":39,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"94a81dd1-3f52-4520-891e-b2440f660945","doc_count":34,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"3cd7e6f3-8fcb-426f-791e-6285150c6994","doc_count":28,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"bc8994ee-9bed-4a9f-bc4c-9fd32b97d326","doc_count":27,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"6982e936-d26f-45bf-9155-9aeed85a5137","doc_count":26,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"b4b713e1-b51f-469b-6a2b-b02b074e7271","doc_count":24,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"f86500d8-2c88-43a2-bb18-0ef531664fd1","doc_count":20,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"d153c6f9-3809-4be9-b0db-113e8e4a0fe4","doc_count":19,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"e6c3a25b-e576-4d28-4ca5-d8c4af3d0113","doc_count":17,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"dda69191-5492-4b7e-88b2-9d9d42f61899","doc_count":14,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"0af5ff0c-c697-4576-a717-0319264e231b","doc_count":12,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"31e661ac-6553-4d6c-6f3c-6528f64b8fcd","doc_count":12,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"539f5833-95f0-453a-821a-6b4516fc93b0","doc_count":12,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"addf84e3-5c30-49d0-aaa7-1564ee109c0b","doc_count":9,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"f783334c-97b3-478e-8788-2d8e76c91d91","doc_count":9,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"709da9f4-c61d-4443-9cdd-8dce20f65052","doc_count":8,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"087cae5a-db8f-4d34-bef2-1175c7e0d394","doc_count":4,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"200d9c02-6e45-41ee-8ad4-8738810ec69f","doc_count":4,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"4c1ce669-7acc-475f-b5ee-b5c826ff5c3c","doc_count":4,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"9a2e490b-63d1-4dfd-830f-c26a8eccfcb0","doc_count":4,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"d000292d-bd9c-44c9-a3f5-598ec3a561f4","doc_count":3,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"45abf2e8-65c8-4787-adef-43c90144edd6","doc_count":2,"Avg_TTR":{"value":{"VERY_HIGH":0}}},{"key":"85b1f8b6-421d-403d-9bb1-ebbd05732958","doc_count":1,"Avg_TTR":{"value":{"VERY_HIGH":0}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":0,"compare_reports":null},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":0,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":0,"compare_reports":null},{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5412-791408daadeb","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":0,"compare_reports":null}]}]`)
		organisation, err := getOrganisationMttr()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getMttrHighComponentComparison("mttrHighComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing mttr high component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_getMttrMediumComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of MTTR medium component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"mttrMediumChart":{"took":244,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":3243,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"mttr_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":661,"Avg_TTR":{"value":{"MEDIUM":25}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":608,"Avg_TTR":{"value":{"MEDIUM":13}}},{"key":"f4180826-bb76-421e-5412-791408daadeb","doc_count":608,"Avg_TTR":{"value":{"MEDIUM":15}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":13,"compare_reports":null},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":0,"value_in_millis":20,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":25,"compare_reports":null},{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5412-791408daadeb","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":15,"compare_reports":null}]}]`)
		organisation, err := getOrganisationMttr()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getMttrMediumComponentComparison("mttrMediumComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing mttr medium component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_getMttrLowComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of MTTR low component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"mttrLowChart":{"took":244,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":3243,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"mttr_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":661,"Avg_TTR":{"value":{"LOW":20}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":608,"Avg_TTR":{"value":{"LOW":13}}},{"key":"f4180826-bb76-421e-5412-791408daadeb","doc_count":608,"Avg_TTR":{"value":{"LOW":25}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":13,"compare_reports":null},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":0,"value_in_millis":23,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":20,"compare_reports":null},{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5412-791408daadeb","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":25,"compare_reports":null}]}]`)
		organisation, err := getOrganisationMttr()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getMttrLowComponentComparison("mttrLowComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing mttr low component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func getOrganisationMttr() (*constants.Organization, error) {
	jsonData := `
	{
	    "id": "Org-1",
	    "name": "Org-1",
	    "sub_orgs":
	    [
	        {
	            "id": "sub-org-1",
	            "name": "sub-org-1",
	            "sub_orgs":
	            [],
	            "components":
	            [
	                {
	                    "id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
	                    "name": "component 11"
	                },
					{
	                    "id": "f4180826-bb76-421e-5412-791408daadeb",
	                    "name": "component 11"
	                }
	            ]
	        }
	    ],
	    "components":
	    [
	        {
	            "id": "f4180826-bb76-421e-5410-791408daadeb",
	            "name": "component 1"
	        }
	    ]
	}`

	// Unmarshal JSON into Organization struct
	var orgData constants.Organization
	if err := json.Unmarshal([]byte(jsonData), &orgData); err != nil {
		log.Errorf(err, "Error parsing JSON : ")
		return nil, err
	}
	return &orgData, nil
}

func Test_GetSastVulnerabilitiesComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of sast vulnerabilities scanner component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"sastVulnerabilityScannerChart":{"took":414,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":0,"hits":[{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1192_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S100_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S3776_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S107_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1871_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1479_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S108_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1135_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_docker:S6471_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_b369a51f-f220-421d-8787-359fb2617193_CODE_SMELL-go:S1192_1710652490840","_score":0}]},"aggregations":{"vul_by_scanner_type_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":2554,"vulByScannerTypeCounts":{"value":{"VERY_HIGH":[{"x":"SAST","y":9},{"x":"DAST","y":0},{"x":"Container","y":0},{"x":"SCA","y":0}],"HIGH":[{"x":"SAST","y":42},{"x":"DAST","y":0},{"x":"Container","y":0},{"x":"SCA","y":0}],"MEDIUM":[{"x":"SAST","y":18},{"x":"DAST","y":0},{"x":"Container","y":0},{"x":"SCA","y":0}],"LOW":[{"x":"SAST","y":1},{"x":"DAST","y":0},{"x":"Container","y":0},{"x":"SCA","y":0}]}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":21022,"vulByScannerTypeCounts":{"value":{"VERY_HIGH":[{"x":"SAST","y":3},{"x":"DAST","y":0},{"x":"Container","y":0},{"x":"SCA","y":0}],"HIGH":[{"x":"SAST","y":4},{"x":"DAST","y":0},{"x":"Container","y":0},{"x":"SCA","y":0}],"MEDIUM":[{"x":"SAST","y":3},{"x":"DAST","y":0},{"x":"Container","y":1},{"x":"SCA","y":0}],"LOW":[{"x":"SAST","y":1},{"x":"DAST","y":0},{"x":"Container","y":84},{"x":"SCA","y":0}]}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":70,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Very high","value":9},{"title":"High","value":42},{"title":"Medium","value":18},{"title":"Low","value":1}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":11,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":11,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Very high","value":3},{"title":"High","value":4},{"title":"Medium","value":3},{"title":"Low","value":1}]}}],"section":{"data":[{"title":"Very high","value":3},{"title":"High","value":4},{"title":"Medium","value":3},{"title":"Low","value":1}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getSastVulnerabilitiesComponentComparison("sastVulnerabilitiesComponentSpec", x, replacements, organisation)

		assert.Nil(t, err, "error processing sast vulnerabilities scanner component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetDastVulnerabilitiesComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of dast vulnerabilities scanner component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"dastVulnerabilityScannerChart":{"took":414,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":0,"hits":[{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1192_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S100_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S3776_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S107_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1871_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1479_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S108_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1135_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_docker:S6471_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_b369a51f-f220-421d-8787-359fb2617193_CODE_SMELL-go:S1192_1710652490840","_score":0}]},"aggregations":{"vul_by_scanner_type_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":2554,"vulByScannerTypeCounts":{"value":{"VERY_HIGH":[{"x":"SAST","y":9},{"x":"DAST","y":5},{"x":"Container","y":7},{"x":"SCA","y":8}],"HIGH":[{"x":"SAST","y":42},{"x":"DAST","y":6},{"x":"Container","y":2},{"x":"SCA","y":5}],"MEDIUM":[{"x":"SAST","y":18},{"x":"DAST","y":6},{"x":"Container","y":3},{"x":"SCA","y":7}],"LOW":[{"x":"SAST","y":1},{"x":"DAST","y":4},{"x":"Container","y":6},{"x":"SCA","y":8}]}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":21022,"vulByScannerTypeCounts":{"value":{"VERY_HIGH":[{"x":"SAST","y":3},{"x":"DAST","y":6},{"x":"Container","y":7},{"x":"SCA","y":3}],"HIGH":[{"x":"SAST","y":4},{"x":"DAST","y":2},{"x":"Container","y":1},{"x":"SCA","y":5}],"MEDIUM":[{"x":"SAST","y":3},{"x":"DAST","y":7},{"x":"Container","y":1},{"x":"SCA","y":4}],"LOW":[{"x":"SAST","y":1},{"x":"DAST","y":7},{"x":"Container","y":84},{"x":"SCA","y":6}]}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":21,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Very high","value":5},{"title":"High","value":6},{"title":"Medium","value":6},{"title":"Low","value":4}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":22,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":22,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Very high","value":6},{"title":"High","value":2},{"title":"Medium","value":7},{"title":"Low","value":7}]}}],"section":{"data":[{"title":"Very high","value":6},{"title":"High","value":2},{"title":"Medium","value":7},{"title":"Low","value":7}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getDastVulnerabilitiesComponentComparison("dastVulnerabilitiesComponentSpec", x, replacements, organisation)

		assert.Nil(t, err, "error processing dast vulnerabilities scanner component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetContainerVulnerabilitiesComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of conatiner vulnerabilities scanner component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"containerVulnerabilityScannerChart":{"took":414,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":0,"hits":[{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1192_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S100_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S3776_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S107_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1871_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1479_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S108_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1135_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_docker:S6471_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_b369a51f-f220-421d-8787-359fb2617193_CODE_SMELL-go:S1192_1710652490840","_score":0}]},"aggregations":{"vul_by_scanner_type_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":2554,"vulByScannerTypeCounts":{"value":{"VERY_HIGH":[{"x":"SAST","y":9},{"x":"DAST","y":5},{"x":"Container","y":7},{"x":"SCA","y":8}],"HIGH":[{"x":"SAST","y":42},{"x":"DAST","y":6},{"x":"Container","y":2},{"x":"SCA","y":5}],"MEDIUM":[{"x":"SAST","y":18},{"x":"DAST","y":6},{"x":"Container","y":3},{"x":"SCA","y":7}],"LOW":[{"x":"SAST","y":1},{"x":"DAST","y":4},{"x":"Container","y":6},{"x":"SCA","y":8}]}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":21022,"vulByScannerTypeCounts":{"value":{"VERY_HIGH":[{"x":"SAST","y":3},{"x":"DAST","y":6},{"x":"Container","y":7},{"x":"SCA","y":3}],"HIGH":[{"x":"SAST","y":4},{"x":"DAST","y":2},{"x":"Container","y":1},{"x":"SCA","y":5}],"MEDIUM":[{"x":"SAST","y":3},{"x":"DAST","y":7},{"x":"Container","y":1},{"x":"SCA","y":4}],"LOW":[{"x":"SAST","y":1},{"x":"DAST","y":7},{"x":"Container","y":84},{"x":"SCA","y":6}]}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":18,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Very high","value":7},{"title":"High","value":2},{"title":"Medium","value":3},{"title":"Low","value":6}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":93,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":93,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Very high","value":7},{"title":"High","value":1},{"title":"Medium","value":1},{"title":"Low","value":84}]}}],"section":{"data":[{"title":"Very high","value":7},{"title":"High","value":1},{"title":"Medium","value":1},{"title":"Low","value":84}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getContainerVulnerabilitiesComponentComparison("containerVulnerabilitiesComponentSpec", x, replacements, organisation)

		assert.Nil(t, err, "error processing conatiner vulnerabilities scanner component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetScaVulnerabilitiesComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of sca vulnerabilities scanner component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"scaVulnerabilityScannerChart":{"took":414,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":0,"hits":[{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1192_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S100_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S3776_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S107_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1871_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1479_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S108_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_CODE_SMELL-go:S1135_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_979cd4e3-4ab1-462a-a720-cd73a2ef6a99_docker:S6471_1710654465310","_score":0},{"_index":"scan_results","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_b369a51f-f220-421d-8787-359fb2617193_CODE_SMELL-go:S1192_1710652490840","_score":0}]},"aggregations":{"vul_by_scanner_type_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":2554,"vulByScannerTypeCounts":{"value":{"VERY_HIGH":[{"x":"SAST","y":9},{"x":"DAST","y":5},{"x":"Container","y":7},{"x":"SCA","y":8}],"HIGH":[{"x":"SAST","y":42},{"x":"DAST","y":6},{"x":"Container","y":2},{"x":"SCA","y":5}],"MEDIUM":[{"x":"SAST","y":18},{"x":"DAST","y":6},{"x":"Container","y":3},{"x":"SCA","y":7}],"LOW":[{"x":"SAST","y":1},{"x":"DAST","y":4},{"x":"Container","y":6},{"x":"SCA","y":8}]}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":21022,"vulByScannerTypeCounts":{"value":{"VERY_HIGH":[{"x":"SAST","y":3},{"x":"DAST","y":6},{"x":"Container","y":7},{"x":"SCA","y":3}],"HIGH":[{"x":"SAST","y":4},{"x":"DAST","y":2},{"x":"Container","y":1},{"x":"SCA","y":5}],"MEDIUM":[{"x":"SAST","y":3},{"x":"DAST","y":7},{"x":"Container","y":1},{"x":"SCA","y":4}],"LOW":[{"x":"SAST","y":1},{"x":"DAST","y":7},{"x":"Container","y":84},{"x":"SCA","y":6}]}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":28,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Very high","value":8},{"title":"High","value":5},{"title":"Medium","value":7},{"title":"Low","value":8}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":18,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":18,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"Very high","value":3},{"title":"High","value":5},{"title":"Medium","value":4},{"title":"Low","value":6}]}}],"section":{"data":[{"title":"Very high","value":3},{"title":"High","value":5},{"title":"Medium","value":4},{"title":"Low","value":6}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getScaVulnerabilitiesComponentComparison("scaVulnerabilitiesComponentSpec", x, replacements, organisation)

		assert.Nil(t, err, "error processing sca vulnerabilities scanner component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_getDoraMttrComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of Dora MTTR component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"mttrHeader":{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"mttr_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":1,"deployments":{"value":{"recoveredTotalDuration":100,"recoveredCount":3}}},{"key":"f4180826-bb76-421e-5412-791408daadeb","doc_count":1,"deployments":{"value":{"recoveredTotalDuration":200,"recoveredCount":5}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":1,"deployments":{"value":{"recoveredTotalDuration":300,"recoveredCount":7}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":43,"compare_reports":null},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":0,"value_in_millis":38,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":33,"compare_reports":null},{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5412-791408daadeb","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":40,"compare_reports":null}]}]`)
		organisation, err := getOrganisationMttr()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getDoraMttrComponentComparison("mttrComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing Dora mttr component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_getDeploymentLeadTimeComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of Deployment Lead Time component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"deploymentLeadTimeHeader":{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"deployment_lead_time_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":1,"deploy_data":{"value":{"totalDuration":57000,"average":57000,"deployments":1}}},{"key":"f4180826-bb76-421e-5412-791408daadeb","doc_count":1,"deploy_data":{"value":{"totalDuration":64000,"average":32000,"deployments":2}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":1,"deploy_data":{"value":{"totalDuration":98000,"average":32666.67,"deployments":3}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":32667,"compare_reports":null},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":0,"value_in_millis":40333,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":57000,"compare_reports":null},{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5412-791408daadeb","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":0,"value_in_millis":32000,"compare_reports":null}]}]`)

		organisation, err := getOrganisationMttr()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getDeploymentLeadTimeComponentComparison("deploymentLeadTimeSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing Deployment Lead Time component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_getFailureRateComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of failure rate component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"averageFailureRateHeader":{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"failure_rate_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":1,"deploy_data":{"value":{"average":"30.0%","deployments":100,"failedDeployments":30}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":1,"deploy_data":{"value":{"average":"0.0%","deployments":90,"failedDeployments":20}}},{"key":"f4180826-bb76-421e-5412-791408daadeb","doc_count":1,"deploy_data":{"value":{"average":"0.0%","deployments":75,"failedDeployments":14}}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":30,"value_in_millis":0,"compare_reports":null},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":21,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":22,"value_in_millis":0,"compare_reports":null},{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5412-791408daadeb","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":19,"value_in_millis":0,"compare_reports":null}]}]`)
		organisation, err := getOrganisationMttr()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getFailureRateComponentComparison("failureRateComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing failure rate component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_transformSuccessfulBuildDuration(t *testing.T) {

	type SuccessfulBuildDurationOutput struct {
		Data []struct {
			X string `json:"x"`
			Y []int  `json:"y"`
		} `json:"data"`
		ID  string `json:"id"`
		Min string `json:"min"`
		Max string `json:"max"`
	}

	// constructing test input in the format the function expects
	inputData := json.RawMessage(`{
		"_shards": {
		  "failed": 0,
		  "skipped": 0,
		  "successful": 2,
		  "total": 2
		},
		"aggregations": {
		  "builds": {
			"value": {
			  "s": [
				26000,
				26000,
				27000,
				48000,
				48000
			  ]
			}
		  }
		},
		"hits": {
		  "hits": [],
		  "max_score": null,
		  "total": {
			"relation": "gte",
			"value": 10000
		  }
		},
		"status": 200,
		"timed_out": false,
		"took": 671
	  }`)
	testInputMap := make(map[string]json.RawMessage)
	testInputMap["successfulBuildDuration"] = inputData

	inputData2 := json.RawMessage(`{
		"_shards": {
		  "failed": 0,
		  "skipped": 0,
		  "successful": 2,
		  "total": 2
		},
		"aggregations": {
		  "builds": {
			"value": {}
		  }
		},
		"hits": {
		  "hits": [],
		  "max_score": null,
		  "total": {
			"relation": "gte",
			"value": 10000
		  }
		},
		"status": 200,
		"timed_out": false,
		"took": 671
	  }`)
	testInputMap2 := make(map[string]json.RawMessage)
	testInputMap2["successfulBuildDuration"] = inputData2

	// expected output from the function
	expected := json.RawMessage(`[{"data":[{"x":"s","y":[26000,26000,27000,48000,48000]}],"id":"Build Duration","max":"6000","min":"0"}]`)
	expected2 := json.RawMessage(`[{"data":[],"id":"Build Duration","max":"6000","min":"0"}]`)

	//creating structs to unmarshall the expected and actual json.RawMessage data into, so that they can be checked for equality, without taking the order of the JSON fields at the outermost level into account
	var assertionActual []SuccessfulBuildDurationOutput

	var assertionExpected []SuccessfulBuildDurationOutput
	err := json.Unmarshal(expected, &assertionExpected)
	if err != nil {
		t.Errorf("error unmarshalling expected output to struct in Test_transformSuccessfulBuildDuration(): %s", err)
		return
	}

	var assertionExpected2 []SuccessfulBuildDurationOutput
	err = json.Unmarshal(expected2, &assertionExpected2)
	if err != nil {
		t.Errorf("error unmarshalling expected output to struct in Test_transformSuccessfulBuildDuration(): %s", err)
		return
	}

	type args struct {
		specKey      string
		data         map[string]json.RawMessage
		replacements map[string]any
	}

	tests := []struct {
		name    string
		args    args
		want    []SuccessfulBuildDurationOutput
		wantErr bool
	}{
		// test cases
		{
			name: "test case 1: succesful transformation case",
			args: args{
				specKey: "successfulBuildDurationSpec",
				data:    testInputMap,
				replacements: map[string]any{
					"min": "0",
					"max": "6000",
				},
			},
			want:    assertionExpected,
			wantErr: false,
		},
		{
			name: "test case 2: transformation when the value field in the input is empty",
			args: args{
				specKey: "successfulBuildDurationSpec",
				data:    testInputMap2,
				replacements: map[string]any{
					"min": "0",
					"max": "6000",
				},
			},
			want:    assertionExpected2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformSuccessfulBuildDuration(tt.args.specKey, tt.args.data, tt.args.replacements)
			if (err != nil) != tt.wantErr {
				t.Errorf("transformSuccessfulBuildDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// unmarshall actual output into a struct for assertion
			err = json.Unmarshal(got, &assertionActual)
			if err != nil {
				t.Errorf("error unmarshalling actual output to struct in Test_transformSuccessfulBuildDuration(): %s", err)
				return
			}

			if !assert.Equal(t, tt.want, assertionActual) {
				t.Errorf("transformSuccessfulBuildDuration() = %v, want %v", assertionActual, tt.want)
			}
		})
	}
}

func Test_transformSummaryLatestTestResultsSection(t *testing.T) {

	// Case 1: constructing test input in the format the function expects
	inputData := json.RawMessage(`{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":58,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"workflows":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"33","doc_count":58,"suites":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"","doc_count":10,"latest_suite_doc":{"hits":{"total":{"value":10,"relation":"eq"},"max_score":null,"hits":[{"_index":"tt","_id":"gg","_score":null,"_source":{"duration":0,"start_time":"2024-06-24 07:12:41","total":13,"run_id":"rr","passed":8,"failed":4,"automation_name":"workflow","skipped":1,"source":"CloudBees"},"fields":{"zoned_start_time":["2024/06/24 12:42:41"],"start_time_in_millis":[1719213161000]},"sort":[1719213161000]}]}}},{"key":"dd","doc_count":10,"latest_suite_doc":{"hits":{"total":{"value":10,"relation":"eq"},"max_score":null,"hits":[{"_index":"cc","_id":"aa","_score":null,"_source":{"duration":600,"start_time":"2024-06-24 07:12:41","total":11,"run_id":"r2","passed":10,"failed":0,"automation_name":"workflow","skipped":1,"source":"CloudBees"},"fields":{"zoned_start_time":["2024/06/24 12:42:41"],"start_time_in_millis":[1719213161000]},"sort":[1719213161000]}]}}},{"key":"k3","doc_count":10,"latest_suite_doc":{"hits":{"total":{"value":10,"relation":"eq"},"max_score":null,"hits":[{"_index":"cb_test_suites","_id":"vv","_score":null,"_source":{"duration":400,"start_time":"2024-06-24 07:12:41","total":7,"run_id":"r2d2","passed":4,"failed":2,"automation_name":"workflow","skipped":1,"source":"CloudBees"},"fields":{"zoned_start_time":["2024/06/24 12:42:41"],"start_time_in_millis":[1719213161000]},"sort":[1719213161000]}]}}},{"key":"ss","doc_count":10,"latest_suite_doc":{"hits":{"total":{"value":10,"relation":"eq"},"max_score":null,"hits":[{"_index":"cc","_id":"s3","_score":null,"_source":{"duration":400,"start_time":"2024-06-24 07:12:41","total":6,"run_id":"r3","passed":4,"failed":2,"automation_name":"workflow","skipped":0,"source":"CloudBees"},"fields":{"zoned_start_time":["2024/06/24 12:42:41"],"start_time_in_millis":[1719213161000]},"sort":[1719213161000]}]}}},{"key":"k3","doc_count":9,"latest_suite_doc":{"hits":{"total":{"value":9,"relation":"eq"},"max_score":null,"hits":[{"_index":"cc","_id":"fg","_score":null,"_source":{"duration":400,"start_time":"2024-06-24 07:12:41","total":7,"run_id":"rr","passed":4,"failed":2,"automation_name":"workflow","skipped":1,"source":"CloudBees"},"fields":{"zoned_start_time":["2024/06/24 12:42:41"],"start_time_in_millis":[1719213161000]},"sort":[1719213161000]}]}}},{"key":"gg","doc_count":9,"latest_suite_doc":{"hits":{"total":{"value":9,"relation":"eq"},"max_score":null,"hits":[{"_index":"cc","_id":"ss","_score":null,"_source":{"duration":400,"start_time":"2024-06-24 07:12:41","total":6,"run_id":"rr","passed":4,"failed":2,"automation_name":"workflow","skipped":0,"source":"CloudBees"},"fields":{"zoned_start_time":["2024/06/24 12:42:41"],"start_time_in_millis":[1719213161000]},"sort":[1719213161000]}]}}}]}}]}}}`)
	testInputMap := make(map[string]json.RawMessage)
	testInputMap["summaryLatestTestResultsSection"] = inputData

	// expected output from the function
	expected := `[{"testSuiteName":"","lastRun":"1970/01/01 05:30:00","lastRunInMillis":0,"totalTestCases":24,"testCasesPassed":18,"testCasesFailed":4,"testCasesSkipped":2,"runTime":"0s","source":"CloudBees","drillDown":{"reportId":"latest-test-results","reportTitle":"Test cases","reportType":"status","reportInfo":{"run_id":"00f2eb27-ab2b-43b7-8f0c-aea4de761526"}}},{"testSuiteName":"github.com/calculi-corp/template-go-testing/test_data","lastRun":"1970/01/01 05:30:00","lastRunInMillis":0,"totalTestCases":11,"testCasesPassed":10,"testCasesFailed":0,"testCasesSkipped":1,"runTime":"0.6s","source":"CloudBees","drillDown":{"reportId":"latest-test-results","reportTitle":"Test cases","reportType":"status","reportInfo":{"run_id":"00f2eb27-ab2b-43b7-8f0c-aea4de761526","test_suite_name":"github.com/calculi-corp/template-go-testing/test_data"}}},{"testSuiteName":"github.com/calculi-corp/template-go-testing/test_suite_1","lastRun":"1970/01/01 05:30:00","lastRunInMillis":0,"totalTestCases":7,"testCasesPassed":4,"testCasesFailed":2,"testCasesSkipped":1,"runTime":"0.4s","source":"CloudBees","drillDown":{"reportId":"latest-test-results","reportTitle":"Test cases","reportType":"status","reportInfo":{"run_id":"00f2eb27-ab2b-43b7-8f0c-aea4de761526","test_suite_name":"github.com/calculi-corp/template-go-testing/test_suite_1"}}},{"testSuiteName":"github.com/calculi-corp/template-go-testing/test_suite_2","lastRun":"1970/01/01 05:30:00","lastRunInMillis":0,"totalTestCases":6,"testCasesPassed":4,"testCasesFailed":2,"testCasesSkipped":0,"runTime":"0.4s","source":"CloudBees","drillDown":{"reportId":"latest-test-results","reportTitle":"Test cases","reportType":"status","reportInfo":{"run_id":"00f2eb27-ab2b-43b7-8f0c-aea4de761526","test_suite_name":"github.com/calculi-corp/template-go-testing/test_suite_2"}}}]`

	// Case 2: constructing test input in the format the function expects
	testInputMap2 := make(map[string]json.RawMessage)
	testInputMap2["summaryLatestTestResultsSectio"] = inputData

	// Case 3: constructing test input in the format the function expects
	inputData3 := json.RawMessage(`{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":9,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"workflows":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"ff","doc_count":4,"suites":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"kk","doc_count":4,"latest_suite_doc":{"hits":{"total":{"value":4,"relation":"eq"},"max_score":null,"hits":[{"_index":"cc","_id":"gg","_score":null,"_source":{"duration":15,"start_time":"1970-01-20 21:46:44","total":8,"run_id":"rr","passed":8,"failed":0,"automation_name":"workflow","skipped":0,"source":"CloudBees"},"fields":{"zoned_start_time":["1970/01/21 03:16:44"],"start_time_in_millis":[1720004000]},"sort":[1720004158000]}]}}}]}},{"key":"tt","doc_count":3,"suites":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"kk","doc_count":3,"latest_suite_doc":{"hits":{"total":{"value":3,"relation":"eq"},"max_score":null,"hits":[{"_index":"cc","_id":"hh","_score":null,"_source":{"duration":14,"start_time":"1970-01-20 14:17:11","total":8,"run_id":"rr","passed":8,"failed":0,"automation_name":"workflow1","skipped":0},"fields":{"zoned_start_time":["1970/01/20 19:47:11"],"start_time_in_millis":[1693031000]},"sort":[1720004110000]}]}}}]}},{"key":"k3","doc_count":2,"suites":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"ww","doc_count":2,"latest_suite_doc":{"hits":{"total":{"value":2,"relation":"eq"},"max_score":null,"hits":[{"_index":"cc","_id":"gg","_score":null,"_source":{"duration":14,"start_time":"1970-01-20 14:17:11","total":8,"run_id":"rr","passed":8,"failed":0,"automation_name":"workflow2","skipped":0},"fields":{"zoned_start_time":["1970/01/20 19:47:11"],"start_time_in_millis":[1693031000]},"sort":[1720004033000]}]}}}]}}]}}}`)
	testInputMap3 := make(map[string]json.RawMessage)
	testInputMap3["summaryLatestTestResultsSection"] = inputData3

	type args struct {
		specKey      string
		data         map[string]json.RawMessage
		replacements map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// test cases
		{
			name: "test case 1: succesful transformation case",
			args: args{
				specKey:      "summaryLatestTestResultsSectionSpec",
				data:         testInputMap,
				replacements: map[string]any{},
			},
			want:    2755,
			wantErr: false,
		},
		{
			name: "test case 2: internal server error when query key is wrong in data",
			args: args{
				specKey:      "summaryLatestTestResultsSectionSpec",
				data:         testInputMap2,
				replacements: map[string]any{},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "test case 3: succesful transformation case for multiple workflows in the same branch",
			args: args{
				specKey:      "summaryLatestTestResultsSectionSpec",
				data:         testInputMap3,
				replacements: map[string]any{},
			},
			want:    1382,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformSummaryLatestTestResultsSection(tt.args.specKey, tt.args.data, tt.args.replacements)
			if (err != nil) != tt.wantErr {
				t.Errorf("transformSummaryLatestTestResultsSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, len(got)) {
				t.Errorf("transformSummaryLatestTestResultsSection() = %v, want %v", string(got), expected)
			}

		})
	}
}

func Test_GetAutomationRunsForTestSuites(t *testing.T) {
	t.Run("Case 1: Successful execution of Automation runs for test suites", func(t *testing.T) {
		replacements := map[string]any{
			"startDate": "2024-06-01",
			"endDate":   "2024-06-30",
			"orgId":     "8509888e-d27f-44fa-46a9-29bc76f5e790",
		}

		responseString := `{"automationRuns":{"took":56,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":7978,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"component_activity":{"value":{"runs":6402}}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":6402}`)
		b, err := getAutomationRunsForTestSuites("totalTestRunsSpec", x, replacements)

		assert.Nil(t, err, "error processing GetAutomationRunsForTestSuites  header")
		assert.Equal(t, expectResult, []byte(b))

		responseString = `{"testSuites":{"took":15,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":707,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"component_activity":{"value":{"runs":182}}}},"automationRuns":{"took":56,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":7978,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"component_activity":{"value":{"runs":6402}}}}}`
		x = map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult = []byte(`{"data":[{"name":"With test suites","value":3},{"name":"Without test suites","value":97}],"info":[{"drillDown":{"reportType":"testSuiteType","reportId":"test-insights-workflowRuns","reportTitle":"Workflow runs"},"title":"With test suites","value":182},{"drillDown":{"reportType":"testSuiteType","reportId":"test-insights-workflowRuns","reportTitle":"Workflow runs"},"title":"Without test suites","value":6220}]}`)
		b, err = getAutomationRunsForTestSuites("testRunsChartSpec", x, replacements)
		assert.Nil(t, err, "error processing GetAutomationRunsForTestSuites section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_TransformTestWorkflowsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of test workflows component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"test workflow component comparison data":{"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a":{"withTestSuites":0,"withoutTestSuites":12},"002bac07-6f60-43da-81c3-86d72f0610c4":{"withTestSuites":0,"withoutTestSuites":1},"f4180826-bb76-421e-5410-791408daadeb":{"withTestSuites":0,"withoutTestSuites":3},"013f8c84-fde9-4172-9a45-f33295ebd09d":{"withTestSuites":0,"withoutTestSuites":7}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":3,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With test suites","value":0},{"title":"Without test suites","value":3}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":1,"total_value":12,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":12,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With test suites","value":0},{"title":"Without test suites","value":12}]}}],"section":{"data":[{"title":"With test suites","value":0},{"title":"Without test suites","value":12}]}}]`)
		organisation, err := getOrganisation()
		assert.Nil(t, err, "error fetching organisations")
		b, err := transformTestWorkflowsComponentComparison("", x, replacements, organisation)

		assert.Nil(t, err, "error processing test workflows component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_GetTestSuiteWorkflowRunsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of test suite workflow runs component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		data := map[string]json.RawMessage{
			"testSuiteRuns":  json.RawMessage(`{"took":10,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":724,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"workflow_runs_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":562,"component_activity":{"value":{"runs":168}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":162,"component_activity":{"value":{"runs":25}}}]}}}`),
			"automationRuns": json.RawMessage(`{"took":75,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":9357,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"workflow_runs_component_comparison":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"2b4d0070-09b7-40e1-8450-687094218174","doc_count":2525,"run_status":{"value":{"runs":1924}}},{"key":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","doc_count":720,"run_status":{"value":{"runs":200}}},{"key":"f9857f58-bb8d-48bf-8a33-3bfa480c791a","doc_count":676,"run_status":{"value":{"runs":578}}},{"key":"62b6124f-4ba6-44d4-a83a-2dbd9aae97ce","doc_count":471,"run_status":{"value":{"runs":307}}},{"key":"f4180826-bb76-421e-5410-791408daadeb","doc_count":445,"run_status":{"value":{"runs":300}}}]}}}`),
		}

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 2","sub_org_count":0,"component_count":1,"total_value":300,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With test suites","value":25},{"title":"Without test suites","value":275}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":1,"component_count":0,"total_value":200,"value_in_millis":0,"compare_reports":[{"is_sub_org":true,"sub_org_id":"sub-org-2","compare_title":"sub-org-2","sub_org_count":1,"component_count":0,"total_value":200,"value_in_millis":0,"compare_reports":[{"is_sub_org":true,"sub_org_id":"sub-org-3","compare_title":"sub-org-3","sub_org_count":0,"component_count":1,"total_value":200,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 3","sub_org_count":0,"component_count":1,"total_value":200,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With test suites","value":168},{"title":"Without test suites","value":32}]}}],"section":{"data":[{"title":"With test suites","value":168},{"title":"Without test suites","value":32}]}}],"section":{"data":[{"title":"With test suites","value":168},{"title":"Without test suites","value":32}]}}],"section":{"data":[{"title":"With test suites","value":168},{"title":"Without test suites","value":32}]}}]`)
		organisation, err := getTestInsightsMultiLevelSubOrg()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getTestSuiteWorkflowRunsComponentComparison("workflowRunsComponentSpec", data, replacements, organisation)
		assert.Nil(t, err, "error processing test suite workflow runs component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func getTestInsightsMultiLevelSubOrg() (*constants.Organization, error) {
	jsonData := `
	{
	    "id": "Org-1",
	    "name": "Org-1",
	    "sub_orgs":
	    [
	        {
	            "id": "sub-org-1",
	            "name": "sub-org-1",
	            "sub_orgs":
	            [{
					"id": "sub-org-2",
					"name": "sub-org-2",
					"sub_orgs":
					[{
						"id": "sub-org-3",
						"name": "sub-org-3",
						"sub_orgs":
						[],
						"components":
						[
							{
								"id": "020c6a27-5680-4d21-b1f2-3fac1a15053e",
								"name": "component 1"
							},
							{
								"id": "e389272a-3ad0-4766-bed8-1eff2211ed70",
								"name": "component 2"
							},
							{
								"id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
								"name": "component 3"
							}
						]
					}],
					"components":
					[
					]
				}],
	            "components":
	            [
	            ]
	        }
	    ],
	    "components":
	    [{
			"id": "f4180826-bb76-421e-5410-791408daadeb",
			"name": "component 2"
		}
	    ]
	}`

	// Unmarshal JSON into Organization struct
	var orgData constants.Organization
	if err := json.Unmarshal([]byte(jsonData), &orgData); err != nil {
		log.Errorf(err, "Error parsing JSON : ")
		return nil, err
	}
	return &orgData, nil
}

func Test_GetTestComponentsComponentComparison(t *testing.T) {
	t.Run("Case 1: Successful execution of test components component comparison", func(t *testing.T) {
		replacements := map[string]any{
			"startDate":            "2024-03-01",
			"endDate":              "2024-03-31",
			"orgId":                "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
			"aggrBy":               "week",
			"duration":             "month",
			"dateHistogramMin":     "2024-03-01",
			"dateHistogramMax":     "2024-03-31",
			"normalizeMonthInSpec": "2024-03-01",
			"userId":               "testUserId",
		}

		responseString := `{"components":{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":2436,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"components":{"value":["24c38dba-ad04-4847-9462-e7a3b13b48f6","73007a00-754e-4b8e-acc9-392a55f5bb2c","fa922acb-1aab-4cde-a1d7-620826895b5e","be41b4ea-d51f-418f-8da2-c35ea75e771f","23a56310-fb50-42cc-ad0d-4909b9639911","9a50fdd9-1af6-4732-b520-fa7e15f82d14","60d545a5-bd2f-447e-ab1d-a71b805501ff","62b6124f-4ba6-44d4-a83a-2dbd9aae97ce","afe0bc34-5113-408d-a8b7-96c00bfec0a1","3333344c-ae0f-4df4-b1a7-efcaacdf449e","ddbe564e-3107-4433-880a-fe88e1798b2a","a9814ddd-f367-4f0b-95ad-310e2dc2123d","8e2d5d2f-99fc-4933-a9ba-81c726324958","b5e7b966-36eb-43f4-ae29-1f82de323bf1","484d5e12-6424-4070-a159-4e5639a807a2","eec43295-207f-4eaf-96ac-864db047eefc","af2d841d-16bb-4fda-bc0f-4d1009feb95b","b4f66638-5075-474b-ab29-901c68f55d9b","2281041f-7381-491b-aec0-ca2e5cd9061d","67aceeaf-605d-4dfc-8802-32235e60c301","28ac919e-1701-4406-9523-759fd033dde1","f4e2d078-0f7b-4a60-9860-36934d7bf009","996c5c4d-b86f-4681-b492-6e4a29b02647","eac5f708-8c7f-4ac4-966d-1fb22730ffa7","3911ded3-6c6f-4a67-96a7-c54ec93fe12b","6ee2c296-5bc2-423a-8bef-2c2cdb344635","52f52762-90d2-44de-a58f-98b4afb2b6ea","1e3f4e35-a38b-4b13-9b0d-0e420e35f73a","6954010b-13be-4264-8141-b75bedf1cf3f","33ea6faa-29e7-4eba-990b-a529623da89e","ac571f1f-4bee-49cf-8971-8c65a75b445f","40ef520b-e2ce-454d-97de-8505eb82d341","1469ac28-9800-4932-94b2-71148e24aa2d","83b428fb-3392-4c76-839a-432ae1219109","691488e4-bcc2-49ad-a6bd-f3346fd2a906","1d6d17b3-c185-4fb3-a6c3-dba0f9e94498","90a9fceb-1c5a-4879-99c9-a66ddc1453d9","7328b75b-4580-4012-843e-1371a958f53c","e2229312-d5ba-486f-935c-5e0228fc858e","3cd793f3-244d-45d9-bd61-b1fb95b0c6a5","395fbf97-a97f-4d7f-9e01-a53fa13efdbd","54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","f4180826-bb76-421e-5410-791408daadeb","a29ae041-66b7-46bd-853a-507a05f01ad9"]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"is_sub_org":false,"sub_org_id":"f4180826-bb76-421e-5410-791408daadeb","compare_title":"component 1","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With test suites","value":1},{"title":"Without test suites","value":0}]}},{"is_sub_org":false,"sub_org_id":"test-inactive-comp","compare_title":"test inactive component","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With test suites","value":0},{"title":"Without test suites","value":1}]}},{"is_sub_org":true,"sub_org_id":"sub-org-1","compare_title":"sub-org-1","sub_org_count":0,"component_count":2,"total_value":2,"value_in_millis":0,"compare_reports":[{"is_sub_org":false,"sub_org_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","compare_title":"component 11","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With test suites","value":1},{"title":"Without test suites","value":0}]}},{"is_sub_org":false,"sub_org_id":"test-inactive-sub-org-comp","compare_title":"test inactive sub org component","sub_org_count":0,"component_count":1,"total_value":1,"value_in_millis":0,"compare_reports":null,"section":{"data":[{"title":"With test suites","value":0},{"title":"Without test suites","value":1}]}}],"section":{"data":[{"title":"With test suites","value":1},{"title":"Without test suites","value":1}]}}]`)
		organisation, err := getOrganisationInactive()
		assert.Nil(t, err, "error fetching organisations")
		b, err := getTestComponentsComponentComparison("testComponentsComponentSpec", x, replacements, organisation)
		assert.Nil(t, err, "error processing test components component comparison")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_OpenFindingsBySeverity(t *testing.T) {
	t.Run("Case 1: Successful execution of open findings by severity", func(t *testing.T) {
		replacements := map[string]any{
			"startDate": "2024-03-01",
			"endDate":   "2024-03-31",
			"orgId":     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		}

		// testing case with INFORMATION severity
		responseStringAllSeverities := `{"openFindingsBySeverity":{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":25,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"open_findings_by_severity":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"MEDIUM","doc_count":14,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2024-0727_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-0853_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-13176_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-13176_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-9143_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-9143_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-9681_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}},{"key":"LOW","doc_count":4,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"CVE-2024-11053_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-2004_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2025-0167_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}},{"key":"HIGH","doc_count":3,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2025-0725_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}},{"key":"INFORMATION","doc_count":2,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-2511_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-2511_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}},{"key":"VERY_HIGH","doc_count":2,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2022-48174_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2022-48174_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}}]}}}}`
		xAllSev := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseStringAllSeverities), &xAllSev)
		expectedResultAllSev := []byte(`[{"id":"Very high","value":2,"percentage":9.52,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"High","value":3,"percentage":14.29,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Medium","value":12,"percentage":57.14,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Low","value":4,"percentage":19.05,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}}]`)
		bAllSev, err := transformOpenFindingsBySeverity("openFindingsBySeveritySpec", xAllSev, replacements)
		assert.Equal(t, expectedResultAllSev, []byte(bAllSev))

		// testing case higher count
		responseString2 := `{"openFindingsBySeverity":{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":58,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"open_findings_by_severity":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"MEDIUM","doc_count":38,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":3,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3}]}}]}},{"key":"LOW","doc_count":10,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":4}]}},{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":4}]}},{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":4}]}},{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":4}]}}]}},{"key":"HIGH","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}},{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}},{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}}]}},{"key":"VERY_HIGH","doc_count":4,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}},{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}},{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}}]}}]}}}}`
		x2 := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString2), &x2)
		expectedResult2 := []byte(`[{"id":"Very high","value":4,"percentage":14.81,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"High","value":4,"percentage":14.81,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Medium","value":15,"percentage":55.56,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Low","value":4,"percentage":14.81,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}}]`)
		b2, err := transformOpenFindingsBySeverity("openFindingsBySeveritySpec", x2, replacements)
		assert.Equal(t, expectedResult2, []byte(b2))

		responseString := `{"openFindingsBySeverity":{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":58,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"open_findings_by_severity":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"MEDIUM","doc_count":38,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":3,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":3}]}}]}},{"key":"LOW","doc_count":10,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":4}]}}]}},{"key":"HIGH","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}}]}},{"key":"VERY_HIGH","doc_count":4,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":2}]}}]}}]}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectedResult := []byte(`[{"id":"Very high","value":2,"percentage":25,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"High","value":2,"percentage":25,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Medium","value":3,"percentage":37.5,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Low","value":1,"percentage":12.5,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}}]`)
		b, err := transformOpenFindingsBySeverity("openFindingsBySeveritySpec", x, replacements)
		assert.Equal(t, expectedResult, []byte(b))

		assert.Nil(t, err, "error processing Open Findings By Severity")

	})
}

func Test_SlaBreachesByAsset(t *testing.T) {
	t.Run("Case 1: Successful execution of sla breaches by asset type", func(t *testing.T) {
		replacements := map[string]any{
			"startDate": "2024-03-01",
			"endDate":   "2024-03-31",
			"orgId":     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		}

		responseString := `{"slaBreachedByAssetType":{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":58,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"remediation_key":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"BINARY","doc_count":52,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":4},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4},{"key":"CVE-2024-11053_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-2004_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":3},{"key":"CVE-2024-9143_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":3},{"key":"CVE-2024-9143_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":3},{"key":"CVE-2024-9681_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-0727_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-0727_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-0853_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-13176_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-13176_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-2466_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2}]}},{"key":"CODE","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":2}]}},{"key":"INFRASTRUCTURE","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":2},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2}]}},{"key":"PIPELINE","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":2},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2}]}},{"key":"CORE","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":2}]}}]}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectedResult := []byte(`[{"assetType":"BINARY","findingsPercentage":62,"total":21,"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"assetTypes=BINARY&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL&sla=BREACHED"}]}},{"assetType":"INFRASTRUCTURE","findingsPercentage":15,"total":5,"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"assetTypes=INFRASTRUCTURE&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL&sla=BREACHED"}]}},{"assetType":"PIPELINE","findingsPercentage":12,"total":4,"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"assetTypes=PIPELINE&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL&sla=BREACHED"}]}},{"assetType":"CODE","findingsPercentage":6,"total":2,"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"assetTypes=CODE&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL&sla=BREACHED"}]}},{"assetType":"CORE","findingsPercentage":6,"total":2,"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"assetTypes=CORE&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL&sla=BREACHED"}]}}]`)
		b, err := transformSlasBreachedByAsset("slaBreachedByAssetSpec", x, replacements)
		assert.Nil(t, err, "error processing sla breaches by asset type")

		var arr1, arr2 []SlaRecord

		// Unmarshal the first JSON string
		err = json.Unmarshal([]byte(b), &arr1)
		assert.NoError(t, err, "Failed to parse actual output")

		// Unmarshal the second JSON string
		err = json.Unmarshal([]byte(expectedResult), &arr2)
		assert.NoError(t, err, "Failed to parse expected result")

		// Sort both arrays by 'Asset'
		sortSlaJSON(arr1)
		sortSlaJSON(arr2)

		// Compare the sorted arrays
		assert.Equal(t, arr1, arr2, "The JSON arrays do not match!")
	})

	t.Run("Case 2: Successful execution of sla breaches by asset type with missing types in query response", func(t *testing.T) {
		replacements := map[string]any{
			"startDate": "2024-03-01",
			"endDate":   "2024-03-31",
			"orgId":     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		}

		responseString := `{"slaBreachedByAssetType":{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":58,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"remediation_key":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"BINARY","doc_count":52,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":4},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4},{"key":"CVE-2024-11053_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-2004_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":3},{"key":"CVE-2024-9143_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":3},{"key":"CVE-2024-9143_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":3},{"key":"CVE-2024-9681_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-0727_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-0727_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-0853_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-13176_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-13176_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-2466_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2}]}},{"key":"CODE","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":2}]}}]}}}}`

		//Header
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectedResult := []byte(`[{"assetType":"BINARY","findingsPercentage":91,"total":21,"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"assetTypes=BINARY&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL&sla=BREACHED"}]}},{"assetType":"CODE","findingsPercentage":9,"total":2,"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"assetTypes=CODE&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL&sla=BREACHED"}]}}]`)
		b, err := transformSlasBreachedByAsset("slaBreachedByAssetSpec", x, replacements)
		assert.Nil(t, err, "error processing sla breaches by asset type")

		var arr1, arr2 []SlaRecord

		// Unmarshal the first JSON string
		err = json.Unmarshal([]byte(b), &arr1)
		assert.NoError(t, err, "Failed to parse actual output")

		// Unmarshal the second JSON string
		err = json.Unmarshal([]byte(expectedResult), &arr2)
		assert.NoError(t, err, "Failed to parse expected result")

		// Sort both arrays by 'Asset'
		sortSlaJSON(arr1)
		sortSlaJSON(arr2)

		// Compare the sorted arrays
		assert.Equal(t, arr1, arr2, "The JSON arrays do not match!")
	})
}

type Element struct {
	ID    string `json:"id"`
	Value int    `json:"value"`
}

type SlaRecord struct {
	AssetType          string  `json:"assetType"`
	FindingsPercentage float64 `json:"findingsPercentage"`
	Total              int     `json:"total"`
	ColorScheme        []struct {
		Color0 string `json:"color0"`
		Color1 string `json:"color1"`
	} `json:"colorScheme"`
}

func sortSlaJSON(input []SlaRecord) {
	sort.Slice(input, func(i, j int) bool {
		return input[i].FindingsPercentage > input[j].FindingsPercentage
	})
}

func Test_OpenFindingsByCategory(t *testing.T) {

	t.Run("Case 1: Successful execution of open findings by category", func(t *testing.T) {
		replacements := map[string]any{
			"startDate": "2024-03-01",
			"endDate":   "2024-03-31",
			"orgId":     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		}

		responseString := `{"openFindingsDistributionByCategory":{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":58,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"category":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CONFIGURATION","doc_count":58,"severity":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"MEDIUM","doc_count":38,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":4},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2}]}},{"key":"HIGH","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2}]}},{"key":"VERY_HIGH","doc_count":4,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2}]}}]}},{"key":"VULNERABILITY","doc_count":58,"severity":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"MEDIUM","doc_count":38,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":4},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":3},{"key":"CVE-2024-9143_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":3},{"key":"CVE-2024-9143_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":3},{"key":"CVE-2024-9681_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":2},{"key":"CVE-2024-0727_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-0727_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-0853_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-13176_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-13176_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-2466_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2}]}},{"key":"LOW","doc_count":10,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4},{"key":"CVE-2024-11053_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-2004_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3}]}},{"key":"HIGH","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2}]}},{"key":"VERY_HIGH","doc_count":4,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2}]}}]}}]}}}}`

		// Unmarshal response
		x := map[string]json.RawMessage{}
		err := json.Unmarshal([]byte(responseString), &x)
		assert.Nil(t, err, "error unmarshalling response")

		// Expected result (as Go structs)
		expectedResultJSON := `[{"total":23,"severityDistribution":{"colorScheme":[{"color0":"#FA6D71","color1":"#D5252A"},{"color0":"#FCC16C","color1":"#FF8307"},{"color0":"#FCFF89","color1":"#FDC913"},{"color0":"#9FB6C1","color1":"#648192"}],"data":[{"title":"Very high","value":2},{"title":"High","value":3},{"title":"Medium","value":15},{"title":"Low","value":3}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"Very high","url":"categories=VULNERABILITY&severities=VERY_HIGH&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"High","url":"categories=VULNERABILITY&severities=HIGH&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"Medium","url":"categories=VULNERABILITY&severities=MEDIUM&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"Low","url":"categories=VULNERABILITY&severities=LOW&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}},"categoryName":"Vulnerability","drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"categories=VULNERABILITY&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}},{"total":6,"severityDistribution":{"colorScheme":[{"color0":"#FA6D71","color1":"#D5252A"},{"color0":"#FCC16C","color1":"#FF8307"},{"color0":"#FCFF89","color1":"#FDC913"},{"color0":"#9FB6C1","color1":"#648192"}],"data":[{"title":"Very high","value":2},{"title":"High","value":1},{"title":"Medium","value":3},{"title":"Low","value":0}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"Very high","url":"categories=CONFIGURATION&severities=VERY_HIGH&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"High","url":"categories=CONFIGURATION&severities=HIGH&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"Medium","url":"categories=CONFIGURATION&severities=MEDIUM&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"Low","url":"categories=CONFIGURATION&severities=LOW&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}},"categoryName":"Configuration","drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"categories=CONFIGURATION&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}}]`

		var expectedResult []constants.OpenFindingsByCategoryData
		err = json.Unmarshal([]byte(expectedResultJSON), &expectedResult)
		assert.Nil(t, err, "error unmarshalling expected result")

		// Transform the response
		b, err := transformOpenFindingsDistributionByCategory("openFindingsDistributionByCategorySpec", x, replacements)
		assert.Nil(t, err, "error processing Open Findings By Category")

		// Unmarshal the transformed response
		var actualResult []constants.OpenFindingsByCategoryData
		err = json.Unmarshal([]byte(b), &actualResult)
		assert.Nil(t, err, "error unmarshalling actual result")

		// Compare sorted results
		assert.Equal(t, expectedResult, actualResult)
	})
}

type FindingByDistribution struct {
	Total                int    `json:"total"`
	SecurityToolName     string `json:"securityToolName"`
	ToolId               string `json:"toolId"`
	SeverityDistribution struct {
		ColorScheme []struct {
			Color0 string `json:"color0"`
			Color1 string `json:"color1"`
		} `json:"colorScheme"`
		Data []struct {
			Title string `json:"title"`
			Value int    `json:"value"`
		} `json:"data"`
	} `json:"severityDistribution"`
}

func sortTestArrByToolId(arr1, arr2 []FindingByDistribution) (array1, array2 []FindingByDistribution) {
	// Sort both arrays by toolId
	sort.Slice(arr1, func(i, j int) bool {
		return arr1[i].ToolId < arr1[j].ToolId
	})
	sort.Slice(arr2, func(i, j int) bool {
		return arr2[i].ToolId < arr2[j].ToolId
	})
	return arr1, arr2
}

func Test_OpenFindingsDistributionBySecurityTool(t *testing.T) {

	replacements := map[string]any{
		"startDate": "2024-03-01",
		"endDate":   "2024-03-31",
		"orgId":     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
	}

	t.Run("Case 1: Successful execution of open findings by severity", func(t *testing.T) {

		responseString := `{"openFindingsDistributionBySecurityTool":{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":202,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"86591520-ba4a-11eb-9cab-0a58a9feac02","doc_count":109,"severity":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"LOW","doc_count":83,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_f3f1c5ec96491c64b7fcdac258d03d3ec9c685becf15dafd8f6086801aa83272","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Gosec","doc_count":1}]}},{"key":"703_f4191362f3185762f058ebc6b785a7046ad9d697417cd9022e963cf88e3ca821","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Gosec","doc_count":1}]}},{"key":"703_f77cd4a226ae6676efff1eff2bef95df589633593a66f1ff749ee419347168f4","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Gosec","doc_count":1}]}},{"key":"703_f90ee2fc215278eac34524957b7005092378d84451718d6dbf325bd927f88b0a","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Gosec","doc_count":1}]}},{"key":"703_fa83a66b3831e7acc0e1058a8b58e34cf7d26064ab965c9a3c1611c5902983d3","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Gosec","doc_count":1}]}}]}},{"key":"MEDIUM","doc_count":25,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"118_43d67a5d96f4044de07a99fd6545c5d9e44d7ef4ddc78922191d3a0fd9bd0ed8","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Gosec","doc_count":1}]}},{"key":"118_459756988cc974d77a1b0424ed3dd767f36a6ec249e2f97efcf658a0cbe96858","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Gosec","doc_count":1}]}},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Gosec","doc_count":1}]}}]}},{"key":"HIGH","doc_count":1,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"295_a9423ad7c387947b187c4d0257f53fa7bb0e3e3eae078a319d33f686dbfc1317","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Gosec","doc_count":1}]}}]}}]}}]}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectedResult := []byte(`[{"total":9,"severityDistribution":{"colorScheme":[{"color0":"#FA6D71","color1":"#D5252A"},{"color0":"#FCC16C","color1":"#FF8307"},{"color0":"#FCFF89","color1":"#FDC913"},{"color0":"#9FB6C1","color1":"#648192"}],"data":[{"title":"Very high","value":0},{"title":"High","value":1},{"title":"Medium","value":3},{"title":"Low","value":5}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"Very high","url":"tools=86591520-ba4a-11eb-9cab-0a58a9feac02&severities=VERY_HIGH&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"High","url":"tools=86591520-ba4a-11eb-9cab-0a58a9feac02&severities=HIGH&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"Medium","url":"tools=86591520-ba4a-11eb-9cab-0a58a9feac02&severities=MEDIUM&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"Low","url":"tools=86591520-ba4a-11eb-9cab-0a58a9feac02&severities=LOW&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}},"securityToolName":"Gosec","toolId":"86591520-ba4a-11eb-9cab-0a58a9feac02","drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"tools=86591520-ba4a-11eb-9cab-0a58a9feac02&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}}]`)
		b, err := transformOpenFindingsDistributionBySecurityTool("openFindingsDistributionBySecurityToolSpec", x, replacements)
		assert.Nil(t, err, "error processing Open Findings Distribution By Severity")

		var arr1, arr2 []FindingByDistribution

		// Unmarshal the first JSON string
		err = json.Unmarshal([]byte(b), &arr1)
		assert.NoError(t, err, "Failed to parse actual output")

		// Unmarshal the second JSON string
		err = json.Unmarshal([]byte(expectedResult), &arr2)
		assert.NoError(t, err, "Failed to parse expected result")

		// Sort both arrays by toolId
		sortTestArrByToolId(arr1, arr2)

		// Compare the sorted arrays
		assert.Equal(t, arr1, arr2, "The JSON arrays do not match!")
	})

	t.Run("Case 2: Successful execution of open findings with multiple severities", func(t *testing.T) {

		responseString := `{"openFindingsDistributionBySecurityTool":{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":58,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"9c242178-b1ae-11eb-b5dd-0a58a9feac02","doc_count":2,"severity":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"MEDIUM","doc_count":2,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":1,"tool_display_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"Trivy","doc_count":1}]}}]}}]}}]}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectedResult := []byte(`[{"total":2,"severityDistribution":{"colorScheme":[{"color0":"#FA6D71","color1":"#D5252A"},{"color0":"#FCC16C","color1":"#FF8307"},{"color0":"#FCFF89","color1":"#FDC913"},{"color0":"#9FB6C1","color1":"#648192"}],"data":[{"title":"Very high","value":0},{"title":"High","value":0},{"title":"Medium","value":2},{"title":"Low","value":0}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"Very high","url":"tools=9c242178-b1ae-11eb-b5dd-0a58a9feac02&severities=VERY_HIGH&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"High","url":"tools=9c242178-b1ae-11eb-b5dd-0a58a9feac02&severities=HIGH&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"Medium","url":"tools=9c242178-b1ae-11eb-b5dd-0a58a9feac02&severities=MEDIUM&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"},{"id":"Low","url":"tools=9c242178-b1ae-11eb-b5dd-0a58a9feac02&severities=LOW&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}},"securityToolName":"Trivy","toolId":"9c242178-b1ae-11eb-b5dd-0a58a9feac02","drillDown":{"reportId":"redirect-url","redirectionInfo":[{"id":"","url":"tools=9c242178-b1ae-11eb-b5dd-0a58a9feac02&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}}]`)
		b, err := transformOpenFindingsDistributionBySecurityTool("openFindingsDistributionBySecurityToolSpec", x, replacements)
		assert.Nil(t, err, "error processing Open Findings Distribution By Severity")

		var arr1, arr2 []FindingByDistribution

		// Unmarshal the first JSON string
		err = json.Unmarshal([]byte(b), &arr1)
		assert.NoError(t, err, "Failed to parse actual output")

		// Unmarshal the second JSON string
		err = json.Unmarshal([]byte(expectedResult), &arr2)
		assert.NoError(t, err, "Failed to parse expected result")

		// Sort both arrays by toolId
		sortTestArrByToolId(arr1, arr2)

		// Compare the sorted arrays
		assert.Equal(t, arr1, arr2, "The JSON arrays do not match!")

	})
}

func Test_OpenFindingsBySecurityTool(t *testing.T) {
	t.Run("Case 1: Successful execution of open findings by security tool", func(t *testing.T) {
		replacements := map[string]any{
			"startDate": "2024-03-01",
			"endDate":   "2024-03-31",
			"orgId":     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		}

		responseString := `{"openFindingsBySecurityTool":{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":58,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"tool_name":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":50,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":3,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"407bee1d-bcc7-436d-8159-6cc5dff96f27","doc_count":2}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":3,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"407bee1d-bcc7-436d-8159-6cc5dff96f27","doc_count":2}]}},{"key":"CVE-2024-11053_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"407bee1d-bcc7-436d-8159-6cc5dff96f27","doc_count":3}]}},{"key":"CVE-2024-2004_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"407bee1d-bcc7-436d-8159-6cc5dff96f27","doc_count":3}]}},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":3,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"407bee1d-bcc7-436d-8159-6cc5dff96f27","doc_count":2}]}},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"407bee1d-bcc7-436d-8159-6cc5dff96f27","doc_count":2}]}}]}},{"key":"gosec","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"86591520-ba4a-11eb-9cab-0a58a9feac02","doc_count":3}]}},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"86591520-ba4a-11eb-9cab-0a58a9feac02","doc_count":2}]}}]}},{"key":"trivy","doc_count":2,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"9c242178-b1ae-11eb-b5dd-0a58a9feac02","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"9c242178-b1ae-11eb-b5dd-0a58a9feac02","doc_count":1}]}}]}}]}}}}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectedResult := []byte(`[{"total":6,"findingsPercentage":60,"securityToolName":"grype","toolId":"407bee1d-bcc7-436d-8159-6cc5dff96f27","colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"url":"tools=407bee1d-bcc7-436d-8159-6cc5dff96f27&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}},{"total":2,"findingsPercentage":20,"securityToolName":"gosec","toolId":"86591520-ba4a-11eb-9cab-0a58a9feac02","colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"url":"tools=86591520-ba4a-11eb-9cab-0a58a9feac02&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}},{"total":2,"findingsPercentage":20,"securityToolName":"trivy","toolId":"9c242178-b1ae-11eb-b5dd-0a58a9feac02","colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"url":"tools=9c242178-b1ae-11eb-b5dd-0a58a9feac02&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]}}]`)

		b, err := transformOpenFindingsBySecurityTool("openFindingsBySecurityToolSpec", x, replacements)
		assert.Nil(t, err, "error processing Open Findings By Security Tool")

		var arr1, arr2 []Element

		err = json.Unmarshal(b, &arr1)
		assert.NoError(t, err, "Failed to parse actual output")

		err = json.Unmarshal(expectedResult, &arr2)
		assert.NoError(t, err, "Failed to parse expected result")

		assert.Equal(t, arr1, arr2, "The JSON arrays do not match!")

		assert.Equal(t, expectedResult, []byte(b))
	})
}

func Test_TransformFindingsIdentifiedSince(t *testing.T) {
	t.Run("Case 1: Successful execution of findings identified since", func(t *testing.T) {
		replacements := map[string]any{
			"startDate": "2024-03-01",
		}

		//Header
		responseString := `{"findingsIdentifiedCount":{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":58,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"remediation_status":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"OPEN","doc_count":58,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4},{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":4},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4},{"key":"CVE-2024-11053_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-2004_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":3},{"key":"CVE-2024-9143_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":3},{"key":"CVE-2024-9143_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":3},{"key":"CVE-2024-9681_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":2},{"key":"CVE-2024-0727_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-0727_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-0853_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-13176_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-13176_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-2466_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":2},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":2},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2}]}},{"key":"OPEN","doc_count":58,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4},{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":4},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":4},{"key":"CVE-2024-11053_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-2004_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":3},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":3},{"key":"CVE-2024-9143_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":3},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2}]}},{"key":"IN_PROGRESS","doc_count":58,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":2}]}},{"key":"RISK_ACCEPTED","doc_count":58,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":4}]}}]}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":23}`)
		b, err := transformFindingsIdentfiedSince("openFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - open")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":0}`)
		b, err = transformFindingsIdentfiedSince("resolvedFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - resolved")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":1}`)
		b, err = transformFindingsIdentfiedSince("riskAcceptedFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":2}`)
		b, err = transformFindingsIdentfiedSince("inProgressFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - in progress")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":0}`)
		b, err = transformFindingsIdentfiedSince("falsePositiveFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - false positive")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`[{"id":"Open","value":23,"percentage":88.46,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Resolved","value":0,"percentage":0,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Risk accepted","value":1,"percentage":3.85,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"In progress","value":2,"percentage":7.69,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"False positive","value":0,"percentage":0,"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}}]`)

		b, err = transformFindingsIdentfiedSince("findingsIndentifiedChartSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince section")
		assert.Equal(t, expectResult, []byte(b))
	})
}

func Test_TransformRiskAcceptedFalsePositiveFindings(t *testing.T) {
	t.Run("Case 1: Successful execution of findings identified since", func(t *testing.T) {
		replacements := map[string]any{}

		responseString := `{"riskAcceptedFalsePositiveFindingsCount":{"took":99,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":11,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"risk_accepted_counts":{"doc_count":5,"RA_NOT_EXPIRING_IN_30_DAYS":{"doc_count":2},"RA_EXPIRING_IN_30_DAYS":{"doc_count":3},"RA_EXPIRED":{"doc_count":0}},"remediation_status":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"OPEN","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"118_608276b67cf0c4ea8ddc1b8d98479566fac0cf05347b2f50edc361494b942647","doc_count":1},{"key":"GLS_0129_c5e7299a029b8315113c4b05e3ce48e7f28c628b91dda74ac5531f17db01bd69","doc_count":1},{"key":"GLS_0129_caf92a1ffbc4a74e949c174285ef0e3211493ad00b1ac33a487a96c2464d00fe","doc_count":1},{"key":"GLS_0129_e50265347a403badd08e458aec042d96d1a982d3030d64235de154af7397c2c7","doc_count":1},{"key":"GLS_0129_f4acde97a55cf9a4bd300d24198bea4898266ef966ed095f4161ce97c731127b","doc_count":1},{"key":"GLS_0129_f7ab31b69b3602093e645c913d71607d0a743083bcfd00822438f26396fe1469","doc_count":1}]}},{"key":"RISK_ACCEPTED","doc_count":5,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_3abef734fe886d5d0ac6c3d322354b0b0d4cd86b7cceec5053aff5679ea692be","doc_count":1},{"key":"703_50317f329419f879a89e82182c402b26f44a3e31fe22eb1db3d4ca9b687e3c1a","doc_count":1},{"key":"703_94dcffc35464da5f3a210e253858961bb3de1d9ae87f7f737c502692d5578a91","doc_count":1},{"key":"703_c3c6246ff1c6bbcd83aeab6cfd27dc914b863005693120da43ff3281a0a5dcf8","doc_count":1},{"key":"703_f2317f93884c8ccd24410cb138ed17d2fb212cf96b845cbc34da8af2739321a6","doc_count":1}]}}]}}}}`
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`{"value":11}`)
		b, err := transformRiskAcceptedAndFalsePositiveFindings("totalFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformRiskAcceptanceFalsePositiveFindings header - open")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":5}`)
		b, err = transformRiskAcceptedAndFalsePositiveFindings("riskAcceptedFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":3}`)
		b, err = transformRiskAcceptedAndFalsePositiveFindings("raExpiringIn30DaysFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":0}`)
		b, err = transformRiskAcceptedAndFalsePositiveFindings("raExpiredFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`{"value":0}`)
		b, err = transformRiskAcceptedAndFalsePositiveFindings("falsePositiveFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`[{"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"findingsPercentage":45.45}]`)
		b, err = transformRiskAcceptedAndFalsePositiveFindings("riskAcceptedChartSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult, []byte(b))

		expectResult = []byte(`[{"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"findingsPercentage":0}]`)
		b, err = transformRiskAcceptedAndFalsePositiveFindings("falsePositiveChartSpec", x, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince section")
		assert.Equal(t, expectResult, []byte(b))

		responseString2 := `{"riskAcceptedFalsePositiveFindingsCount":{"took":99,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":11,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"risk_accepted_counts":{"doc_count":5,"RA_NOT_EXPIRING_IN_30_DAYS":{"doc_count":2},"RA_EXPIRING_IN_30_DAYS":{"doc_count":3},"RA_EXPIRED":{"doc_count":0}},"remediation_status":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"OPEN","doc_count":6,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"118_608276b67cf0c4ea8ddc1b8d98479566fac0cf05347b2f50edc361494b942647","doc_count":1},{"key":"GLS_0129_c5e7299a029b8315113c4b05e3ce48e7f28c628b91dda74ac5531f17db01bd69","doc_count":1},{"key":"GLS_0129_caf92a1ffbc4a74e949c174285ef0e3211493ad00b1ac33a487a96c2464d00fe","doc_count":1},{"key":"GLS_0129_e50265347a403badd08e458aec042d96d1a982d3030d64235de154af7397c2c7","doc_count":1},{"key":"GLS_0129_f4acde97a55cf9a4bd300d24198bea4898266ef966ed095f4161ce97c731127b","doc_count":1},{"key":"GLS_0129_f7ab31b69b3602093e645c913d71607d0a743083bcfd00822438f26396fe1469","doc_count":1}]}},{"key":"RISK_ACCEPTED","doc_count":5,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_3abef734fe886d5d0ac6c3d322354b0b0d4cd86b7cceec5053aff5679ea692be","doc_count":1},{"key":"703_50317f329419f879a89e82182c402b26f44a3e31fe22eb1db3d4ca9b687e3c1a","doc_count":1},{"key":"703_94dcffc35464da5f3a210e253858961bb3de1d9ae87f7f737c502692d5578a91","doc_count":1},{"key":"703_c3c6246ff1c6bbcd83aeab6cfd27dc914b863005693120da43ff3281a0a5dcf8","doc_count":1},{"key":"703_f2317f93884c8ccd24410cb138ed17d2fb212cf96b845cbc34da8af2739321a6","doc_count":1}]}},{"key":"FALSE_POSITIVE","doc_count":3,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"FP_abc123","doc_count":1},{"key":"FP_def456","doc_count":1},{"key":"FP_ghi789","doc_count":1}]}},{"key":"RESOLVED","doc_count":4,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"RSLV_111aaa","doc_count":1},{"key":"RSLV_222bbb","doc_count":1},{"key":"RSLV_333ccc","doc_count":1},{"key":"RSLV_444ddd","doc_count":1}]}}]}}}}`
		x2 := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString2), &x2)

		expectResult2 := []byte(`{"value":14}`)
		b2, err := transformRiskAcceptedAndFalsePositiveFindings("totalFindingsSpec", x2, replacements)
		assert.Nil(t, err, "error processing TransformRiskAcceptanceFalsePositiveFindings header - open")
		assert.Equal(t, expectResult2, []byte(b2))

		expectResult2 = []byte(`{"value":5}`)
		b2, err = transformRiskAcceptedAndFalsePositiveFindings("riskAcceptedFindingsSpec", x2, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult2, []byte(b2))

		expectResult2 = []byte(`{"value":3}`)
		b2, err = transformRiskAcceptedAndFalsePositiveFindings("raExpiringIn30DaysFindingsSpec", x2, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult2, []byte(b2))

		expectResult2 = []byte(`{"value":0}`)
		b2, err = transformRiskAcceptedAndFalsePositiveFindings("raExpiredFindingsSpec", x2, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult2, []byte(b2))

		expectResult2 = []byte(`{"value":3}`)
		b2, err = transformRiskAcceptedAndFalsePositiveFindings("falsePositiveFindingsSpec", x2, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult2, []byte(b2))

		expectResult2 = []byte(`[{"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"findingsPercentage":35.71}]`)
		b2, err = transformRiskAcceptedAndFalsePositiveFindings("riskAcceptedChartSpec", x2, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince header - risk accepted")
		assert.Equal(t, expectResult2, []byte(b2))

		expectResult2 = []byte(`[{"colorScheme":[{"color0":"#4696E5","color1":"#0963BD"}],"findingsPercentage":21.43}]`)
		b2, err = transformRiskAcceptedAndFalsePositiveFindings("falsePositiveChartSpec", x2, replacements)
		assert.Nil(t, err, "error processing TransformFindingsIdentifiedSince section")
		assert.Equal(t, expectResult2, []byte(b2))

	})
}

func TestTransformSlaBreachesBySeverity(t *testing.T) {

	replacements := map[string]any{
		"startDate": "2024-03-01",
		"endDate":   "2024-03-31",
		"orgId":     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
	}
	responseStringWithInformation := `{"slaBreachesBySeverity":{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":7,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"sla_breaches_by_severity":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"VERY_HIGH","doc_count":4,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2022-48174_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2022-48174_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}},{"key":"HIGH","doc_count":3,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}},{"key":"INFORMATION","doc_count":3,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}}]}}}}`
	responseStringWithOnlyHigh := `{"slaBreachesBySeverity":{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":7,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"sla_breaches_by_severity":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"VERY_HIGH","doc_count":4,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2022-48174_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2022-48174_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}},{"key":"HIGH","doc_count":3,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}}]}}}}`
	responseString := `{"slaBreachesBySeverity":{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":41,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"sla_breaches_by_severity":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"MEDIUM","doc_count":29,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1},{"key":"trivy","doc_count":1}]}},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":2,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1},{"key":"trivy","doc_count":1}]}},{"key":"22_149ccafd542df3b5ff472e12351281b0cdd4ae5b72aa4dca961ca1639a054f3e","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"22_1cfd8d8f2586dfda60fec800fd9473c5f32fd63c11ccc436bc82d907dd597d8d","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"22_439a22a7cd849de375c4107569f7465063422de82fc975b2be4dee5aa31624dc","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"22_92a10691fe6c9f43abab3408c72bc8af53a24d3956792087d067af3038d7c603","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"22_9a902345cd3b8a0708080434b6a0fb64b3c7d6688a4fa9dc0b4f339e53d4d624","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"22_c27520c4d8a8f8a7d92aee8907314af1be337f74cb941ef1e61f84da049eac47","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"276_3eabb93dd7a2875a58e9aa6e8fc78916e92f6cdf411525c38232049a7e1925ce","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"276_66be8390ad6374362f16ab0a24fc1e1dfc2b8e003f84ae80e353af8f66193a6e","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"276_a771239c4005d86c0a470c46e322aaf88bfeff3a7c5b6092dbdd0f70158fd5f8","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"409_fe1cc67934e9135875003d3c79cf8015d441b6ca53c7042c0ad3277a77544d60","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"676_215993b9fe197329b2d18671eda3ae513b88960aa00c34a0925fe39c6137f909","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"676_6542a0ef2da4b3aa6644359f748577807c9b6b12a02dda018a9d1f5fb28de569","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"676_9ad945a983ebe8340262fe355bf3ee2df035bb552d10533d243866d8bf5465ad","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"676_9c3ffd6491e7a1e2dc551ab21196e4be3c827b4434dfc1d8cef3467bd60a2529","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"CVE-2024-0727_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-0727_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-0853_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-2466_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-9143_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-9143_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-9681_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}},{"key":"LOW","doc_count":7,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"242_7fd5a8713d31ee53fa8915157d45d7ac7c71da61c81e273594ad3d4ec381b9fb","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"242_c20e54c1c4cb3fd3ac8e4c4cd8c010a4e0e5b452364aad07f0c11a85de9bdcf3","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"703_0ebd2da5faeb89ade819eb63e698539880ce4a11449cedb54404b58c5a384531","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"703_ca59f3444a9ef5dcc6c1d3efbc938cba5abaaf1e6674d4474a5911b64dfa3f5f","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"gosec","doc_count":1}]}},{"key":"CVE-2024-11053_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-2004_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}},{"key":"HIGH","doc_count":3,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}},{"key":"VERY_HIGH","doc_count":2,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1,"tool_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"grype","doc_count":1}]}}]}}]}}} }`
	x := map[string]json.RawMessage{}
	json.Unmarshal([]byte(responseString), &x)

	t.Run("Case 1: Successful execution of sla breaches by severity", func(t *testing.T) {

		//header
		expectedResult := []byte(`{"value":39}`)

		b, err := transformSlaBreachesBySeverity("slaBreachesBySeverityHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing sla breaches by severity")
		assert.Equal(t, expectedResult, []byte(b))
	})

	t.Run("Case 2: Successful execution for subheader by severity", func(t *testing.T) {
		// Expected result
		expectedResult := []byte(`{"subHeader":[{"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"url":"severities=VERY_HIGH&sla=BREACHED&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]},"title":"Very high","value":2,"color":"#EA4F54"},{"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"url":"severities=HIGH&sla=BREACHED&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]},"title":"High","value":3,"color":"#FE9D33"},{"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"url":"severities=MEDIUM&sla=BREACHED&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]},"title":"Medium","value":27,"color":"#FCE44E"},{"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"url":"severities=LOW&sla=BREACHED&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]},"title":"Low","value":7,"color":"#738E9D"}]}`)

		// Actual result
		b, err := transformSlaBreachesBySeverity("slaBreachesBySeveritySubHeaderSpec", x, replacements)
		assert.Nil(t, err, "error processing sla breaches by severity")
		assert.JSONEq(t, string(expectedResult), string(b), "The JSON outputs do not match")
	})

	t.Run("Case 4: Successful execution for section by severity", func(t *testing.T) {

		// //section
		expectedResult := []byte(`[{"name":"Very high","value":5.13},{"name":"High","value":7.69},{"name":"Medium","value":69.23},{"name":"Low","value":17.95}]`)

		b, err := transformSlaBreachesBySeverity("slaBreachesBySeveritySpec", x, replacements)
		assert.Nil(t, err, "error processing sla breaches by severity")
		assert.Equal(t, expectedResult, []byte(b))
	})

	t.Run("Case 5: Successful execution for section by severity with only high", func(t *testing.T) {
		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseStringWithOnlyHigh), &x)

		expectedTotal := []byte(`{"value":7}`)
		expectedSubheader := []byte(`{"subHeader":[{"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"url":"severities=VERY_HIGH&sla=BREACHED&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]},"title":"Very high","value":4,"color":"#EA4F54"},{"drillDown":{"reportId":"redirect-url","redirectionInfo":[{"url":"severities=HIGH&sla=BREACHED&remediationStatus=OPEN+IN_PROGRESS&triageStatus=ALL"}]},"title":"High","value":3,"color":"#FE9D33"}]}`)
		expectedSection := []byte(`[{"name":"Very high","value":57.14},{"name":"High","value":42.86},{"name":"Medium","value":0},{"name":"Low","value":0}]`)

		t.Run("Case 5.1: Successful execution of sla breaches by severity with only high", func(t *testing.T) {
			//header
			b, err := transformSlaBreachesBySeverity("slaBreachesBySeverityHeaderSpec", x, replacements)
			assert.Nil(t, err, "error processing sla breaches by severity")
			assert.Equal(t, expectedTotal, []byte(b))
		})

		t.Run("Case 5.2: Successful execution for subheader by severity with only high", func(t *testing.T) {
			// Actual result
			b, err := transformSlaBreachesBySeverity("slaBreachesBySeveritySubHeaderSpec", x, replacements)
			assert.Nil(t, err, "error processing sla breaches by severity")
			assert.JSONEq(t, string(expectedSubheader), string(b), "The JSON outputs do not match")
		})

		t.Run("Case 5.3: Successful execution for section by severity with only high", func(t *testing.T) {
			//section
			b, err := transformSlaBreachesBySeverity("slaBreachesBySeveritySpec", x, replacements)
			assert.Nil(t, err, "error processing sla breaches by severity")
			assert.Equal(t, expectedSection, []byte(b))
		})
	})

	t.Run("Case 6: Successful execution for section by severity with Information severity level", func(t *testing.T) {
		xInfo := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseStringWithInformation), &xInfo)
		//section
		expectedResult := []byte(`[{"name":"Very high","value":57.14},{"name":"High","value":42.86},{"name":"Medium","value":0},{"name":"Low","value":0}]`)

		bInfo, err := transformSlaBreachesBySeverity("slaBreachesBySeveritySpec", xInfo, replacements)
		assert.Nil(t, err, "error processing sla breaches by severity")
		assert.Equal(t, expectedResult, []byte(bInfo))

	})
}

func TestTransformOpenFindingsBySlaStatus(t *testing.T) {
	replacements := map[string]any{
		"startDate": "2024-03-01",
		"endDate":   "2024-03-31",
		"orgId":     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
	}

	responseString := `{"openFindingsBySLAStatus":{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":31,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"open_findings_by_sla_status":{"buckets":{"non_sla_breached":{"doc_count":8,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2022-48174_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":1},{"key":"CVE-2022-48174_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":1},{"key":"CVE-2024-13176_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-13176_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-2511_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-2511_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2025-0167_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2025-0725_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1}]}},"sla_breached":{"doc_count":23,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":2},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":2},{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":1},{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":1},{"key":"CVE-2024-0727_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-0727_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-0853_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-11053_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-2004_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-2466_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":1},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-9143_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-9143_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-9681_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1}]}}}}}}}`

	x := map[string]json.RawMessage{}
	json.Unmarshal([]byte(responseString), &x)

	t.Run("Case 1: Successful execution of open findings by sla status for total", func(t *testing.T) {

		//total
		expectedResult := []byte(`{"value":29}`)

		b, err := transformOpenFindingsBySlaStatus("openFindingsSpec", x, replacements)
		assert.Nil(t, err, "error processing open findings by sla status")
		assert.Equal(t, expectedResult, []byte(b))
	})

	t.Run("Case 2: Successful execution of open findings by sla status for within sla", func(t *testing.T) {

		//within sla
		expectedResult := []byte(`{"value":8}`)

		b, err := transformOpenFindingsBySlaStatus("withinSLASpec", x, replacements)
		assert.Nil(t, err, "error processing open findings by sla status")
		assert.Equal(t, expectedResult, []byte(b))

	})

	t.Run("Case 2: Successful execution of open findings by sla status for breached sla", func(t *testing.T) {

		//breached sla
		expectedResult := []byte(`{"value":21}`)

		b, err := transformOpenFindingsBySlaStatus("breachedSLASpec", x, replacements)
		assert.Nil(t, err, "error processing open findings by sla status")
		assert.Equal(t, expectedResult, []byte(b))
	})

	t.Run("Case 3: Successful execution of open findings by sla status for section", func(t *testing.T) {

		//section
		expectedResult := []byte(`[{"name":"Within SLA","value":27.59},{"name":"Breached SLA","value":72.41}]`)

		b, err := transformOpenFindingsBySlaStatus("openFindingsBySLAStatusChartSpec", x, replacements)
		assert.Nil(t, err, "error processing open findings by sla status")
		assert.Equal(t, expectedResult, []byte(b))
	})
}

func TestTransformOpenFindingsByReviewStatus(t *testing.T) {

	replacements := map[string]any{
		"startDate": "2024-03-01",
		"endDate":   "2024-03-31",
		"orgId":     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
	}
	responseString_nil := `{"openFindingsByReviewStatus":{"took":338,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":27,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"triage_status":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"UNREVIEWED","doc_count":27,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"676_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":1},{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":1},{"key":"CVE-2022-48174_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":1},{"key":"CVE-2022-48174_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":1},{"key":"CVE-2023-42366_1c4fb39de1d312b62ede343a43cf70a5b24468718370d4920d9ad3ac1435e05b","doc_count":1},{"key":"CVE-2023-42366_dbda5b121719c170b59361d28099205215dfd2f788dd3bd73091f0105a9eae5f","doc_count":1},{"key":"CVE-2024-0727_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-0727_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-0853_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-11053_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-13176_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-13176_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-2004_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-2398_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-2466_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-28182_7e0611c1f2b49f47e4c73061611a14a01dccffd41de53bf9834e061314c7eb4b","doc_count":1},{"key":"CVE-2024-4741_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-4741_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-5535_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-5535_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-7264_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-8096_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2024-9143_0e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-9143_36068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-9681_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2025-0167_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1},{"key":"CVE-2025-0725_31df2e2a991f00c0318442bea284b20355bed2bc2752b3ef10bd0693c114bf9a","doc_count":1}]}}]}}}}`
	responseString := `{"openFindingsByReviewStatus":{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":180,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"triage_status":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"UNREVIEWED","doc_count":133,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"703_a6acf33a3adea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":8},{"key":"GLS_0111_3d864b259c2774f8e2d627109c95e7e234d71df858563ff2107e499b70f5a168","doc_count":1},{"key":"GLS_0140_3f0fec8f644b29db869aa34fbb95d5c34e8d6211f3ce1e1cb6bbe26de1be61fe","doc_count":1}]}},{"key":"IN_REVIEW","doc_count":5,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"676_a6acf33a3a1dea98426d91bbbb45865c16c4528ac3f8b595be5c8288d6e79984f","doc_count":1},{"key":"CVE-2024-0727_236068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-5535_30e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-5535_436068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-9143_536068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1}]}},{"key":"FIX_REQUIRED","doc_count":1,"tracking_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"CVE-2024-4741_636068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-0727_736068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-5535_80e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-5535_936068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1},{"key":"CVE-2024-5535_110e13927b5af36764314cbff0f7d2362dd05ac2471fda99f25aa0a4e77f6e6d16","doc_count":1},{"key":"CVE-2024-5535_1236068c9df051bdec7729464686c0812d4e0835af30041aa23a7d5ff80ddb8ea0","doc_count":1}]}}]}}}}`
	x := map[string]json.RawMessage{}
	json.Unmarshal([]byte(responseString), &x)

	x_nil := map[string]json.RawMessage{}
	json.Unmarshal([]byte(responseString_nil), &x_nil)

	t.Run("Case 1: Successful execution of open findings by review status - unreviewedspec", func(t *testing.T) {

		//header
		expectedResult := []byte(`{"value":3}`)
		b, err := transformOpenFindingsByReviewStatus("unreviewedSpec", x, replacements)
		assert.Nil(t, err, "error processing open findings by review status - unreviewedSpec")
		assert.Equal(t, expectedResult, []byte(b))

		expectedResult_nil := []byte(`{"value":27}`)

		b_nil, err_nil := transformOpenFindingsByReviewStatus("unreviewedSpec", x_nil, replacements)
		assert.Nil(t, err_nil, "error processing open findings by review status - unreviewedSpec")
		assert.Equal(t, expectedResult_nil, []byte(b_nil))
	})

	t.Run("Case 2: Successful execution of open findings by review status - awaitingApprovalSpec", func(t *testing.T) {

		//header
		expectedResult := []byte(`{"value":5}`)
		b, err := transformOpenFindingsByReviewStatus("awaitingApprovalSpec", x, replacements)
		assert.Nil(t, err, "error processing open findings by review status - awaitingApprovalSpec")
		assert.Equal(t, expectedResult, []byte(b))
		expectedResult_nil := []byte(`{"value":0}`)
		b_nil, err_nil := transformOpenFindingsByReviewStatus("awaitingApprovalSpec", x_nil, replacements)
		assert.Nil(t, err_nil, "error processing open findings by review status - awaitingApprovalSpec")
		assert.Equal(t, expectedResult_nil, []byte(b_nil))
	})
	t.Run("Case 3: Successful execution of open findings by review status - fixRequiredSpec", func(t *testing.T) {

		//header
		expectedResult := []byte(`{"value":6}`)
		b, err := transformOpenFindingsByReviewStatus("fixRequiredSpec", x, replacements)
		assert.Nil(t, err, "error processing open findings by review status - fixRequiredSpec")
		assert.Equal(t, expectedResult, []byte(b))
		expectedResult_nil := []byte(`{"value":0}`)
		b_nil, err_nil := transformOpenFindingsByReviewStatus("fixRequiredSpec", x_nil, replacements)
		assert.Nil(t, err_nil, "error processing open findings by review status - fixRequiredSpec")
		assert.Equal(t, expectedResult_nil, []byte(b_nil))
	})

	t.Run("Case 4: Successful execution for section of open findings by review status", func(t *testing.T) {
		//section chart
		expectedResult := []byte(`[{"name":"Unreviewed","value":21.43},{"name":"Fix required","value":42.86},{"name":"Awaiting approval","value":35.71}]`)

		b, err := transformOpenFindingsByReviewStatus("openFindingsByReviewStatusSectionSpec", x, replacements)
		assert.Nil(t, err, "error processing open findings by review status")
		assert.Equal(t, expectedResult, []byte(b))

		expectedResult_nil := []byte(`[{"name":"Unreviewed","value":100}]`)
		b_nil, err_nil := transformOpenFindingsByReviewStatus("openFindingsByReviewStatusSectionSpec", x_nil, replacements)
		assert.Nil(t, err_nil, "error processing open findings by review status")

		assert.Equal(t, expectedResult_nil, []byte(b_nil))

	})
}
func TestTransformFindingsRemediationTrend(t *testing.T) {
	t.Run("Case 1: Successful execution of findings remediation trend", func(t *testing.T) {
		replacements := map[string]any{
			"duration": "week",
		}

		responseString := `{
			"findingsRemediationTrend": {
				"aggregations": {
					"findings_remediation_trend": {
						"value": {
							"2024-03-01": {
								"Open": 10,
								"BreachedSLA": 5,
								"ClosedWithinSLA": 3
							},
							"2024-03-02": {
								"Open": 15,
								"BreachedSLA": 7,
								"ClosedWithinSLA": 4
							}
						}
					}
				}
			}
		}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"id":"Open","data":[{"x":"Friday","y":10},{"x":"Saturday","y":15}],"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Closed within SLA","data":[{"x":"Friday","y":3},{"x":"Saturday","y":4}],"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Breached SLA","data":[{"x":"Friday","y":5},{"x":"Saturday","y":7}],"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}}]`)
		b, err := transformFindingsRemediationTrend("findingsRemediationTrendSpec", x, replacements)
		assert.Nil(t, err, "error processing transformFindingsRemediationTrend")
		assert.Equal(t, expectResult, []byte(b))
	})

	t.Run("Case 2: Missing data key", func(t *testing.T) {
		replacements := map[string]any{
			"duration": "week",
		}

		x := map[string]json.RawMessage{}

		_, err := transformFindingsRemediationTrend("findingsRemediationTrendSpec", x, replacements)
		assert.NotNil(t, err, "expected error due to missing data key")
	})

	t.Run("Case 3: Invalid duration type", func(t *testing.T) {
		replacements := map[string]any{
			"duration": "invalid",
		}

		responseString := `{
			"findingsRemediationTrend": {
				"aggregations": {
					"findings_remediation_trend": {
						"value": {
							"2024-03-01": {
								"Open": 10,
								"BreachedSLA": 5,
								"ClosedWithinSLA": 3
							}
						}
					}
				}
			}
		}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		_, err := transformFindingsRemediationTrend("findingsRemediationTrendSpec", x, replacements)
		assert.NotNil(t, err, "expected error due to invalid duration type")
	})

	t.Run("Case 4: Successful execution of findings remediation trend - App sec", func(t *testing.T) {
		replacements := map[string]any{
			"duration": "week",
		}

		responseString := `{
			"findingsRemediationTrend": {
				"aggregations": {
					"findings_remediation_trend": {
						"value": {
							"2024-03-01": {
								"Open": 10,
								"BreachedSLA": 5,
								"ClosedWithinSLA": 3,
								"New": 2
							},
							"2024-03-02": {
								"Open": 15,
								"BreachedSLA": 7,
								"ClosedWithinSLA": 4,
								"New": 5
							}
						}
					}
				}
			}
		}`

		x := map[string]json.RawMessage{}
		json.Unmarshal([]byte(responseString), &x)

		expectResult := []byte(`[{"id":"Open","data":[{"x":"Friday","y":10},{"x":"Saturday","y":15}],"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Closed within SLA","data":[{"x":"Friday","y":3},{"x":"Saturday","y":4}],"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"New","data":[{"x":"Friday","y":2},{"x":"Saturday","y":5}],"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}},{"id":"Breached SLA","data":[{"x":"Friday","y":5},{"x":"Saturday","y":7}],"yAxisFormatter":{"appendUnitValue":"Findings","type":"APPEND_TEXT"}}]`)
		b, err := transformFindingsRemediationTrend("findingsRemediationTrendAppSec", x, replacements)
		assert.Nil(t, err, "error processing transformFindingsRemediationTrend")
		assert.Equal(t, expectResult, []byte(b))
	})

}
