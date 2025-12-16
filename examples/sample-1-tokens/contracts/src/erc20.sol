// SPDX-License-Identifier: MIT
// XEN Contracts v0.6.0
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract CERC20 is ERC20 {
    
    uint256           private                 _decimals;

    constructor( string memory name, string memory symbol, uint256 decimal ) ERC20(name, symbol){
        // mint 1,400,000,000 ether == 14억개
        if( decimal == 0 )
            _decimals = 18;
        else
            _decimals = decimal;
            
        _mint(msg.sender, 1400000000 * 10 ** _decimals);
    }

    function decimals() public view override returns (uint8) {
        return uint8(_decimals);
    }
    
}