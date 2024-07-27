package agent

import (
	"fmt"
	"github.com/halviet/metrics/internal/storage"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type Agent struct {
	MemStats    runtime.MemStats
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
	baseURL := "http://localhost:8080/update/"
	metrics := make(map[string]string)

	metrics["Alloc"] = strconv.FormatUint(a.MemStats.Alloc, 10)
	metrics["BuckHashSys"] = strconv.FormatUint(a.MemStats.BuckHashSys, 10)
	metrics["Frees"] = strconv.FormatUint(a.MemStats.Frees, 10)
	metrics["GCCPUFraction"] = strconv.FormatFloat(a.MemStats.GCCPUFraction, 'g', -1, 64)
	metrics["GCSys"] = strconv.FormatUint(a.MemStats.GCSys, 10)
	metrics["HeapAlloc"] = strconv.FormatUint(a.MemStats.HeapAlloc, 10)
	metrics["HeapIdle"] = strconv.FormatUint(a.MemStats.HeapIdle, 10)
	metrics["HeapInuse"] = strconv.FormatUint(a.MemStats.HeapInuse, 10)
	metrics["HeapObjects"] = strconv.FormatUint(a.MemStats.HeapObjects, 10)
	metrics["HeapReleased"] = strconv.FormatUint(a.MemStats.HeapReleased, 10)
	metrics["HeapSys"] = strconv.FormatUint(a.MemStats.HeapSys, 10)
	metrics["LastGC"] = strconv.FormatUint(a.MemStats.LastGC, 10)
	metrics["Lookups"] = strconv.FormatUint(a.MemStats.Lookups, 10)
	metrics["MCacheInuse"] = strconv.FormatUint(a.MemStats.MCacheInuse, 10)
	metrics["MCacheSys"] = strconv.FormatUint(a.MemStats.MCacheSys, 10)
	metrics["MSpanInuse"] = strconv.FormatUint(a.MemStats.MSpanInuse, 10)
	metrics["MSpanSys"] = strconv.FormatUint(a.MemStats.MSpanSys, 10)
	metrics["Mallocs"] = strconv.FormatUint(a.MemStats.Mallocs, 10)
	metrics["NextGC"] = strconv.FormatUint(a.MemStats.NextGC, 10)
	metrics["NumForcedGC"] = strconv.FormatUint(uint64(a.MemStats.NumForcedGC), 10)
	metrics["NumGC"] = strconv.FormatUint(uint64(a.MemStats.NumGC), 10)
	metrics["OtherSys"] = strconv.FormatUint(a.MemStats.OtherSys, 10)
	metrics["PauseTotalNs"] = strconv.FormatUint(a.MemStats.PauseTotalNs, 10)
	metrics["StackInuse"] = strconv.FormatUint(a.MemStats.StackInuse, 10)
	metrics["StackSys"] = strconv.FormatUint(a.MemStats.StackSys, 10)
	metrics["Sys"] = strconv.FormatUint(a.MemStats.Sys, 10)
	metrics["TotalAlloc"] = strconv.FormatUint(a.MemStats.TotalAlloc, 10)

	sendMetric := func(metricType, name, v string) error {
		resp, err := http.Post(baseURL+fmt.Sprintf("%q/%q/%q", metricType, name, v), "text/plain", nil)
		if err != nil {
			return err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal("Closing body err:", err)
			}
		}(resp.Body)

		if resp.StatusCode != 200 {
			return fmt.Errorf("request failed, server sent %d status code", resp.StatusCode)
		}

		return nil
	}

	for name, v := range metrics {
		err := sendMetric("gauge", name, v)
		if err != nil {
			return err
		}
	}

	err := sendMetric("counter", "PollCount", strconv.FormatInt(a.PollCount, 10))
	if err != nil {
		return err
	}
	err = sendMetric("gauge", "RandomValue", strconv.FormatFloat(float64(a.RandomValue), 'g', -1, 64))
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
