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

// Select server node by iterating over server pool
func (roundRobin *roundRobin) NextServerNode(serverPool *serverpool.ServerPool, _ *http.Request) *serverpool.ServerNode {
	roundRobin.mu.Lock()
	defer roundRobin.mu.Unlock()

	serverPool.Mu.Lock()
	defer serverPool.Mu.Unlock()

	if roundRobin.index >= len(serverPool.Healthy) {
		roundRobin.index = 0
	}

	server := serverPool.Healthy[roundRobin.index]
	roundRobin.index++

	return server
}
