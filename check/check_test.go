// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package check_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

// customError is a simple error type for testing typed nil interfaces
type customError struct{}

func (e *customError) Error() string {
	return "custom error"
}

// mockTestingT implements the TestingT interface for testing purposes
type mockTestingT struct {
	errors []string
	failed bool
}

func newMockTestingT() *mockTestingT {
	return &mockTestingT{}
}

func (m *mockTestingT) Cleanup(_ func()) {
	panic("not implemented")
}

func (m *mockTestingT) Error(args ...any) {
	m.failed = true
	m.errors = append(m.errors, fmt.Sprint(args...))
}

func (m *mockTestingT) Errorf(format string, args ...any) {
	m.failed = true
	m.errors = append(m.errors, fmt.Sprintf(format, args...))
}

func (m *mockTestingT) Fail() {
	m.failed = true
}

func (m *mockTestingT) FailNow() {
	m.failed = true
	panic("FailNow called")
}

func (m *mockTestingT) Failed() bool {
	return m.failed
}

func (m *mockTestingT) Fatal(_ ...any) {
	panic("not implemented")
}

func (m *mockTestingT) Fatalf(_ string, _ ...any) {
	panic("not implemented")
}

func (m *mockTestingT) Helper() {
}

func (m *mockTestingT) Log(_ ...any) {
	panic("not implemented")
}

func (m *mockTestingT) Logf(_ string, _ ...any) {
	panic("not implemented")
}

func (m *mockTestingT) Name() string {
	panic("not implemented")
}

func (m *mockTestingT) Setenv(_, _ string) {
	panic("not implemented")
}

func (m *mockTestingT) Chdir(_ string) {
	panic("not implemented")
}

func (m *mockTestingT) Skip(_ ...any) {
	panic("not implemented")
}

func (m *mockTestingT) SkipNow() {
	panic("not implemented")
}

func (m *mockTestingT) Skipf(_ string, _ ...any) {
	panic("not implemented")
}

func (m *mockTestingT) Skipped() bool {
	panic("not implemented")
}

func (m *mockTestingT) TempDir() string {
	panic("not implemented")
}

func (m *mockTestingT) Context() context.Context {
	return context.Background()
}

func TestNew(t *testing.T) {
	mock := newMockTestingT()
	checker := check.New(mock)
	if checker.TestingT != mock {
		t.Error("Expected TestingT to be set correctly")
	}
}

func TestEqual(t *testing.T) {
	tests := []struct { //nolint:govet // Don't care about field alignment in tests
		name        string
		expected    any
		actual      any
		shouldFail  bool
		expectedMsg string
	}{
		{"equal strings", "hello", "hello", false, ""},
		{"unequal strings", "hello", "world", true, "Expected hello, got world"},
		{"equal integers", 42, 42, false, ""},
		{"unequal integers", 42, 24, true, "Expected 42, got 24"},
		{"equal nil values", nil, nil, false, ""},
		{"nil vs non-nil", nil, "test", true, "Expected <nil>, got test"},
		{"non-nil vs nil", "test", nil, true, "Expected test, got <nil>"},
		{"equal byte slices", []byte{1, 2, 3}, []byte{1, 2, 3}, false, ""},
		{"unequal byte slices", []byte{1, 2, 3}, []byte{1, 2, 4}, true, "Expected [1 2 3], got [1 2 4]"},
		{"byte slice vs nil", []byte{1, 2, 3}, nil, true, "Expected [1 2 3], got <nil>"},
		{"nil byte slices", []byte(nil), []byte(nil), false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.Equal(tt.expected, tt.actual)

			if tt.shouldFail {
				if !mock.failed {
					t.Errorf("Expected Equal to fail for %v vs %v", tt.expected, tt.actual)
				}
				if len(mock.errors) == 0 {
					t.Error("Expected error message to be recorded")
				} else if !strings.Contains(mock.errors[0], tt.expectedMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.expectedMsg, mock.errors[0])
				}
			} else if mock.failed {
				t.Errorf("Expected Equal to pass for %v vs %v, but it failed with: %v", tt.expected, tt.actual, mock.errors)
			}
		})
	}
}

