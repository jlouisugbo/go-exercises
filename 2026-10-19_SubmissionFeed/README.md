# Submission Feed

## AI

There's no rule against asking Claude for help before you've tried this yourself, but consider it a friendly guideline. The point of this exercise is to build your own instincts — the growth comes from wrestling with the code yourself, not from reading someone else's refactor. Claude is genuinely useful here, it's just better used *after* your own attempt than instead of one.

So give it a real try with your own hands and judgment first. When your tests pass (or you're stuck and want a second opinion), bring it back — see "When you're ready for a review" below.

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep `SubmitAssignment`'s exported name intact.
* Keep the `Listener` type's signature (`func(Submission) error`) intact — every integration that cares about a submission gets wired in as one of these.

## Running Tests

`go test ./...`

## The Challenge

`SubmitAssignment` is the one function standing between a student clicking "turn in" and their assignment actually counting as submitted. Over a couple of semesters, more things have needed to know when a submission happens: the gradebook needs to pick it up, the student gets a confirmation email, and a couple of integrations with other campus systems have been bolted on along the way. Each one got wired in the same way — append another `Listener` to the slice this function walks through.

It works, and it reads cleanly top to bottom: save the submission, then tell everyone who's registered. The trouble shows up once you look at what "tell everyone" actually costs the student. Every listener runs one after another, in the same call, before the student's browser ever sees a response — and if any one of them comes back with a problem, the function stops right there. Whatever was registered after it never runs at all, and the student is left looking at an error for a submission that, underneath, was already saved just fine.

Nobody designed it to work this way on purpose. It's just what falls out of "when X happens, call these functions" once there's more than one function worth calling.

### The Breaking Point

The registrar's office wants to add a plagiarism-detection integration as another listener — it calls an external scanning service over the network. That service is reliable most of the time, but it's occasionally slow, and occasionally unreachable for a few minutes at a stretch. Whenever that happens today, three things go wrong at once: the student is stuck waiting on a service they've never heard of before their submission appears to go through, a slow-but-eventually-successful scan is indistinguishable from a real failure, and the gradebook-sync listener registered right after it on the list silently stops firing for every submission made during the outage — with nothing in the code pointing at why.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change (a new case, a new type, a new rule) easy to make without reopening what you just refactored
* Let a student's submission be confirmed the moment it's safely recorded, without making them wait on every listener that happens to be registered
* Make sure one listener having a bad day can't prevent the other registered listeners from running at all

## Bonus Challenge

* Give the caller a way to find out which specific listener(s) failed for a submission, without changing what a caller who doesn't care about that gets back

## If you get stuck

* Count how many things currently have to succeed, in order, before `SubmitAssignment` can return to its caller. Which of those does the caller actually need to wait on?
* A caller who submitted successfully and a caller whose plagiarism scanner happened to time out get the exact same kind of error today. Should they?
* Ask what a `Listener` failing should mean for the *other* listeners on the list, not just for the caller.

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
