# Device Reading Cache

## AI

There's no rule against asking Claude for help before you've tried this yourself, but consider it a friendly guideline. The point of this exercise is to build your own instincts — the growth comes from wrestling with the code yourself, not from reading someone else's refactor. Claude is genuinely useful here, it's just better used *after* your own attempt than instead of one.

So give it a real try with your own hands and judgment first. When your tests pass (or you're stuck and want a second opinion), bring it back — see "When you're ready for a review" below.

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep `GetReading(deviceID string) (Reading, error)`'s exported name and signature intact — the dashboard code that calls it is out of scope for this exercise.
* Keep the `SensorStore` interface's method signature intact.

## Running Tests

`go test ./...`

## The Challenge

`ReadingCache` sits between the building dashboard and the actual sensor hardware driver. Every thermostat on every floor gets polled through here, and the sensor bus can only handle so many requests before it falls over — so a cache went in front of `SensorStore.FetchReading` a while back, and it did its job. The bus stopped falling over.

Nobody's touched this file since, and lately a few things don't add up. A technician swaps a miscalibrated sensor on the third floor, confirms the hardware is reporting the right number through the manufacturer's test rig, and the building dashboard still shows the old reading — for everyone, indefinitely. Restarting the dashboard process "fixes" it, which is the kind of fix that makes people nervous rather than relieved.

`GetReading` is short: a lookup, and a fetch-and-store on miss. There's no other code around it worth reading — the whole story is in those dozen lines.

### The Breaking Point

Facilities is rolling out automated HVAC overrides that react to readings pulled through this same path. If a recalibrated or swapped sensor's new value never makes it back to a caller, an override could act on a number that stopped being true hours ago. Before that ships, they want a guarantee, in writing, for how out-of-date a reading returned by `GetReading` is allowed to be.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change (a new case, a new type, a new rule) easy to make without reopening what you just refactored
* Give every cached reading a bounded, explicit lifetime instead of an implicit "forever"
* Make that lifetime a visible, adjustable part of the API — not a constant buried in the middle of a method

## Bonus Challenge

* Support forcing a single device's entry to refresh early (e.g. right after a technician recalibrates it) without discarding every other device's cached data

## If you get stuck

* Once a value lands in that map, trace every place it could change after that. Count them.
* A cache that never forgets anything isn't really tracking a moving target — it's memoizing a constant. Which one is a `Reading` actually supposed to be?
* If you add an expiration, think about how you'd test "this went stale and got refreshed" without making your test suite sleep for real wall-clock seconds.

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
