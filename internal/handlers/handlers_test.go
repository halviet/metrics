package handlers

import (
	"fmt"
	"github.com/halviet/metrics/internal/storage"
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
			target := fmt.Sprintf("/update/%q/%q/%q", test.metricType, test.metricName, test.metricValue)
			req := httptest.NewRequest(http.MethodPost, target, nil)
			req.SetPathValue("metricType", test.metricType)
			req.SetPathValue("metricName", test.metricName)
			req.SetPathValue("metricValue", test.metricValue)

			w := httptest.NewRecorder()

			store := storage.New()
			MetricHandler(store)(w, req)

			res := w.Result()

			if res.StatusCode != test.want.statusCode {
				t.Errorf("invalid status code, got: %d, want %d", res.StatusCode, test.want.statusCode)
			}

			if test.want.statusCode == http.StatusOK {
				switch test.metricType {
				case "gauge":
					val, err := strconv.ParseFloat(test.metricValue, 64)
					if err != nil {
						t.Errorf("err on parsing metricValue: %q; err: %v", test.metricValue, err)
					}

					if store.GetGauge(test.metricName) != storage.Gauge(val) {
						t.Errorf("value has been wrote wrong, got: %f; want: %f", store.GetGauge(test.metricName), storage.Gauge(val))
					}

				}
			}
		})
	}
}
