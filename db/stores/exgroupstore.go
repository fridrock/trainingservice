package stores

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	NotDeleted = errors.New("no rows deleted")
	NotUpdated = errors.New("no rows updated")
)

// ExGroup struct that is entity for exercise_groups table
type ExGroup struct {
	Id     int64  `db:"id" json:"id"`
	UserId int64  `db:"user_id" json:"user_id"`
	Name   string `db:"name" json:"name"`
}

// ExGroupStore - interface which contains all methods for working with exercise_groups table
type ExGroupStore interface {
	Save(ExGroup) (int64, error)
	FindById(int64) (ExGroup, error)
	FindByName(userId int64, groupName string) (ExGroup, error)
	DeleteById(int64) error
	DeleteByName(userId int64, groupName string) error
	Update(ExGroup) error
	UpdateByName(userId int64, name string, updated ExGroup) error
	FindByUserId(int64) ([]ExGroup, error)
}

// EGS - standard realization of ExGroupInterface
type EGS struct {
	conn *sqlx.DB
}

// NewEGS - function that creates realization for ExGroupStore interface
func NewEGS(conn *sqlx.DB) *EGS {
	return &EGS{
		conn: conn,
	}
}

func (egs EGS) Save(exGroup ExGroup) (int64, error) {
	var exGroupId int64
	q := `INSERT INTO exercise_groups(user_id, name) VALUES($1, $2) RETURNING id`
	err := egs.conn.Get(&exGroupId, q, exGroup.UserId, exGroup.Name)
	return exGroupId, err
}

func (egs EGS) FindById(id int64) (ExGroup, error) {
	var exGroup ExGroup
	q := `SELECT * FROM exercise_groups WHERE id=$1`
	err := egs.conn.Get(&exGroup, q, id)
	return exGroup, err
}

func (egs EGS) FindByName(userId int64, name string) (ExGroup, error) {
	var exGroup ExGroup
	q := `SELECT * FROM exercise_groups WHERE name=$1 and user_id=$2`
	err := egs.conn.Get(&exGroup, q, name, userId)
	return exGroup, err
}

func (egs EGS) DeleteById(id int64) error {
	q := `DELETE FROM exercise_groups WHERE id=$1`
	res, err := egs.conn.Exec(q, id)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return NotDeleted
	}
	return nil
}

func (egs EGS) DeleteByName(userId int64, name string) error {
	q := `DELETE FROM exercise_groups WHERE user_id=$1 AND name=$2`
	res, err := egs.conn.Exec(q, userId, name)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return NotDeleted
	}
	return nil
}

func (egs EGS) Update(updated ExGroup) error {
	q := `UPDATE exercise_groups SET name=$1, user_id=$2 WHERE id=$3`
	res, err := egs.conn.Exec(q, updated.Name, updated.UserId, updated.Id)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return NotUpdated
	}
	return nil

}

func (egs EGS) UpdateByName(userId int64, name string, updated ExGroup) error {
	q := `UPDATE exercise_groups SET name=$1, user_id=$2 WHERE user_id=$3 AND name=$4`
	res, err := egs.conn.Exec(q, updated.Name, updated.UserId, userId, name)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return NotUpdated
	}
	return nil
}

func (egs EGS) FindByUserId(userId int64) ([]ExGroup, error) {
	var exGroups []ExGroup
	q := `SELECT * FROM exercise_groups WHERE user_id=$1`
	err := egs.conn.Select(&exGroups, q, userId)
	return exGroups, err
}
