// Random load balancing algorithm ;)

package lbalgorithms

import (
	"math/rand"
	"net/http"

	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

type random struct{}

func NewRandom() *random {
	return &random{}
}

func (r *random) GetName() string {
	return "leastConnections"
}

func (r *random) NextServerNode(serverPool *serverpool.ServerPool, _ *http.Request) *serverpool.ServerNode {
	if len(serverPool.Healthy) == 1 {
		return serverPool.Healthy[0]
	}

	return serverPool.Healthy[rand.Intn(len(serverPool.Healthy))]
}
