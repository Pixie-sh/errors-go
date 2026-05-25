package errors

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"os"
	"testing"

	"github.com/pixie-sh/logger-go/env"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	errCode1 := ErrorCode{
		Name:      "TEST",
		Value:     1,
		HTTPError: 1,
	}

	errCode2 := ErrorCode{
		Name:      "TEST",
		Value:     2,
		HTTPError: 2,
	}

	os.Setenv(env.DebugMode, "FALSE")
	err1 := New("test %s", "message #1", errCode1)
	err2 := New("test %s", "message #2", errCode2)

	_ = os.Unsetenv(env.DebugMode)
	err3 := NewWithError(err1, "test %s", "message #3")
	_ = err3.WithNestedError(fmt.Errorf("som go error"))

	assert.Equal(t, errCode1, err1.Code)
	assert.Equal(t, errCode1, err3.Code)
	assert.Equal(t, errCode2, err2.Code)

	assert.NotNil(t, err3.NestedError)
	assert.Nil(t, err1.NestedError)
	assert.Nil(t, err2.NestedError)

	assert.Nil(t, err1.Trace)
	assert.Nil(t, err2.Trace)
	assert.Nil(t, err3.Trace)

	assert.Equal(t, "test message #1", err1.Message)
	assert.Equal(t, "test message #2", err2.Message)
	assert.Equal(t, "test message #3", err3.Message)

	assert.Equal(t, "TEST-1 test message #3; TEST-1 test message #1; som go error", err3.Error())
	assert.Equal(t, "TEST-1 test message #1", err1.Error())
	assert.Equal(t, "TEST-2 test message #2", err2.Error())

	os.Setenv(env.DebugMode, "true")
	err4 := New("test %s", "message #4").WithErrorCode(errCode2)

	os.Setenv(env.DebugMode, "1")
	err5 := NewWithError(err4, "test %s", "message #5")
	err5.WithNestedError(err3)
	err5.WithNestedError(fmt.Errorf("som go error"))

	assert.NotNil(t, err4.Trace)
	assert.NotNil(t, err5.Trace)
	assert.Equal(t, 3, len(err5.NestedError))

	noCodeErr := New("no error %s", "one")
	assert.Equal(t, "no error one", noCodeErr.Error())

	blob, err := err5.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, blob)
	assert.NotEmpty(t, blob)

	var err51 Error
	err = json.Unmarshal(blob, &err51)
	assert.Nil(t, err)
	assert.Equal(t, err5.Code, err51.Code)
	assert.Equal(t, "TEST-2", err51.Code.String())

	originalBlob, _ := json.Marshal(err51)
	assert.Equal(t, originalBlob, blob)

	err6 := NewWithError(nil, "some %s", "a")
	err6.WithNestedError(nil)
	assert.Empty(t, err6.NestedError)
	assert.Equal(t, err6.Code, UnknownErrorCode)
}

func TestErrorAs(t *testing.T) {
	var e error
	e = NewValidationError("test")

	err, ok := As(e)
	assert.True(t, ok)
	assert.NotNil(t, err)

	err, ok = Has(e, InvalidFormDataCode)
	assert.True(t, ok)
	assert.NotNil(t, err)

}

func TestNewWithVariations(t *testing.T) {
	// Test New with different error codes
	err1 := New("simple error")
	assert.Equal(t, UnknownErrorCode, err1.Code)
	assert.Equal(t, "simple error", err1.Message)

	customCode := ErrorCode{Name: "CUSTOM", Value: 123, HTTPError: 400}
	err2 := New("error with custom code", customCode)
	assert.Equal(t, customCode, err2.Code)
	assert.Equal(t, "error with custom code", err2.Message)
	assert.Equal(t, "CUSTOM-123 error with custom code", err2.Error())
	assert.Equal(t, 400, err2.GetHTTPStatus())

	// Test New with format strings
	err3 := New("error with %s and %d", "string", 42)
	assert.Equal(t, "error with string and 42", err3.Message)

	// Test New with debug mode on/off
	os.Setenv(env.DebugMode, "TRUE")
	debugErr := New("debug error")
	assert.NotNil(t, debugErr.Trace)
	assert.NotEmpty(t, debugErr.Trace.Trace)
	assert.NotEmpty(t, debugErr.Trace.CallerPath)

	os.Setenv(env.DebugMode, "FALSE")
	nonDebugErr := New("non-debug error")
	assert.Nil(t, nonDebugErr.Trace)
}

