package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/halviet/metrics/internal/handlers"
	"github.com/halviet/metrics/internal/storage"
	"log"
	"net/http"
)

var store = storage.New()

func main() {
	r := chi.NewRouter()

	r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.MetricHandler(store))
	r.Get("/value/{metricType}/{metricName}", handlers.GetMetricHandle(store))
	r.Get("/", handlers.GetAllMetricsPageHandler(store))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
