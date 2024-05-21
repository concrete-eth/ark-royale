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

	p := perlin.NewPerlin(2, 2, 2, 0)
	for x := 0; x < sizeInTiles.X; x++ {
		for y := 0; y < sizeInTiles.Y; y++ {
			var img *ebiten.Image
			if x == 0 || y == 0 || x == sizeInTiles.X-1 || y == sizeInTiles.Y-1 {
				img = assets.CrackTileSprite
			} else {
				noise := p.Noise2D(float64(x)/3.0, float64(y)/3.0)
				if noise > 0.15 {
					img = assets.CrackTileSprite
				}
			}
			if img != nil {
				setTerrainTile(layer, image.Point{x, y}, img)
			}
		}
	}

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
}
