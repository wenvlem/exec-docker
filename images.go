package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type imageCollector struct {
	cli    *client.Client
	tags   map[string]interface{}
	fields map[string]interface{}
}

func newImageCollector(cli *client.Client) *imageCollector {
	return &imageCollector{cli: cli, tags: make(map[string]interface{}), fields: make(map[string]interface{})}
}

// usedImages is a map of image ids and their names.
var usedImages = map[string]string{}

// imgUsed returns true if an image is in use by a container.
func imgUsed(s string) bool {
	_, ok := usedImages[s]
	return ok
}

// parseID returns just the hash/id of an image.
func parseID(i string) string {
	j := strings.Split(i, ":")
	if len(j) > 0 {
		return j[1]
	}
	return i
}

func (c *imageCollector) collect() error {
	if c.cli == nil {
		return fmt.Errorf("Client not established")
	}
	if c.fields == nil {
		c.fields = make(map[string]interface{})
	}

	images, err := listImages(c.cli, false)
	if err != nil {
		return fmt.Errorf("Failed to list images - %s", err.Error())
	}
	imagesD, err := listImages(c.cli, true)
	if err != nil {
		return fmt.Errorf("Failed to list dangling images - %s", err.Error())
	}

	c.fields["total"] = len(images) + len(imagesD)
	c.fields["dangling"] = len(imagesD)

	if c.fields["unused"] == nil {
		c.fields["unused"] = int(0)
	}
	if c.fields["size"] == nil {
		c.fields["size"] = int64(0)
	}

	for i := range images {
		if !imgUsed(parseID(images[i].ID)) {
			c.fields["unused"] = c.fields["unused"].(int) + 1
		}
		c.fields["size"] = c.fields["size"].(int64) + images[i].Size
	}
	return nil
}

func listImages(cli *client.Client, dangling bool) ([]types.ImageSummary, error) {
	fil := filters.NewArgs()
	fil.Add("dangling", fmt.Sprintf("%t", dangling))

	imgs, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true, Filters: fil})
	if err != nil {
		return nil, fmt.Errorf("Failed to list images - %s", err.Error())
	}

	return imgs, nil
}

func (c *imageCollector) filter() ([]measurement, error) {
	if c.fields == nil {
		c.fields = make(map[string]interface{})
	}

	return []measurement{measurement{name: "images", fields: c.fields}}, nil
}
