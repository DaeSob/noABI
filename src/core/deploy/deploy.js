const fs    = require("fs");
const path  = require("path");
const {
  OpenWalletFromPrivateKey,
  OpenWalletFromKeystoreV3,
  OpenWalletFromMnemonic,
  OpenHDWallet,
}                     = require("../../lib/import.js");
const { _getKeyValue, 
  _getReturnsName }   = require("../command/common/keyValue.js");
const { sleep }       = require("../../lib/utils/sleep.js");
const { EthDeploy }   = require("../../lib/deploy.js");

const { textColor, printError, printSuccess, printGray, printDefault } = require("../log/printScreen.js");

//------------------------------------------------------------------------------------------
//Deploy
function _writeLog(_resKey, _color, _resValue) {
    console.log(textColor.white, _resKey, _color, _resValue, textColor.white);
  }
  
function _genInputValue( _types, _inputs ) {
    if (_inputs === undefined) return undefined;
    let ret = [];
    _inputs.forEach((element, idx) => {
      if( _types[idx] === 'bool' ) {
        if( element === 'false' || element === 'FALSE' || element === '0' )
          ret.push( false );
        else if( element === 'true' || element === 'TRUE' || element === '1' )        
          ret.push( true );
        else 
          throw new String("Error: invalid boolean type value");
      } else {
        let val = _getKeyValue( element );
        ret.push( val.value )
      }
    });
    return ret;
  }
  
  function _printDeployResult(_contract, _receipt) {
    _writeLog(" - alias:", textColor.success, _contract.deploy.alias ? _contract.deploy.alias : "" );
    _writeLog(" - BlockNumber:", textColor.success, _receipt.blockNumber);
    _writeLog(" - Transaction Hash:", textColor.success, _receipt.transactionHash);
    _writeLog(" - Gas Used:", textColor.success, parseInt(_receipt.gasUsed, 16 ));
    _writeLog(" - Contract Address:", textColor.success, _receipt.contractAddress);
    _writeLog(
      " - Status:",
      _receipt.status === "0x1" ? textColor.success : textColor.error,
      _receipt.status
    );
  }

  function _writeDeployLog( _contract, _receipt ) {

    var returns = {}
    if (_contract.deploy.returns !== undefined) {
      _contract.deploy.returns.varNames.forEach((element, index) => {
        receiptElement = _contract.deploy.returns.txReceipts[index];
        returns[element] = _receipt[receiptElement];
      });
    }

    let log = {
      contractName:     _contract.name,
      blockNumber:      _receipt.blockNumber,
      transactionHash:  _receipt.transactionHash,
      gasUsed:          parseInt(_receipt.gasUsed, 16 ),
      contractAddress:  _receipt.contractAddress,
      status:           _receipt.status,
      returns:          returns,
    }

    let today       = new Date();   
    let logFileName = currentSet.chain + '.' + today.getFullYear() + '.' + ( today.getMonth() + 1 ) + '.' + today.getDate() + '.txt';
    let logFile     = logs.deploy + '/' + logFileName;  

    var logData 
    if( fs.existsSync(logFile) ) {
      var stats   = fs.statSync(logFile)
      if( stats.size > 0 ) {
        logData = ',\r\n' + JSON.stringify(log, null, 2) 
      } else {
        logData = JSON.stringify(log, null, 2)
      }
    } else {
      logData = JSON.stringify(log, null, 2)
    }

    const mkPath = path.dirname(logFile);
    fs.mkdirSync(mkPath, { recursive: true });
    fs.appendFileSync(
      logFile,
      logData,
    );

  }

async function _deploy( _objComposefile, _optionalDeploy ) {
    const rawInputObjectfile  = fs.readFileSync(_objComposefile.files.manifest, "utf8");
    const objObjectfile       = JSON.parse(rawInputObjectfile);
    const wallet              = await OpenWalletFromPrivateKey(_objComposefile.wallet.privKey);
  
    let tempDeployedContracts = {};
  
    let sourceFiles = Object.keys(objObjectfile);
    for (let i = 0; i < sourceFiles.length; i++) {
      let contractAlias = Object.keys(objObjectfile[sourceFiles[i]]);
      for (let j = 0; j < contractAlias.length; j++) {
        let contract = objObjectfile[sourceFiles[i]][contractAlias[j]];
        
        if( contract.deploy === undefined || contract.deploy.ignore ) {
          console.log(textColor.warning + "import abi", contract.name);  
          tempDeployedContracts[contract.name] = {
            name: objObjectfile[sourceFiles[i]][contractAlias[j]].name,
            abi:  objObjectfile[sourceFiles[i]][contractAlias[j]].abi,
          }
          continue;
        }

        if( _optionalDeploy.length > 0 && !_optionalDeploy.includes( contractAlias[j] ) ) { //선택적 deploy를 한다
          continue;
        }
        
        console.log(textColor.blue + "deploy", contract.name);
        
        let inputTypes =
          contract.deploy.inputs !== undefined
            ? contract.deploy.inputs.types
            : undefined;
        let inputValue =
          contract.deploy.inputs !== undefined
            ? contract.deploy.inputs.value
            : undefined;
        inputValue = _genInputValue(inputTypes, inputValue);

        //-------------------------------------------------------------
        //send transaction Deploy

        let receipt = await EthDeploy({
          provider: {
            rpcUrl: _objComposefile.chain.rpc,
          },
          wallet: wallet,
          byteCode: contract.object,
          typesArray: inputTypes,
          parameters: inputValue,
        }); 

        //returns 저장
        if (contract.deploy.returns !== undefined) {
          contract.deploy.returns.varNames.forEach((element, index) => {
            receiptElement = contract.deploy.returns.txReceipts[index];
            varStore[element] = receipt[receiptElement];
          });
        }

        tempDeployedContracts[contract.name] = {
          name: objObjectfile[sourceFiles[i]][contractAlias[j]].name,
          abi:  objObjectfile[sourceFiles[i]][contractAlias[j]].abi,
        };
  
        _printDeployResult( contract, receipt );
        _writeDeployLog( contract, receipt );
        await sleep( 1000 );
      }
    }

    tempDeployedContracts = {
      contracts: tempDeployedContracts,
      index: {
        file: _objComposefile.files.index
      }
    }
    deployedContract[currentCollection] = tempDeployedContracts;

    console.log(textColor.white + "Complete!");
  }

  module.exports._deploy = _deploy;