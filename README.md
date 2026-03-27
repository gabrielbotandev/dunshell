# Dunshell

Dunshell is a single-player terminal roguelike built in Go with Bubble Tea, Lip Gloss, and Charm's TUI tooling. You descend through a drowned abbey, fight through procedural dungeon floors, grow stronger with weapons, armor, charms, and consumables, and hunt the buried relic known as the Cinder Crown.

## Overview

- Turn-based dungeon crawling on a procedural grid map
- Fog of war with line-of-sight field of view
- Cleared-room tracking with visual completion feedback
- Multiple themed floors with rising difficulty
- Inventory, equipment, consumables, loot drops, leveling, and room-clearing feedback
- Distinct enemies with wandering, chasing, and attack behaviors
- Bubble Tea-driven title, help, gameplay, victory, and defeat screens
- Reproducible run seeds with an optional CLI flag

## Architecture

The project is split into two main layers:

- `internal/game`: the simulation layer for dungeon generation, map state, actors, combat, AI, items, field of view, progression, and the message log
- `internal/ui`: the Bubble Tea layer for screen state, terminal input, rendering, layout, and Lip Gloss styling

This keeps the game rules independent from the terminal presentation, making the dungeon systems easier to extend without tangling rendering code into the core loop.

## File Structure

```text
.
|-- go.mod
|-- main.go
|-- README.md
`-- internal
    |-- game
    |   |-- actor.go
    |   |-- content.go
    |   |-- fov.go
    |   |-- game.go
    |   |-- generator.go
    |   |-- item.go
    |   |-- map.go
    |   |-- path.go
    |   `-- types.go
    `-- ui
        |-- model.go
        |-- render.go
        `-- styles.go
```

## Controls

### Menus

- `up/down` or `w/s`: move selection
- `Enter`: confirm
- `?`: open help
- `q` or `Esc`: quit or go back where appropriate
- Letter keys accept lowercase and uppercase

### In Game

- `arrow keys` or `W/A/S/D`: move
- `.`: wait one turn
- `c`: quick heal with the weakest healing consumable
- `e`: contextually interact with your tile
- `e` on loot: pick it up
- `e` on stairs: open the descend confirmation prompt
- `i`: open inventory in the side panel
- `left/right` or `A/D`: switch inventory section
- `up/down` or `W/S`: move inside the inventory
- `e` or `Enter`: perform the primary inventory action
- `u`: use a consumable in the inventory
- `?`: help screen
- `q`: safe quit prompt
- Letter keys accept lowercase and uppercase

## How To Run

```bash
go mod tidy
go run .
```

To replay a specific seed:

```bash
go run . -seed 123456789
```

## Design Summary

- Dungeon generation uses room-and-corridor layouts with doors, loot placement, enemy placement, and a special final sanctum floor objective.
- The gameplay loop is tuned around short tactical runs: fight, loot, clear chambers, manage consumables, descend, level up, and survive status effects.
- Bubble Tea manages screen transitions and input while Lip Gloss provides the visual hierarchy for map, sidebar, log, menus, and end screens.
- The map viewport crops around the player so the game still plays well when the terminal is smaller than the generated dungeon.

## Content Highlights

- Enemies: Gutter Rat, Bone Beetle, Lantern Wisp, Mire Hound, Tomb Brigand, Cathedral Knight, and the Ashen Prior
- Items: healing salves, antivenom, Sunbrew tonic, Ember flask, multiple weapons, armor sets, and charms
- Objective: reach the Ember Sanctum and claim the Cinder Crown
