// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xtime_test

import (
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xtime"
)

func TestFormatDuration(t *testing.T) {
	c := check.New(t)

	// Test 1 millisecond with milliseconds included
	c.Equal("0:00:00.001", xtime.FormatDuration(time.Millisecond, true))

	// Test 1 second (1000ms) with milliseconds included
	c.Equal("0:00:01.000", xtime.FormatDuration(1000*time.Millisecond, true))

	// Test 1001 milliseconds without milliseconds (rounds down)
	c.Equal("0:00:01", xtime.FormatDuration(1001*time.Millisecond, false))

	// Test 1999 milliseconds without milliseconds (rounds down)
	c.Equal("0:00:01", xtime.FormatDuration(1999*time.Millisecond, false))

	// Test 61 seconds without milliseconds
	c.Equal("0:01:01", xtime.FormatDuration(61*time.Second, false))

	// Test 61 minutes without milliseconds
	c.Equal("1:01:00", xtime.FormatDuration(61*time.Minute, false))

	// Test 61 hours without milliseconds
	c.Equal("61:00:00", xtime.FormatDuration(61*time.Hour, false))

	// Test negative duration being treated as zero
	c.Equal("0:00:00.000", xtime.FormatDuration(-time.Millisecond, true))
}

func TestParseDuration(t *testing.T) {
	c := check.New(t)

	// Test valid millisecond parsing
	result, err := xtime.ParseDuration("0:00:00.001")
	c.NoError(err)
	c.Equal(time.Millisecond, result)

	// Test invalid format - missing components
	_, err = xtime.ParseDuration("0.001")
	c.HasError(err)

	// Test invalid format - missing seconds
	_, err = xtime.ParseDuration("0:0.001")
	c.HasError(err)

	// Test invalid format - missing seconds with colon
	_, err = xtime.ParseDuration("0:0:.001")
	c.HasError(err)

	// Test valid format with milliseconds
	result, err = xtime.ParseDuration("0:0:0.001")
	c.NoError(err)
	c.Equal(time.Millisecond, result)

	// Test invalid format - negative seconds
	_, err = xtime.ParseDuration("0:0:-1.001")
	c.HasError(err)

	// Test invalid format - negative hours
	_, err = xtime.ParseDuration("-1:0:0.001")
	c.HasError(err)

	// Test invalid format - negative minutes
	_, err = xtime.ParseDuration("0:-1:0.001")
	c.HasError(err)

	// Test invalid format - negative milliseconds
	_, err = xtime.ParseDuration("0:0:0.-001")
	c.HasError(err)

	// Test overflow handling - seconds > 59 should carry over
	result, err = xtime.ParseDuration("0:1:61.001")
	c.NoError(err)
	c.Equal(2*time.Minute+time.Second+time.Millisecond, result)

	// Test too many decimal points in seconds
	_, err = xtime.ParseDuration("0:0:0.001.002")
	c.HasError(err)
}

