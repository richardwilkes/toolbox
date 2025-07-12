// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package tid_test

import (
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/tid"
)

func TestNewTID(t *testing.T) {
	c := check.New(t)

	// Test with valid kinds
	validKinds := []byte{'A', 'Z', 'a', 'z', '0', '9', 'M', 'x', '5'}
	for _, kind := range validKinds {
		result, err := tid.NewTID(kind)
		c.NoError(err)
		c.Equal(17, len(string(result)), "TID should be 17 characters long")
		c.Equal(kind, string(result)[0], "First character should match the kind")
		c.True(tid.IsValid(result), "Generated TID should be valid")
		c.True(tid.IsKind(result, kind), "TID should have the correct kind")
	}
}

func TestNewTIDInvalidKind(t *testing.T) {
	c := check.New(t)

	// Test with invalid kinds
	invalidKinds := []byte{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '+', '=', ' ', '\t', '\n'}
	for _, kind := range invalidKinds {
		_, err := tid.NewTID(kind)
		c.HasError(err, "Should return error for invalid kind: %c", kind)
	}
}

func TestMustNewTID(t *testing.T) {
	c := check.New(t)

	// Test with valid kind
	result := tid.MustNewTID('A')
	c.Equal(17, len(string(result)), "TID should be 17 characters long")
	c.Equal(byte('A'), string(result)[0], "First character should be 'A'")
	c.True(tid.IsValid(result), "Generated TID should be valid")
}

func TestMustNewTIDPanicsOnInvalidKind(t *testing.T) {
	c := check.New(t)

	c.Panics(func() {
		tid.MustNewTID('!')
	}, "Should panic with invalid kind")
}

func TestFromString(t *testing.T) {
	c := check.New(t)

	// Test with valid TID string
	originalTID := tid.MustNewTID('T')
	result, err := tid.FromString(string(originalTID))
	c.NoError(err)
	c.Equal(originalTID, result)
}

func TestFromStringInvalid(t *testing.T) {
	c := check.New(t)

	invalidTIDs := []string{
		"",                         // empty string
		"short",                    // too short
		"!InvalidKindChar1234",     // invalid kind character
		"Atoolongstring1234567890", // too long
		"A123456789012345",         // wrong length (16 chars)
		"A12345678901234567",       // wrong length (18 chars)
		"A123456789012345!",        // invalid base64 character
		"A@#$%^&*()_+=-[]{}",       // invalid base64 characters
	}

	for _, invalidTID := range invalidTIDs {
		_, err := tid.FromString(invalidTID)
		c.HasError(err, "Should return error for invalid TID: %q", invalidTID)
	}
}

func TestFromStringOfKind(t *testing.T) {
	c := check.New(t)

	// Test with valid TID string of correct kind
	originalTID := tid.MustNewTID('K')
	result, err := tid.FromStringOfKind(string(originalTID), 'K')
	c.NoError(err)
	c.Equal(originalTID, result)
}

func TestFromStringOfKindWrongKind(t *testing.T) {
	c := check.New(t)

	// Test with valid TID string but wrong kind
	originalTID := tid.MustNewTID('A')
	_, err := tid.FromStringOfKind(string(originalTID), 'B')
	c.HasError(err, "Should return error when kind doesn't match")
}

func TestFromStringOfKindInvalidTID(t *testing.T) {
	c := check.New(t)

	// Test with invalid TID string
	_, err := tid.FromStringOfKind("invalid", 'A')
	c.HasError(err, "Should return error for invalid TID")
}

func TestIsValid(t *testing.T) {
	c := check.New(t)

	// Test with valid TIDs
	validTID := tid.MustNewTID('V')
	c.True(tid.IsValid(validTID), "Generated TID should be valid")

	// Test with invalid TIDs
	invalidTIDs := []tid.TID{
		"",                        // empty
		"short",                   // too short
		"!ValidBase64String12",    // invalid kind
		"Vtoolongstring123456789", // too long
		"V123456789012345",        // wrong length (16 chars)
		"V12345678901234567",      // wrong length (18 chars)
		"V123456789012345!",       // invalid base64
	}

	for _, invalidTID := range invalidTIDs {
		c.False(tid.IsValid(invalidTID), "Should be invalid: %q", invalidTID)
	}
}

