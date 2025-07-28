package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/calculi-corp/api/go/service"
	pb "github.com/calculi-corp/api/go/vsm/report"
	"github.com/calculi-corp/common/grpc"
	"github.com/calculi-corp/config"
	client "github.com/calculi-corp/grpc-client"
	mock_grpc_client "github.com/calculi-corp/grpc-client/mock"
	testutil "github.com/calculi-corp/grpc-testutil"
	"github.com/calculi-corp/reports-service/constants"
	"github.com/calculi-corp/reports-service/db"
	"github.com/calculi-corp/reports-service/exceptions"
	"github.com/opensearch-project/opensearch-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func init() {
	config.Config.DefineStringFlag("opensearch.endpoint", "", "The Open search Endpoint")
	config.Config.DefineStringFlag("opensearch.user", "", "The Open search Username")
	config.Config.DefineStringFlag("opensearch.pwd", "", "The Open search password")
	config.Config.DefineStringFlag("report.definition.filepath", "../resources/", "Report Definiftion filepath")
	config.Config.Set("logging.level", "INFO")
	testutil.SetUnitTestConfig()
}

const (
	ErrGetDeploymentTypeAssertion  string = "error in type assertion in getDeployments"
	ErrGetCommitTrendTypeAssertion string = "error in type assertion in getCommitTrends"
	ErrEndpointAPIFailure          string = "endpoint api failed"
)

