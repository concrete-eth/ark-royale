package core

import (
	"fmt"
	go_image "image"

	gen_utils "github.com/concrete-eth/archetype/utils"
	"github.com/concrete-eth/ark-rts/client/assets"
	client_utils "github.com/concrete-eth/ark-rts/client/utils"
	"github.com/concrete-eth/ark-rts/gogen/datamod"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Point = go_image.Point

const (
	BorderWidth  = 2
	PanelPadding = 8
)

// IDs to reference UI elements in the UIManager.
const (
	UI_ButtonType_BuildingIcon = iota
	UI_ButtonType_UnitIcon
	UI_ProgressBar_Resource
	UI_ProgressBar_Compute
	UI_Label_InspectedName
	UI_Label_InspectedDetails
	UI_Label_WorkerCount
	UI_Label_Compute
	UI_Label_FPS
	UI_Label_Interpolation
	UI_Container_OverDisplay
)

// Holds a button click event and the clicked button's type and id.
type buttonPress struct {
	ButtonType int
	ButtonId   int
	Args       *widget.ButtonPressedEventArgs
}

type PrototypeMetadata struct {
	Name        string
	Description string
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

	menuUnitPrototypeIds     []uint8
	menuBuildingPrototypeIds []uint8

	unitPrototypeMetadata     map[uint8]PrototypeMetadata
	buildingPrototypeMetadata map[uint8]PrototypeMetadata

	spriteGetter assets.SpriteGetter
}

// Creates a new UI and UIManager instance attached to the given client.
func NewUI(
	client *Client,
	menuUnitPrototypeIds []uint8,
	unitPrototypeMetadata map[uint8]PrototypeMetadata,
	menuBuildingPrototypeIds []uint8,
	buildingPrototypeMetadata map[uint8]PrototypeMetadata,
	spriteGetter assets.SpriteGetter,
) *UIManager {
	uim := &UIManager{
		client:                    client,
		progressBars:              make(map[int]*widget.ProgressBar),
		labels:                    make(map[int]*widget.Text),
		buttons:                   make(map[int]map[int]*widget.Button),
		containers:                make(map[int]*widget.Container),
		menuUnitPrototypeIds:      menuUnitPrototypeIds,
		unitPrototypeMetadata:     unitPrototypeMetadata,
		menuBuildingPrototypeIds:  menuBuildingPrototypeIds,
		buildingPrototypeMetadata: buildingPrototypeMetadata,
		spriteGetter:              spriteGetter,
	}

	// Create the root container with a two-column layout
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
		)),
	)

	// Create the container that shows over the game board display
	overDisplayContainer := newOverDisplayContainer(uim)

	// Create a new container for side UI
	sideContainer := newSideContainer(uim)

	// Add the top left container and the right column container to the root container
	rootContainer.AddChild(overDisplayContainer)
	rootContainer.AddChild(sideContainer)

	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	uim.eui = eui

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
	return NewUI(m.client, m.menuUnitPrototypeIds, m.unitPrototypeMetadata, m.menuBuildingPrototypeIds, m.buildingPrototypeMetadata, m.spriteGetter)
}

// Adds a progress bar widget to the UI.
func (m *UIManager) AddProgressBar(id int, bar *widget.ProgressBar) {
	m.progressBars[id] = bar
}

// Returns a progress bar widget by id.
func (m *UIManager) GetProgressBar(id int) *widget.ProgressBar {
	return m.progressBars[id]
}

