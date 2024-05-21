package rts

import (
	"image"

	"github.com/concrete-eth/archetype/utils"
)

func DeltasToArea(position image.Point, area image.Rectangle) (int, int) {
	if position.In(area) {
		return 0, 0
	}
	var dx, dy int
	if position.X < area.Min.X {
		dx = area.Min.X - position.X
	} else if position.X >= area.Max.X {
		dx = area.Max.X - position.X - 1
	}
	if position.Y < area.Min.Y {
		dy = area.Min.Y - position.Y
	} else if position.Y >= area.Max.Y {
		dy = area.Max.Y - position.Y - 1
	}
	return dx, dy
}

func DeltasToAreaAsPoint(position image.Point, area image.Rectangle) image.Point {
	dx, dy := DeltasToArea(position, area)
	return image.Point{dx, dy}
}

func ChebyshevDistance(a, b image.Point) int {
	return max(utils.Abs(a.X-b.X), utils.Abs(a.Y-b.Y))
}

func ChebyshevDistanceToArea(position image.Point, area image.Rectangle) int {
	if position.In(area) {
		return 0
	}
	dx, dy := DeltasToArea(position, area)
	return max(utils.Abs(dx), utils.Abs(dy))
}

func ManhattanDistance(a, b image.Point) int {
	return utils.Abs(a.X-b.X) + utils.Abs(a.Y-b.Y)
}

func ManhattanDistanceToArea(position image.Point, area image.Rectangle) int {
	if position.In(area) {
		return 0
	}
	dx, dy := DeltasToArea(position, area)
	return utils.Abs(dx) + utils.Abs(dy)
}

func Distance(a image.Point, b image.Point) int {
	return ChebyshevDistance(a, b)
}

func DistanceToArea(position image.Point, area image.Rectangle) int {
	return ChebyshevDistanceToArea(position, area)
}

func StepTowards1D(curPos, targetL, targetSize int) int {
	if curPos < targetL {
		return 1
	} else if curPos >= targetL+targetSize {
		return -1
	} else {
		return 0
	}
}

func StepTowards2D(curPos image.Point, targetArea image.Rectangle) image.Point {
	return image.Point{
		X: StepTowards1D(curPos.X, targetArea.Min.X, targetArea.Size().X),
		Y: StepTowards1D(curPos.Y, targetArea.Min.Y, targetArea.Size().Y),
	}
}

func StepTowards(curPos image.Point, targetArea image.Rectangle) image.Point {
	return StepTowards2D(curPos, targetArea)
}

func GetPositionAsPoint(object interface {
	GetX() uint16
	GetY() uint16
}) image.Point {
	return image.Point{int(object.GetX()), int(object.GetY())}
}

func GetDimensionsAsPoint(object interface {
	GetWidth() uint8
	GetHeight() uint8
}) image.Point {
	return image.Point{int(object.GetWidth()), int(object.GetHeight())}
}

type IIterator_Uint8 interface {
	Next() bool
	Value() (id uint8, val interface{})
}

type Iterator_uint8 struct {
	Current uint8
	Max     uint8
	Get     func(idx uint8) (id uint8, val interface{})
	nextId  uint8
	nextVal interface{}
}

var _ IIterator_Uint8 = &Iterator_uint8{}

func (i *Iterator_uint8) Next() bool {
	for i.Current < i.Max {
		i.Current++
		id, val := i.Get(i.Current)
		if val != nil {
			i.nextId = id
			i.nextVal = val
			return true
		}
	}
	return false
}

func (i *Iterator_uint8) Value() (id uint8, val interface{}) {
	return i.nextId, i.nextVal
}
