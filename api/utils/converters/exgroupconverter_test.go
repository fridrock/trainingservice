package converters

import (
	"testing"

	"github.com/fridrock/trainingservice/db/stores"
	"github.com/google/go-cmp/cmp"
)

var (
	exg = stores.ExGroup{
		Id:     1,
		UserId: 2,
		Name:   "Back",
	}
)

func TestExGroupConverter(t *testing.T) {
	encoded, err := ExGroupToJson(exg)
	if err != nil {
		t.Error(err)
	}
	correctResult := `{"id":1,"user_id":2,"name":"Back"}`
	if correctResult != string(encoded) {
		t.Error("converted value is not correct")
	}
	exgFromJson, err := FromJsonToExGroup(encoded)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(exgFromJson, exg); diff != "" {
		t.Error("got wrong value from json")
	}
}

func TestParseDeleteRequest(t *testing.T) {
	encoded := `{"user_id":3,"name":"Back"}`
	deleteUserId, deleteGroupName, err := ParseDeleteRequest([]byte(encoded))
	if err != nil {
		t.Error(err)
	}
	if deleteUserId != 3 || deleteGroupName != "Back" {
		t.Error("got wrong id")
	}
}
