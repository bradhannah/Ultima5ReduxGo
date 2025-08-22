# Spells

## Kal Xen (Summon Animal)

```pseudocode
FUNCTION cast_kal_xen():
    play_cast_sound(light)
    roll = random(0, 15)
    IF roll < 6 THEN animal = RATS
    ELSE IF roll < 11 THEN animal = SPIDERS
    ELSE IF roll < 14 THEN animal = BATS
    ELSE animal = SNAKES
    FOR tries = 0 TO 7:
        IF NOT rndxy() THEN CONTINUE
        IF NOT legalmove2(tilenum(RATS), tempx, tempy) THEN CONTINUE
        summoned = createentry(animal, 0, tempx, tempy, level)
        IF summoned < 0 THEN BREAK
        fizz_in_monster(tilenum(animal), tempx, tempy)
        set_flag(combatq[summoned], CHARM_MASK)
        RETURN TRUE
    ENDFOR
    RETURN FALSE
ENDFUNCTION
```

## In Xen Mani (Create Food)

```pseudocode
FUNCTION cast_in_xen_mani():
    play_cast_sound(light)
    food = food + random(1, 3)
    mark_stats_changed()
    RETURN TRUE
ENDFUNCTION
```

## In Vas Por Ylem (Earthquake)

```pseudocode
FUNCTION cast_in_vas_por_ylem(caster_index):
    play_cast_sound(mid); earthquake_visuals()
    FOR target = 0 TO 31:
        IF combatq[target].condition == 0 THEN CONTINUE
        IF loyalty(target) == 0 THEN CONTINUE
        IF getattr(target, DEX) <= rolld30() THEN
            explosion_effect_on(target)
            dmg = random(1, 20)
            add_caster_xp(caster_index, apply_damage(target, dmg))
            diagnose(target, caster_index)
        ENDIF
    ENDFOR
ENDFUNCTION
```

## Mani (Heal)

```pseudocode
FUNCTION cast_mani():
    plr = select_target_party_member(); IF plr < 0 THEN RETURN -1
    IF player[plr].status == DEAD THEN RETURN FALSE
    heal_amount = rolld30(); player[plr].hit_points = min(player[plr].hit_points + heal_amount, player[plr].hit_max)
    mark_stats_changed(); play_cast_sound(light)
    RETURN TRUE
ENDFUNCTION
```

## Kal Xen Corp (Summon Daemon)

```pseudocode
FUNCTION cast_kal_xen_corp(is_scroll, caster_index):
    play_cast_sound(is_scroll ? mid : high)
    FOR tries = 0 TO 7:
        IF NOT rndxy() THEN CONTINUE
        IF NOT legalmove2(tilenum(DEVIL), tempx, tempy) THEN CONTINUE
        IF tile_at(tempx, tempy) == NOTHING THEN CONTINUE
        devil_idx = createentry(DEVIL, 0, tempx, tempy, level)
        IF devil_idx < 0 THEN BREAK
        fizz_in_monster(tilenum(DEVIL), tempx, tempy)
        IF NOT is_scroll AND rolld30() >= getattr(caster_index, INT) THEN show_message("Oops...\n"); RETURN -1
        set_flag(combatq[devil_idx], CHARM_MASK); RETURN TRUE
    ENDFOR
    RETURN FALSE
ENDFUNCTION
```

