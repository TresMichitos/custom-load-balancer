// Tests for round_robin.go

package lbalgorithms

import (
	"net/http"
	"testing"

	config "github.com/TresMichitos/custom-load-balancer/internal/config"
	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

// Test that algorithm returns expected server nodes
func TestRoundRobin(t *testing.T) {
	var lbAlgorithm serverpool.LbAlgorithm = NewRoundRobin()

	// For Round Robin this is also the expected order output
	var urls []string = []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	servers := make([]config.Server, 0, 3)
	for _, url := range urls {
		newServer := config.Server{
			URL: url,
		}
		servers = append(servers, newServer)
	}

	var serverPool *serverpool.ServerPool = serverpool.NewServerPool(servers, 1)
	serverPool.Healthy = serverPool.All

	// Dummy request to satisfy params
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	for _, url := range urls {
		var nextServerNode *serverpool.ServerNode = lbAlgorithm.NextServerNode(serverPool, req)
		if url != nextServerNode.URL {
			t.Errorf("Routed url was %s, expected %s", nextServerNode.URL, url)
		}
	}
}
