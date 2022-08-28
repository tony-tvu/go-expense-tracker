package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptAndDecrypt(t *testing.T) {
	t.Parallel()

	key := "ThisKeyStringIs32BytesLongTest01"
	pw := "thisIsAPassword"

	ciphertext, _ := Encrypt(key, pw)

	// should have ciphertext not equal password string
	assert.Equal(t, true, ciphertext != pw)

	// should have decrypyted password equal original password string
	decrypted, _ := Decrypt(key, ciphertext)
	assert.Equal(t, pw, decrypted)
}
