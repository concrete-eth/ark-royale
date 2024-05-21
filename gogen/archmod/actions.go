/* Autogenerated file. Do not edit manually. */

package archmod

import (
	"reflect"

	"github.com/concrete-eth/archetype/arch"

	contract "github.com/concrete-eth/ark-rts/gogen/abigen/actions"
)

var ActionsABIJson = contract.ContractABI

var ActionSchemasJson = `{
    "initialize": {
        "schema": {
            "width": "uint16",
            "height": "uint16"
        }
    },
    "start": {
        "schema": {}
    },
    "createUnit": {
        "schema": {
            "playerId": "uint8",
            "unitType": "uint8"
        }
    },
    "assignUnit": {
        "schema": {
            "playerId": "uint8",
            "unitId": "uint8",
            "command": "uint64",
            "commandExtra": "uint64",
            "commandMeta": "uint8"
        }
    },
    "placeBuilding": {
        "schema": {
            "playerId": "uint8",
            "buildingType": "uint8",
            "x": "uint16",
            "y": "uint16"
        }
    },
    "addPlayer": {
        "schema": {
            "spawnAreaX": "uint16",
            "spawnAreaY": "uint16",
            "workerPortX": "uint16",
            "workerPortY": "uint16"
        }
    },
    "addUnitPrototype": {
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
    "addBuildingPrototype": {
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

var ActionSchemas arch.ActionSchemas

func init() {
	types := map[string]reflect.Type{
		"Initialize":           reflect.TypeOf(ActionData_Initialize{}),
		"Start":                reflect.TypeOf(ActionData_Start{}),
		"CreateUnit":           reflect.TypeOf(ActionData_CreateUnit{}),
		"AssignUnit":           reflect.TypeOf(ActionData_AssignUnit{}),
		"PlaceBuilding":        reflect.TypeOf(ActionData_PlaceBuilding{}),
		"AddPlayer":            reflect.TypeOf(ActionData_AddPlayer{}),
		"AddUnitPrototype":     reflect.TypeOf(ActionData_AddUnitPrototype{}),
		"AddBuildingPrototype": reflect.TypeOf(ActionData_AddBuildingPrototype{}),
		// "Tick": reflect.TypeOf(arch.CanonicalTickAction{}),
	}
	var err error
	if ActionSchemas, err = arch.NewActionSchemasFromRaw(ActionsABIJson, ActionSchemasJson, types); err != nil {
		panic(err)
	}
}

type IActions interface {
	Initialize(action *ActionData_Initialize) error
	Start(action *ActionData_Start) error
	CreateUnit(action *ActionData_CreateUnit) error
	AssignUnit(action *ActionData_AssignUnit) error
	PlaceBuilding(action *ActionData_PlaceBuilding) error
	AddPlayer(action *ActionData_AddPlayer) error
	AddUnitPrototype(action *ActionData_AddUnitPrototype) error
	AddBuildingPrototype(action *ActionData_AddBuildingPrototype) error
	Tick()
}
