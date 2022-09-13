// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Item struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userID"`
	ItemID      string    `json:"itemID"`
	AccessToken string    `json:"accessToken"`
	Cursor      string    `json:"cursor"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Type     string `json:"type"`
}

type PageInfo struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalRows  int `json:"totalRows"`
	TotalPages int `json:"totalPages"`
}

type PublicTokenInput struct {
	PublicToken string `json:"publicToken"`
}

type Session struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Transaction struct {
	ID            uint      `json:"id"`
	ItemID        uint      `json:"itemID"`
	UserID        uint      `json:"userID"`
	TransactionID string    `json:"transactionID"`
	Date          time.Time `json:"date"`
	Amount        float64   `json:"amount"`
	Category      string    `json:"category"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type TransactionConnection struct {
	Nodes    []*Transaction `json:"nodes"`
	PageInfo *PageInfo      `json:"pageInfo"`
}

type TransactionSearchInput struct {
	UserID  uint `json:"userID"`
	Page    int  `json:"page"`
	PerPage int  `json:"perPage"`
}

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Type      UserType  `json:"type"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserType string

const (
	UserTypeAdmin   UserType = "ADMIN"
	UserTypeRegular UserType = "REGULAR"
)

var AllUserType = []UserType{
	UserTypeAdmin,
	UserTypeRegular,
}

func (e UserType) IsValid() bool {
	switch e {
	case UserTypeAdmin, UserTypeRegular:
		return true
	}
	return false
}

func (e UserType) String() string {
	return string(e)
}

func (e *UserType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserType", str)
	}
	return nil
}

func (e UserType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
