package utils

import (
	"reflect"
	"unsafe"
)

func Setchar(s *string, idx int, ch byte) {
	if s == nil || idx < 0 || idx >= len(*s) {
		return
	}
	v := (*reflect.StringHeader)(unsafe.Pointer(s))
	var b []byte
	bs := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bs.Data, bs.Len, bs.Cap = v.Data, v.Len, v.Len
	b[idx] = ch
}
