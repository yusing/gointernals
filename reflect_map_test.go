//go:build go1.24 && go1.25 && !go1.26

package gointernals

import (
	"reflect"
	"testing"
)

func TestInitMap(t *testing.T) {
	t.Run("nil map to non-nil", func(t *testing.T) {
		var m map[string]int
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 5)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}
		if len(m) != 0 {
			t.Errorf("Expected len 0 for new map, got %d", len(m))
		}
		// Test that we can add elements
		m["test"] = 42
		if m["test"] != 42 {
			t.Error("Map should be writable after initialization")
		}
	})

	t.Run("string to int map", func(t *testing.T) {
		var m map[string]int
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 10)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}
		if len(m) != 0 {
			t.Errorf("Expected len 0 for new map, got %d", len(m))
		}

		// Test basic operations
		m["hello"] = 1
		m["world"] = 2
		if m["hello"] != 1 || m["world"] != 2 {
			t.Error("Map should support basic operations")
		}
	})

	t.Run("int to string map", func(t *testing.T) {
		var m map[int]string
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 3)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}
		if len(m) != 0 {
			t.Errorf("Expected len 0 for new map, got %d", len(m))
		}

		// Test basic operations
		m[1] = "one"
		m[2] = "two"
		if m[1] != "one" || m[2] != "two" {
			t.Error("Map should support basic operations")
		}
	})

	t.Run("struct key and value map", func(t *testing.T) {
		type Key struct {
			Name string
			ID   int
		}
		type Value struct {
			Data string
			Code int
		}
		var m map[Key]Value
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 2)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}
		if len(m) != 0 {
			t.Errorf("Expected len 0 for new map, got %d", len(m))
		}

		// Test basic operations
		key1 := Key{Name: "test", ID: 1}
		val1 := Value{Data: "hello", Code: 100}
		m[key1] = val1
		if m[key1] != val1 {
			t.Error("Map should support struct keys and values")
		}
	})

	t.Run("pointer key map", func(t *testing.T) {
		var m map[*int]string
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 1)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}
		if len(m) != 0 {
			t.Errorf("Expected len 0 for new map, got %d", len(m))
		}

		// Test basic operations
		key := new(int)
		*key = 42
		m[key] = "pointer"
		if m[key] != "pointer" {
			t.Error("Map should support pointer keys")
		}
	})

	t.Run("interface value map", func(t *testing.T) {
		var m map[string]any
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 5)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}
		if len(m) != 0 {
			t.Errorf("Expected len 0 for new map, got %d", len(m))
		}

		// Test basic operations
		m["int"] = 42
		m["string"] = "hello"
		m["float"] = 3.14
		if m["int"] != 42 || m["string"] != "hello" || m["float"] != 3.14 {
			t.Error("Map should support interface{} values")
		}
	})

	t.Run("byte slice key map", func(t *testing.T) {
		var m map[string][]byte
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 3)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}
		if len(m) != 0 {
			t.Errorf("Expected len 0 for new map, got %d", len(m))
		}

		// Test basic operations
		m["data"] = []byte{1, 2, 3}
		if len(m["data"]) != 3 || m["data"][0] != 1 {
			t.Error("Map should support slice values")
		}
	})

	t.Run("zero length capacity", func(t *testing.T) {
		var m map[int]int
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 0)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}
		if len(m) != 0 {
			t.Errorf("Expected len 0 for new map, got %d", len(m))
		}

		// Test that we can still add elements
		m[1] = 10
		if m[1] != 10 {
			t.Error("Map with zero capacity should still be writable")
		}
	})

	t.Run("replace existing map", func(t *testing.T) {
		m := map[string]int{"existing": 1, "data": 2}
		originalLen := len(m)
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 10)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}
		if len(m) != 0 {
			t.Errorf("Expected len 0 for new map, got %d", len(m))
		}
		// Original data should be gone
		if m["existing"] != 0 || m["data"] != 0 {
			t.Error("Original map data should be cleared")
		}
		if originalLen == len(m) {
			t.Log("Original map was replaced as expected")
		}
	})

	t.Run("complex key type", func(t *testing.T) {
		type ComplexKey struct {
			Str    string
			Num    int
			Slice  [3]int // Use array instead of slice for comparable type
			Nested struct {
				Field string
			}
		}
		var m map[ComplexKey]string
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 1)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}

		// Test basic operations
		key := ComplexKey{
			Str:   "test",
			Num:   42,
			Slice: [3]int{1, 2, 3},
			Nested: struct {
				Field string
			}{Field: "nested"},
		}
		m[key] = "value"
		if m[key] != "value" {
			t.Error("Map should support complex key types")
		}
	})

	t.Run("bool key map", func(t *testing.T) {
		var m map[bool]string
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 2)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}

		// Test basic operations
		m[true] = "yes"
		m[false] = "no"
		if m[true] != "yes" || m[false] != "no" {
			t.Error("Map should support bool keys")
		}
	})

	t.Run("rune key map", func(t *testing.T) {
		var m map[rune]int
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 3)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}

		// Test basic operations
		m['A'] = 65
		m['B'] = 66
		if m['A'] != 65 || m['B'] != 66 {
			t.Error("Map should support rune keys")
		}
	})

	t.Run("float key map", func(t *testing.T) {
		var m map[float64]string
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 2)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}

		// Test basic operations
		m[3.14] = "pi"
		m[2.71] = "e"
		if m[3.14] != "pi" || m[2.71] != "e" {
			t.Error("Map should support float keys")
		}
	})

	t.Run("large capacity", func(t *testing.T) {
		var m map[int]int
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 1000)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}
		if len(m) != 0 {
			t.Errorf("Expected len 0 for new map, got %d", len(m))
		}

		// Test that we can add many elements efficiently
		for i := range 100 {
			m[i] = i * 2
		}
		if len(m) != 100 {
			t.Errorf("Expected len 100 after adding elements, got %d", len(m))
		}
	})

	t.Run("panic on non-map type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for non-map type")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectInitMap of non-map type" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		i := 42
		v := reflect.ValueOf(&i).Elem()
		ReflectInitMap(v, 1)
	})

	t.Run("panic on slice type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for slice type")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectInitMap of non-map type" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		s := []int{1, 2, 3}
		v := reflect.ValueOf(&s).Elem()
		ReflectInitMap(v, 1)
	})

	t.Run("panic on string type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for string type")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectInitMap of non-map type" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		s := "hello"
		v := reflect.ValueOf(&s).Elem()
		ReflectInitMap(v, 1)
	})

	t.Run("panic on array type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for array type")
			} else if msg, ok := r.(string); !ok || msg != "gointernals.ReflectInitMap of non-map type" {
				t.Errorf("Expected specific panic message, got %v", r)
			}
		}()

		arr := [3]int{1, 2, 3}
		v := reflect.ValueOf(&arr).Elem()
		ReflectInitMap(v, 1)
	})

	t.Run("map with map values", func(t *testing.T) {
		var m map[string]map[int]string
		v := reflect.ValueOf(&m).Elem()

		ReflectInitMap(v, 2)

		if m == nil {
			t.Error("Expected map to be non-nil after initialization")
		}

		// Test nested map operations
		inner := make(map[int]string)
		inner[1] = "one"
		m["numbers"] = inner

		if m["numbers"][1] != "one" {
			t.Error("Map should support nested map values")
		}
	})
}

