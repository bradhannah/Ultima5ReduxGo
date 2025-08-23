# Movement — Overworld

## Movement Permissions

```pseudocode
FUNCTION iswater(terrain):
    RETURN terrain < SWAMP OR (terrain & 0xF0) == RIVER
ENDFUNCTION

FUNCTION onfoot(vehicle, terrain):
    allowed = NOT forbidden_foot_terrain(terrain)
    IF NOT is_player_on_foot(vehicle) AND base_tile(terrain) == CHAIR THEN allowed = FALSE
    RETURN allowed
ENDFUNCTION

FUNCTION legalmove(vehicle, terrain):
    mode = movement_mode_from_vehicle(vehicle)
    SWITCH mode:
        CASE 0: RETURN onfoot(vehicle, terrain)
        CASE 1: RETURN iswater(terrain)
        CASE 2: RETURN ((terrain & 0xF0) == RIVER) OR iswater(terrain) OR onfoot(vehicle,terrain)
        CASE 3: RETURN onfoot(vehicle,terrain) AND terrain != LAVA AND terrain != SWAMP
        CASE 4: RETURN NOT iswater(terrain)
        CASE 5:
            IF base_tile(terrain) == GRASS_WATER THEN RETURN coast_mask_allows_entry(terrain, vehicle)
            IF iswater(terrain) THEN
                IF terrain < RIVER THEN RETURN TRUE
                RETURN river_mask_allows_entry(terrain, vehicle)
            ENDIF
            RETURN FALSE
        CASE 6: RETURN terrain <= SEA
        CASE 7: RETURN terrain == SWAMP
        CASE 8: RETURN terrain == GRASS
        CASE 9: RETURN terrain == DEEP
        CASE 10: RETURN terrain == DESERT
        DEFAULT: RETURN FALSE
    ENDSWITCH
ENDFUNCTION
```

## Terrain-Based Movement Chance

```pseudocode
FUNCTION should_move_this_tick(tile):
    IF tile_in_class(tile, HEAVY_TERRAIN) THEN RETURN random(0, 2) == 2 // 1-in-3
    IF tile_in_class(tile, DIFFICULT_TERRAIN) THEN RETURN random(0, 1) == 0 // 1-in-2
    RETURN TRUE
ENDFUNCTION
```

## Whirlpools and Pirate Ships

```pseudocode
FUNCTION overworld_move_monster(index):
    mon = object[index].tile
    IF is_whirlpool(mon) THEN
        object[index].number = object[index].number XOR 0x01
        IF object[index].number == 0 THEN RETURN
        IF random(0, 1) == 0 THEN rndmove(index) ELSE seek(index, ATTRACT_TOWARD_PLAYER)
        RETURN
    ENDIF
    IF is_pirate_ship(mon) THEN
        IF winds == 0 THEN RETURN
        shipdir = get_ship_dir(mon)
        speed = spdmtrx[shipdir][winds - 1]
        object[index].misc = object[index].misc + 1
        IF speed != MAX_SPEED AND object[index].misc > speed THEN object[index].misc = 0; RETURN
        seek(index, ATTRACT_TOWARD_PLAYER)
        RETURN
    ENDIF
    seek(index, ATTRACT_TOWARD_PLAYER)
ENDFUNCTION

FUNCTION rndmove(index):
    tx, ty = object[index].xpos, object[index].ypos
    FOR i = 0 TO 2:
        dir = random(0, 3)
        SWITCH dir:
            CASE 0: IF chkmove(index, tx, ty-1) THEN moveit(index, 0,-1); RETURN
            CASE 1: IF chkmove(index, tx+1, ty) THEN moveit(index, 1, 0); RETURN
            CASE 2: IF chkmove(index, tx, ty+1) THEN moveit(index, 0, 1); RETURN
            CASE 3: IF chkmove(index, tx-1, ty) THEN moveit(index,-1, 0); RETURN
        END
    ENDFOR
ENDFUNCTION

FUNCTION moveit(index, dx, dy):
    tx = object[index].xpos + dx; ty = object[index].ypos + dy
    terrain = get_tile_at(tx, ty)
    IF is_heavy_terrain(terrain) AND random(0, 2) != 2 THEN RETURN
    object[index].xpos = tx; object[index].ypos = ty
    mark_fast_los_update()
    IF tile_at(tx, ty) == MOONGATE THEN object[index].tile = 0; object[index].shape = 0
ENDFUNCTION
```

