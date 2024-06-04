package converters

import (
	"encoding/json"

	"github.com/fridrock/trainingservice/db/stores"
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

func ParseDeleteRequest(request []byte) (int64, string, error) {
	type DeleteProperties struct {
		UserId int64  `json:"user_id"`
		Name   string `json:"name"`
	}
	var id DeleteProperties
	err := json.Unmarshal(request, &id)
	return id.UserId, id.Name, err
}
func ParseFindByNameRequest(request []byte) (int64, error) {
	return 0, nil
}
