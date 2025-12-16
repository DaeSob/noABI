/*
  Global 변수 관리
*/
const path = require('path');
// cli.js의 디렉토리를 기준으로 설정 (상위 2단계)
let _appDir = path.resolve(__dirname, '../../');
let _settingsFileName = './conf.d/settings.json';
let _objSettings = undefined;
let _bcc = 1000;              // 블럭생성 주기
let _gasBufferPercent = 5;    // gas fee buffer 5%, 이전 블럭의 gasPrice 에 5%를 더하여 제출한다
let _currentSet = undefined;
let _rpcUrl = undefined;
let _keyStore = undefined;
let _passPhrase = undefined;
let _privKey = undefined;
let _user = undefined;
let _currentCollection = '';
let _currentPath = '';
let _deployedContract = {};
let _varStore = {};
let _customChain = undefined;
let _defaultCallGasLimit = '0';
let _outputFormat = 'full';  // "full" or "minimal"
let _logs = undefined;

// 멀티라인 버퍼 상태
let _multilineBuffer = '';
let _multilineOriginalLines = [];  // 멀티라인 원본 라인들 저장 (history용)
let _multilineDepth = { paren: 0, bracket: 0, brace: 0 };
let _multilineInString = { single: false, double: false };
let _multilineEscape = false;
let _isMultilineInput = false;  // 멀티라인 입력 중인지 플래그

// 커스텀 History 관리
let _customHistory = [];
let _historyIndex = -1;
let _historyMaxSize = 1000;
let _historyFile = undefined;

// Global 변수 설정
function initializeGlobals() {
  global.appDir = _appDir.charAt(0).toLowerCase() + _appDir.slice(1);
  global.workPath = _appDir.charAt(0).toLowerCase() + _appDir.slice(1);
  global.settingsFileName = _settingsFileName;
  global.objSettings = _objSettings;
  global.blockCreationCycle = _bcc;
  global.defaultGasBufferPercent = _gasBufferPercent;
  global.gasBufferPercent = _gasBufferPercent;
  global.currentSet = _currentSet;
  global.rpcUrl = _rpcUrl;
  global.keyStore = _keyStore;
  global.passPhrase = _passPhrase;
  global.privKey = _privKey;
  global.user = _user;
  global.currentCollection = _currentCollection;
  global.currentPath = _currentPath;
  global.deployedContract = _deployedContract;
  global.varStore = _varStore;
  global.customChain = _customChain;
  global.defaultCallGasLimit = _defaultCallGasLimit;
  global.outputFormat = _outputFormat;
  global.snapshot = {};
  global.logs = _logs;
  global.customHistory = _customHistory;
}

// Getter/Setter 함수들
function getSettingsFileName() { return _settingsFileName; }
function setSettingsFileName(value) { 
  _settingsFileName = value;
  global.settingsFileName = value;
}

function getObjSettings() { return _objSettings; }
function setObjSettings(value) { 
  _objSettings = value;
  global.objSettings = value;
}

function getBcc() { return _bcc; }
function setBcc(value) { 
  _bcc = value;
  global.blockCreationCycle = value;
}

function getGasBufferPercent() { return _gasBufferPercent; }
function setGasBufferPercent(value) { 
  _gasBufferPercent = value;
  global.gasBufferPercent = value;
}

function getDefaultGasBufferPercent() { return _gasBufferPercent; }

function getCurrentSet() { return _currentSet; }
function setCurrentSet(value) { 
  _currentSet = value;
  global.currentSet = value;
}

function getRpcUrl() { return _rpcUrl; }
function setRpcUrl(value) { 
  _rpcUrl = value;
  global.rpcUrl = value;
}

function getKeyStore() { return _keyStore; }
function setKeyStore(value) { 
  _keyStore = value;
  global.keyStore = value;
}

function getPassPhrase() { return _passPhrase; }
function setPassPhrase(value) { 
  _passPhrase = value;
  global.passPhrase = value;
}

function getPrivKey() { return _privKey; }
function setPrivKey(value) { 
  _privKey = value;
  global.privKey = value;
}

function getUser() { return _user; }
function setUser(value) { 
  _user = value;
  global.user = value;
}

