package gointernals

import (
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

func TypeFor[T any]() *abi.Type {
	return (*PointerType)(abi.NoEscape(unsafe.Pointer(TypeOf((*T)(nil))))).Elem
}

func TypeOf(v any) *abi.Type {
	return abi.TypeOf(v)
}

func EfaceOf(v any) *abi.Eface {
	return (*abi.Eface)(abi.NoEscape(unsafe.Pointer(&v)))
}

func AnyFrom(e *abi.Eface) any {
	return *(*any)(abi.NoEscape(unsafe.Pointer(e)))
}
