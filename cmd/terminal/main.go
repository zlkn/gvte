package main

import (
	"flag"
	"log"

	"gvte/internal/config"
	"gvte/internal/emulator"
	"gvte/internal/ui"
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

	// Compute initial grid dimensions
	cols := cfg.InitialWidth / 9
	rows := cfg.InitialHeight / 18
	state := emulator.NewState(cols, rows)

	win, err := ui.NewWindow(cfg, state)
	if err != nil {
		log.Fatalf("Failed to initialize window UI: %v", err)
	}

	if err := win.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
