## Core principles:
1. Sitekit **extends** ```net/http```;
2. Standard library first;
3. No vendor lock-in;
4. Escape hatch to ```stdlib```;
5. Components must be composable;
6. Public API must remain small;
7. Convenience must not remove control;
8. Sitekit is SSR/application tooling, not an everything-framework;
9. Boring core, powerful edges;
10. Go compatibility is a long-term requirement.

## Before adding a feature, ask:

1. Is this truly application-level?
2. Can it be built on the standard library?
3. Can it be used independently?
4. Does it unnecessarily expand the public API?
5. Does the user retain control?
6. Does Go already provide this capability?