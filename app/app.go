package app

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Users             *mongo.Collection
	Sessions          *mongo.Collection
	Router            *gin.Engine
}
