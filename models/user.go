package models

type Role string

const (
	AdminUser    Role = "Admin"
	ExternalUser Role = "External"
)

type User struct {
	Name     string
	Email    string
	Password string
	Role     Role
	Verified bool
}
