/*




*/
var ethTx = require('ethereumjs-tx');

async function signRawTx( _wallet, _rawTransaction ) {

    try{
        var tx = new ethTx( _rawTransaction );
        tx.sign(_wallet.getPrivateKey())
        return '0x' + tx.serialize().toString('hex');
    } catch( e ) {
        console.log( e )
    } finally {
    }    

}

module.exports.Sign = signRawTx;