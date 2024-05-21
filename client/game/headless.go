package game

import (
	"image"

	"github.com/concrete-eth/ark-rts/client/core"
	"github.com/concrete-eth/ark-rts/rts"
	"github.com/ethereum/go-ethereum/log"
)

type HeadlessClient struct {
	core.IHeadlessClient
}

func (c *HeadlessClient) PlaceBuilding(buildingType uint8, position image.Point) {
	if buildingType == BuildingPrototypeId_Main {
		log.Error("Cannot place main building")
		return
	}
	var (
		proto     = c.Game().GetBuildingPrototype(buildingType)
		size      = rts.GetDimensionsAsPoint(proto)
		buildArea = image.Rectangle{position, position.Add(size)}
	)
	if !IsInPlayerBuildableArea(c.Game(), c.PlayerId(), buildArea) {
		log.Error("Cannot place building outside of player's build area")
		return
	}
	c.IHeadlessClient.PlaceBuilding(buildingType, position)
}

func (c *HeadlessClient) CreateUnit(unitType uint8) {
	hasArmory := c.Game().GetPlayer(c.PlayerId()).GetCurArmories() > 0
	if !hasArmory && NeedsArmory(unitType) {
		log.Error("Cannot create unit without armory")
		return
	}
	c.IHeadlessClient.CreateUnit(unitType)
}
