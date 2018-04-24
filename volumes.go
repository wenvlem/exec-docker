package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// volumeCollector contains the docker client from which to collect various
// volume stats.
type volumeCollector struct {
	cli *client.Client
}

// newVolumeCollector returns a new VolumeCollector.
func newVolumeCollector(cli *client.Client) *volumeCollector {
	return &volumeCollector{cli: cli}
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

// collect collects volume information from a docker container.
func (c *volumeCollector) collect() ([]measurement, error) {
	if c.cli == nil {
		return nil, fmt.Errorf("Client not established")
	}

	usage, err := c.cli.DiskUsage(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Failed to list volumes - %s", err.Error())
	}
	volumes := usage.Volumes

	ms := []measurement{}
	m := measurement{name: "volumes", tags: map[string]string{"volume": "total"}, fields: map[string]interface{}{}}
	m.fields["total"] = len(volumes)

	if m.fields["unused"] == nil {
		m.fields["unused"] = int(0)
	}
	if m.fields["size"] == nil {
		m.fields["size"] = int64(0)
	}

	for i := range volumes {
		me := measurement{name: "volumes", tags: map[string]string{"volume": volumes[i].Name}, fields: map[string]interface{}{}}

		if !volumeUsed(volumes[i].Name) {
			m.fields["unused"] = m.fields["unused"].(int) + 1
		}
		m.fields["size"] = m.fields["size"].(int64) + volumes[i].UsageData.Size
		me.fields["size"] = volumes[i].UsageData.Size
		ms = append(ms, me)
	}
	ms = append(ms, m)

	return ms, nil
}