func TestReflectStrMapAssign_BasicAssignAndOverwrite(t *testing.T) {
	var m map[string]int
	mv := reflect.ValueOf(&m).Elem()

	ReflectInitMap(mv, 0)

	// assign new key
	ev := ReflectStrMapAssign(mv, "a")
	if ev.Kind() != reflect.Int {
		t.Fatalf("unexpected element value kind: %v -> %v", ev.Kind(), ev.Elem().Kind())
	}
	ev.SetInt(1)
	if got := m["a"]; got != 1 {
		t.Fatalf("want 1, got %d", got)
	}

	// overwrite existing key
	ev2 := ReflectStrMapAssign(mv, "a")
	ev2.SetInt(2)
	if got := m["a"]; got != 2 {
		t.Fatalf("want 2, got %d", got)
	}
}

func TestReflectStrMapAssign_StructValues(t *testing.T) {
	type S struct{ X, Y int }
	var m map[string]S
	mv := reflect.ValueOf(&m).Elem()

	ReflectInitMap(mv, 0)

	ev := ReflectStrMapAssign(mv, "p")
	if ev.Kind() != reflect.Struct {
		t.Fatalf("unexpected kind: %v -> %v", ev.Kind(), ev.Elem().Kind())
	}
	ev.Field(0).SetInt(10)
	ev.Field(1).SetInt(20)
	if got := m["p"]; got != (S{10, 20}) {
		t.Fatalf("want {10 20}, got %+v", got)
	}
}

func TestReflectStrMapAssign_MultipleKeys(t *testing.T) {
	var m map[string]string
	mv := reflect.ValueOf(&m).Elem()

	ReflectInitMap(mv, 0)

	keys := []string{"k1", "k2", "k3"}
	vals := []string{"v1", "v2", "v3"}
	for i := range keys {
		ev := ReflectStrMapAssign(mv, keys[i])
		ev.SetString(vals[i])
	}
	if len(m) != 3 || m["k1"] != "v1" || m["k2"] != "v2" || m["k3"] != "v3" {
		t.Fatalf("unexpected map contents: %#v", m)
	}
}

