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

// Update healthy pool at regular interval
func HealthCheckLoop(serverPool *ServerPool) {
	for {
		healthMap := make(map[*ServerNode]bool)
		var mapMu sync.Mutex

		var wg sync.WaitGroup

		serverPool.mu.Lock()
		for _, serverNode := range serverPool.All {
			wg.Add(1)
			go func(server *ServerNode) {
				defer wg.Done()
				healthy := isServerHealthy(server)

				mapMu.Lock()
				healthMap[server] = healthy
				mapMu.Unlock()
			}(serverNode)
		}
		serverPool.mu.Unlock()

		wg.Wait()

		healthyPool := []*ServerNode{}
		for _, node := range serverPool.All {
			if healthMap[node] {
				healthyPool = append(healthyPool, node)
			}
		}

		serverPool.mu.Lock()
		serverPool.Healthy = healthyPool
		serverPool.mu.Unlock()

		time.Sleep(INTERVAL * time.Second)
	}
}

// Attempts to send request and check status
func isServerHealthy(server *ServerNode) bool {
	httpClient := http.Client{Timeout: TIMEOUT * time.Second}

	server.mu.Lock()
	var URL string = server.URL
	server.mu.Unlock()

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		fmt.Printf("[%v] build-req error: %v\n", URL, err)
		return false
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("[%v] res error: %v\n", URL, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("[%v] healthy status code: %d\n", URL, resp.StatusCode)
		return true
	}
	fmt.Printf("[%v] unhealthy status code: %d\n", URL, resp.StatusCode)
	return false
}
