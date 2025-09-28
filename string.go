package gointernals

import "unsafe"

type String struct {
	ptr unsafe.Pointer
	len int
}

func StringUnpack(s string) *String {
	return (*String)(unsafe.Pointer(&s))
}

func StringPack(s *String) string {
	return *(*string)(unsafe.Pointer(s))
}
