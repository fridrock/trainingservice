package stores

import (
	"context"
	"log"
	"log/slog"
	"testing"

	"github.com/fridrock/trainingservice/test"
	"github.com/jmoiron/sqlx"
)

var (
	ts *TS
)

func initTs() {
	ctx := context.Background()
	connString, err := test.GetDatabaseContainer().ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal("error creating connection string" + err.Error())
	}
	conn, err := sqlx.Open("postgres", connString)
	if err != nil {
		log.Fatal("error opening connection" + err.Error())
	}
	slog.Info("successful creating of test container")
	ts = NewTs(conn)
}

func TestTSMain(m *testing.M) {
	//setting up ts
	if ts == nil {
		initTs()
	}
	m.Run()
	defer ts.conn.Close()
}

func TestTSStartTraining(t *testing.T) {

}