func TestConvertTimeFormat(t *testing.T) {
	tests := []struct {
		inputTime  string
		timeFormat string
		expected   string
		expectErr  bool
	}{
		{"2024/10/03 05:37:49", "12h", "2024/10/03 05:37:49 AM", false},
		{"2024/10/03 23:07:49", "12h", "2024/10/03 11:07:49 PM", false},
		{"2024/10/03 05:37:49", "24h", "2024/10/03 05:37:49", false},
		{"2024/10/03 23:07:49", "24h", "2024/10/03 23:07:49", false},
		{"invalid time", "24h", "", true},
		{"2024/10/03 05:37:49", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.inputTime, func(t *testing.T) {
			result, err := convertTimeFormat(tt.inputTime, tt.timeFormat)

			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestGetExceptionByCode(t *testing.T) {
	var ErrErrorUnsupported string
	testCases := []struct {
		key           string
		expectedError error
	}{
		{ErrGetDeploymentTypeAssertion, errors.New(ErrGetDeploymentTypeAssertion)},
		{ErrGetCommitTrendTypeAssertion, errors.New(ErrGetCommitTrendTypeAssertion)},
		{ErrEndpointAPIFailure, errors.New(ErrEndpointAPIFailure)},
		{ErrErrorUnsupported, errors.ErrUnsupported},
	}

	for _, tc := range testCases {
		exceptions.GetExceptionByCode(tc.key)
		exceptions.GetExceptionByCode(tc.key)
	}
}

func TestComponentReport(t *testing.T) {
	// Create a mock GrpcClient for testing.
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	// Case 1: Test when everything goes smoothly.
	t.Run("Case 1: All components are active", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":      "00000000-0000-0000-0000-000000000000",
			"subOrgId":   "00000000-0000-0000-0000-000000000000",
			"component":  []string{"All"},
			"startDate":  "2023-06-01 10:00:00",
			"endDate":    "2023-06-30 10:00:00",
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background() // Use a context for testing.

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"component_activity": {
					"value": {
					  "238ffe68-8cb4-459d-64ac-2e4f752fe8dc": {
						"repo_url": "https://github.com/calculi-corp/dsl-engine-cli.git",
						"last_active_time": "2023-10-18T13:09:06.000Z",
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

		reports, err := ComponentReport(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 3, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
		testReplacements["component"] = []string{"8f8a7d06-26b0-473a-8e0a-c69557c79a53", "238ffe68-8cb4-459d-64ac-2e4f752fe8dc", "138ffe68-8cb4-459d-64ac-2e4f752fe8dc"}
		reports, _ = ComponentReport(testReplacements, ctx, mockGrpcClient)
		assert.Equal(t, len(reports.Values), 3, "Validating response count for automation run drilldown")
	})

	// Case 2: Test when no components are provided.
	t.Run("Case 2: No components provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "00000000-0000-0000-0000-000000000000",
			"subOrgId":  "00000000-0000-0000-0000-000000000000",
			"component": []string{},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return nil, nil
		}

		_, err := ComponentReport(testReplacements, ctx, mockGrpcClient)

		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 3: No orgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"subOrgId":  "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return nil, nil
		}

		_, err := ComponentReport(testReplacements, ctx, mockGrpcClient)

		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 4: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.
		_, err := ComponentReport(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")

	})

	// Case 6: Error when gRPC call fails
	t.Run("Case 5: gRPC call fails", func(t *testing.T) {
		// Define test replacements
		testReplacements := map[string]any{
			"orgId":     "00000000-0000-0000-0000-000000000000",
			"subOrgId":  "00000000-0000-0000-0000-000000000000",
			"component": []string{"Ingestion Service"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		// Mock gRPC call to SendGrpcCtx, return an error
		expectedError := errors.New("gRPC call failed")
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return nil, expectedError
		}

		_, err := ComponentReport(testReplacements, ctx, mockGrpcClient)

		// Assertions
		if err == nil {
			t.Error("Expected an error, but got no error")
		}
	})

	// Case 7: Test when components are inactive
	t.Run("Case 6: All components are inactive", func(t *testing.T) {
		// Define test replacements
		testReplacements := map[string]any{
			"orgId":     "00000000-0000-0000-0000-000000000000",
			"subOrgId":  "00000000-0000-0000-0000-000000000000",
			"component": []string{"Ingestion Service"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		// Mock gRPC call to SendGrpcCtx, return inactive components
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return nil, nil
		}

		// Call the ComponentReport function
		_, err := ComponentReport(testReplacements, ctx, mockGrpcClient)

		// Assertions
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
		// Add assertions to check if the response contains inactive components
	})

}

func TestAutomationReport(t *testing.T) {
	// Create a mock GrpcClient for testing.
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	// Case 1: Test when everything goes smoothly.
	t.Run("Case 1: Successful automation activities retrieval", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "00000000-0000-0000-0000-000000000000",
			"subOrgId":  "00000000-0000-0000-0000-000000000000",
			"component": []string{"Ingestion Service"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
			"timeZone":  "Asia/Calcutta",
		}
		ctx := context.Background()

		responseString := `{
			"aggregations": {
			  "automation_activity": {
				"value": {
				  "2bb25065-e310-4552-91b8-6d4fdf6ac429": {
					"automation_id": "2bb25065-e310-4552-91b8-6d4fdf6ac429",
					"last_active_time": "2023-10-16T06:37:33.000Z",
					"component_id": "65240b81-48d8-44a7-4592-82a9e7a1a132",
					"workflow_name": "workflow",
					"branch_id": "a532adca-3caf-42a4-8478-8d2a38530378",
					"component_name": "common",
					"branch_name": "test"
				  },
				  "99a8e78a-5a3f-4675-85f0-22eab67c50a0": {
					"automation_id": "99a8e78a-5a3f-4675-85f0-22eab67c50a0",
					"last_active_time": "2023-10-12T20:42:42.000Z",
					"component_id": "7cb30415-b101-457d-75f7-9b7651f2632f",
					"workflow_name": "workflow",
					"branch_id": "5f5458f2-920d-4176-8d66-dd8b20c77bd0",
					"component_name": "stats-service",
					"branch_name": "main"
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
					Id:            "7cb30415-b101-457d-75f7-9b7651f2632f",
					Name:          "stats-service",
					RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
				},
				{
					Id:            "65240b81-48d8-44a7-4592-82a9e7a1a132",
					Name:          "common",
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
		_, err := AutomationReport(testReplacements, ctx, mockGrpcClient)

		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		// Add specific assertions to check the structure and content of the 'reports' variable based on your expected results.
	})

	// Case 2: Test when no automation activities are available.
	t.Run("Case 2: No automation activities available", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "00000000-0000-0000-0000-000000000000",
			"subOrgId":  "00000000-0000-0000-0000-000000000000",
			"component": []string{},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background()

		_, err := AutomationReport(testReplacements, ctx, mockGrpcClient)

		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		// Add specific assertions to check that the 'reports' variable is empty or has the expected structure.
		// For example, you can check that the length of reports.Values is 0 to ensure no automation activities were reported.
	})

	// Case 3: Test when the date range is outside the available data.
	t.Run("Case 3: Date range outside available data", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"comp1", "comp2"},
			"startDate": "2024-01-01", // Assuming no data is available beyond this date.
			"endDate":   "2024-12-31",
		}
		ctx := context.Background()

		_, err := AutomationReport(testReplacements, ctx, mockGrpcClient)

		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		// Add specific assertions to check that the 'reports' variable is empty or has the expected structure.
		// For example, you can check that the length of reports.Values is 0 to ensure no automation activities were reported.
	})

	// Case 4: Test when no organization ID is provided.
	t.Run("Case 4: No organization ID provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "", // No organization ID provided.
			"subOrgId":  "suborg123",
			"component": []string{"comp1", "comp2"},
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
		}
		ctx := context.Background()

		_, err := AutomationReport(testReplacements, ctx, mockGrpcClient)
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		// Add specific assertions to check that the 'err' variable contains the expected error related to the missing organization ID.
		// For example, you can use the `assert` library to check if the error message is as expected.
	})

	t.Run("Case 5: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		//mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		_, err := AutomationReport(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")

	})

}

// func TestGetAutomationResponseMap(t *testing.T) {
// 	// Create a new controller for the mock objects
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	// Create a mock GrpcClient for testing
// 	mockGrpcClient := client.NewMockGrpcClient(mockCtrl)

// 	// Define your test input data
// 	orgID := "org123"
// 	components := []string{"comp1", "comp2"}
// 	automationSet := map[string]struct{}{"automation1": {}, "automation2": {}}
// 	ctx := context.Background()

// 	// Create a mock cache (assuming you have a cache package)
// 	mockCoreDataCache := cache.NewMockCoreDataCache(mockCtrl)
// 	cache.SetCoreDataCache(mockCoreDataCache)

// 	// Set up mock responses for helper.GetOrganisationServices
// 	serviceResponse := &api.ServiceResponse{
// 		Service: []*api.Service{
// 			{Id: "comp1", Name: "Component1"},
// 			{Id: "comp2", Name: "Component2"},
// 			{Id: "comp3", Name: "Component3"},
// 		},
// 	}

// 	mockGrpcClient.EXPECT().SendGrpcCtx(
// 		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
// 	).Return(serviceResponse, nil)

// 	// Set up mock responses for coreDataCache methods as needed
// 	// Example: mockCoreDataCache.EXPECT().GetChildrenOfType(gomock.Any(), gomock.Any()).Return(...)
// 	// Example: mockCoreDataCache.EXPECT().Get(gomock.Any()).Return(...)

// 	// Call the function you want to test
// 	automationMap := yourpackage.GetAutomationResponseMap(ctx, mockGrpcClient, orgID, components, automationSet)

// 	// Add assertions to verify the results
// 	if len(automationMap) != 2 {
// 		t.Errorf("Expected 2 automations in the map, but got %d", len(automationMap))
// 	}

// 	// You can add more specific assertions based on the expected results.
// }

func TestAutomationRunReport(t *testing.T) {
	// Create a mock GrpcClient for testing.
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: All inputs are present", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"component":  []string{"comp1", "comp2"},
			"startDate":  "2023-01-01",
			"endDate":    "2023-12-31",
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
			  "automation_run_activity": {
				"value": {
				  "10000000-0000-0000-0000-000000000000": [
					{
					  "automation_id": "10000000-0000-0000-0000-000000000000",
					  "duration": 0,
					  "component_id": "6cbb8b60-16d6-48a4-61fe-3a33c6c1169f",
					  "run_id": "10000000-1000-1000-1000-100000000000",
					  "component_name": "testComp01",
					  "run_number": 212,
					  "status_timestamp": "2023-09-29T11:10:22.000Z",
					  "status": "Success"
					},
					{
					  "automation_id": "10000000-0000-0000-0000-000000000000",
					  "duration": 0,
					  "component_id": "6cbb8b60-16d6-48a4-61fe-3a33c6c1169f",
					  "run_id": "10000000-1000-1000-1000-200000000000",
					  "component_name": "testComp01",
					  "run_number": 211,
					  "status_timestamp": "2023-09-29T11:09:18.000Z",
					  "status": "Success"
					}
				  ],
				  "20000000-0000-0000-0000-000000000000": [
					{
					  "automation_id": "20000000-0000-0000-0000-000000000000",
					  "duration": 5000,
					  "component_id": "d723fed3-4c2b-4a0b-9520-75e000dbe911",
					  "run_id": "25b9deb7-d874-4f06-b448-ad43a999b62b",
					  "component_name": "testComp02",
					  "run_number": 119,
					  "status_timestamp": "2023-09-28T08:36:52.000Z",
					  "status": "Failure"
					}
				  ],
				  "30000000-0000-0000-0000-000000000000": [
					{
					  "automation_id": "30000000-0000-0000-0000-000000000000",
					  "duration": 0,
					  "component_id": "ae694ee5-6234-419c-5035-f3d59622565d",
					  "run_id": "57b74931-7173-402a-bd1b-bf48445bc37d",
					  "component_name": "testComp02",
					  "run_number": 199,
					  "status_timestamp": "2023-09-21T03:26:36.000Z",
					  "status": "Success"
					}
				  ]
				}
			  }
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		reports, err := AutomationRunReport(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 4, "Validating response count for automation run drilldown")

		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

	t.Run("Case 2: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":      "00000000-0000-0000-0000-000000000000",
			"component":  []string{"All"},
			"startDate":  "2023-06-01 10:00:00",
			"endDate":    "2023-06-01 10:00:00",
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := AutomationRunReport(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")

	})

}

// func TestPullRequestsReport(t *testing.T) {
// 	// Create a mock GrpcClient for testing.
// 	mockCtrl := gomock.NewController(t)
// 	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

// 	grpc.SetSharedGrpcClient(mockGrpcClient)
// 	defer grpc.SetSharedGrpcClient(nil)

// 	t.Run("Case 1: All inputs are present", func(t *testing.T) {
// 		// Define your test input data.
// 		testReplacements := map[string]any{
// 			"orgId":     "org123",
// 			"subOrgId":  "suborg123",
// 			"component": []string{"comp1", "comp2"},
// 			"startDate": "2023-01-01",
// 			"endDate":   "2023-12-31",
// 		}
// 		ctx := context.Background() // Use a context for testing.

// 		openSearchClient = func() (*opensearch.Client, error) {
// 			return opensearch.NewDefaultClient()
// 		}
// 		responseString := `{
// 			"aggregations": {
// 				"pullrequests": {
// 					"value": [
// 					  {
// 						"component_id": "7dfaabc0-5eb0-40c6-427c-a2ac4cc01da3",
// 						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
// 						"pr_created_time": "2023-10-10T18:43:33.000Z",
// 						"component_name": "config",
// 						"target_branch": "",
// 						"review_status": "OPEN",
// 						"pull_request_id": "26",
// 						"repository_name": "calculi-corp/config",
// 						"timestamp": "2023-10-10T18:45:39.000Z",
// 						"source_branch": ""
// 					  },
// 					  {
// 						"component_id": "785894c9-91cf-47b9-734d-ae73377cc29c",
// 						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
// 						"pr_created_time": "2023-10-11T13:02:31.000Z",
// 						"component_name": "run-service",
// 						"target_branch": "main",
// 						"review_status": "APPROVED",
// 						"pull_request_id": "125",
// 						"repository_name": "calculi-corp/run-service",
// 						"timestamp": "2023-10-12T14:18:45.000Z",
// 						"source_branch": "SDP-9645-anonymous-id"
// 					  },
// 					  {
// 						"component_id": "716d4922-add7-474b-7f21-240a8e253c38",
// 						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
// 						"pr_created_time": "2023-10-12T16:08:38.000Z",
// 						"component_name": "jwt-validator",
// 						"target_branch": "main",
// 						"review_status": "APPROVED",
// 						"pull_request_id": "46",
// 						"repository_name": "calculi-corp/jwt-validator",
// 						"timestamp": "2023-10-12T16:12:28.000Z",
// 						"source_branch": "SDP-9900"
// 					  },
// 					  {
// 						"component_id": "55f26650-b259-4453-6ff3-c2843d096780",
// 						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
// 						"pr_created_time": "2023-10-04T18:24:10.000Z",
// 						"component_name": "api",
// 						"target_branch": "",
// 						"review_status": "OPEN",
// 						"pull_request_id": "451",
// 						"repository_name": "calculi-corp/api",
// 						"timestamp": "2023-10-04T18:24:23.000Z",
// 						"source_branch": ""
// 					  }
// 					]
// 				}
// 			}
// 		  }`
// 		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
// 			return responseString, nil
// 		}

// 		reports, err := PullRequestsReport(testReplacements, ctx, mockGrpcClient)
// 		assert.Equal(t, len(reports.Values), 4, "Validating response count for automation run drilldown")
// 		// Assertions to check if the result is as expected.
// 		if err != nil {
// 			t.Errorf("Expected no error, but got an error: %v", err)
// 		}
// 	})

// }

func TestCommitsReport(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: All inputs are present", func(t *testing.T) {
		// Define your test input data.
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"component":  []string{"comp1", "comp2"},
			"timeFormat": "12h",
			"startDate":  "2023-01-01",
			"endDate":    "2023-12-31",
		}
		ctx := context.Background() // Use a context for testing.

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"commits": {
					"value": [
					  {
						"component_id": "e0fb1caf-4c53-4be5-644e-1836524bab2a",
						"commit_timestamp": "2023-10-13T21:42:09.000Z",
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"author": "24964627+vwatinteg@users.noreply.github.com",
						"component_name": "keycloak",
						"branch": "main",
						"commit_id": "82ca44df2196f55d29bec50a2e497c74cd0b544c",
						"repository_name": "calculi-corp/keycloak",
						"repository_url": "https://github.com/calculi-corp/keycloak"
					  },
					  {
						"component_id": "238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
						"commit_timestamp": "2023-10-18T15:03:10.000Z",
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"author": "c_asas@cloudbees.com",
						"component_name": "dsl-engine-cli",
						"branch": "SDP-7551",
						"commit_id": "2822c63f86c58f9542d5b3b4307f35ad7c242b47",
						"repository_name": "calculi-corp/dsl-engine-cli",
						"repository_url": "https://github.com/calculi-corp/dsl-engine-cli"
					  },
					  {
						"component_id": "9b707c3e-b9a2-4705-6c44-6f9f250f61c6",
						"commit_timestamp": "2023-10-11T15:24:29.000Z",
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"author": "psarkar@cloudbees.com",
						"component_name": "ng-tekton-dispatch-service",
						"branch": "SDP-9284-runattempt",
						"commit_id": "9b1777e0bb65cd642e46266bc80a16389dd64e11",
						"repository_name": "calculi-corp/ng-tekton-dispatch-service",
						"repository_url": "https://github.com/calculi-corp/ng-tekton-dispatch-service"
					  },
					  {
						"component_id": "341a8fc8-5460-4a06-88ea-7015a31b5e31",
						"commit_timestamp": "2023-10-03T09:17:03.000Z",
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"author": "132610544+cloudbees-platform-app-staging[bot]@users.noreply.github.com",
						"component_name": "template-circleci-actions",
						"branch": "workflow_test",
						"commit_id": "be01781669c388743f7e94a53840211922c5cc03",
						"repository_name": "calculi-corp/template-circleci-actions",
						"repository_url": "https://github.com/calculi-corp/template-circleci-actions"
					  }
					]
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		reports, err := CommitsReport(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 4, "Validating response count for automation run drilldown")
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})
}

func TestGetRunAndDeployedEnvMap(t *testing.T) {
	// Test case with all deployments available
	t.Run("Case 1: All deployments available", func(t *testing.T) {
		response := `{
            "aggregations": {
                "deployments": {
                    "value": {
						"staging":["1"]
					}
                }
            }
        }`

		result := make(map[string]interface{})
		runsEnvMap := getRunAndDeployedEnvMap(result, response)

		if len(runsEnvMap) != 1 {
			t.Errorf("Expected 1 automation in commitRunMap, but got %d", len(runsEnvMap))
		}
	})

	// Test case with no deployments available
	t.Run("Case 2: No deployments available", func(t *testing.T) {
		response := `{
            "aggregation": {
                "deployments": {
                    "value": { "c9be43f4-47e0-463a-a7ff-81b060878231_3dda89cf-30ab-40e6-b0ca-aa2e0a60c3c9": {
						"automation_id": "c9be43f4-47e0-463a-a7ff-81b060878231",
						"start_time": 0,
						"component_id": "89e389cc-3f56-429c-6ee1-1a8bfc8ef5ae",
						"commit_sha": "",
						"run_id": "3dda89cf-30ab-40e6-b0ca-aa2e0a60c3c9",
						"commit_description": "",
						"component_name": "grpc-server",
						"run_number": 30,
						"status_timestamp": "2023-10-12T14:47:50.000Z",
						"org_name": "cloudbees-staging",
						"status": "Failure"
					  }}
                }
            }
        }`

		result := make(map[string]interface{})
		runsEnvMap := getRunAndDeployedEnvMap(result, response)

		expected := map[string][]string{} // Expect an empty map

		if !reflect.DeepEqual(runsEnvMap, expected) {
			t.Errorf("Expected %v, but got %v", expected, runsEnvMap)
		}
	})

	// Test case with mixed deployments and no values
	t.Run("Case 3: Mixed deployments and no values", func(t *testing.T) {
		response := `{
            "aggregation": {
                "deployments": {
                    "value": {
                        "deployment1": [],
                        "deployment2": ["env3", "env4"],
                        "deployment3": []
                    }
                }
            }
        }`

		result := make(map[string]interface{})
		runsEnvMap := getRunAndDeployedEnvMap(result, response)

		expected := map[string][]string{
			"deployment1": {},
			"deployment2": {"env3", "env4"},
			"deployment3": {},
		}

		if reflect.DeepEqual(runsEnvMap, expected) {
			t.Errorf("Expected %v, but got %v", expected, runsEnvMap)
		}
	})
}

func TestCPSRunInitiatingCommits(t *testing.T) {
	// Create a mock GrpcClient for testing.
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: All inputs are present", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"component":  []string{"comp1", "comp2"},
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background() // Use a context for testing.

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"automation_run": {
					"value": {
					  "c9be43f4-47e0-463a-a7ff-81b060878231_3dda89cf-30ab-40e6-b0ca-aa2e0a60c3c9": {
						"automation_id": "c9be43f4-47e0-463a-a7ff-81b060878231",
						"start_time": 0,
						"component_id": "89e389cc-3f56-429c-6ee1-1a8bfc8ef5ae",
						"commit_sha": "",
						"run_id": "3dda89cf-30ab-40e6-b0ca-aa2e0a60c3c9",
						"commit_description": "",
						"component_name": "grpc-server",
						"run_number": 30,
						"status_timestamp": "2023-10-12T14:47:50.000Z",
						"org_name": "cloudbees-staging",
						"status": "Failure"
					  },
					  "048f6115-ba32-4575-a297-3ed3ca4464cb_09f5923c-fdcf-4765-a7f6-05c6abc697d4": {
						"automation_id": "048f6115-ba32-4575-a297-3ed3ca4464cb",
						"start_time": 1696519001000,
						"component_id": "785894c9-91cf-47b9-734d-ae73377cc29c",
						"commit_sha": "477a22b0dc2fdd5856365c607fd2b45a8c1cd7d9",
						"run_id": "09f5923c-fdcf-4765-a7f6-05c6abc697d4",
						"commit_description": "Trigger",
						"component_name": "run-service",
						"run_number": 282,
						"status_timestamp": "2023-10-05T15:22:58.000Z",
						"org_name": "cloudbees-staging",
						"status": "Success"
					  },
					  "fd51dbbe-18a0-434a-bd3a-f16752ab3fdb_11267517-2ce6-4e9e-a524-bd9bb257c1c6": {
						"automation_id": "fd51dbbe-18a0-434a-bd3a-f16752ab3fdb",
						"start_time": 1697039826000,
						"component_id": "9b707c3e-b9a2-4705-6c44-6f9f250f61c6",
						"commit_sha": "1d63c05be16b86f29b713e135fa127505ad5d151",
						"run_id": "11267517-2ce6-4e9e-a524-bd9bb257c1c6",
						"commit_description": "refactor: [SDP-9284] ensure that the run attempt assignment is confirmed before creating meta-pipeline",
						"component_name": "ng-tekton-dispatch-service",
						"run_number": 309,
						"status_timestamp": "2023-10-11T16:07:08.000Z",
						"org_name": "cloudbees-staging",
						"status": "Success"
					  }
					}
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		if components, ok := testReplacements["components"].([]string); ok && len(components) > 0 {
			// Components field is valid, use it.
		} else {

			t.Skip("Skipping test because 'components' are either missing or not provided as a valid slice.")
			return
		}

		// Mock the necessary dependencies and method calls.
		mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		// Call the function and capture the result.
		reports, err := CPSRunInitiatingCommits(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 3, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 2: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":      "00000000-0000-0000-0000-000000000000",
			"component":  []string{"All"},
			"startDate":  "2023-06-01 10:00:00",
			"endDate":    "2023-06-01 10:00:00",
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background()

		_, err := CPSRunInitiatingCommits(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")

	})
}

func TestGetCommitRunMap(t *testing.T) {
	result := map[string]interface{}{
		"aggregations": map[string]interface{}{
			"automation_run": map[string]interface{}{
				"value": map[string]interface{}{
					"key01": map[string]interface{}{
						"component_id":     "comp1",
						"component_name":   "Component1",
						"automation_id":    "automation1",
						"status":           "Success",
						"run_number":       1.0,
						"run_id":           "run1",
						"org_name":         "Org1",
						"start_time":       1634567890000.0,
						"status_timestamp": "2023-10-19T12:34:56Z",
					},
					"key02": map[string]interface{}{
						"component_id":       "comp1",
						"component_name":     "Component1",
						"automation_id":      "automation1",
						"status":             "Failed",
						"run_number":         2.0,
						"run_id":             "run2",
						"org_name":           "Org1",
						"start_time":         1634567900000.0,
						"status_timestamp":   "2023-10-20T12:34:56Z",
						"commit_sha":         "abc123",
						"commit_description": "Fix a bug",
					},
				},
			},
		},
	}

	componentSet := make(map[string]struct{})

	// Call the getCommitRunMap function with the input data.
	commitRunMap := getCommitRunMap(result, componentSet)

	// Assertions to check if the result is as expected.
	if len(commitRunMap) != 1 {
		t.Errorf("Expected 1 automation in commitRunMap, but got %d", len(commitRunMap))
	}

	// Validate the content of the commit run response.
	automationRuns, ok := commitRunMap["automation1"]
	if !ok {
		t.Errorf("Expected automation1 in commitRunMap, but it was not found")
	} else {
		if len(automationRuns) != 2 {
			t.Errorf("Expected 2 commit runs in automation1, but got %d", len(automationRuns))
		}
	}

}

func TestCodeProgressionSnapshotBuilds(t *testing.T) {
	// Create a mock GrpcClient for testing.
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	// Case 1: Test when everything goes smoothly.
	t.Run("Case 1: All inputs are present", func(t *testing.T) {
		// Define your test input data.
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"component":  []string{"comp1", "comp2"},
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background() // Use a context for testing.
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"builds": {
					"value": {
					  "e30e43d2-4076-4388-add8-89dc03ede5cf_SUCCEEDED": {
						"step_kind": "build",
						"component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
						"run_id": "e30e43d2-4076-4388-add8-89dc03ede5cf",
						"workflow_name": "workflow",
						"component_name": "data-ingestion-service",
						"step_id": "s009-4-build-go-binary",
						"automation_id": "ca90f12f-75e0-44a0-9cae-9164a387d2fb",
						"duration": 56000,
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"job_id": "build-deploy-data-ingestion-service",
						"run_number": 499,
						"target_env": "staging",
						"status_timestamp": "2023-10-16T17:56:12.000Z",
						"org_name": "cloudbees-staging",
						"status": "Success",
						"source": "GitHub"
					  },
					  "427959f5-1d4f-4afe-884d-fd737652f179_SUCCEEDED": {
						"step_kind": "build",
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"run_id": "427959f5-1d4f-4afe-884d-fd737652f179",
						"workflow_name": "workflow",
						"component_name": "reports-service",
						"step_id": "s010-5-build-go-binary",
						"automation_id": "b8384705-244a-4ed2-a5fc-ca5e6f837b27",
						"duration": 49000,
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"job_id": "build-deploy-reports-service",
						"run_number": 873,
						"target_env": "staging",
						"status_timestamp": "2023-10-18T14:54:41.000Z",
						"org_name": "cloudbees-staging",
						"status": "Success",
						"source": "GitHub"
					  }
					}
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		reports, err := CodeProgressionSnapshotBuilds(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 2: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":      "00000000-0000-0000-0000-000000000000",
			"component":  []string{"All"},
			"startDate":  "2023-06-01 10:00:00",
			"endDate":    "2023-06-01 10:00:00",
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := CodeProgressionSnapshotBuilds(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")

	})

	t.Run("Case 3: Check env name", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":      "00000000-0000-0000-0000-000000000000",
			"subOrgId":   "suborg123",
			"component":  []string{"All"},
			"startDate":  "2023-06-01 10:00:00",
			"endDate":    "2023-06-01 10:00:00",
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}

		responseString1 := `{
	"aggregations": {
		"automation_run": {
			"value": {
			  "c9be43f4-47e0-463a-a7ff-81b060878231_ff4e39f0-0145-40e2-b5e3-289a89dabfde": {
				"automation_id": "c9be43f4-47e0-463a-a7ff-81b060878231",
				"start_time": 0,
				"component_id": "89e389cc-3f56-429c-6ee1-1a8bfc8ef5ae",
				"commit_sha": "",
				"run_id": "ff4e39f0-0145-40e2-b5e3-289a89dabfde",
				"commit_description": "",
				"component_name": "grpc-server",
				"run_number": 30,
				"status_timestamp": "2023-10-12T14:47:50.000Z",
				"org_name": "cloudbees-staging",
				"status": "Failure"
			  },
			  "048f6115-ba32-4575-a297-3ed3ca4464cb_98e01e22-b711-4172-829e-3b3b5d13b3db": {
				"automation_id": "048f6115-ba32-4575-a297-3ed3ca4464cb",
				"start_time": 1696519001000,
				"component_id": "785894c9-91cf-47b9-734d-ae73377cc29c",
				"commit_sha": "477a22b0dc2fdd5856365c607fd2b45a8c1cd7d9",
				"run_id": "98e01e22-b711-4172-829e-3b3b5d13b3db",
				"commit_description": "Trigger",
				"component_name": "run-service",
				"run_number": 282,
				"status_timestamp": "2023-10-05T15:22:58.000Z",
				"org_name": "cloudbees-staging",
				"status": "Success"
			  },
			  "048f6115-ba32-4575-a297-3ed3ca4464cb_9624r516-b711-4172-829e-3b3b5d13b3db": {
				"automation_id": "048f6115-ba32-4575-a297-3ed3ca4464cb",
				"start_time": 1696519001000,
				"component_id": "785894c9-91cf-47b9-734d-ae73377cc29c",
				"commit_sha": "477a22b0dc2fdd5856365c607fd2b45a8c1cd7d9",
				"run_id": "9624r516-b711-4172-829e-3b3b5d13b3db",
				"commit_description": "Trigger",
				"component_name": "run-service",
				"run_number": 282,
				"status_timestamp": "2023-10-05T15:22:58.000Z",
				"org_name": "cloudbees-staging",
				"status": "Success"
			  }
			}
		}
	}
  }`

		responseString2 := `{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 3,
    "successful": 3,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 1,
      "relation": "eq"
    },
    "max_score": null,
    "hits": []
  },
  "aggregations": {
    "deployments": {
      "value": {
        "ff4e39f0-0145-40e2-b5e3-289a89dabfde": [
          "preprod-us-west-2","preprod-us-east-1"
        ],
		"98e01e22-b711-4172-829e-3b3b5d13b3db": [
  		  "preprod-us-east-1"
		]
      }
    }
  }
}`

		ctx := context.Background()

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			if IndexName == "automation_run_status" {
				return responseString1, nil
			} else if IndexName == "deploy_data" {
				return responseString2, nil
			} else {
				return "", errors.New("Unknown IndexName")
			}
		}

		val, err := CPSRunInitiatingCommits(testReplacements, ctx, mockGrpcClient)

		listVal := val.AsSlice()
		if len(listVal) != 2 {
			t.Errorf("Expected 2, but got %d", len(listVal))
		}

		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})
}

func TestCodeProgressionSnapshotDeployments(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: All inputs are present", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"component":  []string{"comp1", "comp2"},
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background()

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"deployments": {
					"value": {
					  "staging_3f32c118-b3ff-4edd-920b-7a331815fc1d_build-deploy-reports-service_s022-11-helmpush-0-push_SUCCEEDED": {
						"step_kind": "deploy",
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"run_id": "3f32c118-b3ff-4edd-920b-7a331815fc1d",
						"workflow_name": "workflow",
						"component_name": "reports-service",
						"step_id": "s022-11-helmpush-0-push",
						"automation_id": "84c30481-356e-460b-807b-31ee8e14ff09",
						"duration": 5000,
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"job_id": "build-deploy-reports-service",
						"run_number": 824,
						"target_env": "staging",
						"status_timestamp": "2023-10-16T07:58:16.000Z",
						"org_name": "cloudbees-staging",
						"status": "SUCCEEDED"
					  },
					  "staging_d4726bca-7881-446a-84ee-ed7af68ba5ef_build-deploy-data-ingestion-service_s021-10-helmpush-0-push_SUCCEEDED": {
						"step_kind": "deploy",
						"component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
						"run_id": "d4726bca-7881-446a-84ee-ed7af68ba5ef",
						"workflow_name": "workflow",
						"component_name": "data-ingestion-service",
						"step_id": "s021-10-helmpush-0-push",
						"automation_id": "1d64bb32-4aac-4c9b-bfed-0d0574b6c1f1",
						"duration": 0,
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"job_id": "build-deploy-data-ingestion-service",
						"run_number": 453,
						"target_env": "staging",
						"status_timestamp": "2023-10-15T08:38:03.000Z",
						"org_name": "cloudbees-staging",
						"status": "SUCCEEDED"
					  }
					}
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		reports, err := CodeProgressionSnapshotDeployments(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 2: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":      "00000000-0000-0000-0000-000000000000",
			"component":  []string{"All"},
			"startDate":  "2023-06-01 10:00:00",
			"endDate":    "2023-06-01 10:00:00",
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := CodeProgressionSnapshotDeployments(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")

	})

	// Add more test cases as needed.
}

func TestSuccessfulBuildDuration(t *testing.T) {
	// Create a mock GrpcClient for testing.
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: All inputs are present", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"component":  []string{"comp1", "comp2"},
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background()

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"builds": {
					"value": {
					  "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_daf5a489-c83d-4960-b05f-9836dc819db7_build-deploy-reports-service_s010-5-build-go-binary_2023-10-11T16:35:43.000Z": {
						"step_kind": "build",
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"run_id": "daf5a489-c83d-4960-b05f-9836dc819db7",
						"workflow_name": "workflow",
						"component_name": "reports-service",
						"step_id": "s010-5-build-go-binary",
						"automation_id": "a73f77af-a0ce-4f54-a1b6-ac416be813f4",
						"duration": 29000,
						"start_time": 1697042114000,
						"completed_time": 1697042143000,
						"job_id": "build-deploy-reports-service",
						"run_number": 766,
						"target_env": "staging",
						"status_timestamp": "2023-10-11T16:35:43.000Z",
						"status": "SUCCEEDED",
						"source": "GitHub"
					  },
					  "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_c388f0d0-343f-49ac-9581-a77276543f34_build-deploy-reports-service_s010-5-build-go-binary_2023-10-12T11:46:57.000Z": {
						"step_kind": "build",
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"run_id": "c388f0d0-343f-49ac-9581-a77276543f34",
						"workflow_name": "workflow",
						"component_name": "reports-service",
						"step_id": "s010-5-build-go-binary",
						"automation_id": "a73f77af-a0ce-4f54-a1b6-ac416be813f4",
						"duration": 28000,
						"start_time": 1697111189000,
						"completed_time": 1697111217000,
						"job_id": "build-deploy-reports-service",
						"run_number": 783,
						"target_env": "staging",
						"status_timestamp": "2023-10-12T11:46:57.000Z",
						"status": "SUCCEEDED",
						"source": "GitHub"
					  }
					}
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		reports, err := SuccessfulBuildDuration(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 2: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":      "00000000-0000-0000-0000-000000000000",
			"component":  []string{"All"},
			"startDate":  "2023-06-01 10:00:00",
			"endDate":    "2023-06-01 10:00:00",
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := SuccessfulBuildDuration(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")

	})

}

func TestDeploymentOverviewDrilldown(t *testing.T) {
	// Create a mock GrpcClient for testing.
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: All inputs are present", func(t *testing.T) {
		// Define your test input data.
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"component":  []string{"comp1", "comp2"},
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background()

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"deployments": {
					"value": {
					  "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_6de1a874-f52b-4aa5-836f-0aec04b92deb_build-deploy-reports-service_s022-11-helmpush-0-push_staging_SUCCEEDED": {
						"step_kind": "deploy",
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"run_id": "6de1a874-f52b-4aa5-836f-0aec04b92deb",
						"workflow_name": "workflow",
						"component_name": "reports-service",
						"step_id": "s022-11-helmpush-0-push",
						"automation_id": "ed82584d-0758-44b1-9102-a2e9202ccc70",
						"duration": 0,
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"job_id": "build-deploy-reports-service",
						"run_number": 779,
						"target_env": "staging",
						"status_timestamp": "2023-10-12T10:50:25.000Z",
						"org_name": "cloudbees-staging",
						"status": "Success"
					  },
					  "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_24cd899a-6479-47b9-8663-8e4982677d32_build-deploy-reports-service_s020-10-helmpush-0-push_staging_SUCCEEDED": {
						"step_kind": "deploy",
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"run_id": "24cd899a-6479-47b9-8663-8e4982677d32",
						"workflow_name": "workflow",
						"component_name": "reports-service",
						"step_id": "s020-10-helmpush-0-push",
						"automation_id": "5827d8ba-b3e7-4644-a77b-5585ef0d03d7",
						"duration": 0,
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"job_id": "build-deploy-reports-service",
						"run_number": 572,
						"target_env": "staging",
						"status_timestamp": "2023-09-27T04:52:02.000Z",
						"org_name": "cloudbees-staging",
						"status": "Success"
					  }
					}
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		reports, err := DeploymentOverviewDrilldown(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})
	t.Run("Case 2: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":      "00000000-0000-0000-0000-000000000000",
			"component":  []string{"All"},
			"startDate":  "2023-06-01 10:00:00",
			"endDate":    "2023-06-01 10:00:00",
			"timeFormat": "12h",
			"timeZone":   "Asia/Calcutta",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := DeploymentOverviewDrilldown(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")

	})

}

func TestDoraMetricsMttr(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Successful execution", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
			"component": []string{"comp1", "comp2"},
			"timeZone":  "Asia/Calcutta",
		}
		ctx := context.Background()

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"deployments": {
					"value": [
					  {
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"recovered_duration": 243000,
						"failed_run_number": 529,
						"recovered_on": 1695388595000,
						"component_name": "reports-service",
						"recovered_run": "15330b40-50bc-47e5-a4d6-bd0a97f890df",
						"failed_run": "acc935ab-7eec-48fe-b6c4-0ab6d831c18b",
						"failed_on": 1695388352000,
						"recovered_run_number": 530
					  },
					  {
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"recovered_duration": 809000,
						"failed_run_number": 547,
						"recovered_on": 1695706983000,
						"component_name": "reports-service",
						"recovered_run": "5f6020c4-52f6-4ecb-9429-60675537052d",
						"failed_run": "d19543a9-6e75-41c5-8202-b68b0bc1f028",
						"failed_on": 1695706174000,
						"recovered_run_number": 549
					  }
					]
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		reports, err := DoraMetricsMttr(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

	t.Run("Case 5: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
			"timeZone":  "Asia/Calcutta",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := DoraMetricsMttr(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")

	})

}

func TestFailureRate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: Successful execution ", func(t *testing.T) {
		// Define your test input data.
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
			"component": []string{"comp1", "comp2"},
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"deployments": {
					"value": {
					  "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a": {
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"deployments": 316,
						"success": 300,
						"failure": 16,
						"component_name": "reports-service"
					  },
					  "dda69191-5492-4b7e-88b2-9d9d42f61899": {
						"component_id": "dda69191-5492-4b7e-88b2-9d9d42f61899",
						"deployments": 97,
						"success": 96,
						"failure": 1,
						"component_name": "analytics"
					  }

					}
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		ctx := context.Background()
		reports, err := FailureRate(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 5: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background()

		_, err := FailureRate(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "No_suborgId_provided")

	})

}

func TestDeploymentFrequencyAndLeadTime(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Successful execution with multiple components", func(t *testing.T) {
		// Define your test input data.
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"startDate":  "2023-01-01",
			"timeFormat": "12h",
			"endDate":    "2023-12-31",
			"component":  []string{"comp1", "comp2"},
			// Add components if needed
		}
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"deployments": {
					"value": [
					  {
						"step_kind": "deploy",
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"run_id": "6de1a874-f52b-4aa5-836f-0aec04b92deb",
						"workflow_name": "workflow",
						"component_name": "reports-service",
						"run_start_time": 1697107565000,
						"step_id": "s022-11-helmpush-0-push",
						"automation_id": "ed82584d-0758-44b1-9102-a2e9202ccc70",
						"job_id": "build-deploy-reports-service",
						"run_number": 779,
						"target_env": "staging",
						"status_timestamp": "2023-10-12T10:50:25.000Z",
						"status_timestamp_zoned": "2023-10-12T10:50:25.000Z",
						"run_start_time_string_zoned": "2023-10-12 10:50:25",
						"status": "SUCCEEDED"
					  },
					  {
						"step_kind": "deploy",
						"component_id": "54d882eb-06ad-4528-9b7f-7a5db6b7bd7a",
						"run_id": "24cd899a-6479-47b9-8663-8e4982677d32",
						"workflow_name": "workflow",
						"component_name": "reports-service",
						"run_start_time": 1695790002000,
						"step_id": "s020-10-helmpush-0-push",
						"automation_id": "5827d8ba-b3e7-4644-a77b-5585ef0d03d7",
						"job_id": "build-deploy-reports-service",
						"run_number": 572,
						"target_env": "staging",
						"status_timestamp": "2023-09-27T04:52:02.000Z",
						"status": "SUCCEEDED",
						"status_timestamp_zoned": "2023-10-12T10:50:25.000Z",
						"run_start_time_string_zoned": "2023-10-12 10:50:25"
					  }
					]
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		ctx := context.Background()

		reports, err := DeploymentFrequencyAndLeadTime(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})
	t.Run("Case 2: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":      "00000000-0000-0000-0000-000000000000",
			"component":  []string{"All"},
			"timeFormat": "12h",
			"startDate":  "2023-06-01 10:00:00",
			"endDate":    "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := DeploymentFrequencyAndLeadTime(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")

	})

}

func TestSecurityScanTypeWorkflowsDrillDown(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"component": []string{"comp1", "comp2"},
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Define the response for the first call to getSearchResponse
		responseString := `{
			"aggregations": {
				"automation_run_activity": {
					"value": {
					  "9d68bb5b-65fc-4389-b706-d61953a5acb6": [
						{
							"automation_id": "9d68bb5b-65fc-4389-b706-d61953a5acb6",
							"duration": 369000,
							"component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
							"run_id": "c4c2edb0-9319-4fc1-bed9-e214f7f6a750",
							"component_name": "data-ingestion-service",
							"run_number": 306,
							"status_timestamp": "2023-09-21T21:06:08.000Z",
							"status": "Success"
						  }
					  ],
					  "af19399f-7e1d-4cc7-a425-4dc5283993cc": [
						{
							"automation_id": "af19399f-7e1d-4cc7-a425-4dc5283993cc",
							"duration": 0,
							"component_id": "2a4dfe84-0fe4-4426-b796-739073b93b2e",
							"run_id": "1c42e13b-c645-4d7f-8af8-6b80da360b44",
							"component_name": "template_snyk_go",
							"run_number": 33,
							"status_timestamp": "2023-09-25T08:57:42.000Z",
							"status": "Success"
						  }
						]
					}
				}
			}
		}`

		// Define the response for the second call to getSearchResponse
		responseString1 := `{
			"aggregations": {
				"distinct_run": {
					"value": {
						"1c42e13b-c645-4d7f-8af8-6b80da360b44": {
							"scanner_names": [
								"sonarqube",
								"snyksca",
								"checkmarx"
							],
							"scanner_types": [
								"SCA",
								"SAST"
							]
						},
						"c4c2edb0-9319-4fc1-bed9-e214f7f6a750": {
							"scanner_names": [
								"snyksast"
							],
							"scanner_types": [
								"SAST"
							]
						}
					}
				}
			}
		}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			if IndexName == "automation_run_status" {
				return responseString, nil
			} else if IndexName == "scan_results" {
				return responseString1, nil
			} else {
				return "", errors.New("Unknown IndexName")
			}
		}

		reports, err := SecurityScanTypeWorkflowsDrillDown(testReplacements, context.Background(), mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 2: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := SecurityScanTypeWorkflowsDrillDown(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")
	})

}

func TestSecurityAutomationRunDrillDown(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"comp1", "comp2"},
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Define the response for the first call to getSearchResponse
		responseString := `{
			"aggregations": {
				"automation_run_activity": {
					"value": {
					  "9d68bb5b-65fc-4389-b706-d61953a5acb6": [
						{
							"automation_id": "9d68bb5b-65fc-4389-b706-d61953a5acb6",
							"duration": 369000,
							"component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
							"run_id": "c4c2edb0-9319-4fc1-bed9-e214f7f6a750",
							"component_name": "data-ingestion-service",
							"run_number": 306,
							"status_timestamp": "2023-09-21T21:06:08.000Z",
							"status": "Success"
						  }
					  ],
					  "af19399f-7e1d-4cc7-a425-4dc5283993cc": [
						{
							"automation_id": "af19399f-7e1d-4cc7-a425-4dc5283993cc",
							"duration": 0,
							"component_id": "2a4dfe84-0fe4-4426-b796-739073b93b2e",
							"run_id": "1c42e13b-c645-4d7f-8af8-6b80da360b44",
							"component_name": "template_snyk_go",
							"run_number": 33,
							"status_timestamp": "2023-09-25T08:57:42.000Z",
							"status": "Success"
						  }
						]
					}
				}
			}
		}`

		// Define the response for the second call to getSearchResponse
		responseString2 := `{
			"aggregations": {
				"distinct_run": {
					"value": {
						"1c42e13b-c645-4d7f-8af8-6b80da360b44": [{
							"stackhawk": "Success",
							"snyksast": "Failed"

						}],
						"c4c2edb0-9319-4fc1-bed9-e214f7f6a750": [{
							"stackhawk": "Success",
							"snyksast": "Failed",
							"trufflehogsast": "Success"
						}]
					}
				}
			}
		}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			if IndexName == "automation_run_status" {
				return responseString, nil
			} else if IndexName == "scan_results" {
				return responseString2, nil
			} else {
				return "", errors.New("Unknown IndexName")
			}
		}
		reports, err := SecurityAutomationRunDrillDown(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 5: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := SecurityAutomationRunDrillDown(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")
	})

}

func TestSecurityAutomationDrillDown(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		// Define test input data.
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"comp1", "comp2"},
		}

		ctx := context.Background()

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"distinct_automation": {
					"value": {
					  "036d13c2-3009-4dae-8f09-2d86cc066ee8": {
						"Scanner_List": [
						  "checkmarx"
						],
						"run_ids": [
						  "c9e60725-2ed7-4d7a-9c1a-aa5474ec93f7"
						],
						"run_count": 1
					  },
					  "b772fc94-3613-484e-9236-9f53f28b7503": {
						"Scanner_List": [
						  "checkmarx"
						],
						"run_ids": [
						  "f29b04a3-61e2-401e-ad7d-9145992fe8d3",
						  "2309e585-a05b-44d6-9947-7f6efa9c51b8"
						],
						"run_count": 2
					  }
					}
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		reports, err := SecurityAutomationDrillDown(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 0, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 5: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := SecurityAutomationDrillDown(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")
	})

}

func TestTestAutomationDrillDown(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Successful execution for test insights workflow runs", func(t *testing.T) {
		testReplacements := map[string]interface{}{
			"org_id":    "8509888e-d27f-44fa-46a9-29bc76f5e790",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-30 10:00:00",
		}

		response := map[string]json.RawMessage{}
		response["cb_test_suites"] = []byte(`{"aggregations": {"component_activity": {"value": {"3333344c-ae0f-4df4-b1a7-efcaacdf449e_3615d762-6069-4fdd-b4f4-549a41a3aeeb_95bf8556-bcb3-4e49-894d-d24d0d66bcdb": {"automation_id": "3615d762-6069-4fdd-b4f4-549a41a3aeeb","test_suites_set": ["","com.manning.spock.chapter1.MultiplierTest"],"component_id": "3333344c-ae0f-4df4-b1a7-efcaacdf449e","branch_id": "95bf8556-bcb3-4e49-894d-d24d0d66bcdb","org_id": "8509888e-d27f-44fa-46a9-29bc76f5e790","test_suite_name": "com.manning.spock.chapter1.MultiplierTest"}}}}}`)
		response["automation_run_status"] = []byte(`{"aggregations": {"component_activity": {"value": {"3333344c-ae0f-4df4-b1a7-efcaacdf449e_3615d762-6069-4fdd-b4f4-549a41a3aeeb": {"automation_id": "3615d762-6069-4fdd-b4f4-549a41a3aeeb","component_id": "3333344c-ae0f-4df4-b1a7-efcaacdf449e","run_id": "3ae13ee9-3339-419e-ab1a-a4fdee3e72c8","org_id": "8509888e-d27f-44fa-46a9-29bc76f5e790","run_ids": ["3ae13ee9-3339-419e-ab1a-a4fdee3e72c8"]}}}}}`)

		multiSearchResponse = func(queries map[string]db.DbQuery) (map[string]json.RawMessage, error) {
			return response, nil
		}

		ctx := context.Background()

		reports, err := TestAutomationDrilldown(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 0, "Validating response count for test insights workflows run drilldown")

		assert.NoError(t, err, "no error")

	})
}

func TestTestComponentDrilldown(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: Successful execution for With test suites", func(t *testing.T) {
		ctx := context.Background() // Use a context for testing.
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		response := map[string]json.RawMessage{}

		response["cb_test_suites"] = []byte(`{"aggregations": {"component_activity": {"value": {"bb7f71a0-2325-4905-ba0a-805461a50006": {"automation_id": "22ed46c7-09c5-4430-beec-382dbb122d85","component_id": "bb7f71a0-2325-4905-ba0a-805461a50006","run_id": "7fdd0287-8002-495f-ad58-d33730ea4903","org_id": "8509888e-d27f-44fa-46a9-29bc76f5e790","component_name": "template_jenkins_cbci_actions","test_suites_set": ["","com.manning.spock.chapter1.MultiplierSpec"],"test_suite_name": "com.manning.spock.chapter1.MultiplierSpec"},"484d5e12-6424-4070-a159-4e5639a807a2": {"automation_id": "23a2cb0b-b98a-46a6-b18f-d9632f44e5d1","component_id": "484d5e12-6424-4070-a159-4e5639a807a2","run_id": "9bc2d9f1-0b04-45b5-993e-0947da234dda","org_id": "8509888e-d27f-44fa-46a9-29bc76f5e790","component_name": "template-publish-test-results","test_suites_set": ["com.manning.spock.chapter1.AdderSpec","github.com/calculi-corp/go-test-result-conversion/internal/actions/test_data/test_suite_2"],"test_suite_name": "github.com/calculi-corp/go-test-result-conversion/internal/actions/test_data/test_suite_2"}}}}}`)

		multiSearchResponse = func(map[string]db.DbQuery) (map[string]json.RawMessage, error) {
			return response, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{
				{
					Id:            "bb7f71a0-2325-4905-ba0a-805461a50006",
					Name:          "template_jenkins_cbci_actions",
					RepositoryUrl: "https://github.com/calculi-corp/template_jenkins_cbci_actions.git",
				},
				{
					Id:            "484d5e12-6424-4070-a159-4e5639a807a2",
					Name:          "template-publish-test-results",
					RepositoryUrl: "https://github.com/calculi-corp/template-publish-test-results",
				},
			},
		}

		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := TestComponentDrillDown(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for test insights component run drilldown")
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		testReplacements["component"] = []string{"bb7f71a0-2325-4905-ba0a-805461a50006", "484d5e12-6424-4070-a159-4e5639a807a2"}
		reports, _ = TestComponentDrillDown(testReplacements, ctx, mockGrpcClient)
		assert.Equal(t, len(reports.Values), 2, "Validating response count for test insights component run drilldown")

		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return nil, errors.New("server error")
		}
		_, err = TestComponentDrillDown(testReplacements, ctx, mockGrpcClient)
		assert.Equal(t, err.Error(), "server error", "Validating response count for test insights component run drilldown")

	})

}

func TestSecurityComponentDrillDown(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: Successful execution with components - all scenario for Buldled Sonarqube", func(t *testing.T) {
		ctx := context.Background() // Use a context for testing.
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		// Test Data covers where the component Id and scanners are collected from scan_results index
		//  Case 1: Component Id found only in Raw_scan =>  Component Id and Scanner(sonarqube) added to response
		//  Case 2: Component Id found in both Raw_scan and scan result with Sonarqube in Scaaners list => No change to scanner list
		//  Case 3: Component Id found in both Raw_scan and scan result with Sonarqube NOT in Scaaners list => Sonarqube added to Scanner list
		response := map[string]json.RawMessage{}
		response["scan"] = []byte(`{"aggregations":{"distinct_component":{"value":{"80a5fe62-4d9c-4385-b2ad-fb6ad99e9b2d":["snyksca","snyksast","mendsca","checkmarx","mendsast","sonarqube"],"7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a":["checkmarx","snyksast"]}}}}`)
		response["rawScan"] = []byte(`{"took":28989,"timed_out":false,"aggregations":{"distinct_component":{"value":["94a81dd1-3f52-4520-891e-b2440f660945","7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a", "80a5fe62-4d9c-4385-b2ad-fb6ad99e9b2d"]}}}`)

		multiSearchResponse = func(map[string]db.DbQuery) (map[string]json.RawMessage, error) {
			return response, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{
				{
					Id:            "80a5fe62-4d9c-4385-b2ad-fb6ad99e9b2d",
					Name:          "dsl-no-commit-1",
					RepositoryUrl: "https://github.com/sample-gr/dsl-no-commit-1.git",
				},
				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
				{
					Id:            "1cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "test",
					RepositoryUrl: "https://github.com/calculi-corp/test.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := SecurityComponentDrillDown(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 3, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		testReplacements["component"] = []string{"80a5fe62-4d9c-4385-b2ad-fb6ad99e9b2d", "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a", "1cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a"}
		reports, _ = SecurityComponentDrillDown(testReplacements, ctx, mockGrpcClient)
		assert.Equal(t, len(reports.Values), 3, "Validating response count for automation run drilldown")

		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return nil, errors.New("server error")
		}
		_, err = SecurityComponentDrillDown(testReplacements, ctx, mockGrpcClient)
		assert.Equal(t, err.Error(), "server error", "Validating response count for automation run drilldown")

	})
	t.Run("Case 5: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := SecurityComponentDrillDown(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")
	})

}

