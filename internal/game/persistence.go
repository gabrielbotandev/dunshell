package game

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const SaveVersion = GameVersion

var ErrNoRunSave = errors.New("no active run save")

type Profile struct {
	Version    string   `json:"version"`
	Wins       int      `json:"wins"`
	Difficulty int      `json:"difficulty"`
	Settings   Settings `json:"settings,omitempty"`
	UpdatedAt  string   `json:"updated_at"`
}

type RunState struct {
	Version              string        `json:"version"`
	SavedAt              string        `json:"saved_at"`
	Seed                 int64         `json:"seed"`
	GodMode              bool          `json:"god_mode,omitempty"`
	Mode                 GameMode      `json:"mode"`
	FloorIndex           int           `json:"floor_index"`
	MaxFloors            int           `json:"max_floors"`
	Turn                 int           `json:"turn"`
	Player               *Player       `json:"player"`
	Floor                *Floor        `json:"floor"`
	PersistentDifficulty int           `json:"persistent_difficulty"`
	Endless              bool          `json:"endless"`
	PendingRoutes        []RouteChoice `json:"pending_routes,omitempty"`
	VictoryRecorded      bool          `json:"victory_recorded"`
	NextEnemyID          int           `json:"next_enemy_id"`
	NextChestID          int           `json:"next_chest_id"`
	NextMerchantID       int           `json:"next_merchant_id"`
	RNGState             uint64        `json:"rng_state"`
}

func SaveDirectory() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "dunshell")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func LoadProfile() (Profile, error) {
	path, err := profileFilePath()
	if err != nil {
		return Profile{}, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Profile{Version: SaveVersion, Settings: DefaultSettings()}, nil
	}
	if err != nil {
		return Profile{}, err
	}
	profile := Profile{Version: SaveVersion}
	if err := json.Unmarshal(data, &profile); err != nil {
		return Profile{}, err
	}
	if profile.Version == "" {
		profile.Version = SaveVersion
	}
	profile.Settings = profile.Settings.Normalized()
	return profile, nil
}

func SaveProfile(profile Profile) error {
	profile.Version = SaveVersion
	profile.Settings = profile.Settings.Normalized()
	profile.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	path, err := profileFilePath()
	if err != nil {
		return err
	}
	return writeJSON(path, profile)
}

func LoadRun() (*Game, error) {
	path, err := runFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNoRunSave
	}
	if err != nil {
		return nil, err
	}
	var state RunState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return GameFromState(state), nil
}

func SaveRun(game *Game) error {
	path, err := runFilePath()
	if err != nil {
		return err
	}
	backup := path + ".backup"
	if _, err := os.Stat(path); err == nil {
		_ = copyFile(path, backup)
	}
	return writeJSON(path, game.RunState())
}

func ClearRun() error {
	path, err := runFilePath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	backup := path + ".backup"
	if err := os.Remove(backup); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func GameFromState(state RunState) *Game {
	game := &Game{
		Title:                GameTitle,
		Seed:                 state.Seed,
		GodMode:              state.GodMode,
		Log:                  make([]string, 0, 128),
		Mode:                 state.Mode,
		FloorIndex:           state.FloorIndex,
		MaxFloors:            state.MaxFloors,
		Turn:                 state.Turn,
		Player:               state.Player,
		Floor:                state.Floor,
		PersistentDifficulty: state.PersistentDifficulty,
		Endless:              state.Endless,
		PendingRoutes:        state.PendingRoutes,
		VictoryRecorded:      state.VictoryRecorded,
		rng:                  &RNG{State: state.RNGState},
		nextEnemyID:          state.NextEnemyID,
		nextChestID:          state.NextChestID,
		nextMerchantID:       state.NextMerchantID,
	}
	if game.Player == nil {
		return New(NewGameOptions{Seed: state.Seed, PersistentDifficulty: state.PersistentDifficulty, GodMode: state.GodMode})
	}
	hydrateProgressionState(game)
	game.restoreGodModeState()
	if game.Floor != nil {
		ComputeFOV(game.Floor, game.Player.Pos, game.Player.VisionRadius())
	}
	return game
}

func (g *Game) RunState() RunState {
	state := uint64(0)
	if g.rng != nil {
		state = g.rng.State
	}
	return RunState{
		Version:              SaveVersion,
		SavedAt:              time.Now().UTC().Format(time.RFC3339),
		Seed:                 g.Seed,
		GodMode:              g.GodMode,
		Mode:                 g.Mode,
		FloorIndex:           g.FloorIndex,
		MaxFloors:            g.MaxFloors,
		Turn:                 g.Turn,
		Player:               g.Player,
		Floor:                g.Floor,
		PersistentDifficulty: g.PersistentDifficulty,
		Endless:              g.Endless,
		PendingRoutes:        g.PendingRoutes,
		VictoryRecorded:      g.VictoryRecorded,
		NextEnemyID:          g.nextEnemyID,
		NextChestID:          g.nextChestID,
		NextMerchantID:       g.nextMerchantID,
		RNGState:             state,
	}
}

func profileFilePath() (string, error) {
	dir, err := SaveDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "profile.json"), nil
}

func runFilePath() (string, error) {
	dir, err := SaveDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "run.json"), nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	temp := path + ".tmp"
	if err := os.WriteFile(temp, append(data, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(temp, path)
}

func copyFile(src string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}
