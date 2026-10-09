package pricing

import "errors"

var (
	ErrEmptyOrder       = errors.New("no items in order")
	ErrNegativeQuantity = errors.New("negative quantity")
	ErrUnknownDiscount  = errors.New("unknown discount code")
)

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

// Calculator prices one region's orders. TaxRate stays on this value so two
// regions can calculate at the same time without sharing a package variable.
type Calculator struct {
	TaxRate float64
}

type discountFunc func(subtotal float64) float64

func percentOff(rate float64) discountFunc {
	return func(subtotal float64) float64 {
		return subtotal * rate
	}
}

// New codes are added here. CalculateOrder only looks the code up.
var discounts = map[string]discountFunc{
	"SAVE10": percentOff(0.10),
	"SAVE20": percentOff(0.20),
	"VIP":    percentOff(0.30),
}

// CalculateOrder totals up an order, applying an optional discount code and
// this calculator's tax rate.
func (c Calculator) CalculateOrder(items []Item, discountCode string) (OrderResult, error) {
	if len(items) == 0 {
		return OrderResult{}, ErrEmptyOrder
	}

	subtotal := 0.0
	for _, item := range items {
		if item.Quantity < 0 {
			return OrderResult{}, ErrNegativeQuantity
		}
		subtotal += item.Price * float64(item.Quantity)
	}

	discount := 0.0
	if discountCode != "" {
		apply, ok := discounts[discountCode]
		if !ok {
			return OrderResult{}, ErrUnknownDiscount
		}
		discount = apply(subtotal)
	}

	taxable := subtotal - discount
	tax := taxable * c.TaxRate
	total := taxable + tax

	return OrderResult{
		Subtotal: subtotal,
		Discount: discount,
		Tax:      tax,
		Total:    total,
	}, nil
}
