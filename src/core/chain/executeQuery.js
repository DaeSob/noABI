
const Web3Utils                                     = require('web3-utils');
const Web3ABI                                       = require('web3-eth-abi');
const {
    OpenWalletFromPrivateKey,
    OpenWalletFromKeystoreV3,
    OpenWalletFromMnemonic,
    OpenHDWallet,
  } = require("../../lib/import.js");
const {RpcRequest}                                  = require('../../lib/rpc/rpcRequest.js');
  
const { KlaySendTransaction, EthSendTransaction }   = require("../../lib/sendtx.js");
const { EncodeFunctionCall }                        = require("../../lib/utils/encoder.js");
const { _getCurrentPath, _getParseInput }           = require("../command/common/currentPath.js");
const { _getKeyValue, _getReturnsName, _getCatchName }             = require("../command/common/keyValue.js");
const { textColor, printError, printSuccess, printGray, printDefault } = require("../log/printScreen.js");
const vm = require('vm');

  /**
   * ABI outputs와 결과를 매핑하여 이름이 있는 객체로 변환
   * tuple 타입도 재귀적으로 처리
   * @param {Array} outputs - ABI outputs 배열
   * @param {Object} decodedResult - 디코딩된 결과 객체 (Web3ABI.decodeParameters 결과)
   * @returns {Object} 매핑된 결과 객체
   */
  function _mapAbiOutputs(outputs, decodedResult) {
    if (!outputs || !Array.isArray(outputs) || outputs.length === 0) {
      return decodedResult;
    }

    const mapped = {};

    outputs.forEach((output, idx) => {
      // Web3ABI.decodeParameters는 숫자 인덱스로 접근 가능
      const key = output.name || idx.toString();
      let value = decodedResult[idx];

      // tuple 타입 처리
      if (output.type === 'tuple' && output.components) {
        // tuple의 경우 value는 배열이므로 components와 매핑
        if (Array.isArray(value) || (typeof value === 'object' && value !== null && !value.__length__)) {
          // tuple의 각 컴포넌트를 재귀적으로 매핑
          const tupleComponents = output.components;
          const tupleMapped = {};
          
          tupleComponents.forEach((component, compIdx) => {
            const compKey = component.name || compIdx.toString();
            let compValue = Array.isArray(value) ? value[compIdx] : value[compIdx];
            
            // 중첩된 tuple 처리
            if (component.type === 'tuple' && component.components) {
              const nestedTuple = Array.isArray(compValue) ? compValue : [compValue];
              compValue = _mapAbiOutputs([component], nestedTuple)[0];
            }
            
            tupleMapped[compKey] = compValue;
          });
          
          value = tupleMapped;
        }
      } else if (output.type && output.type.startsWith('tuple[')) {
        // tuple 배열 처리
        if (Array.isArray(value) && output.components) {
          value = value.map((item) => {
            if (Array.isArray(item) || (typeof item === 'object' && item !== null && !item.__length__)) {
              const tupleMapped = {};
              output.components.forEach((component, compIdx) => {
                const compKey = component.name || compIdx.toString();
                let compValue = Array.isArray(item) ? item[compIdx] : item[compIdx];
                
                if (component.type === 'tuple' && component.components) {
                  const nestedTuple = Array.isArray(compValue) ? compValue : [compValue];
                  compValue = _mapAbiOutputs([component], nestedTuple)[0];
                }
                
                tupleMapped[compKey] = compValue;
              });
              return tupleMapped;
            }
            return item;
          });
        }
      }

      mapped[key] = value;
    });

    return mapped;
  }

  /**
   * Node.js 런타임에서 .then() 또는 .catch() 블록 코드 실행
   * @param {string} funcBody - 실행할 JavaScript 코드
   * @param {any} res - Contract Call 결과 (성공 시) 또는 undefined (에러 시)
   * @param {object} memory - DSL memory (global.varStore)
   * @param {string} varName - 변수명 (예: res, myResult, err 등)
   * @param {object} err - 에러 객체 (catch 블록에서만 사용)
   */
  function _executeNodeJsBlock(funcBody, res, memory, varName, err) {
    if (!funcBody || funcBody.trim() === '') {
      return;
    }
    
    try {
      // 허용된 Node.js 모듈 목록 (보안을 위해 제한)
      const allowedModules = [
        'assert', 'util', 'path', 'url', 'crypto', 'buffer',
        'fs', 'os', 'events', 
        'zlib', 'querystring', 'string_decoder', 'timers', 'perf_hooks'
      ];
      
      // require 함수 (제한된 모듈만 허용)
      const safeRequire = (moduleName) => {
        if (allowedModules.includes(moduleName)) {
          try {
            return require(moduleName);
          } catch (e) {
            throw new Error(`Failed to require module "${moduleName}": ${e.message}`);
          }
        }
        throw new Error(`Module "${moduleName}" is not allowed. Allowed modules: ${allowedModules.join(', ')}`);
      };
      
      // Sandbox context 생성
      const context = {
        memory: memory,  // DSL memory 접근
        varStore: memory,  // varStore 직접 접근 (동적 키 사용 가능)
        console: console,  // console.log 사용 가능
        require: safeRequire,  // 제한된 require 지원
        // 유틸리티 함수들
        JSON: JSON,
        Math: Math,
        Date: Date,
        String: String,
        Number: Number,
        Boolean: Boolean,
        Array: Array,
        Object: Object,
        parseInt: parseInt,
        parseFloat: parseFloat,
        isNaN: isNaN,
        isFinite: isFinite
      };
      
      // 변수명으로 Contract Call 결과 또는 에러 전달
      if (varName) {
        // catch 블록인 경우 err 전달, then 블록인 경우 res 전달
        context[varName] = err !== undefined ? err : res;
      } else {
        // 변수명이 없으면 기본적으로 'res' 또는 'err' 사용
        if (err !== undefined) {
          context.err = err;
        } else {
          context.res = res;
        }
      }
      
      // VM context 생성
      const vmContext = vm.createContext(context);
      
      // 코드 실행 (타임아웃 10초)
      const script = new vm.Script(funcBody);
      script.runInContext(vmContext, {
        timeout: 10000,
        displayErrors: true
      });
      
    } catch (e) {
      // 더 자세한 에러 메시지 제공
      const errorMessage = e.message || String(e);
      const errorStack = e.stack ? `\nStack: ${e.stack}` : '';
      const blockType = err !== undefined ? '.catch()' : '.then()';
      throw new Error(`Error in ${blockType} block: ${errorMessage}${errorStack}`);
    }
  }

  async function _ethEstimateGas(_rpcUrl, _wallet, _to, _data, _value) {

    let resGasPrice     = {};
    let resMaxFeePerGas = {};
    let resGas          = {};

    if( customChain === undefined || customChain.gasPrice === undefined ) {
        resGasPrice     = await RpcRequest( _rpcUrl, "eth_gasPrice", [], 2 );
        resMaxFeePerGas = await RpcRequest( _rpcUrl, "eth_getBlockByNumber", ["latest", false], 3 );     
        
        const toHex       = (bi) => '0x' + bi.toString(16);
        let   gasPriceBI  = resGasPrice?.result ? BigInt(resGasPrice.result) : undefined;
        let   baseFeeBI   = resMaxFeePerGas?.result?.baseFeePerGas ? BigInt(resMaxFeePerGas.result.baseFeePerGas) : undefined;

        // 둘 중 하나만 있는 경우
        if (!gasPriceBI && baseFeeBI) {
            gasPriceBI = baseFeeBI;  // gasPrice fallback
        }
        if (!baseFeeBI && gasPriceBI) {
            baseFeeBI = gasPriceBI;  // baseFee fallback
        }

        //gas buffer 를 더한다
        const feeBuffer = BigInt(100 + gasBufferPercent);
        const maxGasPriceBI   = gasPriceBI * feeBuffer / 100n;
        const maxFeePerGasBI  = baseFeeBI * feeBuffer / 100n;
        const maxGasPriceHex  = toHex(maxGasPriceBI);
        const maxFeePerGasHex = toHex(maxFeePerGasBI);
        
        resGasPrice.result                    = maxGasPriceHex;
        resMaxFeePerGas.result.baseFeePerGas  = maxFeePerGasHex;

    } else {
      resGasPrice['result']    = Web3Utils.toHex(customChain.gasPrice);
      resMaxFeePerGas['result'] = {
        baseFeePerGas: Web3Utils.toHex(customChain.gasPrice)
      }      
    }

    if( defaultCallGasLimit === "0" ) {
      resGas          = await RpcRequest( _rpcUrl, "eth_estimateGas", [{
                                from:       _wallet.getAddressString(),
                                to:         _to,
                                data:       _data,
                                gasPrice:   "0x0",
                                value:      _value ||'0x0',                                
                            }], 4 );
      if( resGas.error !== undefined )
          throw resGas.error;
    } else {
      resGas['result'] = Web3Utils.toHex(defaultCallGasLimit);
    }
    
    return { gasPrice: resGasPrice.result, gas: resGas.result, maxFeePerGas: resMaxFeePerGas.result.baseFeePerGas };
  }  
 
  /**
   * RPC 에러를 객체로 변환
   * @param {any} error - 에러 객체 또는 RPC 응답
   * @returns {object} 에러 객체
   */
  function _extractError(error) {
    // 이미 객체인 경우
    if (error && typeof error === 'object' && !(error instanceof Error)) {
      // RPC 에러 응답인 경우
      if (error.error) {
        return error.error;
      }
      // 그 외 객체는 그대로 반환
      return error;
    }
    
    // Error 객체인 경우
    if (error instanceof Error) {
      return {
        code: -32603, // Internal error
        message: error.message,
        data: error.stack
      };
    }
    
    // 문자열인 경우
    if (typeof error === 'string') {
      return {
        code: -32603,
        message: error
      };
    }
    
    // 알 수 없는 형태
    return {
      code: -32603,
      message: String(error)
    };
  }

  async function _call(_rpcUrl, _wallet, _to, _data, _outputTypes) {

    let res = await RpcRequest( _rpcUrl, "eth_call", [{
                  from:       _wallet.getAddressString(),
                  to:         _to,
                  data:       _data,
                },"latest"], 1 );
    
    // RPC 에러 체크
    if (res.error) {
      throw res.error;
    }
    
    if (res.result === undefined || res.result === null) {
      throw {
        code: -32603,
        message: "eth_call returned no result"
      };
    }
    
    var dec = Web3ABI.decodeParameters(_outputTypes, res.result );
    return dec;
  }
  
  async function _ethInvoke( _rpcUrl, _wallet, _to, _data, _value) {
    let   res         = await RpcRequest(_rpcUrl, "eth_getTransactionCount", [_wallet.getAddressString(),"pending"], 4 );
    const count       = res.result;    
    const estimateGas = await _ethEstimateGas(_rpcUrl, _wallet, _to, _data, _value);

    return await EthSendTransaction({
      provider: {
        rpcUrl: _rpcUrl
      },
      wallet:   _wallet,
      bcc: blockCreationCycle,  //블록생성 주기
      rawTransaction: {
          from:         _wallet.getAddressString(),
          to:           _to,
          nonce:        count,
          value:        _value.length == 0 ? '0x00' : Web3Utils.toHex(_value),
          gasLimit:     estimateGas.gas,
          maxFeePerGas: estimateGas.maxFeePerGas,
          gasPrice:     estimateGas.gasPrice,
          data:         _data,
      } 
    })
  }  

 async function _query(_inputTokens, _line) {

    if (_inputTokens.length === 0) {
      return;
    }

    let queryTokenized  = undefined;
    let wallet          = undefined;
    try {

      queryTokenized  = await _queryTokenize(_inputTokens, _line);
      wallet          = await OpenWalletFromPrivateKey(privKey);

      if (
        queryTokenized !== undefined &&
        queryTokenized.funcType === "function"
      ) {
        if (queryTokenized.funcStateMutability === "view" || queryTokenized.funcStateMutability === "pure") {
          try {
            let res = await _call(
              rpcUrl,
              wallet,
              queryTokenized.to,
              queryTokenized.data,
              queryTokenized.funcOutputTypes
            );
            
            // RESULT 섹션 (ABI와 매핑)
            const mappedResult = _mapAbiOutputs(queryTokenized.funcOutputTypes, res);
            
            // .then() 블록이 있는지 확인 (funcReturns가 있고 funcArg가 있으면 .then() 블록 존재)
            const hasThenBlock = queryTokenized.funcReturns !== undefined && 
                                  queryTokenized.funcReturns.length > 0 &&
                                  queryTokenized.funcReturns[0].funcArg !== undefined;
            
            // .then() 블록이 없을 때만 결과 출력
            if (!hasThenBlock) {
              // 출력 형식에 따라 결정
              const outputFormat = global.outputFormat || 'full';
              
              if (outputFormat === 'minimal') {
                // minimal: result만 출력
                console.log(JSON.stringify(mappedResult, null, 2));
              } else {
                // full: call + result 출력 (기본값)
                const callObject = {
                  to: queryTokenized.to,
                  function: queryTokenized.funcName
                };
                
                // value가 있을 때만 추가
                if (queryTokenized.value !== undefined && queryTokenized.value !== null && queryTokenized.value !== '') {
                  callObject.value = queryTokenized.value;
                }
                
                // params 추가
                if (queryTokenized.funcParams && queryTokenized.funcParams.length > 0) {
                  callObject.params = queryTokenized.funcParams;
                } else {
                  callObject.params = [];
                }
                
                const outputObject = {
                  call: callObject,
                  result: mappedResult
                };
                
                console.log(JSON.stringify(outputObject, null, 2));
              }
            }

            // .then() 블록 처리
            if( queryTokenized.funcReturns !== undefined ) {
              if( queryTokenized.funcReturns.length > 0 ) {
                const funcReturn = queryTokenized.funcReturns[0];
                
                // funcArg가 있으면 varStore에 저장 (funcBody 유무와 관계없이)
                // funcArg는 사용자가 지정한 변수명 (예: res, myResult, data 등)
                if (funcReturn.funcArg) {
                  varStore[funcReturn.funcArg] = mappedResult;
                }
                
                // funcBody가 있으면 Node.js에서 실행
                // funcArg를 변수명으로 사용하여 context에 전달
                if (funcReturn.funcBody !== undefined && funcReturn.funcBody.trim() !== '') {
                  _executeNodeJsBlock(funcReturn.funcBody, mappedResult, global.varStore, funcReturn.funcArg);
                }
              }
            }
          } catch (callError) {
            // .catch() 블록이 있는지 확인
            const hasCatchBlock = queryTokenized.funcCatch !== undefined && 
                                  queryTokenized.funcCatch.length > 0;
            
            if (hasCatchBlock) {
              // 에러 객체 추출 후 message만 전달
              const errObj = _extractError(callError);
              const errMessage = errObj.message || String(callError);
              
              // .catch() 블록 실행 (message만 전달)
              const funcCatch = queryTokenized.funcCatch[0];
              
              // funcArg가 있으면 varStore에 저장하지 않음 (catch는 에러만 전달)
              // funcBody가 있으면 Node.js에서 실행
              if (funcCatch.funcBody !== undefined && funcCatch.funcBody.trim() !== '') {
                // err.message만 전달 (객체 형태로)
                _executeNodeJsBlock(funcCatch.funcBody, undefined, global.varStore, funcCatch.funcArg, { message: errMessage });
              }
            } else {
              // .catch() 블록이 없으면 에러 메시지만 추출하여 throw
              const errObj = _extractError(callError);
              const errorMessage = errObj.message || String(callError);
              throw new Error(errorMessage);
            }
          }

        } else {
          try {
            let res = await _ethInvoke(
              rpcUrl,
              wallet,
              queryTokenized.to,
              queryTokenized.data,
              queryTokenized.value
            );

            if( res.gasUsed !== undefined ) {
              res.gasUsed = parseInt(res.gasUsed, 16 )
            }
            if( res.cumulativeGasUsed !== undefined ) {
              res.cumulativeGasUsed = parseInt(res.cumulativeGasUsed, 16 )
            }
            //console.log( res );
            textColor
            console.log("status:", ( (res.status === '0x1') ? textColor.success : textColor.error) + res.status + textColor.white);
            console.log("txHash:", res.transactionHash );
          } catch (invokeError) {
            // .catch() 블록이 있는지 확인
            const hasCatchBlock = queryTokenized.funcCatch !== undefined && 
                                  queryTokenized.funcCatch.length > 0;
            
            if (hasCatchBlock) {
              // 에러 객체 추출 후 message만 전달
              const errObj = _extractError(invokeError);
              const errMessage = errObj.message || String(invokeError);
              
              // .catch() 블록 실행 (message만 전달)
              const funcCatch = queryTokenized.funcCatch[0];
              
              // funcBody가 있으면 Node.js에서 실행
              if (funcCatch.funcBody !== undefined && funcCatch.funcBody.trim() !== '') {
                // err.message만 전달 (객체 형태로)
                _executeNodeJsBlock(funcCatch.funcBody, undefined, global.varStore, funcCatch.funcArg, { message: errMessage });
              }
            } else {
              // .catch() 블록이 없으면 에러 메시지만 추출하여 throw
              const errObj = _extractError(invokeError);
              const errorMessage = errObj.message || String(invokeError);
              throw new Error(errorMessage);
            }
          }
        }
      }
    } catch( e ) {
      throw e;
    } finally {
      wallet = undefined;
    }
  }

  //const   regRequestFormat    = /((.*?)\.contracts\.(.*?)\.(.*?)\((.*?)\))/g;
  // 변경: collection.contractname(address).functionName(...) 형식 지원
  // functionName() 또는 functionName(...) 모두 지원 (빈 괄호도 허용)
  // 멀티라인 지원: 괄호와 점 사이의 공백 허용 (\s*)
  const   regRequestFormat    = /((.*?)\.(.*?)\((.*?)\)\s*\.\s*(.*?)\((.*?)\))/g;
  const   regThenFormat       = /(\.then\((.*?)\=\>(.*?)\))/g
  const   regParenthesis      = /(?<=\()(.*?)((?=(?!\)\")(\))))/g;
  const   regFunction         = /(\(\))/g;

  function _compileRequest( _inputLine ) {

    let     tempRequest       = _inputLine.match( regRequestFormat );
    if (tempRequest === null || tempRequest.length === 0) {
      throw new String("Error: invalid query format");
    }
    let     tempRequestParams = tempRequest[0].match( regParenthesis );
    if (tempRequestParams === null) {
      tempRequestParams = [];
    }
    // .then() 블록 파싱 (멀티라인 지원)
    let tempThen = null;
    if (_inputLine.indexOf('.then(') !== -1) {
      // .then() 블록이 있으면 _getReturnsName으로 파싱
      const thenResult = _getReturnsName(_inputLine);
      if (thenResult.funcArg !== undefined || thenResult.funcBody !== undefined) {
        tempThen = [thenResult];
      }
    }
    
    // .catch() 블록 파싱 (멀티라인 지원)
    let tempCatch = null;
    if (_inputLine.indexOf('.catch(') !== -1) {
      // .catch() 블록이 있으면 _getCatchName으로 파싱
      const catchResult = _getCatchName(_inputLine);
      if (catchResult.funcArg !== undefined || catchResult.funcBody !== undefined) {
        tempCatch = [catchResult];
      }
    }
    
    let     tempTokenized     = tempRequest[0].replace( regParenthesis, '' ).split( '.' );

    let     reqParamsIdx      = 0;    
    let     parsedRequest     = {
      address: '',
      value: '',
      request: {
        tokens:[],
        params:{}
      },
      then: [],
      catch: []
    };

    tempTokenized.forEach( (element, idx) => {
      let   token = element.replace(regFunction, '').trimEnd();
      parsedRequest.request.tokens.push(token);

      if( element.includes('()') ) {  
        parsedRequest.request.params[token] = {
          var: tempRequestParams[reqParamsIdx] === undefined ? [] : JSON.parse( '[' + tempRequestParams[reqParamsIdx].trim().trimEnd() + ']' ),
          val: [] 
        }
        reqParamsIdx++;

        parsedRequest.request.params[token].var.forEach( element => {

          if( idx === 1 ) {
            element = element.trim()
            if( element.indexOf( "${" ) != 0 && element.indexOf( "{" ) > -1 ) {
              str = element.replace(/\'/g, "\"");
              parsed = JSON.parse(str)
              if(parsed.value !== undefined && parsed.value.length > 0 ) {
                parsedRequest.value =_getKeyValue( typeof parsed.value === 'number' ?  parsed.value.toString():  parsed.value ).value
                parsedRequest.value = Web3Utils.toHex(parsedRequest.value);
              }
              if(parsed.to !== undefined && parsed.to.length > 0){
                parsedRequest.request.params[token].val.push( _getKeyValue( typeof parsed.to === 'number' ? parsed.to.toString(): parsed.to ) );
                parsedRequest.address = parsedRequest.request.params[token].val[0].value;                
              }              
            } else {
              parsedRequest.request.params[token].val.push( _getKeyValue( typeof element === 'number' ? element.toString(): element ) );
              parsedRequest.address = parsedRequest.request.params[token].val[0].value;
            }
            tempTokenized[idx] = tempTokenized[idx].split( '()' )[0];
          } else {
            parsedRequest.request.params[token].val.push( _getKeyValue( typeof element === 'number' ? element.toString(): element ) );
          }

        })
      } else {
        parsedRequest.request.params[token] = {
          var: undefined,
          val: []
        };
      }
    })
    
    if( tempThen !== null ) { 
      tempThen.forEach( element => {
        // element가 이미 파싱된 객체이면 그대로 사용, 아니면 파싱
        if (typeof element === 'object' && element.funcArg !== undefined) {
          parsedRequest.then.push(element);
        } else {
          parsedRequest.then.push(_getReturnsName(element));
        }
      })
    }
    
    if( tempCatch !== null ) { 
      tempCatch.forEach( element => {
        // element가 이미 파싱된 객체이면 그대로 사용, 아니면 파싱
        if (typeof element === 'object' && element.funcArg !== undefined) {
          parsedRequest.catch.push(element);
        } else {
          parsedRequest.catch.push(_getCatchName(element));
        }
      })
    }
    return parsedRequest;
  }
  
  function _findAbiFromContract( _tokenized ) {

    let objDeployed = deployedContract[_tokenized.request.tokens[0]];
    if (objDeployed === undefined) {
      throw new String("Error: collection " + _tokenized.request.tokens[0] + " is not defined");
    }

    objDeployed = objDeployed.contracts; 
    if (objDeployed === undefined) {
      throw new String("Error: contracts not found in collection");
    }

    objDeployed = objDeployed[_tokenized.request.tokens[1]]; 
    if ( objDeployed === undefined ) {
      throw new String("Error: "+ _tokenized.request.tokens[1] + " is not imported abi");
    }

    let funcName  = _tokenized.request.tokens[2];
    let idxAbi    = objDeployed.abi.findIndex((element) => {
      if ( element.name === funcName && element.inputs.length === _tokenized.request.params[funcName].val.length ) {
        let found = true;
        for( let i=0; i<element.inputs.length; i++ ) {
          if( _tokenized.request.params[funcName].val[i].type === 'anonymous' ) {
          } else if( element.inputs[i].type !== _tokenized.request.params[funcName].val[i].type ) {
            found = false;
            break;
          }
        }
        return found;
      }
    });
    if (idxAbi === -1) {
      throw new String("Error: not found abi");
    }    
    return objDeployed.abi[idxAbi];
  }

  async function _queryTokenize(_inputTokens, _line) {

    let inputLine = _line;//Array.isArray(_inputTokens) ? _inputTokens.join(" ") : _inputTokens;

    if( inputLine.charAt(0) === '.' ) {
      inputLine = _getCurrentPath() + inputLine;
    } else if( currentCollection.length > 0 )  {
      inputLine = inputLine;
    }

    let requestTokenized = _compileRequest( inputLine );

    if( requestTokenized.address.length == 0 ) {
      throw new String("Error: required contract address");
    }
    abi = _findAbiFromContract( requestTokenized );

    let funcInputTypes  = [];
    abi.inputs.forEach((element) => {
      funcInputTypes.push(element.type);
    });
    
    let funcParamsData  = []
    let funcName = requestTokenized.request.tokens[2];
    let funcParams = requestTokenized.request.params[funcName];
    
    if (funcParams === undefined) {
      funcParams = { var: [], val: [] };
    }
    
    funcParams.val.forEach( ( element, index ) =>{
      if( funcInputTypes[index] === 'bool' ) {
        if( element.value === 'false' || element.value === 'FALSE' || element.value === '0' )
          funcParamsData.push( false );
        else if( element.value === 'true' || element.value === 'TRUE' || element.value === '1' )        
          funcParamsData.push( true );
        else 
          throw new String("Error: invalid boolean type value");
      } else {
        funcParamsData.push( element.value );
      }

    })

    let encodedData     = EncodeFunctionCall(abi, funcParamsData);

    return {
      contractName: requestTokenized.request.tokens[1],
      to:  requestTokenized.address,
      value: requestTokenized.value,
      data: encodedData,
      funcType: abi.type,
      funcStateMutability: abi.stateMutability,
      funcName: requestTokenized.request.tokens[2],
      funcInputTypes: funcInputTypes,
      funcParams: funcParamsData,
      funcOutputTypes: abi.outputs,
      funcReturns: requestTokenized.then,
      funcCatch: requestTokenized.catch
    };
  }

  module.exports._query = _query;  
  