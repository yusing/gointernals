//go:build go1.24 && go1.25 && !go1.26

package gointernals

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/yusing/gointernals/abi"
)

func TestReflectValueType(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := reflect.ValueOf(42)
		typ := ReflectValueType(v)
		if typ == nil {
			t.Fatal("Expected non-nil type")
		}
		if typ.Kind() != abi.Int {
			t.Errorf("Expected kind Int, got %v", typ.Kind())
		}
		if typ.Size != unsafe.Sizeof(int(0)) {
			t.Errorf("Expected size %d, got %d", unsafe.Sizeof(int(0)), typ.Size)
		}
	})

	t.Run("string", func(t *testing.T) {
		v := reflect.ValueOf("hello")
		typ := ReflectValueType(v)
		if typ == nil {
			t.Fatal("Expected non-nil type")
		}
		if typ.Kind() != abi.String {
			t.Errorf("Expected kind String, got %v", typ.Kind())
		}
		if typ.Size != unsafe.Sizeof("") {
			t.Errorf("Expected size %d, got %d", unsafe.Sizeof(""), typ.Size)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type testStruct struct {
			A int
			B string
		}
		v := reflect.ValueOf(testStruct{A: 1, B: "test"})
		typ := ReflectValueType(v)
		if typ == nil {
			t.Fatal("Expected non-nil type")
		}
		if typ.Kind() != abi.Struct {
			t.Errorf("Expected kind Struct, got %v", typ.Kind())
		}
		if typ.Size != unsafe.Sizeof(testStruct{}) {
			t.Errorf("Expected size %d, got %d", unsafe.Sizeof(testStruct{}), typ.Size)
		}
	})

	t.Run("slice", func(t *testing.T) {
		v := reflect.ValueOf([]int{1, 2, 3})
		typ := ReflectValueType(v)
		if typ == nil {
			t.Fatal("Expected non-nil type")
		}
		if typ.Kind() != abi.Slice {
			t.Errorf("Expected kind Slice, got %v", typ.Kind())
		}
	})

	t.Run("pointer", func(t *testing.T) {
		x := 42
		v := reflect.ValueOf(&x)
		typ := ReflectValueType(v)
		if typ == nil {
			t.Fatal("Expected non-nil type")
		}
		if typ.Kind() != abi.Pointer {
			t.Errorf("Expected kind Pointer, got %v", typ.Kind())
		}
		if typ.Size != unsafe.Sizeof(uintptr(0)) {
			t.Errorf("Expected size %d, got %d", unsafe.Sizeof(uintptr(0)), typ.Size)
		}
	})

	t.Run("bool", func(t *testing.T) {
		v := reflect.ValueOf(true)
		typ := ReflectValueType(v)
		if typ == nil {
			t.Fatal("Expected non-nil type")
		}
		if typ.Kind() != abi.Bool {
			t.Errorf("Expected kind Bool, got %v", typ.Kind())
		}
		if typ.Size != unsafe.Sizeof(bool(false)) {
			t.Errorf("Expected size %d, got %d", unsafe.Sizeof(bool(false)), typ.Size)
		}
	})

	t.Run("float64", func(t *testing.T) {
		v := reflect.ValueOf(3.14)
		typ := ReflectValueType(v)
		if typ == nil {
			t.Fatal("Expected non-nil type")
		}
		if typ.Kind() != abi.Float64 {
			t.Errorf("Expected kind Float64, got %v", typ.Kind())
		}
		if typ.Size != unsafe.Sizeof(float64(0)) {
			t.Errorf("Expected size %d, got %d", unsafe.Sizeof(float64(0)), typ.Size)
		}
	})
}

