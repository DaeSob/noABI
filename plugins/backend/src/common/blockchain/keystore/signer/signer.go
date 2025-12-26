package signer

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"

	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

type TSigner struct {
	Address common.Address
	ks      *keystore.KeyStore
	chainId *big.Int
}

type TRawTransaction struct {
	To       common.Address `json:"to"`
	From     common.Address `json:"from"`
	Nonce    uint64         `json:"nonce"`
	Value    *big.Int       `json:"value"`
	GasPrice *big.Int       `json:"gasPrice"`
	GasLimit uint64         `json:"gasLimit"`
	Data     []byte         `json:"data"`
}

type Key struct {
	Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	Address common.Address
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	PrivateKey *ecdsa.PrivateKey
}

var once sync.Once
var lockMutex sync.Mutex
var mmKeyStore map[uint64]map[common.Address]*TSigner
var ks *keystore.KeyStore

func NewSigner(_dir string, _chainId *big.Int, _address common.Address) (tSigner *TSigner) {

	defer func() {
		if e := recover(); e != nil {
			tSigner = nil
		}
		return
	}()

	once.Do(func() {
		mmKeyStore = make(map[uint64]map[common.Address]*TSigner)
		ks = keystore.NewKeyStore(_dir, keystore.StandardScryptN, keystore.StandardScryptP)
	})

	lockMutex.Lock()
	defer lockMutex.Unlock()

	tSigner = mmKeyStore[(*_chainId).Uint64()][_address]
	if tSigner == nil {
		tSigner = new(TSigner)
		tSigner.ks = ks // keystore.NewKeyStore(_dir, keystore.StandardScryptN, keystore.StandardScryptP)
		tSigner.Address = _address
		tSigner.chainId = _chainId
		if len(mmKeyStore[(*_chainId).Uint64()]) == 0 {
			mmKeyStore[(*_chainId).Uint64()] = make(map[common.Address]*TSigner)
		}
		mmKeyStore[(*_chainId).Uint64()][_address] = tSigner
	}
	return

}

func (signer *TSigner) PreCheck(auth string) error {
	_, err := signer.ks.Find(accounts.Account{Address: signer.Address})
	return err
}

func (signer *TSigner) GetKey(auth string) (*keystore.Key, error) {
	// Load the key from the keystore and decrypt its contents
	account, err := signer.ks.Find(accounts.Account{Address: signer.Address})

	keyjson, err := os.ReadFile(account.URL.Path)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keyjson, auth)
	if err != nil {
		return nil, err
	}

	// Make sure we're really operating on the requested key (no swap attacks)
	if key.Address != account.Address {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, account)
	}

	return key, nil
}

func (signer *TSigner) SignRawTx(_phrase string, _rawTx TRawTransaction) (*types.Transaction, error) {
	account, err := signer.ks.Find(accounts.Account{Address: signer.Address})
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(_rawTx.Nonce, _rawTx.To, _rawTx.Value, _rawTx.GasLimit, _rawTx.GasPrice, _rawTx.Data)
	return signer.ks.SignTxWithPassphrase(account, _phrase, tx, signer.chainId)
}

func SignedRawTx2HexString(_signedTx *types.Transaction) string {
	ts := types.Transactions{_signedTx}
	b := new(bytes.Buffer)
	ts.EncodeIndex(0, b)
	rawTxBytes := b.Bytes()
	return hex.EncodeToString(rawTxBytes)
}

func (signer *TSigner) SignHash(_phrase string, _hash []byte) (signature []byte, err error) {
	account, err := signer.ks.Find(accounts.Account{Address: signer.Address})
	if err != nil {
		return nil, err
	}

	signature, err = signer.ks.SignHashWithPassphrase(account, _phrase, _hash)

	if bytes.Equal(signature[64:], []byte{0}) {
		signature[64] = 27
	} else if bytes.Equal(signature[64:], []byte{1}) {
		signature[64] = 28
	}

	return
}

