// Tests for server.go

package serverpool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Handler to send HTTP requests
func testRequestHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test Helloo")
}

// Send test request to load balancer server and check response status
func TestLoadBalancerRespondsOK(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()
	testRequestHandler(w, r)
	response := w.Result()

	if response.StatusCode != 200 {
		t.Errorf("Status Code was %d, expected 200", response.StatusCode)
	}
}

