// Least resources load balancing algorithm

package lbalgorithms

import (
	"math"
	"net/http"

	dockerstats "github.com/TresMichitos/custom-load-balancer/internal/dockerstats"
	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

type leastUsedResources struct{}

func NewLeastUsedResources() *leastUsedResources {
	return &leastUsedResources{}
}

func (lur *leastUsedResources) GetName() string {
	return "leastUsedResources"
}

// Selects the healthiest node based on weighted cpu and memory usage
// Using 60/40 weighting rule
func (lur *leastUsedResources) NextServerNode(serverPool *serverpool.ServerPool, _ *http.Request) *serverpool.ServerNode {
	if len(serverPool.Healthy) == 1 {
		return serverPool.Healthy[0]
	}

	stats, err := dockerstats.GetDockerStats()

	if err != nil {
		return serverPool.Healthy[0]
	}

	var bestnode *serverpool.ServerNode
	lowestScore := math.MaxFloat64

	for _, node := range serverPool.Healthy {

		stat, ok := stats[node.ContainerName]
		if !ok {
			continue
		}

		// Applying 60/40 weighting
		score := (stat.CPUPerc * 0.6) + (stat.MemPerc * 0.4)

		if score < lowestScore {
			bestnode = node
			lowestScore = score
		}
	}

	if bestnode == nil {
		return serverPool.Healthy[0]
	}

	return bestnode

}
