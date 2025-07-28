package internal

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	api "github.com/calculi-corp/api/go"
	pb "github.com/calculi-corp/api/go/vsm/report"
	"github.com/calculi-corp/common/grpc"
	"github.com/calculi-corp/config"
	coredataMock "github.com/calculi-corp/core-data-cache/mock"
	mock_grpc_client "github.com/calculi-corp/grpc-client/mock"
	testutil "github.com/calculi-corp/grpc-testutil"
	"github.com/calculi-corp/reports-service/cache"
	"github.com/calculi-corp/reports-service/mocks"
	mock "github.com/calculi-corp/repository-service/mock"

	"github.com/calculi-corp/reports-service/constants"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/structpb"
)

func init() {
	config.Config.DefineStringFlag("opensearch.endpoint", "", "The Open search Endpoint")
	config.Config.DefineStringFlag("opensearch.user", "", "The Open search Username")
	config.Config.DefineStringFlag("opensearch.pwd", "", "The Open search password")
	config.Config.DefineStringFlag("report.definition.filepath", "../resources/", "Report Definiftion filepath")
	testutil.SetUnitTestConfig()
	config.Config.Set("logging.level", "INFO")
}

func TestGetData_3(t *testing.T) {
	mockCtx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockTransport := mocks.NewMockTransport(mockCtrl)
	responseBody := `{"responses":[{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}]}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()

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
		name                 string
		widgetId             string
		replacements         map[string]any
		baseDataReplacements map[string]any
		expectedResp         map[string]json.RawMessage
		expectedPastResp     map[string]json.RawMessage
		expectedFuncResp     map[string]json.RawMessage
		expectError          bool
	}{
		{
			name:     "Successful scenario with query response and function execution",
			widgetId: "cs1",
			replacements: map[string]any{
				"orgId":     "8509666e-d27f-44fa-46a9-29bc76f5e790",
				"component": []string{"8509666e-d27f-44fa-46a9-29bc76f5e790", "8509444e-d27f-44fa-46a9-29bc76f5e790", "8509333e-d27f-44fa-46a9-29bc76f5e790"},
				"branch":    "main",
			},
			baseDataReplacements: map[string]interface{}{
				"key1": []string{"value1", "value2"},
				"key2": []string{"value3", "value4"},
			},
			expectedResp: map[string]json.RawMessage{
				"query1": []byte(`{"result": "data1"}`),
				"query2": []byte(`{"result": "data2"}`),
			},
			expectedPastResp: map[string]json.RawMessage{
				"pastQuery1": []byte(`{"result": "pastData1"}`),
			},
			expectedFuncResp: map[string]json.RawMessage{
				"func1": []byte(`{"result": "funcData1"}`),
				"func2": []byte(`{"result": "funcData2"}`),
			},
			expectError: false,
		},
		{
			name:                 "Failure scenario with error from getWidgetQueries",
			widgetId:             "widget2",
			replacements:         map[string]any{},
			baseDataReplacements: map[string]any{},
			expectedResp:         nil,
			expectedPastResp:     nil,
			expectedFuncResp:     nil,
			expectError:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			GetData(tt.widgetId, tt.replacements, tt.baseDataReplacements, mockCtx, mockGrpcClient, nil)

		})
	}
}

func TestCreateComponentComparisonWidget_nil(t *testing.T) {
	type args struct {
		widgetId     string
		data         map[string]json.RawMessage
		replacements map[string]interface{}
		fdata        map[string]json.RawMessage
		organization *constants.Organization
	}

	tests := []struct {
		name    string
		args    args
		want    *pb.Widget
		wantErr bool
	}{
		{
			name: "Failure case :  nil widget",
			args: args{
				widgetId:     "nil",
				data:         map[string]json.RawMessage{},
				replacements: map[string]interface{}{},
				fdata:        map[string]json.RawMessage{},
				organization: &constants.Organization{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateComponentComparisonWidget(tt.args.widgetId, tt.args.data, tt.args.fdata, tt.args.replacements, tt.args.organization)
		})
	}

}

func TestGetComponentComparisonData_2(t *testing.T) {

	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	responseBody := `{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}`

	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	epClient := mock.NewMockEndpointServiceClient(mockCtrl)

	tests := []struct {
		widgetID     string
		replacements map[string]interface{}
	}{
		{
			widgetID: "cs1",
			replacements: map[string]interface{}{
				"orgId":     "8509666e-d27f-44fa-46a9-29bc76f5e790",
				"component": []string{"8509666e-d27f-44fa-46a9-29bc76f5e790", "8509444e-d27f-44fa-46a9-29bc76f5e790", "8509333e-d27f-44fa-46a9-29bc76f5e790"},
				"branch":    "main",
			},
		},
	}

	for _, tt := range tests {
		GetComponentComparisonData(tt.widgetID, tt.replacements, ctx, mockGrpcClient, epClient)
	}
}

func TestGetComponentComparisonData(t *testing.T) {

	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	epClient := mock.NewMockEndpointServiceClient(mockCtrl)

	tests := []struct {
		widgetID     string
		replacements map[string]interface{}
	}{
		{
			widgetID: "cs1",
			replacements: map[string]interface{}{
				"orgId":     "8509666e-d27f-44fa-46a9-29bc76f5e790",
				"component": []string{"8509666e-d27f-44fa-46a9-29bc76f5e790", "8509444e-d27f-44fa-46a9-29bc76f5e790", "8509333e-d27f-44fa-46a9-29bc76f5e790"},
				"branch":    "main",
			},
		},
	}

	for _, tt := range tests {
		GetComponentComparisonData(tt.widgetID, tt.replacements, ctx, mockGrpcClient, epClient)
	}
}

func TestGetComponentComparisonQueries(t *testing.T) {

	tests := []struct {
		widgetID     string
		replacements map[string]interface{}
	}{
		{
			widgetID: "cs1",
			replacements: map[string]interface{}{
				"orgId":     "8509666e-d27f-44fa-46a9-29bc76f5e790",
				"component": []string{"8509666e-d27f-44fa-46a9-29bc76f5e790", "8509444e-d27f-44fa-46a9-29bc76f5e790", "8509333e-d27f-44fa-46a9-29bc76f5e790"},
				"branch":    "main",
			},
		},
	}

	for _, tt := range tests {
		getComponentComparisonQueries(tt.widgetID, tt.replacements)
	}
}

func TestGetComponentComparisonQueries_qn(t *testing.T) {

	tests := []struct {
		widgetID     string
		replacements map[string]interface{}
	}{
		{
			widgetID: "cs3",
			replacements: map[string]interface{}{
				"orgId":     "8509666e-d27f-44fa-46a9-29bc76f5e790",
				"component": []string{"8509666e-d27f-44fa-46a9-29bc76f5e790", "8509444e-d27f-44fa-46a9-29bc76f5e790", "8509333e-d27f-44fa-46a9-29bc76f5e790"},
				"branch":    "main",
			},
		},
	}

	for _, tt := range tests {
		getComponentComparisonQueries(tt.widgetID, tt.replacements)
	}
}

func TestGetWidgetPastDurationQueries(t *testing.T) {

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
	}

	for _, tt := range tests {
		getWidgetPastDurationQueries(tt.widgetID, tt.replacements)
	}
}

func TestGetWidgetPastDurationQueries_error(t *testing.T) {

	widgetId := "mock_widget_id"
	replacements := map[string]interface{}{
		"component": []string{"component_id_1", "component_id_2"},
	}
	getWidgetPastDurationQueries(widgetId, replacements)

}

func TestGetComputedData(t *testing.T) {

	widgetId := "mock_widget_id"
	replacements := map[string]interface{}{
		"component": []string{"component_id_1", "component_id_2"},
	}
	GetComputedData(widgetId, replacements)

}

func TestSetSection(t *testing.T) {
	qr := map[string]json.RawMessage{
		"query_result_key": []byte(`{"mock_key": "mock_value"}`),
	}

	fr_m := map[string]json.RawMessage{
		"function_result_key": []byte(`{"spec_key": "{\"mock_key\": \"mock_value\"}"}`),
		"Test Function":       []byte(`{"mock_key": "mock_value_from_fr"}`),
	}

	replacements := map[string]interface{}{
		"replacement_key": "replacement_value",
	}

	wb := &WidgetBuilder{
		widget: pb.Widget{
			Content: []*pb.Content{
				{
					Section: []*pb.ChartInfo{
						{
							Title:        "Test Title",
							FunctionName: "Test Function",
						},
					},
				},
			},
		},
	}
	wb.setSection(qr, fr_m, replacements)

}

func TestSetSection_nil(t *testing.T) {
	qr := map[string]json.RawMessage{
		"query_result_key": []byte(`{"mock_key": "mock_value"}`),
	}

	fr := map[string]json.RawMessage{
		"function_result_key": []byte(`{"spec_key": "{\"mock_key\": \"mock_value\"}"}`),
	}

	replacements := map[string]interface{}{
		"replacement_key": "replacement_value",
	}

	wb := &WidgetBuilder{
		widget: pb.Widget{
			Content: []*pb.Content{
				{
					Section: []*pb.ChartInfo{
						{
							Title:        "Test Title",
							FunctionName: "Test Function",
						},
					},
				},
			},
		},
	}
	wb.setSection(qr, fr, replacements)

}

func TestSetHeaders(t *testing.T) {
	qr := map[string]json.RawMessage{
		"query_result_key": []byte(`{"mock_key": "mock_value"}`),
	}

	fr := map[string]json.RawMessage{
		"function_result_key": []byte(`{"spec_key": "{\"mock_key\": \"mock_value\"}"}`),
		"Test Function":       []byte(`{"mock_key": "mock_value_from_fr"}`),
	}

	structData := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"mock_key": {
				Kind: &structpb.Value_StringValue{
					StringValue: "mock_value_from_fr",
				},
			},
		},
	}

	structDataBytes, err := json.Marshal(structData)
	if err != nil {
		t.Fatalf("error marshaling struct data: %v", err)
	}
	fr["function_result_key"] = structDataBytes

	wb := &WidgetBuilder{
		widget: pb.Widget{
			Content: []*pb.Content{
				{
					Header: []*pb.MetricInfo{
						{
							Title:        "Test Title",
							Type:         "Test Type",
							FunctionName: "Test Function",
							SpecKey:      "Test Spec Key",
						},
					},
				},
			},
		},
	}
	wb.setHeaders(qr, fr)
	err1 := wb.setHeaders(qr, fr)
	if err1 != nil {
		t.Fatalf("setHeaders returned error: %v", err1)
	}

}

func TestSetHeaders_nil(t *testing.T) {
	qr := map[string]json.RawMessage{
		"query_result_key": []byte(`{"mock_key": "mock_value"}`),
	}

	fr := map[string]json.RawMessage{
		"function_result_key": []byte(`{"spec_key": "{\"mock_key\": \"mock_value\"}"}`),
	}

	wb := &WidgetBuilder{
		widget: pb.Widget{
			Content: []*pb.Content{
				{
					Header: []*pb.MetricInfo{
						{
							Title:        "Test Title",
							Type:         "Test Type",
							FunctionName: "Test Function",
							SpecKey:      "Test Spec Key",
						},
					},
				},
			},
		},
	}
	wb.setHeaders(qr, fr)

}

func TestCombineJSONs(t *testing.T) {
	json1 := []byte(`{"key1": "value1"}`)
	json2 := []byte(`{"key2": "value2"}`)

	combineJSONs(json1, json2)

}

func TestSetData(t *testing.T) {

	qr := map[string]json.RawMessage{
		"Insight Completed Runs Widget": []byte(`{"field1": "value1", "field2": "value2"}`),
	}
	fr := map[string]json.RawMessage{
		"Insight Completed Runs Widget": []byte(`{"field3": "value3", "field4": "value4"}`),
	}

	wb := &WidgetBuilder{
		widget: pb.Widget{
			Data: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"function_name": structpb.NewStringValue("Insight Completed Runs Widget"),
					"query_key":     structpb.NewStringValue("some_query_key"),
					"spec_key":      structpb.NewStringValue("some_spec_key"),
					"other_field":   structpb.NewStringValue("other_value"),
				},
			},
		},
	}

	err := wb.setData(qr, fr)

	assert.NoError(t, err, "Expected no error from setData")

	data := wb.widget.Data
	assert.NotNil(t, data, "Expected widget.Data to be not nil")
	assert.Len(t, data.Fields, 3, "Expected widget.Data.Fields to have 4 fields")

	assert.Nil(t, data.Fields["function_name"], "Expected function_name field to be deleted")
	assert.Nil(t, data.Fields["query_key"], "Expected query_key field to be deleted")
	assert.Nil(t, data.Fields["spec_key"], "Expected spec_key field to be deleted")

	assert.Equal(t, "value3", data.Fields["field3"].GetStringValue())
	assert.Equal(t, "value4", data.Fields["field4"].GetStringValue())
}

func TestDataInfoSeperation(t *testing.T) {
	mockChartInfo := &pb.ChartInfo{
		DataType: 1,
	}

	mockData := []byte(`{"data":[1,2,3],"info":[4,5,6]}`)

	mockDataValue := &structpb.ListValue{}
	mockWidgetBuilder := &WidgetBuilder{
		widget: pb.Widget{
			Content: []*pb.Content{
				{
					Section: []*pb.ChartInfo{
						{
							Info: &structpb.ListValue{},
						},
					},
				},
			},
		},
	}

	_, err := dataInfoSeperation(mockChartInfo, json.RawMessage(mockData), mockDataValue, mockWidgetBuilder, 0, 0)
	assert.NoError(t, err, "Expected no error from dataInfoSeperation")

}
func TestUnmarshalProtoJSON(t *testing.T) {
	mockData := make([]json.RawMessage, 3)
	mockData[0] = json.RawMessage(`{"key1": "value1"}`)
	mockData[1] = json.RawMessage(`{"key2": "value2"}`)
	mockData[2] = json.RawMessage(`{"key3": "value3"}`)

	result, err := unmarshalProtoJSON(mockData)

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected a non-nil structpb.Struct")

	assert.Equal(t, "value3", result.GetFields()["key3"].GetStringValue())
}

func TestCreateWidget(t *testing.T) {
	type args struct {
		widgetId         string
		data             map[string]json.RawMessage
		replacements     map[string]any
		replacementsSpec map[string]any
		fdata            map[string]json.RawMessage
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.Widget
		wantErr bool
	}{
		{
			name: "Success case :  e10 widget to cover footer part from create widget",
			args: args{
				widgetId: "e10",
				data: map[string]json.RawMessage{
					"avgDevelopmentHeader":   []byte(`{"took":108,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":9130,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"developmentCycleTime":{"value":{"total":"2d 6h 42m ","value":196920000}}}}`),
					"developmentCycleChart":  []byte(`{"took":85,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":9130,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"developmentCycleTime":{"value":{"coding_time":53.0,"review_time":9.0,"pickup_time":38.0,"deploy_time":0.0}}}}`),
					"developmentTimeFooter":  []byte(`{"took":37,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":9130,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"developmentCycleTime":{"value":{"coding_time":"1d 5h 11m ","pickup_time_in_millis":74911000,"deploy_time_in_millis":0,"coding_time_in_millis":105082000,"review_time":"4h 42m ","pickup_time":"20h 48m ","review_time_in_millis":16927000,"deploy_time":""}}}}`),
					"p_avgDevelopmentHeader": []byte(`{"took":1,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},"hits":{"total":{"value":0,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"components_addition":{"value":0.0}}}`),
				},
				replacements: map[string]any{
					"aggrBy":           "week",
					"component":        []string{"All"},
					"dateHistogramMin": "2023-10-01",
					"dateHistogramMax": "2023-10-20",
					"duration":         "month",
					"orgId":            "2cab10cc-cd9d-11ed-afa1-0242ac120002",
					"endDate":          "2023-10-31 23:59:11",
					"startDate":        "2023-10-01 00:00:12",
					"metricName":       "p_avgDevelopmentHeader",
					"p_endDate":        "2023-10-31 23:59:11",
					"p_startDate":      "2023-10-01 00:00:12",
				},
				replacementsSpec: map[string]any{
					"normalizeMonthInSpec": "2023-10-01",
				},
				fdata: map[string]json.RawMessage{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateWidget(tt.args.widgetId, tt.args.data, tt.args.replacements, tt.args.replacementsSpec, tt.args.fdata)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.HasPrefix(tt.name, "Success case :  e3") {
				assert.Equal(t, got.Content[0].Header[0].DrillDown.ReportTitle, "Workflow runs", "Validating success response of create widget")
			} else if strings.HasPrefix(tt.name, "Success case :  e10") {
				assert.Equal(t, got.Content[0].Header[0].Title, "Average Development Cycle Time", "Validating success response of create widget")
			}
		})
	}
}

