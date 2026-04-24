# Order Pricing

## Goal

This document describes the pricing model used by order flow.

The goal is to keep catalog pricing stable and move order-specific calculations into explicit quote logic.

## Core rule

Catalog variant stores only base price.

Order flow never mutates catalog price.
Instead, it computes an order quote.

## Pricing model

Order price is built from three independent parts:

- base price
- promo discount
- payment markup

### Formula

```text
final price = base price - promo discount + payment markup
```

## Important pricing rule

Promo discount is calculated only from base price.

Promo discount must not be calculated from:

- final price
- base price plus payment markup
- any already adjusted price

Payment markup is not discounted by promo logic.

This rule keeps pricing explainable and prevents promo rules from silently reducing payment-specific surcharge.

## Percentage rules

If promo discount is percentage-based and payment markup is percentage-based, both percentages are calculated from the same base price.

Example:

- base price: `1000`
- promo discount: `10%`
- payment markup: `5%`

Calculation:

- promo discount = `100`
- payment markup = `50`
- final price = `1000 - 100 + 50 = 950`

Not allowed:

- `(1000 + 50) - 10%`
- applying promo discount to payment markup
- chaining discount from already adjusted subtotal unless explicitly redesigned later

## Why this rule exists

This project needs clear and predictable pricing for simple everyday usage.

The user, operator, and admin should be able to explain the result in one sentence:

- item has base price
- promo reduces item price
- payment method adds its own markup
- final price is the sum of those parts

This is easier to understand, test, and support.

## Quote inputs

Suggested quote input shape:

```go
type QuoteInput struct {
    UserID         int64
    BotID          int

    CityID         int
    DistrictID     int

    ProductID      int
    VariantID      int

    BasePrice      int64
    PaymentMethod  PaymentMethod
    PromoCode      string
}
```

## Quote output

Suggested quote shape:

```go
type PriceQuote struct {
    BasePrice      int64
    PromoDiscount  int64
    PaymentMarkup  int64
    FinalPrice     int64

    PaymentMethod  PaymentMethod
    PromoCode      string

    AppliedRules   []AppliedRule
}
```

Suggested rule shape:

```go
type AppliedRule struct {
    Code        string
    Kind        string
    Description string
    Amount      int64
}
```

## Planned promo sources

The pricing model should support the following promo sources:

- first-order discount
- loyalty discount
- promo code
- campaign rule

All of them must follow the same base principle:

- discount is computed from base price
- discount does not reduce payment markup

## Rule ordering

Recommended rule order:

1. take variant base price
2. calculate promo discount from base price
3. calculate payment markup from base price
4. compute final price

This keeps the model stable and avoids ambiguous stacking.

## Stacking policy

For the first implementation, use only one active promo rule at a time.

That means only one of the following should be applied to a single quote:

- first-order discount
- loyalty discount
- promo code
- campaign rule

This keeps the model simple and predictable.

If later business requirements need stacking, that should be introduced explicitly with separate documentation and tests.

## Constraints

Pricing logic should enforce at least the following constraints:

- final price must not be negative
- unknown payment method must not silently apply random markup
- invalid promo code must not silently mutate price
- quote must be reproducible from stored pricing snapshot

## Persistence note

Confirmed order should store pricing snapshot, not only final price.

Suggested fields for future order persistence:

```go
type OrderPricing struct {
    BasePrice      int64
    PromoDiscount  int64
    PaymentMarkup  int64
    FinalPrice     int64

    PaymentMethod  string
    PromoCode      string
    PromoKind      string
}
```

This is important because pricing rules may change over time, while stored order must remain explainable.

## Responsibilities

Pricing logic belongs to quote/pricing layer.

It must not be spread across:

- Telegram runtime
- screen rendering
- callback parsing
- generic catalog navigation

Flow should only:

- collect input
- request quote
- display quote
- confirm order