func TestReflectValueData(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		x := 42
		v := reflect.ValueOf(&x).Elem()
		data := ReflectValueData(v)
		if data == nil {
			t.Fatal("Expected non-nil data pointer")
		}
		if *(*int)(data) != 42 {
			t.Errorf("Expected value 42, got %d", *(*int)(data))
		}
		*(*int)(data) = 100
		if x != 100 {
			t.Errorf("Expected modified value 100, got %d", x)
		}
	})

	t.Run("string", func(t *testing.T) {
		s := "hello"
		v := reflect.ValueOf(&s).Elem()
		data := ReflectValueData(v)
		if data == nil {
			t.Fatal("Expected non-nil data pointer")
		}
		if *(*string)(data) != "hello" {
			t.Errorf("Expected value 'hello', got %s", *(*string)(data))
		}
	})

	t.Run("bool", func(t *testing.T) {
		b := true
		v := reflect.ValueOf(&b).Elem()
		data := ReflectValueData(v)
		if data == nil {
			t.Fatal("Expected non-nil data pointer")
		}
		if *(*bool)(data) != true {
			t.Errorf("Expected value true, got %v", *(*bool)(data))
		}
		*(*bool)(data) = false
		if b != false {
			t.Errorf("Expected modified value false, got %v", b)
		}
	})

	t.Run("float64", func(t *testing.T) {
		f := 3.14
		v := reflect.ValueOf(&f).Elem()
		data := ReflectValueData(v)
		if data == nil {
			t.Fatal("Expected non-nil data pointer")
		}
		if *(*float64)(data) != 3.14 {
			t.Errorf("Expected value 3.14, got %f", *(*float64)(data))
		}
		*(*float64)(data) = 2.71
		if f != 2.71 {
			t.Errorf("Expected modified value 2.71, got %f", f)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type testStruct struct {
			A int
			B string
		}
		ts := testStruct{A: 10, B: "test"}
		v := reflect.ValueOf(&ts).Elem()
		data := ReflectValueData(v)
		if data == nil {
			t.Fatal("Expected non-nil data pointer")
		}
		result := *(*testStruct)(data)
		if result.A != 10 || result.B != "test" {
			t.Errorf("Expected {10, test}, got {%d, %s}", result.A, result.B)
		}
	})

	t.Run("slice", func(t *testing.T) {
		slice := []int{1, 2, 3}
		v := reflect.ValueOf(&slice).Elem()
		data := ReflectValueData(v)
		if data == nil {
			t.Fatal("Expected non-nil data pointer")
		}
		result := *(*[]int)(data)
		if len(result) != 3 || result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("Expected [1 2 3], got %v", result)
		}
	})
}

