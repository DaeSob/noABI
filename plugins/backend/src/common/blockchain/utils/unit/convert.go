package unit

import (
	"cia/common/utils"
	"fmt"
	"math/big"
	"strings"
)

// unit kind
// Wei     1                                           1
// Kwei    1,000                                       10^3
// Mwei    1,000,000                                   10^6
// Gwei    1,000,000,000                               10^9
// Szabo   1,000,000,000,000                           10^12
// Finney  1,000,000,000,000,000                       10^15
// Ether   1,000,000,000,000,000,000                   10^18
const (
	Wei    = 1
	KWei   = 1e3
	MWei   = 1e6
	GWei   = 1e9
	Szabo  = 1e12
	Finney = 1e15
	Ether  = 1e18
)

func unitToInt(_unit string) *big.Int {
	unit := strings.ToLower(_unit)

	switch unit {
	case "wei":
		return big.NewInt(Wei)
	case "kwei":
		return big.NewInt(KWei)
	case "mwei":
		return big.NewInt(MWei)
	case "gwei":
		return big.NewInt(GWei)
	case "szabo":
		return big.NewInt(Szabo)
	case "ether":
		return big.NewInt(Ether)
	default:
		return big.NewInt(Wei)
	}
}

func unitToFloat(_unit string) *big.Float {
	unit := strings.ToLower(_unit)

	switch unit {
	case "wei":
		return big.NewFloat(Wei)
	case "kwei":
		return big.NewFloat(KWei)
	case "mwei":
		return big.NewFloat(MWei)
	case "gwei":
		return big.NewFloat(GWei)
	case "szabo":
		return big.NewFloat(Szabo)
	case "ether":
		return big.NewFloat(Ether)
	default:
		return big.NewFloat(Wei)
	}
}

func ConvertFromWei(_value, _toUnit string) string {
	value := utils.StringToBigInt(_value)
	return fmt.Sprintf("%.18f", fromWei(value, _toUnit))
}

func ConvertToWei(_value, _fromUnit string) string {
	value := utils.StringToBigFloat(_value)
	return toWei(value, _fromUnit).String()
}

func SplitEther(_value string) [2]string {
	etherValue := ConvertFromWei(_value, "ether")
	value := utils.StringToBigFloat(etherValue)

	// part of ether
	fracInt, _ := value.Int(nil)
	partEther := fracInt.String()

	// part of wei
	fracStr := strings.Split(fmt.Sprintf("%.18f", value), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	partWei := fracStr

	return [2]string{partEther, partWei}
}

func toWei(_value *big.Float, _fromUnit string) *big.Int {
	// parse integer & * unit
	truncInt, _ := _value.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, unitToInt(_fromUnit))

	// parse floating number part & * 10^18
	fracStr := strings.Split(fmt.Sprintf("%.18f", _value), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)

	// integer + floating number part
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func fromWei(_wei *big.Int, _toUnit string) *big.Float {
	f := new(big.Float)
	// IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)

	fWei := new(big.Float)
	// IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetPrec(236)
	fWei.SetMode(big.ToNearestEven)

	return f.Quo(fWei.SetInt(_wei), unitToFloat(_toUnit))
}
