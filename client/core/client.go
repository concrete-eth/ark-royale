package core

import (
	"fmt"
	"image"

	gen_utils "github.com/concrete-eth/archetype/utils"
	"github.com/concrete-eth/ark-rts/client/assets"
	"github.com/concrete-eth/ark-rts/gogen/datamod"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type UpdatableWithClient interface {
	Update(c *Client)
}

type DrawableWithClient interface {
	Draw(c *Client, screen *ebiten.Image)
}

type KeyFunction int

const (
	KeyFunction_Quit KeyFunction = iota
	KeyFunction_SetPath
	KeyFunction_Up
	KeyFunction_Down
	KeyFunction_Left
	KeyFunction_Right
	KeyFunction_CenterCamera
	KeyFunction_Base
	KeyFunction_ZoomIn
	KeyFunction_ZoomOut
	KeyFunction_ToggleTargetLines
	KeyFunction_ToggleDebugInfo
	KeyFunction_Deselect
)

type KeyMap map[KeyFunction][]ebiten.Key

func (k KeyMap) IsJustPressed(f KeyFunction) bool {
	for _, key := range k[f] {
		if inpututil.IsKeyJustPressed(key) {
			return true
		}
	}
	return false
}

func (k KeyMap) IsPressed(f KeyFunction) bool {
	for _, key := range k[f] {
		if ebiten.IsKeyPressed(key) {
			return true
		}
	}
	return false
}

func (k KeyMap) IsJustReleased(f KeyFunction) bool {
	for _, key := range k[f] {
		if inpututil.IsKeyJustReleased(key) {
			return true
		}
	}
	return false
}

func (k KeyMap) IsPressedWithShift(f KeyFunction) bool {
	return k.IsPressed(f) && ebiten.IsKeyPressed(ebiten.KeyShift)
}

var DefaultKeyMap = KeyMap{
	// KeyFunction_Quit:              {ebiten.KeyEscape},
	KeyFunction_SetPath:           {ebiten.KeyShift},
	KeyFunction_Up:                {ebiten.KeyW, ebiten.KeyUp},
	KeyFunction_Down:              {ebiten.KeyS, ebiten.KeyDown},
	KeyFunction_Left:              {ebiten.KeyA, ebiten.KeyLeft},
	KeyFunction_Right:             {ebiten.KeyD, ebiten.KeyRight},
	KeyFunction_CenterCamera:      {ebiten.KeyC},
	KeyFunction_Base:              {ebiten.KeyV},
	KeyFunction_ZoomIn:            {ebiten.KeyEqual, ebiten.KeyKPAdd, ebiten.KeyE},
	KeyFunction_ZoomOut:           {ebiten.KeyMinus, ebiten.KeyKPSubtract, ebiten.KeyQ},
	KeyFunction_ToggleTargetLines: {ebiten.KeyL},
	KeyFunction_ToggleDebugInfo:   {ebiten.KeyF3, ebiten.KeyK},
	KeyFunction_Deselect:          {ebiten.KeyEscape},
}

// Selection holds the current selection state.
type Selection struct {
	BuildingType uint8
	Units        []uint8
}

// Clear clears the selection.
func (s *Selection) Clear() {
	s.BuildingType = 0
	s.Units = []uint8{}
}

// Main game client object run by ebiten.RunGame.
type Client struct {
	coreRenderer *CoreRenderer

	keyMap KeyMap // Key map

	hud *HudSet // HUD component set

	lastClickedScreenPosition image.Point // Last left clicked screen position
	selected                  Selection   // Current selection
	selecting                 bool        // Draw selection box

	commandPath []image.Point // Command path

	onSelectionChange func() // On selection change callback

	active bool // Actively update state
}

var _ ebiten.Game = (*Client)(nil)

// Instantiates a new Client object.
// Parameters:
// * headlessClient: the headless client to use
// * hinter: the hinter to use
// Returns:
// * the new Client object
// The headless client keeps the local game state in sync with the backend.
// The hinter provides a list of actions expected to be executed in the next block so the can
// be anticipated by the client.
func NewClient(coreRenderer *CoreRenderer, hud *HudSet, active bool) *Client {
	return &Client{
		coreRenderer: coreRenderer,
		hud:          hud,
		keyMap:       DefaultKeyMap,
		active:       active,
	}
}

func (c *Client) CoreRenderer() *CoreRenderer {
	return c.coreRenderer
}

func (c *Client) Hud() *HudSet {
	return c.hud
}

func (c *Client) Game() *rts.Core {
	return c.coreRenderer.Game()
}

func (c *Client) Headless() IHeadlessClient {
	return c.coreRenderer.IHeadlessClient
}

func (c *Client) PlayerId() uint8 {
	return c.coreRenderer.PlayerId()
}

func (c *Client) SetKeyMap(keyMap KeyMap) {
	c.keyMap = keyMap
}

func (c *Client) SetOnSelectionChange(onSelectionChange func()) {
	c.onSelectionChange = onSelectionChange
}

// Clears the selection and hides the building shadows and buildable area.
func (c *Client) ClearSelection() {
	c.selected.Clear()
	c.clearCommandPath()

	if c.onSelectionChange != nil {
		c.onSelectionChange()
	}
}

// Clears the current selection, selects a given building type, and executes the
// selection side effects: show building shadows and buildable area.
func (c *Client) SelectBuildingType(buildingType uint8) {
	c.ClearSelection()
	c.selected.BuildingType = buildingType
	c.toggleShowBuildingsShadow(true)

	if c.onSelectionChange != nil {
		c.onSelectionChange()
	}
}

// Returns the selected building type.
func (c *Client) SelectedBuildingType() uint8 {
	return c.selected.BuildingType
}

// Returns true if a building type is selected.
func (c *Client) IsSelectingBuildingType() bool {
	return c.selected.BuildingType != 0
}

// Selects a given slice of unit ids.
func (c *Client) SelectUnits(unitIds ...uint8) {
	c.ClearSelection()
	c.selected.Units = unitIds

	if c.onSelectionChange != nil {
		c.onSelectionChange()
	}
}

// Returns the selected unit ids.
func (c *Client) SelectedUnits() []uint8 {
	return c.selected.Units
}

// Returns true if any unit ids are selected.
func (c *Client) IsSelectingUnits() bool {
	return len(c.selected.Units) > 0
}

func (c *Client) addCommandPathPoint(point image.Point) {
	c.commandPath = append(c.commandPath, point)
}

func (c *Client) clearCommandPath() {
	c.commandPath = make([]image.Point, 0)
}

func (c *Client) getCommandPath() []image.Point {
	return c.commandPath
}

// Toggles the visibility of the building shadows.
func (c *Client) toggleShowBuildingsShadow(show bool) {
	buildingLayer := c.coreRenderer.Layers().Layer(LayerName_Buildings)
	for _, sprite := range buildingLayer.Sprites() {
		sprite.ToggleShadow(show)
	}
}

// Run a function for every selected unit.
func (c *Client) forEachSelectedUnit(forEach func(unitId uint8, unit *datamod.UnitsRow)) {
	for _, unitId := range c.SelectedUnits() {
		unit := c.Game().GetUnit(c.PlayerId(), unitId)
		forEach(unitId, unit)
	}
}

func (c *Client) IterSelectedUnits(filters ...rts.UnitFilter) *rts.Iterator_uint8 {
	return &rts.Iterator_uint8{
		Current: 0,
		Max:     uint8(len(c.SelectedUnits())),
		Get: func(unitIdxP1 uint8) (uint8, interface{}) {
			unitId := c.SelectedUnits()[unitIdxP1-1]
			unit := c.Game().GetUnit(c.PlayerId(), unitId)
			for _, filter := range filters {
				if !filter(c.PlayerId(), unitId, unit) {
					return 0, nil
				}
			}
			return unitId, unit
		},
	}
}

func (c *Client) handleInput() {
	if c.keyMap.IsJustPressed(KeyFunction_Deselect) {
		c.ClearSelection()
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		c.ClearSelection()
	}

	pathKeyDown := c.keyMap.IsPressed(KeyFunction_SetPath)
	if c.keyMap.IsJustReleased(KeyFunction_SetPath) {
		path := c.getCommandPath()
		if len(path) == 0 {
			return
		}
		position := path[len(path)-1]
		path = path[:len(path)-1]
		command := rts.NewFighterCommandData(rts.FighterCommandType_HoldPosition)
		command.SetTargetPosition(position)
		c.assignSelectedFighters(command, path)
		c.clearCommandPath()
		// return
	}

	cursorScreenPosition := image.Pt(ebiten.CursorPosition())
	selectionScreenRect := image.Rectangle{Min: c.lastClickedScreenPosition, Max: cursorScreenPosition}.Canon()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.selecting = false
		c.lastClickedScreenPosition = cursorScreenPosition
		return
	} else {
		if selectionScreenRect.Dx() > 4 || selectionScreenRect.Dy() > 4 {
			c.selecting = true
		}
	}

	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}

	if c.selecting {
		c.ClearSelection()
		if selectionScreenRect.Dx() < c.coreRenderer.tileDisplaySize/2 && selectionScreenRect.Dy() < c.coreRenderer.tileDisplaySize/2 {
			tilePosition := c.coreRenderer.ScreenCoordToTileCoord(cursorScreenPosition)
			tile := c.Game().GetBoardTile(uint16(tilePosition.X), uint16(tilePosition.Y))
			var unitId uint8
			if tile.GetLandObjectType() == rts.ObjectType_Unit.Uint8() {
				unitId = tile.GetLandObjectId()
			} else if tile.GetHoverPlayerId() == c.PlayerId() {
				unitId = tile.GetHoverUnitId()
			} else if tile.GetAirPlayerId() == c.PlayerId() {
				unitId = tile.GetAirUnitId()
			}
			if unitId != rts.NilUnitId {
				unit := c.Game().GetUnit(c.PlayerId(), unitId)
				unitState := rts.UnitState(unit.GetState())
				if !unitState.IsDeadOrInactive() {
					c.SelectUnits(unitId)
				}
			}
		} else {
			selectedUnits := make([]uint8, 0)
			c.Game().ForEachUnit(c.PlayerId(), func(unitId uint8, unit *datamod.UnitsRow) {
				unitState := rts.UnitState(unit.GetState())
				if unitState.IsDeadOrInactive() {
					return
				}
				unitCenterScreenPosition := c.coreRenderer.GetUnitScreenPosition(c.PlayerId(), unitId).
					Add(image.Point{c.coreRenderer.tileDisplaySize / 2, c.coreRenderer.tileDisplaySize / 2})
				if unitCenterScreenPosition.In(selectionScreenRect) {
					selectedUnits = append(selectedUnits, unitId)
				}
			})
			c.SelectUnits(selectedUnits...)
		}
		return
	}

	if !cursorScreenPosition.In(c.coreRenderer.boardDisplayRect) {
		// Clicked outside the board
		c.clearCommandPath()
		c.selected.Units = []uint8{}
		return
	}

	var isSelectingFighters bool
	// Check if any fighters are selected
	for _, unitId := range c.SelectedUnits() {
		var (
			unit      = c.Game().GetUnit(c.PlayerId(), unitId)
			protoId   = unit.GetUnitType()
			proto     = c.Game().GetUnitPrototype(protoId)
			unitState = rts.UnitState(unit.GetState())
		)
		if !proto.GetIsWorker() && unitState != rts.UnitState_Dead {
			isSelectingFighters = true
			break
		}
	}
	isSettingPath := pathKeyDown && isSelectingFighters

	var (
		tilePosition   = c.coreRenderer.ScreenCoordToTileCoord(cursorScreenPosition)
		tile           = c.Game().GetBoardTile(uint16(tilePosition.X), uint16(tilePosition.Y))
		tilePlayerId   = tile.GetLandPlayerId()
		tileObjectType = rts.ObjectType(tile.GetLandObjectType())
		tileBuildingId = tile.GetLandObjectId()
	)

	if !tilePosition.In(c.Game().BoardRect()) {
		c.ClearSelection()
		return
	}

	if tilePlayerId == rts.NilPlayerId && tileObjectType == rts.ObjectType_Building {
		// Neutral building (mine)
		var (
			building         = c.Game().GetBuilding(tilePlayerId, tileBuildingId)
			buildingPosition = rts.GetPositionAsPoint(building)
			protoId          = building.GetBuildingType()
			proto            = c.Game().GetBuildingPrototype(protoId)
		)
		if proto.GetResourceMine() > 0 {
			command := rts.NewWorkerCommandData(rts.WorkerCommandType_Gather)
			command.SetTargetBuilding(rts.NilPlayerId, tileBuildingId)

			if len(c.SelectedUnits()) == 0 {
				c.ClearSelection()
				// Assign nearest idle worker
				buildingSize := rts.GetDimensionsAsPoint(proto)
				targetArea := image.Rectangle{Min: buildingPosition, Max: buildingPosition.Add(buildingSize)}
				iter := c.Game().IterUnits(c.PlayerId(), func(playerId, unitId uint8, unit *datamod.UnitsRow) bool {
					var (
						unitProtoId = unit.GetUnitType()
						unitProto   = c.Game().GetUnitPrototype(unitProtoId)
						unitState   = rts.UnitState(unit.GetState())
						unitCommand = rts.WorkerCommandData(unit.GetCommand())
					)
					return unitProto.GetIsWorker() &&
						unitState.HasSpawned() &&
						unitState.IsAlive() &&
						unitCommand.Type() == rts.WorkerCommandType_Idle
				})
				match := c.Game().GetNearest(iter, targetArea)
				workerId := match.ObjectId
				if workerId != rts.NilUnitId {
					c.Headless().AssignUnit(workerId, command)
				} else {
					// Assign spawning worker if no spawned worker is available
					iter := c.Game().IterUnits(c.PlayerId(), func(playerId, unitId uint8, unit *datamod.UnitsRow) bool {
						var (
							unitProtoId = unit.GetUnitType()
							unitProto   = c.Game().GetUnitPrototype(unitProtoId)
							unitState   = rts.UnitState(unit.GetState())
							unitCommand = rts.WorkerCommandData(unit.GetCommand())
						)
						return unitProto.GetIsWorker() &&
							unitState.IsSpawning() &&
							unitCommand.Type() == rts.WorkerCommandType_Idle
					})
					match := c.Game().GetNearest(iter, targetArea)
					workerId := match.ObjectId
					if workerId != rts.NilUnitId {
						c.Headless().AssignUnit(workerId, command)
					}
				}
			} else {
				// Assign selected workers
				unitIds := make([]uint8, 0)
				c.forEachSelectedUnit(func(unitId uint8, unit *datamod.UnitsRow) {
					var (
						unitProtoId = unit.GetUnitType()
						unitProto   = c.Game().GetUnitPrototype(unitProtoId)
						unitState   = rts.UnitState(unit.GetState())
					)
					if !unitProto.GetIsWorker() || unitState == rts.UnitState_Dead {
						return
					}
					unitIds = append(unitIds, unitId)
				})
				c.Headless().AssignUnits(unitIds, command)
			}
		}
		return
	}

	if tilePlayerId == rts.NilPlayerId || tileObjectType != rts.ObjectType_Building {
		// Empty tile
		if c.IsSelectingBuildingType() {
			// Build building
			var (
				snapTilePosition = tilePosition.Div(2).Mul(2)
			)
			if snapTilePosition.In(c.Game().BoardRect()) {
				c.Headless().PlaceBuilding(c.SelectedBuildingType(), snapTilePosition)
				c.ClearSelection()
			}
			return
		} else if isSettingPath && len(c.getCommandPath()) < 4 {
			c.addCommandPathPoint(tilePosition)
			return
		} else if isSelectingFighters {
			command := rts.NewFighterCommandData(rts.FighterCommandType_HoldPosition)
			command.SetTargetPosition(tilePosition)
			c.assignSelectedFighters(command, c.getCommandPath())
			c.clearCommandPath()
			return
		} else {
			c.ClearSelection()
		}
	}

	if tilePlayerId == c.PlayerId() && tileObjectType == rts.ObjectType_Building {
		// Own building
		tileBuilding := c.Game().GetBuilding(tilePlayerId, tileBuildingId)
		if tileBuildingId == 1 {
			// Main building
			// Command workers to idle
			c.assignSelectedWorkers(rts.NewWorkerCommandData(rts.WorkerCommandType_Idle), false)
		} else if rts.BuildingState(tileBuilding.GetState()) == rts.BuildingState_Building {
			// Building not main
			// Building not built
			command := rts.NewWorkerCommandData(rts.WorkerCommandType_Build)
			command.SetTargetBuilding(tilePlayerId, tileBuildingId)
			iter := c.IterSelectedUnits(func(_ uint8, unitId uint8, unit *datamod.UnitsRow) bool {
				var (
					unitState   = rts.UnitState(unit.GetState())
					unitProtoId = unit.GetUnitType()
					unitProto   = c.Game().GetUnitPrototype(unitProtoId)
				)
				return unitProto.GetIsWorker() || unitState != rts.UnitState_Dead
			})
			buildingArea := c.Game().GetBuildingArea(tilePlayerId, tileBuildingId)
			match := c.Game().GetNearest(iter, buildingArea)
			workerId := match.ObjectId
			if workerId != rts.NilUnitId {
				c.Headless().AssignUnit(workerId, command)
			}
		}
		return
	}

	if tile.GetLandPlayerId() == c.PlayerId() && rts.ObjectType(tile.GetLandObjectType()) == rts.ObjectType_Unit {
		c.SelectUnits(tile.GetLandObjectId())
		return
	}

	if tile.GetHoverPlayerId() == c.PlayerId() && tile.GetHoverUnitId() != rts.NilObjectId {
		c.SelectUnits(tile.GetHoverUnitId())
		return
	}

	if tile.GetAirPlayerId() == c.PlayerId() && tile.GetAirUnitId() != rts.NilObjectId {
		c.SelectUnits(tile.GetAirUnitId())
		return
	}

	if tilePlayerId != rts.NilPlayerId && tilePlayerId != c.PlayerId() && tileObjectType == rts.ObjectType_Building {
		// Enemy building
		command := rts.NewFighterCommandData(rts.FighterCommandType_AttackBuilding)
		command.SetTargetBuilding(tilePlayerId, tileBuildingId)
		c.assignSelectedFighters(command, c.getCommandPath())
		c.clearCommandPath()
		return
	}

	c.ClearSelection()
}

