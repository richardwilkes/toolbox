// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// ParseDuration parses the duration string, as produced by FormatDuration().
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
	var millis int
	switch len(parts) {
	case 2:
		if millis, err = strconv.Atoi(parts[1]); err != nil || millis < 0 {
			return 0, errs.New("Invalid millisecond format")
		}
		fallthrough
	case 1:
		if seconds, err = strconv.Atoi(parts[0]); err != nil || seconds < 0 {
			return 0, errs.New("Invalid second format")
		}
	default:
		return 0, errs.New("Invalid second format: too many decimal points")
	}
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(millis)*time.Millisecond, nil
}

// FormatDuration formats the duration as either "0:00:00" or "0:00:00.000".
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
// nanoseconds.
func DurationToCode(duration time.Duration) string {
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
