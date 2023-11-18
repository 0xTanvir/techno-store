package pg

import (
	"log"
	"os"
	"testing"

	"techno-store/config"
	"techno-store/internal/domain/definition"
)

var testStore definition.DataStore

func TestMain(m *testing.M) {
	appConfig, err := config.Parse()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Get a datastore instance
	testStore = GetInstance(appConfig.Db)

	// Run the tests
	exitCode := m.Run()

	os.Exit(exitCode)
}
