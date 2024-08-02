package routers

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
	ts     stores.TrainingStore
	routes map[string]func(amqp091.Delivery) string
}

const EXCHANGE_NAME = "sport_bot"

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
	tr.routes = make(map[string]func(amqp091.Delivery) string)
	tr.routes["start"] = tr.handleStart
	tr.routes["finish"] = tr.handleFinish
	tr.routes["get"] = tr.handleGet
	q, err := tr.RConsumer.CreateQueue()
	if err != nil {
		log.Fatal("error creating queue for exgroup consumer")
	}
	err = tr.RConsumer.SetBinding(q, "trainings.training.#", EXCHANGE_NAME)
	if err != nil {
		log.Fatal("error creating binding for exgroup consumer")
	}
	//creating dispatcher
	dispatcher := rs.NewRDispacher()
	for path, f := range tr.routes {
		dispatcher.RegisterHandler("trainings.training."+path, rs.NewHandlerFunc(func(msg amqp091.Delivery) {
			tr.sendResponse(f(msg), path)
		}))
	}
	tr.RConsumer.RegisterDispatcher(q, dispatcher)
}

func (tr *TrainingRouter) sendResponse(response, path string) {
	tr.RProducer.PublishMessage(
		context.Background(),
		EXCHANGE_NAME,
		"tgbot.training."+path,
		response)
}

// methods, that handles all messages and returns response
func (tr *TrainingRouter) handleStart(msg amqp091.Delivery) string {
	body := msg.Body
	userId, err := converters.ParseUserID(body)
	if err != nil {
		return "ERROR: wrong input"
	}
	slog.Info(fmt.Sprintf("request start training with user: %d", userId))
	trainingId, err := tr.ts.StartTraining(userId)
	if err != nil {
		return fmt.Sprintf("ERROR: error starting training: %v", err)
	}
	return fmt.Sprintf("SUCCESS: id:%d", trainingId)
}

func (tr *TrainingRouter) handleFinish(msg amqp091.Delivery) string {
	body := msg.Body
	userId, err := converters.ParseUserID(body)
	if err != nil {
		return "ERROR: wrong input"
	}
	slog.Info(fmt.Sprintf("request finish training with user: %d", userId))
	err = tr.ts.FinishTraining(userId)
	if err != nil {
		return fmt.Sprintf("ERROR: error finishing training: %v", err)
	}
	return "SUCCESS"
}

func (tr *TrainingRouter) handleGet(msg amqp091.Delivery) string {
	body := msg.Body
	userId, err := converters.ParseUserID(body)
	if err != nil {
		return "ERROR: wrong input"
	}
	slog.Info(fmt.Sprintf("request getting trainings with user: %d", userId))
	trainings, err := tr.ts.GetTrainings(userId)
	if err != nil {
		return fmt.Sprintf("ERROR: error getting trainings: %v", err)
	}
	r, err := json.MarshalIndent(trainings, "", "")
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	return fmt.Sprintf("SUCCESS: %v", string(r))
}

// Stop - Closure for closing channels of consumer and producer
func (tr TrainingRouter) Stop() {
	tr.RConsumer.Stop()
	tr.RProducer.Stop()
}
