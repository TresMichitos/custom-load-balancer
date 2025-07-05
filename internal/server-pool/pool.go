// Implementation of server nodes and server pool

package serverpool

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// Struct to represent each node server we forward to
type ServerNode struct {
	URL string
	ReverseProxy *httputil.ReverseProxy
	mu sync.Mutex
}

// Factory function to initialise a new ServerNode object
func NewServerNode (urlInput string) (*ServerNode, error) {
	url, err := url.Parse(urlInput)
	if err != nil {
		return nil, errors.New("Invalid URL")
	}
	return &ServerNode{URL: urlInput, ReverseProxy: httputil.NewSingleHostReverseProxy(url)}, nil
}

// Proxy function to forward HTTP request to server node
func (serverNode *ServerNode) ForwardRequest (w http.ResponseWriter, r *http.Request) {
	serverNode.ReverseProxy.ServeHTTP(w, r)
}

// Struct to contain collection of server nodes
type ServerPool struct {
	Pool []*ServerNode
	mu sync.Mutex
}

// Factory function to initialise a new ServerPool object
func NewServerPool (urls []string) *ServerPool {
	var serverPool ServerPool
	for _, url := range urls {
		newServerNode, err := NewServerNode(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		serverPool.Pool = append(serverPool.Pool, newServerNode)
	}
	return &serverPool
}

