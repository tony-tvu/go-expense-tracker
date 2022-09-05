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
	"github.com/tony-tvu/goexpense/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&entity.Session{})
	db.AutoMigrate(&entity.User{})

	testApp = &app.App{}
	testApp.Start(ctx)
	testApp.Db = db

	// clear tables
	db.Exec("delete from users;")
	db.Exec("delete from sessions;")

	// start test server
	srv = httptest.NewServer(testApp.Router)

	// run tests
	exitVal := m.Run()

	// teardown
	os.Exit(exitVal)
}
