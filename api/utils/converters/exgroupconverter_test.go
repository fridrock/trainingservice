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

func TestParseExGroupProperties(t *testing.T) {
	//positive case
	encoded := `{"user_id":3,"name":"Back"}`
	deleteUserId, deleteGroupName, err := ParseExGroupProperties([]byte(encoded))
	if err != nil {
		t.Error(err)
	}
	if deleteUserId != 3 || deleteGroupName != "Back" {
		t.Error("got wrong id")
	}
	//negative case
	encoded = `{"user_id":2}`
	deleteUserId, deleteGroupName, err = ParseExGroupProperties([]byte(encoded))
	if err == nil || deleteUserId != 0 || deleteGroupName != "" {
		t.Error("no error with empty field name")
	}
}

func TestParseUpdateExGroup(t *testing.T) {
	data := []struct {
		testName              string
		query                 string
		expectedUpdateExGroup UpdateExGroup
		expectedError         error
	}{
		{
			"negative case: wrong data format",
			`{"user_id":0, "name":""}`,
			UpdateExGroup{},
			emptyField,
		},
		{
			"positive case",
			`{"user_id":13,"name":"back","newname":"newBack"}`,
			UpdateExGroup{
				UserId:  13,
				Name:    "back",
				NewName: "newBack",
			},
			nil,
		},
	}
	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			res, err := ParseUpdateExGroup([]byte(d.query))
			if err != d.expectedError {
				t.Error(err)
			}
			if diff := cmp.Diff(res, d.expectedUpdateExGroup); diff != "" {
				t.Errorf("error while parsing, got wrong values: %s", diff)
			}
		})
	}
}

func TestParseFindByUser(t *testing.T) {
	data := []struct {
		testName       string
		query          string
		expectedUserId int64
		expectedError  error
	}{
		{
			"negative case: wrong input",
			`{"userid":5}`,
			0,
			emptyField,
		},
		{
			"positive case",
			`{"user_id":12}`,
			12,
			nil,
		},
	}
	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			userId, err := ParseFindByUser([]byte(d.query))
			if err != d.expectedError {
				t.Error(err)
			}
			if userId != d.expectedUserId {
				t.Error("got wrong user id")
			}
		})
	}
}
