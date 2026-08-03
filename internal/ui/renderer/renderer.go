package renderer

import (
	"gvte/internal/config"
	"gvte/internal/emulator"
	"gvte/internal/ui/font"

	"github.com/hajimehoshi/ebiten/v2"
)

// Renderer draws terminal grid cells and decorations onto screen destination.
type Renderer struct {
	FontMgr *font.FontManager
	Config  *config.Config
}

// New creates a new UI Renderer instance.
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
