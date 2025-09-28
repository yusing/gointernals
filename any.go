package gointernals

import (
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

func EfaceOf[T any](v T) *abi.Eface {
	vAny := any(v)
	return (*abi.Eface)(unsafe.Pointer(&vAny))
}

func AnyFrom(e *abi.Eface) any {
	return *(*any)(unsafe.Pointer(e))
}