func Test_CreateWidget02(t *testing.T) {
	t.Run("Case 1: Successful creation of widget when function returns data (fdata)", func(t *testing.T) {

		replacements := map[string]any{
			"aggrBy":           "week",
			"component":        []string{"All"},
			"dateHistogramMin": "2023-10-01",
			"dateHistogramMax": "2023-10-20",
			"duration":         "month",
			"orgId":            "2cab10cc-cd9d-11ed-afa1-0242ac120002",
			"endDate":          "2023-10-31 23:59:11",
			"startDate":        "2023-10-01 00:00:12",
		}

		fdata := map[string]json.RawMessage{
			"Component Widget Header":  []byte(`{"subTitle":{"title":"repositories","value":1},"value":1}`),
			"Component Widget Section": []byte(`{"data":[{"name":"Active","value":100},{"name":"Inactive","value":0}],"info":[{"drillDown":{"reportId":"component","reportTitle":"Components","reportType":"status"},"title":"Active","value":1},{"drillDown":{"reportId":"component","reportTitle":"Components","reportType":"status"},"title":"Inactive","value":0}]}`),
		}

		got, err := CreateWidget("e1", nil, replacements, map[string]any{}, fdata)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualData, err := json.Marshal(got)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		expectedData := `{"id":"e1","title":"Components","description":"Track both active and inactive components over the selected time frame.","content":[{"header":[{"title":"Total components","description":"Total components","data":{"subTitle":{"title":"repositories","value":1},"value":1},"drill_down":{"report_id":"component","report_title":"Components"}}],"section":[{"type":6,"color_scheme":[{"color0":"#458CD1","color1":"#0D82F6"},{"color0":"#8BA9B8","color1":"#5D7689"}],"data":[{"name":"Active","value":100},{"name":"Inactive","value":0}],"info":[{"drillDown":{"reportId":"component","reportTitle":"Components","reportType":"status"},"title":"Active","value":1},{"drillDown":{"reportId":"component","reportTitle":"Components","reportType":"status"},"title":"Inactive","value":0}],"data_type":1,"drill_down":{"report_id":"component","report_type":"status","report_title":"Components"},"light_color_scheme":[{"color0":"#35B1F9","color1":"#0781D2"},{"color0":"#AEBEC5","color1":"#7591A2"}]}]}],"enable_components_compare":true,"components_compare_id":"components-compare"}`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedData, string(actualData), "Validating transformed data for transformTestInsightsTotalRunsDrillDown")

	})

	t.Run("Case 2: Error when function returns data but the key in the map doesn't match the function name in the widget definition", func(t *testing.T) {

		replacements := map[string]any{
			"aggrBy":           "week",
			"component":        []string{"All"},
			"dateHistogramMin": "2023-10-01",
			"dateHistogramMax": "2023-10-20",
			"duration":         "month",
			"orgId":            "2cab10cc-cd9d-11ed-afa1-0242ac120002",
			"endDate":          "2023-10-31 23:59:11",
			"startDate":        "2023-10-01 00:00:12",
		}

		fdata := map[string]json.RawMessage{
			"Component Widget":         []byte(`{"subTitle":{"title":"repositories","value":1},"value":1}`),
			"Component Widget Section": []byte(`{"data":[{"name":"Active","value":100},{"name":"Inactive","value":0}],"info":[{"drillDown":{"reportId":"component","reportTitle":"Components","reportType":"status"},"title":"Active","value":1},{"drillDown":{"reportId":"component","reportTitle":"Components","reportType":"status"},"title":"Inactive","value":0}]}`),
		}

		_, err := CreateWidget("e1", nil, replacements, map[string]any{}, fdata)

		// Assertions to check if the result is as expected.
		assert.NotNil(t, err, "Validating transformed data for transformTestInsightsTotalRunsDrillDown")

	})

}

