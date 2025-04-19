package utils

import (
	"reflect"
)

// Nil properly checks if the argument is nil.
// More on nil checking here:
// https://mangatmodi.medium.com/go-check-nil-interface-the-right-way-d142776edef1
func Nil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice, reflect.Func:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}

// IsPointer validates if input is pointer
func IsPointer(i interface{}) bool {
	if i == nil {
		return false
	}

	return reflect.TypeOf(i).Kind() == reflect.Ptr
}
