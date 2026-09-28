package pricing

import "testing"

// ---
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
// ---

func TestCalculateOrder(t *testing.T) {
	t.Run("subtotal only, no discount", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: 2}}
		result, panicked, _ := calc(items, "")
		if panicked {
			t.Fatalf("did not expect a panic")
		}
		if result.Subtotal != 20 {
			t.Errorf("expected subtotal 20, got %v", result.Subtotal)
		}
	})

	t.Run("SAVE10 discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, _, _ := calc(items, "SAVE10")
		if result.Discount != 10 {
			t.Errorf("expected discount 10, got %v", result.Discount)
		}
	})

	t.Run("SAVE20 discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, _, _ := calc(items, "SAVE20")
		if result.Discount != 20 {
			t.Errorf("expected discount 20, got %v", result.Discount)
		}
	})

	t.Run("VIP discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, _, _ := calc(items, "VIP")
		if result.Discount != 30 {
			t.Errorf("expected discount 30, got %v", result.Discount)
		}
	})

	t.Run("tax applied after discount", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, _, _ := calc(items, "SAVE10")
		wantTax := (100 - 10) * TaxRate
		if result.Tax != wantTax {
			t.Errorf("expected tax %v, got %v", wantTax, result.Tax)
		}
	})

	t.Run("empty order panics", func(t *testing.T) {
		_, panicked, _ := calc(nil, "")
		if !panicked {
			t.Errorf("expected a panic for an empty order")
		}
	})

	t.Run("negative quantity panics", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: -1}}
		_, panicked, _ := calc(items, "")
		if !panicked {
			t.Errorf("expected a panic for a negative quantity")
		}
	})

	t.Run("unknown discount code panics", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: 1}}
		_, panicked, _ := calc(items, "NOTREAL")
		if !panicked {
			t.Errorf("expected a panic for an unknown discount code")
		}
	})
}
