package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

func Encrypt(key string, data string) (string, error) {
	blockCipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(key string, data string) (string, error) {
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	blockCipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	nonce, ciphertext := dataBytes[:gcm.NonceSize()], dataBytes[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
