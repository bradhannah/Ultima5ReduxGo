# Small-Town Fixtures

Many fixtures are decorative, but the engine can attach optional, data-driven interactions for player "Use" actions. Combat-specific effects for some fixtures (e.g., fireplaces) are covered under Combat Effects.

## Data-Driven Effects

```pseudocode
STRUCT FixtureUseEffect {
    Kind: enum { None, HealSmall, CurePoison, FillWater, ToggleLight, Custom }
    // Optional parameters
    AmountMin: int // min heal, etc.
    AmountMax: int // max heal, etc.
    CustomID: int  // engine-defined hook for special scripted effects
}

// Map tile -> effect mapping supplied by data (per location if needed)
MAP<Tile, FixtureUseEffect> FixtureEffects
```

## Interaction Entry Point

```pseudocode
FUNCTION interact_with_fixture_at(x, y, active_player):
    tile = top_tile_at(x, y)
    effect = FixtureEffects.get(tile, default=FixtureUseEffect{Kind=None})
    SWITCH effect.Kind:
        CASE None:
            show_message("Nothing happens.\n")
        CASE HealSmall:
            amount = random(effect.AmountMin, effect.AmountMax)
            healed = heal_amount(active_player, amount)
            IF healed > 0 THEN show_message("Refreshed!\n"); play_heal_sfx(); mark_stats_changed()
        CASE CurePoison:
            IF player[active_player].status == POISONED THEN
                player[active_player].status = GOOD; show_message("Poison cured!\n"); play_heal_sfx(); mark_stats_changed()
            ELSE show_message("Refreshing!\n")
        CASE FillWater:
            IF inventory_has_empty_flask(active_player) THEN
                fill_flask(active_player); show_message("Filled flask.\n"); mark_stats_changed()
            ELSE show_message("Nothing to fill.\n")
        CASE ToggleLight:
            toggle_local_light_at(x, y); show_message("The light flickers.\n")
        CASE Custom:
            run_custom_fixture_effect(effect.CustomID, x, y, active_player)
    ENDSWITCH
ENDFUNCTION
```

## Common Fixtures

```pseudocode
// Fountain
// Default: CurePoison OR small heal, per data.
FixtureEffects[FOUNTAIN] = { Kind: CurePoison }

// Fireplace / Brazier / Candle / Sconce
// Default: ToggleLight (purely cosmetic); in-combat fire damage already handled via field effects or terrain.
FixtureEffects[FIREPLACE] = { Kind: ToggleLight }
FixtureEffects[BRAZIER]   = { Kind: ToggleLight }
FixtureEffects[CANDLE]    = { Kind: ToggleLight }
FixtureEffects[S_CONE]    = { Kind: ToggleLight }

// Well
// Default: FillWater into an empty flask if present.
FixtureEffects[WELL] = { Kind: FillWater }

// Telescope
// Default: Custom effect to show a long-range "look" in a direction.
FixtureEffects[TELESCOPE] = { Kind: Custom, CustomID: TELESCOPE_LOOK_ID }

FUNCTION run_custom_fixture_effect(customID, x, y, active_player):
    SWITCH customID:
        CASE TELESCOPE_LOOK_ID:
            dir = ask_direction("Look-")
            show_long_range_view_from(x, y, dir)
        DEFAULT:
            show_message("Nothing happens.\n")
    ENDSWITCH
ENDFUNCTION
```

Notes:

- Keep fixture effects deterministic; for small heals, use bounded random driven by the central PRNG.
- Effects can be scoped per-location by layering a location-specific `FixtureEffects` map over the global defaults.
- Combat-related damage for fireplaces and lava is documented in Combat Effects; this file focuses on out-of-combat "Use" interactions.

## Fixture Defaults Table

