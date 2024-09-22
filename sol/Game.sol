// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./solgen/IActions.sol";
import "./solgen/ITables.sol";
import "./solgen/ICore.sol";
import {Arch, NonZeroBoolean_True} from "./solgen/Arch.sol";
import {BoardLib} from "./solgen/BoardLib.sol";

import {UnitPrototypeAdder, UnitType} from "./Units.sol";
import {BuildingPrototypeAdder, BuildingType} from "./Buildings.sol";
import {LibCommand, WorkerCommandType, FighterCommandType} from "./LibCommand.sol";

uint8 constant TopLane_Y = 3;
uint8 constant BottomLane_Y = 4;

uint8 constant MainBuildingId = 1;
uint8 constant TopLaneBuildingId = 2;
uint8 constant BottomLaneBuildingId = 3;

enum UnitState {
    Nil,
    Unpaid,
    Spawning,
    Active,
    Inactive,
    Dead
}

contract Game is Arch {
    address[2] internal players;

    function getPlayerAddress(uint8 playerId) public view returns (address) {
        return players[playerId - 1];
    }

    function getPlayerId(address playerAddress) public view returns (uint8) {
        for (uint8 i = 0; i < players.length; i++) {
            if (players[i] == playerAddress) {
                return i + 1;
            }
        }
        return 0;
    }

    function addPlayers(bytes memory data) internal {
        address[] memory _players = abi.decode(data, (address[]));
        if (players.length != 2) {
            revert("Game: must have exactly 2 players");
        }
        for (uint8 i = 0; i < _players.length; i++) {
            BoardLib.initPlayer(ICore(proxy), i + 1, 3);
            players[i] = _players[i];
        }
    }

    function otherPlayer(uint8 playerId) internal pure returns (uint8) {
        return (playerId % 2) + 1;
    }

    function _initialize(bytes memory data) internal override {
        UnitPrototypeAdder.addUnitPrototypes(ICore(proxy));
        BuildingPrototypeAdder.addBuildingPrototypes(ICore(proxy));
        BoardLib.initialize(ICore(proxy));
        addPlayers(data);
        ICore(proxy).start();
    }

    function start() public virtual override {
        ICore(proxy).start();
    }

    function createUnit(
        ActionData_CreateUnit memory action
    ) public virtual override {
        if (action.unitType == uint8(UnitType.Worker)) {
            revert("Game: only fighters can be created");
        }
        ICore(proxy).createUnit(action);

        ActionData_AssignUnit memory assignUnitData;
        assignUnitData.playerId = action.playerId;
        assignUnitData.unitId = ITables(proxy)
            .getPlayersRow(assignUnitData.playerId)
            .unitCount;

        uint8 targetPlayerId = otherPlayer(action.playerId);
        uint64 command;

        if (action.y >= TopLane_Y && action.y <= BottomLane_Y) {
            command = LibCommand.assignFighterToAttackBuilding(
                targetPlayerId,
                MainBuildingId
            );
        } else {
            uint8 targetUnitId;
            if (action.y < TopLane_Y) {
                targetUnitId = TopLane_Y;
            } else {
                targetUnitId = BottomLane_Y;
            }
            uint8 targetUnitState = ITables(proxy)
                .getUnitsRow(targetPlayerId, targetUnitId)
                .state;
            if (targetUnitState == uint8(UnitState.Dead)) {
                command = LibCommand.assignFighterToAttackBuilding(
                    targetPlayerId,
                    MainBuildingId
                );
            } else {
                command = LibCommand.assignFighterToAttackUnit(
                    targetPlayerId,
                    targetUnitId
                );
            }
        }

        assignUnitData.command = command;

        ICore(proxy).assignUnit(assignUnitData);
    }

    function archTick() public {
        super.tick();
    }

    function tick() public override {
        (bool success, ) = address(this).call{gas: gasleft() - 10_000}(
            abi.encodeWithSignature("archTick()")
        );
        require(success);
        if (needsPurge == NonZeroBoolean_True) {
            return;
        }
        for (uint8 playerId = 1; playerId <= 2; playerId++) {
            if (gasleft() < 10_000) {
                return;
            }
            uint8 targetPlayerId = otherPlayer(playerId);
            uint8 targetMainBuildingIntegrity = ITables(proxy)
                .getBuildingsRow(targetPlayerId, MainBuildingId)
                .integrity;
            if (targetMainBuildingIntegrity == 0) {
                continue;
            }
            uint8 unitCount = ITables(proxy).getPlayersRow(playerId).unitCount;
            for (uint8 unitId = 4; unitId <= unitCount; unitId++) {
                if (gasleft() < 10_000) {
                    return;
                }
                RowData_Units memory unit = ITables(proxy).getUnitsRow(
                    playerId,
                    unitId
                );
                if (unit.state != uint8(UnitState.Active)) {
                    // Unit is not active
                    continue;
                }
                (FighterCommandType unitCommandType, , ) = LibCommand
                    .parseFighterCommand(unit.command);
                if (unitCommandType == FighterCommandType.HoldPosition) {
                    ActionData_AssignUnit memory assignUnitData;
                    assignUnitData.playerId = playerId;
                    assignUnitData.unitId = unitId;
                    assignUnitData.command = LibCommand
                        .assignFighterToAttackBuilding(
                            targetPlayerId,
                            MainBuildingId
                        );

                    ICore(proxy).assignUnit(assignUnitData);
                }
            }
        }
    }
}
