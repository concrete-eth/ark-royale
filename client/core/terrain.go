package core

import (
	"image"

	"github.com/concrete-eth/ark-rts/client/assets"
	"github.com/concrete-eth/ark-rts/client/decren"
	"github.com/concrete-eth/ark-rts/client/perlin"
	"github.com/hajimehoshi/ebiten/v2"
)

// Sets the image of a tile in the terrain layer.
func setTerrainTile(layer *decren.Layer, pos image.Point, img *ebiten.Image) {
	layer.Sprite(pos.X, pos.Y).
		SetPosition(pos.Mul(assets.TileSize)).
		SetImage(img).
		FitToImage()
}

// Initializes the terrain sprite layer by setting the background, borders, and decorative cracks.
func initTerrain(layer *decren.Layer, sizeInTiles image.Point) {
	layer.SetBackgroundColor(assets.TerrainBackgroundColor)

	bg := ebiten.NewImage(1, 1)
	bg.Fill(assets.GroundColor)
	layer.Sprite("ground").
		SetImage(bg).
		SetSize(sizeInTiles.Mul(assets.TileSize))

	// pathImg := ebiten.NewImage(assets.TileSize, assets.TileSize)
	// pathImg.Fill(assets.ChangeHSV(assets.GroundColor, 0, 0.9, 0.75))
	pathSet := map[image.Point]struct{}{}

	for x := 0; x < sizeInTiles.X; x++ {
		if x == 1 || x == sizeInTiles.X-2 {
			for y := 1; y < sizeInTiles.Y-1; y++ {
				setTerrainTile(layer, image.Point{x, y}, assets.BrickTileSprite)
				pathSet[image.Point{x, y}] = struct{}{}
			}
		} else if x > 1 && x < sizeInTiles.X-2 {
			setTerrainTile(layer, image.Point{x, 1}, assets.BrickTileSprite)
			setTerrainTile(layer, image.Point{x, sizeInTiles.Y - 2}, assets.BrickTileSprite)
			pathSet[image.Point{x, 1}] = struct{}{}
			pathSet[image.Point{x, sizeInTiles.Y - 2}] = struct{}{}
		}
	}

	p := perlin.NewPerlin(2, 2, 2, 0)
	for x := 0; x < sizeInTiles.X; x++ {
		for y := 0; y < sizeInTiles.Y; y++ {
			if y == sizeInTiles.Y/2 {
				continue
			}
			if _, ok := pathSet[image.Point{x, y}]; ok {
				continue
			}
			var img *ebiten.Image
			noise := p.Noise2D(float64(x)/3.0, float64(y)/3.0)
			if noise > 0.075 {
				img = assets.CrackTileSprite
			}
			if img != nil {
				setTerrainTile(layer, image.Point{x, y}, img)
			}
		}
	}

	pitY := sizeInTiles.Y / 2

	for x := 0; x < sizeInTiles.X; x++ {
		var y int
		var img *ebiten.Image

		y = -1
		img = assets.BorderTileSet[assets.Direction_Up]
		setTerrainTile(layer, image.Point{x, y}, img)

		y = sizeInTiles.Y
		img = assets.BorderTileSet[assets.Direction_Down]
		setTerrainTile(layer, image.Point{x, y}, img)
	}

	for y := 0; y < sizeInTiles.Y; y++ {
		if y == pitY {
			continue
		}

		var x int
		var img *ebiten.Image

		x = -1
		img = assets.BorderTileSet[assets.Direction_Left]
		setTerrainTile(layer, image.Point{x, y}, img)

		x = sizeInTiles.X
		img = assets.BorderTileSet[assets.Direction_Right]
		setTerrainTile(layer, image.Point{x, y}, img)
	}

	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			dir := assets.DirectionFromDelta(image.Point{x*2 - 1, y*2 - 1})
			img := assets.BorderTileSet[dir]
			pos := image.Point{-1 + x*(sizeInTiles.X+1), -1 + y*(sizeInTiles.Y+1)}
			setTerrainTile(layer, pos, img)
		}
	}

	// pitTilePos := image.Point{0, pitY}
	leftPitOuterEdgePos := image.Point{0, pitY}
	rightPitOuterEdgePos := image.Point{sizeInTiles.X - 1, pitY}
	setTerrainTile(layer, leftPitOuterEdgePos, assets.PitTileSet[2])
	setTerrainTile(layer, rightPitOuterEdgePos, assets.PitTileSet[0])

	leftPitBorderPos := image.Point{-1, pitY}
	rightPitBorderPos := image.Point{sizeInTiles.X, pitY}
	setTerrainTile(layer, leftPitBorderPos, assets.PitTileSet[3])
	setTerrainTile(layer, rightPitBorderPos, assets.PitTileSet[4])

	leftPitInnerEdgePos := image.Point{2, pitY}
	rightPitInnerEdgePos := image.Point{sizeInTiles.X - 3, pitY}
	setTerrainTile(layer, leftPitInnerEdgePos, assets.PitTileSet[0])
	setTerrainTile(layer, rightPitInnerEdgePos, assets.PitTileSet[2])

	for pitX := 3; pitX < sizeInTiles.X-3; pitX++ {
		pitTilePos := image.Point{pitX, pitY}
		setTerrainTile(layer, pitTilePos, assets.PitTileSet[1])
	}
}