func TestGetData(t *testing.T) {
	mockCtx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	tests := []struct {
		name                 string
		widgetId             string
		replacements         map[string]any
		baseDataReplacements map[string]any
		expectedResp         map[string]json.RawMessage
		expectedPastResp     map[string]json.RawMessage
		expectedFuncResp     map[string]json.RawMessage
		expectError          bool
	}{
		{
			name:                 "Successful scenario with query response and function execution",
			widgetId:             "widget1",
			replacements:         map[string]any{},
			baseDataReplacements: map[string]any{},
			expectedResp: map[string]json.RawMessage{
				"query1": []byte(`{"result": "data1"}`),
				"query2": []byte(`{"result": "data2"}`),
			},
			expectedPastResp: map[string]json.RawMessage{
				"pastQuery1": []byte(`{"result": "pastData1"}`),
			},
			expectedFuncResp: map[string]json.RawMessage{
				"func1": []byte(`{"result": "funcData1"}`),
				"func2": []byte(`{"result": "funcData2"}`),
			},
			expectError: false,
		},
		{
			name:                 "Failure scenario with error from getWidgetQueries",
			widgetId:             "widget2",
			replacements:         map[string]any{},
			baseDataReplacements: map[string]any{},
			expectedResp:         nil,
			expectedPastResp:     nil,
			expectedFuncResp:     nil,
			expectError:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			GetData(tt.widgetId, tt.replacements, tt.baseDataReplacements, mockCtx, mockClient, nil)

		})
	}
}

