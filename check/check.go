// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package check

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/errs"
)

// Equal compares two values for equality.
func Equal(t *testing.T, expected, actual any, msgAndArgs ...any) {
	t.Helper()
	if !equal(expected, actual) {
		errMsg(t, fmt.Sprintf("Expected %v, got %v", expected, actual), msgAndArgs...)
	}
}

// NotEqual compares two values for inequality.
func NotEqual(t *testing.T, expected, actual any, msgAndArgs ...any) {
	t.Helper()
	if equal(expected, actual) {
		errMsg(t, fmt.Sprintf("Expected %v to not be %v", expected, actual), msgAndArgs...)
	}
}

func equal(expected, actual any) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}
	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}
	act, ok2 := actual.([]byte)
	if !ok2 {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

// Nil expects value to be nil.
func Nil(t *testing.T, value any, msgAndArgs ...any) {
	t.Helper()
	if !toolbox.IsNil(value) {
		errMsg(t, fmt.Sprintf("Expected nil, instead got %v", value), msgAndArgs...)
	}
}

// NotNil expects value to not be nil.
func NotNil(t *testing.T, value any, msgAndArgs ...any) {
	t.Helper()
	if toolbox.IsNil(value) {
		errMsg(t, "Expected a non-nil value", msgAndArgs...)
	}
}

// True expects value to be true.
func True(t *testing.T, value bool, msgAndArgs ...any) {
	t.Helper()
	if !value {
		errMsg(t, "Expected true", msgAndArgs...)
	}
}

// False expects value to be false.
func False(t *testing.T, value bool, msgAndArgs ...any) {
	t.Helper()
	if value {
		errMsg(t, "Expected false", msgAndArgs...)
	}
}

// Contains expects s to contain substr.
func Contains(t *testing.T, s, substr string, msgAndArgs ...any) {
	t.Helper()
	if !strings.Contains(s, substr) {
		errMsg(t, fmt.Sprintf("Expected string %q to contain %q", s, substr), msgAndArgs...)
	}
}

// NotContains expects s not to contain substr.
func NotContains(t *testing.T, s, substr string, msgAndArgs ...any) {
	t.Helper()
	if strings.Contains(s, substr) {
		errMsg(t, fmt.Sprintf("Expected string %q not to contain %q", s, substr), msgAndArgs...)
	}
}

// NoError expects err to be nil.
func NoError(t *testing.T, err error, msgAndArgs ...any) {
	t.Helper()
	if err != nil {
		errMsg(t, fmt.Sprintf("Expected no error, got %v", err), msgAndArgs...)
	}
}

// Error expects err to not be nil.
func Error(t *testing.T, err error, msgAndArgs ...any) {
	t.Helper()
	if err == nil {
		errMsg(t, "Expected an error", msgAndArgs...)
	}
}

// Panics expects f to panic when called.
func Panics(t *testing.T, f func(), msgAndArgs ...any) {
	t.Helper()
	if err := doesPanic(f); err == nil {
		errMsg(t, "Expected panic, but got none", msgAndArgs...)
	}
}

// NotPanics expects no panic when f is called.
func NotPanics(t *testing.T, f func(), msgAndArgs ...any) {
	t.Helper()
	if err := doesPanic(f); err != nil {
		errMsg(t, fmt.Sprintf("Expected no panic, but does: %v", err), msgAndArgs...)
	}
}

func doesPanic(f func()) (panicErr error) {
	defer errs.Recovery(func(err error) { panicErr = err })
	f()
	return
}

func errMsg(t *testing.T, prefix string, msgAndArgs ...any) {
	t.Helper()
	var buffer strings.Builder
	buffer.WriteString(prefix)
	switch len(msgAndArgs) {
	case 0:
	case 1:
		_, _ = fmt.Fprintf(&buffer, "; %v", msgAndArgs[0])
	default:
		if s, ok := msgAndArgs[0].(string); ok {
			_, _ = fmt.Fprintf(&buffer, "; "+s, msgAndArgs[1:]...)
		} else {
			buffer.WriteByte(';')
			for _, one := range msgAndArgs {
				_, _ = fmt.Fprintf(&buffer, " %v", one)
			}
		}
	}
	t.Error(buffer.String())
}
