# Go Exercises

A personal workspace for repeating Go exercises. Code here is meant to be generated, practiced against, and then pushed back to this repository.

## Workflow

1. Generate or drop an exercise into the repo.
2. Solve it locally on a dated branch.
3. Run the tests, commit, and push when you are happy with the attempt.
4. Submit the branch for review.
5. Record the result in [`learning-log.jsonl`](learning-log.jsonl).

The learning log tracks completed exercises, concepts practiced, review findings, and useful next topics. See [`LEARNING_LOG.md`](LEARNING_LOG.md) for the format and update rules.

## Requirements

- [Go](https://go.dev/dl/) installed and on your `PATH`

## Quick start

```bash
git clone <this-repo>
cd go-exercises
# work through the exercise in its package/directory
go test ./...
```

Each exercise lives in its own directory. The default branch keeps the published starters, while dated solution branches preserve completed attempts and review changes.
