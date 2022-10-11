package jobs

import (
	"context"
	"log"
	"time"

	"github.com/tony-tvu/goexpense/db"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/teller"
	"go.mongodb.org/mongo-driver/bson"
)

type Jobs struct {
	Db                   *db.MongoDb
	Enabled              bool
	BalancesInterval     int
	TransactionsInterval int
	TellerClient         *teller.TellerClient
}

func (j *Jobs) Start(ctx context.Context) {
	if j.Enabled {
		go j.refreshTransactionsTask(ctx)
		go j.refreshBalancesTask(ctx)
	}
}

func (t *Jobs) refreshBalancesTask(ctx context.Context) {
	for {
		time.Sleep(time.Duration(t.BalancesInterval) * time.Second)
		
		var enrollments []*models.Enrollment
		cursor, _ := t.Db.Enrollments.Find(ctx, bson.M{})
		if err := cursor.All(ctx, &enrollments); err != nil {
			log.Printf("error finding enrollments in refreshBalancesTask: %v\n", err)
		}

		log.Printf("refreshing balances for %d enrollments\n", len(enrollments))
		for _, enrollment := range enrollments {
			t.TellerClient.RefreshBalances(&enrollment.AccessToken)
		}
	}
}

func (t *Jobs) refreshTransactionsTask(ctx context.Context) {
	for {
		time.Sleep(time.Duration(t.TransactionsInterval) * time.Second)

		var enrollments []*models.Enrollment
		cursor, _ := t.Db.Enrollments.Find(ctx, bson.M{})
		if err := cursor.All(ctx, &enrollments); err != nil {
			log.Printf("error finding enrollments in refreshBalancesTask: %v\n", err)
		}

		log.Printf("refreshing transactions for %d enrollments\n", len(enrollments))
		for _, enrollment := range enrollments {
			t.TellerClient.RefreshTransactions(&enrollment.AccessToken)
		}
	}
}
