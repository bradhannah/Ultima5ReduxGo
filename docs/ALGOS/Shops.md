# Shops and Inns

This document captures the high-level flows for merchants, healers, reagent sellers, horse sellers, and inns.

## Town Multipliers (Per-Shop-Type)

Each town applies a multiplier by shop type. Final price = `round(BasePrice × TownMultiplier)`.

### Town Multipliers — Matrix

| Town            | Arms | Reagents | Healer | Inn  | Horses | Ships | Guild |
|-----------------|:----:|:--------:|:------:|:----:|:------:|:-----:|:-----:|
| Britain         |  —   |    —     |   —    |  —   |   —    |   —   |   —   |
| Moonglow        |  —   |    —     |   —    |  —   |   —    |   —   |   —   |
| Jhelom          |  —   |    —     |   —    |  —   |   —    |   —   |   —   |
| Yew             |  —   |    —     |   —    |  —   |   —    |   —   |   —   |
| Minoc           |  —   |    —     |   —    |  —   |   —    |   —   |   —   |
| Trinsic         |  —   |    —     |   —    |  —   |   —    |   —   |   —   |
| Skara Brae      |  —   |    —     |   —    |  —   |   —    |   —   |   —   |
| New Magincia    |  —   |    —     |   —    |  —   |   —    |   —   |   —   |
| Cove            |  —   |    —     |   —    |  —   |   —    |   —   |   —   |
| Buccaneer’s Den |  —   |    —     |   —    |  —   |   —    |   —   |   —   |
| Paws            |  —   |    —     |   —    |  —   |   —    |   —   |   —   |

Notes:

- Fill with canonical values (defaults 1.00 where unspecified). “Guild” covers general goods (keys, torches, gems, etc.). Not all towns have all shop types.

## General Notes

- Gems: Vendors may sell gems. Using `View` consumes one gem to render a tactical map; see Commands → View.
- Determinism: All price rolls and stock checks must use the central PRNG; prices should be data-driven where possible.

## Healing Shop

```pseudocode
FUNCTION visit_healing_shop(active_player, shop_rank):
    // shop_rank influences price schedule
    price = compute_heal_price(shop_rank)
    show_message("I can heal thee ")
    IF confirm_pay(price) == NO THEN RETURN
    IF gold < price THEN show_message("Not enough gold!\n"); RETURN
    gold -= price; mark_stats_changed()
    // Apply effect: cure poison if poisoned, otherwise heal HP
    IF player[active_player].status == POISONED THEN
        player[active_player].status = GOOD
        show_message("Poison cured!\n"); play_heal_sfx(); mark_stats_changed()
    ELSE
        healed = heal(active_player) // adds 1..30 HP up to max
        IF healed THEN play_heal_sfx()
    ENDIF
ENDFUNCTION

FUNCTION heal(plr):
    IF player[plr].status == DEAD THEN RETURN FALSE
    incint(&player[plr].hit_points, rolld30(), player[plr].hit_max)
    mark_stats_changed(); RETURN TRUE
ENDFUNCTION
```

## Reagent Shop

```pseudocode
FUNCTION visit_reagent_shop():
    // Reagents: {Sulfurous Ash, Ginseng, Garlic, Spider Silk, Blood Moss, Black Pearl, Nightshade, Mandrake}
    selection = choose_reagent_and_qty()
    total_price = compute_reagent_price(selection)
    IF gold < total_price THEN show_message("Not enough gold!\n"); RETURN
    gold -= total_price
    FOR EACH (reagent, qty) IN selection: inventory.reagents[reagent] += qty
    // Optional: track shop restock timing via reagent_days_left if needed
    mark_stats_changed()
ENDFUNCTION

### Reagents — Items and Base Prices

| Reagent        | Base Price | Notes                         |
|----------------|------------|-------------------------------|
| Sulfurous Ash  | —          |                               |
| Ginseng        | —          |                               |
| Garlic         | —          |                               |
| Spider Silk    | —          |                               |
| Blood Moss     | —          |                               |
| Black Pearl    | —          |                               |
| Nightshade     | —          | Rare; may be limited          |
| Mandrake Root  | —          | Rare; may be limited          |

Notes:

- Final price: `round(Base * TownMuls[reagents])`.
- Stock/availability may depend on day/time or town; reflect per-location overrides if present.
```

## Arms/Armor Shop (simplified)

