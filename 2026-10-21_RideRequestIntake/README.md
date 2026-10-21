# Ride Request Intake

## AI

There's no rule against asking Claude for help before you've tried this yourself, but consider it a friendly guideline. The point of this exercise is to build your own instincts — the growth comes from wrestling with the code yourself, not from reading someone else's refactor. Claude is genuinely useful here, it's just better used *after* your own attempt than instead of one.

So give it a real try with your own hands and judgment first. When your tests pass (or you're stuck and want a second opinion), bring it back — see "When you're ready for a review" below.

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep `SubmitRideRequest`'s exported name and its `(*RequestStore, RideRequest) (Ride, error)` signature intact.
* Keep the `RideRequest.RequestID` field — the rider app already sets it on every request it sends, retries included, so it needs to keep meaning something to whoever reads it.

## Running Tests

`go test ./...`

## The Challenge

`SubmitRideRequest` is the single function standing between a rider tapping "Request Ride" and a car actually getting sent their way. It checks that the request looks sane, builds a `Ride`, and files it in the store that driver-matching reads from. It's short, and it reads cleanly top to bottom.

The rider app doesn't wait around forever for a response, though. If the dispatch call is slow, or the connection drops before a response makes it back, the app assumes something went wrong and sends the exact same request again — down to a `RequestID` it generated once, before the first attempt, specifically so the backend could tell "the same tap, sent twice" apart from "the rider tapped again." Right now nothing on this side ever looks at that field. Every call that reaches `SubmitRideRequest`, no matter how it got here, ends the same way: a new ride, a new ID, a new fare charge.

Most of the time a rider's connection is fine and this never comes up. But urban cell coverage is not most of the time, and support already has a running list of riders who ended up with two drivers circling their pickup spot and two charges on their card for what was, from the rider's side, one tap.

### The Breaking Point

Dispatch wants to turn on an aggressive client-side retry in the rider app: up to three attempts with backoff, triggered on any timeout or 5xx from this call, with no changes expected on this side — the assumption from the mobile team is that calling `SubmitRideRequest` more than once for the same tap is already safe. It isn't, and turning the retry on as-is would take today's occasional double-dispatch and make it routine.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change (a new case, a new type, a new rule) easy to make without reopening what you just refactored
* Make it safe for the rider app's retry to reach `SubmitRideRequest` more than once for the same tap, without dispatching or charging twice
* Leave the behavior for a genuinely new request (a `RequestID` that hasn't been seen before) exactly as it is today

## Bonus Challenge

* Decide, and make explicit in your code, what should happen if the same `RequestID` shows up again but with different ride details (a different fare or pickup zone) than the first time — is that still a retry, or something else?

## If you get stuck

* Walk through what `SubmitRideRequest` actually knows about a request it's seeing for the second time, versus what it currently checks before dispatching.
* A `RequestID` that's arriving for the second time and a `RequestID` that's never been seen before shouldn't look identical to this function.
* Think about *where* "have I already handled this exact request" needs to get answered, relative to where the ride actually gets created and saved.

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
