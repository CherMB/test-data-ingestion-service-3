package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	api "github.com/calculi-corp/api/go"
	coredataMock "github.com/calculi-corp/core-data-cache/mock"
	"github.com/calculi-corp/reports-service/db"
	"github.com/calculi-corp/reports-service/mocks"
	"github.com/calculi-corp/reports-service/models"
	"github.com/google/uuid"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	auth "github.com/calculi-corp/api/go/auth"
	"github.com/calculi-corp/api/go/endpoint"
	"github.com/calculi-corp/api/go/service"
	pb "github.com/calculi-corp/api/go/vsm/report"
	"github.com/calculi-corp/common/grpc"
	"github.com/calculi-corp/config"
	cmock "github.com/calculi-corp/core-data-cache/mock"
	client "github.com/calculi-corp/grpc-client"
	mock_grpc_client "github.com/calculi-corp/grpc-client/mock"
	handler "github.com/calculi-corp/grpc-handler"
	hostflags "github.com/calculi-corp/grpc-hostflags"
	testutil "github.com/calculi-corp/grpc-testutil"
	testutil_setup "github.com/calculi-corp/grpc-testutil/setup"
	"github.com/calculi-corp/reports-service/cache"
	cache2 "github.com/calculi-corp/reports-service/cache"
	"github.com/calculi-corp/reports-service/constants"
	rmock "github.com/calculi-corp/reports-service/mock"
	"github.com/calculi-corp/repository-service/mock"
	"github.com/opensearch-project/opensearch-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func init() {
	config.Config.Set("logging.level", "INFO")
	config.Config.DefineStringFlag("report.definition.filepath", "../resources/", "Path to json for widget configurations")
	testutil.SetUnitTestConfig()
}

type mockServer struct {
	sendFunc func(*pb.StreamCIInsightsCompletedRunsResponse) error
}

func (m *mockServer) Send(resp *pb.StreamCIInsightsCompletedRunsResponse) error {
	return m.sendFunc(resp)
}

func (m *mockServer) RecvMsg(any) error {
	return nil
}

func (m *mockServer) SendHeader(metadata.MD) error {
	return nil
}

func (m *mockServer) SendMsg(any) error {
	return nil
}

func (m *mockServer) SetHeader(metadata.MD) error {
	return nil
}

func (m *mockServer) SetTrailer(metadata.MD) {
}

func (m *mockServer) Context() context.Context {
	return context.Background()
}

type MockOpensearchConnection struct{}

func (m *MockOpensearchConnection) GetOpensearchConnection() (*opensearch.Client, error) {
	return &opensearch.Client{}, nil
}

func TestBuildDrilldownReport_1(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	epClient := mock.NewMockEndpointServiceClient(mockCtrl)
	mockCtrl1 := gomock.NewController(t)
	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)
	rah := &ReportsHandler{
		client:           mockGrpcClient,
		endpointClient:   epClient,
		orgServiceClient: orgClient,
	}

	mockGrpcClient.EXPECT().
		SendGrpcCtx(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		Do(func(ctx context.Context, host, endpoint, method string, request, response interface{}) {
		}).
		AnyTimes()

	serviceResponse := &service.ListServicesResponse{
		Service: []*service.Service{
			{
				Id:            "2cab10cc-cd9d-11ed-afa1-0242ac120002",
				Name:          "dsl-no-commit-1",
				RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
			},
			{
				Id:            "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "dsl-engine-cli",
				RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
			},
			{
				Id:            "138ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "test",
				RepositoryUrl: "https://github.com/calculi-corp/test.git",
			},
		},
	}
	orgClient.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "2cab10cc-cd9d-11ed-afa1-0242ac120002", DisplayName: "Cloudbees Stagging", ParentId: "0"}}, nil).AnyTimes()

	mockCache := cmock.NewMockResourceCacheI(mockCtrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(mockCtrl)
	cache2.SetMockCache(mockCache, mockEpCache, mockGrpcClient)
	mockCache.EXPECT().GetParentIDs(gomock.Any()).Return([]string{"org123", "parent_id_4"}).AnyTimes()

	getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
		return &endpoint.EndpointsResponse{
			Endpoints: []*endpoint.Endpoint{
				{
					Id:         "ciTool123",
					Name:       "CJOC Test",
					ResourceId: "2cab10cc-cd9d-11ed-afa1-0242ac120002",
				},
			},
		}, nil

	}

	getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
		return serviceResponse, nil
	}

	req := &pb.DrilldownRequest{
		ReportId:     "report123",
		OrgId:        "org123",
		SubOrgId:     "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		Component:    []string{},
		DurationType: pb.DurationType_CURRENT_MONTH,
		StartDate:    "2024-01-01",
		EndDate:      "2024-01-31",
		CiToolId:     "ciTool123",
		TimeZone:     "UTC",
		UserId:       "",
	}

	rah.BuildDrilldownReport(ctx, req)

}

func TestBuildDrilldownReport(t *testing.T) {
	ctx := context.Background()

	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	epClient := mock.NewMockEndpointServiceClient(mockCtrl)

	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)
	rah := &ReportsHandler{
		client:           mockGrpcClient,
		endpointClient:   epClient,
		orgServiceClient: orgClient,
	}

	mockGrpcClient.EXPECT().
		SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), "ListServices", gomock.Any(), gomock.Any()).
		Return(nil).
		Do(func(ctx context.Context, host, endpoint, method string, request, response interface{}) {
		}).
		AnyTimes()

	serviceResponse := &service.ListServicesResponse{
		Service: []*service.Service{
			{
				Id:            "2cab10cc-cd9d-11ed-afa1-0242ac120002",
				Name:          "dsl-no-commit-1",
				RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
			},
			{
				Id:            "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "dsl-engine-cli",
				RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
			},
			{
				Id:            "138ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "test",
				RepositoryUrl: "https://github.com/calculi-corp/test.git",
			},
		},
	}
	orgClient.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "2cab10cc-cd9d-11ed-afa1-0242ac120002", DisplayName: "Cloudbees Stagging", ParentId: "0"}}, nil).AnyTimes()

	mockCache := cmock.NewMockResourceCacheI(mockCtrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(mockCtrl)
	cache2.SetMockCache(mockCache, mockEpCache, mockGrpcClient)
	mockCache.EXPECT().GetParentIDs(gomock.Any()).Return([]string{"org123", "parent_id_4"}).AnyTimes()

	getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
		return serviceResponse, nil
	}

	req1 := &pb.DrilldownRequest{
		ReportId:     "report123",
		OrgId:        "org123",
		SubOrgId:     "org123",
		Component:    []string{},
		DurationType: pb.DurationType_CURRENT_MONTH,
		StartDate:    "2024-01-01",
		EndDate:      "2024-01-31",
		CiToolId:     "ciTool123",
		TimeZone:     "UTC",
		UserId:       uuid.NewString(),
	}

	req2 := &pb.DrilldownRequest{
		ReportId:     "report123",
		OrgId:        "org123",
		SubOrgId:     "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		Component:    []string{},
		DurationType: pb.DurationType_CURRENT_MONTH,
		StartDate:    "2024-01-01",
		EndDate:      "2024-01-31",
		CiToolId:     "ciTool123",
		TimeZone:     "UTC",
		UserId:       uuid.NewString(),
	}

	getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
		return &endpoint.EndpointsResponse{
			Endpoints: []*endpoint.Endpoint{
				{
					Id:         "ciTool123",
					Name:       "CJOC Test",
					ResourceId: "2cab10cc-cd9d-11ed-afa1-0242ac120002",
				},
			},
		}, nil

	}

	rah.BuildDrilldownReport(ctx, req1)

	rah.BuildDrilldownReport(ctx, req2)

}

func TestBuildComputedReport(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	epClient := mock.NewMockEndpointServiceClient(mockCtrl)
	rah := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: epClient,
	}

	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(),
		"ListServices", gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	serviceResponse := &service.ListServicesResponse{
		Service: []*service.Service{
			{
				Id:            "2cab10cc-cd9d-11ed-afa1-0242ac120002",
				Name:          "dsl-no-commit-1",
				RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
			},
			{
				Id:            "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "dsl-engine-cli",
				RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
			},
			{
				Id:            "138ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "test",
				RepositoryUrl: "https://github.com/calculi-corp/test.git",
			},
		},
	}

	mockCache := cmock.NewMockResourceCacheI(mockCtrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(mockCtrl)
	cache2.SetMockCache(mockCache, mockEpCache, mockGrpcClient)
	mockCache.EXPECT().GetParentIDs(gomock.Any()).Return([]string{"org123", "parent_id_4"}).AnyTimes()

	getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
		return serviceResponse, nil
	}

	req1 := &pb.ReportServiceRequest{
		WidgetId:     "cs1",
		OrgId:        "org123",
		SubOrgId:     "suborg123",
		Component:    []string{},
		DurationType: pb.DurationType_CURRENT_MONTH,
		StartDate:    "2024-01-01",
		EndDate:      "2024-01-31",
	}
	req2 := &pb.ReportServiceRequest{
		WidgetId:     "cs1",
		OrgId:        "org123",
		SubOrgId:     "suborg123",
		Component:    []string{"All", "2cab10cc-cd9d-11ed-afa1-0242ac120002"},
		DurationType: pb.DurationType_CURRENT_MONTH,
		StartDate:    "2024-01-01",
		EndDate:      "2024-01-31",
	}

	rah.BuildComputedReport(ctx, req1)

	rah.BuildComputedReport(ctx, req2)

}

func TestBuildComputedDrilldownReport_1(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	epClient := mock.NewMockEndpointServiceClient(mockCtrl)

	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)

	mockRbacClient := rmock.NewMockRBACServiceClient(mockCtrl)

	rah := &ReportsHandler{
		client:           mockGrpcClient,
		endpointClient:   epClient,
		orgServiceClient: orgClient,
		rbacClt:          mockRbacClient,
	}

	mockGrpcClient.EXPECT().
		SendGrpcCtx(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		Do(func(ctx context.Context, host, endpoint, method string, request, response interface{}) {
		}).
		AnyTimes()

	serviceResponse := &service.ListServicesResponse{
		Service: []*service.Service{
			{
				Id:            "2cab10cc-cd9d-11ed-afa1-0242ac120002",
				Name:          "dsl-no-commit-1",
				RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
			},
			{
				Id:            "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "dsl-engine-cli",
				RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
			},
			{
				Id:            "138ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "test",
				RepositoryUrl: "https://github.com/calculi-corp/test.git",
			},
		},
	}
	orgClient.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "2cab10cc-cd9d-11ed-afa1-0242ac120002", DisplayName: "Cloudbees Stagging", ParentId: "0"}}, nil).AnyTimes()

	mockCache := cmock.NewMockResourceCacheI(mockCtrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(mockCtrl)
	cache2.SetMockCache(mockCache, mockEpCache, mockGrpcClient)
	mockCache.EXPECT().GetParentIDs(gomock.Any()).Return([]string{"org123", "parent_id_4"}).AnyTimes()

	getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
		return serviceResponse, nil
	}

	mockRbacClient.EXPECT().IsAuthorized(gomock.Any(), gomock.Any()).Return(&auth.IsAuthorizedResponse{Authorized: true}, nil).AnyTimes()

	req1 := &pb.DrilldownRequest{
		ReportId:     "report123",
		OrgId:        "org123",
		SubOrgId:     "org123",
		Component:    []string{},
		DurationType: pb.DurationType_CURRENT_MONTH,
		StartDate:    "2024-01-01",
		EndDate:      "2024-01-31",
		CiToolId:     "ciTool123",
		TimeZone:     "UTC",
		UserId:       uuid.NewString(),
	}

	req2 := &pb.DrilldownRequest{
		ReportId:     "report123",
		OrgId:        "org123",
		SubOrgId:     "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		Component:    []string{},
		DurationType: pb.DurationType_CURRENT_MONTH,
		StartDate:    "2024-01-01",
		EndDate:      "2024-01-31",
		CiToolId:     "ciTool123",
		TimeZone:     "UTC",
		UserId:       uuid.NewString(),
	}

	req3 := &pb.DrilldownRequest{
		ReportId:     "report123",
		OrgId:        "org123",
		SubOrgId:     "suborg123",
		Component:    []string{"comp1", "All", "2cab10cc-cd9d-11ed-afa1-0242ac120002"},
		DurationType: pb.DurationType_CURRENT_MONTH,
		StartDate:    "2024-01-01",
		EndDate:      "2024-01-31",
		CiToolId:     "ciTool123",
		TimeZone:     "UTC",
		UserId:       uuid.NewString(),
	}

	rah.BuildComputedDrilldownReport(ctx, req1)

	rah.BuildComputedDrilldownReport(ctx, req2)

	rah.BuildComputedDrilldownReport(ctx, req3)
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

func TestBuildComputedDrilldownReport(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockRbacClient := rmock.NewMockRBACServiceClient(mockCtrl)

	rah := &ReportsHandler{rbacClt: mockRbacClient}
	req := &pb.DrilldownRequest{
		ReportId:     "report123",
		OrgId:        "org123",
		SubOrgId:     "",
		Component:    []string{},
		DurationType: pb.DurationType_CURRENT_MONTH,
		StartDate:    "2024-01-01",
		EndDate:      "2024-01-31",
		CiToolId:     "ciTool123",
		TimeZone:     "UTC",
		UserId:       uuid.NewString(),
	}

	mockRbacClient.EXPECT().IsAuthorized(gomock.Any(), gomock.Any()).Return(&auth.IsAuthorizedResponse{Authorized: true}, nil).AnyTimes()
	ctx := context.Background()
	rah.BuildComputedDrilldownReport(ctx, req)

}

