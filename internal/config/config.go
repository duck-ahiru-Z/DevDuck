package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
)

type Config struct {
	Level teaching.Level `json:"level"`
}

func Default() Config {
	return Config{
		Level: teaching.Intermediate,
	}
}

func Path() (string, error) {
	base, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(
		base,
		"DevDuck",
		"config.json",
	), nil
}

func CachePath() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "DevDuck", "explanations.json"), nil
}

func Load() (Config, error) {
	path, err := Path()

	if err != nil {
		return Default(), err
	}

	data, err := os.ReadFile(path)

	if errors.Is(err, os.ErrNotExist) {
		return Default(), nil
	}

	if err != nil {
		return Default(), err
	}

	cfg := Default()

	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), err
	}

	level, err := teaching.ParseLevel(
		string(cfg.Level),
	)

	if err != nil {
		return Default(), err
	}

	cfg.Level = level

	return cfg, nil
}

func Save(cfg Config) error {
	path, err := Path()

	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		cfg,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		path,
		data,
		0600,
	)
}
