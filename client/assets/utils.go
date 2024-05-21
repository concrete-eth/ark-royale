package assets

import (
	"embed"
	"fmt"
	"image"
	"image/color"

	"math"
	"path/filepath"

	"github.com/concrete-eth/archetype/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lucasb-eyer/go-colorful"
)

//go:embed all:sprites
var spritesFS embed.FS

func LoadImage(name string) *ebiten.Image {
	path := filepath.Join("sprites", name)
	img, _, err := ebitenutil.NewImageFromFileSystem(spritesFS, path)
	if err != nil {
		panic(err)
	}
	return img
}

func NewBounds(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}

func SubImage(img *ebiten.Image, rect image.Rectangle) *ebiten.Image {
	absBounds := rect.Add(img.Bounds().Min)
	return img.SubImage(absBounds).(*ebiten.Image)
}

func ScaleImage(img *ebiten.Image, x, y float64) *ebiten.Image {
	scaledImg := ebiten.NewImage(
		int(float64(img.Bounds().Dx())*x),
		int(float64(img.Bounds().Dy())*y),
	)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(x), float64(y))
	scaledImg.DrawImage(img, op)
	return scaledImg
}

func ScaleImageLightness(img *ebiten.Image, scale float64) *ebiten.Image {
	colorM := colorm.ColorM{}
	colorM.ChangeHSV(0, 1, scale)
	imgBounds := img.Bounds()
	newImg := ebiten.NewImage(imgBounds.Dx(), imgBounds.Dy())
	colorm.DrawImage(newImg, img, colorM, nil)
	return newImg
}

func ChangeImageHSV(img *ebiten.Image, hueTheta float64, saturationScale float64, valueScale float64) *ebiten.Image {
	colorM := colorm.ColorM{}
	colorM.ChangeHSV(hueTheta, saturationScale, valueScale)
	imgBounds := img.Bounds()
	newImg := ebiten.NewImage(imgBounds.Dx(), imgBounds.Dy())
	colorm.DrawImage(newImg, img, colorM, nil)
	return newImg
}

var stepToDirection = map[image.Point]Direction{
	image.Pt(0, -1):  Direction_Up,
	image.Pt(1, -1):  Direction_UpRight,
	image.Pt(1, 0):   Direction_Right,
	image.Pt(1, 1):   Direction_DownRight,
	image.Pt(0, 1):   Direction_Down,
	image.Pt(-1, 1):  Direction_DownLeft,
	image.Pt(-1, 0):  Direction_Left,
	image.Pt(-1, -1): Direction_UpLeft,
}

func StepToDirection(step image.Point) Direction {
	if _, ok := stepToDirection[step]; !ok {
		panic(fmt.Sprintf("assets: invalid step: %v", step))
	}
	return stepToDirection[step]
}

func DirectionFromDelta(delta image.Point) Direction {
	max := utils.Max(utils.Abs(delta.X), utils.Abs(delta.Y))
	step := delta
	if step.X != 0 {
		if max/step.X > 2 {
			step.X = 0
		} else {
			step.X /= utils.Abs(step.X)
		}
	}
	if step.Y != 0 {
		if max/step.Y > 2 {
			step.Y = 0
		} else {
			step.Y /= utils.Abs(step.Y)
		}
	}
	return StepToDirection(step)
}

func ChangeHSV(clr color.RGBA, hueTheta float64, saturationScale float64, valueScale float64) color.RGBA {
	base := colorful.Color{R: float64(clr.R) / 255, G: float64(clr.G) / 255, B: float64(clr.B) / 255}
	h, s, v := base.Hsv()
	hueTheta360 := math.Mod(hueTheta/(2*math.Pi)*360, 360)
	base = colorful.Hsv(
		h+hueTheta360,
		utils.Clamp(s*saturationScale, 0, 1),
		utils.Clamp(v*valueScale, 0, 1),
	)
	r, g, b := base.RGB255()
	return color.RGBA{r, g, b, clr.A}
}

func ScaleColor(clr color.RGBA, r, g, b float64) color.RGBA {
	return color.RGBA{
		R: uint8(utils.Clamp(float64(clr.R)*r, 0, 255)),
		G: uint8(utils.Clamp(float64(clr.G)*g, 0, 255)),
		B: uint8(utils.Clamp(float64(clr.B)*b, 0, 255)),
		A: clr.A,
	}
}
