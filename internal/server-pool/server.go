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
	URL          string `json:"url"`
	RequestCount int    `json:"requestCount"`
}

// Struct to represent JSON Metrics object
type Metrics struct {
	AverageLatency     string              `json:"averageLatency"`
	ServerNodesMetrics []ServerNodeMetrics `json:"serverNodesMetrics"`
}

func newMetrics(serverPool *ServerPool) *Metrics {
	serverPool.mu.Lock()
	defer serverPool.mu.Unlock()

	var metrics = Metrics{}
	var totalLatency int = 0

	for _, serverNode := range serverPool.Healthy {
		serverNode.mu.Lock()
		totalLatency += serverNode.Latency
		metrics.ServerNodesMetrics = append(metrics.ServerNodesMetrics, ServerNodeMetrics{
			serverNode.URL,
			serverNode.RequestCount})
		serverNode.mu.Unlock()
	}

	metrics.AverageLatency = fmt.Sprintf("%dms", totalLatency/len(serverPool.Healthy))

	return &metrics
}

// Handler function to provide load balancer metrics
func (server *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	var metrics Metrics = *newMetrics(server.ServerPool)
	json.NewEncoder(w).Encode(metrics)
}

// Handler function to route HTTP request using balancing algorithm
func (server *Server) requestHandler(w http.ResponseWriter, r *http.Request) {
	if len(server.ServerPool.Healthy) == 0 {
		http.Error(w, "Service unavailable: no healthy backend servers", http.StatusServiceUnavailable)
		return
	}

	var nextServerNode *ServerNode = server.LbAlgorithm.NextServerNode(server.ServerPool)
	nextServerNode.ForwardRequest(w, r)
}

func (server *Server) StartLoadBalancer() {
	go HealthCheckLoop(server.ServerPool)
	http.HandleFunc("/", server.requestHandler)
	http.HandleFunc("/metrics", server.metricsHandler)
	http.ListenAndServe(":8080", nil)
}
