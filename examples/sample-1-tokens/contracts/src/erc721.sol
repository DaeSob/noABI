// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Strings }                               from '@openzeppelin/contracts/utils/Strings.sol';
import {IERC165}                                from "@openzeppelin/contracts/interfaces/IERC165.sol";
import {ERC721Enumerable, IERC721Enumerable }   from '@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol';
import {ERC721, IERC721 }                       from '@openzeppelin/contracts/token/ERC721/ERC721.sol';
import {Ownable}                                from "@openzeppelin/contracts/access/Ownable.sol";

contract CERC721 is Ownable, ERC721Enumerable {

using Strings for uint256;
using Strings for address;

    bool                                        private                     _bGuard;

    string                                      internal                    _baseURIextended;

    constructor( string memory name, string memory symbol, string memory tokenBaseUri, bool bAuto ) 
        ERC721( name, symbol ) {
            if( bAuto ) {
                _baseURIextended        = string(abi.encodePacked( tokenBaseUri, '/', block.chainid.toHexString(), '/', address(this).toHexString(), '/' ));
            } else {
                _baseURIextended        = string(abi.encodePacked( tokenBaseUri, '/'));
            }
    }

    modifier noReentrancy() {
        require( !_bGuard, "Reentrant call" );
        _bGuard = true;
        _;
        _bGuard = false;
    }    

    function setBaseURI( string calldata baseUri ) public onlyOwner {
        _baseURIextended = baseUri;
    }

    function tokenURI( uint256 tokenId ) public view virtual override returns (string memory) {
        if( _exists( tokenId ) ) {
            return string(abi.encodePacked( _baseURIextended, tokenId.toString()));
        }
        return string(abi.encodePacked( _baseURIextended, '0' ));
    }

    function mint( address to, uint256 tokenId ) public onlyOwner virtual {
        super._safeMint( to, tokenId );
    }
    
    function burn( uint256 tokenId ) public virtual {
        require(_isApprovedOrOwner(_msgSender(), tokenId), "ERC721: caller is not token owner or approved");
        super._burn( tokenId );
    }

}
