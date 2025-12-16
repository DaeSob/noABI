const fs = require("fs");
const { StandardMerkleTree } = require("@openzeppelin/merkle-tree");
const { _getParseInput } = require("./common/currentPath.js");
const { textColor, printError, printSuccess, printGray, printDefault } = require("../log/printScreen.js");

async function _commandMerkleTree(_inputTokens) {

  let parsedInput = _getParseInput( _inputTokens, 1 );
  let optionValue = undefined;

  if( parsedInput.opt['-gen'] ) {
    //white list 입력
    optionValue = parsedInput.opt['--in'];
    if( optionValue !== undefined ) { 
      let   rawWhiteList    = fs.readFileSync(optionValue, "utf8");
      let   jsonWhiteList   = JSON.parse( rawWhiteList );
      const tree            = StandardMerkleTree.of(jsonWhiteList.data, jsonWhiteList.leafEncoding );
      console.log("Tree:\n" + tree.render());
      //출력
      optionValue = parsedInput.opt['--out'];
      if( optionValue !== undefined ) { 
        fs.writeFileSync(optionValue, JSON.stringify(tree.dump()));
      }       
    }
  }

  if( parsedInput.opt['-get'] ) {
    optionValue = parsedInput.opt['--in'];
    if( optionValue !== undefined ) { 
      const merkleTree = StandardMerkleTree.load(JSON.parse(fs.readFileSync(optionValue)));
      optionValue = parsedInput.opt['--leaf'];
      if( optionValue !== undefined ) {         
        for (const [i, v] of merkleTree.entries()) {
          if (v[0] === optionValue ) {
            const proof = merkleTree.getProof(i);
            console.log(proof);
            return;
          }  
        }      
      }
    }
    console.log(textColor.error + "can not make a proof" + textColor.white);
    return; 
  }  
}
module.exports._commandMerkleTree = _commandMerkleTree;