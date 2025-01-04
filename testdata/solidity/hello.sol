// SPDX-License-Identifier: MIT
// compiler version must be greater than or equal to 0.8.24 and less than 0.9.0
pragma solidity ^0.8.28;

contract HelloWorld {
    uint256 public storedValue;

    constructor(uint256 initialValue) {
        storedValue = initialValue;
    }

    function getValue() public view returns (uint256) {
        return storedValue;
    }

    function setValue(uint256 newValue) public {
        storedValue = newValue;
    }
}