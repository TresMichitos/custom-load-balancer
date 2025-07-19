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
	RequestCount      int
	SuccessCount      int
	FailureCount      int
	Latency           int   // Most recent request
	LatencySamples    []int // Rolling buffer
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
func (serverNode *ServerNode) ForwardRequest(w http.ResponseWriter, r *http.Request) {
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
	elapsedTime := int(time.Since(startTime).Milliseconds())

	serverNode.mu.Lock()
	defer serverNode.mu.Unlock()

	serverNode.ActiveConnections--
	serverNode.Latency = elapsedTime
	serverNode.LatencySamples = append(serverNode.LatencySamples, elapsedTime)
	if len(serverNode.LatencySamples) > 100 {
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
	All     []*ServerNode
	Healthy []*ServerNode
	mu      sync.Mutex
}

// Factory function to initialise a new ServerPool object
func NewServerPool(urls []string) *ServerPool {
	var serverPool ServerPool
	for _, url := range urls {
		newServerNode, err := NewServerNode(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		serverPool.All = append(serverPool.All, newServerNode)
	}
	return &serverPool
}
