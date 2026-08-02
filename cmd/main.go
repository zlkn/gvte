package main

import (
	"flag"
	"fmt"
	"gvte/internal/config"
)

func main() {
	debug := flag.Bool("debug", false, "set log level to verbose debug")
	flag.Parse()

	config.Debug = *debug

	fmt.Println("debug:", config.Debug)
}
