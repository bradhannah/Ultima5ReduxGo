# Coding Conventions

This document defines the coding standards for the project. It balances idiomatic Go, public Google Go style guidance, and the project’s practical needs (classic game fidelity, clarity, and incremental development).

Goals:
- Keep core gameplay logic clear and simple.
- Favor readability over cleverness.
- Prefer conventions that make the code testable, deterministic, and reusable.

## Language and Tooling

- Go version: the repository’s go.mod defines the version; keep code compatible with it.
- Formatting: gofmt (via goimports). CI should fail if formatting is not applied.
- Linting: follow common static checks (e.g., golangci-lint defaults) where reasonable.
- Build tags: use sparingly (e.g., debug vs release).

## Package Structure and Boundaries

- Core logic should be independent of the rendering/input backend. Adapter layers (e.g., Ebiten) call into core, not vice versa.
- Keep packages focused and small; avoid cyclic dependencies.
- Exported types and functions should be stable and usable by external consumers where reasonable.

## Naming

- Exported identifiers: PascalCase. Unexported: camelCase.
- Initialisms: follow Go’s convention (ID, URL, HTTP). Example: UserID, not UserId.
- File names: lowercase with underscores only when it improves clarity; prefer consistent grouping by feature/domain.

## Imports and Dependencies

- Standard library first, external second, internal last; keep groups separated by a blank line.
- Do not import the adapter layer (e.g., rendering/input libraries) into core logic. Inject via interfaces if needed.
- Keep external dependencies minimal and pinned via go.mod/go.sum.

## Errors and Logging

- Return errors instead of panicking in normal control flow.
- Wrap underlying errors with context: `fmt.Errorf("loading sprite: %w", err)`.
- Temporary hard exits:
  - During early development, a critical path may use `log.Fatal` as a stub for “must not happen” states to catch regressions quickly.
  - Always add a comment `// TODO: soften to recoverable error` so we can convert to softer handling later.
- Logging should be concise and actionable. Avoid spamming logs in hot paths; prefer rate-limiting if necessary.

## Time, Randomness, and Ticks

- Do not call `time.Now()` inside core logic or animation helpers that affect gameplay determinism.
- Use a central game clock (ticks/elapsedMs) to drive:
  - Animations (waterfalls, clocks, idle cycles).
  - NPC schedules and time-of-day effects.
  - AI cadence and spawn checks.
- Randomness:
  - Centralize PRNG state and seed per session.
  - For deterministic tests, allow fixed seeds.

## Loops and Collections

- Prefer `range` loops for readability and safety when index is not needed:
  - Good: `for _, v := range items { ... }`
  - Use index-based loops when:
    - You need the index for math or spatial mapping.
    - You need to mutate by index efficiently.
    - You’re iterating over a fixed-size array where index is semantic.
- Preallocate slices with capacity when the size is known (`make(T, 0, n)` or `make(T, n)`).
- Avoid unnecessary allocations in hot paths; reuse buffers where appropriate.

## Data-Driven vs Hardcoded Logic

- Favor data-driven definitions for tile attributes, schedules, and references.
- It is acceptable to hardcode small, well-named rules when it materially improves clarity (e.g., special-case mirror behavior, chair/ladder substitutions).
- When hardcoding, isolate rules behind small helpers so they can be replaced by data later if needed.

## Rendering and Animation

- All animations should read from the tick/elapsed time (not `time.Now`) to keep visuals in sync with gameplay.
- Prefer simple, deterministic animation helpers that accept `(spriteIndex, positionHash, elapsedMs)`.
- Keep per-frame allocations to a minimum (reuse images/options where possible).

## AI and Pathfinding

- Cache computed paths (A*) on the unit and consume steps over subsequent ticks.
- Recompute on cooldowns or when blocked; fall back to simple greedy movement if needed.
- Respect terrain passability by agent type (avatar, vehicle, land/water enemy).
- Do not hard-exit on missing paths; use a clear fallback and log once with a TODO to soften behavior later.

## Concurrency

- Prefer single-threaded game-state mutation during the update step.
- If background loading is needed, communicate via channels or synchronized buffers, and apply results during the main update.

## Comments and Documentation

- Use doc comments for exported identifiers; explain the why, not the obvious what.
- Keep comments up-to-date with behavior; incorrect comments are worse than none.
- Add package-level docs for domains with non-trivial rules (e.g., AI, schedule resolution).

## Testing

- Unit tests for deterministic pieces: animation frame selection, schedule resolution, pathfinding.
- Use fixed seeds and a mock clock for deterministic tests.
- Avoid brittle pixel tests; prefer snapshot structures that capture tile IDs, frame indices, and positions.

## Performance and Allocation

- Avoid per-frame heap churn in tight loops.
- Use small structs and pass pointers only when mutation is required or copying is expensive.
- Profile with pprof before optimizing; optimize only bottlenecks.

## Code Review Checklist

- Readability: Is the code straightforward and consistent with these conventions?
- Safety: Are errors handled appropriately? Are temporary hard exits marked with TODO to soften?
- Determinism: Are time and randomness driven by the central clock/PRNG?
- Style: Are loops using `range` where reasonable? Are names idiomatic?
- Boundaries: Does core avoid importing renderer/input packages?
- Tests: Are new behaviors covered by tests or at least easy to test?
- Docs: Are exported symbols documented and comments accurate?

## Action/Command Patterns

Player actions (commands like Look, Push, Get, Klimb) should follow these patterns:

### GameState Action Methods
```go
// Small map actions - require direction parameter
func (g *GameState) Action[Command]SmallMap(direction references.Direction) bool

// Large map actions - require direction parameter  
func (g *GameState) Action[Command]LargeMap(direction references.Direction) bool
```

### Conventions:
- **Naming**: Always prefix with `Action`, use PascalCase command name, suffix with map type
- **Parameters**: Always include `direction references.Direction` parameter
- **Returns**: Return `bool` indicating success/failure
- **Imports**: Use direct import `"github.com/bradhannah/Ultima5ReduxGo/internal/references"`, avoid aliases
- **File organization**: Group related actions in `action_[command].go` files

### Examples:
- `ActionGetSmallMap(direction references.Direction) bool`
- `ActionLookLargeMap(direction references.Direction) bool`
- `ActionPushSmallMap(direction references.Direction) bool`

## Examples

Range vs index-based:
