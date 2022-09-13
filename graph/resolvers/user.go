package resolvers

import (
	"context"
	"time"

	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/graph"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/util"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/crypto/bcrypt"
)

// Verifies if user is currently logged in
func (r *queryResolver) IsLoggedIn(ctx context.Context) (bool, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if _, _, err := auth.VerifyUser(c, r.Db); err != nil {
		return false, nil
	}

	return true, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input graph.NewUserInput) (*graph.User, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if _, uType, err := auth.VerifyUser(c, r.Db); err != nil || graph.UserType(*uType) != graph.UserTypeAdmin {
		return nil, gqlerror.Errorf("not authorized")
	}

	panic("not implemented")
}

func (r *queryResolver) Users(ctx context.Context) ([]*graph.User, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if _, uType, err := auth.VerifyUser(c, r.Db); err != nil || graph.UserType(*uType) != graph.UserTypeAdmin {
		return nil, gqlerror.Errorf("not authorized")
	}

	var users []*graph.User
	r.Db.Raw("SELECT * FROM users").Scan(&users)

	return users, nil
}

func (r *mutationResolver) Login(ctx context.Context, input graph.LoginInput) (bool, error) {
	c := middleware.GetWriterAndCookies(ctx)

	// TODO: scrub user input

	// find existing user account
	var u *entity.User
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
	refreshToken, err := auth.GetEncryptedToken(auth.RefreshToken, u.ID, string(u.Type))
	if err != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	// delete existing sessions
	if result := r.Db.Exec("DELETE FROM sessions WHERE user_id = ?", u.ID); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	// save new session
	if result := r.Db.Create(&entity.Session{
		UserID:       u.ID,
		RefreshToken: refreshToken.Value,
		ExpiresAt:    refreshToken.ExpiresAt,
	}); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	// create access token
	accessToken, err := auth.GetEncryptedToken(auth.AccessToken, u.ID, string(u.Type))
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

	if result := r.Db.Exec("DELETE FROM sessions WHERE user_id = ?", claims.UserID); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	c.SetToken("goexpense_access", "", time.Now())
	c.SetToken("goexpense_refresh", "", time.Now())

	return true, nil
}

func (r *queryResolver) UserInfo(ctx context.Context) (*graph.User, error) {
	c := middleware.GetWriterAndCookies(ctx)

	id, _, err := auth.VerifyUser(c, r.Db)
	if err != nil {
		return nil, gqlerror.Errorf("not authorized")
	}

	var u *graph.User
	if result := r.Db.Where("id = ?", id).First(&u); result.Error != nil {
		return nil, gqlerror.Errorf("user not found")
	}

	return u, nil
}

func (r *queryResolver) Sessions(ctx context.Context) ([]*graph.Session, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if _, uType, err := auth.VerifyUser(c, r.Db); err != nil || graph.UserType(*uType) != graph.UserTypeAdmin {
		return nil, gqlerror.Errorf("not authorized")
	}
	var sessions []*graph.Session
	r.Db.Raw("SELECT * FROM sessions").Scan(&sessions)

	return sessions, nil
}

// TODO: add resolver for first time user logins to set password (from an admin invite / user creation)
