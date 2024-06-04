package main

import (
	"log"
	"log/slog"

	rs "github.com/fridrock/rabbitsimplier"
	"github.com/fridrock/trainingservice/api/consumers"
	amqp "github.com/rabbitmq/amqp091-go"
)

// defining answer event
type ExGroupCreatedEvent struct {
	Event string `json:"event"`
	Text  string `json:"text"`
}

func setupConfigurer() *rs.RConfigurer {
	brc := &rs.RConfigurer{}
	err := brc.Configure(rs.Config{
		Username: "guest",
		Password: "guest",
		Host:     "localhost",
	})
	if err != nil {
		log.Fatal("error creating connection to rabbitmq")
	}
	return brc
}

func setupConsumer(configurer rs.Configurer, producer rs.Producer) *rs.RConsumer {
	consumer := &rs.RConsumer{}
	err := consumer.CreateChannel(configurer.GetConnection())
	if err != nil {
		log.Fatal("error creating channel for consumer")
	}
	q, err := consumer.CreateQueue()
	if err != nil {
		log.Fatal("error creating queue")
	}
	consumer.SetBinding(q, "trainings.exgroup.#", "sport_bot")

	dispatcher := rs.NewRDispacher()
	dispatcher.RegisterHandler("trainings.exgroup.create", rs.NewHandlerFunc(func(msg amqp.Delivery) {
		slog.Info("trainings.exgroup.create event")
	}))
	consumer.RegisterDispatcher(q, &dispatcher)
	return consumer
}

func setupProducer(brc rs.Configurer) *rs.RProducer {
	brProducer := &rs.RProducer{}
	err := brProducer.CreateChannel(brc.GetConnection())
	if err != nil {
		log.Fatal("error creating channel for producer")
	}
	return brProducer
}
func main() {
	brc := setupConfigurer()
	//setting up configurer
	defer brc.Stop()
	router := consumers.NewExGroupRouter(brc)
	router.Setup()
	defer router.Stop()
	// brProducer := setupProducer(brc)
	// defer brProducer.Stop()
	// brConsumer := setupConsumer(brc, brProducer)
	// defer brConsumer.Stop()
	// //infinite work of service
	var forever chan struct{}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
