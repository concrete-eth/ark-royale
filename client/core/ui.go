package core

import (
	"fmt"
	go_image "image"

	gen_utils "github.com/concrete-eth/archetype/utils"
	"github.com/concrete-eth/ark-royale/client/assets"
	client_utils "github.com/concrete-eth/ark-royale/client/utils"
	"github.com/concrete-eth/ark-royale/rts"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Point = go_image.Point

const (
	StandardSpacing = 2
	IconSize        = 64
)

// IDs to reference UI elements in the UIManager.
const (
	UI_ButtonType_UnitIcon = iota
	UI_ProgressBar_Resource
	UI_Id_Count
)

// Holds a button click event and the clicked button's type and id.
type buttonPress struct {
	ButtonType int
	ButtonId   int
	Args       *widget.ButtonPressedEventArgs
}

// Holds references to widgets and containers and exposes methods to update them.
type UIManager struct {
	client *Client      // Client the UI is attached to
	eui    *ebitenui.UI // Ebiten UI instance

	progressBars map[int]*widget.ProgressBar    // Progress bar widgets
	labels       map[int]*widget.Text           // Label widgets
	buttons      map[int]map[int]*widget.Button // Button widgets
	containers   map[int]*widget.Container      // Container widgets
	buttonPress  *buttonPress                   // Last button press event

	menuUnitPrototypeIds []uint8

	spriteGetter assets.SpriteGetter
}

// Creates a new UI and UIManager instance attached to the given client.
func NewUI(
	client *Client,
	menuUnitPrototypeIds []uint8,
	spriteGetter assets.SpriteGetter,
) *UIManager {
	uim := &UIManager{
		client:               client,
		progressBars:         make(map[int]*widget.ProgressBar),
		labels:               make(map[int]*widget.Text),
		buttons:              make(map[int]map[int]*widget.Button),
		containers:           make(map[int]*widget.Container),
		menuUnitPrototypeIds: menuUnitPrototypeIds,
		spriteGetter:         spriteGetter,
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)
	menuContainer := newMenuContainer(uim)
	rootContainer.AddChild(menuContainer)

	uim.eui = &ebitenui.UI{Container: rootContainer}

	return uim
}

func (m *UIManager) Client() *Client {
	return m.client
}

func (m *UIManager) UI() *ebitenui.UI {
	return m.eui
}

func (m *UIManager) SpriteGetter() assets.SpriteGetter {
	return m.spriteGetter
}

func (m *UIManager) Regenerate() *UIManager {
	return NewUI(m.client, m.menuUnitPrototypeIds, m.spriteGetter)
}

// Adds a progress bar widget to the UI.
func (m *UIManager) RegisterProgressBar(id int, bar *widget.ProgressBar) {
	m.progressBars[id] = bar
}

// Returns a progress bar widget by id.
func (m *UIManager) GetProgressBar(id int) *widget.ProgressBar {
	return m.progressBars[id]
}

// Adds a label widget to the UI.
func (m *UIManager) RegisterButton(buttonType, buttonId int, button *widget.Button) {
	if _, ok := m.buttons[buttonType]; !ok {
		m.buttons[buttonType] = make(map[int]*widget.Button)
	}
	m.buttons[buttonType][buttonId] = button
}

// Returns a button widget by type and id.
func (m *UIManager) GetButton(gid, id int) *widget.Button {
	if _, ok := m.buttons[gid]; !ok {
		return nil
	}
	return m.buttons[gid][id]
}

// Adds a label widget to the UI.
func (m *UIManager) RegisterLabel(id int, label *widget.Text) {
	m.labels[id] = label
}

// Returns a label widget by id.
func (m *UIManager) GetLabel(id int) *widget.Text {
	return m.labels[id]
}

// Adds a container widget to the UI.
func (m *UIManager) RegisterContainer(id int, container *widget.Container) {
	m.containers[id] = container
}

// Returns a container widget by id.
func (m *UIManager) GetContainer(id int) *widget.Container {
	return m.containers[id]
}

// Returns the last button press event and clears it.
func (m *UIManager) PopButtonPress() *buttonPress {
	press := m.buttonPress
	m.buttonPress = nil
	return press
}

// Updates the UI state.
func (m *UIManager) Update() {
	m.setResourceIndicators(m.client.PlayerId())
	m.updateCreationMenu()
	m.eui.Update()
}

func (m *UIManager) Draw(screen *ebiten.Image) {
	m.eui.Draw(screen)
}

// Updates the unit menu by enabling/disabling unit icons based on whether the player has an armory.
func (m *UIManager) updateCreationMenu() {
	var (
		player           = m.client.Game().GetPlayer(m.client.PlayerId())
		computeDemand    = player.GetComputeDemand()
		computeSupply    = player.GetComputeSupply()
		computeSurplus   = gen_utils.SafeSubUint8(computeSupply, computeDemand)
		resourceCapacity = player.GetMaxResource()
	)
	for object := range m.client.coreRenderer.anticipatedObjects {
		if object.PlayerId != m.client.PlayerId() {
			continue
		}
		if object.Type == rts.ObjectType_Unit {
			computeSurplus = gen_utils.SafeSubUint8(computeSurplus, 1) // Compute cost is assumed to be 1
		}
	}
	for _, protoId := range m.menuUnitPrototypeIds {
		proto := m.client.Game().GetUnitPrototype(protoId)
		button := m.GetButton(UI_ButtonType_UnitIcon, int(protoId))
		// insufficientCompute := proto.GetComputeCost() > computeSurplus
		insufficientResourceCapacity := proto.GetResourceCost() > resourceCapacity
		button.GetWidget().Disabled = insufficientResourceCapacity
	}
}

// Sets the resource indicators.
func (m *UIManager) setResourceIndicators(playerId uint8) {
	player := m.client.Game().GetPlayer(playerId)

	curResource := player.GetCurResource()
	maxResource := player.GetMaxResource()

	bar := m.GetProgressBar(UI_ProgressBar_Resource)
	bar.Max = int(maxResource)
	bar.SetCurrent(int(curResource))
	label := m.GetLabel(UI_ProgressBar_Resource)
	label.Label = fmt.Sprintf("%d/%d", int(curResource), int(maxResource))
}

// Creates a new button press handler.
func (m *UIManager) newButtonPressHandler(buttonType int, buttonId int) widget.ButtonPressedHandlerFunc {
	return func(args *widget.ButtonPressedEventArgs) {
		m.buttonPress = &buttonPress{
			ButtonType: buttonType,
			ButtonId:   buttonId,
			Args:       args,
		}
	}
}

func newMenuContainer(uim *UIManager) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(2 * StandardSpacing)),
		)),
	)
	container.AddChild(newCenterMenu(uim))
	return container
}

