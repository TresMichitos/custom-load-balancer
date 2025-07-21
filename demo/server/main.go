package main

import (
	"encoding/json"
	"fmt"
	"log"
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

var simulatedLatency int // In milliseconds

var listenPort string

func handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hostname := os.Getenv("SERVER_ID")
	if hostname == "" {
		if h, err := os.Hostname(); err == nil {
			hostname = h
		}
	}

	time.Sleep(time.Duration(simulatedLatency) * time.Millisecond)

	_ = json.NewEncoder(w).Encode(reply{
		Hostname:  hostname,
		Port:      listenPort,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func main() {
	// Simulated latency
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <latency (ms)>", os.Args[0])	
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n <= 0 {
		fmt.Printf("invalid latency value: %s", os.Args[1])
	}
	simulatedLatency = n

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
