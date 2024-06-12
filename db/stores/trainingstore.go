package stores

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type Training struct {
	Id     int64     `db:"id" json:"id"`
	UserId int64     `db:"user_id" json:"user_id"`
	Begins time.Time `db:"started" json:"begins"`
	Finish time.Time `db:"finished" json:"finish"`
}

type TrainingStore interface {
	StartTraining(userId int64) (int64, error)
	FinishTraining(userId int64) error
	FindById(trainingId int64) (Training, error)
	GetLastTraining(userId int64) (Training, error)
	GetTrainings(userId int64) ([]Training, error)
}

var (
	AllTrainingsFinished = errors.New("Empty non-finished trainings list")
)

type TS struct {
	conn *sqlx.DB
}

func NewTs(conn *sqlx.DB) *TS {
	return &TS{
		conn: conn,
	}
}

func (ts TS) StartTraining(userId int64) (id int64, err error) {
	training := Training{
		UserId: userId,
		Begins: time.Now(),
	}
	q := "INSERT INTO trainings(user_id, begins) VALUES ($1, $2) RETURNING id"
	err = ts.conn.Get(&id, q, training.UserId, training.Begins)
	return id, err
}

func (ts TS) FinishTraining(userId int64) error {
	q := "UPDATE trainings SET finish=$1 WHERE finish IS NULL AND user_id=$2"
	res, err := ts.conn.Exec(q, time.Now(), userId)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return AllTrainingsFinished
	}
	return nil
}

func (ts TS) FindById(trainingId int64) (Training, error) {
	q := "SELECT * FROM trainings WHERE id=$1"
	var training Training
	err := ts.conn.Get(&training, q, trainingId)
	return training, err
}

func (ts TS) GetLastTraining(userId int64) (Training, error) {
	q := "SELECT * FROM trainings WHERE user_id=$1 ORDERED BY begins DESC LIMIT 1"
	var training Training
	err := ts.conn.Get(&training, q, userId)
	return training, err
}

func (ts TS) GetTrainings(userId int64) ([]Training, error) {
	var trainings []Training
	q := "SELECT * FROM trainings WHERE user_id=$1"
	err := ts.conn.Select(&trainings, q, userId)
	return trainings, err
}
