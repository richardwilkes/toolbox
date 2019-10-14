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

// EncryptStreamWithPublicKey copies 'in' to 'out', encrypting the bytes along
// the way. Note that the output stream will be larger than the input stream
// by aes.BlockSize + publicKey.Size() bytes.
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
		S: cipher.NewCFBEncrypter(block, iv),
		W: out,
	}, in); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

// DecryptStreamWithPrivateKey copies 'in' to 'out', decrypting the bytes
// along the way. Note that the output stream will be smaller than the input
// stream by aes.BlockSize + publicKey.Size() bytes.
func DecryptStreamWithPrivateKey(in io.Reader, out io.Writer, privateKey *rsa.PrivateKey) error {
	encryptedEncryptionKey := make([]byte, privateKey.PublicKey.Size())
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
		S: cipher.NewCFBDecrypter(block, iv),
		R: in,
	}); err != nil {
		return errs.Wrap(err)
	}
	return nil
}
