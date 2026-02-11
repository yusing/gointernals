//go:build go1.24 && !go1.27

package gointernals

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

//go:noescape
//go:linkname reflect_unsafe_New reflect.unsafe_New
func reflect_unsafe_New(*abi.Type) unsafe.Pointer

//go:noescape
//go:linkname reflect_unsafe_NewArray reflect.unsafe_NewArray
func reflect_unsafe_NewArray(*abi.Type, int) unsafe.Pointer

//go:noescape
//go:linkname reflect_growslice reflect.growslice
func reflect_growslice(t *abi.Type, old Slice, num int) Slice

//go:nosplit
func ReflectTypeToABIType(t reflect.Type) *abi.Type {
	return (*abi.Type)(abi.NoEscape(EfaceOf(t).Data))
}

//go:nosplit
func ReflectValueType(v reflect.Value) *abi.Type {
	return *(**abi.Type)(abi.NoEscape(unsafe.Pointer(&v)))
}

//go:nosplit
func ReflectValueData(v reflect.Value) unsafe.Pointer {
	return (*[2]unsafe.Pointer)(abi.NoEscape(unsafe.Pointer(&v)))[1]
}

//go:nosplit
func ReflectValueSet[T any](v reflect.Value, x T) {
	*(*T)(abi.NoEscape(ReflectValueData(v))) = x
}

//go:nosplit
func ReflectValueAs[T any](v reflect.Value) T {
	return *(*T)(abi.NoEscape(ReflectValueData(v)))
}

//go:nosplit
func ReflectInitPtr(v reflect.Value) {
	t := v.Type()
	switch t.Kind() {
	case reflect.Pointer:
		elemT := t.Elem()
		newElemPtr := reflect_unsafe_New(ReflectTypeToABIType(elemT))
		ReflectValueSet(v, newElemPtr)
		return
	}

	panic(fmt.Errorf("gointernals.ReflectInitPtr: invalid type %s", t.Kind().String()))
}

// ReflectShallowCopy copies the value of src to dst.
// It will panic if the types are not compatible.
// Note: this function does not update type and flag fields in dst.
//
// FIXME: segfault
//
//go:nosplit
func ReflectShallowCopy(dst, src reflect.Value) {
	dstT, srcT := ReflectValueType(dst), ReflectValueType(src)
	// pointer to same type will always be the same pointer
	if dstT == srcT || dst.Type().AssignableTo(src.Type()) {
		typedmemmove(dstT, ReflectValueData(dst), ReflectValueData(src))
		return
	}

	if src.Type().ConvertibleTo(dst.Type()) {
		srcConv := src.Convert(dst.Type())
		typedmemmove(dstT, ReflectValueData(dst), ReflectValueData(srcConv))
		return
	}

	panic(fmt.Errorf("gointernals.ReflectShallowCopy: invalid shallow copy from %s to %s", dstT.Kind().String(), srcT.Kind().String()))
}

func ReflectIsNumeric(v reflect.Value) bool {
	return v.Kind() >= reflect.Int && v.Kind() <= reflect.Float64
}

func ReflectCanInt(v reflect.Value) bool {
	return v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64
}

func ReflectCanUint(v reflect.Value) bool {
	return v.Kind() >= reflect.Uint && v.Kind() <= reflect.Uint64
}

func ReflectCanFloat(v reflect.Value) bool {
	return v.Kind() >= reflect.Float32 && v.Kind() <= reflect.Float64
}
