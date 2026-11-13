# Compatibility Policy

This document defines Sitekit's backward compatibility policy and versioning rules.

It describes how changes to the public API, observable behavior, configuration, and development conventions are handled across releases.

## 1. Compatibility Scope

Compatibility applies to documented public interfaces and behavior.

The following areas may be covered by compatibility guarantees:

- Public Go APIs and exported types.
- HTTP handling and response behavior.
- Routing and middleware semantics.
- Rendering and page conventions.
- Configuration options and environment variables.
- Documented CLI commands and flags.
- Other behavior explicitly documented as part of the public contract.

Internal implementation details are not covered by compatibility guarantees.

The lifecycle status of individual APIs is defined in `api-lifecycle.md`.

## 2. Pre-1.0 Compatibility

Before `v1.0.0`, Sitekit does not guarantee backward compatibility.

During this period:

- Public APIs MAY be changed, renamed, or removed.
- Function signatures and exported types MAY change.
- Internal architecture MAY be redesigned.
- Runtime behavior MAY change.
- Configuration formats and defaults MAY change.
- Packages MAY be reorganized or removed.

Such changes do not require a formal deprecation period.

Release notes and migration information SHOULD describe changes that affect application code or behavior.

The absence of compatibility guarantees does not remove the responsibility to document significant changes.

### Pre-1.0 Development Priorities

Before `v1.0.0`, development prioritizes:

1. Correct architecture.
2. Consistent public API design.
3. Explicit and predictable contracts.
4. Reliable runtime behavior.
5. A maintainable implementation.

Preserving compatibility with earlier pre-1.0 releases MUST NOT prevent necessary improvements in these areas.

## 3. Stable Compatibility (v1.0.0 and Later)

Starting with `v1.0.0`, Sitekit follows Semantic Versioning.

Version format:

`MAJOR.MINOR.PATCH`

- **MAJOR** — incompatible changes to the public API or documented behavior.
- **MINOR** — backward-compatible functionality additions.
- **PATCH** — backward-compatible bug fixes and maintenance changes.

Stable public APIs MUST preserve their documented contracts across compatible releases.

A change is considered breaking when an application relying on documented behavior must modify its code or configuration to continue working correctly.

Examples include:

- Removing or renaming an exported symbol.
- Changing a public function signature incompatibly.
- Changing documented routing or middleware semantics.
- Removing a documented configuration option.
- Changing documented behavior in a way that breaks existing applications.

Internal refactoring that preserves public contracts is not a breaking change.

## 4. Deprecation Policy

After `v1.0.0`, a stable public API SHOULD be deprecated before removal.

A deprecation notice SHOULD include:

- The API being deprecated.
- The reason for deprecation.
- The recommended replacement, if available.
- Migration guidance.
- The conditions or release in which removal is expected, when known.

Deprecated APIs MAY remain available for a transition period.

The exact minimum deprecation period and removal conditions MUST be established before `v1.0.0`.

Experimental APIs are not subject to the stable API deprecation requirements.

## 5. Breaking Changes

After `v1.0.0`, breaking changes MUST be released through a new major version.

Breaking changes SHOULD:

- Have a documented technical justification.
- Include migration guidance.
- Be reflected in `contracts.md`.
- Be reflected in `api-lifecycle.md`.
- Be accompanied by appropriate tests.

Security fixes and changes required to correct behavior that violates documented contracts MAY require special handling. Such cases MUST be documented explicitly.

## 6. Go Version Compatibility

Sitekit's supported Go versions MUST be documented for each release.

The minimum supported Go version MAY change between major releases.

A change to the minimum supported Go version MUST be documented in the release notes and compatibility information.

Sitekit SHOULD use standard-library capabilities without unnecessarily restricting the range of supported Go versions.

## 7. Dependencies and Implementation

Sitekit does not guarantee compatibility of internal implementation details.

The following are not public compatibility commitments unless explicitly documented:

- Internal package structure.
- Private types and functions.
- Internal dependencies.
- Cache implementations.
- Internal synchronization mechanisms.
- Build and development infrastructure.

Applications SHOULD depend on documented public APIs rather than internal implementation details.

## 8. Release Documentation

Every release SHOULD provide enough information to understand its impact on existing applications.

Release notes SHOULD identify:

- New functionality.
- Changed behavior.
- Fixed issues.
- Deprecated APIs.
- Removed APIs.
- Breaking changes and migration requirements.

For pre-1.0 releases, migration notes are recommended whenever a change affects application code or behavior.

For post-1.0 releases, release documentation MUST clearly identify breaking changes.

## 9. Relationship to Other Documents

This document defines the compatibility policy.

Other documents define related aspects of the framework:

- `design-principles.md` — architectural and design principles.
- `architecture.md` — system structure and component responsibilities.
- `contracts.md` — observable behavior and public guarantees.
- `api-lifecycle.md` — lifecycle status of public APIs.
- `roadmap.md` — planned development and release scope.

Compatibility policy MUST be applied together with the contracts and API lifecycle definitions.

Before `v1.0.0`, these documents describe the intended direction of the framework. They MUST be reviewed against the actual implementation before the first stable release.