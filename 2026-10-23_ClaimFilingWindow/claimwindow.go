// Package claimwindow decides how an auto insurance claim gets paid out
// based on how long after the incident it was filed.
package claimwindow

import (
	"fmt"
	"time"
)

// Claim is a claim as it arrives from the intake form.
type Claim struct {
	ID           string
	IncidentDate time.Time
	ClaimAmount  float64
}

// Decision is what the claims engine decided to pay, and why.
type Decision struct {
	Status       string // "approved", "reduced", or "rejected"
	PayoutAmount float64
	DaysElapsed  int
}

// EvaluateClaim applies the filing-window policy: claims filed within 30
// days of the incident are paid in full, claims filed late are paid at a
// reduced rate that drops the longer the policyholder waited, and claims
// filed more than 90 days out are rejected outright.
func EvaluateClaim(c Claim) (Decision, error) {
	if c.ClaimAmount <= 0 {
		return Decision{}, fmt.Errorf("claim amount must be positive, got %.2f", c.ClaimAmount)
	}

	elapsed := time.Since(c.IncidentDate)
	days := int(elapsed.Hours() / 24)
	if days < 0 {
		return Decision{}, fmt.Errorf("incident date %s is in the future", c.IncidentDate.Format("2006-01-02"))
	}

	switch {
	case days <= 30:
		return Decision{Status: "approved", PayoutAmount: c.ClaimAmount, DaysElapsed: days}, nil
	case days <= 60:
		return Decision{Status: "reduced", PayoutAmount: c.ClaimAmount * 0.8, DaysElapsed: days}, nil
	case days <= 90:
		return Decision{Status: "reduced", PayoutAmount: c.ClaimAmount * 0.5, DaysElapsed: days}, nil
	default:
		return Decision{Status: "rejected", PayoutAmount: 0, DaysElapsed: days}, nil
	}
}
