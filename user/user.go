package user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	Client     *mongo.Client
	Database   string
	Collection string
}

func (s *UserService) CreateUser(c *gin.Context) string {
	type CreateUserRequestBody struct {
		Name     string
		Email    string
		Password string
	}

	var requestBody CreateUserRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		fmt.Print("what")
		panic(err)
	}

	fmt.Println(requestBody.Email)
	fmt.Println(requestBody.Name)
	fmt.Println(requestBody.Password)

	coll := s.Client.Database(s.Database).Collection(s.Collection)
	doc := bson.D{
		{Key: "email", Value: requestBody.Email},
		{Key: "name", Value: requestBody.Name},
		{Key: "password", Value: requestBody.Password},
	}

	_, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}

	return "user created"
}
