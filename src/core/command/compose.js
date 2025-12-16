const path = require('path');
const fs = require('fs');
const { _commandSetCollection } = require('./setCollection.js');
const { _compose } = require('../deploy/compose.js');
const { _getParseInput } = require('./common/currentPath.js');
const { resolveWorkspacePath } = require('../utils/resolvePath.js');
const { textColor, printError, printSuccess, printGray, printDefault } = require('../log/printScreen.js');

async function _commandCompose(_inputTokens) {
  const parsedInput = _getParseInput(_inputTokens, 1);

  const composeFile = resolveWorkspacePath(parsedInput.data[0]);
  const rawInputCompose = fs.readFileSync(composeFile, 'utf8');
  const objComposeFile = JSON.parse(rawInputCompose);

  if (objComposeFile.collection === undefined) {
    throw new String('not found a collection in compose-file');
  }
  _commandSetCollection(['use', objComposeFile.collection]);

  let outputPath = path.resolve(objSettings.workspace.output.build, currentSet.chain, currentCollection).replace(/\\/g, '/');
  outputPath = outputPath.charAt(0).toLowerCase() + outputPath.slice(1);

  if (objComposeFile.files === undefined) {
    objComposeFile['files'] = {};
  }

  objComposeFile.files['manifest'] = outputPath + '/output.obj';

  fs.mkdirSync(outputPath, { recursive: true });
  fs.writeFileSync(composeFile, JSON.stringify(objComposeFile, null, 2));

  await _compose(objComposeFile);
}

module.exports._commandCompose = _commandCompose;