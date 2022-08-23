package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	AdminUser    Role = "Admin"
	ExternalUser Role = "External"
)

type User struct {
	ObjectID  primitive.ObjectID `bson:"_id" json:"_id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	Role      Role               `json:"role"`
	Verified  bool               `json:"verified"`
	CreatedAt time.Time          `json:"created_at"`
}
