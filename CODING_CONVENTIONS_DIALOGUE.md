# Dialogue Conventions (UI & Interpreter)

This document describes the conventions used for dialogue handling and modal dialog UI in the codebase. It is intended as a companion to the project's coding conventions and explains the runtime behavior, input flow, script-opcode handling patterns, and UI modal standards used across the project (for example: conversation interpreter, text prompts, and selection modals like the "character selector").

## Overview

- Dialogue execution is handled by a dedicated interpreter that produces output tokens/messages asynchronously (pushed onto a channel) while consuming user input via an input channel.
- Dialogue UI and modal widgets follow consistent visual and interaction patterns:
  - Blue rounded borders with a semi-transparent interior are used for modals.
  - Modals are centered relative to the current game/map screen.
  - Modal vertical size is dynamic and determined by content (e.g., number of list rows).
  - Keyboard navigation (Up/Down/Left/Right/Enter) is the primary control method; Escape may be intentionally disabled for some modals (e.g., forced selection dialogs).

## Conversation lifecycle and channels

- The conversation runs in a separate goroutine and exposes:
  - An output channel for produced ScriptItems (text or control tokens).
  - An input channel where the caller can send text responses or commands.
- The interpreter is non-blocking in the game loop, but the readLine helper will block while waiting for input from the input channel.
- Start/Stop semantics:
  - Start spawns the interpreter goroutine.
  - Stop cancels the context and causes the loop to return and close the output channel.

## Script/opcode handling conventions

- Script items are represented as opcodes (commands) plus optional string or numeric data.
- Opcode handling is conservative: implemented opcodes produce their intended text or state change. Unimplemented opcodes are passed through as a placeholder string so the game still shows something rather than silently failing.
- Where a script modifies game state (e.g., karma, inventory), the opcode handler both updates the state and emits text to the conversation output describing the result when appropriate.
- Many opcodes emit text via helper functions that enqueue strings or formatted strings to the conversation output channel; this keeps text output decoupled from interpreter logic.
- Special opcodes:
  - NewLine-type opcodes produce explicit line breaks in the output stream.
  - Pause opcodes may cause a short delay before continuing output.
  - Runtime toggles (e.g., rune mode) flip interpreter state rather than immediately emitting text.
- Placeholders: For opcodes that are not yet fully implemented, the code emits a short placeholder string so the flow remains visible in the UI. These placeholders are intended to be replaced as functionality is implemented.

## Skip logic, labels, and conditionals

- The interpreter supports skip semantics for branching: skip next, skip after next, skip to label, etc.
- Conditions like "If the avatar is known to the NPC" are handled by setting skip flags and returning early when appropriate.
- Label handling: when a label is defined or a goto occurs, the interpreter manipulates the conversation order and may push new indices into the processing queue.

## Debugging and optional text

- A debug flag controls emission of extra internal/debug text; when enabled, the interpreter emits short debug tokens that help trace behavior (e.g., "KARMA_DEC_ONE").
- Debug tokens are intentionally short and designed to be visible in the conversation output but distinct from normal script text.

## UI modal/dialog conventions

- Visual style:
  - Use the shared Border widget that produces the blue rounded border and interior. This ensures consistent visuals across dialogs.
  - Interior color is typically a translucent dark (semi-transparent black) so the map remains visible under the modal.
  - Dialog width is determined as a percent of the screen; the dialog is centered on the current game-screen center.
- Layout:
  - Vertical size adapts to content (e.g., list length). Calculate the total height from number of rows × row height × spacing + padding.
  - Rows typically have an icon on the left and a name/label on the right; icons should be scaled to fit the row height.
  - Use a subtle selection highlight (e.g., translucent blue overlay) to indicate the current row.
- Selection semantics:
  - Default selection is often index 0 (the Avatar).
  - Keyboard Up/Down (or Left/Right) move the selection; wrap-around behavior is acceptable.
  - Dead/disabled items should be visually distinct — preferred behavior is to gray them out and skip them when navigating (do not allow selection).
  - Enter confirms the selection and invokes a callback; some dialogs intentionally *do not* allow Escape to cancel (e.g., forced selection flows).
- Input handling:
  - The modal reads non-alphanumeric bound keys (Enter, Esc, arrows) similarly to other widgets, using a shared input helper so that repeat behavior and key registration are centralized.
  - For text inputs and command prompts, use the text input widget which supports autocompletion and command matching; its color indicates match status.

## Implementation patterns & helpers

- Use helper functions to:
  - Create percent-based placements and compute pixel rectangles for drawing.
  - Enqueue strings and formatted strings to the conversation/channel rather than directly drawing text in interpreter code.
  - Create row text outputs with a font instance to draw names and to control color per row.
- Keep UI and interpreter decoupled:
  - Interpreter should not directly draw or depend on rendering context; it should only send script items through the output channel and modify game state.
  - UI widgets read game state and renderer resources to display things (like party members), and provide callbacks back to game logic (e.g., onSelect).

## Example: character selection dialog (conventions applied)

- Centered modal with blue border, vertical size computed from party size.
- Each row shows an icon and the character's name; row height and spacing are consistent with existing button/list sizing.
- Avatar (index 0) is the default selection. Up/Down keys move selection and Enter confirms.
- Dead party members are grayed out and are not selectable; navigation skips them.
- There is no Escape path: Enter is required to confirm a valid (alive) selection.

## Best practices

- Keep UI state and game state synchronized: derive modal rows from authoritative game state (party list) at creation time or on open.
- Keep promotion of unimplemented features visible via placeholders to make scripted dialogs testable even before full functionality is implemented.
- Use shared styles and font helpers so the dialogue UI remains consistent across the project.
- Prefer explicit, small helper functions for:
  - computing layout (row positions, icon sizes),
  - drawing icons (scale & optional desaturation when disabled),
  - and emitting conversation text into channels.
