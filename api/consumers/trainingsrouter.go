package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	rs "github.com/fridrock/rabbitsimplier"
	"github.com/fridrock/trainingservice/api/utils/converters"
	"github.com/fridrock/trainingservice/db/core"
	"github.com/fridrock/trainingservice/db/stores"
	"github.com/rabbitmq/amqp091-go"
)

// TrainingRouter- structure, that contains both consumer, and producer for messaging inside Training domain
type TrainingRouter struct {
	rs.RConsumer
	rs.RProducer
	ts stores.TrainingStore
}

// NewTraningRouter - Default method for creation TrainingRouter, requires rs.Configurer to create channels
// for consumer and producer
func NewTrainingRouter(configurer rs.Configurer) *TrainingRouter {
	trainingRouter := TrainingRouter{}
	trainingRouter.CreateConsumer(configurer)
	trainingRouter.CreateProducer(configurer)
	trainingRouter.SetTS(stores.NewTs(core.CreateConnection()))
	return &trainingRouter
}

// CreateConsumer - helper method
func (tr *TrainingRouter) CreateConsumer(configurer rs.Configurer) {
	tr.RConsumer = rs.RConsumer{}
	err := tr.RConsumer.CreateChannel(configurer.GetConnection())
	if err != nil {
		log.Fatal("error creating RConsumer for TrainingRouter")
	}
}

// CreateProducer - helper method
func (tr *TrainingRouter) CreateProducer(configurer rs.Configurer) {
	tr.RProducer = rs.RProducer{}
	err := tr.RProducer.CreateChannel(configurer.GetConnection())
	if err != nil {
		log.Fatal("error creating RProducer for TrainingRouter")
	}
}

// SetEGS - Dependency injection of stores.ExGroupStore
func (tr *TrainingRouter) SetTS(ts stores.TrainingStore) {
	tr.ts = ts
}

// Setup - main method, that sets up all routes and handlers for them
func (tr *TrainingRouter) Setup() {
	q, err := tr.RConsumer.CreateQueue()
	if err != nil {
		log.Fatal("error creating queue for exgroup consumer")
	}
	err = tr.RConsumer.SetBinding(q, "trainings.training.#", "sport_bot")
	if err != nil {
		log.Fatal("error creating binding for exgroup consumer")
	}
	//creating dispatcher
	dispatcher := rs.NewRDispacher()
	//handler for starting training
	dispatcher.RegisterHandler("trainings.training.start", rs.NewHandlerFunc(func(msg amqp091.Delivery) {
		body := msg.Body
		userId, err := converters.ParseUserID(body)
		if err != nil {
			tr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.training.start", "ERROR: wrong input")
			return
		}
		slog.Info(fmt.Sprintf("request start training with user: %d", userId))
		trainingId, err := tr.ts.StartTraining(userId)
		if err != nil {
			tr.RProducer.PublishMessage(context.Background(),
				"sport_bot",
				"tgbot.training.start",
				fmt.Sprintf("ERROR: error starting training: %v", err))
			return
		}
		tr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.training.start", fmt.Sprintf("SUCCESS: id:%d", trainingId))
	}))

	//handler for finishing training
	dispatcher.RegisterHandler("trainings.training.finish", rs.NewHandlerFunc(func(msg amqp091.Delivery) {
		body := msg.Body
		userId, err := converters.ParseUserID(body)
		if err != nil {
			tr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.training.finish", "ERROR: wrong input")
			return
		}
		slog.Info(fmt.Sprintf("request finish training with user: %d", userId))
		err = tr.ts.FinishTraining(userId)
		if err != nil {
			tr.RProducer.PublishMessage(context.Background(),
				"sport_bot",
				"tgbot.training.finish",
				fmt.Sprintf("ERROR: error finishing training: %v", err))
			return
		}
		tr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.training.finish", "SUCCESS")
	}))

	//handler for getting all trainings
	dispatcher.RegisterHandler("trainings.training.get", rs.NewHandlerFunc(func(msg amqp091.Delivery) {
		body := msg.Body
		userId, err := converters.ParseUserID(body)
		if err != nil {
			tr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.training.get", "ERROR: wrong input")
			return
		}
		slog.Info(fmt.Sprintf("request getting trainings with user: %d", userId))
		trainings, err := tr.ts.GetTrainings(userId)
		if err != nil {
			tr.RProducer.PublishMessage(context.Background(),
				"sport_bot",
				"tgbot.training.get",
				fmt.Sprintf("ERROR: error getting trainings: %v", err))
			return
		}
		r, err := json.MarshalIndent(trainings, "", "")
		if err != nil {
			tr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.training.get", fmt.Sprintf("ERROR: %v", err))
			return
		}
		tr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.training.get", fmt.Sprintf("SUCCESS: %v", string(r)))
	}))
	tr.RConsumer.RegisterDispatcher(q, dispatcher)
}

// Stop - Closure for closing channels of consumer and producer
func (tr TrainingRouter) Stop() {
	tr.RConsumer.Stop()
	tr.RProducer.Stop()
}