func Test_processDrilldownQueryAndSpec(t *testing.T) {
	replacement := make(map[string]any)

	type args struct {
		replacements map[string]any
		query        string
		reportId     string
		response     *pb.DrilldownResponse
	}
	tests := []struct {
		name    string
		args    args
		want    *structpb.ListValue
		wantErr bool
	}{
		{
			args: args{
				replacements: replacement,
				query:        "",
				reportId:     "",
				response:     &pb.DrilldownResponse{},
			},
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	responseString := `{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":0.0,"hits":[{"_index":"cb_ci_tool_insight","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_0edebf0d-6797-4ec0-9a50-fac728645a0e_CJOC Test_CJOC","_score":0.0,"_source":{"plugins":[{"requiredCoreVersion":"2.346.3","active":true,"shortName":"metrics","version":"4.2.13-420.vea_2f17932dd6","enabled":true,"longName":"Metrics Plugin","hasUpdate":true,"dependencies":[{"shortName":"ionicons-api","version":"31.v4757b_6987003"},{"shortName":"jackson2-api","version":"2.13.4.20221013-295.v8e29ea_354141"},{"shortName":"variant","version":"59.vf075fe829ccb"},{"optional":true,"shortName":"instance-identity","version":"3.1"}]},{"requiredCoreVersion":"2.361.1","active":false,"shortName":"display-url-api","version":"2.3.7","enabled":false,"hasUpdate":false,"longName":"Display URL API","status":""},{"requiredCoreVersion":"2.361.2"}],"endpoint_id":"0edebf0d-6797-4ec0-9a50-fac728645a0e","created_at":"2023-12-30 15:07:54","type":"CJOC","version":"2.387.2.3","url":"https://cjoc.rosaas.releaseiq.io/","users":[{"name":"noreply","id":"noreply","type":"COMMITTER","email":"noreply@github.com"},{"name":"SYSTEM","id":"SYSTEM","type":"COMMITTER"},{"name":"releaseiq","id":"releaseiq","type":"ACCESS_USER","email":"sshaik@cloudbees.com"}],"latest_version":"","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","system_health":[{"healthy":true,"name":"plugins","message":"No failed plugins","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"thread-deadlock","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"disk-space","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"temporary-space","timestamp":"2023-12-07T12:47:28.128Z"}],"name":"CJOC Test","metrics":[{"metricsType":"gauges","metricsData":{"jenkins.job.count.value":{"value":27},"jenkins.queue.size.value":{},"jenkins.project.enabled.count.value":{"value":27},"jenkins.executor.in-use.value":{"value":1},"jenkins.node.offline.value":{},"jenkins.queue.stuck.value":{},"jenkins.executor.count.value":{"value":2},"jenkins.executor.free.value":{"value":1},"jenkins.project.count.value":{"value":27},"jenkins.queue.pending.value":{},"jenkins.project.disabled.count.value":{},"jenkins.queue.buildable.value":{},"jenkins.node.online.value":{},"jenkins.queue.blocked.value":{},"jenkins.node.count.value":{"value":1}}}],"org_name":"cloudbees-staging","timestamp":"2023-12-30 15:07:54"}}]}}`
	getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		return responseString, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processDrilldownQueryAndSpec(tt.args.replacements, tt.args.query, tt.args.reportId, tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("processDrilldownQueryAndSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processDrilldownQueryAndSpec() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("test query response decode", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "00000000-0000-0000-0000-000000000000",
			"subOrgId":  "00000000-0000-0000-0000-000000000000",
			"component": []string{"Ingestion Service"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
			"timeZone":  "Asia/Calcutta",
		}
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var resp pb.DrilldownResponse

		testQuery := constants.VulnerabilitiesOverviewDrillDownQuery
		responseBody := `{"took":16,"timed_out":false,"_shards":{"total":3,"successful":3,"skipped":0,"failed":0},"hits":{"total":{"value":58,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"drilldowns":{"value":[{"recurrences":5,"drillDown":{"reportId":"cwe-top25-vulnerabilities-view-location","reportTitle":"Vulnerabilities overview","reportInfo":{"component_id":"60d545a5-bd2f-447e-ab1d-a71b805501ff","code":"CWE-312","run_id":"c3c18621-65af-4f61-82dc-f1cf6c0dae72","scanner_name":"snyksast","branch":"https://github.com/calculi-corp/common.git"}},"component":"common","componentId":"60d545a5-bd2f-447e-ab1d-a71b805501ff","sla":"Breached","branch":"https://github.com/calculi-corp/common.git","scannerName":"snyksast","lastDiscovered":"2024/08/30 03:05:50","status":"Open","slaToolTipContent":"Breached on: 2024/05/12 21:39:12"}]}}}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseBody, nil
		}

		listVal, err := processDrilldownQueryAndSpec(testReplacements, testQuery, "vulnerabilitiesOverview", &resp)

		// Check scannerNameMap is working as expected
		valSlice := listVal.AsSlice()
		for _, val := range valSlice {
			valmap := val.(map[string]interface{})
			if scannerName, ok := valmap["scannerName"].(string); ok {
				assert.Equal(t, scannerName, "Snyk SAST")
			}
		}
		assert.NoError(t, err)
	})

	t.Run("check if index name declared on DrillDownAliasDefinitionMap", func(t *testing.T) {
		for key := range db.DrillDownQueryDefinitionMap {
			if _, exists := db.DrillDownAliasDefinitionMap[key]; !exists {
				t.Error("Index name not declared for:", key, "in DrillDownAliasDefinitionMap")
			}
		}
	})
}

func TestCiInsightsPluginsInfo(t *testing.T) {
	type args struct {
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
	}
	tests := []struct {
		name    string
		args    args
		want    *structpb.ListValue
		wantErr bool
	}{
		{
			args: args{
				replacements: map[string]any{
					"orgId":     "00000000-0000-0000-0000-000000000000",
					"ciToolId":  "10000000-0000-0000-0000-000000000000",
					"startDate": "2023-12-01 00:00:00",
					"endDate":   "2023-12-30 23:00:00",
					"jobId":     "20000000-0000-0000-0000-000000000000",
				},
			},
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	responseString := `{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":0.0,"hits":[{"_index":"cb_ci_tool_insight","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_0edebf0d-6797-4ec0-9a50-fac728645a0e_CJOC Test_CJOC","_score":0.0,"_source":{"plugins":[{"requiredCoreVersion":"2.346.3","active":true,"shortName":"metrics","version":"4.2.13-420.vea_2f17932dd6","enabled":true,"longName":"Metrics Plugin","hasUpdate":true,"dependencies":[{"shortName":"ionicons-api","version":"31.v4757b_6987003"},{"shortName":"jackson2-api","version":"2.13.4.20221013-295.v8e29ea_354141"},{"shortName":"variant","version":"59.vf075fe829ccb"},{"optional":true,"shortName":"instance-identity","version":"3.1"}]},{"requiredCoreVersion":"2.361.1","active":false,"shortName":"display-url-api","version":"2.3.7","enabled":false,"hasUpdate":false,"longName":"Display URL API","status":""},{"requiredCoreVersion":"2.361.2"}],"endpoint_id":"0edebf0d-6797-4ec0-9a50-fac728645a0e","created_at":"2023-12-30 15:07:54","type":"CJOC","version":"2.387.2.3","url":"https://cjoc.rosaas.releaseiq.io/","users":[{"name":"noreply","id":"noreply","type":"COMMITTER","email":"noreply@github.com"},{"name":"SYSTEM","id":"SYSTEM","type":"COMMITTER"},{"name":"releaseiq","id":"releaseiq","type":"ACCESS_USER","email":"sshaik@cloudbees.com"}],"latest_version":"","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","system_health":[{"healthy":true,"name":"plugins","message":"No failed plugins","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"thread-deadlock","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"disk-space","timestamp":"2023-12-07T12:47:28.128Z"},{"healthy":true,"name":"temporary-space","timestamp":"2023-12-07T12:47:28.128Z"}],"name":"CJOC Test","metrics":[{"metricsType":"gauges","metricsData":{"jenkins.job.count.value":{"value":27},"jenkins.queue.size.value":{},"jenkins.project.enabled.count.value":{"value":27},"jenkins.executor.in-use.value":{"value":1},"jenkins.node.offline.value":{},"jenkins.queue.stuck.value":{},"jenkins.executor.count.value":{"value":2},"jenkins.executor.free.value":{"value":1},"jenkins.project.count.value":{"value":27},"jenkins.queue.pending.value":{},"jenkins.project.disabled.count.value":{},"jenkins.queue.buildable.value":{},"jenkins.node.online.value":{},"jenkins.queue.blocked.value":{},"jenkins.node.count.value":{"value":1}}}],"org_name":"cloudbees-staging","timestamp":"2023-12-30 15:07:54"}}]}}`
	getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		return responseString, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CiInsightsPluginsInfo(tt.args.replacements, tt.args.ctx, tt.args.clt)
			if (err != nil) != tt.wantErr {
				t.Errorf("CiInsightsPluginsInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, len(got.Values), 3, "Validating plugins count")
		})
	}
}

func TestCiInsightsCompletedRuns(t *testing.T) {
	type args struct {
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
	}
	tests := []struct {
		name    string
		args    args
		want    *structpb.ListValue
		wantErr bool
	}{
		{
			args: args{
				replacements: map[string]any{
					"reportId":     "completedRuns",
					"orgId":        "2cab10cc-cd9d-11ed-afa1-0242ac120002",
					"durationType": "CURRENT_MONTH",
					"startDate":    "2023-12-01 00:00:00",
					"endDate":      "2023-12-30 23:00:00",
					"ciToolId":     "3cfc093d-e586-4bb5-9b67-0b6765f7d031",
					"jobId":        "5d7db1eb-8d24-4f01-9df3-8afbd7734710",
					"timeZone":     "Asia/Calcutta",
				},
			},
		},
	}
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	responseString := `{"aggregations":{"completedRuns":{"value":[{"result":"SUCCESS","duration":51440,"start_time":"2023-12-20T05:50:00.000Z","run_id":2,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"5d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"3cfc093d-e586-4bb5-9b67-0b6765f7d031","start_time_in_millis":1703031600000,"timestamp":"2023-12-30T07:45:16.000Z","url":"https://jenkins.releaseiq.io/job/TestFreestyle/88/"},{"result":"FAILED","duration":51440,"start_time":"2023-12-22T05:50:00.000Z","run_id":3,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"5d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"3cfc093d-e586-4bb5-9b67-0b6765f7d031","start_time_in_millis":1703204400000,"timestamp":"2023-12-30T07:45:16.000Z","url":"https://jenkins.releaseiq.io/job/TestFreestyle/88/"},{"result":"ABORTED","duration":51440,"start_time":"2023-12-19T05:50:00.000Z","run_id":1,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"5d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"3cfc093d-e586-4bb5-9b67-0b6765f7d031","start_time_in_millis":1702945200000,"timestamp":"2023-12-30T07:45:16.000Z","url":"https://jenkins.releaseiq.io/job/TestFreestyle/88/"},{"result":"UNSTABLE","duration":51101,"start_time":"2023-12-30T17:50:00.000Z","run_id":8,"org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"5d7db1eb-8d24-4f01-9df3-8afbd7734710","endpoint_id":"3cfc093d-e586-4bb5-9b67-0b6765f7d031","start_time_in_millis":1703958600000,"timestamp":"2023-12-30T07:47:41.000Z","url":"https://jenkins.releaseiq.io/job/TestFreestyle/88/"}]}}}`
	getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
		if IndexName == constants.CB_CI_JOB_INFO_INDEX {
			return `{"hits":{"hits":[{"_index":"cb_ci_job_info","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_3cfc093d-e586-4bb5-9b67-0b6765f7d031_259ef922-1c51-42e4-97d2-374af3e7d505","_score":0.0,"_source":{"endpoint_id":"3cfc093d-e586-4bb5-9b67-0b6765f7d031","display_name":"Demo_Folder/Demo_Folder_Pipeline","type":"Pipeline","job_name":"Demo_Folder/job/Demo_Folder_Pipeline/","updated_at":"2023-12-30 07:44:00","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","job_id":"5d7db1eb-8d24-4f01-9df3-8afbd7734710","last_completed_run_id":12,"stage_info":null,"org_name":"cloudbees-staging","timestamp":"2023-12-30 07:44:01"}}]}}`, nil
		}
		return responseString, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CiInsightsCompletedRunsAndTime(tt.args.replacements, tt.args.ctx, tt.args.clt)
			if (err != nil) != tt.wantErr {
				t.Errorf("CiInsightsCompletedRuns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, len(got.Values), 1, "Validating completed runs count")
			assert.Equal(t, len(got.Values[0].GetStructValue().Fields["completedRunsData"].GetStructValue().Fields["data"].GetListValue().Values), 4, "Validating completed runs data count")
		})
	}
}

func TestTestAutomationRunDrillDown(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"comp1", "comp2"},
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Define the response for the first call to getSearchResponse
		responseString := `{
			"aggregations": {
				"automation_run_activity": {
					"value": {
					  "9d68bb5b-65fc-4389-b706-d61953a5acb6": [
						{
							"automation_id": "9d68bb5b-65fc-4389-b706-d61953a5acb6",
							"duration": 369000,
							"component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
							"run_id": "c4c2edb0-9319-4fc1-bed9-e214f7f6a750",
							"component_name": "data-ingestion-service",
							"run_number": 306,
							"status_timestamp": "2023-09-21T21:06:08.000Z",
							"status": "Success"
						  }
					  ],
					  "1c8dfb70-5869-4b22-8889-ecfbc89c7930": [
						{
							"automation_id": "1c8dfb70-5869-4b22-8889-ecfbc89c7930",
							"duration": 0,
							"component_id": "484d5e12-6424-4070-a159-4e5639a807a2",
							"run_id": "1c42e13b-c645-4d7f-8af8-6b80da360b44",
							"component_name": "template_snyk_go",
							"run_number": 33,
							"status_timestamp": "2023-09-25T08:57:42.000Z",
							"status": "Success"
						  }
						]
					}
				}
			}
		}`

		// Define the response for the second call to getSearchResponse
		responseString1 := `{
			"aggregations": {
				"test_workflow_drilldown": {
					"value": {
						"484d5e12-6424-4070-a159-4e5639a807a2_1c8dfb70-5869-4b22-8889-ecfbc89c7930_1c42e13b-c645-4d7f-8af8-6b80da360b44": {
						"component_id": "484d5e12-6424-4070-a159-4e5639a807a2",
						"run_id": "1c42e13b-c645-4d7f-8af8-6b80da360b44",
						"component_name": "template-publish-test-results",
						"run_start_time": 1719415722000,
						"automation_id": "1c8dfb70-5869-4b22-8889-ecfbc89c7930",
						"branch_id": "14b2820b-a3c4-41a7-8431-caff68165151",
						"org_id": "8509888e-d27f-44fa-46a9-29bc76f5e790",
						"branch_name": "preprod",
						"run_number": "283",
						"run_status": "SUCCEEDED",
						"automation_name": "workflow",
						"runs": 2,
						"test_suite_name": "xml-report",
						"status": "FAILED"
						}
					},
					"value": {
						"e2f8fef6-5041-4843-b37e-6cdae38099bc_9d68bb5b-65fc-4389-b706-d61953a5acb6_c4c2edb0-9319-4fc1-bed9-e214f7f6a750": {
						"component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
						"run_id": "c4c2edb0-9319-4fc1-bed9-e214f7f6a750",
						"component_name": "template-publish-test-results",
						"run_start_time": 1719415722000,
						"automation_id": "9d68bb5b-65fc-4389-b706-d61953a5acb6",
						"branch_id": "14b2820b-a3c4-41a7-8431-caff68165151",
						"org_id": "8509888e-d27f-44fa-46a9-29bc76f5e790",
						"branch_name": "preprod",
						"run_number": "283",
						"run_status": "SUCCEEDED",
						"automation_name": "workflow",
						"runs": 2,
						"test_suite_name": "xml-report",
						"status": "Passed"
						}
					}
				}
			}
		}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			if IndexName == "automation_run_status" {
				return responseString, nil
			} else if IndexName == "cb_test_suites" {
				return responseString1, nil
			} else {
				return "", errors.New("Unknown IndexName")
			}
		}
		reports, err := TestAutomationRunDrillDown(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, 2, len(reports.Values), "Validating response count for test automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 5: No suborgId provided", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background() // Use a context for testing.

		_, err := TestAutomationRunDrillDown(testReplacements, ctx, mockGrpcClient)

		validationErr := fmt.Errorf(errMissingRequiredField, constants.SUB_ORG_ID)

		expectedValidationStatus := status.Errorf(codes.InvalidArgument, "ReportServiceRequest Validation failed: %s", validationErr.Error())

		assert.NotEqual(t, expectedValidationStatus, err, "org Id should not be null")
	})

}

func TestTestInsightsViewRunActivityDrillDown(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Success case - View Run Activity drill down", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"test1"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Mocking the response from OpenSearch
		responseString := `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"viewRunActivity":{"value":{"headers":{"avg_duration":100,"total":3,"FAILED":0,"duration_array":[100,100,100],"SKIPPED":0,"PASSED":3},"section":[{"test_case_name":"Test","start_time":1718961841000,"run_id":"35","build_id":"173","job_id":"j","test_case_status":"PASSED","test_case_duration":100,"test_suite_name":"test_suite_1"},{"test_case_name":"Test2","start_time":1718966799000,"run_id":"6d","build_id":"186","job_id":"tjj","test_case_status":"PASSED","test_case_duration":100,"test_suite_name":"test_suite_1"},{"test_case_name":"Test3","start_time":1718984187000,"run_id":"3b","build_id":"206","job_id":"j1","test_case_status":"PASSED","test_case_duration":100,"test_suite_name":"test_suite_1"}]}}}}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		reports, err := TestInsightsViewRunActivityDrillDown(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, 1, len(reports.Values), "Validating response count for View Run Activity drill down")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

	t.Run("Case 2: Failure case - View Run Activity drill down unmarshalling error", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Mocking the response from OpenSearch
		responseString := `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":10,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test":{"value":{"headers":{"avg_duration":100,"total":3,"FAILED":0,"duration_array":[100,100,100],"SKIPPED":0,"PASSED":3},"section":[{"test_case_name":"Test","start_time":1718961841000,"run_id":"35","build_id":"173","job_id":"j","test_case_status":"PASSED","test_case_duration":100,"test_suite_name":"test_suite_1"},{"test_case_name":"Test2","start_time":1718966799000,"run_id":"6d","build_id":"186","job_id":"tjj","test_case_status":"PASSED","test_case_duration":100,"test_suite_name":"test_suite_1"},{"test_case_name":"Test3","start_time":1718984187000,"run_id":"3b","build_id":"206","job_id":"j1","test_case_status":"PASSED","test_case_duration":100,"test_suite_name":"test_suite_1"}]}}}}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		_, err := TestInsightsViewRunActivityDrillDown(testReplacements, ctx, mockGrpcClient)

		expectedErr := db.ErrInternalServer

		assert.Equal(t, expectedErr, err, "opensearch response is not as expected")
		// Assertions to check if the result is as expected.
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}

	})

}

func TestTestInsightsViewRunActivityLogsDrillDown(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Success case - View Run Activity Logs drill down", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"test1"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Mocking the response from OpenSearch
		responseString := `{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"1d","_score":0,"_source":{"std_err":"","run_number":173,"error_trace":"Error Trace Dummy","std_out":"Std Out Dummy"}}]}}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		reports, err := TestInsightsViewRunActivityLogsDrillDown(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, 1, len(reports.Values), "Validating response count for View Run Activity Logs drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 2: Failure case - View Run Activity Logs drill down unmarshalling error", func(t *testing.T) {
		testReplacements := map[string]any{
			"OrgId":     "00000000-0000-0000-0000-000000000000",
			"component": []string{"All"},
			"startDate": "2023-06-01 10:00:00",
			"endDate":   "2023-06-01 10:00:00",
		}
		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Mocking the response from OpenSearch
		responseString := `{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"test-hit":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"dr","_score":0,"_source":{"std_err":"","run_number":173,"error_trace":"Error Trace Dummy","std_out":"Std Out Dummy"}}]}}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		_, err := TestInsightsViewRunActivityLogsDrillDown(testReplacements, ctx, mockGrpcClient)

		expectedErr := db.ErrInternalServer

		assert.Equal(t, expectedErr, err, "opensearch response is not as expected")
	})

}

func TestPullRequestsReport(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: All inputs are present", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"component":  []string{"comp1", "comp2"},
			"startDate":  "2023-01-01",
			"timeFormat": "12h",
			"endDate":    "2023-12-31",
		}
		ctx := context.Background()

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"pullrequests": {
					"value": [
					  {
						"component_id": "7dfaabc0-5eb0-40c6-427c-a2ac4cc01da3",
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"pr_created_time": "2023-10-10T18:43:33.000Z",
						"component_name": "config",
						"target_branch": "",
						"review_status": "OPEN",
						"pull_request_id": "26",
						"repository_name": "calculi-corp/config",
						"timestamp": "2023-10-10T18:45:39.000Z",
						"source_branch": "",
						"provider":"p1",
						"repository_url": "https://github.com/calculi-corp/config"
					  },
					  {
						"component_id": "785894c9-91cf-47b9-734d-ae73377cc29c",
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"pr_created_time": "2023-10-11T13:02:31.000Z",
						"component_name": "run-service",
						"target_branch": "main",
						"review_status": "APPROVED",
						"pull_request_id": "125",
						"repository_name": "calculi-corp/run-service",
						"timestamp": "2023-10-12T14:18:45.000Z",
						"source_branch": "SDP-9645-anonymous-id",
						"provider":"p1",
						"repository_url": "https://github.com/calculi-corp/run-service"
					  },
					  {
						"component_id": "716d4922-add7-474b-7f21-240a8e253c38",
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"pr_created_time": "2023-10-12T16:08:38.000Z",
						"component_name": "jwt-validator",
						"target_branch": "main",
						"review_status": "APPROVED",
						"pull_request_id": "46",
						"repository_name": "calculi-corp/jwt-validator",
						"timestamp": "2023-10-12T16:12:28.000Z",
						"source_branch": "SDP-9900",
						"provider":"p1",
						"repository_url": "https://github.com/calculi-corp/jwt-validator"
					  },
					  {
						"component_id": "55f26650-b259-4453-6ff3-c2843d096780",
						"org_id": "2cab10cc-cd9d-11ed-afa1-0242ac120002",
						"pr_created_time": "2023-10-04T18:24:10.000Z",
						"component_name": "api",
						"target_branch": "",
						"review_status": "OPEN",
						"pull_request_id": "451",
						"repository_name": "calculi-corp/api",
						"timestamp": "2023-10-04T18:24:23.000Z",
						"source_branch": "",
						"provider":"p1",
						"repository_url": "https://github.com/calculi-corp/api"
					  }
					]
				}
			}
		  }`
		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		PullRequestsReport(testReplacements, ctx, mockGrpcClient)

	})

}

