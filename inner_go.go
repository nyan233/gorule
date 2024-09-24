package gorule

import (
	"fmt"
	"reflect"
)

func callInnerLen(val interface{}) int {
	return reflect.Indirect(reflect.ValueOf(val)).Len()
}

func callInnerCap(val interface{}) int {
	return reflect.Indirect(reflect.ValueOf(val)).Cap()
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func callTypeConvertNumber[T number](val interface{}) T {
	iv := reflect.Indirect(reflect.ValueOf(val))
	switch iv.Interface().(type) {
	case int, int8, int16, int32, int64:
		return T(iv.Int())
	case uint, uint8, uint16, uint32, uint64:
		return T(iv.Uint())
	case float32, float64:
		return T(iv.Float())
	default:
		var d T
		panic(fmt.Sprintf("old type(%s) not convert new type(%s)", reflect.TypeOf(val), reflect.TypeOf(d)))
	}
}
