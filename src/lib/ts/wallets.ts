/*




*/
const ethJSWallet = require('ethereumjs-wallet');

class CWallets {

    private wallets: { [address: string]: ethJSWallet };

    constructor() {
        this.wallets = {};
    }

    public AddWallet( _wallet: ethJSWallet ): void {
        try{
            const address = _wallet.getAddressString()
            this.wallets[ address ] = _wallet;
        } catch( e ) {
            console.log( e )
        } finally {
        }
    }

    public GetWallet( _address: string ): ethJSWallet {
        try{
            return this.wallets[ _address.toLowerCase()  ];
        } catch( e ) {
            console.log( e )
        } finally {
        }
    }

    public RemoveWallet( _address: string ): void {
        try{
            delete this.wallets[ _address.toLowerCase()  ];
        } catch( e ) {
            console.log( e )
        } finally {
        }
    }
    
    public GetSize(): number {
        try{
            return Object.keys( this.wallets ).length;
        } catch( e ) {
            console.log( e )
        } finally {
        }
    }

    public GetAddresses(): string[]{
        return Object.keys( this.wallets );
    }    

}

module.exports.CWallets = CWallets;