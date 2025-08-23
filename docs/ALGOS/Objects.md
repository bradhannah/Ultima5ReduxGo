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

## Jimmy (Lockpicks) and Skull Keys

### Jimmy — Towns/Small Maps

```pseudocode
FUNCTION jimmy_town():
    IF keys == 0 THEN show_message("No Keys!\n"); RETURN
    IF NOT getdir() THEN RETURN // requires direction
    tx = player_x + tempx; ty = player_y + tempy
    tile = top_tile_at(tx, ty)
    SWITCH tile:
        CASE LOCKED_DOOR, WINDOW_LOCK_DOOR:
            plr = select_character(); IF plr == NONE THEN RETURN
            IF player[plr].dex <= random(0, 29) THEN
                show_message("Key broke!\n"); keys -= 1; RETURN
            ENDIF
            set_tile(tx, ty, tile - 1) // becomes DOOR or WINDOW_DOOR
            fast_los_update(); show_message("Unlocked!\n"); RETURN
        CASE MAGIC_LOCKED_DOOR, WINDOW_MAG_DOOR:
            show_message("Key broke!\n"); keys -= 1; RETURN // cannot jimmy magical locks
        CASE STOCKS, MANACLES:
            IF no_object_present_at(tx, ty) THEN show_message("No one is there!\n"); RETURN
            npc_obj = get_object_at(tx, ty); plr = select_character(); IF plr == NONE THEN RETURN
            IF player[plr].dex <= random(0, 29) THEN show_message("Key broke!\n"); keys -= 1; RETURN
            IF is_town_map() THEN
                npc = obj_to_npc(npc_obj); clear_tlk(npc); if alive(npc) THEN set_follow_schedule(npc); show_message("\n\"I thank thee!\"\n"); inc_karma(2)
                free_npc(npc)
            ELSE
                set_tile(tx, ty, FLOOR); fast_los_update(); show_message("Unlocked\n")
            ENDIF
            RETURN
        DEFAULT:
            obj = find_object_at(tx, ty, CHEST)
            IF obj != NONE THEN jimmy_chest_town(obj)
            ELSE show_message("No lock!\n")
    ENDSWITCH
ENDFUNCTION
```

```pseudocode
FUNCTION jimmy_chest_town(chest_obj_index):
    plr = select_character(); IF plr == NONE THEN RETURN
    num = object[chest_obj_index].number // bit7 indicates locked bit
    IF num < 128 THEN // not locked
        show_message("Key broke!\n"); delay_glide(); keys -= 1; RETURN
    ENDIF
    difficulty = (30 + (num & 0x7F) - player[plr].dex) / 2
    IF random(1, 30) > difficulty THEN
        show_message("Success!\n"); object[chest_obj_index].number &= 0x7F // clear locked bit
    ELSE
        show_message("Key broke!\n"); delay_glide(); keys -= 1
    ENDIF
ENDFUNCTION
```

### Jimmy — Dungeons

```pseudocode
FUNCTION jimmy_dungeon():
    plr = select_character(); IF plr == NONE THEN RETURN
    tx = player_x; ty = player_y; tile = dng_tile_at(level, ty, tx)
    dex = player[plr].dex; difficulty = (30 + 2*level - dex) / 2
    IF is_normal_chest(tile) THEN
        IF keys == 0 THEN show_message("No keys!\n") ELSE show_message("Key broke!\n"); keys -= 1
    ELSE IF is_trapped_chest(tile) THEN
        IF keys == 0 THEN show_message("No keys!\n")
        ELSE IF random(1, 30) > difficulty THEN
            show_message("Chest unlocked\n"); set_dungeon_chest_unlocked(level, ty, tx, tile)
        ELSE show_message("Key broke!\n"); keys -= 1
    ELSE IF is_opened_chest(tile) THEN show_message("Already open!\n")
    ELSE show_message("What?\n")
ENDFUNCTION
```

### Skull Key — Magical Unlock

```pseudocode
FUNCTION use_skull_key():
    skull_keys -= 1
    show_message("Skull Key\n")
    IF is_outdoors() OR is_dungeon() THEN
        IF NOT spelldir() THEN RETURN // sets tempx, tempy toward target
        IF wizard_unlock_magic_at(tempx, tempy) THEN
            IF is_town_map() THEN kapow_xy(tempx, tempy) // visual pop on towns
        ELSE
            // No effect
        ENDIF
    ELSE
        show_message("Not here!\n")
    ENDIF
ENDFUNCTION

FUNCTION wizard_unlock_magic_at(x, y):
    tile = top_tile_at(x, y)
    SWITCH tile:
        CASE MAGIC_LOCKED_DOOR: set_tile(x, y, DOOR); fast_los_update(); RETURN TRUE
        CASE WINDOW_MAG_DOOR:   set_tile(x, y, WINDOW_DOOR); fast_los_update(); RETURN TRUE
        DEFAULT: RETURN FALSE
    ENDSWITCH
ENDFUNCTION
```

