const { printError } = require('../log/printScreen.js');
const {
  getObjSettings,
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
  setIsMultilineInput
} = require('../state/globals.js');

// 멀티라인 depth 추적 함수 (query 전용)
function _checkMultilineComplete(_line) {
  const objSettings = getObjSettings();
  if (!objSettings.options.multiline) {
    return { complete: true, line: _line };
  }

  let depth = { ...getMultilineDepth() };
  let inString = { ...getMultilineInString() };
  let escape = getMultilineEscape();

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

  // 상태 업데이트
  setMultilineDepth(depth);
  setMultilineInString(inString);
  setMultilineEscape(escape);

  const trimmedLine = _line.trim();
  const hasSemicolon = trimmedLine.indexOf(';') !== -1;

  // 멀티라인 입력 시작 플래그 설정
  const isMultilineInput = getIsMultilineInput();
  if (!isMultilineInput && (depth.paren > 0 || depth.bracket > 0 || depth.brace > 0 || !hasSemicolon)) {
    setIsMultilineInput(true);
  }

  // 멀티라인 입력 중인 경우 원본 라인 저장 (빈 줄 제외)
  if (getIsMultilineInput() && trimmedLine.length > 0) {
    const originalLines = getMultilineOriginalLines();
    if (originalLines.length === 0) {
      // 첫 라인 저장
      setMultilineOriginalLines([_line]);
    } else {
      // 이후 라인 저장
      setMultilineOriginalLines([...originalLines, _line]);
    }
  }

  // 버퍼에 추가 (빈 줄도 처리)
  let multilineBuffer = getMultilineBuffer();
  if (trimmedLine.length === 0) {
    // 빈 줄은 공백 하나로 처리
    if (multilineBuffer.length > 0) {
      setMultilineBuffer(multilineBuffer + ' ');
    }
  } else {
    // 세미콜론이 있으면 그 이후 내용은 무시
    let lineToAdd = trimmedLine;
    const semicolonIndex = trimmedLine.indexOf(';');
    if (semicolonIndex > -1) {
      lineToAdd = trimmedLine.substring(0, semicolonIndex).trim();
    }

    if (multilineBuffer.length > 0) {
      // curl 명령어의 경우, URL 경로가 슬래시로 시작하면 공백 없이 연결
      const isCurlCommand = multilineBuffer.trim().toLowerCase().startsWith('curl');
      if (isCurlCommand && lineToAdd.startsWith('/')) {
        setMultilineBuffer(multilineBuffer + lineToAdd);
      } else if (isCurlCommand) {
        // curl 명령어의 경우, JSON 데이터나 헤더는 줄바꿈 보존
        // -d 옵션이 있거나 -H 옵션이 있으면 줄바꿈 보존
        const hasDataOption = multilineBuffer.includes('-d') || lineToAdd.includes('-d');
        const hasHeaderOption = multilineBuffer.includes('-H') || lineToAdd.includes('-H');
        if (hasDataOption || hasHeaderOption || inString.single || inString.double) {
          // JSON 데이터나 헤더는 줄바꿈 보존
          setMultilineBuffer(multilineBuffer + '\n' + lineToAdd);
        } else {
          setMultilineBuffer(multilineBuffer + ' ' + lineToAdd);
        }
      } else {
        // .then() 블록 내부의 경우 줄바꿈 보존 (JavaScript 코드이므로)
        // 버퍼에 .then(이 포함되어 있거나, 중괄호가 열려있으면 줄바꿈 보존
        const isThenBlock = multilineBuffer.includes('.then(') || depth.brace > 0;
        if (isThenBlock) {
          setMultilineBuffer(multilineBuffer + '\n' + lineToAdd);
        } else {
          setMultilineBuffer(multilineBuffer + ' ' + lineToAdd);
        }
      }
    } else {
      setMultilineBuffer(lineToAdd);
    }
  }

  // 세미콜론이 있고 depth가 모두 0이면 완성
  const isComplete = hasSemicolon && depth.paren === 0 && depth.bracket === 0 && depth.brace === 0 && !inString.single && !inString.double;

  // 완성 시점에 문자열이 열린 상태인지 체크
  if (isComplete && (inString.single || inString.double)) {
    printError('Error: Unclosed string. String must be closed before semicolon.');
    _resetMultilineBuffer();
    return { complete: true, line: null, error: true };
  }

  if (isComplete) {
    multilineBuffer = getMultilineBuffer();
    let completeLine = multilineBuffer.trim();
    // 세미콜론과 그 이후의 모든 내용 제거 (세미콜론은 멀티라인 종료 표시용)
    const semicolonIndex = completeLine.indexOf(';');
    if (semicolonIndex > -1) {
      completeLine = completeLine.substring(0, semicolonIndex).trim();
    }
    // curl 명령어의 경우 줄바꿈 보존 (JSON 데이터나 헤더를 위해)
    const isCurlCommand = completeLine.trim().toLowerCase().startsWith('curl');
    if (!isCurlCommand) {
      // 멀티라인으로 합쳐진 불필요한 공백 정리 (curl이 아닌 경우만)
      // 괄호와 점 사이의 공백 제거: ") ." -> ")."
      completeLine = completeLine.replace(/\s*\)\s*\./g, ').');
      // 점과 함수명 사이의 공백 제거: ". functionName" -> ".functionName"
      completeLine = completeLine.replace(/\.\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(/g, '.$1(');
    }

    // 멀티라인 원본도 준비 (history용) - 세미콜론 포함
    const originalLines = getMultilineOriginalLines();
    let originalMultiline = '';
    if (originalLines.length > 1) {
      // 여러 라인으로 입력된 경우 원본 형태로 저장 (세미콜론 포함)
      originalMultiline = originalLines.join('\n');
    } else if (originalLines.length === 1) {
      // 단일 라인이지만 멀티라인으로 처리된 경우 (세미콜론 포함)
      originalMultiline = originalLines[0];
    } else {
      // 원본이 없는 경우 completeLine에 세미콜론 추가
      originalMultiline = completeLine + ';';
    }

    // 버퍼 초기화
    _resetMultilineBuffer();
    return { complete: true, line: completeLine, originalMultiline: originalMultiline, error: false };
  }

  return { complete: false, line: null, error: false };
}

// 멀티라인 버퍼 초기화 함수
function _resetMultilineBuffer() {
  setMultilineBuffer('');
  setMultilineOriginalLines([]);
  setMultilineDepth({ paren: 0, bracket: 0, brace: 0 });
  setMultilineInString({ single: false, double: false });
  setMultilineEscape(false);
  setIsMultilineInput(false);
}

module.exports = {
  _checkMultilineComplete,
  _resetMultilineBuffer
};

