package smoketest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/server"
)

var port string = "5001"
var testURL string = fmt.Sprintf("http://localhost:%s/", port)

func TestMain(m *testing.M) {
	// Setup
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	s := server.Server{}
	s.Init(ctx,
		port,
		"testing",
		"ThisKeyStringIs32BytesLongTest01",
		"jwt_key",
		"10",
		"5",
		"mongodb://localhost:27017/local_db",
		"goexpense_test",
		"plaidClientID",
		"plaidSecret",
		"sandbox",
		"US,CA",
		"auth,transactions",
	)

	go func() {
		if err := s.App.Server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
		cancel()
	}()

	exitVal := m.Run()

	// Teardown
	<-ctx.Done()
	os.Exit(exitVal)
}

func TestA(t *testing.T) {
	type Body struct {
		Message string `json:"message"`
	}

	resp, err := http.Get(fmt.Sprintf("%s/api/health", testURL))
	if err != nil {
		log.Fatalln(err)
	}

	var b Body
	err = json.NewDecoder(resp.Body).Decode(&b)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, "Ok", b.Message)
}

func TestB(t *testing.T) {
	log.Println("TestB running")
}
