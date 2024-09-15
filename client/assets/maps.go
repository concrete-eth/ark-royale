package assets

import (
	"embed"
	"path/filepath"

	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

//go:embed all:maps/assets
var mapsFS embed.FS

const (
	MapTilesetId_Royale = iota
)

var tilemaps map[int]*tiled.Map

func init() {
	tilemaps = make(map[int]*tiled.Map)
}

func LoadMap(mapId int) *tiled.Map {
	if m, ok := tilemaps[mapId]; ok {
		return m
	}
	var name string
	switch mapId {
	case MapTilesetId_Royale:
		name = "royale-map.tmx"
	default:
		panic("unknown map id")
	}
	path := filepath.Join("maps", "assets", name)
	withFs := tiled.WithFileSystem(mapsFS)
	m, err := tiled.LoadFile(path, withFs)
	if err != nil {
		panic(err)
	}
	tilemaps[mapId] = m
	return m
}

func NewMapRenderer(m *tiled.Map) (*render.Renderer, error) {
	return render.NewRendererWithFileSystem(m, mapsFS)
}
