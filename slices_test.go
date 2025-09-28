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
