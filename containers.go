package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// containerCollector contains the docker client from which to collect various
// container stats.
type containerCollector struct {
	cli *client.Client
}

// newContainerCollector returns a new containerCollector.
func newContainerCollector(cli *client.Client) *containerCollector {
	return &containerCollector{cli: cli}
}

// collect collects container information from a docker container.
func (c *containerCollector) collect() ([]measurement, error) {
	if c.cli == nil {
		return nil, fmt.Errorf("Client not established")
	}

	containers, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Size: true})
	if err != nil {
		return nil, fmt.Errorf("Failed to list containers - %s", err.Error())
	}

	m := measurement{name: "containers", fields: map[string]interface{}{}}
	m.fields["total"] = len(containers)
	if m.fields["size_rw"] == nil {
		m.fields["size_rw"] = int64(0)
	}

	for i := range containers {
		m.fields["size_rw"] = m.fields["size_rw"].(int64) + containers[i].SizeRw
		if m.fields[containers[i].State] == nil {
			m.fields[containers[i].State] = 0
		}
		m.fields[containers[i].State] = m.fields[containers[i].State].(int) + 1

		// populate the used images/volumes (for use by the images/volumes collector).
		// todo: pull out into separate functions (listContainers called by populateImages, populateVolumes)
		usedImages[parseID(containers[i].ImageID)] = containers[i].Image
		populateVolumes(containers[i].Mounts)
	}

	return []measurement{m}, nil
}
