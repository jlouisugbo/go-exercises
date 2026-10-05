// Package booking implements trip booking validation for a small travel
// agency's internal tool.
package booking

import (
	"errors"
	"strings"
	"time"
)

// TripBooking represents a round-trip flight booking for a single party.
type TripBooking struct {
	ID             string
	PassengerEmail string
	PassengerCount int
	DepartureDate  time.Time
	ReturnDate     time.Time
	SeatClass      string // "economy", "premium", "business"
	CorporateID    string // required for group bookings of 6+ passengers
	Confirmed      bool
}

// NewTripBooking builds a booking from raw request fields. It doesn't
// validate anything itself -- that happens in whichever Validate* function
// the caller remembers to call at each stage of the booking lifecycle.
func NewTripBooking(id, email string, count int, depart, ret time.Time, seatClass, corporateID string) *TripBooking {
	return &TripBooking{
		ID:             id,
		PassengerEmail: email,
		PassengerCount: count,
		DepartureDate:  depart,
		ReturnDate:     ret,
		SeatClass:      seatClass,
		CorporateID:    corporateID,
	}
}

var validSeatClasses = map[string]bool{
	"economy":  true,
	"premium":  true,
	"business": true,
}

// ValidateForCreate checks a booking is well-formed before it's saved.
func ValidateForCreate(b *TripBooking) error {
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
}

// ValidateForConfirm checks a booking is ready to lock in with the airline
// once payment has cleared.
func ValidateForConfirm(b *TripBooking) error {
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
}

// ValidateForReschedule checks proposed new dates before they're applied to
// an existing booking.
func ValidateForReschedule(b *TripBooking, newDepart, newReturn time.Time) error {
	if newDepart.Before(time.Now()) {
		return errors.New("new departure date must be in the future")
	}
	if !newReturn.After(newDepart) {
		return errors.New("new return date must be after new departure date")
	}
	return nil
}

// Store is a minimal in-memory booking store.
type Store struct {
	bookings map[string]*TripBooking
}

// NewStore returns an empty Store.
func NewStore() *Store {
	return &Store{bookings: make(map[string]*TripBooking)}
}

// Save persists a booking, overwriting any existing booking with the same ID.
func (s *Store) Save(b *TripBooking) {
	s.bookings[b.ID] = b
}

// Find returns the booking with the given ID, or nil if it doesn't exist.
func (s *Store) Find(id string) *TripBooking {
	return s.bookings[id]
}

// CreateBooking validates and stores a new booking.
func CreateBooking(store *Store, b *TripBooking) error {
	if err := ValidateForCreate(b); err != nil {
		return err
	}
	store.Save(b)
	return nil
}

// ConfirmBooking validates a booking is ready and marks it confirmed.
func ConfirmBooking(store *Store, id string) error {
	b := store.Find(id)
	if b == nil {
		return errors.New("booking not found")
	}
	if err := ValidateForConfirm(b); err != nil {
		return err
	}
	b.Confirmed = true
	return nil
}

// RescheduleBooking validates and applies new dates to an existing booking.
func RescheduleBooking(store *Store, id string, newDepart, newReturn time.Time) error {
	b := store.Find(id)
	if b == nil {
		return errors.New("booking not found")
	}
	if err := ValidateForReschedule(b, newDepart, newReturn); err != nil {
		return err
	}
	b.DepartureDate = newDepart
	b.ReturnDate = newReturn
	return nil
}
