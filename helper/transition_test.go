package helper

import (
	"context"
	"testing"

	api "github.com/calculi-corp/api/go"
	"github.com/calculi-corp/api/go/auth"
	"github.com/calculi-corp/api/go/endpoint"
	cmock "github.com/calculi-corp/core-data-cache/mock"
	mock_grpc_client "github.com/calculi-corp/grpc-client/mock"
	"github.com/calculi-corp/reports-service/cache"
	"github.com/calculi-corp/reports-service/constants"
	rmock "github.com/calculi-corp/reports-service/mock"
	remock "github.com/calculi-corp/repository-service/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetComponents(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCtrl1 := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(ctrl)

	orgClient := rmock.NewMockOrganizationsServiceClient(mockCtrl1)
	assert.NotNil(t, orgClient)

	ctx := context.Background()
	orgID := "org123"
	userID := uuid.NewString()

	orgClient.EXPECT().GetOrganizationById(gomock.Any(), gomock.Any()).Return(&auth.GetOrganizationByIdResponse{Organization: &auth.Organization{Id: "2cab10cc-cd9d-11ed-afa1-0242ac120002", DisplayName: "Cloudbees Stagging", ParentId: "0"}}, nil).Times(1)

	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	response := GetComponents(ctx, mockGrpcClient, orgClient, orgID, userID)

	assert.NotNil(t, response, "Expected response to be not nil")

}

func TestGetWorkflowsCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(ctrl)

	mockCoreDataCache := cmock.NewMockResourceCacheI(ctrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(ctrl)
	cache.SetMockCache(mockCoreDataCache, mockEpCache, mockGrpcClient)

	r := &api.Resource{
		Id:   "123",
		Type: api.ResourceType_RESOURCE_TYPE_AUTOMATION,
	}

	mockCoreDataCache.EXPECT().Get(gomock.Any()).Return(r).AnyTimes()
	mockCoreDataCache.EXPECT().GetChildren(gomock.Any()).Return([]string{"res123"}).AnyTimes()

	ctx := context.Background()
	orgID := "org123"
	components := []string{"comp1", "com2"}

	workflowsCount := GetWorkflowsCount(ctx, orgID, components)

	assert.Equal(t, 2, workflowsCount)
}

func TestGetIndexDocCount(t *testing.T) {
	orgID := "org123"
	components := []string{"comp1", "com2"}

	workflowRunCount := GetIndexDocCount(constants.FLOW_METRICS_INDEX, orgID, components)
	assert.Equal(t, 0, workflowRunCount)
}

func TestGetIndexDocCountByOrgId(t *testing.T) {
	orgID := "org123"
	workflowRunCount := GetIndexDocCountByOrgId(constants.FLOW_METRICS_INDEX, orgID)
	assert.Equal(t, 0, workflowRunCount)
}

func TestGetResourceChildrenByType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(ctrl)

	mockCoreDataCache := cmock.NewMockResourceCacheI(ctrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(ctrl)
	cache.SetMockCache(mockCoreDataCache, mockEpCache, mockGrpcClient)

	r := &api.Resource{
		Id:   "123",
		Type: api.ResourceType_RESOURCE_TYPE_AUTOMATION,
	}

	mockCoreDataCache.EXPECT().Get(gomock.Any()).Return(r).AnyTimes()
	mockCoreDataCache.EXPECT().GetChildren(gomock.Any()).Return([]string{"res123"}).AnyTimes()

	resourceIds := []string{"res123", "res456"}
	resourceType := api.ResourceType_RESOURCE_TYPE_AUTOMATION

	resources := GetResourceChildrenByType(resourceIds, 100, false, resourceType)
	assert.Equal(t, 1, len(resources))
}

func TestGetFlowItemMappings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(ctrl)

	mockCoreDataCache := cmock.NewMockResourceCacheI(ctrl)
	mockEpCache := cmock.NewMockEndpointsCacheI(ctrl)
	cache.SetMockCache(mockCoreDataCache, mockEpCache, mockGrpcClient)

	epClinet := remock.NewMockEndpointServiceClient(ctrl)
	ctx := context.Background()

	endpointt := &endpoint.Endpoint{
		Id:   "12345000-0000-0000-0000-000000000000",
		Name: "staging",
	}
	endpointsMock := []*endpoint.Endpoint{endpointt}
	endpointResp := &endpoint.EndpointsResponse{
		Endpoints: endpointsMock,
	}
	epClinet.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(endpointResp, nil).Times(1)
	orgID := "org123"
	endpoints := GetFlowItemMappings(ctx, epClinet, orgID)
	assert.Equal(t, 0, endpoints)
}
