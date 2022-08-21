package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/models"
)

type Claims struct {
	UserId string `json:"user_id"`
	Role   models.Role
	jwt.RegisteredClaims
}

type Token struct {
	Value     string
	ExpiresAt time.Time
}

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

/*
Refresh tokens are saved in 'sessions' collection upon successful login.
They are used to generate new access tokens and verify user has logged in.
These can be revoked by deleting the refresh_token in the collection.

Default expiration time: 24 hours
*/
func CreateRefreshToken(ctx context.Context, a *app.App, u *models.User) (Token, error) {
	exp := time.Now().Add(time.Duration(a.RefreshTokenExp) * time.Second)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserId: u.ObjectID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	refreshTokenStr, err := refreshToken.SignedString(a.JwtKey)
	if err != nil {
		return Token{}, errors.New("error creating refresh token")
	}

	return Token{Value: refreshTokenStr, ExpiresAt: exp}, nil
}

/*
Access tokens are used to protect role-based endpoints. When these expire,
use the refresh_token from the request cookie and verify it is the same token
in the 'sessions' collection and verify its validity. Use the refresh_token to generate a new access token. If the refresh_token has expired or is not valid,
make the user login again to create a new session/refresh_token.

Default expiration time: 15m
*/
func CreateAccessToken(ctx context.Context, a *app.App, u *models.User) (Token, error) {
	exp := time.Now().Add(
		time.Duration(a.AccessTokenExp) * time.Second)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserId: u.ObjectID.Hex(),
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	accessTokenStr, err := accessToken.SignedString(a.JwtKey)
	if err != nil {
		return Token{}, errors.New("error creating access token")
	}

	return Token{Value: accessTokenStr, ExpiresAt: exp}, nil
}

func IsTokenValid(a *app.App, tokenStr string) bool {
	tkn, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return a.JwtKey, nil
		})

	if err != nil {
		return false
	}
	if !tkn.Valid {
		return false
	}

	return true
}

func IsAdmin(a *app.App, tokenStr string) bool {
	tkn, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return a.JwtKey, nil
		})

	if err != nil {
		return false
	}
	if !tkn.Valid {
		return false
	}
	claims := tkn.Claims.(*Claims)
	return claims.Role == models.AdminUser
}

// func AuthenticateToken() (Token, Token) {

// }
