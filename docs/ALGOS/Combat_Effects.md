# Combat Effects

## Poison on Hit (Monsters with `POISONS`)

```pseudocode
FUNCTION maybe_apply_poison_on_monster_hit(attacker_has_POISONS, victim_is_pc, victim_status):
    IF NOT attacker_has_POISONS THEN RETURN
    IF random(0, 3) == 0 THEN RETURN // 1/4 misses, 3/4 applies
    IF victim_is_pc AND victim_status == GOOD_HEALTH THEN
        set_victim_status_poisoned(); announce("<Name> is poisoned!\n"); flag_wound_type_poisoned(); mark_stats_changed()
    ELSE deal_damage_to_victim(getrandom(20)); award_xp_to_attacker_if_pc()
ENDFUNCTION
```

## Fleeing

```pseudocode
FUNCTION update_monster_flee_state(monster_hp, monster_max_hp):
    quarter = monster_max_hp / 4
    IF monster_hp < quarter THEN set_monster_flee_flag(TRUE); RETURN
    half = quarter * 2
    IF monster_hp < half THEN IF getrandom(256) > 251 THEN set_monster_flee_flag(TRUE) ELSE set_monster_flee_flag(FALSE); RETURN
    set_monster_flee_flag(FALSE)
ENDFUNCTION
```

## Field Effects (`infield`)

```pseudocode
FUNCTION apply_in_combat_field_effects(creature_index):
    ptr = combatq[creature_index]; IF ptr.condition == 0 THEN RETURN
    terrain = get_combat_terrain_at(ptr.xpos, ptr.ypos); flag = 0
    IF terrain == LAVA OR terrain == FIREPLACE THEN flag = 100
    ELSE IF terrain == SWAMP THEN flag = 50
    IF flag == 0 THEN
        FOR i = 0 TO 31:
            IF i == ptr.objnum THEN CONTINUE
            IF object[i].xpos == ptr.xpos AND object[i].ypos == ptr.ypos THEN
                IF object[i].tile == FIRE_FIELD THEN flag = 100
                ELSE IF object[i].tile == POISON_FIELD THEN flag = 50
                ELSE IF object[i].tile == SLEEP_FIELD THEN flag = 150
                IF flag != 0 THEN BREAK
            ENDIF
        ENDFOR
    ENDIF
    SWITCH flag:
        CASE 150: put_to_sleep(creature_index)
        CASE 100: play_hit_effect(creature_index); deal_damage(creature_index, getrandom(10)); diagnose(creature_index, NOTHING); mark_stats_changed()
        CASE 50: IF is_pc(creature_index) THEN apply_poison(creature_index); play_hit_effect(creature_index)
    ENDSWITCH
ENDFUNCTION
```

### Field Effects Matrix

| Source           | Effect                               | Applies To        | Notes                                |
|------------------|--------------------------------------|-------------------|--------------------------------------|
| FIRE_FIELD       | Damage `getrandom(10)`               | PCs and monsters  | Also used for fireplaces/lava tiles. |
| POISON_FIELD     | Apply Poisoned status                | PCs only          | Monsters are not poisoned by fields. |
| SLEEP_FIELD      | Apply Sleep status                   | PCs and monsters  |                                      |

## Missiles

```pseudocode
FUNCTION fire_missile(ax, ay, tx, ty, missile_type):
    xline, yline = drawpath((ax*16)+16, (ay*16)+16, (tx*16)+16, (ty*16)+16)
    incr = 13; IF is_town_map() THEN incr = 6; IF missile_type == ROCK THEN incr = 8
    start_offset = 4 IF missile_type == ARROW ELSE 0
    i = start_offset
    WHILE xline[i] != 0xFF:
        findtile(xline[i], yline[i]); IF tempx == -1 THEN BREAK
        animate_missile_step(xline[i], yline[i], missile_type); delay_frame(); erase_missile_trail(xline[i], yline[i])
        i += incr; IF i >= len(xline) OR xline[i] == 0xFF THEN BREAK
        IF blocked(tempx, tempy) AND NOT (tempx==ax AND tempy==ay) THEN RETURN FALSE
    ENDWHILE
    RETURN TRUE
ENDFUNCTION
```

### Town Cannon and Impact Rules

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

FUNCTION attack_town_tile(x,y):
    tile_obj = looklist(x,y,level)
    IF tile_obj < MONSTERS THEN adjust_karma(-5); activate_guards() ELSE IF is_devil(tile_obj) THEN activate_guards()
    SWITCH top_tile_at(x,y):
        CASE BROKEN_MIRROR, BED, STOCKS, MANACLES:
            show_message("Murdered!\n"); adjust_karma(-5); kapow_xy(x,y)
            npc = obj_to_npc_at(x,y); IF npc >= 0 THEN killnpc(npc); outwar(stat[npc].objnum); exterminate_npc(npc)
    ENDSWITCH
