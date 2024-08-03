// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./solgen/IActions.sol";
import "./solgen/ICore.sol";
import {Arch} from "./solgen/Arch.sol";

import {UnitPrototypeAdder, UnitType} from "./Units.sol";
import {BuildingPrototypeAdder, BuildingType} from "./Buildings.sol";

uint8 constant SIZE = 26;
uint8 constant BUILD_RADIUS = 4;

contract Game is Arch {
    address[4] internal players;

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

    function _initialize() internal override {
        ActionData_Initialize memory initializeData;
        initializeData.width = SIZE;
        initializeData.height = SIZE;

        ICore(proxy).initialize(initializeData);

        UnitPrototypeAdder.addUnitPrototypes(ICore(proxy));
        BuildingPrototypeAdder.addBuildingPrototypes(ICore(proxy));

        address[] memory _players = new address[](2);
        _players[0] = tx.origin;
        _players[1] = _players[0];
        _addPlayers(_players);

        start();
    }

    function _addPlayers(address[] memory _players) internal {
        if (players.length > 4) {
            revert("Game: too many players");
        }
        for (uint8 i = 0; i < _players.length; i++) {
            _addPlayer();
            players[i] = _players[i];
        }
    }

    function _addPlayer() internal {
        uint8 playerId = ICore(proxy).getMetaRow().playerCount + 1;

        uint8 x;
        uint8 y;

        ActionData_AddPlayer memory addPlayerData;

        (x, y) = _getSpawnAreaPosition(playerId);
        addPlayerData.spawnAreaX = x;
        addPlayerData.spawnAreaY = y;

        (x, y) = _getWorkerPortPosition(playerId);
        addPlayerData.workerPortX = x;
        addPlayerData.workerPortY = y;

        ICore(proxy).addPlayer(addPlayerData);

        (x, y) = _getMainBuildingPosition(playerId);
        ICore(proxy).placeBuilding(
            ActionData_PlaceBuilding({
                playerId: playerId,
                buildingType: uint8(BuildingType.Main),
                x: x,
                y: y
            })
        );

        (x, y) = _getFirstMinePosition(playerId);
        ICore(proxy).placeBuilding(
            ActionData_PlaceBuilding({
                playerId: 0,
                buildingType: uint8(BuildingType.Mine),
                x: x,
                y: y
            })
        );

        (x, y) = _getSecondMinePosition(playerId);
        ICore(proxy).placeBuilding(
            ActionData_PlaceBuilding({
                playerId: 0,
                buildingType: uint8(BuildingType.Mine),
                x: x,
                y: y
            })
        );

        ICore(proxy).createUnit(
            ActionData_CreateUnit({
                playerId: playerId,
                unitType: uint8(UnitType.Worker)
            })
        );
    }

    function _getPlayerSetupPosition(
        uint8 playerId,
        uint8 _start,
        uint8 _end
    ) internal pure returns (uint8, uint8) {
        if (playerId == 1) {
            return (_start, _start);
        } else if (playerId == 2) {
            return (_end, _end);
        } else if (playerId == 3) {
            return (_start, _end);
        } else if (playerId == 4) {
            return (_end, _start);
        } else {
            revert("Game: invalid player id");
        }
    }

    function _getMainBuildingPosition(
        uint8 playerId
    ) internal pure returns (uint8, uint8) {
        uint8 _start = BUILD_RADIUS;
        uint8 _end = SIZE - BUILD_RADIUS - 2;
        return _getPlayerSetupPosition(playerId, _start, _end);
    }

    function _getSpawnAreaPosition(
        uint8 playerId
    ) internal pure returns (uint8, uint8) {
        uint8 _start = BUILD_RADIUS + 2;
        uint8 _end = SIZE - BUILD_RADIUS - 2 - 2;
        return _getPlayerSetupPosition(playerId, _start, _end);
    }

    function _getWorkerPortPosition(
        uint8 playerId
    ) internal pure returns (uint8, uint8) {
        uint8 _start = BUILD_RADIUS;
        uint8 _end = SIZE - BUILD_RADIUS - 2 + 1;
        return _getPlayerSetupPosition(playerId, _start, _end);
    }

    function _getFirstMinePosition(
        uint8 playerId
    ) internal pure returns (uint8, uint8) {
        uint8 _start = 2;
        uint8 _end = SIZE - 1 - 2;
        return _getPlayerSetupPosition(playerId, _start, _end);
    }

    function _getSecondMinePosition(
        uint8 playerId
    ) internal pure returns (uint8, uint8) {
        uint8 _start = BUILD_RADIUS + 2 + 2 + 1;
        uint8 _end = SIZE - 1 - BUILD_RADIUS - 2 - 2 - 1;
        return _getPlayerSetupPosition(playerId, _start, _end);
    }

    function _isWithinPlayerBuildArea(
        uint8 playerId,
        uint16 x,
        uint16 y,
        uint16 w,
        uint16 h
    ) internal view returns (bool) {
        RowData_Buildings memory mainBuilding = ICore(proxy).getBuildingsRow(
            playerId,
            1
        );
        uint16 mainX = mainBuilding.x;
        uint16 mainY = mainBuilding.y;
        uint16 minX = mainX - BUILD_RADIUS;
        uint16 minY = mainY - BUILD_RADIUS;
        uint16 maxX = mainX + 2 + BUILD_RADIUS;
        uint16 maxY = mainY + 2 + BUILD_RADIUS;
        return x >= minX && x + w <= maxX && y >= minY && y + h <= maxY;
    }

    function placeBuilding(
        ActionData_PlaceBuilding memory action
    ) public virtual override {
        require(
            action.buildingType != uint8(BuildingType.Main),
            "Game: main building cannot be placed"
        );
        require(
            _isWithinPlayerBuildArea(action.playerId, action.x, action.y, 2, 2),
            "Game: building outside of build area"
        );
        ICore(proxy).placeBuilding(action);
    }

    function createUnit(
        ActionData_CreateUnit memory action
    ) public virtual override {
        if (action.unitType != uint8(UnitType.Worker)) {
            RowData_Players memory player = ICore(proxy).getPlayersRow(
                action.playerId
            );
            require(player.curArmories > 0, "Game: no armories");
        }
        ICore(proxy).createUnit(action);
    }

    function assignUnit(
        ActionData_AssignUnit memory action
    ) public virtual override {
        ICore(proxy).assignUnit(action);
    }

    function start() public virtual override {
        ICore(proxy).start();
    }

    function initialize(ActionData_Initialize memory action) public override {
        revert("not allowed");
    }

    function addPlayer(ActionData_AddPlayer memory action) public override {
        revert("not allowed");
    }

    function addUnitPrototype(
        ActionData_AddUnitPrototype memory action
    ) public override {
        revert("not allowed");
    }

    function addBuildingPrototype(
        ActionData_AddBuildingPrototype memory action
    ) public override {
        revert("not allowed");
    }
}
