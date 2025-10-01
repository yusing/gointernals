//go:build go1.24 && go1.25 && !go1.26

package gointernals

import (
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

// Map constants common to several packages
// runtime/runtime-gdb.py:MapTypePrinter contains its own copy
const (
	// Number of bits in the group.slot count.
	MapGroupSlotsBits = 3

	// Number of slots in a group.
	MapGroupSlots = 1 << MapGroupSlotsBits // 8

	// Maximum key or elem size to keep inline (instead of mallocing per element).
	// Must fit in a uint8.
	MapMaxKeyBytes  = 128
	MapMaxElemBytes = 128

	ctrlEmpty = 0b10000000
	bitsetLSB = 0x0101010101010101

	// Value of control word with all empty slots.
	MapCtrlEmpty = bitsetLSB * uint64(ctrlEmpty)
)

// Hmap is the header for a Go map
type Hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
	// Make sure this stays in sync with the compiler's definition.
	Count     int // # live cells == size of map.  Must be first (used by len() builtin)
	Flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	Noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	Hash0     uint32 // hash seed

	Buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	Oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	Nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	Extra *MapExtra // optional fields
}

type MapExtra struct {
	Overflow     unsafe.Pointer
	Oldoverflow  unsafe.Pointer
	NextOverflow unsafe.Pointer
}

type MapType struct {
	abi.Type
	Key   *abi.Type
	Elem  *abi.Type
	Group *abi.Type // internal type representing a slot group
	// function for hashing keys (ptr to key, seed) -> hash
	Hasher    func(unsafe.Pointer, uintptr) uintptr
	GroupSize uintptr // == Group.Size_
	SlotSize  uintptr // size of key/elem slot
	ElemOff   uintptr // offset of elem in key/elem slot
	Flags     uint32
}

// Flag values
const (
	MapNeedKeyUpdate = 1 << iota
	MapHashMightPanic
	MapIndirectKey
	MapIndirectElem
)

func (mt *MapType) NeedKeyUpdate() bool { // true if we need to update key on an overwrite
	return mt.Flags&MapNeedKeyUpdate != 0
}
func (mt *MapType) HashMightPanic() bool { // true if hash function might panic
	return mt.Flags&MapHashMightPanic != 0
}
func (mt *MapType) IndirectKey() bool { // store ptr to key instead of key itself
	return mt.Flags&MapIndirectKey != 0
}
func (mt *MapType) IndirectElem() bool { // store ptr to elem instead of elem itself
	return mt.Flags&MapIndirectElem != 0
}

type Map struct {
	// The number of filled slots (i.e. the number of elements in all
	// tables). Excludes deleted slots.
	// Must be first (known by the compiler, for len() builtin).
	used uint64

	// seed is the hash seed, computed as a unique random number per map.
	seed uintptr

	// The directory of tables.
	//
	// Normally dirPtr points to an array of table pointers
	//
	// dirPtr *[dirLen]*table
	//
	// The length (dirLen) of this array is `1 << globalDepth`. Multiple
	// entries may point to the same table. See top-level comment for more
	// details.
	//
	// Small map optimization: if the map always contained
	// abi.MapGroupSlots or fewer entries, it fits entirely in a
	// single group. In that case dirPtr points directly to a single group.
	//
	// dirPtr *group
	//
	// In this case, dirLen is 0. used counts the number of used slots in
	// the group. Note that small maps never have deleted slots (as there
	// is no probe sequence to maintain).
	dirPtr unsafe.Pointer
	dirLen int

	// The number of bits to use in table directory lookups.
	globalDepth uint8

	// The number of bits to shift out of the hash for directory lookups.
	// On 64-bit systems, this is 64 - globalDepth.
	globalShift uint8

	// writing is a flag that is toggled (XOR 1) while the map is being
	// written. Normally it is set to 1 when writing, but if there are
	// multiple concurrent writers, then toggling increases the probability
	// that both sides will detect the race.
	writing uint8

	// tombstonePossible is false if we know that no table in this map
	// contains a tombstone.
	tombstonePossible bool

	// clearSeq is a sequence counter of calls to Clear. It is used to
	// detect map clears during iteration.
	clearSeq uint64
}

