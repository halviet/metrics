package handlers

import (
	"github.com/halviet/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortedMetricKeys(t *testing.T) {
	type want struct{ gauge, counter []string }

	tests := []struct {
		name   string
		metric storage.ResultMetric
		want
	}{
		{
			name: "unsorted metrics",
			metric: storage.ResultMetric{
				Gauge:   map[string]storage.Gauge{"b": 0, "c": 0, "a": 0},
				Counter: map[string]storage.Counter{"c": 0, "a": 0, "b": 0},
			},
			want: want{
				gauge:   []string{"a", "b", "c"},
				counter: []string{"a", "b", "c"},
			},
		},
		{
			name: "empty metrics map",
			metric: storage.ResultMetric{
				Gauge:   map[string]storage.Gauge{},
				Counter: map[string]storage.Counter{},
			},
			want: want{
				gauge:   []string{},
				counter: []string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gauge := sortedMetricKeys(test.metric.Gauge)
			counter := sortedMetricKeys(test.metric.Counter)

			assert.Equal(t, test.want.gauge, gauge)
			assert.Equal(t, test.want.counter, counter)
		})
	}
}
