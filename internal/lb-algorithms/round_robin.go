// Round Robin load balancing algorithm

package lbalgorithms

import (
	"net/http"
	"sync"

	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

// Struct to implement serverpool.LbAlgorithm interface
type roundRobin struct {
	index int
	mu    sync.Mutex
}

func NewRoundRobin() *roundRobin {
	return &roundRobin{}
}

func (rr *roundRobin) GetName() string {
	return "leastConnections"
}

// Select server node by iterating over server pool
func (rr *roundRobin) NextServerNode(serverPool *serverpool.ServerPool, _ *http.Request) *serverpool.ServerNode {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if rr.index >= len(serverPool.Healthy) {
		rr.index = 0
	}

	server := serverPool.Healthy[rr.index]
	rr.index++

	return server
}
