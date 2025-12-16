const Web3ABI = require('web3-eth-abi');

// types - shoud be array
// value - shoud be encoded data
function _decodeParam(types, value) {
  try {
    return Web3ABI.decodeParameter(types, value);
  } catch (e) {
    console.log(e);
  } finally {
  }
}

function _decodeParams(types, value) {
  try {
    return Web3ABI.decodeParameters(types, value);
  } catch (e) {
    console.log(e);
  } finally {
  }
}

module.exports.DecodeParam = _decodeParam;
module.exports.DecodeParams = _decodeParams;