function getCurrentCollection() { return _currentCollection; }
function setCurrentCollection(value) { 
  _currentCollection = value;
  global.currentCollection = value;
}

function getCurrentPath() { return _currentPath; }
function setCurrentPath(value) { 
  _currentPath = value;
  global.currentPath = value;
}

function getDeployedContract() { return _deployedContract; }
function setDeployedContract(value) { 
  _deployedContract = value;
  global.deployedContract = value;
}

function getVarStore() { return _varStore; }
function setVarStore(value) { 
  _varStore = value;
  global.varStore = value;
}

function getCustomChain() { return _customChain; }
function setCustomChain(value) { 
  _customChain = value;
  global.customChain = value;
}

function getDefaultCallGasLimit() { return _defaultCallGasLimit; }
function setDefaultCallGasLimit(value) { 
  _defaultCallGasLimit = value;
  global.defaultCallGasLimit = value;
}

function getOutputFormat() { return _outputFormat; }
function setOutputFormat(value) { 
  _outputFormat = value;
  global.outputFormat = value;
}

function getLogs() { return _logs; }
function setLogs(value) { 
  _logs = value;
  global.logs = value;
}

// 멀티라인 상태 접근
function getMultilineBuffer() { return _multilineBuffer; }
function setMultilineBuffer(value) { _multilineBuffer = value; }

function getMultilineOriginalLines() { return _multilineOriginalLines; }
function setMultilineOriginalLines(value) { _multilineOriginalLines = value; }

function getMultilineDepth() { return _multilineDepth; }
function setMultilineDepth(value) { _multilineDepth = value; }

function getMultilineInString() { return _multilineInString; }
function setMultilineInString(value) { _multilineInString = value; }

function getMultilineEscape() { return _multilineEscape; }
function setMultilineEscape(value) { _multilineEscape = value; }

function getIsMultilineInput() { return _isMultilineInput; }
function setIsMultilineInput(value) { _isMultilineInput = value; }

// History 상태 접근
function getCustomHistory() { return _customHistory; }
function setCustomHistory(value) { 
  _customHistory = value;
  global.customHistory = value;
}

function getHistoryIndex() { return _historyIndex; }
function setHistoryIndex(value) { _historyIndex = value; }

function getHistoryMaxSize() { return _historyMaxSize; }
function getHistoryFile() { return _historyFile; }
function setHistoryFile(value) { _historyFile = value; }

// workPath는 별도로 관리 (objSettings에서 설정됨)
function getWorkPath() { return global.workPath; }
function setWorkPath(value) { global.workPath = value; }

module.exports = {
  initializeGlobals,
  getSettingsFileName,
  setSettingsFileName,
  getObjSettings,
  setObjSettings,
  getBcc,
  setBcc,
  getGasBufferPercent,
  setGasBufferPercent,
  getDefaultGasBufferPercent,
  getCurrentSet,
  setCurrentSet,
  getRpcUrl,
  setRpcUrl,
  getKeyStore,
  setKeyStore,
  getPassPhrase,
  setPassPhrase,
  getPrivKey,
  setPrivKey,
  getUser,
  setUser,
  getCurrentCollection,
  setCurrentCollection,
  getCurrentPath,
  setCurrentPath,
  getDeployedContract,
  setDeployedContract,
  getVarStore,
  setVarStore,
  getCustomChain,
  setCustomChain,
  getDefaultCallGasLimit,
  setDefaultCallGasLimit,
  getOutputFormat,
  setOutputFormat,
  getLogs,
  setLogs,
  getMultilineBuffer,
  setMultilineBuffer,
  getMultilineOriginalLines,
  setMultilineOriginalLines,
  getMultilineDepth,
  setMultilineDepth,
  getMultilineInString,
  setMultilineInString,
  getMultilineEscape,
  setMultilineEscape,
  getIsMultilineInput,
  setIsMultilineInput,
  getCustomHistory,
  setCustomHistory,
  getHistoryIndex,
  setHistoryIndex,
  getHistoryMaxSize,
  getHistoryFile,
  setHistoryFile,
  getWorkPath,
  setWorkPath
};

