package core

import (
	"errors"
	"image"
	"image/color"
	"time"

	"github.com/concrete-eth/archetype/arch"
	gen_utils "github.com/concrete-eth/archetype/utils"
	"github.com/concrete-eth/ark-rts/client/assets"
	"github.com/concrete-eth/ark-rts/client/decren"
	client_utils "github.com/concrete-eth/ark-rts/client/utils"
	"github.com/concrete-eth/ark-rts/gogen/archmod"
	"github.com/concrete-eth/ark-rts/gogen/datamod"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/ethereum/go-ethereum/concrete/lib"
	"github.com/ethereum/go-ethereum/log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	ErrQuit = errors.New("quit")
)

var (
	LayerName_Land  = rts.LayerId_Land.String()
	LayerName_Hover = rts.LayerId_Hover.String()
	LayerName_Air   = rts.LayerId_Air.String()

	LayerName_Terrain    = "terrain"
	LayerName_Buildings  = "buildings"
	LayerName_HudLines   = "hud-lines"
	LayerName_HudTerrain = "hud-terrain"
	LayerName_Bars       = "health-bars"
)

const (
	InternalTileSize  = 16
	IndicatorTileSize = 32
)

const (
	NeutralZoomTileSize = 2 * assets.TileSize
	MaxTileDisplaySize  = 3 * NeutralZoomTileSize
	MaxZoomLevel        = 2
	MinZoomLevel        = -1
)

type UpdatableWithRenderer interface {
	Update(c *CoreRenderer)
}

type DrawableWithRenderer interface {
	Draw(c *CoreRenderer, screen *ebiten.Image)
}

type ClientConfig struct {
	ScreenSize image.Point
}

type ClientSettings struct {
	Interpolate bool
	FixedCamera bool
}

// Holds a task to be executed at a certain time or block number.
type ScheduledTask struct {
	Func          func()
	Time          time.Time
	SubTickNumber uint32
}

// Holds a set of tasks pending execution.
type TaskSet struct {
	Tasks []*ScheduledTask
}

// Returns a new TaskSet.
func NewTaskSet() *TaskSet {
	return &TaskSet{
		Tasks: []*ScheduledTask{},
	}
}

// Adds a task to the set.
func (t *TaskSet) AddTask(task *ScheduledTask) {
	t.Tasks = append(t.Tasks, task)
}

// Executes all tasks that are due.
func (t *TaskSet) Update(c *CoreRenderer) {
	if len(t.Tasks) == 0 {
		return
	}
	pendingTasks := make([]*ScheduledTask, 0)
	for _, task := range t.Tasks {
		timeThresholdPassed := !task.Time.IsZero() && time.Now().After(task.Time)
		subTickThresholdPassed := task.SubTickNumber != 0 && c.Game().AbsSubTickIndex() >= task.SubTickNumber
		if timeThresholdPassed || subTickThresholdPassed {
			task.Func()
		} else {
			pendingTasks = append(pendingTasks, task)
		}
	}
	t.Tasks = pendingTasks
}

type InternalEvent struct {
	Id   uint8
	Data interface{}
}

type AnticipatedCommand struct {
	Command uint64
	Extra   uint64
	Meta    uint8
}

// Main game client object run by ebiten.RunGame.
type CoreRenderer struct {
	IHeadlessClient // Embedded headless client

	hintNonce uint64 // Hinter nonce

	config   ClientConfig   // Client configuration (immutable)
	settings ClientSettings // Client settings (modifiable)

	anticipatedObjects  map[rts.Object]struct{}           // Anticipated objects
	anticipatedCommands map[rts.Object]AnticipatedCommand // Anticipated commands

	cameraPosition   image.Point     // Camera position in internal scale
	zoomLevel        int             // Zoom level
	tileDisplaySize  int             // Tile display size in pixels
	boardDisplayRect image.Rectangle // Board display rectangle in pixels

	worldLayers *decren.LayerSet // World layers
	animations  *AnimationSet    // Animation set
	tasks       *TaskSet         // Task set

	lastMiddleClickedScreenPosition image.Point // Last right middle clicked screen position
	cameraPositionAtMiddleClick     image.Point // Camera position at last middle click

	position            map[rts.Object]image.Point      // Current pixel position of objects
	nextTilePosition    map[rts.Object]image.Point      // Next tile position of objects in tiles
	tableUpdatedObjects map[rts.Object]struct{}         // Objects whose table has been updated and need to be reset
	direction           map[rts.Object]assets.Direction // Current direction of objects
	anticipating        bool                            // True when actions and ticks are being anticipated

	lastSubTickTime       time.Time // Last tick time
	lastInterpolationTime time.Time // Last sub-tick time

	internalEventQueue []InternalEvent // Internal event buffer

	onNewBatch   func() // On new batch callback
	onCameraMove func() // On camera move callback

	spriteGetter assets.SpriteGetter // Sprite getter
}

var _ ebiten.Game = (*CoreRenderer)(nil)

