// Least resources load balancing algorithm

package lbalgorithms

import (
	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

type leastUsedResources struct{}

func NewLeastUsedResources() *leastUsedResources {
	return &leastUsedResources{}
}

func (leastUsedResources *leastUsedResources) NextServerNode(serverPool *serverpool.ServerPool) *serverpool.ServerNode {
	// Server health check
	if len(serverPool.Healthy) == 1 {
		return serverPool.Healthy[0]
	}

}
