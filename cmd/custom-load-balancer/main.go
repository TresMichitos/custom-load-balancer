package main

import (
	"github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

func main() {
	var server serverpool.Server
	server.StartLoadBalancer()
}

