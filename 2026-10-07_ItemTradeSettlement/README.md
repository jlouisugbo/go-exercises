# Item Trade Settlement

## Setup

This is a standalone Go module with its own go.mod, independent of every other exercise in this repo. `cd` into this directory before running anything.

## Restrictions

* Any modifications to a test should maintain the spirit of the original test.

## Running Tests

`go test ./...`

## The Challenge

PixelForge Exchange lets players trade items for gold directly with each other — no auction house, just a direct offer accepted by both sides. `TradeService.ExecuteTrade` is the only place this happens: it checks that the seller actually owns the item and that the buyer can afford it, then moves the item from seller to buyer and the gold from buyer to seller.

Under the hood, that "moving" is four separate calls into the player store: take the item off the seller, give the seller the gold, take the gold off the buyer, give the buyer the item. Each call touches one player's record and can return an error if something goes wrong. Today, in practice, nothing ever does go wrong partway through — the upfront checks make sure of that. But the four calls are still four separate things happening one after another, and the code has no plan for what to do if it gets through some of them and not the rest.

That's fine until it isn't. Support wants a new rule: every player's bag has a 20-item cap, and a trade that would push someone over it should be rejected. The natural place to enforce that is right where the buyer receives the item — which is the last of the four steps.

### The Breaking Point

Add the 20-item inventory cap (enforced wherever you decide it belongs) so that a trade which would push a player over 20 items fails with a clear error. Then seed a buyer who's already at the cap and run a trade for them. Look at what's left behind for the seller and the buyer once that trade is refused at the last step.

## Goals

* Make the code easier to extend and trust long-term — the point isn't just to make today's flaw disappear, it's to leave the next change easy to make without reopening what you just refactored
* Make sure a trade either fully happens or doesn't happen at all, no matter how many steps it takes or what new checks get added to it later
* Keep the specific error cases (seller not found, buyer not found, item not owned, insufficient gold) distinguishable to a caller — don't flatten them into one generic failure

## Bonus Challenge

* Add a second trade type (e.g. gifting an item for free, or an item-for-item swap with no gold involved) without duplicating the all-or-nothing guarantee you just built
* Write a test that proves a rejected trade leaves both players' gold and inventory completely untouched, not just "mostly" untouched

## If you get stuck

* Look at what `ExecuteTrade` actually does: how many separate "write" calls are there, and what guarantee does it have that either all of them happen or none of them do?
* Think about what a single function would need to look like for a caller to be able to treat "take the item, pay out the gold" as one indivisible move rather than four independent ones.
* It might help to introduce a type whose job is to describe the whole trade before any of it is applied, kept separate from the code that applies it.

## When you're ready for a review

Don't look for a solution here — there isn't one checked in. Once your tests pass and you're happy with your refactor (or you're stuck and want a second opinion), bring your code to Claude and ask for a review of this exercise. Claude will read this README, the original code, and progress.jsonl, then walk through: your approach → bugs/issues → an idiomatic Go alternative → the backend concept behind it → one thing to remember. It'll also update progress.jsonl with what it saw.
