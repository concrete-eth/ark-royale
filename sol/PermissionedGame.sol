// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./solgen/IActions.sol";
import {Game} from "./Game.sol";

contract PermissionedGame is Game {
    modifier onlyPlayer(uint8 playerId) {
        // require(msg.sender == players[playerId - 1], "Game: onlyPlayer");
        _;
    }

    modifier onlyPlayerOne() {
        // require(msg.sender == players[0], "Game: onlyPlayerOne");
        _;
    }

    function start() public override onlyPlayerOne {
        super.start();
    }

    function createUnit(
        ActionData_CreateUnit memory action
    ) public override onlyPlayer(action.playerId) {
        super.createUnit(action);
    }

    function assignUnit(
        ActionData_AssignUnit memory action
    ) public override onlyPlayer(action.playerId) {
        super.assignUnit(action);
    }

    function placeBuilding(
        ActionData_PlaceBuilding memory action
    ) public override onlyPlayer(action.playerId) {
        super.placeBuilding(action);
    }
}
