// Weighted Round Robin load balancing algorithm

package lbalgorithms

import (
	"net/http"
	"slices"
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

// Select next server node according to weight/usage and health state
func (wrr *weightedRoundRobin) NextServerNode(serverPool *serverpool.ServerPool, _ *http.Request) *serverpool.ServerNode {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	// Set initial active server to first healthy server in pool
	if wrr.activeServer == nil {
		for {
			wrr.activeServer = serverPool.All[wrr.index]
			if !slices.Contains(serverPool.Healthy, wrr.activeServer) {
				wrr.index = (wrr.index + 1) % len(serverPool.All)
			} else {
				break
			}
		}
	}

	// If node unhealthy or usage greater than weight move to next healthy server in pool
	if !slices.Contains(serverPool.Healthy, wrr.activeServer) || wrr.currentUsage >= wrr.activeServer.Weight {

		// Find next node in pool that's healthy
		wrr.index = (wrr.index + 1) % len(serverPool.All)
		wrr.activeServer = serverPool.All[wrr.index]
		for !slices.Contains(serverPool.Healthy, wrr.activeServer) {
			wrr.index = (wrr.index + 1) % len(serverPool.All)
			wrr.activeServer = serverPool.All[wrr.index]
		}
		wrr.currentUsage = 0
	}

	wrr.currentUsage++
	return wrr.activeServer
}
