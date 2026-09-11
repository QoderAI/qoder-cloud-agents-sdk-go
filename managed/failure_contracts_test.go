package managed_test

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/testutil"
	"testing"
)

func TestManagedFailureContracts(t *testing.T) {
	var endpoints []testutil.Endpoint
	for _, c := range contracts(t) {
		endpoints = append(endpoints, testutil.Endpoint{Service: c.Service, Method: c.Name})
	}
	testutil.FailureContracts(t, endpoints, func(fn testutil.Transport) any { return testClient(roundTripFunc(fn)) })
}