| Tile            | Default Effect | UI Text               | Sound        | Combat Interaction                    | Notes                                      |
|-----------------|----------------|-----------------------|--------------|---------------------------------------|--------------------------------------------|
| `Fountain1..4`  | CurePoison     | "Poison cured!"       | heal chime   | none                                  | Can be configured to small heal instead.   |
| `Well`          | FillWater      | "Filled flask."       | water fill   | none                                  | Requires empty flask in inventory.         |
| `Fireplace`     | ToggleLight    | "The light flickers." | light toggle | Fire damage handled in Combat Effects | Visual only here.                          |
| `Brazier`       | ToggleLight    | "The light flickers." | light toggle | Fire damage handled in Combat Effects | Visual only here.                          |
| `RightSconce`   | ToggleLight    | "The light flickers." | light toggle | none                                  | Visual only.                               |
| `LeftSconce`    | ToggleLight    | "The light flickers." | light toggle | none                                  | Visual only.                               |
| `CandleOnTable` | ToggleLight    | "The light flickers." | light toggle | none                                  | Visual only.                               |
| `LampPost`      | ToggleLight    | "The light flickers." | light toggle | none                                  | Visual only.                               |
| `Telescope`     | Custom         | "Look-" prompt        | lens sweep   | none                                  | Shows long-range view in chosen direction. |

## Location Overrides Schema

To override fixture behavior for specific locations, supply a per-location mapping:

```pseudocode
STRUCT LocationFixtureOverrides {
    LocationID: string
    Overrides: MAP<Tile, FixtureUseEffect>
}

// Example: A specific town’s fountain heals 5..10 HP instead of curing poison.
LocationFixtureOverrides{
    LocationID: "Town_X",
    Overrides: {
        FOUNTAIN1: { Kind: HealSmall, AmountMin: 5, AmountMax: 10 },
        FOUNTAIN2: { Kind: HealSmall, AmountMin: 5, AmountMax: 10 },
    }
}
```

## Testing Guidance

- Use fixed PRNG seeds and mock inventory to verify FillWater requires an empty flask.
- Verify ToggleLight effects are idempotent and do not affect combat damage paths.
- For fountains, test both CurePoison and HealSmall overrides; ensure UI messages align with effect.

## Wells and Wish (Expanded)

Wells can optionally support a “Wish” interaction when used. By default, wells only fill a flask. Locations may override a well to grant a wish on coin toss.

### Wish Flow

```pseudocode
FUNCTION wish(x, y, level, active_player):
    show_message("Drop a coin?\n")
    IF NOT read_yes_no() THEN show_message("No.\n"); RETURN FALSE
    IF gold <= 0 THEN show_message("No gold!\n"); RETURN FALSE
    gold -= 1; mark_stats_changed()

    // Fetch per-location outcome table; fallback to defaults
    table = get_wish_table_for_location(current_location_id()) // list of (Outcome, Weight)
    outcome = weighted_pick(table)

    SWITCH outcome.Kind:
        CASE HealSmall:
            n = random(outcome.Min, outcome.Max)
            healed = heal_amount(active_player, n)
            IF healed > 0 THEN show_message("Refreshed!\n"); play_heal_sfx(); RETURN TRUE
        CASE CurePoison:
            IF player[active_player].status == POISONED THEN player[active_player].status = GOOD; show_message("Cured!\n"); play_heal_sfx(); RETURN TRUE
        CASE GrantFood:
            inc_food(outcome.Amount); show_message("Food!\n"); mark_stats_changed(); RETURN TRUE
        CASE GrantGold:
            inc_gold(outcome.Amount); show_message("Gold!\n"); mark_stats_changed(); RETURN TRUE
        CASE SpawnHorse:
            IF can_spawn_horse_near(x, y) THEN spawn_horse_near(x, y); show_message("A fine steed!\n"); RETURN TRUE
        CASE BoostStat:
            stat = outcome.Stat // STR|DEX|INT
            IF getattr(active_player, stat) < outcome.MaxStat THEN inc_stat(active_player, stat, 1); show_message("Thy %s increases!\n", stat); RETURN TRUE
        CASE Nothing:
            show_message("Nothing happens.\n"); RETURN TRUE
        CASE Custom:
            run_custom_wish_effect(outcome.CustomID, x, y, active_player); RETURN TRUE
    ENDSWITCH
    RETURN FALSE
ENDFUNCTION
```

Notes:

- Deterministic: Use weighted tables and the central PRNG. Keep weights and effects in data for per-location tuning.
- Wishes consume exactly 1 gold on “Yes”. If gold is 0, the wish is not attempted.

### Default Wish Outcomes (Template)