func Test_RunDetailsTestResults(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: Successful execution for Test Suite view", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"component":  []string{"comp1", "comp2"},
			"startDate":  "2023-01-01",
			"endDate":    "2023-12-31",
			"viewOption": "TEST_SUITES_VIEW",
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Define the response for the first call to getSearchResponse
		responseString := `{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":80,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_suite_buckets":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"cm","doc_count":1,"test_suite_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id1","_score":0,"_source":{"duration":0,"total":1,"component_id":"cid1","run_id":"rid1","passed":1,"failed":0,"skipped":0,"test_suite_name":"common"}}]}}},{"key":"ev","doc_count":1,"test_suite_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id2","_score":0,"_source":{"duration":0,"total":14,"component_id":"cid2","run_id":"rid2","passed":14,"failed":0,"skipped":0,"test_suite_name":"evn"}}]}}}]}}}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		reports, err := RunDetailsTestResults(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for Run Details Test Suites view")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})

	t.Run("Case 2: Successful execution for Test Cases view", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":      "org123",
			"subOrgId":   "suborg123",
			"component":  []string{"comp1", "comp2"},
			"startDate":  "2023-11-01",
			"endDate":    "2023-12-31",
			"viewOption": "TEST_CASES_VIEW",
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Define the response for the first call to getSearchResponse
		responseString := `{"took":48,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1591,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_case_buckets":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"jj","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"jj1","_score":0,"_source":{"test_case_name":"t1","duration":0,"component_id":"c1","run_id":"r1","test_suite_name":"en","status":"PASSED"}}]}}},{"key":"pt","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"psi","_score":0,"_source":{"test_case_name":"ps","duration":0,"component_id":"c2","run_id":"r2","test_suite_name":"en","status":"PASSED"}}]}}}]}}}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		reports, err := RunDetailsTestResults(testReplacements, ctx, mockGrpcClient)

		assert.Equal(t, len(reports.Values), 2, "Validating response count for Run Details Test Cases view")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})
}

