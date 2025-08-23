# RNG and Saves

## Random Number Generation

```pseudocode
// Core random: inclusive bounds
FUNCTION random(min, max) -> int

// 30-sided die: 1..30
FUNCTION rolld30() -> int

// 0..n-1 value
FUNCTION getrandom(n) -> int
```

Usage notes:

- Always call RNG through the central PRNG to preserve determinism and seedability.
- When documenting odds, include both fraction and comment near the RNG call.

## Intelligence Save (`saveint`)

```pseudocode
FUNCTION does_character_resist_spell(attacker, defender):
    attacker_intelligence = attacker.intelligence
    defender_intelligence = defender.intelligence
    target_number = (30 + defender_intelligence - attacker_intelligence) / 2
    dice_roll = rolld30()
    RETURN dice_roll < target_number
ENDFUNCTION
```

Notes:

- No intelligence save applies versus certain direct-damage spells (e.g., Flam Por, Vas Flam, and higher-tier damage effects) per legacy logic; model these as auto‑fail saves or bypass the save entirely for those spell IDs.
