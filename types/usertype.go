package types

type UserType string

const (
	AdminUser   UserType = "ADMIN"
	RegularUser UserType = "REGULAR"
)

func GetUserType(s string) UserType {
	if s == "ADMIN" {
		return AdminUser
	} else {
		return RegularUser
	}
}
