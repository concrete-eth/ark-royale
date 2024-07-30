package rts

import (
	"errors"
	"image"

	"github.com/concrete-eth/archetype/arch"
	"github.com/concrete-eth/archetype/utils"
	"github.com/concrete-eth/ark-rts/gogen/archmod"
	"github.com/concrete-eth/ark-rts/gogen/datamod"
)

const (
	NilObjectId   = uint8(0)
	NilPlayerId   = NilObjectId
	NilBuildingId = NilObjectId
	NilUnitId     = NilObjectId
)

type ObjectType uint8

const (
	ObjectType_Nil ObjectType = iota
	ObjectType_Building
	ObjectType_Unit
	ObjectType_Count
)

func (t ObjectType) Uint8() uint8 {
	return uint8(t)
}

func (t ObjectType) IsNil() bool {
	return t == ObjectType_Nil
}

// TODO: stronger type checking
type Object struct {
	Type     ObjectType
	PlayerId uint8
	ObjectId uint8
}

func (o Object) IsNil() bool {
	return o.Type.IsNil()
}

type ObjectWithRow struct {
	_o Object
	_r interface{}
}

func (o ObjectWithRow) Object() Object {
	return o._o
}

func (o ObjectWithRow) IsNil() bool {
	return o._o.IsNil()
}

func (o ObjectWithRow) PlayerId() uint8 {
	return o._o.PlayerId
}

func (o ObjectWithRow) ObjectId() uint8 {
	return o._o.ObjectId
}

func (o ObjectWithRow) ObjectType() ObjectType {
	return o._o.Type
}

type BuildingObjectWithRow struct {
	ObjectWithRow
}

func (o BuildingObjectWithRow) Building() *datamod.BuildingsRow {
	if o._o.Type != ObjectType_Building {
		return nil
	}
	return o._r.(*datamod.BuildingsRow)
}

type UnitObjectWithRow struct {
	ObjectWithRow
}

func (o UnitObjectWithRow) Unit() *datamod.UnitsRow {
	if o._o.Type != ObjectType_Unit {
		return nil
	}
	return o._r.(*datamod.UnitsRow)
}

type LayerId uint8

const (
	LayerId_Land LayerId = iota
	LayerId_Hover
	LayerId_Air
	LayerId_Count
)

func (l LayerId) Uint8() uint8 {
	return uint8(l)
}

func (l LayerId) String() string {
	switch l {
	case LayerId_Land:
		return "land"
	case LayerId_Hover:
		return "hover"
	case LayerId_Air:
		return "air"
	default:
		return "unknown"
	}
}

const (
	InternalEventId_Shot uint8 = iota
	InternalEventId_Spawned
	InternalEventId_Killed
	InternalEventId_Built
	InternalEventId_Destroyed
)

type InternalEventHandler func(eventId uint8, data interface{})

type InternalEvent_Shot struct {
	Attacker Object
	Target   Object
}

type InternalEvent_Spawned struct {
	Unit Object
}

type InternalEvent_Killed struct {
	Unit Object
}

type InternalEvent_Built struct {
	Building Object
}

type InternalEvent_Destroyed struct {
	Building Object
}

type (
	Tick                      = arch.CanonicalTickAction
	Start                     = archmod.ActionData_Start
	PlayerAddition            = archmod.ActionData_AddPlayer
	Initialization            = archmod.ActionData_Initialize
	UnitCreation              = archmod.ActionData_CreateUnit
	UnitAssignation           = archmod.ActionData_AssignUnit
	BuildingPlacement         = archmod.ActionData_PlaceBuilding
	UnitPrototypeAddition     = archmod.ActionData_AddUnitPrototype
	BuildingPrototypeAddition = archmod.ActionData_AddBuildingPrototype
)

var (
	ErrNotInitialized                = errors.New("not initialized")
	ErrNotStarted                    = errors.New("not started")
	ErrAlreadyInitialized            = errors.New("already initialized")
	ErrAlreadyStarted                = errors.New("already started")
	ErrPlayerLimitReached            = errors.New("player limit reached")
	ErrUnitLimitReached              = errors.New("unit limit reached")
	ErrBuildingLimitReached          = errors.New("building limit reached")
	ErrUnitPrototypeLimitReached     = errors.New("unit prototype limit reached")
	ErrBuildingPrototypeLimitReached = errors.New("building prototype limit reached")
	ErrInvalidPlayerId               = errors.New("invalid player id")
	ErrInvalidUnitId                 = errors.New("invalid unit id")
	ErrInvalidBuildingId             = errors.New("invalid building id")
	ErrInvalidUnitType               = errors.New("invalid unit type")
	ErrInvalidBuildingType           = errors.New("invalid building type")
	ErrNonexistentPlayer             = errors.New("nonexistent player")
	ErrMissingMainBuilding           = errors.New("player has no main building")
	ErrMainBuildingDestroyed         = errors.New("player's main building destroyed")
	ErrAreaNotBuildable              = errors.New("area not buildable")
	ErrUnitDead                      = errors.New("unit dead")
)

type Core struct {
	arch.BaseCore
	eventHandler InternalEventHandler
}

var _ archmod.IActions = &Core{}

func (c *Core) TicksPerBlock() uint64 {
	return 1
}

func (c *Core) SetEventHandler(handler InternalEventHandler) {
	c.eventHandler = handler
}

func (c *Core) EventHandler() InternalEventHandler {
	return c.eventHandler
}

func (c *Core) emitInternalEvent(eventId uint8, data interface{}) {
	if c.eventHandler != nil {
		c.eventHandler(eventId, data)
	}
}

func (c *Core) GetMeta() *datamod.MetaRow {
	return datamod.NewMeta(c.Datastore()).Get()
}

func (c *Core) GetPlayer(playerId uint8) *datamod.PlayersRow {
	return datamod.NewPlayers(c.Datastore()).Get(playerId)
}

func (c *Core) GetBoardTile(x, y uint16) *datamod.BoardRow {
	return datamod.NewBoard(c.Datastore()).Get(x, y)
}

func (c *Core) GetBuilding(playerId uint8, buildingId uint8) *datamod.BuildingsRow {
	return datamod.NewBuildings(c.Datastore()).Get(playerId, buildingId)
}

func (c *Core) GetBuildingObject(playerId uint8, buildingId uint8) BuildingObjectWithRow {
	return BuildingObjectWithRow{
		ObjectWithRow: ObjectWithRow{
			_o: Object{
				Type:     ObjectType_Building,
				PlayerId: playerId,
				ObjectId: buildingId,
			},
			_r: c.GetBuilding(playerId, buildingId),
		},
	}
}

func (c *Core) GetUnit(playerId uint8, unitId uint8) *datamod.UnitsRow {
	return datamod.NewUnits(c.Datastore()).Get(playerId, unitId)
}

func (c *Core) GetUnitObject(playerId uint8, uintId uint8) UnitObjectWithRow {
	return UnitObjectWithRow{
		ObjectWithRow: ObjectWithRow{
			_o: Object{
				Type:     ObjectType_Unit,
				PlayerId: playerId,
				ObjectId: uintId,
			},
			_r: c.GetUnit(playerId, uintId),
		},
	}
}

func (c *Core) GetUnitPrototype(unitTypeId uint8) *datamod.UnitPrototypesRow {
	return datamod.NewUnitPrototypes(c.Datastore()).Get(unitTypeId)
}

func (c *Core) GetBuildingPrototype(prototypeId uint8) *datamod.BuildingPrototypesRow {
	return datamod.NewBuildingPrototypes(c.Datastore()).Get(prototypeId)
}

func (c *Core) AbsSubTickIndex() uint32 {
	return uint32(c.BlockNumber()*c.TicksPerBlock() + c.InBlockTickIndex())
}

type (
	PlayerFilter   func(playerId uint8, player *datamod.PlayersRow) bool
	BuildingFilter func(playerId, buildingId uint8, building *datamod.BuildingsRow) bool
	UnitFilter     func(playerId, unitId uint8, unit *datamod.UnitsRow) bool
)

func (c *Core) IterPlayers(filters ...PlayerFilter) *Iterator_uint8 {
	return &Iterator_uint8{
		Current: 0,
		Max:     c.GetMeta().GetPlayerCount(),
		Get: func(playerId uint8) (uint8, interface{}) {
			player := c.GetPlayer(playerId)
			for _, filter := range filters {
				if !filter(playerId, player) {
					return 0, nil
				}
			}
			return playerId, player
		},
	}
}

