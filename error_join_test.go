package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoin(t *testing.T) {
	t.Run("should return nil when no errors are passed", func(t *testing.T) {
		err := Join()
		assert.Nil(t, err)
	})

	t.Run("should return nil when all nil errors are passed", func(t *testing.T) {
		err := Join(nil, nil)
		assert.Nil(t, err)
	})

	t.Run("should return the error when only one non-nil error is passed", func(t *testing.T) {
		singleErr := New("single error")
		err := Join(singleErr)
		assert.Equal(t, singleErr, err)
	})

	t.Run("should return the non-nil error when mixed with nil errors", func(t *testing.T) {
		singleErr := New("single error")
		err := Join(nil, singleErr, nil)
		assert.Equal(t, singleErr, err)
	})

	t.Run("should join multiple errors under JoinedErrorCode", func(t *testing.T) {
		err1 := New("error 1")
		err2 := New("error 2")

		joinedErr := Join(err1, err2)

		// Check it's of type E
		e, ok := As(joinedErr)
		assert.True(t, ok)
		assert.Equal(t, JoinedErrorCode, e.Code)
		assert.Equal(t, "JoinedError-50300 [error 1; error 2]; error 1; error 2", e.Error())

		// Check the nested errors
		assert.Len(t, e.NestedError, 2)
		assert.Contains(t, e.NestedError, err1)
		assert.Contains(t, e.NestedError, err2)
	})
}

func TestHas(t *testing.T) {
	// Define some test error codes
	testCode1 := NewErrorCode("TEST_ERROR_1", 10001)
	testCode2 := NewErrorCode("TEST_ERROR_2", 10002)

	t.Run("should return false for nil error", func(t *testing.T) {
		e, ok := Has(nil, testCode1)
		assert.False(t, ok)
		assert.Nil(t, e)
	})

	t.Run("should return false when error doesn't match code", func(t *testing.T) {
		err := New("test error", testCode1)
		e, ok := Has(err, testCode2)
		assert.False(t, ok)
		assert.Nil(t, e)
	})

	t.Run("should return true when error matches code", func(t *testing.T) {
		err := New("test error", testCode1)
		e, ok := Has(err, testCode1)
		assert.True(t, ok)
		assert.Equal(t, err, e)
	})

	t.Run("should find nested error in joined errors", func(t *testing.T) {
		err1 := New("error 1", testCode1)
		err2 := New("error 2", testCode2)

		joinedErr := Join(err1, err2)

		// Find first error
		e1, ok1 := Has(joinedErr, testCode1, true)
		assert.True(t, ok1)
		assert.Equal(t, testCode1, e1.Code)
		assert.Equal(t, "TEST_ERROR_1-10001 error 1", e1.Error())

		// Find second error
		e2, ok2 := Has(joinedErr, testCode2, true)
		assert.True(t, ok2)
		assert.Equal(t, testCode2, e2.Code)
		assert.Equal(t, "TEST_ERROR_2-10002 error 2", e2.Error())
	})

	t.Run("should handle deeply nested errors", func(t *testing.T) {
		innerErr := New("inner error", testCode1)
		middleErr := Join(innerErr, New("other error"))
		outerErr := Join(New("outer error"), middleErr)

		e, ok := Has(outerErr, testCode1, true)
		assert.True(t, ok)
		assert.Equal(t, testCode1, e.Code)
		assert.Equal(t, "TEST_ERROR_1-10001 inner error", e.Error())
	})

	t.Run("should return false when code not present in joined errors", func(t *testing.T) {
		err1 := New("error 1", testCode2)
		err2 := New("error 2", testCode1)

		joinedErr := Join(err1, err2)

		e, ok := Has(joinedErr, testCode1, true)
		assert.True(t, ok)
		assert.NotNil(t, e)
	})
}