func Test_transformRunDetailsTestResultsTestSuitesView(t *testing.T) {
	t.Run("Case 1: Successful transformation for Test Suites view", func(t *testing.T) {
		// Response from OpenSearch
		responseString := `{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":80,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_suite_buckets":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"cm","doc_count":1,"test_suite_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id1","_score":0,"_source":{"duration":0,"total":1,"component_id":"cid1","run_id":"rid1","passed":1,"failed":0,"skipped":0,"test_suite_name":"common"}}]}}},{"key":"ev","doc_count":1,"test_suite_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id2","_score":0,"_source":{"duration":0,"total":14,"component_id":"cid2","run_id":"rid2","passed":0,"failed":14,"skipped":0,"test_suite_name":"evn"}}]}}}]}}}`
		data, err := transformRunDetailsTestResultsTestSuitesView(responseString)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualTransformedData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string instead of structpb.ListValue
		expectedTransformedData := `[{"drillDown":{"reportId":"run-details-total-test-cases","reportInfo":{"component_id":"cid2","run_id":"rid2","run_number":"0","test_suite_name":"evn"},"reportTitle":"Test cases - evn","reportType":"status"},"runTime":0,"testCasesFailed":14,"testCasesPassed":0,"testCasesSkipped":0,"testSuiteName":"evn","totalTestCases":14},{"drillDown":{"reportId":"run-details-total-test-cases","reportInfo":{"component_id":"cid1","run_id":"rid1","run_number":"0","test_suite_name":"common"},"reportTitle":"Test cases - common","reportType":"status"},"runTime":0,"testCasesFailed":0,"testCasesPassed":1,"testCasesSkipped":0,"testSuiteName":"common","totalTestCases":1}]`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedTransformedData, string(actualTransformedData), "Validating transformed data for Test Suites view")

	})

}

