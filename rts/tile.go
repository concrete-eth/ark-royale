package rts

import (
	"fmt"

	"github.com/concrete-eth/ark-rts/gogen/datamod"
)

func SetTileLandObject(tile *datamod.BoardRow, objectType ObjectType, playerId uint8, objectId uint8) {
	tile.SetLandObjectType(objectType.Uint8())
	tile.SetLandPlayerId(playerId)
	tile.SetLandObjectId(objectId)
}

func SetTileHoverUnit(tile *datamod.BoardRow, playerId uint8, unitId uint8) {
	tile.SetHoverPlayerId(playerId)
	tile.SetHoverUnitId(unitId)
}

func SetTileAirUnit(tile *datamod.BoardRow, playerId uint8, unitId uint8) {
	tile.SetAirPlayerId(playerId)
	tile.SetAirUnitId(unitId)
}

func SetTileUnit(tile *datamod.BoardRow, layer LayerId, playerId uint8, unitId uint8) {
	switch layer {
	case LayerId_Land:
		SetTileLandObject(tile, ObjectType_Unit, playerId, unitId)
	case LayerId_Hover:
		SetTileHoverUnit(tile, playerId, unitId)
	case LayerId_Air:
		SetTileAirUnit(tile, playerId, unitId)
	default:
		panic(fmt.Sprintf("invalid layer: %d", layer))
	}
}

func EmptyUnitTile(tile *datamod.BoardRow, layer LayerId) {
	switch layer {
	case LayerId_Land:
		tile.SetLandObjectType(ObjectType_Nil.Uint8())
		tile.SetLandPlayerId(NilPlayerId)
		tile.SetLandObjectId(NilObjectId)
	case LayerId_Hover:
		tile.SetHoverPlayerId(NilPlayerId)
		tile.SetHoverUnitId(NilUnitId)
	case LayerId_Air:
		tile.SetAirPlayerId(NilPlayerId)
		tile.SetAirUnitId(NilUnitId)
	default:
		panic(fmt.Sprintf("invalid layer: %d", layer))
	}
}

func IsTileEmpty(tile *datamod.BoardRow, layer LayerId) bool {
	switch layer {
	case LayerId_Land:
		return tile.GetLandObjectType() == ObjectType_Nil.Uint8()
	case LayerId_Hover:
		return tile.GetHoverUnitId() == NilPlayerId
	case LayerId_Air:
		return tile.GetAirUnitId() == NilPlayerId
	default:
		panic(fmt.Sprintf("invalid layer: %d", layer))
	}
}

func IsTileEmptyAllLayers(tile *datamod.BoardRow) bool {
	return IsTileEmpty(tile, LayerId_Land) && IsTileEmpty(tile, LayerId_Hover) && IsTileEmpty(tile, LayerId_Air)
}
