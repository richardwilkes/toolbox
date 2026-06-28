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
	"testing/iotest"

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

// TestEncryptDecryptStreamVariousSizes exercises the chunking boundaries, including empty input, sub-chunk input, and
// inputs that span exact and partial multiples of the internal chunk size.
func TestEncryptDecryptStreamVariousSizes(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	publicKey := &privateKey.PublicKey
	const chunkSize = 64 * 1024
	for _, size := range []int{0, 1, 100, chunkSize - 1, chunkSize, chunkSize + 1, 2 * chunkSize, 3*chunkSize + 17} {
		plaintext := make([]byte, size)
		_, err = crypto_rand.Read(plaintext)
		c.NoError(err)
		var encrypted bytes.Buffer
		c.NoError(xcrypto.EncryptStreamWithPublicKey(bytes.NewReader(plaintext), &encrypted, publicKey))
		var decrypted bytes.Buffer
		c.NoError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(encrypted.Bytes()), &decrypted, privateKey))
		c.True(bytes.Equal(plaintext, decrypted.Bytes()), "size %d round-trip mismatch", size)
	}
}

// TestDecryptStreamWithPartialReads ensures decryption works when the input stream returns data in short reads, as can
// happen with network connections, pipes, and files. The encrypted key and IV must be read with io.ReadFull rather than
// a single Read call.
func TestDecryptStreamWithPartialReads(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	publicKey := &privateKey.PublicKey
	plaintext := []byte("The quick brown fox jumps over the lazy dog.")
	var encrypted bytes.Buffer
	c.NoError(xcrypto.EncryptStreamWithPublicKey(bytes.NewReader(plaintext), &encrypted, publicKey))
	// iotest.OneByteReader forces every Read to return at most a single byte, exercising the short-read path.
	in := iotest.OneByteReader(bytes.NewReader(encrypted.Bytes()))
	var decrypted bytes.Buffer
	c.NoError(xcrypto.DecryptStreamWithPrivateKey(in, &decrypted, privateKey))
	c.Equal(plaintext, decrypted.Bytes())
}

// TestDecryptStreamDetectsBitFlip verifies that flipping a single bit anywhere in the ciphertext body causes
// decryption to fail rather than silently returning attacker-controlled plaintext.
func TestDecryptStreamDetectsBitFlip(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	publicKey := &privateKey.PublicKey
	plaintext := []byte("The quick brown fox jumps over the lazy dog.")
	// Flip a bit at each position past the RSA-sealed key + nonce prefix header and confirm every one is rejected.
	headerSize := privateKey.Size() + 7
	for offset := headerSize; ; offset++ {
		var encrypted bytes.Buffer
		c.NoError(xcrypto.EncryptStreamWithPublicKey(bytes.NewReader(plaintext), &encrypted, publicKey))
		tampered := encrypted.Bytes()
		if offset >= len(tampered) {
			break
		}
		tampered[offset] ^= 0x01
		var decrypted bytes.Buffer
		c.HasError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(tampered), &decrypted, privateKey))
	}
}

// TestDecryptStreamDetectsTruncation verifies that removing trailing bytes (a truncation attack) is detected rather
// than being accepted as a shorter-but-valid stream.
func TestDecryptStreamDetectsTruncation(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	publicKey := &privateKey.PublicKey
	// Use more than one chunk so dropping the final chunk leaves an otherwise well-formed prefix.
	plaintext := make([]byte, 2*64*1024+512)
	_, err = crypto_rand.Read(plaintext)
	c.NoError(err)
	var encrypted bytes.Buffer
	c.NoError(xcrypto.EncryptStreamWithPublicKey(bytes.NewReader(plaintext), &encrypted, publicKey))
	full := encrypted.Bytes()
	// Drop the trailing chunk (plaintext chunk + 16-byte tag) entirely.
	truncated := full[:len(full)-(64*1024+16)]
	var decrypted bytes.Buffer
	c.HasError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(truncated), &decrypted, privateKey))
	// Removing just the final tag byte must also fail.
	var decrypted2 bytes.Buffer
	c.HasError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(full[:len(full)-1]), &decrypted2, privateKey))
}

// TestDecryptStreamDetectsExtraData verifies that appending an extra forged chunk after a legitimate final chunk is
// rejected, since the real final chunk is bound as final by its authentication tag.
func TestDecryptStreamDetectsExtraData(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	publicKey := &privateKey.PublicKey
	plaintext := []byte("The quick brown fox jumps over the lazy dog.")
	var encrypted bytes.Buffer
	c.NoError(xcrypto.EncryptStreamWithPublicKey(bytes.NewReader(plaintext), &encrypted, publicKey))
	extended := append(encrypted.Bytes(), make([]byte, 32)...)
	var decrypted bytes.Buffer
	c.HasError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(extended), &decrypted, privateKey))
}

// TestDecryptStreamDetectsReorderedChunks verifies that swapping two whole chunks is detected, since each chunk's
// position is bound into its authentication tag.
func TestDecryptStreamDetectsReorderedChunks(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	publicKey := &privateKey.PublicKey
	const chunkSize = 64 * 1024
	plaintext := make([]byte, 2*chunkSize)
	_, err = crypto_rand.Read(plaintext)
	c.NoError(err)
	var encrypted bytes.Buffer
	c.NoError(xcrypto.EncryptStreamWithPublicKey(bytes.NewReader(plaintext), &encrypted, publicKey))
	data := encrypted.Bytes()
	headerSize := privateKey.Size() + 7
	encChunkSize := chunkSize + 16
	// Swap the first and second encrypted chunks in place.
	first := append([]byte(nil), data[headerSize:headerSize+encChunkSize]...)
	second := append([]byte(nil), data[headerSize+encChunkSize:headerSize+2*encChunkSize]...)
	copy(data[headerSize:], second)
	copy(data[headerSize+encChunkSize:], first)
	var decrypted bytes.Buffer
	c.HasError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(data), &decrypted, privateKey))
}

// TestDecryptStreamWrongKeyFails verifies that a stream encrypted to one key cannot be decrypted with a different key.
func TestDecryptStreamWrongKeyFails(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	otherKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	plaintext := []byte("The quick brown fox jumps over the lazy dog.")
	var encrypted bytes.Buffer
	c.NoError(xcrypto.EncryptStreamWithPublicKey(bytes.NewReader(plaintext), &encrypted, &privateKey.PublicKey))
	var decrypted bytes.Buffer
	c.HasError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(encrypted.Bytes()), &decrypted, otherKey))
}