func (c *Core) IterBuildings(playerId uint8, filters ...BuildingFilter) *Iterator_uint8 {
	return &Iterator_uint8{
		Current: 0,
		Max:     c.GetPlayer(playerId).GetBuildingCount(),
		Get: func(buildingId uint8) (uint8, interface{}) {
			building := c.GetBuilding(playerId, buildingId)
			for _, filter := range filters {
				if !filter(playerId, buildingId, building) {
					return 0, nil
				}
			}
			return buildingId, building
		},
	}
}

func (c *Core) IterUnits(playerId uint8, filters ...UnitFilter) *Iterator_uint8 {
	return &Iterator_uint8{
		Current: 0,
		Max:     c.GetPlayer(playerId).GetUnitCount(),
		Get: func(unitId uint8) (uint8, interface{}) {
			unit := c.GetUnit(playerId, unitId)
			for _, filter := range filters {
				if !filter(playerId, unitId, unit) {
					return 0, nil
				}
			}
			return unitId, unit
		},
	}
}

func (c *Core) ForEachPlayer(forEach func(playerId uint8, player *datamod.PlayersRow)) {
	nPlayers := c.GetMeta().GetPlayerCount()
	for playerId := uint8(1); playerId < nPlayers+1; playerId++ {
		forEach(playerId, c.GetPlayer(playerId))
	}
}

func (c *Core) ForEachBuilding(playerId uint8, forEach func(buildingId uint8, building *datamod.BuildingsRow)) {
	nBuildings := c.GetPlayer(playerId).GetBuildingCount()
	for buildingId := uint8(1); buildingId < nBuildings+1; buildingId++ {
		forEach(buildingId, c.GetBuilding(playerId, buildingId))
	}
}

func (c *Core) ForEachUnit(playerId uint8, forEach func(unitId uint8, unit *datamod.UnitsRow)) {
	nUnits := c.GetPlayer(playerId).GetUnitCount()
	for unitId := uint8(1); unitId < nUnits+1; unitId++ {
		forEach(unitId, c.GetUnit(playerId, unitId))
	}
}

func (c *Core) ForEachBuildingPrototype(forEach func(prototypeId uint8, prototype *datamod.BuildingPrototypesRow)) {
	nPrototypes := c.GetMeta().GetBuildingPrototypeCount()
	for prototypeId := uint8(1); prototypeId < nPrototypes+1; prototypeId++ {
		forEach(prototypeId, c.GetBuildingPrototype(prototypeId))
	}
}

func (c *Core) ForEachUnitPrototype(forEach func(prototypeId uint8, prototype *datamod.UnitPrototypesRow)) {
	nPrototypes := c.GetMeta().GetUnitPrototypeCount()
	for prototypeId := uint8(1); prototypeId < nPrototypes+1; prototypeId++ {
		forEach(prototypeId, c.GetUnitPrototype(prototypeId))
	}
}

func (c *Core) SearchIter(iter IIterator_Uint8) (uint8, interface{}) {
	for iter.Next() {
		objectId, object := iter.Value()
		return objectId, object
	}
	return NilObjectId, nil
}

func (c *Core) GetWorkerPortPosition(playerId uint8) image.Point {
	player := c.GetPlayer(playerId)
	position := image.Point{
		int(player.GetWorkerPortX()),
		int(player.GetWorkerPortY()),
	}
	return position
}

func (c *Core) GetWorkerPortArea(playerId uint8) image.Rectangle {
	position := c.GetWorkerPortPosition(playerId)
	size := image.Point{1, 1}
	return image.Rectangle{
		Min: position,
		Max: position.Add(size),
	}
}

func (c *Core) GetBuildingArea(playerId uint8, buildingId uint8) image.Rectangle {
	building := c.GetBuilding(playerId, buildingId)
	position := GetPositionAsPoint(building)
	protoId := building.GetBuildingType()
	proto := c.GetBuildingPrototype(protoId)
	size := GetDimensionsAsPoint(proto)
	return image.Rectangle{
		Min: position,
		Max: position.Add(size),
	}
}

func (c *Core) GetMainBuilding(playerId uint8) *datamod.BuildingsRow {
	// First building is the main building
	return c.GetBuilding(playerId, 1)
}

func (c *Core) GetMainBuildingPosition(playerId uint8) image.Point {
	building := c.GetMainBuilding(playerId)
	return GetPositionAsPoint(building)
}

func (c *Core) GetMainBuildingArea(playerId uint8) image.Rectangle {
	return c.GetBuildingArea(playerId, 1)
}

func (c *Core) GetSpawnArea(playerId uint8) image.Rectangle {
	player := c.GetPlayer(playerId)
	position := image.Point{
		int(player.GetSpawnAreaX()),
		int(player.GetSpawnAreaY()),
	}
	size := image.Point{2, 2}
	area := image.Rectangle{
		Min: position,
		Max: position.Add(size),
	}
	return area
}

func (c *Core) GetSpawnPoint(layer LayerId, playerId uint8) (image.Point, bool) {
	spawnArea := c.GetSpawnArea(playerId)
	mainBuildingArea := c.GetMainBuildingArea(playerId)
	var nearestTile image.Point
	nearestDistance := -1
	for x := uint16(spawnArea.Min.X); x < uint16(spawnArea.Max.X); x++ {
		for y := uint16(spawnArea.Min.Y); y < uint16(spawnArea.Max.Y); y++ {
			tile := c.GetBoardTile(x, y)
			if IsTileEmptyAllLayers(tile) {
				distance := DistanceToArea(nearestTile, mainBuildingArea)
				if nearestDistance == -1 || distance < nearestDistance {
					nearestTile = image.Point{int(x), int(y)}
					nearestDistance = distance
				}
			}
		}
	}
	if nearestDistance == -1 {
		return image.Point{}, false
	}
	return nearestTile, true
}

type ObjectMatchByDistance struct {
	ObjectId uint8
	Distance int
}

type PlayerObjectMatchByDistance struct {
	PlayerId uint8
	ObjectMatchByDistance
}

func (m ObjectMatchByDistance) IsNil() bool {
	return m.ObjectId == NilObjectId
}

func (c *Core) GetNearest(iter IIterator_Uint8, area image.Rectangle) ObjectMatchByDistance {
	var match ObjectMatchByDistance
	for iter.Next() {
		objectId, _object := iter.Value()
		object := _object.(interface {
			GetX() uint16
			GetY() uint16
		})
		objectPosition := GetPositionAsPoint(object)
		distance := DistanceToArea(objectPosition, area)
		if match.IsNil() || distance < match.Distance {
			match.ObjectId = objectId
			match.Distance = distance
		}
		if distance == 0 {
			return match
		}
	}
	return match
}

func (c *Core) GetNearestEnemyUnits(playerId uint8, position image.Point, filters ...UnitFilter) []PlayerObjectMatchByDistance {
	matches := make([]PlayerObjectMatchByDistance, LayerId_Count)
	nPlayers := c.GetMeta().GetPlayerCount()
	for enemyPlayerId := uint8(1); enemyPlayerId < nPlayers+1; enemyPlayerId++ {
		if enemyPlayerId == playerId {
			continue
		}
		iter := c.IterUnits(enemyPlayerId, filters...)
		for iter.Next() {
			var (
				unitId, _unit = iter.Value()
				unit          = _unit.(*datamod.UnitsRow)
				unitPosition  = GetPositionAsPoint(unit)
				protoId       = unit.GetUnitType()
				proto         = c.GetUnitPrototype(protoId)
				layer         = proto.GetLayer()
				distance      = Distance(position, unitPosition)
			)
			if matches[layer].IsNil() || distance < matches[layer].Distance {
				matches[layer] = PlayerObjectMatchByDistance{
					PlayerId: enemyPlayerId,
					ObjectMatchByDistance: ObjectMatchByDistance{
						ObjectId: unitId,
						Distance: distance,
					},
				}
			}
		}
	}
	return matches
}

func (c *Core) BoardSize() image.Point {
	return image.Point{int(c.GetMeta().GetBoardWidth()), int(c.GetMeta().GetBoardHeight())}
}

func (c *Core) BoardRect() image.Rectangle {
	return image.Rectangle{Max: c.BoardSize()}
}

func (c *Core) TileIsInBoard(position image.Point) bool {
	return position.In(c.BoardRect())
}