ENDFUNCTION
```

## Regurgitation (`throwup`)

```pseudocode
FUNCTION maybe_regurgitate_swallowed(meal_index):
    IF rolld30() < combatq[meal_index].dexterity THEN
        IF is_pc(meal_index) THEN announce(player_name(combatq[meal_index].number)) ELSE announce(monster_name(combatq[meal_index].number))
        show_message(" regurgitated!\n"); play_white_noise()
        clear_flag(combatq[meal_index], UNDER_MASK)
        object[combatq[meal_index].objnum].shape = object[combatq[meal_index].objnum].tile
    ENDIF
ENDFUNCTION
```

## Per-Turn Player Updates

```pseudocode
FUNCTION update_players_after_turn():
    alive_and_awake = 0
    FOR i = 0 TO party_size-1:
        stat = player[i].status
        IF stat == DEAD AND i == active_player THEN active_player = NOTHING
        IF stat != DEAD AND stat != SLEEP THEN
            IF stat == POISONED THEN damageplayer(i,1)
            alive_and_awake++
        ENDIF
    ENDFOR
    IF hour != oldhour THEN
        IF food==0 THEN show_message("Starving!\n"); damageparty()
        ELSE IF hour==6 OR hour==12 OR hour==18 THEN decrement_food(alive_and_awake)
        oldhour = hour
    ENDIF
    increment(last_gave)
    IF spell_dur != 0 AND spell_dur != 255 THEN spell_dur -= 1; IF spell_dur == 0 THEN active_spell = FALSE; update_stats_window()
    check_ring_effects()
ENDFUNCTION

## Field Expiration (fieldkill)

```pseudocode
FUNCTION expire_magical_fields_each_tick():
    FOR i = 0 TO 31:
        IF (object[i].tile & 0xFC) == MAGIC_FIELD THEN
            IF random(0,255) < 16 THEN remove_object(i)
ENDFUNCTION
```

Notes:

- 1/16 chance per tick for any magic field object to vanish. Uses object list tiles with the MAGIC_FIELD family mask.

## Aiming UI (plraim)

```pseudocode
FUNCTION aim_with_crosshair(turn_index, max_range) -> range_or_0:
    crosshair = TRUE
    target = get_default_target_or_self(turn_index)
    xcross, ycross = combatq[target].xpos, combatq[target].ypos
    WHILE TRUE:
        update_screen()
        key = getkey()
        dx, dy = 0, 0
        SWITCH key:
            CASE DIR_NW: dx=-1; dy=-1
            CASE DIR_NE: dx= 1; dy=-1
            CASE DIR_SW: dx=-1; dy= 1
            CASE DIR_SE: dx= 1; dy= 1
            CASE DIR_NORTH: dy=-1
            CASE DIR_SOUTH: dy= 1
            CASE DIR_EAST:  dx= 1
            CASE DIR_WEST:  dx=-1
            CASE ' ': IF xcross==combatq[turn_index].xpos AND ycross==combatq[turn_index].ypos THEN crosshair=FALSE; newline(); RETURN 0
                      // else fallthrough to accept
            CASE ENTER, 'A': IF xcross!=combatq[turn_index].xpos OR ycross!=combatq[turn_index].ypos THEN
                                 crosshair=FALSE; newline(); play_cast_or_aim_sfx(); RETURN distance(xcross, ycross, combatq[turn_index])
            CASE ESC: crosshair=FALSE; newline(); RETURN 0
        END
        newx = xcross + dx; newy = ycross + dy
        IF within_bounds(newx, newy) AND distance(newx, newy, combatq[turn_index]) <= max_range THEN
            xcross, ycross = newx, newy
    ENDWHILE
ENDFUNCTION
```

Notes:

- Space on self aborts; Enter/'A' accepts if not on self. Starts at current target if valid; otherwise self.
- Returns 0 on abort; otherwise returns range to aimed tile.

## Diagnose Messaging

```pseudocode
FUNCTION diagnose_after_hit(target, attacker):
    cond = combatq[target].condition
    IF wound_type has DEFLECT_WND THEN announce_name(target); show_message(" grazed!\n"); small_glide(); CLEAR DEFLECT_WND
    ELSE IF NOT (wound_type has VANISHED_WND) THEN
        IF cond == 0 OR (cond has OBJECT_MASK) THEN announce_name(target); show_message(" killed!\n"); SET KILLED_WND
        ELSE IF wound_type has SLEPT_WND THEN announce_name(target); show_message(" slept!\n")
        ELSE IF NOT (wound_type has POISONED_WND) THEN
            IF cond has PC_MASK THEN
                IF attacker != NOTHING AND combatq[attacker].number == CORPSER THEN announce_name(target); show_message(" dragged under!\n"); small_glide(); SET UNDER_MASK on target; hide_target_shape()
                ELSE announce_name(target); show_message(" hit!\n")
            ELSE
                // Monster wound tiers by remaining HP fraction
                msg = monster_condition_bucket(target)
                SWITCH msg: 4→" barely wounded!", 3→" lightly wounded!", 2→" heavily wounded!", 1→" critical!"
        ENDIF
        IF cond has PC_MASK THEN update_stats_and_screen()
        CLEAR POISONED_WND|SLEPT_WND
    ELSE CLEAR DEFLECT_WND|VANISHED_WND
ENDFUNCTION
```

