package assets

import (
	"fmt"
	"image"

	"github.com/concrete-eth/ark-rts/rts"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TileSize     = 16
	BuildingSize = 32
	UnitSize     = 24
)

const (
	BuildingSpriteId_Main = iota
	BuildingSpriteId_Storage
	BuildingSpriteId_Lab
	BuildingSpriteId_Armory
	BuildingSpriteId_Mine
	BuildingSpriteId_SmallMine
	BuildingSpriteId_Count
)

const (
	UnitSpriteId_AntiAir = iota
	UnitSpriteId_Air
	UnitSpriteId_Tank
	UnitSpriteId_Turret
	UnitSpriteId_Worker
	UnitSpriteId_Count
)

var (
	spriteSheet = LoadImage("sprite_sheet.png")

	MineSprite       = SubImage(spriteSheet, NewBounds(0, 0, 18, 16))
	SmallMinesSprite = SubImage(spriteSheet, NewBounds(0, 16, 16, 16))

	spawnPointSprites = [4]*ebiten.Image{
		SubImage(spriteSheet, NewBounds(160, 48, 32, 32)),
		SubImage(spriteSheet, NewBounds(160+192, 48, 32, 32)),
		SubImage(spriteSheet, NewBounds(160+2*192, 48, 32, 32)),
		SubImage(spriteSheet, NewBounds(160+3*192, 48, 32, 32)),
	}

	playerBuildingSprites = [4][BuildingSpriteId_Count][rts.BuildingState_Count]*ebiten.Image{
		loadPlayerBuildingSprites(SubImage(spriteSheet, NewBounds(0, 48, 128, 64))),
		loadPlayerBuildingSprites(SubImage(spriteSheet, NewBounds(192, 48, 128, 64))),
		loadPlayerBuildingSprites(SubImage(spriteSheet, NewBounds(2*192, 48, 128, 64))),
		loadPlayerBuildingSprites(SubImage(spriteSheet, NewBounds(3*192, 48, 128, 64))),
	}

	UnitSpriteOrigin = image.Point{4, 4}

	unitSprites = [4][UnitSpriteId_Count][8][3]*ebiten.Image{
		loadPlayerUnitSprites(SubImage(spriteSheet, NewBounds(0, 112, 192, 312))),
		loadPlayerUnitSprites(SubImage(spriteSheet, NewBounds(192, 112, 192, 312))),
		loadPlayerUnitSprites(SubImage(spriteSheet, NewBounds(2*192, 112, 192, 312))),
		loadPlayerUnitSprites(SubImage(spriteSheet, NewBounds(3*192, 112, 192, 312))),
	}

	SelectionSprite = SubImage(spriteSheet, NewBounds(16, 16, 16, 16))
)

var (
	uiSpriteSheet = SubImage(spriteSheet, NewBounds(32, 0, 12, 30))

	UICornerSize_Small = 2
	UICornerSize_Big   = 5

	UIBox_Small           = SubImage(uiSpriteSheet, NewBounds(0, 0, 6, 6))
	UIBox_Big             = SubImage(uiSpriteSheet, NewBounds(0, 6, 12, 12))
	UIPanel_Convex        = SubImage(uiSpriteSheet, NewBounds(6, 0, 6, 6))
	UIProgressBar_Track   = SubImage(uiSpriteSheet, NewBounds(0, 18, 6, 6))
	UIProgressBar_Mineral = SubImage(uiSpriteSheet, NewBounds(0, 24, 6, 6))
	UIProgressBar_Compute = SubImage(uiSpriteSheet, NewBounds(6, 24, 6, 6))
)

func validatePlayerId(playerId uint8) {
	if playerId == 0 && playerId > 4 {
		panic("invalid player id")
	}
}

func validateDirection(direction Direction) {
	if !direction.IsValid() {
		panic("invalid direction")
	}
}

func validateUnitSpriteId(spriteId uint8) {
	if spriteId >= UnitSpriteId_Count {
		panic("invalid unit sprite id")
	}
}

func validateBuildingSpriteId(spriteId uint8) {
	if spriteId >= BuildingSpriteId_Count {
		panic("invalid building sprite id")
	}
}

type SpriteGetter interface {
	GetSpawnPointSprite(playerId uint8) *ebiten.Image
	GetBuildingSprite(playerId uint8, spriteId uint8, buildingState rts.BuildingState) *ebiten.Image
	GetBuildingSpriteOrigin(spriteId uint8) image.Point
	GetUnitSprite(playerId uint8, spriteId uint8, direction Direction) *ebiten.Image
	GetUnitFireFrame(playerId uint8, spriteId uint8, direction Direction, frame uint) *ebiten.Image
}

var _ SpriteGetter = (*DefaultSpriteGetter)(nil)

type DefaultSpriteGetter struct{}

func (d *DefaultSpriteGetter) GetSpawnPointSprite(playerId uint8) *ebiten.Image {
	validatePlayerId(playerId)
	return spawnPointSprites[playerId-1]
}

