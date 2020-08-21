// SPDX-License-Identifier: MIT
// wETHC token
// This is an ERC-20 token following OpenZepplin guidelines
// This token also wraps around some ether
pragma solidity ^0.6.0;

import "./contracts/token/ERC20/ERC20.sol";

contract ETHC is ERC20 {
    constructor() ERC20("Ethcode", "ETHC") public {
        _mint(msg.sender, 100);
    }
}