// Entrypoint of custom load balancer

package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	lbalgorithms "github.com/TresMichitos/custom-load-balancer/internal/lb-algorithms"
	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

// Load urls from servers.json
func parse_addresses() ([]string, error) {
	var address_arr []string

	b, err := os.ReadFile("servers.json")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := json.Unmarshal(b, &address_arr); err != nil {
		return nil, err
	}

	return address_arr, nil
}

// Initialise load balancer server with configured urls and algorithm parameter
func main() {

	// Set algorithm flag
	var lbAlgorithmChoice *string = flag.String("algorithm", "RoundRobin", "The Load Balancing algorithm to use")

	// Parse given flags
	flag.Parse()

	var lbAlgorithm serverpool.LbAlgorithm

	switch *lbAlgorithmChoice {
	case "RoundRobin":
		lbAlgorithm = lbalgorithms.NewRoundRobin()
	case "WeightedRoundRobin":
		lbAlgorithm = lbalgorithms.NewWeightedRoundRobin([]int{1, 2, 1})
	case "LeastConnections":
		lbAlgorithm = lbalgorithms.NewLeastConnections()
	case "Random":
		lbAlgorithm = lbalgorithms.NewRandom()
	default:
		lbAlgorithm = lbalgorithms.NewRoundRobin()
	}

	urls, err := parse_addresses()
	if err != nil {
		log.Fatal(err)
	}

	var server serverpool.Server = serverpool.Server{ServerPool: serverpool.NewServerPool(urls), LbAlgorithm: lbAlgorithm}
	server.StartLoadBalancer()
}