// Instantiates a new Client object.
// Parameters:
// * headlessClient: the headless client to use
// * hinter: the hinter to use
// Returns:
// * the new Client object
// The headless client keeps the local game state in sync with the backend.
// The hinter provides a list of actions expected to be executed in the next block so the can
// be anticipated by the client.
func NewCoreRenderer(headlessClient IHeadlessClient, config ClientConfig, spriteGetter assets.SpriteGetter) *CoreRenderer {
	// Wait for the headless client to be synced up to initial state
	for !headlessClient.Game().IsInitialized() {
		err := headlessClient.SyncUntil(headlessClient.Game().BlockNumber() + 1)
		if err != nil {
			log.Error("Sync error", "error", err)
		}
	}

	// if config.Interpolate {
	// 	headlessClient.StartInterpolation()
	// }

	worldLayers := decren.NewLayerSet()
	animationSet := NewAnimationSet()
	taskSet := NewTaskSet()

	c := &CoreRenderer{
		IHeadlessClient: headlessClient,
		config:          config,
		settings:        ClientSettings{Interpolate: true, FixedCamera: true},

		hintNonce: 0,

		worldLayers: worldLayers,
		animations:  animationSet,
		tasks:       taskSet,

		position:            make(map[rts.Object]image.Point),
		nextTilePosition:    make(map[rts.Object]image.Point),
		tableUpdatedObjects: make(map[rts.Object]struct{}),
		direction:           make(map[rts.Object]assets.Direction),
		anticipating:        false,

		lastSubTickTime:       headlessClient.LastNewBatchTime(),
		lastInterpolationTime: headlessClient.LastNewBatchTime(),
		internalEventQueue:    make([]InternalEvent, 0),

		spriteGetter: spriteGetter,
	}

	c.boardDisplayRect = image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{c.config.ScreenSize.X, c.config.ScreenSize.Y},
		// Max: image.Point{c.config.ScreenSize.X - 254, c.config.ScreenSize.Y},
	}

	c.initLayers()

	camPos := c.Game().BoardSize().Mul(InternalTileSize).Div(2)
	camPos.Y += InternalTileSize / 2
	c.setCamera(camPos, 0)

	c.Game().SetEventHandler(c.onInternalEvent)
	c.anticipateSubTick()

	// Render the current state
	c.setAllBuildingSprites()
	c.setAllUnitSprites()

	return c
}

func (c *CoreRenderer) Config() ClientConfig {
	return c.config
}

func (c *CoreRenderer) Settings() ClientSettings {
	return c.settings
}

func (c *CoreRenderer) Interpolating() bool {
	return c.settings.Interpolate
}

func (c *CoreRenderer) TileDisplaySize() int {
	return c.tileDisplaySize
}

func (c *CoreRenderer) CameraPosition() image.Point {
	return c.cameraPosition
}

func (c *CoreRenderer) ZoomLevel() int {
	return c.zoomLevel
}

func (c *CoreRenderer) BoardDisplayRect() image.Rectangle {
	return c.boardDisplayRect
}

func (c *CoreRenderer) SetOnNewBatch(onNewBatch func()) {
	c.onNewBatch = onNewBatch
}

func (c *CoreRenderer) SetOnCameraMove(onCameraMove func()) {
	c.onCameraMove = onCameraMove
}

func (c *CoreRenderer) setAllBuildingSprites() {
	nPlayers := c.Game().GetMeta().GetPlayerCount()
	for playerId := uint8(0); playerId < nPlayers+1; playerId++ {
		c.Game().ForEachBuilding(playerId, func(buildingId uint8, building *datamod.BuildingsRow) {
			c.setBuildingSprite(playerId, buildingId, building)
		})
	}
}

func (c *CoreRenderer) setAllUnitSprites() {
	nPlayers := c.Game().GetMeta().GetPlayerCount()
	for playerId := uint8(0); playerId < nPlayers+1; playerId++ {
		c.Game().ForEachUnit(playerId, func(unitId uint8, unit *datamod.UnitsRow) {
			c.setUnitSprite(playerId, unitId, unit)
		})
	}
}

// Called when an internal event is emitted by the rts.
func (c *CoreRenderer) onInternalEvent(eventId uint8, data interface{}) {
	event := InternalEvent{
		Id:   eventId,
		Data: data,
	}
	c.internalEventQueue = append(c.internalEventQueue, event)
}

func (c *CoreRenderer) handleInternalEvents() {
	for _, event := range c.internalEventQueue {
		c.handleInternalEvent(event.Id, event.Data)
	}
	c.internalEventQueue = make([]InternalEvent, 0)
}

func (c *CoreRenderer) handleInternalEvent(eventId uint8, data interface{}) {
	if eventId == rts.InternalEventId_Shot {
		c.onShotEvent(data.(*rts.InternalEvent_Shot))
	} else if eventId == rts.InternalEventId_Spawned {
		c.onSpawnedEvent(data.(*rts.InternalEvent_Spawned))
	} else if eventId == rts.InternalEventId_Killed {
		c.onKilledEvent(data.(*rts.InternalEvent_Killed))
	} else if eventId == rts.InternalEventId_Built {
		c.onBuiltEvent(data.(*rts.InternalEvent_Built))
	}
}

func (c *CoreRenderer) Layers() *decren.LayerSet {
	return c.worldLayers
}

func (c *CoreRenderer) setCamera(position image.Point, zoomLevel int) {
	zoomLevel = gen_utils.Clamp(zoomLevel, MinZoomLevel, MaxZoomLevel)
	c.zoomLevel = zoomLevel
	if zoomLevel >= 0 {
		c.tileDisplaySize = NeutralZoomTileSize + zoomLevel*assets.TileSize
	} else {
		c.tileDisplaySize = NeutralZoomTileSize / gen_utils.Pow(2, -zoomLevel)
	}
	position = image.Point{
		X: gen_utils.Clamp(position.X, 0, InternalTileSize*(c.Game().BoardSize().X-1)),
		Y: gen_utils.Clamp(position.Y, 0, InternalTileSize*(c.Game().BoardSize().Y-1)),
	}
	c.cameraPosition = position
	c.positionLayerSources()

	if c.onCameraMove != nil {
		c.onCameraMove()
	}
}

// Sets the pixel position of a game object.
func (c *CoreRenderer) setPosition(object rts.Object, position image.Point) {
	c.position[object] = position
}

// Returns the pixel position of a game object.
func (c *CoreRenderer) getPosition(object rts.Object) image.Point {
	return c.position[object]
}

// Returns true if the pixel position of a game object is set.
func (c *CoreRenderer) hasPosition(object rts.Object) bool {
	_, ok := c.position[object]
	return ok
}

func (c *CoreRenderer) deletePosition(object rts.Object) {
	delete(c.position, object)
}

// Returns the next tile position of a given unit.
func (c *CoreRenderer) getUnitNextTilePosition(playerId uint8, unitId uint8) image.Point {
	return c.nextTilePosition[rts.Object{
		Type:     rts.ObjectType_Unit,
		PlayerId: playerId,
		ObjectId: unitId,
	}]
}

