# Algorithms Pseudocode Task Tracker

This tracker mirrors the coverage checklist and adds a place to note open items and follow-ups.

## Coverage Checklist

- [x] RNG primitives and usage notes
- [x] Encounter probability (overworld) and score calculation
- [x] Monster selection by terrain (water/land/underworld)
- [x] Weighted random selection (`random_monster`)
- [x] Monster placement/valid location selection (off-screen spawn)

- [x] Commands: Yell / Words of Power / Ship Sails
- [x] Commands: Look (surface + dungeon, special tiles, matrices)
- [x] Commands: Open (doors, chests, timed)
- [x] Commands: Push (forward/swap/blocked + matrix)
- [x] Commands: Get (tile/object/chest pickup + matrix)
- [x] Commands: Talk (freed NPC thanks + matrix)
- [x] Commands: Use (carpet, skull key, amulet, crown, sceptre, spyglass)
- [x] Commands: Fire (town cannons & ship broadsides)
- [x] Commands: Hole Up & Camp (repair/camping)
- [x] Commands: Exit/Enter
 - [x] Commands: Klimb (ladders/grates; small maps + dungeon)
 - [x] Commands: Pass Turn (Space)
 - [x] Commands: Search (surface + dungeon ahead)
 - [x] Commands: View (gem map)
 - [x] Commands: Ztats (party member stats)

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
- [x] Spells: Rel Hur, In Lor/Great Light, Rel Tym (Quickness), Negate Time (Scroll)
- [x] Spells: Uus Por/Des Por (dungeon up/down), An Sanct/An Grav, In Ex Por/An Ex Por

### Spells Coverage (Full List)

Reference: canonical spell names from the original data. We will check these off as we document each spell’s behavior, constraints, and links to subsystems. Items already covered are marked.

- [x] In Lor (Light)
- [x] Grav Por
- [x] An Zu
- [x] An Nox (Cure Poison)
- [x] Mani (Heal)
- [x] An Ylem
- [x] An Sanct
- [x] An Xen Cor
- [x] Rel Hur (Change Wind)
- [x] In Wis
- [x] Kal Xen (Summon Animal)
- [x] In Xen Mani (Create Food)
- [x] Vas Lor (Great Light)
- [x] Vas Flam
- [x] In Flam Grav
- [x] In Nox Grav
- [x] In Zu Grav
- [x] In Por
- [x] An Grav (Dispel Field)
 - [x] In Sanct
- [x] In Sanct G
- [x] Uus Por (Up)
- [x] Des Por (Down)
- [x] Wis Quas
- [x] In Bet Xen
- [x] An Ex Por (Wizard-Lock)
- [x] In Ex Por (Wizard-Unlock Magic)
- [x] Vas Mani (Great Heal)
- [x] In Zu (Sleep)
- [x] Rel Tym (Quickness)
- [x] In Vas Por Ylem (Earthquake)
- [x] Quas An Wis
- [x] In An
- [x] Wis An Ylem
- [x] An Xen Ex
- [x] Rel Xen Bet
- [x] Sanct Lor (Invisibility)
- [x] Xen Corp
- [x] In Quas Xen
- [x] In Quas Wis (View)
- [x] In Nox Hur
- [x] In Quas Corp
- [x] In Mani Corp (Resurrection)
- [x] Kal Xen Corp (Summon Daemon)
- [x] In Vas Grav Corp
- [x] In Flam Hur
- [x] Vas Rel Por
- [x] An Tym (Negate Time)
- [ ] Frotz (unimplemented/reserved)

Notes:

- Names above are canonical. Where our prior docs used alternate labels (e.g., Great Light), we will align and cross-link.

### Scrolls Coverage

- [ ] Light (sets light)
- [x] Wind Change (Rel Hur)
- [x] Protection (dur_spell 'P')
- [x] Negate Magic (dur_spell 'N')
- [x] View (minimap/dungeon view)
- [x] Summon Daemon
- [x] Resurrection
- [x] Negate Time

### Potions Coverage

- [x] Blue (Cure Sleep)
- [x] Yellow (Heal)
- [x] Red (Cure Poison)
- [x] Green (Poison)
- [x] Orange (Sleep)
- [x] Purple (Polymorph to Rat, combat)
- [x] Black (Invisible, combat)
- [x] White (X-Ray, surface)

- [x] AI seek movement (combat targeting and stepping)
- [x] Overworld movement quirks (whirlpools, pirate ships)
- [x] Wind system (UI, change odds, Rel Hur)
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
- [x] Light sources & vision (torch vs magic light; persistence)
- [x] Light spells (In Lor, In Bet Lor)
- [x] Town drawbridges/portcullis at night
- [x] Town guard activation and karma effects (attacks, cannons)
- [x] Special guard behavior (alarm, pursuit, CallGuards/GoToJail)
- [x] Jail flow (relocation, cell setup, escape options)
- [x] Town/indoor hazards on small maps (fireplace, lava)
- [x] Small-town fixtures defaults table (fountains, wells, lamps, telescope)
- [x] Hidden objects and search mechanics (daily skull keys, conditional spawns)
- [x] Hazard matrix table (overworld/town/combat for SWAMP/LAVA/FIREPLACE)
- [x] Shop pricing schema examples (base + location multipliers)
- [x] Moongate phase timing and stones template
- [x] Bridge Trolls corrections (12.5%, on-foot only, toll or combat)
- [x] Town cannon rules (line-of-sight, friendly-fire, aggression)
- [x] Fireplace behavior correction (town step causes Burning + damageparty)

## In Progress / Open Items

- [ ] Doors Overview doc: consolidate door state transitions + UI and interactions (Open/Jimmy/Skull/Use).
- [ ] Doc cross-references per section to internal source (no `OLD/*` in public docs).
- [ ] Shop pricing data: per-town multipliers and item lists (full tables).
- [ ] Testing guidance snippets per module (seeded PRNG, deterministic checks).
- [x] Expand spell list coverage (targeting, durations, town/overworld vs dungeon constraints).
- [ ] NPC schedule transitions (time-of-day) deeper dive.
- [ ] Fixture-specific overrides by location (wells/fountains/lamps); expand Wish outcomes.
- [ ] Rare fixtures/objects edge cases (as discovered).

## Notes

- Keep new pseudocode deterministic and ensure all RNG calls are compatible with the central PRNG.
- If you add a new pseudocode block, link it from `README.md` and tick it here.
- [x] Door state transitions table
- [x] Special moves summary table
- [x] Pirate ship speed matrix template
- [x] Dungeon interactions tables (chests/traps/search)
- [x] Jimmy and Skull Key logic (OLD implementation) recorded
- [x] Open command (timed doors, chest flow) documented
- [x] Talk outcomes matrix (surface/overworld) documented
- [x] Hole up & camp (repair/camp flows) documented
- [x] Enter, Ignite Torch, and New Order commands documented