func TestUpdateRawData(t *testing.T) {
	mockConfig := &MockOpensearchConfig{}
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	mockClient := &opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}

	responseBody := `{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()

	mockConfig.GetOpensearchConnectionFunc = func() (*opensearch.Client, error) {
		return mockClient, nil
	}

	req1 := &pb.DataRequest{
		Query:            "test query",
		IndexName:        "test_index",
		IsMappingRequest: true,
		JobName:          "jobName",
	}

	req2 := &pb.DataRequest{
		Query:            "test query",
		IndexName:        "test_index",
		IsMappingRequest: false,
		JobName:          "jobName",
	}

	rah := &ReportsHandler{}
	rah.UpdateRawData(context.TODO(), req1)
	rah.UpdateRawData(context.TODO(), req2)

}

func TestUpdateRawData_MappingRequest(t *testing.T) {

	mockConfig := &MockOpensearchConfig{}
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	mockClient := &opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}

	responseBody := `{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()

	mockConfig.GetOpensearchConnectionFunc = func() (*opensearch.Client, error) {
		return mockClient, nil
	}
	req := &pb.DataRequest{
		IndexName:        "test_index",
		IsMappingRequest: true,
	}

	rah := &ReportsHandler{}
	rah.UpdateRawData(context.TODO(), req)

}

func TestGetRawData(t *testing.T) {
	mockConfig := &MockOpensearchConfig{}
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	mockClient := &opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}

	responseBody := `{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()

	mockConfig.GetOpensearchConnectionFunc = func() (*opensearch.Client, error) {
		return mockClient, nil
	}

	req1 := &pb.DataRequest{
		Query:            "test query",
		IndexName:        "test_index",
		IsMappingRequest: true,
	}

	req2 := &pb.DataRequest{
		Query:            "test query",
		IndexName:        "test_index",
		IsMappingRequest: false,
	}

	rah := &ReportsHandler{}
	rah.GetRawData(context.TODO(), req1)
	rah.GetRawData(context.TODO(), req2)

}

func TestGetRawData_MappingRequest(t *testing.T) {

	mockConfig := &MockOpensearchConfig{}
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	mockClient := &opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}

	responseBody := `{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()

	mockConfig.GetOpensearchConnectionFunc = func() (*opensearch.Client, error) {
		return mockClient, nil
	}
	req := &pb.DataRequest{
		IndexName:        "test_index",
		IsMappingRequest: true,
	}

	rah := &ReportsHandler{}
	rah.GetRawData(context.TODO(), req)

}

func TestUpdateComputeData(t *testing.T) {
	mockConfig := &MockOpensearchConfig{}
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	mockClient := &opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}

	responseBody := `{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()

	mockConfig.GetOpensearchConnectionFunc = func() (*opensearch.Client, error) {
		return mockClient, nil
	}

	req := &pb.ComputeUpdateRequest{
		OrgId:       "org123",
		ComponentId: "All",
		MetricKey:   "metric123",
		Environment: "test",
	}

	mockCtx := context.TODO()

	rah := &ReportsHandler{}
	rah.UpdateComputeData(mockCtx, req)
}

func TestUpdateComputeData_nil(t *testing.T) {
	mockConfig := &MockOpensearchConfig{}
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	mockClient := &opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}

	responseBody := `{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()

	mockConfig.GetOpensearchConnectionFunc = func() (*opensearch.Client, error) {
		return mockClient, nil
	}

	req := &pb.ComputeUpdateRequest{
		OrgId:       "org123",
		ComponentId: "comp456",
		MetricKey:   "metric123",
		Environment: "test",
	}

	mockCtx := context.TODO()

	rah := &ReportsHandler{}
	rah.UpdateComputeData(mockCtx, req)

}

func TestPerformComputeInfo(t *testing.T) {

	mockConfig := &MockOpensearchConfig{}
	mockCtrl := gomock.NewController(t)
	mockTransport := mocks.NewMockTransport(mockCtrl)
	mockClient := &opensearch.Client{Transport: mockTransport, API: opensearchapi.New(mockTransport)}

	responseBody := `{"hits":{"hits":[{"_source":{"field1":"value1","field2":"value2"}}]}}`
	mockTransport.EXPECT().Perform(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(responseBody))}, nil).AnyTimes()

	mockConfig.GetOpensearchConnectionFunc = func() (*opensearch.Client, error) {
		return mockClient, nil
	}

	request := &pb.ReportServiceRequest{
		WidgetId:  "widget123",
		OrgId:     "org123",
		SubOrgId:  "suborg456",
		Component: []string{"comp1", "comp2"},
		StartDate: "2024-01-01",
		EndDate:   "2024-01-31",
	}

	// section
	response1 := &pb.ReportServiceResponse{
		Widget: &pb.Widget{
			Id: "widget123",
			Content: []*pb.Content{

				{
					Section: []*pb.ChartInfo{
						{
							Title:        "Chart Title",
							FunctionName: "Chart Function",
						},
					},
				},
			},
		},
	}

	// header
	response2 := &pb.ReportServiceResponse{
		Widget: &pb.Widget{
			Id: "widget123",
			Content: []*pb.Content{

				{
					Header: []*pb.MetricInfo{
						{
							Title:       "Chart Title",
							Description: "Chart Function",
						},
					},
				},
			},
		},
	}

	// footer
	response3 := &pb.ReportServiceResponse{
		Widget: &pb.Widget{
			Id: "widget123",
			Content: []*pb.Content{

				{

					Footer: []*pb.MetricInfo{
						{
							Title:       "Chart Title",
							Description: "Chart Function",
						},
					},
				},
			},
		},
	}

	// nil content
	response4 := &pb.ReportServiceResponse{
		Widget: &pb.Widget{
			Id: "widget123",
			Content: []*pb.Content{
				{},
			},
		},
	}

	computeResponse1 := ComputeInfo{
		OrgId:       request.OrgId,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		MetricKey:   "some_metric_key",
		MetricValue: response1.Widget,
	}

	computeResponse2 := ComputeInfo{
		OrgId:       request.OrgId,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		MetricKey:   "some_metric_key",
		MetricValue: response2.Widget,
	}

	computeResponse3 := ComputeInfo{
		OrgId:       request.OrgId,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		MetricKey:   "some_metric_key",
		MetricValue: response3.Widget,
	}

	computeResponse4 := ComputeInfo{
		OrgId:       request.OrgId,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		MetricKey:   "some_metric_key",
		MetricValue: response4.Widget,
	}

	performComputeInfo(mockClient, response1, request, computeResponse1)
	performComputeInfo(mockClient, response2, request, computeResponse2)
	performComputeInfo(mockClient, response3, request, computeResponse3)
	performComputeInfo(mockClient, response4, request, computeResponse4)

}

func TestStreamCIInsightsCompletedRun_2(t *testing.T) {

	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	epClient := mock.NewMockEndpointServiceClient(mockCtrl)
	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)

	rah := &ReportsHandler{
		client:           mockGrpcClient,
		endpointClient:   epClient,
		orgServiceClient: orgClient,
		metrics:          handler.NewMap("your_service_name"),
	}

	mockGrpcClient.EXPECT().SendGrpcCtx(ctx, gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mock := &mockServer{
		sendFunc: func(resp *pb.StreamCIInsightsCompletedRunsResponse) error {
			return nil
		},
	}

	endpointt := &endpoint.Endpoint{
		ResourceId:       "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		ContributionType: "cb.platform.environment",
	}
	endpoints := []*endpoint.Endpoint{endpointt}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpoints,
	}
	fakeResponse := &auth.GetOrganizationByIdResponse{}
	orgClient.EXPECT().
		GetOrganizationById(ctx, gomock.Any()).
		Return(fakeResponse, nil).
		AnyTimes()
	epClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil).AnyTimes()

	var srv pb.ReportServiceHandler_StreamCIInsightsCompletedRunServer = mock

	req := &pb.ReportServiceRequest{
		StartDate:  "2024-01-01",
		EndDate:    "2024-01-31",
		OrgId:      "org123",
		SubOrgId:   "suborg456",
		Component:  []string{"comp1", "comp2"},
		CiToolId:   "ci123",
		CiToolType: "jenkins",
		SortBy:     "date",
		FilterType: "type1",
		ViewOption: "option1",
		TimeZone:   "America/New_York",
	}

	mock.sendFunc = func(resp *pb.StreamCIInsightsCompletedRunsResponse) error {
		return errors.New("mock error")
	}

	rah.StreamCIInsightsCompletedRun(req, srv)

}

func TestReportRecordTiming(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	mockRbacCtl := rmock.NewMockRBACServiceClient(mockCtrl1)
	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)

	handler := &ReportsHandler{
		client:           mockGrpcClient,
		endpointClient:   mockEndpointServiceClient,
		rbacClt:          mockRbacCtl,
		orgServiceClient: orgClient,
		metrics:          handler.NewMap("your_service_name"),
	}
	require.NotNil(t, handler)

	var started time.Time
	started = time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)
	handler.RecordTiming("metric", started)
}

func TestCurrentMonthStartAndEndDate(t *testing.T) {
	tests := []struct {
		name         string
		monthStart   string
		monthEndFrom string
	}{
		{
			name:         "Test Current Month Start And End Date",
			monthStart:   "2024-07-01 00:00:00",
			monthEndFrom: "2024-07-31 23:59:59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CurrentMonthStartAndEndDate()
		})
	}
}

func TestCurrentWeekStartAndEndDate(t *testing.T) {
	tests := []struct {
		name          string
		weekStartFrom string
	}{
		{
			name:          "Test Current Week Start And End Date",
			weekStartFrom: "2024-07-14 00:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CurrentWeekStartAndEndDate()
		})
	}
}

func TestWeekStartDate(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Test Monday",
			input:    time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := weekStartDate(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestBuildReportLayout(t *testing.T) {
	ctx := context.Background()

	rah := &ReportsHandler{}

	req := &pb.ReportLayoutRequest{
		OrgId:         "org123",
		Component:     "component456",
		Branch:        "main",
		DashboardName: "Dashboard 1",
	}

	expectedWidgets := []*pb.ReportLayout{
		{
			WidgetId:     "cs1",
			WidgetName:   "components-activity",
			WidgetWidth:  12,
			WidgetHeight: 2,
			MockData:     false,
		},
		{
			WidgetId:     "cs2",
			WidgetName:   "components-builds",
			WidgetWidth:  6,
			WidgetHeight: 3,
			MockData:     false,
		},
		{
			WidgetId:     "cs3",
			WidgetName:   "components-deployments",
			WidgetWidth:  6,
			WidgetHeight: 3,
			MockData:     false,
		},
		{
			WidgetId:     "cs9",
			WidgetName:   "components-no-scanners-configured",
			WidgetWidth:  12,
			WidgetHeight: 2,
			MockData:     false,
		},
	}

	response, err := rah.BuildReportLayout(ctx, req)
	if err != nil {
		t.Fatalf("BuildReportLayout returned error: %v", err)
	}

	if response.Status != pb.Status_success {
		t.Errorf("Expected status to be success, got %v", response.Status)
	}

	if response.ComponentId != req.Component {
		t.Errorf("Expected component ID %s, got %s", req.Component, response.ComponentId)
	}

	if len(response.Widgets) != len(expectedWidgets) {
		t.Errorf("Expected %d widgets, got %d", len(expectedWidgets), len(response.Widgets))
	}

	for i, expected := range expectedWidgets {
		actual := response.Widgets[i]
		if actual.WidgetId != expected.WidgetId {
			t.Errorf("Expected WidgetId: %s, but got: %s", expected.WidgetId, actual.WidgetId)
		}
		if actual.WidgetName != expected.WidgetName {
			t.Errorf("Expected WidgetName: %s, but got: %s", expected.WidgetName, actual.WidgetName)
		}
		if actual.WidgetWidth != expected.WidgetWidth {
			t.Errorf("Expected WidgetWidth: %d, but got: %d", expected.WidgetWidth, actual.WidgetWidth)
		}
		if actual.WidgetHeight != expected.WidgetHeight {
			t.Errorf("Expected WidgetHeight: %d, but got: %d", expected.WidgetHeight, actual.WidgetHeight)
		}
		if actual.MockData != expected.MockData {
			t.Errorf("Expected MockData: %v, but got: %v", expected.MockData, actual.MockData)
		}
	}
}

func TestGetReportLayouts(t *testing.T) {
	widgetLayout := map[string]string{
		"widget1": `{"widgetId": "widget1", "widgetName": "Widget 1", "widgetWidth": 100, "widgetHeight": 50, "mockData": false}`,
		"widget2": `{"widgetId": "widget2", "widgetName": "Widget 2", "widgetWidth": 150, "widgetHeight": 75, "mockData": true}`,
	}

	widgets := []string{"widget1", "widget2"}

	expectedReportLayouts := []*pb.ReportLayout{
		{
			WidgetId:     "widget1",
			WidgetName:   "Widget 1",
			WidgetWidth:  100,
			WidgetHeight: 50,
			MockData:     false,
		},
		{
			WidgetId:     "widget2",
			WidgetName:   "Widget 2",
			WidgetWidth:  150,
			WidgetHeight: 75,
			MockData:     true,
		},
	}

	reportLayouts := getReportLayouts(widgets, widgetLayout)

	if len(reportLayouts) != len(expectedReportLayouts) {
		t.Errorf("Expected %d report layouts, but got %d", len(expectedReportLayouts), len(reportLayouts))
		return
	}

	for i, expected := range expectedReportLayouts {
		actual := reportLayouts[i]
		if actual.WidgetId != expected.WidgetId {
			t.Errorf("Expected WidgetId: %s, but got: %s", expected.WidgetId, actual.WidgetId)
		}
		if actual.WidgetName != expected.WidgetName {
			t.Errorf("Expected WidgetName: %s, but got: %s", expected.WidgetName, actual.WidgetName)
		}
		if actual.WidgetWidth != expected.WidgetWidth {
			t.Errorf("Expected WidgetWidth: %d, but got: %d", expected.WidgetWidth, actual.WidgetWidth)
		}
		if actual.WidgetHeight != expected.WidgetHeight {
			t.Errorf("Expected WidgetHeight: %d, but got: %d", expected.WidgetHeight, actual.WidgetHeight)
		}
		if actual.MockData != expected.MockData {
			t.Errorf("Expected MockData: %v, but got: %v", expected.MockData, actual.MockData)
		}
	}
}

func TestGetEnvironments(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	epClient := mock.NewMockEndpointServiceClient(mockCtrl)
	rah := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: epClient,
	}
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(),
		"ListServices", gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	serviceResponse := &service.ListServicesResponse{
		Service: []*service.Service{
			{
				Id:            "2cab10cc-cd9d-11ed-afa1-0242ac120002",
				Name:          "dsl-no-commit-1",
				RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
			},
			{
				Id:            "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "dsl-engine-cli",
				RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
			},
			{
				Id:            "138ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "test",
				RepositoryUrl: "https://github.com/calculi-corp/test.git",
			},
		},
	}
	getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
		return serviceResponse, nil
	}

	endpointt := &endpoint.Endpoint{
		ResourceId:       "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		ContributionType: "cb.platform.environment",
	}
	endpoints := []*endpoint.Endpoint{endpointt}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpoints,
	}

	epClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil).AnyTimes()

	tests := []struct {
		name           string
		req            *pb.EnvironmentRequest
		expectedError  bool
		expectedLength int
	}{
		{
			name: "Valid request with components",
			req: &pb.EnvironmentRequest{
				OrgId:    "org123",
				SubOrgId: "suborg456",
			},
			expectedError:  false,
			expectedLength: 0,
		},
		{
			name: "Valid request without components",
			req: &pb.EnvironmentRequest{
				OrgId:    "org789",
				SubOrgId: "suborg789",
			},
			expectedError:  false,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			response, err := rah.GetEnvironments(ctx, tt.req)

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedError, err)
				return
			}

			if err == nil {
				if len(response.Environments) != tt.expectedLength {
					t.Errorf("Expected %d environments, but got %d", tt.expectedLength, len(response.Environments))
				}
			}
		})
	}
}

