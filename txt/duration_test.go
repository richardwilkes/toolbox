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
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/txt"
)

func TestFormatDuration(t *testing.T) {
	c := check.New(t)

	// Test 1 millisecond with milliseconds included
	c.Equal("0:00:00.001", txt.FormatDuration(time.Millisecond, true))

	// Test 1 second (1000ms) with milliseconds included
	c.Equal("0:00:01.000", txt.FormatDuration(1000*time.Millisecond, true))

	// Test 1001 milliseconds without milliseconds (rounds down)
	c.Equal("0:00:01", txt.FormatDuration(1001*time.Millisecond, false))

	// Test 1999 milliseconds without milliseconds (rounds down)
	c.Equal("0:00:01", txt.FormatDuration(1999*time.Millisecond, false))

	// Test 61 seconds without milliseconds
	c.Equal("0:01:01", txt.FormatDuration(61*time.Second, false))

	// Test 61 minutes without milliseconds
	c.Equal("1:01:00", txt.FormatDuration(61*time.Minute, false))

	// Test 61 hours without milliseconds
	c.Equal("61:00:00", txt.FormatDuration(61*time.Hour, false))
}

func TestParseDuration(t *testing.T) {
	c := check.New(t)

	// Test valid millisecond parsing
	result, err := txt.ParseDuration("0:00:00.001")
	c.NoError(err)
	c.Equal(time.Millisecond, result)

	// Test invalid format - missing components
	_, err = txt.ParseDuration("0.001")
	c.HasError(err)

	// Test invalid format - missing seconds
	_, err = txt.ParseDuration("0:0.001")
	c.HasError(err)

	// Test invalid format - missing seconds with colon
	_, err = txt.ParseDuration("0:0:.001")
	c.HasError(err)

	// Test valid format with milliseconds
	result, err = txt.ParseDuration("0:0:0.001")
	c.NoError(err)
	c.Equal(time.Millisecond, result)

	// Test invalid format - negative seconds
	_, err = txt.ParseDuration("0:0:-1.001")
	c.HasError(err)

	// Test invalid format - negative hours
	_, err = txt.ParseDuration("-1:0:0.001")
	c.HasError(err)

	// Test invalid format - negative minutes
	_, err = txt.ParseDuration("0:-1:0.001")
	c.HasError(err)

	// Test invalid format - negative milliseconds
	_, err = txt.ParseDuration("0:0:0.-001")
	c.HasError(err)

	// Test overflow handling - seconds > 59 should carry over
	result, err = txt.ParseDuration("0:1:61.001")
	c.NoError(err)
	c.Equal(2*time.Minute+time.Second+time.Millisecond, result)
}

func TestDurationToCode(t *testing.T) {
	c := check.New(t)

	// Test zero duration
	c.Equal("", txt.DurationToCode(0))

	// Test single nanoseconds
	c.Equal("time.Nanosecond", txt.DurationToCode(1))
	c.Equal("999 * time.Nanosecond", txt.DurationToCode(999))

	// Test single microseconds
	c.Equal("time.Microsecond", txt.DurationToCode(time.Microsecond))
	c.Equal("500 * time.Microsecond", txt.DurationToCode(500*time.Microsecond))
	c.Equal("999 * time.Microsecond", txt.DurationToCode(999*time.Microsecond))

	// Test microseconds with nanoseconds
	c.Equal("time.Microsecond + time.Nanosecond", txt.DurationToCode(time.Microsecond+1))
	c.Equal("time.Microsecond + 500 * time.Nanosecond", txt.DurationToCode(time.Microsecond+500))

	// Test single milliseconds
	c.Equal("time.Millisecond", txt.DurationToCode(time.Millisecond))
	c.Equal("500 * time.Millisecond", txt.DurationToCode(500*time.Millisecond))
	c.Equal("999 * time.Millisecond", txt.DurationToCode(999*time.Millisecond))

	// Test milliseconds with microseconds and nanoseconds
	c.Equal("time.Millisecond + time.Microsecond", txt.DurationToCode(time.Millisecond+time.Microsecond))
	c.Equal("time.Millisecond + time.Microsecond + time.Nanosecond", txt.DurationToCode(time.Millisecond+time.Microsecond+1))
	c.Equal("time.Millisecond + 500 * time.Microsecond + 500 * time.Nanosecond", txt.DurationToCode(time.Millisecond+500*time.Microsecond+500))

	// Test single seconds
	c.Equal("time.Second", txt.DurationToCode(time.Second))
	c.Equal("30 * time.Second", txt.DurationToCode(30*time.Second))
	c.Equal("59 * time.Second", txt.DurationToCode(59*time.Second))

	// Test seconds with smaller units
	c.Equal("time.Second + time.Millisecond", txt.DurationToCode(time.Second+time.Millisecond))
	c.Equal("time.Second + time.Millisecond + time.Microsecond", txt.DurationToCode(time.Second+time.Millisecond+time.Microsecond))
	c.Equal("time.Second + time.Millisecond + time.Microsecond + time.Nanosecond", txt.DurationToCode(time.Second+time.Millisecond+time.Microsecond+1))

	// Test single minutes
	c.Equal("time.Minute", txt.DurationToCode(time.Minute))
	c.Equal("30 * time.Minute", txt.DurationToCode(30*time.Minute))
	c.Equal("59 * time.Minute", txt.DurationToCode(59*time.Minute))

	// Test minutes with smaller units
	c.Equal("time.Minute + time.Second", txt.DurationToCode(time.Minute+time.Second))
	c.Equal("time.Minute + time.Second + time.Millisecond", txt.DurationToCode(time.Minute+time.Second+time.Millisecond))
	c.Equal("time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond", txt.DurationToCode(time.Minute+time.Second+time.Millisecond+time.Microsecond+1))

	// Test single hours
	c.Equal("time.Hour", txt.DurationToCode(time.Hour))
	c.Equal("12 * time.Hour", txt.DurationToCode(12*time.Hour))
	c.Equal("24 * time.Hour", txt.DurationToCode(24*time.Hour))

	// Test hours with smaller units
	c.Equal("time.Hour + time.Minute", txt.DurationToCode(time.Hour+time.Minute))
	c.Equal("time.Hour + time.Minute + time.Second", txt.DurationToCode(time.Hour+time.Minute+time.Second))
	c.Equal("time.Hour + time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond", txt.DurationToCode(time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+1))

	// Test complex durations
	c.Equal("2 * time.Hour + 30 * time.Minute + 45 * time.Second", txt.DurationToCode(2*time.Hour+30*time.Minute+45*time.Second))
	c.Equal("3 * time.Hour + 15 * time.Minute + 30 * time.Second + 500 * time.Millisecond", txt.DurationToCode(3*time.Hour+15*time.Minute+30*time.Second+500*time.Millisecond))

	// Test edge cases with larger values
	c.Equal("100 * time.Hour", txt.DurationToCode(100*time.Hour))
	c.Equal("time.Second", txt.DurationToCode(1000*time.Millisecond))
	c.Equal("time.Millisecond", txt.DurationToCode(1000*time.Microsecond))

	// Test fractional combinations
	c.Equal("time.Hour + 30 * time.Minute + 45 * time.Second + 123 * time.Millisecond + 456 * time.Microsecond + 789 * time.Nanosecond", txt.DurationToCode(time.Hour+30*time.Minute+45*time.Second+123*time.Millisecond+456*time.Microsecond+789))
}
