# Objects and Fixtures

Minor object-specific rules and interactions that affect gameplay.

## Mirrors

```pseudocode
// Reflection on small maps when avatar stands in front of a mirror tile.
FUNCTION update_mirror_reflection_if_avatar_adjacent(tile_pos):
    IF is_avatar_at(tile_pos.down()) THEN set_tile_override(tile_pos, MIRROR_AVATAR)
    ELSE set_tile_override(tile_pos, MIRROR)
ENDFUNCTION

// Breaking a mirror in town.
FUNCTION attack_mirror_at(x, y):
    IF top_tile_at(x, y) == MIRROR THEN
        set_tile(x, y, BROKEN_MIRROR)
        show_message("Broken!\n"); play_glass_break_sfx(); fast_los_update()
    ENDIF
ENDFUNCTION
```

## Pushable Objects (Remove with Spell)

```pseudocode
// Remove pushable objects via a targeted spell effect (e.g., remove clutter).
FUNCTION remove_pushable_object_at(x, y):
    t = top_tile_at(x, y)
    IF t IN {PLANT, CHAIR_FACING_N, CHAIR_FACING_S, CHAIR_FACING_E, CHAIR_FACING_W,
             DESK, BARREL, VANITY, PITCHER, DRAWERS, END_TABLE, FOOTLOCKER, MIRROR} THEN
        set_tile(x, y, FLOOR)
        show_message("POOF!\n"); fast_los_update(); animate_remove_effect()
        RETURN TRUE
    ENDIF
    RETURN FALSE
ENDFUNCTION
```

## Doors (Lock/Unlock)

```pseudocode
// Wizard-lock a door at a targeted location (magical lock).
FUNCTION wizard_lock_door_at(x, y):
    SWITCH top_tile_at(x, y):
        CASE DOOR, LOCKED_DOOR: set_tile(x, y, MAGIC_LOCKED_DOOR); fast_los_update(); RETURN TRUE
        CASE WINDOW_DOOR, WINDOW_LOCK_DOOR: set_tile(x, y, WINDOW_MAG_DOOR); fast_los_update(); RETURN TRUE
    ENDSWITCH
    RETURN FALSE
ENDFUNCTION

// Wizard-unlock a standard locked door.
FUNCTION wizard_unlock_door_at(x, y):
    t = top_tile_at(x, y)
    IF t == LOCKED_DOOR OR t == WINDOW_LOCK_DOOR THEN
        set_tile(x, y, t - 1) // toggle to matching unlocked/windowed state
        fast_los_update(); RETURN TRUE
    ENDIF
    RETURN FALSE
ENDFUNCTION
```

## Furniture and Seating

```pseudocode
// Movement rule: non-foot forms cannot occupy chairs.
FUNCTION is_chair_blocked_for_vehicle(vehicle, terrain):
    IF NOT is_player_on_foot(vehicle) AND base_tile(terrain) == CHAIR THEN RETURN TRUE
    RETURN FALSE
ENDFUNCTION

// Avatar seating visuals (contextual tiles).
FUNCTION get_avatar_seating_tile(chair_tile, pos):
    SWITCH chair_tile:
        CASE CHAIR_FACING_RIGHT: RETURN AVATAR_SITTING_RIGHT
        CASE CHAIR_FACING_LEFT:  RETURN AVATAR_SITTING_LEFT
        CASE CHAIR_FACING_UP:
            up = pos.up(); IF top_tile_at(up) IN {TABLE_FOOD_BOTH, TABLE_FOOD_BOTTOM} THEN RETURN AVATAR_EATING_UP
            RETURN AVATAR_SITTING_UP
        CASE CHAIR_FACING_DOWN:
            down = pos.down(); IF top_tile_at(down) IN {TABLE_FOOD_BOTH, TABLE_FOOD_TOP} THEN RETURN AVATAR_EATING_DOWN
            RETURN AVATAR_SITTING_DOWN
    ENDSWITCH
    RETURN chair_tile
ENDFUNCTION
```

## Grates and Ladders (Klimb)

```pseudocode
// Klimb action routing on small maps.
FUNCTION try_klimb_at_current_tile():
    t = avatar_tile()
    IF t IN {AVATAR_ON_LADDER_DOWN, LADDER_DOWN, GRATE} THEN
        IF can_go_down_one_floor() THEN go_down_one_floor(); show_message("Klimb-Down!"); RETURN TRUE
        ELSE show_message("Can't go lower!\n"); RETURN FALSE
    ENDIF
    IF t IN {AVATAR_ON_LADDER_UP, LADDER_UP} THEN
        IF can_go_up_one_floor() THEN go_up_one_floor(); show_message("Klimb-Up!"); RETURN TRUE
        ELSE show_message("Can't go higher!\n"); RETURN FALSE
    ENDIF
    // Otherwise initiate directional Klimb prompt
    show_message("Klimb-"); prompt_direction_then_klimb()
    RETURN TRUE
ENDFUNCTION
```


### Door State Transitions (Summary)

| Current            | Action            | Result               | Notes                                    |
|--------------------|-------------------|----------------------|------------------------------------------|
| DOOR               | Wizard-Lock       | MAGIC_LOCKED_DOOR    | Magical lock applied                     |
| LOCKED_DOOR        | Wizard-Lock       | MAGIC_LOCKED_DOOR    |                                          |
| WINDOW_DOOR        | Wizard-Lock       | WINDOW_MAG_DOOR      |                                          |
| WINDOW_LOCK_DOOR   | Wizard-Lock       | WINDOW_MAG_DOOR      |                                          |
| LOCKED_DOOR        | Wizard-Unlock     | DOOR                 | Reverts to unlocked                      |
| WINDOW_LOCK_DOOR   | Wizard-Unlock     | WINDOW_DOOR          |                                          |
| LOCKED_DOOR        | Use Key/Jimmy     | DOOR                 | If success; else remains locked          |
| WINDOW_LOCK_DOOR   | Use Key/Jimmy     | WINDOW_DOOR          |                                          |
| MAGIC_LOCKED_DOOR  | Key/Jimmy         | MAGIC_LOCKED_DOOR    | Magical locks ignore keys/jimmy          |
| WINDOW_MAG_DOOR    | Key/Jimmy         | WINDOW_MAG_DOOR      |                                          |

Notes:

- “Wizard-Lock”/“Wizard-Unlock” refer to spell-driven actions (see Spells for details).
- Jimmy/Key success depends on skill/availability and is handled by door-opening routines.
- Windowed variants are purely visual; locking behavior follows the base state.

### Door UI/UX Strings and Sounds

| Event                         | UI Text             | Sound cue    |
|-------------------------------|---------------------|--------------|
| Opened successfully           | "Opened!"           | door open    |
| Locked (standard)             | "Locked!"           | thunk/denied |
| Magically Locked              | "Magically Locked!" | magic deny   |
| Bang to open (not a door)     | "Bang to open!"     | thunk        |
| Jimmy success                 | "Unlocked!"         | small click  |
| Jimmy failure/broken pick     | "Key broke!"        | snap         |
| Not a lock (Jimmy wrong tile) | "Not lock!"         | thunk        |
