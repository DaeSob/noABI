const fs = require("fs");
const { OpenWalletFromPrivateKey, } = require("../../lib/import.js");
const { _getParseInput } = require("./common/currentPath.js");
const { textColor, printError, printSuccess, printGray, printDefault } = require("../log/printScreen.js");

async function _commandCurrent(_inputTokens) {
  console.log( textColor.success + rpcUrl + textColor.white )
}
module.exports._commandCurrent = _commandCurrent;