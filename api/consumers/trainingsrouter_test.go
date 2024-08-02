package consumers

import (
	"context"
	"strings"
	"testing"
)

func TestStartTraining(t *testing.T) {
	data := []struct {
		testName       string
		message        string
		expectedResult string
		errMessage     string
	}{
		{
			"Negative case: wrong input",
			`{"userid":2}`,
			wrongInput,
			"Error starting training, received: %v",
		},
		{
			"Positive case",
			`{"user_id":1}`,
			"SUCCESS: id:12",
			"Error starting training, received: %v",
		},
	}
	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			clientProducer.PublishMessage(context.Background(), "sport_bot", "trainings.training.start", d.message)
		})
		body := <-clientConsumer.LastMessageCh
		if body.RoutingKey != "tgbot.training.start" {
			t.Errorf("error wrong result routing key")
		}
		received := string(body.Body)
		if received != d.expectedResult {
			t.Errorf(d.errMessage, received)
		}
	}
}

func TestFinishTraining(t *testing.T) {
	data := []struct {
		testName       string
		message        string
		expectedResult string
		errMessage     string
	}{
		{
			"Negative case: wrong input",
			`{"userid":2}`,
			wrongInput,
			"Error finishing training, received: %v",
		},

		{
			"Negative case: all trainings finished",
			`{"user_id":1}`,
			"ERROR: error finishing training: Empty non-finished trainings list",
			"Error finishing training, received: %v",
		}, {
			"Positive case",
			`{"user_id":2}`,
			"SUCCESS",
			"Error finishing training, received: %v",
		},
	}
	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			clientProducer.PublishMessage(context.Background(), "sport_bot", "trainings.training.finish", d.message)
		})
		body := <-clientConsumer.LastMessageCh
		if body.RoutingKey != "tgbot.training.finish" {
			t.Errorf("error wrong result routing key")
		}
		received := string(body.Body)
		if received != d.expectedResult {
			t.Errorf(d.errMessage, received)
		}
	}
}

func TestGetTrainings(t *testing.T) {
	data := []struct {
		testName       string
		message        string
		expectedResult string
		errMessage     string
	}{
		{
			"Negative case: wrong input",
			`{"userid":2}`,
			wrongInput,
			"Error getting trainings, received: %v",
		},

		{
			"Negative case: no trainings with such user_id",
			`{"user_id":1}`,
			"ERROR: error getting trainings: sql: no rows in result set",
			"Error getting trainings, received: %v",
		},
		{
			"Positive case",
			`{"user_id":2}`,
			"SUCCESS",
			"Error finishing training, received: %v",
		},
	}
	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			clientProducer.PublishMessage(context.Background(), "sport_bot", "trainings.training.get", d.message)
		})
		body := <-clientConsumer.LastMessageCh
		if body.RoutingKey != "tgbot.training.get" {
			t.Errorf("error wrong result routing key")
		}
		received := string(body.Body)
		if d.testName != "Positive case" && received != d.expectedResult {
			t.Errorf(d.errMessage, received)
		}
		if d.testName == "Positive case" {
			parts := strings.Split(received, ":")
			if parts[0] != "SUCCESS" {
				t.Error(d.errMessage, received)
			}
		}
	}
}