func TestReflectShallowCopy(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		src := 42
		dst := 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 42 {
			t.Errorf("Expected dst to be 42, got %d", dst)
		}
	})

	t.Run("string", func(t *testing.T) {
		src := "hello world"
		dst := ""
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != "hello world" {
			t.Errorf("Expected dst to be 'hello world', got %s", dst)
		}
	})

	t.Run("bool", func(t *testing.T) {
		src := true
		dst := false
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != true {
			t.Errorf("Expected dst to be true, got %v", dst)
		}
	})

	t.Run("float64", func(t *testing.T) {
		src := 3.14159
		dst := 0.0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 3.14159 {
			t.Errorf("Expected dst to be 3.14159, got %f", dst)
		}
	})

	t.Run("struct same size", func(t *testing.T) {
		type testStruct struct {
			A int
			B string
		}
		src := testStruct{A: 42, B: "test"}
		dst := testStruct{}
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst.A != 42 || dst.B != "test" {
			t.Errorf("Expected dst to be {42, test}, got {%d, %s}", dst.A, dst.B)
		}
	})

	t.Run("slice shallow copy", func(t *testing.T) {
		src := []int{1, 2, 3}
		dst := []int{}
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if len(dst) != 3 || dst[0] != 1 || dst[1] != 2 || dst[2] != 3 {
			t.Fatalf("Expected dst to be [1 2 3], got %v", dst)
		}
		src[0] = 100
		if dst[0] != 100 {
			t.Errorf("Expected shallow copy: dst[0] should be 100 after src[0] changed, got %d", dst[0])
		}
	})

	t.Run("pointer", func(t *testing.T) {
		x := 42
		src := &x
		var dst *int
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst == nil {
			t.Fatal("Expected dst to be non-nil")
		}
		if *dst != 42 {
			t.Errorf("Expected *dst to be 42, got %d", *dst)
		}
		if dst != src {
			t.Errorf("Expected dst and src to point to the same address")
		}
	})

	t.Run("int8 to int16 - smaller to larger", func(t *testing.T) {
		var src int8 = 42
		var dst int16 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 42 {
			t.Errorf("Expected dst to be 42, got %d", dst)
		}
	})

	t.Run("int16 to int32 - smaller to larger", func(t *testing.T) {
		var src int16 = 1000
		var dst int32 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 1000 {
			t.Errorf("Expected dst to be 1000, got %d", dst)
		}
	})

	t.Run("uint8 to uint64 - smaller to larger", func(t *testing.T) {
		var src uint8 = 255
		var dst uint64 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 255 {
			t.Errorf("Expected dst to be 255, got %d", dst)
		}
	})

	t.Run("float32 to float64 - smaller to larger", func(t *testing.T) {
		var src float32 = 3.14
		var dst float64 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst < 3.1399 || dst > 3.1401 {
			t.Errorf("Expected dst to be approximately 3.14, got %f", dst)
		}
	})

	// t.Run("int32 to int16 - larger to smaller panics", func(t *testing.T) {
	// 	defer func() {
	// 		if r := recover(); r == nil {
	// 			t.Errorf("Expected panic when copying from larger to smaller numeric type")
	// 		}
	// 	}()
	// 	var src int32 = 1000
	// 	var dst int16 = 0
	// 	srcV := reflect.ValueOf(&src).Elem()
	// 	dstV := reflect.ValueOf(&dst).Elem()
	// 	ReflectShallowCopy(dstV, srcV)
	// })

	// t.Run("uint64 to uint8 - larger to smaller panics", func(t *testing.T) {
	// 	defer func() {
	// 		if r := recover(); r == nil {
	// 			t.Errorf("Expected panic when copying from larger to smaller numeric type")
	// 		}
	// 	}()
	// 	var src uint64 = 255
	// 	var dst uint8 = 0
	// 	srcV := reflect.ValueOf(&src).Elem()
	// 	dstV := reflect.ValueOf(&dst).Elem()
	// 	ReflectShallowCopy(dstV, srcV)
	// })

	t.Run("array shallow copy", func(t *testing.T) {
		src := [3]int{10, 20, 30}
		dst := [3]int{}
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst[0] != 10 || dst[1] != 20 || dst[2] != 30 {
			t.Errorf("Expected dst to be [10 20 30], got %v", dst)
		}
		src[0] = 100
		if dst[0] == 100 {
			t.Errorf("Expected array copy to be independent: dst[0] should remain 10, got %d", dst[0])
		}
	})

	t.Run("int to uint same size positive value", func(t *testing.T) {
		var src int32 = 100
		var dst uint32 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 100 {
			t.Errorf("Expected dst to be 100, got %d", dst)
		}
	})

	// t.Run("int to uint same size negative value panics", func(t *testing.T) {
	// 	defer func() {
	// 		if r := recover(); r == nil {
	// 			t.Errorf("Expected panic when copying negative int to uint")
	// 		}
	// 	}()
	// 	var src int32 = -100
	// 	var dst uint32 = 0
	// 	srcV := reflect.ValueOf(&src).Elem()
	// 	dstV := reflect.ValueOf(&dst).Elem()
	// 	ReflectShallowCopy(dstV, srcV)
	// })

	t.Run("uint to int same size", func(t *testing.T) {
		var src uint32 = 100
		var dst int32 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 100 {
			t.Errorf("Expected dst to be 100, got %d", dst)
		}
	})

	t.Run("uint to int same size large value", func(t *testing.T) {
		var src uint32 = 4294967295
		var dst int32 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != -1 {
			t.Errorf("Expected dst to be -1 (overflow), got %d", dst)
		}
	})

	t.Run("int8 to uint16 positive value", func(t *testing.T) {
		var src int8 = 100
		var dst uint16 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 100 {
			t.Errorf("Expected dst to be 100, got %d", dst)
		}
	})

	// t.Run("int8 to uint16 negative value panics", func(t *testing.T) {
	// 	defer func() {
	// 		if r := recover(); r == nil {
	// 			t.Errorf("Expected panic when copying negative int8 to uint16")
	// 		}
	// 	}()
	// 	var src int8 = -50
	// 	var dst uint16 = 0
	// 	srcV := reflect.ValueOf(&src).Elem()
	// 	dstV := reflect.ValueOf(&dst).Elem()
	// 	ReflectShallowCopy(dstV, srcV)
	// })

	t.Run("uint8 to int64", func(t *testing.T) {
		var src uint8 = 255
		var dst int64 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 255 {
			t.Errorf("Expected dst to be 255, got %d", dst)
		}
	})

	// t.Run("int to float panics", func(t *testing.T) {
	// 	defer func() {
	// 		if r := recover(); r == nil {
	// 			t.Errorf("Expected panic when copying int to float")
	// 		}
	// 	}()
	// 	var src int32 = 100
	// 	var dst float32 = 0
	// 	srcV := reflect.ValueOf(&src).Elem()
	// 	dstV := reflect.ValueOf(&dst).Elem()
	// 	ReflectShallowCopy(dstV, srcV)
	// })

	t.Run("string to int panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when copying string to int")
			}
		}()
		var src string = "hello"
		var dst int = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
	})

	t.Run("same size same type int64", func(t *testing.T) {
		var src int64 = 9876543210
		var dst int64 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 9876543210 {
			t.Errorf("Expected dst to be 9876543210, got %d", dst)
		}
	})

	t.Run("same size same type uint64", func(t *testing.T) {
		var src uint64 = 18446744073709551615
		var dst uint64 = 0
		srcV := reflect.ValueOf(&src).Elem()
		dstV := reflect.ValueOf(&dst).Elem()
		ReflectShallowCopy(dstV, srcV)
		if dst != 18446744073709551615 {
			t.Errorf("Expected dst to be 18446744073709551615, got %d", dst)
		}
	})
}

