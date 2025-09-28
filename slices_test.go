//go:build go1.24 && go1.25 && !go1.26

package gointernals

import (
	"testing"
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

func TestSliceUnpack(t *testing.T) {
	t.Run("string slice", func(t *testing.T) {
		slice := []string{"hello"}
		sliceData, elemType := SliceUnpack(slice)
		if sliceData.ptr != unsafe.Pointer(&slice[0]) {
			t.Errorf("Expected sliceData.ptr to be %p, got %p", unsafe.Pointer(&slice[0]), sliceData.ptr)
		}
		if sliceData.len != 1 {
			t.Errorf("Expected sliceData.len to be 1, got %d", sliceData.len)
		}
		if sliceData.cap != 1 {
			t.Errorf("Expected sliceData.cap to be 1, got %d", sliceData.cap)
		}
		if elemType.Size != unsafe.Sizeof("") {
			t.Errorf("Expected elemType.Size_ to be %d, got %d", unsafe.Sizeof(string("")), elemType.Size)
		}
		if elemType.Kind != abi.String {
			t.Errorf("Expected elemType.Kind_ to be String, got %d", elemType.Kind)
		}
	})
}

func TestSliceClone(t *testing.T) {
	// Test with string slice
	t.Run("string slice", func(t *testing.T) {
		slice2 := SliceCloneAs[string](SliceUnpack([]string{"hello"}))

		if len(slice2) != 1 {
			t.Errorf("Expected length 2, got %d", len(slice2))
		} else {
			if slice2[0] != "hello" {
				t.Errorf("Expected [hello world], got %v", slice2)
			}
		}
	})

	// Test with int slice
	t.Run("int slice", func(t *testing.T) {
		slice2 := SliceCloneAs[int](SliceUnpack([]int{1, 2, 3}))

		if len(slice2) != 3 {
			t.Errorf("Expected length 3, got %d", len(slice2))
		} else {
			expected := []int{1, 2, 3}
			for i, exp := range expected {
				if i < len(slice2) && slice2[i] != exp {
					t.Errorf("At index %d: expected %d, got %d", i, exp, slice2[i])
				}
			}
		}
	})

	// Test with struct slice
	t.Run("struct slice", func(t *testing.T) {
		type TestStruct struct {
			Name string
			ID   int
		}

		slice2 := SliceCloneAs[TestStruct](SliceUnpack([]TestStruct{{Name: "Alice", ID: 1}, {Name: "Bob", ID: 2}}))

		if len(slice2) != 2 {
			t.Errorf("Expected length 2, got %d", len(slice2))
		} else {
			if slice2[0].Name != "Alice" || slice2[0].ID != 1 {
				t.Errorf("Original element changed: got %+v", slice2[0])
			}

			if slice2[1].Name != "Bob" || slice2[1].ID != 2 {
				t.Errorf("Expected {Name: Bob, ID: 2}, got %+v", slice2[1])
			}
		}
	})

	// Test with empty slice
	t.Run("empty slice", func(t *testing.T) {
		slice2 := SliceCloneAs[string](SliceUnpack([]string{}))

		if len(slice2) != 0 {
			t.Errorf("Expected length 0, got %d", len(slice2))
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		slice2 := SliceCloneAs[string](SliceUnpack([]string(nil)))

		if len(slice2) != 0 {
			t.Errorf("Expected length 0, got %d", len(slice2))
		}
	})
}

func TestSliceCloneInto(t *testing.T) {
	// Test case 1: Destination has sufficient capacity (need <= 0) - non-pointer type
	t.Run("sufficient capacity int", func(t *testing.T) {
		src := []int{1, 2, 3}
		srcSlice, elemType := SliceUnpack(src)

		dst := make([]int, 0, 5) // capacity 5, length 0
		dstSlice, _ := SliceUnpack(dst)

		SliceCloneInto(dstSlice, srcSlice, elemType)

		// Convert back to check results
		result := *(*[]int)(unsafe.Pointer(dstSlice))

		if len(result) != 3 {
			t.Errorf("Expected length 3, got %d", len(result))
		}
		if cap(result) != 5 {
			t.Errorf("Expected capacity 5, got %d", cap(result))
		}
		for i, expected := range []int{1, 2, 3} {
			if result[i] != expected {
				t.Errorf("At index %d: expected %d, got %d", i, expected, result[i])
			}
		}
	})

	// Test case 2: Destination has sufficient capacity (need <= 0) - pointer type
	t.Run("sufficient capacity string", func(t *testing.T) {
		src := []string{"hello", "world"}
		srcSlice, elemType := SliceUnpack(src)

		dst := make([]string, 0, 5) // capacity 5, length 0
		dstSlice, _ := SliceUnpack(dst)

		SliceCloneInto(dstSlice, srcSlice, elemType)

		// Convert back to check results
		result := *(*[]string)(unsafe.Pointer(dstSlice))

		if len(result) != 2 {
			t.Errorf("Expected length 2, got %d", len(result))
		}
		if cap(result) != 5 {
			t.Errorf("Expected capacity 5, got %d", cap(result))
		}
		for i, expected := range []string{"hello", "world"} {
			if result[i] != expected {
				t.Errorf("At index %d: expected %s, got %s", i, expected, result[i])
			}
		}
	})

	// Test case 3: Destination needs growth (need > 0 and dst.ptr != nil) - non-pointer type
	t.Run("needs growth int", func(t *testing.T) {
		src := []int{1, 2, 3, 4, 5}
		srcSlice, elemType := SliceUnpack(src)

		dst := make([]int, 0, 2) // capacity 2, needs growth
		dstSlice, _ := SliceUnpack(dst)

		SliceCloneInto(dstSlice, srcSlice, elemType)

		// Convert back to check results
		result := *(*[]int)(unsafe.Pointer(dstSlice))

		if len(result) != 5 {
			t.Errorf("Expected length 5, got %d", len(result))
		}
		if cap(result) < 5 {
			t.Errorf("Expected capacity at least 5, got %d", cap(result))
		}
		for i, expected := range []int{1, 2, 3, 4, 5} {
			if result[i] != expected {
				t.Errorf("At index %d: expected %d, got %d", i, expected, result[i])
			}
		}
	})

	// Test case 4: Destination needs growth (need > 0 and dst.ptr != nil) - pointer type
	t.Run("needs growth string", func(t *testing.T) {
		src := []string{"a", "b", "c", "d"}
		srcSlice, elemType := SliceUnpack(src)

		dst := make([]string, 0, 2) // capacity 2, needs growth
		dstSlice, _ := SliceUnpack(dst)

		SliceCloneInto(dstSlice, srcSlice, elemType)

		// Convert back to check results
		result := *(*[]string)(unsafe.Pointer(dstSlice))

		if len(result) != 4 {
			t.Errorf("Expected length 4, got %d", len(result))
		}
		if cap(result) < 4 {
			t.Errorf("Expected capacity at least 4, got %d", cap(result))
		}
		for i, expected := range []string{"a", "b", "c", "d"} {
			if result[i] != expected {
				t.Errorf("At index %d: expected %s, got %s", i, expected, result[i])
			}
		}
	})

	// Test case 5: Destination is nil/empty (dst.ptr == nil) - non-pointer type
	t.Run("nil destination int", func(t *testing.T) {
		src := []int{10, 20, 30}
		srcSlice, elemType := SliceUnpack(src)

		var dst []int
		dstSlice, _ := SliceUnpack(dst)

		SliceCloneInto(dstSlice, srcSlice, elemType)

		// Convert back to check results
		result := *(*[]int)(unsafe.Pointer(dstSlice))

		if len(result) != 3 {
			t.Errorf("Expected length 3, got %d", len(result))
		}
		if cap(result) != 3 {
			t.Errorf("Expected capacity 3, got %d", cap(result))
		}
		for i, expected := range []int{10, 20, 30} {
			if result[i] != expected {
				t.Errorf("At index %d: expected %d, got %d", i, expected, result[i])
			}
		}
	})

	// Test case 6: Destination is nil/empty (dst.ptr == nil) - pointer type
	t.Run("nil destination string", func(t *testing.T) {
		src := []string{"x", "y"}
		srcSlice, elemType := SliceUnpack(src)

		var dst []string
		dstSlice, _ := SliceUnpack(dst)

		SliceCloneInto(dstSlice, srcSlice, elemType)

		// Convert back to check results
		result := *(*[]string)(unsafe.Pointer(dstSlice))

		if len(result) != 2 {
			t.Errorf("Expected length 2, got %d", len(result))
		}
		if cap(result) != 2 {
			t.Errorf("Expected capacity 2, got %d", cap(result))
		}
		for i, expected := range []string{"x", "y"} {
			if result[i] != expected {
				t.Errorf("At index %d: expected %s, got %s", i, expected, result[i])
			}
		}
	})

	// Test case 7: Empty source slice
	t.Run("empty source", func(t *testing.T) {
		src := []int{}
		srcSlice, elemType := SliceUnpack(src)

		dst := make([]int, 0, 5)
		dstSlice, _ := SliceUnpack(dst)

		SliceCloneInto(dstSlice, srcSlice, elemType)

		// Convert back to check results
		result := *(*[]int)(unsafe.Pointer(dstSlice))

		if len(result) != 0 {
			t.Errorf("Expected length 0, got %d", len(result))
		}
		if cap(result) != 5 {
			t.Errorf("Expected capacity 5, got %d", cap(result))
		}
	})

	// Test case 8: Struct slice (contains pointers)
	t.Run("struct slice", func(t *testing.T) {
		type TestStruct struct {
			Name string
			ID   int
		}

		src := []TestStruct{{Name: "Alice", ID: 1}, {Name: "Bob", ID: 2}}
		srcSlice, elemType := SliceUnpack(src)

		dst := make([]TestStruct, 0, 1) // needs growth
		dstSlice, _ := SliceUnpack(dst)

		SliceCloneInto(dstSlice, srcSlice, elemType)

		// Convert back to check results
		result := *(*[]TestStruct)(unsafe.Pointer(dstSlice))

		if len(result) != 2 {
			t.Errorf("Expected length 2, got %d", len(result))
		}
		if cap(result) < 2 {
			t.Errorf("Expected capacity at least 2, got %d", cap(result))
		}

		if result[0].Name != "Alice" || result[0].ID != 1 {
			t.Errorf("Expected {Alice 1}, got %+v", result[0])
		}
		if result[1].Name != "Bob" || result[1].ID != 2 {
			t.Errorf("Expected {Bob 2}, got %+v", result[1])
		}
	})

	// Test case 9: Pointer slice
	t.Run("pointer slice", func(t *testing.T) {
		val1, val2 := 42, 84
		src := []*int{&val1, &val2}
		srcSlice, elemType := SliceUnpack(src)

		var dst []*int
		dstSlice, _ := SliceUnpack(dst)

		SliceCloneInto(dstSlice, srcSlice, elemType)

		// Convert back to check results
		result := *(*[]*int)(unsafe.Pointer(dstSlice))

		if len(result) != 2 {
			t.Errorf("Expected length 2, got %d", len(result))
		}
		if cap(result) != 2 {
			t.Errorf("Expected capacity 2, got %d", cap(result))
		}

		if result[0] == nil || *result[0] != 42 {
			t.Errorf("Expected pointer to 42, got %v", result[0])
		}
		if result[1] == nil || *result[1] != 84 {
			t.Errorf("Expected pointer to 84, got %v", result[1])
		}

		// Verify it's a shallow copy (same pointers)
		if result[0] != &val1 {
			t.Errorf("Expected same pointer reference for val1")
		}
		if result[1] != &val2 {
			t.Errorf("Expected same pointer reference for val2")
		}
	})
}
