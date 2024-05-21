/* Autogenerated file. Do not edit manually. */

package datamod

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/codegen/datamod/codec"
	"github.com/ethereum/go-ethereum/concrete/crypto"
	"github.com/ethereum/go-ethereum/concrete/lib"
	"github.com/holiman/uint256"
)

// Reference imports to suppress errors if they are not used.
var (
	_ = common.Big1
	_ = codec.EncodeAddress
	_ = uint256.NewInt
)

// var (
//	PlayersDefaultKey = crypto.Keccak256([]byte("datamod.v1.Players"))
// )

func PlayersDefaultKey() []byte {
	return crypto.Keccak256([]byte("datamod.v1.Players"))
}

type PlayersRow struct {
	lib.DatastoreStruct
}

func NewPlayersRow(dsSlot lib.DatastoreSlot) *PlayersRow {
	sizes := []int{2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1}
	return &PlayersRow{*lib.NewDatastoreStruct(dsSlot, sizes)}
}

func (v *PlayersRow) Get() (
	spawnAreaX uint16,
	spawnAreaY uint16,
	workerPortX uint16,
	workerPortY uint16,
	curResource uint16,
	maxResource uint16,
	curArmories uint8,
	computeSupply uint8,
	computeDemand uint8,
	unitCount uint8,
	buildingCount uint8,
	buildingPayQueuePointer uint8,
	buildingBuildQueuePointer uint8,
	unitPayQueuePointer uint8,
) {
	return codec.DecodeUint16(2, v.GetField(0)),
		codec.DecodeUint16(2, v.GetField(1)),
		codec.DecodeUint16(2, v.GetField(2)),
		codec.DecodeUint16(2, v.GetField(3)),
		codec.DecodeUint16(2, v.GetField(4)),
		codec.DecodeUint16(2, v.GetField(5)),
		codec.DecodeUint8(1, v.GetField(6)),
		codec.DecodeUint8(1, v.GetField(7)),
		codec.DecodeUint8(1, v.GetField(8)),
		codec.DecodeUint8(1, v.GetField(9)),
		codec.DecodeUint8(1, v.GetField(10)),
		codec.DecodeUint8(1, v.GetField(11)),
		codec.DecodeUint8(1, v.GetField(12)),
		codec.DecodeUint8(1, v.GetField(13))
}

func (v *PlayersRow) Set(
	spawnAreaX uint16,
	spawnAreaY uint16,
	workerPortX uint16,
	workerPortY uint16,
	curResource uint16,
	maxResource uint16,
	curArmories uint8,
	computeSupply uint8,
	computeDemand uint8,
	unitCount uint8,
	buildingCount uint8,
	buildingPayQueuePointer uint8,
	buildingBuildQueuePointer uint8,
	unitPayQueuePointer uint8,
) {
	v.SetField(0, codec.EncodeUint16(2, spawnAreaX))
	v.SetField(1, codec.EncodeUint16(2, spawnAreaY))
	v.SetField(2, codec.EncodeUint16(2, workerPortX))
	v.SetField(3, codec.EncodeUint16(2, workerPortY))
	v.SetField(4, codec.EncodeUint16(2, curResource))
	v.SetField(5, codec.EncodeUint16(2, maxResource))
	v.SetField(6, codec.EncodeUint8(1, curArmories))
	v.SetField(7, codec.EncodeUint8(1, computeSupply))
	v.SetField(8, codec.EncodeUint8(1, computeDemand))
	v.SetField(9, codec.EncodeUint8(1, unitCount))
	v.SetField(10, codec.EncodeUint8(1, buildingCount))
	v.SetField(11, codec.EncodeUint8(1, buildingPayQueuePointer))
	v.SetField(12, codec.EncodeUint8(1, buildingBuildQueuePointer))
	v.SetField(13, codec.EncodeUint8(1, unitPayQueuePointer))
}

func (v *PlayersRow) GetSpawnAreaX() uint16 {
	data := v.GetField(0)
	return codec.DecodeUint16(2, data)
}

func (v *PlayersRow) SetSpawnAreaX(value uint16) {
	data := codec.EncodeUint16(2, value)
	v.SetField(0, data)
}

func (v *PlayersRow) GetSpawnAreaY() uint16 {
	data := v.GetField(1)
	return codec.DecodeUint16(2, data)
}

