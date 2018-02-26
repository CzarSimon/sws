package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var client = newClient()

// BodyRequest signature of request functions containg a body.
type BodyRequest func(route string, body []byte) (*http.Response, error)

// makePostRequest performs a post request with api credentials.
func makePostRequest(route string, body []byte) (*http.Response, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	req := newRequest("POST", route, config, bytes.NewBuffer(body))
	return client.Do(req)
}

// makeDeleteRequest performs a delete request with api credentials.
func makeDeleteRequest(route string, body []byte) (*http.Response, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	req := newRequest("DELETE", route, config, bytes.NewBuffer(body))
	return client.Do(req)
}

// makeGetRequest performs a get request with api credentials.
func makeGetRequest(route string) (*http.Response, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	req := newRequest("GET", route, config, nil)
	return client.Do(req)
}

// newRequest Creates a new http request of a supplied method
// to the apiserver and embends the api access key
func newRequest(method, route string, config Config, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, config.API.ToURL(route), body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", config.Auth.AccessKey)
	return req
}

// newClient sets up a new client for use.
func newClient() *http.Client {
	const TIMEOUT_SECONDS = 5
	return &http.Client{
		Timeout: time.Second * TIMEOUT_SECONDS,
	}
}