// Resets the pixel position of a unit to its canonical position.
func (c *CoreRenderer) resetUnitPosition(playerId uint8, unitId uint8, unit *datamod.UnitsRow) {
	tilePosition := rts.GetPositionAsPoint(unit)
	c.setPosition(rts.Object{
		Type:     rts.ObjectType_Unit,
		PlayerId: playerId,
		ObjectId: unitId,
	}, tilePosition.Mul(InternalTileSize))
}

// Returns the sprite object of a given building.
func (c *CoreRenderer) getBuildingSpriteObject(playerId uint8, buildingId uint8) *decren.Sprite {
	return c.worldLayers.Layer(LayerName_Buildings).Sprite(playerId, buildingId)
}

// Returns the sprite object of a given unit.
func (c *CoreRenderer) getUnitSpriteObject(playerId uint8, unitId uint8, unit *datamod.UnitsRow) *decren.Sprite {
	var (
		protoId   = unit.GetUnitType()
		proto     = c.Game().GetUnitPrototype(protoId)
		layerId   = rts.LayerId(proto.GetLayer())
		layerName = layerId.String()
	)
	return c.worldLayers.Layer(layerName).Sprite(playerId, unitId)
}

func (c *CoreRenderer) deleteUndefinedUnitSpriteObject(playerId uint8, unitId uint8) {
	c.worldLayers.Layer(LayerName_Land).Sprite(playerId, unitId).Delete()
	c.worldLayers.Layer(LayerName_Hover).Sprite(playerId, unitId).Delete()
	c.worldLayers.Layer(LayerName_Air).Sprite(playerId, unitId).Delete()
}

// Returns the sprite object of a given unit's health bar.
func (c *CoreRenderer) getHealthBarSpriteObject(object rts.Object) *decren.Sprite {
	return c.worldLayers.Layer(LayerName_Bars).Sprite("health", object.Type, object.PlayerId, object.ObjectId)
}

// Returns the sprite object of a given unit's build progress bar.
func (c *CoreRenderer) getBuildBarSpriteObject(object rts.Object) *decren.Sprite {
	return c.worldLayers.Layer(LayerName_Bars).Sprite("build", object.Type, object.PlayerId, object.ObjectId)
}

// Sets the pixel position of every unit to the interpolation of its canonical and next tile position.
func (c *CoreRenderer) interpolate() {
	// Compute the expected change in pixel position if interpolation were to be done
	timeDelta := time.Since(c.lastInterpolationTime)
	expectedInternalPositionDelta := int(c.movementSpeed() * timeDelta.Seconds())
	// If the expected change is less than the size unit pixel size, don't interpolate
	if expectedInternalPositionDelta < InternalTileSize/assets.TileSize {
		return
	}
	c.lastInterpolationTime = time.Now()

	tickFraction := c.tickFraction()
	nPlayers := c.Game().GetMeta().GetPlayerCount()
	for playerId := uint8(1); playerId < nPlayers+1; playerId++ {
		c.Game().ForEachUnit(playerId, func(unitId uint8, unit *datamod.UnitsRow) {
			c.interpolateUnitPosition(playerId, unitId, unit, tickFraction)
		})
	}
}

func (c *CoreRenderer) tickFraction() float64 {
	return float64(time.Since(c.lastSubTickTime)) / float64(c.SubTickPeriod())
}

func (c *CoreRenderer) interpolateUnitPosition(playerId uint8, unitId uint8, unit *datamod.UnitsRow, tickFraction float64) {
	object := rts.Object{
		Type:     rts.ObjectType_Unit,
		PlayerId: playerId,
		ObjectId: unitId,
	}
	if unitState := rts.UnitState(unit.GetState()); unitState.IsDeadOrInactive() {
		return
	}
	var (
		canonicalTilePosition     = rts.GetPositionAsPoint(unit)
		canonicalInternalPosition = canonicalTilePosition.Mul(InternalTileSize)
		nextTilePosition          = c.getUnitNextTilePosition(playerId, unitId)
		nextInternalPosition      = nextTilePosition.Mul(InternalTileSize)
		deltaInternalPosition     = nextInternalPosition.Sub(canonicalInternalPosition)
		nextInterpolatedPosition  = canonicalInternalPosition.Add(image.Point{
			X: int(float64(deltaInternalPosition.X) * tickFraction),
			Y: int(float64(deltaInternalPosition.Y) * tickFraction),
		})
	)

	c.setPosition(object, nextInterpolatedPosition)
	c.setUnitSpritePosition(playerId, unitId, unit)
}

// Pre-runs the next tick and sets the next tile position of every unit.
func (c *CoreRenderer) anticipateSubTick() {
	// c.anticipating = true
	c.Simulate(func(_core arch.Core) {
		core := _core.(*rts.Core)
		arch.RunSingleTick(core)
		core.ForEachPlayer(func(playerId uint8, player *datamod.PlayersRow) {
			core.ForEachUnit(playerId, func(unitId uint8, unit *datamod.UnitsRow) {
				c.anticipateUnitSubTick(playerId, unitId, unit)
			})
		})
	})
	// c.anticipating = false
}

func (c *CoreRenderer) anticipateUnitSubTick(playerId uint8, unitId uint8, unit *datamod.UnitsRow) {
	object := rts.Object{
		Type:     rts.ObjectType_Unit,
		PlayerId: playerId,
		ObjectId: unitId,
	}
	tilePosition := rts.GetPositionAsPoint(unit)
	// Set the next tile position
	c.nextTilePosition[object] = tilePosition
}

