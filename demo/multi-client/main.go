/*
 * Utility for simulating a client making requests.
 * Uses TEST-NET-3 IPv4 range.
 */

package main

import (
	"flag"
	"sync"
	"time"

	client "github.com/TresMichitos/custom-load-balancer/demo/multi-client/client"
)

var (
	loadBalancerURL = flag.String("url", "http://localhost:8080", "URL of load balancer")
	numClients      = flag.Int("clients", 1, "Number of concurrent clients")
	duration        = flag.Duration("duration", 30*time.Second, "Duration of test")
	requestRate     = flag.Float64("rate", 1, "Requests per second per client")
	// outputFile      = flag.String("file", "", "CSV output file for results")
)

func main() {
	flag.Parse()
	clientInterval := time.Duration(1.0 / *requestRate * float64(time.Second)) // inverse request rate

	var wg sync.WaitGroup

	for i := 1; i <= *numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			client.SimulateClient(*loadBalancerURL, *duration, clientInterval, i)
		}(i)
	}

	wg.Wait()
}
