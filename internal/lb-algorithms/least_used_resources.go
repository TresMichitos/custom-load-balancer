// Least resources load balancing algorithm

package lbalgorithms

import (
	"net/http"

	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

type leastUsedResources struct{}

func NewLeastUsedResources() *leastUsedResources {
	return &leastUsedResources{}
}

func (leastUsedResources *leastUsedResources) NextServerNode(serverPool *serverpool.ServerPool, _ *http.Request) *serverpool.ServerNode {
	if len(serverPool.Healthy) == 1 {
		return serverPool.Healthy[0]
	}

	stats, err := serverpool.GetDockerStats()
	if err != nil {
		return serverPool.Healthy[0]
	}

	var bestnode *serverpool.ServerNode
	var lowestScore float64 = 1000

	for _, node := range serverPool.Healthy {

		stat, ok := stats[node.ContainerName]
		if !ok {
			continue
		}

		cpu := stat.CPUPerc
		mem := stat.MemPerc

		// Applying 60/40 weighting
		score := (cpu * 0.6) + (mem * 0.4)

		if score < lowestScore {
			bestnode = node
			lowestScore = score
		}
	}

	if bestnode == nil {
		return serverPool.Healthy[0] // fallback
	}

	return bestnode

}
