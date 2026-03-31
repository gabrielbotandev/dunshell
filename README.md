<p align="center">
  <img src=".github/assets/dunshell-banner.png" alt="Dunshell">
</p>

# Dunshell

Dunshell is a single-player terminal roguelike built in Go with Bubble Tea, Lip Gloss, and Charm TUI tooling.

## Play / Install

Download the archive that matches your platform from [GitHub Releases](https://github.com/gabrielbotandev/dunshell/releases/latest), extract it, and run the bundled binary:

- macOS and Linux: `./dunshell`
- Windows PowerShell: `.\dunshell.exe`

Package manager installs will come later. For now, GitHub Releases is the supported way to grab production binaries.

## Build from source

For contributors and local development, the source-build workflow stays the same:

```bash
go mod tidy
go run .
```

Replay a specific seed from the CLI:

```bash
go run . -seed 123456789
```

Start a developer testing run with endgame gear, boosted stats, persisted god-mode saves, and full player invulnerability:

```bash
go run . -god
```

## Font And Terminal Notes

- Nerd Font is recommended for the full glyph language.
- Runtime glyph behavior now lives in the Settings menu from the title screen or with `p` in game.
- If your terminal/font renders symbols poorly, run with pure ASCII fallback:

```bash
DUNSHELL_ASCII=1 ./dunshell
# or, from source:
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

Files created by the game:

- `profile.json`: persistent wins, omen-tier difficulty, and settings
- `run.json`: active run save
- `run.json.backup`: previous active run snapshot

God-mode runs keep their `GOD MODE` state inside the saved run, so a saved test run resumes as god mode while a normal save stays normal.

Run saves no longer persist the transient in-game message log.

## Game Highlights

- A full 20-floor main descent with miniboss floors on `5`, `10`, and `15`, then the Ashen Prior on floor `20`
- Dedicated boss rooms with warnings, lock-in gates, boss health tracking, and reward chests that unlock after victory
- Route-choice maps between floors, letting you steer into merchants, safer recovery paths, loot-heavy routes, cursed runs, or sharper combat
- Bronze, Silver, and Gold reliquary chests tied to a matching key economy instead of free loot
- Merchant stops, rarity-based gear, charms, consumables, curatives, and unique chase items for long-run progression
- Combat pressure built around poison, fire, elite enemies, scaling encounters, and resistance-focused equipment choices
- Endless mode after victory so a successful run can keep descending into harder post-game floors
- Auto-save / continue support, replayable seeds, runtime glyph settings, and ASCII fallback for terminals with weaker symbol support

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

`go test ./...` currently passes.