func TestReflectSetZero(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		x := 42
		v := reflect.ValueOf(&x).Elem()
		ReflectSetZero(v)
		if x != 0 {
			t.Errorf("Expected x to be 0 after ReflectSetZero, got %d", x)
		}
	})

	t.Run("string", func(t *testing.T) {
		s := "hello world"
		v := reflect.ValueOf(&s).Elem()
		ReflectSetZero(v)
		if s != "" {
			t.Errorf("Expected s to be empty string after ReflectSetZero, got %s", s)
		}
	})

	t.Run("bool", func(t *testing.T) {
		b := true
		v := reflect.ValueOf(&b).Elem()
		ReflectSetZero(v)
		if b != false {
			t.Errorf("Expected b to be false after ReflectSetZero, got %v", b)
		}
	})

	t.Run("float64", func(t *testing.T) {
		f := 3.14159
		v := reflect.ValueOf(&f).Elem()
		ReflectSetZero(v)
		if f != 0.0 {
			t.Errorf("Expected f to be 0.0 after ReflectSetZero, got %f", f)
		}
	})

	t.Run("int64", func(t *testing.T) {
		var x int64 = 9876543210
		v := reflect.ValueOf(&x).Elem()
		ReflectSetZero(v)
		if x != 0 {
			t.Errorf("Expected x to be 0 after ReflectSetZero, got %d", x)
		}
	})

	t.Run("uint64", func(t *testing.T) {
		var x uint64 = 18446744073709551615
		v := reflect.ValueOf(&x).Elem()
		ReflectSetZero(v)
		if x != 0 {
			t.Errorf("Expected x to be 0 after ReflectSetZero, got %d", x)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type testStruct struct {
			A int
			B string
			C float64
		}
		ts := testStruct{A: 42, B: "test", C: 3.14}
		v := reflect.ValueOf(&ts).Elem()
		ReflectSetZero(v)
		if ts.A != 0 || ts.B != "" || ts.C != 0.0 {
			t.Errorf("Expected ts to be zero value {0, '', 0.0}, got {%d, %s, %f}", ts.A, ts.B, ts.C)
		}
	})

	t.Run("slice", func(t *testing.T) {
		slice := []int{1, 2, 3}
		v := reflect.ValueOf(&slice).Elem()
		ReflectSetZero(v)
		if slice != nil {
			t.Errorf("Expected slice to be nil after ReflectSetZero, got %v", slice)
		}
	})

	t.Run("map", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		v := reflect.ValueOf(&m).Elem()
		ReflectSetZero(v)
		if m != nil {
			t.Errorf("Expected m to be nil after ReflectSetZero, got %v", m)
		}
	})

	t.Run("pointer", func(t *testing.T) {
		x := 42
		p := &x
		v := reflect.ValueOf(&p).Elem()
		ReflectSetZero(v)
		if p != nil {
			t.Errorf("Expected p to be nil after ReflectSetZero, got %v", p)
		}
	})

	t.Run("interface", func(t *testing.T) {
		var i any = "hello"
		v := reflect.ValueOf(&i).Elem()
		ReflectSetZero(v)
		if i != nil {
			t.Errorf("Expected i to be nil after ReflectSetZero, got %v", i)
		}
	})

	t.Run("array", func(t *testing.T) {
		arr := [3]int{1, 2, 3}
		v := reflect.ValueOf(&arr).Elem()
		ReflectSetZero(v)
		if arr[0] != 0 || arr[1] != 0 || arr[2] != 0 {
			t.Errorf("Expected arr to be [0 0 0] after ReflectSetZero, got %v", arr)
		}
	})

	t.Run("complex64", func(t *testing.T) {
		c := complex64(3 + 4i)
		v := reflect.ValueOf(&c).Elem()
		ReflectSetZero(v)
		if c != 0 {
			t.Errorf("Expected c to be 0+0i after ReflectSetZero, got %v", c)
		}
	})

	t.Run("complex128", func(t *testing.T) {
		c := complex128(3 + 4i)
		v := reflect.ValueOf(&c).Elem()
		ReflectSetZero(v)
		if c != 0 {
			t.Errorf("Expected c to be 0+0i after ReflectSetZero, got %v", c)
		}
	})

	t.Run("nested struct", func(t *testing.T) {
		type Inner struct {
			X int
			Y string
		}
		type Outer struct {
			A Inner
			B int
		}
		outer := Outer{A: Inner{X: 10, Y: "test"}, B: 20}
		v := reflect.ValueOf(&outer).Elem()
		ReflectSetZero(v)
		if outer.A.X != 0 || outer.A.Y != "" || outer.B != 0 {
			t.Errorf("Expected outer to be zero value, got %+v", outer)
		}
	})
}

