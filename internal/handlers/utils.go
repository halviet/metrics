package handlers

import (
	"github.com/halviet/metrics/internal/storage"
	"sort"
)

func sortedMetricKeys[T storage.Gauge | storage.Counter](data map[string]T) []string {
	keys := make([]string, 0, len(data))

	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
