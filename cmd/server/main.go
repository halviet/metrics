package main

import (
	"github.com/halviet/metrics/internal/handlers"
	"github.com/halviet/metrics/internal/storage"
	"log"
	"net/http"
)

var store = storage.New()

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update/{metricType}/{metricName}/{metricValue}", handlers.MetricHandler(store))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
