const { textColor, printError, printSuccess, printGray, printDefault } = require("./core/log/printScreen.js");

function _welcomeHead( rpc, user ) {

  console.log( textColor.darkGolden + "################################################################################" + textColor.white )
  console.log( textColor.darkGolden + "Welcome to noABI REPL" + textColor.white )
  console.log( textColor.darkGolden + "https://github.com/DaeSob/noAbi" + textColor.white )
  console.log( textColor.darkGolden + "################################################################################" + textColor.white )
  console.log("");
  console.log( textColor.gray + "Usage:" + textColor.white );
  console.log( textColor.gray + "   noAbi [--config <file>]" + textColor.white ); 
}


function _welcome( rpc, user ) {

  console.log("");
  console.log( textColor.gray + "Current Environment:" + textColor.white )
  console.log( textColor.success + "   Active RPC   : " + textColor.blue + rpc + textColor.white )
  console.log( textColor.success + "   Active Wallet: " + textColor.blue + user + textColor.white )  

  console.log("");
  console.log( textColor.darkGolden + "type " + textColor.blue + "'help'" + textColor.white + " to see available commands.");
  console.log("");

}

module.exports._welcomeHead = _welcomeHead;
module.exports._welcome = _welcome;