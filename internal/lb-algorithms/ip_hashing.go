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

	// Extract client psuedo ip from header
	ipClient := req.Header.Get("X-Forwarded-For")
	// you can fill this in with port and combine ip and port for greater uniqueness

	if (len(serverPool.Healthy) == 1) || (ipClient == "") {
		return serverPool.Healthy[0]
	}

	// Hashing client ip using FNV-1a a non cryptographic hashing algorithm to priortise speed
	hashAlgorithm := fnv.New64a()
	hashAlgorithm.Write([]byte(ipClient))
	value := hashAlgorithm.Sum64()

	// using modulo on the value creates the index from possible server pool
	serverIndex := int(value % uint64(len(serverPool.Healthy)))

	// assign client to sever here
	return serverPool.Healthy[serverIndex]

}