Notes:

- Corpser special: prints “dragged under!” and sets UNDER_MASK, hiding the target until regurgitation. Matches “throwup”/regurgitation flow.

## Charm Cleanup (“Bad Trip”)

When mass‑charm/confusion leaves no valid enemies to target, the engine applies a cleanup on charmed PCs: remove CHARM, print “passes out!”, unready Chaos Sword, and put the PC to sleep.

```pseudocode
FUNCTION charm_cleanup_if_no_targets():
    FOR i IN 0..5: // party slots
        IF combatq[i].condition has CHARM_MASK AND PC_MASK THEN
            CLEAR CHARM_MASK
            show_message(player_name(combatq[i].number) + " passes out!")
            unready_item_if_equipped(combatq[i].number, CHAOS_SWORD)
            put_pc_to_sleep(i)
            BREAK
    // If none found, no‑op
ENDFUNCTION
```

Notes:

- This behavior is triggered in the AI when no enemies remain (legacy `badtrip`). It ensures a charmed PC doesn’t soft‑lock combat.
## Equipment Resistances & Effects (Display in Ztats)

Certain equipment grants passive effects or resistances that influence combat and hazards. Ztats should display concise tags when equipped items confer such effects.

```pseudocode
// Equipment metadata
STRUCT EquipEffect { id: string, tag: string, applies: function(context) -> bool }

LIST<EquipEffect> EffectsByItem = [
    // Examples; exact items and values are data-driven
    { id: RING_PROTECTION, tag: "Protection", applies: (ctx) => TRUE },
    { id: RING_REGEN,       tag: "Regen",      applies: (ctx) => ctx.in_combat },
    { id: RING_INVIS,       tag: "Invisible",  applies: (ctx) => ctx.in_combat },
    { id: HELM_MAGIC,       tag: "Resist Magic", applies: (ctx) => TRUE },
    { id: SHIELD_MAGIC,     tag: "Resist",     applies: (ctx) => TRUE },
]

FUNCTION equipment_tags_for(player_index, context):
    tags = []
    FOR item IN [player.weapon, player.armor, player.shield, player.helm, player.ring]:
        FOR eff IN EffectsByItem:
            IF item == eff.id AND eff.applies(context) THEN tags.append(eff.tag)
    RETURN tags
```

### Example Item Tags (align with item data)

| Item                  | Slot   | Tag(s)        | Notes                                                                 |
|-----------------------|--------|---------------|-----------------------------------------------------------------------|
| Ring of Protection    | Ring   | Protection    | Small damage reduction bonus; stacks with base armor; engine-defined. |
| Ring of Regeneration  | Ring   | Regen         | Restores a small amount of HP per combat turn.                        |
| Ring of Invisibility  | Ring   | Invisible     | Sets invisible flag; some foes may still detect (specials).           |
| Magic Helm            | Helm   | Resist Magic  | Improves saves or reduces magic damage; data-driven mitigation.       |
| Magic Shield          | Shield | Resist        | Reduces incoming damage (projectiles/magic per design).               |
| Crown of L. British   | Head   | Light, Negate | Persistent effects while equipped/active; see Commands → Use.         |

These are illustrative; enforce exact effects via your item database and balance tables. Ztats should simply display the resulting tags from `EffectsByItem`.

## Chest Traps and Random Traps (`boom`)

```pseudocode
// Random trap effect selection used in various contexts (e.g., traps sprung on party)
FUNCTION trigger_random_trap(target_player_index):
    // Trap distribution: acid 3/8, poison 2/8, bomb 2/8, gas 1/8
    roll = random(0, 7) // 0..7
    SWITCH roll:
        CASE 0,1,2: // 3/8
            show_message("ACID!\n"); damageplayer(target_player_index, rolld30())
        CASE 3,4:   // 2/8
            show_message("POISON!\n"); poison(target_player_index)
        CASE 5,6:   // 2/8
            show_message("BOMB!\n"); damageparty()
        CASE 7:     // 1/8
            show_message("GAS!\n"); FOR i = 0 TO 5: poison(i)
    ENDSWITCH
ENDFUNCTION
```
```
