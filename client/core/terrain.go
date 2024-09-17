package core

import (
	"image"

	"github.com/concrete-eth/ark-rts/client/assets"
	"github.com/concrete-eth/ark-rts/client/decren"
	"github.com/hajimehoshi/ebiten/v2"
)

func renderMap(mapId int) *ebiten.Image {
	m := assets.LoadMap(mapId)
	renderer, err := assets.NewMapRenderer(m)
	if err != nil {
		panic(err)
	}
	if err := renderer.RenderVisibleLayers(); err != nil {
		panic(err)
	}
	img := renderer.Result
	renderer.Clear()
	eimg := ebiten.NewImageFromImage(img)
	return eimg
}

// Initializes the terrain sprite layer by setting the background, borders, and decorative cracks.
func initTerrain(layer *decren.Layer, mapId int) {
	layer.SetBackgroundColor(assets.TerrainBackgroundColor)
	terrainImg := renderMap(mapId)
	terrainOrigin := image.Point{1, 0} // TODO: get from map data
	layer.Sprite("terrain").
		SetImage(terrainImg).
		FitToImage().
		SetPosition(terrainOrigin.Mul(-assets.TileSize))
}
