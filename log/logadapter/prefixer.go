// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package logadapter

import "fmt"

// Prefixer adds a prefix to another logger's output.
//
// Deprecated: Use slog instead. August 28, 2023
type Prefixer struct {
	Logger Logger
	Prefix string
}

// Debug logs a debug message. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Debug(v ...any) {
	p.Logger.Debugf("%s%s", p.Prefix, fmt.Sprint(v...))
}

// Debugf logs a debug message. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Debugf(format string, v ...any) {
	p.Logger.Debugf("%s%s", p.Prefix, fmt.Sprintf(format, v...))
}

// Info logs an informational message. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Info(v ...any) {
	p.Logger.Infof("%s%s", p.Prefix, fmt.Sprint(v...))
}

// Infof logs an informational message. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Infof(format string, v ...any) {
	p.Logger.Infof("%s%s", p.Prefix, fmt.Sprintf(format, v...))
}

// Warn logs a warning message. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Warn(v ...any) {
	p.Logger.Warnf("%s%s", p.Prefix, fmt.Sprint(v...))
}

// Warnf logs a warning message. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Warnf(format string, v ...any) {
	p.Logger.Warnf("%s%s", p.Prefix, fmt.Sprintf(format, v...))
}

// Error logs an error message. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Error(v ...any) {
	p.Logger.Errorf("%s%s", p.Prefix, fmt.Sprint(v...))
}

// Errorf logs an error message. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Errorf(format string, v ...any) {
	p.Logger.Errorf("%s%s", p.Prefix, fmt.Sprintf(format, v...))
}

// Fatal logs a fatal error message. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Fatal(status int, v ...any) {
	p.Logger.Fatalf(status, "%s%s", p.Prefix, fmt.Sprint(v...))
}

// Fatalf logs a fatal error message. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Fatalf(status int, format string, v ...any) {
	p.Logger.Fatalf(status, "%s%s", p.Prefix, fmt.Sprintf(format, v...))
}

// Time starts timing an event and logs an informational message. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Time(v ...any) Timing {
	return p.Logger.Timef("%s%s", p.Prefix, fmt.Sprint(v...))
}

// Timef starts timing an event and logs an informational message. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: Use slog instead. August 28, 2023
func (p *Prefixer) Timef(format string, v ...any) Timing {
	return p.Logger.Timef("%s%s", p.Prefix, fmt.Sprintf(format, v...))
}
