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

func (random *random) NextServerNode(serverPool *serverpool.ServerPool, _ *http.Request) *serverpool.ServerNode {
	serverPool.Mu.Lock()
	defer serverPool.Mu.Unlock()

	if len(serverPool.Healthy) == 1 {
		return serverPool.Healthy[0]
	}

	return serverPool.Healthy[rand.Intn(len(serverPool.Healthy))]
}
