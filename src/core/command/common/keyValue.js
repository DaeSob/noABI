const { OpenWalletFromPrivateKey } = require('../../../lib/import.js');

const reservedVarName = [
  'STORE',
  'INDEX',
  'CONTRACT',
  'USER'
];

/**
 * 중첩된 키로 값을 가져오는 함수
 * @param {string} keyPath - "key.subkey.subsubkey" 형식의 키 경로
 * @param {object} store - varStore 객체
 * @returns {any} 찾은 값 또는 undefined
 */
function _getNestedValue(keyPath, store) {
  const keys = keyPath.split('.');
  let value = store[keys[0]];

  if (value === undefined) {
    return undefined;
  }

  // 나머지 키들을 순회하며 중첩 접근
  for (let i = 1; i < keys.length; i++) {
    const key = keys[i];
    // 배열 인덱스 처리 (예: "key[0]")
    const arrayIndexMatch = key.match(/^(.+)\[(\d+)\]$/);

    if (arrayIndexMatch) {
      const arrayKey = arrayIndexMatch[1];
      const index = parseInt(arrayIndexMatch[2]);
      if (value[arrayKey] !== undefined && Array.isArray(value[arrayKey])) {
        value = value[arrayKey][index];
      } else {
        return undefined;
      }
    } else {
      if (value[key] !== undefined) {
        value = value[key];
      } else {
        return undefined;
      }
    }

    if (value === undefined) {
      return undefined;
    }
  }

  return value;
}

/**
 * => 연산자로 변수명 추출
 * @param {string} _line - 입력 라인
 * @returns {object} { storeVarName: string|null, commandLine: string }
 */
function _extractStoreVarName(_line) {
  let storeVarName = null;
  let commandLine = _line || '';

  if (_line) {
    const arrowMatch = _line.match(/=>\s*(\w+)/);
    if (arrowMatch && arrowMatch[1]) {
      storeVarName = arrowMatch[1].trim();
      // => 이전 부분만 추출
      const arrowIndex = _line.indexOf('=>');
      if (arrowIndex > -1) {
        commandLine = _line.substring(0, arrowIndex).trim();
      }
    }
  }

  return { storeVarName, commandLine };
}

/**
 * varStore에 값을 저장 (예약 변수명 체크 포함)
 * @param {string} storeVarName - 저장할 변수명
 * @param {any} value - 저장할 값
 */
function _saveToVarStore(storeVarName, value) {
  if (storeVarName) {
    // 예약 변수명 체크
    if (reservedVarName.includes(storeVarName)) {
      throw new Error(`Error: "${storeVarName}" is a reserved variable name`);
    }
    // varStore에 저장
    global.varStore[storeVarName] = value;
  }
}

/**
 * Collection 이름 파싱
 * @param {object} parsedInput - 파싱된 입력
 * @returns {string} collection 이름
 */
function _getCollectionName(parsedInput) {
  const optionValue = parsedInput.opt['--collection'];
  let collectionName = undefined;

  if (optionValue !== undefined) {
    collectionName = optionValue;
  } else {
    collectionName = currentCollection;
  }

  if (collectionName === undefined || collectionName === '') {
    throw new Error('collection name selection required');
  }

  return collectionName;
}

module.exports.reservedVarName = reservedVarName;
module.exports._getNestedValue = _getNestedValue;
module.exports._extractStoreVarName = _extractStoreVarName;
module.exports._saveToVarStore = _saveToVarStore;
module.exports._getCollectionName = _getCollectionName;

function _getReservedWorldValue( _varName ) {
  try {
    if( _varName === 'USER' ) {
      return user;
    }
  } catch ( e ) {
    throw new String( 'SynTaxError: "' + _varName + '" is undefined');
  }
  return undefined;
}

function _valueOfKey( _key, _value ) {
  const regChildKey    = /(?<=\[)(.*?)(?=\])/g;
  let arrayIndex = _key.match( regChildKey );
  if( arrayIndex === null )
    return _value[_key];
  if(arrayIndex.length > 0 ) {
    let tempIdx = _key.indexOf( '[' );
    let tempKey = _key.substring( 0, tempIdx );
    let idx     = arrayIndex[0].trim().trimEnd()
    if( Array.isArray( _value[tempKey] ) ) {
      return _value[tempKey][Number(idx)];
    }
  }
  return undefined;
}

function _parseKeyValue( _keyValue ) {
  const regVarName    = /(?<=\$\{)(.*?)(?=\})/g;
  const varName       = _keyValue.match( regVarName );

  if( varName === undefined || varName === null || varName.length === 0 )
    return _keyValue;
  let varValue      = _getReservedWorldValue( varName[0] );
  if( varValue === undefined ) {
    let varTokenized  = varName[0].trim().trimEnd().split('.');      
    if( varTokenized.length > 0 ) {
      varValue = varStore[varTokenized[0]]
      if( varValue === undefined )
        throw new String( 'Error: "' + _keyValue + '" is undefined');
      for( let i=1; i<varTokenized.length; i++ ){
        varValue = _valueOfKey( varTokenized[i], varValue );
        if( varValue === undefined )
          throw new String( 'Error: "' + _keyValue + '" is undefined');
      }
    }
  }
  return varValue;
}

