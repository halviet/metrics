package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/halviet/metrics/internal/handlers"
	"github.com/halviet/metrics/internal/logger"
	mw "github.com/halviet/metrics/internal/middleware"
	"github.com/halviet/metrics/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	l, err := logger.New(logger.Opts{Lvl: cfg.LogLevel})
	if err != nil {
		panic(err)
	}

	store := storage.New()
	r := chi.NewRouter()

	r.Use(mw.Log(l))

	r.Post("/update/", handlers.JSONMetricHandler(store))
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.MetricHandler(store))

	r.Get("/value/", handlers.JSONGetMetricHandle(store))
	r.Get("/value/{metricType}/{metricName}", handlers.GetMetricHandle(store))
	r.Get("/", handlers.GetAllMetricsPageHandler(store))

	if err = http.ListenAndServe(cfg.SrvAddr, r); err != nil {
		l.Fatal("internal error", zap.Error(err))
	}
}
