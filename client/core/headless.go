package core

import (
	"image"
	"time"

	"github.com/concrete-eth/archetype/arch"
	arch_client "github.com/concrete-eth/archetype/client"
	"github.com/concrete-eth/archetype/rpc"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/ethereum/go-ethereum/concrete/lib"
)

// TODO: consistent interface naming
// TODO: read-only player id 0

type IHeadlessClient interface {
	Game() *rts.Core
	PlayerId() uint8
	SetPlayerId(playerId uint8)

	TickPeriod() time.Duration
	SubTickPeriod() time.Duration
	// LastTickTime() time.Time
	LastNewBatchTime() time.Time

	Sync() (bool, bool, error)
	SyncUntil(tickIndex uint64) error
	InterpolatedSync() (bool, bool, error)
	Simulate(f func(core arch.Core))

	SendAction(action arch.Action) error
	Start()
	CreateUnit(unitType uint8)
	AssignUnit(unitId uint8, command rts.UnitCommandData)
	AssignUnitWithPath(unitId uint8, command rts.UnitCommandData, path []image.Point)
	AssignUnits(unitIds []uint8, command rts.UnitCommandData)
	AssignUnitsWithPath(unitIds []uint8, command rts.UnitCommandData, path []image.Point)
	PlaceBuilding(buildingType uint8, position image.Point)

	Interpolating() bool

	// NextSubTickIndex() uint32
	// StartInterpolation()
}

// Implements a headless client that can sync state and send actions.
type HeadlessClient struct {
	arch_client.Client
	hinter    *rpc.TxHinter
	playerId  uint8
	_tickTime time.Duration
}

var _ IHeadlessClient = (*HeadlessClient)(nil)

func NewHeadlessClient(
	kv lib.KeyValueStore,
	io *rpc.IO,
) *HeadlessClient {
	c := &rts.Core{}
	cli := io.NewClient(kv, c)
	hinter := io.Hinter()
	return &HeadlessClient{
		Client:    *cli,
		hinter:    hinter,
		_tickTime: cli.BlockTime() / time.Duration(c.TicksPerBlock()),
	}
}

// TODO: remove [?]
func (c *HeadlessClient) Game() *rts.Core {
	return c.Core().(*rts.Core)
}

// func (c *HeadlessClient) LastTickTime() time.Time {
// 	panic("not implemented")
// }

func (c *HeadlessClient) TickPeriod() time.Duration {
	return c.BlockTime()
}

func (c *HeadlessClient) SubTickPeriod() time.Duration {
	return c._tickTime
}

func (c *HeadlessClient) Interpolating() bool {
	return false
}

func (c *HeadlessClient) SetPlayerId(playerId uint8) {
	c.playerId = playerId
}

func (c *HeadlessClient) PlayerId() uint8 {
	return c.playerId
}

// Sends a Start action to the Tx sender
func (c *HeadlessClient) Start() {
	action := &rts.Start{}
	c.SendAction(action)
}

// Sends a UnitCreation action to the Tx sender
func (c *HeadlessClient) CreateUnit(unitType uint8) {
	action := &rts.UnitCreation{
		PlayerId: c.playerId,
		UnitType: unitType,
	}
	c.SendAction(action)
}

// Sends a UnitAssignation action to the Tx sender
func (c *HeadlessClient) AssignUnit(unitId uint8, command rts.UnitCommandData) {
	c.AssignUnitWithPath(unitId, command, nil)
}

// Sends a UnitAssignation action to the Tx sender
func (c *HeadlessClient) AssignUnitWithPath(unitId uint8, command rts.UnitCommandData, path []image.Point) {
	c.AssignUnitsWithPath([]uint8{unitId}, command, path)
}

func (c *HeadlessClient) AssignUnits(unitIds []uint8, command rts.UnitCommandData) {
	c.AssignUnitsWithPath(unitIds, command, nil)
}

func (c *HeadlessClient) AssignUnitsWithPath(unitIds []uint8, command rts.UnitCommandData, path []image.Point) {
	if len(path) > 4 {
		path = path[:4]
	} else if path == nil {
		path = make([]image.Point, 0)
	}
	commandPath := &rts.CommandPath{}
	commandPath.SetPath(path)

	actions := make([]arch.Action, 0)
	for _, unitId := range unitIds {
		action := &rts.UnitAssignation{
			PlayerId:     c.playerId,
			UnitId:       unitId,
			Command:      command.Uint64(),
			CommandExtra: commandPath.RawPath(),
			CommandMeta:  commandPath.Meta().Uint8(),
		}
		actions = append(actions, action)
	}
	c.SendActions(actions)
}

// Sends a BuildingPlacement action to the Tx sender
func (c *HeadlessClient) PlaceBuilding(buildingType uint8, position image.Point) {
	action := &rts.BuildingPlacement{
		PlayerId:     c.playerId,
		BuildingType: buildingType,
		X:            uint16(position.X),
		Y:            uint16(position.Y),
	}
	c.SendAction(action)
}
