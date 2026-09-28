package delivery

import (
	"context"
	"errors"
	"testing"
)

type stubGateway struct {
	results  []error
	calls    int
	invoices []Invoice
}

func (s *stubGateway) Send(_ context.Context, invoice Invoice) error {
	s.calls++
	s.invoices = append(s.invoices, invoice)
	if len(s.results) == 0 {
		return nil
	}

	err := s.results[0]
	s.results = s.results[1:]
	return err
}

type testApp struct {
	worker  *Worker
	gateway *stubGateway
}

type deliveryResult struct {
	err      error
	calls    int
	invoices []Invoice
}

// This is the construction and public-call seam. If your refactor changes
// constructors or method signatures, adapt these helpers before changing tests.
func newTestApp(maxAttempts int, results ...error) *testApp {
	gateway := &stubGateway{results: append([]error(nil), results...)}
	return &testApp{
		worker:  NewWorker(gateway, maxAttempts),
		gateway: gateway,
	}
}

func (a *testApp) deliver(ctx context.Context, invoice Invoice) deliveryResult {
	err := a.worker.Deliver(ctx, invoice)
	return deliveryResult{
		err:      err,
		calls:    a.gateway.calls,
		invoices: append([]Invoice(nil), a.gateway.invoices...),
	}
}

func validInvoice() Invoice {
	return Invoice{ID: "invoice-123", Recipient: "billing@example.com"}
}

func temporaryGatewayError() error {
	return &GatewayError{StatusCode: 503, Err: ErrGatewayUnavailable}
}

func TestWorkerDeliver(t *testing.T) {
	t.Run("successful delivery", func(t *testing.T) {
		t.Run("sends the invoice once", func(t *testing.T) {
			app := newTestApp(3)
			result := app.deliver(context.Background(), validInvoice())
			if result.err != nil {
				t.Fatalf("deliver: %v", result.err)
			}
			if result.calls != 1 {
				t.Fatalf("calls = %d, want 1", result.calls)
			}
		})

		t.Run("normalizes invoice fields", func(t *testing.T) {
			app := newTestApp(3)
			result := app.deliver(context.Background(), Invoice{
				ID:        "  invoice-123 ",
				Recipient: " billing@example.com ",
			})
			if result.err != nil {
				t.Fatalf("deliver: %v", result.err)
			}
			if got := result.invoices[0].ID; got != "invoice-123" {
				t.Fatalf("invoice id = %q, want %q", got, "invoice-123")
			}
			if got := result.invoices[0].Recipient; got != "billing@example.com" {
				t.Fatalf("recipient = %q, want %q", got, "billing@example.com")
			}
		})
	})

	t.Run("validation", func(t *testing.T) {
		tests := []struct {
			name    string
			invoice Invoice
		}{
			{name: "blank invoice id", invoice: Invoice{Recipient: "billing@example.com"}},
			{name: "blank recipient", invoice: Invoice{ID: "invoice-123"}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				app := newTestApp(3)
				result := app.deliver(context.Background(), tt.invoice)
				if !errors.Is(result.err, ErrInvalidInvoice) {
					t.Fatalf("error = %v, want ErrInvalidInvoice", result.err)
				}
				if result.calls != 0 {
					t.Fatalf("calls = %d, want 0", result.calls)
				}
			})
		}
	})

	t.Run("repeated attempts", func(t *testing.T) {
		t.Run("eventually succeeds after temporary failures", func(t *testing.T) {
			app := newTestApp(3, temporaryGatewayError(), temporaryGatewayError(), nil)
			result := app.deliver(context.Background(), validInvoice())
			if result.err != nil {
				t.Fatalf("deliver: %v", result.err)
			}
			if result.calls != 3 {
				t.Fatalf("calls = %d, want 3", result.calls)
			}
		})

		t.Run("stops at the configured limit", func(t *testing.T) {
			app := newTestApp(3, temporaryGatewayError(), temporaryGatewayError(), temporaryGatewayError())
			result := app.deliver(context.Background(), validInvoice())
			if result.calls != 3 {
				t.Fatalf("calls = %d, want 3", result.calls)
			}
		})

		t.Run("preserves the underlying failure", func(t *testing.T) {
			app := newTestApp(2, temporaryGatewayError(), temporaryGatewayError())
			result := app.deliver(context.Background(), validInvoice())
			if !errors.Is(result.err, ErrGatewayUnavailable) {
				t.Fatalf("error = %v, want ErrGatewayUnavailable", result.err)
			}
		})

		t.Run("uses one attempt when configured with zero", func(t *testing.T) {
			app := newTestApp(0, temporaryGatewayError())
			result := app.deliver(context.Background(), validInvoice())
			if result.calls != 1 {
				t.Fatalf("calls = %d, want 1", result.calls)
			}
		})
	})
}
