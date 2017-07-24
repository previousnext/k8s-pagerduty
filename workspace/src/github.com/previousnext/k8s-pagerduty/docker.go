package main

import (
	"fmt"
	"time"

	"github.com/fsouza/go-dockerclient"
)

// Helper function get the stats from a running container.
func stats(c *docker.Client, id string) (*docker.Stats, error) {
	errChannel := make(chan error, 1)
	statsChannel := make(chan *docker.Stats)

	go func() {
		errChannel <- c.Stats(docker.StatsOptions{
			ID:      id,
			Stats:   statsChannel,
			Stream:  false,
			Timeout: 10 * time.Second,
		})
	}()

	stats, ok := <-statsChannel
	if !ok {
		return stats, fmt.Errorf("Bad response getting stats for container: %s", id)
	}

	err := <-errChannel
	if err != nil {
		return stats, err
	}

	return stats, nil
}
