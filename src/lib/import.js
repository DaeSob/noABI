/*




*/

const bip39       = require("ethereum-cryptography/bip39");
const ethJSWallet = require('ethereumjs-wallet')
const ethUtil     = require('ethereumjs-util');
const fs          = require("fs")

/*
import bip39       from 'ethereum-cryptography/bip39';
import ethJSWallet from 'ethereumjs-wallet'
import ethUtil     from 'ethereumjs-util';
import fs          from 'fs'
*/

const BASE_DERIVATION_PATH = `m/44'/60'/0'/0/`;

async function openWalletFromPrivateKey( _privateKey ) {

  let privateKeyBuffer = undefined;
  try{
    privateKeyBuffer = ethUtil.toBuffer(_privateKey);
    return await ethJSWallet.default.fromPrivateKey( privateKeyBuffer);
  } catch(e) {
    throw e;
  } finally {
    privateKeyBuffer = undefined;
  }

}

async function openWalletFromKeystoreV3( _filename, _password ) {

  let strJson       = undefined;
  try{
    strJson = fs.readFileSync( _filename, 'utf8' );
    return  await ethJSWallet.default.fromV3( strJson, _password );
  } catch(e) {
    throw e;
  } finally {
    strJson       = undefined;
  }

}

async function openWalletFromMnemonic( _mnemonic ) {
  return openHDWallet( _mnemonic, 0 );
}

async function openHDWallet( _mnemonic, _idx ) {
  let seed        = undefined;
  let masterSeed  = undefined;
  try{
    seed        = await bip39.mnemonicToSeed( _mnemonic );
    masterSeed  = ethJSWallet.hdkey.fromMasterSeed(seed);
    return await masterSeed.derivePath( BASE_DERIVATION_PATH + _idx).getWallet();
  } catch(e) {
    console.log( e );
  } finally {
    seed        = undefined;
    masterSeed  = undefined;
  }
}

module.exports.OpenWalletFromPrivateKey = openWalletFromPrivateKey;
module.exports.OpenWalletFromKeystoreV3 = openWalletFromKeystoreV3;
module.exports.OpenHDWallet             = openHDWallet;
module.exports.OpenWalletFromMnemonic   = openWalletFromMnemonic;