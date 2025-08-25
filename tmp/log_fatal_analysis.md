# Log.Fatal Analysis for Task #7

## ✅ KEEP AS FATAL (Legitimate Uses)

### Critical Data Loading (Embedded/Internal Files)
```go
/internal/text/fonts.go:28                    // Loading embedded font data - KEEP FATAL
/internal/sprites/sprite_utils.go:14          // Loading embedded sprite data - KEEP FATAL  
/internal/sprites/image_animations.go:26      // Loading embedded animation data - KEEP FATAL
```

### Original Game Data Parsing (User Files - Data Corruption)
```go
/internal/config/config.go:85                 // Cannot read DATA.OVL (core Ultima V file) - KEEP FATAL
/internal/references/npc_schedule.go:32       // Not enough data to create NPCReference schedule - KEEP FATAL
/internal/references/look.go:37              // Can't read look file - KEEP FATAL
/internal/references/tlk_talkscript.go:244   // Unexpected end of label section - KEEP FATAL
/internal/references/tlk_talkscript.go:325   // Question without answer in TLK data - KEEP FATAL
/internal/references/inventory_item.go:49    // Error unmarshaling JSON data - KEEP FATAL
/internal/references/enemies_raw.go:72       // Error unmarshaling JSON data - KEEP FATAL
/internal/references/tiles.go:52            // Error unmarshaling JSON data - KEEP FATAL
/internal/references/data_ovl.go:89          // Error reading compressed words - KEEP FATAL
/internal/references/data_ovl.go:98          // Error reading filenames - KEEP FATAL
/internal/references/npc_references.go:32    // Error loading NPC reference data - KEEP FATAL
/internal/references/tlk_references.go:34    // Error parsing talk data - KEEP FATAL
/internal/references/tlk_references.go:71    // Error loading talk file - KEEP FATAL
/internal/references/references.go:41        // Error when loading locations - KEEP FATAL
```

### Programming Logic Errors (Contract Violations)
```go
/internal/map_units/npc_vehicle.go:70        // Skiff quantity is 0, cannot decrement - KEEP FATAL
/internal/game_state/game_state.go:293      // RandomIntInRange: min > max - KEEP FATAL
/internal/references/item_stack.go:88       // Can't pop from empty stack - KEEP FATAL
/internal/references/item_stack.go:98       // Can't peek from empty stack - KEEP FATAL
/internal/references/vehicles.go:128        // Bad vehicle type - KEEP FATAL
```

### Critical System State Violations  
```go
/internal/game_state/game_state.go:131      // Exceeded large map tiles (bounds overflow) - KEEP FATAL
/internal/game_state/game_state.go:252      // Expected large map type, got different - KEEP FATAL
/internal/map_state/layered_map.go:243      // Entered nil position (memory safety) - KEEP FATAL
/internal/map_state/large_map.go:12         // Expected large map type, got different - KEEP FATAL
/internal/ai/npc_ai_controller_large_map.go:224 // Unexpected negative position - KEEP FATAL
```

### Main Function and Core Engine Failures
```go
/cmd/ultimav/main.go:37                     // Game engine startup failure - KEEP FATAL
/cmd/ultimav/gamescene.go:103              // Critical system initialization - KEEP FATAL
/cmd/ultimav/gamescene.go:121              // Failed to create MessageCallbacks - KEEP FATAL
/cmd/ultimav/gamescene.go:165              // Failed to create SystemCallbacks - KEEP FATAL
```

### Game State Loading (Save Game Corruption)
```go
/internal/game_state/game_state.go:92       // Error loading saved gam raw data - KEEP FATAL
/internal/game_state/game_state.go:111      // Error loading legacy save game from bytes - KEEP FATAL
```

## ❌ CONVERT TO SOFT ERRORS (Should Be Recoverable)

### UI State Management
```go
/cmd/ultimav/gamescene.go:239               // Debug dialog index not found - CONVERT TO SOFT ERROR
/internal/ui/widgets/dialog_stack.go:64    // Input dialog box not found - CONVERT TO SOFT ERROR  
/internal/ui/widgets/dialog_stack.go:72    // Debug dialog index not found - CONVERT TO SOFT ERROR
```

