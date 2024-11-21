package main

import (
	"cartermel/bruh/bruh"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const DEBUG = true
const CANVAS_SQ_SUZE = 2000

var (
	BRUSH_SIZE      float32 = 10
	MIN_BRUSH_SIZE  float32 = 4
	MAX_BRUSH_SIZE  float32 = 100
	BRUSH_INCREMENT float32 = 5
)

type Game struct {
	previousX, previousY float32
	mainCanvas           *ebiten.Image

	stateStack *StateStack
}

func NewGame() *Game {
	g := &Game{
		mainCanvas: ebiten.NewImage(CANVAS_SQ_SUZE, CANVAS_SQ_SUZE),
	}
	g.mainCanvas.Fill(color.Black)
	g.stateStack = &StateStack{}

	// sync updates with fps, since we dont really do cmoplex logic in the update
	// and we need smooth drawing (mouse updates are calculated from TPS)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	return g
}

func (g *Game) cursorPositionF() (float32, float32) {
	x, y := ebiten.CursorPosition()
	return float32(x), float32(y)
}

func (g *Game) handleMouseWheel() {
	_, dy := ebiten.Wheel()
	if dy > 0.1 {
		BRUSH_SIZE = min(BRUSH_SIZE+BRUSH_INCREMENT, MAX_BRUSH_SIZE)
	}
	if dy < -0.1 {
		BRUSH_SIZE = max(BRUSH_SIZE-BRUSH_INCREMENT, MIN_BRUSH_SIZE)
	}
}

var mouseDragAnchor = image.Point{}
var mouseDragDelta = image.Point{} // for use WHILE dragging
var canvasOffset = image.Point{}   // to apply after delta has been applied

func (g *Game) Update() error {
	g.handleMouseWheel()
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		ebiten.SetCursorShape(ebiten.CursorShapeMove)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mouseDragAnchor = image.Pt(ebiten.CursorPosition())
			mouseDragDelta = image.Pt(0, 0)
		}

		// on drag end, TODO: what if the user lets go of space??? or some other combination where this doest get called...
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			canvasOffset = canvasOffset.Add(mouseDragDelta)
			mouseDragDelta = image.Pt(0, 0)
		}

		// space, and mouse pressed
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			mouseDragDelta.X = mouseDragAnchor.X - x
			mouseDragDelta.Y = mouseDragAnchor.Y - y
		}
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// push previous state just before drawing a new line
		g.stateStack.Push(ebiten.NewImageFromImage(g.mainCanvas))
		g.previousX, g.previousY = g.cursorPositionF()
	}

	// all below processing cannot be done if the mouse button is pressed (drawing)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return nil
	}

	// redo and undo handling, check redo first, then undo. can only do one per update
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		if redoState, ok := g.stateStack.Redo(); ok {
			g.mainCanvas = redoState
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		// push current state before undoing, with false as a param so it doesnt instantly get poppeds
		if lastState, ok := g.stateStack.Undo(g.mainCanvas); ok {
			g.mainCanvas = lastState
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	offset := canvasOffset.Add(mouseDragDelta)

	x, y := g.cursorPositionF()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// check if the cursor moved or not
		stationaryTick := x == g.previousX && y == g.previousY
		if stationaryTick {
			vector.DrawFilledCircle(
				g.mainCanvas,
				x+float32(offset.X),
				y+float32(offset.Y),
				BRUSH_SIZE/2,
				color.White,
				true,
			)
		} else {
			bruh.StrokeLine(
				g.mainCanvas,
				g.previousX+float32(offset.X),
				g.previousY+float32(offset.Y),
				x+float32(offset.X),
				y+float32(offset.Y),
				BRUSH_SIZE,
				color.White,
				true,
			)
		}

		g.previousX, g.previousY = g.cursorPositionF()
	}

	geo := ebiten.GeoM{}
	geo.Translate(-float64(offset.X), -float64(offset.Y))
	screen.DrawImage(g.mainCanvas, &ebiten.DrawImageOptions{
		GeoM: geo,
	})

	// cursor
	vector.StrokeCircle(screen, x, y, BRUSH_SIZE/2+1, 2, color.White, true)

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