// Pre-runs any hinted actions and updates the sprites of the affected objects.
func (c *CoreRenderer) anticipateActions() {
	hinter := c.Hinter()
	if hinter == nil {
		return
	}
	hintNonce := hinter.HintNonce()
	if hintNonce <= c.hintNonce {
		return
	}

	hintNonce, hintsBatch := hinter.GetHints()
	c.hintNonce = hintNonce

	canonGame := c.Game()

	c.Simulate(func(_gg arch.Core) {
		simGame := _gg.(*rts.Core)

		anticipatedObjects := make(map[rts.Object]struct{}, 0)
		commandAnticipatedObject := make(map[rts.Object]struct{}, 0)

		// Set the child instances state update handler to update the present sprites based on anticipated table updates
		simGame.SetSetFieldHandler(func(table arch.TableSchema, rowKey lib.RowKey, columnName string, value []byte) {
			if table.Name == "Buildings" {
				playerId := rowKey[0].(uint8)
				buildingId := rowKey[1].(uint8)
				object := rts.Object{
					Type:     rts.ObjectType_Building,
					PlayerId: playerId,
					ObjectId: buildingId,
				}
				if columnName == "state" {
					anticipatedObjects[object] = struct{}{}
				}
			} else if table.Name == "Units" {
				playerId := rowKey[0].(uint8)
				unitId := rowKey[1].(uint8)
				object := rts.Object{
					Type:     rts.ObjectType_Unit,
					PlayerId: playerId,
					ObjectId: unitId,
				}
				if columnName == "state" {
					anticipatedObjects[object] = struct{}{}
				} else if columnName == "command" {
					commandAnticipatedObject[object] = struct{}{}
				}
			}
		})

		// Execute every hinted non-tick action on the child game instance
		for _, hint := range hintsBatch {
			for _, action := range hint {
				switch action := action.(type) {
				case *rts.Tick, *rts.Initialization, *rts.Start:
					continue
				default:
					err := archmod.ActionSchemas.ExecuteAction(action, simGame)
					if err != nil {
						log.Debug("Anticipating action error", "action", action, "error", err)
					}
				}
			}
		}

		// Anticipate objects
		c.anticipating = true
		for obj := range anticipatedObjects {
			delete(c.anticipatedObjects, obj)
			if obj.Type == rts.ObjectType_Building {
				building := simGame.GetBuilding(obj.PlayerId, obj.ObjectId)
				c.setBuildingSprite(obj.PlayerId, obj.ObjectId, building)
			} else if obj.Type == rts.ObjectType_Unit {
				unit := simGame.GetUnit(obj.PlayerId, obj.ObjectId)
				// All buildings are in the same layer, but units can be in any of multiple layers
				// We must clear them as the layer of the unit that corresponds to this ID may have changed
				c.deleteUndefinedUnitSpriteObject(obj.PlayerId, obj.ObjectId)
				c.setUnitSprite(obj.PlayerId, obj.ObjectId, unit)

				c.resetUnitPosition(obj.PlayerId, obj.ObjectId, unit)
				c.setUnitSpritePosition(obj.PlayerId, obj.ObjectId, unit)
			}
		}
		c.anticipating = false

		// Re-anticipate all previously anticipated object that have not been anticipated this round
		// to clear the sprites of objects that are no longer anticipated
		// Note the state is read from the canonical state, not the anticipated state
		for obj := range c.anticipatedObjects {
			// Skip object that have already been set in the canonical state
			playerId := obj.PlayerId
			if obj.Type == rts.ObjectType_Building {
				buildingId := obj.ObjectId
				nBuildings := canonGame.GetPlayer(playerId).GetBuildingCount()
				if buildingId <= nBuildings {
					continue
				}
				building := canonGame.GetBuilding(playerId, buildingId)
				c.setBuildingSprite(playerId, buildingId, building)
			} else if obj.Type == rts.ObjectType_Unit {
				c.deleteUndefinedUnitSpriteObject(obj.PlayerId, obj.ObjectId)
				unitId := obj.ObjectId
				unit := canonGame.GetUnit(playerId, unitId)
				c.setUnitSprite(playerId, unitId, unit)

				if unitState := rts.UnitState(unit.GetState()); !unitState.IsNil() {
					c.resetUnitPosition(obj.PlayerId, obj.ObjectId, unit)
					c.setUnitSpritePosition(obj.PlayerId, obj.ObjectId, unit)
				}
			}
		}

		c.anticipatedObjects = anticipatedObjects

		c.anticipatedCommands = make(map[rts.Object]AnticipatedCommand, 0)
		for obj := range commandAnticipatedObject {
			unit := simGame.GetUnit(obj.PlayerId, obj.ObjectId)
			c.anticipatedCommands[obj] = AnticipatedCommand{
				Command: unit.GetCommand(),
				Extra:   unit.GetCommandExtra(),
				Meta:    unit.GetCommandMeta(),
			}
		}
	})
}

// Initialize the sprite layers and set the tile display size.
func (c *CoreRenderer) initLayers() {
	for _, layer := range []struct {
		LayerId  string
		DestRect image.Rectangle
		Depth    int
		ColorM   colorm.ColorM
		Cache    bool
	}{
		{
			LayerId:  LayerName_Terrain,
			DestRect: c.boardDisplayRect,
			Depth:    0,
			ColorM:   assets.NewTerrainColorMatrix(),
			Cache:    true,
		},
		{
			LayerId:  LayerName_Buildings,
			DestRect: c.boardDisplayRect,
			Depth:    10,
			ColorM:   assets.NewBuildingColorMatrix(),
			Cache:    true,
		},
		{
			LayerId:  LayerName_Land,
			DestRect: c.boardDisplayRect,
			Depth:    20,
			ColorM:   assets.NewUnitColorMatrix(),
			Cache:    !c.settings.Interpolate,
		},
		{
			LayerId:  LayerName_Hover,
			DestRect: c.boardDisplayRect,
			Depth:    30,
			ColorM:   assets.NewUnitColorMatrix(),
			Cache:    !c.settings.Interpolate,
		},
		{
			LayerId:  LayerName_Air,
			DestRect: c.boardDisplayRect,
			Depth:    40,
			ColorM:   assets.NewUnitColorMatrix(),
			Cache:    !c.settings.Interpolate,
		},
		{
			LayerId:  LayerName_HudLines,
			DestRect: c.boardDisplayRect,
			Depth:    25,
			ColorM:   colorm.ColorM{},
			Cache:    false,
		},
		{
			LayerId:  LayerName_HudTerrain,
			DestRect: c.boardDisplayRect,
			Depth:    1,
			ColorM:   colorm.ColorM{},
			Cache:    !c.settings.Interpolate,
		},
		{
			LayerId:  LayerName_Bars,
			DestRect: c.boardDisplayRect,
			Depth:    50,
			ColorM:   colorm.ColorM{},
			Cache:    false,
		},
	} {
		c.worldLayers.Layer(layer.LayerId).
			SetDestinationRect(layer.DestRect).
			SetDepth(layer.Depth).
			SetColorMatrix(layer.ColorM).
			Cache(layer.Cache)
	}
	c.initTerrainLayer()

	c.worldLayers.Layer(LayerName_Bars).SetVisible(false)
}

