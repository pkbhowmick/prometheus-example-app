package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	latency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "api_latency_second",
		Help: "Time for the request of api",
	})
	histogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "api_latency_second_histogram",
		Help:    "Time for the request of api",
		Buckets: []float64{0.00001, 0.0001, 0.001, 0.01, 0.1, 1.0, 2.0, 5.0, 10.0},
	})
)

func main() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(latency)
	reg.MustRegister(histogram)

	router := http.NewServeMux()

	router.HandleFunc("/hello", func(res http.ResponseWriter, req *http.Request) {
		start := time.Now().Unix()
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Hello"))
		time.Sleep(time.Second)
		latencyTime := float64(time.Now().Unix() - start)
		latency.Observe(latencyTime)
		histogram.Observe(latencyTime)
	})

	router.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	server := http.Server{
		Addr:    ":8085",
		Handler: router,
	}
	fmt.Println("Server is running...")
	log.Fatalln(server.ListenAndServe())
}
