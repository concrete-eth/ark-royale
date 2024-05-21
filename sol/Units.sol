// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./solgen/ICore.sol";

enum UnitType {
    Nil,
    Worker,
    Air,
    AntiAir,
    Tank
}

enum LayerId {
    Land,
    Hover,
    Air,
    Count
}

library UnitPrototypeAdder {
    function addUnitPrototypes(ICore core) internal {
        // Worker
        core.addUnitPrototype(
            ActionData_AddUnitPrototype({
                layer: uint8(LayerId.Hover),
                resourceCost: 150,
                computeCost: 1,
                spawnTime: 1,
                maxIntegrity: 25,
                landStrength: 0,
                hoverStrength: 0,
                airStrength: 0,
                attackCooldown: 0,
                attackRange: 0,
                isAssault: false,
                isWorker: true
            })
        );

        // Air
        core.addUnitPrototype(
            ActionData_AddUnitPrototype({
                layer: uint8(LayerId.Air),
                resourceCost: 100,
                computeCost: 1,
                spawnTime: 2,
                maxIntegrity: 25,
                landStrength: 5,
                hoverStrength: 3,
                airStrength: 3,
                attackCooldown: 2,
                attackRange: 2,
                isAssault: true,
                isWorker: false
            })
        );

        // AntiAir
        core.addUnitPrototype(
            ActionData_AddUnitPrototype({
                layer: uint8(LayerId.Land),
                resourceCost: 150,
                computeCost: 1,
                spawnTime: 4,
                maxIntegrity: 100,
                landStrength: 5,
                hoverStrength: 15,
                airStrength: 15,
                attackCooldown: 3,
                attackRange: 4,
                isAssault: false,
                isWorker: false
            })
        );

        // Tank
        core.addUnitPrototype(
            ActionData_AddUnitPrototype({
                layer: uint8(LayerId.Land),
                resourceCost: 300,
                computeCost: 1,
                spawnTime: 8,
                maxIntegrity: 150,
                landStrength: 15,
                hoverStrength: 3,
                airStrength: 5,
                attackCooldown: 6,
                attackRange: 3,
                isAssault: false,
                isWorker: false
            })
        );
    }
}
