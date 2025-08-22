# Hidden Objects and Search Mechanics

Certain locations contain hidden objects that can be revealed by using the Search action at specific coordinates. Some secrets have conditions and cooldowns.

## Searchable Hidden Objects

```pseudocode
STRUCT HiddenObject {
    Tile: TileID         // object tile to spawn (e.g., SkeletonKey, GlassSword)
    MapID: byte          // location/map identifier
    Level: byte          // floor/level
    X: byte              // tile X
    Y: byte              // tile Y
    OnePerDay: bool      // if true, only once per in-game day
    Condition: func()bool // optional condition (e.g., player has zero keys)
}

LIST<HiddenObject> HiddenObjects
```

## Search Flow

```pseudocode
FUNCTION search_at(x, y):
    found = FALSE
    FOR i = 0 TO len(HiddenObjects)-1:
        h = HiddenObjects[i]
        IF h.MapID == current_map() AND h.Level == level AND h.X == x AND h.Y == y THEN
            // Conditional secrets
            IF i == SKULL_KEYS_INDEX AND day == last_skull_key_day THEN CONTINUE
            IF i == BT_CASTLE_KEYS_INDEX AND keys > 0 THEN CONTINUE
            IF i == GLASS_SWORD_INDEX AND inventory_count(GLASS_SWORD) > 0 THEN CONTINUE

            // Global per-object found bitmask; skip if previously found (except conditional resets below)
            IF NOT is_conditional_index(i) AND object_found_bit(i) THEN CONTINUE

            // Mark daily limitation for skull keys
            IF i == SKULL_KEYS_INDEX THEN last_skull_key_day = day

            spawn_hidden_object(h.Tile, h.X, h.Y, h.Level)
            fast_los_update(); print_found_message(h.Tile)
            set_object_found_bit(i)
            found = TRUE
            BREAK
        ENDIF
    ENDFOR
    IF NOT found THEN show_message("nothing of note.\n")
ENDFUNCTION
```

## Built-in Conditional Secrets (Examples)

- Keys in Castle: If the party has zero normal keys and searches a specific castle location, a small cache of keys appears there.
- Skull Keys (daily): A skull key cache becomes available once per in-game day at a specific location; after searching, it won’t appear again until the day changes.
- Glass Swords: If the party has none, a hidden glass sword may appear at a specific search spot.

## Notes

- Daily limits: Use a simple `last_skull_key_day` (byte) to prevent multiple pickups on the same day.
- Found bitmask: Non-conditional secrets use a bitmask to prevent spawning again after being found.
- Messages: Use object-specific messages when revealing hidden items for flavor.

