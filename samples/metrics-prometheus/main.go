package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reqs = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "app_requests_total", Help: "total requests"}, []string{"path", "method"})
)

func main() {
	prometheus.MustRegister(reqs)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqs.WithLabelValues(r.URL.Path, r.Method).Inc()
		w.Write([]byte("ok"))
	})
	http.Handle("/metrics", promhttp.Handler())
	log.Println("listening on :2112")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
