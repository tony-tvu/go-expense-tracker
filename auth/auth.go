package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/util"
	"gorm.io/gorm"
)

// Function verifies if user is logged in and tokens are valid
// Refreshes access token if it has expired and extends sessions
// Returns user ID and type
func AuthorizeUser(c *gin.Context, db *gorm.DB) (*uint, *string, error) {
	var userID uint
	var userType string

	// no refresh cookie means session has expired or user is not logged in
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
		var session *entity.Session
		if result := db.Where("user_id = ?", userID).First(&session); result.Error != nil {
			return nil, nil, errors.New("not authorized")
		}

		// validate refresh_token from db
		_, err := ValidateTokenAndGetClaims(session.RefreshToken)
		if err != nil {
			return nil, nil, errors.New("not authorized")
		}

		// renew access_token
		renewedAccess, err := GetEncryptedToken(AccessToken, userID, userType)
		if err != nil {
			return nil, nil, errors.New("internal server error")
		}
		util.SetCookie(c.Writer, "goexpense_access", renewedAccess.Value, renewedAccess.ExpiresAt)

		// extend user session
		renewedRefresh, err := GetEncryptedToken(RefreshToken, userID, userType)
		if err != nil {
			return nil, nil, errors.New("internal server error")
		}
		session.RefreshToken = renewedRefresh.Value
		session.ExpiresAt = renewedRefresh.ExpiresAt
		db.Save(&session)

		util.SetCookie(c.Writer, "goexpense_refresh", renewedRefresh.Value, renewedRefresh.ExpiresAt)
	}

	return &userID, &userType, nil
}
