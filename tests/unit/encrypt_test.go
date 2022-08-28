package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/auth"
)

func TestEncryption(t *testing.T) {
	t.Run("should encrpyt and decrypt passwords correctly", func(t *testing.T) {
		t.Parallel()

		key := "ThisKeyStringIs32BytesLongTest01"
		pw := "thisIsAPassword"

		ciphertext, _ := auth.Encrypt(key, pw)

		// should have ciphertext not equal password string
		assert.Equal(t, true, ciphertext != pw)

		// should have decrypyted password equal original password
		decrypted, _ := auth.Decrypt(key, ciphertext)
		assert.Equal(t, pw, decrypted)
	})
}
