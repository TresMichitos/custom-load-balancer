// Tests for weighted_round_robin.go

package lbalgorithms

import (
	"github.com/TresMichitos/custom-load-balancer/internal/server-pool"
	"testing"
)

// Test that algorithm returns expected server nodes
func TestWeightedRoundRobin (t *testing.T) {
	var lbAlgorithm serverpool.LbAlgorithm = NewWeightedRoundRobin([]int {1, 2, 1})

	var urls []string = []string {
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	var serverPool *serverpool.ServerPool = serverpool.NewServerPool(urls)

	var expectedUrlsRouted []string = []string {
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	
	for _, url := range expectedUrlsRouted {
		var nextServerNode *serverpool.ServerNode = lbAlgorithm.NextServerNode(serverPool)
		if url != nextServerNode.URL {
			t.Errorf("Routed url was %s, expected %s", nextServerNode.URL, url)
		}
	}
}