func TestReflectStrMapAssign_MapStringString(t *testing.T) {
	var m map[string]string
	{
		mv := reflect.ValueOf(&m).Elem()

		ReflectInitMap(mv, 0)

		ev := ReflectStrMapAssign(mv, "p")
		ev.SetString("10")
	}
	if got := m["p"]; got != "10" {
		t.Fatalf("want 10, got %s", got)
	}

	for _, v := range m {
		if v != "10" {
			t.Fatalf("want 10, got %s", v)
		}
	}
}

func TestReflectStrMapAssign_MapStringStructPtr(t *testing.T) {
	type S struct{ X int }
	var m map[string]*S
	mv := reflect.ValueOf(&m).Elem()

	ReflectInitMap(mv, 0)

	ev := ReflectStrMapAssign(mv, "p")
	ev.Set(reflect.New(ev.Type().Elem()))
	ev.Elem().Field(0).SetInt(10)
	if got := *m["p"]; got != (S{10}) {
		t.Fatalf("want {10}, got %+v", got)
	}

	for _, v := range m {
		if v == nil {
			t.Fatalf("want non-nil pointer")
		}
		if v.X != 10 {
			t.Fatalf("want 10, got %d", v.X)
		}
	}
}

func TestReflectStrMapAssign_PanicOnNilMap(t *testing.T) {
	var m map[string]int // nil map
	mv := reflect.ValueOf(&m).Elem()

	// Do not initialize to verify behavior on nil map
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil map assign")
		}
	}()
	_ = ReflectStrMapAssign(mv, "x")
}

func TestReflectMapAssign_BasicIntMap(t *testing.T) {
	var m map[int]int
	mv := reflect.ValueOf(&m).Elem()

	ReflectInitMap(mv, 0)

	ev := ReflectMapAssign(mv, 1)
	if ev.Kind() != reflect.Int {
		t.Fatalf("unexpected element value kind: %v -> %v", ev.Kind(), ev.Elem().Kind())
	}
	ev.SetInt(42)
	if got := m[1]; got != 42 {
		t.Fatalf("want 42, got %d", got)
	}

	// overwrite
	ev2 := ReflectMapAssign(mv, 1)
	ev2.SetInt(7)
	if got := m[1]; got != 7 {
		t.Fatalf("want 7, got %d", got)
	}
}

func TestReflectMapAssign_StringKeyViaAny(t *testing.T) {
	var m map[string]int
	mv := reflect.ValueOf(&m).Elem()

	ReflectInitMap(mv, 0)

	ev := ReflectMapAssign(mv, "k")
	ev.SetInt(3)
	if got := m["k"]; got != 3 {
		t.Fatalf("want 3, got %d", got)
	}
}

func TestReflectMapAssign_StructKeyValue(t *testing.T) {
	type K struct{ A int }
	type V struct{ S string }
	var m map[K]V
	mv := reflect.ValueOf(&m).Elem()

	ReflectInitMap(mv, 0)

	key := K{A: 9}
	ev := ReflectMapAssign(mv, key)
	if ev.Kind() != reflect.Struct {
		t.Fatalf("unexpected kind: %v -> %v", ev.Kind(), ev.Elem().Kind())
	}
	ev.Field(0).SetString("x")
	if got := m[key]; got != (V{S: "x"}) {
		t.Fatalf("unexpected value: %+v", got)
	}
}

func TestReflectMapAssign_InterfaceKey(t *testing.T) {
	var m map[any]string
	mv := reflect.ValueOf(&m).Elem()

	ReflectInitMap(mv, 0)

	key1 := 10
	key2 := "s"
	ReflectMapAssign(mv, key1).SetString("ten")
	ReflectMapAssign(mv, key2).SetString("str")

	if m[10] != "ten" || m["s"] != "str" {
		t.Fatalf("unexpected map contents: %#v", m)
	}
}

func TestReflectMapAssign_PointerKey(t *testing.T) {
	var m map[*int]float64
	mv := reflect.ValueOf(&m).Elem()

	ReflectInitMap(mv, 0)

	k := new(int)
	*k = 5
	ReflectMapAssign(mv, k).SetFloat(1.5)
	if m[k] != 1.5 {
		t.Fatalf("want 1.5, got %v", m[k])
	}
}

func TestReflectMapAssign_PanicOnNilMap(t *testing.T) {
	var m map[int]int
	mv := reflect.ValueOf(&m).Elem()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil map assign")
		}
	}()
	_ = ReflectMapAssign(mv, 1)
}

func TestReflectMapAssign_PanicOnNonMap(t *testing.T) {
	i := 0
	iv := reflect.ValueOf(&i).Elem()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on non-map value")
		}
	}()
	_ = ReflectMapAssign(iv, 1)
}
