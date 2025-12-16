const { _getParseInput } = require("./common/currentPath.js");
const { RpcRequest } = require('../../lib/rpc/rpcRequest.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require("../log/printScreen.js");

/**
 * varStore를 깊은 복사로 복제
 * @param {object} source - 복제할 객체
 * @returns {object} 복제된 객체
 */
function _deepCopyVarStore(source) {
  return JSON.parse(JSON.stringify(source));
}

/**
 * snapshot 생성
 * @param {string} snapshotName - snapshot 이름
 * @param {string} rpcUrl - RPC URL
 */
async function _createSnapshot(snapshotName, rpcUrl) {
  try {
    // ganache-cli 또는 anvil의 evm_snapshot RPC 메서드 호출
    const res = await RpcRequest(rpcUrl, "evm_snapshot", [], 1);
    
    if (res.error) {
      throw new Error(res.error.message || "Failed to create snapshot");
    }
    
    const snapshotId = res.result;
    const timestamp = new Date().toISOString();
    
    // global.snapshot에 저장
    if (!global.snapshot) {
      global.snapshot = {};
    }
    
    global.snapshot[snapshotName] = {
      name: snapshotName,
      snapshotId: snapshotId,
      timestamp: timestamp,
      varStore: _deepCopyVarStore(global.varStore || {})
    };
    
    printSuccess(`Snapshot "${snapshotName}" created with ID: ${snapshotId}`);
    return snapshotId;
  } catch (e) {
    throw new Error(`Error creating snapshot: ${e.message}`);
  }
}

/**
 * snapshot 복원
 * @param {string} snapshotNameOrId - snapshot 이름 또는 ID
 * @param {string} rpcUrl - RPC URL
 */
async function _revertSnapshot(snapshotNameOrId, rpcUrl) {
  try {
    let snapshotId = null;
    let snapshotName = null;
    
    // 이름으로 찾기 (우선순위)
    if (global.snapshot && global.snapshot[snapshotNameOrId]) {
      snapshotId = global.snapshot[snapshotNameOrId].snapshotId;
      snapshotName = snapshotNameOrId;
    } else {
      // ID로 찾기
      if (global.snapshot) {
        for (const [name, snapshot] of Object.entries(global.snapshot)) {
          if (snapshot.snapshotId === snapshotNameOrId) {
            snapshotId = snapshot.snapshotId;
            snapshotName = name;
            break;
          }
        }
      }
      
      // 찾지 못했으면 입력값을 ID로 사용
      if (!snapshotId) {
        snapshotId = snapshotNameOrId;
      }
    }
    
    // ganache-cli 또는 anvil의 evm_revert RPC 메서드 호출
    const res = await RpcRequest(rpcUrl, "evm_revert", [snapshotId], 1);
    
    if (res.error) {
      throw new Error(res.error.message || "Failed to revert snapshot");
    }
    
    if (!res.result) {
      throw new Error("Snapshot revert failed: Invalid snapshot ID");
    }
    
    // varStore 복원
    if (snapshotName && global.snapshot[snapshotName]) {
      global.varStore = _deepCopyVarStore(global.snapshot[snapshotName].varStore);
      printSuccess(`Snapshot "${snapshotName}" reverted and varStore restored`);
    } else {
      printSuccess(`Snapshot ID "${snapshotId}" reverted`);
    }
    
    return true;
  } catch (e) {
    throw new Error(`Error reverting snapshot: ${e.message}`);
  }
}

/**
 * snapshot 리스트 출력
 */
function _listSnapshots() {
  if (!global.snapshot || Object.keys(global.snapshot).length === 0) {
    printGray("No snapshots found");
    return;
  }
  
  console.log(textColor.blue + "Snapshot List:" + textColor.white);
  console.log(textColor.gray + "Time".padEnd(30) + "Name".padEnd(20) + "Snapshot ID" + textColor.white);
  console.log(textColor.gray + "-".repeat(80) + textColor.white);
  
  const snapshots = Object.values(global.snapshot).sort((a, b) => 
    new Date(b.timestamp) - new Date(a.timestamp)
  );
  
  snapshots.forEach(snapshot => {
    const time = new Date(snapshot.timestamp).toLocaleString();
    console.log(
      textColor.white + time.padEnd(30) + 
      textColor.blue + snapshot.name.padEnd(20) + 
      textColor.gray + snapshot.snapshotId + 
      textColor.white
    );
  });
}

async function _commandSnapshot(_inputTokens) {
  const parsedInput = _getParseInput(_inputTokens, 1);
  
  // --name 옵션: snapshot 생성
  if (parsedInput.opt['--name'] !== undefined) {
    const snapshotName = parsedInput.opt['--name'];
    if (!snapshotName || snapshotName.trim() === '') {
      throw new Error("Error: snapshot name is required");
    }
    
    await _createSnapshot(snapshotName.trim(), global.rpcUrl);
    return;
  }
  
  // --revert 옵션: snapshot 복원
  if (parsedInput.opt['--revert'] !== undefined) {
    const snapshotNameOrId = parsedInput.opt['--revert'];
    if (!snapshotNameOrId || snapshotNameOrId.trim() === '') {
      throw new Error("Error: snapshot name or ID is required");
    }
    
    await _revertSnapshot(snapshotNameOrId.trim(), global.rpcUrl);
    return;
  }
  
  // -list 옵션: snapshot 리스트
  if (parsedInput.opt['-list'] !== undefined) {
    _listSnapshots();
    return;
  }
  
  // 옵션이 없으면 에러
  throw new Error("Error: Please specify an option: --name, --revert, or -list");
}

module.exports._commandSnapshot = _commandSnapshot;
