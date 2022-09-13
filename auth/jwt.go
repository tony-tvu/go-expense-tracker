package auth

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID uint
	UserType string
	jwt.RegisteredClaims
}

type Token struct {
	Value     string
	ExpiresAt time.Time
}

type TokenType string

const (
	AccessToken  TokenType = "Access"
	RefreshToken TokenType = "Refresh"
)

var jwtKey string
var refreshTokenExp int
var accessTokenExp int

func init() {
	godotenv.Load(".env")
	jwtKey = os.Getenv("JWT_KEY")
	if jwtKey == "" {
		log.Fatal("jwt key is missing")
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

Access tokens are used to protect role-based endpoints. When these expire,
get the claims (username and type) from the request's cookie and query the sessions
collection for an existing session using the username. After verifying the refresh token
has not expired, generate a new access token and return it in the response writer's cookie.
If the refresh_token has expired or is not valid, make the user login again to create a
new session/refresh_token.

Default expiration time: 15m
*/
func GetEncryptedToken(tokenType TokenType, userID uint, userType string) (Token, error) {
	var exp time.Time
	if tokenType == RefreshToken {
		exp = time.Now().Add(time.Duration(refreshTokenExp) * time.Second)
	} else {
		exp = time.Now().Add(time.Duration(accessTokenExp) * time.Second)
	}

	accessTokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}).SignedString([]byte(jwtKey))
	if err != nil {
		return Token{}, errors.New("error signing token")
	}

	encrpyted, err := Encrypt(accessTokenStr)
	if err != nil {
		return Token{}, errors.New("error encrypting token")
	}

	return Token{Value: encrpyted, ExpiresAt: exp}, nil
}

// Function decrypts an encrypted token string, validates the token, then returns the claims.
func ValidateTokenAndGetClaims(encryptedTkn string) (*Claims, error) {
	if encryptedTkn == "" {
		return nil, errors.New("encrypted token is empty")
	}
	decrypted, err := Decrypt(encryptedTkn)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(decrypted, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
	if err != nil {
		return nil, err
	}

	// return error if claims are missing
	claims := token.Claims.(*Claims)
	if claims.UserID == 0 || claims.UserType == "" {
		return nil, errors.New("token is missing claims")
	}

	return claims, nil
}
