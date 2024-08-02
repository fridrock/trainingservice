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

// ExGroupRouter - structure, that contains both consumer, and producer for messaging inside ExGroup domain
type ExGroupRouter struct {
	rs.RConsumer
	rs.RProducer
	egs    stores.ExGroupStore
	routes map[string]func(amqp091.Delivery) string
}

// NewExGroupRouter - Default method for creation ExGroupRouter, requires rs.Configurer to create channels
// for consumer and producer
func NewExGroupRouter(configurer rs.Configurer) *ExGroupRouter {
	exGroupRouter := ExGroupRouter{}
	exGroupRouter.CreateConsumer(configurer)
	exGroupRouter.CreateProducer(configurer)
	exGroupRouter.SetEGS(stores.NewEGS(core.CreateConnection()))
	return &exGroupRouter
}

// CreateConsumer - helper method
func (egr *ExGroupRouter) CreateConsumer(configurer rs.Configurer) {
	egr.RConsumer = rs.RConsumer{}
	err := egr.RConsumer.CreateChannel(configurer.GetConnection())
	if err != nil {
		log.Fatal("error creating RConsumer for ExGroupRouter")
	}
}

// CreateProducer - helper method
func (egr *ExGroupRouter) CreateProducer(configurer rs.Configurer) {
	egr.RProducer = rs.RProducer{}
	err := egr.RProducer.CreateChannel(configurer.GetConnection())
	if err != nil {
		log.Fatal("error creating RProducer for ExGroupRouter")
	}
}

// SetEGS - Dependency injection of stores.ExGroupStore
func (egr *ExGroupRouter) SetEGS(egs stores.ExGroupStore) {
	egr.egs = egs
}

// Setup - main method, that sets up all routes and handlers for them
func (egr *ExGroupRouter) Setup() {
	egr.routes = make(map[string]func(amqp091.Delivery) string)
	egr.routes["create"] = egr.handleCreate
	egr.routes["delete"] = egr.handleDelete
	egr.routes["find"] = egr.handleFind
	egr.routes["update"] = egr.handleUpdate
	egr.routes["findByUser"] = egr.handleFindByUser
	q, err := egr.RConsumer.CreateQueue()
	if err != nil {
		log.Fatal("error creating queue for exgroup consumer")
	}
	err = egr.RConsumer.SetBinding(q, "trainings.exgroup.#", EXCHANGE_NAME)
	if err != nil {
		log.Fatal("error creating binding for exgroup consumer")
	}
	//creating dispatcher
	dispatcher := rs.NewRDispacher()
	for path, f := range egr.routes {
		dispatcher.RegisterHandler("trainings.exgroup."+path, rs.NewHandlerFunc(func(msg amqp091.Delivery) {
			egr.sendResponse(f(msg), path)
		}))
	}
	egr.RConsumer.RegisterDispatcher(q, dispatcher)
}

func (egr *ExGroupRouter) sendResponse(response, path string) {
	egr.RProducer.PublishMessage(
		context.Background(),
		EXCHANGE_NAME, "tgbot.exgroup."+path,
		response)
}

func (egr *ExGroupRouter) handleCreate(msg amqp091.Delivery) string {
	body := msg.Body
	exg, err := converters.FromJsonToExGroup(body)
	if err != nil {
		return "ERROR: wrong input"
	}
	slog.Info(fmt.Sprintf("request to create exgroup: %#v", exg))
	gotId, err := egr.egs.Save(exg)
	if err != nil {
		return fmt.Sprintf("ERROR: internal server error: %v", err)
	}
	return fmt.Sprintf("SUCCESS: id:%d", gotId)
}

func (egr *ExGroupRouter) handleDelete(msg amqp091.Delivery) string {
	body := msg.Body
	userId, name, err := converters.ParseExGroupProperties(body)
	if err != nil {
		return "ERROR: wrong input"
	}
	slog.Info(fmt.Sprintf("request to delete ex group with user_id: %d, name: %v", userId, name))
	err = egr.egs.DeleteByName(userId, name)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	return "SUCCESS"
}

func (egr *ExGroupRouter) handleFind(msg amqp091.Delivery) string {
	body := msg.Body
	userId, name, err := converters.ParseExGroupProperties(body)
	if err != nil {
		return "ERROR: wrong input"
	}
	slog.Info(fmt.Sprintf("request to find ex group with user_id: %d, name: %v", userId, name))
	exGroup, err := egr.egs.FindByName(userId, name)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	r, err := json.Marshal(&exGroup)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	return "SUCCESS: " + string(r)
}

func (egr *ExGroupRouter) handleFindByUser(msg amqp091.Delivery) string {
	body := msg.Body
	userId, err := converters.ParseUserID(body)
	if err != nil {
		return "ERROR: wrong input"
	}

	slog.Info(fmt.Sprintf("request to find by user_id: %d", userId))
	exGroups, err := egr.egs.FindByUserId(userId)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	response, err := json.MarshalIndent(exGroups, "", "")
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	return fmt.Sprintf("SUCCESS: %v", string(response))
}

func (egr *ExGroupRouter) handleUpdate(msg amqp091.Delivery) string {
	body := msg.Body
	updateExGroup, err := converters.ParseUpdateExGroup(body)
	if err != nil {
		return "ERROR: wrong input"
	}
	slog.Info(fmt.Sprintf(
		"request to update ex group with user_id: %d, name: %s, new_name: %s",
		updateExGroup.UserId,
		updateExGroup.Name,
		updateExGroup.NewName))
	err = egr.egs.UpdateByName(updateExGroup.UserId, updateExGroup.Name, updateExGroup.NewName)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	return "SUCCESS"
}

// Stop - Closure for closing channels of consumer and producer
func (egr ExGroupRouter) Stop() {
	egr.RConsumer.Stop()
	egr.RProducer.Stop()
}
