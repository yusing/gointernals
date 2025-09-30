package gointernals

import (
	"reflect"
	"unsafe"
)

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
		ptr: reflect_unsafe_NewArray(EfaceOf(typ.Elem()).Type, cap),
		len: len,
		cap: cap,
	}
}

// ReflectInitSlice initializes a slice with the given length and capacity.
func ReflectInitSlice(dst reflect.Value, len, cap int) {
	if dst.Kind() != reflect.Slice {
		panic("gointernals.ReflectInitSlice of non-slice type")
	}

	s := (*Slice)(ReflectValueData(dst))
	if s.ptr != nil {
		s.len = len
		if s.cap < cap {
			*s = reflect_growslice(EfaceOf(dst.Type().Elem()).Type, *s, cap-s.cap)
		}
		return
	}

	// dst is nil, assign new slice
	s.ptr = reflect_unsafe_NewArray(EfaceOf(dst.Type().Elem()).Type, cap)
	s.len = len
	s.cap = cap
}

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
