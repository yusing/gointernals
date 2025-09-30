//go:build go1.24 && go1.25 && !go1.26

package gointernals

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestMakeSlice(t *testing.T) {
	t.Run("basic int slice", func(t *testing.T) {
		typ := reflect.TypeOf([]int{})
		s := ReflectMakeSlice(typ, 3, 5)

		if s.len != 3 {
			t.Errorf("Expected len 3, got %d", s.len)
		}
		if s.cap != 5 {
			t.Errorf("Expected cap 5, got %d", s.cap)
		}
		if s.ptr == nil {
			t.Error("Expected non-nil ptr")
		}

		result := *(*[]int)(unsafe.Pointer(&s))
		if len(result) != 3 {
			t.Errorf("Expected length 3, got %d", len(result))
		}
		if cap(result) != 5 {
			t.Errorf("Expected capacity 5, got %d", cap(result))
		}
	})

	t.Run("string slice", func(t *testing.T) {
		typ := reflect.TypeOf([]string{})
		s := ReflectMakeSlice(typ, 2, 4)

		if s.len != 2 {
			t.Errorf("Expected len 2, got %d", s.len)
		}
		if s.cap != 4 {
			t.Errorf("Expected cap 4, got %d", s.cap)
		}
		if s.ptr == nil {
			t.Error("Expected non-nil ptr")
		}

		result := *(*[]string)(unsafe.Pointer(&s))
		if len(result) != 2 {
			t.Errorf("Expected length 2, got %d", len(result))
		}
		if cap(result) != 4 {
			t.Errorf("Expected capacity 4, got %d", cap(result))
		}
	})

	t.Run("struct slice", func(t *testing.T) {
		type TestStruct struct {
			Name string
			ID   int
		}
		typ := reflect.TypeOf([]TestStruct{})
		s := ReflectMakeSlice(typ, 1, 3)

		if s.len != 1 {
			t.Errorf("Expected len 1, got %d", s.len)
		}
		if s.cap != 3 {
			t.Errorf("Expected cap 3, got %d", s.cap)
		}
		if s.ptr == nil {
			t.Error("Expected non-nil ptr")
		}

		result := *(*[]TestStruct)(unsafe.Pointer(&s))
		if len(result) != 1 {
			t.Errorf("Expected length 1, got %d", len(result))
		}
		if cap(result) != 3 {
			t.Errorf("Expected capacity 3, got %d", cap(result))
		}
	})

	t.Run("pointer slice", func(t *testing.T) {
		typ := reflect.TypeOf([]*int{})
		s := ReflectMakeSlice(typ, 0, 5)

		if s.len != 0 {
			t.Errorf("Expected len 0, got %d", s.len)
		}
		if s.cap != 5 {
			t.Errorf("Expected cap 5, got %d", s.cap)
		}
		if s.ptr == nil {
			t.Error("Expected non-nil ptr")
		}

		result := *(*[]*int)(unsafe.Pointer(&s))
		if len(result) != 0 {
			t.Errorf("Expected length 0, got %d", len(result))
		}
		if cap(result) != 5 {
			t.Errorf("Expected capacity 5, got %d", cap(result))
		}
	})

	t.Run("zero length and capacity", func(t *testing.T) {
		typ := reflect.TypeOf([]int{})
		s := ReflectMakeSlice(typ, 0, 0)

		if s.len != 0 {
			t.Errorf("Expected len 0, got %d", s.len)
		}
		if s.cap != 0 {
			t.Errorf("Expected cap 0, got %d", s.cap)
		}

		result := *(*[]int)(unsafe.Pointer(&s))
		if len(result) != 0 {
			t.Errorf("Expected length 0, got %d", len(result))
		}
		if cap(result) != 0 {
			t.Errorf("Expected capacity 0, got %d", cap(result))
		}
	})

	t.Run("len equals cap", func(t *testing.T) {
		typ := reflect.TypeOf([]byte{})
		s := ReflectMakeSlice(typ, 10, 10)

		if s.len != 10 {
			t.Errorf("Expected len 10, got %d", s.len)
		}
		if s.cap != 10 {
			t.Errorf("Expected cap 10, got %d", s.cap)
		}

		result := *(*[]byte)(unsafe.Pointer(&s))
		if len(result) != 10 {
			t.Errorf("Expected length 10, got %d", len(result))
		}
		if cap(result) != 10 {
			t.Errorf("Expected capacity 10, got %d", cap(result))
		}
	})

	t.Run("panic on non-slice type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for non-slice type")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectMakeSlice of non-slice type" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		typ := reflect.TypeOf(42)
		ReflectMakeSlice(typ, 1, 1)
	})

	t.Run("panic on negative len", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for negative len")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectMakeSlice: negative len" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		typ := reflect.TypeOf([]int{})
		ReflectMakeSlice(typ, -1, 5)
	})

	t.Run("panic on negative cap", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for negative cap")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectMakeSlice: negative cap" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		typ := reflect.TypeOf([]int{})
		ReflectMakeSlice(typ, 1, -1)
	})

	t.Run("panic on len > cap", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for len > cap")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectMakeSlice: len > cap" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		typ := reflect.TypeOf([]int{})
		ReflectMakeSlice(typ, 10, 5)
	})
}

