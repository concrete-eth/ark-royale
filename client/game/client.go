package game

import (
	"image"

	client_utils "github.com/concrete-eth/ark-rts/client/utils"

	"github.com/concrete-eth/ark-rts/client/assets"
	"github.com/concrete-eth/ark-rts/client/core"
	"github.com/concrete-eth/ark-rts/gogen/datamod"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/ethereum/go-ethereum/log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Client struct {
	*core.Client
	uim             *UI
	shownLoseScreen bool // True if the lose screen is shown
	shownEndScreen  bool // True if the end screen is shown
	buildingGhost   *core.BuildingGhost
}

func NewClient(headlessClient core.IHeadlessClient, config core.ClientConfig, hinter core.Hinter, active bool) *Client {
	hudSet := core.NewHudSet()
	hudSet.AddComponents(
		core.NewTargetLines(),
		core.NewPathLines(),
		core.NewSelectionBox(),
		core.NewSelectionHighlight(),
		core.NewRangeHighlights(),
		core.NewTileDebugInfo(),
	)
	var (
		whl          = &HeadlessClient{headlessClient}
		coreRenderer = core.NewCoreRenderer(whl, config, hinter, assets.DefaultSpriteGetter)
		cli          = core.NewClient(coreRenderer, hudSet, active)
		uim          = NewUI(cli, assets.DefaultSpriteGetter)
	)
	c := &Client{
		Client:        cli,
		uim:           uim,
		buildingGhost: core.NewBuildingGhost(),
	}
	hudSet.AddComponents(c.buildingGhost)

	cli.SetOnSelectionChange(func() {
		c.toggleShowBuildableArea(c.IsSelectingBuildingType())
	})
	cli.CoreRenderer().SetOnCameraMove(func() {
		c.setNonBuildableAreaSpritePosition()
	})
	cli.CoreRenderer().SetOnNewBatch(func() {
		c.checkGameOver()
	})

	return c
}

func (c *Client) UI() *UI {
	return c.uim
}

// Toggles the visibility of the buildable area.
func (c *Client) toggleShowBuildableArea(show bool) {
	hudTerrainLayer := c.CoreRenderer().Layers().Layer(core.LayerName_HudTerrain)
	spriteObj := hudTerrainLayer.Sprite("buildableTerrain")
	if show && spriteObj.Image() == nil {
		c.setNonBuildableAreaSprite()
	}
	spriteObj.SetVisible(show)
}

func (c *Client) setNonBuildableAreaSprite() {
	c.setNonBuildableAreaSpriteImage()
	c.setNonBuildableAreaSpritePosition()
}

// Creates an overlay shadowing out all non-buildable tiles for the client player.
func (c *Client) setNonBuildableAreaSpriteImage() {
	nonBuildableRect := c.Game().BoardRect()
	buildableRect := GetPlayerBuildableArea(c.Game(), c.PlayerId())

	buildableImage := ebiten.NewImage(1, 1)
	nonBuildableImage := ebiten.NewImage(nonBuildableRect.Dx(), nonBuildableRect.Dy())
	nonBuildableImage.Fill(assets.DarkShadowColor)

	op := client_utils.NewDrawOptions(buildableRect.Sub(nonBuildableRect.Min), buildableImage.Bounds())
	op.Blend = ebiten.BlendSourceIn // See: https://ebitengine.org/en/examples/masking.html
	colorm.DrawImage(nonBuildableImage, buildableImage, colorm.ColorM{}, op)

	hudTerrainLayer := c.CoreRenderer().Layers().Layer(core.LayerName_HudTerrain)
	hudTerrainLayer.Sprite("buildableTerrain").SetImage(nonBuildableImage)
}

func (c *Client) setNonBuildableAreaSpritePosition() {
	boardSizePixels := c.Game().BoardSize().Mul(c.CoreRenderer().TileDisplaySize())
	hudTerrainLayer := c.CoreRenderer().Layers().Layer(core.LayerName_HudTerrain)
	hudTerrainLayer.Sprite("buildableTerrain").
		SetPosition(c.CoreRenderer().TileCoordToDisplayCoord(image.Point{0, 0})).
		SetSize(boardSizePixels)
}

func (c *Client) checkGameOver() {
	if c.shownEndScreen {
		return
	}
	activePlayers := []uint8{}
	c.Game().ForEachPlayer(func(playerId uint8, player *datamod.PlayersRow) {
		// Check for active main building
		mainBuilding := c.Game().GetBuilding(playerId, 1)
		if mainBuilding.GetState() == uint8(rts.BuildingState_Destroyed) {
			return
		}
		activePlayers = append(activePlayers, playerId)
	})
	if len(activePlayers) == 0 {
		c.uim.ShowEndScreen(rts.NilPlayerId, false)
		c.shownEndScreen = true
	} else if len(activePlayers) == 1 {
		c.uim.ShowEndScreen(activePlayers[0], false)
		c.shownEndScreen = true
	} else if !c.shownLoseScreen {
		lost := true
		for _, playerId := range activePlayers {
			if playerId == c.PlayerId() {
				lost = false
				break
			}
		}
		if lost {
			c.uim.ShowLoseScreen()
			c.shownLoseScreen = true
		}
	}

	creationBlockNumber := c.Game().GetMeta().GetCreationBlockNumber()
	if uint32(c.Game().BlockNumber())-creationBlockNumber > MaxTicks {
		c.uim.ShowEndScreen(rts.NilPlayerId, true)
		c.shownEndScreen = true
	}
}

func (c *Client) debugChangePlayer() {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		prevPlayerId := c.PlayerId()
		c.CoreRenderer().SetPlayerId(c.PlayerId()%c.Game().GetMeta().GetPlayerCount() + 1)
		c.uim = c.uim.Regenerate()
		c.ClearSelection()
		c.CoreRenderer().Layers().Layer(core.LayerName_HudLines).Clear()
		c.setNonBuildableAreaSprite()
		c.uim.DismissWinScreen()
		log.Debug("Switched player", "from", prevPlayerId, "to", c.PlayerId())
	}
}

