/*
 * Utility for simulating a client making requests.
 * Uses TEST-NET-3 IPv4 range.
 * Usage:
 * 		go run client.go <url>			# Single request
 *		go run client.go <url> <count>	# <count> requests
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	ipgen "github.com/TresMichitos/custom-load-balancer/demo/client/ipgen"
)

type reply struct {
	Hostname  string `json:"hostname"`
	Port      string `json:"port"`
	Timestamp string `json:"timestamp"`
}

const TIMEOUT = 5    // Timeout in seconds
const INTERVAL = 300 // Interval between requests in ms

func main() {
	var url, count, err = readArgs()
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{Timeout: TIMEOUT * time.Second}
	ip := ipgen.GenTestNet3()

	// Request loop
	for i := 1; i <= count; i++ {
		// Build req
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf("[%02d] build-req error: %v\n", i, err)
			continue
		}
		req.Header.Set("X-Forwarded-For", ip)

		// Send req
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[%02d] request error: %v\n", i, err)
			continue
		}

		// Parse response
		var r reply
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			fmt.Printf("[%02d]: Decode error: %v\n", i, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// Log and interval
		fmt.Printf("[%02d]: host=%s port=%s ts=%s\n", i, r.Hostname, r.Port, r.Timestamp)
		time.Sleep(INTERVAL * time.Millisecond)
	}
}

func readArgs() (string, int, error) {
	// URL
	if len(os.Args) < 2 {
		return "", 0, fmt.Errorf("usage: %s <url> <count>", os.Args[0])
	}

	url := os.Args[1]

	// Request count
	count := 1
	if len(os.Args) > 2 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil || n <= 0 {
			return "", 0, fmt.Errorf("invalid count: %s", os.Args[2])
		}
		count = n
	}

	return url, count, nil
}
