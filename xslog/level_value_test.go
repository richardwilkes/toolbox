// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslog_test

import (
	"flag"
	"log/slog"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xslog"
)

func TestLevelValueSet(t *testing.T) {
	var v xslog.LevelValue
	c := check.New(t)
	c.NoError(v.Set("error"))
	c.Equal(v.Level(), slog.LevelError)
}

func TestLevelValueAddFlags(t *testing.T) {
	var v xslog.LevelValue
	v.AddFlags()
	hasLevel := false
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "log-level" {
			hasLevel = true
		}
	})
	c := check.New(t)
	c.True(hasLevel)
}
