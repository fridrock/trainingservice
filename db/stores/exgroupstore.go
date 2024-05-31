package stores

import "github.com/jmoiron/sqlx"

// ExGroup struct that is entity for exercise_groups table
//
type ExGroup struct {
	Id     int64  `db:"id"`
	UserId int64  `db:"user_id"`
	Name   string `db:"name"`
}

// ExGroupStore - interface which contains all methods for working with exercise_groups table
//
type ExGroupStore interface {
	Save(ExGroup) (int64, error)
	FindById(int64) (ExGroup, error)
	FindName(string) (ExGroup, error)
	DeleteById(int64) error
	Update(ExGroup) error
	FindByUserId(int64) ([]ExGroup, error)
}

// EGS - standard realization of ExGroupInterface
//
type EGS struct {
	conn *sqlx.DB
}

// NewEGS - function that creates realization for ExGroupStore interface
//
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

func (egs EGS) FindByName(name string) (ExGroup, error) {
	var exGroup ExGroup
	q := `SELECT * FROM exercise_groups WHERE name=$1`
	err := egs.conn.Get(&exGroup, q, name)
	return exGroup, err
}

func (egs EGS) DeleteById(id int64) error {
	q := `DELETE FROM exercise_groups WHERE id=$1`
	_, err := egs.conn.Exec(q, id)
	return err
}

func (egs EGS) Update(updatedExGroup ExGroup) error {
	q := `UPDATE exercise_groups SET name=$1, user_id=$2 WHERE id=$3`
	_, err := egs.conn.Exec(q, updatedExGroup.Name, updatedExGroup.UserId, updatedExGroup.Id)
	return err
}

func (egs EGS) FindByUserId(userId int64) ([]ExGroup, error) {
	var exGroups []ExGroup
	q := `SELECT * FROM exercise_groups WHERE user_id=$1`
	err := egs.conn.Select(&exGroups, q, userId)
	return exGroups, err
}
