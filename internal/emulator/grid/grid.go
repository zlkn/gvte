package grid

import (
	"image/color"
)

// AttributeFlags represents styling flags (bold, italic, underline, inverse, etc.).
type AttributeFlags uint16

const (
	AttrBold AttributeFlags = 1 << iota
	AttrItalic
	AttrUnderline
	AttrInverse
)

// Cell represents a single character cell in the terminal grid.
type Cell struct {
	Char       rune
	FgColor    color.Color
	BgColor    color.Color
	Attributes AttributeFlags
	Width      int // 1 for standard width, 2 for wide CJK runes, 0 for continuation
}

// Line represents a single row of cells.
type Line struct {
	Cells []Cell
}

// NewLine creates a Line with specified number of columns.
func NewLine(cols int) Line {
	cells := make([]Cell, cols)
	for i := range cells {
		cells[i] = Cell{Char: ' ', Width: 1}
	}
	return Line{Cells: cells}
}

// Grid manages the 2D character matrix and scrollback buffer.
type Grid struct {
	Cols       int
	Rows       int
	Lines      []Line
	Scrollback []Line
	MaxScroll  int
}

// New creates a new character grid with dimensions and maximum scrollback lines.
func New(cols, rows, maxScrollback int) *Grid {
	lines := make([]Line, rows)
	for i := range lines {
		lines[i] = NewLine(cols)
	}
	return &Grid{
		Cols:       cols,
		Rows:       rows,
		Lines:      lines,
		Scrollback: make([]Line, 0, maxScrollback),
		MaxScroll:  maxScrollback,
	}
}

// Resize adjusts grid dimensions.
func (g *Grid) Resize(cols, rows int) {
	g.Cols = cols
	g.Rows = rows
	newLines := make([]Line, rows)
	for i := range newLines {
		if i < len(g.Lines) {
			newLines[i] = g.Lines[i]
			if len(newLines[i].Cells) < cols {
				extra := cols - len(newLines[i].Cells)
				for e := 0; e < extra; e++ {
					newLines[i].Cells = append(newLines[i].Cells, Cell{Char: ' ', Width: 1})
				}
			} else if len(newLines[i].Cells) > cols {
				newLines[i].Cells = newLines[i].Cells[:cols]
			}
		} else {
			newLines[i] = NewLine(cols)
		}
	}
	g.Lines = newLines
}
