package smoketest

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/tony-tvu/goexpense/server"
	// "github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Setup
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(time.Second*5))
	defer cancel()
	s := server.Server{}
	s.Init(ctx,
		"testing",
		"5000",
		"ThisKeyStringIs32BytesLongTest01",
		"jwt_key",
		"10",
		"5",
		"mongodb://localhost:27017/local_db",
		"goexpense_test",
		"sandbox",
	)

	s.App.PlaidEnv = "sandbox"

	// s.Run(ctx)
	exitVal := m.Run()

	// Teardown
	os.Exit(exitVal)
}

func TestA(t *testing.T) {
	fmt.Println("TestA running")
}

func TestB(t *testing.T) {
	log.Println("TestB running")
}
