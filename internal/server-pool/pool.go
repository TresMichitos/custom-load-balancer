// Implementation of server nodes and server pool

package serverpool

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// Struct to represent each node server we forward to
type ServerNode struct {
	URL               string
	ReverseProxy      *httputil.ReverseProxy
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
func NewServerNode(urlInput string) (*ServerNode, error) {
	url, err := url.Parse(urlInput)
	if err != nil {
		return nil, errors.New("invalid URL")
	}
	return &ServerNode{URL: urlInput, ReverseProxy: httputil.NewSingleHostReverseProxy(url)}, nil
}

// Proxy function to forward HTTP request to server node
func (serverNode *ServerNode) ForwardRequest(w http.ResponseWriter, r *http.Request, serverPool *ServerPool) {
	startTime := time.Now()

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
func NewServerPool(urls []string, maxLatencySamples int) *ServerPool {
	var nodes []*ServerNode
	for _, url := range urls {
		newServerNode, err := NewServerNode(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nodes = append(nodes, newServerNode)
	}

	return &ServerPool{
		All:               nodes,
		Healthy:           []*ServerNode{},
		MaxLatencySamples: maxLatencySamples,
	}
}
