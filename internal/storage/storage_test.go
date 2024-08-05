package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_GetAllMetrics(t *testing.T) {
	store := New()

	store.UpdateGauge("gaugeVal", Gauge(100.01))
	store.UpdateCounter("counterVal", Counter(5))

	res := store.GetAllMetrics()
	assert.Equal(t, Gauge(100.01), res.Gauge["gaugeVal"])
	assert.Equal(t, Counter(5), res.Counter["counterVal"])
}
