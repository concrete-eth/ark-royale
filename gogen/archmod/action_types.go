/* Autogenerated file. Do not edit manually. */

package archmod

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ = common.Big1
)

/*
Table                 KeySize  ValueSize
Initialize            0        4
Start                 0        0
CreateUnit            0        6
AssignUnit            0        19
PlaceBuilding         0        6
AddPlayer             0        10
AddUnitPrototype      0        14
AddBuildingPrototype  0        13
*/

type ActionData_Initialize struct {
	Width  uint16 `json:"width"`
	Height uint16 `json:"height"`
}

func (row *ActionData_Initialize) GetWidth() uint16 {
	return row.Width
}

func (row *ActionData_Initialize) GetHeight() uint16 {
	return row.Height
}

type ActionData_Start struct {
}

type ActionData_CreateUnit struct {
	PlayerId uint8  `json:"playerId"`
	UnitType uint8  `json:"unitType"`
	X        uint16 `json:"x"`
	Y        uint16 `json:"y"`
}

func (row *ActionData_CreateUnit) GetPlayerId() uint8 {
	return row.PlayerId
}

func (row *ActionData_CreateUnit) GetUnitType() uint8 {
	return row.UnitType
}

func (row *ActionData_CreateUnit) GetX() uint16 {
	return row.X
}

func (row *ActionData_CreateUnit) GetY() uint16 {
	return row.Y
}

type ActionData_AssignUnit struct {
	PlayerId     uint8  `json:"playerId"`
	UnitId       uint8  `json:"unitId"`
	Command      uint64 `json:"command"`
	CommandExtra uint64 `json:"commandExtra"`
	CommandMeta  uint8  `json:"commandMeta"`
}

func (row *ActionData_AssignUnit) GetPlayerId() uint8 {
	return row.PlayerId
}

func (row *ActionData_AssignUnit) GetUnitId() uint8 {
	return row.UnitId
}

func (row *ActionData_AssignUnit) GetCommand() uint64 {
	return row.Command
}

func (row *ActionData_AssignUnit) GetCommandExtra() uint64 {
	return row.CommandExtra
}

func (row *ActionData_AssignUnit) GetCommandMeta() uint8 {
	return row.CommandMeta
}

type ActionData_PlaceBuilding struct {
	PlayerId     uint8  `json:"playerId"`
	BuildingType uint8  `json:"buildingType"`
	X            uint16 `json:"x"`
	Y            uint16 `json:"y"`
}

func (row *ActionData_PlaceBuilding) GetPlayerId() uint8 {
	return row.PlayerId
}

func (row *ActionData_PlaceBuilding) GetBuildingType() uint8 {
	return row.BuildingType
}

func (row *ActionData_PlaceBuilding) GetX() uint16 {
	return row.X
}

func (row *ActionData_PlaceBuilding) GetY() uint16 {
	return row.Y
}

type ActionData_AddPlayer struct {
	SpawnAreaX      uint16 `json:"spawnAreaX"`
	SpawnAreaY      uint16 `json:"spawnAreaY"`
	SpawnAreaWidth  uint8  `json:"spawnAreaWidth"`
	SpawnAreaHeight uint8  `json:"spawnAreaHeight"`
	WorkerPortX     uint16 `json:"workerPortX"`
	WorkerPortY     uint16 `json:"workerPortY"`
}

func (row *ActionData_AddPlayer) GetSpawnAreaX() uint16 {
	return row.SpawnAreaX
}

func (row *ActionData_AddPlayer) GetSpawnAreaY() uint16 {
	return row.SpawnAreaY
}

func (row *ActionData_AddPlayer) GetSpawnAreaWidth() uint8 {
	return row.SpawnAreaWidth
}

func (row *ActionData_AddPlayer) GetSpawnAreaHeight() uint8 {
	return row.SpawnAreaHeight
}

func (row *ActionData_AddPlayer) GetWorkerPortX() uint16 {
	return row.WorkerPortX
}

func (row *ActionData_AddPlayer) GetWorkerPortY() uint16 {
	return row.WorkerPortY
}

