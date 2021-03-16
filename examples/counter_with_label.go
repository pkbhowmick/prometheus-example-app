package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_call_counts",
		Help: "The number of api calls",
	}, []string{"path", "method"})
)

func main() {
	router := mux.NewRouter()
	reg := prometheus.NewRegistry()
	reg.MustRegister(requests)

	router.HandleFunc("/api/{path}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		path := "/api/" + vars["path"]
		requests.WithLabelValues(path, req.Method).Inc()
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(path))
	})
	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		requests.WithLabelValues("/", req.Method).Inc()
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Hello from homepage"))

	})

	router.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	server := http.Server{
		Addr:    ":8085",
		Handler: router,
	}
	log.Println("Server is running")
	log.Fatal(server.ListenAndServe())
}