// Adds a label widget to the UI.
func (m *UIManager) AddButton(buttonType, buttonId int, button *widget.Button) {
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
func (m *UIManager) AddLabel(id int, label *widget.Text) {
	m.labels[id] = label
}

// Returns a label widget by id.
func (m *UIManager) GetLabel(id int) *widget.Text {
	return m.labels[id]
}

// Adds a container widget to the UI.
func (m *UIManager) AddContainer(id int, container *widget.Container) {
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
	// m.setResourceIndicators(m.client.PlayerId())
	// m.setComputeIndicators(m.client.PlayerId())
	// m.setWorkerIndicator(m.client.PlayerId())
	// m.setComputeIndicator(m.client.PlayerId())

	m.updateCreationMenu()
	m.updateFPSLabel()
	m.updateInterpolationLabel()
	// m.updateInspectionView()

	m.eui.Update()
}

func (m *UIManager) Draw(screen *ebiten.Image) {
	m.eui.Draw(screen)
}

func (m *UIManager) updateInterpolationLabel() {
	label := m.GetLabel(UI_Label_Interpolation)
	label.Label = boolToOnOff(m.client.coreRenderer.Interpolating())
}

// Updates the FPS label.
func (m *UIManager) updateFPSLabel() {
	fps := ebiten.ActualFPS()
	fpsLabel := m.GetLabel(UI_Label_FPS)
	fpsLabel.Label = fmt.Sprintf("%.1f", fps)
}

// Updates the inspection view.
func (m *UIManager) updateInspectionView() {
	inspectedNameLabel := m.GetLabel(UI_Label_InspectedName)
	inspectedNameLabel.Label = ""
	inspectedDetailsLabel := m.GetLabel(UI_Label_InspectedDetails)
	inspectedDetailsLabel.Label = ""

	for _, protoId := range m.menuBuildingPrototypeIds {
		button := m.GetButton(UI_ButtonType_BuildingIcon, int(protoId))
		if client_utils.CursorPosition().In(button.GetWidget().Rect) {
			// If the cursor is over a building icon, show the building's name and details
			proto := m.client.Game().GetBuildingPrototype(protoId)
			inspectedNameLabel.Label = m.buildingPrototypeMetadata[protoId].Name
			inspectedDetailsLabel.Label = m.buildingDetailsString(protoId, proto)
			return
		}
	}

	for _, protoId := range m.menuUnitPrototypeIds {
		button := m.GetButton(UI_ButtonType_UnitIcon, int(protoId))
		if client_utils.CursorPosition().In(button.GetWidget().Rect) {
			// If the cursor is over a unit icon, show the unit's name and details
			proto := m.client.Game().GetUnitPrototype(protoId)
			inspectedNameLabel.Label = m.unitPrototypeMetadata[protoId].Name
			inspectedDetailsLabel.Label = m.unitDetailsString(protoId, proto)
			return
		}
	}

	// if len(m.client.selected.Units) > 0 {
	// 	if len(m.client.selected.Units) > 1 {
	// 		// If there are multiple selected units, show the number
	// 		inspectedNameLabel.Label = fmt.Sprintf("%d units", len(m.client.selected.Units))
	// 	} else {
	// 		// If there is a single selected unit, show the unit's name
	// 		unitId := m.client.selected.Units[0]
	// 		unit := m.client.Game().GetUnit(m.client.PlayerId(), unitId)
	// 		protoId := unit.GetUnitType()
	// 		inspectedNameLabel.Label = m.unitPrototypeMetadata[protoId].Name + " Unit"
	// 	}
	// } else {
	// 	// If there are no selected units, show the name of the building under the cursor
	// 	cursorTilePosition := m.client.coreRenderer.ScreenCoordToTileCoord(client_utils.CursorPosition())
	// 	boardTile := m.client.Game().GetBoardTile(uint16(cursorTilePosition.X), uint16(cursorTilePosition.Y))
	// 	buildingId := boardTile.GetLandObjectId()
	// 	if buildingId != rts.NilBuildingId {
	// 		playerId := boardTile.GetLandPlayerId()
	// 		building := m.client.Game().GetBuilding(playerId, buildingId)
	// 		protoId := building.GetBuildingType()
	// 		if protoId != 0 {
	// 			inspectedNameLabel.Label = m.buildingPrototypeMetadata[protoId].Name
	// 			switch playerId {
	// 			case rts.NilPlayerId:
	// 				inspectedDetailsLabel.Label = "Neutral"
	// 			case m.client.PlayerId():
	// 				inspectedDetailsLabel.Label = "Friendly"
	// 			default:
	// 				inspectedDetailsLabel.Label = "Enemy"
	// 			}
	// 		}
	// 	}
	// }
}

func (m *UIManager) unitDetailsString(protoId uint8, proto *datamod.UnitPrototypesRow) string {
	return m.unitStatsString(protoId, proto)
}

func (m *UIManager) unitStatsString(protoId uint8, proto *datamod.UnitPrototypesRow) string {
	return fmt.Sprintf(""+
		"Type     : %s\n"+
		"Cost     : %d\n"+
		"Range    : %d\n"+
		"Strength : %d land; %d hover; %d air",
		cases.Title(language.English).String(rts.LayerId(proto.GetLayer()).String()),
		proto.GetResourceCost(), proto.GetAttackRange(),
		proto.GetLandStrength(), proto.GetHoverStrength(), proto.GetAirStrength(),
	)
}

func (m *UIManager) buildingDetailsString(protoId uint8, proto *datamod.BuildingPrototypesRow) string {
	description := m.buildingPrototypeMetadata[protoId].Description
	stats := m.buildingStatsString(protoId, proto)
	if description == "" {
		return stats
	}
	return fmt.Sprintf("%s\n%s", description, stats)
}

func (m *UIManager) buildingStatsString(protoId uint8, proto *datamod.BuildingPrototypesRow) string {
	return fmt.Sprintf(""+
		"Cost      : %d\n"+
		"Integrity : %d\n"+
		"Storage   : %d",
		proto.GetResourceCost(), proto.GetMaxIntegrity(), proto.GetResourceCapacity(),
	)
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
	if rts.BuildingState(m.client.Game().GetMainBuilding(m.client.PlayerId()).GetState()) == rts.BuildingState_Destroyed {
		for _, protoId := range m.menuBuildingPrototypeIds {
			button := m.GetButton(UI_ButtonType_BuildingIcon, int(protoId))
			button.GetWidget().Disabled = true
		}
		for _, protoId := range m.menuUnitPrototypeIds {
			button := m.GetButton(UI_ButtonType_UnitIcon, int(protoId))
			button.GetWidget().Disabled = true
		}
		return
	}
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
	// for _, protoId := range m.menuBuildingPrototypeIds {
	// 	proto := m.client.Game().GetBuildingPrototype(protoId)
	// 	button := m.GetButton(UI_ButtonType_BuildingIcon, int(protoId))
	// 	insufficientResourceCapacity := proto.GetResourceCost() > resourceCapacity
	// 	button.GetWidget().Disabled = insufficientResourceCapacity
	// }
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

func (m *UIManager) setComputeIndicators(playerId uint8) {
	player := m.client.Game().GetPlayer(playerId)

	computeDemand := int(player.GetComputeDemand())
	computeSupply := int(player.GetComputeSupply())
	computeSurplus := computeSupply - computeDemand
	computeSurplusMin0 := gen_utils.Max(0, computeSurplus)

	bar := m.GetProgressBar(UI_ProgressBar_Compute)
	bar.Max = int(computeSupply)
	bar.SetCurrent(int(computeSurplusMin0))

	text := fmt.Sprintf("%d/%d", computeSurplus, int(computeSupply))
	if computeSurplus < 0 {
		text += " (!)"
	}

	label := m.GetLabel(UI_ProgressBar_Compute)
	label.Label = text
}

// Sets the worker stat indicator.
func (m *UIManager) setWorkerIndicator(playerId uint8) {
	var (
		nWorkers         int
		nWorkersInactive int
		nWorkersIdle     int
	)
	m.client.Game().ForEachUnit(playerId, func(unitId uint8, unit *datamod.UnitsRow) {
		unitState := rts.UnitState(unit.GetState())
		protoId := unit.GetUnitType()
		proto := m.client.Game().GetUnitPrototype(protoId)
		if !proto.GetIsWorker() {
			return
		}
		switch unitState {
		case rts.UnitState_Active:
			nWorkers++
			if rts.WorkerCommandData(unit.GetCommand()).Type().IsIdle() {
				nWorkersIdle++
			}
		case rts.UnitState_Inactive:
			nWorkers++
			nWorkersInactive++
		}
	})
	workerCountLabel := m.GetLabel(UI_Label_WorkerCount)
	workerCountLabel.Label = fmt.Sprintf(
		""+
			"Inactive : %d\n"+
			"Idle     : %d\n"+
			"Busy     : %d\n"+
			"Total    : %d",
		nWorkersInactive, nWorkersIdle, nWorkers-nWorkersInactive-nWorkersIdle, nWorkers)
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

func newOverDisplayContainer(uim *UIManager) *widget.Container {
	rect := uim.client.coreRenderer.BoardDisplayRect()
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewAnchorLayout(
				widget.AnchorLayoutOpts.Padding(widget.Insets{Top: rect.Min.Y, Left: rect.Min.X}),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(rect.Dx(), rect.Dy()),
		),
	)
	uim.AddContainer(UI_Container_OverDisplay, container)
	return container
}

func newSideContainer(uim *UIManager) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(assets.UIBackgroundColor)),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, true}),
			widget.GridLayoutOpts.Spacing(BorderWidth, BorderWidth),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    BorderWidth * 3 / 2,
				Bottom: BorderWidth,
				Left:   BorderWidth,
				Right:  BorderWidth,
			}),
		)),
	)

	// container.AddChild(newInspectionView(uim))
	container.AddChild(newCenterMenu(uim))
	container.AddChild(newLogisticsView(uim))

	return container
}

