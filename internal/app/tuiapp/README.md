# tuiapp

`tuiapp` is the composition root for the terminal operator interface of the
platform.

This project is evolving not as a simple bot manager, but as an e-commerce
platform with multiple surfaces of operation and sales flow delivery:

- **Telegram** as a customer-facing and operator-facing channel
- **TUI** as an operator control surface
- **Web** as a future admin and store-management surface
- potentially other runtime-specific surfaces later

In this model, `tuiapp` is responsible for assembling dependencies required by
the terminal interface and bot runtime control flow.

---

## Purpose

`tuiapp` exists to keep `cmd/tui` thin and focused.

Instead of assembling PostgreSQL, repositories, services and runtime manager
directly in `main.go`, the package provides a dedicated wiring layer for the
TUI application slice.

---

## Responsibilities

`tuiapp` is responsible for:

- opening shared PostgreSQL resources
- building storage-backed bot configuration service
- building bot runtime manager
- returning a ready-to-use application container for `cmd/tui`

---

## Non-responsibilities

`tuiapp` does **not**:

- load `.env`
- load top-level application config
- configure the top-level logger
- run seed logic
- run bootstrap logic
- run the TUI itself
- contain business logic
- decide product, order, marketplace or bot behavior

Those concerns belong either to entrypoints (`cmd/tui`) or to dedicated
application/domain packages.

---

## Why this package exists

Before `tuiapp`, dependency wiring for the terminal application lived directly
inside `cmd/tui/main.go`.

That worked for the first iteration, but over time it made the entrypoint carry
too many concerns:

- environment handling
- logger setup
- PostgreSQL startup
- repository wiring
- service construction
- runtime manager construction
- startup flow

`tuiapp` extracts only the dependency assembly part, while preserving startup
policy in the entrypoint.

---

## Architectural role

`tuiapp` is a **composition root** for one specific platform surface:
the terminal operator interface.

This matters because the system is expected to grow into multiple application
surfaces, not one monolithic interface. A likely long-term structure is:

```text
internal/app/
├── app.go
├── bootstrap/
├── pgapp/
├── seed/
├── tuiapp/
├── webapp/
└── tgapp/
```

In this structure:

- `tuiapp` wires terminal operator functionality
- `webapp` can later wire admin/storefront web functionality
- `tgapp` can later wire Telegram operator/admin flows

This keeps composition explicit and scoped to the surface that actually uses it.

---

## Current dependencies assembled by tuiapp

At the moment `tuiapp` wires:

- PostgreSQL pool via `internal/app/pgapp`
- bot configuration storage via `internal/botconfig/postgres`
- bot configuration service via `internal/botconfig`
- runtime manager via `internal/manager`

The resulting container is intended for `cmd/tui`.

---

## Package structure

Typical structure:

```text
internal/app/tuiapp/
├── doc.go
├── config.go
├── app.go
└── builders.go
```

### `doc.go`
Package-level documentation and role description.

### `config.go`
Infrastructure configuration for the TUI application container.

### `app.go`
Application container definition and lifecycle helpers.

### `builders.go`
Dependency wiring for PostgreSQL, repositories, services and manager.

---

## App container

The package exposes an `App` container similar to:

```go
type App struct {
    Pool      *pgxpool.Pool
    BotConfig *botconfig.Service
    Manager   *manager.Manager
}
```

This container is intentionally small and only holds top-level dependencies
needed by the TUI entrypoint.

It should not be used as a general service locator across the system.

---

## Configuration

`tuiapp` uses PostgreSQL configuration provided by `pgapp` and may include
startup-related infrastructure parameters, such as database open timeout.

Example shape:

```go
type Config struct {
    Postgres      pgapp.Config
    OpenDBTimeout time.Duration
}
```

The package may provide `LoadConfigFromEnv()` for its own infrastructure needs,
but top-level config loading remains outside of `tuiapp`.

---

## Usage

Typical usage from `cmd/tui` looks like:

```go
cfg := config.MustLoad()

appLogger, err := logger.Setup(cfg.Env)
if err != nil {
    log.Fatalf("setup logger: %v", err)
}

appCfg := tuiapp.LoadConfigFromEnv()

app, err := tuiapp.New(ctx, appCfg, &demoRunner{}, appLogger.Logger)
if err != nil {
    appLogger.Error("build tui app", "err", err)
    os.Exit(1)
}
defer app.Close()
```

After wiring is complete, `cmd/tui` can continue with startup policy:

- ensure demo data
- bootstrap enabled bots
- launch terminal UI

---

## Startup policy boundary

A key design decision is that `tuiapp` performs **wiring only**.

That means the following still belong to `cmd/tui`:

- loading `.env`
- loading platform config
- setting up structured logging
- seeding development data
- bootstrapping enabled bots
- logging startup results
- starting the terminal UI loop

This keeps the package focused and prevents it from becoming a “god package”.

---

## Relation to bot runtime

`tuiapp` does not own a specific bot runtime implementation.

Instead, it accepts a `manager.Runner` from the outside. This keeps the package
agnostic to whether the runtime is:

- a demo runner
- a Telegram runner
- a future marketplace or integration runtime

That matches the current `manager` design where the manager works with a single
`Runner` interface.

---

## Documentation style

Code inside `tuiapp` should be documented in Go doc style.

That means:

- package-level purpose in `doc.go`
- exported types and functions documented with clear intent
- comments describe responsibilities and boundaries, not obvious mechanics only

This package is an infrastructure assembly layer, so clear comments matter.
Future contributors should understand *why* it exists, not only *what* it does.

---

## Design principles

`tuiapp` follows these principles:

- keep entrypoints thin
- keep wiring explicit
- keep business logic outside composition roots
- assemble only the dependencies required for this surface
- prefer narrow, readable application containers
- preserve room for future `webapp` and `tgapp` packages

---

## Future evolution

Possible next steps after introducing `tuiapp`:

- move `demoRunner` into a dedicated runtime package
- introduce a real Telegram runtime implementation
- add richer TUI control over bot lifecycle
- support operator workflows beyond bot monitoring
- evolve the platform toward a unified e-commerce control plane

If the project later grows an online terminal storefront or terminal-assisted
operator sales flow, `tuiapp` remains the right place for wiring that surface.

---

## Summary

`tuiapp` is not “the app package for everything”.

It is the composition root for the **terminal operator surface** of the broader
e-commerce platform.

Its job is simple and important:

- assemble infrastructure
- expose a ready-to-use container
- keep `cmd/tui` readable
- leave business behavior to the proper layers