func TestGetEnvironmentsv2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	epClient := mock.NewMockEndpointServiceClient(mockCtrl)
	rah := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: epClient,
	}
	endpointt := &endpoint.Endpoint{
		ResourceId:       "2cab10cc-cd9d-11ed-afa1-0242ac120002",
		ContributionType: "cb.platform.environment",
		Name:             "test_env",
	}
	endpoints := []*endpoint.Endpoint{endpointt}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpoints,
	}

	epClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil).AnyTimes()

	tests := []struct {
		name           string
		req            *pb.EnvironmentRequest
		expectedError  bool
		expectedLength int
	}{
		{
			name: "Environment does not exist",
			req: &pb.EnvironmentRequest{
				OrgId: "org123",
				Name:  "dev_env",
			},
			expectedError:  false,
			expectedLength: 0,
		},
		{
			name: "Environment exists",
			req: &pb.EnvironmentRequest{
				OrgId: "org789",
				Name:  "test_env",
			},
			expectedError:  false,
			expectedLength: 1,
		},
		{
			name: "List all environments",
			req: &pb.EnvironmentRequest{
				OrgId: "org789",
			},
			expectedError:  false,
			expectedLength: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			response, err := rah.GetEnvironmentsv2(ctx, tt.req)

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedError, err)
				return
			}

			if err == nil {
				if len(response.Environments) != tt.expectedLength {
					t.Errorf("Expected %d environments, but got %d", tt.expectedLength, len(response.Environments))
				}
			}
		})
	}
}

func TestGetReportRequest(t *testing.T) {
	tests := []struct {
		name         string
		req          *pb.ComputeUpdateRequest
		duration     string
		expected     *pb.ReportServiceRequest
		expectWidget string
	}{
		{
			name: "Current month duration",
			req: &pb.ComputeUpdateRequest{
				OrgId:       "org123",
				ComponentId: "comp456",
				MetricKey:   "metric123",
			},
			duration: constants.CURRENT_MONTH,
			expected: &pb.ReportServiceRequest{
				OrgId:        "org123",
				Component:    []string{"comp456"},
				DurationType: pb.DurationType_CURRENT_MONTH,
				WidgetId:     "metric123",
			},
			expectWidget: "metric123",
		},
		{
			name: "Current week duration",
			req: &pb.ComputeUpdateRequest{
				OrgId:       "org456",
				ComponentId: "comp789",
				MetricKey:   "metric456",
			},
			duration: constants.CURRENT_WEEK,
			expected: &pb.ReportServiceRequest{
				OrgId:        "org456",
				Component:    []string{"comp789"},
				DurationType: pb.DurationType_CURRENT_WEEK,
				WidgetId:     "metric456",
			},
			expectWidget: "metric456",
		},
		{
			name: "No duration specified",
			req: &pb.ComputeUpdateRequest{
				OrgId:       "org789",
				ComponentId: "comp123",
				MetricKey:   "metric789",
			},
			duration: "",
			expected: &pb.ReportServiceRequest{
				OrgId:     "org789",
				Component: []string{"comp123"},
				WidgetId:  "metric789",
			},
			expectWidget: "metric789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getReportRequest(tt.req, tt.duration)

			if result.OrgId != tt.expected.OrgId {
				t.Errorf("Expected OrgId: %s, but got: %s", tt.expected.OrgId, result.OrgId)
			}

			if !equalStringArrays(result.Component, tt.expected.Component) {
				t.Errorf("Expected Component: %v, but got: %v", tt.expected.Component, result.Component)
			}

			if result.WidgetId != tt.expectWidget {
				t.Errorf("Expected WidgetId: %s, but got: %s", tt.expectWidget, result.WidgetId)
			}
		})
	}
}

func equalStringArrays(arr1, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i, v := range arr1 {
		if v != arr2[i] {
			return false
		}
	}
	return true
}

func TestConstructComputeResponseQuery(t *testing.T) {
	tests := []struct {
		name        string
		dev         ComputeInfo
		expected    string
		expectError bool
	}{
		{
			name: "ComputeInfo with MetricValue",
			dev: ComputeInfo{
				OrgId:       "org123",
				ComponentId: "comp456",
				StartDate:   "2024-01-01",
				EndDate:     "2024-01-31",
				MetricKey:   "metric123",
				MetricValue: &pb.Widget{
					Title:       "Widget Title",
					Description: "Widget Description",
				},
			},
			expected: fmt.Sprintf("{\"index\":{\"_index\":\"%s\",\"_id\": \"%s_%s_%s_%s_%s\"}}\n%s\n",
				constants.COMPUTE_INDEX,
				"org123",
				"comp456",
				"2024-01-01",
				"2024-01-31",
				"metric123",
				func() string {
					data, err := json.Marshal(ComputeInfo{
						OrgId:       "org123",
						ComponentId: "comp456",
						StartDate:   "2024-01-01",
						EndDate:     "2024-01-31",
						MetricKey:   "metric123",
						MetricValue: &pb.Widget{
							Title:       "Widget Title",
							Description: "Widget Description",
						},
					})
					if err != nil {
						t.Fatalf("error marshalling data: %v", err)
					}
					return string(data)
				}()),
			expectError: false,
		},
		{
			name: "ComputeInfo with nil MetricValue",
			dev: ComputeInfo{
				OrgId:       "org123",
				ComponentId: "comp456",
				StartDate:   "2024-01-01",
				EndDate:     "2024-01-31",
				MetricKey:   "metric123",
				MetricValue: nil,
			},
			expected:    "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := constructComputeResponseQuery(tt.dev)
			if result != tt.expected {
				t.Errorf("Expected query: %s, but got: %s", tt.expected, result)
			}
		})
	}
}

func TestCheckChartNullResponse(t *testing.T) {
	tests := []struct {
		name     string
		response []*pb.ChartInfo
		expected error
	}{
		{
			name:     "Empty response",
			response: []*pb.ChartInfo{},
			expected: nil,
		},
		{
			name: "No nil data elements",
			response: []*pb.ChartInfo{
				{
					Data: &structpb.ListValue{},
				},
				{
					Data: &structpb.ListValue{},
				},
			},
			expected: nil,
		},
		{
			name: "One element with nil data",
			response: []*pb.ChartInfo{
				{
					Data: &structpb.ListValue{},
				},
				{
					Data: nil,
				},
			},
			expected: db.ErrInternalServer,
		},
		{
			name: "Multiple elements, one with nil data",
			response: []*pb.ChartInfo{
				{
					Data: &structpb.ListValue{},
				},
				{
					Data: nil,
				},
				{
					Data: &structpb.ListValue{},
				},
			},
			expected: db.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkChartNullResponse(tt.response)
			if err != tt.expected {
				t.Errorf("Expected error: %v, but got: %v", tt.expected, err)
			}
		})
	}
}

func TestCheckNullResponse(t *testing.T) {
	tests := []struct {
		name     string
		response []*pb.MetricInfo
		expected error
	}{
		{
			name:     "Empty response",
			response: []*pb.MetricInfo{},
			expected: nil,
		},
		{
			name: "No nil data elements",
			response: []*pb.MetricInfo{
				{
					Data: &structpb.Struct{},
				},
				{
					Data: &structpb.Struct{},
				},
			},
			expected: nil,
		},
		{
			name: "One element with nil data",
			response: []*pb.MetricInfo{
				{
					Data: &structpb.Struct{},
				},
				{
					Data: nil,
				},
			},
			expected: db.ErrInternalServer,
		},
		{
			name: "Multiple elements, one with nil data",
			response: []*pb.MetricInfo{
				{
					Data: &structpb.Struct{},
				},
				{
					Data: nil,
				},
				{
					Data: &structpb.Struct{},
				},
			},
			expected: db.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkNullResponse(tt.response)
			if err != tt.expected {
				t.Errorf("Expected error: %v, but got: %v", tt.expected, err)
			}
		})
	}
}

func TestConstructComputeScheduleQuery(t *testing.T) {
	lastRunTime := "2024-07-11T10:00:00Z"
	timestamp := "2024-07-11T11:00:00Z"
	dev := ComputeSchedule{
		JobName:     "job123",
		LastRunTime: &lastRunTime,
		Timestamp:   &timestamp,
	}

	expectedQuery := "{\"index\":{\"_index\":\"" + constants.COMPUTE_INDEX + "\",\"_id\": \"job123\"}}\n" +
		"{\"job_name\":\"job123\",\"last_run_start_time\":\"2024-07-11T10:00:00Z\",\"timestamp\":\"2024-07-11T11:00:00Z\"}\n"

	actualQuery := constructComputeScheduleQuery(dev)

	assert.Equal(t, expectedQuery, actualQuery)

	devWithoutOptional := ComputeSchedule{
		JobName: "job456",
	}

	expectedQueryWithoutOptional := "{\"index\":{\"_index\":\"" + constants.COMPUTE_INDEX + "\",\"_id\": \"job456\"}}\n" +
		"{\"job_name\":\"job456\"}\n"

	actualQueryWithoutOptional := constructComputeScheduleQuery(devWithoutOptional)

	assert.NotEqual(t, expectedQueryWithoutOptional, actualQueryWithoutOptional)

	devWithoutName := ComputeSchedule{
		LastRunTime: &lastRunTime,
	}

	constructComputeScheduleQuery(devWithoutName)

}

func TestGetCiTransitionConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEndpointClient := mock.NewMockEndpointServiceClient(ctrl)

	ctx := context.Background()
	req := &pb.DashboardLayoutRequest{
		OrgId: "org123",
	}

	fakeResponse := &endpoint.EndpointsResponse{
		Endpoints: []*endpoint.Endpoint{
			{Id: "endpoint1", Name: "Endpoint 1", ResourceId: "org123"},
			{Id: "endpoint2", Name: "Endpoint 2", ResourceId: "org123"},
		},
		Pagination: &api.Pagination{
			Page:       1,
			PageLength: 100,
			LastPage:   true,
		},
	}

	mockEndpointClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(fakeResponse, nil).Times(1)

	handler := &ReportsHandler{
		endpointClient: mockEndpointClient,
	}
	displayTransition, ci := handler.GetCiTransitionConfig(ctx, req)

	assert.False(t, displayTransition, "Expected displayTransition to be false")
	assert.NotNil(t, ci, "Expected ci to be not nil")
	assert.NotNil(t, ci.CiInsights, "Expected ci.CiInsights to be not nil")
	assert.True(t, ci.CiInsights.IsCiToolsFound, "Expected IsCiToolsFound to be true")

}

func TestReportServiceHandler_GetReportData_Precomputation01(t *testing.T) {
	// Create a mock context and request
	ctx := context.Background()
	errorReq := &pb.ReportServiceRequest{
		StartDate: "2023-06-01 10:00:00",
		EndDate:   "2023-06-30 10:00:00",
		Component: []string{"All"},
	}
	// Mock the required dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	newhandler := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, newhandler)

	tb, err := testutil_setup.NewTestBed(newhandler)
	defer cleanup(tb)
	if err != nil {
		t.Fatal("Error on creating Testbed", err)
	}
	if tb.Server == nil {
		t.Error("Should have been able to create the server")
	}

	rsh := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: mockEndpointServiceClient,
	}

	// mockGrpcClient.EXPECT().Connect(gomock.Any()).Return(nil)
	// Exception case for Mandatory field validation
	// Check for WidgetId
	validationErr := fmt.Errorf(errMissingRequiredField, widgetIdField)
	expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())
	_, err = rsh.GetReportData(ctx, errorReq)
	assert.Equal(t, expectedValidationStatus, err, "Expected validation error status")

	// Check for TenantId
	errorReq.WidgetId = "1"
	validationErr = fmt.Errorf(errMissingRequiredField, tenantIdField)
	expectedValidationStatus = status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())
	_, err = rsh.GetReportData(ctx, errorReq)
	assert.Equal(t, expectedValidationStatus, err, "Expected validation error status")

	mockGrpcClient.EXPECT().Close()
}

