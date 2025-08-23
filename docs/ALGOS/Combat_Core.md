# Combat Core

## Hit Calculation (`hit`)

```pseudocode
FUNCTION does_attack_hit(attacker, defender, attack_type):
    IF attack_type is Physical THEN aggressor_stat = attacker.dexterity; victim_stat = defender.dexterity
    ELSE aggressor_stat = attacker.intelligence; victim_stat = defender.intelligence
    target_number = (30 + victim_stat - aggressor_stat) / 2
    dice_roll = rolld30()
    RETURN dice_roll >= target_number
ENDFUNCTION
```

## Damage Calculation (`getdamage`)

```pseudocode
FUNCTION calculate_damage(attacker, defender):
    IF attacker is Monster THEN weapon_damage = attacker.base_damage
    ELSE
        IF attacker.weapon is GlassSword THEN unready(attacker.weapon); RETURN 99
        weapon_damage = random(1, attacker.weapon.attack_value)
    ENDIF
    armor_absorption = defender.armor_value IF defender is Monster ELSE defender.armor_class
    IF armor_absorption > 0 THEN mitigation = random(1, armor_absorption); final_damage = weapon_damage - mitigation
    ELSE final_damage = weapon_damage
    RETURN final_damage
ENDFUNCTION
```

## Turn Modifiers

```pseudocode
FUNCTION monster_turn_should_act(plans):
    IF active_spell == TIME_STOP THEN RETURN FALSE
    IF active_spell == QUICKNESS THEN turn1 = turn1 XOR 0x01; IF turn1 != 0 THEN RETURN FALSE
    IF plans < 0 THEN turn2 = turn2 XOR 0x01; IF turn2 != 0 THEN RETURN FALSE
    IF player_form_is_rider_or_carpet() THEN turn3 = turn3 XOR 0x01; IF turn3 != 0 THEN RETURN FALSE
    RETURN TRUE
ENDFUNCTION
```

## Distance Helpers

```pseudocode
FUNCTION distance_squared(ax, ay, bx, by):
    dx = ax - bx; dy = ay - by
    RETURN dx*dx + dy*dy

FUNCTION distance(ax, ay, bx, by):
    // Integer sqrt via successive odd subtraction (legacy)
    res = distance_squared(ax, ay, bx, by)
    odd = 1; sqrt = 0
    WHILE res >= odd:
        res -= odd; odd += 2; sqrt += 1
    RETURN sqrt

FUNCTION combat_distance(i, j):
    RETURN distance(combatq[i].xpos, combatq[i].ypos, combatq[j].xpos, combatq[j].ypos)
```

Notes:

- Distance is used by AI for target selection and by the aiming UI for range checks.
