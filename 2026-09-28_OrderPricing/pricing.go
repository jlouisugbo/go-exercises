package pricing

import "fmt"

// TaxRate applies to every order calculated anywhere in the service.
var TaxRate = 0.08

type Item struct {
	Name     string
	Price    float64
	Quantity int
}

type OrderResult struct {
	Subtotal float64
	Discount float64
	Tax      float64
	Total    float64
}

// CalculateOrder totals up an order, applying an optional discount code and
// the current tax rate.
func CalculateOrder(items []Item, discountCode string) OrderResult {
	if len(items) == 0 {
		panic("no items in order")
	}

	subtotal := 0.0
	for _, item := range items {
		if item.Quantity < 0 {
			panic("negative quantity")
		}
		subtotal += item.Price * float64(item.Quantity)
	}

	discount := 0.0
	if discountCode != "" {
		if discountCode == "SAVE10" {
			discount = subtotal * 0.10
		} else if discountCode == "SAVE20" {
			discount = subtotal * 0.20
		} else if discountCode == "VIP" {
			discount = subtotal * 0.30
		} else {
			panic(fmt.Sprintf("unknown discount code: %s", discountCode))
		}
	}

	taxable := subtotal - discount
	tax := taxable * TaxRate
	total := taxable + tax

	return OrderResult{
		Subtotal: subtotal,
		Discount: discount,
		Tax:      tax,
		Total:    total,
	}
}
