package stores

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTSStartTrainingFindById(t *testing.T) {
	id, err := ts.StartTraining(1)
	if err != nil {
		t.Fatalf("error starting training: %v", err)
	}
	var empty Training
	training, err := ts.FindById(id)
	if err != nil {
		t.Fatalf("error getting added training:%v", err)
	}
	if diff := cmp.Diff(empty, training); diff == "" {
		t.Errorf("got wrong training value:%v", training)
	}
	if training.Begins != training.Finish {
		t.Errorf("finish value isn't equal to begins value")
	}
	t.Cleanup(clearTables)
}

func TestTSFinishTraining(t *testing.T) {
	id, _ := ts.StartTraining(1)
	err := ts.FinishTraining(1)
	if err != nil {
		t.Fatalf("error finishing training: %v", err)
	}
	training, _ := ts.FindById(id)
	if training.Begins == training.Finish {
		t.Errorf("finish value didn't change")
	}
	t.Cleanup(clearTables)
}

func TestTSGetLastTraining(t *testing.T) {
	ts.StartTraining(1)
	ts.StartTraining(1)
	lastId, _ := ts.StartTraining(1)
	lastTraining, err := ts.GetLastTraining(1)
	if err != nil {
		t.Fatalf("error getting last training: %v", err)
	}
	if lastTraining.Id != lastId {
		t.Errorf("getting non-last training")
	}
	t.Cleanup(clearTables)
}

func TestTSGetTrainings(t *testing.T) {
	for i := 0; i < 10; i++ {
		ts.StartTraining(1)
		ts.FinishTraining(1)
	}
	trainings, err := ts.GetTrainings(1)
	if err != nil {
		t.Errorf("error getting trainings of user with id: %d, error: %v", 1, err)
	}
	if len(trainings) != 10 {
		t.Errorf("Getting wrong amount of trainings")
	}
}