Notes:

- Jimmy success on doors/stocks uses a Dex vs random(0,29) check in towns; for chests, use the difficulty formulas shown.
- Magical locks cannot be opened by normal keys/jimmy; only Skull Key (or In Ex Por spell) can unlock magic-locked doors.
- Skull Keys are usable outdoors and in dungeons (not in towns); they consume one key on use regardless of success.

### Jimmy Outcomes Matrix (At-a-Glance)

| Context  | Target                         | Success Condition                              | On Success                                      | On Failure             | Key Consumed? |
|----------|--------------------------------|------------------------------------------------|-------------------------------------------------|------------------------|---------------|
| Town     | LOCKED_DOOR / WINDOW_LOCK_DOOR | Dex > random(0,29) (selected character)        | Tile → DOOR / WINDOW_DOOR; “Unlocked!”          | “Key broke!”           | Yes (−1)      |
| Town     | MAGIC_LOCKED_DOOR (windowed)   | — (cannot jimmy)                               | —                                               | “Key broke!”           | Yes (−1)      |
| Town     | STOCKS / MANACLES              | Dex > random(0,29)                             | Free NPC; “I thank thee!”; Karma +2 (town)      | “Key broke!”           | Yes (−1)      |
| Town     | Chest (locked; object list)    | random(1,30) > (30 + (num&0x7F) − Dex)/2       | Clear locked bit; “Success!”                    | “Key broke!”           | Yes (−1)      |
| Town     | Chest (unlocked; num<128)      | —                                              | —                                               | “Key broke!”           | Yes (−1)      |
| Town     | No lock                        | —                                              | —                                               | “No lock!”             | No            |
| Dungeon  | Normal chest ((tile&0xF7)==0x40) | — (cannot unlock)                             | —                                               | “Key broke!” if keys>0; else “No keys!” | Yes if >0 |
| Dungeon  | Trapped chest ((tile&0xF0)==0x40) | random(1,30) > (30 + 2*level − Dex)/2        | “Chest unlocked”; tile updated to unlocked      | “Key broke!”           | Yes (−1)      |
| Dungeon  | Opened chest ((tile&0xF0)==0x70) | —                                             | “Already open!”                                  | —                      | No            |
| Dungeon  | Other                           | —                                             | —                                               | “What?”                | No            |

Notes:

- Town doors/stocks require a direction (Jimmy prompt) and a selected character for the Dex check.
- Magical locks (town) cannot be jimmy’d; only Skull Key/“In Ex Por” can unlock them.
- Dungeon “normal chest” jimmy always breaks a key if you have one; trapped chest uses the level-based difficulty check.

### Skull Key Outcomes Matrix (At-a-Glance)

| Context  | Target (directional)                 | On Success                                       | On Failure/Invalid Target          | UI Texts                | Key Consumed? |
|----------|--------------------------------------|--------------------------------------------------|------------------------------------|-------------------------|---------------|
| Outdoors | MAGIC_LOCKED_DOOR / WINDOW_MAG_DOOR  | Unlocks to DOOR / WINDOW_DOOR; kapow visual      | No effect                          | "Skull Key"            | Yes (always)  |
| Dungeon  | MAGIC_LOCKED_DOOR / WINDOW_MAG_DOOR  | Unlocks to DOOR / WINDOW_DOOR                    | No effect                          | "Skull Key"            | Yes (always)  |
| Town     | (any)                                 | —                                                | "Not here!"                        | "Skull Key", "Not here!" | Yes (always)  |

Notes:

- Direction is required (spell-targeting) to select a tile; if no direction is given, the key is still consumed (no effect).
- Only magic-locked doors can be opened by Skull Keys; normal locks use Jimmy/keys.
- Outdoors success triggers a small kapow effect; dungeons do not.

### Open (No Key) Outcomes Matrix (At-a-Glance)

The standard Open command attempts to open what’s in the chosen direction without consuming keys.

| Target Tile                 | Outcome           | UI Text              | Notes                                 |
|----------------------------|-------------------|----------------------|---------------------------------------|
| DOOR / WINDOW_DOOR         | Opened            | "Opened!"           | Door changes state as needed          |
| LOCKED_DOOR / WINDOW_LOCK_DOOR | Locked       | "Locked!"           | Requires Jimmy/key or An Sanct        |
| MAGIC_LOCKED_DOOR / WINDOW_MAG_DOOR | Magic lock | "Magically Locked!" | Requires Skull Key or In Ex Por       |
| Not a door (walls, etc.)   | Not applicable    | "Bang to open!"     | Feedback when target isn’t a door     |

Notes:

- This flow does not consume keys; it only reports the door’s state and opens if already available to open.
- Use Jimmy for standard locks (consumes keys). Use Skull Key or In Ex Por for magic locks.
