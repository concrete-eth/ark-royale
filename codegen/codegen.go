package codegen

import (
	_ "embed"
	"html/template"
	"image"
	"math"
	"sort"
	"strings"

	"github.com/concrete-eth/archetype/codegen"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/lafriks/go-tiled"
	"github.com/spf13/cobra"
)

//go:embed templates/board.sol.tpl
var boardTpl string

const (
	BuildingPrototypeId_Main uint8 = iota + 1
	BuildingPrototypeId_Pit
	BuildingPrototypeId_Mine
)

const (
	UnitPrototypeId_AntiAir uint8 = iota + 1
	UnitPrototypeId_Air
	UnitPrototypeId_Tank
	UnitPrototypeId_Worker
	UnitPrototypeId_Turret
)

type Point struct {
	X, Y int
}

type Object struct {
	PrototypeId uint8
	Position    Point
}

type Building struct {
	Object
}

type Unit struct {
	Object
	Command uint64
}

type Player struct {
	Id         uint8
	Buildings  []*Building
	Units      []*Unit
	SpawnArea  image.Rectangle
	WorkerPort image.Point
}

func getLayerByPrefix(m *tiled.Map, prefix string) *tiled.Layer {
	prefix = strings.ToLower(prefix)
	for _, layer := range m.Layers {
		if strings.HasPrefix(strings.ToLower(layer.Name), prefix) {
			return layer
		}
	}
	return nil
}

func getObjGroupByPrefix(m *tiled.Map, prefix string) *tiled.ObjectGroup {
	prefix = strings.ToLower(prefix)
	for _, group := range m.ObjectGroups {
		if strings.HasPrefix(strings.ToLower(group.Name), prefix) {
			return group
		}
	}
	return nil
}

