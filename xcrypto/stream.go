// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"

	"github.com/richardwilkes/toolbox/errs"
)

// EncryptStreamWithPublicKey copies 'in' to 'out', encrypting the bytes along the way. Note that the output stream will
// be larger than the input stream by aes.BlockSize + publicKey.Size() bytes.
func EncryptStreamWithPublicKey(in io.Reader, out io.Writer, publicKey *rsa.PublicKey) error {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return errs.Wrap(err)
	}
	encryptionKey := make([]byte, 32) // aes256
	if _, err := io.ReadFull(rand.Reader, encryptionKey); err != nil {
		return errs.Wrap(err)
	}
	encryptedEncryptionKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, encryptionKey, nil)
	if err != nil {
		return errs.Wrap(err)
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return errs.Wrap(err)
	}
	if _, err = out.Write(encryptedEncryptionKey); err != nil {
		return errs.Wrap(err)
	}
	if _, err = out.Write(iv); err != nil {
		return errs.Wrap(err)
	}
	if _, err = io.Copy(&cipher.StreamWriter{
		S: cipher.NewCTR(block, iv),
		W: out,
	}, in); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

// DecryptStreamWithPrivateKey copies 'in' to 'out', decrypting the bytes along the way. Note that the output stream
// will be smaller than the input stream by aes.BlockSize + publicKey.Size() bytes.
func DecryptStreamWithPrivateKey(in io.Reader, out io.Writer, privateKey *rsa.PrivateKey) error {
	encryptedEncryptionKey := make([]byte, privateKey.Size())
	if _, err := in.Read(encryptedEncryptionKey); err != nil {
		return errs.Wrap(err)
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := in.Read(iv); err != nil {
		return errs.Wrap(err)
	}
	encryptionKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedEncryptionKey, nil)
	if err != nil {
		return errs.Wrap(err)
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return errs.Wrap(err)
	}
	if _, err = io.Copy(out, &cipher.StreamReader{
		S: cipher.NewCTR(block, iv),
		R: in,
	}); err != nil {
		return errs.Wrap(err)
	}
	return nil
}
