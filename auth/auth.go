package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/util"
	"gorm.io/gorm"
)

// Function verifies if user is logged in and tokens are valid
// Refreshes access token if it has expired and returns user ID and type
func VerifyUser(c *gin.Context, db *gorm.DB) (*uint, *string, error) {
	var userID uint
	var userType string

	refreshCookie, err := c.Request.Cookie("goexpense_refresh")
	if err != nil {
		return nil, nil, errors.New("not authorized")
	}

	// validate refresh_token
	refreshClaims, err := ValidateTokenAndGetClaims(refreshCookie.Value)
	if err != nil {
		return nil, nil, errors.New("not authorized")
	}

	userID = refreshClaims.UserID
	userType = refreshClaims.UserType

	// handle expired or missing access_token
	_, err = c.Request.Cookie("goexpense_access")
	if err != nil {

		// find existing session
		var s *entity.Session
		if result := db.Where("user_id = ?", refreshClaims.UserID).First(&s); result.Error != nil {
			return nil, nil, errors.New("not authorized")
		}

		// verify token from db session matches request's token
		if s.RefreshToken != refreshCookie.Value {
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

		util.SetCookie(c.Writer, "goexpense_access", renewed.Value, renewed.ExpiresAt)
	}

	return &userID, &userType, nil
}
