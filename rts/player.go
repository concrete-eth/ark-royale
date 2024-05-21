package rts

import (
	"github.com/concrete-eth/archetype/utils"
	"github.com/concrete-eth/ark-rts/gogen/datamod"
)

func addStorage(player *datamod.PlayersRow, resource uint16) {
	if resource > 0 {
		maxResource := player.GetMaxResource()
		player.SetMaxResource(utils.SafeAddUint16(maxResource, resource))
	}
}

func subStorage(player *datamod.PlayersRow, resource uint16) {
	if resource > 0 {
		maxResource := player.GetMaxResource()
		player.SetMaxResource(utils.SafeSubUint16(maxResource, resource))
		if player.GetCurResource() > maxResource {
			player.SetCurResource(maxResource)
		}
	}
}

func addResource(player *datamod.PlayersRow, resource uint16) {
	if resource > 0 {
		maxResource := player.GetMaxResource()
		curResource := player.GetCurResource()
		player.SetCurResource(utils.Min(utils.SafeAddUint16(curResource, resource), maxResource))
	}
}

func subResource(player *datamod.PlayersRow, resource uint16) {
	if resource > 0 {
		curResource := player.GetCurResource()
		player.SetCurResource(utils.SafeSubUint16(curResource, resource))
	}
}

func addComputeSupply(player *datamod.PlayersRow, compute uint8) {
	if compute > 0 {
		curCompute := player.GetComputeSupply()
		player.SetComputeSupply(utils.SafeAddUint8(curCompute, compute))
	}
}

func subComputeSupply(player *datamod.PlayersRow, compute uint8) {
	if compute > 0 {
		curCompute := player.GetComputeSupply()
		player.SetComputeSupply(utils.SafeSubUint8(curCompute, compute))
	}
}

func addComputeDemand(player *datamod.PlayersRow, compute uint8) {
	if compute > 0 {
		curCompute := player.GetComputeDemand()
		player.SetComputeDemand(utils.SafeAddUint8(curCompute, compute))
	}
}

func subComputeDemand(player *datamod.PlayersRow, compute uint8) {
	if compute > 0 {
		curCompute := player.GetComputeDemand()
		player.SetComputeDemand(utils.SafeSubUint8(curCompute, compute))
	}
}

func addArmory(player *datamod.PlayersRow) {
	curArmories := player.GetCurArmories()
	player.SetCurArmories(utils.SafeAddUint8(curArmories, 1))
}

func subArmory(player *datamod.PlayersRow) {
	curArmories := player.GetCurArmories()
	player.SetCurArmories(utils.SafeSubUint8(curArmories, 1))
}
