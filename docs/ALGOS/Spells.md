# Spells

## Duration Flags & Effects (Reference)

These duration-based effects are applied via `dur_spell(flag, duration, sfx_level)` and consumed by subsystems:

- 'P' (Protection): Reduces incoming damage or improves saves while active.
- 'Q' (Quickness / Rel Tym): Skips alternate monster turns (see Combat Core → Turn Modifiers).
- 'C' (Mass Charm / Quas An Wis): Monsters may attack monsters on failed INT saves during target selection.
- 'N' (Negate Magic): Suppresses certain magical effects (teleports/special moves) while active.
- TIME_STOP (An Tym/Negate Time): Monsters do not act while active; blocked on specific maps.

Notes:

- Duration is in ticks; most uses here mirror legacy values (e.g., 20 or 30).
- Subsystems check these flags explicitly (AI target selection, special moves, turn modifiers).

## Spells Summary (At‑a‑Glance)

Legend: Context = OW (Overworld), T (Town), D (Dungeon), C (Combat). Targeting = Self, Tile, Creature, Area, View.

### Light / Sight

| Spell          | Context      | Targeting | Duration/Flag | Notes                             |
|----------------|--------------|-----------|---------------|-----------------------------------|
| In Lor         | OW/T/D/C     | Self      | magic_light   | Short light                       |
| Vas Lor        | OW/T/D       | Self      | magic_light   | Long light                        |
| Wis An Ylem    | OW/T         | Self      | —             | X‑Ray vision (surface only)       |
| In Quas Wis    | OW/T/D       | View      | —             | View map; see Commands → View     |

### Movement / Travel

| Spell     | Context | Targeting | Duration/Flag | Notes                                                   |
|-----------|---------|-----------|---------------|---------------------------------------------------------|
| In Por    | OW/C    | Blink     | —             | OW: along dir to grass; C: random legal tile (rooms off) |
| Uus Por   | D       | Level     | —             | Go up one dungeon level (blocked in special dungeons)    |
| Des Por   | D       | Level     | —             | Go down one dungeon level (blocked in special dungeons)  |
| Vas Rel Por | OW/T  | Travel    | —             | Long‑range travel per rules (markers/moongates)          |

### Protection / Negation / Time

| Spell     | Context  | Targeting | Duration/Flag | Notes                                           |
|-----------|----------|-----------|---------------|-------------------------------------------------|
| In Sanct  | OW/T/D/C | Self      | 'P'           | Protection; mitigate per engine                  |
| In An     | OW/T/D/C | Self      | 'N'           | Negate magic; suppress teleports/specials        |
| Rel Tym   | C        | Party     | 'Q'           | Quickness; skip alternate monster turns          |
| An Tym    | OW/T/D/C | World     | TIME_STOP     | Monsters don’t act; blocked on special maps      |

### Projectiles

| Spell    | Context | Targeting | Duration/Flag | Notes                                         |
|----------|---------|-----------|---------------|-----------------------------------------------|
| Grav Por | C       | Creature  | —             | Magic missile (weapon_spell); no INT save      |
| Vas Flam | C       | Creature  | —             | Firebolt (weapon_spell); no INT save  |
| Xen Corp | C       | Creature  | —             | Death bolt (weapon_spell); no INT save|

### Fields (Create)

| Spell         | Context | Targeting | Duration/Flag | Notes                               |
|---------------|---------|-----------|---------------|-------------------------------------|
| In Flam Grav  | D/C     | Tile      | field tile    | Fire field                          |
| In Nox Grav   | D/C     | Tile      | field tile    | Poison field                        |
| In Zu Grav    | D/C     | Tile      | field tile    | Sleep field                         |
| In Sanct Grav | D/C     | Tile      | field tile    | Energy field                        |

### Summoning

| Spell       | Context | Targeting | Duration/Flag | Notes                                        |
|-------------|---------|-----------|---------------|----------------------------------------------|
| Kal Xen     | C       | Tile      | —             | Summon small animal (charmed)                |
| Kal Xen Corp| C       | Tile      | —             | Summon daemon (charmed; caster INT mishap)   |
| In Bet Xen  | C       | Tile      | —             | Summon insect swarm (charmed)                |

### Control / Charm

