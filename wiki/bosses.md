# Bosses

Boss floors seal the fight behind a confirmation prompt. Once the player commits, the gate closes and stays sealed until the boss dies.

In `0.4.0`, bosses receive noticeably heavier encounter scaling than normal enemies. Their authored stats are only the starting point; the actual floor encounter is significantly tougher in HP, damage, defense, and burst pressure.

## Minibosses And Final Boss

| Boss | Floor | Authored Base Stats | Signature Pressure | Reward Chest |
| --- | --- | --- | --- | --- |
| Houndmaster Vey | 5 | HP 64, ATK 13, DEF 4 | `Ruinous Pounce` closes distance and punishes weak early armor | Silver |
| Bell Archivist Oria | 10 | HP 92, ATK 16, DEF 5 | `Chime of Ruin` forces healing and punishes ranged drift | Gold |
| Censer Matriarch | 15 | HP 118, ATK 19, DEF 6 | `Incense Storm` layers attrition through poison pressure | Gold |
| Ashen Prior | 20 | HP 165, ATK 23, DEF 8 | `Funeral Litany`, scorch pressure, and enrage at low health | Gold + Crown |

## Encounter Rules

- Boss rooms contain only the boss and its reward chest.
- The player can postpone the chamber and finish the rest of the floor first.
- Confirming entry seals the gate.
- Boss chests remain locked until the boss dies.
- The boss drops the matching key for the reward chest.
- The final reward chest contains the `Cinder Crown` and one random unique item.
- Bosses now scale harder than they did in earlier versions, so entering undergeared should feel dangerous again.

## Visual Presentation

- Boss floors tint the chamber differently from normal rooms.
- Bosses use dedicated glyph treatment and a visible boss HP bar in the sidebar.
- The final boss uses a stronger icon than standard enemies and feels distinct on sight.
