// 멀티라인 파싱 유틸리티 함수
// 스크립트 파일이나 입력 라인에서 멀티라인 쿼리를 처리하는 공통 로직

/**
 * 한 줄을 분석하여 멀티라인 완성 여부를 확인
 * @param {string} _line - 분석할 라인
 * @param {object} _state - 현재 상태 (depth, inString, escape, buffer)
 * @param {boolean} _multilineEnabled - 멀티라인 기능 활성화 여부
 * @returns {object} { complete: boolean, line: string|null, error: boolean, state: object }
 */
function checkMultilineComplete(_line, _state, _multilineEnabled) {
  if (!_multilineEnabled) {
    return { complete: true, line: _line, error: false, state: _state };
  }

  let depth = {
    paren: _state.depth.paren || 0,
    bracket: _state.depth.bracket || 0,
    brace: _state.depth.brace || 0
  };
  let inString = {
    single: _state.inString.single || false,
    double: _state.inString.double || false
  };
  let escape = _state.escape || false;
  let buffer = _state.buffer || '';

  for (let i = 0; i < _line.length; i++) {
    const char = _line[i];

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
      // 문자열 내부는 depth 계산하지 않음
      continue;
    }

    // Depth 계산 (문자열 외부에서만)
    if (char === '(') depth.paren++;
    else if (char === ')') depth.paren--;
    else if (char === '[') depth.bracket++;
    else if (char === ']') depth.bracket--;
    else if (char === '{') depth.brace++;
    else if (char === '}') depth.brace--;
  }

  const trimmedLine = _line.trim();
  const hasSemicolon = trimmedLine.indexOf(';') !== -1;

  // 버퍼에 추가 (빈 줄도 처리)
  if (trimmedLine.length === 0) {
    // 빈 줄은 공백 하나로 처리
    if (buffer.length > 0) {
      buffer += ' ';
    }
  } else {
    // 세미콜론이 있으면 그 이후 내용은 무시
    let lineToAdd = trimmedLine;
    const semicolonIndex = trimmedLine.indexOf(';');
    if (semicolonIndex > -1) {
      lineToAdd = trimmedLine.substring(0, semicolonIndex).trim();
    }
    
    if (buffer.length > 0) {
      // curl 명령어의 경우, URL 경로가 슬래시로 시작하면 공백 없이 연결
      const isCurlCommand = buffer.trim().toLowerCase().startsWith('curl');
      if (isCurlCommand && lineToAdd.startsWith('/')) {
        buffer += lineToAdd;
      } else {
        // .then() 또는 .catch() 블록 내부의 경우 줄바꿈 보존 (JavaScript 코드이므로)
        // 버퍼에 .then( 또는 .catch(이 포함되어 있거나, 중괄호가 열려있으면 줄바꿈 보존
        const isThenOrCatchBlock = buffer.includes('.then(') || buffer.includes('.catch(') || depth.brace > 0;
        if (isThenOrCatchBlock) {
          buffer += '\n' + lineToAdd;
        } else {
          buffer += ' ' + lineToAdd;
        }
      }
    } else {
      buffer = lineToAdd;
    }
  }

  // 세미콜론이 있고 depth가 모두 0이면 완성
  const isComplete = hasSemicolon && depth.paren === 0 && depth.bracket === 0 && depth.brace === 0 && !inString.single && !inString.double;

  // 완성 시점에 문자열이 열린 상태인지 체크
  if (isComplete && (inString.single || inString.double)) {
    return {
      complete: true,
      line: null,
      error: true,
      state: {
        depth: { paren: 0, bracket: 0, brace: 0 },
        inString: { single: false, double: false },
        escape: false,
        buffer: ''
      }
    };
  }

  if (isComplete) {
    let completeLine = buffer.trim();
    // 세미콜론과 그 이후의 모든 내용 제거 (세미콜론은 멀티라인 종료 표시용)
    const semicolonIndex = completeLine.indexOf(';');
    if (semicolonIndex > -1) {
      completeLine = completeLine.substring(0, semicolonIndex).trim();
    }
    // 멀티라인으로 합쳐진 불필요한 공백 정리
    // 괄호와 점 사이의 공백 제거: ") ." -> ")."
    completeLine = completeLine.replace(/\s*\)\s*\./g, ').');
    // 점과 함수명 사이의 공백 제거: ". functionName" -> ".functionName"
    completeLine = completeLine.replace(/\.\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(/g, '.$1(');
    
    return {
      complete: true,
      line: completeLine,
      error: false,
      state: {
        depth: { paren: 0, bracket: 0, brace: 0 },
        inString: { single: false, double: false },
        escape: false,
        buffer: ''
      }
    };
  }

  return {
    complete: false,
    line: null,
    error: false,
    state: {
      depth: depth,
      inString: inString,
      escape: escape,
      buffer: buffer
    }
  };
}

/**
 * 스크립트 파일의 라인들을 멀티라인 쿼리로 합치기
 * @param {Array<string>} _lines - 스크립트 파일의 라인 배열
 * @param {boolean} _multilineEnabled - 멀티라인 기능 활성화 여부
 * @param {Array<string>} _knownCommands - 알려진 명령어 목록
 * @returns {Array<object>} [{ line: string, isCommand: boolean, originalLines: Array<number> }]
 */
