# Customer Snapshot

**Exercise:** 3  
**Category:** Request lifetime and slow dependencies  
**Difficulty:** 3 / 10  
**Estimated time:** 30–45 minutes

## AI

Use this session to practice tracing control and ownership through unfamiliar backend code. Read the implementation, draw the call path, and write the missing behavioral test before asking AI for a proposed refactor.

AI is useful afterward for reviewing whether your solution can leak work or mishandle an error. If it gives you the finished structure before you attempt the test, you lose the part of the exercise that builds the production instinct.

## Scenario

The customer-support dashboard calls `GET /snapshot?customer_id=...` before opening a case. The endpoint combines a customer's profile with recent orders from an internal directory service.

The endpoint is small, and its ordinary behavior is well tested. It rejects malformed requests, distinguishes missing customers from provider failures, and avoids the second lookup when the first one fails.

The directory client already accepts a value that can represent the lifetime of an operation. Trace what the handler receives, what the service accepts, and what the service actually gives each directory call.

## Setup

This directory is a standalone Go module using only the standard library.

```bash
cd 09-25-2026_CustomerSnapshot
go test ./...
go test -race ./...
go vet ./...
```

All supplied tests pass against the starter. They must continue passing after the refactor.

## Current Behavior to Preserve

- Only `GET` is accepted.
- A blank customer ID returns `400 Bad Request` without calling the directory.
- An unknown customer returns `404 Not Found`.
- A directory failure returns `502 Bad Gateway`.
- A successful response contains the customer and recent orders as JSON.
- Recent orders are not requested when the customer lookup fails.
- Wrapped dependency errors remain discoverable with `errors.Is`.

## The Breaking Point

The customer directory sometimes takes 20 seconds to respond during a deployment. Support agents close the case panel after a second and open another customer, which disconnects the original request.

The server keeps every abandoned lookup alive until the directory's own timeout fires. During a busy incident, hundreds of handlers and downstream calls accumulate even though nobody is waiting for their responses. Restarting the API temporarily clears the pressure, but it returns as soon as agents resume work.

The new requirement is behavioral: when the caller says the operation is over, the in-flight directory work must be able to stop promptly. No later lookup should begin for that request, and the reason the operation ended must remain recognizable to upstream code.

## The Challenge

Refactor the request path so the operation's lifetime is preserved from its entry point to every downstream call.

Add deterministic tests that prove the behavior at the service boundary or through the handler. A useful test will arrange for a fake directory call to block, observe that the call has started, end the request, and verify that the operation returns without manually releasing the fake.

Keep the design proportional to this service. The goal is a clear ownership path, not a new framework.

## Acceptance Criteria

- Existing HTTP behavior remains unchanged for normal requests.
- Ending a caller's operation can unblock an in-flight directory call.
- The resulting error remains recognizable with `errors.Is`.
- Recent orders are not requested after the operation has already ended.
- Non-HTTP callers can use the same service behavior with their own operation lifetime.
- Cancellation tests coordinate with channels or equivalent signals rather than depending on arbitrary sleeps.

## Restrictions

- Use the standard library only.
- Do not solve the problem with a package-level flag or mutable global state.
- Do not add a fixed timeout as the only fix. Different callers own different lifetimes.
- Do not remove the existing lifetime parameter from the `Directory` interface.
- Preserve the meaning of the existing errors and status codes.
- Test safety timeouts are allowed as failure guards, but a timeout must not be the primary synchronization mechanism.

## Suggested Workflow

1. Run the existing tests and trace one successful request from `ServeHTTP` to both directory calls.
2. Compare the lifetime owned by the incoming caller with the values used inside `Build`.
3. Write a fake that signals when a directory call starts and can wait for either completion or the end of the operation.
4. Write the missing test before modifying production code.
5. Make the smallest signature and call-site changes that satisfy the test.
6. Run the full suite, race detector, and vet.

## If You Get Stuck

<details>
<summary>Hint 1</summary>

The incoming HTTP request already carries information about whether its caller is still waiting.

</details>

<details>
<summary>Hint 2</summary>

Compare the parameters accepted by `Directory` with the parameters accepted by `Service.Build`. What information disappears between them?

</details>

<details>
<summary>Hint 3</summary>

Creating a fresh root value inside the service disconnects downstream work from the caller that owns it.

</details>

<details>
<summary>Hint 4</summary>

A deterministic test can use one channel to announce that work started and another signal already provided by the operation. A short timer should only prevent a broken test from hanging forever.

</details>

## Bonus Challenge

Give the directory lookup a server-side maximum duration while still respecting an earlier caller cancellation. Make sure every resource created for that limit is released, and test both which deadline wins and whether the original error remains discoverable.

## Questions to Answer Afterward

1. Who owns the lifetime of this operation?
2. At which boundary was that ownership previously lost?
3. Why is creating a fresh root value inside the service dangerous here?
4. Why should the lifetime value usually be a method parameter instead of a field on `Service`?
5. How did your test prove the call had started before ending the operation?
6. What should the HTTP layer do differently, if anything, when the caller has already disconnected?

## Submitting

Create a solution branch, for example:

```bash
git switch -c 09-25-2026_CustomerSnapshot
```

Push the refactor and send the branch back in ChatGPT. The review will follow:

`my approach -> bugs or issues -> idiomatic Go alternative -> backend concept -> one thing to remember`
