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
type BuildingGhost struct {
	ghostImage        *ebiten.Image
	ghostBuildingType uint8
	colorM            colorm.ColorM
}

var _ UpdatableWithClient = (*BuildingGhost)(nil)
var _ DrawableWithClient = (*BuildingGhost)(nil)

// Creates a new BuildingGhost component.
func NewBuildingGhost() *BuildingGhost {
	return &BuildingGhost{}
}

// Creates the ghost image by drawing the building sprite into a shadow box.
func (bg *BuildingGhost) drawGhost(c *Client, playerId uint8, proto *datamod.BuildingPrototypesRow) *ebiten.Image {
	imgSize := rts.GetDimensionsAsPoint(proto).Mul(assets.TileSize)
	img := ebiten.NewImage(imgSize.X, imgSize.Y)
	img.Fill(assets.LightBlueShadowColor)

	spriteId := c.coreRenderer.spriteGetter.GetBuildingSpriteId(playerId, 0, proto)
	sprite := c.coreRenderer.spriteGetter.GetBuildingSprite(playerId, spriteId, rts.BuildingState_Built)
	colorM := assets.NewBuildingColorMatrix()
	colorM.Concat(assets.GhostColorMatrix)
	op := client_utils.NewDrawOptions(image.Rectangle{Max: imgSize}, sprite.Bounds())
	colorm.DrawImage(img, sprite, colorM, op)

	return img
}

func (bg *BuildingGhost) SetColorMatrix(colorM colorm.ColorM) {
	bg.colorM = colorM
}

// Redraws the ghost image it if the selected building type has changed.
func (bg *BuildingGhost) Update(c *Client) {
	if !c.IsSelectingBuildingType() {
		bg.ghostBuildingType = 0
		return
	}
	if bg.ghostBuildingType != c.selected.BuildingType {
		proto := c.Game().GetBuildingPrototype(c.selected.BuildingType)
		bg.ghostImage = bg.drawGhost(c, c.PlayerId(), proto)
		bg.ghostBuildingType = c.selected.BuildingType
	}
}

// Draws the ghost image onto the screen at the cursor position.
func (bg *BuildingGhost) Draw(c *Client, screen *ebiten.Image) {
	if !c.IsSelectingBuildingType() || bg.ghostImage == nil {
		// If no building is selected, or the ghost image has not been created yet, do nothing.
		return
	}
	var (
		screenPosition = client_utils.CursorPosition()
		tilePosition   = c.coreRenderer.ScreenCoordToTileCoord(screenPosition).Div(2).Mul(2)
		proto          = c.Game().GetBuildingPrototype(c.selected.BuildingType)
		size           = rts.GetDimensionsAsPoint(proto)
		screenRect     = image.Rectangle{
			Min: c.coreRenderer.TileCoordToScreenCoord(tilePosition),
			Max: c.coreRenderer.TileCoordToScreenCoord(tilePosition.Add(size)),
		}
		op = client_utils.NewDrawOptions(screenRect, bg.ghostImage.Bounds())
	)
	if !tilePosition.In(c.Game().BoardRect()) {
		return
	}
	colorm.DrawImage(screen, bg.ghostImage, bg.colorM, op)
}

// Renders lines from units to their targets.
type TargetLines struct {
	visible bool
}

var _ UpdatableWithClient = (*TargetLines)(nil)

