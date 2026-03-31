# Contributing to Dunshell

Thanks for helping improve Dunshell.

This guide covers the practical workflow for contributing code, docs, balance changes, and fixes.

## Before You Start

- For larger features, gameplay reworks, or release-process changes, open an issue or start a discussion before spending a lot of time on implementation.
- Prefer small, focused pull requests over mixed changes that touch unrelated systems.
- If you change visible gameplay, controls, UI flows, or release behavior, update the relevant docs in the same branch.

## Development Setup

Dunshell currently targets Go `1.24`.

Install dependencies and run the game locally:

```bash
go mod tidy
go run .
```

Run the test suite before opening a pull request:

```bash
go test ./...
```

Useful development commands:

- Replay a specific seed: `go run . -seed 123456789`
- Start a god-mode testing run: `go run . -god`
- Force ASCII rendering for terminal compatibility checks: `DUNSHELL_ASCII=1 go run .`

## Project Layout

- `main.go`: CLI entrypoint and Bubble Tea program startup
- `internal/game`: game rules, dungeon generation, combat, progression, loot, persistence, and simulation data
- `internal/ui`: Bubble Tea models, rendering, layout, controls, and terminal presentation
- `wiki/`: player-facing reference material for systems, monsters, items, and lore

## Issues

When opening an issue, include enough detail for someone else to reproduce and verify the problem:

- what happened
- what you expected to happen
- steps to reproduce it
- your platform and terminal when relevant
- screenshots, logs, or save details if they help explain the problem

For feature requests, explain the player or developer problem the change would solve.

## Making Changes

- Keep changes as small as you can while still solving the problem cleanly.
- Avoid unrelated refactors in the same pull request.
- Follow the existing project style in the files you touch instead of introducing a new pattern for a one-off change.
- If you add or change flags, controls, save behavior, release flow, or installation instructions, update `README.md` and any affected `wiki/` pages.
- If you touch persistence, be careful with compatibility and verify that saving and resuming still behaves correctly.

## Testing Expectations

Before you submit a pull request, do the checks that match your change:

- Run `go test ./...`
- Manually play the area of the game you changed
- If you changed terminal rendering or glyph behavior, test ASCII mode with `DUNSHELL_ASCII=1`
- If you changed controls, menus, inventory, route maps, or status displays, verify the updated flow in-game
- If you changed saves or progression, verify a fresh run and a resumed save when relevant

## Pull Requests

Include the following in your pull request description:

- what changed
- why the change was needed
- how you tested it
- screenshots or terminal captures for UI-heavy changes when they help review

If your pull request fixes an issue, link it directly.

## Community Expectations

- Be respectful, direct, and constructive in issues, pull requests, and reviews.
- Assume good intent and focus feedback on the change, not the person.
- Keep discussions actionable and grounded in the game, the codebase, and user impact.

## Releases

GitHub Releases are created from semver tags through GoReleaser.

- Release config: `.goreleaser.yaml`
- Release workflow: `.github/workflows/release.yml`

If you change the release process, keep the local snapshot path working:

```bash
goreleaser release --snapshot --clean
```
