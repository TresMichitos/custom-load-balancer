// Integration tests for custom load balancer

package integration

import (
	"net/http"
	"testing"
)

// Send requests to running load balancer and check responses
func TestLoadBalancerFunctionality (t *testing.T) {
	var lburl string = "http://localhost:8080"
	resp, err := http.Get(lburl)
	if err != nil {
		t.Errorf("http.Get returned err %s", err.Error())
	}


	if resp.StatusCode != 200 {
		t.Errorf("Status code was %d, expected 200", resp.StatusCode)
	}
}

