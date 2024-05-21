package utils

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// Returns the current cursor position as a Point.
func CursorPosition() image.Point {
	x, y := ebiten.CursorPosition()
	return image.Point{x, y}
}

// Creates a new DrawImageOptions that scales the source image to fit the destination rectangle.
func NewDrawOptions(dest, src image.Rectangle) *colorm.DrawImageOptions {
	op := &colorm.DrawImageOptions{}
	op.GeoM.Scale(
		float64(dest.Dx())/float64(src.Dx()),
		float64(dest.Dy())/float64(src.Dy()),
	)
	op.GeoM.Translate(float64(dest.Min.X), float64(dest.Min.Y))
	return op
}

// Scales a rectangle by a given ratio pair.
func ScaleRectangle(rect image.Rectangle, x, y float64) image.Rectangle {
	return image.Rect(
		int(float64(rect.Min.X)*x),
		int(float64(rect.Min.Y)*y),
		int(float64(rect.Max.X)*x),
		int(float64(rect.Max.Y)*y),
	)
}

// Returns the ratio of the width and height of two rectangles.
func ScaleRatio(a image.Rectangle, b image.Rectangle) (float64, float64) {
	return float64(a.Dx()) / float64(b.Bounds().Dx()), float64(a.Dy()) / float64(b.Bounds().Dy())
}
