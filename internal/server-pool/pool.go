// Implementation of server nodes and server pool

package serverpool

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	config "github.com/TresMichitos/custom-load-balancer/internal/config"
)

// Struct to represent each node server we forward to
type ServerNode struct {
	URL               string
	ReverseProxy      *httputil.ReverseProxy
	ArtificialLatency time.Duration
	ActiveConnections int
	RequestCount      int64
	SuccessCount      int64
	FailureCount      int64
	Latency           int64   // Most recent request
	LatencySamples    []int64 // Rolling buffer
	mu                sync.Mutex
}

// Response writer wrapper
type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Factory function to initialise a new ServerNode object
func NewServerNode(urlInput string, artificialLatency time.Duration, maxLatencySamples int) (*ServerNode, error) {
	url, err := url.Parse(urlInput)
	if err != nil {
		return nil, errors.New("invalid URL")
	}
	return &ServerNode{
		URL:               urlInput,
		ReverseProxy:      httputil.NewSingleHostReverseProxy(url),
		ArtificialLatency: artificialLatency,
		LatencySamples:    make([]int64, 0, maxLatencySamples),
	}, nil
}

// Proxy function to forward HTTP request to server node
func (serverNode *ServerNode) ForwardRequest(w http.ResponseWriter, r *http.Request, serverPool *ServerPool) {
	startTime := time.Now()

	time.Sleep(serverNode.ArtificialLatency)

	// Pre-metrics
	serverNode.mu.Lock()
	serverNode.ActiveConnections++
	serverNode.RequestCount++
	serverNode.mu.Unlock()

	// Forward request
	rww := &responseWriterWrapper{ResponseWriter: w, status: 200}
	serverNode.ReverseProxy.ServeHTTP(rww, r)

	// Post-metrics
	elapsedTime := int64(time.Since(startTime).Microseconds())

	serverNode.mu.Lock()
	defer serverNode.mu.Unlock()

	serverNode.ActiveConnections--
	serverNode.Latency = elapsedTime
	serverNode.LatencySamples = append(serverNode.LatencySamples, elapsedTime)
	if len(serverNode.LatencySamples) > serverPool.MaxLatencySamples {
		serverNode.LatencySamples = serverNode.LatencySamples[1:]
	}

	if rww.status >= 200 && rww.status < 300 {
		serverNode.SuccessCount++
	} else {
		serverNode.FailureCount++
	}
}

// Struct to contain collection of server nodes
type ServerPool struct {
	All               []*ServerNode
	Healthy           []*ServerNode
	mu                sync.Mutex
	MaxLatencySamples int
}

// Factory function to initialise a new ServerPool object
func NewServerPool(servers []config.Server, maxLatencySamples int) (*ServerPool, error) {
	var nodes []*ServerNode
	var errors []string

	for _, server := range servers {
		newServerNode, err := NewServerNode(server.URL, server.ArtificialLatency, maxLatencySamples)
		if err != nil {
			log.Printf("Failed to create server node for %s: %v", server.URL, err)
			errors = append(errors, fmt.Sprintf("server %s: %v", server.URL, err))
			continue
		}
		nodes = append(nodes, newServerNode)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no valid servers configured: %s", strings.Join(errors, ", "))
	}

	return &ServerPool{
		All:               nodes,
		Healthy:           []*ServerNode{},
		MaxLatencySamples: maxLatencySamples,
	}, nil
}
