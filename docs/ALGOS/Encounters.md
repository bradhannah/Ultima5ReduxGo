# Encounters and Spawning

## Encounter Probability (Overworld)

```pseudocode
FUNCTION should_monster_be_generated():
    probability_score = get_encounter_probability_score()
    dice_roll = random(1, 30)
    RETURN probability_score > dice_roll
ENDFUNCTION

FUNCTION get_encounter_probability_score():
    IF is_in_underworld() THEN RETURN 3
    current_tile = get_player_tile_type()
    IF current_tile is Road THEN base_score = 0
    ELSE IF current_tile is Swamp or Forest or Mountains THEN base_score = 2
    ELSE base_score = 1
    ENDIF
    IF is_night() THEN RETURN base_score + 3
    RETURN base_score
ENDFUNCTION
```

## Monster Selection (`pickmon`)

```pseudocode
FUNCTION select_monster_for_tile(tile_type, is_underworld):
    IF tile_type is Water or River THEN
        IF random(0, 64) < 16 THEN
            IF is_underworld THEN RETURN weighted_random_selection_from(UnderworldWaterMonsters)
            IF tile_type is DeepWater AND random(0, 7) == 7 THEN RETURN Whirlpool
            RETURN weighted_random_selection_from(SurfaceWaterMonsters)
        ENDIF
        RETURN NoMonster
    ENDIF

    IF tile_type is Desert THEN
        IF random(0, 3) == 0 THEN RETURN Sandtrap
        RETURN NoMonster
    ENDIF

    IF tile_type is Swamp AND is_underworld THEN RETURN Rotworm
    IF tile_type is Mountains or Peaks THEN RETURN NoMonster

    IF is_underworld THEN RETURN weighted_random_selection_from(UnderworldLandMonsters)
    RETURN weighted_random_selection_from(SurfaceLandMonsters)
ENDFUNCTION
```

## Weighted Selection (`random_monster`)

```pseudocode
FUNCTION weighted_random_selection_from(frequency_list):
    dice_roll = random(0, 255)
    index = 0
    WHILE dice_roll >= frequency_list[index]:
        dice_roll -= frequency_list[index]
        index += 1
    ENDWHILE
    RETURN get_monster_at(index)
ENDFUNCTION
```

## Placement and Valid Locations

```pseudocode
FUNCTION choose_valid_offscreen_location():
    LOOP FOREVER:
        tempx = (random(0, 31) + xoffset) & 0xFF
        tempy = (random(0, 31) + yoffset) & 0xFF
        IF abs(tempx - player_x) > 6 AND abs(tempy - player_y) > 6 AND
           abs(tempx - player_x) < 0xFA AND abs(tempy - player_y) < 0xFA THEN RETURN
    ENDLOOP
ENDFUNCTION

FUNCTION create_overworld_monster():
    FOR tries = 0 TO 127:
        choose_valid_offscreen_location()
        tx, ty = tempx, tempy
        tile_here = get_tile_at(tx, ty)
        mon = select_monster_for_tile(tile_here, is_underworld = (level != 0xFF))
        IF mon == NONE THEN CONTINUE
        IF mon == PIRATE_SHIP AND is_river(tile_here) THEN CONTINUE
        slot = find_free_object_slot()
        add_object_to_lists(mon, mon, tx, ty, level, 0, slot)
        IF mon == PIRATE_SHIP THEN set_object_hull(slot, 100)
        RETURN TRUE
    ENDFOR
    RETURN FALSE
ENDFUNCTION
```

