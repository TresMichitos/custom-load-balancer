// Implementation of load balancer server struct and definition of LbAlgorithm interface

package serverpool

import (
	"encoding/json"
	"net/http"
)

// Interface for load balancing algorithms
type LbAlgorithm interface {
	NextServerNode(*ServerPool, *http.Request) *ServerNode
}

// Struct to represent load balancer server
type Server struct {
	ServerPool  *ServerPool
	LbAlgorithm LbAlgorithm
}

// Struct to represent JSON serverNodeMetrics object
type ServerNodeMetrics struct {
	URL            string `json:"url"`
	RequestCount   int64  `json:"requestCount"`
	SuccessCount   int64  `json:"successCount"`
	FailureCount   int64  `json:"failureCount"`
	AverageLatency int64  `json:"averageLatency"`
}

// Struct to represent JSON Metrics object
type Metrics struct {
	TotalRequests     int64               `json:"totalRequests"`
	TotalSuccesses    int64               `json:"totalSuccesses"`
	TotalFailures     int64               `json:"totalFailures"`
	OverallLatency    int64               `json:"overallLatency"`
	ServerNodeMetrics []ServerNodeMetrics `json:"serverNodeMetrics"`
}

func newMetrics(serverPool *ServerPool) *Metrics {
	var metrics Metrics
	var totalLatency int64
	var totalRequests int64

	serverPool.Mu.Lock()
	defer serverPool.Mu.Unlock()

	for _, serverNode := range serverPool.Healthy {
		serverNode.mu.Lock()

		// Node avg latency
		var avgLatency int64
		sampleCount := len(serverNode.LatencySamples)
		for _, sample := range serverNode.LatencySamples {
			avgLatency += sample
		}

		if sampleCount > 0 {
			avgLatency /= int64(sampleCount)
			totalRequests += serverNode.RequestCount
			totalLatency += avgLatency * serverNode.RequestCount
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
			AverageLatency: avgLatency,
		})

		serverNode.mu.Unlock()
	}

	if totalRequests > 0 {
		metrics.OverallLatency = totalLatency / int64(totalRequests)
	} else {
		metrics.OverallLatency = -1
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
	server.ServerPool.Mu.Lock()
	if len(server.ServerPool.Healthy) == 0 {
		server.ServerPool.Mu.Unlock()
		http.Error(w, "Service unavailable: no healthy backend servers", http.StatusServiceUnavailable)
		return
	}

	var nextServerNode *ServerNode = server.LbAlgorithm.NextServerNode(server.ServerPool, r)
	server.ServerPool.Mu.Unlock()
	nextServerNode.ForwardRequest(w, r, server.ServerPool)
}

func (server *Server) StartLoadBalancer(enableMetrics bool) {
	http.HandleFunc("/", server.requestHandler)
	if enableMetrics {
		http.HandleFunc("/metrics", server.metricsHandler)
	}
	http.ListenAndServe(":8080", nil)
}
