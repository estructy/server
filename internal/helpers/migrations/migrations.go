// Package migrations provides functions to run database migrations using the golang-migrate package.
package migrations

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrationsLogger struct{}

func newMigrateLogger() *migrationsLogger {
	return &migrationsLogger{}
}

func (l *migrationsLogger) Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func (l *migrationsLogger) Verbose() bool {
	return true
}

func Up() {
	logger := newMigrateLogger()
	m, err := migrate.New(
		"file://internal/infra/database/migrations",
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	m.Log = logger

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to run")
			return
		}

		log.Fatal(err)
		os.Exit(1)
	}
}
