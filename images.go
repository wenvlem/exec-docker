package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// imageCollector contains the docker client from which to collect various
// image stats.
type imageCollector struct {
	cli *client.Client
}

// newImageCollector returns a new ImageCollector.
func newImageCollector(cli *client.Client) *imageCollector {
	return &imageCollector{cli: cli}
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

// collect collects image information from a docker container.
func (c *imageCollector) collect() ([]measurement, error) {
	if c.cli == nil {
		return nil, fmt.Errorf("Client not established")
	}

	images, err := listImages(c.cli, false)
	if err != nil {
		return nil, fmt.Errorf("Failed to list images - %s", err.Error())
	}
	imagesD, err := listImages(c.cli, true)
	if err != nil {
		return nil, fmt.Errorf("Failed to list dangling images - %s", err.Error())
	}

	m := measurement{name: "images", fields: map[string]interface{}{}}
	m.fields["total"] = len(images) + len(imagesD)
	m.fields["dangling"] = len(imagesD)

	if m.fields["unused"] == nil {
		m.fields["unused"] = int(0)
	}
	if m.fields["size"] == nil {
		m.fields["size"] = int64(0)
	}

	for i := range images {
		if !imgUsed(parseID(images[i].ID)) {
			m.fields["unused"] = m.fields["unused"].(int) + 1
		}
		m.fields["size"] = m.fields["size"].(int64) + images[i].Size
	}

	return []measurement{m}, nil
}

// listImages lists docker images.
func listImages(cli *client.Client, dangling bool) ([]types.ImageSummary, error) {
	fil := filters.NewArgs()
	fil.Add("dangling", fmt.Sprintf("%t", dangling))

	imgs, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true, Filters: fil})
	if err != nil {
		return nil, fmt.Errorf("Failed to list images - %s", err.Error())
	}

	return imgs, nil
}