| Spell        | Context | Targeting  | Duration/Flag | Notes                                                           |
|--------------|---------|------------|---------------|-----------------------------------------------------------------|
| An Xen Ex    | C       | Creature   | —             | Charm toggle on failed INT                                      |
| Quas An Wis  | C       | Aura       | 'C'           | Mass charm/confusion; INT save flips monster targeting loyalty  |
| An Xen Corp  | C       | Area/Undead| —             | Repel undead: set HP=1 + flee on failed INT                     |
| In Quas Corp | C       | Area       | —             | Fear: enemies flee on failed resistance                         |
| Rel Xen Bet  | C       | Creature   | —             | Polymorph creature into rat on failed INT                       |
| In Quas Xen  | C       | Creature   | —             | Clone target into adjacent tile                                 |
| Wis Quas     | C       | Global     | —             | Reveal invisible monsters (clear INVISO, restore shapes)        |

### Unlock / Lock / Dispels

| Spell     | Context | Targeting | Duration/Flag | Notes                                                     |
|-----------|---------|-----------|---------------|-----------------------------------------------------------|
| In Ex Por | OW/T/D  | Door      | —             | Wizard‑unlock magic‑locked doors                           |
| An Ex Por | OW/T/D  | Door      | —             | Wizard‑lock standard/windowed doors                        |
| An Sanct  | OW/T/D  | Chest/Door| —             | Disarm dungeon chest / unlock standard locked door         |
| An Grav   | OW/T/D  | Tile      | —             | Dispel energy field (sleep/poison/fire/electric)          |

### Healing / Food / Earthquake / Sense

| Spell         | Context  | Targeting | Duration/Flag | Notes                               |
|---------------|----------|-----------|---------------|-------------------------------------|
| Mani          | OW/T/D/C | Creature  | —             | Heal 1..30 up to max                 |
| Vas Mani      | OW/T/D/C | Creature  | —             | Big heal (engine roll)               |
| In Xen Mani   | OW/T/D/C | Party     | —             | Create 1..3 food                     |
| In Vas Por Ylem | C      | Area      | —             | Earthquake damage via DEX check      |
| In Wis        | OW       | Self      | —             | Locate position (coords)             |

### Area Damage (Storms)

| Spell             | Context | Targeting | Duration/Flag | Notes                        |
|-------------------|---------|-----------|---------------|------------------------------|
| In Nox Hur        | C       | Area      | —             | Poison storm; no INT save    |
| In Flam Hur       | C       | Area      | —             | Firestorm; no INT save       |
| In Vas Grav Corp  | C       | Area      | —             | Energy storm; no INT save    |

### Unimplemented / Reserved

| Spell | Context | Targeting | Duration/Flag | Notes                                 |
|-------|---------|-----------|---------------|---------------------------------------|
| Frotz | —       | —         | —             | Unimplemented; reserved/no‑op         |

Notes:

- Detailed behaviors and corner cases are described in individual sections above; this table is meant as a quick routing reference (where, how, and whether there’s a duration flag).
- For scroll equivalents, see Scrolls Summary below; contexts sometimes differ (e.g., View scroll disallowed in combat).

## Allowed Contexts

This table consolidates where each spell may be cast. Context rules come from legacy flags (OUTD/TOWN/DUNG/COMB) and our pseudocode; exceptions add special‑map constraints. Global overrides apply:

- Stonegate: All magic is absorbed (blocked) regardless of spell.
- Palace of Blackthorn: All magic absorbed if the Avatar does not wear the Crown of Lord British.

Legend: Y = allowed, N = not allowed.

