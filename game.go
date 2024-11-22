package scribble

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	DEBUG           bool    = false // TODO: compile flag for this
	CANVAS_SQ_SIZE  int     = 2000
	MIN_BRUSH_SIZE  float32 = 4
	MAX_BRUSH_SIZE  float32 = 100
	BRUSH_INCREMENT float32 = 5
)

type Game struct {
	previousX, previousY float32
	mainCanvas           *ebiten.Image

	stateStack *StateStack

	brushSize  float32
	brushColor color.Color

	mouseDragAnchor image.Point
	mouseDragDelta  image.Point // for use WHILE dragging
	canvasOffset    image.Point // to apply after delta has been applied

	isDrawing bool

	lastTickTime time.Time

	hasInitialized bool
}

func NewGame() *Game {
	g := &Game{
		mainCanvas:   ebiten.NewImage(CANVAS_SQ_SIZE, CANVAS_SQ_SIZE),
		brushColor:   color.White,
		lastTickTime: time.Now(),
		brushSize:    10,
	}
	g.mainCanvas.Fill(color.Black)
	g.stateStack = &StateStack{}

	// sync updates with fps, since we dont really do cmoplex logic in the update
	// and we need smooth drawing (mouse updates are calculated from TPS)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	return g
}

func (g *Game) cursorPositionF() (float32, float32) {
	x, y := ebiten.CursorPosition()
	return float32(x), float32(y)
}

func (g *Game) handleMouseWheel() {
	_, dy := ebiten.Wheel()
	if dy > 0.1 {
		g.brushSize = min(g.brushSize+BRUSH_INCREMENT, MAX_BRUSH_SIZE)
	}
	if dy < -0.1 {
		g.brushSize = max(g.brushSize-BRUSH_INCREMENT, MIN_BRUSH_SIZE)
	}
}

// retrieves the time since last tick, which is also set here.
// intended to be called from Update()
func (g *Game) DeltaTime() time.Duration {
	deltaTime := time.Since(g.lastTickTime)
	g.lastTickTime = time.Now()
	return deltaTime // TODO: use this for holding keybinds (like holding ctrl z)
}

func (g *Game) Update() error {
	g.handleMouseWheel()
	pressedKeybind := HandleKeyBindRead()
	anyKeybindsPressed := pressedKeybind != KeybindNone && !g.isDrawing

	// if no keybinds are being pressed, proceed with drawing logic
	if !anyKeybindsPressed {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// push previous state just before drawing a new line
			g.stateStack.Push(ebiten.NewImageFromImage(g.mainCanvas))
			g.previousX, g.previousY = g.cursorPositionF()
		}

		g.isDrawing = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	}

	if g.isDrawing {
		return nil
	}

	// ------------------------------------------------- //
	// --------- NON "DRAWING" MODE CODE BELOW --------- //
	// ------------------------------------------------- //

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.mouseDragAnchor = image.Pt(ebiten.CursorPosition())
			g.mouseDragDelta = image.Pt(0, 0)
		}

		// on drag end, TODO: what if the user lets go of space??? or some other combination where this doest get called...
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			g.canvasOffset = g.canvasOffset.Add(g.mouseDragDelta)
			g.mouseDragDelta = image.Pt(0, 0)
		}

		// space, and mouse pressed
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			g.mouseDragDelta.X = g.mouseDragAnchor.X - x
			g.mouseDragDelta.Y = g.mouseDragAnchor.Y - y
		}
	}

	// redo and undo handling, check redo first, then undo. can only do one per update
	if pressedKeybind == KeybindRedo && inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		if redoState, ok := g.stateStack.Redo(); ok {
			g.mainCanvas = redoState
		}
	} else if pressedKeybind == KeybindUndo && inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		// push current state before undoing, with false as a param so it doesnt instantly get poppeds
		if lastState, ok := g.stateStack.Undo(g.mainCanvas); ok {
			g.mainCanvas = lastState
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 20, 255})
	offset := g.canvasOffset.Add(g.mouseDragDelta)

	x, y := g.cursorPositionF()

	if g.isDrawing {
		// check if the cursor moved or not
		stationaryTick := x == g.previousX && y == g.previousY
		if stationaryTick {
			vector.DrawFilledCircle(
				g.mainCanvas,
				x+float32(offset.X),
				y+float32(offset.Y),
				g.brushSize/2,
				g.brushColor,
				true,
			)
		} else {
			StrokeLine(
				g.mainCanvas,
				g.previousX+float32(offset.X),
				g.previousY+float32(offset.Y),
				x+float32(offset.X),
				y+float32(offset.Y),
				g.brushSize,
				g.brushColor,
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
	vector.StrokeCircle(screen, x, y, g.brushSize/2+1, 3, color.Black, true)
	vector.StrokeCircle(screen, x, y, g.brushSize/2+1, 2, color.White, true)

	if DEBUG {
		debug(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// run init logic, center canvas etc
	if !g.hasInitialized {
		g.hasInitialized = true

		// center the canvas in the screen
		centerCanvas := g.mainCanvas.Bounds().Max.Div(2)
		centerScreen := image.Pt(outsideWidth/2, outsideHeight/2)
		diff := centerCanvas.Sub(centerScreen)
		g.canvasOffset = g.canvasOffset.Add(diff)
	}

	return outsideWidth, outsideHeight
}

func debug(screen *ebiten.Image) {
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}
