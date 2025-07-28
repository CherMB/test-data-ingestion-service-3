package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	pb "github.com/calculi-corp/api/go/vsm/report"
	mock "github.com/calculi-corp/repository-service/mock"

	api "github.com/calculi-corp/api/go"
	"github.com/calculi-corp/api/go/endpoint"
	"github.com/calculi-corp/common/grpc"
	client "github.com/calculi-corp/grpc-client"
	mock_grpc_client "github.com/calculi-corp/grpc-client/mock"
	"github.com/calculi-corp/reports-service/constants"
	"github.com/calculi-corp/reports-service/models"
	"github.com/opensearch-project/opensearch-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_getInsightSystemInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
		epClt        endpoint.EndpointServiceClient
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "SystemInformation success case",
			args: args{
				widgetId: "ci2",
				replacements: map[string]any{
					"aggrBy":           "week",
					"ciToolId":         "1edebf0d-6797-4ec0-9a50-fac728645a0e",
					"ciToolType":       "CBCI",
					"dateHistogramMax": "2023-12-30",
					"dateHistogramMin": "2023-12-01",
					"duration":         "month",
					"endDate":          "2023-12-30 23:59:59",
					"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
					"startDate":        "2023-12-01 00:00:00",
				},
			},
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	responseString := `{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":0.0,"hits":[{"_index":"cb_ci_tool_insight","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_0edebf0d-6797-4ec0-9a50-fac728645a0e_CJOC Test_CJOC","_score":0.0,"_source":{"plugins":[{"requiredCoreVersion":"2.346.3","active":true,"shortName":"metrics","version":"4.2.13-420.vea_2f17932dd6","enabled":true,"longName":"Metrics Plugin","dependencies":[{"shortName":"ionicons-api","version":"31.v4757b_6987003"},{"shortName":"jackson2-api","version":"2.13.4.20221013-295.v8e29ea_354141"},{"shortName":"variant","version":"59.vf075fe829ccb"},{"optional":true,"shortName":"instance-identity","version":"3.1"}]},{"requiredCoreVersion":"2.361.1","active":true,"shortName":"display-url-api","version":"2.3.7","enabled":true,"longName":"Display URL API"}],"endpoint_id":"0edebf0d-6797-4ec0-9a50-fac728645a0e","created_at":"2023-12-30 15:07:54","type":"CJOC","version":"2.387.2.3","url":"https://cjoc.rosaas.releaseiq.io/","users":[{"name":"noreply","id":"noreply","type":"COMMITTER","email":"noreply@github.com"},{"name":"SYSTEM","id":"SYSTEM","type":"COMMITTER"},{"name":"releaseiq","id":"releaseiq","type":"ACCESS_USER","email":"sshaik@cloudbees.com"}],"latest_version":"","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","system_health":[{"healthy":true,"name":"plugins","message":"No failed plugins","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"thread-deadlock","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"disk-space","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"temporary-space","timestamp":"2023-12-07T12:47:28.128Z"}],"name":"CJOC Test","metrics":[{"metricsType":"gauges","metricsData":{"jenkins.job.count.value":{"value":27},"jenkins.queue.size.value":{},"jenkins.project.enabled.count.value":{"value":27},"jenkins.executor.in-use.value":{"value":1},"jenkins.node.offline.value":{},"jenkins.queue.stuck.value":{},"jenkins.executor.count.value":{"value":2},"jenkins.executor.free.value":{"value":1},"jenkins.project.count.value":{"value":27},"jenkins.queue.pending.value":{},"jenkins.project.disabled.count.value":{},"jenkins.queue.buildable.value":{},"jenkins.node.online.value":{},"jenkins.queue.blocked.value":{},"jenkins.node.count.value":{"value":1}}}],"org_name":"cloudbees-staging","timestamp":"2023-12-30 15:07:54"}}]}}`
	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		return responseString, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
				return &endpoint.EndpointsResponse{
					Endpoints: []*endpoint.Endpoint{
						{
							Id:   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
							Name: "CJOC Test",
						},
						{
							Id:             "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
							Name:           "CB1 CBCI Test",
							ContributionId: "cjoc_cbci-app-endpoint-type",
							Properties: []*api.Property{
								{
									Name: "status",
									Value: &api.Property_String_{
										String_: "INSTALLED",
									},
								},
								{
									Name: constants.CJOC_ENDPOINT_ID,
									Value: &api.Property_String_{
										String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
									},
								},
							},
						},
						{
							Id:             "a7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
							Name:           "CBCI Test",
							ContributionId: "cjoc_cbci-app-endpoint-type",
							Properties: []*api.Property{
								{
									Name: "status",
									Value: &api.Property_String_{
										String_: "NOT_INSTALLED",
									},
								},
								{
									Name: constants.CJOC_ENDPOINT_ID,
									Value: &api.Property_String_{
										String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
									},
								},
							},
						},
					},
				}, nil

			}
			got, err := getInsightSystemInformation(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, tt.args.epClt)
			responseMap := map[string]interface{}{}
			json.Unmarshal(got, &responseMap)

			if tt.name == "SystemInformation success case" {
				assert.Equal(t, responseMap["freeExecutors"].(float64), float64(1), "Validating Executor count")
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getInsightSystemInformation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_getInsightSystemHealth(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	replacement := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "1edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType":       "CBCI",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
	}
	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
		epClt        endpoint.EndpointServiceClient
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "SystemHealth success case",
			args: args{
				widgetId:     "ci3",
				replacements: replacement,
			},
		},
		{
			name: "SystemHealth warning case",
			args: args{
				widgetId:     "ci3",
				replacements: replacement,
			},
		},
		{
			name: "SystemHealth failure case",
			args: args{
				widgetId:     "ci3",
				replacements: replacement,
			},
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	responseString := `{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":0.0,"hits":[{"_index":"cb_ci_tool_insight","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_0edebf0d-6797-4ec0-9a50-fac728645a0e_CJOC Test_CJOC","_score":0.0,"_source":{"endpoint_id":"0edebf0d-6797-4ec0-9a50-fac728645a0e","created_at":"2023-12-30 15:07:54","type":"CJOC","version":"2.387.2.3","url":"https://cjoc.rosaas.releaseiq.io/","latest_version":"","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","system_health":[{"healthy":true,"name":"plugins","message":"No failed plugins","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"thread-deadlock","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"disk-space","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"temporary-space","timestamp":"2023-12-07T12:47:28.128Z"}],"name":"CJOC Test","org_name":"cloudbees-staging","timestamp":"2023-12-30 15:07:54"}}]}}`
	responseString1 := `{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":0.0,"hits":[{"_index":"cb_ci_tool_insight","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_0edebf0d-6797-4ec0-9a50-fac728645a0e_CJOC Test_CJOC","_score":0.0,"_source":{"endpoint_id":"0edebf0d-6797-4ec0-9a50-fac728645a0e","created_at":"2023-12-30 15:07:54","type":"CJOC","version":"2.387.2.3","url":"https://cjoc.rosaas.releaseiq.io/","latest_version":"","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","system_health":[{"healthy":true,"name":"plugins","message":"No failed plugins","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":false,"name":"thread-deadlock","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"disk-space","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"temporary-space","timestamp":"2023-12-07T12:47:28.128Z"}],"name":"CJOC Test","org_name":"cloudbees-staging","timestamp":"2023-12-30 15:07:54"}}]}}`
	responseString2 := `{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":0.0,"hits":[{"_index":"cb_ci_tool_insight","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_0edebf0d-6797-4ec0-9a50-fac728645a0e_CJOC Test_CJOC","_score":0.0,"_source":{"endpoint_id":"0edebf0d-6797-4ec0-9a50-fac728645a0e","created_at":"2023-12-30 15:07:54","type":"CJOC","version":"2.387.2.3","url":"https://cjoc.rosaas.releaseiq.io/","latest_version":"","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","system_health":[{"healthy":false,"name":"plugins","message":"No failed plugins","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":false,"name":"thread-deadlock","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"disk-space","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":false,"name":"temporary-space","timestamp":"2023-12-07T12:47:28.128Z"}],"name":"CJOC Test","org_name":"cloudbees-staging","timestamp":"2023-12-30 15:07:54"}}]}}`
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
				return &endpoint.EndpointsResponse{
					Endpoints: []*endpoint.Endpoint{
						{
							Id:   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
							Name: "CJOC Test",
						},
						{
							Id:             "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
							Name:           "CB1 CBCI Test",
							ContributionId: "cjoc_cbci-app-endpoint-type",
							Properties: []*api.Property{
								{
									Name: "status",
									Value: &api.Property_String_{
										String_: "INSTALLED",
									},
								},
								{
									Name: constants.CJOC_ENDPOINT_ID,
									Value: &api.Property_String_{
										String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
									},
								},
							},
						},
						{
							Id:             "a7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
							Name:           "CBCI Test",
							ContributionId: "cjoc_cbci-app-endpoint-type",
							Properties: []*api.Property{
								{
									Name: "status",
									Value: &api.Property_String_{
										String_: "NOT_INSTALLED",
									},
								},
								{
									Name: constants.CJOC_ENDPOINT_ID,
									Value: &api.Property_String_{
										String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
									},
								},
							},
						},
					},
				}, nil

			}
			searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
				if tt.name == "SystemHealth warning case" {
					return responseString1, nil
				} else if tt.name == "SystemHealth failure case" {
					return responseString2, nil
				}
				return responseString, nil
			}
			got, err := getInsightSystemHealth(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, tt.args.epClt)
			responseMap := map[string]interface{}{}
			json.Unmarshal(got, &responseMap)
			if tt.name == "SystemHealth success case" {
				assert.Equal(t, responseMap["healthScore"].(string), "100%", "Validating system health score count")
				assert.Equal(t, responseMap["healthStatus"].(string), "success", "Validating system health status")
			} else if tt.name == "SystemHealth warning case" {
				assert.Equal(t, responseMap["healthScore"].(string), "75%", "Validating system health score count")
				assert.Equal(t, responseMap["healthStatus"].(string), "warning", "Validating system health status")
			} else if tt.name == "SystemHealth failure case" {
				assert.Equal(t, responseMap["healthScore"].(string), "25%", "Validating system health score count")
				assert.Equal(t, responseMap["healthStatus"].(string), "failed", "Validating system health status")
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getInsightSystemHealth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type MockStream struct {
	pb.ReportServiceHandler_StreamCIInsightsCompletedRunServer
	ctx     context.Context
	sendErr error
}

func (m *MockStream) Send(response *pb.StreamCIInsightsCompletedRunsResponse) error {
	return m.sendErr
}

func Test_GetInsightCompletedRunsStream(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointClient := mock.NewMockEndpointServiceClient(mockCtrl)
	mockServer := &MockStream{
		ctx: context.Background(),
	}

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	type args struct {
		replacements map[string]any
		ctx          context.Context
		epClt        endpoint.EndpointServiceClient
		srv          pb.ReportServiceHandler_StreamCIInsightsCompletedRunServer
	}
	replacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "8e78c2f7-6082-4cfd-b5f3-575ba92a4d4e",
		"ciToolType":       "CJOC",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
		"timeZone":         "Asia/Calcutta",
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		if IndexName == constants.CB_CI_JOB_INFO_INDEX {
			return `{"aggregations":{"jobs":{"value":[{"job_name":"Demo_Folder/job/Demo_Folder_Pipeline/","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"1d7db1eb-8d24-4f01-9df3-8afbd7734710","last_completed_run_id":12,"endpoint_id":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","type":"Pipeline","display_name":"Demo_Folder/Demo_Folder_Pipeline"}]}}}`, nil
		}

		return `{"aggregations":{"completedRuns":{"value":[{"result":"FAILURE","duration":51440,"start_time":"2023-12-24T05:50:00.000Z","run_id":4,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"1d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","start_time_in_millis":1703377200000,"timestamp":"2023-12-30T07:45:16.000Z"},{"result":"SUCCESS","duration":51440,"start_time":"2023-12-25T05:50:00.000Z","run_id":5,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"1d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","start_time_in_millis":1703483400000,"timestamp":"2023-12-30T07:45:16.000Z"},{"result":"SUCCESS","duration":51440,"start_time":"2023-12-20T05:50:00.000Z","run_id":2,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"1d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","start_time_in_millis":1703031600000,"timestamp":"2023-12-30T07:45:16.000Z"},{"result":"FAILURE","duration":51440,"start_time":"2023-12-22T05:50:00.000Z","run_id":3,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"1d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","start_time_in_millis":1703204400000,"timestamp":"2023-12-30T07:45:16.000Z"}]}}}`, nil
	}

	endpointt := &endpoint.Endpoint{
		ResourceId: "2cab10cc-cd9d-11ed-afa1-0242ac120002",
	}
	endpoints := []*endpoint.Endpoint{endpointt}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpoints,
	}

	mockEndpointClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil).AnyTimes()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Get completed runs stream",
			args: args{
				replacements: replacements,
				ctx:          context.Background(),
				epClt:        mockEndpointClient,
				srv:          mockServer,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetInsightCompletedRunsStream(tt.args.replacements, tt.args.ctx, tt.args.epClt, tt.args.srv)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err, "unexpected end of JSON input")
				return
			}
		})
	}
}

