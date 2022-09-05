package scheduledtasks

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/plaid/plaid-go/plaid"
	"gorm.io/gorm"
)

type ScheduledTasks struct {
	Db     *gorm.DB
	Client *plaid.APIClient
}

// run go routine every x seconds
var taskInterval int

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
