package font

import (
	"image"
	"testing"

	xfont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Go Mono at the config defaults (config.DefaultConfig: 14pt, 72 DPI) with
// HintingFull. The wanted metrics are measured values, kept as a regression
// anchor: they are what the atlas layout and the grid geometry are sized from.
const (
	testSize = 14.0
	testDPI  = 72.0

	wantCellWidth  = 8
	wantCellHeight = 17
	wantAscent     = 14
)

// The atlas bakes exactly this range (support-text-rendering.md, Step 3).
const (
	firstASCII = 0x20
	lastASCII  = 0x7E
)

func newTestManager(t *testing.T) *FontManager {
	t.Helper()

	fm, err := NewManager("Monospace", testSize, testDPI)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	t.Cleanup(func() { fm.Close() })
	return fm
}

// rasterize draws r into a coverage mask of one cell surrounded by pad px of
// slack, with the baseline at fm.Ascent. The slack is what makes ink outside
// the cell visible: without it the rasterizer would clip to Dst and any
// overflow would go unnoticed.
func rasterize(fm *FontManager, r rune, pad int) *image.Alpha {
	img := image.NewAlpha(image.Rect(0, 0, fm.CellWidth+2*pad, fm.CellHeight+2*pad))
	d := &xfont.Drawer{Dst: img, Src: image.Opaque, Face: fm.Face}
	d.Dot = fixed.P(pad, pad+fm.Ascent) // Dot is the baseline origin, not the top-left
	d.DrawString(string(r))
	return img
}

// inkBox returns the bounding box of non-zero coverage in cell coordinates
// (0,0 = cell top-left). ok is false when nothing was drawn at all.
func inkBox(img *image.Alpha, pad int) (box image.Rectangle, ok bool) {
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X-1, b.Min.Y-1

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if img.AlphaAt(x, y).A == 0 {
				continue
			}
			minX, minY = min(minX, x), min(minY, y)
			maxX, maxY = max(maxX, x), max(maxY, y)
		}
	}
	if maxX < minX {
		return image.Rectangle{}, false
	}
	return image.Rect(minX-pad, minY-pad, maxX-pad+1, maxY-pad+1), true
}

func TestNewManagerMetrics(t *testing.T) {
	fm := newTestManager(t)

	if fm.CellWidth != wantCellWidth || fm.CellHeight != wantCellHeight || fm.Ascent != wantAscent {
		t.Errorf("cell %dx%d ascent=%d, want %dx%d ascent=%d",
			fm.CellWidth, fm.CellHeight, fm.Ascent,
			wantCellWidth, wantCellHeight, wantAscent)
	}

	// main.go and Window.onResize divide by these; zero here is a panic there.
	if fm.CellWidth <= 0 || fm.CellHeight <= 0 {
		t.Fatalf("cell metrics must be positive, got %dx%d", fm.CellWidth, fm.CellHeight)
	}
	if fm.Ascent <= 0 || fm.Ascent > fm.CellHeight {
		t.Errorf("baseline at y=%d lies outside the %d px cell", fm.Ascent, fm.CellHeight)
	}
}

func TestCellWidthMatchesAdvance(t *testing.T) {
	fm := newTestManager(t)

	want, ok := fm.Face.GlyphAdvance('M')
	if !ok {
		t.Fatal("face reports no advance for 'M'")
	}
	if want.Ceil() != fm.CellWidth {
		t.Errorf("CellWidth=%d, but advance('M')=%v (%d px)", fm.CellWidth, want, want.Ceil())
	}

	// A terminal grid is only valid if every baked glyph advances identically.
	for r := rune(firstASCII); r <= lastASCII; r++ {
		adv, ok := fm.Face.GlyphAdvance(r)
		if !ok {
			t.Errorf("rune %q (%#x): face has no glyph", r, r)
			continue
		}
		if adv != want {
			t.Errorf("rune %q: advance %v, want %v — face is not monospaced at this size", r, adv, want)
		}
	}
}

func TestGlyphRasterizes(t *testing.T) {
	fm := newTestManager(t)

	img := rasterize(fm, 'A', 0)

	ink, belowBaseline := 0, 0
	for y := range fm.CellHeight {
		row := make([]byte, fm.CellWidth)
		for x := range fm.CellWidth {
			if img.AlphaAt(x, y).A > 128 {
				row[x] = '#'
				ink++
				if y >= fm.Ascent {
					belowBaseline++
				}
			} else {
				row[x] = '.'
			}
		}
		t.Logf("%2d %s", y, row)
	}

	if ink == 0 {
		t.Fatal("'A' rasterized empty: either the face is broken or the baseline sits outside the cell")
	}
	// 'A' has no descender, so ink under the baseline means Ascent is too small
	// and the whole atlas would be drawn a few pixels too low.
	if belowBaseline > 0 {
		t.Errorf("%d px of 'A' below the baseline at y=%d: Ascent is wrong", belowBaseline, fm.Ascent)
	}
}

func TestInkFitsAtlasCell(t *testing.T) {
	fm := newTestManager(t)

	// Measured on Go Mono @14pt/72dpi: ink never leaves the cell vertically,
	// but 21 of the 95 glyphs ('M', '%', '_', 'W', ...) bleed one column past
	// the advance width. The atlas needs that column as a gutter, otherwise a
	// neighbouring slot picks up the spill.
	const (
		pad            = 4
		maxRightBleed  = 1
		wantBleedRunes = 21
		maxReported    = 5 // a broken Ascent breaks all 95 glyphs identically
	)

	cell := image.Rect(0, 0, fm.CellWidth, fm.CellHeight)
	allowed := image.Rect(0, 0, fm.CellWidth+maxRightBleed, fm.CellHeight)

	union := image.Rectangle{}
	bleeding, escaped, reported := 0, 0, 0
	for r := rune(firstASCII); r <= lastASCII; r++ {
		box, ok := inkBox(rasterize(fm, r, pad), pad)
		if !ok {
			continue // blank glyph, e.g. space
		}
		union = union.Union(box)

		if !box.In(allowed) {
			escaped++
			if reported < maxReported {
				reported++
				t.Errorf("rune %q: ink box %v escapes the %v cell (right-bleed budget %d px)",
					r, box, cell, maxRightBleed)
			}
		}
		if !box.In(cell) {
			bleeding++
		}
	}
	if escaped > reported {
		t.Errorf("and %d more runes escape the cell", escaped-reported)
	}

	t.Logf("ink union across %#x-%#x: %v, cell %v, %d glyphs bleed right",
		firstASCII, lastASCII, union, cell, bleeding)

	// Not a correctness requirement — a tripwire for a font or metric change
	// that silently alters how much gutter the atlas needs.
	if bleeding != wantBleedRunes {
		t.Errorf("%d glyphs bleed outside the cell, want %d", bleeding, wantBleedRunes)
	}
}