func TestReportServiceHandler_GetReportData_Precomputation02(t *testing.T) {
	// Create a mock context and request
	ctx := context.Background()
	errorReq := &pb.ReportServiceRequest{
		OrgId:        "00000000-0000-0000-0000-000000000000",
		SubOrgId:     "10000000-0000-0000-0000-000000000000",
		WidgetId:     "s1",
		StartDate:    "2023-06-01 10:00:00",
		EndDate:      "2023-06-30 10:00:00",
		DurationType: 1,
		Component:    []string{"All"},
	}
	// Mock the required dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	handler := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler)

	rsh := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: mockEndpointServiceClient,
	}
	// No data validation for subOrg
	// mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	_, err := rsh.GetReportData(ctx, errorReq)
	expectedErr := fmt.Errorf("rpc error: code = InvalidArgument desc = ReportServiceRequest Validation failed: resource %s does not belong to the organization", errorReq.SubOrgId)
	assert.Equal(t, expectedErr.Error(), err.Error(), "Resource validation failed")

}

func TestReportServiceHandler_GetReportData_Dora01(t *testing.T) {
	// Create a mock context and request
	ctx := context.Background()
	req := &pb.ReportServiceRequest{
		OrgId:        "00000000-0000-0000-0000-000000000000",
		SubOrgId:     "00000000-0000-0000-0000-000000000000",
		WidgetId:     "d1",
		StartDate:    "2023-06-01 10:00:00",
		EndDate:      "2023-06-30 10:00:00",
		DurationType: 1,
		Component:    []string{"All"},
	}
	// Mock the required dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	handler := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler)

	endpointt := &endpoint.Endpoint{
		Id:   "12345000-0000-0000-0000-000000000000",
		Name: "staging",
	}
	endpoints := []*endpoint.Endpoint{endpointt}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpoints,
	}
	mockEndpointServiceClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil).AnyTimes()

	response, _ := handler.GetReportData(ctx, req)
	assert.Equal(t, response.Message, "No Data Found", "Dora metrics no data if environment is not specified validation")

	req.Environment = "12345000-0000-0000-0000-000000000000"
	expectedErr := fmt.Errorf("ReportServiceRequest failed to get data: missing request attributes")
	res, _ := handler.GetReportData(ctx, req)
	assert.Equal(t, res.Error, expectedErr.Error())
}

func TestReportServiceHandler_GetReportData_Dora02(t *testing.T) {
	// Create a mock context and request
	ctx := context.Background()
	errorReq := &pb.ReportServiceRequest{
		SubOrgId:     "10000000-0000-0000-0000-000000000000",
		WidgetId:     "cs1",
		StartDate:    "2023-06-01 10:00:00",
		EndDate:      "2023-06-30 10:00:00",
		DurationType: 1,
		Component:    []string{"All"},
		Environment:  "12345000-0000-0000-0000-000000000000",
	}
	// Mock the required dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	handler := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler)

	rsh := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: mockEndpointServiceClient,
	}
	// mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	expectedErr := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: ReportService Request is missing required field tenantId")
	_, err := rsh.GetReportData(ctx, errorReq)
	assert.Equal(t, err.Error(), expectedErr.Error(), "Basic orgId validation")
}

func TestReportServiceHandler_GetReportData_Dora03(t *testing.T) {
	// Create a mock context and request
	ctx := context.Background()
	req := &pb.ReportServiceRequest{
		OrgId:        "00000000-0000-0000-0000-000000000000",
		SubOrgId:     "00000000-0000-0000-0000-000000000000",
		WidgetId:     "d1",
		StartDate:    "2023-06-01 10:00:00",
		EndDate:      "2023-06-30 10:00:00",
		DurationType: 1,
		Component:    []string{"All"},
		Environment:  "12345000-0000-0000-0000-000000000000",
	}
	// Mock the required dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	handler := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler)

	rsh := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: mockEndpointServiceClient,
	}

	endpointt := &endpoint.Endpoint{
		Id:   "12345000-0000-0000-0000-000000000000",
		Name: "staging",
	}
	endpoints := []*endpoint.Endpoint{endpointt}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpoints,
	}
	mockEndpointServiceClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil).AnyTimes()
	response, _ := rsh.GetReportData(ctx, req)
	assert.Equal(t, response.Status, pb.Status(1), "Validation")
}

func TestReportServiceHandler_GetReportData_Dora04(t *testing.T) {
	// Create a mock context and request
	ctx := context.Background()
	errorReq := &pb.ReportServiceRequest{
		SubOrgId:     "10000000-0000-0000-0000-000000000000",
		WidgetId:     "cs1",
		StartDate:    "2023-06-01 10:00:00",
		EndDate:      "2023-06-30 10:00:00",
		DurationType: 1,
		Component:    []string{"All"},
		Environment:  "12345000-0000-0000-0000-000000000000",
	}
	// Mock the required dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockGrpcClient.EXPECT().Connect(gomock.Any()).Return(nil, nil).AnyTimes()

	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	handler, _ := NewDefaultReportsHandler(mockGrpcClient)
	require.NotNil(t, handler)
	rsh := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: mockEndpointServiceClient,
	}
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	errorReq.OrgId = "00000000-0000-0000-0000-000000000000"
	_, err := rsh.GetReportData(ctx, errorReq)
	expectedErr := fmt.Errorf("rpc error: code = InvalidArgument desc = ReportServiceRequest Validation failed: resource %s does not belong to the organization", errorReq.SubOrgId)
	assert.Equal(t, expectedErr.Error(), err.Error(), "Resource validation failed")
}

func TestReportServiceHandler_GetReportData_05(t *testing.T) {
	// Create a mock context and request
	ctx := context.Background()
	errorReq := &pb.ReportServiceRequest{
		OrgId:        "00000000-0000-0000-0000-000000000000",
		SubOrgId:     "00000000-0000-0000-0000-000000000000",
		WidgetId:     "cs1",
		StartDate:    "2023-06-01 10:00:00",
		EndDate:      "2023-06-30 10:00:00",
		DurationType: 4,
		Component:    []string{"All"},
		Environment:  "12345000-0000-0000-0000-000000000000",
	}
	// Mock the required dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	handler := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler)
	rsh := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: mockEndpointServiceClient,
	}
	//mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	response, _ := rsh.BuildReport(ctx, errorReq)
	assert.Equal(t, response.Status, pb.Status(1), "Validate response got from build reports.")
}

func TestValidateMandatoryFields(t *testing.T) {
	// Test case 1: Missing WdigetId field
	req1 := &pb.ReportServiceRequest{}
	mockCtrl := gomock.NewController(t)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	err := ValidateDataRequest(req1, context.Background(), mockEndpointServiceClient)
	expectedErr1 := fmt.Errorf(errMissingRequiredField, widgetIdField)
	assert.EqualError(t, err, expectedErr1.Error(), "Missing widgetId field should return error")

	// Test case 2:  Missing TenantId field
	req2 := &pb.ReportServiceRequest{
		WidgetId: "1",
	}
	err = ValidateDataRequest(req2, context.Background(), mockEndpointServiceClient)
	expectedErr2 := fmt.Errorf(errMissingRequiredField, tenantIdField)
	assert.EqualError(t, err, expectedErr2.Error(), "Missing TenantId field should return error")

	// Test case 3:  Successful Valiadation
	req6 := &pb.ReportServiceRequest{
		WidgetId:  "1",
		OrgId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		StartDate: "2023-06-01 10:00:00",
		EndDate:   "2023-06-30 10:00:00",
	}
	err = ValidateDataRequest(req6, context.Background(), mockEndpointServiceClient)
	assert.Nil(t, err, "Expected no error")
}

func cleanup(testBed *testutil_setup.TestBed) {
	// Clean up testBed and its client
	testBed.Client.Close()
	testBed.Cleanup()
}

func TestMetricMap(t *testing.T) {
	// Create an instance of ResportServiceHandler
	rsh := &ReportsHandler{
		metrics: &handler.Map{},
	}
	// Call the MetricMap method
	metricMap := rsh.MetricMap()
	// Assert that the returned value is not nil
	assert.NotNil(t, metricMap, "MetricMap should not be nil")
}

func TestHealthy(t *testing.T) {
	// Create an instance of ResportServiceHandler
	rsh := &ReportsHandler{}
	// Call the Healthy method
	err := rsh.Healthy()
	// Assert that the returned value is nil
	assert.NoError(t, err, "Healthy should return nil")
}

func TestDependencies(t *testing.T) {
	// Create an instance of ResportServiceHandler
	rsh := &ReportsHandler{}
	// Call the Dependencies method
	dependencies := rsh.Dependencies()
	// Assert that the returned value is not empty
	assert.NotEmpty(t, dependencies, "Dependencies should not be empty")
}

func TestHealthDependency(t *testing.T) {
	// Create an instance of ResportServiceHandler
	rsh := &ReportsHandler{}
	// Call the HealthDependency method
	healthResponses := rsh.HealthDependency(1, "task")
	// Assert that the returned value is not empty
	assert.NotEmpty(t, healthResponses, "HealthDependency should not be empty")
}

func TestComponentInOrg(t *testing.T) {
	req := &pb.ReportServiceRequest{
		WidgetId:  "1",
		OrgId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		StartDate: "2023-06-01 10:00:00",
		EndDate:   "2023-06-30 10:00:00",
	}
	mockCtrl := gomock.NewController(t)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	//Test 1: Org and Component with same Id
	req.Component = []string{req.OrgId}
	err := ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	assert.Nil(t, err)

	//Test 2: Component All
	req.Component = []string{constants.ALL}
	err = ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	assert.Nil(t, err)
}

// Test for widget builder using mock data in Widget Definition. Supports functions only (not opensearch query)
func TestReportServiceHandler_Widget_FuncSample(t *testing.T) {
	ctx := context.Background()
	req := &pb.ReportServiceRequest{
		WidgetId:     "99",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		Component:    []string{"707f0080-f9bf-4c81-a07d-25cc0fdd9406"},
		DurationType: pb.DurationType_CURRENT_WEEK,
	}

	// Mock the required dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	handler := NewReportsHandler(mockGrpcClient, nil)
	require.NotNil(t, handler)

	tb, err := testutil_setup.NewTestBed(handler)
	defer cleanup(tb)
	if err != nil {
		t.Fatal("error on creating Testbed", err)
	}
	if tb.Server == nil {
		t.Error("should have been able to create the server")
	}

	rsh := &ReportsHandler{
		client: mockGrpcClient,
	}

	w, _ := rsh.BuildReport(ctx, req)
	assert.Equal(t, db.ErrInternalServer.Error(), w.Error)

	mockGrpcClient.EXPECT().Close()
}

func TestReportServiceHandler_BuildReportComponentNonZero(t *testing.T) {
	ctx := context.Background()
	req := &pb.ReportServiceRequest{
		WidgetId:     "99",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		Component:    []string{"707f0080-f9bf-4c81-a07d-25cc0fdd9406", "707f0080-f9bf-4c81-a07d-25cc0fdd9406"},
		DurationType: pb.DurationType_CURRENT_WEEK,
		UserId:       uuid.NewString(),
	}

	// Mock the required dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	handler := NewReportsHandler(mockGrpcClient, nil)
	require.NotNil(t, handler)

	tb, err := testutil_setup.NewTestBed(handler)
	defer cleanup(tb)
	if err != nil {
		t.Fatal("error on creating Testbed", err)
	}
	if tb.Server == nil {
		t.Error("should have been able to create the server")
	}

	rsh := &ReportsHandler{
		client: mockGrpcClient,
	}

	w, _ := rsh.BuildReport(ctx, req)
	// Assert that the Error returned is empty
	assert.Equal(t, db.ErrInternalServer.Error(), w.Error)

	mockGrpcClient.EXPECT().Close()
}

func TestReportServiceHandler_BuildReport(t *testing.T) {
	ctx := context.Background()
	req := &pb.ReportServiceRequest{
		WidgetId:     "99",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		SubOrgId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		Component:    []string{"All"},
		DurationType: pb.DurationType_CURRENT_WEEK,
		UserId:       uuid.NewString(),
	}
	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)

	handler := NewReportsHandler(mockGrpcClient, nil)
	require.NotNil(t, handler)

	tb, err := testutil_setup.NewTestBed(handler)
	defer cleanup(tb)
	if err != nil {
		t.Fatal("error on creating Testbed", err)
	}
	if tb.Server == nil {
		t.Error("should have been able to create the server")
	}
	rsh := &ReportsHandler{
		client:           mockGrpcClient,
		orgServiceClient: orgClient,
	}

	orgClient.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "2cab10cc-cd9d-11ed-afa1-0242ac120002", DisplayName: "Cloudbees Stagging", ParentId: "0"}}, nil).AnyTimes()
	w, _ := rsh.BuildReport(ctx, req)
	assert.Empty(t, w.Error, "widget creation failed")
	mockGrpcClient.EXPECT().Close()
}

