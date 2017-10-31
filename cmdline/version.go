package cmdline

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/richardwilkes/gokit/i18n"
)

const (
	secondBits  = 6 // Needs to hold 0-59
	secondMask  = (1 << secondBits) - 1
	secondShift = 0
	minuteBits  = 6 // Needs to hold 0-59
	minuteMask  = (1 << minuteBits) - 1
	minuteShift = secondShift + secondBits
	hourBits    = 5 // Needs to hold 0-23
	hourMask    = (1 << hourBits) - 1
	hourShift   = minuteShift + minuteBits
	dayBits     = 5 // Needs to hold 1-31
	dayMask     = (1 << dayBits) - 1
	dayShift    = hourShift + hourBits
	monthBits   = 4 // Needs to hold 1-12
	monthMask   = (1 << monthBits) - 1
	monthShift  = dayShift + dayBits
	yearBits    = 14 // Needs to hold 0-9999
	yearMask    = (1 << yearBits) - 1
	yearShift   = monthShift + monthBits
	patchBits   = 8 // Holds 0-255
	patchMask   = (1 << patchBits) - 1
	patchShift  = yearShift + yearBits
	minorBits   = 8 // Holds 0-255
	minorMask   = (1 << minorBits) - 1
	minorShift  = patchShift + patchBits
	majorBits   = 8 // Holds 0-255
	majorMask   = (1 << majorBits) - 1
	majorShift  = minorShift + minorBits
)

// Version holds an encoded version number that contains fields for the major,
// minor, and patch release numbers as well as a build date and time.
type Version uint64

// NewVersion constructs a new version. If 'major', 'minor', or 'patch' are
// out of the permitted ranges, a value of 0 will be returned.
func NewVersion(major, minor, patch int, when time.Time) Version {
	if (major&majorMask) != major || (minor&minorMask) != minor || (patch&patchMask) != patch {
		return 0
	}
	when = when.UTC()
	return Version(major)<<majorShift | Version(minor)<<minorShift | Version(patch)<<patchShift |
		Version(when.Year())<<yearShift | Version(when.Month())<<monthShift | Version(when.Day())<<dayShift |
		Version(when.Hour())<<hourShift | Version(when.Minute())<<minuteShift | Version(when.Second())<<secondShift
}

// NewVersionFromString constructs a new version from a string.
func NewVersionFromString(version string) Version {
	var err error
	var major, minor, patch, year, month, day, hour, minute, second int
	parts := strings.Split(version, "-")
	if len(parts) > 1 {
		dateTime := parts[1]
		// date/time in the form YYYYMMDDHHMMSS
		if len(dateTime) == 14 {
			if year, err = strconv.Atoi(dateTime[:4]); err != nil {
				return 0
			}
			if month, err = strconv.Atoi(dateTime[4:6]); err != nil {
				return 0
			}
			if day, err = strconv.Atoi(dateTime[6:8]); err != nil {
				return 0
			}
			if hour, err = strconv.Atoi(dateTime[8:10]); err != nil {
				return 0
			}
			if minute, err = strconv.Atoi(dateTime[10:12]); err != nil {
				return 0
			}
			if second, err = strconv.Atoi(dateTime[12:]); err != nil {
				return 0
			}
		}
	}
	if len(parts) > 0 {
		parts = strings.Split(parts[0], ".")
		count := len(parts)
		if count > 0 {
			if major, err = strconv.Atoi(parts[0]); err != nil {
				return 0
			}
		}
		if count > 1 {
			if minor, err = strconv.Atoi(parts[1]); err != nil {
				return 0
			}
		}
		if count > 2 {
			if patch, err = strconv.Atoi(parts[2]); err != nil {
				return 0
			}
		}
	}
	major = constrain(major, 0, 255)
	minor = constrain(minor, 0, 255)
	patch = constrain(patch, 0, 255)
	year = constrain(year, 0, 9999)
	month = constrain(month, 1, 12)
	day = constrain(day, 1, 31)
	hour = constrain(hour, 0, 23)
	minute = constrain(minute, 0, 59)
	second = constrain(second, 0, 59)
	if major|minor|patch|year|hour|minute|second == 0 && month == 1 && day == 1 {
		return 0
	}
	return Version(major)<<majorShift | Version(minor)<<minorShift | Version(patch)<<patchShift |
		Version(year)<<yearShift | Version(month)<<monthShift | Version(day)<<dayShift |
		Version(hour)<<hourShift | Version(minute)<<minuteShift | Version(second)<<secondShift
}

func constrain(value, min, max int) int {
	if value > max {
		return max
	}
	if value < min {
		return min
	}
	return value
}

// Major returns the major release number associated with this version.
func (v Version) Major() int {
	return int((v >> majorShift) & majorMask)
}

// Minor returns the minor release number associated with this version.
func (v Version) Minor() int {
	return int((v >> minorShift) & minorMask)
}

// Patch returns the patch release number associated with this version.
func (v Version) Patch() int {
	return int((v >> patchShift) & patchMask)
}

// Year returns the build year associated with this version.
func (v Version) Year() int {
	return int((v >> yearShift) & yearMask)
}

// Month returns the build month associated with this version.
func (v Version) Month() int {
	return int((v >> monthShift) & monthMask)
}

// Day returns the build day associated with this version.
func (v Version) Day() int {
	return int((v >> dayShift) & dayMask)
}

// Hour returns the build hour associated with this version.
func (v Version) Hour() int {
	return int((v >> hourShift) & hourMask)
}

// Minute returns the build minute associated with this version.
func (v Version) Minute() int {
	return int((v >> minuteShift) & minuteMask)
}

// Second returns the build second associated with this version.
func (v Version) Second() int {
	return int((v >> secondShift) & secondMask)
}

// When returns the build date and time in UTC associated with this version.
func (v Version) When() time.Time {
	return time.Date(v.Year(), time.Month(v.Month()), v.Day(), v.Hour(), v.Minute(), v.Second(), 0, time.UTC)
}

// IsDevelopment returns true if this version has not been initialized.
func (v Version) IsDevelopment() bool {
	return v == 0
}

// IsWhenUnset returns true if the build time has not been initialized.
func (v Version) IsWhenUnset() bool {
	return v.Year() == 0 && v.Month() == 1 && v.Day() == 1 && v.Hour() == 0 && v.Minute() == 0 && v.Second() == 0
}

// String implements the fmt.Stringer interface.
func (v Version) String() string {
	return v.Format(true, false)
}

// Format the version. If 'includeVersionWord' is true, then the word
// "Version" will be included as needed. If 'includeBuildDateTime' is true,
// then the build date and time will be included.
func (v Version) Format(includeVersionWord, includeBuildDateTime bool) string {
	if v.IsDevelopment() {
		if includeVersionWord {
			return i18n.Text("Development Version")
		}
		return i18n.Text("Development")
	}
	buffer := &bytes.Buffer{}
	if includeVersionWord {
		buffer.WriteString(i18n.Text("Version "))
	}
	fmt.Fprintf(buffer, "%d.%d", v.Major(), v.Minor())
	if v.Patch() != 0 {
		fmt.Fprintf(buffer, ".%d", v.Patch())
	}
	if includeBuildDateTime && !v.IsWhenUnset() {
		fmt.Fprintf(buffer, "-%04d%02d%02d%02d%02d%02d", v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second())
	}
	return buffer.String()
}
