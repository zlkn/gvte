package window

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	initialWidth  = 900
	initialHeight = 600
)

// game is the Ebitengine adapter. It reconciles one pane (shell) per tab, routes
// input, and draws the active tab.
type game struct {
	window *Window
}

func Run() error {
	log.Println("Mock run")

	g := &game{
		window: New(),
	}

	ebiten.SetWindowSize(initialWidth, initialHeight)
	ebiten.SetWindowTitle("gvte")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetScreenClearedEveryFrame(false)
	return ebiten.RunGame(g)
}