| Outcome      | Weight | Params                | UI/Text            | Notes                               |
|--------------|:------:|-----------------------|--------------------|-------------------------------------|
| HealSmall    |   3    | Min=3, Max=10         | “Refreshed!”       | Small heal                           |
| CurePoison   |   2    | —                     | “Cured!”           | Only if poisoned                     |
| GrantFood    |   2    | Amount=3              | “Food!”            | Adds food                            |
| GrantGold    |   2    | Amount=25             | “Gold!”            | Adds gold                            |
| SpawnHorse   |   1    | —                     | “A fine steed!”    | Spawns horse near well if possible   |
| BoostStat    |   1    | Stat=INT, MaxStat=30  | “Thy INT increases!” | +1 up to cap                      |
| Nothing      |   5    | —                     | “Nothing happens.” | No effect                            |

Notes:

- These weights/values are placeholders; override per location for canonical behavior.
- Some towns may forbid horses or cap BoostStat differently; enforce via per-location logic.

### Location Overrides — Wells

| Location        | Outcome Overrides                                                                 |
|-----------------|------------------------------------------------------------------------------------|
| Britain         | Increase HealSmall weight; disable SpawnHorse                                      |
| Moonglow        | Favor BoostStat: INT; lower GrantGold                                              |
| Jhelom          | Favor BoostStat: STR; allow SpawnHorse                                             |
| Skara Brae      | Favor SpawnHorse if space allows; otherwise Nothing                                |
| Buccaneer’s Den | Favor GrantGold; increase Nothing                                                  |

Provide exact weights/params in a data table once canonical values are established.

### Wells — Per‑Town Wish Weights (Template)

Weights are relative integers. Higher = more likely. Set to 0 to disable an outcome for that location.

| Location         | HealSmall | CurePoison | GrantFood | GrantGold | SpawnHorse | BoostStat | Nothing | Custom | Notes                         |
|------------------|:---------:|:----------:|:---------:|:---------:|:----------:|:---------:|:-------:|:------:|-------------------------------|
| Britain          |    6      |     2      |    2      |    2      |     0      |    1      |   3     |   0    | Focus on heal; no horses       |
| Moonglow         |    2      |     2      |    2      |    1      |     0      |    4      |   3     |   0    | BoostStat (INT) favored        |
| Jhelom           |    2      |     1      |    2      |    2      |     1      |    4      |   3     |   0    | BoostStat (STR) favored        |
| Yew              |    3      |     3      |    2      |    1      |     0      |    1      |   3     |   0    | Balanced, no horses            |
| Minoc            |    3      |     2      |    2      |    2      |     0      |    1      |   3     |   0    | —                             |
| Trinsic          |    2      |     3      |    2      |    2      |     0      |    1      |   3     |   0    | —                             |
| Skara Brae       |    2      |     1      |    2      |    1      |     4      |    1      |   3     |   0    | Horses favored (if space)      |
| New Magincia     |    3      |     2      |    2      |    1      |     0      |    2      |   3     |   0    | —                             |
| Cove             |    2      |     3      |    2      |    1      |     0      |    1      |   4     |   0    | Modest effects                 |
| Buccaneer’s Den  |    1      |     1      |    2      |    5      |     0      |    1      |   4     |   0    | Gold favored, more “Nothing”   |
| Paws             |    3      |     2      |    3      |    2      |     0      |    1      |   3     |   0    | Food favored                   |
| Castle British   |    5      |     2      |    2      |    1      |     0      |    2      |   2     |   0    | Generous heal                  |
| Empath Abbey     |    4      |     2      |    2      |    1      |     0      |    3      |   2     |   0    | BoostStat modest               |

Notes:

- These weights are placeholders; replace with canonical values when known. Pair weights with the Default Wish Outcomes table parameters (e.g., BoostStat=INT in Moonglow, STR in Jhelom).
- Disable outcomes by setting weight 0 for towns where effects are inappropriate (e.g., SpawnHorse in city centers).
 - Disable outcomes by setting weight 0 for towns where effects are inappropriate (e.g., SpawnHorse in city centers).

## Rare Fixtures & Edge Cases

Some tiles are primarily structural or decorative but have specific passability or messaging that can trip implementations. Keep these consistent.

### Structural Openings & Furniture

