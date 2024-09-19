package core

import (
	"image"
	"time"

	"github.com/concrete-eth/archetype/arch"
	arch_client "github.com/concrete-eth/archetype/client"
	"github.com/concrete-eth/archetype/rpc"
	"github.com/concrete-eth/ark-royale/rts"
	"github.com/ethereum/go-ethereum/concrete/lib"
)

// TODO: read-only player id 0

type IHeadlessClient interface {
	Game() *rts.Core
	PlayerId() uint8
	SetPlayerId(playerId uint8)

	// TickPeriod() time.Duration
	SubTickPeriod() time.Duration
	LastNewBatchTime() time.Time

	Sync() (bool, bool, error)
	SyncUntil(tickIndex uint64) error
	InterpolatedSync() (bool, bool, error)
	// Interpolating() bool
	Simulate(f func(core arch.Core))
	Hinter() *rpc.TxHinter

	SendAction(action arch.Action) error
	Start()
	CreateUnit(unitType uint8, position image.Point)
}

// Implements a headless client that can sync state and send actions.
type HeadlessClient struct {
	*arch_client.Client
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
		Client:    cli,
		hinter:    hinter,
		_tickTime: cli.BlockTime() / time.Duration(c.TicksPerBlock()),
	}
}

func (c *HeadlessClient) Game() *rts.Core {
	return c.Core().(*rts.Core)
}

func (c *HeadlessClient) Hinter() *rpc.TxHinter {
	return c.hinter
}

// func (c *HeadlessClient) LastTickTime() time.Time {
// 	panic("not implemented")
// }

// func (c *HeadlessClient) TickPeriod() time.Duration {
// 	return c.BlockTime()
// }

func (c *HeadlessClient) SubTickPeriod() time.Duration {
	return c._tickTime
}

// func (c *HeadlessClient) Interpolating() bool {
// 	return false
// }

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
func (c *HeadlessClient) CreateUnit(unitType uint8, position image.Point) {
	action := &rts.UnitCreation{
		PlayerId: c.playerId,
		UnitType: unitType,
		X:        uint16(position.X),
		Y:        uint16(position.Y),
	}
	c.SendAction(action)
}