func TestEqualWithMessage(t *testing.T) {
	mock := newMockTestingT()
	checker := check.New(mock)

	checker.Equal("expected", "actual", "custom message")

	if !mock.failed {
		t.Error("Expected Equal to fail")
	}
	if len(mock.errors) == 0 {
		t.Error("Expected error message")
	} else if !strings.Contains(mock.errors[0], "custom message") {
		t.Errorf("Expected error message to contain custom message, got: %s", mock.errors[0])
	}
}

func TestEqualWithFormattedMessage(t *testing.T) {
	mock := newMockTestingT()
	checker := check.New(mock)

	checker.Equal("expected", "actual", "test %d: %s", 1, "formatted")

	if !mock.failed {
		t.Error("Expected Equal to fail")
	}
	if len(mock.errors) == 0 {
		t.Error("Expected error message")
	} else if !strings.Contains(mock.errors[0], "test 1: formatted") {
		t.Errorf("Expected error message to contain formatted message, got: %s", mock.errors[0])
	}
}

func TestNotEqual(t *testing.T) {
	tests := []struct { //nolint:govet // Don't care about field alignment in tests
		name        string
		expected    any
		actual      any
		shouldFail  bool
		expectedMsg string
	}{
		{"unequal strings", "hello", "world", false, ""},
		{"equal strings", "hello", "hello", true, "Expected hello to not be hello"},
		{"unequal integers", 42, 24, false, ""},
		{"equal integers", 42, 42, true, "Expected 42 to not be 42"},
		{"nil vs non-nil", nil, "test", false, ""},
		{"equal nil values", nil, nil, true, "Expected <nil> to not be <nil>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.NotEqual(tt.expected, tt.actual)

			if tt.shouldFail {
				if !mock.failed {
					t.Errorf("Expected NotEqual to fail for %v vs %v", tt.expected, tt.actual)
				}
				if len(mock.errors) == 0 {
					t.Error("Expected error message to be recorded")
				} else if !strings.Contains(mock.errors[0], tt.expectedMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.expectedMsg, mock.errors[0])
				}
			} else if mock.failed {
				t.Errorf("Expected NotEqual to pass for %v vs %v, but it failed with: %v", tt.expected, tt.actual, mock.errors)
			}
		})
	}
}

func TestNil(t *testing.T) {
	tests := []struct { //nolint:govet // Don't care about field alignment in tests
		name       string
		value      any
		shouldFail bool
	}{
		{"nil value", nil, false},
		{"nil pointer", (*int)(nil), false},
		{"nil slice", []int(nil), false},
		{"nil map", map[string]int(nil), false},
		{"nil interface", error(nil), false},
		{"non-nil string", "test", true},
		{"non-nil int", 42, true},
		{"non-nil slice", []int{1, 2, 3}, true},
		{"non-nil map", map[string]int{"key": 1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.Nil(tt.value)

			if tt.shouldFail {
				if !mock.failed {
					t.Errorf("Expected Nil to fail for %v", tt.value)
				}
			} else {
				if mock.failed {
					t.Errorf("Expected Nil to pass for %v, but it failed with: %v", tt.value, mock.errors)
				}
			}
		})
	}
}

func TestNotNil(t *testing.T) {
	tests := []struct { //nolint:govet // Don't care about field alignment in tests
		name       string
		value      any
		shouldFail bool
	}{
		{"non-nil string", "test", false},
		{"non-nil int", 42, false},
		{"non-nil slice", []int{1, 2, 3}, false},
		{"nil value", nil, true},
		{"nil pointer", (*int)(nil), true},
		{"nil slice", []int(nil), true},
		{"nil interface", error(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.NotNil(tt.value)

			if tt.shouldFail {
				if !mock.failed {
					t.Errorf("Expected NotNil to fail for %v", tt.value)
				}
			} else {
				if mock.failed {
					t.Errorf("Expected NotNil to pass for %v, but it failed with: %v", tt.value, mock.errors)
				}
			}
		})
	}
}

