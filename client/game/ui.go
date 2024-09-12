package game

import (
	"fmt"

	"github.com/concrete-eth/ark-rts/client/assets"
	"github.com/concrete-eth/ark-rts/client/core"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

const (
	UI_Label_WinScreen = iota + core.UI_Id_Count
	UI_Container_EndScreen
	UI_Id_Count
)

type UI struct {
	*core.UIManager
}

func NewUI(cli *core.Client, spriteGetter assets.SpriteGetter) *UI {
	ui := &UI{
		UIManager: core.NewUI(cli, UnitPrototypeIds, spriteGetter),
	}

	endScreenContainer := newEndScreenContainer(ui)
	rootContainer := ui.UI().Container
	rootContainer.AddChild(endScreenContainer)

	return ui
}

func (m *UI) IsShowingEndScreen() bool {
	container := m.GetContainer(UI_Container_EndScreen)
	return container.GetWidget().Visibility == widget.Visibility_Show
}

func (m *UI) DismissEndScreen() {
	container := m.GetContainer(UI_Container_EndScreen)
	container.GetWidget().Visibility = widget.Visibility_Hide
}

func (m *UI) ShowLoseScreen() {
	container := m.GetContainer(UI_Container_EndScreen)
	label := m.GetLabel(UI_Label_WinScreen)
	label.Label = "You Lose!"
	container.GetWidget().Visibility = widget.Visibility_Show
}

func (m *UI) ShowEndScreen(winnerId uint8, outOfTime bool) {
	container := m.GetContainer(UI_Container_EndScreen)
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

func newEndScreenContainer(ui *UI) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(assets.UIFogColor)),
		widget.ContainerOpts.Layout((widget.NewAnchorLayout())),
	)

	boxContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceSimple(assets.UIBox_Big, assets.UICornerSize_Big, 2)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    24,
				Left:   36,
				Right:  36,
				Bottom: 24,
			}),
			widget.RowLayoutOpts.Spacing(4*core.StandardSpacing),
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
		widget.TextOpts.Text("Click to Dismiss", assets.GetFontFace(assets.Font_PressStart, 8), assets.TextDarkColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
	)

	boxContainer.AddChild(winLabel)
	boxContainer.AddChild(dismissLabel)
	container.AddChild(boxContainer)

	container.GetWidget().Visibility = widget.Visibility_Hide
	ui.AddContainer(UI_Container_EndScreen, container)
	ui.AddLabel(UI_Label_WinScreen, winLabel)

	return container
}
