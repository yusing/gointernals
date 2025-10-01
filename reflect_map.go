package gointernals

import (
	"reflect"
	"unsafe"

	_ "unsafe"

	"github.com/yusing/gointernals/abi"
)

//go:linkname reflect_makemap reflect.makemap
func reflect_makemap(t *MapType, cap int) *Map

//go:linkname reflect_mapclear reflect.mapclear
//go:noescape
func reflect_mapclear(t *MapType, m *Map)

//go:nosplit
func ReflectInitMap(dst reflect.Value, len int) {
	if dst.Kind() != reflect.Map {
		panic("gointernals.ReflectInitMap of non-map type")
	}

	dstT := ReflectTypeToABIType(dst.Type())
	newMap := reflect_makemap(PointerCast[MapType](dstT), len)
	ReflectValueSet(dst, newMap)
}

//go:nosplit
func ReflectMapUnpack(dst reflect.Value) (*Map, *MapType) {
	eface := EfaceOf(dst.Interface())
	return (*Map)(unsafe.Pointer(eface.Data)), (*MapType)(abi.NoEscape(unsafe.Pointer(eface.Type)))
}

// ReflectStrMapAssign assigns a string key to a map and returns the value.
//
// The returned value should satisfy CanSet().
//
//go:nosplit
func ReflectStrMapAssign(dst reflect.Value, key string) reflect.Value {
	if dst.Kind() != reflect.Map || dst.Type().Key().Kind() != reflect.String {
		panic("gointernals.ReflectStrMapAssign of non map or non-string map type")
	}
	if dst.IsNil() {
		panic("gointernals.ReflectStrMapAssign of nil map")
	}

	m, mType := ReflectMapUnpack(dst)
	elemPtr := mapassign_faststr(mType, m, key)
	return reflect.NewAt(dst.Type().Elem(), elemPtr).Elem()
}

// ReflectMapAssign assigns a key to a map and returns the value.
//
// The returned value should satisfy CanSet().
//
//go:nosplit
func ReflectMapAssign(dst reflect.Value, key any) reflect.Value {
	if dst.Kind() != reflect.Map {
		panic("gointernals.ReflectMapAssign of non map type")
	}
	if dst.IsNil() {
		panic("gointernals.ReflectMapAssign of nil map")
	}

	// fast path (same type)
	m, mType := ReflectMapUnpack(dst)
	keyType := dst.Type().Key()
	if keyType.Kind() == dst.Type().Elem().Kind() {
		elemPtr := mapassign(mType, m, EfaceOf(key).Data)
		return reflect.NewAt(dst.Type().Elem(), elemPtr).Elem()
	}

	// slow path (any / interface key)
	keyVal := reflect.ValueOf(key)
	// Ensure key is of the exact map key type
	if !keyVal.Type().AssignableTo(keyType) {
		if keyVal.Type().ConvertibleTo(keyType) {
			keyVal = keyVal.Convert(keyType)
		} else {
			panic("gointernals.ReflectMapAssign key not assignable to map key type")
		}
	}
	keyPtr := reflect.New(keyType)
	keyPtr.Elem().Set(keyVal)
	elemPtr := mapassign(mType, m, keyPtr.UnsafePointer())
	return reflect.NewAt(dst.Type().Elem(), elemPtr).Elem()
}
