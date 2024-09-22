package gorule

import "reflect"

func callInnerLen(val interface{}) int {
	return reflect.Indirect(reflect.ValueOf(val)).Len()
}

func callInnerCap(val interface{}) int {
	return reflect.Indirect(reflect.ValueOf(val)).Cap()
}
