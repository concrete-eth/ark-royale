package assets

import (
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const (
	Font_GoRegular = iota
	Font_PressStart
)

var (
	goRegularTTF, _  = opentype.Parse(goregular.TTF)
	pressStartTTF, _ = opentype.Parse(fonts.PressStart2P_ttf)

	fontsTTF = map[int]*sfnt.Font{
		Font_GoRegular:  goRegularTTF,
		Font_PressStart: pressStartTTF,
	}

	loadedFonts = map[int]map[float64]font.Face{}
)

func GetFontFace(fontId int, size float64) font.Face {
	if _, ok := loadedFonts[fontId]; !ok {
		loadedFonts[fontId] = map[float64]font.Face{}
	}
	if textFont, ok := loadedFonts[fontId][size]; ok {
		return textFont
	}
	fontTTF := fontsTTF[fontId]
	textFont, _ := opentype.NewFace(fontTTF, &opentype.FaceOptions{Size: size, DPI: 72})
	loadedFonts[fontId][size] = textFont
	return textFont
}
