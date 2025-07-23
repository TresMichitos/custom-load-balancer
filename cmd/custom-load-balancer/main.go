// Entrypoint of custom load balancer

package main

import (
	"log"

	config "github.com/TresMichitos/custom-load-balancer/internal/config"
	lbalgorithms "github.com/TresMichitos/custom-load-balancer/internal/lb-algorithms"
	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

// Initialise load balancer server with configured urls and algorithm parameter
func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	urls := make([]string, len(cfg.Servers))
	for i, server := range cfg.Servers {
		urls[i] = server.URL
	}

	var lbAlgorithm serverpool.LbAlgorithm

	switch cfg.LoadBalancer.Algorithm {
	case "RoundRobin":
		lbAlgorithm = lbalgorithms.NewRoundRobin()
	case "WeightedRoundRobin":
		lbAlgorithm = lbalgorithms.NewWeightedRoundRobin([]int{1, 2, 1})
	case "LeastConnections":
		lbAlgorithm = lbalgorithms.NewLeastConnections()
	case "Random":
		lbAlgorithm = lbalgorithms.NewRandom()
	case "IpHashing":
		lbAlgorithm = lbalgorithms.NewIpHashing()
	default:
		lbAlgorithm = lbalgorithms.NewRoundRobin()
	}

	var server serverpool.Server = serverpool.Server{ServerPool: serverpool.NewServerPool(urls), LbAlgorithm: lbAlgorithm}
	server.StartLoadBalancer()
}
