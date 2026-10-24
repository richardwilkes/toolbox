// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package tid_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/tid"
)

func TestGeneratedTIDsAreWellFormed(t *testing.T) {
	c := check.New(t)

	c.Equal("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", tid.KindAlphabet)
	for i := range len(tid.KindAlphabet) {
		kind := tid.KindAlphabet[i]
		result, err := tid.NewTID(kind)
		c.NoError(err, "NewTID should accept kind %c", kind)
		for _, id := range []tid.TID{result, tid.MustNewTID(kind)} {
			c.Equal(17, len(string(id)), "TID should be 17 characters long: %q", id)
			c.Equal(kind, string(id)[0], "First character should match the kind: %q", id)
			c.True(tid.IsValid(id), "Generated TID should be valid: %q", id)
			c.True(tid.IsKind(id, kind), "TID should have the correct kind: %q", id)
			c.False(strings.ContainsAny(string(id)[1:], "+/="), "Body should not contain '+', '/' or '=': %q", id)
		}
	}
}

func TestGenerationRejectsInvalidKind(t *testing.T) {
	c := check.New(t)

	invalidKinds := []byte{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '+', '=', ' ', '\t', '\n'}
	for _, kind := range invalidKinds {
		_, err := tid.NewTID(kind)
		c.HasError(err, "Should return error for invalid kind: %c", kind)
		c.Panics(func() { tid.MustNewTID(kind) }, "MustNewTID should panic for invalid kind: %c", kind)
	}
}

func TestFromStringRoundTrips(t *testing.T) {
	c := check.New(t)

	originalTID := tid.MustNewTID('K')
	result, err := tid.FromStringOfKind(string(originalTID), 'K')
	c.NoError(err)
	c.Equal(originalTID, result)

	result, err = tid.FromString(string(originalTID))
	c.NoError(err)
	c.Equal(originalTID, result)
}

func TestIsKind(t *testing.T) {
	c := check.New(t)

	testTID := tid.MustNewTID('K')
	c.True(tid.IsKind(testTID, 'K'), "Should return true for correct kind")

	c.False(tid.IsKind(testTID, 'X'), "Should return false for wrong kind")

	c.False(tid.IsKind(tid.TID("!"+string(testTID)[1:]), '!'), "Should return false for invalid kind character")

	c.False(tid.IsKind("short", 'K'), "Should return false for short string")
	c.False(tid.IsKind("toolongstring1234567890", 'K'), "Should return false for long string")
}

func TestIsKindAndValid(t *testing.T) {
	c := check.New(t)

	testTID := tid.MustNewTID('V')
	c.True(tid.IsKindAndValid(testTID, 'V'), "Should return true for valid TID of correct kind")

	c.False(tid.IsKindAndValid(testTID, 'W'), "Should return false for valid TID of wrong kind")

	invalidTID := tid.TID("V123456789012345!")
	c.False(tid.IsKindAndValid(invalidTID, 'V'), "Should return false for invalid TID even with correct kind")

	c.False(tid.IsKindAndValid(invalidTID, 'W'), "Should return false for invalid TID of wrong kind")
}

func TestTIDUniqueness(t *testing.T) {
	c := check.New(t)

	const count = 1000
	tids := make(map[string]bool, count)

	for range count {
		testTID := tid.MustNewTID('U')
		tidStr := string(testTID)
		c.False(tids[tidStr], "TID should be unique: %s", tidStr)
		tids[tidStr] = true
	}
}

func TestRejectsMalformedIDs(t *testing.T) {
	c := check.New(t)

	body := string(tid.MustNewTID('A'))[1:]
	ids := []string{
		"A" + strings.Repeat("\n", 16),
		"A" + strings.Repeat("\r", 16),
		"A" + body[:15] + "\n",
		"A" + body[:14] + "\r\n",
		"A\n" + body[:15],
		"",
		"short",
		"!InvalidKindChar1",        // invalid kind
		"Atoolongstring1234567890", // too long
		"A123456789012345",         // 16 characters
		"A12345678901234567",       // 18 characters
		"A123456789012345!",        // invalid body character
	}
	for _, id := range ids {
		c.False(tid.IsValid(tid.TID(id)), "IsValid should reject %q", id)
		c.False(tid.IsKindAndValid(tid.TID(id), 'A'), "IsKindAndValid should reject %q", id)
		_, err := tid.FromString(id)
		c.HasError(err, "FromString should reject %q", id)
		_, err = tid.FromStringOfKind(id, 'A')
		c.HasError(err, "FromStringOfKind should reject %q", id)
	}
}

func TestErrorsNameRejectedInput(t *testing.T) {
	c := check.New(t)

	_, err := tid.FromString("bogus")
	c.HasError(err)
	c.Contains(err.Error(), `"bogus"`, "FromString error should include the rejected string")

	_, err = tid.NewTID('!')
	c.HasError(err)
	c.Contains(err.Error(), `'!'`, "NewTID error should include the rejected kind")

	_, err = tid.FromStringOfKind("bogus", 'A')
	c.HasError(err)
	c.Contains(err.Error(), `"bogus"`, "FromStringOfKind error should include the rejected string")

	id := string(tid.MustNewTID('A'))
	_, err = tid.FromStringOfKind(id, 'B')
	c.HasError(err)
	c.Contains(err.Error(), `'B'`, "wrong-kind error should name the expected kind")
	c.Contains(err.Error(), fmt.Sprintf("%q", id), "wrong-kind error should include the id")
	c.NotContains(err.Error(), "invalid TID", "a well-formed id of the wrong kind should not be reported as invalid")
}
