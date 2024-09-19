// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

// import "./solgen/ICore.sol";
// import {Arch, NonZeroBoolean_True} from "./solgen/Arch.sol";

// import {UnitPrototypeAdder, UnitType} from "./Units.sol";
// import {BuildingPrototypeAdder, BuildingType} from "./Buildings.sol";

// contract GenGame is Arch {
//     address[2] internal players;

//     function getPlayerAddress(uint8 playerId) public view returns (address) {
//         return players[playerId - 1];
//     }

//     function getPlayerId(address playerAddress) public view returns (uint8) {
//         for (uint8 i = 0; i < players.length; i++) {
//             if (players[i] == playerAddress) {
//                 return i + 1;
//             }
//         }
//         return 0;
//     }

//     function _initialize(bytes memory data) internal override {
//         ActionData_Initialize memory initializeData;
//         initializeData.width = WIDTH;
//         initializeData.height = HEIGHT;

//         ICore(proxy).initialize(initializeData);

//         UnitPrototypeAdder.addUnitPrototypes(ICore(proxy));
//         BuildingPrototypeAdder.addBuildingPrototypes(ICore(proxy));

//         address[] memory _players = abi.decode(data, (address[]));
//         _addPlayers(_players);
//         _addPit();

//         ICore(proxy).start();
//     }

//     function _addPlayers(address[] memory _players) internal {
//         if (players.length != 2) {
//             revert("Game: must have exactly 2 players");
//         }
//         for (uint8 i = 0; i < _players.length; i++) {
//             _addPlayer();
//             players[i] = _players[i];
//         }
//     }

//     function _addPlayer() internal {
//         uint8 playerId = ICore(proxy).getMetaRow().playerCount + 1;

//         (uint8 x, uint8 y) = _getMainBuildingPosition(playerId);

//         uint16 spawnAreaX;
//         if (playerId == 1) {
//             spawnAreaX = 2;
//         } else {
//             spawnAreaX = (WIDTH / 2) + 1;
//         }

//         ICore(proxy).addPlayer(
//             ActionData_AddPlayer({
//                 spawnAreaX: spawnAreaX,
//                 spawnAreaY: 0,
//                 spawnAreaWidth: (WIDTH / 2) - 2,
//                 spawnAreaHeight: HEIGHT,
//                 workerPortX: x,
//                 workerPortY: y,
//                 unpurgeableUnitCount: 3
//             })
//         );
//     }
// }
