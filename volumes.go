package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type volumeCollector struct {
	cli   *client.Client
	stats []measurement // tag=volume:total size:+=usagetotal.size; volume:a size:usagetotal.size
}

func newVolumeCollector(cli *client.Client) *volumeCollector {
	// return &volumeCollector{cli: cli, tags: make(map[string]interface{}), fields: make(map[string]interface{})}
	return &volumeCollector{cli: cli, stats: []measurement{}}
}

// usedVolumes is a map of volume names.
var usedVolumes = map[string]struct{}{}

// volumeUsed returns true if a volume is in use by a container.
func volumeUsed(s string) bool {
	_, ok := usedVolumes[s]
	return ok
}

// populateVolumes populates the usedVolumes map.
func populateVolumes(m []types.MountPoint) {
	for i := range m {
		if m[i].Type != "volume" {
			continue
		}
		usedVolumes[m[i].Name] = struct{}{}
	}
}

// func (c *volumeCollector) collect() ([]measurement, error) {
func (c *volumeCollector) collect() error {
	if c.cli == nil {
		return fmt.Errorf("Client not established")
	}
	ms := []measurement{}

	usage, err := c.cli.DiskUsage(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to list volumes - %s", err.Error())
	}
	volumes := usage.Volumes

	m := measurement{name: "volumes", tags: map[string]interface{}{"volume": "total"}, fields: map[string]interface{}{}}
	m.fields["total"] = len(volumes)

	if m.fields["unused"] == nil {
		m.fields["unused"] = int(0)
	}
	if m.fields["size"] == nil {
		m.fields["size"] = int64(0)
	}

	for i := range volumes {
		me := measurement{name: "volumes", tags: map[string]interface{}{"volume": volumes[i].Name}, fields: map[string]interface{}{}}

		if !volumeUsed(volumes[i].Name) {
			m.fields["unused"] = m.fields["unused"].(int) + 1
		}
		m.fields["size"] = m.fields["size"].(int64) + volumes[i].UsageData.Size
		me.fields["size"] = volumes[i].UsageData.Size
		ms = append(ms, me)
	}
	ms = append(ms, m)
	c.stats = ms
	return nil
	// 	return ms, nil
}

func (c *volumeCollector) filter() ([]measurement, error) {
	if c.stats == nil {
		c.stats = []measurement{}
	}

	return c.stats, nil
}
