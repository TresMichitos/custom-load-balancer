// Tests for pool.go

package serverpool

import (
	"net/http/httptest"
	"testing"
)

// Setup test server pool and forward requests
func TestServerPoolFunctionality (t *testing.T) {

	// Initialise server pool
	var urls []string = []string {
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	var serverPool *ServerPool = NewServerPool(urls)

	// Use each server node to forward HTTP request
	for _, serverNode := range serverPool.Pool {	
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		w := httptest.NewRecorder()
		serverNode.ForwardRequest(w, r)
	}
}

