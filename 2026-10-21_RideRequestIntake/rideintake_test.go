package rideintake

import "testing"

// newStore and submit are expected to change as you refactor — update them
// to match your new API.
func newStore() *RequestStore {
	return &RequestStore{}
}

func submit(store *RequestStore, req RideRequest) (Ride, error) {
	return SubmitRideRequest(store, req)
}

// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func validRequest() RideRequest {
	return RideRequest{
		RequestID:  "req-1",
		RiderID:    "rider-1",
		PickupZone: "downtown",
		FareCents:  1200,
	}
}

func TestSubmitRideRequest(t *testing.T) {
	t.Run("missing rider id", func(t *testing.T) {
		store := newStore()
		req := validRequest()
		req.RiderID = ""

		_, err := submit(store, req)

		t.Run("returns an error", func(t *testing.T) {
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
		t.Run("does not record a ride", func(t *testing.T) {
			if len(store.All()) != 0 {
				t.Errorf("expected no recorded rides, got %d", len(store.All()))
			}
		})
	})

	t.Run("missing pickup zone", func(t *testing.T) {
		store := newStore()
		req := validRequest()
		req.PickupZone = ""

		_, err := submit(store, req)

		t.Run("returns an error", func(t *testing.T) {
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
		t.Run("does not record a ride", func(t *testing.T) {
			if len(store.All()) != 0 {
				t.Errorf("expected no recorded rides, got %d", len(store.All()))
			}
		})
	})

	t.Run("non-positive fare", func(t *testing.T) {
		store := newStore()
		req := validRequest()
		req.FareCents = 0

		_, err := submit(store, req)

		t.Run("returns an error", func(t *testing.T) {
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
		t.Run("does not record a ride", func(t *testing.T) {
			if len(store.All()) != 0 {
				t.Errorf("expected no recorded rides, got %d", len(store.All()))
			}
		})
	})

	t.Run("a valid request", func(t *testing.T) {
		store := newStore()
		req := validRequest()

		ride, err := submit(store, req)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("dispatches the ride", func(t *testing.T) {
			if ride.Status != "dispatched" {
				t.Errorf("expected status %q, got %q", "dispatched", ride.Status)
			}
		})
		t.Run("records the ride in the store", func(t *testing.T) {
			if len(store.All()) != 1 {
				t.Errorf("expected 1 recorded ride, got %d", len(store.All()))
			}
		})
	})

	t.Run("the app retries the same request", func(t *testing.T) {
		store := newStore()
		req := validRequest()

		first, err := submit(store, req)
		if err != nil {
			t.Fatalf("first submission: expected no error, got %v", err)
		}

		second, err := submit(store, req)
		if err != nil {
			t.Fatalf("retried submission: expected no error, got %v", err)
		}

		// This pins today's behavior: a retried request with the same
		// RequestID is dispatched all over again instead of being
		// recognized as the same tap. That's expected to change once a
		// retry returns the original ride instead of creating a new one.
		t.Run("dispatches a second ride with a different ID", func(t *testing.T) {
			if second.ID == first.ID {
				t.Errorf("expected a different ride ID on retry, got the same one: %s", first.ID)
			}
		})
		t.Run("records both rides in the store", func(t *testing.T) {
			if len(store.All()) != 2 {
				t.Errorf("expected 2 recorded rides after a retry, got %d", len(store.All()))
			}
		})
	})
}
