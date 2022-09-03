package scheduledtasks

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// run go routine every x seconds
var taskInterval int
var db *gorm.DB

func init() {
	godotenv.Load(".env")
	taskIntervalInt, err := strconv.Atoi(os.Getenv("TASK_INTERVAL"))
	if err != nil {
		// default: 10m
		taskInterval = 600
	} else {
		taskInterval = taskIntervalInt
	}
}

func Start(gDb *gorm.DB) {
	db = gDb
	go refreshTransactions()
}
