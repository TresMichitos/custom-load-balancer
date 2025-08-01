// This file is for connecting to the Docker API and creating a Docker client.
package dockerstats

import (
	"github.com/docker/docker/client"
)

func NewDockerClient() (*client.Client, error) {
	// Create a new Docker client with the default settings
	cli, err := client.NewClientWithOpts(
		// Uses the default Docker host and API version
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return nil, err
	}

	return cli, nil
}
