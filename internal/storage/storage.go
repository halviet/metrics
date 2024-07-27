package storage

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

func (m *MemStorage) GetGauge(name string) Gauge {
	return m.Collection[name].gauge
}

func (m *MemStorage) GetCounter(name string) Counter {
	return m.Collection[name].counter
}
