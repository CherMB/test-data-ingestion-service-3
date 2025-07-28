package internal

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"testing"

	api "github.com/calculi-corp/api/go"
	scanner "github.com/calculi-corp/api/go/scanner"
	"github.com/calculi-corp/api/go/service"
	"github.com/calculi-corp/common/grpc"
	"github.com/calculi-corp/config"
	coredataMock "github.com/calculi-corp/core-data-cache/mock"
	client "github.com/calculi-corp/grpc-client"
	mock_grpc_client "github.com/calculi-corp/grpc-client/mock"
	cache "github.com/calculi-corp/reports-service/cache"
	"github.com/calculi-corp/reports-service/constants"
	"github.com/calculi-corp/reports-service/db"
	"github.com/opensearch-project/opensearch-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func init() {
	config.Config.Set("logging.level", "INFO")
}

func Test_fetchAllDeployedEnvironments(t *testing.T) {
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
		}
		ctx := context.Background()

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"environments": {
				  "value": [
					"stage",
					"staging"
				  ]
				}
			  }
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		reports, _, err := FetchAllDeployedEnvironments("99", testReplacements, ctx)

		assert.Equal(t, len(reports), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

	t.Run("Case 1: Successful execution", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
			"component": []string{"comp1", "comp2"},
		}
		ctx := context.Background()

		var client *opensearch.Client

		openSearchClient = func() (*opensearch.Client, error) {
			return client, errors.New("error")
		}
		responseString := `{
			"aggregations": {
				"environments": {
				  "value": [
					"stage",
					"staging"
				  ]
				}
			  }
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		reports, _, err := FetchAllDeployedEnvironments("99", testReplacements, ctx)

		assert.Equal(t, len(reports), 2, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

	t.Run("Case 1: Successful execution", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
			"component": []string{"comp1", "comp2"},
		}
		ctx := context.Background()

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"environments": {
				  "value": [
					"stage",
					"staging"
				  ]
				}
			  }
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, errors.New("error")
		}

		_, _, _ = FetchAllDeployedEnvironments("99", testReplacements, ctx)

	})

}

func Test_fetchAllBuildedComponents(t *testing.T) {
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
		}
		ctx := context.Background()

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
			  "components": {
				"value": {
				  "components": [
					"1358c14d-6d31-4b81-974e-5287e51b85fd",
					"a1d40803-7978-40ed-6c6e-d6fa114cbfa7",
					"238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
					"addf84e3-5c30-49d0-aaa7-1564ee109c0b",
					"24c2c304-75b1-4076-9117-30d5d375ab15"
				  ],
				  "min": "0",
				  "max": "0"
				}
			  }
			}
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		reports, _, err := fetchAllBuildedComponents("99", testReplacements, ctx)

		assert.Equal(t, len(reports), 5, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

	t.Run("Case 1: Successful execution", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
			"component": []string{"comp1", "comp2"},
		}
		ctx := context.Background()

		var client *opensearch.Client
		openSearchClient = func() (*opensearch.Client, error) {
			return client, errors.New("error")
		}
		responseString := `{
			"aggregations": {
			  "components": {
				"value": {
				  "components": [
					"1358c14d-6d31-4b81-974e-5287e51b85fd",
					"a1d40803-7978-40ed-6c6e-d6fa114cbfa7",
					"238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
					"addf84e3-5c30-49d0-aaa7-1564ee109c0b",
					"24c2c304-75b1-4076-9117-30d5d375ab15"
				  ],
				  "min": "0",
				  "max": "0"
				}
			  }
			}
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		reports, _, err := fetchAllBuildedComponents("99", testReplacements, ctx)

		assert.Equal(t, len(reports), 5, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

	t.Run("Case 1: Successful execution", func(t *testing.T) {
		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"startDate": "2023-01-01",
			"endDate":   "2023-12-31",
			"component": []string{"comp1", "comp2"},
		}
		ctx := context.Background()

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
			  "components": {
				"value": {
				  "components": [
					"1358c14d-6d31-4b81-974e-5287e51b85fd",
					"a1d40803-7978-40ed-6c6e-d6fa114cbfa7",
					"238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
					"addf84e3-5c30-49d0-aaa7-1564ee109c0b",
					"24c2c304-75b1-4076-9117-30d5d375ab15"
				  ],
				  "min": "0",
				  "max": "0"
				}
			  }
			}
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, errors.New("error")
		}

		_, _, _ = fetchAllBuildedComponents("99", testReplacements, ctx)

	})

}

func TestExecuteFunction(t *testing.T) {
	// Mock GrpcClient and any other necessary dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	FunctionMap = map[string]interface{}{
		"Security Widget Section": securityComponentWidgetSection,
		// Add other function mappings as needed
	}

	tests := []struct {
		name         string
		k            string
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          *mock_grpc_client.MockGrpcClient // Use the mock client type
		wantErr      bool
		expected     json.RawMessage
	}{
		{
			name:     "Case 1: Valid function with valid parameters",
			k:        "Security Widget Section",
			widgetId: "widget123",
			replacements: map[string]interface{}{
				constants.ORG_ID: "org124",
				"component":      []string{"testComp"},
			},
			ctx:      context.TODO(),
			clt:      mockGrpcClient,
			wantErr:  false,
			expected: json.RawMessage(`{"key": "value"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result, err := ExecuteFunction(tt.k, tt.widgetId, tt.replacements, tt.ctx, tt.clt, nil)

			if tt.wantErr && err == nil {
				t.Errorf("Expected an error, but got none.")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, but got an error: %v", err)
			}

			if tt.expected != nil {
				assert.Equal(t, json.RawMessage(nil), result, "No error and response validation")
			}
		})
	}
}

func TestExecuteMultiPageBaseFunction(t *testing.T) {
	// Mock GrpcClient and any other necessary dependencies
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	//mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	FunctionMap = map[string]interface{}{
		"e8": fetchAllBuildedComponents,

		// Add other function mappings as needed
	}

	tests := []struct {
		name         string
		k            string
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          *mock_grpc_client.MockGrpcClient // Use the mock client type
		wantErr      bool
		expected     json.RawMessage
	}{
		{
			name:     "Case 1: Valid function with valid parameters",
			k:        "e8",
			widgetId: "e8",
			replacements: map[string]interface{}{
				constants.ORG_ID: "org124",
				"component":      []string{"testComp"},
			},
			ctx:      context.TODO(),
			clt:      mockGrpcClient,
			wantErr:  false,
			expected: json.RawMessage(`{"key": "value"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			openSearchClient = func() (*opensearch.Client, error) {
				return opensearch.NewDefaultClient()
			}
			responseString := `{
				"aggregations": {
				  "components": {
					"value": {
					  "components": [
						"1358c14d-6d31-4b81-974e-5287e51b85fd",
						"a1d40803-7978-40ed-6c6e-d6fa114cbfa7",
						"238ffe68-8cb4-459d-64ac-2e4f752fe8dc",
						"addf84e3-5c30-49d0-aaa7-1564ee109c0b",
						"24c2c304-75b1-4076-9117-30d5d375ab15"
					  ],
					  "min": "0",
					  "max": "0"
					}
				  }
				}
			  }`
			searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
				return responseString, nil
			}

			result, _, err := ExecuteMultiPageBaseFunction(tt.k, tt.widgetId, tt.replacements, tt.ctx)

			if tt.wantErr && err == nil {
				t.Errorf("Expected an error, but got none.")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, but got an error: %v", err)
			}

			if tt.expected != nil {
				assert.Equal(t, len(result), 5, "Validating response count for automation run drilldown")
			}
		})
	}
}

func Test_componentWidgetHeader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	FunctionMap = map[string]interface{}{
		"e9": FetchAllDeployedEnvironments,

		// Add other function mappings as needed
	}

	tests := []struct {
		name         string
		k            string
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          *mock_grpc_client.MockGrpcClient // Use the mock client type
		wantErr      bool
		expected     json.RawMessage
	}{
		{
			name:     "Case 1: Valid function with valid parameters",
			k:        "e9",
			widgetId: "e9",
			replacements: map[string]interface{}{
				constants.ORG_ID:     "org124",
				constants.SUB_ORG_ID: "sb1",
				"component":          []string{"testComp"},
			},
			ctx:      context.TODO(),
			clt:      mockGrpcClient,
			wantErr:  false,
			expected: json.RawMessage(`{"key": "value"}`),
		},
		{
			name:         "Case 2: No orgId",
			k:            "e9",
			widgetId:     "e9",
			replacements: map[string]interface{}{},
			ctx:          context.TODO(),
			clt:          mockGrpcClient,
			wantErr:      false,
			expected:     json.RawMessage(`{"key": "value"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := componentWidgetHeader(tt.widgetId, tt.replacements, tt.ctx, tt.clt, nil)
			if err != nil {
				assert.Equal(t, len(got), 0, "org Id should not be null")
			} else {
				if tt.expected != nil {
					assert.Equal(t, len(got), 109, "Validating response count for automation run drilldown")
				}
			}

		})
	}
}

func Test_componentWidgetSection(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		widgetId := "widget123"

		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"distinct_component": {
					"value": [
						"7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a"
					]
				}
			}
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := componentWidgetSection(widgetId, testReplacements, ctx, mockGrpcClient, nil)

		assert.Equal(t, len(reports), 391, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		widgetId := "widget123"

		testReplacements := map[string]any{
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"distinct_component": {
					"value": [
						"7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a"
					]
				}
			}
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		_, _ = componentWidgetSection(widgetId, testReplacements, ctx, mockGrpcClient, nil)

	})

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		widgetId := "widget123"

		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"distinct_component": {
					"value": [
						"7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a"
					]
				}
			}
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, errors.New("error")
		}

		_, _ = componentWidgetSection(widgetId, testReplacements, ctx, mockGrpcClient, nil)

	})

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		widgetId := "widget123"

		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		ctx := context.Background()

		var Client *opensearch.Client

		openSearchClient = func() (*opensearch.Client, error) {
			return Client, errors.New("error")
		}
		responseString := `{
			"aggregations": {
				"distinct_component": {
					"value": [
						"7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a"
					]
				}
			}
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		_, _ = componentWidgetSection(widgetId, testReplacements, ctx, mockGrpcClient, nil)

	})
	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		widgetId := "widget123"

		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		responseString := `{
			"aggregations": {
				"distinct_component": {
					"value": [
						"7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a"
					]
				}
			}
		  }`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, errors.New("error")
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		_, _ = componentWidgetSection(widgetId, testReplacements, ctx, mockGrpcClient, nil)

	})

}

