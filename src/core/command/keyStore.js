const Web3                              = require('web3');
const promptSync                        = require( 'prompt-sync' );
const { _getKeyValue, _getReturnsName } = require("../command/common/keyValue.js");
const { _getParseInput }                = require("./common/currentPath.js");
const { textColor, printBlue, printError, printSuccess, printGray, printDefault } = require("../log/printScreen.js");

var web3            = new Web3(new Web3.providers.HttpProvider());
const prompt        = promptSync();
  
async function _commandKeystore(_inputTokens) {
    let inputLine = Array.isArray(_inputTokens) ? _inputTokens.join(" ") : _inputTokens;

    let parsedInput     = _getParseInput( _inputTokens, 1 );

    // -create 옵션 체크 (필수)
    const isCreate = parsedInput.opt['-create'] !== undefined;
    
    if (!isCreate) {
        printError("Error: -create option is required");
        return;
    }    

    let pk = parsedInput.opt['--privateKey'];
    if( pk === undefined ) {
        pk  = prompt( 'private key: ', { echo: '*' });
    }

    let pwd = parsedInput.opt['--password'];
    if( pwd === undefined ) {
        pwd  = prompt( 'passphrase: ', { echo: '*' });
    }    
    
    let res = await web3.eth.accounts.encrypt(pk, pwd)
    console.log( JSON.stringify(res,null,2) )
}
  
module.exports._commandKeystore = _commandKeystore;
