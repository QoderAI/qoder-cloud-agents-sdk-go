package managed_test

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/testsupport"
	"testing"
)

func TestManagedFailureContracts(t *testing.T) {
	var endpoints []testsupport.Endpoint
	for _, c := range contracts(t) {
		endpoints = append(endpoints, testsupport.Endpoint{Service: c.Service, Method: c.Name})
	}
	testsupport.FailureContracts(t, endpoints, func(fn testsupport.Transport) any { return testClient(roundTripFunc(fn)) })
}
