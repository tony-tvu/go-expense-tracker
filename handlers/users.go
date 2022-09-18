package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/cache"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	Db    *gorm.DB
	Cache *cache.Configs
}

var v *validator.Validate

func init() {
	v = validator.New()
}

func (h *UserHandler) IsLoggedIn(c *gin.Context) {
	if _, _, err := auth.VerifyUser(c, h.Db); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func (h *UserHandler) IsAdmin(c *gin.Context) {
	if _, userType, err := auth.VerifyUser(c, h.Db); err != nil || *userType != string(entity.AdminUser) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	if _, uType, err := auth.VerifyUser(c, h.Db); err != nil || *uType != string(entity.AdminUser) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var users []*entity.User
	h.Db.Raw("SELECT * FROM users").Scan(&users)

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Login(c *gin.Context) {
	defer c.Request.Body.Close()

	type Input struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var input Input
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = v.Struct(input)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// find existing user account
	var u *entity.User
	result := h.Db.Where("username = ?", input.Username).First(&u)
	if result.Error != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// create refresh token
	refreshToken, err := auth.GetEncryptedToken(auth.RefreshToken, u.ID, string(u.Type))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// delete existing sessions
	if result := h.Db.Exec("DELETE FROM sessions WHERE user_id = ?", u.ID); result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// save new session
	if result := h.Db.Create(&entity.Session{
		UserID:       u.ID,
		RefreshToken: refreshToken.Value,
		ExpiresAt:    refreshToken.ExpiresAt,
	}); result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// create access token
	accessToken, err := auth.GetEncryptedToken(auth.AccessToken, u.ID, string(u.Type))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	util.SetCookie(c.Writer, "goexpense_access", accessToken.Value, accessToken.ExpiresAt)
	util.SetCookie(c.Writer, "goexpense_refresh", refreshToken.Value, refreshToken.ExpiresAt)
}

func (h *UserHandler) Logout(c *gin.Context) {
	userID, _, err := auth.VerifyUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if result := h.Db.Exec("DELETE FROM sessions WHERE user_id = ?", userID); result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	util.SetCookie(c.Writer, "goexpense_access", "", time.Now())
	util.SetCookie(c.Writer, "goexpense_refresh", "", time.Now())
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	id, _, err := auth.VerifyUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var u *entity.User
	if result := h.Db.Where("id = ?", id).First(&u); result.Error != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": u.Username,
		"email":    u.Email,
		"type":     u.Type,
	})
}

func (h *UserHandler) GetSessions(c *gin.Context) {
	if _, uType, err := auth.VerifyUser(c, h.Db); err != nil || *uType != string(entity.AdminUser) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var sessions []*entity.Session
	h.Db.Raw("SELECT * FROM sessions").Scan(&sessions)

	c.JSON(http.StatusOK, sessions)
}

// TODO
func (h UserHandler) RegisterUser(c *gin.Context) {
	panic("not implemented")
}
