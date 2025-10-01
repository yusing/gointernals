package gointernals

import (
	"reflect"
	"unsafe"
)

//go:nosplit
func ReflectMakeSlice(typ reflect.Type, len, cap int) Slice {
	if typ.Kind() != reflect.Slice {
		panic("gointernals.ReflectMakeSlice of non-slice type")
	}
	if len < 0 {
		panic("gointernals.ReflectMakeSlice: negative len")
	}
	if cap < 0 {
		panic("gointernals.ReflectMakeSlice: negative cap")
	}
	if len > cap {
		panic("gointernals.ReflectMakeSlice: len > cap")
	}
	return Slice{
		ptr: reflect_unsafe_NewArray(ReflectTypeToABIType(typ.Elem()), cap),
		len: len,
		cap: cap,
	}
}

// ReflectInitSlice initializes a slice with the given length and capacity.
//
//go:nosplit
func ReflectInitSlice(dst reflect.Value, len, cap int) {
	if dst.Kind() != reflect.Slice {
		panic("gointernals.ReflectInitSlice of non-slice type")
	}

	s := (*Slice)(ReflectValueData(dst))
	if s.ptr != nil {
		s.len = len
		if s.cap < cap {
			*s = reflect_growslice(ReflectTypeToABIType(dst.Type().Elem()), *s, cap-s.cap)
		}
		return
	}

	// dst is nil, assign new slice
	s.ptr = reflect_unsafe_NewArray(ReflectTypeToABIType(dst.Type().Elem()), cap)
	s.len = len
	s.cap = cap
}

//go:nosplit
func ReflectSetSliceAt(dst reflect.Value, index int, value reflect.Value) {
	if dst.Kind() != reflect.Slice {
		panic("gointernals.ReflectSetSliceAt of non-slice type")
	}
	s := (*Slice)(ReflectValueData(dst))
	if index < 0 || index >= s.len {
		panic("gointernals.ReflectSetSliceAt: index out of range")
	}

	elemT := ReflectTypeToABIType(dst.Type().Elem())
	valT := ReflectTypeToABIType(value.Type())
	if elemT.Size != valT.Size {
		panic("gointernals.ReflectSetSliceAt: element type size mismatch")
	}
	typedmemmove(elemT, unsafe.Pointer(uintptr(s.ptr)+uintptr(index)*elemT.Size), ReflectValueData(value))
}