func (signer *TSigner) SignMessage(_phrase string, data ...interface{}) (signature []byte, err error) {
	hashData := solsha3.SoliditySHA3(data...)
	hashPrefix := solsha3.SoliditySHA3WithPrefix(
		solsha3.Bytes32("0x" + hex.EncodeToString(hashData)),
	)
	signature, err = signer.SignHash(_phrase, hashPrefix)
	return
}

func ConvertPrivateKeyToString(key *keystore.Key) string {
	return fmt.Sprintf("0x%x", crypto.FromECDSA(key.PrivateKey))
}

func PublicKeyToString(key *keystore.Key) string {
	return fmt.Sprintf("%v", key.Address)
}

func (s *TSigner) EIP712UserOpSignature(_phrase, _type string, _domainData apitypes.TypedDataDomain, _messageData apitypes.TypedDataMessage) (sigHash []byte, signature []byte, err error) {

	signerData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"TUserOperation": []apitypes.Type{
				{Name: "accessKey", Type: "uint256"},
				{Name: "sender", Type: "address"},
				{Name: "nonce", Type: "uint256"},
				{Name: "unpaidGasFee", Type: "uint256"},
				{Name: "callGasLimit", Type: "uint256"},
				{Name: "preVerificationGas", Type: "uint256"},
				{Name: "maxFeePerGas", Type: "uint256"},
				{Name: "maxPriorityFeePerGas", Type: "uint256"},
				{Name: "initCode", Type: "bytes"},
				{Name: "callData", Type: "bytes"},
				{Name: "paymasterAndData", Type: "bytes"},
			},
			"Auth": []apitypes.Type{
				{Name: "accessKey", Type: "uint256"},
				{Name: "sender", Type: "address"},
				{Name: "nonce", Type: "uint256"},
				{Name: "unpaidGasFee", Type: "uint256"},
				{Name: "callGasLimit", Type: "uint256"},
				{Name: "maxFeePerGas", Type: "uint256"},
				{Name: "maxPriorityFeePerGas", Type: "uint256"},
				{Name: "initCode", Type: "bytes"},
				{Name: "callData", Type: "bytes"},
				{Name: "paymasterAndData", Type: "bytes"},
			},
		},
		PrimaryType: _type, //"TUserOperation",
		Domain:      _domainData,
		Message:     _messageData,
	}

	typedDataHash, _ := signerData.HashStruct(signerData.PrimaryType, signerData.Message)
	domainSeparator, _ := signerData.HashStruct("EIP712Domain", signerData.Domain.Map())

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))

	sigHash = crypto.Keccak256(rawData)

	signature, err = s.SignHash(_phrase, sigHash)

	return
}

func (s *TSigner) EIP712DataSignature(_phrase string, _eip712Data apitypes.TypedData) (sigHash []byte, signature []byte, err error) {

	typedDataHash, _ := _eip712Data.HashStruct(_eip712Data.PrimaryType, _eip712Data.Message)
	domainSeparator, _ := _eip712Data.HashStruct("EIP712Domain", _eip712Data.Domain.Map())

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))

	sigHash = crypto.Keccak256(rawData)

	signature, err = s.SignHash(_phrase, sigHash)

	return
}

func _recover(sigHash []byte, _signature string) {

	hexS, err := hexutil.Decode(_signature)
	hexS[64] -= 27

	sigPubkey, err := crypto.Ecrecover(sigHash, hexS)
	pubkey, err := crypto.UnmarshalPubkey(sigPubkey)
	if err != nil {
		fmt.Println(err)
		return
	}
	address := crypto.PubkeyToAddress(*pubkey)
	fmt.Println("recover", address)

}

func RecoverDataSigner(sigHash []byte, _signature []byte) (common.Address, error) {

	_signature[64] -= 27
	sigPubkey, err := crypto.Ecrecover(sigHash, _signature)
	pubkey, err := crypto.UnmarshalPubkey(sigPubkey)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil

}

func RecoverDataSignerFromString(_msg, _signature string) string {
	// msg
	hashMsg := solsha3.SoliditySHA3(solsha3.String(_msg))

	// sig
	sig := hexutil.MustDecode(_signature)

	msg := accounts.TextHash([]byte(hashMsg))
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		panic(err)
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	return recoveredAddr.Hex()

}
