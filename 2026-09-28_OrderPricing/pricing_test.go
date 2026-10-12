package pricing

import "testing"

// This helper is expected to change as you refactor — update it to match your new API.
func calc(items []Item, code string) (result OrderResult, panicked bool, panicValue any) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicValue = r
		}
	}()
	result = CalculateOrder(items, code)
	return result, false, nil
}

// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func TestCalculateOrder(t *testing.T) {
	t.Run("no discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: 2}}
		result, panicked, _ := calc(items, "")

		t.Run("does not panic", func(t *testing.T) {
			if panicked {
				t.Fatalf("did not expect a panic")
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
		result, _, _ := calc(items, "SAVE10")

		t.Run("returns correct discount", func(t *testing.T) {
			if result.Discount != 10 {
				t.Errorf("discount = %v, want 10", result.Discount)
			}
		})
	})

	t.Run("SAVE20 discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, _, _ := calc(items, "SAVE20")

		t.Run("returns correct discount", func(t *testing.T) {
			if result.Discount != 20 {
				t.Errorf("discount = %v, want 20", result.Discount)
			}
		})
	})

	t.Run("VIP discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, _, _ := calc(items, "VIP")

		t.Run("returns correct discount", func(t *testing.T) {
			if result.Discount != 30 {
				t.Errorf("discount = %v, want 30", result.Discount)
			}
		})
	})

	t.Run("tax applied after discount", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, _, _ := calc(items, "SAVE10")
		wantTax := (100 - 10) * TaxRate

		t.Run("returns correct tax", func(t *testing.T) {
			if result.Tax != wantTax {
				t.Errorf("tax = %v, want %v", result.Tax, wantTax)
			}
		})
	})

	t.Run("empty order", func(t *testing.T) {
		_, panicked, _ := calc(nil, "")

		t.Run("panics", func(t *testing.T) {
			if !panicked {
				t.Errorf("expected a panic for an empty order")
			}
		})
	})

	t.Run("negative quantity", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: -1}}
		_, panicked, _ := calc(items, "")

		t.Run("panics", func(t *testing.T) {
			if !panicked {
				t.Errorf("expected a panic for a negative quantity")
			}
		})
	})

	t.Run("unknown discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: 1}}
		_, panicked, _ := calc(items, "NOTREAL")

		t.Run("panics", func(t *testing.T) {
			if !panicked {
				t.Errorf("expected a panic for an unknown discount code")
			}
		})
	})
}
