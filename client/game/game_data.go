package game

const (
	BuildingPrototypeId_Main uint8 = iota + 1
	BuildingPrototypeId_Pit
	BuildingPrototypeId_Mine
)

const (
	UnitPrototypeId_Air uint8 = iota + 1
	UnitPrototypeId_AntiAir
	UnitPrototypeId_Tank
	UnitPrototypeId_Turret
	UnitPrototypeId_Worker
)

var (
	BuildableBuildingPrototypeIds = []uint8{}
	UnitPrototypeIds              = []uint8{UnitPrototypeId_Air, UnitPrototypeId_AntiAir, UnitPrototypeId_Tank}
)

const (
	MaxTicks = 1800
)
