// exec-docker is a simple binary to be used with telegraf's `exec` plugin.
// It outputs in influx line protocol. It collects and reports several
// image, volume, and container metrics.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/docker/docker/client"
)

// main calls start() and runs the application. It isn't necessary, but
// serves a similar purpose to the start function, making it clear that
// it can stand alone in it's own package.
func main() {
	start()
}

// start runs the collector. It is a separate method (from main) in case
// there was ever cause to export anything and call from a separate package.
func start() {
	cli, err := newClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	dock := dockerEngine{}

	// These collectors could be added dependent on cli flags.
	dock.addCollector(newContainerCollector(cli))
	dock.addCollector(newImageCollector(cli))
	dock.addCollector(newVolumeCollector(cli))

	var mTex sync.RWMutex
	var measurements = []measurement{}

	var wg sync.WaitGroup

	// collect from collectors.
	for i := range dock.collectors {
		wg.Add(1)

		go func(d collector) {
			defer wg.Done()
			m, err := d.collect()
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
			mTex.Lock()
			measurements = append(measurements, m...)
			mTex.Unlock()
		}(dock.collectors[i])
	}

	// wait for all collectors to finish.
	wg.Wait()

	// set the formatter, currently there is only support for an influx formatter
	// but if other's are needed, this can be in a switch statement to set the
	// formatter to another output telegraf supports.
	dock.setFormatter(newInfluxFormatter())

	err = publish(dock.formatter.format(measurements), os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to publish measurements - %s\n", err.Error())
		return
	}
}

// publish writes the formatted metrics (string) to a writer.
func publish(d string, w io.Writer) error {
	_, err := fmt.Fprint(w, d)
	return err
}

type (
	// collector defines what a collector must do.
	collector interface {
		collect() ([]measurement, error)
	}

	// formatter defines what a formatter must do.
	formatter interface {
		format([]measurement) string
	}

	// measurement defines a metric.
	measurement struct {
		name   string
		tags   map[string]string
		fields map[string]interface{}
	}

	// dockerEngine defines a docker engine to gather from.
	dockerEngine struct {
		collectors []collector
		formatter  formatter
	}
)

// newClient creates a new docker client and negotiate's the API version.
func newClient() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to establish new client - %s", err.Error())
	}

	cli.NegotiateAPIVersion(context.Background())

	return cli, nil
}

// addCollector adds a collector to the dockerEngine.
func (d *dockerEngine) addCollector(c collector) {
	d.collectors = append(d.collectors, c)
}

// setFormatter sets the dockerEngine's formatter.
func (d *dockerEngine) setFormatter(f formatter) {
	d.formatter = f
}
