package main

import (
	"flag"
	"log"

	"gvte/internal/config"
	"gvte/internal/emulator"
	"gvte/internal/ui"
	"gvte/internal/ui/font"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	debug := flag.Bool("debug", false, "Enable verbose debug logging")
	flag.Parse()

	config.Debug = *debug

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fontMgr, err := font.NewManager(cfg.Font.Family, cfg.Font.Size, cfg.Font.DPI)
	if err != nil {
		log.Fatalf("Failed to initialize font manager: %v", err)
	}

	// Compute initial grid dimensions from the real cell metrics
	cols := max(1, cfg.InitialWidth/fontMgr.CellWidth)
	rows := max(1, cfg.InitialHeight/fontMgr.CellHeight)
	state := emulator.NewState(cols, rows)

	win, err := ui.NewWindow(cfg, fontMgr, state)
	if err != nil {
		log.Fatalf("Failed to initialize window UI: %v", err)
	}

	if err := win.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
