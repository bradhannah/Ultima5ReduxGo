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
