package core

import (
	"image"
	"image/color"

	"github.com/concrete-eth/ark-rts/client/assets"
	client_utils "github.com/concrete-eth/ark-rts/client/utils"
	"github.com/concrete-eth/ark-rts/gogen/datamod"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

func setAnticipationDrawSettings(dc *gg.Context, anticipating bool) {
	if anticipating {
		dc.SetDash(8, 8)
	} else {
		dc.SetDash()
	}
}

func drawLine(dc *gg.Context, start, end image.Point) {
	dc.SetLineWidth(3)
	dc.DrawLine(float64(start.X), float64(start.Y), float64(end.X), float64(end.Y))
	dc.Stroke()
}

func drawCircle(dc *gg.Context, position image.Point) {
	dc.DrawCircle(float64(position.X), float64(position.Y), float64(3))
	dc.Fill()
}

type HudComponent interface{}

// Holds a set of hud components.
type HudSet struct {
	Components []HudComponent
}

// Creates a new HudSet.
func NewHudSet() *HudSet {
	return &HudSet{
		Components: []HudComponent{},
	}
}

// Adds components to the HudSet.
func (hs *HudSet) AddComponents(cc ...HudComponent) {
	hs.Components = append(hs.Components, cc...)
}

// Calls Update on all components that implement Updatable.
func (hs *HudSet) Update(c *Client) {
	for _, hc := range hs.Components {
		if hc, ok := hc.(UpdatableWithClient); ok {
			hc.Update(c)
		}
	}
}

// Calls Draw on all components that implement Drawable.
func (hs *HudSet) Draw(c *Client, screen *ebiten.Image) {
	// Pass the sub-image that correspond to the board display
	subScreen := screen.SubImage(c.coreRenderer.boardDisplayRect).(*ebiten.Image)
	for _, hc := range hs.Components {
		if hc, ok := hc.(DrawableWithClient); ok {
			hc.Draw(c, subScreen)
		}
	}
}

// Returns the center of the unit in the draw context.
func unitCenterInContext(c *Client, playerId, unitId uint8) image.Point {
	scaledCameraOrigin := c.coreRenderer.cameraPosition.Mul(c.coreRenderer.tileDisplaySize).Div(InternalTileSize).Sub(c.coreRenderer.boardDisplayRect.Size().Div(2))
	unitInternalPosition := c.coreRenderer.getPosition(rts.Object{
		Type:     rts.ObjectType_Unit,
		PlayerId: playerId,
		ObjectId: unitId,
	})
	return unitInternalPosition.
		Mul(c.coreRenderer.tileDisplaySize).Div(InternalTileSize).
		Sub(scaledCameraOrigin).
		Add(image.Point{1, 1}.Mul(c.coreRenderer.tileDisplaySize / 2))
}

func tileCenterInContext(c *Client, tilePosition image.Point) image.Point {
	return c.coreRenderer.TileCoordToDisplayCoord(tilePosition).Add(image.Point{1, 1}.Mul(c.coreRenderer.tileDisplaySize / 2))
}

// Returns the center of the building in the draw context.
func buildingCenterInContext(c *Client, playerId, buildingId uint8) image.Point {
	var (
		building     = c.Game().GetBuilding(playerId, buildingId)
		protoId      = building.GetBuildingType()
		proto        = c.Game().GetBuildingPrototype(protoId)
		size         = rts.GetDimensionsAsPoint(proto)
		tilePosition = rts.GetPositionAsPoint(building)
	)
	return c.coreRenderer.TileCoordToDisplayCoord(tilePosition).Add(size.Mul(c.coreRenderer.tileDisplaySize).Div(2))
}

// Renders a building ghost where the cursor is when a building is selected.
type UnitGhost struct {
	ghostImage    *ebiten.Image
	ghostUnitType uint8
	colorM        colorm.ColorM
}

var _ UpdatableWithClient = (*UnitGhost)(nil)
var _ DrawableWithClient = (*UnitGhost)(nil)

// Creates a new BuildingGhost component.
func NewUnitGhost() *UnitGhost {
	return &UnitGhost{}
}

// Creates the ghost image by drawing the building sprite into a shadow box.
func (bg *UnitGhost) drawGhost(c *Client, playerId uint8, unitTypeId uint8) *ebiten.Image {
	var direction assets.Direction
	if playerId == 1 {
		direction = assets.Direction_Right
	} else {
		direction = assets.Direction_Left
	}
	sprite := c.coreRenderer.spriteGetter.GetUnitSprite(playerId, unitTypeId, direction)

	colorM := assets.NewUnitColorMatrix()
	colorM.Concat(assets.GhostColorMatrix)
	img := ebiten.NewImage(sprite.Bounds().Dx(), sprite.Bounds().Dy())
	colorm.DrawImage(img, sprite, colorM, nil)

	return img
}

func (bg *UnitGhost) SetColorMatrix(colorM colorm.ColorM) {
	bg.colorM = colorM
}

// Redraws the ghost image it if the selected building type has changed.
func (bg *UnitGhost) Update(c *Client) {
	if !c.IsSelectingUnitType() {
		bg.ghostUnitType = 0
		return
	}
	if bg.ghostUnitType != c.selected.UnitType {
		bg.ghostImage = bg.drawGhost(c, c.PlayerId(), c.selected.UnitType)
		bg.ghostUnitType = c.selected.UnitType
	}
}

// Draws the ghost image onto the screen at the cursor position.
func (bg *UnitGhost) Draw(c *Client, screen *ebiten.Image) {
	if !c.IsSelectingUnitType() || bg.ghostImage == nil {
		// If no building is selected, or the ghost image has not been created yet, do nothing.
		return
	}
	var (
		cursorPos = client_utils.CursorPosition()
		tilePos   = c.coreRenderer.ScreenCoordToTileCoord(cursorPos)
		layerPos  = tilePos.Mul(assets.TileSize).Sub(assets.UnitSpriteOrigin)
		screenPos = c.coreRenderer.TileCoordToScreenCoord(image.Point{}).Add(
			layerPos.Mul(c.coreRenderer.tileDisplaySize).Div(InternalTileSize),
		)
		screenRect = image.Rectangle{
			Min: screenPos,
			Max: screenPos.Add(bg.ghostImage.Bounds().Size().Mul(c.coreRenderer.tileDisplaySize).Div(assets.TileSize)),
		}
		op = client_utils.NewDrawOptions(screenRect, bg.ghostImage.Bounds())
	)
	// if !tilePosition.In(c.Game().BoardRect()) {
	// 	return
	// }
	colorm.DrawImage(screen, bg.ghostImage, bg.colorM, op)
}

// Renders the range highlights around selected units.
type RangeHighlights struct {
	highlightTile *ebiten.Image
}

var _ DrawableWithClient = (*RangeHighlights)(nil)

// Creates a new rangeHighlights component.
func NewRangeHighlights() *RangeHighlights {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return &RangeHighlights{
		highlightTile: img,
	}
}

// Redraws the range highlights.
func (rh *RangeHighlights) Draw(c *Client, _ *ebiten.Image) {
	layer := c.coreRenderer.worldLayers.Layer(LayerName_HudTerrain)
	spriteObj := layer.Sprite("rangeHighlights")

	cursorPos := client_utils.CursorPosition()
	tilePos := c.coreRenderer.ScreenCoordToTileCoord(cursorPos)
	if !tilePos.In(c.Game().BoardRect()) {
		spriteObj.SetImage(nil)
		return
	}

	units := make([]rts.Object, 0)
	tile := c.Game().GetBoardTile(uint16(tilePos.X), uint16(tilePos.Y))
	if tile.GetLandObjectType() == rts.ObjectType_Unit.Uint8() {
		units = append(units, rts.Object{
			Type:     rts.ObjectType_Unit,
			PlayerId: tile.GetLandPlayerId(),
			ObjectId: tile.GetLandObjectId(),
		})
	}
	if tile.GetAirPlayerId() != rts.NilPlayerId {
		units = append(units, rts.Object{
			Type:     rts.ObjectType_Unit,
			PlayerId: tile.GetAirPlayerId(),
			ObjectId: tile.GetAirUnitId(),
		})
	}

	colorM := colorm.ColorM{}
	colorM.Scale(1, 1, 1, 0.5)
	// colorM.Translate(1, 1, 1, 0)
	spriteObj.SetColorMatrix(colorM)

	contextSize := c.coreRenderer.boardDisplayRect.Size()
	dc := gg.NewContext(contextSize.X, contextSize.Y)
	dc.SetLineWidth(2)
	dc.SetDash(8, 8)

	for _, unitObj := range units {
		var (
			unit         = c.Game().GetUnit(unitObj.PlayerId, unitObj.ObjectId)
			unitState    = rts.UnitState(unit.GetState())
			unitPosition = rts.GetPositionAsPoint(unit)
			protoId      = unit.GetUnitType()
			proto        = c.Game().GetUnitPrototype(protoId)
			unitRange    = proto.GetAttackRange()
		)
		if proto.GetIsWorker() {
			return
		}
		if unitState.IsDeadOrInactive() {
			return
		}
		chebyshevRadius := int(unitRange)
		area := image.Rectangle{
			Min: c.coreRenderer.TileCoordToDisplayCoord(unitPosition.Sub(image.Point{chebyshevRadius, chebyshevRadius})),
			Max: c.coreRenderer.TileCoordToDisplayCoord(unitPosition.Add(image.Point{chebyshevRadius + 1, chebyshevRadius + 1})),
		}

		var fillColor, borderColor color.RGBA
		if unitObj.PlayerId == 1 {
			fillColor = color.RGBA{0x00, 0x00, 0xFF, 0x20}
			borderColor = color.RGBA{0x00, 0x00, 0xFF, 0xFF}
		} else {
			fillColor = color.RGBA{0xFF, 0x00, 0x00, 0x20}
			borderColor = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
		}

		dc.DrawRectangle(float64(area.Min.X), float64(area.Min.Y), float64(area.Dx()), float64(area.Dy()))
		dc.SetColor(color.RGBA{0xFF, 0xFF, 0xFF, 0x40})
		dc.SetColor(fillColor)
		dc.FillPreserve()
		dc.SetColor(borderColor)
		dc.Stroke()
	}

	img := ebiten.NewImageFromImage(dc.Image())
	spriteObj.SetImage(img).FitToImage()
}

// Renders the tile debug info highlighting non-empty tiles in a layer.
type TileDebugInfo struct {
	mode int
}

var _ UpdatableWithClient = (*TileDebugInfo)(nil)
var _ DrawableWithClient = (*TileDebugInfo)(nil)

// Creates a new tileDebugInfo component.
func NewTileDebugInfo() *TileDebugInfo {
	return &TileDebugInfo{}
}

// Changes the selected layer when the D key is pressed.
func (td *TileDebugInfo) Update(c *Client) {
	if c.keyMap.IsJustPressed(KeyFunction_ToggleDebugInfo) {
		td.mode = (td.mode + 1) % 4
	}
}

// Draws the tile debug info.
func (td *TileDebugInfo) Draw(c *Client, screen *ebiten.Image) {
	boardSize := c.Game().BoardSize()
	for xx := 0; xx < boardSize.X; xx++ {
		for yy := 0; yy < boardSize.Y; yy++ {
			// For every tile, check if it is empty in the selected layer.

			tilePosition := image.Point{xx, yy}
			tile := c.Game().GetBoardTile(uint16(tilePosition.X), uint16(tilePosition.Y))

			var highlight bool
			switch td.mode {
			case 1:
				highlight = !rts.IsTileEmpty(tile, rts.LayerId_Land)
			case 2:
				highlight = !rts.IsTileEmpty(tile, rts.LayerId_Hover)
			case 3:
				highlight = !rts.IsTileEmpty(tile, rts.LayerId_Air)
			default:
				return
			}
			if !highlight {
				continue
			}

			img := ebiten.NewImage(1, 1)
			img.Fill(color.White)

			colorM := colorm.ColorM{}
			colorM.Scale(1, 1, 1, 0.5)

			screenRect := image.Rectangle{
				Min: c.coreRenderer.TileCoordToScreenCoord(tilePosition),
				Max: c.coreRenderer.TileCoordToScreenCoord(tilePosition.Add(image.Point{1, 1})),
			}

			op := client_utils.NewDrawOptions(screenRect, img.Bounds())
			colorm.DrawImage(screen, img, colorM, op)
		}
	}
}

// Renders lines from units to their targets.
type TargetLines struct{}

var _ UpdatableWithClient = (*TargetLines)(nil)

// Creates a new targetLines component.
func NewTargetLines() *TargetLines {
	return &TargetLines{}
}

func (tl *TargetLines) drawTargetLine(c *Client, dc *gg.Context, playerId uint8, unitId uint8, unit *datamod.UnitsRow, command *AnticipatedCommand, anticipated bool) {
	var (
		bounds    = dc.Image().Bounds()
		protoId   = unit.GetUnitType()
		proto     = c.Game().GetUnitPrototype(protoId)
		unitState = rts.UnitState(unit.GetState())
	)
	if protoId == 0 {
		return
	}
	if unitState.IsDeadOrInactive() {
		return
	}

	var unitLayerPosition image.Point
	if c.coreRenderer.hasPosition(rts.Object{
		Type:     rts.ObjectType_Unit,
		PlayerId: playerId,
		ObjectId: unitId,
	}) {
		unitLayerPosition = unitCenterInContext(c, playerId, unitId)
	} else {
		unitLayerPosition = tileCenterInContext(c, rts.GetPositionAsPoint(unit))
	}

	var targetPosition image.Point
	var targetPlayerId uint8
	var targetBuildingId uint8
	var targetUnitId uint8

	if proto.GetIsWorker() {
		var unitCommand rts.WorkerCommandData
		if command == nil {
			unitCommand = rts.WorkerCommandData(unit.GetCommand())
		} else {
			unitCommand = rts.WorkerCommandData(command.Command)
		}
		if unitCommand.Type().IsIdle() {
			targetPosition = c.Game().GetWorkerPortPosition(playerId)
		} else {
			targetPlayerId, targetBuildingId = unitCommand.TargetBuilding()
		}
	} else {
		var unitCommand rts.FighterCommandData
		if command == nil {
			unitCommand = rts.FighterCommandData(unit.GetCommand())
		} else {
			unitCommand = rts.FighterCommandData(command.Command)
		}
		if unitCommand.Type().IsTargetingPosition() {
			targetPosition = unitCommand.TargetPosition()
		} else if unitCommand.Type().IsTargetingBuilding() {
			targetPlayerId, targetBuildingId = unitCommand.TargetBuilding()
			targetBuilding := c.Game().GetBuilding(targetPlayerId, targetBuildingId)
			targetBuildingState := rts.BuildingState(targetBuilding.GetState())
			if targetBuildingState == rts.BuildingState_Destroyed {
				return
			}
		} else if unitCommand.Type().IsTargetingUnit() {
			targetPlayerId, targetUnitId = unitCommand.TargetUnit()
			targetUnit := c.Game().GetUnit(targetPlayerId, targetUnitId)
			targetUnitState := rts.UnitState(targetUnit.GetState())
			if targetUnitState == rts.UnitState_Dead {
				return
			}
		} else {
			return
		}
		// match := c.Game().GetUnitToFireAt(c.Game().GetUnitObject(playerId, unitId))
		// if !match.IsNil() && unitState.IsActive() {
		// 	targetUnitLayerPosition := unitCenterInContext(c, match.PlayerId, match.ObjectId)
		// 	if unitLayerPosition.In(bounds) || targetUnitLayerPosition.In(bounds) {
		// 		setAnticipationDrawSettings(dc, false)
		// 		drawLine(dc, unitLayerPosition, targetUnitLayerPosition)
		// 	}
		// }
	}

	var targetLayerPosition image.Point
	if targetBuildingId != rts.NilBuildingId {
		targetLayerPosition = buildingCenterInContext(c, targetPlayerId, targetBuildingId)
	} else if targetUnitId != rts.NilUnitId {
		targetLayerPosition = unitCenterInContext(c, targetPlayerId, targetUnitId)
	} else {
		targetLayerPosition = tileCenterInContext(c, targetPosition)
	}

	setAnticipationDrawSettings(dc, anticipated)

	var path *rts.CommandPath
	if command == nil {
		path = rts.NewCommandPath(unit.GetCommandExtra(), unit.GetCommandMeta())
	} else {
		path = rts.NewCommandPath(command.Extra, command.Meta)
	}

	layerPath := make([]image.Point, 0)
	layerPath = append(layerPath, unitLayerPosition)
	if path.HasPath() && path.Pointer() < path.PathLen() {
		for _, point := range path.Path()[path.Pointer():] {
			layerPath = append(layerPath, tileCenterInContext(c, point))
		}
	} else {
		if rts.Distance(targetLayerPosition, unitLayerPosition) < c.coreRenderer.tileDisplaySize/2 {
			return
		}
		if !unitLayerPosition.In(bounds) && !targetLayerPosition.In(bounds) {
			return
		}
	}
	layerPath = append(layerPath, targetLayerPosition)

	for ii, point := range layerPath {
		drawCircle(dc, point)
		if ii > 0 {
			prevPoint := layerPath[ii-1]
			drawLine(dc, prevPoint, point)
		}
	}
}

// Redraws the target lines.
func (tl *TargetLines) Update(c *Client) {
	layer := c.coreRenderer.worldLayers.Layer(LayerName_HudLines)
	spriteObj := layer.Sprite("targetLines")

	// if c.coreRenderer.tileDisplaySize < 16 {
	// 	// If the target lines are not visible, or the display tile size is too small (too zoomed out), hide the component.
	// 	spriteObj.SetImage(nil)
	// 	return
	// }

	cursorPos := client_utils.CursorPosition()
	tilePos := c.coreRenderer.ScreenCoordToTileCoord(cursorPos)
	if !tilePos.In(c.Game().BoardRect()) {
		spriteObj.SetImage(nil)
		return
	}

	units := make([]rts.Object, 0)
	tile := c.Game().GetBoardTile(uint16(tilePos.X), uint16(tilePos.Y))
	if tile.GetLandObjectType() == rts.ObjectType_Unit.Uint8() {
		units = append(units, rts.Object{
			Type:     rts.ObjectType_Unit,
			PlayerId: tile.GetLandPlayerId(),
			ObjectId: tile.GetLandObjectId(),
		})
	}
	if tile.GetAirPlayerId() != rts.NilPlayerId {
		units = append(units, rts.Object{
			Type:     rts.ObjectType_Unit,
			PlayerId: tile.GetAirPlayerId(),
			ObjectId: tile.GetAirUnitId(),
		})
	}

	colorM := colorm.ColorM{}
	colorM.Scale(1, 1, 1, 0.75)
	spriteObj.SetColorMatrix(colorM)

	contextSize := c.coreRenderer.boardDisplayRect.Size()
	dc := gg.NewContext(contextSize.X, contextSize.Y)
	dc.SetColor(color.White)

	for _, unitObj := range units {
		var (
			playerId           = unitObj.PlayerId
			unitId             = unitObj.ObjectId
			unit               = c.Game().GetUnit(playerId, unitId)
			anticipatedCommand = c.coreRenderer.anticipatedCommands[unitObj]
		)
		if rts.BuildingState(c.Game().GetMainBuilding(playerId).GetState()) == rts.BuildingState_Destroyed {
			continue
		}
		if anticipatedCommand == (AnticipatedCommand{}) {
			tl.drawTargetLine(c, dc, playerId, unitId, unit, nil, false)
		} else {
			tl.drawTargetLine(c, dc, playerId, unitId, unit, &anticipatedCommand, true)
		}
	}

	img := ebiten.NewImageFromImage(dc.Image())
	spriteObj.SetImage(img).FitToImage()
}
