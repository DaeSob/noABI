const fs = require("fs");
const path = require("path");

const { textColor, printError, printSuccess, printGray, printDefault } = require("../log/printScreen.js");

//------------------------------------------------------------------------------------------
//Compile
//deploy할 contract compile 하기
  
function _resolvePath( _appDir, _path ) {
  return path.resolve( _appDir, _path ).replace(/\\/g, "/");
}

function _compilingError(_errString, _errType, _errKey, _inFile) {
    return {
      errString: _errString,
      type: _errType,
      errKey: _errKey,
      inFile: _inFile,
    };
  }
  
  function _compiling(_binOutputFile, _sourceFiles, _deployContracts) {
    const rawBinOutputFile  = fs.readFileSync(_binOutputFile, "utf8");
    const binOutput         = JSON.parse(rawBinOutputFile);
  
    let err = [];
    let res = {};
    _sourceFiles.forEach((sourceFile) => {
      if (binOutput.contracts[sourceFile] !== undefined) {
        _deployContracts.forEach((deployContract) => {
          if (
            binOutput.contracts[sourceFile][deployContract.name] !== undefined
          ) {
            if (res[sourceFile] === undefined) {
              res[sourceFile] = {};
            }
            deployContract["used"] = true;
            const contractName =
              deployContract.alias === undefined
                ? deployContract.name
                : deployContract.alias;
            res[sourceFile][contractName] = {
              name: deployContract.name,
              abi: binOutput.contracts[sourceFile][deployContract.name].abi,
              object:
                binOutput.contracts[sourceFile][deployContract.name].evm.bytecode
                  .object,
              deploy: {
                extract: deployContract.extract, 
                ignore: deployContract.ignore,
                alias: deployContract.alias,
                metadata: deployContract.metadata,
                returns: deployContract.returns,
                inputs: deployContract.inputs,
              },
            };
          }
        });
      } else {
        err.push(
          _compilingError(
            "not found source file",
            "file",
            sourceFile,
            _binOutputFile
          )
        );
      }
    });
    _deployContracts.forEach((deployContract) => {
      if (deployContract.used === undefined) {
        err.push(
          _compilingError(
            "not found contract name",
            "contract",
            deployContract.name,
            _binOutputFile
          )
        );
      }
    });
    return { compiled: res, err: err };
  }
  
  function _importing(_binOutputFile, _importContracts, _compiled, _err) {
    const rawBinOutputFile = fs.readFileSync(_binOutputFile, "utf8");
    const binOutput = JSON.parse(rawBinOutputFile);
    _importContracts.forEach((contract) => {
      if (binOutput.contracts[contract.file] !== undefined) {
        if (binOutput.contracts[contract.file][contract.name] !== undefined) {
          if (_compiled[contract.file] === undefined) {
            _compiled[contract.file] = {};
          }
          _compiled[contract.file][contract.name] = {
            name: contract.name,
            abi: binOutput.contracts[contract.file][contract.name].abi,
            object:
              binOutput.contracts[contract.file][contract.name].evm.bytecode
                .object,
          };
        } else {
          _err.push(
            _compilingError(
              "not found contract name",
              "contract",
              contract.name,
              _binOutputFile
            )
          );
        }
      } else {
        _err.push(
          _compilingError(
            "not found source file",
            "file",
            contract.file,
            _binOutputFile
          )
        );
      }
    });
  }
  
  function _printCompileError(_err) {
    _err.forEach((err) => {
      console.log(textColor.error, err.errString + ":", err.errKey, "in", err.inFile);
    });
  }

  async function _compose( _objComposeFile ) {
         
    const regImport = /(?<=import)(.*?)(?=\;)/g;
    const regQuotation = /(?<=\')(.*?)(?=\')/g;

    let output = {};
    _objComposeFile.build.packages.forEach((element) => {
      let sourceFiles = [];
      
      let srcPath = element.source.entry;
      let soliditySrc = fs.readFileSync(srcPath, "utf8");
      let arraySrc = soliditySrc.match(regImport);

      let normalizedPath = element.source.normalize.path;
      if (normalizedPath == undefined ){
        console.log(textColor.error + "not found normalized path" + textColor.white);
      }
      console.log(textColor.blue + "compose " + srcPath + textColor.white );
      arraySrc.forEach((sourceFile) => {
        let srcFile  = sourceFile.match(regQuotation)[0];
        let filePath = _resolvePath( normalizedPath, srcFile.replace( /\.\.\//g, "./" ) );
        sourceFiles.push(filePath);
      });

      let resCompile = _compiling(
        element.artifacts.solcOutput,
        sourceFiles,
        element.contracts
      );
      
      if (element.import !== undefined) {
        for( let i=0; i<element.import.length; i++ ) {
          element.import[i].file =  element.import[i].file;
        }
        _importing(
          element.artifacts.solcOutput,
          element.import,
          resCompile.compiled,
          resCompile.err
        );
      }
      output = { ...output, ...resCompile.compiled };
      _printCompileError(resCompile.err);
      
    });
    
    const mkPath = path.dirname(_objComposeFile.files.manifest);
    fs.mkdirSync(mkPath, { recursive: true });
    fs.writeFileSync(_objComposeFile.files.manifest, JSON.stringify(output, null, 2));
    console.log(textColor.white + "written object file to " + _objComposeFile.files.manifest );
  
    console.log(textColor.white + "complete!");
  }
  
  module.exports._compose = _compose;  