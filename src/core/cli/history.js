const fs = require('fs');
const path = require('path');
const {
  getWorkPath,
  getCustomHistory,
  setCustomHistory,
  getHistoryIndex,
  setHistoryIndex,
  getHistoryMaxSize,
  getHistoryFile,
  setHistoryFile
} = require('../state/globals.js');

// 커스텀 History 관리 함수
function _loadHistory() {
  try {
    let historyFile = getHistoryFile();
    if (historyFile === undefined) {
      const workPath = getWorkPath();
      const historyPath = path.join(workPath, '.cli_history');
      setHistoryFile(historyPath);
      historyFile = historyPath;
    }

    if (fs.existsSync(historyFile)) {
      const historyData = fs.readFileSync(historyFile, 'utf8');
      // 멀티라인 구분자로 분리 (각 history 엔트리는 \n---HISTORY_ENTRY---\n로 구분)
      const entries = historyData.split('\n---HISTORY_ENTRY---\n').filter(entry => entry.trim().length > 0);
      setCustomHistory(entries.map(entry => entry.trim()));
      setHistoryIndex(getCustomHistory().length);
    }
  } catch (error) {
    // history 로드 실패 시 무시
    setCustomHistory([]);
    setHistoryIndex(-1);
  }
}

function _saveHistory() {
  try {
    let historyFile = getHistoryFile();
    if (historyFile === undefined) {
      const workPath = getWorkPath();
      const historyPath = path.join(workPath, '.cli_history');
      setHistoryFile(historyPath);
      historyFile = historyPath;
    }

    const customHistory = getCustomHistory();
    const historyMaxSize = getHistoryMaxSize();
    // 최근 _historyMaxSize 개만 저장
    const historyToSave = customHistory.slice(-historyMaxSize);
    const historyData = historyToSave.join('\n---HISTORY_ENTRY---\n');
    fs.writeFileSync(historyFile, historyData, 'utf8');
  } catch (error) {
    // history 저장 실패 시 무시
  }
}

function _addToHistory(line, isMultiline = false) {
  if (!line || line.trim().length === 0) {
    return;
  }

  const customHistory = getCustomHistory();
  // 중복 제거 (바로 이전과 동일하면 추가하지 않음)
  if (customHistory.length > 0 && customHistory[customHistory.length - 1] === line) {
    return;
  }

  const newHistory = [...customHistory, line];
  setCustomHistory(newHistory);
  setHistoryIndex(newHistory.length);

  const historyMaxSize = getHistoryMaxSize();
  // 최대 크기 제한
  if (newHistory.length > historyMaxSize) {
    newHistory.shift();
    setCustomHistory(newHistory);
    setHistoryIndex(newHistory.length);
  }

  // 주기적으로 저장 (매 10개마다)
  if (newHistory.length % 10 === 0) {
    _saveHistory();
  }
}

function _resetHistoryIndex() {
  const customHistory = getCustomHistory();
  setHistoryIndex(customHistory.length);
}

// History에서 멀티라인을 한 줄로 가져오는 함수
function _getHistoryLine(index) {
  const customHistory = getCustomHistory();
  if (index < 0 || index >= customHistory.length) {
    return null;
  }

  const historyEntry = customHistory[index];
  // 멀티라인인 경우 줄바꿈을 공백으로 변환하여 한 줄로 표현
  if (historyEntry.includes('\n')) {
    // 줄바꿈을 공백으로 변환하고, 연속된 공백을 하나로 정리
    return historyEntry.replace(/\n/g, ' ').replace(/\s+/g, ' ').trim();
  }
  return historyEntry;
}

// History 탐색 함수 (상하 화살표용)
function _navigateHistory(direction) {
  const customHistory = getCustomHistory();
  if (customHistory.length === 0) {
    return null;
  }

  let historyIndex = getHistoryIndex();
  if (direction === 'up') {
    // 위로 (이전 history)
    if (historyIndex > 0) {
      historyIndex--;
      setHistoryIndex(historyIndex);
    }
    return _getHistoryLine(historyIndex);
  } else if (direction === 'down') {
    // 아래로 (다음 history)
    if (historyIndex < customHistory.length - 1) {
      historyIndex++;
      setHistoryIndex(historyIndex);
      return _getHistoryLine(historyIndex);
    } else {
      setHistoryIndex(customHistory.length);
      return '';  // 빈 문자열 (현재 입력으로 돌아감)
    }
  }

  return null;
}

module.exports = {
  _loadHistory,
  _saveHistory,
  _addToHistory,
  _resetHistoryIndex,
  _getHistoryLine,
  _navigateHistory
};

