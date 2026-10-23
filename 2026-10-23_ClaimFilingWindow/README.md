# Claim Filing Window

## AI

There's no rule against asking Claude for help before you've tried this yourself, but consider it a friendly guideline. The point of this exercise is to build your own instincts — the growth comes from wrestling with the code yourself, not from reading someone else's refactor. Claude is genuinely useful here, it's just better used *after* your own attempt than instead of one.

So give it a real try with your own hands and judgment first. When your tests pass (or you're stuck and want a second opinion), bring it back — see "When you're ready for a review" below.

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep `EvaluateClaim`'s exported name and its `(Claim) (Decision, error)` signature intact.
* Keep the filing-window policy itself exactly as it is today: claims filed within 30 days pay out in full, 31–60 days pay 80%, 61–90 days pay 50%, and anything past 90 days is rejected.

## Running Tests

`go test ./...`

## The Challenge

`EvaluateClaim` is the one function that decides what an auto insurance claim is worth once it's filed. It takes the incident date and the claimed amount, figures out how many days passed between the incident and right now, and uses that to place the claim in one of four payout tiers. It's short, it reads in one pass, and every branch is easy to follow.

The trouble shows up when you try to pin it down in a test. There's no way to tell `EvaluateClaim` what day it is — it goes and finds out for itself every time it runs. So a test for "a claim filed 45 days after the incident" can't just write down a fixed incident date and a fixed expected result; it has to compute an incident date relative to whatever moment the test suite happens to execute. The assertion is only ever checking "as of right now," never "as of this specific, written-down date."

That's more than an inconvenience for test-writing. The same claim, with the same stored incident date, evaluated twice on two different days, can land in two different payout tiers — not because the policy changed, but because "now" moved. Nobody's noticed yet, because claims are normally evaluated once, close to when they're filed. But claims do get re-evaluated — on appeal, after a correction, during an audit — and right now there's nothing stopping that second evaluation from disagreeing with the first.

### The Breaking Point

Finance wants a nightly regression report: a fixed table of example claims, each with a literal incident date and the exact payout it's supposed to produce, diffed run over run so a reviewer can tell at a glance whether a logic change shifted anyone's payout. That only works if evaluating the same claim always produces the same answer, no matter which night the report happens to run. Today it doesn't — a claim recorded as "filed 31 days after the incident" last month will quietly drift into the 61–90-day tier if that same fixed incident date gets evaluated again today, purely because time has passed.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change (a new case, a new type, a new rule) easy to make without reopening what you just refactored
* Make `EvaluateClaim`'s result for a given claim reproducible: evaluating the exact same claim twice, on two different days, should produce the exact same decision
* Make it possible to write a test that pins an exact, fixed calendar date and asserts an exact payout, without that test's outcome depending on when the suite happens to run

## Bonus Challenge

* Give production callers a way to keep calling `EvaluateClaim` exactly as they do today, with no caller-side changes required, while still making the dependency the tests exposed explicit and overridable.

## If you get stuck

* Find the one place in `EvaluateClaim` that reaches outside its own parameters to decide anything. Every other line here is already a pure function of its inputs.
* If that one thing became something the caller handed in instead, would the test helper still need to guess what "now" will be when the suite runs?
* The question this function is really answering is "how many days between the incident and the moment we're deciding this claim." Is that moment always the instant the function happens to execute?

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
