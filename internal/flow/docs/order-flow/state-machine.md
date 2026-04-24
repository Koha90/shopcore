# Order Flow State Machine

## Goal

This document describes how order flow is embedded into existing `internal/flow` navigation.

Order flow must reuse the current flow model:

- transport-agnostic flow service
- session-based navigation state
- history-based back behavior
- generic catalog navigation

Order flow must not introduce transport-specific branching or SQL-dependent logic.

## MVP actions

Current MVP actions:

- `ActionOrderStart`
- `ActionOrderConfirm`

## MVP screens

Current MVP screens:

- `ScreenOrderConfirm`
- `ScreenOrderDone`

## Entry condition

Order flow can start only from selected catalog variant leaf.

If current selection is not a variant leaf, order start must not proceed.

## Order context

Order flow builds `OrderContext` from current catalog path and selected variant.

Minimal order context for MVP:

- city
- district
- product
- variant
- base price

Suggested shape:

```go
type OrderContext struct {
    CityID       int
    CityName     string

    DistrictID   int
    DistrictName string

    ProductID    int
    ProductName  string

    VariantID    int
    VariantName  string

    BasePrice    int64
}
```

This context is derived from current flow state.
It is not loaded from separate order storage.

## MVP transitions

### Start from variant leaf

Initial state:

- user already selected catalog path
- current screen is variant leaf screen
- current leaf contains selected variant

Transition:

1. user selects `ActionOrderStart`
2. flow validates that current leaf is a variant
3. flow builds `OrderContext`
4. flow renders `ScreenOrderConfirm`

### Confirm order

Transition:

1. user is on `ScreenOrderConfirm`
2. user selects `ActionOrderConfirm`
3. flow completes current order request path
4. flow renders `ScreenOrderDone`

### Done screen

`ScreenOrderDone` is the terminal screen for MVP order branch.

Expected UX:

- short confirmation text
- no extra branching
- clear return path to root or catalog entry

## Back behavior

Back behavior must continue to use session history.

Important rules:

- do not hardcode back transition to catalog root
- do not bypass `Session.History`
- do not reintroduce old hardcoded back mechanics

Expected result:

- back from order confirmation returns user to previous catalog state
- back behavior remains consistent with the rest of flow package

## Rendering responsibilities

Suggested separation:

- `handleOrderAction(...)` handles order-related actions
- `renderOrderScreen(...)` renders order-related screens
- `orderContext(...)` extracts selected order data from current catalog path/leaf

This keeps order-flow logic isolated without mixing it into generic catalog rendering.

## Suggested MVP screen behavior

### `ScreenOrderConfirm`

This screen should show:

- city
- district
- product
- variant
- base price

Actions:

- confirm order
- back

No text input is needed for MVP.

### `ScreenOrderDone`

This screen should show:

- order request accepted
- next expected user-facing outcome
- action to return to root/menu

The screen should stay intentionally simple.

## Session model impact

MVP order flow should avoid new session complexity where possible.

At this step it is enough to rely on:

- existing current screen state
- current catalog path
- existing history

If later steps require payment selection or promo input, session may grow with explicit order draft state.

## Future extension points

The following steps are expected to fit on top of this MVP path.

### Payment selection

Possible future path:

1. `ActionOrderStart`
2. `ScreenPaymentSelect`
3. `ActionPaymentSelect`
4. `ScreenOrderConfirm`
5. `ActionOrderConfirm`
6. `ScreenOrderDone`

### Promo input

Possible future path:

1. payment selected
2. optional promo input screen shown
3. promo applied to quote
4. confirmation screen re-rendered with updated quote

### Order creation service

At persistence stage, `ActionOrderConfirm` may call application service such as:

```go
type OrderCreator interface {
    Create(ctx context.Context, params CreateOrderParams) error
}
```

That wiring must stay outside generic flow navigation logic.

## Non-goals for MVP

The following should not be added in the first order-flow step:

- storage calls inside generic catalog rendering
- direct SQL access from flow
- transport-specific order logic in Telegram runtime
- complex multi-step forms
- operator/admin processing screens
