package test

import (
	"context"
	"log"
	"strings"

	rs "github.com/fridrock/rabbitsimplier"
	"github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

var (
	rmqContainer     *rabbitmq.RabbitMQContainer
	clientConfigurer *rs.RConfigurer
	clientProducer   *rs.RProducer
	clientConsumer   *AllMessagesConsumer
)

type AllMessagesConsumer struct {
	LastMessageCh chan amqp091.Delivery
	rs.RConsumer
}

func setupContainer() {
	ctx := context.Background()
	rabbitmqContainer, err := rabbitmq.RunContainer(ctx,
		testcontainers.WithImage("rabbitmq:3.12.11-management-alpine"),
		rabbitmq.WithAdminUsername("guest"),
		rabbitmq.WithAdminPassword("guest"),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	rmqContainer = rabbitmqContainer
}

func GetRmqContainer() *rabbitmq.RabbitMQContainer {
	if rmqContainer == nil {
		setupContainer()
	}
	return rmqContainer
}

func setupConfigurer() {
	ctx := context.Background()
	container := GetRmqContainer()
	rconfigurer := &rs.RConfigurer{}
	url, err := container.AmqpURL(ctx)
	if err != nil {
		log.Fatalf("error getting rmqcontainer host for testing: %s", err)
	}
	parts := strings.Split(url, "@")
	err = rconfigurer.Configure(rs.Config{
		Username: "guest",
		Password: "guest",
		Host:     parts[1],
	})
	if err != nil {
		log.Fatalf("error creating configurer: %s", err)
	}
	clientConfigurer = rconfigurer
}

func GetClientConfigurer() *rs.RConfigurer {
	if clientConfigurer == nil {
		setupConfigurer()
	}
	return clientConfigurer
}

func setupProducer() {
	producer := &rs.RProducer{}
	producer.CreateChannel(GetClientConfigurer().GetConnection())
	producer.CreateExchange("sport_bot", "topic")
	clientProducer = producer
}

func GetClientProducer() *rs.RProducer {
	if clientProducer == nil {
		setupProducer()
	}
	return clientProducer
}

func setupConsumer() {
	var allmsgcons AllMessagesConsumer
	allmsgcons.LastMessageCh = make(chan amqp091.Delivery)
	allmsgcons.RConsumer = rs.RConsumer{}
	err := allmsgcons.RConsumer.CreateChannel(GetClientConfigurer().GetConnection())
	if err != nil {
		log.Fatalf("error creating connection for consumer in test: %s", err)
	}
	q, err := allmsgcons.RConsumer.CreateQueue()
	if err != nil {
		log.Fatalf("error creating consuming queue in test: %s", err)
	}
	err = allmsgcons.RConsumer.SetBinding(q, "tgbot.#", "sport_bot")
	if err != nil {
		log.Fatalf("error creating binding for queue in test: %s", err)
	}
	//write all messages to LastMessage field
	err = allmsgcons.RegisterDispatcher(q, rs.NewDispactherFunc(func(ch <-chan amqp091.Delivery) {
		for d := range ch {
			allmsgcons.LastMessageCh <- d
		}
	}))
	if err != nil {
		log.Fatalf("error registering dispatcher for consumer in test: %s", err)
	}
	clientConsumer = &allmsgcons
}

func GetClientConsumer() *AllMessagesConsumer {
	if clientConsumer == nil {
		setupConsumer()
	}
	return clientConsumer
}

func Stop() {
	clientConsumer.Stop()
	clientProducer.Stop()
	clientConfigurer.Stop()
	defer func() {
		if err := rmqContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("error stopping test container for rabbitmq: %s", err)
		}
	}()
}