func (c *Client) assignSelectedWorkers(command rts.UnitCommandData, idle bool) {
	unitIds := make([]uint8, 0)
	c.forEachSelectedUnit(func(unitId uint8, unit *datamod.UnitsRow) {
		var (
			unitProtoId = unit.GetUnitType()
			unitProto   = c.Game().GetUnitPrototype(unitProtoId)
			unitState   = rts.UnitState(unit.GetState())
			unitCommand = rts.WorkerCommandData(unit.GetCommand())
		)
		if !unitProto.GetIsWorker() || unitState == rts.UnitState_Dead || unitCommand.Type().IsIdle() != idle {
			return
		}
		unitIds = append(unitIds, unitId)
	})
	c.Headless().AssignUnits(unitIds, command)
}

func (c *Client) assignSelectedFighters(command rts.UnitCommandData, path []image.Point) {
	unitIds := make([]uint8, 0)
	c.forEachSelectedUnit(func(unitId uint8, unit *datamod.UnitsRow) {
		var (
			protoId   = unit.GetUnitType()
			proto     = c.Game().GetUnitPrototype(protoId)
			unitState = rts.UnitState(unit.GetState())
		)
		if proto.GetIsWorker() || unitState == rts.UnitState_Dead {
			return
		}
		unitIds = append(unitIds, unitId)
	})
	c.Headless().AssignUnitsWithPath(unitIds, command, path)
}

