package config

import (
	"fmt"
	"image/color"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

var Debug bool

type ColorScheme struct {
	Background color.Color
	Foreground color.Color
	Cursor     color.Color
	Selection  color.Color
	Ansi       [16]color.Color
}

type FontConfig struct {
	Family string
	Size   float64
	DPI    float64
}

type Config struct {
	Font          FontConfig
	Colors        ColorScheme
	InitialWidth  int
	InitialHeight int
	Shell         string
}

func ColorToWgpu(c color.Color) wgpu.Color {
	r, g, b, a := c.RGBA() // Возвращает uint32 от 0 до 65535
	return wgpu.Color{
		R: float64(r) / 65535.0,
		G: float64(g) / 65535.0,
		B: float64(b) / 65535.0,
		A: float64(a) / 65535.0,
	}
}

func DefaultConfig() *Config {
	return &Config{
		Font: FontConfig{
			Family: "Monospace",
			Size:   14.0,
			DPI:    72.0,
		},
		Colors: ColorScheme{
			Background: color.RGBA{R: 0x1e, G: 0x1e, B: 0x2e, A: 0xff}, // Catppuccin Mocha Base
			Foreground: color.RGBA{R: 0xcd, G: 0xd6, B: 0xf4, A: 0xff}, // Text
			Cursor:     color.RGBA{R: 0xf5, G: 0xe0, B: 0xdc, A: 0xff}, // Rosewater
			Selection:  color.RGBA{R: 0x58, G: 0x5b, B: 0x70, A: 0x80}, // Surface2
		},
		InitialWidth:  900,
		InitialHeight: 600,
		Shell:         "/bin/bash",
	}
}

// TODO: Load loads application configuration from disk or defaults.
func Load(path string) (*Config, error) {
	// Stub for loading config from path (e.g. YAML/TOML/JSON)
	if path == "" {
		return DefaultConfig(), nil
	}
	return nil, fmt.Errorf("custom config path loading not yet implemented")
}
