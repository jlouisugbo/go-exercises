# Payment Status Poller

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep `PaymentGateway` as the dependency boundary — don't swap in a concrete HTTP client type.

## Running Tests

`go test ./...`

## The Challenge

After a customer submits a payment, the checkout service hands off to `StatusPoller` to find out how it resolved: does it settle, or does it fail? The payment processor doesn't push webhooks for this particular flow (an integration gap nobody's gotten around to closing yet), so `WaitForTerminalStatus` falls back to asking the processor over and over — every few seconds, up to a fixed number of times — until it gets back `"settled"` or `"failed"`, or gives up and returns `ErrPollExhausted`.

It works. Run it by hand and it waits, it polls, it returns the right answer. The trouble only shows up once you look at what's calling it: an HTTP handler on a `/payments/{id}/confirm` endpoint, built with the rest of the service's usual request-scoped plumbing — context included. Every other gateway call in this codebase takes that same first parameter and treats it as a live contract: if the request's deadline passes, or the client disconnects, the call is supposed to unwind quickly instead of running to completion regardless.

### The Breaking Point

Product wants `/payments/{id}/confirm` to have a strict 2-second budget — if a settlement decision isn't back by then, the endpoint should return a "still processing, check back later" response immediately rather than making the caller wait. The handler already builds its request context with that 2-second deadline and passes it straight into `WaitForTerminalStatus`. Nothing changes. The poller runs its full schedule regardless of how long ago the deadline passed, the handler blocks until it finally returns, and the "immediately" part of the requirement never happens no matter how the handler is wired.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change easy to make without reopening what you just refactored
* `WaitForTerminalStatus` should stop promptly and return a clear error once its context is done, instead of running its schedule to completion regardless
* Keep the terminal-status and exhausted-attempts behavior intact for callers who pass a context with no deadline

## Bonus Challenge

* Right now a canceled poll and an exhausted poll both just return an error — give the caller a way to tell them apart (e.g. with `errors.Is`) so the handler can respond differently ("still processing" vs. "gave up")

## If you get stuck

* Read `WaitForTerminalStatus` top to bottom and list every place it could, in principle, stop early. How many of them actually look at `ctx`?
* The standard library's `select` statement isn't only for channels you created yourself — it's also how a goroutine waits on more than one thing at once, including a context's own done-ness.
* Think about what should happen to the in-flight wait between polls, specifically, not just the loop as a whole.

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
