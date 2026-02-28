// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xcrypto_test

import (
	bytes "bytes"
	crypto_rand "crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xcrypto"
)

func TestEncryptDecryptStreamWithKeyPair(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	publicKey := &privateKey.PublicKey
	plaintext := []byte("The quick brown fox jumps over the lazy dog.")
	in := bytes.NewReader(plaintext)
	var encrypted bytes.Buffer
	c.NoError(xcrypto.EncryptStreamWithPublicKey(in, &encrypted, publicKey))
	c.True(encrypted.Len() > len(plaintext))
	var decrypted bytes.Buffer
	c.NoError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(encrypted.Bytes()), &decrypted, privateKey))
	c.Equal(plaintext, decrypted.Bytes())
}