function parseScriptLines(_lines, _multilineEnabled, _knownCommands) {
  const result = [];
  let state = {
    depth: { paren: 0, bracket: 0, brace: 0 },
    inString: { single: false, double: false },
    escape: false,
    buffer: ''
  };
  let currentLines = [];
  let isCommand = false;

  for (let i = 0; i < _lines.length; i++) {
    const line = _lines[i];
    const trimmedLine = line.trim();

    // 빈 줄이나 주석은 건너뛰기
    if (trimmedLine.length === 0 || trimmedLine.indexOf('#') === 0) {
      // 현재 버퍼가 있으면 계속 유지
      if (state.buffer.length > 0) {
        currentLines.push(i);
      }
      continue;
    }

    // 첫 번째 토큰 확인 (명령어인지 쿼리인지)
    const firstToken = trimmedLine.split(/\s+/)[0];
    // curl 명령어는 멀티라인을 지원하므로 쿼리처럼 처리
    const multilineCommands = ['curl'];
    const lineIsCommand = _knownCommands.includes(firstToken) && !multilineCommands.includes(firstToken);

    // 명령어가 시작되면 이전 버퍼를 먼저 처리
    if (lineIsCommand && state.buffer.length > 0) {
      // 이전 쿼리가 완성되지 않았으면 에러
      if (!state.buffer.trim().endsWith(';')) {
        throw new Error(`Error: Incomplete multiline query starting at line ${currentLines[0] + 1}`);
      }
      // 완성된 쿼리 추가
      let completeLine = state.buffer.trim();
      completeLine = completeLine.replace(/\s*\)\s*\./g, ').');
      completeLine = completeLine.replace(/\.\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(/g, '.$1(');
      result.push({
        line: completeLine,
        isCommand: false,
        originalLines: [...currentLines]
      });
      state = {
        depth: { paren: 0, bracket: 0, brace: 0 },
        inString: { single: false, double: false },
        escape: false,
        buffer: ''
      };
      currentLines = [];
    }

    if (lineIsCommand) {
      // 명령어는 즉시 추가
      result.push({
        line: trimmedLine,
        isCommand: true,
        originalLines: [i]
      });
      // 버퍼 초기화
      state = {
        depth: { paren: 0, bracket: 0, brace: 0 },
        inString: { single: false, double: false },
        escape: false,
        buffer: ''
      };
      currentLines = [];
    } else {
      // 쿼리 또는 멀티라인을 지원하는 명령어(curl 등)는 멀티라인 처리
      // curl 명령어는 명령어이지만 멀티라인을 지원하므로 쿼리처럼 처리
      const isMultilineCommand = multilineCommands.includes(firstToken);
      // 버퍼에 이미 멀티라인 명령어가 있는지 확인 (예: curl로 시작한 버퍼가 있으면 계속 처리)
      const hasMultilineCommandInBuffer = state.buffer.length > 0 && 
        multilineCommands.some(cmd => state.buffer.trim().toLowerCase().startsWith(cmd));
      const shouldMultiline = isMultilineCommand || hasMultilineCommandInBuffer;
      
      if (!_multilineEnabled && !shouldMultiline) {
        // 멀티라인 비활성화면 각 라인을 독립적으로 처리 (단, curl은 예외)
        result.push({
          line: trimmedLine,
          isCommand: isMultilineCommand, // curl은 명령어로 표시
          originalLines: [i]
        });
      } else {
        // 멀티라인 활성화면 버퍼에 추가
        currentLines.push(i);
        const checkResult = checkMultilineComplete(line, state, _multilineEnabled || shouldMultiline);
        
        if (checkResult.error) {
          throw new Error(`Error: Unclosed string at line ${i + 1}`);
        }
        
        state = checkResult.state;
        
        if (checkResult.complete) {
          // 완성된 쿼리 또는 멀티라인 명령어 추가
          result.push({
            line: checkResult.line,
            isCommand: shouldMultiline, // curl은 명령어로 표시
            originalLines: [...currentLines]
          });
          state = {
            depth: { paren: 0, bracket: 0, brace: 0 },
            inString: { single: false, double: false },
            escape: false,
            buffer: ''
          };
          currentLines = [];
        }
      }
    }
  }

  // 마지막에 버퍼가 남아있으면 처리
  if (state.buffer.length > 0) {
    if (!state.buffer.trim().endsWith(';')) {
      throw new Error(`Error: Incomplete multiline query starting at line ${currentLines[0] + 1}`);
    }
    let completeLine = state.buffer.trim();
    completeLine = completeLine.replace(/\s*\)\s*\./g, ').');
    completeLine = completeLine.replace(/\.\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(/g, '.$1(');
    result.push({
      line: completeLine,
      isCommand: false,
      originalLines: [...currentLines]
    });
  }

  return result;
}

module.exports.checkMultilineComplete = checkMultilineComplete;
module.exports.parseScriptLines = parseScriptLines;




