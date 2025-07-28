package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/calculi-corp/reports-service/db"
	"github.com/calculi-corp/reports-service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/stretchr/testify/assert"
	uberMock "go.uber.org/mock/gomock"
)

type MockOpensearchConnection struct{}

func (m *MockOpensearchConnection) GetOpensearchConnection() (*opensearch.Client, error) {
	return &opensearch.Client{}, nil
}

type MockOpensearchClient struct{}

func (m *MockOpensearchClient) Msearch(ctx context.Context, body interface{}) (string, error) {
	responseData := `{
		"responses": [
			{"hits": {"total": 10}},
			{"hits": {"total": 20}}
		]
	}`
	return responseData, nil
}

type MockDbOperations struct{}

type MockLogger struct{}

type MockOpensearchConfig struct {
	CheckOpensearchClientFunc   func(ctx context.Context, instance *opensearch.Client) bool
	GetOpensearchConnectionFunc func() (*opensearch.Client, error)
}

func (m *MockOpensearchConfig) CheckOpensearchClient(ctx context.Context, instance *opensearch.Client) bool {
	if m.CheckOpensearchClientFunc != nil {
		return m.CheckOpensearchClientFunc(ctx, instance)
	}
	return false
}

func (m *MockOpensearchConfig) GetOpensearchConnection() (*opensearch.Client, error) {
	if m.GetOpensearchConnectionFunc != nil {
		fmt.Println(m.GetOpensearchConnectionFunc())
		return m.GetOpensearchConnectionFunc()
	}
	return nil, errors.New("mock GetOpensearchConnection not implemented")
}

func TestGetMultiQueryResponse(t *testing.T) {
	qm := map[string]db.DbQuery{
		"query1": {QueryString: `{"query": {"match_all": {}}}`, AliasName: "index1"},
		"query2": {QueryString: `{"query": {"match": {"field": "value"}}}`, AliasName: "index2"},
	}

	mockConfig := &MockOpensearchConfig{}
	mockCtrl := uberMock.NewController(t)
	responseBody := `{"responses":[{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}]}`

	mockTransport := mocks.NewMockTransport(mockCtrl)
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()
	mockClient := &opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}
	mockConfig.GetOpensearchConnectionFunc = func() (*opensearch.Client, error) {
		return mockClient, nil
	}

	GetMultiQueryResponse(qm)

}

func TestHasNoHits(t *testing.T) {
	tests := []struct {
		name           string
		responseString string
		expectedResult bool
	}{
		{
			name:           "No hits - empty response",
			responseString: `{"hits":{"total":{"value":0}}}`,
			expectedResult: true,
		},
		{
			name:           "Hits found",
			responseString: `{"hits":{"total":{"value":5}}}`,
			expectedResult: false,
		},
		{
			name:           "Invalid JSON",
			responseString: `{"hits":{"total":{"value":"invalid"}}}`,
			expectedResult: true,
		},
		{
			name:           "No hits field",
			responseString: `{"not_hits":{"total":{"value":0}}}`,
			expectedResult: true,
		},
		{
			name:           "Empty response",
			responseString: `{}`,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasNoHits(tt.responseString)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetMultiQueryResponse_2(t *testing.T) {
	qm := map[string]db.DbQuery{
		"query1": {QueryString: `{"query": {"match_all": {}}}`, AliasName: "index1"},
		"query2": {QueryString: `{"query": {"match": {"field": "value"}}}`, AliasName: "index2"},
	}

	GetMultiQueryResponse(qm)

}

func TestGetBucketDates(t *testing.T) {
	testCases := []struct {
		name           string
		aggrBy         string
		startDate      string
		endDate        string
		expectedResult []map[string]string
		expectedError  error
	}{
		{
			name:      "Week aggregation",
			aggrBy:    "week",
			startDate: "2024-05-01",
			endDate:   "2024-05-31",
			expectedResult: []map[string]string{
				{"startDate": "2024-05-01", "endDate": "2024-05-05"},
				{"startDate": "2024-05-06", "endDate": "2024-05-12"},
				{"startDate": "2024-05-13", "endDate": "2024-05-19"},
				{"startDate": "2024-05-20", "endDate": "2024-05-26"},
				{"startDate": "2024-05-27", "endDate": "2024-05-31"},
			},
			expectedError: nil,
		},
		{
			name:      "Day aggregation",
			aggrBy:    "day",
			startDate: "2024-05-01",
			endDate:   "2024-05-03",
			expectedResult: []map[string]string{
				{"startDate": "2024-05-01", "endDate": "2024-05-01"},
				{"startDate": "2024-05-02", "endDate": "2024-05-02"},
				{"startDate": "2024-05-03", "endDate": "2024-05-03"},
			},
			expectedError: nil,
		},
		{
			name:      "Month aggregation",
			aggrBy:    "month",
			startDate: "2024-05-01",
			endDate:   "2024-07-15",
			expectedResult: []map[string]string{
				{"startDate": "2024-05-01", "endDate": "2024-05-31"},
				{"startDate": "2024-06-01", "endDate": "2024-06-30"},
				{"startDate": "2024-07-01", "endDate": "2024-07-15"},
			},
			expectedError: nil,
		},
		{
			name:           "Invalid aggregation",
			aggrBy:         "invalid",
			startDate:      "2024-05-01",
			endDate:        "2024-05-31",
			expectedResult: nil,
			expectedError:  db.ErrInternalServer,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetBucketDates(tc.aggrBy, tc.startDate, tc.endDate)

			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				assert.Equal(t, len(tc.expectedResult), len(result))
				for i := range result {
					assert.Equal(t, tc.expectedResult[i]["startDate"], result[i]["startDate"])
					assert.Equal(t, tc.expectedResult[i]["endDate"], result[i]["endDate"])
				}
			}
		})
	}
}
func TestAddDateBuckets(t *testing.T) {
	testCases := []struct {
		name           string
		input          json.RawMessage
		newDates       []map[string]string
		expectedResult json.RawMessage
		expectedError  error
	}{
		{
			name: "Valid input and new dates",
			input: json.RawMessage(`[
				{
					"data": [
						{"x": "2024-05-30", "value": 10},
						{"x": "2024-05-31", "value": 15}
					]
				}
			]`),
			newDates: []map[string]string{
				{"startDate": "2024-05-30", "endDate": "2024-06-05"},
			},
			expectedResult: json.RawMessage(`[
				{
					"data": [
						{"x": "2024-05-30", "value": 10, "date": {"startDate": "2024-05-30", "endDate": "2024-06-05"}},
						{"x": "2024-05-31", "value": 15}
					]
				}
			]`),
			expectedError: nil,
		},
		{
			name:           "Empty input",
			input:          json.RawMessage(`[]`),
			newDates:       []map[string]string{},
			expectedResult: nil,
			expectedError:  db.ErrInternalServer,
		},
		{
			name: "Invalid input",
			input: json.RawMessage(`[
				{
					"data": {}
				}
			]`),
			newDates: []map[string]string{
				{"startDate": "2024-05-30", "endDate": "2024-06-05"},
			},
			expectedResult: nil,
			expectedError:  db.ErrInternalServer,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := AddDateBuckets(tc.input, tc.newDates)
			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.JSONEq(t, string(tc.expectedResult), string(result))
			}
		})
	}
}

