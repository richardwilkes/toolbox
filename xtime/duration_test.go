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
	"math"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xtime"
)

func TestFormatDuration(t *testing.T) {
	c := check.New(t)
	c.Equal("0:00:00.001", xtime.FormatDuration(time.Millisecond, true))
	c.Equal("0:00:01.000", xtime.FormatDuration(1000*time.Millisecond, true))
	c.Equal("0:00:01", xtime.FormatDuration(1001*time.Millisecond, false))
	c.Equal("0:00:01", xtime.FormatDuration(1999*time.Millisecond, false))
	c.Equal("0:01:01", xtime.FormatDuration(61*time.Second, false))
	c.Equal("1:01:00", xtime.FormatDuration(61*time.Minute, false))
	c.Equal("61:00:00", xtime.FormatDuration(61*time.Hour, false))
	c.Equal("0:00:00.000", xtime.FormatDuration(-time.Millisecond, true))
}

func TestParseDuration(t *testing.T) {
	c := check.New(t)

	result, err := xtime.ParseDuration("0:00:00.001")
	c.NoError(err)
	c.Equal(time.Millisecond, result)

	_, err = xtime.ParseDuration("0.001")
	c.HasError(err)

	_, err = xtime.ParseDuration("0:0.001")
	c.HasError(err)

	_, err = xtime.ParseDuration("0:0:.001")
	c.HasError(err)

	result, err = xtime.ParseDuration("0:0:0.001")
	c.NoError(err)
	c.Equal(time.Millisecond, result)

	_, err = xtime.ParseDuration("0:0:-1.001")
	c.HasError(err)

	_, err = xtime.ParseDuration("-1:0:0.001")
	c.HasError(err)

	_, err = xtime.ParseDuration("0:-1:0.001")
	c.HasError(err)

	_, err = xtime.ParseDuration("0:0:0.-001")
	c.HasError(err)

	result, err = xtime.ParseDuration("0:1:61.001")
	c.NoError(err)
	c.Equal(2*time.Minute+time.Second+time.Millisecond, result)

	_, err = xtime.ParseDuration("0:0:0.001.002")
	c.HasError(err)
}

func TestParseDurationOverflow(t *testing.T) {
	c := check.New(t)

	_, err := xtime.ParseDuration("9999999999:00:00")
	c.HasError(err)

	_, err = xtime.ParseDuration("0:999999999999:00")
	c.HasError(err)

	_, err = xtime.ParseDuration("0:00:99999999999999999")
	c.HasError(err)

	// Each component is in range, but their sum is one second past the maximum.
	_, err = xtime.ParseDuration("2562047:47:17")
	c.HasError(err)

	// One second less is the largest whole-second duration.
	result, err := xtime.ParseDuration("2562047:47:16")
	c.NoError(err)
	c.Equal(2562047*time.Hour+47*time.Minute+16*time.Second, result)

	// The largest valid duration must still round-trip.
	result, err = xtime.ParseDuration(xtime.FormatDuration(math.MaxInt64, true))
	c.NoError(err)
	// FormatDuration only emits milliseconds, so the parsed value is truncated to whole milliseconds.
	c.Equal(time.Duration(math.MaxInt64)/time.Millisecond*time.Millisecond, result)
}

func TestParseDurationFractionalSeconds(t *testing.T) {
	c := check.New(t)
	for _, tc := range []struct {
		in   string
		want time.Duration
	}{
		{in: "0:00:00.5", want: 500 * time.Millisecond},
		{in: "0:00:00.05", want: 50 * time.Millisecond},
		{in: "0:00:00.005", want: 5 * time.Millisecond},
		{in: "0:00:00.500", want: 500 * time.Millisecond},
		{in: "0:00:00.050", want: 50 * time.Millisecond},
		{in: "0:00:00.1", want: 100 * time.Millisecond},
		{in: "0:00:00.123", want: 123 * time.Millisecond},
		{in: "0:00:00.000001", want: time.Microsecond},
		{in: "0:00:00.000000001", want: time.Nanosecond},
		{in: "0:00:00.1234567899", want: 123456789 * time.Nanosecond}, // 10th digit truncated
		{in: "0:00:01.5", want: time.Second + 500*time.Millisecond},
	} {
		result, err := xtime.ParseDuration(tc.in)
		c.NoError(err, tc.in)
		c.Equal(tc.want, result, tc.in)
	}

	_, err := xtime.ParseDuration("0:00:00.")
	c.HasError(err)
}

func TestParseDurationRoundTrip(t *testing.T) {
	c := check.New(t)
	for _, d := range []time.Duration{
		0,
		500 * time.Millisecond,
		50 * time.Millisecond,
		time.Second + time.Millisecond,
		2*time.Hour + 30*time.Minute + 45*time.Second + 123*time.Millisecond,
	} {
		formatted := xtime.FormatDuration(d, true)
		result, err := xtime.ParseDuration(formatted)
		c.NoError(err, formatted)
		c.Equal(d, result, formatted)
	}
}

