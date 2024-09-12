package assets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/colorm"
)

var (
	LightGray = color.Gray{0xe0}

	TerrainBackgroundColor = color.RGBA{0x26, 0x1b, 0x23, 0xff}
	// GroundColor            = color.RGBA{0xff, 0xe1, 0xae, 0xff}
	GroundColor   = color.RGBA{0xfa, 0xde, 0xb1, 0xff}
	ResourceColor = color.RGBA{0xf9, 0xab, 0x8a, 0xff}
	ComputeColor  = color.RGBA{0x1d, 0xb5, 0x76, 0xff}

	LightShadowColor     = color.RGBA{0x9b, 0xab, 0xb2, 96}
	LightBlueShadowColor = color.RGBA{0x9b, 0xab, 0xb2, 96}
	DarkShadowColor      = color.RGBA{0x00, 0x00, 0x00, 32}
	DarkBlueShadowColor  = color.RGBA{0x3b, 0x33, 0x42, 96}

	TextLightColor   = color.RGBA{0xb6, 0xa8, 0xbf, 0xff}
	BoxTextDarkColor = color.RGBA{0x8e, 0x7b, 0x9e, 0xff}
)

func NewTerrainColorMatrix() colorm.ColorM {
	m := colorm.ColorM{}
	// m.ChangeHSV(0, 0.525, 1.475)
	return m
}

func NewBuildingColorMatrix() colorm.ColorM {
	m := colorm.ColorM{}
	// m.ChangeHSV(0, 1.075, 1)
	return m
}

func NewUnitColorMatrix() colorm.ColorM {
	m := colorm.ColorM{}
	// m.ChangeHSV(0, 1.25, 1)
	return m
}

var (
	AnticipatedColorMatrix  = colorm.ColorM{}
	UnpaidColorMatrix       = colorm.ColorM{}
	SpawningColorMatrix     = colorm.ColorM{}
	GhostColorMatrix        = colorm.ColorM{}
	NonBuildableColorMatrix = colorm.ColorM{}
)

func init() {
	AnticipatedColorMatrix.Scale(1, 1, 1, 0.65)      // Semi-transparent
	UnpaidColorMatrix.Scale(0.8, 1, 1.25, 0.75)      // Blue tint semi-transparent
	SpawningColorMatrix.ChangeHSV(0, 0.5, 0.75)      // De-saturated darkened
	GhostColorMatrix.Scale(1, 1, 1, 0.65)            // Semi-transparent
	NonBuildableColorMatrix.Scale(1.25, 0.8, 0.8, 1) // Red tint
}
