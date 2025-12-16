const path                      = require("path");

function _resolveWorkspacePath(scriptFile) {
  // 이미 절대경로면 그대로 사용
  if (path.isAbsolute(scriptFile)) {
    return scriptFile;
  }

  // 상대경로면 workspace.root 기준으로 자동 조합
  return path.join(workPath, scriptFile);
}

module.exports.resolveWorkspacePath = _resolveWorkspacePath;