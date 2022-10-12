package auth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/db"
	"github.com/tony-tvu/goexpense/types"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	Username     string             `json:"username" bson:"username"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	ExpiresAt    time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// Function verifies if user is logged in and tokens are valid
// Refreshes access token if it has expired and extends sessions
// Returns user ID and type
func AuthorizeUser(c *gin.Context, db *db.MongoDb) (*primitive.ObjectID, *types.UserType, error) {
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
		var session *Session
		if err = db.Sessions.FindOne(ctx, bson.M{"user_id": objID}).Decode(&session); err != nil {
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

		_, err = db.Sessions.UpdateOne(
			ctx,
			bson.M{"_id": session.ID},
			bson.M{
				"$set": bson.M{
					"refresh_token": renewedRefresh.Value,
					"expires_at":    renewedRefresh.ExpiresAt,
					"updated_at":    time.Now(),
				}},
		)
		if err != nil {
			return nil, nil, errors.New("internal server error")
		}

		util.SetCookie(c.Writer, "goexpense_refresh", renewedRefresh.Value, renewedRefresh.ExpiresAt)
	}

	userType := types.GetUserType(userTypeStr)
	return &objID, &userType, nil
}
