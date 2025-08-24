**Overview**
- **Purpose**: Central reference for enemies and how original game data is combined with hand-authored metadata to control spawning, movement, and behavior.
- **Sources**: Core stats and ability flags are read from the U5 OVL; project-specific metadata comes from `internal/references/data/AdditionalEnemyFlags.json` (authored via the adjacent `.csv`).

**EnemyReference**
- **Role**: Canonical, merged record for a single enemy. Built from raw OVL bytes plus `AdditionalEnemyFlags` and tile references.
- **Key fields**:
  - `KeyFrameTile`: Sprite/tile for the enemy's key animation frame.
  - `Armour`, `Damage`, `Dexterity`, `HitPoints`, `Intelligence`, `MaxPerMap`, `Strength`, `TreasureNumber`: Raw stats read from OVL.
  - `EnemyAbilities` (map): Bitflag abilities from OVL. Supported values include: Bludgeons, PossessCharm, Undead, DivideOnHit, Immortal, PoisonAtRange, StealsFood, NoCorpse, RangedMagic, Teleport, DisappearsOnDeath, Invisibility, GatesInDaemon, Poison, InfectWithPlague.
  - `AdditionalEnemyFlags`: Hand-authored metadata (see below).
  - `AttackRange` (1–9): Read from OVL; used in combat targeting logic (when implemented).
  - `Friend`: Link to another `EnemyReference` (friend/ally index in OVL), resolved after all references are created.
- **Behavior helpers**:
  - `CanSpawnToTile(tile)`: True if the enemy may spawn on the given tile, based on environment and passability.
  - `CanMoveToTile(tile)`: True if the enemy can move into the tile, honoring environment constraints and flags like flying/wall-passing.
  - `HasAbility(ability)`: Tests the OVL-derived ability bitflags.
- **Implementation notes**:
  - Creation: `internal/references/enemies.go` builds `EnemyReference` values from `newRawEnemyReferences` and tile refs, then resolves `Friend` back-references.
  - Tile logic: Environment checks rely on tile helpers (e.g., `IsWaterEnemyPassable`, `IsLandEnemyPassable`, `IsDesert`). Flying and wall-passing are layered in via `AdditionalEnemyFlags`.

**AdditionalEnemyFlags**
- **Role**: Hand-edited extension metadata per enemy. Used to constrain spawn/movement, drive world-spawn weighting, and annotate large-map missile behavior. Serialized as JSON and embedded via `go:embed`.
- **Edit flow**:
  - Author in `internal/references/data/AdditionalEnemyFlags.csv` for spreadsheet-friendly editing, then export/update `AdditionalEnemyFlags.json` (the runtime source).
  - JSON is read at startup (`internal/references/enemies_raw.go`) and unmarshaled into `AdditionalEnemyFlags`. The field `LargeMapMissile` string is converted to the internal `MissileType` during unmarshal.
- **Fields** (from `internal/references/enemies_raw.go`):
  - `Name` (string): Human-readable name; also helps verify ordering against sprite indices.
  - `Experience` (int): XP yield when defeated.
  - `IsWaterEnemy` (bool): Treated as a water creature for spawn/movement checks.
  - `IsSandEnemy` (bool): Treated as a desert/sand creature for spawn/movement checks.
  - `DoNotMove` (bool): Marker for immobile entities (e.g., fields, mimics); movement systems may respect this to freeze movement.
  - `CanFlyOverWater` (bool): Permits entering water tiles even if otherwise a land creature.
  - `CanPassThroughWalls` (bool): Permits entering wall tiles.
  - `ActivelyAttacks` (bool): Aggression hint for AI/aggro systems.
  - `LargeMapMissile` (string → `MissileType`): Missile type used on large maps (e.g., `Red`, `CannonBall`, `None`).
  - `LargeMapMissileRange` (int): Range for large-map missile usage.
  - `WaterWeight`, `DesertWeight`, `LandWeight`, `UnderworldWeight` (ints): Environment-based spawn weights used by the large-map AI.
- **Integration points**:
  - Spawn/move rules: `EnemyReference.CanSpawnToTile` and `CanMoveToTile` read `IsWaterEnemy`, `IsSandEnemy`, `CanFlyOverWater`, `CanPassThroughWalls` to gate tile eligibility.
  - Environment weights: `internal/ai/npc_ai_controller_large_map.go` uses the weights when selecting a monster for the detected environment (`Water`, `Desert`, `Land`, `Underworld`). Enemies with weight 0 are never chosen for that environment.
  - Large-map missiles: `LargeMapMissile` and `LargeMapMissileRange` define ranged behavior on overworld/sea maps for those enemies that can fire; the type is resolved via `references.GetMissileTypeFromString`.

**Data Shape and Ordering**
- `AdditionalEnemyFlags.json` contains one entry per enemy, in strict sprite/enemy order (currently 48 total). The index is used to align with OVL-derived stats during `newRawEnemyReferences` construction.
- Missing fields in JSON default to zero-values when unmarshaled (false/0). If a weight is omitted or zero, that enemy will not be selected in that environment.

**File Locations**
- Runtime data: `internal/references/data/AdditionalEnemyFlags.json` (embedded)
- Authoring sheet: `internal/references/data/AdditionalEnemyFlags.csv`
- Types and construction: `internal/references/enemy.go`, `internal/references/enemies.go`, `internal/references/enemies_raw.go`
- Missile types: `internal/references/combat.go`, `internal/references/missile_type.go`
- Large-map AI integration: `internal/ai/npc_ai_controller_large_map.go`

**Notes & Future Wiring**
- `DoNotMove` and `ActivelyAttacks` are primarily intent flags; movement/combat systems should honor these as those systems mature.
- Some fields (e.g., large-map missiles) are defined in `AdditionalEnemyFlags` and resolved at load; specific usage may expand as combat implementation lands.

**Example: JSON → Behavior**
- JSON snippet (Sea Serpent):
  - `{"Name":"SeaSerpent1/SEA SERPENTS", "IsWaterEnemy":true, "LargeMapMissile":"Red", "LargeMapMissileRange":4, "WaterWeight":15, "LandWeight":0, "DesertWeight":0, "UnderworldWeight":8}`
- Effects in code:
  - Spawning: `EnemyReference.CanSpawnToTile(tile)` returns true for water passable tiles; false for land/desert tiles due to `IsWaterEnemy`.
  - Movement: `CanMoveToTile(tile)` allows movement across water. No land movement unless combined with `CanFlyOverWater` (not set here).
  - Selection: In large-map generation, water environments favor Sea Serpents heavily (`WaterWeight=15`); they are excluded from land/desert (`0` weight).
  - Missiles: Large-map AI can use a `Red` missile up to range 4 when missile usage is integrated for overworld encounters.

