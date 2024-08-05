package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/halviet/metrics/internal/handlers"
	"github.com/halviet/metrics/internal/storage"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("a", "localhost:8080", "HTTP-server address to run on")

	flag.Parse()

	store := storage.New()
	r := chi.NewRouter()

	r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.MetricHandler(store))
	r.Get("/value/{metricType}/{metricName}", handlers.GetMetricHandle(store))
	r.Get("/", handlers.GetAllMetricsPageHandler(store))

	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatal(err)
	}
}
