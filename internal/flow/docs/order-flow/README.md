# Order Flow

## Purpose

Order flow allows a user to turn a selected catalog variant into a confirmed order request.

It is a separate flow branch that starts from catalog navigation and reuses the existing flow/session/history model.

Order flow is intentionally simple.
It is designed for busy non-technical operators and customers, not for complex back-office workflows.

## Current scope

Current MVP scope:

- start order from selected catalog variant
- show order confirmation screen
- confirm order
- show done screen

Current MVP does not include:

- order persistence
- payment method selection
- promo code input
- loyalty rules
- first-order discounts
- campaign rules
- operator processing
- order status tracking

These parts are planned as follow-up steps.

## Entry point

Order flow starts only from selected catalog variant leaf.

That means:

- user navigates through catalog
- reaches product variant
- chooses order action
- flow builds order context from current catalog path

## Core principles

### Catalog remains the source of base price

Catalog variant stores base price only.

Catalog price must not be mutated by payment method, promo code, loyalty logic, or campaign rules.

### Order uses computed quote

Order flow works with a computed quote, not with mutated catalog price.

Quote is built from:

- selected variant base price
- selected payment method
- optional promo code
- optional discount rules
- optional campaign rules

### Flow owns navigation, not storage

At MVP stage order flow is only a navigation/use-case path.

It should not depend on SQL, repository details, or transport-specific storage logic.

### Keep the UX small and clear

The initial order flow should stay minimal:

- one obvious entry point
- one confirmation step
- one completion screen

No hidden logic and no long form-filling path.

## Current planned MVP path

1. user reaches catalog variant leaf
2. user selects order action
3. flow builds `OrderContext`
4. flow shows confirmation screen
5. user confirms order
6. flow shows done screen

## Planned extensions

The next order-flow extensions are expected to be added in small steps.

### Payment step

A payment-selection step may be inserted before confirmation:

1. order start
2. payment selection
3. quote recalculation
4. confirmation
5. done

### Promo and discount step

Promo and discount logic will be added on top of quote calculation.

Planned promo sources:

- first-order discount
- loyalty discount
- promo code
- campaign rule

### Persistence

After the Telegram UX becomes stable, order creation can be connected to application service and storage.

That step is expected to introduce:

- `OrderCreator`
- `CreateOrderParams`
- order service
- postgres repository
- orders migration

## Relationship with catalog flow

Order flow is not a replacement for catalog flow.

Catalog flow is responsible for:

- navigation
- path selection
- screen history
- leaf selection

Order flow is responsible for:

- collecting selected order context
- showing quote/confirmation
- confirming request

This separation must stay explicit.
