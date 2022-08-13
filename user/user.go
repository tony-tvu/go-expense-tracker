package user

type Role string

const (
	Admin    Role = "Admin"
	External Role = "External"
)

type User struct {
	Name     string
	Email    string
	Password string
	Role     Role
	Verified bool
}
