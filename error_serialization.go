package errors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/pixie-sh/logger-go/env"
	"github.com/pixie-sh/logger-go/logger"
)

// MarshalJSON implement json marshaller interface
func (ec ErrorCode) MarshalJSON() ([]byte, error) {
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
		ec.HTTPError = GenericErrorCode.HTTPError
		return nil
	}

	value, err := strconv.ParseInt(codeParts[1], 10, 64)
	if err != nil {
		logger.Logger.Warn("unable to parse error code for %s. using default", mErr)
		ec.Name = GenericErrorCode.Name
		ec.Value = GenericErrorCode.Value
		ec.HTTPError = GenericErrorCode.HTTPError
		return nil
	}

	ec.Name = codeParts[0]
	ec.Value = int(value)
	ec.HTTPError = ec.Value % 1000
	if ec.HTTPError < 0 {
		// Ensure ec.HTTPError always represents the last three digits of ec.Value,
		// even when ec.Value is negative
		ec.HTTPError += 1000
	}

	return nil
}

func (ec ErrorCode) String() string {
	return fmt.Sprintf("%s-%d", ec.Name, ec.Value)
}

// Error implements the error interface
func (e Error) Error() string {
	var ecStr = e.Code.String() + " "
	if e.Code == GenericErrorCode || e.Code == UnknownErrorCode {
		ecStr = ""
	}

	if len(e.NestedError) == 0 {
		return fmt.Sprintf("%s%s", ecStr, e.Message)
	}

	errStr := strings.Builder{}
	errStr.WriteString(fmt.Sprintf("%s%s", ecStr, e.Message))
	for _, err := range e.NestedError {
		errStr.WriteString(fmt.Sprintf("; %s", err.Error()))
	}

	return errStr.String()
}

// String implements Stringer interface
func (e Error) String() string {
	return e.Error()
}

// Unwrap only returns the first NestedError
// which is the most common use case
func (e Error) Unwrap() error {
	if len(e.NestedError) == 0 {
		return nil
	}

	return e.NestedError[0]
}

func (e Error) MarshalJSON() ([]byte, error) {
	// Create a custom type for marshaling that won't trigger the MarshalJSON method recursively
	type AliasError struct {
		Code        ErrorCode         `json:"code,omitempty"`
		Message     string            `json:"message,omitempty"`
		NestedError []json.RawMessage `json:"nested_error,omitempty"`
		Trace       *StackTrace       `json:"stack_trace,omitempty"`
		FieldErrors []*FieldError     `json:"field_errors,omitempty"`
	}

	aliasErr := AliasError{
		Code:        e.Code,
		Message:     e.Message,
		FieldErrors: e.FieldErrors,
	}

	if env.IsDebugActive() {
		aliasErr.Trace = e.Trace
	}

	if len(e.NestedError) > 0 {
		aliasErr.NestedError = make([]json.RawMessage, len(e.NestedError))

		for i, nested := range e.NestedError {
			if nested == nil {
				continue
			}

			if customErr, ok := As(nested); ok {
				data, err := json.Marshal(customErr)
				if err != nil {
					return nil, err
				}
				aliasErr.NestedError[i] = data
			} else {
				errStr := nested.Error()
				data, err := json.Marshal(errStr)
				if err != nil {
					return nil, err
				}
				aliasErr.NestedError[i] = data
			}
		}
	}

	return json.Marshal(aliasErr)
}

func (e *Error) UnmarshalJSON(data []byte) error {
	type AliasError struct {
		Code        ErrorCode         `json:"code,omitempty"`
		Message     string            `json:"message,omitempty"`
		NestedError []json.RawMessage `json:"nested_error,omitempty"`
		Trace       *StackTrace       `json:"stack_trace,omitempty"`
		FieldErrors []*FieldError     `json:"field_errors,omitempty"`
	}

	var aliasErr AliasError
	if err := json.Unmarshal(data, &aliasErr); err != nil {
		return err
	}

	e.Code = aliasErr.Code
	e.Message = aliasErr.Message
	e.FieldErrors = aliasErr.FieldErrors
	e.Trace = aliasErr.Trace

	if len(aliasErr.NestedError) > 0 {
		e.NestedError = make([]error, len(aliasErr.NestedError))

		for i, nestedData := range aliasErr.NestedError {
			var customErr Error
			if err := json.Unmarshal(nestedData, &customErr); err == nil {
				e.NestedError[i] = &customErr
				continue
			}

			var errStr string
			if err := json.Unmarshal(nestedData, &errStr); err == nil {
				e.NestedError[i] = fmt.Errorf("%s", errStr)
			}
		}
	}

	return nil
}
