# Special: Bridge Trolls

## Overview

- Trigger: Overworld bridge tile entry while on foot.
- Chance: 1-in-8 (12.5%).
- Flow: Announce trolls, each eligible party member sneaks; on failure, toll or combat; otherwise “Trolls evaded!”.

## Pseudocode

```pseudocode
FUNCTION maybe_trigger_bridge_trolls_on_enter_tile(current_tile, player_state):
    IF NOT is_bridge_tile(current_tile) THEN RETURN
    IF NOT is_on_foot(player_state) THEN RETURN
    bridge_trolls_encounter()
ENDFUNCTION

FUNCTION bridge_trolls_encounter():
    IF random(0, 7) != 0 THEN RETURN
    show_message("\nThou spieth trolls under the bridge!\n\n")
    FOR EACH member IN get_party_members_in_order():
        IF member.status == DEAD OR member.status == ASLEEP THEN CONTINUE
        show_message(member.name + " sneaks across"); show_progress_dots(3); show_message("\n\n")
        IF random(1, 30) > member.dexterity THEN resolve_bridge_troll_caught(member); RETURN
    ENDFOR
    show_message("Trolls evaded!\n")
ENDFUNCTION

FUNCTION resolve_bridge_troll_caught(member):
    toll_gp = 99 - 3 * member.strength
    show_message("Caught!\n\nThe trolls demand a " + format_int(toll_gp) + " gp toll!\n\nDost thou pay? (Y/N)")
    choice = read_yes_no()
    IF choice == YES AND get_party_gold() >= toll_gp THEN adjust_party_gold(-toll_gp); mark_stats_changed(); RETURN
    start_bridge_troll_combat_at_player()
ENDFUNCTION
```

