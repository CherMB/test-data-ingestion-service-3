package db

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"strings"

	"github.com/calculi-corp/log"
	"github.com/calculi-corp/reports-service/exceptions"
	opensearch "github.com/opensearch-project/opensearch-go"

	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

// query OpenSeach and get the required data
func GetOpensearchData(query string, IndexName string, client *opensearch.Client) (string, error) {
	content := strings.NewReader(query)
	search := opensearchapi.SearchRequest{
		Index: []string{IndexName},
		Body:  content,
	}
	searchResponse, err := search.Do(context.Background(), client)
	if searchResponse != nil {
		log.Debugf("Search data response status code:%d and error:%v", searchResponse.StatusCode, searchResponse.IsError())
		if searchResponse.IsError() || (searchResponse.StatusCode != 200 && searchResponse.StatusCode != 201) {
			log.Infof("Search data indexName : %s error response : %v", IndexName, searchResponse)
			return "", ErrInternalServer
		}
	}
	if log.CheckErrorf(err, exceptions.ErrSearchingDocInHelperGetOpenSearchData) {
		return "", ErrInternalServer
	}
	defer searchResponse.Body.Close()
	return formResponse(searchResponse), nil
}

// query OpenSeach and get the required data
func GetOpensearchMappingData(IndexName string, client *opensearch.Client) (string, error) {

	search := opensearchapi.IndicesGetMappingRequest{
		Index: []string{IndexName},
	}

	searchResponse, err := search.Do(context.Background(), client)
	if log.CheckErrorf(err, exceptions.ErrSearchingDocInHelperGetOpenSearchData) {
		return "", ErrInternalServer
	}
	defer searchResponse.Body.Close()
	return formResponse(searchResponse), nil
}

// convert the response to a string. taken from the String() func inside the opensearchapi package and updated
func formResponse(r *opensearchapi.Response) string {

	var (
		out = new(bytes.Buffer)
		b1  = bytes.NewBuffer([]byte{})
		b2  = bytes.NewBuffer([]byte{})
		tr  io.Reader
	)

	if r != nil && r.Body != nil {
		tr = io.TeeReader(r.Body, b1)
		defer r.Body.Close()

		if _, err := io.Copy(b2, tr); err != nil {
			out.WriteString(fmt.Sprintf("<error reading response body: %v>", err))
			return out.String()
		}
		defer func() { r.Body = io.NopCloser(b1) }()
	}

	if r != nil && r.Body != nil {
		out.ReadFrom(b2) // errcheck exclude (*bytes.Buffer).ReadFrom
	}

	return out.String()
}

func InsertBulkData(client *opensearch.Client, bulkQuery string) error {
	_, err := client.Bulk(
		strings.NewReader(bulkQuery),
	)

	if log.CheckErrorf(err, "could not insert bulk data to opensearch") {
		return ErrDataIngestionFailed
	}

	return nil
}

func GetMultiQueryData(client *opensearch.Client, multiQuery string) (string, error) {
	searchResponse, err := client.Msearch(
		strings.NewReader(multiQuery),
	)

	if log.CheckErrorf(err, "Multi query search to opensearch failed") || searchResponse.IsError() {
		log.Errorf(ErrInternalServer, "multi query failed with status: %s error: %v", searchResponse.Status(), err)
		return "", ErrInvalidRequest
	}

	defer searchResponse.Body.Close()
	return formResponse(searchResponse), nil
}

func GetOpensearchCount(query string, IndexName string, client *opensearch.Client) (string, error) {
	content := strings.NewReader(query)
	search := opensearchapi.CountRequest{
		Index: []string{IndexName},
		Body:  content,
	}
	searchResponse, err := search.Do(context.Background(), client)
	if searchResponse != nil {
		log.Debugf("Search data response status code:%d and error:%v", searchResponse.StatusCode, searchResponse.IsError())
		if searchResponse.IsError() || (searchResponse.StatusCode != 200 && searchResponse.StatusCode != 201) {
			log.Infof("Search data indexName : %s error response : %v", IndexName, searchResponse)
			return "", ErrInternalServer
		}
	}
	if log.CheckErrorf(err, exceptions.ErrSearchingDocInHelperGetOpenSearchData) {
		return "", ErrInternalServer
	}
	defer searchResponse.Body.Close()
	return formResponse(searchResponse), nil
}
