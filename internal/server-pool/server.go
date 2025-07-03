// Implementation of load balancer server struct and definition of LbAlgorithm interface

package serverpool

import (
	"net/http"
)

// Interface for load balancing algorithms
type LbAlgorithm interface {
	NextServerNode (*ServerPool) *ServerNode
}

// Struct to represent load balancer server
type Server struct {
	ServerPool *ServerPool
	LbAlgorithm LbAlgorithm
}

// Handler function to route HTTP request using balancing algorithm
func (server *Server) requestHandler (w http.ResponseWriter, r *http.Request) {
	var nextServerNode *ServerNode = server.LbAlgorithm.NextServerNode(server.ServerPool)
	nextServerNode.ForwardRequest(w, r)
}

func (server *Server) StartLoadBalancer () {
	http.HandleFunc("/", server.requestHandler)
	http.ListenAndServe(":8080", nil)
}

