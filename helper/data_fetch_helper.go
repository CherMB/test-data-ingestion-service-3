package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	opensearchconfig "github.com/calculi-corp/opensearch-config"
	db "github.com/calculi-corp/reports-service/db"

	"github.com/calculi-corp/log"
)

type OpenSearchResponse struct {
	key  string
	data string
	err  error
}

type queryResponseJson struct {
	Hits struct {
		Hits []struct {
		} `json:"hits"`
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
	} `json:"hits"`
	Aggregation struct {
	} `json:"aggregation"`
}

type countResponseJson struct {
	Count int `json:"count"`
}

func GetQueryResponse(qm map[string]db.DbQuery) (map[string]json.RawMessage, error) {

	if qm == nil {
		return nil, nil
	}

	finalOut := make(map[string]json.RawMessage)

	//creating a channel to listen to Opensearch responses
	responseChannel := make(chan OpenSearchResponse)

	//get the response for each query and append it to the final output
	//spawn a separate go routine for each query and get the response
	numElements := len(qm)
	log.Debugf("GetQueryResponse - Number of elements in the map: %d", numElements)
	for k, q := range qm {
		go runQuery(k, q, responseChannel)
	}
	//listen for opensearch responses on the response channel and append them to the final output
	for range qm {
		//this is a blocking call
		out := <-responseChannel
		if log.CheckErrorf(out.err, "error received in GetData() response channel") {

			return nil, out.err
		}

		finalOut[out.key] = json.RawMessage(out.data)
	}

	return finalOut, nil
}

func runQuery(key string, query db.DbQuery, responseChannel chan OpenSearchResponse) error {
	startTime := time.Now()
	//create an opensearch client
	client, err := opensearchconfig.GetOpensearchConnection()
	if log.CheckErrorf(err, "Error establishing connection with OpenSearch in getQueryResponse()") {

		responseChannel <- OpenSearchResponse{key: key, data: "", err: err}
		return err
	}
	log.Debugf("Time took to OpenSearch Connection for Alias %s : %v in milliseconds", query.AliasName, time.Since(startTime).Milliseconds())

	//send a request to opensearch with the query that was fetched
	startTime = time.Now()
	searchResponse, err := db.GetOpensearchData(query.QueryString, query.AliasName, client)
	if log.CheckErrorf(err, "Error fetching Opensearch data in service.getResponse():") {

		responseChannel <- OpenSearchResponse{key: key, data: "", err: err}
		return err
	}

	responseChannel <- OpenSearchResponse{key: key, data: searchResponse, err: err}

	log.Debugf("Time took to Get OpenSearch Data for Alias %s : %v in milliseconds", query.AliasName, time.Since(startTime).Milliseconds())

	return nil
}

func IsResponseEmpty(qr map[string]json.RawMessage) bool {

	isEmpty := true
	for _, data := range qr {
		dataString := string(data)
		queryResponse := queryResponseJson{}
		err := json.Unmarshal([]byte(dataString), &queryResponse)
		if err == nil && queryResponse.Hits.Total.Value > 0 {
			isEmpty = false
		}
	}

	return isEmpty
}

// HasNoHits checks if there were no document hits for the OpenSearch query, and returns a bool
func HasNoHits(responseString string) bool {

	isEmpty := true

	queryResponse := queryResponseJson{}
	err := json.Unmarshal([]byte(responseString), &queryResponse)
	if err == nil && queryResponse.Hits.Total.Value > 0 {
		isEmpty = false
	}

	return isEmpty
}

func GetMultiQueryResponse(qm map[string]db.DbQuery) (map[string]json.RawMessage, error) {
	if qm == nil {
		return nil, nil
	}

	client, err := opensearchconfig.GetOpensearchConnection()
	if log.CheckErrorf(err, "Error establishing connection with OpenSearch in GetMultiQueryResponse()") {
		return nil, err
	}
	//Construct the Multi search query request
	queryKey, multiQuery, err := constructMultiQuery(qm)
	if log.CheckErrorf(err, "Error constructing Multi Search query GetMultiQueryResponse()") {
		return nil, err
	}

	//send Multi Search request to opensearch
	dataStr, err := db.GetMultiQueryData(client, multiQuery)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(dataStr), &data)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in GetMultiQueryResponse()") {
		return nil, err
	}

	d := data["responses"].([]interface{})
	//Expect number of queries in qm and number of responses in data.responses to be same. if not error out
	if len(d) != len(qm) {
		log.Error("mismatch in response and query count", db.ErrMultiSearchResponseError)
		return nil, db.ErrMultiSearchResponseError
	}

	finalOut := make(map[string]json.RawMessage)
	var index = 0
	for _, v := range queryKey {
		md, err := json.Marshal(d[index])
		if log.CheckErrorf(err, "Error marshaling response from OpenSearch in GetMultiQueryResponse()") {
			return nil, err
		}
		finalOut[v] = json.RawMessage(md)
		index++
	}

	return finalOut, nil
}

func constructMultiQuery(qm map[string]db.DbQuery) ([]string, string, error) {
	var queryKey []string
	var multiQuery string
	buffer := new(bytes.Buffer)

	for k, v := range qm {
		queryKey = append(queryKey, k)
		buffer.Reset()
		if err := json.Compact(buffer, []byte(v.QueryString)); err != nil {
			log.Errorf(err, "failed to construct query : ", k)
			return nil, "", err
		}
		multiQuery += fmt.Sprintf("%s%v%s", "{\"index\": \""+v.AliasName+"\"}\n", buffer, "\n")

	}

	return queryKey, multiQuery, nil
}

// HasField returns a bool indicating if a field exists in an index's mapping.
func HasField(field, key string, inputMap map[string]map[string]string) (bool, error) {
	// Query OpenSearch instead of storing this information in a map if necessary in the future
	if indices, ok := inputMap[field]; !ok {
		return false, fmt.Errorf("input field not found in FieldIndexMap")
	} else {
		if _, ok := indices[key]; ok {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func GetCountQueryResponse(query db.DbQuery) (map[string]int, error) {

	client, err := opensearchconfig.GetOpensearchConnection()
	if log.CheckErrorf(err, "Error establishing connection with OpenSearch in GetCountQueryResponse()") {
		return nil, err
	}

	//send count Search request to opensearch
	dataStr, err := db.GetOpensearchCount(query.QueryString, query.AliasName, client)
	if err != nil {
		return nil, err
	}

	var data countResponseJson
	err = json.Unmarshal([]byte(dataStr), &data)
	if log.CheckErrorf(err, "Error unmarshaling response from OpenSearch in GetCountQueryResponse()") {
		return nil, err
	}

	finalOut := make(map[string]int)
	finalOut["count"] = data.Count

	return finalOut, nil
}
