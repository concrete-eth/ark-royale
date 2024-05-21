package rts

import (
	"fmt"
	"image"

	"github.com/concrete-eth/ark-rts/gogen/datamod"
)

func GetAttackStrength(proto *datamod.UnitPrototypesRow, layerId LayerId) uint8 {
	switch layerId {
	case LayerId_Land:
		return proto.GetLandStrength()
	case LayerId_Hover:
		return proto.GetHoverStrength()
	case LayerId_Air:
		return proto.GetAirStrength()
	default:
		panic(fmt.Sprintf("invalid layer: %d", layerId))
	}
}

type UnitState uint8

const (
	UnitState_Nil UnitState = iota
	UnitState_Unpaid
	UnitState_Spawning
	UnitState_Active
	UnitState_Inactive
	UnitState_Dead
	UnitState_Count
)

func (c UnitState) Uint8() uint8 {
	return uint8(c)
}

func (c UnitState) IsNil() bool {
	return c == UnitState_Nil
}

func (c UnitState) IsSpawning() bool {
	return c == UnitState_Spawning
}

func (c UnitState) IsActive() bool {
	return c == UnitState_Active
}

func (c UnitState) IsPaid() bool {
	switch c {
	case UnitState_Spawning, UnitState_Active, UnitState_Inactive, UnitState_Dead:
		return true
	default:
		return false
	}
}

func (c UnitState) HasSpawned() bool {
	switch c {
	case UnitState_Active, UnitState_Inactive, UnitState_Dead:
		return true
	default:
		return false
	}
}

func (c UnitState) IsAlive() bool {
	switch c {
	case UnitState_Spawning, UnitState_Active, UnitState_Inactive:
		return true
	default:
		return false
	}
}

func (c UnitState) IsDeadOrInactive() bool {
	switch c {
	case UnitState_Inactive, UnitState_Dead:
		return true
	default:
		return false
	}
}

type UnitCommandType interface {
	Uint8() uint8
	IsBusy() bool
}

type UnitCommandData interface {
	Uint64() uint64
	String() string
}

type WorkerCommandType uint8

var _ UnitCommandType = WorkerCommandType_Idle

const (
	WorkerCommandType_Idle WorkerCommandType = iota
	WorkerCommandType_Gather
	WorkerCommandType_Build
	WorkerCommandType_Count
)

func (c WorkerCommandType) Uint8() uint8 {
	return uint8(c)
}

func (c WorkerCommandType) IsIdle() bool {
	return c == WorkerCommandType_Idle
}

func (c WorkerCommandType) IsBusy() bool {
	return c == WorkerCommandType_Gather || c == WorkerCommandType_Build
}

// <unused> uint40, CommandType uint8, PlayerId uint8, BuildingId uint8
type WorkerCommandData uint64

var _ UnitCommandData = WorkerCommandData(WorkerCommandType_Idle)

func NewWorkerCommandData(commandType WorkerCommandType) WorkerCommandData {
	return WorkerCommandData(uint64(commandType) << 16)
}

func (c WorkerCommandData) Uint64() uint64 {
	return uint64(c)
}

func (c WorkerCommandData) Type() WorkerCommandType {
	return WorkerCommandType(c >> 16 & 0xFF)
}

func (c WorkerCommandData) TargetPlayerId() uint8 {
	return uint8(c >> 8 & 0xFF)
}

func (c *WorkerCommandData) SetTargetPlayerId(playerId uint8) {
	*c &^= 0xFF00
	*c |= WorkerCommandData(uint64(playerId) << 8)
}

func (c WorkerCommandData) TargetBuildingId() uint8 {
	return uint8(c & 0xFF)
}

func (c *WorkerCommandData) SetTargetBuildingId(buildingId uint8) {
	*c &^= 0xFF
	*c |= WorkerCommandData(uint64(buildingId))
}

func (c *WorkerCommandData) SetTargetBuilding(playerId uint8, buildingId uint8) {
	c.SetTargetPlayerId(playerId)
	c.SetTargetBuildingId(buildingId)
}

func (c *WorkerCommandData) TargetBuilding() (uint8, uint8) {
	return c.TargetPlayerId(), c.TargetBuildingId()
}

func (c WorkerCommandData) String() string {
	switch c.Type() {
	case WorkerCommandType_Idle:
		return "Idle"
	case WorkerCommandType_Gather:
		return fmt.Sprintf("Gather [%d, %d]", c.TargetPlayerId(), c.TargetBuildingId())
	case WorkerCommandType_Build:
		return fmt.Sprintf("Build [%d, %d]", c.TargetPlayerId(), c.TargetBuildingId())
	default:
		return "Unknown"
	}
}

type FighterCommandType uint8

var _ UnitCommandType = FighterCommandType_HoldPosition

const (
	FighterCommandType_HoldPosition FighterCommandType = iota
	FighterCommandType_AttackBuilding
	FighterCommandType_Count
)

func (c FighterCommandType) Uint8() uint8 {
	return uint8(c)
}

func (c FighterCommandType) IsIdle() bool {
	return false
}

func (c FighterCommandType) IsBusy() bool {
	return true
}

func (c FighterCommandType) IsTargetingPosition() bool {
	return c == FighterCommandType_HoldPosition
}

func (c FighterCommandType) IsTargetingBuilding() bool {
	return c == FighterCommandType_AttackBuilding
}

// <unused> uint24, Command uint8, Alpha uint16, Beta uint16
type FighterCommandData uint64

