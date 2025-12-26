const fs = require('fs');
const path = require('path');
const { resolveWorkspacePath } = require('../../utils/resolvePath.js');

/**
 * Block 범위 결정
 * @param {string|undefined} block - 단일 block
 * @param {string|undefined} fromBlock - 시작 block
 * @param {string|undefined} toBlock - 종료 block
 * @returns {string[]} block 번호 배열
 */
function _determineBlockRange(block, fromBlock, toBlock) {
  if (block !== undefined) {
    return [block.toString()];
  }
  
  if (fromBlock !== undefined && toBlock !== undefined) {
    const blocks = [];
    const from = parseInt(fromBlock);
    const to = parseInt(toBlock);
    
    if (isNaN(from) || isNaN(to)) {
      throw new Error('Error: --fromBlock and --toBlock must be valid numbers');
    }
    
    if (from > to) {
      throw new Error('Error: --fromBlock must be less than or equal to --toBlock');
    }
    
    for (let i = from; i <= to; i++) {
      blocks.push(i.toString());
    }
    return blocks;
  }
  
  throw new Error('Error: --block or (--fromBlock and --toBlock) is required');
}

/**
 * Block 파일에서 eventLogs 추출 (flatten)
 * @param {object} blockData - Block JSON 데이터
 * @returns {Array} 이벤트 로그 배열
 */
function _flattenEventLogs(blockData) {
  const events = [];
  
  if (blockData && blockData.eventLogs && Array.isArray(blockData.eventLogs)) {
    events.push(...blockData.eventLogs);
  }
  
  return events;
}

/**
 * 필터 적용 (AND 조건)
 * @param {Array} events - 이벤트 배열
 * @param {object} filters - 필터 옵션
 * @returns {Array} 필터링된 이벤트 배열
 */
function _applyFilters(events, filters) {
  return events.filter(event => {
    // contractName 필터
    if (filters.contractName !== undefined) {
      if (!event.contractName || event.contractName !== filters.contractName) {
        return false;
      }
    }
    
    // contractAddress 필터
    if (filters.contractAddress !== undefined) {
      const eventAddr = (event.contractAddress || '').toLowerCase();
      const filterAddr = filters.contractAddress.toLowerCase();
      if (eventAddr !== filterAddr) {
        return false;
      }
    }
    
    // eventName 필터
    if (filters.eventName !== undefined) {
      if (!event.eventName || event.eventName !== filters.eventName) {
        return false;
      }
    }
    
    // txHash 필터
    if (filters.txHash !== undefined) {
      const eventTxHash = (event.transactionHash || '').toLowerCase();
      const filterTxHash = filters.txHash.toLowerCase();
      if (eventTxHash !== filterTxHash) {
        return false;
      }
    }
    
    // logIndex 필터
    if (filters.logIndex !== undefined) {
      const eventLogIndex = event.logIndex ? parseInt(event.logIndex) : undefined;
      const filterLogIndex = parseInt(filters.logIndex);
      if (isNaN(filterLogIndex) || eventLogIndex !== filterLogIndex) {
        return false;
      }
    }
    
    return true;
  });
}

/**
 * Event를 EventSummary 형식으로 변환
 * @param {Array} events - 이벤트 배열
 * @returns {Array} EventSummary 배열
 */
function _convertToEventSummary(events) {
  return events.map(event => {
    // decodeLog에서 method 제외하고 decoded로 변환
    const decoded = {};
    if (event.decodeLog) {
      Object.keys(event.decodeLog).forEach(key => {
        if (key !== 'method') {
          decoded[key] = event.decodeLog[key];
        }
      });
    }
    
    return {
      blockNumber: event.blockNumber?.toString() || '',
      blockHash: event.blockHash || '',
      transactionHash: event.transactionHash || '',
      transactionIndex: event.transactionIndex?.toString() || '',
      logIndex: event.logIndex?.toString() || '',
      contractName: event.contractName || '',
      contractAddress: event.contractAddress || '',
      eventName: event.eventName || '',
      decoded: decoded
    };
  });
}

/**
 * 파일에서 이벤트 로드
 * @param {string} loadSource - 디렉토리 경로
 * @param {string|undefined} block - 단일 block
 * @param {string|undefined} fromBlock - 시작 block
 * @param {string|undefined} toBlock - 종료 block
 * @returns {Promise<Array>} 이벤트 배열
 */
async function _loadFromFile(loadSource, block, fromBlock, toBlock) {
  const resolvedPath = resolveWorkspacePath(loadSource);
  
  // 디렉토리 존재 확인
  if (!fs.existsSync(resolvedPath)) {
    return [];
  }
  
  if (!fs.statSync(resolvedPath).isDirectory()) {
    throw new Error(`Error: ${resolvedPath} is not a directory`);
  }
  
  // Block 범위 결정
  const blocks = _determineBlockRange(block, fromBlock, toBlock);
  
  // 각 block 파일 로드
  const allEvents = [];
  for (const blockNum of blocks) {
    const filePath = path.join(resolvedPath, `${blockNum}.json`);
    
    // 파일이 없으면 무시 (에러 아님)
    if (!fs.existsSync(filePath)) {
      continue;
    }
    
    // .json 파일만 처리
    if (!filePath.toLowerCase().endsWith('.json')) {
      continue;
    }
    
    try {
      const fileContent = fs.readFileSync(filePath, 'utf8');
      const blockData = JSON.parse(fileContent);
      
      // eventLogs flatten
      const events = _flattenEventLogs(blockData);
      allEvents.push(...events);
    } catch (e) {
      // JSON 파싱 실패 시 무시 (에러 아님)
      continue;
    }
  }
  
  return allEvents;
}

/**
 * 이벤트 로드 및 필터링
 * @param {object} options - 옵션 객체
 * @returns {Promise<Array>} 필터링된 EventSummary 배열
 */
async function _loadAndFilterEvents(options) {
  const { loadSource, block, fromBlock, toBlock, ...filters } = options;
  
  // 1. 파일에서 이벤트 로드
  const eventData = await _loadFromFile(loadSource, block, fromBlock, toBlock);
  
  // 2. 필터 적용 (AND 조건)
  const filteredEvents = _applyFilters(eventData, filters);
  
  // 3. Event Summary 변환
  return _convertToEventSummary(filteredEvents);
}

module.exports = {
  _loadAndFilterEvents,
  _determineBlockRange,
  _flattenEventLogs,
  _applyFilters,
  _convertToEventSummary
};

