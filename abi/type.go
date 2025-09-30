//go:build go1.24 && go1.25 && !go1.26

package abi

import "unsafe"

// Type represents basic type information (from internal/abi/type.go)
type Type struct {
	Size       uintptr
	PtrBytes   uintptr // number of (prefix) bytes in the type that can contain pointers
	Hash       uint32  // hash of type; avoids computation in hash tables
	TFlag      TFlag   // extra type information flags
	Align      uint8   // alignment of variable with this type
	FieldAlign uint8   // alignment of struct field with this type
	Kind_      Kind    // what kind of type this is (string, int, ...)
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
	// GCData stores the GC type data for the garbage collector.
	// Normally, GCData points to a bitmask that describes the
	// ptr/nonptr fields of the type. The bitmask will have at
	// least PtrBytes/ptrSize bits.
	// If the TFlagGCMaskOnDemand bit is set, GCData is instead a
	// **byte and the pointer to the bitmask is one dereference away.
	// The runtime will build the bitmask if needed.
	// (See runtime/type.go:getGCMask.)
	// Note: multiple types may have the same value of GCData,
	// including when TFlagGCMaskOnDemand is set. The types will, of course,
	// have the same pointer layout (but not necessarily the same size).
	GCData    *byte
	Str       NameOff // string form
	PtrToThis TypeOff // type for pointer to this type, may be zero
}

// A Kind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type Kind uint8

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Pointer
	Slice
	String
	Struct
	UnsafePointer
)

// TFlag is used by a Type to signal what extra type information is
// available in the memory directly following the Type value.
type TFlag uint8

const (
	// TFlagUncommon means that there is a data with a type, UncommonType,
	// just beyond the shared-per-type common data.  That is, the data
	// for struct types will store their UncommonType at one offset, the
	// data for interface types will store their UncommonType at a different
	// offset.  UncommonType is always accessed via a pointer that is computed
	// using trust-us-we-are-the-implementors pointer arithmetic.
	//
	// For example, if t.Kind() == Struct and t.tflag&TFlagUncommon != 0,
	// then t has UncommonType data and it can be accessed as:
	//
	//	type structTypeUncommon struct {
	//		structType
	//		u UncommonType
	//	}
	//	u := &(*structTypeUncommon)(unsafe.Pointer(t)).u
	TFlagUncommon TFlag = 1 << 0

	// TFlagExtraStar means the name in the str field has an
	// extraneous '*' prefix. This is because for most types T in
	// a program, the type *T also exists and reusing the str data
	// saves binary size.
	TFlagExtraStar TFlag = 1 << 1

	// TFlagNamed means the type has a name.
	TFlagNamed TFlag = 1 << 2

	// TFlagRegularMemory means that equal and hash functions can treat
	// this type as a single region of t.size bytes.
	TFlagRegularMemory TFlag = 1 << 3

	// TFlagGCMaskOnDemand means that the GC pointer bitmask will be
	// computed on demand at runtime instead of being precomputed at
	// compile time. If this flag is set, the GCData field effectively
	// has type **byte instead of *byte. The runtime will store a
	// pointer to the GC pointer bitmask in *GCData.
	TFlagGCMaskOnDemand TFlag = 1 << 4

	// TFlagDirectIface means that a value of this type is stored directly
	// in the data field of an interface, instead of indirectly. Normally
	// this means the type is pointer-ish.
	TFlagDirectIface TFlag = 1 << 5
)

// NameOff is the offset to a name from moduledata.types.  See resolveNameOff in runtime.
type NameOff int32

// TypeOff is the offset to a type from moduledata.types.  See resolveTypeOff in runtime.
type TypeOff int32

// TextOff is an offset from the top of a text section.  See (rtype).textOff in runtime.
type TextOff int32

// String returns the name of k.
func (k Kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return kindNames[0]
}

// Size returns the size of the kind.
func (k Kind) Size() uintptr {
	if int(k) >= len(kindSizes) {
		return 0
	}
	return kindSizes[k]
}

var kindNames = []string{
	Invalid:       "invalid",
	Bool:          "bool",
	Int:           "int",
	Int8:          "int8",
	Int16:         "int16",
	Int32:         "int32",
	Int64:         "int64",
	Uint:          "uint",
	Uint8:         "uint8",
	Uint16:        "uint16",
	Uint32:        "uint32",
	Uint64:        "uint64",
	Uintptr:       "uintptr",
	Float32:       "float32",
	Float64:       "float64",
	Complex64:     "complex64",
	Complex128:    "complex128",
	Array:         "array",
	Chan:          "chan",
	Func:          "func",
	Interface:     "interface",
	Map:           "map",
	Pointer:       "ptr",
	Slice:         "slice",
	String:        "string",
	Struct:        "struct",
	UnsafePointer: "unsafe.Pointer",
}

var kindSizes = []uintptr{
	Invalid:       0,
	Bool:          unsafe.Sizeof(bool(false)),
	Int:           unsafe.Sizeof(int(0)),
	Int8:          unsafe.Sizeof(int8(0)),
	Int16:         unsafe.Sizeof(int16(0)),
	Int32:         unsafe.Sizeof(int32(0)),
	Int64:         unsafe.Sizeof(int64(0)),
	Uint:          unsafe.Sizeof(uint(0)),
	Uint8:         unsafe.Sizeof(uint8(0)),
	Uint16:        unsafe.Sizeof(uint16(0)),
	Uint32:        unsafe.Sizeof(uint32(0)),
	Uint64:        unsafe.Sizeof(uint64(0)),
	Uintptr:       unsafe.Sizeof(uintptr(0)),
	Float32:       unsafe.Sizeof(float32(0)),
	Float64:       unsafe.Sizeof(float64(0)),
	Complex64:     unsafe.Sizeof(complex64(0)),
	Complex128:    unsafe.Sizeof(complex128(0)),
	Array:         0,
	Chan:          0,
	Func:          0,
	Interface:     0,
	Map:           0,
	Pointer:       unsafe.Sizeof(unsafe.Pointer(nil)),
	Slice:         0,
	String:        0,
	Struct:        0,
	UnsafePointer: unsafe.Sizeof(unsafe.Pointer(nil)),
}

const (
	// TODO (khr, drchase) why aren't these in TFlag?  Investigate, fix if possible.
	KindDirectIface Kind = 1 << 5
	KindMask        Kind = (1 << 5) - 1
)

// TypeOf returns the abi.Type of some value.
func TypeOf(a any) *Type {
	eface := *(*Eface)(unsafe.Pointer(&a))
	// Types are either static (for compiler-created types) or
	// heap-allocated but always reachable (for reflection-created
	// types, held in the central map). So there is no need to
	// escape types. noescape here help avoid unnecessary escape
	// of v.
	return (*Type)(NoEscape(unsafe.Pointer(eface.Type)))
}

func (t *Type) Kind() Kind {
	return t.Kind_ & KindMask
}

// CanPointer reports whether t contains pointers.
func (t *Type) CanPointer() bool {
	return t.PtrBytes != 0
}

func (t *Type) HasName() bool {
	return t.TFlag&TFlagNamed != 0
}

// IsDirectIface reports whether t is stored directly in an interface value.
func (t *Type) IsDirectIface() bool {
	return t.TFlag&TFlagDirectIface != 0
}

func (t *Type) IfaceIndir() bool {
	return t.Kind_&KindDirectIface != 0
}