func TestInitSlice(t *testing.T) {
	t.Run("nil slice to non-nil", func(t *testing.T) {
		var s []int
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 3, 5)

		if len(s) != 3 {
			t.Errorf("Expected len 3, got %d", len(s))
		}
		if cap(s) != 5 {
			t.Errorf("Expected cap 5, got %d", cap(s))
		}
	})

	t.Run("string slice initialization", func(t *testing.T) {
		var s []string
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 2, 4)

		if len(s) != 2 {
			t.Errorf("Expected len 2, got %d", len(s))
		}
		if cap(s) != 4 {
			t.Errorf("Expected cap 4, got %d", cap(s))
		}
	})

	t.Run("grow existing slice", func(t *testing.T) {
		s := make([]int, 2, 3)
		s[0] = 10
		s[1] = 20
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 5, 10)

		if len(s) != 5 {
			t.Errorf("Expected len 5, got %d", len(s))
		}
		if cap(s) < 10 {
			t.Errorf("Expected cap at least 10, got %d", cap(s))
		}
	})

	t.Run("set length without reallocation", func(t *testing.T) {
		s := make([]int, 2, 10)
		s[0] = 100
		s[1] = 200
		oldPtr := unsafe.Pointer(&s[0])
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 5, 7)

		if len(s) != 5 {
			t.Errorf("Expected len 5, got %d", len(s))
		}
		if cap(s) != 10 {
			t.Errorf("Expected cap 10, got %d", cap(s))
		}
		newPtr := unsafe.Pointer(&s[0])
		if oldPtr != newPtr {
			t.Error("Expected slice to not reallocate when capacity is sufficient")
		}
	})

	t.Run("reduce length without reallocation", func(t *testing.T) {
		s := make([]int, 5, 10)
		for i := range 5 {
			s[i] = i * 10
		}
		oldPtr := unsafe.Pointer(&s[0])
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 3, 8)

		if len(s) != 3 {
			t.Errorf("Expected len 3, got %d", len(s))
		}
		if cap(s) != 10 {
			t.Errorf("Expected cap 10, got %d", cap(s))
		}
		newPtr := unsafe.Pointer(&s[0])
		if oldPtr != newPtr {
			t.Error("Expected slice to not reallocate when capacity is sufficient")
		}
	})

	t.Run("struct slice", func(t *testing.T) {
		type TestStruct struct {
			Name string
			ID   int
		}
		var s []TestStruct
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 2, 3)

		if len(s) != 2 {
			t.Errorf("Expected len 2, got %d", len(s))
		}
		if cap(s) != 3 {
			t.Errorf("Expected cap 3, got %d", cap(s))
		}
	})

	t.Run("pointer slice", func(t *testing.T) {
		var s []*int
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 1, 5)

		if len(s) != 1 {
			t.Errorf("Expected len 1, got %d", len(s))
		}
		if cap(s) != 5 {
			t.Errorf("Expected cap 5, got %d", cap(s))
		}
	})

	t.Run("zero length and capacity", func(t *testing.T) {
		var s []int
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 0, 0)

		if len(s) != 0 {
			t.Errorf("Expected len 0, got %d", len(s))
		}
		if cap(s) != 0 {
			t.Errorf("Expected cap 0, got %d", cap(s))
		}
	})

	t.Run("len equals cap", func(t *testing.T) {
		var s []byte
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 10, 10)

		if len(s) != 10 {
			t.Errorf("Expected len 10, got %d", len(s))
		}
		if cap(s) != 10 {
			t.Errorf("Expected cap 10, got %d", cap(s))
		}
	})

	t.Run("panic on non-slice type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for non-slice type")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectInitSlice of non-slice type" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		i := 42
		v := reflect.ValueOf(&i).Elem()
		ReflectInitSlice(v, 1, 1)
	})

	t.Run("initialize string slice from existing", func(t *testing.T) {
		s := []string{"hello", "world"}
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 4, 10)

		if len(s) != 4 {
			t.Errorf("Expected len 4, got %d", len(s))
		}
		if cap(s) < 10 {
			t.Errorf("Expected cap at least 10, got %d", cap(s))
		}
	})

	t.Run("existing slice with exact capacity needed", func(t *testing.T) {
		s := make([]int, 2, 5)
		s[0] = 1
		s[1] = 2
		oldPtr := unsafe.Pointer(&s[0])
		v := reflect.ValueOf(&s).Elem()

		ReflectInitSlice(v, 3, 5)

		if len(s) != 3 {
			t.Errorf("Expected len 3, got %d", len(s))
		}
		if cap(s) != 5 {
			t.Errorf("Expected cap 5, got %d", cap(s))
		}
		newPtr := unsafe.Pointer(&s[0])
		if oldPtr != newPtr {
			t.Error("Expected slice to not reallocate when capacity is sufficient")
		}
	})
}

