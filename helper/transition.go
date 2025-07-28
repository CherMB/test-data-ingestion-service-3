package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"

	api "github.com/calculi-corp/api/go"
	"github.com/calculi-corp/api/go/auth"
	"github.com/calculi-corp/api/go/endpoint"
	client "github.com/calculi-corp/grpc-client"
	"github.com/calculi-corp/log"
	"github.com/calculi-corp/reports-service/cache"
	"github.com/calculi-corp/reports-service/constants"
	db "github.com/calculi-corp/reports-service/db"
)

type ProjectMapping struct {
	ProjectKey       string             `protobuf:"bytes,1,opt,name=project_key,json=projectKey,proto3" json:"project_key,omitempty"`
	ProjectName      string             `protobuf:"bytes,2,opt,name=project_name,json=projectName,proto3" json:"project_name,omitempty"`
	FlowItemMapping  []*FlowItemMapping `protobuf:"bytes,3,rep,name=flow_item_mapping,json=flowItemMapping,proto3" json:"flow_item_mapping,omitempty"`
	InProgressStatus []string           `protobuf:"bytes,4,rep,name=in_progress_status,json=inProgressStatus,proto3" json:"in_progress_status,omitempty"`
	WaitingStatus    []string           `protobuf:"bytes,5,rep,name=waiting_status,json=waitingStatus,proto3" json:"waiting_status,omitempty"`
}

type FlowItemMapping struct {
	Name         string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	MappingType  string `protobuf:"bytes,2,opt,name=mapping_type,json=mappingType,proto3" json:"mapping_type,omitempty"`
	MappingValue string `protobuf:"bytes,3,opt,name=mapping_value,json=mappingValue,proto3" json:"mapping_value,omitempty"`
}

func GetComponents(ctx context.Context, clt client.GrpcClient, orgServiceClient auth.OrganizationsServiceClient, orgId, userId string) []string {

	components := []string{}
	// Get components for the org
	getOrganizationByIdResponse, err := GetOrganisationsById(ctx, orgServiceClient, orgId, userId)
	if err != nil {
		return components
	}

	serviceResponse, err := GetOrganisationServices(ctx, clt, orgId)
	if err != nil {
		return components
	}
	for _, service := range serviceResponse.GetService() {
		if orgId == service.OrganizationId {
			components = append(components, service.Id)
		}
	}
	GetComponentsRecursively(getOrganizationByIdResponse.GetOrganization(), serviceResponse, &components)

	return components
}

func GetWorkflowsCount(ctx context.Context, orgId string, components []string) int {
	resources := make([]*api.Resource, 0)

	for _, c := range components {
		childrens := GetResourceChildrenByType([]string{c}, 100, false, api.ResourceType_RESOURCE_TYPE_AUTOMATION)
		resources = append(resources, childrens...)
	}

	return len(resources)
}

func GetResourceChildrenByType(ids []string, depth int32, includeDisabled bool, rType api.ResourceType) []*api.Resource {
	if depth == 0 {
		return []*api.Resource{}
	}
	out := []*api.Resource{}
	for _, id := range ids {
		res := cache.CoreDataResourceCache.Get(id)
		if res == nil {
			log.Errorf(errors.New("Failed to fetch resource"), "Resource not found in cache: %s", id)
		}

		if res == nil || (!includeDisabled && res.GetIsDisabled()) {
			continue
		}
		if res.GetType() == rType {
			out = append(out, res)
			return out
		}
		out = append(out, GetResourceChildrenByType(cache.CoreDataResourceCache.GetChildren(id), depth-1, includeDisabled, rType)...)
		if len(out) > 0 {
			return out
		}
	}
	return out
}

func GetFlowItemMappings(ctx context.Context, clt endpoint.EndpointServiceClient, orgId string) int {
	flowItemMappingsCount := 0
	contributionIds := []string{constants.JIRA_ENDPOINT_CONTRIBUTION_ID}
	endPointsResponse, err := GetAllEndpoints(ctx, clt, orgId, contributionIds, true)
	log.CheckErrorf(err, "endpoint api failed")

	if err == nil && len(endPointsResponse.Endpoints) > 0 {
		for _, endpoint := range endPointsResponse.Endpoints {
			if endpoint.IsDisabled {
				continue
			}
			projectMappings := GetAllProjectMappingsFromEndpoint(endpoint)
			for _, projectMapping := range projectMappings {
				if projectMapping.ProjectKey == "" || len(projectMapping.FlowItemMapping) == 0 {
					continue
				}
				flowItemMappingsCount++
			}
		}
	}
	return flowItemMappingsCount
}

func GetAllProjectMappingsFromEndpoint(ep *endpoint.Endpoint) []ProjectMapping {
	projectMappings := []ProjectMapping{}
	for _, property := range ep.GetProperties() {
		if property.GetName() == "projectMapping" && property.GetData() != nil {
			bytes := property.GetData().GetBytes()
			json.Unmarshal(bytes, &projectMappings)
		}
	}
	return projectMappings
}

func GetIndexDocCount(index, orgId string, components []string) int {
	// Get the filters from request and create map of replacement placeholders
	replacements := map[string]any{
		// Remove orgId since ui only sends suborgId or orgId
		// "orgId":     orgId,
		"component": components,
	}
	baseQuery, err := db.ReplaceJSONplaceholders(replacements, constants.GetWorkflowRunsCount)
	if log.CheckErrorf(err, "could not replace json placeholders in GetIndexDocCount() for %s :", orgId) {
		return 0
	}

	reader := strings.NewReader(baseQuery)

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	queryBytes := buf.String()

	var data map[string]interface{}
	err = json.Unmarshal([]byte(queryBytes), &data)
	if log.CheckErrorf(err, "could not unmarshal queryBytes :") {
		return 0
	}

	filterArray := data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{})

	if replacements["component"].([]string) != nil {
		filter3 := AddTermsFilter("component_id", replacements["component"].([]string))
		filterArray = append(filterArray, filter3)
	}

	data["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filterArray

	updatedQuery, err := json.MarshalIndent(data, "", " ")
	if log.CheckErrorf(err, "error converting to json") {
		return 0
	}

	countQuery := db.DbQuery{AliasName: index, QueryString: string(updatedQuery)}

	resp, err := GetCountQueryResponse(countQuery)
	if log.CheckErrorf(err, "could not get computed data") {
		return 0
	}

	return resp["count"]
}

func GetIndexDocCountByOrgId(index, orgId string) int {
	// Get the filters from request and create map of replacement placeholders
	replacements := map[string]any{
		"orgId": orgId,
	}
	baseQuery, err := db.ReplaceJSONplaceholders(replacements, constants.DocsCountQuery)
	if log.CheckErrorf(err, "could not replace json placeholders in GetIndexDocCount() for %s :", orgId) {
		return 0
	}

	reader := strings.NewReader(baseQuery)

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	queryBytes := buf.String()

	var data map[string]interface{}
	err = json.Unmarshal([]byte(queryBytes), &data)
	if log.CheckErrorf(err, "could not unmarshal queryBytes :") {
		return 0
	}

	updatedQuery, err := json.MarshalIndent(data, "", " ")
	if log.CheckErrorf(err, "error converting to json") {
		return 0
	}
	countQuery := db.DbQuery{AliasName: index, QueryString: string(updatedQuery)}

	resp, err := GetCountQueryResponse(countQuery)
	if log.CheckErrorf(err, "could not get computed data") {
		return 0
	}

	return resp["count"]
}
