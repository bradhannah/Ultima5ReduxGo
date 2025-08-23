# Commands

## Yell

The Yell command accepts a typed word/phrase and triggers context-specific effects.

```pseudocode
FUNCTION command_yell():
    // Ship sail control (small maps / towns & dungeons)
    IF player_form_is_ship() AND is_small_map():
        IF ship_sails_unfurled() THEN
            show_message("FURL!\n"); furl_sails()
        ELSE
            show_message("HOIST!\n"); hoist_sails()
        ENDIF
        RETURN // finishes immediately
    ENDIF

    show_message("what?\n:")
    s = read_line_max(30)
    IF s == "" THEN show_message("Nothing\n"); RETURN
    newline()

    IF is_town_map() THEN
        end_turn = yell_town(s)
    ELSE IF is_overworld() THEN
        yell_overworld(s)
    ELSE
        show_message("\nNo effect!\n")
    ENDIF
ENDFUNCTION
```

### Yell — Towns (Summon Shadowlord)

```pseudocode
FUNCTION yell_town(s) -> end_turn:
    words = ["FAULINEI", "ASTAROTH", "NOSFENTOR"] // castle-specific
    IF current_map_is_one_of(castle_true, castle_love, castle_courage):
        j = index_of_matching_word(words, s)
        IF j == NONE OR near_map_boundary() OR shadowlord_is_dead(j) THEN
            show_message("\nNo effect!\n"); RETURN TRUE
        ENDIF
        IF shadowlord_already_present() THEN show_message("\nNo effect!\n"); RETURN TRUE
        last_called_shadowlord = j
        spawn_shadowlord_at(player_x, player_y - 2)
        set_shadowlord_attack_schedule(target=player_position())
        show_message("\nA shadowlord appears\n"); play_summon_fx()
        fizz_in_visual_at(5, 3)
        RETURN FALSE // delay aggression one turn
    ELSE
        show_message("\nNo effect!\n"); RETURN TRUE
    ENDIF
ENDFUNCTION
```

Notes:

- Only one shadowlord may be present at a time; repeated calls print “No effect!”.
- Summon position is two tiles north of the player.
- Which castle accepts which word is encoded by map and the words array.

### Yell — Overworld (Words of Power)

```pseudocode
FUNCTION yell_overworld(s):
    vocab = [ "FALLAX", "VILIS", "INOPIA", "MALUM",
              "AVIDUS", "INFAMA", "IGNAVUS", "VERAMOCOR" ]
    dng_tile_ids = [ DNG_0x18, DNG_0x16, DNG_0x16, DNG_0x18,
                     DNG_0x18, DNG_0x17, DNG_0x17, DNG_0x16 ]

    j = index_of_matching_word(vocab, s)
    IF j == NONE THEN show_message("\nNo effect!\n"); RETURN

    show_message("\nA word of power is uttered\n"); play_earthquake_fx()

    // Check the four adjacent tiles in this order: left, down, right, up
    for dir in [Left, Down, Right, Up]:
        tx, ty = adjacent_in(dir)
        tile = tile_at(tx, ty)
        IF tile == dng_tile_ids[j] OR tile == SEALED_DUNGEON OR tile == RUINED_SHRINE THEN
            dx, dy = dir_to_delta(dir)
            BREAK
        ENDIF
    END FOR
    IF dx,dy not set THEN show_message("\nNo effect!\n"); RETURN // proper site not nearby

    IF tile_at(player_x+dx, player_y+dy) == RUINED_SHRINE THEN
        restore_shrine(j, player_x+dx, player_y+dy); RETURN
    ENDIF

    // Toggle dungeon entrance seal at the matched location
    toggle_open_dungeon_flag(j)
    toggle_tile_at(player_x+dx, player_y+dy, between(dng_tile_ids[j], SEALED_DUNGEON))
    fast_los_update()
ENDFUNCTION
```

Notes:

- Each word maps to a specific dungeon entrance tile id; the call toggles the seal state at the matching adjacent location.
- If the adjacent location is a ruined shrine, Yell triggers the shrine restoration flow instead.
- An earthquake visual/sound is played on successful invocation.

### Yell — Ship Sails

```pseudocode
FUNCTION ship_sails_unfurled():
    RETURN player_form_in_ship_family() AND ship_form_is_unfurled()

FUNCTION furl_sails():
    set_player_form_to_ship_furled()

FUNCTION hoist_sails():
    set_player_form_to_ship_unfurled()
ENDFUNCTION
```


### Yell Words Reference (Towns)

Each castle responds to a specific name and summons a Shadowlord two tiles north of the avatar. Use data to bind map IDs to words.

| Castle (virtue)      | Word        | Effect                         |
|----------------------|-------------|---------------------------------|
| Castle of Truth      | FAULINEI    | Summons Shadowlord (Falsehood)  |
| Castle of Love       | ASTAROTH    | Summons Shadowlord (Hatred)     |
| Castle of Courage    | NOSFENTOR   | Summons Shadowlord (Cowardice)  |

Notes:

- Only valid at the corresponding castle map. Summoning while one is already present prints “No effect!”.
- If the corresponding Shadowlord has been defeated for good, yelling the word prints “No effect!”.

### Words of Power Reference (Overworld)

Use data to bind each word to the matching dungeon entrance tile or sealed shrine.

| Word       | Target Adjacent Tile                | Action                                  |
|------------|-------------------------------------|-----------------------------------------|
| FALLAX     | Dungeon entrance tile (word-linked) | Toggle open/sealed at the matched tile   |
| VILIS      | Dungeon entrance tile (word-linked) | Toggle open/sealed at the matched tile   |
| INOPIA     | Dungeon entrance tile (word-linked) | Toggle open/sealed at the matched tile   |
| MALUM      | Dungeon entrance tile (word-linked) | Toggle open/sealed at the matched tile   |
| AVIDUS     | Dungeon entrance tile (word-linked) | Toggle open/sealed at the matched tile   |
| INFAMA     | Dungeon entrance tile (word-linked) | Toggle open/sealed at the matched tile   |
| IGNAVUS    | Dungeon entrance tile (word-linked) | Toggle open/sealed at the matched tile   |
| VERAMOCOR  | Dungeon entrance tile (word-linked) | Toggle open/sealed at the matched tile   |

Notes:

- When the adjacent target is a ruined shrine tile, the shrine restoration flow is invoked instead of toggling a dungeon seal.
- The four-adjacent check order is: left, down, right, up.

## Look

### Look — Dungeon

```pseudocode
FUNCTION command_look_dungeon():
    plr = select_character(); IF plr == NONE THEN RETURN
    IF torch_light == 0 AND magic_light == 0 THEN show_message("You see:\ndarkness.\n"); RETURN
    IF NOT dir_dng(face) THEN RETURN // sets tempx, tempy ahead of avatar
    tx, ty = tempx, tempy
    tile = dng_tile_at(level, ty, tx)

    show_message("You see:\n")
    IF tile == 0x61 THEN tile = 0 // special fixup

    IF (tile & 0xF0) == 0x80 THEN
        SWITCH tile:
            CASE 0x80: show_message("A sleep field.\n")
            CASE 0x81: show_message("A poison gas field.\n")
            CASE 0x82: show_message("A wall of fire.\n")
            CASE 0x83: show_message("An electric field.\n")
            DEFAULT:   show_message("An energy field.\n")
        END
    ELSE IF (tile & 0xF0) == 0xC0 THEN
        type = dngtype & 0x0F
        IF type == 0x01 THEN show_message("a dripping stalactite.\n")
        ELSE IF type == 0x02 THEN show_message("a caved in passage.\n")
        ELSE IF random(1,255) == 255 THEN show_message("an unfortunate software pirate.\n")
        ELSE show_message("a less fortunate adventurer.\n")
    ELSE
        SWITCH (tile & 0xF0):
            CASE 0x00: show_message("a passage.\n")
            CASE 0x10: show_message("an up ladder.\n")
            CASE 0x20: show_message("a down ladder.\n")
            CASE 0x30: show_message("a ladder.\n")
            CASE 0x40: show_message("a wooden chest.\n")
            CASE 0x50: show_message("a fountain.\n")
            CASE 0x60: show_message("a pit.\n")
            CASE 0x70: show_message("an open chest.\n")
            CASE 0x80: show_message("an energy field.\n")
            CASE 0x90: show_message("nothing of note.\n")
            CASE 0xA0: show_message("a heavy door.\n")
            CASE 0xB0: show_message("a wall.\n")
            CASE 0xC0: show_message("SPEC WALL ERR.\n")
            CASE 0xD0: show_message("a wall.\n")
            CASE 0xE0: show_message("a heavy door.\n")
            CASE 0xF0: show_message("a heavy door.\n")
        END
        IF (tile & 0xF0) == 0x50 THEN // Fountain interaction
            show_message("Will you drink?\n"); ch = read_yes_no()
            IF ch == 'N' THEN show_message("No.\n") ELSE
                show_message("Yes.  Gulp!\n")
                SWITCH tile:
                    CASE 0x50: show_message("Cured!\n"); player[plr].status = GOOD
                    CASE 0x51: show_message("Healed!\n"); player[plr].hp = player[plr].hp_max
                    CASE 0x52: show_message("Poisoned!\n"); player[plr].status = POISONED
                    DEFAULT:   show_message("Bad taste.\n"); damageplayer(plr, random(0,7))
                END
            ENDIF
        ENDIF
    ENDIF
ENDFUNCTION
```

### Look — Towns/Overworld

