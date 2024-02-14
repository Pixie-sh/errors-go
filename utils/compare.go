package utils

import (
	"encoding/json"
	"fmt"
	"github.com/pixie-sh/errors-go"
	"os"
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

// Must validate error, panic if error
func Must(err error, format string, attr ...interface{}) {
	if err != nil {
		errStr := err.Error()

		castedErr, ok := err.(errors.E)
		if ok {
			blob, _ := json.Marshal(castedErr)
			errStr = string(blob)
		}

		fmt.Fprintf(os.Stderr, "%s\n %s", fmt.Sprintf(format, attr...), errStr)
		os.Exit(1)
	}
}

// IsPointer validates if input is pointer
func IsPointer(i interface{}) bool {
	if i == nil {
		return false
	}

	return reflect.TypeOf(i).Kind() == reflect.Ptr
}