| Tile/Family          | Walk Passable | Missile Passable | Look/Use Text                         | Notes                                      |
|----------------------|---------------|------------------|---------------------------------------|--------------------------------------------|
| Arrow Slit           | No            | Yes              | “an arrow slit”                       | Blocks movement; allow missile LoS         |
| Window               | No            | Yes              | “a window”                            | As above                                   |
| Window Shelf         | No            | Yes              | “a window shelf”                      | Decorative; treat like window              |
| Bookshelf/Crowded    | No            | No               | “a crowded bookshelf”                 | Blocks missiles as well                    |
| Occupied Bed         | No            | No               | “an occupied bed”                     | Prevents in‑bed sleep                      |
| Bed (empty)          | Edge only     | No               | “a bed”                               | Sleep via in‑bed flow only                 |
| Grandfather Clock    | N/A           | N/A              | “a grandfather clock, showing: HH:MM” | Prints time via Look surface               |
| Hourglass            | N/A           | N/A              | “an hourglass”                        | Decorative; no Use effect                  |
| Bellows              | N/A           | N/A              | “(bellows)”                           | Animation only (Look animates)             |
| Anvil                | No            | No               | “an anvil”                            | Decorative; blocks                         |

Guidelines:

- Arrow Slits/Windows: Block walking; allow missiles and line‑of‑sight through. Ensure combat LoS and fire_missile honor this.
- Beds: Only the in‑bed sleep flow should be allowed. Deny sleep on Occupied Bed and eject sleepers at schedule boundaries.
- Clocks: For surface Look, detect clocks and print formatted time after tile description (already in Commands → Look).
- Decorative tiles: Provide Look text; default Use = “Nothing happens.”

## Lamps and Light Overrides (By Location)

While default lamp/sconce/candle behavior is ToggleLight, locations may lock lights or coordinate group toggles.

### Lamps/Sconces Overrides

| Location        | Tile(s)           | Override                       | Notes                                    |
|-----------------|-------------------|---------------------------------|------------------------------------------|
| Castle British  | Sconces           | Locked (no ToggleLight)         | Decorative only                          |
| Empath Abbey    | Sconces           | Group toggle (adjacent pair)    | Both sconces flicker together            |
| Cove            | Street lamps      | Auto-on at night (no manual)    | Night cycle managed by day/night system  |

Implement via `LocationFixtureOverrides` with `Custom` kinds where group/auto behavior is required.

## Fountains (Per‑Town Overrides)

Default fountain behavior is CurePoison. Towns may override to a small heal (bounded) or retain cure behavior. Use `LocationFixtureOverrides` to set per‑tile effects.

### Fountain Overrides — Matrix (Template)

| Location         | Tiles            | Effect      | Params          | Notes                                  |
|------------------|------------------|-------------|-----------------|----------------------------------------|
| Britain          | Fountain1..2     | HealSmall   | Min=3, Max=10   | “Refreshed!”                           |
| Moonglow         | Fountain1        | CurePoison  | —               | “Poison cured!”                        |
| Jhelom           | Fountain1        | HealSmall   | Min=2, Max=8    |                                        |
| Yew              | Fountain1..3     | CurePoison  | —               |                                        |
| Minoc            | Fountain1        | HealSmall   | Min=1, Max=6    |                                        |
| Trinsic          | Fountain1        | CurePoison  | —               |                                        |
| Skara Brae       | Fountain1        | HealSmall   | Min=4, Max=12   |                                        |
| New Magincia     | Fountain1        | HealSmall   | Min=2, Max=8    |                                        |
| Cove             | Fountain1        | CurePoison  | —               |                                        |
| Buccaneer’s Den  | Fountain1        | HealSmall   | Min=1, Max=4    |                                        |
| Paws             | Fountain1        | HealSmall   | Min=1, Max=6    |                                        |
| Castle British   | Fountain1..2     | HealSmall   | Min=5, Max=12   | Audience hall                          |
| Empath Abbey     | Fountain1        | HealSmall   | Min=4, Max=10   | Monastic grounds                       |

Notes:

- The entries above are placeholders. Replace with canonical mapping once verified; leave unspecified towns at the default behavior (CurePoison).
- Implement by assigning `FixtureEffects[FOUNTAINx] = { Kind: HealSmall, AmountMin, AmountMax }` in the corresponding `LocationFixtureOverrides` block.
