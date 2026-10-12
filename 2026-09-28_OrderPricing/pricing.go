package pricing

import (
	"errors"
	"fmt"
)

// TaxRate is the tax rate used by the default calculator (see DefaultCalculator).
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

var (
	ErrEmptyOrder       = errors.New("no items in order")
	ErrNegativeQuantity = errors.New("negative quantity")
)

// UnknownDiscountError is returned when discountCode is non-empty but not registered.
type UnknownDiscountError struct {
	Code string
}

func (e UnknownDiscountError) Error() string {
	return fmt.Sprintf("unknown discount code: %s", e.Code)
}

// Calculator holds tax rate and discount codes for one pricing context (e.g. a region).
type Calculator struct {
	taxRate   float64
	discounts map[string]float64
}

// DefaultCalculator is used by package-level CalculateOrder.
var DefaultCalculator = NewCalculator(TaxRate)

func defaultDiscounts() map[string]float64 {
	return map[string]float64{
		"SAVE10": 0.10,
		"SAVE20": 0.20,
		"VIP":    0.30,
	}
}

// NewCalculator builds a calculator with the given tax rate and standard discount codes.
func NewCalculator(taxRate float64) *Calculator {
	return &Calculator{
		taxRate:   taxRate,
		discounts: defaultDiscounts(),
	}
}

// RegisterDiscount adds or replaces a percentage discount (0–1) for a code without changing CalculateOrder.
func (c *Calculator) RegisterDiscount(code string, rate float64) {
	c.discounts[code] = rate
}

// CalculateOrder totals up an order using this calculator's tax rate and discount table.
func (c *Calculator) CalculateOrder(items []Item, discountCode string) (OrderResult, error) {
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
		rate, ok := c.discounts[discountCode]
		if !ok {
			return OrderResult{}, UnknownDiscountError{Code: discountCode}
		}
		discount = subtotal * rate
	}

	taxable := subtotal - discount
	tax := taxable * c.taxRate
	total := taxable + tax

	return OrderResult{
		Subtotal: subtotal,
		Discount: discount,
		Tax:      tax,
		Total:    total,
	}, nil
}

// CalculateOrder totals up an order using DefaultCalculator.
func CalculateOrder(items []Item, discountCode string) (OrderResult, error) {
	return DefaultCalculator.CalculateOrder(items, discountCode)
}
