package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// InputMapper translates user keyboard and mouse events into ANSI/VT escape sequences.
type InputMapper struct{}

// NewMapper creates a new InputMapper.
func NewMapper() *InputMapper {
	return &InputMapper{}
}

func (im *InputMapper) MapKey() []byte {
	// NOTE: mocked
	return nil
}

func (im *InputMapper) HandleKey(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) {
	// NOTE: mocked
}
