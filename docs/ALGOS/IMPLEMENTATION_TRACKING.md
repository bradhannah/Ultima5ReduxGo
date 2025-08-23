# Implementation Tracking (Go Code vs. Pseudocode)

Status legend:
- Implemented: code path exists and is wired to runtime.
- Similarity: Identical (logic matches), Similar (core idea present; differences or TODOs), Dissimilar/Not Implemented (missing or diverges materially).

Columns:
- Feature: High‑level area or routine.
- Pseudocode Ref: File/section in docs/ALGOS.
- Code Ref: Go file/function(s).
- Implemented: Yes/No/Partial.
- Similarity: Identical/Similar/Dissimilar.
- Notes: Differences, TODOs, gaps.

## Core Systems

| Feature                        | Pseudocode Ref                                                                      | Code Ref                                               | Implemented | Similarity | Notes                                                                                                                  |
|--------------------------------|-------------------------------------------------------------------------------------|--------------------------------------------------------|-------------|------------|------------------------------------------------------------------------------------------------------------------------|
| Look (tile descriptions)       | [Commands.md → Look — Towns/Overworld](./Commands.md#look-—-townsoverworld)         | `internal/references/look.go` (LookReferences)         | Partial     | Similar    | Loads `LOOK` data and returns descriptions; special tiles (telescope, wells) and trace‑to‑sign logic not visible here. |
| Windows/Arrow Slit LoS         | [Fixtures.md → Rare Fixtures & Edge Cases](./Fixtures.md#rare-fixtures--edge-cases) | `internal/map_state/layered_map.go` (comments 170–171) | Partial     | Similar    | Notes treating windows as opaque unless adjacent; ensure missiles pass and LoS aligns with our table.                  |
| Light sources & vision         | [Environment.md → Light Sources & Vision](./Environment.md#light-sources--vision)   | `internal/map_state/lighting.go`                       | Partial     | Similar    | Torch radius and static light sources exist; tie‑ins to commands/spells not wired.                                     |
| Torch duration                 | [Environment.md → Torch Duration](./Environment.md#torch-duration)                  | `internal/map_state/lighting.go`                       | Partial     | Similar    | `LightTorch()`/`AdvanceTurn()` exist; no UI command to ignite/consume torches wired.                                   |
| RNG & INT saves                | [RNG.md](./RNG.md)                                                                  | —                                                      | No          | —          | No central RNG helpers nor `saveint` equivalents in Go tree.                                                           |
| Field expiration (fieldkill)   | [Combat_Effects.md → Field Expiration](./Combat_Effects.md#field-expiration)        | —                                                      | No          | —          | Missing.                                                                                                               |
| Aiming UI (plraim)             | [Combat_Effects.md → Aiming UI](./Combat_Effects.md#aiming-ui)                      | —                                                      | No          | —          | Missing.                                                                                                               |
| Diagnose post‑hit messaging    | [Combat_Effects.md → Diagnose](./Combat_Effects.md#diagnose)                        | —                                                      | No          | —          | Missing.                                                                                                               |
| Combat field effects (infield) | [Combat_Effects.md → Field Effects](./Combat_Effects.md#field-effects)              | —                                                      | No          | —          | Missing.                                                                                                               |
| Distance helpers               | [Combat_Core.md → Distance Helpers](./Combat_Core.md#distance-helpers)              | —                                                      | No          | —          | A* exists; combat distance helpers not present.                                                                        |

## Commands

Note: If a command is not documented in `docs/ALGOS/Commands.md`, review legacy sources under `OLD/` and capture pseudocode first (then update the Pseudocode Ref and tracking here). Examples: Ready (`OLD/ZSTATS.C`), Cast (`OLD/COMBAT.C`, `OLD/SUBS3.C`).

| Feature             | Pseudocode Ref                                                                     | Code Ref                                                                                             | Implemented | Similarity | Notes                                                                                                                                              |
|---------------------|------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------|-------------|------------|----------------------------------------------------------------------------------------------------------------------------------------------------|
| Jimmy Door (Small)  | [Commands.md → Jimmy](./Commands.md#jimmy)                                         | `cmd/ultimav/gamescene_input_smallmap.go:245` + `internal/game_state/action_jimmy_door.go`           | Partial     | Similar    | Doors only; 50% forced success; keys decrement on magic locks; no chest/stock/ladder variants; needs town constraints.                             |
| Open (Small, doors) | [Commands.md → Open — Towns/Overworld](./Commands.md#open-—-townsoverworld)        | `cmd/ultimav/gamescene_input_smallmap.go:212` + `internal/map_state/action_open_door.go`             | Partial     | Similar    | Door state machine with messages; timed doors/other surfaces TBD; includes LB treasure chest special‑case.                                         |
| Open (Large)        | [Commands.md → Open — Towns/Overworld](./Commands.md#open-—-townsoverworld)        | `cmd/ultimav/gamescene_input_largemap.go:64`                                                         | No          | —          | Outputs “Cannot”.                                                                                                                                  |
| Push (Small)        | [Commands.md → Push](./Commands.md#push)                                           | `cmd/ultimav/gamescene_input_smallmap.go:195` + `internal/game_state/action_push.go`                 | Partial     | Similar    | Handles chairs, cannons, generic push/swap; needs full matrix and blocked cases (e.g., walls/NPCs).                                                |
| Get (Small)         | [Commands.md → Get — Towns/Overworld](./Commands.md#get-—-townsoverworld)          | `cmd/ultimav/gamescene_input_smallmap.go:260` + `internal/game_state/action_get.go`                  | Partial     | Similar    | Picks up item stacks, sconces (torch), food with karma hit, crops; broader pickup/chest flows TBD.                                                 |
| Get (Large)         | [Commands.md → Get — Towns/Overworld](./Commands.md#get-—-townsoverworld)          | `cmd/ultimav/gamescene_input_largemap.go:33`                                                         | No          | —          | “Get what?” only.                                                                                                                                  |
| Ready               | [Commands.md → Ready](./Commands.md#ready)                                         | —                                                                                                    | No          | —          | Pseudocode added; not implemented. Enforce slots/weight/ammo; negatives: no armor change in combat, missing ammo, etc.                             |
| Talk (Small)        | [Commands.md → Talk](./Commands.md#talk-freed-npc-nuance)                          | `cmd/ultimav/gamescene_input_smallmap.go:326`                                                        | Partial     | Dissimilar | Uses linear dialog engine; TLK/merchant flows not integrated yet.                                                                                  |
| Talk (Large)        | [Commands.md → Talk](./Commands.md#talk-freed-npc-nuance)                          | `cmd/ultimav/gamescene_input_largemap.go:80`                                                         | No          | —          | “Talk to who?” only.                                                                                                                               |
| Klimb (Small)       | [Commands.md → Klimb](./Commands.md#klimb)                                         | `cmd/ultimav/gamescene_input_smallmap.go:155,189` + `internal/game_state/action_klimb.go`            | Partial     | Similar    | Ladders/grates up/down; directional climb passes fences only; mountain/grapple and dungeon flows not present.                                      |
| Klimb (Large)       | [Commands.md → Klimb](./Commands.md#klimb)                                         | `cmd/ultimav/gamescene_input_largemap.go:45`, `:101`                                                 | No          | —          | Prompt only; no action on secondary input.                                                                                                         |
| Look (Small/Large)  | [Commands.md → Look](./Commands.md#look)                                           | `cmd/ultimav/gamescene_input_common.go:20` + `cmd/ultimav/gamescene_input_largemap.go:106`           | Yes         | Similar    | Directional look with LookReferences; small map adds clock time on clocks; dungeon lighting constraints not applied here.                          |
| Pass Turn           | [Commands.md → Pass Turn (Space)](./Commands.md#pass-turn-space)                   | `cmd/ultimav/gamescene_input_smallmap.go:25,91` + `cmd/ultimav/gamescene_input_largemap.go:20,61,88` | Yes         | Similar    | Adds “Pass” and calls `FinishTurn()` when in PrimaryInput state.                                                                                   |
| Hole Up & Camp      | [Commands.md → Hole Up & Camp](./Commands.md#hole-up--camp)                        | —                                                                                                    | No          | —          | Not implemented. Includes ship repair (anchored/furled), OW camping (`outcamp`), dungeon camp (`dngcamp`/`vcombat`), guard/watch, time/food/regen. |
| Board               | [Commands.md → Board](./Commands.md#board)                                         | `cmd/ultimav/gamescene_actions.go:8`                                                                 | Partial     | Similar    | Boards vehicle at position; messages vary; relies on vehicle presence; ship/skiff nuances TBD.                                                     |
| Exit                | [Commands.md → Exit](./Commands.md#exit-leave-buildingtown)                        | `cmd/ultimav/gamescene_actions.go:29`                                                                | Partial     | Similar    | Exits vehicle; messages wired; broader map exiting (buildings/towns) handled via Enter/Exit actions elsewhere.                                     |
| Enter (Large)       | [Commands.md → Enter](./Commands.md#enter)                                         | `cmd/ultimav/gamescene_input_largemap.go:52` + `internal/game_state/action_enter.go`                 | Partial     | Similar    | Enters building when on a world location; small‑map Enter not wired.                                                                               |
| Enter (Small)       | [Commands.md → Enter](./Commands.md#enter)                                         | `cmd/ultimav/gamescene_input_smallmap.go:59`                                                         | No          | —          | Negative prompt only: prints “Enter what?”.                                                                                                        |
| Ignite Torch        | [Commands.md → Ignite Torch](./Commands.md#ignite-torch)                           | `cmd/ultimav/gamescene_input_*map.go` + `internal/game_state/action_ignite.go`                       | Yes         | Similar    | Decrements torches and lights torch; dungeon/visibility interactions elsewhere. Negative: prints “None owned!” if zero.                            |
| View (Gem Map)      | [Commands.md → View (Gem Map)](./Commands.md#view-gem-map)                         | —                                                                                                    | No          | —          | Not implemented.                                                                                                                                   |
| Ztats               | [Commands.md → Ztats (Party Member Stats)](./Commands.md#ztats-party-member-stats) | —                                                                                                    | No          | —          | Not implemented.                                                                                                                                   |
| Mix Reagents        | [Commands.md → Mix Reagents](./Commands.md#mix-reagents)                           | —                                                                                                    | No          | —          | Not implemented.                                                                                                                                   |
| Use                 | [Commands.md → Use](./Commands.md#use)                                             | —                                                                                                    | No          | —          | Not implemented. Should route to item/fixture context; negative: “Nothing happens.” on unsupported tiles.                                          |
| Attack              | [Commands.md → Attack](./Commands.md#attack)                                       | —                                                                                                    | No          | —          | Not implemented. Surface/town initiation nuances; negative: deny where inappropriate.                                                              |
| Fire                | [Commands.md → Fire — Town/Ship](./Commands.md#fire-cannons)                       | —                                                                                                    | No          | —          | Not implemented. Large map negative should be “Fire broadsides only!” or “What?/Cannot” contextually.                                              |
| Cast                | [Commands.md → Cast](./Commands.md#cast), Spells.md                                | —                                                                                                    | No          | —          | Pseudocode added; not implemented. Context gating (Not here!/Absorbed!), stock/MP checks, success reporting, turn rules.                           |
| New Order           | [Commands.md → New Order](./Commands.md#new-order-swap-party-positions)            | —                                                                                                    | No          | —          | Not implemented. Negative handling per legacy when invalid.                                                                                        |
| Fire (Cannons)      | [Commands.md → Fire — Town/Ship](./Commands.md#fire-cannons)                       | —                                                                                                    | No          | —          | Not implemented.                                                                                                                                   |
| Search              | [Commands.md → Search](./Commands.md#search)                                       | —                                                                                                    | No          | —          | Not implemented.                                                                                                                                   |

### Dungeon Commands

Dungeon interactions often differ from surface/town flows (ahead‑of‑avatar targeting, dungeon tile families, light checks, underfoot chests, etc.).

| Feature                  | Pseudocode Ref                                                         | Code Ref                                              | Implemented | Similarity | Notes                                                                                                            |
|--------------------------|------------------------------------------------------------------------|-------------------------------------------------------|-------------|------------|------------------------------------------------------------------------------------------------------------------|
| Look — Dungeon           | [Commands.md → Look — Dungeon](./Commands.md#look-—-dungeon)           | —                                                     | No          | —          | Dungeon look requires light and prints tile‑family descriptions; includes fountain drink prompt and field types. |
| Open — Dungeon           | [Commands.md → Open — Dungeon](./Commands.md#open-—-dungeon)           | —                                                     | No          | —          | Handles dungeon doors and underfoot chests; integrates with spells (An Sanct/In Ex Por).                         |
| Jimmy — Dungeon          | [Commands.md → Jimmy](./Commands.md#jimmy)                             | —                                                     | No          | —          | Dungeon variant targets ahead tile and supports dungeon chest jimmy odds.                                        |
| Get — Dungeon            | [Commands.md → Get — Dungeon](./Commands.md#get-—-dungeon)             | —                                                     | No          | —          | Picks from underfoot opened chest; distinct from surface object pickup.                                          |
| Search — Dungeon (Ahead) | [Commands.md → Search — Dungeon](./Commands.md#search-—-dungeon-ahead) | —                                                     | No          | —          | Ahead‑of‑avatar search for secret doors/passages; separate flow.                                                 |
| View (Gem Map) — Dungeon | [Commands.md → View (Gem Map)](./Commands.md#view-gem-map)             | —                                                     | No          | —          | Renders dungeon level cell layout; consumes gem.                                                                 |
| Pass Turn — Dungeon Tick | [Commands.md → Pass Turn (Space)](./Commands.md#pass-turn-space)       | `cmd/ultimav/gamescene_input.go` (no Dungeon handler) | No          | —          | Should advance dungeon hazards/lighting per tick.                                                                |

## Fixtures & Environment

| Feature                   | Pseudocode Ref                                                           | Code Ref | Implemented | Similarity | Notes                                         |
|---------------------------|--------------------------------------------------------------------------|----------|-------------|------------|-----------------------------------------------|
| Fixtures: default mapping | [Fixtures.md → Fixture Defaults](./Fixtures.md#fixture-defaults-mapping) | —        | No          | —          | No FixtureEffects map or Use routing present. |
| Wells: Wish               | [Fixtures.md → Wells & Wish](./Fixtures.md#wells--wish)                  | —        | No          | —          | Missing.                                      |
| Fountains                 | [Fixtures.md → Fountains](./Fixtures.md#fountains)                       | —        | No          | —          | Missing.                                      |
| Lamps/Sconces overrides   | [Fixtures.md → Lamps/Sconces](./Fixtures.md#lampssconces-overrides)      | —        | No          | —          | Missing.                                      |
| Overworld hazards         | [Environment.md](./Environment.md)                                       | —        | No          | —          | Missing.                                      |
| Moongates                 | [Moongates.md](./Moongates.md)                                           | —        | No          | —          | Missing.                                      |
| Town drawbridges          | [Towns.md → Drawbridges](./Towns.md#drawbridges)                         | —        | No          | —          | Missing.                                      |

## Schedules & AI

| Feature                           | Pseudocode Ref                                                                 | Code Ref                                                 | Implemented | Similarity | Notes                                                                      |
|-----------------------------------|--------------------------------------------------------------------------------|----------------------------------------------------------|-------------|------------|----------------------------------------------------------------------------|
| NPC schedules (data/model)        | [NPC_Schedules.md → Data Model](./NPC_Schedules.md#data-model)                 | `internal/references/npc_schedule.go`                    | Yes         | Similar    | Schedule model present; details may differ.                                |
| NPC schedule driver (hour change) | [NPC_Schedules.md → Hourly Transitions](./NPC_Schedules.md#hourly-transitions) | `internal/ai/npc_ai_controller_small_map.go` (various)   | Partial     | Similar    | Controller selects behaviors and floors; exact LEAV/ARIV/POP not verbatim. |
| Small map pathfinding             | [NPC_Schedules.md → Pathfinding](./NPC_Schedules.md#pathfinding)               | `internal/astar/*.go`, `internal/ai/npc_ai_controller_*` | Yes         | Similar    | Pathfinding exists; integration with schedules ongoing.                    |
| Combat AI (seek, special moves)   | [Movement_Combat_AI.md](./Movement_Combat_AI.md)                               | —                                                        | No          | —          | Combat not implemented.                                                    |
| Mass charm targeting ('C')        | [Spells.md → Quas An Wis](./Spells.md#quas-an-wis-mass-charmconfusion)         | —                                                        | No          | —          | Not applicable yet.                                                        |

## Spells & Scrolls

| Feature                 | Pseudocode Ref                                            | Code Ref                                        | Implemented | Similarity | Notes                                                |
|-------------------------|-----------------------------------------------------------|-------------------------------------------------|-------------|------------|------------------------------------------------------|
| Spellcasting core       | [Spells.md](./Spells.md)                                  | —                                               | No          | —          | Not present in Go codebase.                          |
| Specific spells/scrolls | [Spells.md (All)](./Spells.md#spells-summary-at-a-glance) | `internal/references/data/InventoryDetails.csv` | No          | —          | Data present for names/info; no runtime casting/use. |

### Spell Inventory Data (FYI)

| Feature                     | Pseudocode Ref                                      | Code Ref                                            | Implemented | Similarity | Notes                                |
|-----------------------------|-----------------------------------------------------|-----------------------------------------------------|-------------|------------|--------------------------------------|
| Spell metadata (names/info) | [Spells.md](./Spells.md)                            | `internal/references/data/InventoryDetails.csv`     | Yes         | N/A        | Data present; no runtime casting.    |
| Inventory quantities        | [SAVED_GAM_STRUCTURE.md](../SAVED_GAM_STRUCTURE.md) | `internal/party_state/inventory.go` (Scrolls, etc.) | Partial     | N/A        | Data structures exist; no use flows. |

### Potions & Scrolls

| Feature        | Pseudocode Ref                                                         | Code Ref | Implemented | Similarity | Notes            |
|----------------|------------------------------------------------------------------------|----------|-------------|------------|------------------|
| Potion effects | [Potions.md](./Potions.md)                                             | —        | No          | —          | Not implemented. |
| Scroll effects | [Spells.md → Scrolls Summary](./Spells.md#scrolls-summary-at-a-glance) | —        | No          | —          | Not implemented. |

## Special Items & Artifacts

| Feature/Item       | Pseudocode Ref                                                                     | Code Ref                                          | Implemented | Similarity | Notes                                          |
|--------------------|------------------------------------------------------------------------------------|---------------------------------------------------|-------------|------------|------------------------------------------------|
| Magic Carpet (Use) | [Commands.md → Use](./Commands.md#use)                                             | `internal/map_units/npc_vehicle.go` (carpet type) | Partial     | Dissimilar | Carpet NPC exists; use/place/pickup logic TBD. |
| Skull Keys (Use)   | [Objects.md → Skull Key — Magical Unlock](./Objects.md#skull-key-—-magical-unlock) | `internal/party_state/inventory.go` (SkullKeys)   | No          | —          | Inventory tracked; no use flow.                |
| Crown (Use)        | [Commands.md → Use](./Commands.md#use)                                             | —                                                 | No          | —          | Not implemented.                               |
| Sceptre (Use)      | [Commands.md → Use](./Commands.md#use)                                             | —                                                 | No          | —          | Not implemented.                               |
| Amulet (Use)       | [Commands.md → Use](./Commands.md#use)                                             | `internal/references/item_equipment.go`           | No          | —          | Item enum exists; effect/use not wired.        |
| Spyglass/Telescope | [Fixtures.md → Telescope](./Fixtures.md#telescope)                                 | —                                                 | No          | —          | Not implemented.                               |
| Gems (View)        | [Commands.md → View (Gem Map)](./Commands.md#view-gem-map)                         | —                                                 | No          | —          | Not implemented.                               |
| Torches (Ignite)   | [Commands.md → Ignite Torch](./Commands.md#ignite-torch)                           | `internal/map_state/lighting.go`                  | Partial     | Similar    | Lighting supports torches; command missing.    |

## Town Systems

| Feature                        | Pseudocode Ref                                                         | Code Ref | Implemented | Similarity | Notes                        |
|--------------------------------|------------------------------------------------------------------------|----------|-------------|------------|------------------------------|
| Guard alarm/pursuit            | [Towns.md → Special Guard Behavior](./Towns.md#special-guard-behavior) | —        | No          | —          | Not implemented.             |
| Jail flow                      | [Towns.md → Jail Flow](./Towns.md#jail-flow)                           | —        | No          | —          | Not implemented.             |
| Cannons (town/ship broadsides) | [Commands.md → Fire](./Commands.md#fire-cannons)                       | —        | No          | —          | Not implemented.             |
| Shops (pricing/services)       | [Shops.md](./Shops.md)                                                 | —        | No          | —          | Pricing tables to be filled. |

## Potions & Scrolls

| Feature       | Pseudocode Ref      | Code Ref                                    | Implemented | Similarity | Notes                                     |
|---------------|---------------------|---------------------------------------------|-------------|------------|-------------------------------------------|
| Potions (use) | Potions.md          | `internal/party_state/inventory.go` (types) | No          | —          | Quantities exist but no use/effect logic. |
| Scrolls (use) | Spells.md (scrolls) | `internal/party_state/inventory.go` (types) | No          | —          | No scroll use flows implemented.          |

## Special Items

| Feature            | Pseudocode Ref          | Code Ref                                                     | Implemented | Similarity | Notes                                                   |
|--------------------|-------------------------|--------------------------------------------------------------|-------------|------------|---------------------------------------------------------|
| Crown (Use)        | Commands.md → Use       | —                                                            | No          | —          | Not implemented.                                        |
| Sceptre (Use)      | Commands.md → Use       | —                                                            | No          | —          | Not implemented.                                        |
| Amulet (Use)       | Commands.md → Use       | `internal/party_state/types.go` (Amulet field)               | No          | —          | Field exists; no use effect logic.                      |
| Carpet (Board/Use) | Commands.md → Use/Board | `internal/map_units/npc_vehicle.go` (CarpetVehicle), map use | Partial     | Dissimilar | Vehicle types exist; no boarding/Use flows as per docs. |
| Spyglass/Telescope | Commands.md → Use/Look  | —                                                            | No          | —          | Not implemented.                                        |
| Gems (View)        | Commands.md → View      | `internal/party_state/inventory.go` (gems qty)               | No          | —          | Quantities exist; View command absent.                  |

## Exhaustive Spell Checklist (48 Spells)

Legend: Implemented = No (unless otherwise noted), Similarity = —, Code Ref column lists any related data structures.

| #  | Spell            | Pseudocode Ref | Code Ref                                        | Implemented | Similarity | Notes                       |
|----|------------------|----------------|-------------------------------------------------|-------------|------------|-----------------------------|
| 0  | In Lor           | Spells.md      | `internal/references/data/InventoryDetails.csv` | No          | —          | —                           |
| 1  | Grav Por         | Spells.md      | same                                            | No          | —          | —                           |
| 2  | An Zu            | Spells.md      | same                                            | No          | —          | —                           |
| 3  | An Nox           | Spells.md      | same                                            | No          | —          | —                           |
| 4  | Mani             | Spells.md      | same                                            | No          | —          | —                           |
| 5  | An Ylem          | Spells.md      | same                                            | No          | —          | —                           |
| 6  | An Sanct         | Spells.md      | same                                            | No          | —          | —                           |
| 7  | An Xen Corp      | Spells.md      | same                                            | No          | —          | —                           |
| 8  | Rel Hur          | Spells.md      | same                                            | No          | —          | —                           |
| 9  | In Wis           | Spells.md      | same                                            | No          | —          | —                           |
| 10 | Kal Xen          | Spells.md      | same                                            | No          | —          | —                           |
| 11 | In Xen Mani      | Spells.md      | same                                            | No          | —          | —                           |
| 12 | Vas Lor          | Spells.md      | same                                            | No          | —          | —                           |
| 13 | Vas Flam         | Spells.md      | same                                            | No          | —          | —                           |
| 14 | In Flam Grav     | Spells.md      | same                                            | No          | —          | —                           |
| 15 | In Nox Grav      | Spells.md      | same                                            | No          | —          | —                           |
| 16 | In Zu Grav       | Spells.md      | same                                            | No          | —          | —                           |
| 17 | In Por           | Spells.md      | same                                            | No          | —          | —                           |
| 18 | An Grav          | Spells.md      | same                                            | No          | —          | —                           |
| 19 | In Sanct         | Spells.md      | same                                            | No          | —          | —                           |
| 20 | In Sanct G       | Spells.md      | same                                            | No          | —          | Energy field create variant |
| 21 | Uus Por          | Spells.md      | same                                            | No          | —          | —                           |
| 22 | Des Por          | Spells.md      | same                                            | No          | —          | —                           |
| 23 | Wis Quas         | Spells.md      | same                                            | No          | —          | —                           |
| 24 | In Bet Xen       | Spells.md      | same                                            | No          | —          | —                           |
| 25 | An Ex Por        | Spells.md      | same                                            | No          | —          | —                           |
| 26 | In Ex Por        | Spells.md      | same                                            | No          | —          | —                           |
| 27 | Vas Mani         | Spells.md      | same                                            | No          | —          | —                           |
| 28 | In Zu            | Spells.md      | same                                            | No          | —          | —                           |
| 29 | Rel Tym          | Spells.md      | same                                            | No          | —          | —                           |
| 30 | In Vas Por Ylem  | Spells.md      | same                                            | No          | —          | —                           |
| 31 | Quas An Wis      | Spells.md      | same                                            | No          | —          | Mass charm aura             |
| 32 | In An            | Spells.md      | same                                            | No          | —          | Negate magic                |
| 33 | Wis An Ylem      | Spells.md      | same                                            | No          | —          | X‑Ray                       |
| 34 | An Xen Ex        | Spells.md      | same                                            | No          | —          | —                           |
| 35 | Rel Xen Bet      | Spells.md      | same                                            | No          | —          | Polymorph                   |
| 36 | Sanct Lor        | Spells.md      | same                                            | No          | —          | Invisibility                |
| 37 | Xen Corp         | Spells.md      | same                                            | No          | —          | —                           |
| 38 | In Quas Xen      | Spells.md      | same                                            | No          | —          | Clone                       |
| 39 | In Quas Wis      | Spells.md      | same                                            | No          | —          | View                        |
| 40 | In Nox Hur       | Spells.md      | same                                            | No          | —          | Poison storm                |
| 41 | In Quas Corp     | Spells.md      | same                                            | No          | —          | Fear                        |
| 42 | In Mani Corp     | Spells.md      | same                                            | No          | —          | Resurrection                |
| 43 | Kal Xen Corp     | Spells.md      | same                                            | No          | —          | Summon daemon               |
| 44 | In Vas Grav Corp | Spells.md      | same                                            | No          | —          | Energy storm                |
| 45 | In Flam Hur      | Spells.md      | same                                            | No          | —          | Firestorm                   |
| 46 | Vas Rel Por      | Spells.md      | same                                            | No          | —          | Gate Travel                 |
| 47 | An Tym           | Spells.md      | same                                            | No          | —          | Negate Time                 |
| 48 | Frotz (reserved) | Spells.md      | same                                            | No          | —          | Unimplemented               |

## Scrolls Checklist (8)

| Scroll        | Pseudocode Ref | Code Ref                                    | Implemented | Similarity | Notes |
|---------------|----------------|---------------------------------------------|-------------|------------|-------|
| Light         | Spells.md      | `internal/party_state/inventory.go` (types) | No          | —          | —     |
| Wind Change   | Spells.md      | same                                        | No          | —          | —     |
| Protection    | Spells.md      | same                                        | No          | —          | —     |
| Negate Magic  | Spells.md      | same                                        | No          | —          | —     |
| View          | Spells.md      | same                                        | No          | —          | —     |
| Summon Daemon | Spells.md      | same                                        | No          | —          | —     |
| Resurrection  | Spells.md      | same                                        | No          | —          | —     |
| Negate Time   | Spells.md      | same                                        | No          | —          | —     |

## Potions Checklist (8)

| Potion Color | Pseudocode Ref | Code Ref                                    | Implemented | Similarity | Notes              |
|--------------|----------------|---------------------------------------------|-------------|------------|--------------------|
| Blue         | Potions.md     | `internal/party_state/inventory.go` (types) | No          | —          | Cure Sleep         |
| Yellow       | Potions.md     | same                                        | No          | —          | Heal               |
| Red          | Potions.md     | same                                        | No          | —          | Cure Poison        |
| Green        | Potions.md     | same                                        | No          | —          | Poison             |
| Orange       | Potions.md     | same                                        | No          | —          | Sleep              |
| Purple       | Potions.md     | same                                        | No          | —          | Polymorph (combat) |
| Black        | Potions.md     | same                                        | No          | —          | Invisible (combat) |
| White        | Potions.md     | same                                        | No          | —          | X‑Ray (surface)    |

## Special Items Checklist

| Item                  | Pseudocode Ref       | Code Ref                                            | Implemented | Similarity | Notes                                      |
|-----------------------|----------------------|-----------------------------------------------------|-------------|------------|--------------------------------------------|
| Crown of Lord British | Commands.md → Use    | —                                                   | No          | —          | —                                          |
| Sceptre of Lord Brit. | Commands.md → Use    | —                                                   | No          | —          | —                                          |
| Amulet of LB          | Commands.md → Use    | `internal/party_state/types.go` (Amulet field)      | No          | —          | Field present only                         |
| Magic Carpet          | Commands.md → Use    | `internal/map_units/npc_vehicle.go` (CarpetVehicle) | Partial     | Dissimilar | Vehicle type present; no Use/boarding flow |
| Skull Keys            | Commands.md → Use    | —                                                   | No          | —          | —                                          |
| Keys                  | Commands.md → Open   | `internal/party_state/inventory.go` (keys qty)      | No          | —          | No Open/door flows                         |
| Torches               | Commands.md → Ignite | `internal/party_state/inventory.go` (torches qty)   | No          | —          | No Ignite Torch command                    |
| Gems                  | Commands.md → View   | `internal/party_state/inventory.go` (gems qty)      | No          | —          | View not implemented                       |
| Spyglass              | Commands.md → Use    | —                                                   | No          | —          | —                                          |
| Telescope             | Commands.md → Look   | —                                                   | No          | —          | —                                          |

## Towns & Special Systems

| Feature                | Pseudocode Ref                      | Code Ref | Implemented | Similarity | Notes      |
|------------------------|-------------------------------------|----------|-------------|------------|------------|
| Drawbridges/Portcullis | Towns.md                            | —        | No          | —          | Not found. |
| Guard alarm & Jail     | Towns.md → Guard Behavior/Jail      | —        | No          | —          | Not found. |
| Cannons (town fire)    | Combat_Effects.md/Towns.md          | —        | No          | —          | Not found. |
| Bridge trolls          | Special_BridgeTrolls.md             | —        | No          | —          | Not found. |
| Wind system            | Movement_Overworld.md → Wind System | —        | No          | —          | Not found. |
| Ships & Sails          | Commands.md / Movement_Overworld.md | —        | No          | —          | Not found. |
| Moongates              | Moongates.md                        | —        | No          | —          | Not found. |

## Shops & Economy

| Feature                    | Pseudocode Ref       | Code Ref                                                            | Implemented | Similarity | Notes                                                                                                |
|----------------------------|----------------------|---------------------------------------------------------------------|-------------|------------|------------------------------------------------------------------------------------------------------|
| Shop pricing & multipliers | Shops.md             | —                                                                   | No          | —          | Tables added in docs; not implemented in code.                                                       |
| Reagent/Healer/Arms Shops  | Shops.md             | —                                                                   | No          | —          | Not present.                                                                                         |
| Inns (stay months)         | Shops.md → Innkeeper | `internal/party_state/types.go` (MonthsAtInn, PartyStatus=AtTheInn) | No          | —          | Data fields present; no inn UI/flow, no gold deduction, no time advancement, no `inn_party` mapping. |
| Horses/Shipwright          | Shops.md             | —                                                                   | No          | —          | Not present.                                                                                         |

## Conversation System (FYI)

| Feature              | Pseudocode Ref           | Code Ref                                          | Implemented | Similarity | Notes                                        |
|----------------------|--------------------------|---------------------------------------------------|-------------|------------|----------------------------------------------|
| Talk engine (linear) | TALK_SYSTEM_STRUCTURE.md | `internal/conversation/linear_engine.go` (+tests) | Yes         | Similar    | Robust implementation; outside combat scope. |

---

This tracker is a starting point. As features land in Go code, update the “Implemented” and “Similarity” columns, and add concrete function references.

## Sleep & Rest (Overview)

| Feature               | Pseudocode Ref                                                                                   | Code Ref                                         | Implemented | Similarity | Notes                                                                                                            |
|-----------------------|--------------------------------------------------------------------------------------------------|--------------------------------------------------|-------------|------------|------------------------------------------------------------------------------------------------------------------|
| Sleep status effect   | Combat_Effects.md → Field Effects, Per‑Turn Updates; Potions.md (Blue/Orange); Spells.md (In Zu) | `internal/party_state/mappings.go` (Status enum) | No          | —          | Status enum includes Sleep; no application/expiration logic, no wake/cure handling wired.                        |
| Sleep field (tiles)   | Combat_Effects.md → Field Effects                                                                | —                                                | No          | —          | Combat field effects not implemented; would apply Sleep on contact (PCs/monsters).                               |
| In Zu (Sleep spell)   | Spells.md                                                                                        | —                                                | No          | —          | Listed in spells table; no casting/effect pipeline.                                                              |
| Potions: Blue/Orange  | Potions.md                                                                                       | `internal/party_state/inventory.go` (types only) | No          | —          | Blue cures Sleep; Orange applies Sleep; items exist as types only.                                               |
| Hole Up & Camp (rest) | Commands.md → Hole Up & Camp                                                                     | —                                                | No          | —          | Camping/repair flows absent; should handle guard watch, time advance, HP/MP regen, food ticks, encounter checks. |
| In‑bed sleep/ejection | Fixtures.md → Beds; NPC_Schedules.md (eject sleepers at schedule boundaries)                     | —                                                | No          | —          | No in‑bed sleep flow; must prevent sleeping in occupied beds and eject sleepers around hour changes.             |
