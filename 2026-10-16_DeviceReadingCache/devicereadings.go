package devicereadings

import (
	"fmt"
	"time"
)

// Reading is a single sensor measurement pulled from building hardware.
type Reading struct {
	DeviceID    string
	Temperature float64
	TakenAt     time.Time
}

// SensorStore is the slow hardware-backed lookup. Hitting it for every
// dashboard refresh floods the building's sensor bus, so callers are
// expected to cache.
type SensorStore interface {
	FetchReading(deviceID string) (Reading, error)
}

// ReadingCache sits in front of a SensorStore so the building dashboard
// (which polls every few seconds, across every floor) doesn't hit the
// hardware bus directly on every request.
type ReadingCache struct {
	store  SensorStore
	cached map[string]Reading
}

// NewReadingCache wires a cache up to its backing store.
func NewReadingCache(store SensorStore) *ReadingCache {
	return &ReadingCache{
		store:  store,
		cached: make(map[string]Reading),
	}
}

// GetReading returns the latest reading for a device, serving from cache
// when possible so the sensor bus doesn't get hammered by every dashboard
// tab refreshing in parallel.
func (c *ReadingCache) GetReading(deviceID string) (Reading, error) {
	if r, ok := c.cached[deviceID]; ok {
		return r, nil
	}

	r, err := c.store.FetchReading(deviceID)
	if err != nil {
		return Reading{}, fmt.Errorf("fetch reading for %s: %w", deviceID, err)
	}

	c.cached[deviceID] = r
	return r, nil
}
