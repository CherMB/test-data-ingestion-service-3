package db

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	api "github.com/calculi-corp/api/go"
	"github.com/calculi-corp/common/grpc"
	"github.com/calculi-corp/config"
	coredataMock "github.com/calculi-corp/core-data-cache/mock"
	mock_grpc_client "github.com/calculi-corp/grpc-client/mock"
	testutil "github.com/calculi-corp/grpc-testutil"
	"github.com/calculi-corp/log"
	"github.com/calculi-corp/reports-service/cache"
	"github.com/calculi-corp/reports-service/mocks"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func init() {
	config.Config.DefineStringFlag("opensearch.endpoint", "", "The Open search Endpoint")
	config.Config.DefineStringFlag("opensearch.user", "", "The Open search Username")
	config.Config.DefineStringFlag("opensearch.pwd", "", "The Open search password")
	config.Config.DefineStringFlag("report.definition.filepath", "../resources/", "Report Definiftion filepath")
	testutil.SetUnitTestConfig()
	config.Config.Set("logging.level", "INFO")
}

func TestGetWidgetDefinitionList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockCache := coredataMock.NewMockResourceCacheI(mockCtrl)
	cache.CoreDataResourceCache = mockCache
	mockCache.EXPECT().GetChildrenOfType(gomock.Any(), gomock.Any()).Return([]string{"1"}).Times(2).AnyTimes()
	resource := api.Resource{Id: "Id1", ParentId: "PId1", Name: "R1", Type: api.ResourceType_RESOURCE_TYPE_DASHBOARD}
	mockCache.EXPECT().Get(gomock.Any()).Return(&resource).Times(2).AnyTimes()

	tests := []struct {
		dashboardName string
		expectError   bool
	}{
		{
			dashboardName: "software-delivery-activity",
			expectError:   false,
		},
		{
			dashboardName: "security-insights",
			expectError:   false,
		},
		{
			dashboardName: "incorrect",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		_, err := GetWidgetDefinitionList(tt.dashboardName)
		if tt.expectError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

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
		return m.GetOpensearchConnectionFunc()
	}
	return nil, errors.New("mock GetOpensearchConnection not implemented")
}

func TestGetOpenSearchClient(t *testing.T) {
	mockConfig := &MockOpensearchConfig{}

	t.Run("Initialize New OpenSearch Client", func(t *testing.T) {
		mockConfig.CheckOpensearchClientFunc = func(ctx context.Context, instance *opensearch.Client) bool {
			return false
		}

		mockClient := &opensearch.Client{}
		mockConfig.GetOpensearchConnectionFunc = func() (*opensearch.Client, error) {
			return mockClient, nil
		}
		GetOpenSearchClient()

	})
}