func Test_getInsightCompletedRuns(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
		epClt        endpoint.EndpointServiceClient
	}
	replacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "0edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType":       "CJOC",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
		"timeZone":         "Asia/Calcutta",
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "Completed runs success case - Non CJOC",
			args: args{
				widgetId: "ci5",
				replacements: map[string]any{
					"aggrBy":           "week",
					"ciToolId":         "1edebf0d-6797-4ec0-9a50-fac728645a0e",
					"ciToolType":       "CBCI",
					"dateHistogramMax": "2023-12-30",
					"dateHistogramMin": "2023-12-01",
					"duration":         "month",
					"endDate":          "2023-12-30 23:59:59",
					"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
					"startDate":        "2023-12-01 00:00:00",
					"timeZone":         "Asia/Calcutta",
				},
			},
		},
		{
			name: "Completed runs success case - CJOC",
			args: args{
				widgetId:     "ci5",
				replacements: replacements,
			},
		},
		{
			name: "Completed runs success case without cjoc endpoint entry - CJOC",
			args: args{
				widgetId:     "ci5",
				replacements: replacements,
			},
		},
		{
			name: "Completed runs failure case - CJOC",
			args: args{
				widgetId: "ci5",
				replacements: map[string]any{
					"aggrBy":           "week",
					"ciToolId":         "0edebf0d-6797-4ec0-9a50-fac728645a0e",
					"ciToolType":       "CJOC",
					"dateHistogramMax": "2023-12-30",
					"dateHistogramMin": "2023-12-01",
					"duration":         "month",
					"endDate":          "2023-12-30 23:59:59",
					"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
					"startDate":        "2023-12-01 00:00:00",
					"timeZone":         "Asia/Calcutta",
				},
			},
			wantErr: true,
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
		return &endpoint.EndpointsResponse{
			Endpoints: []*endpoint.Endpoint{
				{
					Id:   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
					Name: "CJOC Test",
				},
				{
					Id:             "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
					Name:           "CB1 CBCI Test",
					ContributionId: "cjoc_cbci-app-endpoint-type",
					Properties: []*api.Property{
						{
							Name: "status",
							Value: &api.Property_String_{
								String_: "INSTALLED",
							},
						},
						{
							Name: constants.CJOC_ENDPOINT_ID,
							Value: &api.Property_String_{
								String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
							},
						},
					},
				},
				{
					Id:             "a7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
					Name:           "CBCI Test",
					ContributionId: "cjoc_cbci-app-endpoint-type",
					Properties: []*api.Property{
						{
							Name: "status",
							Value: &api.Property_String_{
								String_: "NOT_INSTALLED",
							},
						},
						{
							Name: constants.CJOC_ENDPOINT_ID,
							Value: &api.Property_String_{
								String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
							},
						},
					},
				},
			},
		}, nil

	}
	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		if IndexName == constants.CB_CI_JOB_INFO_INDEX {
			return `{"aggregations":{"jobs":{"value":[{"job_name":"Demo_Folder/job/Demo_Folder_Pipeline/","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"1d7db1eb-8d24-4f01-9df3-8afbd7734710","last_completed_run_id":12,"endpoint_id":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","type":"Pipeline","display_name":"Demo_Folder/Demo_Folder_Pipeline"}]}}}`, nil
		} else if strings.Contains(query, "result_counts") {
			return `{"took":1,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":640,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"result_counts":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"SUCCESS","doc_count":296},{"key":"ABORTED","doc_count":161},{"key":"FAILURE","doc_count":117},{"key":"UNSTABLE","doc_count":44},{"key":"NOT_BUILT","doc_count":22}]}}}`, nil
		}
		return `{"aggregations":{"completedRuns":{"value":[{"result":"FAILURE","duration":51440,"start_time":"2023-12-24T05:50:00.000Z","run_id":4,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"1d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","start_time_in_millis":1703377200000,"timestamp":"2023-12-30T07:45:16.000Z"},{"result":"SUCCESS","duration":51440,"start_time":"2023-12-25T05:50:00.000Z","run_id":5,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"1d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","start_time_in_millis":1703483400000,"timestamp":"2023-12-30T07:45:16.000Z"},{"result":"SUCCESS","duration":51440,"start_time":"2023-12-20T05:50:00.000Z","run_id":2,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"1d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","start_time_in_millis":1703031600000,"timestamp":"2023-12-30T07:45:16.000Z"},{"result":"FAILURE","duration":51440,"start_time":"2023-12-22T05:50:00.000Z","run_id":3,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"1d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","start_time_in_millis":1703204400000,"timestamp":"2023-12-30T07:45:16.000Z"}]}}}`, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Completed runs success case - CJOC" {

				getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
					return &endpoint.EndpointsResponse{
						Endpoints: []*endpoint.Endpoint{
							{
								Id:   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
								Name: "CJOC Test",
							},
							{
								Id:             "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
								Name:           "CB1 CBCI Test",
								ContributionId: "cjoc_cbci-app-endpoint-type",
								Properties: []*api.Property{
									{
										Name: "status",
										Value: &api.Property_String_{
											String_: "INSTALLED",
										},
									},
									{
										Name: constants.CJOC_ENDPOINT_ID,
										Value: &api.Property_String_{
											String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
										},
									},
								},
							},
							{
								Id:             "a7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
								Name:           "CBCI Test",
								ContributionId: "cjoc_cbci-app-endpoint-type",
								Properties: []*api.Property{
									{
										Name: "status",
										Value: &api.Property_String_{
											String_: "NOT_INSTALLED",
										},
									},
									{
										Name: constants.CJOC_ENDPOINT_ID,
										Value: &api.Property_String_{
											String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
										},
									},
								},
							},
						},
					}, nil

				}
			} else if tt.name == "Completed runs success case without cjoc endpoint entry - CJOC" {
				getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
					return &endpoint.EndpointsResponse{
						Endpoints: []*endpoint.Endpoint{},
					}, nil
				}
			} else if tt.name == "Completed runs failure case - CJOC" {
				getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
					return nil, errors.New("error")
				}
			}
			got, err := getInsightCompletedRuns(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, tt.args.epClt)
			responseMap := map[string]interface{}{}
			json.Unmarshal(got, &responseMap)
			if tt.name == "Completed runs success case - Non CJOC" {
				assert.Equal(t, len(responseMap["data"].([]interface{})), 4, "Validating data count")
				assert.Equal(t, len(responseMap[constants.COUNT_INFO].([]interface{})), 5, "Validating data count")
			} else if tt.name == "Completed runs success case - CJOC" {
				assert.Equal(t, len(responseMap["data"].([]interface{})), 4, "Validating data count")
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getInsightCompletedRuns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_getInsightProjectsActivity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
		epClt        endpoint.EndpointServiceClient
	}
	replacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "0edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType":       "CJOC",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
	}
	cbciReplacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "1edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType":       "CBCI",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
	}
	cbciIdleReplacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "1edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType":       "CBCI",
		"filterType":       "IdleFilter",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
	}
	cbciFragileReplacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "1edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType":       "CBCI",
		"filterType":       "FragileFilter",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "Project activity success case - Non CJOC",
			args: args{
				widgetId:     "ci7",
				replacements: cbciReplacements,
			},
		},
		{
			name: "Project activity success case - CJOC",
			args: args{
				widgetId:     "ci7",
				replacements: replacements,
			},
		},
		{
			name: "Project activity success case without cjoc endpoint entry - CJOC",
			args: args{
				widgetId:     "ci7",
				replacements: replacements,
			},
		},
		{
			name: "Project activity success case - Idle filter",
			args: args{
				widgetId:     "ci7",
				replacements: cbciIdleReplacements,
			},
		},
		{
			name: "Project activity success case - Fragile filter",
			args: args{
				widgetId:     "ci7",
				replacements: cbciFragileReplacements,
			},
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
		return &endpoint.EndpointsResponse{
			Endpoints: []*endpoint.Endpoint{
				{
					Id:   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
					Name: "CJOC Test",
				},
				{
					Id:             "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
					Name:           "CB1 CBCI Test",
					ContributionId: "cjoc_cbci-app-endpoint-type",
					Properties: []*api.Property{
						{
							Name: "status",
							Value: &api.Property_String_{
								String_: "INSTALLED",
							},
						},
						{
							Name: constants.CJOC_ENDPOINT_ID,
							Value: &api.Property_String_{
								String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
							},
						},
					},
				},
				{
					Id:             "a7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
					Name:           "CBCI Test",
					ContributionId: "cjoc_cbci-app-endpoint-type",
					Properties: []*api.Property{
						{
							Name: "status",
							Value: &api.Property_String_{
								String_: "NOT_INSTALLED",
							},
						},
						{
							Name: constants.CJOC_ENDPOINT_ID,
							Value: &api.Property_String_{
								String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
							},
						},
					},
				},
			},
		}, nil

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Project activity success case - CJOC" {
				getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
					return &endpoint.EndpointsResponse{
						Endpoints: []*endpoint.Endpoint{
							{
								Id:   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
								Name: "CJOC Test",
							},
							{
								Id:             "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
								Name:           "CB1 CBCI Test",
								ContributionId: "cjoc_cbci-app-endpoint-type",
								Properties: []*api.Property{
									{
										Name: "status",
										Value: &api.Property_String_{
											String_: "INSTALLED",
										},
									},
									{
										Name: constants.CJOC_ENDPOINT_ID,
										Value: &api.Property_String_{
											String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
										},
									},
								},
							},
						},
					}, nil

				}
			} else if tt.name == "Project activity success case without cjoc endpoint entry - CJOC" {
				getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
					return &endpoint.EndpointsResponse{
						Endpoints: []*endpoint.Endpoint{},
					}, nil
				}
			} else if tt.name == "Project activity failure case - CJOC" {
				getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
					return nil, errors.New("error")
				}
			}
			count := 0
			searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
				if IndexName == constants.CB_CI_JOB_INFO_INDEX {
					return `{"took":6,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":328,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"total_hits":{"value":328},"completedRuns":{"doc_count_error_upper_bound":8,"sum_other_doc_count":223,"buckets":[{"key":"592b7c18-6d13-4b7e-9d8e-67fadec3b112","doc_count":105,"result_buckets":{"buckets":{"aborted":{"doc_count":1,"total_duration":{"value":4879},"last_active":{"value":1706511089334}},"failure":{"doc_count":0,"total_duration":{"value":0},"last_active":{"value":0}},"success":{"doc_count":104,"total_duration":{"value":10958926},"last_active":{"value":1706784092335}},"unstable":{"doc_count":0,"total_duration":{"value":0},"last_active":{"value":0}}}},"endpoint_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"a258baa2-f1b3-45e6-874b-0f4d0fe1e459","doc_count":105}]}}]}}}`, nil
				} else if query == `{"aggs":{"completedRuns":{"scripted_metric":{"combine_script":"return state.dataMap;","init_script":"state.dataMap = [:];","map_script":"def map = state.dataMap; def key = doc.org_id.value + '_' + doc.endpoint_id.value + '_' + doc.job_id.value + '_'+doc.run_id; def v = ['orgId': doc.org_id.value, 'endpointId': doc.endpoint_id.value, 'jobId': doc.job_id.value,'runId': doc.run_id.value, 'result': doc.result.value,'duration': doc.duration.value, 'startTime': doc.start_time.value, 'startTimeInMillis': doc.start_time_in_millis.value,'timestamp':doc.timestamp.value]; map.put(key, v);","reduce_script":"def tmpMap = [: ], resultList = new ArrayList();for (response in states) {if (response != null) {for (key in response.keySet()) {tmpMap.put(key, response.get(key));}}}def jobIds = new HashSet();for (key in tmpMap.keySet()) {def record = tmpMap.get(key);jobIds.add(record.jobId);}return jobIds;"}}},"query":{"bool":{"filter":[{"range":{"start_time":{"format":"yyyy-MM-dd HH:mm:ss","gte":"2023-12-01 00:00:00","lte":"2023-12-30 23:59:59"}}},{"term":{"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002"}}],"must":[{"terms":{"endpoint_id":["b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","4be41c9e-bec3-4f5f-82b9-7483f0eab37d"]}}],"must_not":{"term":{"result":"IN_PROGRESS"}}}},"size":0}` {
					return `{"aggregations":{"completedRuns":{"value":["1d7db1eb-8d24-4f01-9df3-8afbd7734710","127ab988-6422-44ba-a08e-a5a35dd65d89"]}}}`, nil
				} else if tt.name == "Project activity success case - Idle filter" && count == 0 {
					count++
					return `{"aggregations":{"completedRuns":{"value":["1d7db1eb-8d24-4f01-9df3-8afbd7734710","127ab988-6422-44ba-a08e-a5a35dd65d89"]}}}`, nil
				} else if tt.name == "Project activity success case - Fragilefilter" {
					return `{"took":6,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":328,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"total_hits":{"value":328},"completedRuns":{"doc_count_error_upper_bound":8,"sum_other_doc_count":223,"buckets":[{"key":"592b7c18-6d13-4b7e-9d8e-67fadec3b112","doc_count":105,"result_buckets":{"buckets":{"aborted":{"doc_count":1,"total_duration":{"value":4879},"last_active":{"value":1706511089334}},"failure":{"doc_count":0,"total_duration":{"value":0},"last_active":{"value":0}},"success":{"doc_count":104,"total_duration":{"value":10958926},"last_active":{"value":1706784092335}},"unstable":{"doc_count":0,"total_duration":{"value":0},"last_active":{"value":0}}}},"endpoint_id":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"a258baa2-f1b3-45e6-874b-0f4d0fe1e459","doc_count":105}]}}]}}}`, nil
				} else {
					return `{"aggregations":{"completedRuns":{"value":{"1d7db1eb-8d24-4f01-9df3-8afbd7734710":{"unstable":0,"totalDuration":391415,"result":"SUCCESS","lastActive":1703958651101,"success":3,"aborted":1,"endpointId":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","executed":8,"failed":4},"127ab988-6422-44ba-a08e-a5a35dd65d89":{"unstable":0,"totalDuration":1844,"result":"SUCCESS","lastActive":1703948386425,"success":1,"aborted":0,"endpointId":"b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7","executed":1,"failed":0}}}}}`, nil
				}
			}
			got, err := getInsightProjectsActivity(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, tt.args.epClt)
			responseMap := map[string]interface{}{}
			json.Unmarshal(got, &responseMap)
			if tt.name == "Project activity success case - Non CJOC" || tt.name == "Project activity success case - CJOC" || tt.name == "Project activity success case - Idle filter" || tt.name == "Project activity success case - Fragile filter" {
				assert.Equal(t, len(responseMap["data"].([]interface{})), 0, "Validating data count")
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getInsightProjectsActivity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestConvertTime(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		timezone      string
		expectedError bool
	}{
		{
			name:     "Valid Input and Timezone",
			input:    "2024-07-09 12:00:00",
			timezone: "America/New_York",
		},
		{
			name:          "Invalid Input Format",
			input:         "2024-07-09T12:00:00Z",
			timezone:      "America/New_York",
			expectedError: true,
		},
		{
			name:          "Invalid Timezone",
			input:         "2024-07-09 12:00:00",
			timezone:      "Invalid/Timezone",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ConvertTime(tt.input, tt.timezone)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else {
					_, parseErr := time.Parse("2006-01-02 15:04:05", output)
					if parseErr != nil {
						t.Errorf("Output format is incorrect: %v", parseErr)
					}
				}
			}
		})
	}
}