func Test_transformRunDetailsTestResultsTestCasesView(t *testing.T) {
	t.Run("Case 1: Successful transformation for Test Cases view", func(t *testing.T) {
		// Response from OpenSearch
		responseString := `{"took":48,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1591,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_case_buckets":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"jj","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"jj1","_score":0,"_source":{"test_case_name":"t1","duration":0,"component_id":"c1","run_id":"r1","test_suite_name":"en","status":"PASSED","std_err":"","error_trace":"","std_out":""}}]}}},{"key":"pt","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"psi","_score":0,"_source":{"test_case_name":"ps","duration":0,"component_id":"c2","run_id":"r2","test_suite_name":"en","status":"PASSED","std_err":"","error_trace":"Checking for isLogReported flag set to false","std_out":""}}]}}}]}}}`
		data, err := transformRunDetailsTestResultsTestCasesView(responseString)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualTransformedData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string instead of structpb.ListValue
		expectedTransformedData := `[{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c1","run_id":"r1","run_number":"0","test_case_name":"t1","test_suite_name":"en"},"reportTitle":"Test case log - t1","reportType":""},"isLogReported":false,"runTime":0,"status":"Passed","testCaseName":"t1","testSuiteName":"en"},{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c2","run_id":"r2","run_number":"0","test_case_name":"ps","test_suite_name":"en"},"reportTitle":"Test case log - ps","reportType":""},"isLogReported":true,"runTime":0,"status":"Passed","testCaseName":"ps","testSuiteName":"en"}]`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedTransformedData, string(actualTransformedData), "Validating transformed data for Test Suites view")

	})

	t.Run("Case 2: Verifying sort order based on the status field", func(t *testing.T) {
		// Response from OpenSearch
		responseString := `{"took":48,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1591,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_case_buckets":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"jj","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"jj1","_score":0,"_source":{"test_case_name":"t1","duration":0,"component_id":"c1","run_id":"r1","test_suite_name":"en","status":"SKP","std_err":"","error_trace":"","std_out":""}}]}}},{"key":"pt","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"psi","_score":0,"_source":{"test_case_name":"ps","duration":0,"component_id":"c2","run_id":"r2","test_suite_name":"en","status":"FAILED","std_err":"","error_trace":"Checking for isLogReported flag set to false","std_out":""}}]}}}]}}}`
		data, err := transformRunDetailsTestResultsTestCasesView(responseString)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualTransformedData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string instead of structpb.ListValue
		expectedTransformedData := `[{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c2","run_id":"r2","run_number":"0","test_case_name":"ps","test_suite_name":"en"},"reportTitle":"Test case log - ps","reportType":""},"isLogReported":true,"runTime":0,"status":"Failed","testCaseName":"ps","testSuiteName":"en"},{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c1","run_id":"r1","run_number":"0","test_case_name":"t1","test_suite_name":"en"},"reportTitle":"Test case log - t1","reportType":""},"isLogReported":false,"runTime":0,"status":"Skp","testCaseName":"t1","testSuiteName":"en"}]`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedTransformedData, string(actualTransformedData), "Validating transformed data for Test Suites view")

	})

}

