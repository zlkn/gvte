package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// InputMapper translates user keyboard and mouse events into ANSI/VT escape sequences.
type InputMapper struct{}

// NewMapper creates a new InputMapper.
func NewMapper() *InputMapper {
	return &InputMapper{}
}

// MapKey converts an Ebitengine keypress into sequence bytes to send to PTY.
func (im *InputMapper) MapKey(key ebiten.Key, mods ebiten.Key) []byte {
	switch key {
	case ebiten.KeyEnter:
		return []byte("\r")
	case ebiten.KeyBackspace:
		return []byte{0x7f}
	case ebiten.KeyTab:
		return []byte("\t")
	case ebiten.KeyEscape:
		return []byte{0x1b}
	case ebiten.KeyUp:
		return []byte("\x1b[A")
	case ebiten.KeyDown:
		return []byte("\x1b[B")
	case ebiten.KeyRight:
		return []byte("\x1b[C")
	case ebiten.KeyLeft:
		return []byte("\x1b[D")
	}
	return nil
}
