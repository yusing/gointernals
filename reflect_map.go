package gointernals

import (
	"reflect"

	_ "unsafe"
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
