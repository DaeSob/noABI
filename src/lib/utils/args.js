
function _parseArgs(argv) {
  const args = {};
  let currentKey = null;

  for (let i = 2; i < argv.length; i++) {
    const token = argv[i];

    // --longParam
    if (token.startsWith("--")) {
      currentKey = token.slice(2);
      args[currentKey] = true; // 기본값 true (flag)
    }
    // -s (short param)
    else if (token.startsWith("-")) {
      currentKey = token.slice(1);
      args[currentKey] = true;
    }
    // value
    else {
      if (currentKey) {
        args[currentKey] = token;
        currentKey = null;
      }
    }
  }

  return args;
}

module.exports.parseArgs = _parseArgs;