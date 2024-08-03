// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import {console2} from "forge-std/Test.sol";
import {Game} from "./Game.sol";
import {TickMaster} from "./TickMaster.sol";
import {Clones} from "openzeppelin/proxy/Clones.sol";

uint256 constant tickAllocGasPerPlayer = 500_000;
uint256 constant tickAllocBlocks = 1800;

contract GameFactory is TickMaster {
    address public immutable gameImplementation;
    address public immutable coreImplementation;

    event GameCreated(
        address gameAddress,
        string lobbyId,
        address sender,
        address origin
    );

    constructor(
        uint256 _maxGasAllocation,
        address _gameImplementation,
        address _coreImplementation
    ) TickMaster(_maxGasAllocation) {
        gameImplementation = _gameImplementation;
        coreImplementation = _coreImplementation;
    }

    function createGame(
        string memory lobbyId,
        address[] memory _players
    ) external returns (address) {
        address gameAddress = Clones.clone(gameImplementation);
        bool start = true;
        Game(gameAddress).initialize(coreImplementation);
        _setGasAlloc(
            gameAddress,
            tickAllocGasPerPlayer * _players.length,
            block.number + tickAllocBlocks
        );
        emit GameCreated(gameAddress, lobbyId, msg.sender, tx.origin);
        return gameAddress;
    }
}
