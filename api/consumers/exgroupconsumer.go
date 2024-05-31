package consumers

import (
	"log"

	rs "github.com/fridrock/rabbitsimplier"
)

type ExGroupRouter struct {
	rs.RConsumer
	rs.RProducer
}

func NewExGroupRouter(configurer rs.Configurer) *ExGroupRouter {
	var exGroupRouter ExGroupRouter
	exGroupRouter.RConsumer = rs.RConsumer{}
	err := exGroupRouter.RConsumer.CreateChannel(configurer.GetConnection())
	if err != nil {
		log.Fatal("error creating RConsumer for ExGroupRouter")
	}
	exGroupRouter.RProducer = rs.RProducer{}
	err = exGroupRouter.RProducer.CreateChannel(configurer.GetConnection())
	if err != nil {
		log.Fatal("error creating RProducer for ExGroupRouter")
	}
	return &exGroupRouter
}
