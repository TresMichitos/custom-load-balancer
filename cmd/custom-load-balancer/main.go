// Entrypoint of custom load balancer

package main

import (
	"github.com/TresMichitos/custom-load-balancer/internal/lb-algorithms"
	"github.com/TresMichitos/custom-load-balancer/internal/server-pool"
	"flag"
)

// Initialise load balancer server with configured urls and algorithm parameter
func main() {

	// Set algorithm flag
	var lbAlgorithmChoice *string = flag.String("algorithm", "RoundRobin", "The Load Balancing algorithm to use")

	// Parse given flags
	flag.Parse()

	var lbAlgorithm serverpool.LbAlgorithm

	switch *lbAlgorithmChoice {
		default:
			lbAlgorithm = &lbalgorithms.RoundRobin{}
	}

	// TODO: Load urls from config file
	var urls []string = []string {
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	var server serverpool.Server = serverpool.Server{ServerPool: serverpool.NewServerPool(urls), LbAlgorithm: lbAlgorithm}
	server.StartLoadBalancer()
}

