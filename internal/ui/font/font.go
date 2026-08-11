package font

import (
	"fmt"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"image"
)

type GlyphCache struct {
	cache map[rune]image.Image
}

type FontManager struct {
	Family     string
	Size       float64
	Face       xfont.Face
	CellWidth  int
	CellHeight int
	Ascent     int // нужен на Step 3: Dot — это базовая линия, не верх слота
	Glyphs     *GlyphCache
}

// NewManager builds a face at size points and dpi
// dots per inch, then derives the terminal cell box from the font's own
// metrics. family is recorded for reporting only: the face is compiled in, so
// the configured family name selects nothing yet.
func NewManager(family string, size, dpi float64) (*FontManager, error) {
	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return nil, fmt.Errorf("parse gomono: %w", err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: size, DPI: dpi, Hinting: xfont.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create face: %w", err)
	}

	// Monospace: every glyph advances identically, so we able to use any letter
	m := face.Metrics()
	adv, ok := face.GlyphAdvance('A')
	if !ok || adv == 0 {
		face.Close()
		return nil, fmt.Errorf("font has no advance for 'M'")
	}

	fm := &FontManager{
		Family:     family,
		Size:       size,
		Face:       face,
		CellWidth:  adv.Ceil(), // fixed.Int26_6 is 1/64 px; never cast directly
		CellHeight: m.Height.Ceil(),
		Ascent:     m.Ascent.Ceil(),
		Glyphs: &GlyphCache{
			cache: make(map[rune]image.Image),
		},
	}
	return fm, nil
}

func (fm *FontManager) Close() error {
	if fm.Face == nil {
		return nil
	}
	face := fm.Face
	fm.Face = nil
	return face.Close()
}

func (fm *FontManager) GetGlyph(r rune) (image.Image, error) {
	if img, ok := fm.Glyphs.cache[r]; ok {
		return img, nil
	}
	//TODO Stub rasterization
	return nil, fmt.Errorf("glyph rasterization stub for rune '%c'", r)
}
