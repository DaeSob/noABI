const https = require('https');
const http = require('http');
const { URL } = require('url');
const { _getParseInput } = require("./common/currentPath.js");
const { _parseKeyValue, reservedVarName } = require("./common/keyValue.js");
const { textColor, printError, printSuccess, printGray, printDefault } = require("../log/printScreen.js");

/**
 * curl 명령어 처리
 * @param {Array<string>} _inputTokens - 입력 토큰 배열
 * @param {string} _line - 원본 라인
 */
async function _commandCurl(_inputTokens, _line) {
  // => 연산자로 변수명 추출 (예: curl GET https://api.example.com => result)
  let storeVarName = null;
  let commandLine = _line || '';
  
  if (_line) {
    const arrowMatch = _line.match(/=>\s*(\w+)/);
    if (arrowMatch && arrowMatch[1]) {
      storeVarName = arrowMatch[1].trim();
      // => 이전 부분만 추출
      const arrowIndex = _line.indexOf('=>');
      if (arrowIndex > -1) {
        commandLine = _line.substring(0, arrowIndex).trim();
      }
    }
  }

  // 메서드 추출 (기본값: GET)
  let method = 'GET';
  const methodMatch = commandLine.match(/curl\s+(GET|POST|PUT|DELETE|PATCH)/i);
  if (methodMatch) {
    method = methodMatch[1].toUpperCase();
  }

  // URL 추출 (commandLine에서 직접 파싱)
  let url = '';
  // curl [METHOD] URL 형식에서 URL 추출
  // 먼저 메서드와 URL 부분만 추출
  let urlPart = commandLine.replace(/^curl\s+/i, '').trim();
  // 메서드 제거
  urlPart = urlPart.replace(/^(GET|POST|PUT|DELETE|PATCH)\s+/i, '');
  
  // URL은 첫 번째 공백 또는 옵션(-H, -d) 전까지
  const urlEndMatch = urlPart.match(/\s+(-[Hd]|=>)/);
  if (urlEndMatch) {
    url = urlPart.substring(0, urlEndMatch.index).trim();
  } else {
    // 옵션이 없으면 전체가 URL (=> 전까지)
    const arrowIndex = urlPart.indexOf('=>');
    if (arrowIndex > -1) {
      url = urlPart.substring(0, arrowIndex).trim();
    } else {
      url = urlPart.trim();
    }
  }

  if (!url || url === '') {
    throw new Error("Error: URL is required");
  }

  // 세미콜론 제거 (멀티라인 종료 표시용)
  url = url.replace(/;\s*$/, '').trim();

  // 따옴표 제거
  if ((url.startsWith('"') && url.endsWith('"')) ||
      (url.startsWith("'") && url.endsWith("'"))) {
    url = url.slice(1, -1);
  }

  // 변수 파싱 (${변수명} 형식)
  try {
    url = _parseKeyValue(url);
    if (typeof url === 'object') {
      url = JSON.stringify(url);
    } else {
      url = String(url);
    }
  } catch (e) {
    throw new Error(e.toString());
  }

  // URL 파싱
  let urlObj;
  try {
    urlObj = new URL(url);
  } catch (e) {
    throw new Error(`Error: Invalid URL: ${url}`);
  }

  // 헤더 파싱 (-H 옵션, 여러 개 지원)
  const headers = {};
  const headerRegex = /-H\s+(["']?)([^"']+?)\1/g;
  let headerMatch;
  while ((headerMatch = headerRegex.exec(commandLine)) !== null) {
    const headerStr = headerMatch[2].trim();
    const colonIndex = headerStr.indexOf(':');
    if (colonIndex > -1) {
      const headerName = headerStr.substring(0, colonIndex).trim();
      const headerValue = headerStr.substring(colonIndex + 1).trim();
      headers[headerName] = headerValue;
    }
  }

  // 데이터 파싱 (-d 옵션, 따옴표 처리)
  let postData = null;
  const dataRegex = /-d\s+(["'])([\s\S]*?)\1|-d\s+([^\s]+)/;
  const dataMatch = commandLine.match(dataRegex);
  if (dataMatch) {
    // 따옴표로 감싸진 경우
    if (dataMatch[2] !== undefined) {
      postData = dataMatch[2];
    } 
    // 따옴표 없이 공백으로 구분된 경우
    else if (dataMatch[3] !== undefined) {
      postData = dataMatch[3];
    }
    
    if (postData) {
      // 변수 파싱
      try {
        const parsedData = _parseKeyValue(postData);
        if (typeof parsedData === 'object') {
          postData = JSON.stringify(parsedData);
        } else {
          postData = String(parsedData);
        }
      } catch (e) {
        // 변수가 없으면 원본 사용
      }
      
      // Content-Type이 없으면 기본값 설정
      if (!headers['Content-Type'] && !headers['content-type']) {
        headers['Content-Type'] = 'application/json';
      }
    }
  }

  // HTTP 요청 실행
  const result = await _executeHttpRequest(urlObj, method, headers, postData);
  
  // => 연산자가 있으면 store에 저장
  if (storeVarName) {
    // 예약 변수명 체크
    if (reservedVarName.includes(storeVarName)) {
      throw new Error(`Error: "${storeVarName}" is a reserved variable name`);
    }
    // varStore에 저장 (JSON이면 객체로, 아니면 문자열로 저장)
    global.varStore[storeVarName] = result;
  }

  // 결과 출력 (객체면 JSON.stringify, 문자열이면 그대로)
  const output = typeof result === 'object' ? JSON.stringify(result, null, 2) : result;
  console.log(textColor.success + output + textColor.white);
}

/**
 * HTTP 요청 실행
 * @param {URL} urlObj - URL 객체
 * @param {string} method - HTTP 메서드
 * @param {object} headers - HTTP 헤더
 * @param {string|null} postData - POST 데이터
 * @returns {Promise<string|object>} 응답 본문 (JSON이면 객체, 아니면 문자열)
 */
function _executeHttpRequest(urlObj, method, headers, postData) {
  return new Promise((resolve, reject) => {
    const isHttps = urlObj.protocol === 'https:';
    const httpModule = isHttps ? https : http;

    const options = {
      hostname: urlObj.hostname,
      port: urlObj.port || (isHttps ? 443 : 80),
      path: urlObj.pathname + urlObj.search,
      method: method,
      headers: headers
    };

    const req = httpModule.request(options, (res) => {
      let data = '';

      res.on('data', (chunk) => {
        data += chunk;
      });

      res.on('end', () => {
        // HTTP 상태 코드 체크
        if (res.statusCode >= 200 && res.statusCode < 300) {
          // JSON 파싱 시도
          let parsedResult = data;
          if (data && data.trim().length > 0) {
            try {
              // JSON으로 파싱 가능한지 시도
              const trimmed = data.trim();
              if ((trimmed.startsWith('{') && trimmed.endsWith('}')) || 
                  (trimmed.startsWith('[') && trimmed.endsWith(']'))) {
                parsedResult = JSON.parse(data);
              }
            } catch (e) {
              // JSON 파싱 실패하면 문자열 그대로 사용
              parsedResult = data;
            }
          }
          resolve(parsedResult);
        } else {
          reject(new Error(`HTTP Error ${res.statusCode}: ${data}`));
        }
      });
    });

    req.on('error', (error) => {
      reject(new Error(`Request failed: ${error.message}`));
    });

    // POST 데이터 전송
    if (postData && (method === 'POST' || method === 'PUT' || method === 'PATCH')) {
      req.write(postData);
    }

    req.end();
  });
}

module.exports._commandCurl = _commandCurl;

