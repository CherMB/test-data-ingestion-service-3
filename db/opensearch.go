package db

import (
	"context"
	"sync"

	"github.com/calculi-corp/log"
	opensearchconfig "github.com/calculi-corp/opensearch-config"
	opensearch "github.com/opensearch-project/opensearch-go"
)

var (
	newConnection = opensearchconfig.GetOpensearchConnection
	instance      *opensearch.Client
	once          sync.Once
)

// GetOpenSearchClient returns the instance of the OpenSearchClient.
func GetOpenSearchClient() *opensearch.Client {
	chk := opensearchconfig.CheckOpensearchClient(context.Background(), instance)
	if !chk {
		clt, err := opensearchconfig.GetOpensearchConnection()
		if err != nil {
			log.Errorf(err, "could not connect to Opensearch")
		}
		instance = clt
	}

	return instance
}
