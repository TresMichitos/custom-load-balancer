// IP Hashing load balancing algorithm

package lbalgorithms

import (
	"net/http"

	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

type ipHashing struct{}

func NewIpHashing() *ipHashing {
	return &ipHashing{}
}

func (ipHashing *ipHashing) NextServerNode(serverPool *serverpool.ServerPool, req *http.Request) *serverpool.ServerNode {
	// Server health check
	if len(serverPool.Healthy) == 1 {
		return serverPool.Healthy[0]
	}

	ipClient := req.Header.Get("X-Forwarded-For")

	if ipClient == "" {
		return serverPool.Healthy[0]
	}

}