// Add decorative cracks to terrain and the spawn points sprites.
func (c *CoreRenderer) initTerrainLayer() {
	terrainLayer := c.worldLayers.Layer(LayerName_Terrain)
	initTerrain(terrainLayer, assets.MapTilesetId_Royale)
	// TODO: verify vs map?
}

// Initialize the sprite layers and set the tile display size.
func (c *CoreRenderer) positionLayerSources() {
	var (
		tileSizeLayerSize    = c.boardDisplayRect.Size().Mul(assets.TileSize).Div(c.tileDisplaySize)
		tileSizeSourceOrigin = c.cameraPosition.Mul(assets.TileSize).Div(InternalTileSize).Sub(tileSizeLayerSize.Div(2))

		terrainLayerRect = image.Rectangle{Min: tileSizeSourceOrigin, Max: tileSizeSourceOrigin.Add(tileSizeLayerSize)}

		indicatorLayerSize   = tileSizeLayerSize.Mul(IndicatorTileSize).Div(assets.TileSize)
		indicatorLayerOrigin = tileSizeSourceOrigin.Mul(IndicatorTileSize).Div(assets.TileSize)
		indicatorLayerRect   = image.Rectangle{Min: indicatorLayerOrigin, Max: indicatorLayerOrigin.Add(indicatorLayerSize)}
	)

	c.worldLayers.Layer(LayerName_Terrain).SetSourceRect(terrainLayerRect)
	c.worldLayers.Layer(LayerName_Buildings).SetSourceRect(terrainLayerRect)
	c.worldLayers.Layer(LayerName_Land).SetSourceRect(terrainLayerRect)
	c.worldLayers.Layer(LayerName_Hover).SetSourceRect(terrainLayerRect)
	c.worldLayers.Layer(LayerName_Air).SetSourceRect(terrainLayerRect)

	c.worldLayers.Layer(LayerName_Bars).SetSourceRect(indicatorLayerRect)

	// Hud layers have static sources
	c.worldLayers.Layer(LayerName_HudTerrain).SetSourceRect(c.boardDisplayRect)
	c.worldLayers.Layer(LayerName_HudLines).SetSourceRect(c.boardDisplayRect)
}

// Set the sprite properties of a given building and its indicators (health and build progress bars).
func (c *CoreRenderer) setBuildingSprite(playerId uint8, buildingId uint8, building *datamod.BuildingsRow) {
	buildingState := rts.BuildingState(building.GetState())

	object := rts.Object{
		Type:     rts.ObjectType_Building,
		PlayerId: playerId,
		ObjectId: buildingId,
	}

	spriteObj := c.getBuildingSpriteObject(playerId, buildingId)
	healthBarSpriteObj := c.getHealthBarSpriteObject(object)
	buildBarSpriteObj := c.getBuildBarSpriteObject(object)

	if buildingState.IsNil() || buildingState == rts.BuildingState_Destroyed {
		// Remove the sprites if the building is nil, cancelled or destroyed
		spriteObj.Delete()
		healthBarSpriteObj.Delete()
		buildBarSpriteObj.Delete()
		return
	}

	if c.anticipating {
		// If anticipating, return unless the anticipated state is unpaid or building
		if buildingState != rts.BuildingState_Unpaid && buildingState != rts.BuildingState_Building {
			return
		}
	}

	var (
		buildingPosition = rts.GetPositionAsPoint(building)
		protoId          = building.GetBuildingType()
		proto            = c.Game().GetBuildingPrototype(protoId)
	)

	// Set building sprite properties
	layerPosition := buildingPosition.Mul(assets.TileSize)
	spriteOrigin := c.spriteGetter.GetBuildingSpriteOrigin(protoId)
	shiftedLayerPosition := layerPosition.Sub(spriteOrigin)

	spriteImg := c.spriteGetter.GetBuildingSprite(playerId, protoId, buildingState)

	colorM := colorm.ColorM{}
	if buildingState == rts.BuildingState_Unpaid {
		colorM.Concat(assets.UnpaidColorMatrix)
	}
	if c.anticipating {
		colorM.Concat(assets.AnticipatedColorMatrix)
	}

	spriteObj.
		SetPosition(shiftedLayerPosition).
		SetImage(spriteImg).
		FitToImage().
		SetColorMatrix(colorM)

	// Set building indicator sprite properties

	size := rts.GetDimensionsAsPoint(proto)
	barLayerPosition := buildingPosition.Mul(IndicatorTileSize).
		Add(image.Point{size.X*IndicatorTileSize/2 - healthBarWidth/2, -1 * healthBarHeight})

	healthBarSpriteObj.SetPosition(barLayerPosition)
	buildBarSpriteObj.SetPosition(barLayerPosition.Add(image.Point{0, healthBarHeight + healthBarHeight/2}))

	integrity := building.GetIntegrity()
	maxIntegrity := proto.GetMaxIntegrity()
	if integrity != maxIntegrity {
		healthBarImg := NewHealthBarImage(float64(integrity) / float64(maxIntegrity))
		healthBarSpriteObj.SetImage(healthBarImg).FitToImage()
	} else {
		// Hide the health bar if the building is at full health
		healthBarSpriteObj.SetImage(nil)
	}

	timestamp := building.GetTimestamp()
	if buildingState == rts.BuildingState_Building && timestamp != 0 {
		var (
			buildTime   = proto.GetBuildingTime()
			progress    = c.Game().AbsSubTickIndex() - timestamp
			buildBarImg = NewSpawnProgressBarImage(float64(progress) / float64(buildTime))
		)
		buildBarSpriteObj.SetImage(buildBarImg).FitToImage()
	} else {
		// Hide the build bar if the building is not being built
		buildBarSpriteObj.SetImage(nil)
	}
}

