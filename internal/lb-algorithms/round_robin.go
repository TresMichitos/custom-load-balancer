// Round Robin load balancing algorithm

package lbalgorithms

import (
	"github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

// Struct to implement serverpool.LbAlgorithm interface
type roundRobin struct {
	index int
}

func NewRoundRobin () *roundRobin {
	return &roundRobin{}
}

// Select server node by iterating over server pool
func (roundRobin *roundRobin) NextServerNode (serverPool *serverpool.ServerPool) *serverpool.ServerNode {
	defer func () {roundRobin.index ++} ()

	if roundRobin.index >= len(serverPool.Pool) {
		roundRobin.index = 0
	}
	return serverPool.Pool[roundRobin.index]
} 

