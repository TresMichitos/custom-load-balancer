// Checks instance of sever's utilisation statistics

package dockerstats

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

// Container struct for containing data
type ContainerStats struct {
	Name    string
	CPUPerc float64
	MemPerc float64
}

func StartStatsPolling(dockerClient *client.Client, interval time.Duration) {
}

// Helper function for seperating statistic fetching and the algorithm
func GetDockerStats() (map[string]ContainerStats, error) {
	// Docker command "docker stats --no-stream --format {{json .}}"
	// this grabs that instant of server stats and formats it to json
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "{{json .}}")
	// Checks if command can be run without error
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// refernce container stats struct
	stats := make(map[string]ContainerStats)

	// Each line represents json output from a different container
	// Split lines to process each container
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		log.Printf("line: %s", line)
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// Temp struct to match docker json output
		var raw struct {
			Name    string `json:"Name"`
			CPUPerc string `json:"CPUPerc"`
			MemPerc string `json:"MemPerc"`
		}

		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			log.Printf("error in unmarshalling %s, %v", line, err)
			continue
		}

		// parse % utility strings to be float64
		stats[raw.Name] = ContainerStats{
			Name:    raw.Name,
			CPUPerc: ParsePercentage(raw.CPUPerc),
			MemPerc: ParsePercentage(raw.MemPerc),
		}
	}

	return stats, nil
}

func ParsePercentage(s string) float64 {
	s = strings.TrimSpace(strings.TrimSuffix(s, "%"))
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return val
}