func (c *Client) moveCamera() {
	if c.coreRenderer.settings.FixedCamera {
		return
	}

	newCameraPosition := c.coreRenderer.cameraPosition
	newZoomLevel := c.coreRenderer.zoomLevel
	moveIncrement := 48 * InternalTileSize * assets.TileSize / c.coreRenderer.tileDisplaySize / ebiten.TPS()

	if c.keyMap.IsPressed(KeyFunction_Up) {
		newCameraPosition.Y -= moveIncrement
	}
	if c.keyMap.IsPressed(KeyFunction_Down) {
		newCameraPosition.Y += moveIncrement
	}
	if c.keyMap.IsPressed(KeyFunction_Left) {
		newCameraPosition.X -= moveIncrement
	}
	if c.keyMap.IsPressed(KeyFunction_Right) {
		newCameraPosition.X += moveIncrement
	}
	if c.keyMap.IsPressed(KeyFunction_CenterCamera) {
		newCameraPosition = c.Game().BoardSize().Mul(InternalTileSize).Div(2)
	} else if c.keyMap.IsPressed(KeyFunction_Base) {
		newCameraPosition = c.Game().GetMainBuildingPosition(c.PlayerId()).Add(image.Point{1, 1}).Mul(InternalTileSize)
	}
	if c.keyMap.IsJustPressed(KeyFunction_ZoomOut) {
		newZoomLevel -= 1
	}
	if c.keyMap.IsJustPressed(KeyFunction_ZoomIn) {
		newZoomLevel += 1
	}
	if _, dy := ebiten.Wheel(); dy != 0 {
		newZoomLevel += gen_utils.Sign(int(dy * 10))
	}
	if newZoomLevel != c.coreRenderer.zoomLevel || newCameraPosition != c.coreRenderer.cameraPosition {
		c.coreRenderer.setCamera(newCameraPosition, newZoomLevel)
	}
}

