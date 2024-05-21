// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./solgen/ICore.sol";

enum BuildingType {
    Nil,
    Main,
    Storage,
    Lab,
    Armory,
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
                resourceCapacity: 100,
                computeCapacity: 2,
                resourceMine: 0,
                mineTime: 0,
                maxIntegrity: 200,
                buildingTime: 0,
                isArmory: false,
                isEnvironment: false
            })
        );
        // Storage
        core.addBuildingPrototype(
            ActionData_AddBuildingPrototype({
                width: 2,
                height: 2,
                resourceCost: 100,
                resourceCapacity: 50,
                computeCapacity: 0,
                resourceMine: 0,
                mineTime: 0,
                maxIntegrity: 150,
                buildingTime: 5,
                isArmory: false,
                isEnvironment: false
            })
        );
        // Lab
        core.addBuildingPrototype(
            ActionData_AddBuildingPrototype({
                width: 2,
                height: 2,
                resourceCost: 150,
                resourceCapacity: 0,
                computeCapacity: 3,
                resourceMine: 0,
                mineTime: 0,
                maxIntegrity: 75,
                buildingTime: 5,
                isArmory: false,
                isEnvironment: false
            })
        );
        // Armory
        core.addBuildingPrototype(
            ActionData_AddBuildingPrototype({
                width: 2,
                height: 2,
                resourceCost: 200,
                resourceCapacity: 0,
                computeCapacity: 0,
                resourceMine: 0,
                mineTime: 0,
                maxIntegrity: 100,
                buildingTime: 5,
                isArmory: true,
                isEnvironment: false
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
                resourceMine: 20,
                mineTime: 2,
                maxIntegrity: 0,
                buildingTime: 0,
                isArmory: false,
                isEnvironment: true
            })
        );
    }
}
