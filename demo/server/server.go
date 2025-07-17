package main

import (
	"encoding/json"
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

const SIMULATED_LATENCY = 6000

var listenPort string

func handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hostname := os.Getenv("SERVER_ID")
	if hostname == "" {
		if h, err := os.Hostname(); err == nil {
			hostname = h
		}
	}

	time.Sleep(SIMULATED_LATENCY * time.Millisecond)

	_ = json.NewEncoder(w).Encode(reply{
		Hostname:  hostname,
		Port:      listenPort,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func main() {
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