func Test_securityComponentWidgetSection(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		widgetId := "widget123"

		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		response := map[string]json.RawMessage{}
		response["scan"] = []byte(`{"took":2149,"timed_out":false,"aggregations":{"distinct_component":{"value":["7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a","94a81dd1-3f52-4520-891e-b2440f660945","c497a177-931e-4ab9-8822-0c4d39d2acfc"]}}}`)
		response["rawScan"] = []byte(`{"took":28989,"timed_out":false,"aggregations":{"distinct_component":{"value":["94a81dd1-3f52-4520-891e-b2440f660945","4c1ce669-7acc-475f-b5ee-b5c826ff5c3c"]}}}`)

		multiSearchResponse = func(map[string]db.DbQuery) (map[string]json.RawMessage, error) {
			return response, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
				{
					Id:            "4c1ce669-7acc-475f-b5ee-b5c826ff5c3c",
					Name:          "dsl-engine",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine.git",
				},
				{
					Id:            "94a81dd1-3f52-4520-891e-b2440f660945",
					Name:          "dsl-engine-cli9",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli9.git",
				},
				{
					Id:            "c497a177-931e-4ab9-8822-0c4d39d2acfc",
					Name:          "dsl-engine7",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine7.git",
				},
				{
					Id:            "6537a7f8-1b1a-4f5d-8413-ca40e9b57893",
					Name:          "dsl-engine-cli5",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli5.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := securityComponentWidgetSection(widgetId, testReplacements, ctx, mockGrpcClient, nil)

		assert.Equal(t, 445, len(reports), "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})
	t.Run("Case 2: Failed response from OS query to Scan Results", func(t *testing.T) {
		widgetId := "widget123"

		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		response := map[string]json.RawMessage{}
		response["scan"] = []byte(`{"took":2149,"timed_out":false}`)
		response["rawScan"] = []byte(`{"took":28989,"timed_out":false,"aggregations":{"distinct_component":{"value":["94a81dd1-3f52-4520-891e-b2440f660945","4c1ce669-7acc-475f-b5ee-b5c826ff5c3c"]}}}`)

		multiSearchResponse = func(map[string]db.DbQuery) (map[string]json.RawMessage, error) {
			return response, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
				{
					Id:            "4c1ce669-7acc-475f-b5ee-b5c826ff5c3c",
					Name:          "dsl-engine",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine.git",
				},
				{
					Id:            "94a81dd1-3f52-4520-891e-b2440f660945",
					Name:          "dsl-engine-cli9",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli9.git",
				},
				{
					Id:            "c497a177-931e-4ab9-8822-0c4d39d2acfc",
					Name:          "dsl-engine7",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine7.git",
				},
				{
					Id:            "6537a7f8-1b1a-4f5d-8413-ca40e9b57893",
					Name:          "dsl-engine-cli5",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli5.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := securityComponentWidgetSection(widgetId, testReplacements, ctx, mockGrpcClient, nil)

		assert.Equal(t, 445, len(reports), "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})
}

func Test_workflowWidgetComponentComparison(t *testing.T) {
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
		name         string
		orgId        string
		component    []string
		replacements map[string]any
		want         json.RawMessage
		wantErr      bool
	}{
		{
			name:      "Case 1: Successful execution with valid parameters",
			orgId:     "org123",
			component: []string{"All"},
			replacements: map[string]any{
				"orgId":     "org123",
				"subOrgId":  "suborg123",
				"component": []string{"All"},
			},
			wantErr: false,
			want:    json.RawMessage(`{"key": "value"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automationMap := map[string]struct{}{
				"auto1": {},
			}

			openSearchClient = func() (*opensearch.Client, error) {
				return opensearch.NewDefaultClient()
			}
			responseString := `{
                "aggregations": {
                    "distinct_automation": {
                        "value": [
                          "2bb25065-e310-4552-91b8-6d4fdf6ac429",
                          "663a7894-7907-4b97-9ff6-b96c8dad9b2f",
                          "99a8e78a-5a3f-4675-85f0-22eab67c50a0"
                          ]
                    }
                }
              }`
			searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
				return responseString, nil
			}

			GetAutomationMap = func(ctx context.Context, clt client.GrpcClient, orgId string, component []string, branch string) map[string]struct{} {
				return automationMap
			}

			got, err := workflowWidgetComponentComparison(tt.orgId, tt.replacements, context.Background(), mockGrpcClient, nil)

			if tt.wantErr && err == nil {
				t.Errorf("Expected an error, but got none.")
			}

			if err != nil {
				assert.Equal(t, len(got), 61, "Validating response count for automation run drilldown")
			}
		})
	}
}

func Test_testInsightsWorkflowsComponentComparison(t *testing.T) {
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
		name         string
		orgId        string
		component    []string
		replacements map[string]any
		want         json.RawMessage
		wantErr      bool
	}{
		{
			name:      "Case 1: Successful execution with valid parameters",
			orgId:     "org123",
			component: []string{"All"},
			replacements: map[string]any{
				"orgId":     "org123",
				"subOrgId":  "suborg123",
				"component": []string{"All"},
			},
			wantErr: false,
			want:    json.RawMessage(`{"key": "value"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automationMap := map[string]struct{}{
				"auto1": {},
			}

			openSearchClient = func() (*opensearch.Client, error) {
				return opensearch.NewDefaultClient()
			}
			responseString := `{"aggregations":{"test_insights_automations":{"value":["e7f77f79-2b61-4564-a43a-4b5974e03b91","1ed52014-45a6-4bc4-aa78-ce1cdb072643","f00f4443-01d0-4b2f-b75c-19d76c02e0ec"]}}}`
			searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
				return responseString, nil
			}

			GetAutomationMap = func(ctx context.Context, clt client.GrpcClient, orgId string, component []string, branch string) map[string]struct{} {
				return automationMap
			}

			got, err := testInsightsWorkflowsComponentComparison(tt.orgId, tt.replacements, context.Background(), mockGrpcClient, nil)
			if tt.wantErr && err == nil {
				t.Errorf("Expected an error, but got none.")
			}

			if err != nil {
				assert.Equal(t, len(got), 61, "Validating response count for automation run drilldown")
			}
		})
	}
}

func Test_workflowComponentComparisonSI(t *testing.T) {
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
		name         string
		orgId        string
		component    []string
		replacements map[string]any
		want         json.RawMessage
		wantErr      bool
	}{
		{
			name:      "Case 1: Successful execution with valid parameters",
			orgId:     "org123",
			component: []string{"All"},
			replacements: map[string]any{
				"orgId":     "org123",
				"subOrgId":  "suborg123",
				"component": []string{"All"},
			},
			wantErr: false,
			want:    json.RawMessage(`{"key": "value"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automationMap := map[string]struct{}{
				"auto1": {},
			}

			openSearchClient = func() (*opensearch.Client, error) {
				return opensearch.NewDefaultClient()
			}
			responseString := `{
                "aggregations": {
                    "distinct_automation": {
                        "value": [
                          "2bb25065-e310-4552-91b8-6d4fdf6ac429",
                          "663a7894-7907-4b97-9ff6-b96c8dad9b2f",
                          "99a8e78a-5a3f-4675-85f0-22eab67c50a0"
                          ]
                    }
                }
              }`
			searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
				return responseString, nil
			}

			GetAutomationMap = func(ctx context.Context, clt client.GrpcClient, orgId string, component []string, branch string) map[string]struct{} {
				return automationMap
			}

			got, err := workflowComponentComparisonSI(tt.orgId, tt.replacements, context.Background(), mockGrpcClient, nil)

			if tt.wantErr && err == nil {
				t.Errorf("Expected an error, but got none.")
			}

			if err != nil {
				assert.Equal(t, len(got), 61, "Validating response count for automation run drilldown")
			}
		})
	}
}

func TestUpdateMustParentFilters(t *testing.T) {

	updatedJSON := `{
		"query": {
			"bool": {
				"must": []
			}
		}
	}`
	replacements := map[string]interface{}{
		"parentIds": []string{"parent1", "parent2"},
	}

	result := UpdateMustParentFilters(updatedJSON, replacements)

	var resultData map[string]interface{}
	err := json.Unmarshal([]byte(result), &resultData)
	if err != nil {
		t.Fatalf("Failed to unmarshal result JSON: %v", err)
	}
	assert.NotNil(t, resultData["query"], "Query field should not be nil")
	query, ok := resultData["query"].(map[string]interface{})
	assert.True(t, ok, "Query field should be a map")
	boolObj, ok := query["bool"].(map[string]interface{})
	assert.True(t, ok, "Bool field should be a map")
	filterArray, ok := boolObj["filter"].([]interface{})
	assert.True(t, ok, "Filter field should be a slice")
	assert.Equal(t, 1, len(filterArray), "There should be one filter added")

	filter := filterArray[0]
	if filterMap, ok := filter.(map[string]interface{}); ok {
		if terms, ok := filterMap["terms"].(map[string]interface{}); ok {
			orgID := terms["org_id"]
			if orgIDSlice, ok := orgID.([]interface{}); ok {
				var orgIDStrings []string
				for _, v := range orgIDSlice {
					if s, ok := v.(string); ok {
						orgIDStrings = append(orgIDStrings, s)
					} else {
						t.Errorf("unexpected type in org_id: %T", v)
					}
				}
				assert.Equal(t, []string{"parent1", "parent2"}, orgIDStrings, "org_id should contain parent1 and parent2")
			} else {
				t.Errorf("org_id field is not a []interface{}")
			}
		} else {
			t.Errorf("terms field is not a map[string]interface{}")
		}
	} else {
		t.Errorf("filter is not a map[string]interface{}")
	}
}

func Test_securityComponentComponentComparison(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		widgetId := "widget123"

		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		response := map[string]json.RawMessage{}
		response["scan"] = []byte(`{"took":2149,"timed_out":false,"aggregations":{"distinct_component":{"value":["7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a","94a81dd1-3f52-4520-891e-b2440f660945","c497a177-931e-4ab9-8822-0c4d39d2acfc"]}}}`)
		response["rawScan"] = []byte(`{"took":28989,"timed_out":false,"aggregations":{"distinct_component":{"value":["94a81dd1-3f52-4520-891e-b2440f660945","4c1ce669-7acc-475f-b5ee-b5c826ff5c3c"]}}}`)

		multiSearchResponse = func(map[string]db.DbQuery) (map[string]json.RawMessage, error) {
			return response, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
				{
					Id:            "4c1ce669-7acc-475f-b5ee-b5c826ff5c3c",
					Name:          "dsl-engine",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine.git",
				},
				{
					Id:            "94a81dd1-3f52-4520-891e-b2440f660945",
					Name:          "dsl-engine-cli9",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli9.git",
				},
				{
					Id:            "c497a177-931e-4ab9-8822-0c4d39d2acfc",
					Name:          "dsl-engine7",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine7.git",
				},
				{
					Id:            "6537a7f8-1b1a-4f5d-8413-ca40e9b57893",
					Name:          "dsl-engine-cli5",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli5.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := securityComponentComponentComparison(widgetId, testReplacements, ctx, mockGrpcClient, nil)

		assert.Equal(t, len(reports), 157, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})
	t.Run("Case 2: Failed response from OS query to Scan Results", func(t *testing.T) {
		widgetId := "widget123"

		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}
		response := map[string]json.RawMessage{}
		response["scan"] = []byte(`{"took":2149,"timed_out":false}`)
		response["rawScan"] = []byte(`{"took":28989,"timed_out":false,"aggregations":{"distinct_component":{"value":["94a81dd1-3f52-4520-891e-b2440f660945","4c1ce669-7acc-475f-b5ee-b5c826ff5c3c"]}}}`)

		multiSearchResponse = func(map[string]db.DbQuery) (map[string]json.RawMessage, error) {
			return response, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
				{
					Id:            "4c1ce669-7acc-475f-b5ee-b5c826ff5c3c",
					Name:          "dsl-engine",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine.git",
				},
				{
					Id:            "94a81dd1-3f52-4520-891e-b2440f660945",
					Name:          "dsl-engine-cli9",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli9.git",
				},
				{
					Id:            "c497a177-931e-4ab9-8822-0c4d39d2acfc",
					Name:          "dsl-engine7",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine7.git",
				},
				{
					Id:            "6537a7f8-1b1a-4f5d-8413-ca40e9b57893",
					Name:          "dsl-engine-cli5",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli5.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := securityComponentComponentComparison(widgetId, testReplacements, ctx, mockGrpcClient, nil)

		assert.Equal(t, len(reports), 79, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})
}

func Test_getComponentWidgetData(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {

		OrgId := "org123"

		Component := []string{"All"}

		ctx := context.Background()

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := getComponentWidgetData(ctx, mockGrpcClient, OrgId, Component)

		assert.Equal(t, len(reports), 109, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

}

func Test_getAutomationsForBranch(t *testing.T) {

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

	mockGrpcClient.EXPECT().RetriableGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	t.Run("Case 1: Successful execution with components", func(t *testing.T) {

		Branch := "1"

		reports := GetAutomationsForBranch(Branch)

		assert.Equal(t, len(reports), 1, "Validating response count for automation run drilldown")

	})

}

func Test_getBranchNameForId(t *testing.T) {
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

	mockGrpcClient.EXPECT().RetriableGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	t.Run("Case 1: Successful execution with components", func(t *testing.T) {

		Branch := "1"

		reports := GetBranchNameForId(Branch)

		assert.Equal(t, len(reports), 2, "Validating response count for automation run drilldown")

	})

}

func Test_getAutomationMap(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockCache := coredataMock.NewMockResourceCacheI(mockCtrl)
	cache.CoreDataResourceCache = mockCache
	mockCache.EXPECT().GetChildrenOfType(gomock.Any(), gomock.Any()).Return([]string{"1"}).Times(2).AnyTimes()
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
	t.Run("Case 1: Successful execution with components", func(t *testing.T) {

		OrgId := "org123"

		Component := []string{"7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a"}

		Branch := "1"

		ctx := context.Background()

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
			},
		}

		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports := getAutomationMap(ctx, mockGrpcClient, OrgId, Component, Branch)

		assert.Equal(t, len(reports), 1, "Validating response count for automation run drilldown")

	})

}

func Test_getAutomationMapComponentComparison(t *testing.T) {

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

	mockGrpcClient.EXPECT().RetriableGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	t.Run("Case 1: Successful execution with components", func(t *testing.T) {

		OrgId := "org123"

		ctx := context.Background()

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
			},
		}
		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := getAutomationMapComponentComparison(ctx, mockGrpcClient, OrgId)

		assert.Equal(t, len(reports), 1, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

}

func Test_getAutomationWidgetDataForBranch(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockCache := coredataMock.NewMockResourceCacheI(mockCtrl)
	cache.CoreDataResourceCache = mockCache
	mockCache.EXPECT().GetChildrenOfType(gomock.Any(), gomock.Any()).Return([]string{"1"}).Times(2).AnyTimes()

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
	t.Run("Case 1: Successful execution with components", func(t *testing.T) {

		OrgId := "org123"

		Component := []string{"7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a"}

		Branch := "Id1"

		ctx := context.Background()

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
			},
		}

		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := getAutomationWidgetDataForBranch(ctx, mockGrpcClient, OrgId, Component, Branch)

		assert.Equal(t, len(reports), 11, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

}

func Test_getAutomationWidgetData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockCache := coredataMock.NewMockResourceCacheI(mockCtrl)
	cache.CoreDataResourceCache = mockCache

	mockCache.EXPECT().GetChildrenOfType(gomock.Any(), gomock.Any()).Return([]string{"1"}).Times(2).AnyTimes()
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

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {

		OrgId := "org123"

		Component := []string{"7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a"}

		ctx := context.Background()

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "7cb5a7f8-1b1a-4f5d-8413-ca40e9b5781a",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
			},
		}

		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := getAutomationWidgetData(ctx, mockGrpcClient, OrgId, Component)

		assert.Equal(t, len(reports), 105, "Validating response count for automation run drilldown")
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

}

func Test_automationWidgetHeader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {

		widgetId := "widget123"
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		_, _ = automationWidgetHeader(widgetId, replacements, context.Background(), mockGrpcClient, nil)

	})

}

func Test_automationWidgetSection(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	//mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	tests := []struct {
		name         string
		orgId        string
		component    []string
		replacements map[string]any
		want         json.RawMessage
		wantErr      bool
	}{
		{
			name:      "Case 1: Successful execution with valid parameters",
			orgId:     "org123",
			component: []string{"All"},
			replacements: map[string]any{
				"orgId":     "org123",
				"subOrgId":  "suborg123",
				"component": []string{"All"},
			},
			wantErr: false,
			// Define the expected JSON response here
			want: json.RawMessage(`{"key": "value"}`),
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automationMap := map[string]struct{}{
				"auto1": {},
			}

			openSearchClient = func() (*opensearch.Client, error) {
				return opensearch.NewDefaultClient()
			}
			responseString := `{
				"aggregations": {
					"distinct_automation": {
						"value": [
						  "2bb25065-e310-4552-91b8-6d4fdf6ac429",
						  "663a7894-7907-4b97-9ff6-b96c8dad9b2f",
						  "99a8e78a-5a3f-4675-85f0-22eab67c50a0"
						  ]
					}
				}
			  }`
			searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
				return responseString, nil
			}

			GetAutomationMap = func(ctx context.Context, clt client.GrpcClient, orgId string, component []string, branch string) map[string]struct{} {
				return automationMap
			}

			got, err := automationWidgetSection(tt.orgId, tt.replacements, context.Background(), mockGrpcClient, nil)

			if tt.wantErr && err == nil {
				t.Errorf("Expected an error, but got none.")
			}

			if !tt.wantErr && err != nil {
				assert.Equal(t, len(got), 102, "Validating response count for automation run drilldown")
			}

			// Add assertions for comparing the result with tt.want
		})
	}
}

func Test_summaryAutomationWidgetSection(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	//mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	tests := []struct {
		name         string
		orgId        string
		component    []string
		replacements map[string]any
		want         json.RawMessage
		wantErr      bool
	}{
		{
			name:      "Case 1: Successful execution with valid parameters",
			orgId:     "org123",
			component: []string{"All"},
			replacements: map[string]any{
				"orgId":     "org123",
				"subOrgId":  "suborg123",
				"component": []string{"All"},
			},
			wantErr: false,
			// Define the expected JSON response here
			want: json.RawMessage(`{"key": "value"}`),
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automationMap := map[string]struct{}{
				"auto1": {},
			}

			openSearchClient = func() (*opensearch.Client, error) {
				return opensearch.NewDefaultClient()
			}
			responseString := `{
				"aggregations": {
					"distinct_automation": {
						"value": [
						  "2bb25065-e310-4552-91b8-6d4fdf6ac429",
						  "663a7894-7907-4b97-9ff6-b96c8dad9b2f",
						  "99a8e78a-5a3f-4675-85f0-22eab67c50a0"
						  ]
					}
				}
			  }`
			searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
				return responseString, nil
			}

			GetAutomationMap = func(ctx context.Context, clt client.GrpcClient, orgId string, component []string, branch string) map[string]struct{} {
				return automationMap
			}

			got, err := summaryAutomationWidgetSection(tt.orgId, tt.replacements, context.Background(), mockGrpcClient, nil)

			if tt.wantErr && err == nil {
				t.Errorf("Expected an error, but got none.")
			}

			if err != nil {
				assert.Equal(t, len(got), 61, "Validating response count for automation run drilldown")
			}

			// Add assertions for comparing the result with tt.want
		})
	}
}

func Test_securityAutomationWidgetSection(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	//mockGrpcClient.EXPECT().SendGrpcCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	tests := []struct {
		name         string
		orgId        string
		component    []string
		replacements map[string]any
		want         json.RawMessage
		wantErr      bool
	}{
		{
			name:      "Case 1: Successful execution with valid parameters",
			orgId:     "org123",
			component: []string{"All"},
			replacements: map[string]any{
				"orgId":     "org123",
				"subOrgId":  "suborg123",
				"component": []string{"All"},
			},
			wantErr: false,
			// Define the expected JSON response here
			want: json.RawMessage(`{"key": "value"}`),
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automationMap := map[string]struct{}{
				"auto1": {},
			}

			openSearchClient = func() (*opensearch.Client, error) {
				return opensearch.NewDefaultClient()
			}
			responseString := `{
				"aggregations": {
					"distinct_automation": {
						"value": [
							"036d13c2-3009-4dae-8f09-2d86cc066ee8",
							"b772fc94-3613-484e-9236-9f53f28b7503",
							"f1ccd216-8283-41e9-9187-925b2535bf83"
						  ]
					}
				}
			  }`
			searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
				return responseString, nil
			}

			GetAutomationMap = func(ctx context.Context, clt client.GrpcClient, orgId string, component []string, branch string) map[string]struct{} {
				return automationMap
			}

			got, err := securityAutomationWidgetSection(tt.orgId, tt.replacements, context.Background(), mockGrpcClient, nil)

			if tt.wantErr && err == nil {
				t.Errorf("Expected an error, but got none.")
			}

			if err != nil {
				assert.Equal(t, len(got), 447, "Validating response count for automation run drilldown")
			}

			// Add assertions for comparing the result with tt.want
		})
	}
}

func Test_getCodeBaseOverviewWidget(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	var Client client.GrpcClient

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		response, _ := getCodeBaseOverviewWidget(widget, replacements, context.Background(), Client, nil)
		mapArray := []map[string]any{}
		json.Unmarshal(response, &mapArray)
		assert.Equal(t, len(mapArray), 3, "Success validation")
	})
	t.Run("Empty Sonar array", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return "", db.ErrInternalServer
		}
		response, err := getCodeBaseOverviewWidget(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, json.RawMessage(json.RawMessage(nil)), "Empty Sonar array")
		assert.Equal(t, err, nil, "Empty Sonar array")
	})

}

func Test_getIssueTypeHeader1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	var Client client.GrpcClient

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
		return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
	}

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		widget := "s10"
		responseStr := `{"subHeader":[{"title":"Code Smell","value":1,"drillDown":{"reportId":"","reportTitle":"","reportType":""}},{"title":"Bug","value":1,"drillDown":{"reportId":"","reportTitle":"","reportType":""}}]}`
		var out json.RawMessage = []byte(responseStr)
		response, _ := getIssueTypeHeader1(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, out, "Success validation")
	})

	t.Run("Empty Sonar array", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		widget := "s10"

		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return "", db.ErrInternalServer
		}
		response, err := getIssueTypeHeader1(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, json.RawMessage(json.RawMessage(nil)), "Empty Sonar array")
		assert.Equal(t, err, nil, "Empty Sonar array")
	})
}

func Test_getIssueTypeHeader2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	var Client client.GrpcClient

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		responseStr := `{"subHeader":[{"title":"Vulnerabilities","value":1,"drillDown":{"reportId":"","reportTitle":"","reportType":""}},{"title":"Security hotspots","value":1,"drillDown":{"reportId":"","reportTitle":"","reportType":""}}]}`
		var out json.RawMessage = []byte(responseStr)
		response, _ := getIssueTypeHeader2(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, out, "Success validation")
	})

	t.Run("Empty Sonar array", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return "", db.ErrInternalServer
		}
		response, err := getIssueTypeHeader2(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, json.RawMessage(json.RawMessage(nil)), "Empty Sonar array")
		assert.Equal(t, err, nil, "Empty Sonar array")

	})

}

func Test_getIssueTypeSection(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	var Client client.GrpcClient
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	t.Run("Success with single value", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		response, _ := getIssueTypeSection(widget, replacements, context.Background(), Client, nil)
		mapArray := []map[string]any{}
		json.Unmarshal(response, &mapArray)
		assert.Equal(t, len(mapArray), 4, "Success validation")
	})

	t.Run("Success with multiple value", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]},{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		response, _ := getIssueTypeSection(widget, replacements, context.Background(), Client, nil)
		mapArray := []map[string]any{}
		json.Unmarshal(response, &mapArray)
		assert.Equal(t, len(mapArray), 4, "Success validation")
	})

	t.Run("Empty Sonar array", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return "", db.ErrInternalServer
		}
		response, err := getIssueTypeSection(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, json.RawMessage(json.RawMessage(nil)), "Empty Sonar array")
		assert.Equal(t, err, nil, "Empty Sonar array")
	})
}

func Test_getCodeBaseOverview(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		response := &scanner.Report{

			CommitSha:  "c4d3f8a",
			BranchName: "main",
			Source:     "cloudbees/SonarQubeAction",
			Size: &scanner.Size{
				Lines:        1500,
				CommentLines: 300,
			},
			Complexity: &scanner.Complexity{
				Cyclomatic: 20,
			},
			Duplication: &scanner.Duplication{
				Lines:  100,
				Blocks: 5,
			}, Files: []*scanner.FileData{
				{
					File: "main.go",
					Coverage: &scanner.Coverage{
						LinesToCover: 100,
						CoveredLines: 300,
					},
					Size: &scanner.Size{
						Lines: 1000,
					},
					Complexity: &scanner.Complexity{
						Cyclomatic: 10,
					},
				},
			},
		}

		_ = getCodeBaseOverview(response)

	})

}

func Test_IsSonarWidgetsApplicable(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		_ = IsSonarWidgetsApplicable(replacements, context.Background(), constants.SECURITY_INDEX, constants.RAW_SCAN_RESULTS_INDEX)

	})

}

func Test_IsSecurityWidgetsApplicable(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		_ = IsSecurityWidgetsApplicable(replacements, context.Background(), constants.SECURITY_INDEX, constants.RAW_SCAN_RESULTS_INDEX)

	})

}

func Test_getIssueTypeCount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		response := &scanner.Report{

			CommitSha:  "c4d3f8a",
			BranchName: "main",
			Source:     "cloudbees/SonarQubeAction",
			Size: &scanner.Size{
				Lines:        1500,
				CommentLines: 300,
			},
			Complexity: &scanner.Complexity{
				Cyclomatic: 20,
			},
			Duplication: &scanner.Duplication{
				Lines:  100,
				Blocks: 5,
			}, Files: []*scanner.FileData{
				{
					File: "main.go",
					Coverage: &scanner.Coverage{
						LinesToCover: 100,
						CoveredLines: 300,
					},
					Size: &scanner.Size{
						Lines: 1000,
					},
					Complexity: &scanner.Complexity{
						Cyclomatic: 10,
					},
					Issues: []*scanner.Issue{
						{
							Code: "code1",
						},
					},
				},
			},
		}

		_ = getIssueTypeCount(response)

	})

	t.Run("Success with prefix", func(t *testing.T) {

		response := &scanner.Report{

			CommitSha:  "c4d3f8a",
			BranchName: "main",
			Source:     "cloudbees/SonarQubeAction",
			Size: &scanner.Size{
				Lines:        1500,
				CommentLines: 300,
			},
			Complexity: &scanner.Complexity{
				Cyclomatic: 20,
			},
			Duplication: &scanner.Duplication{
				Lines:  100,
				Blocks: 5,
			}, Files: []*scanner.FileData{
				{
					File: "main.go",
					Coverage: &scanner.Coverage{
						LinesToCover: 100,
						CoveredLines: 300,
					},
					Size: &scanner.Size{
						Lines: 1000,
					},
					Complexity: &scanner.Complexity{
						Cyclomatic: 10,
					},
					Issues: []*scanner.Issue{
						{
							Code: "CODE_SMELL",
						}, {
							Code: "BUG",
						}, {
							Code: "SECURITY_HOTSPOT",
						}, {
							Code: "VULNERABILITY",
						},
					},
				},
			},
		}

		_ = getIssueTypeCount(response)

	})

}

func Test_getCoverageDataHeader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	var Client client.GrpcClient
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]},{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		responseStr := `{"value":"0%"}`
		var out json.RawMessage = []byte(responseStr)
		response, _ := getCoverageDataHeader(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, out, "Success validation")
	})

	t.Run("Empty Sonar array", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return "", db.ErrInternalServer
		}
		response, err := getCoverageDataHeader(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, json.RawMessage(json.RawMessage(nil)), "Empty Sonar array")
		assert.Equal(t, err, nil, "Empty Sonar array")
	})
}

func Test_getCoverageDataSection1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	var Client client.GrpcClient

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897,"coverage_pct":20},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]},{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897,"coverage_pct":20},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		response, _ := getCoverageDataSection1(widget, replacements, context.Background(), Client, nil)
		mapArray := []map[string]any{}
		json.Unmarshal(response, &mapArray)
		assert.Equal(t, len(mapArray), 1, "Success validation")
	})
	t.Run("Empty Sonar array", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return "", db.ErrInternalServer
		}
		response, err := getCoverageDataSection1(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, json.RawMessage(json.RawMessage(nil)), "Empty Sonar array")
		assert.Equal(t, err, nil, "Empty Sonar array")
	})
}

func Test_getCoverageDataSection2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	var Client client.GrpcClient
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	t.Run("Success", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]},{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		responseStr := `[{"title":"Total lines","value":13783},{"title":"Total code lines","value":12688},{"title":"Lines covered","value":0},{"title":"Lines to cover","value":2897}]`
		var out json.RawMessage = []byte(responseStr)
		response, _ := getCoverageDataSection2(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, out, "Success validation")

	})

	t.Run("Empty Sonar array", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return "", db.ErrInternalServer
		}
		response, err := getCoverageDataSection2(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, json.RawMessage(json.RawMessage(nil)), "Empty Sonar array")
		assert.Equal(t, err, nil, "Empty Sonar array")
	})
}

func Test_getDuplicationDataSection2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	var Client client.GrpcClient
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	t.Run("Success", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]},{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		responseStr := `[{"title":"Files with duplication","value":23},{"title":"Duplicate blocks","value":50},{"title":"Total lines","value":13783},{"title":"Duplicate lines","value":13783}]`
		var out json.RawMessage = []byte(responseStr)
		response, _ := getDuplicationDataSection2(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, out, "Success validation")
	})

	t.Run("Empty Sonar array", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return "", db.ErrInternalServer
		}
		response, err := getDuplicationDataSection2(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, json.RawMessage(json.RawMessage(nil)), "Empty Sonar array")
		assert.Equal(t, err, nil, "Empty Sonar array")
	})
}

func Test_getDuplicationHeader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)
	var Client client.GrpcClient
	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	t.Run("Success", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]},{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		responseStr := `{"value":"108%"}`
		var out json.RawMessage = []byte(responseStr)
		response, _ := getDuplicationHeader(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, out, "Success validation")

	})

	t.Run("Empty Sonar array", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return "", db.ErrInternalServer
		}
		response, err := getDuplicationHeader(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, json.RawMessage(json.RawMessage(nil)), "Empty Sonar array")
		assert.Equal(t, err, nil, "Empty Sonar array")

	})
}
func Test_getDuplicationDataSection1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	var Client client.GrpcClient

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	t.Run("Success with single value", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		responseStr := `[{"data":[{"x":"Current scan","y":108}],"id":"Duplicated line density"}]`
		var out json.RawMessage = []byte(responseStr)
		response, _ := getDuplicationDataSection1(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, out, "Success validation")
	})

	t.Run("Success with multiple value", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return `{"took":7,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":304,"relation":"eq"},"max_score":null,"hits":[{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]},{"_index":"raw_scan_result","_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002_5415ddd1-2e4c-4103-a7f6-7d8235097e88_54d882eb-06ad-4528-9b7f-7a5db6b7bd7a_1698593634070","_score":null,"_source":{"component_id":"54d882eb-06ad-4528-9b7f-7a5db6b7bd7a","run_id":"5415ddd1-2e4c-4103-a7f6-7d8235097e88","action_raw_result":{"coverage":{"lines_to_cover":2897},"complexity":{"cognitive":2352,"cyclomatic":3367},"commit_sha":"0e70d2d4cc51e957bb2a831196cb8591b310be86","size":{"code_lines":12688,"comment_lines":2375,"functions":2239,"classes":2157,"files":23,"lines":13783,"comment_line_density":17.231373},"files":[{"coverage":{"lines_to_cover":87},"complexity":{"cognitive":110,"cyclomatic":41},"file":"internal/widgetBuilder.go","size":{"code_lines":188,"comment_lines":7,"functions":7,"classes":2,"lines":226,"comment_line_density":3.5897436},"issues":[{"severity":4,"code":"CODE_SMELL-go:S1192","domain":"go"},{"severity":4,"code":"VULNERABILITY-go:S3776","domain":"go"},{"severity":4,"code":"SECURITY_HOTSPOT-go:S3776","domain":"go"},{"severity":4,"code":"BUG-go:S3776","domain":"go"}]},{"coverage":{"lines_to_cover":72},"complexity":{"cognitive":30,"cyclomatic":26},"file":"internal/sca-report.go","size":{"code_lines":122,"comment_lines":17,"functions":4,"classes":4,"lines":172,"comment_line_density":12.230216},"duplicate_code":[{"duplicate_blocks":[{"file":"internal/widget.go","text_range":{"start_line":293,"end_line":326}}],"origin_block":{"start_line":122,"end_line":155}}]},{"coverage":{},"complexity":{},"file":"coverage.xml","size":{"code_lines":4115,"lines":4116}}],"source":"SonarQube CLI","duplication":{"line_density":108.6302,"blocks":50,"files":23,"lines":13783}},"component_name":"reports-service","github_branch_id":"9072d709-2d43-4bd3-9627-bbc52b523396","step_id":"s010-1-build-3-0-sonarqube","job_id":"build-deploy-staging","org_id":"2cab10cc-cd9d-11ed-afa1-0242ac120002","github_branch":"main","org_name":"cloudbees-staging","timestamp":"2023-10-29 15:33:54"},"sort":[1698593634000]}]}}`, nil
		}
		response, _ := getDuplicationDataSection1(widget, replacements, context.Background(), Client, nil)
		mapArray := []map[string]any{}
		json.Unmarshal(response, &mapArray)
		dataMap := mapArray[0]
		assert.Equal(t, len(dataMap["data"].([]any)), 2, "Success validation")
	})

	t.Run("Empty Sonar array", func(t *testing.T) {
		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}
		widget := "s10"
		searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
			return "", db.ErrInternalServer
		}
		response, err := getDuplicationDataSection1(widget, replacements, context.Background(), Client, nil)
		assert.Equal(t, response, json.RawMessage(json.RawMessage(nil)), "Empty Sonar array")
		assert.Equal(t, err, nil, "Empty Sonar array")
	})
}

func Test_getCoveragePercetange(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		response := &scanner.Report{

			CommitSha:  "c4d3f8a",
			BranchName: "main",
			Source:     "cloudbees/SonarQubeAction",
			Size: &scanner.Size{
				Lines:        1500,
				CommentLines: 300,
				CodeLines:    5,
			},
			Coverage: &scanner.Coverage{
				LinesToCover: 10,
			},
		}

		_ = getCoveragePercetange(response)

	})

	t.Run("Empty Sonar array", func(t *testing.T) {

		response := &scanner.Report{}

		_ = getCoveragePercetange(response)
	})

}

func Test_getDuplicationDensity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		response := &scanner.Report{

			CommitSha:  "c4d3f8a",
			BranchName: "main",
			Source:     "cloudbees/SonarQubeAction",
			Size: &scanner.Size{
				Lines:        1500,
				CommentLines: 300,
				CodeLines:    5,
			},
			Duplication: &scanner.Duplication{
				Lines: 10,
			},
		}

		_ = getDuplicationDensity(response)

	})

	t.Run("Empty Sonar array", func(t *testing.T) {

		response := &scanner.Report{}

		_ = getDuplicationDensity(response)
	})

}

