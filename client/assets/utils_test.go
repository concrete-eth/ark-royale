package assets

import (
	"image/color"
	"testing"
)

func TestChangeHSV(t *testing.T) {
	clr := color.RGBA{11, 22, 33, 255}
	clrH0 := ChangeHSV(clr, 0, 1, 1)

	if clrH0 != clr {
		t.Errorf("ChangeHSV(clr, 0, 1, 1) != %v, want %v", clrH0, clr)
	}
}
