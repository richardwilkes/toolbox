// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/txt"
	"github.com/stretchr/testify/assert"
)

func TestDigitToValue(t *testing.T) {
	checkDigitToValue('5', 5, t)
	checkDigitToValue('٥', 5, t)
	checkDigitToValue('𑁯', 9, t)
	_, err := txt.DigitToValue('a')
	assert.Error(t, err)
}

func checkDigitToValue(ch rune, expected int, t *testing.T) {
	value, err := txt.DigitToValue(ch)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, value)
}
