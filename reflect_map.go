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

func ReflectStrMapAssign(dst reflect.Value, key string) reflect.Value {
	eface := EfaceOf(dst.Interface())
	m, mType := (*Map)(unsafe.Pointer(eface.Data)), (*MapType)(abi.NoEscape(unsafe.Pointer(eface.Type)))
	elemPtr := mapassign_faststr(mType, m, key)
	return reflect.NewAt(dst.Type().Elem(), elemPtr)
}
