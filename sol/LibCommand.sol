// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

uint8 constant NilPlayerId = 0;
uint8 constant NilBuildingId = 0;

enum WorkerCommandType {
    Idle,
    Gather,
    Build
}

enum FighterCommandType {
    HoldPosition,
    AttackBuilding,
    AttackUnit
}

library LibCommand {
    function newWorkerCommand(
        WorkerCommandType cmdType,
        uint8 playerId,
        uint8 buildingId
    ) internal pure returns (uint64) {
        uint64 cmd;
        cmd |= uint64(cmdType) << 16;
        cmd |= uint64(playerId) << 8;
        cmd |= uint64(buildingId);
        return cmd;
    }

    function assignWorkerToIdle() internal pure returns (uint64) {
        return
            newWorkerCommand(
                WorkerCommandType.Idle,
                NilPlayerId,
                NilBuildingId
            );
    }

    function assignWorkerToGather(
        uint8 buildingId
    ) internal pure returns (uint64) {
        return
            newWorkerCommand(WorkerCommandType.Gather, NilPlayerId, buildingId);
    }

    function assignWorkerToBuild(
        uint8 playerId,
        uint8 buildingId
    ) internal pure returns (uint64) {
        return newWorkerCommand(WorkerCommandType.Build, playerId, buildingId);
    }

    function parseWorkerCommand(
        uint64 cmd
    ) internal pure returns (WorkerCommandType, uint8, uint8) {
        WorkerCommandType cmdType = WorkerCommandType(uint8(cmd >> 16));
        uint8 playerId = uint8(cmd >> 8);
        uint8 buildingId = uint8(cmd);
        return (cmdType, playerId, buildingId);
    }

    function newFighterCommand(
        FighterCommandType cmdType,
        uint16 alpha,
        uint16 beta
    ) internal pure returns (uint64) {
        uint64 cmd;
        cmd |= uint64(cmdType) << 32;
        cmd |= uint64(alpha) << 16;
        cmd |= uint64(beta);
        return cmd;
    }

    function assignFighterToHoldPosition() internal pure returns (uint64) {
        return newFighterCommand(FighterCommandType.HoldPosition, 0, 0);
    }

    function assignFighterToAttackBuilding(
        uint8 playerId,
        uint8 buildingId
    ) internal pure returns (uint64) {
        return
            newFighterCommand(
                FighterCommandType.AttackBuilding,
                uint16(playerId),
                uint16(buildingId)
            );
    }

    function assignFighterToAttackUnit(
        uint8 playerId,
        uint8 unitId
    ) internal pure returns (uint64) {
        return
            newFighterCommand(
                FighterCommandType.AttackUnit,
                uint16(playerId),
                uint16(unitId)
            );
    }

    function parseFighterCommand(
        uint64 cmd
    ) internal pure returns (FighterCommandType, uint16, uint16) {
        FighterCommandType cmdType = FighterCommandType(uint8(cmd >> 32));
        uint16 alpha = uint16(cmd >> 16);
        uint16 beta = uint16(cmd);
        return (cmdType, alpha, beta);
    }
}
