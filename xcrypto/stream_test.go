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

// TestEncryptDecryptStreamVariousSizes exercises the chunking boundaries: empty, sub-chunk, and exact and partial
// multiples of the chunk size.
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

// TestDecryptStreamWithPartialReads ensures decryption works when the input returns short reads, as network
// connections, pipes, and files can. The encrypted key and nonce prefix must be read with io.ReadFull.
func TestDecryptStreamWithPartialReads(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	publicKey := &privateKey.PublicKey
	plaintext := []byte("The quick brown fox jumps over the lazy dog.")
	var encrypted bytes.Buffer
	c.NoError(xcrypto.EncryptStreamWithPublicKey(bytes.NewReader(plaintext), &encrypted, publicKey))
	in := iotest.OneByteReader(bytes.NewReader(encrypted.Bytes()))
	var decrypted bytes.Buffer
	c.NoError(xcrypto.DecryptStreamWithPrivateKey(in, &decrypted, privateKey))
	c.Equal(plaintext, decrypted.Bytes())
}

// TestDecryptStreamDetectsBitFlip verifies that flipping any single bit in the ciphertext body causes decryption to
// fail.
func TestDecryptStreamDetectsBitFlip(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	publicKey := &privateKey.PublicKey
	plaintext := []byte("The quick brown fox jumps over the lazy dog.")
	// Flip a bit at every position past the header (RSA-sealed key + nonce prefix) and confirm each is rejected.
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

// TestDecryptStreamDetectsTruncation verifies that removing trailing bytes is detected rather than accepted as a
// shorter but valid stream.
func TestDecryptStreamDetectsTruncation(t *testing.T) {
	c := check.New(t)
	privateKey, err := rsa.GenerateKey(crypto_rand.Reader, 2048)
	c.NoError(err)
	publicKey := &privateKey.PublicKey
	// Use more than one chunk so that truncation leaves at least one intact chunk ahead of the cut. This encrypts to
	// two 64KB chunks followed by a 512-byte final chunk, each carrying a 16-byte authentication tag.
	plaintext := make([]byte, 2*64*1024+512)
	_, err = crypto_rand.Read(plaintext)
	c.NoError(err)
	var encrypted bytes.Buffer
	c.NoError(xcrypto.EncryptStreamWithPublicKey(bytes.NewReader(plaintext), &encrypted, publicKey))
	full := encrypted.Bytes()
	// Cut exactly on a chunk boundary, dropping only the final chunk. What remains is a sequence of intact, correctly
	// sized chunks, so this can only be caught by the final-chunk flag bound into the nonce.
	truncated := full[:len(full)-(512+16)]
	var decrypted bytes.Buffer
	c.HasError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(truncated), &decrypted, privateKey))
	// Cut mid-chunk, losing the final chunk and most of the one before it.
	truncated = full[:len(full)-(64*1024+16)]
	var decrypted2 bytes.Buffer
	c.HasError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(truncated), &decrypted2, privateKey))
	// Removing just the final tag byte must also fail.
	var decrypted3 bytes.Buffer
	c.HasError(xcrypto.DecryptStreamWithPrivateKey(bytes.NewReader(full[:len(full)-1]), &decrypted3, privateKey))
}

// TestDecryptStreamDetectsExtraData verifies that data appended after the final chunk is rejected, since that chunk
// is bound as final by its authentication tag.
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

// TestDecryptStreamDetectsReorderedChunks verifies that swapping two chunks is detected, since each chunk's position
// is bound into its authentication tag.
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

// TestDecryptStreamWrongKeyFails verifies that a stream cannot be decrypted with a different key.
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