func TestGetActivityOverviewReplacements(t *testing.T) {
	replacements := map[string]any{
		"key": "value",
	}

	startDate := time.Date(2024, time.July, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, time.July, 9, 0, 0, 0, 0, time.UTC)

	expectedStartDate := startDate.Add(-(endDate.Sub(startDate)))
	expectedEndDate := endDate.Add(-(endDate.Sub(startDate)))

	expectedReplacements := map[string]any{
		constants.START_DATE: expectedStartDate.Format(constants.DATE_FORMAT_WITH_HYPHEN),
		constants.END_DATE:   expectedEndDate.Format(constants.DATE_FORMAT_WITH_HYPHEN),
		"key":                "value",
	}

	result := getActivityOverviewReplacements(replacements, endDate, startDate)
	for key, expectedValue := range expectedReplacements {
		actualValue, ok := result[key]
		if !ok {
			t.Errorf("Key %s not found in result map", key)
			continue
		}
		if actualValue != expectedValue {
			t.Errorf("Key %s: Expected value %v, but got %v", key, expectedValue, actualValue)
		}
	}
	for key := range result {
		if _, ok := expectedReplacements[key]; !ok {
			t.Errorf("Unexpected key %s found in result map", key)
		}
	}
}

func Test_jobAndRunCount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		if IndexName == constants.CB_CI_JOB_INFO_INDEX {
			return `{"aggregations":{"jobs_per_endpoint":{"buckets":[{"key":"endpoint123","doc_count":10},{"key":"endpoint456","doc_count":5}]}}}`, nil
		} else if IndexName == constants.CB_CI_RUN_INFO_INDEX {
			return `{"aggregations":{"jobs_per_endpoint":{"buckets":[{"key":"endpoint123","doc_count":8},{"key":"endpoint456","doc_count":4}]}}}`, nil
		}
		return "", fmt.Errorf("unexpected index name: %s", IndexName)
	}

	replacements := map[string]interface{}{
		"key": "value",
	}

	ctx := context.Background()

	jobCounts, runCounts, err := JobAndRunCount(replacements, ctx)

	assert.NoError(t, err, "Expected no error from JobAndRunCount")

	expectedJobCounts := map[string]int{
		"endpoint123": 10,
		"endpoint456": 5,
	}

	expectedRunCounts := map[string]int{
		"endpoint123": 8,
		"endpoint456": 4,
	}

	assert.Equal(t, expectedJobCounts, jobCounts, "Job counts should match expected")
	assert.Equal(t, expectedRunCounts, runCounts, "Run counts should match expected")
}

