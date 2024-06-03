package stores

import (
	"context"
	"log"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TODO make this using test containers
// TODO make this running all migrations

func TestEGSSaveMethod(t *testing.T) {
	// setting up database container
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:latest"),
		// postgres.WithInitScripts(filepath.Join("..", "testscripts", "start.sh")),
		postgres.WithInitScripts(filepath.Join("..", "testscripts", "initexgroup.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal("error creating connection string" + err.Error())
	}
	conn, err := sqlx.Open("postgres", connString)
	if err != nil {
		log.Fatal("error opening connection" + err.Error())
	}
	slog.Info("successful creating of test container")
	egs := NewEGS(conn)
	createExGroup := ExGroup{
		Name:   "BodyBack",
		UserId: 1,
	}
	result, err := egs.Save(createExGroup)
	t.Log(result)
	if err != nil {

		t.Error(err)
	}
}
