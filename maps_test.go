//go:build go1.24 && go1.25 && !go1.26

package gointernals

import (
	"testing"
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

func TestMapUnpack(t *testing.T) {
	t.Run("string map", func(t *testing.T) {
		m := map[string]string{
			"test_key": "test_value",
		}
		mtable, mtype := MapUnpack(m)
		if mtype == nil {
			t.Fatalf("mtype is nil")
		}
		if mtable == nil {
			t.Fatalf("mtable is nil")
		}
		if mtype.Elem.Kind() != abi.String {
			t.Errorf("Expected mtype.Elem.Kind_ to be String, got %d", mtype.Elem.Kind())
		}
	})

	t.Run("int map", func(t *testing.T) {
		m := make(map[string]int)
		mtable, mtype := MapUnpack(m)
		if mtype == nil {
			t.Error("mtype is nil")
		}
		if mtable == nil {
			t.Fatalf("mtable is nil")
		}
		if mtype.Elem.Kind() != abi.Int {
			t.Errorf("Expected mtype.Elem.Kind_ to be Int, got %d", mtype.Elem.Kind())
		}
	})

	t.Run("struct map", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}

		m := make(map[string]TestStruct)
		mtable, mtype := MapUnpack(m)
		if mtype == nil {
			t.Error("mtype is nil")
		}
		if mtable == nil {
			t.Fatalf("mtable is nil")
		}
		if mtype.Elem.Kind() != abi.Struct {
			t.Errorf("Expected mtype.Elem.Kind_ to be Struct, got %d", mtype.Elem.Kind())
		}
	})
}

func TestStrMapSet(t *testing.T) {
	// Test with string map
	t.Run("string map", func(t *testing.T) {
		m := make(map[string]string)

		mtable, mtype := MapUnpack(m)

		value := "test_value"
		StrMapSet(mtable, mtype, "test_key", unsafe.Pointer(&value))

		if val, ok := m["test_key"]; !ok {
			t.Error("Key 'test_key' was not set in map")
		} else if val != "test_value" {
			t.Errorf("Expected 'test_value', got %q", val)
		}
	})

	// Test with int map
	t.Run("int map", func(t *testing.T) {
		m := make(map[string]int)

		mtable, mtype := MapUnpack(m)

		value := 42
		StrMapSet(mtable, mtype, "int_key", unsafe.Pointer(&value))

		if val, ok := m["int_key"]; !ok {
			t.Error("Key 'int_key' was not set in map")
		} else if val != 42 {
			t.Errorf("Expected 42, got %d", val)
		}
	})

	// Test with struct map
	t.Run("struct map", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}

		m := make(map[string]TestStruct)

		mtable, mtype := MapUnpack(m)

		value := TestStruct{Name: "Alice", Age: 30}
		StrMapSet(mtable, mtype, "person", unsafe.Pointer(&value))

		if val, ok := m["person"]; !ok {
			t.Error("Key 'person' was not set in map")
		} else if val.Name != "Alice" || val.Age != 30 {
			t.Errorf("Expected {Name: Alice, Age: 30}, got %+v", val)
		}
	})

	// Test with empty key
	t.Run("empty key", func(t *testing.T) {
		m := make(map[string]string)

		mtable, mtype := MapUnpack(m)

		value := "empty_key_value"
		StrMapSet(mtable, mtype, "", unsafe.Pointer(&value))

		if val, ok := m[""]; !ok {
			t.Error("Empty key was not set in map")
		} else if val != "empty_key_value" {
			t.Errorf("Expected 'empty_key_value', got %q", val)
		}
	})

	// Test overwriting existing key
	t.Run("overwrite existing key", func(t *testing.T) {
		m := make(map[string]string)
		m["existing"] = "old_value"

		mtable, mtype := MapUnpack(m)

		value := "new_value"
		StrMapSet(mtable, mtype, "existing", unsafe.Pointer(&value))

		if val, ok := m["existing"]; !ok {
			t.Error("Key 'existing' was not found in map")
		} else if val != "new_value" {
			t.Errorf("Expected 'new_value', got %q", val)
		}
	})
}

