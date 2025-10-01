//go:build go1.24 && go1.25 && !go1.26

package gointernals

import (
	"math"
	"reflect"
	"testing"
)

func TestReflectStrToNumBool_Ints(t *testing.T) {
	t.Run("int64 from decimal", func(t *testing.T) {
		var v int64
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "9223372036854775807"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != math.MaxInt64 {
			t.Fatalf("want %d, got %d", int64(math.MaxInt64), v)
		}
	})

	t.Run("int32 within range", func(t *testing.T) {
		var v int32
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "2147483647"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != math.MaxInt32 {
			t.Fatalf("want %d, got %d", int32(math.MaxInt32), v)
		}
	})

	t.Run("int32 out of range error", func(t *testing.T) {
		var v int32
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "2147483648"); err == nil {
			t.Fatal("expected range error, got nil")
		}
	})

	t.Run("negative to int16", func(t *testing.T) {
		var v int16
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "-32768"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != math.MinInt16 {
			t.Fatalf("want %d, got %d", int16(math.MinInt16), v)
		}
	})
}

func TestReflectStrToNumBool_Uints(t *testing.T) {
	t.Run("uint64 from decimal", func(t *testing.T) {
		var v uint64
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "18446744073709551615"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != math.MaxUint64 {
			t.Fatalf("want %d, got %d", uint64(math.MaxUint64), v)
		}
	})

	t.Run("uint8 within range", func(t *testing.T) {
		var v uint8
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "255"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != math.MaxUint8 {
			t.Fatalf("want %d, got %d", uint8(math.MaxUint8), v)
		}
	})

	t.Run("uint8 out of range error", func(t *testing.T) {
		var v uint8
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "256"); err == nil {
			t.Fatal("expected range error, got nil")
		}
	})

	t.Run("uint parse negative error", func(t *testing.T) {
		var v uint64
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "-1"); err == nil {
			t.Fatal("expected parse error for negative to uint, got nil")
		}
	})
}

func TestReflectStrToNumBool_Floats(t *testing.T) {
	t.Run("float32 parse", func(t *testing.T) {
		var v float32
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "3.5"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v < 3.499 || v > 3.501 {
			t.Fatalf("want approx 3.5, got %f", v)
		}
	})

	t.Run("float64 parse", func(t *testing.T) {
		var v float64
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "2.718281828"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v < 2.718281827 || v > 2.718281829 {
			t.Fatalf("want approx 2.718281828, got %f", v)
		}
	})

	t.Run("float parse error", func(t *testing.T) {
		var v float64
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "not-a-number"); err == nil {
			t.Fatal("expected parse error, got nil")
		}
	})
}

func TestReflectStrToNumBool_Bool(t *testing.T) {
	t.Run("true parse", func(t *testing.T) {
		var v bool
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != true {
			t.Fatalf("want true, got %v", v)
		}
	})

	t.Run("false parse", func(t *testing.T) {
		var v bool
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "false"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != false {
			t.Fatalf("want false, got %v", v)
		}
	})

	t.Run("invalid bool parse error", func(t *testing.T) {
		var v bool
		rv := reflect.ValueOf(&v).Elem()
		if err := ReflectStrToNumBool(rv, "not-bool"); err == nil {
			t.Fatal("expected parse error, got nil")
		}
	})
}

func TestReflectStrToNumBool_PanicsOnInvalidDstKind(t *testing.T) {
	t.Run("string dst panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for string destination")
			}
		}()
		var v string
		rv := reflect.ValueOf(&v).Elem()
		_ = ReflectStrToNumBool(rv, "123")
	})
}
