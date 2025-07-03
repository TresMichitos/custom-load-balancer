package serverpool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequestHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test Helloo")
}

func TestLoadBalancerRespondsOK(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()
	testRequestHandler(w, r)
	response := w.Result()

	if response.StatusCode != 200 {
		t.Errorf("Status Code was %d, expected 200", response.StatusCode)
	}
}

