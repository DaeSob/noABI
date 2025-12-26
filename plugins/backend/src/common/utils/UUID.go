package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"runtime"
	"strconv"

	uuid "github.com/google/uuid"
)

// create uuid
func CreateUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

// Sha256 : create sha256 hash value
func Sha256(_plainText string) string {
	hashBytes := sha256.Sum256([]byte(_plainText))

	// Wrong code : return string(hashBytes)
	// [32]byte cannot convert string directly
	// use hex.EncodeToString() function
	return hex.EncodeToString(hashBytes[:])
}

// Sha256 : create sha256 hash value
func Sha256b(_plainBytes []byte) string {
	hashBytes := sha256.Sum256(_plainBytes)

	// Wrong code : return string(hashBytes)
	// [32]byte cannot convert string directly
	// use hex.EncodeToString() function
	return hex.EncodeToString(hashBytes[:])
}

// FNV1a64 : create fnv-1a 64bit hash
//   - ref: https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
func FNV1a64(_plainText string) string {
	// using fnv-1a 64 hash
	h := fnv.New64a()
	h.Write([]byte(_plainText))
	sum := h.Sum64()

	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, sum)
	res := hex.EncodeToString(buffer)
	buffer = nil

	return res
}

func MD5(_plainText string) string {
	hashBytes := md5.Sum([]byte(_plainText))
	fmt.Println(hashBytes, len(hashBytes))
	return hex.EncodeToString(hashBytes[:])
}

func MD5b(_plainBytes []byte) string {
	hashBytes := md5.Sum(_plainBytes)
	return hex.EncodeToString(hashBytes[:])
}

// Encrypt : AES256 key, base64 encode - plain text to cipher text
// See alternate IV creation from ciphertext below
// var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
func Encrypt(_k string, _p string) (string, error) {
	// convert string to byte slice
	key := []byte(_k)
	plainText := []byte(_p)

	// create new key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// plaintext - encode base64
	b := base64.StdEncoding.EncodeToString(plainText)

	// create cipher text slice
	ciphertext := make([]byte, aes.BlockSize+len(b))

	// TODO : code view
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	// hex encode
	cipherHex := hex.EncodeToString(ciphertext)

	// return hex string
	return cipherHex, nil
}

// Decrypt : AES256 key, Base64 Decode - cipher text to plain text
func Decrypt(_k string, _c string) (string, error) {
	// convert string to byte slice
	key := []byte(_k)

	// hex decode
	text, err := hex.DecodeString(_c)
	if nil != err {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(text) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return "", err
	}

	// return plain text
	return string(data), nil
}

func Decrypt_base64(_k string, _c string) (string, error) {
	// convert string to byte slice
	key := []byte(_k)

	// hex decode
	// URLEncoding
	//ptb, err := base64.StdEncoding.DecodeString(c)
	ptb, err := base64.StdEncoding.DecodeString(_c)
	if nil != err {
		return "", err
	}

	hx := hex.EncodeToString(ptb)
	text, err := hex.DecodeString(hx)
	if nil != err {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(text) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return "", err
	}

	// return plain text
	return string(data), nil
}

// GORUTINE ID
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
