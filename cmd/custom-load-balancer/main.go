// Entrypoint of custom load balancer

package main

import (
	"log"

	config "github.com/TresMichitos/custom-load-balancer/internal/config"
	dockerstats "github.com/TresMichitos/custom-load-balancer/internal/dockerstats"
	lbalgorithms "github.com/TresMichitos/custom-load-balancer/internal/lb-algorithms"
	serverpool "github.com/TresMichitos/custom-load-balancer/internal/server-pool"
)

// Initialise load balancer server with configured urls and algorithm parameter
func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	dockerClient, err := dockerstats.NewDockerClient()
	if err != nil {
		log.Fatalf("Failed Docker client init: %v", err)
	}
	defer dockerClient.Close()

	dockerstats.StartStatsPolling(dockerClient, cfg.Docker.PollingInterval)

	var lbAlgorithm serverpool.LbAlgorithm

	switch cfg.LoadBalancer.Algorithm {
	case "RoundRobin":
		lbAlgorithm = lbalgorithms.NewRoundRobin()
	case "WeightedRoundRobin":
		lbAlgorithm = lbalgorithms.NewWeightedRoundRobin()
	case "LeastConnections":
		lbAlgorithm = lbalgorithms.NewLeastConnections()
	case "LeastUsedResources":
		lbAlgorithm = lbalgorithms.NewLeastUsedResources()
	case "Random":
		lbAlgorithm = lbalgorithms.NewRandom()
	case "IpHashing":
		lbAlgorithm = lbalgorithms.NewIpHashing()
	default:
		lbAlgorithm = lbalgorithms.NewRoundRobin()
	}

	serverPool, err := serverpool.NewServerPool(cfg.Servers, cfg.Metrics.LatencySamples)
	if err != nil {
		log.Fatalf("Failed to initialise server: %v", err)
	}

	serverpool.StartHealthChecking(
		serverPool,
		cfg.HealthCheck.Timeout,
		cfg.HealthCheck.Interval,
	)

	server := serverpool.Server{
		ServerPool:  serverPool,
		LbAlgorithm: lbAlgorithm,
	}
	server.StartLoadBalancer(cfg.Metrics.Enabled)
}
