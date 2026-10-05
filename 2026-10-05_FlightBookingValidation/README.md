# Flight Booking Validation

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep the three-stage flow intact (create, confirm, reschedule) -- the internal booking tool calls all three at different points in a trip's life.

## Running Tests

`go test ./...`

## The Challenge

`booking.go` handles trip bookings for a small travel agency's internal tool. A `TripBooking` moves through three stages: it's created when the customer submits the trip form, confirmed once payment clears, and occasionally rescheduled if their plans change. At each stage, a `Validate*` function runs its own checks before the transition is allowed to go through.

The three functions were added one at a time, each by whoever happened to be working on that part of the flow that week. `ValidateForCreate` checks the email and passenger count. `ValidateForConfirm` checks the dates and seat class. `ValidateForReschedule` checks the new dates. Nobody has gone back to line the three lists of checks up side by side -- and since every field on `TripBooking` is exported, anything downstream can also just reach in and read or change them directly, without going through any of the three functions at all.

### The Breaking Point

Ops wants a new rule: any booking with 6 or more passengers needs a non-empty `CorporateID`, since group travel goes through a separate billing process. Add that check. Then ask yourself: where did you just put it, and what happens to a booking that was created with 7 passengers and no `CorporateID`, then confirmed, then rescheduled -- did your change actually catch it at every one of those points, or just the one function you edited?

## Goals

* Make the code easier to extend and trust long-term -- the point isn't just to make today's flaw disappear, it's to leave the next change easy to make without reopening what you just refactored
* Make it so a new validation rule can't quietly end up applying to only some of the lifecycle stages
* Make an illegal booking (a zero passenger count, a return date before the departure date) hard to construct in the first place, not just hard to get past validation

## Bonus Challenge

* Add a rule that international itineraries (decide for yourself what marks a booking as international) require a passport number, and make sure it's enforced everywhere a booking can be created or changed -- not just at one call site.

## If you get stuck

* Write down, for each of the three `Validate*` functions, exactly which fields it checks. Where the three lists disagree, is that intentional, or just drift?
* What would happen if a fourth lifecycle stage got added next month -- would it need its own hand-written list of checks too?
* Consider what it would take to make `TripBooking` impossible to construct in an invalid state to begin with, rather than validating it after the fact.

## When you're ready for a review

Don't look for a solution here -- there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach, bugs/issues, an idiomatic Go alternative, the backend concept behind it, and one thing to remember. It'll also update progress.jsonl with what it saw.
