// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

/* Autogenerated file. Do not edit manually. */

import "./IActions.sol";

abstract contract Entrypoint is IActions {
    function executeMultipleActions(
        uint32[] memory actionIds,
        uint8[] memory actionCount,
        bytes[] memory actionData
    ) external {
        uint256 actionIdx = 0;
        for (uint256 i = 0; i < actionIds.length; i++) {
            uint256 numActions = uint256(actionCount[i]);
            for (uint256 j = actionIdx; j < actionIdx + numActions; j++) {
                _executeAction(actionIds[i], actionData[j]);
            }
            actionIdx += numActions;
        }
    }

    function _executeAction(uint32 actionId, bytes memory actionData) private {
        if (actionId == 0x3eaf5d9f) {
            tick();
        } else if (actionId == 0xeaba9837) {
            ActionData_Initialize memory action = abi.decode(
                actionData,
                (ActionData_Initialize)
            );
            initialize(action);
        } else if (actionId == 0xbe9a6555) {
            start();
        } else if (actionId == 0x97b1de79) {
            ActionData_CreateUnit memory action = abi.decode(
                actionData,
                (ActionData_CreateUnit)
            );
            createUnit(action);
        } else if (actionId == 0xf8613b59) {
            ActionData_AssignUnit memory action = abi.decode(
                actionData,
                (ActionData_AssignUnit)
            );
            assignUnit(action);
        } else if (actionId == 0xd74de075) {
            ActionData_PlaceBuilding memory action = abi.decode(
                actionData,
                (ActionData_PlaceBuilding)
            );
            placeBuilding(action);
        } else if (actionId == 0xa76b2b54) {
            ActionData_AddPlayer memory action = abi.decode(
                actionData,
                (ActionData_AddPlayer)
            );
            addPlayer(action);
        } else if (actionId == 0x251eb4c0) {
            ActionData_AddUnitPrototype memory action = abi.decode(
                actionData,
                (ActionData_AddUnitPrototype)
            );
            addUnitPrototype(action);
        } else if (actionId == 0x982d778d) {
            ActionData_AddBuildingPrototype memory action = abi.decode(
                actionData,
                (ActionData_AddBuildingPrototype)
            );
            addBuildingPrototype(action);
        } else {
            revert("Entrypoint: Invalid action ID");
        }
    }

    function tick() public virtual;

    function initialize(ActionData_Initialize memory action) public virtual;

    function start() public virtual;

    function createUnit(ActionData_CreateUnit memory action) public virtual;

    function assignUnit(ActionData_AssignUnit memory action) public virtual;

    function placeBuilding(
        ActionData_PlaceBuilding memory action
    ) public virtual;

    function addPlayer(ActionData_AddPlayer memory action) public virtual;

    function addUnitPrototype(
        ActionData_AddUnitPrototype memory action
    ) public virtual;

    function addBuildingPrototype(
        ActionData_AddBuildingPrototype memory action
    ) public virtual;
}
