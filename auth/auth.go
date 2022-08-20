package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/models"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func Encrypt(key string, data string) (string, error) {
	dataBytes := []byte(data)
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

	ciphertext := gcm.Seal(nonce, nonce, dataBytes, nil)
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

func GetAuthTokens(a *app.App, user *models.User) (string, string, error) {
	jwtKey := []byte(a.JwtKey)

	exp := time.Now().Add(time.Duration(a.RefreshTokenExp) * time.Second)
	refreshClaims := &Claims{
		UserID: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", "", errors.New("error creating jwt tokens")
	}

	return refreshTokenStr, "", nil
}

func IsTokenValid(a *app.App, tokenStr string) bool {
	jwtKey := []byte(a.JwtKey)

	tkn, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	claims := tkn.Claims.(*Claims)

	log.Println(claims)

	if err != nil {
		log.Println("ERROR - NO LONGER VALID")
		return false
	}
	if !tkn.Valid {
		log.Println("NO LONGER VALID")
		return false
	}

	return true
}