func Test_getVersionAndPluginCount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	searchResponseCalled := false
	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		searchResponseCalled = true
		if IndexName == constants.CB_CI_TOOL_INSIGHT_INDEX {
			return `{"aggregations":{"plugin_count":{"value":{"plugin1":10,"plugin2":5}}}}`, nil
		}
		return "", fmt.Errorf("unexpected index name: %s", IndexName)
	}

	replacements := map[string]interface{}{
		"key": "value",
	}

	ctx := context.Background()

	result, err := GetVersionAndPluginCount(replacements, ctx)

	assert.True(t, searchResponseCalled, "searchResponse should be called")
	assert.NoError(t, err, "Expected no error from GetVersionAndPluginCount")

	if err != nil {
		t.Fatalf("Error calling GetVersionAndPluginCount: %v", err)
	}
	assert.Nil(t, result, "Result should match expected")

}

func Test_getInsightCjocControllers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
		epClt        endpoint.EndpointServiceClient
	}
	replacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "0edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType":       "CJOC",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "CJOC controllers success case - CJOC",
			args: args{
				widgetId:     "ci5",
				replacements: replacements,
			},
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		if IndexName == constants.CB_CI_JOB_INFO_INDEX {
			return `{"aggregations":{"completedRuns":{"buckets":[{"key":"job123","doc_count":10,"result_buckets":{"buckets":{"SUCCESS":{"doc_count":5,"last_active":{"value":1632211620000},"total_duration":{"value":1632211620000}},"FAILED":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}},"ABORTED":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}},"UNSTABLE":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}},"NOT_BUILT":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}}}},"endpoint_id":{"buckets":[{"key":"endpoint123"}]}}]},"endpointJobs":{"value": {"job123":["jobId1","jobId2"]}}}}`, nil
		} else if IndexName == constants.CB_CI_CJOC_CONTROLLER_INFO {
			return `{"aggregations":{"cjocControllerInfo":{"value":["https://test-controller/","https://test-controller-1/","https://test-controller-2/"]}}}`, nil
		}
		return `{"aggregations":{"completedRuns":{"buckets":[{"key":"job123","doc_count":10,"result_buckets":{"buckets":{"SUCCESS":{"doc_count":5,"last_active":{"value":1632211620000},"total_duration":{"value":1632211620000}},"FAILED":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}},"ABORTED":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}},"UNSTABLE":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}},"NOT_BUILT":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}}}},"endpoint_id":{"buckets":[{"key":"endpoint123"}]}}]},"endpointJobs":{"value": {"job123":["jobId1","jobId2"]}}}}`, nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "CJOC controllers success case - CJOC" {
				getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
					return &endpoint.EndpointsResponse{
						Endpoints: []*endpoint.Endpoint{
							{
								Id:   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
								Name: "CJOC Test",
							},
							{
								Id:             "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
								Name:           "CB1 CBCI Test",
								ContributionId: "cjoc_cbci-app-endpoint-type",
								Properties: []*api.Property{
									{
										Name: "status",
										Value: &api.Property_String_{
											String_: "INSTALLED",
										},
									},
									{
										Name: constants.CJOC_ENDPOINT_ID,
										Value: &api.Property_String_{
											String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
										},
									},
									{
										Name: constants.TOOL_URL,
										Value: &api.Property_String_{
											String_: "https://test-controller/",
										},
									},
								},
							},
							{
								Id:             "a7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
								Name:           "CBCI Test",
								ContributionId: "cjoc_cbci-app-endpoint-type",
								Properties: []*api.Property{
									{
										Name: "status",
										Value: &api.Property_String_{
											String_: "NOT_INSTALLED",
										}, Audit: &api.Audit{
											When: &timestamppb.Timestamp{},
										},
									},
									{
										Name: constants.CJOC_ENDPOINT_ID,
										Value: &api.Property_String_{
											String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
										},
									},
								},
							},
						},
					}, nil
				}
			}
			got, err := getInsightCjocControllers(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, tt.args.epClt)
			responseMap := map[string]interface{}{}
			json.Unmarshal(got, &responseMap)
			if tt.name == "CJOC controllers success case - CJOC" {
				assert.Equal(t, responseMap["connectedControllers"].(float64), float64(1), "Validating data count")
				assert.Equal(t, len(responseMap["data"].([]interface{})), 1, "Validating data count")
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getInsightCjocControllers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_getInsightRunsOverview(t *testing.T) {
	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
		epClt        endpoint.EndpointServiceClient
	}
	replacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "0edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType":       "CBCI",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
	}
	CJOCReplacements := map[string]any{
		"ciToolId":   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType": "CJOC",
		"endDate":    "2023-12-30 23:59:59",
		"subOrgId":   "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":  "2023-12-01 00:00:00",
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "CBCI controllers success case - CJOC",
			args: args{
				widgetId:     "ci4",
				replacements: replacements,
			},
		},
		{
			name: "CJOC controllers success case - CJOC",
			args: args{
				widgetId:     "ci4",
				replacements: CJOCReplacements,
			},
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		if IndexName == constants.CB_CI_ACTIVITY_OVERVIEW {
			return `{"hits":{"hits":[{"_index":"cb_ci_activity_overview","_id":"2","_score":0,"_source":{"activity_time":12,"current_time_to_idle":77000,"activity_date":"2024-01-02","endpoint_id":"3cfc093d-e586-4bb5-9b67-0b6765f7d031","active_runs":4,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","activity_day":"TUES","org_name":"cloudbees-staging","idle_executor":2,"runs_waiting_to_start":8,"avg_time_waiting_to_start":91000,"timestamp":"2024-01-03 10:58:03"}}]}}`, nil
		}
		return `{"took":652,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10000,"relation":"gte"},"max_score":null,"hits":[]},"aggregations":{"latest_per_org":{"doc_count_error_upper_bound":0,"sum_other_doc_count":7,"buckets":[{"key":"5bb87055-68cc-4fc5-acdb-14490014b9b6","doc_count":127,"latest_doc":{"hits":{"total":{"value":127,"relation":"eq"},"max_score":null,"hits":[{"_index":"cb_ci_runs_activity","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5bb87055-68cc-4fc5-acdb-14490014b9b6_2024-01-31 11:55:55","_score":null,"_source":{"current_time_to_idle":0,"created_time":11,"endpoint_id":"5bb87055-68cc-4fc5-acdb-14490014b9b6","created_at":"2024-01-31 11:55:55","active_runs":0,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","created_date":"2024-01-31","org_name":"cloudbees-staging","idle_executor":2,"runs_waiting_to_start":0,"avg_time_waiting_to_start":0,"timestamp":"2024-01-31 11:55:55"},"sort":[1706702155000]}]}},"latest_timestamp":{"value":1706702155000,"value_as_string":"2024-01-31 11:55:55"}}]}}}`, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
				return &endpoint.EndpointsResponse{
					Endpoints: []*endpoint.Endpoint{
						{
							Id:   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
							Name: "CJOC Test",
						},
						{
							Id:             "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
							Name:           "CB1 CBCI Test",
							ContributionId: "cjoc_cbci-app-endpoint-type",
							Properties: []*api.Property{
								{
									Name: "status",
									Value: &api.Property_String_{
										String_: "INSTALLED",
									},
								},
								{
									Name: constants.CJOC_ENDPOINT_ID,
									Value: &api.Property_String_{
										String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
									},
								},
							},
						},
						{
							Id:             "a7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
							Name:           "CBCI Test",
							ContributionId: "cjoc_cbci-app-endpoint-type",
							Properties: []*api.Property{
								{
									Name: "status",
									Value: &api.Property_String_{
										String_: "NOT_INSTALLED",
									},
								},
								{
									Name: constants.CJOC_ENDPOINT_ID,
									Value: &api.Property_String_{
										String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
									},
								},
							},
						},
					},
				}, nil

			}
			got, err := getInsightRunsOverview(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, tt.args.epClt)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInsightRunsOverview() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			responseMap := map[string]interface{}{}
			json.Unmarshal(got, &responseMap)
			assert.Equal(t, len(responseMap["data"].([]interface{})), 5, "Validating data count")
		})
	}
}

