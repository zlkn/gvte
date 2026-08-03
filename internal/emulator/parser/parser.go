package parser

// ActionType represents an ANSI / VT100 / xterm parse action.
type ActionType int

const (
	ActionPrint ActionType = iota
	ActionExecute
	ActionCSIParam
	ActionCSIDispatch
	ActionOSCDispatch
	ActionEscDispatch
)

// Sequence represents a parsed escape sequence or plain character.
type Sequence struct {
	Type        ActionType
	Char        rune
	Params      []int
	Intermediate []byte
}

// Parser implements an ANSI / VT100 / xterm escape sequence parser state machine.
type Parser struct {
	state int
}

// New creates a new ANSI sequence parser.
func New() *Parser {
	return &Parser{}
}

// Parse processes incoming bytes and returns parsed sequence actions.
func (p *Parser) Parse(buf []byte, handler func(seq Sequence)) {
	for _, b := range buf {
		// Basic stub handling plain text byte output
		handler(Sequence{
			Type: ActionPrint,
			Char: rune(b),
		})
	}
}
