// Package booking implements trip booking validation for a small travel
// agency's internal tool.
package booking

import (
	"errors"
	"strings"
	"time"
)

// Stage identifies which point in a booking's lifecycle it's being
// validated for.
type Stage int

const (
	StageCreate Stage = iota
	StageConfirm
	StageReschedule
)

// TripBooking represents a round-trip flight booking for a single party.
type TripBooking struct {
	PassengerEmail string
	PassengerCount int
	DepartureDate  time.Time
	ReturnDate     time.Time
	SeatClass      string // "economy", "premium", "business"
	CorporateID    string // required for group bookings of 6+ passengers
}

var validSeatClasses = map[string]bool{
	"economy":  true,
	"premium":  true,
	"business": true,
}

// ValidateBooking checks b against the rules for the given lifecycle stage.
// The booking tool calls this once at each of the three stages a booking
// passes through: when it's first created, when payment clears and it's
// confirmed, and whenever it's rescheduled. For a reschedule check, the
// caller is expected to pass a copy of the booking with the proposed new
// dates already applied.
func ValidateBooking(b *TripBooking, stage Stage) error {
	switch stage {
	case StageCreate:
		if b.PassengerEmail == "" || !strings.Contains(b.PassengerEmail, "@") {
			return errors.New("invalid passenger email")
		}
		if b.PassengerCount < 1 || b.PassengerCount > 9 {
			return errors.New("passenger count must be between 1 and 9")
		}
		if !validSeatClasses[b.SeatClass] {
			return errors.New("invalid seat class")
		}
		return nil
	case StageConfirm:
		if b.DepartureDate.IsZero() || b.ReturnDate.IsZero() {
			return errors.New("both departure and return dates are required")
		}
		if !b.ReturnDate.After(b.DepartureDate) {
			return errors.New("return date must be after departure date")
		}
		if !validSeatClasses[b.SeatClass] {
			return errors.New("invalid seat class")
		}
		return nil
	case StageReschedule:
		if b.DepartureDate.Before(time.Now()) {
			return errors.New("new departure date must be in the future")
		}
		if !b.ReturnDate.After(b.DepartureDate) {
			return errors.New("new return date must be after new departure date")
		}
		return nil
	default:
		return errors.New("unknown booking stage")
	}
}
