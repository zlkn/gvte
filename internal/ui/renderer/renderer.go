package renderer

import (
	"gvte/internal/config"
	"gvte/internal/emulator"
	"gvte/internal/ui/font"

	"github.com/hajimehoshi/ebiten/v2"
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

// Render draws the current terminal state onto the given screen surface.
func (r *Renderer) Render(screen *ebiten.Image, st *emulator.State) {
	// Screen rendering stub (draw background, text grid cells, cursor, and selection)
}
