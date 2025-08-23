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

## Special Guard Behavior

Town guards react to crimes and alarms, switch into an aggressive “attack” schedule, and pursue the Avatar. Non‑guard civilians may flee.

```pseudocode
// Global alarm state per small map
STATE guard_alarm_active = FALSE

FUNCTION raise_guard_alarm(reason):
    guard_alarm_active = TRUE
    activate_guards() // sets guards to attack; some civilians flee
    log_alarm_reason(reason)
ENDFUNCTION

FUNCTION clear_guard_alarm():
    guard_alarm_active = FALSE
ENDFUNCTION

// Conversation callback hook from Talk system
FUNCTION CallGuards():
    raise_guard_alarm("Conversation callback")
ENDFUNCTION

// Conversation-driven arrest (teleport to jail cell)
FUNCTION GoToJail():
    send_avatar_to_jail_cell() // See Talk System: GoToJail
    clear_guard_alarm()
ENDFUNCTION

// Per-tick pursuit on small maps when alarm is active
FUNCTION guard_pursuit_tick():
    IF NOT guard_alarm_active THEN RETURN
    FOR each npc IN active_npcs_on_map():
        IF npc_is_guard(npc) THEN
            setattack(npc.index) // ensure attack schedule persists
            pursue_avatar_with_smallmap_ai(npc)
        ELSE IF civilian_should_flee(npc) THEN
            setflee(npc.index)
        ENDIF
    ENDFOR
ENDFUNCTION

// Triggers that raise the alarm
FUNCTION maybe_raise_alarm_on_surface_action(action, x, y, obj):
    SWITCH action:
        CASE AttackNPC:
            adjust_karma(-5); raise_guard_alarm("Assault")
        CASE FireCannon:
            raise_guard_alarm("Cannon fired")
        CASE KillBoundNPC:
            adjust_karma(-5); raise_guard_alarm("Murder")
        CASE SummonDevilOrShadowlord:
            raise_guard_alarm("Hostile presence")
        CASE TheftMajor:
            adjust_karma(-1); maybe_raise_alarm_per_location_rules("Theft")
    ENDSWITCH
ENDFUNCTION

// Reset on leaving the map
FUNCTION on_smallmap_exit():
    clear_guard_alarm()
ENDFUNCTION
```

Behavior notes:

- Triggers: attacking town NPCs, killing restrained NPCs, firing town cannons, or the presence of a Devil/Shadowlord raise the alarm immediately.
- Activation: guards switch to “attack” behavior and path toward the Avatar; other NPCs may switch to “flee”. The alarm persists until leaving the map; some locations may calm after an extended period without further aggression.
- Conversation: scripts may explicitly call guards (CallGuards) or send the Avatar to jail (GoToJail); see Talk System docs for details.
- Karma: violent acts reduce Karma (typical −5); minor theft can reduce Karma slightly (−1) without necessarily raising the alarm, depending on location rules.
- Cannons: firing a town cannon always raises the alarm; killing civilians with a cannon also reduces Karma and removes the NPC.

## Jail Flow

When conversations or scripted events send the Avatar to jail, the engine performs a deterministic relocation into the map’s jail cell and sets up the cell state.

```pseudocode
// Map-provided jail configuration
STRUCT JailConfig {
    cell_x: int
    cell_y: int
    door_x: int
    door_y: int
    door_tile_locked: Tile // e.g., LOCKED_DOOR or WINDOW_LOCK_DOOR
    guard_spawn_points: [Position]
}

FUNCTION send_avatar_to_jail_cell():
    cfg = get_jail_config_for_current_location(); IF cfg == NONE THEN RETURN FALSE
    // Place party in cell (Avatar underfoot; party members co-located or fanned as per engine rules)
    teleport_party_to(cfg.cell_x, cfg.cell_y)
    // Ensure the cell door is closed and locked
    set_tile(cfg.door_x, cfg.door_y, cfg.door_tile_locked); force_los_update()
    // Spawn or retarget nearby guards to watch posts
    FOR p IN cfg.guard_spawn_points: ensure_guard_on_post(p)
    // Clear current alarm; jail is considered a contained state
    clear_guard_alarm()
    show_message("To the cells with thee!\n")
    RETURN TRUE
ENDFUNCTION

// Escape mechanics are shared with door systems (Open/Jimmy/Skull Key/Spells)
FUNCTION try_escape_jail():
    // Player may:
    // 1) Use Skull Key (directional) on magic/locked door
    // 2) Jimmy a non-magical locked door (chance to break key)
    // 3) Cast An Sanct (disarm/unlock) if allowed
    // 4) Obtain a physical key from a container/guard (map-dependent)
    // Doors Overview documents state transitions and UI strings
    RETURN
ENDFUNCTION
```

Jail Effects & Options:

- Location: Each town/keep may define one or more jail cells via a `JailConfig` (door position and cell tile).
- Inventory: No special inventory confiscation occurs by default; any penalties are enforced by the triggering event itself (e.g., extortion callbacks handle gold loss).
- Time: World time continues to advance normally. The player can pass turns or camp subject to local rules.
- Escape: Locked doors follow the standard systems (Open/Jimmy/Skull Key/Spells). Persistent magical locks require Skull Key or a proper spell to remove.
- Alarm: Being placed in jail clears the current guard alarm; leaving the cell does not automatically re-trigger the alarm unless a new crime is committed.
- Guards: Guards may be positioned at watch posts; they will pursue on sight if the player escapes and re-raises the alarm.

See also: Doors state transitions and UI strings (Doors Overview), and Commands → Use/Open/Jimmy for unlocking behavior.

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
