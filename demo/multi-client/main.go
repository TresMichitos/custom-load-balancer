/*
 * Utility for simulating a client making requests.
 * Uses TEST-NET-3 IPv4 range.
 * Usage:
 * 		go run main.go <url> 										# Single request
 * 		go run main.go <url> <request count>						# Multiple requests
 * 		go run main.go <url> <request count> <client clount>		# Multiple requests & clients
 */

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	client "github.com/TresMichitos/custom-load-balancer/demo/multi-client/client"
)

const TIMEOUT = 5     // Timeout in seconds
const INTERVAL = 1000 // Interval between requests in ms

func main() {
	var url, requestCount, clientCount, err = readArgs()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for i := 1; i <= clientCount; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			client.SimulateClient(TIMEOUT, url, requestCount, INTERVAL, i)
		}(i)
	}

	wg.Wait()
}

func readArgs() (string, int, int, error) {
	// URL
	if len(os.Args) < 2 {
		return "", 0, 0, fmt.Errorf("usage: %s <url> <count>", os.Args[0])
	}
	url := os.Args[1]

	// Request count
	requestCount := 1
	if len(os.Args) > 2 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil || n <= 0 {
			return "", 0, 0, fmt.Errorf("invalid request count: %s", os.Args[2])
		}
		requestCount = n
	}

	// Client count
	clientCount := 1
	if len(os.Args) > 3 {
		n, err := strconv.Atoi(os.Args[3])
		if err != nil || n <= 0 {
			return "", 0, 0, fmt.Errorf("invalid client count: %s", os.Args[3])
		}
		clientCount = n
	}

	return url, requestCount, clientCount, nil
}
