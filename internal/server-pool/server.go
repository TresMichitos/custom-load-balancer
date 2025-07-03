package serverpool

import (
	"fmt"
	"net/http"
)

type Server struct {}

func requestHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Helloo")
}

func (server Server) StartLoadBalancer () {
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(":8080", nil)
}

