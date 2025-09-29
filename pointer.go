package gointernals

import "github.com/yusing/gointernals/abi"

type PointerType struct {
	abi.Type
	Elem *abi.Type
}
