package pricing

import (
	"errors"
	"testing"
)

const testTaxRate = 0.08

func calc(items []Item, code string) (OrderResult, error) {
	c := Calculator{TaxRate: testTaxRate}
	return c.CalculateOrder(items, code)
}

func TestCalculateOrder(t *testing.T) {
	t.Run("subtotal only, no discount", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: 2}}
		result, err := calc(items, "")
		if err != nil {
			t.Fatalf("did not expect an error, got %v", err)
		}
		if result.Subtotal != 20 {
			t.Errorf("expected subtotal 20, got %v", result.Subtotal)
		}
	})

	t.Run("SAVE10 discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, err := calc(items, "SAVE10")
		if err != nil {
			t.Fatal(err)
		}
		if result.Discount != 10 {
			t.Errorf("expected discount 10, got %v", result.Discount)
		}
	})

	t.Run("SAVE20 discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, err := calc(items, "SAVE20")
		if err != nil {
			t.Fatal(err)
		}
		if result.Discount != 20 {
			t.Errorf("expected discount 20, got %v", result.Discount)
		}
	})

	t.Run("VIP discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, err := calc(items, "VIP")
		if err != nil {
			t.Fatal(err)
		}
		if result.Discount != 30 {
			t.Errorf("expected discount 30, got %v", result.Discount)
		}
	})

	t.Run("tax applied after discount", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		result, err := calc(items, "SAVE10")
		if err != nil {
			t.Fatal(err)
		}
		wantTax := (100 - 10) * testTaxRate
		if result.Tax != wantTax {
			t.Errorf("expected tax %v, got %v", wantTax, result.Tax)
		}
	})

	t.Run("empty order", func(t *testing.T) {
		_, err := calc(nil, "")
		if !errors.Is(err, ErrEmptyOrder) {
			t.Errorf("expected ErrEmptyOrder, got %v", err)
		}
	})

	t.Run("negative quantity", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: -1}}
		_, err := calc(items, "")
		if !errors.Is(err, ErrNegativeQuantity) {
			t.Errorf("expected ErrNegativeQuantity, got %v", err)
		}
	})

	t.Run("unknown discount code", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 10, Quantity: 1}}
		_, err := calc(items, "NOTREAL")
		if !errors.Is(err, ErrUnknownDiscount) {
			t.Errorf("expected ErrUnknownDiscount, got %v", err)
		}
	})

	t.Run("two tax rates do not interfere", func(t *testing.T) {
		items := []Item{{Name: "Widget", Price: 100, Quantity: 1}}
		us := Calculator{TaxRate: 0.08}
		eu := Calculator{TaxRate: 0.20}

		usResult, err := us.CalculateOrder(items, "")
		if err != nil {
			t.Fatal(err)
		}
		euResult, err := eu.CalculateOrder(items, "")
		if err != nil {
			t.Fatal(err)
		}
		if usResult.Tax != 8 || euResult.Tax != 20 {
			t.Errorf("expected taxes 8 and 20, got %v and %v", usResult.Tax, euResult.Tax)
		}
	})
}
