package helper

import (
	"github.com/calculi-corp/api/go/endpoint"
	"github.com/calculi-corp/reports-service/cache"
)

func GetEndpointsByContributionId(resourceId, contributionId string) []*endpoint.Endpoint {
	epCache := cache.GetEndpointCache()
	return epCache.GetByContributionID(resourceId, contributionId)
}
