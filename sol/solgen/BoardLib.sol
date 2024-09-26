// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./IActions.sol";
import "./ICore.sol";

uint8 constant WIDTH = 15;
uint8 constant HEIGHT = 8;

library BoardLib {
    function initCore(ICore proxy, uint16 w, uint16 h) internal {
        ActionData_Initialize memory initializeData;
        initializeData.width = w;
        initializeData.height = h;
        proxy.initialize(initializeData);
    }

    function addBuilding(ICore proxy, uint8 playerId, uint8 buildingTypeId, uint16 x, uint16 y) internal {
        ActionData_PlaceBuilding memory placeBuildingData;
        placeBuildingData.playerId = playerId;
        placeBuildingData.buildingType = buildingTypeId;
        placeBuildingData.x = x;
        placeBuildingData.y = y;
        proxy.placeBuilding(placeBuildingData);
    }

    function addUnit(ICore proxy, uint8 playerId, uint8 unitTypeId, uint16 x, uint16 y) internal {
        ActionData_CreateUnit memory createUnitData;
        createUnitData.playerId = playerId;
        createUnitData.unitType = unitTypeId;
        createUnitData.x = x;
        createUnitData.y = y;
        proxy.createUnit(createUnitData);
    }
    
    function assignUnit(ICore proxy, uint8 playerId, uint8 unitId, uint64 command) internal {
        ActionData_AssignUnit memory assignUnitData;
        assignUnitData.playerId = playerId;
        assignUnitData.unitId = unitId;
        assignUnitData.command = command;
        proxy.assignUnit(assignUnitData);
    }

    function initPlayer(ICore proxy, uint8 playerId, uint8 unpurgeableUnitCount) internal {
        ActionData_AddPlayer memory addPlayerData;
        addPlayerData.unpurgeableUnitCount = unpurgeableUnitCount;

        if (playerId == 1) {
            addPlayerData.spawnAreaX = 2;
            addPlayerData.spawnAreaY = 0;
            addPlayerData.spawnAreaWidth = 5;
            addPlayerData.spawnAreaHeight = 8;
            addPlayerData.buildAreaX = 0;
            addPlayerData.buildAreaY = 3;
            addPlayerData.buildAreaWidth = 2;
            addPlayerData.buildAreaHeight = 2;
            addPlayerData.workerPortX = 0;
            addPlayerData.workerPortY = 3;
            proxy.addPlayer(addPlayerData);
            addBuilding(proxy, 1, 1, 0, 3);
            addUnit(proxy, 1, 4, 2, 3);
            assignUnit(proxy, 1, 1, 65543);
            addUnit(proxy, 1, 5, 2, 1);
            assignUnit(proxy, 1, 2, 131073);
            addUnit(proxy, 1, 5, 2, 6);
            assignUnit(proxy, 1, 3, 131078);
        } else if (playerId == 2) {
            addPlayerData.spawnAreaX = 8;
            addPlayerData.spawnAreaY = 0;
            addPlayerData.spawnAreaWidth = 5;
            addPlayerData.spawnAreaHeight = 8;
            addPlayerData.buildAreaX = 13;
            addPlayerData.buildAreaY = 3;
            addPlayerData.buildAreaWidth = 2;
            addPlayerData.buildAreaHeight = 2;
            addPlayerData.workerPortX = 14;
            addPlayerData.workerPortY = 3;
            proxy.addPlayer(addPlayerData);
            addBuilding(proxy, 2, 1, 13, 3);
            addUnit(proxy, 2, 4, 12, 3);
            assignUnit(proxy, 2, 1, 65544);
            addUnit(proxy, 2, 5, 12, 1);
            assignUnit(proxy, 2, 2, 786433);
            addUnit(proxy, 2, 5, 12, 6);
            assignUnit(proxy, 2, 3, 786438);
        } else {
            revert();
        }
    }

    function initEnvironment(ICore proxy) internal {
        addBuilding(proxy, 0, 2, 7, 0);
        addBuilding(proxy, 0, 2, 7, 2);
        addBuilding(proxy, 0, 2, 7, 3);
        addBuilding(proxy, 0, 2, 7, 4);
        addBuilding(proxy, 0, 2, 7, 5);
        addBuilding(proxy, 0, 2, 7, 7);
        addBuilding(proxy, 0, 3, 0, 2);
        addBuilding(proxy, 0, 3, 14, 2);
    }

    function initialize(ICore proxy) internal {
        initCore(proxy, WIDTH, HEIGHT);
        initEnvironment(proxy);
    }
}