func (c *Core) IsBuildableArea(area image.Rectangle) bool {
	if !area.In(c.BoardRect()) {
		return false
	}
	nPlayers := c.GetMeta().GetPlayerCount()
	for playerId := uint8(1); playerId < nPlayers+1; playerId++ {
		if area.Overlaps(c.GetSpawnArea(playerId)) {
			return false
		}
	}
	for x := uint16(area.Min.X); x < uint16(area.Max.X); x++ {
		for y := uint16(area.Min.Y); y < uint16(area.Max.Y); y++ {
			tile := c.GetBoardTile(x, y)
			if !IsTileEmpty(tile, LayerId_Land) {
				return false
			}
		}
	}
	return true
}

func (c *Core) ValidatePlayerId(playerId uint8) error {
	nPlayers := c.GetMeta().GetPlayerCount()
	if playerId == NilPlayerId || playerId > nPlayers {
		return ErrNonexistentPlayer
	}
	return nil
}

func (c *Core) ValidatePlayerOrEnvironmentId(playerId uint8) error {
	nPlayers := c.GetMeta().GetPlayerCount()
	if playerId > nPlayers {
		return ErrNonexistentPlayer
	}
	return nil
}

func (c *Core) ValidateUnitId(playerId uint8, unitId uint8) error {
	if unitId == NilUnitId {
		return ErrInvalidUnitId
	}
	if unitId > c.GetPlayer(playerId).GetUnitCount() {
		return ErrInvalidUnitId
	}
	return nil
}

func (c *Core) ValidateBuildingId(playerId uint8, buildingId uint8) error {
	if buildingId == NilBuildingId {
		return ErrInvalidBuildingId
	}
	if buildingId > c.GetPlayer(playerId).GetBuildingCount() {
		return ErrInvalidBuildingId
	}
	return nil
}

func (c *Core) ValidateUnitType(prototypeId uint8) error {
	if prototypeId == 0 {
		return ErrInvalidUnitType
	}
	if prototypeId > c.GetMeta().GetUnitPrototypeCount() {
		return ErrInvalidUnitType
	}
	return nil
}

func (c *Core) ValidateBuildingType(prototypeId uint8) error {
	if prototypeId == 0 {
		return ErrInvalidBuildingType
	}
	if prototypeId > c.GetMeta().GetBuildingPrototypeCount() {
		return ErrInvalidBuildingType
	}
	return nil
}

func (c *Core) setUnitState(obj UnitObjectWithRow, state UnitState) {
	obj.Unit().SetState(state.Uint8())
}

func (c *Core) setUnitSpawning(obj UnitObjectWithRow) {
	var (
		player      = c.GetPlayer(obj.PlayerId())
		unit        = obj.Unit()
		protoId     = unit.GetUnitType()
		proto       = c.GetUnitPrototype(protoId)
		computeCost = proto.GetComputeCost()
	)
	payPointer := obj.ObjectId()
	player.SetUnitPayQueuePointer(utils.SafeAddUint8(payPointer, 1))
	addComputeDemand(player, computeCost)
	c.setUnitState(obj, UnitState_Spawning)
	unit.SetTimestamp(c.AbsSubTickIndex())
}

func (c *Core) setUnitSpawned(obj UnitObjectWithRow) {
	timeNow := c.AbsSubTickIndex()

	c.setUnitState(obj, UnitState_Active)
	obj.Unit().SetTimestamp(timeNow)

	c.emitInternalEvent(InternalEventId_Spawned, &InternalEvent_Spawned{
		Unit: obj.Object(),
	})
}

func (c *Core) setUnitDead(obj UnitObjectWithRow) {
	var (
		player      = c.GetPlayer(obj.PlayerId())
		unit        = obj.Unit()
		proto       = c.GetUnitPrototype(unit.GetUnitType())
		computeCost = proto.GetComputeCost()
		tile        = c.GetBoardTile(unit.GetX(), unit.GetY())
		layer       = LayerId(proto.GetLayer())
	)
	c.setUnitState(obj, UnitState_Dead)
	EmptyUnitTile(tile, layer)

	// Execute death side-effects
	subComputeDemand(player, computeCost)
	c.resetBuildingProcessIfAny(obj)
	c.emitInternalEvent(InternalEventId_Killed, &InternalEvent_Killed{
		Unit: obj.Object(),
	})
}

func (c *Core) resetBuildingProcessIfAny(obj UnitObjectWithRow) {
	var (
		unit         = obj.Unit()
		unitPosition = GetPositionAsPoint(unit)
		proto        = c.GetUnitPrototype(unit.GetUnitType())
	)
	if proto.GetIsWorker() {
		// Unit is worker
		workerCommand := WorkerCommandData(unit.GetCommand())
		if workerCommand.Type() == WorkerCommandType_Build {
			// Worker is building
			var (
				targetBuildingId    = workerCommand.TargetBuildingId()
				targetBuilding      = c.GetBuilding(obj.PlayerId(), targetBuildingId)
				targetBuildingState = BuildingState(targetBuilding.GetState())
			)
			if targetBuildingState == BuildingState_Building {
				// Target building was being built
				targetBuildingArea := c.GetBuildingArea(obj.PlayerId(), targetBuildingId)
				if unitPosition.In(targetBuildingArea) {
					// Worker was at target building
					// Reset the timestamp
					targetBuilding.SetTimestamp(0)
				}
			}
		}
	}
}

func (c *Core) setUnitCommand(obj UnitObjectWithRow, command UnitCommandData) {
	unit := obj.Unit()
	unit.SetCommand(command.Uint64())
	unit.SetCommandExtra(0)
	unit.SetCommandMeta(0)
}

func (c *Core) setBuildingState(obj BuildingObjectWithRow, state BuildingState) {
	building := obj.Building()
	building.SetState(state.Uint8())
}

func (c *Core) setBuildingBuilding(obj BuildingObjectWithRow) {
	player := c.GetPlayer(obj.PlayerId())
	payPointer := obj.ObjectId()
	player.SetBuildingPayQueuePointer(utils.SafeAddUint8(payPointer, 1))
	c.setBuildingState(obj, BuildingState_Building)
}

func (c *Core) setBuildingBuilt(obj BuildingObjectWithRow) {
	var (
		player           = c.GetPlayer(obj.PlayerId())
		building         = obj.Building()
		protoId          = building.GetBuildingType()
		timeNow          = c.AbsSubTickIndex()
		proto            = c.GetBuildingPrototype(protoId)
		resourceCapacity = proto.GetResourceCapacity()
		computeCapacity  = proto.GetComputeCapacity()
	)

	c.setBuildingState(obj, BuildingState_Built)
	building.SetTimestamp(timeNow)

	if !proto.GetIsEnvironment() {
		// Execute building side effects
		addStorage(player, resourceCapacity)
		addComputeSupply(player, computeCapacity)
		if proto.GetIsArmory() {
			addArmory(player)
		}
		if obj.ObjectId() == 1 {
			addResource(player, resourceCapacity)
		}
	}

	c.emitInternalEvent(InternalEventId_Built, &InternalEvent_Built{
		Building: obj.Object(),
	})
}

func (c *Core) setBuildingDestroyed(obj BuildingObjectWithRow) {
	var (
		player   = c.GetPlayer(obj.PlayerId())
		building = obj.Building()
		position = GetPositionAsPoint(building)
		protoId  = building.GetBuildingType()
		proto    = c.GetBuildingPrototype(protoId)
		size     = GetDimensionsAsPoint(proto)
	)

	c.setBuildingState(obj, BuildingState_Destroyed)
	for x := uint16(position.X); x < uint16(position.X+size.X); x++ {
		for y := uint16(position.Y); y < uint16(position.Y+size.Y); y++ {
			tile := c.GetBoardTile(x, y)
			SetTileLandObject(tile, ObjectType_Nil, NilPlayerId, NilObjectId)
		}
	}
	// Execute building destruction side effects
	subStorage(player, proto.GetResourceCapacity())
	subComputeSupply(player, proto.GetComputeCapacity())
	if proto.GetIsArmory() {
		subArmory(player)
	}

	c.emitInternalEvent(InternalEventId_Destroyed, &InternalEvent_Destroyed{
		Building: obj.Object(),
	})
}

func (c *Core) tickPlayer(playerId uint8) {
	player := c.GetPlayer(playerId)
	nArmories := player.GetCurArmories()
	c.payForUnits(playerId)
	for i := uint8(1); i < nArmories; i++ {
		c.payForUnits(playerId)
	}
	c.payForAndAssignBuildings(playerId)
}

