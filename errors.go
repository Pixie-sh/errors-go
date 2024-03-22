package errors

import (
	"encoding/json"
	"fmt"
	"github.com/pixie-sh/logger-go/caller"
	"github.com/pixie-sh/logger-go/env"
	"github.com/pixie-sh/logger-go/logger"
	"runtime/debug"
	"strconv"
	"strings"
)

// ErrorCode error code used by Error to specify some known error
type ErrorCode struct {
	Name      string `json:"name"`
	Value     int    `json:"value"`
	HTTPError int    `json:"-"`
}

// E Error pointer
type E = *Error

// Error error struct to be used
type Error struct {
	Code        ErrorCode     `json:"code"`
	Message     string        `json:"message,omitempty"`
	Trace       *StackTrace   `json:"stack_trace,omitempty"`
	NestedError []*Error      `json:"nested_error,omitempty"`
	FieldErrors []*FieldError `json:"field_errors,omitempty"`
}

// StackTrace trace from debug.Stack with Caller information
type StackTrace struct {
	Trace      string `json:"trace"`
	CallerPath string `json:"caller,omitempty"`
}

// FieldError contains info regarding the field error.
type FieldError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Param   string `json:"rule_param"`
	Message string `json:"message"`
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
	return newWithCallerDepth(caller.TwoHopsCallerDepth, NoErrorCode, format, messages...)
}

// NewWithoutStackTrace creates a error based on messages without trace
func NewWithoutStackTrace(format string, messages ...interface{}) E {
	return &Error{
		Message: fmt.Sprintf(format, messages...),
	}
}

// NewWithError returns a new error with a nested one. uses the nested error code
func NewWithError(err error, format string, messages ...interface{}) E {
	if err == nil {
		return newWithCallerDepth(caller.TwoHopsCallerDepth, NoErrorCode, format, messages...)
	}

	castedErr, ok := err.(E)
	code := NoErrorCode
	if ok {
		code = castedErr.Code
	}

	err2 := newWithCallerDepth(caller.TwoHopsCallerDepth, code, format, messages...)
	return err2.WithNestedError(err)
}

func newWithCallerDepth(depth caller.Depth, code ErrorCode, format string, messages ...interface{}) E {
	var st *StackTrace = nil
	if env.IsDebugActive() {
		st = &StackTrace{
			Trace:      string(debug.Stack()),
			CallerPath: caller.NewCaller(depth).String(),
		}
	}

	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, messages...),
		Trace:   st,
	}
}

// NewValidationFailure returns an error formatted with validations errors
func NewValidationFailure(field string, rule string, message string) E {
	errorResult := new(Error)

	errorResult.Code = InvalidFormDataCode
	errorResult.Message = "Invalid form sent."

	failure := new(FieldError)
	failure.Rule = rule
	failure.Field = field
	failure.Message = message

	errorResult.FieldErrors = append(errorResult.FieldErrors, failure)

	return errorResult
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Code == NoErrorCode {
		return e.Message
	}

	return e.Code.String()
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
func (e *Error) WithNestedError(err error) E {
	if err == nil {
		return e
	}

	if e.NestedError == nil {
		e.NestedError = []*Error{}
	}

	errToAppend, ok := err.(E)
	if !ok {
		errToAppend = New(err.Error())
	}

	e.NestedError = append(e.NestedError, errToAppend)
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
		ec.Name = NoErrorCode.Name
		ec.Value = NoErrorCode.Value
		return nil
	}

	value, err := strconv.ParseInt(codeParts[1], 10, 64)
	if err != nil {
		logger.Logger.Warn("unable to parse error code for %s. using default", mErr)
		ec.Name = NoErrorCode.Name
		ec.Value = NoErrorCode.Value
		return nil
	}

	ec.Name = codeParts[0]
	ec.Value = int(value)
	return nil
}

func (ec *ErrorCode) String() string {
	return fmt.Sprintf("%s-%d", ec.Name, ec.Value)
}
