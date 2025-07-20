// IP Hashing load balancing algorithm

package lbalgorithms

import (
	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

type ipHashing struct{}

func NewIpHashing() *ipHashing {
	return &ipHashing{}
}

func (ipHashing *ipHashing) NextServerNode(serverPool *serverpool.ServerPool) *serverpool.ServerNode {

}
