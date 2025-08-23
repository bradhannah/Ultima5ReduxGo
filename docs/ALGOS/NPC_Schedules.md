# NPC Schedules

High-level overview of how NPCs move and behave on small maps (towns, interiors) based on schedule data and simple wander logic. This summarizes the common patterns.

## Concepts

- Anchor Position: Each NPC typically has a schedule with an anchor coordinate.
- Wander Within N: Many behaviors allow wandering within N tiles of the anchor.
- Behaviors: Enum-like schedule modes define how an NPC moves (e.g., horse wandering).
- Passability: Wander steps only commit if the target tile is passable for NPCs.

## Wander Within N

```pseudocode
FUNCTION wander_one_tile_within_n(npc, anchor_pos, N):
    FOR tries IN 1..4:
        dir = random_direction()
        newPos = position_in_direction(npc.pos, dir)
        IF abs(newPos.x - anchor_pos.x) <= N AND abs(newPos.y - anchor_pos.y) <= N AND is_npc_passable(newPos) THEN
            npc.pos = newPos
            RETURN TRUE
        ENDIF
    ENDFOR
    RETURN FALSE
ENDFUNCTION
```

## Data Model (From Legacy)

Each NPC has a schedule describing up to 4 timed destinations and an action code for each time slot. At runtime, a per-NPC current state (`stat[npc]`) holds the active action, destination index, and current map level.

```pseudocode
STRUCT Schedule {
    times[4]: int        // hours at which transitions occur (0..23)
    x[4], y[4], z[4]: int// destination coords (z=floor)
    action[?]: int       // per-time-slot behavior code (see below)
}

STRUCT CurrentNPCState {
    x, y, z: int       // current position (tile coords + level)
    objnum: int        // object index backing this NPC
    act: enum { INACT, MOVE, LEAVU, LEAVD, ARIVA, ARIVB, POP }
    misc1: int         // cached last destination index
}
```

Action codes (common):

- INACT: idle at destination; may wander locally by behavior.
- MOVE: pathfind from source to destination on same level.
- LEAVU/LEAVD: leave current level (via up/down ladder/stairs) towards destination on different level.
- ARIVA/ARIVB: arrive onto current level from above/below and proceed to destination.
- POP: source and destination are on different levels and neither is current; no on-screen movement needed (handled offscreen).

## Time-of-Day Transition Driver

On each hour tick, evaluate schedule transitions for all active NPCs. The legacy `move(ndex, time)` routine encapsulates this.

```pseudocode
FUNCTION npc_hourly_driver(current_hour):
    FOR each npc IN active_npcs_on_map():
        dest_slot = -1
        S = schedule[npc]
        // Find matching time slot
        FOR i in 0..3:
            IF current_hour == S.times[i] THEN dest_slot = timetodest(npc, current_hour); BREAK
        IF dest_slot == -1 THEN CONTINUE

        C = state[npc]
        // If already at destination, mark inactive
        IF S.x[dest_slot]==C.x AND S.y[dest_slot]==C.y AND S.z[dest_slot]==C.z THEN
            C.act = INACT; CONTINUE
        ENDIF

        // Decide transition action based on floors (z) and current visible level
        IF C.z == current_level():
            IF S.z[dest_slot] == current_level(): C.act = MOVE
            ELSE C.act = (S.z[dest_slot] < current_level() ? LEAVU : LEAVD)
        ELSE
            C.act = (C.z < current_level() ? ARIVA : ARIVB)
        ENDIF

        C.misc1 = dest_slot // cache last dest
    ENDFOR
ENDFUNCTION
```

Notes:

- `timetodest(npc, hour)` maps the hour to the target destination index for the NPC (handles wrap‑around and ordering).
- If both source and destination are off the current level, NPCs are handled with POP (offscreen) to avoid on‑screen animation.

## Executing Transitions (Per Tick)

After marking actions at the hour change, advance NPCs each tick:

```pseudocode
FUNCTION npc_tick_update():
    FOR each npc IN active_npcs_on_map():
        C = state[npc]; S = schedule[npc]
        SWITCH C.act:
            CASE MOVE:
                pathfind_or_randstep_towards(npc, S.x[C.misc1], S.y[C.misc1])
                IF (C.x, C.y) == (S.x[C.misc1], S.y[C.misc1]) THEN C.act = INACT
            CASE LEAVU, LEAVD:
                // Move to the appropriate ladder/stairs and change level when on it
                IF laddercheck(npc, C.misc1) THEN change_level(npc, dir=(C.act==LEAVU?Up:Down)); C.act = (C.act==LEAVU?ARIVB:ARIVA)
                ELSE step_towards_ladder(npc, target_level=S.z[C.misc1])
            CASE ARIVA, ARIVB:
                // NPC has just arrived from offscreen; place on ladder/stairs and mark MOVE
                tportnpc(npc, S.x[C.misc1], S.y[C.misc1], current_level()) // or spawn at ladder and then MOVE
                C.act = MOVE
            CASE POP:
                // No on-screen action; handled implicitly
            CASE INACT:
                // Idle/wander based on behavior
                apply_schedule_behavior(npc, npc.behavior, anchor_pos_from(S, C.misc1))
        ENDSWITCH
    ENDFOR
ENDFUNCTION
```

## Day/Night and Special Transitions

- Drawbridges/portcullis toggle at night/day (see Towns → Drawbridges): Avoid flipping under the player.
- Beds (inbed flow) will eject sleepers on schedule boundaries: at hour 20 or 5 the engine may force “Thrown out of bed!” and resume NPC schedules.
- Shops: Opening/closing may depend on hour; out‑of‑scope here but interacts with NPCs (shopkeepers go to/from counters/back rooms).

