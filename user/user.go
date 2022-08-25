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
	ObjectID  primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Role      Role               `json:"role" bson:"role"`
	Verified  bool               `json:"verified" bson:"verified"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