func TestWrapWithVariations(t *testing.T) {
	// Test with nil error
	err1 := NewWithError(nil, "wrapping nil")
	assert.Equal(t, "wrapping nil", err1.Message)
	assert.Empty(t, err1.NestedError)

	// Test with standard Go error
	stdErr := fmt.Errorf("standard error")
	err2 := NewWithError(stdErr, "wrapped %s", "message")
	assert.Equal(t, "wrapped message", err2.Message)
	assert.Len(t, err2.NestedError, 1)
	assert.Equal(t, stdErr, err2.NestedError[0])

	// Test with our custom error
	customErr := New("inner error", ErrorCode{Name: "INNER", Value: 100})
	err3 := NewWithError(customErr, "outer error")
	assert.Equal(t, "outer error", err3.Message)
	assert.Equal(t, customErr.Code, err3.Code) // Code should be inherited
	assert.Len(t, err3.NestedError, 1)

	// Test with multiple nested errors
	err4 := New("error level 1")
	err4.WithNestedError(
		fmt.Errorf("plain error 1"),
		New("custom error"),
		fmt.Errorf("plain error 2"),
	)
	assert.Len(t, err4.NestedError, 3)

	// Test stacking multiple levels
	level1 := New("level 1")
	level2 := NewWithError(level1, "level 2")
	level3 := NewWithError(level2, "level 3")
	assert.Equal(t, "level 3", level3.Message)
	assert.Len(t, level3.NestedError, 1)
	assert.Equal(t, level2, level3.NestedError[0])

	// Test with error code overriding
	baseErr := New("base", ErrorCode{Name: "BASE", Value: 1})
	overrideErr := NewWithError(baseErr, "override")
	finalErr := overrideErr.WithErrorCode(ErrorCode{Name: "FINAL", Value: 999})
	assert.Equal(t, "FINAL-999 override; BASE-1 base", finalErr.Error())
}

func TestValidationErrorVariations(t *testing.T) {
	// Simple validation error without field errors
	valErr1 := NewValidationError("validation failed")
	assert.Equal(t, InvalidFormDataCode, valErr1.Code)
	assert.Equal(t, "validation failed", valErr1.Message)
	assert.Empty(t, valErr1.FieldErrors)

	// Validation error with field errors
	fieldErr1 := &FieldError{
		Field:   "username",
		Rule:    "required",
		Param:   "",
		Message: "Username is required",
	}

	fieldErr2 := &FieldError{
		Field:   "email",
		Rule:    "format",
		Param:   "email",
		Message: "Invalid email format",
	}

	valErr2 := NewValidationError("validation failed with fields")
	valErr2.FieldErrors = []*FieldError{fieldErr1, fieldErr2}

	assert.Equal(t, InvalidFormDataCode, valErr2.Code)
	assert.Len(t, valErr2.FieldErrors, 2)
	assert.Equal(t, "username", valErr2.FieldErrors[0].Field)
	assert.Equal(t, "email", valErr2.FieldErrors[1].Field)

	// Create validation error with custom error code
	customErrCode := ErrorCode{Name: "CUSTOM_VALIDATION", Value: 422, HTTPError: 422}
	valErr3 := NewValidationError("custom validation").WithErrorCode(customErrCode)
	assert.Equal(t, customErrCode, valErr3.Code)

	// Test JSON serialization of validation error
	jsonData, err := valErr2.MarshalJSON()
	assert.Nil(t, err)

	var unmarshalledErr Error
	err = json.Unmarshal(jsonData, &unmarshalledErr)
	assert.Nil(t, err)
	assert.Equal(t, valErr2.Code, unmarshalledErr.Code)
	assert.Equal(t, valErr2.Message, unmarshalledErr.Message)
	assert.Len(t, unmarshalledErr.FieldErrors, 2)
	assert.Equal(t, "username", unmarshalledErr.FieldErrors[0].Field)
	assert.Equal(t, "email", unmarshalledErr.FieldErrors[1].Field)

	// Test validation error with nested errors
	baseErr := fmt.Errorf("database error")
	valErr4 := NewValidationError("validation with nested")
	valErr4.WithNestedError(baseErr)
	assert.Len(t, valErr4.NestedError, 1)

	// Test As and Has with validation errors
	var e error = valErr4
	detected, ok := As(e)
	assert.True(t, ok)
	assert.NotNil(t, detected)

	detected, ok = Has(e, InvalidFormDataCode)
	assert.True(t, ok)
	assert.NotNil(t, detected)

	detected, ok = Has(e, ErrorCode{Name: "WRONG", Value: 999})
	assert.False(t, ok)
	assert.Nil(t, detected)
}

