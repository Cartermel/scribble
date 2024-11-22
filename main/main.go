package main

import (
	"cartermel/scribble"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const VERSION = "1.0.0"

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			fmt.Println("scribble version", VERSION)
			os.Exit(0)
		}
	}

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("scribble")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(scribble.NewGame()); err != nil {
		log.Fatal(err)
	}
}
