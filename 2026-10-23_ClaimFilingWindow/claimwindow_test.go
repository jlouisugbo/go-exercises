package claimwindow

import (
	"testing"
	"time"
)

// evaluate is expected to change as you refactor — update it to match
// your new API.
func evaluate(amount float64, daysAgo int) (Decision, error) {
	c := Claim{
		ID:           "claim-1",
		IncidentDate: time.Now().Add(-time.Duration(daysAgo) * 24 * time.Hour),
		ClaimAmount:  amount,
	}
	return EvaluateClaim(c)
}

// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func TestEvaluateClaim(t *testing.T) {
	t.Run("non-positive claim amount", func(t *testing.T) {
		_, err := evaluate(0, 10)

		t.Run("returns an error", func(t *testing.T) {
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	})

	t.Run("incident date in the future", func(t *testing.T) {
		_, err := evaluate(500, -5)

		t.Run("returns an error", func(t *testing.T) {
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	})

	t.Run("filed within the 30-day window", func(t *testing.T) {
		decision, err := evaluate(1000, 10)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("is approved", func(t *testing.T) {
			if decision.Status != "approved" {
				t.Errorf("expected status %q, got %q", "approved", decision.Status)
			}
		})
		t.Run("pays the full amount", func(t *testing.T) {
			if decision.PayoutAmount != 1000 {
				t.Errorf("expected payout 1000, got %.2f", decision.PayoutAmount)
			}
		})
	})

	t.Run("filed in the 31-60 day window", func(t *testing.T) {
		decision, err := evaluate(1000, 45)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("is reduced", func(t *testing.T) {
			if decision.Status != "reduced" {
				t.Errorf("expected status %q, got %q", "reduced", decision.Status)
			}
		})
		t.Run("pays 80 percent", func(t *testing.T) {
			if decision.PayoutAmount != 800 {
				t.Errorf("expected payout 800, got %.2f", decision.PayoutAmount)
			}
		})
	})

	t.Run("filed in the 61-90 day window", func(t *testing.T) {
		decision, err := evaluate(1000, 75)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("is reduced", func(t *testing.T) {
			if decision.Status != "reduced" {
				t.Errorf("expected status %q, got %q", "reduced", decision.Status)
			}
		})
		t.Run("pays 50 percent", func(t *testing.T) {
			if decision.PayoutAmount != 500 {
				t.Errorf("expected payout 500, got %.2f", decision.PayoutAmount)
			}
		})
	})

	t.Run("filed more than 90 days out", func(t *testing.T) {
		decision, err := evaluate(1000, 120)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("is rejected", func(t *testing.T) {
			if decision.Status != "rejected" {
				t.Errorf("expected status %q, got %q", "rejected", decision.Status)
			}
		})
		t.Run("pays nothing", func(t *testing.T) {
			if decision.PayoutAmount != 0 {
				t.Errorf("expected payout 0, got %.2f", decision.PayoutAmount)
			}
		})
	})
}
