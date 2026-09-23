# Account Registration

**Category:** Go backend boundaries and design  
**Difficulty:** 1 / 10 (baseline; later exercises will become less explicit)  
**Estimated time:** 30–45 minutes

## AI

There is no rule against using AI, but the purpose of these sessions is to build your own backend instincts. Read the code, describe what feels difficult to change, form a hypothesis, and make the smallest refactor that improves the design. Run the tests after each meaningful step.

Bring in AI afterward for review or comparison. If you ask for the finished structure before wrestling with the problem, you skip the part that teaches you how to recognize it later in production code.

## Setup

This exercise is a standalone Go module. Navigate into this directory before running commands.

No third-party dependencies or external services are required.

## Restrictions

- Any test changes must preserve the behavior the original test protects.
- Keep the existing HTTP endpoint behavior intact.
- Do not add third-party packages, a web framework, or a dependency-injection container.
- Do not use package-level mutable state to make dependencies globally accessible.
- You may change constructors, add files, introduce types, or reorganize the package.
- Keep the required refactor small enough that another engineer could understand it in one review.

## Running Tests

```bash
go test ./...
go test -race ./...
go vet ./...
```

All tests pass against the starter code. They should continue passing after your refactor. The helper block near the top of `registration_test.go` is the intended adjustment point if your public construction API changes.

## The Challenge

`RegistrationHandler` powers the public account-registration endpoint. It parses a request, normalizes and validates the submitted data, checks for an existing account, stores a new user, sends a welcome email, and chooses the HTTP response.

The endpoint has been stable for months. Most changes have been small enough that adding another condition inside `ServeHTTP` felt reasonable, and the in-memory dependencies made the original feature quick to test.

Read `registration.go` before changing anything. Notice which lines understand HTTP, which lines understand what a valid account is, which concrete types the handler knows about, and how a failure becomes a response. Also pay attention to the order in which externally visible actions occur.

### The Breaking Point

The partnerships team is adding two new sources of registrations:

- a nightly import command for accounts received in a CSV file
- a queue consumer for partner sign-up events

Both sources must behave like the public endpoint. Emails must be normalized the same way, the same plans must be accepted, duplicates must be handled consistently, and successful registrations must still send the same welcome message.

The first proposal was to call the HTTP handler with a fake request. The second was to copy the registration block into each caller. Neither approach survives the next rule change safely. A new plan or validation rule could behave differently depending on how the account entered the system.

Refactor the starter so a future non-HTTP caller has a natural way to perform a registration without duplicating these decisions or pretending to be an HTTP client. You do not need to implement either new caller.

## Goals

- Leave one authoritative path for the decisions involved in registering an account.
- Preserve every current HTTP status and successful response covered by the tests.
- Make it possible for a future caller to register an account without constructing HTTP objects.
- Make replacing the in-memory storage or email implementation a local change.
- Keep failure information meaningful enough that each caller can translate it for its own environment.
- Improve the design without building a framework for hypothetical requirements.

## Bonus Challenge

Look at what happens when saving the user succeeds and sending the welcome email fails:

1. The endpoint reports an internal-server error.
2. The user now exists.
3. Retrying the same request returns a duplicate conflict.

Write down two policies a production system could choose for this situation. Examples may involve changing what counts as success, retrying one side effect separately, or changing how work is recorded. Implement one only after the required refactor is complete and tested.

## If You Get Stuck

- Draw a line around the code that only makes sense because the caller is HTTP. What remains on the other side?
- Imagine the queue consumer has only an email and a plan. What single operation would you want it to call?
- Look at the concrete types accepted by the current constructor. What behavior does the registration workflow actually use from each one?
- The endpoint must distinguish invalid input, duplicate data, and infrastructure failure. Consider how code can communicate categories of failure without depending on complete error-message strings.
- If your refactor creates many interfaces or layers, ask what concrete change each abstraction makes easier today.

## Questions to Answer Afterward

1. Where do normalization and plan validation live in your version, and why?
2. How does the HTTP code decide which status to return?
3. Which dependencies can now be replaced in a test or production program?
4. What would the queue consumer call?
5. Which part of your design would change if welcome emails became asynchronous?
6. Did you introduce anything that the current requirements do not justify?

## Submitting

When you are done, push your refactor or paste the diff back into ChatGPT. The review will follow this structure:

`your approach → bugs or issues → idiomatic Go alternative → backend concept → one thing to remember`
