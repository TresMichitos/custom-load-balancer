// Round Robin load balancing algorithm

package lbalgorithms

import (
	"github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

// Struct to implement serverpool.LbAlgorithm interface
type WeightedRoundRobin struct {
	WeightRatio []int  // Slice of the weight ratio for each serverNode at corresponding
					   // index in serverPool
	index int
	weightRatioUseCount int  // Number of times the weight ratio at index
						     // has been used since index incremented
}

// Select server node following weight ratio
func (weightedRoundRobin *WeightedRoundRobin) NextServerNode (serverPool *serverpool.ServerPool) *serverpool.ServerNode {
	defer func () {weightedRoundRobin.weightRatioUseCount ++} ()

	// Increment index if server node is about to be used more times than its ratio value
	// since index incremented (considering that weightRatioUseCount starts at zero)
	if weightedRoundRobin.weightRatioUseCount >= weightedRoundRobin.WeightRatio[weightedRoundRobin.index] {
		weightedRoundRobin.index ++

		if weightedRoundRobin.index >= len(serverPool.Pool) {
			weightedRoundRobin.index = 0
		}

		weightedRoundRobin.weightRatioUseCount = 0
	}

	return serverPool.Pool[weightedRoundRobin.index]
} 

