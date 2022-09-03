package scheduledtasks

import (
	"log"
	"time"

	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/plaidapi"
)

func refreshTransactions() {
	for {
		log.Println("Running refreshTransactions task")

		var items []entity.Item
		if result := db.Raw("SELECT * FROM items;").Scan(&items); result.Error != nil {
			log.Printf("error occurred during refreshTransactions task: %+v\n", result.Error)
		}

		log.Printf("Length of items: %d\n", len(items))

		for _, item := range(items) {
			transactions, cursor, err := plaidapi.GetTransactions(nil, item.AccessToken)
			if err != nil {
				log.Printf("error occurred during refreshTransactions task: %+v\n", err)
			}

			log.Printf("new cursor: %s\n", *cursor)

			log.Printf("%+v\n", transactions)
		}




		time.Sleep(time.Duration(taskInterval) * time.Second)
	}
}