// MapTable is a Swiss MapTable hash MapTable structure.
//
// Each MapTable is a complete hash MapTable implementation.
//
// Map uses one or more tables to store entries. Extendible hashing (hash
// prefix) is used to select the MapTable to use for a specific key. Using
// multiple tables enables incremental growth by growing only one MapTable at a
// time.
type MapTable struct {
	// The number of filled slots (i.e. the number of elements in the table).
	used uint16

	// The total number of slots (always 2^N). Equal to
	// `(groups.lengthMask+1)*abi.MapGroupSlots`.
	capacity uint16

	// The number of slots we can still fill without needing to rehash.
	//
	// We rehash when used + tombstones > loadFactor*capacity, including
	// tombstones so the table doesn't overfill with tombstones. This field
	// counts down remaining empty slots before the next rehash.
	growthLeft uint16

	// The number of bits used by directory lookups above this table. Note
	// that this may be less then globalDepth, if the directory has grown
	// but this table has not yet been split.
	localDepth uint8

	// Index of this table in the Map directory. This is the index of the
	// _first_ location in the directory. The table may occur in multiple
	// sequential indices.
	//
	// index is -1 if the table is stale (no longer installed in the
	// directory).
	index int

	// groups is an array of slot groups. Each group holds abi.MapGroupSlots
	// key/elem slots and their control bytes. A table has a fixed size
	// groups array. The table is replaced (in rehash) when more space is
	// required.
	//
	// TODO(prattmic): keys and elements are interleaved to maximize
	// locality, but it comes at the expense of wasted space for some types
	// (consider uint8 key, uint64 element). Consider placing all keys
	// together in these cases to save space.
	groups unsafe.Pointer
}

//go:linkname mapclone maps.clone
//go:noescape
func mapclone(m any) any

//go:linkname mapaccess1_faststr runtime.mapaccess1_faststr
//go:noescape
func mapaccess1_faststr(t *MapType, m *Map, ky string) unsafe.Pointer

//go:linkname mapaccess2_faststr
//go:noescape
func mapaccess2_faststr(t *MapType, m *Map, ky string) (unsafe.Pointer, bool)

//go:linkname mapassign_faststr runtime.mapassign_faststr
//go:noescape
func mapassign_faststr(t *MapType, m *Map, s string) unsafe.Pointer

//go:linkname mapassign runtime.mapassign
//go:noescape
func mapassign(t *MapType, m *Map, key unsafe.Pointer) unsafe.Pointer

type (
	StrMapGetFunc = func(m *Map, mType *MapType, key string) unsafe.Pointer
	StrMapSetFunc = func(m *Map, mType *MapType, key string, value unsafe.Pointer)
)

//go:nosplit
//go:linkname StrMapGet gointernals.StrMapGet
func StrMapGet(m *Map, mType *MapType, key string) unsafe.Pointer {
	return mapaccess1_faststr(mType, m, key)
}

//go:nosplit
func StrMapGetAs[K comparable, V any](m *Map, mType *MapType, key string) V {
	return *(*V)(StrMapGet(m, mType, key))
}

//go:nosplit
func StrMapTryGet(m *Map, mType *MapType, key string) (unsafe.Pointer, bool) {
	return mapaccess2_faststr(mType, m, key)
}

//go:nosplit
func StrMapTryGetAs[K comparable, V any](m *Map, mType *MapType, key string) (V, bool) {
	v, ok := StrMapTryGet(m, mType, key)
	if !ok || v == nil {
		var zero V
		return zero, ok
	}
	return *(*V)(v), ok
}

//go:nosplit
//go:linkname StrMapSet gointernals.StrMapSet
func StrMapSet(m *Map, mType *MapType, key string, value unsafe.Pointer) {
	dst := mapassign_faststr(mType, m, key)
	typedmemmove(mType.Elem, dst, value)
}

//go:nosplit
//go:linkname MapClone gointernals.MapClone
func MapClone(m *Map, mType *MapType) any {
	return mapclone(MapToAny(m, mType))
}

//go:nosplit
func MapCloneAs[K comparable, V any](m *Map, mType *MapType) map[K]V {
	return mapclone(MapToAny(m, mType)).(map[K]V)
}

//go:nosplit
//go:linkname MapClear gointernals.MapClear
func MapClear(m *Map, mType *MapType) {
	reflect_mapclear(mType, m)
}

//go:nosplit
func MapUnpack[K comparable, V any](m map[K]V) (*Map, *MapType) {
	eface := EfaceOf(m)
	return (*Map)(unsafe.Pointer(eface.Data)), (*MapType)(abi.NoEscape(unsafe.Pointer(eface.Type)))
}

//go:nosplit
func StrMapCast[K ~string, V any, M ~map[K]V](m M) map[string]V {
	return *(*map[string]V)(unsafe.Pointer(&m))
}

//go:nosplit
//go:linkname MapElemType gointernals.MapElemType
func MapElemType(mType *MapType) *abi.Type {
	return mType.Elem
}

func MapToEface(m *Map, mType *MapType) *abi.Eface {
	return &abi.Eface{
		Data: unsafe.Pointer(m),
		Type: (*abi.Type)(abi.NoEscape(unsafe.Pointer(mType))),
	}
}

func MapToAny(m *Map, mType *MapType) any {
	return AnyFrom(MapToEface(m, mType))
}
