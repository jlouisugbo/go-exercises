# Gift Card Redemption

**Category:** Backend state changes and failure handling  
**Difficulty:** 2 / 10  
**Estimated time:** 30–45 minutes

## AI

There is no rule against using AI, but the value of this exercise comes from tracing the state changes yourself. Read the code, write down what must remain true, and form a hypothesis before asking for a proposed structure.

Use AI afterward to review your reasoning or compare designs. The goal is to build the instinct to notice dangerous failure paths during an ordinary code review.

## Setup

This exercise is a standalone Go module with no third-party dependencies or external services.

```bash
go test ./...
go test -race ./...
go vet ./...
```

All supplied tests pass against the starter. They should remain green while you refactor. You will probably want to add at least one test for the new guarantee described below.

## Restrictions

- Test changes must preserve the behavior that the original tests protect.
- Use the standard library only.
- Do not fix the problem by ignoring an error.
- Do not fix it by simply reversing the two writes. That only changes which half can be left behind.
- Keep the in-memory implementation. You do not need SQL, a database driver, or a web handler.
- You may change constructors, method signatures, and the boundary between the service and store.

## The Challenge

`RedemptionService` handles gift-card redemption for an online store. A successful redemption changes two pieces of state: the card becomes unavailable for future use, and its value is added to a user's balance.

The feature began with an in-memory store, and each store method is individually protected by a mutex. The happy path looks sensible: find the card, mark it redeemed, credit the balance, and return a receipt. The existing tests confirm each visible part of that path.

Read `redemption.go` before editing it. Follow the state after every return statement. Pay particular attention to what the caller is told, what has already changed at that moment, and what would happen if the caller tried the same operation again.

### The Breaking Point

The real balance store occasionally becomes unavailable for a few seconds during deployments. One morning it fails immediately after several gift cards have been marked as redeemed.

The API reports an internal failure, so clients retry. Every retry is rejected because the card now looks used, but the customer never received the balance. Support can see both facts in separate admin screens, yet there is no safe automated recovery because the system cannot tell a completed redemption from a partially completed one.

The product requirement is now explicit: if redemption reports failure, this attempt must not leave either piece of state changed. If it reports success, both changes must be visible together.

Refactor the code so this guarantee has one clear owner. Add or update tests to prove the failure path. Keep the design small and make it reasonable to replace the in-memory implementation with a database-backed one later.

## Current Behavior to Preserve

- Blank user IDs are rejected.
- Card codes are trimmed and normalized to uppercase.
- Unknown and previously redeemed cards are rejected with distinguishable errors.
- A successful redemption returns the amount and resulting balance.
- The redeemed card records the user who redeemed it.
- Dependency errors remain discoverable with `errors.Is`.

## Goals

- A successful call makes both related state changes exactly once.
- A failed call leaves no partial result from that attempt.
- The code has one obvious place responsible for preserving that guarantee.
- The service does not need to know the mechanical details of undoing storage operations.
- Tests demonstrate the new failure behavior without locking you into one exact internal design.
- The solution stays proportional to a small Go service.

## Suggested Workflow

1. Run the tests and trace one successful call.
2. Set `CreditErr`, trace the same call, and inspect the card afterward.
3. Write the missing failure-path test before changing production code.
4. Decide which component has enough information and control to protect the whole state change.
5. Refactor in small steps and rerun the race detector at the end.

## Bonus Challenge

Start two goroutines that attempt to redeem the same card for different users at the same time. Exactly one should succeed, the card should identify that winner, and only the winner's balance should change.

Make the test deterministic enough to run repeatedly with `go test -race -count=100 ./...`. Avoid using an arbitrary sleep as synchronization.

## If You Get Stuck

<details>
<summary>Hint 1</summary>

List the valid final states. Is "redeemed card with no credited balance" one of them?

</details>

<details>
<summary>Hint 2</summary>

A mutex around each individual method does not protect the gaps between those method calls.

</details>

<details>
<summary>Hint 3</summary>

Ask which object can see and control both the card record and the balance while the decision is being made.

</details>

<details>
<summary>Hint 4</summary>

Imagine the in-memory store is replaced by SQL. Where would beginning, committing, and rolling back one database transaction naturally belong?

</details>

## Questions to Answer Afterward

1. What invariant does your refactor protect?
2. Which code owns the full state change, and why?
3. What happens if the card is redeemed twice concurrently?
4. How would a SQL implementation preserve the same behavior?
5. Did your solution add an interface or abstraction? What concrete change does it make easier?

## Submitting

Push your refactor or paste the diff into ChatGPT. The review will follow this structure:

`your approach → bugs or issues → idiomatic Go alternative → backend concept → one thing to remember`
