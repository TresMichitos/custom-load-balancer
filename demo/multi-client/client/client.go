package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	ipgen "github.com/TresMichitos/custom-load-balancer/demo/multi-client/ipgen"
)

type reply struct {
	Hostname  string `json:"hostname"`
	Port      string `json:"port"`
	Timestamp string `json:"timestamp"`
}

func SimulateClient(client http.Client, url string, requestCount int, INTERVAL int, clientID int) {
	ip := ipgen.GenTestNet3()

	for i := 1; i <= requestCount; i++ {
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
		fmt.Printf("[Client %02d - %02d]: host=%s port=%s ts=%s\n", clientID, i, r.Hostname, r.Port, r.Timestamp)
		time.Sleep(time.Duration(INTERVAL) * time.Millisecond)
	}
}
