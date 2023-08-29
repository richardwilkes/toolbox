// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package jot provides simple asynchronous logging.
package jot

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/xio/term"
)

// Log levels
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var logChannel = make(chan *record, 100)

// Level holds a log level.
//
// Deprecated: Use slog instead. August 28, 2023
type Level int

type record struct {
	when        time.Time
	level       Level
	msg         string
	writer      io.Writer
	response    chan bool
	setMinLevel bool
}

func init() {
	atexit.Register(Flush)
	// Initialization of `out` was moved outside of the goroutine in case jot isn't actually being used (but is being
	// initialized as a side-effect of bringing in some other package) to prevent races with other threads that might
	// be trying to use os.Stderr.
	out := term.NewANSI(os.Stderr)
	go func() {
		levelNames := []string{"DBG", "INF", "WRN", "ERR", "FTL"}
		levelColors := []term.Color{term.Blue, 0, term.Yellow, term.Red, term.Red}
		levelStyles := []term.Style{term.Bold, 0, term.Bold, term.Bold, term.Bold | term.Blink}
		minLevel := DEBUG
		for rec := range logChannel {
			switch {
			case rec.writer != nil:
				out = term.NewANSI(rec.writer)
			case rec.response != nil:
				rec.response <- true
			case rec.setMinLevel:
				minLevel = rec.level
			case rec.level >= minLevel:
				color := levelColors[rec.level]
				if color != 0 {
					out.Foreground(color, levelStyles[rec.level])
				}
				write(out, levelNames[rec.level])
				if color != 0 {
					out.Reset()
				}
				timeDate := rec.when.Format(" | 2006-01-02 | 15:04:05.000 | ")
				write(out, timeDate)
				parts := strings.Split(rec.msg, "\n")
				write(out, parts[0])
				if len(parts) > 1 {
					prefix := "\n" + strings.Repeat(" ", len(levelNames[0])+len(timeDate))
					for i := 1; i < len(parts); i++ {
						write(out, prefix)
						write(out, parts[i])
					}
				}
				write(out, "\n")
			}
		}
	}()
}

func write(out io.Writer, text string) {
	// The extra code here is just to quiet the linter about not checking for an error.
	if _, err := out.Write([]byte(text)); err != nil {
		return
	}
}

// SetWriter sets the io.Writer to use when writing log messages. Default is os.Stderr.
//
// Deprecated: Use slog instead. August 28, 2023
func SetWriter(w io.Writer) {
	logChannel <- &record{writer: w}
}

// SetMinimumLevel sets the minimum log level that will be output. Default is DEBUG.
//
// Deprecated: Use slog instead. August 28, 2023
func SetMinimumLevel(level Level) {
	logChannel <- &record{
		level:       level,
		setMinLevel: true,
	}
}

// Debug logs a debugging message. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func Debug(v ...any) {
	logChannel <- &record{
		when:  time.Now(),
		level: DEBUG,
		msg:   fmt.Sprint(v...),
	}
}

// Debugf logs a debugging message. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func Debugf(format string, v ...any) {
	logChannel <- &record{
		when:  time.Now(),
		level: DEBUG,
		msg:   fmt.Sprintf(format, v...),
	}
}

// Info logs an informational message. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func Info(v ...any) {
	logChannel <- &record{
		when:  time.Now(),
		level: INFO,
		msg:   fmt.Sprint(v...),
	}
}

// Infof logs an informational message. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func Infof(format string, v ...any) {
	logChannel <- &record{
		when:  time.Now(),
		level: INFO,
		msg:   fmt.Sprintf(format, v...),
	}
}

// Warn logs a warning message. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func Warn(v ...any) {
	logChannel <- &record{
		when:  time.Now(),
		level: WARN,
		msg:   fmt.Sprint(v...),
	}
}

// Warnf logs a warning message. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func Warnf(format string, v ...any) {
	logChannel <- &record{
		when:  time.Now(),
		level: WARN,
		msg:   fmt.Sprintf(format, v...),
	}
}

// Error logs an error message. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func Error(v ...any) {
	logChannel <- &record{
		when:  time.Now(),
		level: ERROR,
		msg:   fmt.Sprint(v...),
	}
}

// Errorf logs an error message. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func Errorf(format string, v ...any) {
	logChannel <- &record{
		when:  time.Now(),
		level: ERROR,
		msg:   fmt.Sprintf(format, v...),
	}
}

// Fatal logs a fatal error message. Arguments other than the status are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func Fatal(status int, v ...any) {
	logChannel <- &record{
		when:  time.Now(),
		level: FATAL,
		msg:   fmt.Sprint(v...),
	}
	atexit.Exit(status)
}

// Fatalf logs a fatal error message. Arguments other than the status are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func Fatalf(status int, format string, v ...any) {
	logChannel <- &record{
		when:  time.Now(),
		level: FATAL,
		msg:   fmt.Sprintf(format, v...),
	}
	atexit.Exit(status)
}

// FatalIfErr calls 'Fatal(1, err)' if 'err' is not nil.
//
// Deprecated: Use slog instead. August 28, 2023
func FatalIfErr(err error) {
	if err != nil {
		Fatal(1, err)
	}
}

// Flush waits for all current log entries to be written before returning.
//
// Deprecated: Use slog instead. August 28, 2023
func Flush() {
	rec := &record{response: make(chan bool)}
	logChannel <- rec
	<-rec.response
}