function _getArrayValue( _keyValue ) {
  let keyValue = []
  _keyValue.forEach( element => {
    if( Array.isArray( element ) ) {
      keyValue.push( _getArrayValue( element ) );
    } else {
      keyValue.push( _parseKeyValue( element ) );
    }
  })
  return keyValue;
}

function _getKeyValue( _keyValue ) {
  if( typeof _keyValue === 'object' ) {
    if( Array.isArray( _keyValue ) ) {
      return { type: "anonymous", value: _getArrayValue( _keyValue ) }
    } else {
      let types = Object.keys( _keyValue );
      if( types.length !== 1 ) {
        throw new String('SyntaxError: invalid parameters')
      }
      return {
        type: types[0],
        value: Array.isArray( _keyValue[types[0]] ) ? _getArrayValue( _keyValue[types[0]] ) : _parseKeyValue( _keyValue[types[0]] )        
      };
    }
  }
  return { type: "anonymous", value: _parseKeyValue( _keyValue )}
}

/**
 * .then() 또는 .catch() 콜백 파싱 (통합 함수)
 * @param {string} _inputLine - 입력 라인
 * @param {string} callbackType - 'then' 또는 'catch'
 * @returns {object} { funcArg: string|undefined, funcBody: string|undefined }
 */
function _getCallbackName(_inputLine, callbackType) {
  const callbackPattern = `.${callbackType}(`;
  const callbackIndex = _inputLine.indexOf(callbackPattern);
  if (callbackIndex === -1) {
    return { funcArg: undefined, funcBody: undefined };
  }

  // callbackType에 따라 시작 위치 계산
  const startPos = callbackIndex + callbackPattern.length;
  let depth = 0;
  let inString = { single: false, double: false };
  let escape = false;
  let funcArg = undefined;
  let funcBodyStart = -1;
  let funcBodyEnd = -1;

  // => 찾기
  let arrowIndex = -1;
  for (let i = startPos; i < _inputLine.length; i++) {
    const char = _inputLine[i];

    // Escape 문자 처리
    if (escape) {
      escape = false;
      continue;
    }
    if (char === '\\') {
      escape = true;
      continue;
    }

    // 문자열 내부 체크
    if (!inString.single && !inString.double) {
      if (char === "'") {
        inString.single = true;
        continue;
      }
      if (char === '"') {
        inString.double = true;
        continue;
      }
    } else {
      if (inString.single && char === "'") {
        inString.single = false;
        continue;
      }
      if (inString.double && char === '"') {
        inString.double = false;
        continue;
      }
      continue;
    }

    // => 찾기 (문자열 외부에서만)
    if (arrowIndex === -1 && char === '=' && i + 1 < _inputLine.length && _inputLine[i + 1] === '>') {
      arrowIndex = i;
      // funcArg 추출
      const argStr = _inputLine.substring(startPos, i).trim();
      funcArg = argStr || undefined;
      i++; // '>' 건너뛰기
      continue;
    }

    // 중괄호 depth 추적
    if (char === '{') {
      if (depth === 0 && arrowIndex !== -1) {
        funcBodyStart = i + 1; // '{' 다음부터
      }
      depth++;
    } else if (char === '}') {
      depth--;
      if (depth === 0 && funcBodyStart !== -1) {
        funcBodyEnd = i;
        break;
      }
    }
  }

  if (funcBodyStart === -1 || funcBodyEnd === -1) {
    // 중괄호가 없거나 완전하지 않음
    return { funcArg: funcArg, funcBody: undefined };
  }

  // funcBody 추출 (앞뒤 공백만 제거, 내부 줄바꿈은 보존)
  let funcBody = _inputLine.substring(funcBodyStart, funcBodyEnd);
  // 앞뒤 공백/줄바꿈 제거
  funcBody = funcBody.replace(/^\s+/, '').replace(/\s+$/, '');

  if (funcArg && reservedVarName.includes(funcArg)) {
    throw new String('SyntaxError: ' + funcArg + ' is a reserved world');
  }

  return { funcArg: funcArg, funcBody: funcBody };
}

function _getReturnsName(_inputLine) {
  return _getCallbackName(_inputLine, 'then');
}

function _getCatchName(_inputLine) {
  return _getCallbackName(_inputLine, 'catch');
}

module.exports._getKeyValue = _getKeyValue;
module.exports._getReturnsName = _getReturnsName;
module.exports._getCatchName = _getCatchName;
module.exports._parseKeyValue = _parseKeyValue;
