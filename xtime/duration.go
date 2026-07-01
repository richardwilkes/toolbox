// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xtime

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// ParseDuration parses the duration string, as produced by FormatDuration(). The fractional seconds component, if
// present, is interpreted as a decimal fraction of a second (so ".5" is 500ms and ".05" is 50ms), accepting more or
// fewer than the three digits FormatDuration emits; digits beyond nanosecond resolution are truncated.
func ParseDuration(duration string) (time.Duration, error) {
	parts := strings.Split(strings.TrimSpace(duration), ":")
	if len(parts) != 3 {
		return 0, errs.New("Invalid format")
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil || hours < 0 {
		return 0, errs.New("Invalid hour format")
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil || minutes < 0 {
		return 0, errs.New("Invalid minute format")
	}
	parts = strings.Split(parts[2], ".")
	var seconds int
	var nanos int
	switch len(parts) {
	case 2:
		if nanos, err = parseFractionalSeconds(parts[1]); err != nil {
			return 0, err
		}
		fallthrough
	case 1:
		if seconds, err = strconv.Atoi(parts[0]); err != nil || seconds < 0 {
			return 0, errs.New("Invalid second format")
		}
	default:
		return 0, errs.New("Invalid second format: too many decimal points")
	}
	total := time.Duration(0)
	var ok bool
	for _, component := range []struct {
		count int
		unit  time.Duration
	}{
		{count: hours, unit: time.Hour},
		{count: minutes, unit: time.Minute},
		{count: seconds, unit: time.Second},
		{count: nanos, unit: time.Nanosecond},
	} {
		if total, ok = addDurationComponent(total, component.count, component.unit); !ok {
			return 0, errs.New("Duration out of range")
		}
	}
	return total, nil
}

// addDurationComponent adds count*unit to sum, reporting whether the multiplication or addition overflowed. Both count
// and sum are assumed to be non-negative, as guaranteed by ParseDuration's component validation.
func addDurationComponent(sum time.Duration, count int, unit time.Duration) (time.Duration, bool) {
	product := time.Duration(count) * unit
	if product/unit != time.Duration(count) {
		return 0, false
	}
	total := sum + product
	if total < sum {
		return 0, false
	}
	return total, true
}

// parseFractionalSeconds interprets the digits following the decimal point of the seconds field as a decimal fraction
// of a second, returning the equivalent number of nanoseconds. The digit positions matter: "5" yields 500ms and "05"
// yields 50ms. Only digits are accepted (no sign) and any beyond nanosecond resolution (more than 9) are truncated.
func parseFractionalSeconds(frac string) (int, error) {
	if frac == "" {
		return 0, errs.New("Invalid fractional second format")
	}
	for _, r := range frac {
		if r < '0' || r > '9' {
			return 0, errs.New("Invalid fractional second format")
		}
	}
	const nanoDigits = 9
	if len(frac) > nanoDigits {
		frac = frac[:nanoDigits]
	} else {
		frac += strings.Repeat("0", nanoDigits-len(frac))
	}
	nanos, err := strconv.Atoi(frac)
	if err != nil {
		return 0, errs.New("Invalid fractional second format")
	}
	return nanos, nil
}

// FormatDuration formats the duration as either "0:00:00" or "0:00:00.000". Negative durations are treated as zero.
func FormatDuration(duration time.Duration, includeMillis bool) string {
	if duration < 0 {
		duration = 0
	}
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute
	duration -= minutes * time.Minute
	seconds := duration / time.Second
	duration -= seconds * time.Second
	if includeMillis {
		return fmt.Sprintf("%d:%02d:%02d.%03d", hours, minutes, seconds, duration/time.Millisecond)
	}
	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
}

// DurationToCode turns a time.Duration into more human-readable text required for code than a simple number of
// nanoseconds. The result is always a valid Go expression: zero yields "0" and negative durations are wrapped in a
// leading negation, e.g. "-(5 * time.Second)".
func DurationToCode(duration time.Duration) string {
	switch {
	case duration == 0:
		return "0"
	case duration == math.MinInt64:
		// Negating math.MinInt64 overflows back to itself, which would recurse forever, so emit the raw value. It
		// remains a valid Go expression.
		return "time.Duration(" + strconv.FormatInt(int64(duration), 10) + ")"
	case duration < 0:
		return "-(" + DurationToCode(-duration) + ")"
	}
	var buffer strings.Builder
	duration = durationToCodePart(&buffer, duration, time.Hour, "Hour")
	duration = durationToCodePart(&buffer, duration, time.Minute, "Minute")
	duration = durationToCodePart(&buffer, duration, time.Second, "Second")
	duration = durationToCodePart(&buffer, duration, time.Millisecond, "Millisecond")
	duration = durationToCodePart(&buffer, duration, time.Microsecond, "Microsecond")
	durationToCodePart(&buffer, duration, time.Nanosecond, "Nanosecond")
	return buffer.String()
}

func durationToCodePart(buffer *strings.Builder, duration, unit time.Duration, unitName string) time.Duration {
	if duration < unit {
		return duration
	}
	if buffer.Len() > 0 {
		buffer.WriteString(" + ")
	}
	value := duration / unit
	if value != 1 {
		fmt.Fprintf(buffer, "%d * ", value)
	}
	buffer.WriteString("time.")
	buffer.WriteString(unitName)
	return duration - (duration/unit)*unit
}
