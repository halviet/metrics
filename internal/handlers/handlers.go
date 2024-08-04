package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/halviet/metrics/internal/storage"
	"html/template"
	"net/http"
	"strconv"
)

func MetricHandler(store *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			metricType  = chi.URLParam(r, "metricType")
			metricName  = chi.URLParam(r, "metricName")
			metricValue = chi.URLParam(r, "metricValue")
		)

		switch metricType {
		case "gauge":
			val, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(w, "Not correct Gauge value", http.StatusBadRequest)
				return
			}

			gauge := storage.Gauge(val)
			store.UpdateGauge(metricName, gauge)
		case "counter":
			val, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(w, "Not correct Counter value", http.StatusBadRequest)
				return
			}

			counter := storage.Counter(val)
			store.UpdateCounter(metricName, counter)
		default:
			http.Error(w, "Provided type doesn't exist", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetMetricHandle(store *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			metricType = chi.URLParam(r, "metricType")
			metricName = chi.URLParam(r, "metricName")
		)

		switch metricType {
		case "gauge":
			gauge, err := store.GetGauge(metricName)
			if err != nil {
				if errors.Is(err, storage.ErrMetricNotFound) {
					http.NotFound(w, r)
					return
				}

				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "text/plain charset=UTF-8")
			w.Write([]byte(strconv.FormatFloat(float64(gauge), 'g', -1, 64)))
			w.WriteHeader(http.StatusOK)
		case "counter":
			counter, err := store.GetCounter(metricName)
			if err != nil {
				if errors.Is(err, storage.ErrMetricNotFound) {
					http.NotFound(w, r)
					return
				}

				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "text/plain charset=UTF-8")
			w.Write([]byte(strconv.FormatInt(int64(counter), 10)))
			w.WriteHeader(http.StatusOK)
		}
	}
}

func GetAllMetricsPageHandler(store *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := store.GetAllMetrics()

		var gauges string
		for _, key := range sortedMetricKeys(data.Gauge) {
			gauges += fmt.Sprintf("<li><b>%v</b>: %v<li>", key, strconv.FormatFloat(float64(data.Gauge[key]), 'f', -1, 64))
		}

		var counters string
		for _, key := range sortedMetricKeys(data.Counter) {
			counters += fmt.Sprintf("<li><b>%v</b>: %v</li>", key, strconv.FormatInt(int64(data.Counter[key]), 10))
		}

		values := struct {
			Gauge   template.HTML
			Counter template.HTML
		}{
			Gauge:   template.HTML(gauges),
			Counter: template.HTML(counters),
		}

		tmpl, err := template.ParseFiles("./internal/templates/metrics.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, values)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
