package stores

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"testing"

	"github.com/fridrock/trainingservice/test"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	egs            *EGS
	defaultExGroup = ExGroup{
		Name:   "BodyBack",
		UserId: 1,
	}
)

// init EGS object, that is tested
func initEGS() {
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
	egs = NewEGS(conn)
}

func createDefaultExGroup() (int64, error) {
	return egs.Save(defaultExGroup)
}

// default functionality before and after each test
func TestMain(m *testing.M) {
	//setting up
	if egs == nil {
		initEGS()
	}
	m.Run()
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
	_, err := egs.FindByName(defaultExGroup.UserId, "Nosuchname")
	if err != nil && err != sql.ErrNoRows {
		t.Error("error, while trying to find unexisting exgroup")
	}
	//positive case
	_, err = createDefaultExGroup()
	if err != nil {
		t.Fatal("error saving exgroup")
	}
	found, err := egs.FindByName(defaultExGroup.UserId, defaultExGroup.Name)
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

func TestEGSDeleteByName(t *testing.T) {
	//negative case
	err := egs.DeleteByName(1, "NoSuchExGroup")
	if err != nil && err != NotDeleted {
		t.Error(err)
	}
	//positive case
	createDefaultExGroup()
	err = egs.DeleteByName(defaultExGroup.UserId, defaultExGroup.Name)
	if err != nil {
		t.Errorf("error deleting exgroup: %s", err.Error())
	}
	found, err := egs.FindByName(defaultExGroup.UserId, defaultExGroup.Name)

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
func TestEGSUpdateByName(t *testing.T) {
	//negative case
	err := egs.UpdateByName(3, "NoSuchGroup", "updated")
	if err != nil && err != NotUpdated {
		t.Error(err)
	}
	//positive case
	result, err := createDefaultExGroup()
	if err != nil {
		t.Error("error creating exgroup")
	}
	updatedExGroup := ExGroup{
		Id:     result,
		UserId: 1,
		Name:   "Updated",
	}
	err = egs.UpdateByName(defaultExGroup.UserId, defaultExGroup.Name, "Updated")
	if err != nil {
		t.Fatal(err)
	}
	found, err := egs.FindByName(defaultExGroup.UserId, "Updated")
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
