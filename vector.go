package scribble

// contents basically cloned from ebiten's vector package
// but only the StrokeLine parts, because their api doesnt let us set options
import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage     = ebiten.NewImage(3, 3)
	whiteSubImage  = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	cachedVertices []ebiten.Vertex
	cachedIndices  []uint16
)

func init() {
	b := whiteImage.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	for i := range pix {
		pix[i] = 0xff
	}
	whiteImage.WritePixels(pix)
}

// cloned from ebiten's `vector.StrokeLine` but with our own custom options because they dont allow them as a param
func StrokeLine(dst *ebiten.Image, x0, y0, x1, y1 float32, strokeWidth float32, clr color.Color, antialias bool) {
	var path vector.Path
	path.MoveTo(x0, y0)
	path.LineTo(x1, y1)
	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = strokeWidth
	strokeOp.LineCap = vector.LineCapRound

	cachedVertices, cachedIndices = path.AppendVerticesAndIndicesForStroke(cachedVertices[:0], cachedIndices[:0], strokeOp)
	drawVerticesForUtil(dst, cachedVertices, cachedIndices, clr, antialias)
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
