// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/txt"
)

func TestFormatDuration(t *testing.T) {
	c := check.New(t)
	for i, one := range []struct {
		Expected      string
		Duration      time.Duration
		IncludeMillis bool
	}{
		{"0:00:00.001", time.Millisecond, true},
		{"0:00:01.000", 1000 * time.Millisecond, true},
		{"0:00:01", 1001 * time.Millisecond, false},
		{"0:00:01", 1999 * time.Millisecond, false},
		{"0:01:01", 61 * time.Second, false},
		{"1:01:00", 61 * time.Minute, false},
		{"61:00:00", 61 * time.Hour, false},
	} {
		c.Equal(one.Expected, txt.FormatDuration(one.Duration, one.IncludeMillis), "Index %d", i)
	}
}

func TestParseDuration(t *testing.T) {
	c := check.New(t)
	for i, one := range []struct {
		Input            string
		ExpectedDuration time.Duration
		ExpectErr        bool
	}{
		{"0:00:00.001", time.Millisecond, false},
		{"0.001", 0, true},
		{"0:0.001", 0, true},
		{"0:0:.001", 0, true},
		{"0:0:0.001", time.Millisecond, false},
		{"0:0:-1.001", 0, true},
		{"-1:0:0.001", 0, true},
		{"0:-1:0.001", 0, true},
		{"0:0:0.-001", 0, true},
		{"0:1:61.001", 2*time.Minute + time.Second + time.Millisecond, false},
	} {
		result, err := txt.ParseDuration(one.Input)
		desc := fmt.Sprintf("Index %d: %s", i, one.Input)
		if one.ExpectErr {
			c.HasError(err, desc)
		} else {
			c.NoError(err, desc)
			c.Equal(one.ExpectedDuration, result, desc)
		}
	}
}
