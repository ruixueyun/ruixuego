// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"reflect"
	"unsafe"
)

// StringToBytes 字符串转字节切片
// 需要注意的是该方法极不安全，使用过程中应足够谨慎，防止各类访问越界的问题
//nolint
func StringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

// BytesToString 字节切片转字符串
//nolint
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