| Spell                | Overworld | Town | Dungeon | Combat | Exceptions |
|----------------------|:---------:|:----:|:-------:|:------:|------------|
| In Lor               |    Y      |  Y   |    Y    |   Y    | — |
| Vas Lor              |    Y      |  Y   |    Y    |   N    | — |
| Wis An Ylem          |    Y      |  N   |    N    |   N    | — |
| In Quas Wis (View)   |    Y      |  Y   |    Y    |   N    | — |
| In Por               |    Y      |  N   |    N    |   Y    | Combat: random legal tile; rooms off |
| Uus Por              |    N      |  N   |    Y    |   N    | Blocked in Doom |
| Des Por              |    N      |  N   |    Y    |   N    | Blocked in Doom |
| Vas Rel Por          |    Y      |  Y   |    N    |   N    | Does not end turn on success |
| In Sanct             |    Y      |  Y   |    Y    |   Y    | — |
| In An                |    Y      |  Y   |    Y    |   Y    | — |
| Rel Tym              |    N      |  N   |    N    |   Y    | — |
| An Tym               |    Y      |  Y   |    Y    |   Y    | Blocked in Doom/Stonegate |
| Grav Por             |    N      |  N   |    N    |   Y    | — |
| Vas Flam             |    N      |  N   |    N    |   Y    | — |
| Xen Corp             |    N      |  N   |    N    |   Y    | — |
| In Flam Grav         |    N      |  N   |    Y    |   Y    | — |
| In Nox Grav          |    N      |  N   |    Y    |   Y    | — |
| In Zu Grav           |    N      |  N   |    Y    |   Y    | — |
| In Sanct Grav        |    N      |  N   |    Y    |   Y    | — |
| Kal Xen              |    N      |  N   |    N    |   Y    | — |
| Kal Xen Corp         |    N      |  N   |    N    |   Y    | — |
| In Bet Xen           |    N      |  N   |    N    |   Y    | — |
| An Xen Ex            |    N      |  N   |    N    |   Y    | — |
| Quas An Wis          |    N      |  N   |    N    |   Y    | — |
| An Xen Corp          |    N      |  N   |    N    |   Y    | — |
| In Quas Corp         |    N      |  N   |    N    |   Y    | — |
| Rel Xen Bet          |    N      |  N   |    N    |   Y    | — |
| In Quas Xen          |    N      |  N   |    N    |   Y    | — |
| Wis Quas             |    N      |  N   |    N    |   Y    | — |
| In Ex Por            |    Y      |  Y   |    Y    |   N    | — |
| An Ex Por            |    Y      |  Y   |    Y    |   N    | — |
| An Sanct             |    Y      |  Y   |    Y    |   N    | — |
| An Grav              |    Y      |  Y   |    Y    |   N    | Dispel tile fields; not used in combat here |
| Mani                 |    Y      |  Y   |    Y    |   Y    | — |
| Vas Mani             |    Y      |  Y   |    Y    |   Y    | — |
| In Xen Mani          |    Y      |  Y   |    Y    |   Y    | — |
| In Vas Por Ylem      |    N      |  N   |    N    |   Y    | — |
| In Wis               |    Y      |  N   |    N    |   N    | — |
| In Nox Hur           |    N      |  N   |    N    |   Y    | — |
| In Flam Hur          |    N      |  N   |    N    |   Y    | — |
| In Vas Grav Corp     |    N      |  N   |    N    |   Y    | — |
| In Mani Corp         |    Y      |  Y   |    Y    |   Y    | Resurrection |
| Sanct Lor            |    N      |  N   |    N    |   Y    | — |

Notes:

- “Blocked in Doom” denotes spells disabled in the Doom dungeon (legacy `onmap==0x28`), per engine rules; confirm case‑by‑case as corresponding mechanics are implemented.
- All rows inherit global overrides (Stonegate/Blackthorn’s crown rule) regardless of Y/N in their context cells.

## Rel Hur (Change Wind)

```pseudocode
FUNCTION cast_rel_hur():
    show_message("Wind change!\n")
    dir = prompt_spell_direction() // N/S/E/W
    IF is_overworld_or_town() THEN wind_change(dir, immediate=TRUE) ELSE show_message("Not here!\n"); RETURN FALSE
    RETURN TRUE
ENDFUNCTION
```

## Rel Tym (Quickness)

Doubles the party’s action cadence by skipping alternate monster turns. Interacts with turn modifiers.

```pseudocode
FUNCTION cast_rel_tym():
    // Quickness: set active spell 'Q'; monsters may lose every other turn
    dur_spell(QUICKNESS, 30, sfx_level=5)
    RETURN TRUE
ENDFUNCTION
```

Notes:

- Combat Core: see Turn Modifiers — checks `active_spell == QUICKNESS` to throttle turns.

## Negate Time (Scroll)

Stops time briefly so monsters do not act. This is available via scroll usage.

```pseudocode
FUNCTION use_scroll_negate_time():
    IF in_special_map({Doom, Stonegate}) THEN show_message("No effect!\n"); glide(800,2000,1,50); RETURN FALSE
    show_message("Negate time!\n"); dur_spell(TIME_STOP, 20, sfx_level=7); RETURN TRUE
ENDFUNCTION
```

Notes:

- Combat Core: monsters do not act while `active_spell == TIME_STOP`.

