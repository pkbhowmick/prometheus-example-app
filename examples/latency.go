package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	latency = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "api_latency_second",
		Help: "Time for the request of api",
	})
)

func main() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(latency)

	router := http.NewServeMux()

	router.HandleFunc("/hello", func(res http.ResponseWriter, req *http.Request) {
		start := time.Now().Unix()
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Hello"))
		latency.Observe(float64(time.Now().Unix() - start))
	})

	router.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	server := http.Server{
		Addr:    ":8085",
		Handler: router,
	}
	fmt.Println("Server is running...")
	log.Fatalln(server.ListenAndServe())
}
