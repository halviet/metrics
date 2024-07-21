package main

import (
	"github.com/halviet/metrics/internal/storage"
	"log"
	"net/http"
	"strconv"
)

var store = storage.New()

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update/{metricType}/{metricName}/{metricValue}", MetricHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func MetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var (
		metricType  = r.PathValue("metricType")
		metricName  = r.PathValue("metricName")
		metricValue = r.PathValue("metricValue")
	)

	if metricName == "" {
		http.NotFound(w, r)
	}

	switch metricType {
	case "gauge":
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Not correct Gauge value", http.StatusBadRequest)
		}

		gauge := storage.Gauge(val)
		store.UpdateGauge(metricName, gauge)
	case "counter":
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Not correct Counter value", http.StatusBadRequest)
		}

		counter := storage.Counter(val)
		store.UpdateCounter(metricName, counter)
	default:
		http.Error(w, "Provided type doesn't exist", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}
