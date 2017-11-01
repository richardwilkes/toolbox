package errs_test

import (
	"errors"
	"fmt"
	"testing"

	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleError() {
	var bad *int
	func() int {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					fmt.Println(errs.NewWithCause("Caught panic", err))
				}
			}
		}()
		return *bad // trigger a panic due to a nil pointer dereference
	}()
}

func TestAppendError(t *testing.T) {
	original := errs.New("foo")
	result := errs.Append(original, errs.New("bar"))
	require.Equal(t, result.Count(), 2, "Wrong length")
}

func TestAppendToEmptyError(t *testing.T) {
	original := &errs.Error{}
	result := errs.Append(original, errs.New("bar"))
	require.Equal(t, result.Count(), 1, "Wrong length")
}

func TestAppendFlattening(t *testing.T) {
	original := errs.New("foo")
	result := errs.Append(original, errs.Append(nil, errs.New("foo"), errs.New("bar")))
	require.Equal(t, result.Count(), 3, "Wrong length")
}

func TestAppendTypedNil(t *testing.T) {
	var e *errs.Error
	result := errs.Append(e, errs.New("bar"))
	require.Equal(t, result.Count(), 1, "Wrong length")
}

func TestAppendNilError(t *testing.T) {
	var err error
	result := errs.Append(err, errs.New("bar"))
	require.Equal(t, result.Count(), 1, "Wrong length")
}

func TestAppendNilErrorArg(t *testing.T) {
	var err error
	var nilErr *errs.Error
	result := errs.Append(err, nilErr)
	require.Equal(t, result.Count(), 0, "Wrong length")
}

func TestAppendNilErrorInterfaceArg(t *testing.T) {
	var err error
	var nilErr error
	result := errs.Append(err, nilErr)
	require.Equal(t, result.Count(), 0, "Wrong length")
}

func TestAppendNonErrorError(t *testing.T) {
	original := errors.New("foo")
	result := errs.Append(original, errs.New("bar"))
	require.Equal(t, result.Count(), 2, "Wrong length")
}

func TestAppendNonErrorErrorWithAppend(t *testing.T) {
	original := errors.New("foo")
	result := errs.Append(original, errs.Append(nil, errors.New("bar")))
	require.Equal(t, result.Count(), 2, "Wrong length")
}

func TestErrorOrNil(t *testing.T) {
	var err errs.Error
	require.Nil(t, err.ErrorOrNil(), "Should have been nil")
}

func TestErrorOrNilPointer(t *testing.T) {
	var err *errs.Error
	require.Nil(t, err.ErrorOrNil(), "Should have been nil")
}

func TestWrap(t *testing.T) {
	notError := errors.New("foo")
	result := errs.Wrap(notError)
	require.Equal(t, 1, strings.Count(result.Error(), "\n"), "Should have 1 embedded return")
}

func TestNew(t *testing.T) {
	result := errs.New("foo")
	require.Equal(t, 1, strings.Count(result.Error(), "\n"), "Should have 1 embedded return")
}

func TestNewWithCause(t *testing.T) {
	cause := errs.New("bar")
	result := errs.NewWithCause("foo", cause)
	require.Equal(t, 3, strings.Count(result.Error(), "\n"), "Should have 3 embedded returns")
}

func TestFormat(t *testing.T) {
	err := errs.New("test")
	assert.Equal(t, "test", fmt.Sprintf("%s", err))
	assert.Equal(t, `"test"`, fmt.Sprintf("%q", err))
	result := fmt.Sprintf("%v", err)
	assert.Contains(t, result, "[github.com/richardwilkes/toolbox/errs_test.TestFormat]")
	assert.NotContains(t, result, "[runtime.goexit]")
	result = fmt.Sprintf("%+v", err)
	assert.Contains(t, result, "[github.com/richardwilkes/toolbox/errs_test.TestFormat]")
	assert.Contains(t, result, "[runtime.goexit]")
}