func newCenterMenu(uim *UIManager) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(StandardSpacing),
		)),
	)
	resourceInfo := newResourceDisplay(uim)
	unitMenu := newUnitMenu(uim, uim.menuUnitPrototypeIds)
	container.AddChild(resourceInfo)
	container.AddChild(unitMenu)
	return container
}

func newUnitMenu(uim *UIManager, unitPrototypeIds []uint8) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceSimple(assets.UIBox_Small, assets.UICornerSize_Small, assets.UICornerSize_Small)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(assets.UICornerSize_Small)),
		)),
	)
	if len(unitPrototypeIds) == 0 {
		return container
	}
	for _, protoId := range unitPrototypeIds {
		container.AddChild(newUnitIcon(uim, protoId, IconSize))
	}
	return container
}

func newUnitIcon(uim *UIManager, protoId uint8, size int) *widget.Button {
	var dir assets.Direction
	if uim.client.PlayerId() == 1 {
		dir = assets.Direction_Right
	} else {
		dir = assets.Direction_Left
	}
	proto := uim.client.Game().GetUnitPrototype(protoId)
	cost := int(proto.GetResourceCost())
	sprite := uim.spriteGetter.GetUnitSprite(uim.client.PlayerId(), protoId, dir)
	buttonImage := newIconButtonImage(sprite, size, assets.UICornerSize_Small*2, cost)
	button := newIconButton(
		buttonImage,
		uim.newButtonPressHandler(UI_ButtonType_UnitIcon, int(protoId)),
		size,
	)
	uim.RegisterButton(UI_ButtonType_UnitIcon, int(protoId), button)
	return button
}

