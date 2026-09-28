# go-exercises

Solo, async "Refactor Party" for learning Go backend patterns — three mornings a week (Mon/Wed/Fri, ~8am ET) a scheduled Claude session drops a new flawed-but-realistic Go exercise here. You work through it on your own; when you're ready, you bring your attempt to Claude for a review.

## How it works

1. A new dated folder shows up automatically: `{yyyy-mm-dd}_{ExerciseName}/`.
2. Read its README — it describes the problem, never the underlying pattern. That's yours to find.
3. `cd` into the folder and `go test ./...` — tests currently pass against the flawed code. Refactor until you're happy, keeping tests green (update the seam helper at the top of the test file if your refactor changes the public API — the rest of the tests shouldn't need to change).
4. There's no solution checked in anywhere in this repo, on purpose.
5. When you're ready (tests passing, or genuinely stuck), open a chat with Claude, paste or describe your solution, and ask for a review of that exercise. Claude will read the exercise's README, the original code, and `progress.jsonl`, then walk through: your approach → bugs/issues found → an idiomatic Go alternative → the backend concept behind it → one thing to remember — and update `progress.jsonl` with the outcome.

## Files

- **EXERCISE_SPEC.md** — the full generation spec (topic rotation, difficulty progression, file/README format). Edit this any time you want to change how exercises get generated; the scheduled task reads it fresh every run.
- **progress.jsonl** — one line per exercise: topic, difficulty, status, and what the review found. This is what makes the difficulty/topic selection adaptive instead of a fixed schedule.
- **`{date}_{ExerciseName}/`** — each individual exercise, a standalone Go module.

## Requirements

- [Go](https://go.dev/dl/) installed and on your `PATH`
