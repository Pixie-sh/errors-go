package errors

import (
	"encoding/json"
	"fmt"
	"github.com/pixie-sh/logger-go/env"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestErrors(t *testing.T) {
	errCode1 := ErrorCode{
		Name:  "TEST",
		Value: 1,
	}

	errCode2 := ErrorCode{
		Name:  "TEST",
		Value: 2,
	}

	os.Setenv(env.DebugMode, "FALSE")
	err1 := New("test %s", "message #1").WithErrorCode(errCode1)
	err2 := New("test %s", "message #2").WithErrorCode(errCode2)

	_ = os.Unsetenv(env.DebugMode)
	err3 := NewWithError(err1, "test %s", "message #3")

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

	assert.Equal(t, "TEST-1", err3.Error())
	assert.Equal(t, "TEST-1", err1.Error())
	assert.Equal(t, "TEST-2", err2.Error())

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

	blob := err5.JSON()
	fmt.Println(string(blob))

	var err51 Error
	err := json.Unmarshal(blob, &err51)
	assert.Nil(t, err)
	assert.Equal(t, err5.Code, err51.Code)
	assert.Equal(t, "TEST-2", err51.Code.String())

	originalBlob, _ := json.Marshal(err5)
	assert.Equal(t, originalBlob, blob)

	err6 := NewWithError(nil, "some %s", "a")
	err6.WithNestedError(nil)
	assert.Empty(t, err6.NestedError)
	assert.Equal(t, err6.Code, GenericErrorCode)
}
