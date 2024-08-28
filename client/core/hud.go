package core

import (
	"image"
	"image/color"

	"github.com/concrete-eth/ark-rts/client/assets"
	client_utils "github.com/concrete-eth/ark-rts/client/utils"
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
		direction = assets.Direction_Down
	} else {
		direction = assets.Direction_Up
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
