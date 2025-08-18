# SAVED.GAM File Structure (Ultima V)

This document describes the structure of the SAVED.GAM file as defined in the original fan-made port (U5SAVED.ASM) and the canonical Ultima V format. Each entry includes the byte offset, Go struct variable, data format, size, a brief description, and any constraints or notes.

---

## Character Records

**Note:** The SAVED.GAM file contains 16 character records, each 32 bytes. Only the first 6 records are used for the active party, but all 16 are present in the file. This matches the original Ultima V format and ensures compatibility with external tools and references.

**Special Note:** For all equipment fields (helmet, armor, weapons, rings, amulets), 0xFF means "none equipped" or "not present." For inventory items and quest items, 0xFF means "not in inventory" or "not present." Reserved bytes are typically set to 7.

| Offset | Type/Format                                      | Size (bytes) | Description/Usage                        | Constraints                |
|--------|--------------------------------------------------|--------------|------------------------------------------|----------------------------|
| 0x0000 | 16-bit unsigned                                  | 2            | Save file version or marker              | Usually 0 or 1             |
| 0x0002 | [NPlayers]PlayerCharacter                        | 512          | 16 character records x 32 bytes each     | See PlayerCharacter struct |
| 0x0202 | ItemQuantityLarge                                | 2            | Party food supply                        | 0-9999                     |
| 0x0204 | ItemQuantityLarge                                | 2            | Party gold                               | 0-9999                     |
| 0x0206 | ItemQuantitySmall                                | 1            | Number of keys                           | 0-99                       |
| 0x0207 | ItemQuantitySmall                                | 1            | Number of gems                           | 0-99                       |
| 0x0208 | ItemQuantitySmall                                | 1            | Number of torches                        | 0-99                       |
| 0x020B | ItemQuantitySmall                                | 1            | Number of skull keys                     | 0-99                       |
| 0x020C | ItemQuantitySmall                                | 1            | Skull key usage day                      | 0-31                       |
| 0x020D | ItemQuantitySmall                                | 1            | Has British's Amulet                     | 0 or 1                     |
| 0x020E | ItemQuantitySmall                                | 1            | Has Crown                                | 0 or 1                     |
| 0x020F | ItemQuantitySmall                                | 1            | Has Sceptre                              | 0 or 1                     |
| 0x0210 | InventoryQuantities[Shard, *ItemQuantitySmall]   | 4            | Shards of Mondain, Minax, Exodus, etc.   | 0 or 1 per shard           |
| 0x0214 | ItemQuantitySmall                                | 1            | Has Spyglass                             | 0 or 1                     |
| 0x0215 | ItemQuantitySmall                                | 1            | Has Plans                                | 0 or 1                     |
| 0x0216 | ItemQuantitySmall                                | 1            | Has Sextant                              | 0 or 1                     |
| 0x0217 | ItemQuantitySmall                                | 1            | Has Watch                                | 0 or 1                     |
| 0x0218 | ItemQuantitySmall                                | 1            | Has Badge                                | 0 or 1                     |
| 0x0219 | ItemQuantitySmall                                | 1            | Has Wooden Box                           | 0 or 1                     |
| 0x021A | InventoryQuantities[Weapon, *ItemQuantitySmall]  | 48           | Weapons inventory                        | 0-99 per weapon            |
| 0x024A | InventoryQuantities[Spell, *ItemQuantitySmall]   | 48           | Spells inventory                         | 0-99 per spell             |
| 0x027A | InventoryQuantities[Scroll, *ItemQuantitySmall]  | 8            | Scrolls inventory                        | 0-99 per scroll            |
| 0x0282 | InventoryQuantities[Potion, *ItemQuantitySmall]  | 8            | Potions inventory                        | 0-99 per potion            |
| 0x028A | [8]byte                                          | 8            | Stones X coords                          | 0-255                      |
| 0x0292 | [8]byte                                          | 8            | Stones Y coords                          | 0-255                      |
| 0x029A | [8]byte                                          | 8            | Stones M coords                          | 0-255                      |
| 0x02A2 | [8]byte                                          | 8            | Stones L coords                          | 0-255                      |
| 0x02AA | InventoryQuantities[Reagent, *ItemQuantitySmall] | 8            | Reagents inventory                       | 0-99 per reagent           |
| 0x02B2 | [3]byte                                          | 3            | Reagent days left                        | 0-255                      |
| 0x02B5 | 8-bit unsigned                                   | 1            | Number of party members                  | 1-6                        |
| 0x02B6 | [24]byte                                         | 24           | Objects found flags (bitfield)           | 0 or 1 per object          |
| 0x02CE | 16-bit unsigned                                  | 2            | In-game year                             | 0-9999                     |
| 0x02D0 | 16-bit unsigned                                  | 2            | Temporary X coordinate (map/party)       | 0-255                      |
| 0x02D2 | 16-bit unsigned                                  | 2            | Temporary Y coordinate (map/party)       | 0-255                      |
| 0x02D4 | 8-bit unsigned                                   | 1            | Currently active spell                   | 0-255                      |
| 0x02D5 | 8-bit unsigned                                   | 1            | Index of active player                   | 0-5                        |
| 0x02D6 | 8-bit unsigned                                   | 1            | Player sprite/shape index                | 0-255                      |
| 0x02D7 | 8-bit unsigned                                   | 1            | In-game month                            | 1-12                       |
| 0x02D8 | 8-bit unsigned                                   | 1            | In-game day                              | 1-31                       |
| 0x02D9 | 8-bit unsigned                                   | 1            | In-game hour                             | 0-23                       |
| 0x02DA | 8-bit unsigned                                   | 1            | Previous hour (for time events)          | 0-23                       |
| 0x02DB | 8-bit unsigned                                   | 1            | In-game minute                           | 0-59                       |
| 0x02DC | 8-bit unsigned                                   | 1            | Turn counter                             | 0-255                      |
| 0x02DD | 8-bit unsigned                                   | 1            | Update counter                           | 0-255                      |
| 0x02DE | 8-bit unsigned                                   | 1            | Gong event flag                          | 0 or 1                     |
| 0x02DF | 8-bit unsigned                                   | 1            | Trammel moon phase                       | 0-7                        |
| 0x02E0 | 8-bit unsigned                                   | 1            | Felucca moon phase                       | 0-7                        |
| 0x02E1 | 8-bit unsigned                                   | 1            | Moongate stage                           | 0-255                      |
| 0x02E2 | 8-bit unsigned                                   | 1            | Avatar's karma                           | 0-99                       |
| 0x02E3 | 8-bit unsigned                                   | 1            | Eat timer                                | 0-255                      |
| 0x02E4 | 8-bit unsigned                                   | 1            | Food consumption counter                 | 0-255                      |
| 0x02E5 | 8-bit unsigned                                   | 1            | Last NPC given to                        | 0-255                      |
| 0x02E6 | 8-bit unsigned                                   | 1            | Last camped location                     | 0-255                      |
| 0x02E7 | 8-bit unsigned                                   | 1            | Last wander event                        | 0-255                      |
| 0x02E8 | 8-bit unsigned                                   | 1            | Spell duration                           | 0-255                      |
| 0x02E9 | 8-bit unsigned                                   | 1            | Last spell cast                          | 0-255                      |
| 0x02EA | 8-bit unsigned                                   | 1            | Magic state/flag                         | 0-255                      |
| 0x02EB | 8-bit unsigned                                   | 1            | Motion state/flag                        | 0-255                      |
| 0x02EC | 8-bit unsigned                                   | 1            | Wind direction/strength                  | 0-255                      |
| 0x02ED | 8-bit unsigned                                   | 1            | Current map index                        | 0-255                      |
| 0x02EE | 8-bit unsigned                                   | 1            | Previous map index                       | 0-255                      |
| 0x02EF | 8-bit unsigned                                   | 1            | Dungeon level                            | 0-255                      |
| 0x02F0 | 8-bit unsigned                                   | 1            | Player X position                        | 0-255                      |
| 0x02F1 | 8-bit unsigned                                   | 1            | Player Y position                        | 0-255                      |
| 0x02F2 | 8-bit unsigned                                   | 1            | Crosshair state                          | 0-255                      |
| 0x02F3 | 8-bit unsigned                                   | 1            | Crosshair X                              | 0-255                      |
| 0x02F4 | 8-bit unsigned                                   | 1            | Crosshair Y                              | 0-255                      |
| 0x02F5 | 8-bit unsigned                                   | 1            | Map X offset                             | 0-255                      |
| 0x02F6 | 8-bit unsigned                                   | 1            | Map Y offset                             | 0-255                      |
| 0x02F7 | 8-bit unsigned                                   | 1            | Current weapon index                     | 0-255                      |
| 0x02F8 | 8-bit unsigned                                   | 1            | Turn number                              | 0-255                      |
| 0x02F9 | 8-bit unsigned                                   | 1            | Turn light state                         | 0-255                      |
| 0x02FA | 8-bit unsigned                                   | 1            | Exit direction                           | 0-255                      |
| 0x02FB | 8-bit unsigned                                   | 1            | Scenario index                           | 0-255                      |
| 0x02FC | 8-bit unsigned                                   | 1            | Wound type                               | 0-255                      |
| 0x02FD | 8-bit unsigned                                   | 1            | Victory flag                             | 0 or 1                     |
| 0x02FE | 8-bit unsigned                                   | 1            | Update flag                              | 0 or 1                     |
| 0x02FF | 8-bit unsigned                                   | 1            | Sight state                              | 0-255                      |
| 0x0300 | 8-bit unsigned                                   | 1            | Magic light state                        | 0-255                      |
| 0x0301 | 8-bit unsigned                                   | 1            | Torch light state                        | 0-255                      |
| 0x0302 | [32]byte                                         | 32           | Enemy data (IDs, states, etc.)           | 0-255 per enemy            |
| 0x0322 | [3]byte                                          | 3            | Lords met/found flags                    | 0 or 1 per lord            |
| 0x0325 | 8-bit unsigned                                   | 1            | Last called NPC                          | 0-255                      |
| 0x0326 | [2]byte                                          | 2            | Quests in progress (bitfield)            | 0 or 1 per quest           |
| 0x0328 | [2]byte                                          | 2            | Completed quests (bitfield)              | 0 or 1 per quest           |
| 0x032A | [8]byte                                          | 8            | Open dungeons flags (bitfield)           | 0 or 1 per dungeon         |
| 0x0332 | [8]byte                                          | 8            | Closed shrines flags (bitfield)          | 0 or 1 per shrine          |
| 0x033A | [14]byte                                         | 14           | Dungeon rooms visited (bitfield)         | 0 or 1 per room            |
| 0x0348 | [32]byte                                         | 32           | Bridge X positions                       | 0-255                      |
| 0x0368 | [32]byte                                         | 32           | Bridge Y positions                       | 0-255                      |
| 0x0388 | [32]byte                                         | 32           | Bridge tile types                        | 0-255                      |
| 0x03A8 | 8-bit unsigned                                   | 1            | Number of bridges                        | 0-32                       |
| 0x03A9 | 8-bit unsigned                                   | 1            | Door type                                | 0-255                      |
| 0x03AA | 8-bit unsigned                                   | 1            | Door X position                          | 0-255                      |
| 0x03AB | 8-bit unsigned                                   | 1            | Door Y position                          | 0-255                      |
| 0x03AC | 8-bit unsigned                                   | 1            | Door timer                               | 0-255                      |
| 0x03AD | 8-bit unsigned                                   | 1            | Ship X position                          | 0-255                      |
| 0x03AE | 8-bit unsigned                                   | 1            | Ship Y position                          | 0-255                      |
| 0x03AF | 8-bit unsigned                                   | 1            | Sailing state                            | 0 or 1                     |
| 0x03B0 | 8-bit unsigned                                   | 1            | Prompt state                             | 0-255                      |
| 0x03B1 | 8-bit unsigned                                   | 1            | Smashed object flag                      | 0 or 1                     |
| 0x03B2 | 8-bit unsigned                                   | 1            | Shadowlord present flag                  | 0 or 1                     |
| 0x03B3 | 8-bit unsigned                                   | 1            | Easy mode flag                           | 0 or 1                     |
| 0x03B4 | [512]byte                                        | 512          | Dungeon map data                         | 0-255 per cell             |
| 0x05B4 | [128]byte                                        | 128          | NPC dead flags (bitfield, 1 bit per NPC) | 0 or 1 per bit             |
| 0x0634 | [128]byte                                        | 128          | NPC met flags (bitfield, 1 bit per NPC)  | 0 or 1 per bit             |
| 0x06B4 | [256]byte                                        | 256          | Object state data                        | 0-255 per object           |
| 0x07B4 | 16-bit unsigned                                  | 2            | Last X position (map/town)               | 0-255                      |
| 0x07B6 | 16-bit unsigned                                  | 2            | Last Y position (map/town)               | 0-255                      |
| 0x07B8 | [512]byte                                        | 512          | NPC schedule data                        | 0-255 per entry            |
| 0x09B8 | [512]byte                                        | 512          | NPC stat data                            | 0-255 per entry            |
| 0x0BB8 | [1024]byte                                       | 1024         | NPC path data                            | 0-255 per entry            |
| 0x0FB8 | [64]word                                         | 64           | NPC pointers (32 x 2 bytes)              | 0-65535                    |
| 0x0FF8 | [32]byte                                         | 32           | NPC shape indices                        | 0-255                      |
| 0x1018 | 8-bit unsigned                                   | 1            | NPC command                              | 0-255                      |
| 0x1019 | 8-bit unsigned                                   | 1            | NPC number                               | 0-255                      |
| 0x101A | 16-bit unsigned                                  | 2            | Current dungeon level                    | 0-255                      |
| 0x101C | [64]byte                                         | 64           | Pause state data                         | 0-255 per entry            |
| 0x105C | 8-bit unsigned                                   | 1            | Last direction                           | 0-255                      |
| 0x105D | 8-bit unsigned                                   | 1            | Face direction                           | 0-255                      |
| 0x105E | 8-bit unsigned                                   | 1            | Dungeon type                             | 0-255                      |
| 0x105F | 8-bit unsigned                                   | 1            | Skiff count                              | 0-255                      |
| 0x1060 | 16-bit unsigned                                  | 2            | Last saved marker                        | 0-65535                    |
| 0x1062 | [1024]byte                                       | 1024         | Map/town data (union)                    | 0-255 per cell             |