// Set the sprite properties of a given unit.
func (c *CoreRenderer) setUnitSprite(playerId uint8, unitId uint8, unit *datamod.UnitsRow) {
	unitState := rts.UnitState(unit.GetState())

	object := rts.Object{
		Type:     rts.ObjectType_Unit,
		PlayerId: playerId,
		ObjectId: unitId,
	}

	spriteObj := c.getUnitSpriteObject(playerId, unitId, unit)
	healthBarSpriteObj := c.getHealthBarSpriteObject(object)

	if unitState.IsNil() || unitState == rts.UnitState_Dead {
		// Remove the sprites if the unit is nil or dead
		c.deletePosition(object)
		healthBarSpriteObj.Delete()
		// if c.anticipating {
		// 	// If the unit does not exist in the game state, try to delete
		// 	// the sprite from all unit layers as it may be in either of them.
		// 	c.deleteUndefinedUnitSpriteObject(playerId, unitId)
		// } else {
		// 	spriteObj.Delete()
		// }
		spriteObj.Delete()
		return
	} else if unitState == rts.UnitState_Inactive {
		c.resetUnitPosition(playerId, unitId, unit)
	}

	if !c.hasPosition(object) {
		// Set the internal position to the canonical position if the position is not set
		c.resetUnitPosition(playerId, unitId, unit)
		c.setUnitSpritePosition(playerId, unitId, unit)
	}
	if c.anticipating || !c.Interpolating() {
		// If anticipating or interpolation is disabled, set the sprite position to the canonical position
		// When interpolating, the sprite position is set by the interpolation function
		c.setUnitSpritePosition(playerId, unitId, unit)
	}

	if c.anticipating {
		if unitState != rts.UnitState_Unpaid && unitState != rts.UnitState_Spawning {
			return
		}
	}

	c.setUnitSpriteImage(playerId, unitId, unit)
	c.setUnitSpriteColorMatrix(playerId, unitId, unit)
}

// Set the sprite and health bar position of a given unit.
func (c *CoreRenderer) setUnitSpritePosition(playerId uint8, unitId uint8, unit *datamod.UnitsRow) {
	object := rts.Object{
		Type:     rts.ObjectType_Unit,
		PlayerId: playerId,
		ObjectId: unitId,
	}
	if !c.Interpolating() {
		// Set the internal position to the canonical position if interpolation is disabled or the position is not set
		c.resetUnitPosition(playerId, unitId, unit)
	}
	spriteObj := c.getUnitSpriteObject(playerId, unitId, unit)
	healthBarSprite := c.getHealthBarSpriteObject(object)

	internalPosition := c.getPosition(object)
	layerPosition := internalPosition.Mul(assets.TileSize).Div(InternalTileSize).Sub(assets.UnitSpriteOrigin)
	spriteObj.SetPosition(layerPosition)

	barLayerPosition := internalPosition.Mul(IndicatorTileSize).Div(assets.TileSize).Add(image.Point{
		X: (IndicatorTileSize - healthBarWidth) / 2,
		Y: -2 * healthBarHeight,
	})
	healthBarSprite.SetPosition(barLayerPosition)
}

// Set the sprite image of a given unit
func (c *CoreRenderer) setUnitSpriteImage(playerId uint8, unitId uint8, unit *datamod.UnitsRow) {
	var (
		unitState = rts.UnitState(unit.GetState())
		protoId   = unit.GetUnitType()
		proto     = c.Game().GetUnitPrototype(protoId)
	)

	object := rts.Object{
		Type:     rts.ObjectType_Unit,
		PlayerId: playerId,
		ObjectId: unitId,
	}

	spriteObj := c.getUnitSpriteObject(playerId, unitId, unit)
	healthBarSpriteObj := c.getHealthBarSpriteObject(object)

	if unitState.IsDeadOrInactive() {
		spriteObj.SetImage(nil)
		healthBarSpriteObj.SetImage(nil)
		return
	}

	var direction assets.Direction

	if dir, ok := c.direction[object]; ok {
		direction = dir
	} else {
		if playerId == 1 {
			direction = assets.Direction_Right
		} else {
			direction = assets.Direction_Left
		}
	}

	var delta image.Point

	if unitState.IsActive() {
		// Compute the delta between the canonical and next tile position
		canonicalTilePosition := rts.GetPositionAsPoint(unit)
		nextTilePosition := c.getUnitNextTilePosition(playerId, unitId)
		delta = nextTilePosition.Sub(canonicalTilePosition)

		if delta == (image.Point{}) && !proto.GetIsWorker() {
			// If the unit is not moving, try to find a target it is firing at
			if targetMatch := c.Game().GetUnitToFireAt(c.Game().GetUnitObject(playerId, unitId)); !targetMatch.IsNil() {
				targetUnit := c.Game().GetUnit(targetMatch.PlayerId, targetMatch.ObjectId)
				targetTilePosition := rts.GetPositionAsPoint(targetUnit)
				delta = targetTilePosition.Sub(canonicalTilePosition)
			} else {
				// If no target is found, use the delta between the unit and its target
				commandData := rts.FighterCommandData(unit.GetCommand())
				if commandData.Type().IsTargetingPosition() {
					targetTilePosition := commandData.TargetPosition()
					delta = targetTilePosition.Sub(canonicalTilePosition)
				} else {
					var targetArea image.Rectangle
					if commandData.Type().IsTargetingBuilding() {
						targetPlayerId, targetBuildingId := commandData.TargetBuilding()
						targetArea = c.Game().GetBuildingArea(targetPlayerId, targetBuildingId)
					} else {
						taretPlayerId, targetUnitId := commandData.TargetUnit()
						targetPos := rts.GetPositionAsPoint(c.Game().GetUnit(taretPlayerId, targetUnitId))
						targetArea = image.Rectangle{Min: targetPos, Max: targetPos.Add(image.Point{1, 1})}
					}
					delta = rts.DeltasToAreaAsPoint(canonicalTilePosition, targetArea)
				}
			}
		}
		if delta != (image.Point{}) {
			direction = assets.DirectionFromDelta(delta)
		}
	}
	c.direction[object] = direction

	spriteImg := c.spriteGetter.GetUnitSprite(playerId, protoId, direction)
	spriteObj.SetImage(spriteImg).FitToImage()
	layerId := rts.LayerId(proto.GetLayer())

	// Set shadows for air and hover units
	if layerId == rts.LayerId_Hover {
		spriteObj.SetShadow(decren.Shadow{
			Enabled: true,
			Type:    decren.ShadowType_Cast,
			Offset:  image.Point{0, assets.TileSize / 2},
			Color:   assets.DarkShadowColor,
		})
	} else if layerId == rts.LayerId_Air {
		spriteObj.SetShadow(decren.Shadow{
			Enabled: true,
			Type:    decren.ShadowType_Cast,
			Offset:  image.Point{0, assets.TileSize / 2},
			Color:   assets.DarkShadowColor,
		})
	}

	// Set the health bar sprite

	if unitState.IsSpawning() {
		var (
			spawnTime    = proto.GetSpawnTime()
			timestamp    = unit.GetTimestamp()
			progress     = c.Game().AbsSubTickIndex() - timestamp
			healthBarImg = NewSpawnProgressBarImage(float64(progress) / float64(spawnTime))
		)
		healthBarSpriteObj.SetImage(healthBarImg).FitToImage()
	} else {
		integrity := unit.GetIntegrity()
		maxIntegrity := proto.GetMaxIntegrity()
		if integrity == maxIntegrity {
			// Hide the spawn bar if the unit is not spawning or damaged
			healthBarSpriteObj.SetImage(nil)
		} else {
			healthBarImg := NewHealthBarImage(float64(integrity) / float64(maxIntegrity))
			healthBarSpriteObj.SetImage(healthBarImg).FitToImage()
		}
	}
}

