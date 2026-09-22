# Account Registration

**Exercise:** 1  
**Difficulty:** 1/5  
**Estimated time:** 30–45 minutes

## Scenario

The product has a small HTTP endpoint for creating accounts. It accepts an email and plan, rejects invalid or duplicate registrations, stores the user, and sends a welcome email.

The endpoint shipped quickly and works. Its tests pass. However, nearly every decision involved in registering an account now lives inside `ServeHTTP`, beside JSON parsing and HTTP responses.

## The Breaking Point

Two new callers need to register accounts next sprint:

- an internal batch import command
- a queue consumer that handles partner sign-ups

Both callers must apply the same normalization, validation, duplicate detection, persistence, and welcome-email behavior. They should not construct fake HTTP requests, and the rules must not be copied into three places.

## Your Challenge

Refactor the code so the existing endpoint continues to behave exactly as it does while the account-registration operation can be reused by a non-HTTP caller.

Keep the solution proportionate to this small service. You should be able to explain what owns each decision and why each dependency is located where it is.

Do not add the batch command or queue consumer. The goal is to leave a clean place for those callers to use later.

## Current Behavior

- Only `POST` is allowed.
- Unknown JSON fields and malformed JSON are rejected.
- Emails are trimmed and normalized to lowercase.
- An email must be non-empty and contain `@`.
- Plans are limited to `free`, `pro`, and `team`.
- Duplicate normalized emails are rejected.
- The user is stored before the welcome email is sent.
- Storage and email failures return an internal-server error.

## Acceptance Criteria

- All existing tests continue to pass.
- The core registration operation is callable without `http.Request` or `http.ResponseWriter`.
- Registration rules have one source of truth.
- HTTP concerns stay at the HTTP boundary.
- Failure information is strong enough for the HTTP boundary to preserve the existing status codes.
- The account-registration code is not tied to the in-memory implementations.
- The refactor does not introduce abstractions without a concrete use in this exercise.

## Suggested Workflow

1. Run the tests and read the production code.
2. Mark which lines understand HTTP and which lines understand account registration.
3. Decide what a non-HTTP caller would need to provide and receive.
4. Move one responsibility at a time while keeping tests green.
5. Revisit the names and dependencies after the behavior is stable.

## Setup

From this directory:

```bash
go test ./...
go test -race ./...
go vet ./...
```

## Restrictions

- Use the standard library only.
- Do not use a package-level mutable variable as a dependency.
- Do not add a framework or dependency-injection container.
- Preserve externally observable behavior unless you clearly document and test an intentional improvement.

## Hints

<details>
<summary>Hint 1</summary>

Imagine the next caller has a plain `context.Context`, an email, and a plan. Which parts of the current method would still make sense?

</details>

<details>
<summary>Hint 2</summary>

Look at the concrete types accepted by `NewRegistrationHandler`. Ask which capabilities the registration operation actually needs.

</details>

<details>
<summary>Hint 3</summary>

The HTTP layer needs to distinguish invalid input, duplicates, and internal failures. String comparison is not the only way to make errors distinguishable.

</details>

## Bonus Challenge

The user is already stored when the email provider fails. A client retry then receives a duplicate error even though the first request reported failure. Explain two production policies that could handle this situation. Implement one only if the required refactor is already clean and tested.

## Questions to Answer Afterwards

1. Which layer owns email normalization and plan validation in your refactor?
2. Where are dependency failures translated into HTTP status codes?
3. Which interfaces did you introduce, and which type consumes each one?
4. What would the queue consumer need to call your registration operation?
5. What behavior would you choose for a welcome-email failure in a real system?

## Submit Your Attempt

Push your refactor or send the diff back for review. The review will cover your approach, correctness issues, an idiomatic Go alternative where useful, the backend concept behind the exercise, and one takeaway.
