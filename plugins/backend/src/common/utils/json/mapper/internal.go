package mapper

import (
	"bytes"
	"encoding/json"
)

// prettyPrint : pretty-print byte to json string
func prettyPrint(_src []byte) (string, error) {
	var dst bytes.Buffer
	err := json.Indent(&dst, _src, "", "  ")
	if nil != err {
		return "", err
	}
	return dst.String(), nil
}
