package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/fridrock/trainingservice/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

// defining answer event
type ExGroupCreatedEvent struct {
	Event string `json:"event"`
	Text  string `json:"text"`
}

func setupConfigurer() *broker.BrokerConfigurerImpl {
	brc := &broker.BrokerConfigurerImpl{}
	err := brc.Configure()
	if err != nil {
		log.Fatal("error creating connection to rabbitmq")
	}
	return brc
}

func setupConsumer(brc *broker.BrokerConfigurerImpl, brProducer *broker.BrokerProducerImpl) *broker.BrokerConsumerImpl {
	brConsumer := &broker.BrokerConsumerImpl{}
	err := brConsumer.CreateChannel(brc.Connection)
	if err != nil {
		log.Fatal("error creating channel for consumer")
	}
	q, err := brConsumer.CreateQueue()
	if err != nil {
		log.Fatal("error creating queue")
	}
	brConsumer.SetBinding(q, "trainings.exgroup", "sport_bot")
	brConsumer.RegisterConsumer(q, func(msgs <-chan amqp.Delivery) {
		for d := range msgs {
			slog.Info(fmt.Sprintf("Received a message: %s", d.Body))
			answer := ExGroupCreatedEvent{
				Event: "created",
				Text:  "ex group created successfully",
			}
			res, err := json.Marshal(&answer)
			if err != nil {
				slog.Error("error marshalling answer")
			}
			err = brProducer.PublishMessage(context.Background(), "sport_bot", "tgbot.exgroup", string(res))
			if err != nil {
				slog.Error(err.Error())
			}
		}
	})
	return brConsumer
}
func setupProducer(brc *broker.BrokerConfigurerImpl) *broker.BrokerProducerImpl {
	brProducer := &broker.BrokerProducerImpl{}
	err := brProducer.CreateChannel(brc.Connection)
	if err != nil {
		log.Fatal("error creating channel for producer")
	}
	return brProducer
}
func main() {
	brc := setupConfigurer()
	//setting up configurer
	defer brc.Stop()
	brProducer := setupProducer(brc)
	defer brProducer.Stop()
	brConsumer := setupConsumer(brc, brProducer)
	defer brConsumer.Stop()
	//infinite work of service
	var forever chan struct{}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
