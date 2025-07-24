/*
 * Simple server template
 */

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

type reply struct {
	Hostname  string `json:"hostname"`
	Port      string `json:"port"`
	Timestamp string `json:"timestamp"`
}

var simulatedLatency = flag.Duration("latency", -1*time.Millisecond, "The server's simulated latency")

var listenPort string

func handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hostname := os.Getenv("SERVER_ID")
	if hostname == "" {
		if h, err := os.Hostname(); err == nil {
			hostname = h
		}
	}

	time.Sleep(*simulatedLatency)

	_ = json.NewEncoder(w).Encode(reply{
		Hostname:  hostname,
		Port:      listenPort,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func main() {
	flag.Parse()

	// Honour $PORT if set; default to 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	listenPort = port
	addr := ":" + port

	http.HandleFunc("/", handler)

	log.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
