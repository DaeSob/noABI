const readline = require('readline');
const { rlCompleter } = require('../../autoCompleter.js');
const { _navigateHistory } = require('./history.js');

// 커스텀 History 키 이벤트 처리
let _currentInput = '';  // 현재 입력 중인 텍스트
let _isNavigatingHistory = false;  // history 탐색 중인지 플래그

let rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  completer: rlCompleter,
  historySize: 0  // readline history 완전 비활성화 (커스텀 history 사용)
});

// readline 인터페이스를 전역으로 export (keystore 등에서 사용)
global.rl = rl;

// readline의 input 스트림을 가로채서 키 이벤트 처리
// 상하 화살표 키를 직접 처리하여 커스텀 history 탐색
rl.input.on('keypress', (str, key) => {
  if (!key) return;

  // 상하 화살표 키 처리
  if (key.name === 'up') {
    if (rl.line) {
      _currentInput = rl.line;  // 현재 입력 저장
    }
    const historyLine = _navigateHistory('up');
    if (historyLine !== null) {
      _isNavigatingHistory = true;
      // 현재 라인 지우고 history 라인으로 교체
      rl.write(null, { ctrl: true, name: 'u' });  // Ctrl+U로 라인 지우기
      rl.write(historyLine);
    }
  } else if (key.name === 'down') {
    const historyLine = _navigateHistory('down');
    if (historyLine !== null) {
      _isNavigatingHistory = true;
      rl.write(null, { ctrl: true, name: 'u' });  // Ctrl+U로 라인 지우기
      if (historyLine === '') {
        // 빈 문자열이면 원래 입력 복원
        rl.write(_currentInput);
        _currentInput = '';
      } else {
        rl.write(historyLine);
      }
    }
  } else {
    // 다른 키 입력 시 history 탐색 모드 해제
    if (_isNavigatingHistory && key.name !== 'up' && key.name !== 'down') {
      _isNavigatingHistory = false;
      _currentInput = '';
    }
  }
});

function getIsNavigatingHistory() {
  return _isNavigatingHistory;
}

function setIsNavigatingHistory(value) {
  _isNavigatingHistory = value;
}

function getCurrentInput() {
  return _currentInput;
}

function setCurrentInput(value) {
  _currentInput = value;
}

module.exports = {
  rl,
  getIsNavigatingHistory,
  setIsNavigatingHistory,
  getCurrentInput,
  setCurrentInput
};