## Pathfinding and Random Stepping

Pathfinding (`pathfind`, `getpath`) is used when a path is available; otherwise `randstep` or targeted `step_towards` is used. For attraction/flee patterns, NPCs evaluate 4 directional options and pick a move that improves or maintains distance with some randomness.

## Schedule Data Authoring Tips

- Use four evenly spaced time anchors (e.g., 6, 12, 18, 0) for simple day → noon → evening → night cycles.
- Align `z` (level/floor) with map floors; ensure ladder placements are reachable from both source and dest using `laddercheck`.
- For NPCs that should not appear on certain floors at night (e.g., shopkeepers), set POP transitions to avoid on‑screen animation between non‑visible floors.

## Testing Guidance (Schedules)

- Use a fixed time sequence (e.g., simulate hour=5→6→…→20→5) to verify actions change as expected.
- Place NPCs at edge cases: at the ladder tile at transition time; confirm LEAVx → ARIVx transitions behave correctly.
- Verify offscreen transitions (POP) do not produce errant on‑screen moves.
- Confirm per‑behavior idle/wander persists during INACT and resumes MOVE at next scheduled time.

## Schedule Actions Reference

| Action | Trigger (at hour)                            | Execution (tick)                                        | Completion Condition                        | Notes                                               |
|--------|----------------------------------------------|---------------------------------------------------------|---------------------------------------------|-----------------------------------------------------|
| INACT  | Already at destination at transition         | Idle; apply behavior (wander/patrol/static)             | Next hour transition                        | Default resting state                              |
| MOVE   | Source z == current level; dest z == current | Pathfind/step toward (dest.x, dest.y)                   | Reaches (dest.x, dest.y) → INACT            | Uses `pathfind` or `randstep` fallback             |
| LEAVU  | Source z == current; dest z < current        | Step toward up ladder; if on ladder then change level   | On ladder → change to ARIVB; else continue  | Ladder detection via `laddercheck`                 |
| LEAVD  | Source z == current; dest z > current        | Step toward down ladder; if on ladder then change level | On ladder → change to ARIVA; else continue  |                                                   |
| ARIVA  | Source z > current; dest z == current        | Teleport/arrive onto current level; mark MOVE           | After spawn → MOVE                           | Spawn at ladder or at dest; then walk to dest      |
| ARIVB  | Source z < current; dest z == current        | Teleport/arrive onto current level; mark MOVE           | After spawn → MOVE                           |                                                   |
| POP    | Source z != current and dest z != current    | No on‑screen action (offscreen transition)              | Handled offscreen                            | Use when both ends are invisible to the player     |

## Example: Shopkeeper Schedule (Template)

| Hour | Dest (x,y,z) | Action at Transition     | Behavior at INACT           | Notes                              |
|------|--------------|--------------------------|-----------------------------|------------------------------------|
| 06   | (12, 08, 0)  | MOVE (same level)        | Static (behind counter)      | Open shop                          |
| 12   | (10, 14, 0)  | MOVE (same level)        | INACT (break)                | Lunch break                        |
| 14   | (12, 08, 0)  | MOVE (same level)        | Static (behind counter)      | Resume open                        |
| 20   | (04, 22, 1)  | LEAVD (to downstairs)    | INACT (home)                 | Close shop; retire for the night   |

Authoring:

- Use POP for transitions where both source and destination are off the current level to avoid on‑screen animation during time skips.
- Patrol routes can be represented by setting multiple dest points on the same level and applying Patrol behavior during INACT.

## Example: Guard Night Patrol (Template)

| Hour | Dest (x,y,z) | Action at Transition       | Behavior at INACT   | Notes                                  |
|------|--------------|----------------------------|---------------------|----------------------------------------|
| 18   | (20, 06, 0)  | ARIVA (from barracks, up)  | Patrol (route A)    | Begin patrol at dusk                   |
| 22   | (28, 10, 0)  | MOVE (same level)          | Patrol (route B)    | Shift to second sector                 |
| 02   | (12, 18, 0)  | MOVE (same level)          | Patrol (route C)    | Rotate routes overnight                 |
| 06   | (04, 04, 1)  | LEAVD (to barracks below)  | INACT (sleep/home)  | Return to barracks at dawn             |

Authoring:

- Patrol routes can be specified via per-NPC route lists (waypoints) and assigned to Patrol behavior.
- Use LEAVD/ARIVA at dawn/dusk to enter/exit barracks across levels. If barracks isn’t on the current level, use POP to avoid on-screen animation.
- When guard alarms are raised (see Towns → Special Guard Behavior), temporarily override Patrol with setattack until alarm clears, then resume schedule.
## Behavior Dispatch

```pseudocode
FUNCTION apply_schedule_behavior(npc, behavior, anchor_pos):
    SWITCH behavior:
        CASE HorseWander:
            // Keep horse near stable/anchor
            IF NOT within_n(npc.pos, anchor_pos, N=2) THEN
                step_towards(npc, anchor_pos)
            ELSE
                wander_one_tile_within_n(npc, anchor_pos, N=2)
            ENDIF
        CASE WanderArea:
            wander_one_tile_within_n(npc, anchor_pos, N=3)
        CASE Patrol:
            follow_patrol_route(npc)
        CASE Static:
            // No movement
        DEFAULT:
            // Safe fallback: small wander
            wander_one_tile_within_n(npc, anchor_pos, N=1)
    ENDSWITCH
ENDFUNCTION
```

## Notes

- When an NPC moves, refresh occupancy and visuals for consistency.
- Schedule transitions are driven externally (e.g., time-of-day) and can change behaviors.
- If an NPC falls out of its permitted radius, bias steps back towards the anchor before wandering again.
