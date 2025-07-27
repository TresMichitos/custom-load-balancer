// Implementation of load balancer server struct and definition of LbAlgorithm interface

package serverpool

import (
	"encoding/json"
	"math"
	"net/http"
	"slices"
	"time"
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
	P1Latency      int64  `json:"p1Latency"`
	P5Latency      int64  `json:"p5Latency"`
	P95Latency     int64  `json:"p95Latency"`
	P99Latency     int64  `json:"p99Latency"`
}

// Struct to represent JSON Metrics object
type Metrics struct {
	Algorithm            string              `json:"algorithm"`
	Timestamp            time.Time           `json:"timestamp"`
	TotalRequests        int64               `json:"totalRequests"`
	TotalSuccesses       int64               `json:"totalSuccesses"`
	TotalFailures        int64               `json:"totalFailures"`
	OverallLatency       int64               `json:"overallLatency"`
	DistributionFairness float64             `json:"distributionFairness"`
	ServerNodeMetrics    []ServerNodeMetrics `json:"serverNodeMetrics"`
}

func newMetrics(serverPool *ServerPool) *Metrics {
	var metrics Metrics
	var totalLatency int64
	var totalRequests int64

	serverPool.mu.Lock()
	defer serverPool.mu.Unlock()

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

		// Latency percentiles
		p1, p5, p95, p99 := serverNode.calculateLatencyPercentiles()

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
			P1Latency:      p1,
			P5Latency:      p5,
			P95Latency:     p95,
			P99Latency:     p99,
		})

		serverNode.mu.Unlock()
	}

	if totalRequests > 0 {
		metrics.OverallLatency = totalLatency / int64(totalRequests)
		metrics.DistributionFairness = serverPool.calculateDistributionFairness()
	} else {
		metrics.OverallLatency = -1
	}

	return &metrics
}

func (server *ServerNode) calculateLatencyPercentiles() (p1, p5, p95, p99 int64) {
	if len(server.LatencySamples) == 0 {
		return 0, 0, 0, 0
	}

	// Sort for perecentile calcs
	sorted := make([]int64, len(server.LatencySamples))
	copy(sorted, server.LatencySamples)
	slices.Sort(sorted)

	p1 = sorted[len(sorted)*1/100]
	p5 = sorted[len(sorted)*5/100]
	p95 = sorted[len(sorted)*95/100]
	p99 = sorted[len(sorted)*99/100]

	return p1, p5, p95, p99
}

func (pool *ServerPool) calculateDistributionFairness() float64 {
	if len(pool.All) == 1 {
		return 0.0
	}

	// Mean and total requests
	var totalRequests, sumDifferential float64
	for _, server := range pool.All {
		totalRequests += float64(server.RequestCount)
	}
	meanRequests := totalRequests / float64(len(pool.All))

	// Sum of differentials squared
	for _, server := range pool.All {
		diff := float64(server.RequestCount) - meanRequests
		sumDifferential += diff * diff
	}

	// Standard deviation
	return math.Sqrt(float64(sumDifferential) / float64(len(pool.All)))
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

	var nextServerNode *ServerNode = server.LbAlgorithm.NextServerNode(server.ServerPool, r)
	server.ServerPool.mu.Unlock()
	nextServerNode.ForwardRequest(w, r, server.ServerPool)
}

func (server *Server) StartLoadBalancer(enableMetrics bool) {
	http.HandleFunc("/", server.requestHandler)
	if enableMetrics {
		http.HandleFunc("/metrics", server.metricsHandler)
	}
	http.ListenAndServe(":8080", nil)
}
