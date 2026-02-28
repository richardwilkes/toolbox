// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package tid

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xos"
)

// TID is a unique identifier. These are similar to v4 UUIDs, but are shorter and have a different format that includes
// a kind byte as the first character. TIDs are 17 characters long, are URL safe, and contain 96 bits of entropy.
type TID string

// KindAlphabet is the set of characters that can be used as the first character of a TID. The kind has no intrinsic
// meaning, but can be used to differentiate between different types of ids.
const KindAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// MustNewTID creates a new TID with a random value and the specified kind. If an error occurs, this function panics.
func MustNewTID(kind byte) TID {
	return xos.Must(NewTID(kind))
}

// NewTID creates a new TID with a random value and the specified kind.
func NewTID(kind byte) (TID, error) {
	if strings.IndexByte(KindAlphabet, kind) == -1 {
		return "", errs.New("invalid kind")
	}
	var buffer [12]byte
	if _, err := rand.Read(buffer[:]); err != nil {
		return "", errs.Wrap(err)
	}
	return TID(fmt.Sprintf("%c%s", kind, base64.RawURLEncoding.EncodeToString(buffer[:]))), nil
}

// FromString converts a string to a TID.
func FromString(id string) (TID, error) {
	tid := TID(id)
	if !IsValid(tid) {
		return "", errs.New("invalid TID")
	}
	return tid, nil
}

// FromStringOfKind converts a string to a TID and verifies that it has the specified kind.
func FromStringOfKind(id string, kind byte) (TID, error) {
	tid := TID(id)
	if !IsKindAndValid(tid, kind) {
		return "", errs.New("invalid TID")
	}
	return tid, nil
}

// IsValid returns true if the TID is a valid TID.
func IsValid(id TID) bool {
	if len(id) != 17 || strings.IndexByte(KindAlphabet, id[0]) == -1 {
		return false
	}
	_, err := base64.RawURLEncoding.DecodeString(string(id[1:]))
	return err == nil
}

// IsKind returns true if the TID has the specified kind.
func IsKind(id TID, kind byte) bool {
	return len(id) == 17 && id[0] == kind && strings.IndexByte(KindAlphabet, kind) != -1
}

// IsKindAndValid returns true if the TID is a valid TID with the specified kind.
func IsKindAndValid(id TID, kind byte) bool {
	if !IsKind(id, kind) {
		return false
	}
	_, err := base64.RawURLEncoding.DecodeString(string(id[1:]))
	return err == nil
}
