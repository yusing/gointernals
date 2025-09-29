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

func (s *Slice) Ptr() unsafe.Pointer {
	return s.ptr
}

func (s *Slice) Len() int {
	return s.len
}

func (s *Slice) Cap() int {
	return s.cap
}

// internal/abi/type.go
type SliceType struct {
	abi.Type
	Elem *abi.Type
}

type SliceCloneFunc = func(src *Slice, elemType *abi.Type) *Slice

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
		slicecopy(newSlice, src.len, src.ptr, src.len, elemType.Size)
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
//go:linkname SliceCloneInto gointernals.SliceCloneInto
func SliceCloneInto(dst *Slice, src *Slice, elemType *abi.Type) {
	need := src.len - dst.cap
	if need <= 0 { // can fit in dst
		dst.len = src.len
	} else if dst.ptr != nil { // can grow dst
		grown := growslice(dst.ptr, src.len, dst.cap, need, elemType)
		*dst = grown
	} else { // dst is nil, assign new slice
		dst.ptr = makeslice(elemType, src.len, src.len)
		dst.len = src.len
		dst.cap = src.len
	}

	if !elemType.CanPointer() {
		slicecopy(dst.ptr, src.len, src.ptr, src.len, elemType.Size)
	} else {
		typedslicecopy(elemType, dst.ptr, src.len, src.ptr, src.len)
	}
}

//go:nosplit
func SliceCloneAs[T any](src *Slice, elemType *abi.Type) []T {
	return *(*[]T)(unsafe.Pointer(SliceClone(src, elemType)))
}

//go:nosplit
func SliceUnpack[T any](s []T) (*Slice, *abi.Type) {
	eface := EfaceOf(s)
	return (*Slice)(unsafe.Pointer(&s)), (*SliceType)(unsafe.Pointer(eface.Type)).Elem
}

//go:nosplit
func SliceHeader[T any](s []T) *Slice {
	return (*Slice)(unsafe.Pointer(&s))
}

//go:nosplit
func SlicePack[T any](s *Slice) []T {
	return *(*[]T)(unsafe.Pointer(s))
}

//go:nosplit
func SliceCast[ToT any, FromT any](src []FromT) (ret []ToT) {
	srcHeader, srcType := SliceUnpack(src)
	_, dstType := SliceUnpack(ret)

	if srcType.Kind != dstType.Kind || srcType.Size != dstType.Size {
		panic("SliceCast: type mismatch")
	}
	ret = SlicePack[ToT](srcHeader)
	return
}

//go:nosplit
func UnsafeSliceCast[ToT any, FromT any](src []FromT) []ToT {
	to := SlicePack[ToT](SliceHeader(src))
	return to
}
