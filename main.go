package main

import (
	"cartermel/bruh/bruh"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const DEBUG = true
const CANVAS_SQ_SUZE = 2000
const BRUH_SIZE = 10

type Game struct {
	previousX, previousY float32
	mainCanvas           *ebiten.Image

	stateStack RedoStack[*ebiten.Image]
}

func NewGame() *Game {
	g := &Game{
		mainCanvas: ebiten.NewImage(CANVAS_SQ_SUZE, CANVAS_SQ_SUZE),
	}
	g.mainCanvas.Fill(color.Black)
	g.stateStack = RedoStack[*ebiten.Image]{}
	g.stateStack.Push(ebiten.NewImageFromImage(g.mainCanvas))

	return g
}

func CursorPositionF() (float32, float32) {
	x, y := ebiten.CursorPosition()
	return float32(x), float32(y)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		ebiten.SetCursorShape(ebiten.CursorShapeMove)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// push previous state just before drawing a new line
		g.stateStack.Push(ebiten.NewImageFromImage(g.mainCanvas))
		g.previousX, g.previousY = CursorPositionF()
	}

	// all below processing cannot be done if the mouse button is pressed (drawing)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return nil
	}

	// undo handling
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		if lastState, ok := g.stateStack.Pop(); ok {
			g.mainCanvas = lastState
		}
	}

	// redo handling
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		if redoState, ok := g.stateStack.Redo(); ok {
			g.mainCanvas = redoState
		}
	}

	// undo functionality, if we're pressing ctrl-z and not currently moused down (drawing)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := CursorPositionF()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		sameJoint := x == g.previousX && y == g.previousY
		if sameJoint {
			vector.StrokeCircle(g.mainCanvas, x, y, 1, BRUH_SIZE-2, color.White, true)
		} else {
			bruh.StrokeLine(g.mainCanvas, g.previousX, g.previousY, x, y, BRUH_SIZE, color.White, true)
		}

		g.previousX, g.previousY = CursorPositionF()
	}

	screen.DrawImage(g.mainCanvas, nil)

	if DEBUG {
		debug(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func debug(screen *ebiten.Image) {
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Paint (Ebitengine Demo)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
