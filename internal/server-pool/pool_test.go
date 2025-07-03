// Tests for server.go

package serverpool

import (
	"net/http/httptest"
	"testing"
)

// Setup test server pool and forward requests
func TestServerPoolFunctionality (t *testing.T) {

	// Initialise server pool
	var urls [3]string = [3]string {
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	var serverPool ServerPool = ServerPool{}
	for _, url := range urls {
		serverPool.Pool = append(serverPool.Pool, NewServerNode(url))
	}

	// Use each server node to forward HTTP request
	for _, serverNode := range serverPool.Pool {	
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		w := httptest.NewRecorder()
		serverNode.ForwardRequest(w, r)
		/*
		response := w.Result()

		if response.StatusCode != 200 {
			t.Errorf("Status Code was %d, expected 200", response.StatusCode)
		}
		*/
	}
}

