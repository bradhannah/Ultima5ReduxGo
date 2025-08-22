# Town Systems

## Guard Activation and Karma

```pseudocode
FUNCTION activate_guards():
    FOR i = 0 TO 31:
        IF npc_is_active(i) THEN
            t = npc_tile(i)
            IF t is GUARDS OR t is SHADOW_LORD OR t is DEVIL THEN setattack(i)
            ELSE IF random(0, 255) < 128 THEN setflee(i)
        ENDIF
    ENDFOR
ENDFUNCTION

FUNCTION attack_town_tile(x,y):
    tile_obj = looklist(x,y,level)
    IF tile_obj < MONSTERS THEN adjust_karma(-5); activate_guards()
    ELSE IF is_devil(tile_obj) THEN activate_guards()
    SWITCH top_tile_at(x,y):
        CASE BROKEN_MIRROR, BED, STOCKS, MANACLES:
            show_message("Murdered!\n"); adjust_karma(-5); kapow_xy(x,y)
            npc = obj_to_npc_at(x,y); IF npc >= 0 THEN killnpc(npc); outwar(stat[npc].objnum); exterminate_npc(npc)
    ENDSWITCH
ENDFUNCTION
```

## Cannons (Town Fire)

```pseudocode
FUNCTION fire_town_cannon():
    dir = infer_cannon_direction_from_tiles_around_player(); IF dir is NONE THEN show_message("What?\n"); RETURN
    show_message("BOOOM!\n"); play_broadside_sfx(); activate_guards()
    tx, ty = player_x, player_y; sx, sy = 5 + dir.dx, 5 + dir.dy; isdoor = isnpc = FALSE; range = 5
    WHILE range > 0 AND NOT isdoor AND NOT isnpc:
        tx += dir.dx; ty += dir.dy; sx += dir.dx; sy += dir.dy
        obj = looklst2(tx, ty, level)
        IF obj == 0 THEN IF tile_is_door(*look(tx,ty)) THEN isdoor = TRUE
        ELSE IF is_npc_or_valid_target(obj) THEN isnpc = TRUE; idx = object_index_from_looklist()
        range -= 1
    ENDWHILE
    missile(cannon_screen_x, cannon_screen_y, sx, sy, CANNONBALL)
    IF isdoor OR isnpc THEN kapow_xy(tx, ty)
    IF isdoor THEN show_message("Door destroyed!\n"); set_tile(tx, ty, FLOOR); losflag = 1; doortyp = 0
    IF isnpc THEN remove_object(idx); losflag |= 2; adjust_karma(-5); IF (npc = npc_index_from_object(idx)) >= 0 THEN killnpc(npc); exterminate_npc(npc)
    IF isnpc AND idx == 0 THEN damageparty()
ENDFUNCTION
```

## Drawbridges and Portcullis (Night Behavior)

```pseudocode
FUNCTION cache_bridge_tiles_on_map_entry():
    bridge_count = 0
    FOR x FROM 0 TO map_width-1:
        FOR y FROM 0 TO map_height-1:
            t = base_tile_at(x, y)
            IF t == BRIDGE_BASE_1 OR t == BRIDGE_BASE_2 THEN
                bridge_x[bridge_count] = x; bridge_y[bridge_count] = y
                bridge_tile[bridge_count] = tile_at(x, y)
                bridge_count += 1
            ENDIF
        ENDFOR
    ENDFOR
ENDFUNCTION

FUNCTION refresh_drawbridges_by_time():
    IF is_player_on_bridge() THEN RETURN
    IF is_night_time() THEN FOR i = 0 TO bridge_count-1: set_tile(bridge_x[i], bridge_y[i], bridge_tile[i])
    ELSE FOR i = 0 TO bridge_count-1: set_tile(bridge_x[i], bridge_y[i], SHOALS)
ENDFUNCTION
```


### Day/Night Bridge State Summary

| Time  | Bridge Tiles Shown                        |
|-------|-------------------------------------------|
| Day   | `SHOALS` (open water at bridge positions) |
| Night | Original cached `BRIDGE`/portcullis tiles |

Notes:

- The player’s current tile is excluded to avoid toggling beneath the player.
- Bridges are cached on map entry and reused for quick toggling each day/night transition.
