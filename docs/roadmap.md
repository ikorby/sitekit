# Roadmap

This document describes the development direction and major milestones of Sitekit.

It defines priorities and expected outcomes rather than fixed release dates.

The roadmap is a planning document, not a compatibility guarantee. Release scope may change as architectural decisions, implementation experience, and testing provide new information.

## 1. Development Strategy

Sitekit development follows three priorities:

1. Establish a reliable architectural foundation.
2. Define and implement predictable public contracts.
3. Expand functionality without compromising the framework's design principles.

Before introducing new features, the existing core must have clearly defined behavior.

The project prioritizes correctness, composability, and a small public API over feature count.

## 2. Current Stage: Pre-1.0 Development

Sitekit is currently in its pre-stable development stage.

Public APIs are provisional unless explicitly marked otherwise in `api-lifecycle.md`.

During this stage:

- Architecture may be redesigned.
- Public APIs may change or be removed.
- Internal implementation may be refactored.
- Runtime behavior may be adjusted.
- Documentation may be revised to reflect better architectural decisions.

Backward compatibility with previous pre-1.0 releases is not a development constraint.

The primary objective is to establish a coherent framework before stabilizing its public surface.

## 3. v0.2.0 — Foundational Contracts

**Primary goal:** establish the fundamental runtime behavior and public contracts of Sitekit.

This release focuses on the correctness of existing functionality rather than expanding the framework with unrelated features.

### 3.1. Response Lifecycle

Establish a consistent response lifecycle.

Expected outcomes:

- Explicit response commitment tracking.
- Predictable status and header behavior.
- Correct handling of writes after commitment.
- Consistent behavior across response helpers.
- A clear escape hatch to the underlying `http.ResponseWriter`.

The `Response` abstraction must remain compatible with the standard HTTP model.

### 3.2. Rendering and JSON

Make response generation predictable and atomic before commitment.

Expected outcomes:

- HTML rendering prepares output before committing the response.
- JSON encoding completes before committing the response.
- Template and encoding errors do not produce partially committed successful responses.
- Network and writer failures after commitment are handled without attempting to replace the response.

### 3.3. Error Handling

Establish one consistent application error boundary.

Expected outcomes:

- Handler errors reach the configured `ErrorHandler`.
- Errors before commitment can produce an HTTP error response.
- Errors after commitment do not trigger an attempt to replace the existing response.
- Default error handling follows the same response lifecycle rules.
- Panic recovery remains a separate middleware responsibility.

### 3.4. Application Lifecycle

Define the lifecycle of an application.

Expected outcomes:

- Application configuration has a clear initialization phase.
- Route registration has explicitly defined lifecycle rules.
- Server startup and shutdown have predictable behavior.
- Graceful shutdown handles in-flight requests.
- Unexpected server failures are returned to the caller.
- Normal shutdown is distinguishable from server failure.

### 3.5. Concurrency and Mutable State

Define which framework components support concurrent use.

Expected outcomes:

- Request-specific state remains isolated.
- Shared framework components are safe for their documented concurrent operations.
- Renderer configuration has explicit mutation rules.
- Application configuration is not silently modified during request processing.
- Framework-owned mutable state has clear synchronization requirements.

### 3.6. Router and Middleware Semantics

Finalize the behavior of routing and middleware composition.

Expected outcomes:

- Routing remains based on `http.ServeMux`.
- Method and path matching follow standard-library semantics.
- Route groups preserve prefixes and middleware inheritance.
- Application HTTP middleware and handler middleware remain distinct.
- Middleware execution and unwinding order are explicitly defined.
- Static mounts and custom not-found handlers have documented behavior.

### 3.7. Configuration

Make configuration parsing and validation predictable.

Expected outcomes:

- Defaults are documented.
- Invalid values produce explicit errors where validation is required.
- Environment variable names and parsing rules are documented.
- Runtime configuration has a clearly defined mutation policy.

### 3.8. Public API Review

Review the public API before treating it as a long-term foundation.

Expected outcomes:

- Unnecessary exported fields and methods are reconsidered.
- Constructors and options follow consistent conventions.
- Package responsibilities remain clearly separated.
- Standard-library escape hatches remain available.
- API lifecycle statuses reflect actual design decisions.

### 3.9. Testing and Verification

Every fundamental contract must have tests protecting its observable behavior.

The release should include:

- Response lifecycle tests.
- Atomic rendering and JSON tests.
- Error handling tests before and after commitment.
- Application lifecycle tests.
- Concurrency tests where applicable.
- Router and middleware behavior tests.
- Configuration validation tests.
- Integration tests for the complete request lifecycle.

### v0.2.0 Completion Criteria

The milestone is considered complete when:

1. Fundamental runtime contracts are implemented.
2. Public behavior is covered by tests.
3. Response and error handling follow consistent rules.
4. Application lifecycle and concurrency semantics are defined.
5. Routing and middleware behavior are predictable.
6. The public API has undergone architectural review.
7. Documentation accurately describes the implemented behavior.

## 4. Post-v0.2 Development

After the foundational milestone, development can focus on improving the framework's application-level capabilities.

Potential areas include:

- Server-side rendering improvements.
- Page and layout composition.
- Static asset handling.
- SEO tooling.
- Development tooling and CLI improvements.
- Better application configuration.
- Developer experience and project conventions.

These areas are directions for future exploration, not commitments to specific features or release dates.

Every proposed feature must remain consistent with `design-principles.md`.

New functionality must not unnecessarily expand the core or duplicate capabilities already provided by Go.

## 5. Pre-1.0 Stabilization

Before the first stable release, Sitekit must undergo a dedicated stabilization phase.

The objective is to turn the evolving pre-1.0 implementation into a predictable public framework.

Expected outcomes:

- Review all public packages and exported symbols.
- Finalize observable behavior in `contracts.md`.
- Review every API lifecycle status.
- Remove or redesign APIs that do not fit the architecture.
- Establish the supported Go version policy.
- Finalize the compatibility and deprecation policy.
- Verify examples and documentation against the actual implementation.
- Ensure important contracts are protected by automated tests.

No API should become Stable merely because it has existed for several releases.

Stability must be an explicit architectural and compatibility decision.

## 6. v1.0.0 — Stable Release

**Primary goal:** publish a coherent and reliable framework with an explicitly supported public API.

The first stable release should provide:

- A finalized public API.
- Documented runtime contracts.
- Defined compatibility guarantees.
- A clear API lifecycle registry.
- Predictable application and response lifecycles.
- Verified concurrency behavior.
- Complete essential documentation and examples.
- A documented Go version support policy.

Starting with `v1.0.0`, changes to stable APIs follow the compatibility policy in `compatibility.md`.

## 7. Non-Goals

The roadmap does not aim to turn Sitekit into an everything-framework.

The following remain outside the intended core scope:

- ORM or database abstraction.
- Dependency injection container.
- Custom template language.
- Frontend framework.
- Replacement for `net/http`.
- Independent HTTP server or routing engine.

Applications remain free to integrate these capabilities independently.

## 8. Roadmap Maintenance

The roadmap is reviewed when architectural decisions or development priorities change.

Changes to the roadmap should not silently redefine public contracts.

When a planned change affects existing behavior, the relevant documentation must be reviewed:

- `architecture.md` for structural changes.
- `contracts.md` for observable behavior.
- `api-lifecycle.md` for public API status.
- `compatibility.md` for compatibility implications.

The roadmap describes where Sitekit is heading. The other documents define what the framework is, how it behaves, and what applications may rely on.