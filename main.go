package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BrokerMessageHandler interface {
	Handle(message string)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	//Connection to rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//AMQP Channel for communication
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	ch1, err := conn.Channel()
	failOnError(err, "failed to open a channel 2")
	defer ch1.Close()
	err = ch.ExchangeDeclare(
		"ex_group",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	//declaring queue for listening
	q, err := ch.QueueDeclare(
		"create", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = ch.QueueBind(
		q.Name,
		"create",
		"ex_group",
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "error binding queue")
	}
	msgs1, err := ch1.Consume(
		"update",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	msgs, err := ch.Consume(
		"create", // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}
	go func() {
		for d := range msgs1 {
			log.Printf("received another queue: %s", d.Body)
		}
	}()
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