func TestIsKind(t *testing.T) {
	c := check.New(t)

	// Test with correct kind
	testTID := tid.MustNewTID('K')
	c.True(tid.IsKind(testTID, 'K'), "Should return true for correct kind")

	// Test with wrong kind
	c.False(tid.IsKind(testTID, 'X'), "Should return false for wrong kind")

	// Test with invalid kind character
	c.False(tid.IsKind(testTID, '!'), "Should return false for invalid kind character")

	// Test with wrong length TIDs
	c.False(tid.IsKind("short", 'K'), "Should return false for short string")
	c.False(tid.IsKind("toolongstring1234567890", 'K'), "Should return false for long string")
}

func TestIsKindAndValid(t *testing.T) {
	c := check.New(t)

	// Test with valid TID of correct kind
	testTID := tid.MustNewTID('V')
	c.True(tid.IsKindAndValid(testTID, 'V'), "Should return true for valid TID of correct kind")

	// Test with valid TID of wrong kind
	c.False(tid.IsKindAndValid(testTID, 'W'), "Should return false for valid TID of wrong kind")

	// Test with invalid TID of correct kind
	invalidTID := tid.TID("V123456789012345!")
	c.False(tid.IsKindAndValid(invalidTID, 'V'), "Should return false for invalid TID even with correct kind")

	// Test with invalid TID of wrong kind
	c.False(tid.IsKindAndValid(invalidTID, 'W'), "Should return false for invalid TID of wrong kind")
}

func TestKindAlphabet(t *testing.T) {
	c := check.New(t)

	// Test that KindAlphabet contains expected characters
	expectedChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	c.Equal(expectedChars, tid.KindAlphabet)
	c.Equal(62, len(tid.KindAlphabet), "KindAlphabet should have 62 characters")

	// Test that all characters in KindAlphabet work as valid kinds
	for i := 0; i < len(tid.KindAlphabet); i++ {
		kind := tid.KindAlphabet[i]
		testTID, err := tid.NewTID(kind)
		c.NoError(err, "All KindAlphabet characters should be valid kinds: %c", kind)
		c.True(tid.IsKind(testTID, kind), "Generated TID should have correct kind: %c", kind)
	}
}

func TestTIDUniqueness(t *testing.T) {
	c := check.New(t)

	// Generate multiple TIDs and ensure they're unique
	const count = 1000
	tids := make(map[string]bool, count)

	for range count {
		testTID := tid.MustNewTID('U')
		tidStr := string(testTID)
		c.False(tids[tidStr], "TID should be unique: %s", tidStr)
		tids[tidStr] = true
	}
}

func TestTIDFormat(t *testing.T) {
	c := check.New(t)

	testTID := tid.MustNewTID('F')
	tidStr := string(testTID)

	// Test length
	c.Equal(17, len(tidStr), "TID should be exactly 17 characters")

	// Test first character is the kind
	c.Equal(byte('F'), tidStr[0], "First character should be the kind")

	// Test remaining characters are valid base64url
	base64Part := tidStr[1:]
	c.Equal(16, len(base64Part), "Base64 part should be 16 characters")

	// Test that it doesn't contain invalid base64url characters
	invalidChars := "+/="
	for _, char := range invalidChars {
		c.False(strings.ContainsRune(base64Part, char), "Base64 part should not contain character: %c", char)
	}
}

func TestTIDStringConversion(t *testing.T) {
	c := check.New(t)

	// Test that TID can be converted to string and back
	originalTID := tid.MustNewTID('S')
	tidStr := string(originalTID)

	convertedTID, err := tid.FromString(tidStr)
	c.NoError(err)
	c.Equal(originalTID, convertedTID)
}
