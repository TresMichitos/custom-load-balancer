// Least connections load balancing algorithm

package lbalgorithms

import (
	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

type leastConnections struct{}

func NewLeastConnections() *leastConnections {
	return &leastConnections{}
}

func (leastConnections *leastConnections) NextServerNode(serverPool *serverpool.ServerPool) *serverpool.ServerNode {
	if len(serverPool.Healthy) == 1 {
		return serverPool.Healthy[0]
	}

	var bestNode *serverpool.ServerNode
	for _, node := range serverPool.Healthy {
		if node.ActiveConnections == 0 {
			return node
		}
		if bestNode == nil || node.ActiveConnections < bestNode.ActiveConnections {
			bestNode = node
		}
	}

	return bestNode
}
