# Order Pricing

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* The public entry point should still be a single function that takes a list of items and a discount code and returns the computed totals.

## Running Tests

`go test ./...`

## The Challenge

`CalculateOrder` totals up a shopping cart: it sums line items, applies a discount code, then applies tax to what's left. It works, and the tests pass. But it was clearly bolted together fast — an empty cart, a negative quantity, or a discount code nobody recognizes all just crash the program outright, mid-calculation, with a bare string message and no way for a caller to tell them apart or recover.

The tax rate lives in a single package-level variable that every call to `CalculateOrder`, from anywhere in the program, reads from — there's no way to hand a different rate to a different call without stepping on every other in-flight calculation.

And the discount lookup is a chain of string comparisons that grows by one `else if` every time marketing invents a new code, with the only way to add one being to go find this function and edit it directly.

### The Breaking Point

This function is about to be called from inside a real checkout API endpoint. When a shopper submits a cart with a typo'd discount code, or a bad request slips through with a negative quantity, the endpoint needs to return an ordinary 4xx JSON error response — not take down the whole request-handling goroutine. Separately, the company is opening in a second region with a different tax rate, and both regions need to run correctly at the same time.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change (a new discount code, a new region) easy to make without reopening what you just refactored
* Callers should be able to distinguish *why* a calculation failed, programmatically, not just that it failed
* Two calls with two different tax rates should be able to run correctly without interfering with each other

## Bonus Challenge

* Make adding a new discount code a change that doesn't require touching `CalculateOrder`'s body at all

## If you get stuck

* What would it take for a caller to tell "empty cart" apart from "bad discount code" without doing string matching on a panic message?
* Where does `TaxRate` actually need to live so two calls don't share it?
* Is the discount lookup really about *comparing strings*, or is each code really a different way of computing the same thing?

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
