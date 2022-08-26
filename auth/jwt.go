package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/models"
)

type Claims struct {
	UserID   string
	UserType string
	jwt.RegisteredClaims
}

type Token struct {
	Value     string
	ExpiresAt time.Time
}

/*
Refresh tokens are saved in 'sessions' collection upon successful login.
They are used to generate new access tokens and verify user has logged in.
These can be revoked by deleting the refresh_token in the collection.

Default expiration time: 24 hours
*/
func CreateRefreshToken(ctx context.Context, a *app.App, u *models.User) (Token, error) {
	exp := time.Now().Add(time.Duration(a.RefreshTokenExp) * time.Second)
	refreshTokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: u.ObjectID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString([]byte(a.JwtKey))
	if err != nil {
		return Token{}, errors.New("error creating refresh token")
	}

	encrpyted, err := Encrypt(a.EncryptionKey, refreshTokenStr)
	if err != nil {
		return Token{}, errors.New("error creating refresh token")
	}

	return Token{Value: encrpyted, ExpiresAt: exp}, nil
}

/*
Access tokens are used to protect role-based endpoints. When these expire,
use the refresh_token from the request cookie and verify it is the same token
in the 'sessions' collection and verify its validity. Use the refresh_token to generate a new access token. If the refresh_token has expired or is not valid,
make the user login again to create a new session/refresh_token.

Default expiration time: 15m
*/
func CreateAccessToken(ctx context.Context, a *app.App, userID, userType string) (Token, error) {
	exp := time.Now().Add(
		time.Duration(a.AccessTokenExp) * time.Second)
	accessTokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString([]byte(a.JwtKey))
	if err != nil {
		return Token{}, errors.New("error creating access token")
	}

	encrpyted, err := Encrypt(a.EncryptionKey, accessTokenStr)
	if err != nil {
		return Token{}, errors.New("error creating refresh token")
	}

	return Token{Value: encrpyted, ExpiresAt: exp}, nil
}
