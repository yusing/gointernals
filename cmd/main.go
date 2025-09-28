package main

import (
	"C"
	"unsafe"

	_ "github.com/yusing/gointernals"
)

func main() {

}

//export StrMapGet
//go:linkname StrMapGet gointernals.StrMapGet
func StrMapGet(m unsafe.Pointer, mType unsafe.Pointer, key string) unsafe.Pointer

//export StrMapSet
//go:linkname StrMapSet gointernals.StrMapSet
func StrMapSet(m unsafe.Pointer, mType unsafe.Pointer, key string, value unsafe.Pointer)

//export MapClone
//go:linkname MapClone gointernals.MapClone
func MapClone(m unsafe.Pointer, mType unsafe.Pointer) unsafe.Pointer

//export MapClear
//go:linkname MapClear gointernals.MapClear
func MapClear(m unsafe.Pointer, mType unsafe.Pointer)

//export SliceClone
//go:linkname SliceClone gointernals.SliceClone
func SliceClone(src unsafe.Pointer, elemType unsafe.Pointer) unsafe.Pointer
