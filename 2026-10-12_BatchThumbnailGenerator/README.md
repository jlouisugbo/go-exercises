# Batch Thumbnail Generator

## AI

There's no rule against asking Claude for help before you've tried this yourself, but consider it a friendly guideline. The point of this exercise is to build your own instincts — the growth comes from wrestling with the code yourself, not from reading someone else's refactor. Claude is genuinely useful here, it's just better used *after* your own attempt than instead of one.

So give it a real try with your own hands and judgment first. When your tests pass (or you're stuck and want a second opinion), bring it back — see "When you're ready for a review" below.

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep `GenerateThumbnails`'s exported name and signature (`photos []Photo, resizer Resizer`) intact. The review is about how the work inside it gets dispatched, not about the shape callers see.

## Running Tests

`go test ./...`

## The Challenge

`GenerateThumbnails` takes a batch of freshly uploaded photos and asks the resizing service for a thumbnail of each one. It does the obvious thing: loop over the photos, kick off a goroutine for each, and wait for all of them to report back before returning the combined results.

For the handful of photos a single user uploads at a time, this is invisible — ten goroutines come and go in a blink, and the resizing service never notices. It also reads as careful code in isolation: there's a `WaitGroup` so nothing returns early, there's a mutex guarding the shared results slice so there's no data race, and every photo gets accounted for. Nothing about the function itself looks wrong.

The resizing service is the part nobody designed around. It's a shared dependency with a hard ceiling on how many requests it can serve at the same instant, and anything past that ceiling comes back rejected. Today's traffic pattern — small batches, a few at a time — keeps every call comfortably under that ceiling, so the rejections never show up in practice. The code has never had to behave any other way.

### The Breaking Point

Product wants a "restore my library" import: a user connects an old backup and the app submits every photo in it as one batch to `GenerateThumbnails` — sometimes several hundred at once. The resizing service's real-world capacity hasn't changed. Suddenly most of a big batch comes back with errors, the import looks badly broken, and nobody can explain why a ten-photo batch never showed this problem.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change (a new case, a new type, a new rule) easy to make without reopening what you just refactored
* Make sure a batch of any size respects the resizing service's real capacity, instead of quietly depending on every caller happening to send small batches
* Keep the work concurrent — the goal isn't to go back to resizing photos one at a time, just to bound how many happen at once

## Bonus Challenge

* Make the concurrency limit a parameter instead of a hardcoded constant, so a caller talking to a bigger (or smaller) resizing deployment doesn't need a code change to match it

## If you get stuck

* Picture the batch at 500 photos. At the instant the function starts, how many goroutines exist, and how many of them are allowed to be mid-`Resize` call at once right now?
* Look at what stands between "a photo is ready to be resized" and "a goroutine for it starts running." Right now, is anything?
* A channel can carry more than values — it can carry permission to proceed.

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
