package gointernals

import "unsafe"

type String struct {
	ptr unsafe.Pointer
	len int
}

func (s *String) Ptr() unsafe.Pointer {
	return s.ptr
}

func (s *String) Len() int {
	return s.len
}

func (s *String) String() string {
	return *(*string)(unsafe.Pointer(s))
}

func StringUnpack(s string) *String {
	return (*String)(unsafe.Pointer(&s))
}

func StringPack(s *String) string {
	return *(*string)(unsafe.Pointer(s))
}
