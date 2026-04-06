# Monsters

## Core Enemies

| Enemy | Floors | Role | Base Stats |
| --- | --- | --- | --- |
| Gutter Rat | 1-3 | Skittish opener | HP 9, ATK 5, DEF 0 |
| Bone Beetle | 1-5 | Durable early brute | HP 15, ATK 6, DEF 2 |
| Lantern Wisp | 2-6 | Poison pressure | HP 13, ATK 7, DEF 1 |
| Mire Hound | 3-8 | Aggressive hunter, opens doors | HP 20, ATK 8, DEF 1 |
| Tomb Brigand | 4-9 | Gold thief | HP 20, ATK 8, DEF 2 |
| Cathedral Knight | 6-12 | Armored sentinel | HP 30, ATK 11, DEF 4 |
| Censer Acolyte | 7-15 | Ranged fire caster | HP 24, ATK 9, DEF 2 |
| Grave Sycophant | 9-16 | Midgame brute | HP 34, ATK 12, DEF 3 |
| Ash Archer | 10-18 | Fast ranged fire threat | HP 26, ATK 11, DEF 2 |
| Reliquary Ogre | 12-20 | Heavy late brute | HP 42, ATK 15, DEF 4 |
| Drowned Abbot | 14-20 | Blessed ranged poison pressure | HP 36, ATK 13, DEF 5 |
| Ember Seraph | 16-20 | Late hunter / fire threat | HP 32, ATK 15, DEF 3 |

## Combat Notes

- Normal enemies now have an explicit encounter level that tracks the floor: a rat on floor 1 is level 1, while a rat on floor 2 is level 2.
- Elites and bosses sit above the floor baseline, making them more rewarding and more dangerous for the XP model.
- Enemies scale upward by floor depth, and `0.4.0` adds a broader global health bump across the run.
- Floor 1-2 still avoid some of the extra attack pressure to keep the opening fairer.
- From the midgame onward, normal enemy attack scaling ramps harder than before.
- `Ashen Hunt` and `Cursed Procession` can promote enemies into elite versions.
- Elite enemies gain extra health, attack, defense, XP, and gold.
- Elites and bosses also deepen status durations and potency.
- Hunters punish sightlines and chase harder.
- Casters use ranged burst actions when given space.
- Cutpurses still threaten gold economy, which matters more now that merchants are rarer and more strategically valuable.

## Status Pressure

- `Poison`: damage over time from wisps, drowned rites, and elite escalations.
- `Fire`: burning damage over time from censer enemies, ash archers, ember seraphs, and late bosses.

## Readability

- Visible enemies now show a floating health bar directly above the sprite.
- The bar has no HP numbers, but it includes the enemy level beside it for quick threat reading.