func TestDurationToCode(t *testing.T) {
	c := check.New(t)

	// Zero must yield a valid expression rather than an empty string.
	c.Equal("0", xtime.DurationToCode(0))

	c.Equal("time.Nanosecond", xtime.DurationToCode(1))
	c.Equal("999 * time.Nanosecond", xtime.DurationToCode(999))

	c.Equal("time.Microsecond", xtime.DurationToCode(time.Microsecond))
	c.Equal("500 * time.Microsecond", xtime.DurationToCode(500*time.Microsecond))
	c.Equal("999 * time.Microsecond", xtime.DurationToCode(999*time.Microsecond))

	c.Equal("time.Microsecond + time.Nanosecond", xtime.DurationToCode(time.Microsecond+1))
	c.Equal("time.Microsecond + 500 * time.Nanosecond", xtime.DurationToCode(time.Microsecond+500))

	c.Equal("time.Millisecond", xtime.DurationToCode(time.Millisecond))
	c.Equal("500 * time.Millisecond", xtime.DurationToCode(500*time.Millisecond))
	c.Equal("999 * time.Millisecond", xtime.DurationToCode(999*time.Millisecond))

	c.Equal("time.Millisecond + time.Microsecond", xtime.DurationToCode(time.Millisecond+time.Microsecond))
	c.Equal("time.Millisecond + time.Microsecond + time.Nanosecond", xtime.DurationToCode(time.Millisecond+time.Microsecond+1))
	c.Equal("time.Millisecond + 500 * time.Microsecond + 500 * time.Nanosecond", xtime.DurationToCode(time.Millisecond+500*time.Microsecond+500))

	c.Equal("time.Second", xtime.DurationToCode(time.Second))
	c.Equal("30 * time.Second", xtime.DurationToCode(30*time.Second))
	c.Equal("59 * time.Second", xtime.DurationToCode(59*time.Second))

	c.Equal("time.Second + time.Millisecond", xtime.DurationToCode(time.Second+time.Millisecond))
	c.Equal("time.Second + time.Millisecond + time.Microsecond", xtime.DurationToCode(time.Second+time.Millisecond+time.Microsecond))
	c.Equal("time.Second + time.Millisecond + time.Microsecond + time.Nanosecond", xtime.DurationToCode(time.Second+time.Millisecond+time.Microsecond+1))

	c.Equal("time.Minute", xtime.DurationToCode(time.Minute))
	c.Equal("30 * time.Minute", xtime.DurationToCode(30*time.Minute))
	c.Equal("59 * time.Minute", xtime.DurationToCode(59*time.Minute))

	c.Equal("time.Minute + time.Second", xtime.DurationToCode(time.Minute+time.Second))
	c.Equal("time.Minute + time.Second + time.Millisecond", xtime.DurationToCode(time.Minute+time.Second+time.Millisecond))
	c.Equal("time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond", xtime.DurationToCode(time.Minute+time.Second+time.Millisecond+time.Microsecond+1))

	c.Equal("time.Hour", xtime.DurationToCode(time.Hour))
	c.Equal("12 * time.Hour", xtime.DurationToCode(12*time.Hour))
	c.Equal("24 * time.Hour", xtime.DurationToCode(24*time.Hour))

	c.Equal("time.Hour + time.Minute", xtime.DurationToCode(time.Hour+time.Minute))
	c.Equal("time.Hour + time.Minute + time.Second", xtime.DurationToCode(time.Hour+time.Minute+time.Second))
	c.Equal("time.Hour + time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond", xtime.DurationToCode(time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+1))

	c.Equal("2 * time.Hour + 30 * time.Minute + 45 * time.Second", xtime.DurationToCode(2*time.Hour+30*time.Minute+45*time.Second))
	c.Equal("3 * time.Hour + 15 * time.Minute + 30 * time.Second + 500 * time.Millisecond", xtime.DurationToCode(3*time.Hour+15*time.Minute+30*time.Second+500*time.Millisecond))

	c.Equal("100 * time.Hour", xtime.DurationToCode(100*time.Hour))
	c.Equal("time.Second", xtime.DurationToCode(1000*time.Millisecond))
	c.Equal("time.Millisecond", xtime.DurationToCode(1000*time.Microsecond))

	c.Equal("time.Hour + 30 * time.Minute + 45 * time.Second + 123 * time.Millisecond + 456 * time.Microsecond + 789 * time.Nanosecond", xtime.DurationToCode(time.Hour+30*time.Minute+45*time.Second+123*time.Millisecond+456*time.Microsecond+789))

	c.Equal("-(5 * time.Second)", xtime.DurationToCode(-5*time.Second))
	c.Equal("-(time.Nanosecond)", xtime.DurationToCode(-1))
	c.Equal("-(time.Hour + time.Minute + time.Second)", xtime.DurationToCode(-(time.Hour + time.Minute + time.Second)))

	// The most negative duration can't be negated, so the raw value is emitted instead of recursing.
	c.Equal("time.Duration(-9223372036854775808)", xtime.DurationToCode(math.MinInt64))
}
