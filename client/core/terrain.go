package core

import (
	"image"
	"math"
	"strings"

	"github.com/concrete-eth/ark-royale/client/assets"
	"github.com/concrete-eth/ark-royale/client/decren"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

func renderTerrain(mapId int) (*ebiten.Image, image.Point) {
	m := assets.LoadMap(mapId)
	renderer, err := assets.NewMapRenderer(m)
	if err != nil {
		panic(err)
	}

	var terrainLayer *tiled.Layer

	for id, layer := range m.Layers {
		namePrefix := strings.ToLower(strings.Split(layer.Name, "_")[0])
		switch namePrefix {
		case "deco", "terrain", "side":
			if err := renderer.RenderLayer(id); err != nil {
				panic(err)
			}
			if namePrefix == "terrain" {
				terrainLayer = layer
			}
		default:
		}
	}

	if terrainLayer == nil {
		panic("terrain layer not found")
	}

	img := renderer.Result
	renderer.Clear()
	eimg := ebiten.NewImageFromImage(img)

	origin := image.Point{math.MaxInt16, math.MaxInt16}
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			tile := terrainLayer.Tiles[x+y*m.Width]
			if !tile.IsNil() {
				if x < origin.X {
					origin.X = x
				}
				if y < origin.Y {
					origin.Y = y
				}
			}
		}
	}

	return eimg, origin
}

// Initializes the terrain sprite layer by setting the background, borders, and decorative cracks.
func initTerrain(layer *decren.Layer, mapId int) {
	layer.SetBackgroundColor(assets.TerrainBackgroundColor)
	terrainImg, terrainOrigin := renderTerrain(mapId)
	layer.Sprite("terrain").
		SetImage(terrainImg).
		FitToImage().
		SetPosition(terrainOrigin.Mul(-assets.TileSize))
}
