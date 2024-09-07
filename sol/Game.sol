// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./solgen/IActions.sol";
import "./solgen/ITables.sol";
import "./solgen/ICore.sol";
import {Arch, NonZeroBoolean_True} from "./solgen/Arch.sol";

import {UnitPrototypeAdder, UnitType} from "./Units.sol";
import {BuildingPrototypeAdder, BuildingType} from "./Buildings.sol";

uint8 constant WIDTH = 15;
uint8 constant HEIGHT = 8;
uint8 constant BUILD_RADIUS = 0;

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

    function _initialize(bytes memory data) internal override {
        ActionData_Initialize memory initializeData;
        initializeData.width = WIDTH;
        initializeData.height = HEIGHT;

        ICore(proxy).initialize(initializeData);

        UnitPrototypeAdder.addUnitPrototypes(ICore(proxy));
        BuildingPrototypeAdder.addBuildingPrototypes(ICore(proxy));

        address[] memory _players = abi.decode(data, (address[]));
        _addPlayers(_players);
        _addPit();

        ICore(proxy).start();
    }

    function _addPlayers(address[] memory _players) internal {
        if (players.length != 2) {
            revert("Game: must have exactly 2 players");
        }
        for (uint8 i = 0; i < _players.length; i++) {
            _addPlayer();
            players[i] = _players[i];
        }
    }

    function _addPlayer() internal {
        uint8 playerId = ICore(proxy).getMetaRow().playerCount + 1;

        (uint8 x, uint8 y) = _getMainBuildingPosition(playerId);

        uint16 spawnAreaX;
        if (playerId == 1) {
            spawnAreaX = 2;
        } else {
            spawnAreaX = (WIDTH / 2) + 1;
        }

        ICore(proxy).addPlayer(
            ActionData_AddPlayer({
                spawnAreaX: spawnAreaX,
                spawnAreaY: 0,
                spawnAreaWidth: (WIDTH / 2) - 2,
                spawnAreaHeight: HEIGHT,
                workerPortX: x,
                workerPortY: y
            })
        );

        ICore(proxy).placeBuilding(
            ActionData_PlaceBuilding({
                playerId: playerId,
                buildingType: uint8(BuildingType.Main),
                x: x,
                y: y
            })
        );

        if (x == 0) {
            x = 2;
        } else {
            x = WIDTH - 3;
        }
        ICore(proxy).createUnit(
            ActionData_CreateUnit({
                playerId: playerId,
                unitType: uint8(UnitType.Turret),
                x: x,
                y: 1
            })
        );
        ICore(proxy).createUnit(
            ActionData_CreateUnit({
                playerId: playerId,
                unitType: uint8(UnitType.Turret),
                x: x,
                y: HEIGHT - 2
            })
        );

        uint16 mineX;
        uint16 mineY = y - 1;
        uint16 workerX;
        uint16 workerY = mineY;
        if (playerId == 1) {
            mineX = 0;
            workerX = 2;
        } else {
            mineX = WIDTH - 1;
            workerX = WIDTH - 3;
        }
        ICore(proxy).placeBuilding(
            ActionData_PlaceBuilding({
                playerId: 0,
                buildingType: uint8(BuildingType.Mine),
                x: mineX,
                y: mineY
            })
        );
        ICore(proxy).createUnit(
            ActionData_CreateUnit({
                playerId: playerId,
                unitType: uint8(UnitType.Worker),
                x: workerX,
                y: workerY
            })
        );

        ActionData_AssignUnit memory assignUnitData;
        assignUnitData.playerId = playerId;
        assignUnitData.unitId = 3;

        uint64 command = 0;
        command |= uint64(1) << 16; // Gather
        command |= uint64(playerId); // Target building

        assignUnitData.command = command;

        ICore(proxy).assignUnit(assignUnitData);
    }

    function _addPit() internal {
        uint16 x = WIDTH / 2;
        for (uint8 y = 0; y < HEIGHT; y++) {
            if (y == 1 || y == HEIGHT - 2) {
                continue;
            }
            ICore(proxy).placeBuilding(
                ActionData_PlaceBuilding({
                    playerId: 0,
                    buildingType: uint8(BuildingType.Pit),
                    x: x,
                    y: y
                })
            );
        }
    }

    function _getMainBuildingPosition(
        uint8 playerId
    ) internal pure returns (uint8, uint8) {
        if (playerId == 1) {
            return (0, 3);
        } else {
            return (WIDTH - 2, 3);
        }
    }

    function start() public virtual override {
        ICore(proxy).start();
    }

    function createUnit(
        ActionData_CreateUnit memory action
    ) public virtual override {
        // if (action.unitType == uint8(UnitType.Worker)) {
        //     revert("Game: only fighters can be created");
        // }
        ICore(proxy).createUnit(action);

        ActionData_AssignUnit memory assignUnitData;

        assignUnitData.playerId = action.playerId;
        assignUnitData.unitId = ITables(proxy)
            .getPlayersRow(assignUnitData.playerId)
            .unitCount;

        uint64 command = 0;

        uint8 targetPlayerId = (action.playerId % 2) + 1;
        command |= uint64(targetPlayerId) << 16; // Target player id

        if (action.y >= 3 && action.y <= 4) {
            command |= uint64(1) << 32; // Target building
            command |= uint64(1); // Target main building
        } else {
            uint8 targetUnitId;
            if (action.y < 3) {
                targetUnitId = 1;
            } else {
                targetUnitId = 2;
            }
            uint8 targetUnitState = ITables(proxy)
                .getUnitsRow(targetPlayerId, targetUnitId)
                .state;
            if (targetUnitState == 5) {
                // Unit is dead
                command |= uint64(1) << 32; // Target building
                command |= uint64(1); // Target main building
            } else {
                command |= uint64(2) << 32; // Target unit
                command |= uint64(targetUnitId); // Target left unit
            }
        }

        assignUnitData.command = command;

        ICore(proxy).assignUnit(assignUnitData);
    }

    function archTick() public {
        super.tick();
    }

    function tick() public override {
        (bool success, ) = address(this).call{gas: gasleft() - 10000}(
            abi.encodeWithSignature("archTick()")
        );
        require(success);
        if (needsPurge == NonZeroBoolean_True) {
            return;
        }
        for (uint8 playerId = 1; playerId <= 2; playerId++) {
            if (gasleft() < 10000) {
                return;
            }
            uint8 targetPlayerId = (playerId % 2) + 1;
            uint8 targetMainBuildingIntegrity = ITables(proxy)
                .getBuildingsRow(targetPlayerId, 1)
                .integrity;
            if (targetMainBuildingIntegrity == 0) {
                continue;
            }
            uint8 unitCount = ITables(proxy).getPlayersRow(playerId).unitCount;
            for (uint8 unitId = 4; unitId <= unitCount; unitId++) {
                if (gasleft() < 10000) {
                    return;
                }
                RowData_Units memory unit = ITables(proxy).getUnitsRow(
                    playerId,
                    unitId
                );
                if (unit.state != 3) {
                    // Unit is not active
                    continue;
                }
                uint64 unitCommandType = unit.command >> 32;
                if (unitCommandType == 0) {
                    uint64 command = 0;
                    command |= uint64(1) << 32; // Target building
                    command |= uint64(targetPlayerId) << 16; // Target player id
                    command |= uint64(1); // Target main building
                    ActionData_AssignUnit memory assignUnitData;
                    assignUnitData.playerId = playerId;
                    assignUnitData.unitId = unitId;
                    assignUnitData.command = command;
                    ICore(proxy).assignUnit(assignUnitData);
                }
            }
        }
    }
}
