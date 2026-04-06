# Changelog

## 0.4.0

### Changed

- Softened early-floor combat pressure by slightly reducing floor 1-2 enemy damage scaling, trimming early enemy counts, and increasing `Healing Salve` recovery so the opening run stays dangerous but fairer.
- Reworked early sustain and starting gear around a heavier opening loadout: five `Healing Salve`s, no starting `Sunbrew Tonic`, stronger low-tier weapons and armor, and broader item stat ranges.
- Slowed level progression by reducing effective XP earned from kills so the run leans more on loot and routing than rapid stat snowballing.
- Rebuilt level progression around tiny percentage-based XP gains per encounter, now factoring player level, enemy level, floor level, enemy threat, and current gear so early leveling no longer snowballs by floor 3.
- Expanded equipment variety with more weapons, armor pieces, and charms across the early and mid floors, plus slightly richer merchant stock and rare premium merchant rolls.
- Made merchant routes much rarer in route maps and reduced incidental merchant spawns so merchants stay a strategic spike instead of a near-every-floor expectation.
- Refactored `Sunbrew Tonic` and similar positive duration effects to work on floor duration instead of turn duration, keeping movement from wasting the item immediately.
- Increased global enemy health scaling and made bosses notably tougher so upgraded gear matters without flattening the game into easy mode.
- Added explicit enemy levels per floor and floating enemy health bars with level tags on the map to make encounter threat easier to read moment to moment.

## 0.3.2

### Added

- Added a CLI `-god` developer mode that starts runs with endgame testing gear, boosted stats, persisted god-mode saves, and full player invulnerability.

### Changed

- Changed chest prompts so reliquaries no longer reveal their spoils before you spend a matching key.
- Moved global version and god-mode indicators into the bottom footer and reduced footer controls to `?` help and `q` quit.
- Stopped persisting transient whisper-log history in run save files.

### Fixed

- Fixed boss gates so they reopen correctly after the boss is defeated.
- Fixed chest interaction copy in the status sidebar to match the hidden-spoils prompt flow.
- Fixed endless-mode enemy selection after floor 20 so late-game enemy pools continue spawning instead of falling back to rats.

## 0.3.1

### Added

- Added a persistent Settings screen from the title menu and with in-game `p`, covering glyph mode, ASCII fallback, descend confirmation, and log length.
- Added stronger poison and fire status support with clearer combat logs, sidebar indicators, curatives, and resistance gear.
- Added an ultra-rare unique enemy drop jackpot while keeping the final boss chest as the premium reliable source of uniques.

### Changed

- Rebuilt the route-choice screen into a real node-map presentation with branch highlighting and a dedicated detail panel.
- Increased combat pressure across the run through harder enemy baselines, sharper scaling, stronger bosses and elites, and lower passive sustain.
- Reworked UI chrome so panels, menus, logs, and gameplay read cleanly without forcing a custom full-screen background.

## 0.3

### Added

- Expanded the main run from 5 floors to 20 floors with themed floor names and escalating intros.
- Added miniboss cadence on floors 5, 10, and 15, plus the Ashen Prior as the final floor-20 boss.
- Added dedicated boss rooms with confirmation prompts, lock-in encounter logic, boss chest unlock flow, and visible boss health display.
- Added route-choice descent maps that modify the next floor with merchant, rest, elite, gold, chest, and loot-focused paths.
- Added Bronze, Silver, and Gold chests with matching key economy and reward preview prompts.
- Added merchant encounters with five curated offers, gold pricing, and progression-aware stock.
- Added many more weapons, armor pieces, charms, and consumables plus rarity tiers: Common, Uncommon, Rare, Legendary, and Unique.
- Added exactly three unique chase items: one weapon, one armor, and one charm.
- Added persistent profile progression with omen-tier difficulty increases after victories.
- Added post-victory endless continuation beyond floor 20.
- Added local JSON auto-save persistence with active-run backups.
- Added in-game seed selection flow with random or manual seed entry while preserving CLI `-seed` support.
- Added glyph fallback support through `DUNSHELL_ASCII=1`.
- Added `wiki/` documentation covering progression, monsters, bosses, items, loot, merchants, chests, and lore.

### Changed

- Increased map size substantially and refreshed the terminal presentation with a stronger glyph language.
- Reworked floor generation to support special rooms, merchants, boss chambers, and richer encounter layouts.
- Rebalanced combat and enemy scaling to keep healing, gear, and tactical pacing relevant through late-game floors.
- Reworked the Bubble Tea flow to include Continue, New Run seed input, merchant browsing, boss prompts, chest prompts, route selection, and richer outcome screens.
- Simplified saves into backup-friendly JSON under the platform config directory.