func TestTrue(t *testing.T) {
	mock := newMockTestingT()
	checker := check.New(mock)

	// Test passing case
	checker.True(true)
	if mock.failed {
		t.Error("Expected True(true) to pass")
	}

	// Test failing case
	mock = newMockTestingT()
	checker = check.New(mock)
	checker.True(false)
	if !mock.failed {
		t.Error("Expected True(false) to fail")
	}
	if len(mock.errors) == 0 || !strings.Contains(mock.errors[0], "Expected true") {
		t.Errorf("Expected specific error message, got: %v", mock.errors)
	}
}

func TestFalse(t *testing.T) {
	mock := newMockTestingT()
	checker := check.New(mock)

	// Test passing case
	checker.False(false)
	if mock.failed {
		t.Error("Expected False(false) to pass")
	}

	// Test failing case
	mock = newMockTestingT()
	checker = check.New(mock)
	checker.False(true)
	if !mock.failed {
		t.Error("Expected False(true) to fail")
	}
	if len(mock.errors) == 0 || !strings.Contains(mock.errors[0], "Expected false") {
		t.Errorf("Expected specific error message, got: %v", mock.errors)
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		substr     string
		shouldFail bool
	}{
		{"contains substring", "hello world", "world", false},
		{"contains at beginning", "hello world", "hello", false},
		{"contains at end", "hello world", "world", false},
		{"contains empty string", "hello", "", false},
		{"does not contain", "hello world", "xyz", true},
		{"case sensitive", "Hello World", "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.Contains(tt.s, tt.substr)

			if tt.shouldFail {
				if !mock.failed {
					t.Errorf("Expected Contains to fail for %q in %q", tt.substr, tt.s)
				}
			} else {
				if mock.failed {
					t.Errorf("Expected Contains to pass for %q in %q, but it failed with: %v", tt.substr, tt.s, mock.errors)
				}
			}
		})
	}
}

func TestNotContains(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		substr     string
		shouldFail bool
	}{
		{"does not contain", "hello world", "xyz", false},
		{"case sensitive", "Hello World", "hello", false},
		{"contains substring", "hello world", "world", true},
		{"contains empty string", "hello", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.NotContains(tt.s, tt.substr)

			if tt.shouldFail {
				if !mock.failed {
					t.Errorf("Expected NotContains to fail for %q in %q", tt.substr, tt.s)
				}
			} else {
				if mock.failed {
					t.Errorf("Expected NotContains to pass for %q in %q, but it failed with: %v", tt.substr, tt.s, mock.errors)
				}
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		prefix     string
		shouldFail bool
	}{
		{"has prefix", "hello world", "hello", false},
		{"exact match", "hello", "hello", false},
		{"empty prefix", "hello", "", false},
		{"does not have prefix", "hello world", "world", true},
		{"case sensitive", "Hello World", "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.HasPrefix(tt.s, tt.prefix)

			if tt.shouldFail {
				if !mock.failed {
					t.Errorf("Expected HasPrefix to fail for prefix %q in %q", tt.prefix, tt.s)
				}
			} else {
				if mock.failed {
					t.Errorf("Expected HasPrefix to pass for prefix %q in %q, but it failed with: %v", tt.prefix, tt.s, mock.errors)
				}
			}
		})
	}
}

func TestNoPrefix(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		prefix     string
		shouldFail bool
	}{
		{"does not have prefix", "hello world", "world", false},
		{"case sensitive", "Hello World", "hello", false},
		{"has prefix", "hello world", "hello", true},
		{"exact match", "hello", "hello", true},
		{"empty prefix", "hello", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.NoPrefix(tt.s, tt.prefix)

			if tt.shouldFail {
				if !mock.failed {
					t.Errorf("Expected NoPrefix to fail for prefix %q in %q", tt.prefix, tt.s)
				}
			} else {
				if mock.failed {
					t.Errorf("Expected NoPrefix to pass for prefix %q in %q, but it failed with: %v", tt.prefix, tt.s, mock.errors)
				}
			}
		})
	}
}

