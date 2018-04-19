package main

import (
	"fmt"
	"strings"
)

type influxFormatter struct{}

func newInfluxFormatter() influxFormatter {
	return influxFormatter{}
}

func (i influxFormatter) format(m map[string]measurement) string {
	s := ""

	for i := range m {
		tags := []string{}
		for k := range m[i].tags {
			tags = append(tags, fmt.Sprintf("%s=%v", k, m[i].tags[k]))
		}
		fields := []string{}
		for k := range m[i].fields {
			fields = append(fields, fmt.Sprintf("%s=%v", k, m[i].fields[k]))
		}

		s += fmt.Sprint(m[i].name)
		if len(tags) > 0 {
			s += fmt.Sprint(strings.Join(tags, ","))
		}
		if len(fields) > 0 {
			s += fmt.Sprintf(" %s", strings.Join(fields, ","))
		}
		s += fmt.Sprint("\n")
	}

	return s
}
