// Tests for weighted_round_robin.go

package lbalgorithms

import (
	"net/http"
	"strconv"
	"testing"

	config "github.com/TresMichitos/custom-load-balancer/internal/config"
	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

// Test that algorithm returns expected server nodes
func TestWeightedRoundRobin(t *testing.T) {
	weights := []int{1, 2, 1}
	var lbAlgorithm serverpool.LbAlgorithm = NewWeightedRoundRobin(weights)

	servers := make([]config.Server, 0, 3)
	for i := 1; i <= 3; i++ {
		newServer := config.Server{
			URL: "http://localhost:808" + strconv.Itoa(i),
		}
		servers = append(servers, newServer)
	}

	var serverPool *serverpool.ServerPool = serverpool.NewServerPool(servers, 1)
	serverPool.Healthy = serverPool.All

	var expectedUrlsRouted []string = []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	// Dummy request to satisfy params
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	for _, url := range expectedUrlsRouted {
		var nextServerNode *serverpool.ServerNode = lbAlgorithm.NextServerNode(serverPool, req)
		if url != nextServerNode.URL {
			t.Errorf("Routed url was %s, expected %s", nextServerNode.URL, url)
		}
	}
}
