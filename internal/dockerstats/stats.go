// Checks instance of sever's utilisation statistics

package dockerstats

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// Container struct for containing data
type ContainerStats struct {
	Name    string
	CPUPerc float64
	MemPerc float64
}

var (
	statsCache   = make(map[string]ContainerStats)
	cacheMutex   sync.RWMutex
	dockerClient *client.Client
)

// Polls stats every inter val and updates cache
func StartStatsPolling(dockerClientInput *client.Client, interval time.Duration) {
	dockerClient = dockerClientInput
	go func() {
		for {
			stats, err := fetchDockerStats(dockerClient)

			// If fetching was successful, update the cache
			if err == nil {
				cacheMutex.Lock()
				statsCache = stats
				cacheMutex.Unlock()
			}
			time.Sleep(interval)
		}
	}()
}

func fetchDockerStats(dockerClient *client.Client) (map[string]ContainerStats, error) {
	// Fetch the list of containers with server role
	ctx := context.Background()
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "role=server")),
	})

	if err != nil {
		return nil, err
	}

	stats := make(map[string]ContainerStats)
	var wg sync.WaitGroup
	statsMutex := sync.Mutex{}

	for _, container := range containers {
		wg.Add(1)
		go func(container types.Container) {
			defer wg.Done()

			// Fetch stats for container
			// boolean false means we don't want to stream stats
			response, err := dockerClient.ContainerStats(ctx, container.ID, false)

			if err != nil {
				return
			}

			// Decode the stat data
			decoder := json.NewDecoder(response.Body)

			var stat types.StatsJSON
			if err := decoder.Decode(&stat); err != nil {
				response.Body.Close()
				return
			}

			// Close the response body
			response.Body.Close()

			cpuContainerDelta := float64(stat.CPUStats.CPUUsage.TotalUsage - stat.PreCPUStats.CPUUsage.TotalUsage)
			cpuSystemDelta := float64(stat.CPUStats.SystemUsage - stat.PreCPUStats.SystemUsage)

			cpuPercent := 0.0

			// Calculate CPU percentage (Usage)
			if (cpuContainerDelta > 0.0) && (cpuSystemDelta > 0.0) {
				cpuPercent = (cpuContainerDelta / cpuSystemDelta) * float64(len(stat.CPUStats.CPUUsage.PercpuUsage)) * 100
			}

			memPercent := 0.0

			// Calculate Memory percentage (Usage)
			if stat.MemoryStats.Limit > 0 {
				memPercent = float64(stat.MemoryStats.Usage) / float64(stat.MemoryStats.Limit) * 100
			}

			statsMutex.Lock()
			stats[container.Names[0][1:]] = ContainerStats{
				Name:    container.Names[0][1:],
				CPUPerc: cpuPercent,
				MemPerc: memPercent,
			}
			statsMutex.Unlock()
		}(container)
	}
	wg.Wait()
	return stats, nil
}

// Returns cached Docker stats
func GetDockerStats() (map[string]ContainerStats, error) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return statsCache, nil
}