func (c *Client) Update() error {
	// Quitting
	// if c.keyMap.IsPressedWithShift(KeyFunction_Quit) {
	// 	return ErrQuit
	// }

	// Camera
	c.moveCamera()

	// Hud
	c.hud.Update(c)

	// Inputs
	c.handleInput()

	// Update
	if c.active {
		return c.coreRenderer.Update()
	}
	return nil
}

// Draw the game on screen.
func (c *Client) Draw(screen *ebiten.Image) {
	c.coreRenderer.Draw(screen)
	c.hud.Draw(c, screen)
}

// Return the layout for ebiten.
func (c *Client) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return c.coreRenderer.Layout(outsideWidth, outsideHeight)
}

type GenericClient struct {
	*Client
	uim *UIManager
}

func NewGenericClient(headlessClient IHeadlessClient, config ClientConfig, active bool) *GenericClient {
	var (
		nUnitPrototypes           = headlessClient.Game().GetMeta().GetUnitPrototypeCount()
		nBuildingPrototypes       = headlessClient.Game().GetMeta().GetBuildingPrototypeCount()
		menuUnitPrototypeIds      = make([]uint8, 0)
		menuBuildingPrototypeIds  = make([]uint8, 0)
		unitPrototypeMetadata     = make(map[uint8]PrototypeMetadata)
		buildingPrototypeMetadata = make(map[uint8]PrototypeMetadata)
	)
	for i := uint8(1); i <= nUnitPrototypes; i++ {
		menuUnitPrototypeIds = append(menuUnitPrototypeIds, i)
	}
	for i := uint8(1); i <= nBuildingPrototypes; i++ {
		proto := headlessClient.Game().GetBuildingPrototype(i)
		if proto.GetIsEnvironment() {
			continue
		}
		menuBuildingPrototypeIds = append(menuBuildingPrototypeIds, i)
	}
	for i := uint8(1); i <= nUnitPrototypes; i++ {
		unitPrototypeMetadata[i] = PrototypeMetadata{Name: fmt.Sprintf("Unit Type #%d", i)}
	}
	for i := uint8(1); i <= nBuildingPrototypes; i++ {
		buildingPrototypeMetadata[i] = PrototypeMetadata{Name: fmt.Sprintf("Building Type #%d", i)}
	}
	hudSet := NewHudSet()
	hudSet.AddComponents(
		NewTargetLines(),
		NewPathLines(),
		NewSelectionBox(),
		NewSelectionHighlight(),
		NewRangeHighlights(),
		NewTileDebugInfo(),
		NewBuildingGhost(),
	)
	var (
		coreRenderer = NewCoreRenderer(headlessClient, config, assets.DefaultSpriteGetter)
		cli          = NewClient(coreRenderer, hudSet, active)
		uim          = NewUI(
			cli,
			menuUnitPrototypeIds,
			unitPrototypeMetadata,
			menuBuildingPrototypeIds,
			buildingPrototypeMetadata,
			assets.DefaultSpriteGetter,
		)
	)
	return &GenericClient{
		Client: cli,
		uim:    uim,
	}
}

func (c *GenericClient) Update() error {
	if err := c.Client.Update(); err != nil {
		return err
	}
	c.uim.Update()

	// Keyboard
	for ii, protoId := range c.uim.menuBuildingPrototypeIds {
		if ebiten.IsKeyPressed(ebiten.Key1 + ebiten.Key(ii)) {
			c.SelectBuildingType(protoId)
			break
		}
	}

	// UI Click
	uiButtonClick := c.uim.PopButtonPress()
	if uiButtonClick != nil {
		switch uiButtonClick.ButtonType {
		case UI_ButtonType_BuildingIcon:
			c.SelectBuildingType(uint8(uiButtonClick.ButtonId))
		case UI_ButtonType_UnitIcon:
			c.CoreRenderer().CreateUnit(uint8(uiButtonClick.ButtonId))
		}
	}

	return nil
}

func (c *GenericClient) Draw(screen *ebiten.Image) {
	c.Client.Draw(screen)
	c.uim.Draw(screen)
}