func newCenterMenu(uim *UIManager) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(BorderWidth),
		)),
	)
	width := uim.client.coreRenderer.Config().ScreenSize.X - uim.client.coreRenderer.BoardDisplayRect().Max.X - BorderWidth*2
	// buildMenu := newBuildMenu(uim, uim.menuBuildingPrototypeIds, width)
	unitMenu := newUnitMenu(uim, uim.menuUnitPrototypeIds, width)
	resourceInfo := newResourceDisplay(uim, width)
	// computeInfo := newComputeDisplay(uim, width)
	// container.AddChild(buildMenu)
	container.AddChild(unitMenu)
	container.AddChild(resourceInfo)
	// container.AddChild(computeInfo)
	return container
}

func newBuildMenu(uim *UIManager, buildingPrototypeIds []uint8, width int) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(BorderWidth),
		)),
	)
	if len(buildingPrototypeIds) == 0 {
		return container
	}
	iconSize := (width - (BorderWidth * (len(buildingPrototypeIds) - 1))) / len(buildingPrototypeIds)
	for _, protoId := range buildingPrototypeIds {
		container.AddChild(newBuildingIcon(uim, protoId, iconSize))
	}
	return container
}

func newUnitMenu(uim *UIManager, unitPrototypeIds []uint8, width int) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(BorderWidth),
		)),
	)
	if len(unitPrototypeIds) == 0 {
		return container
	}
	iconSize := (width - (BorderWidth * (len(unitPrototypeIds) - 1))) / len(unitPrototypeIds)
	for _, protoId := range unitPrototypeIds {
		container.AddChild(newUnitIcon(uim, protoId, iconSize))
	}
	return container
}