**Note:** This table is a partial extract. The full file includes all fields as defined in U5SAVED.ASM, up to the end of the structure. For a complete and detailed breakdown, continue the pattern above for all fields in the file.

## Field Types
- `db` = 8-bit unsigned (byte)
- `dw` = 16-bit unsigned (word)
- `dup (?)` = array of the given type and length

## Bitfields
- Fields like `npc_dead` and `met_player` are 128-byte bitfields, each bit representing the state of a specific NPC.

## Constraints
- Most inventory and flag fields are 0-99 or 0/1.
- Array fields represent collections (e.g., party members, inventory, NPCs).
- Some fields are reserved or unused (see code for details).

---

This documentation is based on the original U5SAVED.ASM source. For further details, consult the source or Ultima V technical references.

## Character Record Structure (32 bytes)

| Offset | Name         | Type/Format     | Size (bytes) | Description/Usage                                                                  | Constraints        |
|--------|--------------|-----------------|--------------|------------------------------------------------------------------------------------|--------------------|
| 0x00   | name         | char[9]         | 9            | Character name (zero-terminated string)                                            | ASCII, max 8 chars |
| 0x09   | gender       | 8-bit unsigned  | 1            | Gender (0xB=male, 0xC=female)                                                      | 0xB or 0xC         |
| 0x0A   | class        | 8-bit unsigned  | 1            | Class ('A'=Avatar, 'B'=Bard, etc.)                                                 | ASCII              |
| 0x0B   | status       | 8-bit unsigned  | 1            | Status ('G'=Good, etc.)                                                            | ASCII              |
| 0x0C   | strength     | 8-bit unsigned  | 1            | Strength stat                                                                      | 1-30               |
| 0x0D   | dexterity    | 8-bit unsigned  | 1            | Dexterity stat                                                                     | 1-30               |
| 0x0E   | intelligence | 8-bit unsigned  | 1            | Intelligence stat                                                                  | 1-30               |
| 0x0F   | mp           | 8-bit unsigned  | 1            | Current MP                                                                         | 0-30               |
| 0x10   | hp           | 16-bit unsigned | 2            | Current hit points                                                                 | 1-240              |
| 0x12   | max_hp       | 16-bit unsigned | 2            | Maximum hit points                                                                 | 1-240              |
| 0x14   | xp           | 16-bit unsigned | 2            | Experience points                                                                  | 0-9999             |
| 0x16   | level        | 8-bit unsigned  | 1            | Character level                                                                    | 1-8                |
| 0x17   | months_inn   | 8-bit unsigned  | 1            | Months at inn                                                                      | 0-25               |
| 0x18   | reserved     | 8-bit unsigned  | 1            | Reserved (usually 7)                                                               | 7                  |
| 0x19   | helmet       | 8-bit unsigned  | 1            | Helmet equipped (item index or 0xFF)                                               | 0-0x2F, 0xFF       |
| 0x1A   | armor        | 8-bit unsigned  | 1            | Armor equipped (item index or 0xFF)                                                | 0-0x2F, 0xFF       |
| 0x1B   | left_hand    | 8-bit unsigned  | 1            | Weapon/shield left (item index or 0xFF)                                            | 0-0x2F, 0xFF       |
| 0x1C   | right_hand   | 8-bit unsigned  | 1            | Weapon/shield right (item index or 0xFF)                                           | 0-0x2F, 0xFF       |
| 0x1D   | ring         | 8-bit unsigned  | 1            | Ring equipped (item index or 0xFF)                                                 | 0-0x2F, 0xFF       |
| 0x1E   | amulet       | 8-bit unsigned  | 1            | Amulet equipped (item index or 0xFF)                                               | 0-0x2F, 0xFF       |
| 0x1F   | inn_party    | 8-bit unsigned  | 1            | Inn/party flag (0=in party, 0xFF=not joined, 0x7F=killed, else=inn location index) | n/a                |

