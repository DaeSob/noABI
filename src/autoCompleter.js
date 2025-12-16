
const {parseLine, ParseLineOptions} = require('parse-line');

const { textColor, printError, printSuccess, printGray, printDefault } = require("./core/log/printScreen.js");
const { _tokenizePathToObject } = require("./core/command/common/tokenize.js");
const { _getCurrentPath, _getParseInput } = require("./core/command/common/currentPath.js");
const { knownCommands } = require("./core/command/knownCommands.js");

const tagEndCompletion = '<!--END--';

function _selectHits( _hits ) {
  if( _hits.length === 0 ) {
    return undefined; 
  }
 return _hits;
}

function _completionMatcher( _completions, _searchWord ) {
  if( _completions.length === 0 || _searchWord == undefined ) {
    return undefined;
  }
  const hits = _completions.filter((c) => c.startsWith(_searchWord));
  return _selectHits( hits );
}

function _execCompleter(_inputTokens, _line) {

  let parentPath  = undefined;
  let tempPath    = [];
  let completions = [];
  
  // 첫 번째 토큰이 명령어인지 확인
  let startIdx = 0;
  if (_inputTokens.length > 0 && knownCommands.includes(_inputTokens[0])) {
    startIdx = 1; // 명령어가 있으면 두 번째 토큰부터 시작
  }
  
  let parsedInput = _getParseInput( _inputTokens, startIdx );
  
  // show --abi 옵션 다음 contract name 자동완성 지원
  if (parsedInput.opt['--abi'] !== undefined) {
    let abiValue = parsedInput.opt['--abi'];
    let searchKey = abiValue || '';
    
    // currentCollection이 있으면 해당 collection의 contract name 목록 반환
    if (currentCollection && currentCollection.length > 0) {
      let tempObject = _tokenizePathToObject(
        currentCollection + '.',
        deployedContract
      );
      
      if (tempObject !== undefined) {
        completions = Object.keys(tempObject);
        completions.push(tagEndCompletion);
        return _execCompletion(completions, searchKey);
      }
    }
    
    // currentCollection이 없으면 collection 목록 반환
    completions = _useCompleter();
    return [ completions, _inputTokens ];
  }
  
  if( parsedInput.path === undefined ) {
    completions = _useCompleter();
    return [ completions, _inputTokens ];
  } else {
    const regParenthesis      = /(?<=\()(.*?)(?=\))/g;

    // 명령어가 있으면 _line에서 명령어 부분 제거
    let pathLine = _line;
    if (startIdx > 0 && _inputTokens.length > 0) {
      let cmdLength = _inputTokens[0].length;
      if (_line.length > cmdLength && _line.charAt(cmdLength) === ' ') {
        pathLine = _line.substring(cmdLength + 1).trim();
      } else if (_line.startsWith(_inputTokens[0])) {
        pathLine = _line.substring(cmdLength).trim();
      }
    }

    let tokenized = undefined
    let endsWithDot = pathLine.endsWith('.');
    if( pathLine.charAt(0) === '.' ) {
      tokenized = (_getCurrentPath() + pathLine);
      let match = pathLine.match( regParenthesis );
      if (match) {
        tokenized = tokenized.replace( '('+match+')', '' ).split( '.' );
      } else {
        tokenized = tokenized.split( '.' );
      }
    } else {
      let match = pathLine.match( regParenthesis );
      if (match) {
        tokenized = pathLine.replace( '('+match+')', '').split( '.' );
      } else {
        tokenized = pathLine.split( '.' );
      }
    }

    // 빈 문자열 제거
    tokenized = tokenized.filter(t => t.length > 0);

    let searchKey = [''];
    // 점으로 끝나는 경우 (자동완성용): 마지막 요소를 제거하지 않고 전체를 parentPath로 사용
    if (endsWithDot) {
      parentPath = tokenized.join('.') + '.';
    } else {
      // 점으로 끝나지 않는 경우: 마지막 요소를 searchKey로 사용
      searchKey = tokenized.length > 0 ? tokenized.splice( tokenized.length - 1, 1 ) : [''];
      parentPath = tokenized.join('.');
    }

    let tempObject  = _tokenizePathToObject(
      parentPath,
      deployedContract
    );

    if( tempObject === undefined ) {
      return [ [], searchKey ]  
    }

    let completions = [];
    if( Array.isArray( tempObject ) ) {
      tempObject.forEach((element) => {
        if( element.name !== undefined )
          completions.push(element.name);
      });
    } else {
      completions     = Object.keys( tempObject );
    }
    completions.push(tagEndCompletion);
    return _execCompletion( completions, searchKey[0] );
  }
}

function _execCompletion( _completions, _searchWord ) {
  let matched = _completionMatcher( _completions, _searchWord );
  if( matched === undefined || _searchWord === undefined  ) {
    return [ _completions, _searchWord ];
  }
  if(  matched.length > 1  ) {
    return [ matched, _searchWord ];
  }    
  return [ [matched[0]], _searchWord ];  
}

function _useCompleter() {
  let tempObject = _tokenizePathToObject(
    "",
    deployedContract
  );
  let completions = [];
  let res = Object.keys(tempObject);
  res.forEach((element) => {
    completions.push(element);
  });
  completions.push(tagEndCompletion)
  return completions;
}

function _cdCompleter(_inputTokens) {
  let parentPath  = undefined;
  if( currentPath.length === 0 ) { //최상위 경로가 아닌 경우
    parentPath      = _getCurrentPath();
  } else {
    let parsedInput = _getParseInput( _inputTokens, 1 );
    if( parsedInput.path === undefined ) {
      parentPath      = _getCurrentPath();
    } else {
      let tokenized   = parsedInput.path.split( '.' );
      tokenized.splice( tokenized.length - 1, 1 );
      parentPath      = tokenized.join('.');
    }
  }

  let tempObject  = _tokenizePathToObject(
                      parentPath,
                      deployedContract
                    );
  
  let subPath     = Object.keys( tempObject );

  for( i=0; i<subPath.length; i++ ) {
    subPath[i] = currentPath.length >  0 ? '.' + subPath[i] : subPath[i];
  }
  subPath.push(tagEndCompletion)  

  return subPath;
}

function _cmdCompletion( _line, _cmd, _completions, _searchWord ) {
  let matched = _completionMatcher( _completions, _searchWord );
  if( matched === undefined || _searchWord === undefined ) {
    return [ _completions, _line ];
  }
  if( matched.length > 1  ) {
    return [ _completions, _line ];
  }
  return [ [_cmd + ' ' + matched[0]], _line ];
}

function rlCompleter( _line ) {
  const inputTokens = parseLine(_line);
  let   completions = [];

  switch ( inputTokens[0] ) {
    case "use": {
      completions = _useCompleter();
      return _cmdCompletion( _line, inputTokens[0], completions, inputTokens[1]);
    }
    break;
    case "cd": {
      completions = _cdCompleter( inputTokens );
      return _cmdCompletion( _line, inputTokens[0], completions, inputTokens[1]);
    }
    break;    
    default: {
      return _execCompleter(inputTokens, _line);
    }
  }

}

module.exports._useCompleter = _useCompleter;
module.exports._execCompleter = _execCompleter;  
module.exports._cdCompleter = _cdCompleter;  
module.exports.rlCompleter = rlCompleter;  