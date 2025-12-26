package types

import (
	"cia/common/blockchain/abi"
	jMapper "cia/common/utils/json/mapper"
	"regexp"
	"strings"
)

func tokenizeParams(_r *regexp.Regexp, _extractedParams string) (res []abi.TArgumentStr) {
	b := _r.FindStringIndex(_extractedParams)
	p := _extractedParams[b[0]+1 : b[1]-1]
	p = strings.Trim(p, " ")
	if len(p) == 0 {
		return
	}
	params := strings.Split(p, ",")
	for _, v := range params {
		tokenized := strings.Split(strings.Trim(v, " "), " ")
		arg := abi.TArgumentStr{}
		if len(tokenized) == 3 {
			arg.Type = tokenized[0]
			arg.InternalType = tokenized[0]
			if tokenized[1] == "indexed" {
				arg.Indexed = true
			}
			arg.Name = tokenized[2]
		} else if len(tokenized) == 2 {
			arg.Type = tokenized[0]
			arg.InternalType = tokenized[0]
			arg.Indexed = false
			arg.Name = tokenized[1]
		} else if len(tokenized) == 1 {
			arg.Type = tokenized[0]
			arg.InternalType = tokenized[0]
			arg.Indexed = false
		}
		res = append(res, arg)
	}
	return
}

func CompileEvent(_function string) *TEvent {
	temp := strings.Trim(_function, " ")
	temp = strings.ReplaceAll(temp, "  ", " ")

	r, _ := regexp.Compile(`(\()(.*?)(\))`)

	h := (r.ReplaceAllString(temp, ""))
	h = strings.Trim(h, " ")
	tokenized := strings.Split(h, " ")

	if tokenized[0] != "event" {
		return nil
	}

	newEvent := TEvent{
		Type: tokenized[0],
		Name: tokenized[1],
	}

	newEvent.Inputs = tokenizeParams(r, temp)

	return &newEvent
}

func CompileFunction(_function string) *TMethod {
	temp := strings.Trim(_function, " ")
	temp = strings.ReplaceAll(temp, "  ", " ")

	r, _ := regexp.Compile(`(\()(.*?)(\))`)

	h := (r.ReplaceAllString(temp, ""))
	h = strings.Trim(h, " ")
	tokenized := strings.Split(h, " ")

	/*
		if tokenized[0] != "function" {
			return nil
		}
	*/

	stateMutability := ""
	for _, v := range tokenized {
		switch v {
		case "view":
			stateMutability = "view"
		case "pure":
			stateMutability = "pure"
		case "nonpayable":
			stateMutability = "nonpayable"
		case "payable":
			stateMutability = "payable"
		}
		if len(stateMutability) > 0 {
			break
		}
	}

	if len(stateMutability) == 0 {
		stateMutability = "nonpayable"
	}

	newFunc := TMethod{
		Type:            "function",
		Name:            tokenized[0],
		StateMutability: stateMutability,
	}

	idx := strings.LastIndex(temp, "returns")
	if idx > 0 {
		newFunc.Outputs = tokenizeParams(r, temp[idx:])
	}

	newFunc.Inputs = tokenizeParams(r, temp)

	return &newFunc
}

func CompileAbi(_abi interface{}) *TMethod {

	bytes, _ := jMapper.ToJson(_abi)

	var inputs []abi.TArgumentStr

	j, _ := jMapper.NewBytes(bytes)

	abiName, _ := j.Find("name").(string)
	abiType, _ := j.Find("type").(string)
	abiStateMutability, _ := j.Find("stateMutability").(string)

	abiInputs, _ := j.Find("inputs").([]interface{})
	for i := 0; i < len(abiInputs); i++ {

		arg := abiInputs[i].(map[string]interface{})
		input := abi.TArgumentStr{
			Name:         arg["name"].(string),
			Type:         arg["type"].(string),
			InternalType: arg["internalType"].(string),
		}

		if arg["components"] != nil {
			input.Components = getComponets(arg["components"])
		}

		if arg["indexed"] != nil {
			input.Indexed = arg["indexed"].(bool)
		}
		inputs = append(inputs, input)

	}

	return CreateMethodEx(
		abiType,
		abiName,
		abiStateMutability,
		inputs,
		[]abi.TArgumentStr{},
	)

}

func getComponets(_components interface{}) []abi.TArgumentStr {
	panic("not support inputs components types yet")
}
