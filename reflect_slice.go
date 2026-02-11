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

// reflectFlagIndir is the flagIndir bit in reflect.Value.flag.
// When set, the ptr field points to the data; when unset, ptr IS the data.
const reflectFlagIndir = uintptr(1 << 7)

// reflectValueDataPtr returns a pointer to the actual data held by v.
// Unlike ReflectValueData, this correctly handles non-indirect values
// (pointers, maps, channels, funcs) where reflect.Value stores the
// value directly in the ptr field rather than as a pointer to the data.
//
//go:nosplit
func reflectValueDataPtr(v *reflect.Value) unsafe.Pointer {
	// reflect.Value layout: [typ *abi.Type, ptr unsafe.Pointer, flag uintptr]
	words := (*[3]uintptr)(unsafe.Pointer(v))
	if words[2]&reflectFlagIndir != 0 {
		// Indirect: ptr points to the data
		return unsafe.Pointer(words[1])
	}
	// Direct: ptr IS the data, return address of the ptr field
	return unsafe.Pointer(&words[1])
}

func ReflectSetSliceAt(dst reflect.Value, index int, value reflect.Value) {
	if dst.Kind() != reflect.Slice {
		panic("gointernals.ReflectSetSliceAt of non-slice type")
	}
	s := (*Slice)(ReflectValueData(dst))
	if index < 0 || index >= s.len {
		panic("gointernals.ReflectSetSliceAt: index out of range")
	}

	elemType := dst.Type().Elem()
	elemDst := unsafe.Add(s.ptr, uintptr(index)*elemType.Size())

	// Interface element types need special handling:
	// the value's concrete type layout differs from eface/iface layout,
	// so we delegate to reflect which handles the conversion correctly.
	if elemType.Kind() == reflect.Interface {
		reflect.NewAt(elemType, elemDst).Elem().Set(value)
		return
	}

	if elemType.Size() != value.Type().Size() {
		panic("gointernals.ReflectSetSliceAt: element type size mismatch")
	}
	// Use reflectValueDataPtr to correctly handle both indirect and
	// direct (pointer-like) values in reflect.Value.
	typedmemmove(ReflectTypeToABIType(elemType), elemDst, reflectValueDataPtr(&value))
}