**Extra Details:**
- Gender: 0xB = male, 0xC = female
- Class: ASCII code for class letter
- Status: ASCII code for status letter
- Reserved byte at 0x18 is always 7
- 0x1F: 0 = in party, 0xFF = not joined, 0x7F = killed, else = inn location index

## Player Record Structure

Each player record in the `character` array is 32 bytes. The structure is as follows:

| Offset | Name         | Type/Format     | Size (bytes) | Description/Usage                                                                  | Constraints        |
|--------|--------------|-----------------|--------------|------------------------------------------------------------------------------------|--------------------|
| 0x00   | name         | char[9]         | 9            | Character name (zero-terminated string)                                            | ASCII, max 8 chars |
| 0x09   | gender       | 8-bit unsigned  | 1            | Gender (0xB=male, 0xC=female)                                                      | 0xB or 0xC         |
| 0x0A   | class        | 8-bit unsigned  | 1            | Class ('A'=Avatar, 'B'=Bard, etc.)                                                 | ASCII              |
| 0x0B   | status       | 8-bit unsigned  | 1            | Status ('G'=Good, etc.)                                                            | ASCII              |
| 0x0C   | strength     | 8-bit unsigned  | 1            | Strength stat                                                                      | 1-30               |
| 0x0D   | dexterity    | 8-bit unsigned  | 1            | Dexterity stat                                                                     | 1-30               |
| 0x0E   | intelligence | 8-bit unsigned  | 1            | Intelligence stat                                                                  | 1-30               |
| 0x0F   | mp           | 8-bit unsigned  | 1            | Current MP                                                                         | 0-30               |
| 0x10   | hp           | 16-bit unsigned | 2            | Current hit points                                                                 | 1-240              |
| 0x12   | max_hp       | 16-bit unsigned | 2            | Maximum hit points                                                                 | 1-240              |
| 0x14   | xp           | 16-bit unsigned | 2            | Experience points                                                                  | 0-9999             |
| 0x16   | level        | 8-bit unsigned  | 1            | Character level                                                                    | 1-8                |
| 0x17   | months_inn   | 8-bit unsigned  | 1            | Months at inn                                                                      | 0-25               |
| 0x18   | reserved     | 8-bit unsigned  | 1            | Reserved (usually 7)                                                               | 7                  |
| 0x19   | helmet       | 8-bit unsigned  | 1            | Helmet equipped (item index or 0xFF)                                               | 0-0x2F, 0xFF       |
| 0x1A   | armor        | 8-bit unsigned  | 1            | Armor equipped (item index or 0xFF)                                                | 0-0x2F, 0xFF       |
| 0x1B   | left_hand    | 8-bit unsigned  | 1            | Weapon/shield left (item index or 0xFF)                                            | 0-0x2F, 0xFF       |
| 0x1C   | right_hand   | 8-bit unsigned  | 1            | Weapon/shield right (item index or 0xFF)                                           | 0-0x2F, 0xFF       |
| 0x1D   | ring         | 8-bit unsigned  | 1            | Ring equipped (item index or 0xFF)                                                 | 0-0x2F, 0xFF       |
| 0x1E   | amulet       | 8-bit unsigned  | 1            | Amulet equipped (item index or 0xFF)                                               | 0-0x2F, 0xFF       |
| 0x1F   | inn_party    | 8-bit unsigned  | 1            | Inn/party flag (0=in party, 0xFF=not joined, 0x7F=killed, else=inn location index) | n/a                |

