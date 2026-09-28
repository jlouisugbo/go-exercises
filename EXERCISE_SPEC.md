# Go Refactor Session — Generation Spec

This is the spec a scheduled Claude session follows to generate each exercise in this repo. If you (Joel) ever want to change how sessions are generated, edit this file — the scheduled task reads it fresh every run, so no other setup needs to change.

## Context

This is a solo, async version of a "Refactor Party": normally a live group session where everyone refactors the same flawed code for 40 minutes, then discusses approaches. Here it's just you and Claude. Three mornings a week (Mon/Wed/Fri) a new exercise shows up in this repo. You work through it on your own time; when you're ready, you bring your attempt back to Claude (any Claude session — chat, Cowork, whatever's handy) for the second half: a code review.

The underlying goal is not "memorize backend architecture." It's **learn Go by wrestling with backend-shaped code**. You already know Go syntax; you're new to backend patterns. Early exercises should have fairly visible problems. Later ones should force real reasoning about ownership, dependencies, API boundaries, failure modes, concurrency, and maintainability.

## Topic rotation

Rotate through these, roughly in order, but adapt to progress.jsonl (see below) rather than following it rigidly:

1. Idiomatic Go structure & error handling
2. Interfaces & dependency injection — without overengineering
3. HTTP handlers, services, and repositories
4. Validation and domain modeling
5. Database boundaries and transaction thinking
6. Context and cancellation
7. Concurrency: goroutines, channels, worker pools
8. Middleware
9. Caching
10. Queues / pub-sub concepts
11. Idempotency and retries
12. Testing backend code effectively
13. When GoF/OOP patterns translate well to Go, and when Go has a simpler idiom
14. Configuration & dependency wiring (functional options, env-driven config)
15. Observability as a boundary concern (structured logging, metrics without leaking into domain logic)
16. Rate limiting and backpressure
17. Graceful shutdown and resource lifecycle
18. API versioning and backward compatibility
19. Generics: when they help vs. when they add noise
20. Package design: avoiding import cycles and god-packages

After topic 20, loop back to 1 — by then difficulty and combinations should keep it fresh, and you can add topics here yourself as ideas come up.

Vary the **domain** each time independently of topic (e-commerce, gaming, logistics, finance, HR/payroll, healthcare admin, IoT, travel, media, education, etc.) so exercises feel like different real systems, not one growing app.

## Progression logic (read progress.jsonl first, every time)

- Default: advance one step in the topic rotation from the last entry.
- If the last entry's `status` is `"struggled"` or `"needs_revision"`: stay on the same topic (or a closely related one) for the next session instead of advancing, and lower or hold difficulty rather than raising it. Don't advance again until a topic lands as `"completed"`.
- Difficulty is 1–5, roughly:
  - **1–2**: the flaw is visible on a read-through; one dominant, obvious problem.
  - **3**: several related things need to change together; requires understanding *why*, not just spotting *what*.
  - **4–5**: real design tension — ownership, boundaries, failure modes, or concurrency correctness are genuinely ambiguous until reasoned through.
  - Raise difficulty gradually across completed sessions on a topic; don't jump straight to 4–5 the first time a topic appears.
- This runs unattended — there's no one to ask clarifying questions of. Don't propose an idea and wait for approval (unlike a live group session); just pick the topic, domain, and difficulty per the rules above and write the files directly.

## What every exercise must be

- Solvable in about 30–40 minutes.
- Real-world shaped: code that looks like something written under deadline pressure, not a textbook example.
- One **primary** flaw tied to the session's topic, optionally 1–2 smaller independent flaws (magic strings/numbers, poor naming, etc.) for realism — never so many they compete for attention.
- Silent about its own pattern. The README describes what the code does and why it's painful. It never names the smell, the GoF pattern, or the idiom being taught. That's for the reviewer (Claude, afterward) to name, and for you to discover.
- Never shipped with a solution. The point is the attempt.

## File layout

Each exercise is a dated directory at the repo root: `{yyyy-mm-dd}_{ExerciseName}/` (PascalCase exercise name), containing:

- **go.mod** — `module github.com/jlouisugbo/go-exercises/{yyyy-mm-dd}_{ExerciseName}` and a `go 1.23` (or current) directive. No external deps unless the exercise genuinely needs one.
- **{lowercase_package}.go** — the flawed implementation. `package {lowercase_package}` at the top. Idiomatic Go otherwise (gofmt'd, sensible naming) except where the flaw itself requires the anti-pattern.
- **{lowercase_package}_test.go** — see rules below.
- **README.md** — see template below.

### Test file rules

- `package {lowercase_package}` at the top; standard library `testing` only, no third-party frameworks.
- One assertion per subtest; nested `t.Run` — outer group by method/scenario, inner per assertion.
- Cover every conditional branch (`if`/`else if`/`else`/`switch` case) with at least one subtest.
- Tests must pass against the original, unrefactored code as written (if the code currently panics on bad input, the test seam should currently assert that panic — see example below — and the comment should say the seam is expected to change once that becomes a returned error).
- **Extract a seam**: one or more small helper funcs near the top of the file, before any `Test*` func, wrapping construction and the operation under test. Label them:
  `// This helper is expected to change as you refactor — update it to match your new API.`
  (or the plural form for multiple helpers), followed by:
  `// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!`
- Every test in the file must go through these helpers — never construct the type or call the function under test directly inline elsewhere in the file.

### README template

```
# {Exercise Name}

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.
* [exercise-specific restrictions, if any]

## Running Tests

`go test ./...`

## The Challenge

[2-3 paragraphs: what the code does, why it's painful to extend or trust. Concrete, maybe a snippet. Never names the smell or pattern.]

### The Breaking Point

[A concrete new requirement that exposes the flaw and motivates refactoring — the kind of change that's trivial after a good refactor and painful before it.]

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change easy to make without reopening what you just refactored
* [exercise-specific goal]
* [exercise-specific goal]

## Bonus Challenge

* [optional stretch goal]

## If you get stuck

* [hint that guides thinking without naming the pattern]
* [hint]
* [hint]

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
```

## progress.jsonl

One JSON object per line, one line per exercise. Schema:

```
{
  "schema_version": 1,
  "date": "yyyy-mm-dd",
  "exercise": <sequential int>,
  "folder": "{yyyy-mm-dd}_{ExerciseName}",
  "title": "...",
  "topic": "<one of the 20 rotation topics>",
  "topic_index": <1-20>,
  "domain": "...",
  "difficulty": <1-5>,
  "status": "assigned" | "completed" | "needs_revision" | "struggled",
  "key_learnings": [...],
  "mistakes": [...],
  "review_changes": [...],
  "next_focus": [...]
}
```

- When generating a new exercise: append a line with `status: "assigned"` and empty arrays for the review-only fields.
- When reviewing a submitted attempt (this part is manual/on-demand, not part of the scheduled generation — see the repo README): update that exercise's line in place to `"completed"` (or `"needs_revision"`/`"struggled"` if it genuinely didn't land) and fill in the fields from the review.
- Always read the last few lines before picking the next topic/difficulty — this file **is** the adaptive-difficulty state.
