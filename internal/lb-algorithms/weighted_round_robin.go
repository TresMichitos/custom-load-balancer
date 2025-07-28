// Weighted Round Robin load balancing algorithm

package lbalgorithms

import (
	"net/http"
	"sync"

	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

// Struct to implement serverpool.LbAlgorithm interface
type weightedRoundRobin struct {
	index        int
	currentUsage int
	activeServer *serverpool.ServerNode
	mu           sync.Mutex
}

func NewWeightedRoundRobin() *weightedRoundRobin {
	return &weightedRoundRobin{}
}

func (wrr *weightedRoundRobin) GetName() string {
	return "weightedRoundRobin"
}

// Select next server node according to weight/usage and health state
func (wrr *weightedRoundRobin) NextServerNode(serverPool *serverpool.ServerPool, _ *http.Request) *serverpool.ServerNode {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	// Set initial active server
	if wrr.activeServer == nil {
		wrr.activeServer = serverPool.All[wrr.index]
	}

	// Map for O(1) subsequent lookups
	healthyMap := make(map[*serverpool.ServerNode]bool, len(serverPool.Healthy))
	for _, server := range serverPool.Healthy {
		healthyMap[server] = true
	}

	// If node unhealthy or usage greater than weight move to next healthy server in pool
	if !healthyMap[wrr.activeServer] || wrr.currentUsage >= wrr.activeServer.Weight {
		// Find next node in pool that's healthy
		for {
			wrr.index = (wrr.index + 1) % len(serverPool.All)
			wrr.activeServer = serverPool.All[wrr.index]
			if healthyMap[wrr.activeServer] {
				break
			}
		}
		wrr.currentUsage = 0
	}

	wrr.currentUsage++
	return wrr.activeServer
}
