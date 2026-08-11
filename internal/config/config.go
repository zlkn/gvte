package config

import (
	"fmt"
	"image/color"
	"math"

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

// IsSrgbFormat reports whether the GPU applies the sRGB encoding transfer
// function when writing to this format.
func IsSrgbFormat(f wgpu.TextureFormat) bool {
	switch f {
	case wgpu.TextureFormat_RGBA8UnormSrgb, wgpu.TextureFormat_BGRA8UnormSrgb:
		return true
	default:
		return false
	}
}

// srgbToLinear applies the inverse sRGB transfer function to one channel.
func srgbToLinear(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// ColorToRGBA normalizes c to [0,1]. Hex color literals are sRGB-encoded, but
// wgpu treats both clear values and fragment output as linear and re-encodes
// them for *UnormSrgb targets, so those must be decoded first or they render
// far too bright. Assumes opaque colors: color.Color.RGBA is alpha-premultiplied
// and this does not unpremultiply before decoding.
func ColorToRGBA(c color.Color, srgbTarget bool) [4]float64 {
	r, g, b, a := c.RGBA() // uint32 in 0..65535
	out := [4]float64{
		float64(r) / 65535.0,
		float64(g) / 65535.0,
		float64(b) / 65535.0,
		float64(a) / 65535.0,
	}
	if srgbTarget {
		for i := range 3 {
			out[i] = srgbToLinear(out[i])
		}
	}
	return out
}

// ColorToWgpu converts c into a wgpu clear value for a target of the given format.
func ColorToWgpu(c color.Color, srgbTarget bool) wgpu.Color {
	v := ColorToRGBA(c, srgbTarget)
	return wgpu.Color{R: v[0], G: v[1], B: v[2], A: v[3]}
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