**Note:** Only the first 6 character records are used for the active party. The remaining records may be unused, reserved, or for future expansion. This matches the original file format and ensures compatibility with external tools and references.

## Object Array Structure

The `object` array in SAVED.GAM is 256 bytes. Each entry represents a world object (chest, door, item, etc.).

| Offset | Name      | Type/Format     | Size (bytes) | Description/Usage                  | Constraints         |
|--------|-----------|-----------------|--------------|------------------------------------|---------------------|
| 0x00   | type      | 8-bit unsigned  | 1            | Object type (chest, door, etc.)    | 0-255               |
| 0x01   | x         | 8-bit unsigned  | 1            | X position on map                  | 0-255               |
| 0x02   | y         | 8-bit unsigned  | 1            | Y position on map                  | 0-255               |
| 0x03   | state     | 8-bit unsigned  | 1            | State/flags (opened, locked, etc.) | 0-255               |

- There are 64 objects, each 4 bytes, for a total of 256 bytes.
- The meaning of the `state` byte varies by object type (e.g., for chests, it may indicate if opened; for doors, if locked/unlocked).
- 0xFF in any field typically means "not present" or "inactive."

**Note:** The exact mapping of object types and state flags can be expanded if you need more detail. Refer to OLD source files for specifics.

## NPC Bitfields Structure

