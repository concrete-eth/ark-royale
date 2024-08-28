/* Autogenerated file. Do not edit manually. */

package archmod

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ = common.Big1
)

/*
Table               KeySize  ValueSize
Meta                0        13
Players             1        22
Board               4        7
Units               2        30
Buildings           2        11
UnitPrototypes      1        13
BuildingPrototypes  1        13
*/

type RowData_Meta struct {
	BoardWidth             uint16 `json:"boardWidth"`
	BoardHeight            uint16 `json:"boardHeight"`
	PlayerCount            uint8  `json:"playerCount"`
	UnitPrototypeCount     uint8  `json:"unitPrototypeCount"`
	BuildingPrototypeCount uint8  `json:"buildingPrototypeCount"`
	IsInitialized          bool   `json:"isInitialized"`
	HasStarted             bool   `json:"hasStarted"`
	CreationBlockNumber    uint32 `json:"creationBlockNumber"`
}

func (row *RowData_Meta) GetBoardWidth() uint16 {
	return row.BoardWidth
}

func (row *RowData_Meta) GetBoardHeight() uint16 {
	return row.BoardHeight
}

func (row *RowData_Meta) GetPlayerCount() uint8 {
	return row.PlayerCount
}

func (row *RowData_Meta) GetUnitPrototypeCount() uint8 {
	return row.UnitPrototypeCount
}

func (row *RowData_Meta) GetBuildingPrototypeCount() uint8 {
	return row.BuildingPrototypeCount
}

func (row *RowData_Meta) GetIsInitialized() bool {
	return row.IsInitialized
}

func (row *RowData_Meta) GetHasStarted() bool {
	return row.HasStarted
}

func (row *RowData_Meta) GetCreationBlockNumber() uint32 {
	return row.CreationBlockNumber
}

type RowData_Players struct {
	SpawnAreaX                uint16 `json:"spawnAreaX"`
	SpawnAreaY                uint16 `json:"spawnAreaY"`
	SpawnAreaWidth            uint8  `json:"spawnAreaWidth"`
	SpawnAreaHeight           uint8  `json:"spawnAreaHeight"`
	WorkerPortX               uint16 `json:"workerPortX"`
	WorkerPortY               uint16 `json:"workerPortY"`
	CurResource               uint16 `json:"curResource"`
	MaxResource               uint16 `json:"maxResource"`
	CurArmories               uint8  `json:"curArmories"`
	ComputeSupply             uint8  `json:"computeSupply"`
	ComputeDemand             uint8  `json:"computeDemand"`
	UnitCount                 uint8  `json:"unitCount"`
	BuildingCount             uint8  `json:"buildingCount"`
	BuildingPayQueuePointer   uint8  `json:"buildingPayQueuePointer"`
	BuildingBuildQueuePointer uint8  `json:"buildingBuildQueuePointer"`
	UnitPayQueuePointer       uint8  `json:"unitPayQueuePointer"`
}

func (row *RowData_Players) GetSpawnAreaX() uint16 {
	return row.SpawnAreaX
}

func (row *RowData_Players) GetSpawnAreaY() uint16 {
	return row.SpawnAreaY
}

func (row *RowData_Players) GetSpawnAreaWidth() uint8 {
	return row.SpawnAreaWidth
}

func (row *RowData_Players) GetSpawnAreaHeight() uint8 {
	return row.SpawnAreaHeight
}

func (row *RowData_Players) GetWorkerPortX() uint16 {
	return row.WorkerPortX
}

func (row *RowData_Players) GetWorkerPortY() uint16 {
	return row.WorkerPortY
}

func (row *RowData_Players) GetCurResource() uint16 {
	return row.CurResource
}

func (row *RowData_Players) GetMaxResource() uint16 {
	return row.MaxResource
}

func (row *RowData_Players) GetCurArmories() uint8 {
	return row.CurArmories
}

func (row *RowData_Players) GetComputeSupply() uint8 {
	return row.ComputeSupply
}

func (row *RowData_Players) GetComputeDemand() uint8 {
	return row.ComputeDemand
}

func (row *RowData_Players) GetUnitCount() uint8 {
	return row.UnitCount
}

func (row *RowData_Players) GetBuildingCount() uint8 {
	return row.BuildingCount
}

func (row *RowData_Players) GetBuildingPayQueuePointer() uint8 {
	return row.BuildingPayQueuePointer
}

func (row *RowData_Players) GetBuildingBuildQueuePointer() uint8 {
	return row.BuildingBuildQueuePointer
}

func (row *RowData_Players) GetUnitPayQueuePointer() uint8 {
	return row.UnitPayQueuePointer
}

type RowData_Board struct {
	LandObjectType uint8 `json:"landObjectType"`
	LandPlayerId   uint8 `json:"landPlayerId"`
	LandObjectId   uint8 `json:"landObjectId"`
	HoverPlayerId  uint8 `json:"hoverPlayerId"`
	HoverUnitId    uint8 `json:"hoverUnitId"`
	AirPlayerId    uint8 `json:"airPlayerId"`
	AirUnitId      uint8 `json:"airUnitId"`
}

func (row *RowData_Board) GetLandObjectType() uint8 {
	return row.LandObjectType
}

func (row *RowData_Board) GetLandPlayerId() uint8 {
	return row.LandPlayerId
}

func (row *RowData_Board) GetLandObjectId() uint8 {
	return row.LandObjectId
}

func (row *RowData_Board) GetHoverPlayerId() uint8 {
	return row.HoverPlayerId
}

func (row *RowData_Board) GetHoverUnitId() uint8 {
	return row.HoverUnitId
}

func (row *RowData_Board) GetAirPlayerId() uint8 {
	return row.AirPlayerId
}

