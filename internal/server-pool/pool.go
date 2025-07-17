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
	Latency           int
	mu                sync.Mutex
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
	serverNode.mu.Lock()
	defer serverNode.mu.Unlock()

	serverNode.ActiveConnections += 1
	serverNode.RequestCount++
	startTime := time.Now()
	serverNode.ReverseProxy.ServeHTTP(w, r)
	serverNode.Latency = int(time.Since(startTime).Milliseconds())
	serverNode.ActiveConnections -= 1
}

// Struct to contain collection of server nodes
type ServerPool struct {
	Healthy   []*ServerNode
	Unhealthy []*ServerNode
	mu        sync.Mutex
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
		serverPool.Unhealthy = append(serverPool.Unhealthy, newServerNode)
	}
	return &serverPool
}
