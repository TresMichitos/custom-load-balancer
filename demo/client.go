/*
 * Usage:
 * 		go run client.go <url>			# Single request
 *		go run client.go <url> <count>	# <count> requests
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type reply struct {
	Hostname  string `json:"hostname"`
	Port      string `json:"port"`
	Timestamp string `json:"timestamp"`
}

const TIMEOUT = 5    // Timeout in seconds
const INTERVAL = 300 // Interval between requests in ms

func main() {
	// Target URL
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <url> <count>\n", os.Args[0])
		os.Exit(1)
	}
	url := os.Args[1]

	// Number of requests
	count := 1
	if len(os.Args) > 2 {
		if n, err := strconv.Atoi(os.Args[2]); err == nil && n > 0 {
			count = n
		}
	}

	client := http.Client{Timeout: TIMEOUT * time.Second}

	// Request loop
	for i := 1; i <= count; i++ {
		// Send GET request to URL
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("[%2d]: %v\n", i, err)
			continue
		}

		// Decode response body as `reply` struct
		var r reply
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			fmt.Printf("[%2d]: Decode error: %v\n", i, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// Log and interval
		fmt.Printf("[%2d]: host=%s port=%s ts=%s\n", i, r.Hostname, r.Port, r.Timestamp)
		time.Sleep(INTERVAL * time.Millisecond)
	}
}
