package main

import (
	"fmt"
	"strings"
)

// influxFormatter allows measurements to be formatted in influx line protocol.
type influxFormatter struct{}

// newInfluxFormatter returns a new influxFormatter.
func newInfluxFormatter() influxFormatter {
	return influxFormatter{}
}

// format returns measurements in influx line protocol format.
func (i influxFormatter) format(m []measurement) string {
	s := ""

	for i := range m {
		tags := []string{}
		for k := range m[i].tags {
			tags = append(tags, fmt.Sprintf("%s=%s", k, m[i].tags[k]))
		}
		fields := []string{}
		for k := range m[i].fields {
			fields = append(fields, fmt.Sprintf("%s=%v", k, m[i].fields[k]))
		}

		s += fmt.Sprint(m[i].name)
		if len(tags) > 0 {
			s += fmt.Sprintf(",%s", strings.Join(tags, ","))
		}
		if len(fields) > 0 {
			s += fmt.Sprintf(" %s", strings.Join(fields, ","))
		}
		s += fmt.Sprint("\n")
	}

	return s
}
