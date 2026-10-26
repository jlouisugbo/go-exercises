package maintenance

import (
	"testing"
	"time"
)

// scheduleService is expected to change as you refactor — update it to
// match your new API.
func scheduleService(equipmentType string, lastServiced time.Time) (time.Time, error) {
	return ScheduleNextService(equipmentType, lastServiced)
}

// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func TestScheduleNextService(t *testing.T) {
	lastServiced := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("empty equipment type", func(t *testing.T) {
		_, err := scheduleService("", lastServiced)

		t.Run("returns an error", func(t *testing.T) {
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	})

	t.Run("tractor", func(t *testing.T) {
		got, err := scheduleService("tractor", lastServiced)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("shortens the interval to 80 days", func(t *testing.T) {
			want := lastServiced.AddDate(0, 0, 80)
			if !got.Equal(want) {
				t.Errorf("expected %v, got %v", want, got)
			}
		})
	})

	t.Run("harvester", func(t *testing.T) {
		got, err := scheduleService("harvester", lastServiced)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("shortens the interval to 70 days", func(t *testing.T) {
			want := lastServiced.AddDate(0, 0, 70)
			if !got.Equal(want) {
				t.Errorf("expected %v, got %v", want, got)
			}
		})
	})

	t.Run("irrigation_pump", func(t *testing.T) {
		got, err := scheduleService("irrigation_pump", lastServiced)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("lengthens the interval to 105 days", func(t *testing.T) {
			want := lastServiced.AddDate(0, 0, 105)
			if !got.Equal(want) {
				t.Errorf("expected %v, got %v", want, got)
			}
		})
	})

	t.Run("unrecognized equipment type", func(t *testing.T) {
		// The current implementation has no case for an equipment type it
		// doesn't recognize, so the lookup returns a nil strategy and the
		// call below panics. This is expected to become a returned error
		// once you refactor.
		t.Run("panics instead of returning an error", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("expected a panic for an unrecognized equipment type, got none")
				}
			}()
			_, _ = scheduleService("drone", lastServiced)
		})
	})
}