func (c *Core) payForUnits(playerId uint8) {
	var (
		player     = c.GetPlayer(playerId)
		nUnits     = player.GetUnitCount()
		payPointer = player.GetUnitPayQueuePointer()
	)
	if payPointer <= nUnits {
		var (
			unitId         = payPointer
			obj            = c.GetUnitObject(playerId, unitId)
			unit           = obj.Unit()
			protoId        = unit.GetUnitType()
			proto          = c.GetUnitPrototype(protoId)
			resources      = player.GetCurResource()
			computeSupply  = player.GetComputeSupply()
			computeDemand  = player.GetComputeDemand()
			computeSurplus = utils.SafeSubUint8(computeSupply, computeDemand)
			resourceCost   = proto.GetResourceCost()
			computeCost    = proto.GetComputeCost()
		)
		if resources >= resourceCost && computeSurplus >= computeCost {
			subResource(player, resourceCost)
			c.setUnitSpawning(obj)
		}
	}
}

func (c *Core) payForAndAssignBuildings(playerId uint8) {
	var (
		player     = c.GetPlayer(playerId)
		nBuildings = player.GetBuildingCount()
		payPointer = player.GetBuildingPayQueuePointer()
	)
	if payPointer <= nBuildings {
		var (
			buildingId   = payPointer
			obj          = c.GetBuildingObject(playerId, buildingId)
			building     = obj.Building()
			protoId      = building.GetBuildingType()
			proto        = c.GetBuildingPrototype(protoId)
			resources    = player.GetCurResource()
			resourceCost = proto.GetResourceCost()
		)
		if resources >= resourceCost {
			subResource(player, resourceCost)
			c.setBuildingBuilding(obj)
		}
	}
	buildPointer := player.GetBuildingBuildQueuePointer()
	for buildingId := buildPointer; buildingId < payPointer; buildingId++ {
		var (
			obj              = c.GetBuildingObject(playerId, buildingId)
			building         = obj.Building()
			buildingPosition = GetPositionAsPoint(building)
			buildingState    = BuildingState(building.GetState())
			protoId          = building.GetBuildingType()
			proto            = c.GetBuildingPrototype(protoId)
			buildingSize     = GetDimensionsAsPoint(proto)
		)
		if buildingState.HasBeenBuilt() {
			if buildingId == buildPointer {
				buildPointer = utils.SafeAddUint8(buildPointer, 1)
			}
			continue
		}
		// Check if there is a worker assigned to build this building already
		assigneeId, _ := c.SearchIter(c.IterUnits(playerId, func(playerId, unitId uint8, unit *datamod.UnitsRow) bool {
			unitProtoId := unit.GetUnitType()
			unitProto := c.GetUnitPrototype(unitProtoId)
			unitState := UnitState(unit.GetState())
			unitCommand := WorkerCommandData(unit.GetCommand())
			return unitProto.GetIsWorker() &&
				unitState.HasSpawned() &&
				unitState.IsAlive() &&
				unitCommand.Type() == WorkerCommandType_Build &&
				unitCommand.TargetBuildingId() == buildingId
		}))
		if assigneeId == NilUnitId {
			// Assign building to the nearest available worker
			iter := c.IterUnits(playerId, func(playerId, unitId uint8, unit *datamod.UnitsRow) bool {
				unitProtoId := unit.GetUnitType()
				unitProto := c.GetUnitPrototype(unitProtoId)
				unitState := UnitState(unit.GetState())
				unitCommand := WorkerCommandData(unit.GetCommand())
				return unitProto.GetIsWorker() &&
					unitState.HasSpawned() &&
					unitState.IsAlive() &&
					!unitCommand.Type().IsBusy()
			})
			match := c.GetNearest(iter, image.Rectangle{
				Min: buildingPosition,
				Max: buildingPosition.Add(buildingSize),
			})
			assigneeId = match.ObjectId
			if assigneeId == NilUnitId {
				// No available workers
				break
			} else {
				assigneeObj := c.GetUnitObject(playerId, assigneeId)
				c.setWorkerToBuild(assigneeObj, buildingId)
			}
		}
	}
	player.SetBuildingBuildQueuePointer(buildPointer)
}

func (c *Core) placeBuilding(playerId uint8, prototypeId uint8, position image.Point) uint8 {
	var (
		player       = c.GetPlayer(playerId)
		nBuildings   = player.GetBuildingCount()
		buildingId   = utils.SafeAddUint8(nBuildings, 1)
		proto        = c.GetBuildingPrototype(prototypeId)
		maxIntegrity = proto.GetMaxIntegrity()
		size         = GetDimensionsAsPoint(proto)
	)

	if buildingId == nBuildings {
		// Building limit reached
		return NilBuildingId
	}

	player.SetBuildingCount(buildingId)
	c.GetBuilding(playerId, buildingId).Set(
		uint16(position.X), uint16(position.Y),
		prototypeId,
		BuildingState_Unpaid.Uint8(),
		maxIntegrity,
		0,
	)

	for x := uint16(position.X); x < uint16(position.X+size.X); x++ {
		for y := uint16(position.Y); y < uint16(position.Y+size.Y); y++ {
			tile := c.GetBoardTile(x, y)
			SetTileLandObject(tile, ObjectType_Building, playerId, buildingId)
		}
	}

	return buildingId
}

func (c *Core) createUnit(playerId uint8, prototypeId uint8) uint8 {
	var (
		player       = c.GetPlayer(playerId)
		nUnits       = player.GetUnitCount()
		unitId       = utils.SafeAddUint8(nUnits, 1)
		proto        = c.GetUnitPrototype(prototypeId)
		layer        = LayerId(proto.GetLayer())
		maxIntegrity = proto.GetMaxIntegrity()
		timeNow      = c.AbsSubTickIndex()
	)

	position, ok := c.GetSpawnPoint(layer, playerId)
	if !ok {
		return NilUnitId
	}

	if unitId == nUnits {
		// Unit limit reached
		return NilUnitId
	}

	player.SetUnitCount(unitId)
	var command UnitCommandData
	if proto.GetIsWorker() {
		command = NewWorkerCommandData(WorkerCommandType_Idle)
	} else {
		_command := NewFighterCommandData(FighterCommandType_HoldPosition)
		_command.SetTargetPosition(position)
		command = _command
	}

	c.GetUnit(playerId, unitId).Set(
		uint16(position.X), uint16(position.Y),
		prototypeId,
		UnitState_Unpaid.Uint8(),
		0,
		maxIntegrity,
		timeNow,
		command.Uint64(),
		0,
		0,
		false,
	)

	tile := c.GetBoardTile(uint16(position.X), uint16(position.Y))
	SetTileUnit(tile, layer, playerId, unitId)

	return unitId
}

func (c *Core) assignUnit(obj UnitObjectWithRow, command UnitCommandData) {
	var (
		unit       = obj.Unit()
		curCommand = unit.GetCommand()
	)
	if curCommand != command.Uint64() {
		c.resetBuildingProcessIfAny(obj)
	}
	c.setUnitCommand(obj, command)
}

func (c *Core) assignUnitExternal(obj UnitObjectWithRow, command UnitCommandData, path *CommandPath) {
	c.tickUnit(obj)
	unit := obj.Unit()
	unit.SetIsPreTicked(true)
	c.assignUnit(obj, command)
	unit.SetCommandExtra(path.RawPath())
	unit.SetCommandMeta(path.Meta().Uint8())
}

func (c *Core) setWorkerToIdle(obj UnitObjectWithRow) {
	command := NewWorkerCommandData(WorkerCommandType_Idle)
	c.assignUnit(obj, command)
}

func (c *Core) setWorkerToGather(obj UnitObjectWithRow, targetBuildingId uint8) {
	command := NewWorkerCommandData(WorkerCommandType_Gather)
	command.SetTargetBuilding(NilPlayerId, targetBuildingId)
	c.assignUnit(obj, command)
}

func (c *Core) setWorkerToBuild(obj UnitObjectWithRow, targetBuildingId uint8) {
	command := NewWorkerCommandData(WorkerCommandType_Build)
	command.SetTargetBuilding(obj.PlayerId(), targetBuildingId)
	c.assignUnit(obj, command)
}