func TestNoError(t *testing.T) {
	mock := newMockTestingT()
	checker := check.New(mock)

	// Test passing case
	checker.NoError(nil)
	if mock.failed {
		t.Error("Expected NoError(nil) to pass")
	}

	// Test failing case
	mock = newMockTestingT()
	checker = check.New(mock)
	testErr := errors.New("test error")
	checker.NoError(testErr)
	if !mock.failed {
		t.Error("Expected NoError(error) to fail")
	}
	if len(mock.errors) == 0 || !strings.Contains(mock.errors[0], "test error") {
		t.Errorf("Expected error message to contain 'test error', got: %v", mock.errors)
	}
}

func TestHasError(t *testing.T) {
	mock := newMockTestingT()
	checker := check.New(mock)

	// Test passing case
	testErr := errors.New("test error")
	checker.HasError(testErr)
	if mock.failed {
		t.Error("Expected HasError(error) to pass")
	}

	// Test failing case
	mock = newMockTestingT()
	checker = check.New(mock)
	checker.HasError(nil)
	if !mock.failed {
		t.Error("Expected HasError(nil) to fail")
	}
	if len(mock.errors) == 0 || !strings.Contains(mock.errors[0], "Expected an error") {
		t.Errorf("Expected specific error message, got: %v", mock.errors)
	}
}

func TestPanics(t *testing.T) {
	mock := newMockTestingT()
	checker := check.New(mock)

	// Test passing case - function that panics
	checker.Panics(func() {
		panic("test panic")
	})
	if mock.failed {
		t.Error("Expected Panics to pass for panicking function")
	}

	// Test failing case - function that doesn't panic
	mock = newMockTestingT()
	checker = check.New(mock)
	checker.Panics(func() {
		// Do nothing - no panic
	})
	if !mock.failed {
		t.Error("Expected Panics to fail for non-panicking function")
	}
	if len(mock.errors) == 0 || !strings.Contains(mock.errors[0], "Expected panic, but got none") {
		t.Errorf("Expected specific error message, got: %v", mock.errors)
	}
}

func TestNotPanics(t *testing.T) {
	mock := newMockTestingT()
	checker := check.New(mock)

	// Test passing case - function that doesn't panic
	checker.NotPanics(func() {
		// Do nothing - no panic
	})
	if mock.failed {
		t.Error("Expected NotPanics to pass for non-panicking function")
	}

	// Test failing case - function that panics
	mock = newMockTestingT()
	checker = check.New(mock)
	checker.NotPanics(func() {
		panic("test panic")
	})
	if !mock.failed {
		t.Error("Expected NotPanics to fail for panicking function")
	}
	if len(mock.errors) == 0 || !strings.Contains(mock.errors[0], "Expected no panic, but does") {
		t.Errorf("Expected specific error message, got: %v", mock.errors)
	}
}

func TestPanicsWithError(t *testing.T) {
	mock := newMockTestingT()
	checker := check.New(mock)

	// Test with error panic
	checker.Panics(func() {
		panic(errors.New("error panic"))
	})
	if mock.failed {
		t.Error("Expected Panics to pass for error panic")
	}

	// Test with string panic
	mock = newMockTestingT()
	checker = check.New(mock)
	checker.Panics(func() {
		panic("string panic")
	})
	if mock.failed {
		t.Error("Expected Panics to pass for string panic")
	}
}

