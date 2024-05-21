package game

import (
	"image"

	"github.com/concrete-eth/ark-rts/rts"
)

const (
	BuildingPrototypeId_Main uint8 = iota + 1
	BuildingPrototypeId_Storage
	BuildingPrototypeId_Lab
	BuildingPrototypeId_Armory
	BuildingPrototypeId_Mine
)

const (
	UnitPrototypeId_Worker uint8 = iota + 1
	UnitPrototypeId_Air
	UnitPrototypeId_AntiAir
	UnitPrototypeId_Tank
)

const (
	BuildRadius = 4
	MaxTicks    = 1800
)

var (
	BuildableBuildingPrototypeIds = []uint8{BuildingPrototypeId_Storage, BuildingPrototypeId_Lab, BuildingPrototypeId_Armory}
	UnitPrototypeIds              = []uint8{UnitPrototypeId_Worker, UnitPrototypeId_Air, UnitPrototypeId_AntiAir, UnitPrototypeId_Tank}
)

func NeedsArmory(unitPrototypeId uint8) bool {
	return unitPrototypeId != UnitPrototypeId_Worker
}

func GetPlayerBuildableArea(core *rts.Core, playerId uint8) image.Rectangle {
	mainBuilding := core.GetMainBuilding(playerId)
	mainBuildingPosition := rts.GetPositionAsPoint(mainBuilding)
	return image.Rectangle{
		Min: mainBuildingPosition.Sub(image.Pt(BuildRadius, BuildRadius)),
		Max: mainBuildingPosition.Add(image.Pt(BuildRadius, BuildRadius)).Add(image.Pt(2, 2)),
	}
}

func IsInPlayerBuildableArea(core *rts.Core, playerId uint8, area image.Rectangle) bool {
	buildArea := GetPlayerBuildableArea(core, playerId)
	return area.In(buildArea)
}
