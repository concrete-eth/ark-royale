package game

import (
	"fmt"

	"github.com/concrete-eth/ark-rts/client/assets"
	"github.com/concrete-eth/ark-rts/client/core"
	"github.com/concrete-eth/ark-rts/gogen/datamod"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

const (
	UI_Label_WinScreen = iota
	UI_Container_WinScreen
)

var (
	UnitPrototypesMetadata = map[uint8]core.PrototypeMetadata{
		UnitPrototypeId_Worker:  {Name: "Worker"},
		UnitPrototypeId_AntiAir: {Name: "Anti-Air"},
		UnitPrototypeId_Air:     {Name: "Air"},
		UnitPrototypeId_Tank:    {Name: "Tank"},
	}
	BuildingPrototypeNames = map[uint8]core.PrototypeMetadata{
		BuildingPrototypeId_Main:    {Name: "Command Center"},
		BuildingPrototypeId_Storage: {Name: "Storage"},
		BuildingPrototypeId_Lab:     {Name: "Lab"},
		BuildingPrototypeId_Armory:  {Name: "Armory"},
		BuildingPrototypeId_Mine:    {Name: "Mineral Mine"},
	}
)

func buildingDescriptionStrings(protoId uint8, proto *datamod.BuildingPrototypesRow) string {
	description := ""
	switch protoId {
	case BuildingPrototypeId_Main:
		description = "The main building of your base."
	case BuildingPrototypeId_Storage:
		description = "Increases your resource capacity."
	case BuildingPrototypeId_Lab:
		description = "Increases your compute capacity (cpu)."
	case BuildingPrototypeId_Armory:
		description = "Allows you to build advanced units."
	}
	return description
}

type UI struct {
	*core.UIManager
}

func NewUI(cli *core.Client, spriteGetter assets.SpriteGetter) *UI {
	buildingPrototypesMetadata := make(map[uint8]core.PrototypeMetadata, len(BuildingPrototypeNames))
	for id, meta := range BuildingPrototypeNames {
		proto := cli.Game().GetBuildingPrototype(id)
		meta.Description = buildingDescriptionStrings(id, proto)
		buildingPrototypesMetadata[id] = meta
	}
	ui := &UI{
		UIManager: core.NewUI(
			cli,
			UnitPrototypeIds,
			UnitPrototypesMetadata,
			BuildableBuildingPrototypeIds,
			buildingPrototypesMetadata,
			spriteGetter,
		),
	}

	winScreenContainer := newOverDisplayContainer(ui)
	overDisplayContainer := ui.GetContainer(core.UI_Container_OverDisplay)
	overDisplayContainer.AddChild(winScreenContainer)

	return ui
}

func (m *UI) IsShowingWinScreen() bool {
	container := m.GetContainer(UI_Container_WinScreen)
	return container.GetWidget().Visibility == widget.Visibility_Show
}

func (m *UI) DismissWinScreen() {
	container := m.GetContainer(UI_Container_WinScreen)
	container.GetWidget().Visibility = widget.Visibility_Hide
}

func (m *UI) ShowLoseScreen() {
	container := m.GetContainer(UI_Container_WinScreen)
	label := m.GetLabel(UI_Label_WinScreen)
	label.Label = "You Lose!"
	container.GetWidget().Visibility = widget.Visibility_Show
}

func (m *UI) ShowEndScreen(winnerId uint8, outOfTime bool) {
	container := m.GetContainer(UI_Container_WinScreen)
	label := m.GetLabel(UI_Label_WinScreen)

	if outOfTime {
		label.Label = "Out of Time!"
	} else if winnerId == rts.NilPlayerId {
		label.Label = "Mutual Annihilation!"
	} else {
		label.Label = fmt.Sprintf("Player %d Wins!", winnerId)
	}
	container.GetWidget().Visibility = widget.Visibility_Show
}

func (ui *UI) Regenerate() *UI {
	return NewUI(ui.Client(), ui.SpriteGetter())
}

func (ui *UI) Update() {
	ui.UIManager.Update()
	var (
		player    = ui.Client().Game().GetPlayer(ui.Client().PlayerId())
		hasArmory = player.GetCurArmories() > 0
	)
	for _, protoId := range UnitPrototypeIds {
		button := ui.GetButton(core.UI_ButtonType_UnitIcon, int(protoId))
		missingArmory := !hasArmory && NeedsArmory(protoId)
		disable := missingArmory
		button.GetWidget().Disabled = button.GetWidget().Disabled || disable
	}
}

func newOverDisplayContainer(ui *UI) *widget.Container {
	rect := ui.Client().CoreRenderer().BoardDisplayRect()
	// subContainer covers the area over the game display
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	subContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal: true,
				StretchVertical:   true,
			}),
			widget.WidgetOpts.MinSize(rect.Dx(), rect.Dy()),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	winScreenContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal: true,
				StretchVertical:   true,
			}),
		),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(assets.DarkBlueShadowColor)),
		widget.ContainerOpts.Layout((widget.NewAnchorLayout())),
	)

	winLabelContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceSimple(assets.UIPanel_Dark, assets.UICornerSize, assets.UICornerSize)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    24,
				Left:   36,
				Right:  36,
				Bottom: 24,
			}),
			widget.RowLayoutOpts.Spacing(4*core.BorderWidth),
		)),
	)

	winLabel := widget.NewText(
		widget.TextOpts.Text("Default", assets.GetFontFace(assets.Font_PressStart, 16), assets.TextLightColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
	)
	dismissLabel := widget.NewText(
		widget.TextOpts.Text("Click to Dismiss", assets.GetFontFace(assets.Font_PressStart, 8), assets.BoxTextDarkColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
	)

	winLabelContainer.AddChild(winLabel)
	winLabelContainer.AddChild(dismissLabel)
	winScreenContainer.AddChild(winLabelContainer)
	subContainer.AddChild(winScreenContainer)
	container.AddChild(subContainer)

	winScreenContainer.GetWidget().Visibility = widget.Visibility_Hide
	ui.AddContainer(UI_Container_WinScreen, winScreenContainer)
	ui.AddLabel(UI_Label_WinScreen, winLabel)

	return container
}