func (v *PlayersRow) SetSpawnAreaY(value uint16) {
	data := codec.EncodeUint16(2, value)
	v.SetField(1, data)
}

func (v *PlayersRow) GetWorkerPortX() uint16 {
	data := v.GetField(2)
	return codec.DecodeUint16(2, data)
}

func (v *PlayersRow) SetWorkerPortX(value uint16) {
	data := codec.EncodeUint16(2, value)
	v.SetField(2, data)
}

func (v *PlayersRow) GetWorkerPortY() uint16 {
	data := v.GetField(3)
	return codec.DecodeUint16(2, data)
}

func (v *PlayersRow) SetWorkerPortY(value uint16) {
	data := codec.EncodeUint16(2, value)
	v.SetField(3, data)
}

func (v *PlayersRow) GetCurResource() uint16 {
	data := v.GetField(4)
	return codec.DecodeUint16(2, data)
}

func (v *PlayersRow) SetCurResource(value uint16) {
	data := codec.EncodeUint16(2, value)
	v.SetField(4, data)
}

func (v *PlayersRow) GetMaxResource() uint16 {
	data := v.GetField(5)
	return codec.DecodeUint16(2, data)
}

func (v *PlayersRow) SetMaxResource(value uint16) {
	data := codec.EncodeUint16(2, value)
	v.SetField(5, data)
}

func (v *PlayersRow) GetCurArmories() uint8 {
	data := v.GetField(6)
	return codec.DecodeUint8(1, data)
}

func (v *PlayersRow) SetCurArmories(value uint8) {
	data := codec.EncodeUint8(1, value)
	v.SetField(6, data)
}

func (v *PlayersRow) GetComputeSupply() uint8 {
	data := v.GetField(7)
	return codec.DecodeUint8(1, data)
}

func (v *PlayersRow) SetComputeSupply(value uint8) {
	data := codec.EncodeUint8(1, value)
	v.SetField(7, data)
}

func (v *PlayersRow) GetComputeDemand() uint8 {
	data := v.GetField(8)
	return codec.DecodeUint8(1, data)
}

func (v *PlayersRow) SetComputeDemand(value uint8) {
	data := codec.EncodeUint8(1, value)
	v.SetField(8, data)
}

func (v *PlayersRow) GetUnitCount() uint8 {
	data := v.GetField(9)
	return codec.DecodeUint8(1, data)
}

func (v *PlayersRow) SetUnitCount(value uint8) {
	data := codec.EncodeUint8(1, value)
	v.SetField(9, data)
}

func (v *PlayersRow) GetBuildingCount() uint8 {
	data := v.GetField(10)
	return codec.DecodeUint8(1, data)
}

func (v *PlayersRow) SetBuildingCount(value uint8) {
	data := codec.EncodeUint8(1, value)
	v.SetField(10, data)
}

func (v *PlayersRow) GetBuildingPayQueuePointer() uint8 {
	data := v.GetField(11)
	return codec.DecodeUint8(1, data)
}

func (v *PlayersRow) SetBuildingPayQueuePointer(value uint8) {
	data := codec.EncodeUint8(1, value)
	v.SetField(11, data)
}

func (v *PlayersRow) GetBuildingBuildQueuePointer() uint8 {
	data := v.GetField(12)
	return codec.DecodeUint8(1, data)
}

func (v *PlayersRow) SetBuildingBuildQueuePointer(value uint8) {
	data := codec.EncodeUint8(1, value)
	v.SetField(12, data)
}

func (v *PlayersRow) GetUnitPayQueuePointer() uint8 {
	data := v.GetField(13)
	return codec.DecodeUint8(1, data)
}

func (v *PlayersRow) SetUnitPayQueuePointer(value uint8) {
	data := codec.EncodeUint8(1, value)
	v.SetField(13, data)
}

type Players struct {
	dsSlot lib.DatastoreSlot
}

func NewPlayers(ds lib.Datastore) *Players {
	dsSlot := ds.Get(PlayersDefaultKey())
	return &Players{dsSlot}
}

func NewPlayersFromSlot(dsSlot lib.DatastoreSlot) *Players {
	return &Players{dsSlot}
}
func (m *Players) Get(
	playerId uint8,
) *PlayersRow {
	dsSlot := m.dsSlot.Mapping().GetNested(
		codec.EncodeUint8(1, playerId),
	)
	return NewPlayersRow(dsSlot)
}