func (c *Core) unitCanMoveTo(position image.Point, layer LayerId) bool {
	tile := c.GetBoardTile(uint16(position.X), uint16(position.Y))
	switch layer {
	case LayerId_Hover:
		// Workers cannot overlap with other workers
		if !IsTileEmpty(tile, layer) {
			return false
		}
	case LayerId_Land:
		// Land units cannot overlap with other land units, air units, or buildings
		if !IsTileEmpty(tile, layer) || !IsTileEmpty(tile, LayerId_Air) {
			return false
		}
	case LayerId_Air:
		// Air units cannot overlap with other air units or land units
		if !IsTileEmpty(tile, layer) {
			return false
		}
		if ObjectType(tile.GetLandObjectType()) == ObjectType_Unit {
			return false
		}
	}
	return true
}

func (c *Core) moveUnitToTarget(obj UnitObjectWithRow, targetArea image.Rectangle) image.Point {
	var (
		unit     = obj.Unit()
		protoId  = unit.GetUnitType()
		proto    = c.GetUnitPrototype(protoId)
		layer    = LayerId(proto.GetLayer())
		position = GetPositionAsPoint(unit)
	)
	if position.In(targetArea) {
		return position
	}

	if !proto.GetIsWorker() {
		path := NewCommandPath(unit.GetCommandExtra(), unit.GetCommandMeta())
		if path.HasPath() {
			for path.Pointer() < path.PathLen() {
				pathPoint := path.CurrentPoint()
				pathPointArea := image.Rectangle{
					Min: pathPoint.Sub(image.Point{1, 1}),
					Max: pathPoint.Add(image.Point{2, 2}),
				}
				if position.In(pathPointArea) {
					path.IncPointer()
				} else {
					targetArea = pathPointArea
					break
				}
			}
			unit.SetCommandMeta(path.Meta().Uint8())
		}
	}

	_, nextPosition, _, _ := c.greedyPathFind(layer, position, targetArea, 1, 0)
	if nextPosition.Eq(position) {
		return position
	}

	curTile := c.GetBoardTile(uint16(position.X), uint16(position.Y))
	nextTile := c.GetBoardTile(uint16(nextPosition.X), uint16(nextPosition.Y))
	EmptyUnitTile(curTile, layer)
	SetTileUnit(nextTile, layer, obj.PlayerId(), obj.ObjectId())
	unit.SetX(uint16(nextPosition.X))
	unit.SetY(uint16(nextPosition.Y))

	return nextPosition
}

func (c *Core) greedyPathFind(layer LayerId, position image.Point, targetArea image.Rectangle, depth int, pathLength int) (image.Point, image.Point, int, int) {
	var (
		bestStep       = image.Point{}
		bestPosition   = position
		bestDistance   = DistanceToArea(position, targetArea)
		bestPathLength = pathLength
	)
	if bestDistance == 0 || depth == 0 {
		return bestStep, bestPosition, bestDistance, pathLength
	}

	steps := []image.Point{
		// Prioritize orthogonal steps
		{0, -1}, {1, 0}, {0, 1}, {-1, 0},
		{1, -1}, {1, 1}, {-1, 1}, {-1, -1},
	}
	for _, _step := range steps {
		_position := position.Add(_step)
		if !c.unitCanMoveTo(_position, layer) {
			continue
		}
		_, _, _distance, _pathLength := c.greedyPathFind(layer, _position, targetArea, depth-1, pathLength+1)

		replace := false
		if _distance < bestDistance {
			// Shorter distance
			replace = true
		} else if _distance == bestDistance {
			if _pathLength < bestPathLength {
				// Same distance but shorter path
				replace = true
			} else if ManhattanDistanceToArea(_position, targetArea) < ManhattanDistanceToArea(bestPosition, targetArea) {
				// Same distance and path length but closer to target in manhattan distance
				if _pathLength == bestPathLength {
					replace = true
				} else if bestStep.Eq(image.Point{0, 0}) && _pathLength == bestPathLength+1 {
					// New path length will be the same as not moving but manhattan distance is smaller
					replace = true
				}
			}
		}

		if replace {
			bestStep = _step
			bestPosition = _position
			bestDistance = _distance
			bestPathLength = _pathLength
			if depth == 1 && _distance == 0 {
				break
			}
		}
	}
	return bestStep, bestPosition, bestDistance, bestPathLength
}

func (c *Core) GetUnitToFireAt(obj UnitObjectWithRow) PlayerObjectMatchByDistance {
	var (
		playerId        = obj.PlayerId()
		unitId          = obj.ObjectId()
		fighter         = c.GetUnit(playerId, unitId)
		fighterPosition = GetPositionAsPoint(fighter)
		protoId         = fighter.GetUnitType()
		proto           = c.GetUnitPrototype(protoId)
	)
	matches := c.GetNearestEnemyUnits(playerId, fighterPosition, func(playerId, unitId uint8, unit *datamod.UnitsRow) bool {
		unitState := UnitState(unit.GetState())
		return unitState.HasSpawned() && !unitState.IsDeadOrInactive()
	})
	var targetMatch PlayerObjectMatchByDistance
	var targetAttackStrength uint8
	for _, layer := range []LayerId{LayerId_Land, LayerId_Air, LayerId_Hover} {
		match := matches[layer]
		if match.IsNil() {
			continue
		}
		if layer == LayerId_Hover && !targetMatch.IsNil() {
			// Hover units (workers) are the lowest priority
			continue
		}
		strength := GetAttackStrength(proto, layer)
		_range := proto.GetAttackRange()
		if match.Distance > int(_range) {
			continue
		}
		if strength == 0 {
			continue
		}
		if targetMatch.IsNil() || strength > targetAttackStrength {
			targetMatch = match
			targetAttackStrength = strength
		}
	}
	return targetMatch
}

func (c *Core) shootUnit(attackerObj UnitObjectWithRow, targetObj UnitObjectWithRow) {
	var (
		timeNow         = c.AbsSubTickIndex()
		target          = targetObj.Unit()
		targetProtoId   = target.GetUnitType()
		targetProto     = c.GetUnitPrototype(targetProtoId)
		targetIntegrity = target.GetIntegrity()
		layer           = LayerId(targetProto.GetLayer())
		attacker        = attackerObj.Unit()
		attackerProtoId = attacker.GetUnitType()
		attackerProto   = c.GetUnitPrototype(attackerProtoId)
		attackStrength  = GetAttackStrength(attackerProto, layer)
	)
	target.SetIntegrity(utils.SafeSubUint8(targetIntegrity, attackStrength))
	attacker.SetTimestamp(timeNow)
	c.emitInternalEvent(InternalEventId_Shot, &InternalEvent_Shot{
		Attacker: attackerObj.Object(),
		Target:   targetObj.Object(),
	})
}

func (c *Core) shootBuilding(attackerObj UnitObjectWithRow, targetObj BuildingObjectWithRow) {
	var (
		layer           = LayerId_Land
		target          = targetObj.Building()
		targetIntegrity = target.GetIntegrity()
		attacker        = attackerObj.Unit()
		attackerProtoId = attacker.GetUnitType()
		attackerProto   = c.GetUnitPrototype(attackerProtoId)
		attackStrength  = GetAttackStrength(attackerProto, layer)
	)
	integrity := utils.SafeSubUint8(targetIntegrity, attackStrength)
	target.SetIntegrity(integrity)
	attacker.SetTimestamp(c.AbsSubTickIndex())
	if integrity == 0 {
		c.setBuildingDestroyed(targetObj)
		// Target is destroyed, hold position
		attackerPosition := GetPositionAsPoint(attacker)
		attackerCommand := NewFighterCommandData(FighterCommandType_HoldPosition)
		attackerCommand.SetTargetPosition(attackerPosition)
		c.assignUnit(attackerObj, attackerCommand)
	}
	c.emitInternalEvent(InternalEventId_Shot, &InternalEvent_Shot{
		Attacker: attackerObj.Object(),
		Target:   targetObj.Object(),
	})
}

