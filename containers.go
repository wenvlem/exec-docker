package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type container struct {
	cli    *client.Client
	tags   map[string]interface{}
	fields map[string]interface{}
}

func newContainerCollector(cli *client.Client) *container {
	return &container{cli: cli, fields: make(map[string]interface{})}
}

func (c *container) collect() error {
	if c.cli == nil {
		return fmt.Errorf("Client not established")
	}
	if c.fields == nil {
		fmt.Println("making thing")
		c.fields = make(map[string]interface{})
	}

	containers, err := listContainers(c.cli)
	if err != nil {
		return fmt.Errorf("Failed to list containers - %s", err.Error())
	}

	c.fields["total"] = len(containers)
	if c.fields["size_rw"] == nil {
		c.fields["size_rw"] = int64(0)
	}

	for i := range containers {
		c.fields["size_rw"] = c.fields["size_rw"].(int64) + containers[i].SizeRw
		if c.fields[containers[i].State] == nil {
			c.fields[containers[i].State] = 0
		}
		c.fields[containers[i].State] = c.fields[containers[i].State].(int) + 1
		usedImages[parseID(containers[i].ImageID)] = containers[i].Image
	}
	return nil
}

func listContainers(cli *client.Client) ([]types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Size: true})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (c *container) filter() ([]measurement, error) {
	if c.fields == nil {
		fmt.Println("filter making")
		c.fields = make(map[string]interface{})
	}

	return []measurement{measurement{name: "containers", fields: c.fields}}, nil
}
