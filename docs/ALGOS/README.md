# Game Algorithms Index

This folder splits the algorithms and odds documentation into logical topics with self-contained pseudocode blocks.

Use this index to find systems quickly. Each file contains focused sections and links back to adjacent logic where helpful.

## Index

- RNG and Saves: `RNG.md`
- Encounters and Spawning: `Encounters.md`
- Movement (Overworld): `Movement_Overworld.md`
- Movement (Combat AI): `Movement_Combat_AI.md`
- NPC Schedules: `NPC_Schedules.md`
- Combat Core (Hit, Damage, Turn Modifiers): `Combat_Core.md`
- Combat Effects (Poison, Flee, Fields, Missiles, Regurgitation, Per-Turn): `Combat_Effects.md`
- Spells: `Spells.md` (Rel Hur, In/Great Light, Rel Tym/Negate Time, Uus/Des Por, An Sanct/An Grav, In/An Ex Por, summons)
- Potions: `Potions.md`
- Environment and Hazards: `Environment.md`
- Town Systems (Guards, Cannons, Karma, Drawbridges): `Towns.md`
- Objects and Fixtures (mirrors, pushables, doors, seating, Klimb): `Objects.md`
- Small-Town Fixtures (fountains, wells, lamps, telescope): `Fixtures.md`
- Hidden Objects and Search Mechanics (daily skull keys, conditional spawns): `Secrets.md`
- Dungeon Interactions (chests, traps, search): `Dungeon.md`
- Commands (Yell, Look, Open, Push, Get, Talk, Use, Fire, Hole Up & Camp, Exit, Klimb, Pass Turn, Search, View, Ztats): `Commands.md` ⚠️ _See Implementation Notes above for Go remake patterns_
- Special: Bridge Trolls: `Special_BridgeTrolls.md`
- Moongates: `Moongates.md`
- Shops and Inns: `Shops.md`
- Task Tracker: `TASKS.md`

## Conventions

- Pseudocode blocks are fenced with ```pseudocode and aim for clarity over literal code.
- Odds are expressed both in fractions and comments inline with RNG calls.
- Time and RNG must be driven by the central game clock and PRNG for determinism.

## Implementation Notes for Go Remake

**Map Type Separation**: Although the pseudocode MAY include map type checks (e.g., `if is_combat_map()`), the Go implementation uses separate action methods for each map type for clarity:

- `ActionCommandSmallMap(direction)` - Small maps (towns, buildings)
- `ActionCommandLargeMap(direction)` - Large maps (overworld) 
- `ActionCommandCombatMap(direction)` - Combat maps
- `ActionCommandDungeon(direction)` - Dungeon maps

**Avoid Map Type Checks**: Within each action method, avoid `IsInCombat()`, `IsOverworld()`, etc. checks since the method dispatch already handles map type separation. This improves code clarity and reduces coupling.

**Common Helper Functions**: Use shared helper functions for common logic elements when it doesn't negatively affect readability (e.g., `IsPushable()`, `ObjectPresentAt()`).