### Gameplay Logic Validation
```go
/internal/datetime/ultima_date.go:69        // Cannot advance more than 9 hours - CONVERT TO SOFT ERROR
/internal/ai/npc_ai_controller_small_map.go:303 // Unknown aiType - CONVERT TO SOFT ERROR
/internal/ai/npc_ai_controller_small_map.go:378 // Unknown aiType - CONVERT TO SOFT ERROR
```

### Config File Operations
```go
/internal/config/config.go:63              // Error writing config file - CONVERT TO SOFT ERROR
/internal/config/config.go:66              // Error reading config file - CONVERT TO SOFT ERROR
```

### Missing Game Object Data
```go
/internal/references/small_location_reference.go:183 // Missing ladder/stair near NPC - CONVERT TO SOFT ERROR
/internal/map_units/map_unit_details.go:55 // NPC has no path calculated - CONVERT TO SOFT ERROR
/internal/game_state/game_state.go:223     // Tried to remove NPC at position but failed - CONVERT TO SOFT ERROR
/internal/references/small_location_reference.go:123 // Missing max tiles - CONVERT TO SOFT ERROR
```

### Development/Debug Code
```go
/cmd/ultimav/gamescene_tiles.go:137        // Item should exist since we checked - CONVERT TO SOFT ERROR
/cmd/ultimav/gamescene_tiles.go:156        // Unexpected tile index for map unit - CONVERT TO SOFT ERROR  
/cmd/ultimav/gamescene_tiles.go:179        // Unexpected map unit index - CONVERT TO SOFT ERROR
/cmd/ultimav/gamescene_tiles.go:203        // Bad index - CONVERT TO SOFT ERROR
/internal/references/location_references.go:56 // "OOf" - CONVERT TO SOFT ERROR
/internal/references/location_references.go:103 // Unhandled small map type - CONVERT TO SOFT ERROR
/internal/references/tlk_references.go:62  // Unhandled default case for small map type - CONVERT TO SOFT ERROR
```

### Text Input/UI Component Errors  
```go
/internal/ui/widgets/textinput.go:133      // Text input component error - CONVERT TO SOFT ERROR
/internal/ui/widgets/textinput.go:142      // Text input component error - CONVERT TO SOFT ERROR
```

## Summary

| Category | Keep Fatal | Convert to Soft | Total |
|----------|------------|-----------------|-------|
| **Critical Data Loading** | 3 | 0 | 3 |
| **Game Data Parsing** | 13 | 0 | 13 |  
| **Programming Logic Errors** | 5 | 0 | 5 |
| **System State Violations** | 5 | 0 | 5 |
| **Core Engine/Startup** | 4 | 0 | 4 |
| **Save Game Loading** | 2 | 0 | 2 |
| **UI State Management** | 0 | 3 | 3 |
| **Gameplay Validation** | 0 | 3 | 3 |
| **Config Operations** | 0 | 2 | 2 |
| **Missing Game Data** | 0 | 4 | 4 |
| **Debug/Development** | 0 | 9 | 9 |
| **UI Components** | 0 | 2 | 2 |
| **TOTALS** | **32** | **23** | **55** |

## Task #7 Implementation Plan

### Phase 1: Add TODO Comments (Immediate)
- Add appropriate TODO comments to all 55 `log.Fatal` calls
- Use categories: `// TODO: KEEP FATAL - [reason]` or `// TODO: CONVERT TO SOFT ERROR - [reason]`

### Phase 2: Convert Soft Errors (Future Tasks)  
- Convert 23 soft error fatals to proper error handling
- Priority: UI State (3) > Gameplay Validation (3) > Debug/Development (9)
- Lower Priority: Config (2), Missing Data (4), UI Components (2)

### Benefits
- **Keep 32 legitimate fatals** that indicate serious problems (data corruption, programming bugs, system failures)
- **Convert 23 recoverable errors** to graceful error handling that improves user experience
- **Clear guidelines** for future development to prevent inappropriate fatal usage