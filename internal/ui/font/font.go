package font

import (
	"fmt"
	"image"
)

type GlyphCache struct {
	cache map[rune]image.Image
}

type FontManager struct {
	Family     string
	Size       float64
	CellWidth  int
	CellHeight int
	Glyphs     *GlyphCache
}

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

func (fm *FontManager) GetGlyph(r rune) (image.Image, error) {
	if img, ok := fm.Glyphs.cache[r]; ok {
		return img, nil
	}
	//TODO Stub rasterization
	return nil, fmt.Errorf("glyph rasterization stub for rune '%c'", r)
}
