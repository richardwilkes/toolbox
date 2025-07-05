// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/v2/xreflect"
)

// Checker provides some simple helpers for testing.
type Checker struct {
	*testing.T
}

// New creates a new Checker from the testing.T.
func New(t *testing.T) Checker {
	return Checker{T: t}
}

// Equal compares two values for equality.
func (c Checker) Equal(expected, actual any, msgAndArgs ...any) {
	c.Helper()
	if !c.equal(expected, actual) {
		c.errMsg(fmt.Sprintf("Expected %v, got %v", expected, actual), msgAndArgs...)
	}
}

// NotEqual compares two values for inequality.
func (c Checker) NotEqual(expected, actual any, msgAndArgs ...any) {
	c.Helper()
	if c.equal(expected, actual) {
		c.errMsg(fmt.Sprintf("Expected %v to not be %v", expected, actual), msgAndArgs...)
	}
}

func (c Checker) equal(expected, actual any) bool {
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
func (c Checker) Nil(value any, msgAndArgs ...any) {
	c.Helper()
	if !xreflect.IsNil(value) {
		c.errMsg(fmt.Sprintf("Expected nil, instead got %v", value), msgAndArgs...)
	}
}

// NotNil expects value to not be nil.
func (c Checker) NotNil(value any, msgAndArgs ...any) {
	c.Helper()
	if xreflect.IsNil(value) {
		c.errMsg("Expected a non-nil value", msgAndArgs...)
	}
}

// True expects value to be true.
func (c Checker) True(value bool, msgAndArgs ...any) {
	c.Helper()
	if !value {
		c.errMsg("Expected true", msgAndArgs...)
	}
}

// False expects value to be false.
func (c Checker) False(value bool, msgAndArgs ...any) {
	c.Helper()
	if value {
		c.errMsg("Expected false", msgAndArgs...)
	}
}

// Contains expects s to contain substr.
func (c Checker) Contains(s, substr string, msgAndArgs ...any) {
	c.Helper()
	if !strings.Contains(s, substr) {
		c.errMsg(fmt.Sprintf("Expected string %q to contain %q", s, substr), msgAndArgs...)
	}
}

// NotContains expects s not to contain substr.
func (c Checker) NotContains(s, substr string, msgAndArgs ...any) {
	c.Helper()
	if strings.Contains(s, substr) {
		c.errMsg(fmt.Sprintf("Expected string %q not to contain %q", s, substr), msgAndArgs...)
	}
}

// HasPrefix expects s to have the prefix substr.
func (c Checker) HasPrefix(s, prefix string, msgAndArgs ...any) {
	c.Helper()
	if !strings.HasPrefix(s, prefix) {
		c.errMsg(fmt.Sprintf("Expected string %q to have prefix %q", s, prefix), msgAndArgs...)
	}
}

// NoPrefix expects s not to have the prefix substr.
func (c Checker) NoPrefix(s, prefix string, msgAndArgs ...any) {
	c.Helper()
	if strings.HasPrefix(s, prefix) {
		c.errMsg(fmt.Sprintf("Expected string %q not to have prefix %q", s, prefix), msgAndArgs...)
	}
}

// NoError expects err to be nil.
func (c Checker) NoError(err error, msgAndArgs ...any) {
	c.Helper()
	if err != nil {
		c.errMsg(fmt.Sprintf("Expected no error, got %v", err), msgAndArgs...)
	}
}

// HasError expects err to not be nil.
func (c Checker) HasError(err error, msgAndArgs ...any) {
	c.Helper()
	if err == nil {
		c.errMsg("Expected an error", msgAndArgs...)
	}
}

// Panics expects f to panic when called.
func (c Checker) Panics(f func(), msgAndArgs ...any) {
	c.Helper()
	if err := c.doesPanic(f); err == nil {
		c.errMsg("Expected panic, but got none", msgAndArgs...)
	}
}

// NotPanics expects no panic when f is called.
func (c Checker) NotPanics(f func(), msgAndArgs ...any) {
	c.Helper()
	if err := c.doesPanic(f); err != nil {
		c.errMsg(fmt.Sprintf("Expected no panic, but does: %v", err), msgAndArgs...)
	}
}

func (c Checker) doesPanic(f func()) (panicErr error) {
	// Can't use xos.PanicRecovery because that would cause an import cycle in some places
	defer func() {
		if recovered := recover(); recovered != nil {
			err, ok := recovered.(error)
			if !ok {
				err = fmt.Errorf("%+v", recovered)
			}
			panicErr = err
		}
	}()
	f()
	return
}

func (c Checker) errMsg(prefix string, msgAndArgs ...any) {
	c.Helper()
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
	c.Error(buffer.String())
}