func Test_RunDetailsTotalTestCasesDrillDown(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: Successful execution fof the Total Test Cases drill down", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"comp1", "comp2"},
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Define the response for the first call to getSearchResponse
		responseString := `{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":38,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_case_buckets":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"tc1","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id1","_score":0,"_source":{"test_case_name":"tc1","duration":0,"component_id":"c1","run_id":"r1","test_suite_name":"ts1","status":"PASSED"}}]}}},{"key":"tc2","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id2","_score":0,"_source":{"test_case_name":"tc2","duration":0,"component_id":"c2","run_id":"r2","test_suite_name":"ts2","status":"PASSED"}}]}}},{"key":"tc3","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id3","_score":0,"_source":{"test_case_name":"tc3","duration":0,"component_id":"c3","run_id":"r3","test_suite_name":"ts3","status":"PASSED"}}]}}}]}}}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		reports, err := RunDetailsTotalTestCasesDrillDown(testReplacements, ctx, mockGrpcClient)
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		// Assertions to check if the result is as expected.
		assert.Equal(t, 3, len(reports.Values), "Validating response count for Run Details Total Test Cases drill down")

	})

}

func Test_transformRunDetailsTotalTestCasesDrillDown(t *testing.T) {
	t.Run("Case 1: Successful transformation for Total Test Cases drill down", func(t *testing.T) {
		// Response from OpenSearch
		responseString := `{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":38,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_case_buckets":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"tc1","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id1","_score":0,"_source":{"test_case_name":"tc1","duration":0,"component_id":"c1","run_id":"r1","test_suite_name":"ts1","status":"PASSED","std_err":"","error_trace":"","std_out":""}}]}}},{"key":"tc2","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id2","_score":0,"_source":{"test_case_name":"tc2","duration":0,"component_id":"c2","run_id":"r2","test_suite_name":"ts2","status":"PASSED","std_err":"","error_trace":"","std_out":""}}]}}},{"key":"tc3","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id3","_score":0,"_source":{"test_case_name":"tc3","duration":0,"component_id":"c3","run_id":"r3","test_suite_name":"ts3","status":"PASSED","std_err":"Logs available","error_trace":"","std_out":""}}]}}}]}}}`
		data, err := transformRunDetailsTotalTestCasesDrillDown(responseString)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualTransformedData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string instead of structpb.ListValue
		expectedTransformedData := `[{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c1","run_id":"r1","run_number":"0","test_case_name":"tc1","test_suite_name":"ts1"},"reportTitle":"Test case log - tc1","reportType":""},"isLogReported":false,"runTime":0,"status":"Passed","testCaseName":"tc1"},{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c2","run_id":"r2","run_number":"0","test_case_name":"tc2","test_suite_name":"ts2"},"reportTitle":"Test case log - tc2","reportType":""},"isLogReported":false,"runTime":0,"status":"Passed","testCaseName":"tc2"},{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c3","run_id":"r3","run_number":"0","test_case_name":"tc3","test_suite_name":"ts3"},"reportTitle":"Test case log - tc3","reportType":""},"isLogReported":true,"runTime":0,"status":"Passed","testCaseName":"tc3"}]`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedTransformedData, string(actualTransformedData), "Validating transformed data for Run Details Total Test Cases drill down")

	})

	t.Run("Case 2: Verifying sort order based on the status field", func(t *testing.T) {
		// Response from OpenSearch
		responseString := `{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":38,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_case_buckets":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"tc1","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id1","_score":0,"_source":{"test_case_name":"tc1","duration":0,"component_id":"c1","run_id":"r1","test_suite_name":"ts1","status":"SKIPPED","std_err":"","error_trace":"","std_out":""}}]}}},{"key":"tc2","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id2","_score":0,"_source":{"test_case_name":"tc2","duration":0,"component_id":"c2","run_id":"r2","test_suite_name":"ts2","status":"PASSED","std_err":"","error_trace":"","std_out":""}}]}}},{"key":"tc3","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id3","_score":0,"_source":{"test_case_name":"tc3","duration":0,"component_id":"c3","run_id":"r3","test_suite_name":"ts3","status":"FAILED","std_err":"Logs available","error_trace":"","std_out":""}}]}}},{"key":"tc4","doc_count":1,"test_case_doc":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id4","_score":0,"_source":{"test_case_name":"tc4","duration":0,"component_id":"c4","run_id":"r4","test_suite_name":"ts4","status":"FAIL","std_err":"Logs available","error_trace":"","std_out":""}}]}}}]}}}`
		data, err := transformRunDetailsTotalTestCasesDrillDown(responseString)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualTransformedData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string instead of structpb.ListValue
		expectedTransformedData := `[{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c3","run_id":"r3","run_number":"0","test_case_name":"tc3","test_suite_name":"ts3"},"reportTitle":"Test case log - tc3","reportType":""},"isLogReported":true,"runTime":0,"status":"Failed","testCaseName":"tc3"},{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c1","run_id":"r1","run_number":"0","test_case_name":"tc1","test_suite_name":"ts1"},"reportTitle":"Test case log - tc1","reportType":""},"isLogReported":false,"runTime":0,"status":"Skipped","testCaseName":"tc1"},{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c2","run_id":"r2","run_number":"0","test_case_name":"tc2","test_suite_name":"ts2"},"reportTitle":"Test case log - tc2","reportType":""},"isLogReported":false,"runTime":0,"status":"Passed","testCaseName":"tc2"},{"drillDown":{"reportId":"run-details-test-case-log","reportInfo":{"component_id":"c4","run_id":"r4","run_number":"0","test_case_name":"tc4","test_suite_name":"ts4"},"reportTitle":"Test case log - tc4","reportType":""},"isLogReported":true,"runTime":0,"status":"Fail","testCaseName":"tc4"}]`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedTransformedData, string(actualTransformedData), "Validating transformed data for Run Details Total Test Cases drill down")

	})

}