func newIconButton(buttonImage *widget.ButtonImage, handler widget.ButtonPressedHandlerFunc, size int) *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.PressedHandler(handler),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(size, size),
		),
	)
	return button
}

func newIconButtonImage(sprite *ebiten.Image, size, margin, cost int) *widget.ButtonImage {
	img := ebiten.NewImage(size, size)
	imgBounds := img.Bounds()

	bgNineSlice := image.NewNineSliceSimple(assets.UIPanel_Convex, assets.UICornerSize_Small, assets.UICornerSize_Small)
	bgNineSlice.Draw(img, size, size, nil)

	cornerBounds := go_image.Point{margin, margin}
	spriteRect := go_image.Rectangle{Min: cornerBounds, Max: imgBounds.Max.Sub(cornerBounds)}
	op := client_utils.NewDrawOptions(spriteRect, sprite.Bounds())
	colorm.DrawImage(img, sprite, colorm.ColorM{}, op)

	text.Draw(img, fmt.Sprint(cost), assets.BitmapFont1, margin, margin+8, assets.TextLightColor)

	borderWidthHeight := cornerBounds.X
	centerWidthHeight := size - borderWidthHeight*2

	defaultNineSlice := image.NewNineSliceSimple(img, borderWidthHeight, centerWidthHeight)
	lightNineSlice := image.NewNineSliceSimple(assets.ChangeImageHSV(img, 0, 1, 1.1), borderWidthHeight, centerWidthHeight)
	mediumDarkNineSlice := image.NewNineSliceSimple(assets.ChangeImageHSV(img, 0, 1, 0.9), borderWidthHeight, centerWidthHeight)
	darkNineSlice := image.NewNineSliceSimple(assets.ChangeImageHSV(img, 0, 0.25, 0.75), borderWidthHeight, centerWidthHeight)

	return &widget.ButtonImage{
		Idle:     defaultNineSlice,
		Hover:    lightNineSlice,
		Pressed:  mediumDarkNineSlice,
		Disabled: darkNineSlice,
	}
}

func newResourceDisplay(uim *UIManager) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(StandardSpacing),
		)),
	)
	container.AddChild(newProgressBar(uim, UI_ProgressBar_Resource, "minerals", assets.UIProgressBar_Mineral))
	return container
}

func newProgressBar(uim *UIManager, id int, name string, sprite *ebiten.Image) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
		),
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)
	resourceBar := widget.NewProgressBar(
		widget.ProgressBarOpts.Images(
			&widget.ProgressBarImage{
				Idle: image.NewNineSliceSimple(assets.UIProgressBar_Track, assets.UICornerSize_Small, assets.UICornerSize_Small),
			},
			&widget.ProgressBarImage{
				Idle: image.NewNineSliceSimple(sprite, assets.UICornerSize_Small, assets.UICornerSize_Small),
			},
		),
		widget.ProgressBarOpts.Values(0, 1000, 0),
	)
	nameText := widget.NewText(
		widget.TextOpts.Text(name, assets.BitmapFont1, assets.LightGray),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
		widget.TextOpts.Insets(widget.Insets{Left: 2 + assets.UICornerSize_Small, Right: 2}),
	)
	barText := widget.NewText(
		widget.TextOpts.Text("", assets.BitmapFont1, assets.LightGray),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.Insets(widget.Insets{Left: 2, Right: 2}),
	)

	container.AddChild(resourceBar)
	container.AddChild(nameText)
	container.AddChild(barText)

	uim.RegisterProgressBar(id, resourceBar)
	uim.RegisterLabel(id, barText)

	return container
}
