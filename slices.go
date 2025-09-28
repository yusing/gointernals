//go:build go1.24 && go1.25 && !go1.26

package gointernals

import (
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

type Slice struct {
	ptr unsafe.Pointer
	len int
	cap int
}

// internal/abi/type.go
type SliceType struct {
	abi.Type
	Elem *abi.Type
}

type SliceCloneFunc = func(src *Slice, elemType *abi.Type) *Slice

var width = unsafe.Sizeof(int(0))

//go:linkname growslice runtime.growslice
//go:noescape
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *abi.Type) Slice

//go:linkname makeslice runtime.makeslice
//go:noescape
func makeslice(et *abi.Type, len, cap int) unsafe.Pointer

// slicecopy is used to copy from a string or slice of pointerless elements into a slice.
//
//go:linkname slicecopy runtime.slicecopy
//go:noescape
func slicecopy(toPtr unsafe.Pointer, toLen int, fromPtr unsafe.Pointer, fromLen int, width uintptr)

//go:linkname typedslicecopy runtime.typedslicecopy
//go:noescape
func typedslicecopy(typ *abi.Type, dstPtr unsafe.Pointer, dstLen int, srcPtr unsafe.Pointer, srcLen int)

//go:nosplit
//go:linkname SliceClone gointernals.SliceClone
func SliceClone(src *Slice, elemType *abi.Type) *Slice {
	newSlice := makeslice(elemType, src.len, src.len)
	if !elemType.CanPointer() {
		slicecopy(newSlice, src.len, src.ptr, src.len, width)
	} else {
		typedslicecopy(elemType, newSlice, src.len, src.ptr, src.len)
	}
	return &Slice{
		ptr: newSlice,
		len: src.len,
		cap: src.len,
	}
}

//go:nosplit
func SliceCloneAs[T any](src *Slice, elemType *abi.Type) []T {
	return *(*[]T)(unsafe.Pointer(SliceClone(src, elemType)))
}

//go:nosplit
func SliceUnpack[T any](s []T) (*Slice, *abi.Type) {
	sAny := any(s)
	eface := EfaceOf(&sAny)
	return (*Slice)(unsafe.Pointer(&s)), (*SliceType)(unsafe.Pointer(eface.Type)).Elem
}
