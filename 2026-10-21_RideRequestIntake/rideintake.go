// Package rideintake accepts ride requests from the rider app and hands
// them off to dispatch.
package rideintake

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// RideRequest is what the rider app sends when someone taps "Request
// Ride." The app tags every request with a RequestID before it leaves the
// device, so retries of the same tap can be told apart from the rider
// tapping again on purpose.
type RideRequest struct {
	RequestID  string
	RiderID    string
	PickupZone string
	FareCents  int
}

// Ride is a ride that's been accepted and handed to the driver-matching
// system.
type Ride struct {
	ID         string
	RiderID    string
	PickupZone string
	FareCents  int
	Status     string
}

var rideSeq int64

func nextRideID() string {
	n := atomic.AddInt64(&rideSeq, 1)
	return fmt.Sprintf("ride_%d", n)
}

// RequestStore is the thinnest possible stand-in for the dispatch
// database -- everything worth looking at here is in how a request turns
// into a ride, not in how it's stored.
type RequestStore struct {
	rides []Ride
}

// Add records a dispatched ride.
func (s *RequestStore) Add(r Ride) {
	s.rides = append(s.rides, r)
}

// All returns every ride recorded so far, oldest first.
func (s *RequestStore) All() []Ride {
	return s.rides
}

// SubmitRideRequest validates an incoming request, dispatches a ride for
// it, and records the result. The rider app retries this call on its own
// whenever the network drops or dispatch answers slowly, so the exact
// same RideRequest -- same RequestID and all -- can land here more than
// once.
func SubmitRideRequest(store *RequestStore, req RideRequest) (Ride, error) {
	if req.RiderID == "" {
		return Ride{}, errors.New("rider id is required")
	}
	if req.PickupZone == "" {
		return Ride{}, errors.New("pickup zone is required")
	}
	if req.FareCents <= 0 {
		return Ride{}, errors.New("fare must be positive")
	}

	ride := Ride{
		ID:         nextRideID(),
		RiderID:    req.RiderID,
		PickupZone: req.PickupZone,
		FareCents:  req.FareCents,
		Status:     "dispatched",
	}

	store.Add(ride)
	return ride, nil
}
