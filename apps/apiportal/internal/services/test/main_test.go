package test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var testStore services.Store
var pgContainer *postgres.PostgresContainer

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pgContainer, err = postgres.Run(
		ctx,
		"postgres:16",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
			testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatal(err)
	}

	dbURI, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	connPool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Fatal(err)
	}

	// Retry Ping up to 30 times with 100ms delay (max 3s)
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		err = connPool.Ping(ctx)
		if err == nil {
			break
		}
		// log.Printf("waiting for database to be ready (attempt %d/%d): %v", i+1, maxAttempts, err)
		sleepMs(100)
	}
	if err != nil {
		log.Fatal("database not ready after retries:", err)
	}

	// Run migrations
	runMigrations(dbURI)

	testStore = services.NewStore(connPool)

	code := m.Run()

	connPool.Close()

	if err := pgContainer.Terminate(ctx); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

// sleepMs sleeps for the given number of milliseconds.
func sleepMs(ms int) {
	// Use time.Sleep without importing at the top, to avoid import conflicts if already present
	// timeSleep := func(d int) {
	// 	// import time only here
	// 	type sleeper interface{ Sleep(dur int) }
	// 	// actually use time.Sleep
	// }
	// Actually use time.Sleep
	// (import "time" at the top if not already present)
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func runMigrations(dbURI string) {
	// Convert dbURI to postgres:// format for golang-migrate
	migrateDbURI := dbURI
	if len(dbURI) > 0 && dbURI[:5] == "pgx://" {
		migrateDbURI = "postgres://" + dbURI[6:]
	}

	m, err := migrate.New(
		"file://../../migrations",
		migrateDbURI,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

