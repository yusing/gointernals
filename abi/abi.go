//go:build go1.24 && go1.25 && !go1.26

package abi

import (
	"unsafe"
)

type Eface struct {
	Type *Type
	Data unsafe.Pointer
}

// from runtime/runtime2.go
type Iface struct {
	Tab  *ITab
	Data unsafe.Pointer
}

// from internal/abi/iface.go
type ITab struct {
	Inter unsafe.Pointer
	Type  *Type
	Hash  uint32     // copy of Type.Hash. Used for type switches.
	Fun   [1]uintptr // variable sized. fun[0]==0 means Type does not implement Inter.
}