func TestReportServiceHandler_BuildReport2(t *testing.T) {
	ctx := context.Background()
	req1 := &pb.ReportServiceRequest{
		WidgetId:     "99",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		SubOrgId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		Component:    []string{"All"},
		DurationType: pb.DurationType_CURRENT_WEEK,
		UserId:       uuid.NewString(),
	}

	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)

	handler := NewReportsHandler(mockGrpcClient, nil)
	require.NotNil(t, handler)

	tb, err := testutil_setup.NewTestBed(handler)
	defer cleanup(tb)
	if err != nil {
		t.Fatal("error on creating Testbed", err)
	}
	if tb.Server == nil {
		t.Error("should have been able to create the server")
	}
	rsh := &ReportsHandler{
		client:           mockGrpcClient,
		orgServiceClient: orgClient,
	}

	serviceResponse := &service.ListServicesResponse{
		Service: []*service.Service{
			{
				Id:            "8f8a7d06-26b0-473a-8e0a-c69557c79a53",
				Name:          "dsl-no-commit-1",
				RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
			},
			{
				Id:            "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "dsl-engine-cli",
				RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
			},
			{
				Id:            "138ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "test",
				RepositoryUrl: "https://github.com/calculi-corp/test.git",
			},
		},
	}
	getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
		return serviceResponse, nil
	}

	orgClient.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "2cab10cc-cd9d-11ed-afa1-0242ac120002", DisplayName: "Cloudbees Stagging", ParentId: "0"}}, nil).AnyTimes()
	w, _ := rsh.BuildReport(ctx, req1)
	assert.Empty(t, w.Error, "widget creation failed")

	mockGrpcClient.EXPECT().Close().Times(1)
	req2 := &pb.ReportServiceRequest{
		WidgetId:     "ci99",
		OrgId:        "a",
		SubOrgId:     "b",
		Component:    []string{"All"},
		DurationType: pb.DurationType_LAST_90_DAYS,
		UserId:       "userid",
	}
	req3 := &pb.ReportServiceRequest{
		WidgetId:     "ci99",
		OrgId:        "a",
		SubOrgId:     "b",
		Component:    []string{"All"},
		DurationType: pb.DurationType_CUSTOM_RANGE,
		UserId:       "userid",
	}
	mockGrpcClient.EXPECT().
		SendGrpc(gomock.Any(), gomock.Any(), "CreateTable", gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	mockCache := cmock.NewMockResourceCacheI(mockCtrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(mockCtrl)

	cache2.SetMockCache(mockCache, mockEpCache, mockGrpcClient)
	mockCache.EXPECT().GetParentIDs(gomock.Any()).Return([]string{"a", "parent_id_4"}).AnyTimes()
	_, err2 := rsh.BuildReport(ctx, req2)
	_, err3 := rsh.BuildReport(ctx, req3)

	assert.Nil(t, err2, "BuildReport should not return an error")
	assert.NotNil(t, err3, "BuildReport should not return an error")

	mockGrpcClient.EXPECT().Close().AnyTimes()
}

func TestReportServiceHandler_calculateDateBydurationType(t *testing.T) {
	// Create a mock context and request
	req := &pb.ReportServiceRequest{
		StartDate:    "2023-30-01 00:00:00",
		EndDate:      "2023-06-30 10:00:00",
		SubOrgId:     "Analytics",
		Component:    []string{"All"},
		WidgetId:     "s6",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		DurationType: pb.DurationType_CURRENT_WEEK,
	}
	mockCtrl := gomock.NewController(t)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	replacements := map[string]any{
		"startDate": req.StartDate,
		"endDate":   req.EndDate,
		"orgId":     req.OrgId,
		"SubOrgId":  req.SubOrgId,
		"component": req.Component,
		"aggrBy":    "day",
		"duration":  "week",
	}

	dateBydurationType := models.CalculateDateBydurationType{
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   map[string]any{"normalizeMonthInSpec": "@x"},
		NormalizeMonthFlag: false,
		IsComputFlag:       false,
	}

	dateBydurationType.CurrentTime = time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC)
	calculateDateBydurationType(dateBydurationType)
	out := time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
	dateBydurationType.CurrentTime = time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC)
	calculateDateBydurationType(dateBydurationType)
	out = time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
	dateBydurationType.CurrentTime = time.Date(2023, 10, 20, 0, 0, 0, 0, time.UTC)
	calculateDateBydurationType(dateBydurationType)
	out = time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
}

func TestReportServiceHandler_calculateDateBydurationTypeMonthCurrent(t *testing.T) {
	// Create a mock context and request
	req := &pb.ReportServiceRequest{
		StartDate:    "2023-30-01 00:00:00",
		EndDate:      "2023-06-30 10:00:00",
		SubOrgId:     "Analytics",
		Component:    []string{"All"},
		WidgetId:     "s6",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		DurationType: pb.DurationType_CURRENT_MONTH,
	}
	mockCtrl := gomock.NewController(t)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	replacements := map[string]any{
		"startDate": req.StartDate,
		"endDate":   req.EndDate,
		"orgId":     req.OrgId,
		"SubOrgId":  req.SubOrgId,
		"component": req.Component,
		"aggrBy":    "day",
		"duration":  "week",
	}

	dateBydurationType := models.CalculateDateBydurationType{
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   map[string]any{"normalizeMonthInSpec": "@x"},
		NormalizeMonthFlag: false,
		IsComputFlag:       false,
		CurrentTime:        time.Date(2023, 10, 20, 0, 0, 0, 0, time.UTC),
	}
	calculateDateBydurationType(dateBydurationType)
	out := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
}

func TestReportServiceHandler_calculateDateBydurationTypeMonthPrevious(t *testing.T) {
	// Create a mock context and request
	req := &pb.ReportServiceRequest{
		StartDate:    "2023-30-01 00:00:00",
		EndDate:      "2023-06-30 10:00:00",
		SubOrgId:     "Analytics",
		Component:    []string{"All"},
		WidgetId:     "s6",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		DurationType: pb.DurationType_PREVIOUS_MONTH,
	}
	mockCtrl := gomock.NewController(t)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	replacements := map[string]any{
		"startDate": req.StartDate,
		"endDate":   req.EndDate,
		"orgId":     req.OrgId,
		"SubOrgId":  req.SubOrgId,
		"component": req.Component,
		"aggrBy":    "day",
		"duration":  "week",
	}

	dateBydurationType := models.CalculateDateBydurationType{
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   map[string]any{"normalizeMonthInSpec": "@x"},
		NormalizeMonthFlag: false,
		IsComputFlag:       false,
	}
	currentTimeJan := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeJan
	calculateDateBydurationType(dateBydurationType)
	out := time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
	currentTimeOthers := time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeOthers
	calculateDateBydurationType(dateBydurationType)
	out = time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
}

func TestReportServiceHandler_calculateDateBydurationTypeWeekPrevious(t *testing.T) {
	// Create a mock context and request
	req := &pb.ReportServiceRequest{
		StartDate:    "2023-30-01 00:00:00",
		EndDate:      "2023-06-30 10:00:00",
		SubOrgId:     "Analytics",
		Component:    []string{"All"},
		WidgetId:     "s6",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		DurationType: pb.DurationType_PREVIOUS_WEEK,
	}
	mockCtrl := gomock.NewController(t)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	replacements := map[string]any{
		"startDate": req.StartDate,
		"endDate":   req.EndDate,
		"orgId":     req.OrgId,
		"SubOrgId":  req.SubOrgId,
		"component": req.Component,
		"aggrBy":    "day",
		"duration":  "week",
	}

	dateBydurationType := models.CalculateDateBydurationType{
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   map[string]any{"normalizeMonthInSpec": "@x"},
		NormalizeMonthFlag: false,
		IsComputFlag:       false,
	}

	currentTimeSunday := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeSunday
	calculateDateBydurationType(dateBydurationType)
	out := time.Date(2023, 9, 18, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
	currentTimeMonday := time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeMonday
	calculateDateBydurationType(dateBydurationType)
	out = time.Date(2023, 9, 25, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
	currentTimeOthers := time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeOthers
	calculateDateBydurationType(dateBydurationType)
	out = time.Date(2023, 9, 25, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
}

func TestReportServiceHandler_calculateDateBydurationTypeTwoMonthBack(t *testing.T) {
	// Create a mock context and request
	req := &pb.ReportServiceRequest{
		StartDate:    "2023-30-01 00:00:00",
		EndDate:      "2023-06-30 10:00:00",
		SubOrgId:     "Analytics",
		Component:    []string{"All"},
		WidgetId:     "s6",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		DurationType: pb.DurationType_TWO_MONTHS_BACK,
	}
	mockCtrl := gomock.NewController(t)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	replacements := map[string]any{
		"startDate": req.StartDate,
		"endDate":   req.EndDate,
		"orgId":     req.OrgId,
		"SubOrgId":  req.SubOrgId,
		"component": req.Component,
		"aggrBy":    "day",
		"duration":  "week",
	}

	dateBydurationType := models.CalculateDateBydurationType{
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   map[string]any{"normalizeMonthInSpec": "@x"},
		NormalizeMonthFlag: false,
		IsComputFlag:       false,
	}
	currentTimeJan := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeJan
	calculateDateBydurationType(dateBydurationType)
	out := time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
	currentTimeFeb := time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeFeb
	calculateDateBydurationType(dateBydurationType)
	out = time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
	currentTimeOthers := time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeOthers
	calculateDateBydurationType(dateBydurationType)
	out = time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
}

func TestReportServiceHandler_calculateDateBydurationTypeTwoWeeksBack(t *testing.T) {
	// Create a mock context and request
	req := &pb.ReportServiceRequest{
		StartDate:    "2023-30-01 00:00:00",
		EndDate:      "2023-06-30 10:00:00",
		SubOrgId:     "Analytics",
		Component:    []string{"All"},
		WidgetId:     "s6",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		DurationType: pb.DurationType_TWO_WEEKS_BACK,
	}
	mockCtrl := gomock.NewController(t)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	replacements := map[string]any{
		"startDate": req.StartDate,
		"endDate":   req.EndDate,
		"orgId":     req.OrgId,
		"SubOrgId":  req.SubOrgId,
		"component": req.Component,
		"aggrBy":    "day",
		"duration":  "week",
	}

	dateBydurationType := models.CalculateDateBydurationType{
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   map[string]any{"normalizeMonthInSpec": "@x"},
		NormalizeMonthFlag: false,
		IsComputFlag:       false,
	}
	currentTimeSunday := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeSunday
	calculateDateBydurationType(dateBydurationType)
	out := time.Date(2023, 9, 11, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
	currentTimeMonday := time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeMonday
	calculateDateBydurationType(dateBydurationType)
	out = time.Date(2023, 9, 18, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
	currentTimeOthers := time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)
	dateBydurationType.CurrentTime = currentTimeOthers
	calculateDateBydurationType(dateBydurationType)
	out = time.Date(2023, 9, 18, 0, 0, 0, 0, time.UTC).Format(timeLayout)
	assert.Equal(t, out, replacements["startDate"], "Expected no error")
}

func TestReportServiceHandler_calculateDateBydurationTypeUnknown(t *testing.T) {
	// Create a mock context and request
	req := &pb.ReportServiceRequest{
		StartDate:    "2023-30-01 00:00:00",
		EndDate:      "2023-06-30 10:00:00",
		SubOrgId:     "Analytics",
		Component:    []string{"All"},
		WidgetId:     "s6",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		DurationType: pb.DurationType_UNKNOWN_DURATION_TYPE,
	}
	replacements := map[string]any{
		"startDate": req.StartDate,
		"endDate":   req.EndDate,
		"orgId":     req.OrgId,
		"SubOrgId":  req.SubOrgId,
		"component": req.Component,
		"aggrBy":    "day",
		"duration":  "week",
	}

	dateBydurationType := models.CalculateDateBydurationType{
		DurationType:       req.DurationType,
		Replacements:       replacements,
		ReplacementsSpec:   map[string]any{"normalizeMonthInSpec": "@x"},
		NormalizeMonthFlag: false,
		IsComputFlag:       false,
		CurrentTime:        time.Now().UTC(),
	}
	// Check for TenantId
	calculateDateBydurationType(dateBydurationType)
	// Assert that the Error returned is empty
	assert.Equal(t, "0001-01-01 00:00:12", replacements["startDate"], "Expected no error")
}

func TestReportServiceHandler_validateDataRequest(t *testing.T) {
	// Create a mock context and request
	req := &pb.ReportServiceRequest{
		StartDate:    "2023-06-01 00:00:00",
		EndDate:      "2023-06-30 10:00:00",
		Component:    []string{"All"},
		WidgetId:     "s6",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		DurationType: pb.DurationType_UNKNOWN_DURATION_TYPE,
	}
	req.StartDate = ""
	mockCtrl := gomock.NewController(t)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	err := ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	expectedErr := fmt.Errorf(errInvalidDateFormatRequest, req.StartDate, req.DurationType.String())
	assert.Equal(t, expectedErr, err, "StartDate validation failed")
	req.StartDate = "2023-01-01 00:00:00"
	req.EndDate = ""
	err = ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	expectedErr = fmt.Errorf(errInvalidDateFormatRequest, req.EndDate, req.DurationType.String())
	assert.Equal(t, expectedErr, err, "EndDate validation failed")
	req.StartDate = "2023-01-01 00:00:00"
	req.EndDate = "2022-12-01 00:00:00"
	err = ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	expectedErr = errors.New(errInvalidRequest)
	assert.Equal(t, expectedErr, err, "EndDate validation failed")
	req.Component = []string{"test"}
	err = ValidateDataRequest(req, context.Background(), mockEndpointServiceClient)
	expectedErr = fmt.Errorf(errResourceNotInOrg, req.Component[0])
	assert.Equal(t, expectedErr, err, "Resource validation failed")

	req.CiToolId = "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7"

}

func TestReportServiceHandler_GetReportData_SubWidget01(t *testing.T) {
	// Create a mock context and request
	ctx := context.Background()
	errorReq := &pb.ReportServiceRequest{
		OrgId:        "00000000-0000-0000-0000-000000000000",
		SubOrgId:     "00000000-0000-0000-0000-000000000000",
		WidgetId:     "e8",
		StartDate:    "2023-06-01 10:00:00",
		EndDate:      "2023-06-30 10:00:00",
		DurationType: 4,
		Component:    []string{"All"},
		Environment:  "12345000-0000-0000-0000-000000000000",
	}
	// Mock the required dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	handler := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler)
	rsh := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: mockEndpointServiceClient,
	}
	response, _ := rsh.BuildReport(ctx, errorReq)
	assert.Equal(t, response.Status, pb.Status(1), "Validate sub widget case.")
}

// e8
func Test_getSubReportWidget(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	handler := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler)
	rsh := &ReportsHandler{
		client:         mockGrpcClient,
		endpointClient: mockEndpointServiceClient,
	}

	specReplacements := make(map[string]any)
	type args struct {
		filterData       []string
		baseData         []string
		req              *pb.ReportServiceRequest
		replacements     map[string]any
		ctx              context.Context
		rah              *ReportsHandler
		replacementsSpec map[string]any
		parentWidget     *pb.Widget
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				filterData: []string{},
				baseData:   []string{},
				req: &pb.ReportServiceRequest{
					WidgetId: "e8",
				},
				replacements: map[string]interface{}{
					"component": []string{"testComp"},
				},
				ctx:              ctx,
				rah:              rsh,
				replacementsSpec: specReplacements,
				parentWidget:     &pb.Widget{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subReportWidget := models.GetSubReportWidget{
				FilterData:       tt.args.filterData,
				BaseData:         tt.args.baseData,
				Req:              tt.args.req,
				Replacements:     tt.args.replacements,
				Ctx:              tt.args.ctx,
				ReplacementsSpec: tt.args.replacementsSpec,
				ParentWidget:     tt.args.parentWidget,
			}

			getSubReportWidget(subReportWidget, tt.args.rah)
			assert.Equal(t, tt.args.parentWidget, &pb.Widget{}, "No error and response validation")
		})
	}
}

func TestReportServiceHandler_BuildDrilldownReport01(t *testing.T) {
	ctx := context.Background()
	req := &pb.DrilldownRequest{
		ReportId:  "component",
		OrgId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		StartDate: "2023-06-01 10:00:00",
		EndDate:   "2023-06-30 10:00:00",
	}
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	handler := NewReportsHandler(mockGrpcClient, nil)
	require.NotNil(t, handler)
	_, err := handler.BuildDrilldownReport(ctx, req)
	assert.Nil(t, err)
}

func TestReportServiceHandler_BuildDrilldownReport02(t *testing.T) {
	ctx := context.Background()
	req := &pb.DrilldownRequest{
		// ReportId: "component", //removed temp
		ReportId:  "workflows",
		OrgId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		StartDate: "2023-06-01 10:00:00",
		EndDate:   "2023-06-30 10:00:00",
		CiToolId:  "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		ReportInfo: &pb.ReportInfo{
			DeploymentEnv: "staging",
			ComponentId:   "compId",
			Code:          "code",
			ScannerName:   "name",
			RunId:         "runId",
		},
	}

	req_ext := &pb.DrilldownRequest{
		ReportId:  "component",
		OrgId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		StartDate: "2023-06-01 10:00:00",
		EndDate:   "2023-06-30 10:00:00",
		CiToolId:  "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		ReportInfo: &pb.ReportInfo{
			DeploymentEnv:   "staging",
			ComponentId:     "compId",
			Code:            "code",
			ScannerName:     "name",
			RunId:           "runId",
			Branch:          "branch",
			Author:          "authorName",
			JobId:           "jobId",
			LicenseType:     "licenseType",
			ScannerNameList: []string{"scanner1", "scanner2", "scanner3"},
			RunIdList:       []string{"runId1", "runId2", "runId3"},
			TestSuiteName:   "testSuiteName",
			TestCaseName:    "testCaseName",
			AutomationId:    "automationId",
		},
	}
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	handler := NewReportsHandler(mockGrpcClient, nil)
	require.NotNil(t, handler)

	getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
		return &endpoint.EndpointsResponse{
			Endpoints: []*endpoint.Endpoint{
				{
					Id:         "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
					Name:       "CJOC Test",
					ResourceId: "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
				},
			},
		}, nil

	}

	getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
		return nil, nil
	}

	// req.SubOrgId = "607f0080-f9bf-4c81-a07d-25cc0fdd9406"
	// mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	// response, _ := handler.BuildDrilldownReport(ctx, req)
	// assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for suborg with no components")

	// response, _ = handler.BuildDrilldownReport(ctx, req)
	// assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for suborg with no components")

	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	req.SubOrgId = req.OrgId
	req.Component = []string{"All"}
	// response, _ := handler.BuildDrilldownReport(ctx, req)
	// assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for suborg with no components")

	req.ReportId = "workflows"
	response, _ := handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for workflows drilldown")

	getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		if IndexName == constants.CB_CI_TOOL_INSIGHT_INDEX || IndexName == constants.CB_CI_JOB_INFO_INDEX {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":null,"hits":[]}}`, nil
		}
		return "", nil
	}
	req.ReportId = "workflowRuns"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for workflowRuns drilldown")

	req_ext.ReportId = "component-summary-workflows"
	response, _ = handler.BuildDrilldownReport(ctx, req_ext)
	assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for component-summary-workflows drilldown")

	req.ReportId = "component-summary-workflows"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for component-summary-workflows drilldown")

	req.ReportId = "component-summary-commits"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for component-summary-commits drilldown")

	req.ReportId = "component-summary-builds"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for component-summary-builds drilldown")

	req.ReportId = "component-summary-deploymentOverview"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for component-summary-deploymentOverview drilldown")

	req.ReportId = "component-summary-workflowRuns"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for component-summary-workflowRuns drilldown")

	// // commented out to prevent breaking of unit test temporarily till time_format added by UI
	// req.ReportId = "security-components"
	// response, _ = handler.BuildDrilldownReport(ctx, req)
	// assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for security-components drilldown")

	req.ReportId = "security-workflows"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for security-workflows drilldown")

	req.ReportId = "security-workflowRuns"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for security-workflowRuns drilldown")

	req.ReportId = "security-scan-type-workflows"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports.Values, []*structpb.Value(nil), "Empty response for security-scan-type-workflows drilldown")

	req.ReportId = "commits"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for commits drilldown")

	req.ReportId = "pullrequests"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for pullrequests drilldown")

	req.ReportId = "runInitiatingCommits"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for runInitiatingCommits drilldown")

	req.ReportId = "builds"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for builds drilldown")

	req.ReportId = "deployments"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for deployments drilldown")

	req.ReportId = "successfulBuildsDuration"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for successfulBuildsDuration drilldown")

	req.ReportId = "deploymentOverview"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for deploymentOverview drilldown")

	req.ReportId = "doraMetrics-deploymentFrequency"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for doraMetrics-deploymentFrequency drilldown")

	req.ReportId = "doraMetrics-deploymentLeadTime"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for doraMetrics-deploymentLeadTime drilldown")

	req.ReportId = "doraMetrics-failureRate"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for doraMetrics-failureRate drilldown")

	req.ReportId = "doraMetrics-mttr"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for doraMetrics-mttr drilldown")

	req.ReportId = "pluginsInfo"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for pluginsInfo drilldown")
	req.ReportId = "runInformation"
	response, _ = handler.BuildDrilldownReport(ctx, req)
	assert.Equal(t, response.Reports, &structpb.ListValue{}, "Empty response for runsInformation drilldown")

}

