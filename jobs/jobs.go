package jobs

import (
	"context"
	"log"
	"time"

	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/teller"
)

type Jobs struct {
	Db           *database.MongoDb
	Interval     int
	Enabled      bool
	TellerClient *teller.TellerClient
}

func (j *Jobs) Start(ctx context.Context) {
	if j.Enabled {
		go j.refreshTransactionsTask(ctx)
		go j.refreshAccountsTask(ctx)
	}
}

func (t *Jobs) refreshAccountsTask(ctx context.Context) {
	for {
		log.Println("refresh accounts called")
		time.Sleep(time.Duration(t.Interval) * time.Second)
	}
}

func (t *Jobs) refreshTransactionsTask(ctx context.Context) {
	for {
		log.Println("refresh transactions called")
		time.Sleep(time.Duration(t.Interval) * time.Second)
	}
}
