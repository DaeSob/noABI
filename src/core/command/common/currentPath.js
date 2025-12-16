function _getCurrentPath() {
    if( currentCollection.length === 0 )
        return '';
    if( currentPath.length === 0 )
        return currentCollection;
    return currentCollection + '.' + currentPath;
}

function _getParseInput( _inputTokens, _idx ) {
    let searchPath      = undefined;
    let arrayData       = [];
    let opt = {};
    for (let i = _idx; i < _inputTokens.length; i++) {
        if (_inputTokens[i].indexOf('--') === 0 ) {//long 옵션 뒤에 옵션의 value 값이 따라 온다
          if( i + 1 >= _inputTokens.length ) {
            throw new Error('Error: invalid options format');
          }
          opt[_inputTokens[i]] = _inputTokens[i+1];
          i++;
        } else if (_inputTokens[i].charAt(0) === "-") {
          opt[_inputTokens[i]] = true;
        } else if (searchPath === undefined) {
          if( _getCurrentPath().length === 0 ) {
            searchPath = _inputTokens[i];
          } else if (_inputTokens[i].charAt(0) === ".") {
            searchPath = _getCurrentPath() + _inputTokens[i];
          } else {
            searchPath = _inputTokens[i];
          }
          arrayData.push( _inputTokens[i] );
        } else {
          arrayData.push( _inputTokens[i] );
        }
    }
    return { path: searchPath, opt: opt, data: arrayData }
}

module.exports._getCurrentPath = _getCurrentPath;
module.exports._getParseInput = _getParseInput;