```pseudocode
FUNCTION visit_arms_shop():
    item = choose_weapon_or_armor()
    price = price_for(item)
    IF gold < price THEN show_message("Not enough gold!\n"); RETURN
    gold -= price; add_item_to_inventory(item); mark_stats_changed()
ENDFUNCTION

### Arms & Armor — Items and Base Prices

Helms/Shields/Armor:

| Item          | Base | Notes           |
|---------------|------|-----------------|
| Leather Helm  | —    |                 |
| Chain Coif    | —    |                 |
| Iron Helm     | —    |                 |
| Spiked Helm   | —    |                 |
| Small Shield  | —    |                 |
| Large Shield  | —    |                 |
| Spiked Shield | —    |                 |
| Magic Shield  | —    |                 |
| Jeweled Shield| —    |                 |
| Cloth Armor   | —    |                 |
| Leather Armor | —    |                 |
| Ring Mail     | —    |                 |
| Scale Armor   | —    |                 |
| Chain Armor   | —    |                 |
| Plate Armor   | —    |                 |
| Mystic Armor  | —    | special         |

Weapons/Ammo:

| Item           | Base | Notes                |
|----------------|------|----------------------|
| Dagger         | —    |                      |
| Sling          | —    | two‑handed           |
| Club           | —    |                      |
| Flame Oil      | —    | thrown consumable    |
| Main‑Gauche    | —    |                      |
| Spear          | —    |                      |
| Throwing Axe   | —    |                      |
| Short Sword    | —    |                      |
| Mace           | —    |                      |
| Morning Star   | —    |                      |
| Bow            | —    | two‑handed           |
| Arrows         | —    | ammo                 |
| Crossbow       | —    | two‑handed           |
| Quarrels       | —    | ammo                 |
| Long Sword     | —    |                      |
| 2H Hammer      | —    | two‑handed           |
| 2H Axe         | —    | two‑handed           |
| 2H Sword       | —    | two‑handed           |
| Halberd        | —    | two‑handed           |
| Chaos Sword    | —    | special              |
| Magic Bow      | —    | special              |
| Silver Sword   | —    | special              |
| Magic Axe      | —    | special              |
| Glass Sword    | —    | one‑use              |
| Jeweled Sword  | —    | special              |
| Mystic Sword   | —    | special              |

Rings/Amulets:

| Item               | Base | Notes                   |
|--------------------|------|-------------------------|
| Invisibility Ring  | —    | effect‑tagged; see Ztats |
| Protection Ring    | —    | effect‑tagged            |
| Regeneration Ring  | —    | effect‑tagged            |
| Amulet of Turning  | —    | special                  |
| Spiritual Collar   | —    | special                  |
| Ankh               | —    | special                  |

Notes:

- Final price: `round(Base * TownMuls[arms])`.
- Per‑shop inventory varies; use per‑shop item whitelists to reflect availability.
```

## Horse Seller

```pseudocode
FUNCTION visit_horse_seller():
    IF confirm_buy_horse() == NO THEN RETURN
    IF already_has_horse() THEN show_message("Stable is full!\n"); RETURN
    spawn_horse_nearby() // MakeAHorse callback
    show_message("A fine steed is thine!\n")
ENDFUNCTION

### Horses — Base Price

| Item  | Base | Notes              |
|-------|------|--------------------|
| Horse | —    | Multiplied by town |

## Shipwright — Items and Services

| Item/Service  | Base | Notes                           |
|---------------|------|---------------------------------|
| Skiff         |  —   | Portable boat                   |
| Ship          |  —   | Full‑size vessel                |
| Repair Hull   |  —   | Per‑point or flat; by location  |

Notes:

- Availability varies by town; coastal locations only.
```

## Innkeeper

```pseudocode
FUNCTION visit_inn():
    months = ask_how_many_months(1..N)
    cost = months * inn_monthly_rate()
    IF gold < cost THEN show_message("Not enough gold!\n"); RETURN
    gold -= cost; mark_stats_changed()
    // Park selected characters at the inn; update saved-game metadata
    FOR EACH member IN party_members_to_stay():
        member.party_status = AtTheInn
        member.months_inn += months
        // inn_party field reflects where the member is staying (by location index)
    ENDFOR
    show_message("Rest well!\n")
ENDFUNCTION

### Inns — Monthly Rates

| Town           | Monthly Rate |
|----------------|--------------|
| Britain        | —            |
| Moonglow       | —            |
| Jhelom         | —            |
| Yew            | —            |
| Minoc          | —            |
| Trinsic        | —            |
| Skara Brae     | —            |
| New Magincia   | —            |
| Cove           | —            |
| Bucs Den       | —            |
| Paws           | —            |

## Guild Shop — General Goods

| Item    | Base | Notes                         |
|---------|------|-------------------------------|
| Key     |  —   | Standard door key             |
| Torch   |  —   | Wall/hand torch               |
| Gem     |  —   | Consumed by View              |

Notes:

- Final price: `round(Base * TownMuls[guild])`. Inventory may include other sundries per location.
```


## Pricing Schema (Examples)

Pricing can be driven by base prices and per-location multipliers to allow towns to vary costs without code changes.

