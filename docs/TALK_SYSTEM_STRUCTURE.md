# Ultima V Conversation System (*.TLK)

This document describes the structure and encoding of the Ultima V conversation system, focusing on the *.TLK files and their relationship to NPCs, DATA.OVL, and game logic. It is based on the canonical format, the wiki (https://wiki.ultimacodex.com/wiki/Ultima_V_internal_formats), and the original fan-made source in OLD (e.g., CHEKTALK.C, COMTALK.C).

---

## Overview

Ultima V's conversation system is driven by *.TLK files, which contain all NPC dialogue scripts for a given map group (CASTLE, DWELLING, KEEP, TOWNE). Each NPC is mapped to a dialogue entry via a dialogue index, and the conversation text is encoded with special codes and compressed words.

---

## *.TLK File Format

Each *.TLK file contains:
- **Header:** Number of entries (NPCs with dialogue)
- **Script Index Table:** For each NPC, a uint16 NPC index and a uint16 offset into the script data
- **Script Data:** Encoded conversation text blocks

### Layout
| Section            | Type/Format          | Description                        |
|--------------------|----------------------|------------------------------------|
| Header             | uint16               | Number of entries                  |
| Script Index Table | [uint16, uint16] × N | NPC index, offset into script data |
| Script Data        | uint8[]              | Encoded conversation text          |

- NPC indices start at 1 and are sorted ascending.
- Offsets point to the start of each NPC's script block.

---

## NPC Mapping

- Each NPC in a map is assigned a dialogue index (see *.NPC files and DATA.OVL).
- The dialogue index maps to the script index in the *.TLK file.
- If the dialogue index is 0, the NPC has no response.

---

## Script Data Encoding

Conversation text is encoded with:
- **Fixed Entries:** Name, Description, Greeting, Job, Bye
- **Keywords:** Encoded as key word blocks
- **Labels:** For branching Q&A
- **Special Codes:** For compressed words, symbols, and control codes

### Fixed Entries
For each NPC, the script block starts with a set of fixed entries:
1. Name
2. Description
3. Greeting
4. Job
5. Bye

Each is a zero-terminated string, possibly containing special codes.

### Keywords and Q&A
- Keywords are encoded as blocks, each with a keyword and an answer text.
- Branching is handled via labels and control codes.

## Keyword Blocks and OR Chains

Keywords are encoded as blocks, each with a keyword and an answer text. Multiple keywords can be chained together using the OR code (135). For example:

| Encoding | Meaning |
|----------|---------|
| [Keyword1][Answer1]\0 | Keyword1 triggers Answer1 |
| [135][0][Keyword2][Answer2]\0 | Keyword1 OR Keyword2 triggers Answer2 |

The OR code (135) is followed by a zero byte and the next keyword. This allows for flexible keyword matching in conversations.

## Label and Branching Mechanics

Labels (codes 145-155) are used for branching Q&A. When a label code is encountered, the conversation can jump to the corresponding answer block. Conditional branches (e.g., code 140 for If/Else knows Avatar's name) allow for dynamic responses based on game state.

**Example:**
- [Label145][Answer for label 1]\0
- [140][Conditional Answer]\0

### Special Codes
| Code Value | Meaning                                             |
|------------|-----------------------------------------------------|
| <129       | Entry to offset table in DATA.OVL (compressed word) |
| 129        | Insert Avatar’s name                                |
| 131        | Conversation pause                                  |
| 135        | OR in the key words                                 |
| 136        | Ask for avatar's name                               |
| 140        | If/Else the NPC knows the Avatar's name             |
| 141        | New line                                            |
| 143        | Key wait                                            |
| 145-155    | Labels 1 to 10                                      |
| >=160,<255 | Subtract 128 for DATA.OVL offset                    |

Other codes may be used for control flow or special symbols.

## Expanded Special Codes Table

| Code Value | Meaning                                             | Example                              |
|------------|-----------------------------------------------------|--------------------------------------|
| <129       | Entry to offset table in DATA.OVL (compressed word) | [128] expands to "Britannia"         |
| 129        | Insert Avatar’s name                                | [129] expands to party leader's name |
| 130        | Reserved/unused                                     |                                      |
| 131        | Conversation pause                                  | [131] pauses dialogue                |
| 132-134    | Reserved/unused                                     |                                      |
| 135        | OR in the key words                                 | [135][0][Keyword2]                   |
| 136        | Ask for avatar's name                               | [136] prompts for name               |
| 137-139    | Reserved/unused                                     |                                      |
| 140        | If/Else the NPC knows the Avatar's name             | [140] branches based on name known   |
| 141        | New line                                            | [141] inserts a newline              |
| 142        | Reserved/unused                                     |                                      |
| 143        | Key wait                                            | [143] waits for keypress             |
| 145-155    | Labels 1 to 10                                      | [145] is label 1                     |
| >=160,<255 | Subtract 128 for DATA.OVL offset                    | [160] expands to DATA.OVL offset 32  |

## Compressed Word Handling Example

Compressed words are encoded as a single byte (<129 or >=160,<255) and expanded at runtime using DATA.OVL. For example:
- [128] in TLK expands to "Britannia" (from DATA.OVL offset table)
- [160] in TLK expands to the word at DATA.OVL offset 32

## NPC Mapping and Dialog Index Details

NPCs are mapped to TLK entries using the dialog_number field in *.NPC files. The mapping process:
- Each NPC has a dialog_number (0 = no response, >0 = TLK entry, 129+ = merchant/special)
- dialog_number maps to the script index in the TLK file
- Special values (e.g., 129 = weapon dealer, 130 = barkeeper, etc.) trigger merchant/innkeeper logic

| dialog_number | Meaning                 |
|---------------|-------------------------|
| 0             | No response             |
| 1-N           | TLK entry index         |
| 129           | Weapon dealer           |
| 130           | Barkeeper               |
| 131           | Horse seller            |
| 132           | Ship seller             |
| 133           | Magic seller            |
| 134           | Guild Master            |
| 135           | Healer                  |
| 136           | Innkeeper               |
| 255           | Guard (harasses player) |

## Edge Cases and Implementation Notes

- NPCs with dialog_number 0 have no conversation.
- Merchants and innkeepers use special dialog logic, not TLK scripts.
- If a TLK entry is missing or malformed, fallback logic should display a default message or skip the NPC.
- Error handling should log and skip invalid entries, ensuring the game does not crash.
- Compatibility with fan remake quirks: some TLK files may have extra or missing entries; robust parsing is recommended.

## Testing and Integration

- Integration tests should load actual data files and simulate conversations with real NPCs.
- Use testdata/britain2_SAVED.GAM for save file state.
- Verify that all fixed entries, keywords, labels, and branching paths are parsed and displayed correctly.
- Test error handling for missing or malformed TLK entries.
- Ensure integration with SAVED.GAM and NPC schedules for correct state-dependent responses.

## Visual Diagrams

Below is a flowchart showing the conversation flow:

```
NPC Mapping (dialog_number) --> TLK Script Index --> TLK Script Block
    |                              |
    v                              v
Merchant/Innkeeper?         Parse Fixed Entries
    |                              |
    v                              v
Special Logic               Parse Keywords/Labels
    |                              |
    v                              v
Display Merchant UI         Branch/Jump/Expand Words
```

## Example Conversation Block
```
[Name]\0[Description]\0[Greeting]\0[Job]\0[Bye]\0
[Keyword1][Answer1]\0[Keyword2][Answer2]\0...
[Label][Branching Answer]\0
```
Special codes may be embedded in any string.

---

## Implementation Notes
- NPCs are mapped to dialogue via *.NPC files and DATA.OVL.
- Text decoding must handle special codes and compressed words.
- Branching and labels allow for complex Q&A and conditional responses.

---

## See Also
- [SAVED_GAM_STRUCTURE.md](./SAVED_GAM_STRUCTURE.md)
- DATA.OVL format and compressed word table

---

## References
- [Ultima V Internal Formats Wiki](https://wiki.ultimacodex.com/wiki/Ultima_V_internal_formats)

---

## NPC Dialog Table

| NPC Name        | Dialog Number | Location           | File         |
|-----------------|--------------|--------------------|--------------|
| Lord British    | 1            | Castle Britannia   | CASTLE.TLK   |
| Geoffrey        | 2            | Castle Britannia   | CASTLE.TLK   |
| Chuckles        | 3            | Castle Britannia   | CASTLE.TLK   |
| Thorne          | 4            | Jhelom             | TOWNE.TLK    |
| Smith the Horse | 5            | Iolo's Hut         | DWELLING.TLK |
| Iolo            | 6            | Iolo's Hut         | DWELLING.TLK |
| Shamino         | 7            | Shamino's Hut      | DWELLING.TLK |
| Dupre           | 8            | Trinsic            | TOWNE.TLK    |
| Katrina         | 9            | New Magincia       | TOWNE.TLK    |
| Jaana           | 10           | Yew                | TOWNE.TLK    |
| Mariah          | 11           | Moonglow           | TOWNE.TLK    |
| Julia           | 12           | Minoc              | TOWNE.TLK    |
| Sentri          | 13           | Serpent's Hold     | KEEP.TLK     |
| Johne           | 14           | Dungeon Deceit     | DWELLING.TLK |
| Gwenno          | 15           | Iolo's Hut         | DWELLING.TLK |
| Toshi           | 16            | Magincia           | TOWNE.TLK    |
| Inamo           | 17            | Magincia           | TOWNE.TLK    |
| Quenton         | 18            | Skara Brae         | TOWNE.TLK    |
| Horance         | 19            | Skara Brae         | TOWNE.TLK    |
| Fiona           | 20            | Skara Brae         | TOWNE.TLK    |
| ...             | ...          | ...                | ...          |

*This table lists all NPCs, their dialog numbers, locations, and TLK files as referenced in the [Ultima V transcript](https://wiki.ultimacodex.com/wiki/Ultima_V_transcript). This is a partial sample. The full table should be expanded to include every NPC from the transcript, with dialog numbers shown without leading zeros.*
