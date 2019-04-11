package utils

import (
	"fmt"
	"bytes"
)

const (
	_base = 34
	_max_code_len = 6
)
var _baseChars = []byte("ABCDEF12345GHJKM6789NPQRSTUV0WXYZL")
var _charIdx map[byte] int

func init() {
	/*
	if len(_baseChars) != _base {
		fmt.Printf("Please give me a string with length %d\n", _base)
	} */
	_charIdx = make(map[byte]int, _base)
	for i, v := range _baseChars {
		_charIdx[v] = i
	}
}

func Int2Code(n uint64) (string, error) {
	d := n
	m := uint64(0)
	buf := bytes.Buffer{}
	for d != 0 {
		m, d = d % _base, d / _base
		buf.WriteByte(_baseChars[int(m)])
	}

	length := buf.Len()
	if length > _max_code_len {
		return "", fmt.Errorf("too big interger")
	}
	if length == _max_code_len {
		return buf.String(), nil
	}

	res := make([]byte, 6)
	paddingLen := 6 - length
	for i:=paddingLen-1; i>=0; i-- {
		res[i] = _baseChars[0]
	}
	b := buf.Bytes()
	for i,j := 0,paddingLen; i<length; i,j = i+1,j+1 {
		res[j] = b[i]
	}
	return string(res), nil
}

func Code2Int(str string) (uint64, error) {
	if str == "" {
		return 0, fmt.Errorf("empty code given")
	}
	length := len(str)
	if length == 0 {
		return 0, nil
	}

	code := bytes.ToUpper([]byte(str))
	res := uint64(0)
	b := uint64(1)
	firstNon0 := 0
	for ; firstNon0 < length; firstNon0++ {
		if code[firstNon0] != _baseChars[0] {
			break
		}
	}
	for i:=firstNon0; i<length; i++ {
		d, ok := _charIdx[code[i]]
		if !ok {
			return 0, fmt.Errorf("bad code")
		}

		res +=  b*uint64(d)
		b *= _base
	}
	return res, nil
}
