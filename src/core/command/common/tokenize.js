
function _tokenizePathToObject(_path, _object) {
  let tempObject = _object;
  if (_path === undefined || _path.length === 0) {
  } else {
    // 끝에 점이 있는지 확인 (자동완성용)
    let endsWithDot = _path.endsWith('.');
    let cleanPath = endsWithDot ? _path.slice(0, -1) : _path;
    let arrayPath = cleanPath.split(".").filter(p => p.length > 0);
    
    // collection.contractname 형식인지 확인 (2개 토큰이고, contracts가 없는 경우)
    if (arrayPath.length === 2 && arrayPath[1] !== "contracts") {
      // collection.contractname -> deployedContract[collection].contracts[contractname].abi
      if (tempObject[arrayPath[0]] === undefined) {
        tempObject = undefined;
        return;
      }
      tempObject = tempObject[arrayPath[0]];
      
      if (tempObject.contracts === undefined) {
        tempObject = undefined;
        return;
      }
      tempObject = tempObject.contracts;
      
      if (tempObject[arrayPath[1]] === undefined) {
        tempObject = undefined;
        return;
      }
      tempObject = tempObject[arrayPath[1]];
      
      // abi 배열 반환
      if (tempObject.abi !== undefined) {
        tempObject = tempObject.abi;
      }
    } else if (arrayPath.length === 1 && endsWithDot) {
      // collection. 입력 시 -> deployedContract[collection].contracts 반환 (자동완성용)
      if (tempObject[arrayPath[0]] === undefined) {
        tempObject = undefined;
        return;
      }
      tempObject = tempObject[arrayPath[0]];
      
      if (tempObject.contracts !== undefined) {
        tempObject = tempObject.contracts;
      }
    } else if (arrayPath.length === 2 && arrayPath[1] !== "contracts" && endsWithDot) {
      // collection.contractname. 입력 시 -> deployedContract[collection].contracts[contractname].abi 반환 (자동완성용)
      if (tempObject[arrayPath[0]] === undefined) {
        tempObject = undefined;
        return;
      }
      tempObject = tempObject[arrayPath[0]];
      
      if (tempObject.contracts === undefined) {
        tempObject = undefined;
        return;
      }
      tempObject = tempObject.contracts;
      
      if (tempObject[arrayPath[1]] === undefined) {
        tempObject = undefined;
        return;
      }
      tempObject = tempObject[arrayPath[1]];
      
      // abi 배열 반환
      if (tempObject.abi !== undefined) {
        tempObject = tempObject.abi;
      }
    } else {
      // 기존 로직: collection.contracts.contractname.abi 형식 등
      for (let i = 0; i < arrayPath.length; i++) {
        if (tempObject[arrayPath[i]] === undefined) {
          tempObject = undefined;
          return;
        }
        tempObject = tempObject[arrayPath[i]];
      }
      
      // collection만 입력한 경우 (자동완성용): contracts 반환
      // 단, 이미 contracts를 거친 경우는 제외
      if (arrayPath.length === 1 && !endsWithDot && tempObject !== undefined) {
        if (tempObject.contracts !== undefined && typeof tempObject.contracts === 'object') {
          tempObject = tempObject.contracts;
        }
      }
    }
  }
  return tempObject;
}

module.exports._tokenizePathToObject = _tokenizePathToObject;