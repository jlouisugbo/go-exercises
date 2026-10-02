# Shipment Status Updates

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep `UpdateStatusHandler` as the HTTP entry point — you can reshape what it calls into, but a `POST /shipments/status` request should still produce the same behavior it does today.

## Running Tests

`go test ./...`

## The Challenge

This is the whole status-update flow for a package-delivery API: a courier app or ops dashboard calls `POST /shipments/status` with a shipment ID and a new status, and `UpdateStatusHandler` is where everything happens. It decodes the JSON body, trims and validates the fields, looks the shipment up in the in-memory store, and then runs the rules for what a status change actually means — a shipment that goes `delayed` long enough past its carrier ETA gets flagged for escalation, and one that reaches `delivered` gets that flag cleared. Then it writes the updated shipment back as JSON.

It reads fine start to bottom, which is part of the problem — there's exactly one function, and it currently knows about the HTTP request/response cycle, field validation, the shipment store, and the escalation rule all at once. Testing any one piece of that means building a fake HTTP request and reading the rules' effects back out of a decoded JSON response, even though none of those rules have anything to do with HTTP.

### The Breaking Point

Ops wants a nightly job that scans every shipment still `in_transit` past its carrier ETA and runs it through the same delayed-status-and-escalation logic as a normal status update — no incoming HTTP request involved, just a cron job walking the store directly. Right now the only place that logic exists is inside a function that takes an `http.ResponseWriter` and an `*http.Request`, so the batch job's options are to reimplement the rule from scratch (and now there are two copies to keep in sync) or to fabricate a fake HTTP request just to call a handler, for a job that was never going to produce an HTTP response in the first place.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change easy to make without reopening what you just refactored
* Make it possible for something other than an HTTP request to trigger a status update and get the same escalation behavior, without duplicating the rule
* Make the escalation rule itself testable without going through JSON encoding/decoding or an `http.ResponseWriter`

## Bonus Challenge

* Write a small `RunEscalationSweep` function (or similar) that the hypothetical nightly job would call directly, and give it its own test that doesn't touch `net/http` at all

## If you get stuck

* List out everything `UpdateStatusHandler` currently has to know about to do its job. Which of those things are about HTTP, and which would still be true if the caller were a cron job instead of a request?
* What would it take to write a test for "a shipment delayed 3 hours past ETA gets escalated" that never constructs an `http.Request`?
* The shipment store, the escalation rule, and the HTTP decoding are three different reasons for this code to change. They don't have to live in the same function to stay related.

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
