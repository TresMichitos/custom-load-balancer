package client

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	ipgen "github.com/TresMichitos/custom-load-balancer/demo/multi-client/ipgen"
)

type reply struct {
	Hostname  string `json:"hostname"`
	Port      string `json:"port"`
	Timestamp string `json:"timestamp"`
}

func SimulateClient(url string, duration time.Duration, interval time.Duration, clientID int) {
	httpClient := http.Client{}
	ip := ipgen.GenTestNet3()
	startTime := time.Now()
	nextRequest := startTime //
	var reqID int

	for {
		// Break on exceeded test duration
		if time.Since(startTime) > duration {
			break
		}
		reqID++

		// Build req
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Printf("[Client %02d - %02d] request builder error: %v\n", clientID, reqID, err)
			continue
		}
		req.Header.Set("X-Forwarded-For", ip)

		// Send req
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Printf("[Client %02d - %02d] request error: %v\n", clientID, reqID, err)
			continue
		}

		// Parse response
		var r reply
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			log.Printf("[Client %02d - %02d]: decode error: %v\n", clientID, reqID, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// Log
		log.Printf(
			"[Client %02d - %02d]: host=%s port=%s\n",
			clientID, reqID, r.Hostname, r.Port,
		)
		delay := time.Since(nextRequest)
		if delay > interval {
			log.Printf("[Client %02d - %02d] WARN: late by %v", clientID, reqID, delay)
		}

		// Wait for next request
		nextRequest = nextRequest.Add(interval)
		time.Sleep(time.Until(nextRequest))
	}
}
