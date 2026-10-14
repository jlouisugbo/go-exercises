# Salary Adjustment Gate

## AI

There's no rule against asking Claude for help before you've tried this yourself, but consider it a friendly guideline. The point of this exercise is to build your own instincts — the growth comes from wrestling with the code yourself, not from reading someone else's refactor. Claude is genuinely useful here, it's just better used *after* your own attempt than instead of one.

So give it a real try with your own hands and judgment first. When your tests pass (or you're stuck and want a second opinion), bring it back — see "When you're ready for a review" below.

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* Keep `RequireRole` and `WithRequestLog` as separate, independently reusable `net/http` middleware — don't collapse the role check into `AdjustSalaryHandler` itself or merge the two middleware into one function. Other endpoints in this service are going to want to wrap themselves in `RequireRole` on their own.
* Keep `NewPayrollHandler`'s signature (`func NewPayrollHandler() http.Handler`) intact.

## Running Tests

`go test ./...`

## The Challenge

The internal tools team ships a small HTTP endpoint finance ops uses to push one-off salary corrections — a signing bonus here, a backdated raise there. Nobody wants this endpoint reachable by just anybody with a valid session, so it's wrapped in two layers: `WithRequestLog`, which writes an access-log line for every hit, and `RequireRole`, which is supposed to turn away anyone whose `X-Employee-Role` header isn't `"admin"` before the request gets anywhere near `AdjustSalaryHandler`.

`RequireRole` reads like a normal auth gate: check the header, and if it doesn't match, log who got turned away and write a 403. `NewPayrollHandler` stacks it in front of the handler the way every middleware tutorial says to:

```go
func NewPayrollHandler() http.Handler {
	s := newStore()
	return WithRequestLog(RequireRole("admin", AdjustSalaryHandler(s)))
}
```

Calling it by hand with the right role works fine, and calling it with a bogus role does get you back a 403 — so at a glance, it looks finished.

### The Breaking Point

Security runs a routine pass over the payroll audit log and flags something odd: several `rejected request from role "viewer", want "admin"` lines are immediately followed by an `adjusted salary for e-100 to ...` line for the exact same request. The adjustment is going through even on requests the gate just said no to. They want this closed before anyone notices it can be triggered on purpose — and they also want a second gated endpoint, for viewing (not editing) salary history, built on the same `RequireRole`. Whatever fix lands needs to hold up for more than just this one handler.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change (a new case, a new type, a new rule) easy to make without reopening what you just refactored
* Make sure a rejected request can never reach the handler it was rejected in front of — for this gate and for any future one built the same way
* Keep the request-logging and role-checking concerns independent of each other and of the salary-adjustment logic itself

## Bonus Challenge

* Add a second small middleware of your own (e.g. requiring a particular HTTP method) and write a test proving that stacking it doesn't reopen the same hole

## If you get stuck

* Trace exactly what happens to the `http.ResponseWriter` and the request after the 403 gets written inside `RequireRole` — is anything stopping execution from continuing past that point?
* Look at how `net/http`'s own documentation and examples write a handler that needs to stop early
* Try writing a test that calls `RequireRole` directly with a mismatched role and checks more than just the status code it returns

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