func TestErrorMessageFormatting(t *testing.T) {
	tests := []struct { //nolint:govet // Don't care about field alignment in tests
		name        string
		msgAndArgs  []any
		expectedMsg string
	}{
		{"no additional message", nil, "Expected hello, got world"},
		{"single string message", []any{"custom message"}, "Expected hello, got world; custom message"},
		{"formatted message", []any{"test %d: %s", 1, "formatted"}, "Expected hello, got world; test 1: formatted"},
		{"multiple args", []any{"arg1", "arg2", 42}, "arg1%!(EXTRA string=arg2, int=42)"},
		{"multiple non-string args", []any{123, "arg2", 42}, "123 arg2 42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.Equal("hello", "world", tt.msgAndArgs...)

			if !mock.failed {
				t.Error("Expected Equal to fail")
			}
			if len(mock.errors) == 0 {
				t.Error("Expected error message")
			} else if !strings.Contains(mock.errors[0], tt.expectedMsg) {
				t.Errorf("Expected error message to contain %q, got %q", tt.expectedMsg, mock.errors[0])
			}
		})
	}
}

// Additional edge case tests
func TestEdgeCases(t *testing.T) {
	// Test with empty slices vs nil slices
	mock := newMockTestingT()
	checker := check.New(mock)

	emptySlice := []int{}
	var nilSlice []int

	checker.Equal(emptySlice, nilSlice)
	if !mock.failed {
		t.Error("Expected Equal to fail for empty slice vs nil slice")
	}

	// Test with different types that have same string representation
	mock = newMockTestingT()
	checker = check.New(mock)

	i := 42
	f := 42.0

	checker.Equal(i, f)
	if !mock.failed {
		t.Error("Expected Equal to fail for int vs float64 with same value")
	}
}

func TestByteSliceHandling(t *testing.T) {
	tests := []struct { //nolint:govet // Don't care about field alignment in tests
		name       string
		a, b       any
		shouldFail bool
	}{
		{"equal byte slices", []byte{1, 2, 3}, []byte{1, 2, 3}, false},
		{"unequal byte slices", []byte{1, 2, 3}, []byte{1, 2, 4}, true},
		{"nil byte slices", []byte(nil), []byte(nil), false},
		{"empty vs nil byte slice", []byte{}, []byte(nil), true},
		{"byte slice vs string", []byte("hello"), "hello", true},
		{"byte slice vs int slice", []byte{1, 2, 3}, []int{1, 2, 3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.Equal(tt.a, tt.b)

			if tt.shouldFail && !mock.failed {
				t.Errorf("Expected Equal to fail for %v vs %v", tt.a, tt.b)
			}
			if !tt.shouldFail && mock.failed {
				t.Errorf("Expected Equal to pass for %v vs %v, but failed with: %v", tt.a, tt.b, mock.errors)
			}
		})
	}
}

func TestNilChecksWithInterfaces(t *testing.T) {
	// Test various nil interface scenarios
	var nilError error
	var nilStringer fmt.Stringer
	var nonNilError error = (*customError)(nil) // typed nil

	mock := newMockTestingT()
	checker := check.New(mock)

	// All these should be considered nil
	checker.Nil(nilError)
	checker.Nil(nilStringer)
	checker.Nil((*int)(nil))
	checker.Nil([]int(nil))
	checker.Nil(map[string]int(nil))

	if mock.failed {
		t.Errorf("Expected all nil checks to pass, but got errors: %v", mock.errors)
	}

	// This typed nil should also be considered nil
	mock = newMockTestingT()
	checker = check.New(mock)
	checker.Nil(nonNilError)
	if mock.failed {
		t.Errorf("Expected typed nil to be considered nil, but got errors: %v", mock.errors)
	}
}

func TestPanicRecovery(t *testing.T) {
	// Test various panic types
	tests := []struct { //nolint:govet // Don't care about field alignment in tests
		name      string
		panicFunc func()
	}{
		{"string panic", func() { panic("test panic") }},
		{"error panic", func() { panic(fmt.Errorf("error panic")) }},
		{"int panic", func() { panic(42) }},
		{"struct panic", func() { panic(struct{ msg string }{"panic"}) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockTestingT()
			checker := check.New(mock)

			checker.Panics(tt.panicFunc)
			if mock.failed {
				t.Errorf("Expected Panics to pass for %s, but got errors: %v", tt.name, mock.errors)
			}
		})
	}
}
