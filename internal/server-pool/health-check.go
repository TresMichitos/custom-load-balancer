/*
 * Manages health checks for the server pool
 */

package serverpool

import (
	"fmt"
	"net/http"
	"time"
)

const TIMEOUT = 5 // Timeout in seconds

func EvalServerPool(pool []*ServerNode, endpoint string) ([]*ServerNode, []*ServerNode) {
	healthyPool := []*ServerNode{}
	unhealthyPool := []*ServerNode{}

	for _, serverNode := range pool {
		isHealthy := checkServerHealth(serverNode, endpoint)

		if isHealthy {
			healthyPool = append(healthyPool, serverNode)
		} else {
			unhealthyPool = append(unhealthyPool, serverNode)
		}
	}

	return healthyPool, unhealthyPool
}

func checkServerHealth(server *ServerNode, endpoint string) bool {
	httpClient := http.Client{Timeout: TIMEOUT * time.Second}

	req, err := http.NewRequest(http.MethodGet, server.URL+endpoint, nil)
	if err != nil {
		fmt.Printf("[%v] build-req error: %v\n", server.URL, err)
		return false
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("[%v] res error: %v\n", server.URL, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true
	}
	fmt.Printf("[%v] unhealthy status code: %d\n", server.URL, resp.StatusCode)
	return false
}
