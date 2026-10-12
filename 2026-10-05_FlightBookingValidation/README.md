# Flight Booking Validation

## AI

There's no rule against asking Claude for help before you've tried this yourself, but consider it a friendly guideline. The point of this exercise is to build your own instincts — the growth comes from wrestling with the code yourself, not from reading someone else's refactor. Claude is genuinely useful here, it's just better used *after* your own attempt than instead of one.

So give it a real try with your own hands and judgment first. When your tests pass (or you're stuck and want a second opinion), bring it back — see "When you're ready for a review" below.

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep all three lifecycle stages -- create, confirm, reschedule -- validated by this code; the internal booking tool calls validation at all three points in a trip's life, so none of them can quietly disappear.

## Running Tests

`go test ./...`

## The Challenge

`booking.go` validates a `TripBooking` for a small travel agency's internal tool. A booking moves through three stages over its life: it's created when the customer submits the trip form, confirmed once payment clears, and occasionally rescheduled if their plans change. `ValidateBooking` takes the booking and a stage flag, and runs a different set of checks depending on which stage is passed in.

The checks for each stage were added one at a time, by whoever happened to be working on that part of the flow that week. The create stage checks the email and passenger count. The confirm stage checks the dates and seat class. The reschedule stage checks the new dates. Nobody has gone back to line the three lists of checks up side by side and ask whether they agree with each other -- and since every field on `TripBooking` is exported, anything downstream can also just reach in and read or change them directly, without going through `ValidateBooking` at any stage at all.

### The Breaking Point

Ops wants a new rule: any booking with 6 or more passengers needs a non-empty `CorporateID`, since group travel goes through a separate billing process. Add that check. Then ask yourself: where exactly did you just put it, and what happens to a booking that was created with 7 passengers and no `CorporateID`, then confirmed, then rescheduled -- did your change actually catch it at every one of those stages, or just the one you edited?

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change (a new case, a new type, a new rule) easy to make without reopening what you just refactored
* Make it so a new validation rule can't quietly end up applying to only some of the lifecycle stages
* Make an illegal booking (a zero passenger count, a return date before the departure date) hard to construct in the first place, not just hard to get past validation

## Bonus Challenge

* Add a rule that international itineraries (decide for yourself what marks a booking as international) require a passport number, and make sure it's enforced everywhere a booking can be created or changed -- not just at one call site.

## If you get stuck

* Write down, for each of the three stages, exactly which fields get checked. Where the lists disagree, is that intentional, or just drift?
* What would happen if a fourth lifecycle stage got added next month -- would it need its own hand-written list of checks too?
* Consider what it would take to make `TripBooking` impossible to construct in an invalid state to begin with, rather than validating it after the fact.

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