func Test_getCoverageInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		response := &scanner.Report{

			CommitSha:  "c4d3f8a",
			BranchName: "main",
			Source:     "cloudbees/SonarQubeAction",
			Size: &scanner.Size{
				Lines:        1500,
				CommentLines: 300,
				CodeLines:    5,
			},
			Coverage: &scanner.Coverage{
				LinesToCover: 10,
			},
		}

		_ = getCoverageInfo(response)

	})

	t.Run("Empty Sonar array", func(t *testing.T) {

		response := &scanner.Report{}

		_ = getCoverageInfo(response)
	})

}

func Test_getDuplicateInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		response := &scanner.Report{

			CommitSha:  "c4d3f8a",
			BranchName: "main",
			Source:     "cloudbees/SonarQubeAction",
			Size: &scanner.Size{
				Lines:        1500,
				CommentLines: 300,
				CodeLines:    5,
			},
			Duplication: &scanner.Duplication{
				Lines: 10,
			},
		}

		_ = getDuplicateInfo(response)

	})

	t.Run("Empty Sonar array", func(t *testing.T) {

		response := &scanner.Report{}

		_ = getDuplicateInfo(response)
	})

}

func Test_getSonarQubeResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		responseString := `{
			"hits": {
				"hits": [
				  {
					"_source": {
					  "component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
					  "action_raw_result": {
						"coverage": {
						  "lines_to_cover": 2034
						},
						"complexity": {
						  "cognitive": 1549,
						  "cyclomatic": 2169
						},
						"size": {
						  "code_lines": 6258,
						  "comment_lines": 1905,
						  "functions": 1805,
						  "classes": 1392,
						  "files": 28,
						  "lines": 7446,
						  "comment_line_density": 25.584206
						}
					 }
					}


				  }

				]
			  }
		  }`

		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		_, _ = getSonarQubeResponse(replacements, context.Background())

	})

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		var client *opensearch.Client
		openSearchClient = func() (*opensearch.Client, error) {
			return client, errors.New("error")
		}

		responseString := `{
			"hits": {
				"hits": [
				  {
					"_source": {
					  "component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
					  "action_raw_result": {
						"coverage": {
						  "lines_to_cover": 2034
						},
						"complexity": {
						  "cognitive": 1549,
						  "cyclomatic": 2169
						},
						"size": {
						  "code_lines": 6258,
						  "comment_lines": 1905,
						  "functions": 1805,
						  "classes": 1392,
						  "files": 28,
						  "lines": 7446,
						  "comment_line_density": 25.584206
						}
					 }
					}


				  }

				]
			  }
		  }`

		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		_, _ = getSonarQubeResponse(replacements, context.Background())

	})

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		responseString := `{
			"hits": {
				"hits": [
				  {
					"_source": {
					  "component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
					  "action_raw_result": {
						"coverage": {
						  "lines_to_cover": 2034
						},
						"complexity": {
						  "cognitive": 1549,
						  "cyclomatic": 2169
						},
						"size": {
						  "code_lines": 6258,
						  "comment_lines": 1905,
						  "functions": 1805,
						  "classes": 1392,
						  "files": 28,
						  "lines": 7446,
						  "comment_line_density": 25.584206
						}
					 }
					}


				  }

				]
			  }
		  }`

		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, errors.New("error")
		}

		_, _ = getSonarQubeResponse(replacements, context.Background())

	})

}

func Test_isSonarOrSecurityDataFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		responseString := `{
			"hits": {
				"hits": [
				  {
					"_source": {
					  "component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
					  "action_raw_result": {
						"coverage": {
						  "lines_to_cover": 2034
						},
						"complexity": {
						  "cognitive": 1549,
						  "cyclomatic": 2169
						},
						"size": {
						  "code_lines": 6258,
						  "comment_lines": 1905,
						  "functions": 1805,
						  "classes": 1392,
						  "files": 28,
						  "lines": 7446,
						  "comment_line_density": 25.584206
						}
					 }
					}


				  }

				]
			  }
		  }`

		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		_, _ = isSonarOrSecurityDataFound(replacements, context.Background(), true, constants.SECURITY_INDEX, constants.RAW_SCAN_RESULTS_INDEX)

	})
	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		var client *opensearch.Client
		openSearchClient = func() (*opensearch.Client, error) {
			return client, errors.New("error")
		}

		responseString := `{
			"hits": {
				"hits": [
				  {
					"_source": {
					  "component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
					  "action_raw_result": {
						"coverage": {
						  "lines_to_cover": 2034
						},
						"complexity": {
						  "cognitive": 1549,
						  "cyclomatic": 2169
						},
						"size": {
						  "code_lines": 6258,
						  "comment_lines": 1905,
						  "functions": 1805,
						  "classes": 1392,
						  "files": 28,
						  "lines": 7446,
						  "comment_line_density": 25.584206
						}
					 }
					}


				  }

				]
			  }
		  }`

		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		_, _ = isSonarOrSecurityDataFound(replacements, context.Background(), true, constants.SECURITY_INDEX, constants.RAW_SCAN_RESULTS_INDEX)

	})

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		responseString := `{
			"hits": {
				"hits": [
				  {
					"_source": {
					  "component_id": "e2f8fef6-5041-4843-b37e-6cdae38099bc",
					  "action_raw_result": {
						"coverage": {
						  "lines_to_cover": 2034
						},
						"complexity": {
						  "cognitive": 1549,
						  "cyclomatic": 2169
						},
						"size": {
						  "code_lines": 6258,
						  "comment_lines": 1905,
						  "functions": 1805,
						  "classes": 1392,
						  "files": 28,
						  "lines": 7446,
						  "comment_line_density": 25.584206
						}
					 }
					}


				  }

				]
			  }
		  }`

		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, errors.New("error")
		}

		_, _ = isSonarOrSecurityDataFound(replacements, context.Background(), true, constants.SECURITY_INDEX, constants.RAW_SCAN_RESULTS_INDEX)

	})

}

