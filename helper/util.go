package helper

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	cutils "github.com/calculi-corp/common/utils"

	api "github.com/calculi-corp/api/go"
	"github.com/calculi-corp/api/go/auth"
	"github.com/calculi-corp/api/go/auth/permission"
	"github.com/calculi-corp/api/go/endpoint"
	"github.com/calculi-corp/api/go/service"
	client "github.com/calculi-corp/grpc-client"
	hostflags "github.com/calculi-corp/grpc-hostflags"
	srvauth "github.com/calculi-corp/grpc-server/auth"
	"github.com/calculi-corp/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/calculi-corp/api/go/vsm/report"
	"github.com/calculi-corp/reports-service/cache"
	"github.com/calculi-corp/reports-service/constants"
	db "github.com/calculi-corp/reports-service/db"
	"golang.org/x/exp/slices"
)

const (
	automationServiceName = "api.service.AutomationService"
	serviceEndpointName   = "api.service.ServiceEndpoint"
	automationMethod      = "ListOrganizationAutomations"
	serviceMethod         = "ListServices"
)

// ciInsightsWidgets is a map of valid CI Insights widgets
var ciInsightsWidgets = map[string]bool{
	"ci1":  true,
	"ci2":  true,
	"ci4":  true,
	"ci3":  true,
	"ci5":  true,
	"ci6":  true,
	"ci7":  true,
	"ci01": true,
}

func GetOrganisationServiceAutomations(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListOrganizationAutomationsResponse, error) {
	response := &service.ListOrganizationAutomationsResponse{}
	request := &service.ListOrganizationAutomationsRequest{
		OrgId: orgId,
		Runs:  -1,
		Pagination: &api.Pagination{
			Page:       1,
			PageLength: 500,
		},
	}
	err := clt.SendGrpcCtx(ctx, hostflags.RepositoryServiceHost(), automationServiceName, automationMethod, request, response) // Send the request
	log.CheckErrorf(err, "Failed to get a response from server : ")
	if response != nil && response.Pagination != nil {
		isLastPage := response.Pagination.LastPage
		if !isLastPage {
			page := response.Pagination.Page
			length := response.Pagination.PageLength
			// Send the request
			getNextPageData(page, orgId, length, clt, ctx, response)
		}
	}
	return response, nil
}

func getNextPageData(page int32, orgId string, length int32, clt client.GrpcClient, ctx context.Context, response *service.ListOrganizationAutomationsResponse) {
	nextResponse := &service.ListOrganizationAutomationsResponse{}
	nextPage := page + 1
	log.Debug("Fetching nextPage automation info for page : " + string(nextPage))
	request := &service.ListOrganizationAutomationsRequest{
		OrgId: orgId,
		Runs:  -1,
		Pagination: &api.Pagination{
			Page:       nextPage,
			PageLength: length,
		},
	}
	err := clt.SendGrpcCtx(ctx, hostflags.RepositoryServiceHost(), automationServiceName, automationMethod, request, nextResponse)
	if err != nil {
		log.Error("Failed to get a response from server for page : "+string(nextPage), err)
	} else if nextResponse != nil && len(nextResponse.GetServices()) != 0 {
		response.Services = append(response.Services, nextResponse.GetServices()...)
		if nextResponse.Pagination != nil && !nextResponse.Pagination.LastPage {
			getNextPageData(nextResponse.Pagination.Page, orgId, nextResponse.Pagination.PageLength, clt, ctx, response)
		}
	}

}

func GetOrganisationServices(ctx context.Context, clt client.GrpcClient, orgId string) (*service.ListServicesResponse, error) {
	if orgId != "" {
		response := &service.ListServicesResponse{}
		request := &service.ListServicesRequest{
			OrgId: orgId,
		}
		ctx = srvauth.SystemUserCtx(ctx, "reports-service:GetOrganisationServices")
		err := clt.SendGrpcCtx(ctx, hostflags.RepositoryServiceHost(), serviceEndpointName, serviceMethod, request, response) // Send the request
		if log.CheckErrorf(err, "Failed to get a response from server : ") {
			return nil, err
		}
		return response, nil
	} else {
		return nil, errors.New("org Id should not be null")
	}
}

