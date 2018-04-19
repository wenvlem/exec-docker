package main

// tags are indexed, fields are not
// <measurement>[,<tag_key>=<tag_value>[,<tag_key>=<tag_value>]] <field_key>=<field_value>[,<field_key>=<field_value>] [<timestamp>]

import (
	"context"
	// "encoding/json"
	"fmt"
	"io"
	// "strings"
	"os"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := newClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dock := dockerEngine{cli: cli}

	c := newContainerCollector(cli)
	dock.addCollector(&c)

	var mTex sync.RWMutex
	var measurements = map[string]measurement{}

	// err = dock.getDanglingImg()
	for i := range dock.collectors {
		err := dock.collectors[i].collect()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		m, err := dock.collectors[i].filter()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for i := range m {
			mTex.Lock()
			measurements[m[i].name] = m[i]
			mTex.Unlock()
		}
	}

	f := newInfluxFormatter()
	dock.addFormatter(f, os.Stdout)

	for i := range dock.formatters {
		// todo: instead of os.Stdout, add writer/closer to
		err = publish(dock.formatters[i].fmtr.format(measurements), dock.formatters[i].writer)
		if err != nil {
			fmt.Printf("Failed to publish measurements - %s\n", err.Error())
			return
		}
	}
	// todo: if WriteCloser, close
}

func publish(d string, w io.WriteCloser) error {
	defer w.Close()
	_, err := fmt.Fprint(w, d)
	return err
}

type (
	collector interface {
		collect() error
		filter() ([]measurement, error)
	}

	formatter interface {
		// publish(map[string]measurement, io.Writer) error
		// format(map[string]measurement) ([]byte, error)
		format(map[string]measurement) string
	}

	fmter struct {
		fmtr   formatter
		writer io.WriteCloser
	}

	dockerEngine struct {
		cli        *client.Client
		collectors []collector
		// formatters []formatter
		formatters []fmter
		// containers []types.Container
		// volumes []types.Volume
		// networks []types.Network
		stats map[string]interface{}
	}

	// "containers",
	// "containers_size_rw",
	// "containers_total",
	// "containers_created",
	// "containers_restarting",
	// "containers_running",
	// "containers_removing",
	// "containers_paused",
	// "containers_exited",
	// "containers_dead",

	// "images_total"
	// "images_unused"
	// "images_dangling"

	// "volumes_total"
	// "volumes_unused"

	// "networks_total"
	// "networks_unused"
	// container struct {
	// 	ID string
	// 	Name string
	// 	State string
	// 	SizeRoot int64
	// 	SizeRw int64
	// 	// holds cont
	// 	stats map[string]interface{}
	// }
)

func newClient() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to establish new client - %s", err.Error())
	}

	cli.NegotiateAPIVersion(context.Background())

	return cli, nil
}

func (d *dockerEngine) addCollector(c collector) error {
	d.collectors = append(d.collectors, c)
	return nil
}

func (d *dockerEngine) addFormatter(f formatter, w io.WriteCloser) error {
	d.formatters = append(d.formatters, fmter{fmtr: f, writer: w})
	return nil
}

// func (d *dockerEngine) listContainers() error {
// 	containers, err := d.cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Size: true})
// 	if err != nil {
// 		return fmt.Errorf("Failed to list containers - %s", err.Error())
// 	}
// 	d.containers = containers
// 	return nil
// }

func (d *dockerEngine) listImages(dangling bool) error {
	fil := filters.NewArgs()
	if dangling {
		fil.Add("dangling", "true")
	} else {
		fil.Add("dangling", "false")
	}

	imgs, err := d.cli.ImageList(context.Background(), types.ImageListOptions{All: true, Filters: fil})
	if err != nil {
		return fmt.Errorf("Failed to list images - %s", err.Error())
	}
	fmt.Printf("IMAGES: %+v\n", imgs)
	for _, i := range imgs {
		fmt.Println(i.Containers)
	}
	return nil
}

// func (d *dockerEngine) collectContainers() {
// 	for _, c := range d.containers {
// 		var v *types.StatsJSON

// 		// func (cli *Client) ContainerStats(ctx context.Context, containerID string, stream bool) (types.ContainerStats, error)
// 		r, err := d.cli.ContainerStats(context.Background(), c.ID, false)
// 		if err != nil {
// 			fmt.Printf("Failed to fetch stats for '%s' - %s\n", c.Names, err.Error())
// 			return
// 		}
// 		defer r.Body.Close()

// 		dec := json.NewDecoder(r.Body)
// 		if err = dec.Decode(&v); err != nil {
// 			if err == io.EOF {
// 				return
// 			}
// 			fmt.Printf("Error decoding: %s\n", err.Error())
// 			return
// 		}

// 		fmt.Printf("Stats - %+v\n", v)
// 	}
// 	return

// 	for _, container := range containers {
// 		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
// 	}
// }