func newBuildingIcon(uim *UIManager, protoId uint8, size int) *widget.Button {
	sprite := uim.spriteGetter.GetBuildingSprite(uim.client.PlayerId(), protoId, rts.BuildingState_Built)
	buttonImage := newIconButtonImage(uim, sprite, size, assets.UICornerSize*3/2, nil)
	button := newIconButton(
		buttonImage,
		uim.newButtonPressHandler(UI_ButtonType_BuildingIcon, int(protoId)),
		size,
	)
	uim.AddButton(UI_ButtonType_BuildingIcon, int(protoId), button)
	return button
}

func newUnitIcon(uim *UIManager, protoId uint8, size int) *widget.Button {
	sprite := uim.spriteGetter.GetUnitSprite(uim.client.PlayerId(), protoId, assets.Direction_Right)
	buttonImage := newIconButtonImage(uim, sprite, size, assets.UICornerSize, nil)
	button := newIconButton(
		buttonImage,
		uim.newButtonPressHandler(UI_ButtonType_UnitIcon, int(protoId)),
		size,
	)
	uim.AddButton(UI_ButtonType_UnitIcon, int(protoId), button)
	return button
}

func newIconButton(buttonImage *widget.ButtonImage, handler widget.ButtonPressedHandlerFunc, size int) *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(size, size),
		),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.PressedHandler(handler),
	)
	return button
}

func newIconButtonImage(uim *UIManager, sprite *ebiten.Image, size, margin int, alerts []*ebiten.Image) *widget.ButtonImage {
	img := ebiten.NewImage(size, size)
	imgBounds := img.Bounds()

	bgNineSlice := image.NewNineSliceSimple(assets.UIBox_LightConvex, assets.UICornerSize, assets.UICornerSize)
	bgNineSlice.Draw(img, size, size, nil)

	cornerBounds := go_image.Point{margin, margin}
	spriteRect := go_image.Rectangle{Min: cornerBounds, Max: imgBounds.Max.Sub(cornerBounds)}
	op := client_utils.NewDrawOptions(spriteRect, sprite.Bounds())
	colorm.DrawImage(img, sprite, colorm.ColorM{}, op)

	if len(alerts) > 0 {
		alertSize := 2 * 8
		maxAlerts := size / (alertSize + BorderWidth)
		for i, alert := range alerts {
			if i >= maxAlerts {
				break
			}
			position := go_image.Point{size - (i+1)*(alertSize+BorderWidth), size - alertSize - BorderWidth}
			alertRect := go_image.Rectangle{
				Min: position,
				Max: position.Add(go_image.Point{alertSize, alertSize}),
			}
			op := client_utils.NewDrawOptions(alertRect, alert.Bounds())
			colorM := colorm.ColorM{}
			colorM.Scale(1, 1, 1, 0.75)
			colorm.DrawImage(img, alert, colorM, op)
		}
	}

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

func newResourceDisplay(uim *UIManager, width int) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(BorderWidth),
		)),
	)
	container.AddChild(newProgressBar(uim, UI_ProgressBar_Resource, "minerals", assets.UIProgressBarRed, width))
	return container
}

