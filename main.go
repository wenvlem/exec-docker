package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/docker/docker/client"
)

func main() {
	cli, err := newClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dock := dockerEngine{cli: cli}

	dock.addCollector(newContainerCollector(cli))
	dock.addCollector(newImageCollector(cli))
	dock.addCollector(newVolumeCollector(cli))

	var mTex sync.RWMutex
	var measurements = []measurement{}

	for i := range dock.collectors {
		// todo: goroutine
		m, err := dock.collectors[i].collect()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		mTex.Lock()
		measurements = append(measurements, m...)
		mTex.Unlock()
	}

	dock.addFormatter(newInfluxFormatter())

	for i := range dock.formatters {
		err = publish(dock.formatters[i].format(measurements), os.Stdout)
		if err != nil {
			fmt.Printf("Failed to publish measurements - %s\n", err.Error())
			return
		}
	}
}

func publish(d string, w io.WriteCloser) error {
	defer w.Close()
	_, err := fmt.Fprint(w, d)
	return err
}

type (
	collector interface {
		collect() ([]measurement, error)
	}

	formatter interface {
		format([]measurement) string
	}

	measurement struct {
		name   string
		tags   map[string]interface{}
		fields map[string]interface{}
	}

	dockerEngine struct {
		cli        *client.Client
		collectors []collector
		formatters []formatter
	}
)

func newClient() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to establish new client - %s", err.Error())
	}

	cli.NegotiateAPIVersion(context.Background())

	return cli, nil
}

func (d *dockerEngine) addCollector(c collector) {
	d.collectors = append(d.collectors, c)
}

func (d *dockerEngine) addFormatter(f formatter) {
	d.formatters = append(d.formatters, f)
}
