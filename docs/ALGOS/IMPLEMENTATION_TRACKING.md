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

| Implemented | Feature                        | Pseudocode Ref                                                                      | Code Ref                                               | Similarity | Notes                                                                                                                  |
|-------------|--------------------------------|-------------------------------------------------------------------------------------|--------------------------------------------------------|------------|------------------------------------------------------------------------------------------------------------------------|
| Partial     | Look (tile descriptions)       | [Commands.md → Look — Towns/Overworld](./Commands.md#look-—-townsoverworld)         | `internal/references/look.go` (LookReferences)         | Similar    | Loads `LOOK` data and returns descriptions; special tiles (telescope, wells) and trace‑to‑sign logic not visible here. |
| Partial     | Windows/Arrow Slit LoS         | [Fixtures.md → Rare Fixtures & Edge Cases](./Fixtures.md#rare-fixtures--edge-cases) | `internal/map_state/layered_map.go` (comments 170–171) | Similar    | Notes treating windows as opaque unless adjacent; ensure missiles pass and LoS aligns with our table.                  |
| Partial     | Light sources & vision         | [Environment.md → Light Sources & Vision](./Environment.md#light-sources--vision)   | `internal/map_state/lighting.go`                       | Similar    | Torch radius and static light sources exist; tie‑ins to commands/spells not wired.                                     |
| Partial     | Torch duration                 | [Environment.md → Torch Duration](./Environment.md#torch-duration)                  | `internal/map_state/lighting.go`                       | Similar    | `LightTorch()`/`AdvanceTurn()` exist; no UI command to ignite/consume torches wired.                                   |
| No          | RNG & INT saves                | [RNG.md](./RNG.md)                                                                  | —                                                      | —          | No central RNG helpers nor `saveint` equivalents in Go tree.                                                           |
| No          | Field expiration (fieldkill)   | [Combat_Effects.md → Field Expiration](./Combat_Effects.md#field-expiration)        | —                                                      | —          | Missing.                                                                                                               |
| No          | Aiming UI (plraim)             | [Combat_Effects.md → Aiming UI](./Combat_Effects.md#aiming-ui)                      | —                                                      | —          | Missing.                                                                                                               |
| No          | Diagnose post‑hit messaging    | [Combat_Effects.md → Diagnose](./Combat_Effects.md#diagnose)                        | —                                                      | —          | Missing.                                                                                                               |
| No          | Combat field effects (infield) | [Combat_Effects.md → Field Effects](./Combat_Effects.md#field-effects)              | —                                                      | —          | Missing.                                                                                                               |
| No          | Distance helpers               | [Combat_Core.md → Distance Helpers](./Combat_Core.md#distance-helpers)              | —                                                      | —          | A* exists; combat distance helpers not present.                                                                        |

## Commands

Note: If a command is not documented in `docs/ALGOS/Commands.md`, review legacy sources under `OLD/` and capture pseudocode first (then update the Pseudocode Ref and tracking here). Examples: Ready (`OLD/ZSTATS.C`), Cast (`OLD/COMBAT.C`, `OLD/SUBS3.C`).

| Implemented | Feature             | Pseudocode Ref                                                                     | Code Ref                                                                                             | Similarity | Notes                                                                                                                                                                  |
|-------------|---------------------|------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------|------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Yes         | Jimmy Door (Small)  | [Commands.md → Jimmy](./Commands.md#jimmy)                                         | `cmd/ultimav/gamescene_input_smallmap.go:285` + `internal/game_state/action_jimmy_door.go`           | Similar    | Fixed key consumption logic, added chest support, proper error messages; dungeon ahead-targeting still TODO.                                                           |
| Yes         | Jimmy Door (Large)  | [Commands.md → Jimmy](./Commands.md#jimmy)                                         | `cmd/ultimav/gamescene_input_largemap.go:116` + `internal/game_state/action_jimmy.go`                | Similar    | Uses same logic as small map; directional input implemented; character selection placeholder added.                                                                    |
| Partial     | Open (Small, doors) | [Commands.md → Open — Towns/Overworld](./Commands.md#open-—-townsoverworld)        | `cmd/ultimav/gamescene_input_smallmap.go:212` + `internal/map_state/action_open_door.go`             | Similar    | Door state machine with messages; timed doors/other surfaces TBD; includes LB treasure chest special‑case.                                                             |
| Yes         | Open (Large)        | [Commands.md → Open — Towns/Overworld](./Commands.md#open-—-townsoverworld)        | `cmd/ultimav/gamescene_input_largemap.go:68` + `internal/game_state/action_open.go:19-55`            | Similar    | Outputs "Open what?" immediately (no directional input). ActionOpenLargeMap handles doors and basic openable items; most interactions on large maps use Enter instead. |
| Partial     | Push (Small)        | [Commands.md → Push](./Commands.md#push)                                           | `cmd/ultimav/gamescene_input_smallmap.go:195` + `internal/game_state/action_push.go`                 | Similar    | Handles chairs, cannons, generic push/swap; needs full matrix and blocked cases (e.g., walls/NPCs).                                                                    |
| Partial     | Get (Small)         | [Commands.md → Get — Towns/Overworld](./Commands.md#get-—-townsoverworld)          | `cmd/ultimav/gamescene_input_smallmap.go:260` + `internal/game_state/action_get.go`                  | Similar    | Picks up item stacks, sconces (torch), food with karma hit, crops; broader pickup/chest flows TBD.                                                                     |
| No          | Get (Large)         | [Commands.md → Get — Towns/Overworld](./Commands.md#get-—-townsoverworld)          | `cmd/ultimav/gamescene_input_largemap.go:33`                                                         | —          | “Get what?” only.                                                                                                                                                      |
| No          | Ready               | [Commands.md → Ready](./Commands.md#ready)                                         | —                                                                                                    | —          | Pseudocode added; not implemented. Enforce slots/weight/ammo; negatives: no armor change in combat, missing ammo, etc.                                                 |
| Partial     | Talk (Small)        | [Commands.md → Talk](./Commands.md#talk-freed-npc-nuance)                          | `cmd/ultimav/gamescene_input_smallmap.go:326`                                                        | Dissimilar | Uses linear dialog engine; TLK/merchant flows not integrated yet.                                                                                                      |
| No          | Talk (Large)        | [Commands.md → Talk](./Commands.md#talk-freed-npc-nuance)                          | `cmd/ultimav/gamescene_input_largemap.go:80`                                                         | —          | “Talk to who?” only.                                                                                                                                                   |
| Partial     | Klimb (Small)       | [Commands.md → Klimb](./Commands.md#klimb)                                         | `cmd/ultimav/gamescene_input_smallmap.go:155,189` + `internal/game_state/action_klimb.go`            | Similar    | Ladders/grates up/down; directional climb passes fences only; mountain/grapple and dungeon flows not present.                                                          |
| No          | Klimb (Large)       | [Commands.md → Klimb](./Commands.md#klimb)                                         | `cmd/ultimav/gamescene_input_largemap.go:45`, `:101`                                                 | —          | Prompt only; no action on secondary input.                                                                                                                             |
| Yes         | Look (Small/Large)  | [Commands.md → Look](./Commands.md#look)                                           | `cmd/ultimav/gamescene_input_common.go:20` + `cmd/ultimav/gamescene_input_largemap.go:106`           | Similar    | Directional look with LookReferences; small map adds clock time on clocks; dungeon lighting constraints not applied here.                                              |
| Yes         | Pass Turn           | [Commands.md → Pass Turn (Space)](./Commands.md#pass-turn-space)                   | `cmd/ultimav/gamescene_input_smallmap.go:25,91` + `cmd/ultimav/gamescene_input_largemap.go:20,61,88` | Similar    | Adds “Pass” and calls `FinishTurn()` when in PrimaryInput state.                                                                                                       |
| No          | Hole Up & Camp      | [Commands.md → Hole Up & Camp](./Commands.md#hole-up--camp)                        | —                                                                                                    | —          | Not implemented. Includes ship repair (anchored/furled), OW camping (`outcamp`), dungeon camp (`dngcamp`/`vcombat`), guard/watch, time/food/regen.                     |
| Partial     | Board               | [Commands.md → Board](./Commands.md#board)                                         | `cmd/ultimav/gamescene_actions.go:8`                                                                 | Similar    | Boards vehicle at position; messages vary; relies on vehicle presence; ship/skiff nuances TBD.                                                                         |
| Partial     | Exit                | [Commands.md → Exit](./Commands.md#exit-leave-buildingtown)                        | `cmd/ultimav/gamescene_actions.go:29`                                                                | Similar    | Exits vehicle; messages wired; broader map exiting (buildings/towns) handled via Enter/Exit actions elsewhere.                                                         |
| Partial     | Enter (Large)       | [Commands.md → Enter](./Commands.md#enter)                                         | `cmd/ultimav/gamescene_input_largemap.go:52` + `internal/game_state/action_enter.go`                 | Similar    | Enters building when on a world location; small‑map Enter not wired.                                                                                                   |
| No          | Enter (Small)       | [Commands.md → Enter](./Commands.md#enter)                                         | `cmd/ultimav/gamescene_input_smallmap.go:59`                                                         | —          | Negative prompt only: prints “Enter what?”.                                                                                                                            |
| Yes         | Ignite Torch        | [Commands.md → Ignite Torch](./Commands.md#ignite-torch)                           | `cmd/ultimav/gamescene_input_*map.go` + `internal/game_state/action_ignite.go`                       | Similar    | Decrements torches and lights torch; dungeon/visibility interactions elsewhere. Negative: prints “None owned!” if zero.                                                |
| No          | View (Gem Map)      | [Commands.md → View (Gem Map)](./Commands.md#view-gem-map)                         | —                                                                                                    | —          | Not implemented.                                                                                                                                                       |
| No          | Ztats               | [Commands.md → Ztats (Party Member Stats)](./Commands.md#ztats-party-member-stats) | —                                                                                                    | —          | Not implemented.                                                                                                                                                       |
| No          | Mix Reagents        | [Commands.md → Mix Reagents](./Commands.md#mix-reagents)                           | —                                                                                                    | —          | Not implemented.                                                                                                                                                       |
| No          | Use                 | [Commands.md → Use](./Commands.md#use)                                             | —                                                                                                    | —          | Not implemented. Should route to item/fixture context; negative: “Nothing happens.” on unsupported tiles.                                                              |
| No          | Attack              | [Commands.md → Attack](./Commands.md#attack)                                       | —                                                                                                    | —          | Not implemented. Surface/town initiation nuances; negative: deny where inappropriate.                                                                                  |
| No          | Fire                | [Commands.md → Fire — Town/Ship](./Commands.md#fire-cannons)                       | —                                                                                                    | —          | Not implemented. Large map negative should be “Fire broadsides only!” or “What?/Cannot” contextually.                                                                  |
| No          | Cast                | [Commands.md → Cast](./Commands.md#cast), Spells.md                                | —                                                                                                    | —          | Pseudocode added; not implemented. Context gating (Not here!/Absorbed!), stock/MP checks, success reporting, turn rules.                                               |
| No          | New Order           | [Commands.md → New Order](./Commands.md#new-order-swap-party-positions)            | —                                                                                                    | —          | Not implemented. Negative handling per legacy when invalid.                                                                                                            |
| No          | Fire (Cannons)      | [Commands.md → Fire — Town/Ship](./Commands.md#fire-cannons)                       | —                                                                                                    | —          | Not implemented.                                                                                                                                                       |
| No          | Search              | [Commands.md → Search](./Commands.md#search)                                       | —                                                                                                    | —          | Not implemented.                                                                                                                                                       |
| No          | Yell                | [Commands.md → Yell](./Commands.md#yell)                                           | `OLD/CMDS.C`                                                                                         | —          | Summons Shadowlords in towns, opens/closes dungeons and restores shrines outside, and furls/hoists sails on ships.                                                     |
| No          | Escape              | [Commands.md → Escape](./Commands.md#escape)                                       | `OLD/CMDS.C`                                                                                         | —          | Exits a combat screen after winning a battle.                                                                                                                          |

### Dungeon Commands

Dungeon interactions often differ from surface/town flows (ahead‑of‑avatar targeting, dungeon tile families, light checks, underfoot chests, etc.).

| Implemented | Feature                  | Pseudocode Ref                                                         | Code Ref                                              | Similarity | Notes                                                                                                            |
|-------------|--------------------------|------------------------------------------------------------------------|-------------------------------------------------------|------------|------------------------------------------------------------------------------------------------------------------|
| No          | Look — Dungeon           | [Commands.md → Look — Dungeon](./Commands.md#look-—-dungeon)           | —                                                     | —          | Dungeon look requires light and prints tile‑family descriptions; includes fountain drink prompt and field types. |
| No          | Open — Dungeon           | [Commands.md → Open — Dungeon](./Commands.md#open-—-dungeon)           | —                                                     | —          | Handles dungeon doors and underfoot chests; integrates with spells (An Sanct/In Ex Por).                         |
| No          | Jimmy — Dungeon          | [Commands.md → Jimmy](./Commands.md#jimmy)                             | —                                                     | —          | Dungeon variant targets ahead tile and supports dungeon chest jimmy odds.                                        |
| No          | Get — Dungeon            | [Commands.md → Get — Dungeon](./Commands.md#get-—-dungeon)             | —                                                     | —          | Picks from underfoot opened chest; distinct from surface object pickup.                                          |
| No          | Search — Dungeon (Ahead) | [Commands.md → Search — Dungeon](./Commands.md#search-—-dungeon-ahead) | —                                                     | —          | Ahead‑of‑avatar search for secret doors/passages; separate flow.                                                 |
| No          | View (Gem Map) — Dungeon | [Commands.md → View (Gem Map)](./Commands.md#view-gem-map)             | —                                                     | —          | Renders dungeon level cell layout; consumes gem.                                                                 |
| No          | Pass Turn — Dungeon Tick | [Commands.md → Pass Turn (Space)](./Commands.md#pass-turn-space)       | `cmd/ultimav/gamescene_input.go` (no Dungeon handler) | —          | Should advance dungeon hazards/lighting per tick.                                                                |

## Fixtures & Environment

| Implemented | Feature                   | Pseudocode Ref                                                           | Code Ref | Similarity | Notes                                         |
|-------------|---------------------------|--------------------------------------------------------------------------|----------|------------|-----------------------------------------------|
| No          | Fixtures: default mapping | [Fixtures.md → Fixture Defaults](./Fixtures.md#fixture-defaults-mapping) | —        | —          | No FixtureEffects map or Use routing present. |
| No          | Wells: Wish               | [Fixtures.md → Wells & Wish](./Fixtures.md#wells--wish)                  | —        | —          | Missing.                                      |
| No          | Fountains                 | [Fixtures.md → Fountains](./Fixtures.md#fountains)                       | —        | —          | Missing.                                      |
| No          | Lamps/Sconces overrides   | [Fixtures.md → Lamps/Sconces](./Fixtures.md#lampssconces-overrides)      | —        | —          | Missing.                                      |
| No          | Overworld hazards         | [Environment.md](./Environment.md)                                       | —        | —          | Missing.                                      |
| No          | Moongates                 | [Moongates.md](./Moongates.md)                                           | —        | —          | Missing.                                      |
| No          | Town drawbridges          | [Towns.md → Drawbridges](./Towns.md#drawbridges)                         | —        | —          | Missing.                                      |

## Schedules & AI

| Implemented | Feature                           | Pseudocode Ref                                                                 | Code Ref                                                 | Similarity | Notes                                                                      |
|-------------|-----------------------------------|--------------------------------------------------------------------------------|----------------------------------------------------------|------------|----------------------------------------------------------------------------|
| Yes         | NPC schedules (data/model)        | [NPC_Schedules.md → Data Model](./NPC_Schedules.md#data-model)                 | `internal/references/npc_schedule.go`                    | Similar    | Schedule model present; details may differ.                                |
| Partial     | NPC schedule driver (hour change) | [NPC_Schedules.md → Hourly Transitions](./NPC_Schedules.md#hourly-transitions) | `internal/ai/npc_ai_controller_small_map.go` (various)   | Similar    | Controller selects behaviors and floors; exact LEAV/ARIV/POP not verbatim. |
| Yes         | Small map pathfinding             | [NPC_Schedules.md → Pathfinding](./NPC_Schedules.md#pathfinding)               | `internal/astar/*.go`, `internal/ai/npc_ai_controller_*` | Similar    | Pathfinding exists; integration with schedules ongoing.                    |
| Yes         | Large map monster generation      | [Movement_Combat_AI.md → Monster Generation](./Movement_Combat_AI.md)          | `internal/ai/npc_ai_controller_large_map.go`            | Similar    | Environment-based monster spawning with tile probability system implemented. Fixed double-gating issue in spawn rates. |
| No          | Combat AI (seek, special moves)   | [Movement_Combat_AI.md](./Movement_Combat_AI.md)                               | —                                                        | —          | Combat not implemented.                                                    |
| No          | Mass charm targeting ('C')        | [Spells.md → Quas An Wis](./Spells.md#quas-an-wis-mass-charmconfusion)         | —                                                        | —          | Not applicable yet.                                                        |

## Spells & Scrolls

| Implemented | Feature                 | Pseudocode Ref                                            | Code Ref                                        | Similarity | Notes                                                |
|-------------|-------------------------|-----------------------------------------------------------|-------------------------------------------------|------------|------------------------------------------------------|
| No          | Spellcasting core       | [Spells.md](./Spells.md)                                  | —                                               | —          | Not present in Go codebase.                          |
| No          | Specific spells/scrolls | [Spells.md (All)](./Spells.md#spells-summary-at-a-glance) | `internal/references/data/InventoryDetails.csv` | —          | Data present for names/info; no runtime casting/use. |

### Spell Inventory Data (FYI)

| Implemented | Feature                     | Pseudocode Ref                                      | Code Ref                                            | Similarity | Notes                                |
|-------------|-----------------------------|-----------------------------------------------------|-----------------------------------------------------|------------|--------------------------------------|
| Yes         | Spell metadata (names/info) | [Spells.md](./Spells.md)                            | `internal/references/data/InventoryDetails.csv`     | N/A        | Data present; no runtime casting.    |
| Partial     | Inventory quantities        | [SAVED_GAM_STRUCTURE.md](../SAVED_GAM_STRUCTURE.md) | `internal/party_state/inventory.go` (Scrolls, etc.) | N/A        | Data structures exist; no use flows. |

### Potions & Scrolls

| Implemented | Feature        | Pseudocode Ref                                                         | Code Ref | Similarity | Notes            |
|-------------|----------------|------------------------------------------------------------------------|----------|------------|------------------|
| No          | Potion effects | [Potions.md](./Potions.md)                                             | —        | —          | Not implemented. |
| No          | Scroll effects | [Spells.md → Scrolls Summary](./Spells.md#scrolls-summary-at-a-glance) | —        | —          | Not implemented. |

## Special Items & Artifacts

| Implemented | Feature/Item       | Pseudocode Ref                                                                     | Code Ref                                          | Similarity | Notes                                          |
|-------------|--------------------|------------------------------------------------------------------------------------|---------------------------------------------------|------------|------------------------------------------------|
| Partial     | Magic Carpet (Use) | [Commands.md → Use](./Commands.md#use)                                             | `internal/map_units/npc_vehicle.go` (carpet type) | Dissimilar | Carpet NPC exists; use/place/pickup logic TBD. |
| No          | Skull Keys (Use)   | [Objects.md → Skull Key — Magical Unlock](./Objects.md#skull-key-—-magical-unlock) | `internal/party_state/inventory.go` (SkullKeys)   | —          | Inventory tracked; no use flow.                |
| No          | Crown (Use)        | [Commands.md → Use](./Commands.md#use)                                             | —                                                 | —          | Not implemented.                               |
| No          | Sceptre (Use)      | [Commands.md → Use](./Commands.md#use)                                             | —                                                 | —          | Not implemented.                               |
| No          | Amulet (Use)       | [Commands.md → Use](./Commands.md#use)                                             | `internal/references/item_equipment.go`           | —          | Item enum exists; effect/use not wired.        |
| No          | Spyglass/Telescope | [Fixtures.md → Telescope](./Fixtures.md#telescope)                                 | —                                                 | —          | Not implemented.                               |
| No          | Gems (View)        | [Commands.md → View (Gem Map)](./Commands.md#view-gem-map)                         | —                                                 | —          | Not implemented.                               |
| Partial     | Torches (Ignite)   | [Commands.md → Ignite Torch](./Commands.md#ignite-torch)                           | `internal/map_state/lighting.go`                  | Similar    | Lighting supports torches; command missing.    |

## Town Systems

| Implemented | Feature                        | Pseudocode Ref                                                         | Code Ref | Similarity | Notes                        |
|-------------|--------------------------------|------------------------------------------------------------------------|----------|------------|------------------------------|
| No          | Guard alarm/pursuit            | [Towns.md → Special Guard Behavior](./Towns.md#special-guard-behavior) | —        | —          | Not implemented.             |
| No          | Jail flow                      | [Towns.md → Jail Flow](./Towns.md#jail-flow)                           | —        | —          | Not implemented.             |
| No          | Cannons (town/ship broadsides) | [Commands.md → Fire](./Commands.md#fire-cannons)                       | —        | —          | Not implemented.             |
| No          | Shops (pricing/services)       | [Shops.md](./Shops.md)                                                 | —        | —          | Pricing tables to be filled. |

## Potions & Scrolls

| Implemented | Feature       | Pseudocode Ref      | Code Ref                                    | Similarity | Notes                                     |
|-------------|---------------|---------------------|---------------------------------------------|------------|-------------------------------------------|
| No          | Potions (use) | Potions.md          | `internal/party_state/inventory.go` (types) | —          | Quantities exist but no use/effect logic. |
| No          | Scrolls (use) | Spells.md (scrolls) | `internal/party_state/inventory.go` (types) | —          | No scroll use flows implemented.          |

## Special Items

| Implemented | Feature            | Pseudocode Ref          | Code Ref                                                     | Similarity | Notes                                                   |
|-------------|--------------------|-------------------------|--------------------------------------------------------------|------------|---------------------------------------------------------|
| No          | Crown (Use)        | Commands.md → Use       | —                                                            | —          | Not implemented.                                        |
| No          | Sceptre (Use)      | Commands.md → Use       | —                                                            | —          | Not implemented.                                        |
| No          | Amulet (Use)       | Commands.md → Use       | `internal/party_state/types.go` (Amulet field)               | —          | Field exists; no use effect logic.                      |
| Partial     | Carpet (Board/Use) | Commands.md → Use/Board | `internal/map_units/npc_vehicle.go` (CarpetVehicle), map use | Dissimilar | Vehicle types exist; no boarding/Use flows as per docs. |
| No          | Spyglass/Telescope | Commands.md → Use/Look  | —                                                            | —          | Not implemented.                                        |
| No          | Gems (View)        | Commands.md → View      | `internal/party_state/inventory.go` (gems qty)               | —          | Quantities exist; View command absent.                  |

## Exhaustive Spell Checklist (48 Spells)

Legend: Implemented = No (unless otherwise noted), Similarity = —, Code Ref column lists any related data structures.

| Implemented | #  | Spell            | Pseudocode Ref | Code Ref                                        | Similarity | Notes                       |
|-------------|----|------------------|----------------|-------------------------------------------------|------------|-----------------------------|
| No          | 0  | In Lor           | Spells.md      | `internal/references/data/InventoryDetails.csv` | —          | —                           |
| No          | 1  | Grav Por         | Spells.md      | same                                            | —          | —                           |
| No          | 2  | An Zu            | Spells.md      | same                                            | —          | —                           |
| No          | 3  | An Nox           | Spells.md      | same                                            | —          | —                           |
| No          | 4  | Mani             | Spells.md      | same                                            | —          | —                           |
| No          | 5  | An Ylem          | Spells.md      | same                                            | —          | —                           |
| No          | 6  | An Sanct         | Spells.md      | same                                            | —          | —                           |
| No          | 7  | An Xen Corp      | Spells.md      | same                                            | —          | —                           |
| No          | 8  | Rel Hur          | Spells.md      | same                                            | —          | —                           |
| No          | 9  | In Wis           | Spells.md      | same                                            | —          | —                           |
| No          | 10 | Kal Xen          | Spells.md      | same                                            | —          | —                           |
| No          | 11 | In Xen Mani      | Spells.md      | same                                            | —          | —                           |
| No          | 12 | Vas Lor          | Spells.md      | same                                            | —          | —                           |
| No          | 13 | Vas Flam         | Spells.md      | same                                            | —          | —                           |
| No          | 14 | In Flam Grav     | Spells.md      | same                                            | —          | —                           |
| No          | 15 | In Nox Grav      | Spells.md      | same                                            | —          | —                           |
| No          | 16 | In Zu Grav       | Spells.md      | same                                            | —          | —                           |
| No          | 17 | In Por           | Spells.md      | same                                            | —          | —                           |
| No          | 18 | An Grav          | Spells.md      | same                                            | —          | —                           |
| No          | 19 | In Sanct         | Spells.md      | same                                            | —          | —                           |
| No          | 20 | In Sanct G       | Spells.md      | same                                            | —          | Energy field create variant |
| No          | 21 | Uus Por          | Spells.md      | same                                            | —          | —                           |
| No          | 22 | Des Por          | Spells.md      | same                                            | —          | —                           |
| No          | 23 | Wis Quas         | Spells.md      | same                                            | —          | —                           |
| No          | 24 | In Bet Xen       | Spells.md      | same                                            | —          | —                           |
| No          | 25 | An Ex Por        | Spells.md      | same                                            | —          | —                           |
| No          | 26 | In Ex Por        | Spells.md      | same                                            | —          | —                           |
| No          | 27 | Vas Mani         | Spells.md      | same                                            | —          | —                           |
| No          | 28 | In Zu            | Spells.md      | same                                            | —          | —                           |
| No          | 29 | Rel Tym          | Spells.md      | same                                            | —          | —                           |
| No          | 30 | In Vas Por Ylem  | Spells.md      | same                                            | —          | —                           |
| No          | 31 | Quas An Wis      | Spells.md      | same                                            | —          | Mass charm aura             |
| No          | 32 | In An            | Spells.md      | same                                            | —          | Negate magic                |
| No          | 33 | Wis An Ylem      | Spells.md      | same                                            | —          | X‑Ray                       |
| No          | 34 | An Xen Ex        | Spells.md      | same                                            | —          | —                           |
| No          | 35 | Rel Xen Bet      | Spells.md      | same                                            | —          | Polymorph                   |
| No          | 36 | Sanct Lor        | Spells.md      | same                                            | —          | Invisibility                |
| No          | 37 | Xen Corp         | Spells.md      | same                                            | —          | —                           |
| No          | 38 | In Quas Xen      | Spells.md      | same                                            | —          | Clone                       |
| No          | 39 | In Quas Wis      | Spells.md      | same                                            | —          | View                        |
| No          | 40 | In Nox Hur       | Spells.md      | same                                            | —          | Poison storm                |
| No          | 41 | In Quas Corp     | Spells.md      | same                                            | —          | Fear                        |
| No          | 42 | In Mani Corp     | Spells.md      | same                                            | —          | Resurrection                |
| No          | 43 | Kal Xen Corp     | Spells.md      | same                                            | —          | Summon daemon               |
| No          | 44 | In Vas Grav Corp | Spells.md      | same                                            | —          | Energy storm                |
| No          | 45 | In Flam Hur      | Spells.md      | same                                            | —          | Firestorm                   |
| No          | 46 | Vas Rel Por      | Spells.md      | same                                            | —          | Gate Travel                 |
| No          | 47 | An Tym           | Spells.md      | same                                            | —          | Negate Time                 |
| No          | 48 | Frotz (reserved) | Spells.md      | same                                            | —          | Unimplemented               |

## Scrolls Checklist (8)

| Implemented | Scroll        | Pseudocode Ref | Code Ref                                    | Similarity | Notes |
|-------------|---------------|----------------|---------------------------------------------|------------|-------|
| No          | Light         | Spells.md      | `internal/party_state/inventory.go` (types) | —          | —     |
| No          | Wind Change   | Spells.md      | same                                        | —          | —     |
| No          | Protection    | Spells.md      | same                                        | —          | —     |
| No          | Negate Magic  | Spells.md      | same                                        | —          | —     |
| No          | View          | Spells.md      | same                                        | —          | —     |
| No          | Summon Daemon | Spells.md      | same                                        | —          | —     |
| No          | Resurrection  | Spells.md      | same                                        | —          | —     |
| No          | Negate Time   | Spells.md      | same                                        | —          | —     |

## Potions Checklist (8)

| Implemented | Potion Color | Pseudocode Ref | Code Ref                                    | Similarity | Notes              |
|-------------|--------------|----------------|---------------------------------------------|------------|--------------------|
| No          | Blue         | Potions.md     | `internal/party_state/inventory.go` (types) | —          | Cure Sleep         |
| No          | Yellow       | Potions.md     | same                                        | —          | Heal               |
| No          | Red          | Potions.md     | same                                        | —          | Cure Poison        |
| No          | Green        | Potions.md     | same                                        | —          | Poison             |
| No          | Orange       | Potions.md     | same                                        | —          | Sleep              |
| No          | Purple       | Potions.md     | same                                        | —          | Polymorph (combat) |
| No          | Black        | Potions.md     | same                                        | —          | Invisible (combat) |
| No          | White        | Potions.md     | same                                        | —          | X‑Ray (surface)    |

## Special Items Checklist

| Implemented | Item                  | Pseudocode Ref       | Code Ref                                            | Similarity | Notes                                      |
|-------------|-----------------------|----------------------|-----------------------------------------------------|------------|--------------------------------------------|
| No          | Crown of Lord British | Commands.md → Use    | —                                                   | —          | —                                          |
| No          | Sceptre of Lord Brit. | Commands.md → Use    | —                                                   | —          | —                                          |
| No          | Amulet of LB          | Commands.md → Use    | `internal/party_state/types.go` (Amulet field)      | —          | Field present only                         |
| Partial     | Magic Carpet          | Commands.md → Use    | `internal/map_units/npc_vehicle.go` (CarpetVehicle) | Dissimilar | Vehicle type present; no Use/boarding flow |
| No          | Skull Keys            | Commands.md → Use    | —                                                   | —          | —                                          |
| No          | Keys                  | Commands.md → Open   | `internal/party_state/inventory.go` (keys qty)      | —          | No Open/door flows                         |
| No          | Torches               | Commands.md → Ignite | `internal/party_state/inventory.go` (torches qty)   | —          | No Ignite Torch command                    |
| No          | Gems                  | Commands.md → View   | `internal/party_state/inventory.go` (gems qty)      | —          | View not implemented                       |
| No          | Spyglass              | Commands.md → Use    | —                                                   | —          | —                                          |
| No          | Telescope             | Commands.md → Look   | —                                                   | —          | —                                          |

## Towns & Special Systems

| Implemented | Feature                | Pseudocode Ref                      | Code Ref | Similarity | Notes      |
|-------------|------------------------|-------------------------------------|----------|------------|------------|
| No          | Drawbridges/Portcullis | Towns.md                            | —        | —          | Not found. |
| No          | Guard alarm & Jail     | Towns.md → Guard Behavior/Jail      | —        | —          | Not found. |
| No          | Cannons (town fire)    | Combat_Effects.md/Towns.md          | —        | —          | Not found. |
| No          | Bridge trolls          | Special_BridgeTrolls.md             | —        | —          | Not found. |
| No          | Wind system            | Movement_Overworld.md → Wind System | —        | —          | Not found. |
| No          | Ships & Sails          | Commands.md / Movement_Overworld.md | —        | —          | Not found. |
| No          | Moongates              | Moongates.md                        | —        | —          | Not found. |

## Shops & Economy

| Implemented | Feature                    | Pseudocode Ref       | Code Ref                                                            | Similarity | Notes                                                                                                |
|-------------|----------------------------|----------------------|---------------------------------------------------------------------|------------|------------------------------------------------------------------------------------------------------|
| No          | Shop pricing & multipliers | Shops.md             | —                                                                   | —          | Tables added in docs; not implemented in code.                                                       |
| No          | Reagent/Healer/Arms Shops  | Shops.md             | —                                                                   | —          | Not present.                                                                                         |
| No          | Inns (stay months)         | Shops.md → Innkeeper | `internal/party_state/types.go` (MonthsAtInn, PartyStatus=AtTheInn) | —          | Data fields present; no inn UI/flow, no gold deduction, no time advancement, no `inn_party` mapping. |
| No          | Horses/Shipwright          | Shops.md             | —                                                                   | —          | Not present.                                                                                         |

## Conversation System (FYI)

| Implemented | Feature              | Pseudocode Ref           | Code Ref                                          | Similarity | Notes                                        |
|-------------|----------------------|--------------------------|---------------------------------------------------|------------|----------------------------------------------|
| Yes         | Talk engine (linear) | TALK_SYSTEM_STRUCTURE.md | `internal/conversation/linear_engine.go` (+tests) | Similar    | Robust implementation; outside combat scope. |

---

This tracker is a starting point. As features land in Go code, update the “Implemented” and “Similarity” columns, and add concrete function references.

## Sleep & Rest (Overview)

| Implemented | Feature               | Pseudocode Ref                                                                                   | Code Ref                                         | Similarity | Notes                                                                                                            |
|-------------|-----------------------|--------------------------------------------------------------------------------------------------|--------------------------------------------------|------------|------------------------------------------------------------------------------------------------------------------|
| No          | Sleep status effect   | Combat_Effects.md → Field Effects, Per‑Turn Updates; Potions.md (Blue/Orange); Spells.md (In Zu) | `internal/party_state/mappings.go` (Status enum) | —          | Status enum includes Sleep; no application/expiration logic, no wake/cure handling wired.                        |
| No          | Sleep field (tiles)   | Combat_Effects.md → Field Effects                                                                | —                                                | —          | Combat field effects not implemented; would apply Sleep on contact (PCs/monsters).                               |
| No          | In Zu (Sleep spell)   | Spells.md                                                                                        | —                                                | —          | Listed in spells table; no casting/effect pipeline.                                                              |
| No          | Potions: Blue/Orange  | Potions.md                                                                                       | `internal/party_state/inventory.go` (types only) | —          | Blue cures Sleep; Orange applies Sleep; items exist as types only.                                               |
| No          | Hole Up & Camp (rest) | Commands.md → Hole Up & Camp                                                                     | —                                                | —          | Camping/repair flows absent; should handle guard watch, time advance, HP/MP regen, food ticks, encounter checks. |
| No          | In‑bed sleep/ejection | Fixtures.md → Beds; NPC_Schedules.md (eject sleepers at schedule boundaries)                     | —                                                | —          | No in‑bed sleep flow; must prevent sleeping in occupied beds and eject sleepers around hour changes.             |