func TestStrMapGet(t *testing.T) {
	t.Run("string map", func(t *testing.T) {
		m := map[string]string{
			"test_key": "test_value",
		}
		mtable, mtype := MapUnpack(m)
		val := StrMapGetAs[string, string](mtable, mtype, "test_key")
		if val != "test_value" {
			t.Errorf("Expected 'test_value', got %q", val)
		}
	})

	t.Run("int map", func(t *testing.T) {
		m := map[string]int{
			"100": 1,
			"200": 2,
		}
		mtable, mtype := MapUnpack(m)
		val := StrMapGetAs[string, int](mtable, mtype, "100")
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}
	})

	t.Run("struct map", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}

		m := map[string]TestStruct{
			"test_key": {Name: "test_value", Age: 30},
		}
		mtable, mtype := MapUnpack(m)
		val := StrMapGetAs[string, TestStruct](mtable, mtype, "test_key")
		if val.Name != "test_value" || val.Age != 30 {
			t.Errorf("Expected {Name: test_value, Age: 30}, got %+v", val)
		}
	})
}

func TestMapClone(t *testing.T) {
	t.Run("string map", func(t *testing.T) {
		m := map[string]string{
			"test_key":  "test_value",
			"test_key2": "test_value2",
		}
		mtable, mtype := MapUnpack(m)
		clone := MapCloneAs[string, string](mtable, mtype)
		if clone == nil {
			t.Fatalf("clone is nil")
		}
		if len(clone) != len(m) {
			t.Errorf("Expected %d, got %d", len(m), len(clone))
		}
		for k, v := range clone {
			if v != m[k] {
				t.Errorf("Expected %q, got %q", m[k], v)
			}
		}
	})

	t.Run("int map", func(t *testing.T) {
		m := map[int]int{
			100: 1,
			200: 2,
		}
		mtable, mtype := MapUnpack(m)
		clone := MapCloneAs[int, int](mtable, mtype)
		if clone == nil {
			t.Fatalf("clone is nil")
		}
		if len(clone) != len(m) {
			t.Errorf("Expected %d, got %d", len(m), len(clone))
		}
		for k, v := range clone {
			if v != m[k] {
				t.Errorf("Expected %d, got %d", m[k], v)
			}
		}
	})

	t.Run("struct map", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}

		m := map[string]TestStruct{
			"test_key":  {Name: "test_value", Age: 30},
			"test_key2": {Name: "test_value2", Age: 40},
		}
		mtable, mtype := MapUnpack(m)
		clone := MapCloneAs[string, TestStruct](mtable, mtype)
		if clone == nil {
			t.Fatalf("clone is nil")
		}
		if len(clone) != len(m) {
			t.Errorf("Expected %d, got %d", len(m), len(clone))
		}
		for k, v := range clone {
			if v != m[k] {
				t.Errorf("Expected %+v, got %+v", m[k], v)
			}
		}
	})
}

func TestMapClear(t *testing.T) {
	t.Run("string map", func(t *testing.T) {
		m := map[string]string{
			"test_key": "test_value",
		}
		mtable, mtype := MapUnpack(m)
		MapClear(mtable, mtype)
		if len(m) != 0 {
			t.Errorf("Expected %d, got %d", 0, len(m))
		}
	})

	t.Run("int map", func(t *testing.T) {
		m := map[int]int{
			100: 1,
			200: 2,
		}
		mtable, mtype := MapUnpack(m)
		MapClear(mtable, mtype)
		if len(m) != 0 {
			t.Errorf("Expected %d, got %d", 0, len(m))
		}
	})

	t.Run("struct map", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}

		m := map[string]TestStruct{
			"test_key":  {Name: "test_value", Age: 30},
			"test_key2": {Name: "test_value2", Age: 40},
		}
		mtable, mtype := MapUnpack(m)
		MapClear(mtable, mtype)
		if len(m) != 0 {
			t.Errorf("Expected %d, got %d", 0, len(m))
		}
	})
}