func TestErrorEdgeCases(t *testing.T) {
	err := New("initial error")
	assert.Equal(t, UnknownErrorCode, err.Code)

	code1 := ErrorCode{Name: "CODE1", Value: 1}
	code2 := ErrorCode{Name: "CODE2", Value: 2}

	err.WithErrorCode(code1)
	assert.Equal(t, code1, err.Code)

	err.WithErrorCode(code2)
	assert.Equal(t, code2, err.Code)

	initialNestedCount := len(err.NestedError)
	err.WithNestedError(nil)
	assert.Equal(t, initialNestedCount, len(err.NestedError))

	originalErr := fmt.Errorf("original")
	wrappedOnce := New("wrapped once").WithNestedError(originalErr)
	wrappedTwice := New("wrapped twice").WithNestedError(wrappedOnce, originalErr)
	assert.Len(t, wrappedTwice.NestedError, 2)

	emptyErr := New("")
	jsonData, merr := emptyErr.MarshalJSON()
	assert.Nil(t, merr)

	var unmarshalled Error
	merr = json.Unmarshal(jsonData, &unmarshalled)
	assert.Nil(t, merr)
	assert.Equal(t, emptyErr.Code, unmarshalled.Code)
	assert.Equal(t, emptyErr.Message, unmarshalled.Message)

	baseErr := New("base error")
	currentErr := baseErr
	for i := 0; i < 10; i++ {
		nextErr := New(fmt.Sprintf("level %d", i))
		nextErr.WithNestedError(currentErr)
		currentErr = nextErr
	}

	jsonData, merr = currentErr.MarshalJSON()
	assert.Nil(t, merr)
	assert.NotNil(t, jsonData)

	var unmarshallCurrentErr Error
	merr = json.Unmarshal(jsonData, &unmarshallCurrentErr)
	assert.Nil(t, merr)
	assert.Equal(t, currentErr.Code, unmarshallCurrentErr.Code)
	assert.Equal(t, currentErr.Message, unmarshallCurrentErr.Message)
	assert.Equal(t, currentErr.NestedError, unmarshallCurrentErr.NestedError)
}

// ErrorCode value layout: NewErrorCode validates that value%1000 is a real
// HTTP status, so the test codes below use 10NNN suffixes mapping to known
// codes (200, 201, 202, ...). This matches the convention in error_join_test.go.

// TestWrapExplicitCodeOverridesNested verifies AC-1: an explicit ErrorCode
// passed to Wrap is preserved on the resulting *Error, even when the wrapped
// error is itself an *Error with a different code.
func TestWrapExplicitCodeOverridesNested(t *testing.T) {
	innerCode := NewErrorCode("INNER_T1", 10200)
	outerCode := NewErrorCode("OUTER_T1", 20201)
	inner := New("inner", innerCode)

	wrapped := Wrap(inner, "wrapped", outerCode)

	assert.Equal(t, outerCode, wrapped.Code)
	assert.Len(t, wrapped.NestedError, 1)
	assert.Equal(t, inner, wrapped.NestedError[0])
}

// TestWrapInheritsNestedCodeWhenNoExplicit verifies AC-2: when no explicit
// ErrorCode is provided and the wrapped error is an *Error, the nested code
// is inherited (regression guard for prior behavior).
func TestWrapInheritsNestedCodeWhenNoExplicit(t *testing.T) {
	innerCode := NewErrorCode("INNER_T2", 10202)
	inner := New("inner", innerCode)

	wrapped := Wrap(inner, "wrapped")

	assert.Equal(t, innerCode, wrapped.Code)
}

// TestWrapFallbackForStdError verifies AC-3: wrapping a non-*Error with no
// explicit ErrorCode falls back to UnknownErrorCode.
func TestWrapFallbackForStdError(t *testing.T) {
	stdErr := fmt.Errorf("plain")

	wrapped := Wrap(stdErr, "wrapped")

	assert.Equal(t, UnknownErrorCode, wrapped.Code)
	assert.Len(t, wrapped.NestedError, 1)
	assert.Equal(t, stdErr, wrapped.NestedError[0])
}

