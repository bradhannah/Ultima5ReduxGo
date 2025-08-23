# Potions

Ultima V has eight potion colors. Potions apply immediate effects to a selected party member, with some effects only working in combat or outside dungeons. There is a small randomness factor that may alter the effect.

## Use Flow

```pseudocode
FUNCTION use_potion(color_index):
    // color_index: 0..7 → Blue, Yellow, Red, Green, Orange, Purple, Black, White
    potions[color_index] -= 1
    show_message("Potion\n")

    // Choose target
    IF in_combat() THEN plr = combatq[turn].number ELSE plr = select_target_party_member()
    IF plr < 0 THEN RETURN FALSE

    // Random variance: 1/16 chance force Orange (Sleep), 1/16 chance random color
    r = random(0, 15)
    IF r == 0 THEN eff = ORANGE
    ELSE IF r == 1 THEN eff = random(0, 7)
    ELSE eff = color_index

    SWITCH eff:
        CASE BLUE:
            IF player[plr].status == SLEEP THEN
                player[plr].status = GOOD
                IF in_combat() AND is_current_actor(plr) AND current_cond_is(PC|SLEEP) THEN wakeup(turn)
                show_message("Awake!\n"); update_stats()
            ELSE show_message("No effect!\n"); RETURN FALSE
        CASE YELLOW:
            IF heal(plr) THEN show_message("Healed!\n"); update_stats() ELSE RETURN FALSE
        CASE RED:
            IF player[plr].status == POISONED THEN player[plr].status = GOOD; show_message("Poison cured!\n"); update_stats() ELSE RETURN FALSE
        CASE GREEN:
            IF player[plr].status == GOOD THEN player[plr].status = POISONED; show_message("POISONED!\n"); update_stats() ELSE RETURN FALSE
        CASE ORANGE:
            IF player[plr].status == GOOD THEN
                IF in_combat() THEN put_actor_to_sleep(plr) ELSE player[plr].status = SLEEP
                show_message("Slept!\n"); update_stats()
            ELSE RETURN FALSE
        CASE PURPLE: // Polymorph to rat (combat only)
            IF in_combat() THEN show_message("Poof!\n"); set_actor_shape(plr, RATS) ELSE show_message("\nNo noticeable effect now!\n")
        CASE BLACK: // Invisibility (combat only)
            IF in_combat() THEN show_message("Invisible!\n"); set_invisible(plr); set_actor_shape(plr, INVIS_PC) ELSE show_message("\nNo noticeable effect now!\n")
        CASE WHITE: // X-Ray (not in dungeons)
            IF NOT is_dungeon() THEN x_ray_vision() ELSE show_message("\nNo noticeable effect now!\n")
    ENDSWITCH
    RETURN TRUE
ENDFUNCTION
```

Notes:

- Selection: In combat, the current actor is the potion user; outside, you pick a target party member.
- Randomness: 1/16 chance the potion acts as Orange (Sleep), and 1/16 chance it acts as a random color.
- UI strings: Matches legacy behavior, including “No noticeable effect now!” for context-mismatched effects.

## Effects Matrix

| Color  | Effect                     | Context constraint                | UI text                       |
|--------|----------------------------|-----------------------------------|-------------------------------|
| Blue   | Cure Sleep (Awaken)        | Target must be Asleep             | “Awake!”                      |
| Yellow | Heal (1..30 HP, up to max) | —                                 | “Healed!”                     |
| Red    | Cure Poison                | Target must be Poisoned           | “Poison cured!”               |
| Green  | Apply Poison               | Target must be Good               | “POISONED!”                   |
| Orange | Apply Sleep                | Target must be Good               | “Slept!”                      |
| Purple | Polymorph to Rat           | Combat only                       | “Poof!”                       |
| Black  | Invisibility (PC)          | Combat only                       | “Invisible!”                  |
| White  | X-Ray Vision               | Not in dungeons (towns/overworld) | “(X-ray effect)” or no-notice |

Tips:

- Black potion sets invisible flags and changes the actor’s shape to the invisible PC tile; AI targeting rules apply.
- Purple potion swaps the actor’s object tile/shape to the RATS tile; duration and gameplay impact are engine-defined.
- White potion calls the same helper as the Sanct Lor spell; it’s suppressed in dungeons.