// Creates a new targetLines component.
func NewTargetLines() *TargetLines {
	return &TargetLines{
		visible: true,
	}
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
	if unitState == rts.UnitState_Dead {
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
		} else {
			return
		}
		match := c.Game().GetUnitToFireAt(c.Game().GetUnitObject(playerId, unitId))
		if !match.IsNil() && unitState.IsActive() {
			targetUnitLayerPosition := unitCenterInContext(c, match.PlayerId, match.ObjectId)
			if unitLayerPosition.In(bounds) || targetUnitLayerPosition.In(bounds) {
				setAnticipationDrawSettings(dc, false)
				drawLine(dc, unitLayerPosition, targetUnitLayerPosition)
			}
		}
	}

	var targetLayerPosition image.Point
	if targetBuildingId == rts.NilBuildingId {
		targetLayerPosition = tileCenterInContext(c, targetPosition)
	} else {
		targetLayerPosition = buildingCenterInContext(c, targetPlayerId, targetBuildingId)
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
	if c.keyMap.IsJustReleased(KeyFunction_ToggleTargetLines) {
		// If the L key is pressed, toggle the visibility of the target lines.
		tl.visible = !tl.visible
	}

	layer := c.coreRenderer.worldLayers.Layer(LayerName_HudLines)
	spriteObj := layer.Sprite("targetLines")

	if !tl.visible || c.coreRenderer.tileDisplaySize < 16 {
		// If the target lines are not visible, or the display tile size is too small (too zoomed out), hide the component.
		spriteObj.SetImage(nil)
		return
	}

	colorM := colorm.ColorM{}
	colorM.Scale(1, 1, 1, 0.75)
	spriteObj.SetColorMatrix(colorM)

	contextSize := c.coreRenderer.boardDisplayRect.Size()
	dc := gg.NewContext(contextSize.X, contextSize.Y)
	dc.SetColor(color.White)

	nPlayers := c.Game().GetMeta().GetPlayerCount()
	for playerId := uint8(1); playerId < nPlayers+1; playerId++ {
		if rts.BuildingState(c.Game().GetMainBuilding(playerId).GetState()) == rts.BuildingState_Destroyed {
			continue
		}
		c.Game().ForEachUnit(playerId, func(unitId uint8, unit *datamod.UnitsRow) {
			tl.drawTargetLine(c, dc, playerId, unitId, unit, nil, false)
		})
	}
	for object, command := range c.coreRenderer.anticipatedCommands {
		var (
			playerId = object.PlayerId
			unitId   = object.ObjectId
			unit     = c.Game().GetUnit(playerId, unitId)
		)
		if rts.BuildingState(c.Game().GetMainBuilding(playerId).GetState()) == rts.BuildingState_Destroyed {
			continue
		}
		tl.drawTargetLine(c, dc, playerId, unitId, unit, &command, true)
	}
	img := ebiten.NewImageFromImage(dc.Image())
	spriteObj.SetImage(img).FitToImage()
}

// Renders the selection box when the left mouse button is held down.
type selectionBox struct{}

var _ DrawableWithClient = (*selectionBox)(nil)

// Creates a new selectionBox component.
func NewSelectionBox() *selectionBox {
	return &selectionBox{}
}

// Draws the selection box when the left mouse button is held down.
func (bg *selectionBox) Draw(c *Client, screen *ebiten.Image) {
	if !c.lastClickedScreenPosition.In(c.coreRenderer.boardDisplayRect) {
		return
	}
	// if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || !c.selecting {
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return
	}
	cursorPosition := client_utils.CursorPosition()
	selectionRect := image.Rectangle{Min: c.lastClickedScreenPosition, Max: cursorPosition}.Canon()
	if selectionRect.Dx() == 0 || selectionRect.Dy() == 0 {
		return
	}

	dc := gg.NewContext(selectionRect.Dx(), selectionRect.Dy())
	dc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xc0})
	dc.SetLineWidth(2)
	dc.DrawRectangle(1, 1, float64(selectionRect.Dx())-2, float64(selectionRect.Dy())-2)
	dc.Stroke()

	img := ebiten.NewImageFromImage(dc.Image())
	op := client_utils.NewDrawOptions(selectionRect, img.Bounds())
	colorm.DrawImage(screen, img, colorm.ColorM{}, op)
}

// Renders the selection highlight box around selected units.
type selectionHighlight struct{}

var _ DrawableWithClient = (*selectionHighlight)(nil)

// Creates a new selectionHighlight component.
func NewSelectionHighlight() *selectionHighlight {
	return &selectionHighlight{}
}

