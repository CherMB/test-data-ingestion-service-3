package cache

import (
	"github.com/calculi-corp/common/defines"
	"github.com/calculi-corp/common/pkg/messaging"
	"github.com/calculi-corp/core-data-cache/secrets"
	client "github.com/calculi-corp/grpc-client"
	"github.com/calculi-corp/log"

	coredata "github.com/calculi-corp/core-data-cache"
)

var (
	CoreDataResourceCache coredata.ResourceCacheI
	GrpcClient            client.GrpcClient
	epCache               coredata.EndpointsCacheI
)

func SetMockCache(mockResourceCache coredata.ResourceCacheI, mockEndpointCache coredata.EndpointsCacheI, clt client.GrpcClient) {
	CoreDataResourceCache = mockResourceCache
	epCache = mockEndpointCache
	GrpcClient = clt
}

func InitializeCache(clt client.GrpcClient, msgClt messaging.Messaging) error {
	CoreDataResourceCache = coredata.CreateResourceCache(msgClt, clt)
	GrpcClient = clt
	err := CoreDataResourceCache.Start()
	if log.CheckErrorf(err, "Failed to start Core data resource cache") {
		return err
	}

	sclt, err := secrets.NewClient(clt)
	if err != nil {
		return err
	}

	epCache = coredata.CreateEndpointsCache(msgClt, clt, CoreDataResourceCache, sclt)
	if epCache == nil {
		return defines.ErrInternal
	}

	return nil

}

func GetCoreDataCache() coredata.ResourceCacheI {
	return CoreDataResourceCache
}

func GetEndpointCache() coredata.EndpointsCacheI {
	return epCache
}

func GetGrpcClient() client.GrpcClient {
	return GrpcClient
}
