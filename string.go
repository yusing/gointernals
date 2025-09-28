package gointernals

import "unsafe"

type String struct {
	ptr unsafe.Pointer
	len int
}
