package tests

import (
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	if util.ContainsEmpty(dbUser, dbPwd, dbHost, dbName, dbPort) {
		log.Fatal("test postgres config envs are missing")
	}

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s port=%s",
		dbUser, dbPwd, dbName, dbHost, dbPort)
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&entity.Session{})
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Item{})
	db.AutoMigrate(&entity.Transaction{})

	testApp = &app.App{}
	testApp.Initialize(ctx)
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
