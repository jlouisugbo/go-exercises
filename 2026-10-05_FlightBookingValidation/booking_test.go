package booking

import (
	"testing"
	"time"
)

// This helper is expected to change as you refactor -- update it to match
// your new API.

func validateBooking(b *TripBooking, stage Stage) error {
	return ValidateBooking(b, stage)
}

// End of helper -- the tests below should not need to change as you
// refactor. But you are welcome to change them if you find you need to!

func TestValidateBooking(t *testing.T) {
	depart := time.Now().AddDate(0, 0, 10)
	ret := depart.AddDate(0, 0, 7)

	t.Run("create", func(t *testing.T) {
		t.Run("missing at sign in email", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "not-an-email", PassengerCount: 2, DepartureDate: depart, ReturnDate: ret, SeatClass: "economy"}
			if err := validateBooking(b, StageCreate); err == nil {
				t.Error("expected an error for an email with no @, got nil")
			}
		})

		t.Run("empty email", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "", PassengerCount: 2, DepartureDate: depart, ReturnDate: ret, SeatClass: "economy"}
			if err := validateBooking(b, StageCreate); err == nil {
				t.Error("expected an error for an empty email, got nil")
			}
		})

		t.Run("zero passengers", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 0, DepartureDate: depart, ReturnDate: ret, SeatClass: "economy"}
			if err := validateBooking(b, StageCreate); err == nil {
				t.Error("expected an error for 0 passengers, got nil")
			}
		})

		t.Run("too many passengers", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 10, DepartureDate: depart, ReturnDate: ret, SeatClass: "economy"}
			if err := validateBooking(b, StageCreate); err == nil {
				t.Error("expected an error for 10 passengers, got nil")
			}
		})

		t.Run("invalid seat class", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 2, DepartureDate: depart, ReturnDate: ret, SeatClass: "first"}
			if err := validateBooking(b, StageCreate); err == nil {
				t.Error("expected an error for an invalid seat class, got nil")
			}
		})

		t.Run("valid booking", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 2, DepartureDate: depart, ReturnDate: ret, SeatClass: "economy"}
			if err := validateBooking(b, StageCreate); err != nil {
				t.Errorf("expected no error for a valid booking, got %v", err)
			}
		})

		t.Run("does not catch a return date before the departure date", func(t *testing.T) {
			// The create branch never looks at the dates at all, so this
			// currently passes. That's the behavior under test right now --
			// it's expected to change once dates are checked consistently
			// across stages.
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 2, DepartureDate: ret, ReturnDate: depart, SeatClass: "economy"}
			if err := validateBooking(b, StageCreate); err != nil {
				t.Errorf("expected no error (create doesn't check dates yet), got %v", err)
			}
		})
	})

	t.Run("confirm", func(t *testing.T) {
		t.Run("zero departure date", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 2, DepartureDate: time.Time{}, ReturnDate: ret, SeatClass: "economy"}
			if err := validateBooking(b, StageConfirm); err == nil {
				t.Error("expected an error for a zero departure date, got nil")
			}
		})

		t.Run("zero return date", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 2, DepartureDate: depart, ReturnDate: time.Time{}, SeatClass: "economy"}
			if err := validateBooking(b, StageConfirm); err == nil {
				t.Error("expected an error for a zero return date, got nil")
			}
		})

		t.Run("return date not after departure date", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 2, DepartureDate: ret, ReturnDate: depart, SeatClass: "economy"}
			if err := validateBooking(b, StageConfirm); err == nil {
				t.Error("expected an error when the return date is before the departure date, got nil")
			}
		})

		t.Run("invalid seat class", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 2, DepartureDate: depart, ReturnDate: ret, SeatClass: "first"}
			if err := validateBooking(b, StageConfirm); err == nil {
				t.Error("expected an error for an invalid seat class, got nil")
			}
		})

		t.Run("valid booking", func(t *testing.T) {
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 2, DepartureDate: depart, ReturnDate: ret, SeatClass: "economy"}
			if err := validateBooking(b, StageConfirm); err != nil {
				t.Errorf("expected no error for a valid booking, got %v", err)
			}
		})

		t.Run("does not catch a zero passenger count", func(t *testing.T) {
			// The confirm branch never looks at PassengerCount, so a booking
			// that never should have been created in the first place still
			// confirms cleanly. Expected to change once the rules are
			// enforced consistently at every stage.
			b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 0, DepartureDate: depart, ReturnDate: ret, SeatClass: "economy"}
			if err := validateBooking(b, StageConfirm); err != nil {
				t.Errorf("expected no error (confirm doesn't check passenger count), got %v", err)
			}
		})
	})

	t.Run("reschedule", func(t *testing.T) {
		t.Run("new departure date in the past", func(t *testing.T) {
			b := &TripBooking{
				PassengerEmail: "a@example.com",
				PassengerCount: 2,
				DepartureDate:  time.Now().AddDate(0, 0, -1),
				ReturnDate:     time.Now().AddDate(0, 0, 5),
				SeatClass:      "economy",
			}
			if err := validateBooking(b, StageReschedule); err == nil {
				t.Error("expected an error for a new departure date in the past, got nil")
			}
		})

		t.Run("new return date not after new departure date", func(t *testing.T) {
			b := &TripBooking{
				PassengerEmail: "a@example.com",
				PassengerCount: 2,
				DepartureDate:  time.Now().AddDate(0, 0, 20),
				ReturnDate:     time.Now().AddDate(0, 0, 15),
				SeatClass:      "economy",
			}
			if err := validateBooking(b, StageReschedule); err == nil {
				t.Error("expected an error when the new return date is before the new departure date, got nil")
			}
		})

		t.Run("valid new dates", func(t *testing.T) {
			b := &TripBooking{
				PassengerEmail: "a@example.com",
				PassengerCount: 2,
				DepartureDate:  time.Now().AddDate(0, 0, 20),
				ReturnDate:     time.Now().AddDate(0, 0, 27),
				SeatClass:      "economy",
			}
			if err := validateBooking(b, StageReschedule); err != nil {
				t.Errorf("expected no error for valid new dates, got %v", err)
			}
		})
	})

	t.Run("unknown stage", func(t *testing.T) {
		b := &TripBooking{PassengerEmail: "a@example.com", PassengerCount: 2, DepartureDate: depart, ReturnDate: ret, SeatClass: "economy"}
		if err := validateBooking(b, Stage(99)); err == nil {
			t.Error("expected an error for an unrecognized stage, got nil")
		}
	})
}
