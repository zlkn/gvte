package font

import (
	"fmt"
	"image"
)

// GlyphCache stores pre-rendered glyph images for quick lookup.
type GlyphCache struct {
	cache map[rune]image.Image
}

// FontManager handles font loading, metrics, glyph rasterization and caching.
type FontManager struct {
	Family    string
	Size      float64
	CellWidth  int
	CellHeight int
	Glyphs    *GlyphCache
}

// NewManager initializes font settings and glyph cache.
func NewManager(family string, size float64) (*FontManager, error) {
	fm := &FontManager{
		Family:     family,
		Size:       size,
		CellWidth:  9,
		CellHeight: 18,
		Glyphs: &GlyphCache{
			cache: make(map[rune]image.Image),
		},
	}
	return fm, nil
}

// GetGlyph returns rasterized glyph image or renders and caches it if missing.
func (fm *FontManager) GetGlyph(r rune) (image.Image, error) {
	if img, ok := fm.Glyphs.cache[r]; ok {
		return img, nil
	}
	// Stub rasterization
	return nil, fmt.Errorf("glyph rasterization stub for rune '%c'", r)
}
