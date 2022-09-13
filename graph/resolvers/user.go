package resolvers

import (
	"context"
	"time"

	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/graph/models"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/util"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/crypto/bcrypt"
)

// Verifies if user is currently logged in
func (r *queryResolver) IsLoggedIn(ctx context.Context) (bool, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if !auth.IsAuthorized(c, r.Db) {
		return false, nil
	}
	
	return true, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input models.NewUserInput) (*models.User, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if !auth.IsAuthorized(c, r.Db) || !auth.IsAdmin(c, r.Db) {
		return nil, gqlerror.Errorf("not authorized")
	}
	
	panic("not implemented")
}

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if !auth.IsAuthorized(c, r.Db) || !auth.IsAdmin(c, r.Db) {
		return nil, gqlerror.Errorf("not authorized")
	}
	var users []*models.User
	r.Db.Raw("SELECT * FROM users").Scan(&users)

	return users, nil
}

func (r *mutationResolver) Login(ctx context.Context, input models.LoginInput) (bool, error) {
	c := middleware.GetWriterAndCookies(ctx)

	// TODO: scrub user input

	// find existing user account
	var u *models.User
	result := r.Db.Where("username = ?", input.Username).First(&u)
	if result.Error != nil {
		return false, gqlerror.Errorf("user not found")
	}

	// verify password
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))
	if err != nil {
		return false, gqlerror.Errorf("not authorized")
	}

	// create refresh token
	refreshToken, err := auth.GetEncryptedToken(auth.RefreshToken, u.Username, string(u.Type))
	if err != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	// delete existing sessions
	if result := r.Db.Exec("DELETE FROM sessions WHERE username = ?", u.Username); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	// save new session
	if result := r.Db.Create(&models.Session{
		Username:     u.Username,
		RefreshToken: refreshToken.Value,
		ExpiresAt:    refreshToken.ExpiresAt,
	}); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	// create access token
	accessToken, err := auth.GetEncryptedToken(auth.AccessToken, u.Username, string(u.Type))
	if err != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	c.SetToken("goexpense_access", accessToken.Value, accessToken.ExpiresAt)
	c.SetToken("goexpense_refresh", refreshToken.Value, refreshToken.ExpiresAt)

	return true, nil
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if util.ContainsEmpty(c.EncryptedRefreshToken) {
		return false, gqlerror.Errorf("not authorized")
	}

	claims, err := auth.ValidateTokenAndGetClaims(c.EncryptedRefreshToken)
	if err != nil {
		return false, gqlerror.Errorf("invalid token")
	}

	if result := r.Db.Exec("DELETE FROM sessions WHERE username = ?", claims.Username); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	c.SetToken("goexpense_access", "", time.Now())
	c.SetToken("goexpense_refresh", "", time.Now())

	return true, nil
}

func (r *queryResolver) UserInfo(ctx context.Context) (*models.User, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if !auth.IsAuthorized(c, r.Db) {
		return nil, gqlerror.Errorf("not authorized")
	}

	claims, err := auth.ValidateTokenAndGetClaims(c.EncryptedRefreshToken)
	if err != nil {
		return nil, gqlerror.Errorf("invalid token")
	}

	var u *models.User
	if result := r.Db.Where("username = ?", claims.Username).First(&u); result.Error != nil {
		return nil, gqlerror.Errorf("user not found")
	}

	return u, nil
}

func (r *queryResolver) Sessions(ctx context.Context) ([]*models.Session, error) {
	c := middleware.GetWriterAndCookies(ctx)
	
	if !auth.IsAuthorized(c, r.Db) || !auth.IsAdmin(c, r.Db) {
		return nil, gqlerror.Errorf("not authorized")
	}
	var sessions []*models.Session
	r.Db.Raw("SELECT * FROM sessions").Scan(&sessions)

	return sessions, nil
}

// TODO: add resolver for first time user logins to set password (from an admin invite / user creation)