func (c *Core) tickUnitPreliminary(obj UnitObjectWithRow) bool {
	var (
		unit          = obj.Unit()
		unitState     = UnitState(unit.GetState())
		unitTimestamp = unit.GetTimestamp()
		protoId       = unit.GetUnitType()
		proto         = c.GetUnitPrototype(protoId)
		spawnTime     = proto.GetSpawnTime()
	)
	if unit.GetIsPreTicked() {
		// Unit was just (externally) assigned, skip this tick
		unit.SetIsPreTicked(false)
		return false
	}
	if unitState.IsSpawning() {
		// Spawn unit if spawn time has elapsed
		timeNow := c.AbsSubTickIndex()
		deltaTime := timeNow - unitTimestamp
		if deltaTime >= uint32(spawnTime) {
			c.setUnitSpawned(obj)
		}
		return false
	}
	if !unitState.IsAlive() || !unitState.IsPaid() {
		// Unit is dead or unpaid, skip this tick
		return false
	}
	if c.GetMainBuilding(obj.PlayerId()).GetIntegrity() == 0 {
		return false
	}
	if proto.GetIsWorker() {
		var (
			workerState       = UnitState(unit.GetState())
			workerCommand     = WorkerCommandData(unit.GetCommand())
			workerCommandType = workerCommand.Type()
			workerPosition    = GetPositionAsPoint(unit)
			workerLoad        = unit.GetLoad()
		)
		if workerState == UnitState_Inactive {
			// Activate worker if busy but inactive
			if workerCommandType.IsBusy() {
				tile := c.GetBoardTile(uint16(workerPosition.X), uint16(workerPosition.Y))
				layer := LayerId(proto.GetLayer())
				if IsTileEmpty(tile, layer) {
					SetTileUnit(tile, layer, obj.PlayerId(), obj.ObjectId())
					c.setUnitState(obj, UnitState_Active)
				}
			}
			return false
		} else if workerCommandType.IsIdle() &&
			workerPosition.Eq(c.GetWorkerPortPosition(obj.PlayerId())) &&
			workerLoad == 0 {
			// Deactivate worker if idle and at worker port with no load
			tile := c.GetBoardTile(uint16(workerPosition.X), uint16(workerPosition.Y))
			layer := LayerId(proto.GetLayer())
			EmptyUnitTile(tile, layer)
			c.setUnitState(obj, UnitState_Inactive)
			return false
		}
	}
	return true
}

func (c *Core) tickWorkerAction(obj UnitObjectWithRow) bool {
	var (
		worker            = obj.Unit()
		workerPosition    = GetPositionAsPoint(worker)
		workerCommand     = WorkerCommandData(worker.GetCommand())
		workerCommandType = workerCommand.Type()
		workerLoad        = uint16(worker.GetLoad())
	)
	if workerLoad > 0 {
		mainBuildingArea := c.GetMainBuildingArea(obj.PlayerId())
		if workerPosition.In(mainBuildingArea) {
			// If worker is carrying resources and at main building, unload
			player := c.GetPlayer(obj.PlayerId())
			computeDemand := player.GetComputeDemand()
			computeSupply := player.GetComputeSupply()
			penaltyDivisor := uint16(1)
			if computeDemand > computeSupply {
				penaltyDivisor = 2
			}
			addResource(player, workerLoad/penaltyDivisor)
			worker.SetLoad(0)
			return false
		}
		return true
	}
	if workerCommandType.IsIdle() {
		// If worker is idle, move on to the movement phase
		return true
	}

	var (
		targetPlayerId, targetBuildingId = workerCommand.TargetBuilding()
		targetBuilding                   = c.GetBuilding(targetPlayerId, targetBuildingId)
		targetPosition                   = GetPositionAsPoint(targetBuilding)
		targetProtoId                    = targetBuilding.GetBuildingType()
		targetProto                      = c.GetBuildingPrototype(targetProtoId)
		targetSize                       = GetDimensionsAsPoint(targetProto)
		targetArea                       = image.Rectangle{Min: targetPosition, Max: targetPosition.Add(targetSize)}
	)
	if workerCommandType == WorkerCommandType_Build &&
		BuildingState(targetBuilding.GetState()) != BuildingState_Building {
		// If assigned building is not in building state, set worker to idle
		c.setWorkerToIdle(obj)
		return true
	}
	if !workerPosition.In(targetArea) {
		// If worker is not at target building, move on to the movement phase
		return true
	}

	var (
		targetBuildingTimestamp = targetBuilding.GetTimestamp()
		timeNow                 = c.AbsSubTickIndex()
		deltaTime               = timeNow - targetBuildingTimestamp
	)
	if workerCommandType == WorkerCommandType_Gather {
		// If worker is gathering resources and at target building (mine), load resources
		if deltaTime < uint32(targetProto.GetMineTime()) {
			// No resources to load yet
			return false
		}
		// Resources are ready, load them
		targetBuilding.SetTimestamp(timeNow)
		worker.SetLoad(targetProto.GetResourceMine())
		return false
	} else if workerCommandType == WorkerCommandType_Build {
		// If worker is building, at target building, and building time has elapsed, build

		if targetBuildingTimestamp == 0 {
			// Building has no timestamp, set it to current time -1
			targetBuilding.SetTimestamp(timeNow - 1)
			deltaTime = 0
			return false
		}
		if deltaTime < uint32(targetProto.GetBuildingTime()) {
			// Building time has not elapsed yet
			return false
		}
		// Building time has elapsed, build
		c.setBuildingBuilt(c.GetBuildingObject(targetPlayerId, targetBuildingId))

		// Building is built, idle
		c.setWorkerToIdle(obj)
		return false
	}

	return true
}

func (c *Core) tickFighterAction(obj UnitObjectWithRow) bool {
	var (
		fighter         = obj.Unit()
		fighterPosition = GetPositionAsPoint(fighter)
		fighterProtoId  = fighter.GetUnitType()
		fighterProto    = c.GetUnitPrototype(fighterProtoId)
		fighterCommand  = FighterCommandData(fighter.GetCommand())
	)

	deltaTime := c.AbsSubTickIndex() - fighter.GetTimestamp()
	if deltaTime < uint32(fighterProto.GetAttackCooldown()) {
		// Cooldown has not elapsed, move on to the movement phase
		// Unit will move even if an enemy is in range if the cooldown has not elapsed
		return true
	}

	// Shoot an enemy unit, if in range
	targetMatch := c.GetUnitToFireAt(obj)
	if !targetMatch.IsNil() {
		targetUnit := c.GetUnitObject(targetMatch.PlayerId, targetMatch.ObjectId)
		c.shootUnit(obj, targetUnit)

		if fighterProto.GetIsAssault() {
			// Assault units can move and attack in the same tick so they move on to
			// the movement phase
			return true
		} else {
			return false
		}
	}

	if !fighterCommand.Type().IsTargetingBuilding() {
		// Nothing to shoot at, move on to the movement phase
		return true
	}

	var (
		targetPlayerId, targetBuildingId = fighterCommand.TargetBuilding()
		targetBuildingObj                = c.GetBuildingObject(targetPlayerId, targetBuildingId)
		targetBuilding                   = targetBuildingObj.Building()
	)
	if BuildingState(targetBuilding.GetState()) == BuildingState_Destroyed {
		// Target building is destroyed, hold position
		fighterCommand = NewFighterCommandData(FighterCommandType_HoldPosition)
		fighterCommand.SetTargetPosition(fighterPosition)
		c.assignUnit(obj, fighterCommand)
		return false
	}

	var (
		targetPosition   = GetPositionAsPoint(targetBuilding)
		targetProtoId    = targetBuilding.GetBuildingType()
		targetProto      = c.GetBuildingPrototype(targetProtoId)
		targetSize       = GetDimensionsAsPoint(targetProto)
		targetArea       = image.Rectangle{Min: targetPosition, Max: targetPosition.Add(targetSize)}
		distanceToTarget = DistanceToArea(fighterPosition, targetArea)
	)
	if distanceToTarget > int(fighterProto.GetAttackRange()) {
		// Not in range, move on to the movement phase
		return true
	}
	// Target building in range, shoot
	c.shootBuilding(obj, targetBuildingObj)

	return false
}

func (c *Core) tickUnitIntermediate(obj UnitObjectWithRow) bool {
	var (
		unit          = obj.Unit()
		unitState     = UnitState(unit.GetState())
		unitIntegrity = unit.GetIntegrity()
	)
	if !unitState.IsAlive() {
		// Unit is not alive, skip this tick
		return false
	}
	if unitIntegrity == 0 {
		c.setUnitDead(obj)
		return false
	}
	return true
}

func (c *Core) tickWorkerMovement(obj UnitObjectWithRow) bool {
	var (
		playerId       = obj.PlayerId()
		worker         = obj.Unit()
		workerPosition = GetPositionAsPoint(worker)
		workerCommand  = WorkerCommandData(worker.GetCommand())
		workerLoad     = uint16(worker.GetLoad())
	)
	var targetArea image.Rectangle

	if workerCommand.Type() == WorkerCommandType_Idle {
		// Worker is idle, move to worker port
		targetArea = c.GetWorkerPortArea(playerId)
	} else if workerLoad > 0 {
		// Worker is carrying resources, move to main building
		targetArea = c.GetMainBuildingArea(playerId)
	} else {
		targetPlayerId, targetBuildingId := workerCommand.TargetBuilding()
		targetArea = c.GetBuildingArea(targetPlayerId, targetBuildingId)
	}

	if workerPosition.In(targetArea) {
		// Already in target area, move on to the next phase (do nothing)
		// NOTE: all workers that get to the movement phase are not in their target area so
		// this should never happen
		return true
	}
	c.moveUnitToTarget(obj, targetArea)

	return false
}

