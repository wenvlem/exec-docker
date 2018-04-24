package main

import (
	"bytes"
	"testing"
)

var container_id string

func TestMain(m *testing.M) {
	m.Run()
}

func TestInfluxFormatter(t *testing.T) {
	iFF := newInfluxFormatter()

	ms := []measurement{
		{
			name:   "test",
			fields: map[string]interface{}{"metricA": "123"},
		},
	}
	out := iFF.format(ms)

	if out != "test metricA=123\n" {
		t.Errorf("Format incorrect! Got: '%s'", out)
		t.FailNow()
	}

	ms = []measurement{
		{
			name: "test",
			tags: map[string]interface{}{"metricA": "123"},
		},
	}
	out = iFF.format(ms)

	if out != "test,metricA=123\n" {
		t.Errorf("Format incorrect! Got: '%s'", out)
		t.FailNow()
	}

	ms = []measurement{
		{
			name:   "test",
			tags:   map[string]interface{}{"metricA": "123"},
			fields: map[string]interface{}{"fieldA": "abc"},
		},
	}
	out = iFF.format(ms)

	if out != "test,metricA=123 fieldA=abc\n" {
		t.Errorf("Format incorrect! Got: '%s'", out)
		t.FailNow()
	}
}

func TestVolumes(t *testing.T) {
	cli, err := newClient()
	if err != nil {
		t.Errorf("Failed to initialize client - '%s'", err.Error())
		t.FailNow()
	}

	vc := newVolumeCollector(cli)
	_, err = vc.collect()
	if err != nil {
		t.Errorf("Failed to collect volumes - '%s'", err.Error())
		t.FailNow()
	}

	vc = newVolumeCollector(nil)
	_, err = vc.collect()
	if err == nil {
		t.Errorf("Failed to fail collecting volumes")
		t.FailNow()
	}
}

func TestContainers(t *testing.T) {
	cli, err := newClient()
	if err != nil {
		t.Errorf("Failed to initialize client - '%s'", err.Error())
		t.FailNow()
	}

	cc := newContainerCollector(cli)
	_, err = cc.collect()
	if err != nil {
		t.Errorf("Failed to collect containers - '%s'", err.Error())
		t.FailNow()
	}

	cc = newContainerCollector(nil)
	_, err = cc.collect()
	if err == nil {
		t.Errorf("Failed to fail collecting containers")
		t.FailNow()
	}
}

func TestImages(t *testing.T) {
	cli, err := newClient()
	if err != nil {
		t.Errorf("Failed to initialize client - '%s'", err.Error())
		t.FailNow()
	}

	ic := newImageCollector(cli)
	_, err = ic.collect()
	if err != nil {
		t.Errorf("Failed to collect images - '%s'", err.Error())
		t.FailNow()
	}

	ic = newImageCollector(nil)
	_, err = ic.collect()
	if err == nil {
		t.Errorf("Failed to fail collecting images")
		t.FailNow()
	}
}

func TestPublish(t *testing.T) {
	var b bytes.Buffer
	publish("test string\n", &b)
	if b.String() != "test string\n" {
		t.Errorf("Failed to publish! Got: '%s'", b.String())
		t.FailNow()
	}
}

func TestEngine(t *testing.T) {
	iFF := newInfluxFormatter()
	cc := newContainerCollector(nil)

	dock := dockerEngine{}
	dock.addCollector(cc)
	dock.addFormatter(iFF)
}

func TestMainFunc(t *testing.T) {
	main()
}
