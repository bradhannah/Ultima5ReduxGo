# Ultima V Conversation System (*.TLK)

This document describes the structure and encoding of the Ultima V conversation system, focusing on the *.TLK files and
their relationship to NPCs, DATA.OVL, and game logic. It is based on the canonical format, the
wiki (https://wiki.ultimacodex.com/wiki/Ultima_V_internal_formats), and the original fan-made source in OLD (e.g.,
CHEKTALK.C, COMTALK.C).

---

## Overview

Ultima V's conversation system is driven by *.TLK files, which contain all NPC dialogue scripts for a given map group (
CASTLE, DWELLING, KEEP, TOWNE). Each NPC is mapped to a dialogue entry via a dialogue index, and the conversation text
is encoded with special codes and compressed words.

---

# TLK Script TalkCommand Reference

This table documents the TalkCommand constants used in TLK scripts. Each command is represented by a byte value and
controls dialogue flow, substitutions, prompts, and engine-internal operations.

---

## Sorted TalkCommand Table (by Type)

| Type of Command | Command Name                   | Byte Value | Description (inferred from code/comments)                                                                                                                                                                                                  | Usage Notes | User Action Required                                  | How/Why User Action Is Required                                    | Requires Callback                  | Callback Description          |
|-----------------|--------------------------------|------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------|-------------------------------------------------------|--------------------------------------------------------------------|------------------------------------|-------------------------------|
| Callback Action | JoinParty                      | 0x84       | NPC joins party                                                                                                                                                                                                                            |             | No                                                    |                                                                    | Yes                                | Triggers party join logic     |
| Callback Action | KarmaPlusOne                   | 0x89       | Increase karma by one                                                                                                                                                                                                                      |             | No                                                    |                                                                    | Yes                                | Triggers karma increase logic |
| Callback Action | KarmaMinusOne                  | 0x8A       | Decrease karma by one                                                                                                                                                                                                                      |             | No                                                    |                                                                    | Yes                                | Triggers karma decrease logic |
| Callback Action | CallGuards                     | 0x8B       | Call guards                                                                                                                                                                                                                                |             | No                                                    |                                                                    | Yes                                | Triggers guard call logic     |
| String Output   | ExtortionAmount                | 0xE0       | Extortion amount (engine-internal)                                                                                                                                                                                                         | Sometimes   | No                                                    | Show user how much extortion is                                    | Triggers extortion logic           |                               |
| Callback Action | GoToJail                       | 0xE1       | Go to jail (engine-internal)                                                                                                                                                                                                               | No          |                                                       |                                                                    | Triggers jail logic                | Reposition Avatar into jail   |
| Callback Action | PayGenericExtortion            | 0xE2       | Pay generic extortion (engine-internal)                                                                                                                                                                                                    | Sometimes   | May prompt user for payment or confirmation           |                                                                    | Triggers extortion payment logic   |                               |
| Callback Action | PayHalfGoldExtortion           | 0xE3       | Pay half gold extortion (engine-internal)                                                                                                                                                                                                  | Sometimes   | May prompt user for payment or confirmation           |                                                                    | Triggers half-gold extortion logic |                               |
| Callback Action | MakeAHorse                     | 0xE4       | Make a horse (engine-internal)                                                                                                                                                                                                             | No          |                                                       | Yes                                                                | Triggers horse creation logic      | Make a horse appear           |
| Input Expected  | UserInputNotRecognized         | 0x7E       | User input not recognized                                                                                                                                                                                                                  | Yes         | User must retry input or is notified of invalid input | No                                                                 |                                    |                               |
| Input Expected  | PromptUserForInputUserInterest | 0x7F       | Prompt user for input (user interest)                                                                                                                                                                                                      |             | Yes                                                   | User is prompted to enter a topic of interest                      | No                                 |                               |
| Input Expected  | PromptUserForInputNpcQuestion  | 0x80       | Prompt user for input (NPC question)                                                                                                                                                                                                       |             | Yes                                                   | User is prompted to answer NPC's question                          | No                                 |                               |
| Input Expected  | GoldPrompt                     | 0x85       | Prompt for gold (was Gold)                                                                                                                                                                                                                 |             | Yes                                                   | User is prompted to enter an amount of gold                        | No                                 |                               |
| Input Expected  | AskName                        | 0x88       | Ask for name. Prompts user for name and matches against party members. If match found, marks NPC as met and responds "A pleasure!", otherwise "If you say so..."                                                                        |             | Yes                                                   | User is prompted to provide a name                                 | Yes                                | Name matching and SetMet logic |
| Input Expected  | KeyWait                        | 0x8F       | Wait for key input                                                                                                                                                                                                                         |             | Yes, user expected to press enter to continue         | User must press a key to continue                                  | No                                 |                               |
| String Output   | PlainString                    | 0x00       | Basic printable string                                                                                                                                                                                                                     |             | No                                                    |                                                                    | No                                 |                               |
| String Output   | AvatarsName                    | 0x81       | Substitute Avatar's name                                                                                                                                                                                                                   |             | No                                                    |                                                                    | No                                 |                               |
| String Output   | NewLine                        | 0x8D       | Insert new line                                                                                                                                                                                                                            |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | OrBranch                       | 0x87       | Branch (was Or). Will require a look ahead to determine if there is an OR condition.                                                                                                                                                       |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | IfElseKnowsName                | 0x8C       | If/else branch based on name knowledge. The next script item (+1) will be what happens if they DO know the Avatar (HasMet), the one after that (+2) will be what happens if they do NOT know the Avatar.                                   |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | StartLabelDef                  | 0x90       | Start label definition.    Defines the beginning or end of the label sections.   If next item (+1) is EndScript, then that is end of all labels and conversation script. If next item (+1) is DefineLabel then it is defining a NEW label. |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | Label1                         | 0x91       | Label 1                                                                                                                                                                                                                                    |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | Label2                         | 0x92       | Label 2                                                                                                                                                                                                                                    |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | Label3                         | 0x93       | Label 3                                                                                                                                                                                                                                    |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | Label4                         | 0x94       | Label 4                                                                                                                                                                                                                                    |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | Label5                         | 0x95       | Label 5                                                                                                                                                                                                                                    |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | Label6                         | 0x96       | Label 6                                                                                                                                                                                                                                    |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | Label7                         | 0x97       | Label 7                                                                                                                                                                                                                                    |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | Label8                         | 0x98       | Label 8                                                                                                                                                                                                                                    |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | Label9                         | 0x99       | Label 9                                                                                                                                                                                                                                    |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | Label10                        | 0x9A       | Label 10                                                                                                                                                                                                                                   |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | EndScript                      | 0x9F       | End of script                                                                                                                                                                                                                              |             | No                                                    | Yes, user will be expected to press ENTER to close dialogue window | No                                 |                               |
| Workflow        | StartNewSection                | 0xA2       | Start new section. Defines the beginning of a new Section                                                                                                                                                                                  |             | No                                                    |                                                                    | No                                 |                               |
| Workflow        | GotoLabel                      | 0xFD       | Go to label. Move pointer to new position immediately.                                                                                                                                                                                     | No          |                                                       | No                                                                 |                                    |                               |
| Workflow        | DefineLabel                    | 0xFE       | Define label.  Defines the beginning of a new label. Next script item (+1) will be a Label Id (ie. Label1).                                                                                                                                | No          |                                                       | No                                                                 |                                    |                               |
| Workflow        | DoNothingSection               | 0xFF       | No operation                                                                                                                                                                                                                               | No          |                                                       | No                                                                 |                                    |                               |
| Other           | EndConversation                | 0x82       | End the current conversation                                                                                                                                                                                                               |             | No                                                    |                                                                    | No                                 |                               |
| Other           | Pause                          | 0x83       | Pause dialogue                                                                                                                                                                                                                             |             | Sort of - the UI will pause for a period of time      | May require user to press a key to continue                        | No                                 |                               |
| Other           | Change                         | 0x86       | Change (context-specific)                                                                                                                                                                                                                  |             | Sometimes                                             | May require user to confirm or act depending on context            | No                                 |                               |
| Other           | Rune                           | 0x8E       | Rune (context-specific)                                                                                                                                                                                                                    |             | No                                                    |                                                                    | No                                 |                               |