func FetchOrganizationAndServices(ctx context.Context, clt client.GrpcClient, client auth.OrganizationsServiceClient, orgId, userId string) (*constants.Organization, []string, error) {
	components := []string{}
	getOrganizationByIdResponse, err := GetOrganisationsById(ctx, client, orgId, userId)
	if log.CheckErrorf(err, "Error getting organization for org - %s", orgId) {
		return nil, nil, err
	}
	serviceResponse, err := GetOrganisationServices(ctx, clt, orgId)
	organization := &constants.Organization{ID: getOrganizationByIdResponse.GetOrganization().GetId(),
		Name: getOrganizationByIdResponse.GetOrganization().GetDisplayName()}
	if log.CheckErrorf(err, "Error getting organization service for org id - %s", orgId) {
		return nil, nil, err
	}
	for _, service := range serviceResponse.GetService() {
		if organization.ID == service.OrganizationId {
			organization.Components = append(organization.Components, &constants.Component{ID: service.Id, Name: service.Name})
			components = append(components, service.Id)
		}
	}
	log.Debugf("Nested organisations : %v", organization)
	getServicesRecursively(organization, getOrganizationByIdResponse.GetOrganization(), serviceResponse, components)
	log.Debugf("Fetched organisation and it's nested sub org and services : %+v", organization)
	return organization, components, nil
}

func GetOrganisationsById(ctx context.Context, clt auth.OrganizationsServiceClient, orgId, userId string) (*auth.GetOrganizationByIdResponse, error) {
	getOrganizationsByIdRequest := &auth.GetOrganizationByIdRequest{UserId: userId, Id: orgId, Nested: true}
	getOrganizationByIdResponse, err := clt.GetOrganizationById(ctx, getOrganizationsByIdRequest)
	if err != nil {
		log.Errorf(err, "Error getting organization for org - %s", orgId)
		return nil, err
	}
	return getOrganizationByIdResponse, nil
}

func getServicesRecursively(parentOrg *constants.Organization, sourceOrg *auth.Organization, serviceResponse *service.ListServicesResponse, components []string) {
	for _, childOrg := range sourceOrg.GetChildOrganizations() {
		log.Debugf("Fetching components from child org : %v", childOrg)
		subOrg := &constants.Organization{ID: childOrg.GetId(), Name: childOrg.GetDisplayName()}
		for _, service := range serviceResponse.GetService() {
			if subOrg.ID == service.OrganizationId {
				subOrg.Components = append(subOrg.Components, &constants.Component{ID: service.Id, Name: service.Name})
				components = append(components, service.Id)
			}
		}
		log.Debugf("Fetched components from child org : %s, %v", subOrg.Name, subOrg.Components)
		parentOrg.SubOrgs = append(parentOrg.SubOrgs, subOrg)
		getServicesRecursively(subOrg, childOrg, serviceResponse, components)
	}
}

func GetComponentsRecursively(sourceOrg *auth.Organization, serviceResponse *service.ListServicesResponse, components *[]string) {
	for _, childOrg := range sourceOrg.GetChildOrganizations() {
		log.Debugf("Fetching components from child org : %v", childOrg)
		for _, service := range serviceResponse.GetService() {
			if childOrg.GetId() == service.OrganizationId {
				*components = append(*components, service.Id)
			}
		}
		log.Debugf("Fetched components from child org : %s, %v", childOrg.GetId(), components)
		GetComponentsRecursively(childOrg, serviceResponse, components)
	}
}

func GetAllEndpoints(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string, contributionIds []string, includeParent bool) (*endpoint.EndpointsResponse, error) {
	request := &endpoint.ListEndpointsRequest{
		ResourceId: orgId,
		Filter: &endpoint.EndpointsFilter{
			ResourceIds: []string{orgId},
		},
		Pagination: &api.Pagination{
			Page:       1,
			PageLength: 100,
		},
		Parents: includeParent,
	}
	if len(contributionIds) > 0 {
		request.Filter.ContributionIds = contributionIds
	}
	response, err := clt.ListEndpoints(ctx, request)
	if err != nil {
		return nil, err
	}
	if response != nil && response.Pagination != nil {
		isLastPage := response.Pagination.LastPage
		if !isLastPage {
			page := response.Pagination.Page
			length := response.Pagination.PageLength
			getNextPageEndpoint(page, orgId, length, clt, ctx, response, contributionIds)
		}
	}
	return response, nil
}

func getNextPageEndpoint(page int32, orgId string, length int32, clt endpoint.EndpointServiceClient, ctx context.Context, response *endpoint.EndpointsResponse, contributionIds []string) {
	nextPage := page + 1
	request := &endpoint.ListEndpointsRequest{
		ResourceId: orgId,
		Filter: &endpoint.EndpointsFilter{
			ResourceIds: []string{orgId},
		},
		Pagination: &api.Pagination{
			Page:       nextPage,
			PageLength: length,
		},
	}
	if len(contributionIds) > 0 {
		request.Filter.ContributionIds = contributionIds
	}
	nextResponse, err := clt.ListEndpoints(ctx, request)
	if err != nil {
		log.Error("Failed to get a response from server for page : "+string(nextPage), err)
	} else if nextResponse != nil && len(nextResponse.GetEndpoints()) != 0 {
		response.Endpoints = append(response.Endpoints, nextResponse.GetEndpoints()...)
		if nextResponse.Pagination != nil && !nextResponse.Pagination.LastPage {
			getNextPageEndpoint(nextResponse.Pagination.Page, orgId, nextResponse.Pagination.PageLength, clt, ctx, response, contributionIds)
		}
	}
}

