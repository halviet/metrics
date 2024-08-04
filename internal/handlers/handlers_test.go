package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/halviet/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		want        want
	}{
		{
			"correct gauge request",
			"gauge",
			"gaugeMetric",
			"100.1",
			want{http.StatusOK},
		},
		{
			"correct counter request",
			"counter",
			"counterMetric",
			"5",
			want{http.StatusOK},
		},
		{
			"with wrong type",
			"anotherType",
			"anotherMetric",
			"0",
			want{http.StatusBadRequest},
		},
		{
			"string as gauge value",
			"gauge",
			"gaugeMetric",
			"Hello",
			want{http.StatusBadRequest},
		},
		{
			"string as counter value",
			"counter",
			"counterValue",
			"Hello",
			want{http.StatusBadRequest},
		},
		{
			"int as gauge value",
			"gauge",
			"gaugeValue",
			"10",
			want{http.StatusOK},
		},
		{
			"float as counter value",
			"counter",
			"counterValue",
			"15.15",
			want{http.StatusBadRequest},
		},
		{
			"counter value overflow",
			"counter",
			"counterValue",
			"9223372036854775808",
			want{http.StatusBadRequest},
		},
		{
			"gauge value overflow",
			"gauge",
			"gaugeValue",
			"1.7E+309",
			want{http.StatusBadRequest},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := storage.New()

			r := chi.NewRouter()
			r.Post("/update/{metricType}/{metricName}/{metricValue}", MetricHandler(store))

			ts := httptest.NewServer(r)

			target := ts.URL + "/update/" + test.metricType + "/" + test.metricName + "/" + test.metricValue
			req, err := http.NewRequest(http.MethodPost, target, nil)
			assert.NoError(t, err)

			resp, err := ts.Client().Do(req)
			assert.NoError(t, err)

			err = resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, test.want.statusCode, resp.StatusCode)

			if test.want.statusCode == http.StatusOK {
				switch test.metricType {
				case "gauge":
					val, err := strconv.ParseFloat(test.metricValue, 64)
					assert.NoError(t, err)

					g, err := store.GetGauge(test.metricName)
					assert.NoError(t, err)
					assert.Equal(t, storage.Gauge(val), g)
				case "counter":
					val, err := strconv.ParseInt(test.metricValue, 10, 64)
					assert.NoError(t, err)

					c, err := store.GetCounter(test.metricName)
					assert.NoError(t, err)
					assert.Equal(t, storage.Counter(val), c)
				}
			}
		})
	}
}

func TestGetMetricHandle(t *testing.T) {
	type want struct {
		statusCode int
		value      string
	}

	tests := []struct {
		name       string
		metricType string
		metricName string
		want       want
	}{
		{
			"gauge value",
			"gauge",
			"gaugeValue",
			want{
				statusCode: http.StatusOK,
				value:      "100.01",
			},
		},
		{
			"counter value",
			"counter",
			"counterValue",
			want{
				statusCode: http.StatusOK,
				value:      "5",
			},
		},
		{
			"not existing gauge value",
			"gauge",
			"notExist",
			want{statusCode: http.StatusNotFound},
		},
		{
			"not existing counter value",
			"counter",
			"notExist",
			want{statusCode: http.StatusNotFound},
		},
	}

	store := storage.New()
	store.UpdateGauge("gaugeValue", storage.Gauge(100.01))
	store.UpdateCounter("counterValue", storage.Counter(5))

	r := chi.NewRouter()
	r.Get("/value/{metricType}/{metricName}", GetMetricHandle(store))

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := ts.URL + "/value/" + test.metricType + "/" + test.metricName
			req, err := http.NewRequest(http.MethodGet, target, nil)
			assert.NoError(t, err)

			resp, err := ts.Client().Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, test.want.statusCode, resp.StatusCode)

			if test.want.statusCode == http.StatusOK {
				respBody, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Equal(t, test.want.value, string(respBody))
			}
		})
	}
}
