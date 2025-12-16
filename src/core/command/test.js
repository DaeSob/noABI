const { _getKeyValue, _getReturnsName } = require("../command/common/keyValue.js");
const { _getParseInput } = require("./common/currentPath.js");
const { textColor, printBlue, printError, printSuccess, printGray, printDefault } = require("../log/printScreen.js");

const Web3 = require('web3');
var web3 = new Web3(new Web3.providers.HttpProvider());

function _parseRequest( _inputLine ) {
    return varStore;
}
  
async function _commandTest(_inputTokens) {
    let inputLine = Array.isArray(_inputTokens) ? _inputTokens.join(" ") : _inputTokens;

    let parsedInput     = _getParseInput( _inputTokens, 1 );

    let pk = parsedInput.opt['--privateKey'];
    let pwd = parsedInput.opt['--password'];
    
    let res = await web3.eth.accounts.encrypt(pk, pwd)
    console.log(res)
}
  
module.exports._commandTest = _commandTest;