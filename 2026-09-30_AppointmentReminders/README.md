# Appointment Reminders

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Standard library only — don't reach for a third-party mocking library to make this easier to test.

## Running Tests

`go test ./...`

## The Challenge

Riverside Family Clinic's scheduling system sends a reminder email to every patient the day before their appointment. `ReminderService.SendReminder` formats the message and hands it to the clinic's SMTP relay, and it's worked fine since launch — reminders go out, patients show up more often, everyone's happy.

The trouble started when a teammate tried to add a "preview" screen to the admin dashboard, so front-desk staff can see what a reminder will say before it goes out. Wiring that up meant reaching into `SendReminder`, because there's no way to get just the formatted message without also being ready to actually hand it to a live relay connection.

Now product wants delivery analytics — which relay, how often it fails, that kind of thing — and every attempt to unit test that logic runs into the same wall: the thing actually worth testing is "did we build the right message and hand it off correctly," but there's no way to exercise that without a relay configuration sitting in the middle of every test.

### The Breaking Point

A meaningful share of Riverside's older patients don't reliably check email, but they always carry a phone. The clinic's ops lead wants: try emailing the reminder, and if the patient has no email on file, send it as a text through the SMS gateway the billing team already integrated elsewhere in the codebase. Today, adding that fallback means either duplicating `SendReminder`'s message-building logic into a near-identical method, or bolting a channel-detection branch onto one function that now has to know the details of two completely different delivery mechanisms.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change easy to make without reopening what you just refactored
* Make it possible to add the SMS fallback without duplicating the message-building logic or teaching `ReminderService` the internals of a second delivery mechanism
* Make it possible to test message construction and delivery-failure handling without any relay configuration involved at all

## Bonus Challenge

* Implement the email-then-SMS fallback described in the Breaking Point, with a test proving it falls back only when the patient has no email on file

## If you get stuck

* Look at exactly where `*SMTPClient` gets created inside `SendReminder`, and ask what would have to change for a test to hand `ReminderService` something else instead
* What's the smallest set of behaviors `ReminderService` actually needs from whatever delivers the message? It's probably narrower than the whole `SMTPClient` struct
* Consider what `NewReminderService` would need to accept as a parameter for the rest of the type to stop caring how messages actually get delivered

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