## Rough Seas and Waterfalls

```pseudocode
FUNCTION apply_rough_seas_if_applicable(current_tile):
    IF current_tile == DEEP AND (player_is_in_skiff() OR player_is_on_carpet()) THEN
        show_message("Rough seas!\n"); kapow_xy(player_x, player_y); hit2()
    ENDIF
ENDFUNCTION

FUNCTION handle_waterfall_collision():
    show_message("F-A-L-L-S!!!\n"); nudge_player_down(); wait_update(1); nudge_player_down(); play_fall_glide_effect()
    old_form = player_form; set_player_form(ON_FOOT); wait_update(1)
    FOR i = 0 TO party_size-1: IF player[i].status != DEAD AND rolld30() >= player[i].dexterity THEN damageplayer(i, 1)
    wait_update(2); set_player_form(old_form)
    IF at_underworld_fall_spot(player_x, player_y) THEN show_message("Falling into underworld!!\n"); enter_underworld_via_waterfall()
ENDFUNCTION
```


### Terrain Classification (Example)

| Class             | Tiles (example)        | Move Chance    | Notes                            |
|-------------------|------------------------|----------------|----------------------------------|
| HEAVY_TERRAIN     | Deep forest, mountains | 1-in-3 (33%)   | Use `random(0,2) == 2`           |
| DIFFICULT_TERRAIN | Swamps, light forest   | 1-in-2 (50%)   | Use `random(0,1) == 0`           |
| OPEN_TERRAIN      | Grass, roads           | Always         | No throttle                      |

Note: Exact tile-class assignments should be data-driven and can vary per map if desired.

### Pirate Ship Speed Matrix (Template)

Pirate ship movement frequency is determined by a speed matrix indexed by the ship’s heading and the global wind. Lower numbers indicate more frequent movement; a tier value of 4 stalls movement on that tick.

| Heading \ Wind | North | South | East | West |
|----------------|:-----:|:-----:|:----:|:----:|
| Up             |   2   |   4   |  3   |  3   |
| Right          |   3   |   3   |  2   |  4   |
| Down           |   4   |   2   |  3   |  3   |
| Left           |   3   |   3   |  4   |  2   |

Notes:

- Example values shown; actual matrix should be supplied by data.
- In code, a non-4 speed gates movement by incrementing a per-ship counter and moving when the counter exceeds the speed; 4 means skip movement.

## Wind System

Wind affects sailing and pirate ship behavior and is visible in the UI.

```pseudocode
FUNCTION set_wind(dir):
    // dir: 1=N, 2=S, 3=E, 4=W, 0=calm, -1=display only
    IF dir != -1 THEN winds = dir; wind_update_counter = 0
    show_wind_ui(winds) // prints “Calm/ North/ South/ East/ West Winds” with arrow indicators

FUNCTION update_wind_tick():
    // 1/64 chance per tick to change wind
    IF random(0, 0x3F) == 0 THEN
        REPEAT
            new_wind = random(0, 4) // 0..4 inclusive
        UNTIL NOT (new_wind == 0 AND random(0, 0xFF) < 0xC0) // bias: calm only 25% when picked
        set_wind(new_wind)
    ENDIF
ENDFUNCTION
```

Notes:

- `winds` drives pirate ship speed via the matrix above and can influence player ship behavior (e.g., auto-sailing).
- Calm (0) is rarer: when the wind changes, calm is accepted only 25% of the times it is rolled.
- The “Rel Hur” spell (Change Wind) can set wind by direction; see Spells → Rel Hur.
