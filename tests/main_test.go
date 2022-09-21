package tests

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testApp         *app.App
	srv             *httptest.Server
	ctx             context.Context
	refreshTokenExp int
	accessTokenExp  int
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}

	refreshExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP"))
	if err != nil {
		log.Fatal(err)
	}
	refreshTokenExp = refreshExp

	accessExp, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	if err != nil {
		log.Fatal(err)
	}
	accessTokenExp = accessExp

	testApp = &app.App{}
	testApp.Initialize(ctx)

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	if util.ContainsEmpty(mongoURI, dbName) {
		log.Fatal("test db configs are missing")
	}
	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			log.Println("mongo has been disconnected: ", err)
		}
	}()
	testApp.Db.Configs = mongoclient.Database(dbName).Collection("configs")
	testApp.Db.Items = mongoclient.Database(dbName).Collection("items")
	testApp.Db.Sessions = mongoclient.Database(dbName).Collection("sessions")
	testApp.Db.Transactions = mongoclient.Database(dbName).Collection("transactions")
	testApp.Db.Users = mongoclient.Database(dbName).Collection("users")

	// Create unique constraints
	database.CreateUniqueConstraints(ctx, testApp.Db)

	// clear tables
	testApp.Db.Items.Drop(ctx)
	testApp.Db.Sessions.Drop(ctx)
	testApp.Db.Transactions.Drop(ctx)
	testApp.Db.Users.Drop(ctx)

	// start test server
	srv = httptest.NewServer(testApp.Router)

	// run tests
	exitVal := m.Run()

	// teardown
	os.Exit(exitVal)
}
