// Tests for round_robin.go

package lbalgorithms

import (
	"github.com/TresMichitos/custom-load-balancer/internal/server-pool"
	"testing"
)

// Test that algorithm returns expected server nodes
func TestRoundRobin (t *testing.T) {
	var lbAlgorithm serverpool.LbAlgorithm = NewRoundRobin()

	// For Round Robin this is also the expected order output
	var urls []string = []string {
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	var serverPool *serverpool.ServerPool = serverpool.NewServerPool(urls)
	
	for _, url := range urls {
		var nextServerNode *serverpool.ServerNode = lbAlgorithm.NextServerNode(serverPool)
		if url != nextServerNode.URL {
			t.Errorf("Routed url was %s, expected %s", nextServerNode.URL, url)
		}
	}
}

