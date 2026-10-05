package booking

import (
	"testing"
	"time"
)

// These helpers are expected to change as you refactor -- update them to
// match your new API.

func newTestBooking(email string, count int, depart, ret time.Time, seatClass string) *TripBooking {
	return NewTripBooking("b1", email, count, depart, ret, seatClass, "")
}

func validateCreate(b *TripBooking) error {
	return ValidateForCreate(b)
}

func validateConfirm(b *TripBooking) error {
	return ValidateForConfirm(b)
}

func validateReschedule(b *TripBooking, newDepart, newReturn time.Time) error {
	return ValidateForReschedule(b, newDepart, newReturn)
}

func newTestStore() *Store {
	return NewStore()
}

func createBooking(store *Store, b *TripBooking) error {
	return CreateBooking(store, b)
}

func confirmBooking(store *Store, id string) error {
	return ConfirmBooking(store, id)
}

func rescheduleBooking(store *Store, id string, newDepart, newReturn time.Time) error {
	return RescheduleBooking(store, id, newDepart, newReturn)
}

// End of helper -- the tests below should not need to change as you
// refactor. But you are welcome to change them if you find you need to!

func TestValidateForCreate(t *testing.T) {
	depart := time.Now().AddDate(0, 0, 10)
	ret := depart.AddDate(0, 0, 7)

	t.Run("missing at sign in email", func(t *testing.T) {
		b := newTestBooking("not-an-email", 2, depart, ret, "economy")
		if err := validateCreate(b); err == nil {
			t.Error("expected an error for an email with no @, got nil")
		}
	})

	t.Run("empty email", func(t *testing.T) {
		b := newTestBooking("", 2, depart, ret, "economy")
		if err := validateCreate(b); err == nil {
			t.Error("expected an error for an empty email, got nil")
		}
	})

	t.Run("zero passengers", func(t *testing.T) {
		b := newTestBooking("a@example.com", 0, depart, ret, "economy")
		if err := validateCreate(b); err == nil {
			t.Error("expected an error for 0 passengers, got nil")
		}
	})

	t.Run("too many passengers", func(t *testing.T) {
		b := newTestBooking("a@example.com", 10, depart, ret, "economy")
		if err := validateCreate(b); err == nil {
			t.Error("expected an error for 10 passengers, got nil")
		}
	})

	t.Run("invalid seat class", func(t *testing.T) {
		b := newTestBooking("a@example.com", 2, depart, ret, "first")
		if err := validateCreate(b); err == nil {
			t.Error("expected an error for an invalid seat class, got nil")
		}
	})

	t.Run("valid booking", func(t *testing.T) {
		b := newTestBooking("a@example.com", 2, depart, ret, "economy")
		if err := validateCreate(b); err != nil {
			t.Errorf("expected no error for a valid booking, got %v", err)
		}
	})

	t.Run("accepts a return date before the departure date", func(t *testing.T) {
		// ValidateForCreate never looks at the dates at all, so this
		// currently passes. That's the behavior under test right now --
		// it's expected to change once dates are checked consistently.
		b := newTestBooking("a@example.com", 2, ret, depart, "economy")
		if err := validateCreate(b); err != nil {
			t.Errorf("expected no error (create doesn't check dates yet), got %v", err)
		}
	})
}

func TestValidateForConfirm(t *testing.T) {
	depart := time.Now().AddDate(0, 0, 10)
	ret := depart.AddDate(0, 0, 7)

	t.Run("zero departure date", func(t *testing.T) {
		b := newTestBooking("a@example.com", 2, time.Time{}, ret, "economy")
		if err := validateConfirm(b); err == nil {
			t.Error("expected an error for a zero departure date, got nil")
		}
	})

	t.Run("zero return date", func(t *testing.T) {
		b := newTestBooking("a@example.com", 2, depart, time.Time{}, "economy")
		if err := validateConfirm(b); err == nil {
			t.Error("expected an error for a zero return date, got nil")
		}
	})

	t.Run("return date not after departure date", func(t *testing.T) {
		b := newTestBooking("a@example.com", 2, ret, depart, "economy")
		if err := validateConfirm(b); err == nil {
			t.Error("expected an error when the return date is before the departure date, got nil")
		}
	})

	t.Run("invalid seat class", func(t *testing.T) {
		b := newTestBooking("a@example.com", 2, depart, ret, "first")
		if err := validateConfirm(b); err == nil {
			t.Error("expected an error for an invalid seat class, got nil")
		}
	})

	t.Run("valid booking", func(t *testing.T) {
		b := newTestBooking("a@example.com", 2, depart, ret, "economy")
		if err := validateConfirm(b); err != nil {
			t.Errorf("expected no error for a valid booking, got %v", err)
		}
	})

	t.Run("does not re-check passenger count", func(t *testing.T) {
		// ValidateForConfirm never looks at PassengerCount, so a booking
		// that never should have been created in the first place still
		// confirms cleanly. Expected to change once the rules are
		// enforced consistently at every stage.
		b := newTestBooking("a@example.com", 0, depart, ret, "economy")
		if err := validateConfirm(b); err != nil {
			t.Errorf("expected no error (confirm doesn't check passenger count), got %v", err)
		}
	})
}

