package errors

import (
	"fmt"
	"github.com/pixie-sh/logger-go/caller"
	"github.com/pixie-sh/logger-go/env"
	"runtime/debug"
)

// Depth caller depth type
type Depth = int

// most used depth for caller
const (
	SelfCallerDepth      Depth = 1
	FnCallerDepth        Depth = 2
	TwoHopsCallerDepth   Depth = 3
	ThreeHopsCallerDepth Depth = 4
	FourHopsCallerDepth  Depth = 5
)

// ErrorCode error code used by Error to specify some known error
type ErrorCode struct {
	Name      string `json:"name"`
	Value     int    `json:"value"`
	HTTPError int    `json:"-"`
}

// E Error pointer
type E = *Error

// Error struct to be used
type Error struct {
	Code        ErrorCode     `json:"code"`
	Message     string        `json:"message,omitempty"`
	Trace       *StackTrace    `json:"stack_trace,omitempty"`
	NestedError []error       `json:"nested_error,omitempty"`
	FieldErrors []*FieldError `json:"field_errors,omitempty"`
}

// StackTrace trace from debug.Stack with Caller information
type StackTrace struct {
	Trace      []byte `json:"trace,omitempty"`
	CallerPath string `json:"caller,omitempty"`
}

// FieldError contains info regarding the field error.
type FieldError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Param   string `json:"rule_param"`
	Message string `json:"message"`
}

func newWithCallerDepth(depth Depth, code ErrorCode, format string, messages ...interface{}) E {
	var st *StackTrace
	if env.IsDebugActive() {
		st = &StackTrace{
			Trace:      debug.Stack(),
			CallerPath: caller.NewCaller(depth).String(),
		}
	}

	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, messages...),
		Trace:   st,
	}
}

// GetHTTPStatus get's the http status for the error
func (e *Error) GetHTTPStatus() int {
	return e.Code.HTTPError
}

// WithNestedError add newWithArgs error to current one, as a wrapped error using the native error wrapper
func (e *Error) WithNestedError(errors ...error) E {
	if len(errors) == 0 {
		return e
	}

	if e.NestedError == nil {
		e.NestedError = make([]error, 0)
	}

	for _, err := range errors {
		if err == nil {
			continue
		}

		e.NestedError = append(e.NestedError, err)
	}

	return e
}

// WithErrorCode add code to Error
func (e *Error) WithErrorCode(code ErrorCode) E {
	e.Code = code
	return e
}
