package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var _ErrInvalidNonceSize = errors.New("invalid nonce size")

type Encryptor interface {
	Encrypt(text string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type aesEncryptor struct {
	key []byte
}

func NewAesEncryptor(key []byte) Encryptor {
	validKey := make([]byte, 32)
	copy(validKey, key)
	return aesEncryptor{
		key: validKey,
	}
}

func (e aesEncryptor) Encrypt(text string) (string, error) {
	c, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	result := gcm.Seal(nonce, nonce, []byte(text), nil)

	return base64.StdEncoding.EncodeToString(result), nil
}

func (e aesEncryptor) Decrypt(text string) (string, error) {
	c, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	buf, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(buf) < nonceSize {
		return "", _ErrInvalidNonceSize
	}

	nonce, ciphertext := buf[:nonceSize], buf[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)

	return string(plaintext), nil
}
