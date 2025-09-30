package gointernals

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/yusing/gointernals/abi"
)

//go:nosplit
func ReflectStrToNumBool(dst reflect.Value, src string) error {
	dstTKind := dst.Kind()
	if !ReflectIsNumeric(dst) || dstTKind == reflect.Bool {
		panic(fmt.Errorf("gointernals.ReflectStrToNumBool: invalid destination type %s", dst.Type()))
	}
	switch {
	case ReflectCanInt(dst):
		i, err := strconv.ParseInt(src, 10, 64)
		if err != nil {
			return err
		}
		switch abi.Kind(dstTKind).Size() {
		case 8:
			ReflectValueSet(dst, int64(i))
		case 4:
			if i > math.MaxInt32 {
				return fmt.Errorf("gointernals.ReflectStrToNumBool: value out of range for int32: %s", src)
			}
			ReflectValueSet(dst, int32(i))
		case 2:
			if i > math.MaxInt16 {
				return fmt.Errorf("gointernals.ReflectStrToNumBool: value out of range for int16: %s", src)
			}
			ReflectValueSet(dst, int16(i))
		case 1:
			if i > math.MaxInt8 {
				return fmt.Errorf("gointernals.ReflectStrToNumBool: value out of range for int8: %s", src)
			}
			ReflectValueSet(dst, int8(i))
		}
	case ReflectCanUint(dst):
		i, err := strconv.ParseUint(src, 10, 64)
		if err != nil {
			return err
		}
		switch abi.Kind(dstTKind).Size() {
		case 8:
			ReflectValueSet(dst, uint64(i))
		case 4:
			if i > math.MaxUint32 {
				return fmt.Errorf("gointernals.ReflectStrToNumBool: value out of range for uint32: %s", src)
			}
			ReflectValueSet(dst, uint32(i))
		case 2:
			if i > math.MaxUint16 {
				return fmt.Errorf("gointernals.ReflectStrToNumBool: value out of range for uint16: %s", src)
			}
			ReflectValueSet(dst, uint16(i))
		case 1:
			if i > math.MaxUint8 {
				return fmt.Errorf("gointernals.ReflectStrToNumBool: value out of range for uint8: %s", src)
			}
			ReflectValueSet(dst, uint8(i))
		}
	case dstTKind == reflect.Float32:
		f, err := strconv.ParseFloat(src, 32)
		if err != nil {
			return err
		}
		ReflectValueSet(dst, float32(f))
	case dstTKind == reflect.Float64:
		f, err := strconv.ParseFloat(src, 64)
		if err != nil {
			return err
		}
		ReflectValueSet(dst, f)
	default:
		b, err := strconv.ParseBool(src)
		if err != nil {
			return err
		}
		ReflectValueSet(dst, b)
	}
	return nil
}
