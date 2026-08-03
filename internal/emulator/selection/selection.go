package selection

// Pos represents a grid coordinate (Column, Row).
type Pos struct {
	Col int
	Row int
}

// Mode defines selection modes (character, word, line).
type Mode int

const (
	ModeNone Mode = iota
	ModeChar
	ModeWord
	ModeLine
)

// Selection handles text selection ranges and system clipboard interaction.
type Selection struct {
	Start    Pos
	End      Pos
	Active   bool
	Mode     Mode
}

// New creates a new Selection tracker.
func New() *Selection {
	return &Selection{
		Mode: ModeNone,
	}
}

// StartSelection begins a new selection at pos.
func (s *Selection) StartSelection(pos Pos, mode Mode) {
	s.Start = pos
	s.End = pos
	s.Active = true
	s.Mode = mode
}

// UpdateSelection updates the current endpoint of selection.
func (s *Selection) UpdateSelection(pos Pos) {
	if s.Active {
		s.End = pos
	}
}

// Clear resets the active selection.
func (s *Selection) Clear() {
	s.Active = false
	s.Mode = ModeNone
}

// Contains checks if a grid coordinate is within current selection bounds.
func (s *Selection) Contains(col, row int) bool {
	if !s.Active {
		return false
	}
	// Simplified boundary check logic
	startRow, endRow := s.Start.Row, s.End.Row
	startCol, endCol := s.Start.Col, s.End.Col
	if startRow > endRow || (startRow == endRow && startCol > endCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	if row < startRow || row > endRow {
		return false
	}
	if row == startRow && col < startCol {
		return false
	}
	if row == endRow && col > endCol {
		return false
	}
	return true
}