func TestReportServiceHandler_BuildDrilldownReport03(t *testing.T) {
	ctx := context.Background()
	req := &pb.DrilldownRequest{
		ReportId:  "component",
		OrgId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		StartDate: "2023-06-01 10:00:00",
		EndDate:   "2023-06-30 10:00:00",
		CiToolId:  "b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7",
	}
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	handler := NewReportsHandler(mockGrpcClient, nil)
	require.NotNil(t, handler)
	getAllEndpoints = func(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, parents bool) (*endpoint.EndpointsResponse, error) {
		return &endpoint.EndpointsResponse{
			Endpoints: []*endpoint.Endpoint{
				{
					Id:         "0edebf0d-6797-4ec0-9a50-fac728645a0e",
					Name:       "CJOC Test",
					ResourceId: "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
				},
			},
		}, nil

	}

	_, err := handler.BuildDrilldownReport(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ci tool b7ad0a49-99f9-4982-9940-5f1d3d2dc0f7 does not belong to the organization")
}

func TestDateDiff(t *testing.T) {

	StartDate := "2023-06-01"
	EndDate := "2023-06-30"

	d, err := getDateDiffInDays(StartDate, EndDate)
	assert.Nil(t, err)
	assert.EqualValues(t, 29, d)
}

func TestReportServiceHandler_UpdateDashboardLayout(t *testing.T) {

	ctx := context.Background()
	w := pb.WidgetLayout{I: "s1", W: 1.0, H: 2.5}
	d := pb.DashboardLayout{}
	d.Xl = append(d.Xl, &w)
	d.Lg = append(d.Lg, &w)
	req := &pb.DashboardLayoutRequest{
		OrgId:           "00000000-0000-0000-0000-000000000000",
		UserId:          uuid.NewString(),
		DashboardName:   "software-delivery-activity",
		DashboardLayout: &d,
	}

	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	// Creates mock RCAB client
	mockRbacCtl := rmock.NewMockRBACServiceClient(mockCtrl1)
	assert.NotNil(t, mockRbacCtl)

	handler := NewReportsHandlerWithRbac(mockGrpcClient, mockRbacCtl)
	require.NotNil(t, handler)

	if req.UserId != "" {
		mockRbacCtl.EXPECT().GetContextUserId(ctx, gomock.Any()).Return(&auth.GetContextUserIdResponse{UserId: req.UserId}, nil).Times(1)
	}

	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), hostflags.EndpointServiceHost(), endpointService,
		updateUserPreferenceMethod, gomock.Any(), gomock.Any()).Return(nil)

	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), hostflags.EndpointServiceHost(), endpointService,
		getUserPreferencesMethod, gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// SDA Dashboard Layout
	_, err := handler.UpdateDashboardLayout(ctx, req)
	assert.Nil(t, err)
}

func TestReportServiceHandler_GetDashboardLayout(t *testing.T) {
	ctx := context.Background()
	req := &pb.DashboardLayoutRequest{
		OrgId:         "00000000-0000-0000-0000-000000000000",
		UserId:        uuid.NewString(),
		DashboardName: "software-delivery-activity",
	}

	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	orgClient1 := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient1)

	orgClient1.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "2cab10cc-cd9d-11ed-afa1-0242ac120002", DisplayName: "Cloudbees Stagging", ParentId: "0"}}, nil).AnyTimes()

	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	// Creates mock RCAB client
	mockRbacCtl := rmock.NewMockRBACServiceClient(mockCtrl1)
	assert.NotNil(t, mockRbacCtl)

	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)

	handler := NewReportsHandlerWithRbac(mockGrpcClient, mockRbacCtl)
	require.NotNil(t, handler)

	epClient := mock.NewMockEndpointServiceClient(mockCtrl)

	handler.orgServiceClient = orgClient
	handler.endpointClient = epClient

	mockCoreDataCache := cmock.NewMockResourceCacheI(mockCtrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(mockCtrl)
	cache2.SetMockCache(mockCoreDataCache, mockEpCache, mockGrpcClient)

	orgClient.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "2cab10cc-cd9d-11ed-afa1-0242ac120002", DisplayName: "Cloudbees Stagging", ParentId: "0"}}, nil).AnyTimes()

	if req.UserId != "" {
		mockRbacCtl.EXPECT().GetContextUserId(ctx, gomock.Any()).Return(&auth.GetContextUserIdResponse{UserId: req.UserId}, nil).AnyTimes()
	}

	endpointt := &endpoint.Endpoint{
		Id:   "12345000-0000-0000-0000-000000000000",
		Name: "staging",
	}
	endpointsMock := []*endpoint.Endpoint{endpointt}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpointsMock,
	}

	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(),
		"ListServices", gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(),
		"GetUserPreferences", gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCoreDataCache.EXPECT().Get(gomock.Any()).Return(nil).AnyTimes()
	mockCoreDataCache.EXPECT().GetChildren(gomock.Any()).Return([]string{"res123"}).AnyTimes()
	mockEpCache.EXPECT().GetByContributionID(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	epClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil).AnyTimes()

	// SDA Dashboard Layout
	_, err := handler.GetDashboardLayout(ctx, req)
	assert.Nil(t, err)

	// Security Insights Dashboard Layout
	req.DashboardName = constants.SECURITY_INSIGHTS_DASHBOARD
	_, err = handler.GetDashboardLayout(ctx, req)
	assert.Nil(t, err)

	// Flow metrics Dashboard Layout
	req.DashboardName = constants.FLOW_METRICS_DASHBOARD
	_, err = handler.GetDashboardLayout(ctx, req)
	assert.Nil(t, err)

	// Dora metrics Dashboard Layout
	req.DashboardName = constants.DORA_METRICS_DASHBOARD
	_, err = handler.GetDashboardLayout(ctx, req)
	assert.Nil(t, err)

	// CI Insights Dashboard Layout
	req.DashboardName = constants.CI_INSIGHTS_DASHBOARD
	_, err = handler.GetDashboardLayout(ctx, req)
	assert.Nil(t, err)

	// Test Insights Dashboard Layout
	req.DashboardName = constants.TEST_INSIGHTS_DASHBOARD
	_, err = handler.GetDashboardLayout(ctx, req)
	assert.Nil(t, err)

	// Component Summary Dashboard Layout
	req.DashboardName = constants.COMPONENT_SUMMARY_DASHBOARD
	_, err = handler.GetDashboardLayout(ctx, req)
	assert.Nil(t, err)

	// Invalide Dashboard Name in Request
	req.DashboardName = ""
	_, err = handler.GetDashboardLayout(ctx, req)
	assert.NotNil(t, err)
}