// IsResourceInOrg gets the parent resource tree for the resource and checks if the organization matches
func IsResourceInOrg(resourceId string, organizationId string) bool {
	if resourceId == organizationId {
		return true
	}
	coreDataCache := cache.GetCoreDataCache()
	if coreDataCache != nil {
		resourceParents := coreDataCache.GetParentIDs(resourceId)
		if len(resourceParents) > 0 {
			return slices.Contains(resourceParents, organizationId)
		}
	}
	return false
}

// Checks if layout for all screens are empty
func CheckEmptyDashboardLayout(layout *db.VsmDashboardLayout) bool {
	if len(layout.DashboardLayout.Lg) == 0 && len(layout.DashboardLayout.Md) == 0 && len(layout.DashboardLayout.Xl) == 0 && len(layout.DashboardLayout.Xs) == 0 {
		return true
	}
	return false
}

// isAuthorizedToAccessCIInsights checks if the user is authorized to access CI Insights
func isAuthorizedToAccessCIInsights(ctx context.Context, orgId string, rbacClient auth.RBACServiceClient) (bool, error) {

	// Only admin users are allowed to access CI Insights
	// admin users must have CREATE permission on CI Insights
	isAuthorizedRequest := auth.IsAuthorizedRequest{
		ResourceId: orgId,
		Permissions: []*permission.Permission{
			{
				Action: permission.PermissionAction_CREATE,
				Type:   permission.ApiType_CI_INSIGHTS,
			},
		},
	}

	isAuthorizedResponse, err := rbacClient.IsAuthorized(ctx, &isAuthorizedRequest)
	if err != nil {
		log.Errorf(err, "Error checking rbac")
		return false, status.Error(codes.Internal, "Internal error")
	}
	if isAuthorizedResponse == nil {
		log.Debugf("Empty response from rbac service for org: %s", orgId)
		return false, status.Error(codes.Internal, "Internal error")
	}
	return isAuthorizedResponse.GetAuthorized(), nil

}

// ValidateCIInsightsReportAccess validates if the user has access to CI Insights reports
func ValidateCIInsightsReportAccess(ctx context.Context, rbacClient auth.RBACServiceClient, req *pb.ReportServiceRequest) error {
	// validate only if widget id is present and it is a valid CI Insights widget or CiToolId is present
	if (req.WidgetId != "" && ciInsightsWidgets[req.WidgetId]) || req.CiToolId != "" {
		if isAuthorized, err := isAuthorizedToAccessCIInsights(ctx, req.OrgId, rbacClient); err != nil {
			return err
		} else {
			if !isAuthorized {
				return status.Error(codes.PermissionDenied, "permission denied")
			}
		}
	}
	return nil
}

// ValidateCIInsightsDrilldownReportAccess validates if the user has access to CI Insights drilldown reports
func ValidateCIInsightsDrilldownReportAccess(ctx context.Context, rbacClient auth.RBACServiceClient, req *pb.DrilldownRequest) error {
	// validate only if CiToolId is present
	// CiToolId is present only for CI Insights drilldown reports
	if req.CiToolId != "" {
		if isAuthorized, err := isAuthorizedToAccessCIInsights(ctx, req.OrgId, rbacClient); err != nil {
			return err
		} else {
			if !isAuthorized {
				return status.Error(codes.PermissionDenied, "permission denied")
			}
		}
	}
	return nil
}

// Truncate to at most one decimal place while preserving the format
// Input: 20.437457 -> Output: 20.4
// Input: 20.4 -> Output: 20.4
func TruncateFloat(val float32) float32 {
	truncated := float32(math.Trunc(float64(val)*10) / 10) // Convert to float64, truncate, and convert back

	// Check if the original value already had exactly 1 decimal place
	if float64(val)*10 == math.Floor(float64(val)*10) {
		return val // Keep it unchanged if it originally had 1 decimal place
	}

	return truncated
}

// Function to convert "00.0%" string to float32
// Input: "25.4%" -> Output: 25.4
func ConvertPercentageToFloat(value string) float32 {
	cleaned := strings.TrimSuffix(value, "%")
	floatVal, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		fmt.Println("Error converting string to float for getTestsOverview():", err)
		return 0
	}
	return float32(floatVal)
}

func GetWorkflowDisplayName(workflowName string) string {
	// Get display name for the workflow
	// This is required because the automation resource names are stored in different formats internally for external workflows (Ex: GHA workflows)
	workflowDisplayName, _, err := cutils.GetDisplayNameAndOrigin(workflowName)
	if err != nil {
		log.Error("Error getting display name for the workflow", err)
		// If error occurs, use the workflow name as is
		workflowDisplayName = workflowName
	}
	return workflowDisplayName
}

func ContainsString(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
