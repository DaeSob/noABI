const Web3ABI = require('web3-eth-abi');
const EtherABI = require('ethereumjs-abi');

function _encodeAbi(signature, ...args) {
  try {
    return '0x' + EtherABI.simpleEncode(signature, ...args).toString('hex');
  } catch (e) {
    console.log(e);
  } finally {
  }
}

function _encodeFunctionCall(abi, params) {
  try {
    return Web3ABI.encodeFunctionCall(abi, params);
  } catch (e) {
    console.log(e);
  } finally {
  }
}

function _encodeFunctionSignature(name) {
  try {
    return Web3ABI.encodeFunctionSignature(name);
  } catch (e) {
    console.log(e);
  } finally {
  }
}

module.exports.EncodeAbi = _encodeAbi;
module.exports.EncodeFunctionCall = _encodeFunctionCall;
module.exports.EncodeFunctionSignature = _encodeFunctionSignature

