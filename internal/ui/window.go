package ui

import (
	"log"

	"gvte/internal/config"
	"gvte/internal/emulator"
	"gvte/internal/ui/font"
	"gvte/internal/ui/input"
	"gvte/internal/ui/renderer"

	"github.com/hajimehoshi/ebiten/v2"
)

// AppWindow manages the main application GUI window, loop, rendering, and input dispatching.
type AppWindow struct {
	Config   *config.Config
	State    *emulator.State
	Renderer *renderer.Renderer
	Input    *input.InputMapper
	FontMgr  *font.FontManager
}

// NewWindow initializes the UI window component.
func NewWindow(cfg *config.Config, st *emulator.State) (*AppWindow, error) {
	fm, err := font.NewManager(cfg.Font.Family, cfg.Font.Size)
	if err != nil {
		return nil, err
	}

	rend := renderer.New(fm, cfg)
	inp := input.NewMapper()

	return &AppWindow{
		Config:   cfg,
		State:    st,
		Renderer: rend,
		Input:    inp,
		FontMgr:  fm,
	}, nil
}

// Update handles frame logic and input.
func (w *AppWindow) Update() error {
	// Handle window resize or input processing
	return nil
}

// Draw renders frame content.
func (w *AppWindow) Draw(screen *ebiten.Image) {
	w.Renderer.Render(screen, w.State)
}

// Layout calculates screen dimensions for logical resolution.
func (w *AppWindow) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// Run starts the graphical application event loop.
func (w *AppWindow) Run() error {
	ebiten.SetWindowSize(w.Config.InitialWidth, w.Config.InitialHeight)
	ebiten.SetWindowTitle("Terminal")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetScreenClearedEveryFrame(true)

	log.Println("Starting terminal window UI loop...")
	return ebiten.RunGame(w)
}