func TestReportServiceHandler_GetWidgets(t *testing.T) {
	ctx := context.Background()
	req := &pb.ManageWidgetRequest{
		OrgId:         "00000000-0000-0000-0000-000000000000",
		DashboardName: "software-delivery-activity",
	}

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	handler := NewReportsHandler(mockGrpcClient, nil)
	require.NotNil(t, handler)

	// SDA Dashboard Layout
	_, err := handler.GetWidgets(ctx, req)
	assert.Nil(t, err)

	// Software Delivery Activity Dashboard
	req.DashboardName = constants.SDA_DASHBOARD
	_, err = handler.GetWidgets(ctx, req)
	assert.Nil(t, err)

	// Security Insights Dashboard
	req.DashboardName = constants.SECURITY_INSIGHTS_DASHBOARD
	_, err = handler.GetWidgets(ctx, req)
	assert.Nil(t, err)

	// Flow metrics Dashboard Layout
	req.DashboardName = constants.FLOW_METRICS_DASHBOARD
	_, err = handler.GetWidgets(ctx, req)
	assert.Nil(t, err)

	// Dora metrics Dashboard Layout
	req.DashboardName = constants.DORA_METRICS_DASHBOARD
	_, err = handler.GetWidgets(ctx, req)
	assert.Nil(t, err)

	// Component Security Dashboard
	req.DashboardName = constants.COMPONENT_SECURITY_DASHBOARD
	_, err = handler.GetWidgets(ctx, req)
	assert.Nil(t, err)

	// Application Security Dashboard
	req.DashboardName = constants.APPLICATION_SECURITY_DASHBOARD
	_, err = handler.GetWidgets(ctx, req)
	assert.Nil(t, err)

	// Component Summary Dashboard
	req.DashboardName = constants.COMPONENT_SUMMARY_DASHBOARD
	_, err = handler.GetWidgets(ctx, req)
	assert.Nil(t, err)

	// Invalid/Incorrect Dashboard Name in Request
	req.DashboardName = "invalid_dashboard_name"
	_, err = handler.GetWidgets(ctx, req)
	assert.NotNil(t, err)

	// Test validation error handling
	req.DashboardName = "valid_dashboard_name"
	req.OrgId = "" // OrgId empty to trigger a validation error
	_, err = handler.GetWidgets(ctx, req)
	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), "ReportServiceRequest Validation failed")
}

func TestReportServiceHandler_BuildComponenComparisonReport(t *testing.T) {
	ctx := context.Background()
	req := &pb.ComponentComparisonRequest{
		WidgetId:     "deployment-frequency-compare",
		OrgId:        "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		SubOrgId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		UserId:       uuid.NewString(),
		StartDate:    "2024-05-01",
		EndDate:      "2024-05-31",
		TimeZone:     "Asia/Calcutta",
		DurationType: pb.DurationType_CURRENT_MONTH,
		Environment:  "12345000-0000-0000-0000-000000000000",
	}
	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)

	defer grpc.SetSharedGrpcClient(nil)

	handler1 := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler1)

	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	endpointt := &endpoint.Endpoint{
		Id:   "12345000-0000-0000-0000-000000000000",
		Name: "staging",
	}
	endpoints := []*endpoint.Endpoint{endpointt}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpoints,
	}
	mockEndpointServiceClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil)

	mockRbacCtl := rmock.NewMockRBACServiceClient(mockCtrl1)
	assert.NotNil(t, mockRbacCtl)

	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)

	handler := &ReportsHandler{
		client:           mockGrpcClient,
		endpointClient:   mockEndpointServiceClient,
		rbacClt:          mockRbacCtl,
		orgServiceClient: orgClient,
	}

	require.NotNil(t, handler)

	mockRbacCtl.EXPECT().GetContextUserId(ctx, gomock.Any()).Return(&auth.GetContextUserIdResponse{UserId: req.UserId}, nil).AnyTimes()

	orgClient.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "2cab10cc-cd9d-11ed-afa1-0242ac120002", DisplayName: "Cloudbees Stagging", ParentId: "0"}}, nil).AnyTimes()
	req.WidgetId = "deployment-frequency-compare"
	response, _ := handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.SubOrgId = req.OrgId
	req.WidgetId = "workflow-runs-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "commit-trends-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "components-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "deployment-success-rate-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "development-cycle-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "pull-requests-trend-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "cycle-time-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "velocity-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "average-active-work-time-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "work-wait-time-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "work-load-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "security-components-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "vulnerabilities-scanner-type-container-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "vulnerabilities-scanner-type-DAST-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "vulnerabilities-very-high-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "vulnerabilities-very-low-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "mttr-vulnerabilities-very-high-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "mttr-vulnerabilities-high-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "mttr-vulnerabilities-low-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "open-vulnerabilities-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "vulnerabilities-scanner-type-SAST-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "vulnerabilities-scanner-type-SCA-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "vulnerabilities-overview-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "build-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "run-initiating-commits-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "security-workflow-runs-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "mttr-vulnerabilities-medium-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "builds-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "successful-deployments-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "security-workflows-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "test-insights-workflows-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)

	req.WidgetId = "test-insights-workflow-runs-compare"
	response, _ = handler.BuildComponentComparisonReport(ctx, req)
	assert.NotNil(t, response)
}

func TestReportGetControllersInfo(t *testing.T) {
	ctx := context.Background()
	req := &pb.CiControllerInfoRequest{
		OrgId: "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
	}
	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)

	defer grpc.SetSharedGrpcClient(nil)

	handler1 := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler1)

	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	endpointt := &endpoint.Endpoint{
		Id:             "12345000-0000-0000-0000-000000000001",
		Name:           "staging",
		ResourceId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		ContributionId: constants.JENKINS_ENDPOINT,
		Audit: &api.Audit{
			When: &timestamppb.Timestamp{
				Seconds: 8727726,
			},
		},
	}

	endpoint1 := &endpoint.Endpoint{
		Id:             "12345000-0000-0000-0000-000000000000",
		Name:           "staging",
		ResourceId:     "12345000-0000-0000-0000-000000000000",
		ContributionId: constants.CBCI_ENDPOINT,
		Properties: []*api.Property{
			{
				Name: "toolUrl",
				Value: &api.Property_String_{
					String_: "https://test2.rosaas.releaseiq.io/",
				},
			},
			{
				Name: "status",
				Value: &api.Property_String_{
					String_: "FAILED",
				},
			},
		},
		Audit: &api.Audit{
			When: &timestamppb.Timestamp{
				Seconds: 8727726,
			},
		},
	}

	endpoint2 := &endpoint.Endpoint{
		Id:             "85bf5c22-7960-440f-8024-492c6ac55845",
		Name:           "staging",
		ResourceId:     "12345000-0000-0000-0000-000000000000",
		ContributionId: constants.CJOC_ENDPOINT,
		Properties: []*api.Property{
			{
				Name: "toolUrl",
				Value: &api.Property_String_{
					String_: "http://example.com/CJOC",
				},
			},
			{
				Name: "status",
				Value: &api.Property_String_{
					String_: "FAILED",
				},
			},
		},
		Audit: &api.Audit{
			When: &timestamppb.Timestamp{
				Seconds: 8727726,
			},
		},
	}

	endpoint3 := &endpoint.Endpoint{
		Id:             "12345000-0000-0000-0000-000000000002",
		Name:           "staging",
		ResourceId:     "12345000-0000-0000-0000-000000000000",
		ContributionId: constants.CBCI_ENDPOINT,
		Properties: []*api.Property{
			{
				Name: "toolUrl",
				Value: &api.Property_String_{
					String_: "",
				},
			},
			{
				Name: "status",
				Value: &api.Property_String_{
					String_: "FAILED",
				},
			},
		},
		Audit: &api.Audit{
			When: &timestamppb.Timestamp{
				Seconds: 8727726,
			},
		},
	}

	endpoints := []*endpoint.Endpoint{endpointt, endpoint1, endpoint2, endpoint3}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpoints,
	}
	mockEndpointServiceClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil).Times(1)

	mockRbacCtl := rmock.NewMockRBACServiceClient(mockCtrl1)
	assert.NotNil(t, mockRbacCtl)

	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)

	handler := &ReportsHandler{
		client:           mockGrpcClient,
		endpointClient:   mockEndpointServiceClient,
		rbacClt:          mockRbacCtl,
		orgServiceClient: orgClient,
	}

	require.NotNil(t, handler)

	controllerMap := make(map[string][]string)

	controllerMap["85bf5c22-7960-440f-8024-492c6ac55845"] = []string{
		"https://testcontroller.rosaas.releaseiq.io/",
		"https://test2.rosaas.releaseiq.io/",
		"https://cb3.rosaas.releaseiq.io/",
		"https://cb2.rosaas.releaseiq.io/",
		"https://controllertest01.rosaas.releaseiq.io/",
	}

	getControllerUrlMap = func(replacements map[string]any, ctx context.Context) (map[string][]string, error) {
		return controllerMap, nil
	}

	jobMap := make(map[string]int)
	runMap := make(map[string]int)
	jobMap["12345000-0000-0000-0000-000000000000"] = 2
	runMap["12345000-0000-0000-0000-000000000000"] = 4

	jobAndRunCount = func(replacements map[string]any, ctx context.Context) (map[string]int, map[string]int, error) {
		return jobMap, runMap, nil
	}

	versionPluginMap := make(map[string]interface{})
	versionPluginMap["12345000-0000-0000-0000-000000000000"] = map[string]interface{}{
		"version": "2.1.1",
		"count":   30.0,
	}
	getVersionAndPluginCount = func(replacements map[string]any, ctx context.Context) (map[string]interface{}, error) {
		return versionPluginMap, nil
	}
	response, _ := handler.GetControllersInfo(ctx, req)
	assert.NotNil(t, response)

	req1 := &pb.CiControllerInfoRequest{
		OrgId:           "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		ContributionIds: []string{constants.JAAS_ENDPOINT},
	}
	endpoint4 := &endpoint.Endpoint{
		Id:             "12345000-0000-0000-0000-000000000002",
		Name:           "staging",
		ResourceId:     "12345000-0000-0000-0000-000000000000",
		ContributionId: constants.CBCI_ENDPOINT,
		Properties: []*api.Property{
			{
				Name: "toolUrl",
				Value: &api.Property_String_{
					String_: "",
				},
			},
			{
				Name: "status",
				Value: &api.Property_String_{
					String_: "FAILED",
				},
			},
		},
		Audit: &api.Audit{
			When: &timestamppb.Timestamp{
				Seconds: 8727726,
			},
		},
	}
	endpoints1 := []*endpoint.Endpoint{endpoint4}
	endpointResp1 := &endpoint.EndpointsResponse{
		Endpoints: endpoints1,
	}
	mockEndpointServiceClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp1, nil).Times(1)

	response1, _ := handler.GetControllersInfo(ctx, req1)
	assert.NotNil(t, response1)

}

func TestReportGetControllersInfoError(t *testing.T) {
	ctx := context.Background()
	req := &pb.CiControllerInfoRequest{
		OrgId: "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
	}
	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)

	defer grpc.SetSharedGrpcClient(nil)

	handler1 := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler1)

	mockEndpointServiceClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)

	mockRbacCtl := rmock.NewMockRBACServiceClient(mockCtrl1)
	assert.NotNil(t, mockRbacCtl)

	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)

	handler := &ReportsHandler{
		client:           mockGrpcClient,
		endpointClient:   mockEndpointServiceClient,
		rbacClt:          mockRbacCtl,
		orgServiceClient: orgClient,
	}

	require.NotNil(t, handler)

	response, err := handler.GetControllersInfo(ctx, req)
	assert.NotNil(t, response)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "test")

}

