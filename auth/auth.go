package auth

import (
	"errors"

	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/util"
	"gorm.io/gorm"
)

// Function verifies if user is logged in and tokens are valid
// Refreshes access token if it has expired and returns user ID and type
func VerifyUser(c *middleware.WriterAndCookies, db *gorm.DB) (*uint, *string, error) {
	var userID uint
	var userType string

	if util.ContainsEmpty(c.EncryptedRefreshToken) {
		return nil, nil, errors.New("not authorized")
	}

	// validate refresh_token
	refreshClaims, err := ValidateTokenAndGetClaims(c.EncryptedRefreshToken)
	if err != nil {
		return nil, nil, errors.New("not authorized")
	}

	userID = refreshClaims.UserID
	userType = refreshClaims.UserType

	// check if access_token has expired
	_, err = ValidateTokenAndGetClaims(c.EncryptedAccessToken)

	// handle expired or missing access_token
	if err != nil {

		// find existing session
		var s *entity.Session
		if result := db.Where("user_id = ?", refreshClaims.UserID).First(&s); result.Error != nil {
			return nil, nil, errors.New("not authorized")
		}

		// verify token from db session matches request's token
		if s.RefreshToken != c.EncryptedRefreshToken {
			return nil, nil, errors.New("not authorized")
		}

		// validate refresh_token from db
		_, err := ValidateTokenAndGetClaims(s.RefreshToken)
		if err != nil {
			return nil, nil, errors.New("not authorized")
		}

		// renew access_token
		renewed, err := GetEncryptedToken(AccessToken, refreshClaims.UserID, refreshClaims.UserType)
		if err != nil {
			return nil, nil, errors.New("not authorized")
		}

		c.SetToken("goexpense_access", renewed.Value, renewed.ExpiresAt)
	}

	return &userID, &userType, nil
}
