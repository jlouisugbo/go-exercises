# Invoice Delivery

**Exercise:** 4  
**Category:** Queue worker failure handling  
**Difficulty:** 3 / 10  
**Estimated time:** 30–45 minutes

## AI

Read the worker and tests before asking AI what to change. Write down which failures could improve on a later attempt and which ones would return the same answer forever. Then add the missing tests from that reasoning.

Use AI afterward as a reviewer. The important skill here is learning to decide what a failure means to the caller, not memorizing a particular helper or interface.

## Scenario

A background worker sends invoice emails through a third-party delivery gateway. Jobs are removed from the queue when `Deliver` returns `nil`; failures are returned to the queue system for later handling.

The gateway exposes HTTP-like status codes through `GatewayError`. The current worker validates the invoice, sends it, and tries again up to its configured limit whenever sending fails. This works well in tests where the gateway is briefly unavailable and then recovers.

Read `delivery.go` and follow the loop for every kind of error. Notice what information is available to the worker and which information currently affects its decision.

## Setup

This directory is a standalone Go module using only the standard library.

```bash
cd 09-28-2026_InvoiceDelivery
go test ./...
go test -race ./...
go vet ./...
```

All supplied tests pass against the starter. They must remain green after the refactor.

## Current Behavior to Preserve

- Blank invoice IDs and recipients are rejected before calling the gateway.
- Invoice fields are trimmed before delivery.
- Successful delivery stops immediately.
- Temporary gateway failures can be attempted again up to the configured limit.
- At least one attempt is made even when the worker is configured with zero.
- The final returned error keeps its underlying cause discoverable through `errors.Is`.

## The Breaking Point

A customer mistypes an email address. The gateway responds with status `422` and `ErrRecipientRejected`. The worker immediately sends the same unchanged request two more times, then the queue redelivers the job and the cycle begins again.

Another invoice is already waiting behind it. During a larger import, thousands of requests known to be rejected consume the same gateway capacity as requests that could succeed after a short outage.

The delivery team defines these rules:

- Status `408`, status `429`, and status codes from `500` through `599` may be attempted again.
- Other status codes from `400` through `499` must stop after the first response.
- An ended operation must stop without another gateway call.
- Unknown failures should be handled conservatively and must not loop forever.

## The Challenge

Refactor the worker so another attempt is made only when the failure gives a reasonable reason to expect a different result.

Add tests for the `422` incident and for an operation that has already ended. Keep error causes discoverable, retain the configured attempt limit, and keep the decision understandable when another status is added later.

The production code should not need to know how your test fake stores its result sequence. The tests should assert observable call counts and errors rather than one exact internal design.

## Acceptance Criteria

- A `422` rejection results in exactly one gateway call.
- Status `408`, status `429`, and `5xx` failures may be attempted again.
- Other `4xx` responses stop after one call.
- An ended operation does not trigger an additional gateway call.
- A successful later attempt still returns `nil` immediately.
- The final error remains compatible with `errors.Is` and, where useful, `errors.As`.
- The logic has one obvious place that decides whether another attempt is allowed.

## Restrictions

- Use the standard library only.
- Do not compare complete error-message strings.
- Do not discard the original error when wrapping it.
- Do not add `time.Sleep` to the required solution or tests.
- Do not special-case only status `422`; the stated status families should behave consistently.
- You may change method signatures, constructors, or error types when the behavior remains compatible.

## Suggested Workflow

1. Run the existing suite and trace a sequence of `503`, `503`, then success.
2. Add a test using `GatewayError{StatusCode: 422, Err: ErrRecipientRejected}`.
3. Add a test where the supplied operation has already ended before delivery starts.
4. List each possible outcome before modifying the loop.
5. Put the decision in one place and rerun the tests after each change.
6. Run vet and the race detector when the suite is green.

## If You Get Stuck

<details>
<summary>Hint 1</summary>

The loop currently asks only whether an error exists. What other facts are available before it chooses `continue`?

</details>

<details>
<summary>Hint 2</summary>

Use `errors.As` when behavior depends on structured information carried by an error. Use `errors.Is` when checking whether an error chain contains a particular cause.

</details>

<details>
<summary>Hint 3</summary>

The caller's operation can end between any two attempts, including before the first one.

</details>

<details>
<summary>Hint 4</summary>

Try expressing the decision as a small question that accepts an error and answers whether another call is justified.

</details>

## Bonus Challenge

Add waiting between allowed attempts without making tests slow or flaky. Support a gateway-provided delay for status `429`, otherwise increase the delay between attempts. The wait itself must stop promptly when the operation ends, and tests should control time without relying on real sleeps.

## Questions to Answer Afterward

1. Which failures can reasonably produce a different result on another attempt?
2. Why is repeating a rejected request harmful even when the gateway can handle the traffic?
3. Where does your code check whether the operation has ended?
4. How do `errors.Is` and `errors.As` serve different purposes in your solution?
5. What happens when the worker receives an error type it does not recognize?
6. What extra information would you need before adding delays between attempts?

## Submitting

Create a solution branch, for example:

```bash
git switch -c 09-28-2026_InvoiceDelivery
```

Push the refactor and send the branch back in ChatGPT. The review will follow:

`my approach -> bugs or issues -> idiomatic Go alternative -> backend concept -> one thing to remember`