The `npc_dead` and `met_player` arrays are each 128 bytes (1024 bits). Each bit represents the state of a specific NPC:

- **npc_dead**: 1 = NPC is dead, 0 = NPC is alive
- **met_player**: 1 = NPC has been met, 0 = NPC has not been met

| Bit Index | NPC | Description |
|-----------|-----|-------------|
| 0         | 0   | NPC 0       |
| 1         | 1   | NPC 1       |
| ...       | ... | ...         |
| 127       | 127 | NPC 127     |

**Location Mapping:**
- NPC bitfields are per-settlement. To determine which location a bitfield refers to, use the location index mapping from DATA.OVL and the wiki. For example, bit 0 in the bitfield for Britain refers to NPC 0 in Britain.
- See DATA.OVL and the wiki for the full location index table.

## Monster Table Structure

The monster table is located at offset 0x6B4 and is 256 bytes (32 monsters × 8 bytes each).

| Offset | Name      | Type/Format    | Size (bytes) | Description/Usage                          |
|--------|-----------|----------------|--------------|--------------------------------------------|
| 0x00   | tile      | 8-bit unsigned | 1            | Tile (first frame of animated group)       |
| 0x01   | anim_tile | 8-bit unsigned | 1            | Tile of current animation frame            |
| 0x02   | x         | 8-bit unsigned | 1            | X-coordinate                               |
| 0x03   | y         | 8-bit unsigned | 1            | Y-coordinate                               |
| 0x04   | z         | 8-bit unsigned | 1            | Z-coordinate (level)                       |
| 0x05   | value1    | 8-bit unsigned | 1            | Value 1 (item number, hull strength, etc.) |
| 0x06   | value2    | 8-bit unsigned | 1            | Value 2 (bitmap, monster type, etc.)       |
| 0x07   | value3    | 8-bit unsigned | 1            | Value 3 (skiffs on board, etc.)            |

**Notes:**
- The monster table contains monsters, inanimate objects (e.g., ships), and the party (slot 0).
- Non-empty entries do not have to be contiguous.
- For frigates, value1 is hull strength, value3 is number of skiffs.
- Value2 may be a bitmap or animation state.

## References
- [Ultima V Internal Formats Wiki](https://wiki.ultimacodex.com/wiki/Ultima_V_internal_formats)
