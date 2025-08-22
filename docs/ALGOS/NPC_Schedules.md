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