func TestDurationToCode(t *testing.T) {
	c := check.New(t)

	// Test zero duration
	c.Equal("", xtime.DurationToCode(0))

	// Test single nanoseconds
	c.Equal("time.Nanosecond", xtime.DurationToCode(1))
	c.Equal("999 * time.Nanosecond", xtime.DurationToCode(999))

	// Test single microseconds
	c.Equal("time.Microsecond", xtime.DurationToCode(time.Microsecond))
	c.Equal("500 * time.Microsecond", xtime.DurationToCode(500*time.Microsecond))
	c.Equal("999 * time.Microsecond", xtime.DurationToCode(999*time.Microsecond))

	// Test microseconds with nanoseconds
	c.Equal("time.Microsecond + time.Nanosecond", xtime.DurationToCode(time.Microsecond+1))
	c.Equal("time.Microsecond + 500 * time.Nanosecond", xtime.DurationToCode(time.Microsecond+500))

	// Test single milliseconds
	c.Equal("time.Millisecond", xtime.DurationToCode(time.Millisecond))
	c.Equal("500 * time.Millisecond", xtime.DurationToCode(500*time.Millisecond))
	c.Equal("999 * time.Millisecond", xtime.DurationToCode(999*time.Millisecond))

	// Test milliseconds with microseconds and nanoseconds
	c.Equal("time.Millisecond + time.Microsecond", xtime.DurationToCode(time.Millisecond+time.Microsecond))
	c.Equal("time.Millisecond + time.Microsecond + time.Nanosecond", xtime.DurationToCode(time.Millisecond+time.Microsecond+1))
	c.Equal("time.Millisecond + 500 * time.Microsecond + 500 * time.Nanosecond", xtime.DurationToCode(time.Millisecond+500*time.Microsecond+500))

	// Test single seconds
	c.Equal("time.Second", xtime.DurationToCode(time.Second))
	c.Equal("30 * time.Second", xtime.DurationToCode(30*time.Second))
	c.Equal("59 * time.Second", xtime.DurationToCode(59*time.Second))

	// Test seconds with smaller units
	c.Equal("time.Second + time.Millisecond", xtime.DurationToCode(time.Second+time.Millisecond))
	c.Equal("time.Second + time.Millisecond + time.Microsecond", xtime.DurationToCode(time.Second+time.Millisecond+time.Microsecond))
	c.Equal("time.Second + time.Millisecond + time.Microsecond + time.Nanosecond", xtime.DurationToCode(time.Second+time.Millisecond+time.Microsecond+1))

	// Test single minutes
	c.Equal("time.Minute", xtime.DurationToCode(time.Minute))
	c.Equal("30 * time.Minute", xtime.DurationToCode(30*time.Minute))
	c.Equal("59 * time.Minute", xtime.DurationToCode(59*time.Minute))

	// Test minutes with smaller units
	c.Equal("time.Minute + time.Second", xtime.DurationToCode(time.Minute+time.Second))
	c.Equal("time.Minute + time.Second + time.Millisecond", xtime.DurationToCode(time.Minute+time.Second+time.Millisecond))
	c.Equal("time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond", xtime.DurationToCode(time.Minute+time.Second+time.Millisecond+time.Microsecond+1))

	// Test single hours
	c.Equal("time.Hour", xtime.DurationToCode(time.Hour))
	c.Equal("12 * time.Hour", xtime.DurationToCode(12*time.Hour))
	c.Equal("24 * time.Hour", xtime.DurationToCode(24*time.Hour))

	// Test hours with smaller units
	c.Equal("time.Hour + time.Minute", xtime.DurationToCode(time.Hour+time.Minute))
	c.Equal("time.Hour + time.Minute + time.Second", xtime.DurationToCode(time.Hour+time.Minute+time.Second))
	c.Equal("time.Hour + time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond", xtime.DurationToCode(time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+1))

	// Test complex durations
	c.Equal("2 * time.Hour + 30 * time.Minute + 45 * time.Second", xtime.DurationToCode(2*time.Hour+30*time.Minute+45*time.Second))
	c.Equal("3 * time.Hour + 15 * time.Minute + 30 * time.Second + 500 * time.Millisecond", xtime.DurationToCode(3*time.Hour+15*time.Minute+30*time.Second+500*time.Millisecond))

	// Test edge cases with larger values
	c.Equal("100 * time.Hour", xtime.DurationToCode(100*time.Hour))
	c.Equal("time.Second", xtime.DurationToCode(1000*time.Millisecond))
	c.Equal("time.Millisecond", xtime.DurationToCode(1000*time.Microsecond))

	// Test fractional combinations
	c.Equal("time.Hour + 30 * time.Minute + 45 * time.Second + 123 * time.Millisecond + 456 * time.Microsecond + 789 * time.Nanosecond", xtime.DurationToCode(time.Hour+30*time.Minute+45*time.Second+123*time.Millisecond+456*time.Microsecond+789))
}
