// Round Robin load balancing algorithm

package lbalgorithms

import (
	"sync"
	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

// Struct to implement serverpool.LbAlgorithm interface
type weightedRoundRobin struct {
	weightRatio []int // Slice of the weight ratio for each serverNode at corresponding
	// index in serverPool
	index               int
	weightRatioUseCount int // Number of times the weight ratio at index
	// has been used since index incremented
	mu sync.Mutex
}

func NewWeightedRoundRobin(weightRatio []int) *weightedRoundRobin {
	return &weightedRoundRobin{weightRatio: weightRatio, weightRatioUseCount: 1}
}

// Select server node following weight ratio
func (weightedRoundRobin *weightedRoundRobin) NextServerNode(serverPool *serverpool.ServerPool) *serverpool.ServerNode {
	defer func() { weightedRoundRobin.weightRatioUseCount++; weightedRoundRobin.mu.Unlock() }()

	weightedRoundRobin.mu.Lock()

	// Increment index if server node is about to be used more times than its ratio value
	// since index incremented
	if weightedRoundRobin.weightRatioUseCount > weightedRoundRobin.weightRatio[weightedRoundRobin.index] {
		weightedRoundRobin.index++

		if weightedRoundRobin.index >= len(serverPool.Healthy) {
			weightedRoundRobin.index = 0
		}

		weightedRoundRobin.weightRatioUseCount = 1
	}

	return serverPool.Healthy[weightedRoundRobin.index]
}
