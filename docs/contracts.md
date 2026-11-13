# Contracts

This document defines the observable guarantees provided by Sitekit's public API.

It describes what applications may rely on, independently of the framework's internal implementation.

The contracts in this document define the target behavior for v0.2.0. Existing implementation details are not considered guarantees unless explicitly stated here.

Normative terms:

- **MUST** — required behavior.
- **MUST NOT** — prohibited behavior.
- **SHOULD** — recommended behavior.
- **MAY** — permitted behavior.

## 1. HTTP Contract

Sitekit extends Go's `net/http` and MUST preserve its fundamental HTTP model.

- Applications MUST retain access to standard HTTP request and response types.
- Sitekit handlers MUST be interoperable with `http.Handler`.
- HTTP semantics MUST remain consistent with the standard library.
- Sitekit MUST NOT require applications to use a proprietary HTTP server or protocol.

## 2. Application Contract

The application is the root of runtime configuration and request handling.

- An application MUST be configurable before serving requests.
- Application configuration and route registration MUST be completed before the application starts serving.
- The application MUST support concurrent request processing.
- Server shutdown MUST allow in-flight requests to finish within the configured shutdown period.
- Unexpected server failures MUST be returned to the caller.
- Normal graceful shutdown MUST NOT be reported as an unexpected server failure.

Application lifecycle operations and their exact API are defined by the application package.

## 3. Handler Contract

A Sitekit handler processes one HTTP request and returns an error.

- Each handler invocation MUST receive a request-scoped Context.
- A handler MUST be able to access the underlying HTTP request through `Context.Request` and the Sitekit response boundary through `Context.Writer`.
The underlying `http.ResponseWriter` MUST remain accessible through `Context.Writer.Unwrap()`.
Direct writes through the unwrapped writer bypass Sitekit response tracking and MUST be treated as an explicit escape hatch.
- Returning `nil` indicates successful handler completion.
- Returning an error transfers error-response responsibility to the application error boundary.
- A handler MUST NOT assume that returning an error can replace a response that has already been committed.

Applications MAY use standard `http.Handler` implementations instead of Sitekit handlers.

## 4. Context Contract

Context represents one request and MUST NOT be shared between independent requests.

- The underlying request MUST remain accessible.
- The underlying standard context MUST preserve its cancellation, deadline, and value semantics.
- Request parameters MUST reflect the matched route.
- Query access MUST follow standard URL query semantics.
- Values stored in Context MUST remain scoped to the request.
- Context MUST NOT be used as global application storage.

The lifetime of Context ends when its request handling completes.

## 5. Router Contract

The router MUST remain an abstraction over `http.ServeMux`.

- Route matching MUST follow standard-library routing semantics.
- HTTP method and path patterns MUST be supported through the underlying mux.
- Route groups MUST preserve their configured prefixes.
- Child groups MUST inherit the handler middleware registered on their parent group at the time the child group is created.
Middleware registered on a parent group after child group creation MUST NOT be assumed to propagate to existing child groups unless explicitly documented otherwise.
Additional group middleware MUST be composed in registration order.
- Additional group middleware MUST be composed in registration order.
- Conflicting route registrations MUST follow the underlying mux's conflict behavior.

The router MUST NOT introduce an independent routing engine or silently change standard-library matching semantics.

## 6. Middleware Contract

Sitekit provides two distinct middleware types:

- HTTP middleware operating on `http.Handler`;
- Handler middleware operating on Sitekit handlers.

Middleware MUST execute in registration order.

For middleware A, B, and C, the execution order MUST be:

```text id="6b2h7x"
A → B → C → Handler
```

The unwinding order MUST be reversed.

Application HTTP middleware MUST wrap the application handler. Route middleware MUST apply only to the routes where it is registered or inherited.

The two middleware layers MUST remain distinguishable and composable.

## 7. Response Contract

An HTTP response has a commitment boundary.

Before commitment:

- Status and headers MAY be changed.
- The response body MUST NOT be considered final.

After commitment:

- The status code MUST NOT be replaced.
- Previously sent headers MUST NOT be modified.
- The response MUST NOT be replaced with a different response.

A response becomes committed when its headers are written or its body begins to be written. If no status is explicitly written, the standard HTTP behavior of an implicit `200 OK` applies.

All Sitekit response helpers MUST respect this boundary.

### Atomic response preparation

`Render` and `JSON` MUST prepare their output before committing a successful response.

- Template execution failures MUST NOT produce a partially committed successful response.
- JSON encoding failures MUST NOT produce a partially committed successful response.
- Once preparation succeeds, the response may be committed and written.
- Network or writer failures after commitment cannot be rolled back.

Direct writes through `http.ResponseWriter` remain permitted, but applications assume responsibility for their response state.

## 8. Error Contract

Handler errors MUST pass through one explicit application error boundary.

- An error returned before response commitment MUST be eligible for conversion into an HTTP error response.
- The configured ErrorHandler MUST receive the original handler error.
- A default error handler MUST exist when no custom handler is configured.
- An error returned after response commitment MUST NOT trigger an attempt to replace the existing response.
- Error handling MUST NOT silently violate the response commitment rules.

Recovery from panics is a separate middleware responsibility and MUST respect the same response commitment boundary.

## 9. Rendering and Page Contract

Sitekit rendering MUST use Go's `html/template`.

- Page data MUST remain application-owned.
- A Page MUST provide the information required by the renderer.
- Layouts and templates MUST follow the documented rendering conventions.
- Template execution failures MUST be returned as errors.
- Rendering MUST satisfy the atomic response preparation rules.
- Template function registration and renderer configuration MUST be completed before concurrent rendering begins.

Sitekit MUST NOT introduce a proprietary template language or bypass the standard template escaping model.

## 10. Configuration and Static Files

Configuration MUST have documented defaults and validation rules.

- Invalid configuration MUST produce an explicit error where validation is required.
- Environment-based configuration MUST follow documented variable names and parsing rules.
- Runtime configuration MUST NOT be silently mutated by unrelated requests.

Static-file helpers MUST preserve standard HTTP file-serving behavior and MUST NOT bypass application security boundaries.

## 11. Concurrency Contract

Sitekit MUST support concurrent HTTP requests.

- Shared framework components MUST be safe for their documented concurrent operations.
- Request-specific state MUST NOT leak between requests.
- Renderer operations MUST be safe for concurrent use after configuration is complete.
- Applications are responsible for synchronization of their own shared mutable data.
- Configuration and template functions MUST NOT be modified concurrently with request processing unless explicitly documented as safe.

The framework MUST NOT imply that application-owned data is automatically concurrency-safe.

## 12. Compatibility and Lifecycle

These contracts define intended behavior, not an automatic promise that every current API is already stable.

The lifecycle status of each API is defined in `api-lifecycle.md`.

Compatibility guarantees and breaking-change policy are defined in `compatibility.md`.

A change to observable behavior MUST be reviewed against this document before release.