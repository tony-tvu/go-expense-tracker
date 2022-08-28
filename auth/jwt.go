package auth

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/util"
)

type Claims struct {
	Email    string
	UserType string
	jwt.RegisteredClaims
}

type Token struct {
	Value     string
	ExpiresAt time.Time
}

var encryptionKey string
var jwtKey string
var refreshTokenExp int
var accessTokenExp int

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}

	encryptionKey = os.Getenv("ENCRYPTION_KEY")
	jwtKey = os.Getenv("JWT_KEY")
	if util.ContainsEmpty(encryptionKey, jwtKey) {
		log.Fatal("auth keys are missing")
	}

	refreshExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP"))
	if err != nil {
		refreshTokenExp = 86400
	} else {
		refreshTokenExp = refreshExp
	}
	accessExp, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	if err != nil {
		accessTokenExp = 900
	} else {
		accessTokenExp = accessExp
	}
}

/*
Refresh tokens are saved in 'sessions' collection upon successful login.
They are used to generate new access tokens and verify user has logged in.
These can be revoked by deleting the refresh_token in the collection.

Default expiration time: 24 hours
*/
func CreateRefreshToken(ctx context.Context, a *app.App, u *models.User) (Token, error) {
	exp := time.Now().Add(time.Duration(refreshTokenExp) * time.Second)
	refreshTokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString([]byte(jwtKey))
	if err != nil {
		return Token{}, errors.New("error creating refresh token")
	}

	encrpyted, err := Encrypt(encryptionKey, refreshTokenStr)
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
func CreateAccessToken(ctx context.Context, a *app.App, email, userType string) (Token, error) {
	exp := time.Now().Add(
		time.Duration(accessTokenExp) * time.Second)
	accessTokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email:    email,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString([]byte(jwtKey))
	if err != nil {
		return Token{}, errors.New("error creating access token")
	}

	encrpyted, err := Encrypt(encryptionKey, accessTokenStr)
	if err != nil {
		return Token{}, errors.New("error creating refresh token")
	}

	return Token{Value: encrpyted, ExpiresAt: exp}, nil
}