// Draws the selection highlight box around selected units.
func (sh *selectionHighlight) Draw(c *Client, screen *ebiten.Image) {
	if c.coreRenderer.tileDisplaySize < 16 {
		// If the display tile size is too small (too zoomed out), do nothing.
		return
	}
	if len(c.selected.Units) == 0 {
		// If no units are selected, do nothing.
		return
	}
	c.forEachSelectedUnit(func(unitId uint8, unit *datamod.UnitsRow) {
		unitState := rts.UnitState(unit.GetState())
		if unitState.IsDeadOrInactive() {
			// If the unit is hidden, do nothing.
			return
		}
		pos := c.coreRenderer.GetUnitScreenPosition(c.PlayerId(), unitId)
		img := assets.SelectionSprite
		pos = pos.Sub(image.Point{1, 1}.Mul(c.coreRenderer.tileDisplaySize / 2))
		size := image.Point{1, 1}.Mul(c.coreRenderer.tileDisplaySize * 2)
		op := client_utils.NewDrawOptions(image.Rectangle{Min: pos, Max: pos.Add(size)}, img.Bounds())
		colorm.DrawImage(screen, img, colorm.ColorM{}, op)
	})
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

	if !c.IsSelectingUnits() {
		spriteObj.SetImage(nil)
		return
	}

	colorM := colorm.ColorM{}
	colorM.Scale(1, 1, 1, 0.75)
	colorM.Translate(1, 1, 1, 0)
	spriteObj.SetColorMatrix(colorM)

	contextSize := c.coreRenderer.boardDisplayRect.Size()
	dc := gg.NewContext(contextSize.X, contextSize.Y)
	dc.SetLineWidth(2)
	dc.SetDash(8, 8)

	c.forEachSelectedUnit(func(unitId uint8, unit *datamod.UnitsRow) {
		var (
			unitState    = rts.UnitState(unit.GetState())
			unitPosition = rts.GetPositionAsPoint(unit)
			protoId      = unit.GetUnitType()
			proto        = c.Game().GetUnitPrototype(protoId)
			unitRange    = proto.GetAttackRange()
		)
		if proto.GetIsWorker() {
			return
		}
		if !unitState.IsAlive() {
			return
		}
		chebyshevRadius := int(unitRange)
		area := image.Rectangle{
			Min: c.coreRenderer.TileCoordToDisplayCoord(unitPosition.Sub(image.Point{chebyshevRadius, chebyshevRadius})),
			Max: c.coreRenderer.TileCoordToDisplayCoord(unitPosition.Add(image.Point{chebyshevRadius + 1, chebyshevRadius + 1})),
		}
		dc.DrawRectangle(float64(area.Min.X), float64(area.Min.Y), float64(area.Dx()), float64(area.Dy()))
		dc.SetColor(color.RGBA{0x00, 0x00, 0x00, 0x40})
		dc.FillPreserve()
		dc.SetColor(color.Black)
		dc.Stroke()
	})

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

type pathLines struct{}

var _ UpdatableWithClient = (*pathLines)(nil)

func NewPathLines() *pathLines {
	return &pathLines{}
}

func (pl *pathLines) Update(c *Client) {
	layer := c.coreRenderer.worldLayers.Layer(LayerName_HudLines)
	spriteObj := layer.Sprite("pathLines")

	if c.coreRenderer.tileDisplaySize < 16 {
		// If the target lines are not visible, or the display tile size is too small (too zoomed out), hide the component.
		spriteObj.SetImage(nil)
		return
	}

	colorM := colorm.ColorM{}
	colorM.Scale(1, 1, 1, 0.75)
	spriteObj.SetColorMatrix(colorM)

	contextSize := c.coreRenderer.boardDisplayRect.Size()
	dc := gg.NewContext(contextSize.X, contextSize.Y)
	dc.SetColor(color.White)
	setAnticipationDrawSettings(dc, true)

	path := c.commandPath
	if len(path) > 0 {
		cursorTilePosition := c.coreRenderer.ScreenCoordToTileCoord(client_utils.CursorPosition())
		if !cursorTilePosition.Eq(path[len(path)-1]) {
			path = append(path, cursorTilePosition)
		}
	}

	for ii, point := range path {
		contextPosition := tileCenterInContext(c, point)
		drawCircle(dc, contextPosition)
		if ii > 0 {
			prevContextPosition := tileCenterInContext(c, path[ii-1])
			drawLine(dc, prevContextPosition, contextPosition)
		}
	}

	if len(path) > 0 {
		c.forEachSelectedUnit(func(unitId uint8, unit *datamod.UnitsRow) {
			unitState := rts.UnitState(unit.GetState())
			if unitState.IsDeadOrInactive() {
				return
			}
			unitLayerPosition := unitCenterInContext(c, c.PlayerId(), unitId)
			drawCircle(dc, unitLayerPosition)
			if unitLayerPosition.In(dc.Image().Bounds()) {
				drawLine(dc, unitLayerPosition, tileCenterInContext(c, path[0]))
			}
		})
	}

	img := ebiten.NewImageFromImage(dc.Image())
	spriteObj.SetImage(img).FitToImage()
}
