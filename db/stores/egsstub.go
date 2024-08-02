package stores

import (
	"database/sql"
)

type EGSStub struct{}

func (egss EGSStub) Save(group ExGroup) (int64, error) {
	return 1, nil
}
func (egss EGSStub) FindById(id int64) (ExGroup, error) {
	var group ExGroup
	return group, nil
}

func (egss EGSStub) FindByName(userId int64, name string) (group ExGroup, err error) {
	if name == "Unexisting" {
		err = sql.ErrNoRows
	} else {
		group = ExGroup{
			Id:     1,
			UserId: 2,
			Name:   "Back",
		}
	}
	return group, err
}
func (egss EGSStub) DeleteById(id int64) error {
	return nil
}
func (egss EGSStub) DeleteByName(userId int64, name string) error {
	if name == "Unexisting" {
		return NotDeleted
	} else {
		return nil
	}
}
func (egss EGSStub) Update(group ExGroup) error {
	return nil
}
func (egss EGSStub) UpdateByName(userId int64, name string, newName string) error {
	if name == "Unexisting" {
		return NotUpdated
	}
	return nil
}

func (egss EGSStub) FindByUserId(userId int64) ([]ExGroup, error) {
	if userId == 1 {
		return nil, sql.ErrNoRows
	}
	groups := []ExGroup{
		{
			Id:     1,
			Name:   "Back",
			UserId: userId,
		},
		{
			Id:     2,
			Name:   "Front",
			UserId: userId,
		},
		{
			Id:     3,
			Name:   "Chest",
			UserId: userId,
		},
	}
	return groups, nil
}
