// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xcrypto

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"math"

	"github.com/richardwilkes/toolbox/v2/errs"
)

const (
	streamKeySize         = 32        // AES-256
	streamChunkSize       = 64 * 1024 // Plaintext bytes sealed per authenticated chunk
	streamNoncePrefixSize = 7         // Random per-stream nonce prefix
	streamCounterSize     = 4         // Per-chunk counter; prefix + counter + 1-byte final flag fill the 12-byte nonce
)

// EncryptStreamWithPublicKey copies 'in' to 'out', encrypting the bytes along the way using authenticated encryption. A
// fresh random AES-256 key is generated for each call and sealed to publicKey with RSA-OAEP. The data is then processed
// in chunks, each sealed with AES-GCM, so any tampering, reordering, or truncation of the output stream is detected
// when it is decrypted. The output stream is larger than the input stream by publicKey.Size() + 7 bytes, plus a 16-byte
// authentication tag for every 64KB chunk of input (with a minimum of one chunk, so empty input still produces output).
func EncryptStreamWithPublicKey(in io.Reader, out io.Writer, publicKey *rsa.PublicKey) error {
	encryptionKey := make([]byte, streamKeySize)
	if _, err := io.ReadFull(rand.Reader, encryptionKey); err != nil {
		return errs.Wrap(err)
	}
	noncePrefix := make([]byte, streamNoncePrefixSize)
	if _, err := io.ReadFull(rand.Reader, noncePrefix); err != nil {
		return errs.Wrap(err)
	}
	encryptedEncryptionKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, encryptionKey, nil)
	if err != nil {
		return errs.Wrap(err)
	}
	gcm, err := newStreamAEAD(encryptionKey)
	if err != nil {
		return err
	}
	if _, err = out.Write(encryptedEncryptionKey); err != nil {
		return errs.Wrap(err)
	}
	if _, err = out.Write(noncePrefix); err != nil {
		return errs.Wrap(err)
	}
	reader := bufio.NewReader(in)
	nonce := make([]byte, gcm.NonceSize())
	copy(nonce, noncePrefix)
	plaintext := make([]byte, streamChunkSize)
	var ciphertext []byte
	for counter := uint32(0); ; counter++ {
		n, readErr := io.ReadFull(reader, plaintext)
		if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
			return errs.Wrap(readErr)
		}
		// The chunk is the last one if reading it hit EOF, or if no further data follows it.
		last := readErr == io.EOF || readErr == io.ErrUnexpectedEOF
		if !last {
			if _, peekErr := reader.Peek(1); peekErr != nil {
				if peekErr != io.EOF {
					return errs.Wrap(peekErr)
				}
				last = true
			}
		}
		if !last && counter == math.MaxUint32 {
			return errs.New("stream too large to encrypt")
		}
		fillStreamNonce(nonce, counter, last)
		ciphertext = gcm.Seal(ciphertext[:0], nonce, plaintext[:n], nil)
		if _, err = out.Write(ciphertext); err != nil {
			return errs.Wrap(err)
		}
		if last {
			return nil
		}
	}
}

// DecryptStreamWithPrivateKey copies 'in' to 'out', decrypting the bytes along the way and verifying their integrity.
// It reverses EncryptStreamWithPublicKey: the AES-256 key is recovered with RSA-OAEP and each chunk is opened with
// AES-GCM. An error is returned, and no further plaintext is written, if any chunk fails authentication, which happens
// if the stream was modified, reordered, or truncated. The output stream is smaller than the input stream by
// privateKey.Size() + 7 bytes, plus a 16-byte authentication tag for every chunk.
func DecryptStreamWithPrivateKey(in io.Reader, out io.Writer, privateKey *rsa.PrivateKey) error {
	encryptedEncryptionKey := make([]byte, privateKey.Size())
	if _, err := io.ReadFull(in, encryptedEncryptionKey); err != nil {
		return errs.Wrap(err)
	}
	noncePrefix := make([]byte, streamNoncePrefixSize)
	if _, err := io.ReadFull(in, noncePrefix); err != nil {
		return errs.Wrap(err)
	}
	encryptionKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedEncryptionKey, nil)
	if err != nil {
		return errs.Wrap(err)
	}
	gcm, err := newStreamAEAD(encryptionKey)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(in)
	nonce := make([]byte, gcm.NonceSize())
	copy(nonce, noncePrefix)
	ciphertext := make([]byte, streamChunkSize+gcm.Overhead())
	var plaintext []byte
	for counter := uint32(0); ; counter++ {
		n, readErr := io.ReadFull(reader, ciphertext)
		if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
			return errs.Wrap(readErr)
		}
		// The chunk is the last one if reading it hit EOF, or if no further data follows it.
		last := readErr == io.EOF || readErr == io.ErrUnexpectedEOF
		if !last {
			if _, peekErr := reader.Peek(1); peekErr != nil {
				if peekErr != io.EOF {
					return errs.Wrap(peekErr)
				}
				last = true
			}
		}
		if n < gcm.Overhead() {
			return errs.New("truncated or corrupt encrypted stream")
		}
		if !last && counter == math.MaxUint32 {
			return errs.New("truncated or corrupt encrypted stream")
		}
		fillStreamNonce(nonce, counter, last)
		if plaintext, err = gcm.Open(plaintext[:0], nonce, ciphertext[:n], nil); err != nil {
			return errs.Wrap(err)
		}
		if _, err = out.Write(plaintext); err != nil {
			return errs.Wrap(err)
		}
		if last {
			return nil
		}
	}
}

// newStreamAEAD creates the AES-GCM AEAD used to seal and open individual stream chunks.
func newStreamAEAD(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return gcm, nil
}

// fillStreamNonce writes the per-chunk counter and final-chunk flag into the portion of the nonce that follows the
// random per-stream prefix. Binding the counter and flag into the nonce makes reordering, dropping, duplicating, or
// truncating chunks fail authentication.
func fillStreamNonce(nonce []byte, counter uint32, last bool) {
	binary.BigEndian.PutUint32(nonce[streamNoncePrefixSize:streamNoncePrefixSize+streamCounterSize], counter)
	if last {
		nonce[streamNoncePrefixSize+streamCounterSize] = 1
	} else {
		nonce[streamNoncePrefixSize+streamCounterSize] = 0
	}
}
