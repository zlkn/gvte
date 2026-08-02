package main

import (
	"flag"
	"fmt"
)

func main() {
	debug := flag.Bool("debug", false, "set log level to verbose debug")
	flag.Parse()

	fmt.Println("debug:", *debug)
}
