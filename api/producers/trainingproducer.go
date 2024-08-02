package producers

import "github.com/rabbitmq/amqp091-go"

type TrainingProducer interface {
	StartTraining(amqp091.Delivery)
	FinishTraining(amqp091.Delivery)
	GetTraining(amqp091.Delivery)
}
