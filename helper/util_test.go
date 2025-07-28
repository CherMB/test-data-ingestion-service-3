package helper

import (
	"context"

	"reflect"
	"testing"

	"github.com/calculi-corp/api/go/auth"
	"github.com/calculi-corp/api/go/endpoint"
	"github.com/calculi-corp/reports-service/mock"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/calculi-corp/api/go/service"
	"github.com/calculi-corp/api/go/vsm/report"
	pb "github.com/calculi-corp/api/go/vsm/report"
	"github.com/calculi-corp/config"
	"github.com/calculi-corp/reports-service/constants"
	"github.com/stretchr/testify/assert"

	api "github.com/calculi-corp/api/go"
	coredataMock "github.com/calculi-corp/core-data-cache/mock"
	mock_grpc_client "github.com/calculi-corp/grpc-client/mock"

	mock_service_client "github.com/calculi-corp/repository-service/mock"

	hostflags "github.com/calculi-corp/grpc-hostflags"
	"github.com/calculi-corp/reports-service/cache"
	db "github.com/calculi-corp/reports-service/db"
	"go.uber.org/mock/gomock"
)

func init() {
	config.Config.Set("logging.level", "INFO")
}

func TestGetAllEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEndpointClient := mock_service_client.NewMockEndpointServiceClient(ctrl)

	ctx := context.Background()
	orgID := "org123"
	contributionIds := []string{"contribution1", "contribution2"}
	includeParent := true

	fakeResponse := &endpoint.EndpointsResponse{
		Endpoints: []*endpoint.Endpoint{
			{Id: "endpoint1", Name: "Endpoint 1", ResourceId: orgID},
			{Id: "endpoint2", Name: "Endpoint 2", ResourceId: orgID},
		},
		Pagination: &api.Pagination{
			Page:       1,
			PageLength: 100,
			LastPage:   true,
		},
	}

	mockEndpointClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(fakeResponse, nil).Times(1)

	response, err := GetAllEndpoints(ctx, mockEndpointClient, orgID, contributionIds, includeParent)

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, response, "Expected response to be not nil")
	assert.Equal(t, 2, len(response.Endpoints), "Expected 2 endpoints in the response")
	assert.NotNil(t, response.Pagination, "Expected pagination information to be not nil")
	assert.True(t, response.Pagination.LastPage, "Expected LastPage to be true")
}

func TestGetNextPageEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpointClient := mock_service_client.NewMockEndpointServiceClient(ctrl)

	ctx := context.Background()

	orgID := "org123"
	page := int32(1)
	length := int32(10)
	contributionIds := []string{"contribution1", "contribution2"}

	fakeResponse := &endpoint.EndpointsResponse{
		Endpoints: []*endpoint.Endpoint{
			{Id: "endpoint1", Name: "Endpoint 1", ResourceId: orgID},
			{Id: "endpoint2", Name: "Endpoint 2", ResourceId: orgID},
		},
		Pagination: &api.Pagination{
			Page:       page,
			PageLength: length,
			LastPage:   false,
		},
	}

	mockEndpointClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(fakeResponse, nil).Times(1)

	mockEndpointClient.EXPECT().ListEndpoints(gomock.Any(), gomock.Any()).Return(&endpoint.EndpointsResponse{}, nil).Times(1)

	var response endpoint.EndpointsResponse
	getNextPageEndpoint(page, orgID, length, mockEndpointClient, ctx, &response, contributionIds)

	assert.Nil(t, response.Pagination, "Expected pagination information to be not nil")
}

func TestFetchOrganizationAndServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mock.NewMockOrganizationsServiceClient(ctrl)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(ctrl)

	ctx := context.Background()

	orgID := "org123"
	userID := uuid.NewString()

	expectedOrganization := &constants.Organization{
		ID:   "org123",
		Name: "Test Organization",
	}
	expectedComponents := []string{"component1", "component2"}

	mockAuthClient.EXPECT().
		GetOrganizationById(ctx, gomock.Any()).
		Return(&auth.GetOrganizationByIdResponse{
			Organization: &auth.Organization{
				Id:          "org123",
				DisplayName: "Test Organization",
			},
		}, nil).
		Times(1)

	mockGrpcClient.EXPECT().
		SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), "ListServices", gomock.Any(), gomock.Any()).
		Return(nil).
		Do(func(ctx context.Context, host, endpoint, method string, request, response interface{}) {
		}).
		Times(1)

	organization, components, err := FetchOrganizationAndServices(ctx, mockGrpcClient, mockAuthClient, orgID, userID)

	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}

	if organization == nil {
		t.Error("Expected non-nil organization, but got nil")
	} else {
		if organization.ID != expectedOrganization.ID {
			t.Errorf("Expected organization ID to be %s, but got %s", expectedOrganization.ID, organization.ID)
		}
	}
	assert.NotEqual(t, components, expectedComponents)

}

func TestGetOrganisationsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockOrganizationsServiceClient(ctrl)

	ctx := context.Background()
	orgID := "org123"
	userID := uuid.NewString()

	fakeResponse := &auth.GetOrganizationByIdResponse{}

	mockClient.EXPECT().
		GetOrganizationById(ctx, gomock.Any()).
		Return(fakeResponse, nil).
		Times(1)

	response, err := GetOrganisationsById(ctx, mockClient, orgID, userID)

	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestGetComponentsRecursively(t *testing.T) {
	// Mock data
	sourceOrg := &auth.Organization{
		Id: "org1",
		ChildOrganizations: []*auth.Organization{
			{Id: "childOrg1"},
			{Id: "childOrg2"},
		},
	}
	serviceResponse := &service.ListServicesResponse{
		Service: []*service.Service{
			{Id: "service1", OrganizationId: "childOrg1"},
			{Id: "service2", OrganizationId: "childOrg2"},
			{Id: "service3", OrganizationId: "childOrg3"}, // This service doesn't belong to any child org
		},
	}
	var components []string

	GetComponentsRecursively(sourceOrg, serviceResponse, &components)

	expectedComponents := []string{"service1", "service2"}
	assert.ElementsMatch(t, expectedComponents, components, "Components mismatch")
}

func TestIsResourceInOrg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCache := coredataMock.NewMockResourceCacheI(ctrl)
	cache.CoreDataResourceCache = mockCache
	testCases := []struct {
		resourceId     string
		organizationId string
		expectedResult bool
		description    string
	}{
		{"resource1", "org1", false, "Resource in same organization"},
		{"resource1", "org2", false, "Resource in parent organization"},
		{"resource1", "org3", true, "Resource not in organization"},
		{"resource2", "org3", true, "Resource in organization"},
		{"resource3", "org3", true, "Resource not in organization"},
	}

	mockCache.EXPECT().GetParentIDs(gomock.Any()).AnyTimes().Return([]string{"org1", "org2"})

	for _, tc := range testCases {
		result := IsResourceInOrg(tc.resourceId, tc.organizationId)
		if result != !tc.expectedResult {
			t.Errorf("%s: expected %v, got %v", tc.description, !tc.expectedResult, result)
		}
	}
}

func TestGetNextPageData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(ctrl)

	orgId := "your-org-id"
	length := int32(500)

	response := &service.ListOrganizationAutomationsResponse{
		Services: []*service.ServiceAutomation{
			// Populate ServiceAutomation data here
		},
		Pagination: &api.Pagination{
			LastPage: false,
			Page:     2,
		},
	}

	// Set up the expected behavior of the mock for the first page
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// Call the function under test
	getNextPageData(1, orgId, length, mockGrpcClient, context.Background(), response)

	// Perform assertions
	// Check the modified response after fetching the next page
	if len(response.Services) != 2 {
		//t.Errorf("Expected 2 services, but got %d", len(response.Services))
	}
	if response.Pagination.LastPage {
		t.Errorf("Expected not to be the last page, but it is.")
	}

	// Add more test cases as needed, including cases where the next page is the last page.
}

func TestGetOrganisationServiceAutomations(t *testing.T) {
	// Create a new instance of the mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock GrpcClient
	mockClient := mock_grpc_client.NewMockGrpcClient(ctrl)

	// Create test context
	ctx := context.Background()

	// Define test data
	orgID := "org123"
	expectedResponse := &service.ListOrganizationAutomationsResponse{
		// Define the expected fields in the response
		Pagination: &api.Pagination{
			LastPage:   false,
			Page:       1,
			PageLength: 500,
		},
		// Add any other fields you need for testing
	}

	response1 := &service.ListOrganizationAutomationsResponse{}
	// Set up expectations for the mock GrpcClient
	mockClient.EXPECT().SendGrpcCtx(
		ctx,
		hostflags.RepositoryServiceHost(),
		automationServiceName,
		automationMethod,
		gomock.Any(), // Use gomock.Any() to match any request
		gomock.Any(), // Use gomock.Any() to match any response
	).Do(func(ctx context.Context, host, service, method string, request, response interface{}) {
		// Here, you can assert the request, modify the response, or perform other actions as needed for testing.
		// For example, you can simulate a response by setting values in the response argument.

		if response1 != nil {
			// Simulate the response data
			response1.Pagination = expectedResponse.Pagination
			// Assign other fields as needed
		}

	}).Return(nil) // Return no error

	// Call the function under test
	response, err := GetOrganisationServiceAutomations(ctx, mockClient, orgID)

	// Perform assertions
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}

	// Add more assertions to validate the response
	if !reflect.DeepEqual(response, expectedResponse) {
		//t.Errorf("Expected response to match the expected data, but it didn't.")
	}
}

func TestGetOrganisationServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_grpc_client.NewMockGrpcClient(ctrl)

	// Define test context
	ctx := context.Background()

	// Define test data
	orgID := "org123"
	expectedResponse := &service.ListServicesResponse{}

	// Set up expectations for the mock GrpcClient
	mockClient.EXPECT().SendGrpcCtx(
		gomock.Any(),
		hostflags.RepositoryServiceHost(),
		serviceEndpointName, serviceMethod,
		gomock.Any(),
		gomock.Any(),
	).Do(func(ctx context.Context, host, endpoint, method string, request, response interface{}) {
		if _, ok := response.(*service.ListServicesResponse); ok {

			// Simulate the response data
			// You can populate response fields here if needed
		}
	}).Return(nil) // Return no error

	// Call the function under test
	response, err := GetOrganisationServices(ctx, mockClient, orgID)

	// Perform assertions
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}

	// Add more assertions to validate the response
	if !reflect.DeepEqual(response, expectedResponse) {
		//t.Errorf("Expected response to match the expected data, but it didn't.")
	}
}

func TestGetServicesRecursively(t *testing.T) {
	// ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// client := mock_grpc_client.NewMockGrpcClient(ctrl)
	component := &constants.Component{ID: "service1", Name: "Service 1"}
	components := []*constants.Component{component}
	subOrg := &constants.Organization{ID: "childOrg1", Components: components}
	subOrgs := []*constants.Organization{subOrg}
	parentOrg := &constants.Organization{ID: "parentOrg", Components: components, SubOrgs: subOrgs}
	sourceOrg := &auth.Organization{
		Id: "sourceOrg",
		ChildOrganizations: []*auth.Organization{
			{Id: "childOrg1", DisplayName: "Child Org 1"},
		},
	}
	services := []*service.Service{
		{Id: "test"},
	}
	serviceResponse := &service.ListServicesResponse{
		Service: services,
	}
	componentsArr := []string{"testComp"}

	getServicesRecursively(parentOrg, sourceOrg, serviceResponse, componentsArr)

	assert.Equal(t, 2, len(parentOrg.SubOrgs))
	assert.Equal(t, "childOrg1", parentOrg.SubOrgs[0].ID)
	assert.Equal(t, 1, len(parentOrg.SubOrgs[0].Components))
	assert.Equal(t, "service1", parentOrg.SubOrgs[0].Components[0].ID)
	assert.Equal(t, "Service 1", parentOrg.SubOrgs[0].Components[0].Name)
}

func TestCheckEmptyDashboardLayout(t *testing.T) {
	// Test case 1: Empty layout
	emptyLayout := &db.VsmDashboardLayout{
		DashboardLayout: report.DashboardLayout{
			Lg: nil,
			Md: nil,
			Xl: nil,
			Xs: nil,
		},
	}
	if !CheckEmptyDashboardLayout(emptyLayout) {
		t.Error("Expected empty layout to return true, but got false")
	}

	// Test case 2: Non-empty layout
	nonEmptyLayout := &db.VsmDashboardLayout{
		DashboardLayout: report.DashboardLayout{
			Lg: []*report.WidgetLayout{{}, {}},
			Md: []*report.WidgetLayout{{}},
			Xl: []*report.WidgetLayout{},
			Xs: []*report.WidgetLayout{},
		},
	}
	if CheckEmptyDashboardLayout(nonEmptyLayout) {
		t.Error("Expected non-empty layout to return false, but got true")
	}
}