func TestDateDiff(t *testing.T) {
	qm := map[string]db.DbQuery{
		"1": {AliasName: "commit_data", QueryString: "{\"size\":0,\"query\":{\"bool\":{\"filter\":[{\"range\":{\"commit_timestamp\":{\"gte\":\"2024-02-01 00:00:00\",\"lte\":\"2024-02-30 20:59:59\",\"format\":\"yyyy-MM-dd HH:mm:ss\"}}},{\"term\":{\"org_id\":\"2cab10cc-cd9d-11ed-afa1-0242ac120002\"}},{\"bool\":{\"must_not\":{\"term\":{\"author\":\"github-actions[bot]\"}}}},{\"bool\":{\"must_not\":{\"prefix\":{\"branch\":{\"value\":\"dependabot\"}}}}}]}},\"aggs\":{\"commits\":{\"scripted_metric\":{\"init_script\":\"state.statusMap = [:];\",\"map_script\":\"def map = state.statusMap;def key = doc.repository_name.value+'_'+doc.commit_id.value+'_'+doc.branch.value;def v = ['branch': doc.branch.value, 'component_id': doc.component_id.value, 'commit_id': doc.commit_id.value, 'org_id': doc.org_id.value, 'commit_timestamp': doc.commit_timestamp.value, 'component_name': doc.component_name.value, 'repository_name': doc.repository_name.value, 'author': doc.author.value];map.put(key, v);\",\"combine_script\":\"return state.statusMap;\",\"reduce_script\":\"def tmpMap=[:],resultMap=new HashMap(),commits=new ArrayList();for(response in states){if(response!=null){for(key in response.keySet()){tmpMap.put(key,response.get(key));}}}def authorSet=new HashSet();for(key in tmpMap.keySet()){authorSet.add(tmpMap.get(key).author);}resultMap.put('commits_count',tmpMap.size());def average=0;if(authorSet.size()>0){average=tmpMap.size()/authorSet.size();}def devRecord=['title':'Active Developers','value':authorSet.size()];def commitAvgRecord=['title':'Commits / active dev','value':average];resultMap.put('avg',commitAvgRecord);resultMap.put('dev',devRecord);return resultMap;\"}}}}"},
		"2": {AliasName: "automation_run_status", QueryString: "{\"size\":0,\"query\":{\"bool\":{\"filter\":[{\"range\":{\"status_timestamp\":{\"gte\":\"2024-02-01 00:00:00\",\"lte\":\"2024-02-30 20:59:59\",\"format\":\"yyyy-MM-dd HH:mm:ss\"}}},{\"term\":{\"org_id\":\"2cab10cc-cd9d-11ed-afa1-0242ac120002\"}},{\"term\":{\"job_id\":\"\"}},{\"term\":{\"step_id\":\"\"}},{\"bool\":{\"should\":[{\"term\":{\"status\":\"SUCCEEDED\"}},{\"term\":{\"status\":\"FAILED\"}},{\"term\":{\"status\":\"TIMED_OUT\"}},{\"term\":{\"status\":\"ABORTED\"}}]}}]}},\"aggs\":{\"automation_run\":{\"scripted_metric\":{\"init_script\":\"state.data_map=[:];\",\"map_script\":\"def map = state.data_map;def key = doc.run_id.value + '_' + doc.status.value;def v = ['run_id': doc.run_id.value, 'status': doc.status.value];map.put(key, v);\",\"combine_script\":\"return state.data_map;\",\"reduce_script\":\"def tmpMap = [: ], out = [: ], resultMap = new HashMap(), countMap = new HashMap(), totalCount = 0.0;for (response in states) {if (response != null) {for (key in response.keySet()) {def record = response.get(key);if (record.status == 'SUCCEEDED') {record.status = 'Success';} else if (record.status == 'FAILED' || record.status == 'TIMED_OUT' || record.status == 'ABORTED') {record.status = 'Failure';}tmpMap.put(key, record);}}}for (key in tmpMap.keySet()) {def mapRecord = tmpMap.get(key);if (countMap.containsKey(mapRecord.status)) {def count = countMap.get(mapRecord.status);countMap.put(mapRecord.status, count + 1);} else {countMap.put(mapRecord.status, 1);}}def statusArray = ['Success', 'Failure'];for (key in statusArray) {def count = 0, dataMap = new HashMap();if (countMap.containsKey(key)) {count = countMap.get(key);}totalCount += count;}return totalCount;\"}}}}"},
		// "3"  :db.DbQuery { AliasName : "automation_run_status", QueryString : "" },
	}
	_, mqStr, _ := constructMultiQuery(qm)
	// mqStr, err := GetMultiQueryResponse(qm)
	// assert.NotNil(t, err)
	assert.NotEmpty(t, mqStr)
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	resp, err := Get(server.URL, "dummy_token")
	if err != nil {
		t.Errorf("Get returned error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetQueryResponse(t *testing.T) {
	queryMap := map[string]db.DbQuery{
		"query1": {
			AliasName:   "",
			QueryString: "",
		},
		"query2": {
			AliasName:   "",
			QueryString: "",
		},
	}
	GetQueryResponse(queryMap)
}

func TestRunQuery(t *testing.T) {
	query := db.DbQuery{
		QueryString: "sample_query",
		AliasName:   "sample_alias",
	}

	responseChannel := make(chan OpenSearchResponse)

	go func() {
		err := runQuery("sample_key", query, responseChannel)

		if err != nil && err.Error() == "internal server error" {
			return
		}

		t.Errorf("RunQuery returned unexpected error: %v", err)
	}()

	select {
	case <-responseChannel:
	case <-time.After(5 * time.Second):
		t.Error("Test timed out")
	}
}

func TestIsResponseEmpty_EmptyMap(t *testing.T) {
	qr := make(map[string]json.RawMessage)

	isEmpty := IsResponseEmpty(qr)

	if !isEmpty {
		t.Error("Expected IsResponseEmpty to return true for an empty map, but it returned false")
	}
}

func TestIsResponseEmpty_NonEmptyMap(t *testing.T) {
	qr := make(map[string]json.RawMessage)
	qr["key1"] = json.RawMessage(`{"hits":{"total":{"value":10}}}`)
	qr["key2"] = json.RawMessage(`{"hits":{"total":{"value":0}}}`)

	isEmpty := IsResponseEmpty(qr)

	if isEmpty {
		t.Error("Expected IsResponseEmpty to return false for a non-empty map, but it returned true")
	}
}

func TestIsResponseEmpty_InvalidJSON(t *testing.T) {
	qr := make(map[string]json.RawMessage)
	qr["key1"] = json.RawMessage(`{"hits":{"total":{"value":10}`) // Missing closing brace

	isEmpty := IsResponseEmpty(qr)

	if !isEmpty {
		t.Error("Expected IsResponseEmpty to return true for invalid JSON data, but it returned false")
	}
}

func TestGetCountQueryResponse(t *testing.T) {
	query := db.DbQuery{
		AliasName:   "sample_alias1",
		QueryString: "sample_query1",
	}

	GetCountQueryResponse(query)
}

func TestHasField(t *testing.T) {
	tests := []struct {
		field    string
		index    string
		expected bool
	}{
		{"branch_id", "cb_test_suites", true},
		{"branch_id", "cb_test_cases", true},
		{"branch_id", "non_existing_index", false},
		{"non_existing_field", "cb_test_suites", false},
		{"non_existing_field", "non_existing_index", false}, // New test case for non-existing field
	}

	for _, test := range tests {
		actual, err := HasField(test.field, test.index, db.FieldIndexMap)

		if err != nil {
			expectedError := fmt.Sprintf("input field not found in FieldIndexMap")
			if err.Error() != expectedError {
				t.Errorf("For field '%s' and index '%s', expected error '%s' but got '%v'",
					test.field, test.index, expectedError, err)
			}
		} else {
			if actual != test.expected {
				t.Errorf("For field '%s' and index '%s', expected %v but got %v",
					test.field, test.index, test.expected, actual)
			}
		}
	}
}
