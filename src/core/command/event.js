const { _getParseInput } = require('./common/currentPath.js');
const { _extractStoreVarName, _saveToVarStore, _parseKeyValue } = require('./common/keyValue.js');
const { _loadAndFilterEvents } = require('./common/eventParser.js');
const { printError, printSuccess, printDefault } = require('../log/printScreen.js');

/**
 * event 명령어 처리
 * @param {Array<string>} _inputTokens - 입력 토큰 배열
 * @param {string} _line - 원본 라인 (멀티라인 포함)
 */
async function _commandEvent(_inputTokens, _line) {
  try {
    // 1. => 연산자로 varStoreKey 추출
    const { storeVarName, commandLine } = _extractStoreVarName(_line);
    
    // 2. 옵션 파싱
    const parsedInput = _getParseInput(_inputTokens, 1);
    
    // 변수 파싱 헬퍼 함수
    const parseOptionValue = (value) => {
      if (value === undefined) return undefined;
      try {
        const parsed = _parseKeyValue(value);
        return typeof parsed === 'object' ? JSON.stringify(parsed) : String(parsed);
      } catch (e) {
        throw new Error(`Error parsing option value "${value}": ${e.message}`);
      }
    };
    
    // 3. 필수 옵션 체크 및 변수 파싱
    // -load는 플래그이고, 경로는 data[0]에서 가져옴
    if (!parsedInput.opt['-load']) {
      throw new Error('Error: -load option is required');
    }
    
    if (!parsedInput.data || parsedInput.data.length === 0) {
      throw new Error('Error: Event source path is required');
    }
    
    const loadSource = parseOptionValue(parsedInput.data[0]);
    if (!loadSource) {
      throw new Error('Error: Event source path is required');
    }
    
    // 4. Block 범위 옵션 체크 및 변수 파싱
    const block = parseOptionValue(parsedInput.opt['--block']);
    const fromBlock = parseOptionValue(parsedInput.opt['--fromBlock']);
    const toBlock = parseOptionValue(parsedInput.opt['--toBlock']);
    
    // --block과 --fromBlock/--toBlock 동시 사용 불가
    if (block !== undefined && (fromBlock !== undefined || toBlock !== undefined)) {
      throw new Error('Error: --block and (--fromBlock/--toBlock) cannot be used together');
    }
    
    // 5. 필터 옵션 및 변수 파싱
    const contractName = parseOptionValue(parsedInput.opt['--contractName']);
    const contractAddress = parseOptionValue(parsedInput.opt['--contractAddress']);
    const eventName = parseOptionValue(parsedInput.opt['--eventName']);
    const txHash = parseOptionValue(parsedInput.opt['--txHash']);
    const logIndex = parseOptionValue(parsedInput.opt['--logIndex']);
    
    // 6. 이벤트 로드 및 필터링
    const events = await _loadAndFilterEvents({
      loadSource,
      block,
      fromBlock,
      toBlock,
      contractName,
      contractAddress,
      eventName,
      txHash,
      logIndex
    });
    
    // 7. varStore에 저장 (항상 배열)
    // events가 undefined이거나 null이면 빈 배열로 저장
    const resultArray = Array.isArray(events) ? events : [];
    
    if (storeVarName) {
      _saveToVarStore(storeVarName, resultArray);
      printSuccess(`Event loaded: ${resultArray.length} event(s) saved to varStore["${storeVarName}"]`);
    } else {
      // varStoreKey가 없으면 CLI에 출력
      if (resultArray.length === 0) {
        printDefault('[]');
      } else {
        console.log(JSON.stringify(resultArray, null, 2));
      }
    }
  } catch (e) {
    // 에러 발생 시 빈 배열 저장 (DSL 중단 안 함)
    const { storeVarName } = _extractStoreVarName(_line);
    if (storeVarName) {
      _saveToVarStore(storeVarName, []);
      printError(`Warning: ${e.message} (empty array saved to varStore["${storeVarName}"])`);
    } else {
      printError(`Error: ${e.message}`);
      console.log('[]');
    }
  }
}

module.exports = {
  _commandEvent
};