var _ UnitCommandData = FighterCommandData(FighterCommandType_HoldPosition)

func NewFighterCommandData(command FighterCommandType) FighterCommandData {
	return FighterCommandData(uint64(command) << 32)
}

func (c FighterCommandData) Uint64() uint64 {
	return uint64(c)
}

func (c FighterCommandData) Type() FighterCommandType {
	return FighterCommandType(c >> 32 & 0xFF)
}

func (c FighterCommandData) Alpha() uint16 {
	return uint16(c >> 16 & 0xFFFF)
}

func (c *FighterCommandData) SetAlpha(alpha uint16) {
	*c &^= 0xFFFF0000
	*c |= FighterCommandData(uint64(alpha) << 16)
}

func (c FighterCommandData) Beta() uint16 {
	return uint16(c & 0xFFFF)
}

func (c *FighterCommandData) SetBeta(beta uint16) {
	*c &^= 0xFFFF
	*c |= FighterCommandData(uint64(beta))
}

func (c FighterCommandData) TargetPosition() image.Point {
	return image.Point{
		X: int(c.Alpha()),
		Y: int(c.Beta()),
	}
}

func (c *FighterCommandData) SetTargetPosition(p image.Point) {
	c.SetAlpha(uint16(p.X))
	c.SetBeta(uint16(p.Y))
}

func (c FighterCommandData) TargetPlayerId() uint8 {
	return uint8(c.Alpha())
}

func (c *FighterCommandData) SetTargetPlayerId(playerId uint8) {
	c.SetAlpha(uint16(playerId))
}

func (c FighterCommandData) TargetBuildingId() uint8 {
	return uint8(c.Beta())
}

func (c *FighterCommandData) SetTargetBuildingId(buildingId uint8) {
	c.SetBeta(uint16(buildingId))
}

func (c *FighterCommandData) SetTargetBuilding(playerId uint8, buildingId uint8) {
	c.SetTargetPlayerId(playerId)
	c.SetTargetBuildingId(buildingId)
}

func (c *FighterCommandData) TargetBuilding() (uint8, uint8) {
	return c.TargetPlayerId(), c.TargetBuildingId()
}

func (c FighterCommandData) String() string {
	switch c.Type() {
	case FighterCommandType_HoldPosition:
		return fmt.Sprintf("HoldPosition [%d, %d]", c.TargetPosition().X, c.TargetPosition().Y)
	case FighterCommandType_AttackBuilding:
		return fmt.Sprintf("AttackBuilding [%d, %d]", c.TargetPlayerId(), c.TargetBuildingId())
	default:
		return "Unknown"
	}
}

type commandPathMeta uint8

func (c commandPathMeta) Uint8() uint8 {
	return uint8(c)
}

func (c *commandPathMeta) SetPathLen(len uint8) {
	*c &^= 0x0F
	*c |= commandPathMeta(len & 0x0F)
}

func (c commandPathMeta) PathLen() uint8 {
	return uint8(c) & 0x0F
}

func (c commandPathMeta) HasPath() bool {
	return c.PathLen() > 0
}

func (c *commandPathMeta) setPointer(pointer uint8) {
	*c &^= 0xF0
	*c |= commandPathMeta(pointer << 4)
}

func (c *commandPathMeta) IncPointer() {
	pointer := c.Pointer()
	if pointer < c.PathLen() {
		c.setPointer(pointer + 1)
	}
}

func (c commandPathMeta) Pointer() uint8 {
	return uint8(c) >> 4
}

type CommandPath struct {
	path uint64
	meta commandPathMeta
}

func NewCommandPath(pathRaw uint64, metaRaw uint8) *CommandPath {
	meta := commandPathMeta(metaRaw)
	return &CommandPath{
		path: pathRaw,
		meta: meta,
	}
}

func (c *CommandPath) SetPath(path []image.Point) {
	if len(path) > 4 {
		path = path[:4]
	}
	c.meta = commandPathMeta(0)
	c.meta.SetPathLen(uint8(len(path)))
	c.path = 0
	for i, p := range path {
		x := uint8(p.X)
		y := uint8(p.Y)
		c.path |= uint64(x) << (i * 16)
		c.path |= uint64(y) << (i*16 + 8)
	}
}

func (c *CommandPath) GetPathPoint(i int) image.Point {
	if i >= int(c.meta.PathLen()) {
		return image.Point{}
	}
	x := uint8(c.path >> (i * 16) & 0xFF)
	y := uint8(c.path >> (i*16 + 8) & 0xFF)
	return image.Point{int(x), int(y)}
}

func (c *CommandPath) RawPath() uint64 {
	return c.path
}

func (c *CommandPath) Meta() commandPathMeta {
	return c.meta
}

func (c *CommandPath) Path() []image.Point {
	path := make([]image.Point, c.meta.PathLen())
	for i := range path {
		path[i] = c.GetPathPoint(i)
	}
	return path
}

func (c *CommandPath) PathLen() uint8 {
	return c.meta.PathLen()
}

func (c *CommandPath) HasPath() bool {
	return c.meta.HasPath()
}

func (c *CommandPath) IncPointer() {
	c.meta.IncPointer()
}

func (c *CommandPath) Pointer() uint8 {
	return c.meta.Pointer()
}

func (c *CommandPath) CurrentPoint() image.Point {
	return c.GetPathPoint(int(c.Pointer()))
}
