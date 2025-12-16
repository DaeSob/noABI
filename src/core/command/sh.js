const fs = require('fs');
const { sleep } = require('../../lib/utils/sleep.js');
const { _getParseInput } = require('./common/currentPath.js');
const { _commandSet } = require('./set.js');
const { _commandUse } = require('./use.js');
const { _commandCurl } = require('./curl.js');
const { _commandCompose } = require('./compose.js');
const { _commandDeploy } = require('./deploy.js');
const { _commandImport } = require('./import.js');
const { _commandExport } = require('./export.js');
const { _commandEcho, _mapEchoVariables } = require('./echo.js');
const { _commandSnapshot } = require('./snapshot.js');
const { _commandToUTC } = require('./time.js');
const { _commandGetBalance } = require('./getBalance.js');
const { _query } = require('../chain/executeQuery.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');
const { parseScriptLines } = require('../utils/multilineParser.js');
const { resolveWorkspacePath } = require('../utils/resolvePath.js');
const { knownCommands } = require('./knownCommands.js');

function _prompt(_script, _inputTokenized, _line) {
  prompt = textColor.pink + currentSet.wallet + textColor.white + '> ';

  if (_inputTokenized[0] === 'echo') {
    // echo 명령어의 경우 변수를 매핑하여 프롬프트에 표시
    const echoValue = _mapEchoVariables(_inputTokenized, _line);
    if (echoValue !== '') {
      prompt = prompt + textColor.gray + echoValue + textColor.white;
    }
  } else {
    prompt = prompt + textColor.gray + _script + textColor.white;
  }
  console.log(prompt);
}

async function _commandSh(_inputTokens) {
  const parsedInput = _getParseInput(_inputTokens, 1);
  const scriptFile = resolveWorkspacePath(parsedInput.data[0]);

  if (parsedInput.opt['--repeat'] !== undefined) {
    parsedInput.opt['--repeat'] = Number(parsedInput.opt['--repeat']);
  }

  const rawScripts = fs.readFileSync(scriptFile, 'utf8');  
  // Windows와 Unix 줄바꿈 모두 지원
  let arrayScript = rawScripts.split(/\r?\n/);
  
  // 멀티라인 설정 확인
  const multilineEnabled = global.objSettings && global.objSettings.options && global.objSettings.options.multiline !== false;
  
  // 스크립트 라인들을 멀티라인 쿼리로 파싱
  let parsedLines = [];
  try {
    parsedLines = parseScriptLines(arrayScript, multilineEnabled, knownCommands);
  } catch (e) {
    throw e;
  }
  
  while ( true ) {
    for( let i=0; i<parsedLines.length; i++ ){
      const parsedLine = parsedLines[i];
      const line = parsedLine.line;
      const isCommand = parsedLine.isCommand;
      
      // 원본 라인 정보로 프롬프트 표시
      // 멀티라인으로 합쳐진 경우 모든 원본 라인 표시
      let originalLine = '';
      if (parsedLine.originalLines.length > 1) {
        // 멀티라인: 모든 원본 라인을 개행으로 연결
        originalLine = parsedLine.originalLines
          .map(idx => arrayScript[idx])
          .filter(l => l !== undefined)
          .join('\n');
      } else {
        // 단일 라인: 첫 번째 원본 라인 사용
        const originalLineIndex = parsedLine.originalLines[0];
        originalLine = arrayScript[originalLineIndex] || line;
      }
      
      let inputTokenized = line.trimStart().split(/\s+/);
      _prompt( originalLine, inputTokenized, line );

      if (isCommand) {
        // 명령어 처리
        switch( inputTokenized[0] ) {
          case 'set': 
          {
            await _commandSet( inputTokenized, line );
          }
          break;
          case 'use': 
          {
            await _commandUse( inputTokenized );
          }
          break;    
          case "snapshot": 
          {
            await _commandSnapshot(_inputTokens);
          }  
          break;          
          case 'compose': 
          {
            await _commandCompose(inputTokenized);
          }
          break;
          case 'deploy': 
          {
            await _commandDeploy(inputTokenized);
          }
          break;
          case 'import':
          {
            await _commandImport(inputTokenized);
          }
          break;
          case 'export':
          {
            await _commandExport(inputTokenized);
          }
          break;       
          case "toUtc": {
            await _commandToUTC(_inputTokens, line);
          }
          break;       
          case "getBalance": {
            await _commandGetBalance(_inputTokens, line);
          }
          break;                  
          case "curl":
          {
            await _commandCurl(_inputTokens, line);
          }
          break;
          case 'echo':
          {
            // echo 명령어는 프롬프트에만 표시하고 출력하지 않음
            await _commandEcho(inputTokenized, line);
          }
          break;
          case 'sleep':
          {
            await sleep( inputTokenized[1] );
          }
          break;          
          default:
          {
            await _query( inputTokenized, line );
          }
          break; 
        }
      } else {
        // 쿼리 처리 (멀티라인으로 합쳐진 쿼리)
        await _query( inputTokenized, line );
      }
    }

    if( parsedInput.opt['--repeat'] === undefined )
      break;
    else if( parsedInput.opt['--repeat'] === -1 )
      continue;
    else {
      parsedInput.opt['--repeat']--;
      if( parsedInput.opt['--repeat'] <= 0)  
        break;
    }
  }
}
  
module.exports._commandSh = _commandSh;