```pseudocode
STRUCT PriceEntry { ItemID: string, Base: int }
STRUCT LocationPriceMul { LocationID: string, Multiplier: float }

LIST<PriceEntry> PriceTable = [
    { ItemID: HEAL_VISIT, Base: 50 },
    { ItemID: REAGENT_SULFUROUS_ASH, Base: 4 },
    { ItemID: REAGENT_GINSENG, Base: 6 },
    // ... etc.
]

LIST<LocationPriceMul> LocationMuls = [
    { LocationID: "Britain", Multiplier: 1.00 },
    { LocationID: "Moonglow", Multiplier: 1.10 },
    { LocationID: "Jhelom", Multiplier: 0.90 },
]

FUNCTION compute_price(item_id, location_id):
    base = PriceTable[item_id].Base
    mul = LocationMuls.get(location_id, default=1.00)
    return round_to_int(base * mul)
ENDFUNCTION

// Example: healing shop price choices by shop rank
LIST<int> HealPriceByRank = [35, 40, 45, 50, 55, 60, 65, 70]

FUNCTION compute_heal_price(shop_rank):
    return HealPriceByRank[clamp(shop_rank, 0, len(HealPriceByRank)-1)]
ENDFUNCTION

// Reagent total
FUNCTION compute_reagent_price(selection):
    total = 0
    FOR EACH (reagent, qty) IN selection:
        total += qty * compute_price(reagent, current_town())
    RETURN total
ENDFUNCTION
```

## Hours of Operation (Template)

Some shops close at night or on specific days. Capture canonical hours here; apply by shop type and town.

| Town            | Arms | Reagents | Healer | Inn | Horses | Shipwright | Guild |
|-----------------|:----:|:--------:|:------:|:---:|:------:|:----------:|:-----:|
| Britain         |  —   |    —     |   —    |  —  |   —    |     —      |   —   |
| Moonglow        |  —   |    —     |   —    |  —  |   —    |     —      |   —   |
| Jhelom          |  —   |    —     |   —    |  —  |   —    |     —      |   —   |
| Yew             |  —   |    —     |   —    |  —  |   —    |     —      |   —   |
| Minoc           |  —   |    —     |   —    |  —  |   —    |     —      |   —   |
| Trinsic         |  —   |    —     |   —    |  —  |   —    |     —      |   —   |
| Skara Brae      |  —   |    —     |   —    |  —  |   —    |     —      |   —   |
| New Magincia    |  —   |    —     |   —    |  —  |   —    |     —      |   —   |
| Cove            |  —   |    —     |   —    |  —  |   —    |     —      |   —   |
| Buccaneer’s Den |  —   |    —     |   —    |  —  |   —    |     —      |   —   |
| Paws            |  —   |    —     |   —    |  —  |   —    |     —      |   —   |

Legend: list as 24h ranges (e.g., 8–20), or Closed.

## Per‑Town Inventories (Templates)

Define available items per town and shop type. Use ✓ where sold; leave blank otherwise. Keep rare/special items clearly flagged.

### Arms & Armor Inventory

| Item            | Britain | Moonglow | Jhelom | Yew | Minoc | Trinsic | Skara Brae | New Magincia | Cove | Bucs Den | Paws |
|-----------------|:-------:|:--------:|:------:|:---:|:-----:|:-------:|:----------:|:------------:|:----:|:--------:|:----:|
| Dagger          |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Sling           |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Club            |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Bow/Arrows      |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Crossbow/Quarrels|  —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Longsword       |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Halberd         |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| 2H Weapons      |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Magic Axe       |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Glass Sword     |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Mystic Sword    |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |

### Reagents Inventory

| Reagent        | Britain | Moonglow | Jhelom | Yew | Minoc | Trinsic | Skara Brae | New Magincia | Cove | Bucs Den | Paws |
|----------------|:-------:|:--------:|:------:|:---:|:-----:|:-------:|:----------:|:------------:|:----:|:--------:|:----:|
| Sulfurous Ash  |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Ginseng        |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Garlic         |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Spider Silk    |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Blood Moss     |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Black Pearl    |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Nightshade     |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Mandrake Root  |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |

### Guild Shop Inventory

| Item    | Britain | Moonglow | Jhelom | Yew | Minoc | Trinsic | Skara Brae | New Magincia | Cove | Bucs Den | Paws |
|---------|:-------:|:--------:|:------:|:---:|:-----:|:-------:|:----------:|:------------:|:----:|:--------:|:----:|
| Key     |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Torch   |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |
| Gem     |   —     |    —     |   —    | —   |   —   |    —    |     —      |      —       |  —   |    —     |  —   |

### Healer Services

| Town            | Heal (HP) | Cure Poison | Resurrect |
|-----------------|:---------:|:-----------:|:---------:|
| Britain         |     —     |      —      |     —     |
| Moonglow        |     —     |      —      |     —     |
| …               |     —     |      —      |     —     |

### Shipwright Services

| Town            | Skiff | Ship | Repair Hull |
|-----------------|:-----:|:----:|:-----------:|
| Britain         |   —   |  —   |      —      |
| Moonglow        |   —   |  —   |      —      |
| …               |   —   |  —   |      —      |

Notes:

- Cross‑link: using Gems triggers [Commands → View](./Commands.md#view-gem-map).
- Cross‑link: mixing spells uses reagents; see [Commands → Mix Reagents](./Commands.md#mix-reagents) and [Spells](./Spells.md).
- Deeper shoppes investigation required; keep these matrices as placeholders to be filled with canonical data.