func runBoardLibGen(cmd *cobra.Command, args []string) {
	mapPath, err := cmd.Flags().GetString("map")
	if err != nil {
		panic(err)
	}
	if mapPath == "" {
		panic("map path is required")
	}
	m, err := tiled.LoadFile(mapPath)
	if err != nil {
		panic(err)
	}

	outPath, err := cmd.Flags().GetString("out")
	if err != nil {
		panic(err)
	}

	terrainLayer := getLayerByPrefix(m, "Terrain")
	if terrainLayer == nil {
		panic("terrain layer not found")
	}
	buildingLayer := getLayerByPrefix(m, "Buildings")
	unitLayer := getLayerByPrefix(m, "Units")
	spawnAreasGroup := getObjGroupByPrefix(m, "SpawnAreas")
	workerPortsGroup := getObjGroupByPrefix(m, "WorkerPorts")

	minX, minY := m.Width, m.Height
	maxX, maxY := 0, 0

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tile := terrainLayer.Tiles[x+y*m.Width]
			if !tile.IsNil() {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	players := make(map[uint8]*Player, 0)
	players[0] = &Player{
		Id:        0,
		Buildings: []*Building{},
		Units:     []*Unit{},
	}

	for x := minX; x < maxX+1; x++ {
		for y := minY; y < maxY+1; y++ {
			tile := terrainLayer.Tiles[x+y*m.Width]
			bX, bY := x-minX, y-minY
			if tile.IsNil() {
				players[0].Buildings = append(players[0].Buildings, &Building{
					Object: Object{
						PrototypeId: BuildingPrototypeId_Pit,
						Position:    Point{X: bX, Y: bY},
					},
				})
			}
		}
	}

	if buildingLayer != nil {
		for x := minX; x < maxX+1; x++ {
			for y := minY; y < maxY+1; y++ {
				tile := buildingLayer.Tiles[x+y*m.Width]
				bX, bY := x-minX, y-minY-1
				if !tile.IsNil() {
					tileId := tile.ID
					tileRect := tile.Tileset.GetTileRect(tileId)
					playerId := uint8(tileRect.Min.Y / tile.Tileset.TileHeight)
					player, ok := players[playerId]
					if !ok {
						player = &Player{
							Id:        playerId,
							Buildings: []*Building{},
							Units:     []*Unit{},
						}
						players[playerId] = player
					}
					var building *Building
					if playerId == 0 {
						protoId := BuildingPrototypeId_Mine
						building = &Building{
							Object: Object{
								PrototypeId: protoId,
								Position:    Point{X: bX, Y: bY},
							},
						}
						player.Buildings = append(player.Buildings, building)
					} else {
						protoId := uint8(tileRect.Min.X/tile.Tileset.TileWidth) + 1
						building = &Building{
							Object: Object{
								PrototypeId: protoId,
								Position:    Point{X: bX, Y: bY},
							},
						}
						if protoId == BuildingPrototypeId_Main {
							player.Buildings = append([]*Building{building}, player.Buildings...)
						} else {
							player.Buildings = append(player.Buildings, building)
						}
					}
				}
			}
		}
	}

	if unitLayer != nil {
		for x := minX; x < maxX+1; x++ {
			for y := minY; y < maxY+1; y++ {
				tile := unitLayer.Tiles[x+y*m.Width]
				bX, bY := x-minX, y-minY
				if !tile.IsNil() {
					tileId := tile.ID
					tileRect := tile.Tileset.GetTileRect(tileId)
					playerId := uint8(tileRect.Min.Y/tile.Tileset.TileHeight) + 1
					protoId := uint8(tileRect.Min.X/tile.Tileset.TileWidth) + 1
					player, ok := players[playerId]
					if !ok {
						player = &Player{
							Id:        playerId,
							Buildings: []*Building{},
							Units:     []*Unit{},
						}
						players[playerId] = player
					}
					player.Units = append(player.Units, &Unit{
						Object: Object{
							PrototypeId: protoId,
							Position:    Point{X: bX, Y: bY},
						},
					})
				}
			}
		}

		workerCommandType := rts.WorkerCommandType_Idle
		if environment, ok := players[0]; ok {
			for _, building := range environment.Buildings {
				if building.PrototypeId == BuildingPrototypeId_Mine {
					workerCommandType = rts.WorkerCommandType_Gather
					break
				}
			}
		}
		for playerId, player := range players {
			if playerId == 0 {
				continue
			}
			for _, unit := range player.Units {
				var rawCommand uint64
				if unit.PrototypeId == UnitPrototypeId_Worker {
					command := rts.NewWorkerCommandData(workerCommandType)
					if workerCommandType == rts.WorkerCommandType_Gather {
						var mineBuildingId uint8
						distance := math.MaxInt32
						for buildingIdx, building := range players[0].Buildings {
							if building.PrototypeId == BuildingPrototypeId_Mine {
								d := rts.DistanceToArea(
									image.Point{unit.Position.X, unit.Position.Y},
									image.Rectangle{
										Min: image.Point{building.Position.X, building.Position.Y},
										Max: image.Point{building.Position.X + 1, building.Position.Y + 1},
									},
								)
								if d < distance {
									distance = d
									mineBuildingId = uint8(buildingIdx) + 1
								}
							}
						}
						command.SetTargetBuildingId(mineBuildingId)
					}
					rawCommand = command.Uint64()
				} else {
					command := rts.NewFighterCommandData(rts.FighterCommandType_HoldPosition)
					command.SetTargetPosition(image.Point{unit.Position.X, unit.Position.Y})
					rawCommand = command.Uint64()
				}
				unit.Command = rawCommand
			}
		}
	}

	for playerIdx := 1; playerIdx < len(players); playerIdx++ {
		if _, ok := players[uint8(playerIdx)]; !ok {
			for nextPlayerIdx := playerIdx + 1; nextPlayerIdx < len(players)+1; nextPlayerIdx++ {
				if _, ok := players[uint8(nextPlayerIdx)]; ok {
					players[uint8(playerIdx)] = players[uint8(nextPlayerIdx)]
					delete(players, uint8(nextPlayerIdx))
					break
				}
			}
		}
	}

	if spawnAreasGroup != nil {
		for _, obj := range spawnAreasGroup.Objects {
			playerId := uint8(obj.Properties.GetInt("PlayerId"))
			if playerId == 0 {
				panic("PlayerId is required for SpawnAreas")
			}
			player, ok := players[playerId]
			if !ok {
				player = &Player{
					Id:        playerId,
					Buildings: []*Building{},
					Units:     []*Unit{},
				}
				players[playerId] = player
			}
			x := int(obj.X) / m.TileWidth
			y := int(obj.Y) / m.TileHeight
			w := int(obj.Width) / m.TileWidth
			h := int(obj.Height) / m.TileHeight
			player.SpawnArea = image.Rectangle{
				Min: image.Point{X: x - minX, Y: y - minY},
				Max: image.Point{X: x - minX + w, Y: y - minY + h},
			}
		}
	}

	if workerPortsGroup != nil {
		for _, obj := range workerPortsGroup.Objects {
			x := int(obj.X) / m.TileWidth
			y := int(obj.Y) / m.TileHeight
			position := image.Point{X: x - minX, Y: y - minY}
			var match bool
			for playerId, player := range players {
				if playerId == 0 {
					continue
				}
				if len(player.Buildings) == 0 {
					continue
				}
				firstBuilding := player.Buildings[0]
				if firstBuilding.PrototypeId != BuildingPrototypeId_Main {
					continue
				}
				mainBuilding := firstBuilding
				mainBuildingArea := image.Rectangle{
					Min: image.Point{mainBuilding.Position.X, mainBuilding.Position.Y},
					Max: image.Point{mainBuilding.Position.X + 2, mainBuilding.Position.Y + 2},
				}
				if !position.In(mainBuildingArea) {
					continue
				}
				player.WorkerPort = position
				match = true
				break
			}
			if !match {
				panic("WorkerPort must be inside a main building area")
			}
		}
	}

	for _, player := range players {
		sort.Slice(player.Buildings, func(i, j int) bool {
			return player.Buildings[i].PrototypeId < player.Buildings[j].PrototypeId
		})
		sort.Slice(player.Units, func(i, j int) bool {
			return player.Units[i].PrototypeId < player.Units[j].PrototypeId
		})
	}

	// b, err := json.MarshalIndent(players, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }
	// println(string(b))

	data := map[string]interface{}{
		"Map":     m,
		"OriginX": minX,
		"OriginY": minY,
		"Width":   maxX - minX + 1,
		"Height":  maxY - minY + 1,
		"Players": players,
		"Imports": []string{
			"./IActions.sol",
			"./ICore.sol",
		},
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	if err := codegen.ExecuteTemplate(boardTpl, "", outPath, data, funcMap); err != nil {
		panic(err)
	}
}

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{Use: "ark-codegen"}
	boardLibGenCmd := &cobra.Command{Use: "board", Short: "Generate a solidity board library", Run: runBoardLibGen}
	// boardLibGenCmd.Flags().Uint8P("pit", "p", 1, "pit building prototype id")
	boardLibGenCmd.Flags().StringP("out", "o", "./BoardLib.sol", "output file path")
	boardLibGenCmd.Flags().StringP("map", "m", "", "input map file path")
	rootCmd.AddCommand(boardLibGenCmd)
	return rootCmd
}

func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
