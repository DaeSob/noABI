const _textColor = {
    blue: "\x1b[34m",
    white: "\x1b[37m",
    success: "\x1b[32m",
    warning: "\x1b[33m",
    error: "\x1b[31m",
    gray: "\x1b[90m",
    pink: "\x1b[35m",
    darkGolden : "\x1b[38;5;94m",
    
}

function _writeLog(_resKey, _resValue) {
    console.log(_textColor.white + _resKey, _textColor.success, _resValue);
}

function _printBlue( _msgString ) {
    console.log(_textColor.blue + _msgString + _textColor.white);
}

function _printError(_msgString) {
    console.log(_textColor.error + _msgString + _textColor.white);
}

function _printWarning(_msgString) {
    console.log(_textColor.warning + _msgString + _textColor.white);
}

function _printSuccess(_msgString) {
    console.log(_textColor.success + _msgString + _textColor.white);
}

function _printGray(_msgString) {
    console.log(_textColor.gray + _msgString + _textColor.white);
}

function _printDefault(_msgString) {
    console.log(_textColor.white, _msgString);
}

module.exports.textColor = _textColor;
module.exports.writeLog = _writeLog;
module.exports.printError = _printError;
module.exports.printWarning = _printWarning;
module.exports.printSuccess = _printSuccess;
module.exports.printBlue = _printBlue;
module.exports.printGray = _printGray;
module.exports.printDefault = _printDefault;
