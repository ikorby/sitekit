# API Lifecycle

This document defines the lifecycle and compatibility commitment of Sitekit's public API.

It answers two questions:

- How important is this API to the framework's architecture?
- How much does Sitekit commit to preserving it?

`contracts.md` defines API behavior. `api-lifecycle.md` defines how APIs may evolve.

## Lifecycle Levels

### Stable

The API is part of Sitekit's long-term public surface.

Its documented source and behavioral contracts must be preserved. Breaking changes require an explicit compatibility policy.

### Provisional

The API is public and usable, but its contract is not final.

It may change before the first stable release. Changes affecting users must be documented.

### Experimental

The API is intended for experimentation. It may change substantially or be removed without preserving compatibility.

### Deprecated

The API remains available temporarily but is no longer recommended.

Deprecation must include a reason, migration guidance, and a replacement when available.

### Internal

The implementation detail is not part of the public API contract and may change freely.

## API Roles

Architectural role and lifecycle status are independent classifications.

### Core

Defines Sitekit's fundamental programming model.

Examples: application handling, routing, pages, and rendering.

### Supporting

Provides additional functionality around the core without defining its fundamental programming model.

Examples: configuration, SEO, static files, and convenience middleware.

### Tooling

Supports development workflows rather than application runtime behavior.

Examples: CLI, project generators, and live reload.

A Core API may remain Provisional. A Supporting API may become Stable. Tooling may evolve independently of runtime APIs.

## Pre-stable API Evolution

Before `v1.0.0`, Provisional APIs may undergo source-level, behavioral, and structural changes.

Experimental APIs may change or be removed without prior deprecation.

Changes affecting application code or observable behavior must be documented.

### v0.2.0 milestone

Version `v0.2.0` is a foundational architecture and contract milestone.

It covers:

- Response commitment and response helpers;
- Atomic rendering and JSON responses;
- Error handling after response commitment;
- Application lifecycle;
- Concurrency guarantees;
- Public API boundaries and mutable state;
- Middleware semantics;
- Router semantics.

Reaching `v0.2.0` does not automatically make any API Stable.

## Registry

The registry identifies the architectural role and intended lifecycle of individual public API symbols.

Each entry represents an explicit project decision. Exported Go identifiers are not automatically considered stable public contracts.

Package-level documentation remains the source for detailed API descriptions. This registry defines how strongly Sitekit commits to preserving each listed API.