// Set the sprite color matrix of a given unit.
func (c *CoreRenderer) setUnitSpriteColorMatrix(playerId uint8, unitId uint8, unit *datamod.UnitsRow) {
	spriteObj := c.getUnitSpriteObject(playerId, unitId, unit)
	colorM := colorm.ColorM{}

	unitState := rts.UnitState(unit.GetState())

	if unitState == rts.UnitState_Unpaid {
		colorM.Concat(assets.UnpaidColorMatrix)
	} else if unitState == rts.UnitState_Spawning {
		colorM.Concat(assets.SpawningColorMatrix)
	}
	if c.anticipating {
		colorM.Concat(assets.AnticipatedColorMatrix)
	}

	spriteObj.SetColorMatrix(colorM)
}

// Trigger animations for shots
func (c *CoreRenderer) onShotEvent(shot *rts.InternalEvent_Shot) {
	var (
		targetType     = shot.Target.Type
		targetPlayerId = shot.Target.PlayerId
		targetId       = shot.Target.ObjectId
	)

	// Trigger a flash animation on the target
	var targetSpriteObj *decren.Sprite
	var imgOverride *ebiten.Image

	if targetType == rts.ObjectType_Building {
		building := c.Game().GetBuilding(targetPlayerId, targetId)
		protoId := building.GetBuildingType()
		imgOverride = c.spriteGetter.GetBuildingSprite(targetPlayerId, protoId, rts.BuildingState_Built)
		targetSpriteObj = c.getBuildingSpriteObject(targetPlayerId, targetId)
	} else if targetType == rts.ObjectType_Unit {
		unit := c.Game().GetUnit(targetPlayerId, targetId)
		protoId := unit.GetUnitType()
		unitDirection := c.direction[shot.Target]
		imgOverride = c.spriteGetter.GetUnitSprite(targetPlayerId, protoId, unitDirection)
		targetSpriteObj = c.getUnitSpriteObject(targetPlayerId, targetId, unit)
	}
	_ = imgOverride
	c.animations.RunAnimation(NewFlashAnimation(targetSpriteObj, imgOverride), AnimationConfig{
		FPS:  8,
		Mode: AnimationMode_Once,
	})

	// Trigger a fire animation on the attacker
	var (
		attackerDirection = c.direction[shot.Attacker]
		attackerUnit      = c.Game().GetUnit(shot.Attacker.PlayerId, shot.Attacker.ObjectId)
		attackerProtoId   = attackerUnit.GetUnitType()
		attackerSprite    = c.getUnitSpriteObject(shot.Attacker.PlayerId, shot.Attacker.ObjectId, attackerUnit)
	)
	c.animations.RunAnimation(NewFireAnimation(attackerSprite, shot.Attacker.PlayerId, attackerProtoId, attackerDirection, c.spriteGetter), AnimationConfig{
		FPS:  8,
		Mode: AnimationMode_Once,
	})
}

// Remove unit spawn bar when a unit spawns
func (c *CoreRenderer) onSpawnedEvent(spawn *rts.InternalEvent_Spawned) {
	healthBarSpriteObj := c.getHealthBarSpriteObject(spawn.Unit)
	spawnBarImg := NewSpawnProgressBarImage(1)
	healthBarSpriteObj.SetImageOverride(spawnBarImg)
	c.tasks.AddTask(&ScheduledTask{
		Time: time.Now().Add(500 * time.Millisecond),
		Func: func() {
			healthBarSpriteObj.SetImageOverride(nil)
		},
	})
}

func (c *CoreRenderer) onKilledEvent(kill *rts.InternalEvent_Killed) {
}

// Remove building build bar when a building is built
func (c *CoreRenderer) onBuiltEvent(built *rts.InternalEvent_Built) {
	buildBarSpriteObj := c.getBuildBarSpriteObject(built.Building)
	buildBarImg := NewSpawnProgressBarImage(1)
	buildBarSpriteObj.SetImageOverride(buildBarImg)
	c.tasks.AddTask(&ScheduledTask{
		Time: time.Now().Add(500 * time.Millisecond),
		Func: func() {
			buildBarSpriteObj.SetImageOverride(nil)
		},
	})
}

// Returns the movement speed of units in internal pixels per second.
func (c *CoreRenderer) movementSpeed() float64 {
	return float64(InternalTileSize) / c.SubTickPeriod().Seconds()
}