func (c *Client) setBuildingGhostColor() {
	if !c.IsSelectingBuildingType() {
		return
	}
	var (
		cursorScreenPosition = client_utils.CursorPosition()
		tilePosition         = c.CoreRenderer().ScreenCoordToTileCoord(cursorScreenPosition).Div(2).Mul(2)
		proto                = c.Game().GetBuildingPrototype(c.SelectedBuildingType())
		size                 = rts.GetDimensionsAsPoint(proto)
		buildArea            = image.Rectangle{Min: tilePosition, Max: tilePosition.Add(size)}
	)
	if IsInPlayerBuildableArea(c.Game(), c.PlayerId(), buildArea) {
		c.buildingGhost.SetColorMatrix(colorm.ColorM{})
	} else {
		c.buildingGhost.SetColorMatrix(assets.NonBuildableColorMatrix)
	}
}

func (c *Client) Update() error {
	c.debugChangePlayer()

	if err := c.Client.Update(); err != nil {
		return err
	}
	c.uim.Update()

	c.setBuildingGhostColor()

	// Keyboard
	for ii, protoId := range BuildableBuildingPrototypeIds {
		if ebiten.IsKeyPressed(ebiten.Key1 + ebiten.Key(ii)) {
			// This is enforced at the client level instead of the game level as it only serves
			// to protect users from placing buildings they don't have the resource capacity to pay for,
			// as placing them would permanently block their building payment queue.
			resourceCapacity := c.Game().GetPlayer(c.PlayerId()).GetMaxResource()
			proto := c.Game().GetBuildingPrototype(protoId)
			insufficientResourceCapacity := proto.GetResourceCost() > resourceCapacity
			if insufficientResourceCapacity {
				c.ClearSelection()
			} else {
				c.SelectBuildingType(protoId)
			}
			break
		}
	}

	// Win screen
	if c.uim.IsShowingWinScreen() {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			c.uim.DismissWinScreen()
		}
		return nil
	}

	// UI Click
	uiButtonClick := c.uim.PopButtonPress()
	if uiButtonClick != nil {
		switch uiButtonClick.ButtonType {
		case core.UI_ButtonType_BuildingIcon:
			c.SelectBuildingType(uint8(uiButtonClick.ButtonId))
		case core.UI_ButtonType_UnitIcon:
			c.CoreRenderer().CreateUnit(uint8(uiButtonClick.ButtonId))
		}
	}

	return nil
}

func (c *Client) Draw(screen *ebiten.Image) {
	c.Client.Draw(screen)
	c.uim.Draw(screen)
}
