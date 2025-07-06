// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package i18n provides internationalization support for applications.
package i18n

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xfilepath"
	"github.com/richardwilkes/toolbox/v2/xio"
)

const (
	// Extension is the file name extension required on localization files.
	Extension = ".i18n"
)

var (
	// Dir is the directory to scan for localization files. This will occur only once, the first time a call to Text()
	// is made. If you do not set this prior to the first call, a directory in the same location as the executable with
	// "_i18n" appended to the executable name (sans any extension) will be used.
	Dir string
	// Language is the language that should be used for text returned from calls to Text(). It is initialized to the
	// result of calling Locale(). You may set this at runtime, forcing a particular language for all subsequent calls
	// to Text().
	Language = Locale()
	// Languages is a slice of languages to fall back to should the one specified in the Language variable not be
	// available. It is initialized to the value of the LANGUAGE environment variable.
	Languages    = strings.Split(os.Getenv("LANGUAGE"), ":")
	altLocalizer atomic.Pointer[localizer]
	once         sync.Once
	langMap      = make(map[string]map[string]string)
	hierLock     sync.Mutex
	hierMap      = make(map[string][]string)
)

type localizer struct {
	Text func(string) string
}

// SetLocalizer sets the function to use for localizing text. If this is not set or explicitly set to nil, the default
// localization mechanism will be used.
func SetLocalizer(f func(string) string) {
	var trampoline *localizer
	if f != nil {
		trampoline = &localizer{Text: f}
	}
	altLocalizer.Store(trampoline)
}

// Text returns a localized version of the text if one exists, or the original text if not.
func Text(text string) string {
	if f := altLocalizer.Load(); f != nil {
		return f.Text(text)
	}
	once.Do(func() {
		if Dir == "" {
			path, err := os.Executable()
			if err != nil {
				return
			}
			path, err = filepath.EvalSymlinks(path)
			if err != nil {
				return
			}
			path, err = filepath.Abs(xfilepath.TrimExtension(path) + "_i18n")
			if err != nil {
				return
			}
			Dir = path
		}
		dirEntry, err := os.ReadDir(Dir)
		if err != nil {
			return
		}
		for _, one := range dirEntry {
			if one.IsDir() {
				continue
			}
			if name := one.Name(); filepath.Ext(name) == Extension {
				load(name)
			}
		}
	})

	var result string
	if result = lookup(text, Language); result != "" {
		return result
	}
	for _, language := range Languages {
		if result = lookup(text, language); result != "" {
			return result
		}
	}
	return text
}

func lookup(text, language string) string {
	for _, lang := range hierarchy(language) {
		if translations := langMap[lang]; translations != nil {
			if str, ok := translations[text]; ok {
				return str
			}
		}
	}
	return ""
}

func hierarchy(language string) []string {
	lang := strings.ToLower(language)
	hierLock.Lock()
	defer hierLock.Unlock()
	if s, ok := hierMap[lang]; ok {
		return s
	}
	one := strings.ReplaceAll(strings.ReplaceAll(lang, "-", "_"), ".", "_")
	var s []string
	for {
		s = append(s, one)
		if i := strings.LastIndex(one, "_"); i != -1 {
			one = one[:i]
		} else {
			break
		}
	}
	hierMap[lang] = s
	return s
}

func load(name string) {
	path := filepath.Join(Dir, name)
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		errs.Log(errs.NewWithCause("i18n: unable to load", err), "path", path)
		return
	}
	defer xio.CloseIgnoringErrors(f)
	lineNum := 1
	lastKeyLineStart := 1
	translations := make(map[string]string)
	var key, value string
	var hasKey, hasValue bool
	s := bufio.NewScanner(f)
	const lineLogKey = "line"
	const fileLogKey = "file"
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "k:") {
			if hasValue {
				if _, exists := translations[key]; !exists {
					translations[key] = value
				} else {
					slog.Warn("i18n: ignoring duplicate key", lineLogKey, lastKeyLineStart, fileLogKey, path)
				}
				hasKey = false
				hasValue = false
			}
			var buffer string
			if _, err = fmt.Sscanf(line, "k:%q", &buffer); err != nil {
				slog.Warn("i18n: ignoring invalid key", lineLogKey, lineNum, fileLogKey, path)
			} else {
				if hasKey {
					key += "\n" + buffer
				} else {
					key = buffer
					hasKey = true
					lastKeyLineStart = lineNum
				}
			}
		} else if strings.HasPrefix(line, "v:") {
			if hasKey {
				var buffer string
				if _, err = fmt.Sscanf(line, "v:%q", &buffer); err != nil {
					slog.Warn("i18n: ignoring invalid value", lineLogKey, lineNum, fileLogKey, path)
				} else {
					if hasValue {
						value += "\n" + buffer
					} else {
						value = buffer
						hasValue = true
					}
				}
			} else {
				slog.Warn("i18n: ignoring value with no previous key", lineLogKey, lineNum, fileLogKey, path)
			}
		}
		lineNum++
	}
	if hasKey {
		if hasValue {
			if _, exists := translations[key]; !exists {
				translations[key] = value
			} else {
				slog.Warn("i18n: ignoring duplicate key", lineLogKey, lastKeyLineStart, fileLogKey, path)
			}
		} else {
			slog.Warn("i18n: ignoring key with missing value", lineLogKey, lastKeyLineStart, fileLogKey, path)
		}
	}
	key = strings.ToLower(name[:len(name)-len(Extension)])
	key = strings.ReplaceAll(strings.ReplaceAll(key, "-", "_"), ".", "_")
	langMap[key] = translations
}
