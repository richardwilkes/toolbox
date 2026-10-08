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

// TID is a unique identifier similar to a v4 UUID, but shorter: a kind character followed by 96 bits of entropy, for
// 17 URL-safe characters in all.
type TID string

// KindAlphabet is the set of characters permitted as the kind (first character) of a TID. The kind has no intrinsic
// meaning, but can be used to distinguish different types of ids.
const KindAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// bodyAlphabet is the unpadded URL-safe base64 alphabet, which the 16 characters after the kind must come from.
// It is checked directly rather than by decoding, since the base64 decoder silently skips CR and LF.
const bodyAlphabet = KindAlphabet + "-_"

// MustNewTID is like NewTID, but panics on error.
func MustNewTID(kind byte) TID {
	return xos.Must(NewTID(kind))
}

// NewTID creates a new TID with a random value and the specified kind, which must be in KindAlphabet.
func NewTID(kind byte) (TID, error) {
	if strings.IndexByte(KindAlphabet, kind) == -1 {
		return "", errs.Newf("invalid kind %q", kind)
	}
	var buffer [12]byte
	if _, err := rand.Read(buffer[:]); err != nil {
		return "", errs.Wrap(err)
	}
	return TID(fmt.Sprintf("%c%s", kind, base64.RawURLEncoding.EncodeToString(buffer[:]))), nil
}

// FromString converts a string to a TID, returning an error if it is not valid.
func FromString(id string) (TID, error) {
	tid := TID(id)
	if !IsValid(tid) {
		return "", errs.Newf("invalid TID %q", id)
	}
	return tid, nil
}

// FromStringOfKind converts a string to a TID, returning an error if it is not valid or not of the specified kind.
func FromStringOfKind(id string, kind byte) (TID, error) {
	tid, err := FromString(id)
	if err != nil {
		return "", err
	}
	if !IsKind(tid, kind) {
		return "", errs.Newf("TID %q is not of kind %q", id, kind)
	}
	return tid, nil
}

// IsValid returns true if id is a well-formed TID.
func IsValid(id TID) bool {
	return len(id) == 17 && strings.IndexByte(KindAlphabet, id[0]) != -1 && hasValidBody(id)
}

// IsKind returns true if the TID has the specified kind.
func IsKind(id TID, kind byte) bool {
	return len(id) == 17 && id[0] == kind && strings.IndexByte(KindAlphabet, kind) != -1
}

// IsKindAndValid returns true if the TID is a valid TID with the specified kind.
func IsKindAndValid(id TID, kind byte) bool {
	return IsKind(id, kind) && hasValidBody(id)
}

// hasValidBody returns true if every character after the kind is in bodyAlphabet. The caller must have already
// checked the length.
func hasValidBody(id TID) bool {
	for i := 1; i < len(id); i++ {
		if strings.IndexByte(bodyAlphabet, id[i]) == -1 {
			return false
		}
	}
	return true
}