func TestReflectInitPtr(t *testing.T) {
	t.Run("pointer to int", func(t *testing.T) {
		var p *int
		v := reflect.ValueOf(&p).Elem()
		ReflectInitPtr(v)
		if p == nil {
			t.Fatal("Expected p to be non-nil after ReflectInitPtr")
		}
		if *p != 0 {
			t.Errorf("Expected *p to be 0 (zero value), got %d", *p)
		}
	})

	t.Run("pointer to string", func(t *testing.T) {
		var p *string
		v := reflect.ValueOf(&p).Elem()
		ReflectInitPtr(v)
		if p == nil {
			t.Fatal("Expected p to be non-nil after ReflectInitPtr")
		}
		if *p != "" {
			t.Errorf("Expected *p to be empty string (zero value), got %s", *p)
		}
	})

	t.Run("pointer to struct", func(t *testing.T) {
		type testStruct struct {
			A int
			B string
		}
		var p *testStruct
		v := reflect.ValueOf(&p).Elem()
		ReflectInitPtr(v)
		if p == nil {
			t.Fatal("Expected p to be non-nil after ReflectInitPtr")
		}
		if p.A != 0 || p.B != "" {
			t.Errorf("Expected *p to be zero value {0, ''}, got {%d, %s}", p.A, p.B)
		}
	})

	t.Run("pointer to slice", func(t *testing.T) {
		var p *[]int
		v := reflect.ValueOf(&p).Elem()
		ReflectInitPtr(v)
		if p == nil {
			t.Fatal("Expected p to be non-nil after ReflectInitPtr")
		}
		if *p != nil {
			t.Errorf("Expected *p to be nil slice (zero value), got %v", *p)
		}
	})

	t.Run("pointer to map", func(t *testing.T) {
		var p *map[string]int
		v := reflect.ValueOf(&p).Elem()
		ReflectInitPtr(v)
		if p == nil {
			t.Fatal("Expected p to be non-nil after ReflectInitPtr")
		}
		if *p != nil {
			t.Errorf("Expected *p to be nil map (zero value), got %v", *p)
		}
	})

	t.Run("pointer to bool", func(t *testing.T) {
		var p *bool
		v := reflect.ValueOf(&p).Elem()
		ReflectInitPtr(v)
		if p == nil {
			t.Fatal("Expected p to be non-nil after ReflectInitPtr")
		}
		if *p != false {
			t.Errorf("Expected *p to be false (zero value), got %v", *p)
		}
	})

	t.Run("pointer to float64", func(t *testing.T) {
		var p *float64
		v := reflect.ValueOf(&p).Elem()
		ReflectInitPtr(v)
		if p == nil {
			t.Fatal("Expected p to be non-nil after ReflectInitPtr")
		}
		if *p != 0.0 {
			t.Errorf("Expected *p to be 0.0 (zero value), got %f", *p)
		}
	})

	t.Run("pointer to pointer", func(t *testing.T) {
		var p **int
		v := reflect.ValueOf(&p).Elem()
		ReflectInitPtr(v)
		if p == nil {
			t.Fatal("Expected p to be non-nil after ReflectInitPtr")
		}
		if *p != nil {
			t.Errorf("Expected *p to be nil (zero value for *int), got %v", *p)
		}
	})

	t.Run("panic on string", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when calling ReflectInitPtr on string type")
			}
		}()
		var s string = "hello"
		v := reflect.ValueOf(&s).Elem()
		ReflectInitPtr(v)
	})

	t.Run("panic on slice", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when calling ReflectInitPtr on slice type")
			}
		}()
		var slice []int
		v := reflect.ValueOf(&slice).Elem()
		ReflectInitPtr(v)
	})

	t.Run("panic on struct", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when calling ReflectInitPtr on struct type")
			}
		}()
		type testStruct struct {
			A int
		}
		var ts testStruct
		v := reflect.ValueOf(&ts).Elem()
		ReflectInitPtr(v)
	})
}