func Test_getInsightUsagePatterns(t *testing.T) {
	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
		epClt        endpoint.EndpointServiceClient
	}
	replacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "0edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType":       "CBCI",
		"viewOption":       "ActiveRuns",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
		"timeZone":         "Asia/Calcutta",
		"timeFormat":       "12h",
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "Usage pattern success case - CJOC",
			args: args{
				widgetId:     "ci6",
				replacements: replacements,
			},
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		return `{
			"took": 604,
			"timed_out": false,
			"_shards": {
			  "total": 2,
			  "successful": 2,
			  "skipped": 0,
			  "failed": 0
			},
			"hits": {
			  "total": {
				"value": 120,
				"relation": "eq"
			  },
			  "max_score": null,
			  "hits": []
			},
			"aggregations": {
			  "activities": {
				"buckets": [
				  {
					"key_as_string": "2024-05-02 20:00:00",
					"key": 1714660200000,
					"doc_count": 120,
					"endpoint_ids": {
					  "doc_count_error_upper_bound": 0,
					  "sum_other_doc_count": 0,
					  "buckets": [
						{
						  "key": "d8732390-8a50-4987-9d41-993de47bf90e",
						  "doc_count": 120,
						  "avg_idle_executor": {
							"value": 0
						  },
						  "avg_active_runs": {
							"value": 0
						  },
						  "avg_current_time_to_idle": {
							"value": 0
						  },
						  "avg_time_waiting_to_start": {
							"value": 0
						  },
						  "avg_runs_waiting_to_start": {
							"value": 0
						  }
						}
					  ]
					}
				  }
				]
			  }
			}
		  }`, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
				return &endpoint.EndpointsResponse{
					Endpoints: []*endpoint.Endpoint{
						{
							Id:   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
							Name: "CJOC Test",
						},
						{
							Id:             "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
							Name:           "CB1 CBCI Test",
							ContributionId: "cjoc-app-endpoint-type",
							Properties: []*api.Property{
								{
									Name: "status",
									Value: &api.Property_String_{
										String_: "INSTALLED",
									},
								},
								{
									Name: constants.CJOC_ENDPOINT_ID,
									Value: &api.Property_String_{
										String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
									},
								},
							},
						},
						{
							Id:             "a7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
							Name:           "CBCI Test",
							ContributionId: "cjoc-app-endpoint-type",
							Properties: []*api.Property{
								{
									Name: "status",
									Value: &api.Property_String_{
										String_: "NOT_INSTALLED",
									},
								},
								{
									Name: constants.CJOC_ENDPOINT_ID,
									Value: &api.Property_String_{
										String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
									},
								},
							},
						},
					},
				}, nil
			}
			got, err := getInsightUsagePatterns(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, tt.args.epClt)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInsightUsagePatterns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			responseMap := map[string]interface{}{}
			json.Unmarshal(got, &responseMap)
			assert.Equal(t, len(responseMap["data"].([]interface{})), 7, "Validating data count")
		})
	}
}

