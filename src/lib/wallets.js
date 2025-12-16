"use strict";
/*




*/
const ethJSWallet = require('ethereumjs-wallet');
class CWallets {
    constructor() {
        this.wallets = {};
    }
    AddWallet(_wallet) {
        try {
            const address = _wallet.getAddressString();
            this.wallets[address] = _wallet;
        }
        catch (e) {
            console.log(e);
        }
        finally {
        }
    }
    GetWallet(_address) {
        try {
            return this.wallets[_address.toLowerCase()];
        }
        catch (e) {
            console.log(e);
        }
        finally {
        }
    }
    RemoveWallet(_address) {
        try {
            delete this.wallets[_address.toLowerCase()];
        }
        catch (e) {
            console.log(e);
        }
        finally {
        }
    }
    GetSize() {
        try {
            return Object.keys(this.wallets).length;
        }
        catch (e) {
            console.log(e);
        }
        finally {
        }
    }
    GetAddresses() {
        return Object.keys(this.wallets);
    }
}
module.exports.CWallets = CWallets;
