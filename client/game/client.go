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
	unitGhost       *core.UnitGhost
}

func NewClient(headlessClient core.IHeadlessClient, config core.ClientConfig, active bool) *Client {
	hudSet := core.NewHudSet()
	hudSet.AddComponents(
		// core.NewSelectionBox(),
		// core.NewSelectionHighlight(),
		core.NewTargetLines(),
		core.NewRangeHighlights(),
		core.NewTileDebugInfo(),
	)
	var (
		whl          = headlessClient
		coreRenderer = core.NewCoreRenderer(whl, config, SpriteGetter)
		cli          = core.NewClient(coreRenderer, hudSet, active)
		uim          = NewUI(cli, SpriteGetter)
	)
	c := &Client{
		Client:    cli,
		uim:       uim,
		unitGhost: core.NewUnitGhost(),
	}
	hudSet.AddComponents(c.unitGhost)

	cli.SetOnSelectionChange(func() {
		c.toggleShowSpawnArea(c.IsSelectingUnitType())
	})
	cli.CoreRenderer().SetOnCameraMove(func() {
		c.setSpawnAreaSpriteRect()
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
func (c *Client) toggleShowSpawnArea(show bool) {
	hudTerrainLayer := c.CoreRenderer().Layers().Layer(core.LayerName_HudTerrain)
	spriteObj := hudTerrainLayer.Sprite("spawnTerrain")
	if show && spriteObj.Image() == nil {
		c.setSpawnAreaSprite()
	}
	spriteObj.SetVisible(show)
}

func (c *Client) setSpawnAreaSprite() {
	c.setSpawnAreaSpriteImage()
	c.setSpawnAreaSpriteRect()
}

// Creates an overlay shadowing out all non-buildable tiles for the client player.
func (c *Client) setSpawnAreaSpriteImage() {
	spawnAreaImage := ebiten.NewImage(1, 1)
	spawnAreaImage.Fill(assets.LightBlueShadowColor)
	hudTerrainLayer := c.CoreRenderer().Layers().Layer(core.LayerName_HudTerrain)
	hudTerrainLayer.Sprite("spawnTerrain").SetImage(spawnAreaImage)
}

func (c *Client) setSpawnAreaSpriteRect() {
	spawnArea := c.Game().GetSpawnArea(c.PlayerId())
	hudTerrainLayer := c.CoreRenderer().Layers().Layer(core.LayerName_HudTerrain)
	hudTerrainLayer.Sprite("spawnTerrain").SetRect(
		image.Rectangle{
			Min: c.CoreRenderer().TileCoordToScreenCoord(spawnArea.Min),
			Max: c.CoreRenderer().TileCoordToScreenCoord(spawnArea.Max),
		},
	)
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
		c.setSpawnAreaSprite()
		c.uim.DismissWinScreen()
		log.Debug("Switched player", "from", prevPlayerId, "to", c.PlayerId())
	}
}

func (c *Client) setUnitGhostColor() {
	if !c.IsSelectingUnitType() {
		return
	}
	var (
		cursorScreenPosition = client_utils.CursorPosition()
		tilePosition         = c.CoreRenderer().ScreenCoordToTileCoord(cursorScreenPosition).Div(2).Mul(2)
		size                 = image.Point{1, 1}
		spawnArea            = image.Rectangle{Min: tilePosition, Max: tilePosition.Add(size)}
	)
	if spawnArea.In(c.Game().GetSpawnArea(c.PlayerId())) {
		c.unitGhost.SetColorMatrix(colorm.ColorM{})
	} else {
		c.unitGhost.SetColorMatrix(assets.NonBuildableColorMatrix)
	}
}

func (c *Client) Update() error {
	c.debugChangePlayer()

	if err := c.Client.Update(); err != nil {
		return err
	}
	c.uim.Update()

	c.setUnitGhostColor()

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
		// case core.UI_ButtonType_BuildingIcon:
		// 	c.SelectUnitType(uint8(uiButtonClick.ButtonId))
		case core.UI_ButtonType_UnitIcon:
			c.SelectUnitType(uint8(uiButtonClick.ButtonId))
		}
	}

	return nil
}

func (c *Client) Draw(screen *ebiten.Image) {
	c.Client.Draw(screen)
	c.uim.Draw(screen)
}