func Test_getInsightProjectTypes(t *testing.T) {
	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
		epClt        endpoint.EndpointServiceClient
	}
	replacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "0edebf0d-6797-4ec0-9a50-fac728645a0e",
		"ciToolType":       "CJOC",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "Usage pattern success case - CJOC",
			args: args{
				widgetId:     "ci1",
				replacements: replacements,
			},
		},
	}
	getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
		return &endpoint.EndpointsResponse{
			Endpoints: []*endpoint.Endpoint{
				{
					Id:   "0edebf0d-6797-4ec0-9a50-fac728645a0e",
					Name: "CJOC Test",
				},
				{
					Id:             "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
					Name:           "CB1 CBCI Test",
					ContributionId: "cjoc_cbci-app-endpoint-type",
					Properties: []*api.Property{
						{
							Name: "status",
							Value: &api.Property_String_{
								String_: "INSTALLED",
							},
						},
						{
							Name: constants.CJOC_ENDPOINT_ID,
							Value: &api.Property_String_{
								String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
							},
						},
					},
				},
				{
					Id:             "a7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
					Name:           "CBCI Test",
					ContributionId: "cjoc_cbci-app-endpoint-type",
					Properties: []*api.Property{
						{
							Name: "status",
							Value: &api.Property_String_{
								String_: "NOT_INSTALLED",
							},
						},
						{
							Name: constants.CJOC_ENDPOINT_ID,
							Value: &api.Property_String_{
								String_: "0edebf0d-6797-4ec0-9a50-fac728645a0e",
							},
						},
					},
				},
			},
		}, nil
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		return `{"aggregations":{"projectTypes":{"value":[{"name":"Freestyle","value":1},{"name":"Pipeline","value":2}]}}}`, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getInsightProjectTypes(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, tt.args.epClt)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInsightProjectTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			responseMap := map[string]interface{}{}
			json.Unmarshal(got, &responseMap)
			assert.Equal(t, len(responseMap["data"].([]interface{})), 2, "Validating data count")
			runInfoByte, err := json.Marshal(responseMap[constants.RUNS_INFO].(interface{}))
			runInfo := &RunsInfo{}
			json.Unmarshal(runInfoByte, runInfo)
			assert.Equal(t, runInfo.TotalNumberofProjects, int64(3), "Validating total number of Projects")

		})
	}

	testsEmptyProject := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "Usage pattern success case empty Project - CJOC",
			args: args{
				widgetId:     "ci1",
				replacements: replacements,
			},
		},
	}

	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		return `{"aggregations":{"projectTypes":{"value":[{"name":"Freestyle","value":0},{"name":"Pipeline","value":0}]}}}`, nil
	}

	for _, tt := range testsEmptyProject {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getInsightProjectTypes(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, tt.args.epClt)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInsightProjectTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			responseMap := map[string]interface{}{}
			json.Unmarshal(got, &responseMap)
			assert.Equal(t, len(responseMap["data"].([]interface{})), 2, "Validating data count")
			runInfoByte, err := json.Marshal(responseMap[constants.RUNS_INFO].(interface{}))
			runInfo := &RunsInfo{}
			json.Unmarshal(runInfoByte, runInfo)
			assert.Equal(t, runInfo.TotalNumberofProjects, int64(0), "Validating total number of Projects")

		})
	}
}

