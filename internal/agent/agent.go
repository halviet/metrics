package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/halviet/metrics/internal/storage"
	"github.com/halviet/metrics/internal/storage/models"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Agent struct {
	MemStats    runtime.MemStats
	SrvAddr     string
	PollCount   int64
	RandomValue storage.Gauge
}

func New() Agent {
	return Agent{
		PollCount: 0,
	}
}

func (a *Agent) Update() {
	runtime.ReadMemStats(&a.MemStats)
	a.PollCount += 1
	a.RandomValue = storage.Gauge(rand.Float64())
}

func (a *Agent) Send() error {
	baseURL := "http://" + a.SrvAddr + "/update/"
	metrics := make(map[string]float64)

	metrics["Alloc"] = float64(a.MemStats.Alloc)
	metrics["BuckHashSys"] = float64(a.MemStats.BuckHashSys)
	metrics["Frees"] = float64(a.MemStats.Frees)
	metrics["GCCPUFraction"] = a.MemStats.GCCPUFraction
	metrics["GCSys"] = float64(a.MemStats.GCSys)
	metrics["HeapAlloc"] = float64(a.MemStats.HeapAlloc)
	metrics["HeapIdle"] = float64(a.MemStats.HeapIdle)
	metrics["HeapInuse"] = float64(a.MemStats.HeapInuse)
	metrics["HeapObjects"] = float64(a.MemStats.HeapObjects)
	metrics["HeapReleased"] = float64(a.MemStats.HeapReleased)
	metrics["HeapSys"] = float64(a.MemStats.HeapSys)
	metrics["LastGC"] = float64(a.MemStats.LastGC)
	metrics["Lookups"] = float64(a.MemStats.Lookups)
	metrics["MCacheInuse"] = float64(a.MemStats.MCacheInuse)
	metrics["MCacheSys"] = float64(a.MemStats.MCacheSys)
	metrics["MSpanInuse"] = float64(a.MemStats.MSpanInuse)
	metrics["MSpanSys"] = float64(a.MemStats.MSpanSys)
	metrics["Mallocs"] = float64(a.MemStats.Mallocs)
	metrics["NextGC"] = float64(a.MemStats.NextGC)
	metrics["NumForcedGC"] = float64(a.MemStats.NumForcedGC)
	metrics["NumGC"] = float64(a.MemStats.NumGC)
	metrics["OtherSys"] = float64(a.MemStats.OtherSys)
	metrics["PauseTotalNs"] = float64(a.MemStats.PauseTotalNs)
	metrics["StackInuse"] = float64(a.MemStats.StackInuse)
	metrics["StackSys"] = float64(a.MemStats.StackSys)
	metrics["Sys"] = float64(a.MemStats.Sys)
	metrics["TotalAlloc"] = float64(a.MemStats.TotalAlloc)

	sendMetric := func(metric models.Metrics) error {
		body := bytes.NewBuffer(nil)
		wBody := gzip.NewWriter(body)

		err := json.NewEncoder(wBody).Encode(metric)
		if err != nil {
			return err
		}

		wBody.Close()

		r, err := http.NewRequest(http.MethodPost, baseURL, body)
		if err != nil {
			return err
		}
		r.Header.Set("Content-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			return err
		}
		err = resp.Body.Close()
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("request failed, server sent %d status code", resp.StatusCode)
		}

		return nil
	}

	for name, v := range metrics {
		err := sendMetric(models.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &v,
		})
		if err != nil {
			return err
		}
	}

	err := sendMetric(models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &a.PollCount,
	})
	if err != nil {
		return err
	}

	rv := float64(a.RandomValue)
	err = sendMetric(models.Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: &rv,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *Agent) Start(pollInterval, reportInterval time.Duration) error {
	go func() {
		for {
			a.Update()
			time.Sleep(pollInterval)
		}
	}()

	sender := make(chan error)
	go func() {
		time.Sleep(reportInterval)
		for {
			err := a.Send()
			if err != nil {
				sender <- err
				return
			}
			time.Sleep(reportInterval)
		}
	}()

	s := <-sender
	return s
}