func TestGetWidgetQueries(t *testing.T) {
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
		{
			widgetID: "aso2",
			replacements: map[string]interface{}{
				"orgId":     "8509666e-d27f-44fa-46a9-29bc76f5e790",
				"application_id": []string{"app1"},
			},
		},
	}

	for _, tt := range tests {
		_, _, err := getWidgetQueries(tt.widgetID, tt.replacements, nil)
		assert.Nil(t, err)
	}
}

func TestSortCompareReports(t *testing.T) {
	testCases := []struct {
		name          string
		input         []*pb.CompareReports
		expectedOrder []*pb.CompareReports
	}{
		{
			name: "Sorting with IsSubOrg condition",
			input: []*pb.CompareReports{
				{IsSubOrg: true},
				{IsSubOrg: false},
				{IsSubOrg: true},
				{IsSubOrg: false},
			},
			expectedOrder: []*pb.CompareReports{
				{IsSubOrg: false},
				{IsSubOrg: false},
				{IsSubOrg: true},
				{IsSubOrg: true},
			},
		},
		{
			name:          "Empty input",
			input:         []*pb.CompareReports{},
			expectedOrder: []*pb.CompareReports{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputCopy := make([]*pb.CompareReports, len(tc.input))
			copy(inputCopy, tc.input)
			sortCompareReports(inputCopy)
		})
	}
}

func TestSortCompareReportsAlphabetically(t *testing.T) {
	mockCompareReports := []*pb.CompareReports{
		{CompareTitle: "D", IsSubOrg: false},
		{CompareTitle: "A", IsSubOrg: true},
		{CompareTitle: "C", IsSubOrg: true},
		{CompareTitle: "B", IsSubOrg: true},
	}
	expectedOrder := []*pb.CompareReports{
		{CompareTitle: "A", IsSubOrg: true},
		{CompareTitle: "B", IsSubOrg: true},
		{CompareTitle: "C", IsSubOrg: true},
		{CompareTitle: "D", IsSubOrg: false},
	}
	sortCompareReportsAlphabetically(mockCompareReports)
	if !compareReportsSliceEqual(mockCompareReports, expectedOrder) {
		t.Errorf("SortCompareReportsAlphabetically() did not sort as expected.\nExpected: %+v\nGot: %+v", expectedOrder, mockCompareReports)
	}
}

func compareReportsSliceEqual(a, b []*pb.CompareReports) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !compareReportsEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func compareReportsEqual(a, b *pb.CompareReports) bool {
	return a.IsSubOrg == b.IsSubOrg && a.CompareTitle == b.CompareTitle
}
