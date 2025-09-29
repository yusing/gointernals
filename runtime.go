package gointernals

import (
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

// Functions below pushed from runtime.

//go:linkname fatal runtime.fatal
func fatal(s string)

//go:linkname rand runtime.rand
func rand() uint64

//go:linkname memmove runtime.memmove
//go:noescape
func memmove(dst, src unsafe.Pointer, n uintptr)

//go:linkname typedmemmove runtime.typedmemmove
//go:noescape
func typedmemmove(typ *abi.Type, dst, src unsafe.Pointer)

//go:linkname typedmemclr runtime.typedmemclr
//go:noescape
func typedmemclr(typ *abi.Type, ptr unsafe.Pointer)

//go:linkname newarray runtime.newarray
//go:noescape
func newarray(typ *abi.Type, n int) unsafe.Pointer

//go:linkname newobject runtime.newobject
//go:noescape
func newobject(typ *abi.Type) unsafe.Pointer
