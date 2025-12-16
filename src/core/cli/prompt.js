const { textColor } = require('../log/printScreen.js');
const { getCurrentSet } = require('../state/globals.js');

function _prompt(_echoValue) {
  // keystore 등에서 사용할 수 있도록 전역으로 export
  global._prompt = _prompt;
  let prompt = '';

  const currentSet = getCurrentSet();
  prompt = textColor.pink + currentSet.wallet + textColor.white + '> ';

  // echo 명령어의 매핑된 값이 있으면 프롬프트에 추가
  if (_echoValue !== undefined && _echoValue !== null && _echoValue !== '') {
    prompt += textColor.gray + _echoValue + textColor.white;
  }

  return prompt;
}

module.exports = {
  _prompt
};

