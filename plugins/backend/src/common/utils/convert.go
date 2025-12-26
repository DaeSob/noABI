package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	jMapper "cia/common/utils/json/mapper"
)

// I64UnixTimeTo : Unix Time To String
func I64UnixTimeTo(_i64Time int64) (string, error) {
	str := strconv.FormatInt(_i64Time, 10)
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return "", err
	}
	t := time.Unix(i, 0)
	//return t.Format(time.UnixDate), nil
	return t.Format("2006-01-02 15:04:05"), nil
}

// StringUnixTimeToString : string Unix Time To String
func StringUnixTimeToString(_strTime string) (string, error) {
	i, err := strconv.ParseInt(_strTime, 10, 64)
	if err != nil {
		return "", err
	}
	t := time.Unix(i, 0)
	//return t.Format(time.UnixDate), nil
	return t.Format("2006-01-02 15:04:05"), nil
}

// StringToI64 : String to Int64
func StringToI64(_str string) (int64, error) {
	i, err := strconv.ParseInt(_str, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func StringToINT(_str string) (int, error) {
	i, err := strconv.Atoi(_str)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func StringToI32(_str string) (int32, error) {
	i, err := StringToI64(_str)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func StringToUint32(_str string, base int) uint32 {
	i, err := strconv.ParseUint(_str, base, 32)
	if err != nil {
		panic(err)
	}

	return uint32(i)
}

func StringToUint8(_str string, base int) uint8 {
	i, err := strconv.ParseUint(_str, base, 8)
	if err != nil {
		panic(err)
	}

	return uint8(i)
}

func StringToUint64(_str string, base int) uint64 {
	i, err := strconv.ParseUint(_str, base, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func StringToUint64V2(_str string) uint64 {

	if len(_str) > 0 {
		if strings.Index(_str, "0x") == 0 || strings.Index(_str, "0X") == 0 {
			return HexStringToUint64(_str)
		} else {
			return StringToUint64(_str, 10)
		}
	}

	return 0
}

func IsHexString(_str string) bool {
	if len(_str) > 0 {
		if strings.Index(_str, "0x") == 0 || strings.Index(_str, "0X") == 0 {
			return true
		}
	}
	return false
}

func StringToFloat(_str string) float64 {
	f, err := strconv.ParseFloat(_str, 64)
	if err != nil {
		panic(err)
	}

	return f
}

func HexToInt64(_str string) int64 {
	// if has '0x' prefix, remove it
	_str = remove0xPrefix(_str)

	i, err := strconv.ParseInt(_str, 16, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func HexToUint64(_str string) uint64 {
	// if has '0x' prefix, remove it
	_str = remove0xPrefix(_str)

	return StringToUint64(_str, 16)
}

func HexStringToUint64(_str string) uint64 {
	// if has '0x' prefix, remove it
	_str = remove0xPrefix(_str)

	return StringToUint64(_str, 16)
}

func HexStringToDecimalString(_str string) string {
	_str = remove0xPrefix(_str)

	i, err := strconv.ParseUint(_str, 16, 64)
	if err != nil {
		panic(err)
	}
	return strconv.FormatUint(i, 10)
}

func DecimalStringToHexString(_str string) string {
	i, err := strconv.ParseUint(_str, 10, 64)
	if err != nil {
		panic(err)
	}
	return "0x" + strconv.FormatUint(i, 16)
}

func StringToDecimalString(_str string) string {
	if IsHexString(_str) {
		return HexStringToDecimalString(_str)
	}
	return _str
}

func StringToHexString(_str string) string {
	if IsHexString(_str) {
		return _str
	}
	return DecimalStringToHexString(_str)
}

// string to bool
func StringToBool(_str string) bool {
	res, _ := strconv.ParseBool(_str)
	return res
}

// I64ToString : String to Int64
func I64ToString(_i64 int64) string {
	return strconv.FormatInt(_i64, 10)
}

func I64ToHexString(_i64 int64) string {
	return "0x" + strconv.FormatInt(_i64, 16)
}

// I32ToString : String to Int64
func I32ToString(_i32 int32) string {
	return strconv.FormatInt(int64(_i32), 10)
}

// Uint64ToString : String to Uint64
func Uint64ToString(_uint64 uint64) string {
	return strconv.FormatUint(_uint64, 10)
}

// Uint64ToString : String to Uint64
func Uint8ToString(_uint8 uint8) string {
	return strconv.FormatUint(uint64(_uint8), 10)
}

func Uint64ToHexString(_uint64 uint64) string {
	return "0x" + strconv.FormatUint(_uint64, 16)
}

func FloatToString(_float64 float64) string {
	return strconv.FormatFloat(_float64, 'f', 0, 64)
}

// _arrFloat : []float 이어야 함
func FloatArrayToStringArray(_arrFloat []interface{}) []string {
	if _arrFloat != nil {
		result := []string{}
		for _, f := range _arrFloat {
			result = append(result, FloatToString(f.(float64)))
		}

		return result
	}

	return []string{}
}

// string to big int
func StringToBigInt(_value string) *big.Int {
	i := new(big.Int)

	if len(_value) > 0 {
		_, err := fmt.Sscan(_value, i)
		if err != nil {
			panic(err)
		}
	} else {
		_, err := fmt.Sscan("0", i)
		if err != nil {
			panic(err)
		}
	}

	return i
}

// string to BigFloat parse string value to big.Float
func StringToBigFloat(_value string) *big.Float {
	f := new(big.Float)
	// IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)

	_, err := fmt.Sscan(_value, f)
	if err != nil {
		panic(err)
	}

	return f
}

func hexToBytes(_str string) []byte {
	bytes, _ := hex.DecodeString(_str)
	return bytes
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(_str string) bool {
	return len(_str) >= 2 && _str[0] == '0' && (_str[1] == 'x' || _str[1] == 'X')
}

func remove0xPrefix(_str string) string {
	// if has '0x' prefix, remove it
	if has0xPrefix(_str) {
		_str = _str[2:]
	}
	return _str
}

func HexToBytes(_str string) []byte {
	// if has '0x' prefix, remove it
	_str = remove0xPrefix(_str)
	if len(_str)%2 == 1 {
		_str = "0" + _str
	}
	return hexToBytes(_str)
}

// byte to hex string
func BytesToHexString(_bytes []byte) string {
	return "0x" + hex.EncodeToString(_bytes)
}

// string map to interface array
func StringMapToInterfaceArray(_map map[string]interface{}) (result []interface{}) {
	for _, value := range _map {
		result = append(result, value)
	}
	return
}

// string array to json string
func StringArrayToJsonString(_arr []string, _pretty bool) string {
	json, _ := json.Marshal(_arr)
	return fmt.Sprint(string(json))
}

// interface to json string
func InterfaceToJsonString(_interface interface{}, _pretty bool) string {
	bytes, err := jMapper.ToJson(_interface)
	if err != nil {
		panic(err)
	}

	jMap, err := jMapper.NewBytes(bytes)
	if err != nil {
		panic(err)
	}

	if _pretty {
		return jMap.PPrint()
	} else {
		return jMap.Print()
	}
}

// interface array to string array
func InterfaceArrayToStringArray(_interface []interface{}) []string {
	if _interface != nil {
		result := []string{}
		for _, inter := range _interface {
			result = append(result, inter.(string))
		}
		return result
	}
	return []string{}
}
