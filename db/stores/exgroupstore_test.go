package stores

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TODO remake this running init script from migrations
var (
	pgContainer    *postgres.PostgresContainer
	egs            *EGS
	defaultExGroup = ExGroup{
		Name:   "BodyBack",
		UserId: 1,
	}
)

func initPostgresContainer() {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	containerCreated, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:latest"),
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
	pgContainer = containerCreated
}

// init EGS object, that is tested
func initEGS() {
	ctx := context.Background()
	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal("error creating connection string" + err.Error())
	}
	conn, err := sqlx.Open("postgres", connString)
	if err != nil {
		log.Fatal("error opening connection" + err.Error())
	}
	slog.Info("successful creating of test container")
	egs = NewEGS(conn)
}

func createDefaultExGroup() (int64, error) {
	return egs.Save(defaultExGroup)
}

// default functionality before and after each test
func TestMain(m *testing.M) {
	//setting up
	if pgContainer == nil {
		initPostgresContainer()
	}
	if egs == nil {
		initEGS()
	}
	m.Run()
	// removing database container after all tests passed
	defer func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
}
func clearTable() {
	egs.conn.Exec("DELETE FROM exercise_groups")
}
func TestEGSSaveMethod(t *testing.T) {
	result, err := createDefaultExGroup()
	fmt.Println(result)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(clearTable)
}

func TestEGSFindById(t *testing.T) {
	//negative case
	found, err := egs.FindById(20)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
	//positive case
	result, err := createDefaultExGroup()
	if err != nil {
		t.Error("error saving exgroup")
	}
	found, err = egs.FindById(result)
	if err != nil {
		t.Errorf("failed to find exgroup with id:%d", result)
	}
	if found.Name != defaultExGroup.Name || found.UserId != defaultExGroup.UserId {
		t.Error("got wrong data")
	}
	t.Cleanup(clearTable)
}

func TestEGSFindByName(t *testing.T) {
	//negative case
	_, err := egs.FindByName("Nosuchname")
	if err != nil && err != sql.ErrNoRows {
		t.Error("error, while trying to find unexisting exgroup")
	}
	//positive case
	_, err = createDefaultExGroup()
	if err != nil {
		t.Fatal("error saving exgroup")
	}
	found, err := egs.FindByName("BodyBack")
	if err != nil {
		t.Fatalf("error finding exgroup by name")
	}
	slog.Info(string(fmt.Sprintf("successfully found: %#v", &found)))
	t.Cleanup(clearTable)
}

func TestEGSDeleteById(t *testing.T) {
	//negative case
	err := egs.DeleteById(100)
	if err != nil && err != NotDeleted {
		t.Errorf("error deleting unexisted exGroup:%s", err.Error())
	}
	//positive case
	result, err := createDefaultExGroup()
	if err != nil {
		t.Error(err)
	}
	slog.Info(fmt.Sprintf("created exgroup id's %d", result))
	err = egs.DeleteById(result)
	if err != nil {
		t.Errorf("error deleting exgroup: %s", err.Error())
	}
	found, err := egs.FindById(result)
	if err != nil && err != sql.ErrNoRows {
		t.Errorf("found some exgroup after deletion: %#v", &found)
	}

	t.Cleanup(clearTable)
}

func TestEGSUpdate(t *testing.T) {
	//negative case
	updatedExGroup := ExGroup{
		Name:   "Updated",
		UserId: 3,
		Id:     0,
	}
	err := egs.Update(updatedExGroup)
	if err != nil && err != NotUpdated {
		t.Error(err)
	}

	//positive case
	result, err := createDefaultExGroup()
	if err != nil {
		t.Error("error creating exgroup")
	}
	updatedExGroup = ExGroup{
		Name:   "Updated",
		UserId: 3,
		Id:     result,
	}
	err = egs.Update(updatedExGroup)
	if err != nil {
		t.Fatal(err)
	}
	found, err := egs.FindById(result)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(found, updatedExGroup); diff != "" {
		t.Fatal(diff)
	}
	t.Cleanup(clearTable)
}

func TestFindByUserId(t *testing.T) {
	//negative case
	_, err := egs.FindByUserId(2)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
	//positive case
	for i := 0; i < 3; i++ {
		createDefaultExGroup()
	}
	res, err := egs.FindByUserId(defaultExGroup.UserId)
	if err != nil {
		t.Error(err)
	}
	if len(res) != 3 {
		t.Errorf("error finding exgroups of user")
	}
}
