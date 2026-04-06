# Progression

## Main Run Structure

Dunshell `0.4.0` is built around a 20-floor main descent.

| Floor Range | Theme Direction | Notes |
| --- | --- | --- |
| 1-4 | Flooded crypts and pilgrim corridors | Early gear ramp, Bronze/Silver chest economy starts here |
| 5 | Kennel Reliquary | First miniboss floor |
| 6-9 | Cloisters, trenches, archive depths | Merchant and elite routes start to matter more |
| 10 | Resonant Nave | Second miniboss floor |
| 11-14 | Black refectory to thorn processional | Rare and Legendary loot become realistic chase targets |
| 15 | Saint's Wake | Third miniboss floor |
| 16-19 | Char vaults and crown approach | Late-game attrition and elite pressure spike hard |
| 20 | Throne of Ash | Ashen Prior, final boss chest, Cinder Crown |

Early runs in `0.4.0` are intentionally less level-driven than before:

- kill XP is lower, so gear and routing matter more than fast early level spikes
- floor 1-2 still give breathing room, but the broader run scales harder through enemy health, damage, and bosses

## Experience Model

`0.4.0` uses percentage-based XP progress inside each level instead of the older raw-integer pacing.

- each level now represents `100%` progress
- kills award only a small slice of that bar
- the reward is influenced by player level, enemy level, area level, enemy threat, and current gear power
- higher-threat targets like elites, casters, and bosses grant more progress than weak trash like rats
- as the player level rises, the same fight contributes a smaller portion of the next level

The goal is to keep leveling relevant without letting floor 2-3 snowball the run into free sustain.

## Boss Cadence

- Floor 5: `Houndmaster Vey`
- Floor 10: `Bell Archivist Oria`
- Floor 15: `Censer Matriarch`
- Floor 20: `Ashen Prior`

Every boss floor includes a dedicated sealed chamber with a boss chest. Boss rooms are skipped until the player deliberately confirms entry.

## Route Map

After confirming stairs, the game opens a route-choice map instead of instantly generating the next floor. Each route changes the next floor while preserving the fixed boss cadence.

Route pool in `0.4.0`:

| Route | Effect |
| --- | --- |
| Gilded Way | More gold and a Bronze key on arrival |
| Broker's Lantern | Merchant guaranteed on next floor |
| Pilgrim's Rest | Heal, cleanse, lighter enemy pressure |
| Reliquary Breach | Extra chest, stronger loot, Silver key on arrival |
| Ashen Hunt | More elites and better drops |
| Cursed Procession | Hardest modifier, more gold, extra chest, richer loot |

Merchant routing is no longer close to guaranteed:

- `Broker's Lantern` only appears on about `25%` of route maps
- normal merchant spawns are rarer than before
- a merchant route is now a meaningful economic event, not a default expectation

## Endless Mode

Claiming the Cinder Crown unlocks a victory choice:

- End the run and return to the title screen
- Continue into endless floors beyond 20

Endless floors continue route-choice generation and keep spawning boss depths on later multiples of five.

## Omen Tier

Each successful crown claim increases the persistent `omen tier` stored in `profile.json`.

- Wins increase omen tier by `+1`
- Deaths do not reduce it
- New runs start at the current omen tier

Omen tier increases enemy durability and pressure across the whole run.
