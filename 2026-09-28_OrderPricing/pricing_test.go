package pricing

import (
	"errors"
	"testing"
)

func calc(items []Item, code string) (OrderResult, error) {
	return CalculateOrder(items, code)
}

func TestCalculateOrder(t *testing.T) {
	t.Run("no discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: 2}}
		result, err := calc(items, "")

		t.Run("does not error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("did not expect an error: %v", err)
			}
		})
		t.Run("returns correct subtotal", func(t *testing.T) {
			if result.Subtotal != 20 {
				t.Errorf("subtotal = %v, want 20", result.Subtotal)
			}
		})
	})

	t.Run("SAVE10 discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, err := calc(items, "SAVE10")
		if err != nil {
			t.Fatal(err)
		}

		t.Run("returns correct discount", func(t *testing.T) {
			if result.Discount != 10 {
				t.Errorf("discount = %v, want 10", result.Discount)
			}
		})
	})

	t.Run("SAVE20 discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, err := calc(items, "SAVE20")
		if err != nil {
			t.Fatal(err)
		}

		t.Run("returns correct discount", func(t *testing.T) {
			if result.Discount != 20 {
				t.Errorf("discount = %v, want 20", result.Discount)
			}
		})
	})

	t.Run("VIP discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, err := calc(items, "VIP")
		if err != nil {
			t.Fatal(err)
		}

		t.Run("returns correct discount", func(t *testing.T) {
			if result.Discount != 30 {
				t.Errorf("discount = %v, want 30", result.Discount)
			}
		})
	})

	t.Run("tax applied after discount", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, err := calc(items, "SAVE10")
		if err != nil {
			t.Fatal(err)
		}
		wantTax := (100 - 10) * TaxRate

		t.Run("returns correct tax", func(t *testing.T) {
			if result.Tax != wantTax {
				t.Errorf("tax = %v, want %v", result.Tax, wantTax)
			}
		})
	})

	t.Run("empty order", func(t *testing.T) {
		_, err := calc(nil, "")

		t.Run("returns ErrEmptyOrder", func(t *testing.T) {
			if !errors.Is(err, ErrEmptyOrder) {
				t.Errorf("expected ErrEmptyOrder, got %v", err)
			}
		})
	})

	t.Run("negative quantity", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: -1}}
		_, err := calc(items, "")

		t.Run("returns ErrNegativeQuantity", func(t *testing.T) {
			if !errors.Is(err, ErrNegativeQuantity) {
				t.Errorf("expected ErrNegativeQuantity, got %v", err)
			}
		})
	})

	t.Run("unknown discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: 1}}
		_, err := calc(items, "NOTREAL")

		t.Run("returns UnknownDiscountError", func(t *testing.T) {
			var discErr UnknownDiscountError
			if !errors.As(err, &discErr) {
				t.Fatalf("expected UnknownDiscountError, got %T: %v", err, err)
			}
			if discErr.Code != "NOTREAL" {
				t.Errorf("code = %q, want NOTREAL", discErr.Code)
			}
		})
	})
}

func TestCalculatorIndependentTaxRates(t *testing.T) {
	items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}

	lowTax := NewCalculator(0.05)
	highTax := NewCalculator(0.15)

	low, err := lowTax.CalculateOrder(items, "")
	if err != nil {
		t.Fatal(err)
	}
	high, err := highTax.CalculateOrder(items, "")
	if err != nil {
		t.Fatal(err)
	}

	if low.Tax == high.Tax {
		t.Fatalf("expected different tax amounts, both got %v", low.Tax)
	}
	if low.Tax != 5.0 {
		t.Errorf("low tax = %v, want 5", low.Tax)
	}
	if high.Tax != 15.0 {
		t.Errorf("high tax = %v, want 15", high.Tax)
	}
}

func TestRegisterDiscount(t *testing.T) {
	c := NewCalculator(0.08)
	c.RegisterDiscount("HALF", 0.50)

	items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
	result, err := c.CalculateOrder(items, "HALF")
	if err != nil {
		t.Fatal(err)
	}
	if result.Discount != 50 {
		t.Errorf("discount = %v, want 50", result.Discount)
	}
}
