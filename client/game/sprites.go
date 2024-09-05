package game

import (
	"fmt"
	"image"

	"github.com/concrete-eth/ark-rts/client/assets"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/hajimehoshi/ebiten/v2"
)

var _ assets.SpriteGetter = (*spriteGetter)(nil)

var SpriteGetter = &spriteGetter{}

type spriteGetter struct {
	assets.DefaultSpriteGetter
}

var (
	buildingProtoIdToSpriteId = map[uint8]uint8{
		BuildingPrototypeId_Main: assets.BuildingSpriteId_Main,
		BuildingPrototypeId_Mine: assets.BuildingSpriteId_Mine,
	}
	unitProtoIdToSpriteId = map[uint8]uint8{
		UnitPrototypeId_AntiAir: assets.UnitSpriteId_AntiAir,
		UnitPrototypeId_Air:     assets.UnitSpriteId_Air,
		UnitPrototypeId_Tank:    assets.UnitSpriteId_Tank,
		UnitPrototypeId_Turret:  assets.UnitSpriteId_Turret,
		UnitPrototypeId_Worker:  assets.UnitSpriteId_Worker,
	}
)

var (
	emptyImage = ebiten.NewImage(1, 1)
)

func (d *spriteGetter) GetBuildingSprite(playerId uint8, buildingId uint8, buildingState rts.BuildingState) *ebiten.Image {
	if buildingId == BuildingPrototypeId_Pit {
		return emptyImage
	}
	if spriteId, ok := buildingProtoIdToSpriteId[buildingId]; ok {
		return d.DefaultSpriteGetter.GetBuildingSprite(playerId, spriteId, buildingState)
	}
	panic(fmt.Sprintf("invalid building id %v", buildingId))
}

func (d *spriteGetter) GetBuildingSpriteOrigin(buildingId uint8) image.Point {
	if buildingId == BuildingPrototypeId_Pit {
		return image.Point{0, 0}
	}
	if spriteId, ok := buildingProtoIdToSpriteId[buildingId]; ok {
		return d.DefaultSpriteGetter.GetBuildingSpriteOrigin(spriteId)
	}
	panic(fmt.Sprintf("invalid building id %v", buildingId))
}

func (d *spriteGetter) GetUnitSprite(playerId uint8, spriteId uint8, direction assets.Direction) *ebiten.Image {
	if spriteId, ok := unitProtoIdToSpriteId[spriteId]; ok {
		return d.DefaultSpriteGetter.GetUnitSprite(playerId, spriteId, direction)
	}
	panic(fmt.Sprintf("invalid unit id %v", spriteId))
}

func (d *spriteGetter) GetUnitFireFrame(playerId uint8, spriteId uint8, direction assets.Direction, frame uint) *ebiten.Image {
	if spriteId, ok := unitProtoIdToSpriteId[spriteId]; ok {
		return d.DefaultSpriteGetter.GetUnitFireFrame(playerId, spriteId, direction, frame)
	}
	panic(fmt.Sprintf("invalid unit id %v", spriteId))
}
