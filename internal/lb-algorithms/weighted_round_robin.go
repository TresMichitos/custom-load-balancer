// Round Robin load balancing algorithm
// NOTE: If servers go down or come back online, index alignment will change.
// Consider building a stable server list internally if stronger guarantees are needed.

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

func (wrr *weightedRoundRobin) GetName() string {
	return "leastConnections"
}

// Select next server node according to weight/usage and health state
func (wrr *weightedRoundRobin) NextServerNode(serverPool *serverpool.ServerPool, _ *http.Request) *serverpool.ServerNode {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	if len(serverPool.Healthy) == 1 {
		return serverPool.Healthy[0]
	}

	// Set initial active server
	if wrr.activeServer == nil {
		wrr.activeServer = serverPool.Healthy[0]
	}

	// If node unhealthy or usage greater than weight move to next server
	if wrr.currentUsage >= wrr.activeServer.Weight || !slices.Contains(serverPool.Healthy, wrr.activeServer) {
		wrr.index = (wrr.index + 1) % len(serverPool.Healthy)
		wrr.currentUsage = 0
		wrr.activeServer = serverPool.Healthy[wrr.index]
	}

	wrr.currentUsage++
	return wrr.activeServer
}