func TestGetReadableDuration(t *testing.T) {
	response := getReadableDuration(100000)
	assert.NotNil(t, response)
	response = getReadableDuration(10000000)
	assert.NotNil(t, response)
	response = getReadableDuration(0)
	assert.NotNil(t, response)
}

func TestGetProcessedJobType(t *testing.T) {
	response := GetProcessedJobType(constants.FREESTYLE_JOB)
	assert.NotNil(t, response)
	response = GetProcessedJobType(constants.PIPELINE_JOB_TEMPLATE)
	assert.NotNil(t, response)
	response = GetProcessedJobType(constants.WORKFLOW_JOB)
	assert.NotNil(t, response)
	response = GetProcessedJobType(constants.MATRIX_JOB)
	assert.NotNil(t, response)
	response = GetProcessedJobType(constants.JENKINS_MULTI_JOB)
	assert.NotNil(t, response)
	response = GetProcessedJobType(constants.JENKINS_FOLDER)
	assert.NotNil(t, response)
	response = GetProcessedJobType(constants.JENKINS_BRANCH)
	assert.NotNil(t, response)
}

func TestConvertMilliSecToTime(t *testing.T) {
	t.Run("for 100000000000 input", func(t *testing.T) {
		response := convertMilliSecToTime(100000000000)
		assert.NotNil(t, response)
	})

	t.Run("empty input ", func(t *testing.T) {
		response := convertMilliSecToTime(0)
		assert.Equal(t, "0s", response)
	})

	t.Run("negative value ", func(t *testing.T) {
		response := convertMilliSecToTime(-1)
		assert.Equal(t, "0s", response)
	})

}

