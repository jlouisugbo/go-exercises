# Worker Pool Basics

Build a small worker pool in Go.

## Goal

Implement `ProcessJobs` in `worker_pool.go`.

The function should:

- accept a list of integer jobs
- process them with exactly `workers` goroutines
- square each job value
- preserve the input order in the returned slice
- return an error when `workers <= 0`
- return an empty slice for empty input

## Why this exercise

This practices a few backend-oriented Go fundamentals:

- goroutines
- channels
- coordinating concurrent work
- avoiding goroutine leaks
- preserving deterministic output
- table-driven tests

## Run it

```bash
go test ./...
```

## Stretch goals

1. Add context cancellation with `context.Context`.
2. Make the worker function injectable instead of always squaring.
3. Add a benchmark comparing 1 worker vs multiple workers.