func newComputeDisplay(uim *UIManager, width int) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(BorderWidth),
		)),
	)
	container.AddChild(newProgressBar(uim, UI_ProgressBar_Compute, "cpu", assets.UIProgressBarGreen, width))
	return container
}

func newProgressBar(uim *UIManager, id int, name string, sprite *ebiten.Image, width int) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
		),
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)
	resourceBar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(width, 0),
		),
		widget.ProgressBarOpts.Images(
			&widget.ProgressBarImage{
				Idle: image.NewNineSliceSimple(assets.UIBox_DarkConcave, assets.UICornerSize_ProgressBar, assets.UICornerSize_ProgressBar),
			},
			&widget.ProgressBarImage{
				Idle: image.NewNineSliceSimple(sprite, assets.UICornerSize_ProgressBar, assets.UICornerSize_ProgressBar),
			},
		),
		widget.ProgressBarOpts.Values(0, 1000, 0),
	)
	nameText := widget.NewText(
		widget.TextOpts.Text(name, assets.BitmapFont1, assets.LightGray),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
		widget.TextOpts.Insets(widget.Insets{Left: assets.UICornerSize, Right: assets.UICornerSize}),
	)
	barText := widget.NewText(
		widget.TextOpts.Text("", assets.BitmapFont1, assets.LightGray),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.Insets(widget.Insets{Left: assets.UICornerSize, Right: assets.UICornerSize}),
	)

	container.AddChild(resourceBar)
	container.AddChild(nameText)
	container.AddChild(barText)

	uim.AddProgressBar(id, resourceBar)
	uim.AddLabel(id, barText)

	return container
}

func newInspectionView(uim *UIManager) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(0, 96),
		),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceSimple(assets.UIBox_DarkConcave, assets.UICornerSize, assets.UICornerSize)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(BorderWidth),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(PanelPadding)),
		)),
	)
	inspectedNameLabel := widget.NewText(
		widget.TextOpts.Text("", assets.GetFontFace(assets.Font_PressStart, 8), assets.TextLightColor),
	)
	inspectedDetailsLabel := widget.NewText(
		widget.TextOpts.Text("", assets.BitmapFont1, assets.TextLightColor),
	)
	container.AddChild(inspectedNameLabel)
	container.AddChild(inspectedDetailsLabel)
	uim.AddLabel(UI_Label_InspectedName, inspectedNameLabel)
	uim.AddLabel(UI_Label_InspectedDetails, inspectedDetailsLabel)
	return container
}

func newLogisticsView(uim *UIManager) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceSimple(assets.UIBox_DarkConcave, assets.UICornerSize, assets.UICornerSize)),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, false}),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(PanelPadding)),
		)),
	)
	statsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{false}),
		)),
	)
	workerStatsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(4),
		)),
	)
	workerStatsTitleLabel := widget.NewText(
		widget.TextOpts.Text("Workers", assets.GetFontFace(assets.Font_PressStart, 8), assets.TextLightColor),
	)
	workerStatsLabel := widget.NewText(
		widget.TextOpts.Text("", assets.BitmapFont1, assets.TextLightColor),
	)

	bottomLine := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, false}, []bool{false, false, false}),
		)),
	)

	workerStatsContainer.AddChild(workerStatsTitleLabel)
	workerStatsContainer.AddChild(workerStatsLabel)
	statsContainer.AddChild(workerStatsContainer)

	bottomLine.AddChild(newSettingsText(uim, "Sync.     : "))
	bottomLine.AddChild(newSettingsLabel(uim, -1, boolToOnOff(uim.client.active)))
	bottomLine.AddChild(newSettingsText(uim, "Interpol. : "))
	bottomLine.AddChild(newSettingsLabel(uim, UI_Label_Interpolation, boolToOnOff(false)))
	bottomLine.AddChild(newSettingsText(uim, "FPS       : "))
	bottomLine.AddChild(newSettingsLabel(uim, UI_Label_FPS, "0.0"))

	container.AddChild(statsContainer)
	container.AddChild(bottomLine)

	uim.AddLabel(UI_Label_WorkerCount, workerStatsLabel)

	return container
}

func newSettingsLabel(uim *UIManager, id int, text string) *widget.Text {
	label := newSettingsText(uim, text)
	uim.AddLabel(id, label)
	return label
}

func newSettingsText(uim *UIManager, text string) *widget.Text {
	return widget.NewText(
		widget.TextOpts.Text(text, assets.BitmapFont1, assets.BoxTextDarkColor),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	)
}

func boolToOnOff(b bool) string {
	if b {
		return "ON"
	}
	return "OFF"
}
