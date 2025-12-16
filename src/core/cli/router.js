const { printError } = require('../log/printScreen.js');
const { _commandVersion } = require('../command/version.js');
const { _commandSet } = require('../command/set.js');
const { _commandImport } = require('../command/import.js');
const { _commandExport } = require('../command/export.js');
const { _commandShow } = require('../command/show.js');
const { _commandWhoAmI } = require('../command/whoami.js');
const { _commandCompose } = require('../command/compose.js');
const { _commandDeploy } = require('../command/deploy.js');
const { _commandSh } = require('../command/sh.js');
const { _commandHelp } = require('../command/help.js');
const { _commandMemory } = require('../command/memory.js');
const { _commandToUTC } = require('../command/time.js');
const { _commandGetBalance } = require('../command/getBalance.js');
const { _commandEstimateGas } = require('../command/estimateGas.js');
const { _commandUse } = require('../command/use.js');
const { _commandAbi } = require('../command/abi.js');
const { _commandKeystore } = require('../command/keyStore.js');
const { _commandSnapshot } = require('../command/snapshot.js');
const { _commandEcho, _mapEchoVariables } = require('../command/echo.js');
const { _commandCurl } = require('../command/curl.js');
const { _query } = require('../chain/executeQuery.js');
const { _commandTest } = require('../command/test.js');

async function _do(_inputTokens, _line) {
  try {
    switch (_inputTokens[0]) {
      case 'version': {
        await _commandVersion(_inputTokens);
      }
      break;
      case 'set':
        {
          await _commandSet(_inputTokens, _line);
        }
        break;
      case 'import':
        {
          await _commandImport(_inputTokens, _line);
        }
        break;
      case 'export':
        {
          await _commandExport(_inputTokens);
        }
        break;
      case 'ls':
      case 'show':
        {
          await _commandShow(_inputTokens);
        }
        break;
      case 'whoami':
        {
          await _commandWhoAmI(_inputTokens);
        }
        break;
      case 'help':
        {
          await _commandHelp(_inputTokens);
        }
        break;
      case 'compose':
        {
          await _commandCompose(_inputTokens);
        }
        break;
      case 'deploy':
        {
          await _commandDeploy(_inputTokens);
        }
        break;
      case 'run':
      case 'sh':
        {
          await _commandSh(_inputTokens);
        }
        break;
      case 'mem':
      case 'memory':
        {
          await _commandMemory(_inputTokens);
        }
        break;
      case 'use':
        {
          await _commandUse(_inputTokens);
        }
        break;
      case 'snapshot':
        {
          await _commandSnapshot(_inputTokens);
        }
        break;
      case 'toUtc': {
        await _commandToUTC(_inputTokens, _line);
      }
      break;
      case 'getBalance': {
        await _commandGetBalance(_inputTokens, _line);
      }
      break;
      case 'estimateGas': {
        await _commandEstimateGas(_inputTokens);
      }
      break;
      case 'abi':
        {
          await _commandAbi(_inputTokens, _line);
        }
        break;
      case 'keystore': {
        await _commandKeystore(_inputTokens);
      }
      break;
      case 'curl':
        {
          await _commandCurl(_inputTokens, _line);
        }
        break;
      case 'test':
        {
          await _commandTest(_inputTokens);
        }
        break;
      case 'echo':
        {
          // sh 전용
        }
        break;
      default: {
        // 명령어가 아니고 쿼리 형식도 아닌 경우, 입력된 값을 그대로 출력 (echo처럼 처리)
        const trimmedLine = _line.trim();
        const regQueryFormat = /((.*?)\.(.*?)\((.*?)\)\s*\.\s*(.*?)\((.*?)\))/; // 쿼리 형식인지 확인

        if (!regQueryFormat.test(trimmedLine)) {
          // 쿼리 형식이 아닌 경우 echo처럼 처리하여 프롬프트에 표시
          const { _prompt } = require('./prompt.js');
          const echoValue = _mapEchoVariables(_inputTokens, trimmedLine);
          console.log(_prompt(echoValue));
          await _commandEcho(_inputTokens, trimmedLine);
        } else {
          // 쿼리 형식인 경우 기존대로 처리
          await _query(_inputTokens, _line);
        }
      }
    }
  } catch (e) {
    printError('\x1b[31m' + e.toString());
  }
}

module.exports = {
  _do
};