func TestReportGetInsightsIntegration(t *testing.T) {
	ctx := context.Background()
	req := &pb.CiInsightIntegrationRequest{
		OrgId: "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
	}

	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)

	defer grpc.SetSharedGrpcClient(nil)

	handler1 := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler1)

	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	endpointt := &endpoint.Endpoint{
		Id:             "12345000-0000-0000-0000-000000000001",
		Name:           "staging",
		ResourceId:     "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		ContributionId: constants.JENKINS_ENDPOINT,
		Audit: &api.Audit{
			When: &timestamppb.Timestamp{
				Seconds: 8727726,
			},
		},
	}

	endpoint1 := &endpoint.Endpoint{
		Id:             "12345000-0000-0000-0000-000000000000",
		Name:           "staging",
		ResourceId:     "12345000-0000-0000-0000-000000000000",
		ContributionId: constants.CBCI_ENDPOINT,
		Properties: []*api.Property{
			{
				Name: "toolUrl",
				Value: &api.Property_String_{
					String_: "https://test2.rosaas.releaseiq.io/",
				},
			},
			{
				Name: "status",
				Value: &api.Property_String_{
					String_: "FAILED",
				},
			},
		},
		Audit: &api.Audit{
			When: &timestamppb.Timestamp{
				Seconds: 8727726,
			},
		},
	}

	endpoint2 := &endpoint.Endpoint{
		Id:             "85bf5c22-7960-440f-8024-492c6ac55845",
		Name:           "staging",
		ResourceId:     "12345000-0000-0000-0000-000000000000",
		ContributionId: constants.CJOC_ENDPOINT,
		Properties: []*api.Property{
			{
				Name: "toolUrl",
				Value: &api.Property_String_{
					String_: "http://example.com/CJOC",
				},
			},
			{
				Name: "status",
				Value: &api.Property_String_{
					String_: "FAILED",
				},
			},
		},
		Audit: &api.Audit{
			When: &timestamppb.Timestamp{
				Seconds: 8727726,
			},
		},
	}

	endpoint3 := &endpoint.Endpoint{
		Id:             "12345000-0000-0000-0000-000000000002",
		Name:           "staging",
		ResourceId:     "12345000-0000-0000-0000-000000000000",
		ContributionId: constants.CBCI_ENDPOINT,
		Properties: []*api.Property{
			{
				Name: "toolUrl",
				Value: &api.Property_String_{
					String_: "",
				},
			},
			{
				Name: "status",
				Value: &api.Property_String_{
					String_: "FAILED",
				},
			},
		},
		Audit: &api.Audit{
			When: &timestamppb.Timestamp{
				Seconds: 8727726,
			},
		},
	}

	endpoints := []*endpoint.Endpoint{endpointt, endpoint1, endpoint2, endpoint3}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpoints,
	}
	mockEndpointServiceClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil).Times(1)

	mockRbacCtl := rmock.NewMockRBACServiceClient(mockCtrl1)
	assert.NotNil(t, mockRbacCtl)

	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)

	handler := &ReportsHandler{
		client:           mockGrpcClient,
		endpointClient:   mockEndpointServiceClient,
		rbacClt:          mockRbacCtl,
		orgServiceClient: orgClient,
	}

	require.NotNil(t, handler)

	controllerMap := make(map[string][]string)

	controllerMap["85bf5c22-7960-440f-8024-492c6ac55845"] = []string{
		"https://testcontroller.rosaas.releaseiq.io/",
		"https://test2.rosaas.releaseiq.io/",
		"https://cb3.rosaas.releaseiq.io/",
		"https://cb2.rosaas.releaseiq.io/",
		"https://controllertest01.rosaas.releaseiq.io/",
	}

	getControllerUrlMap = func(replacements map[string]any, ctx context.Context) (map[string][]string, error) {
		return controllerMap, nil
	}

	response, _ := handler.GetInsightsIntegration(ctx, req)
	assert.NotNil(t, response)

	req1 := &pb.CiInsightIntegrationRequest{
		OrgId:           "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
		ContributionIds: []string{constants.JAAS_ENDPOINT},
	}
	endpoint4 := &endpoint.Endpoint{
		Id:             "12345000-0000-0000-0000-000000000002",
		Name:           "staging",
		ResourceId:     "12345000-0000-0000-0000-000000000000",
		ContributionId: constants.CBCI_ENDPOINT,
		Properties: []*api.Property{
			{
				Name: "toolUrl",
				Value: &api.Property_String_{
					String_: "",
				},
			},
			{
				Name: "status",
				Value: &api.Property_String_{
					String_: "FAILED",
				},
			},
		},
		Audit: &api.Audit{
			When: &timestamppb.Timestamp{
				Seconds: 8727726,
			},
		},
	}
	endpoints1 := []*endpoint.Endpoint{endpoint4}
	endpointResp1 := &endpoint.EndpointsResponse{
		Endpoints: endpoints1,
	}
	mockEndpointServiceClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp1, nil).Times(1)

	response1, _ := handler.GetInsightsIntegration(ctx, req1)
	assert.NotNil(t, response1)

}

func TestReportGetInsightsIntegrationError(t *testing.T) {
	ctx := context.Background()
	req := &pb.CiInsightIntegrationRequest{
		OrgId: "707f0080-f9bf-4c81-a07d-25cc0fdd9406",
	}

	mockCtrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	mockEndpointServiceClient := mock.NewMockEndpointServiceClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)

	defer grpc.SetSharedGrpcClient(nil)

	handler1 := NewReportsHandler(mockGrpcClient, mockEndpointServiceClient)
	require.NotNil(t, handler1)

	mockEndpointServiceClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)

	mockRbacCtl := rmock.NewMockRBACServiceClient(mockCtrl1)
	assert.NotNil(t, mockRbacCtl)

	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)

	handler := &ReportsHandler{
		client:           mockGrpcClient,
		endpointClient:   mockEndpointServiceClient,
		rbacClt:          mockRbacCtl,
		orgServiceClient: orgClient,
	}

	require.NotNil(t, handler)

	response, err := handler.GetInsightsIntegration(ctx, req)
	assert.NotNil(t, response)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "test")

}

func Test_getAutomationResponseMapForBranch(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache := coredataMock.NewMockResourceCacheI(mockCtrl)
	cache.CoreDataResourceCache = mockCache
	ctx := context.TODO()

	orgID := "org123"
	components := []string{"238ffe68-8cb4-459d-64ac-2e4f752fe8dc", "8f8a7d06-26b0-473a-8e0a-c69557c79a53"}
	automationSet := map[string]struct{}{"auto1": {}, "auto2": {}}
	excludeDisabledBranch := true
	branchID := "branch123"

	mockCache.EXPECT().GetChildrenOfType(gomock.Any(), gomock.Any()).Return([]string{"1"}).Times(3).AnyTimes()
	resource := api.Resource{Id: "Id1", ParentId: "PId1", Name: "R1", Type: api.ResourceType_RESOURCE_TYPE_AUTOMATION}
	mockCache.EXPECT().Get(gomock.Any()).Return(&resource).Times(2).AnyTimes()

	mockGrpcClient.EXPECT().RetriableGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	responseString := `{
		"aggregations": {
			"component_activity": {
				"value": {
				  "238ffe68-8cb4-459d-64ac-2e4f752fe8dc": {
					"repo_url": "https://github.com/calculi-corp/dsl-engine-cli.git",
					"last_active_time": "2023-10-18T22:09:06.000Z",
					"component_id": "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
					"component_name": "dsl-engine-cli"
				  },
				  "8f8a7d06-26b0-473a-8e0a-c69557c79a53": {
					"repo_url": "https://github.com/sample-gr/dsl-no-commit-1.git",
					"last_active_time": "2023-06-20T23:37:04.000Z",
					"component_id": "8f8a7d06-26b0-473a-8e0a-c69557c79a53",
					"component_name": "dsl-no-commit-1"
				  }
				}
			}
		}
	  }`
	getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		return responseString, nil
	}

	serviceResponse := &service.ListServicesResponse{
		Service: []*service.Service{
			{
				Id:            "8f8a7d06-26b0-473a-8e0a-c69557c79a53",
				Name:          "dsl-no-commit-1",
				RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
			},
			{
				Id:            "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "dsl-engine-cli",
				RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
			},
			{
				Id:            "138ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "test",
				RepositoryUrl: "https://github.com/calculi-corp/test.git",
			},
		},
	}
	getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
		return serviceResponse, nil
	}

	getAutomationResponseMapForBranch(ctx, mockGrpcClient, orgID, components, automationSet, excludeDisabledBranch, branchID)

}

func Test_getAutomationResponseMap(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache := coredataMock.NewMockResourceCacheI(mockCtrl)
	cache.CoreDataResourceCache = mockCache
	ctx := context.TODO()

	orgID := "org123"
	components := []string{"238ffe68-8cb4-459d-64ac-2e4f752fe8dc", "8f8a7d06-26b0-473a-8e0a-c69557c79a53"}
	automationSet := map[string]struct{}{"auto1": {}, "auto2": {}}

	mockCache.EXPECT().GetChildrenOfType(gomock.Any(), gomock.Any()).Return([]string{"1"}).Times(3).AnyTimes()
	audit := &api.Audit{
		Why: "TestReason",
	}
	resource := api.Resource{
		Id:         "Id1",
		ParentId:   "PId1",
		Name:       "R1",
		Type:       api.ResourceType_RESOURCE_TYPE_DASHBOARD,
		Audit:      audit,
		IsDisabled: false,
	}
	mockCache.EXPECT().Get(gomock.Any()).Return(&resource).Times(2).AnyTimes()

	mockGrpcClient.EXPECT().RetriableGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	responseString := `{
		"aggregations": {
			"component_activity": {
				"value": {
				  "238ffe68-8cb4-459d-64ac-2e4f752fe8dc": {
					"repo_url": "https://github.com/calculi-corp/dsl-engine-cli.git",
					"last_active_time": "2023-10-18T22:09:06.000Z",
					"component_id": "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
					"component_name": "dsl-engine-cli"
				  },
				  "8f8a7d06-26b0-473a-8e0a-c69557c79a53": {
					"repo_url": "https://github.com/sample-gr/dsl-no-commit-1.git",
					"last_active_time": "2023-06-20T23:37:04.000Z",
					"component_id": "8f8a7d06-26b0-473a-8e0a-c69557c79a53",
					"component_name": "dsl-no-commit-1"
				  }
				}
			}
		}
	  }`
	getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		return responseString, nil
	}

	serviceResponse := &service.ListServicesResponse{
		Service: []*service.Service{
			{
				Id:            "8f8a7d06-26b0-473a-8e0a-c69557c79a53",
				Name:          "dsl-no-commit-1",
				RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
			},
			{
				Id:            "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "dsl-engine-cli",
				RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
			},
			{
				Id:            "138ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "test",
				RepositoryUrl: "https://github.com/calculi-corp/test.git",
			},
		},
	}
	getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
		return serviceResponse, nil
	}

	getAutomationResponseMap(ctx, mockGrpcClient, orgID, components, automationSet, true)

}

type OpenSearchResponse struct {
	key  string
	data string
	err  error
}

func Test_automationReportForBranch(t *testing.T) {

	testReplacements := map[string]any{
		"orgId":     "00000000-0000-0000-0000-000000000000",
		"subOrgId":  "00000000-0000-0000-0000-000000000000",
		"component": []string{"All"},
		"startDate": "2023-06-01 10:00:00",
		"endDate":   "2023-06-30 10:00:00",
		"timeZone":  "Asia/Calcutta",
		"branch":    "main",
	}
	ctx := context.Background() // Use a context for testing.

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCache := coredataMock.NewMockResourceCacheI(mockCtrl)
	cache.CoreDataResourceCache = mockCache

	mockCache.EXPECT().GetChildrenOfType(gomock.Any(), gomock.Any()).Return([]string{"1"}).Times(3).AnyTimes()
	resource := api.Resource{Id: "Id1", ParentId: "PId1", Name: "R1", Type: api.ResourceType_RESOURCE_TYPE_AUTOMATION}
	mockCache.EXPECT().Get(gomock.Any()).Return(&resource).Times(2).AnyTimes()

	mockGrpcClient.EXPECT().RetriableGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	responseString := `{
		"aggregations": {
			"component_activity": {
				"value": {
				  "238ffe68-8cb4-459d-64ac-2e4f752fe8dc": {
					"repo_url": "https://github.com/calculi-corp/dsl-engine-cli.git",
					"last_active_time": "2023-10-18T22:09:06.000Z",
					"component_id": "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
					"component_name": "dsl-engine-cli"
				  },
				  "8f8a7d06-26b0-473a-8e0a-c69557c79a53": {
					"repo_url": "https://github.com/sample-gr/dsl-no-commit-1.git",
					"last_active_time": "2023-06-20T23:37:04.000Z",
					"component_id": "8f8a7d06-26b0-473a-8e0a-c69557c79a53",
					"component_name": "dsl-no-commit-1"
				  }
				}
			}
		}
	  }`
	getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		return responseString, nil
	}

	serviceResponse := &service.ListServicesResponse{
		Service: []*service.Service{
			{
				Id:            "8f8a7d06-26b0-473a-8e0a-c69557c79a53",
				Name:          "dsl-no-commit-1",
				RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
			},
			{
				Id:            "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "dsl-engine-cli",
				RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
			},
			{
				Id:            "138ffe68-8cb4-459d-64ac-2e4f752fe8dc",
				Name:          "test",
				RepositoryUrl: "https://github.com/calculi-corp/test.git",
			},
		},
	}
	getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
		return serviceResponse, nil
	}

	AutomationReportForBranch(testReplacements, ctx, mockGrpcClient)

}

func TestGetSiTransitionConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(ctrl)
	mockOrgServiceClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	ctx := context.Background()
	req := &pb.DashboardLayoutRequest{
		OrgId:     "org123",
		UserId:    uuid.NewString(),
		Component: "comp1",
	}

	mockOrgServiceClient.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(
		&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "org123"}}, nil).Times(1)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockCoreDataCache := cmock.NewMockResourceCacheI(ctrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(ctrl)
	cache.SetMockCache(mockCoreDataCache, mockEpCache, mockGrpcClient)

	r := &api.Resource{
		Id:   "123",
		Type: api.ResourceType_RESOURCE_TYPE_AUTOMATION,
	}
	mockCoreDataCache.EXPECT().Get(gomock.Any()).Return(r).AnyTimes()
	mockCoreDataCache.EXPECT().GetChildren(gomock.Any()).Return([]string{"res123"}).AnyTimes()

	handler := &ReportsHandler{
		client:           mockGrpcClient,
		orgServiceClient: mockOrgServiceClient,
	}
	handler.GetSiTransitionConfig(ctx, req)
}

func TestGetTestInsightsTransitionConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(ctrl)
	mockOrgServiceClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	ctx := context.Background()
	req := &pb.DashboardLayoutRequest{
		OrgId:     "org123",
		UserId:    uuid.NewString(),
		Component: "comp1",
	}

	mockOrgServiceClient.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(
		&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "org123"}}, nil).Times(1)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockCoreDataCache := cmock.NewMockResourceCacheI(ctrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(ctrl)
	cache.SetMockCache(mockCoreDataCache, mockEpCache, mockGrpcClient)

	r := &api.Resource{
		Id:   "123",
		Type: api.ResourceType_RESOURCE_TYPE_AUTOMATION,
	}
	mockCoreDataCache.EXPECT().Get(gomock.Any()).Return(r).AnyTimes()
	mockCoreDataCache.EXPECT().GetChildren(gomock.Any()).Return([]string{"res123"}).AnyTimes()

	handler := &ReportsHandler{
		client:           mockGrpcClient,
		orgServiceClient: mockOrgServiceClient,
	}
	handler.GetTestInsightsTransitionConfig(ctx, req)
}
