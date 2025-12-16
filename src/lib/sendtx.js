/*




*/
const Web3Utils     = require('web3-utils');
const EthereumTx    = require("ethereumjs-tx").Transaction;
const Common        = require('ethereumjs-common').default;
const {RpcRequest}  = require('./rpc/rpcRequest.js');
const { sleep }     = require( "./utils/sleep.js");

async function _ethGetTransactionReceipt( _rpcUrl, _txHashCode, _bcc, _retryCnt ) {
    for( let i=0; i<_retryCnt; i++ ){ 
        await sleep( _bcc * i  );            
        let res = await RpcRequest( _rpcUrl, "eth_getTransactionReceipt", [_txHashCode], 1 );   
        if (res.result){
            return res;
        }
    }
    return undefined;
}

async function _ethSendTransaction( _object ) {
    let res         = undefined;
    try{
        let common = customChain !== undefined ? Common.forCustomChain(
            customChain.baseChain,
            customChain.param,
            customChain.hardFork,
            customChain.supportedHardFork
        ) : undefined;
        //0x 를 제거 하여야 한다
        let     privateKey  = Buffer.from(_object.wallet.getPrivateKeyString().substring(2), 'hex')
        const   tx          = new EthereumTx( _object.rawTransaction, common !== undefined ? {common} : undefined );
        tx.sign( privateKey );
        let serializedTx = '0x' + tx.serialize().toString('hex');

        const txHash = Web3Utils.keccak256(serializedTx);
        res = await RpcRequest( _object.provider.rpcUrl, "eth_sendRawTransaction", [ serializedTx ], 1 );
        if( res.error !== undefined ) {
            const errMsg = res.error.message || "";
            if (
                errMsg.includes("nonce too low") ||
                errMsg.includes("already known") ||
                errMsg.includes("already used")
            ) {
                return txHash;
            }            
            throw Error(res.error.message);
        }
        return res.result;
    }
    catch( e ) {
        throw e;
    } finally {
        res     = undefined;
    }
}

async function ethSendTransaction( _object ) {
    let res     = undefined;
    try{
        let txHash  = await _ethSendTransaction( _object );
        res = await _ethGetTransactionReceipt( _object.provider.rpcUrl, txHash, _object.bcc === undefined ? 1000 : _object.bcc, 10 );
        if( res.error !== undefined ) {
            throw Error(res.error.message);
        }        
        return res.result;
    } catch( e ) {
        throw e;
    } finally {
        res   = undefined
    }
}

module.exports.EthSendTransaction  = ethSendTransaction;