type ActionData_AddUnitPrototype struct {
	Layer             uint8  `json:"layer"`
	ResourceCost      uint16 `json:"resourceCost"`
	ComputeCost       uint8  `json:"computeCost"`
	SpawnTime         uint8  `json:"spawnTime"`
	MaxIntegrity      uint8  `json:"maxIntegrity"`
	LandStrength      uint8  `json:"landStrength"`
	HoverStrength     uint8  `json:"hoverStrength"`
	AirStrength       uint8  `json:"airStrength"`
	AttackRange       uint8  `json:"attackRange"`
	AttackCooldown    uint8  `json:"attackCooldown"`
	IsAssault         bool   `json:"isAssault"`
	IsConfrontational bool   `json:"isConfrontational"`
	IsWorker          bool   `json:"isWorker"`
}

func (row *ActionData_AddUnitPrototype) GetLayer() uint8 {
	return row.Layer
}

func (row *ActionData_AddUnitPrototype) GetResourceCost() uint16 {
	return row.ResourceCost
}

func (row *ActionData_AddUnitPrototype) GetComputeCost() uint8 {
	return row.ComputeCost
}

func (row *ActionData_AddUnitPrototype) GetSpawnTime() uint8 {
	return row.SpawnTime
}

func (row *ActionData_AddUnitPrototype) GetMaxIntegrity() uint8 {
	return row.MaxIntegrity
}

func (row *ActionData_AddUnitPrototype) GetLandStrength() uint8 {
	return row.LandStrength
}

func (row *ActionData_AddUnitPrototype) GetHoverStrength() uint8 {
	return row.HoverStrength
}

func (row *ActionData_AddUnitPrototype) GetAirStrength() uint8 {
	return row.AirStrength
}

func (row *ActionData_AddUnitPrototype) GetAttackRange() uint8 {
	return row.AttackRange
}

func (row *ActionData_AddUnitPrototype) GetAttackCooldown() uint8 {
	return row.AttackCooldown
}

func (row *ActionData_AddUnitPrototype) GetIsAssault() bool {
	return row.IsAssault
}

func (row *ActionData_AddUnitPrototype) GetIsConfrontational() bool {
	return row.IsConfrontational
}

func (row *ActionData_AddUnitPrototype) GetIsWorker() bool {
	return row.IsWorker
}

type ActionData_AddBuildingPrototype struct {
	Width            uint8  `json:"width"`
	Height           uint8  `json:"height"`
	ResourceCost     uint16 `json:"resourceCost"`
	ResourceCapacity uint16 `json:"resourceCapacity"`
	ComputeCapacity  uint8  `json:"computeCapacity"`
	ResourceMine     uint8  `json:"resourceMine"`
	MineTime         uint8  `json:"mineTime"`
	MaxIntegrity     uint8  `json:"maxIntegrity"`
	BuildingTime     uint8  `json:"buildingTime"`
	IsArmory         bool   `json:"isArmory"`
	IsEnvironment    bool   `json:"isEnvironment"`
}

func (row *ActionData_AddBuildingPrototype) GetWidth() uint8 {
	return row.Width
}

func (row *ActionData_AddBuildingPrototype) GetHeight() uint8 {
	return row.Height
}

func (row *ActionData_AddBuildingPrototype) GetResourceCost() uint16 {
	return row.ResourceCost
}

func (row *ActionData_AddBuildingPrototype) GetResourceCapacity() uint16 {
	return row.ResourceCapacity
}

func (row *ActionData_AddBuildingPrototype) GetComputeCapacity() uint8 {
	return row.ComputeCapacity
}

func (row *ActionData_AddBuildingPrototype) GetResourceMine() uint8 {
	return row.ResourceMine
}

func (row *ActionData_AddBuildingPrototype) GetMineTime() uint8 {
	return row.MineTime
}

func (row *ActionData_AddBuildingPrototype) GetMaxIntegrity() uint8 {
	return row.MaxIntegrity
}

func (row *ActionData_AddBuildingPrototype) GetBuildingTime() uint8 {
	return row.BuildingTime
}

func (row *ActionData_AddBuildingPrototype) GetIsArmory() bool {
	return row.IsArmory
}

func (row *ActionData_AddBuildingPrototype) GetIsEnvironment() bool {
	return row.IsEnvironment
}