| Area | API | Role | Lifecycle | Intended Direction |
|---|---|---|---|---|
| `app` | `App` | Core | Provisional | Application runtime and composition root. |
| `app` | `New` | Core | Provisional | Primary application constructor. |
| `app` | `Option` | Core | Provisional | Application constructor configuration. |
| `app` | `App.Config` | Core | Provisional | Application configuration access. |
| `app` | `App.Renderer` | Core | Provisional | Renderer access and configuration. |
| `app` | `App.Mux` | Core | Provisional | Underlying standard-library router access. |
| `app` | `App.Logger` | Core | Provisional | Application logger access. |
| `app` | `App.ErrorHandler` | Core | Provisional | Application error-boundary configuration. |
| `app` | `WithRenderer` | Core | Provisional | Renderer injection. |
| `app` | `WithLogger` | Core | Provisional | Logger injection. |
| `app` | `WithErrorHandler` | Core | Provisional | Error-handler injection. |
| `app` | `WithMiddleware` | Core | Provisional | Application HTTP middleware registration. |
| `app` | `Middleware` | Core | Provisional | Standard HTTP middleware abstraction. |
| `app` | `HandlerFunc` | Core | Provisional | Primary Sitekit handler abstraction. |
| `app` | `HandlerMiddleware` | Core | Provisional | Handler-level middleware abstraction. |
| `app` | `ErrorHandler` | Core | Provisional | Explicit application error boundary. |
| `app` | `Wrap` | Core | Provisional | Handler middleware composition. |
| `app` | `ToHandler` | Core | Provisional | Conversion to standard HTTP handler. |
| `app` | `Context` | Core | Provisional | Request-scoped application context. |
| `app` | `Context.Writer` | Core | Provisional | Sitekit response boundary. |
| `app` | `Context.Request` | Core | Provisional | Underlying HTTP request access. |
| `app` | `Context.Context` | Core | Provisional | Underlying request context access. |
| `app` | `WithValue` / `Context.Value` | Core | Provisional | Request-scoped value storage and retrieval. |
| `app` | `Context.Param` / `Context.Query` | Core | Provisional | Request parameter and query accessors. |
| `app` | `Context.Render` / `Context.JSON` / `Context.Redirect` | Core | Provisional | Response helpers. |
| `app` | `Response` | Core | Provisional | HTTP response boundary and state tracking. |
| `app` | `NewResponse` | Core | Provisional | Response boundary constructor. |
| `app` | `Response.Header` | Core | Provisional | Response header access. |
| `app` | `Response.WriteHeader` | Core | Provisional | Response status commitment. |
| `app` | `Response.Write` | Core | Provisional | Response body writing and byte tracking. |
| `app` | `Response.Status` | Core | Provisional | Response status inspection. |
| `app` | `Response.BytesWritten` | Core | Provisional | Response byte-count inspection. |
| `app` | `Response.Committed` | Core | Provisional | Response commitment inspection. |
| `app` | `Response.Unwrap` | Core | Provisional | Escape hatch to the underlying HTTP response writer. |
| `router` | `Router` | Core | Provisional | Routing adapter over `http.ServeMux`. |
| `router` | `New` | Core | Provisional | Router constructor. |
| `router` | `Group` / `Use` | Core | Provisional | Route grouping and middleware composition. |
| `router` | `Handle` / `Get` / `Post` / `Put` / `Patch` / `Delete` / `Head` / `Options` / `Any` / `NotFound` | Core | Provisional | Route registration API. |
| `router` | `Static` | Core | Provisional | Static handler mounting. |
| `page` | `Meta` | Core | Provisional | Page metadata model. |
| `page` | `Page` | Core | Provisional | Renderer input contract. |
| `page` | `TypedPage[T]` | Core | Provisional | Generic typed page implementation. |
| `page` | `New` / `WithLayout` / `WithMeta` | Core | Provisional | Typed page construction and configuration. |
| `render` | `Renderer` | Core | Provisional | Server-side rendering abstraction. |
| `render` | `New` | Core | Provisional | Renderer constructor. |
| `render` | `Funcs` | Core | Provisional | Template function registration. |
| `render` | `Render` | Core | Provisional | HTML rendering entry point. |
| `config` | `Env` / `Development` / `Production` | Supporting | Provisional | Environment model. |
| `config` | `Config` | Supporting | Provisional | Runtime configuration model. |
| `config` | `Load` | Supporting | Provisional | Environment-based configuration loading. |
| `config` | `IsDevelopment` / `IsProduction` / `Addr` | Supporting | Provisional | Configuration helpers. |
| `errors` | `HTTPError` | Supporting | Provisional | Structured HTTP error representation. |
| `errors` | `New` / `Wrap` / `NotFoundError` | Supporting | Provisional | HTTP error constructors. |
| `errors` | `NotFound` / `Handler` | Supporting | Provisional | HTTP error and error-page handlers. |
| `seo` | `Defaults` / `DefaultsFromConfig` / `Apply` / `CanonicalURL` | Supporting | Provisional | SEO metadata helpers. |
| `static` | `New` | Supporting | Provisional | Standard-library static file handler construction. |
| `middleware` | `Logger` / `Recovery` / `Security` | Supporting | Provisional | Convenience HTTP middleware. |
| `livereload` | `StartWatcher` / `SSEHandler` / `Middleware` | Tooling | Experimental | Development-only live reload functionality. |
| `cmd/sk` | CLI commands and generated project structure | Tooling | Experimental | Project creation and development tooling. |

The registry describes intentional project decisions. Exported Go identifiers are not automatically considered stable public contracts.

Package-level documentation remains the source for identifying individual exported APIs.

## Promotion and Removal

An API may become Stable only when:

1. Its architectural purpose is established;
2. Its behavior is documented in `contracts.md`;
3. Its compatibility guarantees are defined in `compatibility.md`;
4. Its implementation satisfies those guarantees;
5. Tests protect its observable behavior.

Removal rules:

- Experimental APIs may be removed directly.
- Provisional APIs may be changed or removed before `v1.0.0`, with documented migration information.
- Stable APIs must follow the compatibility policy.
- Deprecated APIs may be removed according to their documented removal conditions.

Lifecycle transitions must be reflected in this registry.

## Relationship with Other Documents

The documentation dependency is:

```text
Design Principles
       ↓
Architecture
       ↓
Contracts
       ↓
Compatibility
       ↓
Roadmap
```

`api-lifecycle.md` crosses this chain by tracking the lifecycle of concrete APIs.

A lifecycle change may require reviewing the corresponding architecture, contracts, compatibility rules, and roadmap.

## Source of Truth

Use:

- `architecture.md` for architectural boundaries and responsibilities;
- `contracts.md` for observable API guarantees;
- `compatibility.md` for compatibility policy;
- `roadmap.md` for planned changes and release milestones;
- `api-lifecycle.md` for API lifecycle status.

Lifecycle status does not describe implementation quality. It describes how strongly users may rely on an API remaining compatible.