```pseudocode
FUNCTION command_look_surface():
    show_message("Look-")
    IF NOT getdir() THEN RETURN
    tx = player_x + tempx; ty = player_y + tempy
    obj  = *look(tx,ty)
    obj1 = looklist(tx,ty,level) // object on tile (if any)

    // Signs: trace back to holder to reach actual sign tile
    WHILE obj IN {SIGN_HOLDER_N, SIGN_HOLDER_E, SIGN_HOLDER_W}:
        SWITCH obj:
            CASE SIGN_HOLDER_N: ty -= 1
            CASE SIGN_HOLDER_E: tx += 1
            CASE SIGN_HOLDER_W: tx -= 1
        END
        obj = *look(tx,ty)
    ENDWHILE

    show_message("\nThou dost see\n")

    // Special fixtures
    IF obj == TELESCOPE THEN astronomy()
    ELSE IF obj == WELL THEN wish(player_x, player_y, level)
    ELSE IF base_is_fountain(obj) THEN fountain()
    ELSE
        tile_msg(obj) // generic description from data
        IF base_is_clock(obj) THEN
            hr = hour % 12; IF hr == 0 THEN hr = 12
            print_time(hr, minute) // prints HH:MM and AM/PM
        ENDIF
    ENDIF
ENDFUNCTION
```

Notes:

- Surface Look requires a direction and prints a header (“Thou dost see”).
- Signs are handled by tracing back from the sign-holder tile to the sign face before printing text.
- Special tiles (telescope, well, fountains, clocks) trigger dedicated interaction flows; otherwise, a data-driven description (tile_msg) is printed.
- Dungeon Look depends on active light; see Environment → Light Sources & Vision: [Environment.md#light-sources--vision](Environment.md#light-sources--vision).

### Look Special Tiles Matrix (Surface)

| Tile            | Action/Effect                | UI/Text                      | Notes                                     |
|-----------------|------------------------------|------------------------------|-------------------------------------------|
| Telescope       | astronomy()                  | “Looking...” (sky view)      | Night-only stars view logic inside        |
| Well            | wish(x,y,level)              | “Drop a coin?” + wish prompt | May spawn horse in specific towns         |
| Fountain        | fountain()                   | “a gurgling fountain!”       | Simple flavor on surface                  |
| Clock           | print_time(hr, minute)       | HH:MM + AM/PM                | After tile_msg()                          |
| Sign (holders)  | trace to sign, tile_msg()    | sign text                    | Walk back from holder to face             |

Notes:

- For surface fountains, the dungeon-style drink prompt does not occur (handled only in dungeon look flow).
- Sign holders (N/E/W) are traced back to the sign face before printing the sign’s message.

### Look Dungeon Tiles Matrix

| Tile Group (tile & 0xF0) | Description           | Example UI Text             |
|--------------------------|-----------------------|-----------------------------|
| 0x00                     | Passage               | "a passage."               |
| 0x10                     | Up ladder             | "an up ladder."            |
| 0x20                     | Down ladder           | "a down ladder."           |
| 0x30                     | Ladder (unspecified)  | "a ladder."                |
| 0x40                     | Wooden chest          | "a wooden chest."          |
| 0x50                     | Fountain              | "a fountain."              |
| 0x60                     | Pit                   | "a pit."                   |
| 0x70                     | Open chest            | "an open chest."           |
| 0x80                     | Energy field (see below) | "an energy field."     |
| 0x90                     | Nothing of note       | "nothing of note."         |
| 0xA0 / 0xE0 / 0xF0       | Heavy door            | "a heavy door."            |
| 0xB0 / 0xD0              | Wall                  | "a wall."                  |
| 0xC0                     | Special wall (see below) | varies                  |

#### Energy Fields Detail

| Exact Tile | Description             |
|------------|-------------------------|
| 0x80       | A sleep field.          |
| 0x81       | A poison gas field.     |
| 0x82       | A wall of fire.         |
| 0x83       | An electric field.      |
| other 0x8x | An energy field.        |

#### Special Wall Detail (0xC0 group)

- Uses `dngtype & 0x0F` to choose a description:
  - 0x01 → "a dripping stalactite."
  - 0x02 → "a caved in passage."
  - Otherwise: 1/255 chance for "an unfortunate software pirate.", else "a less fortunate adventurer."

#### Fountain Drink Outcomes (Dungeon)

Prompt: "Will you drink?" → Y/N. On Y, prints "Yes.  Gulp!" then outcome by exact tile value:

| Exact Tile | Outcome          |
|------------|------------------|
| 0x50       | Cured! (status = Good)               |
| 0x51       | Healed! (HP = max)                  |
| 0x52       | Poisoned! (status = Poisoned)       |
| other 0x5x | Bad taste. (damage random(0,7))     |


## Open

### Open — Towns/Overworld

```pseudocode
FUNCTION command_open_surface():
    // Close any previously opened timed door state
    shut_previous_door(doortyp, doorx, doory)
    IF NOT getdir() THEN RETURN
    tx = player_x + tempx; ty = player_y + tempy
    obj = *look(tx, ty)
    SWITCH obj:
        CASE FOOTLOCKER: show_message("It's open!\n"); RETURN
        CASE PORTCULLIS: show_message("Too heavy!\n"); RETURN
        CASE LOCKED_DOOR, WINDOW_LOCK_DOOR, MAGIC_LOCKED_DOOR, WINDOW_MAG_DOOR:
            show_message("Locked!\n"); RETURN
        CASE DOOR, WINDOW_DOOR:
            doortyp = obj; doortim = 4; doorx = tx; doory = ty
            set_tile(tx, ty, FLOOR); force_los_update()
            show_message("Opened!\n"); RETURN
        DEFAULT:
            open_chest_surface(tx, ty, level) // includes non-chest handling
            RETURN
    ENDSWITCH
ENDFUNCTION
```

## Jimmy

Picks standard locks on doors and chests. Does not work on magic-locked doors. Dungeon variant targets the tile ahead; surface variant targets the facing door or chest.

```pseudocode
FUNCTION command_jimmy():
    show_message("Jimmy-")
    IF is_dungeon() THEN jimmy_dungeon() ELSE jimmy_surface()
ENDFUNCTION

FUNCTION jimmy_surface():
    IF NOT getdir() THEN RETURN
    tx = player_x + tempx; ty = player_y + tempy
    tile = top_tile_at(tx, ty)
    SWITCH tile:
        CASE LOCKED_DOOR, WINDOW_LOCK_DOOR:
            IF have_keys() THEN
                IF random_key_break() THEN show_message("Key broke!\n"); consume_one_key(); RETURN
                set_tile(tx, ty, tile - 1) // to DOOR/WINDOW_DOOR
                show_message("Unlocked!\n"); force_los_update()
            ELSE show_message("No keys!\n")
        CASE MAGIC_LOCKED_DOOR, WINDOW_MAG_DOOR:
            show_message("Magically Locked!\n")
        DEFAULT:
            // Chests are objects on surface; use Open to handle them
            show_message("Not lock!\n")
    ENDSWITCH

FUNCTION jimmy_dungeon():
    // Dungeon doors and chests use ahead-of-avatar targeting
    IF NOT dir_dng(face) THEN RETURN
    tx, ty = tempx, tempy
    t = dng_tile_at(level, ty, tx)
    IF (t & 0xF0) == 0x40 THEN // closed chest
        // Separate helper for dungeon chest jimmy odds/flow
        jimmy_chest_dungeon(tx, ty, level)
    ELSE IF (t & 0xF0) IN {0xA0, 0xE0, 0xF0} THEN // heavy doors
        IF is_magic_locked_door(t) THEN show_message("Magically Locked!\n") ELSE unlock_dungeon_door(tx, ty)
    ELSE show_message("Not lock!\n")
ENDFUNCTION
```

Notes:

- Magic-locked doors require In Ex Por (or Skull Key) to unlock; see Spells → In Ex Por and Commands → Use (Skull Key). Standard locks can be picked with keys/Jimmy.
- An Sanct can also unlock standard doors or disarm chests (dungeon); see Spells → An Sanct.
- Dungeon chest jimmy has its own trap/odds handling; see Dungeon → Chest Traps and Interactions.

```pseudocode
FUNCTION open_chest_surface(x, y, lvl):
    // Find chest object at x,y on current level (or surface if outdoors)
    obj_idx = find_object_at(x, y, lvl)
    IF obj_idx == NONE THEN show_message("Nothing to open!\n"); RETURN
    IF object[obj_idx].tile == BOX THEN show_message("Can't!\n"); RETURN

    plr = select_character(); IF plr == NONE THEN RETURN

    maxval = object[obj_idx].number // encodes treasure value and trap bit 7
    remove_object_from_lists(obj_idx); fast_los_update()
    IF is_town_map() THEN dec_karma(2)

    IF maxval > 127 THEN
        maxval &= 0x7F; show_message("Trapped!\n"); trigger_random_trap(plr)
        IF is_combat_map() AND player_dead(plr) THEN
            mark_pc_dead_in_combat_queue(plr)
            IF plr == active_player THEN active_player = NOTHING
            update_screen()
        ENDIF
    ENDIF
    found = FALSE; chkmisc(maxval, x, y, lvl, &found); chkarms(maxval, x, y, lvl, &found)
    IF NOT found THEN show_message("Chest empty!\n")
ENDFUNCTION
```

### Open — Dungeon

```pseudocode
FUNCTION command_open_dungeon():
    tile = dng_tile_at(level, player_y & 7, player_x & 7)
    IF (tile & 0xF0) == 0x40 THEN // closed chest underfoot
        plr = select_character(); IF plr == NONE THEN RETURN
        IF (tile & 7) != 0 THEN trigger_random_trap(plr)
        set_dungeon_tile(level, player_y & 7, player_x & 7, 0x70 + (tile & 8))
        show_message("\nChest opened\n")
    ELSE IF (tile & 0xF0) == 0x70 THEN show_message("Already Open!\n")
    ELSE show_message("What?\n")
ENDFUNCTION
```

Notes:

- Locked and magic-locked doors: Open reports state (“Locked!”), leaving unlocking to Jimmy/Keys (standard locks) or Skull Key/In Ex Por (magic locks).
- Spell interactions: An Sanct can unlock standard doors or disarm dungeon chests (see Spells → An Sanct). In/An Ex Por handle magical locks (see Spells → In Ex Por / An Ex Por). Skull Key behaves like In Ex Por.
- Town chest opening reduces Karma by 2; dungeons have no karma change.
- Surface chest open handles trap checks and loot distribution; dungeons use the underfoot tile variant.

#### Timed Doors (Surface)

Some opened doors on surface maps auto-close after a short duration.

```pseudocode
// Called each turn to restore a previously opened door
FUNCTION shut_previous_door(doortyp, doorx, doory):
    IF doortim > 0 THEN
        doortim = doortim - 1
        IF doortim == 0 THEN
            // Restore the original door tile at its position
            set_tile(doorx, doory, doortyp)
            force_los_update()
            doortyp = 0; doorx = 0; doory = 0
        ENDIF
    ENDIF
ENDFUNCTION
```

#### Open Chest Flow (Surface) — Diagram

1. Find chest object at target tile; if none → “Nothing to open!”; if BOX → “Can’t!”.
2. Ask “Who will?”; if none selected → return.
3. Remove chest object from object lists; LOS fast update.
4. If town map → Karma −2 (min 0).
5. If chest trapped (bit7 set):
   - Clear bit7 (retain treasure value 0..127);
   - Print “Trapped!”; trigger_random_trap(selected character);
   - If in combat and the character dies, mark dead in combat queue, update, continue.
6. Distribute loot via `chkmisc(maxval)` and `chkarms(maxval)`; track if any items found.
7. If no items found → “Chest empty!”.


## Push

```pseudocode
FUNCTION command_push():
    // Auto-shut doors if needed
    shut_previous_door(doortyp, doorx, doory)
    IF NOT getdir() THEN RETURN
    dx, dy = tempx, tempy

    // Combat: temporarily map player position to active combatant for push resolution
    IF is_combat_map() THEN save_avatar_pos(); map_to_active_combatant()

    tx = player_x + dx; ty = player_y + dy
    IF is_combat_map() AND checktrap(tx, ty) THEN RETURN // stepping trap guarding push

    tile2 = *look(tx, ty)
    IF object_present_at(tx, ty, level) OR NOT is_pushable(tile2) THEN
        show_message("Won't budge!\n"); RETURN
    ENDIF

    // Legal floor ahead depends on object: cannons require hex metal floor
    legal_floor = (base_tile(tile2) == CANNON) ? HEX_FLOOR : FLOOR

    t1x = tx + dx; t1y = ty + dy // tile one past the object
    tile3 = *look(t1x, t1y)
    tile1 = *look(player_x, player_y)

    IF NOT object_present_at(t1x, t1y, level) AND tile3 == legal_floor THEN
        // Push object forward into tile3
        push_it(tile2, tile3, tx, ty, t1x, t1y, dx, dy)
    ELSE
        // Try swapping player with object if player stands on legal floor
        IF tile1 == legal_floor THEN
            swap_it(tile1, tile2, player_x, player_y, tx, ty, dx, dy)
        ELSE
            show_message("Won't budge\n"); RETURN
        ENDIF
    ENDIF

    // Move player forward one tile
    player_x += dx; player_y += dy; force_los_update()

    IF is_combat_map() THEN
        sync_combatant_and_object_positions(dx, dy)
        restore_avatar_pos(); update_screen()
    ENDIF
ENDFUNCTION
```

```pseudocode
FUNCTION is_pushable(tile):
    SWITCH base_tile(tile):
        CASE PLANT, CHAIR, CHAIR+1, CHAIR+2, CHAIR+3,
             DESK, BARREL, VANITY, PITCHER, DRAWERS,
             END_TABLE, FOOTLOCKER, CANNON, CANNON+1, CANNON+2, CANNON+3:
            RETURN TRUE
    END
    RETURN FALSE
ENDFUNCTION
```

```pseudocode
FUNCTION push_it(tile2, tile3, tx, ty, t1x, t1y, dx, dy):
    show_message("Pushed!\n")
    set_tile(t1x, t1y, tile2)
    set_tile(tx,  ty,  tile3)
    bt = base_tile(tile2)
    IF bt == CHAIR OR bt == CANNON THEN
        set_tile(t1x, t1y, turn_shape(dx, dy, bt, reverse=0))
    ENDIF
ENDFUNCTION
```

```pseudocode
FUNCTION swap_it(tile1, tile2, px, py, tx, ty, dx, dy):
    show_message("Pulled!\n")
    set_tile(px, py, tile2)
    set_tile(tx, ty, tile1)
    bt = base_tile(tile2)
    IF bt == CHAIR OR bt == CANNON THEN
        set_tile(px, py, turn_shape(dx, dy, bt, reverse=1))
    ENDIF
ENDFUNCTION
```

```pseudocode
FUNCTION turn_shape(dx, dy, shape, reverse):
    // Returns oriented shape based on movement direction; reverse flips N<->S or E<->W
    IF dx == 0 AND dy == -1 THEN shape += NORTH
    IF dx == 1 AND dy == 0 THEN shape += EAST
    IF dx == 0 AND dy == 1 THEN shape += SOUTH
    IF dx == -1 AND dy == 0 THEN shape += WEST
    IF reverse THEN shape ^= 0x02
    RETURN shape
ENDFUNCTION
```

### Push Outcomes Matrix (At-a-Glance)

| Case                                   | Conditions                                                               | Result          | UI Text       |
|----------------------------------------|---------------------------------------------------------------------------|-----------------|---------------|
| Push forward (object moves)            | Next tile (t1x,t1y) is legal floor for object AND unoccupied             | Object → t1x,t1y; tile under object filled by prior tile | “Pushed!”     |
| Swap with player (pull)                | Next tile blocked OR not legal; player stands on legal floor for object  | Player and object swap positions                         | “Pulled!”     |
| Blocked (won’t budge)                  | Target tile not pushable OR occupied; or player not on legal floor       | No movement     | “Won’t budge!”|
| Combat trap check before push          | Combat map and checktrap(tx,ty) true                                     | Trap triggers   | (trap messages)|

Notes:

- Legal floor: `FLOOR` for standard objects; `HEX_FLOOR` for cannon family.
- Orientation: Chairs and cannons are oriented to movement direction (`turn_shape`), reversed on swaps.
- Combat: Player position is mapped to active combatant during push and restored afterward; positions are kept in sync.

## Get

### Get — Towns/Overworld

```pseudocode
FUNCTION command_get_surface():
    IF NOT getdir() THEN RETURN
    dx, dy = tempx, tempy; tx = player_x + dx; ty = player_y + dy
    show_message("\n")

    // Try to pick up an object from the object list at (tx,ty)
    obj_idx = find_object_at_surface(tx, ty, level)
    IF obj_idx != NONE THEN
        obj = object[obj_idx].tile; num = object[obj_idx].number
        // Only certain object tiles are “gettable” on surface
        IF obj < HORSE OR obj == MOONSTONE OR obj == CARPET OR base_is_shard_family(obj) THEN
            get_chest(obj, num, obj_idx) // generic handler for lootable objects
            RETURN
        ENDIF
    ENDIF

    // Otherwise, inspect the map tile
    tile = *look(tx, ty)
    SWITCH tile:
        CASE TORCH1, TORCH2:
            set_tile(tx, ty, FLOOR); force_los_update();
            IF is_town_map() THEN findlights()
            torch_light = 100; show_message("Borrowed!\n"); delay_glide(); update_screen()
        CASE CROPS:
            set_tile(tx, ty, PLOWED_FIELD); force_los_update(); show_message("Crops picked!\n")
            inc_food(1); mark_stats_changed(); dec_karma(1)
        CASE FOOD_UP:
            IF dy == -1 THEN show_message("Can't reach plate!\n") ELSE
                set_tile(tx, ty, TABLE); force_los_update(); show_message("Mmmmm...!\n")
                inc_food(1); mark_stats_changed(); dec_karma(1)
        CASE FOOD_DOWN:
            IF dy == 1 THEN show_message("Can't reach plate!\n") ELSE
                set_tile(tx, ty, TABLE); force_los_update(); show_message("Mmmmm...!\n")
                inc_food(1); dec_karma(1)
        CASE FOOD_BOTH:
            IF dx != 0 THEN show_message("Can't reach plate!\n") ELSE
                set_tile(tx, ty, (dy == 1) ? FOOD_DOWN : FOOD_UP); force_los_update(); show_message("Mmmmm...!\n")
                inc_food(1); dec_karma(1)
        DEFAULT:
            show_message("Nothing to get!\n"); RETURN
    ENDSWITCH
ENDFUNCTION
```

```pseudocode
FUNCTION get_chest(obj, num, idx):
    SWITCH obj:
        CASE CHEST: show_message("Open it first!\n"); RETURN
        CASE MOONSTONE: show_message("A moonstone!\n"); mark_stone_taken(num)
        CASE CARPET: show_message("A magic carpet!\n"); inc_carpets(); maybe_exterminate_carpet_npc()
        CASE FOOD: print_amount(num); show_message(" food!\n"); inc_food(num)
        CASE BOX: show_message("A sandalwood box!\n"); set_wooden_box_acquired(); mark_related_npc_dead()
        CASE TORCH: print_amount(num); show_pluralized(" torch"); inc_torches(num)
        CASE GEM: print_amount(num); show_pluralized(" gem"); inc_gems(num)
        CASE KEY:
            IF num > 127 THEN skull_keys += (num & 0x7F); print_amount(num&0x7F); show_message(" odd key"); show_plural_mark()
            ELSE keys += num; print_amount(num); show_message(" key"); show_plural_mark()
        CASE SCROLL:
            IF num == 0xFF THEN show_message("The plans for the HMS Cape!\n"); plans = 0xFF
            ELSE show_message("A scroll: "); print_scroll_abbrev(num); inc_scrolls(num)
        CASE GOLD: print_amount(num); show_message(" gold!\n"); inc_gold(num)
        CASE POTION: show_message("A "); print_potion_color(num); show_message(" potion!\n"); inc_potions(num)
        CASE ARROWS, QUARRELS, WEAPON_FAMILY: add_weapons(num)
        CASE SHARD_FAMILY: acquire_shard_variant(num)
        CASE CROWN: acquire_crown(); show_message("The Crown of Lord British!\n"); kill_and_exterminate_related_npc(idx)
        CASE SCEPTRE: acquire_sceptre(); show_message("The Sceptre of Lord British!\n")
        CASE AMULET: acquire_amulet(); show_message("The Amulet of Lord British!\n")
        DEFAULT: show_message("Nothing to get!\n"); RETURN
    END
    // Clear the picked-up object
    IF idx < 32 THEN remove_object_from_lists(idx)
    fast_los_update(); mark_stats_changed()
ENDFUNCTION
```

### Get — Dungeon

```pseudocode
FUNCTION command_get_dungeon():
    show_message("Get\n")
    tile = dng_tile_at(level, player_y & 7, player_x & 7)
    IF (tile & 0xF0) == 0x40 THEN show_message("Must open first!\n")
    ELSE IF (tile & 0xF0) == 0x70 THEN
        // Loot an opened chest underfoot; probabilities vary by level
        mark_tile_as_looted(level, player_y & 7, player_x & 7)
        show_message("contents\nof chest\nYou find:\n")
        for i in 0..6:
            IF random(1, (level + 1) * 4) >= idif[i] THEN
                IF i == POTION_BUCKET THEN get_chest(POTION, random(0,7), 32)
                ELSE IF i == SCROLL_BUCKET THEN get_chest(SCROLL, random(0,7), 32)
                ELSE
                    num = (i == GOLD_BUCKET) ? random(1, level * 8) : random(1, ichance[i])
                    get_chest(items[i], num, 32)
                ENDIF
        END FOR
    ELSE show_message("Not here!\n")
ENDFUNCTION
```

Notes:

- Surface Get prioritizes object-list items at the target tile; otherwise acts on the map tile for special cases (torches/crops/food).
- Dungeon Get only works on an opened chest underfoot; closed chests must be opened first, and an empty/non-chest tile rejects the action.
- Picking up certain items in towns (crops/food) reduces Karma by 1.

### Get Outcomes Matrix (Surface)

| Target (priority)               | Conditions                               | Result/Effect                       | UI Text             | Karma |
|---------------------------------|-------------------------------------------|-------------------------------------|---------------------|-------|
| Object-list item at (tx,ty)     | obj < HORSE OR obj in {MOONSTONE,CARPET,SHARD_FAMILY} | get_chest(obj,num,idx); clears object | varies (e.g., “A moonstone!”) | 0     |
| TORCH1/TORCH2 (wall-mounted)    | —                                         | Tile→FLOOR; torch_light=100; relight| “Borrowed!”         | 0     |
| CROPS                           | —                                         | Tile→PLOWED_FIELD; food+1           | “Crops picked!”     | −1    |
| FOOD_UP (north plate)           | dy == 1 (player south of table)           | Tile→TABLE; food+1                   | “Mmmmm...!”         | −1    |
| FOOD_DOWN (south plate)         | dy == −1 (player north of table)          | Tile→TABLE; food+1                   | “Mmmmm...!”         | −1    |
| FOOD_BOTH (center food)         | dx==0; dy picks UP/DOWN plate             | Tile→FOOD_UP/DOWN; food+1            | “Mmmmm...!”         | −1    |
| Other                           | —                                         | None                                 | “Nothing to get!”   | 0     |

Notes:

- Object-list pickup takes precedence over tile interactions; if nothing “gettable” exists in objects, tile fallback is attempted.
- Picking table food or crops reduces Karma by 1.
- Dungeon Get is separate (underfoot opened chest); see the dungeon section above.

## Talk (Freed NPC Nuance)

```pseudocode
FUNCTION talk_surface(direction):
    tx, ty = player_x + direction.dx, player_y + direction.dy
    tile = top_tile_at(tx, ty)
    IF tile IN {STOCKS, MANACLES} THEN
        IF no_object_present_at(tx, ty) THEN show_message("No one is there!\n"); RETURN
        npc_obj = get_object_at(tx, ty)
        npc = obj_to_npc(npc_obj)
        // If recently freed (via Jimmy), ensure the NPC acknowledges and follows
        IF npc_valid(npc) THEN
            clear_tlk(npc)
            IF npc_is_alive(npc) THEN set_follow_schedule(npc); show_message("\n\"I thank thee!\"\n"); inc_karma(2)
            killnpc(npc) // release from stocks/manacles
        ENDIF
        RETURN
    ENDIF
    // Otherwise, standard conversation flow (see Talk system docs)
    start_conversation_with_target(direction)
ENDFUNCTION
```

Notes:

- Stocks/manacles tiles use Talk to trigger an acknowledgement when the NPC has been freed; karma +2 is awarded on thanks.
- Normal Talk defers to the conversation system; this snippet only covers the freed NPC nuance.

## Search

### Search — Towns/Overworld

```pseudocode
FUNCTION command_search_surface():
    IF NOT getdir() THEN RETURN
    tx = player_x + tempx; ty = player_y + tempy

    // Moongate tile exclusion (handled elsewhere)
    IF top_tile_at(tx, ty) == MOONGATE THEN RETURN

    // Stone cache, reagent cache, and general hidden object list
    IF search_stones_at(tx, ty, level) THEN RETURN
    IF search_reagents_at(tx, ty) THEN RETURN
    search_hidden_objects_at(tx, ty, level) // daily skull keys, castle keys, glass swords (Secrets.md)
ENDFUNCTION
```

### Search — Dungeon (Ahead)

```pseudocode
FUNCTION command_search_dungeon():
    // Delegates to the detailed dungeon search flow
    search_dungeon_ahead(face)
ENDFUNCTION
```

Notes:

- Surface Search prioritizes: stone caches → reagents → location-specific HiddenObjects (see Secrets.md for daily/conditional spawns).
- Dungeon Search requires light and reports tile-ahead features (trap hints, ladders, fountains), with chance-based outcomes.

## Klimb

Klimb handles ladders and grates to change floors or traverse vertical connectors. On small maps, it supports both on-tile and directional usage.

```pseudocode
FUNCTION command_klimb():
    // Small maps (towns/overworld buildings/underworld small maps)
    IF is_small_map():
        // Try immediate on-tile Klimb first (ladder/grate under avatar)
        IF try_klimb_at_current_tile() THEN RETURN // see Objects: try_klimb_at_current_tile

        // Otherwise prompt for direction and attempt Klimb in that direction
        show_message("Klimb-")
        IF NOT getdir() THEN RETURN
        tx = player_x + tempx; ty = player_y + tempy
        t = top_tile_at(tx, ty)
        IF t IN {LADDER_UP, AVATAR_ON_LADDER_UP} THEN
            IF can_go_up_one_floor() THEN go_up_one_floor(); show_message("Klimb-Up!\n"); RETURN
            ELSE show_message("Can't go higher!\n"); RETURN
        ELSE IF t IN {LADDER_DOWN, AVATAR_ON_LADDER_DOWN, GRATE} THEN
            IF can_go_down_one_floor() THEN go_down_one_floor(); show_message("Klimb-Down!\n"); RETURN
            ELSE show_message("Can't go lower!\n"); RETURN
        ELSE
            show_message("Nothing to klimbe!\n"); RETURN
        ENDIF
    ENDIF

    // Dungeon (first-person)
    IF is_dungeon():
        // Use the tile ahead; ladders change level
        IF NOT dir_dng(face) THEN RETURN // sets tempx, tempy ahead
        tx, ty = tempx, tempy
        tile = dng_tile_at(level, ty, tx)
        group = tile & 0xF0
        IF group == 0x10 THEN // Up ladder
            IF can_go_up_one_floor() THEN go_up_one_floor(); show_message("Klimb-Up!\n") ELSE show_message("Can't go higher!\n")
        ELSE IF group == 0x20 THEN // Down ladder
            IF can_go_down_one_floor() THEN go_down_one_floor(); show_message("Klimb-Down!\n") ELSE show_message("Can't go lower!\n")
        ELSE
            show_message("Nowhere to klimbe.\n")
        ENDIF
        RETURN
    ENDIF

    // Overworld top-level map: no vertical transitions by Klimb
    show_message("No effect!\n")
ENDFUNCTION
```

### Klimb Context Matrix

| Context            | Target tile(s)                            | Result                  | Notes                                  |
|--------------------|-------------------------------------------|-------------------------|----------------------------------------|
| Small map (on-tile)| `LADDER_UP`, `LADDER_DOWN`, `GRATE`        | Up/Down one floor       | Uses on-tile first; see Objects.md      |
| Small map (dir)    | Tile ahead is ladder/grate                | Up/Down one floor       | Prompts “Klimb-” for direction          |
| Dungeon            | `(tile & 0xF0)==0x10/0x20` (Up/Down)       | Up/Down one floor       | Uses forward cell in view               |
| Overworld          | —                                         | No effect               | No vertical connectors on overworld     |

### Klimb — Overworld (Mountains with Grapple)

```pseudocode
FUNCTION command_klimb_overworld_mountain():
    IF grapples == 0 THEN show_message("With what?\n"); RETURN
    IF NOT on_foot() THEN show_message("On foot!\n"); RETURN
    show_message("Klimb-")
    IF NOT getdir() THEN RETURN
    dx, dy = tempx, tempy; tx = player_x + dx; ty = player_y + dy
    tile = top_tile_at(tx, ty)
    IF tile == PEAKS THEN show_message("Impassable!\n"); RETURN
    IF tile != MOUNTAINS THEN show_message("Not climbable!\n"); RETURN

    // Per-member dexterity check for falls (1..5 damage on fail)
    FOR i = 0 TO party_size-1:
        IF player[i].status == DEAD THEN CONTINUE
        IF player[i].dexterity < random(1, 30) THEN show_message("Fell!\n"); damageplayer(i, random(1,5))
    ENDFOR
    move_out(dx, dy) // advance into the mountain tile
ENDFUNCTION
```

Notes:

- Requires at least one grapple in inventory; does not consume a grapple.
- Only works on foot; target tile must be `MOUNTAINS` and not `PEAKS`.
- Magical locks or sealed transitions are out of scope here; doors/grates follow door logic elsewhere.
- Some grates behave as “down” transitions per map rules.
 - Dungeon alternatives: Uus Por/Des Por move up/down floors via magic; see Spells → Uus Por / Des Por.

## Pass Turn (Space)

Pressing Space consumes a turn without moving or acting. This advances time and processes world/combat updates deterministically.

```pseudocode
FUNCTION command_pass_turn():
    IF in_combat():
        // End the current actor's turn immediately
        end_current_actor_turn()
        advance_combat_initiative() // next actor
        RETURN
    ENDIF

    // Special-case: outdoors ship under sail — stop sailing
    IF is_overworld() AND ship_is_sailing() THEN show_message("Sheets in irons!\n"); stop_ship_sailing()

    // Non-combat: advance one tick without moving
    advance_clock_by_one_tick()
    apply_environment_hazards_under_party()      // swamp/lava/fireplace (Environment.md)
    decay_effect_durations()                     // spells, fields, statuses (Combat_Effects.md)
    update_npc_schedules_and_ai()                // towns and small maps (NPC_Schedules.md)
    move_world_features()                        // whirlpools, pirates, waterfalls (Movement_Overworld.md)
    try_random_encounters_if_applicable()        // overworld spawn checks (Encounters.md)
    consume_light_sources_if_needed()            // torch/magic light in dungeon (Environment.md)
    refresh_fov_and_ui()
ENDFUNCTION
```

### Pass Turn Effects Matrix

| Mode       | Primary effect                | Systems updated                          |
|------------|-------------------------------|------------------------------------------|
| Combat     | End actor turn                | Initiative queue, field ticks, statuses  |
| Small map  | Advance one tick              | NPC AI/schedules, hazards, fixtures      |
| Overworld  | Advance one tick              | Encounters, whirlpools/pirates, hazards  |
| Dungeon    | Advance one tick              | Torch decay, hazards, dungeon fixtures   |

Notes:

- Passing a turn never moves the party; it only advances time and world state.
- Deterministic RNG use: any rolls triggered here must consume from the central PRNG in documented order.

## View (Gem Map)

Consumes a gem to reveal a tactical map of the current area. On surface/small maps, shows a minimap centered near the party; in dungeons, displays the current level’s cell layout.

```pseudocode
FUNCTION command_view():
    show_message("View\n")
    IF gems <= 0 THEN show_message("No gems!\n"); RETURN
    gems -= 1; mark_stats_changed()

    IF is_dungeon():
        // Dungeon gem view: render current level layout (8x8 grid), avatar position, ladders/rooms
        render_dungeon_gem_map(level)
    ELSE
        // Surface/small-map gem view: minimap around the party with terrain and notable fixtures
        cx, cy = player_x, player_y
        render_surface_minimap(center=(cx,cy), radius=R)
    ENDIF

    update_screen_after_map_view()
ENDFUNCTION
```

Notes:

- Always consumes exactly one gem on invocation, regardless of context.
- Does not move the party; time advancement follows standard item/command rules for your engine.
- Dungeon rendering highlights up/down ladders and the party position; surface rendering emphasizes terrain and landmarks.
 - A View scroll reproduces this effect without consuming a gem; see Spells → View (Scroll): [Spells.md#view-scroll](Spells.md#view-scroll).

### View Outcomes Matrix

| Context     | Requirement | Effect                            | UI/Text      |
|-------------|-------------|-----------------------------------|--------------|
| Surface     | gems > 0    | Minimap render around party       | “View”       |
| Small map   | gems > 0    | Minimap render of local building  | “View”       |
| Dungeon     | gems > 0    | Level layout render (8x8 cells)   | “View”       |
| Any         | gems == 0   | None                              | “No gems!”   |

## Ztats (Party Member Stats)

Displays detailed stats for a selected party member; allows cycling through members.

```pseudocode
FUNCTION command_zstats():
    show_message("Ztats\n\n")
    idx = selchar_or_cycle() // select or cycle through party members
    IF idx == -1 THEN show_message("nobody!\n"); RETURN
    p = player[idx]

    // Header
    print_name(p.name); print_space(); print_class(p.class); print_space(); print_sex(p.sex); newline()
    print_status_line(p.status) // Good/Poisoned/Sleep/Dead/etc.

    // Core stats
    print_field("Level", p.level)
    print_field("Exp", p.experience)
    print_field("HP", p.hit_points, "/", p.hit_max)
    print_field("MP", p.magic_points, "/", p.magic_max)
    print_field("STR", p.str)
    print_field("DEX", p.dex)
    print_field("INT", p.int)

    // Equipment
    print_field("Weapon", weapon_name(p.weapon))
    print_field("Armor", armor_name(p.armor))
    print_field("Shield", shield_name(p.shield))
    print_field("Helm", helm_name(p.helm))
    IF p.ring != NONE THEN print_field("Ring", ring_name(p.ring))

    // Misc
    print_field("Condition", condition_flags_to_text(p.conditions))
    print_field("Moves", p.movement_points) // if applicable
    newline()
ENDFUNCTION
```

Notes:

- `selchar_or_cycle()` should support quick cycling through members without leaving the panel.
- Equipment fields print “None” for empty slots; ring is optional.
- Values reflect current state; if called in combat, show the combatant’s current HP/MP and conditions.
 - Shows equipment effect tags when applicable (e.g., Protection, Regen, Invisible). See Combat Effects → Equipment tags: [Combat_Effects.md#equipment-resistances--effects-display-in-ztats](Combat_Effects.md#equipment-resistances--effects-display-in-ztats).

## Attack

### Attack — Towns

```pseudocode
FUNCTION command_attack_town(direction):
    tx, ty = player_x + direction.dx, player_y + direction.dy
    tile = top_tile_at(tx, ty)

    // Mirror breaking special-case
    IF tile == MIRROR THEN set_tile(tx, ty, BROKEN_MIRROR); show_message("Broken!\n"); play_glass_break_sfx(); fast_los_update(); RETURN

    obj = looklist(tx, ty, level)
    nothing = TRUE
    IF obj != 0 THEN
        // NPC or valid target (exclude monsters, fields, shadow lords, horses)
        IF obj >= NPC AND (obj < MAGIC_FIELD OR obj >= MONGBAT_FAMILY) AND (obj & 0xFC) != SHADOW_LORD THEN nothing = FALSE
    ENDIF

    IF nothing THEN show_message("Nothing to attack!\n"); RETURN

    // Aggression consequences and guard activation
    IF obj < MONSTERS THEN dec_karma(5); activate_guards()
    ELSE IF is_devil(obj) THEN activate_guards()

    SWITCH tile:
        CASE BED, STOCKS, MANACLES:
            // Killing bound NPCs
            if tile != LOCKED case then
                show_message("Murdered!\n"); dec_karma(5); kapow_xy(tx, ty)
                npc = obj_to_npc_at(tx, ty)
                IF npc >= 0 THEN killnpc(npc); outwar(stat[npc].objnum); exterminate_npc(npc)
        DEFAULT:
            // Hand off to combat initiation if applicable
            initiate_town_combat_if_target(obj)
    ENDSWITCH
ENDFUNCTION
```

### Attack — Overworld

```pseudocode
FUNCTION command_attack_overworld(direction):
    // Typically defers to movement/encounter systems; direct attack handling is minimal on the overworld
    attempt_overworld_attack(direction) // engine-specific, often not supported directly
ENDFUNCTION
```

Notes:

- Town Attack affects karma and can activate town guards; breaking mirrors is supported directly.
- For monsters encountered in towns/dungeons, Combat mode handles attack resolution; this section only covers the surface/town initiation nuances.

### Talk Outcomes Matrix (At-a-Glance)

| Context            | Target                               | Conditions                               | Result/Effect                                   | UI Text                         | Karma |
|--------------------|--------------------------------------|-------------------------------------------|-------------------------------------------------|----------------------------------|-------|
| Town/Small Map     | Stocks / Manacles (occupied)         | NPC object present at tile                | Frees NPC; sets follow; acknowledges            | “I thank thee!”                 | +2    |
| Town/Small Map     | Stocks / Manacles (empty)            | No object present at tile                 | No interaction                                  | “No one is there!”              | 0     |
| Town/Small Map     | NPC at target tile                   | —                                         | Start conversation (Talk system)                | —                                | 0     |
| Town/Small Map     | Empty tile / non-interactive         | —                                         | No interaction                                  | “No effect!”                    | 0     |
| Overworld          | Any                                  | —                                         | Not supported                                   | “Talk-Funny, no response!”      | 0     |

Notes:

- The freed-NPC flow (stocks/manacles) is triggered via Talk if the NPC is still present at the restraint tile and has been unlocked.
- Normal NPC dialogue is deferred to the conversation system; this table only covers surface-level routing and messages.

## Board

```pseudocode
FUNCTION command_board():
    IF is_dungeon_map() THEN show_message("\nNot here!\n"); RETURN TRUE

    obj = looklist(player_x, player_y, level) // object underfoot
    idx = last_index_from_looklist()

    SWITCH base_tile(obj):
        CASE HORSE:
            IF is_town_map():
                npc = obj_to_npc(idx)
                IF npc != -1 AND npc_refuses_mount(npc) THEN show_message("\"Nay!\"\n"); RETURN TRUE
            ENDIF
            IF NOT on_foot() THEN RETURN TRUE
            show_message("horse\n"); player_shape = obj + 2 // riding variant
        CASE CARPET:
            IF NOT on_foot() THEN RETURN TRUE
            show_message("carpet\n"); player_shape = CARPET_RIDER
        CASE SKIFF:
            IF NOT on_foot() THEN RETURN TRUE
            show_message("skiff\n"); player_shape = obj // inherit facing
        CASE ANCHORED_SHIP:
            IF NOT can_board() THEN RETURN TRUE
            show_message("Ship\n")
            ship_hp = object[idx].number; object[0].number = ship_hp
            IF ship_hp < 10 THEN show_message("\nDANGER: SHIP BADLY DAMAGED!\n")
            skiffs_onboard = object[idx].misc
            IF player_shape_is_carpet_rider() THEN carpets++
            IF player_shape_is_skiff() THEN skiffs_onboard++
            IF skiffs_onboard == 0 THEN show_message("\nWARNING: NO SKIFFS ON BOARD!\n")
            player_shape = obj; object[0].misc = skiffs_onboard; mark_stats_changed()
        DEFAULT:
            show_message("What?\n"); RETURN FALSE
    ENDSWITCH

    // Remove the boarded object from the lists
    remove_object_from_lists(idx); fast_los_update()
    RETURN TRUE
ENDFUNCTION
```

Notes:

- on_foot(): prints “On foot” and returns 0 if not on foot; prevents boarding mounts/vehicles when already mounted.
- can_board(): allows boarding for PLAYER, INVIS_PC, CARPET_RIDER, SKIFF family; prints “On foot” and returns 0 otherwise.
- Boarding a ship sets player_shape to the ship tile, copies hull points to object[0].number for UI, and updates onboard skiffs; warnings are printed for low hull and zero skiffs.

## Exit (Leave Building/Town)

Stepping beyond the bounds of a small map (buildings/towns) prompts to leave. This is movement-triggered rather than a dedicated command key.

```pseudocode
FUNCTION check_smallmap_exit(new_position):
    IF is_small_map() AND position_out_of_bounds(new_position) THEN
        show_modal_prompt("Dost thou wish to leave?")
        IF answer_yes() THEN
            perform_exit_transition() // return to parent map (overworld or town)
        ELSE
            cancel_move_and_stay()
        ENDIF
        RETURN TRUE
    ENDIF
    RETURN FALSE
ENDFUNCTION
```

Notes:

- Exit transitions restore the parent map context and player position as per location linkage.
- Some exits may be blocked by quests or special conditions; those display appropriate text.

## Use

Based on the classic Use handler for notable items (IDs in comments for context).

```pseudocode
FUNCTION command_use(item_id):
    SWITCH item_id:
        CASE MAGIC_CARPET (16):
            show_message("Carpet\n\n")
            IF is_overworld() AND player_on_mountains() THEN show_message("Not here!\n"); RETURN
            IF player_is_on_foot() THEN
                show_message("Boarded!\n")
                set_player_form_to_carpet_rider_random_facing()
                decrement_carpets()
            ELSE IF player_form_is_ship() THEN
                show_message("X-it ship first!\n")
            ELSE
                show_message("Only on foot!\n")
            ENDIF
        CASE SKULL_KEY (17):
            show_message("Skull Key\n")
            IF is_overworld() OR is_dungeon() THEN
                IF NOT spelldir() THEN RETURN // sets tempx,tempy
                success = wizard_unlock_magic_at(tempx, tempy) // In Ex Por semantics
                IF success AND is_town_map() THEN kapow_xy(tempx, tempy)
            ELSE
                show_message("Not here!\n")
            ENDIF
            decrement_skull_keys()
        CASE AMULET (18):
            show_message("Amulet\n\n")
            IF remove_from_inventory(AMULET) THEN
                show_message("Wearing the Amulet of LB")
                dur_spell(LB_AMULET, 255, sfx_level=9) // permanent while equipped
            ENDIF
        CASE CROWN (19):
            show_message("Crown\n\n")
            IF remove_from_inventory(CROWN) THEN
                show_message("Thou dost don the Crown of LB")
                dur_spell(CROWN, 255, sfx_level=9)
            ENDIF
        CASE SCEPTRE (20):
            show_message("Sceptre\n\nWielding the Sceptre of LB")
            play_sceptre_pulse()
            if is_surface_or_dungeon():
                // Try to clear adjacent opened chests (0x70 group) to GRASS on surface; otherwise dispel energy field
                cleared = clear_open_chests_around_or_dispel_field()
                IF NOT cleared THEN
                    temp = an_grav(sound=0) // dispel energy field helper
                    IF temp == TRUE THEN show_message("Field dissolved!\n")
                    ELSE IF temp == 0 THEN show_message("No effect!\n")
                ENDIF
        CASE SPYGLASS:
            show_message("Spyglass\n\n")
            IF is_overworld() AND is_daytime() THEN show_message("No stars!\n")
            ELSE show_message("Looking...\n"); astronomy()
        DEFAULT:
            show_message("What?\n")
    ENDSWITCH
ENDFUNCTION
```

Helper:

```pseudocode
FUNCTION wizard_unlock_magic_at(x, y):
    tile = top_tile_at(x, y)
    SWITCH tile:
        CASE MAGIC_LOCKED_DOOR: set_tile(x, y, DOOR); force_los_update(); RETURN TRUE
        CASE WINDOW_MAG_DOOR:   set_tile(x, y, WINDOW_DOOR); force_los_update(); RETURN TRUE
        DEFAULT: RETURN FALSE
    END
ENDFUNCTION
```

Notes:

- Carpet: only boards on foot and not on ship; disallowed on mountains; consumes a carpet charge.
- Skull Key: behaves like a directional In Ex Por to unlock magic-locked doors; not usable in towns (“Not here!”); always consumes one.
- Amulet/Crown: enable persistent effects (e.g., light/negation) with duration set to 255. For light behavior, see [Environment.md#light-sources--vision](Environment.md#light-sources--vision).
- Sceptre: dissolves nearby open chest tiles on surface or uses An Grav to dispel an energy field; prints “Field dissolved!” or “No effect!”.
- Spyglass: astronomy view only at night on the overworld.

### Use Outcomes Matrix (At-a-Glance)

| Item            | Context              | Conditions/Checks                                 | Effect/Result                                        | UI Text                              |
|-----------------|----------------------|---------------------------------------------------|------------------------------------------------------|--------------------------------------|
| Magic Carpet    | Overworld/Surface    | On foot; not on mountains; not on ship            | Board carpet; set rider form; decrement carpets      | “Boarded!”, else “Only on foot!”/“X-it ship first!”/“Not here!” |
| Skull Key       | Overworld/Dungeon    | Direction given; target is magic-locked door      | Unlock to DOOR/WINDOW_DOOR; kapow in towns           | “Skull Key”; town success shows kapow; “Not here!” in towns    |
| Amulet          | Any                  | Removed from inventory                            | Equip; set persistent effect duration (255)          | “Wearing the Amulet of LB”           |
| Crown           | Any                  | Removed from inventory                            | Equip; set persistent effect duration (255)          | “Thou dost don the Crown of LB”      |
| Sceptre         | Overworld/Dungeon    | —                                                 | Try clear nearby open chests or dispel field         | “Wielding the Sceptre of LB”; “Field dissolved!”/“No effect!”  |
| Spyglass        | Overworld            | Nighttime                                         | Astronomy sky view                                   | “Looking...”; day: “No stars!”       |

Notes:

- Skull Key always consumes one use regardless of success; Amulet/Crown/Sceptre “Use” remove from inventory and apply effects.
- Carpet “Use” consumes a carpet instance; boarding conditions enforced.

## Fire (Cannons)

### Fire — Town Cannons

```pseudocode
FUNCTION command_fire_town():
    IF is_dungeon_map() THEN show_message("What?\n"); RETURN
    IF is_overworld() THEN fire_ship(); RETURN

    // Determine a cannon adjacent to the player and its firing axis
    shut_previous_door(doortyp, doorx, doory)
    dir = infer_cannon_direction_from_adjacent_tiles()
    IF dir == NONE THEN show_message("What?\n"); RETURN

    tx = player_x + dir.dx; ty = player_y + dir.dy
    cannon_sx = 5 + dir.dx; cannon_sy = 5 + dir.dy

    // Cannon shape encodes its facing; compute dx,dy accordingly
    dx, dy = facing_from_cannon_tile_at(4,5,...)
    show_message("BOOOM!\n"); play_cannon_fx(); activate_guards()

    // Trace outward up to range 5; stop on first door or NPC
    isdoor = FALSE; isnpc = FALSE; range = 5
    WHILE NOT isdoor AND NOT isnpc AND (--range > 0):
        tx += dx; ty += dy; cannon_sx += dx; cannon_sy += dy
        obj = looklst2(tx, ty, level)
        IF obj == 0 THEN
            IF tile_is_any_door_or_portcullis(top_tile_at(tx,ty)) THEN isdoor = TRUE
        ELSE IF is_npc_or_valid_target(obj) THEN
            isnpc = TRUE; idx = last_index_from_looklist()
        ENDIF
    ENDWHILE

    missile(cannon_origin_x, cannon_origin_y, cannon_sx, cannon_sy, CANNONBALL)
    IF isdoor OR isnpc THEN kapow_xy(tx, ty)

    IF isdoor THEN
        show_message("Door destroyed!\n"); set_tile(tx, ty, FLOOR); losflag = 1; clear_door_timer()
    ENDIF

    IF isnpc THEN
        remove_object_from_lists(idx); losflag |= 2; dec_karma(5)
        IF npc_index_from_object(idx) >= 0 THEN killnpc(npc); exterminate_npc(npc)
        IF idx == 0 THEN damageparty() // mishap
    ENDIF
ENDFUNCTION
```

### Fire — Ship (Broadsides)

```pseudocode
FUNCTION fire_ship():
    IF NOT player_form_in_ship_family() THEN show_message("What?\n"); RETURN
    IF NOT getdir() THEN RETURN
    dx, dy = tempx, tempy
    // Require broadsides (fire only perpendicular to ship’s axis)
    IF (dx == 0 AND ship_facing_is_vertical() == FALSE) OR (dx != 0 AND ship_facing_is_vertical() == TRUE) THEN
        show_message("Fire broadsides only!\n"); RETURN
    ENDIF

    tx, ty = player_x, player_y
    glide_cannon_anim()
    FOR k IN 1..3:
        tx += dx; ty += dy
        obj = looklist(tx, ty, level)
        IF is_monster(obj) AND NOT is_whirlpool(obj) THEN
            index = last_index_from_looklist()
            target_sx = 5 + (object[index].xpos - player_x)
            target_sy = 5 + (object[index].ypos - player_y)
            not_blocked = missile(5, 5, target_sx, target_sy, CANNONBALL)
            IF not_blocked THEN
                update_screen(); kapow_xy(tx, ty)
                object[index].number -= random(1, 20) // ship HP down by 1..20
                IF object[index].number <= 127 THEN remove_object_from_lists(index); losflag |= 2
            ENDIF
            RETURN
        ENDIF
    END FOR
    // If no target, still fire for show down range 3
    missile(5, 5, 5 + 3*dx, 5 + 3*dy, CANNONBALL)
ENDFUNCTION
```

### Fire Outcomes Matrix

| Context     | Target hit               | Result/Effect                                   | UI Text            | Karma |
|-------------|--------------------------|-------------------------------------------------|--------------------|-------|
| Town        | Door                     | Door destroyed → FLOOR                          | “Door destroyed!”  | 0     |
| Town        | NPC                      | NPC removed; guards activated                   | —                  | −5    |
| Town        | None                     | No effect beyond visuals                        | —                  | 0     |
| Overworld   | Enemy Ship/Monster       | Cannonball hit; target HP −random(1,20)         | —                  | 0     |

Notes:

- Town cannon fire always activates guards.
- Ship broadsides must align with ship facing; otherwise “Fire broadsides only!”.
- Projectile visuals use missile() with erase/kapow handling on impact.

## Hole Up & Camp

Camp to rest and heal, or repair ship hull when aboard and anchored/furled.

```pseudocode
FUNCTION command_hole_up_and_camp():
    ptr = &object[0] // player object slot
    plr_shape = ptr.shape
    show_message("Hole up & ")

    // Aboard ship: repair flow
    IF player_form_in_ship_family():
        show_message("\nrepair...\n\n")
        IF ship_sails_unfurled() THEN
            show_message("Sails must be\nlowered!\n\n"); RETURN
        ENDIF
        // Allow monsters to move a bit while repairing; enforce still aboard
        FOR i IN 1..5:
            check_update(); monster_main()
            IF NOT player_form_is_anchored_ship() THEN RETURN
            addtime(5)
        END FOR
        // Increment hull until safe (>=10)
        DO:
            ptr.number += random(1,3)
            IF ptr.number > 99 THEN ptr.number = 99
        WHILE ptr.number < 10
        show_message("Hull now "); print_int(ptr.number,2); show_message("!\n\n")
        mark_stats_changed(); RETURN
    ENDIF

    // On land/foot: camping
    show_message("camp!\n\n")
    IF is_overworld() THEN
        plrshp_tile = *look(player_x, player_y)
        IF plrshp_tile > 0 AND plrshp_tile < 4 THEN
            show_message("On land or ship!\n\n"); RETURN // illegal spot
        ENDIF
        IF ptr.shape != PLAYER THEN show_message("On foot!\n"); RETURN
    ENDIF

    show_message("For how many hours? (1-9) ")
    answer = read_key_digit_or_space()
    IF answer == ' ' OR answer == '0' THEN RETURN
    time = answer - '0'

    // Optional watch if more than one eligible party member
    num_alive = count_members_with_status({'G','P'})
    guard = -1
    IF num_alive > 1 THEN
        show_message("\nWilt thou set a watch? ")
        yn = read_yes_no()
        IF yn == 'Y' THEN
            show_message("Yes\n\nWho will stand guard? ")
            guard = selchar()
            newline()
            IF guard == -1 OR player[guard].status != 'G' THEN guard = -1; show_message("None posted!\n\n")
        ELSE show_message("No\n\n")
        ENDIF
    ENDIF

    IF is_combat_or_dungeon_map():
        scenario = IN_CAMP_SCEN | DNG_COMBAT_SCEN
        init_tiles(); dngcamp(); vcombat(scenario, guard, time)
        restore_dungeon_ui_and_objects()
    ELSE
        outcamp(guard, time)
    ENDIF
    force_los_update()
ENDFUNCTION
```

### Hole Up Outcomes Matrix

| Context           | Conditions                                 | Effect                                       | UI Text                          |
|-------------------|--------------------------------------------|----------------------------------------------|----------------------------------|
| Ship repair       | Sails furled; anchored ship                | 5 ticks of world time; hull to at least 10   | “repair...”; “Sails must be lowered!” if not |
| Overworld camp    | On foot; legal tile (not sea/ship)         | outcamp(guard,time)                          | Prompts hours; optional watch    |
| Dungeon/Combat    | —                                          | dngcamp()+vcombat(scenario,guard,time)       | Prompts hours; optional watch    |

Notes:

- Ship hull increases by random(1,3) per loop until at least 10; capped at 99.
- Repair loop interleaves monster movement/time; repair aborts if player leaves ship.
- Camping selects guard only if at least two party members are in good/poisoned status.

## Enter

```pseudocode
FUNCTION command_enter():
    IF is_overworld() THEN
        loc = lookup_world_location_at(player_x, player_y)
        IF loc != NONE THEN
            enter_location(loc) // swap to small map; print entering text for location
        ELSE
            show_message("Enter what?\n")
        ENDIF
    ELSE
        show_message("Enter what?\n")
    ENDIF
ENDFUNCTION
```

## Ignite Torch

```pseudocode
FUNCTION command_ignite_torch():
    show_message("Ignite Torch!\n")
    IF NOT ignite_torch() THEN show_message("None owned!\n")
ENDFUNCTION
```

See Environment → Light Sources & Vision for timers and visibility rules: [Environment.md#light-sources--vision](Environment.md#light-sources--vision), and Torch Duration: [Environment.md#torch-duration](Environment.md#torch-duration).

## New Order (Swap Party Positions)

```pseudocode
FUNCTION command_new_order():
    show_message("\n\nSwap ")
    ch1 = selchar()
    IF ch1 == -1 THEN show_message("nobody!\n"); RETURN
    show_message(get_player_name(ch1))
    IF ch1 == 0 THEN show_message("\n\n"); show_message(get_player_name(0)); show_message(" must lead!\n"); RETURN

    show_message("\nwith ")
    ch2 = selchar()
    IF ch2 == -1 THEN show_message("nobody!\n"); RETURN
    show_message(get_player_name(ch2))
    IF ch2 == 0 THEN show_message("\n\n"); show_message(get_player_name(0)); show_message(" must lead!\n"); RETURN
    show_message("!\n")

    swap_party_positions(ch1, ch2)
ENDFUNCTION
```

## Mix Reagents

## Ready

Select a party member and ready or unready equipment into appropriate slots, enforcing weight, slot, ammo, and context rules.

```pseudocode
FUNCTION command_ready():
    plr = select_character(); IF plr == NONE THEN RETURN
    IF count_owned_arms() == 0 THEN show_message("Thou art empty-\nhanded!\n"); RETURN
    show_message("Item: ")
    selection = select_from_scroll_list(arms, highlight_readied_for(plr)) // supports paging, shows icons for readied
    IF selection == ESC THEN show_message("Done\n"); RETURN
    try_to_ready(plr, selection)

FUNCTION try_to_ready(plr, item) -> ring_vanished:
    IF item IN {ARROWS, QUARRELS} THEN RETURN FALSE // cannot directly ready ammo
    IF item IN ARMOR_FAMILY AND in_combat() THEN show_message("Thou canst not change armour in heated battle!"); RETURN FALSE
    IF is_item_already_readied(plr, item) THEN
        unready_item(plr, item); increment_inventory(item)
        IF item == INVISO_RING AND in_combat() THEN restore_player_sprite(plr)
        RETURN FALSE
    ENDIF
    IF item REQUIRES_AMMO AND no_ammo_for(item) THEN show_message("Thou hast no ammunition for that weapon!"); RETURN FALSE

    sum = strength_req_of_readied(plr) + strength_req[item]
    IF sum > player[plr].strength THEN show_message("Thou art not strong enough!"); RETURN FALSE

    SWITCH weapon_usage[item]:
        CASE HEAD_USE:  IF helm_slot_occupied(plr) THEN show_message("Remove first thy present helm!"); RETURN FALSE ELSE slot = HELM
        CASE BODY_USE:  IF armor_slot_occupied(plr) THEN show_message("Thou must first remove thine other armour!"); RETURN FALSE ELSE slot = ARMOR
        CASE ONE_HAND_USE:
            hand = freehand(plr) // returns 0 or 1 if that hand free; 2 if both free; NOTHING if none
            IF hand == NOTHING THEN show_message("Thou must free one of thy hands first!"); RETURN FALSE
            IF hand == 2 THEN hand = 0
            slot = WEAPON[hand]
        CASE TWO_HAND_USE:
            IF freehand(plr) != 2 THEN show_message("Both hands must be free before thou canst wield that!"); RETURN FALSE
            slot = WEAPON[0] // occupies both hands per engine
        CASE NECK_USE:  IF amulet_slot_occupied(plr) THEN show_message("Thou must remove thine other amulet!"); RETURN FALSE ELSE slot = AMULET
        CASE RING_USE:  IF ring_slot_occupied(plr) THEN show_message("Only one magic ring may be worn at a time!"); RETURN FALSE ELSE slot = RING
    END SWITCH

    set_slot(plr, slot, item); decrement_inventory(item)

    // Special ring behavior
    IF item IN {INVISO_RING, REGEN_RING} AND random(0,15) == 0 THEN
        show_message("\n\nRing vanishes!\n"); clear_ring_slot(plr); delay_fx(); RETURN TRUE
    IF item == INVISO_RING AND in_combat() THEN set_player_sprite_invisible(plr)
    RETURN FALSE
```

Notes:

- Ammo gating: Bows require ARROWS; Crossbows require QUARRELS; message: “Thou hast no ammunition for that weapon!”.
- Weight/strength: The sum of strength requirements of all readied items plus the candidate must be ≤ the character’s Strength.
- Negative contexts: Armor changes are disallowed in combat; enforce single ring/amulet; enforce free hands for one/two‑handed weapons.
- Ring vanish: Invisibility and Regeneration rings have a 1‑in‑16 chance to vanish when equipped; Invisibility ring sets invisible sprite in combat.

## Cast

Prompt for a spell by name, validate context and resources, then dispatch to the appropriate spell handler. Uses spell level = floor(index/6)+1.

```pseudocode
FUNCTION command_cast():
    caster = select_character(); IF caster == NONE THEN RETURN
    show_message("Spell name:\n:")
    spellnum = getspell() // maps typed name to index; returns -1 none, -2 “no effect”
    IF spellnum == -1 THEN show_message("None!\n"); RETURN
    IF spellnum == -2 THEN show_message("No effect!\n"); RETURN

    spell_level = (spellnum / 6) + 1

    // Context gating
    IF is_overworld() AND NOT spell_allowed_outdoors(spellnum) THEN not_here()
    ELSE IF in_combat() AND NOT spell_allowed_combat(spellnum) THEN not_here()
    ELSE IF in_blackthorns_castle_without_crown() OR in_stonegate() THEN show_message("Absorbed!\n"); absorb_fx(); RETURN
    ELSE IF is_town() AND NOT spell_allowed_town(spellnum) THEN not_here()
    ELSE IF is_dungeon() AND NOT spell_allowed_dungeon(spellnum) THEN not_here()

    // Reagents/stock and MP
    IF spells_stock[spellnum] == 0 THEN show_message("None mixed!\n"); RETURN
    spells_stock[spellnum] -= 1
    IF player[caster].mp < spell_level THEN show_message("M.P. too low!\n"); print_failed(); RETURN
    player[caster].mp -= spell_level
    IF player[caster].level < spell_level THEN print_failed(); RETURN

    success = TRUE; end_player_turn = TRUE
    SWITCH spellnum:
        CASE 0:  light_short();
        CASE 1:  weapon_spell(FLAM_POR);
        CASE 2:  success = cure_sleep_single();
        CASE 3:  success = cure_poison_single();
        CASE 4:  success = heal_small();
        CASE 5:  success = locate();
        CASE 6:  success = disarm_unlock_standard();
        CASE 7:  summon_daemon_charmed();
        CASE 8:  wind_change(prompt_dir(), immediate=FALSE);
        CASE 9:  xray_surface();
        ...     // Continue per Spells.md mapping (see detailed behaviors there)
        CASE 39: IF is_surface() THEN view_area() ELSE dng_view();
        CASE 46: success = gate_travel(); IF success THEN end_player_turn = FALSE
        CASE 47: success = time_stop();
    END SWITCH

    IF success == TRUE THEN show_message("Success!\n") ELSE print_failed()
    IF end_player_turn THEN finish_turn()

FUNCTION not_here():
    show_message("Not here!\n"); small_delay(); RETURN
```

Notes:

- Context negatives: Prints “Not here!” when disallowed by map type; prints “Absorbed!” in prohibited areas (e.g., Stonegate, or Blackthorn’s castle without LB’s Crown).
- Resources: Requires mixed spell stock and sufficient MP; prints “None mixed!” or “M.P. too low!” when lacking.
- Success/Failure: Some handlers return success/failure; unsuccessful casts print “Failed!”. Certain spells do not end the player’s turn (e.g., Gate Travel).

### Context Gating Table

| Context      | Allow Rule                         | Failure Text |
|--------------|------------------------------------|--------------|
| Overworld    | `OUTD_SPELL` flag required         | “Not here!”  |
| Town         | `TOWN_SPELL` flag required         | “Not here!”  |
| Dungeon      | `DUNG_SPELL` flag required         | “Not here!”  |
| Combat       | `COMB_SPELL` flag required         | “Not here!”  |
| Blackthorn’s | Block if no LB Crown               | “Absorbed!”  |
| Stonegate    | Always block (no magic allowed)    | “Absorbed!”  |

Reference (OLD/CAST1.C): `onmap==0x12 && !crown` or `onmap==0x1d` → “Absorbed!”. 0x12 = Palace_of_Blackthorn; 0x1d = Stonegate.

### Dispatch Map (from legacy)

The legacy engine dispatches by `spellnum` (0..47). The table below shows representative cases; full behavior lives in Spells.md.

- 0: Light (short) → `light_spell(100)`
- 1: Vas Flam (bolt) → `weapon_spell(FLAM_POR)`
- 2: An Zu (Awaken) → `an_zu()` success/fail
- 3: An Nox (Cure Poison) → `an_nox()`
- 4: Mani (Heal small) → `mani()`
- 6: An Sanct (Disarm/Unlock std) → `an_sanct()`
- 8: Rel Hur (Change Wind) → `wind_change(spelldir(), FALSE)`
- 10: Kal Xen (Summon Animal) → `kal_xen()`
- 14..16,20: Create Field (Poison/Sleep/Fire/Energy) → `field_spell(type)`
- 17: In Por (Blink/Teleport) → `in_por()`
- 18: An Grav (Dispel Field) → `an_grav(1)`
- 19: In Sanct (Protection) → `dur_spell('P', 20, 4)`
- 21/22: Uus Por/Des Por (Dungeon level up/down) with Doom block `onmap==0x28`
- 23: Wis Quas (Reveal Invisible) → `wis_quas()`
- 24: In Bet Xen (Summon Insects) → `in_bet_xen()`
- 25/26: An Ex Por / In Ex Por (Lock/Unlock Magic)
- 27: Vas Mani (Great Heal)
- 28/40/44/45: Storms (Purple/Green/Blue/Red) → `nukem(turn, kind, color)`
- 29: Rel Tym (Quickness) → `dur_spell('Q', 30, 5)`
- 30: In Vas Por Ylem (Earthquake) → `in_vas_por_ylem(caster)`
- 31: Quas An Wis (Mass Charm) → `dur_spell('C', 20, 6)`
- 32: In An (Negate Magic) → `dur_spell('N', 10, 6)`
- 33: Wis An Ylem (X‑Ray) → `x_ray_vision()`
- 34: An Xen Ex (Charm toggle) → `an_xen_ex()`
- 35: Rel Xen Bet (Polymorph) → `rel_xen_bet()`
- 36: Sanct Lor (Invisibility) → `sanct_lor()`
- 37: Xen Corp (Death bolt) → `weapon_spell(IN_CORP)`
- 38: In Quas Xen (Clone) → `in_quas_xen()`
- 39: In Quas Wis (View) → surface `view_area(...)`, dungeon `dng_view()`
- 42: In Mani Corp (Resurrection) → `resurrect(onwho(), FALSE)`
- 43: Kal Xen Corp (Summon Daemon) → `summon_daemon(FALSE)`
- 46: Vas Rel Por (Gate Travel) → `vas_rel_por()`; does not end turn on success
- 47: An Tym (Negate Time) → `an_tym()`

Notes:

- Doom/Stonegate blocks: Some spells have additional map‑specific blocks (e.g., Uus/Des Por in Doom via `onmap==0x28`); see each spell’s section in Spells.md for exact rules and texts.


Mix combines owned reagents into a spell component stock, typically outside combat. It validates inventory and prints simple prompts.

```pseudocode
FUNCTION command_mix_reagents():
    show_message("Mix Reagents\n\n")
    recipe = choose_spell_recipe() // user selects target spell; shows needed reagents
    IF recipe == NONE THEN RETURN
    IF NOT has_reagents_for(recipe) THEN show_message("Not enough reagents!\n"); RETURN
    consume_reagents(recipe)
    increment_spell_stock(recipe.spell_id)
    mark_stats_changed(); show_message("Mixed.\n")
ENDFUNCTION
```

Notes:

- UI may present per-spell recipes with quantities and total batch count. Selection can be per unit or in multiples.
- Deterministic: no RNG rolls; purely inventory transformation. Errors if insufficient reagents.
- Some engines support mixing only out of combat; follow local rules.