func TestValidateForReschedule(t *testing.T) {
	depart := time.Now().AddDate(0, 0, 10)
	ret := depart.AddDate(0, 0, 7)
	existing := newTestBooking("a@example.com", 2, depart, ret, "economy")

	t.Run("new departure date in the past", func(t *testing.T) {
		pastDepart := time.Now().AddDate(0, 0, -1)
		futureReturn := time.Now().AddDate(0, 0, 5)
		if err := validateReschedule(existing, pastDepart, futureReturn); err == nil {
			t.Error("expected an error for a new departure date in the past, got nil")
		}
	})

	t.Run("new return date not after new departure date", func(t *testing.T) {
		newDepart := time.Now().AddDate(0, 0, 20)
		newReturn := time.Now().AddDate(0, 0, 15)
		if err := validateReschedule(existing, newDepart, newReturn); err == nil {
			t.Error("expected an error when the new return date is before the new departure date, got nil")
		}
	})

	t.Run("valid new dates", func(t *testing.T) {
		newDepart := time.Now().AddDate(0, 0, 20)
		newReturn := time.Now().AddDate(0, 0, 27)
		if err := validateReschedule(existing, newDepart, newReturn); err != nil {
			t.Errorf("expected no error for valid new dates, got %v", err)
		}
	})
}

func TestCreateBooking(t *testing.T) {
	depart := time.Now().AddDate(0, 0, 10)
	ret := depart.AddDate(0, 0, 7)

	t.Run("rejects an invalid booking", func(t *testing.T) {
		store := newTestStore()
		b := newTestBooking("not-an-email", 2, depart, ret, "economy")
		if err := createBooking(store, b); err == nil {
			t.Error("expected an error for an invalid booking, got nil")
		}
	})

	t.Run("saves a valid booking", func(t *testing.T) {
		store := newTestStore()
		b := newTestBooking("a@example.com", 2, depart, ret, "economy")
		err := createBooking(store, b)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Errorf("expected no error for a valid booking, got %v", err)
			}
		})

		t.Run("the booking is findable in the store", func(t *testing.T) {
			if found := store.Find(b.ID); found == nil {
				t.Error("expected the booking to be saved, but it wasn't found in the store")
			}
		})
	})
}

func TestConfirmBooking(t *testing.T) {
	depart := time.Now().AddDate(0, 0, 10)
	ret := depart.AddDate(0, 0, 7)

	t.Run("booking not found", func(t *testing.T) {
		store := newTestStore()
		if err := confirmBooking(store, "missing"); err == nil {
			t.Error("expected an error for a booking that doesn't exist, got nil")
		}
	})

	t.Run("confirms a valid booking", func(t *testing.T) {
		store := newTestStore()
		b := newTestBooking("a@example.com", 2, depart, ret, "economy")
		store.Save(b)
		err := confirmBooking(store, b.ID)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})

		t.Run("marks the booking confirmed", func(t *testing.T) {
			if !b.Confirmed {
				t.Error("expected the booking to be marked confirmed")
			}
		})
	})

	t.Run("rejects a booking that fails confirm validation", func(t *testing.T) {
		store := newTestStore()
		b := newTestBooking("a@example.com", 2, ret, depart, "economy") // return before departure
		store.Save(b)
		if err := confirmBooking(store, b.ID); err == nil {
			t.Error("expected an error for a booking with return date before departure date, got nil")
		}
	})
}

func TestRescheduleBooking(t *testing.T) {
	depart := time.Now().AddDate(0, 0, 10)
	ret := depart.AddDate(0, 0, 7)

	t.Run("booking not found", func(t *testing.T) {
		store := newTestStore()
		newDepart := time.Now().AddDate(0, 0, 20)
		newReturn := time.Now().AddDate(0, 0, 27)
		if err := rescheduleBooking(store, "missing", newDepart, newReturn); err == nil {
			t.Error("expected an error for a booking that doesn't exist, got nil")
		}
	})

	t.Run("applies valid new dates", func(t *testing.T) {
		store := newTestStore()
		b := newTestBooking("a@example.com", 2, depart, ret, "economy")
		store.Save(b)
		newDepart := time.Now().AddDate(0, 0, 20)
		newReturn := time.Now().AddDate(0, 0, 27)
		err := rescheduleBooking(store, b.ID, newDepart, newReturn)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})

		t.Run("updates the departure date", func(t *testing.T) {
			if !b.DepartureDate.Equal(newDepart) {
				t.Error("expected the departure date to be updated")
			}
		})
	})

	t.Run("rejects invalid new dates", func(t *testing.T) {
		store := newTestStore()
		b := newTestBooking("a@example.com", 2, depart, ret, "economy")
		store.Save(b)
		pastDepart := time.Now().AddDate(0, 0, -1)
		if err := rescheduleBooking(store, b.ID, pastDepart, ret); err == nil {
			t.Error("expected an error for a new departure date in the past, got nil")
		}
	})
}
