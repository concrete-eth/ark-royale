// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import {console2} from "forge-std/Test.sol";
import "solmate/auth/Owned.sol";

interface ITick {
    function tick() external;
    function setTarget(address target) external;
}

struct Tickee {
    uint64 gas;
    uint64 index;
    uint64 evictionBlock;
}

uint256 constant tickGasMargin = 75_000;
uint256 constant returnGasMargin = 5_000;

contract TickMaster is Owned {
    mapping(address => Tickee) internal tickees;
    mapping(uint256 => address) internal tickAddresses;

    uint256 public nActiveTickees;
    uint256 public totalGasAllocation;
    uint256 public maxGasAllocation;

    uint256 public lastBlock;

    event GasAllocSet(
        address indexed tickee,
        uint256 gas,
        uint256 evictionBlock
    );

    constructor(uint256 _maxGasAllocation) Owned(msg.sender) {
        maxGasAllocation = _maxGasAllocation;
    }

    function tick() external {
        require(block.number > lastBlock, "TickMaster: only once per block");
        lastBlock = block.number;

        // Evict tickees if due
        uint256 n = nActiveTickees;
        for (uint256 i = 0; i < n; i++) {
            console2.log("Tickee:", i);
            if (gasleft() < tickGasMargin) {
                // This will only happen if we are very tight on gas. The call should still succeed so
                // that any evictions that can be done are done.
                console2.log("Under gas margin");
                break;
            }
            address t = tickAddresses[i];
            uint256 evictionBlock = tickees[t].evictionBlock;
            if (block.number >= evictionBlock) {
                _setGasAlloc(t, 0, evictionBlock);
            }
        }

        // Tick active tickees
        for (uint256 i = 0; i < nActiveTickees; i++) {
            console2.log("Gasleft:", gasleft());
            console2.log("Tickee:", i);
            if (gasleft() < tickGasMargin) {
                console2.log("Under gas margin");
                break;
            }
            address t = tickAddresses[i];
            uint256 gas = tickees[t].gas;
            if (gasleft() >= gas + returnGasMargin) {
                (bool success, ) = t.call{gas: gas}(
                    abi.encodeWithSignature("tick()")
                );
            }
        }
    }

    function _setGasAlloc(
        address t,
        uint256 gas,
        uint256 evictionBlock
    ) internal {
        console2.log("Setting gas alloc:", t, gas, evictionBlock);

        Tickee storage tickee = tickees[t];

        require(
            totalGasAllocation + gas - tickee.gas < maxGasAllocation,
            "TickMaster: gas allocation exceeds max"
        );

        totalGasAllocation -= tickee.gas;
        totalGasAllocation += gas;

        if (gas == 0) {
            if (tickee.gas == 0) {
                // Tickee is already inactive
                return;
            } else {
                // Tickee is active
                // Swap with last active tickee and pop
                nActiveTickees--;
                uint256 lastIndex = nActiveTickees;
                address lastAddress = tickAddresses[lastIndex];
                tickAddresses[tickee.index] = lastAddress;
                tickees[lastAddress].index = tickee.index;
                // Remove allocation
                tickee.gas = 0;
                tickee.index = 0;
            }
        } else {
            if (tickee.gas == 0) {
                // Tickee is not active
                // Add to end of active tickees
                uint256 lastIndex = nActiveTickees;
                nActiveTickees++;
                tickee.index = uint64(lastIndex);
                tickAddresses[lastIndex] = t;
            }
            // Set gas
            tickee.gas = uint64(gas);
            tickee.evictionBlock = uint64(evictionBlock);
        }

        emit GasAllocSet(t, gas, evictionBlock);
    }

    function setGasAlloc(
        address t,
        uint256 gas,
        uint256 evictionBlock
    ) external onlyOwner {
        _setGasAlloc(t, gas, evictionBlock);
    }

    function getGasAllocOf(address t) public view returns (uint256) {
        return tickees[t].gas;
    }

    function getIndexOf(address t) public view returns (uint256) {
        return tickees[t].index;
    }

    function getAddressOf(uint256 idx) public view returns (address) {
        return tickAddresses[idx];
    }
}
