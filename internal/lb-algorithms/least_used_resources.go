// Least resources load balancing algorithm

package lbalgorithms

import (
	"log"
	"net/http"

	dockerstats "github.com/TresMichitos/custom-load-balancer/internal/dockerstats"
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

	stats, err := dockerstats.GetDockerStats()

	if err != nil {
		log.Printf("There was an error, %v", err)
		return serverPool.Healthy[0]
	}

	var bestnode *serverpool.ServerNode
	var lowestScore float64 = 1000

	for _, node := range serverPool.Healthy {
		log.Printf("container nhame: %s ", node.ContainerName)

		for k := range stats {
			log.Printf("container stats: %s ", k)
		}

		stat, ok := stats[node.ContainerName]
		if !ok {
			continue
		}
		log.Printf("Container %s CPU: %.2f%%, Mem: %.2f%%", stat.Name, stat.CPUPerc, stat.MemPerc)
		cpu := stat.CPUPerc
		mem := stat.MemPerc

		// Applying 60/40 weighting
		score := (cpu * 0.6) + (mem * 0.4)
		log.Printf("score: %f", score)

		if score < lowestScore {
			bestnode = node
			lowestScore = score
			log.Printf("lowestscore: %f", lowestScore)
		}
	}

	if bestnode == nil {
		log.Printf("No suitable node found by resource stats; using fallback node")
		return serverPool.Healthy[0]
	}

	return bestnode

}
