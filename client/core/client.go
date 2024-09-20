package core

import (
	"image"

	gen_utils "github.com/concrete-eth/archetype/utils"
	"github.com/concrete-eth/ark-royale/client/assets"
	"github.com/concrete-eth/ark-royale/rts"
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
	UnitType uint8
}

// Clear clears the selection.
func (s *Selection) Clear() {
	s.UnitType = 0
}

// Main game client object run by ebiten.RunGame.
type Client struct {
	coreRenderer      *CoreRenderer
	hud               *HudSet   // HUD component set
	keyMap            KeyMap    // Key map
	selected          Selection // Current selection
	onSelectionChange func()    // On selection change callback
	active            bool      // Actively update state
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

	if c.onSelectionChange != nil {
		c.onSelectionChange()
	}
}

// Clears the current selection, selects a given building type, and executes the
// selection side effects: show building shadows and buildable area.
func (c *Client) SelectUnitType(unitType uint8) {
	c.ClearSelection()
	c.selected.UnitType = unitType

	if c.onSelectionChange != nil {
		c.onSelectionChange()
	}
}

// Returns the selected building type.
func (c *Client) SelectedUnitType() uint8 {
	return c.selected.UnitType
}

// Returns true if a building type is selected.
func (c *Client) IsSelectingUnitType() bool {
	return c.selected.UnitType != 0
}

func (c *Client) CreateUnit(unitType uint8, position image.Point) {
	queue := c.coreRenderer.internalEventQueue
	c.Headless().CreateUnit(unitType, position)
	c.coreRenderer.internalEventQueue = queue
}

func (c *Client) handleInput() {
	if c.keyMap.IsJustPressed(KeyFunction_Deselect) {
		c.ClearSelection()
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		c.ClearSelection()
	}

	cursorScreenPosition := image.Pt(ebiten.CursorPosition())

	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}
	if !cursorScreenPosition.In(c.coreRenderer.boardDisplayRect) {
		return
	}
	terrainDisplayRect := image.Rectangle{
		Min: c.coreRenderer.TileCoordToDisplayCoord(image.Pt(0, 0)),
		Max: c.coreRenderer.TileCoordToDisplayCoord(c.Game().BoardSize()),
	}
	if !cursorScreenPosition.In(terrainDisplayRect) {
		return
	}

	var (
		tilePosition = c.coreRenderer.ScreenCoordToTileCoord(cursorScreenPosition)
		tile         = c.Game().GetBoardTile(uint16(tilePosition.X), uint16(tilePosition.Y))
	)
	if !rts.IsTileEmptyAllLayers(tile) {
		return
	}

	c.CreateUnit(c.SelectedUnitType(), tilePosition)
	c.ClearSelection()
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