func TestGetWeekDayByName(t *testing.T) {
	response := getWeekDayByName(constants.SATURDAY)
	assert.NotNil(t, response)
	response = getWeekDayByName(constants.MONDAY)
	assert.NotNil(t, response)
	response = getWeekDayByName(constants.TUESDAY)
	assert.NotNil(t, response)
	response = getWeekDayByName(constants.WEDNESDAY)
	assert.NotNil(t, response)
	response = getWeekDayByName(constants.THURSDAY)
	assert.NotNil(t, response)
	response = getWeekDayByName(constants.FRIDAY)
	assert.NotNil(t, response)
	response = getWeekDayByName(constants.SUNDAY)
	assert.NotNil(t, response)
	response = getWeekDayByName("")
	assert.NotNil(t, response)
}

func Test_GetControllerUrlMap(t *testing.T) {
	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
		epClt        endpoint.EndpointServiceClient
	}
	replacements := map[string]any{
		"aggrBy":           "week",
		"ciToolId":         "42e98b10-f58b-4b64-8645-d42c00f06e5f",
		"ciToolType":       "CJOC",
		"viewOption":       "ActiveRuns",
		"dateHistogramMax": "2023-12-30",
		"dateHistogramMin": "2023-12-01",
		"duration":         "month",
		"endDate":          "2023-12-30 23:59:59",
		"subOrgId":         "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		"startDate":        "2023-12-01 00:00:00",
		"timeFormat":       "12h",
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			name: "Usage pattern success case - CJOC",
			args: args{
				widgetId:     "ci6",
				replacements: replacements,
			},
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		return `{
			"aggregations": {
				"cjocControllerInfo": {
				  "value": {
					"42e98b10-f58b-4b64-8645-d42c00f06e5f": [
					  "https://cb2.rosaas.releaseiq.io/",
					  "https://controllertest01.rosaas.releaseiq.io/",
					  "https://cb1.rosaas.releaseiq.io/",
					  "https://test2.rosaas.releaseiq.io/",
					  "https://testcontroller.rosaas.releaseiq.io/",
					  "https://cb3.rosaas.releaseiq.io/"
					]
				  }
				}
			  }
		  }`, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetControllerUrlMap(tt.args.replacements, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInsightUsagePatterns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, len(got["42e98b10-f58b-4b64-8645-d42c00f06e5f"]), 6, "Validating data count")
		})
	}
}

func TestConvertToAMPM(t *testing.T) {
	//Case1 : 24h format 0
	response := convertToAMPM(0, "24h")
	assert.Equal(t, response, "00:00")

	//Case2 : 24h format 9
	response = convertToAMPM(9, "24h")
	assert.Equal(t, response, "09:00")

	//Case3 : 24h format 15
	response = convertToAMPM(15, "24h")
	assert.Equal(t, response, "15:00")

	//Case4 : 12h format 0
	response = convertToAMPM(0, "12h")
	assert.Equal(t, response, "12am")

	//Case5 : 12h format 9
	response = convertToAMPM(9, "12h")
	assert.Equal(t, response, "9am")

	//Case6 : 12h format 15
	response = convertToAMPM(15, "12h")
	assert.Equal(t, response, "3pm")

	//Case7 : 12h format 12
	response = convertToAMPM(12, "12h")
	assert.Equal(t, response, "12pm")

}

func Test_updateProjectActivityResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	mockOpenSearchClient, _ := openSearchClient()

	mockReplacements := make(map[string]interface{})
	mockOutputResponse := make(map[string]interface{})
	mockControllerMap := make(map[string]ControllerInfo)

	updateProjectActivity := models.UpdateProjectActivityResponse{
		Response:       `{"aggregations":{"completedRuns":{"buckets":[{"key":"job123","doc_count":10,"result_buckets":{"buckets":{"SUCCESS":{"doc_count":5,"last_active":{"value":1632211620000},"total_duration":{"value":1632211620000}},"FAILED":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}},"ABORTED":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}},"UNSTABLE":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}},"NOT_BUILT":{"doc_count":3,"last_active":{"value":1632211520000},"total_duration":{"value":1632211620000}}}},"endpoint_id":{"buckets":[{"key":"endpoint123"}]}}]}}}`,
		ReportId:       "report123",
		Replacements:   mockReplacements,
		Client:         mockOpenSearchClient,
		OutputResponse: mockOutputResponse,
		IsIdle:         false,
		JobIds:         []string{"job123"},
		IsFragile:      true,
	}

	updateProjectActivityResponse(updateProjectActivity, mockControllerMap)

}
