package auth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Function verifies if user is logged in and tokens are valid
// Refreshes access token if it has expired and extends sessions
// Returns user ID and type
func AuthorizeUser(c *gin.Context, db *database.MongoDb) (*primitive.ObjectID, *models.Type, error) {
	ctx := c.Request.Context()
	var userIDHex string
	var userTypeStr string

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

	userIDHex = refreshClaims.UserID
	userTypeStr = refreshClaims.UserType
	objID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, nil, errors.New("internal server error")
	}

	// handle expired or missing access_token
	_, err = c.Request.Cookie("goexpense_access")
	if err != nil {

		// find existing session
		var session *models.Session
		if err = db.Sessions.FindOne(ctx, bson.D{{Key: "user_id", Value: objID}}).Decode(&session); err != nil {
			return nil, nil, errors.New("not authorized")
		}

		// validate refresh_token from db
		_, err := ValidateTokenAndGetClaims(session.RefreshToken)
		if err != nil {
			return nil, nil, errors.New("not authorized")
		}

		// renew access_token
		renewedAccess, err := GetEncryptedToken(AccessToken, userIDHex, userTypeStr)
		if err != nil {
			return nil, nil, errors.New("internal server error")
		}
		util.SetCookie(c.Writer, "goexpense_access", renewedAccess.Value, renewedAccess.ExpiresAt)

		// extend user session
		renewedRefresh, err := GetEncryptedToken(RefreshToken, userIDHex, userTypeStr)
		if err != nil {
			return nil, nil, errors.New("internal server error")
		}

		_, err = db.Users.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "refresh_token", Value: renewedRefresh.Value},
					{Key: "expires_at", Value: renewedRefresh.ExpiresAt},
					{Key: "updated_at", Value: time.Now()},
				}},
			},
		)
		if err != nil {
			return nil, nil, errors.New("internal server error")
		}
		
		util.SetCookie(c.Writer, "goexpense_refresh", renewedRefresh.Value, renewedRefresh.ExpiresAt)
	}

	userType := models.GetUserType(userTypeStr)
	return &objID, &userType, nil
}
