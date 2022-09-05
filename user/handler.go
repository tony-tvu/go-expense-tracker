package user

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	Db *gorm.DB
}

type CredentialsInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

var v *validator.Validate

func init() {
	v = validator.New()
}

func (h UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	defer c.Request.Body.Close()

	var input CredentialsInput

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

	// validate
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
	refreshToken, err := auth.GetEncryptedRefreshToken(ctx, u)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// delete existing sessions
	if result := h.Db.Exec("DELETE FROM sessions WHERE username = ?", u.Username); result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// save new session
	if result := h.Db.Create(&entity.Session{
		Username:     u.Username,
		RefreshToken: refreshToken.Value,
		ExpiresAt:    refreshToken.ExpiresAt,
	}); result.Error != nil {
		log.Fatal(err)
	}

	// create access token
	accessToken, err := auth.GetEncryptedAccessToken(ctx, u.Username, string(u.Type))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "goexpense_access",
		Value:    accessToken.Value,
		Expires:  accessToken.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "goexpense_refresh",
		Value:    refreshToken.Value,
		Expires:  refreshToken.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h UserHandler) Logout(c *gin.Context) {
	cookie, err := c.Request.Cookie("goexpense_refresh")

	// missing refresh_token
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	claims, err := auth.ValidateTokenAndGetClaims(cookie.Value)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if result := h.Db.Exec("DELETE FROM sessions WHERE username = ?", claims.Username); result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h UserHandler) GetUserInfo(c *gin.Context) {
	cookie, err := c.Request.Cookie("goexpense_refresh")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateTokenAndGetClaims(cookie.Value)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var u *entity.User
	if result := h.Db.Where("username = ?", claims.Username).First(&u); result.Error != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// do not send back hashed password
	u.Password = ""
	c.JSON(200, u)
}

// TODO: Send email invite to new user
func (h UserHandler) InviteUser(c *gin.Context) {
	type Body struct {
		Email string `json:"email"`
	}
	defer c.Request.Body.Close()

	var b Body
	err := json.NewDecoder(c.Request.Body).Decode(&b)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// validate email address
	_, err = mail.ParseAddress(b.Email)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	to := []string{b.Email}

	from := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		log.Println("error - email sender configs are missing")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := "Subject: This is the subject of the mail\n"
	body := "This is the body of the mail"
	message := []byte(subject + body)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h UserHandler) GetSessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User sessions called"})
}

// Handler is used to check if user is logged in
func (h UserHandler) Ping(c *gin.Context) {
	c.Writer.WriteHeader(200)
}
