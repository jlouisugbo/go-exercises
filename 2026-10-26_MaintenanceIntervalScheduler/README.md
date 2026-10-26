# Maintenance Interval Scheduler

## AI

There's no rule against asking Claude for help before you've tried this yourself, but consider it a friendly guideline. The point of this exercise is to build your own instincts — the growth comes from wrestling with the code yourself, not from reading someone else's refactor. Claude is genuinely useful here, it's just better used *after* your own attempt than instead of one.

So give it a real try with your own hands and judgment first. When your tests pass (or you're stuck and want a second opinion), bring it back — see "When you're ready for a review" below.

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep `ScheduleNextService`'s exported name and its `(string, time.Time) (time.Time, error)` signature intact.
* Keep the equipment-specific numbers themselves exactly as they are today: a 90-day baseline, tractors at -10 days, harvesters at -20 days, and irrigation pumps at +15 days.

## Running Tests

`go test ./...`

## The Challenge

`ScheduleNextService` figures out when a piece of farm equipment is next due for maintenance. It starts from a 90-day baseline and nudges that interval up or down depending on how hard the equipment type actually runs — tractors and harvesters get serviced sooner because of heavy seasonal use, irrigation pumps can go a little longer since they run light and steady.

The way that nudge gets picked is where it gets heavier than it needs to be. Each equipment type has its own tiny type, and that type's only job is to implement a single method that returns the adjusted number. A lookup function takes the equipment's type string, matches it against a short list of cases, and hands back whichever one of those tiny types fits — and `ScheduleNextService` just calls that method without looking at what it got back. Three equipment types means three small types, three one-line methods, and one switch statement tying them all together, to do what amounts to picking one of three numbers.

The fleet only has three equipment types today, all wired in by hand, so nobody's hit the gap yet. But the lookup function doesn't actually promise to find a match — it just falls through to doing nothing useful if the type string doesn't ring a bell, and `ScheduleNextService` never checks for that case before using what it got.

### The Breaking Point

The co-op is adding drones for crop monitoring to the fleet, and the equipment-onboarding form now lets a technician type in any equipment type string directly, free text, no longer limited to the three the system already knows about. The moment someone enters "drone" — or just fat-fingers "Tractor" with a capital T — `ScheduleNextService` doesn't reject it and doesn't fall back to the baseline. It crashes the request outright, taking down whatever called it, because nothing between the typo and the number it was about to use ever stopped to check whether a match was actually found.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change (a new case, a new type, a new rule) easy to make without reopening what you just refactored
* Make an unrecognized equipment type a normal, handled outcome — reported back to the caller as an error, never a crash
* Make adding a fourth equipment type (like the incoming drones) a change in exactly one place, not a change spread across a new type, a new method, and a switch statement that all have to agree with each other

## Bonus Challenge

* Add a way to list every equipment type the scheduler currently knows about and the adjustment each one gets, so it could back a simple admin page without anyone having to go read the source to find out what's supported

## If you get stuck

* Count the number of places you'd have to touch to add the drone case. Does each of those places really need to exist, or could the same information live somewhere more boring?
* What is each of those tiny per-equipment types actually storing or doing, beyond returning one fixed number? Is the type itself carrying any weight, or is it just a name wrapped around a value?
* Once you can't rely on the lookup always finding a match, where's the one place that result should be checked, and what should happen there instead of plowing ahead?

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
