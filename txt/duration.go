package txt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/richardwilkes/gokit/errs"
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
