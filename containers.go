package main

import (
	"context"
	"fmt"
	// "io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type container struct {
	cli    *client.Client
	tags   map[string]interface{}
	fields map[string]interface{}
}

func newContainerCollector(cli *client.Client) container {
	return container{cli: cli, fields: make(map[string]interface{})}
}

var usedImages = map[string]string{} // "nanobox/logvac:latest": "abc123..."
// func (img image) isUsed(i) {
// 	// if _, ok := usedImages[i]; !ok {
// 	// 	return false
// 	// }
// 	_, ok := usedImages[i]
// 	return ok
// }

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

func parseID(i string) string {
	j := strings.Split(i, ":")
	if len(j) > 0 {
		return j[1]
	}
	return i
}

func listContainers(cli *client.Client) ([]types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Size: true})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// func collectContainer(c types.Container) (map[string]int, error) {
// 		// "containers",
// 		// "containers_size_rw",
// 		// "containers_total",
// 		// "containers_created",
// 		// "containers_restarting",
// 		// "containers_running",
// 		// "containers_removing",
// 		// "containers_paused",
// 		// "containers_exited",
// 		// "containers_dead",
// "size_rw"
// "total"
// "created"
// "restarting"
// "running"
// "removing"
// "paused"
// "exited"
// "dead"
// 	c.State

// 	// for _, c := range d.containers {
// 	var v *types.StatsJSON

// 	// func (cli *Client) ContainerStats(ctx context.Context, containerID string, stream bool) (types.ContainerStats, error)
// 	r, err := d.cli.ContainerStats(context.Background(), c.ID, false)
// 	if err != nil {
// 		fmt.Printf("Failed to fetch stats for '%s' - %s\n", c.Names, err.Error())
// 		return
// 	}
// 	defer r.Body.Close()

// 	dec := json.NewDecoder(r.Body)
// 	if err = dec.Decode(&v); err != nil {
// 		if err == io.EOF {
// 			return
// 		}
// 		fmt.Printf("Error decoding: %s\n", err.Error())
// 		return
// 	}

// 	fmt.Printf("Stats - %+v\n", v)
// 	// }
// 	return

// 	for _, container := range containers {
// 		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
// 	}
// }

func (c *container) filter() ([]measurement, error) {
	if c.cli == nil {
		return nil, fmt.Errorf("Client not established.")
	}
	if c.fields == nil {
		fmt.Println("filter making")
		c.fields = make(map[string]interface{})
	}

	// todo: create a "measurement"
	return []measurement{measurement{name: "containers", fields: c.fields}}, nil
}

// func (c *container) publish(w io.Writer) error {
// 	if c.cli == nil {
// 		return fmt.Errorf("Client not established.")
// 	}

// 	// todo: write a measurement
// 	fmt.Printf("%+v\n",c.fields)
// 	fmt.Printf("%+v\n",usedImages)
// 	return nil
// }
