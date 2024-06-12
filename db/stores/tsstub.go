package stores

import (
	"database/sql"
	"time"
)

type TrainingStoreStub struct {
	TrainingStore
}

func (tss TrainingStoreStub) StartTraining(userId int64) (int64, error) {
	return 12, nil
}

func (tss TrainingStoreStub) FinishTraining(userId int64) error {
	if userId == 1 {
		return AllTrainingsFinished
	}
	return nil
}
func (tss TrainingStoreStub) GetTrainings(userId int64) ([]Training, error) {
	var trainings []Training
	if userId == 1 {
		return trainings, sql.ErrNoRows
	}
	trainings = []Training{
		{
			Id:     1,
			UserId: userId,
			Begins: time.Now(),
			Finish: time.Now(),
		},
		{
			Id:     2,
			UserId: userId,
			Begins: time.Now(),
			Finish: time.Now(),
		},
		{
			Id:     3,
			UserId: userId,
			Begins: time.Now(),
			Finish: time.Now(),
		},
	}
	return trainings, nil
}
