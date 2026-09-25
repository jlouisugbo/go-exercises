package redemption

import (
	"errors"
	"testing"
)

type testApp struct {
	service *RedemptionService
	store   *MemoryStore
}

type redemptionResult struct {
	receipt Receipt
	err     error
	card    GiftCard
	balance int
}

// These helpers are expected to change as you refactor. Update them to match
// your new API. Tests below this seam should need only narrow changes.
func newTestApp(card GiftCard) *testApp {
	store := NewMemoryStore(card)
	return &testApp{
		service: NewRedemptionService(store),
		store:   store,
	}
}

func (a *testApp) redeem(userID, cardCode string) redemptionResult {
	receipt, err := a.service.Redeem(userID, cardCode)
	return redemptionResult{
		receipt: receipt,
		err:     err,
		card:    a.store.GiftCard("WELCOME25"),
		balance: a.store.Balance("user-123"),
	}
}

func activeCard() GiftCard {
	return GiftCard{Code: "WELCOME25", ValueCents: 2500}
}

func TestRedemptionService(t *testing.T) {
	t.Run("successful redemption", func(t *testing.T) {
		t.Run("returns the redeemed amount", func(t *testing.T) {
			app := newTestApp(activeCard())
			result := app.redeem("user-123", "welcome25")
			if result.receipt.AmountCents != 2500 {
				t.Fatalf("amount = %d, want 2500", result.receipt.AmountCents)
			}
		})

		t.Run("returns the new balance", func(t *testing.T) {
			app := newTestApp(activeCard())
			result := app.redeem("user-123", "WELCOME25")
			if result.receipt.NewBalance != 2500 {
				t.Fatalf("balance = %d, want 2500", result.receipt.NewBalance)
			}
		})

		t.Run("marks the card redeemed", func(t *testing.T) {
			app := newTestApp(activeCard())
			result := app.redeem("user-123", "WELCOME25")
			if !result.card.Redeemed {
				t.Fatal("card was not marked redeemed")
			}
		})

		t.Run("records who redeemed the card", func(t *testing.T) {
			app := newTestApp(activeCard())
			result := app.redeem("user-123", "WELCOME25")
			if result.card.RedeemedBy != "user-123" {
				t.Fatalf("redeemed by = %q, want %q", result.card.RedeemedBy, "user-123")
			}
		})
	})

	t.Run("input and card rules", func(t *testing.T) {
		t.Run("rejects a blank user", func(t *testing.T) {
			app := newTestApp(activeCard())
			result := app.redeem("   ", "WELCOME25")
			if !errors.Is(result.err, ErrInvalidUser) {
				t.Fatalf("error = %v, want ErrInvalidUser", result.err)
			}
		})

		t.Run("rejects an unknown card", func(t *testing.T) {
			app := newTestApp(activeCard())
			result := app.redeem("user-123", "MISSING")
			if !errors.Is(result.err, ErrCardNotFound) {
				t.Fatalf("error = %v, want ErrCardNotFound", result.err)
			}
		})

		t.Run("rejects a redeemed card", func(t *testing.T) {
			card := activeCard()
			card.Redeemed = true
			card.RedeemedBy = "user-456"
			app := newTestApp(card)
			result := app.redeem("user-123", "WELCOME25")
			if !errors.Is(result.err, ErrCardRedeemed) {
				t.Fatalf("error = %v, want ErrCardRedeemed", result.err)
			}
		})

		t.Run("normalizes the card code", func(t *testing.T) {
			app := newTestApp(activeCard())
			result := app.redeem("user-123", "  welcome25 ")
			if result.receipt.CardCode != "WELCOME25" {
				t.Fatalf("card code = %q, want %q", result.receipt.CardCode, "WELCOME25")
			}
		})
	})

	t.Run("dependency failures", func(t *testing.T) {
		t.Run("returns the credit failure", func(t *testing.T) {
			app := newTestApp(activeCard())
			app.store.CreditErr = ErrUnavailable
			result := app.redeem("user-123", "WELCOME25")
			if !errors.Is(result.err, ErrUnavailable) {
				t.Fatalf("error = %v, want ErrUnavailable", result.err)
			}
		})

		t.Run("does not increase the balance when crediting fails", func(t *testing.T) {
			app := newTestApp(activeCard())
			app.store.CreditErr = ErrUnavailable
			result := app.redeem("user-123", "WELCOME25")
			if result.balance != 0 {
				t.Fatalf("balance = %d, want 0", result.balance)
			}
		})

		t.Run("does not redeem the card when crediting fails", func(t *testing.T) {
			app := newTestApp(activeCard())
			app.store.CreditErr = ErrUnavailable
			result := app.redeem("user-123", "WELCOME25")
			if result.card.Redeemed {
				t.Fatal("card was marked redeemed after a failed credit")
			}
			if result.card.RedeemedBy != "" {
				t.Fatalf("redeemed by = %q, want empty", result.card.RedeemedBy)
			}
		})
	})
}