func Test_UpdateMustNotFilters(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Success with Job IDs", func(t *testing.T) {
		replacements := map[string]any{
			"jobIds": []string{"job1", "job2"},
		}

		updatedJSON := `{"query":{"bool":{"must_not":[]}}}`

		result := UpdateMustNotFilters(updatedJSON, replacements)

		expectedResult := `{
			"query": {
				"bool": {
					"must_not": [
						{"terms": {"job_id": ["job1", "job2"]}}
					]
				}
			}
		}`
		assert.JSONEq(t, expectedResult, result)
	})

}

func Test_UpdateFiltersForDrilldown(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
			"branch":    "main",
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		branches := []string{"main"}
		getAutomationsForBranch = func(branch string) []string {
			return branches
		}

		updatedJSON, _ := db.ReplaceJSONplaceholders(replacements, constants.ComponentFilterQuery)

		_ = UpdateFiltersForDrilldown(updatedJSON, replacements, true, true)

	})

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
			"branch":    "main",
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		branches := "main"
		getBranchNameForId = func(branch string) string {
			return branches
		}

		updatedJSON, _ := db.ReplaceJSONplaceholders(replacements, constants.ComponentFilterQuery)

		_ = UpdateFiltersForDrilldown(updatedJSON, replacements, false, true)

	})

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
			"branch":    "main",
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		updatedJSON, _ := db.ReplaceJSONplaceholders(replacements, constants.ComponentFilterQuery)

		_ = UpdateFiltersForDrilldown(updatedJSON, replacements, false, false)

	})

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"hello"},
			"branch":    "main",
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		updatedJSON, _ := db.ReplaceJSONplaceholders(replacements, constants.ComponentFilterQuery)

		_ = UpdateFiltersForDrilldown(updatedJSON, replacements, false, false)

	})

}

func Test_UpdateFiltersForSonar(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
			"branch":    "main",
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		branches := []string{"main"}
		getAutomationsForBranch = func(branch string) []string {
			return branches
		}

		updatedJSON, _ := db.ReplaceJSONplaceholders(replacements, constants.ComponentFilterQuery)

		_ = UpdateFiltersForSonar(updatedJSON, replacements)

	})

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
			"branch":    "main",
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		branches := "main"
		getBranchNameForId = func(branch string) string {
			return branches
		}

		updatedJSON, _ := db.ReplaceJSONplaceholders(replacements, constants.ComponentFilterQuery)

		_ = UpdateFiltersForSonar(updatedJSON, replacements)

	})

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
			"branch":    "main",
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		updatedJSON, _ := db.ReplaceJSONplaceholders(replacements, constants.ComponentFilterQuery)

		_ = UpdateFiltersForSonar(updatedJSON, replacements)

	})

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"hello"},
			"branch":    "main",
		}

		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		updatedJSON, _ := db.ReplaceJSONplaceholders(replacements, constants.ComponentFilterQuery)

		_ = UpdateFiltersForSonar(updatedJSON, replacements)

	})

}

func Test_mock_content(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	var client client.GrpcClient

	t.Run("Success", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":        "org123",
			"subOrgId":     "suborg123",
			"component":    []string{"All"},
			"branch":       "main",
			"FunctionName": "test",
		}

		_, _ = mockContent("s9", replacements, context.Background(), client, nil)

	})

	t.Run("failure", func(t *testing.T) {

		replacements := map[string]any{
			"orgId":        "org123",
			"subOrgId":     "suborg123",
			"component":    []string{"All"},
			"branch":       "main",
			"FunctionName": "test",
		}

		_, _ = mockContent("test", replacements, context.Background(), client, nil)

	})

}

