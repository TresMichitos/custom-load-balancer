/*
 * Manages health checks for the server pool
 */

package serverpool

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const TIMEOUT = 5  // Timeout in seconds
const INTERVAL = 5 // Health check interval

func HealthCheckLoop(serverPool *ServerPool) {
	for {
		healthyPool := []*ServerNode{}
		unhealthyPool := []*ServerNode{}

		pool := append(serverPool.Healthy, serverPool.Unhealthy...)

		var wg sync.WaitGroup
		for _, serverNode := range pool {
			wg.Add(1)
			go sortServer(serverNode, &healthyPool, &unhealthyPool, &wg)
		}
		wg.Wait()

		serverPool.mu.Lock()
		serverPool.Healthy, serverPool.Unhealthy = healthyPool, unhealthyPool
		serverPool.mu.Unlock()

		time.Sleep(INTERVAL * time.Second)
	}
}

func sortServer(server *ServerNode, healthyPool *[]*ServerNode, unhealthyPool *[]*ServerNode, wg *sync.WaitGroup) {
	isHealthy := isServerHealthy(server)

	if isHealthy {
		*healthyPool = append(*healthyPool, server)
	} else {
		*unhealthyPool = append(*unhealthyPool, server)
	}

	defer wg.Done()
}

func isServerHealthy(server *ServerNode) bool {
	httpClient := http.Client{Timeout: TIMEOUT * time.Second}

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
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
		fmt.Printf("[%v] healthy status code: %d\n", server.URL, resp.StatusCode)
		return true
	}
	fmt.Printf("[%v] unhealthy status code: %d\n", server.URL, resp.StatusCode)
	return false
}
