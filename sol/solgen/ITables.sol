// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

/* Autogenerated file. Do not edit manually. */

struct RowData_Meta {
    uint16 boardWidth;
    uint16 boardHeight;
    uint8 playerCount;
    uint8 unitPrototypeCount;
    uint8 buildingPrototypeCount;
    bool isInitialized;
    bool hasStarted;
    uint32 creationBlockNumber;
}

struct RowData_Players {
    uint16 spawnAreaX;
    uint16 spawnAreaY;
    uint16 workerPortX;
    uint16 workerPortY;
    uint16 curResource;
    uint16 maxResource;
    uint8 curArmories;
    uint8 computeSupply;
    uint8 computeDemand;
    uint8 unitCount;
    uint8 buildingCount;
    uint8 buildingPayQueuePointer;
    uint8 buildingBuildQueuePointer;
    uint8 unitPayQueuePointer;
}

struct RowData_Board {
    uint8 landObjectType;
    uint8 landPlayerId;
    uint8 landObjectId;
    uint8 hoverPlayerId;
    uint8 hoverUnitId;
    uint8 airPlayerId;
    uint8 airUnitId;
}

struct RowData_Units {
    uint16 x;
    uint16 y;
    uint8 unitType;
    uint8 state;
    uint8 load;
    uint8 integrity;
    uint32 timestamp;
    uint64 command;
    uint64 commandExtra;
    uint8 commandMeta;
    bool isPreTicked;
}

struct RowData_Buildings {
    uint16 x;
    uint16 y;
    uint8 buildingType;
    uint8 state;
    uint8 integrity;
    uint32 timestamp;
}

struct RowData_UnitPrototypes {
    uint8 layer;
    uint16 resourceCost;
    uint8 computeCost;
    uint8 spawnTime;
    uint8 maxIntegrity;
    uint8 landStrength;
    uint8 hoverStrength;
    uint8 airStrength;
    uint8 attackRange;
    uint8 attackCooldown;
    bool isAssault;
    bool isWorker;
}

struct RowData_BuildingPrototypes {
    uint8 width;
    uint8 height;
    uint16 resourceCost;
    uint16 resourceCapacity;
    uint8 computeCapacity;
    uint8 resourceMine;
    uint8 mineTime;
    uint8 maxIntegrity;
    uint8 buildingTime;
    bool isArmory;
    bool isEnvironment;
}

interface ITables {
    function getMetaRow() external view returns (RowData_Meta memory);
    function getPlayersRow(
        uint8 playerId
    ) external view returns (RowData_Players memory);
    function getBoardRow(
        uint16 x,
        uint16 y
    ) external view returns (RowData_Board memory);
    function getUnitsRow(
        uint8 playerId,
        uint8 unitId
    ) external view returns (RowData_Units memory);
    function getBuildingsRow(
        uint8 playerId,
        uint8 buildingId
    ) external view returns (RowData_Buildings memory);
    function getUnitPrototypesRow(
        uint8 unitType
    ) external view returns (RowData_UnitPrototypes memory);
    function getBuildingPrototypesRow(
        uint8 buildingType
    ) external view returns (RowData_BuildingPrototypes memory);
}
