// Implementation of load balancer server struct and definition of LbAlgorithm interface

package serverpool

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Interface for load balancing algorithms
type LbAlgorithm interface {
	NextServerNode(*ServerPool) *ServerNode
}

// Struct to represent load balancer server
type Server struct {
	ServerPool  *ServerPool
	LbAlgorithm LbAlgorithm
}

// Struct to represent JSON serverNodeMetrics object
type ServerNodeMetrics struct {
	URL            string `json:"url"`
	RequestCount   int    `json:"requestCount"`
	SuccessCount   int    `json:"successCount"`
	FailureCount   int    `json:"failureCount"`
	AverageLatency string `json:"averageLatency"`
}

// Struct to represent JSON Metrics object
type Metrics struct {
	TotalRequests     int                 `json:"totalRequests"`
	TotalSuccesses    int                 `json:"totalSuccesses"`
	TotalFailures     int                 `json:"totalFailures"`
	OverallLatency    string              `json:"overallLatency"`
	ServerNodeMetrics []ServerNodeMetrics `json:"serverNodeMetrics"`
}

func newMetrics(serverPool *ServerPool) *Metrics {
	var metrics Metrics
	var totalLatency int
	var totalSamples int

	serverPool.mu.Lock()
	defer serverPool.mu.Unlock()

	for _, serverNode := range serverPool.Healthy {
		serverNode.mu.Lock()

		// Node avg latency
		var avgLatency int
		sampleCount := len(serverNode.LatencySamples)
		for _, sample := range serverNode.LatencySamples {
			avgLatency += sample
		}

		if sampleCount > 0 {
			avgLatency /= sampleCount
			totalSamples++
			totalLatency += avgLatency
		}

		// Overall metrics
		metrics.TotalRequests += serverNode.RequestCount
		metrics.TotalSuccesses += serverNode.SuccessCount
		metrics.TotalFailures += serverNode.FailureCount

		metrics.ServerNodeMetrics = append(metrics.ServerNodeMetrics, ServerNodeMetrics{
			URL:            serverNode.URL,
			RequestCount:   serverNode.RequestCount,
			SuccessCount:   serverNode.SuccessCount,
			FailureCount:   serverNode.FailureCount,
			AverageLatency: fmt.Sprintf("%dms", avgLatency),
		})

		serverNode.mu.Unlock()
	}

	if totalSamples > 0 {
		metrics.OverallLatency = fmt.Sprintf("%dms", totalLatency/totalSamples)
	} else {
		metrics.OverallLatency = "N/A: no samples"
	}

	return &metrics
}

// Handler function to provide load balancer metrics
func (server *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := newMetrics(server.ServerPool)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// Handler function to route HTTP request using balancing algorithm
func (server *Server) requestHandler(w http.ResponseWriter, r *http.Request) {
	server.ServerPool.mu.Lock()
	if len(server.ServerPool.Healthy) == 0 {
		server.ServerPool.mu.Unlock()
		http.Error(w, "Service unavailable: no healthy backend servers", http.StatusServiceUnavailable)
		return
	}
	server.ServerPool.mu.Unlock()

	server.ServerPool.mu.Lock()
	var nextServerNode *ServerNode = server.LbAlgorithm.NextServerNode(server.ServerPool)
	server.ServerPool.mu.Unlock()
	nextServerNode.ForwardRequest(w, r)
}

func (server *Server) StartLoadBalancer() {
	go HealthCheckLoop(server.ServerPool)
	http.HandleFunc("/", server.requestHandler)
	http.HandleFunc("/metrics", server.metricsHandler)
	http.ListenAndServe(":8080", nil)
}
