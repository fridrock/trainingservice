package converters

import (
	"encoding/json"
	"errors"

	"github.com/fridrock/trainingservice/db/stores"
)

type ExGroupPropeties struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
}

type UpdateExGroup struct {
	UserId  int64  `json:"user_id"`
	Name    string `json:"name"`
	NewName string `json:"newname"`
}

type FindByUser struct {
	UserId int64 `json:"user_id"`
}

var (
	emptyField = errors.New("empty Field")
)

func ExGroupToJson(exg stores.ExGroup) ([]byte, error) {
	r, err := json.Marshal(&exg)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func FromJsonToExGroup(exGroupEncoded []byte) (stores.ExGroup, error) {
	var exg stores.ExGroup
	err := json.Unmarshal(exGroupEncoded, &exg)
	return exg, err
}

func ParseExGroupProperties(request []byte) (int64, string, error) {
	var properties ExGroupPropeties
	err := json.Unmarshal(request, &properties)
	if err != nil {
		return 0, "", err
	}
	if properties.UserId == 0 || properties.Name == "" {
		return 0, "", emptyField
	}
	return properties.UserId, properties.Name, err
}

// TODO refactor
func ParseUpdateExGroup(request []byte) (updateQuery UpdateExGroup, err error) {
	err = json.Unmarshal(request, &updateQuery)
	//TODO write handling this in test
	if err != nil {
		return updateQuery, err
	}
	if updateQuery.UserId == 0 || updateQuery.Name == "" || updateQuery.NewName == "" {
		err = emptyField
	}
	return updateQuery, err
}

func ParseFindByUser(request []byte) (int64, error) {
	var userIdRequest FindByUser
	err := json.Unmarshal(request, &userIdRequest)
	//TODO write handling this in test
	if err != nil {
		return 0, err
	}
	if userIdRequest.UserId == 0 {
		return 0, emptyField
	}
	return userIdRequest.UserId, nil
}
