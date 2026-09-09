# Architecture

Sitekit is an opinionated server-side rendering framework built on top of Go's `net/http` and standard-library primitives.

This document describes the architectural model that follows from the design principles in [`design-principles.md`](./design-principles.md). It is concerned with responsibilities, dependency direction, and request flow rather than individual API signatures. Concrete public APIs are tracked separately by the API lifecycle registry.

## 1. Architectural position

Sitekit is an application-layer framework, not a replacement for the Go HTTP stack.

The architectural boundary is therefore:

```text
application
    |
    v
Sitekit application layer
    |
    +-- app
    +-- router
    +-- page
    +-- render
    +-- middleware
    +-- config
    +-- static / seo / errors
    |
    v
Go standard library
    |
    +-- net/http
    +-- html/template
    +-- context
    +-- io/fs
    +-- log/slog
    +-- encoding/json
```

Sitekit adds application structure around standard-library capabilities. It must not make the standard library inaccessible to the application.

## 2. Core architectural units

### `app`

`app` owns the application runtime boundary.

Its responsibilities include:

- application configuration;
- the underlying `http.ServeMux`;
- request-context creation;
- conversion between Sitekit handlers and `http.Handler`;
- application middleware;
- error handling;
- server startup and graceful shutdown.

`App` is the composition root. It should not become a container for unrelated application features.

### `router`

`router` provides a higher-level registration API over `http.ServeMux`.

The router is an adapter, not an independent HTTP server. Route matching ultimately belongs to the standard library. Grouping and handler middleware are composition conveniences around that primitive.

### `page`

`page` defines the application-facing representation of a renderable page.

A page carries the template identity, optional layout, metadata, and data supplied to the renderer. Generic typed pages provide compile-time structure for application data without introducing a custom template language.

### `render`

`render` owns server-side HTML rendering.

The renderer uses `html/template` and `io/fs` rather than defining a Sitekit template engine. Template functions, layouts, caching, and development reload behavior are renderer concerns.

### `middleware`

Middleware is split conceptually into two layers:

1. HTTP middleware compatible with `func(http.Handler) http.Handler`;
2. handler middleware compatible with Sitekit's `HandlerFunc` abstraction.

These layers must remain explicit. The framework should not hide the transition between them when doing so would reduce control or make execution order unclear.

### supporting packages

`config`, `errors`, `static`, `seo`, and development tooling provide focused application-level capabilities. They should remain independently understandable and should not become prerequisites for the core request path unless a documented contract requires it.

## 3. Request lifecycle

The normal request path is:

```text
HTTP request
    |
    v
http.Server
    |
    v
App HTTP middleware
    |
    v
http.ServeMux
    |
    v
route + handler middleware
    |
    v
Sitekit Context
    |
    v
application HandlerFunc
    |
    +----> JSON response
    |
    +----> redirect
    |
    +----> page.Page -> Renderer -> html/template
    |
    +----> error -> ErrorHandler
    v
HTTP response
```

The important property is that Sitekit does not replace the HTTP lifecycle. It adds a small framework lifecycle inside it.

## 4. Context boundary

`app.Context` is an application convenience around `http.ResponseWriter` and `*http.Request`.

The underlying request remains available. `Context.Context()` returns the request's standard `context.Context`, request parameters come from `Request.PathValue`, and query access is based on `Request.URL.Query()`.

This boundary must remain transparent enough that an application can drop down to `net/http` whenever necessary.

## 5. Rendering boundary

Rendering is deliberately separated from routing and request handling:

```text
HandlerFunc
   |
   v
page.Page
   |
   v
Renderer
   |
   v
html/template + io/fs
```

A handler decides *what* page to render. The renderer decides *how* that page becomes HTML. The page abstraction carries the data and template identity between the two.

The renderer writes the HTTP status and content type as part of rendering. Therefore rendering is an HTTP response operation, not a pure string-generation utility.

## 6. Error boundary

Handlers return `error` instead of writing framework-specific error responses directly.

The application converts a handler error into an HTTP response through `ErrorHandler`. A default handler exists when the application does not configure one.

This establishes one important invariant:

> Returning an error from a Sitekit handler transfers error-response responsibility to the application error boundary.

Applications that need lower-level control can use `http.Handler` directly or write to the response themselves.

## 7. Middleware composition

Middleware order is part of behavior and therefore part of the architectural contract.

For a sequence `A, B, C`, wrapping is performed from the end of the list so that execution is:

```text
A -> B -> C -> handler
```

and unwinding follows the reverse order.

Route groups inherit parent handler middleware and append child middleware. Application-level HTTP middleware wraps the final mux handler.

Future changes must preserve this model unless an explicit compatibility decision says otherwise.

## 8. Dependency direction

The preferred dependency direction is toward lower-level primitives:

```text
app -> config / render / livereload
router -> app
render -> page
page -> standard library types only
middleware -> app
application -> any public Sitekit package
```

The core must avoid circular conceptual dependencies. In particular:

- `page` must not depend on `app`;
- `render` must not own routing concerns;
- `router` must not become a rendering layer;
- supporting packages must not silently become global application state.

## 9. Composition over centralization

Sitekit should expose small composable units rather than a single global framework object containing every capability.

The intended composition model is:

```go
cfg := config.Load()
app := app.New(cfg)
r := router.New(app)
renderer := render.New(templateFS, cfg.IsDevelopment())

app = app.WithRenderer(renderer)

r.Get("/", handler)
```

The exact API may evolve; the architectural property must not: components should be replaceable and independently understandable.

## 10. Standard-library escape hatches

Every major abstraction must have an obvious path back to the standard library.

Examples include:

- `App.Mux` exposes the underlying `*http.ServeMux`;
- handlers receive the underlying `*http.Request` and `http.ResponseWriter` through `Context`;
- middleware can be written directly as `func(http.Handler) http.Handler`;
- templates are standard `html/template` templates;
- template files are supplied through `io/fs`.

An abstraction that prevents this escape path is architecturally suspect.

## 11. Development-only behavior

Development conveniences, such as live reload, must remain orthogonal to production request semantics.

Development mode may add watchers, endpoints, or middleware, but the production architecture must not depend on them. A feature is not considered core merely because it is enabled automatically during development.

## 12. Architectural invariants

The following invariants are normative:

1. `net/http` remains the underlying HTTP model.
2. Routing remains based on `http.ServeMux` rather than a separate routing engine.
3. Server-side HTML rendering remains based on `html/template`.
4. Page data is application-owned; Sitekit does not introduce a global model container.
5. Middleware composition is explicit and order-sensitive.
6. Handler errors cross one explicit application error boundary.
7. Public abstractions retain an escape hatch to lower-level Go primitives.
8. Supporting features must not force unrelated dependencies into the core.

These invariants are the architectural consequences of the principles documented in `design-principles.md`. Public API details and their compatibility status belong in `contracts.md` and the API lifecycle registry.
