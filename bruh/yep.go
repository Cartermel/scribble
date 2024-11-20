// TODO: clean this up, cloned straight from vector.StrokeLine() or whatever
// just so i could set the options for the stroke which they dont expose...

package bruh

import (
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

var (
	cachedVertices []ebiten.Vertex
	cachedIndices  []uint16
	cacheM         sync.Mutex
)

func init() {
	b := whiteImage.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	for i := range pix {
		pix[i] = 0xff
	}
	// This is hacky, but WritePixels is better than Fill in term of automatic texture packing.
	whiteImage.WritePixels(pix)
}

func useCachedVerticesAndIndices(fn func([]ebiten.Vertex, []uint16) (vs []ebiten.Vertex, is []uint16)) {
	cacheM.Lock()
	defer cacheM.Unlock()
	cachedVertices, cachedIndices = fn(cachedVertices[:0], cachedIndices[:0])
}

// StrokeLine strokes a line (x0, y0)-(x1, y1) with the specified width and color.
// clr has be to be a solid (non-transparent) color.
func StrokeLine(dst *ebiten.Image, x0, y0, x1, y1 float32, strokeWidth float32, clr color.Color, antialias bool) {
	var path vector.Path
	path.MoveTo(x0, y0)
	path.LineTo(x1, y1)
	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = strokeWidth
	strokeOp.LineCap = vector.LineCapRound

	useCachedVerticesAndIndices(func(vs []ebiten.Vertex, is []uint16) ([]ebiten.Vertex, []uint16) {
		vs, is = path.AppendVerticesAndIndicesForStroke(vs, is, strokeOp)
		drawVerticesForUtil(dst, vs, is, clr, antialias)
		return vs, is
	})
}

func drawVerticesForUtil(dst *ebiten.Image, vs []ebiten.Vertex, is []uint16, clr color.Color, antialias bool) {
	r, g, b, a := clr.RGBA()
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha
	op.AntiAlias = antialias
	dst.DrawTriangles(vs, is, whiteSubImage, op)
}