func TestSetSliceAt(t *testing.T) {
	t.Run("set int slice element", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 2, reflect.ValueOf(100))

		if s[2] != 100 {
			t.Errorf("Expected s[2] to be 100, got %d", s[2])
		}
		if s[0] != 1 || s[1] != 2 || s[3] != 4 || s[4] != 5 {
			t.Error("Other elements should remain unchanged")
		}
	})

	t.Run("set string slice element", func(t *testing.T) {
		s := []string{"a", "b", "c"}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 1, reflect.ValueOf("hello"))

		if s[1] != "hello" {
			t.Errorf("Expected s[1] to be 'hello', got %s", s[1])
		}
		if s[0] != "a" || s[2] != "c" {
			t.Error("Other elements should remain unchanged")
		}
	})

	t.Run("set struct slice element", func(t *testing.T) {
		type TestStruct struct {
			Name string
			ID   int
		}
		s := []TestStruct{
			{Name: "first", ID: 1},
			{Name: "second", ID: 2},
			{Name: "third", ID: 3},
		}
		v := reflect.ValueOf(&s).Elem()

		newVal := TestStruct{Name: "updated", ID: 99}
		ReflectSetSliceAt(v, 1, reflect.ValueOf(newVal))

		if s[1].Name != "updated" || s[1].ID != 99 {
			t.Errorf("Expected s[1] to be {updated, 99}, got %+v", s[1])
		}
		if s[0].Name != "first" || s[2].Name != "third" {
			t.Error("Other elements should remain unchanged")
		}
	})

	t.Run("set pointer slice element", func(t *testing.T) {
		a, b, c := 1, 2, 3
		newVal := 100
		s := []*int{&a, &b, &c}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 0, reflect.ValueOf(&newVal))

		if s[0] != &newVal || *s[0] != 100 {
			t.Errorf("Expected s[0] to point to newVal (100), got %d", *s[0])
		}
		if s[1] != &b || s[2] != &c {
			t.Error("Other elements should remain unchanged")
		}
	})

	t.Run("set byte slice element", func(t *testing.T) {
		s := []byte{0, 1, 2, 3}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 3, reflect.ValueOf(byte(255)))

		if s[3] != 255 {
			t.Errorf("Expected s[3] to be 255, got %d", s[3])
		}
		if s[0] != 0 || s[1] != 1 || s[2] != 2 {
			t.Error("Other elements should remain unchanged")
		}
	})

	t.Run("set first element", func(t *testing.T) {
		s := []int{10, 20, 30}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 0, reflect.ValueOf(999))

		if s[0] != 999 {
			t.Errorf("Expected s[0] to be 999, got %d", s[0])
		}
	})

	t.Run("set last element", func(t *testing.T) {
		s := []int{10, 20, 30}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 2, reflect.ValueOf(999))

		if s[2] != 999 {
			t.Errorf("Expected s[2] to be 999, got %d", s[2])
		}
	})

	t.Run("set float64 slice element", func(t *testing.T) {
		s := []float64{1.1, 2.2, 3.3}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 1, reflect.ValueOf(9.9))

		if s[1] != 9.9 {
			t.Errorf("Expected s[1] to be 9.9, got %f", s[1])
		}
	})

	t.Run("set bool slice element", func(t *testing.T) {
		s := []bool{true, false, true}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 1, reflect.ValueOf(true))

		if s[1] != true {
			t.Errorf("Expected s[1] to be true, got %t", s[1])
		}
	})

	t.Run("set complex slice element", func(t *testing.T) {
		s := []complex128{1 + 2i, 3 + 4i}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 0, reflect.ValueOf(5+6i))

		if s[0] != 5+6i {
			t.Errorf("Expected s[0] to be 5+6i, got %v", s[0])
		}
	})

	t.Run("set interface slice element", func(t *testing.T) {
		s := []any{1, "hello", 3.14}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 1, reflect.ValueOf("world"))

		if s[1] != "world" {
			t.Errorf("Expected s[1] to be 'world', got %v", s[1])
		}
	})

	t.Run("panic on non-slice type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for non-slice type")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectSetSliceAt of non-slice type" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		i := 42
		v := reflect.ValueOf(&i).Elem()
		ReflectSetSliceAt(v, 0, reflect.ValueOf(100))
	})

	t.Run("panic on negative index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for negative index")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectSetSliceAt: index out of range" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		s := []int{1, 2, 3}
		v := reflect.ValueOf(&s).Elem()
		ReflectSetSliceAt(v, -1, reflect.ValueOf(100))
	})

	t.Run("panic on index out of range high", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for index out of range")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectSetSliceAt: index out of range" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		s := []int{1, 2, 3}
		v := reflect.ValueOf(&s).Elem()
		ReflectSetSliceAt(v, 3, reflect.ValueOf(100))
	})

	t.Run("panic on index equal to length", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for index equal to length")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectSetSliceAt: index out of range" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		s := []int{1, 2, 3}
		v := reflect.ValueOf(&s).Elem()
		ReflectSetSliceAt(v, 3, reflect.ValueOf(100))
	})

	t.Run("panic on element type size mismatch", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for element type size mismatch")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectSetSliceAt: element type size mismatch" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		s := []int64{1, 2, 3}
		v := reflect.ValueOf(&s).Elem()
		ReflectSetSliceAt(v, 0, reflect.ValueOf(int32(100)))
	})

	t.Run("set element in single element slice", func(t *testing.T) {
		s := []int{42}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 0, reflect.ValueOf(999))

		if s[0] != 999 {
			t.Errorf("Expected s[0] to be 999, got %d", s[0])
		}
	})

	t.Run("set nested struct slice element", func(t *testing.T) {
		type Inner struct {
			Value int
		}
		type Outer struct {
			Data Inner
			Name string
		}
		s := []Outer{
			{Data: Inner{Value: 1}, Name: "first"},
			{Data: Inner{Value: 2}, Name: "second"},
		}
		v := reflect.ValueOf(&s).Elem()

		newVal := Outer{Data: Inner{Value: 99}, Name: "updated"}
		ReflectSetSliceAt(v, 0, reflect.ValueOf(newVal))

		if s[0].Data.Value != 99 || s[0].Name != "updated" {
			t.Errorf("Expected s[0] to be updated, got %+v", s[0])
		}
		if s[1].Data.Value != 2 || s[1].Name != "second" {
			t.Error("Other elements should remain unchanged")
		}
	})

	t.Run("set array slice element", func(t *testing.T) {
		s := [][3]int{{1, 2, 3}, {4, 5, 6}}
		v := reflect.ValueOf(&s).Elem()

		ReflectSetSliceAt(v, 1, reflect.ValueOf([3]int{7, 8, 9}))

		if s[1] != [3]int{7, 8, 9} {
			t.Errorf("Expected s[1] to be [7, 8, 9], got %v", s[1])
		}
		if s[0] != [3]int{1, 2, 3} {
			t.Error("Other elements should remain unchanged")
		}
	})

	t.Run("set map slice element", func(t *testing.T) {
		s := []map[string]int{
			{"a": 1},
			{"b": 2},
		}
		v := reflect.ValueOf(&s).Elem()

		newMap := map[string]int{"c": 3, "d": 4}
		ReflectSetSliceAt(v, 0, reflect.ValueOf(newMap))

		if s[0]["c"] != 3 || s[0]["d"] != 4 {
			t.Errorf("Expected s[0] to be updated map, got %v", s[0])
		}
	})
}
