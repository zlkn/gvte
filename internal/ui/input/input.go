package input

import ()

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
