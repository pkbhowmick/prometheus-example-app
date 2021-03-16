package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

var (
	foos = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "foos_total",
		Help: "The number of foo calls",
	})
	reg = prometheus.NewRegistry()
)

func Foo() {
	foos.Inc()
}

func getMetricValue(col prometheus.Collector) float64 {
	c := make(chan prometheus.Metric, 1) // 1 for metric with no vector
	col.Collect(c)                       // collect current metric value into the channel
	m := dto.Metric{}
	_ = (<-c).Write(&m) // read metric value from the channel
	return *m.Counter.Value
}

func TestFoo(t *testing.T) {
	reg.MustRegister(foos)
	before := getMetricValue(foos)
	Foo()
	after := getMetricValue(foos)
	assert.Equal(t, 1, int(after-before))
}
