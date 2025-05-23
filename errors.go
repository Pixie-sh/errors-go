package errors

import (
	"encoding/json"
	goErrors "errors"
	"fmt"
	"github.com/pixie-sh/logger-go/caller"
	"github.com/pixie-sh/logger-go/env"
	"github.com/pixie-sh/logger-go/logger"
	"runtime/debug"
	"strconv"
	"strings"
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
	Trace       *StackTrace   `json:"stack_trace,omitempty"`
	NestedError []*Error      `json:"nested_error,omitempty"`
	FieldErrors []*FieldError `json:"field_errors,omitempty"`
}

// StackTrace trace from debug.Stack with Caller information
type StackTrace struct {
	Trace      []byte `json:"trace"`
	CallerPath string `json:"caller,omitempty"`
}

// FieldError contains info regarding the field error.
type FieldError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Param   string `json:"rule_param"`
	Message string `json:"message"`
}

func As(err error) (E, bool) {
	var e E

	switch {
	case goErrors.As(err, &e):
		return e, true
	}

	return nil, false
}

func Has(err error, ec ErrorCode) (E, bool) {
	e, valid := As(err)
	if valid && e.Code == ec {
		return e, true
	}

	return nil, false
}

// NewErrorCode returns an ErrorCode
func NewErrorCode(name string, value int, httpCode int) ErrorCode {
	return ErrorCode{
		Name:      name,
		Value:     value,
		HTTPError: httpCode,
	}
}

// New creates a error based on messages
func New(format string, messages ...interface{}) E {
	return newWithCallerDepth(TwoHopsCallerDepth, GenericErrorCode, format, messages...)
}

// NewWithoutStackTrace creates a error based on messages without trace
func NewWithoutStackTrace(format string, messages ...interface{}) E {
	return &Error{
		Message: fmt.Sprintf(format, messages...),
		Code:    GenericErrorCode,
	}
}

// NewWithError returns a new error with a nested one. uses the nested error code
func NewWithError(err error, format string, messages ...interface{}) E {
	if err == nil {
		return newWithCallerDepth(TwoHopsCallerDepth, GenericErrorCode, format, messages...)
	}

	castedErr, ok := As(err)
	code := UnknownErrorCode
	if ok {
		code = castedErr.Code
	}

	err2 := newWithCallerDepth(TwoHopsCallerDepth, code, format, messages...)
	return err2.WithNestedError(err)
}

// NewWithCallerDepth returns a new error with caller depth
func NewWithCallerDepth(depth Depth, format string, messages ...interface{}) E {
	return newWithCallerDepth(depth, GenericErrorCode, format, messages...)
}

func newWithCallerDepth(depth Depth, code ErrorCode, format string, messages ...interface{}) E {
	var st *StackTrace = nil
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

// NewValidationError returns an error formatted with validations errors
func NewValidationError(message string, fields ...*FieldError) E {
	var errorResult Error
	errorResult.Code = InvalidFormDataCode
	errorResult.Message = message
	errorResult.FieldErrors = fields
	return &errorResult
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Code == GenericErrorCode {
		return e.Message
	}

	return fmt.Sprintf("%s: %s", e.Code.String(), e.Message)
}

// JSON return a json string
func (e *Error) JSON() []byte {
	s, _ := json.Marshal(&e)
	return s
}

// GetHTTPStatus get's the http status for the error
func (e *Error) GetHTTPStatus() int {
	return e.Code.HTTPError
}

// WithNestedError add new error to current one, ignores the nested error code
func (e *Error) WithNestedError(errors ...error) E {
	if len(errors) == 0 {
		return e
	}

	if e.NestedError == nil {
		e.NestedError = []*Error{}
	}

	for _, err := range errors {
		if err == nil {
			continue
		}

		errToAppend, ok := err.(E)
		if !ok {
			errToAppend = New(err.Error())
		}

		e.NestedError = append(e.NestedError, errToAppend)
	}

	return e
}

// WithErrorCode add code to Error
func (e *Error) WithErrorCode(code ErrorCode) E {
	e.Code = code
	return e
}

// MarshalJSON implement json marshaller interface
func (ec *ErrorCode) MarshalJSON() ([]byte, error) {
	c := ec.String()
	return json.Marshal(&c)
}

// UnmarshalJSON implement json marshaller interface
func (ec *ErrorCode) UnmarshalJSON(data []byte) error {
	var mErr string
	err := json.Unmarshal(data, &mErr)
	if err != nil {
		return err
	}

	codeParts := strings.Split(mErr, "-")
	if len(codeParts) != 2 {
		logger.Logger.Warn("unable to parse error code for %s. using default", mErr)
		ec.Name = GenericErrorCode.Name
		ec.Value = GenericErrorCode.Value
		return nil
	}

	value, err := strconv.ParseInt(codeParts[1], 10, 64)
	if err != nil {
		logger.Logger.Warn("unable to parse error code for %s. using default", mErr)
		ec.Name = GenericErrorCode.Name
		ec.Value = GenericErrorCode.Value
		return nil
	}

	ec.Name = codeParts[0]
	ec.Value = int(value)
	ec.HTTPError = ec.Value % 1000
	return nil
}

func (ec *ErrorCode) String() string {
	return fmt.Sprintf("%s-%d", ec.Name, ec.Value)
}
