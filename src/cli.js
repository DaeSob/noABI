const { initializeGlobals } = require('./core/state/globals.js');
const { _setConfiguration } = require('./core/cli/config.js');
const { _welcome } = require('./welcome.js');
const { _loadHistory } = require('./core/cli/history.js');
const { _prompt } = require('./core/cli/prompt.js');
const { rl } = require('./core/cli/readline.js');
const { enableLineHandler } = require('./core/cli/lineHandler.js');
const { getRpcUrl, getUser } = require('./core/state/globals.js');

// 전역 변수 초기화
initializeGlobals();

// CLI 시작
_setConfiguration().then(() => {
  _welcome(getRpcUrl(), getUser());
  _loadHistory();  // History 로드
  rl.setPrompt(_prompt());
  rl.prompt();
  enableLineHandler(); // 라인 핸들러 활성화
});
