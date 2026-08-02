package main

import (
	"flag"
	"fmt"
	"log"

	"gvte/internal/config"
	"gvte/internal/window"
)

func main() {
	debug := flag.Bool("debug", false, "set log level to verbose debug")
	flag.Parse()

	config.Debug = *debug
	log.Println("debug:", config.Debug)

	window.Run()

	if err := window.Run(); err != nil {
		log.Fatal(err)
	}

}
