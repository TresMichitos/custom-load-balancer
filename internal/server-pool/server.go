// Implementation of load balancer server struct

package serverpool

import (
	"fmt"
	"net/http"
)

// Struct to represent load balancer server
type Server struct {}

// Handler function to send HTTP requests
func requestHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Helloo")
}

func (server Server) StartLoadBalancer () {
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(":8080", nil)
}

