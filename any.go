package gointernals

import (
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

func TypeFor[T any]() *abi.Type {
	var v T
	return EfaceOf(v).Type
}

func TypeOf[T any](v T) *abi.Type {
	return EfaceOf(v).Type
}

func EfaceOf[T any](v T) *abi.Eface {
	vAny := any(v)
	return (*abi.Eface)(unsafe.Pointer(&vAny))
}

func AnyFrom(e *abi.Eface) any {
	return *(*any)(unsafe.Pointer(e))
}
