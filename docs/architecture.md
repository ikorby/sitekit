# Architecture

Sitekit is an opinionated server-side rendering framework built on Go's `net/http` and standard-library primitives.

This document defines the architectural model, component responsibilities, dependency direction, and system boundaries. It describes the intended architecture rather than individual API signatures or implementation details.

The model follows the principles defined in [`design-principles.md`](./design-principles.md).

## 1. Architectural Position

Sitekit is an application-layer framework, not a replacement for Go's HTTP stack.

Its architecture is built around three layers:

```text
Application
     |
     v
Sitekit Application Layer
     |
     +-- app
     +-- router
     +-- page
     +-- render
     +-- supporting packages
     |
     v
Go Standard Library
     |
     +-- net/http
     +-- html/template
     +-- context
     +-- io/fs
     +-- log/slog
     +-- encoding/json
```

Sitekit adds structure and convenience around standard-library capabilities while preserving access to their underlying abstractions.

## 2. Architectural Components

### Core

Core packages define the fundamental programming model of Sitekit.

**`app`**

Owns the application runtime and composition boundary. It coordinates configuration, request handling, HTTP middleware, handler integration, error handling, and server lifecycle.

The application object must remain focused on runtime responsibilities rather than becoming a container for unrelated features.

**`router`**

Provides route registration and grouping over `http.ServeMux`.

It is an adapter, not an independent routing engine. Route matching remains the responsibility of the standard library.

**`page`**

Defines the application-facing representation of a renderable page, including its template identity, layout, metadata, and application-owned data.

Typed page implementations provide compile-time structure without introducing a custom template language.

**`render`**

Owns server-side HTML rendering using `html/template` and `io/fs`.

Template loading, execution, layouts, function registration, and caching belong to this component. Rendering must remain independent of routing concerns.

### Supporting Packages

Supporting packages provide additional application-level functionality without defining the core programming model.

These include:

- `config` — application configuration;
- `errors` — structured HTTP errors;
- `middleware` — reusable middleware;
- `static` — static-file utilities;
- `seo` — metadata and SEO helpers.

Supporting packages must remain focused and must not introduce unnecessary dependencies into the core.

### Tooling

Development tools, including the CLI and live reload, support development workflows rather than production runtime requirements.

They may evolve independently of the core architecture.

## 3. Request Lifecycle

The normal request flow is:

```text
Incoming HTTP Request
        |
        v
http.Server
        |
        v
Application HTTP Middleware
        |
        v
http.ServeMux
        |
        v
Route and Handler Middleware
        |
        v
Sitekit Handler
        |
        v
     Context
    |       |
    v       v
 Request  Response
             |
             v
     http.ResponseWriter
```

Sitekit introduces an application-level lifecycle inside the standard HTTP lifecycle. It does not replace the underlying request and response model.

For each incoming request, Sitekit creates a request-scoped `Response` wrapper around the underlying `http.ResponseWriter`. The handler accesses this wrapper through `Context.Writer`, while the original request is available through `Context.Request`.

The `Response` wrapper tracks response commitment and write operations while delegating actual HTTP response operations to the underlying `http.ResponseWriter`.

Calling `Response.Unwrap()` exposes the underlying writer directly. Writes through the unwrapped writer bypass the tracking guarantees provided by the wrapper.
## 4. Context and HTTP Boundaries

`app.Context` provides request-scoped access to the HTTP request, Sitekit response boundary, and application-level request helpers.

It exposes three distinct boundaries:
- `Context.Request` provides access to the underlying `*http.Request`.
- `Context.Writer` provides access to Sitekit's `*app.Response` abstraction.
- `Context.Writer.Unwrap()` provides an escape hatch to the underlying `http.ResponseWriter`.
  The Response abstraction tracks response status, bytes written, and commitment state while preserving the standard `http.ResponseWriter` model.

Applications may use the underlying standard-library response writer when necessary, but direct writes through the unwrapped writer bypass Sitekit's response tracking and must be treated as an explicit escape hatch.

## 5. Rendering and Error Boundaries

Rendering is separated from request routing and application handler logic.

A handler determines which page should be rendered. The page carries template identity and application data. The renderer converts that page into HTML.

The renderer uses standard Go template primitives and does not own application routing or business logic.

Handler errors are delegated to the application error boundary. This provides one explicit integration point for translating application errors into HTTP responses.

The exact response and error-handling guarantees belong in `contracts.md`.

## 6. Middleware Composition

Sitekit supports two distinct middleware layers:

1. HTTP middleware operating on `http.Handler`;
2. Handler middleware operating on Sitekit handler functions.

HTTP middleware wraps the application HTTP handler. Handler middleware composes around individual route handlers.

Middleware composition is explicit and order-sensitive. Route groups may inherit parent middleware and extend it with additional middleware.

The exact ordering and inheritance rules are defined in `contracts.md`.

## 7. Dependency Direction

Dependencies should follow a clear direction:

```text
app         -> config, render, livereload
router      -> app
render      -> page
middleware  -> config
static      -> config
seo         -> config, page
errors      -> app, page
livereload  -> config
page        -> standard library
```

Arrows represent Go package import dependencies. An arrow from package A to package B means that package A imports package B.

Applications depend on Sitekit's public packages, while Sitekit packages depend on the standard library and, where necessary, on other Sitekit packages. The dependency graph must remain acyclic.

The architecture must avoid circular dependencies and unnecessary coupling.

In particular:

- `page` must not depend on `app`;
- `render` must not own routing;
- `router` must not become a rendering layer;
- supporting packages must not introduce implicit global application state.

## 8. Composition and Extensibility

Sitekit favors small, composable components over a centralized framework object.

Applications should be able to configure, replace, or bypass individual components where appropriate.

Major abstractions must preserve standard-library escape hatches:

- HTTP handlers remain compatible with `net/http`;
- applications retain access to the underlying request and response writer;
- routing is based on `http.ServeMux`;
- templates use `html/template`;
- template files are provided through `io/fs`.

Convenience must not remove control over the underlying Go primitives.

## 9. Development and Production Separation

Development-only functionality, such as live reload, must remain independent of production runtime requirements.

Development mode may introduce additional middleware, endpoints, or background processes. Production behavior must not depend on these features.

Development tooling is not part of the core runtime merely because it is automatically enabled in development mode.

## 10. Architectural Invariants

The following rules define the intended architecture:

1. `net/http` remains the underlying HTTP model.
2. Routing remains based on `http.ServeMux`.
3. Server-side rendering remains based on `html/template`.
4. Page data remains application-owned.
5. Middleware composition remains explicit and order-sensitive.
6. Handler errors pass through one explicit application error boundary.
7. Public abstractions preserve access to lower-level Go primitives.
8. Supporting functionality must not impose unrelated dependencies on the core.
9. Development tooling must remain independent of production requirements.

These invariants follow from the principles in `design-principles.md`.

Public API guarantees are defined in `contracts.md`. API lifecycle and compatibility commitments are tracked in `api-lifecycle.md` and `compatibility.md`.