func (d *DefaultSpriteGetter) GetBuildingSprite(playerId uint8, spriteId uint8, buildingState rts.BuildingState) *ebiten.Image {
	validatePlayerId(playerId)
	validateBuildingSpriteId(spriteId)
	switch {
	case spriteId == BuildingSpriteId_SmallMine:
		return SmallMinesSprite
	case spriteId == BuildingSpriteId_Mine:
		return MineSprite
	default:
		return playerBuildingSprites[playerId-1][spriteId][buildingState]
	}
}

func (d *DefaultSpriteGetter) GetBuildingSpriteOrigin(spriteId uint8) image.Point {
	if spriteId == BuildingSpriteId_Mine {
		return image.Point{-1, 0}
	}
	return image.Point{0, 0}
}

func (d *DefaultSpriteGetter) GetUnitSprite(playerId uint8, spriteId uint8, direction Direction) *ebiten.Image {
	validatePlayerId(playerId)
	validateDirection(direction)
	return unitSprites[playerId-1][spriteId][direction][0]
}

func (d *DefaultSpriteGetter) GetUnitFireFrame(playerId uint8, spriteId uint8, direction Direction, frame uint) *ebiten.Image {
	validatePlayerId(playerId)
	validateUnitSpriteId(spriteId)
	validateDirection(direction)
	if int(frame) >= len(unitSprites[0][0][0])-1 {
		panic(fmt.Sprintf("invalid frame %v", frame))
	}
	return unitSprites[playerId-1][spriteId][direction][frame+1]
}

func loadPlayerBuildingSprites(spriteSheet *ebiten.Image) [BuildingSpriteId_Count][rts.BuildingState_Count]*ebiten.Image {
	sprites := [BuildingSpriteId_Count][rts.BuildingState_Count]*ebiten.Image{}
	for ii, buildingType := range []int{
		BuildingSpriteId_Main,
		BuildingSpriteId_Lab,
		BuildingSpriteId_Armory,
		BuildingSpriteId_Storage,
	} {
		sprites[buildingType] = loadBuildingStateSprites(SubImage(spriteSheet, NewBounds(ii*BuildingSize, 0, BuildingSize, 2*BuildingSize)))
	}
	return sprites
}

func loadBuildingStateSprites(spriteSheet *ebiten.Image) [rts.BuildingState_Count]*ebiten.Image {
	sprites := [rts.BuildingState_Count]*ebiten.Image{}
	sprites[rts.BuildingState_Built] = SubImage(spriteSheet, NewBounds(0, 0, BuildingSize, BuildingSize))
	sprites[rts.BuildingState_Building] = SubImage(spriteSheet, NewBounds(0, BuildingSize, BuildingSize, BuildingSize))
	return sprites
}

func loadPlayerUnitSprites(spriteSheet *ebiten.Image) [UnitSpriteId_Count][8][3]*ebiten.Image {
	sprites := [UnitSpriteId_Count][8][3]*ebiten.Image{}
	for ii, unitType := range []int{
		UnitSpriteId_AntiAir,
		UnitSpriteId_Air,
		UnitSpriteId_Tank,
		UnitSpriteId_Turret,
		UnitSpriteId_Worker,
	} {
		sprites[unitType] = loadUnitDirectionSpriteSets(SubImage(spriteSheet, NewBounds(0, 3*ii*UnitSize, 8*UnitSize, 3*UnitSize)))
	}
	return sprites
}

func loadUnitDirectionSpriteSets(spriteSheet *ebiten.Image) [8][3]*ebiten.Image {
	sprites := [8][3]*ebiten.Image{}
	for ii, direction := range []Direction{
		Direction_Left,
		Direction_UpLeft,
		Direction_Up,
		Direction_UpRight,
		Direction_Right,
		Direction_DownRight,
		Direction_Down,
		Direction_DownLeft,
	} {
		sprites[direction] = loadUnitDirectionSprites(SubImage(spriteSheet, NewBounds(ii*UnitSize, 0, UnitSize, 3*UnitSize)))
	}
	return sprites
}

func loadUnitDirectionSprites(spriteSheet *ebiten.Image) [3]*ebiten.Image {
	sprites := [3]*ebiten.Image{}
	sprites[0] = SubImage(spriteSheet, NewBounds(0, 0, UnitSize, UnitSize))
	sprites[1] = SubImage(spriteSheet, NewBounds(0, UnitSize, UnitSize, UnitSize))
	sprites[2] = SubImage(spriteSheet, NewBounds(0, 2*UnitSize, UnitSize, UnitSize))
	return sprites
}

type Direction int

const (
	Direction_Up Direction = iota
	Direction_UpRight
	Direction_Right
	Direction_DownRight
	Direction_Down
	Direction_DownLeft
	Direction_Left
	Direction_UpLeft
)

func (d Direction) IsValid() bool {
	return d >= 0 && d < 8
}

func (d Direction) IsVertical() bool {
	return d == Direction_Up || d == Direction_Down
}

func (d Direction) IsHorizontal() bool {
	return d == Direction_Left || d == Direction_Right
}

func (d Direction) IsStraight() bool {
	return d.IsVertical() || d.IsHorizontal()
}

func (d Direction) IsDiagonal() bool {
	return d == Direction_UpRight || d == Direction_DownRight || d == Direction_DownLeft || d == Direction_UpLeft
}