func (c *Core) tickFighterMovement(obj UnitObjectWithRow) bool {
	var (
		fighter         = obj.Unit()
		fighterPosition = GetPositionAsPoint(fighter)
		fighterCommand  = FighterCommandData(fighter.GetCommand())
		fighterProtoId  = fighter.GetUnitType()
		fighterProto    = c.GetUnitPrototype(fighterProtoId)
	)
	var targetArea image.Rectangle

	if fighterCommand.Type().IsTargetingBuilding() {
		targetPlayerId, targetBuildingId := fighterCommand.TargetBuilding()
		targetArea = c.GetBuildingArea(targetPlayerId, targetBuildingId)
	} else {
		targetPosition := fighterCommand.TargetPosition()
		targetSize := image.Point{1, 1}
		targetArea = image.Rectangle{Min: targetPosition, Max: targetPosition.Add(targetSize)}
	}

	if fighterPosition.In(targetArea) {
		// Already in target area, move on to the next phase (do nothing)
		return true
	}

	if fighterCommand.Type().IsTargetingBuilding() {
		if DistanceToArea(fighterPosition, targetArea) <= int(fighterProto.GetAttackRange()) {
			// In range, move on to the next phase (do nothing)
			return true
		}
	}

	c.moveUnitToTarget(obj, targetArea)

	return false
}

func (c *Core) tickWorker(obj UnitObjectWithRow) {
	var pre, act, inter bool
	pre = c.tickUnitPreliminary(obj)
	if pre {
		act = c.tickWorkerAction(obj)
	}
	inter = c.tickUnitIntermediate(obj)
	if act && inter {
		c.tickWorkerMovement(obj)
	}
}

func (c *Core) tickFighter(obj UnitObjectWithRow) {
	var pre, act, inter bool
	pre = c.tickUnitPreliminary(obj)
	if pre {
		act = c.tickFighterAction(obj)
	}
	inter = c.tickUnitIntermediate(obj)
	if act && inter {
		c.tickFighterMovement(obj)
	}
}

func (c *Core) tickUnit(obj UnitObjectWithRow) {
	unit := obj.Unit()
	protoId := unit.GetUnitType()
	proto := c.GetUnitPrototype(protoId)
	if proto.GetIsWorker() {
		c.tickWorker(obj)
	} else {
		c.tickFighter(obj)
	}
}

func (c *Core) Tick() {
	if !c.IsInitialized() {
		return
	}
	if !c.HasStarted() {
		return
	}

	nPlayers := c.GetMeta().GetPlayerCount()
	if nPlayers == 0 {
		return
	}
	for playerId := uint8(1); playerId < nPlayers+1; playerId++ {
		c.tickPlayer(playerId)
	}

	// Rotate starting player
	startingPlayerIdx := uint8(c.AbsSubTickIndex() % uint32(nPlayers))

	prePassers := make([][]uint8, nPlayers)
	interPassers := make([][]uint8, nPlayers)
	passedAct := map[uint8]map[uint8]struct{}{}

	for ii := uint8(0); ii < nPlayers; ii++ {
		prePassers[ii] = []uint8{}
		interPassers[ii] = []uint8{}
		passedAct[ii+1] = make(map[uint8]struct{})
	}

	// Preliminary phase [all]
	for ii := uint8(0); ii < nPlayers; ii++ {
		playerIdx := (startingPlayerIdx + ii) % nPlayers
		playerId := playerIdx + 1
		nUnits := c.GetPlayer(playerId).GetUnitCount()
		for unitId := uint8(1); unitId < nUnits+1; unitId++ {
			ok := c.tickUnitPreliminary(c.GetUnitObject(playerId, unitId))
			if ok {
				prePassers[playerIdx] = append(prePassers[playerIdx], unitId)
			}
		}
	}

	// Action phase [fighters that passed preliminary phase]
	for ii := uint8(0); ii < nPlayers; ii++ {
		playerIdx := (startingPlayerIdx + ii) % nPlayers
		playerId := playerIdx + 1
		for _, unitId := range prePassers[playerIdx] {
			unit := c.GetUnit(playerId, unitId)
			unitProtoId := unit.GetUnitType()
			unitProto := c.GetUnitPrototype(unitProtoId)
			if unitProto.GetIsWorker() {
				continue
			}
			ok := c.tickFighterAction(c.GetUnitObject(playerId, unitId))
			if ok {
				passedAct[playerId][unitId] = struct{}{}
			}
		}
	}

	// Action phase [workers that passed preliminary phase]
	for ii := uint8(0); ii < nPlayers; ii++ {
		playerIdx := (startingPlayerIdx + ii) % nPlayers
		playerId := playerIdx + 1
		for _, unitId := range prePassers[playerIdx] {
			unit := c.GetUnit(playerId, unitId)
			unitProtoId := unit.GetUnitType()
			unitProto := c.GetUnitPrototype(unitProtoId)
			if !unitProto.GetIsWorker() {
				continue
			}
			ok := c.tickWorkerAction(c.GetUnitObject(playerId, unitId))
			if ok {
				passedAct[playerId][unitId] = struct{}{}
			}
		}
	}

	// Intermediate phase [all]
	for ii := uint8(0); ii < nPlayers; ii++ {
		playerIdx := (startingPlayerIdx + ii) % nPlayers
		playerId := playerIdx + 1
		nUnits := c.GetPlayer(playerId).GetUnitCount()
		for unitId := uint8(1); unitId < nUnits+1; unitId++ {
			ok := c.tickUnitIntermediate(c.GetUnitObject(playerId, unitId))
			if ok {
				interPassers[playerIdx] = append(interPassers[playerIdx], unitId)
			}
		}
	}

	// Movement phase [fighters that passed action and intermediate phase]
	for ii := uint8(0); ii < nPlayers; ii++ {
		playerIdx := (startingPlayerIdx + ii) % nPlayers
		playerId := playerIdx + 1
		for _, unitId := range interPassers[playerIdx] {
			if _, ok := passedAct[playerId][unitId]; !ok {
				continue
			}
			unit := c.GetUnit(playerId, unitId)
			unitProtoId := unit.GetUnitType()
			unitProto := c.GetUnitPrototype(unitProtoId)
			if unitProto.GetIsWorker() {
				continue
			}
			c.tickFighterMovement(c.GetUnitObject(playerId, unitId))
		}
	}

	// Movement phase [workers that passed action and intermediate phase]
	for ii := uint8(0); ii < nPlayers; ii++ {
		playerIdx := (startingPlayerIdx + ii) % nPlayers
		playerId := playerIdx + 1
		for _, unitId := range interPassers[playerIdx] {
			if _, ok := passedAct[playerId][unitId]; !ok {
				continue
			}
			unit := c.GetUnit(playerId, unitId)
			unitProtoId := unit.GetUnitType()
			unitProto := c.GetUnitPrototype(unitProtoId)
			if !unitProto.GetIsWorker() {
				continue
			}
			c.tickWorkerMovement(c.GetUnitObject(playerId, unitId))
		}
	}
}

func (c *Core) IsInitialized() bool {
	return c.GetMeta().GetIsInitialized()
}

func (c *Core) HasStarted() bool {
	return c.GetMeta().GetHasStarted()
}

func (c *Core) Initialize(action *Initialization) error {
	if c.IsInitialized() {
		return ErrAlreadyInitialized
	}
	bn := c.BlockNumber()
	c.GetMeta().Set(action.Width, action.Height, 0, 0, 0, true, false, uint32(bn))
	return nil
}

func (c *Core) AddPlayer(action *PlayerAddition) error {
	if c.HasStarted() {
		return ErrAlreadyStarted
	}
	nPlayers := c.GetMeta().GetPlayerCount()
	playerId := utils.SafeAddUint8(nPlayers, 1)
	if playerId == NilPlayerId {
		return ErrPlayerLimitReached
	}
	c.GetMeta().SetPlayerCount(playerId)
	player := c.GetPlayer(playerId)
	player.SetSpawnAreaX(action.SpawnAreaX)
	player.SetSpawnAreaY(action.SpawnAreaY)
	player.SetWorkerPortX(action.WorkerPortX)
	player.SetWorkerPortY(action.WorkerPortY)
	player.SetBuildingPayQueuePointer(1)
	player.SetBuildingBuildQueuePointer(1)
	player.SetUnitPayQueuePointer(1)

	return nil
}