## An Tym (Negate Time)

Identical effect to the scroll variant, applied via spellcasting.

```pseudocode
FUNCTION cast_an_tym():
    show_message("Negate time!\n")
    dur_spell(TIME_STOP, 20, sfx_level=7)
    RETURN TRUE
ENDFUNCTION
```

## View (Scroll)

Renders the same tactical map as the `View` command, but consumes a scroll instead of a gem.

```pseudocode
FUNCTION use_scroll_view():
    show_message("View!\n")
    IF in_combat() THEN show_message("Not here!\n"); RETURN FALSE
    IF is_dungeon() THEN dng_view() ELSE view_area(player_x, player_y)
    play_cast_sound(view)
    RETURN TRUE
ENDFUNCTION
```

See also: Commands → View (gem map): [Commands.md#view-gem-map](Commands.md#view-gem-map).

## In Quas Wis (View)

Spell-based tactical map view identical to the `View` command, without consuming a gem.

```pseudocode
FUNCTION cast_in_quas_wis():
    show_message("View!\n")
    IF is_dungeon() THEN dng_view() ELSE view_area(player_x, player_y)
    play_cast_sound(view)
    RETURN TRUE
ENDFUNCTION
```

See also: Commands → View and Spells → View (Scroll).

## In Lor (Light)

```pseudocode
FUNCTION cast_in_lor():
    // Short-lived magical light source
    play_cast_sound(light)
    magic_light = 64 + random(0, 31)
    RETURN TRUE
ENDFUNCTION
```

## An Nox (Cure Poison)

```pseudocode
FUNCTION cast_an_nox():
    plr = select_target_party_member(); IF plr < 0 THEN RETURN -1
    IF player[plr].status == POISONED THEN player[plr].status = GOOD; mark_stats_changed(); play_cast_sound(light); RETURN TRUE
    show_message("No effect!\n"); RETURN FALSE
ENDFUNCTION
```

See also: Combat Effects → Per-Turn Player Updates (poison damage), Environment → Fountain cure.

## An Zu (Awaken)

```pseudocode
FUNCTION cast_an_zu():
    plr = select_target_party_member(); IF plr < 0 THEN RETURN -1
    IF player[plr].status == SLEEP THEN player[plr].status = GOOD; mark_stats_changed(); play_cast_sound(light); RETURN TRUE
    show_message("No effect!\n"); RETURN FALSE
ENDFUNCTION
```

Notes:

- Wakes a sleeping party member. Does not affect monsters.

## In Bet Lor (Great Light)

```pseudocode
FUNCTION cast_in_bet_lor():
    // Long-lived magical light source
    play_cast_sound(light)
    magic_light = 255
    RETURN TRUE
ENDFUNCTION
```

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

## Vas Mani (Great Heal)

```pseudocode
FUNCTION cast_vas_mani():
    plr = select_target_party_member(); IF plr < 0 THEN RETURN -1
    IF player[plr].status == DEAD THEN RETURN FALSE
    heal_amount = vas_mani_heal_roll() // engine-defined; greater than Mani
    player[plr].hit_points = min(player[plr].hit_points + heal_amount, player[plr].hit_max)
    mark_stats_changed(); play_cast_sound(light)
    RETURN TRUE
ENDFUNCTION
```

Notes:

- Heal amount exceeds Mani; use engine’s `vas_mani()` roll. Both respect max HP and fail on Dead targets.

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

## Grav Por (Magic Missile)

Fires a basic magic projectile at a selected target, using weapon-style missile handling.

```pseudocode
FUNCTION cast_grav_por(target):
    // Uses engine's weapon_spell routing for projectiles
    return weapon_spell(MAGIC_MISSILE, target)
ENDFUNCTION
```

Notes:

- Legal in combat only; see spell flags. Projectile rules in Combat Effects → Missiles.
- Direct damage; bypasses intelligence save per legacy `saveint` rules.

## Vas Flam (Firebolt)

Throws a stronger fire-based projectile at a target.

```pseudocode
FUNCTION cast_vas_flam(target):
    return weapon_spell(FIREBOLT, target) // higher damage than Grav Por
ENDFUNCTION
```

Notes:

- Direct damage; bypasses intelligence save per legacy `saveint` rules.

## Field Spells (Create Fields)

Create a persistent field at a targeted tile. Field effects resolve per tick; see Combat Effects → Field Effects.

```pseudocode
FUNCTION cast_in_flam_grav(target_xy):
    return field_spell(FIRE_FIELD)

FUNCTION cast_in_nox_grav(target_xy):
    return field_spell(POISON_FIELD)

FUNCTION cast_in_zu_grav(target_xy):
    return field_spell(SLEEP_FIELD)

FUNCTION cast_in_sanct_grav(target_xy):
    return field_spell(ENERGY_FIELD)
ENDFUNCTION
```

Notes:

- Legal in combat/dungeons; exact flags per spell. Fields affect anyone standing on them according to field type.

## In Por (Blink)

Short-range teleport. In combat moves the caster to a legal nearby tile; outdoors teleports a short distance subject to terrain rules.

```pseudocode
FUNCTION cast_in_por():
    IF is_combat() THEN
        IF pick_random_legal_adjacent_tile(&tx, &ty) THEN move_caster_to(tx, ty); RETURN TRUE ELSE RETURN FALSE
    ELSE IF is_overworld() THEN
        dir = prompt_direction(); steps = random(1, 3); try_blink_overworld(dir, steps)
        RETURN TRUE
    ELSE RETURN FALSE
ENDFUNCTION
```

## In Wis (Sense)

Reveals information about the immediate area/world context (e.g., direction hints, notable features). Outdoors only.

```pseudocode
FUNCTION cast_in_wis():
    IF NOT is_overworld() THEN show_message("Not here!\n"); RETURN FALSE
    locate_position() // print coordinates/region info
    RETURN TRUE
ENDFUNCTION
```

## In Zu (Sleep)

Attempts to put a target to sleep in combat.

```pseudocode
FUNCTION cast_in_zu(target):
    IF save_vs_intelligence(caster, target, bonus=0) THEN RETURN FALSE
    apply_status(target, SLEEP); announce_name(target); show_message(" sleeps!\n")
    RETURN TRUE
ENDFUNCTION
```

## An Ylem (Disintegrate Object)

Removes a pushable object in a chosen direction (chairs, barrels, mirrors, etc.).

```pseudocode
FUNCTION cast_an_ylem():
    IF NOT getdir() THEN RETURN -1
    IF remove_pushable_object_at(player_x + tempx, player_y + tempy) THEN show_message("POOF!\\n"); fast_los_update(); glide() ELSE RETURN FALSE
    RETURN TRUE
ENDFUNCTION
```

See also: Objects → Pushables; this is the spell-driven removal.

## An Xen Corp (Repel Undead)

For each undead monster in combat: on failed INT save, set HP to 1 and force fleeing.

```pseudocode
FUNCTION cast_an_xen_corp():
    FOR each actor IN combat_queue:
        IF is_monster(actor) AND has_flag(actor, UNDEAD) AND NOT immune(actor) THEN
            IF NOT save_vs_intelligence(caster, actor, 0) THEN set_hp(actor, 1); set_flee(actor)
    RETURN TRUE
ENDFUNCTION
```

## Wis Quas (Reveal)

Reveals invisible monsters in combat.

```pseudocode
FUNCTION cast_wis_quas():
    FOR each actor IN combat_queue:
        IF is_monster(actor) AND is_invisible(actor) THEN clear_invisible(actor); restore_object_shape(actor)
    refresh_combat_view()
    RETURN TRUE
ENDFUNCTION
```

## In Bet Xen (Insect Plague)

Summons a small swarm of charmed insects near the caster if legal tiles exist.

```pseudocode
FUNCTION cast_in_bet_xen():
    IF pick_legal_summon_xy(INSECTS, &tx, &ty) THEN
        FOR i in 1..4: idx = createentry(INSECTS, 0, tx, ty, level); IF idx >= 0 THEN charm(idx)
        RETURN TRUE
    RETURN FALSE
ENDFUNCTION
```

## An Xen Ex (Charm)

Charms a targeted creature (monster or player) in combat on failed INT save; toggles charm state.

```pseudocode
FUNCTION cast_an_xen_ex():
    target = aim_creature_with_crosshair(max_range=15)
    IF target == NONE OR immune(target) OR loyalty(target) == 0 THEN RETURN FALSE
    IF save_vs_intelligence(caster, target, 0) THEN RETURN FALSE
    toggle_charm(target); RETURN TRUE
ENDFUNCTION
```

## Rel Xen Bet (Polymorph)

Polymorphs a valid target into a rat on success.

```pseudocode
FUNCTION cast_rel_xen_bet():
    target = aim_creature_with_crosshair(max_range=15)
    IF target == NONE OR immune(target) OR loyalty(target) == 0 THEN RETURN FALSE
    IF save_vs_intelligence(caster, target, 0) THEN RETURN FALSE
    set_actor_shape(target, RATS); RETURN TRUE
ENDFUNCTION
```

## In Quas Xen (Clone)

Creates a clone of a targeted player or monster.

```pseudocode
FUNCTION cast_in_quas_xen():
    target = aim_creature_with_crosshair(max_range=15)
    IF target == NONE OR immune(target) THEN RETURN FALSE
    RETURN clone_target_into_adjacent_tile(target)
ENDFUNCTION
```

## In Quas Corp (Fear)

Inflicts fear; targets flee if they fail resistance.

```pseudocode
FUNCTION cast_in_quas_corp():
    FOR each enemy IN visible_enemies():
        IF NOT immune(enemy) AND NOT save_vs_intelligence(caster, enemy, 0) THEN set_flee(enemy)
    RETURN TRUE
ENDFUNCTION
```

## Xen Corp (Death Bolt)

Fires a powerful projectile at a target.

```pseudocode
FUNCTION cast_xen_corp(target):
    return weapon_spell(DEATH_BOLT, target) // heavy-damage projectile
ENDFUNCTION
```

Notes:

- Direct damage; bypasses intelligence save per legacy `saveint` rules.

## Storm Spells (Area Effects)

Large area-of-effect damage/status spells; details are engine-tuned via `nukem`.

```pseudocode
FUNCTION cast_in_nox_hur():   return nukem(caster, variant=2, color=GREEN)
FUNCTION cast_in_flam_hur():  return nukem(caster, variant=3, color=RED)
FUNCTION cast_in_vas_grav_corp(): return nukem(caster, variant=4, color=BLUE)
```

Notes:

- Colors/variants map to damage type and radius; poison/fire/energy respectively.
- Area damage; bypasses intelligence save per legacy `saveint` rules.

## Intelligence Save Bypass (Reference)

The following damaging spells do not allow an INT save and should apply effects directly:

- Grav Por (Magic Missile / Flam Por)
- Vas Flam (Firebolt)
- In Nox Hur (Poison Storm)
- In Flam Hur (Firestorm)
- In Vas Grav Corp (Energy Storm)
- Xen Corp (Death Bolt)

## Scrolls: Light and Wind Change

```pseudocode
FUNCTION use_scroll_light():
    light_spell(duration=240); play_cast_sound(light); RETURN TRUE

FUNCTION use_scroll_wind_change():
    show_message("Wind change!\\n"); dir = prompt_spell_direction(); IF is_overworld_or_town() THEN wind_change(dir, TRUE) ELSE RETURN FALSE
    RETURN TRUE
ENDFUNCTION
```

## Scrolls: Summon Daemon and Resurrection

```pseudocode
FUNCTION use_scroll_summon_daemon():
    IF in_combat() THEN RETURN summon_daemon(is_scroll=TRUE)
    show_message("Not here!\n"); RETURN FALSE

FUNCTION use_scroll_resurrection():
    IF in_combat() THEN show_message("Not here!\n"); RETURN FALSE
    plr = select_target_party_member(); IF plr < 0 THEN RETURN -1
    RETURN resurrect(plr, is_scroll=TRUE)
ENDFUNCTION
```

## Scrolls Summary (At-a-Glance)

| Scroll            | Context                      | Effect                                    | Notes                                  |
|-------------------|------------------------------|-------------------------------------------|----------------------------------------|
| Light             | Any                           | Sets strong `magic_light`                 | `light_spell(240)`                     |
| Wind Change       | Overworld/Town                | Sets wind to chosen direction             | Uses `wind_change(dir, TRUE)`          |
| Protection        | Any                           | `dur_spell('P', 100)`                     | Longer than spell variant              |
| Negate Magic      | Any                           | `dur_spell('N', 20)`                      | Suppresses certain magic effects       |
| View              | Surface/Dungeon (not combat)  | Minimap or dungeon view                   | Mirrors `View` command                 |
| Summon Daemon     | Combat only                   | Summons a daemon; charmed                 | Fails “Not here!” outside combat       |
| Resurrection      | Surface/Town (not combat)     | Revive selected party member              | Fails “Not here!” in combat            |
| Negate Time       | Any except special maps       | `dur_spell(TIME_STOP, 20)`                | Blocked in Doom/Stonegate              |


## Quas An Wis (Mass Charm/Confusion)

Applies a temporary aura that can make monsters fight each other. While active, each monster performs an INT save when choosing a target; on a failed save, its effective allegiance flips so it treats other monsters as enemies for that decision.

```pseudocode
// Cast: set the global aura and duration
CONST MASS_CHARM_FLAG = 'C'

FUNCTION cast_quas_an_wis():
    // Duration matches legacy: 20 ticks; sfx_level from data
    dur_spell(MASS_CHARM_FLAG, duration=20, sfx_level=6)
    RETURN TRUE
ENDFUNCTION

// AI Integration: called from targeting (e.g., whosnear / select_nearest_enemy)
FUNCTION effective_loyalty_for_targeting(attacker_index):
    base = loyalty(attacker_index) // >0 for monsters, 0 for PCs
    IF active_spell == MASS_CHARM_FLAG AND is_monster(attacker_index) THEN
        // INT save: failure when d30 > INT (lower INT → more likely to fail)
        IF rolld30() > getattr(attacker_index, STAT_INT) THEN
            RETURN 0 // flip to “PC” side; now enemies are monsters (loyalty != 0)
        ENDIF
    ENDIF
    RETURN base
ENDFUNCTION

// Example: nearest-enemy selection loop (simplified)
FUNCTION pick_target(attacker_index):
    atk_loyalty = effective_loyalty_for_targeting(attacker_index)
    chosen = NONE; best = +INF
    FOR each i IN 0..31:
        IF i == attacker_index OR !is_alive(i) OR is_object(i) THEN CONTINUE
        IF side_of(atk_loyalty) == side_of(loyalty(i)) THEN CONTINUE // must be opposite sides
        IF is_invisible(i) AND NOT is_shadowlord(i) AND NOT in_special_map() THEN CONTINUE
        IF is_underground(i) THEN CONTINUE
        d = manhattan(attacker.pos, i.pos)
        IF d < best THEN best = d; chosen = i
    RETURN chosen
```

Behavior details:

- Scope: Affects monsters only; PCs are not flipped. Each monster rolls independently and re-rolls every time it chooses a target (i.e., per target selection), so behavior can vary turn-to-turn.
- Save: 1d30 vs INT; if roll > INT the monster fails and flips side for that decision.
- Targeting: Flipped monsters will not attack PCs while flipped (PCs’ loyalty is 0, which matches the flipped side and is excluded by the “opposite sides” check). They will target other monsters instead.
- Duration: Aura ends when `spell_dur` for 'C' reaches 0; behavior instantly reverts.
- Distinct from Charm: Unlike An Xen Ex (single-target charm toggle), Quas An Wis never sets the CHARM flag; it only perturbs targeting while active.

Testing tips:

- Seed PRNG and pit monsters with varying INT; log target selections with and without the aura to verify flip frequency ~ P(fail) = max(0, (30-INT)/30).
- Confirm invisibility and underground suppression still apply to chosen targets.

## Frotz (Unimplemented)

Present in the canonical spell list but not implemented in legacy behavior. Treat as reserved: do not expose to players, and assigning any effect risks diverging from authentic gameplay.

Notes:

- Keep inert/no-op in reimplementations to match original.
## Wis An Ylem (X-Ray Vision)

Grants X-ray vision that reveals features beyond walls; enhances Look/Search. Surface/towns only.

```pseudocode
FUNCTION cast_wis_an_ylem():
    IF is_dungeon() OR in_combat() THEN show_message("Not here!\n"); RETURN FALSE
    x_ray_vision(); RETURN TRUE
ENDFUNCTION
```

See also: Commands → Look, Dungeon → Search and Hidden Features, Potions → White.

## In An (Negate Magic)

Reduces or nullifies magical effects on the party for a short duration.

```pseudocode
FUNCTION cast_in_an():
    dur_spell(NEGATE_MAGIC, 10, sfx_level=6)
    RETURN TRUE
ENDFUNCTION
```

## In Sanct (Protection)

Applies a temporary protection buff that reduces incoming damage or improves saves (engine-defined).

```pseudocode
FUNCTION cast_in_sanct():
    dur_spell(PROTECTION, 20, sfx_level=4)
    RETURN TRUE
ENDFUNCTION
```

Notes:

- The exact mitigation model is data-driven; typically reduces damage or increases defense rolls.

## Scroll: Protection / Negate Magic

Scrolls that apply the same temporary effects without using MP.

```pseudocode
FUNCTION use_scroll_protection():
    show_message("Protection!\n")
    dur_spell(PROTECTION, 100, sfx_level=2) // longer than spell variant
    RETURN TRUE

FUNCTION use_scroll_negate_magic():
    show_message("Negate magic!\n")
    dur_spell(NEGATE_MAGIC, 20, sfx_level=3)
    RETURN TRUE
ENDFUNCTION
```

See also: Combat Core → Turn Modifiers and Effects for how these flags influence combat.

## Sanct Lor (Invisibility)

Turns the caster invisible in combat; AI targeting is affected.

```pseudocode
FUNCTION cast_sanct_lor():
    set_actor_shape(current_actor, INVIS_PC)
    set_flag(combatq[turn], INVISO_MASK)
    play_cast_sound(high)
    RETURN TRUE
ENDFUNCTION
```

Notes:

- Combat AI respects invisibility except for special cases (e.g., Shadowlords). See Combat AI → Special Moves/Targeting.

## Resurrect (Spell)

```pseudocode
FUNCTION cast_resurrect():
    plr = select_target_party_member(); IF plr < 0 THEN RETURN -1
    RETURN resurrect(plr, is_scroll=FALSE)
ENDFUNCTION
```

Notes:

- On success, restores a dead party member; consumes MP and mixed spell stock.

## Vas Rel Por (Gate Travel)

Long-range teleport; prompts for a destination and moves the party if legal.

```pseudocode
FUNCTION cast_vas_rel_por():
    dest = choose_travel_destination() // constrained by game rules (moongates/markers)
    IF dest == NONE THEN RETURN FALSE
    IF validate_travel_destination(dest) THEN perform_travel(dest); RETURN TRUE ELSE show_message("No effect!\n"); RETURN FALSE
ENDFUNCTION
```

Notes:

- Often restricted by quest flags or special maps; may not consume the player’s turn on success.

## An Grav (Dispel Field)

```pseudocode
FUNCTION cast_an_grav():
    // Dispel a targeted energy field (sleep/poison/fire/electric)
    success = an_grav(targeted=TRUE)
    IF success THEN show_message("Field dissolved!\n") ELSE show_message("No effect!\n")
    RETURN success
ENDFUNCTION
```

See also: Combat Effects → Field Effects, Commands → Use (Sceptre attempts An Grav first).

## An Sanct (Disarm/Unlock)

```pseudocode
FUNCTION cast_an_sanct():
    // Disarm dungeon chest or unlock standard locked door
    tx, ty = target_tile()
    IF disarm_or_unlock_target_tile(tx, ty) THEN RETURN TRUE ELSE show_message("No effect!\n"); RETURN FALSE
ENDFUNCTION
```

See also: Dungeon → Disarm/Unlock, Commands → Open.

## Uus Por (Up) and Des Por (Down)

```pseudocode
FUNCTION cast_uus_por():
    IF is_dungeon_map_of({Doom}) THEN RETURN FALSE
    neat_sound(view); IF newlvl(-1, TRUE) THEN exit_dng(); RETURN TRUE
    RETURN FALSE

FUNCTION cast_des_por():
    IF is_dungeon_map_of({Doom}) THEN RETURN FALSE
    neat_sound(view); IF newlvl(1, TRUE) THEN exit_dng(); RETURN TRUE
    RETURN FALSE
ENDFUNCTION
```

Notes:

- Only valid in dungeons; blocked in certain special dungeons.
- On success, performs the same level transition as a ladder/grate with magic flag.

## In Ex Por / An Ex Por (Magical Lock/Unlock)

```pseudocode
FUNCTION cast_in_ex_por(x, y):
    // Wizard-unlock magic-locked doors
    RETURN wizard_unlock_magic_at(x, y)

FUNCTION cast_an_ex_por(x, y):
    // Wizard-lock doors (apply magical lock)
    RETURN wizard_lock_door_at(x, y)
ENDFUNCTION
```

See also: Objects → Doors (Lock/Unlock), Commands → Use (Skull Key uses In Ex Por behavior).
