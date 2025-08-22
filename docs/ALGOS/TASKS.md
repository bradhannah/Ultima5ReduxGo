# Algorithms Pseudocode Task Tracker

This tracker mirrors the coverage checklist and adds a place to note open items and follow-ups.

## Coverage Checklist

- [x] RNG primitives and usage notes
- [x] Encounter probability (overworld) and score calculation
- [x] Monster selection by terrain (water/land/underworld)
- [x] Weighted random selection (`random_monster`)
- [x] Monster placement/valid location selection (off-screen spawn)

- [x] Hit calculation (physical/magical)
- [x] Damage calculation (weapons/armor; Glass Sword)
- [x] Combat turn modifiers (time stop, quickness, riding)
- [x] Treasure drop and trapped chest checks
- [x] Ship damage (enemy hits) and player cannon damage
- [x] Pirate cannon fire
- [x] Poison on hit from monsters with `POISONS`
- [x] Monster fleeing thresholds and 1/64 chance at <50% HP
- [x] In-combat field effects (lava/fireplace/fire/poison/sleep)
- [x] Monster special moves: Possess/Invisible/Gate Daemon
- [x] Monster teleport behavior (boxed-in or 75% otherwise)
- [x] Player flee/try-move/klimb in combat (same-exit rule)
- [x] Regurgitation (`throwup`) when swallowed

- [x] Intelligence save (`saveint`)
- [x] Magical field expiration (`fieldkill`)
- [x] Spells: Kal Xen, In Xen Mani, In Vas Por Ylem, Mani, Kal Xen Corp

- [x] AI seek movement (combat targeting and stepping)
- [x] Overworld movement quirks (whirlpools, pirate ships)
- [x] Random movement (`rndmove`) and `moveit` terrain throttle
- [x] Terrain-based movement chance helper
- [x] Movement permissions (onfoot/iswater/legalmove)

- [x] Bridge Trolls encounter (chance, sneak, toll/combat)
- [x] Overworld hazards (swamp/lava)
- [x] Waterfalls
- [x] Underworld earthquakes
- [x] Darkness (sight suppression)
- [x] Moongates
- [x] Torch duration (dungeon)
- [x] Town drawbridges/portcullis at night
- [x] Town guard activation and karma effects (attacks, cannons)
- [x] Town/indoor hazards on small maps (fireplace, lava)
- [x] Small-town fixtures defaults table (fountains, wells, lamps, telescope)
- [x] Hidden objects and search mechanics (daily skull keys, conditional spawns)
- [x] Hazard matrix table (overworld/town/combat for SWAMP/LAVA/FIREPLACE)
- [x] Shop pricing schema examples (base + location multipliers)

## Open Items / Nice-to-Haves

- [ ] Doc index cross-references for each section (internal source pointers).
- [ ] Pricing schema examples for shops (data table format + example).
- [ ] NPC schedule transitions (time-of-day) deeper dive.
- [ ] Rare fixtures/objects edge cases (if discovered).

## Notes

- Keep new pseudocode deterministic and ensure all RNG calls are compatible with the central PRNG.
- If you add a new pseudocode block, link it from `README.md` and tick it here.
- [x] Door state transitions table
- [x] Special moves summary table
- [x] Pirate ship speed matrix template
- [x] Dungeon interactions tables (chests/traps/search)
