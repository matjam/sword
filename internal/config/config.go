package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

var globalConfig *Config

type Assets struct {
	Images   map[string]string        `json:"images"`
	Fonts    map[string]FontConfig    `json:"fonts"`
	Tilesets map[string]TilesetConfig `json:"tilesets"`
}

type FontConfig struct {
	Path string  `json:"path"`
	Size float64 `json:"size"`
}

type TilesetConfig struct {
	Path      string            `json:"path"`
	TileSize  int               `json:"tile_size"`
	Columns   int               `json:"columns"`
	Rows      int               `json:"rows"`
	Autotiles [][2]int          `json:"autotiles"`
	Fixtures  map[string][2]int `json:"fixtures"`
}

type Config struct {
	Assets Assets `json:"assets"`
}

func Load() *Config {
	if globalConfig != nil {
		return globalConfig
	}

	assetsData, err := os.ReadFile("assets.json")
	if err != nil {
		slog.Info("error reading assets.json", err)
		panic(err)
	}

	config := Config{}
	err = json.Unmarshal(assetsData, &config.Assets)
	if err != nil {
		slog.Info("error reading assets.json", err)
		panic(err)
	}

	globalConfig = &config

	return globalConfig
}
