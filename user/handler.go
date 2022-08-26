package user

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	App *app.App
}

func (h UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie("goexpense_access")

	// no access token - make user log in
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type Claims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}

	decrypted, err := auth.Decrypt(h.App.EncryptionKey, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	/*
		No need to check if access token is valid/expired because LoggedIn middleware
		has already validated it and might've refreshed the access token already, which
		would be in w http.ResponseWriter. We only care about the user_id from the
		original r *http.Request here.
	*/
	tkn, _ := jwt.ParseWithClaims(decrypted, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(h.App.JwtKey), nil
		})
	claims := tkn.Claims.(*Claims)

	var u *models.User
	objID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.App.Users.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// do not send back hashed password
	u.Password = ""
	json.NewEncoder(w).Encode(u)
}

// TODO: Admin-only handler to send email invitation for new user