func TestGetOpensearchData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	responseBody := `{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()
	client := opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}
	_, err := GetOpensearchData("query", "1", &client)
	require.Nil(t, err)
}

func TestGetOpensearchDataFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	mockTransport.EXPECT().Perform(gomock.Any()).Return(nil, errors.New("error")).Times(1)
	client := opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}
	_, err := GetOpensearchData("query", "1", &client)
	require.NotNil(t, err)
}

func TestFormResponse(t *testing.T) {
	// Test case 1: Response is nil
	result := formResponse(nil)
	require.Equal(t, "", result)

	// Test case 2: Response body is nil
	response := &opensearchapi.Response{Body: nil}
	result = formResponse(response)
	require.Equal(t, "", result)

	// Test case 3: Error reading response body
	response = &opensearchapi.Response{Body: ioutil.NopCloser(&errorReader{})}
	expectedErrMsg := "<error reading response body:"
	result = formResponse(response)
	require.Contains(t, result, expectedErrMsg)

	// Test case 4: Successful case
	responseBody := "Test response body"
	response = &opensearchapi.Response{Body: ioutil.NopCloser(bytes.NewBufferString(responseBody))}
	result = formResponse(response)
	require.Equal(t, responseBody, result)
}

// errorReader is a mock reader that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading response body")
}

func TestInsertBulkData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200}, nil).Times(1)
	client := opensearch.Client{Transport: mockTransport,
		API: opensearchapi.New(mockTransport),
	}
	err := InsertBulkData(&client, "query")
	require.Nil(t, err, "No Error")

	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 300, Status: "Issue occuring, errors:=true"}, errors.New("error")).Times(1)
	client1 := opensearch.Client{Transport: mockTransport,
		API: opensearchapi.New(mockTransport),
	}
	err = InsertBulkData(&client1, "Query")
	require.Equal(t, err.Error(), "data ingestion failed")

	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 300, Body: io.NopCloser(bytes.NewReader([]byte("\"errors\":true")))}, nil).Times(1)
	client1 = opensearch.Client{Transport: mockTransport,
		API: opensearchapi.New(mockTransport),
	}
	err = InsertBulkData(&client1, "Query")
	require.Nil(t, err)

}

func TestGetMultiQueryData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	responseBody := `{"responses":[{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}]}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()
	client := opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}

	multiQuery := `
	{"index": "index1"}
	{"query": {"match_all": {}}}
	{"index": "index2"}
	{"query": {"match_all": {}}}
	`

	_, err := GetMultiQueryData(&client, multiQuery)
	require.Nil(t, err)
}

func TestGetOpensearchMappingData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	responseBody := `{"index1":{"mappings":{"properties":{"field1":{"type":"text"},"field2":{"type":"integer"}}}}}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()
	client := opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}

	indexName := "index1"
	result, err := GetOpensearchMappingData(indexName, &client)
	require.Nil(t, err)

	expectedResult := `{"index1":{"mappings":{"properties":{"field1":{"type":"text"},"field2":{"type":"integer"}}}}}`
	require.Equal(t, expectedResult, result)
}

func TestGetOpensearchCount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	responseBody := `{"count":10}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()
	client := opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}

	query := `{"query": {"match_all": {}}}`
	indexName := "index1"
	result, err := GetOpensearchCount(query, indexName, &client)
	require.Nil(t, err)

	expectedResult := `{"count":10}`
	require.Equal(t, expectedResult, result)
}

func TestGetWidgetEntity(t *testing.T) {
	widgetId := "e1"
	replacements := map[string]interface{}{
		"e1": "swdelivery/components.json",
	}

	result, err := GetWidgetEntity(widgetId, replacements)
	log.Debugf("Error", err)
	require.NotNil(t, result)
}

func TestGetWidgetEntity_NilWidget(t *testing.T) {
	widgetId := "nilwidget"
	replacements := map[string]interface{}{
		"e1": "swdelivery/components.json",
	}

	result, err := GetWidgetEntity(widgetId, replacements)
	log.Debugf("Error", err)
	require.Nil(t, result)
}

func TestGetComponentComparisonConfig(t *testing.T) {
	widgetId := "e1"
	replacements := map[string]interface{}{
		"e1": "swdelivery/components.json",
	}
	readWidgetConfigFile(widgetId)
	result, err := GetComponentComparisonConfig(widgetId, replacements)
	log.Debugf("Error", err)
	require.NotNil(t, result)
}

