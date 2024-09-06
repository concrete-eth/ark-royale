// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./solgen/ICore.sol";

enum BuildingType {
    Nil,
    Main,
    Pit,
    Mine
}

library BuildingPrototypeAdder {
    function addBuildingPrototypes(ICore core) internal {
        // Main
        core.addBuildingPrototype(
            ActionData_AddBuildingPrototype({
                width: 2,
                height: 2,
                resourceCost: 0,
                resourceCapacity: 300,
                computeCapacity: 8,
                resourceMine: 0,
                mineTime: 0,
                maxIntegrity: 250,
                buildingTime: 0,
                isArmory: false,
                isEnvironment: false
            })
        );

        // Pit
        core.addBuildingPrototype(
            ActionData_AddBuildingPrototype({
                width: 1,
                height: 1,
                resourceCost: 0,
                resourceCapacity: 0,
                computeCapacity: 0,
                resourceMine: 0,
                mineTime: 0,
                maxIntegrity: 0,
                buildingTime: 0,
                isArmory: false,
                isEnvironment: true
            })
        );

        // Mine
        core.addBuildingPrototype(
            ActionData_AddBuildingPrototype({
                width: 1,
                height: 1,
                resourceCost: 0,
                resourceCapacity: 0,
                computeCapacity: 0,
                resourceMine: 25,
                mineTime: 0,
                maxIntegrity: 0,
                buildingTime: 0,
                isArmory: false,
                isEnvironment: true
            })
        );
    }
}
