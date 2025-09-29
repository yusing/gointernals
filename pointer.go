package gointernals

import (
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

type PointerType struct {
	abi.Type
	Elem *abi.Type
}

//go:nosplit
func PointerCast[ToT any, FromT any](src *FromT) *ToT {
	return (*ToT)(unsafe.Pointer(src))
}
