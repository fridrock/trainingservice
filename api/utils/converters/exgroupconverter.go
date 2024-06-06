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
		return 0, "", errors.New("empty field")
	}
	return properties.UserId, properties.Name, err
}
