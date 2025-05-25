package errors

import (
	goErrors "errors"
	"github.com/pixie-sh/errors-go/utils"
	"strings"
)

// New Package errors provides utility methods and custom error handling mechanisms,
// including support for error codes, nested errors, validation errors, and more.
// It enhances the standard `errors` package by allowing structured error creation
// and formatting with additional context and metadata.
func New(message string, args ...interface{}) E {
	return newWithArgs(TwoHopsCallerDepth, message, args...)
}

func Wrap(err error, message string, args ...interface{}) E {
	return newWithArgs(ThreeHopsCallerDepth, message, append(args, err)...)
}

func Must(err error) {
	if err != nil {
		e, ok := As(err)
		if ok {
			panic(e)
		}

		panic(err)
	}
}

func As(err error) (E, bool) {
	var e E

	switch {
	case goErrors.As(err, &e):
		return e, true
	}

	return nil, false
}

// Has checks if the given error includes an error with the specified ErrorCode.
// It traverses through nested errors if the error is of a joined type (JoinedErrorCode).
// If the ErrorCode is found, it returns the corresponding error (E) and true; otherwise, it returns nil and false.
func Has(err error, ec ErrorCode, evalNested ...bool) (E, bool) {
	e, valid := As(err)
	if valid && e.Code == ec {
		return e, true
	}

	if valid && e.Code == JoinedErrorCode && len(evalNested) > 0 && evalNested[0] {
		for _, nestedErr := range e.NestedError {
			nestedE, nestedOk := Has(nestedErr, ec, evalNested...)
			if nestedOk {
				return nestedE, true
			}
		}
	}

	return nil, false
}

// NewErrorCode returns an ErrorCode
// HttpCode is the last three digits of Value
func NewErrorCode(name string, value int) ErrorCode {
	httpError := value % 1000
	if httpError < 0 {
		// Ensure ec.HTTPError always represents the last three digits of ec.Value,
		// even when ec.Value is negative
		httpError += 1000
	}

	return ErrorCode{
		Name:      name,
		Value:     value,
		HTTPError: httpError,
	}
}

// NewWithError returns a newWithArgs error with a nested one. uses the nested error code
func NewWithError(err error, format string, args ...interface{}) E {
	return newWithArgs(ThreeHopsCallerDepth, format, append(args, err)...)
}

// NewValidationError returns an error formatted with validations errors
func NewValidationError(message string, fields ...*FieldError) E {
	var args = mapSlice(fields, func(field *FieldError) interface{} {
		return field
	})
	args = append(args, InvalidFormDataCode)

	return newWithArgs(ThreeHopsCallerDepth, message, args...)
}

// Join combines multiple errors into a single Error.
// It returns nil if no non-nil errors are passed.
// If only one non-nil error is passed, it will return that error.
// If all entries are nil, nil is returned
// If only one valid is passed, that one is returned instead of JoinedError
// Otherwise, it creates a new Error that contains all non-nil errors as nested errors.
func Join(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	if len(errs) == 1 {
		return errs[0]
	}

	baseErr := &Error{
		Code:        JoinedErrorCode,
		NestedError: make([]error, 0),
	}

	messageBuilder := strings.Builder{}
	for i := 0; i < len(errs); {
		if errs[i] == nil {
			errs = append(errs[:i], errs[i+1:]...)
			continue
		}

		if i == 0 {
			messageBuilder.WriteString("[")
		}
		messageBuilder.WriteString(errs[i].Error())
		if i < len(errs)-1 {
			messageBuilder.WriteString("; ")
		} else {
			messageBuilder.WriteString("]")
		}

		baseErr.NestedError = append(baseErr.NestedError, errs[i])
		i++
	}

	if len(baseErr.NestedError) == 0 {
		return nil
	}

	if len(baseErr.NestedError) == 1 {
		return baseErr.NestedError[0]
	}

	baseErr.Message = messageBuilder.String()
	return baseErr
}

func newWithArgs(depth Depth, message string, args ...interface{}) E {
	var code = UnknownErrorCode
	var fields []*FieldError
	var toWrap error

	for i := 0; i < len(args); {
		if args[i] == nil {
			args = append(args[:i], args[i+1:]...)
			continue
		}

		switch v := args[i].(type) {
		case ErrorCode:
			code = v
			args = append(args[:i], args[i+1:]...)
		case FieldError:
			fields = append(fields, &v)
			args = append(args[:i], args[i+1:]...)
		case *FieldError:
			fields = append(fields, v)
			args = append(args[:i], args[i+1:]...)
		case error:
			toWrap = v
			args = append(args[:i], args[i+1:]...)
		default:
			i++
		}
	}

	e := newWithCallerDepth(depth, code, message, args...)
	if len(fields) > 0 {
		e.FieldErrors = fields
		if code == UnknownErrorCode {
			e.Code = InvalidFormDataCode
		}
	}

	if toWrap != nil {
		toWrapCasted, ok := As(toWrap)
		if ok {
			_ = e.WithErrorCode(toWrapCasted.Code)
		}

		e = e.WithNestedError(toWrap)
	}

	return e
}

func mapSlice[S ~[]E, E any, R any](model S, f func(item E) R) []R {
	var result []R
	for _, item := range model {
		res := f(item)
		if !utils.Nil(res) {
			result = append(result, res)
		}
	}

	return result
}
