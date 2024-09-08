// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./solgen/ICore.sol";

enum UnitType {
    Nil,
    Air,
    AntiAir,
    Tank,
    Turret,
    Worker
}

enum LayerId {
    Land,
    Hover,
    Air,
    Count
}

library UnitPrototypeAdder {
    function addUnitPrototypes(ICore core) internal {
        // Air
        core.addUnitPrototype(
            ActionData_AddUnitPrototype({
                layer: uint8(LayerId.Air),
                resourceCost: 100,
                computeCost: 1,
                spawnTime: 2,
                maxIntegrity: 25,
                landStrength: 5,
                hoverStrength: 0,
                airStrength: 3,
                attackCooldown: 2,
                attackRange: 2,
                isAssault: true,
                isConfrontational: true,
                isWorker: false,
                isPurgeable: true
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
                hoverStrength: 10,
                airStrength: 15,
                attackCooldown: 3,
                attackRange: 4,
                isAssault: false,
                isConfrontational: true,
                isWorker: false,
                isPurgeable: true
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
                landStrength: 10,
                hoverStrength: 0,
                airStrength: 3,
                attackCooldown: 4,
                attackRange: 3,
                isAssault: false,
                isConfrontational: false,
                isWorker: false,
                isPurgeable: true
            })
        );

        // Turret
        core.addUnitPrototype(
            ActionData_AddUnitPrototype({
                layer: uint8(LayerId.Land),
                resourceCost: 300,
                computeCost: 0,
                spawnTime: 8,
                maxIntegrity: 150,
                landStrength: 3,
                hoverStrength: 0,
                airStrength: 3,
                attackCooldown: 1,
                attackRange: 3,
                isAssault: false,
                isConfrontational: true,
                isWorker: false,
                isPurgeable: true
            })
        );

        // Worker
        core.addUnitPrototype(
            ActionData_AddUnitPrototype({
                layer: uint8(LayerId.Hover),
                resourceCost: 0,
                computeCost: 0,
                spawnTime: 0,
                maxIntegrity: 1,
                landStrength: 0,
                hoverStrength: 0,
                airStrength: 0,
                attackCooldown: 0,
                attackRange: 0,
                isAssault: false,
                isConfrontational: true,
                isWorker: true,
                isPurgeable: true
            })
        );
    }
}
