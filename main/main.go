package main

import (
	"cartermel/scribble"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("scribble")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(scribble.NewGame()); err != nil {
		log.Fatal(err)
	}
}
