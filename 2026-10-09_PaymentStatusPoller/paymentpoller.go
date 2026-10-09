// Package paymentpoller waits for a submitted payment to reach a
// terminal status by repeatedly polling the payment gateway.
package paymentpoller

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrPollExhausted is returned when a payment never reaches a terminal
// status within the configured number of attempts.
var ErrPollExhausted = errors.New("payment status poll exhausted all attempts without reaching a terminal status")

// PaymentGateway is the subset of the external payment processor's API
// the poller depends on.
type PaymentGateway interface {
	// CheckStatus returns the current status of a previously submitted
	// payment: "pending", "settled", or "failed".
	CheckStatus(paymentID string) (string, error)
}

// StatusPoller repeatedly checks a payment gateway until a payment
// reaches a terminal state, or polling is exhausted.
type StatusPoller struct {
	gateway     PaymentGateway
	interval    time.Duration
	maxAttempts int
}

// NewStatusPoller builds a StatusPoller that checks the gateway every
// interval, up to maxAttempts times, before giving up.
func NewStatusPoller(gateway PaymentGateway, interval time.Duration, maxAttempts int) *StatusPoller {
	return &StatusPoller{
		gateway:     gateway,
		interval:    interval,
		maxAttempts: maxAttempts,
	}
}

// WaitForTerminalStatus blocks until the payment reaches a terminal
// status ("settled" or "failed"), or returns ErrPollExhausted once
// maxAttempts polls have passed without one.
func (p *StatusPoller) WaitForTerminalStatus(ctx context.Context, paymentID string) (string, error) {
	for attempt := 0; attempt < p.maxAttempts; attempt++ {
		status, err := p.gateway.CheckStatus(paymentID)
		if err != nil {
			return "", fmt.Errorf("checking status for payment %s: %w", paymentID, err)
		}

		if status == "settled" {
			return status, nil
		} else if status == "failed" {
			return status, nil
		}

		time.Sleep(p.interval)
	}

	return "", ErrPollExhausted
}
