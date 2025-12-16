const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');
const { _tokenizePathToObject } = require('./common/tokenize.js');
const { _getCurrentPath, _getParseInput } = require('./common/currentPath.js');

function _compileAbi(element) {
  let params = '';
  let returns = '';
  if (element.inputs !== undefined) {
    element.inputs.forEach((inputs, idx) => {
      if (idx > 0) {
        params += ',';
      }
      params += inputs.type;
    });
  }

  if (element.outputs !== undefined) {
    element.outputs.forEach((outputs, idx) => {
      if (idx > 0) {
        returns += ',';
      }
      returns += outputs.type;
    });
  }

  if (element.type === 'constructor') {
    console.log(
      textColor.blue,
      'constructor(' + textColor.white + params + textColor.blue + ')',
      (element.stateMutability !== undefined
        ? textColor.warning + element.stateMutability
        : '') + textColor.white
    );
  } else {
    console.log(
      textColor.blue,
      element.type,
      textColor.success +
      (element.name === undefined ? '(' : element.name + '(') +
        textColor.white +
        params +
        textColor.success +
        ')',
      element.stateMutability !== undefined
        ? textColor.warning + element.stateMutability
        : '',
      (returns.length > 0 ? textColor.gray + 'returns(' + returns + ')' : '') +
      textColor.white
    );
  }
}

function _filter(filter, _tempObject) {
  const filtering = [];
  for (let i = 0; i < filter.data.length; i++) {
    _tempObject.forEach(element => {
      if (filter.data[i] === element.name) {
        filtering.push(element);
      }
    });
  }
  return filtering;
}

async function _commandShow(_inputTokens) {
  let res = undefined;
  const parsedInput = _getParseInput(_inputTokens, 1);

  if (parsedInput.opt['-rpcUrl']) {
    printGray(rpcUrl);
    return;
  } else if (parsedInput.opt['-currentSet']) {
    printGray(JSON.stringify(currentSet));
    return;
  } else if (parsedInput.opt['-gasLimit']) {
    printGray(defaultCallGasLimit);
    return;
  }

  let filter = [];
  const filterIdx = parsedInput.data.findIndex(element => {
    if (element === '|') {
      return true;
    }
  });

  if (filterIdx > -1) {
    filter = parsedInput.data.slice(filterIdx);
  }

  if (parsedInput.path !== undefined && parsedInput.path.charAt(parsedInput.path.length - 1) === '.') {
    parsedInput.path = _getCurrentPath();
  }

  const tempObject = _tokenizePathToObject(
    parsedInput.path,
    deployedContract
  );

  if (tempObject === undefined) {
    throw new String('not exist');
  }

  if (typeof tempObject === 'string') {
    printSuccess(tempObject);
  } else {
    if (Array.isArray(tempObject)) {
      let filtering = [];
      if (filter.length > 0) {
        tempObject.forEach(element => {
          if (filter.includes(element.name)) {
            filtering.push(element);
          }
        });
      } else {
        filtering = tempObject;
      }

      if (parsedInput.opt['-json']) {
        printSuccess(JSON.stringify(filtering, null, 2));
      } else {
        filtering.forEach((element) => {
          _compileAbi(element);
        });
      }
    } else {
      res = Object.keys(tempObject);
      res.forEach((element) => {
        if (filter.length > 0) {
          if (filter.includes(element)) {
            printSuccess(element);
          }
        } else {
          printSuccess(element);
        }
      });
    }
  }
}

module.exports._commandShow = _commandShow;