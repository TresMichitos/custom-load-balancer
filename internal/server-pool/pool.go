package serverpool

import (
	"fmt"
	"net/http"
)

type Server struct {}

func handler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Helloo")
}

func (server Server) StartLoadBalancer () {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

