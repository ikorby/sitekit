# API Lifecycle

This document is the registry for Sitekit's public API lifecycle.

It answers a different question from `contracts.md`:

> **How permanent is this API, and how may Sitekit change it?**

`contracts.md` defines what an API promises to do. This document defines how strongly Sitekit commits to preserving that API over time.

## Lifecycle levels

### Stable

The API is part of Sitekit's long-term public surface.

Stable APIs must preserve their documented source and behavioral contracts. Breaking changes require an explicit compatibility policy and are not introduced as ordinary maintenance.

Stable does **not** mean that the implementation can never change. It means that observable behavior and the public API are protected while the implementation remains free to evolve behind them.

### Provisional

The API is public and usable, but its contract is not yet considered final.

Provisional APIs may be changed before the first stable release when experience, testing, or architectural review shows that a different design is better. Changes should be documented when they affect users.

### Experimental

The API exists primarily to validate an idea or provide functionality whose final shape is intentionally undecided.

Experimental APIs may change substantially or disappear without preserving source compatibility. They should not be treated as foundations for long-lived application code.

### Deprecated

The API remains available for compatibility but is no longer the preferred direction.

Deprecated APIs should have a documented replacement or an explicit reason for removal. New code should not depend on them.

### Internal

The implementation detail is not part of Sitekit's public API contract.

Internal code may change freely as long as the documented public contracts remain valid.

## Registry

The registry below describes the intended lifecycle of the current public surface. A lifecycle label is an explicit project decision; it is not inferred merely from the fact that a Go identifier is exported.

| Area | API | Lifecycle | Intended direction |
|---|---|---|---|
| `app` | `App` | Provisional | Core application object; likely long-term API, but exported fields and lifecycle semantics need tightening before stabilization. |
| `app` | `New` | Provisional | Primary application constructor. |
| `app` | `Option` | Provisional | Constructor configuration mechanism. |
| `app` | `WithRenderer` | Provisional | Renderer injection/configuration. |
| `app` | `WithLogger` | Provisional | Logger injection/configuration. |
| `app` | `WithErrorHandler` | Provisional | Application error-boundary configuration. |
| `app` | `WithMiddleware` | Provisional | HTTP middleware configuration. |
| `app` | `Middleware` | Provisional | Standard `net/http` middleware shape; intended to become stable. |
| `app` | `HandlerFunc` | Provisional | Primary Sitekit handler abstraction; intended to become stable. |
| `app` | `HandlerMiddleware` | Provisional | Handler-level middleware abstraction; intended to become stable. |
| `app` | `ErrorHandler` | Provisional | Explicit application error boundary; intended to become stable after response semantics are finalized. |
| `app` | `Wrap` | Provisional | Handler middleware composition primitive. |
| `app` | `ToHandler` | Provisional | Escape hatch from Sitekit handlers to `http.HandlerFunc`. |
| `app` | `Context` | Provisional | Request-scoped Sitekit convenience layer over `net/http`. |
| `app` | `Context.Writer` / `Context.Request` | Provisional | Direct stdlib escape hatches; intended to remain available. |
| `app` | `Context.Context` | Provisional | Access to the underlying request context. |
| `app` | `WithValue` / `Context.Value` | Provisional | Request-scoped application values; semantics require further contract definition. |
| `app` | `Context.Param` / `Context.Query` | Provisional | Request convenience accessors built on stdlib request data. |
| `app` | `Context.Render` / `Context.JSON` / `Context.Redirect` | Provisional | Response helpers; behavior must be aligned with the final Response Contract. |
| `router` | `Router` | Provisional | Thin routing layer over `http.ServeMux`; intended long-term. |
| `router` | `New` | Provisional | Router constructor. |
| `router` | `Group` / `Use` | Provisional | Route composition and middleware inheritance. |
| `router` | `Handle` / `Get` / `Post` / `Put` / `Patch` / `Delete` / `Head` / `Options` / `Any` / `NotFound` | Provisional | Main route registration surface; semantics should remain based on `http.ServeMux`. |
| `router` | `Static` | Provisional | Static-handler mounting convenience; exact contract still needs review. |
| `page` | `Meta` | Provisional | Page metadata model; intended long-term but still subject to SEO/API review. |
| `page` | `Page` | Provisional | Core renderer input contract; intended long-term. |
| `page` | `TypedPage[T]` | Provisional | Generic typed page implementation; intended long-term. |
| `page` | `New` / `WithLayout` / `WithMeta` | Provisional | Typed page construction helpers. |
| `render` | `Renderer` | Provisional | SSR renderer; core abstraction but rendering/response semantics are not yet final. |
| `render` | `New` | Provisional | Renderer constructor. |
| `render` | `Funcs` | Provisional | Template function registration; exact lifecycle and mutation semantics need review. |
| `render` | `Render` | Provisional | Rendering entrypoint; must satisfy the final Response and Rendering Contracts. |
| `config` | `Env` / `Development` / `Production` | Provisional | Environment model; intended long-term. |
| `config` | `Config` | Provisional | Runtime configuration model; field-level compatibility needs explicit policy. |
| `config` | `Load` | Provisional | Environment-based configuration loading. |
| `config` | `IsDevelopment` / `IsProduction` / `Addr` | Provisional | Configuration helpers. |
| `errors` | `HTTPError` | Provisional | Structured HTTP error; intended long-term after error semantics are finalized. |
| `errors` | `New` / `Wrap` / `NotFoundError` | Provisional | HTTP error constructors. |
| `errors` | `NotFound` | Provisional | Error-page handler convenience; exact fallback behavior needs contract review. |
| `seo` | Public SEO helpers and types | Provisional | Useful application-level tooling; exact API surface requires contract review before stabilization. |
| `static` | Public static-file helpers | Provisional | Application-level convenience; exact API surface requires contract review before stabilization. |
| `middleware` | Public middleware helpers | Provisional | Convenience middleware; each API must be reviewed individually before stabilization. |
| `livereload` | Public development helpers | Experimental | Development tooling; implementation and API shape may change substantially. |
| `cmd/sk` | CLI commands and generated project shape | Experimental | Tooling is intentionally allowed to evolve quickly before the framework reaches a stable release. |