func (row *RowData_Board) GetAirUnitId() uint8 {
	return row.AirUnitId
}

type RowData_Units struct {
	X            uint16 `json:"x"`
	Y            uint16 `json:"y"`
	UnitType     uint8  `json:"unitType"`
	State        uint8  `json:"state"`
	Load         uint8  `json:"load"`
	Integrity    uint8  `json:"integrity"`
	Timestamp    uint32 `json:"timestamp"`
	Command      uint64 `json:"command"`
	CommandExtra uint64 `json:"commandExtra"`
	CommandMeta  uint8  `json:"commandMeta"`
	IsPreTicked  bool   `json:"isPreTicked"`
}

func (row *RowData_Units) GetX() uint16 {
	return row.X
}

func (row *RowData_Units) GetY() uint16 {
	return row.Y
}

func (row *RowData_Units) GetUnitType() uint8 {
	return row.UnitType
}

func (row *RowData_Units) GetState() uint8 {
	return row.State
}

func (row *RowData_Units) GetLoad() uint8 {
	return row.Load
}

func (row *RowData_Units) GetIntegrity() uint8 {
	return row.Integrity
}

func (row *RowData_Units) GetTimestamp() uint32 {
	return row.Timestamp
}

func (row *RowData_Units) GetCommand() uint64 {
	return row.Command
}

func (row *RowData_Units) GetCommandExtra() uint64 {
	return row.CommandExtra
}

func (row *RowData_Units) GetCommandMeta() uint8 {
	return row.CommandMeta
}

func (row *RowData_Units) GetIsPreTicked() bool {
	return row.IsPreTicked
}

type RowData_Buildings struct {
	X            uint16 `json:"x"`
	Y            uint16 `json:"y"`
	BuildingType uint8  `json:"buildingType"`
	State        uint8  `json:"state"`
	Integrity    uint8  `json:"integrity"`
	Timestamp    uint32 `json:"timestamp"`
}

func (row *RowData_Buildings) GetX() uint16 {
	return row.X
}

func (row *RowData_Buildings) GetY() uint16 {
	return row.Y
}

func (row *RowData_Buildings) GetBuildingType() uint8 {
	return row.BuildingType
}

func (row *RowData_Buildings) GetState() uint8 {
	return row.State
}

func (row *RowData_Buildings) GetIntegrity() uint8 {
	return row.Integrity
}

func (row *RowData_Buildings) GetTimestamp() uint32 {
	return row.Timestamp
}

type RowData_UnitPrototypes struct {
	Layer          uint8  `json:"layer"`
	ResourceCost   uint16 `json:"resourceCost"`
	ComputeCost    uint8  `json:"computeCost"`
	SpawnTime      uint8  `json:"spawnTime"`
	MaxIntegrity   uint8  `json:"maxIntegrity"`
	LandStrength   uint8  `json:"landStrength"`
	HoverStrength  uint8  `json:"hoverStrength"`
	AirStrength    uint8  `json:"airStrength"`
	AttackRange    uint8  `json:"attackRange"`
	AttackCooldown uint8  `json:"attackCooldown"`
	IsAssault      bool   `json:"isAssault"`
	IsWorker       bool   `json:"isWorker"`
}

func (row *RowData_UnitPrototypes) GetLayer() uint8 {
	return row.Layer
}

func (row *RowData_UnitPrototypes) GetResourceCost() uint16 {
	return row.ResourceCost
}

func (row *RowData_UnitPrototypes) GetComputeCost() uint8 {
	return row.ComputeCost
}

func (row *RowData_UnitPrototypes) GetSpawnTime() uint8 {
	return row.SpawnTime
}

func (row *RowData_UnitPrototypes) GetMaxIntegrity() uint8 {
	return row.MaxIntegrity
}

func (row *RowData_UnitPrototypes) GetLandStrength() uint8 {
	return row.LandStrength
}

func (row *RowData_UnitPrototypes) GetHoverStrength() uint8 {
	return row.HoverStrength
}

func (row *RowData_UnitPrototypes) GetAirStrength() uint8 {
	return row.AirStrength
}

func (row *RowData_UnitPrototypes) GetAttackRange() uint8 {
	return row.AttackRange
}

func (row *RowData_UnitPrototypes) GetAttackCooldown() uint8 {
	return row.AttackCooldown
}

func (row *RowData_UnitPrototypes) GetIsAssault() bool {
	return row.IsAssault
}

func (row *RowData_UnitPrototypes) GetIsWorker() bool {
	return row.IsWorker
}

type RowData_BuildingPrototypes struct {
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

func (row *RowData_BuildingPrototypes) GetWidth() uint8 {
	return row.Width
}

func (row *RowData_BuildingPrototypes) GetHeight() uint8 {
	return row.Height
}

func (row *RowData_BuildingPrototypes) GetResourceCost() uint16 {
	return row.ResourceCost
}

func (row *RowData_BuildingPrototypes) GetResourceCapacity() uint16 {
	return row.ResourceCapacity
}

func (row *RowData_BuildingPrototypes) GetComputeCapacity() uint8 {
	return row.ComputeCapacity
}

func (row *RowData_BuildingPrototypes) GetResourceMine() uint8 {
	return row.ResourceMine
}

func (row *RowData_BuildingPrototypes) GetMineTime() uint8 {
	return row.MineTime
}

func (row *RowData_BuildingPrototypes) GetMaxIntegrity() uint8 {
	return row.MaxIntegrity
}

func (row *RowData_BuildingPrototypes) GetBuildingTime() uint8 {
	return row.BuildingTime
}

func (row *RowData_BuildingPrototypes) GetIsArmory() bool {
	return row.IsArmory
}

func (row *RowData_BuildingPrototypes) GetIsEnvironment() bool {
	return row.IsEnvironment
}
