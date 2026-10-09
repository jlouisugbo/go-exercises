package paymentpoller

import (
	"context"
	"errors"
	"testing"
	"time"
)

// fakeGateway returns a scripted sequence of statuses to CheckStatus, one
// per call, repeating the final entry once the script runs out. If err is
// set, every call returns that error instead.
type fakeGateway struct {
	statuses []string
	err      error
	calls    int
}

func (f *fakeGateway) CheckStatus(paymentID string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	idx := f.calls
	if idx >= len(f.statuses) {
		idx = len(f.statuses) - 1
	}
	f.calls++
	return f.statuses[idx], nil
}

// This helper is expected to change as you refactor — update it to match your new API.
func newTestPoller(gw PaymentGateway) *StatusPoller {
	return NewStatusPoller(gw, time.Millisecond, 3)
}

// This helper is expected to change as you refactor — update it to match your new API.
func pollStatus(ctx context.Context, p *StatusPoller, paymentID string) (string, error) {
	return p.WaitForTerminalStatus(ctx, paymentID)
}

// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func TestWaitForTerminalStatus(t *testing.T) {
	t.Run("settled on first check", func(t *testing.T) {
		gw := &fakeGateway{statuses: []string{"settled"}}
		p := newTestPoller(gw)

		status, err := pollStatus(context.Background(), p, "pay_1")

		t.Run("returns settled", func(t *testing.T) {
			if status != "settled" {
				t.Fatalf("status = %q, want %q", status, "settled")
			}
		})
		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("err = %v, want nil", err)
			}
		})
	})

	t.Run("failed on first check", func(t *testing.T) {
		gw := &fakeGateway{statuses: []string{"failed"}}
		p := newTestPoller(gw)

		status, err := pollStatus(context.Background(), p, "pay_2")

		t.Run("returns failed", func(t *testing.T) {
			if status != "failed" {
				t.Fatalf("status = %q, want %q", status, "failed")
			}
		})
		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("err = %v, want nil", err)
			}
		})
	})

	t.Run("pending then settled", func(t *testing.T) {
		gw := &fakeGateway{statuses: []string{"pending", "pending", "settled"}}
		p := newTestPoller(gw)

		status, err := pollStatus(context.Background(), p, "pay_3")

		t.Run("returns settled", func(t *testing.T) {
			if status != "settled" {
				t.Fatalf("status = %q, want %q", status, "settled")
			}
		})
		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("err = %v, want nil", err)
			}
		})
		t.Run("checked the gateway three times", func(t *testing.T) {
			if gw.calls != 3 {
				t.Fatalf("calls = %d, want 3", gw.calls)
			}
		})
	})

	t.Run("gateway error", func(t *testing.T) {
		gw := &fakeGateway{err: errors.New("gateway unreachable")}
		p := newTestPoller(gw)

		_, err := pollStatus(context.Background(), p, "pay_4")

		t.Run("returns a non-nil error", func(t *testing.T) {
			if err == nil {
				t.Fatal("err = nil, want non-nil")
			}
		})
	})

	t.Run("never reaches a terminal status", func(t *testing.T) {
		gw := &fakeGateway{statuses: []string{"pending"}}
		p := newTestPoller(gw)

		_, err := pollStatus(context.Background(), p, "pay_5")

		t.Run("returns ErrPollExhausted", func(t *testing.T) {
			if !errors.Is(err, ErrPollExhausted) {
				t.Fatalf("err = %v, want ErrPollExhausted", err)
			}
		})
	})
}
