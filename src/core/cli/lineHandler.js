const { parseLine } = require('parse-line');
const { knownCommands } = require('../command/knownCommands.js');
const { _do } = require('./router.js');
const { _prompt } = require('./prompt.js');
const { _checkMultilineComplete, _resetMultilineBuffer } = require('./multiline.js');
const { _addToHistory, _resetHistoryIndex, _saveHistory } = require('./history.js');
const { _mapEchoVariables } = require('../command/echo.js');
const { _query } = require('../chain/executeQuery.js');
const { getIsMultilineInput } = require('../state/globals.js');
const { rl, getIsNavigatingHistory, setIsNavigatingHistory, getCurrentInput, setCurrentInput } = require('./readline.js');

let lineHandler;

async function enableLineHandler() {
  if (lineHandler) {
    rl.off('line', lineHandler);
  }
  lineHandler = defaultLineHandler;
  rl.on('line', lineHandler);
}

function defaultLineHandler(_line) {
  if (_line == 'exit') {
    _saveHistory();  // 종료 시 history 저장
    rl.close();
    console.log('good bye');
    process.exit(0);
  }

  // history 탐색 모드 해제
  if (getIsNavigatingHistory()) {
    setIsNavigatingHistory(false);
    setCurrentInput('');
  }

  // 빈 입력은 무시 (단, 멀티라인 입력 중이면 허용)
  const trimmedLine = _line.trim();
  if (trimmedLine.length === 0 && !getIsMultilineInput()) {
    rl.setPrompt(_prompt());
    rl.prompt();
    return;
  }

  // 멀티라인 입력 중이면 명령어 체크를 건너뛰고 멀티라인 처리만 수행
  const isMultilineInput = getIsMultilineInput();
  
  if (isMultilineInput) {
    // 멀티라인 입력 중에는 명령어 체크를 하지 않고 멀티라인 처리만 수행
    const multilineResult = _checkMultilineComplete(_line);
    
    if (multilineResult.error) {
      rl.setPrompt(_prompt());
      rl.prompt();
      return;
    }
    
    if (!multilineResult.complete) {
      rl.setPrompt('... ');
      rl.prompt();
      return;
    }
    
    // 완성된 멀티라인 처리
    let completeLine = multilineResult.line;
    let originalMultiline = multilineResult.originalMultiline || completeLine;
    if (completeLine === null) {
      rl.setPrompt(_prompt());
      rl.prompt();
      return;
    }
    
    const semicolonIndex = completeLine.indexOf(';');
    if (semicolonIndex > -1) {
      completeLine = completeLine.substring(0, semicolonIndex).trim();
    }
    
    const isMultiline = originalMultiline.includes('\n');
    _addToHistory(isMultiline ? originalMultiline : (completeLine + ';'), isMultiline);
    _resetHistoryIndex();
    
    const completeTokens = parseLine(completeLine);
    _do(completeTokens, completeLine).then(() => {
      rl.setPrompt(_prompt());
      rl.prompt();
    });
    return;
  }

  const inputTokens = parseLine(_line);

  // curl, event 명령어는 멀티라인 입력 지원 (특별 처리)
  const isMultilineCommand = inputTokens.length > 0 && 
    (inputTokens[0] === 'curl' || inputTokens[0] === 'event');

  if (isMultilineCommand) {
    // curl 명령어는 멀티라인 처리
    const multilineResult = _checkMultilineComplete(_line);

    // 에러가 발생한 경우 (문자열이 닫히지 않음)
    if (multilineResult.error) {
      rl.setPrompt(_prompt());
      rl.prompt();
      return;
    }

    if (!multilineResult.complete) {
      // 아직 완성되지 않음 - 다음 라인 대기
      rl.setPrompt('... ');
      rl.prompt();
      return;
    }

    // 완성된 curl 명령어 처리
    let completeLine = multilineResult.line;
    let originalMultiline = multilineResult.originalMultiline || completeLine;
    if (completeLine === null) {
      rl.setPrompt(_prompt());
      rl.prompt();
      return;
    }

    // 세미콜론과 그 이후의 모든 내용 제거 (이미 _checkMultilineComplete에서 처리되었지만 안전을 위해)
    const semicolonIndex = completeLine.indexOf(';');
    if (semicolonIndex > -1) {
      completeLine = completeLine.substring(0, semicolonIndex).trim();
    }

    // History에 추가 (멀티라인인 경우 원본 형태로 저장, 세미콜론 포함)
    const isMultiline = originalMultiline.includes('\n');
    // originalMultiline에는 이미 세미콜론이 포함되어 있음
    _addToHistory(isMultiline ? originalMultiline : (completeLine + ';'), isMultiline);
    _resetHistoryIndex();

    const completeTokens = parseLine(completeLine);
    _do(completeTokens, completeLine).then(() => {
      rl.setPrompt(_prompt());
      rl.prompt();
    });
    return;
  }

  // 명령어가 있으면 즉시 처리 (멀티라인 버퍼 초기화)
  if (inputTokens.length > 0 && knownCommands.includes(inputTokens[0])) {
    _resetMultilineBuffer();

    // History에 추가 (명령어도 저장)
    _addToHistory(_line);
    _resetHistoryIndex();

    _do(inputTokens, _line).then(() => {
      rl.setPrompt(_prompt());
      rl._promptDisplayed = false; // 프롬프트 표시 플래그 리셋
      rl.prompt();
    });
    return;
  }

  // 명령어가 없으면 _query로 처리될 것 - 멀티라인 처리
  const multilineResult = _checkMultilineComplete(_line);

  // 에러가 발생한 경우 (문자열이 닫히지 않음)
  if (multilineResult.error) {
    rl.setPrompt(_prompt());
    rl.prompt();
    return;
  }

  if (!multilineResult.complete) {
    // 아직 완성되지 않음 - 다음 라인 대기
    rl.setPrompt('... ');
    rl.prompt();
    return;
  }

  // 완성된 라인 처리
  const completeLine = multilineResult.line;
  let originalMultiline = multilineResult.originalMultiline || completeLine;
  if (completeLine === null) {
    // 에러로 인한 완성 (실제로는 처리하지 않음)
    rl.setPrompt(_prompt());
    rl.prompt();
    return;
  }

  // History에 추가 (멀티라인 완성 시 원본 형태로 저장, 세미콜론 포함)
  const isMultiline = originalMultiline.includes('\n');
  // originalMultiline에는 이미 세미콜론이 포함되어 있음
  _addToHistory(isMultiline ? originalMultiline : (completeLine + ';'), isMultiline);
  _resetHistoryIndex();

  const completeTokens = parseLine(completeLine);

  _do(completeTokens, completeLine).then(() => {
    rl.setPrompt(_prompt());
    rl.prompt();
  });
}

module.exports = {
  enableLineHandler,
  defaultLineHandler
};

