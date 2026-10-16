package devicereadings

import (
	"errors"
	"testing"
	"time"
)

// fakeStore is a test double for SensorStore with scriptable responses so
// tests can control what the "hardware" would return on each call.
type fakeStore struct {
	readings map[string]Reading
	err      error
	calls    int
}

func (f *fakeStore) FetchReading(deviceID string) (Reading, error) {
	f.calls++
	if f.err != nil {
		return Reading{}, f.err
	}
	r, ok := f.readings[deviceID]
	if !ok {
		return Reading{}, errors.New("unknown device")
	}
	return r, nil
}

// newTestCache and getReading are expected to change as you refactor —
// update them to match your new API.
func newTestCache(store *fakeStore) *ReadingCache {
	return NewReadingCache(store)
}

func getReading(c *ReadingCache, deviceID string) (Reading, error) {
	return c.GetReading(deviceID)
}

// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func TestGetReading(t *testing.T) {
	t.Run("first call for a device", func(t *testing.T) {
		store := &fakeStore{readings: map[string]Reading{
			"thermo-1": {
				DeviceID:    "thermo-1",
				Temperature: 68.5,
				TakenAt:     time.Date(2026, 10, 16, 9, 0, 0, 0, time.UTC),
			},
		}}
		cache := newTestCache(store)

		r, err := getReading(cache, "thermo-1")

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("returns the fetched temperature", func(t *testing.T) {
			if r.Temperature != 68.5 {
				t.Errorf("expected temperature 68.5, got %v", r.Temperature)
			}
		})
		t.Run("hits the store exactly once", func(t *testing.T) {
			if store.calls != 1 {
				t.Errorf("expected 1 store call, got %d", store.calls)
			}
		})
	})

	t.Run("store returns an error", func(t *testing.T) {
		store := &fakeStore{err: errors.New("sensor bus timeout")}
		cache := newTestCache(store)

		_, err := getReading(cache, "thermo-2")

		t.Run("returns an error", func(t *testing.T) {
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	})

	t.Run("second call for the same device", func(t *testing.T) {
		store := &fakeStore{readings: map[string]Reading{
			"thermo-3": {
				DeviceID:    "thermo-3",
				Temperature: 70.0,
				TakenAt:     time.Date(2026, 10, 16, 9, 0, 0, 0, time.UTC),
			},
		}}
		cache := newTestCache(store)

		if _, err := getReading(cache, "thermo-3"); err != nil {
			t.Fatalf("unexpected error priming the cache: %v", err)
		}
		r, err := getReading(cache, "thermo-3")

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("still returns the previously fetched temperature", func(t *testing.T) {
			if r.Temperature != 70.0 {
				t.Errorf("expected temperature 70.0, got %v", r.Temperature)
			}
		})
		t.Run("does not call the store again", func(t *testing.T) {
			if store.calls != 1 {
				t.Errorf("expected store to be called once total, got %d", store.calls)
			}
		})
	})
}