Here are a few snippets from a C# project.

```C#
/// <summary>
///     Does this line represent the end of all Labels in the NPC talk script (end of script)
/// </summary>
/// <returns></returns>
public bool IsEndOfLabelSection =>
    GetScriptItem(0).Command == TalkCommand.StartLabelDefinition &&
    GetScriptItem(1).Command == TalkCommand.EndScript;

/// <summary>
///     Does this line represent a new label definition
/// </summary>
/// <returns></returns>
public bool IsLabelDefinition() =>
    GetScriptItem(0).Command == TalkCommand.StartLabelDefinition &&
    GetScriptItem(1).Command == TalkCommand.DefineLabel;

```

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

Keywords are encoded as blocks, each with a keyword and an answer text. Multiple keywords can be chained together using
the OR code (135). For example:

| Encoding                      | Meaning                               |
|-------------------------------|---------------------------------------|
| [Keyword1][Answer1]\0         | Keyword1 triggers Answer1             |
| [135][0][Keyword2][Answer2]\0 | Keyword1 OR Keyword2 triggers Answer2 |

The OR code (135) is followed by a zero byte and the next keyword. This allows for flexible keyword matching in
conversations.

## Label and Branching Mechanics

Labels (codes 145-155) are used for branching Q&A. When a label code is encountered, the conversation can jump to the
corresponding answer block. Conditional branches (e.g., code 140 for If/Else knows Avatar's name) allow for dynamic
responses based on game state.

**Example:**

- [Label145][Answer for label 1]\0
- [140][Conditional Answer]\0

### Special Codes

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

## TLK File Conversation Processing Algorithm

This section describes the general algorithm for reading and processing TLK files during NPC conversations. It abstracts the pointer-driven logic used in the original TALKNPC.C implementation, suitable for re-implementation in other languages or systems.

### 1. Loading TLK Data
- Load the TLK file data for the current NPC into a buffer.
- Use the script index table to locate the start of the NPC's script block.

### 2. Pointer-Based Script Traversal
- Initialize a pointer to the start of the NPC's script block.
- Traverse the TLK script data using the pointer, reading one byte/command at a time.

### 3. Searching for Dialogue Blocks
- Use search routines to locate specific markers, string numbers, or command codes within the TLK data.
- Advance the pointer to the desired block or response using delimiter bytes or command codes.

### 4. Outputting Dialogue
- Output printable strings by reading bytes from the buffer until a null or delimiter is found.
- Handle text formatting and line wrapping as needed.

### 5. Handling User Input and Branching
- When prompted, wait for user input and match it against known keywords or triggers.
- Use conditional branching to advance the pointer to the appropriate response block based on input or script logic.

### 6. Executing Scripted Actions
- When a TLK command requiring an action is encountered (e.g., join party, give item, adjust karma), trigger the corresponding game logic.
- Actions are executed by interpreting specific command bytes in the TLK data.

### 7. Looping and Termination
- Continue processing commands and advancing the pointer until an end-of-script marker is reached or the conversation is terminated.

#### Abstracted Algorithm (for LLM Re-implementation)

1. Load TLK data for the current NPC into a buffer.
2. Initialize a pointer to the start of the NPC’s script block.
3. While not at end-of-script:
    - Read the next byte/command.
    - If it’s a printable string, output it.
    - If it’s a prompt, wait for user input and match against keywords.
    - If it’s a branch or label, advance the pointer to the correct block.
    - If it’s an action command, trigger the corresponding game logic.
    - Advance the pointer as needed.
4. Repeat until the conversation ends.

---
