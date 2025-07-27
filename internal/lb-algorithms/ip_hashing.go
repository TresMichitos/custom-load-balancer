// IP Hashing load balancing algorithm

package lbalgorithms

import (
	"hash/fnv"
	"net/http"

	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

type ipHashing struct{}

func NewIpHashing() *ipHashing {
	return &ipHashing{}
}

func (ipHashing *ipHashing) NextServerNode(serverPool *serverpool.ServerPool, req *http.Request) *serverpool.ServerNode {
	ipClient := req.Header.Get("X-Forwarded-For")

	if (len(serverPool.Healthy) == 1) || (ipClient == "") {
		return serverPool.Healthy[0]
	}

	// Hashing client ip using FNV-1a
	hashAlgorithm := fnv.New64a()
	hashAlgorithm.Write([]byte(ipClient))
	value := hashAlgorithm.Sum64()

	// assign hash value to server index
	serverIndex := int(value % uint64(len(serverPool.Healthy)))

	return serverPool.Healthy[serverIndex]

}