func Test_getOpenIssuesSection(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	testReplacements := map[string]any{
		"orgId":     "org123",
		"subOrgId":  "suborg123",
		"startDate": "2023-01-01",
		"endDate":   "2023-12-31",
		"component": []string{"All"},
	}
	responseStr := `[{"dateOfDiscovery":"2023:10:27 09:36","drillDown":{"reportId":"nested-drilldown-view-location","reportInfo":{"branch":"registry.saas-dev.beescloud.com/staging/compliance-snykcontainer:33","code":"CVE-2023-30590","component_id":"e29bc352-20f8-49a7-b457-639dcca94efc","run_id":"53fed825-270b-4ac6-836c-94bd5289e252","scanner_name":"snykcontainer"},"reportTitle":"Open Issues"},"recurrences":2,"scannerName":"snykcontainer","severity":"Medium","sla":"At Risk","vulnerabilityId":"CVE-2023-30590"},{"dateOfDiscovery":"2023:10:27 09:36","drillDown":{"reportId":"nested-drilldown-view-location","reportInfo":{"branch":"registry.saas-dev.beescloud.com/staging/compliance-snykcontainer:33","code":"CVE-2018-5709","component_id":"e29bc352-20f8-49a7-b457-639dcca94efc","run_id":"53fed825-270b-4ac6-836c-94bd5289e252","scanner_name":"snykcontainer"},"reportTitle":"Open Issues"},"recurrences":8,"scannerName":"snykcontainer","severity":"Low","sla":"At Risk","vulnerabilityId":"CVE-2018-5709"},{"dateOfDiscovery":"2023:10:27 08:27","drillDown":{"reportId":"nested-drilldown-view-location","reportInfo":{"branch":"registry.saas-dev.beescloud.com/staging/compliance-snykcontainer:17","code":"CVE-2016-2781","component_id":"e29bc352-20f8-49a7-b457-639dcca94efc","run_id":"05622053-9e3d-4309-b055-b670edaf12d4","scanner_name":"snykcontainer"},"reportTitle":"Open Issues"},"recurrences":2,"scannerName":"snykcontainer","severity":"Low","sla":"At Risk","vulnerabilityId":"CVE-2016-2781"}]`

	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			args: args{
				widgetId:     "test",
				replacements: testReplacements,
				ctx:          context.Background(),
				clt:          mockGrpcClient,
			},
			want: []byte(responseStr),
		},
	}
	searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
		return `{"hits":{"total":{"value":642,"relation":"eq"}},"aggregations":{"drilldowns":{"value":[{"drillDown":{"reportId":"nested-drilldown-view-location","reportTitle":"Open Issues","reportInfo":{"component_id":"e29bc352-20f8-49a7-b457-639dcca94efc","code":"CVE-2023-30590","run_id":"53fed825-270b-4ac6-836c-94bd5289e252","scanner_name":"snykcontainer","branch":"registry.saas-dev.beescloud.com/staging/compliance-snykcontainer:33"}},"severity":"Medium","recurrences":2,"vulnerabilityId":"CVE-2023-30590","dateOfDiscovery":"2023:10:27 09:36","sla":"At Risk","scannerName":"snykcontainer"},{"drillDown":{"reportId":"nested-drilldown-view-location","reportTitle":"Open Issues","reportInfo":{"component_id":"e29bc352-20f8-49a7-b457-639dcca94efc","code":"CVE-2018-5709","run_id":"53fed825-270b-4ac6-836c-94bd5289e252","scanner_name":"snykcontainer","branch":"registry.saas-dev.beescloud.com/staging/compliance-snykcontainer:33"}},"severity":"Low","recurrences":8,"vulnerabilityId":"CVE-2018-5709","dateOfDiscovery":"2023:10:27 09:36","sla":"At Risk","scannerName":"snykcontainer"},{"drillDown":{"reportId":"nested-drilldown-view-location","reportTitle":"Open Issues","reportInfo":{"component_id":"e29bc352-20f8-49a7-b457-639dcca94efc","code":"CVE-2016-2781","run_id":"05622053-9e3d-4309-b055-b670edaf12d4","scanner_name":"snykcontainer","branch":"registry.saas-dev.beescloud.com/staging/compliance-snykcontainer:17"}},"severity":"Low","recurrences":2,"vulnerabilityId":"CVE-2016-2781","dateOfDiscovery":"2023:10:27 08:27","sla":"At Risk","scannerName":"snykcontainer"}]}}}`, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getOpenIssuesSection(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOpenIssuesSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("got : %s", got)
				t.Logf("want : %s", tt.want)
				t.Errorf("getOpenIssuesSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getComponentOpenIssuesHeader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	testReplacements := map[string]any{
		"orgId":     "org123",
		"subOrgId":  "suborg123",
		"startDate": "2023-01-01",
		"endDate":   "2023-12-31",
		"component": []string{"All"},
	}
	responseStr := `{"HIGH":{"value":193},"LOW":{"value":121},"MEDIUM":{"value":220},"TOTAL":{"value":664},"VERY_HIGH":{"value":130}}`

	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			args: args{
				widgetId:     "test",
				replacements: testReplacements,
				ctx:          context.Background(),
				clt:          mockGrpcClient,
			},
			want: []byte(responseStr),
		},
	}
	searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
		return `{"hits":{"total":{"value":642,"relation":"eq"}},"aggregations":{"severityCounts":{"value":{"TOTAL":664,"VERY_HIGH":130,"HIGH":193,"MEDIUM":220,"LOW":121}}}}`, nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getComponentOpenIssuesHeader(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("getComponentOpenIssuesHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("got : %s", got)
				t.Logf("want : %s", tt.want)
				t.Errorf("getComponentOpenIssuesHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOpenIssuesSection01(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}
	testReplacements := map[string]any{
		"orgId":     "org123",
		"subOrgId":  "suborg123",
		"startDate": "2023-01-01",
		"endDate":   "2023-12-31",
		"component": []string{"All"},
	}

	type args struct {
		widgetId     string
		replacements map[string]any
		ctx          context.Context
		clt          client.GrpcClient
	}
	tests := []struct {
		name    string
		args    args
		want    json.RawMessage
		wantErr bool
	}{
		{
			args: args{
				widgetId:     "test",
				replacements: testReplacements,
				ctx:          context.Background(),
				clt:          mockGrpcClient,
			},
			wantErr: true,
		},
	}
	searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
		return "", db.ErrInternalServer
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getOpenIssuesSection(tt.args.widgetId, tt.args.replacements, tt.args.ctx, tt.args.clt, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOpenIssuesSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOpenIssuesSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergedDefaultBranchCommitsSection(t *testing.T) {

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {

		responseString := `{
			"automationRunsCount": {
			  "aggregations": {
				"automation_run": {
				  "value": {
					"data": {
					  "848a228d-3244-471f-9292-8df08e908099": {
						"staging": 6
					  },
					  "f1ccd216-8283-41e9-9187-925b2535bf83": {
						"staging": 2
					  }
					},
					"totalCount": 9851
				  }
				}
			  }
			},
			"deployedAutomationCount": {
			  "aggregations": {
				"automation_run": {
				  "value": {
					"data": {
					  "221bc7e2-a00d-408d-a5b6-b72a9477f2df_5229adcc-9f60-4be1-99aa-f2d375bb38cc": {
						"staging": 1
					  },
					  "221bc7e2-a00d-408d-a5b6-b72a9477f2df_1c4512fe-ec87-4f1d-9069-134db09f71a8": {
						"staging": 1
					  },
					  "221bc7e2-a00d-408d-a5b6-b72a9477f2df_9aff4071-5436-475b-8590-85cb1fcd914c": {
						"staging": 1
					  },
					  "221bc7e2-a00d-408d-a5b6-b72a9477f2df_ba0e7143-e111-47c6-927a-f7384b47a953": {
						"staging": 1
					  }
					}
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

		reports, err := mergedDefaultBranchCommitsSection("", x, nil)
		assert.Equal(t, len(reports), 207, "Validating response count for automation run")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
	})
	t.Run("Case 2: Error on invalid query key", func(t *testing.T) {
		responseString := `{
				"automationRunsCount_fail": {
				  "aggregations": {
					"automation_run": {
					  "value": {
						"data": {
						  "848a228d-3244-471f-9292-8df08e908099": {
							"staging": 6
						  },
						  "f1ccd216-8283-41e9-9187-925b2535bf83": {
							"staging": 2
						  }
						},
						"totalCount": 9851
					  }
					}
				  }
				},
				"deployedAutomationCount": {
				  "aggregations": {
					"automation_run": {
					  "value": {
						"data": {
						  "221bc7e2-a00d-408d-a5b6-b72a9477f2df_5229adcc-9f60-4be1-99aa-f2d375bb38cc": {
							"staging": 1
						  },
						  "221bc7e2-a00d-408d-a5b6-b72a9477f2df_1c4512fe-ec87-4f1d-9069-134db09f71a8": {
							"staging": 1
						  },
						  "221bc7e2-a00d-408d-a5b6-b72a9477f2df_9aff4071-5436-475b-8590-85cb1fcd914c": {
							"staging": 1
						  },
						  "221bc7e2-a00d-408d-a5b6-b72a9477f2df_ba0e7143-e111-47c6-927a-f7384b47a953": {
							"staging": 1
						  }
						}
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

		_, err := mergedDefaultBranchCommitsSection("", x, nil)
		assert.Error(t, err, db.ErrInternalServer)

	})

}

func Test_mergedDefaultBranchCommitsHeader(t *testing.T) {

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		responseString := `{
			"automationRunsCount": {
				"aggregations": {
					"automation_run": {
						"value": {
					  		"data": {
								"848a228d-3244-471f-9292-8df08e908099": {
						  		"staging": 6
							},
							"f1ccd216-8283-41e9-9187-925b2535bf83": {
						  		"staging": 2
							}

						},
						"totalCount": 9851
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

		reports, err := mergedDefaultBranchCommitsHeader("", x, nil)

		assert.Equal(t, len(reports), 14, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

	t.Run("Case 2: Error om missing query key", func(t *testing.T) {

		responseString := `{
			"aggregations": {
				"automation_run_fail": {
					"value": {
					  "data": {
						"c9be43f4-47e0-463a-a7ff-81b060878231_3dda89cf-30ab-40e6-b0ca-aa2e0a60c3c9": 1,
						"5713172e-af92-4db0-b15a-132abc258dff_972708bc-6ad0-43f6-aea2-701d14f56dbb": 1,
						"fd51dbbe-18a0-434a-bd3a-f16752ab3fdb_11267517-2ce6-4e9e-a524-bd9bb257c1c6": 1
					  },
					  "totalCount": 3
					}
				}
			}
		  }`

		x := map[string]json.RawMessage{}

		json.Unmarshal([]byte(responseString), &x)

		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		_, err := mergedDefaultBranchCommitsHeader("", x, nil)
		assert.Error(t, err, db.ErrInternalServer)

	})

}

func Test_testInsightsWidgetSection(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	tests := []struct {
		name         string
		orgId        string
		component    []string
		replacements map[string]any
		want         json.RawMessage
		wantErr      bool
	}{
		{
			name:      "Case 1: Successful execution with valid parameters",
			orgId:     "org123",
			component: []string{"All"},
			replacements: map[string]any{
				"orgId":     "org123",
				"subOrgId":  "suborg123",
				"component": []string{"All"},
			},
			wantErr: false,
			want:    json.RawMessage(`{"key": "value"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			automationMap := map[string]struct{}{
				"auto1": {},
			}

			openSearchClient = func() (*opensearch.Client, error) {
				return opensearch.NewDefaultClient()
			}
			responseString := `{
				"aggregations": {
					"test_insights_automations": {
						"value": [
							"036d13c2-3009-4dae-8f09-2d86cc066ee8",
							"b772fc94-3613-484e-9236-9f53f28b7503",
							"f1ccd216-8283-41e9-9187-925b2535bf83"
						  ]
					}
				}
			  }`
			searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
				return responseString, nil
			}

			GetAutomationMap = func(ctx context.Context, clt client.GrpcClient, orgId string, component []string, branch string) map[string]struct{} {
				return automationMap
			}

			got, err := testInsightsWidgetSection(tt.orgId, tt.replacements, context.Background(), mockGrpcClient, nil)

			if tt.wantErr && err == nil {
				t.Errorf("Expected an error, but got none.")
			}

			if err != nil {
				assert.Equal(t, len(got), 447, "Validating response count for test insights workflows")
			}
		})
	}
}

func Test_testComponentWidgetSection(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)
	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	t.Run("Case 1: Successful execution with components", func(t *testing.T) {
		widgetId := "widget123"

		testReplacements := map[string]any{
			"orgId":     "org123",
			"subOrgId":  "suborg123",
			"component": []string{"All"},
		}

		ctx := context.Background()
		openSearchClient = func() (*opensearch.Client, error) {
			return opensearch.NewDefaultClient()
		}

		responseString := `{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":705,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"distinct_component":{"value":["3333344c-ae0f-4df4-b1a7-efcaacdf449e","484d5e12-6424-4070-a159-4e5639a807a2"]}}}`
		searchResponse = func(query string, IndexName string, client *opensearch.Client) (string, error) {
			return responseString, nil
		}

		serviceResponse := &service.ListServicesResponse{
			Service: []*service.Service{

				{
					Id:            "484d5e12-6424-4070-a159-4e5639a807a2",
					Name:          "dsl-engine-cli",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli.git",
				},
				{
					Id:            "3333344c-ae0f-4df4-b1a7-efcaacdf449e",
					Name:          "dsl-engine",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine.git",
				},
				{
					Id:            "94a81dd1-3f52-4520-891e-b2440f660945",
					Name:          "dsl-engine-cli9",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli9.git",
				},
				{
					Id:            "c497a177-931e-4ab9-8822-0c4d39d2acfc",
					Name:          "dsl-engine7",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine7.git",
				},
				{
					Id:            "6537a7f8-1b1a-4f5d-8413-ca40e9b57893",
					Name:          "dsl-engine-cli5",
					RepositoryUrl: "https://github.com/calculi-corp/dsl-engine-cli5.git",
				},
			},
		}

		getOrganisationServices = func(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
			return serviceResponse, nil
		}

		reports, err := testComponentWidgetSection(widgetId, testReplacements, ctx, mockGrpcClient, nil)
		assert.Equal(t, len(reports), 477, "Validating response count for automation run drilldown")
		// Assertions to check if the result is as expected.
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

	})

}

func Test_UpdateTestSuitesOverviewResponse(t *testing.T) {
	t.Run("Case 1: Successful execution of test suites overview", func(t *testing.T) {

		responseString := `{"took":1,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":6,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"testSuitesOverview":{"value":{"apple":{"skipped_cases_count":0,"component_id":"apple","run_id":"apple","total_cases":42,"component_name":"ci-insights-service","duration_in_millis":60,"average_duration":60,"automation_id":"apple","duration":60,"successful_cases_count":42,"start_time":"2024/08/01 11:06","total_duration_in_millis":60,"failure_rate_for_last_run":"0.0%","branch_id":"apple","run_id_set":["apple"],"org_id":"apple","failed_cases_count":0,"branch_name":"main","workflow_runs":1,"start_time_in_millis":1722490571000,"automation_name":"workflow","test_suite_runs":1,"test_suite_name":"github.com"}}}}}`

		expectResult := `[{"testSuiteName":"github.com","componentName":"ci-insights-service","workflow":"workflow","source":"CloudBees","branch":"main","defaultBranch":"main","lastRun":"2024/08/01 11:06","lastRunInMillis":1722490571000,"totalTestCases":{"value":42,"drillDown":{"reportId":"test-overview-total-tests-cases","reportTitle":"Test cases - github.com","reportType":"","reportInfo":{"branch":"apple","test_suite_name":"github.com","automation_id":"apple","component_name":"ci-insights-service","workflow_name":"workflow","branch_name":"main","source":"CloudBees"}}},"totalTestCasesValue":42,"avgRunTime":60,"totalRuns":{"value":1,"drillDown":{"reportId":"test-overview-total-runs","reportTitle":"Runs - github.com","reportType":"","reportInfo":{"branch":"apple","test_suite_name":"github.com","automation_id":"apple"}}},"totalRunsValue":1,"failureRate":{"colorScheme":[{"color0":"#009C5B","color1":"#62CA9D"},{"color0":"#D32227","color1":"#FB6E72"},{"color0":"#F2A414","color1":"#FFE6C1"}],"data":[{"title":"Successful test cases","value":42},{"title":"Failed test cases","value":0},{"title":"Skipped test cases","value":0}],"lightColorScheme":[{"color0":"#0C9E61","color1":"#79CAA8"},{"color0":"#E83D39","color1":"#F39492"},{"color0":"#F2A414","color1":"#FFE6C1"}],"type":"SINGLE_BAR","value":"0.0%"},"failureRateValue":0}]`
		b, err := updateTestSuitesOverviewResponse(responseString)
		if err != nil {
			return
		}
		bString, _ := json.Marshal(b)
		assert.Nil(t, err, "error processing test suites overview")
		assert.Equal(t, expectResult, string(bString))
	})
}

func Test_UpdateTestComponentsViewResponse(t *testing.T) {
	t.Run("Case 1: Successful execution of test components overview", func(t *testing.T) {

		responseString := `{"took":3,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":72,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"workflow_buckets":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"2281041f-7381-491b-aec0-ca2e5cd9061d_99b1055f-92d2-4cdd-a327-10edd986d98e","doc_count":44,"total_test_cases_runs":{"value":1102},"latest_doc":{"hits":{"total":{"value":44,"relation":"eq"},"max_score":null,"hits":[{"_index":"cb_test_suites","_id":"8509888e-d27f-44fa-46a9-29bc76f5e790_99b1055f-92d2-4cdd-a327-10edd986d98e_2281041f-7381-491b-aec0-ca2e5cd9061d_7bfad64a-6d10-4f09-b608-112ef74b6579_7dad5a3f-45f5-47d8-9bf4-ae375f9cb9d9_github.com/calculi-corp/rbac-service","_score":null,"_source":{"automation_id":"99b1055f-92d2-4cdd-a327-10edd986d98e","total":1,"run_id":"7dad5a3f-45f5-47d8-9bf4-ae375f9cb9d9","branch_id":"7bfad64a-6d10-4f09-b608-112ef74b6579","branch_name":"main","component_name":"rbac-service","run_start_time":"2024-07-02 17:47:50","automation_name":"workflow"},"fields":{"zoned_run_start_time":["2024/07/02 23:17:50"],"run_start_time_in_millis":[1719942470000]},"sort":[1719942470000]}]}},"total_duration":{"value":148069},"skipped_count":{"value":10},"success_count":{"value":1092},"failure_count":{"value":0},"failure_rate":{"value":10},"avg_run_time":{"value":134.3638838475499},"total_test_cases_count":{"value":639}},{"key":"3911ded3-6c6f-4a67-96a7-c54ec93fe12b_02c8817a-07cc-4db3-bce9-6afd13766a4d","doc_count":15,"total_test_cases_runs":{"value":1542},"latest_doc":{"hits":{"total":{"value":15,"relation":"eq"},"max_score":null,"hits":[{"_index":"cb_test_suites","_id":"8509888e-d27f-44fa-46a9-29bc76f5e790_02c8817a-07cc-4db3-bce9-6afd13766a4d_3911ded3-6c6f-4a67-96a7-c54ec93fe12b_94de10c7-e441-498b-b7f3-863684d06990_207c77fe-f555-4670-b100-b390bffa4eef_github.com/calculi-corp/reports-service/helper","_score":null,"_source":{"automation_id":"02c8817a-07cc-4db3-bce9-6afd13766a4d","total":51,"run_id":"207c77fe-f555-4670-b100-b390bffa4eef","branch_id":"94de10c7-e441-498b-b7f3-863684d06990","branch_name":"main","component_name":"reports-service","run_start_time":"2024-07-02 17:18:26","automation_name":"workflow"},"fields":{"zoned_run_start_time":["2024/07/02 22:48:26"],"run_start_time_in_millis":[1719940706000]},"sort":[1719940706000]}]}},"total_duration":{"value":700},"skipped_count":{"value":3},"success_count":{"value":1539},"failure_count":{"value":0},"failure_rate":{"value":8},"avg_run_time":{"value":0.45395590142671854},"total_test_cases_count":{"value":639}},{"key":"62b6124f-4ba6-44d4-a83a-2dbd9aae97ce_5823e667-31c1-4dc9-a4a5-cd26f11328c1","doc_count":7,"total_test_cases_runs":{"value":186},"latest_doc":{"hits":{"total":{"value":7,"relation":"eq"},"max_score":null,"hits":[{"_index":"cb_test_suites","_id":"8509888e-d27f-44fa-46a9-29bc76f5e790_5823e667-31c1-4dc9-a4a5-cd26f11328c1_62b6124f-4ba6-44d4-a83a-2dbd9aae97ce_4f37435c-efda-4374-8aa6-80b3fa124105_40707dcd-ff4e-42b5-b56b-7d81437766c6_github.com/calculi-corp/asset-service/listeners/endpoints","_score":null,"_source":{"automation_id":"5823e667-31c1-4dc9-a4a5-cd26f11328c1","total":2,"run_id":"40707dcd-ff4e-42b5-b56b-7d81437766c6","branch_id":"4f37435c-efda-4374-8aa6-80b3fa124105","branch_name":"main","component_name":"asset-service","run_start_time":"2024-07-02 09:24:06","automation_name":"workflow"},"fields":{"zoned_run_start_time":["2024/07/02 14:54:06"],"run_start_time_in_millis":[1719912246000]},"sort":[1719912246000]}]}},"total_duration":{"value":80},"skipped_count":{"value":0},"success_count":{"value":186},"failure_count":{"value":0},"failure_rate":{"value":7},"avg_run_time":{"value":0.43010752688172044},"total_test_cases_count":{"value":639}},{"key":"23a56310-fb50-42cc-ad0d-4909b9639911_c636f9b0-dbda-437e-a652-290d926fa858","doc_count":4,"total_test_cases_runs":{"value":202},"latest_doc":{"hits":{"total":{"value":4,"relation":"eq"},"max_score":null,"hits":[{"_index":"cb_test_suites","_id":"8509888e-d27f-44fa-46a9-29bc76f5e790_c636f9b0-dbda-437e-a652-290d926fa858_23a56310-fb50-42cc-ad0d-4909b9639911_91227d14-036f-4963-8081-8e611f656f46_ea3e5b1a-5360-4a1a-8865-c0cfbf211bed_github.com/calculi-corp/projectmgmt-service/dbutil","_score":null,"_source":{"automation_id":"c636f9b0-dbda-437e-a652-290d926fa858","total":35,"run_id":"ea3e5b1a-5360-4a1a-8865-c0cfbf211bed","branch_id":"91227d14-036f-4963-8081-8e611f656f46","branch_name":"main","component_name":"projectmgmt-service","run_start_time":"2024-07-02 17:18:29","automation_name":"workflow"},"fields":{"zoned_run_start_time":["2024/07/02 22:48:29"],"run_start_time_in_millis":[1719940709000]},"sort":[1719940709000]}]}},"total_duration":{"value":4060},"skipped_count":{"value":1},"success_count":{"value":201},"failure_count":{"value":0},"failure_rate":{"value":5},"avg_run_time":{"value":20.099009900990097},"total_test_cases_count":{"value":639}},{"key":"6ee2c296-5bc2-423a-8bef-2c2cdb344635_a8252bce-394c-4f37-b2d8-13e5be78d6dc","doc_count":2,"total_test_cases_runs":{"value":43},"latest_doc":{"hits":{"total":{"value":2,"relation":"eq"},"max_score":null,"hits":[{"_index":"cb_test_suites","_id":"8509888e-d27f-44fa-46a9-29bc76f5e790_a8252bce-394c-4f37-b2d8-13e5be78d6dc_6ee2c296-5bc2-423a-8bef-2c2cdb344635_663d7e40-6700-4705-97f9-236be4681530_00077e8c-9706-401b-9295-f1c0e6801ce2_github.com/calculi-corp/cbci-actions/internal/actions","_score":null,"_source":{"automation_id":"a8252bce-394c-4f37-b2d8-13e5be78d6dc","total":40,"run_id":"00077e8c-9706-401b-9295-f1c0e6801ce2","branch_id":"663d7e40-6700-4705-97f9-236be4681530","branch_name":"main","component_name":"cbci-actions","run_start_time":"2024-07-02 11:51:55","automation_name":"workflow"},"fields":{"zoned_run_start_time":["2024/07/02 17:21:55"],"run_start_time_in_millis":[1719921115000]},"sort":[1719921115000]}]}},"total_duration":{"value":20020},"skipped_count":{"value":0},"success_count":{"value":43},"failure_count":{"value":0},"failure_rate":{"value":0},"avg_run_time":{"value":465.5813953488372},"total_test_cases_count":{"value":639}}]}}}`

		expectResult := `[{"avgRunTime":134.36,"defaultBranch":"main","componentName":"rbac-service","failureRate":{"colorScheme":[{"color0":"#009C5B","color1":"#62CA9D"},{"color0":"#D32227","color1":"#FB6E72"},{"color0":"#F2A414","color1":"#FFE6C1"}],"data":[{"title":"Successful test case runs","value":1092},{"title":"Failed test case runs","value":0},{"title":"Skipped test case runs","value":10}],"lightColorScheme":[{"color0":"#0C9E61","color1":"#79CAA8"},{"color0":"#E83D39","color1":"#F39492"},{"color0":"#F2A414","color1":"#FFE6C1"}],"type":"SINGLE_BAR","value":"10.0%"},"failureRateValue":10,"lastRun":"2024/07/02 23:17:50","lastRunInMillis":1719942470000,"totalTestCasesValue":639,"totalTestCases":{"value":639,"drillDown":{"reportId":"test-overview-total-tests-cases","reportTitle":"Test cases - rbac-service","reportType":"","reportInfo":{"branch":"7bfad64a-6d10-4f09-b608-112ef74b6579","automation_id":"99b1055f-92d2-4cdd-a327-10edd986d98e","component_name":"rbac-service","workflow_name":"workflow","branch_name":"main","source":"CloudBees"}}},"workflow":"workflow","source":"CloudBees"},{"avgRunTime":0.45,"defaultBranch":"main","componentName":"reports-service","failureRate":{"colorScheme":[{"color0":"#009C5B","color1":"#62CA9D"},{"color0":"#D32227","color1":"#FB6E72"},{"color0":"#F2A414","color1":"#FFE6C1"}],"data":[{"title":"Successful test case runs","value":1539},{"title":"Failed test case runs","value":0},{"title":"Skipped test case runs","value":3}],"lightColorScheme":[{"color0":"#0C9E61","color1":"#79CAA8"},{"color0":"#E83D39","color1":"#F39492"},{"color0":"#F2A414","color1":"#FFE6C1"}],"type":"SINGLE_BAR","value":"8.0%"},"failureRateValue":8,"lastRun":"2024/07/02 22:48:26","lastRunInMillis":1719940706000,"totalTestCasesValue":639,"totalTestCases":{"value":639,"drillDown":{"reportId":"test-overview-total-tests-cases","reportTitle":"Test cases - reports-service","reportType":"","reportInfo":{"branch":"94de10c7-e441-498b-b7f3-863684d06990","automation_id":"02c8817a-07cc-4db3-bce9-6afd13766a4d","component_name":"reports-service","workflow_name":"workflow","branch_name":"main","source":"CloudBees"}}},"workflow":"workflow","source":"CloudBees"},{"avgRunTime":0.43,"defaultBranch":"main","componentName":"asset-service","failureRate":{"colorScheme":[{"color0":"#009C5B","color1":"#62CA9D"},{"color0":"#D32227","color1":"#FB6E72"},{"color0":"#F2A414","color1":"#FFE6C1"}],"data":[{"title":"Successful test case runs","value":186},{"title":"Failed test case runs","value":0},{"title":"Skipped test case runs","value":0}],"lightColorScheme":[{"color0":"#0C9E61","color1":"#79CAA8"},{"color0":"#E83D39","color1":"#F39492"},{"color0":"#F2A414","color1":"#FFE6C1"}],"type":"SINGLE_BAR","value":"7.0%"},"failureRateValue":7,"lastRun":"2024/07/02 14:54:06","lastRunInMillis":1719912246000,"totalTestCasesValue":639,"totalTestCases":{"value":639,"drillDown":{"reportId":"test-overview-total-tests-cases","reportTitle":"Test cases - asset-service","reportType":"","reportInfo":{"branch":"4f37435c-efda-4374-8aa6-80b3fa124105","automation_id":"5823e667-31c1-4dc9-a4a5-cd26f11328c1","component_name":"asset-service","workflow_name":"workflow","branch_name":"main","source":"CloudBees"}}},"workflow":"workflow","source":"CloudBees"},{"avgRunTime":20.1,"defaultBranch":"main","componentName":"projectmgmt-service","failureRate":{"colorScheme":[{"color0":"#009C5B","color1":"#62CA9D"},{"color0":"#D32227","color1":"#FB6E72"},{"color0":"#F2A414","color1":"#FFE6C1"}],"data":[{"title":"Successful test case runs","value":201},{"title":"Failed test case runs","value":0},{"title":"Skipped test case runs","value":1}],"lightColorScheme":[{"color0":"#0C9E61","color1":"#79CAA8"},{"color0":"#E83D39","color1":"#F39492"},{"color0":"#F2A414","color1":"#FFE6C1"}],"type":"SINGLE_BAR","value":"5.0%"},"failureRateValue":5,"lastRun":"2024/07/02 22:48:29","lastRunInMillis":1719940709000,"totalTestCasesValue":639,"totalTestCases":{"value":639,"drillDown":{"reportId":"test-overview-total-tests-cases","reportTitle":"Test cases - projectmgmt-service","reportType":"","reportInfo":{"branch":"91227d14-036f-4963-8081-8e611f656f46","automation_id":"c636f9b0-dbda-437e-a652-290d926fa858","component_name":"projectmgmt-service","workflow_name":"workflow","branch_name":"main","source":"CloudBees"}}},"workflow":"workflow","source":"CloudBees"},{"avgRunTime":465.58,"defaultBranch":"main","componentName":"cbci-actions","failureRate":{"colorScheme":[{"color0":"#009C5B","color1":"#62CA9D"},{"color0":"#D32227","color1":"#FB6E72"},{"color0":"#F2A414","color1":"#FFE6C1"}],"data":[{"title":"Successful test case runs","value":43},{"title":"Failed test case runs","value":0},{"title":"Skipped test case runs","value":0}],"lightColorScheme":[{"color0":"#0C9E61","color1":"#79CAA8"},{"color0":"#E83D39","color1":"#F39492"},{"color0":"#F2A414","color1":"#FFE6C1"}],"type":"SINGLE_BAR","value":"0.0%"},"failureRateValue":0,"lastRun":"2024/07/02 17:21:55","lastRunInMillis":1719921115000,"totalTestCasesValue":639,"totalTestCases":{"value":639,"drillDown":{"reportId":"test-overview-total-tests-cases","reportTitle":"Test cases - cbci-actions","reportType":"","reportInfo":{"branch":"663d7e40-6700-4705-97f9-236be4681530","automation_id":"a8252bce-394c-4f37-b2d8-13e5be78d6dc","component_name":"cbci-actions","workflow_name":"workflow","branch_name":"main","source":"CloudBees"}}},"workflow":"workflow","source":"CloudBees"}]`
		b, err := updateTestComponentsViewResponse(responseString)
		if err != nil {
			return
		}
		bString, _ := json.Marshal(b)
		assert.Nil(t, err, "error processing test components overview")
		assert.Equal(t, expectResult, string(bString))
	})
}

func Test_UpdateTestCasesOverviewResponse(t *testing.T) {
	t.Run("Case 1: Successful execution of test cases overview", func(t *testing.T) {

		responseString := `{"took":21,"timed_out":false,"_shards":{"total":2,"successful":2,"skipped":0,"failed":0},"hits":{"total":{"value":391,"relation":"eq"},"max_score":null,"hits":[]},"aggregations":{"testCasesOverview":{"value":{"bb7f71a0-2325-4905-ba0a-805461a50006_bb2dff3a-e5f3-4d86-abb4-38e3d6196ecc_a17c0ea1-4bc0-4f81-8828-0041962984ad__TestAdd":{"component_id":"bb7f71a0-2325-4905-ba0a-805461a50006","skipped_count":0,"component_name":"template-go-test-result-conversion","duration_in_millis":100000,"success_count":0,"failure_count":7,"start_time":"2024-06-17T15:01:39.000Z","average_duration":100000,"automation_id":"bb2dff3a-e5f3-4d86-abb4-38e3d6196ecc","test_case_name":"TestAdd","duration":100,"failure_rate":"100.0%","total_duration_in_millis":700000,"branch_id":"a17c0ea1-4bc0-4f81-8828-0041962984ad","org_id":"8509888e-d27f-44fa-46a9-29bc76f5e790","branch_name":"preprod","start_time_in_millis":1718636499000,"automation_name":"workflow","total_exec_count":7,"runs":7,"test_suite_name":"github.com/calculi-corp/template-go-testing/test_suite_2","status":"FAILED"}}}}}`

		expectResult := `[{"testSuiteName":"github.com/calculi-corp/template-go-testing/test_suite_2","testCaseName":"TestAdd","componentName":"template-go-test-result-conversion","workflow":"workflow","source":"CloudBees","branch":"preprod","lastRun":"2024-06-17T15:01:39.000Z","lastRunInMillis":1718636499000,"avgRunTime":100000,"totalRuns":{"value":7,"drillDown":{"reportId":"test-overview-view-run-activity","reportTitle":"Runs - TestAdd","reportType":"","reportInfo":{"component_id":"bb7f71a0-2325-4905-ba0a-805461a50006","branch":"a17c0ea1-4bc0-4f81-8828-0041962984ad","test_suite_name":"github.com/calculi-corp/template-go-testing/test_suite_2","test_case_name":"TestAdd","automation_id":"bb2dff3a-e5f3-4d86-abb4-38e3d6196ecc","workflow_name":"workflow","source":"CloudBees"}}},"totalRunsValue":7,"failureRate":{"type":"SINGLE_BAR","colorScheme":[{"color0":"#009C5B","color1":"#62CA9D"},{"color0":"#D32227","color1":"#FB6E72"}],"lightColorScheme":[{"color0":"#0C9E61","color1":"#79CAA8"},{"color0":"#E83D39","color1":"#F39492"}],"value":"100.0%","data":[{"title":"Successful runs","value":0},{"title":"Failed runs","value":7}]},"failureRateValue":100}]`
		b, err := updateTestCasesOverviewResponse(responseString)
		if err != nil {
			return
		}
		bString, _ := json.Marshal(b)
		assert.Nil(t, err, "error processing test cases overview")
		assert.Equal(t, expectResult, string(bString))
	})
}

func Test_millisecondsToSecondsString(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{input: 1000.0, expected: "1s"},
		{input: 1500.0, expected: "1.5s"},
		{input: 250.123, expected: "0.250123s"},
		{input: 0.0, expected: "0s"},
		{input: 987654321.0, expected: "987654.321s"},
	}

	for _, tt := range tests {
		t.Run(strconv.FormatFloat(tt.input, 'f', -1, 64), func(t *testing.T) {
			result := millisecondsToSecondsString(tt.input)
			if result != tt.expected {
				t.Errorf("millisecondsToSecondsString(%f) = %s; expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func Test_getTestsOverview_2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGrpcClient := mock_grpc_client.NewMockGrpcClient(mockCtrl)

	grpc.SetSharedGrpcClient(mockGrpcClient)
	defer grpc.SetSharedGrpcClient(nil)

	openSearchClient = func() (*opensearch.Client, error) {
		return opensearch.NewDefaultClient()
	}

	testReplacements := map[string]interface{}{
		"orgId":      "org123",
		"subOrgId":   "suborg123",
		"startDate":  "2023-01-01",
		"endDate":    "2023-12-31",
		"component":  []string{"All"},
		"viewOption": "testSuite",
	}

	testCaseReplacements := map[string]interface{}{
		"orgId":      "org123",
		"subOrgId":   "suborg123",
		"startDate":  "2023-01-01",
		"endDate":    "2023-12-31",
		"component":  []string{"All"},
		"viewOption": "testCase",
	}

	tests := []struct {
		name             string
		replacements     map[string]interface{}
		expectedResponse string
		expectedError    bool
	}{
		{
			name:             "TestSuite View Option",
			replacements:     testReplacements,
			expectedResponse: `<expected response for test suite>`,
			expectedError:    false,
		},
		{
			name:             "TestCase View Option",
			replacements:     testCaseReplacements,
			expectedResponse: `<expected response for test case>`,
			expectedError:    false,
		},
		{
			name:             "Invalid View Option",
			replacements:     map[string]interface{}{"viewOption": "invalidOption"},
			expectedResponse: "",
			expectedError:    true,
		},
	}

	suiteResponse := `{"mocked": "suite response"}`
	caseResponse := `{"mocked": "case response"}`

	searchResponse = func(query, IndexName string, client *opensearch.Client) (string, error) {
		if IndexName == constants.TEST_SUITE_INDEX {
			return suiteResponse, nil
		} else if IndexName == constants.TEST_CASES_INDEX {
			return caseResponse, nil
		} else {
			return "", errors.New("unexpected index name")
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getTestsOverview("widget123", tt.replacements, context.Background(), mockGrpcClient, nil)
		})
	}
}
