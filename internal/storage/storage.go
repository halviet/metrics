package storage

import (
	"errors"
)

type Gauge float64
type Counter int64

type MetricCollection struct {
	gauge   Gauge
	counter Counter
}

func (c *MetricCollection) UpdateGauge(v Gauge) {
	c.gauge = v
}

func (c *MetricCollection) UpdateCounter(v Counter) {
	c.counter += v
}

func (c *MetricCollection) GetGauge() Gauge {
	return c.gauge
}

func (c *MetricCollection) GetCounter() Counter {
	return c.counter
}

type MemStorage struct {
	Collection map[string]MetricCollection
}

func New() *MemStorage {
	collection := make(map[string]MetricCollection)
	return &MemStorage{collection}
}

func (m *MemStorage) UpdateGauge(name string, v Gauge) {
	collection, ok := m.Collection[name]
	if !ok {
		m.Collection[name] = MetricCollection{
			gauge: v,
		}
		return
	}

	collection.UpdateGauge(v)
	m.Collection[name] = collection
}

func (m *MemStorage) UpdateCounter(name string, v Counter) {
	collection, ok := m.Collection[name]
	if !ok {
		m.Collection[name] = MetricCollection{
			counter: v,
		}
		return
	}

	collection.UpdateCounter(v)
	m.Collection[name] = collection
}

var ErrMetricNotFound = errors.New("metric not found")

func (m *MemStorage) GetGauge(name string) (Gauge, error) {
	g, ok := m.Collection[name]
	if !ok {
		return Gauge(0), ErrMetricNotFound
	}

	return g.gauge, nil
}

func (m *MemStorage) GetCounter(name string) (Counter, error) {
	c, ok := m.Collection[name]
	if !ok {
		return Counter(0), ErrMetricNotFound
	}

	return c.counter, nil
}

type ResultMetric struct {
	Gauge   map[string]Gauge
	Counter map[string]Counter
}

func (m *MemStorage) GetAllMetrics() ResultMetric {
	res := ResultMetric{
		Gauge:   map[string]Gauge{},
		Counter: map[string]Counter{},
	}

	for name, collection := range m.Collection {
		if val := collection.GetGauge(); val != Gauge(0) {
			res.Gauge[name] = val
		}
		if val := collection.GetCounter(); val != Counter(0) {
			res.Counter[name] = val
		}
	}

	return res
}
