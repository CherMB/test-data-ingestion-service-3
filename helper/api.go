package helper

import (
	"crypto/tls"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	// Create a TLS config that accepts any certificate
	tlsConfig := &tls.Config{InsecureSkipVerify: false}
	// Create a new transport that uses the TLS config
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	// Create an HTTP client that uses the transport
	Client = &http.Client{Transport: transport, Timeout: 15 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return nil
	}}
}

func Get(url string, encodedToken string) (*http.Response, error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Error while making the GET request "+err.Error())
	}
	if encodedToken != "" {
		// Create a Basic string by appending string access token
		var basicToken = "Basic " + encodedToken
		// add authorization header to the req
		req.Header.Add("Authorization", basicToken)
	}
	resp, err := Client.Do(req)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Error while Invoking the GET request "+err.Error())
	}
	return resp, nil
}
