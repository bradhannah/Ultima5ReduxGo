# Environment and Hazards

## Overworld Hazards (Swamp, Lava)

```pseudocode
FUNCTION apply_overworld_tile_hazards(current_tile):
    SWITCH current_tile:
        CASE LAVA: show_message("Burning!\n"); damage_party_on_land()
        CASE SWAMP:
            FOR EACH member IN party_members_in_order():
                IF member.status != DEAD AND member.status != POISONED THEN
                    IF random(1, 30) > member.dexterity THEN member.status = POISONED; show_message("Poisoned!\n")
                ENDIF
            ENDFOR
    ENDSWITCH
ENDFUNCTION

## Town/Indoor Hazards (Small Maps)

```pseudocode
FUNCTION apply_town_tile_hazards(current_tile):
    IF current_tile == FIREPLACE OR current_tile == LAVA THEN
        show_message("Burning!\n")
        damageparty() // each living member takes 1..8
    ENDIF
ENDFUNCTION
```

## Hazard Matrix (Summary)

| Tile       | Overworld (move onto)                | Town/Indoor (move onto)              | Combat (standing each tick)                  |
|------------|--------------------------------------|--------------------------------------|----------------------------------------------|
| SWAMP      | Dex check; on fail: Poisoned         | Dex check; on fail: Poisoned         | PCs only poisoned (monsters unaffected)      |
| LAVA       | “Burning!” + damage party (1..8)     | “Burning!” + damage party (1..8)     | Fire damage `getrandom(10)`                  |
| FIREPLACE  | n/a                                  | “Burning!” + damage party (1..8)     | Fire damage `getrandom(10)`                  |

Notes:

- “Damage party” is the whole-party 1..8 routine, not per-entity field damage.
- In combat, fireplaces and lava are handled via the field-effects system (see Combat Effects → Field Effects).

## Underworld Earthquakes

```pseudocode
FUNCTION maybe_trigger_underworld_earthquake():
    IF is_underworld() AND random(0, 255) == 0x69 THEN
        show_message("EARTHQUAKE!\n"); quake_visuals(); damageparty()
    ENDIF
ENDFUNCTION
```

## Darkness (Sight Suppression)

```pseudocode
FUNCTION update_darkness(lost):
    IF tile_at_player() == NOTHING AND active_spell != LB_AMULET THEN
        sight = 0; IF NOT lost THEN update_screen(); lost = TRUE
    ELSE lost = FALSE; addtime(0)
    RETURN lost
ENDFUNCTION
```

## Light Sources & Vision

Light enables vision in dark contexts (dungeons and underworld interiors). Two timers control light: `torch_light` and `magic_light`. Some items provide a persistent light effect.

```pseudocode
// Returns TRUE if no active light is present
FUNCTION no_light_sources():
    RETURN (torch_light <= 0 AND magic_light <= 0 AND NOT is_persistent_light_active())

FUNCTION is_persistent_light_active():
    // Long-lived item effects (e.g., LB Amulet or Crown) grant light while equipped
    RETURN is_effect_active(LB_AMULET) OR is_effect_active(CROWN)

// Called each world tick in dark contexts to age light durations
FUNCTION consume_light_sources_if_needed():
    IF NOT is_in_dungeon() AND NOT is_underworld_small_map() THEN RETURN
    IF magic_light > 0 THEN magic_light -= 1
    IF torch_light > 0 THEN torch_light -= 1
ENDFUNCTION
```

### Light Sources Table

- Torch: inventory item; `Ignite Torch` consumes one torch and sets `torch_light` duration.
- Wall Torch: surface/town fixture; `Get`ing a wall torch sets `torch_light = 100` and removes the torch tile (`Borrowed!`).
- In Lor (Light): spell; sets `magic_light` to a short duration.
- In Bet Lor (Great Light): spell; sets `magic_light` to a long duration.
- LB Amulet / Crown: `Use` equips a persistent effect (duration 255) which counts as light.

### Torch Duration

```pseudocode
FUNCTION ignite_torch():
    IF torches == 0 THEN show_message("None owned!\n"); RETURN FALSE
    torches -= 1
    IF is_in_dungeon() THEN torch_light = 112 + random(0, 15) ELSE torch_light = 240
    RETURN TRUE
ENDFUNCTION
```

Notes:

- Torches provide light anywhere but are primarily useful in dungeons/underworld; outside, light has no gameplay effect beyond visuals.
- If both `magic_light` and `torch_light` are active, visibility is granted as long as either is non-zero; both timers tick independently.
- Equipping the Amulet or Crown activates a persistent effect treated as light (does not decrement until unequipped or cleared).
