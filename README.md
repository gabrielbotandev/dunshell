# Dunshell

Dunshell is a single-player terminal roguelike built in Go with Bubble Tea, Lip Gloss, and Charm TUI tooling.

## Install And Run

```bash
go mod tidy
go run .
```

Replay a specific seed from the CLI:

```bash
go run . -seed 123456789
```

## Font And Terminal Notes

- Nerd Font is recommended for the full glyph language.
- Runtime glyph behavior now lives in the Settings menu from the title screen or with `p` in game.
- If your terminal/font renders symbols poorly, run with pure ASCII fallback:

```bash
DUNSHELL_ASCII=1 go run .
```

- `DUNSHELL_ASCII=1` still works as an override and will force ASCII even if the runtime menu is set differently.

## Controls

### Menus

- `up/down` or `w/s`: move selection
- `left/right` or `a/d`: switch mode where shown
- `Enter`: confirm
- `p`: settings
- `?`: help
- `Esc`: back
- `q`: quit where appropriate

### In Game

- `arrow keys` or `W/A/S/D`: move
- `.`: wait
- `c`: quick heal with the weakest healing consumable
- `e`: contextual interact
- `i`: open inventory
- `p`: open settings
- `left/right` or `A/D`: switch inventory pane
- `up/down` or `W/S`: move inside inventory or merchant stock
- `Enter` or `e`: confirm / buy / use primary action
- `u`: use selected consumable in the pack
- `?`: help
- `q`: safe quit prompt

## Save Location

Dunshell stores human-readable JSON saves under the platform config directory:

- Linux: `~/.config/dunshell/`
- macOS: `~/Library/Application Support/dunshell/`
- Windows: `%AppData%\dunshell\`

Files created by `0.3.1`:

- `profile.json`: persistent wins, omen-tier difficulty, and settings
- `run.json`: active run save
- `run.json.backup`: previous active run snapshot

## Version 0.3.1 Highlights

- 20-floor main progression with miniboss floors at `5`, `10`, `15`, and the Ashen Prior on `20`
- Bigger maps with improved Unicode glyph presentation and runtime glyph settings
- Dedicated boss rooms with confirmation prompts, lock-in logic, and visible boss health bars
- Route-choice map rendered as a proper node graph with detail preview
- Bronze, Silver, and Gold chests with matching keys and reward previews
- Merchant floors with curated five-slot stock and gold pricing
- Expanded weapons, armor, charms, consumables, rarities, and three chase unique items
- Stronger scaling, poison and fire pressure, curatives, and resistance gear
- Ultra-rare unique enemy drops while the final boss remains the premium reliable source
- Cleaner theme-friendly panels that respect the terminal background
- Auto-save / load persistence and seed entry flow from the UI

## Architecture

- `internal/game`: simulation, content, dungeon generation, AI, combat, progression, route modifiers, economy, and persistence data
- `internal/ui`: Bubble Tea screen state, terminal input, rendering, layout, glyph handling, and Lip Gloss styling

## Wiki

The detailed game reference now lives in `wiki/`:

- `wiki/progression.md`
- `wiki/monsters.md`
- `wiki/bosses.md`
- `wiki/items.md`
- `wiki/loot-and-rarities.md`
- `wiki/merchant-and-economy.md`
- `wiki/keys-and-chests.md`
- `wiki/lore.md`

## Build Status

`go test ./...` passes on `0.3.1`.