## Current commitment level

At the current pre-v1 stage, **no runtime API is declared Stable yet**.

This is deliberate. The architecture is being documented before the public surface is frozen, and the response, error, rendering, concurrency, and lifecycle contracts still need to be made explicit.

The intended progression is:

```text
Experimental
     ↓
Provisional
     ↓
Stable
```

An API should move to Stable only when:

1. its purpose is understood and consistent with the architecture;
2. its behavior is explicitly defined in `contracts.md`;
3. its compatibility expectations are documented in `compatibility.md`;
4. its implementation satisfies the contract;
5. its tests protect the important observable behavior.

## What lifecycle status does not mean

Lifecycle status does not describe implementation quality.

An Experimental API can be well implemented. A Provisional API can already be heavily tested. A Stable API may still receive internal refactors.

The status answers only this question:

> **How much should an application depend on the exact public shape and behavior of this API?**

## Relationship with other documents

The intended documentation dependency is:

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

`api-lifecycle.md` crosses this chain as the registry of concrete APIs:

```text
             API Lifecycle
            /      |      \
           ↓       ↓       ↓
Architecture → Contracts → Compatibility → Roadmap
```

When an API changes lifecycle status, the corresponding contract and compatibility rules must be reviewed as well.

## Rules for lifecycle changes

Lifecycle changes should be explicit.

### Promoting an API

Moving an API from Experimental or Provisional to Stable requires a contract review. The project should be able to state exactly what users may rely on.

### Deprecating an API

A deprecated API should not be silently changed into a different abstraction. Deprecation should preserve enough behavior for existing users to migrate deliberately.

### Removing an API

Removal is a compatibility event, not merely a cleanup. The API must first be marked Deprecated unless an explicitly documented pre-stable policy permits direct removal.

## Source of truth

This document is the lifecycle registry, not the API reference.

If a question concerns **what an API does**, use `contracts.md` and the package documentation.

If a question concerns **whether the API may change**, use this document and `compatibility.md`.

If a question concerns **when a planned change will happen**, use `roadmap.md`.
