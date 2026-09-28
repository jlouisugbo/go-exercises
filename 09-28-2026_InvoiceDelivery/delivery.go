package delivery

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidInvoice     = errors.New("invalid invoice")
	ErrGatewayUnavailable = errors.New("delivery gateway unavailable")
	ErrRecipientRejected  = errors.New("recipient rejected")
)

type Invoice struct {
	ID        string
	Recipient string
}

type Gateway interface {
	Send(ctx context.Context, invoice Invoice) error
}

type GatewayError struct {
	StatusCode int
	Err        error
}

func (e *GatewayError) Error() string {
	return fmt.Sprintf("gateway returned status %d: %v", e.StatusCode, e.Err)
}

func (e *GatewayError) Unwrap() error {
	return e.Err
}

type Worker struct {
	gateway     Gateway
	maxAttempts int
}

func NewWorker(gateway Gateway, maxAttempts int) *Worker {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	return &Worker{gateway: gateway, maxAttempts: maxAttempts}
}

func (w *Worker) Deliver(ctx context.Context, invoice Invoice) error {
	invoice.ID = strings.TrimSpace(invoice.ID)
	invoice.Recipient = strings.TrimSpace(invoice.Recipient)
	if invoice.ID == "" || invoice.Recipient == "" {
		return ErrInvalidInvoice
	}

	var lastErr error
	for attempt := 1; attempt <= w.maxAttempts; attempt++ {
		if err := w.gateway.Send(ctx, invoice); err != nil {
			lastErr = err
			continue
		}
		return nil
	}

	return fmt.Errorf("deliver invoice after %d attempts: %w", w.maxAttempts, lastErr)
}