func (c *Core) Start(action *Start) error {
	_ = action
	if c.HasStarted() {
		return ErrAlreadyStarted
	}
	c.GetMeta().SetHasStarted(true)
	return nil
}

func (c *Core) CreateUnit(action *UnitCreation) error {
	if !c.IsInitialized() {
		return ErrNotInitialized
	}
	var (
		playerId = action.PlayerId
		protoId  = action.UnitType
	)
	if err := c.ValidatePlayerId(playerId); err != nil {
		return err
	}
	if err := c.ValidateUnitType(protoId); err != nil {
		return err
	}
	if c.GetMainBuilding(playerId).GetIntegrity() == 0 {
		return ErrMainBuildingDestroyed
	}
	unitId := c.createUnit(playerId, protoId)
	if unitId == NilUnitId {
		return ErrUnitLimitReached
	}
	if !c.HasStarted() {
		// Spawn unit immediately during genesis phase
		obj := c.GetUnitObject(playerId, unitId)
		c.setUnitSpawning(obj)
		c.setUnitSpawned(obj)
	}
	return nil
}

func (c *Core) AssignUnit(action *UnitAssignation) error {
	if !c.IsInitialized() {
		return ErrNotInitialized
	}
	if !c.HasStarted() {
		return ErrNotStarted
	}
	var (
		playerId = action.PlayerId
		unitId   = action.UnitId
		unit     = c.GetUnit(playerId, unitId)
		protoId  = unit.GetUnitType()
		proto    = c.GetUnitPrototype(protoId)
	)
	if err := c.ValidatePlayerId(playerId); err != nil {
		return err
	}
	if c.GetMainBuilding(playerId).GetIntegrity() == 0 {
		return ErrMainBuildingDestroyed
	}
	if UnitState(unit.GetState()) == UnitState_Dead {
		return ErrUnitDead
	}
	if proto.GetIsWorker() {
		command := WorkerCommandData(action.Command)
		commandType := command.Type()
		if commandType.IsIdle() {
			// Make sure the command is canonical
			command = NewWorkerCommandData(WorkerCommandType_Idle)
		} else {
			var (
				targetPlayerId   = command.TargetPlayerId()
				targetBuildingId = command.TargetBuildingId()
				targetBuilding   = c.GetBuilding(targetPlayerId, targetBuildingId)
			)
			if commandType == WorkerCommandType_Gather {
				var (
					targetProtoId = targetBuilding.GetBuildingType()
					targetProto   = c.GetBuildingPrototype(targetProtoId)
				)
				if !targetProto.GetIsEnvironment() {
					return errors.New("target must be environment")
				}
			} else if commandType == WorkerCommandType_Build {
				if targetPlayerId != playerId {
					return errors.New("target must be self")
				}
				if BuildingState(targetBuilding.GetState()) != BuildingState_Building {
					return errors.New("target must be in building state")
				}
			} else {
				return errors.New("command not assignable")
			}
		}
		c.assignUnitExternal(c.GetUnitObject(playerId, unitId), command, &CommandPath{})
	} else {
		command := FighterCommandData(action.Command)
		commandType := command.Type()
		if commandType.IsTargetingBuilding() {
			targetPlayerId := command.TargetPlayerId()
			targetBuildingId := command.TargetBuildingId()
			if err := c.ValidatePlayerId(targetPlayerId); err != nil {
				return err
			}
			if err := c.ValidateBuildingId(targetPlayerId, targetBuildingId); err != nil {
				return err
			}
			if targetPlayerId == playerId {
				return errors.New("target must be an enemy")
			}
			targetBuilding := c.GetBuilding(targetPlayerId, targetBuildingId)
			targetProtoId := targetBuilding.GetBuildingType()
			targetProto := c.GetBuildingPrototype(targetProtoId)
			if targetProto.GetIsEnvironment() {
				return errors.New("target must be player-owned")
			}
			targetBuildingState := BuildingState(targetBuilding.GetState())
			if targetBuildingState == BuildingState_Destroyed {
				return errors.New("target must not be destroyed")
			}
			if targetBuildingState == BuildingState_Unpaid {
				return errors.New("target must not be unpaid")
			}
		} else if commandType.IsTargetingPosition() {
			targetPosition := command.TargetPosition()
			if !c.TileIsInBoard(targetPosition) {
				return errors.New("target must be in board")
			}
		} else {
			return errors.New("command not assignable")
		}
		path := NewCommandPath(action.CommandExtra, action.CommandMeta)
		c.assignUnitExternal(c.GetUnitObject(playerId, unitId), command, path)
	}
	return nil
}

func (c *Core) PlaceBuilding(action *BuildingPlacement) error {
	// fmt.Println("PlaceBuilding!!!")

	if !c.IsInitialized() {
		return ErrNotInitialized
	}
	var (
		playerId = action.PlayerId
		protoId  = action.BuildingType
		proto    = c.GetBuildingPrototype(protoId)
		size     = GetDimensionsAsPoint(proto)
		position = image.Point{int(action.X), int(action.Y)}
	)
	if err := c.ValidatePlayerOrEnvironmentId(playerId); err != nil {
		return err
	}
	if err := c.ValidateBuildingType(protoId); err != nil {
		return err
	}
	if proto.GetIsEnvironment() {
		if playerId != NilPlayerId {
			return errors.New("environment buildings must be assigned to nil player")
		}
	} else {
		if playerId == NilPlayerId {
			return errors.New("player buildings must be assigned to a player")
		}
		if c.GetPlayer(playerId).GetBuildingCount() != 0 {
			if c.GetMainBuilding(playerId).GetIntegrity() == 0 {
				return ErrMainBuildingDestroyed
			}
		}
	}
	buildArea := image.Rectangle{Min: position, Max: position.Add(size)}
	if !c.IsBuildableArea(buildArea) {
		return ErrAreaNotBuildable
	}
	buildingId := c.placeBuilding(playerId, protoId, position)
	if buildingId == NilBuildingId {
		return ErrBuildingLimitReached
	}
	if !c.HasStarted() || proto.GetIsEnvironment() {
		// Build building immediately during genesis phase
		obj := c.GetBuildingObject(playerId, buildingId)
		c.setBuildingBuilding(obj)
		c.setBuildingBuilt(obj)
	}
	return nil
}

func (c *Core) AddUnitPrototype(action *UnitPrototypeAddition) error {
	if !c.IsInitialized() {
		return ErrNotInitialized
	}
	nUnitPrototypes := c.GetMeta().GetUnitPrototypeCount()
	protoId := utils.SafeAddUint8(nUnitPrototypes, 1)
	if protoId == nUnitPrototypes {
		return ErrUnitPrototypeLimitReached
	}
	c.GetMeta().SetUnitPrototypeCount(protoId)
	proto := c.GetUnitPrototype(protoId)
	proto.Set(
		action.Layer,
		action.ResourceCost,
		action.ComputeCost,
		action.SpawnTime,
		action.MaxIntegrity,
		action.LandStrength,
		action.HoverStrength,
		action.AirStrength,
		action.AttackRange,
		action.AttackCooldown,
		action.IsAssault,
		action.IsWorker,
	)
	return nil
}

func (c *Core) AddBuildingPrototype(action *BuildingPrototypeAddition) error {
	if !c.IsInitialized() {
		return ErrNotInitialized
	}
	nBuildingPrototypes := c.GetMeta().GetBuildingPrototypeCount()
	protoId := utils.SafeAddUint8(nBuildingPrototypes, 1)
	if protoId == nBuildingPrototypes {
		return ErrBuildingPrototypeLimitReached
	}
	c.GetMeta().SetBuildingPrototypeCount(protoId)
	proto := c.GetBuildingPrototype(protoId)
	proto.Set(
		action.Width,
		action.Height,
		action.ResourceCost,
		action.ResourceCapacity,
		action.ComputeCapacity,
		action.ResourceMine,
		action.MineTime,
		action.MaxIntegrity,
		action.BuildingTime,
		action.IsArmory,
		action.IsEnvironment,
	)
	return nil
}