// TestErrorIsGoCompat verifies AC-4/AC-5: stdlib errors.Is matches *Error
// values by ErrorCode across single-nested and multi-nested (joined) chains,
// and falls through to Unwrap-based traversal for non-*Error targets.
func TestErrorIsGoCompat(t *testing.T) {
	codeA := NewErrorCode("CODE_A", 10203)
	codeB := NewErrorCode("CODE_B", 10204)

	sentinel := New("sentinel", codeA)

	// Single-nested wrap with code preservation (inherited via newWithArgs).
	wrapped := Wrap(sentinel, "ctx")
	assert.True(t, goerrors.Is(wrapped, sentinel))

	// Explicit override: outer code != sentinel code; sentinel is still found
	// via Unwrap traversal of the nested chain.
	overridden := Wrap(sentinel, "override", codeB)
	assert.True(t, goerrors.Is(overridden, sentinel))

	// Joined / multi-nested: both branches are reachable via Has traversal
	// from inside the custom Is method.
	other := New("other", codeB)
	joined := Join(sentinel, other)
	assert.True(t, goerrors.Is(joined, sentinel))
	assert.True(t, goerrors.Is(joined, other))

	// Non-*Error target: our Is returns false, stdlib continues via Unwrap()
	// and matches by pointer identity.
	plain := fmt.Errorf("plain")
	wrappedPlain := Wrap(plain, "ctx")
	assert.True(t, goerrors.Is(wrappedPlain, plain))

	// Nil target: stdlib short-circuits without invoking our Is.
	assert.False(t, goerrors.Is(wrapped, nil))

	// Non-matching sentinel: not found anywhere in the chain.
	unrelated := New("unrelated", NewErrorCode("UNRELATED", 10205))
	assert.False(t, goerrors.Is(wrapped, unrelated))
}

// TestErrorAsGoCompat verifies AC-5: stdlib errors.As resolves an *Error
// target through a single-nested chain (regression guard).
func TestErrorAsGoCompat(t *testing.T) {
	inner := New("inner", NewErrorCode("AS_TEST", 10206))
	wrapped := Wrap(inner, "wrapped")

	var got *Error
	assert.True(t, goerrors.As(wrapped, &got))
	assert.NotNil(t, got)
}

// TestWrapWVerbRendersInline verifies %w in the format string renders the
// wrapped error's message inline, matching fmt.Errorf semantics. The wrapped
// error is also captured in NestedError and reachable via Unwrap().
func TestWrapWVerbRendersInline(t *testing.T) {
	inner := fmt.Errorf("disk full")

	wrapped := Wrap(inner, "save failed: %w")

	assert.Equal(t, "save failed: disk full", wrapped.Message)
	assert.Len(t, wrapped.NestedError, 1)
	assert.Equal(t, inner, wrapped.NestedError[0])
	assert.Equal(t, inner, goerrors.Unwrap(wrapped))
	assert.True(t, goerrors.Is(wrapped, inner))
}

// TestNewWVerbRendersInline verifies %w works via New as well.
func TestNewWVerbRendersInline(t *testing.T) {
	inner := fmt.Errorf("connection refused")

	err := New("dial failed: %w", inner)

	assert.Equal(t, "dial failed: connection refused", err.Message)
	assert.Len(t, err.NestedError, 1)
	assert.Equal(t, inner, goerrors.Unwrap(err))
}

// TestWVerbWithErrorCode verifies %w coexists with explicit ErrorCode args.
func TestWVerbWithErrorCode(t *testing.T) {
	code := NewErrorCode("W_VERB_CODE", 10300)
	inner := fmt.Errorf("upstream timeout")

	err := New("call failed: %w", inner, code)

	assert.Equal(t, code, err.Code)
	assert.Equal(t, "call failed: upstream timeout", err.Message)
	assert.Equal(t, inner, goerrors.Unwrap(err))
}

// TestMultipleWVerbs verifies multiple %w directives render all referenced
// errors and capture all of them in NestedError, preserving order.
func TestMultipleWVerbs(t *testing.T) {
	a := fmt.Errorf("alpha")
	b := fmt.Errorf("beta")

	err := New("combined: %w and %w", a, b)

	assert.Equal(t, "combined: alpha and beta", err.Message)
	assert.Len(t, err.NestedError, 2)
	assert.Equal(t, a, err.NestedError[0])
	assert.Equal(t, b, err.NestedError[1])
	// Unwrap() returns NestedError[0] per the library contract.
	assert.Equal(t, a, goerrors.Unwrap(err))

	jsonb, _ := json.Marshal(err)
	fmt.Println(string(jsonb))
}

// TestNoWVerbBackwardCompat ensures calls without %w behave identically to
// the pre-change code path.
func TestNoWVerbBackwardCompat(t *testing.T) {
	inner := fmt.Errorf("root cause")

	err := Wrap(inner, "operation %s", "failed")

	assert.Equal(t, "operation failed", err.Message)
	assert.Len(t, err.NestedError, 1)
	assert.Equal(t, inner, err.NestedError[0])
	assert.Equal(t, inner, goerrors.Unwrap(err))
}
