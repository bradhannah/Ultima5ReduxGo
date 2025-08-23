# Dungeon Interactions

This file summarizes common dungeon fixtures, chests/traps, and search behavior.

## Chest Traps and Disarming

```pseudocode
// When a chest is opened in town/dungeon contexts
FUNCTION open_chest_at(x, y, opener):
    // Disarm via spell (An Sanct) toggles chest tile to opened state ahead of time
    IF chest_is_disarmed(x, y) THEN show_message("Chest opened!\n"); RETURN
    // Random trap effect selection (see Combat Effects → Chest Traps)
    trigger_random_trap(opener)
    // After trap, award chest contents per location data
    grant_chest_loot(x, y)
ENDFUNCTION

// Disarm/Unlock via spell (An Sanct)
FUNCTION disarm_or_unlock_target_tile(x, y):
    IF tile_is_dungeon_chest(x, y) THEN
        show_message("Chest opened!\n"); set_tile(x, y, CHEST_OPENED); RETURN TRUE
    ELSE IF tile_is_locked_door(x, y) THEN
        unlock_standard_door(x, y); show_message("Opened!\n"); RETURN TRUE
    ENDIF
    RETURN FALSE
ENDFUNCTION
```

Light requirement: `no_light_sources()` gates visibility; see [Environment.md#light-sources--vision](Environment.md#light-sources--vision).

### Chest Trap Probabilities (Drops)

Dropped monster chests use a treasure-value check to decide if trapped:

| Check                          | Probability                   |
|--------------------------------|-------------------------------|
| Chest drops (on monster death) | `rolld30() <= treasure_value` |
| Chest trapped (if dropped)     | `rolld30() < treasure_value`  |

Trap effects and their odds (where used) are covered under Combat Effects → Chest Traps.

## Search and Hidden Features

```pseudocode
FUNCTION search_dungeon_ahead(turning_face):
    IF no_light_sources() THEN show_message("You find:\ndarkness.\n"); RETURN
    IF not_valid_direction(turning_face) THEN RETURN
    tx, ty = tile_in_direction(turning_face)
    tile = get_dungeon_tile(level, ty, tx)

    // Print generic header
    show_message("You find:\n")
    SWITCH tile_category(tile):
        CASE Ladder:
            show_message("Nothing hidden on the ladder.\n")
        CASE Trap:
            // Reveal trap difficulty (simple/complex) using Dex vs difficulty
            difficulty = (30 + 2*level - player_dex) / 2
            IF random(1, 30) > difficulty THEN
                trap_level = level // no trap or reveal fixed trap
            ELSE
                trap_level = random(1, 8)
            ENDIF
            IF trap_level < 4 THEN show_message("A simple trap\n")
            ELSE IF trap_level >= 7 THEN show_message("A complex trap\n")
            ELSE show_message("A trap\n")
        CASE Fountain:
            show_message("Nothing hidden on the fountain.\n")
        DEFAULT:
            show_message("Nothing of note.\n")
    ENDSWITCH
ENDFUNCTION
```

## Dungeon Fixtures Quick Table

| Fixture/Tile             | Use/Effect (Out of Combat)              | Notes                                      |
|--------------------------|-----------------------------------------|--------------------------------------------|
| Ladder Up/Down           | Klimb moves a floor up/down             | Also accessible via directional Klimb      |
| Grate                    | Klimb Down (if matched by map rules)    |                                            |
| Fountain                 | No hidden search reward by default      | Some locations may override to heal/cure   |
| Torch/Light sources      | Increases `torch_light`/`magic_light`   | See Environment → [Light Sources & Vision](Environment.md#light-sources--vision); [Torch Duration](Environment.md#torch-duration) |
| Trap tiles               | Search can reveal “simple/complex” trap | Disarming logic handled by specific spells |
| Fire/Poison/Sleep fields | Field effects per tick in combat        | See Combat Effects → Field Effects         |

Note: Uus Por (Up) and Des Por (Down) provide magical floor changes equivalent to ladders; see Spells → Uus Por / Des Por.


### Treasure Drop/Trap Matrix (Buckets)

Chest drops and trap checks on monster death use a 1d30 roll compared to the monster's `treasure_value`.

- Drop check: drops if `rolld30() <= treasure_value` → P(drop) = treasure_value / 30
- Trap check (if dropped): trapped if `rolld30() < treasure_value` → P(trap|drop) ≈ (treasure_value - 1) / 30

| Treasure Value | P(Drop) | P(Trap | Drop) | Notes                         |
|----------------|:-------:|:--------------:|-------------------------------|
| 0              |  0.0%   |      —         | No chest ever drops           |
| 5              | 16.7%   |     13.3%      | Low-tier monsters             |
| 10             | 33.3%   |     30.0%      | Common low/mid                |
| 15             | 50.0%   |     46.7%      | Mid-tier                      |
| 20             | 66.7%   |     63.3%      | Mid/high                      |
| 25             | 83.3%   |     80.0%      | High-tier                     |
| 30             | 100.0%  |     96.7%      | Always drops; almost always trapped |

Balancing tips:

- Use treasure_value as the single knob for both drop and trap likelihood; higher value → more drops and more traps.
- If you want frequent drops but fewer traps, consider clamping the trap check (e.g., `min(treasure_value - k, 30)`).
- For very low treasure_value monsters, you can eliminate traps entirely by setting a minimum threshold before traps are allowed.
