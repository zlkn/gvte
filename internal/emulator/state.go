package emulator

import (
	"gvte/internal/emulator/grid"
	"gvte/internal/emulator/parser"
	"gvte/internal/emulator/selection"
)

// CursorShape represents visual cursor style.
type CursorShape int

const (
	CursorBlock CursorShape = iota
	CursorUnderline
	CursorBar
)

// Cursor defines current position and visual properties of the terminal cursor.
type Cursor struct {
	X       int
	Y       int
	Visible bool
	Shape   CursorShape
	Blink   bool
}

// State manages the core terminal state machine.
type State struct {
	Grid      *grid.Grid
	Parser    *parser.Parser
	Selection *selection.Selection
	Cursor    Cursor
	Title     string
}

// NewState creates a new terminal emulator state instance.
func NewState(cols, rows int) *State {
	return &State{
		Grid:      grid.New(cols, rows, 10000),
		Parser:    parser.New(),
		Selection: selection.New(),
		Cursor: Cursor{
			X:       0,
			Y:       0,
			Visible: true,
			Shape:   CursorBlock,
		},
		Title: "Terminal",
	}
}

// ProcessBytes processes raw bytes received from PTY.
func (s *State) ProcessBytes(data []byte) {
	s.Parser.Parse(data, func(seq parser.Sequence) {
		// Route parsed action sequence to grid update / state updates
		if seq.Type == parser.ActionPrint {
			if s.Cursor.X < s.Grid.Cols && s.Cursor.Y < s.Grid.Rows {
				s.Grid.Lines[s.Cursor.Y].Cells[s.Cursor.X].Char = seq.Char
				s.Cursor.X++
				if s.Cursor.X >= s.Grid.Cols {
					s.Cursor.X = 0
					s.Cursor.Y++
					if s.Cursor.Y >= s.Grid.Rows {
						s.Cursor.Y = s.Grid.Rows - 1
					}
				}
			}
		}
	})
}
