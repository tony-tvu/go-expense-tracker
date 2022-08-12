package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	// given
	key := []byte("ThisKeyStringIs32BytesLongTest01")
	pw := "thisIsAPassword"

	// when
	ciphertext, _ := Encrypt(key, pw)

	// then: verify encrypted pw isn't the same as pw string
	assert.Equal(t, true, ciphertext != pw)

	// and: decrypted pw matches original pw string
	decrypted, _ := Decrypt(key, ciphertext)
	assert.Equal(t, pw, decrypted)
}
