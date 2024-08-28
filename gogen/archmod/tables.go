/* Autogenerated file. Do not edit manually. */

package archmod

import (
	"reflect"

	"github.com/concrete-eth/archetype/arch"

	contract "github.com/concrete-eth/ark-rts/gogen/abigen/tables"
	"github.com/concrete-eth/ark-rts/gogen/datamod"
)

var TablesABIJson = contract.ContractABI

var TableSchemasJson = `{
    "meta": {
        "schema": {
            "boardWidth": "uint16",
            "boardHeight": "uint16",
            "playerCount": "uint8",
            "unitPrototypeCount": "uint8",
            "buildingPrototypeCount": "uint8",
            "isInitialized": "bool",
            "hasStarted": "bool",
            "creationBlockNumber": "uint32"
        }
    },
    "players": {
        "keySchema": {
            "playerId": "uint8"
        },
        "schema": {
            "spawnAreaX": "uint16",
            "spawnAreaY": "uint16",
            "spawnAreaWidth": "uint8",
            "spawnAreaHeight": "uint8",
            "workerPortX": "uint16",
            "workerPortY": "uint16",
            "curResource": "uint16",
            "maxResource": "uint16",
            "curArmories": "uint8",
            "computeSupply": "uint8",
            "computeDemand": "uint8",
            "unitCount": "uint8",
            "buildingCount": "uint8",
            "buildingPayQueuePointer": "uint8",
            "buildingBuildQueuePointer": "uint8",
            "unitPayQueuePointer": "uint8"
        }
    },
    "board": {
        "keySchema": {
            "x": "uint16",
            "y": "uint16"
        },
        "schema": {
            "landObjectType": "uint8",
            "landPlayerId": "uint8",
            "landObjectId": "uint8",
            "hoverPlayerId": "uint8",
            "hoverUnitId": "uint8",
            "airPlayerId": "uint8",
            "airUnitId": "uint8"
        }
    },
    "units": {
        "keySchema": {
            "playerId": "uint8",
            "unitId": "uint8"
        },
        "schema": {
            "x": "uint16",
            "y": "uint16",
            "unitType": "uint8",
            "state": "uint8",
            "load": "uint8",
            "integrity": "uint8",
            "timestamp": "uint32",
            "command": "uint64",
            "commandExtra": "uint64",
            "commandMeta": "uint8",
            "isPreTicked": "bool"
        }
    },
    "buildings": {
        "keySchema": {
            "playerId": "uint8",
            "buildingId": "uint8"
        },
        "schema": {
            "x": "uint16",
            "y": "uint16",
            "buildingType": "uint8",
            "state": "uint8",
            "integrity": "uint8",
            "timestamp": "uint32"
        }
    },
    "unitPrototypes": {
        "keySchema": {
            "unitType": "uint8"
        },
        "schema": {
            "layer": "uint8",
            "resourceCost": "uint16",
            "computeCost": "uint8",
            "spawnTime": "uint8",
            "maxIntegrity": "uint8",
            "landStrength": "uint8",
            "hoverStrength": "uint8",
            "airStrength": "uint8",
            "attackRange": "uint8",
            "attackCooldown": "uint8",
            "isAssault": "bool",
            "isWorker": "bool"
        }
    },
    "buildingPrototypes": {
        "keySchema": {
            "buildingType": "uint8"
        },
        "schema": {
            "width": "uint8",
            "height": "uint8",
            "resourceCost": "uint16",
            "resourceCapacity": "uint16",
            "computeCapacity": "uint8",
            "resourceMine": "uint8",
            "mineTime": "uint8",
            "maxIntegrity": "uint8",
            "buildingTime": "uint8",
            "isArmory": "bool",
            "isEnvironment": "bool"
        }
    }
}`

var TableSchemas arch.TableSchemas

func init() {
	types := map[string]reflect.Type{
		"Meta":               reflect.TypeOf(RowData_Meta{}),
		"Players":            reflect.TypeOf(RowData_Players{}),
		"Board":              reflect.TypeOf(RowData_Board{}),
		"Units":              reflect.TypeOf(RowData_Units{}),
		"Buildings":          reflect.TypeOf(RowData_Buildings{}),
		"UnitPrototypes":     reflect.TypeOf(RowData_UnitPrototypes{}),
		"BuildingPrototypes": reflect.TypeOf(RowData_BuildingPrototypes{}),
	}
	getters := map[string]interface{}{
		"Meta":               datamod.NewMeta,
		"Players":            datamod.NewPlayers,
		"Board":              datamod.NewBoard,
		"Units":              datamod.NewUnits,
		"Buildings":          datamod.NewBuildings,
		"UnitPrototypes":     datamod.NewUnitPrototypes,
		"BuildingPrototypes": datamod.NewBuildingPrototypes,
	}
	var err error
	if TableSchemas, err = arch.NewTableSchemasFromRaw(TablesABIJson, TableSchemasJson, types, getters); err != nil {
		panic(err)
	}
}
