/*




*/
const Web3ABI       = require('web3-eth-abi');
const Web3Utils     = require('web3-utils');
const { KlaySendTransaction, EthSendTransaction } = require( "./sendtx.js")
const {RpcRequest}  = require('./rpc/rpcRequest.js');

/*
async function klayDeploy( _object ) {

    let encodedParameters               = undefined;
    let bytecodeWithEncodedParameters   = undefined;

    let count                           = undefined;
    let gasPrice                        = undefined;
    let gasLimit                        = undefined;

    try{
        if( _object.typesArray !== undefined && _object.parameters !== undefined ) {
            encodedParameters               = Web3ABI.encodeParameters( _object.typesArray, _object.parameters ).slice(2);
            bytecodeWithEncodedParameters   = _object.byteCode + encodedParameters;
        } else {
            bytecodeWithEncodedParameters   = _object.byteCode;
        }

        console.log(_object.provider.rpcUrl);
        let   res           = await RpcRequest(_object.provider.rpcUrl, "klay_getTransactionCount", [_object.wallet.getAddressString(),"latest"], 1 );
        count               = res.result;          
        res                 = await RpcRequest( _object.provider.rpcUrl, "klay_gasPrice", [], 2 );
        gasPrice            = res.result;        
        res                 = await RpcRequest( _object.provider.rpcUrl, "klay_estimateGas", [{
                                                from:       _object.wallet.getAddressString(),
                                                input:      `0x${bytecodeWithEncodedParameters}`,
                                                gasPrice:   gasPrice
                                            }], 3 );
        if( res.error !== undefined )
            throw new Error( res.error.message );
        gasLimit = res.result;                         

        return await KlaySendTransaction( {
            provider: _object.provider,
            wallet:   _object.wallet,
            rawTransaction: {
                from:       _object.wallet.getAddressString(),
                nonce:      count,
                value:      '0x0',
                gasLimit:   gasLimit,
                gasPrice:   gasPrice,
                input:      `0x${bytecodeWithEncodedParameters}`,
            }
        })
        
    } catch( e ) {
        throw e;
    }finally{
        encodedParameters               = undefined;
        bytecodeWithEncodedParameters   = undefined;
    }
}
*/
  
async function ethDeploy( _object ) {

    let encodedParameters               = undefined;
    let bytecodeWithEncodedParameters   = undefined;
    
    let count                           = undefined;
    let maxFeePerGas                    = undefined;
    let gasPrice                        = undefined;
    let gasLimit                        = undefined;

    try{
        if( _object.typesArray !== undefined && _object.parameters !== undefined ) {
            encodedParameters               = Web3ABI.encodeParameters( _object.typesArray, _object.parameters ).slice(2);
            bytecodeWithEncodedParameters   = _object.byteCode + encodedParameters;
        } else {
            bytecodeWithEncodedParameters   = _object.byteCode;
        }

        let res         = await RpcRequest( _object.provider.rpcUrl, "eth_getTransactionCount", [_object.wallet.getAddressString(), "latest"], 1 );
        count           = res.result;

        if( customChain === undefined || customChain.gasPrice === undefined ) {
            res         = await RpcRequest( _object.provider.rpcUrl, "eth_gasPrice", [], 2 );
            if( res.error !== undefined )
                throw new Error( res.error.message );   
            gasPrice    = res.result;

            gasPriceAndTip = Number( gasPrice ) + Math.round( Number( gasPrice ) * 0.01 )
            console.log("gas price:", gasPriceAndTip)
            gasPrice = Web3Utils.toHex(gasPriceAndTip)

            res             = await RpcRequest( _object.provider.rpcUrl, "eth_getBlockByNumber", ["latest", false], 3 );
            if( res.error !== undefined )
                throw new Error( res.error.message );            
            maxFeePerGas    = res.result.baseFeePerGas;  

        } else {
            gasPrice    = Web3Utils.toHex( customChain.gasPrice );
        }
        
        if( customChain === undefined || customChain.gasLimit === undefined ) {
            res         = await RpcRequest( _object.provider.rpcUrl, "eth_estimateGas", [{
                                    from:       _object.wallet.getAddressString(),
                                    data:       `0x${bytecodeWithEncodedParameters}`,
                                    value:      '0x0',
                                    gasPrice:   gasPrice
                                }], 4 );
            if( res.error !== undefined )
                throw new Error( res.error.message );
            gasLimit = res.result;

        } else {
            gasLimit = Web3Utils.toHex( customChain.gasLimit );
        }

        return await EthSendTransaction({
            provider: _object.provider,
            wallet:   _object.wallet,
            bcc: blockCreationCycle,  //블록생성 주기
            rawTransaction: {
                nonce:          count,
                from:           _object.wallet.getAddressString(),
                gasLimit:       gasLimit,
                maxFeePerGas:   maxFeePerGas,
                gasPrice:       gasPrice,
                value:          '0x0',
                data:           `0x${bytecodeWithEncodedParameters}`,
            } 
        })

    } catch( e ) {
        console.log(e)
        throw e;
    }finally{
        encodedParameters               = undefined;
        bytecodeWithEncodedParameters   = undefined;
    }    
}

//module.exports.KlayDeploy  = klayDeploy;
module.exports.EthDeploy  = ethDeploy;