func TestGetWidgetEntity_ReplaceWidget(t *testing.T) {
	widgetId := "d3"
	replacements := map[string]interface{}{
		"e1": "dorametrics/failureRate.json",
	}

	result, err := GetWidgetEntity(widgetId, replacements)
	log.Debugf("Error", err)
	require.NotNil(t, result)
}
func TestGetWidgetEntity_3(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockCache := coredataMock.NewMockResourceCacheI(mockCtrl)
	cache.CoreDataResourceCache = mockCache
	mockCache.EXPECT().GetChildrenOfType(gomock.Any(), gomock.Any()).Return([]string{"1"}).Times(2).AnyTimes()
	resource := api.Resource{Id: "Id1", ParentId: "PId1", Name: "R1", Type: api.ResourceType_RESOURCE_TYPE_DASHBOARD}
	mockCache.EXPECT().Get(gomock.Any()).Return(&resource).Times(2).AnyTimes()

	tests := []struct {
		widgetID     string
		replacements map[string]interface{}
	}{
		{
			widgetID: "s3",
			replacements: map[string]interface{}{
				"orgId":     "8509666e-d27f-44fa-46a9-29bc76f5e790",
				"component": []string{"8509666e-d27f-44fa-46a9-29bc76f5e790", "8509444e-d27f-44fa-46a9-29bc76f5e790", "8509333e-d27f-44fa-46a9-29bc76f5e790"},
				"branch":    "main",
			},
		},
		{
			widgetID: "e3",
			replacements: map[string]interface{}{
				"orgId":     "8509666e-d27f-44fa-46a9-29bc76f5e790",
				"component": []string{"8509666e-d27f-44fa-46a9-29bc76f5e790", "8509444e-d27f-44fa-46a9-29bc76f5e790", "8509333e-d27f-44fa-46a9-29bc76f5e790"},
				"branch":    "main",
			},
		},
		{
			widgetID: "cs1",
			replacements: map[string]interface{}{
				"orgId":     "8509666e-d27f-44fa-46a9-29bc76f5e790",
				"component": []string{"8509666e-d27f-44fa-46a9-29bc76f5e790", "8509444e-d27f-44fa-46a9-29bc76f5e790", "8509333e-d27f-44fa-46a9-29bc76f5e790"},
				"branch":    "main",
			},
		},
		{
			widgetID: "cs2",
			replacements: map[string]interface{}{
				"orgId":     "8509666e-d27f-44fa-46a9-29bc76f5e790",
				"component": []string{"8509666e-d27f-44fa-46a9-29bc76f5e790", "8509444e-d27f-44fa-46a9-29bc76f5e790", "8509333e-d27f-44fa-46a9-29bc76f5e790"},
				"branch":    "main",
			},
		},
	}
	for _, tt := range tests {
		we, err := GetWidgetEntity(tt.widgetID, tt.replacements)
		assert.NotNil(t, we)
		assert.Nil(t, err)
	}

}

func TestGetComponentComparisonConfig_Nil(t *testing.T) {
	widgetId := "nil"
	replacements := map[string]interface{}{
		"e1": "swdelivery/components.json",
	}
	readWidgetConfigFile(widgetId)
	result, err := GetComponentComparisonConfig(widgetId, replacements)
	log.Debugf("Error", err)
	require.Nil(t, result)
}

func TestGetReportEntity_Nil(t *testing.T) {
	var reportId int64
	reportId = 0
	result, err := GetReportEntity(reportId)
	log.Debugf("Error", err)
	require.Nil(t, result)
}

func TestGetReportEntity(t *testing.T) {
	var reportId int64
	reportId = 1
	result, err := GetReportEntity(reportId)
	log.Debugf("Error", err)
	require.NotNil(t, result)
}

func TestGetReportEntity_New(t *testing.T) {
	var reportId int64
	reportId = 1
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConfig := mocks.NewMockConfiguration(mockCtrl)
	config.Config = mockConfig
	mockConfig.EXPECT().GetString("report.definition.filepath").Return("db/widget_defintion_map.go").AnyTimes()
	result, err := GetReportEntity(reportId)
	log.Debugf("Error", err)
	require.Nil(t, result)
}

func TestReplaceJSONplaceholders(t *testing.T) {
	widget := WidgetDefinition{
		Id: "widget123",
		W:  10.5,
		H:  20.0,
	}
	query := `{"id": "{{.Id}}", "width": "{{json .W}}", "height": "{{json .H}}"}`
	expected := `{"id": "widget123", "width": "10.5", "height": "20"}`
	result, err := ReplaceJSONplaceholders(widget, query)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Unexpected result. Got: %s, Expected: %s", result, expected)
	}
}
