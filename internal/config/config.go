package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Server struct {
	URL               string        `yaml:"url"`
	Weight            int           `yaml:"weight"`
	ArtificialLatency time.Duration `yaml:"artificial_latency"`
}

type Config struct {
	Server struct {
		Port    int           `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"server"`

	HealthCheck struct {
		Interval time.Duration `yaml:"interval"`
		Timeout  time.Duration `yaml:"timeout"`
	} `yaml:"health_check"`

	LoadBalancer struct {
		Algorithm string `yaml:"algorithm"`
	} `yaml:"load_balancer"`

	Servers []Server `yaml:"servers"`

	Metrics struct {
		Enabled        bool `yaml:"enabled"`
		LatencySamples int  `yaml:"latency_samples"`
	} `yaml:"metrics"`

	Docker struct {
		Enabled         bool          `yaml:"enabled"`
		PollingInterval time.Duration `yaml:"polling_interval"`
		ContainerStats  bool          `yaml:"container_stats"`
	} `yaml:"docker"`

	Clients struct {
		Timeout  time.Duration `yaml:"timeout"`
		Interval time.Duration `yaml:"interval"`
	} `yaml:"clients"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
