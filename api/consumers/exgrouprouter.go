package consumers

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/fridrock/auth_service/db/core"
	rs "github.com/fridrock/rabbitsimplier"
	"github.com/fridrock/trainingservice/api/utils/converters"
	"github.com/fridrock/trainingservice/db/stores"
	"github.com/rabbitmq/amqp091-go"
)

type ExGroupRouter struct {
	rs.RConsumer
	rs.RProducer
	egs stores.ExGroupStore
}

// each consumer have only one queue, but we can have binding with # symbols
func NewExGroupRouter(configurer rs.Configurer) *ExGroupRouter {
	exGroupRouter := ExGroupRouter{}
	exGroupRouter.CreateConsumer(configurer)
	exGroupRouter.CreateProducer(configurer)
	exGroupRouter.SetEGS(stores.NewEGS(core.CreateConnection()))
	exGroupRouter.egs = stores.NewEGS(core.CreateConnection())
	return &exGroupRouter
}

func (egr *ExGroupRouter) CreateConsumer(configurer rs.Configurer) {
	egr.RConsumer = rs.RConsumer{}
	err := egr.RConsumer.CreateChannel(configurer.GetConnection())
	if err != nil {
		log.Fatal("error creating RConsumer for ExGroupRouter")
	}
}

func (egr *ExGroupRouter) CreateProducer(configurer rs.Configurer) {
	egr.RProducer = rs.RProducer{}
	err := egr.RProducer.CreateChannel(configurer.GetConnection())
	if err != nil {
		log.Fatal("error creating RProducer for ExGroupRouter")
	}
}

func (egr *ExGroupRouter) SetEGS(egs stores.ExGroupStore) {
	egr.egs = egs
}

func (egr *ExGroupRouter) Setup() {
	q, err := egr.RConsumer.CreateQueue()
	if err != nil {
		log.Fatal("error creating queue for exgroup consumer")
	}
	err = egr.RConsumer.SetBinding(q, "trainings.exgroup.#", "sport_bot")
	if err != nil {
		log.Fatal("error creating binding for exgroup consumer")
	}
	//creating dispatcher
	dispatcher := rs.NewRDispacher()
	//handler for creating exgroup
	dispatcher.RegisterHandler("trainings.exgroup.create", rs.NewHandlerFunc(func(msg amqp091.Delivery) {
		body := msg.Body
		exg, err := converters.FromJsonToExGroup(body)
		if err != nil {
			egr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.exgroup.create", "ERROR: wrong input")
			return
		}
		slog.Info(fmt.Sprintf("request to create exgroup: %#v", exg))
		gotId, err := egr.egs.Save(exg)
		if err != nil {
			egr.RProducer.PublishMessage(context.Background(),
				"sport_bot",
				"tgbot.exgroup.create",
				fmt.Sprintf("ERROR: internal server error: %v", err))
			return
		}
		egr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.exgroup.create", fmt.Sprintf("SUCCESS: id:%d", gotId))
	}))

	//handler for deleting exgroup
	dispatcher.RegisterHandler("trainings.exgroup.delete", rs.NewHandlerFunc(func(msg amqp091.Delivery) {
		body := msg.Body
		userId, name, err := converters.ParseDeleteRequest(body)
		if err != nil {
			egr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.exgroup.delete", "ERROR: wrong input")
			return
		}
		slog.Info(fmt.Sprintf("request to delete ex group with user_id: %d, name: %v", userId, name))
		err = egr.egs.DeleteByName(userId, name)
		if err != nil {
			egr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.exgroup.delete", fmt.Sprintf("ERROR: %v", err))
			return
		}
		egr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.exgroup.delete", "SUCCESS")
	}))

	//handler for finding exgroup by id
	dispatcher.RegisterHandler("trainings.exgroup.find", rs.NewHandlerFunc(func(msg amqp091.Delivery) {
		body := msg.Body
		userId, name, err := converters.ParseDeleteRequest(body)
		if err != nil {
			egr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.exgroup.delete", "ERROR: wrong input")
			return
		}
		slog.Info(fmt.Sprintf("request to delete ex group with user_id: %d, name: %v", userId, name))
		err = egr.egs.DeleteByName(userId, name)
		if err != nil {
			egr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.exgroup.delete", fmt.Sprintf("ERROR: %v", err))
			return
		}
		egr.RProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.exgroup.delete", "SUCCESS")

	}))

	egr.RConsumer.RegisterDispatcher(q, dispatcher)

}
func (egr ExGroupRouter) Stop() {
	egr.RConsumer.Stop()
	egr.RProducer.Stop()
}
