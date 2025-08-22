# Moongates

## On-Gate Travel and Phases

```pseudocode
FUNCTION on_moongate_if_present():
    IF tile_at_player() != MOONGATE THEN RETURN FALSE
    wait_update(1); play_moongate_pulse(); old_form = player_form; set_player_form(FIZZ); fizz_tile_effect(MOONGATE, 5, 5); set_player_form(ON_FOOT); wait_update(1)
    FOR gate_stage = 15 DOWNTO 1: show_moongate_animation(gate_stage, 5, 5); tick(2)
    set_tile_at_player(GRASS); force_los_update(); draw_tile(GRASS, 5, 5)
    IF hour == 0 AND minute < 10 THEN enter_shrine_of_spirituality(); set_player_form(old_form); RETURN TRUE
    IF hour < 12 THEN gate_travel(TRAMMEL_PHASE) ELSE gate_travel(FELUCCA_PHASE)
    set_player_form(old_form); RETURN TRUE
ENDFUNCTION

FUNCTION gate_travel(phase):
    IF stones_map_for_phase(phase) == 0xFF THEN RETURN FALSE
    IF is_outdoors() THEN save_overworld_objects()
    oldmap = current_map_id(); onmap = stones_map_for_phase(phase); player_x = stones_x_for_phase(phase); player_y = stones_y_for_phase(phase); level = stones_level_for_phase(phase)
    IF is_town_map(onmap) AND is_town_map(oldmap) THEN init_town(reuse=True)
    ELSE IF is_overworld(onmap) AND is_overworld(oldmap) THEN load_overworld_objects(); init_overworld()
    RETURN TRUE
ENDFUNCTION
```


## Phase Timing (Summary)

| Time Window      | Destination Phase | Notes                                  |
|------------------|-------------------|----------------------------------------|
| 00:00–00:09      | Shrine window     | Enters Shrine of Spirituality directly |
| 00:10–11:59      | Trammel           | Uses Trammel-linked gate travel        |
| 12:00–23:59      | Felucca           | Uses Felucca-linked gate travel        |

## Stones Mapping Template

Moongate destinations are defined via four arrays (per phase index): map, X, Y, level. Populate these via data.

```pseudocode
STRUCT StonesPhase {
    Map: byte
    X: byte
    Y: byte
    Level: byte
}

LIST<StonesPhase> TrammelStones // 8 entries
LIST<StonesPhase> FeluccaStones // 8 entries
```
