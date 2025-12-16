/*




*/
const fs = require('fs');

async function exportToKeystoreV3( _wallet, _password, _pretty ) {

  let fileName          = undefined;
  let keyStoreJson      = undefined;
  try{
    fileName  = _wallet.getV3Filename();
    
    if( _pretty == true ){
      keyStoreJson   = await _wallet.toV3(_password);
      await fs.writeFileSync(fileName, JSON.stringify(keyStoreJson, null, 2));      
    } else {
      keyStoreJson   = await _wallet.toV3String(_password);
      await fs.writeFileSync(fileName, keyStoreJson);
    }
    return true;
  } catch(e) {
    console.log( e )
  } finally {
    keyStoreJson    = undefined;
    fileName        = undefined;
  }

}

module.exports.ExportToKeystoreV3 = exportToKeystoreV3;