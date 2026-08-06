package renderer

import (
	"gvte/internal/config"
	"gvte/internal/emulator"
	"gvte/internal/ui/font"
)

type Renderer struct {
	FontMgr *font.FontManager
	Config  *config.Config
}

func New(fm *font.FontManager, cfg *config.Config) *Renderer {
	return &Renderer{
		FontMgr: fm,
		Config:  cfg,
	}
}

// Render draws the current terminal state.
func (r *Renderer) Render(st *emulator.State) {
	// Screen rendering stub (draw background, text grid cells, cursor, and selection)
}
