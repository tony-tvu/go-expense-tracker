package auth

import (
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/util"
	"gorm.io/gorm"
)

// Function verifies if user is logged in and tokens are valid
// Refreshes access tokens if they've expired and refresh token is valid
func IsAuthorized(c *middleware.WriterAndCookies, db *gorm.DB) bool {
	if util.ContainsEmpty(c.EncryptedRefreshToken) {
		return false
	}

	// validate refresh_token
	refreshClaims, err := ValidateTokenAndGetClaims(c.EncryptedRefreshToken)
	if err != nil {
		return false
	}

	// check if access_token has expired
	_, err = ValidateTokenAndGetClaims(c.EncryptedAccessToken)

	// handle expired or missing access_token
	if err != nil {

		// find existing session
		var s *entity.Session
		if result := db.Where("username = ?", refreshClaims.Username).First(&s); result.Error != nil {
			return false
		}

		// verify token from db session matches request's token
		if s.RefreshToken != c.EncryptedRefreshToken {
			return false
		}

		// validate refresh_token from db
		_, err := ValidateTokenAndGetClaims(s.RefreshToken)
		if err != nil {
			return false
		}

		// renew access_token
		renewed, err := GetEncryptedToken(AccessToken, refreshClaims.Username, refreshClaims.UserType)
		if err != nil {
			return false
		}

		c.SetToken("goexpense_access", renewed.Value, renewed.ExpiresAt)
	}

	return true
}

// Function verifies if user is an admin
func IsAdmin(c *middleware.WriterAndCookies) bool {
	if util.ContainsEmpty(c.EncryptedRefreshToken) {
		return false
	}

	claims, err := ValidateTokenAndGetClaims(c.EncryptedRefreshToken)
	if err != nil {
		return false
	}

	if claims.UserType != string(entity.AdminUser) {
		return false
	}
	return true
}
