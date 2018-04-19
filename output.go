package main

import (
	"sync"
)

type measurement struct {
	name   string
	tags   map[string]interface{}
	fields map[string]interface{}
}

var mTex sync.RWMutex
var measurements = map[string]measurement{}

// func (m *measurements) Add(name string, metric measurement) {
// 	mTex.Lock()
// 	defer mTex.Unlock()
// 	// measurements = append(measurements, metric)
// 	measurements[name] = metric
// }
