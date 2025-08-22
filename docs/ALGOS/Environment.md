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
        update_screen()

## Hazard Matrix (Summary)

| Tile       | Overworld (move onto)                | Town/Indoor (move onto)              | Combat (standing each tick)                  |
|------------|--------------------------------------|--------------------------------------|----------------------------------------------|
| SWAMP      | Dex check; on fail: Poisoned         | Dex check; on fail: Poisoned         | PCs only poisoned (monsters unaffected)      |
| LAVA       | “Burning!” + damage party (1..8)     | “Burning!” + damage party (1..8)     | Fire damage `getrandom(10)`                  |
| FIREPLACE  | n/a                                  | “Burning!” + damage party (1..8)     | Fire damage `getrandom(10)`                  |

Notes:

- “Damage party” is the whole-party 1..8 routine, not per-entity field damage.
- Combat tile effects are handled via the field-effects system (see Combat Effects → Field Effects).

        show_message("Burning!\n")
        damageparty() // each living member takes 1..8
    ENDIF
ENDFUNCTION
```

Notes:

- On town/small maps, stepping on a `FIREPLACE` or `LAVA` tile immediately burns the party (same text as overworld lava).
- In combat, fireplaces and lava are also hazardous via the field-effects system (fire flag). See Combat Effects → Field Effects.

```

Note: Interior fixtures like `FIREPLACE` are not hazardous outside combat; they only cause damage during combat via the field-effects system (see Combat Effects → Field Effects).

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

## Torch Duration (Dungeon)

```pseudocode
FUNCTION ignite_torch():
    IF torches == 0 THEN show_message("None owned!\n"); RETURN FALSE
    torches -= 1
    IF is_in_dungeon() THEN torch_light = 112 + random(0, 15) ELSE torch_light = 240
    RETURN TRUE
ENDFUNCTION
```
