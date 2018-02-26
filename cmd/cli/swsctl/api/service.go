package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/CzarSimon/sws/pkg/service"
)

const serviceRoute = "v1/service"

// GetServices gets the list of users current active services.
func GetServices() ([]service.Service, error) {
	resp, err := makeGetRequest(serviceRoute)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	services := make([]service.Service, 0)
	err = json.NewDecoder(resp.Body).Decode(&services)
	return services, err
}

// PostService sends a service specification to the apiserver.
func PostService(svc service.Service) (string, error) {
	return serviceChangeRequest(svc, makePostRequest)
}

// DeleteService send a delete request to the apiserver.
func DeleteService(svc service.Service) (string, error) {
	return serviceChangeRequest(svc, makeDeleteRequest)
}

// serviceChangeRequest executes a change request method for a supplied service.
func serviceChangeRequest(svc service.Service, req BodyRequest) (string, error) {
	bytes, err := json.Marshal(svc)
	if err != nil {
		return "", err
	}
	resp, err := req(serviceRoute, bytes)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}