func Test_RunDetailsTestCaseLogDrillDown(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: Successful execution fof the Test Case Log drill down", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"comp1", "comp2"},
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Define the response for the first call to getSearchResponse
		responseString := `{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id1","_score":0,"_source":{"std_err":"Test log 1","error_trace":"Test log 2","std_out":"Test log 3"}}]}}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}
		reports, err := RunDetailsTestCaseLogDrillDown(testReplacements, ctx, mockGrpcClient)
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		// Assertions to check if the result is as expected.
		assert.Equal(t, 1, len(reports.Values), "Validating response count for Run Details Test Case Log drill down")

	})

}

func Test_transformRunDetailsTestCaseLogDrillDown(t *testing.T) {
	t.Run("Case 1: Successful transformation for Total Test Cases drill down", func(t *testing.T) {
		// Response from OpenSearch
		responseString := `{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"test","_id":"id1","_score":0,"_source":{"std_err":"Test log 1","error_trace":"Test log 2","std_out":"Test log 3"}}]}}`
		data, err := transformRunDetailsTestCaseLogDrillDown(responseString)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualTransformedData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string instead of structpb.ListValue
		expectedTransformedData := `[{"logDetails":{"message":"Test log 3\nTest log 1\nTest log 2\n"}}]`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedTransformedData, string(actualTransformedData), "Validating transformed data for Run Details Test Case Log drill down")

	})

}

func Test_transformTestInsightsTotalRunsDrillDown(t *testing.T) {
	t.Run("Case 1: Successful transformation for Total Test Cases drill down - no sub rows", func(t *testing.T) {
		// Response from OpenSearch
		responseString := `{"took":8,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":70,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_cases_that_failed_at_least_once":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[]},"runs":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"k1","doc_count":35,"run_details":{"hits":{"total":{"value":35,"relation":"eq"},"max_score":0,"hits":[{"_index":"tc","_id":"id1","_score":0,"_source":{"automation_id":"a1","component_id":"c1","run_id":"r1","branch_id":"b1","org_id":"o1","run_number":527,"run_status":"SUCCEEDED"},"fields":{"zoned_run_start_time":["2024/09/12 18:04"]}}]}},"failed_docs":{"doc_count":0,"failed_test_cases":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[]}},"total_test_cases":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"tc1","doc_count":1},{"key":"tc2","doc_count":1}]},"total_test_cases_count":{"value":35}},{"key":"k2","doc_count":35,"run_details":{"hits":{"total":{"value":35,"relation":"eq"},"max_score":0,"hits":[{"_index":"tc","_id":"id2","_score":0,"_source":{"automation_id":"a2","component_id":"c2","run_id":"r2","branch_id":"b2","org_id":"o1","run_number":542,"run_status":"SUCCEEDED"},"fields":{"zoned_run_start_time":["2024/09/16 14:41"]}}]}},"failed_docs":{"doc_count":0,"failed_test_cases":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[]}},"total_test_cases":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"tc1","doc_count":1},{"key":"tc2","doc_count":1}]},"total_test_cases_count":{"value":35}}]}}}`
		data, err := transformTestInsightsTotalRunsDrillDown(responseString)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualTransformedData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string instead of structpb.ListValue
		expectedTransformedData := `[{"automationId":"a2","branchId":"b2","build":"542","componentId":"c2","failedTests":0,"organizationId":"o1","runId":"r2","runStatus":"Success","runTime":"2024/09/16 14:41","subRows":[],"totalTests":35},{"automationId":"a1","branchId":"b1","build":"527","componentId":"c1","failedTests":0,"organizationId":"o1","runId":"r1","runStatus":"Success","runTime":"2024/09/12 18:04","subRows":[],"totalTests":35}]`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedTransformedData, string(actualTransformedData), "Validating transformed data for transformTestInsightsTotalRunsDrillDown")

	})

	t.Run("Case 2: Successful transformation for Total Test Cases drill down - with sub rows", func(t *testing.T) {
		// Response from OpenSearch
		responseString := `{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":6,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_cases_that_failed_at_least_once":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"tc1","doc_count":1,"test_case_name":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"tc","_id":"id1","_score":0,"_source":{"test_case_name":"tc1","test_suite_name":"ts1"}}]}},"test_case_status_history":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"FAILED","doc_count":1}]},"failed_docs":{"doc_count":1}},{"key":"tc2","doc_count":1,"test_case_name":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"tc","_id":"id2","_score":0,"_source":{"test_case_name":"tc2","test_suite_name":"ts2"}}]}},"test_case_status_history":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"FAILED","doc_count":1}]},"failed_docs":{"doc_count":1}}]},"runs":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"r1","doc_count":6,"run_details":{"hits":{"total":{"value":6,"relation":"eq"},"max_score":0,"hits":[{"_index":"tc","_id":"id1","_score":0,"_source":{"automation_id":"a1","component_id":"c1","run_id":"r1","branch_id":"b1","org_id":"o1","run_number":733,"run_status":"SUCCEEDED"},"fields":{"zoned_run_start_time":["2024/09/12 12:50"]}}]}},"failed_docs":{"doc_count":2,"failed_test_cases":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"tc1","doc_count":1},{"key":"tc2","doc_count":1}]}},"total_test_cases":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"tc1","doc_count":1},{"key":"tc2","doc_count":1},{"key":"tc3","doc_count":1}]},"total_test_cases_count":{"value":3}}]}}}`
		data, err := transformTestInsightsTotalRunsDrillDown(responseString)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualTransformedData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string instead of structpb.ListValue
		expectedTransformedData := `[{"automationId":"a1","branchId":"b1","build":"733","componentId":"c1","failedTests":2,"organizationId":"o1","runId":"r1","runStatus":"Success","runTime":"2024/09/12 12:50","subRows":[{"failedTestName":"tc1","failureRate":{"colorScheme":[{"color0":"#009C5B","color1":"#62CA9D"},{"color0":"#D32227","color1":"#FB6E72"},{"color0":"#F2A414","color1":"#FFE6C1"}],"data":[{"title":"Successful runs","value":0},{"title":"Failed runs","value":1},{"title":"Skipped runs","value":0}],"lightColorScheme":[{"color0":"#0C9E61","color1":"#79CAA8"},{"color0":"#E83D39","color1":"#F39492"},{"color0":"#F2A414","color1":"#FFE6C1"}],"type":"SINGLE_BAR","value":"100.0%"},"failureRateValue":100,"viewRunActivity":{"drillDown":{"reportId":"test-overview-view-run-activity","reportInfo":{"automation_id":"a1","branch":"b1","component_id":"c1","run_id":"r1","run_number":"733","test_case_name":"tc1","test_suite_name":"ts1"},"reportTitle":"Test case activity - tc1","reportType":""}}},{"failedTestName":"tc2","failureRate":{"colorScheme":[{"color0":"#009C5B","color1":"#62CA9D"},{"color0":"#D32227","color1":"#FB6E72"},{"color0":"#F2A414","color1":"#FFE6C1"}],"data":[{"title":"Successful runs","value":0},{"title":"Failed runs","value":1},{"title":"Skipped runs","value":0}],"lightColorScheme":[{"color0":"#0C9E61","color1":"#79CAA8"},{"color0":"#E83D39","color1":"#F39492"},{"color0":"#F2A414","color1":"#FFE6C1"}],"type":"SINGLE_BAR","value":"100.0%"},"failureRateValue":100,"viewRunActivity":{"drillDown":{"reportId":"test-overview-view-run-activity","reportInfo":{"automation_id":"a1","branch":"b1","component_id":"c1","run_id":"r1","run_number":"733","test_case_name":"tc2","test_suite_name":"ts2"},"reportTitle":"Test case activity - tc2","reportType":""}}}],"totalTests":3}]`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedTransformedData, string(actualTransformedData), "Validating transformed data for transformTestInsightsTotalRunsDrillDown")

	})

	t.Run("Case 3: Error when there's no data found", func(t *testing.T) {
		// Response from OpenSearch
		responseString := `{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":0,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_cases_that_failed_at_least_once":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[]},"runs":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[]}}}`
		_, err := transformTestInsightsTotalRunsDrillDown(responseString)

		// Assertions to check if the result is as expected.
		assert.Equal(t, err, db.ErrNoDataFound, "Validating error message for no data found in transformTestInsightsTotalRunsDrillDown")

	})

}

func Test_RunDetailsTestResultsIndicators(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	t.Run("Case 1: Successful execution when there's no data for the component", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"comp1"},
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		reports, err := RunDetailsTestResultsIndicators(testReplacements, ctx, mockGrpcClient)
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
		actualData, err := json.Marshal(reports)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string
		expectedData := `[{"isTestInsightsDataFound":false,"testCasesFailed":0,"testCasesPassed":0,"testCasesSkipped":0}]`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedData, string(actualData), "Validating response count for RunDetailsTestResultsIndicators")

	})

	t.Run("Case 2: Successful execution when there's data for the component", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"comp1"},
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		// Define the response to getSearchResponse
		responseString := `{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":2,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_suites":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"ts1","doc_count":1,"statuses":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"t","_id":"id1","_score":0,"_source":{"passed":1,"failed":0,"skipped":2}}]}}},{"key":"ts2","doc_count":1,"statuses":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"t","_id":"id2","_score":0,"_source":{"passed":14,"failed":1,"skipped":0}}]}}}]}}}`

		getSearchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		getDocCount = func(IndexName string, orgId string, components []string) int {
			return 2
		}

		reports, err := RunDetailsTestResultsIndicators(testReplacements, ctx, mockGrpcClient)
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
		actualData, err := json.Marshal(reports)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string
		expectedData := `[{"isTestInsightsDataFound":true,"testCasesFailed":1,"testCasesPassed":15,"testCasesSkipped":2}]`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedData, string(actualData), "Validating response count for RunDetailsTestResultsIndicators")

	})

}

func Test_transformRunDetailsTestResultsIndicators(t *testing.T) {
	t.Run("Case 1: Successful transformation for when the aggregation returns data for a particular run", func(t *testing.T) {
		// Response from OpenSearch
		output := RunDetailsTestResultsIndicatorsResponse{
			IsTestInsightsDataFound: true,
		}
		responseString := `{"took":4,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":2,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_suites":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"ts1","doc_count":1,"statuses":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"t","_id":"id1","_score":0,"_source":{"passed":1,"failed":0,"skipped":2}}]}}},{"key":"ts2","doc_count":1,"statuses":{"hits":{"total":{"value":1,"relation":"eq"},"max_score":0,"hits":[{"_index":"t","_id":"id2","_score":0,"_source":{"passed":14,"failed":1,"skipped":0}}]}}}]}}}`
		err := transformRunDetailsTestResultsIndicators(responseString, &output)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualTransformedData, err := json.Marshal(output)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string
		expectedTransformedData := `{"testCasesFailed":1,"testCasesPassed":15,"testCasesSkipped":2,"isTestInsightsDataFound":true}`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedTransformedData, string(actualTransformedData), "Validating transformed data for transformRunDetailsTestResultsIndicators")

	})

	t.Run("Case 2: Successful transformation for when the aggregation returns no data for a particular", func(t *testing.T) {
		// Response from OpenSearch
		output := RunDetailsTestResultsIndicatorsResponse{
			IsTestInsightsDataFound: true,
		}
		responseString := `{"took":2,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":0,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"test_suites":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[]}}}`
		err := transformRunDetailsTestResultsIndicators(responseString, &output)
		if err != nil {
			t.Errorf("Failed to transform data: %v", err)
		}

		actualTransformedData, err := json.Marshal(output)
		if err != nil {
			t.Errorf("Failed to marshal transformed data: %v", err)
		}

		// Expected data after transformation as a string
		expectedTransformedData := `{"testCasesFailed":0,"testCasesPassed":0,"testCasesSkipped":0,"isTestInsightsDataFound":true}`

		// Assertions to check if the result is as expected.
		assert.Equal(t, expectedTransformedData, string(actualTransformedData), "Validating transformed data for transformRunDetailsTestResultsIndicators")

	})

}