// Returns the tile position that corresponds to the given screen position.
func (c *CoreRenderer) ScreenCoordToTileCoord(screenPosition image.Point) image.Point {
	layerSourceOrigin := c.worldLayers.Layer(LayerName_Terrain).SourceRect().Min
	scaledSourceOrigin := layerSourceOrigin.Mul(c.tileDisplaySize).Div(assets.TileSize)
	tilePosition := screenPosition.
		Sub(c.boardDisplayRect.Min).
		Add(scaledSourceOrigin).
		Div(c.tileDisplaySize)
	return tilePosition
}

func (c *CoreRenderer) TileCoordToDisplayCoord(tilePosition image.Point) image.Point {
	layerSourceOrigin := c.worldLayers.Layer(LayerName_Terrain).SourceRect().Min
	scaledSourceOrigin := layerSourceOrigin.Mul(c.tileDisplaySize).Div(assets.TileSize)
	displayCoord := tilePosition.Mul(c.tileDisplaySize).Sub(scaledSourceOrigin)
	return displayCoord
}

// Returns the screen position that corresponds to the given tile position.
func (c *CoreRenderer) TileCoordToScreenCoord(tilePosition image.Point) image.Point {
	displayCoord := c.TileCoordToDisplayCoord(tilePosition)
	screenPosition := displayCoord.Add(c.boardDisplayRect.Min)
	return screenPosition
}

// Returns the screen position of a given unit.
func (c *CoreRenderer) GetUnitScreenPosition(playerId uint8, unitId uint8) image.Point {
	object := rts.Object{
		Type:     rts.ObjectType_Unit,
		PlayerId: playerId,
		ObjectId: unitId,
	}
	internalPosition := c.getPosition(object)
	screenPosition := c.TileCoordToScreenCoord(image.Point{0, 0}).
		Add(internalPosition.Mul(c.tileDisplaySize).Div(InternalTileSize))
	return screenPosition
}

func (c *CoreRenderer) dragMoveCamera() {
	if c.settings.FixedCamera {
		return
	}
	var (
		cursorScreenPosition = client_utils.CursorPosition()
		newCameraPosition    = c.cameraPosition
	)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		c.lastMiddleClickedScreenPosition = cursorScreenPosition
		c.cameraPositionAtMiddleClick = c.cameraPosition
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		delta := cursorScreenPosition.Sub(c.lastMiddleClickedScreenPosition)
		newCameraPosition = c.cameraPositionAtMiddleClick.Sub(delta.Mul(InternalTileSize).Div(c.tileDisplaySize))
	}
	if newCameraPosition != c.cameraPosition {
		c.setCamera(newCameraPosition, c.zoomLevel)
	}
}

func (c *CoreRenderer) Update() error {
	// Syncing
	var newBatch, subTicked bool
	var err error
	if c.Game().HasStarted() {
		newBatch, subTicked, err = c.InterpolatedSync()
	} else {
		newBatch, subTicked, err = c.Sync()
	}
	if err != nil {
		return err
	}

	c.anticipateActions()
	c.handleInternalEvents()
	if newBatch || subTicked {
		c.lastSubTickTime = time.Now()
		c.anticipateSubTick()
		c.setAllBuildingSprites()
		c.setAllUnitSprites()
	}
	if newBatch {
		if c.onNewBatch != nil {
			c.onNewBatch()
		}
	}

	// Interpolation
	if c.Interpolating() {
		c.interpolate()
	}
	// Animations
	c.animations.Update(c)
	// Tasks
	c.tasks.Update(c)

	return nil
}

// Draw the game on screen.
func (c *CoreRenderer) Draw(screen *ebiten.Image) {
	c.dragMoveCamera()
	c.worldLayers.Draw(screen)

	// Simple hack to get progress bars to look good
	for _, sprite := range c.worldLayers.Layer(LayerName_Bars).Sprites() {
		img := sprite.ImageOverride()
		if img == nil {
			img = sprite.Image()
		}
		if img == nil {
			continue
		}
		positionInLayer := sprite.Position()
		centerInLayer := positionInLayer.Add(image.Point{healthBarWidth / 2, healthBarHeight / 2})
		centerInScreen := c.TileCoordToScreenCoord(image.Point{0, 0}).
			Add(centerInLayer.Mul(c.tileDisplaySize).Div(IndicatorTileSize))
		sizeInScreen := img.Bounds().Size()
		positionInScreen := centerInScreen.Sub(sizeInScreen.Div(2))
		op := client_utils.NewDrawOptions(image.Rectangle{
			Min: positionInScreen,
			Max: positionInScreen.Add(sizeInScreen),
		}, img.Bounds())
		colorm.DrawImage(screen, img, colorm.ColorM{}, op)
	}
}

// Return the layout for ebiten.
func (c *CoreRenderer) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return c.config.ScreenSize.X, c.config.ScreenSize.Y
}

const (
	healthBarWidth  = 32
	healthBarHeight = 4
)

// Create a new health bar image.
func NewHealthBarImage(integrity float64) *ebiten.Image {
	return NewProgressBar(image.Point{healthBarWidth, healthBarHeight}, integrity, color.RGBA{0xff, 0x00, 0x00, 0xff})
}

// Create a new spawn progress bar image.
func NewSpawnProgressBarImage(progress float64) *ebiten.Image {
	progress = gen_utils.Clamp(progress, 0, 1)
	return NewProgressBar(image.Point{healthBarWidth, healthBarHeight}, progress, color.RGBA{0xff, 0xff, 0xff, 0xff})
}

// Create a new progress bar image.
func NewProgressBar(size image.Point, value float64, clr color.Color) *ebiten.Image {
	track := ebiten.NewImage(healthBarWidth, healthBarHeight)
	bar := ebiten.NewImage(1, 1)

	track.Fill(color.Black)
	bar.Fill(clr)

	op := client_utils.NewDrawOptions(image.Rectangle{
		Min: image.Point{1, 1},
		Max: image.Point{1 + int(float64(healthBarWidth-2)*value), healthBarHeight - 1},
	}, bar.Bounds())

	colorm.DrawImage(track, bar, colorm.ColorM{}, op)
	return track
}
