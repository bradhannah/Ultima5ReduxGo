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