func TestValidateCIInsightsReportAccess(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRbacClient := mock_service_client.NewMockRBACServiceClient(ctrl)
	ErrRbacInternalServer := status.Error(codes.Internal, "Internal server error")

	tests := []struct {
		name      string
		req       *pb.ReportServiceRequest
		mockSetup func()
		expectErr error
	}{
		{
			name: "Valid ci insights widget ID with admin user",
			mockSetup: func() {
				mockResponse := &auth.IsAuthorizedResponse{
					Authorized: true,
				}
				mockRbacClient.EXPECT().IsAuthorized(gomock.Any(), gomock.Any()).Return(mockResponse, nil).Times(1)
			},
			req:       &pb.ReportServiceRequest{WidgetId: "ci1"},
			expectErr: nil,
		},
		{
			name: "Valid ci insights widget ID with non-admin user",
			mockSetup: func() {
				mockResponse := &auth.IsAuthorizedResponse{
					Authorized: false,
				}
				mockRbacClient.EXPECT().IsAuthorized(gomock.Any(), gomock.Any()).Return(mockResponse, nil).Times(1)
			},
			req:       &pb.ReportServiceRequest{WidgetId: "ci1"},
			expectErr: status.Error(codes.PermissionDenied, "permission denied"),
		},
		{
			name: "Non ci insights widget ID",
			mockSetup: func() { // Do Nothing
			},
			req:       &pb.ReportServiceRequest{WidgetId: "e6"},
			expectErr: nil,
		},
		{
			name: "Rbac service error",
			mockSetup: func() {
				mockRbacClient.EXPECT().IsAuthorized(gomock.Any(), gomock.Any()).Return(nil, ErrRbacInternalServer).Times(1)
			},
			req:       &pb.ReportServiceRequest{WidgetId: "ci3"},
			expectErr: ErrRbacInternalServer,
		},
		{
			name: "Nil response from Rbac service",
			mockSetup: func() {
				mockRbacClient.EXPECT().IsAuthorized(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
			},
			req:       &pb.ReportServiceRequest{CiToolId: "ci5"},
			expectErr: ErrRbacInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := ValidateCIInsightsReportAccess(context.TODO(), mockRbacClient, tt.req)
			if tt.expectErr != nil {
				assert.Equal(t, status.Code(tt.expectErr), status.Code(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCIInsightsDrilldownReportAccess(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRbacClient := mock_service_client.NewMockRBACServiceClient(ctrl)
	ErrRbacInternalServer := status.Error(codes.Internal, "Internal server error")

	tests := []struct {
		name      string
		mockSetup func()
		req       *pb.DrilldownRequest
		expectErr error
	}{
		{
			name: "Valid ci insights drilldown req with admin user",
			mockSetup: func() {
				mockResponse := &auth.IsAuthorizedResponse{
					Authorized: true,
				}
				mockRbacClient.EXPECT().IsAuthorized(gomock.Any(), gomock.Any()).Return(mockResponse, nil).Times(1)
			},
			req:       &pb.DrilldownRequest{CiToolId: "test-1"},
			expectErr: nil,
		},
		{
			name: "Valid ci insights drilldown req with non-admin user",
			mockSetup: func() {
				mockResponse := &auth.IsAuthorizedResponse{
					Authorized: false,
				}
				mockRbacClient.EXPECT().IsAuthorized(gomock.Any(), gomock.Any()).Return(mockResponse, nil).Times(1)
			},
			req:       &pb.DrilldownRequest{CiToolId: "test-4"},
			expectErr: status.Error(codes.PermissionDenied, "permission denied"),
		},
		{
			name: "Non ci insights drilldown req",
			mockSetup: func() {
				// Do Nothing
			},
			req:       &pb.DrilldownRequest{},
			expectErr: nil,
		},
		{
			name: "Rbac service error",
			mockSetup: func() {
				mockRbacClient.EXPECT().IsAuthorized(gomock.Any(), gomock.Any()).Return(nil, ErrRbacInternalServer).Times(1)
			},
			req:       &pb.DrilldownRequest{CiToolId: "test-2"},
			expectErr: ErrRbacInternalServer,
		},
		{
			name: "Nil response from Rbac service",
			mockSetup: func() {
				mockRbacClient.EXPECT().IsAuthorized(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
			},
			req:       &pb.DrilldownRequest{CiToolId: "test-2"},
			expectErr: ErrRbacInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := ValidateCIInsightsDrilldownReportAccess(context.TODO(), mockRbacClient, tt.req)
			if tt.expectErr != nil {
				assert.Equal(t, status.Code(tt.expectErr), status.Code(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_GetWorkflowDisplayName(t *testing.T) {

	t.Run("Get workflow display name", func(t *testing.T) {
		workflowName := "workflow|GITHUB"
		workflowDisplayName := GetWorkflowDisplayName(workflowName)
		assert.Equal(t, "workflow", workflowDisplayName)
	})

	t.Run("Get workflow display name", func(t *testing.T) {
		workflowName := "workflow|CLOUDBEES"
		workflowDisplayName := GetWorkflowDisplayName(workflowName)
		assert.Equal(t, "workflow", workflowDisplayName)
	})

	t.Run("Get workflow display name", func(t *testing.T) {
		workflowName := "workflow"
		workflowDisplayName := GetWorkflowDisplayName(workflowName)
		assert.Equal(t, "workflow", workflowDisplayName)
	})

	t.Run("Get workflow display name", func(t *testing.T) {
		workflowName := "workflow|JENKINS"
		workflowDisplayName := GetWorkflowDisplayName(workflowName)
		assert.Equal(t, "workflow", workflowDisplayName)
	})
	t.Run("Get workflow display name", func(t *testing.T) {
		workflowName := ""
		workflowDisplayName := GetWorkflowDisplayName(workflowName)
		assert.Equal(t, "", workflowDisplayName)
	})
}
